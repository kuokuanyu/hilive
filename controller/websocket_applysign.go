package controller

import (
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"time"

	"github.com/gin-gonic/gin"
)

// ApplysignParam 報名簽到簽到人員資訊參數
type ApplysignParam struct {
	IsSign          bool                    `json:"is_sign" example:"true"` // 是否簽到
	ApplysignStaffs []models.ApplysignModel // 報名簽到簽到人員
	Status          string                  `json:"status" example:"sign"`         // 簽到狀態
	Name            string                  `json:"name" example:"name"`           // 名稱
	Method          string                  `json:"method" example:"find"`         // 方法
	AllPeople       int64                   `json:"all_people" example:"100"`      // 所有人數
	NoPeople        int64                   `json:"no_people" example:"100"`       // 未報名人數
	ReviewPeople    int64                   `json:"review_people" example:"100"`   // 審核人數
	ApplyPeople     int64                   `json:"apply_people" example:"100"`    // 報名人數
	RefusePeople    int64                   `json:"refuse_people" example:"100"`   // 拒絕人數
	SignPeople      int64                   `json:"sign_people" example:"100"`     // 簽到人數
	People          int64                   `json:"people" example:"100"`          // 人數
	Limit           int64                   `json:"limit" example:"100"`           // 需要的資料數量
	Offset          int64                   `json:"offset" example:"100"`          // 跳過的資料數量
	Error           string                  `json:"error" example:"error message"` // 錯誤訊息
	IsSendPeople    bool                    `json:"is_send_people" example:"true"` // 是否傳遞人數
}

// @Summary 即時回傳報名簽到人員資訊(平台管理員端報名簽到頁面)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param user_id query string true "user ID"
// @Param isredis query bool true "redis"
// @param body body ApplysignParam true "applysign param"
// @Success 200 {array} ApplysignParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/applysign [get]
func (h *Handler) ApplysignWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		userID            = ctx.Query("user_id")
		isredis           = ctx.Query("isredis")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// oldStaffs         = make([]models.ApplysignModel, 0)
		// isOPen = true
	)
	// log.Println("開啟即時回傳報名簽到人員資訊(平台管理員端報名簽到頁面)", activityID, userID)
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || isredis == "" {
		b, _ := json.Marshal(SignStaffParam{
			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// ###優化前，定頻###
	go func() {
		for {
			if isredis == "true" {
				// 使用SET格式處理
				// 判斷用戶是否簽到完成
				isSign := h.IsSign(config.SIGN_STAFFS_2_REDIS, activityID, userID)

				b, _ := json.Marshal(ApplysignParam{
					IsSign: isSign,
				})
				conn.WriteMessage(b)
			} else if isredis == "false" {
				// 查詢該活動不同狀態的人員數量(主持端)
				applysignModel, _ := models.DefaultApplysignModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindHostApplysignAmount(activityID)

				b, _ := json.Marshal(ApplysignParam{
					ApplyPeople:  applysignModel.ApplyPeople,
					SignPeople:   applysignModel.SignPeople,
					AllPeople:    applysignModel.AllPeople,
					ReviewPeople: applysignModel.ReviewPeople,
					NoPeople:     applysignModel.NoPeople,
					RefusePeople: applysignModel.RefusePeople,
					IsSendPeople: true,
				})
				conn.WriteMessage(b)
			}

			// ws關閉
			if conn.isClose {
				return
			}

			time.Sleep(time.Second * 5)
		}
	}()
	// ###優化前，定頻###

	for {
		var (
			result ApplysignParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉即時回傳報名簽到人員資訊(平台管理員端報名簽到頁面)")
			// 取消訂閱
			// h.redisConn.Unsubscribe(config.CHANNEL_SIGN_STAFFS_2_REDIS + activityID)
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offset != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得報名簽到資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		if isredis == "false" && result.Method == "find" {
			// 使用資料庫
			// 報名簽到人員資料(包含簽到、未簽到、審核中...人員)
			staffs, _ := h.getAllStaffs(activityID, "", result.Name, result.Status,
				result.Limit, result.Offset)

			amount, _ := models.DefaultApplysignModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindApplysignAmount(activityID, result.Name, result.Status)

			// 回傳簽到人員訊息至前端
			b, _ := json.Marshal(ApplysignParam{
				ApplysignStaffs: staffs,
				IsSendPeople:    false,
				People:          amount,
			})
			conn.WriteMessage(b)
		}
	}
}

// // 所有報名簽到人員數量(包含簽到、未簽到、審核中...人員)
// allAnount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "")

// // 未報名人員數量
// noAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "no")

// // 審核人員數量
// reviewAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "review")

// // 報名完成人員數量
// applyAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "apply")

// // 拒絕簽到人員數量
// refuseAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "refuse")

// // 簽到完成人員數量
// signAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "", "sign")

// ###pub/sub###
// if isOPen {
// 	// 開啟時傳送訊息
// 	if isredis == "true" {
// 		// 優化後，使用SET格式處理
// 		isSign := h.IsSign(config.SIGN_STAFFS_2_REDIS, activityID, userID)

// 		b, _ := json.Marshal(ApplysignParam{
// 			IsSign: isSign,
// 		})
// 		conn.WriteMessage(b)
// 	} else if isredis == "false" {
// 		// 所有報名簽到人員數量(包含簽到、未簽到、審核中...人員)
// 		allAnount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "")

// 		// 未報名人員數量
// 		noAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "no")

// 		// 審核人員數量
// 		reviewAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "review")

// 		// 報名完成人員數量
// 		applyAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "apply")

// 		// 拒絕簽到人員數量
// 		refuseAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "refuse")

// 		// 簽到完成人員數量
// 		signAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "sign")

// 		b, _ := json.Marshal(ApplysignParam{
// 			ApplyPeople:  applyAmount,
// 			SignPeople:   signAmount,
// 			AllPeople:    allAnount,
// 			ReviewPeople: reviewAmount,
// 			NoPeople:     noAmount,
// 			RefusePeople: refuseAmount,
// 		})
// 		conn.WriteMessage(b)
// 	}

// 	isOPen = false
// }
// ###pub/sub###

// 用來控制 Redis 訂閱開關判斷
// context, cancel := context.WithCancel(ctx.Request.Context())
// defer cancel() // 確保函式退出時取消訂閱

// 啟用redis訂閱
// go h.redisConn.Subscribe(context, config.CHANNEL_SIGN_STAFFS_2_REDIS+activityID, func(channel, message string) {
// 	// 資料改變時，回傳最新資料至前端
// 	// log.Println("資料改變? ", message)

// 	// 開啟時傳送訊息
// 	if isredis == "true" {
// 		// 優化後，使用SET格式處理
// 		isSign := h.IsSign(config.SIGN_STAFFS_2_REDIS, activityID, userID)

// 		b, _ := json.Marshal(ApplysignParam{
// 			IsSign: isSign,
// 		})
// 		conn.WriteMessage(b)
// 	} else if isredis == "false" {
// 		// 所有報名簽到人員數量(包含簽到、未簽到、審核中...人員)
// 		allAnount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "")

// 		// 未報名人員數量
// 		noAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "no")

// 		// 審核人員數量
// 		reviewAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "review")

// 		// 報名完成人員數量
// 		applyAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "apply")

// 		// 拒絕簽到人員數量
// 		refuseAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "refuse")

// 		// 簽到完成人員數量
// 		signAmount, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).FindApplysignAmount(activityID, "sign")

// 		b, _ := json.Marshal(ApplysignParam{
// 			ApplyPeople:  applyAmount,
// 			SignPeople:   signAmount,
// 			AllPeople:    allAnount,
// 			ReviewPeople: reviewAmount,
// 			NoPeople:     noAmount,
// 			RefusePeople: refuseAmount,
// 		})
// 		conn.WriteMessage(b)
// 	}
// })

// getAllStaffs 所有人員資訊(包含簽到、未簽到、審核中...人員)
func (h *Handler) getAllStaffs(activityID string, userID string, name string, status string, limit, offset int64) (signStaffs []models.ApplysignModel, err error) {
	signStaffs, err = models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAll(activityID, userID, name, status, limit, offset)
	if err != nil {
		return signStaffs, err
	}
	return
}

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
// if err := conn.WriteMess
