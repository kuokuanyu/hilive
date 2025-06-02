package models

import (
	"encoding/json"
	"errors"
	"unicode/utf8"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// ScheduleModel 資料表欄位
type ScheduleModel struct {
	Base         `json:"-"`
	ID           int64  `json:"id"`
	ActivityID   string `json:"activity_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ScheduleDate string `json:"schedule_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

// EditScheduleModel 資料表欄位
type EditScheduleModel struct {
	ID           string `json:"id"`
	ActivityID   string `json:"activity_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ScheduleDate string `json:"schedule_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

// EditScheduleSettingModel 資料表欄位
type EditScheduleSettingModel struct {
	ActivityID            string `json:"activity_id"`
	ScheduleTitle         string `json:"schedule_title"`
	ScheduleDisplayDate   string `json:"schedule_display_date"`
	ScheduleDisplayDetail string `json:"schedule_display_detail"`
}

// DefaultScheduleModel ScheduleModel
func DefaultScheduleModel() ScheduleModel {
	return ScheduleModel{Base: Base{TableName: config.ACTIVITY_SCHEDULE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (s ScheduleModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ScheduleModel {
	s.DbConn = dbconn
	s.RedisConn = cacheconn
	s.MongoConn = mongoconn
	return s
}

// SetDbConn 設定connection
// func (s ScheduleModel) SetDbConn(conn db.Connection) ScheduleModel {
// 	s.DbConn = conn
// 	return s
// }

// Find 尋找資料
func (s ScheduleModel) Find(field string, value interface{}) ([]ScheduleModel, error) {
	items, err := s.Table(s.Base.TableName).Where(field, "=", value).
		OrderBy("schedule_date", "asc", "start_time", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得活動行程資訊，請重新查詢")
	}
	return s.MapToModel(items), nil
}

// Add 增加資料
func (s ScheduleModel) Add(model EditScheduleModel) error {
	var (
		fields = []string{
			"activity_id",
			"title",
			"content",
			"schedule_date",
			"start_time",
			"end_time",
		}
	)

	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 主題上限為20個字元，請輸入有效的主題名稱")
	}
	if utf8.RuneCountInString(model.Content) > 200 {
		return errors.New("錯誤: 內容上限為200個字元，請輸入有效的內容")
	}
	if !CompareTime(model.StartTime, model.EndTime) {
		return errors.New("錯誤: 時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	_, err := s.Table(s.TableName).Insert(FilterFields(data, fields))
	return err
}

// Update 更新資料
func (s ScheduleModel) Update(model EditScheduleModel) error {
	var (
		fieldValues = command.Value{
			"title":         model.Title,
			"content":       model.Content,
			"schedule_date": model.ScheduleDate}
	)

	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 主題上限為20個字元，請輸入有效的主題名稱")
	}
	if utf8.RuneCountInString(model.Content) > 200 {
		return errors.New("錯誤: 內容上限為200個字元，請輸入有效的內容")
	}
	if model.StartTime != "" && model.EndTime != "" {
		if !CompareTime(model.StartTime, model.EndTime) {
			return errors.New("錯誤: 時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
		}
		fieldValues["start_time"] = model.StartTime
		fieldValues["end_time"] = model.EndTime
	}

	return s.Table(s.Base.TableName).Where("id", "=", model.ID).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// UpdateSchdule 更新活動行程基本設置資料
func (a ActivityModel) UpdateSchdule(model EditScheduleSettingModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"schedule_title", "schedule_display_date", "schedule_display_detail"}
	)

	if utf8.RuneCountInString(model.ScheduleTitle) > 20 {
		return errors.New("錯誤: 手機頁面標題上限為20個字元，請輸入有效的頁面標題")
	}

	if model.ScheduleDisplayDate != "" {
		if model.ScheduleDisplayDate != "open" && model.ScheduleDisplayDate != "close" {
			return errors.New("錯誤: 顯示日期資料發生問題，請輸入有效的資料")
		}
	}

	if model.ScheduleDisplayDetail != "" {
		if model.ScheduleDisplayDetail != "open" && model.ScheduleDisplayDetail != "close" {
			return errors.New("錯誤: 顯示詳情資料發生問題，請輸入有效的資料")
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
	return a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// MapToModel map轉換[]ScheduleModel
func (s ScheduleModel) MapToModel(items []map[string]interface{}) []ScheduleModel {
	var schedules = make([]ScheduleModel, 0)
	for _, item := range items {
		var schedule ScheduleModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &schedule)

		schedules = append(schedules, schedule)
	}
	return schedules
}

// schedule.ID, _ = item["id"].(int64)
// schedule.ActivityID, _ = item["activity_id"].(string)
// schedule.Title, _ = item["title"].(string)
// schedule.Content, _ = item["content"].(string)
// schedule.ScheduleDate, _ = item["schedule_date"].(string)
// schedule.StartTime, _ = item["start_time"].(string)
// schedule.EndTime, _ = item["end_time"].(string)

// command.Value{
// 	"activity_id":   model.ActivityID,
// 	"title":         model.Title,
// 	"content":       model.Content,
// 	"schedule_date": model.ScheduleDate,
// 	"start_time":    model.StartTime,
// 	"end_time":      model.EndTime,
// }

// fieldValues = command.Value{}
// fields      = []string{"title", "content", "schedule_date"}
// values      = []string{model.Title, model.Content, model.ScheduleDate}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
// if len(fieldValues) == 0 {
// 	return nil
// }

// values      = []string{model.ScheduleTitle, model.ScheduleDisplayDate, model.ScheduleDisplayDetail}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
