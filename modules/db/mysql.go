package db

import (
	"database/sql"
	"log"
	"regexp"
	"strings"
	"time"

	"hilive/modules/config"
)

// Mysql 資料庫引擎
type Mysql struct {
	Base
}

// GetDefaultMysql 設置Mysql
func GetDefaultMysql() *Mysql {
	return &Mysql{
		Base: Base{
			DbList: make(map[string]*sql.DB),
		},
	}
}

// Name 引擎名稱
func (db *Mysql) Name() string {
	return "mysql"
}

// InitDB 初始化資料庫引擎
func (db *Mysql) InitDB(cfgs map[string]config.Database) Connection {
	db.Once.Do(func() {
		for conn, cfg := range cfgs {
			if cfg.Dsn == "" {
				cfg.Dsn = cfg.User + ":" + cfg.Pwd + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" +
					cfg.Name
			}
			sqlDB, err := sql.Open("mysql", cfg.Dsn)
			if err != nil {
				if sqlDB != nil {
					_ = sqlDB.Close()
				}
				panic("連接資料庫引擎發生錯誤")
			} else {
				// SetMaxIdleConns：設定閒置連線量，如果設定0或是小於0就是沒有閒置連線，一但作完了連線就不回連線池待下次取用直接廢棄掉，變成每次都要起新的連線，沒辦法節省資源
				sqlDB.SetMaxIdleConns(cfg.MaxIdleCon)

				// SetMaxOpenConns：設定最大連線數，如果設定0或是小於0就是是無限大，但是基本上不會設定成無限大，因為db會收不下過量的連線
				sqlDB.SetMaxOpenConns(cfg.MaxOpenCon)

				// SetConnMaxLifetime：設置了連接可重用的最大時間，如果設定0或是小於0就是沒有生命週期，設定時間到期後這些連線就會廢棄無法重用，需要重新起連線。
				sqlDB.SetConnMaxLifetime(60 * time.Second)
				db.DbList[conn] = sqlDB
			}
			if err := sqlDB.Ping(); err != nil {
				panic("執行資料庫引擎發生錯誤")
			}
		}
	})
	return db
}

// Query 查詢
func (db *Mysql) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	return commonQuery(db.DbList["hilives"], query, args...)
}

// Exec 執行
func (db *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	return commonExec(db.DbList["hilives"], query, args...)
}

// GetTx 取得Tx
// func (db *Mysql) GetTx() *sql.Tx {
// 	return commonBeginTxWithLevel(db.Base.DbList["hilives"], sql.LevelDefault)
// }

// QueryWithConnection 查詢(給定conn)
func (db *Mysql) QueryWithConnection(con string, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// fmt.Println("query: ", query)
	return commonQuery(db.DbList[con], query, args...)
}

// ExecWithConnection 執行(給定conn)
func (db *Mysql) ExecWithConnection(con string, query string, args ...interface{}) (sql.Result, error) {
	return commonExec(db.DbList[con], query, args...)
}

// GetDelimiter 分隔符號
func (db *Mysql) GetDelimiter() string {
	return "`"
}

// GetDelimiter2 分隔符號
func (db *Mysql) GetDelimiter2() string {
	return "`"
}

// QueryWithTx 利用tx查詢
// func (db *Mysql) QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {
// 	return commonQueryWithTx(tx, query, args...)
// }

// ExecWithTx 利用tx執行
// func (db *Mysql) ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
// 	return commonExecWithTx(tx, query, args...)
// }

// commonQuery 查詢
func commonQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	var (
		rs      *sql.Rows
		typeVal []*sql.ColumnType
		col     []string
		results = make([]map[string]interface{}, 0)
		err     error
	)
	// log.Println("select command: ", query, "args: ", args)
	if rs, err = db.Query(query, args...); err != nil {
		// log.Println("資料庫命令: ", query)
		// log.Println("資料庫值", args)
		log.Println("執行db.commonQuery函式查詢資料發生錯誤: ", err)

		// 插入error資料到資料表
		// db.Exec(
		// 	"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
		// 	"後端除錯", fmt.Sprintf("查詢資料發生問題: %s", err.Error()))

		return nil, err
		// panic(fmt.Sprintf("執行db.commonQuery函式查詢資料發生錯誤: %s", err))
	}
	defer func() {
		if rs != nil {
			_ = rs.Close()
		}
	}()
	if col, err = rs.Columns(); err != nil {
		log.Println("執行db.commonQuery函式中的rs.Columns發生錯誤: ", err)
		return nil, err
	}
	if typeVal, err = rs.ColumnTypes(); err != nil {
		log.Println("執行db.commonQuery函式中的rs.ColumnTypes發生錯誤: ", err)
		return nil, err
	}
	for rs.Next() {
		var colVar = make([]interface{}, len(col))
		r, _ := regexp.Compile(`\\((.*)\\)`)
		for i := 0; i < len(col); i++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[i].DatabaseTypeName(), ""))
			SetColVarType(&colVar, i, typeName)
		}
		result := make(map[string]interface{})
		if err := rs.Scan(colVar...); err != nil {
			log.Println("執行db.commonQuery函式中的rs.Scan發生錯誤: ", err)
			return nil, err
		}
		for j := 0; j < len(col); j++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[j].DatabaseTypeName(), ""))
			SetResultValue(&result, col[j], colVar[j], typeName)
		}
		results = append(results, result)
	}
	if err := rs.Err(); err != nil {
		log.Println("執行db.commonQuery函式中的rs.Err發生錯誤: ", err)
		return nil, err
	}
	return results, nil
}

// commonExec 執行
func commonExec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	rs, err := db.Exec(query, args...)
	if err != nil {
		// log.Println("資料庫命令: ", query)
		// log.Println("資料庫值", args)
		log.Println("執行db.commonExec函式發生錯誤: ", err)
		return nil, err
	}

	return rs, nil
}

// isDeadlockError 判斷錯誤訊息是否為死鎖
func isDeadlockError(err error) bool {
	if err == nil {
		return false
	}

	// 檢查 MySQL 錯誤代碼或錯誤訊息是否與死鎖相關
	return strings.Contains(err.Error(), "1205") || strings.Contains(err.Error(), "lock wait timeout exceeded")
}

// CommonBeginTxWithLevel 取得Tx
// func commonBeginTxWithLevel(db *sql.DB, level sql.IsolationLevel) *sql.Tx {
// 	if tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: level}); err != nil {
// 		log.Println("執行db.commonBeginTxWithLevel函式發生錯誤: ", err)
// 		panic(err)
// 	} else {
// 		return tx
// 	}
// }

// commonExecWithTx 利用tx執行
// func commonExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
// 	rs, err := tx.Exec(query, args...)
// 	if err != nil {
// 		log.Println("資料庫命令: ", query)
// 		log.Println("資料庫值", args)
// 		log.Println("執行db.commonExecWithTx函式發生錯誤: ", err)
// 		return nil, err
// 	}
// 	return rs, nil
// }

// commonQueryWithTx 利用tx查詢
// func commonQueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {
// 	var (
// 		rs      *sql.Rows
// 		typeVal []*sql.ColumnType
// 		col     []string
// 		results = make([]map[string]interface{}, 0)
// 		err     error
// 	)
// 	if rs, err = tx.Query(query, args...); err != nil {
// 		log.Println("資料庫命令: ", query)
// 		log.Println("資料庫值", args)
// 		panic(fmt.Sprintf("執行db.commonQueryWithTx函式查詢資料發生錯誤: %s", err))
// 	}
// 	defer func() {
// 		if rs != nil {
// 			_ = rs.Close()
// 		}
// 	}()
// 	col, colErr := rs.Columns()
// 	if colErr != nil {
// 		log.Println("執行db.commonQueryWithTx函式中的rs.Columns發生錯誤: ", err)
// 		return nil, err
// 	}
// 	typeVal, err = rs.ColumnTypes()
// 	if err != nil {
// 		log.Println("執行db.commonQueryWithTx函式中的rs.ColumnTypes發生錯誤: ", err)
// 		return nil, err
// 	}
// 	for rs.Next() {
// 		var colVar = make([]interface{}, len(col))
// 		r, _ := regexp.Compile(`\\((.*)\\)`)
// 		for i := 0; i < len(col); i++ {
// 			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[i].DatabaseTypeName(), ""))
// 			SetColVarType(&colVar, i, typeName)
// 		}
// 		result := make(map[string]interface{})
// 		if err = rs.Scan(colVar...); err != nil {
// 			log.Println("執行db.commonQueryWithTx函式中的rs.Scan發生錯誤: ", err)
// 			return nil, err
// 		}
// 		for j := 0; j < len(col); j++ {
// 			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[j].DatabaseTypeName(), ""))
// 			SetResultValue(&result, col[j], colVar[j], typeName)
// 		}
// 		results = append(results, result)
// 	}
// 	if err := rs.Err(); err != nil {
// 		log.Println("執行db.commonQueryWithTx函式中的rs.Err發生錯誤: ", err)
// 		return nil, err
// 	}
// 	return results, nil
// }
