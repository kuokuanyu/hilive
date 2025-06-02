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
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// PermissionModel 資料表欄位
type PermissionModel struct {
	Base
	ID         int64    `json:"permission_id"` // ID
	Permission string   `json:"permission"`    // 權限名稱
	HTTPMethod string   `json:"http_method"`   // GET、POST、PUT、DELETE
	HTTPPath   []string `json:"http_path"`     // 可訪問路徑
	// CreatedAt  string   // 建立時間
	// UpdatedAt  string   // 更新時間
}

// EditPermissionModel 資料表欄位
type EditPermissionModel struct {
	ID         string `json:"id"`          // ID
	Permission string `json:"permission"`  // 權限名稱
	HTTPMethod string `json:"http_method"` // GET、POST、PUT、DELETE
	HTTPPath   string `json:"http_path"`   // 可訪問路徑
}

// DefaultPermissionModel 預設PermissionModel
func DefaultPermissionModel() PermissionModel {
	return PermissionModel{Base: Base{TableName: config.PERMISSION_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (p PermissionModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) PermissionModel {
	p.DbConn = dbconn
	p.RedisConn = cacheconn
	p.MongoConn = mongoconn
	return p
}

// SetDbConn 設定connection
// func (p PermissionModel) SetDbConn(conn db.Connection) PermissionModel {
// 	p.DbConn = conn
// 	return p
// }

// Find 尋找資料
func (p PermissionModel) Find(id int64) (interface{}, error) {
	var (
		data interface{}
	)

	if id != 0 {
		// 查詢單個permission資料
		item, err := p.Table(p.Base.TableName).
			Select("permission.id as permission_id", "permission", "http_method", "http_path").
			Where("permission.id", "=", id).First()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得權限資訊，請重新查詢")
		}
		data = p.MapToModel(item)
	} else if id == 0 {
		// 查詢所有permission資料
		items, err := p.Table(p.Base.TableName).
			Select("permission.id as permission_id", "permission", "http_method", "http_path").
			All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得權限資訊，請重新查詢")
		}
		data = MapToPermissionModelModel(items)
	}

	return data, nil
}

// Add 新增權限資料
func (p PermissionModel) Add(model EditPermissionModel) error {
	var (
		fields = []string{
			"permission",
			"http_method",
			"http_path",
		}
	)

	// 欄位值判斷
	if utf8.RuneCountInString(model.Permission) > 20 {
		return errors.New("錯誤: 權限名稱上限為20個字元，請輸入有效的權限名稱")
	}
	if model.HTTPPath != "*" && model.HTTPPath != "admin" {
		// model.HTTPPath += ",/admin/user,/admin/activity"
	} else {
		// 所有權限
		model.HTTPMethod = "POST,PUT,DELETE"
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	if _, err := p.Table(p.Base.TableName).Insert(FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 無法新增權限資料，請重新操作")
	}
	return nil
}

// Update 更新權限資料
func (p PermissionModel) Update(model EditPermissionModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"permission", "http_method", "http_path"}
		// values      = []string{model.Permission, model.HTTPMethod, model.HTTPPath}
	)

	// 欄位值判斷
	id, err := strconv.Atoi(model.ID)
	if err != nil {
		return errors.New("錯誤: ID發生問題，請輸入有效的ID")
	}
	if utf8.RuneCountInString(model.Permission) > 20 {
		return errors.New("錯誤: 權限名稱上限為20個字元，請輸入有效的權限名稱")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 更新權限資料
	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			if key == "http_path" {
				if val != "*" && val != "admin" {
					// value += ",/admin/user"
				} else {
					// 所有權限
					fieldValues["http_method"] = "POST,PUT,DELETE"
				}
			}

			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		if err := p.Table(p.Base.TableName).
			Where("id", "=", id).
			Update(fieldValues); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 無法更新權限資料，請重新操作")
		}
	}
	return nil
}

// MapToPermissionModelModel map轉換[]PermissionModel
func MapToPermissionModelModel(items []map[string]interface{}) []PermissionModel {
	var permissions = make([]PermissionModel, 0)
	for _, item := range items {
		var (
			permission PermissionModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &permission)

		paths, _ := item["http_path"].(string)
		if paths != "" {
			permission.HTTPPath = strings.Split(paths, ",")
		} else {
			permission.HTTPPath = []string{""}
		}

		permissions = append(permissions, permission)
	}
	return permissions
}

// MapToModel map轉換model
func (permission PermissionModel) MapToModel(m map[string]interface{}) PermissionModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &permission)

	paths, _ := m["http_path"].(string)
	if paths != "" {
		permission.HTTPPath = strings.Split(paths, ",")
	} else {
		permission.HTTPPath = []string{""}
	}

	return permission
}

// CheckPermission 透過url判斷活動是否有權限訪問頁面
func (user UserModel) CheckPermission(activityID, path string) bool {
	// 路徑判斷處理
	if path == "" {
		return false
	}
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// 登出
	logoutCheck, _ := regexp.Compile("/admin/logout" + "(.*?)")
	if logoutCheck.MatchString(path) {
		return true
	}

	urls := make([]string, 0) // 可用路徑菜單
	if strings.Contains(path, "/admin/manager") ||
		strings.Contains(path, "/admin/permission") ||
		strings.Contains(path, "/admin/menu") ||
		strings.Contains(path, "/admin/overview") ||
		strings.Contains(path, "/admin/log") ||
		strings.Contains(path, "/admin/error_log") {
		// 活動ID為空，管理員頁面，取得用戶權限
		for _, permission := range user.Permissions {
			for _, httppath := range permission.HTTPPath {
				if !utils.InArray(urls, httppath) && httppath != "" {
					urls = append(urls, httppath)
				}
			}
		}

		// 判斷是否有最大權限
		if utils.InArray(urls, "admin") {
			// 管理員權限
			urls = []string{"admin"}
		} else if utils.InArray(urls, "*") {
			// 活動所有權限
			urls = []string{"*"}
		}
	} else if activityID != "" {
		// 活動可用路徑菜單
		urls = user.ActivityMenus[activityID]
	}

	// 管理員權限或活動所有權限
	if len(urls) > 0 {
		if urls[0] == "*" {
			// 只有活動所有權限
			if strings.Contains(path, "/admin/manager") ||
				strings.Contains(path, "/admin/permission") ||
				strings.Contains(path, "/admin/menu") ||
				strings.Contains(path, "/admin/overview") ||
				strings.Contains(path, "/admin/log") ||
				strings.Contains(path, "/admin/error_log") {
				// fmt.Println("*權限，沒有後端權限管理功能")
				return false
			}
			return true
		} else if urls[0] == "admin" {
			// 管理員權限
			return true
		}
	} else {
		// 無任何權限
		return false
	}

	// fmt.Println("path: ", path, "menu: ", activityMenus)
	// 檢查可用路徑裡是否有權限
	for _, url := range urls {
		// if menu == path {
		// 	return true
		// }
		if strings.Contains(path, url) && url != "" {
			return true
		}
	}
	return false
}

// permission.ID, _ = item["permission_id"].(int64)
// permission.Permission, _ = item["permission"].(string)

// permission.HTTPMethod, _ = item["http_method"].(string)
// methods, _ := item["http_method"].(string)
// if methods != "" {
// 	permission.HTTPMethod = strings.Split(methods, ",")
// } else {
// 	permission.HTTPMethod = []string{""}
// }

// permission.ID, _ = m["permission_id"].(int64)
// permission.Permission, _ = m["permission"].(string)

// permission.HTTPMethod, _ = m["http_method"].(string)
// methods, _ := m["http_method"].(string)
// if methods != "" {
// 	permission.HTTPMethod = strings.Split(methods, ",")
// } else {
// 	permission.HTTPMethod = []string{""}
// }

// path, params := GetURLParam(path)
// for key, value := range formParams {
// 	if len(value) > 0 {
// 		params.Add(key, value[0])
// 	}
// }

// GetPermissions 取得權限
// func (user UserModel) GetPermissions() UserModel {
// 	var permissions = make([]map[string]interface{}, 0)
// 	roleIDs := user.GetAllRole()
// 	if len(roleIDs) > 0 {
// 		permissions, _ = user.Base.Table(config.ROLE_PERMISSIONS_TABLE).
// 			LeftJoin(command.Join{
// 				FieldA:    config.PERMISSIONS_TABLE + ".id",
// 				FieldB:    config.ROLE_PERMISSIONS_TABLE + ".permission_id",
// 				Table:     config.PERMISSIONS_TABLE,
// 				Operation: "=",
// 			}).
// 			WhereIn("role_id", roleIDs).Select(config.PERMISSIONS_TABLE+".http_method", config.PERMISSIONS_TABLE+".http_path",
// 			config.PERMISSIONS_TABLE+".id", config.PERMISSIONS_TABLE+".name", config.PERMISSIONS_TABLE+".slug",
// 			config.PERMISSIONS_TABLE+".created_at", config.PERMISSIONS_TABLE+".updated_at").All()
// 	}
// 	userPermissions, _ := user.Base.Table(config.USER_PERMISSIONS_TABLE).
// 		LeftJoin(command.Join{
// 			FieldA:    config.PERMISSIONS_TABLE + ".id",
// 			FieldB:    config.USER_PERMISSIONS_TABLE + ".permission_id",
// 			Table:     config.PERMISSIONS_TABLE,
// 			Operation: "=",
// 		}).
// 		Where("user_id", "=", user.ID).Select(config.PERMISSIONS_TABLE+".http_method", config.PERMISSIONS_TABLE+".http_path",
// 		config.PERMISSIONS_TABLE+".id", config.PERMISSIONS_TABLE+".name", config.PERMISSIONS_TABLE+".slug",
// 		config.PERMISSIONS_TABLE+".created_at", config.PERMISSIONS_TABLE+".updated_at").All()
// 	permissions = append(permissions, userPermissions...)

// 	for i := 0; i < len(permissions); i++ {
// 		exist := false
// 		for j := 0; j < len(user.Permissions); j++ {
// 			if user.Permissions[j].ID == permissions[i]["id"] {
// 				exist = true
// 				break
// 			}
// 		}
// 		if exist {
// 			continue
// 		}
// 		user.Permissions = append(user.Permissions, DefaultPermissionModel().MapToModel(permissions[i]))
// 	}
// 	return user
// }

// IsSuperAdmin 是否為超級管理員
// func (user UserModel) IsSuperAdmin() bool {
// 	for _, permission := range user.Permissions {
// 		if len(permission.HTTPPath) > 0 && permission.HTTPPath[0] == "*" && permission.HTTPMethod[0] == "" {
// 			return true
// 		}
// 	}
// 	return false
// }

// GetCheckPermission 檢查用戶權限
// func (user UserModel) GetCheckPermission(path, method string) string {
// 	if !user.CheckPermission(path, method, url.Values{}) {
// 		return ""
// 	}
// 	return path
//

// command.Value{
// 	"permission":  model.Permission,
// 	"http_method": model.HTTPMethod,
// 	// "http_method": model.HTTPMethod + ",GET",
// 	"http_path": model.HTTPPath,
// }

// if model.Permission != "" {
// 	fieldValues["permission"] = model.Permission
// }
// if model.HTTPMethod != "" {
// 	// fieldValues["http_method"] = model.HTTPMethod + ",GET"
// 	fieldValues["http_method"] = model.HTTPMethod
// }
// if model.HTTPPath != "" {
// 	if model.HTTPPath != "*" && model.HTTPPath != "admin" {
// 		// model.HTTPPath += ",/admin/user,/admin/activity"
// 	} else {
// 		// 所有權限
// 		fieldValues["http_method"] = "POST,PUT,DELETE"
// 	}
// 	fieldValues["http_path"] = model.HTTPPath
// }

// 更新權限資料
// for i, value := range values {
// 	if value != "" {
// 		if fields[i] == "http_path" {
// 			if value != "*" && value != "admin" {
// 				// value += ",/admin/user"
// 			} else {
// 				// 所有權限
// 				fieldValues["http_method"] = "POST,PUT,DELETE"
// 			}
// 		}

// 		fieldValues[fields[i]] = value
// 	}
// }
