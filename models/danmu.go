package models

import (
	"errors"
	"strconv"

	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
)

// EditDanmuModel 資料表欄位
type EditDanmuModel struct {
	ActivityID              string `json:"activity_id"`
	DanmuLoop               string `json:"danmu_loop"`
	DanmuTop                string `json:"danmu_top"`
	DanmuMid                string `json:"danmu_mid"`
	DanmuBottom             string `json:"danmu_bottom"`
	DanmuDisplayName        string `json:"danmu_display_name"`
	DanmuDisplayAvatar      string `json:"danmu_display_avatar"`
	DanmuSize               string `json:"danmu_size"`
	DanmuSpeed              string `json:"danmu_speed"`
	DanmuDensity            string `json:"danmu_density"`
	DanmuOpacity            string `json:"danmu_opacity"`
	DanmuBackground         string `json:"danmu_background"`
	DanmuNewBackgroundColor string `json:"danmu_new_background_color"`
	DanmuNewTextColor       string `json:"danmu_new_text_color"`
	DanmuOldBackgroundColor string `json:"danmu_old_background_color"`
	DanmuOldTextColor       string `json:"danmu_old_text_color"`
	DanmuStyle              string `json:"danmu_style"`
	DanmuCustomLeftNew      string `json:"danmu_custom_left_new"`
	DanmuCustomCenterNew    string `json:"danmu_custom_center_new"`
	DanmuCustomRightNew     string `json:"danmu_custom_right_new"`
	DanmuCustomLeftOld      string `json:"danmu_custom_left_old"`
	DanmuCustomCenterOld    string `json:"danmu_custom_center_old"`
	DanmuCustomRightOld     string `json:"danmu_custom_right_old"`
}

// EditSpecialDanmuModel 資料表欄位
type EditSpecialDanmuModel struct {
	ActivityID               string `json:"activity_id"`
	SpecialDanmuMessageCheck string `json:"special_danmu_message_check"`
	SpecialDanmuGeneralPrice string `json:"special_danmu_general_price"`
	SpecialDanmuLargePrice   string `json:"special_danmu_large_price"`
	SpecialDanmuOverPrice    string `json:"special_danmu_over_price"`
	SpecialDanmuTopic        string `json:"special_danmu_topic"`
	SpecialDanmuBanSecond    string `json:"special_danmu_ban_second"`
}

// UpdateDanmu 更新彈幕牆基本設置資料
func (a ActivityModel) UpdateDanmu(model EditDanmuModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"danmu_loop",
			"danmu_top",
			"danmu_mid",
			"danmu_bottom",
			"danmu_display_name",
			"danmu_display_avatar",
			"danmu_size",
			"danmu_speed",
			"danmu_density",
			"danmu_opacity",
			"danmu_background"}

		fieldValues2 = command.Value{}
		fields2      = []string{
			"danmu_new_background_color",
			"danmu_new_text_color",
			"danmu_old_background_color",
			"danmu_old_text_color",

			"danmu_style",
			"danmu_custom_left_new",
			"danmu_custom_center_new",
			"danmu_custom_right_new",
			"danmu_custom_left_old",
			"danmu_custom_center_old",
			"danmu_custom_right_old",
		}

		err error
	)

	if model.DanmuLoop != "" {
		if model.DanmuLoop != "open" && model.DanmuLoop != "close" {
			return errors.New("錯誤: 彈幕循環資料發生問題，請輸入有效的彈幕循環資料")
		}
	}
	if model.DanmuTop != "" {
		if model.DanmuTop != "open" && model.DanmuTop != "close" {
			return errors.New("錯誤: 頂部位置資料發生問題，請輸入有效的頂部位置資料")
		}
	}
	if model.DanmuMid != "" {
		if model.DanmuMid != "open" && model.DanmuMid != "close" {
			return errors.New("錯誤: 中間位置資料發生問題，請輸入有效的中間位置資料")
		}
	}
	if model.DanmuBottom != "" {
		if model.DanmuBottom != "open" && model.DanmuBottom != "close" {
			return errors.New("錯誤: 底部位置資料發生問題，請輸入有效的底部位置資料")
		}
	}
	if model.DanmuDisplayName != "" {
		if model.DanmuDisplayName != "open" && model.DanmuDisplayName != "close" {
			return errors.New("錯誤: 顯示暱稱資料發生問題，請輸入有效的顯示暱稱資料")
		}
	}
	if model.DanmuDisplayAvatar != "" {
		if model.DanmuDisplayAvatar != "open" && model.DanmuDisplayAvatar != "close" {
			return errors.New("錯誤: 顯示人物頭像資料發生問題，請輸入有效的顯示人物頭像資料")
		}
	}
	if model.DanmuSize != "" {
		if sizeFloat, err := strconv.ParseFloat(model.DanmuSize, 32); sizeFloat > 1.5 ||
			sizeFloat < 0.5 || err != nil {
			return errors.New("錯誤: 彈幕大小資料發生問題，請輸入有效的彈幕大小資料")
		}
	}
	if model.DanmuSpeed != "" {
		if speedFloat, err := strconv.ParseFloat(model.DanmuSpeed, 32); speedFloat > 5 ||
			speedFloat < 1 || err != nil {
			return errors.New("錯誤: 彈幕速度資料發生問題，請輸入有效的彈幕速度資料")
		}
	}
	if model.DanmuDensity != "" {
		if densityFloat, err := strconv.ParseFloat(model.DanmuDensity, 32); densityFloat > 4.5 ||
			densityFloat < 0.5 || err != nil {
			return errors.New("錯誤: 彈幕密度資料發生問題，請輸入有效的彈幕密度資料")
		}
	}
	if model.DanmuOpacity != "" {
		if opacityFloat, err := strconv.ParseFloat(model.DanmuOpacity, 32); opacityFloat > 1 ||
			opacityFloat < 0.1 || err != nil {
			return errors.New("錯誤: 彈幕不透明度資料發生問題，請輸入有效的彈幕不透明度資料")
		}
	}
	if model.DanmuBackground != "" {
		if model.DanmuBackground != "1" && model.DanmuBackground != "2" {
			return errors.New("錯誤: 風格選擇資料發生問題，請輸入有效的樣式資料")
		}
	}
	if model.DanmuStyle != "" {
		if model.DanmuStyle != "classic" && model.DanmuStyle != "custom" {
			return errors.New("錯誤: 樣式選擇資料發生問題，請輸入有效的樣式資料")
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

// UpdateSpecialDanmu 更新特殊彈幕牆基本設置資料
func (a ActivityModel) UpdateSpecialDanmu(model EditSpecialDanmuModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"special_danmu_message_check",
			"special_danmu_general_price",
			"special_danmu_large_price",
			"special_danmu_over_price",
			"special_danmu_topic"}

		fieldValues2 = command.Value{}
		fields2      = []string{"special_danmu_ban_second"}

		err error
	)

	if model.SpecialDanmuMessageCheck != "" {
		if model.SpecialDanmuMessageCheck != "open" && model.SpecialDanmuMessageCheck != "close" {
			return errors.New("錯誤: 特殊彈幕訊息審核資料發生問題，請輸入有效的訊息審核資料")
		}
	}
	if model.SpecialDanmuGeneralPrice != "" {
		if _, err := strconv.Atoi(model.SpecialDanmuGeneralPrice); err != nil {
			return errors.New("錯誤: 一般效果資料發生問題，請輸入有效的一般效果資料")
		}
	}
	if model.SpecialDanmuLargePrice != "" {
		if _, err := strconv.Atoi(model.SpecialDanmuLargePrice); err != nil {
			return errors.New("錯誤: 巨大效果資料發生問題，請輸入有效的巨大效果資料")
		}
	}
	if model.SpecialDanmuOverPrice != "" {
		if _, err := strconv.Atoi(model.SpecialDanmuOverPrice); err != nil {
			return errors.New("錯誤: 重疊效果資料發生問題，請輸入有效的重疊效果資料")
		}
	}
	if model.SpecialDanmuTopic != "" {
		if model.SpecialDanmuTopic != "open" && model.SpecialDanmuTopic != "close" {
			return errors.New("錯誤: 主題資料發生問題，請輸入有效的主題資料")
		}
	}
	if model.SpecialDanmuBanSecond != "" {
		if _, err := strconv.Atoi(model.SpecialDanmuBanSecond); err != nil {
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

// values = []string{model.DanmuLoop, model.DanmuTop, model.DanmuMid,
// 	model.DanmuBottom, model.DanmuDisplayName, model.DanmuDisplayAvatar,
// 	model.DanmuSize, model.DanmuSpeed, model.DanmuDensity, model.DanmuOpacity,
// 	model.DanmuBackground}
// values2 = []string{model.DanmuNewBackgroundColor, model.DanmuNewTextColor,
// 	model.DanmuOldBackgroundColor, model.DanmuOldTextColor,
// 	model.DanmuStyle,
// 	model.DanmuCustomLeftNew,
// 	model.DanmuCustomCenterNew,
// 	model.DanmuCustomRightNew,
// 	model.DanmuCustomLeftOld,
// 	model.DanmuCustomCenterOld,
// 	model.DanmuCustomRightOld,}

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

// values = []string{model.SpecialDanmuMessageCheck, model.SpecialDanmuGeneralPrice,
// 	model.SpecialDanmuLargePrice, model.SpecialDanmuOverPrice, model.SpecialDanmuTopic}
// values2      = []string{model.SpecialDanmuBanSecond}

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
