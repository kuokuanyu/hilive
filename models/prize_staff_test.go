package models

import (
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"log"
	"testing"
)

// 將資料表所有資料加上game欄位資料
func Test_Prize_Staff(t *testing.T) {
	items, _ := db.Table("activity_staff_prize").WithConn(conn).
		Where("game", "=", "").
		Limit(20000).
		OrderBy("id", "desc").
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
				SetConn(conn, redis, mongoConn).
				FindGameType(true, gameID)

				// 將該遊戲ID下的所有資料更新
			db.Table("activity_staff_prize").WithConn(conn).
				Where("game_id", "=", gameID).Update(command.Value{"game": gameType})
		}
		// }()
	}

	// wg.Wait() //等待計數器歸0
}

// 更新快問快答遊戲歷史中獎人員名次資料(更新至正式區時需要執行)
// func Test_Update_QA_Rank_Data(t *testing.T) {
// 	var (
// 		rank   = 1
// 		gameID string
// 		round  int64
// 	)

// 	// 先取得activity_staff_prize所有資料(快問快答)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "QA").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy(
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.score", "desc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
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
// 		err = db.Table("activity_staff_prize").WithConn(conn).
// 			Where("activity_staff_prize.activity_id", "=", activityID).
// 			Where("activity_staff_prize.game_id", "=", newGameID).
// 			Where("activity_staff_prize.user_id", "=", userID).
// 			Where("activity_staff_prize.round", "=", newRound).
// 			Update(command.Value{
// 				"rank": rank,
// 			})

// 		gameID = newGameID
// 		round = newRound
// 		rank++
// 	}
// 	if err != nil {
// 		t.Error("更新中獎人員名次資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新中獎人員名次資料資料完成")
// 	}
// }

// // 更新鑑定師遊戲歷史中獎人員名次資料(更新至正式區時需要執行)
// func Test_Update_Monopoly_Rank_Data(t *testing.T) {
// 	var (
// 		rank   = 1
// 		gameID string
// 		round  int64
// 	)

// 	// 先取得activity_staff_prize所有資料(鑑定師)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "monopoly").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy(
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.score", "desc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
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
// 		err = db.Table("activity_staff_prize").WithConn(conn).
// 			Where("activity_staff_prize.activity_id", "=", activityID).
// 			Where("activity_staff_prize.game_id", "=", newGameID).
// 			Where("activity_staff_prize.user_id", "=", userID).
// 			Where("activity_staff_prize.round", "=", newRound).
// 			Update(command.Value{
// 				"rank": rank,
// 			})

// 		gameID = newGameID
// 		round = newRound
// 		rank++
// 	}
// 	if err != nil {
// 		t.Error("更新中獎人員名次資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新中獎人員名次資料資料完成")
// 	}
// }

// // 更新敲敲樂遊戲歷史中獎人員名次資料(更新至正式區時需要執行)
// func Test_Update_Whack_Mole_Rank_Data(t *testing.T) {
// 	var (
// 		rank   = 1
// 		gameID string
// 		round  int64
// 	)

// 	// 先取得activity_staff_prize所有資料(敲敲樂)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "whack_mole").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy(
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.score", "desc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
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
// 		err = db.Table("activity_staff_prize").WithConn(conn).
// 			Where("activity_staff_prize.activity_id", "=", activityID).
// 			Where("activity_staff_prize.game_id", "=", newGameID).
// 			Where("activity_staff_prize.user_id", "=", userID).
// 			Where("activity_staff_prize.round", "=", newRound).
// 			Update(command.Value{
// 				"rank": rank,
// 			})

// 		gameID = newGameID
// 		round = newRound
// 		rank++
// 	}
// 	if err != nil {
// 		t.Error("更新中獎人員名次資料發生錯誤: ", err)
// 	} else {
// 		t.Log("更新中獎人員名次資料資料完成")
// 	}
// }

// // 更新快問快答遊戲歷史中獎人員分數資料(更新至正式區時需要執行)
// func Test_Update_QA_Score_Data(t *testing.T) {

// 	// 先取得activity_staff_prize所有資料(快問快答)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "QA").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy("activity_game.game", "asc",
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"activity_staff_prize.status", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.prize_id", "asc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
// 	}

// 	// 併發處理
// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	wg.Add(len(items))    //計數器
// 	for i := 0; i < len(items); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			var (
// 				activityID = utils.GetString(items[i]["activity_id"], "")
// 				gameID     = utils.GetString(items[i]["game_id"], "")
// 				userID     = utils.GetString(items[i]["user_id"], "")
// 				round      = utils.GetInt64(items[i]["round"], 0)
// 			)

// 			// 查詢用戶分數資料
// 			score, err := db.Table("activity_score").WithConn(conn).
// 				Where("activity_score.activity_id", "=", activityID).
// 				Where("activity_score.game_id", "=", gameID).
// 				Where("activity_score.user_id", "=", userID).
// 				Where("activity_score.round", "=", round).
// 				First()
// 			if err != nil {
// 				t.Error("查詢activity_score所有資料發生錯誤")
// 			} else {
// 				t.Log("查詢activity_score所有資料完成")
// 			}

// 			score1 := utils.GetInt64(score["score"], 0)
// 			score2 := utils.GetInt64(score["score_2"], 0)

// 			if score1 != 0 || score2 != 0 {
// 				// 更新中獎人員分數資料
// 				err = db.Table("activity_staff_prize").WithConn(conn).
// 					Where("activity_staff_prize.activity_id", "=", activityID).
// 					Where("activity_staff_prize.game_id", "=", gameID).
// 					Where("activity_staff_prize.user_id", "=", userID).
// 					Where("activity_staff_prize.round", "=", round).
// 					Update(command.Value{
// 						"score":   score1,
// 						"score_2": score2,
// 					})
// 			}
// 			if err != nil {
// 				t.Error("更新中獎人員分數資料發生錯誤: ", err)
// 			} else {
// 				t.Log("更新中獎人員分數資料資料完成")
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0
// }

// // 更新鑑定師遊戲歷史中獎人員分數資料(更新至正式區時需要執行)
// func Test_Update_Monopoly_Score_Data(t *testing.T) {

// 	// 先取得activity_staff_prize所有資料(鑑定師)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "monopoly").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy("activity_game.game", "asc",
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"activity_staff_prize.status", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.prize_id", "asc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
// 	}

// 	// 併發處理
// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	wg.Add(len(items))    //計數器
// 	for i := 0; i < len(items); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			var (
// 				activityID = utils.GetString(items[i]["activity_id"], "")
// 				gameID     = utils.GetString(items[i]["game_id"], "")
// 				userID     = utils.GetString(items[i]["user_id"], "")
// 				round      = utils.GetInt64(items[i]["round"], 0)
// 			)

// 			// 查詢用戶分數資料
// 			score, err := db.Table("activity_score").WithConn(conn).
// 				Where("activity_score.activity_id", "=", activityID).
// 				Where("activity_score.game_id", "=", gameID).
// 				Where("activity_score.user_id", "=", userID).
// 				Where("activity_score.round", "=", round).
// 				First()
// 			if err != nil {
// 				t.Error("查詢activity_score所有資料發生錯誤")
// 			} else {
// 				t.Log("查詢activity_score所有資料完成")
// 			}

// 			score1 := utils.GetInt64(score["score"], 0)
// 			score2 := utils.GetInt64(score["score_2"], 0)

// 			if score1 != 0 || score2 != 0 {
// 				// 更新中獎人員分數資料
// 				err = db.Table("activity_staff_prize").WithConn(conn).
// 					Where("activity_staff_prize.activity_id", "=", activityID).
// 					Where("activity_staff_prize.game_id", "=", gameID).
// 					Where("activity_staff_prize.user_id", "=", userID).
// 					Where("activity_staff_prize.round", "=", round).
// 					Update(command.Value{
// 						"score":   score1,
// 						"score_2": score2,
// 					})
// 			}
// 			if err != nil {
// 				t.Error("更新中獎人員分數資料發生錯誤: ", err)
// 			} else {
// 				t.Log("更新中獎人員分數資料資料完成")
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0
// }

// // 更新敲敲樂遊戲歷史中獎人員分數資料(更新至正式區時需要執行)
// func Test_Update_Whack_Mole_Score_Data(t *testing.T) {

// 	// 先取得activity_staff_prize所有資料(敲敲樂)
// 	items, err := db.Table("activity_staff_prize").WithConn(conn).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_game.game", "=", "whack_mole").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").
// 		OrderBy("activity_game.game", "asc",
// 			"activity_staff_prize.game_id", "asc",
// 			"activity_staff_prize.round", "asc",
// 			"activity_staff_prize.status", "asc",
// 			"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 			"activity_staff_prize.prize_id", "asc",
// 		).All()
// 	if err != nil {
// 		t.Error("查詢activity_staff_prize所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_staff_prize所有資料完成")
// 	}

// 	// 併發處理
// 	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	wg.Add(len(items))    //計數器
// 	for i := 0; i < len(items); i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			var (
// 				activityID = utils.GetString(items[i]["activity_id"], "")
// 				gameID     = utils.GetString(items[i]["game_id"], "")
// 				userID     = utils.GetString(items[i]["user_id"], "")
// 				round      = utils.GetInt64(items[i]["round"], 0)
// 			)

// 			// 查詢用戶分數資料
// 			score, err := db.Table("activity_score").WithConn(conn).
// 				Where("activity_score.activity_id", "=", activityID).
// 				Where("activity_score.game_id", "=", gameID).
// 				Where("activity_score.user_id", "=", userID).
// 				Where("activity_score.round", "=", round).
// 				First()
// 			if err != nil {
// 				t.Error("查詢activity_score所有資料發生錯誤")
// 			} else {
// 				t.Log("查詢activity_score所有資料完成")
// 			}

// 			score1 := utils.GetInt64(score["score"], 0)
// 			score2 := utils.GetInt64(score["score_2"], 0)

// 			if score1 != 0 || score2 != 0 {
// 				// 更新中獎人員分數資料
// 				err = db.Table("activity_staff_prize").WithConn(conn).
// 					Where("activity_staff_prize.activity_id", "=", activityID).
// 					Where("activity_staff_prize.game_id", "=", gameID).
// 					Where("activity_staff_prize.user_id", "=", userID).
// 					Where("activity_staff_prize.round", "=", round).
// 					Update(command.Value{
// 						"score":   score1,
// 						"score_2": score2,
// 					})
// 			}
// 			if err != nil {
// 				t.Error("更新中獎人員分數資料發生錯誤: ", err)
// 			} else {
// 				t.Log("更新中獎人員分數資料資料完成")
// 			}
// 		}(i)
// 	}
// 	wg.Wait() //等待計數器歸0
// }

// 測試查詢所有中獎人員資訊(過濾.複選功能)
// func Test_Prize_Staff_FindAll(t *testing.T) {
// 	var (
// 		activityID = "Bfon6SaV6ORhmuDQUioI"
// 		gameID     = "Pb1eZzzoeWMTOM5GSRCC"
// 		game       = "monopoly,QA"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = "187,188,189,190"
// 		status = "yes"
// 		staffids1 = make([]int64, 0)
// 		staffids2  = make([]int64, 0)
// 		staffids3  = make([]int64, 0)
// 		staffids4  = make([]int64, 0)
// 		staffids5  = make([]int64, 0)
// 	)

// 	// 活動下所有資料
// 	staffs1, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindAll(activityID, "", "", "", "", "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 活動下所有資料錯誤", err)
// 	}
// 	for _, staff := range staffs1 {
// 		staffids1 = append(staffids1, staff.ID)
// 	}
// 	fmt.Println("活動下所有資料成功(應該有56筆): ", len(staffs1), staffids1)

// 	// 活動下特定遊戲種類所有資料
// 	staffs2, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindAll(activityID, "", "", game, "", "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 活動下特定遊戲種類所有資料錯誤", err)
// 	}
// 	for _, staff := range staffs2 {
// 		staffids2 = append(staffids2, staff.ID)
// 	}
// 	fmt.Println("活動下特定遊戲種類所有資料成功(應該有19筆): ", len(staffs2), staffids2)

// 	// 活動下特定遊戲場次所有資料
// 	staffs3, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, "", game, "", "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 活動下特定遊戲場次所有資料錯誤", err)
// 	}
// 	for _, staff := range staffs3 {
// 		staffids3 = append(staffids3, staff.ID)
// 	}
// 	fmt.Println("活動下特定遊戲場次所有資料成功(應該有7筆): ", len(staffs3), staffids3)

// 	// 活動下特定用戶特定輪次所有資料
// 	staffs4, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, userID, game, round, "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 活動下特定遊戲場次所有資料錯誤", err)
// 	}
// 	for _, staff := range staffs4 {
// 		staffids4 = append(staffids4, staff.ID)
// 	}
// 	fmt.Println("活動下特定遊戲場次所有資料成功(應該有4筆): ", len(staffs4), staffids4)

// 	// 活動下特定用戶特定輪次特定status所有資料
// 	staffs5, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, userID, game, round, status, 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 活動下特定用戶特定輪次特定status所有資料錯誤", err)
// 	}
// 	for _, staff := range staffs5 {
// 		staffids5 = append(staffids5, staff.ID)
// 	}
// 	fmt.Println("活動下特定用戶特定輪次特定status所有資料成功(應該有1筆): ", len(staffs5), staffids5)
// }

// // 比較資料表與redis查詢100筆資料效率差異
// func Test_Compare_DB_Redis_PrizeStaff(t *testing.T) {
// 	var (
// 		// activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID = "whack_mole"
// 		// prizeID    = "redpack"
// 		// userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round = 1
// 	)

// 	// 中獎人員資訊(資料表查詢)，20-30ms左右
// 	for i := 0; i < 5; i++ {
// 		startTime := time.Now()

// 		_, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 			SetRedisConn(redis).FindWinningRecordsByGameID(false, 0, gameID, int64(round), 10)
// 		if err != nil {
// 			t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 		}

// 		endTime := time.Now()
// 		log.Println("資料表查詢花費時間: ", endTime.Sub(startTime))
// 	}

// 	log.Println("redis查詢")

// 	_, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindWinningRecordsByGameID(true, 0, gameID, int64(round), 10)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 	}

// 	// 中獎人員資訊(redis查詢)，1ms左右
// 	for i := 0; i < 5; i++ {
// 		startTime := time.Now()

// 		_, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 			SetRedisConn(redis).FindWinningRecordsByGameID(true, 0, gameID, int64(round), 100)
// 		if err != nil {
// 			t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 		}

// 		endTime := time.Now()
// 		log.Println("redis查詢花費時間: ", endTime.Sub(startTime))
// 	}

// 	// 清除redis資訊
// 	redis.DelCache(config.WINNING_STAFFS_REDIS + gameID)
// }

// // 利用redis查詢該遊戲場次的中獎紀錄(join line_users、activity_prize join)。redis沒有資料則由資料表查詢並加入redis中
// func Test_PrizeStaff_FindWinningRecordsByGameID_Redis(t *testing.T) {
// 	var (
// 		// activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID = "8Q48czHOxJiJt122AQj5"
// 		// prizeID    = "redpack"
// 		// userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round = 19
// 	)

// 	// 中獎人員資訊(特定輪次)
// 	staffs, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindWinningRecordsByGameID(true, 0, gameID, int64(round), 10)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 	}

// 	fmt.Println("中獎人員(特定輪次): ", staffs)

// 	// 清除redis資訊
// 	redis.DelCache(config.WINNING_STAFFS_REDIS + gameID)

// 	// 中獎人員資訊(不分輪次)
// 	staffs, err = DefaultPrizeStaffModel().SetDbConn(conn).
// 		SetRedisConn(redis).FindWinningRecordsByGameID(true, 0, gameID, 0, 3)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 	}

// 	fmt.Println("中獎人員(不分輪次): ", staffs)

// 	// 清除redis資訊
// 	redis.DelCache(config.WINNING_STAFFS_REDIS + gameID)
// }

// // 測試新增1000筆中獎人員資料
// func Test_Add_PrizeStaff_Data(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "Crmjs2jktiFbkjIbASb0"
// 		prizeID    = "uQfMjmkpShha0sdPevLk"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		// avatar     = "https://profile.line-scdn.net/0hsZ7yPLSYLHUKADgjjH1TIjZFIhh9Lio9cjExFngIdxUmOW50NGY2GipUIEEmOW8nMm5kFicIIUx3"
// 		round = "44"
// 		wg    sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	)

// 	for i := 1; i < 1001; i++ {
// 		go func(i int) {
// 			defer wg.Done()

// 			if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 				Add(NewPrizeStaffModel{
// 					UserID:     userID,
// 					ActivityID: activityID,
// 					GameID:     gameID,
// 					PrizeID:    prizeID,
// 					Round:      round,
// 					Status:     "no",
// 					// White:      "no",
// 				}); err != nil {
// 				t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 			}
// 		}(i)
// 	}
// 	wg.Add(1000) //計數器+1
// 	wg.Wait()    //等待計數器歸0

// 	// 刪除測試資料
// 	roundInt, _ := strconv.Atoi(round)
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).
// 		Where("game_id", "=", gameID).Where("round", "=", roundInt-1).Delete()

// 	// db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).
// 	// 	Where("game_id", "=", gameID).Where("round", "=", roundInt-1).Delete()
// 	// db.Table(config.ACTIVITY_APPLYSIGN_TABLE).WithConn(conn).Where("activity_id", "=", "test").Delete()
// }

// // 透過輪次查詢中獎紀錄並依照分數高低排序(join line_users、activity_prize、activity_whack_mole_score join)
// func Test_PrizeStaff_FindWinningRecordsOrderByScore(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "whack_mole"
// 		prizeID    = "whack_mole"
// 		// userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round = "1"
// 	)

// 	// 新增獎品資料
// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, "whack_mole", prizeID, NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "敲敲樂",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "666",
// 			PrizePrice:    "10654",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增獎品資料發生問題", err)
// 	}

// 	// 新增計分資料
// 	// if _, err := db.Table("activity_score").WithConn(conn).
// 	// 	Insert(command.Value{
// 	// 		"activity_id": activityID,
// 	// 		"game_id":     gameID,
// 	// 		"user_id":     userID,
// 	// 		"round":       round,
// 	// 		"score":       100,
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增計分資料發生問題")
// 	// }
// 	for i := 1; i <= 10; i++ {
// 		if _, err := db.Table("activity_score").WithConn(conn).
// 			Insert(command.Value{
// 				"activity_id": activityID,
// 				"game_id":     gameID,
// 				"user_id":     i,
// 				"round":       round,
// 				"score":       i,
// 			}); err != nil {
// 			t.Error("錯誤: 新增計分資料發生問題")
// 		}
// 	}

// 	// 新增中獎人員資料
// 	// if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	Add(NewPrizeStaffModel{
// 	// 		UserID:     userID,
// 	// 		ActivityID: activityID,
// 	// 		GameID:     gameID,
// 	// 		PrizeID:    prizeID,
// 	// 		Round:      round,
// 	// 		Status:     "no",
// 	// 		White:      "no",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	// }
// 	for i := 1; i <= 10; i++ {
// 		if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 			Add(NewPrizeStaffModel{
// 				UserID:     strconv.Itoa(i),
// 				ActivityID: activityID,
// 				GameID:     gameID,
// 				PrizeID:    prizeID,
// 				Round:      round,
// 				Status:     "no",
// 				// White:      "no",
// 			}); err != nil {
// 			t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 		}
// 	}

// 	// prizes, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	FindWinningRecordsOrderByScore(gameID, round)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 無法取得中獎紀錄，請重新查詢", err)
// 	// }

// 	// for i, prize := range prizes {
// 	// 	fmt.Println("排名", i+1, ": ", prize)
// 	// 	fmt.Println("分數: ", prize.Score)
// 	// }

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_SCORE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 透過用戶ID或輪次資訊查詢該遊戲場次的中獎紀錄(join line_users、activity_prize join)
// func Test_PrizeStaff_FindWinningRecords(t *testing.T) {
// 	// var (
// 	// 	activityID = "FEVNqoH9Vv3iDlo7byeB"
// 	// 	gameID     = "redpack"
// 	// 	prizeID    = "redpack"
// 	// 	userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	// 	round      = 1
// 	// )

// 	// 新增獎品資料
// 	// if err := DefaultPrizeModel().SetDbConn(conn).
// 	// 	Add(false,"redpack", prizeID, NewPrizeModel{
// 	// 		ActivityID:    activityID,
// 	// 		GameID:        gameID,
// 	// 		PrizeName:     "獎品",
// 	// 		PrizeType:     "general",
// 	// 		PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 	// 		PrizeAmount:   "80",
// 	// 		PrizePrice:    "10000",
// 	// 		PrizeMethod:   "site",
// 	// 		PrizePassword: "password",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增獎品資料發生問題", err)
// 	// }

// 	// 新增中獎人員資料，第一輪中獎
// 	// if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	Add(NewPrizeStaffModel{
// 	// 		UserID:     userID,
// 	// 		ActivityID: activityID,
// 	// 		GameID:     gameID,
// 	// 		PrizeID:    prizeID,
// 	// 		Round:      "1",
// 	// 		Status:     "no",
// 	// 		White:      "no",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	// }
// 	// 第二輪沒中
// 	// if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	Add(NewPrizeStaffModel{
// 	// 		UserID:     userID,
// 	// 		ActivityID: activityID,
// 	// 		GameID:     gameID,
// 	// 		PrizeID:    "",
// 	// 		Round:      "2",
// 	// 		Status:     "no",
// 	// 		White:      "no",
// 	// 	}); err != nil {
// 	// 	t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	// }

// 	// // 輪次
// 	// prizes, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	FindWinningRecords(false, gameID,  int64(round), 100)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(輪次查詢)", err)
// 	// }

// 	// for _, prize := range prizes {
// 	// 	fmt.Println("prize: ", prize)
// 	// }

// 	// fmt.Println("------------------------------")

// 	// 用戶ID
// 	// prizes, err = DefaultPrizeStaffModel().SetDbConn(conn).
// 	// 	FindWinningRecords(gameID, "activity_staff_prize.user_id", userID)
// 	// if err != nil {
// 	// 	t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(用戶ID查詢)", err)
// 	// }

// 	// for _, prize := range prizes {
// 	// 	fmt.Println("prize: ", prize)
// 	// }

// 	// 刪除測試資料
// 	// db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	// db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 透過game_id查詢該遊戲場次的中獎紀錄(join line_users、activity_prize join)，用在沒有輪次分別的遊戲(遊戲抽獎)
// func Test_PrizeStaff_FindWinningRecordsByGameID(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "ropepack"
// 		prizeID    = "ropepack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增獎品資料
// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, "ropepack", prizeID, NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "恭喜",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "80",
// 			PrizePrice:    "10000",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增獎品資料發生問題", err)
// 	}

// 	// 新增中獎人員資料，第一輪中獎
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "1",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}
// 	// 第二輪中獎
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "2",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}

// 	prizes, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindWinningRecordsByGameID(false, 0, gameID, 0, 1000)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(FindWinningRecordsByGameID)", err)
// 	}

// 	for _, prize := range prizes {
// 		fmt.Println("prize: ", prize)
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 透過輪次或用戶ID資訊查詢中獎及未中獎紀錄資料(join line_users、activity_prize join)
// func Test_PrizeStaff_LeftJoinPrizeAndLineUsers(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "ropepack"
// 		prizeID    = "ropepack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = 1
// 	)

// 	// 新增獎品資料
// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, "ropepack", prizeID, NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "恭喜",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "80",
// 			PrizePrice:    "10000",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增獎品資料發生問題", err)
// 	}

// 	// 新增中獎人員資料，第一輪中獎
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "1",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}
// 	// 第二輪沒中
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    "",
// 			Round:      "2",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}

// 	// 輪次
// 	prizes, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindWinningRecordsByGameID(false, 0, gameID, int64(round), 1000)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(輪次查詢)", err)
// 	}

// 	for _, prize := range prizes {
// 		fmt.Println("prize: ", prize)
// 	}

// 	fmt.Println("------------------------------")

// 	// 用戶ID
// 	prizes, err = DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindUserWinningRecords(userID, "activity_staff_prize.game_id", gameID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(用戶ID查詢)", err)
// 	}

// 	for _, prize := range prizes {
// 		fmt.Println("prize: ", prize)
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 透過用戶ID查詢中獎紀錄資料(join activity_game、activity_prize join)
// func Test_PrizeStaff_FindUserWinningRecords(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "ropepack"
// 		prizeID    = "ropepack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	// 新增獎品資料
// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, "ropepack", prizeID, NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "恭喜",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "80",
// 			PrizePrice:    "10000",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增獎品資料發生問題", err)
// 	}

// 	// 新增中獎人員資料，第一輪中獎
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "1",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}
// 	// 第二輪
// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "2",
// 			Status:     "no",
// 			// White:      "no",
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}

// 	prizes, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		FindUserWinningRecords(userID, "activity_staff_prize.activity_id", activityID)
// 	if err != nil {
// 		t.Error("錯誤: 無法取得中獎紀錄，請重新查詢(FindUserWinningRecords)", err)
// 	}

// 	for _, prize := range prizes {
// 		fmt.Println("prize: ", prize)
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 測試更新中獎人員資料(管理員、手機用戶)
// func Test_PrizeStaff_Update(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "ropepack"
// 		prizeID    = "ropepack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		status     = "yes"
// 		// black      = "no"
// 	)

// 	// 新增獎品資料
// 	if err := DefaultPrizeModel().SetDbConn(conn).
// 		Add(false, "ropepack", prizeID, NewPrizeModel{
// 			ActivityID:    activityID,
// 			GameID:        gameID,
// 			PrizeName:     "恭喜",
// 			PrizeType:     "first",
// 			PrizePicture:  "/admin/uploads/system/img-prize-pic.png",
// 			PrizeAmount:   "80",
// 			PrizePrice:    "10000",
// 			PrizeMethod:   "site",
// 			PrizePassword: "password",
// 		}); err != nil {
// 		t.Error("錯誤: 新增獎品資料發生問題", err)
// 	}

// 	// 新增中獎人員資料，第一輪中獎
// 	id1, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "1",
// 			Status:     "no",
// 			// White:      "no",
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增中獎人員1資料發生問題", err)
// 	}
// 	// 第二輪
// 	id2, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      "2",
// 			Status:     "no",
// 			// White:      "no",
// 		})
// 	if err != nil {
// 		t.Error("錯誤: 新增中獎人員2資料發生問題", err)
// 	}

// 	// 管理員(不用輸入密碼)
// 	if err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Update(EditPrizeStaffModel{
// 			ID:       strconv.Itoa(int(id1)),
// 			Role:     "admin",
// 			Password: "",
// 			Status:   status,
// 			// White:    black,
// 		}); err != nil {
// 		t.Error("錯誤: 更新中獎人員資料發生問題(管理員，不用輸入密碼驗證)", err)
// 	}

// 	// 手機用戶(需要輸入密碼驗證)
// 	if err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Update(EditPrizeStaffModel{
// 			ID:       strconv.Itoa(int(id2)),
// 			Role:     "guest",
// 			Password: "password",
// 			Status:   status,
// 			// White:    black,
// 		}); err != nil {
// 		t.Error("錯誤: 更新中獎人員資料發生問題(手機用戶，需要密碼驗證)", err)
// 	}

// 	// 刪除測試資料
// 	db.Table(config.ACTIVITY_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", gameID).Delete()
// }

// // 測試刪除測試中獎人員資料
// func Test_PrizeStaff_Delete(t *testing.T) {
// 	db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).Where("game_id", "=", "redpack").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "ropepack").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "whack_mole").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "turntable").Delete()
// 	// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "jiugongge").Delete()
// }

// // 測試新增中獎人員資料
// func Test_PrizeStaff_Add(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "redpack"
// 		prizeID    = "redpack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = "4"
// 		status     = "no"
// 		// black      = "no"
// 	)

// 	if _, err := DefaultPrizeStaffModel().SetDbConn(conn).
// 		Add(NewPrizeStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			PrizeID:    prizeID,
// 			Round:      round,
// 			Status:     status,
// 			// White:      black,
// 		}); err != nil {
// 		t.Error("錯誤: 新增中獎人員資料發生問題", err)
// 	}
// }

// 測試查詢所有遊戲先union後再與activity_staff_prize join
// func Test_PrizeStaffModel_FindByLeftJoin(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 	)

// 	items, err := db.Table(config.ACTIVITY_STAFF_PRIZE_TABLE).WithConn(conn).
// 		Select("game.title", "activity_game.game", "activity_staff_prize.picture",
// 			"activity_staff_prize.prize_name", "activity_staff_prize.win_time",
// 			"activity_staff_prize.status", "activity_staff_prize.method",
// 			"activity_staff_prize.password", "activity_staff_prize.id").
// 		LeftJoin(command.Join{
// 			FieldA:          "activity_staff_prize.game_id",
// 			FieldA1:          "game.game_id",
// 			Table:           "",
// 			Operation:       "=",
// 			UnionAliasTable: "game",
// 			UnionStatement: `SELECT game_id,title FROM activity_redpack
// 			union all SELECT game_id,title FROM activity_ropepack
// 			union all SELECT game_id,title FROM activity_lottery
// 			union all SELECT game_id,title FROM activity_whack_mole`}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:    "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		Where("activity_staff_prize.activity_id", "=", activityID).
// 		Where("activity_staff_prize.user_id", "=", userID).
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_staff_prize.method", "!=", "thanks").All()
// 	if err != nil {
// 		t.Error("錯誤: 無法查詢union、join資料 ", err)
// 	}

// 	for _, item := range items {
// 		// fmt.Println(item)
// 		id, _ := item["id"].(int64)
// 		title, _ := item["title"].(string)
// 		game, _ := item["game"].(string)
// 		password, _ := item["password"].(string)
// 		status, _ :
