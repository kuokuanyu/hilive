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
	"strconv"
	"strings"
)

// GameQARecordModel 資料表欄位
type GameQARecordModel struct {
	Base `json:"-"`
}

// EditGameQARecordModel 資料表欄位
type EditGameQARecordModel struct {
	UserID      string `json:"user_id" example:"user_id"`
	GameID      string `json:"game_id" example:"game_id"`
	QARound     int64  `json:"qa_round" example:"1"`       //遊戲題數
	QAOption    string `json:"qa_option" example:"A"`      //遊戲選項
	Score       int64  `json:"score" example:"100"`        //遊戲分數
	OriginScore int64  `json:"origin_score" example:"100"` // 原始分數
	AddScore    int64  `json:"add_score" example:"100"`    // 加成分數
}

// DefaultGameQARecordModel 預設GameQARecordModel
func DefaultGameQARecordModel() GameQARecordModel {
	return GameQARecordModel{Base: Base{TableName: config.ACTIVITY_GAME_QA_RECORD_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameQARecordModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameQARecordModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameQARecordModel) SetDbConn(conn db.Connection) GameQARecordModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameQARecordModel) SetRedisConn(conn cache.Connection) GameQARecordModel {
// 	a.RedisConn = conn
// 	return a
// }

// Find 查詢用戶答題紀錄，從redis中查詢
func (a GameQARecordModel) Find(isRedis bool, gameID, userID string) (map[string]interface{}, error) {
	var record = make(map[string]interface{}, 0)
	if isRedis {
		data, err := a.RedisConn.HashGetCache(config.QA_RECORD_REDIS+gameID, userID)
		if err != nil {
			return record, errors.New("錯誤: 取得用戶答題紀錄快取資料發生問題")
		}

		// 解碼
		json.Unmarshal([]byte(data), &record)
	}
	return record, nil
}

// FindAllByRedis 查詢所有用戶答題紀錄，從redis中查詢
func (a GameQARecordModel) FindAllByRedis(isRedis bool,
	gameID string) (map[string]map[string]interface{}, error) {
	var records = make(map[string]map[string]interface{}, 0)
	if isRedis {
		dataMap, err := a.RedisConn.HashGetAllCache(config.QA_RECORD_REDIS + gameID)
		if err != nil {
			return records, errors.New("錯誤: 取得所有用戶答題紀錄快取資料發生問題")
		}

		for userID, data := range dataMap {
			var record map[string]interface{}

			// 解碼
			json.Unmarshal([]byte(data), &record)

			records[userID] = record
		}
	}
	return records, nil
}

// FindAllByDB 查詢所有用戶答題紀錄，從資料表中查詢
func (a GameQARecordModel) FindAllByDB(gameID string, round int64) ([]map[string]interface{}, error) {
	var (
		sql = a.Table(a.Base.TableName).
			LeftJoin(command.Join{
				FieldA:    "activity_game_qa_record.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_game_qa_record.game_id",
				FieldA1:   "activity_game_qa.game_id",
				Table:     "activity_game_qa",
				Operation: "="}).
			Where("activity_game_qa_record.game_id", "=", gameID).
			OrderBy("activity_game_qa_record.round", "asc",
				"activity_game_qa_record.score", "desc")
		items = make([]map[string]interface{}, 0)
		err   error
	)

	if round == 0 {
		// 不分輪次
	} else {
		// 特定輪次
		sql = sql.Where("activity_game_qa_record.round", "=", round)
	}

	// 查詢資料表所有答題紀錄資料
	items, err = sql.All()
	if err != nil {
		return items, errors.New("錯誤: 無法取得所有用戶答題紀錄，請重新查詢")
	}
	// fmt.Println("items:", items)
	return items, nil
}

// Adds 批量新增用戶答題紀錄資料(加入資料庫)
func (a GameQARecordModel) Adds(dataAmount int, activityID, gameID string, round int,
	records map[string]map[string]interface{}) error {
	// 判斷資料數
	if dataAmount > 0 {
		var (
			activityIDs  = make([]string, dataAmount)   // 活動ID
			gameIDs      = make([]string, dataAmount)   // 活動ID
			userIDs      = make([]string, dataAmount)   // 用戶ID
			rounds       = make([]string, dataAmount)   // 輪次
			roundStr     = strconv.Itoa(round)          // 輪次(string)
			scores       = make([]string, dataAmount)   // 分數
			recordParams = make(map[string][]string, 0) // 批量新增須要使用參數(key為資料表欄位名稱)
			keys         = []string{
				"qa_1_option", "qa_2_option", "qa_3_option", "qa_4_option", "qa_5_option",
				"qa_6_option", "qa_7_option", "qa_8_option", "qa_9_option", "qa_10_option",
				"qa_11_option", "qa_2_option", "qa_13_option", "qa_14_option", "qa_15_option",
				"qa_16_option", "qa_17_option", "qa_18_option", "qa_19_option", "qa_20_option",

				"qa_1_origin_score", "qa_2_origin_score", "qa_3_origin_score", "qa_4_origin_score", "qa_5_origin_score",
				"qa_6_origin_score", "qa_7_origin_score", "qa_8_origin_score", "qa_9_origin_score", "qa_10_origin_score",
				"qa_11_origin_score", "qa_2_origin_score", "qa_13_origin_score", "qa_14_origin_score", "qa_15_origin_score",
				"qa_16_origin_score", "qa_17_origin_score", "qa_18_origin_score", "qa_19_origin_score", "qa_20_origin_score",

				"qa_1_add_score", "qa_2_add_score", "qa_3_add_score", "qa_4_add_score", "qa_5_add_score",
				"qa_6_add_score", "qa_7_add_score", "qa_8_add_score", "qa_9_add_score", "qa_10_add_score",
				"qa_11_add_score", "qa_2_add_score", "qa_13_add_score", "qa_14_add_score", "qa_15_add_score",
				"qa_16_add_score", "qa_17_add_score", "qa_18_add_score", "qa_19_add_score", "qa_20_add_score",
			}
		)

		// 將recordParams寫入預設長度陣列資料
		for _, key := range keys {
			recordParams[key] = make([]string, dataAmount)
		}

		// 處理活動.遊戲ID.輪次陣列資料(固定資料)
		for i := range dataAmount {
			activityIDs[i] = activityID
			gameIDs[i] = gameID
			rounds[i] = roundStr
		}

		// 處理答題記錄
		var count int64 // 答題記錄index資料
		for userID, record := range records {
			userIDs[count] = userID // 用戶ID

			for key, value := range record {
				if key == "score" {
					scores[count] = utils.GetString(value, "0") // 分數
				} else {
					recordParams[key][count] = utils.GetString(value, "") // 其他答題記錄欄位
				}
			}

			count++
		}

		// 處理recordParams參數，將origin_score、add_score相關欄位資料預設值設為0
		for key, values := range recordParams {
			if strings.Contains(key, "origin_score") || strings.Contains(key, "add_score") {
				for i, value := range values {
					if value == "" {
						// 如果資料為空字串，則改為"0"
						values[i] = "0"
					}
				}
			}
		}

		// 將資料匹量寫入activity_game_qa_record表中
		err := a.Table(a.TableName).BatchInsert(
			dataAmount,
			"activity_id,game_id,user_id,round,score,qa_1_option,qa_2_option,qa_3_option,qa_4_option,qa_5_option,qa_6_option,qa_7_option,qa_8_option,qa_9_option,qa_10_option,qa_11_option,qa_12_option,qa_13_option,qa_14_option,qa_15_option,qa_16_option,qa_17_option,qa_18_option,qa_19_option,qa_20_option,qa_1_origin_score,qa_2_origin_score,qa_3_origin_score,qa_4_origin_score,qa_5_origin_score,qa_6_origin_score,qa_7_origin_score,qa_8_origin_score,qa_9_origin_score,qa_10_origin_score,qa_11_origin_score,qa_12_origin_score,qa_13_origin_score,qa_14_origin_score,qa_15_origin_score,qa_16_origin_score,qa_17_origin_score,qa_18_origin_score,qa_19_origin_score,qa_20_origin_score,qa_1_add_score,qa_2_add_score,qa_3_add_score,qa_4_add_score,qa_5_add_score,qa_6_add_score,qa_7_add_score,qa_8_add_score,qa_9_add_score,qa_10_add_score,qa_11_add_score,qa_12_add_score,qa_13_add_score,qa_14_add_score,qa_15_add_score,qa_16_add_score,qa_17_add_score,qa_18_add_score,qa_19_add_score,qa_20_add_score",
			[][]string{activityIDs, gameIDs, userIDs, rounds, scores,
				recordParams["qa_1_option"], recordParams["qa_2_option"], recordParams["qa_3_option"], recordParams["qa_4_option"], recordParams["qa_5_option"],
				recordParams["qa_6_option"], recordParams["qa_7_option"], recordParams["qa_8_option"], recordParams["qa_9_option"], recordParams["qa_10_option"],
				recordParams["qa_11_option"], recordParams["qa_2_option"], recordParams["qa_13_option"], recordParams["qa_14_option"], recordParams["qa_15_option"],
				recordParams["qa_16_option"], recordParams["qa_17_option"], recordParams["qa_18_option"], recordParams["qa_19_option"], recordParams["qa_20_option"],

				recordParams["qa_1_origin_score"], recordParams["qa_2_origin_score"], recordParams["qa_3_origin_score"], recordParams["qa_4_origin_score"], recordParams["qa_5_origin_score"],
				recordParams["qa_6_origin_score"], recordParams["qa_7_origin_score"], recordParams["qa_8_origin_score"], recordParams["qa_9_origin_score"], recordParams["qa_10_origin_score"],
				recordParams["qa_11_origin_score"], recordParams["qa_2_origin_score"], recordParams["qa_13_origin_score"], recordParams["qa_14_origin_score"], recordParams["qa_15_origin_score"],
				recordParams["qa_16_origin_score"], recordParams["qa_17_origin_score"], recordParams["qa_18_origin_score"], recordParams["qa_19_origin_score"], recordParams["qa_20_origin_score"],

				recordParams["qa_1_add_score"], recordParams["qa_2_add_score"], recordParams["qa_3_add_score"], recordParams["qa_4_add_score"], recordParams["qa_5_add_score"],
				recordParams["qa_6_add_score"], recordParams["qa_7_add_score"], recordParams["qa_8_add_score"], recordParams["qa_9_add_score"], recordParams["qa_10_add_score"],
				recordParams["qa_11_add_score"], recordParams["qa_2_add_score"], recordParams["qa_13_add_score"], recordParams["qa_14_add_score"], recordParams["qa_15_add_score"],
				recordParams["qa_16_add_score"], recordParams["qa_17_add_score"], recordParams["qa_18_add_score"], recordParams["qa_19_add_score"], recordParams["qa_20_add_score"]})
		if err != nil {
			return errors.New("錯誤: 匹量新增玩家答題記錄發生問題(activity_game_qa_record)，請重新操作")
		}
	}

	return nil
}

// Add 新增用戶答題紀錄資料(加入資料庫)
func (a GameQARecordModel) Add(userID, activityID, gameID string, round int64,
	record map[string]interface{}) error {
	var (
		fieldValues = command.Value{
			"user_id":     userID,
			"activity_id": activityID,
			"game_id":     gameID,
			"round":       round,
		}
	)
	for key, value := range record {
		fieldValues[key] = value
	}

	if _, err := a.Table(a.TableName).Insert(fieldValues); err != nil {
		return errors.New("錯誤: 新增用戶答題紀錄發生問題")
	}

	return nil
}

// Update 更新用戶答題紀錄，更新redis中資料
func (a GameQARecordModel) Update(isRedis bool, model EditGameQARecordModel) error {
	if isRedis {
		var (
			record, _ = DefaultGameQARecordModel().
				SetConn(a.DbConn, a.RedisConn, a.MongoConn).
				Find(true, model.GameID, model.UserID) // 用戶答題紀錄
			// fieldValues = command.Value{"score": model.Score}
			// recordJson = utils.JSON(record) // 編碼用戶答題紀錄
			// mapRecord  = make(map[string]interface{}, 0)
		)
		// json.Unmarshal([]byte(recordJson), &mapRecord)

		// 更新用戶分數、答題紀錄
		record["score"] = model.Score
		record[fmt.Sprintf("qa_%d_option", model.QARound)] = model.QAOption
		record[fmt.Sprintf("qa_%d_origin_score", model.QARound)] = model.OriginScore
		record[fmt.Sprintf("qa_%d_add_score", model.QARound)] = model.AddScore

		// 將新的答題紀錄JSON編碼
		// recordJson := utils.JSON(mapRecord)

		// 更新答題紀錄redis資訊
		a.RedisConn.HashSetCache(config.QA_RECORD_REDIS+model.GameID,
			model.UserID, utils.JSON(record))
	}
	return nil
}

// UserID     string `json:"user_id"`
// ActivityID string `json:"activity_id" example:"activity_id"`
// GameID     string `json:"game_id" example:"game_id"`
// Round      int64  `json:"round" example:"1"`
// Score int64 `json:"score" example:"1"`

// 答題紀錄
// QA1Option       string `json:"qa_1_option"`
// QA1OriginScore  int64  `json:"qa_1_origin_score"`
// QA1AddScore     int64  `json:"qa_1_add_score"`
// QA2Option       string `json:"qa_2_option"`
// QA2OriginScore  int64  `json:"qa_2_origin_score"`
// QA2AddScore     int64  `json:"qa_2_add_score"`
// QA3Option       string `json:"qa_3_option"`
// QA3OriginScore  int64  `json:"qa_3_origin_score"`
// QA3AddScore     int64  `json:"qa_3_add_score"`
// QA4Option       string `json:"qa_4_option"`
// QA4OriginScore  int64  `json:"qa_4_origin_score"`
// QA4AddScore     int64  `json:"qa_4_add_score"`
// QA5Option       string `json:"qa_5_option"`
// QA5OriginScore  int64  `json:"qa_5_origin_score"`
// QA5AddScore     int64  `json:"qa_5_add_score"`
// QA6Option       string `json:"qa_6_option"`
// QA6OriginScore  int64  `json:"qa_6_origin_score"`
// QA6AddScore     int64  `json:"qa_6_add_score"`
// QA7Option       string `json:"qa_7_option"`
// QA7OriginScore  int64  `json:"qa_7_origin_score"`
// QA7AddScore     int64  `json:"qa_7_add_score"`
// QA8Option       string `json:"qa_8_option"`
// QA8OriginScore  int64  `json:"qa_8_origin_score"`
// QA8AddScore     int64  `json:"qa_8_add_score"`
// QA9Option       string `json:"qa_9_option"`
// QA9OriginScore  int64  `json:"qa_9_origin_score"`
// QA9AddScore     int64  `json:"qa_9_add_score"`
// QA10Option      string `json:"qa_10_option"`
// QA10OriginScore int64  `json:"qa_10_origin_score"`
// QA10AddScore    int64  `json:"qa_10_add_score"`
// QA11Option      string `json:"qa_11_option"`
// QA11OriginScore int64  `json:"qa_11_origin_score"`
// QA11AddScore    int64  `json:"qa_11_add_score"`
// QA12Option      string `json:"qa_12_option"`
// QA12OriginScore int64  `json:"qa_12_origin_score"`
// QA12AddScore    int64  `json:"qa_12_add_score"`
// QA13Option      string `json:"qa_13_option"`
// QA13OriginScore int64  `json:"qa_13_origin_score"`
// QA13AddScore    int64  `json:"qa_13_add_score"`
// QA14Option      string `json:"qa_14_option"`
// QA14OriginScore int64  `json:"qa_14_origin_score"`
// QA14AddScore    int64  `json:"qa_14_add_score"`
// QA15Option      string `json:"qa_15_option"`
// QA15OriginScore int64  `json:"qa_15_origin_score"`
// QA15AddScore    int64  `json:"qa_15_add_score"`
// QA16Option      string `json:"qa_16_option"`
// QA16OriginScore int64  `json:"qa_16_origin_score"`
// QA16AddScore    int64  `json:"qa_16_add_score"`
// QA17Option      string `json:"qa_17_option"`
// QA17OriginScore int64  `json:"qa_17_origin_score"`
// QA17AddScore    int64  `json:"qa_17_add_score"`
// QA18Option      string `json:"qa_18_option"`
// QA18OriginScore int64  `json:"qa_18_origin_score"`
// QA18AddScore    int64  `json:"qa_18_add_score"`
// QA19Option      string `json:"qa_19_option"`
// QA19OriginScore int64  `json:"qa_19_origin_score"`
// QA19AddScore    int64  `json:"qa_19_add_score"`
// QA20Option      string `json:"qa_20_option"`
// QA20OriginScore int64  `json:"qa_20_origin_score"`
// QA20AddScore    int64  `json:"qa_20_add_score"`

// command.Value{
// 	"user_id":     userID,
// 	"activity_id": activityID,
// 	"game_id":     gameID,
// 	"round":       round,
// 	"score":       model.Score,

// 	"qa_1_option":       model.QA1Option,
// 	"qa_1_origin_score": model.QA1OriginScore,
// 	"qa_1_add_score":    model.QA1AddScore,

// 	"qa_2_option":       model.QA2Option,
// 	"qa_2_origin_score": model.QA2OriginScore,
// 	"qa_2_add_score":    model.QA2AddScore,

// 	"qa_3_option":       model.QA3Option,
// 	"qa_3_origin_score": model.QA3OriginScore,
// 	"qa_3_add_score":    model.QA3AddScore,

// 	"qa_4_option":       model.QA4Option,
// 	"qa_4_origin_score": model.QA4OriginScore,
// 	"qa_4_add_score":    model.QA4AddScore,

// 	"qa_5_option":       model.QA5Option,
// 	"qa_5_origin_score": model.QA5OriginScore,
// 	"qa_5_add_score":    model.QA5AddScore,

// 	"qa_6_option":       model.QA6Option,
// 	"qa_6_origin_score": model.QA6OriginScore,
// 	"qa_6_add_score":    model.QA6AddScore,

// 	"qa_7_option":       model.QA7Option,
// 	"qa_7_origin_score": model.QA7OriginScore,
// 	"qa_7_add_score":    model.QA7AddScore,

// 	"qa_8_option":       model.QA8Option,
// 	"qa_8_origin_score": model.QA8OriginScore,
// 	"qa_8_add_score":    model.QA8AddScore,

// 	"qa_9_option":       model.QA9Option,
// 	"qa_9_origin_score": model.QA9OriginScore,
// 	"qa_9_add_score":    model.QA9AddScore,

// 	"qa_10_option":       model.QA10Option,
// 	"qa_10_origin_score": model.QA10OriginScore,
// 	"qa_10_add_score":    model.QA10AddScore,

// 	"qa_11_option":       model.QA11Option,
// 	"qa_11_origin_score": model.QA11OriginScore,
// 	"qa_11_add_score":    model.QA11AddScore,

// 	"qa_12_option":       model.QA12Option,
// 	"qa_12_origin_score": model.QA12OriginScore,
// 	"qa_12_add_score":    model.QA12AddScore,

// 	"qa_13_option":       model.QA13Option,
// 	"qa_13_origin_score": model.QA13OriginScore,
// 	"qa_13_add_score":    model.QA13AddScore,

// 	"qa_14_option":       model.QA14Option,
// 	"qa_14_origin_score": model.QA14OriginScore,
// 	"qa_14_add_score":    model.QA14AddScore,

// 	"qa_15_option":       model.QA15Option,
// 	"qa_15_origin_score": model.QA15OriginScore,
// 	"qa_15_add_score":    model.QA15AddScore,

// 	"qa_16_option":       model.QA16Option,
// 	"qa_16_origin_score": model.QA16OriginScore,
// 	"qa_16_add_score":    model.QA16AddScore,

// 	"qa_17_option":       model.QA17Option,
// 	"qa_17_origin_score": model.QA17OriginScore,
// 	"qa_17_add_score":    model.QA17AddScore,

// 	"qa_18_option":       model.QA18Option,
// 	"qa_18_origin_score": model.QA18OriginScore,
// 	"qa_18_add_score":    model.QA18AddScore,

// 	"qa_19_option":       model.QA19Option,
// 	"qa_19_origin_score": model.QA19OriginScore,
// 	"qa_19_add_score":    model.QA19AddScore,

// 	"qa_20_option":       model.QA20Option,
// 	"qa_20_origin_score": model.QA20OriginScore,
// 	"qa_20_add_score":    model.QA20AddScore,
// }

// NewGameQARecordModel 資料表欄位
// type NewGameQARecordModel struct {
// 	UserID     string `json:"user_id"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	GameID     string `json:"game_id" example:"game_id"`
// 	Round      int64  `json:"round" example:"1"`
// 	Score      int64  `json:"score" example:"1"`

// 	// 答題紀錄
// 	QA1Option       string `json:"qa_1_option"`
// 	QA1OriginScore  int64  `json:"qa_1_origin_score"`
// 	QA1AddScore     int64  `json:"qa_1_add_score"`
// 	QA2Option       string `json:"qa_2_option"`
// 	QA2OriginScore  int64  `json:"qa_2_origin_score"`
// 	QA2AddScore     int64  `json:"qa_2_add_score"`
// 	QA3Option       string `json:"qa_3_option"`
// 	QA3OriginScore  int64  `json:"qa_3_origin_score"`
// 	QA3AddScore     int64  `json:"qa_3_add_score"`
// 	QA4Option       string `json:"qa_4_option"`
// 	QA4OriginScore  int64  `json:"qa_4_origin_score"`
// 	QA4AddScore     int64  `json:"qa_4_add_score"`
// 	QA5Option       string `json:"qa_5_option"`
// 	QA5OriginScore  int64  `json:"qa_5_origin_score"`
// 	QA5AddScore     int64  `json:"qa_5_add_score"`
// 	QA6Option       string `json:"qa_6_option"`
// 	QA6OriginScore  int64  `json:"qa_6_origin_score"`
// 	QA6AddScore     int64  `json:"qa_6_add_score"`
// 	QA7Option       string `json:"qa_7_option"`
// 	QA7OriginScore  int64  `json:"qa_7_origin_score"`
// 	QA7AddScore     int64  `json:"qa_7_add_score"`
// 	QA8Option       string `json:"qa_8_option"`
// 	QA8OriginScore  int64  `json:"qa_8_origin_score"`
// 	QA8AddScore     int64  `json:"qa_8_add_score"`
// 	QA9Option       string `json:"qa_9_option"`
// 	QA9OriginScore  int64  `json:"qa_9_origin_score"`
// 	QA9AddScore     int64  `json:"qa_9_add_score"`
// 	QA10Option      string `json:"qa_10_option"`
// 	QA10OriginScore int64  `json:"qa_10_origin_score"`
// 	QA10AddScore    int64  `json:"qa_10_add_score"`
// 	QA11Option      string `json:"qa_11_option"`
// 	QA11OriginScore int64  `json:"qa_11_origin_score"`
// 	QA11AddScore    int64  `json:"qa_11_add_score"`
// 	QA12Option      string `json:"qa_12_option"`
// 	QA12OriginScore int64  `json:"qa_12_origin_score"`
// 	QA12AddScore    int64  `json:"qa_12_add_score"`
// 	QA13Option      string `json:"qa_13_option"`
// 	QA13OriginScore int64  `json:"qa_13_origin_score"`
// 	QA13AddScore    int64  `json:"qa_13_add_score"`
// 	QA14Option      string `json:"qa_14_option"`
// 	QA14OriginScore int64  `json:"qa_14_origin_score"`
// 	QA14AddScore    int64  `json:"qa_14_add_score"`
// 	QA15Option      string `json:"qa_15_option"`
// 	QA15OriginScore int64  `json:"qa_15_origin_score"`
// 	QA15AddScore    int64  `json:"qa_15_add_score"`
// 	QA16Option      string `json:"qa_16_option"`
// 	QA16OriginScore int64  `json:"qa_16_origin_score"`
// 	QA16AddScore    int64  `json:"qa_16_add_score"`
// 	QA17Option      string `json:"qa_17_option"`
// 	QA17OriginScore int64  `json:"qa_17_origin_score"`
// 	QA17AddScore    int64  `json:"qa_17_add_score"`
// 	QA18Option      string `json:"qa_18_option"`
// 	QA18OriginScore int64  `json:"qa_18_origin_score"`
// 	QA18AddScore    int64  `json:"qa_18_add_score"`
// 	QA19Option      string `json:"qa_19_option"`
// 	QA19OriginScore int64  `json:"qa_19_origin_score"`
// 	QA19AddScore    int64  `json:"qa_19_add_score"`
// 	QA20Option      string `json:"qa_20_option"`
// 	QA20OriginScore int64  `json:"qa_20_origin_score"`
// 	QA20AddScore    int64  `json:"qa_20_add_score"`
// }

// // MapToModel 將值設置至GameQARecordModel
// func (a GameQARecordModel) MapToModel(m map[string]interface{}) GameQARecordModel {
// 	a.ID, _ = m["id"].(int64)
// 	// a.UserID, _ = m["user_id"].(string)
// 	a.ActivityID, _ = m["activity_id"].(string)
// 	a.GameID, _ = m["game_id"].(string)
// 	a.Game, _ = m["game"].(string)
// 	a.Title, _ = m["title"].(string)
// 	a.GameType, _ = m["game_type"].(string)
// 	a.LimitTime, _ = m["limit_time"].(string)
// 	a.Second, _ = m["second"].(int64)
// 	a.MaxPeople, _ = m["max_people"].(int64)
// 	a.People, _ = m["people"].(int64)
// 	a.MaxTimes, _ = m["max_times"].(int64)
// 	a.Allow, _ = m["allow"].(string)
// 	a.Percent, _ = m["percent"].(int64)
// 	a.FirstPrize, _ = m["first_prize"].(int64)
// 	a.SecondPrize, _ = m["second_prize"].(int64)
// 	a.ThirdPrize, _ = m["third_prize"].(int64)
// 	a.GeneralPrize, _ = m["general_prize"].(int64)
// 	a.Topic, _ = m["topic"].(string)
// 	a.Skin, _ = m["skin"].(string)
// 	a.DisplayName, _ = m["display_name"].(string)
// 	a.GameRound, _ = m["game_round"].(int64)
// 	a.GameSecond, _ = m["game_second"].(int64)
// 	a.GameStatus, _ = m["game_status"].(string)
// 	a.GameAttend, _ = m["game_attend"].(int64)

// 	// 用戶
// 	// a.MaxActivityPeople, _ = m["max_activity_people"].(int64)
// 	// a.MaxGamePeople, _ = m["max_game_people"].(int64)

// 	// 遊戲基本設置
// 	a.LotteryGameAllow, _ = m["lottery_game_allow"].(string)
// 	a.RedpackGameAllow, _ = m["redpack_game_allow"].(string)
// 	a.RopepackGameAllow, _ = m["ropepack_game_allow"].(string)
// 	a.WhackMoleGameAllow, _ = m["whack_mole_game_allow"].(string)
// 	a.MonopolyGameAllow, _ = m["monopoly_game_allow"].(string)
// 	a.QAGameAllow, _ = m["qa_game_allow"].(string)
// 	a.DrawNumbersGameAllow, _ = m["draw_numbers_game_allow"].(string)
// 	a.AllGameAllow, _ = m["all_game_allow"].(string)

// 	// 敲敲樂自定義圖片
// 	a.WhackMoleHostBackground, _ = m["whack_mole_host_background"].(string)
// 	a.WhackMoleGuestBackground, _ = m["whack_mole_guest_background"].(string)
// 	a.WhackMoleDollarRatPicture, _ = m["whack_mole_dollar_rat_picture"].(string)
// 	a.WhackMoleRedpackRatPicture, _ = m["whack_mole_redpack_rat_picture"].(string)
// 	a.WhackMoleBombPicture, _ = m["whack_mole_bomb_picture"].(string)
// 	a.WhackMoleRatHolePicture, _ = m["whack_mole_rat_hole_picture"].(string)
// 	a.WhackMoleRockPicture, _ = m["whack_mole_rock_picture"].(string)
// 	a.WhackMoleRankPicture, _ = m["whack_mole_rank_picture"].(string)
// 	a.WhackMoleRankBackground, _ = m["whack_mole_rank_background"].(string)

// 	// 搖號抽獎自定義圖片
// 	a.DrawNumbersBackground, _ = m["draw_numbers_background"].(string)
// 	a.DrawNumbersGiftInsidePicture, _ = m["draw_numbers_gift_inside_picture"].(string)
// 	a.DrawNumbersGiftOutsidePicture, _ = m["draw_numbers_gift_outside_picture"].(string)
// 	a.DrawNumbersWinnerBackground, _ = m["draw_numbers_winner_background"].(string)

// 	// 抓偽鈔自定義圖片
// 	a.MonopolyGuestBackground, _ = m["monopoly_guest_background"].(string)
// 	a.MonopolyHostRank, _ = m["monopoly_host_rank"].(string)
// 	a.MonopolyHostWinnerList, _ = m["monopoly_host_winner_list"].(string)
// 	a.MonopolyHostMascot, _ = m["monopoly_host_mascot"].(string)
// 	a.MonopolyPeople, _ = m["monopoly_people"].(string)

// 	// 自定義圖片參數(陣列)
// 	var pictures = make([]string, 0)
// 	if a.Game == "whack_mole" {
// 		// 敲敲樂
// 		pictures = append(pictures, a.WhackMoleHostBackground)
// 		pictures = append(pictures, a.WhackMoleGuestBackground)
// 		pictures = append(pictures, a.WhackMoleDollarRatPicture)
// 		pictures = append(pictures, a.WhackMoleRedpackRatPicture)
// 		pictures = append(pictures, a.WhackMoleBombPicture)
// 		pictures = append(pictures, a.WhackMoleRatHolePicture)
// 		pictures = append(pictures, a.WhackMoleRockPicture)
// 		pictures = append(pictures, a.WhackMoleRankPicture)
// 		pictures = append(pictures, a.WhackMoleRankBackground)
// 	} else if a.Game == "draw_numbers" {
// 		// 搖號抽獎
// 		pictures = append(pictures)
// 		pictures = append(pictures, a.DrawNumbersBackground)
// 		pictures = append(pictures, a.DrawNumbersGiftInsidePicture)
// 		pictures = append(pictures, a.DrawNumbersGiftOutsidePicture)
// 		pictures = append(pictures, a.DrawNumbersWinnerBackground)
// 	} else if a.Game == "monopoly" {
// 		// 抓偽鈔
// 		pictures = append(pictures)
// 		pictures = append(pictures, a.MonopolyGuestBackground)
// 		pictures = append(pictures, a.MonopolyHostRank)
// 		pictures = append(pictures, a.MonopolyHostWinnerList)
// 		pictures = append(pictures, a.MonopolyHostMascot)
// 		pictures = append(pictures, a.MonopolyPeople)
// 	}
// 	a.CustomizePictures = pictures

// 	// 快問快答
// 	a.TotalQA, _ = m["total_qa"].(int64)
// 	a.QASecond, _ = m["qa_second"].(int64)
// 	a.QARound, _ = m["qa_round"].(int64)

// 	var questions = make([]QuestionModel, a.TotalQA)
// 	for i := 0; i < int(a.TotalQA); i++ {
// 		questions[i].Question, _ = m["qa_"+strconv.Itoa(i+1)].(string)
// 		// questions[i].Picture, _ = m["qa_"+strconv.Itoa(i+1)+"_picture"].(string)
// 		options, _ := m["qa_"+strconv.Itoa(i+1)+"_options"].(string)
// 		questions[i].Options = strings.Split(options, "&&&")
// 		questions[i].Answer, _ = m["qa_"+strconv.Itoa(i+1)+"_answer"].(string)
// 		questions[i].Score, _ = m["qa_"+strconv.Itoa(i+1)+"_score"].(int64)
// 	}
// 	a.Questions = questions

// 	a.QA1, _ = m["qa_1"].(string)
// 	qa1, _ := m["qa_1_options"].(string)
// 	a.QA1Options = strings.Split(qa1, "&&&")
// 	a.QA1Answer, _ = m["qa_1_answer"].(string)
// 	a.QA1Score, _ = m["qa_1_score"].(int64)

// 	a.QA2, _ = m["qa_2"].(string)
// 	qa2, _ := m["qa_2_options"].(string)
// 	a.QA2Options = strings.Split(qa2, "&&&")
// 	a.QA2Answer, _ = m["qa_2_answer"].(string)
// 	a.QA2Score, _ = m["qa_1_score"].(int64)

// 	a.QA3, _ = m["qa_3"].(string)
// 	qa3, _ := m["qa_3_options"].(string)
// 	a.QA3Options = strings.Split(qa3, "&&&")
// 	a.QA3Answer, _ = m["qa_3_answer"].(string)
// 	a.QA3Score, _ = m["qa_1_score"].(int64)

// 	a.QA4, _ = m["qa_4"].(string)
// 	qa4, _ := m["qa_4_options"].(string)
// 	a.QA4Options = strings.Split(qa4, "&&&")
// 	a.QA4Answer, _ = m["qa_4_answer"].(string)
// 	a.QA4Score, _ = m["qa_4_score"].(int64)

// 	a.QA5, _ = m["qa_5"].(string)
// 	qa5, _ := m["qa_5_options"].(string)
// 	a.QA5Options = strings.Split(qa5, "&&&")
// 	a.QA5Answer, _ = m["qa_5_answer"].(string)
// 	a.QA5Score, _ = m["qa_5_score"].(int64)

// 	a.QA6, _ = m["qa_6"].(string)
// 	qa6, _ := m["qa_6_options"].(string)
// 	a.QA6Options = strings.Split(qa6, "&&&")
// 	a.QA6Answer, _ = m["qa_6_answer"].(string)
// 	a.QA6Score, _ = m["qa_6_score"].(int64)

// 	a.QA7, _ = m["qa_7"].(string)
// 	qa7, _ := m["qa_7_options"].(string)
// 	a.QA7Options = strings.Split(qa7, "&&&")
// 	a.QA7Answer, _ = m["qa_7_answer"].(string)
// 	a.QA7Score, _ = m["qa_7_score"].(int64)

// 	a.QA8, _ = m["qa_8"].(string)
// 	qa8, _ := m["qa_8_options"].(string)
// 	a.QA8Options = strings.Split(qa8, "&&&")
// 	a.QA8Answer, _ = m["qa_8_answer"].(string)
// 	a.QA8Score, _ = m["qa_8_score"].(int64)

// 	a.QA9, _ = m["qa_9"].(string)
// 	qa9, _ := m["qa_9_options"].(string)
// 	a.QA9Options = strings.Split(qa9, "&&&")
// 	a.QA9Answer, _ = m["qa_9_answer"].(string)
// 	a.QA9Score, _ = m["qa_9_score"].(int64)

// 	a.QA10, _ = m["qa_10"].(string)
// 	qa10, _ := m["qa_10_options"].(string)
// 	a.QA10Options = strings.Split(qa10, "&&&")
// 	a.QA10Answer, _ = m["qa_10_answer"].(string)
// 	a.QA10Score, _ = m["qa_10_score"].(int64)

// 	a.QA11, _ = m["qa_11"].(string)
// 	qa11, _ := m["qa_11_options"].(string)
// 	a.QA11Options = strings.Split(qa11, "&&&")
// 	a.QA11Answer, _ = m["qa_11_answer"].(string)
// 	a.QA11Score, _ = m["qa_11_score"].(int64)

// 	a.QA12, _ = m["qa_12"].(string)
// 	qa12, _ := m["qa_12_options"].(string)
// 	a.QA12Options = strings.Split(qa12, "&&&")
// 	a.QA12Answer, _ = m["qa_12_answer"].(string)
// 	a.QA12Score, _ = m["qa_12_score"].(int64)

// 	a.QA13, _ = m["qa_13"].(string)
// 	qa13, _ := m["qa_13_options"].(string)
// 	a.QA13Options = strings.Split(qa13, "&&&")
// 	a.QA13Answer, _ = m["qa_13_answer"].(string)
// 	a.QA13Score, _ = m["qa_13_score"].(int64)

// 	a.QA14, _ = m["qa_14"].(string)
// 	qa14, _ := m["qa_14_options"].(string)
// 	a.QA14Options = strings.Split(qa14, "&&&")
// 	a.QA14Answer, _ = m["qa_14_answer"].(string)
// 	a.QA14Score, _ = m["qa_14_score"].(int64)

// 	a.QA15, _ = m["qa_15"].(string)
// 	qa15, _ := m["qa_15_options"].(string)
// 	a.QA15Options = strings.Split(qa15, "&&&")
// 	a.QA15Answer, _ = m["qa_15_answer"].(string)
// 	a.QA15Score, _ = m["qa_15_score"].(int64)

// 	a.QA16, _ = m["qa_16"].(string)
// 	qa16, _ := m["qa_16_options"].(string)
// 	a.QA16Options = strings.Split(qa16, "&&&")
// 	a.QA16Answer, _ = m["qa_16_answer"].(string)
// 	a.QA16Score, _ = m["qa_16_score"].(int64)

// 	a.QA17, _ = m["qa_17"].(string)
// 	qa17, _ := m["qa_17_options"].(string)
// 	a.QA17Options = strings.Split(qa17, "&&&")
// 	a.QA17Answer, _ = m["qa_17_answer"].(string)
// 	a.QA17Score, _ = m["qa_17_score"].(int64)

// 	a.QA18, _ = m["qa_18"].(string)
// 	qa18, _ := m["qa_18_options"].(string)
// 	a.QA18Options = strings.Split(qa18, "&&&")
// 	a.QA18Answer, _ = m["qa_18_answer"].(string)
// 	a.QA18Score, _ = m["qa_18_score"].(int64)

// 	a.QA19, _ = m["qa_19"].(string)
// 	qa19, _ := m["qa_19_options"].(string)
// 	a.QA19Options = strings.Split(qa19, "&&&")
// 	a.QA19Answer, _ = m["qa_19_answer"].(string)
// 	a.QA19Score, _ = m["qa_19_score"].(int64)

// 	a.QA20, _ = m["qa_20"].(string)
// 	qa20, _ := m["qa_20_options"].(string)
// 	a.QA20Options = strings.Split(qa20, "&&&")
// 	a.QA20Answer, _ = m["qa_20_answer"].(string)
// 	a.QA20Score, _ = m["qa_20_score"].(int64)

// 	// 活動資訊(join activity)
// 	a.UserID, _ = m["user_id"].(string)
// 	return a
// }

// // MapToGameQARecordModel map轉換[]GameQARecordModel
// func MapToGameQARecordModel(items []map[string]interface{}) []GameQARecordModel {
// 	var games = make([]GameQARecordModel, 0)
// 	for _, item := range items {
// 		var (
// 			game GameQARecordModel
// 		)
// 		game.ID, _ = item["id"].(int64)
// 		// game.UserID, _ = item["user_id"].(string)
// 		game.ActivityID, _ = item["activity_id"].(string)
// 		game.GameID, _ = item["game_id"].(string)
// 		game.Game, _ = item["game"].(string)
// 		game.Title, _ = item["title"].(string)
// 		game.GameType, _ = item["game_type"].(string)
// 		game.LimitTime, _ = item["limit_time"].(string)
// 		game.Second, _ = item["second"].(int64)
// 		game.MaxPeople, _ = item["max_people"].(int64)
// 		game.People, _ = item["people"].(int64)
// 		game.MaxTimes, _ = item["max_times"].(int64)
// 		game.Allow, _ = item["allow"].(string)
// 		game.Percent, _ = item["percent"].(int64)
// 		game.FirstPrize, _ = item["first_prize"].(int64)
// 		game.SecondPrize, _ = item["second_prize"].(int64)
// 		game.ThirdPrize, _ = item["third_prize"].(int64)
// 		game.GeneralPrize, _ = item["general_prize"].(int64)
// 		game.Topic, _ = item["topic"].(string)
// 		game.Skin, _ = item["skin"].(string)
// 		game.DisplayName, _ = item["display_name"].(string)
// 		game.GameRound, _ = item["game_round"].(int64)
// 		game.GameSecond, _ = item["game_second"].(int64)
// 		game.GameStatus, _ = item["game_status"].(string)
// 		game.GameAttend, _ = item["game_attend"].(int64)

// 		// 用戶
// 		// game.MaxActivityPeople, _ = item["max_activity_people"].(int64)
// 		// game.MaxGamePeople, _ = item["max_game_people"].(int64)
// 		games = append(games, game)
// 	}
// 	return games
// }

// fields = []string{"title", "game_type", "limit_time",
// 	"max_times", "allow", "percent", "first_prize",
// 	"second_prize", "third_prize", "general_prize", "topic", "skin", "display_name",

// 	// 敲敲樂自定義
// 	"whack_mole_host_background", "whack_mole_guest_background", "whack_mole_dollar_rat_picture",
// 	"whack_mole_redpack_rat_picture", "whack_mole_bomb_picture", "whack_mole_rat_hole_picture",
// 	"whack_mole_rock_picture", "whack_mole_rank_picture", "whack_mole_rank_background",

// 	// 搖號抽獎自定義
// 	"draw_numbers_background", "draw_numbers_gift_inside_picture",
// 	"draw_numbers_gift_outside_picture", "draw_numbers_winner_background",

// 	// 抓偽鈔自定義
// 	"monopoly_guest_background", "monopoly_host_rank",
// 	"monopoly_host_winner_list", "monopoly_host_mascot",
// 	"monopoly_people"}
// options = []string{record.QA1Option, record.QA2Option, record.QA3Option,
// 	record.QA4Option, record.QA5Option, record.QA6Option, record.QA7Option,
// 	record.QA8Option, record.QA9Option, record.QA10Option, record.QA11Option,
// 	record.QA12Option, record.QA13Option, record.QA14Option, record.QA15Option,
// 	record.QA16Option, record.QA17Option, record.QA18Option, record.QA19Option,
// 	record.QA20Option} // 選項
// origins = []int64{record.QA1OriginScore, record.QA2OriginScore, record.QA3OriginScore,
// 	record.QA4OriginScore, record.QA5OriginScore, record.QA6OriginScore, record.QA7OriginScore,
// 	record.QA8OriginScore, record.QA9OriginScore, record.QA10OriginScore, record.QA11OriginScore,
// 	record.QA12OriginScore, record.QA13OriginScore, record.QA14OriginScore, record.QA15OriginScore,
// 	record.QA16OriginScore, record.QA17OriginScore, record.QA18OriginScore, record.QA19OriginScore,
// 	record.QA20OriginScore} // 原始分數
// adds = []int64{record.QA1AddScore, record.QA2AddScore, record.QA3AddScore,
// 	record.QA4AddScore, record.QA5AddScore, record.QA6AddScore, record.QA7AddScore,
// 	record.QA8AddScore, record.QA9AddScore, record.QA10AddScore, record.QA11AddScore,
// 	record.QA12AddScore, record.QA13AddScore, record.QA14AddScore, record.QA15AddScore,
// 	record.QA16AddScore, record.QA17AddScore, record.QA18AddScore, record.QA19AddScore,
// 	record.QA20AddScore} // 加成分數

// qa1Option = make([]string, dataAmount)
// qa2Option = make([]string, dataAmount)
// qa3Option = make([]string, dataAmount)
// qa4Option = make([]string, dataAmount)
// qa5Option = make([]string, dataAmount)
// qa6Option = make([]string, dataAmount)
// qa7Option = make([]string, dataAmount)
// qa8Option = make([]string, dataAmount)
// qa9Option = make([]string, dataAmount)
// qa10Option = make([]string, dataAmount)
// qa11Option = make([]string, dataAmount)
// qa12Option = make([]string, dataAmount)
// qa13Option = make([]string, dataAmount)
// qa14Option = make([]string, dataAmount)
// qa15Option = make([]string, dataAmount)
// qa16Option = make([]string, dataAmount)
// qa17Option = make([]string, dataAmount)
// qa18Option = make([]string, dataAmount)
// qa19Option = make([]string, dataAmount)
// qa20Option = make([]string, dataAmount)
// qa1OriginScore = make([]string, dataAmount)
// qa2OriginScore = make([]string, dataAmount)
// qa3OriginScore = make([]string, dataAmount)
// qa4OriginScore = make([]string, dataAmount)
// qa5OriginScore = make([]string, dataAmount)
// qa6OriginScore = make([]string, dataAmount)
// qa7OriginScore = make([]string, dataAmount)
// qa8OriginScore = make([]string, dataAmount)
// qa9OriginScore = make([]string, dataAmount)
// qa10OriginScore = make([]string, dataAmount)
// qa11OriginScore = make([]string, dataAmount)
// qa12OriginScore = make([]string, dataAmount)
// qa13OriginScore = make([]string, dataAmount)
// qa14OriginScore = make([]string, dataAmount)
// qa15OriginScore = make([]string, dataAmount)
// qa16OriginScore = make([]string, dataAmount)
// qa17OriginScore = make([]string, dataAmount)
// qa18OriginScore = make([]string, dataAmount)
// qa19OriginScore = make([]string, dataAmount)
// qa20OriginScore = make([]string, dataAmount)
// qa1AddScore = make([]string, dataAmount)
// qa2AddScore = make([]string, dataAmount)
// qa3AddScore = make([]string, dataAmount)
// qa4AddScore = make([]string, dataAmount)
// qa5AddScore = make([]string, dataAmount)
// qa6AddScore = make([]string, dataAmount)
// qa7AddScore = make([]string, dataAmount)
// qa8AddScore = make([]string, dataAmount)
// qa9AddScore = make([]string, dataAmount)
// qa10AddScore = make([]string, dataAmount)
// qa11AddScore = make([]string, dataAmount)
// qa12AddScore = make([]string, dataAmount)
// qa13AddScore = make([]string, dataAmount)
// qa14AddScore = make([]string, dataAmount)
// qa15AddScore = make([]string, dataAmount)
// qa16AddScore = make([]string, dataAmount)
// qa17AddScore = make([]string, dataAmount)
// qa18AddScore = make([]string, dataAmount)
// qa19AddScore =
