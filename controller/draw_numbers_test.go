package controller

// 測試搖號抽獎中獎判斷(測試允許以及不允許重複中獎下的判斷)
// func Test_GetDrawNumbers(t *testing.T) {
// var (
// 	game       = "draw_numbers"
// 	userID     = "admin"
// 	activityID = "activity"
// 	gameID     = "game"
// 	prizeID    = "prize"
// 	peopleStr  = "100"
// 	allow      = "close"
// 	staffs     = make([]models.ApplysignModel, 0) // 簽到人員資訊
// 	// gameModel  = h.getGameInfo(gameID)            // 遊戲資訊
// 	numbers = make([]int64, 0)
// 	params  = []interface{}{config.WINNING_STAFFS_REDIS + gameID} // redis參數(不允許重複中獎時需要添加中獎人員進redis中)
// 	wg      sync.WaitGroup
// 	lock    sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
// )

// // 新增測試活動
// if _, err := db.Table(config.ACTIVITY_TABLE).WithConn(conn).Insert(command.Value{
// 	"activity_id":   activityID,
// 	"user_id":       userID,
// 	"activity_name": activityID,
// 	"activity_type": "",
// 	"people":        200,
// 	"attend":        200,
// 	"city":          "",
// 	"town":          "",
// 	"start_time":    "2022-01-01 00:00:00",
// 	"end_time":      "2022-12-31 00:00:00",
// 	"number":        1,
// 	// 活動總覽
// 	"overview_message":       "open",
// 	"overview_topic":         "open",
// 	"overview_question":      "open",
// 	"overview_danmu":         "open",
// 	"overview_special_danmu": "open",
// 	"overview_picture":       "close",
// 	"overview_holdscreen":    "open",
// 	"overview_general":       "open",
// 	"overview_threed":        "close",
// 	"overview_countdown":     "close",
// 	"overview_lottery":       "open",
// 	"overview_redpack":       "open",
// 	"overview_ropepack":      "open",
// 	"overview_whack_mole":    "open",
// 	"overview_draw_numbers":  "open",
// 	"overview_monopoly":      "open",
// 	"overview_qa":            "open",
// 	// 活動介紹
// 	"introduce_title": "活動介紹",
// 	// 活動行程
// 	"schedule_title":          "活動行程",
// 	"schedule_display_date":   "open",
// 	"schedule_display_detail": "open",
// 	// 活動嘉賓
// 	"guest_title":      "活動嘉賓",
// 	"guest_background": config.UPLOAD_SYSTEM_URL + "img-guest-bg.png",
// 	// 活動資料
// 	"material_title": "活動資料",
// 	// 報名
// 	"apply_check": "open",
// 	// 簽到
// 	"sign_check":   "open",
// 	"sign_allow":   "close",
// 	"sign_minutes": 30,
// 	"sign_manual":  "open",
// 	// 訊息牆
// 	"message_picture":                  "close",
// 	"message_auto":                     "close",
// 	"message_ban":                      "close",
// 	"message_ban_second":               1,
// 	"message_refresh_second":           20,
// 	"message_open":                     "open",
// 	"message":                          "歡迎來到活動聊天室!",
// 	"message_background":               config.UPLOAD_SYSTEM_URL + "img-topic-bg.png",
// 	"message_text_color":               "#404b79",
// 	"message_screen_title_color":       "#ffffff",
// 	"message_text_frame_color":         "#81869b",
// 	"message_frame_color":              "#e7eaf4",
// 	"message_screen_title_frame_color": "#81869b",
// 	// 主題牆
// 	"topic_background": config.UPLOAD_SYSTEM_URL + "img-topic-bg.png",
// 	// 提問牆
// 	// "question_message_check":  "close",
// 	"question_anonymous": "close",
// 	// "question_hide_answered":  "close",
// 	"question_qrcode":     "open",
// 	"question_background": config.UPLOAD_SYSTEM_URL + "img-request-bg.png",
// 	// "question_guest_answered": "close",
// 	// 彈幕
// 	"danmu_loop":           "close",
// 	"danmu_top":            "open",
// 	"danmu_mid":            "open",
// 	"danmu_bottom":         "open",
// 	"danmu_display_name":   "open",
// 	"danmu_display_avatar": "open",
// 	"danmu_size":           100,
// 	"danmu_speed":          100,
// 	"danmu_density":        100,
// 	"danmu_opacity":        100,
// 	"danmu_background":     1,
// 	// 特殊彈幕
// 	"special_danmu_message_check": "close",
// 	"special_danmu_general_price": 0,
// 	"special_danmu_large_price":   0,
// 	"special_danmu_over_price":    0,
// 	"special_danmu_topic":         "open",
// 	// 圖片牆
// 	"picture_start_time":    "2022-01-01 00:00:00",
// 	"picture_end_time":      "2022-12-31 00:00:00",
// 	"picture_hide_time":     "close",
// 	"picture_switch_second": 3,
// 	"picture_play_order":    "order",
// 	"picture":               "system/img-pic.png",
// 	"picture_background":    config.UPLOAD_SYSTEM_URL + "img-picture-bg.png",
// 	"picture_animation":     1,
// 	// 霸屏
// 	"holdscreen_price":          0,
// 	"holdscreen_message_check":  "close",
// 	"holdscreen_only_picture":   "close",
// 	"holdscreen_display_second": 10,
// 	"holdscreen_birthday_topic": "open",
// 	// 一般簽到
// 	"general_display_people": "open",
// 	"general_style":          1,
// 	"general_background":     config.UPLOAD_SYSTEM_URL + "img-sign-bg.png",
// 	// 立體簽到
// 	"threed_avatar":           config.UPLOAD_SYSTEM_URL + "img-sign3d-headpic.png",
// 	"threed_avatar_shape":     "circle",
// 	"threed_display_people":   "open",
// 	"threed_background_style": 1,
// 	"threed_background":       config.UPLOAD_SYSTEM_URL + "img-sign3d-bg.png",
// 	"threed_image":            config.UPLOAD_SYSTEM_URL + "img-sign3d-design.png",
// 	"threed_image_size":       50,
// 	// 倒數計時
// 	"countdown_second":       5,
// 	"countdown_url":          "current",
// 	"countdown_avatar":       config.UPLOAD_SYSTEM_URL + "img-countdown-headpic.png",
// 	"countdown_avatar_shape": "circle",
// 	"countdown_background":   1,
// }); err != nil {
// 	t.Error("錯誤: 無法新增活動資料")
// }

// // 新增測試報名簽到人員資料(user1~200)
// for i := 1; i <= 200; i++ {
// 	if _, err := db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Insert(command.Value{
// 		"user_id":     i,
// 		"activity_id": activityID,
// 		"status":      "sign",
// 		"number":      i,
// 	}); err != nil {
// 		t.Error("錯誤: 無法新增報名簽到人員資料")
// 	}
// }

// // 新增測試搖號抽獎遊戲(不允許重複中獎)
// if err := models.DefaultGameModel().SetDbConn(conn).
// 	Add(game, gameID, models.NewGameModel{
// 		ActivityID:   activityID,
// 		Title:        "搖號抽獎",
// 		GameType:     "",
// 		LimitTime:    "open",
// 		Second:       "16",
// 		MaxPeople:    "0",
// 		MaxTimes:     "0",
// 		Allow:        allow,
// 		Percent:      "0",
// 		FirstPrize:   "0",
// 		SecondPrize:  "0",
// 		ThirdPrize:   "0",
// 		GeneralPrize: "0",
// 		Topic:        "classic",
// 		Skin:         "",
// 	}); err != nil {
// 	t.Error("錯誤: 新增遊戲資料發生問題", err)
// }

// // 新增測試獎品資訊(100個獎品)
// if err := models.DefaultPrizeModel().SetDbConn(conn).
// 	Add(false, game, prizeID, models.NewPrizeModel{
// 		ActivityID:    activityID,
// 		GameID:        gameID,
// 		PrizeName:     "no.1",
// 		PrizeType:     "first",
// 		PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 		PrizeAmount:   "100",
// 		PrizePrice:    "10654",
// 		PrizeMethod:   "site",
// 		PrizePassword: "password",
// 	}); err != nil {
// 	t.Error("錯誤: 新增獎品資料發生問題", err)
// }

// // 新增測試中獎人員資料(第101~200用戶中獎)
// for i := 101; i <= 200; i++ {
// 	_, err := models.DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(models.NewPrizeStaffModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     strconv.Itoa(i),
// 			PrizeID:    prizeID,
// 			Round:      "1",
// 			Status:     "no",
// 			// White:      "no",
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增中獎人員發生問題")
// 	}
// }

// // 取得遊戲資訊
// gameModel, err := models.DefaultGameModel().SetDbConn(conn).
// 	SetRedisConn(redis).FindGame(true, gameID)
// if err != nil || gameModel.ID == 0 {
// 	t.Error("錯誤: 取得遊戲資訊資料發生問題")
// }

// // 活動資訊
// activityModel, _ := models.DefaultActivityModel().
// 	SetDbConn(conn).SetRedisConn(redis).Find(true, activityID)

// // 搖號抽獎遊戲需要判斷是否重複中獎，如果不允許重複中獎則需要處理可抽獎人數資訊
// if gameModel.Allow == "open" {
// 	// 允許重複中獎，取得目前活動的簽到人數
// 	log.Println("可以抽獎的人數: ", activityModel.Attend)
// } else if gameModel.Allow == "close" {
// 	// 不允許重複中獎
// 	prizeStaffs, _ := models.DefaultPrizeStaffModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindWinningStaffs(true, gameID) // 中獎人員資料

// 	log.Println("中獎人員是否為100人? ", len(prizeStaffs) == 100)

// 	// 扣除已中獎人員
// 	log.Println("扣除已中獎人員，可以抽獎的人數: ", activityModel.Attend-int64(len(prizeStaffs)),
// 		activityModel.Attend-int64(len(prizeStaffs)) == 100)
// }

// startTime := time.Now()
// people, err := strconv.Atoi(peopleStr)
// if err != nil || people == 0 {
// 	t.Error("錯誤: 無法辨識中獎人數資訊，請輸入有效的中獎人數")

// }

// // 獎品資訊
// prizeModel, err := models.DefaultPrizeModel().SetDbConn(conn).
// 	FindPrize(prizeID)
// if err != nil {
// 	t.Error(err.Error())
// }

// // 簽到人員資訊
// if staffs, err = models.DefaultApplysignModel().SetDbConn(conn).
// 	SetRedisConn(redis).FindSignStaffs(false, activityID); err != nil {
// 	t.Error(err.Error())
// }
// log.Println("一開始取得的簽到人員人數，是否為200人: ", len(staffs) == 200)

// // 不允許重複中獎，將簽到人員裡已中獎人員資料去除
// if gameModel.Allow == "close" {
// 	newStaffs := make([]models.ApplysignModel, 0) // 簽到人員資訊(整理過後)
// 	for _, staff := range staffs {
// 		if !redis.SetIsMember(config.WINNING_STAFFS_REDIS+gameID, staff.UserID) {
// 			// 未中獎人員，加入抽獎陣列中
// 			newStaffs = append(newStaffs, staff)
// 		}
// 	}
// 	staffs = newStaffs
// 	log.Println("不允許重複中獎，處理簽到人員後的抽獎人數，是否剩100人: ", len(staffs) == 100)
// 	// log.Println("可抽獎的人員: ", staffs)
// }

// // 抽到的index值(不是抽中的號碼，而是先抽報名簽到人員陣列中的index值)
// for i := 0; len(numbers) < people; i++ {
// 	// 隨機值
// 	rand.Seed(time.Now().UnixNano())
// 	random := rand.Intn(len(staffs))

// 	// 判斷是否已經抽過該值，如果沒抽過則加入陣列中
// 	if !utils.IntInArray(numbers, int64(random)) {
// 		numbers = append(numbers, int64(random))
// 	}
// }

// // 發獎
// // 多線呈將資料加入資料表
// wg.Add(len(numbers)) //計數器
// for i := 0; i < len(numbers); i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		staff := staffs[numbers[i]] // 中獎人員

// 		// 新增中獎人員名單
// 		if _, err = models.DefaultPrizeStaffModel().SetDbConn(conn).
// 			Add(models.NewPrizeStaffModel{
// 				UserID:     staff.UserID,
// 				ActivityID: activityID,
// 				GameID:     gameID,
// 				PrizeID:    prizeID,
// 				Round:      strconv.Itoa(int(gameModel.GameRound)),
// 				Status:     "no",
// 				// White:      "no",
// 			}); err != nil {
// 			t.Error("錯誤: 新增中獎人員發生問題，請重新操作")
// 		}

// 		lock.Lock()                           //佔有資源
// 		numbers[i] = staff.Number             // 抽到的號碼，更新至陣列中
// 		params = append(params, staff.UserID) // 將中獎人員加入redis中
// 		lock.Unlock()                         //釋放資源
// 	}(i)
// }
// wg.Wait() //等待計數器歸0

// // 更新獎品數量
// if err = models.DefaultPrizeModel().SetDbConn(conn).
// 	UpdateRemain(prizeID, prizeModel.PrizeRemain-int64(people)); err != nil {
// 	t.Error("錯誤: 獎品無庫存，請重新抽獎")
// }

// // 更新遊戲輪次
// if err = models.DefaultGameModel().SetDbConn(conn).SetRedisConn(redis).
// 	UpdateGameStatus(true, gameID, strconv.Itoa(int(gameModel.GameRound+1)),
// 		"", "", ""); err != nil {
// 	t.Error("錯誤: 無法更新遊戲狀態")
// }

// // 不允許重複中獎，將中獎人員加入redis中
// if gameModel.Allow == "close" {
// 	redis.SetAdd(params)
// 	redis.SetExpire(config.WINNING_STAFFS_REDIS+gameID,
// 		config.REDIS_EXPIRE)
// }

// endTime := time.Now()
// log.Println("統計中獎人員花費時間: ", endTime.Sub(startTime))
// log.Println("剩餘獎品是否為0: ", prizeModel.PrizeRemain-int64(people) == 0)
// log.Println("中獎人員是否為100人: ", len(numbers) == 100)
// log.Println("中獎人員: ", numbers)

// // 刪除測試資料
// db.Table(config.ACTIVITY_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// db.Table(config.ACTIVITY_GAME_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()

// // 清除redis測試資料
// redis.DelCache(config.ACTIVITY_REDIS + activityID)   // 活動資訊
// redis.DelCache(config.GAME_REDIS + gameID)           // 遊戲資訊
// redis.DelCache(config.WINNING_STAFFS_REDIS + gameID) // 獎品資訊
// }
