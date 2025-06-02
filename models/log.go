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

// LogModel 資料表欄位
type LogModel struct {
	Base   `json:"-"`
	ID     int64  `json:"id"`
	UserID string `json:"user_id"`
	Method string `json:"method"`
	Path   string `json:"path"`

	// 用戶
	Name string `json:"name"`
}

// EditLogModel 資料表欄位
type EditLogModel struct {
	UserID string `json:"user_id"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

// DefaultLogModel 預設LogModel
func DefaultLogModel() LogModel {
	return LogModel{Base: Base{TableName: config.OPERATION_LOG_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (l LogModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) LogModel {
	l.DbConn = dbconn
	l.RedisConn = cacheconn
	l.MongoConn = mongoconn
	return l
}

// SetDbConn 設定connection
// func (l LogModel) SetDbConn(conn db.Connection) LogModel {
// 	l.DbConn = conn
// 	return l
// }

// FindAll 查詢所有操作日誌資料
func (l LogModel) FindAll() ([]LogModel, error) {
	items, err := l.Table(l.Base.TableName).
		Select("operation_log.id", "operation_log.user_id",
			"operation_log.method", "operation_log.path",

			// 用戶
			"users.name").
		LeftJoin(command.Join{
			FieldA:    "operation_log.user_id",
			FieldA1:   "users.user_id",
			Table:     "users",
			Operation: "="}).
		OrderBy("operation_log.id", "desc").
		Limit(100).
		All()
	if err != nil {
		return []LogModel{}, errors.New("錯誤: 無法取得操作日誌資訊，請重新查詢")
	}

	return MapToLogModel(items), nil
}

// Add 新增用戶
func (l LogModel) Add(model EditLogModel) error {
	_, err := l.Table(l.TableName).Insert(command.Value{
		"user_id": model.UserID,
		"method":  model.Method,
		"path":    model.Path,
	})
	if err != nil {
		return errors.New("錯誤: 新增操作日誌發生問題，請重新操作")
	}
	return nil
}

// MapToLogModel map轉換[]LogModel
func MapToLogModel(items []map[string]interface{}) []LogModel {
	var logs = make([]LogModel, 0)
	for _, item := range items {
		var (
			log LogModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &log)

		logs = append(logs, log)
	}
	return logs
}

// log.ID, _ = item["id"].(int64)
// log.UserID, _ = item["user_id"].(string)
// log.Method, _ = item["method"].(string)
// log.Path, _ = item["path"].(string)

// 用戶
// log.Name, _ = item["name"].(string)
