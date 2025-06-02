package models

import (
	"encoding/json"
	"errors"
	"strings"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
)

// PhoneModel 資料表欄位
type PhoneModel struct {
	Base     `json:"-"`
	ID       int64  `json:"id"`        // ID
	Phone    string `json:"phone"`     // 電話
	Status   string `json:"status"`    // 驗證狀態
	Times    int64  `json:"times"`     // 驗證次數
	SendTime string `json:"send_time"` // 驗證日期

}

// EditPhoneModel 資料表欄位
type EditPhoneModel struct {
	Phone    string `json:"phone"`     // 用戶電話
	Status   string `json:"status"`    // 驗證狀態
	Times    string `json:"times"`     // 驗證次數
	SendTime string `json:"send_time"` // 驗證日期
}

// DefaultPhoneModel 預設PhoneModel
func DefaultPhoneModel() PhoneModel {
	return PhoneModel{Base: Base{TableName: config.USER_PHONE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (phone PhoneModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) PhoneModel {
	phone.DbConn = dbconn
	phone.RedisConn = cacheconn
	phone.MongoConn = mongoconn
	return phone
}

// SetDbConn 設定connection
// func (phone PhoneModel) SetDbConn(conn db.Connection) PhoneModel {
// 	phone.DbConn = conn
// 	return phone
// }

// Find 尋找電話資料
func (phone PhoneModel) Find(number string) (PhoneModel, error) {
	item, err := phone.Table(phone.Base.TableName).
		Where("phone", "=", number).First()
	if err != nil {
		return phone, errors.New("錯誤: 無法取得電話資訊，請重新查詢")
	}

	return phone.MapToModel(item), nil
}

// Add 增加電話資料
func (phone PhoneModel) Add(model EditPhoneModel) error {
	// 電話
	if len(model.Phone) > 2 {
		if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	} else {
		return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
	}

	_, err := phone.Table(phone.TableName).Insert(command.Value{
		"phone":  model.Phone,
		"status": model.Status,
		"times":  model.Times,
		// "send_time":      time.Now().Format("2006-01-02"),
	})
	if err != nil {
		return errors.New("錯誤: 新增電話資料發生問題，請重新操作")
	}

	return nil
}

// Update 更新資料
func (phone PhoneModel) Update(model EditPhoneModel) error {
	var (
		fields      = []string{"status", "times", "send_time"}
		values      = []string{model.Status, model.Times, model.SendTime}
		fieldValues = command.Value{}
	)

	// 判斷需要更新欄位
	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	err := phone.Table(phone.TableName).
		Where("phone", "=", model.Phone).Update(fieldValues)
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新電話資料發生問題，請重新操作")
	}
	return nil
}

// MapToModel map轉換model
func (phone PhoneModel) MapToModel(m map[string]interface{}) PhoneModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &phone)

	return phone
}

// phone.ID, _ = m["id"].(int64)
// phone.Phone, _ = m["phone"].(string)
// phone.Status, _ = m["status"].(string)
// phone.Times, _ = m["times"].(int64)
// phone.SendTime, _ = m["send_time"].(string)
