package models

import (
	"errors"

	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
)

// EditGeneralModel 資料表欄位
type EditGeneralModel struct {
	ActivityID           string `json:"activity_id"`
	GeneralDisplayPeople string `json:"general_display_people"`
	GeneralStyle         string `json:"general_style"`
	GeneralBackground    string `json:"general_background"`
	GeneralTopic         string `json:"general_topic"`
	GeneralMusic         string `json:"general_music"`
	GeneralLoop          string `json:"general_loop"`
	GeneralLatest        string `json:"general_latest"`

	// 一般簽到牆自定義圖片
	GeneralClassicHPic01 string `json:"general_classic_h_pic_01" example:"picture"`
	GeneralClassicHPic02 string `json:"general_classic_h_pic_02" example:"picture"`
	GeneralClassicHPic03 string `json:"general_classic_h_pic_03" example:"picture"`
	GeneralClassicHAni01 string `json:"general_classic_h_ani_01" example:"picture"`

	// 音樂
	GeneralBgm string `json:"general_bgm" example:"picture"`
}

// EditThreeDModel 資料表欄位
type EditThreeDModel struct {
	ActivityID             string `json:"activity_id"`
	ThreedAvatar           string `json:"threed_avatar"`
	ThreedAvatarShape      string `json:"threed_avatar_shape"`
	ThreedDisplayPeople    string `json:"threed_display_people"`
	ThreedBackgroundStyle  string `json:"threed_background_style"`
	ThreedBackground       string `json:"threed_background"`
	ThreedImageLogo        string `json:"threed_image_logo"`
	ThreedImageCircle      string `json:"threed_image_circle"`
	ThreedImageSpiral      string `json:"threed_image_spiral"`
	ThreedImageRectangle   string `json:"threed_image_rectangle"`
	ThreedImageSquare      string `json:"threed_image_square"`
	ThreedImageLogoPicture string `json:"threed_image_logo_picture"`
	ThreedTopic            string `json:"threed_topic"`
	ThreedMusic            string `json:"threed_music"`

	// 音樂
	ThreedBgm string `json:"threed_bgm" example:"picture"`
}

// UpdateGeneral 更新一般簽到基本設置資料
func (a ActivityModel) UpdateGeneral(isRedis bool, model EditGeneralModel) error {
	var (
		activityfieldValues  = command.Value{}
		activity2fieldValues = command.Value{"general_edit_times": "general_edit_times + 1"}
		fields               = []string{
			"general_display_people",
			"general_style",
			"general_background"}
		fields2 = []string{
			"general_topic",
			"general_music",
			"general_loop",
			"general_latest",

			"general_classic_h_pic_01",
			"general_classic_h_pic_02",
			"ageneral_classic_h_pic_03",
			"general_classic_h_ani_01",

			"general_bgm",
		}
	)

	if model.GeneralDisplayPeople != "" {
		if model.GeneralDisplayPeople != "open" && model.GeneralDisplayPeople != "close" {
			return errors.New("錯誤: 顯示人數資料發生問題，請輸入有效的資料")
		}
	}
	if model.GeneralStyle != "" {
		// && model.GeneralStyle != "2" && model.GeneralStyle != "3"
		if model.GeneralStyle != "1" {
			return errors.New("錯誤: 展示方式資料發生問題，請輸入有效的展示方式")
		}
	}

	if model.GeneralTopic != "" {
		if model.GeneralTopic != "01_classic" {
			return errors.New("錯誤: 主題資料發生問題，請輸入有效的資料")
		}
	}

	if model.GeneralMusic != "" && (model.GeneralMusic != "classic" && model.GeneralMusic != "customize") {
		return errors.New("錯誤: 音樂資料發生問題，請輸入有效的音樂")
	}

	if model.GeneralLoop != "" {
		if model.GeneralLoop != "open" && model.GeneralLoop != "close" {
			return errors.New("錯誤: 是否輪播資料發生問題，請輸入有效的資料")
		}
	}

	if model.GeneralLatest != "" {
		if model.GeneralLatest != "open" && model.GeneralLatest != "close" {
			return errors.New("錯誤: 跳轉至最新簽到資料發生問題，請輸入有效的資料")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			activityfieldValues[key] = val
		}
	}
	if len(activityfieldValues) == 0 {
		return nil
	}

	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			activity2fieldValues[key] = val
		}
	}
	if len(activity2fieldValues) == 0 {
		return nil
	}

	// 清除簽到牆設置資料(redis)
	if isRedis {
		a.RedisConn.DelCache(config.ACTIVITY_REDIS + model.ActivityID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GENERAL_EDIT_TIMES_REDIS+model.ActivityID, "修改資料")
	}

	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(activityfieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新資料發生問題(activity)")
	}

	if err := a.Table(config.ACTIVITY_2_TABLE).
		Where("activity_id", "=", model.ActivityID).Update(activity2fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新資料發生問題(activity_2)")
	}

	return nil
}

// UpdateThreed 更新立體簽到基本設置資料
func (a ActivityModel) UpdateThreed(isRedis bool, model EditThreeDModel) error {
	var (
		activityfieldValues  = command.Value{}
		activity2fieldValues = command.Value{"threed_edit_times": "threed_edit_times + 1"}
		fields               = []string{
			"threed_avatar",
			"threed_avatar_shape",
			"threed_display_people",
			"threed_background_style",
			"threed_background",
			"threed_image_logo",
			"threed_image_circle",
			"threed_image_spiral",
			"threed_image_rectangle",
			"threed_image_square",
			"threed_image_logo_picture"}
		fields2 = []string{
			"threed_topic",
			"threed_music",
			"threed_bgm",
		}
	)

	if model.ThreedAvatarShape != "" {
		if model.ThreedAvatarShape != "circle" && model.ThreedAvatarShape != "square" {
			return errors.New("錯誤: 頭像形狀資料發生問題，請輸入有效的頭像形狀資料")
		}
	}
	if model.ThreedDisplayPeople != "" {
		if model.ThreedDisplayPeople != "open" && model.ThreedDisplayPeople != "close" {
			return errors.New("錯誤: 顯示人數資料發生問題，請輸入有效的資料")
		}
	}
	if model.ThreedBackgroundStyle != "" {
		if model.ThreedBackgroundStyle != "1" && model.ThreedBackgroundStyle != "2" &&
			model.ThreedBackgroundStyle != "3" {
			return errors.New("錯誤: 背景樣式資料發生問題，請輸入有效的背景樣式")
		}
	}
	// LOGO
	if model.ThreedImageLogo != "" {
		if model.ThreedImageLogo != "open" && model.ThreedImageLogo != "close" {
			return errors.New("錯誤: 3D頭像(LOGO)資料發生問題，請輸入有效的資料")
		}
	}
	// 圓形
	if model.ThreedImageCircle != "" {
		if model.ThreedImageCircle != "open" && model.ThreedImageCircle != "close" {
			return errors.New("錯誤: 3D頭像(圓形)資料發生問題，請輸入有效的資料")
		}
	}
	// 螺旋
	if model.ThreedImageSpiral != "" {
		if model.ThreedImageSpiral != "open" && model.ThreedImageSpiral != "close" {
			return errors.New("錯誤: 3D頭像(螺旋)資料發生問題，請輸入有效的資料")
		}
	}
	// 矩形
	if model.ThreedImageRectangle != "" {
		if model.ThreedImageRectangle != "open" && model.ThreedImageRectangle != "close" {
			return errors.New("錯誤: 3D頭像(矩形)資料發生問題，請輸入有效的資料")
		}
	}
	// 正方形
	if model.ThreedImageSquare != "" {
		if model.ThreedImageSquare != "open" && model.ThreedImageSquare != "close" {
			return errors.New("錯誤: 3D頭像(正方形)資料發生問題，請輸入有效的資料")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			activityfieldValues[key] = val
		}
	}
	if len(activityfieldValues) == 0 {
		return nil
	}

	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			activity2fieldValues[key] = val
		}
	}
	if len(activity2fieldValues) == 0 {
		return nil
	}

	// 清除簽名牆設置資料(redis)
	if isRedis {
		a.RedisConn.DelCache(config.ACTIVITY_REDIS + model.ActivityID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_THREED_EDIT_TIMES_REDIS+model.ActivityID, "修改資料")
	}

	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(activityfieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新資料發生問題(activity)")
	}

	if err := a.Table(config.ACTIVITY_2_TABLE).
		Where("activity_id", "=", model.ActivityID).Update(activity2fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新資料發生問題(activity_2)")
	}

	return nil
}

// values  = []string{model.GeneralDisplayPeople, model.GeneralStyle, model.GeneralBackground}
// values2 = []string{
// 	model.GeneralTopic, model.GeneralMusic,
// 	model.GeneralLoop, model.GeneralLatest,

// 	model.GeneralClassicHPic01,
// 	model.GeneralClassicHPic02,
// 	model.GeneralClassicHPic03,
// 	model.GeneralClassicHAni01,

// 	model.GeneralBgm,
// }

// for i, value := range values {
// 	if value != "" {
// 		activityfieldValues[fields[i]] = value
// 	}
// }

// for i, value2 := range values2 {
// 	if value2 != "" {
// 		activity2fieldValues[fields2[i]] = value2
// 	}
// }

// values = []string{model.ThreedAvatar, model.ThreedAvatarShape, model.ThreedDisplayPeople,
// 	model.ThreedBackgroundStyle, model.ThreedBackground, model.ThreedImageLogo,
// 	model.ThreedImageCircle, model.ThreedImageSpiral, model.ThreedImageRectangle, model.ThreedImageSquare,
// 	model.ThreedImageLogoPicture}
// values2 = []string{
// 	model.ThreedTopic, model.ThreedMusic,
// 	model.ThreedBgm,
// }

// for i, value := range values {
// 	if value != "" {
// 		activityfieldValues[fields[i]] = value
// 	}
// }
// if len(activityfieldValues) == 0 {
// 	return nil
// }

// for i, value2 := range values2 {
// 	if value2 != "" {
// 		activity2fieldValues[fields2[i]] = value2
// 	}
// }
// if len(activity2fieldValues) == 0 {
// 	return nil
// }
