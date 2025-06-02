package models

import (
	"errors"
	"strconv"

	"hilive/modules/db/command"
)

// EditPictureModel 資料表欄位
type EditPictureModel struct {
	ActivityID          string `json:"activity_id"`
	PictureStartTime    string `json:"picture_start_time"`
	PictureEndTime      string `json:"picture_end_time"`
	PictureHideTime     string `json:"picture_hide_time"`
	PictureSwitchSecond string `json:"picture_switch_second"`
	PicturePlayOrder    string `json:"picture_play_order"`
	Picture             string `json:"picture"`
	PictureBackground   string `json:"picture_background"`
	PictureAnimation    string `json:"picture_animation"`
}

// UpdatePictrure 更新圖片牆基本設置資料
func (a ActivityModel) UpdatePictrure(model EditPictureModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"picture_hide_time", "picture_switch_second", "picture_play_order",
			"picture", "picture_background", "picture_animation"}
		values = []string{model.PictureHideTime, model.PictureSwitchSecond, model.PicturePlayOrder,
			model.Picture, model.PictureBackground, model.PictureAnimation}
	)

	if model.PictureStartTime != "" && model.PictureEndTime != "" {
		if !CompareTDatetime(model.PictureStartTime, model.PictureEndTime) {
			return errors.New("錯誤: 時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
		}
		fieldValues["picture_start_time"] = model.PictureStartTime
		fieldValues["picture_end_time"] = model.PictureEndTime
	}
	if model.PictureHideTime != "" {
		if model.PictureHideTime != "open" && model.PictureHideTime != "close" {
			return errors.New("錯誤: 隱藏時間資料發生問題，請輸入有效的資料")
		}
	}
	if model.PictureSwitchSecond != "" {
		if _, err := strconv.Atoi(model.PictureSwitchSecond); err != nil {
			return errors.New("錯誤: 切換圖片秒數資料發生問題，請輸入有效的秒數")
		}
	}
	if model.PicturePlayOrder != "" {
		if model.PicturePlayOrder != "order" && model.PicturePlayOrder != "random" {
			return errors.New("錯誤: 播放順序資料發生問題，請輸入有效的資料")
		}
	}
	if model.PictureAnimation != "" {
		if model.PictureAnimation != "1" && model.PictureAnimation != "2" &&
			model.PictureAnimation != "3" {
			return errors.New("錯誤: 動畫樣式資料發生問題，請輸入有效的樣式")
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
