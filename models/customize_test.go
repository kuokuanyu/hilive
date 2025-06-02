package models

// 測試更新activity_customize資料時是否有清空(報名簽到人員、參加遊戲人員、中獎人員)ext_n欄位資料
// func Test_Customize_Update(t *testing.T) {
// 	// 新增一筆測試資料
// 	_, err := db.Table("activity_customize").WithConn(conn).
// 		Insert(command.Value{"activity_id": "test"})
// 	if err != nil {
// 		t.Error("加入資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試自定義資料完成")
// 	}

// 	// 新增一筆報名簽到、參加遊戲人員、中獎人員資料
// 	_, err = db.Table("activity_applysign").WithConn(conn).
// 		Insert(command.Value{
// 			"user_id":     "test1",
// 			"activity_id": "test",
// 			"name":        "test",
// 			"avatar":      "test",
// 			"status":      "sign",
// 			"ext_1":       "123",
// 		})
// 	if err != nil {
// 		t.Error("加入報名簽到人員資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試報名簽到人員資料完成")
// 	}
// 	_, err = db.Table("activity_staff_game").WithConn(conn).
// 		Insert(command.Value{
// 			"user_id":     "test1",
// 			"activity_id": "test",
// 			"game_id":     "test",
// 			"name":        "test",
// 			"avatar":      "test",
// 			"round":       1,
// 			"status":      "no",
// 			"black":       "no",
// 			"ext_1":       "123",
// 		})
// 	if err != nil {
// 		t.Error("加入參加遊戲人員資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試參加遊戲人員資料完成")
// 	}
// 	_, err = db.Table("activity_staff_prize").WithConn(conn).
// 		Insert(command.Value{
// 			"user_id":     "test1",
// 			"activity_id": "test",
// 			"game_id":     "test",
// 			"prize_id":    "test",
// 			"name":        "test",
// 			"avatar":      "test",
// 			"prize_name":  "test",
// 			"picture":     "test",
// 			"price":       100,
// 			"method":      "site",
// 			"password":    "test",
// 			"round":       1,
// 			"status":      "no",
// 			"white":       "no",
// 			"prize_type":  "first",
// 			"ext_1":       "123",
// 		})
// 	if err != nil {
// 		t.Error("加入中獎人員資料發生錯誤:", err)
// 	} else {
// 		t.Log("新增一筆測試參加中獎人員資料完成")
// 	}

// 	// 更新自定義欄位資料
// 	if err = DefaultCustomizeModel().SetDbConn(conn).
// 		Update(EditCustomizeModel{
// 			ActivityID: "test",
// 			Field:      "ext_1",
// 			Name:       "生日",
// 			Type:       "text",
// 			Options:    "2022/01/01",
// 			Required:   "true",
// 		}); err != nil {
// 		t.Error("更新自定義欄位資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新自定義欄位資料完成")
// 	}

// 	// 檢查更新自定義欄位後，該欄位是否有清空資料
// 	applysignModel, err := DefaultApplysignModel().SetDbConn(conn).
// 	Find(0, "test", "test1", false)
// 	if err != nil {
// 		t.Error("查詢報名簽到資料發生錯誤")
// 	} else if applysignModel.Ext1 != "" {
// 		t.Error("更新報名簽到自定義欄位資料發生錯誤: 沒有清空ext_1欄位資料")
// 	}
// 	mapModel, err := db.Table("activity_staff_game").WithConn(conn).
// 		Where("activity_id", "=", "test").First()
// 	if err != nil {
// 		t.Error("查詢參加遊戲人員資料發生錯誤")
// 	} else if mapModel["ext_1"] != "" {
// 		t.Error("更新參加遊戲人員自定義欄位資料發生錯誤: 沒有清空ext_1欄位資料")
// 	}
// 	mapModel, err = db.Table("activity_staff_prize").WithConn(conn).
// 		Where("activity_id", "=", "test").First()
// 	if err != nil {
// 		t.Error("查詢參加遊戲人員資料發生錯誤")
// 	} else if mapModel["ext_1"] != "" {
// 		t.Error("更新參加中獎人員自定義欄位資料發生錯誤: 沒有清空ext_1欄位資料")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// 測試將activity表所有活動資料加入activity_customize表(更新至正式區時需要執行)
// func Test_Add_Customize_Data(t *testing.T) {
// 	var data int

// 	// 先取得activity所有資料
// 	items, err := db.Table("activity").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("查詢activity所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity所有資料完成")
// 	}

// 	// 將所有活動資料加入activity_customize
// 	for _, item := range items {
// 		err = DefaultCustomizeModel().SetDbConn(conn).
// 			Add(EditCustomizeModel{
// 				ActivityID: utils.GetString(item["activity_id"], ""),
// 			})
// 		if err != nil {
// 			t.Error("加入資料發生錯誤")
// 		} else {
// 			data++
// 		}
// 	}

// 	if data == len(items) {
// 		t.Log("所有活動資料加入activity_customize完成，資料所有筆數: ", data)
// 	} else {
// 		t.Error("所有活動資料加入activity_customize發生異常，資料數不正確")
// 	}
// }
