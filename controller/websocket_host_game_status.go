package controller

import (
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// jmeter -n -t 1207.jmx -l result.jtl

// jmeter -g result.jtl -o report

// hset game_Pb1eZzzoeWMTOM5GSRCC game_status open

// @Summary 主持人端遊戲狀態，更新資料表遊戲狀態
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, lottery, monopoly, QA, tugofwar, bingo, 3DGachaMachine, vote)
// @Param role query string true "role" Enums(host, guest)
// @param body body GameParam true "game param"
// @Success 200 {array} GameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/status/host [get]
func (h *Handler) HostGameStatusWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		role              = ctx.Query("role")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            GameParam
		oldStatus string
	)
	defer wsConn.Close()
	defer conn.Close()

	// log.Println("開啟主持端遊戲狀態ws", activityID, gameID, game)

	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if game == "3DGachaMachine" && (role != "host" && role != "guest") {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識扭蛋機遊戲資訊"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result                   GameParam
			newStatus, round, second string
			amount                   int64
		)

		data, err := conn.ReadMessage()

		// 遊戲資訊
		gameModel := h.getGameInfo(gameID, game)

		// 關閉網頁或連線，遊戲狀態關閉
		if err != nil {
			// 關閉主持端遊戲狀態ws，清除資料(資料表.redis資訊)、更新遊戲狀態、歸零遊戲人數
			h.closeHostGameStatusWebsocket(activityID, game, gameID, role, gameModel)
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)
		newStatus = result.Game.GameStatus

		// 更新遊戲狀態
		if newStatus != "" {
			// log.Println("有收到嗎?", result.Game.GameStatus)
			if newStatus == "open" {
				second = h.openStatus(gameID, game, gameModel)

				if game == "3DGachaMachine" {
					// 遊戲獎品數(redis處理)
					amount, _ = h.getPrizesAmount(gameID)
				}

				// #####加入測試資料#####start
				// if game == "tugofwar" {
				// 	leftplayers := []interface{}{config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID}
				// 	rightplayers := []interface{}{config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID}

				// 	// 兩隊人員資料
				// 	// 左隊人員
				// 	for i := 1; i <= 500; i++ {
				// 		leftplayers = append(leftplayers, strconv.Itoa(i))
				// 	}
				// 	// 右隊人員
				// 	for i := 501; i <= 999; i++ {
				// 		rightplayers = append(rightplayers, strconv.Itoa(i))
				// 	}

				// 	log.Println("人數多少人: ", len(leftplayers)-1, len(rightplayers)-1)

				// 	// 隊伍資料更新至redis
				// 	h.redisConn.SetAdd(leftplayers)
				// 	h.redisConn.SetAdd(rightplayers)
				// }
				// #####加入測試資料#####end
			} else if newStatus == "end" {
				h.endStatus(gameID, game)
			} else if newStatus == "start" {

				if oldStatus != newStatus {
					round = h.startStatus(gameID, game)
				}

			} else if newStatus == "gaming" {
				// 判斷是否為結算狀態
				if gameModel.GameStatus == "calculate" {
					// 當下狀態已經是結算，不更新為遊戲中
					newStatus = "calculate"
				} else if gameModel.GameStatus != "calculate" {
					// 當下不是結算狀態，更新遊戲秒數
					second = strconv.Itoa(int(result.Game.GameSecond))

					// 扭蛋機遊戲
					if game == "3DGachaMachine" {
						// 扭蛋機遊戲沒有秒數參數
						second = ""

						// 判斷是否為關閉狀態
						if gameModel.GameStatus == "close" {
							// 扭蛋機遊戲已關閉，不更新遊戲狀態
							newStatus = ""
						} else if gameModel.GameStatus != "close" {
							// 更新是否搖晃資料
							models.DefaultGameModel().
								SetConn(h.dbConn, h.redisConn, h.mongoConn).
								UpdateShake(true, gameID, result.Game.IsShake)
						}
					}
				}

				// #####加入測試資料#####start
				// if game == "whack_mole" || game == "monopoly" || game == "QA" {
				// 	rand.Seed(time.Now().UnixNano())
				// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
				// 	n := 999
				// 	wg.Add(n) //計數器
				// 	for i := 1; i <= n; i++ {
				// 		go func(i int) {
				// 			defer wg.Done()
				// 			if strconv.Itoa(i) == "1" {
				// 				h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, strconv.Itoa(i),
				// 					int64(rand.Intn(5)))
				// 			} else {
				// 				// redis
				// 				h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, strconv.Itoa(i),
				// 					int64(rand.Intn(5)))
				// 				h.redisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
				// 					float64(i)+0.1)
				// 			}
				// 		}(i)
				// 	}
				// 	wg.Wait() //等待計數器歸0
				// } else if game == "tugofwar" {
				// 	leftplayers := make([]string, 0)
				// 	rightplayers := make([]string, 0)
				// 	// 左隊人員
				// 	for i := 1; i <= 500; i++ {
				// 		leftplayers = append(leftplayers, strconv.Itoa(i))
				// 	}
				// 	// 右隊人員
				// 	for i := 501; i <= 999; i++ {
				// 		rightplayers = append(rightplayers, strconv.Itoa(i))
				// 	}

				// 	rand.Seed(time.Now().UnixNano())
				// 	for _, userID := range leftplayers {
				// 		// idint, _ := strconv.Atoi(userID)
				// 		// 加入分數資料
				// 		h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, userID, int64(rand.Intn(5)))
				// 	}
				// 	for _, userID := range rightplayers {
				// 		// idint, _ := strconv.Atoi(userID)
				// 		// 加入分數資料
				// 		h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, userID, int64(rand.Intn(5)))
				// 	}
				// }
				// #####加入測試資料#####end

			} else if newStatus == "question" {
				// 清空redis題目相關資訊(所有資訊歸零)
				if game == "QA" {
					h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
						"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

					// 設置過期時間
					// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

					// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
					h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")
				}

				// #####加入測試資料#####start
				// if game == "QA" {
				// 	rand.Seed(time.Now().UnixNano())
				// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
				// 	n := 200
				// 	wg.Add(n) //計數器
				// 	for i := 1; i <= n; i++ {
				// 		go func(i int) {
				// 			defer wg.Done()

				// 			option := ""
				// 			x := i % 4
				// 			if x == 0 {
				// 				option = "A"
				// 			} else if x == 1 {
				// 				option = "B"
				// 			} else if x == 2 {
				// 				option = "C"
				// 			} else if x == 3 {
				// 				option = "D"
				// 			}
				// 			// 遞增redis裡總答題人數、各選項答題人數資訊
				// 			h.redisConn.HashIncrCache(config.QA_REDIS+gameID, "Total")
				// 			h.redisConn.HashIncrCache(config.QA_REDIS+gameID, option)
				// 		}(i)
				// 	}
				// 	wg.Wait() //等待計數器歸0
				// }
				// #####加入測試資料#####end

			} else if newStatus == "result" {
				if game == "QA" && oldStatus != newStatus {
					if gameModel.QARound < gameModel.TotalQA {
						// 快問快答公布答案後，遞增qa_round(目前進行的題數)資訊
						if err = models.DefaultGameModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							UpdateQARound(true, gameID, "+1"); err != nil {
							b, _ := json.Marshal(GameParam{Error: "錯誤: 無法更新目前進行的題數資料"})
							conn.WriteMessage(b)

							models.DefaultErrorLogModel().
								SetConn(h.dbConn, h.redisConn, h.mongoConn).
								Add(models.EditErrorLogModel{
									UserID:    utils.ClientIP(ctx.Request),
									Code:      http.StatusBadRequest,
									Method:    ctx.Request.Method,
									Path:      ctx.Request.URL.Path,
									Message:   "錯誤: 無法更新目前進行的題數資料",
									PathQuery: ctx.Request.URL.RawQuery,
								})
							return
						}
					}
				}
			} else if newStatus == "calculate" {
				// 投票遊戲
				if game == "vote" {

					// 收到結算狀態時，將投票結束時間改為目前時間
					if err = models.DefaultGameModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						UpdateVoteEndTime(true, gameID); err != nil {
						b, _ := json.Marshal(GameParam{Error: err.Error()})
						conn.WriteMessage(b)

						// 寫入error log表
						models.DefaultErrorLogModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							Add(models.EditErrorLogModel{
								UserID:    utils.ClientIP(ctx.Request),
								Code:      http.StatusBadRequest,
								Method:    ctx.Request.Method,
								Path:      ctx.Request.URL.Path,
								Message:   err.Error(),
								PathQuery: ctx.Request.URL.RawQuery,
							})
						return
					}
				}
			} else if newStatus == "close" {
			} else if newStatus == "order" {
			} else if newStatus == "adjust" {
				// 調整雙方隊伍資訊
				if err = h.adjustTeam(gameID, game); err != nil {
					b, _ := json.Marshal(GameParam{Error: err.Error()})
					conn.WriteMessage(b)

					models.DefaultErrorLogModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						Add(models.EditErrorLogModel{
							UserID:    utils.ClientIP(ctx.Request),
							Code:      http.StatusBadRequest,
							Method:    ctx.Request.Method,
							Path:      ctx.Request.URL.Path,
							Message:   err.Error(),
							PathQuery: ctx.Request.URL.RawQuery,
						})
					return
				}
			} else if newStatus == "openball" {
				// 扭蛋機遊戲
			} else if newStatus == "showprize" {
				// 扭蛋機遊戲
				if game == "3DGachaMachine" {
					// 遊戲獎品數(redis處理)
					amount, _ = h.getPrizesAmount(gameID)
				}
			}

			if err = models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateGameStatus(true, gameID, round, second, newStatus); err != nil {
				b, _ := json.Marshal(GameParam{Error: "錯誤: 無法更新遊戲狀態"})
				conn.WriteMessage(b)

				models.DefaultErrorLogModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Add(models.EditErrorLogModel{
						UserID:    utils.ClientIP(ctx.Request),
						Code:      http.StatusBadRequest,
						Method:    ctx.Request.Method,
						Path:      ctx.Request.URL.Path,
						Message:   "錯誤: 無法更新遊戲狀態",
						PathQuery: ctx.Request.URL.RawQuery,
					})
				return
			}

			// 用來避免重複收到一樣的遊戲狀態判斷
			oldStatus = newStatus

			// 回傳遊戲狀態訊息給前端
			b, _ := json.Marshal(GameParam{
				Game: GameModel{
					GameStatus:  newStatus,
					PrizeAmount: amount,
				},
			})
			conn.WriteMessage(b)
		}
	}
}

// startStatus 主持端接收start狀態
func (h *Handler) startStatus(gameID string, game string) string {
	var (
		round string
	)
	// 開始時更新輪次與清空遊戲人數資訊
	if game == "redpack" || game == "ropepack" ||
		game == "whack_mole" || game == "monopoly" ||
		game == "tugofwar" || game == "bingo" {
		round = "+1"
		// attend = "0"
	}

	// 快問快答
	if game == "QA" {
		// 清空redis題目相關資訊(所有資訊歸零)
		h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
			"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

		// 設置過期時間
		// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")
	}

	// 拔河遊戲
	if game == "tugofwar" {
		var (
			teamModel        = h.getTeamInfo(gameID)      // 隊伍資訊
			leftTeamPlayers  = teamModel.LeftTeamPlayers  // 左方隊伍玩家
			rightTeamPlayers = teamModel.RightTeamPlayers // 右方隊伍玩家
			leftTeamLeader   = teamModel.LeftTeamLeader   // 左方隊伍隊長
			rightTeamLeader  = teamModel.RightTeamLeader  // 右方隊伍隊長
		)

		// 判斷雙方隊伍是否有隊長(沒有則隨機挑選隊長)
		// 判斷左方隊伍是否有隊長
		if leftTeamLeader != "" {
			// 左方隊伍有隊長，不用隨機挑選隊長
		} else if leftTeamLeader == "" {
			// 左方隊伍無隊長，隨機挑選隊長

			if len(leftTeamPlayers) > 0 {
				// 隨機挑選
				rand.Seed(time.Now().UnixNano())
				random := rand.Intn(len(leftTeamPlayers))
				leftTeamLeader = leftTeamPlayers[random]
			}
		}
		// 判斷右方隊伍是否有隊長
		if rightTeamLeader != "" {
			// 右方隊伍有隊長，不用隨機挑選隊長
		} else if rightTeamLeader == "" {
			// 右方隊伍無隊長，隨機挑選隊長

			if len(rightTeamPlayers) > 0 {
				// 隨機挑選
				rand.Seed(time.Now().UnixNano())
				random := rand.Intn(len(rightTeamPlayers))
				rightTeamLeader = rightTeamPlayers[random]
			}
		}

		// 更新雙方隊伍隊長資料(redis)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateTeamLeader(true, gameID, leftTeamLeader, rightTeamLeader)
	}

	// if game == "bingo" {
	// 	// fmt.Println("清除資訊")
	// 	// 賓果遊戲
	// 	h.redisConn.DelCache(config.SCORES_REDIS + gameID)            // 分數
	// 	h.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID) // 紀錄抽過的號碼，LIST
	// 	h.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)   // 賓果中獎人員，ZSET
	// }

	// 扭蛋機遊戲
	// if game == "3DGachaMachine" {
	// }

	return round
}

// endStatus 主持端接收end狀態
func (h *Handler) endStatus(gameID string, game string) {
	// 快問快答
	if game == "QA" {
		// 清空qa_people參數(redis)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateQAPeople(true, gameID, 0)
	}

	// 拔河遊戲
	// if game == "tugofwar" {
	// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		UpdateBothTeamPeople(true, gameID, 0, 0)

	// 清空遊戲隊伍人員資料
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

	// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

	// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	// h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	// }

	// 賓果遊戲
	if game == "bingo" {
		// 歸零賓果回合數(redis)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateBingoRound(true, gameID, "0")
	}

	// 清除分數資訊
	h.redisConn.DelCache(config.SCORES_REDIS + gameID)            // 分數
	h.redisConn.DelCache(config.SCORES_2_REDIS + gameID)          // 第二分數
	h.redisConn.DelCache(config.CORRECT_REDIS + gameID)           // 清除答對題數資訊
	h.redisConn.DelCache(config.QA_REDIS + gameID)                // 快問快答題目資訊
	h.redisConn.DelCache(config.QA_RECORD_REDIS + gameID)         // 清除答對題數資訊
	h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)    // 中獎人員
	h.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + gameID) // 未中獎人員
	// 賓果遊戲
	h.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID)     // 紀錄抽過的號碼，LIST
	h.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)       // 賓果中獎人員，ZSET
	h.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + gameID)      // 紀錄玩家的號碼排序，HASH
	h.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + gameID) // 紀錄玩家是否即將中獎，HASH

	// 清空遊戲隊伍人員資料
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

	h.redisConn.DelCache(config.GAME_ATTEND_REDIS + gameID) // 遊戲人員

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
	h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
	h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

}

// openStatus 主持端接收open狀態
func (h *Handler) openStatus(gameID string, game string, gameModel models.GameModel) string {

	var (
		second = strconv.Itoa(int(gameModel.Second)) // 重置遊戲秒數
		// attend string
	)
	// 搖號抽獎遊戲判斷是否被開啟
	// if game == "draw_numbers" && gameModel.GameStatus != "close" {
	// 	b, _ := json.Marshal(GameParam{Error: "錯誤: 該場次遊戲已被開啟，無法重複開啟遊戲"})
	// 	conn.WriteMessage(b)
	// 	return
	// }

	if game == "redpack" || game == "ropepack" ||
		game == "whack_mole" || game == "monopoly" ||
		game == "tugofwar" || game == "bingo" {
		// 打開遊戲頁面，清空遊戲參加人數
		// attend = "0"
		h.redisConn.DelCache(config.GAME_ATTEND_REDIS + gameID) // 遊戲人員

		// 將已報名遊戲的人員資料刪除
		// models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 	DeleteGame(gameID, gameModel.GameRound)

		// 清除中獎人員資料
		h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
		h.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + gameID) // 未中獎人員
		// 清除分數資料
		h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
		h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
		// 清除遊戲隊伍資訊
		h.redisConn.DelCache(config.GAME_TEAM_REDIS + gameID)
		h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
		h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

		// 賓果遊戲
		h.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID)     // 紀錄抽過的號碼，LIST
		h.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)       // 賓果中獎人員，ZSET
		h.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + gameID)      // 紀錄玩家的號碼排序，HASH
		h.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + gameID) // 紀錄玩家是否即將中獎，HASH

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資料)
		h.redisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		// #####加入測試資料(加入遊戲人員)#####start
		// params := []interface{}{config.GAME_ATTEND_REDIS + gameID}

		// // 加入遊戲人員
		// for i := 1; i <= 999; i++ {
		// 	params = append(params, strconv.Itoa(i))
		// }

		// h.redisConn.SetAdd(params)

		// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		// h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		// #####加入測試資料(加入遊戲人員)#####end
	}

	// 快問快答
	if game == "QA" {
		// 清空redis題目相關資訊(所有資訊歸零)
		h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
			"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

		// 設置過期時間
		// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")

		// 清空qa_people參數(redis)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateQAPeople(true, gameID, 0)

		// #####加入測試資料(加入遊戲人員)#####start

		// params := []interface{}{config.GAME_ATTEND_REDIS + gameID}

		// // 加入遊戲人員
		// for i := 1; i <= 999; i++ {
		// 	params = append(params, strconv.Itoa(i))
		// }

		// h.redisConn.SetAdd(params)

		// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		// h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		// #####加入測試資料(加入遊戲人員)#####end
	}

	// 拔河遊戲
	// if game == "tugofwar" {
	// 	// 歸零左右方隊伍人數(資料表、redis)
	// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		UpdateBothTeamPeople(true, gameID, 0, 0)

	// 清空遊戲隊伍人員資料
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
	// h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	// h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	// }

	// 賓果遊戲
	if game == "bingo" {
		// 歸零賓果回合數(redis)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateBingoRound(true, gameID, "0")

	}

	// 扭蛋機遊戲
	if game == "3DGachaMachine" {
		second = ""
		// 搖晃資料改回預設值false
		h.redisConn.HashSetCache(config.GAME_REDIS+gameID, "is_shake", "false")

		// 設置過期時間
		// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		// 清空redis中中獎紀錄資料
		h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
	}

	return second
}

// adjustTeam 調整雙方隊伍資訊
func (h *Handler) adjustTeam(gameID string, game string) error {
	// 拔河遊戲
	if game == "tugofwar" {
		var (
			teamModel                             = h.getTeamInfo(gameID)            // 隊伍資訊
			leftTeamPeople                        = len(teamModel.LeftTeamPlayers)   // 左方隊伍人數
			rightTeamPeople                       = len(teamModel.RightTeamPlayers)  // 右方隊伍人數
			totalPeople                           = leftTeamPeople + rightTeamPeople // 雙方隊伍人數
			leftTeamPlayers                       = teamModel.LeftTeamPlayers        // 左方隊伍玩家
			rightTeamPlayers                      = teamModel.RightTeamPlayers       // 右方隊伍玩家
			leftTeamLeader                        = teamModel.LeftTeamLeader         // 左方隊伍隊長
			rightTeamLeader                       = teamModel.RightTeamLeader        // 右方隊伍隊長
			newLeftTeamPeople, newRightTeamPeople int                                // 調整後的雙方隊伍人數
			isEven                                = totalPeople%2 != 1               // 是否為偶數
		)

		// 判斷雙方隊伍人數
		if isEven && leftTeamPeople == rightTeamPeople {
			// 雙方隊伍人數相同(總人數偶數)，不需調整人員
			newLeftTeamPeople = leftTeamPeople
			newRightTeamPeople = rightTeamPeople
		} else if !isEven &&
			(leftTeamPeople-rightTeamPeople == 1 || rightTeamPeople-leftTeamPeople == 1) {
			// 雙方隊伍人數相差1(總人數奇數)，不需調整人員
			newLeftTeamPeople = leftTeamPeople
			newRightTeamPeople = rightTeamPeople
		} else if leftTeamPeople != rightTeamPeople {
			// 雙方隊伍人數不同，判斷是否為偶數，設置左右方調整後的人數
			if isEven {
				// 偶數，左右方調整後的人數一樣
				newLeftTeamPeople = totalPeople / 2
				newRightTeamPeople = totalPeople / 2
			} else if !isEven {
				// 奇數，原本人數多的隊伍方會多一人
				if leftTeamPeople > rightTeamPeople {
					// 原本左方隊伍人數較多，調整後的左方隊伍多一人
					newLeftTeamPeople = totalPeople/2 + 1
					newRightTeamPeople = totalPeople / 2
				} else if leftTeamPeople < rightTeamPeople {
					// 原本右方隊伍人數較多，調整後的右方隊伍多一人
					newLeftTeamPeople = totalPeople / 2
					newRightTeamPeople = totalPeople/2 + 1
				}
			}

			// 調整雙方隊伍人員(先判斷哪一方人數較多)
			if leftTeamPeople > rightTeamPeople {
				// 左方隊伍人數較多，調整左方人數
				for i := leftTeamPeople - 1; i >= 0; i-- {
					player := leftTeamPlayers[i]
					if len(leftTeamPlayers) == newLeftTeamPeople {
						// 調整人數完成
						break
					}

					// 判斷是否為隊長
					if player == leftTeamLeader {
						// 隊長，不能移至另一隊
						continue
					} else if player != leftTeamLeader {
						// 不為隊長，移至另一隊
						leftTeamPlayers = append((leftTeamPlayers)[:i], leftTeamPlayers[i+1:]...)
						rightTeamPlayers = append(rightTeamPlayers, player)
					}
				}
			} else if leftTeamPeople < rightTeamPeople {
				// 右方隊伍人數較多，調整右方人數
				for i := rightTeamPeople - 1; i >= 0; i-- {
					player := rightTeamPlayers[i]
					if len(rightTeamPlayers) == newRightTeamPeople {
						// 調整人數完成
						break
					}

					// 判斷是否為隊長
					if player == rightTeamLeader {
						// 隊長，不能移至另一隊
						continue
					} else if player != rightTeamLeader {
						// 不為隊長，移至另一隊
						leftTeamPlayers = append(leftTeamPlayers, player)
						rightTeamPlayers = append((rightTeamPlayers)[:i], rightTeamPlayers[i+1:]...)
					}
				}
			}
		}

		// 更新雙方隊伍人數資料(資料庫.redis)
		// if err := models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 	UpdateBothTeamPeople(true, gameID, int64(newLeftTeamPeople), int64(newRightTeamPeople)); err != nil {
		// 	return err
		// }

		// 更新雙方隊伍玩家資料(redis)
		if err := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateTeamPlayers(true, gameID, leftTeamPlayers, rightTeamPlayers); err != nil {
			return err
		}
	}
	return nil
}

// closeHostGameStatusWebsocket 關閉主持端遊戲狀態ws，清除資料(資料表.redis資訊)、更新遊戲狀態、歸零遊戲人數
func (h *Handler) closeHostGameStatusWebsocket(activityID, game string, gameID string, role string,
	gameModel models.GameModel) {
	// var (
	// attend string
	// )

	// 紅包遊戲
	if game == "redpack" || game == "ropepack" {
		// 判斷遊戲狀態
		if gameModel.GameStatus == "gaming" {
			// 紅包遊戲正在遊戲中
			// log.Println("目前正在遊戲中，將遊戲人員資料寫入資料表")

			// 取得該場次所有遊戲人員資訊
			users, _ := h.redisConn.SetGetMembers(config.GAME_ATTEND_REDIS + gameID)
			// log.Println("將遊戲人員紀錄寫入資料表: ", len(users))

			// 批量將所有遊戲人員資料寫入資料表(紅包遊戲的遊戲中輪次已遞增)
			models.DefaultGameStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Adds(len(users), activityID, gameID, game, users, int(gameModel.GameRound-1))

			// log.Println("將遊戲人員紀錄寫入資料表完成")

			var (
				staffs   = make([]models.PrizeStaffModel, 0) // 中獎人資料
				nostaffs = make([]models.PrizeStaffModel, 0) // 未中獎人員資料
				wg       sync.WaitGroup                      // 宣告WaitGroup 用以等待執行序
			)

			// 取得紅包遊戲中獎人員資料
			staffs, _, _, _ = h.gaming(gameModel, gameID, game, gameModel.GameRound-1, 1)

			// 取得紅包遊戲未中獎人員資料
			var recordsJson []string // 未中獎人員json資料
			recordsJson, _ = h.redisConn.ListRange(config.NO_WINNING_STAFFS_REDIS+gameID, 0, 0)
			if len(recordsJson) > 0 {
				for _, record := range recordsJson {
					var staff models.PrizeStaffModel
					// 解碼
					json.Unmarshal([]byte(record), &staff)

					nostaffs = append(nostaffs, staff)
				}
			}

			// 批量寫入所有中獎人員資料
			// log.Println("中斷遊戲，批量寫入所有中獎人員資料: ", len(staffs))

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
			models.DefaultPrizeStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Adds(len(staffs), activityID, gameID, game, userIDs, peizeIDs,
					int(gameModel.GameRound-1), scores1, scores2, ranks, teams, leaders, mvps)

			// log.Println("中斷遊戲，批量寫入所有中獎人員完成")

			// 批量寫入所有未中獎人員資料
			// log.Println("中斷遊戲，批量寫入所有未中獎人員資料: ", len(nostaffs))

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
			models.DefaultPrizeStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Adds(len(nostaffs), activityID, gameID, game, userIDs, peizeIDs,
					int(gameModel.GameRound-1), scores1, scores2, ranks, teams, leaders, mvps)

			// log.Println("中斷遊戲，批量寫入所有未中獎人員完成")

			// 處理剩餘獎品
			// 從獎品redis(hash格式)中取得所有獎品剩餘數量，並更新至資料表中(多線呈處理)
			prizes, _ := h.getPrizes(gameID)

			// log.Println("中斷遊戲，更新紅包遊戲剩餘獎品: ", len(prizes))

			// 多線呈更新資料表獎品資訊
			wg.Add(len(prizes)) //計數器
			for i := 0; i < len(prizes); i++ {
				go func(i int) {
					defer wg.Done()

					// 更新獎品剩餘數量
					models.DefaultPrizeModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						UpdateRemain(
							gameID, prizes[i].PrizeID, prizes[i].PrizeRemain)
				}(i)
			}
			wg.Wait() //等待計數器歸0

			// log.Println("中斷遊戲，更新紅包遊戲剩餘完成")

		} else {
			// log.Println("目前不在遊戲中，不用將遊戲人員資料寫入資料表: ", gameModel.GameStatus)
		}
	}

	if game == "redpack" || game == "ropepack" ||
		game == "whack_mole" || game == "monopoly" ||
		game == "tugofwar" || game == "bingo" {
		// 清空遊戲人數、隊伍人數、關閉遊戲
		// attend = "0"

		// 清空遊戲人員
		h.redisConn.DelCache(config.GAME_ATTEND_REDIS + gameID)                     // 遊戲人員
		h.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左方隊伍
		h.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(遊戲人員)
		h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		// 將已報名遊戲的人員資料刪除
		// models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 	DeleteGame(gameID, gameModel.GameRound)
	}

	// 快問快答(不用將加入遊戲人員資料清除，只需要把qa_people歸零)
	if game == "QA" {
		// 清空redis題目相關資訊(所有資訊歸零)
		h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
			"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

		// 設置過期時間
		// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")

		// 關閉頁面時快問快答遊戲關閉遊戲狀態還有qa_people參數
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateQAPeople(true, gameID, 0)
	}

	// 拔河遊戲
	// if game == "tugofwar" {
	// 	// 關閉頁面時歸零左右方隊伍人數
	// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		UpdateBothTeamPeople(true, gameID, 0, 0)
	// }

	// 賓果遊戲
	if game == "bingo" {
		// 關閉頁面時歸零賓果回合數
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateBingoRound(true, gameID, "0")
	}

	// 扭蛋機遊戲
	if game == "3DGachaMachine" {
		// 搖晃資料改回預設值false
		h.redisConn.HashSetCache(config.GAME_REDIS+gameID, "is_shake", "false")

		// 設置過期時間
		// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		if role == "host" {
			// 主持端關閉遊戲頁面時，將遊戲狀態關閉(close)
			models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateGameStatus(true, gameID, "", "", "close")

			// 清空redis中中獎紀錄資料
			h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
		} else if role == "guest" {
			// 玩家端關閉遊戲頁面時，先判斷目前該場次的遊戲狀態
			if gameModel.GameStatus == "close" {
				// 主持端已關閉遊戲頁面，不更新遊戲狀態
			} else if gameModel.GameStatus != "close" {
				// 主持端未關閉遊戲頁面，將遊戲狀態更新為open，讓其他玩家可以進入遊戲頁面
				models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateGameStatus(true, gameID, "", "", "open")

				// 清空redis中中獎紀錄資料
				h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
			}
		}
	}

	// 扭蛋機遊戲需要利用role參數判斷遊戲狀態的處理(上面)
	// 投票遊戲不需要更新遊戲狀態
	if game != "3DGachaMachine" && game != "vote" {
		// 更新遊戲狀態(關閉)
		models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateGameStatus(true, gameID, "", strconv.Itoa(int(gameModel.Second)), "close")
	}

	// 刪除redis所有遊戲相關資訊
	// h.redisConn.DelCache(config.GAME_REDIS + gameID)              // 遊戲資訊
	// h.redisConn.DelCache(config.PRIZES_REDIS + gameID)            // 獎品數量
	// h.redisConn.DelCache(config.BLACK_STAFFS_GAME_REDIS + gameID) // 黑名單人員
	// h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)    // 中獎人員

	// 刪除分數redis資訊(快問快答不清除分數資訊)
	if game == "redpack" || game == "ropepack" ||
		game == "whack_mole" || game == "monopoly" ||
		game == "tugofwar" || game == "bingo" {
		h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)    // 中獎人員
		h.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + gameID) // 未中獎人員
		h.redisConn.DelCache(config.SCORES_REDIS + gameID)            // 分數
		h.redisConn.DelCache(config.SCORES_2_REDIS + gameID)          // 第二分數

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

		if game == "tugofwar" {
			// 清除遊戲隊伍資訊
			h.redisConn.DelCache(config.GAME_TEAM_REDIS + gameID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資料)
			h.redisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")
		}

		if game == "bingo" {
			// 賓果遊戲
			h.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID)     // 紀錄抽過的號碼，LIST
			h.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)       // 賓果中獎人員，ZSET
			h.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + gameID)      // 紀錄玩家的號碼排序，HASH
			h.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + gameID) // 紀錄玩家是否即將中獎，HASH

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, "修改資料")
		}
	} else if game == "QA" {
		h.redisConn.DelCache(config.QA_REDIS + gameID) // 清除快問快答題目資訊

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")
	} else if game == "vote" {
		// 投票遊戲
		// h.redisConn.DelCache(config.SCORES_REDIS + gameID) // 分數

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		// h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
	}
}

// 測試增加1名遊戲人員---start
// if newStatus == "open" {
// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateGameStatus(true, gameID, "", "", "1", newStatus)
// }
// 測試增加1名遊戲人員---end

// 賓果遊戲測試資料-----start
// for i := 11; i <= 14; i++ {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 5)          // 賓果連線數資料
// 	h.redisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID, strconv.Itoa(i), 2) // 完成賓果的回合數
// }
// 賓果遊戲測試資料-----end

// fmt.Println("主持人端收到calculate狀態")

// if oldStatus != newStatus {
// 	// 快問快答結算時更新輪次以及題目進行題數資訊(清空qa_people人數)
// 	if game == "QA" {
// 		round = "+1" // 輪次加1
// 		attend = "0" // 清空遊戲人數
// 		// 清空redis題目相關資訊(所有資訊歸零)
// 		h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
// 			"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})
// 		// 更新題目進行題數資料(qa_round=1)
// 		if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).UpdateQARound(true, gameID, "1"); err != nil {
// 			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法更新題目進行題數資料(qa_round=1)"})
// 		conn.WriteMessage(b)
// 			return
// 		}
// 	}
// }

// 加入測試人員(競技遊戲)-----start
// rand.Seed(time.Now().UnixNano())
// var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// n := 10
// wg.Add(n) //計數器
// for i := 1; i <= n; i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		if strconv.Itoa(i) == "1" {
// 			h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 				rand.Int63n(100))
// 		} else {
// 			// redis
// 			h.redisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 				rand.Int63n(100))
// 			h.redisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
// 				float64(i)+0.1)
// 		}

// 		// option := ""
// 		// x := i % 4
// 		// if x == 0 {
// 		// 	option = "A"
// 		// } else if x == 1 {
// 		// 	option = "B"
// 		// } else if x == 2 {
// 		// 	option = "C"
// 		// } else if x == 3 {
// 		// 	option = "D"
// 		// }
// 		// 遞增redis裡總答題人數、各選項答題人數資訊
// 		// fmt.Println("有執行嗎")
// 		// h.redisConn.HashIncrCache(config.QA_REDIS+gameID, "Total")
// 		// h.redisConn.HashIncrCache(config.QA_REDIS+gameID, option)
// 	}(i)
// }
// wg.Wait() //等待計數器歸0
// 加入測試人員(競技遊戲)-----end

// 賓果遊戲測試人員，加入賓果人員、即將賓果人員(設置五條線賓果)-----start
// 賓果人員(第一輪中獎人員，前50名)
// for i := 1; i <= 4; i++ {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 5)          // 賓果連線數資料
// 	h.redisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID, strconv.Itoa(i), 1) // 完成賓果的回合數
// }

// 賓果人員(第二輪中獎人員，後50名)
// for i := 11; i <= 14; i++ {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 4)          // 賓果連線數資料
// 	h.redisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID, strconv.Itoa(i), 2) // 完成賓果的回合數
// }

// 即將賓果人員(四條線，10名)
// for i := 3; i <= 10; i++ {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 1) // 賓果連線數資料
// }
// 賓果遊戲測試人員，加入賓果人員、即將賓果人員-----end

// 套紅包測試人員，加入中獎人員-----start
// for i := 1; i <= 10; i++ {
// 	now, _ := time.ParseInLocation("2006-01-02 15:04:05",
// 		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)

// 	record := models.PrizeStaffModel{
// 		ID:            int64(i),
// 		ActivityID:    activityID,
// 		GameID:        gameID,
// 		PrizeID:       strconv.Itoa(i),
// 		UserID:        strconv.Itoa(i),
// 		Name:          "我的名字需要十五個字還缺五個字",
// 		Avatar:        "https://dev.hilives.net/admin/uploads/0aSRpSne.png",
// 		Game:          game,
// 		PrizeName:     "獎品",
// 		PrizeType:     "first",
// 		PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 		PrizePrice:    1,
// 		PrizeMethod:   "site",
// 		PrizePassword: "12345",
// 		Round:         1,
// 		WinTime:       now.Format("2006-01-02 15:04:05"),
// 		Status:        "no",
// 	}

// 	// 依照中獎的順序將資料推送至list格式的redis中(沒有按照獎品類型排序)
// 	h.redisConn.ListRPush(config.WINNING_STAFFS_REDIS+gameID, utils.JSON(record))
// }
// 套紅包測試人員-----end

// 加入測試人員(拔河遊戲)-----start
// leftplayers := []string{"1", "2", "3",
// 	"4", "Ua4fc0eda353d11d6a1af106e4fff53e6"}
// rightplayers := []string{"Ua8782dfb9073e4dac11c9d5c2744e9d8", "12", "13",
// 	"14", "15"}

// // 更新至redis
// h.redisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID,
// 	"left_team_players", utils.JSON(leftplayers))
// h.redisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID,
// 	"right_team_players", utils.JSON(rightplayers))

// rand.Seed(time.Now().UnixNano())
// for _, userID := range leftplayers {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, userID, rand.Int63n(20))
// }
// for _, userID := range rightplayers {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, userID, rand.Int63n(20))
// }
// 加入測試人員(拔河遊戲)-----end

// 加入測試人員(拔河遊戲)---start
// leftplayers := []string{"1", "2", "3",
// 	"4", "5", "6", "7", "8", "9"}
// rightplayers := []string{"11", "12", "13",
// 	"14", "15", "16", "17", "18", "19", "20"}

// // 更新至redis
// h.redisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID,
// 	"left_team_players", utils.JSON(leftplayers))
// h.redisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID,
// 	"right_team_players", utils.JSON(rightplayers))

// rand.Seed(time.Now().UnixNano())
// for _, userID := range leftplayers {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, userID, 0)
// }
// for _, userID := range rightplayers {
// 	// 加入分數資料
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, userID, 0)
// }

// // 關閉頁面時歸零左右方隊伍人數
// models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateBothTeamPeople(true, gameID, int64(len(leftplayers)), int64(len(rightplayers)))
// 加入測試人員(拔河遊戲)---end

// 原本result時處理
// 清空redis題目相關資訊(所有資訊歸零)
// h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
// 	"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})

// 原本start時處理
// 清除redis資料
// if game == "redpack" || game == "ropepack" {
// 	// 清除中獎人員資料
// 	h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
// } else if game == "whack_mole" || game == "monopoly" {
// 	// 清除中獎人員資料
// 	h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
// 	// 清除分數資料
// 	h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
// 	h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
// } else if game == "QA" {
// 	// 清空redis題目相關資訊(所有資訊歸零)
// 	h.redisConn.HashMultiSetCache([]interface{}{config.QA_REDIS + gameID,
// 		"A", 0, "B", 0, "C", 0, "D", 0, "Total", 0})
// }

// 原本end時處理
// 快問快答結算發獎後更新輪次以及題目進行題數資訊(清空qa_people人數)
// if game == "QA" {
// round = "+1" // 輪次加1
// attend = "0" // 清空遊戲人數
// // 更新題目進行題數資料(qa_round=1)
// if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).UpdateQARound(true, gameID, "1"); err != nil {
// 	b, _ := json.Marshal(GameParam{Error: "錯誤: 無法更新題目進行題數資料(qa_round=1)"})
// 	conn.WriteMessage(b)
// 	return
// }

// 原本的遊戲狀態判斷
// if newStatus != "open" && newStatus != "start" && newStatus != "question" &&
// 	newStatus != "gaming" && newStatus != "result" && newStatus != "calculate" &&
// 	newStatus != "end" && newStatus != "close" && newStatus != "order" {
// 	b, _ := json.Marshal(GameParam{Error: "錯誤: 無法辨識遊戲狀態、輪次資訊"})
// 	conn.WriteMessage(b)
// 	return
// }

// 原本start時處理
// 新增測試用戶分數紀錄(測試用，999筆分數緩存資料)
// for i := 1; i < 800; i++ {
// 	// redis
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 		int64(i))
// 	h.redisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
// 		float64(i)+0.1)
// }

// 快問快答等待開始時更新題目進行題數資料(更新qa_round=1)
// if game == "QA" {
// 	if err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).UpdateQARound(true, gameID, "1"); err != nil {
// 		b, _ := json.Marshal(GameParam{Error: "錯誤: 無法更新目前進行的題數資料(qa_round=1)"})
// 		conn.WriteMessage(b)
// 		return
// 	}
// }

// 設置過期時間(儲存玩家分數的redis)
// if game == "whack_mole" || game == "monopoly" || game == "QA" {
// 	h.redisConn.SetExpire(config.SCORES_REDIS+gameID, config.REDIS_EXPIRE)
// 	h.redisConn.SetExpire(config.SCORES_2_REDIS+gameID, config.REDIS_EXPIRE)
// }

// 更新題目進行題數資料(更新qa_round=1)
// models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateQARound(false, gameID, "1")

// if status == "start" {
// 	// 新增測試用戶分數紀錄(測試用，999筆分數緩存資料)
// 	for i := 1; i < 1000; i++ {
// 		// redis
// 		h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 			10)
// 		h.redisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
// 			float64(i)+0.1)
// 	}
// }

// 新增測試用戶分數紀錄(測試用，999筆分數緩存資料)
// for i := 1; i < 1000; i++ {
// 	// redis
// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 		10)
// 	h.redisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
// 		float64(i)+0.1)
// }

// gameModel, _ := models.DefaultGameModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindGame(true, gameID)

// gameModel, err1   = models.DefaultGameModel().SetDbConn(h.dbConn).Find(gameID)
// var (
// 	editGameModel models.EditGameModel
// )

// GamePrizeParam 場次、中獎人員相關參數資訊
// type GamePrizeParam struct {
// 	Game GameModel // 場次資訊
// 	// PrizeStaffs []PrizeStaffModel // 中獎人員
// 	Error string `json:"error" example:"error message"` // 錯誤訊息
// }

// 刪除redis所有遊戲相關資訊
// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊
// 刪除redis所有遊戲相關資訊
// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊

// status            string
// second            int64

// status = gameModel.GameStatus
// second = gameModel.GameSecond

// fmt.Println("玩家端遊戲狀態: ", status)

// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, "", strconv.Itoa(int(gameModel.Second)),
// 		"", result.Game.GameStatus)
// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, "", "", "", result.Game.GameStatus)
// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, strconv.Itoa(int(gameModel.GameRound+1)),
// 		"", "0", result.Game.GameStatus)
// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, "", strconv.Itoa(int(result.Game.GameSecond)),
// 		"", result.Game.GameStatus)
// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, "", "",
// 		"", result.Game.GameStatus)
// err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateGameStatus(true, gameID, strconv.Itoa(int(gameModel.GameRound+1)), "",
// 		"", result.Game.GameStatus)

// getGameModel 取得GameModel
// func (h *Handler) getGameModel(url, activityID string) (gameModel models.GameModel, err error) {
// 	if strings.Contains(url, "redpack") {
// 		gameModel, err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 			FindGameOpen(activityID, "搖紅包")
// 	} else if strings.Contains(url, "ropepack") {
// 		gameModel, err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 			FindGameOpen(activityID, "套紅包")
// 	}
// 	return gameModel, err
// }

// @Summary 回傳搖紅包主持端展示的遊戲場次資訊(玩家端)
// @Tags Redpack Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Success 200 {array} WebsocketMessage
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/redpack/status/open [get]
// func (h *Handler) GameOpenWebsocket(ctx *gin.Context) {
// 	var (
// 		activityID        = ctx.Query("activity_id")
// 		wsConn, conn, err = NewWebsocketConn(ctx)
// 		gameModel, err1   = h.getGameModel(ctx.Request.URL.Path, activityID)
// 		gameID            = gameModel.GameID
// 	)
// 	defer wsConn.Close()
// 	defer conn.Close()
// 	fmt.Println("開啟主持端展示的遊戲場次資訊(玩家端)ws")
// 	// 判斷活動、遊戲資訊是否有效
// 	if err != nil || err1 != nil || activityID == "" {
// 		b, _ := json.Marshal(WebsocketMessage{
// 			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
// 		conn.WriteMessage(b)
// 		return
// 	}

// 	// 初次傳送展示中的場次資訊
// 	fmt.Println("初次傳送展示中的場次ID: ", gameID)
// 	b, _ := json.Marshal(WebsocketMessage{
// 		Game: GameModel{
// 			GameID: gameID,
// 		},
// 	})
// 	if err := conn.WriteMessage(b); err != nil {
// 		return
// 	}

// 	go func() {
// 		for {
// 			// 即時展示中的場次資訊
// 			if gameModel, _ = h.getGameModel(ctx.Request.URL.Path, activityID); gameModel.GameID != "" &&
// 				gameModel.GameID != gameID { // 如果人數或輪次更新才會回傳前端
// 				gameID = gameModel.GameID
// 				fmt.Println("切換展示中的場次ID: ", gameID)

// 				b, _ := json.Marshal(WebsocketMessage{
// 					Game: GameModel{
// 						GameID: gameID,
// 					},
// 				})
// 				if err := conn.WriteMessage(b); err != nil {
// 					return
// 				}
// 			}

// 			if conn.isClose { // ws關閉
// 				return
// 			}
// 			time.Sleep(time.Second * 2)
// 		}
// 	}
