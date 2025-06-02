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

// MenuModel 資料表欄位
type MenuModel struct {
	Base         `json:"-"`
	ID           int64       `json:"id"`
	ParentID     int64       `json:"parent_id"`
	Title        string      `json:"title"`         // menu title 欄位資料
	URL          string      `json:"url"`           // menu URL
	SidebarBtnL  string      `json:"sidebarbtnl"`   // SidebarBtnL欄位資料
	ChildrenList []MenuModel `json:"children_list"` // 放子選單
}

// EditMenuModel 資料表欄位
type EditMenuModel struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	Title       string `json:"title"`       // menu title 欄位資料
	URL         string `json:"url"`         // menu URL
	SidebarBtnL string `json:"sidebarbtnl"` // SidebarBtnL欄位資料
}

// DefaultMenuModel 預設MenuModel
func DefaultMenuModel() MenuModel {
	return MenuModel{Base: Base{TableName: config.MENU_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a MenuModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) MenuModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a MenuModel) SetDbConn(conn db.Connection) MenuModel {
// 	a.DbConn = conn
// 	return a
// }

// Find 尋找資料
func (a MenuModel) Find(id int64) (interface{}, error) {
	var (
		// item  map[string]interface{}
		// items = make([]map[string]interface{}, 0)
		data interface{}
		// err   error
	)

	if id != 0 {
		// 查詢單個overview資料
		item, err := a.Table(a.Base.TableName).
			Where("id", "=", id).First()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得菜單資訊，請重新查詢")
		}
		data = a.MapToModel(item)
	} else {
		// 查詢所有overview資料
		items, err := a.Table(a.Base.TableName).All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得菜單資訊，請重新查詢")
		}
		data = MapToMenuModel(items, 0)
	}
	return data, nil
}

// Add 新增菜單資料
func (a MenuModel) Add(model EditMenuModel) error {
	var (
		fields = []string{
			"parent_id",
			"title",
			"url",
			"sidebarbtnl",
		}
	)

	_, err := strconv.Atoi(model.ParentID)
	if err != nil {
		return errors.New("錯誤: parent_id資料發生問題，請輸入有效的ID")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	if _, err := a.Table(a.Base.TableName).Insert(FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 無法新增菜單資料，請重新操作")
	}
	return nil
}

//	command.Value{
//		"parent_id":   parentID,
//		"title":       model.Title,
//		"url":         model.URL,
//		"sidebarbtnl": model.SidebarBtnL,
//	}
//

// Update 更新菜單資料
func (a MenuModel) Update(model EditMenuModel) error {
	var (
		fieldValues = command.Value{"url": model.URL}
		fields      = []string{"title", "sidebarbtnl"}
		// values      = []string{model.Title, model.SidebarBtnL}
	)

	id, err := strconv.Atoi(model.ID)
	if err != nil {
		return errors.New("錯誤: ID發生問題，請輸入有效的ID")
	}
	if model.ParentID != "" {
		parentID, err := strconv.Atoi(model.ParentID)
		if err != nil {
			return errors.New("錯誤: parent_id資料發生問題，請輸入有效的ID")
		}
		fieldValues["parent_id"] = parentID
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 更新菜單資料
	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		if err := a.Table(a.Base.TableName).
			Where("id", "=", id).
			Update(fieldValues); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 無法更新菜單資料，請重新操作")
		}
	}
	return nil
}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

// MapToModel map轉換model
func (a MenuModel) MapToModel(item map[string]interface{}) MenuModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &a)

	return a
}

// MapToMenuModel map轉換[]MenuModel
func MapToMenuModel(items []map[string]interface{}, parentID int64) []MenuModel {
	menus := make([]MenuModel, 0)
	for j := 0; j < len(items); j++ {
		var (
			menu MenuModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(items[j])
		json.Unmarshal(b, &menu)

		if parentID == menu.ParentID {

			menu.ChildrenList = MapToMenuModel(items, menu.ID)

			menus = append(menus, menu)
		}
	}
	return menus
}

// menuParentID, _ := items[j]["parent_id"].(int64)
// menu.ID, _ = items[j]["id"].(int64)
// menu.ParentID = menuParentID
// menu.Title, _ = items[j]["title"].(string)
// menu.URL, _ = items[j]["url"].(string)
// menu.SidebarBtnL, _ = items[j]["sidebarbtnl"].(string)

// a.ID, _ = item["id"].(int64)
// a.ParentID = item["parent_id"].(int64)
// a.Title, _ = item["title"].(string)
// a.URL, _ = item["url"].(string)
// a.SidebarBtnL, _ = item["sidebarbtnl"].(string)
