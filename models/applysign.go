package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"

	"github.com/xuri/excelize/v2"
)

// ApplysignModel 資料表欄位
type ApplysignModel struct {
	Base `json:"-"`
	ID   int64 `json:"id"`
	// JoinID     int64  `json:"join_id"`
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Status     string `json:"status"`
	Role       string `json:"role"`
	Number     int64  `json:"number"`
	ApplyTime  string `json:"apply_time"`
	ReviewTime string `json:"review_time"`
	SignTime   string `json:"sign_time"`
	Ext1       string `json:"ext_1"`
	Ext2       string `json:"ext_2"`
	Ext3       string `json:"ext_3"`
	Ext4       string `json:"ext_4"`
	Ext5       string `json:"ext_5"`
	Ext6       string `json:"ext_6"`
	Ext7       string `json:"ext_7"`
	Ext8       string `json:"ext_8"`
	Ext9       string `json:"ext_9"`
	Ext10      string `json:"ext_10"`

	// 過濾參數
	AllPeople    int64 `json:"all_people" example:"100"`    // 所有人數
	NoPeople     int64 `json:"no_people" example:"100"`     // 未報名人數
	ReviewPeople int64 `json:"review_people" example:"100"` // 審核人數
	ApplyPeople  int64 `json:"apply_people" example:"100"`  // 報名人數
	RefusePeople int64 `json:"refuse_people" example:"100"` // 拒絕人數
	SignPeople   int64 `json:"sign_people" example:"100"`   // 簽到人數

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

	// 黑名單原因
	Reason string `json:"reason"`

	// 活動資訊(join activity)
	ActivityUserID          string `json:"activity_user_id"`
	ActivityName            string `json:"activity_name"`
	StartTime               string `json:"start_time"`
	EndTime                 string `json:"end_time"`
	SignCheck               string `json:"sign_check"`
	SignAllow               string `json:"sign_allow"`
	SignMinutes             int64  `json:"sign_minutes"`
	MessageBan              string `json:"message_ban"`
	MessageBanSecond        int64  `json:"message_ban_second"`
	HoldscreenBirthdayTopic string `json:"holdscreen_birthday_topic"`
	ActivityDevice          string `json:"activity_device"` // 活動驗證裝置

	// 活動總覽
	OverviewMessage string `json:"overview_message"`
	// OverviewTopic        string `json:"overview_topic"`
	OverviewQuestion     string `json:"overview_question"`
	OverviewDanmu        string `json:"overview_danmu"`
	OverviewSpecialDanmu string `json:"overview_special_danmu"`
	// OverviewPicture      string `json:"overview_picture"`
	OverviewHoldscreen string `json:"overview_holdscreen"`
	// OverviewGeneral      string `json:"overview_general"`
	// OverviewThreed       string `json:"overview_threed"`
	// OverviewCountdown    string `json:"overview_countdown"`
	OverviewLottery   string `json:"overview_lottery"`
	OverviewRedpack   string `json:"overview_redpack"`
	OverviewRopepack  string `json:"overview_ropepack"`
	OverviewWhackMole string `json:"overview_whack_mole"`
	// OverviewDrawNumbers  string `json:"overview_draw_numbers"`
	OverviewMonopoly       string `json:"overview_monopoly"`
	OverviewQA             string `json:"overview_qa"`
	OverviewTugofwar       string `json:"overview_tugofwar"`
	OverviewBingo          string `json:"overview_bingo"`
	OverviewSignname       string `json:"overview_signname"`
	Overview3DGachaMachine string `json:"overview_3d_gacha_machine"`
	OverviewVote           string `json:"overview_vote"`

	// 用戶資訊(join line_users)
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ExtEmail    string `json:"ext_email"`
	Friend      string `json:"friend"`
	Identify    string `json:"identify"`
	UserDevice  string `json:"user_device"`  // 用戶驗證裝置
	ExtPassword string `json:"ext_password"` // 驗證碼
	AdminID     string `json:"admin_id"`
	// LineID        string `json:"line_id"`        // line_id
	// ChannelID     string `json:"channel_id"`     // channel_id
	// ChannelSecret string `json:"channel_secret"` // channel_secret
	// ChatbotSecret string `json:"chatbot_secret"` // chatbot_secret
	// ChatbotToken  string `json:"chatbot_token"`  // chatbot_token

	// 自定義(join自定義表)
	Ext1Name     string `json:"ext_1_name"`
	Ext1Type     string `json:"ext_1_type"`
	Ext1Options  string `json:"ext_1_options"`
	Ext1Required string `json:"ext_1_required"`
	Ext1Unique   string `json:"ext_1_unique"`

	Ext2Name     string `json:"ext_2_name"`
	Ext2Type     string `json:"ext_2_type"`
	Ext2Options  string `json:"ext_2_options"`
	Ext2Required string `json:"ext_2_required"`
	Ext2Unique   string `json:"ext_2_unique"`

	Ext3Name     string `json:"ext_3_name"`
	Ext3Type     string `json:"ext_3_type"`
	Ext3Options  string `json:"ext_3_options"`
	Ext3Required string `json:"ext_3_required"`
	Ext3Unique   string `json:"ext_3_unique"`

	Ext4Name     string `json:"ext_4_name"`
	Ext4Type     string `json:"ext_4_type"`
	Ext4Options  string `json:"ext_4_options"`
	Ext4Required string `json:"ext_4_required"`
	Ext4Unique   string `json:"ext_4_unique"`

	Ext5Name     string `json:"ext_5_name"`
	Ext5Type     string `json:"ext_5_type"`
	Ext5Options  string `json:"ext_5_options"`
	Ext5Required string `json:"ext_5_required"`
	Ext5Unique   string `json:"ext_5_unique"`

	Ext6Name     string `json:"ext_6_name"`
	Ext6Type     string `json:"ext_6_type"`
	Ext6Options  string `json:"ext_6_options"`
	Ext6Required string `json:"ext_6_required"`
	Ext6Unique   string `json:"ext_6_unique"`

	Ext7Name     string `json:"ext_7_name"`
	Ext7Type     string `json:"ext_7_type"`
	Ext7Options  string `json:"ext_7_options"`
	Ext7Required string `json:"ext_7_required"`
	Ext7Unique   string `json:"ext_7_unique"`

	Ext8Name     string `json:"ext_8_name"`
	Ext8Type     string `json:"ext_8_type"`
	Ext8Options  string `json:"ext_8_options"`
	Ext8Required string `json:"ext_8_required"`
	Ext8Unique   string `json:"ext_8_unique"`

	Ext9Name     string `json:"ext_9_name"`
	Ext9Type     string `json:"ext_9_type"`
	Ext9Options  string `json:"ext_9_options"`
	Ext9Required string `json:"ext_9_required"`
	Ext9Unique   string `json:"ext_9_unique"`

	Ext10Name     string `json:"ext_10_name"`
	Ext10Type     string `json:"ext_10_type"`
	Ext10Options  string `json:"ext_10_options"`
	Ext10Required string `json:"ext_10_required"`
	Ext10Unique   string `json:"ext_10_unique"`

	ExtEmailRequired string `json:"ext_email_required"`
	ExtPhoneRequired string `json:"ext_phone_required"`
	InfoPicture      string `json:"info_picture"`
}

// NewApplysignModel 資料表欄位
type NewApplysignModel struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Status     string `json:"status"`
}

// EditApplysignModel 資料表欄位
type EditApplysignModel struct {
	LineUsers  []string `json:"line_users"` // 批量更新用戶
	ActivityID string   `json:"activity_id"`
	Status     string   `json:"status"`
	Role       string   `json:"role"`
	ReviewTime string   `json:"review_time"`
	SignTime   string   `json:"sign_time"`
}

// EditApplyModel 資料表欄位
type EditApplyModel struct {
	ActivityID             string `json:"activity_id"`
	ApplyCheck             string `json:"apply_check"`
	CustomizePassword      string `json:"customize_password"`       // 是否自定義設置驗證碼
	AllowCustomizeApply    string `json:"allow_customize_apply"`    // 用戶是否允許自定義報名
	CustomizeDefaultAvatar string `json:"customize_default_avatar"` // 自定義人員預設頭像
}

// EditSignModel 資料表欄位
type EditSignModel struct {
	ActivityID  string `json:"activity_id"`
	SignCheck   string `json:"sign_check"`
	SignAllow   string `json:"sign_allow"`
	SignMinutes string `json:"sign_minutes"`
	SignManual  string `json:"sign_manual"`
}

// DefaultApplysignModel 預設ApplysignModel
func DefaultApplysignModel() ApplysignModel {
	return ApplysignModel{Base: Base{TableName: config.ACTIVITY_APPLYSIGN_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (s ApplysignModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ApplysignModel {
	s.DbConn = dbconn
	s.RedisConn = cacheconn
	s.MongoConn = mongoconn
	return s
}

// SetDbConn 設定connection
// func (s ApplysignModel) SetDbConn(conn db.Connection) ApplysignModel {
// 	s.DbConn = conn
// 	return s
// }

// // SetRedisConn 設定connection
// func (s ApplysignModel) SetRedisConn(conn cache.Connection) ApplysignModel {
// 	s.RedisConn = conn
// 	return s
// }

// FindHostApplysignAmount 查詢該活動不同狀態的人員數量(主持端)
func (s ApplysignModel) FindHostApplysignAmount(activityID string) (ApplysignModel, error) {
	var (
		sql = s.Table(s.Base.TableName).
			Select(
				"COUNT(*) AS all_people",
				"SUM(CASE WHEN status = 'sign' THEN 1 ELSE 0 END) AS sign_people",
				"SUM(CASE WHEN status = 'apply' THEN 1 ELSE 0 END) AS apply_people",
				"SUM(CASE WHEN status = 'review' THEN 1 ELSE 0 END) AS review_people",
				"SUM(CASE WHEN status = 'refuse' THEN 1 ELSE 0 END) AS refuse_people",
				"SUM(CASE WHEN status = 'no' THEN 1 ELSE 0 END) AS no_people",
			)
	)

	if activityID != "" {
		sql = sql.Where("activity_applysign.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return s, errors.New("錯誤: 無法取得該活動報名簽到人員數量，請重新查詢")
	}

	// 回傳的資料為[]byte格式，需轉換為int64
	for key, value := range item {
		item[key] = utils.GetInt64(value, 0)
	}

	s = s.MapToModel(item)

	return s, nil
}

// FindApplysignAmount 查詢該活動報名簽到人員數量
func (s ApplysignModel) FindApplysignAmount(activityID, name, status string) (int64, error) {
	var (
		sql = s.Table(s.Base.TableName).
			Select("count(*)").
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="})
	)

	if activityID != "" {
		sql = sql.Where("activity_applysign.activity_id", "=", activityID)
	}
	if name != "" {
		sql = sql.WhereIn("line_users.name", interfaces(strings.Split(name, ",")))
	}
	if status != "" {
		sql = sql.WhereIn("activity_applysign.status", interfaces(strings.Split(status, ",")))
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法該活動報名簽到人員數量，請重新查詢")
	}

	return item["count(*)"].(int64), nil
}

// FindAll 查詢所有報名簽到人員資料
func (s ApplysignModel) FindAll(activityID, userID,
	name, status string,
	limit, offset int64) ([]ApplysignModel, error) {
	var (
		sql = s.Table(s.Base.TableName).
			Select("activity_applysign.id", "activity_applysign.user_id",
				"activity_applysign.activity_id", "activity_applysign.number",
				"activity_applysign.status", "activity_applysign.apply_time",
				"activity_applysign.review_time", "activity_applysign.sign_time",
				"activity_applysign.role",

				// 自定義欄位
				"activity_applysign.ext_1", "activity_customize.ext_1_name", "activity_customize.ext_1_type",
				"activity_customize.ext_1_options", "activity_customize.ext_1_required",
				"activity_applysign.ext_2", "activity_customize.ext_2_name", "activity_customize.ext_2_type",
				"activity_customize.ext_2_options", "activity_customize.ext_2_required",
				"activity_applysign.ext_3", "activity_customize.ext_3_name", "activity_customize.ext_3_type",
				"activity_customize.ext_3_options", "activity_customize.ext_3_required",
				"activity_applysign.ext_4", "activity_customize.ext_4_name", "activity_customize.ext_4_type",
				"activity_customize.ext_4_options", "activity_customize.ext_4_required",
				"activity_applysign.ext_5", "activity_customize.ext_5_name", "activity_customize.ext_5_type",
				"activity_customize.ext_5_options", "activity_customize.ext_5_required",
				"activity_applysign.ext_6", "activity_customize.ext_6_name", "activity_customize.ext_6_type",
				"activity_customize.ext_6_options", "activity_customize.ext_6_required",
				"activity_applysign.ext_7", "activity_customize.ext_7_name", "activity_customize.ext_7_type",
				"activity_customize.ext_7_options", "activity_customize.ext_7_required",
				"activity_applysign.ext_8", "activity_customize.ext_8_name", "activity_customize.ext_8_type",
				"activity_customize.ext_8_options", "activity_customize.ext_8_required",
				"activity_applysign.ext_9", "activity_customize.ext_9_name", "activity_customize.ext_9_type",
				"activity_customize.ext_9_options", "activity_customize.ext_9_required",
				"activity_applysign.ext_10", "activity_customize.ext_10_name", "activity_customize.ext_10_type",
				"activity_customize.ext_10_options", "activity_customize.ext_10_required",

				"activity_customize.ext_1_unique", "activity_customize.ext_2_unique", "activity_customize.ext_3_unique", "activity_customize.ext_4_unique",
				"activity_customize.ext_5_unique", "activity_customize.ext_6_unique", "activity_customize.ext_7_unique", "activity_customize.ext_8_unique",
				"activity_customize.ext_9_unique", "activity_customize.ext_10_unique",

				"activity_customize.ext_email_required", "activity_customize.ext_phone_required",
				"activity_customize.info_picture",

				// 用戶資訊
				"line_users.name", "line_users.avatar", "line_users.phone",
				"line_users.email", "line_users.ext_email", "line_users.device as user_device",
				"line_users.ext_password", "line_users.admin_id",

				// 活動資訊
				"activity.user_id as activity_user_id", "activity.activity_name",
				"activity.start_time", "activity.end_time",
				"activity.device as activity_device").
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "=",
			}).
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "=",
			}).
			OrderBy(
				"field", "(activity_applysign.status, 'sign', 'apply', 'review', 'cancel', 'refuse', 'no')",
			)
	)

	// 判斷參數是否為空
	if userID != "" {
		sql = sql.Where("activity_applysign.user_id", "=", userID)
	}
	if activityID != "" {
		sql = sql.Where("activity_applysign.activity_id", "=", activityID)
	}
	if name != "" {
		sql = sql.WhereIn("line_users.name", interfaces(strings.Split(name, ",")))
	}
	if status != "" {
		sql = sql.WhereIn("activity_applysign.status", interfaces(strings.Split(status, ",")))
	}
	if limit != 0 {
		sql = sql.Limit(limit)
	}
	if offset != 0 {
		sql = sql.Offset(offset)
	}

	items, err := sql.All()
	if err != nil {
		return []ApplysignModel{}, errors.New("錯誤: 無法取得報名簽到人員資訊，請重新查詢")
	}

	return MapToApplysignModel(items), nil
}

// Export 匯出報名簽到人員資料
func (s ApplysignModel) Export(activityID, status string) (*excelize.File, error) {
	var (
		file                 = excelize.NewFile()      // 開啟EXCEL檔案
		index, _             = file.NewSheet("Sheet1") // 創建一個工作表
		customizeModel, err1 = DefaultCustomizeModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					Find(activityID) // 自定義欄位
		rowNames = []string{"A", "B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
			"S", "T", "U", "V", "W", "X", "Y", "Z"}
		colNames = []string{
			// 用戶資訊
			"用戶名稱", "抽獎號碼", "驗證裝置", "驗證碼",

			// 報名簽到資訊
			"報名簽到狀態", "報名時間", "審核時間", "簽到時間",
		}

		// 自定義欄位是否必填
		extRequireds = []string{
			// 電話.信箱是否必填
			customizeModel.ExtPhoneRequired, customizeModel.ExtEmailRequired,
			// 自定義欄位
			customizeModel.Ext1Name, customizeModel.Ext2Name, customizeModel.Ext3Name,
			customizeModel.Ext4Name, customizeModel.Ext5Name, customizeModel.Ext6Name,
			customizeModel.Ext7Name, customizeModel.Ext8Name, customizeModel.Ext9Name,
			customizeModel.Ext10Name,
		}
		applySigns, err2 = DefaultApplysignModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					FindAll(activityID, "", "", status, 0, 0)
		sheet    = "Sheet1"
		statusCN string
	)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err1 != nil || err2 != nil {
		if err1 != nil {
			return file, err1
		} else if err2 != nil {
			return file, err2
		}
	}

	// 設置活頁簿的默認工作表
	file.SetActiveSheet(index)

	// 判斷自定義欄位
	for i, extRequired := range extRequireds {
		if i == 0 {
			if extRequired == "true" {
				colNames = append(colNames, "電話")
			}
		} else if i == 1 {
			if extRequired == "true" {
				colNames = append(colNames, "信箱")
			}
		} else {
			if extRequired != "" {
				colNames = append(colNames, extRequired)
			}
		}
	}

	// 將欄位中文名稱寫入EXCEL第一行
	for i, name := range colNames {
		file.SetCellValue(sheet, rowNames[i]+"1", name) // 設置存儲格的值
	}

	// 將所有報名簽到人員寫入EXCEL中(從第二行開始)
	for i, applySign := range applySigns {
		var (
			statusCN  string
			extValues = []string{
				// 電話.信箱
				applySign.Phone, applySign.ExtEmail,
				// 自定義欄位
				applySign.Ext1, applySign.Ext2, applySign.Ext3,
				applySign.Ext4, applySign.Ext5, applySign.Ext6,
				applySign.Ext7, applySign.Ext8, applySign.Ext9,
				applySign.Ext10,
			}
		)

		// 判斷報名簽到狀態('sign', 'apply', 'review', 'cancel', 'refuse', 'no')
		if applySign.Status == "review" {
			statusCN = "審核中"
		} else if applySign.Status == "apply" {
			statusCN = "已報名"
		} else if applySign.Status == "sign" {
			statusCN = "已簽到"
		} else if applySign.Status == "refuse" {
			statusCN = "拒絕"
		} else if applySign.Status == "no" {
			statusCN = "資料不齊"
		} else if applySign.Status == "cancel" {
			statusCN = "取消簽到"
		}

		values := []interface{}{
			// 用戶資訊
			applySign.Name, applySign.Number, applySign.UserDevice, applySign.ExtPassword,

			// 報名簽到資訊
			statusCN, applySign.ApplyTime, applySign.ReviewTime,
			applySign.SignTime,
		}

		// 判斷自定義欄位
		for i, extRequired := range extRequireds {
			if i == 0 {
				if extRequired == "true" {
					values = append(values, extValues[i])
				}
			} else if i == 1 {
				if extRequired == "true" {
					values = append(values, extValues[i])
				}
			} else {
				if extRequired != "" {
					values = append(values, extValues[i])
				}
			}
		}

		for n, value := range values {
			file.SetCellValue(sheet, rowNames[n]+strconv.Itoa(i+2), value) // 設置存儲格的值
		}
	}

	// 處理excel檔案名稱
	// 判斷是否取得特定狀態資料
	if status != "" {
		statusCN = status
	} else if status == "" {
		statusCN = "所有狀態"
	}
	// 儲存EXCEL
	if err := file.SaveAs(config.STORE_PATH + "/excel/報名簽到-" + activityID + "-" + statusCN + ".xlsx"); err != nil {
		return file, err
	}

	return file, nil
}

// FindSignStaff 查詢活動所有簽到人員資料(redis處理)
func (s ApplysignModel) FindSignStaffs(isRedis bool, isUserInfo bool,
	redisName, activityID string,
	limit, offset int64) ([]ApplysignModel, error) {
	var (
		users      = make([]string, 0) // 包含test測試資料，避免重複執行資料庫
		applysigns = make([]ApplysignModel, 0)
		err        error
	)

	if isRedis {
		if redisName == config.SIGN_STAFFS_2_REDIS { // 簽到人員2(更新資料時，修改redis裡的資料，james用)，SET
			// 判斷redis裡是否有簽到人員資訊，有則不執行查詢資料表功能
			users, err = s.RedisConn.SetGetMembers(redisName + activityID) // 包含test測試資料，避免重複執行資料庫
			if err != nil {
				return applysigns, errors.New("錯誤: 從redis中取得簽到人員資訊發生問題")
			}
		}
		//  else if redisName == config.SIGN_STAFFS_1_REDIS { // 簽到人員1(更新資料時，不修改redis裡的資料，威翔用)，LIST
		// 	// 查詢redis裡是否有簽到人員資料，有則不查詢資料表
		// 	users, err = s.RedisConn.ListRange(redisName+activityID, 0, 0)
		// 	if err != nil {
		// 		return applysigns, errors.New("錯誤: 從redis中取得簽到人員資訊發生問題")
		// 	}

		// log.Println("redis裡的資料: ", len(users))

		// 處理簽到人員資訊
		for _, userID := range users {
			if isUserInfo { // 需要用戶詳細資料
				// 判斷redis裡是否有用戶資訊，有則不查詢資料表
				user, err := DefaultLineModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					Find(true, "", "user_id", userID)
				if err != nil {
					return applysigns, errors.New("錯誤: 無法取得用戶資訊")
				}

				if userID != "test" {
					// 過濾測試資料
					applysigns = append(applysigns, ApplysignModel{
						ID:         user.ID,
						UserID:     userID,
						Name:       user.Name,
						Avatar:     user.Avatar,
						Phone:      user.Phone,
						UserDevice: user.Device,
					})
				}

			} else { // 不需要用戶詳細資料
				if userID != "test" {
					// 過濾測試資料
					applysigns = append(applysigns, ApplysignModel{
						UserID: userID,
					})
				}
			}
		}
	}

	if len(users) == 0 { // redis裡連測試資料都沒有，查詢資料表
		// log.Println("redis裡連測試資料都沒有，查詢資料表: ")

		var (
			items  = make([]map[string]interface{}, 0)
			params = []interface{}{redisName + activityID} // redis資料
		)

		sql := s.Table(s.Base.TableName).
			// Select("activity_applysign.id", "activity_applysign.user_id",
			// 	"activity_id", "status", "apply_sign", "review_time", "sign_time", "number", "ext_1",
			// 	"ext_2", "ext_3", "ext_4", "ext_5", "ext_6", "ext_7", "ext_8", "ext_9", "ext_10",
			// 	"name", "avatar", "email", "identify", "friend", "phone", "email", "ext_email").
			Select("activity_applysign.id", "activity_applysign.user_id",
				"activity_applysign.activity_id", "activity_applysign.number",

				"line_users.name", "line_users.avatar", "line_users.phone",
				"line_users.device as user_device", "line_users.ext_password").
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_applysign.activity_id", "=", activityID).
			Where("activity_applysign.status", "=", "sign").
			OrderBy("activity_applysign.id", "asc")
		if limit != 0 {
			sql = sql.Limit(limit)
		}
		if offset != 0 {
			sql = sql.Offset(offset)
		}

		items, err := sql.All()
		if err != nil {
			return applysigns, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
		}
		applysigns = MapToApplysignModel(items)

		if limit == 0 && offset == 0 { // 無過濾的資料庫查詢時，才將查詢資料寫入redis中
			for _, applysign := range applysigns {
				params = append(params, applysign.UserID)
			}

			if len(applysigns) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				params = append(params, "test")
			}

			// 將資料寫入redis中
			if redisName == config.SIGN_STAFFS_2_REDIS { // 簽到人員2(更新資料時，修改redis裡的資料，james用)，SET
				// 將簽到人員資訊加入redis中(SET)
				s.RedisConn.SetAdd(params)
			}
			// else  if redisName == config.SIGN_STAFFS_1_REDIS { // 簽到人員1(更新資料時，不修改redis裡的資料，威翔用)，LIST
			// 	// 將簽到人員資訊加入redis中(LIST)
			// 	s.RedisConn.ListMultiRPush(staffs)
			// } else

			// s.RedisConn.SetExpire(redisName+activityID,
			// 	config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+activityID, "FindSignStaffs")
		} else {
			// log.Println("有過濾，不加入redis")
		}
	}
	return applysigns, err
}

// Find 查詢報名簽到人員資料
func (s ApplysignModel) Find(id int64, activityID, userID string, isTime bool) (ApplysignModel, error) {
	var (
		sql = s.Table(s.Base.TableName).
			Select("activity_applysign.id", "activity_applysign.user_id",
				"activity_applysign.activity_id", "activity_applysign.status",
				"activity_applysign.number", "activity_applysign.role",

				// 活動資訊
				"activity.user_id as activity_user_id", "activity.activity_name",
				"activity.start_time", "activity.end_time",
				"activity.sign_check", "activity.sign_allow", "activity.sign_minutes",
				"activity.overview_message", "activity.message_ban", "activity.message_ban_second",
				"activity.overview_danmu", "activity.overview_special_danmu", "activity.overview_holdscreen",
				"activity.holdscreen_birthday_topic", "activity.overview_redpack",
				"activity.overview_ropepack", "activity.overview_whack_mole", "activity.overview_lottery",
				"activity.overview_question", "activity.overview_monopoly", "activity.overview_qa",
				"activity.overview_tugofwar", "activity.overview_bingo", "activity.overview_signname",
				"activity.overview_3d_gacha_machine",
				"activity.device as activity_device",

				"activity_2.overview_vote",

				// 活動官方帳號
				"activity.line_id as activity_line_id",
				"activity.channel_id as activity_channel_id",
				"activity.channel_secret as activity_channel_secret",
				"activity.chatbot_secret as activity_chatbot_secret",
				"activity.chatbot_token as activity_chatbot_token",

				// 用戶官方帳號
				"users.line_id as user_line_id",
				"users.channel_id as user_channel_id",
				"users.channel_secret as user_channel_secret",
				"users.chatbot_secret as user_chatbot_secret",
				"users.chatbot_token as user_chatbot_token",

				// 自定義欄位
				"activity_applysign.ext_1", "activity_customize.ext_1_name", "activity_customize.ext_1_type",
				"activity_customize.ext_1_options", "activity_customize.ext_1_required",
				"activity_applysign.ext_2", "activity_customize.ext_2_name", "activity_customize.ext_2_type",
				"activity_customize.ext_2_options", "activity_customize.ext_2_required",
				"activity_applysign.ext_3", "activity_customize.ext_3_name", "activity_customize.ext_3_type",
				"activity_customize.ext_3_options", "activity_customize.ext_3_required",
				"activity_applysign.ext_4", "activity_customize.ext_4_name", "activity_customize.ext_4_type",
				"activity_customize.ext_4_options", "activity_customize.ext_4_required",
				"activity_applysign.ext_5", "activity_customize.ext_5_name", "activity_customize.ext_5_type",
				"activity_customize.ext_5_options", "activity_customize.ext_5_required",
				"activity_applysign.ext_6", "activity_customize.ext_6_name", "activity_customize.ext_6_type",
				"activity_customize.ext_6_options", "activity_customize.ext_6_required",
				"activity_applysign.ext_7", "activity_customize.ext_7_name", "activity_customize.ext_7_type",
				"activity_customize.ext_7_options", "activity_customize.ext_7_required",
				"activity_applysign.ext_8", "activity_customize.ext_8_name", "activity_customize.ext_8_type",
				"activity_customize.ext_8_options", "activity_customize.ext_8_required",
				"activity_applysign.ext_9", "activity_customize.ext_9_name", "activity_customize.ext_9_type",
				"activity_customize.ext_9_options", "activity_customize.ext_9_required",
				"activity_applysign.ext_10", "activity_customize.ext_10_name", "activity_customize.ext_10_type",
				"activity_customize.ext_10_options", "activity_customize.ext_10_required",

				"activity_customize.ext_1_unique", "activity_customize.ext_2_unique", "activity_customize.ext_3_unique", "activity_customize.ext_4_unique",
				"activity_customize.ext_5_unique", "activity_customize.ext_6_unique", "activity_customize.ext_7_unique", "activity_customize.ext_8_unique",
				"activity_customize.ext_9_unique", "activity_customize.ext_10_unique",

				"activity_customize.ext_email_required", "activity_customize.ext_phone_required",
				"activity_customize.info_picture",

				// LINE用戶資訊
				"line_users.name", "line_users.avatar",
				"line_users.phone", "line_users.email", "line_users.ext_email",
				"line_users.friend", "line_users.identify",
				"line_users.device as user_device", "line_users.ext_password",
				"line_users.admin_id").
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.activity_id",
				FieldA1:   "activity_2.activity_id",
				Table:     "activity_2",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "users.user_id",
				FieldA1:   "activity.user_id",
				Table:     "users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_applysign.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="})
		item   map[string]interface{}
		err    error
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)

	// 判斷參數是否為空
	if id != 0 {
		sql = sql.Where("activity_applysign.id", "=", id)
	}
	if userID != "" {
		sql = sql.Where("activity_applysign.user_id", "=", userID)
	}
	if activityID != "" {
		sql = sql.Where("activity_applysign.activity_id", "=", activityID)
	}
	if isTime {
		sql = sql.Where("activity.end_time", ">=", now)
	}

	item, err = sql.First()
	if err != nil {
		return ApplysignModel{}, errors.New("錯誤: 無法取得報名簽到人員資訊，請重新查詢")
	}
	if item == nil {
		return ApplysignModel{}, nil
	}
	return s.MapToModel(item), nil
}

// IsSign 是否簽到(目前都是從redis判斷，目前只有判斷SIGN_STAFFS_2_REDIS裡的資料)
func (s ApplysignModel) IsSign(redisName, activityID, userID string) (isSign bool) {
	var (
		// staffs = make([]string, 0)
		people int64
	)

	// 判斷redis裡是否有簽到人員資訊，有則不執行查詢資料表功能
	people = s.RedisConn.SetCard(redisName + activityID)

	// redis中沒有簽到人員資訊，查詢資料表並加入redis中
	if people == 0 {
		// log.Println("redis中目前沒有資料，執行FindSignStaffs")

		// 從redis取得資料，確定redis中有該場活動報名簽到人員資料(sign_staffs_2_activityID)
		s.FindSignStaffs(true, false, config.SIGN_STAFFS_2_REDIS, activityID, 0, 0)
	} else {
		// log.Println("reids中有資料，直接判斷是否簽到")
	}

	// 判斷是否簽到
	isSign = s.RedisConn.SetIsMember(redisName+activityID, userID)

	return
}

// GetAttend 查詢資料
func (s ApplysignModel) GetAttend(activity string) (int, error) {
	items, err := s.Table(s.Base.TableName).Where("activity_id", "=", activity).
		Where("status", "=", "sign").All()
	if err != nil {
		return 0, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
	}
	return len(items), nil
}

// Add 增加資料(傳送訊息)
func (s ApplysignModel) Add(isRedis bool, model NewApplysignModel) (int64, error) {
	id, err := s.Table(s.TableName).Insert(command.Value{
		"user_id":     model.UserID,
		"activity_id": model.ActivityID,
		"status":      model.Status,
		"number":      0,
	})
	if err != nil {
		return id, errors.New("錯誤: 無法新增報名簽到人員資料，請重新操作")
	}

	if isRedis {
		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "Add")
	}

	return id, nil
}

// UpdateExt 更新報名簽到人員Ext欄位資料
func (s ApplysignModel) UpdateExt(activityID, userID string, values []string) error {
	var (
		fields = []string{"ext_1", "ext_2", "ext_3",
			"ext_4", "ext_5", "ext_6", "ext_7", "ext_8", "ext_9", "ext_10"}
		fieldValues = command.Value{}
	)

	// 取得自定義欄位資料
	customizeModel, err := DefaultCustomizeModel().
		SetConn(s.DbConn, s.RedisConn, s.MongoConn).
		Find(activityID)
	if err != nil || customizeModel.ID == 0 {
		return errors.New("錯誤: 查詢自定義欄位資料發生問題，請重新操作")
	}

	// 自定義匯入報名簽到人員
	var (
		uniques = []string{
			customizeModel.Ext1Unique, customizeModel.Ext2Unique, customizeModel.Ext3Unique, customizeModel.Ext4Unique,
			customizeModel.Ext5Unique, customizeModel.Ext6Unique, customizeModel.Ext7Unique, customizeModel.Ext8Unique,
			customizeModel.Ext9Unique, customizeModel.Ext10Unique,
		}
		requires = []string{
			customizeModel.Ext1Required, customizeModel.Ext2Required, customizeModel.Ext3Required, customizeModel.Ext4Required,
			customizeModel.Ext5Required, customizeModel.Ext6Required, customizeModel.Ext7Required, customizeModel.Ext8Required,
			customizeModel.Ext9Required, customizeModel.Ext10Required,
		}

		exts = make([][]string, 10)
	)

	// 取得該活動所有報名簽到人員資料
	applysigns, err := DefaultApplysignModel().
		SetConn(s.DbConn, s.RedisConn, s.MongoConn).
		FindAll(activityID, "", "", "", 0, 0)
	if err != nil {
		return errors.New("錯誤: 查詢簽到人員資料發生問題，請重新操作")
	}

	// 取得原有報名簽到資料的唯一值
	for _, applysign := range applysigns {
		applysginExts := []string{
			applysign.Ext1, applysign.Ext2, applysign.Ext3, applysign.Ext4,
			applysign.Ext5, applysign.Ext6, applysign.Ext7, applysign.Ext8,
			applysign.Ext9, applysign.Ext10,
		}

		for n, ext := range applysginExts {
			// 判斷自定義欄位是否為唯一值，是的話加入陣列中
			if uniques[n] == "true" {
				exts[n] = append(exts[n], ext)
			}
		}
	}

	for n, uniques := range uniques {
		if uniques == "true" {
			if values[n] == "" {
				return errors.New("錯誤: 唯一值欄位不能為空，請輸入有效資料")
			}

			if utils.InArray(exts[n], values[n]) {
				return errors.New("錯誤: 唯一值欄位不能重複，請輸入有效資料")
			}
		}

		if requires[n] == "true" && values[n] == "" {
			return errors.New("錯誤: 必填欄位不能為空，請輸入有效資料")
		}
	}

	for i, field := range fields {
		if values[i] != "" {
			fieldValues[field] = values[i]
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	if err := s.Table(s.Base.TableName).
		Where("activity_id", "=", activityID).
		Where("user_id", "=", userID).
		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新報名簽到人員Ext欄位資料發生問題，請重新操作")
	}
	return nil
}

// DeleteExt 清空報名簽到人員Ext欄位資料
func (s ApplysignModel) DeleteExt(activityID, field string) error {
	if err := s.Table(s.Base.TableName).Where("activity_id", "=", activityID).
		Update(command.Value{field: ""}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新報名簽到人員Ext欄位資料發生問題，請重新操作")
	}
	return nil
}

// UpdateRole 更新簽到人員角色資料
func (s ApplysignModel) UpdateRole(isRedis bool, model EditApplysignModel) error {

	// 更新角色資料
	if err := s.Table(s.Base.TableName).
		Where("activity_id", "=", model.ActivityID).
		WhereIn("user_id", interfaces(model.LineUsers)).
		Update(command.Value{"role": model.Role}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新報名簽到人角色資料發生問題，請重新操作")
	}

	return nil
}

// UpdateStatus 更新報名簽到狀態(傳送訊息)
func (s ApplysignModel) UpdateStatus(isRedis bool, host string, model EditApplysignModel, isSign bool) error {
	var (
		fieldValues = command.Value{"status": model.Status}
		hosts       = []string{config.HILIVES_NET_URL}
		status      = []string{"review", "apply", "sign", "refuse", "cancel"}
		// number                               int64
		hostBool, statusBool, isPushMessage  bool
		chatbotSecret, chatbotToken, message string
		err                                  error
		// url    = fmt.Sprintf("https://%s/applysign?activity_id=%s",
		// 	host, model.ActivityID)
		// hostBool, statusBool                          bool
	)
	for i := range hosts {
		if hosts[i] == host {
			hostBool = true
			break
		}
	}
	if !hostBool {
		return errors.New("錯誤: 網域發生問題，請輸入有效的網域名稱")
	}
	// if host == config.HILIVES_NET_URL {
	// 	liffURL = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, model.ActivityID)
	// }

	for i := range status {
		if status[i] == model.Status {
			statusBool = true
			break
		}
	}
	if !statusBool {
		return errors.New("錯誤: 報名簽到狀態發生問題，請輸入有效的報名簽到狀態")
	}

	// 取得活動資訊
	activityModel, err := DefaultActivityModel().
		SetConn(s.DbConn, s.RedisConn, s.MongoConn).
		Find(false, model.ActivityID)
	// number = activityModel.Number
	if err != nil || activityModel.ID == 0 {
		return errors.New("錯誤: 無法取得活動資訊，請重新查詢")
	}

	// 判斷是否綁定官方帳號
	if activityModel.ActivityChannelID != "" && activityModel.ActivityChannelSecret != "" &&
		activityModel.ActivityChatbotSecret != "" && activityModel.ActivityChatbotToken != "" &&
		activityModel.ActivityLineID != "" {
		// 活動官方帳號
		chatbotSecret = activityModel.ActivityChatbotSecret
		chatbotToken = activityModel.ActivityChatbotToken
		isPushMessage = true
	} else if activityModel.UserChannelID != "" && activityModel.UserChannelSecret != "" &&
		activityModel.UserChatbotSecret != "" && activityModel.UserChatbotToken != "" &&
		activityModel.UserLineID != "" {
		// 用戶活動官方帳號
		chatbotSecret = activityModel.UserChatbotSecret
		chatbotToken = activityModel.UserChatbotToken
		isPushMessage = true
	} else {
		// 我方公司官方帳號
		chatbotSecret = config.CHATBOT_SECRET
		chatbotToken = config.CHATBOT_TOKEN

		if activityModel.PushMessage == "open" {
			isPushMessage = true
		}
	}

	// 報名成功、拒絕
	if model.Status == "apply" || model.Status == "refuse" {
		fieldValues["review_time"] = model.ReviewTime

		if model.Status == "apply" { // 通過報名

			var (
				mailAmount    = activityModel.MailAmount    // 郵件數量
				messageAmount = activityModel.MessageAmount // 簡訊數量
				// remain        int64
				// sendAmount    int64
			)

			// 郵件處理
			if activityModel.SendMail == "open" && mailAmount > 0 {
				var (
				// userIDs   = make([]string, 0)
				// names     = make([]string, 0)
				// devices   = make([]string, 0)
				// emails    = make([]string, 0)
				// passwords = make([]string, 0)
				)
				for _, userID := range model.LineUsers {
					if mailAmount > 0 { // 剩餘數量大於0，傳送郵件
						// 取得用戶資料
						lineModel, err := DefaultLineModel().
							SetConn(s.DbConn, s.RedisConn, s.MongoConn).
							Find(true, "", "user_id", userID)
						if err != nil || lineModel.UserID == "" {
							return errors.New("錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
						}

						if lineModel.ExtEmail != "" {
							// 郵件處理
							url := "" // 報名簽到連結
							qrcode := fmt.Sprintf(config.HTTPS_APPLYSIGN_QRCODE_URL, config.HILIVES_NET_URL, fmt.Sprintf("activity_id=%s.user_id=%s", model.ActivityID, userID))
							qrcodeMessage := "<p></p>" // qrcode訊息

							if lineModel.Device == "facebook" || lineModel.Device == "gmail" || lineModel.Device == "customize" {
								// 一般url
								url = fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, model.ActivityID, userID, "") // 報名簽到連結
							} else if lineModel.Device == "line" {
								// liff url
								url = fmt.Sprintf(config.HILIVES_APPLYSIGN_URL_LIFF_URL, model.ActivityID, userID)
							}

							subject := fmt.Sprintf("Subject: %s 活動報名簽到訊息\r\n", activityModel.ActivityName)
							mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

							// 判斷主持人掃描qrcode開關是否開啟
							if activityModel.HostScan == "open" {
								qrcodeMessage = fmt.Sprintf(`
						<p>以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到判斷(避免資料洩漏，連結勿提供給他人使用)：</p>
						<img src="https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=%s" alt="簽到QRcode"/>
						`, qrcode)
							}

							body := ""
							if lineModel.Device == "facebook" || lineModel.Device == "gmail" || lineModel.Device == "line" {
								body = fmt.Sprintf(`
	    					<html>
								<body>
									<p>%s 您好，</p>
									<p>感謝您報名參加 %s 活動！以下是您的簽到信息：</p>
	    							<p>可利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
									<p><a href="%s"> %s 活動簽到連結</a></p>
									%s
									<p>如有任何問題，請隨時與我們聯繫。</p>
									<p>祝您一切順利！</p>
									<p>此致，</p>
									<p>活動團隊Hilives</p>
								</body>
							</html>`, lineModel.Name, activityModel.ActivityName, url, activityModel.ActivityName, qrcodeMessage)
							} else if lineModel.Device == "customize" {
								body = fmt.Sprintf(`
	    					<html>
								<body>
									<p>%s 您好，</p>
									<p>感謝您報名參加 %s 活動！以下是您的簽到信息：</p>
									<p>該活動的驗證碼為%s，可透過驗證碼進行簽到，</p>
	    							<p>或利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
									<p><a href="%s"> %s 活動簽到連結</a></p>
									%s
									<p>如有任何問題，請隨時與我們聯繫。</p>
									<p>祝您一切順利！</p>
									<p>此致，</p>
									<p>活動團隊Hilives</p>
								</body>
							</html>`, lineModel.Name, activityModel.ActivityName, lineModel.ExtPassword, url, activityModel.ActivityName, qrcodeMessage)
							}
							message = subject + mime + body

							// 發送郵件
							err = SendMail(lineModel.ExtEmail, message)
							if err != nil {
								return err
							}

							mailAmount-- // 剩餘數量等於0
						}
					} else {
						break // 數量等於0
					}
				}

				// 更新郵件剩餘數量
				err = DefaultActivityModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					UpdateMailAmount(true, model.ActivityID, mailAmount)
				if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					// 其他錯誤
					return err
				}
			}

			// 簡訊處理
			if activityModel.PushPhoneMessage == "open" && messageAmount > 0 {
				var (
				// userIDs   = make([]string, 0)
				// names     = make([]string, 0)
				// devices   = make([]string, 0)
				// phones    = make([]string, 0)
				// passwords = make([]string, 0)
				)

				for _, userID := range model.LineUsers {
					if messageAmount > 0 { // 剩餘數量大於0，傳送簡訊
						// 取得用戶資料
						lineModel, err := DefaultLineModel().
							SetConn(s.DbConn, s.RedisConn, s.MongoConn).
							Find(true, "", "user_id", userID)
						if err != nil || lineModel.UserID == "" {
							return errors.New("錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
						}

						if lineModel.Phone != "" {
							// 簡訊訊息處理(簡訊無法有net的字串連結，因此使用liff url)
							url := fmt.Sprintf(config.HILIVES_APPLYSIGN_URL_LIFF_URL, model.ActivityID, userID)

							if lineModel.Device == "facebook" || lineModel.Device == "gmail" || lineModel.Device == "line" {
								message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，可利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): %s",
									lineModel.Name, activityModel.ActivityName, url)
							} else if lineModel.Device == "customize" {
								message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，該活動的驗證碼為%s，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): %s",
									lineModel.Name, activityModel.ActivityName, lineModel.ExtPassword, url)
							}

							// 傳送簡訊
							err = SendMessage(lineModel.Phone, message)
							if err != nil {
								return err
							}

							messageAmount-- // 遞減簡訊數量
						}
					} else {
						break // 剩餘數量等於0
					}
				}

				// 更新簡訊剩餘數量
				err = DefaultActivityModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					UpdateMessageAmount(true, model.ActivityID, messageAmount)
				if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					// 其他錯誤
					return err
				}
			}

			for _, userID := range model.LineUsers {
				// 取得用戶資料
				lineModel, err := DefaultLineModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					Find(true, "", "user_id", userID)
				if err != nil || lineModel.UserID == "" {
					return errors.New("錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
				}

				if lineModel.Device == "line" && isPushMessage {
					// line裝置
					message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，可利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): ",
						lineModel.Name, activityModel.ActivityName)
					// url := fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, model.ActivityID) + "&sign=open"
					url := fmt.Sprintf(config.HILIVES_APPLYSIGN_URL_LIFF_URL, model.ActivityID, userID)

					// 使用第三方的官方帳號才需傳送訊息
					if err = pushMessage(chatbotSecret, chatbotToken, userID, message+url); err != nil {
						return err
					}
				}
			}

			// 更新報名狀態
			err = s.Table(s.Base.TableName).
				Where("activity_id", "=", model.ActivityID).
				WhereIn("user_id", interfaces(model.LineUsers)).
				Update(fieldValues)

			// 判斷更新用戶超過一人時(後臺批量處理)檢查簽到人數並更新
			if len(model.LineUsers) > 1 {
				// 取得新的已簽到人數
				attend, err := DefaultApplysignModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					GetAttend(model.ActivityID)
				if err != nil {
					return nil
				}

				// 更新活動人數
				if err = DefaultActivityModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					UpdateAttendAndNumber(false, model.ActivityID, attend, 0, []string{}); err != nil &&
					err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					return err
				}
			}

			// 刪除redis活動相關資訊
			// s.RedisConn.DelCache(config.SIGN_STAFFS_1_REDIS + model.ActivityID) // 簽到人員資訊
			// s.RedisConn.DelCache(config.SIGN_STAFFS_2_REDIS + model.ActivityID) // 簽到人員資訊

		} else if model.Status == "refuse" { // 拒絕
			// 拒絕只有後端平台可以處理
			// refuse有人數增減問題，先更新人數再更新簽到人數
			err = s.Table(s.Base.TableName).
				Where("activity_id", "=", model.ActivityID).
				WhereIn("user_id", interfaces(model.LineUsers)).
				Update(fieldValues)

			// 取得新的已簽到人數
			attend, err := DefaultApplysignModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				GetAttend(model.ActivityID)
			if err != nil {
				return nil
			}

			// 更新活動人數
			if err = DefaultActivityModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				UpdateAttendAndNumber(false, model.ActivityID, attend, 0, []string{}); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}

			// 刪除redis活動相關資訊
			// s.RedisConn.DelCache(config.SIGN_STAFFS_1_REDIS + model.ActivityID) // 簽到人員資訊
			// s.RedisConn.DelCache(config.SIGN_STAFFS_2_REDIS + model.ActivityID) // 簽到人員資訊
		}

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "apply與refuse")
		return err
	}

	if model.Status == "sign" || model.Status == "cancel" {
		// 簽到、取消簽到
		if model.SignTime != "" {
			fieldValues["sign_time"] = model.SignTime
		}

		if model.Status == "sign" {
			var (
				lineUsers  = model.LineUsers
				attend     = activityModel.Attend // 參加活動人數
				maxPeople  = activityModel.People // 活動上限人數
				dataAmount = len(model.LineUsers) // 更新的資料數
				newNumber  int64                  // 要更新的號碼資料
				signs      = make([]string, 0)    // 簽到人員(寫入redis中)
			)

			// log.Println("判斷簽到人數是否額滿")
			if maxPeople-attend <= 0 {
				// log.Println("人數額滿，回傳錯誤")
				return errors.New("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
			}

			// 取得活動號碼資料
			number, err := DefaultActivityModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				GetActivityNumber(true, model.ActivityID)
			if err != nil {
				return err
			}

			if maxPeople-attend < int64(dataAmount) { // 可簽到人數小於更新資料數
				signPeople := maxPeople - attend   // 可簽到人數
				lineUsers = lineUsers[:signPeople] // 可簽到人員

				// 增加活動參加人數跟號碼資訊(增加至上限值)
				attend = attend + signPeople
				newNumber = number + signPeople
			} else {
				// 增加活動參加人數跟號碼資訊
				attend = attend + int64(dataAmount)
				newNumber = number + int64(dataAmount)
			}

			// 將簽到狀態人員寫入陣列中
			// for _, userID := range lineUsers {
			signs = append(signs, lineUsers...)
			// }

			// 更新參加人數跟號碼資訊
			if err = DefaultActivityModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				UpdateAttendAndNumber(true, model.ActivityID, int(attend), int(newNumber), signs); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}

			// 用戶分別更新(需更新號碼資料)
			for i, userID := range lineUsers {
				// 發放抽獎號碼給用戶
				fieldValues["number"] = number + int64(i)

				// 分別更新並發放號碼給用戶(號碼加一)
				s.Table(s.Base.TableName).
					Where("activity_id", "=", model.ActivityID).
					Where("user_id", "=", userID).
					Update(fieldValues)
			}

			// 判斷更新用戶超過一人時(後臺批量處理)檢查簽到人數並更新
			if len(model.LineUsers) > 1 {
				// 取得新的已簽到人數
				attendInt, err := DefaultApplysignModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					GetAttend(model.ActivityID)
				if err != nil {
					return nil
				}

				// 更新活動人數
				if err = DefaultActivityModel().
					SetConn(s.DbConn, s.RedisConn, s.MongoConn).
					UpdateAttendAndNumber(false, model.ActivityID, attendInt, 0, []string{}); err != nil &&
					err.Error() != "錯誤: 無更新任何資料，請重新操作" {
					return err
				}
			}

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "sign")

		} else if model.Status == "cancel" {
			// 取消簽到只有後端平台能批量更新

			// var attend = activityModel.Attend // 參加活動人數

			// 更新活動人數(不更新號碼)
			// if err = DefaultActivityModel().SetDbConn(s.DbConn).SetRedisConn(s.RedisConn).
			// 	UpdateAttendAndNumber(false, model.ActivityID,
			// 		int(attend)-len(model.LineUsers), 0, []string{}); err != nil &&
			// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			// 	return err
			// }

			// for _, userID := range model.LineUsers {
			// 	// 遞減活動人數
			// 	if err = activityModel.DecrAttend(true, model.ActivityID,
			// 		userID); err != nil {
			// 		return err
			// 	}
			// }

			// 取消簽到，將狀態改回apply(報名成功)
			fieldValues["status"] = "apply"

			// 更新所有用戶狀態
			err = s.Table(s.Base.TableName).
				Where("activity_id", "=", model.ActivityID).
				WhereIn("user_id", interfaces(model.LineUsers)).
				Update(fieldValues)
			if err != nil {
				return nil
			}

			// 取得新的已簽到人數
			attend, err := DefaultApplysignModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				GetAttend(model.ActivityID)
			if err != nil {
				return nil
			}

			// 更新活動人數
			if err = DefaultActivityModel().
				SetConn(s.DbConn, s.RedisConn, s.MongoConn).
				UpdateAttendAndNumber(false, model.ActivityID, attend, 0, []string{}); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}

			// 刪除redis活動相關資訊(取消簽到時會變更簽到人員redis裡的資料)
			// s.RedisConn.DelCache(config.SIGN_STAFFS_1_REDIS + model.ActivityID) // 簽到人員資訊
			s.RedisConn.DelCache(config.SIGN_STAFFS_2_REDIS + model.ActivityID) // 簽到人員資訊

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "cancel")

			return err
		}
	}

	if model.Status == "review" {
		// 審核
		err = s.Table(s.Base.TableName).
			Where("activity_id", "=", model.ActivityID).
			WhereIn("user_id", interfaces(model.LineUsers)).
			Update(fieldValues)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "review")

		return err
	}

	return nil
}

// UpdateApply 更新報名牆基本設置資料
func (a ActivityModel) UpdateApply(model EditApplyModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"apply_check", "customize_password", "allow_customize_apply"}
		// values      = []string{model.ApplyCheck, model.CustomizePassword, model.AllowCustomizeApply}

		fieldValues2 = command.Value{}
		fields2      = []string{"customize_default_avatar"}
		// values2      = []string{model.CustomizeDefaultAvatar}
		err error
	)

	if model.ApplyCheck != "" {
		if model.ApplyCheck != "open" && model.ApplyCheck != "close" {
			return errors.New("錯誤: 報名審核資料發生問題，請輸入有效的報名審核資料")
		}
	}
	if model.CustomizePassword != "" {
		if model.CustomizePassword != "open" && model.CustomizePassword != "close" {
			return errors.New("錯誤: 自定義設置驗證碼資料發生問題，請輸入有效的資料")
		}
	}
	if model.AllowCustomizeApply != "" {
		if model.AllowCustomizeApply != "open" && model.AllowCustomizeApply != "close" {
			return errors.New("錯誤: 是否允許用戶自定義報名資料發生問題，請輸入有效的資料")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		if err = a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
	}

	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			fieldValues2[key] = val
		}
	}

	if len(fieldValues2) != 0 {
		// log.Println("更新自定義人員預設頭像: ", model.CustomizeDefaultAvatar)

		if err = a.Table(config.ACTIVITY_2_TABLE).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues2); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}

		if model.CustomizeDefaultAvatar != "" {
			// log.Println("修改自定義人員預設頭像，必須執行其他處理")
			var (
				userIDs = make([]interface{}, 0) // 需要更新的用戶
			)

			// 取得所有舊的匯入人員資料，將頭像改為新的資料路徑
			// 判斷舊的人員資料是否更換為自己的頭像(判斷is_modify參數)
			users, err := DefaultLineModel().
				SetConn(a.DbConn, a.RedisConn, a.MongoConn).
				FindAllCustomizeUsers(model.ActivityID)
			if err != nil {
				return err
			}

			for _, user := range users {
				if user.IsModify == "no" {
					// log.Println("未更換為自己的頭像，將該用戶頭像改為新的資料路徑")

					// 未更換為自己的頭像，將該用戶頭像改為新的資料路徑
					userIDs = append(userIDs, user.UserID)

					// 需要修改的用戶必須清除舊的redis資料
					a.RedisConn.DelCache(config.AUTH_USERS_REDIS + user.UserID)
				} else {
					// log.Println("已更換為自己的頭像，不用更新資料")
				}
			}

			if len(userIDs) > 0 {
				// log.Println("批量更新頭像資料")
				// 需要更新的資料數量大於0，進行頭像批量更新
				err = DefaultLineModel().
					SetConn(a.DbConn, a.RedisConn, a.MongoConn).
					UpdateAvatar(userIDs, model.CustomizeDefaultAvatar)
				if err != nil {
					return err
				}

			} else {
				// log.Println("不用批量更新頭像資料")
			}
		}
	}
	return nil
}

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

// UpdateSign 更新簽到強基本設置資料
func (a ActivityModel) UpdateSign(model EditSignModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"sign_check", "sign_allow", "sign_minutes", "sign_manual"}
		// values      = []string{model.SignCheck, model.SignAllow,
		// 	model.SignMinutes, model.SignManual}
	)
	if model.SignCheck != "" {
		if model.SignCheck != "open" && model.SignCheck != "close" {
			return errors.New("錯誤: 簽到審核資料發生問題，請輸入有效的簽到審核資料")
		}
	}
	if model.SignAllow != "" {
		if model.SignAllow != "open" && model.SignAllow != "close" {
			return errors.New("錯誤: 允許讓參加人員即刻簽到資料發生問題，請輸入有效的即刻簽到資料")
		}
	}
	if model.SignMinutes != "" {
		if _, err := strconv.Atoi(model.SignMinutes); err != nil {
			return errors.New("錯誤: 簽到時間發生問題，請輸入有效的簽到時間")
		}
	}
	if model.SignManual != "" {
		if model.SignManual != "open" && model.SignManual != "close" {
			return errors.New("錯誤: 手動簽到資料發生問題，請輸入有效的手動簽到資料")
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

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

// MapToModel 設置ApplysignModel
func (s ApplysignModel) MapToModel(m map[string]interface{}) ApplysignModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &s)

	if s.InfoPicture != "" {
		s.InfoPicture = "/admin/uploads/" + s.ActivityUserID + "/" + s.ActivityID + "/applysign/customize/" + s.InfoPicture
	}

	return s
}

// MapToApplysignModel map轉換[]ApplysignModel
func MapToApplysignModel(items []map[string]interface{}) []ApplysignModel {
	var applySigns = make([]ApplysignModel, 0)
	for _, item := range items {
		var (
			applySign ApplysignModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &applySign)

		if applySign.InfoPicture != "" {
			applySign.InfoPicture = "/admin/uploads/" + applySign.ActivityUserID + "/" + applySign.ActivityID + "/applysign/customize/" + applySign.InfoPicture
		}

		applySigns = append(applySigns, applySign)

	}
	return applySigns
}

// applySign.ID, _ = item["id"].(int64)
// applySign.UserID, _ = item["user_id"].(string)
// applySign.ActivityID, _ = item["activity_id"].(string)
// applySign.Status, _ = item["status"].(string)
// applySign.Number, _ = item["number"].(int64)
// applySign.ApplyTime, _ = item["apply_time"].(string)
// applySign.ReviewTime, _ = item["review_time"].(string)
// applySign.SignTime, _ = item["sign_time"].(string)

// // 自定義用戶資料
// applySign.Ext1, _ = item["ext_1"].(string)
// applySign.Ext2, _ = item["ext_2"].(string)
// applySign.Ext3, _ = item["ext_3"].(string)
// applySign.Ext4, _ = item["ext_4"].(string)
// applySign.Ext5, _ = item["ext_5"].(string)
// applySign.Ext6, _ = item["ext_6"].(string)
// applySign.Ext7, _ = item["ext_7"].(string)
// applySign.Ext8, _ = item["ext_8"].(string)
// applySign.Ext9, _ = item["ext_9"].(string)
// applySign.Ext10, _ = item["ext_10"].(string)

// // 活動資訊
// applySign.ActivityUserID, _ = item["activity_user_id"].(string)
// applySign.ActivityName, _ = item["activity_name"].(string)
// applySign.StartTime, _ = item["start_time"].(string)
// applySign.EndTime, _ = item["end_time"].(string)
// applySign.ActivityDevice, _ = item["activity_device"].(string)

// // 用戶資訊
// applySign.Name, _ = item["name"].(string)
// applySign.Avatar, _ = item["avatar"].(string)
// applySign.Phone, _ = item["phone"].(string)
// applySign.Email, _ = item["email"].(string)
// applySign.ExtEmail, _ = item["ext_email"].(string)
// applySign.UserDevice, _ = item["user_device"].(string)
// applySign.ExtPassword, _ = item["ext_password"].(string)

// // 自定義欄位資訊
// applySign.Ext1Name, _ = item["ext_1_name"].(string)
// applySign.Ext1Type, _ = item["ext_1_type"].(string)
// applySign.Ext1Options, _ = item["ext_1_options"].(string)
// applySign.Ext1Required, _ = item["ext_1_required"].(string)

// applySign.Ext2Name, _ = item["ext_2_name"].(string)
// applySign.Ext2Type, _ = item["ext_2_type"].(string)
// applySign.Ext2Options, _ = item["ext_2_options"].(string)
// applySign.Ext2Required, _ = item["ext_2_required"].(string)

// applySign.Ext3Name, _ = item["ext_3_name"].(string)
// applySign.Ext3Type, _ = item["ext_3_type"].(string)
// applySign.Ext3Options, _ = item["ext_3_options"].(string)
// applySign.Ext3Required, _ = item["ext_3_required"].(string)

// applySign.Ext4Name, _ = item["ext_4_name"].(string)
// applySign.Ext4Type, _ = item["ext_4_type"].(string)
// applySign.Ext4Options, _ = item["ext_4_options"].(string)
// applySign.Ext4Required, _ = item["ext_4_required"].(string)

// applySign.Ext5Name, _ = item["ext_5_name"].(string)
// applySign.Ext5Type, _ = item["ext_5_type"].(string)
// applySign.Ext5Options, _ = item["ext_5_options"].(string)
// applySign.Ext5Required, _ = item["ext_5_required"].(string)

// applySign.Ext6Name, _ = item["ext_6_name"].(string)
// applySign.Ext6Type, _ = item["ext_6_type"].(string)
// applySign.Ext6Options, _ = item["ext_6_options"].(string)
// applySign.Ext6Required, _ = item["ext_6_required"].(string)

// applySign.Ext7Name, _ = item["ext_7_name"].(string)
// applySign.Ext7Type, _ = item["ext_7_type"].(string)
// applySign.Ext7Options, _ = item["ext_7_options"].(string)
// applySign.Ext7Required, _ = item["ext_7_required"].(string)

// applySign.Ext8Name, _ = item["ext_8_name"].(string)
// applySign.Ext8Type, _ = item["ext_8_type"].(string)
// applySign.Ext8Options, _ = item["ext_8_options"].(string)
// applySign.Ext8Required, _ = item["ext_8_required"].(string)

// applySign.Ext9Name, _ = item["ext_9_name"].(string)
// applySign.Ext9Type, _ = item["ext_9_type"].(string)
// applySign.Ext9Options, _ = item["ext_9_options"].(string)
// applySign.Ext9Required, _ = item["ext_9_required"].(string)

// applySign.Ext10Name, _ = item["ext_10_name"].(string)
// applySign.Ext10Type, _ = item["ext_10_type"].(string)
// applySign.Ext10Options, _ = item["ext_10_options"].(string)
// applySign.Ext10Required, _ = item["ext_10_required"].(string)

// applySign.Ext1Unique, _ = item["ext_1_unique"].(string)
// applySign.Ext2Unique, _ = item["ext_2_unique"].(string)
// applySign.Ext3Unique, _ = item["ext_3_unique"].(string)
// applySign.Ext4Unique, _ = item["ext_4_unique"].(string)
// applySign.Ext5Unique, _ = item["ext_5_unique"].(string)
// applySign.Ext6Unique, _ = item["ext_6_unique"].(string)
// applySign.Ext7Unique, _ = item["ext_7_unique"].(string)
// applySign.Ext8Unique, _ = item["ext_8_unique"].(string)
// applySign.Ext9Unique, _ = item["ext_9_unique"].(string)
// applySign.Ext10Unique, _ = item["ext_10_unique"].(string)

// applySign.ExtEmailRequired, _ = item["ext_email_required"].(string)
// applySign.ExtPhoneRequired, _ = item["ext_phone_required"].(string)
// applySign.InfoPicture, _ = item["info_picture"].(string)
// if applySign.InfoPicture != "" {
// 	applySign.InfoPicture = "/admin/uploads/" + applySign.ActivityUserID + "/" + applySign.ActivityID + "/applysign/customize/" + applySign.InfoPicture
// }

// log.Println("人數資料: ", s.AllPeople, s.SignPeople, s.ApplyPeople, s.ReviewPeople, s.RefusePeople, s.NoPeople)

// s.ID, _ = m["id"].(int64)
// s.JoinID, _ = m["join_id"].(int64)
// s.UserID, _ = m["user_id"].(string)
// s.ActivityID, _ = m["activity_id"].(string)
// s.Status, _ = m["status"].(string)
// s.Number, _ = m["number"].(int64)
// s.Ext1, _ = m["ext_1"].(string)
// s.Ext2, _ = m["ext_2"].(string)
// s.Ext3, _ = m["ext_3"].(string)
// s.Ext4, _ = m["ext_4"].(string)
// s.Ext5, _ = m["ext_5"].(string)
// s.Ext6, _ = m["ext_6"].(string)
// s.Ext7, _ = m["ext_7"].(string)
// s.Ext8, _ = m["ext_8"].(string)
// s.Ext9, _ = m["ext_9"].(string)
// s.Ext10, _ = m["ext_10"].(string)

// s.AllPeople, _ = m["all_people"].(int64)
// s.SignPeople, _ = m["sign_people"].(int64)
// s.ApplyPeople, _ = m["apply_people"].(int64)
// s.ReviewPeople, _ = m["review_people"].(int64)
// s.RefusePeople, _ = m["refuse_people"].(int64)
// s.NoPeople, _ = m["no_people"].(int64)

// log.Println("人數資料: ", s.AllPeople, s.SignPeople, s.ApplyPeople, s.ReviewPeople, s.RefusePeople, s.NoPeople)

// 用戶官方帳號綁定
// s.UserLineID, _ = m["user_line_id"].(string)
// s.UserChannelID, _ = m["user_channel_id"].(string)
// s.UserChannelSecret, _ = m["user_channel_secret"].(string)
// s.UserChatbotSecret, _ = m["user_chatbot_secret"].(string)
// s.UserChatbotToken, _ = m["user_chatbot_token"].(string)
// 活動官方帳號綁定
// s.ActivityLineID, _ = m["activity_line_id"].(string)
// s.ActivityChannelID, _ = m["activity_channel_id"].(string)
// s.ActivityChannelSecret, _ = m["activity_channel_secret"].(string)
// s.ActivityChatbotSecret, _ = m["activity_chatbot_secret"].(string)
// s.ActivityChatbotToken, _ = m["activity_chatbot_token"].(string)

// 活動資訊(join activity)
// s.ActivityUserID, _ = m["activity_user_id"].(string)
// s.ActivityName, _ = m["activity_name"].(string)
// s.StartTime, _ = m["start_time"].(string)
// s.EndTime, _ = m["end_time"].(string)
// s.SignCheck, _ = m["sign_check"].(string)
// s.SignAllow, _ = m["sign_allow"].(string)
// s.SignMinutes, _ = m["sign_minutes"].(int64)
// s.ActivityDevice, _ = m["activity_device"].(string)

// 活動總覽
// s.OverviewMessage, _ = m["overview_message"].(string)
// // s.OverviewTopic, _ = m["overview_topic"].(string)
// s.OverviewQuestion, _ = m["overview_question"].(string)
// s.OverviewDanmu, _ = m["overview_danmu"].(string)
// s.OverviewSpecialDanmu, _ = m["overview_special_danmu"].(string)
// // s.OverviewPicture, _ = m["overview_picture"].(string)
// s.OverviewHoldscreen, _ = m["overview_holdscreen"].(string)
// // s.OverviewGeneral, _ = m["overview_general"].(string)
// // s.OverviewThreed, _ = m["overview_threed"].(string)
// // s.OverviewCountdown, _ = m["overview_countdown"].(string)
// s.OverviewLottery, _ = m["overview_lottery"].(string)
// s.OverviewRedpack, _ = m["overview_redpack"].(string)
// s.OverviewRopepack, _ = m["overview_ropepack"].(string)
// s.OverviewWhackMole, _ = m["overview_whack_mole"].(string)
// s.OverviewMonopoly, _ = m["overview_monopoly"].(string)
// s.OverviewQA, _ = m["overview_qa"].(string)
// s.OverviewTugofwar, _ = m["overview_tugofwar"].(string)
// s.OverviewBingo, _ = m["overview_bingo"].(string)
// s.OverviewSignname, _ = m["overview_signname"].(string)
// // s.OverviewDrawNumbers, _ = m["overview_draw_numbers"].(string)
// s.Overview3DGachaMachine, _ = m["overview_3d_gacha_machine"].(string)
// s.OverviewVote, _ = m["overview_vote"].(string)

// 訊息牆
// s.MessageBan, _ = m["message_ban"].(string)
// s.MessageBanSecond, _ = m["message_ban_second"].(int64)

// 霸屏
// s.HoldscreenBirthdayTopic, _ = m["holdscreen_birthday_topic"].(string)

// 用戶資訊
// s.Name, _ = m["name"].(string)
// s.Avatar, _ = m["avatar"].(string)
// s.Phone, _ = m["phone"].(string)
// s.Email, _ = m["email"].(string)
// s.ExtEmail, _ = m["ext_email"].(string)
// s.Friend, _ = m["friend"].(string)
// s.Identify, _ = m["identify"].(string)
// s.UserDevice, _ = m["user_device"].(string)
// s.ExtPassword, _ = m["ext_password"].(string)

// 自定義
// s.Ext1Name, _ = m["ext_1_name"].(string)
// s.Ext1Type, _ = m["ext_1_type"].(string)
// s.Ext1Options, _ = m["ext_1_options"].(string)
// s.Ext1Required, _ = m["ext_1_required"].(string)

// s.Ext2Name, _ = m["ext_2_name"].(string)
// s.Ext2Type, _ = m["ext_2_type"].(string)
// s.Ext2Options, _ = m["ext_2_options"].(string)
// s.Ext2Required, _ = m["ext_2_required"].(string)

// s.Ext3Name, _ = m["ext_3_name"].(string)
// s.Ext3Type, _ = m["ext_3_type"].(string)
// s.Ext3Options, _ = m["ext_3_options"].(string)
// s.Ext3Required, _ = m["ext_3_required"].(string)

// s.Ext4Name, _ = m["ext_4_name"].(string)
// s.Ext4Type, _ = m["ext_4_type"].(string)
// s.Ext4Options, _ = m["ext_4_options"].(string)
// s.Ext4Required, _ = m["ext_4_required"].(string)

// s.Ext5Name, _ = m["ext_5_name"].(string)
// s.Ext5Type, _ = m["ext_5_type"].(string)
// s.Ext5Options, _ = m["ext_5_options"].(string)
// s.Ext5Required, _ = m["ext_5_required"].(string)

// s.Ext6Name, _ = m["ext_6_name"].(string)
// s.Ext6Type, _ = m["ext_6_type"].(string)
// s.Ext6Options, _ = m["ext_6_options"].(string)
// s.Ext6Required, _ = m["ext_6_required"].(string)

// s.Ext7Name, _ = m["ext_7_name"].(string)
// s.Ext7Type, _ = m["ext_7_type"].(string)
// s.Ext7Options, _ = m["ext_7_options"].(string)
// s.Ext7Required, _ = m["ext_7_required"].(string)

// s.Ext8Name, _ = m["ext_8_name"].(string)
// s.Ext8Type, _ = m["ext_8_type"].(string)
// s.Ext8Options, _ = m["ext_8_options"].(string)
// s.Ext8Required, _ = m["ext_8_required"].(string)

// s.Ext9Name, _ = m["ext_9_name"].(string)
// s.Ext9Type, _ = m["ext_9_type"].(string)
// s.Ext9Options, _ = m["ext_9_options"].(string)
// s.Ext9Required, _ = m["ext_9_required"].(string)

// s.Ext10Name, _ = m["ext_10_name"].(string)
// s.Ext10Type, _ = m["ext_10_type"].(string)
// s.Ext10Options, _ = m["ext_10_options"].(string)
// s.Ext10Required, _ = m["ext_10_required"].(string)

// s.Ext1Unique, _ = m["ext_1_unique"].(string)
// s.Ext2Unique, _ = m["ext_2_unique"].(string)
// s.Ext3Unique, _ = m["ext_3_unique"].(string)
// s.Ext4Unique, _ = m["ext_4_unique"].(string)
// s.Ext5Unique, _ = m["ext_5_unique"].(string)
// s.Ext6Unique, _ = m["ext_6_unique"].(string)
// s.Ext7Unique, _ = m["ext_7_unique"].(string)
// s.Ext8Unique, _ = m["ext_8_unique"].(string)
// s.Ext9Unique, _ = m["ext_9_unique"].(string)
// s.Ext10Unique, _ = m["ext_10_unique"].(string)

// s.ExtEmailRequired, _ = m["ext_email_required"].(string)
// s.ExtPhoneRequired, _ = m["ext_phone_required"].(string)
// s.InfoPicture, _ = m["info_picture"].(string)
// if s.InfoPicture != "" {
// 	s.InfoPicture = "/admin/uploads/" + s.ActivityUserID + "/" + s.ActivityID + "/applysign/customize/" + s.InfoPicture
// }

// if model.Status == "apply" { // 通過報名
// 	// %s 您好: 您已報名 %s 活動，請於活動當天依照主辦單位指示完成活動簽到
// 	message = "%s 您好: 您已報名 %s 活動，請於活動當天依照主辦單位指示完成活動簽到"
// }
// else if model.Status == "refuse" { // 拒絕報名
// 	message = "活動報名失敗，可能原因為審核資格不符，如有疑問，請聯繫主辦單位"
// }

// 取得用戶資料
// lineModel, err := DefaultLineModel().SetDbConn(s.DbConn).
// 	SetRedisConn(s.RedisConn).Find(true, "", "user_id", userID)
// if err != nil || lineModel.UserID == "" {
// 	return errors.New("錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
// }

// message = fmt.Sprintf("%s 您好: 您已簽到 %s 活動，可利用以下連結進入活動(避免資料洩漏，連結勿提供給他人使用): ",
// 	lineModel.Name, activityModel.ActivityName)

// if lineModel.Device == "line" && isPushMessage {
// 	url := fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, model.ActivityID) + "&sign=open"

// 	// 使用第三方的官方帳號才需傳送訊息
// 	if err = pushMessage(chatbotSecret, chatbotToken, userID, message+url); err != nil {
// 		return err
// 	}
// } else if lineModel.Device == "facebook" || lineModel.Device == "gmail" || lineModel.Device == "customize" {
// 	url := fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, model.ActivityID, userID) + "&sign=open"

// 	if lineModel.Phone != "" {
// 		// 傳送簡訊
// 		err = sendMessage(lineModel.Phone, message+url)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// FindUserApplysignStatus 查詢用戶報名簽到狀態
// func (s ApplysignModel) FindUserApplysignStatus(activity, user string) (ApplysignModel, error) {
// 	item, err := s.Table(s.Base.TableName).Where("activity_id", "=", activity).
// 		Where("user_id", "=", user).First()
// 	if err != nil {
// 		return ApplysignModel{}, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
// 	}
// 	if item == nil {
// 		return ApplysignModel{}, nil
// 	}

// 	return s.MapToModel(item), nil
// }

// FindApplysignLeftJoinActivityCustomize 查詢資料(join line_users、activity_customize)
// func (s ApplysignModel) FindApplysignLeftJoinActivityCustomize(activity, userID string) (ApplysignModel, error) {
// 	item, err := s.Table(s.Base.TableName).
// 		Select("activity_applysign.id", "activity_applysign.user_id",
// 			"activity_applysign.activity_id",

// 			"line_users.phone", "line_users.ext_email",

// 			"activity_applysign.ext_1", "activity_customize.ext_1_name", "activity_customize.ext_1_type",
// 			"activity_customize.ext_1_options", "activity_customize.ext_1_required",
// 			"activity_applysign.ext_2", "activity_customize.ext_2_name", "activity_customize.ext_2_type",
// 			"activity_customize.ext_2_options", "activity_customize.ext_2_required",
// 			"activity_applysign.ext_3", "activity_customize.ext_3_name", "activity_customize.ext_3_type",
// 			"activity_customize.ext_3_options", "activity_customize.ext_3_required",
// 			"activity_applysign.ext_4", "activity_customize.ext_4_name", "activity_customize.ext_4_type",
// 			"activity_customize.ext_4_options", "activity_customize.ext_4_required",
// 			"activity_applysign.ext_5", "activity_customize.ext_5_name", "activity_customize.ext_5_type",
// 			"activity_customize.ext_5_options", "activity_customize.ext_5_required",
// 			"activity_applysign.ext_6", "activity_customize.ext_6_name", "activity_customize.ext_6_type",
// 			"activity_customize.ext_6_options", "activity_customize.ext_6_required",
// 			"activity_applysign.ext_7", "activity_customize.ext_7_name", "activity_customize.ext_7_type",
// 			"activity_customize.ext_7_options", "activity_customize.ext_7_required",
// 			"activity_applysign.ext_8", "activity_customize.ext_8_name", "activity_customize.ext_8_type",
// 			"activity_customize.ext_8_options", "activity_customize.ext_8_required",
// 			"activity_applysign.ext_9", "activity_customize.ext_9_name", "activity_customize.ext_9_type",
// 			"activity_customize.ext_9_options", "activity_customize.ext_9_required",
// 			"activity_applysign.ext_10", "activity_customize.ext_10_name", "activity_customize.ext_10_type",
// 			"activity_customize.ext_10_options", "activity_customize.ext_10_required",

// 			"activity_customize.ext_email_required", "activity_customize.ext_phone_required",
// 			"activity_customize.info_picture").
// 		LeftJoin(command.Join{
// 			FieldA: "activity_applysign.user_id",
// 			// FieldA:    "activity_applysign.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_applysign.activity_id",
// 			FieldA1:   "activity_customize.activity_id",
// 			Table:     "activity_customize",
// 			Operation: "=",
// 		}).
// 		Where("activity_applysign.activity_id", "=", activity).
// 		Where("line_users.user_id", "=", userID).First()
// 	if err != nil {
// 		return ApplysignModel{}, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
// 	}
// 	if item == nil {
// 		return ApplysignModel{}, nil
// 	}
// 	return s.MapToModel(item), nil
// }

// FindByID 查詢資料
// func (s ApplysignModel) FindByID(id int64) (ApplysignModel, error) {
// 	item, err := s.Table(s.Base.TableName).
// 		Select("id", "user_id", "activity_id").
// 		Where("id", "=", id).First()
// 	if err != nil {
// 		return ApplysignModel{}, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
// 	}
// 	if item == nil {
// 		return ApplysignModel{}, nil
// 	}
// 	return s.MapToModel(item), nil
// }

// // FindByIDs 查詢資料(多個id)
// func (s ApplysignModel) FindByIDs(ids []interface{}) ([]ApplysignModel, error) {
// 	items, err := s.Table(s.Base.TableName).
// 		Select("id", "user_id", "activity_id").
// 		WhereIn("id", ids).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得簽到人員資訊，請重新查詢")
// 	}

// 	return MapToApplysignModel(items), nil
// }

// FindUserAllApplysignStatus 查詢用戶所有活動報名簽到狀態
// func (s ApplysignModel) FindUserAllApplysignStatus(userID string) ([]ApplysignModel, error) {
// 	var (
// 	// item   map[string]interface{}
// 	// err    error
// 	// now, _ = time.ParseInLocation("2006-01-02 15:04:05",
// 	// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// 	)

// 	items, err := s.Table(s.Base.TableName).
// 		Select("activity_applysign.id", "activity_applysign.user_id",
// 			"activity_applysign.activity_id", "activity_applysign.status",
// 			"activity_applysign.number",

// 			// activity
// 			"activity.user_id as activity_user_id", "activity.activity_name",
// 			"activity.start_time", "activity.end_time").
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_applysign.activity_id",
// 			FieldA1:    "activity.activity_id",
// 			Table:     "activity",
// 			Operation: "=",
// 		}).
// 		Where("activity_applysign.user_id", "=", userID).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得用戶所有活動報名簽到狀態資訊，請重新查詢")
// 	}

// 	return MapToApplysignModel(items), nil
// }

// LINE相關資訊
// if activityModel.ChatbotSecret != "" && activityModel.ChatbotToken != "" {
// 	chatbotSecret = activityModel.ChatbotSecret
// 	chatbotToken = activityModel.ChatbotToken
// 	// isPushMessage = true
// } else {
// 	chatbotSecret = config.CHATBOT_SECRET
// 	chatbotToken = config.CHATBOT_TOKEN
// }

// var number int
// 已簽到人數
// attend, err := DefaultApplysignModel().SetDbConn(s.DbConn).GetAttend(model.ActivityID)
// if err != nil {
// 	return err
// }
// attend++ // 簽到人數+1
// number++ // 簽到完成時發放號碼並遞增number資料
// attend-- // 取消簽到人數-1
// 更新活動人數
// if err = activityModel.UpdateAttend(model.ActivityID, attend); err != nil {
// 	return err
// }
// url := fmt.Sprintf(liffURL, model.ActivityID)
// 遞增抽獎號碼紀錄
// if err = activityModel.IncrNumber(model.ActivityID); err != nil {
// 	return err
// }

// "name":        model.Name,
// "avatar":      model.Avatar,
// "email":       model.Email,
// "phone":       model.Phone,
// "ext_email":   model.ExtEmail,
// "apply_time":  now,

// status = []string{"review", "apply", "sign", "refuse", "cancel"}
// now, _ = time.ParseInLocation("2006-01-02 15:04:05",
// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)

// ID         string `json:"id"`
// Name       string `json:"name"`
// Avatar     string `json:"avatar"`
// ApplyTime  string `json:"apply_time"`
// ReviewTime string `json:"review_time"`
// SignTime   string `json:"sign_time"`
// Email    string `json:"email"`
// ExtEmail string `json:"ext_email"`
// Phone    string `json:"phone"`
// Name       string `json:"name"`
// Avatar     string `json:"avatar"`
// ApplyTime  string `json:"apply_time"`
// Phone      string `json:"phone"`
// Email      string `json:"email"`
// ExtEmail   string `json:"ext_email"`
// Ext1       string `json:"ext_1"`
// Ext2       string `json:"ext_2"`
// Ext3       string `json:"ext_3"`
// Ext4       string `json:"ext_4"`
// Ext5       string `json:"ext_5"`
// Ext6       string `json:"ext_6"`
// Ext7       string `json:"ext_7"`
// Ext8       string `json:"ext_8"`
// Ext9       string `json:"ext_9"`
// Ext10      string `json:"ext_10"`

// UpdateUser 更新報名簽到人員姓名、頭像資料
// func (s ApplysignModel) UpdateUser(userID, name, avatar, email, extEmail, phone string) error {
// 	var (
// 		fields      = []string{"phone", "ext_email"}
// 		values      = []string{phone, extEmail}
// 		fieldValues = command.Value{"name": name, "avatar": avatar, "email": email}
// 	)
// 	for i, field := range fields {
// 		if values[i] != "" {
// 			fieldValues[field] = values[i]
// 		}
// 	}

// 	if err := s.Tabl
