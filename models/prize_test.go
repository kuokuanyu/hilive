package models

// 利用redis取得獎品數
// func Test_123(t *testing.T) {
// 	// 刪除遊戲紀錄: del user_game_records_Bfon6SaV6ORhmuDQUioI
// 	// 刪除獎品資訊: del game_prizes_amount_bCIX7MJaFen3v4BtxP5E

// 	var (
// 		// game       = "redpack"
// 		activityID = "Bfon6SaV6ORhmuDQUioI"
// 		gameID     = "bCIX7MJaFen3v4BtxP5E"
// 	)
// 	// prizes, _ := redis.HashGetAllCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
// 	// fmt.Println(prizes)

// 	// 刪除中獎人員名單
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).
// 		Where("activity_id", "=", activityID).Where("game_id", "=", gameID).Delete()

// 	// 修改獎品數
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1104).Update(command.Value{"prize_remain": 1})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1105).Update(command.Value{"prize_remain": 2})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1106).Update(command.Value{"prize_remain": 3})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1107).Update(command.Value{"prize_remain": 4})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1108).Update(command.Value{"prize_remain": 5})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1109).Update(command.Value{"prize_remain": 6})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1110).Update(command.Value{"prize_remain": 7})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1111).Update(command.Value{"prize_remain": 8})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1112).Update(command.Value{"prize_remain": 9})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1113).Update(command.Value{"prize_remain": 10})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1114).Update(command.Value{"prize_remain": 11})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1115).Update(command.Value{"prize_remain": 12})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1116).Update(command.Value{"prize_remain": 13})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1117).Update(command.Value{"prize_remain": 14})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1118).Update(command.Value{"prize_remain": 15})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1119).Update(command.Value{"prize_remain": 16})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1120).Update(command.Value{"prize_remain": 17})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1121).Update(command.Value{"prize_remain": 18})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1122).Update(command.Value{"prize_remain": 19})
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("id", "=", 1123).Update(command.Value{"prize_remain": 110})

// 	// DefaultPrizeModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).FindPrizesAmount(true, gameID)

// 	// prizes, _ = redis.HashGetAllCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
// 	// fmt.Println(prizes)

// 	// 刪除redis資訊
// 	// _, err := redis.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
// 	// fmt.Println(err)

// 	// prizes, _ = redis.HashGetAllCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
// 	// fmt.Println(prizes)

// }

// // 利用redis取得獎品數
// func Test_Pieze_FindPrizesCount_Redis(t *testing.T) {
// 	var (
// 		// game       = "redpack"
// 		// activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID = "ZPmlZQCOhLd42GKHtpxZ"
// 	)

// 	count, err := DefaultPrizeModel().SetDbConn(conn).SetRedisConn(redis).
// 		FindPrizesAmount(true, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 獎品數量查詢發生問題(FindPrizesCount)", err)
// 	}
// 	fmt.Println("獎品數量: ", count)

// 	// 刪除redis資訊
// 	redis.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
// }

// // 測試獎品查詢功能(FindExistPrizes)
// func Test_Pieze_FindExistPrizes(t *testing.T) {
// 	var (
// 		game       = "redpack"
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "redpack"
// 	)

// 	prizes, err := DefaultPrizeModel().SetDbConn(conn).
// 		FindExistPrizes(game, activityID)
// 	if err != nil {
// 		t.Error("錯誤: 獎品查詢發生問題(FindExistPrizes)", err)
// 	}
// 	fmt.Println("prizes(activity_id): ", prizes)

// 	prizes, err = DefaultPrizeModel().SetDbConn(conn).
// 		FindPrizes(false, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 獎品查詢發生問題(FindExistPrizes)", err)
// 	}
// 	fmt.Println("prizes(game_id): ", prizes)
// }

// // 測試獎品查詢功能(FindPrizes)
// func Test_Pieze_FindPrizes(t *testing.T) {
// 	var (
// 		game       = "redpack"
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "redpack"
// 	)

// 	prizes, err := DefaultPrizeModel().SetDbConn(conn).
// 		FindExistPrizes(game, activityID)
// 	if err != nil {
// 		t.Error("錯誤: 獎品查詢發生問題(FindPrizes)", err)
// 	}
// 	fmt.Println("prizes(activity_id): ", prizes)

// 	prizes, err = DefaultPrizeModel().SetDbConn(conn).
// 		FindPrizes(false, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 獎品查詢發生問題(FindPrizes)", err)
// 	}
// 	fmt.Println("prizes(game_id): ", prizes)
// }

// // 測試獎品查詢功能(FindPrize)
// func Test_Pieze_FindPrize(t *testing.T) {
// 	var (
// 		prizeID = "redpack"
// 	)

// 	prize, err := DefaultPrizeModel().SetDbConn(conn).
// 		FindPrize(false, prizeID)
// 	if err != nil || prize.ID == 0 {
// 		t.Error("錯誤: 獎品查詢發生問題(FindPrize)", err)
// 	}
// 	fmt.Println("prize: ", prize)
// }

// // 測試獎品數量遞減功能
// func Test_Pieze_DecrRemain(t *testing.T) {
// 	var (
// 		gameID  = ""
// 		prizeID = "redpack"
// 	)

// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		DecrRemain(gameID, prizeID); err != nil {
// 		t.Error("錯誤: 獎品遞減發生問題", err)
// 	}
// }

// // 測試更新獎品資料
// func Test_Pieze_Update(t *testing.T) {
// 	var (
// 		prizeID = "redpack"
// 	)

// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Update(false, "redpack", EditPrizeModel{
// 			PrizeID:       prizeID,
// 			PrizeName:     "獎品更新",
// 			PrizeType:     "thanks",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "100",
// 			PrizePrice:    "500",
// 			PrizeMethod:   "mail",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 更新獎品資料發生問題", err)
// 	}
// }

// // 測試刪除測試獎品資料
// func Test_Prize_Delete(t *testing.T) {
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("prize_id", "=", "redpack").Delete()
// 	// db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("prize_id", "=", "whack_mole").Delete()
// }

// // 測試新增敲敲樂獎品資料
// // func Test_Pieze_Add_Whack_Mole(t *testing.T) {
// // 	var (
// // 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// // 		gameID     = "whack_mole"
// // 		prizeID    = "whack_mole"
// // 	)

// // 	if err := DefaultPrizeModel().SetDbConn(conn).
// // 		Add(false, "whack_mole", prizeID, NewPrizeModel{
// // 			ActivityID:    activityID,
// // 			GameID:        gameID,
// // 			PrizeName:     "敲敲樂",
// // 			PrizeType:     "third",
// // 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// // 			PrizeAmount:   "666",
// // 			PrizePrice:    "10654",
// // 			PrizeMethod:   "site",
// // 			PrizePassword: "password",
// // 		}); err != nil {
// // 		t.Error("錯誤: 新增獎品資料發生問題", err)
// // 	}
// // }

// // 測試新增搖紅包獎品資料
// // func Test_Pieze_Add_Redpack(t *testing.T) {
// // 	var (
// // 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// // 		gameID     = "redpack"
// // 		prizeID    = "redpack"
// // 	)

// // 	if err := DefaultPrizeModel().SetDbConn(conn).
// // 		Add(false, "redpack", prizeID, NewPrizeModel{
// // 			ActivityID:    activityID,
// // 			GameID:        gameID,
// // 			PrizeName:     "獎品",
// // 			PrizeType:     "money",
// // 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// // 			PrizeAmount:   "80",
// // 			PrizePrice:    "10000",
// // 			PrizeMethod:   "site",
// // 			PrizePassword: "password",
// // 		}); err != nil {
// // 		t.Error("錯誤: 新增獎品資料發生問題", err)
// // 	}
// // }

// // 測試將所有獎品資料加入activity_prize1表(更新至正式區時需要執行)
// func Test_Add_Prize_Data(t *testing.T) {
// 	// 搖紅包所有獎品
// 	items, err := db.Table("activity_redpack_prize").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("錯誤: 查詢所有搖紅包獎品資料發生問題")
// 	} else {
// 		t.Log("查詢所有搖紅包獎品資料完成")
// 	}

// 	// 將所有搖紅包獎品資料加入activity_prize1
// 	for _, item := range items {

// 		if _, err = db.Table("activity_prize1").WithConn(conn).Insert(
// 			command.Value{
// 				"activity_id":    utils.GetStringFromMap(item, "activity_id", ""),
// 				"game_id":        utils.GetStringFromMap(item, "game_id", ""),
// 				"prize_id":       utils.GetStringFromMap(item, "prize_id", ""),
// 				"game":           "redpack",
// 				"prize_name":     utils.GetStringFromMap(item, "name", "prize"),
// 				"prize_type":     "general",
// 				"prize_picture":  utils.GetStringFromMap(item, "picture", ""),
// 				"prize_amount":   utils.GetInt64FromMap(item, "amount", 0),
// 				"prize_remain":   utils.GetInt64FromMap(item, "remain", 0),
// 				"prize_price":    utils.GetInt64FromMap(item, "price", 0),
// 				"prize_method":   utils.GetStringFromMap(item, "method", ""),
// 				"prize_password": string(utils.Encode([]byte(utils.GetStringFromMap(item, "password", "password")))),
// 			}); err != nil {
// 			t.Error("錯誤: 加入搖紅包獎品資料至activity_prize1發生問題")
// 		}
// 	}
// 	fmt.Println("將所有搖紅包獎品新增至合併獎品資料表中完成")

// 	// 套紅包所有獎品
// 	items, err = db.Table("activity_ropepack_prize").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("錯誤: 查詢所有套紅包獎品資料發生問題")
// 	} else {
// 		t.Log("查詢所有套紅包獎品資料完成")
// 	}

// 	// 將所有套紅包獎品資料加入activity_prize1
// 	for _, item := range items {

// 		if _, err = db.Table("activity_prize1").WithConn(conn).Insert(
// 			command.Value{
// 				"activity_id":    utils.GetStringFromMap(item, "activity_id", ""),
// 				"game_id":        utils.GetStringFromMap(item, "game_id", ""),
// 				"prize_id":       utils.GetStringFromMap(item, "prize_id", ""),
// 				"game":           "ropepack",
// 				"prize_name":     utils.GetStringFromMap(item, "name", "prize"),
// 				"prize_type":     "general",
// 				"prize_picture":  utils.GetStringFromMap(item, "picture", ""),
// 				"prize_amount":   utils.GetInt64FromMap(item, "amount", 0),
// 				"prize_remain":   utils.GetInt64FromMap(item, "remain", 0),
// 				"prize_price":    utils.GetInt64FromMap(item, "price", 0),
// 				"prize_method":   utils.GetStringFromMap(item, "method", ""),
// 				"prize_password": string(utils.Encode([]byte(utils.GetStringFromMap(item, "password", "password")))),
// 			}); err != nil {
// 			t.Error("錯誤: 加入套紅包獎品資料至activity_prize1發生問題")
// 		}
// 	}
// 	fmt.Println("將所有套紅包獎品新增至合併獎品資料表中完成")

// 	// 敲敲樂所有獎品
// 	items, err = db.Table("activity_whack_mole_prize").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("錯誤: 查詢所有敲敲樂獎品資料發生問題")
// 	} else {
// 		t.Log("查詢所有敲敲樂獎品資料完成")
// 	}

// 	// 將所有敲敲樂獎品資料加入activity_prize1
// 	for _, item := range items {

// 		if _, err = db.Table("activity_prize1").WithConn(conn).Insert(
// 			command.Value{
// 				"activity_id":    utils.GetStringFromMap(item, "activity_id", ""),
// 				"game_id":        utils.GetStringFromMap(item, "game_id", ""),
// 				"prize_id":       utils.GetStringFromMap(item, "prize_id", ""),
// 				"game":           "whack_mole",
// 				"prize_name":     utils.GetStringFromMap(item, "name", "prize"),
// 				"prize_type":     "general",
// 				"prize_picture":  utils.GetStringFromMap(item, "picture", ""),
// 				"prize_amount":   utils.GetInt64FromMap(item, "amount", 0),
// 				"prize_remain":   utils.GetInt64FromMap(item, "remain", 0),
// 				"prize_price":    utils.GetInt64FromMap(item, "price", 0),
// 				"prize_method":   utils.GetStringFromMap(item, "method", ""),
// 				"prize_password": string(utils.Encode([]byte(utils.GetStringFromMap(item, "password", "password")))),
// 			}); err != nil {
// 			t.Error("錯誤: 加入敲敲樂獎品資料至activity_prize1發生問題")
// 		}
// 	}
// 	fmt.Println("將所有敲敲樂獎品新增至合併獎品資料表中完成")

// 	// 搖號抽獎所有獎品
// 	items, err = db.Table("activity_draw_numbers_prize").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("錯誤: 查詢所有搖號抽獎獎品資料發生問題")
// 	} else {
// 		t.Log("查詢所有搖號抽獎獎品資料完成")
// 	}

// 	// 將所有搖號抽獎獎品資料加入activity_prize1
// 	for _, item := range items {

// 		if _, err = db.Table("activity_prize1").WithConn(conn).Insert(
// 			command.Value{
// 				"activity_id":    utils.GetStringFromMap(item, "activity_id", ""),
// 				"game_id":        utils.GetStringFromMap(item, "game_id", ""),
// 				"prize_id":       utils.GetStringFromMap(item, "prize_id", ""),
// 				"game":           "draw_numbers",
// 				"prize_name":     utils.GetStringFromMap(item, "name", "prize"),
// 				"prize_type":     "general",
// 				"prize_picture":  utils.GetStringFromMap(item, "picture", ""),
// 				"prize_amount":   utils.GetInt64FromMap(item, "amount", 0),
// 				"prize_remain":   utils.GetInt64FromMap(item, "remain", 0),
// 				"prize_price":    utils.GetInt64FromMap(item, "price", 0),
// 				"prize_method":   utils.GetStringFromMap(item, "method", ""),
// 				"prize_password": string(utils.Encode([]byte(utils.GetStringFromMap(item, "password", "password")))),
// 			}); err != nil {
// 			t.Error("錯誤: 加入搖號抽獎獎品資料至activity_prize1發生問題")
// 		}
// 	}
// 	fmt.Println("將所有搖號抽獎獎品新增至合併獎品資料表中完成")

// 	// 遊戲抽獎所有獎品
// 	items, err = db.Table("activity_lottery_prize").WithConn(conn).All()
// 	if err != nil {
// 		t.Error("錯誤: 查詢所有遊戲抽獎獎品資料發生問題")
// 	} else {
// 		t.Log("查詢所有遊戲抽獎獎品資料完成")
// 	}

// 	// 將所有遊戲抽獎獎品資料加入activity_prize1
// 	for _, item := range items {

// 		if _, err = db.Table("activity_prize1").WithConn(conn).Insert(
// 			command.Value{
// 				"activity_id":    utils.GetStringFromMap(item, "activity_id", ""),
// 				"game_id":        utils.GetStringFromMap(item, "game_id", ""),
// 				"prize_id":       utils.GetStringFromMap(item, "prize_id", ""),
// 				"game":           "lottery",
// 				"prize_name":     utils.GetStringFromMap(item, "name", "prize"),
// 				"prize_type":     utils.GetStringFromMap(item, "prize_type", ""),
// 				"prize_picture":  utils.GetStringFromMap(item, "picture", ""),
// 				"prize_amount":   utils.GetInt64FromMap(item, "amount", 0),
// 				"prize_remain":   utils.GetInt64FromMap(item, "remain", 0),
// 				"prize_price":    utils.GetInt64FromMap(item, "price", 0),
// 				"prize_method":   utils.GetStringFromMap(item, "method", ""),
// 				"prize_password": string(utils.Encode([]byte(utils.GetStringFromMap(item, "password", "password")))),
// 			}); err != nil {
// 			t.Error("錯誤: 加入遊戲抽獎獎品資料至activity_prize1發生問題")
// 		}
// 	}
// 	fmt.Println("將所有遊戲抽獎獎品新增至合併獎品資料表")
// }
