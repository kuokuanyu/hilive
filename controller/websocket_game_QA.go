package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// HostQAWriteParam 主持端快問快答回傳前端訊息
type HostQAWriteParam struct {
	Game  GameModel // 場次資訊
	A     int64     `json:"A" example:"1"`                 // 答案A人數
	B     int64     `json:"B" example:"2"`                 // 答案B人數
	C     int64     `json:"C" example:"3"`                 // 答案C人數
	D     int64     `json:"D" example:"4"`                 // 答案D人數
	Total int64     `json:"Total" example:"10"`            // 總答題人數
	Error string    `json:"error" example:"error message"` // 錯誤訊息
}

// #####接下來不需要#####
// GuestQAReadParam 玩家端快問快答接收前端訊息
// type GuestQAReadParam struct {
// User   UserModel // 用戶資訊
// Game   GameModel // 遊戲資訊
// Answer string `json:"answer" example:"A、B、C、D"` // 玩家選擇答案

// QARound     int64  `json:"qa_round" example:"1"`       //遊戲題數
// QAOption    string `json:"qa_option" example:"A"`      //遊戲選項
// Score       int64  `json:"score" example:"100"`        //遊戲分數
// OriginScore int64  `json:"origin_score" example:"100"` // 原始分數
// AddScore    int64  `json:"add_score" example:"100"`    // 加成分數

// }

// #####到這都不需要#####

// @Summary 主持端取得總答題人數、各選項答題人數資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(QA)
// @Success 200 {array} HostQAWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/host/QA [get]
func (h *Handler) HostGameQAWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		isOPen            = true
	)
	// fmt.Println("開啟主持端取得總答題人數、各選項答題人數資訊ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 開啟時傳送訊息
		var (
			gameModel = h.getGameInfo(gameID, game)
		)

		dataMap, _ := h.redisConn.HashGetAllCache(config.QA_REDIS + gameID)

		conn.WriteMessage([]byte(utils.JSON(
			HostQAWriteParam{
				Game: GameModel{
					QARound: gameModel.QARound,
				},
				A:     utils.GetInt64FromStringMap(dataMap, "A", 0),
				B:     utils.GetInt64FromStringMap(dataMap, "B", 0),
				C:     utils.GetInt64FromStringMap(dataMap, "C", 0),
				D:     utils.GetInt64FromStringMap(dataMap, "D", 0),
				Total: utils.GetInt64FromStringMap(dataMap, "Total", 0)})))

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribes(context,
		[]string{config.CHANNEL_QA_REDIS + gameID,
			config.CHANNEL_GAME_REDIS + gameID,
		}, func(channel, message string) {
			// 資料改變時，回傳最新資料至前端
			var (
				gameModel = h.getGameInfo(gameID, game)
			)

			dataMap, _ := h.redisConn.HashGetAllCache(config.QA_REDIS + gameID)

			conn.WriteMessage([]byte(utils.JSON(
				HostQAWriteParam{
					Game: GameModel{
						QARound:    gameModel.QARound,
						GameStatus: gameModel.GameStatus,
					},
					A:     utils.GetInt64FromStringMap(dataMap, "A", 0),
					B:     utils.GetInt64FromStringMap(dataMap, "B", 0),
					C:     utils.GetInt64FromStringMap(dataMap, "C", 0),
					D:     utils.GetInt64FromStringMap(dataMap, "D", 0),
					Total: utils.GetInt64FromStringMap(dataMap, "Total", 0)})))
		})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉主持端取得總答題人數、各選項答題人數資訊ws")
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_QA_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_REDIS + gameID)
			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			gameModel = h.getGameInfo(gameID)
// 			// a, b, c, d, total int64
// 		)

// 		dataMap, _ := h.redisConn.HashGetAllCache(config.QA_REDIS + gameID)

// 		conn.WriteMessage([]byte(utils.JSON(
// 			HostQAWriteParam{
// 				Game: GameModel{
// 					QARound: gameModel.QARound,
// 				},
// 				A:     utils.GetInt64FromStringMap(dataMap, "A", 0),
// 				B:     utils.GetInt64FromStringMap(dataMap, "B", 0),
// 				C:     utils.GetInt64FromStringMap(dataMap, "C", 0),
// 				D:     utils.GetInt64FromStringMap(dataMap, "D", 0),
// 				Total: utils.GetInt64FromStringMap(dataMap, "Total", 0)})))
// 		// }

// 		// 壓力測試用
// 		// if err := conn.WriteMessage([]byte(utils.JSON(
// 		// 	HostQAWriteParam{}))); err != nil {
// 		// 	return
// 		// }
// 		// 壓力測試結束
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// @Summary 玩家端更新總答題人數、各選項答題人數、答題紀錄、遊戲分數資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(QA)
// @param body body UserGameParam true "user、game param"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/guest/QA [get]
func (h *Handler) GuestGameQAWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            UserGameParam
	)
	// fmt.Println("開啟玩家端更新總答題人數、各選項答題人數、答題紀錄資訊ws")
	defer wsConn.Close()
	defer conn.Close()
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(UserGameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result UserGameParam
			option string
		)

		// var option string
		data, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉玩家端更新總答題人數、各選項答題人數、答題紀錄資訊ws")
			conn.Close()
			return
		}
		// fmt.Println("收到QA訊息")

		json.Unmarshal(data, &result)

		// fmt.Println("更新答題紀錄")

		if result.User.UserID != "" &&
			result.Game.GameStatus != "" && result.Game.QARound != 0 &&
			(result.Game.QAOption == "0" || result.Game.QAOption == "1" ||
				result.Game.QAOption == "2" || result.Game.QAOption == "3") {

			// fmt.Println("result.User.UserID: ", result.User.UserID)
			// fmt.Println("result.Game.GameStatus: ", result.Game.GameStatus)
			// fmt.Println(" result.Game.QARound: ", result.Game.QARound)
			// fmt.Println("result.Game.QAOption: ", result.Game.QAOption)
			// fmt.Println("原始分數: ", result.Game.OriginScore)
			// fmt.Println("加成分數: ", result.Game.AddScore)
			// fmt.Println("第一分數: ", result.Game.GameScore)
			// fmt.Println("第二分數: ", result.Game.GameScore2)
			// fmt.Println("答對題數: ", result.Game.GameCorrect)

			if result.Game.QAOption == "0" {
				option = "A"
			} else if result.Game.QAOption == "1" {
				option = "B"
			} else if result.Game.QAOption == "2" {
				option = "C"
			} else if result.Game.QAOption == "3" {
				option = "D"
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

			// 更新redis中score_2資訊
			if err = models.DefaultScoreModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateScore2(true, gameID, result.User.UserID, result.Game.GameScore2); err != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: "錯誤: 更新人員第二分數資料發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 更新redis中答對題數資訊
			if err = models.DefaultScoreModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateCorrect(true, gameID, result.User.UserID, result.Game.GameCorrect); err != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: "錯誤: 更新人員答對題數資料發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 更新用戶答題紀錄
			models.DefaultGameQARecordModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Update(true, models.EditGameQARecordModel{
					UserID:      result.User.UserID,
					GameID:      gameID,
					QARound:     result.Game.QARound,
					QAOption:    option,
					Score:       result.Game.GameScore,
					OriginScore: result.Game.OriginScore,
					AddScore:    result.Game.AddScore,
				})

			// 遞增redis裡總答題人數、各選項答題人數資訊
			h.redisConn.HashIncrCache(config.QA_REDIS+gameID, "Total")
			h.redisConn.HashIncrCache(config.QA_REDIS+gameID, option)

			// 設置過期時間
			// h.redisConn.SetExpire(config.QA_REDIS+gameID, config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_QA_REDIS+gameID, "修改資料")

			rank := h.redisConn.ZSetRevRank(config.SCORES_REDIS+gameID,
				result.User.UserID) // 排名資料
			// fmt.Println("排名: ", rank+1)
			b, _ := json.Marshal(UserGameParam{
				Game: GameModel{
					GameRank: rank + 1,
				}})
			conn.WriteMessage(b)
		}
	}
}

// if result.User.UserID == "" ||
// 	result.Game.GameStatus == "" || result.Game.QARound == 0 ||
// 	(result.Game.QAOption != "0" && result.Game.QAOption != "1" &&
// 		result.Game.QAOption != "2" && result.Game.QAOption != "3") {
// 	b, _ := json.Marshal(GameModel{
// 		Error: "錯誤: 無法取得答題資訊"})
// conn.WriteMessage(b)
// 	return
// }

// #####接下來不執行#####
// if result.Answer != "0" && result.Answer != "1" &&
// 	result.Answer != "2" && result.Answer != "3" {
// 	b, _ := json.Marshal(GameModel{
// 		Error: "錯誤: 無法取得答題資訊"})
// 	conn.WriteMessage(b)
// 	return
// }
// if result.Answer == "0" {
// 	answer = "A"
// } else if result.Answer == "1" {
// 	answer = "B"
// } else if result.Answer == "2" {
// 	answer = "C"
// } else if result.Answer == "3" {
// 	answer = "D"
// }
// #####到這都不執行#####

// #####接下來不執行#####
// 遞增redis裡總答題人數、各選項答題人數資訊
// h.redisConn.HashIncrCache(config.QA_REDIS+gameID, "Total")
// h.redisConn.HashIncrCache(config.QA_REDIS+gameID, answer)
// #####到這都不執行#####

// 接收到結算的遊戲狀態時，關閉ws
// if result.Game.GameStatus == "calculate" {
// 	fmt.Println("遊戲結算，關閉玩家端更新總答題人數、各選項答題人數資訊ws")
// 	return
// }

// 更新redis中score_2資訊
// if err = models.DefaultScoreModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateScore2(true, gameID, result.User.UserID, result.Game.GameScore2); err != nil {
// 	b, _ := json.Marshal(UserGameParam{
// 		Error: "錯誤: 更新人員第二分數資料發生問題"})
// 	conn.WriteMessage(b)
// 	return
// }

// 回傳遊戲狀態訊息給前端
// b, _ := json.Marshal(GameParam{
// 	Game: GameModel{GameStatus: status},
// })
// conn.WriteMessage(b)
