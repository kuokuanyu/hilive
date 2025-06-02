package models

import (
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"testing"
)

// 測試利用redis查詢活動資料(包含查詢、遞增遞減人數功能)
// func Test_Activity_Find_Redis(t *testing.T) {
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

// 	// 活動資訊
// 	_, err := DefaultActivityModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動資料發生問題", err)
// 	}

// 	// value, _ := redis.GetCache(config.ACTIVITY_REDIS + activityID)
// 	// fmt.Println("目前redis無活動資訊資料，查詢後並設置redis應為0: ", value == "0")

// 	// 新增報名人員
// 	id, err := DefaultApplysignModel().SetDbConn(conn).Add(
// 		NewApplysignModel{
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

// 	// 活動資訊
// 	activityModel, err := DefaultActivityModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動資料發生問題", err)
// 	}
// 	fmt.Println("簽到後取得redis中的活動人數，應為1: ", activityModel.Attend == 1)

// 	// 更新為未簽到
// 	if err = DefaultApplysignModel().SetDbConn(conn).SetRedisConn(redis).
// 		UpdateStatus(config.HILIVES_NET_URL, EditApplysignModel{
// 			ID:         strconv.Itoa(int(id)),
// 			ActivityID: activityID, UserID: userID,
// 			Status: "cancel"}, false); err != nil {
// 		t.Error("錯誤: 更新報名簽到人員發生問題")
// 	}

// 	// 活動資訊
// 	activityModel, err = DefaultActivityModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動資料發生問題", err)
// 	}
// 	fmt.Println("取消簽到後取得redis中的活動人數，應為0: ", activityModel.Attend == 0)

// 	// 刪除redis中活動資訊
// 	// redis.DelCache(config.ACTIVITY_REDIS + activityID)
// 	redis.DelCache(config.SIGN_STAFFS_REDIS + activityID)
// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }

// 測試更新活動人數上限判斷(更新的人數不能小於目前的簽到人數attend)
func Test_Activity_UpdatePeople(t *testing.T) {
	var (
		activityID = "lubvQvBqyB7Kp3ckTaHR"
		// fieldValues = command.Value{"attend": 100}
		people = 101
	)

	// 更新的people值必須要大於attend(簽到完成)的數量
	if err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).
		Where("activity_id", "=", activityID).
		Where("attend", "<=", people).
		Update(command.Value{"people": people}); err != nil {
		t.Error("錯誤: 活動人數上限數值必須大於目前的參加人數: ", err)
	} else {
		fmt.Println("更新成功")
	}
}

// 測試更新活動人數語法判斷
func Test_Activity_UpdateAttend(t *testing.T) {
	// var (
	// 	activityID = "lubvQvBqyB7Kp3ckTaHR"
	// 	// fieldValues = command.Value{"attend": 101}
	// 	attend = 102
	// )

	// if err := DefaultActivityModel().SetDbConn(conn).
	// 	UpdateAttend(activityID, attend); err != nil {
	// 	t.Error("更新活動人數發生錯誤: ", err)
	// }
	// fmt.Println("更新活動人數成功")

	// _, err := conn.Exec(
	// 	fmt.Sprintf("update `activity` set `attend` = %d where `activity_id` = %s and `attend` <= `people`",
	// 		attend, fmt.Sprintf("'%s'", activityID)))
	// if err != nil {
	// 	t.Error("更新活動人數發生錯誤: ", err)
	// }
	// if affectRow, _ := res.RowsAffected(); affectRow < 1 {
	// 	t.Error("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
	// }

	// if err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).
	// 	WhereRaw("`activity_id` = ? and `attend` <= `people`", activityID).
	// 	Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	t.Error("更新活動人數發生錯誤: ", err)
	// }
}

// 測試遞減活動人數語法判斷
func Test_Activity_DecrAttend(t *testing.T) {
	// var (
	// 	activityID = "lubvQvBqyB7Kp3ckTaHR"
	// 	// fieldValues = command.Value{"attend": "attend - 1"}
	// )

	// if err := DefaultActivityModel().SetDbConn(conn).
	// 	DecrAttend(false, activityID, ""); err != nil {
	// 	t.Error("遞減活動人數發生錯誤: ", err)
	// }
	// fmt.Println("遞減活動人數成功")

	// _, err := conn.Exec(
	// 	fmt.Sprintf("update `activity` set `attend` = attend - 1 where `activity_id` = %s and `attend` <= `people`",
	// 		fmt.Sprintf("'%s'", activityID)))
	// if err != nil {
	// 	t.Error("減少活動人數發生錯誤: ", err)
	// }
	// if affectRow, _ := res.RowsAffected(); affectRow < 1 {
	// 	t.Error("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
	// }

	// if err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).
	// 	WhereRaw("`activity_id` = ? and `attend` <= `people`", activityID).
	// 	Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	t.Error("更新活動人數發生錯誤: ", err)
	// }

}

// 測試遞增活動人數語法判斷
func Test_Activity_IncrAttend(t *testing.T) {
	// var (
	// 	activityID = "lubvQvBqyB7Kp3ckTaHR"
	// 	// fieldValues = command.Value{"attend": "attend + 1", "number": "number + 1"}
	// )

	// if err := DefaultActivityModel().SetDbConn(conn).
	// 	IncrAttend(false, activityID, ""); err != nil {
	// 	t.Error("遞增活動人數發生錯誤: ", err)
	// }
	// fmt.Println("遞增活動人數成功")
	// res, err := conn.Exec(
	// 	fmt.Sprintf("update `activity` set `attend` = attend + 1, `number` = number + 1 where `activity_id` = %s and `attend` < `people`",
	// 		fmt.Sprintf("'%s'", activityID)))
	// if err != nil {
	// 	t.Error("更新活動人數發生錯誤: ", err)
	// }
	// if affectRow, _ := res.RowsAffected(); affectRow < 1 {
	// 	t.Error("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
	// }

	// if err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).
	// 	WhereRaw("`activity_id` = ? and `attend` < `people`", activityID).
	// 	Update(fieldValues); err != nil {
	// 	t.Error("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
	// }
}

// 測試新增活動資料，是否有添加該活動自定義欄位資料
// func Test_Activity_LeftJoinActivityCustomize(t *testing.T) {
// 	// 新增測試資料
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
// 	if err := DefaultCustomizeModel().SetDbConn(conn).
// 		Update(EditCustomizeModel{
// 			ActivityID: "test",
// 			Field:      "ext_2",
// 			Name:       "性別",
// 			Type:       "checkbox",
// 			Options:    "男&&&女",
// 			Required:   "true",
// 		}); err != nil {
// 		t.Error("更新自定義欄位資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新自定義欄位資料完成")
// 	}

// 	// 測試LeftJoinActivityCustomize查詢資料功能
// 	item, err := DefaultActivityModel().SetDbConn(conn).FindActivityLeftJoinCustomize("test")
// 	if err != nil {
// 		t.Error("查詢join資料發生錯誤(LeftJoinActivityCustomize): ", err)
// 	} else if item.ID == 0 {
// 		t.Error("查詢join資料發生錯誤(LeftJoinActivityCustomize): 活動已過期")
// 	} else {
// 		fmt.Println("item: ", item)
// 		fmt.Println("item.Ext1Name: ", item.Ext1Name, "item.Ext2Name: ", item.Ext2Name)
// 		t.Log("查詢join資料發生成功(LeftJoinActivityCustomize)")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// 測試新增活動資料，是否有添加該活動自定義欄位資料
// func Test_Activity_Add(t *testing.T) {
// 	// 新增活動資料
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

// 	// 查詢自定義欄位是否有相關資料
// 	item, err := DefaultCustomizeModel().SetDbConn(conn).
// 		Find("test")
// 	if err != nil || item == nil {
// 		t.Error("該活動沒有自定義欄位資料")
// 	} else {
// 		t.Log("成功查詢該活動自定義欄位資料")
// 	}

// 	// 刪除活動資料
// 	db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// 	db.Table(config.ACTIVITY_CUSTOMIZE_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }
