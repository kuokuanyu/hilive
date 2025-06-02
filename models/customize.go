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
)

// CustomizeModel 資料表欄位
type CustomizeModel struct {
	Base         `json:"-"`
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	ActivityID   string `json:"activity_id"`
	Ext1Name     string `json:"ext_1_name"`
	Ext1Type     string `json:"ext_1_type"`
	Ext1Options  string `json:"ext_1_options"`
	Ext1Required string `json:"ext_1_required"`
	Ext1Unique   string `json:"ext_1_unique"`

	Ext2Name     string `json:"ext_2_name"`
	Ext2Type     string `json:"ext_2_type"`
	Ext2Options  string `json:"ext_2_options"`
	Ext2Required string `json:"ext_2_required"`
	Ext2Unique   string `json:"ext_2_unique"`

	Ext3Name     string `json:"ext_3_name"`
	Ext3Type     string `json:"ext_3_type"`
	Ext3Options  string `json:"ext_3_options"`
	Ext3Required string `json:"ext_3_required"`
	Ext3Unique   string `json:"ext_3_unique"`

	Ext4Name     string `json:"ext_4_name"`
	Ext4Type     string `json:"ext_4_type"`
	Ext4Options  string `json:"ext_4_options"`
	Ext4Required string `json:"ext_4_required"`
	Ext4Unique   string `json:"ext_4_unique"`

	Ext5Name     string `json:"ext_5_name"`
	Ext5Type     string `json:"ext_5_type"`
	Ext5Options  string `json:"ext_5_options"`
	Ext5Required string `json:"ext_5_required"`
	Ext5Unique   string `json:"ext_5_unique"`

	Ext6Name     string `json:"ext_6_name"`
	Ext6Type     string `json:"ext_6_type"`
	Ext6Options  string `json:"ext_6_options"`
	Ext6Required string `json:"ext_6_required"`
	Ext6Unique   string `json:"ext_6_unique"`

	Ext7Name     string `json:"ext_7_name"`
	Ext7Type     string `json:"ext_7_type"`
	Ext7Options  string `json:"ext_7_options"`
	Ext7Required string `json:"ext_7_required"`
	Ext7Unique   string `json:"ext_7_unique"`

	Ext8Name     string `json:"ext_8_name"`
	Ext8Type     string `json:"ext_8_type"`
	Ext8Options  string `json:"ext_8_options"`
	Ext8Required string `json:"ext_8_required"`
	Ext8Unique   string `json:"ext_8_unique"`

	Ext9Name     string `json:"ext_9_name"`
	Ext9Type     string `json:"ext_9_type"`
	Ext9Options  string `json:"ext_9_options"`
	Ext9Required string `json:"ext_9_required"`
	Ext9Unique   string `json:"ext_9_unique"`

	Ext10Name     string `json:"ext_10_name"`
	Ext10Type     string `json:"ext_10_type"`
	Ext10Options  string `json:"ext_10_options"`
	Ext10Required string `json:"ext_10_required"`
	Ext10Unique   string `json:"ext_10_unique"`

	ExtEmailRequired string `json:"ext_email_required"`
	ExtPhoneRequired string `json:"ext_phone_required"`
	InfoPicture      string `json:"info_picture"`

	ActivityName string `json:"activity_name"`
	Device       string `json:"device"`

	MessageAmount    int64  `json:"message_amount"`     // 簡訊數量
	PushPhoneMessage string `json:"push_phone_message"` // 發送手機簡訊

	SendMail   string `json:"send_mail"`   // 發送郵件
	MailAmount int64  `json:"mail_amount"` // 郵件數量

	CustomizePassword string `json:"customize_password"` // 是否自定義設置驗證碼
	SignCheck         string `json:"sign_check"`         // 簽到審核

	HostScan string `json:"host_scan"` // 主持人掃qrcode

	Number int64 `json:"number"` // 號碼

	CustomizeDefaultAvatar string `json:"customize_default_avatar"` // 自定義人員預設頭像
}

// EditCustomizeModel 資料表欄位
type EditCustomizeModel struct {
	ActivityID             string `json:"activity_id"`
	Field                  string `json:"field"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Options                string `json:"options"`
	Required               string `json:"required"`
	Unique                 string `json:"unique"`
	ExtEmailRequired       string `json:"ext_email_required"`
	ExtPhoneRequired       string `json:"ext_phone_required"`
	InfoPicture            string `json:"info_picture"`
	InfoPictureDefaultFlag string `json:"info_picture__default_flag"`
	IsDelete               string `json:"is_delete"`
}

// DefaultCustomizeModel 預設CustomizeModel
func DefaultCustomizeModel() CustomizeModel {
	return CustomizeModel{Base: Base{TableName: config.ACTIVITY_CUSTOMIZE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m CustomizeModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) CustomizeModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m CustomizeModel) SetDbConn(conn db.Connection) CustomizeModel {
// 	m.DbConn = conn
// 	return m
// }

// Find 尋找資料後轉換成CustomizeModel
func (m CustomizeModel) Find(activityID string) (CustomizeModel, error) {
	item, err := m.Table(m.Base.TableName).
		Select("activity_customize.id", "activity_customize.activity_id", "ext_email_required",
			"ext_phone_required", "info_picture",

			// 自定義
			"activity_customize.ext_1_name", "activity_customize.ext_1_type", "activity_customize.ext_1_options", "activity_customize.ext_1_required",
			"activity_customize.ext_2_name", "activity_customize.ext_2_type", "activity_customize.ext_2_options", "activity_customize.ext_2_required",
			"activity_customize.ext_3_name", "activity_customize.ext_3_type", "activity_customize.ext_3_options", "activity_customize.ext_3_required",
			"activity_customize.ext_4_name", "activity_customize.ext_4_type", "activity_customize.ext_4_options", "activity_customize.ext_4_required",
			"activity_customize.ext_5_name", "activity_customize.ext_5_type", "activity_customize.ext_5_options", "activity_customize.ext_5_required",
			"activity_customize.ext_6_name", "activity_customize.ext_6_type", "activity_customize.ext_6_options", "activity_customize.ext_6_required",
			"activity_customize.ext_7_name", "activity_customize.ext_7_type", "activity_customize.ext_7_options", "activity_customize.ext_7_required",
			"activity_customize.ext_8_name", "activity_customize.ext_8_type", "activity_customize.ext_8_options", "activity_customize.ext_8_required",
			"activity_customize.ext_9_name", "activity_customize.ext_9_type", "activity_customize.ext_9_options", "activity_customize.ext_9_required",
			"activity_customize.ext_10_name", "activity_customize.ext_10_type", "activity_customize.ext_10_options", "activity_customize.ext_10_required",

			"activity_customize.ext_1_unique", "activity_customize.ext_2_unique", "activity_customize.ext_3_unique", "activity_customize.ext_4_unique",
			"activity_customize.ext_5_unique", "activity_customize.ext_6_unique", "activity_customize.ext_7_unique", "activity_customize.ext_8_unique",
			"activity_customize.ext_9_unique", "activity_customize.ext_10_unique",

			// 活動
			"activity.user_id", "activity.device", "activity.activity_name",
			"activity.message_amount", "activity.push_phone_message",
			"activity.send_mail", "activity.mail_amount",
			"activity.host_scan",
			"activity.customize_password",
			"activity.sign_check", // 簽到審核
			"activity.number",     // 號碼

			"activity_2.customize_default_avatar", // 自定義人員預設頭像
		).
		LeftJoin(command.Join{
			FieldA:    "activity_customize.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity",
			Operation: "="}).
		LeftJoin(command.Join{
			FieldA:    "activity_customize.activity_id",
			FieldA1:   "activity_2.activity_id",
			Table:     "activity_2",
			Operation: "="}).
		Where("activity_customize.activity_id", "=", activityID).First()
	if err != nil || item == nil {
		return CustomizeModel{}, errors.New("錯誤: 無法取得自定義欄位資訊，請重新查詢")
	}

	return m.MapToModel(item), nil
}

// Add 新增自定義欄位資料
func (m CustomizeModel) Add(model EditCustomizeModel) error {
	if _, err := m.Table(m.Base.TableName).Insert(command.Value{
		"activity_id": model.ActivityID,
	}); err != nil {
		return errors.New("錯誤: 無法新增自定義欄位資料，請重新操作")
	}
	return nil
}

// Update 更新自定義欄位資料
func (m CustomizeModel) Update(model EditCustomizeModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"ext_email_required", "ext_phone_required"}
		// values      = []string{model.ExtEmailRequired, model.ExtPhoneRequired}
	)

	// 判斷是否上傳圖片或清空圖片
	if model.InfoPicture != "" || model.InfoPictureDefaultFlag == "1" {
		fieldValues["info_picture"] = model.InfoPicture
	}

	if model.Field != "" {
		if model.IsDelete != "true" && model.IsDelete != "false" {
			return errors.New("錯誤: 是否清除欄位資料的資料發生問題，請輸入有效的資料")
		}

		fieldValues[model.Field+"_name"] = model.Name
		fieldValues[model.Field+"_type"] = model.Type
		fieldValues[model.Field+"_options"] = model.Options
		fieldValues[model.Field+"_required"] = model.Required
		fieldValues[model.Field+"_unique"] = model.Unique

		if model.IsDelete == "true" {
			// 更新自定義欄位資料時必須將activity_applysign資料表該欄位(ext_n)資料清空
			if err := DefaultApplysignModel().
				SetConn(m.DbConn, m.RedisConn, m.MongoConn).
				DeleteExt(model.ActivityID, model.Field); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 自定義開關欄位
	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			if val != "true" && val != "false" {
				return errors.New("錯誤: 自定義開關資料(是否必填信箱、是否必填電話)發生問題，請輸入有效的資料")
			}
			fieldValues[key] = val
		}
	}

	if len(fieldValues) == 0 {
		return nil
	}

	if err := m.Table(m.Base.TableName).Where("activity_id", "=", model.ActivityID).
		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 無法更新自定義欄位資料，請重新操作")
	}
	return nil
}

// for i, value := range values {
// 	if value != "" {
// 		if value != "true" && value != "false" {
// 			return errors.New("錯誤: 自定義開關資料(是否必填信箱、是否必填電話)發生問題，請輸入有效的資料")
// 		}
// 		fieldValues[fields[i]] = value
// 	}
// }

// MapToModel map轉換model
func (m CustomizeModel) MapToModel(item map[string]interface{}) CustomizeModel {
	// 活動

	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &m)

	// m.InfoPicture, _ = item["info_picture"].(string)
	if m.InfoPicture != "" {
		m.InfoPicture = "/admin/uploads/" + m.UserID + "/" + m.ActivityID + "/applysign/customize/" + m.InfoPicture
	}

	return m
}

// m.ActivityName, _ = item["activity_name"].(string)
// m.Device, _ = item["device"].(string)

// m.MessageAmount, _ = item["message_amount"].(int64)
// m.PushPhoneMessage, _ = item["push_phone_message"].(string)

// m.SendMail, _ = item["send_mail"].(string)
// m.MailAmount, _ = item["mail_amount"].(int64)

// m.HostScan, _ = item["host_scan"].(string)
// m.CustomizePassword, _ = item["customize_password"].(string)
// m.SignCheck, _ = item["sign_check"].(string)
// m.Number, _ = item["number"].(int64)

// m.CustomizeDefaultAvatar, _ = item["customize_default_avatar"].(string) // 自定義人員預設頭像

// m.ID, _ = item["id"].(int64)
// 	m.UserID, _ = item["user_id"].(string)
// 	m.ActivityID, _ = item["activity_id"].(string)

// 	m.Ext1Name, _ = item["ext_1_name"].(string)
// 	m.Ext1Type, _ = item["ext_1_type"].(string)
// 	m.Ext1Options, _ = item["ext_1_options"].(string)
// 	m.Ext1Required, _ = item["ext_1_required"].(string)

// 	m.Ext2Name, _ = item["ext_2_name"].(string)
// 	m.Ext2Type, _ = item["ext_2_type"].(string)
// 	m.Ext2Options, _ = item["ext_2_options"].(string)
// 	m.Ext2Required, _ = item["ext_2_required"].(string)

// 	m.Ext3Name, _ = item["ext_3_name"].(string)
// 	m.Ext3Type, _ = item["ext_3_type"].(string)
// 	m.Ext3Options, _ = item["ext_3_options"].(string)
// 	m.Ext3Required, _ = item["ext_3_required"].(string)

// 	m.Ext4Name, _ = item["ext_4_name"].(string)
// 	m.Ext4Type, _ = item["ext_4_type"].(string)
// 	m.Ext4Options, _ = item["ext_4_options"].(string)
// 	m.Ext4Required, _ = item["ext_4_required"].(string)

// 	m.Ext5Name, _ = item["ext_5_name"].(string)
// 	m.Ext5Type, _ = item["ext_5_type"].(string)
// 	m.Ext5Options, _ = item["ext_5_options"].(string)
// 	m.Ext5Required, _ = item["ext_5_required"].(string)

// 	m.Ext6Name, _ = item["ext_6_name"].(string)
// 	m.Ext6Type, _ = item["ext_6_type"].(string)
// 	m.Ext6Options, _ = item["ext_6_options"].(string)
// 	m.Ext6Required, _ = item["ext_6_required"].(string)

// 	m.Ext7Name, _ = item["ext_7_name"].(string)
// 	m.Ext7Type, _ = item["ext_7_type"].(string)
// 	m.Ext7Options, _ = item["ext_7_options"].(string)
// 	m.Ext7Required, _ = item["ext_7_required"].(string)

// 	m.Ext8Name, _ = item["ext_8_name"].(string)
// 	m.Ext8Type, _ = item["ext_8_type"].(string)
// 	m.Ext8Options, _ = item["ext_8_options"].(string)
// 	m.Ext8Required, _ = item["ext_8_required"].(string)

// 	m.Ext9Name, _ = item["ext_9_name"].(string)
// 	m.Ext9Type, _ = item["ext_9_type"].(string)
// 	m.Ext9Options, _ = item["ext_9_options"].(string)
// 	m.Ext9Required, _ = item["ext_9_required"].(string)

// 	m.Ext10Name, _ = item["ext_10_name"].(string)
// 	m.Ext10Type, _ = item["ext_10_type"].(string)
// 	m.Ext10Options, _ = item["ext_10_options"].(string)
// 	m.Ext10Required, _ = item["ext_10_required"].(string)

// 	m.Ext1Unique, _ = item["ext_1_unique"].(string)
// 	m.Ext2Unique, _ = item["ext_2_unique"].(string)
// 	m.Ext3Unique, _ = item["ext_3_unique"].(string)
// 	m.Ext4Unique, _ = item["ext_4_unique"].(string)
// 	m.Ext5Unique, _ = item["ext_5_unique"].(string)
// 	m.Ext6Unique, _ = item["ext_6_unique"].(string)
// 	m.Ext7Unique, _ = item["ext_7_unique"].(string)
// 	m.Ext8Unique, _ = item["ext_8_unique"].(string)
// 	m.Ext9Unique, _ = item["ext_9_unique"].(string)
// 	m.Ext10Unique, _ = item["ext_10_unique"].(string)

// 	m.ExtEmailRequired, _ = item["ext_email_required"].(string)
// 	m.ExtPhoneRequired, _ = item["ext_phone_required"].(string)

// Find 尋找資料
// func (m CustomizeModel) Find(activityID string) (map[string]interface{}, error) {
// 	item, err := m.Table(m.Base.TableName).
// 		Where("activity_id", "=", activityID).First()
// 	if err != nil || item == nil {
// 		return nil, errors.New("錯誤: 無法取得自定義欄位資訊，請重新查詢")
// 	}
// 	return item, nil
// }

// if model.Field == "" {
// 	return errors.New("錯誤: field資料發生問題，請輸入有效的field資料")
// }

// if model.Name == "" {
// 	return errors.New("錯誤: name資料發生問題，請輸入有效的name資料")
// }
// if model.Type != "text" && model.Type != "radio" &&
// 	model.Type != "select" && model.Type != "checkbox" && model.Type != "textarea" &&
// 	model.Type != "date" && model.Type != "time" {
// 	return errors.New("錯誤: type資料發生問題，請輸入有效的type資料")
// }
// if model.Required != "true" && model.Required != "false" {
// 	return errors.New("錯誤: required資料發生問題，請輸入有效的required資料")
// }

// if err := DefaultGameStaffModel().SetDbConn(m.DbConn).
// 	DeleteExt(model.ActivityID, model.Field); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return err
// }
// if err := DefaultPrizeStaffModel().SetDbConn(m.DbConn).
// 	DeleteExt(model.ActivityID, model.Field); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	re
