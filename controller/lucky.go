package controller

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"hilive/modules/utils"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// var (
// 	luckyLock sync.Mutex
// )

// PrizesParam 中獎紀錄、是否中獎參數
type PrizesParam struct {
	// IsWin       bool `json:"is_win" example:"true"` // 是否獲勝
	Game        GameModel
	PrizeStaff  models.PrizeStaffModel
	PrizeStaffs []models.PrizeStaffModel
	Error       string `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 判斷是否中獎，回傳中獎紀錄及抽獎紀錄(玩家端)
// @Tags Lucky
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param session_id query string true "session ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, lottery, monopoly, QA, tugofwar, bingo, 3DGachaMachine)
// @Param mode query string false "game" Enums(order)
// @Success 200 {array} PrizesParam
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /game/lucky [get]
func (h *Handler) GetLucky(ctx *gin.Context) {
	var (
		host        = ctx.Request.Host
		activityID  = ctx.Query("activity_id")
		gameID      = ctx.Query("game_id")
		game        = ctx.Query("game")
		sessionID   = ctx.Query("session_id")
		mode        = ctx.Query("mode")
		gameModel   = h.getGameInfo(gameID, game)
		allRecords  = make([]models.PrizeStaffModel, 0) // 用戶所有遊戲紀錄
		gameRecords = make([]models.PrizeStaffModel, 0) // 用戶遊戲紀錄
		winRecords  = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄
		prize       = models.PrizeModel{ID: 0, PrizeName: "未中獎",
			PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
			PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(預設未中獎)
		record = models.PrizeStaffModel{ID: 0, PrizeName: "未中獎",
			PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
			PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(預設未中獎)
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
		user    models.LineModel
		team    string
		prizeID string
		round   = int(gameModel.GameRound) - 1
		id      int64
		err     error
	)

	// 競技遊戲不需要互斥鎖，因為競技遊戲時該API是查詢用戶中獎資訊，不是發獎
	// if game == "redpack" || game == "ropepack" ||
	// game == "lottery" || game == "3DGachaMachine" {
	// luckyLock.Lock()
	// defer luckyLock.Unlock()
	// }

	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if game == "" || activityID == "" || gameID == "" || gameModel.ID == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID",
		})
		return
	}

	// 順序抽獎不需要取得用戶資訊、不需要取得用戶遊戲紀錄
	if !(game == "lottery" && mode == "order") {
		// 利用session值取得用戶資料
		if user, _, err = h.getUser(sessionID, gameModel.UserID); err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}

		// 用戶遊戲紀錄(包含中獎與未中獎)
		allRecords, _, _, gameRecords, winRecords, err = h.getUserGameRecords(activityID, gameID, game, user.UserID)
		if err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	}

	// 抽獎判斷
	// 競技遊戲不需要互斥鎖，因為競技遊戲時該API是查詢用戶中獎資訊，不是發獎
	if game == "redpack" || game == "ropepack" || game == "lottery" || game == "3DGachaMachine" {
		for l := 0; l < MaxRetries; l++ {
			// 上鎖，抽獎
			ok, _ := h.acquireLock(config.LUCKY_LOCK_REDIS+gameID, LockExpiration)
			if ok == "OK" {
				// 釋放鎖
				// defer h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)

				// 隨機抽獎
				if game == "redpack" || game == "ropepack" || game == "3DGachaMachine" {
					if game == "3DGachaMachine" {
						// 扭蛋機沒有輪次
						round = len(gameRecords) + 1 // 第n次抽獎
					}

					// 判斷是否中獎
					if prizeID, err = h.getEnvelopePrize(gameID, gameModel.Percent); err != nil {
						response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  "",
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: err.Error(),
						})

						// 釋放鎖
						h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)
						return
					}
				} else if game == "lottery" {
					// 遊戲抽獎
					// 遊戲抽獎的round為第n次抽獎，先取得用戶的抽獎紀錄後再計算round參數
					if mode == "" {
						round = len(gameRecords) + 1 // 第n次抽獎
					}

					// 判斷是否中獎
					if prizeID, err = h.getLotteryPrize(gameID); err != nil {
						response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  "",
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: err.Error(),
						})

						// 釋放鎖
						h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)
						return
					}
				}

				// 獎品處理
				if game == "lottery" || game == "3DGachaMachine" {
					// 遊戲抽獎跟扭蛋機抽到獎時立即處理獎品數(資料表、redis)
					// 中獎，更新獎品數量(id != 0)
					if prizeID != "" {
						err = models.DefaultPrizeModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							DecrRemain(gameID, prizeID)
						if err != nil {
							// 獎品無庫存，改為未中獎資訊
							prizeID = ""
						}

						// 獎品資訊(redis處理)
						prize, err = models.DefaultPrizeModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							FindPrize(true, prizeID)
						if err != nil {
							response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
								UserID:  "",
								Method:  ctx.Request.Method,
								Path:    ctx.Request.URL.Path,
								Message: err.Error(),
							})

							// 釋放鎖
							h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)
							return
						}

					}
					// else if prizeID == "" {
					// }

				} else if game == "redpack" || game == "ropepack" {
					// 紅包遊戲抽到獎時不處理獎品數(資料表)，結算統一整理

					// 中獎，更新獎品數量(peize_id != "")
					if prizeID != "" {
						// 遞減redis獎品數量，結算時才將剩餘獎品數輛寫入資料表
						h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

						// 設置過期時間
						// h.redisConn.SetExpire(config.GAME_PRIZES_AMOUNT_REDIS+gameID, config.REDIS_EXPIRE)

						// 獎品資訊(redis處理)
						prize, err = models.DefaultPrizeModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							FindPrize(true, prizeID)
						if err != nil {
							response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
								UserID:  "",
								Method:  ctx.Request.Method,
								Path:    ctx.Request.URL.Path,
								Message: err.Error(),
							})

							// 釋放鎖
							h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)
							return
						}
					}
				}

				// 釋放鎖
				h.releaseLock(config.LUCKY_LOCK_REDIS + gameID)
				break
			}

			// 鎖被佔用，稍微延遲後重試
			time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
		}
	}

	// 獎品、中獎人員資料處理
	if game == "redpack" || game == "ropepack" || game == "lottery" || game == "3DGachaMachine" {
		// 順序抽獎不需要新增中獎紀錄(只須扣除獎品數量)
		// 紅包遊戲不需要新增中獎紀錄(結算時統一寫入)
		// 扭蛋機跟遊戲抽獎同時抽獎模式下新增中獎紀錄(資料表)
		if !(game == "lottery" && mode == "order") && game != "redpack" && game != "ropepack" {
			// 新增中獎人員
			if id, err = models.DefaultPrizeStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(models.NewPrizeStaffModel{
					UserID:     user.UserID,
					ActivityID: activityID,
					GameID:     gameID,
					Game:       game,
					PrizeID:    prizeID,
					Round:      strconv.Itoa(round),
					Status:     "no",
					// White:      "no",
					Score:  0,
					Score2: 0,
					Rank:   0,
				}); err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 新增中獎人員發生問題，請重新操作",
				})
				return
			}
		}

		// 遊戲紀錄(前端利用method判斷是否中獎)
		record = models.PrizeStaffModel{
			ID:            id,
			ActivityID:    activityID,
			GameID:        gameID,
			PrizeID:       prize.PrizeID,
			UserID:        user.UserID,
			Name:          user.Name,
			Avatar:        user.Avatar,
			Game:          game,
			PrizeName:     prize.PrizeName,
			PrizeType:     prize.PrizeType,
			PrizePicture:  prize.PrizePicture,
			PrizePrice:    prize.PrizePrice,
			PrizeMethod:   prize.PrizeMethod,
			PrizePassword: prize.PrizePassword,
			Round:         int64(round),
			WinTime:       now.Format("2006-01-02 15:04:05"),
			Status:        "no",
		}

		// 順序抽獎不需要更新redis用戶遊戲紀錄
		// 紅包遊戲需要更新redis用戶遊戲紀錄
		// 扭蛋機跟遊戲抽獎同時抽獎模式下需要更新redis用戶遊戲紀錄
		if !(game == "lottery" && mode == "order") {
			// 更新redis裡用戶遊戲紀錄
			allRecords = append(allRecords, record)

			h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
				user.UserID, utils.JSON(allRecords))
			// 設置過期時間
			// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)
		}

		// 將中獎紀錄加入參數中
		if prizeID != "" {
			winRecords = append(winRecords, record)
		}
	}

	// 競技遊戲才需要判斷該輪次是否中獎
	if game == "whack_mole" || game == "monopoly" || game == "QA" || game == "bingo" {
		// 其他競技遊戲
		for i := 0; i < len(winRecords); i++ {
			if winRecords[i].Round == int64(round) {
				record = winRecords[i]
				break
			}
		}
	} else if game == "tugofwar" {
		// 第三種方式，不管什麼發獎方式，查詢該輪次遊戲中獎資訊並可以取得獲勝隊伍
		// 取得雙方隊伍資訊
		teamModel, _ := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindTeam(true, gameID)

		// 第二種方式: 取得玩家隊伍，從redis中抓取獲勝隊伍
		if utils.InArray(teamModel.LeftTeamPlayers, user.UserID) {
			team = "left_team"
		} else if utils.InArray(teamModel.RightTeamPlayers, user.UserID) {
			team = "right_team"
		}

		// log.Println("獲勝隊: ", teamModel.WinTeam, teamModel.LeftTeamPlayers, teamModel.RightTeamPlayers)
		// log.Println("哪隊? ", team)

		// 判斷玩家隊伍是否獲勝
		if team == teamModel.WinTeam {
			// 獲勝
			record = models.PrizeStaffModel{ID: 0, PrizeName: "中獎",
				PrizeType: "first", PrizePassword: "no", PrizeMethod: "first",
				PrizePicture: "/admin/uploads/system/img-prize-pic.png"} // 抽獎紀錄(中獎)
		}
	}

	// 將中獎資訊加入redis中(競技遊戲與遊戲抽獎的順序抽獎不需要加入redis中)
	if game == "redpack" || game == "ropepack" ||
		(game == "lottery" && mode == "") {
		if record.PrizeType != "thanks" && record.PrizeMethod != "thanks" {
			// 依照中獎的順序將資料推送至list格式的redis中(沒有按照獎品類型排序)
			h.redisConn.ListRPush(config.WINNING_STAFFS_REDIS+gameID, utils.JSON(record))

			// 設置過期時間
			// h.redisConn.SetExpire(config.WINNING_STAFFS_REDIS+gameID, config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
		} else if (game == "redpack" || game == "ropepack") && (record.PrizeType == "thanks" || record.PrizeMethod == "thanks") {
			// 紅包遊戲，未中獎時，將未中獎紀錄寫入redis中
			h.redisConn.ListRPush(config.NO_WINNING_STAFFS_REDIS+gameID, utils.JSON(record))

			// 設置過期時間
			// h.redisConn.SetExpire(config.NO_WINNING_STAFFS_REDIS+gameID, config.REDIS_EXPIRE)
		}
	}

	// 扭蛋機遊戲
	if game == "3DGachaMachine" {
		// 不管是否中獎都將遊戲紀錄資訊加入redis中(winning_staffs_)
		h.redisConn.ListRPush(config.WINNING_STAFFS_REDIS+gameID, utils.JSON(record))

		// 設置過期時間
		// h.redisConn.SetExpire(config.WINNING_STAFFS_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
	}

	// 回傳訊息
	response.OkWithData(ctx, PrizesParam{
		PrizeStaffs: winRecords,
		PrizeStaff:  record,
		Game:        GameModel{},
	})
}

// getEnvelopePrize 紅包遊戲抽獎
func (h *Handler) getEnvelopePrize(gameID string, percent int64) (prizeID string, err error) {
	var (
		prizes                          = make([]models.PrizeModel, 0)
		prizeIDs                        = make([]string, 0)
		random, prizeTyepLength, sample int
	)

	// 獎品資訊(redis)
	if prizes, err = h.getPrizes(gameID); err != nil {
		return
	}

	// 處理獎品資訊(獎品數量大於0)
	for _, prize := range prizes {
		if prize.PrizeRemain > 0 {
			prizeIDs = append(prizeIDs, prize.PrizeID)
			prizeTyepLength++ // 獎品種類長度
		}
	}
	if prizeTyepLength == 0 {
		return
	}

	sample = (prizeTyepLength * 100 / int(percent)) // 樣本數
	// 避免獎品種類等於樣本數的情況下中獎機率會百分之百(機率不為100%)
	if percent < 100 && sample == prizeTyepLength {
		sample += 1
	}

	// 隨機抽獎
	rand.Seed(time.Now().UnixNano())
	random = rand.Intn(sample)

	// 未中獎
	if random >= prizeTyepLength {
		return
	}

	// 中獎，回傳獎品ID
	return prizeIDs[random], nil
}

// getLotteryPrize 遊戲抽獎中獎判斷
func (h *Handler) getLotteryPrize(gameID string) (prizeID string, err error) {
	var (
		prizes   = make([]models.PrizeModel, 0) // 包含剩餘數為0的獎品
		PrizeIDs = make([]string, 0)            // 數量大於0的獎品資訊(補足未中獎資訊)
		length   int64
	)

	// 所有獎品資訊(不管剩餘數量是否為0)、獎品總數
	if prizes, err = h.getPrizes(gameID); err != nil {
		return
	}

	// 清除沒有剩餘的獎品資訊並補齊需要的獎品數量
	for i := 0; i < len(prizes); i++ {
		if prizes[i].PrizeRemain > 0 {
			PrizeIDs = append(PrizeIDs, prizes[i].PrizeID)
		}
	}
	// 補未中獎資訊
	length = int64(8 - len(prizes)) // 需要補足的獎品數
	for i := 0; i < int(length); i++ {
		PrizeIDs = append(PrizeIDs, "")
	}

	// 隨機抽獎
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(len(PrizeIDs))
	prizeID = PrizeIDs[random]

	return
}

// lockLucky 上鎖，抽獎時需要使用，避免用戶併發多次抽獎
// func lockLucky(redis cache.Connection, id string) bool {
// 	return lockLuckyServ(redis, id)
// }

// UnLockLucky 解鎖，抽獎時需要使用，避免用戶併發多次抽獎
// func unLockLucky(redis cache.Connection, id string) bool {
// 	return unlockLuckyServ(redis, id)
// }

// getLuckyLockKey 取得緩存名稱
// func getLuckyLockKey(id string) string {
// 	return fmt.Sprintf("lucky_lock_%s", id)
// }

// lockLuckyServ 上鎖
// func lockLuckyServ(redis cache.Connection, id string) bool {
// 	key := getLuckyLockKey(id)
// 	if rs, _ := redis.SetCache(key, 1, "EX", 3, "NX"); rs == "OK" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// unlockLuckyServ 解鎖
// func unlockLuckyServ(redis cache.Connection, id string) bool {
// 	key := getLuckyLockKey(id)
// 	if rs, _ := redis.DelCache(key); rs == "OK" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// 取得該輪次拔河遊戲中獎紀錄
// prizeStaffModel, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	FindAll(activityID, gameID, "", game, strconv.Itoa(round), "", "win",
// 		1, 0)
// if err != nil {
// 	response.BadRequest(ctx, err.Error())
// 	return
// }
// if len(prizeStaffModel) > 0 {
// 	winTeam = prizeStaffModel[0].Team
// }
// fmt.Println("獲勝方資訊: ", prizeStaffModel)
// fmt.Println("獲勝方隊伍: ", winTeam)

// 拔河遊戲，判斷發獎方式(第一、二種方式)
// if gameModel.Prize == "uniform" {
// 	fmt.Println("測試1，統一發放獎品")
// 	// 統一發放
// 	var (
// 		teamModel = h.getTeamInfo(gameID) // 隊伍資訊
// 		// leaderID  string
// 		team string
// 	)

// fmt.Println("左方隊伍玩家: ", teamModel.LeftTeamPlayers)
// fmt.Println("右方隊伍玩家: ", teamModel.RightTeamPlayers)
// fmt.Println("玩家ID: ", user.UserID)

// 	// 第一種方式: 取得玩家的隊伍，判斷隊長是否有中獎紀錄
// 	// 取得用戶隊伍方的隊長資訊
// 	// 判斷用戶是否為左方隊伍
// 	// 	for _, id := range teamModel.LeftTeamPlayers {
// 	// 		if id == user.UserID {
// 	// 			leaderID = teamModel.LeftTeamLeader
// 	// 			break
// 	// 		}
// 	// 	}
// 	// 	// 判斷用戶是否為右方隊伍
// 	// 	if leaderID == "" {
// 	// 		for _, id := range teamModel.RightTeamPlayers {
// 	// 			if id == user.UserID {
// 	// 				leaderID = teamModel.RightTeamLeader
// 	// 				break
// 	// 			}
// 	// 		}
// 	// 	}

// 	// 	// 隊長遊戲紀錄(包含中獎與未中獎)
// 	// 	_, _, _, _, leaderWinRecords, err := h.getUserGameRecords(activityID, gameID, game, leaderID)
// 	// 	if err != nil {
// 	// 		response.BadRequest(ctx, err.Error())
// 	// 		return
// 	// 	}

// 	// 	for i := 0; i < len(leaderWinRecords); i++ {
// 	// 		if leaderWinRecords[i].Round == int64(round) {
// 	// 			record = leaderWinRecords[i]
// 	// 		}
// 	// 	}
// } else if gameModel.Prize == "all" {
// 	fmt.Println("測試2，全部發放獎品")

// 	// 全部發放(與其他競技類遊戲相同)
// 	for i := 0; i < len(winRecords); i++ {
// 		if winRecords[i].Round == int64(round) {
// 			record = winRecords[i]
// 			isWin = true
// 			break
// 		}
// 	}
// }

// if game == "lottery" && (mode != "order" && mode != "simultaneous") {
// 	response.BadRequest(ctx, "錯誤: 無法判斷遊戲抽獎規則資訊，請輸入有效的抽獎規則")
// 	return
// }

// 用戶中獎紀錄(順序抽獎不需要取得用戶中獎紀錄)
// if !(game == "lottery" && mode == "order") {
// 	if winRecords, err = h.getUserWinningRecords(gameID, user.UserID); err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// }

// 舊版: 依照獎品類型分類至redis中(zset)
// if game == "redpack" || game == "ropepack" ||
// 	(game == "lottery" && mode == "simultaneous") {
// 	var prizeType int64
// 	if record.PrizeType == "first" {
// 		prizeType = 1
// 	} else if record.PrizeType == "second" {
// 		prizeType = 2
// 	} else if record.PrizeType == "third" {
// 		prizeType = 3
// 	} else if record.PrizeType == "general" {
// 		prizeType = 4
// 	}

// 	if prizeType > 0 && prizeType <= 4 {
// 		recordJson := utils.JSON(record) // 中獎人員json編碼資料
// 		// 將中獎人員加入redis
// 		h.redisConn.ZSetAdd(config.WINNING_STAFFS_REDIS+gameID, recordJson, prizeType)
// 	}
// }

// 用户抽獎分布式鎖
// ok := lockLucky(h.redisConn, user.UserID)
// if ok {
// 	defer unLockLucky(h.redisConn, sessionID)
// } else {
// 	fmt.Println("這裡有錯? ")
// 	// response.BadRequest(ctx, "錯誤: 用戶正在抽獎，請稍後再試")
// }

// if game == "lottery" {
// 	response.BadRequest(ctx, "錯誤: 已無獎品，無法繼續抽獎")
// 	return
// }

// prize = models.PrizeModel{ID: 0, PrizeName: "未中獎",
// 	PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
// 	PrizePicture: "/admin/uploads/system/img-prize-pic.png"}
//  else {
// 	// 遞減redis獎品數量
// 	h.redisConn.HashDecrCache(config.PRIZES_REDIS+gameID, prizeID)
// }

// if game == "QA" {
// 	round = 1 // 快問快答沒有輪次分別(取round = 1資料)
// }

// PrizeModel 獎品訊息
// type PrizeModel struct {
// 	ID            int64  `json:"id" example:"1"`
// 	ActivityID    string `json:"activity_id" example:"activity_id"`
// 	PrizeID       string `json:"prize_id" example:"prize_id"`
// 	PrizeName     string `json:"prize_name" example:"獎品名稱"`
// 	PrizeType     string `json:"prize_type" example:"獎品類型"`
// 	PrizePicture  string `json:"prize_picture" example:"獎品圖片"`
// 	PrizeAmount   int64  `json:"prize_amount" example:"5"`
// 	PrizeRemain   int64  `json:"prize_remain" example:"5"`
// 	PrizePrice    int64  `json:"prize_price" example:"1000"`
// 	PrizeMethod   string `json:"prize_method" example:"兌換方式"`
// 	PrizePassword string `json:"prize_password" example:"兌獎密碼"`
// }

// now, _         = time.ParseInLocation("2006-01-02 15:04:05",
// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)

// if game == "redpack" {
// 	err = models.DefaultRedpackPrizeModel().SetDbConn(h.dbConn).
// 		DecrRemain(prize.PrizeID, "1")
// } else if game == "ropepack" {
// 	err = models.DefaultRopepackPrizeModel().SetDbConn(h.dbConn).
// 		DecrRemain(prize.PrizeID, "1")
// } else if game == "lottery" {
// 	err = models.DefaultLotteryPrizeModel().SetDbConn(h.dbConn).
// 		DecrRemain(prize.PrizeID, "1")
// }
// // 獎品無庫存，改為未中獎資訊
// if err != nil {
// 	fmt.Println("無庫存")
// 	fmt.Println(err)
// 	if game == "lottery" {
// 		response.BadRequest(ctx, "錯誤: 已無獎品，無法繼續抽獎")
// 		return
// 	}
// 	prize = PrizeModel{ID: 0, PrizeName: "未中獎",
// 		PrizeType: "thanks", PrizePassword: "no", PrizeMethod: "thanks",
// 		PrizePicture: "/admin/uploads/system/img-prize-pic.png"}
// }

// Name:       user.Name,
// Avatar:     user.Avatar,
// PrizeName:  prize.Name,
// PrizeType:  prize.PrizeType,
// Picture:    prize.Picture,
// Price:      strconv.Itoa(int(prize.Price)),
// Method:     prize.Method,
// Password:   prize.Password,
// WinTime:    now.Format("2006-01-02 15:04:05"),

// if prizeStaffs, err = h.getWinningRecords(gameID, "activity_staff_prize.user_id", user.UserID); err != nil {
// 	response.BadRequest(ctx, err.Error())
// 	return
// }

// 抽獎紀錄(前端利用id判斷是否中獎)
// prizeStaff = PrizeStaffModel{
// 	ID:         id,
// 	ActivityID: activityID,
// 	GameID:     gameID,
// 	UserID:     user.UserID,
// 	Name:       user.Name,
// 	Avatar:     user.Avatar,
// 	PrizeName:  prize.PrizeName,
// 	Picture:    prize.Picture,
// 	Price:      prize.Price,
// 	Method:     prize.Method,
// 	Password:   prize.Password,
// 	Round:      gameModel.Round - 1,
// 	WinTime:    now.Format("2006-01-02 15:04:05"),
// 	Status:     "no",
//
