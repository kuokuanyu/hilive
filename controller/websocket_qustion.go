package controller

import (
	"encoding/json"
	"hilive/models"
	"time"

	"github.com/gin-gonic/gin"
)

// HostQuestionReadParam 主持端提問牆接收訊息
type HostQuestionReadParam struct {
	// SessionID string `json:"session_id" example:"session"`  // session ID
	ID            int64  `json:"id" example:"1"`                         // 提問ID
	Action        string `json:"action" example:"like、unlike"`           // 訊息
	Limit         int64  `json:"limit" example:"100"`                    // 需要的資料數量
	Offdet        int64  `json:"offdet" example:"100"`                   // 跳過的資料數量
	MessageStatus string `json:"message_status" example:"yes、no、review"` // 訊息
	Like          string `json:"like" example:"yes、no"`                  // 主持端是否按讚
	// Method string `json:"method" example:"sign"`        // 方法
}

// HostQuestionWriteParam 主持端提問牆回傳訊息
type HostQuestionWriteParam struct {
	QuestionsOrderByTime []models.QuestionUserModel // 按時間排序
	// QuestionsOrderByLikes []models.QuestionUserModel // 按熱門排序
	Error                string `json:"error" example:"error message"`       // 錯誤訊息
	YesQuestionAmount    int64  `json:"yes_question_amount" example:"10"`    // 已通過提問數量
	NoQuestionAmount     int64  `json:"no_question_amount" example:"10"`     // 未通過提問數量
	ReviewQuestionAmount int64  `json:"review_question_amount" example:"10"` // 審核提問數量
	// UserQuestionAmount       int64 `json:"user_question_amount" example:"10"`        // 該用戶訊息數量
	QuestionLikesAmount      int64 `json:"question_likes_amount" example:"10"`       // 所有提問總讚數數量
	HostQuestionLikeAmount   int64 `json:"host_question_like_amount" example:"10"`   // 主持人已按讚提問數量
	HostQuestionUnlikeAmount int64 `json:"host_question_unlike_amount" example:"10"` // 主持人未按讚提問數量
	IsSendPeople             bool  `json:"is_send_people" example:"true"`            // 是否傳遞人數
}

// GuestQuestionReadParam 用戶端提問牆接收訊息
type GuestQuestionReadParam struct {
	// SessionID string `json:"session_id" example:"session"`  // session ID
	ID            int64  `json:"id" example:"1"`                         // 提問ID
	Action        string `json:"action" example:"like、unlike、question"`  // 按讚、取消讚、提問
	Message       string `json:"message" example:"message"`              // 訊息
	MessageStatus string `json:"message_status" example:"yes、no、review"` // 訊息
	Like          string `json:"like" example:"yes、no"`                  // 主持端是否按讚
	Limit         int64  `json:"limit" example:"100"`                    // 需要的資料數量
	Offset        int64  `json:"offset" example:"100"`                   // 跳過的資料數量
	// Method        string `json:"method" example:"sign"`                  // 方法

}

// GuestQuestionWriteParam 用戶端提問牆回傳訊息
type GuestQuestionWriteParam struct {
	Questions []models.QuestionUserModel // 所有提問
	// UserQuestions []models.QuestionUserModel        // 我的提問
	LikesRecords             []models.QuestionLikesRecordModel // 按讚紀錄
	Error                    string                            `json:"error" example:"error message"`            // 錯誤訊息
	YesQuestionAmount        int64                             `json:"yes_question_amount" example:"10"`         // 已通過提問數量
	NoQuestionAmount         int64                             `json:"no_question_amount" example:"10"`          // 未通過提問數量
	ReviewQuestionAmount     int64                             `json:"review_question_amount" example:"10"`      // 審核提問數量
	UserQuestionAmount       int64                             `json:"user_question_amount" example:"10"`        // 該用戶提問數量
	UserLikesRecordAmount    int64                             `json:"user_likes_record_amount" example:"10"`    // 該用戶按讚紀錄數量
	QuestionLikesAmount      int64                             `json:"question_likes_amount" example:"10"`       // 所有提問總讚數數量
	HostQuestionLikeAmount   int64                             `json:"host_question_like_amount" example:"10"`   // 主持人已按讚提問數量
	HostQuestionUnlikeAmount int64                             `json:"host_question_unlike_amount" example:"10"` // 主持人未按讚提問數量
	IsSendPeople             bool                              `json:"is_send_people" example:"true"`            // 是否傳遞人數
}

// @Summary 即時回傳用戶端提問資料(所有提問資料、我的提問資料)、按讚紀錄
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param session_id query string true "session ID"
// @param body body GuestQuestionReadParam true "id、message、content param"
// @Success 200 {array} GuestQuestionWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/guest/question [get]
func (h *Handler) GuestQuestionWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		sessionID         = ctx.Query("session_id")
		wsConn, conn, err = NewWebsocketConn(ctx)
		user              models.LineModel
		// result            GuestQuestionReadParam
		// oldQuestions = make([]models.QuestionUserModel, 0)        // 舊的提問紀錄
		// oldRecords   = make([]models.QuestionLikesRecordModel, 0) // 就的按讚紀錄
		// limit      = 20
	)
	// log.Println("開啟用戶端回傳即時提問資料(所有提問資料、我的提問資料)、按讚紀錄ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || sessionID == "" {
		b, _ := json.Marshal(HostQuestionWriteParam{
			Error: "錯誤: 無法辨識活動.用戶、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if user, _, err = h.getUser(sessionID, ""); err != nil {
		// 利用session值取得用戶資料
		b, _ := json.Marshal(HostQuestionWriteParam{
			Error: "錯誤: 無法辨識用戶，請重新操作"})
		conn.WriteMessage(b)
		return
	}

	go func() {
		for {
			// ###優化前，定頻###
			var (
				questionUserModel     models.QuestionUserModel
				userLikesRecordAmount int64
				// questions                                                 = make([]models.QuestionUserModel, 0) // 所有提問紀錄，未過濾的資料

			)

			// 查詢不同類型聊天訊息數量(玩家端)
			questionUserModel, _ = models.DefaultQuestionUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindGuestQuestionAmount(activityID, user.UserID)

			// 用戶按讚記錄數量
			userLikesRecordAmount, _ = models.DefaultQuestionLikesRecordModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindUserLikesRecordAmount(activityID, user.UserID)

			b, _ := json.Marshal(GuestQuestionWriteParam{
				YesQuestionAmount:        questionUserModel.YesQuestionAmount,
				NoQuestionAmount:         questionUserModel.NoQuestionAmount,
				ReviewQuestionAmount:     questionUserModel.ReviewQuestionAmount,
				UserQuestionAmount:       questionUserModel.UserQuestionAmount,
				UserLikesRecordAmount:    userLikesRecordAmount,
				QuestionLikesAmount:      questionUserModel.QuestionLikesAmount,
				HostQuestionLikeAmount:   questionUserModel.HostQuestionLikeAmount,
				HostQuestionUnlikeAmount: questionUserModel.HostQuestionUnlikeAmount,
				IsSendPeople:             true,
			})
			conn.WriteMessage(b)
			// ###優化前，定頻###

			if conn.isClose { // ws關閉
				// 關閉用戶提問牆ws時，清除redis裡的用戶按讚紀錄資料
				// h.redisConn.HashDelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID, user.UserID)
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var (
			result GuestQuestionReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// 關閉用戶提問牆ws時，清除redis裡的用戶按讚紀錄資料
			// h.redisConn.HashDelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID, user.UserID)
			// log.Println("關閉用戶端回傳即時提問資料(所有提問資料、我的提問資料)、按讚紀錄ws")
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offset != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得提問訊息資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		if result.Action != "" {
			var (
				status, like string
			)

			if result.Action == "like" || result.Action == "unlike" {
				// 按讚或收回讚
				if result.ID == 0 {
					b, _ := json.Marshal(HostQuestionWriteParam{
						Error: "錯誤: 取得提問資料發生問題"})
					conn.WriteMessage(b)
					return
				}

				// 更新讚數資料
				if err = models.DefaultQuestionUserModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateGuestLikes(true, result.ID,
						activityID, result.Action); err != nil {
					b, _ := json.Marshal(HostQuestionWriteParam{
						Error: "錯誤: 更新讚數資料發生問題"})
					conn.WriteMessage(b)
					return
				}

				if result.Action == "like" {
					if err = models.DefaultQuestionLikesRecordModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						Add(
							true, models.NewQuestionLikesRecordModel{
								ActivityID: activityID,
								QuestionID: result.ID,
								UserID:     user.UserID,
							}); err != nil {
						b, _ := json.Marshal(HostQuestionWriteParam{
							Error: "錯誤: 新增用戶按讚紀錄資料發生問題"})
						conn.WriteMessage(b)
						return
					}
				} else if result.Action == "unlike" {
					if err = models.DefaultQuestionLikesRecordModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						Delete(true, result.ID,
							activityID, user.UserID); err != nil {
						b, _ := json.Marshal(HostQuestionWriteParam{
							Error: "錯誤: 刪除用戶按讚紀錄資料發生問題"})
						conn.WriteMessage(b)
						return
					}
				}
			} else if result.Action == "question" {
				if (result.MessageStatus != "yes" && result.MessageStatus != "no" &&
					result.MessageStatus != "review") || result.Message == "" {
					b, _ := json.Marshal(HostQuestionWriteParam{
						Error: "錯誤: 訊息資料發生問題，請輸入有效的訊息資料"})
					conn.WriteMessage(b)
					return
				}

				// 用戶提問
				if _, err = models.DefaultQuestionUserModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Add(true,
						models.EditQuestionUserModel{
							ActivityID:    activityID,
							UserID:        user.UserID,
							Message:       result.Message,
							MessageStatus: result.MessageStatus,
						}); err != nil {
					b, _ := json.Marshal(HostQuestionWriteParam{
						Error: "錯誤: 新增提問資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			} else if result.Action == "find" {
				// 查詢提問資料
				status = result.MessageStatus
				like = result.Like

				// var isRedis bool
				// if result.Limit != 0 || result.Offset != 0 ||
				// 	status != "" || like != "" {
				// 	// 取部分資料，用資料表
				// 	isRedis = false
				// } else {
				// 	// 取全部資料，用redis
				// 	isRedis = true
				// }

				// 所有提問資料(時間排序)
				questions, _ := models.DefaultQuestionUserModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, "",
						status, like,
						result.Limit, result.Offset)

				// 用戶按讚紀錄
				records, _ := models.DefaultQuestionLikesRecordModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, user.UserID, 0, 0)

				// 回傳提問資料、按讚紀錄
				b, _ := json.Marshal(GuestQuestionWriteParam{
					Questions:    questions,
					LikesRecords: records,
					IsSendPeople: false,
				})
				conn.WriteMessage(b)
			}
		}
	}
}

// @Summary 即時回傳主持人端提問資料(新舊排序、熱門排序)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @param body body HostQuestionReadParam true "id、message param"
// @Success 200 {array} HostQuestionWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/host/question [get]
func (h *Handler) HostQuestionWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            HostQuestionReadParam
		// n                        = 1
		// oldQuestions = make([]models.QuestionUserModel, 0) // 舊的提問紀錄
	)
	// log.Println("開啟主持人端回傳即時提問資料(新舊排序、熱門排序)ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" {
		b, _ := json.Marshal(HostQuestionWriteParam{
			Error: "錯誤: 無法辨識活動、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	// 測試活動:dQSrpxSJlexj6G9EIph4
	// h.redisConn.DelCache(config.QUESTION_REDIS + activityID)
	// h.redisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID)
	// h.redisConn.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + activityID)
	// h.redisConn.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + activityID)
	// 測試活動:dQSrpxSJlexj6G9EIph4

	go func() {
		for {
			// 測試活動:dQSrpxSJlexj6G9EIph4
			// if activityID == "dQSrpxSJlexj6G9EIph4" {
			// 	if n <= 8000 {
			// 		for i := 0; i < 5; i++ {
			// 			// 用戶提問
			// 			models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
			// 				SetRedisConn(h.redisConn).Add(true,
			// 				models.EditQuestionUserModel{
			// 					ActivityID:    activityID,
			// 					UserID:        "U3140fd73cfd35dd992668ab3b6efdae9",
			// 					Message:       "result.Message",
			// 					MessageStatus: "yes",
			// 				})
			// 			n++
			// 		}
			// 	}
			// }
			// 測試活動:dQSrpxSJlexj6G9EIph4

			// ###優化前，定頻###
			var (
				questionUserModel models.QuestionUserModel
			)

			// 查詢不同類型聊天訊息數量(主持端)
			questionUserModel, _ = models.DefaultQuestionUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindHostQuestionAmount(activityID)

			b, _ := json.Marshal(HostQuestionWriteParam{
				YesQuestionAmount:        questionUserModel.YesQuestionAmount,
				NoQuestionAmount:         questionUserModel.NoQuestionAmount,
				ReviewQuestionAmount:     questionUserModel.ReviewQuestionAmount,
				QuestionLikesAmount:      questionUserModel.QuestionLikesAmount,
				HostQuestionLikeAmount:   questionUserModel.HostQuestionLikeAmount,
				HostQuestionUnlikeAmount: questionUserModel.HostQuestionUnlikeAmount,
				IsSendPeople:             true,
			})
			conn.WriteMessage(b)
			// ###優化前，定頻###

			if conn.isClose { // ws關閉
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		var (
			result HostQuestionReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉主持人端回傳即時提問資料(新舊排序、熱門排序)ws")
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Offdet != 0 && result.Limit == 0 {
			b, _ := json.Marshal(GameParam{Error: "錯誤: 無法取得提問訊息資料(limit.offset參數發生錯誤)"})
			conn.WriteMessage(b)
			return
		}

		if (result.ID != 0 && (result.Action == "like" || result.Action == "unlike")) ||
			result.Action == "find" {
			var like string
			if result.Action == "like" {
				like = "yes"
			} else if result.Action == "unlike" {
				like = "no"
			}
			// else if result.Action == "find" {
			// 查詢提問資料
			// }

			if result.Action == "like" || result.Action == "unlike" {
				// 更新主持人按讚紀錄
				if err = models.DefaultQuestionUserModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateHostLikes(true, result.ID,
						activityID, like); err != nil {
					b, _ := json.Marshal(HostQuestionWriteParam{
						Error: "錯誤: 更新主持人按讚資料發生問題"})
					conn.WriteMessage(b)
					return
				}
			} else if result.Action == "find" {
				// var isRedis bool
				// if result.Limit != 0 || result.Offdet != 0 ||
				// 	result.MessageStatus != "" || result.Like != "" {
				// 	// 取部分資料，用資料表
				// 	isRedis = false
				// } else {
				// 	// 取全部資料，用redis
				// 	isRedis = true
				// }

				// 所有提問資料(時間排序)
				questions, _ := models.DefaultQuestionUserModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, activityID, "",
						result.MessageStatus, result.Like,
						result.Limit, result.Offdet)

				// 回傳提問資料(時間排序)
				b, _ := json.Marshal(HostQuestionWriteParam{
					QuestionsOrderByTime: questions,
					IsSendPeople:         false,
				})
				conn.WriteMessage(b)
			}
		}
	}
}

// yesQuestionAmount, noQuestionAmount, reviewQuestionAmount int64
// questionLikesAmount                                       int64
// userQuestionAmount, userLikesRecordAmount                 int64
// hostQuestionLikeAmount, hostQuestionUnlikeAmount          int64

// // 通過資料數量
// yesQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "yes", "")

// // 未通過資料數量
// noQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "no", "")

// // 審核資料數量
// reviewQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "review", "")

// 用戶提問資料數量
// userQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, user.UserID, "", "")

// // 所有提問總讚數資料數量
// questionLikesAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionLikesAmount(activityID)

// // 主持人已按讚提問資料數量
// hostQuestionLikeAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "", "yes")

// // 主持人未按讚提問資料數量
// hostQuestionUnlikeAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "", "no")

// yesQuestionAmount, noQuestionAmount, reviewQuestionAmount int64
// questionLikesAmount                                       int64
// hostQuestionLikeAmount, hostQuestionUnlikeAmount          int64

// 通過資料數量
// yesQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "yes", "")

// // 未通過資料數量
// noQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "no", "")

// // 審核資料數量
// reviewQuestionAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "review", "")

// // 所有提問總讚數資料數量
// questionLikesAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionLikesAmount(activityID)

// // 主持人已按讚提問資料數量
// hostQuestionLikeAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "", "yes")

// // 主持人未按讚提問資料數量
// hostQuestionUnlikeAmount, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).FindQuestionAmount(activityID, "", "", "no")

// ###優化前，定頻###
// var (
// 	questions    = make([]models.QuestionUserModel, 0) // 所有提問紀錄，未過濾的資料
// 	newQuestions = make([]models.QuestionUserModel, 0) // 過濾後的資料
// 	isSend       bool                                  // 是否回傳訊息至前端
// )

// // 所有提問資料(時間排序)，包含通過、未通過、審核中
// questions, _ = models.DefaultQuestionUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).Find(true, activityID)

// // 用戶按讚紀錄
// newRecords, _ := models.DefaultQuestionLikesRecordModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).Find(true, activityID, user.UserID)

// // log.Println("玩家端提問變化: ", len(questions), len(oldQuestions))
// // log.Println("玩家端按讚變化: ", len(records), len(oldRecords))

// // 過濾黑名單資料
// for _, question := range questions {
// 	if !h.IsBlackStaff(activityID, "", "question", question.UserID) {
// 		// question.MessageStatus == "yes" {
// 		// 不是黑名單
// 		newQuestions = append(newQuestions, question)
// 	}
// }

// // 判斷資料是否有變化(提問資料、按讚資料)
// if len(newQuestions) != len(oldQuestions) ||
// 	len(newRecords) != len(oldRecords) {
// 	// log.Println("玩家端提問資料長度不一樣")
// 	isSend = true
// } else if len(newQuestions) == len(oldQuestions) {
// 	// 判斷原始資料是否有變動
// 	for i := 0; i < len(oldQuestions); i++ {
// 		if oldQuestions[i].MessageStatus != newQuestions[i].MessageStatus || // 判斷提問狀態是否改變
// 			oldQuestions[i].Likes != newQuestions[i].Likes { // 判斷讚數是否改變
// 			// log.Println("玩家端提問資料有變動: ", i)
// 			isSend = true
// 			break
// 		}
// 	}
// }

// if isSend { // 回傳資料至前端
// 	// 將新的資料放進old參數中
// 	oldQuestions = newQuestions
// 	oldRecords = newRecords

// 	b, _ := json.Marshal(GuestQuestionWriteParam{
// 		Questions:    newQuestions,
// 		LikesRecords: newRecords,
// 	})
// 	conn.WriteMessage(b)

// }
// ###優化前，定頻###

// ###優化前，定頻###
// var (
// 	questions    = make([]models.QuestionUserModel, 0) // 所有提問紀錄，未過濾的資料
// 	newQuestions = make([]models.QuestionUserModel, 0) // 過濾後的資料
// 	isSend       bool                                  // 是否回傳訊息至前端
// )

// // 所有提問資料(時間排序)
// questions, _ = models.DefaultQuestionUserModel().
// 	SetDbConn(h.dbConn).SetRedisConn(h.redisConn).Find(true, activityID)

// // 過濾資料，判斷資料是否為黑名單(時間排序)
// for _, question := range questions {
// 	if !h.IsBlackStaff(activityID, "", "question", question.UserID) {
// 		newQuestions = append(newQuestions, question)
// 	}
// }

// // 判斷提問資料長度是否出現變化
// if len(newQuestions) != len(oldQuestions) {
// 	// log.Println("主持端提問資料長度不一樣")
// 	isSend = true
// } else if len(newQuestions) == len(oldQuestions) {
// 	// 判斷原始資料是否有變動
// 	for i := 0; i < len(oldQuestions); i++ {
// 		if oldQuestions[i].MessageStatus != newQuestions[i].MessageStatus || // 判斷提問狀態是否改變
// 			oldQuestions[i].Likes != newQuestions[i].Likes { // 判斷讚數是否改變
// 			// log.Println("主持端提問資料有變動: ", i)

// 			isSend = true
// 			break
// 		}
// 	}
// }

// if isSend { // 回傳資料至前端
// 	// 將新的資料放進old參數中
// 	oldQuestions = newQuestions

// 	b, _ := json.Marshal(HostQuestionWriteParam{
// 		QuestionsOrderByTime: newQuestions,
// 	})
// 	conn.WriteMessage(b)

// }
// ###優化前，定頻###

// var (
// 	questionsOrderByTime  = make([]models.QuestionUserModel, 0)
// 	questionsOrderByLikes = make([]models.QuestionUserModel, 0)
// )

// var (
// 	questionsOrderByTime  = make([]models.QuestionUserModel, 0)
// 	questionsOrderByLikes = make([]models.QuestionUserModel, 0)
// )

// likesRecordModel  = models.DefaultQuestionLikesRecordModel().SetDbConn(h.dbConn)
// questions         = make([]models.QuestionUserModel, 0)
// userQuestions     = make([]models.QuestionUserModel, 0)
// records           = make([]models.QuestionLikesRecordModel, 0)

// questionModel         = models.DefaultQuestionUserModel().SetDbConn(h.dbConn)
// questionsOrderByTime  = make([]models.QuestionUserModel, 0)
// questionsOrderByLikes = make([]models.QuestionUserModel, 0)

// // 提問資料(時間排序)
// questionsOrderByTime, _ = h.getQuestions("activity_id", activityID, "send_time")
// // 提問資料(熱門排序)
// questionsOrderByLikes, _ = h.getQuestions("activity_id", activityID, "likes")
// // 回傳提問資料(時間排序、讚數排序)
// b, _ := json.Marshal(HostQuestionWriteParam{
// 	QuestionsOrderByTime:  questionsOrderByTime,
// 	QuestionsOrderByLikes: questionsOrderByLikes,
// })
// conn.WriteMessage(b)

// QuestionUserModel 用戶提問資料欄位
// type QuestionUserModel struct {
// 	ID         int64  `json:"id" example:"1"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	UserID     string `json:"user_id" example:"user_id"`
// 	Name       string `json:"name" example:"name"`
// 	Avatar     string `json:"avatar" example:"avatar"`
// 	Content    string `json:"content" example:"content"`
// 	Likes      int64  `json:"likes" example:"1"`
// 	SendTime   string `json:"send_time" example:"2022/01/01 00:00:00"`
// 	Like       string `json:"like" example:"yes、no"`
// }

// QuestionLikesRecordModel 資料表欄位
// type QuestionLikesRecordModel struct {
// 	ID         int64  `json:"id" example:"1"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	QuestionID int64  `json:"question_id" example:"1"`
// 	UserID     string `json:"user_i
