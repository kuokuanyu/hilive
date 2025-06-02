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

// ErrorLogModel 資料表欄位
type ErrorLogModel struct {
	Base      `json:"-"`
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	Code      string `json:"code"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	PathQuery string `json:"path_query"`

	// 用戶
	Name string `json:"name"`
}

// EditErrorLogModel 資料表欄位
type EditErrorLogModel struct {
	UserID    string `json:"user_id"`
	Code      int64  `json:"code"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	PathQuery string `json:"path_query"`
}

// DefaultErrorLogModel 預設ErrorLogModel
func DefaultErrorLogModel() ErrorLogModel {
	return ErrorLogModel{Base: Base{TableName: config.OPERATION_ERROR_LOG_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (l ErrorLogModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ErrorLogModel {
	l.DbConn = dbconn
	l.RedisConn = cacheconn
	l.MongoConn = mongoconn
	return l
}

// SetDbConn 設定connection
// func (l ErrorLogModel) SetDbConn(conn db.Connection) ErrorLogModel {
// 	l.DbConn = conn
// 	return l
// }

// FindAll 查詢所有錯誤日誌資料
func (l ErrorLogModel) FindAll() ([]ErrorLogModel, error) {
	items, err := l.Table(l.Base.TableName).
		Select("operation_error_log.id", "operation_error_log.user_id",
			"operation_error_log.code", "operation_error_log.method",
			"operation_error_log.path", "operation_error_log.message",
			"operation_error_log.path_query",

			// 用戶
			"users.name").
		LeftJoin(command.Join{
			FieldA:    "operation_error_log.user_id",
			FieldA1:   "users.user_id",
			Table:     "users",
			Operation: "="}).
		OrderBy("operation_error_log.id", "desc").
		Limit(100).
		All()
	if err != nil {
		return []ErrorLogModel{}, errors.New("錯誤: 無法取得錯誤日誌資訊，請重新查詢")
	}

	return MapToErrorLogModel(items), nil
}

// Add 新增錯誤日誌資料
func (l ErrorLogModel) Add(model EditErrorLogModel) error {
	_, err := l.Table(l.TableName).Insert(command.Value{
		"user_id":    model.UserID,
		"code":       model.Code,
		"method":     model.Method,
		"path":       model.Path,
		"message":    model.Message,
		"path_query": model.PathQuery,
	})
	if err != nil {
		return errors.New("錯誤: 新增錯誤日誌資料發生問題，請重新操作")
	}
	return nil
}

// MapToErrorLogModel map轉換[]ErrorLogModel
func MapToErrorLogModel(items []map[string]interface{}) []ErrorLogModel {
	var logs = make([]ErrorLogModel, 0)
	for _, item := range items {
		var (
			log ErrorLogModel
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
// log.Code, _ = item["cpde"].(string)
// log.Method, _ = item["method"].(string)
// log.Path, _ = item["path"].(string)
// log.Message, _ = item["message"].(string)
// log.PathQuery, _ = item["path_query"].(string)

// // 用戶
// log.Name, _ = item["name"].(string)
