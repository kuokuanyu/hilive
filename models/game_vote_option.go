package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GameVoteOptionModel 資料表欄位
type GameVoteOptionModel struct {
	Base            `json:"-"`
	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id" example:"activity_id"`
	GameID          string `json:"game_id" example:"game_id"`
	OptionID        string `json:"option_id" example:"option_id"`
	OptionName      string `json:"option_name" example:"option_name"`
	OptionPicture   string `json:"option_picture" example:"option_picture"`
	OptionIntroduce string `json:"option_introduce" example:"option_introduce"`
	OptionScore     int64  `json:"option_score" example:"1"`

	// 活動
	UserID string `json:"user_id" example:"user_id"`
}

// EditGameVoteOptionModel 資料表欄位
type EditGameVoteOptionModel struct {
	ActivityID      string `json:"activity_id" example:"activity_id"`
	GameID          string `json:"game_id" example:"game_id"`
	OptionID        string `json:"option_id" example:"option_id"`
	OptionName      string `json:"option_name" example:"option_name"`
	OptionPicture   string `json:"option_picture" example:"option_picture"`
	OptionIntroduce string `json:"option_introduce" example:"option_introduce"`
	OptionScore     string `json:"option_score" example:"1"`
}

// DefaultGameVoteOptionModel 預設GameVoteOptionModel
func DefaultGameVoteOptionModel() GameVoteOptionModel {
	return GameVoteOptionModel{Base: Base{TableName: config.ACTIVITY_GAME_VOTE_OPTION_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameVoteOptionModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameVoteOptionModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameVoteOptionModel) SetDbConn(conn db.Connection) GameVoteOptionModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameVoteOptionModel) SetRedisConn(conn cache.Connection) GameVoteOptionModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (a GameVoteOptionModel) SetMongoConn(conn mongo.Connection) GameVoteOptionModel {
// 	a.MongoConn = conn
// 	return a
// }

// IncrOptionScore 遞增投票選項分數資料
func (a GameVoteOptionModel) IncrOptionScore(gameID string, optionID string, score int64) error {
	// 從redis取得資料，確定redis中有該場遊戲的投票選項資料，避免分數計算錯誤
	a.FindOrderByScore(true, gameID)

	// ### redis鎖處理
	for l := 0; l < MaxRetries; l++ {
		// 上鎖
		ok, _ := acquireLock(a.RedisConn, config.VOTE_OPTION_LOCK_REDIS, LockExpiration)
		if ok == "OK" {
			// 加入funciton(如果有提早return一定要釋放鎖)

			// 遞增資料表選項分數資料
			err := a.Table(a.Base.TableName).
				Where("game_id", "=", gameID).
				Where("option_id", "=", optionID).
				Update(command.Value{"option_score": fmt.Sprintf("`option_score` + %s", strconv.Itoa(int(score)))})
			if err != nil {
				releaseLock(a.RedisConn, config.VOTE_OPTION_LOCK_REDIS)
				return errors.New("錯誤: 遞增投票選項分數發生問題(資料庫)")
			}

			// 遞增該選項分數資料(redis)
			err = a.RedisConn.ZSetIncrCache(config.SCORES_REDIS+gameID, optionID, score)
			if err != nil {
				releaseLock(a.RedisConn, config.VOTE_OPTION_LOCK_REDIS)
				return errors.New("錯誤: 遞增投票選項分數發生問題(redis)")
			}

			// 手動釋放鎖
			releaseLock(a.RedisConn, config.VOTE_OPTION_LOCK_REDIS)

			break
		}

		// 鎖被佔用，稍微延遲後重試
		time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
	}

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

	return nil
}

// FindOrderByScore 查詢遊戲選項資料(分數排序，從redis取資料，如果沒資料的話才從資料庫取資料)
func (a GameVoteOptionModel) FindOrderByScore(isRedis bool, gameID string) ([]GameVoteOptionModel, error) {
	// log.Println("執行FindOrderByScore")
	var (
		optionModels = make([]GameVoteOptionModel, 0)
		// err    error
	)

	if isRedis {
		// 分數從高至低(回傳資料回option_id陣列)
		options, err := a.RedisConn.ZSetRevRange(config.SCORES_REDIS+gameID, 0, 0) // 包含test測試資料，避免重複執行資料庫
		if err != nil {
			return optionModels, errors.New("錯誤: 從redis中取得分數由高至低的投票選項資訊發生問題")
		}

		// log.Println("查詢redis中排名紀錄(包含測試資料): ", len(options))

		// 處理分數排名資訊
		for _, optionID := range options {
			// 過濾測試資料
			if optionID != "test" {
				score := a.RedisConn.ZSetIntScore(config.SCORES_REDIS+gameID, optionID) // 分數資料

				optionModels = append(optionModels, GameVoteOptionModel{
					OptionID:    optionID,
					OptionScore: score,
				})
			}
		}

		// log.Println("查詢redis後的排名紀錄(不包含測試資料): ", len(optionModels))

		if len(options) == 0 {
			// log.Println("redis裡連測試資料都沒有，查詢資料表: ")

			// 從資料表中取得資料
			items, err := a.Table(a.Base.TableName).
				Where("game_id", "=", gameID).
				OrderBy("option_score", "desc").All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得投票選項資訊(mysql)，請重新查詢")
			}

			optionModels = MapToVoteOptionModel(items)

			// 將資料寫入redis中
			params := []interface{}{config.SCORES_REDIS + gameID} // redis參數(票數排序)

			for _, option := range optionModels {
				params = append(params, option.OptionScore, option.OptionID)
			}

			if len(optionModels) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				params = append(params, 0, "test")
			}

			// 票數排序
			if err = a.RedisConn.ZSetMultiAdd(params); err != nil {
				return nil, errors.New("錯誤: 將投票選項資訊寫入redis中發生問題")
			}

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
		}
	}

	// else { // 資料庫查詢資料
	// 	items, err := a.Table(a.Base.TableName).
	// 		Where("game_id", "=", gameID).
	// 		OrderBy("option_score", "desc").All()
	// 	if err != nil {
	// 		return nil, errors.New("錯誤: 無法取得投票選項資訊(mysql)，請重新查詢")
	// 	}

	// 	optionModels = MapToVoteOptionModel(items)
	// }

	return optionModels, nil
}

// FindOrderByID 查詢遊戲選項資料(id排序，從資料庫取資料後將資料寫入redis)
func (a GameVoteOptionModel) FindOrderByID(isRedis bool, gameID string) ([]GameVoteOptionModel, error) {
	// log.Println("執行FindOrderByID")
	var (
		optionModels = make([]GameVoteOptionModel, 0)
		// datas        map[string]string
		// err          error
	)

	// if isRedis {
	// }

	// if len(datas) == 0 {
	// log.Println("redis裡連測試資料都沒有，查詢資料表: ")

	// 資料庫取得資料
	items, err := a.Table(a.Base.TableName).
		Select(
			"activity_game_vote_option.id",
			"activity_game_vote_option.activity_id",
			"activity_game_vote_option.game_id",
			"activity_game_vote_option.option_id",
			"activity_game_vote_option.option_name",
			"activity_game_vote_option.option_picture",
			"activity_game_vote_option.option_introduce",
			"activity_game_vote_option.option_score",

			// 活動
			"activity.user_id",
		).
		LeftJoin(command.Join{
			FieldA:    "activity_game_vote_option.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity",
			Operation: "="}).
		Where("game_id", "=", gameID).
		OrderBy("id", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得投票選項資訊，請重新查詢")
	}
	optionModels = MapToVoteOptionModel(items)

	// 判斷角色時需執行，避免過度寫入redis
	// 將資料寫入redis中
	// params := []interface{}{config.SCORES_REDIS + gameID} // redis參數(票數排序)

	// for _, option := range optionModels {
	// 	params = append(params, option.OptionScore, option.OptionID)
	// }

	// if len(optionModels) == 0 {
	// 	// log.Println("資料表一樣沒資料，加入測試資料")

	// 	// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
	// 將測試資料寫入redis，避免重複查詢資料庫
	// 	params = append(params, 0, "test")
	// }

	// // 票數排序
	// if err = a.RedisConn.ZSetMultiAdd(params); err != nil {
	// 	return nil, errors.New("錯誤: 將投票選項資訊寫入redis中發生問題")
	// }

	// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	// a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

	return optionModels, nil
}

// MapToVoteOptionModel map轉換[]GameVoteOptionModel
func MapToVoteOptionModel(items []map[string]interface{}) []GameVoteOptionModel {
	var options = make([]GameVoteOptionModel, 0)
	for _, item := range items {
		var (
			option GameVoteOptionModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &option)

		if !strings.Contains(option.OptionPicture, "system") {
			option.OptionPicture = "/admin/uploads/" + option.UserID + "/" + option.ActivityID + "/interact/sign/" + "vote" + "/" + option.GameID + "/" + option.OptionPicture
		}

		options = append(options, option)
	}
	return options
}

// Add 新增投票選項資料(加入資料庫)
func (a GameVoteOptionModel) Add(model EditGameVoteOptionModel) error {
	var (
		optionID    = utils.UUID(20)
		fieldValues = command.Value{
			"activity_id":      model.ActivityID,
			"game_id":          model.GameID,
			"option_id":        optionID,
			"option_name":      model.OptionName,
			"option_picture":   model.OptionPicture,
			"option_introduce": model.OptionIntroduce,
			"option_score":     0,
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

	// 將資料寫入資料庫
	if _, err := a.Table(a.TableName).Insert(fieldValues); err != nil {
		return errors.New("錯誤: 新增投票選項資料發生問題")
	}

	// 從redis取得資料，確定redis中有該場遊戲的投票選項資料
	a.FindOrderByScore(true, model.GameID)

	// 將新增的投票選項資料寫入redis中
	a.RedisConn.ZSetAddInt(config.SCORES_REDIS+model.GameID, optionID, 0)

	// 編輯次數更新
	a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

	return nil
}

// Update 更新選項資料
func (a GameVoteOptionModel) Update(model EditGameVoteOptionModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"option_name", "option_picture", "option_introduce"}
		values      = []string{model.OptionName, model.OptionPicture, model.OptionIntroduce}
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

	// 從redis取得資料，確定redis中有該場遊戲的投票選項資料
	a.FindOrderByScore(true, model.GameID)

	err := a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).
		Where("game_id", "=", model.GameID).
		Where("option_id", "=", model.OptionID).
		Update(fieldValues)
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return err
	}

	// 編輯次數更新
	a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

	return nil
}

// Reset 重新計算所有選項投票紀錄
func (a GameVoteOptionModel) Reset(activityID, gameID string) error {
	// 取得所有票數分數資料
	options, err := a.FindOrderByID(false, gameID)
	if err != nil {
		return err
	}

	if len(options) > 0 {
		// 取得遊戲輪次紀錄
		gameModel, err := DefaultGameModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			Find(false, gameID, "vote")
		if err != nil {
			return err
		}

		// 重新計算所有選項分數資料
		for _, option := range options {
			// 選項分數資料歸零
			var score int64
			// options[i].OptionScore = 0

			// 取得該選項投票紀錄
			records, err := DefaultGameVoteRecordModel().
				SetConn(a.DbConn, a.RedisConn, a.MongoConn).
				FindAll(activityID, gameID, option.OptionID, strconv.Itoa(int(gameModel.GameRound)), 0, 0)
			if err != nil {
				return err
			}

			for _, record := range records {
				// 將投票紀錄分數加入選項資料中
				score += record.Score
			}

			// 將重新計算後的分數資料更新至資料表中
			err = a.Table(a.Base.TableName).
				Where("activity_id", "=", activityID).
				Where("game_id", "=", gameID).
				Where("option_id", "=", option.OptionID).
				Update(command.Value{"option_score": strconv.Itoa(int(score))})
			if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
		}

		// 清除舊的redis資訊
		a.RedisConn.DelCache(config.SCORES_REDIS + gameID)

		// 將新的資料重新寫入redis
		a.FindOrderByScore(true, gameID)
	}
	return nil
}

// option.ID, _ = item["id"].(int64)
// option.ActivityID, _ = item["activity_id"].(string)
// option.UserID, _ = item["user_id"].(string)
// option.GameID, _ = item["game_id"].(string)
// option.OptionID, _ = item["option_id"].(string)
// option.OptionName, _ = item["option_name"].(string)
// option.OptionPicture, _ = item["option_picture"].(string)
// option.OptionIntroduce, _ = item["option_introduce"].(string)
// option.OptionScore, _ = item["option_score"].(int64)

// 判斷redis裡是否有投票選項資訊，有則不執行查詢資料表功能
// datas, err = a.RedisConn.HashGetAllCache(config.VOTE_OPTION_REDIS + gameID) // 包含test測試資料，避免重複執行資料庫
// if err != nil {
// 	return optionModels, errors.New("錯誤: 取得投票選項快取資料發生問題")
// }

// log.Println("查詢redis中投票選項資料(包含測試資料): ", len(datas))

// for optionID, data := range datas {
// 	var (
// 		optionModel GameVoteOptionModel
// 	)
// 	// 解碼，取得投票選項資訊
// 	json.Unmarshal([]byte(data), &optionModel)

// 	if optionID != "test" {
// 		// 過濾測試資料
// 		optionModels = append(optionModels, optionModel)
// 	}
// }

// log.Println("查詢redis後投票選項資料(不包含測試資料): ", len(optionModels))

// 將資料寫入redis中
// params := []interface{}{config.VOTE_OPTION_REDIS + gameID} // redis參數(投票選項資訊)

// for _, option := range optionModels {
// 	if option.OptionID != "" {
// 		// 加密，寫入redis中
// 		params = append(params, option.OptionID, utils.JSON(option))
// 	}
// }

// if len(optionModels) == 0 {
// 	log.Println("資料表一樣沒資料，加入測試資料")

// 	// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
// 	params = append(params, "test", utils.JSON(GameVoteOptionModel{}))
// }

// // 投票選項資訊
// if err = a.RedisConn.HashMultiSetCache(params); err != nil {
// 	return nil, errors.New("錯誤: 將投票選項資訊寫入redis中發生問題")
// }

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資
