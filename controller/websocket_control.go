package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HostControlParam 主持端遠端遙控參數
type HostControlParam struct {
	Game          GameModel // 場次資訊
	ChannelID     string    `json:"channel_id" example:"channel_id"` // 頻道id
	ChannelStatus string    `json:"channel_status" example:"open"`   // 頻道狀態
	// GameID    string    `json:"game_id" example:"game_id"`       // 遊戲id
	Page   string `json:"page" example:"massage"`        // 點擊頁面
	Action string `json:"action" example:"next page"`    // 點擊動作
	Error  string `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 主持端遙控螢幕
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param device query string true "device" Enums(pc, mobile)
// @Param role query string true "role" Enums(host, guest)
// @Param channel_id query string true "channel ID"
// @param body body HostControlParam true "page、action"
// @Success 200 {array} HostControlParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/host/control [get]
func (h *Handler) HostControlWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		device            = ctx.Query("device")
		channelID         = "channel_1" // channel_1、ctx.Query("session_id")
		role              = ctx.Query("role")
		redisName         = config.HOST_CONTROL_REDIS + activityID + "_" + channelID // host_control_id_channel_n
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result               HostControlParam
		page, action, status string
		isOPen               = true
	)
	// log.Println("開啟主持端遙控螢幕ws")
	defer wsConn.Close()
	defer conn.Close()

	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || channelID == "" ||
		(device != "pc" && device != "mobile") {
		b, _ := json.Marshal(HostControlParam{
			Error: "錯誤: 無法辨識活動、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if device == "mobile" && role == "host" {
		// 主持人操作手機端遙控頁面，重置page、action(hash方式)
		values := []interface{}{redisName}
		values = append(values, "page", "")
		values = append(values, "action", "")
		values = append(values, "game_status", "")
		values = append(values, "game_id", "")
		values = append(values, "game_id_2", "")

		// 更新主持端遙控資訊(host_control_活動ID_sessionID)
		if err := h.redisConn.HashMultiSetCache(values); err != nil {
			b, _ := json.Marshal(HostControlParam{
				Error: "錯誤: 設置聊天室快取資料發生問題"})
			conn.WriteMessage(b)
			return
		}

		// 設置過期時間
		// h.redisConn.SetExpire(redisName, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+activityID+"_"+channelID, "修改資料")
	} else if device == "pc" && gameID == "" {
		// 電腦端主持人大螢幕遙控頁面

		// 頻道更新為open狀態(mysql)
		// models.DefaultActivityChannelModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 	Update(true, activityID, channelID, "open")

		// 頻道更新為open狀態(mongo)
		models.DefaultActivityChannelModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateByMongo(true, activityID, channelID, "open")
	} else if device == "pc" && gameID != "" {
		// 電腦端遊戲遙控頁面

		// 開啟遊戲時先將遊戲設置redis資訊裡的遊戲狀態清空
		h.redisConn.HashSetCache(config.GAME_REDIS+gameID,
			"control_game_status", "")

		// 設置過期時間
		// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	} else if device == "mobile" && role == "guest" {
		// 手機玩家端隨時監控主持端所在頁面
	}

	if isOPen {
		// 開啟時傳送訊息
		var (
			newPage, newAction, newStatus, newGameID, newGameID2 string
		)

		// 取得redis中的page、action資料
		dataMap, _ := h.redisConn.HashGetAllCache(redisName)
		newPage, _ = dataMap["page"]
		newAction, _ = dataMap["action"]
		newStatus, _ = dataMap["game_status"]
		newGameID, _ = dataMap["game_id"]
		newGameID2, _ = dataMap["game_id_2"]

		if newPage != "" || newAction != "" || newStatus != "" {
			// page、action參數變動
			if newPage != page || newAction != action {
				b, _ := json.Marshal(HostControlParam{
					Game: GameModel{
						GameID:  newGameID,
						GameID2: newGameID2,
					},
					Page:   newPage,
					Action: newAction,
				})
				conn.WriteMessage(b)

				page = newPage
				action = newAction
			} else if newStatus != status {
				// game_status參數變動
				b, _ := json.Marshal(HostControlParam{
					Game: GameModel{
						GameID:     newGameID,
						GameID2:    newGameID2,
						GameStatus: newStatus,
					},
				})
				conn.WriteMessage(b)

				status = newStatus
			}
		}

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, config.CHANNEL_HOST_CONTROL_REDIS+activityID+"_"+channelID, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		var (
			newPage, newAction, newStatus, newGameID, newGameID2 string
		)

		// 取得redis中的page、action資料
		dataMap, _ := h.redisConn.HashGetAllCache(redisName)
		newPage, _ = dataMap["page"]
		newAction, _ = dataMap["action"]
		newStatus, _ = dataMap["game_status"]
		newGameID, _ = dataMap["game_id"]
		newGameID2, _ = dataMap["game_id_2"]

		if newPage != "" || newAction != "" || newStatus != "" {
			// page、action參數變動
			if newPage != page || newAction != action {
				b, _ := json.Marshal(HostControlParam{
					Game: GameModel{
						GameID:  newGameID,
						GameID2: newGameID2,
					},
					Page:   newPage,
					Action: newAction,
				})
				conn.WriteMessage(b)

				page = newPage
				action = newAction
			} else if newStatus != status {
				// game_status參數變動
				b, _ := json.Marshal(HostControlParam{
					Game: GameModel{
						GameID:     newGameID,
						GameID2:    newGameID2,
						GameStatus: newStatus,
					},
				})
				conn.WriteMessage(b)

				status = newStatus
			}
		}
	})

	for {
		var (
			result HostControlParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉手主持端遙控螢幕ws", err)

			// 電腦端大螢幕頁面裝置，關閉時清除可遙控的redis資訊
			if device == "pc" && gameID == "" {
				// 電腦端主持人大螢幕遙控頁面

				// 頻道更新為close狀態(mysql)
				// models.DefaultActivityChannelModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	Update(true, activityID, channelID, "close")

				// 頻道更新為close狀態(mongo)
				models.DefaultActivityChannelModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateByMongo(true, activityID, channelID, "close")

				// 清除主持端遙控資訊(host_control_活動ID_sessionID)
				h.redisConn.DelCache(redisName)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+activityID+"_"+channelID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				// h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
			}

			// 電腦端遊戲遙控頁面

			if device == "pc" && gameID != "" {
				// log.Println("電腦端遊戲頁面遙控裝置，關閉遊戲時將遊戲設置redis資訊裡的遊戲狀態清空")

				// 關閉遊戲時將遊戲設置redis資訊裡的遊戲狀態清空
				h.redisConn.HashSetCache(config.GAME_REDIS+gameID,
					"control_game_status", "")

				// 設置過期時間
				// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
			}

			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_" + channelID)
			conn.Close()
			return
		}

		// 收到訊息，修改redis資訊
		json.Unmarshal(data, &result)

		var (
			keys   = []string{"page", "action", "game_status", "game_id", "game_id_2"}
			values = []string{result.Page, result.Action, result.Game.GameStatus, result.Game.GameID, result.Game.GameID2}
			params = []interface{}{redisName}
		)
		for i, value := range values {
			if value != "" {
				params = append(params, keys[i], value)
			}
		}

		if len(params) > 1 {
			// 更新主持端遙控資訊(host_control_活動ID_sessionID)
			if err := h.redisConn.HashMultiSetCache(params); err != nil {
				b, _ := json.Marshal(HostControlParam{
					Error: "錯誤: 設置聊天室快取資料發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 設置過期時間
			// h.redisConn.SetExpire(redisName, config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+activityID+"_"+channelID, "修改資料")

			// 判斷game_id跟game_status是否為空，不為空的話更新遙控端遊戲狀態
			if result.Game.GameID != "" && result.Game.GameStatus != "" {
				if err := h.redisConn.HashSetCache(config.GAME_REDIS+gameID,
					"control_game_status", result.Game.GameStatus); err != nil {
					b, _ := json.Marshal(HostControlParam{
						Error: "錯誤: 設置遊戲快取資料發生問題"})
					conn.WriteMessage(b)
					return
				}

				// 設置過期時間
				// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
			}
		}

	}
}

// @Summary 該活動所有可控制的頻道資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Success 200 {array} []HostControlParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/host/control/channel [get]
func (h *Handler) HostControlSessionWebsocket(ctx *gin.Context) {
	var (
		activityID = ctx.Query("activity_id")
		// redisName            = config.HOST_CONTROL_CHANNEL_REDIS + activityID // 主持端所有可遙控的channel，hash
		wsConn, conn, err    = NewWebsocketConn(ctx)
		page, action, status string
		gameID, gameID2      string
		isOPen               = true
	)
	// log.Println("開啟取得所有可控制裝置資訊ws", activityID)
	defer wsConn.Close()
	defer conn.Close()

	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" {
		b, _ := json.Marshal(HostControlParam{
			Error: "錯誤: 無法辨識活動、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 開啟時傳送訊息
		var (
			channelModel, _ = models.DefaultActivityChannelModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindByMongo(true, activityID)
			controls = make([]HostControlParam, channelModel.ChannelAmount)
		)

		// 取得裝置的page.action.game_status...等資訊
		for i := range channelModel.ChannelAmount {
			channel := "channel_" + strconv.Itoa(int(i+1))
			channelStatus := ""

			if channel == "channel_1" {
				channelStatus = channelModel.Channel1
			} else if channel == "channel_2" {
				channelStatus = channelModel.Channel2
			} else if channel == "channel_3" {
				channelStatus = channelModel.Channel3
			} else if channel == "channel_4" {
				channelStatus = channelModel.Channel4
			} else if channel == "channel_5" {
				channelStatus = channelModel.Channel5
			} else if channel == "channel_6" {
				channelStatus = channelModel.Channel6
			} else if channel == "channel_7" {
				channelStatus = channelModel.Channel7
			} else if channel == "channel_8" {
				channelStatus = channelModel.Channel8
			} else if channel == "channel_9" {
				channelStatus = channelModel.Channel9
			} else if channel == "channel_10" {
				channelStatus = channelModel.Channel10
			}

			// 取得redis中的page、action資料，hash
			dataMap, _ := h.redisConn.HashGetAllCache(config.HOST_CONTROL_REDIS + activityID + "_" + channel) // host_control_id_channel_n
			// log.Println("頻道: ", config.HOST_CONTROL_REDIS+activityID+"_"+channel)
			// log.Println("取得資料: ", dataMap)
			// log.Println("dataMap[page]", dataMap["page"])

			page, _ = dataMap["page"]
			action, _ = dataMap["action"]
			status, _ = dataMap["game_status"]
			gameID, _ = dataMap["game_id"]
			gameID2, _ = dataMap["game_id_2"]

			controls[i] = HostControlParam{
				ChannelID:     channel,
				ChannelStatus: channelStatus,
				Game: GameModel{
					GameID:     gameID,
					GameID2:    gameID2,
					GameStatus: status,
				},
				Page:   page,
				Action: action,
			}
		}

		b, _ := json.Marshal(controls)
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribes(context,
		[]string{
			config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS + activityID, // 偵測所有頻道的開關狀態
			// 以下偵測所有頻道的頁面使用狀態
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_1",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_2",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_3",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_4",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_5",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_6",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_7",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_8",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_9",
			config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_10",
		},
		func(channel, message string) {
			// 資料改變時，回傳最新資料至前端
			var (
				channelModel, _ = models.DefaultActivityChannelModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						FindByMongo(true, activityID)
				controls = make([]HostControlParam, channelModel.ChannelAmount)
			)

			// 取得裝置的page.action.game_status...等資訊
			for i := range channelModel.ChannelAmount {
				channel := "channel_" + strconv.Itoa(int(i+1))
				channelStatus := ""

				if channel == "channel_1" {
					channelStatus = channelModel.Channel1
				} else if channel == "channel_2" {
					channelStatus = channelModel.Channel2
				} else if channel == "channel_3" {
					channelStatus = channelModel.Channel3
				} else if channel == "channel_4" {
					channelStatus = channelModel.Channel4
				} else if channel == "channel_5" {
					channelStatus = channelModel.Channel5
				} else if channel == "channel_6" {
					channelStatus = channelModel.Channel6
				} else if channel == "channel_7" {
					channelStatus = channelModel.Channel7
				} else if channel == "channel_8" {
					channelStatus = channelModel.Channel8
				} else if channel == "channel_9" {
					channelStatus = channelModel.Channel9
				} else if channel == "channel_10" {
					channelStatus = channelModel.Channel10
				}

				// 取得redis中的page、action資料，hash
				dataMap, _ := h.redisConn.HashGetAllCache(config.HOST_CONTROL_REDIS + activityID + "_" + channel) // host_control_id_channel_n
				page, _ = dataMap["page"]
				action, _ = dataMap["action"]
				status, _ = dataMap["game_status"]
				gameID, _ = dataMap["game_id"]
				gameID2, _ = dataMap["game_id_2"]

				controls[i] = HostControlParam{
					ChannelID:     channel,
					ChannelStatus: channelStatus,
					Game: GameModel{
						GameID:     gameID,
						GameID2:    gameID2,
						GameStatus: status,
					},
					Page:   page,
					Action: action,
				}
			}

			b, _ := json.Marshal(controls)
			conn.WriteMessage(b)
		})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉取得所有可控制裝置資訊ws")
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS + activityID)

			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_1")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_2")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_3")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_4")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_5")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_6")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_7")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_8")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_9")
			h.redisConn.Unsubscribe(config.CHANNEL_HOST_CONTROL_REDIS + activityID + "_channel_10")

			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			datas, _ = h.redisConn.SetGetMembers(redisName) // 該活動所有可控制裝置session資訊
// 			controls = make([]HostControlParam, 0)
// 		)

// 		// 取得裝置的page.action.game_status...等資訊
// 		for _, data := range datas {
// 			// 取得redis中的page、action資料
// 			dataMap, _ := h.redisConn.HashGetAllCache(config.HOST_CONTROL_REDIS + activityID + "_" + data)
// 			page, _ = dataMap["page"]
// 			action, _ = dataMap["action"]
// 			status, _ = dataMap["game_status"]

// 			controls = append(controls, HostControlParam{
// 				SessionID: data,
// 				Game: GameModel{
// 					GameStatus: status,
// 				},
// 				Page:   page,
// 				Action: action,
// 			})

// 		}

// 		b, _ := json.Marshal(controls)
// 		conn.WriteMessage(b)

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 5)
// 	}
// }()
// ###優化前，定頻###

// 回傳資料
// b, _ := json.Marshal(HostChatroomParam{
// })
// conn.WriteMessage(b)

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			newPage, newAction, newStatus, newGameID, newGameID2 string
// 		)

// 		// 取得redis中的page、action資料
// 		dataMap, _ := h.redisConn.HashGetAllCache(redisName)
// 		newPage, _ = dataMap["page"]
// 		newAction, _ = dataMap["action"]
// 		newStatus, _ = dataMap["game_status"]
// 		newGameID, _ = dataMap["game_id"]
// 		newGameID2, _ = dataMap["game_id_2"]

// 		if newPage != "" || newAction != "" || newStatus != "" {
// 			// page、action參數變動
// 			if newPage != page || newAction != action {
// 				b, _ := json.Marshal(HostControlParam{
// 					Game: GameModel{
// 						GameID:  newGameID,
// 						GameID2: newGameID2,
// 					},
// 					Page:   newPage,
// 					Action: newAction,
// 				})
// 				conn.WriteMessage(b)

// 				page = newPage
// 				action = newAction
// 			} else if newStatus != status {
// 				// game_status參數變動
// 				// fmt.Println("newStatus: ", newStatus)
// 				b, _ := json.Marshal(HostControlParam{
// 					Game: GameModel{
// 						GameID:     newGameID,
// 						GameID2:    newGameID2,
// 						GameStatus: newStatus,
// 					},
// 				})
// 				conn.WriteMessage(b)

// 				status = newStatus
// 			}
// 		}

// 		// 參數不為空，並且參數變動
// 		// if (newPage != "" || newAction != "" || newStatus != "") &&
// 		// 	(newPage != page || newAction != action || newStatus != status) {
// 		// 	// fmt.Println("newStatus: ", newStatus)
// 		// 	b, _ := json.Marshal(HostChatroomParam{
// 		// 		Page:   newPage,
// 		// 		Action: newAction,
// 		// 		Game: GameModel{
// 		// 			GameStatus: newStatus,
// 		// 		},
// 		// 	})
// 		// 	if err := conn.WriteMessage(b); err != nil {
// 		// 		return
// 		// 	}

// 		// 	page = newPage
// 		// 	action = newAction
// 		// 	status = newStatus
// 		// }

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// 測試
// fmt.Println("測試錯誤時會不會直接關閉ws")
// b, _ := json.Marshal(HostChatroomParam{
// 	Error: "收到關閉訊息"})
// conn.WriteMessage(b)
// return
// fmt.Println("
