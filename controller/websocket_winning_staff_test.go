package controller

// 測試敲敲樂發獎多線程功能(status = calculate)
// func Test_Calculate(t *testing.T) {
// 	var (
// 		game                           = "QA"
// 		activityID                     = "test"
// 		gameID                         = "QA"
// 		prize1, prize2, prize3, prize4 = "prize1", "prize2", "prize3", "prize4"
// 		round                          = 1
// 	)
// 	// 建立測試遊戲
// 	if err := models.DefaultGameModel().SetDbConn(conn).
// 		Add(true, game, gameID, models.NewGameModel{
// 			ActivityID:   activityID,
// 			Title:        "敲敲樂",
// 			GameType:     "",
// 			LimitTime:    "open",
// 			Second:       "16",
// 			MaxPeople:    "50",
// 			MaxTimes:     "0",
// 			Allow:        "open",
// 			Percent:      "0",
// 			FirstPrize:   "50",
// 			SecondPrize:  "50",
// 			ThirdPrize:   "100",
// 			GeneralPrize: "800",
// 			Topic:        "classic",
// 			Skin:         "rat",
// 			TotalQA:      "1",
// 		}); err != nil {
// 		t.Error("錯誤: 新增遊戲資料發生問題", err)
// 	}

// 	// 建立測試頭獎(50)、二獎(50)、三獎(100)、四獎(800)獎品
// 	if err := models.DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, game, prize1, models.NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "no.1",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "50",
// 			PrizePrice:    "10654",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增頭獎資料發生問題", err)
// 	}
// 	if err := models.DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, game, prize2, models.NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "no.2",
// 			PrizeType:     "second",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "50",
// 			PrizePrice:    "2222",
// 			PrizeMethod:   "mail",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增二獎資料發生問題", err)
// 	}
// 	if err := models.DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, game, prize3, models.NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "no.3",
// 			PrizeType:     "third",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "100",
// 			PrizePrice:    "1333",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增三獎資料發生問題", err)
// 	}
// 	if err := models.DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, game, prize4, models.NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "no.4",
// 			PrizeType:     "general",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "800",
// 			PrizePrice:    "400",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增普通獎資料發生問題", err)
// 	}

// 	// 新增測試分數資料
// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	wg.Add(1000)          //計數器
// 	for i := 1; i <= 1000; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			// redis(分數)
// 			if err := redis.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 				int64(i)); err != nil {
// 				t.Error("錯誤: 新增redis分數資料發生問題")
// 			}

// 			// redis(第2分數)
// 			if err := redis.ZSetAddFloat(config.SCORES_2_REDIS+gameID, strconv.Itoa(i),
// 				float64(i)+0.1); err != nil {
// 				t.Error("錯誤: 新增redis第二分數資料發生問題")
// 			}
// 		}(i)
// 	}
// 	wg.Wait()

// 	// 結算開始時間
// 	log.Println("遊戲結束，統計中獎人員資料")
// 	startTime := time.Now()

// 	// 將所有人員分數資料加入資料表中(分數大於0)
// 	staffs, err := models.DefaultScoreModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(true, gameID, int64(round), 0, "")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	wg.Add(len(staffs)) //計數器
// 	for i := 0; i < len(staffs); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			if err = models.DefaultScoreModel().SetDbConn(conn).SetRedisConn(redis).
// 				Add(models.NewScoreModel{
// 					ActivityID: activityID,
// 					GameID:     gameID,
// 					UserID:     staffs[i].UserID,
// 					Round:      int64(round),
// 					Score:      staffs[i].Score,
// 					Score2:     staffs[i].Score2,
// 				}); err != nil {
// 				t.Error("錯誤: 新增分數資料發生問題，請重新操作")
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0

// 	// 取得所有獎項設置的人數
// 	gameModel, err := models.DefaultGameModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindGame(true, gameID)
// 	if err != nil || gameModel.ID == 0 {
// 		t.Error("錯誤: 取得遊戲資訊資料發生問題")
// 	}

// 	var first, second, third, general, limit int64
// 	first, second, third, general = gameModel.FirstPrize, gameModel.SecondPrize, gameModel.ThirdPrize, gameModel.GeneralPrize
// 	limit = first + second + third + general
// 	fmt.Println("取得所有獎項設置的人數: ", first, ", ", second, ", ", third, ", ", general)

// 	var order string
// 	if game == "whack_mole" || game == "monopoly" {
// 		order = "desc" // 正確率由高至低
// 	} else if game == "QA" {
// 		order = "asc" // 耗時排序由低至高
// 	}
// 	// 取得資料表中用戶排名資料(大於0)
// 	scores, err := models.DefaultScoreModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(false, gameID, int64(round), limit, order)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	// 如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量
// 	if len(scores) < int(limit) {
// 		limit = int64(len(scores))
// 		if limit <= first {
// 			first, second, third, general = limit, 0, 0, 0
// 		} else if limit > first && limit <= first+second {
// 			second, third, general = limit-first, 0, 0
// 		} else if limit > first+second && limit <= first+second+third {
// 			third, general = limit-first-second, 0
// 		} else if limit > first+second+third && limit < first+second+third+general {
// 			general = limit - first - second - third
// 		}
// 	}

// 	var (
// 		allprizes           = make([]models.PrizeModel, 0)   // 發放獎品資訊
// 		allTypePrizes       = make([][]models.PrizeModel, 0) // 各獎項獎品資訊
// 		allTypePrizesAmount = make([]int64, 0)               // 各獎項獎品總數
// 		allPrizesPeople     = make([]int64, 0)               // 各獎項需要發獎的人數
// 	)
// 	allPrizesPeople = []int64{first, second, third, general}
// 	fmt.Println("如果玩遊戲人數少於所有獎項總和，取較低值並重新計算各個獎項需要發的獎品數量: ",
// 		first, ", ", second, ", ", third, ", ", general)

// 	// 取得各獎項所有資訊、數量
// 	prizeModel, err := models.DefaultPrizeModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindPrizes(false, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 查詢遊戲獎品資訊發生問題")
// 	}

// 	// 判斷頭獎、二獎、三獎、普通獎
// 	var (
// 		firstPrizes                                           = make([]models.PrizeModel, 0)
// 		secondPrizes                                          = make([]models.PrizeModel, 0)
// 		thirdPrizes                                           = make([]models.PrizeModel, 0)
// 		generalPrizes                                         = make([]models.PrizeModel, 0)
// 		firstAmount, secondAmount, thirdAmount, generalAmount int64
// 	)
// 	for i := 0; i < len(prizeModel); i++ {
// 		if prizeModel[i].PrizeRemain == 0 {
// 			continue
// 		}

// 		// 將獎項加入陣列中
// 		if prizeModel[i].PrizeType == "first" {
// 			firstPrizes = append(firstPrizes, prizeModel[i])
// 			firstAmount += prizeModel[i].PrizeRemain
// 		} else if prizeModel[i].PrizeType == "second" {
// 			secondPrizes = append(secondPrizes, prizeModel[i])
// 			secondAmount += prizeModel[i].PrizeRemain
// 		} else if prizeModel[i].PrizeType == "third" {
// 			thirdPrizes = append(thirdPrizes, prizeModel[i])
// 			thirdAmount += prizeModel[i].PrizeRemain
// 		} else if prizeModel[i].PrizeType == "general" {
// 			generalPrizes = append(generalPrizes, prizeModel[i])
// 			generalAmount += prizeModel[i].PrizeRemain
// 		}
// 	}

// 	allTypePrizes = append(allTypePrizes, firstPrizes, secondPrizes, thirdPrizes, generalPrizes)
// 	allTypePrizesAmount = append(allTypePrizesAmount, firstAmount, secondAmount, thirdAmount, generalAmount)
// 	fmt.Println("各獎項獎品: ", allTypePrizes, ", 各獎項獎品數: ", allTypePrizesAmount)

// 	// 判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值
// 	for i := 0; i < len(allPrizesPeople); i++ {
// 		if allTypePrizesAmount[i] < allPrizesPeople[i] {
// 			if i < len(allPrizesPeople)-1 {
// 				// 不夠發放，將多餘的值設置至下一個獎項
// 				allPrizesPeople[i+1] = allPrizesPeople[i+1] -
// 					allTypePrizesAmount[i] + allPrizesPeople[i]
// 			}
// 			allPrizesPeople[i] = allTypePrizesAmount[i]
// 		}
// 	}
// 	limit = allPrizesPeople[0] + allPrizesPeople[1] +
// 		allPrizesPeople[2] + allPrizesPeople[3]
// 	fmt.Println("判斷各獎項獎品數量是否足夠發放，如果不夠，取較低值: ",
// 		allPrizesPeople[0], ", ", allPrizesPeople[1], ", ", allPrizesPeople[2], ", ", allPrizesPeople[3])

// 	// 將需要發放的獎品都append allprizes
// 	for i := 0; i < len(allPrizesPeople); i++ {
// 		var stopValue int64 // 當len(allprizes)=stopValue，停止該獎項發放並發放下一個獎項
// 		for m := 0; m <= i; m++ {
// 			stopValue += allPrizesPeople[m]
// 		}

// 		for _, prize := range allTypePrizes[i] {
// 			for n := 0; n < int(prize.PrizeRemain); n++ {
// 				allprizes = append(allprizes, prize)
// 				if int64(len(allprizes)) == stopValue {
// 					break
// 				}
// 			}
// 			if int64(len(allprizes)) == stopValue {
// 				break
// 			}
// 		}
// 	}

// 	// 所有獎品資訊
// 	prizes, err := models.DefaultPrizeModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindPrizes(true, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得獎品資訊，請重新操作")
// 	}

// 	// 發放獎品(多線程發放)
// 	var (
// 		winRecords = make([]models.PrizeStaffModel, limit)
// 		// wg          sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	)
// 	wg.Add(int(limit)) //計數器
// 	for i := 0; i < int(limit); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			// 遞減獎品數量
// 			// if err = models.DefaultPrizeModel().SetDbConn(conn).
// 			// 	DecrRemain(allprizes[i].PrizeID); err != nil {
// 			// 	// 獎品無庫存
// 			// 	t.Error("錯誤: 獎品無庫存")
// 			// }

// 			// 遞減redis獎品數量
// 			redis.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, allprizes[i].PrizeID)

// 			// 新增中獎人員名單
// 			id, err := models.DefaultPrizeStaffModel().SetDbConn(conn).
// 				Add(models.NewPrizeStaffModel{
// 					ActivityID: scores[i].ActivityID,
// 					GameID:     gameID,
// 					UserID:     scores[i].UserID,
// 					PrizeID:    allprizes[i].PrizeID,
// 					Round:      strconv.Itoa(int(round)),
// 					Status:     "no",
// 					// White:      "no",
// 				})
// 			if err != nil {
// 				t.Error("錯誤: 新增中獎人員發生問題，請重新操作")
// 			}

// 			winRecords[i] = models.PrizeStaffModel{
// 				ID:            id,
// 				UserID:        scores[i].UserID,
// 				ActivityID:    activityID,
// 				GameID:        gameID,
// 				PrizeID:       allprizes[i].PrizeID,
// 				Name:          scores[i].Name,
// 				Avatar:        scores[i].Avatar,
// 				PrizeName:     allprizes[i].PrizeName,
// 				PrizeType:     allprizes[i].PrizeType,
// 				PrizePicture:  allprizes[i].PrizePicture,
// 				PrizePrice:    allprizes[i].PrizePrice,
// 				PrizeMethod:   allprizes[i].PrizeMethod,
// 				PrizePassword: allprizes[i].PrizePassword,
// 				Round:         int64(round),
// 				// WinTime:    now.Format("2006-01-02 15:04:05"),
// 				Status: "no",
// 				Score:  scores[i].Score,
// 				Score2: scores[i].Score2,
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0

// 	// 從獎品redis(hash格式)中取得所有獎品剩餘數量，並更新至資料表中
// 	var (
// 		amount int64
// 		lock   sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
// 	)
// 	// 所有獎品資訊
// 	prizes, err = models.DefaultPrizeModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindPrizes(true, gameID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得獎品資訊，請重新操作")
// 	}
// 	wg.Add(len(prizes)) //計數器
// 	for i := 0; i < len(prizes); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			// 計算剩餘獎品數
// 			lock.Lock() //佔有資源
// 			amount += prizes[i].PrizeRemain
// 			lock.Unlock() //釋放資源
// 			// 更新獎品剩餘數量
// 			if err = models.DefaultPrizeModel().SetDbConn(conn).UpdateRemain(
// 				prizes[i].PrizeID, prizes[i].PrizeRemain); err != nil {
// 				t.Error(err.Error())
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0
// 	fmt.Println("中獎人員: ", len(winRecords) == 1000)
// 	fmt.Println("剩餘獎品: ", amount == 0)

// 	endTime := time.Now()
// 	log.Println("統計中獎人員花費時間: ", endTime.Sub(startTime))

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_GAME_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_GAME_QA_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	// 清除redis測試資料
// 	redis.DelCache(config.SCORES_REDIS + gameID)   // 分數
// 	redis.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
// 	redis.DelCache(config.GAME_REDIS + gameID)     // 遊戲資訊
// 	redis.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)   // 獎品資訊
// }

// // 分數前n高人員
// // var staffs = make([]models.PrizeStaffModel, 0)
// // scoreStaffs, err := models.DefaultScoreModel().SetDbConn(conn).
// // 	SetRedisConn(redis).Find(true, gameID, limit, order)
// // if err != nil {
// // 	t.Error("錯誤: 無法取得分數前n高的人員資料")
// // }
// // for _, score := range scoreStaffs {
// // 	staffs = append(staffs, models.PrizeStaffModel{
// // 		// ID:         score.ID,
// // 		// ActivityID: score.ActivityID,
// // 		// GameID:     score.GameID,
// // 		UserID: score.UserID,
// // 		Name:   score.Name,
// // 		Avatar: score.Avatar,
// // 		// Round:      score.Round,
// // 		Score: score.Score,
// // 	})
// // }

// // var (
// // 	// prizeStaffs         = make([]models.PrizeStaffModel, 0) // 中獎人員(可能因獎品不夠，導致有些人員是沒有中獎的)
// // 	allprizes           = make([]models.PrizeModel, 0)   // 發放獎品資訊
// // 	allTypePrizes       = make([][]models.PrizeModel, 0) // 各獎項獎品資訊
// // 	allTypePrizesAmount = make([]int64, 0)               // 各獎項獎品總數
// // 	allPrizesPeople     = make([]int64, 0)               // 各獎項需要發獎的人數
// // 	// now, _              = time.ParseInLocation("2006-01-02 15:04:05",
// // 	// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// // 	first, second, third, general, limit int64
// // )

// // prizeStaffs, err = models.DefaultPrizeStaffModel().SetDbConn(conn).
// // 	FindWinningRecordsOrderByScore(gameID, strconv.Itoa(round))
// // if err != nil {
// // 	t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(輪次查詢)", err)
// // }
// // for i, prize := range prizeStaffs {
// // 	fmt.Println("排名", i+1, ": ", prize)
// // 	fmt.Println("分數: ", prize.Score)
// // }

// // 結算後的獎品數(從redis取得，如果沒有才執行資料表查詢)
// // amount, err := models.DefaultPrizeModel().SetDbConn(conn).
// // 	SetRedisConn(redis).FindPrizesAmount(true, gameID)
// // if err != nil {
// // 	 t.Error("錯誤: 無法取得獎品數量資訊，請重新操作")
// // }

// // 	if _, err := db.Table("activity_score").WithConn(conn).
// // 		Insert(command.Value{
// // 			"activity_id": activityID,
// // 			"game_id":     gameID,
// // 			"user_id":     i,
// // 			"round":       round,
// // 			"score":       i,
// // 		}); err != nil {
// // 		t.Error("錯誤: 新增計分資料發生問題")
// // 	}
// // }(i)
