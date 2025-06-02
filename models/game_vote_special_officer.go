package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// GameVoteSpecialOfficerModel 資料表欄位
type GameVoteSpecialOfficerModel struct {
	Base             `json:"-"`
	ID               int64  `json:"id"`
	ActivityID       string `json:"activity_id" example:"activity_id"`
	GameID           string `json:"game_id" example:"game_id"`
	UserID           string `json:"user_id" example:"user_id"`
	Score            int64  `json:"score" example:"0"`              // 特殊人員權重分數
	Times            int64  `json:"times" example:"0"`              // 特殊人員投票次數
	VoteMethodPlayer int64  `json:"vote_method_player" example:"0"` // 特殊人員投票方式

	Name string `json:"name" example:"name"`
}

// EditGameVoteSpecialOfficerModel 資料表欄位
type EditGameVoteSpecialOfficerModel struct {
	ActivityID       string `json:"activity_id" example:"activity_id"`
	GameID           string `json:"game_id" example:"game_id"`
	UserID           string `json:"user_id" example:"user_id"`
	Score            string `json:"score" example:"0"`              // 特殊人員權重分數
	Times            string `json:"times" example:"0"`              // 特殊人員投票次數
	VoteMethodPlayer string `json:"vote_method_player" example:"0"` // 特殊人員投票方式

	People string `json:"people"` // 匯入資料數
}

// DefaultGameVoteSpecialOfficerModel 預設GameVoteSpecialOfficerModel
func DefaultGameVoteSpecialOfficerModel() GameVoteSpecialOfficerModel {
	return GameVoteSpecialOfficerModel{Base: Base{TableName: config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameVoteSpecialOfficerModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameVoteSpecialOfficerModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameVoteSpecialOfficerModel) SetDbConn(conn db.Connection) GameVoteSpecialOfficerModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameVoteSpecialOfficerModel) SetRedisConn(conn cache.Connection) GameVoteSpecialOfficerModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (a GameVoteSpecialOfficerModel) SetMongoConn(conn mongo.Connection) GameVoteSpecialOfficerModel {
// 	a.MongoConn = conn
// 	return a
// }

// Find 查詢特殊人員資料
func (a GameVoteSpecialOfficerModel) Find(isRedis bool, gameID, userID string) ([]GameVoteSpecialOfficerModel, error) {
	// log.Println("查詢特殊人員資料")

	var (
		sql = a.Table(a.Base.TableName).
			Select(
				"activity_game_vote_special_officer.id",
				"activity_game_vote_special_officer.activity_id",
				"activity_game_vote_special_officer.game_id",
				"activity_game_vote_special_officer.user_id",
				"activity_game_vote_special_officer.score",
				"activity_game_vote_special_officer.times",
				"activity_game_vote_special_officer.vote_method_player",

				// 用戶
				"line_users.name", "line_users.avatar",
			).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_special_officer.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("game_id", "=", gameID).
			OrderBy("id", "asc")
		officers = make([]GameVoteSpecialOfficerModel, 0)
		// datas    map[string]string
		// err      error
	)

	if userID != "" { // 查詢特定用戶
		sql = sql.Where("activity_game_vote_special_officer.user_id", "=", userID)
	}

	if isRedis && userID != "" {
		// redis處理
		// 取得特殊人員分數票數權重資料
		data, err := a.RedisConn.HashGetCache(config.VOTE_SPECIAL_OFFICER_REDIS+gameID, userID) // 包含test測試資料，避免重複執行資料庫
		if err != nil {
			// redis查詢不到用戶資料
			// 從資料庫讀取特殊人員資料
			// log.Println("redis無法取得特殊人員資料，查詢資料表")

			items, err := sql.All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得特殊人員資訊，請重新查詢")
			}
			officers = MapToVoteSpecialOfficerModel(items) // 應該為1筆資料

			if len(officers) > 0 {
				// log.Println("資料表有資料，加入特殊人員資料至redis")

				// 將特殊人員分數權重寫入redis
				err = a.RedisConn.HashSetCache(config.VOTE_SPECIAL_OFFICER_REDIS+gameID, userID, utils.JSON(officers[0]))
				if err != nil {
					return nil, errors.New("錯誤: 將特殊人員分數權重寫入redis中發生問題")
				}

				// a.RedisConn.SetExpire(config.VOTE_SPECIAL_OFFICER_REDIS+gameID,
				// 	config.REDIS_EXPIRE)
			} else if len(officers) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				err = a.RedisConn.HashSetCache(config.VOTE_SPECIAL_OFFICER_REDIS+gameID, userID,
					utils.JSON(GameVoteSpecialOfficerModel{
						UserID: "test",
					}))
				if err != nil {
					return nil, errors.New("錯誤: 將特殊人員分數權重寫入redis中發生問題")
				}

				// a.RedisConn.SetExpire(config.VOTE_SPECIAL_OFFICER_REDIS+gameID,
				// 	config.REDIS_EXPIRE)
			}
		} else {
			// log.Println("redis有資料")

			// redis成功取得用戶分數權重資訊
			var gameSpecialOfficerModel GameVoteSpecialOfficerModel
			// 解碼，取得特殊人員資料
			json.Unmarshal([]byte(data), &gameSpecialOfficerModel)

			// log.Println("資料: ", gameSpecialOfficerModel)

			if gameSpecialOfficerModel.UserID != "test" {
				// 過濾測試資料
				officers = append(officers, gameSpecialOfficerModel)
			} else {
				// log.Println("目前為測試資料，跳過")
			}

		}
	}

	// 查詢資料表
	if !isRedis {
		// log.Println("查詢資料表: ")

		// 資料庫取得資料
		items, err := sql.All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得特殊人員資訊，請重新查詢")
		}
		officers = MapToVoteSpecialOfficerModel(items)

		// 將資料寫入redis中(redis的value值為struct，可以將所有特殊人員資料一起寫入)
		var params = []interface{}{config.VOTE_SPECIAL_OFFICER_REDIS + gameID} // 特殊人員資料

		for _, officer := range officers {
			// if officer.UserID != "" {
			params = append(params, officer.UserID, utils.JSON(officer))
			// }
		}

		if len(officers) == 0 {
			// log.Println("資料庫無資料，加入test")

			// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
			params = append(params, "test", utils.JSON(GameVoteSpecialOfficerModel{
				UserID: "test",
			}))
		}

		// 將特殊人員分數權重寫入redis
		err = a.RedisConn.HashMultiSetCache(params)
		if err != nil {
			return nil, errors.New("錯誤: 將特殊人員分數權重寫入redis中發生問題")
		}

		// a.RedisConn.SetExpire(config.VOTE_SPECIAL_OFFICER_REDIS+gameID,
		// 	config.REDIS_EXPIRE)

	}

	return officers, nil
}

// MapToVoteSpecialOfficerModel map轉換[]GameVoteSpecialOfficerModel
func MapToVoteSpecialOfficerModel(items []map[string]interface{}) []GameVoteSpecialOfficerModel {
	var specialOfficers = make([]GameVoteSpecialOfficerModel, 0)
	for _, item := range items {
		var (
			specialOfficer GameVoteSpecialOfficerModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &specialOfficer)

		specialOfficers = append(specialOfficers, specialOfficer)
	}
	return specialOfficers
}

// Adds 新增多筆特殊人員
func (a GameVoteSpecialOfficerModel) Adds(isRedis bool, model EditGameVoteSpecialOfficerModel) error {
	// 取得匯入總資料數
	dataAmount, err := strconv.Atoi(model.People)
	if err != nil {
		return errors.New("錯誤: 人數資料發生問題，請重新操作")
	}

	if dataAmount > 0 {

		// 處理匯入的參數資料
		var (
			activityID       = make([]string, dataAmount)                 // 活動ID
			gameID           = make([]string, dataAmount)                 // 遊戲ID
			userID           = strings.Split(model.UserID, ",")           // 用戶ID
			score            = strings.Split(model.Score, ",")            // 權重分數
			times            = strings.Split(model.Times, ",")            // 權重分數
			voteMethodPlayer = strings.Split(model.VoteMethodPlayer, ",") // 權重分數
		)

		// 處理活動.遊戲ID陣列資料
		for i := range userID {
			activityID[i] = model.ActivityID
			gameID[i] = model.GameID
		}

		// 更新遊戲場次編輯次數(mysq，會影響遊戲所以需重整)
		// if err := a.Table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", model.GameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return err
		// }

		// 更新遊戲場次編輯次數(mongo，會影響遊戲所以需重整)
		filter := bson.M{"game_id": model.GameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 將資料匹量寫入activity_game_vote_special_officer表中
		err = a.Table(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).BatchInsert(
			dataAmount, "activity_id,game_id,user_id,score,times,vote_method_player",
			[][]string{activityID, gameID, userID, score, times, voteMethodPlayer})

		if err != nil {
			return errors.New("錯誤: 匹量新增特殊人員發生問題(activity_game_vote_special_officer)，請重新操作")
		}

		if isRedis {
			// redis處理

			// 從redis取得資料，確定redis中有該場遊戲的投票特殊人員資料
			a.Find(false, model.GameID, "")

			// 將特殊人員分數權重寫入redis
			params := []interface{}{config.VOTE_SPECIAL_OFFICER_REDIS + model.GameID}

			for i, userIDStr := range userID {
				scoreInt, _ := strconv.Atoi(score[i])
				timesInt, _ := strconv.Atoi(times[i])
				voteMethodPlayerInt, _ := strconv.Atoi(voteMethodPlayer[i])

				// 將用戶權重分數加入redis中
				params = append(params, userIDStr, utils.JSON(GameVoteSpecialOfficerModel{
					ActivityID:       model.ActivityID,
					GameID:           model.GameID,
					UserID:           userIDStr,
					Score:            int64(scoreInt),
					Times:            int64(timesInt),
					VoteMethodPlayer: int64(voteMethodPlayerInt),
				}))
			}

			err := a.RedisConn.HashMultiSetCache(params)
			if err != nil {
				return errors.New("錯誤: 將特殊人員分數權重寫入redis中發生問題")
			}
		}

		// 編輯次數更新
		a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")
	}
	return nil
}

// Add 新增特殊人員資料(加入資料庫)
func (a GameVoteSpecialOfficerModel) Add(isRedis bool, model EditGameVoteSpecialOfficerModel) error {
	var (
		fieldValues = command.Value{
			"activity_id":        model.ActivityID,
			"game_id":            model.GameID,
			"user_id":            model.UserID,
			"score":              model.Score,
			"times":              model.Times,
			"vote_method_player": model.VoteMethodPlayer,
		}
	)

	// 更新遊戲場次編輯次數(mysql，會影響遊戲所以需重整)
	// if err := a.Table(config.ACTIVITY_GAME_TABLE).
	// 	Where("game_id", "=", model.GameID).
	// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return err
	// }

	// 更新遊戲場次編輯次數(mongo，會影響遊戲所以需重整)
	filter := bson.M{"game_id": model.GameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		// "$set": fieldValues,
		// "$unset": bson.M{},                // 移除不需要的欄位
		"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return err
	}

	id, err := a.Table(a.TableName).Insert(fieldValues)
	if err != nil {
		return errors.New("錯誤: 新增特殊人員資料發生問題")
	}

	if isRedis {
		// redis處理

		// 從redis取得資料，確定redis中有該場遊戲的投票特殊人員資料
		a.Find(false, model.GameID, "")

		scoreInt, _ := strconv.Atoi(model.Score)
		timesInt, _ := strconv.Atoi(model.Times)
		voteMethodPlayerInt, _ := strconv.Atoi(model.VoteMethodPlayer)

		// 將特殊人員分數權重寫入redis
		err := a.RedisConn.HashSetCache(config.VOTE_SPECIAL_OFFICER_REDIS+model.GameID,
			model.UserID, utils.JSON(GameVoteSpecialOfficerModel{
				ID:               id,
				ActivityID:       model.ActivityID,
				GameID:           model.GameID,
				UserID:           model.UserID,
				Score:            int64(scoreInt),
				Times:            int64(timesInt),
				VoteMethodPlayer: int64(voteMethodPlayerInt),
			}))
		if err != nil {
			return errors.New("錯誤: 將特殊人員分數權重寫入redis中發生問題")
		}
	}

	// 編輯次數更新
	a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

	return nil
}

// Update 更新特殊人員資料
func (a GameVoteSpecialOfficerModel) Update(isRedis bool, model EditGameVoteSpecialOfficerModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"score", "times", "vote_method_player"}
		values      = []string{model.Score, model.Times, model.VoteMethodPlayer}
	)

	// 更新遊戲場次編輯次數(mysql，會影響遊戲所以需重整)
	// if err := a.Table(config.ACTIVITY_GAME_TABLE).
	// 	Where("game_id", "=", model.GameID).
	// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return err
	// }

	// 更新遊戲場次編輯次數(mongo，會影響遊戲所以需重整)
	filter := bson.M{"game_id": model.GameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		// "$set": fieldValues,
		// "$unset": bson.M{},                // 移除不需要的欄位
		"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return err
	}

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	if isRedis {
		// redis處理

		// 從redis取得資料，確定redis中有該場遊戲的投票特殊人員資料
		a.Find(false, model.GameID, "")

		scoreInt, _ := strconv.Atoi(model.Score)
		timesInt, _ := strconv.Atoi(model.Times)
		voteMethodPlayerInt, _ := strconv.Atoi(model.VoteMethodPlayer)

		// 更新redis中特殊人員分數權重資訊
		a.RedisConn.HashSetCache(config.VOTE_SPECIAL_OFFICER_REDIS+model.GameID,
			model.UserID, utils.JSON(GameVoteSpecialOfficerModel{
				ActivityID:       model.ActivityID,
				GameID:           model.GameID,
				UserID:           model.UserID,
				Score:            int64(scoreInt),
				Times:            int64(timesInt),
				VoteMethodPlayer: int64(voteMethodPlayerInt),
			}))
	}

	err := a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).
		Where("game_id", "=", model.GameID).
		Where("user_id", "=", model.UserID).
		Update(fieldValues)
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return err
	}

	// 編輯次數更新
	a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

	return nil
}

// specialOfficer.ID, _ = item["id"].(int64)
// specialOfficer.ActivityID, _ = item["activity_id"].(string)
// specialOfficer.GameID, _ = item["game_id"].(string)
// specialOfficer.UserID, _ = item["user_id"].(string)
// specialOfficer.Score, _ = item["score"].(int64)
// specialOfficer.Times, _ = item["times"].(int64)
// specialOfficer.VoteMethodPlayer, _ = item["vote_method_player"].(int64)

// specialOfficer.Name, _ = item["name"].(string)

// log.Println("查詢redis中特殊人員(包含測試資料): ", len(datas))

// for userID, data := range datas {
// 	var gameSpecialOfficerModel GameVoteSpecialOfficerModel
// 	// 解碼，取得特殊人員資料
// 	json.Unmarshal([]byte(data), &gameSpecialOfficerModel)

// 	if userID != "" && userID != "test" {
// 		// 過濾測試資料
// 		officers = append(officers, gameSpecialOfficerModel)
// 	}
// }

// log.Println("查詢redis中特殊�
