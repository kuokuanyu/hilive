package controller

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary LINE MESSAGE CALLBACK
// @Tags LINE Auth
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param user_id query string false "用戶ID"
// @Success 200 {array} response.Response
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /message/line/callback [post]
func (h *Handler) LineMessageCallback(ctx *gin.Context) {
	// var client = NewLineBotClient(ctx)
	var (
		activityID                        = ctx.Request.FormValue("activity_id")
		userID                            = ctx.Request.FormValue("user_id")
		chatbotSecret, chatbotToken, line string
		err                               error
	)

	// LineBotlock.Lock()
	// defer LineBotlock.Unlock()
	if ctx.Request.Host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	// 取得LINE機器人資訊
	if activityID != "" {
		// 取得活動資訊
		activityModel, err := h.getActivityInfo(false, activityID)
		if err != nil || activityModel.ID == 0 {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得活動資訊，請重新查詢",
			})
			return
		}

		chatbotSecret = activityModel.ActivityChatbotSecret
		chatbotToken = activityModel.ActivityChatbotToken
	} else if userID != "" {
		// 取得用戶資訊
		userModel, err := models.DefaultUserModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, true, "users.user_id", userID)
		if err != nil || userModel.ID == 0 {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID",
			})

			return
		}

		if userModel.ChannelID != "" && userModel.ChannelSecret != "" &&
			userModel.ChatbotSecret != "" && userModel.ChatbotToken != "" &&
			userModel.LineID != "" {
			chatbotSecret = userModel.ChatbotSecret
			chatbotToken = userModel.ChatbotToken
			line = userModel.LineID
		} else {
			chatbotSecret = config.CHATBOT_SECRET
			chatbotToken = config.CHATBOT_TOKEN
		}
	} else {
		chatbotSecret = config.CHATBOT_SECRET
		chatbotToken = config.CHATBOT_TOKEN
	}

	client := NewLineBotClient(ctx, chatbotSecret, chatbotToken)
	events, err := client.ParseRequest(ctx.Request)
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})

		return
	}

	for _, event := range events {
		var (
			lineModel models.LineModel
			// message   = "歡迎將本帳號設為好友!"
			lineID string
		)
		if event.Source != nil {
			lineID = event.Source.UserID
		}

		// log.Println("user_id: ", userID)
		// log.Println("lineID: ", lineID)

		// 取得用戶資料
		res, _, err := getProfile(client, lineID)
		if err != nil {
			// 被封鎖，將用戶資料更新為未加入好友
			models.DefaultLineModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				UpdateUser(models.LineModel{
					UserID: lineID,
					Friend: "no",
				})

			// 刪除redis資訊
			// h.redisConn.DelCache(config.AUTH_USERS_REDIS + lineID)

			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})

			return
		}

		for l := 0; l < MaxRetries; l++ {
			// 上鎖
			ok, _ := h.acquireLock(config.LINE_USERS_LOCK_REDIS+lineID, LockExpiration)
			if ok == "OK" {
				// 釋放鎖
				// defer h.releaseLock(config.LINE_USERS_LOCK_REDIS + lineID)

				// 加入funciton

				if lineModel, err = models.DefaultLineModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(false, "", "user_id", lineID); err != nil {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 取得用戶資料發生問題",
					})

					// 釋放鎖
					h.releaseLock(config.LINE_USERS_LOCK_REDIS + lineID)
					return
				}

				// 是否存在用戶資料
				if lineModel.UserID == "" {

					// message = "歡迎將本帳號設為好友!"
					_, err = models.DefaultLineModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						Add(models.LineModel{
							UserID:   lineID,
							Name:     res.DisplayName,
							Avatar:   res.PictureURL,
							Email:    "",
							Identify: lineID,
							Friend:   "yes",
							Line:     line,
							Device:   "line",
						})

				} else {
					// 只更新friend資料
					if lineModel.Friend == "no" {
						// log.Println("好友，更新friend資料")
						err = models.DefaultLineModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							UpdateUser(models.LineModel{
								UserID: lineID,
								// Name:   res.DisplayName,
								// Avatar: res.PictureURL,
								Friend: "yes",
							})

						// if value, _ := h.redisConn.HashGetCache(config.AUTH_USERS_REDIS+lineID, "user_id"); value != "" {
						// 	h.redisConn.HashMultiSetCache([]interface{}{config.AUTH_USERS_REDIS + lineID,
						// 		"name", res.DisplayName, "avatar", res.PictureURL, "friend", "yes"})
						// }
					}
				}

				// 釋放鎖
				h.releaseLock(config.LINE_USERS_LOCK_REDIS + lineID)
				break
			}

			// 鎖被佔用，稍微延遲後重試
			time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
		}

		if err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID + "+" + lineID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}

		// 自動回覆訊息
		// if event.Type == linebot.EventTypeMessage {
		// 	// 訊息
		// 	message = "自動回覆"

		// 	switch message := event.Message.(type) {
		// 	case *linebot.TextMessage:
		// 		switch {
		// 		case message.Text == "1":
		// 			pushmessagewithbutton(ctx, chatbotSecret, chatbotToken, lineID,
		// 				"按鈕", "訊息")
		// 		case message.Text == "2":
		// 			pushmessagewithbuttons(ctx, chatbotSecret, chatbotToken, lineID,
		// 				"按鈕", "按鈕2", "訊息")
		// 		case message.Text == "3":
		// 			pushmessagebyimage(ctx, chatbotSecret, chatbotToken, lineID,
		// 				"https://profile.line-scdn.net/0hCSeMHjdLHHxYNAiRsHFjK2RxEhEvGho0IFYBGS8yER5xBQkuZFJXEn9hFh8nVF0qZFtaGy41Ekxz",
		// 				"訊息")
		// 		case message.Text == "4":
		// 			pushmessagewithcamara(ctx, chatbotSecret, chatbotToken, lineID,
		// 				"按鈕", "按鈕2", "訊息")
		// 		}
		// 	}
		// } else if event.Type == linebot.EventTypeFollow {
		// 	// 追隨
		// 	message = "加好友"
		// 	if err = pushMessage(ctx, chatbotSecret, chatbotToken, lineID, message); err != nil {
		// 		response.Error(ctx, err.Error())
		// 		return
		// 	}
		// } else if event.Type == linebot.EventTypeUnfollow {
		// 	// 取消追隨
		// 	fmt.Println("取消好友")
		// 	message = "取消好友"
		// 	if err = pushMessage(ctx, chatbotSecret, chatbotToken, lineID, message); err != nil {
		// 		response.Error(ctx, err.Error())
		// 		return
		// 	}
		// }

		// if err = pushMessage(ctx, chatbotSecret, chatbotToken, lineID, message); err != nil {
		// 	response.Error(ctx, err.Error())
		// 	return
		// }
	}
}

// func pushmessagewithbutton(ctx *gin.Context, secret, token, userID, btn1, message string) error {
// 	if _, err := NewLineBotClient(ctx, secret, token).PushMessage(
// 		userID,
// 		linebot.NewTextMessage(message).
// 			WithQuickReplies(linebot.NewQuickReplyItems(
// 				linebot.NewQuickReplyButton(
// 					"",
// 					linebot.NewMessageAction(btn1, btn1)),
// 			))).Do(); err != nil {
// 		return errors.New("錯誤: 傳送Line訊息發生問題(一個按鈕)，請重新傳送")
// 	}
// 	return nil
// }

// func pushmessagewithbuttons(ctx *gin.Context, secret, token, userID, btn1, btn2, message string) error {
// 	if _, err := NewLineBotClient(ctx, secret, token).PushMessage(userID,
// 		linebot.NewTextMessage(message).WithQuickReplies(linebot.NewQuickReplyItems(
// 			linebot.NewQuickReplyButton(
// 				"https://profile.line-scdn.net/0hCSeMHjdLHHxYNAiRsHFjK2RxEhEvGho0IFYBGS8yER5xBQkuZFJXEn9hFh8nVF0qZFtaGy41Ekxz",
// 				linebot.NewMessageAction(btn1, "按鈕1")),
// 			linebot.NewQuickReplyButton(
// 				"",
// 				linebot.NewMessageAction(btn2, btn2)),
// 		)),
// 	).Do(); err != nil {
// 		return errors.New("錯誤: 傳送Line訊息發生問題(兩個按鈕)，請重新傳送")
// 	}
// 	return nil
// }

// func pushmessagebyimage(ctx *gin.Context, secret, token, userID, image, message string) error {
// 	if _, err := NewLineBotClient(ctx, secret, token).PushMessage(userID,
// 		linebot.NewImageMessage(image, image),
// 		linebot.NewTextMessage(message).
// 			WithQuickReplies(linebot.NewQuickReplyItems(
// 				linebot.NewQuickReplyButton(
// 					"",
// 					linebot.NewMessageAction("取消付款", "取消付款")),
// 			))).Do(); err != nil {
// 		return errors.New("錯誤: 傳送Line訊息發生問題(圖片)，請重新傳送")
// 	}
// 	return nil
// }

// func pushmessagewithcamara(ctx *gin.Context, secret, token, userID, but1, but2, message string) error {
// 	if _, err := NewLineBotClient(ctx, secret, token).PushMessage(
// 		userID,
// 		linebot.NewTextMessage(message).
// 			WithQuickReplies(linebot.NewQuickReplyItems(
// 				linebot.NewQuickReplyButton(
// 					"",
// 					linebot.NewCameraAction(but1)),
// 				linebot.NewQuickReplyButton(
// 					"",
// 					linebot.NewCameraRollAction(but2)),
// 			))).Do(); err != nil {
// 		return errors.New("錯誤: 傳送Line訊息發生問題(相機)，請重新傳送")
// 	}
// 	return nil
// }
