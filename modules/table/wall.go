package table

import (
	"encoding/json"
	"errors"

	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"hilive/template/form"

	"github.com/gin-gonic/gin"
)

// GetMessagePanel 訊息牆
func (s *SystemTable) GetMessagePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}
		var background string
		if values.Get("message_background"+DEFAULT_FALG) == "1" {
			background = UPLOAD_SYSTEM_URL + "img-topic-bg.jpg"
		} else if values.Get("message_background") != "" {
			background = values.Get("message_background")
		} else {
			background = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMessageModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.MessageBackground = background

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateMessage(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetMessageCheckPanel 訊息審核
func (s *SystemTable) GetMessageCheckPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMessageCheckModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateMessageCheck(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetMessageSensitivityPanel 訊息敏感詞
func (s *SystemTable) GetMessageMessageSensitivityPanel() (table Table) {
	// fmt.Println("執行GetMessageMessageSensitivityPanel")
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	// 刪除
	info.SetTable(config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	formList := table.GetForm()
	// 新增敏感詞
	formList.SetTable(config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE).
		SetInsertFunc(func(values form2.Values) error {
			if values.IsEmpty("activity_id", "sensitivity_word") {
				return errors.New("錯誤: 活動ID、敏感詞資料發生問題，請輸入有效的資料")
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditMessageSensitivityModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultMessageSensitivityModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Add(model); err != nil {
				return err
			}
			return nil
		})

	// 編輯敏感詞
	formList.SetTable(config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE).
		SetUpdateFunc(func(values form2.Values) error {
			if values.IsEmpty("id") {
				return errors.New("錯誤: ID資料發生問題，請輸入有效的ID資料")
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditMessageSensitivityModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultMessageSensitivityModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
			return nil
		})
	return
}

// GetTopicPanel 主題牆
func (s *SystemTable) GetTopicPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.AddField("背景圖片", "topic_background", db.Varchar, form.File)

	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}
		var background string
		if values.Get("topic_background"+DEFAULT_FALG) == "1" {
			background = UPLOAD_SYSTEM_URL + "img-topic-bg.jpg"
		} else if values.Get("topic_background") != "" {
			background = values.Get("topic_background")
		} else {
			background = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditTopicModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.TopicBackground = background

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateTopic(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetQuestionPanel 提問牆
func (s *SystemTable) GetQuestionPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 提問嘉賓資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_QUESTION_GUEST_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_QUESTION_GUEST_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 基本設置表單資訊
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		var background string
		if values.Get("question_background"+DEFAULT_FALG) == "1" {
			background = UPLOAD_SYSTEM_URL + "img-request-bg.png"
		} else if values.Get("question_background") != "" {
			background = values.Get("question_background")
		} else {
			background = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditQuestionModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.QuestionBackground = background

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateQuestion(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetDanmuPanel 彈幕
func (s *SystemTable) GetDanmuPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		var (
			pictures = []PictureField{
				{FieldName: "danmu_custom_left_new", Path: "danmu/custom/left_new.svg"},
				{FieldName: "danmu_custom_center_new", Path: "danmu/custom/center_new.svg"},
				{FieldName: "danmu_custom_right_new", Path: "danmu/custom/right_new.svg"},
				{FieldName: "danmu_custom_left_old", Path: "danmu/custom/left_old.svg"},
				{FieldName: "danmu_custom_center_old", Path: "danmu/custom/center_old.svg"},
				{FieldName: "danmu_custom_right_old", Path: "danmu/custom/right_old.svg"},
			}
		)

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditDanmuModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(pictures, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateDanmu(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetSpecialDanmuPanel 特殊彈幕
func (s *SystemTable) GetSepcialDanmuPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).
		SetUpdateFunc(func(values form2.Values) error {
			if values.IsEmpty("activity_id") {
				return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditSpecialDanmuModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultActivityModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				UpdateSpecialDanmu(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
			return nil
		})
	return
}

// GetHoldScreenPanel 霸屏
func (s *SystemTable) GetHoldScreenPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: id不能為空")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditHoldScreenModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateHoldScreen(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetPicturePanel 圖片牆
// func (s *SystemTable) GetPicturePanel() (table Table) {
// 	table = DefaultTable(DefaultConfig())
// 	settingFormList := table.GetSettingForm()
// 	settingFormList.AddField("活動開始", "picture_start_time", db.Datetime, form.DatetimeRange)
// 	// .
// 	// SetDisplayFunc(func(model types.FieldModel) interface{} {
// 	// 	var timeRange interface{}
// 	// 	if model.ID == "" {
// 	// 		return timeRange
// 	// 	}
// 	// 	res, _ := s.table(config.ACTIVITY_TABLE).Select("picture_start_time, picture_end_time").
// 	// 		Find("activity_id", model.ID)
// 	// 	start := strings.Split(fmt.Sprintf("%s", res["picture_start_time"]), ":")
// 	// 	end := strings.Split(fmt.Sprintf("%s", res["picture_end_time"]), ":")
// 	// 	return fmt.Sprintf("%s,%s", fmt.Sprintf("%s:%s", start[0], start[1]),
// 	// 		fmt.Sprintf("%s:%s", end[0], end[1]))
// 	// })
// 	settingFormList.AddField("隱藏時間", "picture_hide_time", db.Varchar, form.Checkbox)
// 	settingFormList.AddField("切換圖片秒數", "picture_switch_second", db.Int, form.Number)
// 	settingFormList.AddField("播放順序", "picture_play_order", db.Varchar, form.Radio).
// 		SetFieldOptions(types.FieldOptions{
// 			{Text: "順序播放", Value: "order"}, {Text: "隨機播放", Value: "random"},
// 		})
// 	settingFormList.AddField("圖片", "picture", db.Varchar, form.Multifile).
// 		SetFieldHelpMsg(template.HTML("選取多個圖片輪流播放"))
// 	// .
// 	// 	SetDisplayFunc(func(model types.FieldModel) interface{} {
// 	// 		res, _ := s.table(config.ACTIVITY_TABLE).Select("picture").Find("activity_id", model.ID)
// 	// 		return strings.Split(fmt.Sprintf("%s", res["picture"]), "\n")
// 	// 	})
// 	settingFormList.AddField("背景", "picture_background", db.Varchar, form.File)
// 	settingFormList.AddField("動畫樣式", "picture_animation", db.Int, form.Radio).SetFieldOptions(types.FieldOptions{
// 		{Text: "樣式1", Value: "1", Text2: UPLOAD_RADIO_URL + "picture-background-style1.png"},
// 		{Text: "樣式2", Value: "2", Text2: UPLOAD_RADIO_URL + "picture-background-style2.png"},
// 		{Text: "樣式3", Value: "3", Text2: UPLOAD_RADIO_URL + "picture-background-style3.png"},
// 	})

// 	// 未使用
// 	settingFormList.AddField("picture_end_time", "picture_end_time", db.Datetime, form.NoUse)

// 	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
// 		if values.IsEmpty("activity_id") {
// 			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
// 		}
// 		var pictures string
// 		pics := values.Gets("picture")
// 		if len(pics) > 0 {
// 			for i := 0; i < len(pics); i++ {
// 				if i == len(pics)-1 {
// 					pictures += pics[i]
// 				} else {
// 					pictures += pics[i] + "\n"
// 				}
// 			}
// 		} else {
// 			pictures = ""
// 		}
// 		var background string
// 		// 圖片儲存需要處理
// 		if values.Get("picture_background"+DEFAULT_FALG) == "1" {
// 			background = UPLOAD_SYSTEM_URL + "img-picture-bg.png"
// 		} else if values.Get("picture_background") != "" {
// 			background = values.Get("picture_background")
// 		} else {
// 			background = ""
// 		}

// 		if err := models.DefaultActivityModel().SetDbConn(s.dbConn).
// 			UpdatePictrure(models.EditPictureModel{
// 				ActivityID:          values.Get("activity_id"),
// 				PictureStartTime:    values.Get("picture_start_time"),
// 				PictureEndTime:      values.Get("picture_end_time"),
// 				PictureHideTime:     values.Get("picture_hide_time"),
// 				PictureSwitchSecond: values.Get("picture_switch_second"),
// 				PicturePlayOrder:    values.Get("picture_play_order"),
// 				Picture:             pictures,
// 				PictureBackground:   background,
// 				PictureAnimation:    values.Get("picture_animation"),
// 			}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 			return err
// 		}
// 		return nil
// 	})
// return
// }

// 提問嘉賓表單資訊
// formList := table.GetForm()
// formList.SetTable(config.ACTIVITY_QUESTION_GUEST_TABLE).SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "name") {
// 		return errors.New("錯誤: 活動ID、姓名發生問題，請輸入有效的資料")
// 	}

// 	var avatar string
//  圖片儲存需要處理
// 	if values.Get("avatar") != "" {
// 		avatar = values.Get("avatar")
// 	} else {
// 		avatar = UPLOAD_SYSTEM_URL + "img-guest-pic.png"
// 	}

// 	if err := models.DefaultQuestionGuestModel().SetDbConn(s.dbConn).Add(
// 		models.EditQuestionGuestModel{
// 			ActivityID: values.Get("activity_id"),
// 			Name:       values.Get("name"),
// 			Avatar:     avatar,
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// })

// formList.SetUpdateFunc(func(values form2.Values) error {
// 	if values.IsEmpty("id", "activity_id") {
// 		return errors.New("錯誤: ID發生問題，請輸入有效的資料")
// 	}
// 	var avatar string
//  圖片儲存需要處理
// 	if values.Get("avatar"+DEFAULT_FALG) == "1" {
// 		avatar = UPLOAD_SYSTEM_URL + "img-guest-pic.png"
// 	} else if values.Get("avatar") != "" {
// 		avatar = values.Get("avatar")
// 	} else {
// 		avatar = ""
// 	}

// 	if err := models.DefaultQuestionGuestModel().SetDbConn(s.dbConn).
// 		Update(models.EditQuestionGuestModel{
// 			ID:         values.Get("id"),
// 			ActivityID: values.Get("activity_id"),
// 			Name:       values.Get("name"),
// 			Avatar:     avatar,
// 		}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// 	return nil
// })

// @Summary 新增訊息敏感詞資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param sensitivity_word formData string true "敏感詞" minlength(1)
// @param replace_word formData string false "替代詞"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message_sensitivity [post]
func POSTMessageSensitivity(ctx *gin.Context) {
}

// @Summary 新增提問嘉賓資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param name formData string true "姓名(上限為10個字元)" minlength(1) maxlength(10)
// @param avatar formData file false "頭像"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/question/guest [post]
func POSTQuestionGuest(ctx *gin.Context) {
}

// @Summary 訊息牆設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param message_picture formData string false "訊息區可以上傳圖片" Enums(open, close)
// @param message_auto formData string false "圖片自動放大顯示" Enums(open, close)
// @param message_ban formData string false "禁止連續傳送訊息" Enums(open, close)
// @param message_ban_second formData integer false "禁止秒數"
// @param message_refresh_second formData integer false "訊息換頁秒數"
// @param message_open formData string false "跑馬燈開啟" Enums(open, close)
// @param message formData string false "展示訊息設置(一行設置一個跑馬燈訊息，若要輸入新的跑馬燈請換行)"
// @param message_background formData file false "背景圖片"
// @param message_text_color formData string false "訊息文字顏色"
// @param message_screen_title_color formData string false "大螢幕標題顏色"
// @param message_text_frame_color formData string false "訊息文字外框顏色"
// @param message_frame_color formData string false "訊息框顏色"
// @param message_screen_title_frame_color formData string false "大螢幕標題邊框顏色"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message [patch]
func PATCHMessage(ctx *gin.Context) {
}

// @Summary 訊息審核設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param message_check_manual_check formData string false "手動審核" Enums(open, close)
// @param message_check_sensitivity formData string false "敏感詞" Enums(replace, refuse)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message_check [patch]
func PATCHMessageCheck(ctx *gin.Context) {
}

// @Summary 編輯訊息敏感詞資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID"
// @param sensitivity_word formData string false "敏感詞"
// @param replace_word formData string false "替代詞"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message_sensitivity [put]
func PUTMessageSensitivity(ctx *gin.Context) {
}

// @Summary 主題牆設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param topic_background formData file false "背景圖片"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/topic [patch]
func PATCHTopic(ctx *gin.Context) {
}

// @Summary 提問牆設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param question_anonymous formData string false "匿名提問" Enums(open, close)
// @param question_qrcode formData string false "二維碼" Enums(open, close)
// @param question_background formData file false "背景圖片"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/question [patch]
func PATCHQuestion(ctx *gin.Context) {
}

// @Summary 編輯提問嘉賓資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param name formData string false "姓名(上限為10個字元)" minlength(1) maxlength(10)
// @param avatar formData file false "頭像"
// @param token formData string false "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/question/guest [put]
func PUTQuestionGuest(ctx *gin.Context) {
}

// @Summary 彈幕設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param danmu_loop formData string false "彈幕循環" Enums(open, close)
// @param danmu_top formData string false "頂部位置" Enums(open, close)
// @param danmu_mid formData string false "中間位置" Enums(open, close)
// @param danmu_bottom formData string false "底部位置" Enums(open, close)
// @param danmu_display_name formData string false "顯示暱稱" Enums(open, close)
// @param danmu_display_avatar formData string false "顯示人物頭像" Enums(open, close)
// @param danmu_size formData integer false "彈幕大小(1-200)" minimum(1) maximum(200)
// @param danmu_speed formData integer false "彈幕速度(1-200)" minimum(1) maximum(200)
// @param danmu_density formData integer false "彈幕密度(1-100)" minimum(1) maximum(100)
// @param danmu_opacity formData integer false "彈幕不透明度(1-100)" minimum(1) maximum(100)
// @param danmu_background formData integer false "樣式選擇" Enums(1)
// @param danmu_new_background_color formData string false "danmu_new_background_color"
// @param danmu_new_text_color formData string false "danmu_new_text_color"
// @param danmu_old_background_color formData string false "danmu_old_background_color"
// @param danmu_old_text_color formData string false "danmu_old_text_color"
// @param danmu_style formData string false "danmu_style"
// @param danmu_custom_left_new formData string false "danmu_custom_left_new"
// @param danmu_custom_center_new formData string false "danmu_custom_center_new"
// @param danmu_custom_right_new formData string false "danmu_custom_right_new"
// @param danmu_custom_left_old formData string false "danmu_custom_left_old"
// @param danmu_custom_center_old formData string false "danmu_custom_center_old"
// @param danmu_custom_right_old formData string false "danmu_custom_right_old"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/danmu [patch]
func PATCHDanmu(ctx *gin.Context) {
}

// @Summary 特殊彈幕設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param special_danmu_message_check formData string false "特殊彈幕訊息審核" Enums(open, close)
// @param special_danmu_general_price formData integer false "一般效果"
// @param special_danmu_large_price formData integer false "巨大效果"
// @param special_danmu_over_price formData integer false "重疊效果"
// @param special_danmu_topic formData string false "主題" Enums(open, close)
// @param special_danmu_ban_second formData integer false "special_danmu_ban_second"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/specialdanmu [patch]
func PATCHSuperDanmu(ctx *gin.Context) {
}

// @Summary 圖片牆設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param picture_start_time formData string false "開始時間(西元年-月-日T時:分)"
// @param picture_end_time formData string false "結束時間(西元年-月-日T時:分)"
// @param picture_hide_time formData string false "隱藏時間" Enums(open, close)
// @param picture_switch_second formData integer false "切換圖片秒數"
// @param picture_play_order formData string false "播放順序" Enums(order, random)
// @param picture formData file false "圖片"
// @param picture_background formData file false "背景"
// @param picture_animation formData integer false "動畫樣式" Enums(1, 2, 3)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/picture [patch]
func PATCHPicture(ctx *gin.Context) {
}

// @Summary 霸屏設置
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param holdscreen_message_check formData string false "霸佔彈幕訊息審核" Enums(open, close)
// @param holdscreen_only_picture formData string false "只允許以此類型彈幕發送圖片" Enums(open, close)
// @param holdscreen_price formData integer false "彈幕展示價格(秒)"
// @param holdscreen_display_second formData integer false "彈幕展示時間"
// @param holdscreen_birthday_topic formData string false "生日主題" Enums(open, close)
// @param holdscreen_ban_second formData integer false "holdscreen_ban_second"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/holdscreen [patch]
func PATCHHoldscreen(ctx *gin.Context) {
}

// @Summary 刪除訊息敏感詞資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message_sensitivity [delete]
func DELETEMessageSensitivity(ctx *gin.Context) {
}

// @Summary 刪除提問嘉賓資料
// @Tags Wall
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/question/guest [delete]
func DELETEQuestionGuest(ctx *gin.Context) {
}

// @Summary 訊息敏感詞JSON資料
// @Tags Wall
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/wall/message_sensitivity [get]
func MessageSensitivityJSON(ctx *gin.Context) {
}

// models.EditMessageModel{
// 	ActivityID:                   values.Get("activity_id"),
// 	MessagePicture:               values.Get("message_picture"),
// 	MessageAuto:                  values.Get("message_auto"),
// 	MessageBan:                   values.Get("message_ban"),
// 	MessageBanSecond:             values.Get("message_ban_second"),
// 	MessageRefreshSecond:         values.Get("message_refresh_second"),
// 	MessageOpen:                  values.Get("message_open"),
// 	Message:                      values.Get("message"),
// 	MessageBackground:            background,
// 	MessageTextColor:             values.Get("message_text_color"),
// 	MessageScreenTitleColor:      values.Get("message_screen_title_color"),
// 	MessageTextFrameColor:        values.Get("message_text_frame_color"),
// 	MessageFrameColor:            values.Get("message_frame_color"),
// 	MessageScreenTitleFrameColor: values.Get("message_screen_title_frame_color"),
// }

// models.EditMessageCheckModel{
// 	ActivityID:              values.Get("activity_id"),
// 	MessageCheckManualCheck: values.Get("message_check_manual_check"),
// 	MessageCheckSensitivity: values.Get("message_check_sensitivity"),
// }

// models.EditMessageSensitivityModel{
// 	ActivityID:      values.Get("activity_id"),
// 	SensitivityWord: values.Get("sensitivity_word"),
// 	ReplaceWord:     values.Get("replace_word"),
// }

// models.EditMessageSensitivityModel{
// 	ID:              values.Get("id"),
// 	SensitivityWord: values.Get("sensitivity_word"),
// 	ReplaceWord:     values.Get("replace_word"),
// }

// models.EditTopicModel{
// 	ActivityID:      values.Get("activity_id"),
// 	TopicBackground: background,
// }

// models.EditQuestionModel{
// 	ActivityID: values.Get("activity_id"),
// 	// QuestionMessageCheck:  values.Get("question_message_check"),
// 	QuestionAnonymous: values.Get("question_anonymous"),
// 	// QuestionHideAnswered:  values.Get("question_hide_answered"),
// 	QuestionQrcode:     values.Get("question_qrcode"),
// 	QuestionBackground: background,
// 	// QuestionGuestAnswered: values.Get("question_guest_answered"),
// }

// pics = []string{
// 	"danmu/custom/left_new.svg",
// 	"danmu/custom/center_new.svg",
// 	"danmu/custom/right_new.svg",
// 	"danmu/custom/left_old.svg",
// 	"danmu/custom/center_old.svg",
// 	"danmu/custom/right_old.svg",
// }
// fields = []string{
// 	"danmu_custom_left_new",
// 	"danmu_custom_center_new",
// 	"danmu_custom_right_new",
// 	"danmu_custom_left_old",
// 	"danmu_custom_center_old",
// 	"danmu_custom_right_old",
// }
// update = make([]string, 100)

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditDanmuModel{
// 	ActivityID:              values.Get("activity_id"),
// 	DanmuLoop:               values.Get("danmu_loop"),
// 	DanmuTop:                values.Get("danmu_top"),
// 	DanmuMid:                values.Get("danmu_mid"),
// 	DanmuBottom:             values.Get("danmu_bottom"),
// 	DanmuDisplayName:        values.Get("danmu_display_name"),
// 	DanmuDisplayAvatar:      values.Get("danmu_display_avatar"),
// 	DanmuSize:               values.Get("danmu_size"),
// 	DanmuSpeed:              values.Get("danmu_speed"),
// 	DanmuDensity:            values.Get("danmu_density"),
// 	DanmuOpacity:            values.Get("danmu_opacity"),
// 	DanmuBackground:         values.Get("danmu_background"),
// 	DanmuNewBackgroundColor: values.Get("danmu_new_background_color"),
// 	DanmuNewTextColor:       values.Get("danmu_new_text_color"),
// 	DanmuOldBackgroundColor: values.Get("danmu_old_background_color"),
// 	DanmuOldTextColor:       values.Get("danmu_old_text_color"),

// 	DanmuStyle:           values.Get("danmu_style"),
// 	DanmuCustomLeftNew:   update[0],
// 	DanmuCustomCenterNew: update[1],
// 	DanmuCustomRightNew:  update[2],
// 	DanmuCustomLeftOld:   update[3],
// 	DanmuCustomCenterOld: update[4],
// 	DanmuCustomRightOld:  update[5],
// }

// models.EditSpecialDanmuModel{
// 	ActivityID:               values.Get("activity_id"),
// 	SpecialDanmuMessageCheck: values.Get("special_danmu_message_check"),
// 	SpecialDanmuGeneralPrice: values.Get("special_danmu_general_price"),
// 	SpecialDanmuLargePrice:   values.Get("special_danmu_large_price"),
// 	SpecialDanmuOverPrice:    values.Get("special_danmu_over_price"),
// 	SpecialDanmuTopic:        values.Get("special_danmu_topic"),
// 	SpecialDanmuBanSecond:    values.Get("special_danmu_ban_second"),
// }

// models.EditHoldScreenModel{
// 	ActivityID:              values.Get("activity_id"),
// 	HoldScreenPrice:         values.Get("holdscreen_price"),
// 	HoldScreenMessageCheck:  values.Get("holdscreen_message_check"),
// 	HoldScreenOnlyPicture:   values.Get("holdscreen_only_picture"),
// 	HoldScreenDisplaySecond: values.Get("holdscreen_display_second"),
// 	HoldScreenBirthdayTopic: values.Get("holdscreen_birthday_topic"),
// 	HoldscreenBanSecond:     values.Get("holdscreen_ban_second"),
// }
