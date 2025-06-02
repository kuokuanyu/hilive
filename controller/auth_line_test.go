package controller

// // 測試activity join activity_customize表時參數顯示方式
// func Test_Activity_Join_Customize_Data(t *testing.T) {
// 	now, _ := time.ParseInLocation("2006-01-02 15:04:05",
// 		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// 	item, err := db.Table("activity").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity.activity_id",
// 			FieldB:    "activity_customize.activity_id",
// 			Table:     "activity_customize",
// 			Operation: "=",
// 		}).
// 		Where("end_time", ">", now).Where("activity.activity_id", "=", "FEVNqoH9Vv3iDlo7byeB").First()
// 	id, _ := item["id"].(int64)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤")
// 	} else if id == 0 {
// 		t.Error("活動已結束")
// 	}
// 	ext, _ := item["ext_1"].(string)
// 	fmt.Println(item)
// 	fmt.Println(ext)
// 	fmt.Println(item["apply_check"].(string))
// 	fmt.Println(item["activity_name"].(string))
// 	// fmt.Println(item["ext_1"].(string))
// }

// // 測試新增報名人員資料
// func Test_Add_Applysign_Data(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增活動資料
// 	if err := models.DefaultActivityModel().SetDbConn(conn).
// 		Add(false, models.EditActivityModel{
// 			UserID:       "admin",
// 			ActivityID:   activityID,
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "100",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2020-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	} else {
// 		t.Log("新增活動資料完成")
// 	}

// 	// 更新自定義欄位資料
// 	// if err := models.DefaultCustomizeModel().SetDbConn(conn).
// 	// 	Update(models.EditCustomizeModel{
// 	// 		ActivityID: "test",
// 	// 		Field:      "ext_1",
// 	// 		Name:       "生日",
// 	// 		Type:       "text",
// 	// 		Required:   "true",
// 	// 	}); err != nil {
// 	// 	t.Error("更新自定義欄位1資料發生錯誤: ", err)
// 	// } else {
// 	// 	t.Log("更新自定義欄位資料完成")
// 	// }
// 	// 更新自定義欄位資料
// 	// if err := models.DefaultCustomizeModel().SetDbConn(conn).
// 	// 	Update(models.EditCustomizeModel{
// 	// 		ActivityID: "test",
// 	// 		Field:      "ext_2",
// 	// 		Name:       "地址",
// 	// 		Type:       "text",
// 	// 		Required:   "false",
// 	// 	}); err != nil {
// 	// 	t.Error("更新自定義欄位2資料發生錯誤: ", err)
// 	// } else {
// 	// 	t.Log("更新自定義欄位資料完成")
// 	// }

// 	// 取得活動狀態資訊
// 	activityModel, err := models.DefaultActivityModel().SetDbConn(conn).
// 		FindActivityLeftJoinCustomize(activityID)
// 	if err != nil || activityModel.ID == 0 {
// 		t.Error("錯誤: 活動已結束，謝謝參與")
// 		return
// 	}

// 	// 報名簽到人員資料
// 	applysign, err := models.DefaultApplysignModel().SetDbConn(conn).
// 		FindUserApplysignStatus(activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得報名簽到人員資料")
// 	}

// 	// 判斷報名審核開關
// 	if applysign.ID == 0 { // 無人員資料
// 		var status string
// 		if activityModel.ApplyCheck == "open" { // 報名審核開啟
// 			status = "review"
// 			// message = fmt.Sprintf("已提交%s的活動報名，請等待主辦單位審核", activityModel.ActivityName)
// 		} else { // 報名審核關閉
// 			status = "apply"
// 			// message = fmt.Sprintf("完成%s的活動報名，請於活動當天依照主辦單位指示完成活動簽到", activityModel.ActivityName)
// 		}

// 		// 新增人員資料
// 		_, err := models.DefaultApplysignModel().SetDbConn(conn).Add(
// 			models.NewApplysignModel{
// 				UserID: userID, ActivityID: activityID, Status: status,
// 			})
// 		if err != nil {
// 			t.Error("錯誤: 新增報名簽到人員發生問題")
// 		}
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }

// // 測試更新報名人員資料
// func Test_Update_Applysign_Data(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增活動資料
// 	if err := models.DefaultActivityModel().SetDbConn(conn).
// 		Add(false, models.EditActivityModel{
// 			UserID:       "admin",
// 			ActivityID:   activityID,
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "100",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2020-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	} else {
// 		t.Log("新增活動資料完成")
// 	}

// 	// 修改報名審核開關
// 	if err := models.DefaultActivityModel().SetDbConn(conn).
// 		UpdateApply(models.EditApplyModel{
// 			ActivityID: activityID, ApplyCheck: "close",
// 		}); err != nil {
// 		t.Error("修改活動報名審核開關發生錯誤")
// 	}

// 	// 新增報名人員人員資料
// 	id, err := models.DefaultApplysignModel().SetDbConn(conn).Add(
// 		models.NewApplysignModel{
// 			UserID: userID, ActivityID: activityID, Status: "review",
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增報名簽到人員發生問題")
// 	}

// 	// 取得活動狀態資訊
// 	activityModel, err := models.DefaultActivityModel().SetDbConn(conn).
// 		FindActivityLeftJoinCustomize(activityID)
// 	if err != nil || activityModel.ID == 0 {
// 		t.Error("錯誤: 活動已結束，謝謝參與")
// 		return
// 	}

// 	// 報名簽到人員資料
// 	applysign, err := models.DefaultApplysignModel().SetDbConn(conn).
// 		FindUserApplysignStatus(activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得報名簽到人員資料")
// 	}

// 	if applysign.ID != 0 {
// 		// 原本資料為審核中但管理員已將報名審核關閉，因此更新為報名完成狀態
// 		if activityModel.ApplyCheck == "close" && applysign.Status == "review" {
// 			fmt.Println("原本資料為審核中但管理員已將報名審核關閉，因此更新為報名完成狀態")
// 			if err = db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("id", "=", id).
// 				Update(command.Value{"status": "apply"}); err != nil {
// 				t.Error("錯誤: 更新報名簽到人員資料發生問題")
// 			}
// 		}
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// // 測試判斷用戶資料是否補齊
// func Test_Data_Check(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增活動資料
// 	if err := models.DefaultActivityModel().SetDbConn(conn).
// 		Add(false, models.EditActivityModel{
// 			UserID:       "admin",
// 			ActivityID:   activityID,
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "100",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2020-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	} else {
// 		t.Log("新增活動資料完成")
// 	}
// 	// 更新自定義欄位資料
// 	if err := models.DefaultCustomizeModel().SetDbConn(conn).
// 		Update(models.EditCustomizeModel{
// 			ActivityID: activityID,
// 			Field:      "ext_1",
// 			Name:       "生日",
// 			Type:       "text",
// 			Required:   "false",
// 		}); err != nil {
// 		t.Error("更新自定義欄位1資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新自定義欄位資料完成")
// 	}
// 	// 更新自定義欄位資料
// 	if err := models.DefaultCustomizeModel().SetDbConn(conn).
// 		Update(models.EditCustomizeModel{
// 			ActivityID: activityID,
// 			Field:      "ext_2",
// 			Name:       "地址",
// 			Type:       "text",
// 			Required:   "false",
// 		}); err != nil {
// 		t.Error("更新自定義欄位2資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新自定義欄位資料完成")
// 	}

// 	if _, err := db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Insert(command.Value{
// 		"user_id":     "U3140fd73cfd35dd992668ab3b6efdae9",
// 		"activity_id": activityID,
// 		"status":      "sign",
// 		"ext_1":       "test",
// 	}); err != nil {
// 		t.Error("錯誤: 新增報名簽到人員發生問題")
// 	}

// 	// 取得活動狀態資訊
// 	activityModel, err := models.DefaultActivityModel().SetDbConn(conn).
// 		FindActivityLeftJoinCustomize(activityID)
// 	if err != nil || activityModel.ID == 0 {
// 		t.Error("錯誤: 活動已結束，謝謝參與")
// 		return
// 	}

// 	// 報名簽到人員資料
// 	applysign, err := models.DefaultApplysignModel().SetDbConn(conn).
// 		FindUserApplysignStatus(activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得報名簽到人員資料")
// 	}

// 	// 判斷自定義欄位的必填資料是否補齊
// 	var exts = []string{activityModel.Ext1Name, activityModel.Ext2Name, activityModel.Ext3Name,
// 		activityModel.Ext4Name, activityModel.Ext5Name, activityModel.Ext6Name, activityModel.Ext7Name,
// 		activityModel.Ext8Name, activityModel.Ext9Name, activityModel.Ext10Name}
// 	var values = []string{applysign.Ext1, applysign.Ext2, applysign.Ext3,
// 		applysign.Ext4, applysign.Ext5, applysign.Ext6, applysign.Ext7,
// 		applysign.Ext8, applysign.Ext9, applysign.Ext10}

// 	for i, ext := range exts {
// 		if ext != "" && values[i] == "" {
// 			fmt.Println("未補齊")
// 		} else {
// 			fmt.Println("補齊")
// 		}
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// // 測試更新用戶資料
// func Test_Update_LINE_User_Data(t *testing.T) {
// 	var payload = &social.Payload{
// 		Sub:     "test",
// 		Name:    "test123",
// 		Picture: "test123",
// 		Email:   "test123@gmail.com",
// 	}

// 	// 新增測試資料
// 	if _, err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).
// 		Insert(command.Value{
// 			"user_id": "test", "name": "test", "avatar": "test",
// 			"email": "test", "identify": "test", "friend": "no",
// 			"phone": "0912345678", "ext_email": "ext_test",
// 		}); err != nil {
// 		t.Error("錯誤: 新增用戶資料發生問題")
// 	} else {
// 		t.Log("完成新增測試LINE用戶")
// 	}
// 	// // 報名簽到
// 	// if _, err := models.DefaultApplysignModel().SetDbConn(conn).Add(
// 	// 	models.NewApplysignModel{
// 	// 		UserID: "test", ActivityID: "test", Status: "review",
// 	// 	}, "測試"); err != nil {
// 	// 	t.Error("錯誤: 新增報名簽到資料發生錯誤")
// 	// }
// 	// // 聊天紀錄
// 	// if err := models.DefaultChatroomRecordModel().SetDbConn(conn).Add(
// 	// 	models.EditChatroomRecordModel{
// 	// 		UserID: "test", ActivityID: "test", Name: "test", Avatar: "test",
// 	// 		MessageType: "normal-barrage", MessageStyle: "default", Message: "test",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增聊天紀錄資料發生錯誤")
// 	// }
// 	// // 遊戲參加人員
// 	// if _, err := db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Insert(command.Value{
// 	// 	"user_id": "test", "activity_id": "test", "game_id": "test",
// 	// 	"name": "test", "avatar": "test", "round": 1,
// 	// 	"status": "no", "black": "no"}); err != nil {
// 	// 	t.Error("錯誤: 新增遊戲參加人員資料發生錯誤")
// 	// }
// 	// // 中獎人員
// 	// if _, err := db.Table("activity_staff_prize").WithConn(conn).
// 	// 	Insert(command.Value{
// 	// 		"user_id": "test", "activity_id": "test", "game_id": "test", "prize_id": "test",
// 	// 		"name": "test", "avatar": "test", "prize_name": "test", "picture": "test",
// 	// 		"price": 100, "method": "site", "password": "test", "round": 1,
// 	// 		"status": "no", "white": "no", "prize_type": "first"}); err != nil {
// 	// 	t.Error("新增中獎人員資料發生錯誤:", err)
// 	// }
// 	// // 打地鼠計分表
// 	// if err := models.DefaultWhackMoleScoreModel().SetDbConn(conn).Add(models.NewWhackMoleScoreModel{
// 	// 	UserID: "test", ActivityID: "test", GameID: "test", Name: "test", Avatar: "test", Round: "1"}); err != nil {
// 	// 	t.Error("錯誤: 新增打地鼠計分表資料發生錯誤")
// 	// }
// 	// // 提問訊息
// 	// if err := models.DefaultQuestionUserModel().SetDbConn(conn).Add(
// 	// 	models.NewQuestionUserModel{
// 	// 		UserID: "test", ActivityID: "test", Name: "test", Avatar: "test", Content: "test",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增提問訊息資料發生錯誤")
// 	// }
// 	// // 按讚紀錄
// 	// if err := models.DefaultQuestionLikesRecordModel().SetDbConn(conn).Add(
// 	// 	models.NewQuestionLikesRecordModel{
// 	// 		UserID: "test", ActivityID: "test", Name: "test", Avatar: "test", QuestionID: 1,
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增按讚紀錄資料發生錯誤")
// 	// }

// 	// 是否存在用戶資料
// 	lineModel, err := models.DefaultLineModel().SetDbConn(conn).
// 		Find(false, "", "user_id", "test")
// 	if err != nil {
// 		t.Error("錯誤: 取得用戶資料發生問題")
// 	} else if lineModel.ID == 0 {
// 		t.Error("錯誤: 無法查詢LINE用戶")
// 	} else if lineModel.ID != 0 {
// 		t.Log("完成查詢LINE用戶")
// 	}

// 	fmt.Println("判斷用戶資訊是否有變化")
// 	fmt.Println(lineModel.Name != payload.Name)
// 	fmt.Println(lineModel.Email != payload.Email)
// 	fmt.Println(lineModel.Avatar != payload.Picture)

// 	if lineModel.Name != payload.Name || lineModel.Email != payload.Email ||
// 		lineModel.Avatar != payload.Picture {
// 		fmt.Println("個人資料有異動，更新所有用戶相關資料")
// 		if err = models.DefaultLineModel().SetDbConn(conn).UpdateUser(models.LineModel{
// 			UserID: payload.Sub, Name: payload.Name, Avatar: payload.Picture,
// 			Email: payload.Email}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 			t.Error("錯誤: 更新用戶資料發生問題")
// 			return
// 		} else {
// 			t.Log("完成更新LINE用戶")
// 		}
// 		// 更新LINE用戶最新的名稱、頭像、email資料
// 		// 聊天紀錄
// 		// if err := models.DefaultChatroomRecordModel().SetDbConn(conn).UpdateUser(
// 		// 	payload.Sub, payload.Name, payload.Picture); err != nil {
// 		// 	t.Error("錯誤: 更新聊天紀錄資料發生錯誤，請重新報名簽到")
// 		// 	return
// 		// }
// 		// // 打地鼠計分表
// 		// if err := models.DefaultWhackMoleScoreModel().SetDbConn(conn).UpdateUser(
// 		// 	payload.Sub, payload.Name, payload.Picture); err != nil {
// 		// 	t.Error("錯誤: 更新打地鼠計分表資料發生錯誤，請重新報名簽到")
// 		// 	return
// 		// }
// 		// // 提問訊息
// 		// if err := models.DefaultQuestionUserModel().SetDbConn(conn).UpdateUser(
// 		// 	payload.Sub, payload.Name, payload.Picture); err != nil {
// 		// 	t.Error("錯誤: 更新提問訊息資料發生錯誤，請重新報名簽到")
// 		// 	return
// 		// }
// 		// // 按讚紀錄
// 		// if err := models.DefaultQuestionLikesRecordModel().SetDbConn(conn).UpdateUser(
// 		// 	payload.Sub, payload.Name, payload.Picture); err != nil {
// 		// 	t.Error("錯誤: 更新按讚紀錄資料發生錯誤，請重新報名簽到")
// 		// 	return
// 		// }
// 	}
// 	// 報名簽到
// 	// if err := models.DefaultApplysignModel().SetDbConn(conn).UpdateUser(
// 	// 	payload.Sub, payload.Name, payload.Picture, payload.Email,
// 	// 	lineModel.ExtEmail, lineModel.Phone); err != nil {
// 	// 	t.Error("錯誤: 更新報名簽到資料發生錯誤，請重新報名簽到")
// 	// 	return
// 	// }

// 	// 遊戲參加人員
// 	// if err := models.DefaultGameStaffModel().SetDbConn(conn).UpdateUser(
// 	// 	payload.Sub, payload.Name, payload.Picture, payload.Email,
// 	// 	lineModel.ExtEmail, lineModel.Phone); err != nil {
// 	// 	t.Error("錯誤: 更新遊戲參加人員資料發生錯誤，請重新報名簽到")
// 	// 	return
// 	// }
// 	// // 中獎人員
// 	// if err := models.DefaultPrizeStaffModel().SetDbConn(conn).UpdateUser(
// 	// 	payload.Sub, payload.Name, payload.Picture, payload.Email,
// 	// 	lineModel.ExtEmail, lineModel.Phone); err != nil {
// 	// 	t.Error("錯誤: 更新中獎人員資料發生錯誤，請重新報名簽到")
// 	// 	return
// 	// }

// 	// 刪除測試資料
// 	db.Table(config.LINE_USERS_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_CHATROOM_RECORD_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_WHACK_MOLE_SCORE_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_QUESTION_USER_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// }

// // 測試新增用戶資料
// func Test_Insert_LINE_User_Data(t *testing.T) {
// 	var payload = &social.Payload{
// 		Sub:     "test",
// 		Name:    "test",
// 		Picture: "test",
// 		Email:   "test@gmail.com",
// 	}

// 	// 是否存在用戶資料
// 	lineModel, err := models.DefaultLineModel().SetDbConn(conn).
// 		Find(false, "", "user_id", payload.Sub)
// 	if err != nil {
// 		t.Error("錯誤: 取得用戶資料發生問題")
// 		return
// 	} else {
// 		t.Log("完成查詢LINE用戶")
// 	}

// 	if lineModel.ID == 0 {
// 		fmt.Println("新增用戶資料")
// 		if _, err = models.DefaultLineModel().SetDbConn(conn).Add(models.LineModel{
// 			UserID:   payload.Sub,
// 			Name:     payload.Name,
// 			Avatar:   payload.Picture,
// 			Email:    payload.Email,
// 			Identify: "test",
// 			Friend:   "no",
// 		}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 			t.Error("錯誤: 新增用戶資料發生問題")
// 			return
// 		} else {
// 			t.Log("完成新增LINE用戶")
// 		}
// 	} else {
// 		t.Error("錯誤: 已有用戶資料")
// 	}

// 	// 刪除測試資料
// 	db.Table(config.LINE_USERS_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// }
