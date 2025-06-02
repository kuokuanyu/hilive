package controller

import (
	"context"
	"encoding/json"
	"hilive/modules/config"
	"time"

	"github.com/gin-gonic/gin"
)

// ResetParam 重新整理頁面相關參數資訊
type ResetParam struct {
	Reset bool   `json:"reset" example:"true"`          // 頁面是否需要重整
	Error string `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 即時判斷頁面是否需要重整ws api
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game query string true "game" Enums(signname, general, threed)
// @Success 200 {array} ResetParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/activity/reset [get]
func (h *Handler) ActivityResetWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// activityModel, err1 = h.getActivityInfo(true, activityID)
		// oldEditTimes        int64
		channel string
		// isOPen  = true
	)

	// redis key名稱
	if game == "signname" {
		// 活動
		channel = config.CHANNEL_SIGNNAME_EDIT_TIMES_REDIS + activityID
	} else if game == "general" {
		// 訊息
		channel = config.CHANNEL_GENERAL_EDIT_TIMES_REDIS + activityID
	} else if game == "threed" {
		// 提問
		channel = config.CHANNEL_THREED_EDIT_TIMES_REDIS + activityID
	}

	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || game == "" {
		b, _ := json.Marshal(ResetParam{
			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// if isOPen {
	// 	// 開啟時傳送訊息
	// 	b, _ := json.Marshal(ResetParam{
	// 		Reset: true,
	// 	})
	// 	conn.WriteMessage(b)
	// 	isOPen = false
	// }

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, channel, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		b, _ := json.Marshal(ResetParam{
			Reset: true,
		})
		conn.WriteMessage(b)
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉reset ws api", err)
			// 取消訂閱
			h.redisConn.Unsubscribe(channel)
			conn.Close()
			return
		}
	}
}

// @Summary 即時判斷頁面是否需要重整ws api
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, monopoly, draw_numbers, QA, tugofwar, bingo)
// @Param role query string true "role" Enums(host, guest)
// @Success 200 {array} ResetParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/reset [get]
func (h *Handler) GameResetWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		role              = ctx.Query("role")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// gameModel                = h.getGameInfo(gameID)
		// oldEditTimes             = gameModel.EditTimes
		// isReset, isClose, isOpen bool
	)
	// log.Println("開啟遊戲reset ws api: ", activityID, gameID,game)

	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" ||
		game == "" || role == "" {
		b, _ := json.Marshal(ResetParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// if isOPen {
	// 	// 開啟時傳送訊息
	// 	isOPen = false
	// 	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		if role == "host" {
			// 主持端
		} else if role == "guest" {
			// 玩家端
			// 主持端必須接收到編輯訊息後並且遊戲為關閉後再打開的狀態才執行刷新頁面
			// 延後5秒重置
			time.Sleep(time.Second * 5)
		}

		b, _ := json.Marshal(ResetParam{
			Reset: true,
		})
		conn.WriteMessage(b)
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_EDIT_TIMES_REDIS + gameID)
			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			gameModel    = h.getGameInfo(gameID) // 遊戲資訊
// 			newEditTimes = gameModel.EditTimes
// 		)

// 		// 判斷後端是否更新遊戲設置
// 		if newEditTimes != oldEditTimes {
// 			// fmt.Println("遊戲設置更新", newEditTimes, oldEditTimes)
// 			isReset = true

// 			oldEditTimes = newEditTimes
// 		} else {
// 			// fmt.Println("遊戲設置未更新", newEditTimes, oldEditTimes)
// 		}

// 		if role == "host" {
// 			// 主持端
// 			if isReset == true {
// 				b, _ := json.Marshal(ResetParam{
// 					Reset: true,
// 				})
// 				conn.WriteMessage(b)

// 				isReset = false
// 			}
// 		} else if role == "guest" {
// 			// 玩家端
// 			// 必須接收到編輯後並且遊戲為關閉後再打開的狀態才執行刷新頁面
// 			if isReset == true {
// 				if gameModel.GameStatus == "close" {
// 					isClose = true
// 				}

// 				if isClose == true && gameModel.GameStatus == "open" {
// 					isOpen = true
// 				}

// 				if isClose == true && isOpen == true {
// 					b, _ := json.Marshal(ResetParam{
// 						Reset: true,
// 					})
// 					conn.WriteMessage(b)

// 					isReset = false
// 				}
// 			}
// 		}

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			activityModel, _ = h.getActivityInfo(true, activityID) // 活動資訊
// 			newEditTimes     int64
// 		)

// 		if game == "signname" {
// 			newEditTimes = activityModel.SignnameEditTimes
// 		} else if game == "general" {
// 			newEditTimes = activityModel.GeneralEditTimes
// 		} else if game == "threed" {
// 			newEditTimes = activityModel.ThreedEditTimes
// 		}

// 		// 判斷後端是否更新活動設置
// 		if newEditTimes != oldEditTimes {
// 			// fmt.Println("活動設置更新", newEditTimes, oldEditTimes)
// 			b, _ := json.Marshal(ResetParam{
// 				Reset: true,
// 			})
// 			conn.WriteMessage(b)

// 			oldEditTimes = newEditTimes
// 		} else {
// 			// log.Println("活動設置未更新", newEditTimes, oldEditTimes)
// 			// b, _ := json.Marshal(ResetParam{
// 			// 	Reset: true,
// 			// })
// 			// conn.WriteMessage(b)
// 		}

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###
