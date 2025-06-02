package models

// 測試利用redis處理提問牆用戶端的提問、按讚、收回讚功能
// func Test_QuestionUser_Guest_Redis(t *testing.T) {
// 	var (
// 		activityID = "activity"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		id1, id2   int64
// 		err        error
// 	)

// 	// 查詢提問資料
// 	// questionsByTime, questionsByLikes, err := DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("查詢資料表並因為資料數為0所以不會執行加入redis功能")
// 	// log.Println("第一次測試結束-------------------------------------------")

// 	// 新增提問測試資料
// 	id1, err = DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).Add(true, EditQuestionUserModel{
// 		ActivityID: activityID, UserID: userID, Message: "測試", MessageStatus: "yes",
// 	})
// 	if err != nil {
// 		t.Error("錯誤: 新增提問資料發生問題")
// 	}
// 	log.Println("新增提問資料並將資料加入redis中")
// 	log.Println("第二次測試結束-------------------------------------------", id1)

// 	// 查詢提問資料
// 	// questionsByTime, questionsByLikes, err := DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("查詢資料表並因為資料數不為0所以會執行加入redis功能")
// 	// log.Println("questionsByTime: ", questionsByTime,
// 	// 	",questionsByLikes: ", questionsByLikes)
// 	// log.Println("第三次測試結束-------------------------------------------")

// 	// 再新增提問測試資料
// 	time.Sleep(1 * time.Second)
// 	id2, err = DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).Add(true, EditQuestionUserModel{
// 		ActivityID: activityID, UserID: userID, Message: "測試2", MessageStatus: "yes",
// 	})
// 	if err != nil {
// 		t.Error("錯誤: 新增提問資料發生問題")
// 	}
// 	log.Println("再次新增提問資料並將資料加入redis中")
// 	log.Println("第四次測試結束-------------------------------------------", id2)

// 	// 查詢用戶按讚紀錄
// 	// _, err = DefaultQuestionLikesRecordModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID, userID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢用戶按讚紀錄資料發生問題")
// 	// }
// 	// log.Println("會查詢資料表並因為資料數為0所以不會執行加入redis功能")
// 	// log.Println("第五次測試結束-------------------------------------------")

// 	// 測試用戶按讚功能
// 	// 新增用戶按讚紀錄
// 	err = DefaultQuestionLikesRecordModel().
// 		SetDbConn(conn).SetRedisConn(redis).Add(true,
// 		NewQuestionLikesRecordModel{
// 			ActivityID: activityID, QuestionID: id1, UserID: userID,
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增按讚紀錄資料發生問題")
// 	}
// 	// 遞增提問讚數資料
// 	err = DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).UpdateGuestLikes(true, id1,
// 		activityID, "like")
// 	if err != nil {
// 		t.Error("錯誤: 遞增提問讚數資料發生問題")
// 	}
// 	log.Println("新增用戶按讚紀錄資料並執行加入redis功能，遞增提問讚數資料")
// 	log.Println("第六次測試結束-------------------------------------------")
// 	// 查詢用戶按讚紀錄
// 	records, err := DefaultQuestionLikesRecordModel().
// 		SetDbConn(conn).SetRedisConn(redis).Find(true, activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢用戶按讚紀錄資料發生問題")
// 	}
// 	log.Println("不會查詢資料表並直接執行redis查詢")
// 	log.Println("records: ", records)
// 	log.Println("第七次測試結束-------------------------------------------")

// 	// 第二次新增用戶按讚紀錄
// 	err = DefaultQuestionLikesRecordModel().
// 		SetDbConn(conn).SetRedisConn(redis).Add(true,
// 		NewQuestionLikesRecordModel{
// 			ActivityID: activityID, QuestionID: id2, UserID: userID,
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增按讚紀錄資料發生問題")
// 	}
// 	// 遞增提問讚數資料
// 	err = DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).UpdateGuestLikes(true, id2,
// 		activityID, "like")
// 	if err != nil {
// 		t.Error("錯誤: 遞增提問讚數資料發生問題")
// 	}
// 	log.Println("新增用戶按讚紀錄資料並執行加入redis功能，遞增提問讚數資料")
// 	log.Println("第八次測試結束-------------------------------------------")
// 	// // 查詢用戶按讚紀錄
// 	// records, err = DefaultQuestionLikesRecordModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID, userID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢用戶按讚紀錄資料發生問題")
// 	// }
// 	// log.Println("不會查詢資料表並直接redis查詢")
// 	// log.Println("records: ", records)
// 	// log.Println("第九次測試結束-------------------------------------------")

// 	// // 查詢提問資料
// 	// questionsByTime, questionsByLikes, err = DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("從redis查詢提問資料", questionsByTime[0].Like == "no", questionsByTime[1].Like == "no",
// 	// 	questionsByTime[0].Likes == 1, questionsByTime[1].Likes == 1)
// 	// log.Println("questionsByTime: ", questionsByTime,
// 	// 	",questionsByLikes: ", questionsByLikes)
// 	// log.Println("第十次測試結束-------------------------------------------")

// 	// 測試收回讚
// 	// 刪除用戶按讚紀錄
// 	err = DefaultQuestionLikesRecordModel().
// 		SetDbConn(conn).SetRedisConn(redis).Delete(true,
// 		id1, activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 刪除用戶按讚紀錄資料發生問題")
// 	}
// 	err = DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).UpdateGuestLikes(true, id1,
// 		activityID, "unlike")
// 	if err != nil {
// 		t.Error("錯誤: 遞減提問讚數資料發生問題")
// 	}
// 	log.Println("刪除用戶按讚紀錄資料並執行加入redis功能，遞減提問讚數資料")
// 	log.Println("第十一次測試結束-------------------------------------------")

// 	// // 查詢用戶按讚紀錄
// 	// records, err = DefaultQuestionLikesRecordModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID, userID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢用戶按讚紀錄資料發生問題")
// 	// }
// 	// log.Println("不會查詢資料表並直接執行redis查詢", len(records) == 1)
// 	// log.Println("records: ", records)
// 	// log.Println("第十二次測試結束-------------------------------------------")

// 	// 查詢提問資料
// 	questionsByTime, questionsByLikes, err := DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢提問資料發生問題")
// 	}
// 	log.Println("從redis查詢提問資料", questionsByTime[0].Like == "no", questionsByTime[1].Like == "no",
// 		questionsByTime[0].Likes == 1, questionsByTime[1].Likes == 0)
// 	log.Println("questionsByTime: ", questionsByTime,
// 		",questionsByLikes: ", questionsByLikes)
// 	log.Println("第十三次測試結束-------------------------------------------")

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_QUESTION_USER_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	db.Table(config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	// 刪除緩存資料
// 	redis.DelCache(config.QUESTION_REDIS + activityID)                   // 提問資訊(所有提問資料)
// 	redis.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID)     // 提問資訊(時間)
// 	redis.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + activityID)    // 提問資訊(讚數)
// 	redis.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + activityID) // 提問資訊(用戶按讚紀錄)
// }

// 測試利用redis處理提問牆主持端的按讚、收回讚功能
// func Test_QuestionUser_Host_Redis(t *testing.T) {
// 	var (
// 		activityID = "activity"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 查詢提問資料
// 	// questionsByTime, questionsByLikes, err := DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("查詢資料表並因為資料數為0所以不會執行加入redis功能")
// 	// log.Println("questionsByTime: ", questionsByTime,
// 	// 	",questionsByLikes: ", questionsByLikes)
// 	// log.Println("第一次測試結束-------------------------------------------")

// 	// 新增提問測試資料
// 	// for i := 0; i < 3; i++ {
// 	time.Sleep(1 * time.Second)
// 	id, _ := db.Table(config.ACTIVITY_QUESTION_USER_TABLE).WithConn(conn).
// 		Insert(command.Value{
// 			"activity_id": activityID,
// 			"user_id":     userID,
// 			"content":     "測試: ",
// 			"likes":       0,
// 			"like":        "no",
// 		})
// 	// }

// 	// 查詢提問資料
// 	questionsByTime, questionsByLikes, err := DefaultQuestionUserModel().
// 		SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢提問資料發生問題")
// 	}
// 	log.Println("會查詢資料表並因為資料數不為0所以會將資料加入redis中 ")
// 	log.Println("questionsByTime: ", questionsByTime,
// 		",questionsByLikes: ", questionsByLikes)
// 	log.Println("第二次測試結束-------------------------------------------")

// 	// 查詢提問資料
// 	// questionsByTime, questionsByLikes, err = DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("不會查詢資料表並且從redis查詢提問資料")
// 	// log.Println("questionsByTime: ", questionsByTime,
// 	// 	",questionsByLikes: ", questionsByLikes)
// 	// log.Println("第三次測試結束-------------------------------------------")

// 	// 測試主持人按讚功能
// 	err = DefaultQuestionUserModel().SetDbConn(conn).SetRedisConn(redis).
// 		UpdateHostLikes(true, id, activityID, "yes")
// 	if err != nil {
// 		t.Error("錯誤: 更新按讚資料發生問題")
// 	}
// 	// 查詢按讚後的提問資料是否正確
// 	// questionsByTime, questionsByLikes, err = DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("從redis查詢提問資料並檢查redis資料是否更新: ",
// 	// 	questionsByTime[0].Like == "yes",
// 	// 	questionsByLikes[0].Likes == 1)
// 	// log.Println("questionsByTime: ", questionsByTime[0],
// 	// 	",questionsByLikes: ", questionsByLikes[0])
// 	log.Println("第四次測試結束(按讚)-------------------------------------------")

// 	// 測試主持人收回讚功能
// 	err = DefaultQuestionUserModel().SetDbConn(conn).SetRedisConn(redis).
// 		UpdateHostLikes(true, id, activityID, "no")
// 	if err != nil {
// 		t.Error("錯誤: 更新收回讚資料發生問題")
// 	}
// 	// 查詢按讚後的提問資料是否正確
// 	// questionsByTime, questionsByLikes, err = DefaultQuestionUserModel().
// 	// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 查詢提問資料發生問題")
// 	// }
// 	// log.Println("從redis查詢提問資料並檢查redis資料是否更新: ",
// 	// 	questionsByTime[0].Like == "no",
// 	// 	questionsByLikes[0].Likes == 0)
// 	// log.Println("questionsByTime: ", questionsByTime[0],
// 	// 	",questionsByLikes: ", questionsByLikes[0])
// 	log.Println("第五次測試結束(收回讚)-------------------------------------------")

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_QUESTION_USER_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// 	// 刪除緩存資料
// 	redis.DelCache(config.QUESTION_REDIS + activityID)                   // 提問資訊(時間)
// 	redis.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID)     // 提問資訊(所有提問資料)
// 	redis.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + activityID)    // 提問資訊(讚數)
// 	redis.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + activityID) // 提問資訊(用戶按讚紀錄)
// }

// 測試按讚紀錄所有查詢資料功能(包含新增按讚紀錄功能)
// func Test_QuestionLikesRecord_LeftJoinLineUsers(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增三筆按讚紀錄測試資料
// 	for i := 0; i < 3; i++ {
// 		if err := DefaultQuestionLikesRecordModel().SetDbConn(conn).
// 			Add(NewQuestionLikesRecordModel{
// 				ActivityID: activityID,
// 				UserID:     userID,
// 				QuestionID: int64(i),
// 			}); err != nil {
// 			t.Error("錯誤: 新增按讚紀錄發生問題")
// 		} else {
// 			t.Log("新增按讚紀錄完成")
// 		}
// 	}

// 	// LeftJoinLineUsersOrderByTime(時間排列)
// 	records, err := DefaultQuestionLikesRecordModel().SetDbConn(conn).
// 		LeftJoinLineUsers(activityID, userID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢按讚紀錄發生問題")
// 	} else {
// 		t.Log("查詢資料完成(LeftJoinLineUsers)")
// 		fmt.Println("按讚紀錄: ", records)
// 	}

// 	fmt.Println("按讚紀錄查詢功能正常執行")

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }
