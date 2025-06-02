package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"

	"github.com/gin-gonic/gin"
)

// @Summary 主持端即時更新賓果遊戲號碼
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(bingo)
// @param body body GameParam true "game param"
// @Success 200 {array} GameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/number/host [get]
func (h *Handler) HostGameBingoWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            GameParam
		// count int64
	)
	// fmt.Println("開啟主持人端即時更新賓果遊戲號碼ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result GameParam
		)

		data, err := conn.ReadMessage()
		// fmt.Println("收到主持端賓果號碼訊息")
		// 關閉網頁或連線
		if err != nil {
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		// #####加入測試資料#####start
		// count++
		// if count == 1 {
		// 	// 賓果人員(第一輪中獎人員，前5名)
		// 	for i := 1; i <= 1000; i++ {
		// 		// 加入分數資料
		// 		h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 1)          // 賓果連線數資料
		// 		h.redisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID, strconv.Itoa(i), 1) // 完成賓果的回合數
		// 	}

		// 	// 即將中獎人員
		// 	// for i := 501; i <= 999; i++ {
		// 	// 	// 加入分數資料
		// 	// 	h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 0) // 賓果連線數資料
		// 	// 	h.redisConn.HashSetCache(config.GAME_BINGO_USER_GOING_BINGO+gameID, strconv.Itoa(i), "true")
		// 	// }
		// }
		// else if count == 5 {
		// 	// 賓果人員(第二輪中獎人員，後10名)
		// 	for i := 6; i <= 7; i++ {
		// 		// 加入分數資料
		// 		h.redisConn.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i), 1)          // 賓果連線數資料
		// 		h.redisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID, strconv.Itoa(i), 15) // 完成賓果的回合數
		// 	}
		// }
		// #####加入測試資料#####end

		// if result.Game.Number == 0 {
		// 	b, _ := json.Marshal(GameParam{
		// 		Error: "錯誤: 請傳遞賓果遊戲抽獎號碼"})
		// conn.WriteMessage(b)
		// 	return
		// }

		if len(result.Game.Numbers) > 0 {
			// 更新賓果遊戲抽獎號碼資料(redis)
			models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateBingoNumber(true, gameID, result.Game.Numbers)

			// 遞增賓果遊戲進行回合數資料(redis)
			models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateBingoRound(true, gameID, "+1")
		}
	}
}

// @Summary 玩家端即時回傳賓果遊戲所有號碼、更新玩家號碼排序
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(bingo)
// @param body body UserGameParam true "UserGameParam param"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/number/guest [get]
func (h *Handler) GuestGameBingoWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            UserGameParam
		isOPen = true
		// userID string
	)
	// fmt.Println("開啟玩家端即時回傳賓果遊戲號碼、更新玩家號碼排序ws")
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

		// 查詢賓果遊戲所有抽獎號碼(redis)
		numbers, _ := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindBingoNumbers(true, gameID)

		// log.Println("賓果遊戲玩家接收號碼資料: ", numbers, "用戶: ", userID)

		b, _ := json.Marshal(UserGameParam{
			Game: GameModel{
				Numbers: numbers,
			},
		})
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		// 查詢賓果遊戲所有抽獎號碼(redis)
		numbers, _ := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindBingoNumbers(true, gameID)

		// log.Println("賓果遊戲玩家接收號碼資料: ", numbers, "用戶: ", userID)

		b, _ := json.Marshal(UserGameParam{
			Game: GameModel{
				Numbers: numbers,
			},
		})
		conn.WriteMessage(b)
	})

	for {
		var (
			result UserGameParam
		)

		data, err := conn.ReadMessage()
		// 關閉網頁或連線
		if err != nil {
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_BINGO_NUMBER_REDIS + gameID)
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.User.UserID != "" && len(result.Game.Numbers) != 0 {
			// 更新玩家的號碼排序
			models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateUserNumber(true, gameID, result.User.UserID, result.Game.Numbers)

			// userID = result.User.UserID
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		// var numbers = make([]string, 0)

// 		// 查詢賓果遊戲所有抽獎號碼(redis)
// 		numbers, _ := models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 			FindBingoNumbers(true, gameID)
// 		// fmt.Println("玩家端接收號碼: ", numbers)

// 		b, _ := json.Marshal(UserGameParam{
// 			Game: GameModel{
// 				Numbers: numbers,
// 			},
// 		})
// 		conn.WriteMessage(b)

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// if result.User.UserID == "" || len(result.Game.Numbers) == 0 {
// 	b, _ := json.Marshal(GameParam{
// 		Error: "錯誤: 無法辨識玩家資訊、號碼排序資訊"})
// 	conn.WriteMessage(b)
// 	return
// }
// fmt.Println("收到玩家端賓果號碼排序: ",
