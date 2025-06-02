package controller

import (
	"errors"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"hilive/modules/utils"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

var (
	userLock sync.Mutex
)

// GameVoteAvatarModel 資料表欄位
type GameVoteAvatarModel struct {
	OptionID string `json:"option_id" example:"option_id"`
	Avatar   string `json:"avatar" example:"avatar"`
}

// GameVoteRecordModel 資料表欄位
type GameVoteRecordModel struct {
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	OptionID   string `json:"option_id" example:"option_id"`
	Round      int64  `json:"round" example:"0"`
	Score      int64  `json:"score" example:"0"`
}

// GameVoteOptionModel 資料表欄位
type GameVoteOptionModel struct {
	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id" example:"activity_id"`
	GameID          string `json:"game_id" example:"game_id"`
	OptionID        string `json:"option_id" example:"option_id"`
	OptionName      string `json:"option_name" example:"option_name"`
	OptionPicture   string `json:"option_picture" example:"option_picture"`
	OptionIntroduce string `json:"option_introduce" example:"option_introduce"`
	OptionScore     int64  `json:"option_score" example:"1"`

	// 活動
	UserID string `json:"user_id" example:"user_id"`
}

// UserModel 資料表欄位
type UserModel struct {
	UserID string `json:"user_id" example:"1"`
	Name   string `json:"name" example:"User Name"`
	Avatar string `json:"avatar" example:"https://..."`
	Role   string `json:"role" example:"host、guest"`

	// 隊伍資訊
	Team     string `json:"team" example:"left_team.right_team"` // 玩家隊伍
	Action   string `json:"action" example:"join.exit"`          // 玩家動作
	IsLeader bool   `json:"is_leader" example:"true"`            // 玩家是否為隊長
	Score    int64  `json:"score" example:"1"`                   // 玩家分數

	// 賓果遊戲
	// Numbers []string `json:"numbers" example:"1"` // 玩家號碼排序
	// UserBingoRound int64    `json:"user_bingo_round" example:"1"` // 完成賓果的回合數
}

// UserParam 判斷用戶角色API請求參數
type UserParam struct {
	User        UserModel
	Game        GameModel
	PrizeStaffs []models.PrizeStaffModel // 用戶中獎紀錄
	// Score       int64                    // 用戶分數
	Times int64 `json:"times" example:"1"` // 用戶可玩遊戲次數
}

// @Summary 判斷用戶角色、是否為黑名單，回傳用戶、遊戲輪次、中獎記錄資訊
// @Tags User
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string false "game ID"
// @Param session_id query string true "session ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, lottery, monopoly, QA, tugofwar, bingo, signname, 3DGachaMachine, vote)
// @Param device query string false "device"
// @Param login query string false "是否需要登入才能執行遊戲頁面" Enums(true, false)
// @Success 200 {array} UserParam
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /user [get]
func (h *Handler) GetUserRole(ctx *gin.Context) {
	var (
		host                                 = ctx.Request.Host
		activityID                           = ctx.Query("activity_id")
		gameID                               = ctx.Query("game_id")
		sessionID                            = ctx.Query("session_id")
		game                                 = ctx.Query("game")
		login                                = ctx.Query("login")
		device                               = ctx.Query("device")
		signnameModel                        models.SignnameSettingModel
		gameModel                            models.GameModel
		winRecords                           = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄
		user                                 models.LineModel
		isExist                              bool
		amount, times                        int64
		winPrizeAmount, losePrizeAmount      int64 // 拔河遊戲勝敗方剩餘獎品數
		role, hostUserID, gachaMachineQRcode string
		voteOptions                          = make([]GameVoteOptionModel, 0) // 投票選項資訊
		voteRecords                          = make([]GameVoteRecordModel, 0) // 玩家投票紀錄
		voteTimes                            int64                            // 剩餘投票次數
		voteScore                            int64                            // 玩家投票權重
		voteMethodPlayer                     int64                            // 玩家投票方式，0代表自由多選，不為0代表一個選項n票
		err                                  error
	)
	// 清除redis(晚點刪除)
	// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊
	// log.Println("判斷角色: ", activityID, game)

	// 扭蛋機遊戲需上鎖，避免玩家同時進入遊戲頁面
	if game == "3DGachaMachine" {
		userLock.Lock()
		defer userLock.Unlock()
	}

	if host != config.API_URL {
		response.BadRequestWithRole(ctx, "guest", h.dbConn, models.EditErrorLogModel{
			UserID:  utils.ClientIP(ctx.Request),
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if game == "" {
		response.BadRequestWithRole(ctx, "guest", h.dbConn, models.EditErrorLogModel{
			UserID:  utils.ClientIP(ctx.Request),
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識遊戲資訊，請輸入有效的遊戲參數",
		})
		return
	}

	// 簽名牆沒有game_id參數，其他遊戲都有
	if gameID != "" {
		gameModel = h.getGameInfo(gameID, game) // 遊戲資訊(redis處理)
		if gameModel.ID == 0 {
			response.BadRequestWithRole(ctx, "guest", h.dbConn, models.EditErrorLogModel{
				UserID:  utils.ClientIP(ctx.Request),
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識遊戲資訊，請輸入有效的遊戲ID",
			})
			return
		}

		hostUserID = gameModel.UserID
	}

	// 簽名牆
	if game == "signname" {
		// 取得活動主持人user_id資訊
		activityModel, err := h.getActivityInfo(true, activityID) // 活動資訊(redis處理)
		if err != nil {
			// logger.LoggerError(ctx, "錯誤: 無法辨識遊戲資訊，請輸入有效的遊戲ID")
			response.BadRequestWithRole(ctx, "guest", h.dbConn, models.EditErrorLogModel{
				UserID:  utils.ClientIP(ctx.Request),
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}

		// 活動主持人ID
		hostUserID = activityModel.UserID

		// 簽名牆設置資料(redis)
		signnameModel, err = h.getSignnameSetting(true, activityID)
		if err != nil {
			response.BadRequestWithRole(ctx, "guest", h.dbConn, models.EditErrorLogModel{
				UserID:  utils.ClientIP(ctx.Request),
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	}

	if device == "terminal" {
		// 不需要登入就可以執行簽名牆頁面
		role = "host" // 主持人
	} else if login != "false" {
		// 利用session值解碼取得用戶資料(沒有使用redis、資料表)
		if user, role, err = h.getUser(sessionID, hostUserID); err != nil {
			response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
				UserID:  utils.ClientIP(ctx.Request),
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	} else if login == "false" {
		// 不需要登入就可以執行遊戲頁面
		role = "host" // 主持人
	}

	if role == "guest" {
		// 優化後，使用SET格式處理，判斷玩家是否簽到
		isExist = h.IsSign(config.SIGN_STAFFS_2_REDIS, activityID, user.UserID)

		// redis中無簽到資料，查詢資料表中是否有簽到資料
		if !isExist {
			var staffs = make([]models.ApplysignModel, 0)

			// 所有人員資訊(包含簽到、未簽到、審核中...人員)
			staffs, _ = h.getAllStaffs(activityID, "", "", "", 0, 0)
			for _, staff := range staffs {
				if staff.UserID == user.UserID && staff.Status == "sign" {
					isExist = true
				}
			}
		}

		if !isExist {
			// log.Println("玩家未完成簽到，自動導向報名簽到頁面")

			// 報名簽到頁面
			// redirect := fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, config.HILIVES_NET_URL, activityID, gameID) + "&is_exist=false&sign=open"

			// ctx.Redirect(http.StatusFound, redirect)

			// logger.LoggerError(ctx, "錯誤: 玩家未完成報名簽到，請重新報名簽到")
			response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
				UserID:  user.UserID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 玩家未完成報名簽到，請重新報名簽到",
			})
			return
		}

		// 遊戲類型才需要取得中獎名單
		if gameID != "" {
			if game == "vote" {
				// log.Println("投票遊戲，判斷黑名單")

				// 判斷是否為活動黑名單、遊戲黑名單(redis處理)
				if h.IsBlackStaff(activityID, gameID, game, user.UserID) {
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 用戶為黑名單人員(活動、遊戲黑名單)，無法參加遊戲。如有疑問，請聯繫主辦單位",
					})
					return
				}

				// 判斷遊戲狀態
				var gameStatus string

				// 判斷投票結束時間
				now, _ := time.ParseInLocation("2006-01-02 15:04:05",
					time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local) // 目前時間
				startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteStartTime, time.Local)
				endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteEndTime, time.Local)

				// 比較時間，判斷遊戲狀態
				if now.Before(startTime) { // 現在時間在開始時間之前
					gameStatus = "close" // 關閉
				} else if now.Before(endTime) { // 現在時間在截止時間之前
					gameStatus = "gaming" // 遊戲中
				} else { // 現在時間在截止時間之後
					gameStatus = "calculate" // 結算狀態
				}

				if gameStatus != gameModel.GameStatus {
					// 狀態不一樣，更新為最新狀態
					gameModel.GameStatus = gameStatus

					// 更新遊戲狀態(mysql)
					// if err := db.Table(config.ACTIVITY_GAME_TABLE).WithConn(h.dbConn).
					// 	Where("game_id", "=", gameModel.GameID).
					// 	Update(command.Value{"game_status": gameStatus}); err != nil &&
					// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					// 	response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					// 		UserID:  user.UserID,
					// 		Method:  ctx.Request.Method,
					// 		Path:    ctx.Request.URL.Path,
					// 		Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
					// 	})
					// 	return
					// }

					// 更新遊戲狀態(mongo，activity_game)
					filter := bson.M{"game_id": gameModel.GameID} // 過濾條件
					// 要更新的值
					update := bson.M{
						"$set": bson.M{"game_status": gameStatus},
						// "$unset": bson.M{},                // 移除不需要的欄位
						// "$inc":   bson.M{"edit_times": 1}, // edit 欄位遞增 1
					}

					if _, err := h.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
						response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  "",
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
						})
					}

					// 修改redis中的遊戲資訊
					h.redisConn.HashSetCache(config.GAME_REDIS+gameID, "game_status", gameStatus)

					// 設置過期時間
					// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

					// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
					h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "遊戲狀態改變")

					// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
					h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
				}

				if gameModel.GameStatus == "close" { // 關閉
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 該投票遊戲場次尚未開始，請等待遊戲開始",
					})
					return
				} else if gameModel.GameStatus == "calculate" { // 關閉
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 該投票遊戲場次已結算，請等待遊戲重新開始",
					})
					return
				}

				// 玩家投票權重處理
				var isSpecialOfficer bool
				voteScore = 1                   // 玩家權重分數
				voteTimes = gameModel.VoteTimes // 玩家投票次數

				// 判斷投票方式
				if gameModel.VoteMethodPlayer == "one_vote" {
					voteMethodPlayer = 1
				} else if gameModel.VoteMethodPlayer == "free_vote" {
					voteMethodPlayer = 0
				}

				// 取得所有特殊人員資料
				officers, err := models.DefaultGameVoteSpecialOfficerModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(true, gameID, user.UserID)
				if err != nil {
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: err.Error(),
					})
					return
				}

				// 查詢玩家是否為特殊人員
				if len(officers) > 0 {
					//特殊人員
					isSpecialOfficer = true

					voteScore = officers[0].Score                   // 玩家投票權重
					voteTimes = officers[0].Times                   // 玩家投票次數
					voteMethodPlayer = officers[0].VoteMethodPlayer // 玩家投票方式
				}
				// for _, office := range officers {
				// }

				// log.Println("玩家權重分數: ", voteScore)

				// 判斷投票限制
				if gameModel.VoteRestriction == "special_officer" {
					// 特殊人員才可以投票
					if !isSpecialOfficer {
						response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
							UserID:  user.UserID,
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 該投票遊戲場次只有特殊人員能夠參與",
						})
						return
					}
				}

				// 投票遊戲取得玩家投票紀錄
				voteRecordModels, err := models.DefaultGameVoteRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindUserRecord(true, gameID, user.UserID)
				if err != nil {
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: err.Error(),
					})
					return
				}

				for _, recordModel := range voteRecordModels {
					voteRecords = append(voteRecords, GameVoteRecordModel{
						ID:         recordModel.ID,
						ActivityID: recordModel.ActivityID,
						GameID:     recordModel.GameID,
						UserID:     recordModel.UserID,
						OptionID:   recordModel.OptionID,
						Round:      recordModel.Round,
						Score:      recordModel.Score,
					})
				}
				// log.Println("投票次數: ", voteTimes, int64(len(voteRecords)))
				// 計算玩家投票剩餘次數
				voteTimes = voteTimes - int64(len(voteRecords))

				if voteTimes <= 0 {
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 玩家剩餘投票次數為0，請等待遊戲結束",
					})
					return
				}

			} else {
				// log.Println("投票遊戲之外的遊戲，判斷黑名單")
				// 投票遊戲之外的遊戲
				// 處理玩家端判斷(遊戲狀態、黑名單、用戶遊戲紀錄、是否允許重複中獎、遊戲抽獎可玩次數)
				winRecords, times, err = h.getUserGuest(activityID, gameID, user.UserID, game, gameModel)
				if err != nil {
					response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
						UserID:  user.UserID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: err.Error(),
					})
					return
				}

				// log.Println("中獎紀錄: ", len(winRecords))
			}
		}

		// 簽名牆
		if game == "signname" {
			// log.Println("玩家簽名牆，判斷黑名單")

			// 判斷是否為活動黑名單、簽名黑名單(redis處理)
			if h.IsBlackStaff(activityID, gameID, game, user.UserID) {
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 用戶為黑名單人員(活動、簽名牆黑名單)，無法參加遊戲。如有疑問，請聯繫主辦單位",
				})
				return
			}

			// 簽名模式必須為手機簽名模式
			if signnameModel.SignnameMode == "terminal" {
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 簽名牆為終端機簽名模式，用戶無法進入頁面",
				})
				return
			}

			// 判斷黑名單

			// 活動下所有用戶簽名牆資料
			signnames, err := h.getSignnameDatas(true, activityID, 0, 0)
			if err != nil {
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: err.Error(),
				})
				return
			}

			// 判斷是否限制簽名次數
			if signnameModel.SignnameLimitTimes == "open" {
				// fmt.Println("限制用戶簽名次數")

				// 用戶可簽名次數
				times = signnameModel.SignnameTimes

				for i := 0; i < len(signnames); i++ {
					if signnames[i].UserID == user.UserID {
						// 該筆資料為用戶的簽名資料，可簽名次數-1
						times--
					}
				}
				// 簽名次數小於零，修改為0
				if times < 0 {
					times = 0
				}
			} else {
				// fmt.Println("沒有限制用戶簽名次數")
			}
		}
	} else if role == "host" {

		// 判斷為主持人時，取得所有搖號抽獎場次中獎人員資料，確保寫入redis中
		h.getDrawNumbersAllWinningStaffs(true, activityID)

		// 遊戲類型才需要取得獎品數量
		if gameID != "" {
			// 處理主持端判斷(獎品數量、遊戲狀態)
			amount, winPrizeAmount, losePrizeAmount, err = h.getUserHost(gameID, game, gameModel)
			if err != nil {
				// logger.LoggerError(ctx, err.Error())
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: err.Error(),
				})
				return
			}
		}

		// 簽名牆
		if game == "signname" {
			// fmt.Println(signnameModel.SignnameMode, device)
			// 簽名模式必須為終端機簽名模式
			if signnameModel.SignnameMode == "phone" && device == "terminal" {
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 簽名牆為手機簽名模式，主持人無法進入終端機簽名頁面",
				})
				return
			}
		}

		// 扭蛋機QRcode處理
		if game == "3DGachaMachine" {
			activityModel, err := h.getActivityInfo(false, activityID)
			if err != nil || activityModel.ID == 0 {
				h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
				return
			}

			if activityModel.Device == "line" {
				// 只開啟line裝置驗證，liff url
				gachaMachineQRcode = fmt.Sprintf(config.HILIVES_GACHAMACHINE_GAME_LIFF_URL, activityID) + gameID
			} else {
				// 開啟多個裝置驗證，一般url
				gachaMachineQRcode = fmt.Sprintf(config.HTTPS_GACHAMACHINE_GAME_URL, config.HILIVES_NET_URL, activityID, "guest") + gameID
			}
		}

		// 投票遊戲判斷
		if game == "vote" {
			// 遊戲狀態處理在get api中處理了，這邊暫不處理
			// 判斷遊戲狀態
			var gameStatus string

			// 判斷投票結束時間
			now, _ := time.ParseInLocation("2006-01-02 15:04:05",
				time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local) // 目前時間
			startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteStartTime, time.Local)
			endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteEndTime, time.Local)

			log.Println()

			// 比較時間，判斷遊戲狀態
			if now.Before(startTime) { // 現在時間在開始時間之前
				gameStatus = "close" // 關閉
			} else if now.Before(endTime) { // 現在時間在截止時間之前
				gameStatus = "gaming" // 遊戲中
			} else { // 現在時間在截止時間之後
				gameStatus = "calculate" // 結算狀態
			}

			if gameStatus != gameModel.GameStatus {
				// 狀態不一樣，更新為最新狀態
				gameModel.GameStatus = gameStatus

				// 更新遊戲狀態(mysql)
				// if err := db.Table(config.ACTIVITY_GAME_TABLE).WithConn(h.dbConn).
				// 	Where("game_id", "=", gameModel.GameID).
				// 	Update(command.Value{"game_status": gameStatus}); err != nil &&
				// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				// 	response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
				// 		UserID:  user.UserID,
				// 		Method:  ctx.Request.Method,
				// 		Path:    ctx.Request.URL.Path,
				// 		Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
				// 	})
				// 	return
				// }

				// 更新遊戲狀態(mongo，activity_game)
				filter := bson.M{"game_id": gameModel.GameID} // 過濾條件
				// 要更新的值
				update := bson.M{
					"$set": bson.M{"game_status": gameStatus},
					// "$unset": bson.M{},                // 移除不需要的欄位
					// "$inc":   bson.M{"edit_times": 1}, // edit 欄位遞增 1
				}

				if _, err := h.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
					response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  "",
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
					})
				}

				// 修改redis中的遊戲資訊
				h.redisConn.HashSetCache(config.GAME_REDIS+gameID, "game_status", gameStatus)

				// 設置過期時間
				// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "遊戲狀態改變")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
			}

			if gameModel.GameStatus == "close" { // 關閉
				response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
					UserID:  user.UserID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 該投票遊戲場次尚未開始，請等待遊戲開始",
				})
				return
			}

		}
	}

	if game == "vote" {
		// 主持玩家端都需要投票選項資訊(照id順序排列)
		// 取得投票選項資訊
		voteOptionModels, err := models.DefaultGameVoteOptionModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindOrderByID(true, gameID)
		if err != nil {
			response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得投票選項資訊，請輸入有效的遊戲ID",
			})
			return
		}

		for _, optionModel := range voteOptionModels {
			voteOptions = append(voteOptions, GameVoteOptionModel{
				ID:              optionModel.ID,
				ActivityID:      optionModel.ActivityID,
				GameID:          optionModel.GameID,
				UserID:          optionModel.UserID,
				OptionID:        optionModel.OptionID,
				OptionName:      optionModel.OptionName,
				OptionPicture:   optionModel.OptionPicture,
				OptionIntroduce: optionModel.OptionIntroduce,
				OptionScore:     optionModel.OptionScore,
			})
		}

		if len(voteOptions) == 0 {
			response.BadRequestWithRole(ctx, role, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無投票選項資訊，請主持人設置投票選項才能遊戲",
			})
			return
		}
	}

	// 傳送角色訊息、獎品數給前端
	response.OkWithData(ctx, UserParam{
		User: UserModel{
			UserID: user.UserID,
			Name:   user.Name,
			Avatar: user.Avatar,
			Role:   role,
		},
		Game: GameModel{
			PrizeAmount:      amount,              // 獎品數
			WinPrizeAmount:   winPrizeAmount,      // 拔河遊戲勝方獎品數
			LosePrizeAmount:  losePrizeAmount,     // 拔河遊戲敗方獎品數
			GameRound:        gameModel.GameRound, // 輪次
			QARound:          gameModel.QARound,   // 題目進行題數
			QRcode:           gachaMachineQRcode,
			GameStatus:       gameModel.GameStatus,
			VoteOptions:      voteOptions,      // 投票選項資訊
			VoteTimes:        voteTimes,        // 剩餘投票次數
			VoteRecords:      voteRecords,      // 玩家投票紀錄
			VoteScore:        voteScore,        // 投票權重
			VoteMethodPlayer: voteMethodPlayer, // 投票方式
		},
		PrizeStaffs: winRecords,
		Times:       times,
	})
}

// getUserHost 處理主持端判斷(獎品數量、遊戲狀態)
func (h *Handler) getUserHost(gameID string, game string, gameModel models.GameModel) (int64, int64, int64, error) {
	var (
		winPrizeAmount, losePrizeAmount int64
	)

	// 遊戲獎品數(redis處理)
	amount, err := h.getPrizesAmount(gameID)
	if err != nil || amount == 0 {
		if game == "vote" {
			// 投票遊戲不一定要有獎品
		} else {
			// 其他遊戲
			// 清除獎品數量redis(重新計算)
			h.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)

			return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 獎品數為0，無法進入遊戲，請管理員設置獎品後才能開始遊戲")
		}
	}

	if game == "bingo" && gameModel.RoundPrize > amount {
		// 清除獎品數量redis(重新計算)
		h.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)

		return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 獎品數量必須大於每輪發獎人數，請重新設置")
	}

	// 快問快答題目數為0，不能進入遊戲畫面
	if game == "QA" && gameModel.TotalQA == 0 {
		return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 快問快答題目數為0，無法進入遊戲，請管理員設置題目後才能開始遊戲")
	}

	if game == "tugofwar" {
		// 拔河遊戲需判斷勝敗兩方剩餘獎品數是否為0
		prizes, err := models.DefaultPrizeModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindPrizes(false, gameID)
		// log.Println("獎品數: ", len(prizes))
		if err != nil || len(prizes) != 2 {
			return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 拔河遊戲獎品數有誤，無法進入遊戲")
		}

		for _, prize := range prizes {
			if prize.PrizeRemain == 0 {
				return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 拔河遊戲獎品數為0，無法進入遊戲，請管理員設置獎品後才能開始遊戲")
			}

			if gameModel.Prize == "uniform" {
				// log.Println("統一發獎")
				// 統一發獎，獎品數需大於0
				if prize.PrizeRemain < 0 {
					return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 拔河遊戲統一發獎模式時獎品數需大於0，無法進入遊戲，請管理員設置獎品後才能開始遊戲")
				}
			} else if gameModel.Prize == "all" {
				// log.Println("全部發獎")
				if prize.PrizeRemain < gameModel.People {
					return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 拔河遊戲全部發獎模式時獎品數需大於每隊上限人數，無法進入遊戲，請管理員設置獎品後才能開始遊戲")
				}
			}

			if prize.TeamType == "win" {
				winPrizeAmount = prize.PrizeRemain
				// log.Println("勝方: ", winPrizeAmount)
			} else if prize.TeamType == "lose" {
				losePrizeAmount = prize.PrizeRemain
				// log.Println("敗方: ", losePrizeAmount)
			}
		}
	}

	if game == "vote" {
		// 投票遊戲可以重複開啟頁面
	} else {
		// 遊戲已開啟，不能重複開啟
		if gameModel.GameStatus != "close" {
			// 清除遊戲資訊
			// h.redisConn.DelCache(config.GAME_REDIS + gameID)

			// 清除遊戲資訊
			// h.redisConn.DelCache(config.GAME_TYPE_REDIS + gameID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			// h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

			return amount, winPrizeAmount, losePrizeAmount, errors.New("錯誤: 該場次遊戲已被開啟，無法重複開啟遊戲")
		}
	}

	return amount, winPrizeAmount, losePrizeAmount, nil
}

// getUserGuest 處理玩家端判斷(遊戲狀態、黑名單、用戶遊戲紀錄、是否允許重複中獎、遊戲抽獎可玩次數)
func (h *Handler) getUserGuest(activityID string, gameID string, userID string,
	game string, gameModel models.GameModel) ([]models.PrizeStaffModel, int64, error) {
	var (
		winRecords  = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄
		gameRecords = make([]models.PrizeStaffModel, 0) // 用戶遊戲紀錄
		times       int64
		err         error
	)

	// 遊戲狀態未開啟
	// 遊戲抽獎不管什麼遊戲狀態下都能進入遊戲頁面(其他遊戲不行)
	if game != "lottery" && gameModel.GameStatus == "close" {
		return winRecords, times, errors.New("錯誤: 主持人未開啟遊戲，無法參加遊戲，請等待主持人開放該場次遊戲")
	}

	// 判斷遊戲玩家資訊(用戶遊戲紀錄、中獎紀錄、黑名單、是否允許重複中獎)
	winRecords, gameRecords, err = h.checkGameUser(activityID, gameID, userID, game, gameModel)
	if err != nil {
		return winRecords, times, err
	}

	// 幸運轉盤判斷可玩次數
	// 扭蛋機判斷可玩次數
	if game == "lottery" || game == "3DGachaMachine" {
		times = gameModel.MaxTimes - int64(len(gameRecords))
		if times <= 0 {
			return winRecords, times, errors.New("錯誤: 該場次抽獎次數已達上限，無法參加抽獎")
		}
	}

	// 扭蛋機遊戲
	if game == "3DGachaMachine" {
		if gameModel.GameStatus == "open" {
			// 主持端開啟遊戲畫面，玩家可進入遊戲
			// 將遊戲狀態改為start，避免其他玩家同時進入頁面
			if err = models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateGameStatus(true, gameID, "", "", "start"); err != nil {
				return winRecords, times, errors.New("錯誤: 無法更新遊戲狀態")
			}
		} else if gameModel.GameStatus != "open" {
			// 其他玩家正在遊戲中，玩家不能進入遊戲
			return winRecords, times, errors.New("錯誤: 其他玩家正在遊戲中，無法參加遊戲，請等待遊戲結束")
		}
	}

	return winRecords, times, nil
}

// checkGameUser 判斷遊戲玩家資訊(用戶遊戲紀錄、中獎紀錄、黑名單、是否允許重複中獎)
func (h *Handler) checkGameUser(activityID string, gameID string, userID string,
	game string, gameModel models.GameModel) ([]models.PrizeStaffModel, []models.PrizeStaffModel, error) {
	// 用戶遊戲紀錄(redis處理)
	// var (
	// 	activityPrizes = make([]models.PrizeStaffModel, 0) // 用戶活動中獎紀錄
	// 	gamePrizes     = make([]models.PrizeStaffModel, 0) // 用戶相同類型遊戲中獎紀錄
	// 	gameRecords    = make([]models.PrizeStaffModel, 0) // 用戶遊戲紀錄
	// )
	_, activityPrizes, gamePrizes, gameRecords, winRecords,
		err := h.getUserGameRecords(activityID, gameID, game, userID)
	if err != nil {
		return winRecords, gameRecords, err
	}

	// 判斷是否為活動黑名單、遊戲黑名單(redis處理)
	if h.IsBlackStaff(activityID, gameID, game, userID) {
		return winRecords, gameRecords, errors.New("錯誤: 用戶為黑名單人員(活動、遊戲黑名單)，無法參加遊戲。如有疑問，請聯繫主辦單位")
	}

	var allow string
	if game == "lottery" {
		allow = gameModel.LotteryGameAllow
	} else if game == "redpack" {
		allow = gameModel.RedpackGameAllow
	} else if game == "ropepack" {
		allow = gameModel.RopepackGameAllow
	} else if game == "whack_mole" {
		allow = gameModel.WhackMoleGameAllow
	} else if game == "monopoly" {
		allow = gameModel.MonopolyGameAllow
	} else if game == "QA" {
		allow = gameModel.QAGameAllow
	} else if game == "tugofwar" {
		allow = gameModel.TugofwarGameAllow
	} else if game == "bingo" {
		allow = gameModel.BingoGameAllow
	} else if game == "3DGachaMachine" {
		allow = gameModel.GachaMachineGameAllow
	}
	// 搖號抽獎目前沒有使用判斷角色api
	// else if game == "draw_numbers" {
	// 	allow = gameModel.DrawNumbersGameAllow
	// }

	if (gameModel.AllGameAllow == "close" && len(activityPrizes) > 0) ||
		(allow == "close" && len(gamePrizes) > 0) ||
		(gameModel.Allow == "close" && len(winRecords) > 0) {

		return winRecords, gameRecords, errors.New("錯誤: 不允許重複中獎(活動不允許重複中獎、相同類型遊戲不許允重複中獎、同場次遊戲不允許重複中獎)，您已中獎，無法參加遊戲")
	}

	// 判斷是否所有活動都不允許重複中獎
	if gameModel.AllGameAllow == "close" {
		// 當所有活動都不允許中獎時，判斷搖號抽獎所有場次裡的資料是否有玩家中獎資料
		// 因為玩家遊戲紀錄redis中不一定有所有搖號抽獎中獎紀錄因為抽獎時是將中獎紀錄直接寫入該活動所有搖號抽獎場次中獎人員redis中(draw_numbers_winning_staffs_ 活動ID)

		// 取得所有搖號抽獎場次中獎人員資料，確保寫入redis中(draw_numbers_winning_staffs_活動ID)
		allPrizeStaffs, err := h.getDrawNumbersAllWinningStaffs(true, activityID)
		if err != nil {
			return winRecords, gameRecords, err
		}

		if utils.InArray(allPrizeStaffs, userID) {
			// log.Println("玩家有搖號抽獎中獎紀錄")
			return winRecords, gameRecords, errors.New("錯誤: 該場次活動不允許重複中獎，您已中獎，無法參加遊戲")
		}
	}

	return winRecords, gameRecords, nil
}

// if isBlack := h.IsBlackStaff(activityID, "", "activity", user.UserID); isBlack {
// 	response.BadRequest(ctx, "錯誤: 用戶為活動黑名單人員，無法參加遊戲。如有疑問，請聯繫主辦單位")
// 	return
// }
// // 判斷是否為遊戲黑名單
// if isBlack := h.IsBlackStaff(activityID, gameID, game, user.UserID); isBlack {
// 	response.BadRequest(ctx, "錯誤: 用戶為遊戲黑名單人員，無法參加遊戲。如有疑問，請聯繫主辦單位")
// 	return
// }

// if allow == "close" && len(gamePrizes) > 0 {
// 	response.BadRequest(ctx, "錯誤: 該活動不允許相同類型遊戲重複中獎，您已中獎，無法參加遊戲")
// 	return
// }

// 判斷遊戲是否允許重複中獎
// if gameModel.Allow == "close" && len(winRecords) > 0 {
// 	response.BadRequest(ctx, "錯誤: 該遊戲不允許重複中獎，您已中獎，無法參加下一輪遊戲")
// 	return
// }

// winRecords  = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄

// for i := 0; i < len(gameRecords); i++ {
// 	if (gameRecords[i].PrizeMethod != "" && gameRecords[i].PrizeMethod != "thanks") &&
// 		(gameRecords[i].PrizeType != "" && gameRecords[i].PrizeType != "thanks") &&
// 		gameRecords[i].PrizeID != "" {
// 		winRecords = append(winRecords, gameRecords[i])
// 	}
// }

// if game == "redpack" || game == "ropepack" || game == "whack_mole" ||
// 	game == "lottery" || game == "monopoly" {
// }

// GameScore:   score,               // 分數
// GameRank:    rank + 1,            // 排名
// GameScore2:  score2,              // 玩家花費時間
// Score:       score,

// 取得用戶先前分數
// if game == "QA" {
// 	score = h.redisConn.ZSetIntScore(config.SCORES_REDIS+gameID, user.UserID)      // 分數資料
// 	rank = h.redisConn.ZSetRevRank(config.SCORES_REDIS+gameID, user.UserID)        // 排名資料
// 	score2 = h.redisConn.ZSetFloatScore(config.SCORES_2_REDIS+gameID, user.UserID) // 答題秒數
// }
// else if game == "monopoly" || game == "whack_mole" {
// 	score = h.redisConn.ZSetIntScore(config.SCORES_REDIS+gameID, user.UserID) // 分數資料
// }

// lotteryModel, err := models.DefaultLotteryModel().SetDbConn(h.dbConn).Find(gameID)
// if err != nil || lotteryModel.ID == 0 {
// 	response.BadRequest(ctx, err.Error())
// 	return
// }

// 用戶中獎紀錄
// if prizeStaffs, err = h.getPrizeStaffs(gameID, []string{"user_id", "method"},
// 	[]interface{}{user.UserID, "thanks"}); err != nil {
// 	response.BadRequest(ctx, err.Error())
// 	return
// }

// 快問快答遊戲進行中
// if game == "QA" && gameModel.GameStatus == "gaming" {
// 	response.BadRequest(ctx, "錯誤: 快問快答遊戲進行中，如要加入遊戲，請於主持人公布下一道題目前進入遊戲中")
// 	return
// }

// 遊戲已結束
// if game == "QA" && gameModel.TotalQA < gameModel.GameRound {
// 	response.BadRequest(ctx, "錯誤: 該場次快問快答遊戲已結束，謝謝")
// 	// 刪除redis所有遊戲相關資訊
// 	// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊
// 	return
// }

// GetUser 用戶的json資料 GET API
// func (h *Handler) GetUser(ctx *gin.Context) {
// 	var (
// 		userID    = ctx.Query("user_id")
// 		user, err = models.DefaultUserModel().SetDbConn(h.dbConn).Find("user_id", userID)
// 	)
// 	if userID == "" || err != nil {
// 		response.BadRequest(ctx, "發生錯誤: 無法辨識用戶資訊")
// 		return
// 	}
// 	response.OkWithData(ctx, user)
// }

// 執行登入POST API
// client    = &http.Client{}
// data      = url.Values{}
// result    interface{}
// data.Set("phone", user.Phone)
// data.Set("password", user.Password)
// req, err := http.NewRequest("POST", fmt.Sprintf(LOGIN_URL, host), strings.NewReader(data.Encode()))
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 無法登入")
// 	return
// }
// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// resp, err := client.Do(req)
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 無法登入")
// 	return
// }
// defer resp.Body.Close()
// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 無法讀取內容")
// 	return
// }
// json.Unmarshal(body, &result)
// if result.(map[string]interface{})["code"].(float64) == 400 {
// 	response.Error(ctx, result.(map[string]interface{})["msg"].(string))
// 	return
// }
// for _, cookie := range resp.Cookies() {
// 	ctx.SetCookie(cookie.N
