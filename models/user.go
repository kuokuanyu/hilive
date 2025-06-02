package models

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// LoginUser 登入用戶資訊(合併users、line_users)
type LoginUser struct {
	UserID string `json:"user_id"` // 用戶ID
	Name   string `json:"name"`    // 用戶名稱
	Phone  string `json:"phone"`   // 用戶電話
	Email  string `json:"email"`   // 用戶信箱
	// Password string // 密碼
	Avatar string `json:"avatar"` // 用戶照片
	// Bind     string `json:"bind"`     // 綁定
	Cookie   string `json:"cookie"`   // 接受cookie
	Identify string `json:"identify"` // line用戶辨識資訊
	Friend   string `json:"friend"`   // line用戶是否加為好友
	Table    string `json:"table"`    // 資料表
	// *****可以同時登入(確定拿除)*****
	// Ip                string              `json:"-"`                   // ip
	// *****可以同時登入(確定拿除)*****
	MaxActivity       int64               `json:"max_activity"`        // 活動場次上限
	MaxActivityPeople int64               `json:"max_activity_people"` // 活動人數上限
	MaxGamePeople     int64               `json:"max_game_people"`     // 遊戲人數上限
	Permissions       []PermissionModel   `json:"permissions"`         // 用戶權限
	Activitys         []ActivityModel     `json:"activitys"`           // 活動資訊(包含活動權限)
	ActivityMenus     map[string][]string `json:"activity_menus"`      // 活動可用路徑菜單
	// Menus       []string          // 菜單

	LineBind  string `json:"line_bind"`  // LINE綁定
	FbBind    string `json:"fb_bind"`    // FB綁定
	GmailBind string `json:"gmail_bind"` // GMAIL綁定
}

// UserModel 資料表欄位
type UserModel struct {
	Base          `json:"-"`
	ID            int64               `json:"id"`             // ID
	UserID        string              `json:"user_id"`        // 用戶ID
	Name          string              `json:"name"`           // 用戶名稱
	Phone         string              `json:"phone"`          // 用戶電話
	Email         string              `json:"email"`          // 用戶信箱
	Password      string              `json:"password"`       // 密碼
	Avatar        string              `json:"avatar"`         // 用戶照片
	Permissions   []PermissionModel   `json:"permissions"`    // 用戶權限
	Activitys     []ActivityModel     `json:"activitys"`      // 活動資訊(包含活動權限)
	ActivityMenus map[string][]string `json:"activity_menus"` // 活動可用路徑菜單
	// CreatedAt   string            `json:"-"`           // 建立時間
	// UpdatedAt   string            `json:"-"`           // 更新時間
	// Bind   string `json:"bind"`   // 綁定
	Cookie string `json:"cookie"` // 接受cookie
	// *****可以同時登入(暫時拿除)*****
	// Ip                string `json:"-"`                   // ip
	// *****可以同時登入(暫時拿除)*****
	MaxActivity       int64 `json:"max_activity"`        // 活動場次上限
	MaxActivityPeople int64 `json:"max_activity_people"` // 活動人數上限
	MaxGamePeople     int64 `json:"max_game_people"`     // 遊戲人數上限

	// 官方帳號綁定
	LineID        string `json:"line_id"`        // line_id
	ChannelID     string `json:"channel_id"`     // channel_id
	ChannelSecret string `json:"channel_secret"` // channel_secret
	ChatbotSecret string `json:"chatbot_secret"` // chatbot_secret
	ChatbotToken  string `json:"chatbot_token"`  // chatbot_token
	// Roles       []RoleModel       `json:"-"`       // 角色

	CreatedAt string `json:"created_at" example:"0"`

	LineBind  string `json:"line_bind"`  // LINE綁定
	FbBind    string `json:"fb_bind"`    // FB綁定
	GmailBind string `json:"gmail_bind"` // GMAIL綁定
}

// EditUserModel 資料表欄位
type EditUserModel struct {
	// ID            string `json:"ID"`             // ID
	UserID            string `json:"user_id"`             // 用戶ID
	Name              string `json:"name"`                // 用戶名稱
	Phone             string `json:"phone"`               // 用戶電話
	Email             string `json:"email"`               // 用戶信箱
	Password          string `json:"password"`            // 密碼
	PasswordAgain     string `json:"password_again"`      // 再次輸入密碼
	Avatar            string `json:"avatar"`              // 用戶照片
	// Bind              string `json:"bind"`                // 綁定
	Cookie            string `json:"cookie"`              // 接受cookie
	MaxActivity       string `json:"max_activity"`        // 活動場次上限
	MaxActivityPeople string `json:"max_activity_people"` // 活動人數上限
	MaxGamePeople     string `json:"max_game_people"`     // 遊戲人數上限
	Permissions       string `json:"permissions"`         // 權限

	// 官方帳號綁定
	LineID        string `json:"line_id"`        // line_id
	ChannelID     string `json:"channel_id"`     // channel_id
	ChannelSecret string `json:"channel_secret"` // channel_secret
	ChatbotSecret string `json:"chatbot_secret"` // chatbot_secret
	ChatbotToken  string `json:"chatbot_token"`  // chatbot_token

	LineBind  string `json:"line_bind"`  // LINE綁定
	FbBind    string `json:"fb_bind"`    // FB綁定
	GmailBind string `json:"gmail_bind"` // GMAIL綁定
}

// DefaultUserModel 預設UserModel
func DefaultUserModel() UserModel {
	return UserModel{Base: Base{TableName: config.USERS_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (user UserModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) UserModel {
	user.DbConn = dbconn
	user.RedisConn = cacheconn
	user.MongoConn = mongoconn
	return user
}

// SetDbConn 設定connection
// func (user UserModel) SetDbConn(conn db.Connection) UserModel {
// 	user.DbConn = conn
// 	return user
// }

// // SetRedisConn 設定connection
// func (user UserModel) SetRedisConn(conn cache.Connection) UserModel {
// 	user.RedisConn = conn
// 	return user
// }

// // SetMongoConn 設定connection
// func (user UserModel) SetMongoConn(conn mongo.Connection) UserModel {
// 	user.MongoConn = conn
// 	return user
// }

// FindAll 尋找所有用戶資料
func (user UserModel) FindUsers() ([]UserModel, error) {
	items, err := user.Table(user.Base.TableName).
		Select("users.id", "users.user_id", "users.name", "users.phone", "users.email",
			"users.password", "users.avatar", "users.bind", "users.cookie", "users.created_at",
			// *****可以同時登入(暫時拿除)*****
			// "users.ip",
			// *****可以同時登入(暫時拿除)*****
			"users.max_activity", "users.max_activity_people", "users.max_game_people",
			"line_id", "channel_id", "channel_secret", "chatbot_secret", "chatbot_token",
			"users.line_bind","users.fb_bind","users.gmail_bind",

			//權限
			"permission.id as permission_id",
			"permission.permission", "permission.http_method", "permission.http_path").
		OrderBy("users.id", "asc", "user_permissions.permission_id", "asc").
		LeftJoin(command.Join{
			FieldA:    "user_permissions.user_id",
			FieldA1:   "users.user_id",
			Table:     "user_permissions",
			Operation: "="}).
		LeftJoin(command.Join{
			FieldA:    "permission.id",
			FieldA1:   "user_permissions.permission_id",
			Table:     "permission",
			Operation: "="}).All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得用戶資訊，請重新查詢")
	}

	return MapToUserModel(items), nil
}

// Find 尋找資料(從redis取得，如果沒有才執行資料表查詢)
// *****可以同時登入(暫時拿除)*****
// func (user UserModel) Find(isRedis, isActivityPermission bool, ip string, key string, value string) (UserModel, error) {
// *****可以同時登入(暫時拿除)*****
// *****可以同時登入(新)*****
func (user UserModel) Find(isRedis, isActivityPermission bool, key string, value string) (UserModel, error) {
	// *****可以同時登入(新)*****
	if isRedis {
		// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
		dataMap, err := user.RedisConn.HashGetAllCache(config.HILIVE_USERS_REDIS + value)
		if err != nil {
			return UserModel{}, errors.New("錯誤: 取得用戶快取資料發生問題")
		}
		id, _ := strconv.Atoi(dataMap["id"])
		user.ID = int64(id)
		user.UserID, _ = dataMap["user_id"]
		user.Name, _ = dataMap["name"]
		user.Avatar, _ = dataMap["avatar"]
		user.Phone, _ = dataMap["phone"]
		user.Email, _ = dataMap["email"]
		// user.Bind, _ = dataMap["bind"]
		user.Cookie, _ = dataMap["cookie"]
		user.CreatedAt, _ = dataMap["created_at"]
		// *****可以同時登入(暫時拿除)*****
		// user.Ip, _ = dataMap["ip"]
		// *****可以同時登入(暫時拿除)*****
		user.LineBind, _ = dataMap["line_bind"]
		user.FbBind, _ = dataMap["fb_bind"]
		user.GmailBind, _ = dataMap["gmail_bind"]

		// 活動場次上限
		maxActivity, _ := strconv.Atoi(dataMap["max_activity"])
		user.MaxActivity = int64(maxActivity)
		// 活動人數上限
		maxActivityPeople, _ := strconv.Atoi(dataMap["max_activity_people"])
		user.MaxActivityPeople = int64(maxActivityPeople)
		// 遊戲人數上限
		maxGamePeople, _ := strconv.Atoi(dataMap["max_game_people"])
		user.MaxGamePeople = int64(maxGamePeople)

		// 用戶權限
		var permissions = make([]PermissionModel, 0)
		permissionsStr, _ := dataMap["permissions"]
		// 解碼
		json.Unmarshal([]byte(permissionsStr), &permissions)
		user.Permissions = permissions

		// LINE資訊
		user.LineID, _ = dataMap["line_id"]
		user.ChannelID, _ = dataMap["channel_id"]
		user.ChannelSecret, _ = dataMap["channel_secret"]
		user.ChatbotSecret, _ = dataMap["chatbot_secret"]
		user.ChatbotToken, _ = dataMap["chatbot_token"]

		// 活動權限
		var activitys = make([]ActivityModel, 0)
		activitysStr, _ := dataMap["activitys"]
		// 解碼
		json.Unmarshal([]byte(activitysStr), &activitys)
		user.Activitys = activitys

		// 菜單
		var menus = make(map[string][]string, 0)
		menuStr, _ := dataMap["activity_menus"]
		// 解碼
		json.Unmarshal([]byte(menuStr), &menus)
		user.ActivityMenus = menus
	}

	if user.UserID == "" {
		items, err := user.Table(user.Base.TableName).
			Select("users.id", "users.user_id", "users.name", "users.phone", "users.email",
				"users.password", "users.avatar", "users.bind", "users.cookie", "users.created_at",
				// *****可以同時登入(暫時拿除)*****
				// "users.ip",
				// *****可以同時登入(暫時拿除)*****
				"users.max_activity", "users.max_activity_people", "users.max_game_people",
				"line_id", "channel_id", "channel_secret", "chatbot_secret", "chatbot_token", // LINE相關資訊
				"users.line_bind","users.fb_bind","users.gmail_bind",

				//權限
				"permission.id as permission_id",
				"permission.permission", "permission.http_method", "permission.http_path").
			Where(key, "=", value).
			OrderBy("users.id", "asc", "user_permissions.permission_id", "asc").
			LeftJoin(command.Join{
				FieldA:    "user_permissions.user_id",
				FieldA1:   "users.user_id",
				Table:     "user_permissions",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "permission.id",
				FieldA1:   "user_permissions.permission_id",
				Table:     "permission",
				Operation: "="}).
			All()
		if err != nil {
			return UserModel{}, errors.New("錯誤: 無法取得用戶資訊，請重新查詢")
		}

		user = user.MapToModel(items, isActivityPermission)
		// 將用戶資訊加入redis
		if isRedis {
			values := []interface{}{config.HILIVE_USERS_REDIS + user.UserID}
			values = append(values, "id", user.ID)
			values = append(values, "user_id", user.UserID)
			values = append(values, "name", user.Name)
			values = append(values, "phone", user.Phone)
			values = append(values, "email", user.Email)
			values = append(values, "avatar", user.Avatar)
			// values = append(values, "bind", user.Bind)
			values = append(values, "cookie", user.Cookie)
			values = append(values, "created_at", user.CreatedAt)
			// *****可以同時登入(暫時拿除)*****
			// values = append(values, "ip", ip)
			// *****可以同時登入(暫時拿除)*****
			values = append(values, "max_activity", user.MaxActivity)
			values = append(values, "max_activity_people", user.MaxActivityPeople)
			values = append(values, "max_game_people", user.MaxGamePeople)

			// LINE資訊
			values = append(values, "line_id", user.LineID)
			values = append(values, "channel_id", user.ChannelID)
			values = append(values, "channel_secret", user.ChannelSecret)
			values = append(values, "chatbot_secret", user.ChatbotSecret)
			values = append(values, "chatbot_token", user.ChatbotToken)

			// 綁定
			values = append(values, "line_bind", user.LineBind)
			values = append(values, "fb_bind", user.FbBind)
			values = append(values, "gmail_bind", user.GmailBind)

			// 用戶權限
			// json編碼
			permissions := utils.JSON(user.Permissions)
			values = append(values, "permissions", permissions)

			// 活動權限
			// json編碼
			activitys := utils.JSON(user.Activitys)
			values = append(values, "activitys", activitys)

			// 菜單
			// json編碼
			menus := utils.JSON(user.ActivityMenus)
			values = append(values, "activity_menus", menus)

			if err := user.RedisConn.HashMultiSetCache(values); err != nil {
				return user, errors.New("錯誤: 設置用戶快取資料發生問題")
			}
			// 設置過期時間
			// user.RedisConn.SetExpire(config.HILIVE_USERS_REDIS+user.UserID,
			// 	config.REDIS_EXPIRE)
		}
	}
	return user, nil
}

// FindPhoneAndEmail 尋找資料
func (user UserModel) FindPhoneAndEmail(phone, combination, email string) (UserModel, error) {
	items, err := user.Table(user.Base.TableName).
		Select("users.id", "users.user_id", "users.name", "users.phone", "users.email").
		// "users.password", "users.avatar", "users.bind", "users.cookie", "users.ip",
		// "users.max_activity", "users.max_activity_people", "users.max_game_people").

		//權限
		// "permission.id",
		// "permission.permission", "permission.http_method", "permission.http_path").
		Where("phone", "=", phone, combination).Where("email", "=", email).
		// OrderBy("users.id", "asc").
		// LeftJoin(command.Join{
		// 	FieldA:    "user_permissions.user_id",
		// 	FieldA1:    "users.user_id",
		// 	Table:     "user_permissions",
		// 	Operation: "="}).
		// LeftJoin(command.Join{
		// 	FieldA:    "permission.id",
		// 	FieldA1:    "user_permissions.permission_id",
		// 	Table:     "permission",
		// 	Operation: "="}).
		All()
	if err != nil {
		return UserModel{}, errors.New("錯誤: 無法取得用戶資訊，請重新查詢")
	}
	if items == nil {
		return UserModel{}, nil
	}
	return user.MapToModel(items, false), nil
}

// Add 增加用戶
func (user UserModel) Add(model EditUserModel) (UserModel, error) {
	var (
		// avatar                                        string
		// maxActivity, maxActivityPeople, maxGamePeople int
		err error

		fields = []string{
			"user_id",
			"name",
			"phone",
			"email",
			"password",
			"avatar",
			// "bind",
			"cookie",
			"max_activity",
			"max_activity_people",
			"max_game_people",
			"line_id",
			"channel_id",
			"channel_secret",
			"chatbot_secret",
			"chatbot_token",
			"line_bind",
			"fb_bind",
			"gmail_bind",
		}
	)

	// 欄位值判斷
	// 姓名
	if utf8.RuneCountInString(model.Name) > 20 {
		return user, errors.New("錯誤: 用戶名稱不能為空並且不能超過20個字元，請輸入有效的用戶名稱")
	}

	// 電話
	if len(model.Phone) > 2 {
		if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
			return user, errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	} else {
		return user, errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
	}

	// 密碼
	if model.Password != model.PasswordAgain {
		return user, errors.New("錯誤: 輸入密碼不一致，請輸入有效的密碼")
	}

	// 信箱
	if !strings.Contains(model.Email, "@") {
		return user, errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	// 判斷是否已註冊
	if userModel, err := DefaultUserModel().
		SetConn(user.DbConn, user.RedisConn, user.MongoConn).
		FindPhoneAndEmail(model.Phone, "or", model.Email); userModel.Phone != "" || err != nil {
		return user, errors.New("錯誤: 電話號碼或電子郵件地址已被註冊過，請輸入有效的手機號碼與電子郵件地址")
	}

	// 活動場次上限
	if model.MaxActivity == "" {
		model.MaxActivity = "1"
	} else {
		_, err = strconv.Atoi(model.MaxActivity)
		if err != nil {
			model.MaxActivity = "1"
		}
	}

	// 活動人數上限
	if model.MaxActivityPeople == "" {
		model.MaxActivityPeople = "5"
	} else {
		_, err = strconv.Atoi(model.MaxActivityPeople)
		if err != nil {
			model.MaxActivityPeople = "5"
		}
	}

	// 遊戲人數場次上限
	if model.MaxGamePeople == "" {
		model.MaxGamePeople = "5"
	} else {
		_, err = strconv.Atoi(model.MaxGamePeople)
		if err != nil {
			model.MaxGamePeople = "5"
		}
	}

	userID := utils.UUID(32) // 隨機產生user_id

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	data["user_id"] = userID
	data["password"] = EncodePassword([]byte(model.Password))

	id, err := user.Table(user.TableName).
		Insert(FilterFields(data, fields))
	if err != nil {
		return user, errors.New("錯誤: 新增用戶發生問題，請重新操作")
	}

	user.ID = id
	user.UserID = userID
	user.Name = model.Name
	user.Phone = model.Phone
	user.Email = model.Email
	user.Password = model.Password
	user.Avatar = model.Avatar
	// user.Bind = "no"
	user.Cookie = "no"
	user.LineBind = "no"
	user.FbBind = "no"
	user.GmailBind = "no"

	// 權限為空，給予預設權限
	if model.Permissions == "" {
		model.Permissions = "3"
	}
	if err = user.AddPermission(strings.Split(model.Permissions, ",")); err != nil {
		return user, err
	}

	// 新增一場測試活動
	if err := DefaultActivityModel().
		SetConn(user.DbConn, user.RedisConn, user.MongoConn).
		Add(false, EditActivityModel{
			UserID:           userID,
			ActivityID:       utils.UUID(20),
			ActivityName:     "測試活動",
			ActivityType:     "其他",
			MaxPeople:        model.MaxActivityPeople,
			People:           model.MaxActivityPeople,
			City:             "台中市",
			Town:             "南區",
			StartTime:        "2022-01-01T00:00",
			EndTime:          "2025-12-31T00:00",
			LoginRequired:    "close",
			PushMessage:      "close",
			PasswordRequired: "open",
			// EditTimes:        "0",
			Permissions: "2",
			Device:      "line,facebook,gmail,customize",
		}); err != nil {
		return user, err
	}

	// 建立用戶資料夾
	os.MkdirAll(config.STORE_PATH+"/"+userID, os.ModePerm)

	return user, nil
}

// Update 更新資料
func (user UserModel) Update(isRedis bool, model EditUserModel,
	key, value string) error {
	var (
		// 第一步：一定要加入的欄位（即使為空）
		fieldValues = map[string]interface{}{
			// 官方帳號(可能會清空數值)
			"line_id":        model.LineID,
			"channel_id":     model.ChannelID,
			"channel_secret": model.ChannelSecret,
			"chatbot_secret": model.ChatbotSecret,
			"chatbot_token":  model.ChatbotToken,
		}

		fields = []string{
			"name", "phone", "email",
			"avatar",  "cookie",
			"max_activity", "max_activity_people", "max_game_people",
			"line_bind","fb_bind","gmail_bind",
		}
	)

	// 欄位值判斷
	// 姓名
	if utf8.RuneCountInString(model.Name) > 20 {
		return errors.New("錯誤: 用戶姓名上限為20個字元，請輸入有效的用戶名型")
	}

	// 電話
	if model.Phone != "" {
		if len(model.Phone) > 2 {
			if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
				return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
			}
		} else {
			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	}
	if model.Email != "" && !strings.Contains(model.Email, "@") {
		return errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	// 密碼
	if model.Password != "" && model.PasswordAgain != "" {
		if model.Password != model.PasswordAgain {
			return errors.New("錯誤: 密碼不相符，請輸入有效的密碼")
		}
		fieldValues["password"] = EncodePassword([]byte(model.Password))
	}

	//綁定
	if model.LineBind != "" {
		if model.LineBind != "yes" && model.LineBind != "no" {
			return errors.New("錯誤: 綁定資料發生問題，請輸入有效的資料")
		}
	}
	if model.FbBind != "" {
		if model.FbBind != "yes" && model.FbBind != "no" {
			return errors.New("錯誤: 綁定資料發生問題，請輸入有效的資料")
		}
	}
	if model.GmailBind != "" {
		if model.GmailBind != "yes" && model.GmailBind != "no" {
			return errors.New("錯誤: 綁定資料發生問題，請輸入有效的資料")
		}
	}

	// cookie
	if model.Cookie != "" {
		if model.Cookie != "yes" && model.Cookie != "no" {
			return errors.New("錯誤: cookie訊息資料發生問題，請輸入有效的資料")
		}
	}

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

	err := user.Table(user.Base.TableName).
		Where(key, "=", value).Update(fieldValues)
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新用戶資訊發生問題，請重新操作")
	}

	// 新增權限
	if model.Permissions != "" {
		// id, _ := strconv.Atoi(model.ID)
		// user.ID = int64(id)

		user.UserID = model.UserID
		if err = user.AddPermission(strings.Split(model.Permissions, ",")); err != nil {
			return err
		}
	}

	if isRedis {
		// 清除用戶redis資訊
		user.RedisConn.DelCache(config.HILIVE_USERS_REDIS + model.UserID)
	}
	return nil
}

// AddPermission 增加權限(多個權限)
func (user UserModel) AddPermission(ids []string) error {
	if len(ids) > 0 {
		// sql := db.Table(config.USER_PERMISSIONS_TABLE).WithConn(user.DbConn)

		// 先刪除原有的角色權限資料
		if err := db.Table(config.USER_PERMISSIONS_TABLE).
			WithConn(user.DbConn).Where("user_id", "=", user.UserID).
			Delete(); err != nil && err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
			return errors.New("錯誤: 刪除用戶權限資料發生問題，請重新操作")
		}

		// 加入新的權限資料
		for _, id := range ids {
			if _, err := db.Table(config.USER_PERMISSIONS_TABLE).
				WithConn(user.DbConn).Insert(command.Value{
				"permission_id": id, "user_id": user.UserID,
			}); err != nil {
				return errors.New("錯誤: 新增用戶權限資料發生問題，請重新操作")
			}
		}
	}

	return nil
}

// MapToModel 設置UserModel
func (user UserModel) MapToModel(items []map[string]interface{}, isActivityPermission bool) UserModel {
	// fmt.Println("user", len(items), items)

	if len(items) > 0 {
		// json解碼，轉換成strcut
		b, _ := json.Marshal(items[0])
		json.Unmarshal(b, &user)

		if !strings.Contains(user.Avatar, "system") {
			user.Avatar = "/admin/uploads/" + user.UserID + "/" + user.Avatar
		}
	}

	// 權限
	user.Permissions = MapToPermissionModelModel(items) // 用戶權限

	// 是否需要活動權限資料
	if isActivityPermission {
		user.Activitys, _ = DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			FindActivityPermissions("activity.user_id", user.UserID) // 活動權限

		// var canUseAll bool
		activityMenus := make(map[string][]string, 0) // 活動可用路徑菜單
		userURLs := []string{}                        // 用戶權限的可用路徑

		// 用戶權限
		for _, permission := range user.Permissions {
			// if len(permission.HTTPPath) > 0 && (permission.HTTPPath[0] == "*" ||
			// 	permission.HTTPPath[0] == "admin") {
			// 		userURLs = []string{permission.HTTPPath[0]} // 所有權限

			// 		// 將所有活動的權限設為所有權限
			// 		for _, activity := range user.Activitys {
			// 			menus[activity.ActivityID] = userURLs
			// 		}

			// 		// 可使用所有權限
			// 		canUseAll = true
			// 	break
			// }

			// 將可用路徑加入userURLs中
			for _, path := range permission.HTTPPath {
				if !utils.InArray(userURLs, path) && path != "" {
					userURLs = append(userURLs, path)
				}
			}
		}

		// fmt.Println("用戶權限: ", userURLs)
		// 活動權限
		for _, activity := range user.Activitys {
			activityMenus[activity.ActivityID] = userURLs // 加入用戶權限
			// fmt.Println("用戶權限: ", activityMenus[activity.ActivityID], len(activityMenus[activity.ActivityID]))

			// 活動權限路徑處理
			for _, permission := range activity.Permissions {
				for _, path := range permission.HTTPPath {
					// fmt.Println("path: ", path)
					if !utils.InArray(activityMenus[activity.ActivityID], path) &&
						path != "" {
						activityMenus[activity.ActivityID] = append(activityMenus[activity.ActivityID], path)
					}
				}
			}
			// fmt.Println("活動權限: ", activityMenus[activity.ActivityID], len(activityMenus[activity.ActivityID]))
		}

		// 最後判斷是否有最大權限
		for activity, menu := range activityMenus {
			// fmt.Println(activity, menu)
			if utils.InArray(menu, "admin") {
				// 管理員權限
				activityMenus[activity] = []string{"admin"}
			} else if utils.InArray(menu, "*") {
				// 活動所有權限
				activityMenus[activity] = []string{"*"}
			}
		}

		// 沒有所有權限，判斷活動權限
		// if !canUseAll {

		// }
		user.ActivityMenus = activityMenus
	}
	return user
}

// MapToUserModel map轉換[]UserModel
func MapToUserModel(items []map[string]interface{}) []UserModel {
	var (
		users = make([]UserModel, 0)
		user  = UserModel{}
	)

	for i := 0; i < len(items); i++ {
		userID, _ := items[i]["user_id"].(string)

		// 第一筆資料不用比較
		if i == 0 {

			// json解碼，轉換成strcut
			b, _ := json.Marshal(items[i])
			json.Unmarshal(b, &user)

			if !strings.Contains(user.Avatar, "system") {
				user.Avatar = "/admin/uploads/" + user.UserID + "/" + user.Avatar
			}

			user.Permissions = append(user.Permissions, DefaultPermissionModel().MapToModel(items[i]))
		} else if user.UserID == userID {
			// 還是同一個用戶的權限資料，將權限資料加入用戶中
			user.Permissions = append(user.Permissions, DefaultPermissionModel().MapToModel(items[i]))
		} else if user.UserID != userID {
			// 不是上一個用戶的權限資料，將上一個用戶資料放入users中
			users = append(users, user)

			// 清空user參數資料
			user = UserModel{}

			// json解碼，轉換成strcut
			b, _ := json.Marshal(items[i])
			json.Unmarshal(b, &user)

			if !strings.Contains(user.Avatar, "system") {
				user.Avatar = "/admin/uploads/" + user.UserID + "/" + user.Avatar
			}

			user.Permissions = append(user.Permissions, DefaultPermissionModel().MapToModel(items[i]))
		}

		if i == len(items)-1 {
			// 最後一筆資料
			users = append(users, user)
		}
	}
	return users
}

// UpdatePassword 更新密碼
func (user UserModel) UpdatePassword(password string) UserModel {
	user.Base.Table(user.Base.TableName).Where("id", "=", user.ID).Update(command.Value{
		"password": password,
	})
	user.Password = password
	return user
}

// if model.UserID == "" {
// 	// 查詢用戶ID
// 	userModel, _ := DefaultUserModel().SetDbConn(user.DbConn).
// 		SetRedisConn(user.RedisConn).Find(false, "", "users.id", model.ID)
// 	model.UserID = userModel.UserID
// }

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

// values = []string{model.Name, model.Phone, model.Email,
// 	model.Avatar, model.Bind, model.Cookie,
// 	model.MaxActivity, model.MaxActivityPeople, model.MaxGamePeople,
// }

// *****可以同時登入(暫時拿除)*****
// UpdateIP 更新用戶IP資訊
// func (user UserModel) UpdateIP(userID, ip string) error {
// 	err := user.Table(user.TableName).Where("user_id", "=", userID).
// 		Update(command.Value{
// 			"ip": ip,
// 		})
// 	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新用戶資訊發生問題，請重新操作")
// 	}

// 	return nil
// }
// *****可以同時登入(暫時拿除)*****

// user.Bind, _ = items[i]["bind"].(string)
// user.Cookie, _ = items[i]["cookie"].(string)
// user.CreatedAt, _ = items[i]["created_at"].(string)
// *****可以同時登入(暫時拿除)*****
// user.Ip, _ = items[i]["ip"].(string)
// *****可以同時登入(暫時拿除)*****
// user.MaxActivity, _ = items[i]["max_activity"].(int64)
// user.MaxActivityPeople, _ = items[i]["max_activity_people"].(int64)
// user.MaxGamePeople, _ = items[i]["max_game_people"].(int64)
// user.LineID, _ = items[i]["line_id"].(string)
// user.ChannelID, _ = items[i]["channel_id"].(string)
// user.ChannelSecret, _ = items[i]["channel_secret"].(string)
// user.ChatbotSecret, _ = items[i]["chatbot_secret"].(string)
// user.ChatbotToken, _ = items[i]["chatbot_token"].(string)

// user.ID, _ = items[i]["id"].(int64)
// user.UserID, _ = items[i]["user_id"].(string)
// user.Name, _ = items[i]["name"].(string)
// user.Phone, _ = items[i]["phone"].(string)
// user.Email, _ = items[i]["email"].(string)
// user.Password, _ = items[i]["password"].(string)
// *****舊*****
// user.Avatar, _ = items[i]["avatar"].(string)
// *****舊*****
// *****新*****
// avatar, _ := items[i]["avatar"].(string)
// *****新*****

// user.Bind, _ = items[i]["bind"].(string)
// user.Cookie, _ = items[i]["cookie"].(string)
// user.CreatedAt, _ = items[i]["created_at"].(string)
// *****可以同時登入(暫時拿除)*****
// user.Ip, _ = items[i]["ip"].(string)
// *****可以同時登入(暫時拿除)*****
// user.MaxActivity, _ = items[i]["max_activity"].(int64)
// user.MaxActivityPeople, _ = items[i]["max_activity_people"].(int64)
// user.MaxGamePeople, _ = items[i]["max_game_people"].(int64)
// user.LineID, _ = items[i]["line_id"].(string)
// user.ChannelID, _ = items[i]["channel_id"].(string)
// user.ChannelSecret, _ = items[i]["channel_secret"].(string)
// user.ChatbotSecret, _ = items[i]["chatbot_secret"].(string)
// user.ChatbotToken, _ = items[i]["chatbot_token"].(string)

// user.ID, _ = items[i]["id"].(int64)
// user.UserID, _ = items[i]["user_id"].(string)
// user.Name, _ = items[i]["name"].(string)
// user.Phone, _ = items[i]["phone"].(string)
// user.Email, _ = items[i]["email"].(string)
// user.Password, _ = items[i]["password"].(string)

// *****舊*****
// user.Avatar, _ = items[i]["avatar"].(string)
// *****舊*****
// *****新*****
// avatar, _ := items[i]["avatar"].(string)
// *****新*****

// user.Bind, _ = items[0]["bind"].(string)
// user.Cookie, _ = items[0]["cookie"].(string)
// user.CreatedAt, _ = items[0]["created_at"].(string)
// *****可以同時登入(暫時拿除)*****
// user.Ip, _ = items[0]["ip"].(string)
// *****可以同時登入(暫時拿除)*****
// user.MaxActivity, _ = items[0]["max_activity"].(int64)
// user.MaxActivityPeople, _ = items[0]["max_activity_people"].(int64)
// user.MaxGamePeople, _ = items[0]["max_game_people"].(int64)
// user.LineID, _ = items[0]["line_id"].(string)
// user.ChannelID, _ = items[0]["channel_id"].(string)
// user.ChannelSecret, _ = items[0]["channel_secret"].(string)
// user.ChatbotSecret, _ = items[0]["chatbot_secret"].(string)
// user.ChatbotToken, _ = items[0]["chatbot_token"].(string)

// user.ID, _ = items[0]["id"].(int64)
// user.UserID, _ = items[0]["user_id"].(string)
// user.Name, _ = items[0]["name"].(string)
// user.Phone, _ = items[0]["phone"].(string)
// user.Email, _ = items[0]["email"].(string)
// user.Password, _ = items[0]["password"].(string)
// *****舊*****
// user.Avatar, _ = items[0]["avatar"].(string)
// *****舊*****
// *****新*****
// avatar, _ := items[0]["avatar"].(string)
// *****新*****
// AddRole 增加角色
// func (user UserModel) AddRole(id string) (int64, error) {
// 	if checkRole, _ := user.Table(config.ROLE_USERS_TABLE).Where("role_id", "=", id).
// 		Where("user_id", "=", user.ID).First(); id != "" && checkRole == nil {
// 		user.Table(config.ROLE_USERS_TABLE).Insert(command.Value{
// 			"role_id": id, "user_id": user.ID,
// 		})
// 	}
// 	return 0, nil
// }

// RoleModel 資料表欄位
// type RoleModel struct {
// 	Base
// 	ID        int64  // ID
// 	Name      string // 角色名稱
// 	Slug      string // 角色slogan
// 	CreatedAt string // 建立時間
// 	UpdatedAt string // 更新時間
// }
// // DefaultRoleModel 預設RoleModel
// func DefaultRoleModel() RoleModel {
// 	return RoleModel{Base: Base{TableName: config.ROLES_TABLE}}
// }

// MapToModel 設置RoleModel
// func (role RoleModel) MapToModel(m map[string]interface{}) RoleModel {
// 	role.ID, _ = m["id"].(int64)
// 	role.Name, _ = m["name"].(string)
// 	role.Slug, _ = m["slug"].(string)
// 	role.CreatedAt, _ = m["created_at"].(string)
// 	role.UpdatedAt, _ = m["updated_at"].(string)
// 	return role
// }

// GetRoles 取得角色
// func (user UserModel) GetRoles() UserModel {
// 	roleModel, _ := user.Base.Table(config.ROLE_USERS_TABLE).
// 		LeftJoin(command.Join{
// 			FieldA:    config.ROLES_TABLE + ".id",
// 			FieldA1:    config.ROLE_USERS_TABLE + ".role_id",
// 			Table:     config.ROLES_TABLE,
// 			Operation: "=",
// 		}).
// 		Where("user_id", "=", user.ID).
// 		Select(config.ROLES_TABLE+".id", config.ROLES_TABLE+".name", config.ROLES_TABLE+".slug",
// 			config.ROLES_TABLE+".created_at", config.ROLES_TABLE+".updated_at").All()
// 	for _, role := range roleModel {
// 		user.Roles = append(user.Roles, DefaultRoleModel().MapToModel(role))
// 	}
// 	return user
// }

// GetMenus 取得可用menu
// func (user UserModel) GetMenus() UserModel {
// 	var (
// 		menuIDsModel []map[string]interface{}
// 		menuIDs      []int64
// 	)
// 	if user.IsSuperAdmin() {
// 		menuIDsModel, _ = user.Base.Table(config.ROLE_MENU_TABLE).
// 			LeftJoin(command.Join{
// 				FieldA:    config.MENU_TABLE + ".id",
// 				FieldA1:    config.ROLE_MENU_TABLE + ".menu_id",
// 				Table:     config.MENU_TABLE,
// 				Operation: "=",
// 			}).
// 			Select("menu_id", "parent_id").All()
// 	} else {
// 		rolesID := user.GetAllRole()
// 		if len(rolesID) > 0 {
// 			menuIDsModel, _ = user.Base.Table(config.ROLE_MENU_TABLE).
// 				LeftJoin(command.Join{
// 					FieldA:    config.MENU_TABLE + ".id",
// 					FieldA1:    config.ROLE_MENU_TABLE + ".menu_id",
// 					Table:     config.MENU_TABLE,
// 					Operation: "=",
// 				}).
// 				WhereIn(config.ROLE_MENU_TABLE+".role_id", rolesID).Select("menu_id", "parent_id").All()
// 		}
// 	}
// 	for _, mid := range menuIDsModel {
// 		menuIDs = append(menuIDs, mid["menu_id"].(int64))
// 	}
// 	user.MenuIDs = menuIDs
// 	return user
// }

// GetAllRole 用戶所有role_id
// func (user UserModel) GetAllRole() []interface{} {
// 	var ids = make([]interface{}, len(user.Roles))
// 	for key, role := range user.Roles {
// 		ids[key] = role.ID
// 	}
// 	return ids
// }

// if model.Name == "" || model.Phone == "" || model.Email == "" ||
// 	model.Password == "" || model.PasswordAgain == "" {
// 	return user, errors.New("錯誤: 欄位資訊都不能為空，請輸入有效的欄位資訊")
// }

// hash, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
// if err != nil {
// 	return user, errors.New("錯誤: 加密發生錯誤，請重新註冊用戶")
// }

// 頭像
// if model.Avatar == "" {
// 	avatar = config.UPLOAD_SYSTEM_URL + "img-user-pic.png"
// }

// command.Value{
// 	"user_id":             userID,
// 	"name":                model.Name,
// 	"phone":               model.Phone,
// 	"email":               model.Email,
// 	"password":            EncodePassword([]byte(model.Password)),
// 	"avatar":              model.Avatar,
// 	"bind":                "no",
// 	"cookie":              "no",
// 	"max_activity":        maxActivity,
// 	"max_activity_people": maxActivityPeople,
// 	"max_game_people":     maxGamePeople,
// 	"line_id":             model.LineID,
// 	"channel_id":          model.ChannelID,
// 	"channel_secret":      model.ChannelSecret,
// 	"chatbot_secret":      model.ChatbotSecret,
// 	"chatbot_token":       model.ChatbotToken,
// }
