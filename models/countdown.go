package models

import (
	"errors"
	"strconv"

	"hilive/modules/db/command"
)

// EditCountdownModel 資料表欄位
type EditCountdownModel struct {
	ActivityID           string `json:"activity_id"`
	CountdownSecond      string `json:"countdown_second"`
	CountdownURL         string `json:"countdown_url"`
	CountdownAvatar      string `json:"countdown_avatar"`
	CountdownAvatarShape string `json:"countdown_avatar_shape"`
	CountdownBackground  string `json:"countdown_background"`
}

// UpdateCountdown 更新倒數計時牆基本設置資料
func (a ActivityModel) UpdateCountdown(model EditCountdownModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"countdown_second", "countdown_url", "countdown_avatar",
			"countdown_avatar_shape", "countdown_background"}
		values = []string{model.CountdownSecond, model.CountdownURL, model.CountdownAvatar,
			model.CountdownAvatarShape, model.CountdownBackground}
	)

	if model.CountdownSecond != "" {
		if _, err := strconv.Atoi(model.CountdownSecond); err != nil {
			return errors.New("錯誤: 倒數秒數資料發生問題，請輸入有效的倒數秒數資料")
		}
	}
	if model.CountdownURL != "" {
		if model.CountdownURL != "current" && model.CountdownURL != "threed" {
			return errors.New("錯誤: 倒數後進入頁面資料發生問題，請輸入有效的倒數後進入頁面資料")
		}
	}
	if model.CountdownAvatarShape != "" {
		if model.CountdownAvatarShape != "circle" && model.CountdownAvatarShape != "square" {
			return errors.New("錯誤: 頭像形狀資料發生問題，請輸入有效的頭像形狀資料")
		}
	}
	if model.CountdownBackground != "" {
		if model.CountdownBackground != "1" && model.CountdownBackground != "2" &&
			model.CountdownBackground != "3" {
			return errors.New("錯誤: 背景資料發生問題，請輸入有效的背景資料")
		}
	}

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}
	return a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}
