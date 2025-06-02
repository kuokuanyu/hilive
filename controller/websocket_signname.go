package controller

import (
	"encoding/json"
	"hilive/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SignNameParam 簽名牆資訊參數
type SignNameParam struct {
	Signnames []models.SignnameModel // 簽名牆資料
	// ChatroomRecords []models.ChatroomRecordModel // 聊天紀錄資料
	SignnameAmount int64  `json:"signname_amount" example:"10"`  // 簽名牆數量
	IsAdd          bool   `json:"is_add" example:"true"`         // 是否儲存資料
	Error          string `json:"error" example:"error message"` // 錯誤訊息
	IsSendPeople   bool   `json:"is_send_people" example:"true"` // 是否傳遞人數
	IsDelete       bool   `json:"is_delete" example:"true"`      // 是否刪除
}

// SignNameReadParam 簽名牆接收訊息
type SignNameReadParam struct {
	ID      string `json:"id" example:"100"`          // id
	Role    string `json:"role" example:"guest"`      // 角色
	UserID  string `json:"user_id" example:"user_id"` // 用戶ID
	Picture string `json:"picture" example:"picture"` // 簽名資料
	Limit   int64  `json:"limit" example:"100"`       // 需要的資料數量
	Offset  int64  `json:"offset" example:"100"`      // 跳過的資料數量
	Method  string `json:"method" example:"sign"`     // 方法
}

// @Summary 即時新增簽名牆資訊(玩家端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param isredis query bool true "redis"
// @param body body SignNameReadParam true "picture"
// @Success 200 {array} SignNameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/guest/signname [get]
func (h *Handler) SignnameGuestWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		isredis           = ctx.Query("isredis")
		wsConn, conn, err = NewWebsocketConn(ctx)
	)

	defer wsConn.Close()
	defer conn.Close()

	// 判斷參數是否有效
	if err != nil || activityID == "" || isredis == "" {
		b, _ := json.Marshal(SignNameParam{
			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result SignNameReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// 刪除redis簽名牆相關資訊
			// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)               // 簽名牆資訊(HASH)
			// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資訊(LIST)

			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if (result.Role == "host" || result.Role == "guest") &&
			result.UserID != "" && result.Picture != "" {
			// 簽名牆設置資料(redis)
			signnameModel, err := h.getSignnameSetting(true, activityID)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			if result.Role == "guest" {
				// 角色為玩家時，判斷簽名牆是否為手機簽名模式
				// 終端機簽名模式
				if signnameModel.SignnameMode == "terminal" {
					b, _ := json.Marshal(SignNameParam{
						Error: "錯誤: 簽名牆為終端機簽名模式，用戶無法新增簽名"})
					conn.WriteMessage(b)
					return
				}
			}

			// 將資料寫入資料表.redis
			if err := models.DefaultSignnameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(true,
					models.NewSignnameModel{
						ActivityID:     activityID,
						ActivityUserID: signnameModel.UserID, // 活動主辦人user_id
						UserID:         result.UserID,
						Picture:        result.Picture, // base64格式
					}); err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳新的簽名牆資料
			b, _ := json.Marshal(SignNameParam{
				IsAdd: true,
			})
			conn.WriteMessage(b)
		}
	}
}

// @Summary 即時回傳所有簽名牆資訊(主持人端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param isredis query bool true "redis"
// @Param content query bool true "write、message"
// @param body body SignNameReadParam true "picture"
// @Success 200 {array} SignNameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/host/signname [get]
func (h *Handler) SignnameHostWebsocket(ctx *gin.Context) {
	var (
		activityID = ctx.Query("activity_id")
		isredis    = ctx.Query("isredis")
		// content           = ctx.Query("content")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// oldSignnames      = make([]models.SignnameModel, 0)
		// n = 1
	)

	defer wsConn.Close()
	defer conn.Close()

	// 判斷參數是否有效
	if err != nil || activityID == "" || isredis == "" {
		//  || content == "" {
		b, _ := json.Marshal(SignNameParam{
			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// #####壓力測試#####
	// db.Table(config.ACTIVITY_SIGNNAME_TABLE).WithConn(h.dbConn).
	// 	Where("activity_id", "=", activityID).Where("id", ">", 5000).Delete()

	// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)
	// h.redisConn.DelCache(config.SIGNNAME_DATAS_REDIS + activityID)
	// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID)
	// #####壓力測試#####

	go func() {
		for {

			// #####壓力測試#####
			// if activityID == "PDn2lrPfw0txzTCelctW" || activityID == "LevJ8qMfsJb9TlC0UXcj" {
			// 	// if n <= 500 {
			// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
			// 	wg.Add(100)           //計數器

			// 	activity, _ := models.DefaultActivityModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
			// 		Find(false, activityID)

			// 	for i := 0; i < 100; i++ {
			// 		go func(i int) {
			// 			defer wg.Done()
			// 			models.DefaultSignnameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
			// 				Add(true, models.NewSignnameModel{
			// 					ActivityID:     activityID,
			// 					ActivityUserID: activity.UserID,
			// 					UserID:         "U3140fd73cfd35dd992668ab3b6efdae9",
			// 					Picture:        "",
			// 				})
			// 		}(i)
			// 	}

			// 	wg.Wait() //等待計數器歸0

			// 	// n = n + 500
			// 	// }
			// }
			// #####壓力測試#####

			var (
				signnameAmount int64
			)

			// 簽名牆資料處理
			// 簽名數量
			signnameAmount, _ = models.DefaultSignnameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindSignnameAmount(activityID)

			// if content == "write" {
			// 	// 簽名牆資料處理
			// 	// 簽名數量
			// 	signnameAmount, _ = models.DefaultSignnameModel().SetDbConn(h.dbConn).
			// 		SetRedisConn(h.redisConn).FindSignnameAmount(activityID)
			// } else if content == "message" {
			// 	// 聊天紀錄資料處理
			// 	// 聊天紀錄通過的資料數量
			// 	signnameAmount, _ = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).
			// 		SetRedisConn(h.redisConn).FindMessageAmount(activityID, "", "yes", "", "")
			// }

			b, _ := json.Marshal(SignNameParam{
				Signnames: []models.SignnameModel{},
				// ChatroomRecords: []models.ChatroomRecordModel{},
				SignnameAmount: signnameAmount,
				IsSendPeople:   true,
			})
			conn.WriteMessage(b)

			// ws關閉
			if conn.isClose {
				// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)               // 簽名牆資訊(HASH)
				// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資訊(LIST)
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var (
			result SignNameReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// 刪除redis簽名牆相關資訊
			// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)               // 簽名牆資訊(HASH)
			// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資訊(LIST)

			// #####壓力測試#####
			// db.Table(config.ACTIVITY_SIGNNAME_TABLE).WithConn(h.dbConn).
			// 	Where("activity_id", "=", activityID).Where("id", ">", 5000).Delete()

			// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)
			// h.redisConn.DelCache(config.SIGNNAME_DATAS_REDIS + activityID)
			// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID)
			// #####壓力測試#####

			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offset != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得報名簽到資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		// 簽名牆資料處理
		if result.Method == "find" {
			// var isRedis bool
			// if result.Limit != 0 || result.Offset != 0 {
			// 	// 取部分資料，用資料表
			// 	isRedis = false
			// } else {
			// 	// 取全部資料，用redis
			// 	isRedis = true
			// }

			// 簽名資料
			signnames, _ := h.getSignnameDatas(false, activityID, result.Limit, result.Offset)

			// 簽名數量
			// signnameAmount, _ := models.DefaultSignnameModel().SetDbConn(h.dbConn).
			// 	SetRedisConn(h.redisConn).FindSignnameAmount(activityID)

			// 回傳簽到人員訊息至前端
			b, _ := json.Marshal(SignNameParam{
				Signnames: signnames,
				// SignnameAmount: signnameAmount,
				IsSendPeople: false,
			})
			conn.WriteMessage(b)
		} else if result.Method == "delete" && result.ID != "" {
			// ids := strings.Split(result.ID, ",") // 多個id

			// 刪除簽名資料
			err = models.DefaultSignnameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Delete(true, activityID, strings.Split(result.ID, ","))
			if err != nil {
				// 回傳簽到人員訊息至前端
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error(),
				})
				conn.WriteMessage(b)
			}

			// 回傳簽到人員訊息至前端
			b, _ := json.Marshal(SignNameParam{
				IsDelete: true,
			})
			conn.WriteMessage(b)
		}

		// if content == "write" {
		// 	// 簽名牆資料處理
		// 	if result.Method == "find" {
		// 		// var isRedis bool
		// 		// if result.Limit != 0 || result.Offset != 0 {
		// 		// 	// 取部分資料，用資料表
		// 		// 	isRedis = false
		// 		// } else {
		// 		// 	// 取全部資料，用redis
		// 		// 	isRedis = true
		// 		// }

		// 		// 簽名資料
		// 		signnames, _ := h.getSignnameDatas(false, activityID, result.Limit, result.Offset)

		// 		// 簽名數量
		// 		// signnameAmount, _ := models.DefaultSignnameModel().SetDbConn(h.dbConn).
		// 		// 	SetRedisConn(h.redisConn).FindSignnameAmount(activityID)

		// 		// 回傳簽到人員訊息至前端
		// 		b, _ := json.Marshal(SignNameParam{
		// 			Signnames: signnames,
		// 			// SignnameAmount: signnameAmount,
		// 			IsSendPeople: false,
		// 		})
		// 		conn.WriteMessage(b)
		// 	}
		// } else if content == "message" {
		// 	// 聊天紀錄資料處理
		// 	if result.Method == "find" {
		// 		// 所有聊天紀錄資料(從資料表取得，只取得通過資料)
		// 		records, _ := models.DefaultChatroomRecordModel().
		// 			SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
		// 			Find(false, activityID, result.UserID, "yes", "", "",
		// 				result.Limit, result.Offset)

		// 			// 回傳簽到人員訊息至前端
		// 		b, _ := json.Marshal(SignNameParam{
		// 			ChatroomRecords: records,
		// 			IsSendPeople:    false,
		// 		})
		// 		conn.WriteMessage(b)
		// 	}
		// }
	}
}

// getSignnameDatas 簽名牆資訊
func (h *Handler) getSignnameDatas(isRedis bool, activityID string, limit, offset int64) (signnames []models.SignnameModel, err error) {
	signnames, err = models.DefaultSignnameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(isRedis, activityID, limit, offset)
	if err != nil {
		return signnames, err
	}
	return
}

// getSignnameSetting 簽名牆設置資訊
func (h *Handler) getSignnameSetting(isRedis bool, activityID string) (signname models.SignnameSettingModel, err error) {
	signname, err = models.DefaultSignnameSettingModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(isRedis, activityID)
	if err != nil {
		return signname, err
	}
	return
}

// ###優化前，定頻###
// var (
// 	signnames = make([]models.SignnameModel, 0)
// )

// if isredis == "true" {
// 	// 簽名牆資料
// 	signnames, _ = h.getSignnameDatas(true, activityID)
// }

// if len(signnames) != len(oldSignnames) {
// 	b, _ := json.Marshal(SignNameParam{
// 		Signnames: signnames,
// 	})
// 	conn.WriteMessage(b)

// 	oldSignnames = signnames
// }

// // ws關閉
// if conn.isClose {
// 	// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)               // 簽名牆資訊(HASH)
// 	// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資訊(LIST)
// 	return
// }
// time.Sleep(time.Second * 5)
// ###優化前，定頻###

// if conn.isClose { // ws關閉
// 	return
// }

// 測試活動:Bfon6SaV6ORhmuDQUioI
// h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)
// h.redisConn.DelCache(config.SIGNNAME_DATAS_REDIS + activityID)
// h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID)
// 測試活動:Bfon6SaV6ORhmuDQUioI

// 測試活動:Bfon6SaV6ORhmuDQUioI
// if activityID == "Bfon6SaV6ORhmuDQUioI" {
// 	for i := 0; i < 5; i++ {
// 		models.DefaultSignnameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 			Add(true, models.NewSignnameModel{
// 				ActivityID:     activityID,
// 				ActivityUserID: "admin",
// 				UserID:         "U3140fd73cfd35dd992668ab3b6efdae9",
// 				Picture:        "",
// 			})
// 	}
// }
// 測試活動:Bfon6SaV6ORhmuDQUioI

// if (result.Role != "host" && result.Role != "guest") ||
// 	result.UserID == "" || result.Picture == "" {
// 	b, _ := json.Marshal(SignNameParam{
// 		Error: "錯誤: 取得簽名牆資料發生問題"})
// 	conn.WriteMessage(b)
// 	return
// }

// go func() {
// 	for {
// 		// if conn.isClose { // ws關閉
// 		// 	return
// 		// }

// 		if isredis == "true" {
// 			// 簽名牆資料
// 			signnames, _ = h.getSignnames(true, activityID)
// 		}

// 		// fmt.Println("資料數: ", len(signnames))

// 		b, _ := json.Marshal(SignNameParam{
// 			Signnames: signnames,
// 		})
// 		conn.WriteMessage(b)

// 		// ws關閉
// 		if conn.isClose {
// 			h.redisConn.DelCache(config.SIGNNAME_REDIS + activityID)               // 簽名牆資訊(HASH)
// 			h.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資訊(LIST)
// 			return
// 		}
// 		time.Sleep(time.Second * 5)
// 	}
// }()

//
