package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// OverviewModel 資料表欄位
type OverviewModel struct {
	Base         `json:"-"`
	ID           int64  `json:"id"`
	OverviewType string `json:"overview_type"`
	OverviewName string `json:"overview_name"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ClassID      string `json:"class_id"`
	DivID        string `json:"div_id"`
	URL          string `json:"url"`
}

// EditOverviewModel 資料表欄位
type EditOverviewModel struct {
	ID           string `json:"id"`
	OverviewType string `json:"overview_type"`
	OverviewName string `json:"overview_name"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ClassID      string `json:"class_id"`
	DivID        string `json:"div_id"`
	URL          string `json:"url"`
}

// EditActivityOverviewModel 資料表欄位(編輯活動總覽開關)
type EditActivityOverviewModel struct {
	ActivityID             string `json:"activity_id"`
	OverviewMessage        string `json:"overview_message"`
	OverviewTopic          string `json:"overview_topic"`
	OverviewQuestion       string `json:"overview_question"`
	OverviewDanmu          string `json:"overview_danmu"`
	OverviewSpecialDanmu   string `json:"overview_special_danmu"`
	OverviewPicture        string `json:"overview_picture"`
	OverviewHoldscreen     string `json:"overview_holdscreen"`
	OverviewGeneral        string `json:"overview_general"`
	OverviewThreed         string `json:"overview_threed"`
	OverviewCountdown      string `json:"overview_countdown"`
	OverviewLottery        string `json:"overview_lottery"`
	OverviewRedpack        string `json:"overview_redpack"`
	OverviewRopepack       string `json:"overview_ropepack"`
	OverviewWhackMole      string `json:"overview_whack_mole"`
	OverviewDrawNumbers    string `json:"overview_draw_numbers"`
	OverviewMonopoly       string `json:"overview_monopoly"`
	OverviewQA             string `json:"overview_qa"`
	OverviewTugofwar       string `json:"overview_tugofwar"`
	OverviewBingo          string `json:"overview_bingo"`
	OverviewSignname       string `json:"overview_signname"`
	Overview3DGachaMachine string `json:"overview_3d_gacha_machine"`
	OverviewVote           string `json:"overview_vote"`
	OverviewChatroom       string `json:"overview_chatroom"`
	OverviewInteract       string `json:"overview_interact"`
	OverviewInfo           string `json:"overview_info"`
}

// DefaultOverviewModel 預設OverviewModel
func DefaultOverviewModel() OverviewModel {
	return OverviewModel{Base: Base{TableName: config.ACTIVITY_OVERVIEW_GAME_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a OverviewModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) OverviewModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a OverviewModel) SetDbConn(conn db.Connection) OverviewModel {
// 	a.DbConn = conn
// 	return a
// }

// Find 尋找資料
func (a OverviewModel) Find(id int64) (interface{}, error) {
	var (
		data interface{}
	// 	items = make([]map[string]interface{}, 0)
	// 	err   error
	)

	if id != 0 {
		// 查詢單個overview資料
		item, err := a.Table(a.Base.TableName).
			Where("id", "=", id).First()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得總覽資訊，請重新查詢")
		}
		data = a.MapToModel(item)
	} else if id == 0 {
		// 查詢所有overview資料
		items, err := a.Table(a.Base.TableName).All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得總覽資訊，請重新查詢")
		}
		data = MapToOverviewModel(items)
	}
	// if err != nil {
	// 	return nil, errors.New("錯誤: 無法取得總覽資訊，請重新查詢")
	// }

	return data, nil
}

// Add 新增總覽資料
func (a OverviewModel) Add(model EditOverviewModel) error {
	var (
		fields = []string{
			"overview_type",
			"overview_name",
			"name",
			"description",
			"class_id",
			"div_id",
			"url",
		}
	)

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	if _, err := a.Table(a.Base.TableName).Insert(FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 無法新增總覽資料，請重新操作")
	}
	return nil
}

// command.Value{
// 	"overview_type": model.OverviewType,
// 	"overview_name": model.OverviewName,
// 	"name":          model.Name,
// 	"description":   model.Description,
// 	"class_id":      model.ClassID,
// 	"div_id":        model.DivID,
// 	"url":           model.URL,
// }

// Update 更新總覽資料
func (a OverviewModel) Update(model EditOverviewModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"overview_type", "overview_name", "name",
			"description", "class_id", "div_id", "url"}
		// values = []string{model.OverviewType, model.OverviewName, model.Name,
		// 	model.Description, model.ClassID, model.DivID,
		// 	model.URL}
	)

	id, err := strconv.Atoi(model.ID)
	if err != nil {
		return errors.New("錯誤: ID發生問題，請輸入有效的ID")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 更新總覽資料
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
			return errors.New("錯誤: 無法更新總覽資料，請重新操作")
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
func (a OverviewModel) MapToModel(item map[string]interface{}) OverviewModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &a)

	return a
}

// MapToOverviewModel map轉換[]OverviewModel
func MapToOverviewModel(items []map[string]interface{}) []OverviewModel {
	var overviews = make([]OverviewModel, 0)
	for _, item := range items {
		var (
			overview OverviewModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &overview)

		overviews = append(overviews, overview)
	}
	return overviews
}

// UpdateOverview 更新活動總覽資料
func (a ActivityModel) UpdateOverview(model EditActivityOverviewModel) error {
	var (
		fieldValues  = command.Value{}
		fieldValues2 = command.Value{}
		fields       = []string{"overview_message", "overview_topic", "overview_question",
			"overview_danmu", "overview_special_danmu", "overview_picture", "overview_holdscreen",
			"overview_general", "overview_threed", "overview_countdown", "overview_lottery",
			"overview_redpack", "overview_ropepack", "overview_whack_mole", "overview_draw_numbers",
			"overview_monopoly", "overview_qa", "overview_tugofwar", "overview_bingo", "overview_signname",
			"overview_3d_gacha_machine",
		}
		fields2 = []string{
			"overview_vote",
			"overview_chatroom",
			"overview_interact",
			"overview_info",
		}
	)

	// activity
	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			if val != "open" && val != "close" {
				return errors.New("錯誤: 活動總覽資料發生問題，請輸入有效的資料")
			}

			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		err := a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 更新活動總覽資料發生錯誤(1)")
		}
	}

	// activity_2
	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			if val != "open" && val != "close" {
				return errors.New("錯誤: 活動總覽資料發生問題，請輸入有效的資料")
			}

			fieldValues2[key] = val
		}
	}

	if len(fieldValues2) != 0 {
		err := a.Table(config.ACTIVITY_2_TABLE).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues2)
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 更新活動總覽資料發生錯誤(2)")
		}
	}

	return nil
}

// values = []string{model.OverviewMessage, model.OverviewTopic, model.OverviewQuestion,
// 	model.OverviewDanmu, model.OverviewSpecialDanmu, model.OverviewPicture, model.OverviewHoldscreen,
// 	model.OverviewGeneral, model.OverviewThreed, model.OverviewCountdown, model.OverviewLottery,
// 	model.OverviewRedpack, model.OverviewRopepack, model.OverviewWhackMole, model.OverviewDrawNumbers,
// 	model.OverviewMonopoly, model.OverviewQA, model.OverviewTugofwar, model.OverviewBingo,
// 	model.OverviewSignname, model.Overview3DGachaMachine,
// }
// values2 = []string{
// 	model.OverviewVote,
// 	model.OverviewChatroom,
// 	model.OverviewInteract,
// 	model.OverviewInfo}

// for i, value := range values {
// 	if value != "" {
// 		if value != "open" && value != "close" {
// 			return errors.New("錯誤: 活動總覽資料發生問題，請輸入有效的資料")
// 		}
// 		fieldValues[fields[i]] = value
// 	}
// }
// if len(fieldValues) == 0 {
// 	return nil
// }

// for i, value := range values2 {
// 	if value != "" {
// 		if value != "open" && value != "close" {
// 			return errors.New("錯誤: 活動總覽資料發生問題，請輸入有效的資料")
// 		}
// 		fieldValues2[fields2[i]] = value
// 	}
// }
// if len(fieldValues2) == 0 {
// 	return nil
// }

// overview.ID, _ = item["id"].(int64)
// overview.OverviewType, _ = item["overview_type"].(string)
// overview.OverviewName, _ = item["overview_name"].(string)
// overview.Name, _ = item["name"].(string)
// overview.Description, _ = item["description"].(string)
// overview.ClassID, _ = item["class_id"].(string)
// overview.DivID, _ = item["div_id"].(string)
// overview.URL, _ = item["url"].(string)

// a.ID, _ = item["id"].(int64)
// a.OverviewType, _ = item["overview_type"].(string)
// a.OverviewName, _ = item["overview_name"].(string)
// a.Name, _ = item["name"].(string)
// a.Description, _ = item["description"].(string)
// a.ClassID, _ = item["class_id"].(string)
// a.DivID, _ = item["div_id"].(string)
// a.URL, _ = item["url"].(string)

// NewOverviewModel 資料表欄位
// type NewOverviewModel struct {
// 	OverviewType string  `json:"overview_type"`
// 	OverviewName string  `json:"overview_name"`
// 	Name string  `json:"name"`
// 	Description string  `json:"description"`
// 	ClassID string  `json:"class_id"`
// 	DivID string  `json:"div_id"`
// 	URL string  `json:"url"`
// }

// NewOverviewModel 資料表欄位
// type NewOverviewModel struct {
// 	ActivityID string `json:"activity_id"`
// 	OverviewID string `json:"overview_id"`
// 	Open       string `json:"open"`
// }

// DefaultOverviewModel 預設OverviewModel
// func DefaultOverviewModel() OverviewModel {
// 	return OverviewModel{Base: Base{TableName: config.ACTIVITY_OVERVIEW_TABLE}}
// }

// SetDbConn 設定connection
// func (o OverviewModel) SetDbConn(conn db.Connection) OverviewModel {
// 	o.Conn = conn
// 	return o
// }

// Find 尋找資料
// func (o OverviewModel) Find(field string, value interface{}) ([]OverviewModel, error) {
// 	items, err := o.Table(o.Base.TableName).Where(field, "=", value).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得活動總覽資訊，請重新查詢")
// 	}
// 	return o.MapToModel(items), nil
// }

// Add 增加資料
// func (o OverviewModel) Add(model NewOverviewModel) error {
// 	if model.ActivityID == "" || model.OverviewID == "" {
// 		return errors.New("錯誤: id資料不能為空，請輸入有效的id資料")
// 	}
// 	if model.Open != "open" && model.Open != "close" {
// 		return errors.New("錯誤: 是否開啟資料發生問題，請輸入有效的資料")
// 	}
// 	_, err := o.Table(o.TableName).Insert(command.Value{
// 		"activity_id": model.ActivityID,
// 		"overview_id": model.OverviewID,
// 		"open":        model.Open,
// 	})
// 	return err
// }

// Update 更新資料
// func (o OverviewModel) Update(model EditOverviewModel) error {
// 	if model.Open != "" {
// 		if model.Open != "open" && model.Open != "close" {
// 			return errors.New("錯誤: 是否開啟資料發生問題，請輸入有效的資料")
// 		}
// 		update := command.Value{
// 			"open": model.Open,
// 		}
// 		return o.Table(o.Base.TableName).Where("activity_id", "=", model.ActivityID).
// 			Where("overview_id", "=", model.OverviewID).Update(update)
// 	}
// 	return nil
// }

// MapToModel map轉換[]OverviewModel
// func (o OverviewModel) MapToModel(items []map[string]interface{}) []OverviewModel {
// 	var overviews []OverviewModel
// 	for _, item := range items {
// 		var overview OverviewModel
// 		overview.ID, _ = item["id"].(int64)
// 		overview.ActivityID, _ = item["activity_id"].(string)
// 		overview.OverviewID, _ = item["overview_id"].(int64)
// 		overview.Open, _ = item["open"].(string)
// 		overviews = append(overviews, overview)
// 	}
// 	return overviews
// }

//---------------------------下次修改活動總覽時，上半部都刪除
