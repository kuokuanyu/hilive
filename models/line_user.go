package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
)

// LineModel 資料表欄位
type LineModel struct {
	Base     `json:"-"`
	ID       int64  `json:"id"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	ExtEmail string `json:"ext_email"`
	Phone    string `json:"phone"`
	Identify string `json:"identify"`
	Friend   string `json:"friend"`
	// *****可以同時登入(暫時拿除)*****
	// Ip       string `json:"-"`
	// *****可以同時登入(暫時拿除)*****
	Line     string `json:"line"`
	Device   string `json:"device"`
	IsModify string `json:"is_modify"`
	AdminID  string `json:"admin_id"`

	// ActivityID  string `json:"activity_id"`  // 自定義簽到人員活動主辦方
	// ExtAccount  string `json:"ext_account"`  // 自定義簽到人員用戶
	ExtPassword string `json:"ext_password"` // 自定義簽到人員密碼
}

// DefaultLineModel 預設LineModel
func DefaultLineModel() LineModel {
	return LineModel{Base: Base{TableName: config.LINE_USERS_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (user LineModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) LineModel {
	user.DbConn = dbconn
	user.RedisConn = cacheconn
	user.MongoConn = mongoconn
	return user
}

// SetDbConn 設定connection
// func (user LineModel) SetDbConn(conn db.Connection) LineModel {
// 	user.DbConn = conn
// 	return user
// }

// // SetRedisConn 設定connection
// func (user LineModel) SetRedisConn(conn cache.Connection) LineModel {
// 	user.RedisConn = conn
// 	return user
// }

// Find 尋找資料(從redis取得，如果沒有才執行資料表查詢)
func (user LineModel) Find(isRedis bool, ip string, field string, value string) (LineModel, error) {
	if isRedis {
		// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
		dataMap, err := user.RedisConn.HashGetAllCache(config.AUTH_USERS_REDIS + value)
		if err != nil {
			return LineModel{}, errors.New("錯誤: 取得用戶快取資料發生問題")
		}
		id, _ := strconv.Atoi(dataMap["id"])
		user.ID = int64(id)
		user.UserID, _ = dataMap["user_id"]
		user.Name, _ = dataMap["name"]
		user.Avatar, _ = dataMap["avatar"]
		user.Phone, _ = dataMap["phone"]
		user.ExtEmail, _ = dataMap["ext_email"]
		user.Email, _ = dataMap["email"]
		user.Identify, _ = dataMap["identify"]
		user.Friend, _ = dataMap["friend"]
		user.Device, _ = dataMap["device"]
		user.ExtPassword, _ = dataMap["ext_password"]
		user.IsModify, _ = dataMap["is_modify"]
		user.AdminID, _ = dataMap["admin_id"]
		// *****可以同時登入(暫時拿除)*****
		// user.Ip, _ = dataMap["ip"]
		// *****可以同時登入(暫時拿除)*****
	}

	if user.UserID == "" {
		item, err := user.Table(user.Base.TableName).Where(field, "=", value).First()
		if err != nil {
			return LineModel{}, errors.New("錯誤: 無法取得用戶資訊，請重新查詢")
		}
		if item == nil {
			return LineModel{}, nil
		}

		user = user.MapToModel(item)
		// 將用戶資訊加入redis
		if isRedis {
			values := []interface{}{config.AUTH_USERS_REDIS + user.UserID}
			values = append(values, "id", user.ID)
			values = append(values, "user_id", user.UserID)
			values = append(values, "name", user.Name)
			values = append(values, "phone", user.Phone)
			values = append(values, "ext_email", user.ExtEmail)
			values = append(values, "email", user.Email)
			values = append(values, "avatar", user.Avatar)
			values = append(values, "identify", user.Identify)
			values = append(values, "friend", user.Friend)
			values = append(values, "device", user.Device)
			values = append(values, "ext_password", user.ExtPassword)
			values = append(values, "is_modify", user.IsModify)
			values = append(values, "admin_id", user.AdminID)
			// *****可以同時登入(暫時拿除)*****
			// values = append(values, "ip", ip)
			// *****可以同時登入(暫時拿除)*****

			if err := user.RedisConn.HashMultiSetCache(values); err != nil {
				return user, errors.New("錯誤: 設置用戶快取資料發生問題")
			}
			// 設置過期時間
			// user.RedisConn.SetExpire(config.AUTH_USERS_REDIS+user.UserID,
			// 	config.REDIS_EXPIRE)
		}
	}
	return user, nil
}

// FindAllCustomizeUsers 尋找該活動所有自定義人員資料(資料庫)
func (user LineModel) FindAllCustomizeUsers(activityID string) ([]LineModel, error) {
	var (
		users = make([]LineModel, 0)
	)

	items, err := user.Table(user.Base.TableName).
		Where("activity_id", "=", activityID).All()
	if err != nil {
		return users, errors.New("錯誤: 無法取得所有自定義人員資訊，請重新查詢")
	}

	users = MapToLineModel(items)

	return users, nil
}

// UpdateUser 更新用戶(姓名、頭像、信箱、好友欄位)，會清除舊的redis資料
func (user LineModel) UpdateUser(model LineModel) error {
	var (
		fields = []string{
			"name",
			"avatar",
			"email",
			"friend",
			"is_modify"}

		fieldValues = command.Value{}
	)

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) == 0 {
		return nil
	}

	err := user.Table(user.TableName).
		Where("user_id", "=", model.UserID).
		Update(fieldValues)
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新用戶資訊發生問題，請重新操作")
	}

	// 需要修改的用戶必須清除舊的redis資料
	user.RedisConn.DelCache(config.AUTH_USERS_REDIS + model.UserID)

	return nil
}

// UpdatePhoneAndEmail 更新電話以及自定義電子信箱資訊
func (user LineModel) UpdatePhoneAndEmail(model LineModel) error {
	var (
		fields      = []string{"phone", "ext_email"}
		values      = []string{model.Phone, model.ExtEmail}
		fieldValues = command.Value{}
	)
	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	if err := user.Table(user.TableName).
		Where("user_id", "=", model.UserID).
		Update(fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新LINE用戶電話、信箱資料發生問題，請重新操作")
	}
	return nil
}

// UpdateAvatar 更新頭像資訊(批量)
func (user LineModel) UpdateAvatar(userIDs []interface{}, avatar string) error {

	if err := user.Table(user.TableName).
		WhereIn("user_id", userIDs).
		Update(command.Value{
			"avatar": avatar,
		}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新自定義人員頭像資料發生問題，請重新操作")
	}
	return nil
}

// UpdateAdminID 更新用戶admin_id欄位資料(綁定處理)
func (user LineModel) UpdateAdminID(userID, identify string) error {

	return user.Table(user.TableName).Where("user_id", "=", userID).
		Update(command.Value{
			"admin_id": identify,
		})
}

// Add 新增用戶
func (user LineModel) Add(model LineModel) (int64, error) {
	var line string
	if model.Line == "" {
		line = "晶橙"
	} else {
		line = model.Line
	}

	id, err := user.Table(user.TableName).Insert(command.Value{
		"user_id":  model.UserID,
		"name":     model.Name,
		"avatar":   model.Avatar,
		"email":    model.Email,
		"identify": model.Identify,
		"friend":   model.Friend,
		// *****可以同時登入(暫時拿除)*****
		// "ip":       model.Ip,
		// *****可以同時登入(暫時拿除)*****
		"line":      line,
		"device":    model.Device,
		"is_modify": "no",
		"admin_id":  "",
	})
	if err != nil {
		return 0, errors.New("錯誤: 新增用戶發生問題，請重新操作")
	}
	return id, nil
}

// MapToLineModel map轉換[]ApplysignModeLineModell
func MapToLineModel(items []map[string]interface{}) []LineModel {
	var users = make([]LineModel, 0)
	for _, item := range items {
		var (
			lineModel LineModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &lineModel)

		users = append(users, lineModel)
	}
	return users
}

// MapToModel 設置LineModel
func (user LineModel) MapToModel(m map[string]interface{}) LineModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &user)

	return user
}

// lineModel.ID, _ = item["id"].(int64)
// lineModel.UserID, _ = item["user_id"].(string)
// lineModel.Name, _ = item["name"].(string)
// lineModel.Avatar, _ = item["avatar"].(string)
// lineModel.Email, _ = item["email"].(string)
// lineModel.Phone, _ = item["phone"].(string)
// lineModel.ExtEmail, _ = item["ext_email"].(string)
// lineModel.Identify, _ = item["identify"].(string)
// lineModel.Friend, _ = item["friend"].(string)
// *****可以同時登入(暫時拿除)*****
// lineModel.Ip, _ = item["ip"].(string)
// *****可以同時登入(暫時拿除)*****
// lineModel.Line, _ = item["line"].(string)
// lineModel.Device, _ = item["device"].(string)
// lineModel.IsModify, _ = item["is_modify"].(string)

// 自定義匯入簽到人員
// lineModel.ExtPassword, _ = item["ext_password"].(string)

// UpdateIP 更新用戶IP資訊
// func (user LineModel) UpdateIP(userID, ip string) error {
// 	err := user.Table(user.TableName).Where("user_id", "=", userID).
// 		Update(command.Value{
// 			"ip": ip,
// 		})
// 	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新用戶資訊發生問題，請重新操作")
// 	}

// 	return nil
// }

// user.ID, _ = m["id"].(int64)
// user.UserID, _ = m["user_id"].(string)
// user.Name, _ = m["name"].(string)
// user.Avatar, _ = m["avatar"].(string)
// user.Email, _ = m["email"].(string)
// user.Phone, _ = m["phone"].(string)
// user.ExtEmail, _ = m["ext_email"].(string)
// user.Identify, _ = m["identify"].(string)
// user.Friend, _ = m["friend"].(string)
// *****可以同時登入(暫時拿除)*****
// user.Ip, _ = m["ip"].(string)
// *****可以同時登入(暫時拿除)*****
// user.Line, _ = m["line"].(string)
// user.Device, _ = m["device"].(string)
// user.IsModify, _ = m["is_modify"].(string)

// 自定義匯入簽到人員
// user.ActivityID, _ = m["activity_id"].(string)
// user.ExtAccount, _ = m["ext_account"].(string)
// user.ExtPassword, _ = m["ex

// *****可以同時登入(暫時拿除)*****
// "ip"
// *****可以同時登入(暫時拿除)*****
// values = []string{model.Name, model.Avatar, model.Email, model.Friend, model.IsModify}
// *****可以同時登入(暫時拿除)*****
// model.Ip,
// *****可以同時登入(暫時拿除)*****

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
