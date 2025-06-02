package controller

import (
	"encoding/json"
	"hilive/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ChatroomRecordWriteParam 回傳聊天紀錄訊息
type ChatroomRecordWriteParam struct {
	ChatroomRecords     []models.ChatroomRecordModel // 聊天紀錄
	YesMessageAmount    int64                        `json:"yes_message_amount" example:"10"`    // 通過訊息數量
	NoMessageAmount     int64                        `json:"no_message_amount" example:"10"`     // 未通過訊息數量
	ReviewMessageAmount int64                        `json:"review_message_amount" example:"10"` // 審核訊息數量
	UserMessageAmount   int64                        `json:"user_message_amount" example:"10"`   // 用戶訊息數量
	YesPlayedAmount     int64                        `json:"yes_played_amount" example:"10"`     // 已播放數量
	NoPlayedAmount      int64                        `json:"no_played_amount" example:"10"`      // 未播放數量

	// 不同種類聊天訊息
	YesNormalMessageAmount     int64 `json:"yes_normal_message_amount" example:"10"`
	YesNormalBarrageAmount     int64 `json:"yes_normal_barrage_amount" example:"10"`
	YesSpecialBarrageAmount    int64 `json:"yes_special_barrage_amount" example:"10"`
	YesOccupyBarrageAmount     int64 `json:"yes_occupy_barrage_amount" example:"10"`
	NoNormalMessageAmount      int64 `json:"no_normal_message_amount" example:"10"`
	NoNormalBarrageAmount      int64 `json:"no_normal_barrage_amount" example:"10"`
	NoSpecialBarrageAmount     int64 `json:"no_special_barrage_amount" example:"10"`
	NoOccupyBarrageAmount      int64 `json:"no_occupy_barrage_amount" example:"10"`
	ReviewNormalMessageAmount  int64 `json:"review_normal_message_amount" example:"10"`
	ReviewNormalBarrageAmount  int64 `json:"review_normal_barrage_amount" example:"10"`
	ReviewSpecialBarrageAmount int64 `json:"review_special_barrage_amount" example:"10"`
	ReviewOccupyBarrageAmount  int64 `json:"review_occupy_barrage_amount" example:"10"`

	MessageSensitivityAmount int64 `json:"message_sensitivity_amount" example:"10"`

	Error        string `json:"error" example:"error message"` // 錯誤訊息
	IsSendPeople bool   `json:"is_send_people" example:"true"` // 是否傳遞人數

	MessageMethod string `json:"message_method" example:"add、update"` // 執行方式

	NoMessage []int64 `json:"no_message" example:"10"` // 未通過訊息陣列

	SensitivityWords []string `json:"sensitivity_words" example:"10"` // 敏感詞資料
	ReplaceWords     []string `json:"replace_words" example:"10"`     // 替代資料
}

// ChatroomRecordReadParam 接收聊天紀錄訊息
type ChatroomRecordReadParam struct {
	ID     string `json:"id" example:"id"`           // 聊天紀錄ID
	UserID string `json:"user_id" example:"user_id"` // 用戶ID
	// Role          string `json:"role" example:"role"`       // 腳色
	MessageType   string `json:"message_type" example:"normal-message"`
	MessageStyle  string `json:"message_style" example:"default"`
	MessagePrice  string `json:"message_price" example:"10"`             // 訊息價格
	Message       string `json:"message" example:"message"`              // 訊息
	MessageEffect string `json:"message_effect" example:"effect"`        // 效果
	MessageStatus string `json:"message_status" example:"yes、no、review"` // 狀態
	MessagePlayed string `json:"message_played" example:"yes、no"`        // 播放狀態
	MessageMethod string `json:"message_method" example:"add、update"`    // 執行方式
	Limit         int64  `json:"limit" example:"100"`                    // 需要的資料數量
	Offset        int64  `json:"offset" example:"100"`                   // 跳過的資料數輛
}

// @Summary 即時回傳聊天紀錄資料(包含通過.未通過.審核中)、新增聊天紀錄資料
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param user_id query string false "user ID"
// @Param role query string true "角色" Enums(host, guest)
// @param body body ChatroomRecordReadParam true "user_id、message... param"
// @Success 200 {array} ChatroomRecordWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/chatroom/record [get]
func (h *Handler) ChatroomRecordWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		userID            = ctx.Query("user_id")
		role              = ctx.Query("role")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// oldRecords        = make([]models.ChatroomRecordModel, 0) // 舊的資料
		// n int64
	)
	// log.Println("開啟即時回傳聊天紀錄資料ws", activityID)
	defer wsConn.Close()
	defer conn.Close()

	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || role == "" ||
		(role == "guest" && userID == "") {
		b, _ := json.Marshal(ChatroomRecordWriteParam{
			Error: "錯誤: 無法辨識活動、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	go func() {
		for {
			var (
				chatroomRecordModel      models.ChatroomRecordModel
				messageSensitivityAmount int64
				noMessage                = make([]int64, 0)
			)

			if role == "host" {
				// 查詢不同類型聊天訊息數量(主持端)
				chatroomRecordModel, _ = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindHostMessageAmount(activityID)

				// 敏感詞
				messageSensitivityAmount, _ = models.DefaultMessageSensitivityModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindAmount(activityID)

				// 狀態為no的資料(從資料表取得)
				records, _ := models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, "", "no", "", "", 0, 0)
				for _, record := range records {
					noMessage = append(noMessage, record.ID)
				}

			} else if role == "guest" && userID != "" {
				// 查詢不同類型聊天訊息數量(玩家端)
				chatroomRecordModel, _ = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindGuestMessageAmount(activityID, userID)

				// 敏感詞
				messageSensitivityAmount, _ = models.DefaultMessageSensitivityModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindAmount(activityID)
			}

			b, _ := json.Marshal(ChatroomRecordWriteParam{
				YesMessageAmount:    chatroomRecordModel.YesMessageAmount,
				NoMessageAmount:     chatroomRecordModel.NoMessageAmount,
				ReviewMessageAmount: chatroomRecordModel.ReviewMessageAmount,

				UserMessageAmount: chatroomRecordModel.UserMessageAmount,

				YesPlayedAmount: chatroomRecordModel.YesPlayedAmount,
				NoPlayedAmount:  chatroomRecordModel.NoPlayedAmount,

				// 不同種類訊息
				YesNormalMessageAmount:  chatroomRecordModel.YesNormalMessageAmount,
				YesNormalBarrageAmount:  chatroomRecordModel.YesNormalBarrageAmount,
				YesSpecialBarrageAmount: chatroomRecordModel.YesSpecialBarrageAmount,
				YesOccupyBarrageAmount:  chatroomRecordModel.YesOccupyBarrageAmount,

				NoNormalMessageAmount:  chatroomRecordModel.NoNormalMessageAmount,
				NoNormalBarrageAmount:  chatroomRecordModel.NoNormalBarrageAmount,
				NoSpecialBarrageAmount: chatroomRecordModel.NoSpecialBarrageAmount,
				NoOccupyBarrageAmount:  chatroomRecordModel.NoOccupyBarrageAmount,

				ReviewNormalMessageAmount:  chatroomRecordModel.ReviewNormalMessageAmount,
				ReviewNormalBarrageAmount:  chatroomRecordModel.ReviewNormalBarrageAmount,
				ReviewSpecialBarrageAmount: chatroomRecordModel.ReviewSpecialBarrageAmount,
				ReviewOccupyBarrageAmount:  chatroomRecordModel.ReviewOccupyBarrageAmount,
				IsSendPeople:               true,

				// 敏感詞
				MessageSensitivityAmount: messageSensitivityAmount,

				NoMessage: noMessage, // 審核狀態未通過陣列(id)
			})

			conn.WriteMessage(b)

			if conn.isClose { // ws關閉
				return
			}

			// #####壓力測試#####start
			// if activityID == "LevJ8qMfsJb9TlC0UXcj" {
			// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
			// 	wg.Add(1)             //計數器

			// 	for i := 1; i <= 1; i++ {
			// 		go func(i int) {
			// 			defer wg.Done()
			// 			// 新增聊天紀錄資料
			// 			models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
			// 				SetRedisConn(h.redisConn).Add(true, models.EditChatroomRecordModel{
			// 				UserID:        "U172765b72f9e72583cfcb25e9fd6605f",
			// 				ActivityID:    activityID,
			// 				MessageType:   "normal-message",
			// 				MessageStyle:  "default",
			// 				MessagePrice:  "0",
			// 				MessageStatus: "yes",
			// 				MessageEffect: "",
			// 				Message:       strconv.Itoa(int(n)),
			// 			})
			// 			n++
			// 			// 新增聊天紀錄資料
			// 			models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
			// 				SetRedisConn(h.redisConn).Add(true, models.EditChatroomRecordModel{
			// 				UserID:        "U172765b72f9e72583cfcb25e9fd6605f",
			// 				ActivityID:    activityID,
			// 				MessageType:   "normal-barrage",
			// 				MessageStyle:  "default",
			// 				MessagePrice:  "0",
			// 				MessageStatus: "yes",
			// 				MessageEffect: "",
			// 				Message:       strconv.Itoa(int(n)),
			// 			})
			// 			n++
			// 			// 新增聊天紀錄資料
			// 			models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
			// 				SetRedisConn(h.redisConn).Add(true, models.EditChatroomRecordModel{
			// 				UserID:        "U172765b72f9e72583cfcb25e9fd6605f",
			// 				ActivityID:    activityID,
			// 				MessageType:   "special-barrage",
			// 				MessageStyle:  "default",
			// 				MessagePrice:  "0",
			// 				MessageStatus: "yes",
			// 				MessageEffect: "",
			// 				Message:       strconv.Itoa(int(n)),
			// 			})
			// 			n++
			// 			// 新增聊天紀錄資料
			// 			models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
			// 				SetRedisConn(h.redisConn).Add(true, models.EditChatroomRecordModel{
			// 				UserID:        "U172765b72f9e72583cfcb25e9fd6605f",
			// 				ActivityID:    activityID,
			// 				MessageType:   "occupy-barrage",
			// 				MessageStyle:  "default",
			// 				MessagePrice:  "0",
			// 				MessageStatus: "yes",
			// 				MessageEffect: "",
			// 				Message:       strconv.Itoa(int(n)),
			// 			})
			// 			n++
			// 		}(i)
			// 	}

			// 	// for i := 1; i <= 5; i++ {
			// 	// 	go func(i int) {
			// 	// 		defer wg.Done()

			// 	// 	}(i)
			// 	// }

			// }

			// time.Sleep(time.Second * 5)
			// #####壓力測試#####end

			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var (
			result ChatroomRecordReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// #####壓力測試#####start
			// db.Table(config.ACTIVITY_CHATROOM_RECORD_TABLE).WithConn(h.dbConn).
			// 	Where("activity_id", "=", activityID).Delete()
			// #####壓力測試#####end

			// log.Println("關閉即時回傳聊天紀錄資料ws")
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offset != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得聊天訊息資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		// 用戶新增訊息或修改訊息
		if result.UserID != "" &&
			(result.MessageMethod == "add" || result.MessageMethod == "update") {
			if result.MessageMethod == "add" {
				// 新增聊天紀錄資料
				if err = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Add(true, models.EditChatroomRecordModel{
						UserID:        result.UserID,
						ActivityID:    activityID,
						MessageType:   result.MessageType,
						MessageStyle:  result.MessageStyle,
						MessagePrice:  result.MessagePrice,
						MessageStatus: result.MessageStatus,
						MessageEffect: result.MessageEffect,
						Message:       result.Message,
					}); err != nil {
					b, _ := json.Marshal(ChatroomRecordWriteParam{
						Error: "錯誤: 新增聊天紀錄資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			} else if result.MessageMethod == "update" {
				// 更新聊天紀錄播放資料
				if err = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Update(
						true, models.EditChatroomRecordModel{
							ID:            result.ID,
							ActivityID:    activityID,
							MessagePlayed: result.MessagePlayed,
						}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					b, _ := json.Marshal(ChatroomRecordWriteParam{
						Error: "錯誤: 更新聊天紀錄資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			}

			// 用戶聊天紀錄資料(從資料表取得，包含通過、未通過、審核中)
			// records, _ := models.DefaultChatroomRecordModel().
			// 	SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
			// 	Find(false, activityID, result.UserID, "", "", result.Limit, result.Offset)

			// b, _ := json.Marshal(ChatroomRecordWriteParam{
			// 	ChatroomRecords: records,
			// })
			// conn.WriteMessage(b)
		} else if result.MessageMethod == "find_1" || result.MessageMethod == "find_2" {
			// log.Println("收到find: ", result.Limit, result.Offset)
			// 取得聊天紀錄資料回傳前端
			var (
				records = make([]models.ChatroomRecordModel, 0) // 聊天資料
			)
			if role == "host" {
				// var isRedis bool
				// if result.Limit != 0 || result.Offset != 0 ||
				// 	result.MessageStatus != "" || result.MessagePlayed != ""||
				// 	result.UserID != "" {
				// 	// 取部分資料，用資料表
				// 	isRedis = false
				// } else {
				// 	// 取全部資料，用redis
				// 	isRedis = true
				// }

				// 所有聊天紀錄資料(從資料表取得，包含通過、未通過、審核中)
				records, _ = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, result.UserID, result.MessageStatus,
						result.MessagePlayed, result.MessageType,
						result.Limit, result.Offset)

				// log.Println("主持端回傳資料數: ", activityID, len(records))
			} else if role == "guest" && userID != "" {
				// 玩家
				// 用戶聊天紀錄資料(從資料表取得，包含通過、未通過、審核中)
				records, _ = models.DefaultChatroomRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, result.UserID, result.MessageStatus,
						result.MessagePlayed, result.MessageType,
						result.Limit, result.Offset)
			}

			// 回傳資料至前端
			b, _ := json.Marshal(ChatroomRecordWriteParam{
				ChatroomRecords: records,
				IsSendPeople:    false,
				MessageMethod:   result.MessageMethod,
			})
			conn.WriteMessage(b)
		} else if result.MessageMethod == "find_sensitivity" {
			// 敏感詞
			// 取得聊天紀錄資料回傳前端
			var (
				sensitivityWords = make([]string, 0) // 敏感詞資料
				replaceWords     = make([]string, 0) // 替代資料
			)

			// 敏感詞
			messages, _ := models.DefaultMessageSensitivityModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(activityID)

			for _, message := range messages {
				sensitivityWords = append(sensitivityWords, message.SensitivityWord)
				replaceWords = append(replaceWords, message.ReplaceWord)
			}

			// 回傳資料至前端
			b, _ := json.Marshal(ChatroomRecordWriteParam{
				SensitivityWords: sensitivityWords,
				ReplaceWords:     replaceWords,
				MessageMethod:    result.MessageMethod,
			})

			conn.WriteMessage(b)
		}
	}
}

// yesMessageAmount, noMessageAmount, reviewMessageAmount, userMessageAmount                                   int64
// yesPlayedAmount, noPlayedAmount                                                                             int64
// yesNormalMessageAmount, yesNormalBarrageAmount, yesSpecialBarrageAmount, yesOccupyBarrageAmount             int64
// noNormalMessageAmount, noNormalBarrageAmount, noSpecialBarrageAmount, noOccupyBarrageAmount                 int64
// reviewNormalMessageAmount, reviewNormalBarrageAmount, reviewSpecialBarrageAmount, reviewOccupyBarrageAmount int64

// // 通過資料數量
// yesMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "")

// // 未通過資料數量
// noMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "no", "", "")

// // 審核資料數量
// reviewMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "review", "", "")

// // 已播放數量
// yesPlayedAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "", "yes", "")

// // 未播放數量
// noPlayedAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "", "no", "")

// // 不同種類的聊天訊息
// yesNormalMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "normal-message")
// yesNormalBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "normal-barrage")
// yesSpecialBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "special-barrage")
// yesOccupyBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "occupy-barrage")

// noNormalMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "no", "", "normal-message")
// noNormalBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "no", "", "normal-barrage")
// noSpecialBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "no", "", "special-barrage")
// noOccupyBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "no", "", "occupy-barrage")

// reviewNormalMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "review", "", "normal-message")
// reviewNormalBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "review", "", "normal-barrage")
// reviewSpecialBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "review", "", "special-barrage")
// reviewOccupyBarrageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "review", "", "occupy-barrage")

// // 通過資料數量
// yesMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, userID, "yes", "", "")

// // 未通過資料數量
// noMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, userID, "no", "", "")

// // 審核資料數量
// reviewMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, userID, "review", "", "")

// // 用戶資料數量
// userMessageAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindMessageAmount(activityID, userID, "", "", "")

// ###優化前，定頻###
// var (
// 	records    = make([]models.ChatroomRecordModel, 0) // 未過濾的資料
// 	newRecords = make([]models.ChatroomRecordModel, 0) // 過濾後的資料
// 	isSend     bool                                    // 是否回傳訊息至前端
// )

// // 所有聊天紀錄資料(包含通過、未通過、審核中)，未過濾的資料
// records, _ = models.DefaultChatroomRecordModel().
// 	SetDbConn(h.dbConn).SetRedisConn(h.redisConn).Find(true, activityID)

// // 過濾黑名單資料
// for _, record := range records {
// 	if !h.IsBlackStaff(activityID, "", "message", record.UserID) { // 不是黑名單
// 		if role == "host" {
// 			newRecords = append(newRecords, record)
// 		} else if role == "guest" && record.UserID == userID { // 判斷是否為該用戶資料
// 			// 用戶端，回傳用戶的聊天紀錄(不論審核狀態)
// 			newRecords = append(newRecords, record)
// 		}
// 	}
// }

// if len(newRecords) != len(oldRecords) {
// 	// 資料長度出現變化，回傳新資料至前端
// 	isSend = true
// } else if len(newRecords) == len(oldRecords) {
// 	// 判斷原始資料是否有變動
// 	for i := 0; i < len(oldRecords); i++ {
// 		if oldRecords[i].MessageStatus != newRecords[i].MessageStatus || // 判斷聊天紀錄狀態是否改變
// 			oldRecords[i].MessagePlayed != newRecords[i].MessagePlayed { // 判斷播放紀錄狀態是否改變
// 			isSend = true
// 			break
// 		}
// 	}
// }

// if isSend {
// 	// 將新的資料放進old參數中
// 	oldRecords = newRecords

// 	b, _ := json.Marshal(ChatroomRecordWriteParam{
// 		ChatroomRecords: newRecords,
// 	})
// 	conn.WriteMessage(b)
// }
// ###優化前，定頻###

// HostChatroomParam 主持端聊天室參數
// type HostChatroomParam struct {
// 	Game   GameModel // 場次資訊
// 	Page   string    `json:"page" example:"massage"`        // 點擊頁面
// 	Action string    `json:"action" example:"next page"`    // 點擊動作
// 	Error  string    `json:"error" example:"error message"` // 錯誤訊息
// }

// if result.UserID == "" || result.MessageType == "" ||
// 	result.MessageStyle == "" || result.MessagePrice == "" ||
// 	result.Message == "" || result.MessageStatus == "" {
// 	// fmt.Println(result.MessagePrice == "")
// 	log.Println("Chatroom錯誤? ", result)
// 	// b, _ := json.Marshal(ChatroomRecordWriteParam{
// 	// 	Error: "錯誤: 取得聊天紀錄資料發生問題"})
// 	// conn.WriteMessage(b)
// 	// return
// }

// if role == "guest" && len(newRecords) > limit { // 玩家端聊天紀錄資料大於上限值
// 	newRecords = newRecords[len(newRecords)-limit:]
// 	// log.Println("玩家端聊天紀錄資料大於limit", len(newRecords))
// }

// 測試活動:dQSrpxSJlexj6G9EIph4
// if activityID == "dQSrpxSJlexj6G9EIph4" && role == "host" {
// 	h.redisConn.DelCache(config.CHATROOM_REDIS + activityID)
// }
// 測試活動:dQSrpxSJlexj6G9EIph4
// 測試活動:dQSrpxSJlexj6G9EIph4
// if activityID == "dQSrpxSJlexj6G9EIph4" && role == "host" {
// 	if n <= 1000 {
// 		for i := 0; i < 50; i++ {
// 			models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).Add(true, models.EditChatroomRecordModel{
// 				ActivityID:    activityID,
// 				UserID:        "U3140fd73cfd35dd992668ab3b6efdae9",
// 				MessageType:   "normal-barrage",
// 				MessageStyle:  "default",
// 				Message:       "test",
// 				MessageStatus: "yes",
// 				MessagePrice:  "100",
// 			})

// 			n++
// 		}
// 	}
// }
// 測試活動:dQSrpxSJlexj6G9EIph4

// // @@@Summary 手機控制主持人聊天室螢幕
// // @@@Tags Websocket
// // @@@version 1.0
// // @@@Accept  json
// // @@@Param activity_id query string true "activity ID"
// // @@@Param device query string true "device" Enums(pc, mobile)
// // @@@param body body HostChatroomParam true "page、action"
// // @@@Success 200 {array} HostChatroomParam
// // @@@Failure 500 {array} response.ResponseInternalServerError
// // @@@Router /ws/v1/host/chatroom [get]
// func (h *Handler) HostChatroomWebsocket(ctx *gin.Context) {
// 	var (
// 		activityID           = ctx.Query("activity_id")
// 		device               = ctx.Query("device")
// 		wsConn, conn, err    = NewWebsocketConn(ctx)
// 		result               HostChatroomParam
// 		page, action, status string
// 	)
// 	// log.Println("開啟手機控制主持人聊天室螢幕ws")
// 	defer wsConn.Close()
// 	defer conn.Close()
// 	// 判斷活動、遊戲資訊是否有效
// 	if err != nil || activityID == "" {
// 		b, _ := json.Marshal(HostChatroomParam{
// 			Error: "錯誤: 無法辨識活動、錯誤的網域請求"})
// 		conn.WriteMessage(b)
// 		return
// 	}

// 	// 手機端裝置，重置page、action(hash方式)
// 	if device == "mobile" {
// 		values := []interface{}{config.HOST_CHATROOM_REDIS + activityID}
// 		values = append(values, "page", "")
// 		values = append(values, "action", "")
// 		values = append(values, "game_status", "")
// 		if err := h.redisConn.HashMultiSetCache(values); err != nil {
// 			b, _ := json.Marshal(HostChatroomParam{
// 				Error: "錯誤: 設置聊天室快取資料發生問題"})
// 			conn.WriteMessage(b)
// 			return
// 		}

// 		// 設置redis過期時間
// 		h.redisConn.SetExpire(config.HOST_CHATROOM_REDIS+activityID,
// 			config.REDIS_EXPIRE)
// 	}

// 	go func() {
// 		for {
// 			var (
// 				newPage, newAction, newStatus string
// 			)

// 			// 取得redis中的page、action資料
// 			dataMap, _ := h.redisConn.HashGetAllCache(config.HOST_CHATROOM_REDIS + activityID)
// 			newPage, _ = dataMap["page"]
// 			newAction, _ = dataMap["action"]
// 			newStatus, _ = dataMap["game_status"]

// 			if newPage != "" || newAction != "" || newStatus != "" {
// 				// page、action參數變動
// 				if newPage != page || newAction != action {
// 					b, _ := json.Marshal(HostChatroomParam{
// 						Page:   newPage,
// 						Action: newAction,
// 					})
// 					conn.WriteMessage(b)

// 					page = newPage
// 					action = newAction
// 				} else if newStatus != status {
// 					// game_status參數變動
// 					// fmt.Println("newStatus: ", newStatus)
// 					b, _ := json.Marshal(HostChatroomParam{
// 						Game: GameModel{
// 							GameStatus: newStatus,
// 						},
// 					})
// 					conn.WriteMessage(b)

// 					status = newStatus
// 				}
// 			}

// 			// 參數不為空，並且參數變動
// 			// if (newPage != "" || newAction != "" || newStatus != "") &&
// 			// 	(newPage != page || newAction != action || newStatus != status) {
// 			// 	// fmt.Println("newStatus: ", newStatus)
// 			// 	b, _ := json.Marshal(HostChatroomParam{
// 			// 		Page:   newPage,
// 			// 		Action: newAction,
// 			// 		Game: GameModel{
// 			// 			GameStatus: newStatus,
// 			// 		},
// 			// 	})
// 			// 	if err := conn.WriteMessage(b); err != nil {
// 			// 		return
// 			// 	}

// 			// 	page = newPage
// 			// 	action = newAction
// 			// 	status = newStatus
// 			// }

// 			// ws關閉
// 			if conn.isClose {
// 				return
// 			}
// 			time.Sleep(time.Second * 1)
// 		}
// 	}()

// 	for {
// 		data, err := conn.ReadMessage()
// 		// log.Println("收到搖控訊息")
// 		if err != nil {
// 			// log.Println("關閉手機控制主持人聊天室螢幕ws", err)

// 			// 手機端裝置，關閉時清除redis資訊
// 			// if device == "mobile" {
// 			// 	h.redisConn.DelCache(config.HOST_CHATROOM_REDIS + activityID)
// 			// }

// 			conn.Close()
// 			return
// 		}

// 		// 收到訊息，修改redis資訊
// 		json.Unmarshal(data, &result)
// 		var (
// 			keys   = []string{"page", "action", "game_status"}
// 			values = []string{result.Page, result.Action, result.Game.GameStatus}
// 			params = []interface{}{config.HOST_CHATROOM_REDIS + activityID}
// 		)
// 		for i, value := range values {
// 			if value != "" {
// 				params = append(params, keys[i], value)
// 			}
// 		}

// 		if len(params) > 1 {
// 			if err := h.redisConn.HashMultiSetCache(params); err != nil {
// 				b, _ := json.Marshal(HostChatroomParam{
// 					Error: "錯誤: 設置聊天室快取資料發生問題"})
// 				conn.WriteMessage(b)
// 				return
// 			}
// 		}

// 		// 回傳資料
// 		// b, _ := json.Marshal(HostChatroomParam{
// 		// })
//
