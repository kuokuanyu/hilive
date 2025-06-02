package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// MessageSensitivityModel 資料表欄位
type MessageSensitivityModel struct {
	Base            `json:"-"`
	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id"`
	SensitivityWord string `json:"sensitivity_word"`
	ReplaceWord     string `json:"replace_word"`
}

// EditMessageSensitivityModel 資料表欄位
type EditMessageSensitivityModel struct {
	ID              string `json:"id"`
	ActivityID      string `json:"activity_id"`
	SensitivityWord string `json:"sensitivity_word"`
	ReplaceWord     string `json:"replace_word"`
}

// DefaultMessageSensitivityModel 預設MessageSensitivityModel
func DefaultMessageSensitivityModel() MessageSensitivityModel {
	return MessageSensitivityModel{Base: Base{TableName: config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m MessageSensitivityModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) MessageSensitivityModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m MessageSensitivityModel) SetDbConn(conn db.Connection) MessageSensitivityModel {
// 	m.DbConn = conn
// 	return m
// }

// FindAmount 查詢該活動敏感詞數量
func (m MessageSensitivityModel) FindAmount(activityID string) (int64, error) {
	var (
		sql = m.Table(m.Base.TableName).
			Select("count(*)")
	)

	if activityID != "" {
		sql = sql.Where("activity_message_sensitivity.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法取得該活動敏感詞數量，請重新查詢")
	}

	return item["count(*)"].(int64), nil
}

// Find 尋找資料
func (m MessageSensitivityModel) Find(activityID string) ([]MessageSensitivityModel, error) {
	items, err := m.Table(m.Base.TableName).
		Where("activity_id", "=", activityID).
		OrderBy("id", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得訊息敏感詞資料，請重新查詢")
	}
	return MapToMessageSensitivityModel(items), nil
}

// Add 增加敏感詞資料
func (m MessageSensitivityModel) Add(model EditMessageSensitivityModel) error {
	var (
		fields = []string{
			"activity_id",
			"sensitivity_word",
			"replace_word",
		}
	)
	// if model.SensitivityWord == "" {
	// 	return errors.New("錯誤: 敏感詞資料發生問題，請輸入有效的敏感詞資料")
	// }

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	_, err := m.Table(m.TableName).Insert(FilterFields(data, fields))

	return err
}

// Update 更新敏感詞資料
func (m MessageSensitivityModel) Update(model EditMessageSensitivityModel) error {
	var (
		fieldValues = command.Value{"replace_word": model.ReplaceWord}
		fields      = []string{"sensitivity_word"}
	)

	if _, err := strconv.Atoi(model.ID); err != nil {
		return errors.New("錯誤: ID資料發生問題，請輸入有效的ID")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	return m.Table(m.Base.TableName).
		Where("id", "=", model.ID).Update(fieldValues)
}

// MapToMessageSensitivityModel map轉換[]MessageSensitivityModel
func MapToMessageSensitivityModel(items []map[string]interface{}) []MessageSensitivityModel {
	var records = make([]MessageSensitivityModel, 0)
	for _, item := range items {
		var record MessageSensitivityModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &record)

		records = append(records, record)
	}
	return records
}

// record.ID, _ = item["id"].(int64)
// record.ActivityID, _ = item["activity_id"].(string)
// record.SensitivityWord, _ = item["sensitivity_word"].(string)
// record.ReplaceWord, _ = item["replace_word"].(string)

// command.Value{
// 	"activity_id":      model.ActivityID,
// 	"sensitivity_word": model.SensitivityWord,
// 	"replace_word":     model.ReplaceWord,
// }

// values      = []string{model.SensitivityWord}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
