package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
)

// QuestionLikesRecordModel 資料表欄位
type QuestionLikesRecordModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id" example:"1"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	QuestionID int64  `json:"question_id" example:"1"`
	UserID     string `json:"user_id" example:"user_id"`

	// 用戶資訊
	Name   string `json:"name" example:"name"`
	Avatar string `json:"avatar" example:"avatar"`
}

// NewQuestionLikesRecordModel 資料表欄位
type NewQuestionLikesRecordModel struct {
	ActivityID string `json:"activity_id"`
	QuestionID int64  `json:"question_id"`
	UserID     string `json:"user_id"`
}

// DefaultQuestionLikesRecordModel 預設QuestionLikesRecordModel
func DefaultQuestionLikesRecordModel() QuestionLikesRecordModel {
	return QuestionLikesRecordModel{Base: Base{TableName: config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (g QuestionLikesRecordModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) QuestionLikesRecordModel {
	g.DbConn = dbconn
	g.RedisConn = cacheconn
	g.MongoConn = mongoconn
	return g
}

// SetDbConn 設定connection
// func (g QuestionLikesRecordModel) SetDbConn(conn db.Connection) QuestionLikesRecordModel {
// 	g.DbConn = conn
// 	return g
// }

// // SetRedisConn 設定connection
// func (g QuestionLikesRecordModel) SetRedisConn(conn cache.Connection) QuestionLikesRecordModel {
// 	g.RedisConn = conn
// 	return g
// }

// FindUserLikesRecordAmount 查詢用戶按讚紀錄數量
func (g QuestionLikesRecordModel) FindUserLikesRecordAmount(activityID, userID string) (int64, error) {
	var (
		sql = g.Table(g.Base.TableName).
			Select("count(*)")
	)

	if activityID != "" {
		sql = sql.Where("activity_question_likes_record.activity_id", "=", activityID)
	}
	if userID != "" {
		sql = sql.Where("activity_question_likes_record.user_id", "=", userID)
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法取得用戶按讚紀錄數量，請重新查詢")
	}

	// log.Println("item: ", item)

	return item["count(*)"].(int64), nil
}

// Find 尋找用戶在提問牆中的按讚紀錄(通過)
func (g QuestionLikesRecordModel) Find(isRedis bool, activityID, userID string,
	limit, offset int64) (records []QuestionLikesRecordModel, err error) {
	// 判斷redis裡是否有按讚紀錄資訊
	if isRedis {
		// var data string
		// data, err = g.RedisConn.HashGetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID,
		// 	userID)
		// if err == nil {
		// 	// 解碼
		// 	json.Unmarshal([]byte(data), &records)
		// }
	}

	// redis中無用戶按讚紀錄資料，從資料表中查詢
	if len(records) == 0 {
		var items = make([]map[string]interface{}, 0)
		// 該活動所有的按讚紀錄
		sql := g.Table(g.Base.TableName).
			Select("activity_question_likes_record.id", "activity_question_likes_record.activity_id",
				"activity_question_likes_record.user_id", "activity_question_likes_record.question_id",

				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_question_likes_record.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_question_likes_record.activity_id", "=", activityID).
			Where("activity_question_likes_record.user_id", "=", userID).
			// OrderBy("activity_question_likes_record.question_user_id", "desc").
			OrderBy("activity_question_likes_record.question_id", "desc")

		if limit != 0 {
			sql = sql.Limit(limit)
		}
		if offset != 0 {
			sql = sql.Offset(offset)
		}

		items, err := sql.All()
		if err != nil {
			return records, errors.New("錯誤: 無法取得用戶按讚紀錄，請重新查詢")
		}

		records = MapToQuestionLikesRecordModel(items)

		// 將按讚紀錄資訊加入redis中
		if isRedis && len(records) > 0 {
			// 將用戶按讚紀錄資料加入redis中
			// err = g.RedisConn.HashSetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID,
			// 	userID, utils.JSON(records))
			// if err != nil {
			// 	return records, errors.New("錯誤: 設置用戶按讚記錄快取資料發生問題")
			// }
		}
	}
	return
}

// Add 增加資料(用戶按讚時)
func (g QuestionLikesRecordModel) Add(isRedis bool,
	model NewQuestionLikesRecordModel) (err error) {
	// var id int64
	_, err = g.Table(g.TableName).Insert(command.Value{
		"activity_id": model.ActivityID,
		"question_id": model.QuestionID,
		"user_id":     model.UserID,
	})

	// 判斷redis裡是否有按讚紀錄資訊
	if isRedis {
		// var (
		// 	data    string
		// 	records = make([]QuestionLikesRecordModel, 0)
		// )
		// data, err = g.RedisConn.HashGetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+model.ActivityID,
		// 	model.UserID)
		// if err == nil {
		// 	// 解碼
		// 	json.Unmarshal([]byte(data), &records)
		// }

		// // 判斷redis裡是否有用戶資訊，有則不查詢資料表
		// var user LineModel
		// user, err = DefaultLineModel().SetDbConn(g.DbConn).
		// 	SetRedisConn(g.RedisConn).Find(true, "", "user_id", model.UserID)
		// if err != nil {
		// 	return errors.New("錯誤: 無法取得用戶資訊")
		// }

		// records = append(records, QuestionLikesRecordModel{
		// 	ID:         id,
		// 	ActivityID: model.ActivityID,
		// 	QuestionID: model.QuestionID,
		// 	UserID:     model.UserID,
		// 	Name:       user.Name,
		// 	Avatar:     user.Avatar,
		// })

		// // 將最新的用戶按讚紀錄資料更新至redis中
		// err = g.RedisConn.HashSetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+model.ActivityID,
		// 	model.UserID, utils.JSON(records))
		// if err != nil {
		// 	return errors.New("錯誤: 更新用戶按讚記錄快取資料發生問題")
		// }
	}
	return err
}

// Delete 刪除按讚紀錄(用戶收回讚時)
func (g QuestionLikesRecordModel) Delete(isRedis bool, questionID int64,
	activityID, userID string) (err error) {
	if err = g.Table(g.Base.TableName).Where("activity_id", "=", activityID).
		Where("user_id", "=", userID).Where("question_id", "=", questionID).Delete(); err != nil &&
		err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
		return err
	}

	// 判斷redis裡是否有按讚紀錄資訊
	if isRedis {
		// var (
		// 	data    string
		// 	records = make([]QuestionLikesRecordModel, 0)
		// )
		// data, err = g.RedisConn.HashGetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID,
		// 	userID)
		// if err == nil {
		// 	// 解碼
		// 	json.Unmarshal([]byte(data), &records)
		// }

		// for i, record := range records {
		// 	if record.QuestionID == questionID {
		// 		records = append(records[:i], records[i+1:]...) // 清除被刪除按讚紀錄
		// 	}
		// }

		// // 將最新的用戶按讚紀錄資料更新至redis中
		// err = g.RedisConn.HashSetCache(config.QUESTION_USER_LIKE_RECORDS_REDIS+activityID,
		// 	userID, utils.JSON(records))
		// if err != nil {
		// 	return errors.New("錯誤: 更新用戶按讚記錄快取資料發生問題")
		// }
	}
	return nil
}

// MapToQuestionLikesRecordModel map轉換[]QuestionLikesRecordModel
func MapToQuestionLikesRecordModel(items []map[string]interface{}) []QuestionLikesRecordModel {
	var records = make([]QuestionLikesRecordModel, 0)
	for _, item := range items {
		var (
			record QuestionLikesRecordModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &record)

		records = append(records, record)
	}
	return records
}

// record.ID, _ = item["id"].(int64)
// record.ActivityID, _ = item["activity_id"].(string)
// record.QuestionID, _ = item["question_id"].(int64)
// record.UserID, _ = item["user_id"].(string)

// 用戶資訊
// record.Name, _ = item["name"].(string)
// record.Avatar, _ = item["avatar"].(string)
