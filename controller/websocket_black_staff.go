package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"

	"github.com/gin-gonic/gin"
)

// BlackStaffParam 黑名單人員參數資訊
type BlackStaffParam struct {
	IsBlack bool     `json:"is_black" example:"true"` // 是否黑名單
	Staffs  []string // 黑名單人員
	Error   string   `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 即時回傳黑名單資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string false "game ID"
// @Param user_id query string false "user ID"
// @Param game query string true "game" Enums(activity, message, question)
// @Success 200 {array} BlackStaffParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/black/staff [get]
func (h *Handler) BlackStaffWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		userID            = ctx.Query("user_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// oldBlackStaffs    = make([]models.BlackStaffModel, 0)
		channel string
		isOPen  = true
	)
	// log.Println("開啟即時回傳黑名單資訊ws", activityID, userID)
	defer wsConn.Close()
	defer conn.Close()
	if err != nil || activityID == "" || game == "" {
		b, _ := json.Marshal(BlackStaffParam{
			Error: "錯誤: 無法辨識資訊、錯誤的網域請求"})
		// fmt.Println("錯誤: ", err != nil, err)
		conn.WriteMessage(b)
		return
	}

	// redis key名稱
	if game == "activity" {
		// 活動
		channel = config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + activityID
	} else if game == "message" {
		// 訊息
		channel = config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + activityID
	} else if game == "question" {
		// 提問
		channel = config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + activityID
	} else if game == "signname" {
		// 簽名
		channel = config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + activityID
	} else if gameID != "" {
		// 遊戲
		channel = config.CHANNEL_BLACK_STAFFS_GAME_REDIS + gameID
	}

	if isOPen {
		// 開啟時傳送訊息
		if userID != "" { // 判斷用戶是否為黑名單
			// 判斷是否為黑名單(redis處理)
			isBlack := h.IsBlackStaff(activityID, gameID, game, userID)

			b, _ := json.Marshal(BlackStaffParam{IsBlack: isBlack})
			conn.WriteMessage(b)
		} else {
			// 取得所有黑名單人員資料
			var staffs = make([]string, 0)
			blackStaffs, _ := models.DefaultBlackStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindAll(true, activityID, gameID, game)

			for _, staffModel := range blackStaffs {
				staffs = append(staffs, staffModel.UserID)
			}

			b, _ := json.Marshal(BlackStaffParam{Staffs: staffs})
			conn.WriteMessage(b)
		}

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, channel, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		if userID != "" { // 判斷用戶是否為黑名單
			// 判斷是否為黑名單(redis處理)
			isBlack := h.IsBlackStaff(activityID, gameID, game, userID)

			b, _ := json.Marshal(BlackStaffParam{IsBlack: isBlack})
			conn.WriteMessage(b)
		} else {
			// 取得所有黑名單人員資料
			var staffs = make([]string, 0)
			blackStaffs, _ := models.DefaultBlackStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindAll(true, activityID, gameID, game)

			for _, staffModel := range blackStaffs {
				staffs = append(staffs, staffModel.UserID)
			}

			b, _ := json.Marshal(BlackStaffParam{Staffs: staffs})
			conn.WriteMessage(b)
		}
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉即時回傳黑名單資訊ws")
			// 刪除redis黑名單相關資訊
			// h.redisConn.DelCache(config.ACTIVITY_REDIS + activityID) // 活動資訊

			// 取消訂閱
			h.redisConn.Unsubscribe(channel)

			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		if userID != "" { // 判斷用戶是否為黑名單
// 			// 判斷是否為黑名單(redis處理)
// 			isBlack := h.IsBlackStaff(activityID, gameID, game, userID)

// 			b, _ := json.Marshal(BlackStaffParam{IsBlack: isBlack})
// 			conn.WriteMessage(b)
// 		} else { // 取得所有黑名單人員資料
// 			var staffs = make([]string, 0)
// 			newBlackStaffs, _ := models.DefaultBlackStaffModel().
// 				SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 				FindAll(true, activityID, gameID, game)

// 			if len(newBlackStaffs) != len(oldBlackStaffs) {
// 				// 將新的資料放進old參數中
// 				oldBlackStaffs = newBlackStaffs

// 				for _, staffModel := range newBlackStaffs {
// 					staffs = append(staffs, staffModel.UserID)
// 				}

// 				b, _ := json.Marshal(BlackStaffParam{Staffs: staffs})
// 				conn.WriteMessage(b)
// 			}
// 		}

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}

// 		if userID != "" { // 判斷用戶是否為黑名單
// 			time.Sleep(time.Second * 30)
// 		} else { // 取得所有黑名單人員資料
// 			time.Sleep(time.Second * 5)
// 		}
// 	}
// }()
// ###優化前，定頻###

// activity.ActivityName = activityModel.ActivityName
// activity.ActivityType = activityModel.ActivityType
// activity.People = activityModel.People
// activity.Attend = activityModel.Attend
// activity.City = activityModel.City
// activity.Town = activityModel.Town
// activity.StartTime = activityModel.StartTime
// activity.EndTime = activityModel.
