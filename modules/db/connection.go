package db

import (
	"database/sql"
	"fmt"

	"hilive/modules/config"
)

// Connection 資料庫處理程序
type Connection interface {
	Name() string
	InitDB(cfg map[string]config.Database) Connection                                              // 初始化資料庫
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)                     // 查詢
	Exec(query string, args ...interface{}) (sql.Result, error)                                    // 執行
	QueryWithConnection(conn, query string, args ...interface{}) ([]map[string]interface{}, error) // 查詢(給定conn)
	ExecWithConnection(conn, query string, args ...interface{}) (sql.Result, error)                // 執行(給定conn)
	GetDelimiter() string                                                                          // 分隔符號
	GetDelimiter2() string                                                                         // 分隔符號
	// GetTx() *sql.Tx                                                                                // 取得Tx
	// QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error)   // 利用tx查詢
	// ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)                  // 利用tx執行
}

// GetConnectionByService 取得Connection
// func GetConnectionByService(srvs service.List) Connection {
// 	if conn, ok := srvs.Get(config.GetDatabases().GetHilive().Driver).(Connection); ok {
// 		return conn
// 	}
// 	panic("錯誤的Service")
// }

// GetConnectionByDriver 取得Connection
func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case "mysql":
		return GetDefaultMysql()
	default:
		panic("找不到資料庫引擎")
	}
}

// GetConnectionFromService 取得Connection
func GetConnectionFromService(s interface{}) Connection {
	if c, ok := s.(Connection); ok {
		return c
	}
	panic("錯誤的Service")
}

// GetAggregationExpression 判斷資料庫引擎取得表達式
func GetAggregationExpression(driver, field, headField, delimiter string) string {
	switch driver {
	case "mysql":
		return fmt.Sprintf("group_concat(%s separator '%s') as %s", field, delimiter, headField)
	default:
		panic("錯誤的引擎")
	}
}
