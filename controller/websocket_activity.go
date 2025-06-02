package controller

// ActivityParam 活動相關參數資訊
// type ActivityParam struct {
// 	Activity models.ActivityModel // 活動資訊
// 	Error    string               `json:"error" example:"error message"` // 錯誤訊息
// }

// ActivityModel 資料表欄位
// type ActivityModel struct {
// 	ActivityName string `json:"activity_name" example:"Activity Name"`
// 	ActivityType string `json:"activity_type" example:"Activity Type"`
// 	People       int64  `json:"people" example:"100"`
// 	Attend       int64  `json:"attend" example:"100"`
// 	City         string `json:"city" example:"台中市"`
// 	Town         string `json:"town" example:"南區"`
// 	StartTime    string `json:"start_time" example:"2022/01/01 00:00:00"`
// 	EndTime      string `json:"end_time" example:"2022/01/01 00:00:00"`
// }

// @Summary 即時回傳活動資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Success 200 {array} ActivityParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/activity [get]
// func (h *Handler) ActivityWebsocket(ctx *gin.Context) {
// 	var (
// 		activityID        = ctx.Query("activity_id")
// 		wsConn, conn, err = NewWebsocketConn(ctx)
// 	)
// 	fmt.Println("開啟即時活動資訊ws")
// 	defer wsConn.Close()
// 	defer conn.Close()
// 	if err != nil || activityID == "" {
// 		b, _ := json.Marshal(ActivityParam{
// 			Error: "錯誤: 無法辨識活動資訊、錯誤的網域請求"})
// 		conn.WriteMessage(b)
// 		return
// 	}

// 	go func() {
// 		for {
// 			activityModel, _ := models.DefaultActivityModel().
// 				SetDbConn(h.dbConn).SetRedisConn(h.redisConn).Find(true, activityID)
// 			b, _ := json.Marshal(ActivityParam{Activity: activityModel})
// 			if err := conn.WriteMessage(b); err != nil {
// 				return
// 			}

// 			// ws關閉
// 			if conn.isClose {
// 				// 刪除redis活動相關資訊
// 				// h.redisConn.DelCache(config.ACTIVITY_REDIS + activityID) // 活動資訊
// 				return
// 			}
// 			time.Sleep(time.Second * 5)
// 		}
// 	}()

// 	for {
// 		_, err := conn.ReadMessage()
// 		if err != nil {
// 			// 刪除redis活動相關資訊
// 			// h.redisConn.DelCache(config.ACTIVITY_REDIS + activityID) // 活動資訊

// 			fmt.Println("關閉即時活動資訊ws")
// 			conn.Close()
// 			return
// 		}
// 	}
// }

// activity.ActivityName = activityModel.ActivityName
// activity.ActivityType = activityModel.ActivityType
// activity.People = activityModel.People
// activity.Attend = activityModel.Attend
// activity.City = activityModel.City
// activity.Town = activityModel.Town
// activity.StartTime = activityModel.StartTime
// activity.EndTime = activityModel.EndTime
