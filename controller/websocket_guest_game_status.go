package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"

	"github.com/gin-gonic/gin"
)

// @Summary 玩家端遊戲狀態，判斷資料表遊戲狀態變化
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, lottery, monopoly, QA, tugofwar, bingo, 3DGachaMachine, vote)
// @Success 200 {array} GameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/status/guest [get]
func (h *Handler) GuestGameStatusWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		isOPen            = true
		channels          = []string{config.CHANNEL_GUEST_GAME_STATUS_REDIS + gameID} // 遊戲資訊
	)

	// log.Println("開啟玩家端遊戲狀態ws", activityID, gameID, game)
	defer wsConn.Close()
	defer conn.Close()
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// 開啟時傳送訊息
	if isOPen {
		var (
			gameModel  = h.getGameInfo(gameID, game) // 即時人數、輪次
			prizeStaff models.PrizeStaffModel        // 中獎紀錄
			// amount     int64
			isShake bool
		)

		// 扭蛋機遊戲
		if game == "3DGachaMachine" {
			if gameModel.GameStatus == "calculate" {
				// 結算狀態
				var recordsJson []string // 中獎人員json資料

				// 從redis中取得中獎紀錄(winning_staffs_)
				recordsJson, _ = h.redisConn.ListRange(config.WINNING_STAFFS_REDIS+gameID, 0, 1)

				if len(recordsJson) > 0 {
					for _, record := range recordsJson {
						// 解碼
						json.Unmarshal([]byte(record), &prizeStaff)
					}
				}
			} else if gameModel.GameStatus == "gaming" {
				// 遊戲中狀態

				// 是否搖晃資料判斷
				if gameModel.IsShake == "true" {
					isShake = true
				} else {
					isShake = false
				}
			}
		}

		// 測試
		// if game == "QA" {
		// 	log.Println("玩家端遊戲狀態: ", gameModel.GameStatus)
		// }

		b, _ := json.Marshal(GameParam{
			Game: GameModel{
				ControlGameStatus:   gameModel.ControlGameStatus,
				GameStatus:          gameModel.GameStatus,
				GameRound:           gameModel.GameRound,
				GameSecond:          gameModel.GameSecond,
				QARound:             gameModel.QARound,
				BingoRound:          gameModel.BingoRound,
				LeftTeamGameAttend:  gameModel.LeftTeamGameAttend,
				RightTeamGameAttend: gameModel.RightTeamGameAttend,
				// PrizeAmount:         amount,
				IsShake: isShake,
			},
			PrizeStaff: prizeStaff,
		})
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	if game == "3DGachaMachine" {
		channels = append(channels, config.CHANNEL_WINNING_STAFFS_REDIS+gameID)
	}

	// 啟用redis訂閱
	go h.redisConn.Subscribes(context, channels, func(channel, message string) {
		// log.Println("訂閱者收到資料變動的訊息")

		var (
			gameModel  = h.getGameInfo(gameID, game) // 即時人數、輪次
			prizeStaff models.PrizeStaffModel        // 中獎紀錄
			// amount     int64
			isShake bool
		)

		// 扭蛋機遊戲
		if game == "3DGachaMachine" {
			if gameModel.GameStatus == "calculate" {
				// 結算狀態
				var recordsJson []string // 中獎人員json資料

				// 從redis中取得中獎紀錄(winning_staffs_)
				recordsJson, _ = h.redisConn.ListRange(config.WINNING_STAFFS_REDIS+gameID, 0, 1)

				if len(recordsJson) > 0 {
					for _, record := range recordsJson {
						// 解碼
						json.Unmarshal([]byte(record), &prizeStaff)
					}
				}
			} else if gameModel.GameStatus == "gaming" {
				// 遊戲中狀態

				// 是否搖晃資料判斷
				if gameModel.IsShake == "true" {
					isShake = true
				} else {
					isShake = false
				}
			}
		}

		// 測試
		// if game == "QA" {
		// 	log.Println("玩家端遊戲狀態: ", gameModel.GameStatus)
		// }

		// log.Println("玩家端遊戲狀態: ", gameModel.GameStatus)

		b, _ := json.Marshal(GameParam{
			Game: GameModel{
				ControlGameStatus:   gameModel.ControlGameStatus,
				GameStatus:          gameModel.GameStatus,
				GameRound:           gameModel.GameRound,
				GameSecond:          gameModel.GameSecond,
				QARound:             gameModel.QARound,
				BingoRound:          gameModel.BingoRound,
				LeftTeamGameAttend:  gameModel.LeftTeamGameAttend,
				RightTeamGameAttend: gameModel.RightTeamGameAttend,
				// PrizeAmount:         amount,
				IsShake: isShake,
			},
			PrizeStaff: prizeStaff,
		})
		conn.WriteMessage(b)
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉玩家端遊戲狀態ws")

			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GUEST_GAME_STATUS_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_WINNING_STAFFS_REDIS + gameID)

			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			gameModel  = h.getGameInfo(gameID) // 即時人數、輪次
// 			prizeStaff models.PrizeStaffModel  // 中獎紀錄
// 			// amount     int64
// 			isShake bool
// 		)

// 		// 扭蛋機遊戲
// 		if game == "3DGachaMachine" {
// 			if gameModel.GameStatus == "calculate" {
// 				// 結算狀態
// 				var recordsJson []string // 中獎人員json資料

// 				// 從redis中取得中獎紀錄(winning_staffs_)
// 				recordsJson, _ = h.redisConn.ListRange(config.WINNING_STAFFS_REDIS+gameID, 0, 1)

// 				if len(recordsJson) > 0 {
// 					for _, record := range recordsJson {
// 						// 解碼
// 						json.Unmarshal([]byte(record), &prizeStaff)
// 					}
// 				}
// 			} else if gameModel.GameStatus == "gaming" {
// 				// 遊戲中狀態

// 				// 是否搖晃資料判斷
// 				if gameModel.IsShake == "true" {
// 					isShake = true
// 				} else {
// 					isShake = false
// 				}
// 			}
// 		}

// 		b, _ := json.Marshal(GameParam{
// 			Game: GameModel{
// 				ControlGameStatus:   gameModel.ControlGameStatus,
// 				GameStatus:          gameModel.GameStatus,
// 				GameRound:           gameModel.GameRound,
// 				GameSecond:          gameModel.GameSecond,
// 				QARound:             gameModel.QARound,
// 				BingoRound:          gameModel.BingoRound,
// 				LeftTeamGameAttend:  gameModel.LeftTeamGameAttend,
// 				RightTeamGameAttend: gameModel.RightTeamGameAttend,
// 				// PrizeAmount:         amount,
// 				IsShake: isShake,
// 			},
// 			PrizeStaff: prizeStaff,
// 		})
// 		conn.WriteMessage(b)

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}

// 		// 扭蛋機遊戲
// 		if game == "3DGachaMachine" {
// 			time.Sleep(time.Second * 1 / 10)
// 		} else {
// 			// 其他遊戲
// 			time.Sleep(time.Second * 1 / 2)
// 		}

// 		// time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

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
// 	}()

// 	for {
// 		_, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("關閉主持端展示的遊戲場次資訊(玩家端)ws")
// 			conn.Close()
//
