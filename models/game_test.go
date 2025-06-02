package models

import (
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"testing"
)

// 將簽名牆自定義欄位填入預設值(更新至正式區時需要執行)
func Test_Update_Signname_Customize_Data(t *testing.T) {

	// 先取得activity_2所有資料
	items, err := db.Table("activity_2").WithConn(conn).
		Select("activity_id", "signname_topic").All()
	if err != nil {
		t.Error("查詢activity_2所有資料發生錯誤")
	} else {
		t.Log("查詢activity_2所有資料完成")
	}

	for _, item := range items {
		var (
			activityID = utils.GetString(item["activity_id"], "")
			// topics     = strings.Split(utils.GetString(item["signname_topic"], ""), "_")
			// topic      = ""
			route = "/admin/uploads/system/signname/%s/%s"
		)

		// if len(topics) == 2 {
		// 	topic = topics[1]
		// } else if len(topics) == 3 {
		// 	topic = topics[1] + "_" + topics[2]
		// }

		// 更新預設值
		err = db.Table("activity_2").WithConn(conn).
			Where("activity_id", "=", activityID).Update(command.Value{
			// "signname_classic_h_pic_01": fmt.Sprintf(route, "classic", "signname_classic_h_pic_01.png"),
			// "signname_classic_h_pic_02": fmt.Sprintf(route, "classic", "signname_classic_h_pic_02.png"),
			// "signname_classic_h_pic_03": fmt.Sprintf(route, "classic", "signname_classic_h_pic_03.png"),
			// "signname_classic_h_pic_04": fmt.Sprintf(route, "classic", "signname_classic_h_pic_04.png"),
			// "signname_classic_h_pic_05": fmt.Sprintf(route, "classic", "signname_classic_h_pic_05.png"),
			"signname_classic_h_pic_06": fmt.Sprintf(route, "classic", "signname_classic_h_pic_06.jpg"),
			"signname_classic_h_pic_07": fmt.Sprintf(route, "classic", "signname_classic_h_pic_07.png"),
			"signname_classic_g_pic_01": fmt.Sprintf(route, "classic", "signname_classic_g_pic_01.jpg"),
			// "signname_classic_c_pic_01": fmt.Sprintf(route, "classic", "signname_classic_c_pic_01.jpg"),
		})
		if err != nil {
			t.Error("更新自定義預設值發生錯誤", err)
		}

	}
}

// 更新遊戲資料表user_id欄位
// func Test_Update_Game_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "activity_id").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID     = utils.GetString(item["game_id"], "")
// 			activityID = utils.GetString(item["activity_id"], "")
// 		)

// 		// 查詢活動資料
// 		activity, err := db.Table("activity").WithConn(conn).
// 			Select("user_id", "activity_id").
// 			Where("activity_id", "=", activityID).First()
// 		if err != nil {
// 			t.Error("查詢activity_game所有資料發生錯誤")
// 		} else {
// 			t.Log("查詢activity_game所有資料完成")
// 		}

// 		// 更新user_id
// 		err = db.Table("activity_game").WithConn(conn).
// 			Where("game_id", "=", gameID).Update(command.Value{
// 			"user_id": utils.GetString(activity["user_id"], ""),
// 		})
// 		if err != nil {
// 			t.Error("更新user_id發生錯誤")
// 		}

// 	}
// }

// 將投票遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Vote_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/vote/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "vote" {
// 			// 更新預設值
// 			err = db.Table("activity_game_vote_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				// 投票遊戲自定義
// 				"vote_classic_h_pic_37": fmt.Sprintf(route, "classic", "vote_classic_h_pic_37.png"),
// 				"vote_classic_g_pic_07": fmt.Sprintf(route, "classic", "vote_classic_g_pic_07.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 修改line_users資料表自定義簽到人員資料
// func Test_LINE_Users_Update(t *testing.T) {
// 	// 先取得line_users所有資料
// 	items, err := db.Table("line_users").WithConn(conn).
// 		Select("user_id", "activity_id", "identify").All()
// 	if err != nil {
// 		t.Error("查詢line_users所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢line_users所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			activityID = utils.GetString(item["activity_id"], "")
// 			userID     = utils.GetString(item["user_id"], "")
// 			identify = utils.GetString(item["identify"], "")
// 		)

// 		// 判斷是否為自定義用戶
// 		if strings.Contains(userID, "_") && activityID != ""&& userID != identify {
// 			// 更新identify
// 			err = db.Table("line_users").WithConn(conn).
// 				Where("user_id", "=", userID).
// 				Update(command.Value{
// 					"identify": userID,
// 				})
// 			if err != nil {
// 				t.Error("更新預設值發生錯誤", err)
// 			}
// 		}
// 	}
// }

// 將一般簽到牆自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_General_Customize_Data(t *testing.T) {

// 	// 先取得activity_2所有資料
// 	items, err := db.Table("activity_2").WithConn(conn).
// 		Select("activity_id", "general_topic").All()
// 	if err != nil {
// 		t.Error("查詢activity_2所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_2所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			activityID = utils.GetString(item["activity_id"], "")
// 			// topics     = strings.Split(utils.GetString(item["general_topic"], ""), "_")
// 			// topic      = ""
// 			route = "/admin/uploads/system/general/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		// 更新預設值
// 		err = db.Table("activity_2").WithConn(conn).
// 			Where("activity_id", "=", activityID).Update(command.Value{
// 			"general_classic_h_pic_01": fmt.Sprintf(route, "classic", "general_classic_h_pic_01.png"),
// 			"general_classic_h_pic_02": fmt.Sprintf(route, "classic", "general_classic_h_pic_02.png"),
// 			"general_classic_h_pic_03": fmt.Sprintf(route, "classic", "general_classic_h_pic_03.jpg"),
// 			"general_classic_h_ani_01": fmt.Sprintf(route, "classic", "general_classic_h_ani_01.png"),
// 		})
// 		if err != nil {
// 			t.Error("更新自定義預設值發生錯誤", err)
// 		}

// 	}
// }

// 將賓果遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Bingo_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/bingo/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "bingo" {
// 			// 更新預設值
// 			err = db.Table("activity_game_bingo_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				// 賓果遊戲自定義
// 				"bingo_classic_h_pic_16": fmt.Sprintf(route, "classic", "bingo_classic_h_pic_16.png"),
// 				"bingo_classic_g_pic_08": fmt.Sprintf(route, "classic", "bingo_classic_g_pic_08.png"),
// 				// "bingo_classic_c_pic_05": fmt.Sprintf(route, "classic", "bingo_classic_c_pic_05.png"),
// 				"bingo_classic_g_ani_01": fmt.Sprintf(route, "classic", "bingo_classic_g_ani_01.png"),

// 				"bingo_newyear_dragon_g_pic_08": fmt.Sprintf(route, "newyear_dragon", "bingo_newyear_dragon_g_pic_08.png"),
// 				"bingo_newyear_dragon_g_ani_01": fmt.Sprintf(route, "newyear_dragon", "bingo_newyear_dragon_g_ani_01.png"),

// 				"bingo_cherry_h_pic_19": fmt.Sprintf(route, "cherry", "bingo_cherry_h_pic_19.png"),
// 				"bingo_cherry_g_pic_06": fmt.Sprintf(route, "cherry", "bingo_cherry_g_pic_06.png"),

// 				"bingo_cherry_g_ani_02": fmt.Sprintf(route, "cherry", "bingo_cherry_g_ani_02.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將所有場次資料加入排序預設值(更新至正式區時需要執行)
// func Test_Add_Game_Order_Data(t *testing.T) {
// 	var (
// 		activityID string
// 		count      int64
// 	)

// 	// 先取得activity_game所有資料(activity_id排序)
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("activity_id", "game_id", "game").
// 		OrderBy("activity_id", "asc").
// 		All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			id       = utils.GetString(item["activity_id"], "")
// 			gameID   = utils.GetString(item["game_id"], "")
// 			game     = utils.GetString(item["game"], "")
// 			overview string
// 		)

// 		if id != activityID {
// 			// 活動場次不一樣，count回復預設值1
// 			count = 1

// 			activityID = id
// 		} else if id == activityID {
// 			// 活動場次相同，count++
// 			count++
// 		}

// 		// 更新排序資料
// 		db.Table("activity_game").WithConn(conn).
// 			Where("game_id", "=", gameID).Update(command.Value{
// 			// 排序
// 			"game_order": count,
// 		})

// 		if game == "redpack" {
// 			overview = "overview_redpack"
// 		} else if game == "ropepack" {
// 			overview = "overview_ropepack"
// 		} else if game == "whack_mole" {
// 			overview = "overview_whack_mole"
// 		} else if game == "draw_numbers" {
// 			overview = "overview_draw_numbers"
// 		} else if game == "monopoly" {
// 			overview = "overview_monopoly"
// 		} else if game == "lottery" {
// 			overview = "overview_lottery"
// 		} else if game == "QA" {
// 			overview = "overview_qa"
// 		} else if game == "tugofwar" {
// 			overview = "overview_tugofwar"
// 		} else if game == "bingo" {
// 			overview = "overview_bingo"
// 		} else if game == "3DGachaMachine" {
// 			overview = "overview_3d_gacha_machine"
// 		}

// 		// 更新活動總覽資料
// 		db.Table("activity").WithConn(conn).
// 			Where("activity_id", "=", id).Update(command.Value{
// 			// 排序
// 			overview: "open",
// 		})

// 	}
// }

// 將扭蛋機遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_3d_Gacha_Machine_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/3DGachaMachine/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "3DGachaMachine" {
// 			// 更新預設值
// 			err = db.Table("activity_game_3d_gacha_machine_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{

// 				// 扭蛋機自定義
// 				"3d_gacha_machine_classic_g_pic_02": fmt.Sprintf(route, "classic", "3d_gacha_machine_classic_g_pic_02.png"),
// 				"3d_gacha_machine_classic_g_pic_03": fmt.Sprintf(route, "classic", "3d_gacha_machine_classic_g_pic_03.png"),
// 				"3d_gacha_machine_classic_g_pic_04": fmt.Sprintf(route, "classic", "3d_gacha_machine_classic_g_pic_04.png"),
// 				"3d_gacha_machine_classic_g_pic_05": fmt.Sprintf(route, "classic", "3d_gacha_machine_classic_g_pic_05.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將立體簽到牆自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Threed_Customize_Data(t *testing.T) {

// 	// 先取得activity_2所有資料
// 	items, err := db.Table("activity_2").WithConn(conn).
// 		Select("activity_id", "threed_topic").All()
// 	if err != nil {
// 		t.Error("查詢activity_2所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_2所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			activityID = utils.GetString(item["activity_id"], "")
// 			topics     = strings.Split(utils.GetString(item["threed_topic"], ""), "_")
// 			topic      = ""
// 			route      = "/admin/uploads/system/threed/%s/%s"
// 		)

// 		if len(topics) == 2 {
// 			topic = topics[1]
// 		} else if len(topics) == 3 {
// 			topic = topics[1] + "_" + topics[2]
// 		}

// 		// 更新預設值
// 		err = db.Table("activity_2").WithConn(conn).
// 			Where("activity_id", "=", activityID).Update(command.Value{

// 			"threed_bgm": fmt.Sprintf(route, topic, "bgm/bgm.mp3"),
// 		})
// 		if err != nil {
// 			t.Error("更新自定義預設值發生錯誤", err)
// 		}

// 	}
// }

// 將一般簽到牆自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_General_Customize_Data(t *testing.T) {

// 	// 先取得activity_2所有資料
// 	items, err := db.Table("activity_2").WithConn(conn).
// 		Select("activity_id", "general_topic").All()
// 	if err != nil {
// 		t.Error("查詢activity_2所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_2所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			activityID = utils.GetString(item["activity_id"], "")
// 			topics     = strings.Split(utils.GetString(item["general_topic"], ""), "_")
// 			topic      = ""
// 			route      = "/admin/uploads/system/general/%s/%s"
// 		)

// 		if len(topics) == 2 {
// 			topic = topics[1]
// 		} else if len(topics) == 3 {
// 			topic = topics[1] + "_" + topics[2]
// 		}

// 		// 更新預設值
// 		err = db.Table("activity_2").WithConn(conn).
// 			Where("activity_id", "=", activityID).Update(command.Value{

// 			"general_bgm": fmt.Sprintf(route, topic, "bgm/bgm.mp3"),
// 		})
// 		if err != nil {
// 			t.Error("更新自定義預設值發生錯誤", err)
// 		}

// 	}
// }

// 將快問快答遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_QA_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/qa/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "QA" {
// 			// 更新預設值
// 			err = db.Table("activity_game_qa_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				"qa_classic_h_pic_01": fmt.Sprintf(route, "classic", "qa_classic_h_pic_01.png"),
// 				"qa_classic_h_pic_02": fmt.Sprintf(route, "classic", "qa_classic_h_pic_02.png"),
// 				"qa_classic_h_pic_03": fmt.Sprintf(route, "classic", "qa_classic_h_pic_03.jpg"),
// 				"qa_classic_h_pic_04": fmt.Sprintf(route, "classic", "qa_classic_h_pic_04.jpg"),
// 				"qa_classic_h_pic_05": fmt.Sprintf(route, "classic", "qa_classic_h_pic_05.png"),
// 				"qa_classic_h_pic_06": fmt.Sprintf(route, "classic", "qa_classic_h_pic_06.png"),
// 				"qa_classic_h_pic_07": fmt.Sprintf(route, "classic", "qa_classic_h_pic_07.png"),
// 				"qa_classic_h_pic_08": fmt.Sprintf(route, "classic", "qa_classic_h_pic_08.png"),
// 				"qa_classic_h_pic_09": fmt.Sprintf(route, "classic", "qa_classic_h_pic_09.png"),
// 				"qa_classic_h_pic_10": fmt.Sprintf(route, "classic", "qa_classic_h_pic_10.png"),
// 				"qa_classic_h_pic_11": fmt.Sprintf(route, "classic", "qa_classic_h_pic_11.png"),
// 				"qa_classic_h_pic_12": fmt.Sprintf(route, "classic", "qa_classic_h_pic_12.png"),
// 				"qa_classic_h_pic_13": fmt.Sprintf(route, "classic", "qa_classic_h_pic_13.png"),
// 				"qa_classic_h_pic_14": fmt.Sprintf(route, "classic", "qa_classic_h_pic_14.png"),
// 				"qa_classic_h_pic_15": fmt.Sprintf(route, "classic", "qa_classic_h_pic_15.png"),
// 				"qa_classic_h_pic_16": fmt.Sprintf(route, "classic", "qa_classic_h_pic_16.png"),
// 				"qa_classic_h_pic_17": fmt.Sprintf(route, "classic", "qa_classic_h_pic_17.png"),
// 				"qa_classic_h_pic_18": fmt.Sprintf(route, "classic", "qa_classic_h_pic_18.png"),
// 				"qa_classic_h_pic_19": fmt.Sprintf(route, "classic", "qa_classic_h_pic_19.png"),
// 				"qa_classic_h_pic_20": fmt.Sprintf(route, "classic", "qa_classic_h_pic_20.png"),
// 				"qa_classic_h_pic_21": fmt.Sprintf(route, "classic", "qa_classic_h_pic_21.png"),
// 				"qa_classic_h_pic_22": fmt.Sprintf(route, "classic", "qa_classic_h_pic_22.png"),
// 				"qa_classic_g_pic_01": fmt.Sprintf(route, "classic", "qa_classic_g_pic_01.jpg"),
// 				"qa_classic_g_pic_02": fmt.Sprintf(route, "classic", "qa_classic_g_pic_02.jpg"),
// 				"qa_classic_g_pic_03": fmt.Sprintf(route, "classic", "qa_classic_g_pic_03.png"),
// 				"qa_classic_g_pic_04": fmt.Sprintf(route, "classic", "qa_classic_g_pic_04.png"),
// 				"qa_classic_g_pic_05": fmt.Sprintf(route, "classic", "qa_classic_g_pic_05.png"),
// 				"qa_classic_c_pic_01": fmt.Sprintf(route, "classic", "qa_classic_c_pic_01.png"),
// 				"qa_classic_h_ani_01": fmt.Sprintf(route, "classic", "qa_classic_h_ani_01.png"),
// 				"qa_classic_h_ani_02": fmt.Sprintf(route, "classic", "qa_classic_h_ani_02.png"),
// 				"qa_classic_g_ani_01": fmt.Sprintf(route, "classic", "qa_classic_g_ani_01.png"),
// 				"qa_classic_g_ani_02": fmt.Sprintf(route, "classic", "qa_classic_g_ani_02.png"),

// 				"qa_electric_h_pic_01": fmt.Sprintf(route, "electric", "qa_electric_h_pic_01.png"),
// 				"qa_electric_h_pic_02": fmt.Sprintf(route, "electric", "qa_electric_h_pic_02.png"),
// 				"qa_electric_h_pic_03": fmt.Sprintf(route, "electric", "qa_electric_h_pic_03.png"),
// 				"qa_electric_h_pic_04": fmt.Sprintf(route, "electric", "qa_electric_h_pic_04.jpg"),
// 				"qa_electric_h_pic_05": fmt.Sprintf(route, "electric", "qa_electric_h_pic_05.png"),
// 				"qa_electric_h_pic_06": fmt.Sprintf(route, "electric", "qa_electric_h_pic_06.png"),
// 				"qa_electric_h_pic_07": fmt.Sprintf(route, "electric", "qa_electric_h_pic_07.png"),
// 				"qa_electric_h_pic_08": fmt.Sprintf(route, "electric", "qa_electric_h_pic_08.png"),
// 				"qa_electric_h_pic_09": fmt.Sprintf(route, "electric", "qa_electric_h_pic_09.png"),
// 				"qa_electric_h_pic_10": fmt.Sprintf(route, "electric", "qa_electric_h_pic_10.png"),
// 				"qa_electric_h_pic_11": fmt.Sprintf(route, "electric", "qa_electric_h_pic_11.png"),
// 				"qa_electric_h_pic_12": fmt.Sprintf(route, "electric", "qa_electric_h_pic_12.png"),
// 				"qa_electric_h_pic_13": fmt.Sprintf(route, "electric", "qa_electric_h_pic_13.png"),
// 				"qa_electric_h_pic_14": fmt.Sprintf(route, "electric", "qa_electric_h_pic_14.png"),
// 				"qa_electric_h_pic_15": fmt.Sprintf(route, "electric", "qa_electric_h_pic_15.jpg"),
// 				"qa_electric_h_pic_16": fmt.Sprintf(route, "electric", "qa_electric_h_pic_16.png"),
// 				"qa_electric_h_pic_17": fmt.Sprintf(route, "electric", "qa_electric_h_pic_17.png"),
// 				"qa_electric_h_pic_18": fmt.Sprintf(route, "electric", "qa_electric_h_pic_18.png"),
// 				"qa_electric_h_pic_19": fmt.Sprintf(route, "electric", "qa_electric_h_pic_19.png"),
// 				"qa_electric_h_pic_20": fmt.Sprintf(route, "electric", "qa_electric_h_pic_20.jpg"),
// 				"qa_electric_h_pic_21": fmt.Sprintf(route, "electric", "qa_electric_h_pic_21.png"),
// 				"qa_electric_h_pic_22": fmt.Sprintf(route, "electric", "qa_electric_h_pic_22.png"),
// 				"qa_electric_h_pic_23": fmt.Sprintf(route, "electric", "qa_electric_h_pic_23.png"),
// 				"qa_electric_h_pic_24": fmt.Sprintf(route, "electric", "qa_electric_h_pic_24.png"),
// 				"qa_electric_h_pic_25": fmt.Sprintf(route, "electric", "qa_electric_h_pic_25.png"),
// 				"qa_electric_h_pic_26": fmt.Sprintf(route, "electric", "qa_electric_h_pic_26.png"),
// 				"qa_electric_g_pic_01": fmt.Sprintf(route, "electric", "qa_electric_g_pic_01.png"),
// 				"qa_electric_g_pic_02": fmt.Sprintf(route, "electric", "qa_electric_g_pic_02.png"),
// 				"qa_electric_g_pic_03": fmt.Sprintf(route, "electric", "qa_electric_g_pic_03.png"),
// 				"qa_electric_g_pic_04": fmt.Sprintf(route, "electric", "qa_electric_g_pic_04.png"),
// 				"qa_electric_g_pic_05": fmt.Sprintf(route, "electric", "qa_electric_g_pic_05.jpg"),
// 				"qa_electric_g_pic_06": fmt.Sprintf(route, "electric", "qa_electric_g_pic_06.png"),
// 				"qa_electric_g_pic_07": fmt.Sprintf(route, "electric", "qa_electric_g_pic_07.jpg"),
// 				"qa_electric_g_pic_08": fmt.Sprintf(route, "electric", "qa_electric_g_pic_08.png"),
// 				"qa_electric_g_pic_09": fmt.Sprintf(route, "electric", "qa_electric_g_pic_09.png"),
// 				"qa_electric_c_pic_01": fmt.Sprintf(route, "electric", "qa_electric_c_pic_01.png"),
// 				"qa_electric_h_ani_01": fmt.Sprintf(route, "electric", "qa_electric_h_ani_01.png"),
// 				"qa_electric_h_ani_02": fmt.Sprintf(route, "electric", "qa_electric_h_ani_02.png"),
// 				"qa_electric_h_ani_03": fmt.Sprintf(route, "electric", "qa_electric_h_ani_03.png"),
// 				"qa_electric_h_ani_04": fmt.Sprintf(route, "electric", "qa_electric_h_ani_04.png"),
// 				"qa_electric_h_ani_05": fmt.Sprintf(route, "electric", "qa_electric_h_ani_05.png"),
// 				"qa_electric_g_ani_01": fmt.Sprintf(route, "electric", "qa_electric_g_ani_01.png"),
// 				"qa_electric_g_ani_02": fmt.Sprintf(route, "electric", "qa_electric_g_ani_02.png"),
// 				"qa_electric_c_ani_01": fmt.Sprintf(route, "electric", "qa_electric_c_ani_01.png"),

// 				"qa_moonfestival_h_pic_01": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_01.png"),
// 				"qa_moonfestival_h_pic_02": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_02.png"),
// 				"qa_moonfestival_h_pic_03": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_03.png"),
// 				"qa_moonfestival_h_pic_04": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_04.png"),
// 				"qa_moonfestival_h_pic_05": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_05.jpg"),
// 				"qa_moonfestival_h_pic_06": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_06.png"),
// 				"qa_moonfestival_h_pic_07": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_07.png"),
// 				"qa_moonfestival_h_pic_08": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_08.png"),
// 				"qa_moonfestival_h_pic_09": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_09.png"),
// 				"qa_moonfestival_h_pic_10": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_10.png"),
// 				"qa_moonfestival_h_pic_11": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_11.png"),
// 				"qa_moonfestival_h_pic_12": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_12.png"),
// 				"qa_moonfestival_h_pic_13": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_13.png"),
// 				"qa_moonfestival_h_pic_14": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_14.png"),
// 				"qa_moonfestival_h_pic_15": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_15.png"),
// 				"qa_moonfestival_h_pic_16": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_16.png"),
// 				"qa_moonfestival_h_pic_17": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_17.png"),
// 				"qa_moonfestival_h_pic_18": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_18.png"),
// 				"qa_moonfestival_h_pic_19": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_19.png"),
// 				"qa_moonfestival_h_pic_20": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_20.png"),
// 				"qa_moonfestival_h_pic_21": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_21.png"),
// 				"qa_moonfestival_h_pic_22": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_22.png"),
// 				"qa_moonfestival_h_pic_23": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_23.png"),
// 				"qa_moonfestival_h_pic_24": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_pic_24.png"),
// 				"qa_moonfestival_g_pic_01": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_pic_01.png"),
// 				"qa_moonfestival_g_pic_02": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_pic_02.png"),
// 				"qa_moonfestival_g_pic_03": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_pic_03.jpg"),
// 				"qa_moonfestival_g_pic_04": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_pic_04.png"),
// 				"qa_moonfestival_g_pic_05": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_pic_05.png"),
// 				"qa_moonfestival_c_pic_01": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_c_pic_01.png"),
// 				"qa_moonfestival_c_pic_02": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_c_pic_02.png"),
// 				"qa_moonfestival_c_pic_03": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_c_pic_03.png"),
// 				"qa_moonfestival_h_ani_01": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_ani_01.png"),
// 				"qa_moonfestival_h_ani_02": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_h_ani_02.png"),
// 				"qa_moonfestival_g_ani_01": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_ani_01.png"),
// 				"qa_moonfestival_g_ani_02": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_ani_02.png"),
// 				"qa_moonfestival_g_ani_03": fmt.Sprintf(route, "moonfestival", "qa_moonfestival_g_ani_03.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將快問快答遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_QA_2_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/qa/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "QA" {
// 			// 更新預設值
// 			err = db.Table("activity_game_qa_picture_2").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				"qa_newyear_dragon_h_pic_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_01.png"),
// 				"qa_newyear_dragon_h_pic_02": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_02.png"),
// 				"qa_newyear_dragon_h_pic_03": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_03.png"),
// 				"qa_newyear_dragon_h_pic_04": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_04.png"),
// 				"qa_newyear_dragon_h_pic_05": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_05.png"),
// 				"qa_newyear_dragon_h_pic_06": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_06.jpg"),
// 				"qa_newyear_dragon_h_pic_07": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_07.png"),
// 				"qa_newyear_dragon_h_pic_08": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_08.png"),
// 				"qa_newyear_dragon_h_pic_09": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_09.png"),
// 				"qa_newyear_dragon_h_pic_10": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_10.png"),
// 				"qa_newyear_dragon_h_pic_11": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_11.png"),
// 				"qa_newyear_dragon_h_pic_12": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_12.png"),
// 				"qa_newyear_dragon_h_pic_13": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_13.png"),
// 				"qa_newyear_dragon_h_pic_14": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_14.png"),
// 				"qa_newyear_dragon_h_pic_15": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_15.png"),
// 				"qa_newyear_dragon_h_pic_16": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_16.png"),
// 				"qa_newyear_dragon_h_pic_17": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_17.png"),
// 				"qa_newyear_dragon_h_pic_18": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_18.png"),
// 				"qa_newyear_dragon_h_pic_19": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_19.png"),
// 				"qa_newyear_dragon_h_pic_20": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_20.png"),
// 				"qa_newyear_dragon_h_pic_21": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_21.png"),
// 				"qa_newyear_dragon_h_pic_22": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_22.png"),
// 				"qa_newyear_dragon_h_pic_23": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_23.png"),
// 				"qa_newyear_dragon_h_pic_24": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_pic_24.png"),
// 				"qa_newyear_dragon_g_pic_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_01.png"),
// 				"qa_newyear_dragon_g_pic_02": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_02.png"),
// 				"qa_newyear_dragon_g_pic_03": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_03.png"),
// 				"qa_newyear_dragon_g_pic_04": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_04.png"),
// 				"qa_newyear_dragon_g_pic_05": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_05.jpg"),
// 				"qa_newyear_dragon_g_pic_06": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_pic_06.png"),
// 				"qa_newyear_dragon_c_pic_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_c_pic_01.png"),
// 				"qa_newyear_dragon_h_ani_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_ani_01.png"),
// 				"qa_newyear_dragon_h_ani_02": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_h_ani_02.png"),
// 				"qa_newyear_dragon_g_ani_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_ani_01.png"),
// 				"qa_newyear_dragon_g_ani_02": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_ani_02.png"),
// 				"qa_newyear_dragon_g_ani_03": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_g_ani_03.png"),
// 				"qa_newyear_dragon_c_ani_01": fmt.Sprintf(route, "newyear_dragon", "qa_newyear_dragon_c_ani_01.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將搖號抽獎自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Draw_Numbers_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/draw_numbers/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "draw_numbers" {
// 			// 更新預設值
// 			err = db.Table("activity_game_draw_numbers_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				// 遊戲抽獎自定義
// 				"draw_numbers_3D_space_h_pic_01": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_01.png"),
// 				"draw_numbers_3D_space_h_pic_02": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_02.png"),
// 				"draw_numbers_3D_space_h_pic_03": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_03.png"),
// 				"draw_numbers_3D_space_h_pic_04": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_04.png"),
// 				"draw_numbers_3D_space_h_pic_05": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_05.png"),
// 				"draw_numbers_3D_space_h_pic_06": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_06.png"),
// 				"draw_numbers_3D_space_h_pic_07": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_07.png"),
// 				"draw_numbers_3D_space_h_pic_08": fmt.Sprintf(route, "3D_space", "draw_numbers_3D_space_h_pic_08.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將拔河遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Tugofwar_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/tugofwar/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "tugofwar" {
// 			// 更新預設值
// 			err = db.Table("activity_game_tugofwar_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{

// 				// 拔河遊戲自定義
// 				"tugofwar_classic_h_pic_01": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_01.png"),
// 				"tugofwar_classic_h_pic_02": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_02.png"),
// 				"tugofwar_classic_h_pic_03": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_03.png"),
// 				"tugofwar_classic_h_pic_04": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_04.png"),
// 				"tugofwar_classic_h_pic_05": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_05.png"),
// 				"tugofwar_classic_h_pic_06": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_06.png"),
// 				"tugofwar_classic_h_pic_07": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_07.png"),
// 				"tugofwar_classic_h_pic_08": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_08.jpg"),
// 				"tugofwar_classic_h_pic_09": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_09.png"),
// 				"tugofwar_classic_h_pic_10": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_10.png"),
// 				"tugofwar_classic_h_pic_11": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_11.png"),
// 				"tugofwar_classic_h_pic_12": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_12.jpg"),
// 				"tugofwar_classic_h_pic_13": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_13.png"),
// 				"tugofwar_classic_h_pic_14": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_14.png"),
// 				"tugofwar_classic_h_pic_15": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_15.png"),
// 				"tugofwar_classic_h_pic_16": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_16.png"),
// 				"tugofwar_classic_h_pic_17": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_17.png"),
// 				"tugofwar_classic_h_pic_18": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_18.png"),
// 				"tugofwar_classic_h_pic_19": fmt.Sprintf(route, "classic", "tugofwar_classic_h_pic_19.png"),
// 				"tugofwar_classic_g_pic_01": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_01.png"),
// 				"tugofwar_classic_g_pic_02": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_02.png"),
// 				"tugofwar_classic_g_pic_03": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_03.png"),
// 				"tugofwar_classic_g_pic_04": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_04.png"),
// 				"tugofwar_classic_g_pic_05": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_05.png"),
// 				"tugofwar_classic_g_pic_06": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_06.png"),
// 				"tugofwar_classic_g_pic_07": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_07.jpg"),
// 				"tugofwar_classic_g_pic_08": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_08.png"),
// 				"tugofwar_classic_g_pic_09": fmt.Sprintf(route, "classic", "tugofwar_classic_g_pic_09.png"),
// 				"tugofwar_classic_h_ani_01": fmt.Sprintf(route, "classic", "tugofwar_classic_h_ani_01.png"),
// 				"tugofwar_classic_h_ani_02": fmt.Sprintf(route, "classic", "tugofwar_classic_h_ani_02.png"),
// 				"tugofwar_classic_h_ani_03": fmt.Sprintf(route, "classic", "tugofwar_classic_h_ani_03.png"),
// 				"tugofwar_classic_c_ani_01": fmt.Sprintf(route, "classic", "tugofwar_classic_c_ani_01.png"),

// 				"tugofwar_school_h_pic_01": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_01.png"),
// 				"tugofwar_school_h_pic_02": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_02.png"),
// 				"tugofwar_school_h_pic_03": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_03.png"),
// 				"tugofwar_school_h_pic_04": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_04.png"),
// 				"tugofwar_school_h_pic_05": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_05.png"),
// 				"tugofwar_school_h_pic_06": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_06.png"),
// 				"tugofwar_school_h_pic_07": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_07.png"),
// 				"tugofwar_school_h_pic_08": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_08.png"),
// 				"tugofwar_school_h_pic_09": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_09.png"),
// 				"tugofwar_school_h_pic_10": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_10.png"),
// 				"tugofwar_school_h_pic_11": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_11.png"),
// 				"tugofwar_school_h_pic_12": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_12.png"),
// 				"tugofwar_school_h_pic_13": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_13.png"),
// 				"tugofwar_school_h_pic_14": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_14.png"),
// 				"tugofwar_school_h_pic_15": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_15.png"),
// 				"tugofwar_school_h_pic_16": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_16.png"),
// 				"tugofwar_school_h_pic_17": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_17.png"),
// 				"tugofwar_school_h_pic_18": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_18.png"),
// 				"tugofwar_school_h_pic_19": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_19.png"),
// 				"tugofwar_school_h_pic_20": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_20.png"),
// 				"tugofwar_school_h_pic_21": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_21.png"),
// 				"tugofwar_school_h_pic_22": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_22.png"),
// 				"tugofwar_school_h_pic_23": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_23.png"),
// 				"tugofwar_school_h_pic_24": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_24.png"),
// 				"tugofwar_school_h_pic_25": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_25.png"),
// 				"tugofwar_school_h_pic_26": fmt.Sprintf(route, "school", "tugofwar_school_h_pic_26.png"),
// 				"tugofwar_school_g_pic_01": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_01.png"),
// 				"tugofwar_school_g_pic_02": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_02.jpg"),
// 				"tugofwar_school_g_pic_03": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_03.png"),
// 				"tugofwar_school_g_pic_04": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_04.png"),
// 				"tugofwar_school_g_pic_05": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_05.png"),
// 				"tugofwar_school_g_pic_06": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_06.png"),
// 				"tugofwar_school_g_pic_07": fmt.Sprintf(route, "school", "tugofwar_school_g_pic_07.png"),
// 				"tugofwar_school_c_pic_01": fmt.Sprintf(route, "school", "tugofwar_school_c_pic_01.png"),
// 				"tugofwar_school_c_pic_02": fmt.Sprintf(route, "school", "tugofwar_school_c_pic_02.png"),
// 				"tugofwar_school_c_pic_03": fmt.Sprintf(route, "school", "tugofwar_school_c_pic_03.png"),
// 				"tugofwar_school_c_pic_04": fmt.Sprintf(route, "school", "tugofwar_school_c_pic_04.png"),
// 				"tugofwar_school_h_ani_01": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_01.png"),
// 				"tugofwar_school_h_ani_02": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_02.png"),
// 				"tugofwar_school_h_ani_03": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_03.png"),
// 				"tugofwar_school_h_ani_04": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_04.png"),
// 				"tugofwar_school_h_ani_05": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_05.png"),
// 				"tugofwar_school_h_ani_06": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_06.png"),
// 				"tugofwar_school_h_ani_07": fmt.Sprintf(route, "school", "tugofwar_school_h_ani_07.png"),

// 				"tugofwar_christmas_h_pic_01": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_01.png"),
// 				"tugofwar_christmas_h_pic_02": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_02.png"),
// 				"tugofwar_christmas_h_pic_03": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_03.png"),
// 				"tugofwar_christmas_h_pic_04": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_04.png"),
// 				"tugofwar_christmas_h_pic_05": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_05.png"),
// 				"tugofwar_christmas_h_pic_06": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_06.png"),
// 				"tugofwar_christmas_h_pic_07": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_07.jpg"),
// 				"tugofwar_christmas_h_pic_08": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_08.png"),
// 				"tugofwar_christmas_h_pic_09": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_09.png"),
// 				"tugofwar_christmas_h_pic_10": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_10.png"),
// 				"tugofwar_christmas_h_pic_11": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_11.png"),
// 				"tugofwar_christmas_h_pic_12": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_12.png"),
// 				"tugofwar_christmas_h_pic_13": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_13.png"),
// 				"tugofwar_christmas_h_pic_14": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_14.png"),
// 				"tugofwar_christmas_h_pic_15": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_15.png"),
// 				"tugofwar_christmas_h_pic_16": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_16.png"),
// 				"tugofwar_christmas_h_pic_17": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_17.png"),
// 				"tugofwar_christmas_h_pic_18": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_18.png"),
// 				"tugofwar_christmas_h_pic_19": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_19.png"),
// 				"tugofwar_christmas_h_pic_20": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_20.png"),
// 				"tugofwar_christmas_h_pic_21": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_pic_21.png"),
// 				"tugofwar_christmas_g_pic_01": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_01.png"),
// 				"tugofwar_christmas_g_pic_02": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_02.png"),
// 				"tugofwar_christmas_g_pic_03": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_03.png"),
// 				"tugofwar_christmas_g_pic_04": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_04.png"),
// 				"tugofwar_christmas_g_pic_05": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_05.png"),
// 				"tugofwar_christmas_g_pic_06": fmt.Sprintf(route, "christmas", "tugofwar_christmas_g_pic_06.jpg"),
// 				"tugofwar_christmas_c_pic_01": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_pic_01.png"),
// 				"tugofwar_christmas_c_pic_02": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_pic_02.png"),
// 				"tugofwar_christmas_c_pic_03": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_pic_03.png"),
// 				"tugofwar_christmas_c_pic_04": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_pic_04.png"),
// 				"tugofwar_christmas_h_ani_01": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_ani_01.png"),
// 				"tugofwar_christmas_h_ani_02": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_ani_02.png"),
// 				"tugofwar_christmas_h_ani_03": fmt.Sprintf(route, "christmas", "tugofwar_christmas_h_ani_03.png"),
// 				"tugofwar_christmas_c_ani_01": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_ani_01.png"),
// 				"tugofwar_christmas_c_ani_02": fmt.Sprintf(route, "christmas", "tugofwar_christmas_c_ani_02.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將鑑定師遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Monopoly_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/monopoly/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "monopoly" {
// 			// 更新預設值
// 			err = db.Table("activity_game_monopoly_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{

// 				// 鑑定師自定義
// 				"monopoly_classic_h_pic_01": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_01.png"),
// 				"monopoly_classic_h_pic_02": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_02.png"),
// 				"monopoly_classic_h_pic_03": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_03.png"),
// 				"monopoly_classic_h_pic_04": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_04.png"),
// 				"monopoly_classic_h_pic_05": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_05.jpg"),
// 				"monopoly_classic_h_pic_06": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_06.png"),
// 				"monopoly_classic_h_pic_07": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_07.png"),
// 				"monopoly_classic_h_pic_08": fmt.Sprintf(route, "classic", "monopoly_classic_h_pic_08.png"),
// 				"monopoly_classic_g_pic_01": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_01.png"),
// 				"monopoly_classic_g_pic_02": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_02.png"),
// 				"monopoly_classic_g_pic_03": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_03.png"),
// 				"monopoly_classic_g_pic_04": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_04.png"),
// 				"monopoly_classic_g_pic_05": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_05.png"),
// 				"monopoly_classic_g_pic_06": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_06.png"),
// 				"monopoly_classic_g_pic_07": fmt.Sprintf(route, "classic", "monopoly_classic_g_pic_07.jpg"),
// 				"monopoly_classic_c_pic_01": fmt.Sprintf(route, "classic", "monopoly_classic_c_pic_01.png"),
// 				"monopoly_classic_c_pic_02": fmt.Sprintf(route, "classic", "monopoly_classic_c_pic_02.png"),
// 				"monopoly_classic_g_ani_01": fmt.Sprintf(route, "classic", "monopoly_classic_g_ani_01.png"),
// 				"monopoly_classic_g_ani_02": fmt.Sprintf(route, "classic", "monopoly_classic_g_ani_02.png"),
// 				"monopoly_classic_c_ani_01": fmt.Sprintf(route, "classic", "monopoly_classic_c_ani_01.png"),

// 				"monopoly_redpack_h_pic_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_01.png"),
// 				"monopoly_redpack_h_pic_02": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_02.png"),
// 				"monopoly_redpack_h_pic_03": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_03.png"),
// 				"monopoly_redpack_h_pic_04": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_04.png"),
// 				"monopoly_redpack_h_pic_05": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_05.png"),
// 				"monopoly_redpack_h_pic_06": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_06.png"),
// 				"monopoly_redpack_h_pic_07": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_07.png"),
// 				"monopoly_redpack_h_pic_08": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_08.jpg"),
// 				"monopoly_redpack_h_pic_09": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_09.png"),
// 				"monopoly_redpack_h_pic_10": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_10.jpg"),
// 				"monopoly_redpack_h_pic_11": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_11.png"),
// 				"monopoly_redpack_h_pic_12": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_12.png"),
// 				"monopoly_redpack_h_pic_13": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_13.png"),
// 				"monopoly_redpack_h_pic_14": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_14.png"),
// 				"monopoly_redpack_h_pic_15": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_15.png"),
// 				"monopoly_redpack_h_pic_16": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_pic_16.png"),
// 				"monopoly_redpack_g_pic_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_01.png"),
// 				"monopoly_redpack_g_pic_02": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_02.png"),
// 				"monopoly_redpack_g_pic_03": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_03.png"),
// 				"monopoly_redpack_g_pic_04": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_04.jpg"),
// 				"monopoly_redpack_g_pic_05": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_05.png"),
// 				"monopoly_redpack_g_pic_06": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_06.png"),
// 				"monopoly_redpack_g_pic_07": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_07.png"),
// 				"monopoly_redpack_g_pic_08": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_08.png"),
// 				"monopoly_redpack_g_pic_09": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_09.png"),
// 				"monopoly_redpack_g_pic_10": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_pic_10.jpg"),
// 				"monopoly_redpack_c_pic_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_c_pic_01.png"),
// 				"monopoly_redpack_c_pic_02": fmt.Sprintf(route, "redpack", "monopoly_redpack_c_pic_02.png"),
// 				"monopoly_redpack_c_pic_03": fmt.Sprintf(route, "redpack", "monopoly_redpack_c_pic_03.png"),
// 				"monopoly_redpack_h_ani_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_ani_01.png"),
// 				"monopoly_redpack_h_ani_02": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_ani_02.png"),
// 				"monopoly_redpack_h_ani_03": fmt.Sprintf(route, "redpack", "monopoly_redpack_h_ani_03.png"),
// 				"monopoly_redpack_g_ani_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_ani_01.png"),
// 				"monopoly_redpack_g_ani_02": fmt.Sprintf(route, "redpack", "monopoly_redpack_g_ani_02.png"),
// 				"monopoly_redpack_C_ani_01": fmt.Sprintf(route, "redpack", "monopoly_redpack_c_ani_01.png"),

// 				"monopoly_newyear_rabbit_h_pic_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_01.png"),
// 				"monopoly_newyear_rabbit_h_pic_02": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_02.png"),
// 				"monopoly_newyear_rabbit_h_pic_03": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_03.png"),
// 				"monopoly_newyear_rabbit_h_pic_04": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_04.png"),
// 				"monopoly_newyear_rabbit_h_pic_05": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_05.png"),
// 				"monopoly_newyear_rabbit_h_pic_06": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_06.png"),
// 				"monopoly_newyear_rabbit_h_pic_07": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_07.png"),
// 				"monopoly_newyear_rabbit_h_pic_08": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_08.png"),
// 				"monopoly_newyear_rabbit_h_pic_09": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_09.png"),
// 				"monopoly_newyear_rabbit_h_pic_10": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_10.png"),
// 				"monopoly_newyear_rabbit_h_pic_11": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_11.png"),
// 				"monopoly_newyear_rabbit_h_pic_12": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_pic_12.png"),
// 				"monopoly_newyear_rabbit_g_pic_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_01.png"),
// 				"monopoly_newyear_rabbit_g_pic_02": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_02.jpg"),
// 				"monopoly_newyear_rabbit_g_pic_03": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_03.png"),
// 				"monopoly_newyear_rabbit_g_pic_04": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_04.png"),
// 				"monopoly_newyear_rabbit_g_pic_05": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_05.png"),
// 				"monopoly_newyear_rabbit_g_pic_06": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_06.png"),
// 				"monopoly_newyear_rabbit_g_pic_07": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_pic_07.png"),
// 				"monopoly_newyear_rabbit_c_pic_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_c_pic_01.png"),
// 				"monopoly_newyear_rabbit_c_pic_02": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_c_pic_02.png"),
// 				"monopoly_newyear_rabbit_c_pic_03": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_c_pic_03.png"),
// 				"monopoly_newyear_rabbit_h_ani_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_ani_01.png"),
// 				"monopoly_newyear_rabbit_h_ani_02": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_h_ani_02.png"),
// 				"monopoly_newyear_rabbit_g_ani_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_ani_01.png"),
// 				"monopoly_newyear_rabbit_g_ani_02": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_g_ani_02.png"),
// 				"monopoly_newyear_rabbit_c_ani_01": fmt.Sprintf(route, "newyear_rabbit", "monopoly_newyear_rabbit_c_ani_01.png"),

// 				"monopoly_sashimi_h_pic_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_pic_01.png"),
// 				"monopoly_sashimi_h_pic_02": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_pic_02.png"),
// 				"monopoly_sashimi_h_pic_03": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_pic_03.jpg"),
// 				"monopoly_sashimi_h_pic_04": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_pic_04.jpg"),
// 				"monopoly_sashimi_h_pic_05": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_pic_05.png"),
// 				"monopoly_sashimi_g_pic_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_01.png"),
// 				"monopoly_sashimi_g_pic_02": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_02.jpg"),
// 				"monopoly_sashimi_g_pic_03": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_03.png"),
// 				"monopoly_sashimi_g_pic_04": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_04.png"),
// 				"monopoly_sashimi_g_pic_05": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_05.png"),
// 				"monopoly_sashimi_g_pic_06": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_pic_06.png"),
// 				"monopoly_sashimi_c_pic_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_c_pic_01.png"),
// 				"monopoly_sashimi_c_pic_02": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_c_pic_02.png"),
// 				"monopoly_sashimi_h_ani_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_ani_01.png"),
// 				"monopoly_sashimi_h_ani_02": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_h_ani_02.png"),
// 				"monopoly_sashimi_g_ani_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_ani_01.png"),
// 				"monopoly_sashimi_g_ani_02": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_g_ani_02.png"),
// 				"monopoly_sashimi_c_ani_01": fmt.Sprintf(route, "sashimi", "monopoly_sashimi_c_ani_01.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將遊戲抽獎自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Lottery_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/lottery/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "lottery" {
// 			// 更新預設值
// 			err = db.Table("activity_game_lottery_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				// 遊戲抽獎自定義
// 				"lottery_jiugongge_classic_h_pic_01": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_h_pic_01.png"),
// 				"lottery_jiugongge_classic_h_pic_02": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_h_pic_02.png"),
// 				"lottery_jiugongge_classic_h_pic_03": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_h_pic_03.jpg"),
// 				"lottery_jiugongge_classic_h_pic_04": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_h_pic_04.png"),
// 				"lottery_jiugongge_classic_g_pic_01": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_g_pic_01.jpg"),
// 				"lottery_jiugongge_classic_g_pic_02": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_g_pic_02.png"),
// 				"lottery_jiugongge_classic_c_pic_01": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_pic_01.png"),
// 				"lottery_jiugongge_classic_c_pic_02": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_pic_02.png"),
// 				"lottery_jiugongge_classic_c_pic_03": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_pic_03.png"),
// 				"lottery_jiugongge_classic_c_pic_04": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_pic_04.png"),
// 				"lottery_jiugongge_classic_c_ani_01": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_ani_01.png"),
// 				"lottery_jiugongge_classic_c_ani_02": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_ani_02.png"),
// 				"lottery_jiugongge_classic_c_ani_03": fmt.Sprintf(route, "classic", "lottery_jiugongge_classic_c_ani_03.png"),

// 				"lottery_turntable_classic_h_pic_01": fmt.Sprintf(route, "classic", "lottery_turntable_classic_h_pic_01.png"),
// 				"lottery_turntable_classic_h_pic_02": fmt.Sprintf(route, "classic", "lottery_turntable_classic_h_pic_02.png"),
// 				"lottery_turntable_classic_h_pic_03": fmt.Sprintf(route, "classic", "lottery_turntable_classic_h_pic_03.jpg"),
// 				"lottery_turntable_classic_h_pic_04": fmt.Sprintf(route, "classic", "lottery_turntable_classic_h_pic_04.png"),
// 				"lottery_turntable_classic_g_pic_01": fmt.Sprintf(route, "classic", "lottery_turntable_classic_g_pic_01.jpg"),
// 				"lottery_turntable_classic_g_pic_02": fmt.Sprintf(route, "classic", "lottery_turntable_classic_g_pic_02.png"),
// 				"lottery_turntable_classic_c_pic_01": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_01.png"),
// 				"lottery_turntable_classic_c_pic_02": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_02.png"),
// 				"lottery_turntable_classic_c_pic_03": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_03.png"),
// 				"lottery_turntable_classic_c_pic_04": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_04.png"),
// 				"lottery_turntable_classic_c_pic_05": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_05.png"),
// 				"lottery_turntable_classic_c_pic_06": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_pic_06.png"),
// 				"lottery_turntable_classic_c_ani_01": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_ani_01.png"),
// 				"lottery_turntable_classic_c_ani_02": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_ani_02.png"),
// 				"lottery_turntable_classic_c_ani_03": fmt.Sprintf(route, "classic", "lottery_turntable_classic_c_ani_03.png"),

// 				"lottery_jiugongge_starrysky_h_pic_01": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_01.png"),
// 				"lottery_jiugongge_starrysky_h_pic_02": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_02.png"),
// 				"lottery_jiugongge_starrysky_h_pic_03": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_03.png"),
// 				"lottery_jiugongge_starrysky_h_pic_04": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_04.png"),
// 				"lottery_jiugongge_starrysky_h_pic_05": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_05.png"),
// 				"lottery_jiugongge_starrysky_h_pic_06": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_06.png"),
// 				"lottery_jiugongge_starrysky_h_pic_07": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_h_pic_07.jpg"),
// 				"lottery_jiugongge_starrysky_g_pic_01": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_g_pic_01.png"),
// 				"lottery_jiugongge_starrysky_g_pic_02": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_g_pic_02.png"),
// 				"lottery_jiugongge_starrysky_g_pic_03": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_g_pic_03.jpg"),
// 				"lottery_jiugongge_starrysky_g_pic_04": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_g_pic_04.png"),
// 				"lottery_jiugongge_starrysky_c_pic_01": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_pic_01.png"),
// 				"lottery_jiugongge_starrysky_c_pic_02": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_pic_02.png"),
// 				"lottery_jiugongge_starrysky_c_pic_03": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_pic_03.png"),
// 				"lottery_jiugongge_starrysky_c_pic_04": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_pic_04.png"),
// 				"lottery_jiugongge_starrysky_c_ani_01": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_01.png"),
// 				"lottery_jiugongge_starrysky_c_ani_02": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_02.png"),
// 				"lottery_jiugongge_starrysky_c_ani_03": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_03.png"),
// 				"lottery_jiugongge_starrysky_c_ani_04": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_04.png"),
// 				"lottery_jiugongge_starrysky_c_ani_05": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_05.png"),
// 				"lottery_jiugongge_starrysky_c_ani_06": fmt.Sprintf(route, "starrysky", "lottery_jiugongge_starrysky_c_ani_06.png"),

// 				"lottery_turntable_starrysky_h_pic_01": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_01.png"),
// 				"lottery_turntable_starrysky_h_pic_02": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_02.png"),
// 				"lottery_turntable_starrysky_h_pic_03": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_03.png"),
// 				"lottery_turntable_starrysky_h_pic_04": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_04.png"),
// 				"lottery_turntable_starrysky_h_pic_05": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_05.png"),
// 				"lottery_turntable_starrysky_h_pic_06": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_06.png"),
// 				"lottery_turntable_starrysky_h_pic_07": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_07.png"),
// 				"lottery_turntable_starrysky_h_pic_08": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_h_pic_08.jpg"),
// 				"lottery_turntable_starrysky_g_pic_01": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_g_pic_01.png"),
// 				"lottery_turntable_starrysky_g_pic_02": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_g_pic_02.png"),
// 				"lottery_turntable_starrysky_g_pic_03": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_g_pic_03.png"),
// 				"lottery_turntable_starrysky_g_pic_04": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_g_pic_04.jpg"),
// 				"lottery_turntable_starrysky_g_pic_05": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_g_pic_05.png"),
// 				"lottery_turntable_starrysky_c_pic_01": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_pic_01.png"),
// 				"lottery_turntable_starrysky_c_pic_02": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_pic_02.png"),
// 				"lottery_turntable_starrysky_c_pic_03": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_pic_03.png"),
// 				"lottery_turntable_starrysky_c_pic_04": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_pic_04.png"),
// 				"lottery_turntable_starrysky_c_pic_05": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_pic_05.png"),
// 				"lottery_turntable_starrysky_c_ani_01": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_01.png"),
// 				"lottery_turntable_starrysky_c_ani_02": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_02.png"),
// 				"lottery_turntable_starrysky_c_ani_03": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_03.png"),
// 				"lottery_turntable_starrysky_c_ani_04": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_04.png"),
// 				"lottery_turntable_starrysky_c_ani_05": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_05.png"),
// 				"lottery_turntable_starrysky_c_ani_06": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_06.png"),
// 				"lottery_turntable_starrysky_c_ani_07": fmt.Sprintf(route, "starrysky", "lottery_turntable_starrysky_c_ani_07.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將套紅包遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Ropepack_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/ropepack/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "ropepack" {
// 			// 更新預設值
// 			err = db.Table("activity_game_ropepack_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{
// 				// 套紅包自定義
// 				"ropepack_classic_h_pic_01": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_01.png"),
// 				"ropepack_classic_h_pic_02": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_02.png"),
// 				"ropepack_classic_h_pic_03": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_03.jpg"),
// 				"ropepack_classic_h_pic_04": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_04.png"),
// 				"ropepack_classic_h_pic_05": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_05.png"),
// 				"ropepack_classic_h_pic_06": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_06.png"),
// 				"ropepack_classic_h_pic_07": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_07.png"),
// 				"ropepack_classic_h_pic_08": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_08.png"),
// 				"ropepack_classic_h_pic_09": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_09.png"),
// 				"ropepack_classic_h_pic_10": fmt.Sprintf(route, "classic", "ropepack_classic_h_pic_10.png"),
// 				"ropepack_classic_g_pic_01": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_01.png"),
// 				"ropepack_classic_g_pic_02": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_02.png"),
// 				"ropepack_classic_g_pic_03": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_03.jpg"),
// 				"ropepack_classic_g_pic_04": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_04.png"),
// 				"ropepack_classic_g_pic_05": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_05.png"),
// 				"ropepack_classic_g_pic_06": fmt.Sprintf(route, "classic", "ropepack_classic_g_pic_06.jpg"),
// 				"ropepack_classic_h_ani_01": fmt.Sprintf(route, "classic", "ropepack_classic_h_ani_01.png"),
// 				"ropepack_classic_g_ani_01": fmt.Sprintf(route, "classic", "ropepack_classic_g_ani_01.png"),
// 				"ropepack_classic_g_ani_02": fmt.Sprintf(route, "classic", "ropepack_classic_g_ani_02.png"),
// 				"ropepack_classic_c_ani_01": fmt.Sprintf(route, "classic", "ropepack_classic_c_ani_01.png"),

// 				"ropepack_newyear_rabbit_h_pic_01": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_01.png"),
// 				"ropepack_newyear_rabbit_h_pic_02": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_02.jpg"),
// 				"ropepack_newyear_rabbit_h_pic_03": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_03.png"),
// 				"ropepack_newyear_rabbit_h_pic_04": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_04.png"),
// 				"ropepack_newyear_rabbit_h_pic_05": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_05.png"),
// 				"ropepack_newyear_rabbit_h_pic_06": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_06.png"),
// 				"ropepack_newyear_rabbit_h_pic_07": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_07.png"),
// 				"ropepack_newyear_rabbit_h_pic_08": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_08.png"),
// 				"ropepack_newyear_rabbit_h_pic_09": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_pic_09.png"),
// 				"ropepack_newyear_rabbit_g_pic_01": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_pic_01.png"),
// 				"ropepack_newyear_rabbit_g_pic_02": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_pic_02.jpg"),
// 				"ropepack_newyear_rabbit_g_pic_03": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_pic_03.png"),
// 				"ropepack_newyear_rabbit_h_ani_01": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_h_ani_01.png"),
// 				"ropepack_newyear_rabbit_g_ani_01": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_ani_01.png"),
// 				"ropepack_newyear_rabbit_g_ani_02": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_ani_02.png"),
// 				"ropepack_newyear_rabbit_g_ani_03": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_g_ani_03.png"),
// 				"ropepack_newyear_rabbit_c_ani_01": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_c_ani_01.png"),
// 				"ropepack_newyear_rabbit_c_ani_02": fmt.Sprintf(route, "newyear_rabbit", "ropepack_newyear_rabbit_c_ani_02.png"),

// 				"ropepack_moonfestival_h_pic_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_01.png"),
// 				"ropepack_moonfestival_h_pic_02": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_02.png"),
// 				"ropepack_moonfestival_h_pic_03": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_03.png"),
// 				"ropepack_moonfestival_h_pic_04": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_04.png"),
// 				"ropepack_moonfestival_h_pic_05": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_05.jpg"),
// 				"ropepack_moonfestival_h_pic_06": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_06.png"),
// 				"ropepack_moonfestival_h_pic_07": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_07.png"),
// 				"ropepack_moonfestival_h_pic_08": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_08.jpg"),
// 				"ropepack_moonfestival_h_pic_09": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_pic_09.png"),
// 				"ropepack_moonfestival_g_pic_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_g_pic_01.png"),
// 				"ropepack_moonfestival_g_pic_02": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_g_pic_02.jpg"),
// 				"ropepack_moonfestival_c_pic_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_c_pic_01.png"),
// 				"ropepack_moonfestival_h_ani_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_h_ani_01.png"),
// 				"ropepack_moonfestival_g_ani_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_g_ani_01.png"),
// 				"ropepack_moonfestival_g_ani_02": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_g_ani_02.png"),
// 				"ropepack_moonfestival_c_ani_01": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_c_ani_01.png"),
// 				"ropepack_moonfestival_c_ani_02": fmt.Sprintf(route, "moonfestival", "ropepack_moonfestival_c_ani_02.png"),

// 				"ropepack_3D_h_pic_01": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_01.png"),
// 				"ropepack_3D_h_pic_02": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_02.png"),
// 				"ropepack_3D_h_pic_03": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_03.png"),
// 				"ropepack_3D_h_pic_04": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_04.jpg"),
// 				"ropepack_3D_h_pic_05": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_05.png"),
// 				"ropepack_3D_h_pic_06": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_06.png"),
// 				"ropepack_3D_h_pic_07": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_07.png"),
// 				"ropepack_3D_h_pic_08": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_08.png"),
// 				"ropepack_3D_h_pic_09": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_09.png"),
// 				"ropepack_3D_h_pic_10": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_10.png"),
// 				"ropepack_3D_h_pic_11": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_11.png"),
// 				"ropepack_3D_h_pic_12": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_12.png"),
// 				"ropepack_3D_h_pic_13": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_13.png"),
// 				"ropepack_3D_h_pic_14": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_14.png"),
// 				"ropepack_3D_h_pic_15": fmt.Sprintf(route, "3D", "ropepack_3D_h_pic_15.png"),
// 				"ropepack_3D_g_pic_01": fmt.Sprintf(route, "3D", "ropepack_3D_g_pic_01.png"),
// 				"ropepack_3D_g_pic_02": fmt.Sprintf(route, "3D", "ropepack_3D_g_pic_02.jpg"),
// 				"ropepack_3D_g_pic_03": fmt.Sprintf(route, "3D", "ropepack_3D_g_pic_03.png"),
// 				"ropepack_3D_g_pic_04": fmt.Sprintf(route, "3D", "ropepack_3D_g_pic_04.png"),
// 				"ropepack_3D_h_ani_01": fmt.Sprintf(route, "3D", "ropepack_3D_h_ani_01.png"),
// 				"ropepack_3D_h_ani_02": fmt.Sprintf(route, "3D", "ropepack_3D_h_ani_02.png"),
// 				"ropepack_3D_h_ani_03": fmt.Sprintf(route, "3D", "ropepack_3D_h_ani_03.png"),
// 				"ropepack_3D_g_ani_01": fmt.Sprintf(route, "3D", "ropepack_3D_g_ani_01.png"),
// 				"ropepack_3D_g_ani_02": fmt.Sprintf(route, "3D", "ropepack_3D_g_ani_02.png"),
// 				"ropepack_3D_c_ani_01": fmt.Sprintf(route, "3D", "ropepack_3D_c_ani_01.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將搖紅包遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Redpack_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/redpack/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "redpack" {
// 			// 更新預設值
// 			err = db.Table("activity_game_redpack_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{

// 				// 搖紅包自定義
// 				"redpack_classic_h_pic_01": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_01.png"),
// 				"redpack_classic_h_pic_02": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_02.jpg"),
// 				"redpack_classic_h_pic_03": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_03.png"),
// 				"redpack_classic_h_pic_04": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_04.png"),
// 				"redpack_classic_h_pic_05": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_05.png"),
// 				"redpack_classic_h_pic_06": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_06.png"),
// 				"redpack_classic_h_pic_07": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_07.png"),
// 				"redpack_classic_h_pic_08": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_08.png"),
// 				"redpack_classic_h_pic_09": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_09.png"),
// 				"redpack_classic_h_pic_10": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_10.png"),
// 				"redpack_classic_h_pic_11": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_11.png"),
// 				"redpack_classic_h_pic_12": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_12.png"),
// 				"redpack_classic_h_pic_13": fmt.Sprintf(route, "classic", "redpack_classic_h_pic_13.jpg"),
// 				"redpack_classic_g_pic_01": fmt.Sprintf(route, "classic", "redpack_classic_g_pic_01.png"),
// 				"redpack_classic_g_pic_02": fmt.Sprintf(route, "classic", "redpack_classic_g_pic_02.jpg"),
// 				"redpack_classic_g_pic_03": fmt.Sprintf(route, "classic", "redpack_classic_g_pic_03.png"),
// 				"redpack_classic_h_ani_01": fmt.Sprintf(route, "classic", "redpack_classic_h_ani_01.png"),
// 				"redpack_classic_h_ani_02": fmt.Sprintf(route, "classic", "redpack_classic_h_ani_02.png"),
// 				"redpack_classic_g_ani_01": fmt.Sprintf(route, "classic", "redpack_classic_g_ani_01.png"),
// 				"redpack_classic_g_ani_02": fmt.Sprintf(route, "classic", "redpack_classic_g_ani_02.png"),
// 				"redpack_classic_g_ani_03": fmt.Sprintf(route, "classic", "redpack_classic_g_ani_03.png"),

// 				"redpack_cherry_h_pic_01": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_01.png"),
// 				"redpack_cherry_h_pic_02": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_02.png"),
// 				"redpack_cherry_h_pic_03": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_03.png"),
// 				"redpack_cherry_h_pic_04": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_04.png"),
// 				"redpack_cherry_h_pic_05": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_05.png"),
// 				"redpack_cherry_h_pic_06": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_06.png"),
// 				"redpack_cherry_h_pic_07": fmt.Sprintf(route, "cherry", "redpack_cherry_h_pic_07.png"),
// 				"redpack_cherry_g_pic_01": fmt.Sprintf(route, "cherry", "redpack_cherry_g_pic_01.png"),
// 				"redpack_cherry_g_pic_02": fmt.Sprintf(route, "cherry", "redpack_cherry_g_pic_02.jpg"),
// 				"redpack_cherry_h_ani_01": fmt.Sprintf(route, "cherry", "redpack_cherry_h_ani_01.png"),
// 				"redpack_cherry_h_ani_02": fmt.Sprintf(route, "cherry", "redpack_cherry_h_ani_02.png"),
// 				"redpack_cherry_g_ani_01": fmt.Sprintf(route, "cherry", "redpack_cherry_g_ani_01.png"),
// 				"redpack_cherry_g_ani_02": fmt.Sprintf(route, "cherry", "redpack_cherry_g_ani_02.png"),

// 				"redpack_christmas_h_pic_01": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_01.png"),
// 				"redpack_christmas_h_pic_02": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_02.jpg"),
// 				"redpack_christmas_h_pic_03": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_03.png"),
// 				"redpack_christmas_h_pic_04": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_04.png"),
// 				"redpack_christmas_h_pic_05": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_05.png"),
// 				"redpack_christmas_h_pic_06": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_06.png"),
// 				"redpack_christmas_h_pic_07": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_07.png"),
// 				"redpack_christmas_h_pic_08": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_08.png"),
// 				"redpack_christmas_h_pic_09": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_09.png"),
// 				"redpack_christmas_h_pic_10": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_10.png"),
// 				"redpack_christmas_h_pic_11": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_11.png"),
// 				"redpack_christmas_h_pic_12": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_12.png"),
// 				"redpack_christmas_h_pic_13": fmt.Sprintf(route, "christmas", "redpack_christmas_h_pic_13.png"),
// 				"redpack_christmas_g_pic_01": fmt.Sprintf(route, "christmas", "redpack_christmas_g_pic_01.png"),
// 				"redpack_christmas_g_pic_02": fmt.Sprintf(route, "christmas", "redpack_christmas_g_pic_02.jpg"),
// 				"redpack_christmas_g_pic_03": fmt.Sprintf(route, "christmas", "redpack_christmas_g_pic_03.png"),
// 				"redpack_christmas_g_pic_04": fmt.Sprintf(route, "christmas", "redpack_christmas_g_pic_04.png"),
// 				"redpack_christmas_c_pic_01": fmt.Sprintf(route, "christmas", "redpack_christmas_c_pic_01.png"),
// 				"redpack_christmas_c_pic_02": fmt.Sprintf(route, "christmas", "redpack_christmas_c_pic_02.png"),
// 				"redpack_christmas_c_ani_01": fmt.Sprintf(route, "christmas", "redpack_christmas_c_ani_01.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }

// 將敲敲樂遊戲自定義欄位填入預設值(更新至正式區時需要執行)
// func Test_Update_Game_Whack_Mole_Customize_Data(t *testing.T) {

// 	// 先取得activity_game所有資料
// 	items, err := db.Table("activity_game").WithConn(conn).
// 		Select("game_id", "game", "topic", "skin").All()
// 	if err != nil {
// 		t.Error("查詢activity_game所有資料發生錯誤")
// 	} else {
// 		t.Log("查詢activity_game所有資料完成")
// 	}

// 	for _, item := range items {
// 		var (
// 			gameID = utils.GetString(item["game_id"], "")
// 			game   = utils.GetString(item["game"], "")
// 			// topics = strings.Split(utils.GetString(item["topic"], ""), "_")
// 			// topic  = ""
// 			route = "/admin/uploads/system/whackmole/%s/%s"
// 		)

// 		// if len(topics) == 2 {
// 		// 	topic = topics[1]
// 		// } else if len(topics) == 3 {
// 		// 	topic = topics[1] + "_" + topics[2]
// 		// }

// 		if game == "whack_mole" {
// 			// 更新預設值
// 			err = db.Table("activity_game_whack_mole_picture").WithConn(conn).
// 				Where("game_id", "=", gameID).Update(command.Value{

// 				// 敲敲樂自定義
// 				"whackmole_classic_h_pic_01": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_01.png"),
// 				"whackmole_classic_h_pic_02": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_02.jpg"),
// 				"whackmole_classic_h_pic_03": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_03.png"),
// 				"whackmole_classic_h_pic_04": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_04.png"),
// 				"whackmole_classic_h_pic_05": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_05.png"),
// 				"whackmole_classic_h_pic_06": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_06.png"),
// 				"whackmole_classic_h_pic_07": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_07.png"),
// 				"whackmole_classic_h_pic_08": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_08.png"),
// 				"whackmole_classic_h_pic_09": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_09.png"),
// 				"whackmole_classic_h_pic_10": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_10.png"),
// 				"whackmole_classic_h_pic_11": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_11.png"),
// 				"whackmole_classic_h_pic_12": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_12.png"),
// 				"whackmole_classic_h_pic_13": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_13.png"),
// 				"whackmole_classic_h_pic_14": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_14.png"),
// 				"whackmole_classic_h_pic_15": fmt.Sprintf(route, "classic", "whackmole_classic_h_pic_15.png"),
// 				"whackmole_classic_g_pic_01": fmt.Sprintf(route, "classic", "whackmole_classic_g_pic_01.png"),
// 				"whackmole_classic_g_pic_02": fmt.Sprintf(route, "classic", "whackmole_classic_g_pic_02.jpg"),
// 				"whackmole_classic_g_pic_03": fmt.Sprintf(route, "classic", "whackmole_classic_g_pic_03.png"),
// 				"whackmole_classic_g_pic_04": fmt.Sprintf(route, "classic", "whackmole_classic_g_pic_04.png"),
// 				"whackmole_classic_g_pic_05": fmt.Sprintf(route, "classic", "whackmole_classic_g_pic_05.png"),
// 				"whackmole_classic_c_pic_01": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_01.png"),
// 				"whackmole_classic_c_pic_02": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_02.png"),
// 				"whackmole_classic_c_pic_03": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_03.png"),
// 				"whackmole_classic_c_pic_04": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_04.png"),
// 				"whackmole_classic_c_pic_05": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_05.png"),
// 				"whackmole_classic_c_pic_06": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_06.png"),
// 				"whackmole_classic_c_pic_07": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_07.png"),
// 				"whackmole_classic_c_pic_08": fmt.Sprintf(route, "classic", "whackmole_classic_c_pic_08.png"),
// 				"whackmole_classic_c_ani_01": fmt.Sprintf(route, "classic", "whackmole_classic_c_ani_01.png"),

// 				"whackmole_halloween_h_pic_01": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_01.png"),
// 				"whackmole_halloween_h_pic_02": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_02.jpg"),
// 				"whackmole_halloween_h_pic_03": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_03.png"),
// 				"whackmole_halloween_h_pic_04": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_04.png"),
// 				"whackmole_halloween_h_pic_05": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_05.png"),
// 				"whackmole_halloween_h_pic_06": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_06.png"),
// 				"whackmole_halloween_h_pic_07": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_07.png"),
// 				"whackmole_halloween_h_pic_08": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_08.png"),
// 				"whackmole_halloween_h_pic_09": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_09.png"),
// 				"whackmole_halloween_h_pic_10": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_10.png"),
// 				"whackmole_halloween_h_pic_11": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_11.png"),
// 				"whackmole_halloween_h_pic_12": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_12.png"),
// 				"whackmole_halloween_h_pic_13": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_13.png"),
// 				"whackmole_halloween_h_pic_14": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_14.png"),
// 				"whackmole_halloween_h_pic_15": fmt.Sprintf(route, "halloween", "whackmole_halloween_h_pic_15.png"),
// 				"whackmole_halloween_g_pic_01": fmt.Sprintf(route, "halloween", "whackmole_halloween_g_pic_01.png"),
// 				"whackmole_halloween_g_pic_02": fmt.Sprintf(route, "halloween", "whackmole_halloween_g_pic_02.jpg"),
// 				"whackmole_halloween_g_pic_03": fmt.Sprintf(route, "halloween", "whackmole_halloween_g_pic_03.png"),
// 				"whackmole_halloween_g_pic_04": fmt.Sprintf(route, "halloween", "whackmole_halloween_g_pic_04.png"),
// 				"whackmole_halloween_g_pic_05": fmt.Sprintf(route, "halloween", "whackmole_halloween_g_pic_05.png"),
// 				"whackmole_halloween_c_pic_01": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_01.png"),
// 				"whackmole_halloween_c_pic_02": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_02.png"),
// 				"whackmole_halloween_c_pic_03": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_03.png"),
// 				"whackmole_halloween_c_pic_04": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_04.png"),
// 				"whackmole_halloween_c_pic_05": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_05.png"),
// 				"whackmole_halloween_c_pic_06": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_06.png"),
// 				"whackmole_halloween_c_pic_07": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_07.png"),
// 				"whackmole_halloween_c_pic_08": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_pic_08.png"),
// 				"whackmole_halloween_c_ani_01": fmt.Sprintf(route, "halloween", "whackmole_halloween_c_ani_01.png"),

// 				"whackmole_christmas_h_pic_01": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_01.png"),
// 				"whackmole_christmas_h_pic_02": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_02.png"),
// 				"whackmole_christmas_h_pic_03": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_03.jpg"),
// 				"whackmole_christmas_h_pic_04": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_04.png"),
// 				"whackmole_christmas_h_pic_05": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_05.png"),
// 				"whackmole_christmas_h_pic_06": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_06.png"),
// 				"whackmole_christmas_h_pic_07": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_07.png"),
// 				"whackmole_christmas_h_pic_08": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_08.png"),
// 				"whackmole_christmas_h_pic_09": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_09.png"),
// 				"whackmole_christmas_h_pic_10": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_10.png"),
// 				"whackmole_christmas_h_pic_11": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_11.png"),
// 				"whackmole_christmas_h_pic_12": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_12.png"),
// 				"whackmole_christmas_h_pic_13": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_13.png"),
// 				"whackmole_christmas_h_pic_14": fmt.Sprintf(route, "christmas", "whackmole_christmas_h_pic_14.png"),
// 				"whackmole_christmas_g_pic_01": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_01.png"),
// 				"whackmole_christmas_g_pic_02": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_02.jpg"),
// 				"whackmole_christmas_g_pic_03": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_03.png"),
// 				"whackmole_christmas_g_pic_04": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_04.png"),
// 				"whackmole_christmas_g_pic_05": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_05.png"),
// 				"whackmole_christmas_g_pic_06": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_06.png"),
// 				"whackmole_christmas_g_pic_07": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_07.png"),
// 				"whackmole_christmas_g_pic_08": fmt.Sprintf(route, "christmas", "whackmole_christmas_g_pic_08.png"),
// 				"whackmole_christmas_c_pic_01": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_01.png"),
// 				"whackmole_christmas_c_pic_02": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_02.png"),
// 				"whackmole_christmas_c_pic_03": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_03.png"),
// 				"whackmole_christmas_c_pic_04": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_04.png"),
// 				"whackmole_christmas_c_pic_05": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_05.png"),
// 				"whackmole_christmas_c_pic_06": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_06.png"),
// 				"whackmole_christmas_c_pic_07": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_07.png"),
// 				"whackmole_christmas_c_pic_08": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_pic_08.png"),
// 				"whackmole_christmas_c_ani_01": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_ani_01.png"),
// 				"whackmole_christmas_c_ani_02": fmt.Sprintf(route, "christmas", "whackmole_christmas_c_ani_02.png"),
// 			})
// 			if err != nil {
// 				t.Error("更新自定義預設值發生錯誤")
// 			}
// 		}
// 	}
// }
