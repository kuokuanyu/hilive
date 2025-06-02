package controller

import (
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 即時更新遊戲分數(玩家端)，接收並更新分數
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(whack_mole, monopoly, tugofwar, bingo)
// @param body body UserGameParam true "game param"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/score [get]
func (h *Handler) ScoreWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            UserGameParam
	)
	// log.Println("開啟即時更新遊戲分數(玩家端)ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		// fmt.Println("??????")
		b, _ := json.Marshal(GameModel{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result UserGameParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉即時更新遊戲分數(玩家端)ws")
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)
		// if result.User.UserID == "" || result.Game.GameStatus == "" {
		// 	b, _ := json.Marshal(UserGameParam{
		// 		Error: "錯誤: 無法辨識用戶、遊戲狀態資訊"})
		// 	conn.WriteMessage(b)
		// 	return
		// }

		// log.Println("收到?", result.User.UserID, result.Game.GameStatus, result.Game.UserBingoRound)
		if result.User.UserID != "" && result.Game.GameStatus != "" {

			if game == "tugofwar" {
				// 拔河遊戲
				// log.Println("分數: ", result.User.UserID, result.Game.GameStatus, result.Game.GameScore)
				// 判斷玩家分數，如果接收前端資料為玩家分數為0時，查詢redis的分數資料是否正確
				if result.Game.GameScore == 0 {
					score := h.redisConn.ZSetIntScore(config.SCORES_REDIS+gameID, result.User.UserID) // 分數資料

					// 更新分數資料
					result.Game.GameScore = score

					if result.Game.GameScore != 0 {
						// 如果確定redis中分數不為0時，回傳資料給前端
						b, _ := json.Marshal(UserGameParam{
							Game: GameModel{
								IsZero:    true,
								GameScore: result.Game.GameScore,
							},
						})
						conn.WriteMessage(b)
					}
				}

				// 更新分數資料(score)
				if err = models.DefaultScoreModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateScore(true, gameID, result.User.UserID, result.Game.GameScore); err != nil {
					b, _ := json.Marshal(UserGameParam{
						Error: "錯誤: 更新人員計分資料發生問題"})
					conn.WriteMessage(b)
					return
				}

			} else {
				// 其他遊戲
				// 更新分數資料(score)
				if err = models.DefaultScoreModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateScore(true, gameID, result.User.UserID, result.Game.GameScore); err != nil {
					b, _ := json.Marshal(UserGameParam{
						Error: "錯誤: 更新人員計分資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			}

			if game == "whack_mole" || game == "monopoly" {
				// log.Println("競技遊戲接收分數資料: ", result.Game.GameScore, "輪次: ", result.Game.GameRound)

				// 更新redis中score_2資訊
				if err = models.DefaultScoreModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateScore2(true, gameID, result.User.UserID, result.Game.GameScore2); err != nil {
					b, _ := json.Marshal(UserGameParam{
						Error: "錯誤: 更新人員第二分數資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			} else if game == "bingo" {
				// 更新玩家是否即將中獎資料(redis)
				if err = models.DefaultScoreModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateGoingBingo(true, gameID, result.User.UserID, result.Game.GoingBingo); err != nil {
					b, _ := json.Marshal(UserGameParam{
						Error: err.Error()})
					conn.WriteMessage(b)
					return
				}

				if result.Game.UserBingoRound != 0 {
					gameModel := h.getGameInfo(gameID, game) // 單輪中獎人數

					// 玩家賓果，上鎖判斷是否中獎
					// 計算剩餘獎品數
					for l := 0; l < MaxRetries; l++ {
						// 上鎖
						ok, _ := h.acquireLock(config.BINGO_LOCK_REDIS+gameID, LockExpiration)
						if ok == "OK" {
							// 釋放鎖
							// defer h.releaseLock(config.BINGO_LOCK_REDIS + gameID)

							var (
								isBingo  bool
								users, _ = h.redisConn.ZSetRange(config.GAME_BINGO_USER_REDIS+gameID, 0, 0) // 取得已賓果玩家資料
							)

							// 判斷中獎人數是否已達上限
							if int(gameModel.RoundPrize)-len(users) > 0 {
								// 未達上限，中獎，更新玩家賓果資料
								// 更新玩家賓果的回合數(redis)
								if err = models.DefaultScoreModel().
									SetConn(h.dbConn, h.redisConn, h.mongoConn).
									UpdateBingoRound(true, gameID, result.User.UserID,
										result.Game.UserBingoRound); err != nil {
									b, _ := json.Marshal(UserGameParam{
										Error: err.Error()})
									conn.WriteMessage(b)

									// 釋放鎖
									h.releaseLock(config.BINGO_LOCK_REDIS + gameID)
									return
								}

								// 中獎
								isBingo = true

								// 測試用
								// users, _ = h.redisConn.ZSetRange(config.GAME_BINGO_USER_REDIS+gameID, 0, 0) // 取得已賓果玩家資料
								// log.Println("賓果人數: ", len(users))
							}
							// else {
							// 人數已達上限
							// 遞增未中獎人數
							// h.redisConn.IncrCache("unlucky")
							// amount, _ := h.redisConn.GetCache("unlucky")
							// log.Println("未中獎人數: ", amount)
							// }

							// 回傳是否中獎資料至前端
							b, _ := json.Marshal(UserGameParam{
								Game: GameModel{
									IsBingo: isBingo,
								},
							})
							conn.WriteMessage(b)

							// 釋放鎖
							h.releaseLock(config.BINGO_LOCK_REDIS + gameID)
							break
						}

						// 鎖被佔用，稍微延遲後重試
						time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
					}

				}
			}

			// 接收到結算的遊戲狀態時，關閉ws
			if result.Game.GameStatus == "calculate" {
				// fmt.Println("有執行?")
				// fmt.Println("遊戲結算，關閉即時更新遊戲分數(玩家端)ws")
				return
			}
		}
	}
}

// #####接下來QA不執行#####
// if game == "QA" {
// 	// 更新redis中答對題數資訊
// 	if err = models.DefaultScoreModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateCorrect(true, gameID, result.User.UserID, result.Game.GameCorrect); err != nil {
// 		b, _ := json.Marshal(UserGameParam{
// 			Error: "錯誤: 更新人員答對題數資料發生問題"})
// 		conn.WriteMessage(b)
// 		return
// 	}

// 	rank := h.redisConn.ZSetRevRank(config.SCORES_REDIS+gameID,
// 		result.User.UserID) // 排名資料
// 	b, _ := json.Marshal(GameParam{
// 		Game: GameModel{
// 			GameRank: rank + 1,
// 		}})
// 	conn.WriteMessage(b)
// }
// #####到這都不執行#####

// 	models.EditScoreModel{
// 	GameID: gameID,
// 	UserID: result.User.UserID,
// 	// Round:  strconv.Itoa(int(result.Game.GameRound)),
// 	Score:  strconv.Itoa(int(result.Game.GameScore)),
// }
