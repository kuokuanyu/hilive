package models

import (
	"errors"
	"strconv"

	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
)

// EditHoldScreenModel 資料表欄位
type EditHoldScreenModel struct {
	ActivityID              string `json:"activity_id"`
	HoldScreenPrice         string `json:"holdscreen_price"`
	HoldScreenMessageCheck  string `json:"holdscreen_message_check"`
	HoldScreenOnlyPicture   string `json:"holdscreen_only_picture"`
	HoldScreenDisplaySecond string `json:"holdscreen_display_second"`
	HoldScreenBirthdayTopic string `json:"holdscreen_birthday_topic"`
	HoldscreenBanSecond     string `json:"holdscreen_ban_second"`
}

// UpdateHoldScreen 更新霸屏基本設置資料
func (a ActivityModel) UpdateHoldScreen(model EditHoldScreenModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"holdscreen_price",
			"holdscreen_message_check",
			"holdscreen_only_picture",
			"holdscreen_display_second",
			"holdscreen_birthday_topic"}

		fieldValues2 = command.Value{}
		fields2      = []string{"holdscreen_ban_second"}

		err error
	)

	if model.HoldScreenPrice != "" {
		if _, err := strconv.Atoi(model.HoldScreenPrice); err != nil {
			return errors.New("錯誤: 彈幕展示價格資料發生問題，請輸入有效的彈幕展示價格")
		}
	}
	if model.HoldScreenMessageCheck != "" {
		if model.HoldScreenMessageCheck != "open" && model.HoldScreenMessageCheck != "close" {
			return errors.New("錯誤: 霸佔彈幕訊息審核資料發生問題，請輸入有效的訊息審核資料")
		}
	}
	if model.HoldScreenOnlyPicture != "" {
		if model.HoldScreenOnlyPicture != "open" && model.HoldScreenOnlyPicture != "close" {
			return errors.New("錯誤: 只允許以此類型彈幕發送圖片資料發生問題，請輸入有效的資料")
		}
	}
	if model.HoldScreenDisplaySecond != "" {
		if _, err := strconv.Atoi(model.HoldScreenDisplaySecond); err != nil {
			return errors.New("錯誤: 彈幕展示時間資料發生問題，請輸入有效的展示時間資料")
		}
	}
	if model.HoldScreenBirthdayTopic != "" {
		if model.HoldScreenBirthdayTopic != "open" && model.HoldScreenBirthdayTopic != "close" {
			return errors.New("錯誤: 生日主題資料發生問題，請輸入有效的主題資料")
		}
	}
	if model.HoldscreenBanSecond != "" {
		if _, err := strconv.Atoi(model.HoldscreenBanSecond); err != nil {
			return errors.New("錯誤: 訊息傳送間隔秒數資料發生問題，請輸入有效的秒數")
		}
	}

	// 更新activity
	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) == 0 {
	} else {
		err = a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
	}

	// 更新activity_2
	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			fieldValues2[key] = val
		}
	}

	if len(fieldValues2) == 0 {
	} else {
		err = a.Table(config.ACTIVITY_2_TABLE).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues2)
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
	}

	return nil
}

// values = []string{model.HoldScreenPrice, model.HoldScreenMessageCheck, model.HoldScreenOnlyPicture,
// 	model.HoldScreenDisplaySecond, model.HoldScreenBirthdayTopic}

// values2      = []string{model.HoldscreenBanSecond}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

// for i, value2 := range values2 {
// 	if value2 != "" {
// 		fieldValues2[fields2[i]] = value2
// 	}
// }
