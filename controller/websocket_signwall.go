package controller

import (
	"encoding/json"
	"hilive/models"
	"time"

	"github.com/gin-gonic/gin"
)

// SignStaffModel 資料表欄位
// type SignStaffModel struct {
// 	ID         int64  `json:"id" example:"1"`
// 	UserID     string `json:"user_id" example:"user_id"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	Name       string `json:"name" example:"username"`
// 	Avatar     string `json:"avatar" example:"https://..."`
// 	Number     int64  `json:"number" example:"5"`
// 	// Status     string `json:"status" example:"sign"`
// 	// ApplyTime  string `json:"apply_time"`
// 	// ReviewTime string `json:"review_time"`
// 	// SignTime   string `json:"sign_time"`
// }

// SignStaffParam 簽到人員資訊參數
type SignStaffParam struct {
	SignStaffs   []models.ApplysignModel // 簽到人員
	SignPeople   int64                   `json:"sign_people" example:"10"`      // 簽到人數
	Limit        int64                   `json:"limit" example:"100"`           // 需要的資料數量
	Offset       int64                   `json:"offset" example:"100"`          // 跳過的資料數量
	Error        string                  `json:"error" example:"error message"` // 錯誤訊息
	Method       string                  `json:"method" example:"sign"`         // 方法
	IsSendPeople bool                    `json:"is_send_people" example:"true"` // 是否傳遞人數
}

// @Summary 即時活動簽到人數、人員資訊(主持人端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @param body body SignStaffParam true "applysign param"
// @Success 200 {array} SignStaffParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/host/signwall [get]
func (h *Handler) SignWallWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// oldStaffs         = make([]models.ApplysignModel, 0)
		// n = 1
	)
	// fmt.Println("開啟活動簽到即時人數、人員資訊ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" {
		b, _ := json.Marshal(SignStaffParam{
			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// #####壓力測試#####
	// delete data
	// if activityID == "PDn2lrPfw0txzTCelctW" {
	// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(h.dbConn).
	// 		Where("activity_id", "=", activityID).
	// 		Where("id", ">", 200000).Delete()

	// 	h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊
	// }

	// #####壓力測試#####

	go func() {
		for {
			// #####壓力測試#####
			// if activityID == "PDn2lrPfw0txzTCelctW" {
			// 	if n <= 10000 {
			// 		var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
			// 		wg.Add(100)           //計數器

			// 		for i := n; i < n+100; i++ {
			// 			go func(i int) {
			// 				defer wg.Done()
			// 				// 遞增活動人數
			// 				// models.DefaultActivityModel().SetDbConn(h.dbConn).
			// 				// 	SetRedisConn(h.redisConn).IncrAttend(true, activityID,
			// 				// 	strconv.Itoa(int(i)))

			// 				// 測試定時加入簽到人員
			// 				models.DefaultApplysignModel().
			// 					SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
			// 					Add(true, models.NewApplysignModel{
			// 						UserID:     strconv.Itoa(int(i)),
			// 						ActivityID: activityID,
			// 						Status:     "sign",
			// 					})
			// 			}(i)
			// 		}

			// 		wg.Wait() //等待計數器歸0

			// 		n = n + 100
			// 	}
			// }

			// #####壓力測試#####

			// 簽到完成人員數量
			signAmount, _ := models.DefaultApplysignModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindApplysignAmount(activityID, "", "sign")

			b, _ := json.Marshal(SignStaffParam{
				SignStaffs:   []models.ApplysignModel{},
				SignPeople:   signAmount,
				IsSendPeople: true,
			})
			conn.WriteMessage(b)

			// ws關閉
			if conn.isClose {
				// h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var (
			result SignStaffParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// #####壓力測試#####
			// 刪除資料
			// delete data
			// if activityID == "PDn2lrPfw0txzTCelctW" {
			// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(h.dbConn).
			// 		Where("activity_id", "=", activityID).
			// 		Where("id", ">", 200000).Delete()

			// 	h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊
			// }

			// // #####壓力測試#####

			// // 刪除redis簽到相關資訊
			// h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊

			// fmt.Println("關閉活動簽到即時人數、人員資訊ws")
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offset != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得報名簽到資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		if result.Method == "find" {
			// var isRedis bool
			// if result.Limit != 0 || result.Offset != 0 {
			// 	// 取部分資料，用資料表
			// 	isRedis = false
			// } else {
			// 	// 取全部資料，用redis
			// 	isRedis = true
			// }

			// 聊天紀錄資料(包含通過、未通過、審核中)
			staffs, _ := h.getSignStaffs(false, false, "", activityID, result.Limit, result.Offset)

			// 簽到完成人員數量
			// signAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
			// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "sign")

			// 回傳簽到人員訊息至前端
			b, _ := json.Marshal(SignStaffParam{
				SignStaffs: staffs,
				// SignPeople: signAmount,
				IsSendPeople: false,
			})
			conn.WriteMessage(b)
		}
	}
}

// ###優化前，定頻###
// 簽到人員資料(需要用戶詳細資訊)
// newStaffs, _ := h.getSignStaffs(true, true, config.SIGN_STAFFS_1_REDIS, activityID)

// // 判斷簽到人數是否改變，改變才回傳至前端
// if len(newStaffs) != len(oldStaffs) {
// 	b, _ := json.Marshal(SignStaffParam{
// 		SignStaffs: newStaffs,
// 		SignPeople: int64(len(newStaffs)),
// 	})
// 	conn.WriteMessage(b)

// 	oldStaffs = newStaffs
// }
// ###優化前，定頻###

// 原版
// 簽到人員資料(需要用戶詳細資訊)
// staffs, _ = h.getSignStaffs(true, true, config.SIGN_STAFFS_1_REDIS, activityID)
// // fmt.Println("人數: ", int64(len(staffs)))

// b, _ := json.Marshal(SignStaffParam{
// 	SignStaffs: staffs,
// 	SignPeople: int64(len(staffs)),
// })
// conn.WriteMessage(b)

// if conn.isClose { // ws關閉
// 	h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊
// 	return
// }
// time.Sleep(time.Second * 5)
// 原版

// 測試活動:dQSrpxSJlexj6G9EIph4
// if activityID == "dQSrpxSJlexj6G9EIph4" {
// 	h.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID)
// }
// 測試活動:dQSrpxSJlexj6G9EIph4

// if conn.isClose { // ws關閉
// 	return
// }

// 測試活動:dQSrpxSJlexj6G9EIph4
// if activityID == "dQSrpxSJlexj6G9EIph4" {
// 	if n <= 1000 {
// 		for i := n; i < n+10; i++ {
// 			// 測試定時加入簽到人員
// 			models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 				Add(models.NewApplysignModel{
// 					UserID:     strconv.Itoa(i),
// 					ActivityID: activityID,
// 					Status:     "sign",
// 				})

// 				// 遞增活動人數
// 			models.DefaultActivityModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).IncrAttend(true, "dQSrpxSJlexj6G9EIph4",
// 				strconv.Itoa(i))

// 		}

// 		n = n + 10
// 	}
// }
// 測試活動:dQSrpxSJlexj6G9EIph4

// 測試資料
// if n == 0 {
// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	wg.Add(36)            //計數器
// 	for i := 1; i <= 36; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			h.redisConn.ListRPush(config.SIGN_STAFFS_1_REDIS+activityID,
// 				strconv.Itoa(i))
// 			// n++
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0
// 	n++
// }

// if newStaffs, _ := h.getSignStaffs(activityID); len(newStaffs) !=
// 	int(staffsLength) { // 新增簽到人員，將新的推送至陣列
// 	for i := int(staffsLength); i < len(newStaffs); i++ {
// 		staffs = append(staffs, SignStaffModel{
// 			ID:         newStaffs[i].ID,
// 			UserID:     newStaffs[i].UserID,
// 			ActivityID: newStaffs[i].ActivityID,
// 			Name:       newStaffs[i].Name,
// 			Avatar:     newStaffs[i].Avatar,
// 			Number:     newStaffs[i].Number,
// 			// Status:     newStaffs[i].Status,
// 			// ApplyTime:  newStaffs[i].ApplyTime,
// 			// ReviewTime: newStaffs[i].ReviewTime,
// 			// SignTime:   newStaffs[i].SignTime,
// 		})
// 	}
// staffsLength = len(newStaffs)
// fmt.Println("簽到人員數量: ", staffsLength)
// fmt.Println("簽到人員資訊: ", staffs)

// 簽到人員資訊
// if staffs, err = h.getSignStaffs(activityID); err != nil {
// 	b, _ := json.Marshal(SignStaffParam{Error: err.Error()})
// 	conn.WriteMessage(b)
// 	return
// }
// staffsLength = len(staffs)

// b, _ := json.Marshal(SignStaffParam{
// 	SignStaffs: staffs,
// 	SignPeople: int64(staffsLength),
// })
// fmt.Println("簽到人員數量: ", staffsLength)
// fmt.Println("簽到人員資訊: ", staffs)
// if err := conn.WriteMessage(b); err != nil {
// 	return
