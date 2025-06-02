package auth

// DBDriver 紀錄cookie的資料表資訊
// type DBDriver struct {
// 	conn      db.Connection
// 	tableName string
// }

// Session session資訊
// type Session struct {
// Expires       time.Duration          // session 時效
// Cookie        string                 // session 名稱
// SessionID     string                 // session ID
// SessionValues map[string]interface{} // session value
// Driver        PersistenceDriver      // session 方法
// 	Context *gin.Context
// }

// PersistenceDriver 處理session資料
// type PersistenceDriver interface {
// 	Load(string) (map[string]interface{}, error)       // 透過cookie值取得session資料表資料
// 	Update(bool, string, map[string]interface{}) error // 更新session資料表資料，如果沒有資料則新增資料
// }

// table 設置SQL(struct)
// func (driver *DBDriver) table() *db.SQL {
// 	return db.Table(driver.tableName).WithConn(driver.conn)
// }

// DefaultDBDriver 預設DBDriver
// func DefaultDBDriver(conn db.Connection) *DBDriver {
// 	return &DBDriver{
// 		conn:      conn,
// 		tableName: config.SESSION_TABLE,
// 	}
// }

// InitSession 初始化session並取得資料表的session資料。如果沒有session，新增session值
// func InitSession(ctx *gin.Context, conn db.Connection, sessionName string) (bool, *Session, error) {
// 	var (
// 		session = new(Session)
// 		isExist bool // 判斷資料表裡是否有session資料
// 	)

// 	// cookie相關設置
// 	session.Expires = time.Second * time.Duration(config.GetSessionLifeTime())
// 	session.Cookie = sessionName
// 	session.Driver = DefaultDBDriver(conn)
// 	session.SessionValues = make(map[string]interface{})

// 	if c, err := ctx.Request.Cookie(session.Cookie); err == nil && c.Value != "" {
// 		session.SessionID = c.Value
// 		if valueFromDriver, err := session.Driver.Load(c.Value); err != nil {
// 			return isExist, nil, err
// 		} else if len(valueFromDriver) > 0 {
// 			// 資料表裡已有session資料
// 			isExist = true
// 			session.SessionValues = valueFromDriver
// 		}
// 	} else {
// 		// 網頁端目前沒有cookie value
// 		session.SessionID = utils.UUID(32)
// 	}
// 	session.Context = ctx
// 	return isExist, session, nil
// }

// deleteOverdueSession 刪除資料表中超過時效的session資料
// func (driver *DBDriver) deleteOverdueSession() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	if config.GetDatabases().GetHilive().Driver == "mysql" {
// 		driver.table().WhereRaw(`unix_timestamp(updated_at) < unix_timestamp() - ` +
// 			strconv.Itoa(config.GetSessionLifeTime())).Delete()
// 	}
// }

// Load 透過cookie值取得session資料表資料
// func (driver *DBDriver) Load(sessionid string) (map[string]interface{}, error) {
// 	go driver.deleteOverdueSession() // 刪除過期session
// 	go driver.deleteOverdueToken()   // 刪除過期token

// 	var (
// 		sessionModel map[string]interface{}
// 		values       map[string]interface{}
// 		err          error
// 	)
// 	if sessionModel, err = driver.table().
// 		Where("session_id", "=", sessionid).First(); err != nil {
// 		return nil, err
// 	} else if sessionModel == nil {
// 		return map[string]interface{}{}, nil
// 	}

// 	// 更新資料表session時效
// 	// if err = driver.table().Where("session_id", "=", sessionid).Update(
// 	// 	command.Value{"updated_at": time.Now().Add(time.Hour * time.Duration(8))}); err != nil &&
// 	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	// 	return values, err
// 	// }
// 	if err = json.Unmarshal([]byte(sessionModel["session_values"].(string)),
// 		&values); err != nil {
// 		return values, errors.New("錯誤: json解碼發生問題，請重新操作")
// 	}
// 	return values, nil
// }

// Update 更新session資料表資料，如果沒有資料則新增資料
// func (driver *DBDriver) Update(isExist bool, sessionid string, values map[string]interface{}) error {
// 	var (
// 		err          error
// 		sessionValue string
// 	)
// 	if sessionid != "" {
// 		if len(values) == 0 {
// 			// delete cookie時執行，刪除session資料後return
// 			if err = driver.table().Where("session_id", "=", sessionid).Delete(); err != nil &&
// 				err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
// 				return err
// 			}
// 			return nil
// 		}

// 		valueByte, err := json.Marshal(values)
// 		if err != nil {
// 			return fmt.Errorf("json編碼發生錯誤: %s", err)
// 		}
// 		sessionValue = string(valueByte)

// 		if !isExist {
// 			// session 資料不存在
// 			if !config.GetNoLimitLoginIP() { // 限制IP數量
// 				if err = driver.table().Where("session_values", "=", sessionValue).Delete(); err != nil &&
// 					err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
// 					return err
// 				}
// 			}
// 			if _, err = driver.table().Insert(command.Value{
// 				"session_values": sessionValue, "session_id": sessionid}); err != nil {
// 				return err
// 			}
// 		} else {
// 			// session資料存在
// 			if err = driver.table().Where("session_id", "=", sessionid).Update(command.Value{
// 				"session_values": sessionValue}); err != nil &&
// 				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 				return err
// 			}
// 		}

// if sesModel, _ := driver.table().Where("session_id", "=", sessionid).First(); sesModel == nil {
// 	if !config.GetNoLimitLoginIP() { // 限制IP數量
// 		if err = driver.table().Where("session_values", "=", sessionValue).Delete(); err != nil &&
// 			err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
// 			return err
// 		}
// 	}
// 	if _, err = driver.table().Insert(command.Value{
// 		"session_values": sessionValue, "session_id": sessionid}); err != nil {
// 		return err
// 	}
// } else {
// 	if err = driver.table().Where("session_id", "=", sessionid).Update(command.Value{
// 		"session_values": sessionValue}); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }
// 	}
// 	return nil
// }
