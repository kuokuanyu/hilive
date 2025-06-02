package models

import (
	"encoding/json"
	"errors"
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

// GuestModel 資料表欄位
type GuestModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Introduce  string `json:"introduce"`
	Detail     string `json:"detail"`
	GuestOrder int64  `json:"guest_order"`
}

// EditGuestModel 資料表欄位
type EditGuestModel struct {
	ID         string `json:"id"`
	ActivityID string `json:"activity_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Introduce  string `json:"introduce"`
	Detail     string `json:"detail"`
	GuestOrder string `json:"guest_order"`
}

// EditGuestSettingModel 資料表欄位
type EditGuestSettingModel struct {
	ActivityID      string `json:"activity_id"`
	GuestTitle      string `json:"guest_title"`
	GuestBackground string `json:"guest_background"`
}

// DefaultGuestModel 預設GuestModel
func DefaultGuestModel() GuestModel {
	return GuestModel{Base: Base{TableName: config.ACTIVITY_GUEST_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (g GuestModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GuestModel {
	g.DbConn = dbconn
	g.RedisConn = cacheconn
	g.MongoConn = mongoconn
	return g
}

// SetDbConn 設定connection
// func (g GuestModel) SetDbConn(conn db.Connection) GuestModel {
// 	g.DbConn = conn
// 	return g
// }

// Find 尋找資料
func (g GuestModel) Find(field string, value interface{}) ([]GuestModel, error) {
	items, err := g.Table(g.Base.TableName).
		Select("activity_guest.id", "activity_guest.activity_id", "avatar",
			"name", "introduce", "detail", "guest_order",

			// 活動
			"activity.user_id").
		Where(field, "=", value).
		LeftJoin(command.Join{
			FieldA:    "activity_guest.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity",
			Operation: "="}).
		OrderBy("guest_order", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得嘉賓資訊，請重新查詢")
	}
	return g.MapToModel(items), nil
}

// Add 增加資料
func (g GuestModel) Add(model EditGuestModel) error {
	var (
		fields = []string{
			"activity_id",
			"avatar",
			"name",
			"introduce",
			"detail",
			"guest_order",
		}
	)

	if utf8.RuneCountInString(model.Name) > 20 {
		return errors.New("錯誤: 姓名上限為20個字元，請輸入有效的姓名")
	}

	if utf8.RuneCountInString(model.Introduce) > 20 {
		return errors.New("錯誤: 簡介上限為20個字元，請輸入有效的簡介資料")
	}

	if utf8.RuneCountInString(model.Detail) > 200 {
		return errors.New("詳情上限為200個字元，請輸入有效的詳情資料")
	}

	if _, err := strconv.Atoi(model.GuestOrder); err != nil {
		return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	_, err := g.Table(g.TableName).Insert(FilterFields(data, fields))
	return err
}

// Update 更新資料
func (g GuestModel) Update(model EditGuestModel) error {
	var (
		fieldValues = command.Value{
			"name":        model.Name,
			"guest_order": model.GuestOrder,
			"introduce":   model.Introduce,
			"detail":      model.Detail}
		fields = []string{"avatar"}
		values = []string{model.Avatar}
	)

	if utf8.RuneCountInString(model.Name) > 20 {
		return errors.New("錯誤: 姓名上限為20個字元，請輸入有效的姓名")
	}

	if utf8.RuneCountInString(model.Introduce) > 20 {
		return errors.New("錯誤: 簡介上限為20個字元，請輸入有效的簡介資料")
	}

	if utf8.RuneCountInString(model.Detail) > 200 {
		return errors.New("詳情上限為200個字元，請輸入有效的詳情資料")
	}

	if model.GuestOrder != "" {
		if _, err := strconv.Atoi(model.GuestOrder); err != nil {
			return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
		}
	}

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}

	return g.Table(g.Base.TableName).Where("id", "=", model.ID).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// UpdateGuest 更新活動嘉賓基本設置資料
func (a ActivityModel) UpdateGuest(model EditGuestSettingModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"guest_title", "guest_background"}
	)
	if utf8.RuneCountInString(model.GuestTitle) > 20 {
		return errors.New("錯誤: 手機頁面標題上限為20個字元，請輸入有效的頁面標題")
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
	return a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// MapToModel map轉換[]GuestModel
func (g GuestModel) MapToModel(items []map[string]interface{}) []GuestModel {
	var guests = make([]GuestModel, 0)
	for _, item := range items {
		var guest GuestModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &guest)

		// avatar, _ := item["avatar"].(string)
		if !strings.Contains(guest.Avatar, "system") {
			guest.Avatar = "/admin/uploads/" + guest.UserID + "/" + guest.ActivityID + "/info/guest/" + guest.Avatar
		}

		guests = append(guests, guest)
	}
	return guests
}

// guest.ID, _ = item["id"].(int64)
// guest.ActivityID, _ = item["activity_id"].(string)
// guest.UserID, _ = item["user_id"].(string)
// *****舊*****
// guest.Avatar, _ = item["avatar"].(string)
// *****舊*****

// guest.Name, _ = item["name"].(string)
// guest.Introduce, _ = item["introduce"].(string)
// guest.Detail, _ = item["detail"].(string)
// guest.GuestOrder, _ = item["guest_order"].(int64)

// command.Value{
// 	"activity_id": model.ActivityID,
// 	"avatar":      model.Avatar,
// 	"name":        model.Name,
// 	"introduce":   model.Introduce,
// 	"detail":      model.Detail,
// 	"guest_order": model.GuestOrder,
// }

// fieldValues = command.Value{}
// fields      = []string{"name", "avatar", "guest_order", "introduce", "detail"}
// values      = []string{model.Name, model.Avatar, model.GuestOrder, model.Introduce, model.Detail}

// if len(fieldValues) == 0 {
// 	return nil
// }

// values      = []string{model.GuestTitle, model.GuestBackground}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
