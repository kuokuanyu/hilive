package db

import (
	"database/sql"
	"sync"
)

// Base 資料庫引擎
type Base struct {
	DbList map[string]*sql.DB // 資料庫引擎
	Once   sync.Once          // sync.Once為唯一鎖，在代碼需要被執行時，只會被執行一次
}

// Close 關閉引擎
func (db *Base) Close() []error {
	errs := make([]error, 0)
	for _, d := range db.DbList {
		errs = append(errs, d.Close())
	}
	return errs
}

// GetDB 取得引擎
func (db *Base) GetDB(key string) *sql.DB {
	return db.DbList[key]
}

