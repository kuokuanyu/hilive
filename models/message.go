package models

import (
	"errors"
	"strconv"

	"hilive/modules/db/command"
	"hilive/modules/utils"
)

// EditMessageCheckModel 資料表欄位
type EditMessageCheckModel struct {
	ActivityID              string `json:"activity_id"`
	MessageCheckManualCheck string `json:"message_check_manual_check"` // 手動審核
	MessageCheckSensitivity string `json:"message_check_sensitivity"`  // 敏感詞
}

// EditMessageModel 資料表欄位
type EditMessageModel struct {
	ActivityID                   string `json:"activity_id"`
	MessagePicture               string `json:"message_picture"`
	MessageAuto                  string `json:"message_auto"`
	MessageBan                   string `json:"message_ban"`
	MessageBanSecond             string `json:"message_ban_second"`
	MessageRefreshSecond         string `json:"message_refresh_second"`
	MessageOpen                  string `json:"message_open"`
	Message                      string `json:"message"`
	MessageBackground            string `json:"message_background"`
	MessageTextColor             string `json:"message_text_color"`
	MessageScreenTitleColor      string `json:"message_screen_title_color"`
	MessageTextFrameColor        string `json:"message_text_frame_color"`
	MessageFrameColor            string `json:"message_frame_color"`
	MessageScreenTitleFrameColor string `json:"message_screen_title_frame_color"`
}

// UpdateMessage 更新訊息牆基本設置資料
func (a ActivityModel) UpdateMessage(model EditMessageModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"message_picture", "message_auto", "message_ban",
			"message_ban_second", "message_refresh_second", "message_open",
			"message", "message_background", "message_text_color",
			"message_screen_title_color", "message_text_frame_color",
			"message_frame_color", "message_screen_title_frame_color"}
	)

	if model.MessagePicture != "" {
		if model.MessagePicture != "open" && model.MessagePicture != "close" {
			return errors.New("錯誤: 訊息區可以上傳圖片資料發生問題，請輸入有效的資料")
		}
	}
	if model.MessageAuto != "" {
		if model.MessageAuto != "open" && model.MessageAuto != "close" {
			return errors.New("錯誤: 圖片自動放大顯示資料發生問題，請輸入有效的資料")
		}
	}
	if model.MessageBan != "" {
		if model.MessageBan != "open" && model.MessageBan != "close" {
			return errors.New("錯誤: 禁止連續傳送訊息資料發生問題，請輸入有效的資料")
		}
	}
	if model.MessageBanSecond != "" {
		if _, err := strconv.Atoi(model.MessageBanSecond); err != nil {
			return errors.New("錯誤: 禁止秒數資料發生問題，請輸入有效的秒數")
		}
	}
	if model.MessageRefreshSecond != "" {
		if _, err := strconv.Atoi(model.MessageRefreshSecond); err != nil {
			return errors.New("錯誤: 訊息換頁秒數資料發生問題，請輸入有效的秒數")
		}
	}
	if model.MessageOpen != "" {
		if model.MessageOpen != "open" && model.MessageOpen != "close" {
			return errors.New("錯誤: 跑馬燈開啟資料發生問題，請輸入有效的資料")
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

// UpdateMessageCheck 更新訊息審核資料
func (a ActivityModel) UpdateMessageCheck(model EditMessageCheckModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"message_check_manual_check", "message_check_sensitivity"}
	)

	if model.MessageCheckManualCheck != "" {
		if model.MessageCheckManualCheck != "open" && model.MessageCheckManualCheck != "close" {
			return errors.New("錯誤: 手動審核資料發生問題，請輸入有效的資料")
		}
	}
	if model.MessageCheckSensitivity != "" {
		if model.MessageCheckSensitivity != "replace" && model.MessageCheckSensitivity != "refuse" {
			return errors.New("錯誤: 敏感詞資料發生問題，請輸入有效的資料")
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

// fields      = []string{"message_picture", "message_auto", "message_ban",
// 	"message_ban_second", "message_refresh_second", "message_open",
// 	"message", "message_background"}
// values = []string{model.MessagePicture, model.MessageAuto, model.MessageBan,
// 	model.MessageBanSecond, model.MessageRefreshSecond, model.MessageOpen,
// 	model.Message, model.MessageBackground}

// values = []string{model.MessagePicture, model.MessageAuto, model.MessageBan,
// 	model.MessageBanSecond, model.MessageRefreshSecond, model.MessageOpen,
// 	model.Message, model.MessageBackground, model.MessageTextColor,
// 	model.MessageScreenTitleColor, model.MessageTextFrameColor,
// 	model.MessageFrameColor, model.MessageScreenTitleFrameColor}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

// values      = []string{model.MessageCheckManualCheck, model.MessageCheckSensitivity}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
