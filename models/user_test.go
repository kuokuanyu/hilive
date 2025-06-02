package models

import (
	"errors"
	"io"
	"os"
	"sync"
)

// 建立資料夾、整理所有圖片檔案
// func Test_Add_File(t *testing.T) {
// 	// 查詢所有活動
// 	activitys, err := DefaultActivityModel().SetDbConn(conn).FindOpenActivitys()
// 	if err != nil {
// 		t.Error("錯誤: 取得用戶下的所有活動資料錯誤", err)
// 	} else {
// 		fmt.Println("所有活動數: ", len(activitys))
// 	}

// 	for _, activity := range activitys {
// 		// 建立簽名牆簽到資料夾
// 		os.MkdirAll("../uploads/"+activity.UserID+"/"+activity.ActivityID+"/interact/game/3DGachaMachine", os.ModePerm)
// 	}
// }

// 建立資料夾、整理所有圖片檔案
// func Test_Add_File(t *testing.T) {
// 	// 所有用戶資料
// 	users, err := DefaultUserModel().SetDbConn(conn).FindUsers()
// 	if err != nil {
// 		t.Error("錯誤: 取得所有用戶資料錯誤", err)
// 	} else {
// 		// fmt.Println("所有用戶數: ", len(users))
// 	}

// 	// os.Mkdir("uploads/", os.ModePerm)

// 	for _, user := range users {
// 		// if user.UserID == "admin" { // 最後要拿掉
// 		// 建立用戶資料夾
// 		// os.MkdirAll("../uploads/"+user.UserID, os.ModePerm)

// 		// 查詢用戶下的所有活動
// 		activitys, err := DefaultActivityModel().SetDbConn(conn).FindActivityPermissions("activity.user_id", user.UserID)
// 		if err != nil {
// 			t.Error("錯誤: 取得用戶下的所有活動資料錯誤", err)
// 		} else {
// 			// fmt.Println("所有活動數: ", len(activitys))
// 		}

// 		for _, activity := range activitys {
// 			// if activity.ActivityID == "Bfon6SaV6ORhmuDQUioI" { // 最後要拿掉
// 			// // 建立活動資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID, os.ModePerm)

// 			// // 建立活動介紹資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/info/introduce", os.ModePerm)
// 			// // 建立活動嘉賓資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/info/guest", os.ModePerm)

// 			// // 建立自定義欄位設置資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/applysign/customize", os.ModePerm)
// 			// // 建立QRcode自定義欄位設置資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/applysign/qrcode", os.ModePerm)

// 			// // 建立訊息區設置資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/wall/message", os.ModePerm)
// 			// // 建立主題區資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/wall/topic", os.ModePerm)
// 			// // 建立提問區資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/wall/question", os.ModePerm)

// 			// // 建立一般簽到資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/sign/general", os.ModePerm)
// 			// // 建立3D簽到資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/sign/threed", os.ModePerm)

// 			// // 建立遊戲抽獎資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/lottery", os.ModePerm)
// 			// // 建立搖紅包資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/redpack", os.ModePerm)
// 			// // 建立套紅包資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/ropepack", os.ModePerm)
// 			// // 建立搖號抽獎資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/draw_numbers", os.ModePerm)
// 			// // 建立敲敲樂資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/whack_mole", os.ModePerm)
// 			// // 建立鑑定師資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/monopoly", os.ModePerm)
// 			// // 建立快問快答資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/QA", os.ModePerm)
// 			// // 建立拔河遊戲資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/tugofwar", os.ModePerm)
// 			// // 建立賓果遊戲資料夾
// 			// os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/bingo", os.ModePerm)

// 			// 查詢活動下的所有遊戲場次
// 			// games, err := db.Table("activity_game").WithConn(conn).
// 			// 	Where("activity_id", "=", activity.ActivityID).All()
// 			// if err != nil {
// 			// 	t.Error("錯誤: 取得活動下的所有遊戲場次資料錯誤", err)
// 			// } else {
// 			// 	// 27
// 			// 	// fmt.Println("所有遊戲數: ", len(games))
// 			// }
// 			// for _, game := range games {
// 			// 	gameType, _ := game["game"].(string)
// 			// 	gameID, _ := game["game_id"].(string)

// 			// 	// 建立遊戲的資料夾
// 			// 	os.MkdirAll("../uploads/"+user.UserID+"/"+activity.ActivityID+"/interact/game/"+gameType+"/"+gameID, os.ModePerm)

// 			// 	// 查詢遊戲下的所有獎品資料
// 			// 	prizes, err := DefaultPrizeModel().SetDbConn(conn).
// 			// 		FindPrizes(false, gameID)
// 			// 	if err != nil {
// 			// 		t.Error("錯誤: 取得遊戲下的所有獎品資料錯誤", err)
// 			// 	}

// 			// 	for _, prize := range prizes {
// 			// 		if !strings.Contains(prize.PrizePicture, "system") {
// 			// 			// 原始圖片路徑
// 			// 			oldRoute := "../uploads/" + strings.Split(prize.PrizePicture, "/admin/uploads/")[1]
// 			// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" + gameID + "/" + strings.Split(prize.PrizePicture, "/admin/uploads/")[1]
// 			// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// 			err = copy(oldRoute, newRoute)
// 			// 			if err != nil {
// 			// 				t.Error(err)
// 			// 			}

// 			// 			// 更新圖片資料
// 			// 			DefaultPrizeModel().SetDbConn(conn).
// 			// 				Update(false, gameType, EditPrizeModel{
// 			// 					PrizeID:      prize.PrizeID,
// 			// 					PrizePicture: strings.Split(prize.PrizePicture, "/admin/uploads/")[1],
// 			// 				})
// 			// 		}
// 			// 	}
// 			// }

// 			// 用戶頭像
// 			// if !strings.Contains(user.Avatar, "system") {
// 			// 	file := strings.Split(user.Avatar, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/user." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultUserModel().SetDbConn(conn).
// 			// 		Update(false, EditUserModel{
// 			// 			Avatar: "user." + strings.Split(file, ".")[1],
// 			// 		}, "user_id", user.UserID)
// 			// }

// 			// 查詢活動下的活動介紹資料
// 			// introduces, err := DefaultIntroduceModel().SetDbConn(conn).
// 			// 	Find("activity_introduce.activity_id", activity.ActivityID)
// 			// if err != nil {
// 			// 	t.Error("錯誤: 取得活動下的活動介紹資料錯誤", err)
// 			// } else {
// 			// 	// fmt.Println("活動介紹數: ", len(introduces))
// 			// }
// 			// for _, introduce := range introduces {
// 			// 	// 判斷活動介紹類型
// 			// 	if introduce.IntroduceType == "picture" {
// 			// 		// 原始圖片路徑
// 			// 		oldRoute := "../uploads/" + introduce.Content
// 			// 		newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/info/introduce/" + introduce.Content

// 			// 		err = copy(oldRoute, newRoute)
// 			// 		if err != nil {
// 			// 			t.Error(err)
// 			// 		}

// 			// 		// 更新圖片資料
// 			// 		// DefaultIntroduceModel().SetDbConn(conn).
// 			// 		// 	Update(EditIntroduceModel{
// 			// 		// 		ID:         strconv.Itoa(int(introduce.ID)),
// 			// 		// 		ActivityID: activity.ActivityID,
// 			// 		// 		Content:    "/admin/uploads/" + user.UserID + "/" + activity.ActivityID + "/info/introduce/" + introduce.Content,
// 			// 		// 	})

// 			// 	}
// 			// }

// 			// 查詢活動下的活動嘉賓資料
// 			guests, err := DefaultGuestModel().SetDbConn(conn).
// 				Find("activity_guest.activity_id", activity.ActivityID)
// 			if err != nil {
// 				t.Error("錯誤: 取得活動下的活動嘉賓資料錯誤", err)
// 			} else {
// 				// fmt.Println("活動嘉賓數: ", len(guests))
// 			}
// 			for _, guest := range guests {
// 				if !strings.Contains(guest.Avatar, "system") {
// 					// 原始圖片路徑
// 					oldRoute := "../uploads/" + strings.Split(guest.Avatar, "/admin/uploads/")[1]
// 					newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/info/guest/" + strings.Split(guest.Avatar, "/admin/uploads/")[1]
// 					// fmt.Println("oldRoute: ", oldRoute)
// 					err = copy(oldRoute, newRoute)
// 					if err != nil {
// 						t.Error(err)
// 					}

// 					// fmt.Println("活動嘉賓資料: ", strings.Split(guest.Avatar, "/admin/uploads/")[1], strconv.Itoa(int(guest.ID)), activity.ActivityID)
// 					// 更新圖片資料
// 					DefaultGuestModel().SetDbConn(conn).
// 						Update(EditGuestModel{
// 							ID:         strconv.Itoa(int(guest.ID)),
// 							ActivityID: activity.ActivityID,
// 							Name:       guest.Name,
// 							Avatar:     strings.Split(guest.Avatar, "/admin/uploads/")[1],
// 							Introduce:  guest.Introduce,
// 							Detail:     guest.Detail,
// 							GuestOrder: "1",
// 						})
// 				}
// 			}

// 			// 查詢活動下的自定義資料
// 			// customize, err := DefaultCustomizeModel().SetDbConn(conn).
// 			// 	Find(activity.ActivityID)
// 			// if err != nil {
// 			// 	t.Error("錯誤: 取得活動下的定自訂欄位資料錯誤", err)
// 			// } else {
// 			// 	// fmt.Println("活動嘉賓數: ", len(guests))
// 			// }
// 			// if customize.InfoPicture != "" {
// 			// 	file := strings.Split(customize.InfoPicture, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/applysign/customize/info_picture." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	// fmt.Println("newRoute: ", newRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultCustomizeModel().SetDbConn(conn).
// 			// 		Update(EditCustomizeModel{
// 			// 			ActivityID:  activity.ActivityID,
// 			// 			InfoPicture: "info_picture." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 查詢活動下的QRcode自定義資料
// 			// activityModel, err := DefaultActivityModel().SetDbConn(conn).
// 			// 	Find(activity.ActivityID)
// 			// if err != nil {
// 			// 	t.Error("錯誤: 取得活動資料錯誤", err)
// 			// } else {
// 			// 	// fmt.Println("活動嘉賓數: ", len(guests))
// 			// }
// 			// if !strings.Contains(activityModel.QRcodeLogoPicture, "system") {
// 			// 	file := strings.Split(activityModel.QRcodeLogoPicture, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/applysign/qrcode/qrcode_logo_picture." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateQRcode(EditQRcodeModel{
// 			// 			ActivityID:        activity.ActivityID,
// 			// 			QRcodeLogoPicture: "qrcode_logo_picture." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 訊息區
// 			// if !strings.Contains(activityModel.MessageBackground, "system") {
// 			// 	file := strings.Split(activityModel.MessageBackground, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/wall/message/message_background." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateMessage(EditMessageModel{
// 			// 			ActivityID:        activity.ActivityID,
// 			// 			MessageBackground: "message_background." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 主題區
// 			// if !strings.Contains(activityModel.TopicBackground, "system") {
// 			// 	file := strings.Split(activityModel.TopicBackground, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/wall/topic/topic_background." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateTopic(EditTopicModel{
// 			// 			ActivityID:      activity.ActivityID,
// 			// 			TopicBackground: "topic_background." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 提問區
// 			// if !strings.Contains(activityModel.QuestionBackground, "system") {
// 			// 	file := strings.Split(activityModel.QuestionBackground, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/wall/question/question_background." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateQuestion(EditQuestionModel{
// 			// 			ActivityID:         activity.ActivityID,
// 			// 			QuestionBackground: "question_background." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 一般簽到
// 			// if !strings.Contains(activityModel.GeneralBackground, "system") {
// 			// 	file := strings.Split(activityModel.GeneralBackground, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/sign/general/general_background." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateGeneral(EditGeneralModel{
// 			// 			ActivityID:        activity.ActivityID,
// 			// 			GeneralBackground: "general_background." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 			// 3D簽到
// 			// if !strings.Contains(activityModel.ThreedAvatar, "system") {
// 			// 	file := strings.Split(activityModel.ThreedAvatar, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/sign/threed/threed_avatar." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateThreed(EditThreeDModel{
// 			// 			ActivityID:   activity.ActivityID,
// 			// 			ThreedAvatar: "threed_avatar." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }
// 			// if !strings.Contains(activityModel.ThreedBackground, "system") {
// 			// 	file := strings.Split(activityModel.ThreedBackground, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/sign/threed/threed_background." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateThreed(EditThreeDModel{
// 			// 			ActivityID:       activity.ActivityID,
// 			// 			ThreedBackground: "threed_background." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }
// 			// if !strings.Contains(activityModel.ThreedImageLogoPicture, "system") {
// 			// 	file := strings.Split(activityModel.ThreedImageLogoPicture, "/admin/uploads/")[1]
// 			// 	// fmt.Println("file: ", file)
// 			// 	// 原始圖片路徑
// 			// 	oldRoute := "../uploads/" + file
// 			// 	newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/sign/threed/threed_image_logo_picture." + strings.Split(file, ".")[1]
// 			// 	// fmt.Println("oldRoute: ", oldRoute)
// 			// 	err = copy(oldRoute, newRoute)
// 			// 	if err != nil {
// 			// 		t.Error(err)
// 			// 	}

// 			// 	// 更新圖片資料
// 			// 	DefaultActivityModel().SetDbConn(conn).
// 			// 		UpdateThreed(EditThreeDModel{
// 			// 			ActivityID:             activity.ActivityID,
// 			// 			ThreedImageLogoPicture: "threed_image_logo_picture." + strings.Split(file, ".")[1],
// 			// 		})
// 			// }

// 		}
// 	}
// }

// for _, game := range games {
// gameType, _ := game["game"].(string)
// gameID, _ := game["game_id"].(string)

// 遊戲資訊
// gameModel, _ := DefaultGameModel().SetDbConn(conn).Find(false, gameID)

// if gameType == "lottery" {
// 	fields := []string{
// 		// 遊戲抽獎自定義
// 		"lottery_screen_prizer",
// 		"lottery_screen_mascot",
// 		"lottery_screen_background",
// 		"lottery_screen_prize_notify",
// 		"lottery_screen_select_input",
// 		"lottery_screen_close_prize_notify_button",
// 		"lottery_player_background",
// 		"lottery_player_rules",
// 		"lottery_jiugongge_grid",
// 		"lottery_jiugongge_start_button",
// 		"lottery_turntable_start_button",
// 		"lottery_turntable_roulette",
// 		// 動圖
// 		"lottery_get_prize",
// 		"lottery_jiugongge_border",
// 		"lottery_jiugongge_title",
// 		"lottery_turntable_border",
// 		"lottery_turntable_title",
// 	}

// 	values := []string{
// 		// 遊戲抽獎自定義
// 		gameModel.LotteryScreenPrizer,
// 		gameModel.LotteryScreenMascot,
// 		gameModel.LotteryScreenBackground,
// 		gameModel.LotteryScreenPrizeNotify,
// 		gameModel.LotteryScreenSelectInput,
// 		gameModel.LotteryScreenClosePrizeNotifyButton,
// 		gameModel.LotteryPlayerBackground,
// 		gameModel.LotteryPlayerRules,
// 		gameModel.LotteryJiugonggeGrid,
// 		gameModel.LotteryJiugonggeStartButton,
// 		gameModel.LotteryTurntableStartButton,
// 		gameModel.LotteryTurntableRoulette,
// 		// 動圖
// 		gameModel.LotteryGetPrize,
// 		gameModel.LotteryJiugonggeBorder,
// 		gameModel.LotteryJiugonggeTitle,
// 		gameModel.LotteryTurntableBorder,
// 		gameModel.LotteryTurntableTitle,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("遊戲抽獎file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_2_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "redpack" {
// 	fields := []string{
// 		// 搖紅包自定義
// 		"redpack_screen_again_button",
// 		"redpack_screen_background",
// 		"redpack_screen_end_button",
// 		"redpack_screen_prize_list",
// 		"redpack_screen_prize_redpack",
// 		"redpack_screen_start_button",
// 		"redpack_screen_title",
// 		"redpack_screen_gaming_list",
// 		"redpack_screen_gaming_list_background",
// 		"redpack_screen_ema",
// 		"redpack_screen_new_list",
// 		"redpack_screen_lantern1",
// 		"redpack_screen_lantern2",
// 		"redpack_player_join_button",
// 		"redpack_player_search_prize_background",
// 		"redpack_player_background",
// 		"redpack_player_title",
// 		"redpack_player_lantern",
// 		// 動圖
// 		"redpack_screen_lucky_bag",
// 		"redpack_screen_money_piles",
// 		"redpack_player_shake",
// 		"redpack_player_lucky_bag",
// 		"redpack_player_money_piles",
// 		"redpack_screen_background_dynamic",
// 		"redpack_player_background_dynamic",
// 		"redpack_player_ready",
// 		// 音樂
// 		"redpack_bgm_start",
// 		"redpack_bgm_gaming",
// 		"redpack_bgm_end",
// 	}

// 	values := []string{
// 		// 搖紅包自定義
// 		gameModel.RedpackScreenAgainButton,
// 		gameModel.RedpackScreenBackground,
// 		gameModel.RedpackScreenEndButton,
// 		gameModel.RedpackScreenPrizeList,
// 		gameModel.RedpackScreenPrizeRedpack,
// 		gameModel.RedpackScreenStartButton,
// 		gameModel.RedpackScreenTitle,
// 		gameModel.RedpackScreenGamingList,
// 		gameModel.RedpackScreenGamingListBackground,
// 		gameModel.RedpackScreenEma,
// 		gameModel.RedpackScreenNewList,
// 		gameModel.RedpackScreenLantern1,
// 		gameModel.RedpackScreenLantern2,
// 		gameModel.RedpackPlayerJoinButton,
// 		gameModel.RedpackPlayerSearchPrizeBackground,
// 		gameModel.RedpackPlayerBackground,
// 		gameModel.RedpackPlayerTitle,
// 		gameModel.RedpackPlayerLantern,
// 		// 動圖
// 		gameModel.RedpackScreenLuckyBag,
// 		gameModel.RedpackScreenMoneyPiles,
// 		gameModel.RedpackPlayerShake,
// 		gameModel.RedpackPlayerLuckyBag,
// 		gameModel.RedpackPlayerMoneyPiles,
// 		gameModel.RedpackScreenBackgroundDynamic,
// 		gameModel.RedpackPlayerBackgroundDynamic,
// 		gameModel.RedpackPlayerReady,
// 		// 音樂
// 		gameModel.RedpackBgmStart,
// 		gameModel.RedpackBgmGaming,
// 		gameModel.RedpackBgmEnd,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_2_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}

// } else if gameType == "ropepack" {
// 	fields := []string{
// 		// 套紅包自定義
// 		"ropepack_screen_prize_list",
// 		"ropepack_screen_again_button",
// 		"ropepack_screen_background",
// 		"ropepack_screen_decoration",
// 		"ropepack_screen_end_button",
// 		"ropepack_screen_end_prize_list",
// 		"ropepack_screen_prize_redpack",
// 		"ropepack_screen_start_logo",
// 		"ropepack_screen_start_button",
// 		"ropepack_screen_prize_skin_red",
// 		"ropepack_screen_prize_skin_green",
// 		"ropepack_player_join_logo",
// 		"ropepack_player_join_button",
// 		"ropepack_player_background",
// 		"ropepack_player_ready_redpack1",
// 		"ropepack_player_ready_redpack2",
// 		"ropepack_player_ready_background",
// 		"ropepack_player_title",
// 		"ropepack_screen_background_foreground",     // 主持端背景前景
// 		"ropepack_screen_background_decoration",     // 主持端背景裝飾
// 		"ropepack_screen_gaming_prize_list",         // 主持端遊戲中中獎列表
// 		"ropepack_player_game_start",                // 玩家端遊戲即將開始
// 		"ropepack_screen_end_game_button",           // 主持端結束遊戲按鈕
// 		"ropepack_screen_end_prize_frame",           // 主持端結算得獎框
// 		"ropepack_screen_start_end_background",      // 主持端開始和結算背景
// 		"ropepack_screen_gaming_background",         // 主持端遊戲中背景
// 		"ropepack_screen_gaming_yellow_prize_frame", // 主持端遊戲中遊戲中黃得獎框
// 		"ropepack_screen_gaming_blue_prize_frame",   // 主持端遊戲中遊戲中藍得獎框
// 		"ropepack_game_title",                       // 遊戲名稱
// 		// 動圖
// 		"ropepack_screen_background_effect",
// 		"ropepack_player_ropepack_button",
// 		"ropepack_player_finger",
// 		"ropepack_redpack",
// 		"ropepack_screen_circle_effects",    // 主持端圓圈特效
// 		"ropepack_rabbit",                   // 兔子
// 		"ropepack_player_button",            // 玩家端按鈕
// 		"ropepack_player_ready_effects",     // 玩家端等待遊戲開始特效
// 		"ropepack_start_effects",            // 開始特效
// 		"ropepack_screen_start_end_effects", // 主持端開始和結算特效
// 		"ropepack_chang_e",                  // 嫦娥
// 		"ropepack_chang_e_silhouette",       // 嫦娥剪影
// 		// 音樂
// 		"ropepack_bgm_start",
// 		"ropepack_bgm_gaming",
// 		"ropepack_bgm_end",
// 	}

// 	values := []string{
// 		// 套紅包自定義
// 		gameModel.RopepackScreenPrizeList,
// 		gameModel.RopepackScreenAgainButton,
// 		gameModel.RopepackScreenBackground,
// 		gameModel.RopepackScreenDecoration,
// 		gameModel.RopepackScreenEndButton,
// 		gameModel.RopepackScreenEndPrizeList,
// 		gameModel.RopepackScreenPrizeRedpack,
// 		gameModel.RopepackScreenStartLogo,
// 		gameModel.RopepackScreenStartButton,
// 		gameModel.RopepackScreenPrizeSkinRed,
// 		gameModel.RopepackScreenPrizeSkinGreen,
// 		gameModel.RopepackPlayerJoinLogo,
// 		gameModel.RopepackPlayerJoinButton,
// 		gameModel.RopepackPlayerBackground,
// 		gameModel.RopepackPlayerReadyRedpack1,
// 		gameModel.RopepackPlayerReadyRedpack2,
// 		gameModel.RopepackPlayerReadyBackground,
// 		gameModel.RopepackPlayerTitle,
// 		gameModel.RopepackScreenBackgroundForeground,   // 主持端背景前景
// 		gameModel.RopepackScreenBackgroundDecoration,   // 主持端背景裝飾
// 		gameModel.RopepackScreenGamingPrizeList,        // 主持端遊戲中中獎列表
// 		gameModel.RopepackPlayerGameStart,              // 玩家端遊戲即將開始
// 		gameModel.RopepackScreenEndGameButton,          // 主持端結束遊戲按鈕
// 		gameModel.RopepackScreenEndPrizeFrame,          // 主持端結算得獎框
// 		gameModel.RopepackScreenStartEndBackground,     // 主持端開始和結算背景
// 		gameModel.RopepackScreenGamingBackground,       // 主持端遊戲中背景
// 		gameModel.RopepackScreenGamingYellowPrizeFrame, // 主持端遊戲中遊戲中黃得獎框
// 		gameModel.RopepackScreenGamingBluePrizeFrame,   // 主持端遊戲中遊戲中藍得獎框
// 		gameModel.RopepackGameTitle,                    // 遊戲名稱
// 		// 動圖
// 		gameModel.RopepackScreenBackgroundEffect,
// 		gameModel.RopepackPlayerRopepackButton,
// 		gameModel.RopepackPlayerFinger,
// 		gameModel.RopepackRedpack,
// 		gameModel.RopepackScreenCircleEffects,   // 主持端圓圈特效
// 		gameModel.RopepackRabbit,                // 兔子
// 		gameModel.RopepackPlayerButton,          // 玩家端按鈕
// 		gameModel.RopepackPlayerReadyEffects,    // 玩家端等待遊戲開始特效
// 		gameModel.RopepackStartEffects,          // 開始特效
// 		gameModel.RopepackScreenStartEndEffects, // 主持端開始和結算特效
// 		gameModel.RopepackChangE,                // 嫦娥
// 		gameModel.RopepackChangESilhouette,      // 嫦娥剪影
// 		// 音樂
// 		gameModel.RopepackBgmStart,
// 		gameModel.RopepackBgmGaming,
// 		gameModel.RopepackBgmEnd,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("套紅包file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_2_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "draw_numbers" {
// 	fields := []string{
// 		// 搖號抽獎自定義
// 		"draw_numbers_background",
// 		"draw_numbers_title",
// 		"draw_numbers_gift_inside_picture",
// 		"draw_numbers_gift_outside_picture",
// 		"draw_numbers_prize_left_button",
// 		"draw_numbers_prize_right_button",
// 		"draw_numbers_prize_leftright_button",
// 		"draw_numbers_addpeople_no_button",
// 		"draw_numbers_addpeople_yes_button",
// 		"draw_numbers_people_background",
// 		"draw_numbers_add_people",
// 		"draw_numbers_reduce_people",
// 		"draw_numbers_winner_background",
// 		"draw_numbers_blackground",
// 		"draw_numbers_go_button",
// 		"draw_numbers_open_winner_button",
// 		"draw_numbers_close_winner_button",
// 		"draw_numbers_current_people",
// 		// 動圖
// 		"draw_numbers_gacha_machine",
// 		"draw_numbers_hood",
// 		"draw_numbers_body",
// 		"draw_numbers_gacha",
// 		// 音樂
// 		"draw_numbers_bgm_gaming",
// 	}

// 	values := []string{
// 		// 搖號抽獎自定義
// 		gameModel.DrawNumbersBackground,
// 		gameModel.DrawNumbersTitle,
// 		gameModel.DrawNumbersGiftInsidePicture,
// 		gameModel.DrawNumbersGiftOutsidePicture,
// 		gameModel.DrawNumbersPrizeLeftButton,
// 		gameModel.DrawNumbersPrizeRightButton,
// 		gameModel.DrawNumbersPrizeLeftrightButton,
// 		gameModel.DrawNumbersAddpeopleNoNutton,
// 		gameModel.DrawNumbersAddpeopleYesButton,
// 		gameModel.DrawNumbersPeopleBackground,
// 		gameModel.DrawNumbersAddPeople,
// 		gameModel.DrawNumbersReducePeople,
// 		gameModel.DrawNumbersWinnerBackground,
// 		gameModel.DrawNumbersBlackground,
// 		gameModel.DrawNumbersGoButton,
// 		gameModel.DrawNumbersOpenWinnerButton,
// 		gameModel.DrawNumbersCloseWinnerButton,
// 		gameModel.DrawNumbersCurrentPeople,
// 		// 動圖
// 		gameModel.DrawNumbersGachaMachine,
// 		gameModel.DrawNumbersHood,
// 		gameModel.DrawNumbersBody,
// 		gameModel.DrawNumbersGacha,
// 		// 音樂
// 		gameModel.DrawNumbersBgmGaming,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_4_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "monopoly" {
// 	fields := []string{
// 		// 抓偽鈔自定義
// 		"monopoly_screen_again_button",                // 主持端再來一輪按鈕
// 		"monopoly_screen_top6_title",                  // 主持端前六名標題
// 		"monopoly_screen_end_button",                  // 主持端結束遊戲按鈕
// 		"monopoly_screen_start_button",                // 主持端開始遊戲按鈕
// 		"monopoly_screen_gaming_background_png",       // 主持端遊戲中背景
// 		"monopoly_screen_round_countdown",             // 主持端遊戲中輪次倒數
// 		"monopoly_screen_winner_list",                 // 主持端遊戲和結算中獎列表
// 		"monopoly_screen_rank_border",                 // 主持端遊戲和結算名次框
// 		"monopoly_player_carton",                      // 玩家端遊戲中下滑紙箱
// 		"monopoly_player_any_start_text",              // 玩家端按任意處開始文字
// 		"monopoly_player_scoreboard",                  // 玩家端計分和時間板
// 		"monopoly_player_wait_start_text",             // 玩家端計分和時間板
// 		"monopoly_player_transparent_background",      // 玩家端開始場景半透明黑底
// 		"monopoly_player_pile_objects",                // 玩家端滑動物件堆
// 		"monopoly_player_gaming_background",           // 玩家端遊戲中背景
// 		"monopoly_add_points",                         // 上滑小標示
// 		"monopoly_deduct_points",                      // 下滑小標示
// 		"monopoly_player_background_dynamic",          // 玩家端背景
// 		"monopoly_player_answer_effect",               // 玩家端遊戲中答對或錯特效
// 		"monopoly_background_and_gold",                // 主持端背景和玩家端金銅條
// 		"monopoly_screen_redpack_seal",                // 主持端紅包袋封口
// 		"monopoly_screen_again_button_background",     // 主持端再來一輪按鈕底
// 		"monopoly_screen_end_info_skin",               // 主持端結算3名後玩家資訊木框
// 		"monopoly_screen_end_npc",                     // 主持端結算吉祥物
// 		"monopoly_screen_top_stair",                   // 主持端結算前三名台階
// 		"monopoly_screen_top_info_skin",               // 主持端結算前三名資訊框
// 		"monopoly_screen_top_avatar_skin",             // 主持端結算前三名頭像框
// 		"monopoly_screen_end_background",              // 主持端結算背景
// 		"monopoly_screen_start_npc_dialog",            // 主持端開始畫面人物對話框
// 		"monopoly_screen_leaderboard",                 // 主持端遊戲中排行榜
// 		"monopoly_screen_round_background",            // 主持端遊戲中輪次底
// 		"monopoly_screen_start_end_button_background", // 主持端開始和結束按鈕底
// 		"monopoly_screen_start_background",            // 主持端開始背景
// 		"monopoly_screen_start_right_top_decoration",  // 主持端開始右上裝飾
// 		"monopoly_player_tip_arrow",                   // 玩家端遊戲中提示箭頭
// 		"monopoly_player_npc_dialog",                  // 玩家端人物對話框
// 		"monopoly_player_join_button_background",      // 玩家端加入遊戲按鈕底
// 		"monopoly_player_join_background",             // 玩家端加入遊戲背景
// 		"monopoly_player_redpack_space",               // 玩家端紅包袋白底
// 		"monopoly_player_redpack_seal",                // 玩家端紅包袋封口
// 		"monopoly_player_redpack_background",          // 玩家端紅包袋背景
// 		"monopoly_player_money_piles",                 // 玩家端鈔票堆
// 		"monopoly_player_background",                  // 玩家端遊戲背景
// 		"monopoly_player_title",                       // 玩家端遊戲標題
// 		"monopoly_npc",                                // 代表人物
// 		"monopoly_button",                             // 按鈕
// 		"monopoly_screen_top_light",                   // 主持端前三名發亮
// 		"monopoly_screen_end_revolving_light",         // 主持端結算背景旋轉燈
// 		"monopoly_screen_end_ribbon",                  // 主持端結算彩帶
// 		"monopoly_player_gaming_redpack",              // 玩家端遊戲中紅包
// 		"monopoly_screen_gaming_redpack",              // 主持端遊戲中紅包和玩家端紅包
// 		"monopoly_screen_top_after_player_info",       // 主持端結算3名後玩家資訊框
// 		"monopoly_screen_top_front_player_info",       // 主持端結算前三名資訊框
// 		"monopoly_screen_rank",                        // 主持端遊戲中排行榜
// 		"monopoly_screen_npc_dialog",                  // 主持端對話框
// 		"monopoly_screen_left_bottom_decoration",      // 主持端遊戲中裝飾小物件左下
// 		"monopoly_player_basket_background",           // 玩家端竹籃背景
// 		"monopoly_player_gaming_carrots",              // 玩家端遊戲中紅蘿蔔堆
// 		"monopoly_button_background",                  // 按鈕背景
// 		"monopoly_screen_end_background_dynamic",      // 主持端遊戲中和結算背景
// 		"monopoly_screen_start_background_dynamic",    // 主持端開始背景
// 		"monopoly_player_gaming_background_dynamic",   // 玩家端遊戲背景
// 		"monopoly_picking_carrots_and_carrots",        // 主持端遊戲中採蘿蔔和玩家端蘿蔔
// 		"monopoly_player_top_info",                    // 玩家端上方資訊
// 		"monopoly_player_search_prize_background",     // 玩家端查看獎品背景
// 		"monopoly_player_food_waste_bin",              // 玩家端遊戲中廚餘桶
// 		"monopoly_screen_end_dynamic",                 // 主持端結算動圖
// 		"monopoly_screen_timer",                       // 主持端遊戲中計時器
// 		"monopoly_player_start_gaming_eyecatch",       // 玩家端開始和遊戲過場
// 		"monopoly_gaming_dynamic_and_fish",            // 主持端遊戲中動圖和玩家端魚
// 		"monopoly_screen_gaming_background_jpg",       // 主持端遊戲中背景
// 		// 音樂
// 		"monopoly_bgm_start",
// 		"monopoly_bgm_gaming",
// 		"monopoly_bgm_end",
// 	}

// 	values := []string{
// 		// 抓偽鈔自定義
// 		gameModel.MonopolyScreenAgainButton,              // 主持端再來一輪按鈕
// 		gameModel.MonopolyScreenTop6Title,                // 主持端前六名標題
// 		gameModel.MonopolyScreenEndButton,                // 主持端結束遊戲按鈕
// 		gameModel.MonopolyScreenStartButton,              // 主持端開始遊戲按鈕
// 		gameModel.MonopolyScreenGamingBackgroundPng,      // 主持端遊戲中背景
// 		gameModel.MonopolyScreenRoundCountdown,           // 主持端遊戲中輪次倒數
// 		gameModel.MonopolyScreenWinnerList,               // 主持端遊戲和結算中獎列表
// 		gameModel.MonopolyScreenRankBorder,               // 主持端遊戲和結算名次框
// 		gameModel.MonopolyPlayerCarton,                   // 玩家端遊戲中下滑紙箱
// 		gameModel.MonopolyPlayerAnyStartText,             // 玩家端按任意處開始文字
// 		gameModel.MonopolyPlayerScoreboard,               // 玩家端計分和時間板
// 		gameModel.MonopolyPlayerWaitStartText,            // 玩家端計分和時間板
// 		gameModel.MonopolyPlayerTransparentBackground,    // 玩家端開始場景半透明黑底
// 		gameModel.MonopolyPlayerPileObjects,              // 玩家端滑動物件堆
// 		gameModel.MonopolyPlayerGamingBackground,         // 玩家端遊戲中背景
// 		gameModel.MonopolyAddPoints,                      // 上滑小標示
// 		gameModel.MonopolyDeductPoints,                   // 下滑小標示
// 		gameModel.MonopolyPlayerBackgroundDynamic,        // 玩家端背景
// 		gameModel.MonopolyPlayerAnswerEffect,             // 玩家端遊戲中答對或錯特效
// 		gameModel.MonopolyBackgroundAndGold,              // 主持端背景和玩家端金銅條
// 		gameModel.MonopolyScreenRedpackSeal,              // 主持端紅包袋封口
// 		gameModel.MonopolyScreenAgainButtonBackground,    // 主持端再來一輪按鈕底
// 		gameModel.MonopolyScreenEndInfoSkin,              // 主持端結算3名後玩家資訊木框
// 		gameModel.MonopolyScreenEndNpc,                   // 主持端結算吉祥物
// 		gameModel.MonopolyScreenTopStair,                 // 主持端結算前三名台階
// 		gameModel.MonopolyScreenTopInfoSkin,              // 主持端結算前三名資訊框
// 		gameModel.MonopolyScreenTopAvatarSkin,            // 主持端結算前三名頭像框
// 		gameModel.MonopolyScreenEndBackground,            // 主持端結算背景
// 		gameModel.MonopolyScreenStartNpcDialog,           // 主持端開始畫面人物對話框
// 		gameModel.MonopolyScreenLeaderboard,              // 主持端遊戲中排行榜
// 		gameModel.MonopolyScreenRoundBackground,          // 主持端遊戲中輪次底
// 		gameModel.MonopolyScreenStartEndButtonBackground, // 主持端開始和結束按鈕底
// 		gameModel.MonopolyScreenStartBackground,          // 主持端開始背景
// 		gameModel.MonopolyScreenStartRightTopDecoration,  // 主持端開始右上裝飾
// 		gameModel.MonopolyPlayerTipArrow,                 // 玩家端遊戲中提示箭頭
// 		gameModel.MonopolyPlayerNpcDialog,                // 玩家端人物對話框
// 		gameModel.MonopolyPlayerJoinButtonBackground,     // 玩家端加入遊戲按鈕底
// 		gameModel.MonopolyPlayerJoinBackground,           // 玩家端加入遊戲背景
// 		gameModel.MonopolyPlayerRedpackSpace,             // 玩家端紅包袋白底
// 		gameModel.MonopolyPlayerRedpackSeal,              // 玩家端紅包袋封口
// 		gameModel.MonopolyPlayerRedpackBackground,        // 玩家端紅包袋背景
// 		gameModel.MonopolyPlayerMoneyPiles,               // 玩家端鈔票堆
// 		gameModel.MonopolyPlayerBackground,               // 玩家端遊戲背景
// 		gameModel.MonopolyPlayerTitle,                    // 玩家端遊戲標題
// 		gameModel.MonopolyNpc,                            // 代表人物
// 		gameModel.MonopolyButton,                         // 按鈕
// 		gameModel.MonopolyScreenTopLight,                 // 主持端前三名發亮
// 		gameModel.MonopolyScreenEndRevolvingLight,        // 主持端結算背景旋轉燈
// 		gameModel.MonopolyScreenEndRibbon,                // 主持端結算彩帶
// 		gameModel.MonopolyPlayerGamingRedpack,            // 玩家端遊戲中紅包
// 		gameModel.MonopolyScreenGamingRedpack,            // 主持端遊戲中紅包和玩家端紅包
// 		gameModel.MonopolyScreenTopAfterPlayerInfo,       // 主持端結算3名後玩家資訊框
// 		gameModel.MonopolyScreenTopFrontPlayerInfo,       // 主持端結算前三名資訊框
// 		gameModel.MonopolyScreenRank,                     // 主持端遊戲中排行榜
// 		gameModel.MonopolyScreenNpcDialog,                // 主持端對話框
// 		gameModel.MonopolyScreenLeftBottomDecoration,     // 主持端遊戲中裝飾小物件左下
// 		gameModel.MonopolyPlayerBasketBackground,         // 玩家端竹籃背景
// 		gameModel.MonopolyPlayerGamingCarrots,            // 玩家端遊戲中紅蘿蔔堆
// 		gameModel.MonopolyButtonBackground,               // 按鈕背景
// 		gameModel.MonopolyScreenEndBackgroundDynamic,     // 主持端遊戲中和結算背景
// 		gameModel.MonopolyScreenStartBackgroundDynamic,   // 主持端開始背景
// 		gameModel.MonopolyPlayerGamingBackgroundDynamic,  // 玩家端遊戲背景
// 		gameModel.MonopolyPickingCarrotsAndCarrots,       // 主持端遊戲中採蘿蔔和玩家端蘿蔔
// 		gameModel.MonopolyPlayerTopInfo,                  // 玩家端上方資訊
// 		gameModel.MonopolyPlayerSearchPrizeBackground,    // 玩家端查看獎品背景
// 		gameModel.MonopolyPlayerFoodWasteBin,             // 玩家端遊戲中廚餘桶
// 		gameModel.MonopolyScreenEndDynamic,               // 主持端結算動圖
// 		gameModel.MonopolyScreenTimer,                    // 主持端遊戲中計時器
// 		gameModel.MonopolyPlayerStartGamingEyecatch,      // 玩家端開始和遊戲過場
// 		gameModel.MonopolyGamingDynamicAndFish,           // 主持端遊戲中動圖和玩家端魚
// 		gameModel.MonopolyScreenGamingBackgroundJpg,      // 主持端遊戲中背景
// 		// 音樂
// 		gameModel.MonopolyBgmStart,
// 		gameModel.MonopolyBgmGaming,
// 		gameModel.MonopolyBgmEnd,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_3_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "whack_mole" {
// 	fields := []string{
// 		// 敲敲樂自定義
// 		"whackmole_classic_01",
// 		"whackmole_classic_02",
// 		"whackmole_classic_03",
// 		"whackmole_classic_04",
// 		"whackmole_classic_05",
// 		"whackmole_classic_06",
// 		"whackmole_classic_07",
// 		"whackmole_classic_08",
// 		"whackmole_classic_09",
// 		"whackmole_classic_10",
// 		"whackmole_classic_11",
// 		"whackmole_classic_12",
// 		"whackmole_classic_13",
// 		"whackmole_classic_14",
// 		"whackmole_classic_15",
// 		"whackmole_classic_16",
// 		"whackmole_classic_17",
// 		"whackmole_classic_18",
// 		"whackmole_classic_19",
// 		"whackmole_classic_20",
// 		"whackmole_classic_21",
// 		"whackmole_classic_22",
// 		"whackmole_classic_23",
// 		"whackmole_classic_24",
// 		"whackmole_classic_25",
// 		"whackmole_classic_26",
// 		"whackmole_classic_27",
// 		"whackmole_classic_28",
// 		"whackmole_classic_29",
// 		"whackmole_classic_30",

// 		"whackmole_halloween_01",
// 		"whackmole_halloween_02",
// 		"whackmole_halloween_03",
// 		"whackmole_halloween_04",
// 		"whackmole_halloween_05",
// 		"whackmole_halloween_06",
// 		"whackmole_halloween_07",
// 		"whackmole_halloween_08",
// 		"whackmole_halloween_09",
// 		"whackmole_halloween_10",
// 		"whackmole_halloween_11",
// 		"whackmole_halloween_12",
// 		"whackmole_halloween_13",
// 		"whackmole_halloween_14",
// 		"whackmole_halloween_15",
// 		"whackmole_halloween_16",
// 		"whackmole_halloween_17",
// 		"whackmole_halloween_18",
// 		"whackmole_halloween_19",
// 		"whackmole_halloween_20",
// 		"whackmole_halloween_21",
// 		"whackmole_halloween_22",
// 		"whackmole_halloween_23",
// 		"whackmole_halloween_24",
// 		"whackmole_halloween_25",
// 		"whackmole_halloween_26",
// 		"whackmole_halloween_27",
// 		"whackmole_halloween_28",
// 		"whackmole_halloween_29",
// 		"whackmole_halloween_30",

// 		"whackmole_christmas_01",
// 		"whackmole_christmas_02",
// 		"whackmole_christmas_03",
// 		"whackmole_christmas_04",
// 		"whackmole_christmas_05",
// 		"whackmole_christmas_06",
// 		"whackmole_christmas_07",
// 		"whackmole_christmas_08",
// 		"whackmole_christmas_09",
// 		"whackmole_christmas_10",
// 		"whackmole_christmas_11",
// 		"whackmole_christmas_12",
// 		"whackmole_christmas_13",
// 		"whackmole_christmas_14",
// 		"whackmole_christmas_15",
// 		"whackmole_christmas_16",
// 		"whackmole_christmas_17",
// 		"whackmole_christmas_18",
// 		"whackmole_christmas_19",
// 		"whackmole_christmas_20",
// 		"whackmole_christmas_21",
// 		"whackmole_christmas_22",
// 		"whackmole_christmas_23",
// 		"whackmole_christmas_24",
// 		"whackmole_christmas_25",
// 		"whackmole_christmas_26",
// 		"whackmole_christmas_27",
// 		"whackmole_christmas_28",
// 		"whackmole_christmas_29",
// 		"whackmole_christmas_30",
// 		"whackmole_christmas_31",

// 		// 動圖
// 		"whackmole_classic_31",
// 		"whackmole_halloween_31",
// 		"whackmole_christmas_32",
// 		"whackmole_christmas_33",

// 		// 音樂
// 		"whackmole_bgm_start",
// 		"whackmole_bgm_gaming",
// 		"whackmole_bgm_end",
// 	}

// 	values := []string{
// 		// 敲敲樂自定義
// 		gameModel.WhackmoleClassic01,
// 		gameModel.WhackmoleClassic02,
// 		gameModel.WhackmoleClassic03,
// 		gameModel.WhackmoleClassic04,
// 		gameModel.WhackmoleClassic05,
// 		gameModel.WhackmoleClassic06,
// 		gameModel.WhackmoleClassic07,
// 		gameModel.WhackmoleClassic08,
// 		gameModel.WhackmoleClassic09,
// 		gameModel.WhackmoleClassic10,
// 		gameModel.WhackmoleClassic11,
// 		gameModel.WhackmoleClassic12,
// 		gameModel.WhackmoleClassic13,
// 		gameModel.WhackmoleClassic14,
// 		gameModel.WhackmoleClassic15,
// 		gameModel.WhackmoleClassic16,
// 		gameModel.WhackmoleClassic17,
// 		gameModel.WhackmoleClassic18,
// 		gameModel.WhackmoleClassic19,
// 		gameModel.WhackmoleClassic20,
// 		gameModel.WhackmoleClassic21,
// 		gameModel.WhackmoleClassic22,
// 		gameModel.WhackmoleClassic23,
// 		gameModel.WhackmoleClassic24,
// 		gameModel.WhackmoleClassic25,
// 		gameModel.WhackmoleClassic26,
// 		gameModel.WhackmoleClassic27,
// 		gameModel.WhackmoleClassic28,
// 		gameModel.WhackmoleClassic29,
// 		gameModel.WhackmoleClassic30,

// 		gameModel.WhackmoleHalloween01,
// 		gameModel.WhackmoleHalloween02,
// 		gameModel.WhackmoleHalloween03,
// 		gameModel.WhackmoleHalloween04,
// 		gameModel.WhackmoleHalloween05,
// 		gameModel.WhackmoleHalloween06,
// 		gameModel.WhackmoleHalloween07,
// 		gameModel.WhackmoleHalloween08,
// 		gameModel.WhackmoleHalloween09,
// 		gameModel.WhackmoleHalloween10,
// 		gameModel.WhackmoleHalloween11,
// 		gameModel.WhackmoleHalloween12,
// 		gameModel.WhackmoleHalloween13,
// 		gameModel.WhackmoleHalloween14,
// 		gameModel.WhackmoleHalloween15,
// 		gameModel.WhackmoleHalloween16,
// 		gameModel.WhackmoleHalloween17,
// 		gameModel.WhackmoleHalloween18,
// 		gameModel.WhackmoleHalloween19,
// 		gameModel.WhackmoleHalloween20,
// 		gameModel.WhackmoleHalloween21,
// 		gameModel.WhackmoleHalloween22,
// 		gameModel.WhackmoleHalloween23,
// 		gameModel.WhackmoleHalloween24,
// 		gameModel.WhackmoleHalloween25,
// 		gameModel.WhackmoleHalloween26,
// 		gameModel.WhackmoleHalloween27,
// 		gameModel.WhackmoleHalloween28,
// 		gameModel.WhackmoleHalloween29,
// 		gameModel.WhackmoleHalloween30,

// 		gameModel.WhackmoleChristmas01,
// 		gameModel.WhackmoleChristmas02,
// 		gameModel.WhackmoleChristmas03,
// 		gameModel.WhackmoleChristmas04,
// 		gameModel.WhackmoleChristmas05,
// 		gameModel.WhackmoleChristmas06,
// 		gameModel.WhackmoleChristmas07,
// 		gameModel.WhackmoleChristmas08,
// 		gameModel.WhackmoleChristmas09,
// 		gameModel.WhackmoleChristmas10,
// 		gameModel.WhackmoleChristmas11,
// 		gameModel.WhackmoleChristmas12,
// 		gameModel.WhackmoleChristmas13,
// 		gameModel.WhackmoleChristmas14,
// 		gameModel.WhackmoleChristmas15,
// 		gameModel.WhackmoleChristmas16,
// 		gameModel.WhackmoleChristmas17,
// 		gameModel.WhackmoleChristmas18,
// 		gameModel.WhackmoleChristmas19,
// 		gameModel.WhackmoleChristmas20,
// 		gameModel.WhackmoleChristmas21,
// 		gameModel.WhackmoleChristmas22,
// 		gameModel.WhackmoleChristmas23,
// 		gameModel.WhackmoleChristmas24,
// 		gameModel.WhackmoleChristmas25,
// 		gameModel.WhackmoleChristmas26,
// 		gameModel.WhackmoleChristmas27,
// 		gameModel.WhackmoleChristmas28,
// 		gameModel.WhackmoleChristmas29,
// 		gameModel.WhackmoleChristmas30,
// 		gameModel.WhackmoleChristmas31,

// 		// 動圖
// 		gameModel.WhackmoleClassic31,
// 		gameModel.WhackmoleHalloween31,
// 		gameModel.WhackmoleChristmas32,
// 		gameModel.WhackmoleChristmas33,

// 		// 音樂
// 		gameModel.WhackmoleBgmStart,
// 		gameModel.WhackmoleBgmGaming,
// 		gameModel.WhackmoleBgmEnd,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_7_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "QA" {
// 	fields4 := []string{ // 快問快答自定義
// 		"qa_mascot",
// 		"qa_host_start_background",
// 		"qa_host_game_background",
// 		"qa_host_end_background",
// 		"qa_game_top_1",
// 		"qa_game_top_2",
// 		"qa_game_top_3",
// 		"qa_game_top_4",
// 		"qa_game_top_5",
// 		"qa_end_top_1",
// 		"qa_end_top_2",
// 		"qa_end_top_3",
// 		"qa_end_top",
// 		"qa_host_start_game_button",
// 		"qa_host_pause_countdown_button",
// 		"qa_host_continue_countdown_button",
// 		"qa_host_start_answer_button",
// 		"qa_host_see_answer_button",
// 		"qa_host_next_question_button",
// 		"qa_host_end_game_button",
// 		"qa_host_again_game_button",
// 		"qa_player_start_background",
// 		"qa_player_game_background",
// 		"qa_player_join_game_button",
// 		"qa_player_select_answer_button",
// 		"qa_player_confirm_answer_button",
// 		"qa_player_confirm_status_button",
// 		"qa_host_countdown_left_bottom_frame",  // 主持端左下倒數框
// 		"qa_host_background",                   // 主持端背景
// 		"qa_host_countdown_small_frame",        // 主持端倒數數字小框
// 		"qa_host_question_frame",               // 主持端提問中提問框
// 		"qa_host_end_rank_frame",               // 主持端結算排行榜框
// 		"qa_host_end_rank_frame_sign",          // 主持端結算排行榜標示
// 		"qa_answer_bar",                        // 主持端遊戲中答題人數條
// 		"qa_correct_answer_bar",                // 主持端遊戲中答題正確人數條
// 		"qa_pause_countdown_button",            // 主持端暫停倒數按鈕
// 		"qa_continue_countdown_button",         // 主持端繼續倒數按鈕
// 		"qa_player_background",                 // 玩家端背景
// 		"qa_player_small_background",           // 玩家端遊戲中答題小背景
// 		"qa_title",                             // 標題
// 		"qa_screen_gaming_question_number_bar", // 主持端遊戲中答題人數條
// 		"qa_screen_gaming_answer_number_bar",   // 主持端遊戲中答題正確人數條

// 		// 動圖
// 		"qa_host_effect",           // 主持端特效
// 		"qa_host_options",          // 主持端選項
// 		"qa_player_current_option", // 玩家端目前選擇的選項
// 		"qa_player_effect",         // 玩家端特效
// 		"qa_player_options",        // 玩家端選項
// 		"qa_screen_end_light",      // 主持端結算燈光
// 		"qa_screen_option",         // 主持端選項
// 		"qa_player_option",         // 玩家端選項

// 		// 音樂
// 		"qa_bgm_start",  // 遊戲開始
// 		"qa_bgm_gaming", // 遊戲進行中
// 		"qa_bgm_end",    // 遊戲結束
// 	}

// 	values4 := []string{
// 		// 快問快答自定義
// 		gameModel.QAMascot,
// 		gameModel.QAHostStartBackground,
// 		gameModel.QAHostGameBackground,
// 		gameModel.QAhostEndBackground,
// 		gameModel.QAGameTop1,
// 		gameModel.QAGameTop2,
// 		gameModel.QAGameTop3,
// 		gameModel.QAGameTop4,
// 		gameModel.QAGameTop5,
// 		gameModel.QAEndTop1,
// 		gameModel.QAEndTop2,
// 		gameModel.QAEndTop3,
// 		gameModel.QAEndTop,
// 		gameModel.QAHostStartGameButton,
// 		gameModel.QAHostPauseCountdownButton,
// 		gameModel.QAHostContinueCountdownButton,
// 		gameModel.QAHostStartAnswerButton,
// 		gameModel.QAHostSeeAnswerButton,
// 		gameModel.QAHostNextQuestionButton,
// 		gameModel.QAHostEndGameButton,
// 		gameModel.QAHostAgainGameButton,
// 		gameModel.QAPlayerStartBackground,
// 		gameModel.QAPlayerGameBackground,
// 		gameModel.QAPlayerJoinGameButton,
// 		gameModel.QAPlayerSelectAnswerButton,
// 		gameModel.QAPlayerConfirmAnswerButton,
// 		gameModel.QAPlayerConfirmStatusButton,
// 		gameModel.QAHostCountdownLeftBottomFrame,  // 主持端左下倒數框
// 		gameModel.QAHostBackground,                // 主持端背景
// 		gameModel.QAHostCountdownSmallFrame,       // 主持端倒數數字小框
// 		gameModel.QAHostQuestionFrame,             // 主持端提問中提問框
// 		gameModel.QAHostEndRankFrame,              // 主持端結算排行榜框
// 		gameModel.QAHostEndRankFrameSign,          // 主持端結算排行榜標示
// 		gameModel.QAAnswerBar,                     // 主持端遊戲中答題人數條
// 		gameModel.QACorrectAnswerBar,              // 主持端遊戲中答題正確人數條
// 		gameModel.QAPauseCountdownButton,          // 主持端暫停倒數按鈕
// 		gameModel.QAContinueCountdownButton,       // 主持端繼續倒數按鈕
// 		gameModel.QAPlayerBackground,              // 玩家端背景
// 		gameModel.QAPlayerSmallBackground,         // 玩家端遊戲中答題小背景
// 		gameModel.QATitle,                         // 標題
// 		gameModel.QAScreenGamingQuestionNumberBar, // 主持端遊戲中答題人數條
// 		gameModel.QAScreenGamingAnswerNumberBar,   // 主持端遊戲中答題正確人數條

// 		// 動圖
// 		gameModel.QAHostEffect,          // 主持端特效
// 		gameModel.QAHostOptions,         // 主持端選項
// 		gameModel.QAPlayerCurrentOption, // 玩家端目前選擇的選項
// 		gameModel.QAPlayerEffect,        // 玩家端特效
// 		gameModel.QAPlayerOptions,       // 玩家端選項
// 		gameModel.QAScreenEndLight,      // 主持端結算燈光
// 		gameModel.QAScreenOption,        // 主持端選項
// 		gameModel.QAPlayerOption,        // 玩家端選項

// 		// 音樂
// 		gameModel.QABgmStart,  // 遊戲開始
// 		gameModel.QABgmGaming, // 遊戲進行中
// 		gameModel.QABgmEnd,    // 遊戲結束
// 	}

// 	for i, value4 := range values4 {
// 		if !strings.Contains(value4, "system") {
// 			file := strings.Split(value4, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields4[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_4_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields4[i]: fields4[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}

// 	fields6 := []string{ // 快問快答自定義
// 		// 快問快答自定義
// 		"qa_screen_next_question_button",                // 主持端下一題按鈕
// 		"qa_screen_another_round_button",                // 主持端再來一輪按鈕
// 		"qa_screen_another_round_button_hover",          // 主持端再來一輪按鈕懸停
// 		"qa_screen_question_end_background",             // 主持端提問和結算背景
// 		"qa_screen_answer_reveal_button",                // 主持端答案揭曉按鈕
// 		"qa_screen_end_winner_list",                     // 主持端結算中獎列表
// 		"qa_screen_end_button",                          // 主持端結算按鈕
// 		"qa_screen_end_box_after_fourth",                // 主持端結算框四名後
// 		"qa_screen_end_first_place_box",                 // 主持端結算第一名框
// 		"qa_screen_end_first_place_icon",                // 主持端結算第一名圖標
// 		"qa_screen_end_first_place_label",               // 主持端結算第一名標示
// 		"qa_screen_end_second_place_box",                // 主持端結算第二名框
// 		"qa_screen_end_second_place_icon",               // 主持端結算第二名圖標
// 		"qa_screen_end_second_place_label",              // 主持端結算第二名標示
// 		"qa_screen_end_third_place_box",                 // 主持端結算第三名框
// 		"qa_screen_end_third_place_icon",                // 主持端結算第三名圖標
// 		"qa_screen_end_third_place_label",               // 主持端結算第三名標示
// 		"qa_screen_start_answer_button",                 // 主持端開始作答按鈕
// 		"qa_screen_start_answer_button_hover",           // 主持端開始作答按鈕懸停
// 		"qa_screen_start_background",                    // 主持端開始背景
// 		"qa_screen_start_background_normal_light",       // 主持端開始背景常態光
// 		"qa_screen_start_background_mask",               // 主持端開始背景遮罩
// 		"qa_screen_start_game_button",                   // 主持端開始遊戲按鈕
// 		"qa_screen_start_game_hover",                    // 主持端開始遊戲懸停
// 		"qa_screen_gaming_background",                   // 主持端遊戲中背景
// 		"qa_screen_pause_countdown_button",              // 主持端暫停倒數按鈕
// 		"qa_screen_pause_countdown_button_hover",        // 主持端暫停倒數按鈕懸停
// 		"qa_screen_continue_countdown_button",           // 主持端繼續倒數按鈕
// 		"qa_screen_continue_countdown_hover",            // 主持端繼續倒數懸停
// 		"qa_player_join_game_screen_start_light_source", // 玩家端加入遊戲和主持端開始遊戲光源
// 		"qa_player_join_game_button_pressed",            // 玩家端加入遊戲按鈕按下
// 		"qa_player_answer_confirmed_button",             // 玩家端答案已確認按鈕
// 		"qa_player_answer_confirm_button",               // 玩家端答案確認按鈕
// 		"qa_player_start_background_mask",               // 玩家端開始背景遮罩
// 		"qa_player_gaming_background",                   // 玩家端遊戲中背景
// 		"qa_player_gaming_mask",                         // 玩家端遊戲中遮罩
// 		"qa_player_confirm_answer_halo_weak",            // 玩家端確認答案光暈弱
// 		"qa_player_confirm_answer_halo_strong",          // 玩家端確認答案光暈強

// 		// 動圖
// 		"qa_screen_question_box",       // 主持端提問框
// 		"qa_screen_start_bee",          // 主持端開始場景蜜蜂
// 		"qa_screen_gaming_ranking",     // 主持端遊戲中排名
// 		"qa_game_name_score_end_title", // 遊戲名稱和得分結算標題
// 	}

// 	values6 := []string{
// 		// 快問快答自定義
// 		gameModel.QAScreenNextQuestionButton,             // 主持端下一題按鈕
// 		gameModel.QAScreenAnotherRoundButton,             // 主持端再來一輪按鈕
// 		gameModel.QAScreenAnotherRoundButtonHover,        // 主持端再來一輪按鈕懸停
// 		gameModel.QAScreenQuestionEndBackground,          // 主持端提問和結算背景
// 		gameModel.QAScreenAnswerRevealButton,             // 主持端答案揭曉按鈕
// 		gameModel.QAScreenEndWinnerList,                  // 主持端結算中獎列表
// 		gameModel.QAScreenEndButton,                      // 主持端結算按鈕
// 		gameModel.QAScreenEndBoxAfterFourth,              // 主持端結算框四名後
// 		gameModel.QAScreenEndFirstPlaceBox,               // 主持端結算第一名框
// 		gameModel.QAScreenEndFirstPlaceIcon,              // 主持端結算第一名圖標
// 		gameModel.QAScreenEndFirstPlaceLabel,             // 主持端結算第一名標示
// 		gameModel.QAScreenEndSecondPlaceBox,              // 主持端結算第二名框
// 		gameModel.QAScreenEndSecondPlaceIcon,             // 主持端結算第二名圖標
// 		gameModel.QAScreenEndSecondPlaceLabel,            // 主持端結算第二名標示
// 		gameModel.QAScreenEndThirdPlaceBox,               // 主持端結算第三名框
// 		gameModel.QAScreenEndThirdPlaceIcon,              // 主持端結算第三名圖標
// 		gameModel.QAScreenEndThirdPlaceLabel,             // 主持端結算第三名標示
// 		gameModel.QAScreenStartAnswerButton,              // 主持端開始作答按鈕
// 		gameModel.QAScreenStartAnswerButtonHover,         // 主持端開始作答按鈕懸停
// 		gameModel.QAScreenStartBackground,                // 主持端開始背景
// 		gameModel.QAScreenStartBackgroundNormalLight,     // 主持端開始背景常態光
// 		gameModel.QAScreenStartBackgroundMask,            // 主持端開始背景遮罩
// 		gameModel.QAScreenStartGameButton,                // 主持端開始遊戲按鈕
// 		gameModel.QAScreenStartGameHover,                 // 主持端開始遊戲懸停
// 		gameModel.QAScreenGamingBackground,               // 主持端遊戲中背景
// 		gameModel.QAScreenPauseCountdownButton,           // 主持端暫停倒數按鈕
// 		gameModel.QAScreenPauseCountdownButtonHover,      // 主持端暫停倒數按鈕懸停
// 		gameModel.QAScreenContinueCountdownButton,        // 主持端繼續倒數按鈕
// 		gameModel.QAScreenContinueCountdownHover,         // 主持端繼續倒數懸停
// 		gameModel.QAPlayerJoinGameScreenStartLightSource, // 玩家端加入遊戲和主持端開始遊戲光源
// 		gameModel.QAPlayerJoinGameButtonPressed,          // 玩家端加入遊戲按鈕按下
// 		gameModel.QAPlayerAnswerConfirmedButton,          // 玩家端答案已確認按鈕
// 		gameModel.QAPlayerAnswerConfirmButton,            // 玩家端答案確認按鈕
// 		gameModel.QAPlayerStartBackgroundMask,            // 玩家端開始背景遮罩
// 		gameModel.QAPlayerGamingBackground,               // 玩家端遊戲中背景
// 		gameModel.QAPlayerGamingMask,                     // 玩家端遊戲中遮罩
// 		gameModel.QAPlayerConfirmAnswerHaloWeak,          // 玩家端確認答案光暈弱
// 		gameModel.QAPlayerConfirmAnswerHaloStrong,        // 玩家端確認答案光暈強

// 		// 動圖
// 		gameModel.QAScreenQuestionBox,     // 主持端提問框
// 		gameModel.QAScreenStartBee,        // 主持端開始場景蜜蜂
// 		gameModel.QAScreenGamingRanking,   // 主持端遊戲中排名
// 		gameModel.QAGameNameScoreEndTitle, // 遊戲名稱和得分結算標題
// 	}

// 	for i, value6 := range values6 {
// 		if !strings.Contains(value6, "system") {
// 			file := strings.Split(value6, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields6[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_6_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields6[i]: fields6[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "tugofwar" {
// 	fields := []string{
// 		// 拔河遊戲自定義
// 		"tugofwar_screen_score_and_countdown_info_column",               // 主持端分數和倒數資訊欄
// 		"tugofwar_screen_next_round_button",                             // 主持端再來一輪按鈕
// 		"tugofwar_screen_background",                                    // 主持端背景
// 		"tugofwar_screen_winning_MVP_label",                             // 主持端勝方MVP標示
// 		"tugofwar_screen_winning_captain_label",                         // 主持端勝方隊長標示
// 		"tugofwar_screen_end_game_button",                               // 主持端結束遊戲按鈕
// 		"tugofwar_screen_end_MVP_winner_avatar_frame",                   // 主持端結算MVP勝方頭像框
// 		"tugofwar_screen_end_background",                                // 主持端結算背景
// 		"tugofwar_screen_end_losing_MVP_label",                          // 主持端結算敗方MVP標示
// 		"tugofwar_screen_end_losing_prize_list",                         // 主持端結算敗方中獎列表
// 		"tugofwar_screen_end_losing_score_frame",                        // 主持端結算敗方分數框
// 		"tugofwar_screen_end_losing_player_info_frame",                  // 主持端結算敗方玩家資訊框
// 		"tugofwar_screen_end_losing_captain_label",                      // 主持端結算敗方隊長標示
// 		"tugofwar_screen_end_winning_prize_list",                        // 主持端結算勝方中獎列表
// 		"tugofwar_screen_end_winning_score_frame",                       // 主持端結算勝方分數框
// 		"tugofwar_screen_end_winning_player_info_frame",                 // 主持端結算勝方玩家資訊框
// 		"tugofwar_screen_end_captain_MVP_nameplate_score_losing_frame",  // 主持端結算隊長和MVP名牌分數敗方框
// 		"tugofwar_screen_end_captain_MVP_nameplate_score_winning_frame", // 主持端結算隊長和MVP名牌分數勝方框
// 		"tugofwar_screen_end_captain_MVP_losing_avatar_frame",           // 主持端結算隊長和MVP敗方頭像框
// 		"tugofwar_screen_end_captain_winning_avatar_frame",              // 主持端結算隊長勝方頭像框
// 		"tugofwar_screen_start_scene_adjust_number_background",          // 主持端開始場景調整人數背景
// 		"tugofwar_screen_start_game_button",                             // 主持端開始遊戲按鈕
// 		"tugofwar_screen_current_participants_frame",                    // 主持端當前參加人數框
// 		"tugofwar_screen_gaming_captain_avatar_frame",                   // 主持端遊戲中隊長頭像框
// 		"tugofwar_screen_game_captain_name_frame",                       // 主持端遊戲隊長名稱框
// 		"tugofwar_screen_adjust_number_button",                          // 主持端調整人數按鈕
// 		"tugofwar_right_team_name_frame",                                // 右隊名稱框
// 		"tugofwar_left_team_name_frame",                                 // 左隊名稱框
// 		"tugofwar_player_join_faction_button",                           // 玩家端加入陣營按鈕
// 		"tugofwar_player_background",                                    // 玩家端背景
// 		"tugofwar_player_tutorial_diagram_example",                      // 玩家端教學示意圖示
// 		"tugofwar_player_captain_hint_frame",                            // 玩家端隊長提示框
// 		"tugofwar_player_gaming_score_frame",                            // 玩家端遊戲中分數框
// 		"tugofwar_player_game_about_to_start",                           // 玩家端遊戲即將開始
// 		"tugofwar_player_title",                                         // 玩家端標題
// 		"tugofwar_team_logo_picture_frame",                              // 隊伍LOGO圖片框
// 		"tugofwar_team_number_of_people_frame",                          // 隊伍人數框
// 		"tugofwar_game_start_button",                                    // 加入遊戲按鈕

// 		// 動圖
// 		"tugofwar_screen_end_winning_confetti",    // 主持端結算勝方彩帶
// 		"tugofwar_screen_end_winning_fireworks",   // 主持端結算勝方煙火
// 		"tugofwar_screen_end_win_loss_label",      // 主持端結算勝敗標示
// 		"tugofwar_screen_gaming_sign_holding_bee", // 主持端遊戲中舉牌蜜蜂
// 		"tugofwar_character_sportsman",            // 拔河人物-運動男
// 		"tugofwar_character_bee",                  // 拔河人物-蜜蜂
// 		"tugofwar_middle_ribbon",                  // 拔河中間緞帶

// 		// 音樂
// 		"tugofwar_bgm_start",  // 遊戲開始
// 		"tugofwar_bgm_gaming", // 遊戲進行中
// 		"tugofwar_bgm_end",    // 遊戲結束
// 	}

// 	values := []string{
// 		// 拔河遊戲自定義
// 		gameModel.TugofwarScreenScoreAndCountdownInfoColumn,             // 主持端分數和倒數資訊欄
// 		gameModel.TugofwarScreenNextRoundButton,                         // 主持端再來一輪按鈕
// 		gameModel.TugofwarScreenBackground,                              // 主持端背景
// 		gameModel.TugofwarScreenWinningMVPLabel,                         // 主持端勝方MVP標示
// 		gameModel.TugofwarScreenWinningCaptainLabel,                     // 主持端勝方隊長標示
// 		gameModel.TugofwarScreenEndGameButton,                           // 主持端結束遊戲按鈕
// 		gameModel.TugofwarScreenEndMVPWinnerAvatarFrame,                 // 主持端結算MVP勝方頭像框
// 		gameModel.TugofwarScreenEndBackground,                           // 主持端結算背景
// 		gameModel.TugofwarScreenEndLosingMVPLabel,                       // 主持端結算敗方MVP標示
// 		gameModel.TugofwarScreenEndLosingPrizeList,                      // 主持端結算敗方中獎列表
// 		gameModel.TugofwarScreenEndLosingScoreFrame,                     // 主持端結算敗方分數框
// 		gameModel.TugofwarScreenEndLosingPlayerInfoFrame,                // 主持端結算敗方玩家資訊框
// 		gameModel.TugofwarScreenEndLosingCaptainLabel,                   // 主持端結算敗方隊長標示
// 		gameModel.TugofwarScreenEndWinningPrizeList,                     // 主持端結算勝方中獎列表
// 		gameModel.TugofwarScreenEndWinningScoreFrame,                    // 主持端結算勝方分數框
// 		gameModel.TugofwarScreenEndWinningPlayerInfoFrame,               // 主持端結算勝方玩家資訊框
// 		gameModel.TugofwarScreenEndCaptainMVPNameplateScoreLosingFrame,  // 主持端結算隊長和MVP名牌分數敗方框
// 		gameModel.TugofwarScreenEndCaptainMVPNameplateScoreWinningFrame, // 主持端結算隊長和MVP名牌分數勝方框
// 		gameModel.TugofwarScreenEndCaptainMVPLosingAvatarFrame,          // 主持端結算隊長和MVP敗方頭像框
// 		gameModel.TugofwarScreenEndCaptainWinningAvatarFrame,            // 主持端結算隊長勝方頭像框
// 		gameModel.TugofwarScreenStartSceneAdjustNumberBackground,        // 主持端開始場景調整人數背景
// 		gameModel.TugofwarScreenStartGameButton,                         // 主持端開始遊戲按鈕
// 		gameModel.TugofwarScreenCurrentParticipantsFrame,                // 主持端當前參加人數框
// 		gameModel.TugofwarScreenGamingCaptainAvatarFrame,                // 主持端遊戲中隊長頭像框
// 		gameModel.TugofwarScreenGameCaptainNameFrame,                    // 主持端遊戲隊長名稱框
// 		gameModel.TugofwarScreenAdjustNumberButton,                      // 主持端調整人數按鈕
// 		gameModel.TugofwarRightTeamNameFrame,                            // 右隊名稱框
// 		gameModel.TugofwarLeftTeamNameFrame,                             // 左隊名稱框
// 		gameModel.TugofwarPlayerJoinFactionButton,                       // 玩家端加入陣營按鈕
// 		gameModel.TugofwarPlayerBackground,                              // 玩家端背景
// 		gameModel.TugofwarPlayerTutorialDiagramExample,                  // 玩家端教學示意圖示
// 		gameModel.TugofwarPlayerCaptainHintFrame,                        // 玩家端隊長提示框
// 		gameModel.TugofwarPlayerGamingScoreFrame,                        // 玩家端遊戲中分數框
// 		gameModel.TugofwarPlayerGameAboutToStart,                        // 玩家端遊戲即將開始
// 		gameModel.TugofwarPlayerTitle,                                   // 玩家端標題
// 		gameModel.TugofwarTeamLogoPictureFrame,                          // 隊伍LOGO圖片框
// 		gameModel.TugofwarTeamNumberOfPeopleFrame,                       // 隊伍人數框
// 		gameModel.TugofwarGameStartButton,                               // 加入遊戲按鈕

// 		// 動圖
// 		gameModel.TugofwarScreenEndWinningConfetti,   // 主持端結算勝方彩帶
// 		gameModel.TugofwarScreenEndWinningFireworks,  // 主持端結算勝方煙火
// 		gameModel.TugofwarScreenEndWinLossLabel,      // 主持端結算勝敗標示
// 		gameModel.TugofwarScreenGamingSignHoldingEee, // 主持端遊戲中舉牌蜜蜂
// 		gameModel.TugofwarCharacterSportsman,         // 拔河人物-運動男
// 		gameModel.TugofwarCharacterBee,               // 拔河人物-蜜蜂
// 		gameModel.TugofwarMiddleRibbon,               // 拔河中間緞帶

// 		// 音樂
// 		gameModel.TugofwarBgmStart,  // 遊戲開始
// 		gameModel.TugofwarBgmGaming, // 遊戲進行中
// 		gameModel.TugofwarBgmEnd,    // 遊戲結束
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_5_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// } else if gameType == "bingo" {
// 	fields := []string{
// 		// 賓果遊戲
// 		"bingo_classic_01",
// 		"bingo_classic_02",
// 		"bingo_classic_03",
// 		"bingo_classic_04",
// 		"bingo_classic_05",
// 		"bingo_classic_06",
// 		"bingo_classic_07",
// 		"bingo_classic_08",
// 		"bingo_classic_09",
// 		"bingo_classic_10",
// 		"bingo_classic_11",
// 		"bingo_classic_12",
// 		"bingo_classic_13",
// 		"bingo_classic_14",
// 		"bingo_classic_15",
// 		"bingo_classic_16",
// 		"bingo_classic_17",
// 		"bingo_classic_18",
// 		"bingo_classic_19",
// 		"bingo_classic_20",
// 		"bingo_classic_21",
// 		"bingo_classic_22",
// 		"bingo_classic_23",
// 		"bingo_classic_24",
// 		"bingo_classic_25",
// 		"bingo_classic_26",
// 		"bingo_classic_27",

// 		// 動圖
// 		"bingo_classic_28",
// 		"bingo_classic_29",
// 		"bingo_classic_30",

// 		// 音樂
// 		"bingo_bgm_start",
// 		"bingo_bgm_gaming",
// 		"bingo_bgm_end",
// 	}

// 	values := []string{
// 		// 賓果遊戲自定義
// 		gameModel.BingoClassic01,
// 		gameModel.BingoClassic02,
// 		gameModel.BingoClassic03,
// 		gameModel.BingoClassic04,
// 		gameModel.BingoClassic05,
// 		gameModel.BingoClassic06,
// 		gameModel.BingoClassic07,
// 		gameModel.BingoClassic08,
// 		gameModel.BingoClassic09,
// 		gameModel.BingoClassic10,
// 		gameModel.BingoClassic11,
// 		gameModel.BingoClassic12,
// 		gameModel.BingoClassic13,
// 		gameModel.BingoClassic14,
// 		gameModel.BingoClassic15,
// 		gameModel.BingoClassic16,
// 		gameModel.BingoClassic17,
// 		gameModel.BingoClassic18,
// 		gameModel.BingoClassic19,
// 		gameModel.BingoClassic20,
// 		gameModel.BingoClassic21,
// 		gameModel.BingoClassic22,
// 		gameModel.BingoClassic23,
// 		gameModel.BingoClassic24,
// 		gameModel.BingoClassic25,
// 		gameModel.BingoClassic26,
// 		gameModel.BingoClassic27,

// 		// 動圖
// 		gameModel.BingoClassic28,
// 		gameModel.BingoClassic29,
// 		gameModel.BingoClassic30,

// 		// 音樂
// 		gameModel.BingoBgmStart,
// 		gameModel.BingoBgmGaming,
// 		gameModel.BingoBgmEnd,
// 	}

// 	for i, value := range values {
// 		if !strings.Contains(value, "system") {
// 			file := strings.Split(value, "/admin/uploads/")[1]
// 			// fmt.Println("file: ", file)
// 			// 原始圖片路徑
// 			oldRoute := "../uploads/" + file
// 			newRoute := "../uploads/" + user.UserID + "/" + activity.ActivityID + "/interact/game/" + gameType + "/" +
// 				gameID + "/" + fields[i] + "." + strings.Split(file, ".")[1]
// 			// fmt.Println("oldRoute: ", oldRoute)
// 			// fmt.Println("newRoute: ", newRoute)
// 			err = copy(oldRoute, newRoute)
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// 更新圖片資料
// 			db.Table(config.ACTIVITY_GAME_8_TABLE).WithConn(conn).
// 				Where("game_id", "=", gameID).
// 				Update(command.Value{
// 					fields[i]: fields[i] + "." + strings.Split(file, ".")[1]})

// 		}
// 	}
// }
// }
// 	}

// }

// }
// }

// 將舊的圖片檔案新增至新路徑中
func copy(old string, new string) error {
	// 打開原始檔案
	f, err := os.Open(old)
	if err != nil {
		return errors.New("錯誤: 開啟圖片發生錯誤")
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = err2
		}
	}()

	// 建立新檔案
	ff, err := os.Create(new)
	if err != nil {
		return errors.New("錯誤: 建立新檔案發生錯誤")
	}
	defer func() {
		if err2 := ff.Close(); err2 != nil {
			err = err2
		}
	}()

	// 將資料寫入新檔案中
	_, err = copyZeroAlloc(ff, f)
	if err != nil {
		return errors.New("錯誤: 資料寫入新檔案中發生錯誤")
	}

	return nil
}

// copyZeroAlloc 將用戶上傳的檔案資料寫入遠端的新檔案中
func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	buf := copyBufPool.Get().([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(buf)
	return n, err
}

var copyBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}
