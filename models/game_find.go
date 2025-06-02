package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/config"
	"hilive/modules/utils"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

// FindFindGameType 查詢遊戲種類資訊，redis有資料先從reids取得，沒有則呼叫資料表並加入redis中(mongo)
func (a GameModel) FindGameType(isRedis bool, gameID string) (string, error) {
	var (
		game string
		err  error
	)

	if isRedis {
		// 判斷redis裡是否有遊戲種類資訊，有則不執行查詢資料表功能
		game, err = a.RedisConn.GetCache(config.GAME_TYPE_REDIS + gameID)
		if err != nil {
			return "", errors.New("錯誤: 取得遊戲種類快取資料發生問題")
		}
	}

	if game == "" {
		// 查詢mongo的遊戲資料
		result, err := a.MongoConn.FindOne(a.TableName,
			bson.M{"game_id": gameID},
			options.FindOne().SetProjection(bson.M{"game": 1})) // 只取得game欄位
		if err != nil {
			return "", errors.New("錯誤: 無法取得遊戲種類資訊(mongo)，請重新查詢")
		}

		// 取得遊戲種類
		game, _ = result["game"].(string)

		// 將遊戲種類資訊加入redis
		if _, err := a.RedisConn.SetCache(config.GAME_TYPE_REDIS+gameID, game); err != nil {
			return "", errors.New("錯誤: 設置遊戲種類快取資料發生問題")
		}

	}
	return game, nil
}

// FindAll 查詢所有場次遊戲資訊(mongo)
func (a GameModel) FindAll(activityID, game string) ([]GameModel, error) {
	var (
		// 查詢活動下所有遊戲場次
		filter = bson.M{"activity_id": activityID} // 過濾條件

		// sql = a.Table(config.ACTIVITY_GAME_QA_TABLE).
		// 	Select(
		// 		"activity_game_qa.total_qa", "activity_game_qa.qa_second",
		// 		"activity_game_qa.qa_round", "activity_game_qa.qa_people",
		// 	).
		// 	Where("activity_game.activity_id", "=", activityID)
		// games = make([]GameModel, 0)
	)

	if game != "" {
		// mongo過濾game參數
		filter["game"] = bson.M{"$in": strings.Split(game, ",")}

		// mysql過濾game參數
		// sql = sql.WhereIn("game", interfaces(strings.Split(game, ",")))
	}

	// 查詢mysql的所有遊戲場次資料
	// items, err := sql.All()
	// if err != nil {
	// 	return nil, errors.New("錯誤: 無法取得所有場次遊戲資訊(mysql)，請重新查詢")
	// }

	// mysqlModel := MapToGameModelByMysql(items)

	// 查詢mongo的所有遊戲場次資料
	result, err := a.MongoConn.FindMany(a.TableName,
		filter,
		options.Find().SetSort(bson.M{"game_order": 1})) // 依照game_order升冪排列
	if err != nil {
		return nil, errors.New("錯誤: 無法取得所有場次遊戲資訊(mongo)，請重新查詢")
	}

	mongoModel := MapToGameModelByMongo(result)

	// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
	// dataMap := make(map[string]GameModel)
	// for _, m := range mongoModel {
	// 	dataMap[m.GameID] = m
	// }

	// for _, item := range mongoModel {
	// 	if dataItem, found := dataMap[item.GameID]; found {
	// 		// 合併欄位（可依需要擴充）
	// 		item.TotalQA = dataItem.TotalQA
	// 		item.QASecond = dataItem.QASecond
	// 		item.QARound = dataItem.QARound
	// 	}
	// 	games = append(games, item)
	// }

	return mongoModel, nil
}

// Find 查詢遊戲資訊，redis有資料先從reids取得，沒有則呼叫資料表並加入redis中(mongo)
func (a GameModel) Find(isRedis bool, gameID, game string) (GameModel, error) {
	if isRedis {
		var (
			// 自定義圖片參數(陣列)
			hpics              = make([]string, 0)        // 主持靜態
			gpics              = make([]string, 0)        // 玩家靜態
			cpics              = make([]string, 0)        // 共用靜態
			hanis              = make([]string, 0)        // 主持動態
			ganis              = make([]string, 0)        // 玩家動態
			canis              = make([]string, 0)        // 共用動態
			musics             = make([]string, 0)        // 音樂
			questions          = make([]QuestionModel, 0) // 題目資訊
			customizeSceneData map[string]interface{}
		)

		// 判斷redis裡是否有遊戲資訊，有則不執行查詢資料表功能
		dataMap, err := a.RedisConn.HashGetAllCache(config.GAME_REDIS + gameID)
		if err != nil {
			return GameModel{}, errors.New("錯誤: 取得遊戲快取資料發生問題")
		}

		a.ID = utils.GetInt64FromStringMap(dataMap, "id", 0)
		a.UserID, _ = dataMap["user_id"]
		a.ActivityID, _ = dataMap["activity_id"]
		a.GameID, _ = dataMap["game_id"]
		a.Game, _ = dataMap["game"]
		a.Title, _ = dataMap["title"]
		a.GameType, _ = dataMap["game_type"]
		a.LimitTime, _ = dataMap["limit_time"]
		a.Second = utils.GetInt64FromStringMap(dataMap, "second", 0)
		a.MaxPeople = utils.GetInt64FromStringMap(dataMap, "max_people", 0)
		a.People = utils.GetInt64FromStringMap(dataMap, "people", 0)
		a.MaxTimes = utils.GetInt64FromStringMap(dataMap, "max_times", 0)
		a.Allow, _ = dataMap["allow"]
		a.Percent = utils.GetInt64FromStringMap(dataMap, "percent", 0)
		a.FirstPrize = utils.GetInt64FromStringMap(dataMap, "first_prize", 0)
		a.SecondPrize = utils.GetInt64FromStringMap(dataMap, "second_prize", 0)
		a.ThirdPrize = utils.GetInt64FromStringMap(dataMap, "third_prize", 0)
		a.GeneralPrize = utils.GetInt64FromStringMap(dataMap, "general_prize", 0)

		a.Topic, _ = dataMap["topic"]
		a.Skin, _ = dataMap["skin"]
		a.Music, _ = dataMap["music"]
		// 解碼，取得自定義圖片資料
		json.Unmarshal([]byte(dataMap["customize_scene_data"]), &customizeSceneData)
		a.CustomizeSceneData = customizeSceneData

		a.GameRound = utils.GetInt64FromStringMap(dataMap, "game_round", 0)
		a.GameSecond = utils.GetInt64FromStringMap(dataMap, "game_second", 0)
		a.GameStatus, _ = dataMap["game_status"]
		a.ControlGameStatus, _ = dataMap["control_game_status"]
		a.GameAttend = a.RedisConn.SetCard(config.GAME_ATTEND_REDIS + gameID) // 遊戲人數(redis)
		a.DisplayName, _ = dataMap["display_name"]
		a.GameOrder = utils.GetInt64FromStringMap(dataMap, "game_order", 0) // 快問快答人數
		a.BoxReflection, _ = dataMap["box_reflection"]
		a.SamePeople, _ = dataMap["same_people"]

		if a.Game == "QA" {
			// 快問快答
			a.TotalQA = utils.GetInt64FromStringMap(dataMap, "total_qa", 0) // 快問快答總題數
			a.QASecond = utils.GetInt64FromStringMap(dataMap, "qa_second", 0)
			a.QARound = utils.GetInt64FromStringMap(dataMap, "qa_round", 0)   // 快問快答進行題數
			a.QAPeople = utils.GetInt64FromStringMap(dataMap, "qa_people", 0) // 快問快答人數

			// 解碼，取得題目資訊
			json.Unmarshal([]byte(dataMap["questions"]), &questions)
			a.Questions = questions
		}

		if a.Game == "tugofwar" {
			// 拔河遊戲
			a.AllowChooseTeam = utils.GetStringFromStringMap(dataMap, "allow_choose_team", "")                 // 允許玩家選擇隊伍
			a.LeftTeamName = utils.GetStringFromStringMap(dataMap, "left_team_name", "")                       // 左邊隊伍名稱
			a.LeftTeamPicture = utils.GetStringFromStringMap(dataMap, "left_team_picture", "")                 // 左邊隊伍照片
			a.RightTeamName = utils.GetStringFromStringMap(dataMap, "right_team_name", "")                     // 右邊隊伍名稱
			a.RightTeamPicture = utils.GetStringFromStringMap(dataMap, "right_team_picture", "")               // 右邊隊伍照片
			a.LeftTeamGameAttend = a.RedisConn.SetCard(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)   // 左方隊伍參加遊戲人數(redis)
			a.RightTeamGameAttend = a.RedisConn.SetCard(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右方隊伍參加遊戲人數(redis)
			a.Prize = utils.GetStringFromStringMap(dataMap, "prize", "")                                       // 獎品發放
		}

		if a.Game == "bingo" {
			// 賓果遊戲
			a.MaxNumber = utils.GetInt64FromStringMap(dataMap, "max_number", 0)   // 最大號碼
			a.BingoLine = utils.GetInt64FromStringMap(dataMap, "bingo_line", 0)   // 賓果連線數
			a.RoundPrize = utils.GetInt64FromStringMap(dataMap, "round_prize", 0) // 每輪發獎人數
			a.BingoRound = utils.GetInt64FromStringMap(dataMap, "bingo_round", 0) // 賓果回合數
		}

		// 扭蛋機遊戲
		if a.Game == "3DGachaMachine" {
			a.GachaMachineReflection = utils.GetFloat64FromStringMap(dataMap, "gacha_machine_reflection", 0) // 球的反射度
			a.IsShake = utils.GetStringFromStringMap(dataMap, "is_shake", "false")                           // 是否搖晃
			a.ReflectiveSwitch = utils.GetStringFromStringMap(dataMap, "reflective_switch", "")              // 反射開關
		}

		// 投票遊戲
		if a.Game == "vote" {
			a.VoteScreen = utils.GetStringFromStringMap(dataMap, "vote_screen", "")              // 投票畫面(長條圖顯示、排名顯示、詳細資訊顯示)
			a.VoteTimes = utils.GetInt64FromStringMap(dataMap, "vote_times", 0)                  // 人員投票次數
			a.VoteMethod = utils.GetStringFromStringMap(dataMap, "vote_method", "")              // 投票模式(全選項投票)
			a.VoteMethodPlayer = utils.GetStringFromStringMap(dataMap, "vote_method_player", "") // 玩家投票方式(一個選項一票、自由投票)
			a.VoteRestriction = utils.GetStringFromStringMap(dataMap, "vote_restriction", "")    // 投票限制(所有人員都能投票、特殊人員才能投票)
			a.AvatarShape = utils.GetStringFromStringMap(dataMap, "avatar_shape", "")            // 選項框是圓形還是方形
			a.VoteStartTime = utils.GetStringFromStringMap(dataMap, "vote_start_time", "")       // 投票結束時間
			a.VoteEndTime = utils.GetStringFromStringMap(dataMap, "vote_end_time", "")           // 投票結束時間
			a.AutoPlay = utils.GetStringFromStringMap(dataMap, "auto_play", "")                  // 自動輪播
			a.ShowRank = utils.GetStringFromStringMap(dataMap, "show_rank", "")                  // 排名展示
			a.TitleSwitch = utils.GetStringFromStringMap(dataMap, "title_switch", "")            // 場次名稱開關
			a.ArrangementGuest = utils.GetStringFromStringMap(dataMap, "arrangement_guest", "")  // 玩家端選項排列方式
		}

		// 編輯次數
		a.EditTimes = utils.GetInt64FromStringMap(dataMap, "edit_times", 0) // 編輯次數

		// 遊戲基本設置
		a.LotteryGameAllow, _ = dataMap["lottery_game_allow"]
		a.RedpackGameAllow, _ = dataMap["redpack_game_allow"]
		a.RopepackGameAllow, _ = dataMap["ropepack_game_allow"]
		a.WhackMoleGameAllow, _ = dataMap["whack_mole_game_allow"]
		a.MonopolyGameAllow, _ = dataMap["monopoly_game_allow"]
		a.QAGameAllow, _ = dataMap["qa_game_allow"]
		a.DrawNumbersGameAllow, _ = dataMap["draw_numbers_game_allow"]
		a.TugofwarGameAllow, _ = dataMap["tugofwar_game_allow"]
		a.BingoGameAllow, _ = dataMap["bingo_game_allow"]
		a.GachaMachineGameAllow, _ = dataMap["3d_gacha_machine_game_allow"]
		a.VoteGameAllow, _ = dataMap["vote_game_allow"]
		a.AllGameAllow, _ = dataMap["all_game_allow"]

		// 活動資訊(join activity)
		// a.UserID, _ = dataMap["user_id"]
		// a.Device, _ = dataMap["device"]

		// 自定義圖片
		// 解碼，取得自定義圖片資料
		json.Unmarshal([]byte(dataMap["customize_host_pictures"]), &hpics)
		json.Unmarshal([]byte(dataMap["customize_guest_pictures"]), &gpics)
		json.Unmarshal([]byte(dataMap["customize_common_pictures"]), &cpics)
		json.Unmarshal([]byte(dataMap["customize_host_anipictures"]), &hanis)
		json.Unmarshal([]byte(dataMap["customize_guest_anipictures"]), &ganis)
		json.Unmarshal([]byte(dataMap["customize_common_anipictures"]), &canis)
		json.Unmarshal([]byte(dataMap["customize_musics"]), &musics)

		a.CustomizeHostPictures = hpics
		a.CustomizeGuestPictures = gpics
		a.CustomizeCommonPictures = cpics
		a.CustomizeHostAnipictures = hanis
		a.CustomizeGuestAnipictures = ganis
		a.CustomizeCommonAnipictures = canis
		a.CustomizeMusics = musics
	}

	if a.ID == 0 {
		// log.Println("redis無資料")

		var (
			sql = a.Table(config.ACTIVITY_GAME_SETTING_TABLE)

			fields = []string{
				// 基本設置
				"activity_game_setting.lottery_game_allow", "activity_game_setting.redpack_game_allow",
				"activity_game_setting.ropepack_game_allow", "activity_game_setting.whack_mole_game_allow",
				"activity_game_setting.monopoly_game_allow", "activity_game_setting.qa_game_allow",
				"activity_game_setting.draw_numbers_game_allow", "activity_game_setting.tugofwar_game_allow",
				"activity_game_setting.bingo_game_allow",
				"activity_game_setting.all_game_allow",
				"activity_game_setting.3d_gacha_machine_game_allow",
				"activity_game_setting.vote_game_allow",
			}

			filter = bson.M{"game_id": gameID} // 過濾條件
		)

		// 查詢mongo的所有遊戲場次資料
		result, err := a.MongoConn.FindOne(a.TableName, filter)
		if err != nil {
			return GameModel{}, errors.New("錯誤: 無法取得遊戲資訊(mongo)，請重新查詢")
		}

		a = a.MapToModel(result)

		// 查詢mysql的所有遊戲場次資料
		item, err := sql.Select(fields...).Where("activity_id", "=", a.ActivityID).First()
		if err != nil {
			return GameModel{}, errors.New("錯誤: 無法取得遊戲資訊(mysql)，請重新查詢")
		}

		// 將mysql裡的參數資料寫入model中
		utils.ApplyInterfaceMapToStruct(item, &a)

		// 將遊戲資訊加入redis
		if isRedis {
			values := []interface{}{config.GAME_REDIS + gameID}
			values = append(values, "id", a.ID)
			values = append(values, "user_id", a.UserID)
			values = append(values, "activity_id", a.ActivityID)
			values = append(values, "game_id", a.GameID)
			values = append(values, "game", a.Game)
			values = append(values, "title", a.Title)
			values = append(values, "game_type", a.GameType)
			values = append(values, "limit_time", a.LimitTime)
			values = append(values, "second", a.Second)
			values = append(values, "max_people", a.MaxPeople)
			values = append(values, "people", a.People)
			values = append(values, "max_times", a.MaxTimes)
			values = append(values, "allow", a.Allow)
			values = append(values, "percent", a.Percent)
			values = append(values, "first_prize", a.FirstPrize)
			values = append(values, "second_prize", a.SecondPrize)
			values = append(values, "third_prize", a.ThirdPrize)
			values = append(values, "general_prize", a.GeneralPrize)

			values = append(values, "topic", a.Topic)
			values = append(values, "skin", a.Skin)
			values = append(values, "music", a.Music)
			values = append(values, "customize_scene_data", utils.JSON(a.CustomizeSceneData))

			values = append(values, "game_round", a.GameRound)
			values = append(values, "game_second", a.GameSecond)
			values = append(values, "game_status", a.GameStatus)
			values = append(values, "control_game_status", "")
			values = append(values, "game_attend", a.RedisConn.SetCard(config.GAME_ATTEND_REDIS+gameID)) // 遊戲人數(redis)
			values = append(values, "display_name", a.DisplayName)
			values = append(values, "game_order", a.GameOrder)
			values = append(values, "box_reflection", a.BoxReflection)
			values = append(values, "same_people", a.SamePeople)

			if a.Game == "QA" {
				// 快問快答
				values = append(values, "total_qa", a.TotalQA)
				values = append(values, "qa_second", a.QASecond)
				values = append(values, "qa_round", a.QARound)
				values = append(values, "qa_people", a.QAPeople)

				// 加密，將題目資訊設置至redis中
				values = append(values, "questions", utils.JSON(a.Questions))
			}

			if a.Game == "tugofwar" {
				// 拔河遊戲
				values = append(values, "allow_choose_team", a.AllowChooseTeam)
				values = append(values, "left_team_name", a.LeftTeamName)
				values = append(values, "left_team_picture", a.LeftTeamPicture)
				values = append(values, "right_team_name", a.RightTeamName)
				values = append(values, "right_team_picture", a.RightTeamPicture)
				values = append(values, "left_team_game_attend", a.RedisConn.SetCard(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID))   // 左方隊伍參加遊戲人數(redis)
				values = append(values, "right_team_game_attend", a.RedisConn.SetCard(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID)) // 右方隊伍參加遊戲人數(redis)
				values = append(values, "prize", a.Prize)
			}

			if a.Game == "bingo" {
				// 賓果遊戲
				values = append(values, "max_number", a.MaxNumber)
				values = append(values, "bingo_line", a.BingoLine)
				values = append(values, "round_prize", a.RoundPrize)
				values = append(values, "bingo_round", a.BingoRound)
			}

			// 扭蛋機遊戲
			if a.Game == "3DGachaMachine" {
				a.IsShake = "false"
				values = append(values, "gacha_machine_reflection", a.GachaMachineReflection)
				values = append(values, "is_shake", "false")
				values = append(values, "reflective_switch", a.ReflectiveSwitch)
			}

			// 投票遊戲
			if a.Game == "vote" {
				values = append(values, "vote_screen", a.VoteScreen)              // 投票畫面(長條圖顯示、排名顯示、詳細資訊顯示)
				values = append(values, "vote_times", a.VoteTimes)                // 人員投票次數
				values = append(values, "vote_method", a.VoteMethod)              // 投票模式(全選項投票)
				values = append(values, "vote_method_player", a.VoteMethodPlayer) // 玩家投票方式(一個選項一票、自由投票)
				values = append(values, "vote_restriction", a.VoteRestriction)    // 投票限制(所有人員都能投票、特殊人員才能投票)
				values = append(values, "avatar_shape", a.AvatarShape)            // 選項框是圓形還是方形
				values = append(values, "vote_start_time", a.VoteStartTime)       // 投票開始時間
				values = append(values, "vote_end_time", a.VoteEndTime)           // 投票結束時間
				values = append(values, "auto_play", a.AutoPlay)                  // 自動輪播
				values = append(values, "show_rank", a.ShowRank)                  // 排名展示
				values = append(values, "title_switch", a.TitleSwitch)            // 場次名稱開關
				values = append(values, "arrangement_guest", a.ArrangementGuest)  // 玩家端選項排列方式
			}

			// 編輯次數
			values = append(values, "edit_times", a.EditTimes)

			// 基本設置
			values = append(values, "lottery_game_allow", a.LotteryGameAllow)
			values = append(values, "redpack_game_allow", a.RedpackGameAllow)
			values = append(values, "ropepack_game_allow", a.RopepackGameAllow)
			values = append(values, "whack_mole_game_allow", a.WhackMoleGameAllow)
			values = append(values, "monopoly_game_allow", a.MonopolyGameAllow)
			values = append(values, "qa_game_allow", a.QAGameAllow)
			values = append(values, "draw_numbers_game_allow", a.DrawNumbersGameAllow)
			values = append(values, "tugofwar_game_allow", a.TugofwarGameAllow)
			values = append(values, "bingo_game_allow", a.BingoGameAllow)
			values = append(values, "3d_gacha_machine_game_allow", a.GachaMachineGameAllow)
			values = append(values, "vote_game_allow", a.VoteGameAllow)
			values = append(values, "all_game_allow", a.AllGameAllow)
			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
				return a, errors.New("錯誤: 設置遊戲快取資料發生問題")
			}

			// 自定義圖片
			// 加密，將自定義圖片資料設置至redis中
			values = append(values, "customize_host_pictures", utils.JSON(a.CustomizeHostPictures))
			values = append(values, "customize_guest_pictures", utils.JSON(a.CustomizeGuestPictures))
			values = append(values, "customize_common_pictures", utils.JSON(a.CustomizeCommonPictures))
			values = append(values, "customize_host_anipictures", utils.JSON(a.CustomizeHostAnipictures))
			values = append(values, "customize_guest_anipictures", utils.JSON(a.CustomizeGuestAnipictures))
			values = append(values, "customize_common_anipictures", utils.JSON(a.CustomizeCommonAnipictures))
			values = append(values, "customize_musics", utils.JSON(a.CustomizeMusics))

			// 遊戲資訊
			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
				return a, errors.New("錯誤: 設置遊戲快取資料發生問題")
			}

			// 遊戲種類
			if _, err := a.RedisConn.SetCache(config.GAME_TYPE_REDIS+gameID, a.Game); err != nil {
				return a, errors.New("錯誤: 設置遊戲種類快取資料發生問題")
			}

			// 設置過期時間
			// a.RedisConn.SetExpire(config.GAME_REDIS+gameID,
			// 	config.REDIS_EXPIRE)

			// 設置過期時間
			// a.RedisConn.SetExpire(config.GAME_TYPE_REDIS+gameID,
			// 	config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		}
	}
	return a, nil
}

// MapToModel 將值設置至GameModel
func (a GameModel) MapToModel(m map[string]interface{}) GameModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	// a.LeftTeamPicture, _ = m["left_team_picture"].(string)
	if !strings.Contains(a.LeftTeamPicture, "system") {
		a.LeftTeamPicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + a.LeftTeamPicture
	}

	// a.RightTeamPicture, _ = m["right_team_picture"].(string)
	if !strings.Contains(a.RightTeamPicture, "system") {
		a.RightTeamPicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + a.RightTeamPicture
	}

	// 自定義圖片參數(陣列)
	var (
		hpics  = make([]string, 0) // 主持靜態
		gpics  = make([]string, 0) // 玩家靜態
		cpics  = make([]string, 0) // 共用靜態
		hanis  = make([]string, 0) // 主持動態
		ganis  = make([]string, 0) // 玩家動態
		canis  = make([]string, 0) // 共用動態
		musics = make([]string, 0)
	)

	if a.Skin == "classic" || a.Skin == "customize" {
		if a.Game == "whack_mole" {
			// 敲敲樂自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.WhackmoleClassicHPic01)
					hpics = append(hpics, a.WhackmoleClassicHPic02)
					hpics = append(hpics, a.WhackmoleClassicHPic03)
					hpics = append(hpics, a.WhackmoleClassicHPic04)
					hpics = append(hpics, a.WhackmoleClassicHPic05)
					hpics = append(hpics, a.WhackmoleClassicHPic06)
					hpics = append(hpics, a.WhackmoleClassicHPic07)
					hpics = append(hpics, a.WhackmoleClassicHPic08)
					hpics = append(hpics, a.WhackmoleClassicHPic09)
					hpics = append(hpics, a.WhackmoleClassicHPic10)
					hpics = append(hpics, a.WhackmoleClassicHPic11)
					hpics = append(hpics, a.WhackmoleClassicHPic12)
					hpics = append(hpics, a.WhackmoleClassicHPic13)
					hpics = append(hpics, a.WhackmoleClassicHPic14)
					hpics = append(hpics, a.WhackmoleClassicHPic15)
					gpics = append(gpics, a.WhackmoleClassicGPic01)
					gpics = append(gpics, a.WhackmoleClassicGPic02)
					gpics = append(gpics, a.WhackmoleClassicGPic03)
					gpics = append(gpics, a.WhackmoleClassicGPic04)
					gpics = append(gpics, a.WhackmoleClassicGPic05)
					cpics = append(cpics, a.WhackmoleClassicCPic01)
					cpics = append(cpics, a.WhackmoleClassicCPic02)
					cpics = append(cpics, a.WhackmoleClassicCPic03)
					cpics = append(cpics, a.WhackmoleClassicCPic04)
					cpics = append(cpics, a.WhackmoleClassicCPic05)
					cpics = append(cpics, a.WhackmoleClassicCPic06)
					cpics = append(cpics, a.WhackmoleClassicCPic07)
					cpics = append(cpics, a.WhackmoleClassicCPic08)
					canis = append(canis, a.WhackmoleClassicCAni01)
				} else if a.Topic == "02_halloween" {
					hpics = append(hpics, a.WhackmoleHalloweenHPic01)
					hpics = append(hpics, a.WhackmoleHalloweenHPic02)
					hpics = append(hpics, a.WhackmoleHalloweenHPic03)
					hpics = append(hpics, a.WhackmoleHalloweenHPic04)
					hpics = append(hpics, a.WhackmoleHalloweenHPic05)
					hpics = append(hpics, a.WhackmoleHalloweenHPic06)
					hpics = append(hpics, a.WhackmoleHalloweenHPic07)
					hpics = append(hpics, a.WhackmoleHalloweenHPic08)
					hpics = append(hpics, a.WhackmoleHalloweenHPic09)
					hpics = append(hpics, a.WhackmoleHalloweenHPic10)
					hpics = append(hpics, a.WhackmoleHalloweenHPic11)
					hpics = append(hpics, a.WhackmoleHalloweenHPic12)
					hpics = append(hpics, a.WhackmoleHalloweenHPic13)
					hpics = append(hpics, a.WhackmoleHalloweenHPic14)
					hpics = append(hpics, a.WhackmoleHalloweenHPic15)
					gpics = append(gpics, a.WhackmoleHalloweenGPic01)
					gpics = append(gpics, a.WhackmoleHalloweenGPic02)
					gpics = append(gpics, a.WhackmoleHalloweenGPic03)
					gpics = append(gpics, a.WhackmoleHalloweenGPic04)
					gpics = append(gpics, a.WhackmoleHalloweenGPic05)
					cpics = append(cpics, a.WhackmoleHalloweenCPic01)
					cpics = append(cpics, a.WhackmoleHalloweenCPic02)
					cpics = append(cpics, a.WhackmoleHalloweenCPic03)
					cpics = append(cpics, a.WhackmoleHalloweenCPic04)
					cpics = append(cpics, a.WhackmoleHalloweenCPic05)
					cpics = append(cpics, a.WhackmoleHalloweenCPic06)
					cpics = append(cpics, a.WhackmoleHalloweenCPic07)
					cpics = append(cpics, a.WhackmoleHalloweenCPic08)
					canis = append(canis, a.WhackmoleHalloweenCAni01)
				} else if a.Topic == "03_christmas" {
					hpics = append(hpics, a.WhackmoleChristmasHPic01)
					hpics = append(hpics, a.WhackmoleChristmasHPic02)
					hpics = append(hpics, a.WhackmoleChristmasHPic03)
					hpics = append(hpics, a.WhackmoleChristmasHPic04)
					hpics = append(hpics, a.WhackmoleChristmasHPic05)
					hpics = append(hpics, a.WhackmoleChristmasHPic06)
					hpics = append(hpics, a.WhackmoleChristmasHPic07)
					hpics = append(hpics, a.WhackmoleChristmasHPic08)
					hpics = append(hpics, a.WhackmoleChristmasHPic09)
					hpics = append(hpics, a.WhackmoleChristmasHPic10)
					hpics = append(hpics, a.WhackmoleChristmasHPic11)
					hpics = append(hpics, a.WhackmoleChristmasHPic12)
					hpics = append(hpics, a.WhackmoleChristmasHPic13)
					hpics = append(hpics, a.WhackmoleChristmasHPic14)
					gpics = append(gpics, a.WhackmoleChristmasGPic01)
					gpics = append(gpics, a.WhackmoleChristmasGPic02)
					gpics = append(gpics, a.WhackmoleChristmasGPic03)
					gpics = append(gpics, a.WhackmoleChristmasGPic04)
					gpics = append(gpics, a.WhackmoleChristmasGPic05)
					gpics = append(gpics, a.WhackmoleChristmasGPic06)
					gpics = append(gpics, a.WhackmoleChristmasGPic07)
					gpics = append(gpics, a.WhackmoleChristmasGPic08)
					cpics = append(cpics, a.WhackmoleChristmasCPic01)
					cpics = append(cpics, a.WhackmoleChristmasCPic02)
					cpics = append(cpics, a.WhackmoleChristmasCPic03)
					cpics = append(cpics, a.WhackmoleChristmasCPic04)
					cpics = append(cpics, a.WhackmoleChristmasCPic05)
					cpics = append(cpics, a.WhackmoleChristmasCPic06)
					cpics = append(cpics, a.WhackmoleChristmasCPic07)
					cpics = append(cpics, a.WhackmoleChristmasCPic08)
					canis = append(canis, a.WhackmoleChristmasCAni01)
					canis = append(canis, a.WhackmoleChristmasCAni02)
				}
			}
			// 音樂
			musics = append(musics, a.WhackmoleBgmStart)
			musics = append(musics, a.WhackmoleBgmGaming)
			musics = append(musics, a.WhackmoleBgmEnd)
		} else if a.Game == "draw_numbers" {
			// 搖號抽獎自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.DrawNumbersClassicHPic01)
					hpics = append(hpics, a.DrawNumbersClassicHPic02)
					hpics = append(hpics, a.DrawNumbersClassicHPic03)
					hpics = append(hpics, a.DrawNumbersClassicHPic04)
					hpics = append(hpics, a.DrawNumbersClassicHPic05)
					hpics = append(hpics, a.DrawNumbersClassicHPic06)
					hpics = append(hpics, a.DrawNumbersClassicHPic07)
					hpics = append(hpics, a.DrawNumbersClassicHPic08)
					hpics = append(hpics, a.DrawNumbersClassicHPic09)
					hpics = append(hpics, a.DrawNumbersClassicHPic10)
					hpics = append(hpics, a.DrawNumbersClassicHPic11)
					hpics = append(hpics, a.DrawNumbersClassicHPic12)
					hpics = append(hpics, a.DrawNumbersClassicHPic13)
					hpics = append(hpics, a.DrawNumbersClassicHPic14)
					hpics = append(hpics, a.DrawNumbersClassicHPic15)
					hpics = append(hpics, a.DrawNumbersClassicHPic16)
					hanis = append(hanis, a.DrawNumbersClassicHAni01)
				} else if a.Topic == "02_gold" {
					hpics = append(hpics, a.DrawNumbersGoldHPic01)
					hpics = append(hpics, a.DrawNumbersGoldHPic02)
					hpics = append(hpics, a.DrawNumbersGoldHPic03)
					hpics = append(hpics, a.DrawNumbersGoldHPic04)
					hpics = append(hpics, a.DrawNumbersGoldHPic05)
					hpics = append(hpics, a.DrawNumbersGoldHPic06)
					hpics = append(hpics, a.DrawNumbersGoldHPic07)
					hpics = append(hpics, a.DrawNumbersGoldHPic08)
					hpics = append(hpics, a.DrawNumbersGoldHPic09)
					hpics = append(hpics, a.DrawNumbersGoldHPic10)
					hpics = append(hpics, a.DrawNumbersGoldHPic11)
					hpics = append(hpics, a.DrawNumbersGoldHPic12)
					hpics = append(hpics, a.DrawNumbersGoldHPic13)
					hpics = append(hpics, a.DrawNumbersGoldHPic14)
					hanis = append(hanis, a.DrawNumbersGoldHAni01)
					hanis = append(hanis, a.DrawNumbersGoldHAni02)
					hanis = append(hanis, a.DrawNumbersGoldHAni03)
				} else if a.Topic == "03_newyear_dragon" {
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic01)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic02)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic03)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic04)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic05)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic06)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic07)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic08)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic09)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic10)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic11)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic12)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic13)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic14)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic15)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic16)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic17)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic18)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic19)
					hpics = append(hpics, a.DrawNumbersNewyearDragonHPic20)
					hanis = append(hanis, a.DrawNumbersNewyearDragonHAni01)
					hanis = append(hanis, a.DrawNumbersNewyearDragonHAni02)
				} else if a.Topic == "04_cherry" {
					hpics = append(hpics, a.DrawNumbersCherryHPic01)
					hpics = append(hpics, a.DrawNumbersCherryHPic02)
					hpics = append(hpics, a.DrawNumbersCherryHPic03)
					hpics = append(hpics, a.DrawNumbersCherryHPic04)
					hpics = append(hpics, a.DrawNumbersCherryHPic05)
					hpics = append(hpics, a.DrawNumbersCherryHPic06)
					hpics = append(hpics, a.DrawNumbersCherryHPic07)
					hpics = append(hpics, a.DrawNumbersCherryHPic08)
					hpics = append(hpics, a.DrawNumbersCherryHPic09)
					hpics = append(hpics, a.DrawNumbersCherryHPic10)
					hpics = append(hpics, a.DrawNumbersCherryHPic11)
					hpics = append(hpics, a.DrawNumbersCherryHPic12)
					hpics = append(hpics, a.DrawNumbersCherryHPic13)
					hpics = append(hpics, a.DrawNumbersCherryHPic14)
					hpics = append(hpics, a.DrawNumbersCherryHPic15)
					hpics = append(hpics, a.DrawNumbersCherryHPic16)
					hpics = append(hpics, a.DrawNumbersCherryHPic17)
					hanis = append(hanis, a.DrawNumbersCherryHAni01)
					hanis = append(hanis, a.DrawNumbersCherryHAni02)
					hanis = append(hanis, a.DrawNumbersCherryHAni03)
					hanis = append(hanis, a.DrawNumbersCherryHAni04)
				} else if a.Topic == "05_3D_space" {
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic01)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic02)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic03)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic04)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic05)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic06)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic07)
					hpics = append(hpics, a.DrawNumbers3DSpaceHPic08)
				}
			}
			// 音樂
			musics = append(musics, a.DrawNumbersBgmGaming)
		} else if a.Game == "monopoly" {
			// 鑑定師自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.MonopolyClassicHPic01)
					hpics = append(hpics, a.MonopolyClassicHPic02)
					hpics = append(hpics, a.MonopolyClassicHPic03)
					hpics = append(hpics, a.MonopolyClassicHPic04)
					hpics = append(hpics, a.MonopolyClassicHPic05)
					hpics = append(hpics, a.MonopolyClassicHPic06)
					hpics = append(hpics, a.MonopolyClassicHPic07)
					hpics = append(hpics, a.MonopolyClassicHPic08)
					gpics = append(gpics, a.MonopolyClassicGPic01)
					gpics = append(gpics, a.MonopolyClassicGPic02)
					gpics = append(gpics, a.MonopolyClassicGPic03)
					gpics = append(gpics, a.MonopolyClassicGPic04)
					gpics = append(gpics, a.MonopolyClassicGPic05)
					gpics = append(gpics, a.MonopolyClassicGPic06)
					gpics = append(gpics, a.MonopolyClassicGPic07)
					cpics = append(cpics, a.MonopolyClassicCPic01)
					cpics = append(cpics, a.MonopolyClassicCPic02)
					ganis = append(ganis, a.MonopolyClassicGAni01)
					ganis = append(ganis, a.MonopolyClassicGAni02)
					canis = append(canis, a.MonopolyClassicCAni01)
				} else if a.Topic == "02_redpack" {
					hpics = append(hpics, a.MonopolyRedpackHPic01)
					hpics = append(hpics, a.MonopolyRedpackHPic02)
					hpics = append(hpics, a.MonopolyRedpackHPic03)
					hpics = append(hpics, a.MonopolyRedpackHPic04)
					hpics = append(hpics, a.MonopolyRedpackHPic05)
					hpics = append(hpics, a.MonopolyRedpackHPic06)
					hpics = append(hpics, a.MonopolyRedpackHPic07)
					hpics = append(hpics, a.MonopolyRedpackHPic08)
					hpics = append(hpics, a.MonopolyRedpackHPic09)
					hpics = append(hpics, a.MonopolyRedpackHPic10)
					hpics = append(hpics, a.MonopolyRedpackHPic11)
					hpics = append(hpics, a.MonopolyRedpackHPic12)
					hpics = append(hpics, a.MonopolyRedpackHPic13)
					hpics = append(hpics, a.MonopolyRedpackHPic14)
					hpics = append(hpics, a.MonopolyRedpackHPic15)
					hpics = append(hpics, a.MonopolyRedpackHPic16)
					gpics = append(gpics, a.MonopolyRedpackGPic01)
					gpics = append(gpics, a.MonopolyRedpackGPic02)
					gpics = append(gpics, a.MonopolyRedpackGPic03)
					gpics = append(gpics, a.MonopolyRedpackGPic04)
					gpics = append(gpics, a.MonopolyRedpackGPic05)
					gpics = append(gpics, a.MonopolyRedpackGPic06)
					gpics = append(gpics, a.MonopolyRedpackGPic07)
					gpics = append(gpics, a.MonopolyRedpackGPic08)
					gpics = append(gpics, a.MonopolyRedpackGPic09)
					gpics = append(gpics, a.MonopolyRedpackGPic10)
					cpics = append(cpics, a.MonopolyRedpackCPic01)
					cpics = append(cpics, a.MonopolyRedpackCPic02)
					cpics = append(cpics, a.MonopolyRedpackCPic03)
					hanis = append(hanis, a.MonopolyRedpackHAni01)
					hanis = append(hanis, a.MonopolyRedpackHAni02)
					hanis = append(hanis, a.MonopolyRedpackHAni03)
					ganis = append(ganis, a.MonopolyRedpackGAni01)
					ganis = append(ganis, a.MonopolyRedpackGAni02)
					canis = append(canis, a.MonopolyRedpackCAni01)
				} else if a.Topic == "03_newyear_rabbit" {
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic01)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic02)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic03)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic04)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic05)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic06)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic07)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic08)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic09)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic10)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic11)
					hpics = append(hpics, a.MonopolyNewyearRabbitHPic12)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic01)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic02)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic03)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic04)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic05)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic06)
					gpics = append(gpics, a.MonopolyNewyearRabbitGPic07)
					cpics = append(cpics, a.MonopolyNewyearRabbitCPic01)
					cpics = append(cpics, a.MonopolyNewyearRabbitCPic02)
					cpics = append(cpics, a.MonopolyNewyearRabbitCPic03)
					hanis = append(hanis, a.MonopolyNewyearRabbitHAni01)
					hanis = append(hanis, a.MonopolyNewyearRabbitHAni02)
					ganis = append(ganis, a.MonopolyNewyearRabbitGAni01)
					ganis = append(ganis, a.MonopolyNewyearRabbitGAni02)
					canis = append(canis, a.MonopolyNewyearRabbitCAni01)
				} else if a.Topic == "04_sashimi" {
					hpics = append(hpics, a.MonopolySashimiHPic01)
					hpics = append(hpics, a.MonopolySashimiHPic02)
					hpics = append(hpics, a.MonopolySashimiHPic03)
					hpics = append(hpics, a.MonopolySashimiHPic04)
					hpics = append(hpics, a.MonopolySashimiHPic05)
					gpics = append(gpics, a.MonopolySashimiGPic01)
					gpics = append(gpics, a.MonopolySashimiGPic02)
					gpics = append(gpics, a.MonopolySashimiGPic03)
					gpics = append(gpics, a.MonopolySashimiGPic04)
					gpics = append(gpics, a.MonopolySashimiGPic05)
					gpics = append(gpics, a.MonopolySashimiGPic06)
					cpics = append(cpics, a.MonopolySashimiCPic01)
					cpics = append(cpics, a.MonopolySashimiCPic02)
					hanis = append(hanis, a.MonopolySashimiHAni01)
					hanis = append(hanis, a.MonopolySashimiHAni02)
					ganis = append(ganis, a.MonopolySashimiGAni01)
					ganis = append(ganis, a.MonopolySashimiGAni02)
					canis = append(canis, a.MonopolySashimiCAni01)
				}
			}
			// 音樂
			musics = append(musics, a.MonopolyBgmStart)
			musics = append(musics, a.MonopolyBgmGaming)
			musics = append(musics, a.MonopolyBgmEnd)
		} else if a.Game == "QA" {
			// 快問快答自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.QAClassicHPic01)
					hpics = append(hpics, a.QAClassicHPic02)
					hpics = append(hpics, a.QAClassicHPic03)
					hpics = append(hpics, a.QAClassicHPic04)
					hpics = append(hpics, a.QAClassicHPic05)
					hpics = append(hpics, a.QAClassicHPic06)
					hpics = append(hpics, a.QAClassicHPic07)
					hpics = append(hpics, a.QAClassicHPic08)
					hpics = append(hpics, a.QAClassicHPic09)
					hpics = append(hpics, a.QAClassicHPic10)
					hpics = append(hpics, a.QAClassicHPic11)
					hpics = append(hpics, a.QAClassicHPic12)
					hpics = append(hpics, a.QAClassicHPic13)
					hpics = append(hpics, a.QAClassicHPic14)
					hpics = append(hpics, a.QAClassicHPic15)
					hpics = append(hpics, a.QAClassicHPic16)
					hpics = append(hpics, a.QAClassicHPic17)
					hpics = append(hpics, a.QAClassicHPic18)
					hpics = append(hpics, a.QAClassicHPic19)
					hpics = append(hpics, a.QAClassicHPic20)
					hpics = append(hpics, a.QAClassicHPic21)
					hpics = append(hpics, a.QAClassicHPic22)
					gpics = append(gpics, a.QAClassicGPic01)
					gpics = append(gpics, a.QAClassicGPic02)
					gpics = append(gpics, a.QAClassicGPic03)
					gpics = append(gpics, a.QAClassicGPic04)
					gpics = append(gpics, a.QAClassicGPic05)
					cpics = append(cpics, a.QAClassicCPic01)
					hanis = append(hanis, a.QAClassicHAni01)
					hanis = append(hanis, a.QAClassicHAni02)
					ganis = append(ganis, a.QAClassicGAni01)
					ganis = append(ganis, a.QAClassicGAni02)
				} else if a.Topic == "02_electric" {
					hpics = append(hpics, a.QAElectricHPic01)
					hpics = append(hpics, a.QAElectricHPic02)
					hpics = append(hpics, a.QAElectricHPic03)
					hpics = append(hpics, a.QAElectricHPic04)
					hpics = append(hpics, a.QAElectricHPic05)
					hpics = append(hpics, a.QAElectricHPic06)
					hpics = append(hpics, a.QAElectricHPic07)
					hpics = append(hpics, a.QAElectricHPic08)
					hpics = append(hpics, a.QAElectricHPic09)
					hpics = append(hpics, a.QAElectricHPic10)
					hpics = append(hpics, a.QAElectricHPic11)
					hpics = append(hpics, a.QAElectricHPic12)
					hpics = append(hpics, a.QAElectricHPic13)
					hpics = append(hpics, a.QAElectricHPic14)
					hpics = append(hpics, a.QAElectricHPic15)
					hpics = append(hpics, a.QAElectricHPic16)
					hpics = append(hpics, a.QAElectricHPic17)
					hpics = append(hpics, a.QAElectricHPic18)
					hpics = append(hpics, a.QAElectricHPic19)
					hpics = append(hpics, a.QAElectricHPic20)
					hpics = append(hpics, a.QAElectricHPic21)
					hpics = append(hpics, a.QAElectricHPic22)
					hpics = append(hpics, a.QAElectricHPic23)
					hpics = append(hpics, a.QAElectricHPic24)
					hpics = append(hpics, a.QAElectricHPic25)
					hpics = append(hpics, a.QAElectricHPic26)
					gpics = append(gpics, a.QAElectricGPic01)
					gpics = append(gpics, a.QAElectricGPic02)
					gpics = append(gpics, a.QAElectricGPic03)
					gpics = append(gpics, a.QAElectricGPic04)
					gpics = append(gpics, a.QAElectricGPic05)
					gpics = append(gpics, a.QAElectricGPic06)
					gpics = append(gpics, a.QAElectricGPic07)
					gpics = append(gpics, a.QAElectricGPic08)
					gpics = append(gpics, a.QAElectricGPic09)
					cpics = append(cpics, a.QAElectricCPic01)
					hanis = append(hanis, a.QAElectricHAni01)
					hanis = append(hanis, a.QAElectricHAni02)
					hanis = append(hanis, a.QAElectricHAni03)
					hanis = append(hanis, a.QAElectricHAni04)
					hanis = append(hanis, a.QAElectricHAni05)
					ganis = append(ganis, a.QAElectricGAni01)
					ganis = append(ganis, a.QAElectricGAni02)
					canis = append(canis, a.QAElectricCAni01)
				} else if a.Topic == "03_moonfestival" {
					hpics = append(hpics, a.QAMoonfestivalHPic01)
					hpics = append(hpics, a.QAMoonfestivalHPic02)
					hpics = append(hpics, a.QAMoonfestivalHPic03)
					hpics = append(hpics, a.QAMoonfestivalHPic04)
					hpics = append(hpics, a.QAMoonfestivalHPic05)
					hpics = append(hpics, a.QAMoonfestivalHPic06)
					hpics = append(hpics, a.QAMoonfestivalHPic07)
					hpics = append(hpics, a.QAMoonfestivalHPic08)
					hpics = append(hpics, a.QAMoonfestivalHPic09)
					hpics = append(hpics, a.QAMoonfestivalHPic10)
					hpics = append(hpics, a.QAMoonfestivalHPic11)
					hpics = append(hpics, a.QAMoonfestivalHPic12)
					hpics = append(hpics, a.QAMoonfestivalHPic13)
					hpics = append(hpics, a.QAMoonfestivalHPic14)
					hpics = append(hpics, a.QAMoonfestivalHPic15)
					hpics = append(hpics, a.QAMoonfestivalHPic16)
					hpics = append(hpics, a.QAMoonfestivalHPic17)
					hpics = append(hpics, a.QAMoonfestivalHPic18)
					hpics = append(hpics, a.QAMoonfestivalHPic19)
					hpics = append(hpics, a.QAMoonfestivalHPic20)
					hpics = append(hpics, a.QAMoonfestivalHPic21)
					hpics = append(hpics, a.QAMoonfestivalHPic22)
					hpics = append(hpics, a.QAMoonfestivalHPic23)
					hpics = append(hpics, a.QAMoonfestivalHPic24)
					gpics = append(gpics, a.QAMoonfestivalGPic01)
					gpics = append(gpics, a.QAMoonfestivalGPic02)
					gpics = append(gpics, a.QAMoonfestivalGPic03)
					gpics = append(gpics, a.QAMoonfestivalGPic04)
					gpics = append(gpics, a.QAMoonfestivalGPic05)
					cpics = append(cpics, a.QAMoonfestivalCPic01)
					cpics = append(cpics, a.QAMoonfestivalCPic02)
					cpics = append(cpics, a.QAMoonfestivalCPic03)
					hanis = append(hanis, a.QAMoonfestivalHAni01)
					hanis = append(hanis, a.QAMoonfestivalHAni02)
					ganis = append(ganis, a.QAMoonfestivalGAni01)
					ganis = append(ganis, a.QAMoonfestivalGAni02)
					ganis = append(ganis, a.QAMoonfestivalGAni03)
				} else if a.Topic == "04_newyear_dragon" {
					hpics = append(hpics, a.QANewyearDragonHPic01)
					hpics = append(hpics, a.QANewyearDragonHPic02)
					hpics = append(hpics, a.QANewyearDragonHPic03)
					hpics = append(hpics, a.QANewyearDragonHPic04)
					hpics = append(hpics, a.QANewyearDragonHPic05)
					hpics = append(hpics, a.QANewyearDragonHPic06)
					hpics = append(hpics, a.QANewyearDragonHPic07)
					hpics = append(hpics, a.QANewyearDragonHPic08)
					hpics = append(hpics, a.QANewyearDragonHPic09)
					hpics = append(hpics, a.QANewyearDragonHPic10)
					hpics = append(hpics, a.QANewyearDragonHPic11)
					hpics = append(hpics, a.QANewyearDragonHPic12)
					hpics = append(hpics, a.QANewyearDragonHPic13)
					hpics = append(hpics, a.QANewyearDragonHPic14)
					hpics = append(hpics, a.QANewyearDragonHPic15)
					hpics = append(hpics, a.QANewyearDragonHPic16)
					hpics = append(hpics, a.QANewyearDragonHPic17)
					hpics = append(hpics, a.QANewyearDragonHPic18)
					hpics = append(hpics, a.QANewyearDragonHPic19)
					hpics = append(hpics, a.QANewyearDragonHPic20)
					hpics = append(hpics, a.QANewyearDragonHPic21)
					hpics = append(hpics, a.QANewyearDragonHPic22)
					hpics = append(hpics, a.QANewyearDragonHPic23)
					hpics = append(hpics, a.QANewyearDragonHPic24)
					gpics = append(gpics, a.QANewyearDragonGPic01)
					gpics = append(gpics, a.QANewyearDragonGPic02)
					gpics = append(gpics, a.QANewyearDragonGPic03)
					gpics = append(gpics, a.QANewyearDragonGPic04)
					gpics = append(gpics, a.QANewyearDragonGPic05)
					gpics = append(gpics, a.QANewyearDragonGPic06)
					cpics = append(cpics, a.QANewyearDragonCPic01)
					hanis = append(hanis, a.QANewyearDragonHAni01)
					hanis = append(hanis, a.QANewyearDragonHAni02)
					ganis = append(ganis, a.QANewyearDragonGAni01)
					ganis = append(ganis, a.QANewyearDragonGAni02)
					ganis = append(ganis, a.QANewyearDragonGAni03)
					canis = append(canis, a.QANewyearDragonCAni01)
				}
			}
			// 音樂
			musics = append(musics, a.QABgmStart)
			musics = append(musics, a.QABgmGaming)
			musics = append(musics, a.QABgmEnd)
		} else if a.Game == "redpack" {
			// 搖紅包自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.RedpackClassicHPic01)
					hpics = append(hpics, a.RedpackClassicHPic02)
					hpics = append(hpics, a.RedpackClassicHPic03)
					hpics = append(hpics, a.RedpackClassicHPic04)
					hpics = append(hpics, a.RedpackClassicHPic05)
					hpics = append(hpics, a.RedpackClassicHPic06)
					hpics = append(hpics, a.RedpackClassicHPic07)
					hpics = append(hpics, a.RedpackClassicHPic08)
					hpics = append(hpics, a.RedpackClassicHPic09)
					hpics = append(hpics, a.RedpackClassicHPic10)
					hpics = append(hpics, a.RedpackClassicHPic11)
					hpics = append(hpics, a.RedpackClassicHPic12)
					hpics = append(hpics, a.RedpackClassicHPic13)
					gpics = append(gpics, a.RedpackClassicGPic01)
					gpics = append(gpics, a.RedpackClassicGPic02)
					gpics = append(gpics, a.RedpackClassicGPic03)
					hanis = append(hanis, a.RedpackClassicHAni01)
					hanis = append(hanis, a.RedpackClassicHAni02)
					ganis = append(ganis, a.RedpackClassicGAni01)
					ganis = append(ganis, a.RedpackClassicGAni02)
					ganis = append(ganis, a.RedpackClassicGAni03)
				} else if a.Topic == "02_cherry" {
					hpics = append(hpics, a.RedpackCherryHPic01)
					hpics = append(hpics, a.RedpackCherryHPic02)
					hpics = append(hpics, a.RedpackCherryHPic03)
					hpics = append(hpics, a.RedpackCherryHPic04)
					hpics = append(hpics, a.RedpackCherryHPic05)
					hpics = append(hpics, a.RedpackCherryHPic06)
					hpics = append(hpics, a.RedpackCherryHPic07)
					gpics = append(gpics, a.RedpackCherryGPic01)
					gpics = append(gpics, a.RedpackCherryGPic02)
					hanis = append(hanis, a.RedpackCherryHAni01)
					hanis = append(hanis, a.RedpackCherryHAni02)
					ganis = append(ganis, a.RedpackCherryGAni01)
					ganis = append(ganis, a.RedpackCherryGAni02)
				} else if a.Topic == "03_christmas" {
					hpics = append(hpics, a.RedpackChristmasHPic01)
					hpics = append(hpics, a.RedpackChristmasHPic02)
					hpics = append(hpics, a.RedpackChristmasHPic03)
					hpics = append(hpics, a.RedpackChristmasHPic04)
					hpics = append(hpics, a.RedpackChristmasHPic05)
					hpics = append(hpics, a.RedpackChristmasHPic06)
					hpics = append(hpics, a.RedpackChristmasHPic07)
					hpics = append(hpics, a.RedpackChristmasHPic08)
					hpics = append(hpics, a.RedpackChristmasHPic09)
					hpics = append(hpics, a.RedpackChristmasHPic10)
					hpics = append(hpics, a.RedpackChristmasHPic11)
					hpics = append(hpics, a.RedpackChristmasHPic12)
					hpics = append(hpics, a.RedpackChristmasHPic13)
					gpics = append(gpics, a.RedpackChristmasGPic01)
					gpics = append(gpics, a.RedpackChristmasGPic02)
					gpics = append(gpics, a.RedpackChristmasGPic03)
					gpics = append(gpics, a.RedpackChristmasGPic04)
					cpics = append(cpics, a.RedpackChristmasCPic01)
					cpics = append(cpics, a.RedpackChristmasCPic02)
					canis = append(canis, a.RedpackChristmasCAni01)
				}
			}
			// 音樂
			musics = append(musics, a.RedpackBgmStart)
			musics = append(musics, a.RedpackBgmGaming)
			musics = append(musics, a.RedpackBgmEnd)
		} else if a.Game == "ropepack" {
			// 套紅包自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.RopepackClassicHPic01)
					hpics = append(hpics, a.RopepackClassicHPic02)
					hpics = append(hpics, a.RopepackClassicHPic03)
					hpics = append(hpics, a.RopepackClassicHPic04)
					hpics = append(hpics, a.RopepackClassicHPic05)
					hpics = append(hpics, a.RopepackClassicHPic06)
					hpics = append(hpics, a.RopepackClassicHPic07)
					hpics = append(hpics, a.RopepackClassicHPic08)
					hpics = append(hpics, a.RopepackClassicHPic09)
					hpics = append(hpics, a.RopepackClassicHPic10)
					gpics = append(gpics, a.RopepackClassicGPic01)
					gpics = append(gpics, a.RopepackClassicGPic02)
					gpics = append(gpics, a.RopepackClassicGPic03)
					gpics = append(gpics, a.RopepackClassicGPic04)
					gpics = append(gpics, a.RopepackClassicGPic05)
					gpics = append(gpics, a.RopepackClassicGPic06)
					hanis = append(hanis, a.RopepackClassicHAni01)
					ganis = append(ganis, a.RopepackClassicGAni01)
					ganis = append(ganis, a.RopepackClassicGAni02)
					canis = append(canis, a.RopepackClassicCAni01)
				} else if a.Topic == "02_newyear_rabbit" {
					hpics = append(hpics, a.RopepackNewyearRabbitHPic01)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic02)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic03)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic04)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic05)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic06)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic07)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic08)
					hpics = append(hpics, a.RopepackNewyearRabbitHPic09)
					gpics = append(gpics, a.RopepackNewyearRabbitGPic01)
					gpics = append(gpics, a.RopepackNewyearRabbitGPic02)
					gpics = append(gpics, a.RopepackNewyearRabbitGPic03)
					hanis = append(hanis, a.RopepackNewyearRabbitHAni01)
					ganis = append(ganis, a.RopepackNewyearRabbitGAni01)
					ganis = append(ganis, a.RopepackNewyearRabbitGAni02)
					ganis = append(ganis, a.RopepackNewyearRabbitGAni03)
					canis = append(canis, a.RopepackNewyearRabbitCAni01)
					canis = append(canis, a.RopepackNewyearRabbitCAni02)
				} else if a.Topic == "03_moonfestival" {
					hpics = append(hpics, a.RopepackMoonfestivalHPic01)
					hpics = append(hpics, a.RopepackMoonfestivalHPic02)
					hpics = append(hpics, a.RopepackMoonfestivalHPic03)
					hpics = append(hpics, a.RopepackMoonfestivalHPic04)
					hpics = append(hpics, a.RopepackMoonfestivalHPic05)
					hpics = append(hpics, a.RopepackMoonfestivalHPic06)
					hpics = append(hpics, a.RopepackMoonfestivalHPic07)
					hpics = append(hpics, a.RopepackMoonfestivalHPic08)
					hpics = append(hpics, a.RopepackMoonfestivalHPic09)
					gpics = append(gpics, a.RopepackMoonfestivalGPic01)
					gpics = append(gpics, a.RopepackMoonfestivalGPic02)
					cpics = append(cpics, a.RopepackMoonfestivalCPic01)
					hanis = append(hanis, a.RopepackMoonfestivalHAni01)
					ganis = append(ganis, a.RopepackMoonfestivalGAni01)
					ganis = append(ganis, a.RopepackMoonfestivalGAni02)
					canis = append(canis, a.RopepackMoonfestivalCAni01)
					canis = append(canis, a.RopepackMoonfestivalCAni02)
				} else if a.Topic == "04_3D" {
					hpics = append(hpics, a.Ropepack3DHPic01)
					hpics = append(hpics, a.Ropepack3DHPic02)
					hpics = append(hpics, a.Ropepack3DHPic03)
					hpics = append(hpics, a.Ropepack3DHPic04)
					hpics = append(hpics, a.Ropepack3DHPic05)
					hpics = append(hpics, a.Ropepack3DHPic06)
					hpics = append(hpics, a.Ropepack3DHPic07)
					hpics = append(hpics, a.Ropepack3DHPic08)
					hpics = append(hpics, a.Ropepack3DHPic09)
					hpics = append(hpics, a.Ropepack3DHPic10)
					hpics = append(hpics, a.Ropepack3DHPic11)
					hpics = append(hpics, a.Ropepack3DHPic12)
					hpics = append(hpics, a.Ropepack3DHPic13)
					hpics = append(hpics, a.Ropepack3DHPic14)
					hpics = append(hpics, a.Ropepack3DHPic15)
					gpics = append(gpics, a.Ropepack3DGPic01)
					gpics = append(gpics, a.Ropepack3DGPic02)
					gpics = append(gpics, a.Ropepack3DGPic03)
					gpics = append(gpics, a.Ropepack3DGPic04)
					hanis = append(hanis, a.Ropepack3DHAni01)
					hanis = append(hanis, a.Ropepack3DHAni02)
					hanis = append(hanis, a.Ropepack3DHAni03)
					ganis = append(ganis, a.Ropepack3DGAni01)
					ganis = append(ganis, a.Ropepack3DGAni02)
					canis = append(canis, a.Ropepack3DCAni01)
				}
			}
			// 音樂
			musics = append(musics, a.RopepackBgmStart)
			musics = append(musics, a.RopepackBgmGaming)
			musics = append(musics, a.RopepackBgmEnd)
		} else if a.Game == "lottery" {
			// 遊戲抽獎自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					if a.GameType == "jiugongge" {
						hpics = append(hpics, a.LotteryJiugonggeClassicHPic01)
						hpics = append(hpics, a.LotteryJiugonggeClassicHPic02)
						hpics = append(hpics, a.LotteryJiugonggeClassicHPic03)
						hpics = append(hpics, a.LotteryJiugonggeClassicHPic04)
						gpics = append(gpics, a.LotteryJiugonggeClassicGPic01)
						gpics = append(gpics, a.LotteryJiugonggeClassicGPic02)
						cpics = append(cpics, a.LotteryJiugonggeClassicCPic01)
						cpics = append(cpics, a.LotteryJiugonggeClassicCPic02)
						cpics = append(cpics, a.LotteryJiugonggeClassicCPic03)
						cpics = append(cpics, a.LotteryJiugonggeClassicCPic04)
						canis = append(canis, a.LotteryJiugonggeClassicCAni01)
						canis = append(canis, a.LotteryJiugonggeClassicCAni02)
						canis = append(canis, a.LotteryJiugonggeClassicCAni03)
					} else if a.GameType == "turntable" {
						hpics = append(hpics, a.LotteryTurntableClassicHPic01)
						hpics = append(hpics, a.LotteryTurntableClassicHPic02)
						hpics = append(hpics, a.LotteryTurntableClassicHPic03)
						hpics = append(hpics, a.LotteryTurntableClassicHPic04)
						gpics = append(gpics, a.LotteryTurntableClassicGPic01)
						gpics = append(gpics, a.LotteryTurntableClassicGPic02)
						cpics = append(cpics, a.LotteryTurntableClassicCPic01)
						cpics = append(cpics, a.LotteryTurntableClassicCPic02)
						cpics = append(cpics, a.LotteryTurntableClassicCPic03)
						cpics = append(cpics, a.LotteryTurntableClassicCPic04)
						cpics = append(cpics, a.LotteryTurntableClassicCPic05)
						cpics = append(cpics, a.LotteryTurntableClassicCPic06)
						canis = append(canis, a.LotteryTurntableClassicCAni01)
						canis = append(canis, a.LotteryTurntableClassicCAni02)
						canis = append(canis, a.LotteryTurntableClassicCAni03)
					}

				} else if a.Topic == "02_starrysky" {
					if a.GameType == "jiugongge" {
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic01)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic02)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic03)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic04)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic05)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic06)
						hpics = append(hpics, a.LotteryJiugonggeStarryskyHPic07)
						gpics = append(gpics, a.LotteryJiugonggeStarryskyGPic01)
						gpics = append(gpics, a.LotteryJiugonggeStarryskyGPic02)
						gpics = append(gpics, a.LotteryJiugonggeStarryskyGPic03)
						gpics = append(gpics, a.LotteryJiugonggeStarryskyGPic04)
						cpics = append(cpics, a.LotteryJiugonggeStarryskyCPic01)
						cpics = append(cpics, a.LotteryJiugonggeStarryskyCPic02)
						cpics = append(cpics, a.LotteryJiugonggeStarryskyCPic03)
						cpics = append(cpics, a.LotteryJiugonggeStarryskyCPic04)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni01)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni02)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni03)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni04)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni05)
						canis = append(canis, a.LotteryJiugonggeStarryskyCAni06)
					} else if a.GameType == "turntable" {
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic01)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic02)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic03)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic04)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic05)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic06)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic07)
						hpics = append(hpics, a.LotteryTurntableStarryskyHPic08)
						gpics = append(gpics, a.LotteryTurntableStarryskyGPic01)
						gpics = append(gpics, a.LotteryTurntableStarryskyGPic02)
						gpics = append(gpics, a.LotteryTurntableStarryskyGPic03)
						gpics = append(gpics, a.LotteryTurntableStarryskyGPic04)
						gpics = append(gpics, a.LotteryTurntableStarryskyGPic05)
						cpics = append(cpics, a.LotteryTurntableStarryskyCPic01)
						cpics = append(cpics, a.LotteryTurntableStarryskyCPic02)
						cpics = append(cpics, a.LotteryTurntableStarryskyCPic03)
						cpics = append(cpics, a.LotteryTurntableStarryskyCPic04)
						cpics = append(cpics, a.LotteryTurntableStarryskyCPic05)
						canis = append(canis, a.LotteryTurntableStarryskyCAni01)
						canis = append(canis, a.LotteryTurntableStarryskyCAni02)
						canis = append(canis, a.LotteryTurntableStarryskyCAni03)
						canis = append(canis, a.LotteryTurntableStarryskyCAni04)
						canis = append(canis, a.LotteryTurntableStarryskyCAni05)
						canis = append(canis, a.LotteryTurntableStarryskyCAni06)
						canis = append(canis, a.LotteryTurntableStarryskyCAni07)
					}
				}
			}
			// 音樂
			musics = append(musics, a.LotteryBgmGaming)
		} else if a.Game == "tugofwar" {
			// 拔河遊戲自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.TugofwarClassicHPic01)
					hpics = append(hpics, a.TugofwarClassicHPic02)
					hpics = append(hpics, a.TugofwarClassicHPic03)
					hpics = append(hpics, a.TugofwarClassicHPic04)
					hpics = append(hpics, a.TugofwarClassicHPic05)
					hpics = append(hpics, a.TugofwarClassicHPic06)
					hpics = append(hpics, a.TugofwarClassicHPic07)
					hpics = append(hpics, a.TugofwarClassicHPic08)
					hpics = append(hpics, a.TugofwarClassicHPic09)
					hpics = append(hpics, a.TugofwarClassicHPic10)
					hpics = append(hpics, a.TugofwarClassicHPic11)
					hpics = append(hpics, a.TugofwarClassicHPic12)
					hpics = append(hpics, a.TugofwarClassicHPic13)
					hpics = append(hpics, a.TugofwarClassicHPic14)
					hpics = append(hpics, a.TugofwarClassicHPic15)
					hpics = append(hpics, a.TugofwarClassicHPic16)
					hpics = append(hpics, a.TugofwarClassicHPic17)
					hpics = append(hpics, a.TugofwarClassicHPic18)
					hpics = append(hpics, a.TugofwarClassicHPic19)
					gpics = append(gpics, a.TugofwarClassicGPic01)
					gpics = append(gpics, a.TugofwarClassicGPic02)
					gpics = append(gpics, a.TugofwarClassicGPic03)
					gpics = append(gpics, a.TugofwarClassicGPic04)
					gpics = append(gpics, a.TugofwarClassicGPic05)
					gpics = append(gpics, a.TugofwarClassicGPic06)
					gpics = append(gpics, a.TugofwarClassicGPic07)
					gpics = append(gpics, a.TugofwarClassicGPic08)
					gpics = append(gpics, a.TugofwarClassicGPic09)
					hanis = append(hanis, a.TugofwarClassicHAni01)
					hanis = append(hanis, a.TugofwarClassicHAni02)
					hanis = append(hanis, a.TugofwarClassicHAni03)
					canis = append(canis, a.TugofwarClassicCAni01)
				} else if a.Topic == "02_school" {
					hpics = append(hpics, a.TugofwarSchoolHPic01)
					hpics = append(hpics, a.TugofwarSchoolHPic02)
					hpics = append(hpics, a.TugofwarSchoolHPic03)
					hpics = append(hpics, a.TugofwarSchoolHPic04)
					hpics = append(hpics, a.TugofwarSchoolHPic05)
					hpics = append(hpics, a.TugofwarSchoolHPic06)
					hpics = append(hpics, a.TugofwarSchoolHPic07)
					hpics = append(hpics, a.TugofwarSchoolHPic08)
					hpics = append(hpics, a.TugofwarSchoolHPic09)
					hpics = append(hpics, a.TugofwarSchoolHPic10)
					hpics = append(hpics, a.TugofwarSchoolHPic11)
					hpics = append(hpics, a.TugofwarSchoolHPic12)
					hpics = append(hpics, a.TugofwarSchoolHPic13)
					hpics = append(hpics, a.TugofwarSchoolHPic14)
					hpics = append(hpics, a.TugofwarSchoolHPic15)
					hpics = append(hpics, a.TugofwarSchoolHPic16)
					hpics = append(hpics, a.TugofwarSchoolHPic17)
					hpics = append(hpics, a.TugofwarSchoolHPic18)
					hpics = append(hpics, a.TugofwarSchoolHPic19)
					hpics = append(hpics, a.TugofwarSchoolHPic20)
					hpics = append(hpics, a.TugofwarSchoolHPic21)
					hpics = append(hpics, a.TugofwarSchoolHPic22)
					hpics = append(hpics, a.TugofwarSchoolHPic23)
					hpics = append(hpics, a.TugofwarSchoolHPic24)
					hpics = append(hpics, a.TugofwarSchoolHPic25)
					hpics = append(hpics, a.TugofwarSchoolHPic26)
					gpics = append(gpics, a.TugofwarSchoolGPic01)
					gpics = append(gpics, a.TugofwarSchoolGPic02)
					gpics = append(gpics, a.TugofwarSchoolGPic03)
					gpics = append(gpics, a.TugofwarSchoolGPic04)
					gpics = append(gpics, a.TugofwarSchoolGPic05)
					gpics = append(gpics, a.TugofwarSchoolGPic06)
					gpics = append(gpics, a.TugofwarSchoolGPic07)
					cpics = append(cpics, a.TugofwarSchoolCPic01)
					cpics = append(cpics, a.TugofwarSchoolCPic02)
					cpics = append(cpics, a.TugofwarSchoolCPic03)
					cpics = append(cpics, a.TugofwarSchoolCPic04)
					hanis = append(hanis, a.TugofwarSchoolHAni01)
					hanis = append(hanis, a.TugofwarSchoolHAni02)
					hanis = append(hanis, a.TugofwarSchoolHAni03)
					hanis = append(hanis, a.TugofwarSchoolHAni04)
					hanis = append(hanis, a.TugofwarSchoolHAni05)
					hanis = append(hanis, a.TugofwarSchoolHAni06)
					hanis = append(hanis, a.TugofwarSchoolHAni07)
				} else if a.Topic == "03_christmas" {
					hpics = append(hpics, a.TugofwarChristmasHPic01)
					hpics = append(hpics, a.TugofwarChristmasHPic02)
					hpics = append(hpics, a.TugofwarChristmasHPic03)
					hpics = append(hpics, a.TugofwarChristmasHPic04)
					hpics = append(hpics, a.TugofwarChristmasHPic05)
					hpics = append(hpics, a.TugofwarChristmasHPic06)
					hpics = append(hpics, a.TugofwarChristmasHPic07)
					hpics = append(hpics, a.TugofwarChristmasHPic08)
					hpics = append(hpics, a.TugofwarChristmasHPic09)
					hpics = append(hpics, a.TugofwarChristmasHPic10)
					hpics = append(hpics, a.TugofwarChristmasHPic11)
					hpics = append(hpics, a.TugofwarChristmasHPic12)
					hpics = append(hpics, a.TugofwarChristmasHPic13)
					hpics = append(hpics, a.TugofwarChristmasHPic14)
					hpics = append(hpics, a.TugofwarChristmasHPic15)
					hpics = append(hpics, a.TugofwarChristmasHPic16)
					hpics = append(hpics, a.TugofwarChristmasHPic17)
					hpics = append(hpics, a.TugofwarChristmasHPic18)
					hpics = append(hpics, a.TugofwarChristmasHPic19)
					hpics = append(hpics, a.TugofwarChristmasHPic20)
					hpics = append(hpics, a.TugofwarChristmasHPic21)
					gpics = append(gpics, a.TugofwarChristmasGPic01)
					gpics = append(gpics, a.TugofwarChristmasGPic02)
					gpics = append(gpics, a.TugofwarChristmasGPic03)
					gpics = append(gpics, a.TugofwarChristmasGPic04)
					gpics = append(gpics, a.TugofwarChristmasGPic05)
					gpics = append(gpics, a.TugofwarChristmasGPic06)
					cpics = append(cpics, a.TugofwarChristmasCPic01)
					cpics = append(cpics, a.TugofwarChristmasCPic02)
					cpics = append(cpics, a.TugofwarChristmasCPic03)
					cpics = append(cpics, a.TugofwarChristmasCPic04)
					hanis = append(hanis, a.TugofwarChristmasHAni01)
					hanis = append(hanis, a.TugofwarChristmasHAni02)
					hanis = append(hanis, a.TugofwarChristmasHAni03)
					canis = append(canis, a.TugofwarChristmasCAni01)
					canis = append(canis, a.TugofwarChristmasCAni02)
				}
			}

			// 音樂
			musics = append(musics, a.TugofwarBgmStart)
			musics = append(musics, a.TugofwarBgmGaming)
			musics = append(musics, a.TugofwarBgmEnd)
		} else if a.Game == "bingo" {
			// 賓果遊戲自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.BingoClassicHPic01)
					hpics = append(hpics, a.BingoClassicHPic02)
					hpics = append(hpics, a.BingoClassicHPic03)
					hpics = append(hpics, a.BingoClassicHPic04)
					hpics = append(hpics, a.BingoClassicHPic05)
					hpics = append(hpics, a.BingoClassicHPic06)
					hpics = append(hpics, a.BingoClassicHPic07)
					hpics = append(hpics, a.BingoClassicHPic08)
					hpics = append(hpics, a.BingoClassicHPic09)
					hpics = append(hpics, a.BingoClassicHPic10)
					hpics = append(hpics, a.BingoClassicHPic11)
					hpics = append(hpics, a.BingoClassicHPic12)
					hpics = append(hpics, a.BingoClassicHPic13)
					hpics = append(hpics, a.BingoClassicHPic14)
					hpics = append(hpics, a.BingoClassicHPic15)
					hpics = append(hpics, a.BingoClassicHPic16)
					gpics = append(gpics, a.BingoClassicGPic01)
					gpics = append(gpics, a.BingoClassicGPic02)
					gpics = append(gpics, a.BingoClassicGPic03)
					gpics = append(gpics, a.BingoClassicGPic04)
					gpics = append(gpics, a.BingoClassicGPic05)
					gpics = append(gpics, a.BingoClassicGPic06)
					gpics = append(gpics, a.BingoClassicGPic07)
					gpics = append(gpics, a.BingoClassicGPic08)
					cpics = append(cpics, a.BingoClassicCPic01)
					cpics = append(cpics, a.BingoClassicCPic02)
					cpics = append(cpics, a.BingoClassicCPic03)
					cpics = append(cpics, a.BingoClassicCPic04)
					// cpics = append(cpics, a.BingoClassicCPic05)
					hanis = append(hanis, a.BingoClassicHAni01)
					ganis = append(ganis, a.BingoClassicGAni01)
					canis = append(canis, a.BingoClassicCAni01)
					canis = append(canis, a.BingoClassicCAni02)
				} else if a.Topic == "02_newyear_dragon" {
					hpics = append(hpics, a.BingoNewyearDragonHPic01)
					hpics = append(hpics, a.BingoNewyearDragonHPic02)
					hpics = append(hpics, a.BingoNewyearDragonHPic03)
					hpics = append(hpics, a.BingoNewyearDragonHPic04)
					hpics = append(hpics, a.BingoNewyearDragonHPic05)
					hpics = append(hpics, a.BingoNewyearDragonHPic06)
					hpics = append(hpics, a.BingoNewyearDragonHPic07)
					hpics = append(hpics, a.BingoNewyearDragonHPic08)
					hpics = append(hpics, a.BingoNewyearDragonHPic09)
					hpics = append(hpics, a.BingoNewyearDragonHPic10)
					hpics = append(hpics, a.BingoNewyearDragonHPic11)
					hpics = append(hpics, a.BingoNewyearDragonHPic12)
					hpics = append(hpics, a.BingoNewyearDragonHPic13)
					hpics = append(hpics, a.BingoNewyearDragonHPic14)
					// hpics = append(hpics, a.BingoNewyearDragonHPic15)
					hpics = append(hpics, a.BingoNewyearDragonHPic16)
					hpics = append(hpics, a.BingoNewyearDragonHPic17)
					hpics = append(hpics, a.BingoNewyearDragonHPic18)
					hpics = append(hpics, a.BingoNewyearDragonHPic19)
					hpics = append(hpics, a.BingoNewyearDragonHPic20)
					hpics = append(hpics, a.BingoNewyearDragonHPic21)
					hpics = append(hpics, a.BingoNewyearDragonHPic22)
					gpics = append(gpics, a.BingoNewyearDragonGPic01)
					gpics = append(gpics, a.BingoNewyearDragonGPic02)
					gpics = append(gpics, a.BingoNewyearDragonGPic03)
					gpics = append(gpics, a.BingoNewyearDragonGPic04)
					gpics = append(gpics, a.BingoNewyearDragonGPic05)
					gpics = append(gpics, a.BingoNewyearDragonGPic06)
					gpics = append(gpics, a.BingoNewyearDragonGPic07)
					gpics = append(gpics, a.BingoNewyearDragonGPic08)
					cpics = append(cpics, a.BingoNewyearDragonCPic01)
					cpics = append(cpics, a.BingoNewyearDragonCPic02)
					cpics = append(cpics, a.BingoNewyearDragonCPic03)
					hanis = append(hanis, a.BingoNewyearDragonHAni01)
					ganis = append(ganis, a.BingoNewyearDragonGAni01)
					canis = append(canis, a.BingoNewyearDragonCAni01)
					canis = append(canis, a.BingoNewyearDragonCAni02)
					canis = append(canis, a.BingoNewyearDragonCAni03)
				} else if a.Topic == "03_cherry" {
					hpics = append(hpics, a.BingoCherryHPic01)
					hpics = append(hpics, a.BingoCherryHPic02)
					hpics = append(hpics, a.BingoCherryHPic03)
					hpics = append(hpics, a.BingoCherryHPic04)
					hpics = append(hpics, a.BingoCherryHPic05)
					hpics = append(hpics, a.BingoCherryHPic06)
					hpics = append(hpics, a.BingoCherryHPic07)
					hpics = append(hpics, a.BingoCherryHPic08)
					hpics = append(hpics, a.BingoCherryHPic09)
					hpics = append(hpics, a.BingoCherryHPic10)
					hpics = append(hpics, a.BingoCherryHPic11)
					hpics = append(hpics, a.BingoCherryHPic12)
					// hpics = append(hpics, a.BingoCherryHPic13)
					hpics = append(hpics, a.BingoCherryHPic14)
					hpics = append(hpics, a.BingoCherryHPic15)
					// hpics = append(hpics, a.BingoCherryHPic16)
					hpics = append(hpics, a.BingoCherryHPic17)
					hpics = append(hpics, a.BingoCherryHPic18)
					hpics = append(hpics, a.BingoCherryHPic19)
					gpics = append(gpics, a.BingoCherryGPic01)
					gpics = append(gpics, a.BingoCherryGPic02)
					gpics = append(gpics, a.BingoCherryGPic03)
					gpics = append(gpics, a.BingoCherryGPic04)
					gpics = append(gpics, a.BingoCherryGPic05)
					gpics = append(gpics, a.BingoCherryGPic06)
					cpics = append(cpics, a.BingoCherryCPic01)
					cpics = append(cpics, a.BingoCherryCPic02)
					cpics = append(cpics, a.BingoCherryCPic03)
					cpics = append(cpics, a.BingoCherryCPic04)
					// hanis = append(hanis, a.BingoCherryHAni01)
					hanis = append(hanis, a.BingoCherryHAni02)
					hanis = append(hanis, a.BingoCherryHAni03)
					ganis = append(ganis, a.BingoCherryGAni01)
					ganis = append(ganis, a.BingoCherryGAni02)
					canis = append(canis, a.BingoCherryCAni01)
					canis = append(canis, a.BingoCherryCAni02)
				}
			}

			// 音樂
			musics = append(musics, a.BingoBgmStart)
			musics = append(musics, a.BingoBgmGaming)
			musics = append(musics, a.BingoBgmEnd)
		} else if a.Game == "3DGachaMachine" {
			// 扭蛋機遊戲自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.GachaMachineClassicHPic02)
					hpics = append(hpics, a.GachaMachineClassicHPic03)
					hpics = append(hpics, a.GachaMachineClassicHPic04)
					hpics = append(hpics, a.GachaMachineClassicHPic05)
					gpics = append(gpics, a.GachaMachineClassicGPic01)
					gpics = append(gpics, a.GachaMachineClassicGPic02)
					gpics = append(gpics, a.GachaMachineClassicGPic03)
					gpics = append(gpics, a.GachaMachineClassicGPic04)
					gpics = append(gpics, a.GachaMachineClassicGPic05)
					cpics = append(cpics, a.GachaMachineClassicCPic01)

				}
			}

			// 音樂
			musics = append(musics, a.GachaMachineBgmGaming)
		} else if a.Game == "vote" {
			// 投票遊戲自定義
			if a.Skin == "customize" {
				if a.Topic == "01_classic" {
					hpics = append(hpics, a.VoteClassicHPic01)
					hpics = append(hpics, a.VoteClassicHPic02)
					hpics = append(hpics, a.VoteClassicHPic03)
					hpics = append(hpics, a.VoteClassicHPic04)
					hpics = append(hpics, a.VoteClassicHPic05)
					hpics = append(hpics, a.VoteClassicHPic06)
					hpics = append(hpics, a.VoteClassicHPic07)
					hpics = append(hpics, a.VoteClassicHPic08)
					hpics = append(hpics, a.VoteClassicHPic09)
					hpics = append(hpics, a.VoteClassicHPic10)
					hpics = append(hpics, a.VoteClassicHPic11)
					// hpics = append(hpics, a.VoteClassicHPic12)
					hpics = append(hpics, a.VoteClassicHPic13)
					hpics = append(hpics, a.VoteClassicHPic14)
					hpics = append(hpics, a.VoteClassicHPic15)
					hpics = append(hpics, a.VoteClassicHPic16)
					hpics = append(hpics, a.VoteClassicHPic17)
					hpics = append(hpics, a.VoteClassicHPic18)
					hpics = append(hpics, a.VoteClassicHPic19)
					hpics = append(hpics, a.VoteClassicHPic20)
					hpics = append(hpics, a.VoteClassicHPic21)
					// hpics = append(hpics, a.VoteClassicHPic22)
					hpics = append(hpics, a.VoteClassicHPic23)
					hpics = append(hpics, a.VoteClassicHPic24)
					hpics = append(hpics, a.VoteClassicHPic25)
					hpics = append(hpics, a.VoteClassicHPic26)
					hpics = append(hpics, a.VoteClassicHPic27)
					hpics = append(hpics, a.VoteClassicHPic28)
					hpics = append(hpics, a.VoteClassicHPic29)
					hpics = append(hpics, a.VoteClassicHPic30)
					hpics = append(hpics, a.VoteClassicHPic31)
					hpics = append(hpics, a.VoteClassicHPic32)
					hpics = append(hpics, a.VoteClassicHPic33)
					hpics = append(hpics, a.VoteClassicHPic34)
					hpics = append(hpics, a.VoteClassicHPic35)
					hpics = append(hpics, a.VoteClassicHPic36)
					hpics = append(hpics, a.VoteClassicHPic37)
					gpics = append(gpics, a.VoteClassicGPic01)
					gpics = append(gpics, a.VoteClassicGPic02)
					gpics = append(gpics, a.VoteClassicGPic03)
					gpics = append(gpics, a.VoteClassicGPic04)
					gpics = append(gpics, a.VoteClassicGPic05)
					gpics = append(gpics, a.VoteClassicGPic06)
					gpics = append(gpics, a.VoteClassicGPic07)
					cpics = append(cpics, a.VoteClassicCPic01)
					cpics = append(cpics, a.VoteClassicCPic02)
					cpics = append(cpics, a.VoteClassicCPic03)
					cpics = append(cpics, a.VoteClassicCPic04)
				}
			}

			// 音樂
			musics = append(musics, a.VoteBgmGaming)
		}
	} else if a.Skin == "customize_scene" {
		// 自定義場景時

		// 查詢自定義場景資料
		scene, _ := DefaultUserCustomizeSceneModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			Find(a.UserID, a.Topic)

		a.CustomizeSceneData = scene.CustomizeSceneData
	}

	// 判斷是否為自定義上傳資料
	// 主持靜態
	for i, picture := range hpics {
		if !strings.Contains(picture, "system") && a.Game != "vote" {
			hpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + picture
		} else if !strings.Contains(picture, "system") && a.Game == "vote" {
			hpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + picture
		}
	}
	// 玩家靜態
	for i, picture := range gpics {
		if !strings.Contains(picture, "system") && a.Game != "vote" {
			gpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + picture
		} else if !strings.Contains(picture, "system") && a.Game == "vote" {
			gpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + picture
		}
	}
	// 共用靜態
	for i, picture := range cpics {
		if !strings.Contains(picture, "system") && a.Game != "vote" {
			cpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + picture
		} else if !strings.Contains(picture, "system") && a.Game == "vote" {
			cpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + picture
		}
	}

	// 主持動態
	for i, anipicture := range hanis {
		if !strings.Contains(anipicture, "system") && a.Game != "vote" {
			hanis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + anipicture
		} else if !strings.Contains(anipicture, "system") && a.Game == "vote" {
			hanis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + anipicture
		}
	}
	// 玩家動態
	for i, anipicture := range ganis {
		if !strings.Contains(anipicture, "system") && a.Game != "vote" {
			ganis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + anipicture
		} else if !strings.Contains(anipicture, "system") && a.Game == "vote" {
			ganis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + anipicture
		}
	}
	// 共用動態
	for i, anipicture := range canis {
		if !strings.Contains(anipicture, "system") && a.Game != "vote" {
			canis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + anipicture
		} else if !strings.Contains(anipicture, "system") && a.Game == "vote" {
			canis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + anipicture
		}
	}

	// 音樂
	for i, music := range musics {
		if !strings.Contains(music, "system") && a.Game != "vote" {
			musics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + music
		} else if !strings.Contains(music, "system") && a.Game == "vote" {
			musics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + music
		}
	}

	// if a.Skin == "classic" {
	// 	// 經典，自定義圖片回傳空值
	// 	hpics = []string{}
	// 	gpics = []string{}
	// 	cpics = []string{}
	// 	hanis = []string{}
	// 	ganis = []string{}
	// 	canis = []string{}
	// }

	a.CustomizeHostPictures = hpics      // 主持靜態
	a.CustomizeGuestPictures = gpics     // 玩家靜態
	a.CustomizeCommonPictures = cpics    // 共用靜態
	a.CustomizeHostAnipictures = hanis   // 主持動態
	a.CustomizeGuestAnipictures = ganis  // 玩家動態
	a.CustomizeCommonAnipictures = canis // 共用動態
	a.CustomizeMusics = musics           // 音樂

	var questions = make([]QuestionModel, a.TotalQA)
	for i := 0; i < int(a.TotalQA); i++ {
		questions[i].Question, _ = m["qa_"+strconv.Itoa(i+1)].(string)
		// questions[i].Picture, _ = m["qa_"+strconv.Itoa(i+1)+"_picture"].(string)
		options, _ := m["qa_"+strconv.Itoa(i+1)+"_options"].(string)
		questions[i].Options = strings.Split(options, "&&&")
		questions[i].Answer, _ = m["qa_"+strconv.Itoa(i+1)+"_answer"].(string)
		questions[i].Score, _ = m["qa_"+strconv.Itoa(i+1)+"_score"].(int64)
	}
	a.Questions = questions

	// 快問快答
	qa1, _ := m["qa_1_options"].(string)
	a.QA1Options = strings.Split(qa1, "&&&")

	qa2, _ := m["qa_2_options"].(string)
	a.QA2Options = strings.Split(qa2, "&&&")

	qa3, _ := m["qa_3_options"].(string)
	a.QA3Options = strings.Split(qa3, "&&&")

	qa4, _ := m["qa_4_options"].(string)
	a.QA4Options = strings.Split(qa4, "&&&")

	qa5, _ := m["qa_5_options"].(string)
	a.QA5Options = strings.Split(qa5, "&&&")

	qa6, _ := m["qa_6_options"].(string)
	a.QA6Options = strings.Split(qa6, "&&&")

	qa7, _ := m["qa_7_options"].(string)
	a.QA7Options = strings.Split(qa7, "&&&")

	qa8, _ := m["qa_8_options"].(string)
	a.QA8Options = strings.Split(qa8, "&&&")

	qa9, _ := m["qa_9_options"].(string)
	a.QA9Options = strings.Split(qa9, "&&&")

	qa10, _ := m["qa_10_options"].(string)
	a.QA10Options = strings.Split(qa10, "&&&")

	qa11, _ := m["qa_11_options"].(string)
	a.QA11Options = strings.Split(qa11, "&&&")

	qa12, _ := m["qa_12_options"].(string)
	a.QA12Options = strings.Split(qa12, "&&&")

	qa13, _ := m["qa_13_options"].(string)
	a.QA13Options = strings.Split(qa13, "&&&")

	qa14, _ := m["qa_14_options"].(string)
	a.QA14Options = strings.Split(qa14, "&&&")

	qa15, _ := m["qa_15_options"].(string)
	a.QA15Options = strings.Split(qa15, "&&&")

	qa16, _ := m["qa_16_options"].(string)
	a.QA16Options = strings.Split(qa16, "&&&")

	qa17, _ := m["qa_17_options"].(string)
	a.QA17Options = strings.Split(qa17, "&&&")

	qa18, _ := m["qa_18_options"].(string)
	a.QA18Options = strings.Split(qa18, "&&&")

	qa19, _ := m["qa_19_options"].(string)
	a.QA19Options = strings.Split(qa19, "&&&")

	qa20, _ := m["qa_20_options"].(string)
	a.QA20Options = strings.Split(qa20, "&&&")

	return a
}

// MapToGameModelByMongo bson.M轉換[]GameModel
func MapToGameModelByMongo(items []bson.M) []GameModel {
	var games = make([]GameModel, 0)
	for _, item := range items {
		var (
			game GameModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &game)

		// game.LeftTeamPicture, _ = item["left_team_picture"].(string)
		if !strings.Contains(game.LeftTeamPicture, "system") {
			game.LeftTeamPicture = "/admin/uploads/" + game.UserID + "/" + game.ActivityID + "/interact/game/" + game.Game + "/" + game.GameID + "/" + game.LeftTeamPicture
		}

		// game.RightTeamPicture, _ = item["right_team_picture"].(string)
		if !strings.Contains(game.RightTeamPicture, "system") {
			game.RightTeamPicture = "/admin/uploads/" + game.UserID + "/" + game.ActivityID + "/interact/game/" + game.Game + "/" + game.GameID + "/" + game.RightTeamPicture
		}

		games = append(games, game)
	}
	return games
}

// MapToGameModelByMysql map轉換[]GameModel
func MapToGameModelByMysql(items []map[string]interface{}) []GameModel {
	var games = make([]GameModel, 0)
	for _, item := range items {
		var (
			game GameModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &game)

		// game.LeftTeamPicture, _ = item["left_team_picture"].(string)
		if !strings.Contains(game.LeftTeamPicture, "system") {
			game.LeftTeamPicture = "/admin/uploads/" + game.UserID + "/" + game.ActivityID + "/interact/game/" + game.Game + "/" + game.GameID + "/" + game.LeftTeamPicture
		}

		// game.RightTeamPicture, _ = item["right_team_picture"].(string)
		if !strings.Contains(game.RightTeamPicture, "system") {
			game.RightTeamPicture = "/admin/uploads/" + game.UserID + "/" + game.ActivityID + "/interact/game/" + game.Game + "/" + game.GameID + "/" + game.RightTeamPicture
		}

		games = append(games, game)
	}
	return games
}

// game.LeftTeamGameAttend, _ = item["left_team_game_attend"].(int64)
// game.RightTeamGameAttend, _ = item["right_team_game_attend"].(int64)
// game.Prize, _ = item["prize"].(string)

// game.MaxNumber, _ = item["max_number"].(int64)
// game.BingoLine, _ = item["bingo_line"].(int64)
// game.RoundPrize, _ = item["round_prize"].(int64)
// game.BingoRound, _ = item["bingo_round"].(int64)

// 扭蛋機遊戲
// game.GachaMachineReflection, _ = item["gacha_machine_reflection"].(float64)
// game.ReflectiveSwitch, _ = item["reflective_switch"].(string)

// 投票遊戲
// game.VoteScreen, _ = item["vote_screen"].(string)
// game.VoteTimes, _ = item["vote_times"].(int64)
// game.VoteMethod, _ = item["vote_method"].(string)
// game.VoteMethodPlayer, _ = item["vote_method_player"].(string)
// game.VoteRestriction, _ = item["vote_restriction"].(string)
// game.AvatarShape, _ = item["avatar_shape"].(string)
// game.VoteStartTime, _ = item["vote_start_time"].(string)
// game.VoteEndTime, _ = item["vote_end_time"].(string)
// game.AutoPlay, _ = item["auto_play"].(string)
// game.ShowRank, _ = item["show_rank"].(string)

// 編輯次數
// game.EditTimes, _ = item["edit_times"].(int64)

// 快問快答
// game.TotalQA, _ = item["total_qa"].(int64)
// game.QASecond, _ = item["qa_second"].(int64)
// game.QARound, _ = item["qa_round"].(int64)

// 用戶
// game.MaxActivityPeople, _ = item["max_activity_people"].(int64)
// game.MaxGamePeople, _ = item["max_game_people"].(int64)

// game.ID, _ = item["id"].(int64)
// game.UserID, _ = item["user_id"].(string)
// game.Device, _ = item["device"].(string)
// game.ActivityID, _ = item["activity_id"].(string)
// game.GameID, _ = item["game_id"].(string)
// game.Game, _ = item["game"].(string)
// game.Title, _ = item["title"].(string)
// game.GameType, _ = item["game_type"].(string)
// game.LimitTime, _ = item["limit_time"].(string)
// game.Second, _ = item["second"].(int64)
// game.MaxPeople, _ = item["max_people"].(int64)
// game.People, _ = item["people"].(int64)
// game.MaxTimes, _ = item["max_times"].(int64)
// game.Allow, _ = item["allow"].(string)
// game.Percent, _ = item["percent"].(int64)
// game.FirstPrize, _ = item["first_prize"].(int64)
// game.SecondPrize, _ = item["second_prize"].(int64)
// game.ThirdPrize, _ = item["third_prize"].(int64)
// game.GeneralPrize, _ = item["general_prize"].(int64)
// game.Topic, _ = item["topic"].(string)
// game.Skin, _ = item["skin"].(string)
// game.DisplayName, _ = item["display_name"].(string)
// game.GameRound, _ = item["game_round"].(int64)
// game.GameSecond, _ = item["game_second"].(int64)
// game.GameStatus, _ = item["game_status"].(string)
// game.GameAttend, _ = item["game_attend"].(int64)
// game.GameOrder, _ = item["game_order"].(int64)
// game.BoxReflection, _ = item["box_reflection"].(string)
// game.SamePeople, _ = item["same_people"].(string)

// game.AllowChooseTeam, _ = item["allow_choose_team"].(string)
// game.LeftTeamName, _ = item["left_team_name"].(string)
// game.RightTeamName, _ = item["right_team_name"].(string)

// a.TotalQA, _ = m["total_qa"].(int64)
// a.QASecond, _ = m["qa_second"].(int64)
// a.QARound, _ = m["qa_round"].(int64)

// a.QA1, _ = m["qa_1"].(string)
// a.QA1Answer, _ = m["qa_1_answer"].(string)
// a.QA1Score, _ = m["qa_1_score"].(int64)

// a.QA2, _ = m["qa_2"].(string)
// a.QA2Answer, _ = m["qa_2_answer"].(string)
// a.QA2Score, _ = m["qa_1_score"].(int64)

// a.QA3, _ = m["qa_3"].(string)
// a.QA3Answer, _ = m["qa_3_answer"].(string)
// a.QA3Score, _ = m["qa_1_score"].(int64)
// a.QA4, _ = m["qa_4"].(string)
// a.QA4Answer, _ = m["qa_4_answer"].(string)
// a.QA4Score, _ = m["qa_4_score"].(int64)

// a.QA5, _ = m["qa_5"].(string)
// a.QA5Answer, _ = m["qa_5_answer"].(string)
// a.QA5Score, _ = m["qa_5_score"].(int64)

// a.QA6, _ = m["qa_6"].(string)
// a.QA6Answer, _ = m["qa_6_answer"].(string)
// a.QA6Score, _ = m["qa_6_score"].(int64)

// a.QA7, _ = m["qa_7"].(string)
// a.QA7Answer, _ = m["qa_7_answer"].(string)
// a.QA7Score, _ = m["qa_7_score"].(int64)

// a.QA8, _ = m["qa_8"].(string)
// a.QA8Answer, _ = m["qa_8_answer"].(string)
// a.QA8Score, _ = m["qa_8_score"].(int64)

// a.QA9, _ = m["qa_9"].(string)

// a.QA9Answer, _ = m["qa_9_answer"].(string)
// a.QA9Score, _ = m["qa_9_score"].(int64)

// a.QA10, _ = m["qa_10"].(string)

// a.QA10Answer, _ = m["qa_10_answer"].(string)
// a.QA10Score, _ = m["qa_10_score"].(int64)

// a.QA11, _ = m["qa_11"].(string)

// a.QA11Answer, _ = m["qa_11_answer"].(string)
// a.QA11Score, _ = m["qa_11_score"].(int64)

// a.QA12, _ = m["qa_12"].(string)
// a.QA12Answer, _ = m["qa_12_answer"].(string)
// a.QA12Score, _ = m["qa_12_score"].(int64)

// a.QA13, _ = m["qa_13"].(string)
// a.QA13Answer, _ = m["qa_13_answer"].(string)
// a.QA13Score, _ = m["qa_13_score"].(int64)

// a.QA14, _ = m["qa_14"].(string)
// a.QA14Answer, _ = m["qa_14_answer"].(string)
// a.QA14Score, _ = m["qa_14_score"].(int64)

// a.QA15, _ = m["qa_15"].(string)
// a.QA15Answer, _ = m["qa_15_answer"].(string)
// a.QA15Score, _ = m["qa_15_score"].(int64)

// a.QA16, _ = m["qa_16"].(string)
// a.QA16Answer, _ = m["qa_16_answer"].(string)
// a.QA16Score, _ = m["qa_16_score"].(int64)

// a.QA17, _ = m["qa_17"].(string)
// a.QA17Answer, _ = m["qa_17_answer"].(string)
// a.QA17Score, _ = m["qa_17_score"].(int64)

// a.QA18, _ = m["qa_18"].(string)
// a.QA18Answer, _ = m["qa_18_answer"].(string)
// a.QA18Score, _ = m["qa_18_score"].(int64)

// a.QA19, _ = m["qa_19"].(string)
// a.QA19Answer, _ = m["qa_19_answer"].(string)
// a.QA19Score, _ = m["qa_19_score"].(int64)

// a.QA20, _ = m["qa_20"].(string)
// a.QA20Answer, _ = m["qa_20_answer"].(string)
// a.QA20Score, _ = m["qa_20_score"].(int64)

// a.ID, _ = m["id"].(int64)
// a.UserID, _ = m["user_id"].(string)
// a.ActivityID, _ = m["activity_id"].(string)
// a.GameID, _ = m["game_id"].(string)
// a.Game, _ = m["game"].(string)
// a.Title, _ = m["title"].(string)
// a.GameType, _ = m["game_type"].(string)
// a.LimitTime, _ = m["limit_time"].(string)
// a.Second, _ = m["second"].(int64)
// a.MaxPeople, _ = m["max_people"].(int64)
// a.People, _ = m["people"].(int64)
// a.MaxTimes, _ = m["max_times"].(int64)
// a.Allow, _ = m["allow"].(string)
// a.Percent, _ = m["percent"].(int64)
// a.FirstPrize, _ = m["first_prize"].(int64)
// a.SecondPrize, _ = m["second_prize"].(int64)
// a.ThirdPrize, _ = m["third_prize"].(int64)
// a.GeneralPrize, _ = m["general_prize"].(int64)
// a.Topic, _ = m["topic"].(string)
// a.Skin, _ = m["skin"].(string)
// a.Music, _ = m["music"].(string)
// a.DisplayName, _ = m["display_name"].(string)
// a.GameRound, _ = m["game_round"].(int64)
// a.GameSecond, _ = m["game_second"].(int64)
// a.GameStatus, _ = m["game_status"].(string)
// a.GameAttend, _ = m["game_attend"].(int64)
// a.GameOrder, _ = m["game_order"].(int64)
// a.BoxReflection, _ = m["box_reflection"].(string)
// a.SamePeople, _ = m["same_people"].(string)

// 活動資訊(join activity)
// a.UserID, _ = m["user_id"].(string)
// a.Device, _ = m["device"].(string)

// 拔河遊戲
// a.AllowChooseTeam, _ = m["allow_choose_team"].(string)
// a.LeftTeamName, _ = m["left_team_name"].(string)
// a.RightTeamName, _ = m["right_team_name"].(string)

// a.LeftTeamGameAttend, _ = m["left_team_game_attend"].(int64)
// a.RightTeamGameAttend, _ = m["right_team_game_attend"].(int64)
// a.Prize, _ = m["prize"].(string)

// 賓果遊戲
// a.MaxNumber, _ = m["max_number"].(int64)
// a.BingoLine, _ = m["bingo_line"].(int64)
// a.RoundPrize, _ = m["round_prize"].(int64)
// a.BingoRound, _ = m["bingo_round"].(int64)

// 扭蛋機遊戲
// a.GachaMachineReflection, _ = m["gacha_machine_reflection"].(float64)
// a.ReflectiveSwitch, _ = m["reflective_switch"].(string)

// 投票遊戲
// a.VoteScreen, _ = m["vote_screen"].(string)
// a.VoteTimes, _ = m["vote_times"].(int64)
// a.VoteMethod, _ = m["vote_method"].(string)
// a.VoteMethodPlayer, _ = m["vote_method_player"].(string)
// a.VoteRestriction, _ = m["vote_restriction"].(string)
// a.AvatarShape, _ = m["avatar_shape"].(string)
// a.VoteStartTime, _ = m["vote_start_time"].(string)
// a.VoteEndTime, _ = m["vote_end_time"].(string)
// a.AutoPlay, _ = m["auto_play"].(string)
// a.ShowRank, _ = m["show_rank"].(string)
// a.TitleSwitch, _ = m["title_switch"].(string)
// a.ArrangementGuest, _ = m["arrangement_guest"].(string)

// 編輯次數
// a.EditTimes, _ = m["edit_times"].(int64)

// 用戶
// a.MaxActivityPeople, _ = m["max_activity_people"].(int64)
// a.MaxGamePeople, _ = m["max_game_people"].(int64)

// 遊戲基本設置
// a.LotteryGameAllow, _ = m["lottery_game_allow"].(string)
// a.RedpackGameAllow, _ = m["redpack_game_allow"].(string)
// a.RopepackGameAllow, _ = m["ropepack_game_allow"].(string)
// a.WhackMoleGameAllow, _ = m["whack_mole_game_allow"].(string)
// a.MonopolyGameAllow, _ = m["monopoly_game_allow"].(string)
// a.QAGameAllow, _ = m["qa_game_allow"].(string)
// a.DrawNumbersGameAllow, _ = m["draw_numbers_game_allow"].(string)
// a.TugofwarGameAllow, _ = m["tugofwar_game_allow"].(string)
// a.BingoGameAllow, _ = m["bingo_game_allow"].(string)
// a.GachaMachineGameAllow, _ = m["3d_gacha_machine_game_allow"].(string)
// a.VoteGameAllow, _ = m["vote_game_allow"].(string)
// a.AllGameAllow, _ = m["all_game_allow"].(string)

// if a.Game == "ropepack" {
// 套紅包自定義
// a.RopepackClassicHPic01, _ = m["ropepack_classic_h_pic_01"].(string)
// a.RopepackClassicHPic02, _ = m["ropepack_classic_h_pic_02"].(string)
// a.RopepackClassicHPic03, _ = m["ropepack_classic_h_pic_03"].(string)
// a.RopepackClassicHPic04, _ = m["ropepack_classic_h_pic_04"].(string)
// a.RopepackClassicHPic05, _ = m["ropepack_classic_h_pic_05"].(string)
// a.RopepackClassicHPic06, _ = m["ropepack_classic_h_pic_06"].(string)
// a.RopepackClassicHPic07, _ = m["ropepack_classic_h_pic_07"].(string)
// a.RopepackClassicHPic08, _ = m["ropepack_classic_h_pic_08"].(string)
// a.RopepackClassicHPic09, _ = m["ropepack_classic_h_pic_09"].(string)
// a.RopepackClassicHPic10, _ = m["ropepack_classic_h_pic_10"].(string)
// a.RopepackClassicGPic01, _ = m["ropepack_classic_g_pic_01"].(string)
// a.RopepackClassicGPic02, _ = m["ropepack_classic_g_pic_02"].(string)
// a.RopepackClassicGPic03, _ = m["ropepack_classic_g_pic_03"].(string)
// a.RopepackClassicGPic04, _ = m["ropepack_classic_g_pic_04"].(string)
// a.RopepackClassicGPic05, _ = m["ropepack_classic_g_pic_05"].(string)
// a.RopepackClassicGPic06, _ = m["ropepack_classic_g_pic_06"].(string)
// a.RopepackClassicHAni01, _ = m["ropepack_classic_h_ani_01"].(string)
// a.RopepackClassicGAni01, _ = m["ropepack_classic_g_ani_01"].(string)
// a.RopepackClassicGAni02, _ = m["ropepack_classic_g_ani_02"].(string)
// a.RopepackClassicCAni01, _ = m["ropepack_classic_c_ani_01"].(string)

// a.RopepackNewyearRabbitHPic01, _ = m["ropepack_newyear_rabbit_h_pic_01"].(string)
// a.RopepackNewyearRabbitHPic02, _ = m["ropepack_newyear_rabbit_h_pic_02"].(string)
// a.RopepackNewyearRabbitHPic03, _ = m["ropepack_newyear_rabbit_h_pic_03"].(string)
// a.RopepackNewyearRabbitHPic04, _ = m["ropepack_newyear_rabbit_h_pic_04"].(string)
// a.RopepackNewyearRabbitHPic05, _ = m["ropepack_newyear_rabbit_h_pic_05"].(string)
// a.RopepackNewyearRabbitHPic06, _ = m["ropepack_newyear_rabbit_h_pic_06"].(string)
// a.RopepackNewyearRabbitHPic07, _ = m["ropepack_newyear_rabbit_h_pic_07"].(string)
// a.RopepackNewyearRabbitHPic08, _ = m["ropepack_newyear_rabbit_h_pic_08"].(string)
// a.RopepackNewyearRabbitHPic09, _ = m["ropepack_newyear_rabbit_h_pic_09"].(string)
// a.RopepackNewyearRabbitGPic01, _ = m["ropepack_newyear_rabbit_g_pic_01"].(string)
// a.RopepackNewyearRabbitGPic02, _ = m["ropepack_newyear_rabbit_g_pic_02"].(string)
// a.RopepackNewyearRabbitGPic03, _ = m["ropepack_newyear_rabbit_g_pic_03"].(string)
// a.RopepackNewyearRabbitHAni01, _ = m["ropepack_newyear_rabbit_h_ani_01"].(string)
// a.RopepackNewyearRabbitGAni01, _ = m["ropepack_newyear_rabbit_g_ani_01"].(string)
// a.RopepackNewyearRabbitGAni02, _ = m["ropepack_newyear_rabbit_g_ani_02"].(string)
// a.RopepackNewyearRabbitGAni03, _ = m["ropepack_newyear_rabbit_g_ani_03"].(string)
// a.RopepackNewyearRabbitCAni01, _ = m["ropepack_newyear_rabbit_c_ani_01"].(string)
// a.RopepackNewyearRabbitCAni02, _ = m["ropepack_newyear_rabbit_c_ani_02"].(string)

// a.RopepackMoonfestivalHPic01, _ = m["ropepack_moonfestival_h_pic_01"].(string)
// a.RopepackMoonfestivalHPic02, _ = m["ropepack_moonfestival_h_pic_02"].(string)
// a.RopepackMoonfestivalHPic03, _ = m["ropepack_moonfestival_h_pic_03"].(string)
// a.RopepackMoonfestivalHPic04, _ = m["ropepack_moonfestival_h_pic_04"].(string)
// a.RopepackMoonfestivalHPic05, _ = m["ropepack_moonfestival_h_pic_05"].(string)
// a.RopepackMoonfestivalHPic06, _ = m["ropepack_moonfestival_h_pic_06"].(string)
// a.RopepackMoonfestivalHPic07, _ = m["ropepack_moonfestival_h_pic_07"].(string)
// a.RopepackMoonfestivalHPic08, _ = m["ropepack_moonfestival_h_pic_08"].(string)
// a.RopepackMoonfestivalHPic09, _ = m["ropepack_moonfestival_h_pic_09"].(string)
// a.RopepackMoonfestivalGPic01, _ = m["ropepack_moonfestival_g_pic_01"].(string)
// a.RopepackMoonfestivalGPic02, _ = m["ropepack_moonfestival_g_pic_02"].(string)
// a.RopepackMoonfestivalCPic01, _ = m["ropepack_moonfestival_c_pic_01"].(string)
// a.RopepackMoonfestivalHAni01, _ = m["ropepack_moonfestival_h_ani_01"].(string)
// a.RopepackMoonfestivalGAni01, _ = m["ropepack_moonfestival_g_ani_01"].(string)
// a.RopepackMoonfestivalGAni02, _ = m["ropepack_moonfestival_g_ani_02"].(string)
// a.RopepackMoonfestivalCAni01, _ = m["ropepack_moonfestival_c_ani_01"].(string)
// a.RopepackMoonfestivalCAni02, _ = m["ropepack_moonfestival_c_ani_02"].(string)

// a.Ropepack3DHPic01, _ = m["ropepack_3D_h_pic_01"].(string)
// a.Ropepack3DHPic02, _ = m["ropepack_3D_h_pic_02"].(string)
// a.Ropepack3DHPic03, _ = m["ropepack_3D_h_pic_03"].(string)
// a.Ropepack3DHPic04, _ = m["ropepack_3D_h_pic_04"].(string)
// a.Ropepack3DHPic05, _ = m["ropepack_3D_h_pic_05"].(string)
// a.Ropepack3DHPic06, _ = m["ropepack_3D_h_pic_06"].(string)
// a.Ropepack3DHPic07, _ = m["ropepack_3D_h_pic_07"].(string)
// a.Ropepack3DHPic08, _ = m["ropepack_3D_h_pic_08"].(string)
// a.Ropepack3DHPic09, _ = m["ropepack_3D_h_pic_09"].(string)
// a.Ropepack3DHPic10, _ = m["ropepack_3D_h_pic_10"].(string)
// a.Ropepack3DHPic11, _ = m["ropepack_3D_h_pic_11"].(string)
// a.Ropepack3DHPic12, _ = m["ropepack_3D_h_pic_12"].(string)
// a.Ropepack3DHPic13, _ = m["ropepack_3D_h_pic_13"].(string)
// a.Ropepack3DHPic14, _ = m["ropepack_3D_h_pic_14"].(string)
// a.Ropepack3DHPic15, _ = m["ropepack_3D_h_pic_15"].(string)
// a.Ropepack3DGPic01, _ = m["ropepack_3D_g_pic_01"].(string)
// a.Ropepack3DGPic02, _ = m["ropepack_3D_g_pic_02"].(string)
// a.Ropepack3DGPic03, _ = m["ropepack_3D_g_pic_03"].(string)
// a.Ropepack3DGPic04, _ = m["ropepack_3D_g_pic_04"].(string)
// a.Ropepack3DHAni01, _ = m["ropepack_3D_h_ani_01"].(string)
// a.Ropepack3DHAni02, _ = m["ropepack_3D_h_ani_02"].(string)
// a.Ropepack3DHAni03, _ = m["ropepack_3D_h_ani_03"].(string)
// a.Ropepack3DGAni01, _ = m["ropepack_3D_g_ani_01"].(string)
// a.Ropepack3DGAni02, _ = m["ropepack_3D_g_ani_02"].(string)
// a.Ropepack3DCAni01, _ = m["ropepack_3D_c_ani_01"].(string)

// 音樂
// a.RopepackBgmStart, _ = m["ropepack_bgm_start"].(string)
// a.RopepackBgmGaming, _ = m["ropepack_bgm_gaming"].(string)
// a.RopepackBgmEnd, _ = m["ropepack_bgm_end"].(string)
// }

// if a.Game == "redpack" {
// 搖紅包自定義
// 	a.RedpackClassicHPic01, _ = m["redpack_classic_h_pic_01"].(string)
// 	a.RedpackClassicHPic02, _ = m["redpack_classic_h_pic_02"].(string)
// 	a.RedpackClassicHPic03, _ = m["redpack_classic_h_pic_03"].(string)
// 	a.RedpackClassicHPic04, _ = m["redpack_classic_h_pic_04"].(string)
// 	a.RedpackClassicHPic05, _ = m["redpack_classic_h_pic_05"].(string)
// 	a.RedpackClassicHPic06, _ = m["redpack_classic_h_pic_06"].(string)
// 	a.RedpackClassicHPic07, _ = m["redpack_classic_h_pic_07"].(string)
// 	a.RedpackClassicHPic08, _ = m["redpack_classic_h_pic_08"].(string)
// 	a.RedpackClassicHPic09, _ = m["redpack_classic_h_pic_09"].(string)
// 	a.RedpackClassicHPic10, _ = m["redpack_classic_h_pic_10"].(string)
// 	a.RedpackClassicHPic11, _ = m["redpack_classic_h_pic_11"].(string)
// 	a.RedpackClassicHPic12, _ = m["redpack_classic_h_pic_12"].(string)
// 	a.RedpackClassicHPic13, _ = m["redpack_classic_h_pic_13"].(string)
// 	a.RedpackClassicGPic01, _ = m["redpack_classic_g_pic_01"].(string)
// 	a.RedpackClassicGPic02, _ = m["redpack_classic_g_pic_02"].(string)
// 	a.RedpackClassicGPic03, _ = m["redpack_classic_g_pic_03"].(string)
// 	a.RedpackClassicHAni01, _ = m["redpack_classic_h_ani_01"].(string)
// 	a.RedpackClassicHAni02, _ = m["redpack_classic_h_ani_02"].(string)
// 	a.RedpackClassicGAni01, _ = m["redpack_classic_g_ani_01"].(string)
// 	a.RedpackClassicGAni02, _ = m["redpack_classic_g_ani_02"].(string)
// 	a.RedpackClassicGAni03, _ = m["redpack_classic_g_ani_03"].(string)

// 	a.RedpackCherryHPic01, _ = m["redpack_cherry_h_pic_01"].(string)
// 	a.RedpackCherryHPic02, _ = m["redpack_cherry_h_pic_02"].(string)
// 	a.RedpackCherryHPic03, _ = m["redpack_cherry_h_pic_03"].(string)
// 	a.RedpackCherryHPic04, _ = m["redpack_cherry_h_pic_04"].(string)
// 	a.RedpackCherryHPic05, _ = m["redpack_cherry_h_pic_05"].(string)
// 	a.RedpackCherryHPic06, _ = m["redpack_cherry_h_pic_06"].(string)
// 	a.RedpackCherryHPic07, _ = m["redpack_cherry_h_pic_07"].(string)
// 	a.RedpackCherryGPic01, _ = m["redpack_cherry_g_pic_01"].(string)
// 	a.RedpackCherryGPic02, _ = m["redpack_cherry_g_pic_02"].(string)
// 	a.RedpackCherryHAni01, _ = m["redpack_cherry_h_ani_01"].(string)
// 	a.RedpackCherryHAni02, _ = m["redpack_cherry_h_ani_02"].(string)
// 	a.RedpackCherryGAni01, _ = m["redpack_cherry_g_ani_01"].(string)
// 	a.RedpackCherryGAni02, _ = m["redpack_cherry_g_ani_02"].(string)

// 	a.RedpackChristmasHPic01, _ = m["redpack_christmas_h_pic_01"].(string)
// 	a.RedpackChristmasHPic02, _ = m["redpack_christmas_h_pic_02"].(string)
// 	a.RedpackChristmasHPic03, _ = m["redpack_christmas_h_pic_03"].(string)
// 	a.RedpackChristmasHPic04, _ = m["redpack_christmas_h_pic_04"].(string)
// 	a.RedpackChristmasHPic05, _ = m["redpack_christmas_h_pic_05"].(string)
// 	a.RedpackChristmasHPic06, _ = m["redpack_christmas_h_pic_06"].(string)
// 	a.RedpackChristmasHPic07, _ = m["redpack_christmas_h_pic_07"].(string)
// 	a.RedpackChristmasHPic08, _ = m["redpack_christmas_h_pic_08"].(string)
// 	a.RedpackChristmasHPic09, _ = m["redpack_christmas_h_pic_09"].(string)
// 	a.RedpackChristmasHPic10, _ = m["redpack_christmas_h_pic_10"].(string)
// 	a.RedpackChristmasHPic11, _ = m["redpack_christmas_h_pic_11"].(string)
// 	a.RedpackChristmasHPic12, _ = m["redpack_christmas_h_pic_12"].(string)
// 	a.RedpackChristmasHPic13, _ = m["redpack_christmas_h_pic_13"].(string)
// 	a.RedpackChristmasGPic01, _ = m["redpack_christmas_g_pic_01"].(string)
// 	a.RedpackChristmasGPic02, _ = m["redpack_christmas_g_pic_02"].(string)
// 	a.RedpackChristmasGPic03, _ = m["redpack_christmas_g_pic_03"].(string)
// 	a.RedpackChristmasGPic04, _ = m["redpack_christmas_g_pic_04"].(string)
// 	a.RedpackChristmasCPic01, _ = m["redpack_christmas_c_pic_01"].(string)
// 	a.RedpackChristmasCPic02, _ = m["redpack_christmas_c_pic_02"].(string)
// 	a.RedpackChristmasCAni01, _ = m["redpack_christmas_c_ani_01"].(string)

// 	// 音樂
// 	a.RedpackBgmStart, _ = m["redpack_bgm_start"].(string)
// 	a.RedpackBgmGaming, _ = m["redpack_bgm_gaming"].(string)
// 	a.RedpackBgmEnd, _ = m["redpack_bgm_end"].(string)
// }

// if a.Game == "lottery" {
// 	// 遊戲抽獎自定義
// 	a.LotteryJiugonggeClassicHPic01, _ = m["lottery_jiugongge_classic_h_pic_01"].(string)
// 	a.LotteryJiugonggeClassicHPic02, _ = m["lottery_jiugongge_classic_h_pic_02"].(string)
// 	a.LotteryJiugonggeClassicHPic03, _ = m["lottery_jiugongge_classic_h_pic_03"].(string)
// 	a.LotteryJiugonggeClassicHPic04, _ = m["lottery_jiugongge_classic_h_pic_04"].(string)
// 	a.LotteryJiugonggeClassicGPic01, _ = m["lottery_jiugongge_classic_g_pic_01"].(string)
// 	a.LotteryJiugonggeClassicGPic02, _ = m["lottery_jiugongge_classic_g_pic_02"].(string)
// 	a.LotteryJiugonggeClassicCPic01, _ = m["lottery_jiugongge_classic_c_pic_01"].(string)
// 	a.LotteryJiugonggeClassicCPic02, _ = m["lottery_jiugongge_classic_c_pic_02"].(string)
// 	a.LotteryJiugonggeClassicCPic03, _ = m["lottery_jiugongge_classic_c_pic_03"].(string)
// 	a.LotteryJiugonggeClassicCPic04, _ = m["lottery_jiugongge_classic_c_pic_04"].(string)
// 	a.LotteryJiugonggeClassicCAni01, _ = m["lottery_jiugongge_classic_c_ani_01"].(string)
// 	a.LotteryJiugonggeClassicCAni02, _ = m["lottery_jiugongge_classic_c_ani_02"].(string)
// 	a.LotteryJiugonggeClassicCAni03, _ = m["lottery_jiugongge_classic_c_ani_03"].(string)

// 	a.LotteryTurntableClassicHPic01, _ = m["lottery_turntable_classic_h_pic_01"].(string)
// 	a.LotteryTurntableClassicHPic02, _ = m["lottery_turntable_classic_h_pic_02"].(string)
// 	a.LotteryTurntableClassicHPic03, _ = m["lottery_turntable_classic_h_pic_03"].(string)
// 	a.LotteryTurntableClassicHPic04, _ = m["lottery_turntable_classic_h_pic_04"].(string)
// 	a.LotteryTurntableClassicGPic01, _ = m["lottery_turntable_classic_g_pic_01"].(string)
// 	a.LotteryTurntableClassicGPic02, _ = m["lottery_turntable_classic_g_pic_02"].(string)
// 	a.LotteryTurntableClassicCPic01, _ = m["lottery_turntable_classic_c_pic_01"].(string)
// 	a.LotteryTurntableClassicCPic02, _ = m["lottery_turntable_classic_c_pic_02"].(string)
// 	a.LotteryTurntableClassicCPic03, _ = m["lottery_turntable_classic_c_pic_03"].(string)
// 	a.LotteryTurntableClassicCPic04, _ = m["lottery_turntable_classic_c_pic_04"].(string)
// 	a.LotteryTurntableClassicCPic05, _ = m["lottery_turntable_classic_c_pic_05"].(string)
// 	a.LotteryTurntableClassicCPic06, _ = m["lottery_turntable_classic_c_pic_06"].(string)
// 	a.LotteryTurntableClassicCAni01, _ = m["lottery_turntable_classic_c_ani_01"].(string)
// 	a.LotteryTurntableClassicCAni02, _ = m["lottery_turntable_classic_c_ani_02"].(string)
// 	a.LotteryTurntableClassicCAni03, _ = m["lottery_turntable_classic_c_ani_03"].(string)

// 	a.LotteryJiugonggeStarryskyHPic01, _ = m["lottery_jiugongge_starrysky_h_pic_01"].(string)
// 	a.LotteryJiugonggeStarryskyHPic02, _ = m["lottery_jiugongge_starrysky_h_pic_02"].(string)
// 	a.LotteryJiugonggeStarryskyHPic03, _ = m["lottery_jiugongge_starrysky_h_pic_03"].(string)
// 	a.LotteryJiugonggeStarryskyHPic04, _ = m["lottery_jiugongge_starrysky_h_pic_04"].(string)
// 	a.LotteryJiugonggeStarryskyHPic05, _ = m["lottery_jiugongge_starrysky_h_pic_05"].(string)
// 	a.LotteryJiugonggeStarryskyHPic06, _ = m["lottery_jiugongge_starrysky_h_pic_06"].(string)
// 	a.LotteryJiugonggeStarryskyHPic07, _ = m["lottery_jiugongge_starrysky_h_pic_07"].(string)
// 	a.LotteryJiugonggeStarryskyGPic01, _ = m["lottery_jiugongge_starrysky_g_pic_01"].(string)
// 	a.LotteryJiugonggeStarryskyGPic02, _ = m["lottery_jiugongge_starrysky_g_pic_02"].(string)
// 	a.LotteryJiugonggeStarryskyGPic03, _ = m["lottery_jiugongge_starrysky_g_pic_03"].(string)
// 	a.LotteryJiugonggeStarryskyGPic04, _ = m["lottery_jiugongge_starrysky_g_pic_04"].(string)
// 	a.LotteryJiugonggeStarryskyCPic01, _ = m["lottery_jiugongge_starrysky_c_pic_01"].(string)
// 	a.LotteryJiugonggeStarryskyCPic02, _ = m["lottery_jiugongge_starrysky_c_pic_02"].(string)
// 	a.LotteryJiugonggeStarryskyCPic03, _ = m["lottery_jiugongge_starrysky_c_pic_03"].(string)
// 	a.LotteryJiugonggeStarryskyCPic04, _ = m["lottery_jiugongge_starrysky_c_pic_04"].(string)
// 	a.LotteryJiugonggeStarryskyCAni01, _ = m["lottery_jiugongge_starrysky_c_ani_01"].(string)
// 	a.LotteryJiugonggeStarryskyCAni02, _ = m["lottery_jiugongge_starrysky_c_ani_02"].(string)
// 	a.LotteryJiugonggeStarryskyCAni03, _ = m["lottery_jiugongge_starrysky_c_ani_03"].(string)
// 	a.LotteryJiugonggeStarryskyCAni04, _ = m["lottery_jiugongge_starrysky_c_ani_04"].(string)
// 	a.LotteryJiugonggeStarryskyCAni05, _ = m["lottery_jiugongge_starrysky_c_ani_05"].(string)
// 	a.LotteryJiugonggeStarryskyCAni06, _ = m["lottery_jiugongge_starrysky_c_ani_06"].(string)

// 	a.LotteryTurntableStarryskyHPic01, _ = m["lottery_turntable_starrysky_h_pic_01"].(string)
// 	a.LotteryTurntableStarryskyHPic02, _ = m["lottery_turntable_starrysky_h_pic_02"].(string)
// 	a.LotteryTurntableStarryskyHPic03, _ = m["lottery_turntable_starrysky_h_pic_03"].(string)
// 	a.LotteryTurntableStarryskyHPic04, _ = m["lottery_turntable_starrysky_h_pic_04"].(string)
// 	a.LotteryTurntableStarryskyHPic05, _ = m["lottery_turntable_starrysky_h_pic_05"].(string)
// 	a.LotteryTurntableStarryskyHPic06, _ = m["lottery_turntable_starrysky_h_pic_06"].(string)
// 	a.LotteryTurntableStarryskyHPic07, _ = m["lottery_turntable_starrysky_h_pic_07"].(string)
// 	a.LotteryTurntableStarryskyHPic08, _ = m["lottery_turntable_starrysky_h_pic_08"].(string)
// 	a.LotteryTurntableStarryskyGPic01, _ = m["lottery_turntable_starrysky_g_pic_01"].(string)
// 	a.LotteryTurntableStarryskyGPic02, _ = m["lottery_turntable_starrysky_g_pic_02"].(string)
// 	a.LotteryTurntableStarryskyGPic03, _ = m["lottery_turntable_starrysky_g_pic_03"].(string)
// 	a.LotteryTurntableStarryskyGPic04, _ = m["lottery_turntable_starrysky_g_pic_04"].(string)
// 	a.LotteryTurntableStarryskyGPic05, _ = m["lottery_turntable_starrysky_g_pic_05"].(string)
// 	a.LotteryTurntableStarryskyCPic01, _ = m["lottery_turntable_starrysky_c_pic_01"].(string)
// 	a.LotteryTurntableStarryskyCPic02, _ = m["lottery_turntable_starrysky_c_pic_02"].(string)
// 	a.LotteryTurntableStarryskyCPic03, _ = m["lottery_turntable_starrysky_c_pic_03"].(string)
// 	a.LotteryTurntableStarryskyCPic04, _ = m["lottery_turntable_starrysky_c_pic_04"].(string)
// 	a.LotteryTurntableStarryskyCPic05, _ = m["lottery_turntable_starrysky_c_pic_05"].(string)
// 	a.LotteryTurntableStarryskyCAni01, _ = m["lottery_turntable_starrysky_c_ani_01"].(string)
// 	a.LotteryTurntableStarryskyCAni02, _ = m["lottery_turntable_starrysky_c_ani_02"].(string)
// 	a.LotteryTurntableStarryskyCAni03, _ = m["lottery_turntable_starrysky_c_ani_03"].(string)
// 	a.LotteryTurntableStarryskyCAni04, _ = m["lottery_turntable_starrysky_c_ani_04"].(string)
// 	a.LotteryTurntableStarryskyCAni05, _ = m["lottery_turntable_starrysky_c_ani_05"].(string)
// 	a.LotteryTurntableStarryskyCAni06, _ = m["lottery_turntable_starrysky_c_ani_06"].(string)
// 	a.LotteryTurntableStarryskyCAni07, _ = m["lottery_turntable_starrysky_c_ani_07"].(string)

// 	// 音樂
// 	a.LotteryBgmGaming, _ = m["lottery_bgm_gaming"].(string)
// }

// if a.Game == "draw_numbers" {
// 	// 搖號抽獎自定義圖片
// 	a.DrawNumbersClassicHPic01, _ = m["draw_numbers_classic_h_pic_01"].(string)
// 	a.DrawNumbersClassicHPic02, _ = m["draw_numbers_classic_h_pic_02"].(string)
// 	a.DrawNumbersClassicHPic03, _ = m["draw_numbers_classic_h_pic_03"].(string)
// 	a.DrawNumbersClassicHPic04, _ = m["draw_numbers_classic_h_pic_04"].(string)
// 	a.DrawNumbersClassicHPic05, _ = m["draw_numbers_classic_h_pic_05"].(string)
// 	a.DrawNumbersClassicHPic06, _ = m["draw_numbers_classic_h_pic_06"].(string)
// 	a.DrawNumbersClassicHPic07, _ = m["draw_numbers_classic_h_pic_07"].(string)
// 	a.DrawNumbersClassicHPic08, _ = m["draw_numbers_classic_h_pic_08"].(string)
// 	a.DrawNumbersClassicHPic09, _ = m["draw_numbers_classic_h_pic_09"].(string)
// 	a.DrawNumbersClassicHPic10, _ = m["draw_numbers_classic_h_pic_10"].(string)
// 	a.DrawNumbersClassicHPic11, _ = m["draw_numbers_classic_h_pic_11"].(string)
// 	a.DrawNumbersClassicHPic12, _ = m["draw_numbers_classic_h_pic_12"].(string)
// 	a.DrawNumbersClassicHPic13, _ = m["draw_numbers_classic_h_pic_13"].(string)
// 	a.DrawNumbersClassicHPic14, _ = m["draw_numbers_classic_h_pic_14"].(string)
// 	a.DrawNumbersClassicHPic15, _ = m["draw_numbers_classic_h_pic_15"].(string)
// 	a.DrawNumbersClassicHPic16, _ = m["draw_numbers_classic_h_pic_16"].(string)
// 	a.DrawNumbersClassicHAni01, _ = m["draw_numbers_classic_h_ani_01"].(string)

// 	a.DrawNumbersGoldHPic01, _ = m["draw_numbers_gold_h_pic_01"].(string)
// 	a.DrawNumbersGoldHPic02, _ = m["draw_numbers_gold_h_pic_02"].(string)
// 	a.DrawNumbersGoldHPic03, _ = m["draw_numbers_gold_h_pic_03"].(string)
// 	a.DrawNumbersGoldHPic04, _ = m["draw_numbers_gold_h_pic_04"].(string)
// 	a.DrawNumbersGoldHPic05, _ = m["draw_numbers_gold_h_pic_05"].(string)
// 	a.DrawNumbersGoldHPic06, _ = m["draw_numbers_gold_h_pic_06"].(string)
// 	a.DrawNumbersGoldHPic07, _ = m["draw_numbers_gold_h_pic_07"].(string)
// 	a.DrawNumbersGoldHPic08, _ = m["draw_numbers_gold_h_pic_08"].(string)
// 	a.DrawNumbersGoldHPic09, _ = m["draw_numbers_gold_h_pic_09"].(string)
// 	a.DrawNumbersGoldHPic10, _ = m["draw_numbers_gold_h_pic_10"].(string)
// 	a.DrawNumbersGoldHPic11, _ = m["draw_numbers_gold_h_pic_11"].(string)
// 	a.DrawNumbersGoldHPic12, _ = m["draw_numbers_gold_h_pic_12"].(string)
// 	a.DrawNumbersGoldHPic13, _ = m["draw_numbers_gold_h_pic_13"].(string)
// 	a.DrawNumbersGoldHPic14, _ = m["draw_numbers_gold_h_pic_14"].(string)
// 	a.DrawNumbersGoldHAni01, _ = m["draw_numbers_gold_h_ani_01"].(string)
// 	a.DrawNumbersGoldHAni02, _ = m["draw_numbers_gold_h_ani_02"].(string)
// 	a.DrawNumbersGoldHAni03, _ = m["draw_numbers_gold_h_ani_03"].(string)

// 	a.DrawNumbersNewyearDragonHPic01, _ = m["draw_numbers_newyear_dragon_h_pic_01"].(string)
// 	a.DrawNumbersNewyearDragonHPic02, _ = m["draw_numbers_newyear_dragon_h_pic_02"].(string)
// 	a.DrawNumbersNewyearDragonHPic03, _ = m["draw_numbers_newyear_dragon_h_pic_03"].(string)
// 	a.DrawNumbersNewyearDragonHPic04, _ = m["draw_numbers_newyear_dragon_h_pic_04"].(string)
// 	a.DrawNumbersNewyearDragonHPic05, _ = m["draw_numbers_newyear_dragon_h_pic_05"].(string)
// 	a.DrawNumbersNewyearDragonHPic06, _ = m["draw_numbers_newyear_dragon_h_pic_06"].(string)
// 	a.DrawNumbersNewyearDragonHPic07, _ = m["draw_numbers_newyear_dragon_h_pic_07"].(string)
// 	a.DrawNumbersNewyearDragonHPic08, _ = m["draw_numbers_newyear_dragon_h_pic_08"].(string)
// 	a.DrawNumbersNewyearDragonHPic09, _ = m["draw_numbers_newyear_dragon_h_pic_09"].(string)
// 	a.DrawNumbersNewyearDragonHPic10, _ = m["draw_numbers_newyear_dragon_h_pic_10"].(string)
// 	a.DrawNumbersNewyearDragonHPic11, _ = m["draw_numbers_newyear_dragon_h_pic_11"].(string)
// 	a.DrawNumbersNewyearDragonHPic12, _ = m["draw_numbers_newyear_dragon_h_pic_12"].(string)
// 	a.DrawNumbersNewyearDragonHPic13, _ = m["draw_numbers_newyear_dragon_h_pic_13"].(string)
// 	a.DrawNumbersNewyearDragonHPic14, _ = m["draw_numbers_newyear_dragon_h_pic_14"].(string)
// 	a.DrawNumbersNewyearDragonHPic15, _ = m["draw_numbers_newyear_dragon_h_pic_15"].(string)
// 	a.DrawNumbersNewyearDragonHPic16, _ = m["draw_numbers_newyear_dragon_h_pic_16"].(string)
// 	a.DrawNumbersNewyearDragonHPic17, _ = m["draw_numbers_newyear_dragon_h_pic_17"].(string)
// 	a.DrawNumbersNewyearDragonHPic18, _ = m["draw_numbers_newyear_dragon_h_pic_18"].(string)
// 	a.DrawNumbersNewyearDragonHPic19, _ = m["draw_numbers_newyear_dragon_h_pic_19"].(string)
// 	a.DrawNumbersNewyearDragonHPic20, _ = m["draw_numbers_newyear_dragon_h_pic_20"].(string)
// 	a.DrawNumbersNewyearDragonHAni01, _ = m["draw_numbers_newyear_dragon_h_ani_01"].(string)
// 	a.DrawNumbersNewyearDragonHAni02, _ = m["draw_numbers_newyear_dragon_h_ani_02"].(string)

// 	a.DrawNumbersCherryHPic01, _ = m["draw_numbers_cherry_h_pic_01"].(string)
// 	a.DrawNumbersCherryHPic02, _ = m["draw_numbers_cherry_h_pic_02"].(string)
// 	a.DrawNumbersCherryHPic03, _ = m["draw_numbers_cherry_h_pic_03"].(string)
// 	a.DrawNumbersCherryHPic04, _ = m["draw_numbers_cherry_h_pic_04"].(string)
// 	a.DrawNumbersCherryHPic05, _ = m["draw_numbers_cherry_h_pic_05"].(string)
// 	a.DrawNumbersCherryHPic06, _ = m["draw_numbers_cherry_h_pic_06"].(string)
// 	a.DrawNumbersCherryHPic07, _ = m["draw_numbers_cherry_h_pic_07"].(string)
// 	a.DrawNumbersCherryHPic08, _ = m["draw_numbers_cherry_h_pic_08"].(string)
// 	a.DrawNumbersCherryHPic09, _ = m["draw_numbers_cherry_h_pic_09"].(string)
// 	a.DrawNumbersCherryHPic10, _ = m["draw_numbers_cherry_h_pic_10"].(string)
// 	a.DrawNumbersCherryHPic11, _ = m["draw_numbers_cherry_h_pic_11"].(string)
// 	a.DrawNumbersCherryHPic12, _ = m["draw_numbers_cherry_h_pic_12"].(string)
// 	a.DrawNumbersCherryHPic13, _ = m["draw_numbers_cherry_h_pic_13"].(string)
// 	a.DrawNumbersCherryHPic14, _ = m["draw_numbers_cherry_h_pic_14"].(string)
// 	a.DrawNumbersCherryHPic15, _ = m["draw_numbers_cherry_h_pic_15"].(string)
// 	a.DrawNumbersCherryHPic16, _ = m["draw_numbers_cherry_h_pic_16"].(string)
// 	a.DrawNumbersCherryHPic17, _ = m["draw_numbers_cherry_h_pic_17"].(string)
// 	a.DrawNumbersCherryHAni01, _ = m["draw_numbers_cherry_h_ani_01"].(string)
// 	a.DrawNumbersCherryHAni02, _ = m["draw_numbers_cherry_h_ani_02"].(string)
// 	a.DrawNumbersCherryHAni03, _ = m["draw_numbers_cherry_h_ani_03"].(string)
// 	a.DrawNumbersCherryHAni04, _ = m["draw_numbers_cherry_h_ani_04"].(string)

// 	a.DrawNumbers3DSpaceHPic01, _ = m["draw_numbers_3D_space_h_pic_01"].(string)
// 	a.DrawNumbers3DSpaceHPic02, _ = m["draw_numbers_3D_space_h_pic_02"].(string)
// 	a.DrawNumbers3DSpaceHPic03, _ = m["draw_numbers_3D_space_h_pic_03"].(string)
// 	a.DrawNumbers3DSpaceHPic04, _ = m["draw_numbers_3D_space_h_pic_04"].(string)
// 	a.DrawNumbers3DSpaceHPic05, _ = m["draw_numbers_3D_space_h_pic_05"].(string)
// 	a.DrawNumbers3DSpaceHPic06, _ = m["draw_numbers_3D_space_h_pic_06"].(string)
// 	a.DrawNumbers3DSpaceHPic07, _ = m["draw_numbers_3D_space_h_pic_07"].(string)
// 	a.DrawNumbers3DSpaceHPic08, _ = m["draw_numbers_3D_space_h_pic_08"].(string)

// 	// 音樂
// 	a.DrawNumbersBgmGaming, _ = m["draw_numbers_bgm_gaming"].(string)
// }

// if a.Game == "whack_mole" {
// 	// 敲敲樂自定義圖片
// 	a.WhackmoleClassicHPic01, _ = m["whackmole_classic_h_pic_01"].(string)
// 	a.WhackmoleClassicHPic02, _ = m["whackmole_classic_h_pic_02"].(string)
// 	a.WhackmoleClassicHPic03, _ = m["whackmole_classic_h_pic_03"].(string)
// 	a.WhackmoleClassicHPic04, _ = m["whackmole_classic_h_pic_04"].(string)
// 	a.WhackmoleClassicHPic05, _ = m["whackmole_classic_h_pic_05"].(string)
// 	a.WhackmoleClassicHPic06, _ = m["whackmole_classic_h_pic_06"].(string)
// 	a.WhackmoleClassicHPic07, _ = m["whackmole_classic_h_pic_07"].(string)
// 	a.WhackmoleClassicHPic08, _ = m["whackmole_classic_h_pic_08"].(string)
// 	a.WhackmoleClassicHPic09, _ = m["whackmole_classic_h_pic_09"].(string)
// 	a.WhackmoleClassicHPic10, _ = m["whackmole_classic_h_pic_10"].(string)
// 	a.WhackmoleClassicHPic11, _ = m["whackmole_classic_h_pic_11"].(string)
// 	a.WhackmoleClassicHPic12, _ = m["whackmole_classic_h_pic_12"].(string)
// 	a.WhackmoleClassicHPic13, _ = m["whackmole_classic_h_pic_13"].(string)
// 	a.WhackmoleClassicHPic14, _ = m["whackmole_classic_h_pic_14"].(string)
// 	a.WhackmoleClassicHPic15, _ = m["whackmole_classic_h_pic_15"].(string)
// 	a.WhackmoleClassicGPic01, _ = m["whackmole_classic_g_pic_01"].(string)
// 	a.WhackmoleClassicGPic02, _ = m["whackmole_classic_g_pic_02"].(string)
// 	a.WhackmoleClassicGPic03, _ = m["whackmole_classic_g_pic_03"].(string)
// 	a.WhackmoleClassicGPic04, _ = m["whackmole_classic_g_pic_04"].(string)
// 	a.WhackmoleClassicGPic05, _ = m["whackmole_classic_g_pic_05"].(string)
// 	a.WhackmoleClassicCPic01, _ = m["whackmole_classic_c_pic_01"].(string)
// 	a.WhackmoleClassicCPic02, _ = m["whackmole_classic_c_pic_02"].(string)
// 	a.WhackmoleClassicCPic03, _ = m["whackmole_classic_c_pic_03"].(string)
// 	a.WhackmoleClassicCPic04, _ = m["whackmole_classic_c_pic_04"].(string)
// 	a.WhackmoleClassicCPic05, _ = m["whackmole_classic_c_pic_05"].(string)
// 	a.WhackmoleClassicCPic06, _ = m["whackmole_classic_c_pic_06"].(string)
// 	a.WhackmoleClassicCPic07, _ = m["whackmole_classic_c_pic_07"].(string)
// 	a.WhackmoleClassicCPic08, _ = m["whackmole_classic_c_pic_08"].(string)
// 	a.WhackmoleClassicCAni01, _ = m["whackmole_classic_c_ani_01"].(string)

// 	a.WhackmoleHalloweenHPic01, _ = m["whackmole_halloween_h_pic_01"].(string)
// 	a.WhackmoleHalloweenHPic02, _ = m["whackmole_halloween_h_pic_02"].(string)
// 	a.WhackmoleHalloweenHPic03, _ = m["whackmole_halloween_h_pic_03"].(string)
// 	a.WhackmoleHalloweenHPic04, _ = m["whackmole_halloween_h_pic_04"].(string)
// 	a.WhackmoleHalloweenHPic05, _ = m["whackmole_halloween_h_pic_05"].(string)
// 	a.WhackmoleHalloweenHPic06, _ = m["whackmole_halloween_h_pic_06"].(string)
// 	a.WhackmoleHalloweenHPic07, _ = m["whackmole_halloween_h_pic_07"].(string)
// 	a.WhackmoleHalloweenHPic08, _ = m["whackmole_halloween_h_pic_08"].(string)
// 	a.WhackmoleHalloweenHPic09, _ = m["whackmole_halloween_h_pic_09"].(string)
// 	a.WhackmoleHalloweenHPic10, _ = m["whackmole_halloween_h_pic_10"].(string)
// 	a.WhackmoleHalloweenHPic11, _ = m["whackmole_halloween_h_pic_11"].(string)
// 	a.WhackmoleHalloweenHPic12, _ = m["whackmole_halloween_h_pic_12"].(string)
// 	a.WhackmoleHalloweenHPic13, _ = m["whackmole_halloween_h_pic_13"].(string)
// 	a.WhackmoleHalloweenHPic14, _ = m["whackmole_halloween_h_pic_14"].(string)
// 	a.WhackmoleHalloweenHPic15, _ = m["whackmole_halloween_h_pic_15"].(string)
// 	a.WhackmoleHalloweenGPic01, _ = m["whackmole_halloween_g_pic_01"].(string)
// 	a.WhackmoleHalloweenGPic02, _ = m["whackmole_halloween_g_pic_02"].(string)
// 	a.WhackmoleHalloweenGPic03, _ = m["whackmole_halloween_g_pic_03"].(string)
// 	a.WhackmoleHalloweenGPic04, _ = m["whackmole_halloween_g_pic_04"].(string)
// 	a.WhackmoleHalloweenGPic05, _ = m["whackmole_halloween_g_pic_05"].(string)
// 	a.WhackmoleHalloweenCPic01, _ = m["whackmole_halloween_c_pic_01"].(string)
// 	a.WhackmoleHalloweenCPic02, _ = m["whackmole_halloween_c_pic_02"].(string)
// 	a.WhackmoleHalloweenCPic03, _ = m["whackmole_halloween_c_pic_03"].(string)
// 	a.WhackmoleHalloweenCPic04, _ = m["whackmole_halloween_c_pic_04"].(string)
// 	a.WhackmoleHalloweenCPic05, _ = m["whackmole_halloween_c_pic_05"].(string)
// 	a.WhackmoleHalloweenCPic06, _ = m["whackmole_halloween_c_pic_06"].(string)
// 	a.WhackmoleHalloweenCPic07, _ = m["whackmole_halloween_c_pic_07"].(string)
// 	a.WhackmoleHalloweenCPic08, _ = m["whackmole_halloween_c_pic_08"].(string)
// 	a.WhackmoleHalloweenCAni01, _ = m["whackmole_halloween_c_ani_01"].(string)

// 	a.WhackmoleChristmasHPic01, _ = m["whackmole_christmas_h_pic_01"].(string)
// 	a.WhackmoleChristmasHPic02, _ = m["whackmole_christmas_h_pic_02"].(string)
// 	a.WhackmoleChristmasHPic03, _ = m["whackmole_christmas_h_pic_03"].(string)
// 	a.WhackmoleChristmasHPic04, _ = m["whackmole_christmas_h_pic_04"].(string)
// 	a.WhackmoleChristmasHPic05, _ = m["whackmole_christmas_h_pic_05"].(string)
// 	a.WhackmoleChristmasHPic06, _ = m["whackmole_christmas_h_pic_06"].(string)
// 	a.WhackmoleChristmasHPic07, _ = m["whackmole_christmas_h_pic_07"].(string)
// 	a.WhackmoleChristmasHPic08, _ = m["whackmole_christmas_h_pic_08"].(string)
// 	a.WhackmoleChristmasHPic09, _ = m["whackmole_christmas_h_pic_09"].(string)
// 	a.WhackmoleChristmasHPic10, _ = m["whackmole_christmas_h_pic_10"].(string)
// 	a.WhackmoleChristmasHPic11, _ = m["whackmole_christmas_h_pic_11"].(string)
// 	a.WhackmoleChristmasHPic12, _ = m["whackmole_christmas_h_pic_12"].(string)
// 	a.WhackmoleChristmasHPic13, _ = m["whackmole_christmas_h_pic_13"].(string)
// 	a.WhackmoleChristmasHPic14, _ = m["whackmole_christmas_h_pic_14"].(string)
// 	a.WhackmoleChristmasGPic01, _ = m["whackmole_christmas_g_pic_01"].(string)
// 	a.WhackmoleChristmasGPic02, _ = m["whackmole_christmas_g_pic_02"].(string)
// 	a.WhackmoleChristmasGPic03, _ = m["whackmole_christmas_g_pic_03"].(string)
// 	a.WhackmoleChristmasGPic04, _ = m["whackmole_christmas_g_pic_04"].(string)
// 	a.WhackmoleChristmasGPic05, _ = m["whackmole_christmas_g_pic_05"].(string)
// 	a.WhackmoleChristmasGPic06, _ = m["whackmole_christmas_g_pic_06"].(string)
// 	a.WhackmoleChristmasGPic07, _ = m["whackmole_christmas_g_pic_07"].(string)
// 	a.WhackmoleChristmasGPic08, _ = m["whackmole_christmas_g_pic_08"].(string)
// 	a.WhackmoleChristmasCPic01, _ = m["whackmole_christmas_c_pic_01"].(string)
// 	a.WhackmoleChristmasCPic02, _ = m["whackmole_christmas_c_pic_02"].(string)
// 	a.WhackmoleChristmasCPic03, _ = m["whackmole_christmas_c_pic_03"].(string)
// 	a.WhackmoleChristmasCPic04, _ = m["whackmole_christmas_c_pic_04"].(string)
// 	a.WhackmoleChristmasCPic05, _ = m["whackmole_christmas_c_pic_05"].(string)
// 	a.WhackmoleChristmasCPic06, _ = m["whackmole_christmas_c_pic_06"].(string)
// 	a.WhackmoleChristmasCPic07, _ = m["whackmole_christmas_c_pic_07"].(string)
// 	a.WhackmoleChristmasCPic08, _ = m["whackmole_christmas_c_pic_08"].(string)
// 	a.WhackmoleChristmasCAni01, _ = m["whackmole_christmas_c_ani_01"].(string)
// 	a.WhackmoleChristmasCAni02, _ = m["whackmole_christmas_c_ani_02"].(string)

// 	// 音樂
// 	a.WhackmoleBgmStart, _ = m["whackmole_bgm_start"].(string)
// 	a.WhackmoleBgmGaming, _ = m["whackmole_bgm_gaming"].(string)
// 	a.WhackmoleBgmEnd, _ = m["whackmole_bgm_end"].(string)
// }

// if a.Game == "monopoly" {
// 	// 鑑定師自定義
// 	a.MonopolyClassicHPic01, _ = m["monopoly_classic_h_pic_01"].(string)
// 	a.MonopolyClassicHPic02, _ = m["monopoly_classic_h_pic_02"].(string)
// 	a.MonopolyClassicHPic03, _ = m["monopoly_classic_h_pic_03"].(string)
// 	a.MonopolyClassicHPic04, _ = m["monopoly_classic_h_pic_04"].(string)
// 	a.MonopolyClassicHPic05, _ = m["monopoly_classic_h_pic_05"].(string)
// 	a.MonopolyClassicHPic06, _ = m["monopoly_classic_h_pic_06"].(string)
// 	a.MonopolyClassicHPic07, _ = m["monopoly_classic_h_pic_07"].(string)
// 	a.MonopolyClassicHPic08, _ = m["monopoly_classic_h_pic_08"].(string)
// 	a.MonopolyClassicGPic01, _ = m["monopoly_classic_g_pic_01"].(string)
// 	a.MonopolyClassicGPic02, _ = m["monopoly_classic_g_pic_02"].(string)
// 	a.MonopolyClassicGPic03, _ = m["monopoly_classic_g_pic_03"].(string)
// 	a.MonopolyClassicGPic04, _ = m["monopoly_classic_g_pic_04"].(string)
// 	a.MonopolyClassicGPic05, _ = m["monopoly_classic_g_pic_05"].(string)
// 	a.MonopolyClassicGPic06, _ = m["monopoly_classic_g_pic_06"].(string)
// 	a.MonopolyClassicGPic07, _ = m["monopoly_classic_g_pic_07"].(string)
// 	a.MonopolyClassicCPic01, _ = m["monopoly_classic_c_pic_01"].(string)
// 	a.MonopolyClassicCPic02, _ = m["monopoly_classic_c_pic_02"].(string)
// 	a.MonopolyClassicGAni01, _ = m["monopoly_classic_g_ani_01"].(string)
// 	a.MonopolyClassicGAni02, _ = m["monopoly_classic_g_ani_02"].(string)
// 	a.MonopolyClassicCAni01, _ = m["monopoly_classic_c_ani_01"].(string)

// 	a.MonopolyRedpackHPic01, _ = m["monopoly_redpack_h_pic_01"].(string)
// 	a.MonopolyRedpackHPic02, _ = m["monopoly_redpack_h_pic_02"].(string)
// 	a.MonopolyRedpackHPic03, _ = m["monopoly_redpack_h_pic_03"].(string)
// 	a.MonopolyRedpackHPic04, _ = m["monopoly_redpack_h_pic_04"].(string)
// 	a.MonopolyRedpackHPic05, _ = m["monopoly_redpack_h_pic_05"].(string)
// 	a.MonopolyRedpackHPic06, _ = m["monopoly_redpack_h_pic_06"].(string)
// 	a.MonopolyRedpackHPic07, _ = m["monopoly_redpack_h_pic_07"].(string)
// 	a.MonopolyRedpackHPic08, _ = m["monopoly_redpack_h_pic_08"].(string)
// 	a.MonopolyRedpackHPic09, _ = m["monopoly_redpack_h_pic_09"].(string)
// 	a.MonopolyRedpackHPic10, _ = m["monopoly_redpack_h_pic_10"].(string)
// 	a.MonopolyRedpackHPic11, _ = m["monopoly_redpack_h_pic_11"].(string)
// 	a.MonopolyRedpackHPic12, _ = m["monopoly_redpack_h_pic_12"].(string)
// 	a.MonopolyRedpackHPic13, _ = m["monopoly_redpack_h_pic_13"].(string)
// 	a.MonopolyRedpackHPic14, _ = m["monopoly_redpack_h_pic_14"].(string)
// 	a.MonopolyRedpackHPic15, _ = m["monopoly_redpack_h_pic_15"].(string)
// 	a.MonopolyRedpackHPic16, _ = m["monopoly_redpack_h_pic_16"].(string)
// 	a.MonopolyRedpackGPic01, _ = m["monopoly_redpack_g_pic_01"].(string)
// 	a.MonopolyRedpackGPic02, _ = m["monopoly_redpack_g_pic_02"].(string)
// 	a.MonopolyRedpackGPic03, _ = m["monopoly_redpack_g_pic_03"].(string)
// 	a.MonopolyRedpackGPic04, _ = m["monopoly_redpack_g_pic_04"].(string)
// 	a.MonopolyRedpackGPic05, _ = m["monopoly_redpack_g_pic_05"].(string)
// 	a.MonopolyRedpackGPic06, _ = m["monopoly_redpack_g_pic_06"].(string)
// 	a.MonopolyRedpackGPic07, _ = m["monopoly_redpack_g_pic_07"].(string)
// 	a.MonopolyRedpackGPic08, _ = m["monopoly_redpack_g_pic_08"].(string)
// 	a.MonopolyRedpackGPic09, _ = m["monopoly_redpack_g_pic_09"].(string)
// 	a.MonopolyRedpackGPic10, _ = m["monopoly_redpack_g_pic_10"].(string)
// 	a.MonopolyRedpackCPic01, _ = m["monopoly_redpack_c_pic_01"].(string)
// 	a.MonopolyRedpackCPic02, _ = m["monopoly_redpack_c_pic_02"].(string)
// 	a.MonopolyRedpackCPic03, _ = m["monopoly_redpack_c_pic_03"].(string)
// 	a.MonopolyRedpackHAni01, _ = m["monopoly_redpack_h_ani_01"].(string)
// 	a.MonopolyRedpackHAni02, _ = m["monopoly_redpack_h_ani_02"].(string)
// 	a.MonopolyRedpackHAni03, _ = m["monopoly_redpack_h_ani_03"].(string)
// 	a.MonopolyRedpackGAni01, _ = m["monopoly_redpack_g_ani_01"].(string)
// 	a.MonopolyRedpackGAni02, _ = m["monopoly_redpack_g_ani_02"].(string)
// 	a.MonopolyRedpackCAni01, _ = m["monopoly_redpack_c_ani_01"].(string)

// 	a.MonopolyNewyearRabbitHPic01, _ = m["monopoly_newyear_rabbit_h_pic_01"].(string)
// 	a.MonopolyNewyearRabbitHPic02, _ = m["monopoly_newyear_rabbit_h_pic_02"].(string)
// 	a.MonopolyNewyearRabbitHPic03, _ = m["monopoly_newyear_rabbit_h_pic_03"].(string)
// 	a.MonopolyNewyearRabbitHPic04, _ = m["monopoly_newyear_rabbit_h_pic_04"].(string)
// 	a.MonopolyNewyearRabbitHPic05, _ = m["monopoly_newyear_rabbit_h_pic_05"].(string)
// 	a.MonopolyNewyearRabbitHPic06, _ = m["monopoly_newyear_rabbit_h_pic_06"].(string)
// 	a.MonopolyNewyearRabbitHPic07, _ = m["monopoly_newyear_rabbit_h_pic_07"].(string)
// 	a.MonopolyNewyearRabbitHPic08, _ = m["monopoly_newyear_rabbit_h_pic_08"].(string)
// 	a.MonopolyNewyearRabbitHPic09, _ = m["monopoly_newyear_rabbit_h_pic_09"].(string)
// 	a.MonopolyNewyearRabbitHPic10, _ = m["monopoly_newyear_rabbit_h_pic_10"].(string)
// 	a.MonopolyNewyearRabbitHPic11, _ = m["monopoly_newyear_rabbit_h_pic_11"].(string)
// 	a.MonopolyNewyearRabbitHPic12, _ = m["monopoly_newyear_rabbit_h_pic_12"].(string)
// 	a.MonopolyNewyearRabbitGPic01, _ = m["monopoly_newyear_rabbit_g_pic_01"].(string)
// 	a.MonopolyNewyearRabbitGPic02, _ = m["monopoly_newyear_rabbit_g_pic_02"].(string)
// 	a.MonopolyNewyearRabbitGPic03, _ = m["monopoly_newyear_rabbit_g_pic_03"].(string)
// 	a.MonopolyNewyearRabbitGPic04, _ = m["monopoly_newyear_rabbit_g_pic_04"].(string)
// 	a.MonopolyNewyearRabbitGPic05, _ = m["monopoly_newyear_rabbit_g_pic_05"].(string)
// 	a.MonopolyNewyearRabbitGPic06, _ = m["monopoly_newyear_rabbit_g_pic_06"].(string)
// 	a.MonopolyNewyearRabbitGPic07, _ = m["monopoly_newyear_rabbit_g_pic_07"].(string)
// 	a.MonopolyNewyearRabbitCPic01, _ = m["monopoly_newyear_rabbit_c_pic_01"].(string)
// 	a.MonopolyNewyearRabbitCPic02, _ = m["monopoly_newyear_rabbit_c_pic_02"].(string)
// 	a.MonopolyNewyearRabbitCPic03, _ = m["monopoly_newyear_rabbit_c_pic_03"].(string)
// 	a.MonopolyNewyearRabbitHAni01, _ = m["monopoly_newyear_rabbit_h_ani_01"].(string)
// 	a.MonopolyNewyearR

// 查詢資料表
// var (
// sql = a.Table(a.Base.TableName).Select("game").Where("activity_game.game_id", "=", gameID)
// )

// item, err := sql.First()
// if err != nil {
// 	return "", errors.New("錯誤: 無法取得遊戲種類資訊，請重新查詢")
// }
// if item == nil {
// 	return "", nil
// }

// 取得遊戲種類
// game, _ = item["game"].(string)

// 設置過期時間
// a.RedisConn.SetExpire(config.GAME_TYPE_REDIS+gameID,
// 	config.REDIS_EXPIRE)

// sql = a.Table(a.Base.TableName).
// 	Select(
// 		"activity_game.id", "activity_game.activity_id", "activity_game.user_id",
// 		"activity_game.game_id", "activity_game.game", "activity_game.title",
// 		"activity_game.game_type", "activity_game.limit_time", "activity_game.second",
// 		"activity_game.max_people", "activity_game.people",
// 		"activity_game.max_times", "activity_game.allow",
// 		"activity_game.percent", "activity_game.first_prize", "activity_game.second_prize",
// 		"activity_game.third_prize", "activity_game.general_prize", "activity_game.topic",
// 		"activity_game.skin", "activity_game.music",
// 		"activity_game.game_round", "activity_game.game_second",
// 		"activity_game.game_status", "activity_game.game_attend", "activity_game.display_name",
// 		"activity_game.game_order",
// 		"activity_game.box_reflection",
// 		"activity_game.same_people",

// 		// 投票遊戲
// 		"activity_game.vote_screen",
// 		"activity_game.vote_times",
// 		"activity_game.vote_method",
// 		"activity_game.vote_method_player",
// 		"activity_game.vote_restriction",
// 		"activity_game.avatar_shape",
// 		"activity_game.vote_start_time",
// 		"activity_game.vote_end_time",
// 		"activity_game.auto_play",
// 		"activity_game.show_rank",
// 		"activity_game.title_switch",
// 		"activity_game.arrangement_guest",

// 		"activity_game.allow_choose_team", "activity_game.left_team_name", "activity_game.left_team_picture",
// 		"activity_game.right_team_name", "activity_game.right_team_picture", "activity_game.left_team_game_attend",
// 		"activity_game.right_team_game_attend", "activity_game.prize",

// 		"activity_game.max_number", "activity_game.bingo_line",
// 		"activity_game.round_prize", "activity_game.bingo_round",

// 		"activity_game.edit_times",

// 		"activity_game_qa.total_qa", "activity_game_qa.qa_second",
// 		"activity_game_qa.qa_round", "activity_game_qa.qa_people",

// 		"activity_game.gacha_machine_reflection",
// 		"activity_game.reflective_switch",
// 	).
// 	LeftJoin(command.Join{
// 		FieldA:    "activity_game.game_id",
// 		FieldA1:   "activity_game_qa.game_id",
// 		Table:     "activity_game_qa",
// 		Operation: "="}).
// 	Where("activity_game.activity_id", "=", activityID).
// 	OrderBy("activity_game.game_order", "asc")
// OrderBy("activity_game.game", "asc")
// items = make([]map[string]interface{}, 0)
// err   error

// 查詢特定遊戲種類所有遊戲場次
// sql = sql.WhereIn("game", interfaces(strings.Split(game, ",")))

// items, err = sql.All()
// if err != nil {
// 	return nil, errors.New("錯誤: 無法取得所有場次遊戲資訊，請重新查詢")
// }
