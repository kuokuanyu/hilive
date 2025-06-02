package models

import (
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ActivityRequireModel 資料表欄位
type ActivityRequireModel struct {
	// 活動需求
	Base         `json:"-"`
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`       // LINE用戶ID
	Name         string `json:"name"`          // 聯絡人
	Phone        string `json:"phone"`         // 連絡電話
	Email        string `json:"email"`         // 聯絡信箱
	CompanyName  string `json:"company_name"`  // 公司名稱
	ActivityType string `json:"activity_type"` // 活動類型
	People       int64  `json:"people"`        // 活動人數
	StartTime    string `json:"start_time"`    // 活動開始時間
	EndTime      string `json:"end_time"`      // 活動結束時間
	Needs        string `json:"needs"`         // 需求
	Other        string `json:"other"`         // 其他
}

// EditActivityRequireModel 資料表欄位
type EditActivityRequireModel struct {
	UserID       string `json:"user_id"`       // LINE用戶ID
	Name         string `json:"name"`          // 聯絡人
	Phone        string `json:"phone"`         // 連絡電話
	Email        string `json:"email"`         // 聯絡信箱
	CompanyName  string `json:"company_name"`  // 公司名稱
	ActivityType string `json:"activity_type"` // 活動類型
	People       string `json:"people"`        // 活動人數
	StartTime    string `json:"start_time"`    // 活動開始時間
	EndTime      string `json:"end_time"`      // 活動結束時間
	Needs        string `json:"needs"`         // 需求
	Other        string `json:"other"`         // 其他
}

// DefaultActivityRequireModel 預設ActivityRequireModel
func DefaultActivityRequireModel() ActivityRequireModel {
	return ActivityRequireModel{Base: Base{TableName: config.ACTIVITY_REQUIRE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a ActivityRequireModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ActivityRequireModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a ActivityRequireModel) SetDbConn(conn db.Connection) ActivityRequireModel {
// 	a.DbConn = conn
// 	return a
// }

// FindAll 尋找舉辦活動需求資料
// func (a ActivityRequireModel) FindAll() ([]ActivityRequireModel, error) {
// 		items, err := a.Table(a.Base.TableName).
// 			LeftJoin(command.Join{
// 				FieldA:    "line_users.user_id",
// 				FieldB:    "activity_require.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 				All()
// 		if err != nil {
// 			return nil, errors.New("錯誤: 無法取得舉辦活動需求資訊，請重新查詢")
// 		}
// 		return MapToActivityRequireModel(items), nil
// }

// Add 增加資料
func (a ActivityRequireModel) Add(model EditActivityRequireModel) error {
	var (
		fields = []string{
			"user_id",
			"name",
			"phone",
			"email",
			"company_name",
			"activity_type",
			"people",
			"start_time",
			"end_time",
			"needs",
			"other",
		}
	)

	if utf8.RuneCountInString(model.Name) > 20 {
		return errors.New("錯誤: 聯絡人上限為20個字元，請輸入有效的聯絡人")
	}

	if len(model.Phone) > 2 {
		if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	} else {
		return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
	}

	if !strings.Contains(model.Email, "@") {
		return errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	if utf8.RuneCountInString(model.CompanyName) > 20 {
		return errors.New("錯誤: 公司名稱上限為20個字元，請輸入有效的公司名稱")
	}

	// 判斷活動人數上限
	_, err := strconv.Atoi(model.People)
	if err != nil {
		return errors.New("錯誤: 活動人數資料發生問題，請輸入有效的活動人數")
	}

	if !CompareTDatetime(model.StartTime, model.EndTime) {
		return errors.New("錯誤: 活動時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	if _, err := a.Table(a.TableName).Insert(FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 無法新增舉辦活動需求資料，請重新操作")
	}

	return nil
}

// command.Value{
// 	"user_id":       model.UserID,
// 	"name":          model.Name,
// 	"phone":         model.Phone,
// 	"email":         model.Email,
// 	"company_name":  model.CompanyName,
// 	"activity_type": model.ActivityType,
// 	"people":        peopleInt,
// 	"start_time":    model.StartTime,
// 	"end_time":      model.EndTime,
// 	"needs":         model.Needs,
// 	"other":         model.Other,
// }

// MapToActivityRequireModel map轉換[]ActivityRequireModel
// func MapToActivityRequireModel(items []map[string]interface{}) []ActivityRequireModel {
// 	var (
// 		activitys = make([]ActivityRequireModel, 0)
// 	)

// 	// fmt.Println("活動權限數量: ", len(items), items)

// 	for _, item := range items {
// 		var (
// 			activity ActivityRequireModel
// 		)

// 		// 活動
// 		activity.ID, _ = item["activity"].(int64)
// 		activity.UserID, _ = item["user_id"].(string)
// 		activity.Name, _ = item["name"].(string)
// 		activity.Phone, _ = item["phone"].(string)
// 		activity.Email, _ = items[0]["email"].(string)
// 		activity.CompanyName, _ = item["company_name"].(string)
// 		activity.ActivityType, _ = item["activity_type"].(string)
// 		activity.People, _ = item["people"].(int64)
// 		activity.StartTime, _ = item["start_time"].(string)
// 		activity.EndTime, _ = item["end_time"].(string)
// 		activity.Needs, _ = item["needs"].(string)
// 		activity.Other, _ = item["other"].(string)
// 		activitys = append(activitys, activity)
// 	}

// 	return activitys
// }

// MapToModel map轉換model
// func (a ActivityRequireModel) MapToModel(m map[string]interface{}) ActivityRequireModel {
// }
