package models

// 第一步: 測試更新用戶答題紀錄(redis)
// 第二步: 測試查詢功能(redis)
// 第三步: 測試將資料加入資料表
// 第四步: 測試查詢功能(資料表)

// 測試更新用戶答題紀錄(redis)
// func Test_Update_Game_QA_Record_Data(t *testing.T) {
// 	var (
// 		// activityID = "activity"
// 		gameID = "game"
// 		userID = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round  = 3
// 		// score      = 100
// 		editParam = EditGameQARecordModel{
// 			UserID:   userID,
// 			GameID:   gameID,
// 			QARound:  int64(round),
// 			QAOption: "A",
// 			// Score:       int64(round),
// 			OriginScore: int64(round),
// 			AddScore:    int64(round) * int64(round),
// 		}
// 	)

// 	// 更新用戶答題紀錄(redis)
// 	err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		Update(true, editParam)
// 	if err == nil {
// 		t.Log("更新用戶答題紀錄(redis)完成")
// 	} else if err != nil {
// 		t.Error("更新用戶答題紀錄(redis)發生問題: ", err)
// 	}

// 	// 先取得用戶答題紀錄(redis)
// 	record, err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		Find(true, gameID, userID)
// 	if err == nil {
// 		log.Println("record: ", record)
// 		t.Log("查詢用戶答題紀錄(redis)完成")
// 	} else if err != nil {
// 		t.Error("查詢用戶答題紀錄(redis)完成發生問題: ", err)
// 	}

// 	// redis.DelCache("QA_record_game")
// }

// // 測試查詢用戶答題紀錄(redis)
// func Test_Find_Game_QA_Record_By_Redis_Data(t *testing.T) {
// 	var (
// 		// activityID = "activity"
// 		gameID = "game"
// 		userID = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		// round      = 1
// 		// score      = 100
// 	)

// 	// 先查詢用戶答題紀錄(redis)，FindAllByRedis
// 	records, err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		FindAllByRedis(true, gameID)
// 	if err == nil {
// 		t.Log("查詢用戶答題紀錄(redis，FindAllByRedis)完成")
// 		log.Println("records: ", records)
// 	} else if err != nil {
// 		t.Error("查詢用戶答題紀錄(redis，FindAllByRedis)發生問題: ", err)
// 	}

// 	// 查詢用戶答題紀錄(redis)，Find
// 	record, err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		Find(true, gameID, userID)
// 	if err == nil {
// 		t.Log("查詢用戶答題紀錄(redis，Find)完成")
// 		log.Println("record: ", record)
// 	} else if err != nil {
// 		t.Error("查詢用戶答題紀錄(redis，Find)發生問題: ", err)
// 	}
// }

// // 測試將用戶答題紀錄加入activity_game_qa_record表
// func Test_Add_Game_QA_Record_Data(t *testing.T) {
// 	var (
// 		activityID = "activity"
// 		gameID     = "game"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = 1
// 		// score      = 100
// 	)

// 	// 先取得用戶答題紀錄(redis)
// 	record, _ := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		Find(true, gameID, userID)

// 	// 將答題紀錄加入activity_game_qa_record表
// 	err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		Add(userID, activityID, gameID, int64(round), record)
// 	if err == nil {
// 		t.Log("用戶答題紀錄加入activity_game_qa_record表完成")
// 	} else if err != nil {
// 		t.Error("用戶答題紀錄加入activity_game_qa_record表發生問題: ", err)
// 	}
// }

// // 測試查詢用戶答題紀錄(資料表)
// func Test_Find_Game_QA_Record_By_DB_Data(t *testing.T) {
// 	var (
// 		// activityID = "activity"
// 		gameID = "game"
// 		// userID     = "user"
// 		// round      = 1
// 		// score      = 100
// 	)

// 	// 取得所有用戶答題紀錄(資料表)
// 	records, err := DefaultGameQARecordModel().SetDbConn(conn).SetRedisConn(redis).
// 		FindAllByDB(gameID, 0)
// 	if err == nil {
// 		t.Log("查詢所有用戶答題紀錄(資料表，FindAllByDB)完成")
// 		log.Println("records: ", records)
// 	} else if err != nil {
// 		t.Error("查詢所有用戶答題紀錄(資料表，FindAllByDB)發生問題: ", err)
// 	}
// }

// // NewGameQARecordModel{
// // 	UserID:     userID,
// // 	ActivityID: activityID,
// // 	GameID:     gameID,
// // 	Round:      int64(round),
// // 	Score:      int64(score),

// // 	QA1Option:      record.QA1Option,
// // 	QA1OriginScore: record.QA1OriginScore,
// // 	QA1AddScore:    record.QA1AddScore,

// // 	QA2Option:      record.QA2Option,
// // 	QA2OriginScore: record.QA2OriginScore,
// // 	QA2AddScore:    record.QA2AddScore,

// // 	QA3Option:      record.QA3Option,
// // 	QA3OriginScore: record.QA3OriginScore,
// // 	QA3AddScore:    record.QA3AddScore,

// // 	QA4Option:      record.QA4Option,
// // 	QA4OriginScore: record.QA4OriginScore,
// // 	QA4AddScore:    record.QA4AddScore,

// // 	QA5Option:      record.QA5Option,
// // 	QA5OriginScore: record.QA5OriginScore,
// // 	QA5AddScore:    record.QA5AddScore,

// // 	QA6Option:      record.QA6Option,
// // 	QA6OriginScore: record.QA6OriginScore,
// // 	QA6AddScore:    record.QA6AddScore,

// // 	QA7Option:      record.QA7Option,
// // 	QA7OriginScore: record.QA7OriginScore,
// // 	QA7AddScore:    record.QA7AddScore,

// // 	QA8Option:      record.QA8Option,
// // 	QA8OriginScore: record.QA8OriginScore,
// // 	QA8AddScore:    record.QA8AddScore,

// // 	QA9Option:      record.QA9Option,
// // 	QA9OriginScore: record.QA9OriginScore,
// // 	QA9AddScore:    record.QA9AddScore,

// // 	QA10Option:      record.QA10Option,
// // 	QA10OriginScore: record.QA10OriginScore,
// // 	QA10AddScore:    record.QA10AddScore,

// // 	QA11Option:      record.QA11Option,
// // 	QA11OriginScore: record.QA11OriginScore,
// // 	QA11AddScore:    record.QA11AddScore,

// // 	QA12Option:      record.QA12Option,
// // 	QA12OriginScore: record.QA12OriginScore,
// // 	QA12AddScore:    record.QA12AddScore,

// // 	QA13Option:      record.QA13Option,
// // 	QA13OriginScore: record.QA13OriginScore,
// // 	QA13AddScore:    record.QA13AddScore,

// // 	QA14Option:      record.QA14Option,
// // 	QA14OriginScore: record.QA14OriginScore,
// // 	QA14AddScore:    record.QA14AddScore,

// // 	QA15Option:      record.QA15Option,
// // 	QA15OriginScore: record.QA15OriginScore,
// // 	QA15AddScore:    record.QA15AddScore,

// // 	QA16Option:      record.QA16Option,
// // 	QA16OriginScore: record.QA16OriginScore,
// // 	QA16AddScore:    record.QA16AddScore,

// // 	QA17Option:      record.QA17Option,
// // 	QA17OriginScore: record.QA17OriginScore,
// // 	QA17AddScore:    record.QA17AddScore,

// // 	QA18Option:      record.QA18Option,
// // 	QA18OriginScore: record.QA18OriginScore,
// // 	QA18AddScore:    record.QA18AddScore,

// // 	QA19Option:      record.QA19Option,
// // 	QA19OriginScore: record.QA19OriginScore,
// // 	QA19AddScore:    record.QA19AddScore,

// // 	QA20Option:      record.QA20Option,
// // 	QA20OriginScore: record.QA20OriginScore,
// // 	QA20AddScore:    record.QA20AddScore,
// // }
