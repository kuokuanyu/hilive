package models

// 測試查詢聊天紀錄資料功能(包含新增聊天紀錄功能)
// func Test_ChatroomRecord_LeftJoinLineUsers(t *testing.T) {
// 	var (
// 		activityID = "dQSrpxSJlexj6G9EIph4"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增三筆聊天紀錄測試資料
// 	for i := 0; i < 500; i++ {
// 		if err := DefaultChatroomRecordModel().SetDbConn(conn).
// 			Add(false, EditChatroomRecordModel{
// 				ActivityID:    activityID,
// 				UserID:        userID,
// 				MessageType:   "normal-barrage",
// 				MessageStyle:  "default",
// 				Message:       strconv.Itoa(i),
// 				MessageStatus: "yes",
// 				MessagePrice: "100",
// 			}); err != nil {
// 			t.Error("錯誤: 新增聊天紀錄發生問題", err)
// 		}
// 	}

// 	// records, err := DefaultChatroomRecordModel().SetDbConn(conn).FindAll(activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢聊天紀錄發生問題")
// 	// } else if len(records) != 3 {
// 	// 	t.Error("錯誤: 查詢聊天紀錄發生問題，聊天紀錄長度錯誤")
// 	// } else {
// 	// 	t.Log("查詢資料完成(LeftJoinLineUsers)")
// 	// 	fmt.Println(records)
// 	// }

// 	// 刪除測試資料
// 	// db.Table(config.ACTIVITY_CHATROOM_RECORD_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }
