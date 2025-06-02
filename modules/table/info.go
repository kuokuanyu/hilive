package table

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// GetOverviewPanel 活動總覽
func (s *SystemTable) GetOverviewPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_OVERVIEW_GAME_TABLE)

	// 表單資訊
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
			var model models.EditActivityOverviewModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultActivityModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				UpdateOverview(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
			return nil
		})
	return
}

// GetIntroducePanel 介紹
func (s *SystemTable) GetIntroducePanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 面板資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_INTRODUCE_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_INTRODUCE_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_INTRODUCE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID、標題發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditIntroduceModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultIntroduceModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "activity_id") {
			return errors.New("錯誤: ID、活動ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditIntroduceModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultIntroduceModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
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

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditIntroduceSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateIntroduce(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetSchedulePanel 行程
func (s *SystemTable) GetSchedulePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_SCHEDULE_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_SCHEDULE_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_SCHEDULE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "title", "schedule_date",
			"start_time", "end_time") {
			return errors.New("錯誤: ID、事件名稱、日期、時間發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditScheduleModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultScheduleModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}

		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "activity_id") {
			return errors.New("錯誤: ID、活動ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditScheduleModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultScheduleModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})

	// 基本設置的表單資訊
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
		var model models.EditScheduleSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateSchdule(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetGuestPanel 嘉賓
func (s *SystemTable) GetGuestPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_GUEST_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		if err := s.table(config.ACTIVITY_GUEST_TABLE).
			WhereIn("id", interfaces(idArr)).Delete(); err != nil {
			return err
		}
		return nil
	})

	// 表單資訊
	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_GUEST_TABLE).
		SetInsertFunc(func(values form2.Values) error {
			if values.IsEmpty("activity_id", "name") {
				return errors.New("錯誤: 活動ID、嘉賓名稱資料發生問題，請輸入有效的資料")
			}

			var avatar string
			if values.Get("avatar") != "" {
				avatar = values.Get("avatar")
			} else {
				avatar = UPLOAD_SYSTEM_URL + "img-guest-pic.png"
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditGuestModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			// 手動處理
			model.Avatar = avatar

			if err := models.DefaultGuestModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Add(model); err != nil {
				return err
			}
			return nil
		})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "activity_id") {
			return errors.New("錯誤: ID、活動ID發生問題，請輸入有效的資料")
		}
		var avatar string
		if values.Get("avatar"+DEFAULT_FALG) == "1" {
			avatar = UPLOAD_SYSTEM_URL + "img-guest-pic.png"
		} else if values.Get("avatar") != "" {
			avatar = values.Get("avatar")
		} else {
			avatar = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGuestModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.Avatar = avatar

		if err := models.DefaultGuestModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
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
		if values.Get("guest_background"+DEFAULT_FALG) == "1" {
			background = UPLOAD_SYSTEM_URL + "img-guest-bg.png"
		} else if values.Get("guest_background") != "" {
			background = UPLOAD_URL + values.Get("guest_background")
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
		var model models.EditGuestSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.GuestBackground = background

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateGuest(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}

		return nil
	})
	return
}

// GetMaterialPanel 資料
func (s *SystemTable) GetMaterialPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_MATERIAL_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_MATERIAL_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_MATERIAL_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "title") {
			return errors.New("錯誤: 活動ID、資料標題發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMaterialModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultMaterialModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "activity_id") {
			return errors.New("錯誤: ID、活動ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMaterialModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultMaterialModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
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

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMaterialSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateMaterial(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增活動介紹資料
// @Tags Introduce
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce_type formData string true "介紹類型" Enums(text, picture)
// @param content formData string false "內容(上限為200個字元)" maxlength(200)
// @param introduce_order formData integer true "介紹排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/introduce [post]
func POSTIntroduce(ctx *gin.Context) {
}

// @Summary 新增活動行程資料
// @Tags Schedule
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "主題(上限為20個字元)" minlength(1) maxlength(20)
// @param content formData string false "內容(上限為200個字元)" maxlength(200)
// @param schedule_date formData string true "日期(西元年-月-日)"
// @param start_time formData string true "開始時間(時:分)"
// @param end_time formData string true "結束時間(時:分)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/schedule [post]
func POSTSchedule(ctx *gin.Context) {
}

// @Summary 新增活動嘉賓資料
// @Tags Guest
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param name formData string true "姓名(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce formData string false "簡介(上限為20個字元)" maxlength(20)
// @param detail formData string false "詳情(上限為200個字元)" maxlength(200)
// @param avatar formData file false "照片"
// @param guest_order formData integer true "嘉賓排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/guest [post]
func POSTGuest(ctx *gin.Context) {
}

// @Summary 新增活動資料
// @Tags Material
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "資料標題(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce formData string false "資料說明(上限為200個字元)" maxlength(200)
// @param link formData string false "資料連結"
// @param material_order formData integer true "資料排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/material [post]
func POSTMaterial(ctx *gin.Context) {
}

// @Summary 活動總覽設置
// @Tags Overview
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param overview_message formData string false "訊息區是否開啟" Enums(open, close)
// @param overview_topic formData string false "主題區是否開啟" Enums(open, close)
// @param overview_question formData string false "提問區是否開啟" Enums(open, close)
// @param overview_danmu formData string false "一般彈幕是否開啟" Enums(open, close)
// @param overview_special_danmu formData string false "特殊彈幕是否開啟" Enums(open, close)
// @param overview_picture formData string false "圖片播放區是否開啟" Enums(open, close)
// @param overview_holdscreen formData string false "霸佔彈幕是否開啟" Enums(open, close)
// @param overview_general formData string false "一般簽到是否開啟" Enums(open, close)
// @param overview_threed formData string false "立體簽到是否開啟" Enums(open, close)
// @param overview_countdown formData string false "倒數簽到是否開啟" Enums(open, close)
// @param overview_signname formData string false "簽名牆是否開啟" Enums(open, close)
// @param overview_lottery formData string false "遊戲抽獎是否開啟" Enums(open, close)
// @param overview_redpack formData string false "搖搖紅包是否開啟" Enums(open, close)
// @param overview_ropepack formData string false "抓紅包是否開啟" Enums(open, close)
// @param overview_whack_mole formData string false "敲敲樂是否開啟" Enums(open, close)
// @param overview_draw_numbers formData string false "搖號抽獎是否開啟" Enums(open, close)
// @param overview_monopoly formData string false "鑑定師是否開啟" Enums(open, close)
// @param overview_tugofwar formData string false "拔河遊戲是否開啟" Enums(open, close)
// @param overview_bingo formData string false "賓果遊戲是否開啟" Enums(open, close)
// @param overview_3d_gacha_machine formData string false "扭蛋機遊戲是否開啟" Enums(open, close)
// @param overview_qa formData string false "快問快答是否開啟" Enums(open, close)
// @param overview_vote formData string false "投票遊戲是否開啟" Enums(open, close)
// @param overview_chatroom formData string false "overview_chatroom" Enums(open, close)
// @param overview_interact formData string false "overview_interact" Enums(open, close)
// @param overview_info formData string false "overview_info" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/overview [patch]
func PATCHOverview(ctx *gin.Context) {
}

// @Summary 編輯活動介紹資料
// @Tags Introduce
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param id formData integer true "ID"
// @param activity_id formData string true "活動ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce_type formData string false "介紹類型" Enums(text, picture)
// @param content formData string false "內容(上限為200個字元)" maxlength(200)
// @param introduce_order formData integer false "資料排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/introduce [put]
func PUTIntroduce(ctx *gin.Context) {
}

// @Summary 活動介紹設置
// @Tags Introduce
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param introduce_title formData string false "手機頁面標題(上限為20個字元)" maxlength(20)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/introduce [patch]
func PATCHIntroduce(ctx *gin.Context) {
}

// @Summary 編輯活動行程資料
// @Tags Schedule
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "ID"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string false "主題(上限為20個字元)" minlength(1) maxlength(20)
// @param content formData string false "內容(上限為200個字元)" maxlength(200)
// @param schedule_date formData string false "日期(西元年-月-日)"
// @param start_time formData string false "開始時間(時:分)"
// @param end_time formData string false "結束時間(時:分)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/schedule [put]
func PUTSchedule(ctx *gin.Context) {
}

// @Summary 活動行程設置
// @Tags Schedule
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param schedule_title formData string false "手機頁面標題(上限為20個字元)" maxlength(20)
// @param schedule_display_date formData string false "顯示日期" Enums(open, close)
// @param schedule_display_detail formData string false "顯示詳情" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/schedule [patch]
func PATCHSchedule(ctx *gin.Context) {
}

// @Summary 編輯活動嘉賓資料
// @Tags Guest
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "ID"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param name formData string false "姓名(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce formData string false "簡介(上限為20個字元)" maxlength(20)
// @param detail formData string false "詳情(上限為200個字元)" maxlength(200)
// @param avatar formData file false "照片"
// @param guest_order formData integer false "資料排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/guest [put]
func PUTGuest(ctx *gin.Context) {
}

// @Summary 活動嘉賓設置
// @Tags Guest
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param guest_title formData string false "手機頁面標題(上限為20個字元)" maxlength(20)
// @param guest_background formData file false "背景圖"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/guest [patch]
func PATCHGuest(ctx *gin.Context) {
}

// @Summary 編輯活動資料
// @Tags Material
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "ID"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string false "資料標題(上限為20個字元)" minlength(1) maxlength(20)
// @param introduce formData string false "資料說明(上限為200個字元)" maxlength(200)
// @param link formData string false "資料連結"
// @param material_order formData integer false "資料排序"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/material [put]
func PUTMaterial(ctx *gin.Context) {
}

// @Summary 活動資料設置
// @Tags Material
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param material_title formData string false "手機頁面標題(上限為20個字元)" maxlength(20)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/material [patch]
func PATCHMaterial(ctx *gin.Context) {
}

// @Summary 刪除活動介紹資料
// @Tags Introduce
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/introduce [delete]
func DELETEIntroduce(ctx *gin.Context) {
}

// @Summary 刪除活動行程資料
// @Tags Schedule
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/schedule [delete]
func DELETESchedule(ctx *gin.Context) {
}

// @Summary 刪除活動嘉賓資料
// @Tags Guest
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/guest [delete]
func DELETEGuest(ctx *gin.Context) {
}

// @Summary 刪除活動資料
// @Tags Material
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/material [delete]
func DELETEMaterial(ctx *gin.Context) {
}

// @Summary 活動介紹JSON資料
// @Tags Introduce
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/introduce [get]
func ActivityIntroduceJSON(ctx *gin.Context) {
}

// @Summary 活動行程JSON資料
// @Tags Schedule
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/schedule [get]
func ActivityScheduleJSON(ctx *gin.Context) {
}

// @Summary 活動嘉賓JSON資料
// @Tags Guest
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/guest [get]
func ActivityGuestJSON(ctx *gin.Context) {
}

// @Summary 活動資料JSON資料
// @Tags Material
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /info/material [get]
func ActivityMaterialJSON(ctx *gin.Context) {
}

// models.EditActivityOverviewModel{
// 	ActivityID:             values.Get("activity_id"),
// 	OverviewMessage:        values.Get("overview_message"),
// 	OverviewTopic:          values.Get("overview_topic"),
// 	OverviewQuestion:       values.Get("overview_question"),
// 	OverviewDanmu:          values.Get("overview_danmu"),
// 	OverviewSpecialDanmu:   values.Get("overview_special_danmu"),
// 	OverviewPicture:        values.Get("overview_picture"),
// 	OverviewHoldscreen:     values.Get("overview_holdscreen"),
// 	OverviewGeneral:        values.Get("overview_general"),
// 	OverviewThreed:         values.Get("overview_threed"),
// 	OverviewCountdown:      values.Get("overview_countdown"),
// 	OverviewLottery:        values.Get("overview_lottery"),
// 	OverviewRedpack:        values.Get("overview_redpack"),
// 	OverviewRopepack:       values.Get("overview_ropepack"),
// 	OverviewWhackMole:      values.Get("overview_whack_mole"),
// 	OverviewDrawNumbers:    values.Get("overview_draw_numbers"),
// 	OverviewMonopoly:       values.Get("overview_monopoly"),
// 	OverviewQA:             values.Get("overview_qa"),
// 	OverviewTugofwar:       values.Get("overview_tugofwar"),
// 	OverviewBingo:          values.Get("overview_bingo"),
// 	OverviewSignname:       values.Get("overview_signname"),
// 	Overview3DGachaMachine: values.Get("overview_3d_gacha_machine"),
// 	OverviewVote:           values.Get("overview_vote"),
// 	OverviewChatroom:       values.Get("overview_chatroom"),
// 	OverviewInteract:       values.Get("overview_interact"),
// 	OverviewInfo:           values.Get("overview_info"),
// }

// models.EditIntroduceModel{
// 	ActivityID:     values.Get("activity_id"),
// 	Title:          values.Get("title"),
// 	IntroduceType:  values.Get("introduce_type"),
// 	Content:        values.Get("content"),
// 	IntroduceOrder: values.Get("introduce_order"),
// }

// models.EditIntroduceModel{
// 	ID:             values.Get("id"),
// 	ActivityID:     values.Get("activity_id"),
// 	Title:          values.Get("title"),
// 	IntroduceType:  values.Get("introduce_type"),
// 	Content:        values.Get("content"),
// 	IntroduceOrder: values.Get("introduce_order"),
// }

// models.EditIntroduceSettingModel{
// 	ActivityID:     values.Get("activity_id"),
// 	IntroduceTitle: values.Get("introduce_title"),
// }

// models.EditScheduleModel{
// 	ActivityID:   values.Get("activity_id"),
// 	Title:        values.Get("title"),
// 	Content:      values.Get("content"),
// 	ScheduleDate: values.Get("schedule_date"),
// 	StartTime:    values.Get("start_time"),
// 	EndTime:      values.Get("end_time"),
// }

// models.EditScheduleModel{
// 	ID:           values.Get("id"),
// 	ActivityID:   values.Get("activity_id"),
// 	Title:        values.Get("title"),
// 	Content:      values.Get("content"),
// 	ScheduleDate: values.Get("schedule_date"),
// 	StartTime:    values.Get("start_time"),
// 	EndTime:      values.Get("end_time"),
// }

// models.EditScheduleSettingModel{
// 	ActivityID:            values.Get("activity_id"),
// 	ScheduleTitle:         values.Get("schedule_title"),
// 	ScheduleDisplayDate:   values.Get("schedule_display_date"),
// 	ScheduleDisplayDetail: values.Get("schedule_display_detail"),
// }

// models.EditGuestModel{
// 	ActivityID: values.Get("activity_id"),
// 	Name:       values.Get("name"),
// 	Avatar:     avatar,
// 	Introduce:  values.Get("introduce"),
// 	Detail:     values.Get("detail"),
// 	GuestOrder: values.Get("guest_order"),
// }

// models.EditGuestModel{
// 	ID:         values.Get("id"),
// 	ActivityID: values.Get("activity_id"),
// 	Name:       values.Get("name"),
// 	Avatar:     avatar,
// 	Introduce:  values.Get("introduce"),
// 	Detail:     values.Get("detail"),
// 	GuestOrder: values.Get("guest_order"),
// }

// models.EditGuestSettingModel{
// 	ActivityID:      values.Get("activity_id"),
// 	GuestTitle:      values.Get("guest_title"),
// 	GuestBackground:

// models.EditMaterialModel{
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	Introduce:     values.Get("introduce"),
// 	Link:          values.Get("link"),
// 	MaterialOrder: values.Get("material_order"),
// }

// models.EditMaterialModel{
// 	ID:            values.Get("id"),
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	Introduce:     values.Get("introduce"),
// 	Link:          values.Get("link"),
// 	MaterialOrder: values.Get("material_order"),
// }

// models.EditMaterialSettingModel{
// 	ActivityID:    values.Get("activity_id"),
// 	MaterialTitle: values.Get("material_title"),
// }
