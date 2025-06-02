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

var (
	signnamePictures = []PictureField{
		{FieldName: "signname_classic_h_pic_01", Path: "signname/classic/signname_classic_h_pic_01.png"},
		{FieldName: "signname_classic_h_pic_02", Path: "signname/classic/signname_classic_h_pic_02.png"},
		{FieldName: "signname_classic_h_pic_03", Path: "signname/classic/signname_classic_h_pic_03.png"},
		{FieldName: "signname_classic_h_pic_04", Path: "signname/classic/signname_classic_h_pic_04.png"},
		{FieldName: "signname_classic_h_pic_05", Path: "signname/classic/signname_classic_h_pic_05.png"},
		{FieldName: "signname_classic_h_pic_06", Path: "signname/classic/signname_classic_h_pic_06.jpg"},
		{FieldName: "signname_classic_h_pic_07", Path: "signname/classic/signname_classic_h_pic_07.png"},
		{FieldName: "signname_classic_g_pic_01", Path: "signname/classic/signname_classic_g_pic_01.jpg"},
		{FieldName: "signname_classic_c_pic_01", Path: "signname/classic/signname_classic_c_pic_01.jpg"},

		{FieldName: "signname_bgm", Path: "signname/%s/bgm/bgm.mp3"},
	}

	generalPictures = []PictureField{
		{FieldName: "general_classic_h_pic_01", Path: "general/classic/general_classic_h_pic_01.png"},
		{FieldName: "general_classic_h_pic_02", Path: "general/classic/general_classic_h_pic_02.png"},
		{FieldName: "general_classic_h_pic_03", Path: "general/classic/general_classic_h_pic_03.jpg"},
		{FieldName: "general_classic_h_ani_01", Path: "general/classic/general_classic_h_ani_01.png"},

		{FieldName: "general_bgm", Path: "general/%s/bgm/bgm.mp3"},
	}

	threedPictures = []PictureField{
		{FieldName: "threed_bgm", Path: "threed/%s/bgm/bgm.mp3"},
	}
)

// GetGeneralPanel 簽名牆
func (s *SystemTable) GetSignnamePanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	info := table.GetInfo()
	// formList := table.GetForm()
	settingFormList := table.GetSettingForm()

	info.SetTable(config.ACTIVITY_SIGNNAME_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_SIGNNAME_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}

			// 簽名牆
			// s.redisConn.DelCache(config.SIGNNAME_DATAS_REDIS + activityID)         // 簽名牆資料
			// s.redisConn.DelCache(config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID) // 簽名牆資料
			return nil
		})

	settingFormList.SetTable(config.ACTIVITY_2_TABLE).SetUpdateFunc(func(values form2.Values) error {
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
		var model models.EditSignnameSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(signnamePictures, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateSignnameSetting(true, model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetGeneralPanel 一般簽到
func (s *SystemTable) GetGeneralPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		var background string
		if values.Get("general_background"+DEFAULT_FALG) == "1" {
			background = UPLOAD_SYSTEM_URL + "img-sign-bg.png"
		} else if values.Get("general_background") != "" {
			background = values.Get("general_background")
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
		var model models.EditGeneralModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(generalPictures, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.GeneralBackground = background

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateGeneral(true, model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetThreedPanel 3D簽到牆
func (s *SystemTable) GetThreedPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	settingFormList := table.GetSettingForm()
	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		var (
			threedPics   = []string{"img-sign3d-headpic.png", "background/orange_map.jpg", "mask/mask-hilives.svg"}
			threedFields = []string{"threed_avatar", "threed_background", "threed_image_logo_picture"}
			threedUpdate = make([]string, 3)
		)
		for i, field := range threedFields {
			if values.Get(field+DEFAULT_FALG) == "1" {
				threedUpdate[i] = UPLOAD_SYSTEM_URL + threedPics[i]
			} else if values.Get(field) != "" {
				threedUpdate[i] = values.Get(field)
			} else {
				threedUpdate[i] = ""
			}
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditThreeDModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(threedPictures, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.ThreedAvatar = threedUpdate[0]
		model.ThreedBackground = threedUpdate[1]
		model.ThreedImageLogoPicture = threedUpdate[2]

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateThreed(true, model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetCountdownPanel 倒數簽到
func (s *SystemTable) GetCountdownPanel() (table Table) {

	return
}

// @Summary 一般簽到設置
// @Tags Sign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param general_display_people formData string false "顯示人數" Enums(open, close)
// @param general_style formData integer false "展示方式" Enums(1)
// @param general_background formData file false "背景"
// @param general_loop formData string false "是否輪播" Enums(open, close)
// @param general_latest formData string false "跳轉至最新簽名" Enums(open, close)
// @param general_topic formData string false "主題" Enums(01_classic)
// @param general_music formData string false "音樂" Enums(classic)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/general [patch]
func PATCHGeneral(ctx *gin.Context) {
}

// @Summary 立體簽到設置
// @Tags Sign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param threed_avatar formData file false "簽到頭像"
// @param threed_avatar_shape formData string false "頭像形狀" Enums(circle, square)
// @param threed_display_people formData string false "顯示人數" Enums(open, close)
// @param threed_background_style formData integer false "背景" Enums(1, 2, 3)
// @param threed_background formData file false "背景圖片"
// @param threed_image_logo formData string false "3D圖像(LOGO)" Enums(open, close)
// @param threed_image_logo_picture formData file false "3D圖像LOGO圖片"
// @param threed_image_circle formData string false "3D圖像(圓形)" Enums(open, close)
// @param threed_image_spiral formData string false "3D圖像(螺旋)" Enums(open, close)
// @param threed_image_rectangle formData string false "3D圖像(矩形)" Enums(open, close)
// @param threed_image_square formData string false "3D圖像(正方形)" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/threed [patch]
func PATCHThreed(ctx *gin.Context) {
}

// @Summary 倒數簽到設置
// @Tags Sign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param countdown_second formData integer false "倒數秒數"
// @param countdown_url formData string false "倒數後進入頁面" Enums(current, threed)
// @param countdown_avatar formData file false "頭像"
// @param countdown_avatar_shape formData string false "頭像形狀" Enums(circle, square)
// @param countdown_background formData integer false "背景" Enums(1, 2, 3)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/countdown [patch]
func PATCHCountdown(ctx *gin.Context) {
}

// @Summary 簽名牆設置
// @Tags Sign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param signname_mode formData string false "簽名模式" Enums(terminal, phone)
// @param signname_times formData integer false "簽名次數"
// @param signname_display formData string false "簽名牆顯示方式" Enums(dynamics, static)
// @param signname_limit_times formData string false "限制簽名次數" Enums(open, close)
// @param signname_loop formData string false "是否輪播" Enums(open, close)
// @param signname_latest formData string false "跳轉至最新簽名" Enums(open, close)
// @param signname_show_name formData string false "signname_show_name" Enums(open, close)
// @param signname_topic formData string false "主題" Enums(01_classic)
// @param signname_music formData string false "音樂" Enums(classic)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/signname [patch]
func PATCHSignname(ctx *gin.Context) {
}

// @Summary 刪除簽名牆資料
// @Tags Sign
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/signname [delete]
func DELETESignname(ctx *gin.Context) {
}

// @Summary 所有遊戲JSON資料(sign分類下)
// @Tags Sign
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game query string false "遊戲種類(用,間隔)"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign [get]
func SignGamesJSON(ctx *gin.Context) {
}

// table = DefaultTable(DefaultConfig())
// settingFormList := table.GetSettingForm()
// settingFormList.AddField("倒數秒數", "countdown_second", db.Int, form.Number)
// settingFormList.AddField("倒數後進入頁面", "countdown_url", db.Varchar, form.Radio).
// 	SetFieldOptions(types.FieldOptions{
// 		{Value: "current", Text: "當前頁面"}, {Value: "threed", Text: "3D簽到牆"},
// 	})
// settingFormList.AddField("頭像", "countdown_avatar", db.Varchar, form.File)
// settingFormList.AddField("頭像形狀", "countdown_avatar_shape", db.Varchar, form.Radio).
// 	SetFieldOptions(types.FieldOptions{
// 		{Value: "circle", Text: "圓形"}, {Value: "square", Text: "方形"},
// 	})
// settingFormList.AddField("背景", "countdown_background", db.Int, form.Radio).
// 	SetFieldOptions(types.FieldOptions{
// 		{Value: "1", Text: "預設1", Text2: UPLOAD_RADIO_URL + "countdown-background-style1.png"},
// 		{Value: "2", Text: "預設2", Text2: UPLOAD_RADIO_URL + "countdown-background-style2.png"},
// 		{Value: "3", Text: "預設3", Text2: UPLOAD_RADIO_URL + "countdown-background-style3.png"},
// 	})

// settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id") {
// 		return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
// 	}
// 	var avatar string
// 	//  圖片儲存需要處理
// 	if values.Get("countdown_avatar"+DEFAULT_FALG) == "1" {
// 		avatar = UPLOAD_SYSTEM_URL + "img-countdown-headpic.png"
// 	} else if values.Get("countdown_avatar") != "" {
// 		avatar = UPLOAD_URL + values.Get("countdown_avatar")
// 	} else {
// 		avatar = ""
// 	}

// 	if err := models.DefaultActivityModel().SetDbConn(s.dbConn).
// 		UpdateCountdown(models.EditCountdownModel{
// 			ActivityID:           values.Get("activity_id"),
// 			CountdownSecond:      values.Get("countdown_second"),
// 			CountdownURL:         values.Get("countdown_url"),
// 			CountdownAvatar:      avatar,
// 			CountdownAvatarShape: values.Get("countdown_avatar_shape"),
// 			CountdownBackground:  values.Get("countdown_background"),
// 		}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// 	return nil
// })

// formList.SetTable(config.ACTIVITY_SIGNNAME_TABLE).SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "user_id") {
// 		return errors.New("錯誤: 活動ID、用戶ID發生問題，請輸入有效的資料")
// 	}

// 	if err := models.DefaultSignnameModel().SetDbConn(s.dbConn).
// 		SetRedisConn(s.redisConn).Add(true,
// 		models.NewSignnameModel{
// 			ActivityID: values.Get("activity_id"),
// 			UserID:     values.Get("user_id"),
// 			Picture:    values.Get("picture"),
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// })

// var (
// 	pics = []string{
// 		// 簽名牆自定義
// 		"signname/classic/signname_classic_h_pic_01.png",
// 		"signname/classic/signname_classic_h_pic_02.png",
// 		"signname/classic/signname_classic_h_pic_03.png",
// 		"signname/classic/signname_classic_h_pic_04.png",
// 		"signname/classic/signname_classic_h_pic_05.png",
// 		"signname/classic/signname_classic_h_pic_06.jpg",
// 		"signname/classic/signname_classic_h_pic_07.png",
// 		"signname/classic/signname_classic_g_pic_01.jpg",
// 		"signname/classic/signname_classic_c_pic_01.jpg",

// 		"signname/%s/bgm/bgm.mp3",
// 	}

// 	fields = []string{
// 		// 簽名牆自定義
// 		"signname_classic_h_pic_01",
// 		"signname_classic_h_pic_02",
// 		"signname_classic_h_pic_03",
// 		"signname_classic_h_pic_04",
// 		"signname_classic_h_pic_05",
// 		"signname_classic_h_pic_06",
// 		"signname_classic_h_pic_07",
// 		"signname_classic_g_pic_01",
// 		"signname_classic_c_pic_01",

// 		"signname_bgm",
// 	}

// 	update = make([]string, 200)
// )

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("signname_topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditSignnameSettingModel{
// 	ActivityID:         values.Get("activity_id"),
// 	SignnameMode:       values.Get("signname_mode"),
// 	SignnameTimes:      values.Get("signname_times"),
// 	SignnameDisplay:    values.Get("signname_display"),
// 	SignnameLimitTimes: values.Get("signname_limit_times"),
// 	SignnameTopic:      values.Get("signname_topic"),
// 	SignnameMusic:      values.Get("signname_music"),
// 	SignnameLoop:       values.Get("signname_loop"),
// 	SignnameLatest:     values.Get("signname_latest"),
// 	SignnameContent:    values.Get("signname_content"),
// 	SignnameShowName:   values.Get("signname_show_name"),

// 	SignnameClassicHPic01: update[0],
// 	SignnameClassicHPic02: update[1],
// 	SignnameClassicHPic03: update[2],
// 	SignnameClassicHPic04: update[3],
// 	SignnameClassicHPic05: update[4],
// 	SignnameClassicHPic06: update[5],
// 	SignnameClassicHPic07: update[6],
// 	SignnameClassicGPic01: update[7],
// 	SignnameClassicCPic01: update[8],

// 	SignnameBgm: update[9],
// }

// var (
// 	pics = []string{
// 		// 一般簽到牆自定義
// 		"general/classic/general_classic_h_pic_01.png",
// 		"general/classic/general_classic_h_pic_02.png",
// 		"general/classic/general_classic_h_pic_03.jpg",
// 		"general/classic/general_classic_h_ani_01.png",

// 		"general/%s/bgm/bgm.mp3",
// 	}

// 	fields = []string{
// 		// 一般簽到牆自定義
// 		"general_classic_h_pic_01",
// 		"general_classic_h_pic_02",
// 		"general_classic_h_pic_03",
// 		"general_classic_h_ani_01",

// 		"general_bgm",
// 	}

// 	update = make([]string, 200)
// )

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("general_topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditGeneralModel{
// 	ActivityID:           values.Get("activity_id"),
// 	GeneralDisplayPeople: values.Get("general_display_people"),
// 	GeneralStyle:         values.Get("general_style"),
// 	GeneralBackground:    background,
// 	GeneralTopic:         values.Get("general_topic"),
// 	GeneralMusic:         values.Get("general_music"),
// 	GeneralLoop:          values.Get("general_loop"),
// 	GeneralLatest:        values.Get("general_latest"),

// 	GeneralClassicHPic01: update[0],
// 	GeneralClassicHPic02: update[1],
// 	GeneralClassicHPic03: update[2],
// 	GeneralClassicHAni01: update[3],

// 	GeneralBgm: update[4],
// }

// var (
// 	pics = []string{
// 		// 立體簽到牆自定義
// 		"threed/%s/bgm/bgm.mp3",
// 	}

// 	fields = []string{
// 		// 立體簽到牆自定義
// 		"threed_bgm",
// 	}

// 	update = make([]string, 200)
// )

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("threed_topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditThreeDModel{
// 	ActivityID:             values.Get("activity_id"),
// 	ThreedAvatar:           threedUpdate[0],
// 	ThreedAvatarShape:      values.Get("threed_avatar_shape"),
// 	ThreedDisplayPeople:    values.Get("threed_display_people"),
// 	ThreedBackgroundStyle:  values.Get("threed_background_style"),
// 	ThreedBackground:       threedUpdate[1],
// 	ThreedImageLogo:        values.Get("threed_image_logo"),
// 	ThreedImageCircle:      values.Get("threed_image_circle"),
// 	ThreedImageSpiral:      values.Get("threed_image_spiral"),
// 	ThreedImageRectangle:   values.Get("threed_image_rectangle"),
// 	ThreedImageSquare:      values.Get("threed_image_square"),
// 	ThreedImageLogoPicture: threedUpdate[2],
// 	ThreedTopic:            values.Get("threed_topic"),
// 	ThreedMusic:            values.Get("threed_music"),

// 	ThreedBgm: update[0],
// }
