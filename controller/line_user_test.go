package controller

// 測試報名簽到人員自定義欄位更新成功
// func Test_Update_Ext_Data(t *testing.T) {
// 	// 新增測試資料
// 	// 報名簽到
// 	if _, err := db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).
// 		Insert(command.Value{"user_id": "test", "activity_id": "test", "status": "review"}); err != nil {
// 		t.Error("錯誤: 新增遊戲參加人員資料發生錯誤")
// 	}
// 	// 遊戲參加人員
// 	// if _, err := db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Insert(command.Value{
// 	// 	"user_id": "test", "activity_id": "test", "game_id": "test",
// 	// 	"name": "test", "avatar": "test", "round": 1,
// 	// 	"status": "no", "black": "no"}); err != nil {
// 	// 	t.Error("錯誤: 新增遊戲參加人員資料發生錯誤")
// 	// }
// 	// 中獎人員
// 	// if _, err := db.Table("activity_staff_prize").WithConn(conn).
// 	// 	Insert(command.Value{
// 	// 		"user_id": "test", "activity_id": "test", "game_id": "test", "prize_id": "test",
// 	// 		"name": "test", "avatar": "test", "prize_name": "test", "picture": "test",
// 	// 		"price": 100, "method": "site", "password": "test", "round": 1,
// 	// 		"status": "no", "white": "no", "prize_type": "first"}); err != nil {
// 	// 	t.Error("新增中獎人員資料發生錯誤:", err)
// 	// }

// 	var values = []string{"ext1", "ext2",
// 		"ext3", "ext4", "ext5", "", "", "", "", ""}
// 	// 更新報名簽到人員、參加遊戲人員、中獎人員ext欄位資料
// 	if err := models.DefaultApplysignModel().SetDbConn(conn).UpdateExt(
// 		"test", "test", values); err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	// if err := models.DefaultGameStaffModel().SetDbConn(conn).UpdateExt(
// 	// 	"test", "test", values); err != nil {
// 	// 	t.Error(err.Error())
// 	// 	return
// 	// }
// 	// if err := models.DefaultPrizeStaffModel().SetDbConn(conn).UpdateExt(
// 	// 	"test", "test", values); err != nil {
// 	// 	t.Error(err.Error())
// 	// 	return
// 	// }

// 	// 刪除測試資料
// 	// db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// }

// // 測試更新LINE用戶信箱與電話資料
// func Test_Line_UpdatePhoneAndEmail(t *testing.T) {
// 	// 新增測試資料
// 	if _, err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).Insert(command.Value{
// 		"user_id": "test", "name": "test", "avatar": "test", "email": "test",
// 		"identify": "test", "friend": "no", "phone": "", "ext_email": "",
// 	}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		t.Error("錯誤: 新增用戶1資料發生問題")
// 	} else {
// 		t.Log("完成新增LINE用戶1")
// 	}
// 	// if _, err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).Insert(command.Value{
// 	// 	"user_id": "test2", "name": "test2", "avatar": "test2", "email": "test2",
// 	// 	"identify": "test2", "friend": "no"}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	// 	t.Error("錯誤: 新增用戶2資料發生問題")
// 	// } else {
// 	// 	t.Log("完成新增LINE用戶2")
// 	// }

// 	// LINE用戶資訊
// 	// lineModel, err := models.DefaultLineModel().SetDbConn(conn).Find("identify", "test2")
// 	// if err != nil || lineModel.ID == 0 {
// 	// 	t.Error("錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
// 	// }

// 	// 檢查電話是否被使用過
// 	var phone = "0987654321"
// 	var email = "test@gmail.com"
// 	if phone != "" {
// 		if len(phone) > 2 {
// 			if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
// 				t.Error("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 			}
// 		} else {
// 			t.Error("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 		}
// 		// if user, err := models.DefaultLineModel().SetDbConn(conn).
// 		// 	FindByPhone(phone); user.Phone != "" || err != nil {
// 		// 	t.Error("錯誤: 電話號碼已被註冊過，請輸入有效的手機號碼: ", user.ID)
// 		// }
// 	} else {
// 		t.Error("未輸入電話")
// 	}
// 	if email != "" && !strings.Contains(email, "@") {
// 		t.Error("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
// 	} else if email == "" {
// 		fmt.Println("未輸入電子信箱")
// 	}

// 	// 更新電話以及自定義電子信箱資訊
// 	if err := models.DefaultLineModel().SetDbConn(conn).UpdatePhoneAndEmail(
// 		models.LineModel{UserID: "test", Phone: phone, ExtEmail: email}); err != nil {
// 		t.Error(err.Error())
// 	}

// 	// 刪除測試資料
// 	db.Table(config.LINE_USERS_TABLE).WithConn(conn).Where("user_id", "=", "test").Delete()
// 	// db.Table(config.LINE_USERS_TABLE).WithConn(conn).Where("user_id", "=", "test2").Delete()
// }
