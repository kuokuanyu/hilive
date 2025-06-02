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

// IntroduceModel 資料表欄位
type IntroduceModel struct {
	Base           `json:"-"`
	ID             int64  `json:"id"`
	UserID         string `json:"user_id"`
	ActivityID     string `json:"activity_id"`
	Title          string `json:"title"`
	IntroduceType  string `json:"introduce_type"`
	Content        string `json:"content"`
	IntroduceOrder int64  `json:"introduce_order"`
}

// EditIntroduceModel 資料表欄位
type EditIntroduceModel struct {
	ID             string `json:"id"`
	ActivityID     string `json:"activity_id"`
	Title          string `json:"title"`
	IntroduceType  string `json:"introduce_type"`
	Content        string `json:"content"`
	IntroduceOrder string `json:"introduce_order"`
}

// EditIntroduceSettingModel 資料表欄位
type EditIntroduceSettingModel struct {
	ActivityID     string `json:"activity_id"`
	IntroduceTitle string `json:"introduce_title"`
}

// DefaultIntroduceModel 預設IntroduceModel
func DefaultIntroduceModel() IntroduceModel {
	return IntroduceModel{Base: Base{TableName: config.ACTIVITY_INTRODUCE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (i IntroduceModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) IntroduceModel {
	i.DbConn = dbconn
	i.RedisConn = cacheconn
	i.MongoConn = mongoconn
	return i
}

// SetDbConn 設定connection
// func (i IntroduceModel) SetDbConn(conn db.Connection) IntroduceModel {
// 	i.DbConn = conn
// 	return i
// }

// Find 尋找資料
func (i IntroduceModel) Find(field string, value interface{}) ([]IntroduceModel, error) {
	items, err := i.Table(i.Base.TableName).
		Select("activity_introduce.id", "activity_introduce.activity_id", "title",
			"introduce_type", "content", "introduce_order",

			// 活動
			"activity.user_id").
		Where(field, "=", value).
		LeftJoin(command.Join{
			FieldA:    "activity_introduce.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity",
			Operation: "="}).
		OrderBy("introduce_order", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得活動介紹資訊，請重新查詢")
	}
	return i.MapToModel(items), nil
}

// Add 增加資料
func (i IntroduceModel) Add(model EditIntroduceModel) error {
	var (
		fields = []string{
			"activity_id",
			"title",
			"introduce_type",
			"content",
			"introduce_order",
		}
	)

	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題不能為空並且上限為20個字元，請輸入有效的標題名稱")
	}
	if model.IntroduceType != "text" && model.IntroduceType != "picture" {
		return errors.New("錯誤: 活動介紹類型發生問題，請輸入有效的介紹類型")
	}
	if utf8.RuneCountInString(model.Content) > 200 {
		return errors.New("錯誤: 內容上限為200個字元，請輸入有效的介紹內容")
	}
	if _, err := strconv.Atoi(model.IntroduceOrder); err != nil {
		return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	_, err := i.Table(i.TableName).Insert(FilterFields(data, fields))
	return err
}

// Update 更新資料
func (i IntroduceModel) Update(model EditIntroduceModel) error {
	var (
		fieldValues = command.Value{
			"title":           model.Title,
			"introduce_type":  model.IntroduceType,
			"introduce_order": model.IntroduceOrder}
	)
	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題不能為空並且上限為20個字元，請輸入有效的標題名稱")
	}
	if model.IntroduceType != "" {
		if model.IntroduceType != "text" && model.IntroduceType != "picture" {
			return errors.New("錯誤: 活動介紹類型發生問題，請輸入有效的介紹類型")
		}
	}
	if utf8.RuneCountInString(model.Content) > 200 {
		return errors.New("錯誤: 內容上限為200個字元，請輸入有效的介紹內容")
	}
	if model.IntroduceOrder != "" {
		if _, err := strconv.Atoi(model.IntroduceOrder); err != nil {
			return errors.New("錯誤: 排序資料發生問題，請輸入有效的排序資料")
		}
	}

	if model.IntroduceType == "text" || (model.IntroduceType == "picture" && model.Content != "") {
		fieldValues["content"] = model.Content
	}

	return i.Table(i.Base.TableName).Where("activity_id", "=", model.ActivityID).
		Where("id", "=", model.ID).Update(fieldValues)
}

// UpdateIntroduce 更新活動介紹基本設置資料
func (a ActivityModel) UpdateIntroduce(model EditIntroduceSettingModel) error {
	if model.IntroduceTitle != "" {
		if utf8.RuneCountInString(model.IntroduceTitle) > 20 {
			return errors.New("錯誤: 手機頁面標題上限為20個字元，請輸入有效的頁面標題")
		}
		fieldValues := command.Value{
			"introduce_title": model.IntroduceTitle,
		}
		return a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
	}
	return nil
}

// MapToModel map轉換[]IntroduceModel
func (i IntroduceModel) MapToModel(items []map[string]interface{}) []IntroduceModel {
	var introduces = make([]IntroduceModel, 0)
	for _, item := range items {
		var introduce IntroduceModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &introduce)

		// content, _ := item["content"].(string)
		if introduce.IntroduceType == "picture" {
			introduce.Content = "/admin/uploads/" + introduce.UserID + "/" + introduce.ActivityID + "/info/introduce/" + introduce.Content
		}

		introduces = append(introduces, introduce)
	}
	return introduces
}

// introduce.ID, _ = item["id"].(int64)
// introduce.UserID, _ = item["user_id"].(string)
// introduce.ActivityID, _ = item["activity_id"].(string)
// introduce.Title, _ = item["title"].(string)
// introduce.IntroduceType, _ = item["introduce_type"].(string)
// *****舊*****
// introduce.Content, _ = item["content"].(string)
// *****舊*****

// introduce.IntroduceOrder, _ = item["introduce_order"].(int64

// command.Value{
// 	"activity_id":     model.ActivityID,
// 	"title":           model.Title,
// 	"introduce_type":  model.IntroduceType,
// 	"content":         model.Content,
// 	"introduce_order": model.IntroduceOrder,
// }

// fieldValues = command.Value{}
// fields      = []string{"title", "introduce_type", "introduce_order", "content"}
// values      = []string{model.Title, model.IntroduceType, model.IntroduceOrder, model.Content}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
// if len(fieldValues) == 0 {
// 	return nil
// }
