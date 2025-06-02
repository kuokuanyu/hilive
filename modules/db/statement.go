package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"hilive/modules/db/command"
)

const (
	CONNECT_NAME = "hilives"
)

// TxFn func(tx *sql.Tx) (error, map[string]interface{})
type TxFn func(tx *sql.Tx) (error, map[string]interface{})

// SQLPool 降低內存負擔，用於需要被重複分配、回收內存的地方
var SQLPool = sync.Pool{
	New: func() interface{} {
		return &SQL{
			filter: command.Filter{
				Fields:    make([]string, 0),
				TableName: "",
				Args:      make([]interface{}, 0),
				Wheres:    make([]command.Where, 0),
				Leftjoins: make([]command.Join, 0),
				Order:     "",
				Group:     "",
				Limit:     "",
			},
			conn:    nil,
			command: nil,
		}
	},
}

// newSQL 取得SQL
func newSQL() *SQL {
	return SQLPool.Get().(*SQL)
}

// SQL 過濾條件、CRUD方法、Conn...等資訊
type SQL struct {
	filter   command.Filter  // 篩選條件
	conn     Connection      // 資料庫的處理程序
	connName string          // connection名稱
	command  command.Command // 處理資料庫CRUD命令
	tx       *sql.Tx
}

// Table 設置SQL
func Table(table string) *SQL {
	sqlpool := newSQL()
	sqlpool.filter.TableName = table
	sqlpool.connName = CONNECT_NAME
	return sqlpool
}

// Conn 設置SQL
func Conn(conn Connection) *SQL {
	sql := newSQL()
	sql.conn = conn
	sql.command = command.GetCommand(conn.Name())
	sql.connName = CONNECT_NAME
	return sql
}

// WithConn 設置connection
func (s *SQL) WithConn(conn Connection) *SQL {
	s.conn = conn
	s.command = command.GetCommand(conn.Name())
	return s
}

// WithConnName 設置connName
func (s *SQL) WithConnName(conn string) *SQL {
	s.connName = conn
	return s
}

// Table 清除過濾條件後設置table
func (s *SQL) Table(table string) *SQL {
	s.clean()
	s.filter.TableName = table
	return s
}

// ShowColumns 所有欄位資訊
func (s *SQL) ShowColumns() ([]map[string]interface{}, error) {
	defer RecycleSQL(s)
	return s.conn.QueryWithConnection(s.connName, s.command.ShowColumns(s.filter.TableName))
}

// BatchInsert 匹量插入一千筆
func (s *SQL) BatchInsert(dataAmount int, fields string, values [][]string) error {
	defer RecycleSQL(s)

	batchSize := 1000 // 匹量插入設置一千筆

	if dataAmount > 0 && len(values) > 0 {
		// 比對value參數陣列長度是否跟欄位數量相等
		if len(values) != len(strings.Split(fields, ",")) {
			return errors.New("錯誤: 欄位數量與值數量不匹配")
		}

		// 比對value陣列裡的陣列資料長度是否跟dataAmount數量相等
		for i := 0; i < len(values); i++ {
			if len(values[i]) != dataAmount {
				// log.Println("錯誤: 資料不完整或與資料數量不匹配")
				return errors.New("錯誤: 資料不完整或與資料數量不匹配")
			}
		}

		// 開始批量處理
		for start := 0; start < dataAmount; start += batchSize {
			end := start + batchSize
			if end > dataAmount {
				end = dataAmount
			}

			// 提取當前批次資料
			currentValues := make([][]string, len(values))
			for i := range values {
				currentValues[i] = values[i][start:end]
			}

			// 呼叫插入函式
			if err := s.Inserts(end-start, fields, currentValues); err != nil {
				// log.Println("錯誤3")
				return err
			}
		}
	}
	return nil
}

// Inserts 匹量插入資料
func (s *SQL) Inserts(dataAmount int, fields string, values [][]string) error {
	// defer RecycleSQL(s)
	var (
		res sql.Result
		err error
	)

	err = s.command.Inserts(&s.filter, dataAmount, fields, values)
	if err != nil {
		log.Println("錯誤: 處理匹量新增資料庫語法處理發生問題")
		return errors.New("錯誤: 處理匹量新增資料庫語法處理發生問題")
	}

	// log.Println("insert command: ", s.filter.Statement)
	// if s.tx != nil {
	// 	res, err = s.conn.ExecWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// } else {
	res, err = s.conn.ExecWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// }
	if err != nil {
		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("匹量新增資料發生問題: %s", err.Error()))
		return errors.New("錯誤: 匹量新增資料發生問題，請重新操作")
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return errors.New("錯誤: 無新增任何資料，請重新操作")
	}

	return nil
}

// Insert 單筆插入資料
func (s *SQL) Insert(values command.Value) (int64, error) {
	defer RecycleSQL(s)
	var (
		res sql.Result
		err error
	)
	s.filter.Values = values
	s.command.Insert(&s.filter)
	// log.Println("insert command: ", s.filter.Statement)
	// if s.tx != nil {
	// 	res, err = s.conn.ExecWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// } else {
	res, err = s.conn.ExecWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// }
	if err != nil {
		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("單筆新增資料發生問題: %s, %s", err.Error(), s.filter.Statement))
		return 0, errors.New("錯誤: 新增資料發生問題，請重新操作")
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("錯誤: 無新增任何資料，請重新操作")
	}
	return res.LastInsertId()
}

// Update 資料更新
func (s *SQL) Update(values command.Value) error {
	var err error
	defer RecycleSQL(s)
	s.filter.Values = values
	s.command.Update(&s.filter)
	// if s.tx != nil {
	// 	_, err = s.conn.ExecWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// } else {
	res, err := s.conn.ExecWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// log.Println("update command: ", s.filter.Statement)
	// log.Println(s.filter.Args...)
	// }
	if err != nil {
		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("編輯資料發生問題: %s, %s", err.Error(), s.filter.Statement))
		return errors.New("錯誤: 更新資料發生問題，請重新操作")
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return errors.New("錯誤: 無更新任何資料，請重新操作")
	}
	return nil
}

// Delete 資料刪除
func (s *SQL) Delete() error {
	defer RecycleSQL(s)
	var (
		err error
	)
	s.command.Delete(&s.filter)
	// log.Println("delete command: ", s.filter.Statement)
	// if s.tx != nil {
	// 	_, err = s.conn.ExecWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// } else {
	res, err := s.conn.ExecWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// }
	if err != nil {
		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("刪除資料發生問題: %s, %s", err.Error(), s.filter.Statement))
		return errors.New("錯誤: 刪除資料發生問題，請重新操作")
	}

	affectRow, _ := res.RowsAffected()
	if affectRow < 1 {
		return errors.New("錯誤: 無刪除任何資料，請重新操作")
	}
	return nil
}

// All 所有資料
func (s *SQL) All() ([]map[string]interface{}, error) {
	defer RecycleSQL(s)
	s.command.Select(&s.filter)
	// log.Println("select command: ", s.filter.Statement)
	// if s.tx != nil {
	// 	return s.conn.QueryWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// }
	res, err := s.conn.QueryWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// }
	if err != nil {
		// log.Println("錯誤???")

		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("查詢多筆資料發生問題: %s", err.Error()))
		return nil, errors.New("錯誤: 查詢資料發生問題，請重新操作")
	}

	return res, nil
}

// Find 尋找資料
// func (s *SQL) Find(field string, arg interface{}) (map[string]interface{}, error) {
// 	return s.Where(field, "=", arg, "and").First()
// }

// First 查詢單筆資料
func (s *SQL) First() (map[string]interface{}, error) {
	defer RecycleSQL(s)
	var (
		res []map[string]interface{}
		err error
	)
	s.command.Select(&s.filter)
	// log.Println("select command: ", s.filter.Statement)
	// if s.tx != nil {
	// 	res, err = s.conn.QueryWithTx(s.tx, s.filter.Statement, s.filter.Args...)
	// } else {
	res, err = s.conn.QueryWithConnection(s.connName, s.filter.Statement, s.filter.Args...)
	// }
	if err != nil {
		// 插入error資料到資料表
		s.conn.ExecWithConnection(s.connName,
			"INSERT INTO operation_error_log (user_id, message) VALUES (?,?,?,?,?,?)",
			"後端除錯", fmt.Sprintf("查詢單筆資料發生問題: %s", err.Error()))
		return nil, errors.New("錯誤: 查詢資料發生問題，請重新操作")
	}
	if len(res) < 1 {
		return nil, nil
	}
	return res[0], nil
}

// Select 處理欄位
func (s *SQL) Select(fields ...string) *SQL {
	s.filter.Fields = fields
	return s
}

// Where where條件
func (s *SQL) Where(field, operation string,
	arg interface{}, combinations ...string) *SQL {
	var (
		combination string
	)
	// and、or
	if len(combinations) > 0 {
		combination = combinations[0]
	} else {
		combination = "and"
	}
	s.filter.Wheres = append(s.filter.Wheres, command.Where{
		Field:       field,
		Operation:   operation,
		Value:       "?",
		Combination: combination,
	})
	s.filter.Args = append(s.filter.Args, arg)
	return s
}

// WhereRaw set WhereRaws and arguments.
func (s *SQL) WhereRaw(raw string, args ...interface{}) *SQL {
	s.filter.WhereRaws = raw
	s.filter.Args = append(s.filter.Args, args...)
	return s
}

// WhereIn where in 多個數值
func (s *SQL) WhereIn(field string, arg []interface{}, combinations ...string) *SQL {
	var (
		combination string
	)

	if len(arg) == 0 {
		panic("錯誤: wherein函式arg參數不能為空")
	}
	// and、or
	if len(combinations) > 0 {
		combination = combinations[0]
	} else {
		combination = "and"
	}
	s.filter.Wheres = append(s.filter.Wheres,
		command.Where{
			Field:       field,
			Operation:   "in",
			Value:       "(" + strings.Repeat("?,", len(arg)-1) + "?)",
			Combination: combination,
		})
	s.filter.Args = append(s.filter.Args, arg...)
	return s
}

// LeftJoin Join語法
func (s *SQL) LeftJoin(join command.Join) *SQL {
	s.filter.Leftjoins = append(s.filter.Leftjoins, join)
	return s
}

// Limit set limit value.
func (sql *SQL) Limit(limit int64) *SQL {
	sql.filter.Limit = strconv.Itoa(int(limit))
	return sql
}

// Offset set offset value.
func (sql *SQL) Offset(offset int64) *SQL {
	sql.filter.Offset = strconv.Itoa(int(offset))
	return sql
}

// OrderBy 資料排序
func (s *SQL) OrderBy(fields ...string) *SQL {
	if len(fields) == 0 {
		panic("錯誤: OrderBy資料排序語法")
	} else if len(fields)%2 == 1 {
		// 單數
		panic("錯誤: OrderBy參數不能為單數")
	}

	for i := 0; i < len(fields); i++ {
		s.filter.Order += " " + fields[i]
		if i != len(fields)-1 && i%2 == 1 {
			s.filter.Order += ","
		}
	}

	// if len(fields) == 2 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1]
	// } else if len(fields) == 4 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1] + ", " +
	// 		fields[2] + " " + fields[3]
	// } else if len(fields) == 6 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1] + ", " +
	// 		fields[2] + " " + fields[3] + ", " +
	// 		fields[4] + " " + fields[5]
	// } else if len(fields) == 8 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1] + ", " +
	// 		fields[2] + " " + fields[3] + ", " +
	// 		fields[4] + " " + fields[5] + ", " +
	// 		fields[6] + " " + fields[7]
	// } else if len(fields) == 10 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1] + ", " +
	// 		fields[2] + " " + fields[3] + ", " +
	// 		fields[4] + " " + fields[5] + ", " +
	// 		fields[6] + " " + fields[7] + ", " +
	// 		fields[8] + " " + fields[9]
	// } else if len(fields) == 12 {
	// 	s.filter.Order += " " + fields[0] + " " + fields[1] + ", " +
	// 		fields[2] + " " + fields[3] + ", " +
	// 		fields[4] + " " + fields[5] + ", " +
	// 		fields[6] + " " + fields[7] + ", " +
	// 		fields[8] + " " + fields[9] + ", " +
	// 		fields[10] + " " + fields[11]
	// }
	return s
}

// RecycleSQL 清空SQL
func RecycleSQL(s *SQL) {
	s.clean()
	// s.conn = nil
	s.command = nil
	s.tx = nil
	SQLPool.Put(s)
}

// clean 清空過濾條件
func (s *SQL) clean() {
	s.filter.Group = ""
	s.filter.Values = make(map[string]interface{})
	s.filter.Fields = make([]string, 0)
	s.filter.TableName = ""
	s.filter.Wheres = make([]command.Where, 0)
	s.filter.WhereRaws = ""
	s.filter.Leftjoins = make([]command.Join, 0)
	s.filter.Args = make([]interface{}, 0)
	s.filter.Order = ""
	s.filter.Offset = ""
	s.filter.Limit = ""
	s.filter.Statement = ""
}

// SetTx 設置Tx
// func (s *SQL) SetTx(tx *sql.Tx) *SQL {
// 	s.tx = tx
// 	return s
// }

// WithTransaction 取得Tx，持續執行commit、rollback
// func (s *SQL) WithTransaction(fn TxFn) (res map[string]interface{}, err error) {
// 	tx := s.conn.GetTx()
// 	defer func() {
// 		if p := recover(); p != nil {
// 			_ = tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			_ = tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()
// 	err, res = fn(tx)
