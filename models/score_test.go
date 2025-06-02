package models

import (
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"log"
	"testing"
)

// 將資料表所有資料加上game欄位資料
func Test_Score(t *testing.T) {
	items, _ := db.Table("activity_score").WithConn(conn).
		Where("game", "=", "").
		Limit(20000).
		All()
	log.Println("資料數: ", len(items))

	// var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
	// wg.Add(len(items))    //計數器

	for _, item := range items {
		// go func() {
		// defer wg.Done()
		var (
			// activityID = utils.GetString(item["activity_id"], "")
			gameID = utils.GetString(item["game_id"], "")
			game   = utils.GetString(item["game"], "")
			// fields = make([]string, 0)
		)

		// game為空時執行
		if game == "" {
			// 取得遊戲種類
			gameType, _ := DefaultGameModel().
				SetConn(conn, redis,mongoConn).
				FindGameType(true, gameID)

				// 將該遊戲ID下的所有資料更新
			db.Table("activity_score").WithConn(conn).
				Where("game_id", "=", gameID).Update(command.Value{"game": gameType})
		}
		// }()
	}

	// wg.Wait() //等待計數器歸0
}

// 測試查詢所有遊戲紀錄資訊(過濾.複選功能)
// func Test_Score_FindAll(t *testing.T) {
// 	var (
// 		activityID = "Bfon6SaV6ORhmuDQUioI"
// 		gameID     = "Pb1eZzzoeWMTOM5GSRCC"
// 		game       = "monopoly,QA"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = "186,187,188,189,190,191,192,193,194"
// 		staffids1 = make([]int64, 0)
// 		// staffids2  = make([]int64, 0)
// 		// staffids3  = make([]int64, 0)
// 		// staffids4  = make([]int64, 0)
// 		// staffids5  = make([]int64, 0)
// 	)

// 	staffs1, err := DefaultScoreModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, userID, game, round, 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 資料錯誤", err)
// 	}
// 	for _, staff := range staffs1 {
// 		staffids1 = append(staffids1, staff.ID)
// 	}
// 	fmt.Println("資料成功: ", len(staffs1), staffids1)

// }

// 更新遊戲紀錄名次資料
// func Test_Update_Score_Rank_Data(t *testing.T) {
// 	var (
// 		rank   = 1
// 		gameID string
// 		round  int64
// 	)

// 	// 先取得activity_score所有資料
// 	items, err := db.Table("activity_score").WithConn(conn).
// 		Where("activity_score.score", ">", 0).
// 		OrderBy(
// 			"activity_score.game_id", "asc",
// 			"activity_score.round", "asc",
// 			"activity_score.score", "desc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_score所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_score所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			activityID = utils.GetString(item["activity_id"], "")
// 			userID     = utils.GetString(item["user_id"], "")
// 			newGameID  = utils.GetString(item["game_id"], "")
// 			newRound   = utils.GetInt64(item["round"], 0)
// 		)

// 		if newGameID != gameID || newRound != round {
// 			rank = 1
// 		}

// 		// 更新中獎人員名次資料
// 		err = db.Table("activity_score").WithConn(conn).
// 			Where("activity_score.activity_id", "=", activityID).
// 			Where("activity_score.game_id", "=", newGameID).
// 			Where("activity_score.user_id", "=", userID).
// 			Where("activity_score.round", "=", newRound).
// 			Update(command.Value{
// 				"rank": rank,
// 			})

// 		gameID = newGameID
// 		round = newRound
// 		rank++
// 	}
// 	if err != nil {
// 		t.Error("更新activity_score資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新activity_score資料資料完成")
// 	}
// }

// 測試當多個用戶分數相同時，回傳的排名順序顯示是否一樣(redis，zset格式)
// func Test_Score_Rank_Redis(t *testing.T) {
// 	var (
// 		gameID = "rank"
// 	)

// 	// 新增多個同分用戶分數資訊
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "a", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "b", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "c", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "z", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "x", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "v", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).Update(true, gameID, "n", 200); err != nil {
// 	// 	t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	// }

// 	// 查詢遊戲排名資訊
// 	// var users []string
// 	users, err := redis.ZSetRevRange(config.SCORES_REDIS+gameID, 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 從redis中取得分數由高至低的玩家資訊發生問題")
// 	}
// 	fmt.Println("排名: ", users)
// 	fmt.Println("-------------------------")
// 	users, err = redis.ZSetRevRange(config.SCORES_REDIS+gameID, 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 從redis中取得分數由高至低的玩家資訊發生問題")
// 	}
// 	fmt.Println("排名: ", users)
// 	fmt.Println("-------------------------")

// 	fmt.Println("用整數取分數: ", redis.ZSetIntScore(config.SCORES_REDIS+gameID, "a"))
// 	fmt.Println("用小數取分數: ", redis.ZSetFloatScore(config.SCORES_REDIS+gameID, "a"))
// 	// redis.DelCache(config.SCORES_REDIS + gameID)
// }

// 比較資料表與redis查詢100筆資料效率差異
// func Test_Compare_DB_Redis_Score(t *testing.T) {
// 	var (
// 		// activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID = "whack_mole"
// 		// userID1 = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		// userID2 = "Ua4fc0eda353d11d6a1af106e4fff53e6"
// 		round = 1
// 		limit = 100
// 	)

// 	// 新增測試遊戲資料
// 	// if err := DefaultGameModel().SetDbConn(conn).
// 	// 	Add("whack_mole", "whack_mole", NewGameModel{
// 	// 		ActivityID:   activityID,
// 	// 		Title:        "敲敲樂",
// 	// 		GameType:     "",
// 	// 		LimitTime:    "open",
// 	// 		Second:       "16",
// 	// 		MaxPeople:    "50",
// 	// 		MaxTimes:     "0",
// 	// 		Allow:        "open",
// 	// 		Percent:      "0",
// 	// 		FirstPrize:   "1",
// 	// 		SecondPrize:  "2",
// 	// 		ThirdPrize:   "3",
// 	// 		GeneralPrize: "4",
// 	// 		Topic:        "classic",
// 	// 		Skin:         "rat",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增敲敲樂遊戲資料發生問題", err)
// 	// }

// 	// 新增測試人員資料
// 	for i := 500; i < 1000; i++ {
// 		if _, err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).Insert(command.Value{
// 			"user_id": i, "name": i, "avatar": "https://dev.hilives.net/admin/uploads/yqwjQPj8.png",
// 			"email": i, "friend": "", "phone": "", "identify": i,
// 		}); err != nil {
// 			t.Error("錯誤: 新增人員發生問題")
// 		}
// 	}

// 	// 新增測試用戶分數紀錄
// 	for i := 1; i < 101; i++ {
// 		// redis
// 		if err := redis.ZSetAddInt(config.SCORES_REDIS+gameID, strconv.Itoa(i),
// 			int64(i)); err != nil {
// 			t.Error("錯誤: 新增redis分數資料發生問題")
// 		}

// 		// 資料表
// 		// if _, err := db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).Insert(command.Value{
// 		// 	"activity_id": activityID, "game_id": gameID, "user_id": i,
// 		// 	"round": round, "score": i,
// 		// }); err != nil {
// 		// 	t.Error("錯誤: 新增資料表分數資料發生問題")
// 		// }
// 	}

// 	// redis查詢遊戲排名資訊(平均10ms內)
// 	for i := 0; i < 5; i++ {
// 		startTime := time.Now()

// 		_, err := DefaultScoreModel().SetDbConn(conn).
// 			SetRedisConn(redis).Find(true, gameID, 0, int64(limit), "desc")
// 		if err != nil {
// 			t.Error("錯誤: 查詢遊戲排名發生問題(redis)")
// 		}
// 		// fmt.Println("排名: ", scores)

// 		endTime := time.Now()
// 		log.Println("redis查詢花費時間: ", endTime.Sub(startTime))
// 	}

// 	// 資料表查詢遊戲排名資訊(平均20ms)
// 	for i := 0; i < 5; i++ {
// 		startTime := time.Now()
// 		_, err := db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).
// 			Select("activity_score.id", "activity_score.activity_id",
// 				"activity_score.game_id", "activity_score.user_id",
// 				"activity_score.round", "activity_score.score",

// 				"line_users.name", "line_users.avatar").
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_score.user_id",
// 				FieldA1:   "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			Where("activity_score.game_id", "=", gameID).
// 			Where("activity_score.round", "=", round).
// 			Where("activity_score.score", ">", 0).
// 			Limit(int64(limit)).
// 			OrderBy("activity_score.score", "desc", "activity_score.id", "asc").All()
// 		if err != nil {
// 			t.Error("錯誤: 無法取得分數排名資訊")
// 		}

// 		endTime := time.Now()
// 		log.Println("資料表查詢花費時間: ", endTime.Sub(startTime))
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_GAME_TABLE).WithConn(conn).Where("game_id", "=", "whack_mole").Delete()
// 	redis.DelCache(config.SCORES_REDIS + gameID)
// }

// 測試用戶分數的新增、更新、查詢功能(redis，zset格式)
// func Test_Score_Redis(t *testing.T) {
// 	var (
// 		gameID  = "4CZ8bjMt1bp53yDIKAOp"
// 		userID1 = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		userID2 = "Ua4fc0eda353d11d6a1af106e4fff53e6"
// 	)

// 	// 新增用戶分數紀錄
// 	// 用戶1
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).UpdateScore(true, gameID, userID1); err != nil {
// 	// 	t.Error("錯誤: 新增用戶分數資料發生問題(redis)")
// 	// }
// 	// // 用戶2
// 	// if err := DefaultScoreModel().SetDbConn(conn).
// 	// 	SetRedisConn(redis).UpdateScore(true, gameID, userID2); err != nil {
// 	// 	t.Error("錯誤: 新增用戶分數資料發生問題(redis)")
// 	// }

// 	// 更新用戶分數資訊
// 	// 用戶1
// 	if err := DefaultScoreModel().SetDbConn(conn).
// 		SetRedisConn(redis).UpdateScore(true, gameID, userID1, 200); err != nil {
// 		t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	}
// 	// 用戶2
// 	if err := DefaultScoreModel().SetDbConn(conn).
// 		SetRedisConn(redis).UpdateScore(true, gameID, userID2, 10); err != nil {
// 		t.Error("錯誤: 更新用戶分數資料發生問題(redis)")
// 	}

// 	// 查詢遊戲排名資訊
// 	scores, err := DefaultScoreModel().SetDbConn(conn).
// 		SetRedisConn(redis).Find(true, gameID, 0, 5, "desc")
// 	if err != nil {
// 		t.Error("錯誤: 查詢遊戲排名發生問題(redis)")
// 	}
// 	fmt.Println("排名: ", scores)
// }

// 測試查詢打地鼠計分資料功能(包含新增、更新資料功能)
// func Test_Score_LeftJoinLineUsers(t *testing.T) {
// 	var (
// 		activityID = "test"
// 		gameID     = "test"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = 1
// 	)

// 	// 新增地鼠計分測試資料
// 	if err := DefaultScoreModel().SetDbConn(conn).
// 		Add(false, NewScoreModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     userID,
// 			Round:      strconv.Itoa(round),
// 		}); err != nil {
// 		t.Error("錯誤: 新增用戶打地鼠資料發生問題")
// 	} else {
// 		t.Log("新增用戶打地鼠資料成功")
// 	}

// 	// 更新計分資料
// 	if err := DefaultScoreModel().SetDbConn(conn).
// 		Update(false, EditScoreModel{
// 			GameID: gameID,
// 			UserID: userID,
// 			Round:  "1",
// 			Score:  "100",
// 		}); err != nil {
// 		t.Error("錯誤: 更新用戶計分資料發生問題")
// 	} else {
// 		t.Log("更新用戶計分資料成功")
// 	}

// 	scores, err := DefaultScoreModel().SetDbConn(conn).
// 		LeftJoinLineUsers(gameID, int64(round), 100)
// 	if err != nil {
// 		t.Error("錯誤: 查詢打地鼠計分紀錄發生問題")
// 	} else {
// 		t.Log("查詢資料完成(LeftJoinLineUsers)")

// 	}

// 	fmt.Println("打地鼠計分單元測試完成: ", scores)
// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).Where("activity_id", "=", activityID).Delete()
// }

// 測試新增打地鼠計分功能並測試結算效能(新增多筆資料測試)
// func Test_Add_Score_Data(t *testing.T) {
// 	var (
// 		activityID = "H3gM13lE60GAsPEJ6bSv"
// 		gameID     = "6NdWmDGMajCoBNtlHvCV"
// 		round      = "110"
// 		wg         sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	)

// 	for i := 1; i < 11; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			// 新增計分資料
// 			if _, err := db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).Insert(command.Value{
// 				"activity_id": activityID, "game_id": gameID, "user_id": i,
// 				"round": round, "score": i,
// 			}); err != nil {
// 				t.Error("錯誤: 新增人員發生問題")
// 			}
// 			if err := DefaultScoreModel().SetDbConn(conn).
// 				Add(false, NewScoreModel{
// 					ActivityID: activityID,
// 					GameID:     gameID,
// 					UserID:     strconv.Itoa(i),
// 					// Name:       strconv.Itoa(i),
// 					// Avatar:     "https://profile.line-scdn.net/0hsZ7yPLSYLHUKADgjjH1TIjZFIhh9Lio9cjExFngIdxUmOW50NGY2GipUIEEmOW8nMm5kFicIIUx3",
// 					Round: round,
// 				}); err != nil {
// 				t.Error("錯誤: 新增人員計分表資料發生問題", err)
// 			}

// 			// 更新分數資料
// 			// if err := DefaultScoreModel().SetDbConn(conn).
// 			// 	Update(EditScoreModel{
// 			// 		GameID: gameID,
// 			// 		UserID: strconv.Itoa(i),
// 			// 		Round:  round,
// 			// 		Score:  strconv.Itoa(i),
// 			// 	}); err != nil {
// 			// 	t.Error("錯誤: 更新人員計分表資料發生問題", err)
// 			// }

// 			// 新增人員資料
// 			// if _, err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).Insert(command.Value{
// 			// 	"user_id": i, "name": i, "avatar": "https://profile.line-scdn.net/0hCSeMHjdLHHxYNAiRsHFjK2RxEhEvGho0IFYBGS8yER5xBQkuZFJXEn9hFh8nVF0qZFtaGy41Ekxz",
// 			// 	"email": i, "friend": "", "phone": "", "identify": i,
// 			// }); err != nil {
// 			// 	t.Error("錯誤: 新增人員發生問題")
// 			// }
// 		}(i)
// 	}
// 	wg.Add(10) //計數器+1
// 	wg.Wait()  //等待計數器歸0

// 	// 刪除測試資料
// 	// roundInt, _ := strconv.Atoi(round)
// 	// db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).
// 	// 	Where("game_id", "=", gameID).Where("round", "=", roundInt-1).Delete()
// 	// db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).
// 	// 	Where("game_id", "=", gameID).Where("round", "=", roundInt-1).Delete()
// 	// db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }
