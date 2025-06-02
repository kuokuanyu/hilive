package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strings"
	"time"
)

// ActivityModel 資料表欄位
type ActivityModel struct {
	// 活動
	Base             `json:"-"`
	ID               int64  `json:"id"`
	ActivityID       string `json:"activity_id"`
	UserID           string `json:"user_id"`
	ActivityName     string `json:"activity_name"`
	ActivityType     string `json:"activity_type"`
	MaxPeople        int64  `json:"max_people"`
	People           int64  `json:"people"`
	Attend           int64  `json:"attend"`
	City             string `json:"city"`
	Town             string `json:"town"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Number           int64  `json:"number"`
	LoginRequired    string `json:"login_required"`     // 是否需要登入才可以進入聊天室
	LoginPassword    string `json:"login_password"`     // 當不需要登入就可以進入聊天室時，需要密碼驗證
	PasswordRequired string `json:"password_required"`  // 是否需要設置密碼
	PushMessage      string `json:"push_message"`       // 官方帳號發送訊息
	PushPhoneMessage string `json:"push_phone_message"` // 發送手機簡訊
	MessageAmount    int64  `json:"message_amount"`     // 簡訊數量

	SendMail   string `json:"send_mail"`   // 發送郵件
	MailAmount int64  `json:"mail_amount"` // 郵件數量

	HostScan string `json:"host_scan"` // 主持人掃qrcode

	CustomizePassword   string `json:"customize_password"`    // 是否自定義設置驗證碼
	AllowCustomizeApply string `json:"allow_customize_apply"` // 用戶是否允許自定義報名

	// 官方帳號綁定
	// 活動
	ActivityLineID        string `json:"activity_line_id"`        // line_id
	ActivityChannelID     string `json:"activity_channel_id"`     // channel_id
	ActivityChannelSecret string `json:"activity_channel_secret"` // channel_secret
	ActivityChatbotSecret string `json:"activity_chatbot_secret"` // chatbot_secret
	ActivityChatbotToken  string `json:"activity_chatbot_token"`  // chatbot_token
	// 用戶
	UserLineID        string `json:"user_line_id"`        // line_id
	UserChannelID     string `json:"user_channel_id"`     // channel_id
	UserChannelSecret string `json:"user_channel_secret"` // channel_secret
	UserChatbotSecret string `json:"user_chatbot_secret"` // chatbot_secret
	UserChatbotToken  string `json:"user_chatbot_token"`  // chatbot_token

	Device string `json:"device"` // 驗證裝置

	// 官方帳號綁定
	// 	LineID            string `json:"line_id"`             // line_id
	// ChannelID         string `json:"channel_id"`          // channel_id
	// ChannelSecret     string `json:"channel_secret"`      // channel_secret
	// ChatbotSecret     string `json:"chatbot_secret"`      // chatbot_secret
	// ChatbotToken      string `json:"chatbot_token"`       // chatbot_token

	// 權限
	Permissions []PermissionModel

	// 活動總覽
	OverviewMessage        string `json:"overview_message"`
	OverviewTopic          string `json:"overview_topic"`
	OverviewQuestion       string `json:"overview_question"`
	OverviewDanmu          string `json:"overview_danmu"`
	OverviewSpecialDanmu   string `json:"overview_special_danmu"`
	OverviewPicture        string `json:"overview_picture"`
	OverviewHoldscreen     string `json:"overview_holdscreen"`
	OverviewGeneral        string `json:"overview_general"`
	OverviewThreed         string `json:"overview_threed"`
	OverviewSignname       string `json:"overview_signname"`
	OverviewCountdown      string `json:"overview_countdown"`
	OverviewLottery        string `json:"overview_lottery"`
	OverviewRedpack        string `json:"overview_redpack"`
	OverviewRopepack       string `json:"overview_ropepack"`
	OverviewWhackMole      string `json:"overview_whack_mole"`
	OverviewDrawNumbers    string `json:"overview_draw_numbers"`
	OverviewMonopoly       string `json:"overview_monopoly"`
	OverviewQA             string `json:"overview_qa"`
	OverviewTugofwar       string `json:"overview_tugofwar"`
	OverviewBingo          string `json:"overview_bingo"`
	Overview3DGachaMachine string `json:"overview_3d_gacha_machine"`
	OverviewVote           string `json:"overview_vote"`
	OverviewChatroom       string `json:"overview_chatroom"`
	OverviewInteract       string `json:"overview_interact"`
	OverviewInfo           string `json:"overview_info"`

	// 活動介紹
	IntroduceTitle string `json:"introduce_title"`
	// 活動行程
	ScheduleTitle         string `json:"schedule_title"`
	ScheduleDisplayDate   string `json:"schedule_display_date"`
	ScheduleDisplayDetail string `json:"schedule_display_detail"`
	// 活動嘉賓
	GuestTitle      string `json:"guest_title"`
	GuestBackground string `json:"guest_background"`
	// 活動資料
	MaterialTitle string `json:"material_title"`
	// 報名
	ApplyCheck             string `json:"apply_check"`
	CustomizeDefaultAvatar string `json:"customize_default_avatar"` // 自定義人員預設頭像

	// 簽到
	SignCheck   string `json:"sign_check"`
	SignAllow   string `json:"sign_allow"`
	SignMinutes int64  `json:"sign_minutes"`
	SignManual  string `json:"sign_manual"`
	//QRcode
	QRcodeLogoPicture     string  `json:"qrcode_logo_picture"`
	QRcodeLogoSize        float64 `json:"qrcode_logo_size"`
	QRcodePicturePoint    string  `json:"qrcode_picture_point"`
	QRcodeWhiteDistance   int64   `json:"qrcode_white_distance"`
	QRcodePointColor      string  `json:"qrcode_point_color"`
	QRcodeBackgroundColor string  `json:"qrcode_background_color"`
	// 訊息牆
	MessagePicture               string   `json:"message_picture"`
	MessageAuto                  string   `json:"message_auto"`
	MessageBan                   string   `json:"message_ban"`
	MessageBanSecond             int64    `json:"message_ban_second"`
	MessageRefreshSecond         int64    `json:"message_refresh_second"`
	MessageOpen                  string   `json:"message_open"`
	Message                      []string `json:"message"`
	MessageBackground            string   `json:"message_background"`
	MessageTextColor             string   `json:"message_text_color"`
	MessageScreenTitleColor      string   `json:"message_screen_title_color"`
	MessageTextFrameColor        string   `json:"message_text_frame_color"`
	MessageFrameColor            string   `json:"message_frame_color"`
	MessageScreenTitleFrameColor string   `json:"message_screen_title_frame_color"`
	// 主題牆
	TopicBackground string `json:"topic_background"`
	// 提問牆
	// QuestionMessageCheck  string `json:"question_message_check"`
	QuestionAnonymous string `json:"question_anonymous"`
	// QuestionHideAnswered  string `json:"question_hide_answered"`
	QuestionQrcode     string `json:"question_qrcode"`
	QuestionBackground string `json:"question_background"`
	// QuestionGuestAnswered string `json:"question_guest_answered"`
	// 彈幕
	DanmuLoop               string  `json:"danmu_loop"`
	DanmuTop                string  `json:"danmu_top"`
	DanmuMid                string  `json:"danmu_mid"`
	DanmuBottom             string  `json:"danmu_bottom"`
	DanmuDisplayName        string  `json:"danmu_display_name"`
	DanmuDisplayAvatar      string  `json:"danmu_display_avatar"`
	DanmuSize               float64 `json:"danmu_size"`
	DanmuSpeed              float64 `json:"danmu_speed"`
	DanmuDensity            float64 `json:"danmu_density"`
	DanmuOpacity            float64 `json:"danmu_opacity"`
	DanmuBackground         int64   `json:"danmu_background"`
	DanmuNewBackgroundColor string  `json:"danmu_new_background_color"`
	DanmuNewTextColor       string  `json:"danmu_new_text_color"`
	DanmuOldBackgroundColor string  `json:"danmu_old_background_color"`
	DanmuOldTextColor       string  `json:"danmu_old_text_color"`

	// DanmuBackground string `json:"danmu_background"`
	DanmuStyle           string `json:"danmu_style"`
	DanmuCustomLeftNew   string `json:"danmu_custom_left_new"`
	DanmuCustomCenterNew string `json:"danmu_custom_center_new"`
	DanmuCustomRightNew  string `json:"danmu_custom_right_new"`
	DanmuCustomLeftOld   string `json:"danmu_custom_left_old"`
	DanmuCustomCenterOld string `json:"danmu_custom_center_old"`
	DanmuCustomRightOld  string `json:"danmu_custom_right_old"`

	// 特殊彈幕
	SpecialDanmuMessageCheck string `json:"special_danmu_message_check"`
	SpecialDanmuGeneralPrice int64  `json:"special_danmu_general_price"`
	SpecialDanmuLargePrice   int64  `json:"special_danmu_large_price"`
	SpecialDanmuOverPrice    int64  `json:"special_danmu_over_price"`
	SpecialDanmuTopic        string `json:"special_danmu_topic"`
	SpecialDanmuBanSecond    int64  `json:"special_danmu_ban_second"`

	// 圖片牆
	PictureStartTime    string   `json:"picture_start_time"`
	PictureEndTime      string   `json:"picture_end_time"`
	PictureHideTime     string   `json:"picture_hide_time"`
	PictureSwitchSecond int64    `json:"picture_switch_second"`
	PicturePlayOrder    string   `json:"picture_play_order"`
	Picture             []string `json:"picture"`
	PictureBackground   string   `json:"picture_background"`
	PictureAnimation    int64    `json:"picture_animation"`

	// 霸屏
	HoldscreenPrice         int64  `json:"holdscreen_price"`
	HoldscreenMessageCheck  string `json:"holdscreen_message_check"`
	HoldscreenOnlyPicture   string `json:"holdscreen_only_picture"`
	HoldscreenDisplaySecond int64  `json:"holdscreen_display_second"`
	HoldscreenBirthdayTopic string `json:"holdscreen_birthday_topic"`
	HoldscreenBanSecond     int64  `json:"holdscreen_ban_second"`

	// 一般簽到
	GeneralDisplayPeople string `json:"general_display_people"`
	GeneralStyle         int64  `json:"general_style"`
	GeneralBackground    string `json:"general_background"`
	GeneralTopic         string `json:"general_topic"`
	GeneralMusic         string `json:"general_music"`
	GeneralLoop          string `json:"general_loop"`
	GeneralLatest        string `json:"general_latest"`
	// 編輯次數
	GeneralEditTimes int64 `json:"general_edit_times" example:"0"` // 編輯次數

	// 一般簽到牆自定義圖片
	GeneralClassicHPic01 string `json:"general_classic_h_pic_01" example:"picture"`
	GeneralClassicHPic02 string `json:"general_classic_h_pic_02" example:"picture"`
	GeneralClassicHPic03 string `json:"general_classic_h_pic_03" example:"picture"`
	GeneralClassicHAni01 string `json:"general_classic_h_ani_01" example:"picture"`

	// 音樂
	GeneralBgm string `json:"general_bgm" example:"picture"`

	// 一般簽到自定義圖片陣列參數
	GeneralCustomizeHostPictures      []string `json:"general_customize_host_pictures" example:"pictures"`      // 一般圖片
	GeneralCustomizeGuestPictures     []string `json:"general_customize_guest_pictures" example:"pictures"`     // 一般圖片
	GeneralCustomizeCommonPictures    []string `json:"general_customize_common_pictures" example:"pictures"`    // 一般圖片
	GeneralCustomizeHostAnipictures   []string `json:"general_customize_host_anipictures" example:"pictures"`   // 動圖
	GeneralCustomizeGuestAnipictures  []string `json:"general_customize_guest_anipictures" example:"pictures"`  // 動圖
	GeneralCustomizeCommonAnipictures []string `json:"general_customize_common_anipictures" example:"pictures"` // 動圖
	GeneralCustomizeMusics            []string `json:"general_customize_musics" example:"pictures"`             // 音樂

	// 立體簽到
	ThreedAvatar           string `json:"threed_avatar"`
	ThreedAvatarShape      string `json:"threed_avatar_shape"`
	ThreedDisplayPeople    string `json:"threed_display_people"`
	ThreedBackgroundStyle  int64  `json:"threed_background_style"`
	ThreedBackground       string `json:"threed_background"`
	ThreedImageLogo        string `json:"threed_image_logo"`
	ThreedImageCircle      string `json:"threed_image_circle"`
	ThreedImageSpiral      string `json:"threed_image_spiral"`
	ThreedImageRectangle   string `json:"threed_image_rectangle"`
	ThreedImageSquare      string `json:"threed_image_square"`
	ThreedImageLogoPicture string `json:"threed_image_logo_picture"`
	ThreedTopic            string `json:"threed_topic"`
	ThreedMusic            string `json:"threed_music"`

	// 編輯次數
	ThreedEditTimes int64 `json:"threed_edit_times" example:"0"` // 編輯次數

	// 音樂
	ThreedBgm string `json:"threed_bgm" example:"picture"`

	// 立體簽到自定義圖片陣列參數
	ThreedCustomizePictures    []string `json:"threed_customize_pictures" example:"pictures"`    // 一般圖片
	ThreedCustomizeAnipictures []string `json:"threed_customize_anipictures" example:"pictures"` // 動圖
	ThreedCustomizeMusics      []string `json:"threed_customize_musics" example:"pictures"`      // 音樂

	// 倒數計時
	CountdownSecond      int64  `json:"countdown_second"`
	CountdownURL         string `json:"countdown_url"`
	CountdownAvatar      string `json:"countdown_avatar"`
	CountdownAvatarShape string `json:"countdown_avatar_shape"`
	CountdownBackground  int64  `json:"countdown_background"`

	// 簽名牆
	SignnameMode       string `json:"signname_mode"`
	SignnameTimes      int64  `json:"signname_times"`
	SignnameDisplay    string `json:"signname_display"`
	SignnameLimitTimes string `json:"signname_limit_times"`
	SignnameTopic      string `json:"signname_topic"`
	SignnameMusic      string `json:"signname_music"`
	SignnameLoop       string `json:"signname_loop"`
	SignnameLatest     string `json:"signname_latest"`
	SignnameContent    string `json:"signname_content"`
	SignnameShowName   string `json:"signname_show_name"`
	// 編輯次數
	SignnameEditTimes int64 `json:"signname_edit_times" example:"0"` // 編輯次數

	// 簽名牆自定義圖片
	SignnameClassicHPic01 string `json:"signname_classic_h_pic_01" example:"picture"`
	SignnameClassicHPic02 string `json:"signname_classic_h_pic_02" example:"picture"`
	SignnameClassicHPic03 string `json:"signname_classic_h_pic_03" example:"picture"`
	SignnameClassicHPic04 string `json:"signname_classic_h_pic_04" example:"picture"`
	SignnameClassicHPic05 string `json:"signname_classic_h_pic_05" example:"picture"`
	SignnameClassicHPic06 string `json:"signname_classic_h_pic_06" example:"picture"`
	SignnameClassicHPic07 string `json:"signname_classic_h_pic_07" example:"picture"`
	SignnameClassicGPic01 string `json:"signname_classic_g_pic_01" example:"picture"`
	SignnameClassicCPic01 string `json:"signname_classic_c_pic_01" example:"picture"`

	// 音樂
	SignnameBgm string `json:"signname_bgm" example:"picture"`

	// 簽名牆自定義圖片陣列參數
	SignnameCustomizeHostPictures      []string `json:"signname_customize_host_pictures" example:"pictures"`      // 一般圖片
	SignnameCustomizeGuestPictures     []string `json:"signname_customize_guest_pictures" example:"pictures"`     // 一般圖片
	SignnameCustomizeCommonPictures    []string `json:"signname_customize_common_pictures" example:"pictures"`    // 一般圖片
	SignnameCustomizeHostAnipictures   []string `json:"signname_customize_host_anipictures" example:"pictures"`   // 動圖
	SignnameCustomizeGuestAnipictures  []string `json:"signname_customize_guest_anipictures" example:"pictures"`  // 動圖
	SignnameCustomizeCommonAnipictures []string `json:"signname_customize_common_anipictures" example:"pictures"` // 動圖
	SignnameCustomizeMusics            []string `json:"signname_customize_musics" example:"pictures"`             // 音樂

	// 訊息審核
	MessageCheckManualCheck string `json:"message_check_manual_check"` // 手動審核
	MessageCheckSensitivity string `json:"message_check_sensitivity"`  // 敏感詞

	// 自定義
	Ext1Name     string `json:"ext_1_name"`
	Ext1Type     string `json:"ext_1_type"`
	Ext1Options  string `json:"ext_1_options"`
	Ext1Required string `json:"ext_1_required"`

	Ext2Name     string `json:"ext_2_name"`
	Ext2Type     string `json:"ext_2_type"`
	Ext2Options  string `json:"ext_2_options"`
	Ext2Required string `json:"ext_2_required"`

	Ext3Name     string `json:"ext_3_name"`
	Ext3Type     string `json:"ext_3_type"`
	Ext3Options  string `json:"ext_3_options"`
	Ext3Required string `json:"ext_3_required"`

	Ext4Name     string `json:"ext_4_name"`
	Ext4Type     string `json:"ext_4_type"`
	Ext4Options  string `json:"ext_4_options"`
	Ext4Required string `json:"ext_4_required"`

	Ext5Name     string `json:"ext_5_name"`
	Ext5Type     string `json:"ext_5_type"`
	Ext5Options  string `json:"ext_5_options"`
	Ext5Required string `json:"ext_5_required"`

	Ext6Name     string `json:"ext_6_name"`
	Ext6Type     string `json:"ext_6_type"`
	Ext6Options  string `json:"ext_6_options"`
	Ext6Required string `json:"ext_6_required"`

	Ext7Name     string `json:"ext_7_name"`
	Ext7Type     string `json:"ext_7_type"`
	Ext7Options  string `json:"ext_7_options"`
	Ext7Required string `json:"ext_7_required"`

	Ext8Name     string `json:"ext_8_name"`
	Ext8Type     string `json:"ext_8_type"`
	Ext8Options  string `json:"ext_8_options"`
	Ext8Required string `json:"ext_8_required"`

	Ext9Name     string `json:"ext_9_name"`
	Ext9Type     string `json:"ext_9_type"`
	Ext9Options  string `json:"ext_9_options"`
	Ext9Required string `json:"ext_9_required"`

	Ext10Name     string `json:"ext_10_name"`
	Ext10Type     string `json:"ext_10_type"`
	Ext10Options  string `json:"ext_10_options"`
	Ext10Required string `json:"ext_10_required"`

	ExtEmailRequired string `json:"ext_email_required"`
	ExtPhoneRequired string `json:"ext_phone_required"`
	InfoPicture      string `json:"info_picture"`

	// 用戶
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	MaxActivityPeople int64  `json:"max_activity_people"`
	MaxGamePeople     int64  `json:"max_game_people"`
}

// EditActivityModel 資料表欄位
type EditActivityModel struct {
	UserID       string `json:"user_id"`
	ActivityID   string `json:"activity_id"`
	ActivityName string `json:"activity_name"`
	ActivityType string `json:"activity_type"`
	MaxPeople    string `json:"max_people"`
	People       string `json:"people"`
	// Attend       string `json:"attend"`
	City             string `json:"city"`
	Town             string `json:"town"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	LoginRequired    string `json:"login_required"`    // 是否需要登入才可以進入聊天室
	LoginPassword    string `json:"login_password"`    // 當不需要登入就可以進入聊天室時，需要密碼驗證
	PasswordRequired string `json:"password_required"` // 是否需要設置密碼
	PushMessage      string `json:"push_message"`      // 官方帳號發送訊息
	// EditTimes        string `json:"edit_times" example:"0"` // 編輯次數
	PushPhoneMessage string `json:"push_phone_message"` // 發送手機簡訊
	MessageAmount    string `json:"message_amount"`     // 簡訊數量
	SendMail         string `json:"send_mail"`          // 發送郵件
	MailAmount       string `json:"mail_amount"`        // 郵件數量

	HostScan string `json:"host_scan"` // 主持人掃qrcode

	Permissions string `json:"permissions"` // 權限

	Device string `json:"device"` // 驗證裝置

	// 官方帳號綁定
	LineID        string `json:"line_id"`        // line_id
	ChannelID     string `json:"channel_id"`     // channel_id
	ChannelSecret string `json:"channel_secret"` // channel_secret
	ChatbotSecret string `json:"chatbot_secret"` // chatbot_secret
	ChatbotToken  string `json:"chatbot_token"`  // chatbot_token

	ChannelAmount string `json:"channel_amount"` // 頻道數量
}

// DefaultActivityModel 預設ActivityModel
func DefaultActivityModel() ActivityModel {
	return ActivityModel{Base: Base{TableName: config.ACTIVITY_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a ActivityModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ActivityModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a ActivityModel) SetDbConn(conn db.Connection) ActivityModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (m ActivityModel) SetRedisConn(conn cache.Connection) ActivityModel {
// 	m.RedisConn = conn
// 	return m
// }

// // SetMongoConn 設定connection
// func (m ActivityModel) SetMongoConn(conn mongo.Connection) ActivityModel {
// 	m.MongoConn = conn
// 	return m
// }

// FindActivityPermissions 查詢活動資料、活動權限(join activity_permissions)
func (a ActivityModel) FindActivityPermissions(key, value string) ([]ActivityModel, error) {
	items, err := a.Table(a.Base.TableName).
		Select("activity.id", "activity.activity_id", "activity.user_id",
			"activity.activity_name", "activity.activity_type",
			"activity.max_people", "activity.people", "activity.attend",
			"activity.city", "activity.town", "activity.start_time", "activity.end_time", "activity.number",
			"activity.login_required", "activity.login_password", "activity.password_required",
			"activity.push_message", "activity.message_amount", "activity.push_phone_message",
			"activity.send_mail", "activity.mail_amount",
			"activity.host_scan",
			"activity.customize_password",
			"activity.allow_customize_apply",

			// 用戶官方帳號
			"users.line_id as user_line_id",
			"users.channel_id as user_channel_id",
			"users.channel_secret as user_channel_secret",
			"users.chatbot_secret as user_chatbot_secret",
			"users.chatbot_token as user_chatbot_token",

			// 活動官方帳號
			"activity.line_id as activity_line_id",
			"activity.channel_id as activity_channel_id",
			"activity.channel_secret as activity_channel_secret",
			"activity.chatbot_secret as activity_chatbot_secret",
			"activity.chatbot_token as activity_chatbot_token",

			"activity.device",

			// 活動總覽
			"activity.overview_message", "activity.overview_topic", "activity.overview_question",
			"activity.overview_danmu", "activity.overview_special_danmu", "activity.overview_picture",
			"activity.overview_holdscreen", "activity.overview_general", "activity.overview_threed",
			"activity.overview_countdown", "activity.overview_lottery", "activity.overview_redpack",
			"activity.overview_ropepack", "activity.overview_whack_mole", "activity.overview_draw_numbers",
			"activity.overview_monopoly", "activity.overview_qa", "activity.overview_tugofwar",
			"activity.overview_bingo", "activity.overview_signname", "activity.overview_3d_gacha_machine",
			"activity_2.overview_vote",
			"activity_2.overview_chatroom", "activity_2.overview_interact", "activity_2.overview_info",

			// 活動介紹
			"activity.introduce_title",
			// 活動行程
			"activity.schedule_title", "activity.schedule_display_date", "activity.schedule_display_detail",
			// 活動嘉賓
			"activity.guest_title", "activity.guest_background",
			// 活動資料
			"activity.material_title",
			// 報名
			"activity.apply_check",
			"activity_2.customize_default_avatar",
			// 簽到
			"activity.sign_check", "activity.sign_allow", "activity.sign_minutes", "activity.sign_manual",

			// QRcode自定義
			"activity.qrcode_logo_picture",
			"activity.qrcode_logo_size",
			"activity.qrcode_picture_point",
			"activity.qrcode_white_distance",
			"activity.qrcode_point_color",
			"activity.qrcode_background_color",

			// 訊息牆
			"activity.message_picture", "activity.message_auto",
			"activity.message_ban", "activity.message_ban_second", "activity.message_refresh_second",
			"activity.message_open", "activity.message", "activity.message_background",
			"activity.message_text_color", "activity.message_screen_title_color", "activity.message_text_frame_color",
			"activity.message_frame_color", "activity.message_screen_title_frame_color",
			// 主題牆
			"activity.topic_background",
			// 提問牆
			"activity.question_anonymous", "activity.question_qrcode", "activity.question_background",
			// 彈幕
			"activity.danmu_loop", "activity.danmu_top", "activity.danmu_mid",
			"activity.danmu_bottom", "activity.danmu_display_name", "activity.danmu_display_avatar",
			"activity.danmu_size", "activity.danmu_speed", "activity.danmu_density",
			"activity.danmu_opacity", "activity.danmu_background",
			"activity_2.danmu_new_background_color",
			"activity_2.danmu_new_text_color",
			"activity_2.danmu_old_background_color",
			"activity_2.danmu_old_text_color",

			"activity_2.danmu_style",
			"activity_2.danmu_custom_left_new",
			"activity_2.danmu_custom_center_new",
			"activity_2.danmu_custom_right_new",
			"activity_2.danmu_custom_left_old",
			"activity_2.danmu_custom_center_old",
			"activity_2.danmu_custom_right_old",

			// 特殊彈幕
			"activity.special_danmu_message_check",
			"activity.special_danmu_general_price", "activity.special_danmu_large_price", "activity.special_danmu_over_price",
			"activity.special_danmu_topic",
			"activity_2.special_danmu_ban_second",

			// 圖片牆
			"activity.picture_start_time", "activity.picture_end_time",
			"activity.picture_hide_time", "activity.picture_switch_second", "activity.picture_play_order",
			"activity.picture", "activity.picture_background",
			"activity.picture_animation",

			// 霸屏
			"activity.holdscreen_price", "activity.holdscreen_message_check",
			"activity.holdscreen_only_picture", "activity.holdscreen_display_second",
			"activity.holdscreen_birthday_topic",
			"activity_2.holdscreen_ban_second",

			// 一般簽到
			"activity.general_display_people", "activity.general_style",
			"activity.general_background",
			"activity_2.general_topic",
			"activity_2.general_music",
			"activity_2.general_loop",
			"activity_2.general_latest",
			"activity_2.general_edit_times",

			"activity_2.general_classic_h_pic_01",
			"activity_2.general_classic_h_pic_02",
			"activity_2.general_classic_h_pic_03",
			"activity_2.general_classic_h_ani_01",

			"activity_2.general_bgm",

			// 立體簽到
			"activity.threed_avatar",
			"activity.threed_avatar_shape", "activity.threed_display_people", "activity.threed_background_style",
			"activity.threed_background", "activity.threed_image_logo",
			"activity.threed_image_circle", "activity.threed_image_spiral", "activity.threed_image_rectangle",
			"activity.threed_image_square", "activity.threed_image_logo_picture",
			"activity_2.threed_topic",
			"activity_2.threed_music",
			"activity_2.threed_edit_times",

			"activity_2.threed_bgm",

			// 倒數計時
			"activity.countdown_second", "activity.countdown_url", "activity.countdown_avatar",
			"activity.countdown_avatar_shape", "activity.countdown_background",

			// 簽名牆
			"activity_2.signname_mode",
			"activity_2.signname_times",
			"activity_2.signname_display",
			"activity_2.signname_limit_times",
			"activity_2.signname_topic",
			"activity_2.signname_music",
			"activity_2.signname_loop",
			"activity_2.signname_latest",
			"activity_2.signname_content",
			"activity_2.signname_show_name",
			"activity_2.signname_edit_times",

			"activity_2.signname_classic_h_pic_01",
			"activity_2.signname_classic_h_pic_02",
			"activity_2.signname_classic_h_pic_03",
			"activity_2.signname_classic_h_pic_04",
			"activity_2.signname_classic_h_pic_05",
			"activity_2.signname_classic_h_pic_06",
			"activity_2.signname_classic_h_pic_07",
			"activity_2.signname_classic_g_pic_01",
			"activity_2.signname_classic_c_pic_01",

			"activity_2.signname_bgm",

			// 訊息審核
			"activity.message_check_manual_check", "activity.message_check_sensitivity",

			//權限
			"permission.id as permission_id",
			"permission.permission", "permission.http_method", "permission.http_path",

			// 用戶
			"users.name", "users.avatar",
			"users.max_activity_people", "users.max_game_people",
		).
		OrderBy("activity.id", "asc", "activity_permissions.permission_id", "asc").
		LeftJoin(command.Join{
			FieldA:    "activity_2.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity_2",
			Operation: "="}).
		LeftJoin(command.Join{
			FieldA:    "activity_permissions.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity_permissions",
			Operation: "="}).
		LeftJoin(command.Join{
			FieldA:    "permission.id",
			FieldA1:   "activity_permissions.permission_id",
			Table:     "permission",
			Operation: "="}).
		LeftJoin(command.Join{
			FieldA:    "users.user_id",
			FieldA1:   "activity.user_id",
			Table:     "users",
			Operation: "="}).
		Where(key, "=", value).
		All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得用戶的活動資訊，請重新查詢")
	}

	return MapToActivityModel(items), nil
}

// Find 尋找資料(單個活動)
func (a ActivityModel) Find(isRedis bool, activityID string) (ActivityModel, error) {
	var (
	// item map[string]interface{}
	// value string
	// attend int
	// err    error
	)
	if isRedis {
		// var value string
		// 判斷redis裡是否有活動資訊，有則不執行查詢資料表功能
		dataMap, err := a.RedisConn.HashGetAllCache(config.ACTIVITY_REDIS + activityID)
		if err != nil {
			return ActivityModel{}, errors.New("錯誤: 取得遊戲快取資料發生問題")
		}

		a.ID = utils.GetInt64FromStringMap(dataMap, "id", 0)
		a.ActivityID, _ = dataMap["activity_id"]
		a.UserID, _ = dataMap["user_id"]
		a.SignnameEditTimes = utils.GetInt64FromStringMap(dataMap, "signname_edit_times", 0) // 編輯次數
		a.GeneralEditTimes = utils.GetInt64FromStringMap(dataMap, "general_edit_times", 0)   // 編輯次數
		a.ThreedEditTimes = utils.GetInt64FromStringMap(dataMap, "threed_edit_times", 0)     // 編輯次數
		a.MessageAmount = utils.GetInt64FromStringMap(dataMap, "message_amount", 0)          // 簡訊數量
		a.MessageAmount = utils.GetInt64FromStringMap(dataMap, "mail_amount", 0)             // 郵件數量

		// value, err = a.RedisConn.GetCache(config.ACTIVITY_REDIS + activityID)
		// if err != nil {
		// 	return ActivityModel{}, errors.New("錯誤: 取得活動快取資料發生問題")
		// }

		// if value != "" {
		// 	// reids中已有資料
		// 	attend, err = strconv.Atoi(value)
		// 	if err != nil {
		// 		return ActivityModel{}, errors.New("錯誤: 轉換活動快取資料發生問題")
		// 	}
		// 	a.Attend = int64(attend)
		// }
	}

	// redis中無活動快取資料，從資料表中查詢
	if a.ID == 0 {
		item, err := a.Table(a.Base.TableName).Select(
			"activity.id", "activity.activity_id", "activity.user_id",
			"activity.activity_name", "activity.activity_type",
			"activity.max_people", "activity.people", "activity.attend",
			"activity.city", "activity.town", "activity.start_time", "activity.end_time", "activity.number",
			"activity.login_required", "activity.login_password", "activity.password_required",
			"activity.push_message", "activity.message_amount", "activity.push_phone_message",
			"activity.send_mail", "activity.mail_amount",
			"activity.host_scan",
			"activity.customize_password",
			"activity.allow_customize_apply",

			// 用戶官方帳號
			"users.line_id as user_line_id",
			"users.channel_id as user_channel_id",
			"users.channel_secret as user_channel_secret",
			"users.chatbot_secret as user_chatbot_secret",
			"users.chatbot_token as user_chatbot_token",

			// 活動官方帳號
			"activity.line_id as activity_line_id",
			"activity.channel_id as activity_channel_id",
			"activity.channel_secret as activity_channel_secret",
			"activity.chatbot_secret as activity_chatbot_secret",
			"activity.chatbot_token as activity_chatbot_token",

			"activity.device",

			// 活動總覽
			"activity.overview_message", "activity.overview_topic", "activity.overview_question",
			"activity.overview_danmu", "activity.overview_special_danmu", "activity.overview_picture",
			"activity.overview_holdscreen", "activity.overview_general", "activity.overview_threed",
			"activity.overview_countdown", "activity.overview_lottery", "activity.overview_redpack",
			"activity.overview_ropepack", "activity.overview_whack_mole", "activity.overview_draw_numbers",
			"activity.overview_monopoly", "activity.overview_qa", "activity.overview_tugofwar",
			"activity.overview_bingo", "activity.overview_signname", "activity.overview_3d_gacha_machine",
			"activity_2.overview_vote",
			"activity_2.overview_chatroom", "activity_2.overview_interact", "activity_2.overview_info",

			// 活動介紹
			"activity.introduce_title",
			// 活動行程
			"activity.schedule_title", "activity.schedule_display_date", "activity.schedule_display_detail",
			// 活動嘉賓
			"activity.guest_title", "activity.guest_background",
			// 活動資料
			"activity.material_title",
			// 報名
			"activity.apply_check",
			"activity_2.customize_default_avatar",

			// 簽到
			"activity.sign_check", "activity.sign_allow", "activity.sign_minutes", "activity.sign_manual",
			// QRcode自定義
			"activity.qrcode_logo_picture",
			"activity.qrcode_logo_size",
			"activity.qrcode_picture_point",
			"activity.qrcode_white_distance",
			"activity.qrcode_point_color",
			"activity.qrcode_background_color",
			// 訊息牆
			"activity.message_picture", "activity.message_auto",
			"activity.message_ban", "activity.message_ban_second", "activity.message_refresh_second",
			"activity.message_open", "activity.message", "activity.message_background",
			"activity.message_text_color", "activity.message_screen_title_color", "activity.message_text_frame_color",
			"activity.message_frame_color", "activity.message_screen_title_frame_color",
			// 主題牆
			"activity.topic_background",
			// 提問牆
			"activity.question_anonymous", "activity.question_qrcode", "activity.question_background",
			// 彈幕
			"activity.danmu_loop", "activity.danmu_top", "activity.danmu_mid",
			"activity.danmu_bottom", "activity.danmu_display_name", "activity.danmu_display_avatar",
			"activity.danmu_size", "activity.danmu_speed", "activity.danmu_density",
			"activity.danmu_opacity", "activity.danmu_background",
			"activity_2.danmu_new_background_color",
			"activity_2.danmu_new_text_color",
			"activity_2.danmu_old_background_color",
			"activity_2.danmu_old_text_color",

			"activity_2.danmu_style",
			"activity_2.danmu_custom_left_new",
			"activity_2.danmu_custom_center_new",
			"activity_2.danmu_custom_right_new",
			"activity_2.danmu_custom_left_old",
			"activity_2.danmu_custom_center_old",
			"activity_2.danmu_custom_right_old",

			// 特殊彈幕
			"activity.special_danmu_message_check",
			"activity.special_danmu_general_price", "activity.special_danmu_large_price", "activity.special_danmu_over_price",
			"activity.special_danmu_topic",
			"activity_2.special_danmu_ban_second",

			// 圖片牆
			"activity.picture_start_time", "activity.picture_end_time",
			"activity.picture_hide_time", "activity.picture_switch_second", "activity.picture_play_order",
			"activity.picture", "activity.picture_background",
			"activity.picture_animation",
			// 霸屏
			"activity.holdscreen_price", "activity.holdscreen_message_check",
			"activity.holdscreen_only_picture", "activity.holdscreen_display_second",
			"activity.holdscreen_birthday_topic",
			"activity_2.holdscreen_ban_second",

			// 一般簽到
			"activity.general_display_people", "activity.general_style",
			"activity.general_background",
			"activity_2.general_topic",
			"activity_2.general_music",
			"activity_2.general_loop",
			"activity_2.general_latest",
			"activity_2.general_edit_times",

			"activity_2.general_classic_h_pic_01",
			"activity_2.general_classic_h_pic_02",
			"activity_2.general_classic_h_pic_03",
			"activity_2.general_classic_h_ani_01",

			"activity_2.general_bgm",

			// 立體簽到
			"activity.threed_avatar",
			"activity.threed_avatar_shape", "activity.threed_display_people", "activity.threed_background_style",
			"activity.threed_background", "activity.threed_image_logo",
			"activity.threed_image_circle", "activity.threed_image_spiral", "activity.threed_image_rectangle",
			"activity.threed_image_square", "activity.threed_image_logo_picture",
			"activity_2.threed_topic",
			"activity_2.threed_music",
			"activity_2.threed_edit_times",

			"activity_2.threed_bgm",

			// 倒數計時
			"activity.countdown_second", "activity.countdown_url", "activity.countdown_avatar",
			"activity.countdown_avatar_shape", "activity.countdown_background",

			// 簽名牆
			"activity_2.signname_mode",
			"activity_2.signname_times",
			"activity_2.signname_display",
			"activity_2.signname_limit_times",
			"activity_2.signname_topic",
			"activity_2.signname_music",
			"activity_2.signname_loop",
			"activity_2.signname_latest",
			"activity_2.signname_content",
			"activity_2.signname_show_name",
			"activity_2.signname_edit_times",

			"activity_2.signname_classic_h_pic_01",
			"activity_2.signname_classic_h_pic_02",
			"activity_2.signname_classic_h_pic_03",
			"activity_2.signname_classic_h_pic_04",
			"activity_2.signname_classic_h_pic_05",
			"activity_2.signname_classic_h_pic_06",
			"activity_2.signname_classic_h_pic_07",
			"activity_2.signname_classic_g_pic_01",
			"activity_2.signname_classic_c_pic_01",

			"activity_2.signname_bgm",

			// 訊息審核
			"activity.message_check_manual_check", "activity.message_check_sensitivity",

			// 用戶
			"users.name", "users.avatar",
			"users.max_activity_people", "users.max_game_people",
		).
			Where("activity.activity_id", "=", activityID).
			LeftJoin(command.Join{
				FieldA:    "activity_2.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity_2",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "users.user_id",
				FieldA1:   "activity.user_id",
				Table:     "users",
				Operation: "="}).
			First()
		if err != nil {
			return ActivityModel{}, errors.New("錯誤: 無法取得活動資訊，請重新查詢")
		}
		if item == nil {
			return ActivityModel{}, nil
		}
		a = a.MapToModel(item)

		if isRedis {
			values := []interface{}{config.ACTIVITY_REDIS + activityID}
			values = append(values, "id", a.ID)
			values = append(values, "activity_id", a.ActivityID)
			values = append(values, "user_id", a.UserID)
			values = append(values, "signname_edit_times", a.SignnameEditTimes)
			values = append(values, "general_edit_times", a.GeneralEditTimes)
			values = append(values, "threed_edit_times", a.ThreedEditTimes)
			values = append(values, "message_amount", a.MessageAmount)
			values = append(values, "mail_amount", a.MessageAmount)

			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
				return a, errors.New("錯誤: 設置活動快取資料發生問題")
			}

			// 設置過期時間
			// a.RedisConn.SetExpire(config.ACTIVITY_REDIS+activityID,
			// 	config.REDIS_EXPIRE)

			// if _, err = a.RedisConn.SetCache(config.ACTIVITY_REDIS+activityID,
			// 	a.Attend, "EX", config.REDIS_EXPIRE); err != nil {
			// 	return ActivityModel{}, errors.New("錯誤: 設置活動快取資料發生問題")
			// }
		}
	}
	return a, nil
}

// FindOpenActivitys 尋找資料(進行中所有活動)
func (a ActivityModel) FindOpenActivitys() ([]ActivityModel, error) {
	var (
	// now, _ = time.ParseInLocation("2006-01-02 15:04:05",
	// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)

	items, err := a.Table(a.Base.TableName).
		// Where("start_time", "<=", now).
		// Where("end_time", ">=", now).
		All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得進行中的活動資訊，請重新查詢")
	}

	return MapToActivityModel(items), nil
}

// FindActivityLeftJoinCustomize 查詢活動資料、自定義資料(join activity_customize，單個活動)
func (a ActivityModel) FindActivityLeftJoinCustomize(activityID string) (ActivityModel, error) {
	var (
		item   map[string]interface{}
		err    error
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)

	item, err = a.Table(a.Base.TableName).
		Select("activity.id",
			// 報名
			"activity.apply_check",

			// 用戶官方帳號
			"users.line_id as user_line_id",
			"users.channel_id as user_channel_id",
			"users.channel_secret as user_channel_secret",
			"users.chatbot_secret as user_chatbot_secret",
			"users.chatbot_token as user_chatbot_token",

			// 活動官方帳號
			"activity.line_id as activity_line_id",
			"activity.channel_id as activity_channel_id",
			"activity.channel_secret as activity_channel_secret",
			"activity.chatbot_secret as activity_chatbot_secret",
			"activity.chatbot_token as activity_chatbot_token",

			// 官方帳號發送訊息
			"activity.push_message",

			"activity.device",

			// 自定義
			"activity_customize.ext_1_required",
			"activity_customize.ext_2_required",
			"activity_customize.ext_3_required",
			"activity_customize.ext_4_required",
			"activity_customize.ext_5_required",
			"activity_customize.ext_6_required",
			"activity_customize.ext_7_required",
			"activity_customize.ext_8_required",
			"activity_customize.ext_9_required",
			"activity_customize.ext_10_required",
			"activity_customize.ext_email_required",
			"activity_customize.ext_phone_required",
		).
		LeftJoin(command.Join{
			FieldA:    "activity.activity_id",
			FieldA1:   "activity_customize.activity_id",
			Table:     "activity_customize",
			Operation: "=",
		}).
		LeftJoin(command.Join{
			FieldA:    "users.user_id",
			FieldA1:   "activity.user_id",
			Table:     "users",
			Operation: "="}).
		Where("activity.activity_id", "=", activityID).
		Where("end_time", ">=", now).First()
	if err != nil {
		return ActivityModel{}, errors.New("錯誤: 無法取得活動資訊，請重新查詢")
	}
	if item == nil {
		return ActivityModel{}, nil
	}
	return a.MapToModel(item), nil
}

// IsEnd 活動是否結束
func (a ActivityModel) IsEnd(activityID string) (ActivityModel, bool) {
	var (
		item   map[string]interface{}
		err    error
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
		activityModel ActivityModel
	)

	item, err = a.Table(a.Base.TableName).Where("activity_id", "=", activityID).
		Where("end_time", ">=", now).First()
	if err != nil {
		// return ActivityModel{}, errors.New("錯誤: 無法取得活動狀態，請重新查詢")
		return activityModel, true
	}
	if item == nil {
		return activityModel, true
	}

	// 判斷活動是否結束
	activityModel = a.MapToModel(item)
	if activityModel.ID == 0 {
		return activityModel, true
	}

	return activityModel, false
}

// MapToActivityModel map轉換[]ActivityModel
func MapToActivityModel(items []map[string]interface{}) []ActivityModel {
	var (
		activitys = make([]ActivityModel, 0)
		activity  = ActivityModel{}
	)

	// fmt.Println("活動權限數量: ", len(items), items)

	for i := 0; i < len(items); i++ {
		activityID, _ := items[i]["activity_id"].(string)

		// 第一筆資料不用比較
		if i == 0 {
			// 活動

			// json解碼，轉換成strcut
			b, _ := json.Marshal(items[i])
			json.Unmarshal(b, &activity)

			activity.Permissions = append(activity.Permissions, DefaultPermissionModel().MapToModel(items[i]))

			// QRcode自定義
			// qrcodeLogoPicture, _ := items[i]["qrcode_logo_picture"].(string)
			if !strings.Contains(activity.QRcodeLogoPicture, "system") {
				activity.QRcodeLogoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/applysign/qrcode/" + activity.QRcodeLogoPicture
			}

			// 訊息陣列處理
			message, _ := items[i]["message"].(string)
			activity.Message = strings.Split(message, "\n")

			// 圖片
			// messageBackground, _ := items[i]["message_background"].(string)
			if !strings.Contains(activity.MessageBackground, "system") {
				activity.MessageBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/message/" + activity.MessageBackground
			}

			// 主題牆
			// 圖片
			// topicBackground, _ := items[i]["topic_background"].(string)
			if !strings.Contains(activity.TopicBackground, "system") {
				activity.TopicBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/topic/" + activity.TopicBackground
			}

			// 提問牆
			// 圖片
			// questionBackground, _ := items[i]["question_background"].(string)
			if !strings.Contains(activity.QuestionBackground, "system") {
				activity.QuestionBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/question/" + activity.QuestionBackground
			}

			// 彈幕
			// danmuCustomLeftNew, _ := items[i]["danmu_custom_left_new"].(string)
			if !strings.Contains(activity.DanmuCustomLeftNew, "system") {
				activity.DanmuCustomLeftNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomLeftNew
			}

			// danmuCustomCenterNew, _ := items[i]["danmu_custom_center_new"].(string)
			if !strings.Contains(activity.DanmuCustomCenterNew, "system") {
				activity.DanmuCustomCenterNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomCenterNew
			}

			// danmuCustomRightNew, _ := items[i]["danmu_custom_right_new"].(string)
			if !strings.Contains(activity.DanmuCustomRightNew, "system") {
				activity.DanmuCustomRightNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomRightNew
			}

			// danmuCustomLeftOld, _ := items[i]["danmu_custom_left_old"].(string)
			if !strings.Contains(activity.DanmuCustomLeftOld, "system") {
				activity.DanmuCustomLeftOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomLeftOld
			}

			// danmuCustomCenterOld, _ := items[i]["danmu_custom_center_old"].(string)
			if !strings.Contains(activity.DanmuCustomCenterOld, "system") {
				activity.DanmuCustomCenterOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomCenterOld
			}

			// danmuCustomRightOld, _ := items[i]["danmu_custom_right_old"].(string)
			if !strings.Contains(activity.DanmuCustomRightOld, "system") {
				activity.DanmuCustomRightOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomRightOld
			}

			// 圖片牆
			picture, _ := items[i]["picture"].(string)
			activity.Picture = strings.Split(picture, "\n")

			// 一般簽到
			// 圖片
			// generalBackground, _ := items[i]["general_background"].(string)
			if !strings.Contains(activity.GeneralBackground, "system") {
				activity.GeneralBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + activity.GeneralBackground
			}

			// 自定義圖片參數(陣列)
			var (
				hpics  = make([]string, 0) // 主持靜態
				gpics  = make([]string, 0) // 玩家靜態
				cpics  = make([]string, 0) // 共用靜態
				hanis  = make([]string, 0) // 主持動態
				ganis  = make([]string, 0) // 玩家動態
				canis  = make([]string, 0) // 共用動態
				musics = make([]string, 0)
			)

			// 自定義圖片
			if activity.GeneralTopic == "01_classic" {
				hpics = append(hpics, activity.GeneralClassicHPic01)
				hpics = append(hpics, activity.GeneralClassicHPic02)
				hpics = append(hpics, activity.GeneralClassicHPic03)
				hanis = append(hanis, activity.GeneralClassicHAni01)
			}
			musics = append(musics, activity.GeneralBgm)

			// 主持靜態
			for i, picture := range hpics {
				if !strings.Contains(picture, "system") {
					hpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range gpics {
				if !strings.Contains(picture, "system") {
					gpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 共用靜態
			for i, picture := range cpics {
				if !strings.Contains(picture, "system") {
					cpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}

			// 主持靜態
			for i, picture := range hanis {
				if !strings.Contains(picture, "system") {
					hanis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range ganis {
				if !strings.Contains(picture, "system") {
					ganis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 共用靜態
			for i, picture := range canis {
				if !strings.Contains(picture, "system") {
					canis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}

			// 音樂
			for i, music := range musics {
				if !strings.Contains(music, "system") {
					musics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + music
				}
			}

			activity.GeneralCustomizeHostPictures = hpics      // 主持靜態
			activity.GeneralCustomizeGuestPictures = gpics     // 玩家靜態
			activity.GeneralCustomizeCommonPictures = cpics    // 共用靜態
			activity.GeneralCustomizeHostAnipictures = hanis   // 主持動態
			activity.GeneralCustomizeGuestAnipictures = ganis  // 玩家動態
			activity.GeneralCustomizeCommonAnipictures = canis // 共用動態
			activity.GeneralCustomizeMusics = musics           // 音樂

			// 立體簽到
			// 圖片
			// threedAvatar, _ := items[i]["threed_avatar"].(string)
			if !strings.Contains(activity.ThreedAvatar, "system") {
				activity.ThreedAvatar = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedAvatar
			}

			// 圖片
			// threedBackground, _ := items[i]["threed_background"].(string)
			if !strings.Contains(activity.ThreedBackground, "system") {
				activity.ThreedBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedBackground
			}

			// 圖片
			// threedImageLogoPicture, _ := items[i]["threed_image_logo_picture"].(string)
			if !strings.Contains(activity.ThreedImageLogoPicture, "system") {
				activity.ThreedImageLogoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedImageLogoPicture
			}

			// 自定義圖片參數(陣列)
			var (
				threedPictures    = make([]string, 0)
				threedAnipictures = make([]string, 0)
				threedMusics      = make([]string, 0)
			)

			// 自定義圖片
			if activity.ThreedTopic == "01_classic" {
				threedPictures = append(threedPictures, activity.ThreedBackground)
			}
			threedMusics = append(threedMusics, activity.ThreedBgm)

			// 判斷是否為自定義上傳資料
			for i, picture := range threedPictures {
				if !strings.Contains(picture, "system") {
					threedPictures[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + picture
				}
			}
			for i, music := range threedMusics {
				if !strings.Contains(music, "system") {
					threedMusics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + music
				}
			}

			activity.ThreedCustomizePictures = threedPictures
			activity.ThreedCustomizeAnipictures = threedAnipictures
			activity.ThreedCustomizeMusics = threedMusics

			// 自定義圖片參數(陣列)
			hpics = make([]string, 0) // 主持靜態
			gpics = make([]string, 0) // 玩家靜態
			cpics = make([]string, 0) // 共用靜態
			hanis = make([]string, 0) // 主持動態
			ganis = make([]string, 0) // 玩家動態
			canis = make([]string, 0) // 共用動態
			musics = make([]string, 0)

			// 自定義圖片
			if activity.SignnameTopic == "01_classic" {
				hpics = append(hpics, activity.SignnameClassicHPic01)
				hpics = append(hpics, activity.SignnameClassicHPic02)
				hpics = append(hpics, activity.SignnameClassicHPic03)
				hpics = append(hpics, activity.SignnameClassicHPic04)
				hpics = append(hpics, activity.SignnameClassicHPic05)
				hpics = append(hpics, activity.SignnameClassicHPic06)
				hpics = append(hpics, activity.SignnameClassicHPic07)
				gpics = append(cpics, activity.SignnameClassicGPic01)

			}
			musics = append(musics, activity.SignnameBgm)

			// 主持靜態
			for i, picture := range hpics {
				if !strings.Contains(picture, "system") {
					hpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range gpics {
				if !strings.Contains(picture, "system") {
					gpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 共用靜態
			for i, picture := range cpics {
				if !strings.Contains(picture, "system") {
					cpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}

			// 主持靜態
			for i, picture := range hanis {
				if !strings.Contains(picture, "system") {
					hanis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range ganis {
				if !strings.Contains(picture, "system") {
					ganis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 共用靜態
			for i, picture := range canis {
				if !strings.Contains(picture, "system") {
					canis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}

			// 音樂
			for i, music := range musics {
				if !strings.Contains(music, "system") {
					musics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + music
				}
			}

			activity.SignnameCustomizeHostPictures = hpics      // 主持靜態
			activity.SignnameCustomizeGuestPictures = gpics     // 玩家靜態
			activity.SignnameCustomizeCommonPictures = cpics    // 共用靜態
			activity.SignnameCustomizeHostAnipictures = hanis   // 主持動態
			activity.SignnameCustomizeGuestAnipictures = ganis  // 玩家動態
			activity.SignnameCustomizeCommonAnipictures = canis // 共用動態
			activity.SignnameCustomizeMusics = musics           // 音樂

			if activity.InfoPicture != "" {
				activity.InfoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/applysign/customize/" + activity.InfoPicture
			}
		} else if activity.ActivityID == activityID {
			// 還是同一個活動的權限資料，將權限資料加入活動中
			activity.Permissions = append(activity.Permissions, DefaultPermissionModel().MapToModel(items[i]))
		} else if activity.ActivityID != activityID {
			// 不是上一個活動的權限資料，將上一個活動資料放入activitys中
			activitys = append(activitys, activity)

			// 清空activity參數資料
			activity = ActivityModel{}

			// 活動

			// json解碼，轉換成strcut
			b, _ := json.Marshal(items[i])
			json.Unmarshal(b, &activity)

			activity.Permissions = append(activity.Permissions, DefaultPermissionModel().MapToModel(items[i]))

			// QRcode自定義
			// qrcodeLogoPicture, _ := items[i]["qrcode_logo_picture"].(string)
			if !strings.Contains(activity.QRcodeLogoPicture, "system") {
				activity.QRcodeLogoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/applysign/qrcode/" + activity.QRcodeLogoPicture
			}

			// 訊息陣列處理
			message, _ := items[i]["message"].(string)
			activity.Message = strings.Split(message, "\n")

			// 圖片
			// messageBackground, _ := items[i]["message_background"].(string)
			if !strings.Contains(activity.MessageBackground, "system") {
				activity.MessageBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/message/" + activity.MessageBackground
			}

			// 主題牆
			// 圖片
			// topicBackground, _ := items[i]["topic_background"].(string)
			if !strings.Contains(activity.TopicBackground, "system") {
				activity.TopicBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/topic/" + activity.TopicBackground
			}

			// 提問牆
			// 圖片
			// questionBackground, _ := items[i]["question_background"].(string)
			if !strings.Contains(activity.QuestionBackground, "system") {
				activity.QuestionBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/question/" + activity.QuestionBackground
			}

			// 彈幕
			// danmuCustomLeftNew, _ := items[i]["danmu_custom_left_new"].(string)
			if !strings.Contains(activity.DanmuCustomLeftNew, "system") {
				activity.DanmuCustomLeftNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomLeftNew
			}

			// danmuCustomCenterNew, _ := items[i]["danmu_custom_center_new"].(string)
			if !strings.Contains(activity.DanmuCustomCenterNew, "system") {
				activity.DanmuCustomCenterNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomCenterNew
			}

			// danmuCustomRightNew, _ := items[i]["danmu_custom_right_new"].(string)
			if !strings.Contains(activity.DanmuCustomRightNew, "system") {
				activity.DanmuCustomRightNew = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomRightNew
			}

			// danmuCustomLeftOld, _ := items[i]["danmu_custom_left_old"].(string)
			if !strings.Contains(activity.DanmuCustomLeftOld, "system") {
				activity.DanmuCustomLeftOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomLeftOld
			}

			// danmuCustomCenterOld, _ := items[i]["danmu_custom_center_old"].(string)
			if !strings.Contains(activity.DanmuCustomCenterOld, "system") {
				activity.DanmuCustomCenterOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomCenterOld
			}

			// danmuCustomRightOld, _ := items[i]["danmu_custom_right_old"].(string)
			if !strings.Contains(activity.DanmuCustomRightOld, "system") {
				activity.DanmuCustomRightOld = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/wall/danmu/" + activity.DanmuCustomRightOld
			}

			// 圖片牆
			picture, _ := items[i]["picture"].(string)
			activity.Picture = strings.Split(picture, "\n")

			// 一般簽到
			// 圖片
			// generalBackground, _ := items[i]["general_background"].(string)
			if !strings.Contains(activity.GeneralBackground, "system") {
				activity.GeneralBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + activity.GeneralBackground
			}

			// 自定義圖片參數(陣列)
			var (
				hpics  = make([]string, 0) // 主持靜態
				gpics  = make([]string, 0) // 玩家靜態
				cpics  = make([]string, 0) // 共用靜態
				hanis  = make([]string, 0) // 主持動態
				ganis  = make([]string, 0) // 玩家動態
				canis  = make([]string, 0) // 共用動態
				musics = make([]string, 0)
			)

			// 自定義圖片
			if activity.GeneralTopic == "01_classic" {
				hpics = append(hpics, activity.GeneralClassicHPic01)
				hpics = append(hpics, activity.GeneralClassicHPic02)
				hpics = append(hpics, activity.GeneralClassicHPic03)
				hanis = append(hanis, activity.GeneralClassicHAni01)
			}
			musics = append(musics, activity.GeneralBgm)

			// 主持靜態
			for i, picture := range hpics {
				if !strings.Contains(picture, "system") {
					hpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range gpics {
				if !strings.Contains(picture, "system") {
					gpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 共用靜態
			for i, picture := range cpics {
				if !strings.Contains(picture, "system") {
					cpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}

			// 主持靜態
			for i, picture := range hanis {
				if !strings.Contains(picture, "system") {
					hanis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range ganis {
				if !strings.Contains(picture, "system") {
					ganis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}
			// 共用靜態
			for i, picture := range canis {
				if !strings.Contains(picture, "system") {
					canis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + picture
				}
			}

			// 音樂
			for i, music := range musics {
				if !strings.Contains(music, "system") {
					musics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/general/" + music
				}
			}

			activity.GeneralCustomizeHostPictures = hpics      // 主持靜態
			activity.GeneralCustomizeGuestPictures = gpics     // 玩家靜態
			activity.GeneralCustomizeCommonPictures = cpics    // 共用靜態
			activity.GeneralCustomizeHostAnipictures = hanis   // 主持動態
			activity.GeneralCustomizeGuestAnipictures = ganis  // 玩家動態
			activity.GeneralCustomizeCommonAnipictures = canis // 共用動態
			activity.GeneralCustomizeMusics = musics           // 音樂

			// 立體簽到
			// 圖片
			// threedAvatar, _ := items[i]["threed_avatar"].(string)
			if !strings.Contains(activity.ThreedAvatar, "system") {
				activity.ThreedAvatar = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedAvatar
			}

			// 圖片
			// threedBackground, _ := items[i]["threed_background"].(string)
			if !strings.Contains(activity.ThreedBackground, "system") {
				activity.ThreedBackground = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedBackground
			}

			// 圖片
			// threedImageLogoPicture, _ := items[i]["threed_image_logo_picture"].(string)
			if !strings.Contains(activity.ThreedImageLogoPicture, "system") {
				activity.ThreedImageLogoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + activity.ThreedImageLogoPicture
			}

			// 自定義圖片參數(陣列)
			var (
				threedPictures    = make([]string, 0)
				threedAnipictures = make([]string, 0)
				threedMusics      = make([]string, 0)
			)

			// 自定義圖片
			if activity.ThreedTopic == "01_classic" {
				threedPictures = append(threedPictures, activity.ThreedBackground)
			}
			threedMusics = append(threedMusics, activity.ThreedBgm)

			// 判斷是否為自定義上傳資料
			for i, picture := range threedPictures {
				if !strings.Contains(picture, "system") {
					threedPictures[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + picture
				}
			}
			for i, music := range threedMusics {
				if !strings.Contains(music, "system") {
					threedMusics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/threed/" + music
				}
			}

			activity.ThreedCustomizePictures = threedPictures
			activity.ThreedCustomizeAnipictures = threedAnipictures
			activity.ThreedCustomizeMusics = threedMusics

			// 自定義圖片參數(陣列)
			hpics = make([]string, 0) // 主持靜態
			gpics = make([]string, 0) // 玩家靜態
			cpics = make([]string, 0) // 共用靜態
			hanis = make([]string, 0) // 主持動態
			ganis = make([]string, 0) // 玩家動態
			canis = make([]string, 0) // 共用動態
			musics = make([]string, 0)

			// 自定義圖片
			if activity.SignnameTopic == "01_classic" {
				hpics = append(hpics, activity.SignnameClassicHPic01)
				hpics = append(hpics, activity.SignnameClassicHPic02)
				hpics = append(hpics, activity.SignnameClassicHPic03)
				hpics = append(hpics, activity.SignnameClassicHPic04)
				hpics = append(hpics, activity.SignnameClassicHPic05)
				hpics = append(hpics, activity.SignnameClassicHPic06)
				hpics = append(hpics, activity.SignnameClassicHPic07)
				gpics = append(cpics, activity.SignnameClassicGPic01)

			}
			musics = append(musics, activity.SignnameBgm)

			// 主持靜態
			for i, picture := range hpics {
				if !strings.Contains(picture, "system") {
					hpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range gpics {
				if !strings.Contains(picture, "system") {
					gpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 共用靜態
			for i, picture := range cpics {
				if !strings.Contains(picture, "system") {
					cpics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}

			// 主持靜態
			for i, picture := range hanis {
				if !strings.Contains(picture, "system") {
					hanis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 玩家靜態
			for i, picture := range ganis {
				if !strings.Contains(picture, "system") {
					ganis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}
			// 共用靜態
			for i, picture := range canis {
				if !strings.Contains(picture, "system") {
					canis[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + picture
				}
			}

			// 音樂
			for i, music := range musics {
				if !strings.Contains(music, "system") {
					musics[i] = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/interact/sign/signname/" + music
				}
			}

			activity.SignnameCustomizeHostPictures = hpics      // 主持靜態
			activity.SignnameCustomizeGuestPictures = gpics     // 玩家靜態
			activity.SignnameCustomizeCommonPictures = cpics    // 共用靜態
			activity.SignnameCustomizeHostAnipictures = hanis   // 主持動態
			activity.SignnameCustomizeGuestAnipictures = ganis  // 玩家動態
			activity.SignnameCustomizeCommonAnipictures = canis // 共用動態
			activity.SignnameCustomizeMusics = musics           // 音樂

			if activity.InfoPicture != "" {
				activity.InfoPicture = "/admin/uploads/" + activity.UserID + "/" + activity.ActivityID + "/applysign/customize/" + activity.InfoPicture
			}
		}

		if i == len(items)-1 {
			// 最後一筆資料
			activitys = append(activitys, activity)
		}
	}
	return activitys
}

// MapToModel map轉換model
func (a ActivityModel) MapToModel(m map[string]interface{}) ActivityModel {
	// 活動

	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	// qrcodeLogoPicture, _ := m["qrcode_logo_picture"].(string)
	if !strings.Contains(a.QRcodeLogoPicture, "system") {
		a.QRcodeLogoPicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/applysign/qrcode/" + a.QRcodeLogoPicture
	}

	// 訊息牆
	message, _ := m["message"].(string)
	a.Message = strings.Split(message, "\n")

	// 圖片
	// messageBackground, _ := m["message_background"].(string)
	if !strings.Contains(a.MessageBackground, "system") {
		a.MessageBackground = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/message/" + a.MessageBackground
	}

	// 主題牆
	// topicBackground, _ := m["topic_background"].(string)
	if !strings.Contains(a.TopicBackground, "system") {
		a.TopicBackground = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/topic/" + a.TopicBackground
	}

	// 提問牆
	// 圖片
	// questionBackground, _ := m["question_background"].(string)
	if !strings.Contains(a.QuestionBackground, "system") {
		a.QuestionBackground = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/question/" + a.QuestionBackground
	}

	// 彈幕
	// danmuCustomLeftNew, _ := m["danmu_custom_left_new"].(string)
	if !strings.Contains(a.DanmuCustomLeftNew, "system") {
		a.DanmuCustomLeftNew = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomLeftNew
	}

	// danmuCustomCenterNew, _ := m["danmu_custom_center_new"].(string)
	if !strings.Contains(a.DanmuCustomCenterNew, "system") {
		a.DanmuCustomCenterNew = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomCenterNew
	}

	// danmuCustomRightNew, _ := m["danmu_custom_right_new"].(string)
	if !strings.Contains(a.DanmuCustomRightNew, "system") {
		a.DanmuCustomRightNew = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomRightNew
	}

	// danmuCustomLeftOld, _ := m["danmu_custom_left_old"].(string)
	if !strings.Contains(a.DanmuCustomLeftOld, "system") {
		a.DanmuCustomLeftOld = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomLeftOld
	}

	// danmuCustomCenterOld, _ := m["danmu_custom_center_old"].(string)
	if !strings.Contains(a.DanmuCustomCenterOld, "system") {
		a.DanmuCustomCenterOld = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomCenterOld
	}

	// danmuCustomRightOld, _ := m["danmu_custom_right_old"].(string)
	if !strings.Contains(a.DanmuCustomRightOld, "system") {
		a.DanmuCustomRightOld = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/wall/danmu/" + a.DanmuCustomRightOld
	}

	// 圖片牆
	picture, _ := m["picture"].(string)
	a.Picture = strings.Split(picture, "\n")

	// 一般簽到
	// 圖片
	// generalBackground, _ := m["general_background"].(string)
	if !strings.Contains(a.GeneralBackground, "system") {
		a.GeneralBackground = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + a.GeneralBackground
	}

	a.GeneralBgm, _ = m["general_bgm"].(string)

	// 自定義圖片參數(陣列)
	var (
		hpics  = make([]string, 0) // 主持靜態
		gpics  = make([]string, 0) // 玩家靜態
		cpics  = make([]string, 0) // 共用靜態
		hanis  = make([]string, 0) // 主持動態
		ganis  = make([]string, 0) // 玩家動態
		canis  = make([]string, 0) // 共用動態
		musics = make([]string, 0)
	)

	// 自定義圖片
	if a.GeneralTopic == "01_classic" {
		hpics = append(hpics, a.GeneralClassicHPic01)
		hpics = append(hpics, a.GeneralClassicHPic02)
		hpics = append(hpics, a.GeneralClassicHPic03)
		hanis = append(hanis, a.GeneralClassicHAni01)
	}
	musics = append(musics, a.GeneralBgm)

	// 主持靜態
	for i, picture := range hpics {
		if !strings.Contains(picture, "system") {
			hpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}
	// 玩家靜態
	for i, picture := range gpics {
		if !strings.Contains(picture, "system") {
			gpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}
	// 共用靜態
	for i, picture := range cpics {
		if !strings.Contains(picture, "system") {
			cpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}

	// 主持靜態
	for i, picture := range hanis {
		if !strings.Contains(picture, "system") {
			hanis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}
	// 玩家靜態
	for i, picture := range ganis {
		if !strings.Contains(picture, "system") {
			ganis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}
	// 共用靜態
	for i, picture := range canis {
		if !strings.Contains(picture, "system") {
			canis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + picture
		}
	}

	// 音樂
	for i, music := range musics {
		if !strings.Contains(music, "system") {
			musics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/general/" + music
		}
	}

	a.GeneralCustomizeHostPictures = hpics      // 主持靜態
	a.GeneralCustomizeGuestPictures = gpics     // 玩家靜態
	a.GeneralCustomizeCommonPictures = cpics    // 共用靜態
	a.GeneralCustomizeHostAnipictures = hanis   // 主持動態
	a.GeneralCustomizeGuestAnipictures = ganis  // 玩家動態
	a.GeneralCustomizeCommonAnipictures = canis // 共用動態
	a.GeneralCustomizeMusics = musics           // 音樂

	// 立體簽到
	// 圖片
	// threedAvatar, _ := m["threed_avatar"].(string)
	if !strings.Contains(a.ThreedAvatar, "system") {
		a.ThreedAvatar = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/threed/" + a.ThreedAvatar
	}

	// 圖片
	// threedBackground, _ := m["threed_background"].(string)
	if !strings.Contains(a.ThreedBackground, "system") {
		a.ThreedBackground = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/threed/" + a.ThreedBackground
	}

	// 圖片
	// threedImageLogoPicture, _ := m["threed_image_logo_picture"].(string)
	if !strings.Contains(a.ThreedImageLogoPicture, "system") {
		a.ThreedImageLogoPicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/threed/" + a.ThreedImageLogoPicture
	}

	// 自定義圖片參數(陣列)
	var (
		threedPictures    = make([]string, 0)
		threedAnipictures = make([]string, 0)
		threedMusics      = make([]string, 0)
	)

	// 自定義圖片
	if a.ThreedTopic == "01_classic" {
		threedPictures = append(threedPictures, a.ThreedBackground)
	}
	threedMusics = append(threedMusics, a.ThreedBgm)

	// 判斷是否為自定義上傳資料
	for i, picture := range threedPictures {
		if !strings.Contains(picture, "system") {
			threedPictures[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/threed/" + picture
		}
	}
	for i, music := range threedMusics {
		if !strings.Contains(music, "system") {
			threedMusics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/threed/" + music
		}
	}

	a.ThreedCustomizePictures = threedPictures
	a.ThreedCustomizeAnipictures = threedAnipictures
	a.ThreedCustomizeMusics = threedMusics

	// 自定義圖片參數(陣列)
	hpics = make([]string, 0) // 主持靜態
	gpics = make([]string, 0) // 玩家靜態
	cpics = make([]string, 0) // 共用靜態
	hanis = make([]string, 0) // 主持動態
	ganis = make([]string, 0) // 玩家動態
	canis = make([]string, 0) // 共用動態
	musics = make([]string, 0)

	// 自定義圖片
	if a.SignnameTopic == "01_classic" {
		hpics = append(hpics, a.SignnameClassicHPic01)
		hpics = append(hpics, a.SignnameClassicHPic02)
		hpics = append(hpics, a.SignnameClassicHPic03)
		hpics = append(hpics, a.SignnameClassicHPic04)
		hpics = append(hpics, a.SignnameClassicHPic05)
		hpics = append(hpics, a.SignnameClassicHPic06)
		hpics = append(hpics, a.SignnameClassicHPic07)
		gpics = append(cpics, a.SignnameClassicGPic01)
	}
	musics = append(musics, a.SignnameBgm)

	// 主持靜態
	for i, picture := range hpics {
		if !strings.Contains(picture, "system") {
			hpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}
	// 玩家靜態
	for i, picture := range gpics {
		if !strings.Contains(picture, "system") {
			gpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}
	// 共用靜態
	for i, picture := range cpics {
		if !strings.Contains(picture, "system") {
			cpics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}

	// 主持靜態
	for i, picture := range hanis {
		if !strings.Contains(picture, "system") {
			hanis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}
	// 玩家靜態
	for i, picture := range ganis {
		if !strings.Contains(picture, "system") {
			ganis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}
	// 共用靜態
	for i, picture := range canis {
		if !strings.Contains(picture, "system") {
			canis[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + picture
		}
	}

	// 音樂
	for i, music := range musics {
		if !strings.Contains(music, "system") {
			musics[i] = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/signname/" + music
		}
	}

	a.SignnameCustomizeHostPictures = hpics      // 主持靜態
	a.SignnameCustomizeGuestPictures = gpics     // 玩家靜態
	a.SignnameCustomizeCommonPictures = cpics    // 共用靜態
	a.SignnameCustomizeHostAnipictures = hanis   // 主持動態
	a.SignnameCustomizeGuestAnipictures = ganis  // 玩家動態
	a.SignnameCustomizeCommonAnipictures = canis // 共用動態
	a.SignnameCustomizeMusics = musics           // 音樂

	if a.InfoPicture != "" {
		a.InfoPicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/applysign/customize/" + a.InfoPicture
	}

	return a
}

// 判斷redis裡是否有簽到人員資訊(sign_staffs_2_activityID，SET)
// staffs2, err := a.RedisConn.SetGetMembers(config.SIGN_STAFFS_2_REDIS + activityID)
// if err != nil {
// 	return errors.New("錯誤: 從redis中取得簽到人員資訊發生問題(sign_staffs_2_activityID)")
// }

// 將用戶資料加入redis中(新增時處理)
// for _, userID := range users {
// #####威翔暫時沒用到#####
// 判斷redis裡是否有簽到人員資訊(sign_staffs_1_activityID)
// staffs1, err := a.RedisConn.ListRange(
// 	config.SIGN_STAFFS_1_REDIS+activityID, 0, 0)
// if err != nil {
// 	return errors.New("錯誤: 從redis中取得簽到人員資訊發生問題(sign_staffs_1_activityID)")
// }
// #####威翔暫時沒用到#####

// #####威翔暫時沒用到#####
// 判斷用戶是否已在redis中(sign_staffs_1_activityID)
// if len(staffs1) != 0 && !utils.InArray(staffs1, userID) {
// 	// redis中不存在簽到人員快取資訊，新增簽到人員資訊
// 	if err := a.RedisConn.ListRPush(config.SIGN_STAFFS_1_REDIS+activityID,
// 		userID); err != nil {
// 		return errors.New("錯誤: 新增簽到人員快取資料發生問題(sign_staffs_1_activityID)")
// 	}

// 	staffs1 = append(staffs1, userID)
// }
// #####威翔暫時沒用到#####

// 判斷redis裡是否有簽到人員資訊(sign_staffs_2_activityID，SET)
// staffs2, err := a.RedisConn.SetGetMembers(config.SIGN_STAFFS_2_REDIS + activityID)
// if err != nil {
// 	return errors.New("錯誤: 從redis中取得簽到人員資訊發生問題(sign_staffs_2_activityID)")
// }

// 	// 判斷用戶是否已在redis中(sign_staffs_2_activityID，SET)
// 	if len(staffs2) != 0 && !utils.InArray(staffs2, userID) {
// 		// redis中不存在簽到人員快取資訊，新增簽到人員資訊
// 		if err := a.RedisConn.SetAdd([]interface{}{config.SIGN_STAFFS_2_REDIS + activityID, userID}); err != nil {
// 			return errors.New("錯誤: 新增簽到人員快取資料發生問題(sign_staffs_2_activityID)")
// 		}

// 		staffs2 = append(staffs2, userID)
// 	}
// }

// #####威翔暫時沒用到#####
// 判斷redis裡是否有簽到人員資訊(sign_staffs_1_activityID)
// staffs1, err := a.RedisConn.ListRange(
// 	config.SIGN_STAFFS_1_REDIS+activityID, 0, 0)
// if err != nil {
// 	return errors.New("錯誤: 從redis中取得簽到人員資訊發生問題(sign_staffs_1_activityID)")
// }
// // 判斷用戶是否已在redis中(sign_staffs_1_activityID)
// if len(staffs1) != 0 && !utils.InArray(staffs1, userID) {
// 	// redis中不存在簽到人員快取資訊，新增簽到人員資訊
// 	if err := a.RedisConn.ListRPush(config.SIGN_STAFFS_1_REDIS+activityID,
// 		userID); err != nil {
// 		return errors.New("錯誤: 新增簽到人員快取資料發生問題(sign_staffs_1_activityID)")
// 	}
// }
// #####威翔暫時沒用到#####

// DecrAttend 遞減活動參加人數(修改活動人數快取資料)
// func (a ActivityModel) DecrAttend(isRedis bool, activityID, userID string) error {
// 	if err := a.Table(a.Base.TableName).
// 		WhereRaw("`activity_id` = ? and `attend` <= `people`", activityID).
// 		Update(command.Value{"attend": "attend - 1"}); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 減少活動人數發生錯誤, 請重新操作")
// 	}

// 	if isRedis {
// 		// 判斷redis裡是否有活動資訊，有則不執行查詢資料表功能
// 		// value, err := a.RedisConn.GetCache(config.ACTIVITY_REDIS + activityID)
// 		// if err != nil {
// 		// 	return errors.New("錯誤: 取得活動快取資料發生問題")
// 		// }
// 		// if value != "" {
// 		// 	// redis中已存在活動人數快取資訊，遞減人數資訊(有錯誤代表沒有快取資訊)
// 		// 	a.RedisConn.DecrCache(config.ACTIVITY_REDIS + activityID)
// 		// }

// 		// TODO: 當redis已存在簽到人員快取資訊時，暫時不刪除取消簽到人員的資料(威翔判斷順序會發生錯誤)
// 		// 判斷redis裡是否有簽到人員資訊
// 		// if count := a.RedisConn.ListLen(config.SIGN_STAFFS_REDIS + activityID); count != 0 {
// 		// redis中已存在簽到人員快取資訊，新增簽到人員資訊
// 		// if err = a.RedisConn.SetRem([]interface{}{config.SIGN_STAFFS_REDIS + activityID,
// 		// 	userID}); err != nil {
// 		// 	return errors.New("錯誤: 移除簽到人員快取資料發生問題")
// 		// }
// 		// }

// 		// fmt.Println("執行刪除")
// 		// 移除redis中簽到人員資訊(sign_staffs_2_activityID，SET)
// 		if err := a.RedisConn.SetRem([]interface{}{config.SIGN_STAFFS_2_REDIS + activityID, userID}); err != nil {
// 			return errors.New("錯誤: 移除簽到人員快取資料發生問題(sign_staffs_2_activityID)")
// 		}

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 		// a.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+activityID, "DecrAttend")
// 	}

// 	return nil
// }

// 判斷redis裡是否有簽到人員資訊(sign_staffs_2_activityID，SET)
// staffs2, err := a.RedisConn.SetGetMembers(config.SIGN_STAFFS_2_REDIS + activityID)
// if err != nil {
// 	return errors.New("錯誤: 從redis中取得簽到人員資訊發生問題(sign_staffs_2_activityID)")
// }

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+activityID, "IncrAtte

// activity.ID, _ = items[i]["id"].(int64)
// activity.ActivityID, _ = items[i]["activity_id"].(string)
// activity.UserID, _ = items[i]["user_id"].(string)
// activity.ActivityName, _ = items[i]["activity_name"].(string)
// activity.ActivityType, _ = items[i]["activity_type"].(string)
// activity.MaxPeople, _ = items[i]["max_people"].(int64)
// activity.People, _ = items[i]["people"].(int64)
// activity.Attend, _ = items[i]["attend"].(int64)
// activity.City, _ = items[i]["city"].(string)
// activity.Town, _ = items[i]["town"].(string)
// activity.StartTime, _ = items[i]["start_time"].(string)
// activity.EndTime, _ = items[i]["end_time"].(string)
// activity.Number, _ = items[i]["number"].(int64)
// activity.LoginRequired, _ = items[i]["login_required"].(string)
// activity.LoginPassword, _ = items[i]["login_password"].(string)
// activity.PasswordRequired, _ = items[i]["password_required"].(string)
// activity.MessageAmount, _ = items[i]["message_amount"].(int64)
// activity.PushPhoneMessage, _ = items[i]["push_phone_message"].(string)

// activity.SendMail, _ = items[i]["send_mail"].(string)
// activity.MailAmount, _ = items[i]["mail_amount"].(int64)

// activity.HostScan, _ = items[i]["host_scan"].(string)
// activity.CustomizePassword, _ = items[i]["customize_password"].(string)
// activity.AllowCustomizeApply, _ = items[i]["allow_customize_apply"].(string)

// 用戶官方帳號綁定
// activity.UserLineID, _ = items[i]["user_line_id"].(string)
// activity.UserChannelID, _ = items[i]["user_channel_id"].(string)
// activity.UserChannelSecret, _ = items[i]["user_channel_secret"].(string)
// activity.UserChatbotSecret, _ = items[i]["user_chatbot_secret"].(string)
// activity.UserChatbotToken, _ = items[i]["user_chatbot_token"].(string)
// 活動官方帳號綁定
// activity.ActivityLineID, _ = items[i]["activity_line_id"].(string)
// activity.ActivityChannelID, _ = items[i]["activity_channel_id"].(string)
// activity.ActivityChannelSecret, _ = items[i]["activity_channel_secret"].(string)
// activity.ActivityChatbotSecret, _ = items[i]["activity_chatbot_secret"].(string)
// activity.ActivityChatbotToken, _ = items[i]["activity_chatbot_token"].(string)

// 官方帳號傳送訊息
// activity.PushMessage, _ = items[i]["push_message"].(string)

// 驗證裝置
// activity.Device, _ = items[i]["device"].(string)

// 活動總覽
// activity.OverviewMessage, _ = items[i]["overview_message"].(string)
// activity.OverviewTopic, _ = items[i]["overview_topic"].(string)
// activity.OverviewQuestion, _ = items[i]["overview_question"].(string)
// activity.OverviewDanmu, _ = items[i]["overview_danmu"].(string)
// activity.OverviewSpecialDanmu, _ = items[i]["overview_special_danmu"].(string)
// activity.OverviewPicture, _ = items[i]["overview_picture"].(string)
// activity.OverviewHoldscreen, _ = items[i]["overview_holdscreen"].(string)
// activity.OverviewGeneral, _ = items[i]["overview_general"].(string)
// activity.OverviewThreed, _ = items[i]["overview_threed"].(string)
// activity.OverviewCountdown, _ = items[i]["overview_countdown"].(string)
// activity.OverviewLottery, _ = items[i]["overview_lottery"].(string)
// activity.OverviewRedpack, _ = items[i]["overview_redpack"].(string)
// activity.OverviewRopepack, _ = items[i]["overview_ropepack"].(string)
// activity.OverviewWhackMole, _ = items[i]["overview_whack_mole"].(string)
// activity.OverviewDrawNumbers, _ = items[i]["overview_draw_numbers"].(string)
// activity.OverviewMonopoly, _ = items[i]["overview_monopoly"].(string)
// activity.OverviewQA, _ = items[i]["overview_qa"].(string)
// activity.OverviewTugofwar, _ = items[i]["overview_tugofwar"].(string)
// activity.OverviewBingo, _ = items[i]["overview_bingo"].(string)
// activity.OverviewSignname, _ = items[i]["overview_signname"].(string)
// activity.Overview3DGachaMachine, _ = items[i]["overview_3d_gacha_machine"].(string)
// activity.OverviewVote, _ = items[i]["overview_vote"].(string)
// activity.OverviewChatroom, _ = items[i]["overview_chatroom"].(string)
// activity.OverviewInteract, _ = items[i]["overview_interact"].(string)
// activity.OverviewInfo, _ = items[i]["overview_info"].(string)

// 活動介紹
// activity.IntroduceTitle, _ = items[i]["introduce_title"].(string)
// // 活動行程
// activity.ScheduleTitle, _ = items[i]["schedule_title"].(string)
// activity.ScheduleDisplayDate, _ = items[i]["schedule_display_date"].(string)
// activity.ScheduleDisplayDetail, _ = items[i]["schedule_display_detail"].(string)
// // 活動嘉賓
// activity.GuestTitle, _ = items[i]["guest_title"].(string)
// activity.GuestBackground, _ = items[i]["guest_background"].(string)
// // 活動資料
// activity.MaterialTitle, _ = items[i]["material_title"].(string)
// // 報名
// activity.ApplyCheck, _ = items[i]["apply_check"].(string)
// activity.CustomizeDefaultAvatar, _ = items[i]["customize_default_avatar"].(string)

// 簽到
// activity.SignCheck, _ = items[i]["sign_check"].(string)
// activity.SignAllow, _ = items[i]["sign_allow"].(string)
// activity.SignMinutes, _ = items[i]["sign_minutes"].(int64)
// activity.SignManual, _ = items[i]["sign_manual"].(string)

// activity.QRcodeLogoSize, _ = items[i]["qrcode_logo_size"].(float64)
// activity.QRcodePicturePoint, _ = items[i]["qrcode_picture_point"].(string)
// activity.QRcodeWhiteDistance, _ = items[i]["qrcode_white_distance"].(int64)
// activity.QRcodePointColor, _ = items[i]["qrcode_point_color"].(string)
// activity.QRcodeBackgroundColor, _ = items[i]["qrcode_background_color"].(string)

// 訊息牆
// activity.MessagePicture, _ = items[i]["message_picture"].(string)
// activity.MessageAuto, _ = items[i]["message_auto"].(string)
// activity.MessageBan, _ = items[i]["message_ban"].(string)
// activity.MessageBanSecond, _ = items[i]["message_ban_second"].(int64)
// activity.MessageRefreshSecond, _ = items[i]["message_refresh_second"].(int64)
// activity.MessageOpen, _ = items[i]["message_open"].(string)

// activity.MessageTextColor, _ = items[i]["message_text_color"].(string)
// activity.MessageScreenTitleColor, _ = items[i]["message_screen_title_color"].(string)
// activity.MessageTextFrameColor, _ = items[i]["message_text_frame_color"].(string)
// activity.MessageFrameColor, _ = items[i]["message_frame_color"].(string)
// activity.MessageScreenTitleFrameColor, _ = items[i]["message_screen_title_frame_color"].(string)

// 提問牆
// activity.QuestionMessageCheck, _ = items[i]["question_message_check"].(string)
// activity.QuestionAnonymous, _ = items[i]["question_anonymous"].(string)
// activity.QuestionHideAnswered, _ = items[i]["question_hide_answered"].(string)
// activity.QuestionQrcode, _ = items[i]["question_qrcode"].(string)

// activity.QuestionGuestAnswered, _ = items[i]["question_guest_answered"].(string)

// 彈幕
// activity.DanmuLoop, _ = items[i]["danmu_loop"].(string)
// activity.DanmuTop, _ = items[i]["danmu_top"].(string)
// activity.DanmuMid, _ = items[i]["danmu_mid"].(string)
// activity.DanmuBottom, _ = items[i]["danmu_bottom"].(string)
// activity.DanmuDisplayName, _ = items[i]["danmu_display_name"].(string)
// activity.DanmuDisplayAvatar, _ = items[i]["danmu_display_avatar"].(string)
// activity.DanmuSize, _ = items[i]["danmu_size"].(float64)
// activity.DanmuSpeed, _ = items[i]["danmu_speed"].(float64)
// activity.DanmuDensity, _ = items[i]["danmu_density"].(float64)
// activity.DanmuOpacity, _ = items[i]["danmu_opacity"].(float64)
// activity.DanmuBackground, _ = items[i]["danmu_background"].(int64)
// activity.DanmuNewBackgroundColor, _ = items[i]["danmu_new_background_color"].(string)
// activity.DanmuNewTextColor, _ = items[i]["danmu_new_text_color"].(string)
// activity.DanmuOldBackgroundColor, _ = items[i]["danmu_old_background_color"].(string)
// activity.DanmuOldTextColor, _ = items[i]["danmu_old_text_color"].(string)

// activity.DanmuStyle, _ = items[i]["danmu_style"].(string)

// 特殊彈幕
// activity.SpecialDanmuMessageCheck, _ = items[i]["special_danmu_message_check"].(string)
// activity.SpecialDanmuGeneralPrice, _ = items[i]["special_danmu_general_price"].(int64)
// activity.SpecialDanmuLargePrice, _ = items[i]["special_danmu_large_price"].(int64)
// activity.SpecialDanmuOverPrice, _ = items[i]["special_danmu_over_price"].(int64)
// activity.SpecialDanmuTopic, _ = items[i]["special_danmu_topic"].(string)
// activity.SpecialDanmuBanSecond, _ = items[i]["special_danmu_ban_second"].(int64)

// activity.PictureStartTime, _ = items[i]["picture_start_time"].(string)
// activity.PictureEndTime, _ = items[i]["picture_end_time"].(string)
// activity.PictureHideTime, _ = items[i]["picture_hide_time"].(string)
// activity.PictureSwitchSecond, _ = items[i]["picture_switch_second"].(int64)
// activity.PicturePlayOrder, _ = items[i]["picture_play_order"].(string)
// activity.PictureBackground, _ = items[i]["picture_background"].(string)
// activity.PictureAnimation, _ = items[i]["picture_animation"].(int64)

// 霸屏
// activity.HoldscreenPrice, _ = items[i]["holdscreen_price"].(int64)
// activity.HoldscreenMessageCheck, _ = items[i]["holdscreen_message_check"].(string)
// activity.HoldscreenOnlyPicture, _ = items[i]["holdscreen_only_picture"].(string)
// activity.HoldscreenDisplaySecond, _ = items[i]["holdscreen_display_second"].(int64)
// activity.HoldscreenBirthdayTopic, _ = items[i]["holdscreen_birthday_topic"].(string)
// activity.HoldscreenBanSecond, _ = items[i]["holdscreen_ban_second"].(int64)

// activity.GeneralDisplayPeople, _ = items[i]["general_display_people"].(string)
// activity.GeneralStyle, _ = items[i]["general_style"].(int64)
// activity.GeneralTopic, _ = items[i]["general_topic"].(string)
// activity.GeneralMusic, _ = items[i]["general_music"].(string)
// activity.GeneralLoop, _ = items[i]["general_loop"].(string)
// activity.GeneralLatest, _ = items[i]["general_latest"].(string)
// activity.GeneralEditTimes, _ = items[i]["general_edit_times"].(int64)

// activity.GeneralClassicHPic01, _ = items[i]["general_classic_h_pic_01"].(string)
// activity.GeneralClassicHPic02, _ = items[i]["general_classic_h_pic_02"].(string)
// activity.GeneralClassicHPic03, _ = items[i]["general_classic_h_pic_03"].(string)
// activity.GeneralClassicHAni01, _ = items[i]["general_classic_h_ani_01"].(string)

// activity.GeneralBgm, _ = items[i]["general_bgm"].(string)

// activity.ThreedAvatarShape, _ = items[i]["threed_avatar_shape"].(string)
// activity.ThreedDisplayPeople, _ = items[i]["threed_display_people"].(string)
// activity.ThreedBackgroundStyle, _ = items[i]["threed_background_style"].(int64)

// 倒數計時
// activity.CountdownSecond, _ = items[i]["countdown_second"].(int64)
// activity.CountdownURL, _ = items[i]["countdown_url"].(string)
// activity.CountdownAvatar, _ = items[i]["countdown_avatar"].(string)
// activity.CountdownAvatarShape, _ = items[i]["countdown_avatar_shape"].(string)
// activity.CountdownBackground, _ = items[i]["countdown_background"].(int64)

// activity.ThreedImageLogo, _ = items[i]["threed_image_logo"].(string)
// activity.ThreedImageCircle, _ = items[i]["threed_image_circle"].(string)
// activity.ThreedImageSpiral, _ = items[i]["threed_image_spiral"].(string)
// activity.ThreedImageRectangle, _ = items[i]["threed_image_rectangle"].(string)
// activity.ThreedImageSquare, _ = items[i]["threed_image_square"].(string)

// activity.ThreedTopic, _ = items[i]["threed_topic"].(string)
// activity.ThreedMusic, _ = items[i]["threed_music"].(string)
// activity.ThreedEditTimes, _ = items[i]["threed_edit_times"].(int64)

// activity.ThreedBgm, _ = items[i]["threed_bgm"].(string)

// 簽名牆
// activity.SignnameMode, _ = items[i]["signname_mode"].(string)
// activity.SignnameTimes, _ = items[i]["signname_times"].(int64)
// activity.SignnameDisplay, _ = items[i]["signname_display"].(string)
// activity.SignnameLimitTimes, _ = items[i]["signname_limit_times"].(string)
// activity.SignnameTopic, _ = items[i]["signname_topic"].(string)
// activity.SignnameMusic, _ = items[i]["signname_music"].(string)
// activity.SignnameLoop, _ = items[i]["signname_loop"].(string)
// activity.SignnameLatest, _ = items[i]["signname_latest"].(string)
// activity.SignnameContent, _ = items[i]["signname_content"].(string)
// activity.SignnameShowName, _ = items[i]["signname_show_name"].(string)
// activity.SignnameEditTimes, _ = items[i]["signname_edit_times"].(int64)

// activity.SignnameClassicHPic01, _ = items[i]["signname_classic_h_pic_01"].(string)
// activity.SignnameClassicHPic02, _ = items[i]["signname_classic_h_pic_02"].(string)
// activity.SignnameClassicHPic03, _ = items[i]["signname_classic_h_pic_03"].(string)
// activity.SignnameClassicHPic04, _ = items[i]["signname_classic_h_pic_04"].(string)
// activity.SignnameClassicHPic05, _ = items[i]["signname_classic_h_pic_05"].(string)
// activity.SignnameClassicHPic06, _ = items[i]["signname_classic_h_pic_06"].(string)
// activity.SignnameClassicHPic07, _ = items[i]["signname_classic_h_pic_07"].(string)
// activity.SignnameClassicGPic01, _ = items[i]["signname_classic_g_pic_01"].(string)
// activity.SignnameClassicCPic01, _ = items[i]["signname_classic_c_pic_01"].(string)

// activity.SignnameBgm, _ = items[i]["signname_bgm"].(string)

// 自定義
// activity.Ext1Name, _ = items[i]["ext_1_name"].(string)
// activity.Ext1Type, _ = items[i]["ext_1_type"].(string)
// activity.Ext1Options, _ = items[i]["ext_1_options"].(string)
// activity.Ext1Required, _ = items[i]["ext_1_required"].(string)

// activity.Ext2Name, _ = items[i]["ext_2_name"].(string)
// activity.Ext2Type, _ = items[i]["ext_2_type"].(string)
// activity.Ext2Options, _ = items[i]["ext_2_options"].(string)
// activity.Ext2Required, _ = items[i]["ext_2_required"].(string)

// activity.Ext3Name, _ = items[i]["ext_3_name"].(string)
// activity.Ext3Type, _ = items[i]["ext_3_type"].(string)
// activity.Ext3Options, _ = items[i]["ext_3_options"].(string)
// activity.Ext3Required, _ = items[i]["ext_3_required"].(string)

// activity.Ext4Name, _ = items[i]["ext_4_name"].(string)
// activity.Ext4Type, _ = items[i]["ext_4_type"].(string)
// activity.Ext4Options, _ = items[i]["ext_4_options"].(string)
// activity.Ext4Required, _ = items[i]["ext_4_required"].(string)

// activity.Ext5Name, _ = items[i]["ext_5_name"].(string)
// activity.Ext5Type, _ = items[i]["ext_5_type"].(string)
// activity.Ext5Options, _ = items[i]["ext_5_options"].(string)
// activity.Ext5Required, _ = items[i]["ext_5_required"].(string)

// activity.Ext6Name, _ = items[i]["ext_6_name"].(string)
// activity.Ext6Type, _ = items[i]["ext_6_type"].(string)
// activity.Ext6Options, _ = items[i]["ext_6_options"].(string)
// activity.Ext6Required, _ = items[i]["ext_6_required"].(string)

// activity.Ext7Name, _ = items[i]["ext_7_name"].(string)
// activity.Ext7Type, _ = items[i]["ext_7_type"].(string)
// activity.Ext7Options, _ = items[i]["ext_7_options"].(string)
// activity.Ext7Required, _ = items[i]["ext_7_required"].(string)

// activity.Ext8Name, _ = items[i]["ext_8_name"].(string)
// activity.Ext8Type, _ = items[i]["ext_8_type"].(string)
// activity.Ext8Options, _ = items[i]["ext_8_options"].(string)
// activity.Ext8Required, _ = items[i]["ext_8_required"].(string)

// activity.Ext9Name, _ = items[i]["ext_9_name"].(string)
// activity.Ext9Type, _ = items[i]["ext_9_type"].(string)
// activity.Ext9Options, _ = items[i]["ext_9_options"].(string)
// activity.Ext9Required, _ = items[i]["ext_9_required"].(string)

// activity.Ext10Name, _ = items[i]["ext_10_name"].(string)
// activity.Ext10Type, _ = items[i]["ext_10_type"].(string)
// activity.Ext10Options, _ = items[i]["ext_10_options"].(string)
// activity.Ext10Required, _ = items[i]["ext_10_required"].(string)

// activity.ExtEmailRequired, _ = items[i]["ext_email_required"].(string)
// activity.ExtPhoneRequired, _ = items[i]["ext_phone_required"].(string)
// activity.InfoPicture, _ = items[i]["info_picture"].(string)

// 訊息審核
// activity.MessageCheckManualCheck, _ = items[i]["message_check_manual_check"].(string)
// activity.MessageCheckSensitivity, _ = items[i]["message_check_sensitivity"].(string)

// 用戶
// activity.Name, _ = items[i]["name"].(string)
// activity.Avatar, _ = items[i]["avatar"].(string)
// activity.MaxActivityPeople, _ = items[i]["max_activity_people"].(int64)
// activity.MaxGamePeople, _ = items[i]["max_game_people"].(int64)
