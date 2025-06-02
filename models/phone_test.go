package models

// // 將users資料表所有資料加入user_phone資料表
// func Test_Add_Phone_Data(t *testing.T) {

// 	// 先取得users資料表所有資料
// 	items, err := db.Table("users").WithConn(conn).
// 		Select("phone").All()
// 	if err != nil {
// 		t.Error("查詢users所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢users所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			phone = utils.GetString(item["phone"], "")
// 		)

// 		// 新增電話資料
// 		err = DefaultPhoneModel().SetDbConn(conn).Add(
// 			EditPhoneModel{
// 				Phone:  phone,
// 				Status: "yes",
// 				Times:  "0",
// 			})
// 		if err != nil {
// 			t.Error("新增電話資料發生錯誤")
// 		}
// 	}
// }
