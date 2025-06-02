package models

import (
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"log"
	"testing"
)

// 將資料表所有資料加上game欄位資料
func Test_Game_FindAll(t *testing.T) {
	items, _ := db.Table("activity_staff_game").WithConn(conn).
		Where("game", "=", "").
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
			db.Table("activity_staff_game").WithConn(conn).
				Where("game_id", "=", gameID).Update(command.Value{"game": gameType})
		}
		// }()
	}

	// wg.Wait() //等待計數器歸0
}

// 測試查詢所有遊戲人員資訊(過濾.複選功能)
// func Test_Game_FindAll(t *testing.T) {
// 	var (
// 		activityID = "Bfon6SaV6ORhmuDQUioI"
// 		gameID     = "yVvGVWEHHQ6zCgkK84xM"
// 		game       = "ropepack,QA"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9,U2c979486326dc4157140fa396f0711d7"
// 		round      = "32,33,34"
// 		staffids1  = make([]int64, 0)
// 		staffids2  = make([]int64, 0)
// 		staffids3  = make([]int64, 0)
// 		staffids4  = make([]int64, 0)
// 		staffids5  = make([]int64, 0)
// 	)

// 	// 查詢活動下所有遊戲人員資料
// 	staffs1, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(activityID, "", "", "", "", 10, 0)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動下所有遊戲人員資料錯誤", err)
// 	}
// 	for _, staff := range staffs1 {
// 		staffids1 = append(staffids1, staff.ID)
// 	}
// 	fmt.Println("查詢活動下所有遊戲人員資料成功(應該有10筆): ", len(staffs1), staffids1)

// 	// 查詢活動下套紅包.快問快答遊戲人員資料
// 	staffs2, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(activityID, "", "", game, "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 查詢活動下套紅包.快問快答遊戲人員資料錯誤", err)
// 	}
// 	for _, staff := range staffs2 {
// 		staffids2 = append(staffids2, staff.ID)
// 	}
// 	fmt.Println("查詢活動下套紅包.快問快答遊戲人員資料成功(應該有5筆): ", len(staffs2), staffids2)

// 	// 查詢特定遊戲資料
// 	staffs3, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, "", game, "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 查詢特定遊戲資料錯誤", err)
// 	}
// 	for _, staff := range staffs3 {
// 		staffids3 = append(staffids3, staff.ID)
// 	}
// 	fmt.Println("查詢特定遊戲資料成功(應該有4筆): ", len(staffs3), staffids3)

// 	// 查詢特定遊戲特定人員資料
// 	staffs4, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, userID, game, "", 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 查詢特定遊戲特定人員資料錯誤", err)
// 	}
// 	for _, staff := range staffs4 {
// 		staffids4 = append(staffids4, staff.ID)
// 	}
// 	fmt.Println("查詢特定遊戲特定人員資料成功(應該有4筆): ", len(staffs4), staffids4)

// 	// 查詢特定遊戲特定人員特定輪次資料
// 	staffs5, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(activityID, gameID, userID, game, round, 0, 0)
// 	if err != nil {
// 		t.Error("錯誤: 查詢特定遊戲特定人員特定輪次資料錯誤", err)
// 	}
// 	for _, staff := range staffs5 {
// 		staffids5 = append(staffids5, staff.ID)
// 	}
// 	fmt.Println("查詢特定遊戲特定人員特定輪次資料成功(應該有3筆): ", len(staffs5), staffids5)
// }

// 測試新增遊戲人員資料(redis遊戲人數也新增)
// func Test_GameStaff_Add_Redis(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "4CZ8bjMt1bp53yDIKAOp"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = "1"
// 		status     = "success"
// 		// black      = "no"
// 	)

// 	game, err := DefaultGameModel().SetDbConn(conn).SetRedisConn(redis).
// 		Find(true, gameID)
// 	if err != nil || game.ID == 0 {
// 		t.Error("錯誤: 查詢遊戲資料發生問題(單個)", err)
// 	}
// 	fmt.Println("新增前人數: ", game.GameAttend)

// 	if err := DefaultGameStaffModel().SetDbConn(conn).SetRedisConn(redis).
// 		Add(NewGameStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			Round:      round,
// 			Status:     status,
// 			// Black:      black,
// 		}); err != nil {
// 		t.Error("錯誤: 新增遊戲人員資料發生問題", err)
// 	}

// 	game, err = DefaultGameModel().SetDbConn(conn).SetRedisConn(redis).
// 		Find(true, gameID)
// 	if err != nil || game.ID == 0 {
// 		t.Error("錯誤: 查詢遊戲資料發生問題(單個)", err)
// 	}
// 	fmt.Println("新增後人數: ", game.GameAttend)

// 	// 刪除redis中遊戲資訊
// 	redis.DelCache(config.GAME_REDIS + gameID)
// }

// IsBlack 利用redis判斷是否為黑名單(如果redis沒資料則從資料表取得並加入redis中)
// func Test_GameStaff_IsBlack_Redis(t *testing.T) {
// 	var (
// 		gameID = "WopzeMGctga0v7VJrGBb"
// 		userID = "Uc09f0955f4eb70e7d7bab794e15517ec"
// 		// round  = "1"
// 	)

// 	isBlack := DefaultGameStaffModel().SetDbConn(conn).SetRedisConn(redis).
// 		IsBlack(true, gameID, userID)
// 	fmt.Println("是否為黑名單: ", isBlack)

// 	redis.DelCache(config.BLACK_STAFFS_REDIS + gameID)
// }

// 查詢輪次所有遊戲人員資料
// func Test_GameStaff_LeftJoinLineUsers(t *testing.T) {
// 	var (
// 		gameID = "redpack"
// 		round  = "1"
// 	)

// 	staffs, err := DefaultGameStaffModel().SetDbConn(conn).
// 		FindAll(gameID, round)
// 	if err != nil {
// 		t.Error("錯誤: 查詢輪次所有遊戲人員資料發生問題", err)
// 	}
// 	fmt.Println("staffs: ", staffs)
// }

// 判斷用戶資料是否存在遊戲資料中
// func Test_GameStaff_IsUserExist(t *testing.T) {
// 	var (
// 		gameID = "redpack"
// 		userID = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round  = 1
// 	)

// 	isbool := DefaultGameStaffModel().SetDbConn(conn).
// 		IsUserExist(gameID, userID, int64(round))
// 	fmt.Println("是否存在: ", isbool)
// }

// 測試刪除測試遊戲人員資料
// func Test_GameStaff_Delete(t *testing.T) {
// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "redpack").Delete()
// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "ropepack").Delete()
// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "whack_mole").Delete()
// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "turntable").Delete()
// db.Table(config.ACTIVITY_STAFF_GAME_TABLE).WithConn(conn).Where("game_id", "=", "jiugongge").Delete()
// }

// 更新遊戲人員資料
// func Test_GameStaff_Update(t *testing.T) {
// 	var (
// 		id     = "2599"
// 		status = "success"
// 		// black  = "no"
// 	)

// 	if err := DefaultGameStaffModel().SetDbConn(conn).
// 		Update(false, EditGameStaffModel{
// 			ID:     id,
// 			Status: status,
// 			// Black:  black,
// 		}); err != nil {
// 		t.Error("錯誤: 更新遊戲人員資料發生問題", err)
// 	}
// }

// 測試新增遊戲人員資料
// func Test_GameStaff_Add(t *testing.T) {
// 	var (
// 		activityID = "FEVNqoH9Vv3iDlo7byeB"
// 		gameID     = "redpack"
// 		userID     = "U3140fd73cfd35dd992668ab3b6efdae9"
// 		round      = "1"
// 		status     = "success"
// 		// black      = "no"
// 	)

// 	if err := DefaultGameStaffModel().SetDbConn(conn).
// 		Add(NewGameStaffModel{
// 			UserID:     userID,
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			Round:      round,
// 			Status:     status,
// 			// Black:      black,
// 		}); err != nil {
// 		t.Error("錯誤: 新增遊戲人員資料
