package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CustomizeTemplateModel 自定義模板資料
type CustomizeTemplateModel struct {
	Base                  `json:"-" bson:"-"`
	ID                    int64              `json:"id" bson:"id"`                                           // id
	TemplateID            string             `json:"template_id" bson:"template_id"`                         // 模板id
	Game                  string             `json:"game" bson:"game"`                                       // 遊戲類型
	TemplateName          string             `json:"template_name" bson:"template_name"`                     // 模板名稱
	CustomizeTemplateData map[string]interface{} `json:"customize_template_data" bson:"customize_template_data"` // 這裡會包含所有畫面設定資料
}

// DefaultCustomizeTemplateModel 預設CustomizeTemplateModel
func DefaultCustomizeTemplateModel() CustomizeTemplateModel {
	return CustomizeTemplateModel{Base: Base{TableName: config.CUSTOMIZE_TEMPLATE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a CustomizeTemplateModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) CustomizeTemplateModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// FindAll 查詢該遊戲所有自定義模板資訊
func (a CustomizeTemplateModel) FindAll(game string) ([]CustomizeTemplateModel, error) {
	// log.Println("查詢該遊戲所有自定義模板資料(mongo，多筆)")

	items, err := a.MongoConn.FindMany(config.CUSTOMIZE_TEMPLATE,
		bson.M{"game": game},
		options.Find().SetProjection(bson.M{
			"id":            1,
			"template_id":   1,
			"game":          1,
			"template_name": 1,
		})) // 先不取得customize_template_data欄位
	if err != nil {
		return []CustomizeTemplateModel{}, errors.New("錯誤: 無法取得自定義模板資訊(mongo)，請重新查詢")
	}

	// log.Println("items: (檢查customize_template_data是否有資料)", items)

	return MapToCustomizeTemplateModel(items), nil
}

// Find 查詢自定義模板資訊
func (a CustomizeTemplateModel) Find(templateID, templateName string) (CustomizeTemplateModel, error) {
	// log.Println("查詢該用戶自定義模板資料(mongo)，單筆")

	item, err := a.MongoConn.FindOne(config.CUSTOMIZE_TEMPLATE,
		bson.M{"template_id": templateID, "template_name": templateName})
	if err != nil {
		return CustomizeTemplateModel{}, errors.New("錯誤: 無法取得自定義模板資訊(mongo)，請重新查詢")
	}

	a = a.MapToModel(item)

	return a, nil
}

// MapToModel map轉換model
func (a CustomizeTemplateModel) MapToModel(m map[string]interface{}) CustomizeTemplateModel {

	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	return a
}

// MapToCustomizeTemplateModel map轉換[]CustomizeTemplateModel
func MapToCustomizeTemplateModel(items []bson.M) []CustomizeTemplateModel {
	var templates = make([]CustomizeTemplateModel, 0)

	// json解碼，轉換成strcut
	b, _ := json.Marshal(items)
	json.Unmarshal(b, &templates)

	// for _, item := range items {
	// 	var (
	// 		scene CustomizeTemplateModel
	// 	)

	// 	// json解碼，轉換成strcut
	// 	b, _ := json.Marshal(item)
	// 	json.Unmarshal(b, &scene)

	// 	scenes = append(scenes, scene)
	// }
	return templates
}

// Add 新增自定義模板資料(mongo)
func (a CustomizeTemplateModel) Add(game string, templateName string,
	data map[string]interface{}) error {
	// log.Println("新增自定義模板")

	// 取得mongo中的id資料(遞增處理)
	mongoID, _ := a.MongoConn.GetNextSequence(config.CUSTOMIZE_TEMPLATE)

	_, err := a.MongoConn.InsertOne(config.CUSTOMIZE_TEMPLATE, bson.M{
		"id":                      mongoID,
		"template_id":             utils.UUID(20),
		"game":                    game,
		"template_name":           templateName,
		"customize_template_data": data,
	})

	if err != nil {
		return errors.New("錯誤: 無法新增自定義模板資料(mongo)，請重新操作")
	}

	return nil
}

// Update 更新自定義模板資料(mongo)
func (a CustomizeTemplateModel) Update(templateID string, templateName string,
	data map[string]interface{}) error {
	// log.Println("更新自定義模板")

	// 更新資料庫
	_, err := a.MongoConn.UpdateOne(config.CUSTOMIZE_TEMPLATE,
		bson.M{"template_id": templateID}, // 過濾參數
		bson.M{
			"$set": bson.M{
				"template_name":           templateName,
				"customize_template_data": data,
			}, // 更新參數
		})
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 無法更新自定義模板資料(mongo)，請重新操作")
	}

	return nil
}
