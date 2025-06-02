package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"unicode/utf8"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// MaterialModel 資料表欄位
type MaterialModel struct {
	Base          `json:"-"`
	ID            int64  `json:"id"`
	ActivityID    string `json:"activity_id"`
	Title         string `json:"title"`
	Introduce     string `json:"introduce"`
	Link          string `json:"link"`
	MaterialOrder int64  `json:"material_order"`
}

// EditMaterialModel 資料表欄位
type EditMaterialModel struct {
	ID            string `json:"id"`
	ActivityID    string `json:"activity_id"`
	Title         string `json:"title"`
	Introduce     string `json:"introduce"`
	Link          string `json:"link"`
	MaterialOrder string `json:"material_order"`
}

// EditMaterialSettingModel 資料表欄位
type EditMaterialSettingModel struct {
	ActivityID    string `json:"activity_id"`
	MaterialTitle string `json:"material_title"`
}

// DefaultMaterialModel 預設MaterialModel
func DefaultMaterialModel() MaterialModel {
	return MaterialModel{Base: Base{TableName: config.ACTIVITY_MATERIAL_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m MaterialModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) MaterialModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m MaterialModel) SetDbConn(conn db.Connection) MaterialModel {
// 	m.DbConn = conn
// 	return m
// }

// Find 尋找資料
func (m MaterialModel) Find(field string, value interface{}) ([]MaterialModel, error) {
	items, err := m.Table(m.Base.TableName).Where(field, "=", value).
		OrderBy("material_order", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得活動資料，請重新查詢")
	}
	return m.MapToModel(items), nil
}

// Add 增加資料
func (m MaterialModel) Add(model EditMaterialModel) error {
	var (
		fields = []string{
			"activity_id",
			"title",
			"introduce",
			"link",
			"material_order",
		}
	)

	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題上限為20個字元，請輸入有效的標題名稱")
	}

	if utf8.RuneCountInString(model.Introduce) > 200 {
		return errors.New("錯誤: 資料說明上限為200個字元，請輸入有效的說明")
	}

	if _, err := strconv.Atoi(model.MaterialOrder); err != nil {
		return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	_, err := m.Table(m.TableName).Insert(FilterFields(data, fields))
	return err
}

// Update 更新資料
func (m MaterialModel) Update(model EditMaterialModel) error {
	var (
		fieldValues = command.Value{
			"introduce":      model.Introduce,
			"link":           model.Link,
			"title":          model.Title,
			"material_order": model.MaterialOrder}
	)
	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題上限為20個字元，請輸入有效的標題名稱")
	}
	if utf8.RuneCountInString(model.Introduce) > 200 {
		return errors.New("錯誤: 資料說明上限為200個字元，請輸入有效的說明")
	}
	if model.MaterialOrder != "" {
		if _, err := strconv.Atoi(model.MaterialOrder); err != nil {
			return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
		}
	}

	return m.Table(m.Base.TableName).Where("id", "=", model.ID).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// UpdateMaterial 更新活動資料基本設置資料
func (a ActivityModel) UpdateMaterial(model EditMaterialSettingModel) error {
	if model.MaterialTitle != "" {
		if utf8.RuneCountInString(model.MaterialTitle) > 20 {
			return errors.New("錯誤: 手機頁面標題上限為20個字元，請輸入有效的頁面標題")
		}

		fieldValues := command.Value{
			"material_title": model.MaterialTitle,
		}

		return a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
	}
	return nil
}

// MapToModel map轉換[]MaterialModel
func (m MaterialModel) MapToModel(items []map[string]interface{}) []MaterialModel {
	var materials = make([]MaterialModel, 0)
	for _, item := range items {
		var material MaterialModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &material)

		materials = append(materials, material)
	}
	return materials
}

// material.ID, _ = item["id"].(int64)
// material.ActivityID, _ = item["activity_id"].(string)
// material.Title, _ = item["title"].(string)
// material.Introduce, _ = item["introduce"].(string)
// material.Link, _ = item["link"].(string)
// material.MaterialOrder, _ = item["material_order"].(int64)

// command.Value{
// 	"activity_id":    model.ActivityID,
// 	"title":          model.Title,
// 	"introduce":      model.Introduce,
// 	"link":           model.Link,
// 	"material_order": model.MaterialOrder,
// }

// fieldValues = command.Value{"introduce": model.Introduce, "link": model.Link}
// fields      = []string{"title", "material_order"}
// values      = []string{model.Title, model.MaterialOrder}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
// if len(fieldValues) == 0 {
// 	return nil
// }
