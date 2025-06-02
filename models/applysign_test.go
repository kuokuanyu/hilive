package models

import (
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"testing"
)

// 今天要使用到的匹量更新自定義人員簽到資料
func Test_Applysign_Update(t *testing.T) {
	// 資料表的ip資料記得修改

	// 測試活動: An9d4maX6UtcPAQQgZCM
	// 正式活動: HLYWytZmkengbOfD9fJi

	var (
		activityID = "ciE3NFlYyMkY0Yye6BvE"
	)

	// 查詢活動資料
	activityModel, err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).
		Where("activity_id", "=", activityID).First()
	if err != nil {
		t.Error("錯誤: 查詢活動資訊發生問題")
	}

	attend := utils.GetInt64(activityModel["attend"], 0)
	number := utils.GetInt64(activityModel["number"], 0)

	// 抓取所有自定義人員資料(status: apply)
	applysigns, err := db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).
		Where("activity_id", "=", activityID).
		Where("status", "=", "apply").
		All()
	if err != nil {
		t.Error("錯誤: 取得所有自定義人員資料發生問題(status: apply)")
	}

	// log.Println("資料數: ", len(applysigns))

	// 更新活動資料
	err = db.Table(config.ACTIVITY_TABLE).WithConn(conn).
		Where("activity_id", "=", activityID).Update(command.Value{
		"attend": int(attend) + len(applysigns),
		"number": int(number) + len(applysigns),
	})
	if err != nil {
		// log.Println(err)
		t.Error("錯誤: 更新活動人數及號碼發生錯誤(都更新至2000)")
	}

	// // log.Println("資料數: ", len(applysigns))

	// number := 1
	for _, applysign := range applysigns {
		userID := utils.GetString(applysign["user_id"], "")

		err = db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).
			Where("activity_id", "=", activityID).
			Where("user_id", "=", userID).
			Update(command.Value{
				"status": "sign",
				"number": number,
			})
		if err != nil {
			t.Error("錯誤: 更新所有自定義人員資料發生問題(status: sign)")
		}

		number++
	}


	// 正式區ssh執行 telnet 10.93.112.11 6379
	// del sign_staffs_2_An9d4maX6UtcPAQQgZCM 測試活動
	// del sign_staffs_1_HLYWytZmkengbOfD9fJi 正式活動
	// del sign_staffs_2_HLYWytZmkengbOfD9fJi 正式活動

	// del sign_staffs_1_ciE3NFlYyMkY0Yye6BvE
	// del sign_staffs_2_ciE3NFlYyMkY0Yye6BvE
}

// 測試查詢所有報名簽到人員資訊功能
// func Test_Applysign_FindAll(t *testing.T) {
// 	var (
// 		activityID    = "Bfon6SaV6ORhmuDQUioI"
// 		userID        = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		status        = "apply,sign"
// 		sign          = "sign"
// 		applysignids1 = make([]int64, 0)
// 		applysignids2 = make([]int64, 0)
// 		applysignids3 = make([]int64, 0)
// 		applysignids4 = make([]int64, 0)
// 	)

// 	// 查詢該用戶所有報名簽到資料
// 	applysign1, err := DefaultApplysignModel().SetDbConn(conn).
// 		FindAll("", userID, "")
// 	if err != nil {
// 		t.Error("錯誤: 查詢該用戶所有報名簽到資料錯誤", err)
// 	}
// 	for _, applysign := range applysign1 {
// 		applysignids1 = append(applysignids1, applysign.ID)
// 	}
// 	fmt.Println("查詢該用戶所有報名簽到資料成功(應該有13筆): ", len(applysign1), applysignids1)

// 	// 查詢活動下所有報名簽到資料
// 	applysign2, err := DefaultApplysignModel().SetDbConn(conn).
// 		FindAll(activityID, "", "")
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動下所有報名簽到資料錯誤", err)
// 	}
// 	for _, applysign := range applysign2 {
// 		applysignids2 = append(applysignids2, applysign.ID)
// 	}
// 	fmt.Println("查詢活動下所有報名簽到資料成功(應該有5筆): ", len(applysign2), applysignids2)

// 	// 查詢活動下特定狀態資料(單選)
// 	applysign3, err := DefaultApplysignModel().SetDbConn(conn).
// 		FindAll(activityID, "", sign)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動下特定狀態(單選)資料錯誤", err)
// 	}
// 	for _, applysign := range applysign3 {
// 		applysignids3 = append(applysignids3, applysign.ID)
// 	}
// 	fmt.Println("查詢活動下特定狀態資料成功(單選, 應該有1筆): ", len(applysign3), applysignids3, applysign3)

// 	// 查詢活動下特定狀態資料(複選)
// 	applysign4, err := DefaultApplysignModel().SetDbConn(conn).
// 		FindAll(activityID, "", status)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動下特定狀態資料(複選)錯誤", err)
// 	}
// 	for _, applysign := range applysign4 {
// 		applysignids4 = append(applysignids4, applysign.ID)
// 	}
// 	fmt.Println("查詢活動下特定狀態資料成功(複選, 應該有5筆): ", len(applysign4), applysignids4)
// }

// 測試利用redis查詢簽到人員資料(包含查詢、遞增遞減人數功能)
// func Test_Applysign_FindSignStaff_Redis(t *testing.T) {
// 	var (
// 		activityID = "activity"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		// fieldValues = command.Value{"attend": 100}
// 		// people = 101
// 	)

// 	// 新增測試資料
// 	if err := DefaultActivityModel().SetDbConn(conn).
// 		Add(false, EditActivityModel{
// 			UserID:       "admin",
// 			ActivityID:   activityID,
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "100",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2022-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	}

// 	// 簽到人員資訊
// 	// fmt.Println("簽到人數為0，沒有設置redis資訊")
// 	// _, err := DefaultApplysignModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).FindSignStaff(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢簽到人員資料發生問題", err)
// 	// }

// 	// count := redis.ListLen(config.SIGN_STAFFS_REDIS + activityID)
// 	// fmt.Println("目前redis無簽到人員資料，redis長度應為0: ", count == 0)

// 	// 新增報名人員
// 	id, err := DefaultApplysignModel().SetDbConn(conn).Add(
//  NewApplysignModel{
// 			UserID: userID, ActivityID: activityID, Status: "apply",
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增報名簽到人員發生問題")
// 	}
// 	// 更新為簽到
// 	now, _ := time.ParseInLocation("2006-01-02 15:04:05",
// 		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// 	if err = DefaultApplysignModel().SetDbConn(conn).SetRedisConn(redis).
// 		UpdateStatus(config.HILIVES_NET_URL, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			ActivityID: activityID, UserID: userID,
// 			SignTime: now.Format("2006-01-02 15:04:05"),
// 			Status:   "sign"}, false); err != nil {
// 		t.Error("錯誤: 更新報名簽到人員發生問題")
// 	}

// 	// 簽到人員資訊(這時候會將簽到人員資訊加入redis中)
// 	fmt.Println("這時候簽到人數為1，這次的執行會查詢資料表並將簽到人員資訊加入redis中")
// 	_, err = DefaultApplysignModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindSignStaffs(true, activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢簽到人員資料發生問題", err)
// 	}

// 	fmt.Println("加入redis後，檢查是否正確")
// 	count := redis.ListLen(config.SIGN_STAFFS_REDIS + activityID)
// 	users, _ := redis.ListRange(config.SIGN_STAFFS_REDIS+activityID, 0, 0)
// 	fmt.Println("簽到後取得redis中的簽到人員，應為1: ", count == 1, users)

// 	// fmt.Println("再查詢一次簽到人員，這次會直接從redis中取得，並從redis中取得用戶資訊")
// 	// // 簽到人員資訊(這時候會將簽到人員資訊加入redis中)
// 	// staffs, err := DefaultApplysignModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).FindSignStaff(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢簽到人員資料發生問題", err)
// 	// }
// 	// fmt.Println("從redis中取得的簽到人員資訊: ", staffs)

// 	// 更新為未簽到
// 	if err = DefaultApplysignModel().SetDbConn(conn).SetRedisConn(redis).
// 		UpdateStatus(config.HILIVES_NET_URL, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			ActivityID: activityID, UserID: userID,
// 			Status: "cancel"}, false); err != nil {
// 		t.Error("錯誤: 更新報名簽到人員發生問題")
// 	}

// 	// fmt.Println("取消簽到後，會將redis中的資訊清除，redis裡為空")
// 	// 簽到人員資訊
// 	// _, err = DefaultApplysignModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).FindSignStaff(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢簽到人員資料發生問題", err)
// 	// }

// 	// fmt.Println("檢查是否正確")
// 	// count = redis.ListLen(config.SIGN_STAFFS_REDIS + activityID)
// 	// users, _ = redis.ListRange(config.SIGN_STAFFS_REDIS + activityID, 0, 0)
// 	// fmt.Println("取消簽到後取得redis中的簽到人員，應為0: ", count == 0, users)

// 	// 刪除redis中活動資訊
// 	// redis.DelCache(config.ACTIVITY_REDIS + activityID)
// 	redis.DelCache(config.SIGN_STAFFS_REDIS + activityID)
// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }

// // 測試所有查詢報名人員資料功能
// func Test_Applysign_Find(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		identify   = "yo"
// 	)

// 	// 新增測試活動資料
// 	if err := DefaultActivityModel().SetDbConn(conn).
// 		Add(false, EditActivityModel{
// 			UserID:       "admin",
// 			ActivityID:   "test",
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "100",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2022-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	} else {
// 		t.Log("新增活動資料完成")
// 	}

// 	// 更新自定義欄位資料
// 	if err := DefaultCustomizeModel().SetDbConn(conn).
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

// 	// 新增一筆測試報名簽到人員資料
// 	id, err := DefaultApplysignModel().SetDbConn(conn).
// 		Add(NewApplysignModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			Status:     "apply",
// 		})
// 	if err != nil {
// 		t.Error("加入資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試報名資料完成")
// 	}

// 	// FindSignStaff
// 	applysignModels, err := DefaultApplysignModel().SetDbConn(conn).
// 		FindSignStaffs(false, activityID)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤(FindSignStaff)")
// 	} else {
// 		fmt.Println("FindSignStaff: ", applysignModels)
// 		t.Log("查詢資料完成(FindSignStaff)")
// 	}

// 	// FindByUserID
// 	item, err := DefaultApplysignModel().SetDbConn(conn).FindApplysignLeftJoinActivityCustomize(activityID, userID)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤(FindByUserID)")
// 	} else {
// 		fmt.Println("FindByUserID: ", item)
// 		fmt.Println("ext1相關設置(正常為空)", item.Ext1Name, item.Ext1Type, item.Ext1Options, item.Ext1Required)
// 		// fmt.Println("item.ID: ", item.ID, "item.JoinID: ", item.JoinID)
// 		t.Log("查詢資料完成(FindByUserID)")
// 	}

// 	// LeftJoinActivity
// 	applysignModel, err := DefaultApplysignModel().SetDbConn(conn).FindApplysignLeftJoinActivityCustomize(activityID, userID)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤(FindByUserIDToModel)")
// 	} else if applysignModel.ID == 0 {
// 		t.Error("查詢資料發生錯誤(FindByUserIDToModel): 活動已結束")
// 	} else if applysignModel.Friend == "no" {
// 		t.Error("查詢資料發生錯誤(FindByUserIDToModel): 未加入LINE官方帳號")
// 	} else if applysignModel.ActivityUserID == applysignModel.Identify {
// 		t.Log("查詢資料完成(FindByUserIDToModel): 管理員")
// 	} else {
// 		fmt.Println("LeftJoinActivity: ", applysignModel)
// 		fmt.Println("LeftJoinActivity.ActivityUserID: ", applysignModel.ActivityUserID,
// 			"LeftJoinActivity.start_time: ", applysignModel.StartTime,
// 			"LeftJoinActivity.end_time: ", applysignModel.EndTime,
// 			"LeftJoinActivity.sign_check: ", applysignModel.SignCheck,
// 			"LeftJoinActivity.sign_allow: ", applysignModel.SignAllow,
// 			"LeftJoinActivity.sign_minutes: ", applysignModel.SignMinutes,
// 			"LeftJoinActivity.identify: ", applysignModel.Identify,
// 			"LeftJoinActivity.friend: ", applysignModel.Friend)
// 		t.Log("查詢資料完成(FindByUserIDToModel)")
// 	}

// 	// LeftJoinActivityCustomize
// 	applysignModel, err = DefaultApplysignModel().SetDbConn(conn).FindApplysignLeftJoinActivityCustomize(activityID, identify)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤(LeftJoinActivityCustomize)")
// 	} else {
// 		fmt.Println("LeftJoinActivityCustomize: ", applysignModel)
// 		fmt.Println("ext1相關設置", applysignModel.Ext1Name, applysignModel.Ext1Type, applysignModel.Ext1Options, applysignModel.Ext1Required)
// 		fmt.Println("LeftJoinActivityCustomize.phone: ", applysignModel.Phone,
// 			"LeftJoinActivityCustomize.ext_email: ", applysignModel.ExtEmail)
// 		t.Log("查詢資料完成(LeftJoinActivityCustomize)")
// 	}

// 	// FindByID
// 	applysignModel, err = DefaultApplysignModel().SetDbConn(conn).FindByID(id)
// 	if err != nil {
// 		t.Error("查詢資料發生錯誤(FindByID)")
// 	} else {
// 		fmt.Println("FindByID: ", applysignModel)
// 		fmt.Println("applysignModel.Name(正常為空): ", applysignModel.Name)
// 		fmt.Println("applysignModel.UserID: ", applysignModel.UserID,
// 			"applysignModel.ActivityID: ", applysignModel.ActivityID)
// 		t.Log("查詢資料完成(FindByID)")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }

// // 測試新增一筆報名簽到人員資料並更新報名簽到狀態
// func Test_Applysign_Add_UpdateStatus(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		host       = config.HILIVES_NET_URL
// 	)

// 	// 新增活動資料(人數上限0)
// 	if err := DefaultActivityModel().SetDbConn(conn).
// 		Add(false, EditActivityModel{
// 			UserID:       userID,
// 			ActivityID:   activityID,
// 			ActivityName: "test",
// 			ActivityType: "其他",
// 			People:       "0",
// 			City:         "台中市",
// 			Town:         "南屯區",
// 			StartTime:    "2022-01-01T00:00",
// 			EndTime:      "2022-12-31T00:00",
// 		}); err != nil {
// 		t.Error("新增活動資料發生錯誤")
// 	} else {
// 		t.Log("新增活動資料完成")
// 	}

// 	// 新增一筆測試資料
// 	id, err := DefaultApplysignModel().SetDbConn(conn).
// 		Add(NewApplysignModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			Status:     "review",
// 		})
// 	if err != nil {
// 		t.Error("加入資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試報名資料完成")
// 	}

// 	// 更新為報名完成狀態
// 	now, _ := time.ParseInLocation("2006-01-02 15:04:05",
// 		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// 	err = DefaultApplysignModel().SetDbConn(conn).
// 		UpdateStatus(host, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			ReviewTime: now.Format("2006-01-02 15:04:05"),
// 			Status:     "apply",
// 		}, true)
// 	if err != nil {
// 		t.Error("更新報名完成狀態發生錯誤")
// 	} else {
// 		t.Log("更新報名完成狀態資料完成")
// 	}

// 	// 更新簽到完成狀態(會失敗，因為報名活動人數已達上限)
// 	err = DefaultApplysignModel().SetDbConn(conn).
// 		UpdateStatus(host, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			SignTime:   now.Format("2006-01-02 15:04:05"),
// 			Status:     "sign",
// 		}, true)
// 	if err != nil {
// 		t.Error("更新簽到完成狀態發生錯誤: ", err)
// 	} else {
// 		t.Log("更新簽到完成狀態資料完成")
// 	}

// 	// 更新活動人數上限資料
// 	if err := DefaultActivityModel().SetDbConn(conn).
// 		UpdateActivity(false, EditActivityModel{
// 			ActivityID: activityID,
// 			People:     "100",
// 		}); err != nil {
// 		t.Error("更新活動人數上限資料發生錯誤")
// 	} else {
// 		t.Log("更新活動人數上限資料完成")
// 	}

// 	// 更新簽到完成狀態(會成功)
// 	err = DefaultApplysignModel().SetDbConn(conn).
// 		UpdateStatus(host, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			SignTime:   now.Format("2006-01-02 15:04:05"),
// 			Status:     "sign",
// 		}, true)
// 	if err != nil {
// 		t.Error("更新簽到完成狀態發生錯誤: ", err)
// 	} else {
// 		t.Log("更新簽到完成狀態資料完成")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// // 測試更新報名簽到人員自定義資料並且清空自定義資料功能
// func Test_Applysign_UpdateExt_DeleteExt(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "test"
// 		values     = []string{"test1", "test2", "test3", "test4", "test5",
// 			"test6", "test7", "test8", "test9", "test10"}
// 		field = "ext_2"
// 	)

// 	// 新增一筆測試資料
// 	_, err := DefaultApplysignModel().SetDbConn(conn).
// 		Add( NewApplysignModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			Status:     "apply",
// 		})
// 	if err != nil {
// 		t.Error("加入資料發生錯誤")
// 	} else {
// 		t.Log("新增一筆測試報名資料完成")
// 	}

// 	// 更新自定義資料
// 	err = DefaultApplysignModel().SetDbConn(conn).
// 		UpdateExt(activityID, userID, values)
// 	if err != nil {
// 		t.Error("更新自定義資料發生錯誤")
// 	} else {
// 		t.Log("更新自定義資料完成")
// 	}

// 	// 清空自定義資料
// 	err = DefaultApplysignModel().SetDbConn(conn).
// 		DeleteExt(activityID, field)
// 	if err != nil {
// 		t.Error("清空自定義資料發生錯誤")
// 	} else {
// 		t.Log("清空自定義資料完成")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }
