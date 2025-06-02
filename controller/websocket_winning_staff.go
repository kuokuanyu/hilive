package controller

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PrizeStaffModel 資料表欄位
// type PrizeStaffModel struct {
// 	ID            int64  `json:"id" example:"1"`
// 	ActivityID    string `json:"activity_id" example:"activity_id"`
// 	GameID        string `json:"game_id" example:"game_id"`
// 	PrizeID       string `json:"prize_id" example:"prize_id"`
// 	UserID        string `json:"user_id" example:"user_id"`
// 	Name          string `json:"name" example:"User Name"`
// 	Avatar        string `json:"avatar" example:"https://..."`
// 	PrizeName     string `json:"prize_name" example:"Prize Name"`
// 	PrizeType     string `json:"prize_type" example:"money、prize、thanks、first、second、third、special"`
// 	PrizePicture  string `json:"prize_picture" example:"https://..."`
// 	PrizePrice    int64  `json:"prize_price" example:"100"`
// 	PrizeMethod   string `json:"prize_method" example:"site、mail、thanks"`
// 	PrizePassword string `json:"prize_password" example:"password"`
// 	Round         int64  `json:"round" example:"1"`
// 	WinTime       string `json:"win_time" example:"2021-01-01 00:00"`
// 	Status        string `json:"status" example:"yes、no"`
// 	Score         int64  `json:"score" example:"5"`
// }

// WinningStaffParam 中獎人員資訊
type WinningStaffParam struct {
	Game        GameModel
	PrizeStaffs []models.PrizeStaffModel
	PKStaffs    []models.ScoreModel
	Message     bool   `json:"message" example:"true"`        // 是否完成結算
	Error       string `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 中獎人員與獎品數量資訊(主持人端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, lottery, monopoly, QA, bingo)
// @param body body WinningStaffParam true "param"
// @Success 200 {array} WinningStaffParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/winning/staff [get]
func (h *Handler) WinningStaffWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		gameModel         = h.getGameInfo(gameID, game) // 遊戲資訊
		staffs            = make([]models.PrizeStaffModel, 0)
		// result WinningStaffParam
		PKStaffs                                                        = make([]models.ScoreModel, 0)
		first, second, third, general, limit, amount, round, redisCount int64
		// n                                                               = 1
		// prizeStaffsLength, amount, round int64
	)
	// fmt.Println("開啟即時中獎人員、獎品資訊(主持端)ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" || gameModel.ID == 0 {
		b, _ := json.Marshal(WinningStaffParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result WinningStaffParam
		)

		data, err := conn.ReadMessage()
		// log.Println("收到中獎人員訊息")
		if err != nil {
			// fmt.Println("關閉即時中獎人員、獎品資訊(主持端)ws", err)
			if game == "redpack" || game == "ropepack" ||
				game == "whack_mole" || game == "monopoly" || game == "bingo" {

				// 清除中獎人員資料
				// h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
				// h.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + gameID) // 未中獎人員

				// 清除分數資料
				h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
				h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
				// 清除遊戲隊伍資訊
				h.redisConn.DelCache(config.GAME_TEAM_REDIS + gameID)
				// h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
				// h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

				// 賓果遊戲
				h.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID)     // 紀錄抽過的號碼，LIST
				h.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)       // 賓果中獎人員，ZSET
				h.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + gameID)      // 紀錄玩家的號碼排序，HASH
				h.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + gameID) // 紀錄玩家是否即將中獎，HASH

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				// h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
			}
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Game.GameStatus != "" {
			// 判斷遊戲狀態(從redis取得，如果沒有才執行資料表查詢)
			gameModel = h.getGameInfo(gameID, game)
			// round、limit相關判斷設置
			if game == "redpack" || game == "ropepack" {
				round = gameModel.GameRound - 1
			} else if game == "whack_mole" || game == "monopoly" || game == "QA" {
				if game == "whack_mole" || game == "monopoly" {
					round = gameModel.GameRound - 1 // 結算狀態時其他遊戲的輪次資訊還已加1
				} else if game == "QA" {
					round = gameModel.GameRound // 結算狀態時快問快答的輪次資訊還未加1
				}
				first, second, third, general = gameModel.FirstPrize, gameModel.SecondPrize,
					gameModel.ThirdPrize, gameModel.GeneralPrize
				limit = first + second + third + general
			} else if game == "lottery" {
				round = 0
			} else if game == "bingo" {
				// 賓果遊戲
				round = gameModel.GameRound - 1 // 結算狀態時其他遊戲的輪次資訊還已加1
				limit = gameModel.RoundPrize    // 每輪發獎人數
			}

			if result.Game.GameStatus != "calculate" {
				// #####加入測試資料#####start
				// if game == "redpack" || game == "ropepack" || game == "lottery" {
				// 	if n < 1000 {
				// 		for i := 1; i <= 333; i++ {
				// 			var (
				// 				record = models.PrizeStaffModel{ID: 0, PrizeName: "未中獎",
				// 					PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
				// 					PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(預設未中獎)
				// 				userID = strconv.Itoa(n)
				// 				prize  = models.PrizeModel{ID: 0, PrizeName: "未中獎",
				// 					PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
				// 					PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(預設未中獎)
				// 				prizeID string
				// 				now, _  = time.ParseInLocation("2006-01-02 15:04:05",
				// 					time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
				// 			)

				// 			// 取得LINE用戶資訊
				// 			user, _ := models.DefaultLineModel().SetDbConn(h.dbConn).
				// 				SetRedisConn(h.redisConn).Find(false, "", "user_id", userID)

				// 			// 判斷是否中獎
				// 			prizeID, _ = h.getEnvelopePrize(gameID, gameModel.Percent)

				// 			// 中獎，更新獎品數量(id != 0)
				// 			if prizeID != "" {
				// 				// 遞減redis獎品數量
				// 				h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

				// 				// 獎品資訊(redis處理)
				// 				prize, _ = models.DefaultPrizeModel().SetDbConn(h.dbConn).
				// 					SetRedisConn(h.redisConn).FindPrize(true, prizeID)
				// 			}

				// 			// 遊戲紀錄(前端利用method判斷是否中獎)
				// 			record = models.PrizeStaffModel{
				// 				ID:            0,
				// 				ActivityID:    gameModel.ActivityID,
				// 				GameID:        gameID,
				// 				PrizeID:       prize.PrizeID,
				// 				UserID:        user.UserID,
				// 				Name:          user.Name,
				// 				Avatar:        user.Avatar,
				// 				Game:          game,
				// 				PrizeName:     prize.PrizeName,
				// 				PrizeType:     prize.PrizeType,
				// 				PrizePicture:  prize.PrizePicture,
				// 				PrizePrice:    prize.PrizePrice,
				// 				PrizeMethod:   prize.PrizeMethod,
				// 				PrizePassword: prize.PrizePassword,
				// 				Round:         int64(round),
				// 				WinTime:       now.Format("2006-01-02 15:04:05"),
				// 				Status:        "no",
				// 			}

				// 			// 將中獎資訊加入redis中(競技遊戲與遊戲抽獎的順序抽獎不需要加入redis中)
				// 			if record.PrizeType != "thanks" && record.PrizeMethod != "thanks" {
				// 				// log.Println("中獎")
				// 				// 依照中獎的順序將資料推送至list格式的redis中(沒有按照獎品類型排序)
				// 				h.redisConn.ListRPush(config.WINNING_STAFFS_REDIS+gameID, utils.JSON(record))
				// 			} else if (game == "redpack" || game == "ropepack") && (record.PrizeType == "thanks" || record.PrizeMethod == "thanks") {
				// 				// log.Println("未中獎")
				// 				// 未中獎時，將未中獎紀錄寫入redis中
				// 				h.redisConn.ListRPush(config.NO_WINNING_STAFFS_REDIS+gameID, utils.JSON(record))
				// 			}

				// 			n++
				// 		}
				// 	}
				// }
				// #####加入測試資料#####end

				// 遊戲中中獎人員資訊、獎品數量資訊
				staffs, PKStaffs, amount, err = h.gaming(gameModel, gameID, game, round, redisCount)
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				// 回傳中獎人員資訊、獎品數
				b, _ := json.Marshal(WinningStaffParam{
					PrizeStaffs: staffs,
					PKStaffs:    PKStaffs,
					Game: GameModel{
						PrizeAmount: amount,
						LastPeople:  int64(len(PKStaffs)), // 即將賓果人數
						BingoPeople: int64(len(staffs)),   // 賓果人數
					},
				})
				conn.WriteMessage(b)

				redisCount++
			} else if (game == "whack_mole" || game == "monopoly" || game == "QA") &&
				gameModel.GameStatus == "calculate" {
				// log.Println("敲敲樂、鑑定師、快問快答遊戲結算")
				// time.Sleep(60 * time.Second)

				// 取得該場次所有遊戲人員資訊
				users, _ := h.redisConn.SetGetMembers(config.GAME_ATTEND_REDIS + gameID)
				// log.Println("批量將遊戲人員寫入資料表: ", len(users))

				// 批量將所有遊戲人員資料寫入資料表(敲敲樂鑑定師遊戲的遊戲中輪次已遞增、快問快答輪次未遞增)
				err = models.DefaultGameStaffModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Adds(len(users), activityID, gameID, game, users, int(round))
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				// log.Println("批量將遊戲人員寫入資料表完成")

				var winRecords = make([]models.PrizeStaffModel, 0)

				// 結算發獎(競技遊戲)，處理答題紀錄資料、分數紀錄資料、中獎人員資料、判斷發獎數量、獎品數量資料、重置快問快答題數、更新遊戲狀態資料
				winRecords, amount, err = h.calculate(gameModel, activityID, gameID, game, round,
					limit, first, second, third, general)
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				b, _ := json.Marshal(WinningStaffParam{
					PrizeStaffs: winRecords,
					Game:        GameModel{PrizeAmount: amount},
					Message:     true,
				})
				conn.WriteMessage(b)
			} else if game == "bingo" && gameModel.GameStatus == "calculate" {
				// log.Println("賓果遊戲結算")
				// time.Sleep(60 * time.Second)

				// 取得該場次所有遊戲人員資訊
				users, _ := h.redisConn.SetGetMembers(config.GAME_ATTEND_REDIS + gameID)
				// log.Println("批量將遊戲人員寫入資料表: ", len(users))

				// 批量將所有遊戲人員資料寫入資料表(賓果遊戲的遊戲中輪次已遞增)
				err = models.DefaultGameStaffModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Adds(len(users), activityID, gameID, game, users, int(round))
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				// log.Println("批量將遊戲人員寫入資料表完成")

				var (
					winRecords = make([]models.PrizeStaffModel, 0)
				)

				// 賓果遊戲結算發獎，處理分數紀錄資料、中獎人員資料、PK人員資料、發獎、重置賓果回合數
				winRecords, PKStaffs, amount, err = h.BingoCalculate(gameModel, activityID,
					gameID, game, round)
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				b, _ := json.Marshal(WinningStaffParam{
					PrizeStaffs: winRecords,
					PKStaffs:    PKStaffs,
					Game:        GameModel{PrizeAmount: amount},
					Message:     true,
				})
				conn.WriteMessage(b)
			} else if (game == "redpack" || game == "ropepack") && gameModel.GameStatus == "calculate" {
				// 間隔兩秒再執行結算
				time.Sleep(time.Second * 2)
				// log.Println("紅包遊戲結算")

				// 取得該場次所有遊戲人員資訊
				users, _ := h.redisConn.SetGetMembers(config.GAME_ATTEND_REDIS + gameID)
				// log.Println("批量將遊戲人員寫入資料表: ", len(users))

				// 批量將所有遊戲人員資料寫入資料表(紅包遊戲的遊戲中輪次已遞增)
				err = models.DefaultGameStaffModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Adds(len(users), activityID, gameID, game, users, int(round))
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				// log.Println("批量將遊戲人員寫入資料表完成")

				// 中獎人員資訊、獎品數量資訊
				staffs, amount, err = h.RedpackCalculate(gameModel, activityID, gameID, game, round)
				if err != nil {
					b, _ := json.Marshal(WinningStaffParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				// log.Println("中獎人員人數: ", len(staffs))

				// 回傳中獎人員資訊、獎品數
				b, _ := json.Marshal(WinningStaffParam{
					PrizeStaffs: staffs,
					Game:        GameModel{PrizeAmount: amount},
					Message:     true,
				})
				conn.WriteMessage(b)
			}
		}
	}
}

// gaming 遊戲中中獎人員資訊、獎品數量資訊
func (h *Handler) gaming(gameModel models.GameModel, gameID string, game string, round int64, redisCount int64) (
	[]models.PrizeStaffModel, []models.ScoreModel, int64, error) {
	var (
		staffs   = make([]models.PrizeStaffModel, 0)
		PKStaffs = make([]models.ScoreModel, 0)
		// lastPeople int64 // 即將賓果人數
		amount int64 // 獎品數量
		length int64 // 資料數量
		err    error
	)

	if game == "whack_mole" || game == "monopoly" || game == "QA" {
		// 競技類型都只傳前6名
		length = 6
	}

	// 中獎人員(從redis取得，如果沒有才執行資料表查詢)
	staffs, PKStaffs, err = h.getWinningRecords(gameModel, game, gameID, round, length, redisCount)
	if err != nil {
		return staffs, PKStaffs, amount, err
	}

	// 獎品數(從redis取得，如果沒有才執行資料表查詢)
	amount, err = h.getPrizesAmount(gameID)
	if err != nil {
		return staffs, PKStaffs, amount, err
	}
	// fmt.Println("中獎人數: ", len(staffs), ", 剩餘獎品數: ", amount)

	return staffs, PKStaffs, amount, nil
}

// RedpackCalculate 紅包遊戲結算發獎，處理中獎人員資料、未中獎人員資料、獎品資料
func (h *Handler) RedpackCalculate(gameModel models.GameModel, activityID string, gameID string, game string,
	round int64) ([]models.PrizeStaffModel, int64, error) {
	var (
		staffs   = make([]models.PrizeStaffModel, 0) // 中獎人資料
		nostaffs = make([]models.PrizeStaffModel, 0) // 未中獎人員資料
		wg       sync.WaitGroup                      // 宣告WaitGroup 用以等待執行序
		amount   int64
		err      error
	)

	// 取得紅包遊戲中獎人員資料
	staffs, _, amount, err = h.gaming(gameModel, gameID, game, round, 1)
	if err != nil {
		return staffs, amount, err
	}

	// 取得紅包遊戲未中獎人員資料
	var recordsJson []string // 未中獎人員json資料
	recordsJson, err = h.redisConn.ListRange(config.NO_WINNING_STAFFS_REDIS+gameID, 0, 0)
	if err != nil {
		return staffs, amount, errors.New("錯誤: 從redis中取得未中獎人員資料發生問題")
	}
	if len(recordsJson) > 0 {
		for _, record := range recordsJson {
			var staff models.PrizeStaffModel
			// 解碼
			json.Unmarshal([]byte(record), &staff)

			nostaffs = append(nostaffs, staff)
		}
	}

	// 批量寫入所有中獎人員資料
	// log.Println("批量寫入所有中獎人員資料: ", len(staffs))

	var (
		userIDs  = make([]string, len(staffs)) // 用戶ID
		peizeIDs = make([]string, len(staffs)) // 獎品ID
		scores1  = make([]string, len(staffs)) // 分數(0)
		scores2  = make([]string, len(staffs)) // 分數2(0)
		ranks    = make([]string, len(staffs)) // 排名資訊(0)
		teams    = make([]string, len(staffs)) // 隊伍資訊(空)
		leaders  = make([]string, len(staffs)) // 隊長資訊(空)
		mvps     = make([]string, len(staffs)) // mvp資訊(空)
	)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(staffs); i++ {
		userIDs[i] = staffs[i].UserID
		peizeIDs[i] = staffs[i].PrizeID
		scores1[i] = strconv.Itoa(0)
		scores2[i] = strconv.Itoa(0)
		ranks[i] = strconv.Itoa(0)
	}

	// 將所有中獎人員資料批量寫入資料表中
	if err = models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(len(staffs), activityID, gameID, game, userIDs, peizeIDs,
			int(round), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
		return staffs, amount, errors.New("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
	}

	// log.Println("批量寫入所有中獎人員完成")

	// 批量寫入所有未中獎人員資料
	// log.Println("批量寫入所有未中獎人員資料: ", len(nostaffs))

	// 清空舊的資料
	userIDs = make([]string, len(nostaffs))  // 用戶ID
	peizeIDs = make([]string, len(nostaffs)) // 獎品ID(空)
	scores1 = make([]string, len(nostaffs))  // 分數(0)
	scores2 = make([]string, len(nostaffs))  // 分數2(0)
	ranks = make([]string, len(nostaffs))    // 排名資訊(0)
	teams = make([]string, len(nostaffs))    // 隊伍資訊(空)
	leaders = make([]string, len(nostaffs))  // 隊長資訊(空)
	mvps = make([]string, len(nostaffs))     // mvp資訊(空)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(nostaffs); i++ {
		userIDs[i] = nostaffs[i].UserID
		scores1[i] = strconv.Itoa(0)
		scores2[i] = strconv.Itoa(0)
		ranks[i] = strconv.Itoa(0)
	}

	// 將所有未中獎人員資料批量寫入資料表中
	if err = models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(len(nostaffs), activityID, gameID, game, userIDs, peizeIDs,
			int(round), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有未中獎人員資料寫入資料表發生問題，請重新操作")
		return staffs, amount, errors.New("錯誤: 將所有未中獎人員資料寫入資料表發生問題，請重新操作")
	}

	// log.Println("批量寫入所有未中獎人員完成")

	// 處理剩餘獎品
	// 從獎品redis(hash格式)中取得所有獎品剩餘數量，並更新至資料表中(多線呈處理)
	prizes, err := h.getPrizes(gameID)
	if err != nil {
		return []models.PrizeStaffModel{}, 0, err
	}
	// log.Println("更新紅包遊戲剩餘獎品: ", len(prizes))

	// 多線呈更新資料表獎品資訊
	wg.Add(len(prizes)) //計數器
	for i := 0; i < len(prizes); i++ {
		go func(i int) {
			defer wg.Done()

			// 更新獎品剩餘數量
			if err = models.DefaultPrizeModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateRemain(
					gameID, prizes[i].PrizeID, prizes[i].PrizeRemain); err != nil {
				// log.Println(err)
				return
			}
		}(i)
	}
	wg.Wait() //等待計數器歸0

	// log.Println("更新紅包遊戲剩餘完成")

	return staffs, amount, nil
}

// BingoCalculate 賓果遊戲結算發獎，處理分數紀錄資料、中獎人員資料、PK人員資料、發獎、重置賓果回合數
func (h *Handler) BingoCalculate(gameModel models.GameModel, activityID string, gameID string, game string,
	round int64) ([]models.PrizeStaffModel, []models.ScoreModel, int64, error) {
	// 間隔1秒後開始統計中獎人員，確保玩家端都接收到遊戲狀態並完成最後一次的分數更新
	time.Sleep(time.Second * 1)

	var (
		scores   = make([]models.ScoreModel, 0)
		bingos   = make([]models.ScoreModel, 0)
		winUsers = make([]models.ScoreModel, 0)
		pkUsers  = make([]models.ScoreModel, 0)
		now, _   = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
		// lock   sync.Mutex     // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
		wg     sync.WaitGroup // 宣告WaitGroup 用以等待執行序
		prize  models.PrizeModel
		amount int64
		err    error
	)

	// 查詢賓果遊戲所有人員分數資料、即將賓果人員人數資料(從redis取得)
	scores, _, _ = models.DefaultScoreModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindBingoScore(true, gameID, gameModel.BingoLine)

	// 查詢賓果人員資料(從redis取得)
	bingos, _ = models.DefaultScoreModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindBingoUser(true, gameID)

	// 優化後不會有賓果人數大於每輪發獎人數的場景
	winUsers = bingos

	// log.Println("賓果遊戲批量寫入分數人員資料: ", len(scores))
	var (
		userIDs = make([]string, len(scores)) // 用戶ID
		scores1 = make([]string, len(scores)) // 分數
		scores2 = make([]string, len(scores)) // 分數2
		teams   = make([]string, len(scores)) // 隊伍資訊(空)
		leaders = make([]string, len(scores)) // 隊長資訊(空)
		mvps    = make([]string, len(scores)) // mvp資訊(空)
	)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(scores); i++ {
		userIDs[i] = scores[i].UserID
		scores1[i] = strconv.Itoa(int(scores[i].Score))
		scores2[i] = strconv.Itoa(int(scores[i].Score2))
	}

	// 將所有人員分數資料批量寫入資料表中
	if err = models.DefaultScoreModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(len(scores), activityID, gameID, game, userIDs,
			int(round), scores1, scores2, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有計分資料寫入資料表發生問題，請重新操作")
		return []models.PrizeStaffModel{}, pkUsers, amount, errors.New("錯誤: 將所有計分資料寫入資料表發生問題，請重新操作")
	}
	// log.Println("賓果遊戲批量寫入分數人員完成")

	// 發放獎品給中獎人員
	var winRecords = make([]models.PrizeStaffModel, len(winUsers))

	// 取得賓果遊戲獎品資訊(資料表)
	prizeModel, err := models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindPrizes(false, gameID)
	if err != nil || len(prizeModel) < 1 {
		return winRecords, pkUsers, amount, errors.New("錯誤: 無法取得獎品資訊，請重新操作")
	}
	prize = prizeModel[0] // 賓果遊戲獎品數只有一個
	// 剩餘獎品數
	amount = prize.PrizeRemain - int64(len(winUsers))

	// 批量寫入所有中獎人員資料
	// log.Println("批量寫入所有中獎人員資料: ", len(winUsers))

	// 重新清空參數
	userIDs = make([]string, len(winUsers))   // 用戶ID
	peizeIDs := make([]string, len(winUsers)) // 獎品ID
	scores1 = make([]string, len(winUsers))   // 分數
	scores2 = make([]string, len(winUsers))   // 分數2
	ranks := make([]string, len(winUsers))    // 排名資訊
	teams = make([]string, len(winUsers))     // 隊伍資訊(空)
	leaders = make([]string, len(winUsers))   // 隊長資訊(空)
	mvps = make([]string, len(winUsers))      // mvp資訊(空)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(winUsers); i++ {
		userIDs[i] = winUsers[i].UserID
		peizeIDs[i] = prize.PrizeID
		scores1[i] = strconv.Itoa(int(winUsers[i].Score))
		scores2[i] = strconv.Itoa(int(winUsers[i].Score2))
		ranks[i] = strconv.Itoa(i + 1)
	}

	// 將所有中獎人員資料批量寫入資料表中
	if err = models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(len(winUsers), activityID, gameID, game, userIDs, peizeIDs,
			int(round), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
		return []models.PrizeStaffModel{}, pkUsers, amount, errors.New("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
	}

	// log.Println("批量寫入所有中獎人員完成")

	// 多線呈更新所有玩家遊戲紀錄(redis)
	wg.Add(len(winUsers)) //計數器
	for i := 0; i < len(winUsers); i++ {
		go func(i int) {
			defer wg.Done()
			// 用戶遊戲紀錄
			record := models.PrizeStaffModel{
				ID:            0,
				ActivityID:    activityID,
				GameID:        gameID,
				PrizeID:       prize.PrizeID,
				UserID:        winUsers[i].UserID,
				Name:          winUsers[i].Name,
				Avatar:        winUsers[i].Avatar,
				Game:          game,
				PrizeName:     prize.PrizeName,
				PrizeType:     prize.PrizeType,
				PrizePicture:  prize.PrizePicture,
				PrizePrice:    prize.PrizePrice,
				PrizeMethod:   prize.PrizeMethod,
				PrizePassword: prize.PrizePassword,
				Round:         round,
				WinTime:       now.Format("2006-01-02 15:04:05"),
				Status:        "no",
				Score:         winUsers[i].Score,
				Score2:        winUsers[i].Score2,
				Rank:          int64(i + 1),
			}

			// 用戶遊戲紀錄(包含中獎與未中獎)
			allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, winUsers[i].UserID)
			if err != nil {
				// log.Println(err)
				return
			}
			allRecords = append(allRecords, record)

			// 更新redis裡用戶遊戲紀錄
			h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
				winUsers[i].UserID, utils.JSON(allRecords))

			// 設置過期時間
			// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)

			// 中獎名單
			winRecords[i] = record
		}(i)
	}
	wg.Wait() //等待計數器歸0

	// 更新獎品剩餘數量(資料表、redis)
	if err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		UpdateRemain(
			gameID, prize.PrizeID, amount); err != nil {
		return winRecords, pkUsers, amount, err
	}

	return winRecords, pkUsers, amount, nil
}

// calculate 結算發獎(競技遊戲)，處理答題紀錄資料、分數紀錄資料、中獎人員資料、判斷發獎數量、獎品數量資料、重置快問快答題數、更新遊戲狀態資料
func (h *Handler) calculate(gameModel models.GameModel, activityID string, gameID string, game string,
	round int64, limit int64, first int64, second int64,
	third int64, general int64) ([]models.PrizeStaffModel, int64, error) {
	// 間隔1秒後開始統計中獎人員，確保玩家端都接收到遊戲狀態並完成最後一次的分數更新
	time.Sleep(time.Second * 1)

	var (
		wg        sync.WaitGroup                               // 宣告WaitGroup 用以等待執行序
		qaRecords = make(map[string]map[string]interface{}, 0) // 答題紀錄
		now, _    = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)
	// startTime := time.Now()
	// fmt.Println("遊戲結束，統計中獎人員資料")
	// fmt.Println("取得所有獎項設置的人數: ", first, ", ", second, ", ", third, ", ", general)

	// 取得分數前n高的人員資料(從redis取得)
	staffs, _, err := h.getWinningRecords(gameModel, game, gameID, round, 0, 0)
	if err != nil {
		return []models.PrizeStaffModel{}, 0, err
	}

	if game == "QA" {
		// 取得快問快答答題紀錄
		qaRecords, err = models.DefaultGameQARecordModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAllByRedis(true, gameID)
		if err != nil {
			return []models.PrizeStaffModel{}, 0, err
		}
	}

	// log.Println("競技遊戲批量寫入分數人員資料: ", len(staffs))
	var (
		userIDs = make([]string, len(staffs)) // 用戶ID
		scores1 = make([]string, len(staffs)) // 分數
		scores2 = make([]string, len(staffs)) // 分數2
		teams   = make([]string, len(staffs)) // 隊伍資訊(空)
		leaders = make([]string, len(staffs)) // 隊長資訊(空)
		mvps    = make([]string, len(staffs)) // mvp資訊(空)
	)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(staffs); i++ {
		userIDs[i] = staffs[i].UserID
		scores1[i] = strconv.Itoa(int(staffs[i].Score))
		scores2[i] = strconv.Itoa(int(staffs[i].Score2))
	}

	// 將所有人員分數資料批量寫入資料表中
	if err = models.DefaultScoreModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(len(staffs), activityID, gameID, game, userIDs,
			int(round), scores1, scores2, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有計分資料寫入資料表發生問題，請重新操作")
		return []models.PrizeStaffModel{}, 0, errors.New("錯誤: 將所有計分資料寫入資料表發生問題，請重新操作")
	}
	// log.Println("競技遊戲批量寫入分數人員資料完成")

	if game == "QA" {
		// log.Println("快問快答遊戲批量寫入答題紀錄資料: ", len(qaRecords))

		// 批量新增快問快答答題紀錄至資料表
		err = models.DefaultGameQARecordModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Adds(len(qaRecords), activityID,
				gameID, int(round), qaRecords)
		if err != nil {
			// log.Println(err)
			return []models.PrizeStaffModel{}, 0, errors.New("錯誤: 將所有玩家答題記錄資料寫入資料表發生問題，請重新操作")
		}

		// log.Println("快問快答遊戲批量寫入答題紀錄完成")
	}

	var order string
	if game == "whack_mole" || game == "monopoly" {
		order = "desc" // 正確率由高至低
	} else if game == "QA" {
		order = "asc" // 耗時排序由低至高
	}
	// 取得資料表中用戶排名資料(大於0)
	scores, err := models.DefaultScoreModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindTopUser(false, gameID, game, round, limit, order)
	if err != nil {
		return []models.PrizeStaffModel{}, 0, err
	}

	var (
		allprizes           = make([]models.PrizeModel, 0)   // 發放獎品資訊
		allTypePrizes       = make([][]models.PrizeModel, 0) // 各獎項獎品資訊
		allTypePrizesAmount = make([]int64, 0)               // 各獎項獎品總數
		allPrizesPeople     = make([]int64, 0)               // 各獎項需要發獎的人數
	)

	// 如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量
	if len(scores) < int(limit) {
		limit = int64(len(scores))
		if limit <= first {
			first, second, third, general = limit, 0, 0, 0
		} else if limit > first && limit <= first+second {
			second, third, general = limit-first, 0, 0
		} else if limit > first+second && limit <= first+second+third {
			third, general = limit-first-second, 0
		} else if limit > first+second+third && limit < first+second+third+general {
			general = limit - first - second - third
		}
	}
	allPrizesPeople = []int64{first, second, third, general}
	// fmt.Println("如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量: ",
	// 	first, ", ", second, ", ", third, ", ", general)

	// 取得各獎項所有資訊、數量
	if allTypePrizes, allTypePrizesAmount, err = h.getAllTypePrizes(gameID); err != nil {
		return []models.PrizeStaffModel{}, 0, err
	}

	// 判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值
	for i := 0; i < len(allPrizesPeople); i++ {
		if allTypePrizesAmount[i] < allPrizesPeople[i] {
			if i < len(allPrizesPeople)-1 {
				// 不夠發放，將多餘的值設置至下一個獎項
				allPrizesPeople[i+1] = allPrizesPeople[i+1] -
					allTypePrizesAmount[i] + allPrizesPeople[i]
			}
			allPrizesPeople[i] = allTypePrizesAmount[i]
		}
	}
	limit = allPrizesPeople[0] + allPrizesPeople[1] +
		allPrizesPeople[2] + allPrizesPeople[3]
	// fmt.Println("判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值: ",
	// 	allPrizesPeople[0], ", ", allPrizesPeople[1], ", ", allPrizesPeople[2], ", ", allPrizesPeople[3])

	// 將需要發放的獎品都append allprizes
	for i := 0; i < len(allPrizesPeople); i++ {
		var stopValue int64 // 當len(allprizes)=stopValue，停止該獎項發放並發放下一個獎項
		for m := 0; m <= i; m++ {
			stopValue += allPrizesPeople[m]
		}
		// fmt.Println("i: ", i, ", stopValue: ", stopValue)

		for _, prize := range allTypePrizes[i] {
			// 該獎項發獎人數為0，直接不執行
			if int64(len(allprizes)) >= stopValue {
				break
			}

			for n := 0; n < int(prize.PrizeRemain); n++ {
				allprizes = append(allprizes, prize)
				if int64(len(allprizes)) >= stopValue {
					break
				}
			}

			if int64(len(allprizes)) >= stopValue {
				break
			}
		}
	}
	// fmt.Println("allprizes: ", allprizes)

	// 批量寫入所有中獎人員資料
	// log.Println("批量寫入所有中獎人員資料: ", int(limit))

	// 重新清空參數
	userIDs = make([]string, int(limit))   // 用戶ID
	peizeIDs := make([]string, int(limit)) // 獎品ID
	scores1 = make([]string, int(limit))   // 分數
	scores2 = make([]string, int(limit))   // 分數2
	ranks := make([]string, int(limit))    // 排名資訊
	teams = make([]string, int(limit))     // 隊伍資訊(空)
	leaders = make([]string, int(limit))   // 隊長資訊(空)
	mvps = make([]string, int(limit))      // mvp資訊(空)

	// 處理批量寫入資料的相關參數
	for i := 0; i < int(limit); i++ {
		userIDs[i] = scores[i].UserID
		peizeIDs[i] = allprizes[i].PrizeID
		scores1[i] = strconv.Itoa(int(scores[i].Score))
		scores2[i] = strconv.Itoa(int(scores[i].Score2))
		ranks[i] = strconv.Itoa(i + 1)
	}

	// 將所有中獎人員資料批量寫入資料表中
	if err = models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Adds(int(limit), activityID, gameID, game, userIDs, peizeIDs,
			int(round), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
		// log.Println("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
		return []models.PrizeStaffModel{}, 0, errors.New("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
	}

	// log.Println("批量寫入所有中獎人員完成")

	// 多線呈更新所有玩家遊戲紀錄(redis)
	var winRecords = make([]models.PrizeStaffModel, limit)
	wg.Add(int(limit)) //計數器
	for i := 0; i < int(limit); i++ {
		go func(i int) {
			defer wg.Done()
			// 遞減redis獎品數量
			h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, allprizes[i].PrizeID)

			// 設置過期時間
			// h.redisConn.SetExpire(config.GAME_PRIZES_AMOUNT_REDIS+gameID, config.REDIS_EXPIRE)

			// 用戶遊戲紀錄
			record := models.PrizeStaffModel{
				ID:            0,
				ActivityID:    activityID,
				GameID:        gameID,
				PrizeID:       allprizes[i].PrizeID,
				UserID:        scores[i].UserID,
				Name:          scores[i].Name,
				Avatar:        scores[i].Avatar,
				Game:          game,
				PrizeName:     allprizes[i].PrizeName,
				PrizeType:     allprizes[i].PrizeType,
				PrizePicture:  allprizes[i].PrizePicture,
				PrizePrice:    allprizes[i].PrizePrice,
				PrizeMethod:   allprizes[i].PrizeMethod,
				PrizePassword: allprizes[i].PrizePassword,
				Round:         round,
				WinTime:       now.Format("2006-01-02 15:04:05"),
				Status:        "no",
				Score:         scores[i].Score,
				Score2:        scores[i].Score2,
				Rank:          int64(i + 1),
			}

			// 用戶遊戲紀錄(包含中獎與未中獎)
			allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, scores[i].UserID)
			if err != nil {
				// log.Println(err)
				return
			}
			allRecords = append(allRecords, record)

			// 更新redis裡用戶遊戲紀錄
			h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
				scores[i].UserID, utils.JSON(allRecords))

			// 設置過期時間
			// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)

			// 中獎名單
			winRecords[i] = record
		}(i)
	}
	wg.Wait() //等待計數器歸0

	// log.Println("玩家遊戲紀錄處理完成")

	// 從獎品redis(hash格式)中取得所有獎品剩餘數量，並更新至資料表中(多線呈處理)
	var (
		amount int64
		lock   sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
	)
	prizes, err := h.getPrizes(gameID)
	if err != nil {
		return []models.PrizeStaffModel{}, 0, err
	}

	// 多線呈更新資料表獎品資訊
	wg.Add(len(prizes)) //計數器
	for i := 0; i < len(prizes); i++ {
		go func(i int) {
			defer wg.Done()

			// 計算剩餘獎品數
			// for l := 0; l < MaxRetries; l++ {
			// 上鎖
			// ok, _ := h.acquireLock(config.PRIZE_LOCK_REDIS+gameID, LockExpiration)
			// if ok == "OK" {
			// 釋放鎖
			// defer h.releaseLock(config.PRIZE_LOCK_REDIS + gameID)

			lock.Lock() //佔有資源
			amount += prizes[i].PrizeRemain
			lock.Unlock() //釋放資源

			// 釋放鎖
			// h.releaseLock(config.PRIZE_LOCK_REDIS + gameID)
			// break
			// }

			// 鎖被佔用，稍微延遲後重試
			// time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
			// }

			// 更新獎品剩餘數量
			if err = models.DefaultPrizeModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateRemain(
					gameID, prizes[i].PrizeID, prizes[i].PrizeRemain); err != nil {
				// log.Println(err)
				return
			}
		}(i)
	}
	wg.Wait() //等待計數器歸0

	// log.Println("獎品數量處理完成: ", len(prizes))

	// 快問快答結算後更新輪次以及題目進行題數資訊(清空qa_people人數)
	if game == "QA" {
		// 清空redis題目相關資訊(所有資訊歸零)
		h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
			"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

		// 設置過期時間
		// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")

		// 更新題目進行題數資料(qa_round=1)
		if err = models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateQARound(true, gameID, "1"); err != nil {
			return []models.PrizeStaffModel{}, 0, errors.New("錯誤: 無法更新題目進行題數資料(qa_round=1)")
		}

		// 清空遊戲人員資料
		h.redisConn.DelCache(config.GAME_ATTEND_REDIS + gameID)

		if err = models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateGameStatus(true, gameID, "+1", "", "calculate"); err != nil {
			return []models.PrizeStaffModel{}, 0, errors.New("錯誤: 無法更新遊戲狀態")
		}
	}

	// log.Println("回傳前端")

	return winRecords, amount, nil
}

// getWinningRecords 查詢該遊戲場次的中獎紀錄(從redis取得，如果沒有才執行資料表查詢)
func (h *Handler) getWinningRecords(gameModel models.GameModel, game, gameID string, round,
	limit, redisCount int64) ([]models.PrizeStaffModel, []models.ScoreModel, error) {
	var (
		staffs   = make([]models.PrizeStaffModel, 0)
		PKStaffs = make([]models.ScoreModel, 0)
		err      error
	)
	if game == "redpack" || game == "ropepack" || game == "lottery" {
		staffs, err = models.DefaultPrizeStaffModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindRedpackAndLotteryWinningStaffs(true, redisCount, gameID, round, limit)
	} else if game == "whack_mole" || game == "monopoly" || game == "QA" || game == "tugofwar" {
		var scores = make([]models.ScoreModel, 0)
		// 取得分數前n高的人員資料(從redis取得)
		scores, err = models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindTopUser(true, gameID, game, round, limit, "")

		for _, score := range scores {
			staffs = append(staffs, models.PrizeStaffModel{
				ID:     score.ID,
				UserID: score.UserID,
				Name:   score.Name,
				Avatar: score.Avatar,
				Score:  score.Score,
				Score2: score.Score2,
			})
		}
	} else if game == "bingo" {
		// 查詢賓果遊戲所有人員分數資料、即將賓果人員人數資料(從redis中取得分數由高至低的玩家資訊)
		_, PKStaffs, err = models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindBingoScore(true, gameID, gameModel.BingoLine)
		if err != nil {
			return staffs, PKStaffs, err
		}
		// fmt.Println("即將賓果的人數", len(PKStaffs))

		// 取得賓果人員資料(從redis取得)
		var scores = make([]models.ScoreModel, 0)
		scores, err = models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindBingoUser(true, gameID)

		// log.Println("賓果人數: ", len(scores))

		for _, score := range scores {
			staffs = append(staffs, models.PrizeStaffModel{
				ID:     score.ID,
				UserID: score.UserID,
				Name:   score.Name,
				Avatar: score.Avatar,
				Score:  score.Score,
				Score2: score.Score2,
			})
		}
	}
	if err != nil {
		return staffs, PKStaffs, err
	}
	return staffs, PKStaffs, err
}

// getAllTypePrizes 取得打地鼠、超級大富翁遊戲所有獎品資訊(各獎項數量)
func (h *Handler) getAllTypePrizes(gameID string) ([][]models.PrizeModel,
	[]int64, error) {
	var (
		firstPrizes                                           = make([]models.PrizeModel, 0)
		secondPrizes                                          = make([]models.PrizeModel, 0)
		thirdPrizes                                           = make([]models.PrizeModel, 0)
		generalPrizes                                         = make([]models.PrizeModel, 0)
		allTypePrizes                                         = make([][]models.PrizeModel, 0)
		allTypePrizesAmount                                   = make([]int64, 0)
		firstAmount, secondAmount, thirdAmount, generalAmount int64
	)
	prizeModel, err := models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindPrizes(false, gameID)
	if err != nil {
		return allTypePrizes, allTypePrizesAmount, err
	}

	// 判斷頭獎、二獎、三獎、普通獎
	for i := 0; i < len(prizeModel); i++ {
		if prizeModel[i].PrizeRemain == 0 {
			continue
		}

		// 將獎項加入陣列中
		if prizeModel[i].PrizeType == "first" {
			firstPrizes = append(firstPrizes, prizeModel[i])
			firstAmount += prizeModel[i].PrizeRemain
		} else if prizeModel[i].PrizeType == "second" {
			secondPrizes = append(secondPrizes, prizeModel[i])
			secondAmount += prizeModel[i].PrizeRemain
		} else if prizeModel[i].PrizeType == "third" {
			thirdPrizes = append(thirdPrizes, prizeModel[i])
			thirdAmount += prizeModel[i].PrizeRemain
		} else if prizeModel[i].PrizeType == "general" {
			generalPrizes = append(generalPrizes, prizeModel[i])
			generalAmount += prizeModel[i].PrizeRemain
		}
	}
	allTypePrizes = append(allTypePrizes, firstPrizes, secondPrizes, thirdPrizes, generalPrizes)
	allTypePrizesAmount = append(allTypePrizesAmount, firstAmount, secondAmount, thirdAmount, generalAmount)
	// fmt.Println("各獎項獎品: ", allTypePrizes, ", 各獎項獎品數: ", allTypePrizesAmount)
	return allTypePrizes, allTypePrizesAmount, err
}

// 新增中獎人員名單
// id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	Add(models.NewPrizeStaffModel{
// 		ActivityID: activityID,
// 		GameID:     gameID,
// 		UserID:     scores[i].UserID,
// 		PrizeID:    allprizes[i].PrizeID,
// 		Round:      strconv.Itoa(int(round)),
// 		Status:     "no",
// 		Score:      scores[i].Score,
// 		Score2:     scores[i].Score2,
// 		Rank:       int64(i + 1),
// 		// White:      "no",
// 	})
// if err != nil {
// 	log.Println(err)
// 	return
// }

// 多線呈將所有人員分數資料加入資料表中(分數大於0)
// wg.Add(len(staffs)) //計數器
// for i := 0; i < len(staffs); i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		if err = models.DefaultScoreModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).Add(models.NewScoreModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     staffs[i].UserID,
// 			Round:      round,
// 			Score:      staffs[i].Score,
// 			Score2:     staffs[i].Score2,
// 			// Rank:       int64(i + 1),
// 		}); err != nil {
// 			log.Println("錯誤: 新增計分資料發生問題，請重新操作")
// 			return
// 		}

// 		if game == "QA" {
// 			// 新增快問快答答題紀錄至資料表
// 			err = models.DefaultGameQARecordModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).Add(staffs[i].UserID, activityID,
// 				gameID, round, qaRecords[staffs[i].UserID])
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 		}
// 	}(i)
// }
// wg.Wait() //等待計數器歸0

// if game == "bingo" {
// 	// 更新賓果回合資料(bingo_round=0)
// 	if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).UpdateBingoRound(true, gameID, "0"); err != nil {
// 		return winRecords, pkUsers, amount, errors.New("錯誤: 無法更新賓果回合資料")
// 	}
// }

// 遞減redis獎品數量
// h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prize.PrizeID)

// 計算剩餘獎品數
// for l := 0; l < MaxRetries; l++ {
// 	// 上鎖
// 	ok, _ := h.acquireLock(config.PRIZE_LOCK_REDIS+gameID, LockExpiration)
// 	if ok == "OK" {
// 		// 釋放鎖
// 		// defer h.releaseLock(config.PRIZE_LOCK_REDIS + gameID)
// 		// lock.Lock() //佔有資源
// 		amount--
// 		// lock.Unlock() //釋放資源

// 		// 釋放鎖
// 		h.releaseLock(config.PRIZE_LOCK_REDIS + gameID)
// 		break
// 	}

// 	// 鎖被佔用，稍微延遲後重試
// 	time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
// }

// 新增中獎人員名單
// id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	Add(models.NewPrizeStaffModel{
// 		ActivityID: activityID,
// 		GameID:     gameID,
// 		UserID:     winUsers[i].UserID,
// 		PrizeID:    prize.PrizeID,
// 		Round:      strconv.Itoa(int(round)),
// 		Status:     "no",
// 		Score:      winUsers[i].Score,
// 		Score2:     winUsers[i].Score2,
// 		Rank:       int64(i + 1),
// 		// White:      "no",
// 	})
// if err != nil {
// 	log.Println(err)
// 	return
// }

// 多線呈將所有人員分數資料加入資料表中
// wg.Add(len(scores)) //計數器
// for i := 0; i < len(scores); i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		if err = models.DefaultScoreModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).Add(models.NewScoreModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     scores[i].UserID,
// 			Round:      round,
// 			Score:      scores[i].Score,
// 			Score2:     scores[i].Score2,
// 		}); err != nil {
// 			log.Println("錯誤: 新增計分資料發生問題，請重新操作")
// 			return
// 		}
// 	}(i)
// }
// wg.Wait() //等待計數器歸0

// 多線呈將PK人員資料加入資料表中
// wg.Add(len(pkUsers)) //計數器
// for i := 0; i < len(pkUsers); i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		if err = models.DefaultPKStaffModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).Add(models.NewPKStaffModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     pkUsers[i].UserID,
// 			Round:      round,
// 		}); err != nil {
// 			log.Println("錯誤: 新增PK人員資料發生問題，請重新操作")
// 			return
// 		}
// 	}(i)
// }
// wg.Wait() //等待計數器歸0

// 間隔n秒傳遞中獎訊息版本
// go func() {
// 	for {
// 		// fmt.Println("3")
// 		// 避免前端開啟沒收到資訊
// 		if game == "QA" {
// 			time.Sleep(time.Second * 1 / 2)
// 		} else {
// 			time.Sleep(time.Second * 7 / 2)
// 		}

// 		// 判斷遊戲狀態(從redis取得，如果沒有才執行資料表查詢)
// 		gameModel = h.getGameInfo(gameID)
// 		// round、limit相關判斷設置
// 		if game == "redpack" || game == "ropepack" {
// 			round = gameModel.GameRound - 1
// 			// limit = 0
// 		} else if game == "whack_mole" || game == "monopoly" || game == "QA" {
// 			if game == "whack_mole" || game == "monopoly" {
// 				round = gameModel.GameRound - 1 // 結算狀態時其他遊戲的輪次資訊還已加1
// 			} else if game == "QA" {
// 				round = gameModel.GameRound // 結算狀態時快問快答的輪次資訊還未加1
// 			}
// 			first, second, third, general = gameModel.FirstPrize, gameModel.SecondPrize,
// 				gameModel.ThirdPrize, gameModel.GeneralPrize
// 			limit = first + second + third + general
// 		} else if game == "lottery" {
// 			round = 0
// 			// limit = 0
// 		}

// 		// fmt.Println("遊戲狀態: ", gameModel.GameStatus)
// 		if gameModel.GameStatus != "calculate" {
// 			// startTime := time.Now()
// 			var length int64 // 取得資料數量
// 			// 中獎人員(從redis取得，如果沒有才執行資料表查詢)
// 			if game == "whack_mole" || game == "monopoly" || game == "QA" {
// 				// 競技類型都只傳前6名
// 				length = 6
// 			}
// 			staffs, _ = h.getWinningRecords(game, gameID, round, length, redisCount)

// 			// 獎品數(從redis取得，如果沒有才執行資料表查詢)
// 			amount, _ = h.getPrizesAmount(gameID)

// 			// 回傳中獎人員資訊、獎品數
// 			b, _ := json.Marshal(WinningStaffParam{
// 				PrizeStaffs: staffs,
// 				Game:        GameModel{PrizeAmount: amount},
// 			})
// 			// fmt.Println("中獎人數: ", len(staffs), ", 剩餘獎品數: ", amount)
// 			if err := conn.WriteMessage(b); err != nil {
// 				return
// 			}
// 			redisCount++
// 			// endTime := time.Now()
// 			// log.Println("回傳排名花費時間: ", endTime.Sub(startTime))
// 		} else if (game == "whack_mole" || game == "monopoly" || game == "QA") &&
// 			gameModel.GameStatus == "calculate" {
// 			var (
// 				wg        sync.WaitGroup                               // 宣告WaitGroup 用以等待執行序
// 				qaRecords = make(map[string]map[string]interface{}, 0) // 答題紀錄
// 				now, _    = time.ParseInLocation("2006-01-02 15:04:05",
// 					time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// 			)
// 			// 結算
// 			// 間隔1秒後開始統計中獎人員，確保玩家端都接收到遊戲狀態並完成最後一次的分數更新
// 			time.Sleep(time.Second * 1)
// 			// startTime := time.Now()
// 			// fmt.Println("遊戲結束，統計中獎人員資料")
// 			// fmt.Println("取得所有獎項設置的人數: ", first, ", ", second, ", ", third, ", ", general)

// 			// 取得分數前n高的人員資料(從redis取得)
// 			staffs, err = h.getWinningRecords(game, gameID, round, 0, 0)
// 			if err != nil {
// 				b, _ := json.Marshal(WinningStaffParam{
// 					Error: err.Error()})
// 				conn.WriteMessage(b)
// 				return
// 			}

// 			if game == "QA" {
// 				// 取得快問快答答題紀錄
// 				qaRecords, err = models.DefaultGameQARecordModel().SetDbConn(h.dbConn).
// 					SetRedisConn(h.redisConn).FindAllByRedis(true, gameID)
// 				if err != nil {
// 					b, _ := json.Marshal(WinningStaffParam{
// 						Error: err.Error()})
// 					conn.WriteMessage(b)
// 					return
// 				}
// 			}

// 			// 多線呈將所有人員分數資料加入資料表中(分數大於0)
// 			wg.Add(len(staffs)) //計數器
// 			for i := 0; i < len(staffs); i++ {
// 				go func(i int) {
// 					defer wg.Done()
// 					if err = models.DefaultScoreModel().SetDbConn(h.dbConn).
// 						SetRedisConn(h.redisConn).Add(models.NewScoreModel{
// 						ActivityID: activityID,
// 						GameID:     gameID,
// 						UserID:     staffs[i].UserID,
// 						Round:      round,
// 						Score:      staffs[i].Score,
// 						Score2:     staffs[i].Score2,
// 					}); err != nil {
// 						b, _ := json.Marshal(WinningStaffParam{
// 							Error: "錯誤: 新增計分資料發生問題，請重新操作"})
// 						conn.WriteMessage(b)
// 						return
// 					}

// 					if game == "QA" {
// 						// 新增快問快答答題紀錄至資料表
// 						err = models.DefaultGameQARecordModel().SetDbConn(h.dbConn).
// 							SetRedisConn(h.redisConn).Add(staffs[i].UserID, activityID,
// 							gameID, round, qaRecords[staffs[i].UserID])
// 						if err != nil {
// 							b, _ := json.Marshal(WinningStaffParam{
// 								Error: err.Error()})
// 							conn.WriteMessage(b)
// 							return
// 						}
// 					}
// 				}(i)
// 			}
// 			wg.Wait() //等待計數器歸0

// 			var order string
// 			if game == "whack_mole" || game == "monopoly" {
// 				order = "desc" // 正確率由高至低
// 			} else if game == "QA" {
// 				order = "asc" // 耗時排序由低至高
// 			}
// 			// 取得資料表中用戶排名資料(大於0)
// 			scores, err := models.DefaultScoreModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).Find(false, gameID, round, limit, order)
// 			if err != nil {
// 				b, _ := json.Marshal(WinningStaffParam{
// 					Error: err.Error()})
// 				conn.WriteMessage(b)
// 				return
// 			}

// 			var (
// 				allprizes           = make([]models.PrizeModel, 0)   // 發放獎品資訊
// 				allTypePrizes       = make([][]models.PrizeModel, 0) // 各獎項獎品資訊
// 				allTypePrizesAmount = make([]int64, 0)               // 各獎項獎品總數
// 				allPrizesPeople     = make([]int64, 0)               // 各獎項需要發獎的人數
// 			)

// 			// 如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量
// 			if len(scores) < int(limit) {
// 				limit = int64(len(scores))
// 				if limit <= first {
// 					first, second, third, general = limit, 0, 0, 0
// 				} else if limit > first && limit <= first+second {
// 					second, third, general = limit-first, 0, 0
// 				} else if limit > first+second && limit <= first+second+third {
// 					third, general = limit-first-second, 0
// 				} else if limit > first+second+third && limit < first+second+third+general {
// 					general = limit - first - second - third
// 				}
// 			}
// 			allPrizesPeople = []int64{first, second, third, general}
// 			// fmt.Println("如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量: ",
// 			// 	first, ", ", second, ", ", third, ", ", general)

// 			// 取得各獎項所有資訊、數量
// 			if allTypePrizes, allTypePrizesAmount, err = h.getAllTypePrizes(gameID); err != nil {
// 				b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 				conn.WriteMessage(b)
// 				return
// 			}

// 			// 判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值
// 			for i := 0; i < len(allPrizesPeople); i++ {
// 				if allTypePrizesAmount[i] < allPrizesPeople[i] {
// 					if i < len(allPrizesPeople)-1 {
// 						// 不夠發放，將多餘的值設置至下一個獎項
// 						allPrizesPeople[i+1] = allPrizesPeople[i+1] -
// 							allTypePrizesAmount[i] + allPrizesPeople[i]
// 					}
// 					allPrizesPeople[i] = allTypePrizesAmount[i]
// 				}
// 			}
// 			limit = allPrizesPeople[0] + allPrizesPeople[1] +
// 				allPrizesPeople[2] + allPrizesPeople[3]
// 			// fmt.Println("判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值: ",
// 			// 	allPrizesPeople[0], ", ", allPrizesPeople[1], ", ", allPrizesPeople[2], ", ", allPrizesPeople[3])

// 			// 將需要發放的獎品都append allprizes
// 			for i := 0; i < len(allPrizesPeople); i++ {
// 				var stopValue int64 // 當len(allprizes)=stopValue，停止該獎項發放並發放下一個獎項
// 				for m := 0; m <= i; m++ {
// 					stopValue += allPrizesPeople[m]
// 				}

// 				for _, prize := range allTypePrizes[i] {
// 					for n := 0; n < int(prize.PrizeRemain); n++ {
// 						allprizes = append(allprizes, prize)
// 						if int64(len(allprizes)) == stopValue {
// 							break
// 						}
// 					}
// 					if int64(len(allprizes)) == stopValue {
// 						break
// 					}
// 				}
// 			}

// 			// 發放獎品給中獎人員(多線程發放)
// 			var winRecords = make([]models.PrizeStaffModel, limit)
// 			wg.Add(int(limit)) //計數器
// 			for i := 0; i < int(limit); i++ {
// 				go func(i int) {
// 					defer wg.Done()
// 					// 遞減redis獎品數量
// 					h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, allprizes[i].PrizeID)

// 					// 新增中獎人員名單
// 					id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 						Add(models.NewPrizeStaffModel{
// 							ActivityID: activityID,
// 							GameID:     gameID,
// 							UserID:     scores[i].UserID,
// 							PrizeID:    allprizes[i].PrizeID,
// 							Round:      strconv.Itoa(int(round)),
// 							Status:     "no",
// 							// White:      "no",
// 						})
// 					if err != nil {
// 						b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 						conn.WriteMessage(b)
// 						return
// 					}

// 					// 用戶遊戲紀錄
// 					record := models.PrizeStaffModel{
// 						ID:            id,
// 						ActivityID:    activityID,
// 						GameID:        gameID,
// 						PrizeID:       allprizes[i].PrizeID,
// 						UserID:        scores[i].UserID,
// 						Name:          scores[i].Name,
// 						Avatar:        scores[i].Avatar,
// 						Game:          game,
// 						PrizeName:     allprizes[i].PrizeName,
// 						PrizeType:     allprizes[i].PrizeType,
// 						PrizePicture:  allprizes[i].PrizePicture,
// 						PrizePrice:    allprizes[i].PrizePrice,
// 						PrizeMethod:   allprizes[i].PrizeMethod,
// 						PrizePassword: allprizes[i].PrizePassword,
// 						Round:         round,
// 						WinTime:       now.Format("2006-01-02 15:04:05"),
// 						Status:        "no",
// 						Score:         scores[i].Score,
// 						Score2:        scores[i].Score2,
// 					}

// 					// 更新redis裡用戶遊戲紀錄
// 					// 用戶遊戲紀錄(包含中獎與未中獎)
// 					allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, scores[i].UserID)
// 					if err != nil {
// 						b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 						conn.WriteMessage(b)
// 						return
// 					}
// 					allRecords = append(allRecords, record)
// 					// allRecords = append(allRecords, models.PrizeStaffModel{
// 					// 	ID:            id,
// 					// 	ActivityID:    activityID,
// 					// 	GameID:        gameID,
// 					// 	PrizeID:       allprizes[i].PrizeID,
// 					// 	UserID:        scores[i].UserID,
// 					// 	Name:          scores[i].Name,
// 					// 	Avatar:        scores[i].Avatar,
// 					// 	Game:          game,
// 					// 	PrizeName:     allprizes[i].PrizeName,
// 					// 	PrizeType:     allprizes[i].PrizeType,
// 					// 	PrizePicture:  allprizes[i].PrizePicture,
// 					// 	PrizePrice:    allprizes[i].PrizePrice,
// 					// 	PrizeMethod:   allprizes[i].PrizeMethod,
// 					// 	PrizePassword: allprizes[i].PrizePassword,
// 					// 	Round:         round,
// 					// 	WinTime:       now.Format("2006-01-02 15:04:05"),
// 					// 	Status:        "no",
// 					// })
// 					h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
// 						scores[i].UserID, utils.JSON(allRecords))
// 					h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
// 						config.REDIS_EXPIRE)

// 					// 中獎名單
// 					winRecords[i] = record
// 					// winRecords[i] = models.PrizeStaffModel{
// 					// 	ID:            id,
// 					// 	UserID:        scores[i].UserID,
// 					// 	ActivityID:    activityID,
// 					// 	GameID:        gameID,
// 					// 	PrizeID:       allprizes[i].PrizeID,
// 					// 	Name:          scores[i].Name,
// 					// 	Avatar:        scores[i].Avatar,
// 					// 	PrizeName:     allprizes[i].PrizeName,
// 					// 	PrizeType:     allprizes[i].PrizeType,
// 					// 	PrizePicture:  allprizes[i].PrizePicture,
// 					// 	PrizePrice:    allprizes[i].PrizePrice,
// 					// 	PrizeMethod:   allprizes[i].PrizeMethod,
// 					// 	PrizePassword: allprizes[i].PrizePassword,
// 					// 	Round:         round,
// 					// 	Status:        "no",
// 					// 	Score:         scores[i].Score,
// 					// 	Score2:        scores[i].Score2,
// 					// }
// 				}(i)
// 			}
// 			wg.Wait() //等待計數器歸0

// 			// 從獎品redis(hash格式)中取得所有獎品剩餘數量，並更新至資料表中(多線呈處理)
// 			var (
// 				amount int64
// 				lock   sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
// 			)
// 			prizes, err := h.getPrizes(gameID)
// 			if err != nil {
// 				b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 				conn.WriteMessage(b)
// 				return
// 			}
// 			// 多線呈更新資料表獎品資訊
// 			wg.Add(len(prizes)) //計數器
// 			for i := 0; i < len(prizes); i++ {
// 				go func(i int) {
// 					defer wg.Done()
// 					// 計算剩餘獎品數
// 					lock.Lock() //佔有資源
// 					amount += prizes[i].PrizeRemain
// 					lock.Unlock() //釋放資源
// 					// 更新獎品剩餘數量
// 					if err = models.DefaultPrizeModel().SetDbConn(h.dbConn).
// 						UpdateRemain(
// 							prizes[i].PrizeID, prizes[i].PrizeRemain); err != nil {
// 						b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 						conn.WriteMessage(b)
// 						return
// 					}
// 				}(i)
// 			}
// 			wg.Wait() //等待計數器歸0

// 			// 快問快答結算後更新輪次以及題目進行題數資訊(清空qa_people人數)
// 			if game == "QA" {
// 				// 清空redis題目相關資訊(所有資訊歸零)
// 				h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
// 					"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

// 				// 更新題目進行題數資料(qa_round=1)
// 				if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 					SetRedisConn(h.redisConn).UpdateQARound(true, gameID, "1"); err != nil {
// 					b, _ := json.Marshal(WinningStaffParam{Error: "錯誤: 無法更新題目進行題數資料(qa_round=1)"})
// 					conn.WriteMessage(b)
// 					return
// 				}

// 				if err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 					UpdateGameStatus(true, gameID, "+1", "", "0", "calculate"); err != nil {
// 					b, _ := json.Marshal(WinningStaffParam{Error: "錯誤: 無法更新遊戲狀態"})
// 					conn.WriteMessage(b)
// 					return
// 				}
// 			}

// 			// endTime := time.Now()
// 			// log.Println("統計中獎人員花費時間: ", endTime.Sub(startTime))

// 			// 結算後刪除redis遊戲相關資訊
// 			h.redisConn.DelCache(config.SCORES_REDIS + gameID)         // 分數
// 			h.redisConn.DelCache(config.SCORES_2_REDIS + gameID)       // 第二分數
// 			h.redisConn.DelCache(config.CORRECT_REDIS + gameID)        // 答對題數
// 			h.redisConn.DelCache(config.QA_REDIS + gameID)             // 快問快答題目資訊
// 			h.redisConn.DelCache(config.QA_RECORD_REDIS + gameID)      // 快問快答答題紀錄
// 			h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID) // 中獎人員

// 			b, _ := json.Marshal(WinningStaffParam{
// 				PrizeStaffs: winRecords,
// 				Game:        GameModel{PrizeAmount: amount},
// 				Message:     true,
// 			})
// 			if err := conn.WriteMessage(b); err != nil {
// 				return
// 			}

// 			// conn.Close()
// 			// conn.isClose = true
// 			return // 停止傳送訊息
// 		}

// 		// ws關閉
// 		if conn.isClose {
// 			if game == "redpack" || game == "ropepack" ||
// 				game == "whack_mole" || game == "monopoly" {
// 				// 清除中獎人員資料
// 				h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
// 				// 清除分數資料
// 				h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
// 				h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
// 			}

// 			// h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID) // 中獎人員

// 			// // 刪除分數redis資訊
// 			// if game == "whack_mole" || game == "monopoly" {
// 			// 	h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
// 			// 	h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
// 			// }
// 			return
// 		}
// 	}
// }()
// 間隔n秒傳遞中獎訊息版本

// #####壓力測試專用開始，新增中獎人員(n名)#####
// if gameID == "bCIX7MJaFen3v4BtxP5E" && game == "redpack" {
// 	var (
// 		// wg         sync.WaitGroup                      // 宣告WaitGroup 用以等待執行序
// 		winRecords = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄
// 		record     = models.PrizeStaffModel{ID: 0, PrizeName: "未中獎",
// 			PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
// 			PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(預設未中獎)
// 		// prizeID = "usCosDb5frEeK58FAgKu"
// 		userID = "test"
// 		// n      = 5
// 		// id      int64
// 		prizeID string
// 	)
// 	// wg.Add(n) //計數器

// 	// for i := 0; i < n; i++ {
// 	// 加入判斷是否中獎相關function
// 	// go func(i int) {
// 	// defer wg.Done()

// 	// 判斷是否中獎
// 	if prizeID, err = h.getEnvelopePrize(gameID, gameModel.Percent); err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	// 中獎，更新獎品數量(id != 0)
// 	if prizeID != "" {
// 		err = models.DefaultPrizeModel().SetDbConn(h.dbConn).
// 			DecrRemain(prizeID)
// 		if err != nil {
// 			// 獎品無庫存，改為未中獎資訊
// 			prizeID = ""
// 		} else {
// 			// 遞減redis獎品數量
// 			h.redisConn.HashDecrCache(config.PRIZES_REDIS+gameID, prizeID)
// 		}
// 	} else if prizeID == "" {
// 		// fmt.Println("未中獎")
// 	}

// 	// 新增中獎人員名單
// 	if _, err = models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		Add(models.NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      strconv.Itoa(int(round)),
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		response.BadRequest(ctx, "錯誤: 新增中獎人員發生問題，請重新操作")
// 		return
// 	}

// 	if winRecords, err = h.getUserWinningRecords(gameID, userID); err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	for i := 0; i < len(winRecords); i++ {
// 		if winRecords[i].Round == int64(round) {
// 			// 該輪次遊戲中獎
// 			record = winRecords[i]
// 		}
// 	}

// 	// 依照中獎的順序將資料推送至list格式的redis中(沒有按照獎品類型排序)
// 	h.redisConn.ListRPush(config.WINNING_STAFFS_REDIS+gameID, utils.JSON(record))
// 	// 		}(i)
// }

// // 	wg.Wait() //等待計數器歸0
// // }
// #####壓力測試專用結束#####

// else if game == "QA" {
// 	h.redisConn.DelCache(config.QA_REDIS + gameID) // 快問快答題目資訊
// }
// else if game == "QA" {
// 	h.redisConn.DelCache(config.QA_REDIS + gameID) // 快問快答題目資訊
// }

// 獎品數(從redis取得，如果沒有才執行資料表查詢)
// amount, err = h.getPrizesAmount(gameID)
// if err != nil {
// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 	conn.WriteMessage(b)
// 	return
// }

// 結算後的獎品數(從redis取得，如果沒有才執行資料表查詢)
// amount, _ = h.getPrizesAmount(gameID)

// fmt.Println("i: ", i)
// 遞減獎品數量
// if err = models.DefaultPrizeModel().SetDbConn(h.dbConn).
// 	DecrRemain(allprizes[i].PrizeID); err != nil {
// 	// 獎品無庫存
// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 	conn.WriteMessage(b)
// 	return
// }

// if game == "redpack" || game == "ropepack" || game == "lottery" {
// PrizeStaffs: staffs,
// Game:        GameModel{PrizeAmount: amount},
// })
// if err := conn.WriteMessage(b); err != nil {
// 	return
// }
// } else
// var newStaffs = make([]models.PrizeStaffModel, 0)
// if game == "redpack" || game == "ropepack" {
// 	newStaffs, _ = h.getWinningRecords(gameID,
// 		"activity_staff_prize.round", round)
// } else if game == "lottery" {
// 	newStaffs, _ = h.getWinningRecordsByGameID(gameID)
// }

// 中獎人員資料
// if len(newStaffs) != int(prizeStaffsLength) {
// 	// 新增中獎人員，將新的推送至陣列
// 	for i := int(prizeStaffsLength); i < len(newStaffs); i++ {
// 		staffs = append(staffs, models.PrizeStaffModel{
// 			ID:         newStaffs[i].ID,
// 			ActivityID: newStaffs[i].ActivityID,
// 			GameID:     newStaffs[i].GameID,
// 			UserID:     newStaffs[i].UserID,
// 			PrizeID:    newStaffs[i].PrizeID,
// 			Name:       newStaffs[i].Name,
// 			Avatar:     newStaffs[i].Avatar,
// 			WinTime:    newStaffs[i].WinTime,
// 			PrizeName:  newStaffs[i].PrizeName,
// 			PrizeType:  newStaffs[i].PrizeName,
// 		})
// 	}

// 獎品數量
// if game != "lottery" {
// 	prizeLength = prizeLength - int64(len(newStaffs)) + prizeStaffsLength
// }
// prizeStaffsLength = int64(len(newStaffs))
// fmt.Println("中獎人員數量: ", prizeStaffsLength, ", 獎品數: ", int64(prizeLength))

// }
// } else if game == "whack_mole" {
// 中獎人員(從redis取得)
// staffs, _ = h.getWinningRecords(game, gameID, round, limit)

// 獎品數(從redis取得，如果沒有才執行資料表查詢)
// amount, _ = h.getPrizesCount(gameID)

// 判斷遊戲狀態(從redis取得，如果沒有才執行資料表查詢)
// gameModel = h.getGameInfo(gameID)

// if gameModel.GameStatus != "calculate" {
// 敲敲樂遊戲計算分數前n高人員(從redis取得，如果沒有才執行資料表查詢)
// staffs, _ = h.getTopScoreStaffs(gameID, round, limit)
// ; err != nil {
// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// conn.WriteMessage(b)
// 	return
// }
// fmt.Println("取得分數前n高的人員: ", limit, ", 人員: ", staffs)

// b, _ := json.Marshal(WinningStaffParam{

// 取得所有獎項設置的人數
// first, second, third, general = h.getWhackMoleInfo(gameID)
// limit = first + second + third + general
// 分數前n高人員
// fmt.Println("檢查輪次資訊是否正確: ", gameModel.Round-1)

// ; err != nil {
// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 	conn.WriteMessage(b)
// 	return
// }
// staffs, _ = h.getTopScoreStaffs(gameID, round, limit)
// prizeStaffs         = make([]models.PrizeStaffModel, 0) // 中獎人員(可能因獎品不夠，導致有些人員是沒有中獎的)
// first, second, third, general int64
// now, _              = time.ParseInLocation("2006-01-02 15:04:05",
// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// first, second, third, general, totalPrizes int64

// Name:       staffs[i].Name,
// Avatar:     staffs[i].Avatar,
// PrizeName:  allprizes[i].Name,
// PrizeType:  allprizes[i].PrizeType,
// Picture:    allprizes[i].Picture,
// Price:      strconv.Itoa(int(allprizes[i].Price)),
// Method:     allprizes[i].Method,
// Password:   allprizes[i].Password,
// WinTime:    now.Format("2006-01-02 15:04:05"),
// WinTime:    now.Format("2006-01-02 15:04:05"),

// if prizeStaffs, err = h.getWinningRecordsOrderByScore(gameID,
// 	strconv.Itoa(int(round))); err != nil {
// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 	conn.WriteMessage(b)
// 	return
// }
//  else if len(prizeStaffs) == 0 {
// 	prizeStaffs = make([]models.PrizeStaffModel, 0)
// }

// 計算中獎人員、分數前n高人員
// 注意: 遊戲已開始，後端輪次已更新，要查詢該輪的輪次資訊需要使用gameModel.Round - 1
// if game == "redpack" || game == "ropepack" || game == "lottery" {
// 	// 中獎人員(紅包遊戲)
// 	if staffs, err = h.getPrizeStaffs(gameID, []string{"round", "method"},
// 		[]interface{}{round, "thanks"}); err != nil {
// 		b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 		conn.WriteMessage(b)
// 		return
// 	}
// 	prizeStaffsLength = int64(len(staffs))
// } else if game == "whack_mole" {

// if game == "whack_mole" {
// 	// 取得所有獎項設置的人數
// 	// gameModel := h.getGameInfo(gameID)
// 	first, second, third, general = gameModel.FirstPrize, gameModel.SecondPrize, gameModel.ThirdPrize, gameModel.GeneralPrize
// 	limit = first + second + third + general

// 	// 敲敲樂遊戲計算分數前n高人員
// 	// if staffs, err = h.getTopScoreStaffs(gameID, round, limit); err != nil {
// 	// 	b, _ := json.Marshal(WinningStaffParam{Error: err.Error()})
// 	// 	conn.WriteMessage(b)
// 	// 	return
// 	// }
// 	// fmt.Println("取得分數前n高的人員: ", limit, ", 人員: ", staffs)
// }

// time.Sleep(1 * time.Second)
// b, _ := json.Marshal(WinningStaffParam{
// 	PrizeStaffs: staffs,
// 	Game:        GameModel{Amount: prizeLength},
// })
// fmt.Println("中獎人員數量: ", prizeStaffsLength, ", 獎品數: ", int64(prizeLength))
// if err := conn.WriteMessage(b); err != nil {

// 比較賓果人數與每輪發獎人數
// if int64(len(bingos)) > gameModel.RoundPrize {
// 	// fmt.Println("賓果人數>每輪發獎人數")
// 	// 賓果人數>每輪發獎人數，判斷中獎人員與PK人員名單
// 	for _, bingo := range bingos {
// 		// 判斷賓果的回合數
// 		if bingo.Score2 == float64(gameModel.BingoRound) {
// 			// 最後一回合賓果的人員為PK人員
// 			pkUsers = append(pkUsers, bingo)
// 		} else if bingo.Score2 < float64(gameModel.BingoRound) {
// 			// 確定中獎人員
// 			winUsers = append(winUsers, bingo)
// 		}
// 	}
// 	// fmt.Println("賓果人數>每輪發獎人數: ", len(winUsers), len(pkUsers))
// } else {
// fmt.Println("賓果人數<=每輪發獎人數")
// 賓果人數<=每輪發獎人數，所有人員都有獎
// winUsers = bin
