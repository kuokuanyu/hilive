package controller

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var (
// staffLock sync.Mutex
)

// @Summary 處理遊戲人員資訊(玩家端)，回傳用戶中獎紀錄、是否玩過該輪次遊戲
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, monopoly, QA, tugofwar, bingo)
// @param body body UserGameParam true "user、game param"
// @Success 200 {array} PrizesParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/staff [get]
func (h *Handler) GameStaffWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            UserGameParam
		userID string
	)
	// log.Println("開啟處理遊戲人員資訊(玩家端)ws", activityID, gameID, game)
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(PrizesParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			// gameStaffModel              = models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn)
			result                                  UserGameParam
			record                                  = models.PrizeStaffModel{ID: 0}     // id為0代表無該輪次遊戲紀錄
			gameRecords                             = make([]models.PrizeStaffModel, 0) // 用戶遊戲紀錄(包含中獎與未中獎)
			winRecords                              = make([]models.PrizeStaffModel, 0) // 用戶中獎紀錄
			guestNumbers                            = make([]int64, 0)                  // 玩家賓果號碼排序
			numbers                                 = make([]int64, 0)                  // 賓果遊戲抽球號碼
			score, rank, round, correct, bingoRound int64
			score2                                  float64
			// goingBingo                              bool
		)

		data, err := conn.ReadMessage()

		gameModel := h.getGameInfo(gameID, game) // 遊戲資訊(redis處理)
		// 遊戲資訊
		if gameModel.ID == 0 {
			b, _ := json.Marshal(PrizesParam{Error: "錯誤: 無法辨識遊戲資訊"})
			conn.WriteMessage(b)
			return
		}

		// 如果遊戲狀態是start or gaming時，輪次參數需要減1(因為start狀態時輪次已加1)
		if (game == "redpack" || game == "ropepack" ||
			game == "whack_mole" || game == "monopoly" ||
			game == "tugofwar" || game == "bingo") &&
			(gameModel.GameStatus == "start" || gameModel.GameStatus == "gaming") {
			round = gameModel.GameRound - 1
		} else {
			round = gameModel.GameRound
		}

		if err != nil {
			// log.Println("關閉處理遊戲人員資訊(玩家端)ws")

			// 關閉玩家端遊戲人員ws，清除遊戲人員資料(資料表.redis)、更新遊戲人數(資料表.redis)、隊伍資料處理(redis)
			h.closeGameStaffWebsocket(gameID, userID, game,
				round, gameModel)

			conn.Close()
			return
		}

		json.Unmarshal(data, &result)
		// fmt.Println("result參數: ", result.User)
		// if result.User.UserID == "" {
		// 	b, _ := json.Marshal(PrizesParam{
		// 		Error: "錯誤: 無法辨識用戶資訊"})
		// 	conn.WriteMessage(b)
		// 	return
		// }

		// log.Println("收到玩家端加入遊戲訊息: ", result)

		if result.User.UserID != "" {
			// log.Println("執行加入遊戲")
			userID = result.User.UserID

			if game == "tugofwar" {
				if result.User.Team != "left_team" && result.User.Team != "right_team" {
					// fmt.Println("沒接收到隊伍參數?")
					// fmt.Println("隊伍參數: ", result.User.Team)
					b, _ := json.Marshal(PrizesParam{
						PrizeStaffs: winRecords,
						Error:       "錯誤: 無法辨識隊伍資訊"})
					conn.WriteMessage(b)
					return
				}

				if result.User.Action != "join" && result.User.Action != "change" {
					b, _ := json.Marshal(PrizesParam{
						PrizeStaffs: winRecords,
						Error:       "錯誤: 無法辨識動作資訊"})
					conn.WriteMessage(b)
					return
				}
			}

			// 判斷遊戲玩家資訊(用戶遊戲紀錄、中獎紀錄、黑名單、是否允許重複中獎)
			winRecords, gameRecords, err = h.checkGameUser(activityID, gameID, result.User.UserID, game, gameModel)
			if err != nil {
				b, _ := json.Marshal(PrizesParam{
					PrizeStaffs: winRecords,
					Error:       err.Error()})
				conn.WriteMessage(b)
				// return
			}

			// 判斷是否有該輪次遊戲紀錄
			if err == nil {
				for i := 0; i < len(gameRecords); i++ {
					if gameRecords[i].Round == round {
						record = gameRecords[i]
						// log.Println("該輪次已玩過", record)
						err = errors.New("錯誤: 該輪次已玩過，無法參加本輪遊戲。請等待主辦方開啟新的一輪遊戲")
						b, _ := json.Marshal(PrizesParam{
							PrizeStaffs: winRecords,
							Error:       err.Error()})
						conn.WriteMessage(b)
						// return
					}
				}
			}

			// 當快問快答遊戲進入答題狀態時，玩家不能執行加入遊戲(避免出現一直答同一題的情況)
			if err == nil {
				if game == "QA" && gameModel.GameStatus == "gaming" {
					err = errors.New("錯誤: 快問快答遊戲已進入答題狀態，請在下一道題目開始之前加入遊戲")
					b, _ := json.Marshal(PrizesParam{
						PrizeStaffs: winRecords,
						Error:       err.Error()})
					conn.WriteMessage(b)
				}
			}

			// 當拔河遊戲已開始時(start、gaming)，玩家不能加入遊戲(已加入過的玩家可以)
			if err == nil {
				if game == "tugofwar" &&
					(gameModel.GameStatus == "start" || gameModel.GameStatus == "gaming") {
					// log.Println("玩家重新加入遊戲，原始隊伍為: ", result.User.Team)

					var key string // redis key
					if result.User.Team == "left_team" {
						key = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲左隊人數資訊，SET
					} else if result.User.Team == "right_team" {
						key = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
					}

					// 判斷是否有玩家加入遊戲的紀錄
					if !h.redisConn.SetIsMember(key, userID) {
						// log.Println("沒有玩家資料，無法重新進入拔河遊戲")
						err = errors.New("錯誤: 拔河遊戲已開始，無法參加本輪遊戲。請等待主辦方開啟新的一輪遊戲")
						b, _ := json.Marshal(PrizesParam{
							PrizeStaffs: winRecords,
							Error:       err.Error()})
						conn.WriteMessage(b)
					}
				}
			}

			// 當賓果遊戲已開始時(start、gaming)，玩家不能加入遊戲(已加入過的玩家可以)
			if err == nil {
				if game == "bingo" &&
					(gameModel.GameStatus == "start" || gameModel.GameStatus == "gaming") {

					// 判斷是否有玩家加入遊戲的紀錄
					if !h.redisConn.SetIsMember(config.GAME_ATTEND_REDIS+gameID, userID) {
						// log.Println("沒有玩家資料，無法重新進入賓果遊戲")
						err = errors.New("錯誤: 賓果遊戲已開始，無法參加本輪遊戲。請等待主辦方開啟新的一輪遊戲")
						b, _ := json.Marshal(PrizesParam{
							PrizeStaffs: winRecords,
							Error:       err.Error()})
						conn.WriteMessage(b)
					}
				}
			}

			// 當遊戲進入結算或結束狀態時，玩家不能執行加入遊戲
			if err == nil {
				if gameModel.GameStatus == "calculate" || gameModel.GameStatus == "end" {
					err = errors.New("錯誤: 遊戲已進入結算或結束狀態，請等待主辦方開啟新的一輪遊戲")
					b, _ := json.Marshal(PrizesParam{
						PrizeStaffs: winRecords,
						Error:       err.Error()})
					conn.WriteMessage(b)
				}
			}

			if err == nil {
				// 加入遊戲，處理遊戲人員資訊(redis)
				err = h.addGameStaff(activityID, gameID, result.User.UserID, game,
					result.User.Team, round, result.User.Action)
				if err != nil {
					b, _ := json.Marshal(PrizesParam{
						PrizeStaffs: winRecords,
						Error:       err.Error()})
					conn.WriteMessage(b)

					if err.Error() != "錯誤: 該隊伍遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝" &&
						err.Error() != "錯誤: 遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝" {
						// 如果是其他錯誤，直接關閉
						return
					}
				}
			}

			if err == nil {
				// 快問快答需要更新redis中的qa_people人數資訊(不管是否已經加入過遊戲，只要加入遊戲就遞增qa_people)
				if game == "QA" {
					// fmt.Println("快問快答遞增qa_people人數")
					models.DefaultGameModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						UpdateQAPeople(true, gameID, 1)
				}

				// 分數資料處理
				if game == "monopoly" || game == "whack_mole" || game == "QA" ||
					game == "tugofwar" || game == "bingo" {
					score = h.redisConn.ZSetIntScore(config.SCORES_REDIS+gameID, result.User.UserID)      // 分數資料
					score2 = h.redisConn.ZSetFloatScore(config.SCORES_2_REDIS+gameID, result.User.UserID) // 第二分數資料
					correct = h.redisConn.ZSetIntScore(config.CORRECT_REDIS+gameID, result.User.UserID)   // 答對題數
					rank = h.redisConn.ZSetRevRank(config.SCORES_REDIS+gameID, result.User.UserID)        // 排名資料

					if game == "tugofwar" && score == 0 {
						// 加入分數資料
						h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, result.User.UserID, score)

						// 設置過期時間
						// h.redisConn.SetExpire(config.SCORES_REDIS+gameID, config.REDIS_EXPIRE)

						// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
						h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
					}

					if game == "bingo" {
						// 完成賓果的回合數
						bingoRound = h.redisConn.ZSetIntScore(config.GAME_BINGO_USER_REDIS+gameID, result.User.UserID)

						// 玩家號碼排序
						guestNumbers, _ = models.DefaultGameModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							FindUserNumber(true, gameID, result.User.UserID)

							// 查詢賓果遊戲所有抽獎號碼(redis)
						numbers, _ = models.DefaultGameModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							FindBingoNumbers(true, gameID)

						// 玩家是否即將賓果
						// goingBingoStr, _ := h.redisConn.HashGetCache(config.GAME_BINGO_USER_GOING_BINGO+gameID, result.User.UserID)
						// if goingBingoStr == "true" {
						// 	goingBingo = true
						// }
						// fmt.Println("玩家的號碼排序: ", numbers)
					}
				}

				// log.Println("加入遊戲時接收到的分數資料: ", score, "輪次: ", gameModel.GameRound)

				// log.Println("controller package GameStaffWebsocket function回傳中獎紀錄: ", winRecords)
				// 回傳中獎紀錄、是否有參與該輪次遊戲紀錄
				b, _ := json.Marshal(PrizesParam{
					PrizeStaffs: winRecords,
					PrizeStaff:  record,
					Game: GameModel{
						GameScore:      score,                // 分數
						GameScore2:     score2,               // 第二分數
						GameCorrect:    correct,              // 答對題數
						GameRank:       rank + 1,             // 排名
						QARound:        gameModel.QARound,    // 題目進行題數
						BingoRound:     gameModel.BingoRound, // 賓果進行回合數
						UserBingoRound: bingoRound,           // 完成賓果的回合數
						// GoingBingo:     goingBingo,           // 玩家是否即將賓果
						Numbers:      numbers,      // 賓果遊戲抽球號碼
						GuestNumbers: guestNumbers, // 玩家賓果號碼排序
					},
				})
				conn.WriteMessage(b)
				// return
			}
		}
	}
}

// gameStaff 加入或退出遊戲，處理遊戲人員資訊(資料表.redis)
func (h *Handler) addGameStaff(activityID string, gameID string, userID string,
	game string, team string, round int64, action string) error {
	var key string // redis key
	if game == "tugofwar" {
		// 拔河遊戲
		if team == "left_team" {
			key = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲左隊人數資訊，SET
		} else if team == "right_team" {
			key = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
		}
	} else {
		// 其他遊戲
		key = config.GAME_ATTEND_REDIS + gameID // 遊戲人員資訊，SET
	}

	// 判斷玩家是否已加入遊戲(redis)
	if !h.redisConn.SetIsMember(key, userID) {
		// 未加入遊戲
		for l := 0; l < MaxRetries; l++ {
			// 上鎖，加入遊戲人員
			ok, _ := h.acquireLock(config.ADD_STAFF_LOCK_REDIS+gameID, LockExpiration)
			if ok == "OK" {
				// 釋放鎖
				// defer h.releaseLock(config.ADD_STAFF_LOCK_REDIS + gameID)

				// 判斷遊戲人數是否額滿並新增遊戲人員(redis處理，結算才將所有資料寫入資料表)
				if game == "tugofwar" {
					// 拔河遊戲
					// 判斷人數是否額滿
					gameModel := h.getGameInfo(gameID, game) // 遊戲資訊(redis處理)

					if h.redisConn.SetCard(key) >= gameModel.People {
						// 額滿
						// log.Println("拔河遊戲額滿")
						return errors.New("錯誤: 該隊伍遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝")
					} else {
						// 未額滿，將玩家寫入隊伍陣列中(redis)
						// log.Println("拔河遊戲未額滿，加入人員")
						if err := models.DefaultGameModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							AddPlayer(true, gameID, userID, team); err != nil {
							return err
						}
					}
				} else {
					// 其他遊戲
					// 判斷人數是否額滿
					gameModel := h.getGameInfo(gameID, game) // 遊戲資訊(redis處理)

					if h.redisConn.SetCard(key) >= gameModel.People {
						// 額滿
						// log.Println("其他遊戲額滿")
						return errors.New("錯誤: 遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝")
					} else {
						// log.Println("其他遊戲未額滿，加入人員")
						// 未額滿，將玩家寫入陣列中(redis)
						h.redisConn.SetAdd([]interface{}{key, userID})

						// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(game_attend人數)
						h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
					}
				}

				// 釋放鎖
				h.releaseLock(config.ADD_STAFF_LOCK_REDIS + gameID)
				break
			}

			// 鎖被佔用，稍微延遲後重試
			time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
		}
	}

	return nil
}

// closeGameStaffWebsocket 關閉玩家端遊戲人員ws，清除遊戲人員資料(資料表.redis)、更新遊戲人數(資料表.redis)、隊伍資料處理(redis)
func (h *Handler) closeGameStaffWebsocket(gameID string, userID string,
	game string, round int64, gameModel models.GameModel) {
	// 快問快答遊戲
	if game == "QA" {
		// 退出遊戲時，快問快答需要更新redis中的qa_people人數資訊(如果已存在玩家遊戲紀錄)
		// redis處理，判斷redis中是否有玩家加入遊戲的紀錄(game_attend_遊戲ID)
		if h.redisConn.SetIsMember(config.GAME_ATTEND_REDIS+gameID, userID) {
			// log.Println("存在玩家加入遊戲紀錄，遞減資料")

			// 存在玩家紀錄，遞減qa_people資料(redis處理)
			models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateQAPeople(true, gameID, -1)
		}

		// ###原本是用資料庫處理###
		// if models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 	IsUserExist(gameID, userID, round) {
		// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 		UpdateQAPeople(true, gameID, -1)
		// }
	}

	// 拔河遊戲
	if game == "tugofwar" {
		var (
			// teamModel = h.getTeamInfo(gameID) // 隊伍資訊
			team string
		)
		// fmt.Println("用戶ID: ", userID, ", 隊伍: ", team)

		// 判斷用戶的隊伍(redis處理，game_tugofwar_left_team_attend_遊戲ID、game_tugofwar_right_team_attend_遊戲ID)
		if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID, userID) {
			team = "left_team"
		}
		// 判斷右方隊伍
		if team == "" {
			if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID, userID) {
				team = "right_team"
			}
		}
		// log.Println("拔河遊戲退出遊戲，用戶隊伍: ", team)

		// ###原本是用資料庫處理###
		// 判斷左方隊伍
		// if utils.InArray(teamModel.LeftTeamPlayers, userID) {
		// 	team = "left_team"
		// }
		// // 判斷右方隊伍
		// if team == "" {
		// 	if utils.InArray(teamModel.RightTeamPlayers, userID) {
		// 		team = "right_team"
		// 	}
		// }
		// fmt.Println("用戶隊伍: ", team)

		// 判斷玩家是否為其中一隊人員(隊伍參數不為空)
		if team != "" {
			if gameModel.GameStatus == "open" || gameModel.GameStatus == "adjust" {
				// log.Println("拔河遊戲未開始，遞減隊伍人數、隊伍資料處理、判斷隊長")

				// 遊戲未開始前退出，遞減人數(redis)、隊伍資料處理(redis)
				// 判斷是否為隊長，清除redis中的隊長資訊

				// 遞減隊伍人數(redis)
				// models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	UpdateTeamPeople(true, gameID, team, -1, -1)

				// 隊伍資料處理(redis)，刪除玩家隊伍資料
				models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					DeletePlayer(true, gameID, team, userID)

				// 判斷是否為隊長(redis)
				var (
					teamModel = h.getTeamInfo(gameID) // 隊伍資訊
					leader    string
				)
				if team == "left_team" {
					leader = teamModel.LeftTeamLeader
				} else if team == "right_team" {
					leader = teamModel.RightTeamLeader
				}

				if userID == leader {
					// 清除隊長資訊(redis)
					models.DefaultGameModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						DeleteLeader(true, gameID, team)
				}

				// ###原本是用資料庫處理###
				// 遊戲未開始前退出,遞減人數(資料庫.redis)、刪除資料表資料、隊伍資料處理(redis)
				// 判斷是否為隊長，清除redis中的隊長資訊
				// 判斷是否有玩家加入遊戲的紀錄
				// if models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	IsUserExist(gameID, userID, round) {
				// 	// 將已報名遊戲的人員資料刪除
				// 	models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 		DeleteUser(gameID, userID, round)

				// 	// 遞減隊伍人數(資料庫.redis)
				// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 		UpdateTeamPeople(true, gameID, team, -1, -1)

				// 	// 隊伍資料處理(redis)
				// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 		DeletePlayer(true, gameID, team, userID)

				// 	// 判斷是否為隊長
				// 	var (
				// 		teamModel = h.getTeamInfo(gameID) // 隊伍資訊
				// 		leader    string
				// 	)
				// 	if team == "left_team" {
				// 		leader = teamModel.LeftTeamLeader
				// 	} else if team == "right_team" {
				// 		leader = teamModel.RightTeamLeader
				// 	}

				// 	if userID == leader {
				// 		// 清除隊長資訊(redis)
				// 		models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 			DeleteLeader(true, gameID, team)
				// 	}
				// }
			} else if gameModel.GameStatus == "start" || gameModel.GameStatus == "gaming" {
				// log.Println("拔河遊戲開始，不做任何處理")

				// 遊戲開始後退出，遞減人數(redis，玩家在回到遊戲時還可以加入)
				// 遞減隊伍人數(redis)
				// models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	UpdateTeamPeople(true, gameID, team, -1, -1)

				// 不進行隊伍資料處理(玩家在回到遊戲時還可以加入)

				// ###原本是用資料庫處理###
				// 遊戲開始後退出,遞減人數、不刪除資料表資料(玩家在回到遊戲時可以加入)
				// 判斷是否有玩家加入遊戲的紀錄
				// if models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	IsUserExist(gameID, userID, round) {
				// 	// 遞減隊伍人數(redis)
				// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 		UpdateTeamPeople(true, gameID, team, 0, -1)
				// }
			}
		}
	}
}

// #####拔河遊戲其他處理方式#####
// }
// } else if game != "tugofwar" {
// 其他遊戲
// 新增遊戲人員資料
// if !models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	IsUserExist(gameID, userID, round) {
// 	// log.Println("controller package GameStaffWebsocket function 新增遊戲人員資料")

// 	// 判斷遊戲人數是否額滿並新增遊戲人員(並更新redis中的遊戲參加人數)
// 	if err := models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		Add(models.NewGameStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			Round:      strconv.Itoa(int(round)),
// 			Status:     "success",
// 			Game:       game,
// 			Team:       team,
// 			// Black:      "no",
// 		}); err != nil {
// 		return err
// 	}
// }
// else {
// 		// log.Println("controller package GameStaffWebsocket function 已存在遊戲紀錄")
// 	}
// }

// else if action == "change" {
// }
// } else if action == "exit" {
// 	// 將已報名遊戲的人員資料刪除
// 	if err := models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		DeleteUser(gameID, userID, round); err != nil {
// 		return err
// 	}

// 	// 遞減隊伍人數(資料庫.redis)
// 	if err := models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateTeamPeople(true, gameID, team, -1, -1); err != nil {
// 		return err
// 	}

// 	// 隊伍資料處理(redis)
// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		DeletePlayer(true, gameID, team, userID)

// 	// 判斷是否為隊長
// 	var (
// 		teamModel = h.getTeamInfo(gameID) // 隊伍資訊
// 		leader    string
// 	)
// 	if team == "left_team" {
// 		leader = teamModel.LeftTeamLeader
// 	} else if team == "right_team" {
// 		leader = teamModel.RightTeamLeader
// 	}

// 	if userID == leader {
// 		// 清除隊長資訊(redis)
// 		models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 			DeleteLeader(true, gameID, team)
// 	}
// }

// if allow == "close" && len(gamePrizes) > 0 {
// 	err = errors.New("錯誤: 該活動不允許相同類型遊戲重複中獎，您已中獎，無法參加遊戲")
// 	b, _ := json.Marshal(PrizesParam{
// 		PrizeStaffs: winRecords,
// 		Error:       err.Error()})
// 	conn.WriteMessage(b)
// 	// return
// }
// 判斷遊戲是否允許重複中獎
// if err == nil {
// 	if gameModel.Allow == "close" && len(winRecords) > 0 {
// 		err = errors.New("錯誤: 該遊戲不允許重複中獎，您已中獎，無法參加下一輪遊戲")
// 		b, _ := json.Marshal(PrizesParam{
// 			PrizeStaffs: winRecords,
// 			Error:       err.Error()})
// 		conn.WriteMessage(b)
// 		// return
// 	}
// }
// #####拔河遊戲其他處理方式#####

// // 判斷是否為活動黑名單
// if isBlack := h.IsBlackStaff(activityID, "", "activity", result.User.UserID); isBlack {
// 	err = errors.New("錯誤: 用戶為活動黑名單人員，無法參加遊戲。如有疑問，請聯繫主辦單位")
// 	b, _ := json.Marshal(PrizesParam{
// 		PrizeStaffs: winRecords,
// 		Error:       err.Error()})
// 	conn.WriteMessage(b)
// 	// return
// }

// // 判斷是否為遊戲黑名單
// if err == nil {
// 	if isBlack := h.IsBlackStaff(activityID, gameID, game, result.User.UserID); isBlack {
// 		err = errors.New("錯誤: 用戶為遊戲黑名單人員，無法參加遊戲。如有疑問，請聯繫主辦單位")
// 		b, _ := json.Marshal(PrizesParam{
// 			PrizeStaffs: winRecords,
// 			Error:       err.Error()})
// 		conn.WriteMessage(b)
// 		// return
// 	}
// }

// 判斷該活動是否允許重複中獎
// if err == nil {
// 	if gameModel.AllGameAllow == "close" && len(activityPrizes) > 0 {
// 		err = errors.New("錯誤: 該活動不允許重複中獎，您已中獎，無法參加遊戲")
// 		b, _ := json.Marshal(PrizesParam{
// 			PrizeStaffs: winRecords,
// 			Error:       err.Error()})
// 		conn.WriteMessage(b)
// 		// return
// 	}
// }

// 判斷是否允許重複中獎
// if gameModel.Allow == "close" && len(winRecords) > 0 {
// 	err = errors.New("錯誤: 該場次遊戲不允許重複中獎，您已中獎，無法參加下一輪遊戲")
// 	log.Println("不允許重複中獎、已有中獎記錄: 無法參加遊戲")
// 	b, _ := json.Marshal(PrizesParam{
// 		PrizeStaffs: winRecords,
// 		Error:       err.Error()})
// 	conn.WriteMessage(b)
// 	// return // 關閉ws
// }

// // 判斷是否有該輪次遊戲紀錄
// for i := 0; i < len(gameRecords); i++ {
// 	if gameRecords[i].Round == round {
// 		record = gameRecords[i]
// 		log.Println("該輪次已玩過", record)
// 	}
// 	if (gameRecords[i].PrizeMethod != "" && gameRecords[i].PrizeMethod != "thanks") &&
// 		(gameRecords[i].PrizeType != "" && gameRecords[i].PrizeType != "thanks") &&
// 		gameRecords[i].PrizeID != "" {
// 		// 添加中獎記錄
// 		winRecords = append(winRecords, gameRecords[i])
// 	}
// }

// 新增玩家分數資料(redis)
// if game == "whack_mole" || game == "monopoly" {
// 	if err = models.DefaultScoreModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateScore(true, gameID, result.User.UserID, 0); err != nil {
// 		b, _ := json.Marshal(PrizesParam{
// 			Error: err.Error()})
// 		conn.WriteMessage(b)
// 		return
// 	}

// 	if err = models.DefaultScoreModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateScore2(true, gameID, result.User.UserID, 0); err != nil {
// 		b, _ := json.Marshal(PrizesParam{
// 			Error: err.Error()})
// 		conn.WriteMessage(b)
// 		return
// 	}
// }

// allow             string
// people            int64
// staffLock.Lock()
// defer staffLock.Unlock()

// models.NewScoreModel{
// ActivityID: activityID,
// GameID: gameID,
// UserID: result.User.UserID,
// Round:      strconv.Itoa(int(result.Game.GameRound)),
// }

// } else {
// 	log.Println("存在計分表: 更新分數為0")
// 	if err = scoreModel.Update(true, models.EditScoreModel{
// 		GameID: gameID,
// 		UserID: result.User.UserID,
// 		// Round:  strconv.Itoa(int(result.Game.GameRound)),
// 		Score: "0",
// 	}); err != nil {
// 		b, _ := json.Marshal(PrizesParam{
// 			Error: "錯誤: 更新人員計分表資料發生問題"})
// 		conn.WriteMessage(b)
// 		return
// 	}
// }

// scoreModel := models.DefaultScoreModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn)
// log.Println("判斷是否有計分資料表")
// if !scoreModel.IsUserExist(gameID, result.User.UserID, result.Game.GameRound) {
// log.Println("不存在計分表: 新建")

// Name:       result.User.Name,
// Avatar:     result.User.Avatar,

// ---測試沒問題的話就可以刪除了
// 新增遊戲人數並判斷是否額滿
// if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 	UpdateAttend(gameID, people); err != nil {
// 	log.Println("遊戲報名人數額滿")
// 	err = errors.New("錯誤: 遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝")
// 	b, _ := json.Marshal(WebsocketMessage{PrizeStaffs: prizes, Error: err.Error()})
// 	conn.WriteMessage(b)
// 	// return
