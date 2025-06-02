package models

import (
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"log"
	"maps"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Add 增加資料
// #####新增用戶時也會自動新增活動場次，修改參數時須注意#####
func (a ActivityModel) Add(isRedis bool, model EditActivityModel) error {
	var (
		activityTypes = []string{"企業會議", "其他", "商業活動", "培訓/教育", "婚禮", "尾牙春酒",
			"校園活動", "競技賽事", "論壇會議", "酒吧/餐飲娛樂", "電視/媒體"}
		isbool bool

		// activity資料表欄位
		activityFields1 = []string{
			"activity_id",
			"user_id",
			"activity_name",
			"activity_type",
			"max_people",
			"people",
			"attend",
			"city",
			"town",
			"start_time",
			"end_time",
			"number",
			"login_required",
			"login_password",
			"push_message",
			"password_required",
			"device",
			"push_phone_message",
			"message_amount",
			"send_mail",
			"mail_amount",
			"host_scan",
			"customize_password",
			"allow_customize_apply",

			// 官方帳號
			"line_id",
			"channel_id",
			"channel_secret",
			"chatbot_secret",
			"chatbot_token",

			// 活動總覽
			"overview_message",
			"overview_topic",
			"overview_question",
			"overview_danmu",
			"overview_special_danmu",
			"overview_picture",
			"overview_holdscreen",
			"overview_general",
			"overview_threed",
			"overview_countdown",
			"overview_lottery",
			"overview_redpack",
			"overview_ropepack",
			"overview_whack_mole",
			"overview_draw_numbers",
			"overview_monopoly",
			"overview_qa",
			"overview_tugofwar",
			"overview_bingo",
			"overview_3d_gacha_machine",
			"overview_signname",

			// 活動介紹
			"introduce_title",

			// 活動行程
			"schedule_title",
			"schedule_display_date",
			"schedule_display_detail",

			// 活動嘉賓
			"guest_title",
			"guest_background",

			// 活動資料
			"material_title",

			// 報名
			"apply_check",

			// 簽到
			"sign_check",
			"sign_allow",
			"sign_minutes",
			"sign_manual",

			// QRcode
			"qrcode_logo_picture",
			"qrcode_logo_size",
			"qrcode_picture_point",
			"qrcode_white_distance",
			"qrcode_point_color",
			"qrcode_background_color",

			// 訊息牆
			"message_picture",
			"message_auto",
			"message_ban",
			"message_ban_second",
			"message_refresh_second",
			"message_open",
			"message",
			"message_background",
			"message_text_color",
			"message_screen_title_color",
			"message_text_frame_color",
			"message_frame_color",
			"message_screen_title_frame_color",

			// 主題牆
			"topic_background",

			// 提問牆
			"question_anonymous",
			"question_qrcode",
			"question_background",

			// 彈幕
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
			"danmu_background",

			// 特殊彈幕
			"special_danmu_message_check",
			"special_danmu_general_price",
			"special_danmu_large_price",
			"special_danmu_over_price",
			"special_danmu_topic",

			// 圖片牆
			"picture_start_time",
			"picture_end_time",
			"picture_hide_time",
			"picture_switch_second",
			"picture_play_order",
			"picture",
			"picture_background",
			"picture_animation",

			// 霸屏
			"holdscreen_price",
			"holdscreen_message_check",
			"holdscreen_only_picture",
			"holdscreen_display_second",
			"holdscreen_birthday_topic",

			// 一般簽到
			"general_display_people",
			"general_style",
			"general_background",

			// 3D簽到
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
			"threed_image_logo_picture",

			// 倒數計時
			"countdown_second",
			"countdown_url",
			"countdown_avatar",
			"countdown_avatar_shape",
			"countdown_background",

			// 訊息審核
			"message_check_manual_check",
			"message_check_sensitivity",
		}

		// activity_2資料表欄位
		activityFields2 = []string{
			"activity_id",
			"user_id",

			// 一般簽到
			"general_topic",
			"general_music",
			"general_loop",
			"general_latest",
			"general_edit_times",
			"general_classic_h_pic_01",
			"general_classic_h_pic_02",
			"general_classic_h_pic_03",
			"general_classic_h_ani_01",
			"general_bgm",

			// 立體簽到
			"threed_topic",
			"threed_music",
			"threed_edit_times",
			"threed_bgm",

			// 簽名牆
			"signname_mode",
			"signname_times",
			"signname_display",
			"signname_limit_times",
			"signname_topic",
			"signname_music",
			"signname_loop",
			"signname_latest",
			"signname_edit_times",
			"signname_content",
			"signname_show_name",
			"signname_classic_h_pic_01",
			"signname_classic_h_pic_02",
			"signname_classic_h_pic_03",
			"signname_classic_h_pic_04",
			"signname_classic_h_pic_05",
			"signname_classic_h_pic_06",
			"signname_classic_h_pic_07",
			"signname_classic_g_pic_01",
			"signname_classic_c_pic_01",
			"signname_bgm",

			// 總覽功能
			"overview_vote",
			"overview_chatroom",
			"overview_interact",
			"overview_info",

			// 一般彈幕
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

			// 特殊彈幕
			"special_danmu_ban_second",

			// 霸屏
			"holdscreen_ban_second",

			// 報名
			"customize_default_avatar",
		}

		// activity_channel資料表欄位
		activityChannelFields = []string{
			"id",
			"activity_id",
			"user_id",
			"channel_1",
			"channel_2",
			"channel_3",
			"channel_4",
			"channel_5",
			"channel_6",
			"channel_7",
			"channel_8",
			"channel_9",
			"channel_10",

			"channel_amount",
		}
	)

	// 欄位值判斷
	if utf8.RuneCountInString(model.ActivityName) > 20 {
		return errors.New("錯誤: 活動名稱上限為20個字元，請輸入有效的活動名稱")
	}

	// 判斷活動人數上限
	maxPeopleInt, err1 := strconv.Atoi(model.MaxPeople)
	peopleInt, err2 := strconv.Atoi(model.People)
	if err1 != nil || err2 != nil || peopleInt > maxPeopleInt {
		return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
	}

	for i := range activityTypes {
		if model.ActivityType == activityTypes[i] {
			isbool = true
			break
		}
	}
	if isbool == false {
		return errors.New("錯誤: 活動類型發生問題，請輸入有效的活動類型")
	}

	if !CompareTDatetime(model.StartTime, model.EndTime) {
		return errors.New("錯誤: 活動時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
	}

	if model.LoginRequired != "open" && model.LoginRequired != "close" {
		return errors.New("錯誤: 是否需要登入才可進入聊天室資料發生問題，請輸入有效的資料")
	}

	if utf8.RuneCountInString(model.LoginPassword) > 8 {
		return errors.New("錯誤: 聊天室驗證碼上限為8個字元，請輸入有效的密碼")
	}
	if model.PasswordRequired != "open" && model.PasswordRequired != "close" {
		return errors.New("錯誤: 是否需要設置密碼才可進入聊天室資料發生問題，請輸入有效的資料")
	}

	if model.PushMessage != "open" && model.PushMessage != "close" {
		return errors.New("錯誤: 官方帳號發送訊息資料發生問題，請輸入有效的資料")
	}

	if model.PushPhoneMessage == "" {
		// 預設關閉
		model.PushPhoneMessage = "close"
	}

	if model.MessageAmount == "" {
		// 簡訊數量預設
		model.MessageAmount = "0"
	}

	if model.SendMail == "" {
		// 預設關閉
		model.SendMail = "close"
	}
	if model.MailAmount == "" {
		// 郵件數量預設
		model.MailAmount = "0"
	}

	if model.HostScan == "" {
		// 預設關閉
		model.HostScan = "close"
	}

	if model.Device == "" {
		// 預設line
		model.Device = "line"
	}

	if model.ChannelAmount == "" {
		// 頻道數量
		model.ChannelAmount = "5"
	} else {
		if _, err := strconv.Atoi(model.ChannelAmount); err != nil {
			return errors.New("錯誤: 頻道數量資料發生問題，請輸入有效的頻道數量")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 處理欄位預設值
	dataDefaults := command.Value{
		"attend":                0,
		"number":                1,
		"customize_password":    "close",
		"allow_customize_apply": "close",

		// 活動總覽
		"overview_message":          "open",
		"overview_topic":            "close",
		"overview_question":         "close",
		"overview_danmu":            "open",
		"overview_special_danmu":    "close",
		"overview_picture":          "close",
		"overview_holdscreen":       "close",
		"overview_general":          "close",
		"overview_threed":           "close",
		"overview_countdown":        "close",
		"overview_lottery":          "close",
		"overview_redpack":          "close",
		"overview_ropepack":         "close",
		"overview_whack_mole":       "close",
		"overview_draw_numbers":     "close",
		"overview_monopoly":         "close",
		"overview_qa":               "close",
		"overview_tugofwar":         "close",
		"overview_bingo":            "close",
		"overview_3d_gacha_machine": "close",
		"overview_signname":         "close",

		// 活動介紹
		"introduce_title": "活動介紹",

		// 活動行程
		"schedule_title":          "活動行程",
		"schedule_display_date":   "open",
		"schedule_display_detail": "open",

		// 活動嘉賓
		"guest_title":      "活動嘉賓",
		"guest_background": config.UPLOAD_SYSTEM_URL + "img-guest-bg.png",

		// 活動資料
		"material_title": "活動資料",

		// 報名
		"apply_check": "close",

		// 簽到
		"sign_check":   "close",
		"sign_allow":   "close",
		"sign_minutes": 30,
		"sign_manual":  "open",

		// QRCode
		"qrcode_logo_picture":     config.UPLOAD_SYSTEM_URL + "qr-logo.svg",
		"qrcode_logo_size":        0.5,
		"qrcode_picture_point":    "close",
		"qrcode_white_distance":   0,
		"qrcode_point_color":      "#424242",
		"qrcode_background_color": "#FAFAFA",

		// 訊息牆
		"message_picture":                  "close",
		"message_auto":                     "close",
		"message_ban":                      "close",
		"message_ban_second":               1,
		"message_refresh_second":           20,
		"message_open":                     "open",
		"message":                          "歡迎來到活動聊天室!",
		"message_background":               config.UPLOAD_SYSTEM_URL + "img-topic-bg.jpg",
		"message_text_color":               "#404b79",
		"message_screen_title_color":       "#ffffff",
		"message_text_frame_color":         "#81869b",
		"message_frame_color":              "#e7eaf4",
		"message_screen_title_frame_color": "#81869b",

		// 主題牆
		"topic_background": config.UPLOAD_SYSTEM_URL + "img-topic-bg.jpg",

		// 提問牆
		"question_anonymous":  "close",
		"question_qrcode":     "open",
		"question_background": config.UPLOAD_SYSTEM_URL + "img-request-bg.png",

		// 彈幕
		"danmu_loop":           "open",
		"danmu_top":            "open",
		"danmu_mid":            "open",
		"danmu_bottom":         "open",
		"danmu_display_name":   "open",
		"danmu_display_avatar": "open",
		"danmu_size":           0.75,
		"danmu_speed":          3,
		"danmu_density":        2.5,
		"danmu_opacity":        0.9,
		"danmu_background":     1,

		// 特殊彈幕
		"special_danmu_message_check": "close",
		"special_danmu_general_price": 0,
		"special_danmu_large_price":   0,
		"special_danmu_over_price":    0,
		"special_danmu_topic":         "open",

		// 圖片牆
		"picture_start_time":    model.StartTime,
		"picture_end_time":      model.EndTime,
		"picture_hide_time":     "close",
		"picture_switch_second": 3,
		"picture_play_order":    "order",
		"picture":               "system/img-pic.png",
		"picture_background":    config.UPLOAD_SYSTEM_URL + "img-picture-bg.png",
		"picture_animation":     1,

		// 霸屏
		"holdscreen_price":          0,
		"holdscreen_message_check":  "close",
		"holdscreen_only_picture":   "close",
		"holdscreen_display_second": 10,
		"holdscreen_birthday_topic": "open",

		// 一般簽到
		"general_display_people": "open",
		"general_style":          1,
		"general_background":     config.UPLOAD_SYSTEM_URL + "img-sign-bg.png",

		// 3D 簽到
		"threed_avatar":             config.UPLOAD_SYSTEM_URL + "img-sign3d-headpic.png",
		"threed_avatar_shape":       "circle",
		"threed_display_people":     "open",
		"threed_background_style":   1,
		"threed_background":         config.UPLOAD_SYSTEM_URL + "background/orange_map.jpg",
		"threed_image_logo":         "open",
		"threed_image_circle":       "open",
		"threed_image_spiral":       "open",
		"threed_image_rectangle":    "open",
		"threed_image_square":       "open",
		"threed_image_logo_picture": config.UPLOAD_SYSTEM_URL + "mask/mask-hilives.svg",

		// 倒數計時
		"countdown_second":       5,
		"countdown_url":          "current",
		"countdown_avatar":       config.UPLOAD_SYSTEM_URL + "img-countdown-headpic.png",
		"countdown_avatar_shape": "circle",
		"countdown_background":   1,

		// 訊息審核
		"message_check_manual_check": "close",
		"message_check_sensitivity":  "refuse",

		// 一般簽到
		"general_topic":            "01_classic",
		"general_music":            "classic",
		"general_loop":             "open",
		"general_latest":           "open",
		"general_edit_times":       0,
		"general_classic_h_pic_01": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_01.png",
		"general_classic_h_pic_02": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_02.png",
		"general_classic_h_pic_03": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_03.jpg",
		"general_classic_h_ani_01": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_ani_01.png",
		"general_bgm":              config.UPLOAD_SYSTEM_URL + "general/classic/bgm/bgm.mp3",

		// 立體簽到
		"threed_topic":      "01_classic",
		"threed_music":      "classic",
		"threed_edit_times": 0,
		"threed_bgm":        config.UPLOAD_SYSTEM_URL + "threed/classic/bgm/bgm.mp3",

		// 簽名牆
		"signname_mode":             "phone",
		"signname_times":            1,
		"signname_display":          "dynamics",
		"signname_limit_times":      "open",
		"signname_topic":            "01_classic",
		"signname_music":            "classic",
		"signname_loop":             "open",
		"signname_latest":           "open",
		"signname_edit_times":       0,
		"signname_content":          "write",
		"signname_show_name":        "close",
		"signname_classic_h_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_01.png",
		"signname_classic_h_pic_02": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_02.png",
		"signname_classic_h_pic_03": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_03.png",
		"signname_classic_h_pic_04": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_04.png",
		"signname_classic_h_pic_05": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_05.png",
		"signname_classic_h_pic_06": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_06.jpg",
		"signname_classic_h_pic_07": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_07.png",
		"signname_classic_g_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_g_pic_01.jpg",
		"signname_classic_c_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_c_pic_01.jpg",
		"signname_bgm":              config.UPLOAD_SYSTEM_URL + "signname/classic/bgm/bgm.mp3",

		// 活動總覽其他項目
		"overview_vote":     "close",
		"overview_chatroom": "open",
		"overview_interact": "open",
		"overview_info":     "open",

		// 一般彈幕樣式
		"danmu_new_background_color": "#90ace7",
		"danmu_new_text_color":       "#1e359d",
		"danmu_old_background_color": "#fdc1f8",
		"danmu_old_text_color":       "#9d1e60",
		"danmu_style":                "classic",
		"danmu_custom_left_new":      "/admin/uploads/system/danmu/custom/left_new.svg",
		"danmu_custom_center_new":    "/admin/uploads/system/danmu/custom/center_new.svg",
		"danmu_custom_right_new":     "/admin/uploads/system/danmu/custom/right_new.svg",
		"danmu_custom_left_old":      "/admin/uploads/system/danmu/custom/left_old.svg",
		"danmu_custom_center_old":    "/admin/uploads/system/danmu/custom/center_old.svg",
		"danmu_custom_right_old":     "/admin/uploads/system/danmu/custom/right_old.svg",

		// 特殊彈幕
		"special_danmu_ban_second": 20,

		// 霸屏
		"holdscreen_ban_second": 20,

		// 報名
		"customize_default_avatar": fmt.Sprintf("%s/admin/uploads/system/img-user-pic.png", config.HTTP_HILIVES_NET_URL),

		// 頻道
		"channel_1":      "close",
		"channel_2":      "close",
		"channel_3":      "close",
		"channel_4":      "close",
		"channel_5":      "close",
		"channel_6":      "close",
		"channel_7":      "close",
		"channel_8":      "close",
		"channel_9":      "close",
		"channel_10":     "close",
		"channel_amount": utils.GetInt64(model.ChannelAmount, 0),
	}

	// log.Println("處理前: ", data)

	// 合併預設值進 data
	maps.Copy(data, dataDefaults)

	// 將id資料手動寫入data
	// 取得mongo中的id資料(遞增處理)
	mongoID, _ := a.MongoConn.GetNextSequence(config.ACTIVITY_CHANNEL_TABLE)
	data["id"] = mongoID

	// log.Println("處理後: ", data)

	if _, err := a.Table(a.TableName).Insert(FilterFields(data, activityFields1)); err != nil {
		return errors.New("錯誤: 無法新增活動資料(activity)，請重新操作")
	}

	if _, err := a.Table(config.ACTIVITY_2_TABLE).Insert(FilterFields(data, activityFields2)); err != nil {
		return errors.New("錯誤: 無法新增活動資料(activity_2)，請重新操作")
	}

	log.Println("新增活動api, 新增頻道數量(mongo): ", model.ChannelAmount)

	// log.Println("mongoID: ", mongoID)

	_, err := a.MongoConn.InsertOne(config.ACTIVITY_CHANNEL_TABLE, FilterFields(data, activityChannelFields))
	if err != nil {
		return errors.New("錯誤: 無法新增活動資料(activity_channel，mongo)，請重新操作")
	}

	// 建立自定義資料
	if err := DefaultCustomizeModel().
	SetConn(a.DbConn, a.RedisConn, a.MongoConn).
	Add(
		EditCustomizeModel{
			ActivityID: model.ActivityID,
		}); err != nil {
		return err
	}

	// 建立遊戲基本設置資料
	if err := DefaultGameSettingModel().
	SetConn(a.DbConn, a.RedisConn, a.MongoConn).
	Add(
		EditGameSettingModel{
			ActivityID: model.ActivityID,
		}); err != nil {
		return err
	}

	// 新增活動權限
	if model.Permissions != "" {
		a.ActivityID = model.ActivityID
		if err := a.AddPermission(strings.Split(model.Permissions, ",")); err != nil {
			return err
		}
	}

	// 準備要建立的資料夾相對路徑
	folders := []string{
		"", // 活動根資料夾
		"info/introduce",
		"info/guest",

		"applysign/apply",
		"applysign/customize",
		"applysign/qrcode",

		"interact/wall/message",
		"interact/wall/topic",
		"interact/wall/question",
		"interact/wall/danmu",

		"interact/sign/general",
		"interact/sign/threed",
		"interact/sign/signname",
		"interact/sign/vote",

		"interact/game/lottery",
		"interact/game/redpack",
		"interact/game/ropepack",
		"interact/game/draw_numbers",
		"interact/game/whack_mole",
		"interact/game/monopoly",
		"interact/game/QA",
		"interact/game/tugofwar",
		"interact/game/bingo",
		"interact/game/3DGachaMachine",
	}

	// 逐一建立資料夾
	for _, folder := range folders {
		path := fmt.Sprintf("%s/%s/%s/%s", config.STORE_PATH, model.UserID, model.ActivityID, folder)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Printf("建立資料夾失敗: %s, 錯誤: %v", path, err)
		}
	}

	// 新增活動時需要清除用戶redis資訊(重新判斷活動權限資訊)
	if isRedis {
		a.RedisConn.DelCache(config.HILIVE_USERS_REDIS + model.UserID)
	}
	return nil
}

// 設定時區為 UTC+8
// loc, _ := time.LoadLocation("Asia/Taipei")
// now := time.Now().In(loc)
// // 格式化為 "YYYY-MM-DD HH:MM:SS"
// formattedTime := now.Format("2006-01-02 15:04:05")
// "created_at":     formattedTime, // 設置 `created_at`
// "updated_at":     formattedTime, // 設置 `updated_at`

// log.Println("新增活動api, 新增頻道數量(mysql): ", model.ChannelAmount)

// if _, err := a.Table(config.ACTIVITY_CHANNEL_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"user_id":     model.UserID,

// 	"channel_amount": model.ChannelAmount,
// }); err != nil {
// 	return errors.New("錯誤: 無法新增活動資料(activity_channel,mysql)，請重新操作")
// }

// 用戶資訊
// userModel, err := DefaultUserModel().SetDbConn(a.DbConn).
// 	SetRedisConn(a.RedisConn).Find(true, true, "",
// 	"users.user_id", model.UserID)
// if err != nil {
// 	return err
// }
// peopleInt, err := strconv.Atoi(model.People)
// if err != nil {
// 	return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
// }
// if err != nil || peopleInt > int(userModel.MaxActivityPeople) {
// 	return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
// }

// command.Value{
// "activity_id":           model.ActivityID,
// "user_id":               model.UserID,
// "activity_name":         model.ActivityName,
// "activity_type":         model.ActivityType,
// "max_people":            maxPeopleInt,
// "people":                peopleInt,
// "attend":                0,
// "city":                  model.City,
// "town":                  model.Town,
// "start_time":            model.StartTime,
// "end_time":              model.EndTime,
// "number":                1,
// "login_required":        model.LoginRequired,
// "login_password":        model.LoginPassword,
// "push_message":          model.PushMessage,
// "password_required":     model.PasswordRequired,
// "device":                model.Device,
// "push_phone_message":    model.PushPhoneMessage,
// "message_amount":        model.MessageAmount,
// "send_mail":             model.SendMail,
// "mail_amount":           model.MailAmount,
// "host_scan":             model.HostScan,
// "customize_password":    "close",
// "allow_customize_apply": "close",

// !!!!!修改以上參數時，要注意新增用戶時(user.go Add function)的活動相關必填參數!!!!!

// 官方帳號
// "line_id":        model.LineID,
// "channel_id":     model.ChannelID,
// "channel_secret": model.ChannelSecret,
// "chatbot_secret": model.ChatbotSecret,
// "chatbot_token":  model.ChatbotToken,

// 活動總覽
// "overview_message":          "open",
// "overview_topic":            "close",
// "overview_question":         "close",
// "overview_danmu":            "open",
// "overview_special_danmu":    "close",
// "overview_picture":          "close",
// "overview_holdscreen":       "close",
// "overview_general":          "close",
// "overview_threed":           "close",
// "overview_countdown":        "close",
// "overview_lottery":          "close",
// "overview_redpack":          "close",
// "overview_ropepack":         "close",
// "overview_whack_mole":       "close",
// "overview_draw_numbers":     "close",
// "overview_monopoly":         "close",
// "overview_qa":               "close",
// "overview_tugofwar":         "close",
// "overview_bingo":            "close",
// "overview_3d_gacha_machine": "close",
// "overview_signname":         "close",

// 活動介紹
// "introduce_title": "活動介紹",
// 活動行程
// "schedule_title":          "活動行程",
// "schedule_display_date":   "open",
// "schedule_display_detail": "open",
// 活動嘉賓
// "guest_title":      "活動嘉賓",
// "guest_background": config.UPLOAD_SYSTEM_URL + "img-guest-bg.png",
// 活動資料
// "material_title": "活動資料",
// 報名
// "apply_check": "close",
// 簽到
// "sign_check":   "close",
// "sign_allow":   "close",
// "sign_minutes": 30,
// "sign_manual":  "open",
// QRcode
// "qrcode_logo_picture":     config.UPLOAD_SYSTEM_URL + "qr-logo.svg",
// "qrcode_logo_size":        0.5,
// "qrcode_picture_point":    "close",
// "qrcode_white_distance":   0,
// "qrcode_point_color":      "#424242",
// "qrcode_background_color": "#FAFAFA",
// 訊息牆
// "message_picture":                  "close",
// "message_auto":                     "close",
// "message_ban":                      "close",
// "message_ban_second":               1,
// "message_refresh_second":           20,
// "message_open":                     "open",
// "message":                          "歡迎來到活動聊天室!",
// "message_background":               config.UPLOAD_SYSTEM_URL + "img-topic-bg.jpg",
// "message_text_color":               "#404b79",
// "message_screen_title_color":       "#ffffff",
// "message_text_frame_color":         "#81869b",
// "message_frame_color":              "#e7eaf4",
// "message_screen_title_frame_color": "#81869b",
// 主題牆
// "topic_background": config.UPLOAD_SYSTEM_URL + "img-topic-bg.jpg",
// 提問牆
// "question_message_check":  "close",
// "question_anonymous": "close",
// "question_hide_answered":  "close",
// "question_qrcode":     "open",
// "question_background": config.UPLOAD_SYSTEM_URL + "img-request-bg.png",
// "question_guest_answered": "close",
// 彈幕
// "danmu_loop":           "open",
// "danmu_top":            "open",
// "danmu_mid":            "open",
// "danmu_bottom":         "open",
// "danmu_display_name":   "open",
// "danmu_display_avatar": "open",
// "danmu_size":           0.75,
// "danmu_speed":          3,
// "danmu_density":        2.5,
// "danmu_opacity":        0.9,
// "danmu_background":     1,

// 特殊彈幕
// "special_danmu_message_check": "close",
// "special_danmu_general_price": 0,
// "special_danmu_large_price":   0,
// "special_danmu_over_price":    0,
// "special_danmu_topic":         "open",
// 圖片牆
// "picture_start_time":    model.StartTime,
// "picture_end_time":      model.EndTime,
// "picture_hide_time":     "close",
// "picture_switch_second": 3,
// "picture_play_order":    "order",
// "picture":               "system/img-pic.png",
// "picture_background":    config.UPLOAD_SYSTEM_URL + "img-picture-bg.png",
// "picture_animation":     1,
// 霸屏
// "holdscreen_price":          0,
// "holdscreen_message_check":  "close",
// "holdscreen_only_picture":   "close",
// "holdscreen_display_second": 10,
// "holdscreen_birthday_topic": "open",
// 一般簽到
// "general_display_people": "open",
// "general_style":          1,
// "general_background":     config.UPLOAD_SYSTEM_URL + "img-sign-bg.png",

// 3D簽到
// "threed_avatar":             config.UPLOAD_SYSTEM_URL + "img-sign3d-headpic.png",
// "threed_avatar_shape":       "circle",
// "threed_display_people":     "open",
// "threed_background_style":   1,
// "threed_background":         config.UPLOAD_SYSTEM_URL + "background/orange_map.jpg",
// "threed_image_logo":         "open",
// "threed_image_circle":       "open",
// "threed_image_spiral":       "open",
// "threed_image_rectangle":    "open",
// "threed_image_square":       "open",
// "threed_image_logo_picture": config.UPLOAD_SYSTEM_URL + "mask/mask-hilives.svg",

// 倒數計時
// "countdown_second":       5,
// "countdown_url":          "current",
// "countdown_avatar":       config.UPLOAD_SYSTEM_URL + "img-countdown-headpic.png",
// "countdown_avatar_shape": "circle",
// "countdown_background":   1,

// 訊息審核
// "message_check_manual_check": "close",
// "message_check_sensitivity":  "refuse",
// }

// command.Value{
// "activity_id": model.ActivityID,
// "user_id":     model.UserID,

// 一般簽到
// "general_topic":      "01_classic",
// "general_music":      "classic",
// "general_loop":       "open",
// "general_latest":     "open",
// "general_edit_times": 0,

// "general_classic_h_pic_01": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_01.png",
// "general_classic_h_pic_02": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_02.png",
// "general_classic_h_pic_03": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_pic_03.jpg",
// "general_classic_h_ani_01": config.UPLOAD_SYSTEM_URL + "general/classic/general_classic_h_ani_01.png",

// "general_bgm": config.UPLOAD_SYSTEM_URL + "general/classic/bgm/bgm.mp3",

// 立體簽到
// "threed_topic":      "01_classic",
// "threed_music":      "classic",
// "threed_edit_times": 0,

// "threed_bgm": config.UPLOAD_SYSTEM_URL + "threed/classic/bgm/bgm.mp3",

// 簽名牆
// "signname_mode":        "phone",
// "signname_times":       1,
// "signname_display":     "dynamics",
// "signname_limit_times": "open",
// "signname_topic":       "01_classic",
// "signname_music":       "classic",
// "signname_loop":        "open",
// "signname_latest":      "open",
// "signname_edit_times":  0,
// "signname_content":     "write",
// "signname_show_name":   "close",

// "signname_classic_h_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_01.png",
// "signname_classic_h_pic_02": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_02.png",
// "signname_classic_h_pic_03": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_03.png",
// "signname_classic_h_pic_04": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_04.png",
// "signname_classic_h_pic_05": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_05.png",
// "signname_classic_h_pic_06": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_06.jpg",
// "signname_classic_h_pic_07": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_h_pic_07.png",
// "signname_classic_g_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_g_pic_01.jpg",
// "signname_classic_c_pic_01": config.UPLOAD_SYSTEM_URL + "signname/classic/signname_classic_c_pic_01.jpg",

// "signname_bgm": config.UPLOAD_SYSTEM_URL + "signname/classic/bgm/bgm.mp3",

// "overview_vote":     "close",
// "overview_chatroom": "open",
// "overview_interact": "open",
// "overview_info":     "open",

// 一般彈幕
// "danmu_new_background_color": "#90ace7",
// "danmu_new_text_color":       "#1e359d",
// "danmu_old_background_color": "#fdc1f8",
// "danmu_old_text_color":       "#9d1e60",

// "danmu_style":             "classic",
// "danmu_custom_left_new":   "/admin/uploads/system/danmu/custom/left_new.svg",
// "danmu_custom_center_new": "/admin/uploads/system/danmu/custom/center_new.svg",
// "danmu_custom_right_new":  "/admin/uploads/system/danmu/custom/right_new.svg",
// "danmu_custom_left_old":   "/admin/uploads/system/danmu/custom/left_old.svg",
// "danmu_custom_center_old": "/admin/uploads/system/danmu/custom/center_old.svg",
// "danmu_custom_right_old":  "/admin/uploads/system/danmu/custom/right_old.svg",

// 特殊彈幕
// "special_danmu_ban_second": 20,
//霸屏
// "holdscreen_ban_second": 20,

// 報名
// "customize_default_avatar": fmt.Sprintf("%s/admin/uploads/system/img-user-pic.png", config.HTTP_HILIVES_NET_URL),
// }

// 建立活動資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID, os.ModePerm)

// 建立活動介紹資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/info/introduce", os.ModePerm)
// 建立活動嘉賓資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/info/guest", os.ModePerm)

// 建立報名設置資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/applysign/apply", os.ModePerm)
// 建立自定義欄位設置資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/applysign/customize", os.ModePerm)
// 建立QRcode自定義欄位設置資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/applysign/qrcode", os.ModePerm)

// 建立訊息區設置資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/wall/message", os.ModePerm)
// 建立主題區資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/wall/topic", os.ModePerm)
// 建立提問區資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/wall/question", os.ModePerm)
// 建立一般彈幕資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/wall/danmu", os.ModePerm)

// 建立一般簽到資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/sign/general", os.ModePerm)
// 建立3D簽到資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/sign/threed", os.ModePerm)
// 建立簽名牆資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/sign/signname", os.ModePerm)

// 建立遊戲抽獎資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/lottery", os.ModePerm)
// 建立搖紅包資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/redpack", os.ModePerm)
// 建立套紅包資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/ropepack", os.ModePerm)
// 建立搖號抽獎資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/draw_numbers", os.ModePerm)
// 建立敲敲樂資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/whack_mole", os.ModePerm)
// 建立鑑定師資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/monopoly", os.ModePerm)
// 建立快問快答資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/QA", os.ModePerm)
// 建立拔河遊戲資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/tugofwar", os.ModePerm)
// 建立賓果遊戲資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/bingo", os.ModePerm)
// 建立扭蛋機遊戲資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/game/3DGachaMachine", os.ModePerm)
// 建立投票遊戲資料夾
// os.MkdirAll(config.STORE_PATH+"/"+model.UserID+"/"+model.ActivityID+"/interact/sign/vote", os.ModePerm)
