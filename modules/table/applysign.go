package table

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"hilive/template/form"
	"hilive/template/types"

	"github.com/gin-gonic/gin"
)

// GetApplysignUsersPanel 自定義匯入報名簽到人員(匹量匯入)
func (s *SystemTable) GetApplysignUsersPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.LINE_USERS_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "name") {
			return errors.New("錯誤: 資料發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.ApplysignUserModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultApplysignUserModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Adds(true, model); err != nil {
			return err
		}
		return nil
	})

	return
}

// GetApplysignUserPanel 自定義匯入報名簽到人員(單筆匯入)
func (s *SystemTable) GetApplysignUserPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()

	info.SetTable(config.LINE_USERS_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		// delete data(line_user表)
		if err := s.table(config.LINE_USERS_TABLE).
			Where("activity_id", "=", activityID).
			WhereIn("user_id", interfaces(idArr)).Delete(); err != nil {
			return err
		}

		// 刪除用戶redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.AUTH_USERS_REDIS + id) // 簽到人員資訊
		}

		return nil
	})

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.LINE_USERS_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "name") {
			return errors.New("錯誤: 資料發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.ApplysignUserModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if _, err := models.DefaultApplysignUserModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, model); err != nil {
			return err
		}
		return nil
	})

	return
}

// 更新自定義人員資料
// formList.SetUpdateFunc(func(values form2.Values) error {
// 	if values.IsEmpty("applysign_user_id", "activity_id", "ext_password") {
// 		return errors.New("錯誤: 資料發生問題，請輸入有效的ID")
// 	}

// 	if err := models.DefaultApplysignUserModel().SetDbConn(s.dbConn).
// 		SetRedisConn(s.redisConn).Update(
// 		models.ApplysignUserModel{
// 			ActivityID: values.Get("activity_id"),
// 			UserID:     values.Get("applysign_user_id"),
// 			Name:       values.Get("name"),
// 			Phone:      values.Get("phone"),
// 			ExtEmail:   values.Get("ext_email"),
// 			// ExtAccount:  values.Get("ext_account"),
// 			ExtPassword: values.Get("ext_password"),
// 		}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}

// 	// 刪除舊的用戶資訊
// 	s.redisConn.DelCache(config.AUTH_USERS_REDIS + values.Get("applysign_user_id"))
// 	return nil
// })

// GetApplysignPanel 報名簽到人員
func (s *SystemTable) GetApplysignPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊
	info := table.GetInfo()
	info.AddField("ID", "id", "INT").FieldHide()
	info.AddField("姓名", "name", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: config.LINE_USERS_TABLE, JoinField: "user_id",
			BaseTable: config.ACTIVITY_APPLYSIGN_TABLE, Field: "user_id",
		})
	info.AddField("頭像", "avatar", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: config.LINE_USERS_TABLE, JoinField: "user_id",
			BaseTable: config.ACTIVITY_APPLYSIGN_TABLE, Field: "user_id",
		})
	info.AddField("簽到狀態", "status", db.Varchar)

	info.AddField("報名時間", "apply_time", db.Datetime)
	info.AddField("審核時間", "review_time", db.Datetime)
	info.AddField("簽到時間", "sign_time", db.Datetime)
	info.AddField("電話", "phone", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: config.LINE_USERS_TABLE, JoinField: "user_id",
			BaseTable: config.ACTIVITY_APPLYSIGN_TABLE, Field: "user_id",
		})
	info.AddField("電子信箱", "email", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: config.LINE_USERS_TABLE, JoinField: "user_id",
			BaseTable: config.ACTIVITY_APPLYSIGN_TABLE, Field: "user_id",
		})
	info.AddField("自定義電子信箱", "ext_email", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: config.LINE_USERS_TABLE, JoinField: "user_id",
			BaseTable: config.ACTIVITY_APPLYSIGN_TABLE, Field: "user_id",
		})
	info.AddField("ext_1", "ext_1", db.Varchar)
	info.AddField("ext_2", "ext_2", db.Varchar)
	info.AddField("ext_3", "ext_3", db.Varchar)
	info.AddField("ext_4", "ext_4", db.Varchar)
	info.AddField("ext_5", "ext_5", db.Varchar)
	info.AddField("ext_6", "ext_6", db.Varchar)
	info.AddField("ext_7", "ext_7", db.Varchar)
	info.AddField("ext_8", "ext_8", db.Varchar)
	info.AddField("ext_9", "ext_9", db.Varchar)
	info.AddField("ext_10", "ext_10", db.Varchar)

	info.SetTable(config.ACTIVITY_APPLYSIGN_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var (
			applySignModel models.ApplysignModel
			err            error
		)

		// 取得activity_id
		if len(idArr) > 0 {
			id, err := strconv.Atoi(idArr[0])
			if err != nil {
				return errors.New("錯誤: ID發生問題，請輸入有效的ID")
			}

			applySignModel, err = models.DefaultApplysignModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Find(int64(id), "", "", false)
			if err != nil || applySignModel.ID == 0 {
				return errors.New("錯誤: ID發生問題(查詢報名簽到資料)，請輸入有效的ID")
			}
		}

		// delete data(報名簽到資料表)
		if err := s.table(config.ACTIVITY_APPLYSIGN_TABLE).
			WhereIn("id", interfaces(idArr)).Delete(); err != nil {
			return err
		}

		// 已簽到人數
		attend, err := models.DefaultApplysignModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			GetAttend(applySignModel.ActivityID)
		if err != nil {
			return nil
		}

		// 更新活動人數
		if err = models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateAttendAndNumber(false, applySignModel.ActivityID, attend, 0, []string{}); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}

		// 刪除redis活動相關資訊
		// s.redisConn.DelCache(config.ACTIVITY_REDIS + applySignModel.ActivityID)    // 活動資訊
		// s.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + applySignModel.ActivityID) // 簽到人員資訊
		s.redisConn.DelCache(config.SIGN_STAFFS_2_REDIS + applySignModel.ActivityID) // 簽到人員資訊

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+applySignModel.ActivityID, "刪除人員")

		return nil
	})

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_APPLYSIGN_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("line_users", "activity_id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 拿掉陣列參數(避免json解碼發生問題)
		delete(flattened, "line_users")

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditApplysignModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.LineUsers = strings.Split(values.Get("line_users"), ",")

		if model.Status != "" {
			// 更新報名簽到狀態資料
			if err := models.DefaultApplysignModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				UpdateStatus(true, values.Get("host"), model, false); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
		} else if model.Role != "" {
			// 更新角色資料
			if err := models.DefaultApplysignModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				UpdateRole(true, model); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
		}

		return nil
	})
	return
}

// GetApplyPanel 報名
func (s *SystemTable) GetApplyPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 基本設置表單資訊
	settingFormList := table.GetSettingForm()
	settingFormList.AddField("報名審核", "apply_check", db.Varchar, form.Checkbox)

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
		var model models.EditApplyModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		var (
			avatar string
		)

		if values.Get("customize_default_avatar"+DEFAULT_FALG) == "1" {
			avatar = fmt.Sprintf("%s/admin/uploads/system/img-user-pic.png", config.HTTP_HILIVES_NET_URL)
		} else if values.Get("customize_default_avatar") != "" {
			avatar = fmt.Sprintf("%s/admin/uploads/%s/%s/applysign/apply/%s",
				config.HTTP_HILIVES_NET_URL, values.Get("user_id"), values.Get("activity_id"), values.Get("customize_default_avatar"))
		}

		// 手動處理
		model.CustomizeDefaultAvatar = avatar

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateApply(model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}

		return nil
	})
	return
}

// GetSignPanel 簽到
func (s *SystemTable) GetSignPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 基本設置表單資訊
	settingFormList := table.GetSettingForm()
	settingFormList.AddField("簽到時需要掃描QRcode審核", "sign_check", db.Varchar, form.Checkbox)
	settingFormList.AddField("允許讓參加人員即刻簽到", "sign_allow", db.Varchar, form.Checkbox)
	settingFormList.AddField("活動開始前", "sign_minutes", db.Int, form.Text)
	settingFormList.AddField("手動設置簽到狀態", "sign_manual", db.Varchar, form.Checkbox)

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
		var model models.EditSignModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateSign(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetCustomizePanel 自定義欄位設置
func (s *SystemTable) GetCustomizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 基本設置表單資訊
	settingFormList := table.GetSettingForm()

	settingFormList.SetTable(config.ACTIVITY_CUSTOMIZE_TABLE).SetUpdateFunc(func(values form2.Values) error {
		// fmt.Println("有問題: ", values.Get("activity_id"))
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
		var model models.EditCustomizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultCustomizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}

		return nil
	})
	return
}

// GetQRcodePanel QRcode自定義
func (s *SystemTable) GetQRcodePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 基本設置表單資訊
	settingFormList := table.GetSettingForm()

	settingFormList.SetTable(config.ACTIVITY_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		var (
			logoPicture string
		)
		if values.Get("qrcode_logo_picture"+DEFAULT_FALG) == "1" {
			logoPicture = UPLOAD_SYSTEM_URL + "qr-logo.svg"
		} else if values.Get("qrcode_logo_picture") != "" {
			logoPicture = values.Get("qrcode_logo_picture")
		} else {
			logoPicture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditQRcodeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.QRcodeLogoPicture = logoPicture

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateQRcode(model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// models.EditQRcodeModel{
// 	ActivityID:            values.Get("activity_id"),
// 	QRcodeLogoPicture:     logoPicture,
// 	QRcodeLogoSize:        values.Get("qrcode_logo_size"),
// 	QRcodePicturePoint:    values.Get("qrcode_picture_point"),
// 	QRcodeWhiteDistance:   values.Get("qrcode_white_distance"),
// 	QRcodePointColor:      values.Get("qrcode_point_color"),
// 	QRcodeBackgroundColor: values.Get("qrcode_background_color"),
// }

// @Summary 新增單筆自定義報名簽到人員資料
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @@@param applysign_user_id formData string true "user_id"
// @param name formData string true "name"
// @param phone formData string false "phone"
// @param ext_email formData string false "ext_email"
// @param ext_password formData string false "ext_password"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/user [post]
func POSTApplysignUser(ctx *gin.Context) {
}

// @Summary 新增多筆自定義報名簽到人員資料
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param people formData integer true "people"
// @@@param applysign_user_id formData string true "user_id"
// @param name formData string true "name，用逗點間隔，陣列長度需與people一樣"
// @param phone formData string false "phone，用逗點間隔，陣列長度需與people一樣"
// @param ext_email formData string false "ext_email，用逗點間隔，陣列長度需與people一樣"
// @param ext_password formData string false "ext_password，用逗點間隔，陣列長度需與people一樣"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/users [post]
func POSTApplysignUsers(ctx *gin.Context) {
}

// @Summary 編輯報名簽到資料
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @@@param ids formData string true "ID(以,區隔多筆ID)"
// @param line_users formData string true "LINE ID(以,區隔多筆LINE ID)"
// @param activity_id formData string true "活動ID"
// @param host formData string false "網域" Enums(hilives.net, www.hilives.net, dev.hilives.net)
// @param status formData string false "狀態" Enums(review, apply, sign, refuse, cancel)
// @param role formData string false "role" Enums(admin, guest)
// @param review_time formData string false "審核時間(西元年-月-日 時:分)"
// @param sign_time formData string false "簽到時間(西元年-月-日 時:分)"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign [put]
func PUTApplysign(ctx *gin.Context) {
}

// @Summary 報名牆設置
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param apply_check formData string false "報名審核" Enums(open, close)
// @param customize_password formData string false "是否自定義設置驗整碼" Enums(open, close)
// @param allow_customize_apply formData string false "是否允許用戶自定義報名" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/apply [patch]
func PATCHApply(ctx *gin.Context) {
}

// @Summary 簽到牆設置
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param sign_check formData string false "簽到時需要掃描QRcode審核" Enums(open, close)
// @param sign_allow formData string false "允許讓參加人員即刻簽到" Enums(open, close)
// @param sign_minutes formData integer false "活動開始前n分鐘"
// @param sign_manual formData string false "手動設置簽到狀態" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/sign [patch]
func PATCHSign(ctx *gin.Context) {
}

// @Summary 自定義欄位設置
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param ext_email_required formData string false "是否必填電子信箱" Enums(true, false)
// @param ext_phone_required formData string false "是否必填電話號碼" Enums(true, false)
// @param info_picture formData file false "表單資訊圖"
// @param field formData string true "ext_n" Enums(ext_1,ext_2,ext_3,ext_4,ext_5,ext_6,ext_7,ext_8,ext_9,ext_10)
// @param name formData string true "ext_n_name"
// @param type formData string true "ext_n_type" Enums(text, radio, select, checkbox, textarea, date, time)
// @param options formData string true "ext_n_options，用分行區隔選單"
// @param required formData string true "ext_n_required" Enums(true, false)
// @param unique formData string true "ext_n_unique" Enums(true, false)
// @param is_delete formData string true "是否清空欄位資料" Enums(true, false)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/customize [patch]
func PATCHCustomize(ctx *gin.Context) {
}

// @Summary QRcode自定義欄位設置
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param qrcode_logo_picture formData file false "LOGO圖片"
// @param qrcode_logo_size formData number false "LOGO尺寸" minimum(0.2) mxnimum(0.5)
// @param qrcode_picture_point formData string false "LOGO圖片後的點" Enums(open, close)
// @param qrcode_white_distance formData integer false "LOGO留白距離" minimum(0) mxnimum(80)
// @param qrcode_point_color formData string false "QR點的顏色"
// @param qrcode_background_color formData string false "背景顏色"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/qrcode [patch]
func PATCHQRcode(ctx *gin.Context) {
}

// @Summary 刪除報名簽到資料
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign [delete]
func DELETEApplysign(ctx *gin.Context) {
}

// @Summary 刪除自定義報名簽到人員資料
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/user [delete]
func DELETEApplysignUser(ctx *gin.Context) {
}

// @Summary 報名簽到人員JSON資料
// @Tags ApplySign
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param user_id query string false "用戶ID(用戶所有報名簽到狀態)"
// @param status query string false "報名簽到狀態(用,間隔)" Enums(sign, apply, review, cancel, refuse, no)
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign [get]
func ApplySignJSON(ctx *gin.Context) {
}

// @Summary 報名簽到自定義JSON資料
// @Tags ApplySign
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/customize [get]
func ApplySignCustomizeJSON(ctx *gin.Context) {
}

// models.ApplysignUserModel{
// 	ActivityID: values.Get("activity_id"),
// 	People:     values.Get("people"),
// 	Name:       values.Get("name"),
// 	Phone:      values.Get("phone"),
// 	ExtEmail:   values.Get("ext_email"),

// 	Ext1:  values.Get("ext_1"),
// 	Ext2:  values.Get("ext_2"),
// 	Ext3:  values.Get("ext_3"),
// 	Ext4:  values.Get("ext_4"),
// 	Ext5:  values.Get("ext_5"),
// 	Ext6:  values.Get("ext_6"),
// 	Ext7:  values.Get("ext_7"),
// 	Ext8:  values.Get("ext_8"),
// 	Ext9:  values.Get("ext_9"),
// 	Ext10: values.Get("ext_10"),

// 	ExtPassword: values.Get("ext_password"),
// 	Source:      "excel",
// }

// delete data(報名簽到資料表)
// if err := s.table(config.ACTIVITY_APPLYSIGN_TABLE).
// 	Where("activity_id", "=", activityID).
// 	WhereIn("user_id", interfaces(idArr)).Delete(); err != nil {
// 	// 有可能該用戶只有被匯入但還沒報名簽到，所以刪除時不會有資料變化
// 	// &&err.Error() != "錯誤: 無刪除任何資料，請重新操作"
// 	return err
// }

// 已簽到人數
// attend, err := models.DefaultApplysignModel().
// 	SetDbConn(s.dbConn).GetAttend(activityID)
// if err != nil {
// 	return nil
// }

// 更新活動人數
// if err = models.DefaultActivityModel().SetDbConn(s.dbConn).
// 	UpdateAttend(activityID, attend); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return err
// }

// 刪除redis活動相關資訊
// s.redisConn.DelCache(config.ACTIVITY_REDIS + applySignModel.ActivityID)    // 活動資訊
// s.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + activityID) // 簽到人員資訊
// s.redisConn.DelCache(config.SIGN_STAFFS_2_REDIS + activityID) // 簽到人員資訊

// models.ApplysignUserModel{
// 	ActivityID: values.Get("activity_id"),
// 	// UserID:     values.Get("applysign_user_id"),
// 	Name:     values.Get("name"),
// 	Phone:    values.Get("phone"),
// 	ExtEmail: values.Get("ext_email"),

// 	Ext1:  values.Get("ext_1"),
// 	Ext2:  values.Get("ext_2"),
// 	Ext3:  values.Get("ext_3"),
// 	Ext4:  values.Get("ext_4"),
// 	Ext5:  values.Get("ext_5"),
// 	Ext6:  values.Get("ext_6"),
// 	Ext7:  values.Get("ext_7"),
// 	Ext8:  values.Get("ext_8"),
// 	Ext9:  values.Get("ext_9"),
// 	Ext10: values.Get("ext_10"),

// 	ExtPassword: values.Get("ext_password"),
// 	// Source:      "excel",
// }

// formList.SetTable(config.ACTIVITY_APPLYSIGN_TABLE).SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "ext_password") {
// 		return errors.New("錯誤: 資料發生問題，請輸入有效的ID")
// 	}

// 	if err := models.DefaultApplysignUserModel().SetDbConn(s.dbConn).Add(
// 		models.ApplysignUserModel{
// 			ActivityID: values.Get("activity_id"),
// 			Name:       values.Get("name"),
// 			Phone:      values.Get("phone"),
// 			ExtEmail:   values.Get("ext_email"),
// 			// ExtAccount:  values.Get("ext_account"),
// 			ExtPassword: values.Get("ext_password"),
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// })

// models.EditApplysignModel{
// 	ActivityID: values.Get("activity_id"),
// 	LineUsers:  strings.Split(values.Get("line_users"), ","),
// 	Status:     values.Get("status"),
// 	ReviewTime: values.Get("review_time"),
// 	SignTime:   values.Get("sign_time"),
// }

// models.EditApplyModel{
// 	ActivityID:          values.Get("activity_id"),
// 	ApplyCheck:          values.Get("apply_check"),
// 	CustomizePassword:   values.Get("customize_password"),
// 	AllowCustomizeApply: values.Get("allow_customize_apply"),
// }

// models.EditSignModel{
// 	ActivityID:  values.Get("activity_id"),
// 	SignCheck:   values.Get("sign_check"),
// 	SignAllow:   values.Get("sign_allow"),
// 	SignMinutes: values.Get("sign_minutes"),
// 	SignManual:  values.Get("sign_manual"),
// }

// models.EditCustomizeModel{
// 	ActivityID:             values.Get("activity_id"),
// 	Field:                  values.Get("field"),
// 	Name:                   values.Get("name"),
// 	Type:                   values.Get("type"),
// 	Options:                values.Get("options"),
// 	Required:               values.Get("required"),
// 	Unique:                 values.Get("unique"),
// 	ExtEmailRequired:       values.Get("ext_email_required"),
// 	ExtPhoneRequired:       values.Get("ext_phone_required"),
// 	InfoPicture:            infoPicture,
// 	InfoPictureDefaultFlag: values.Get("info_picture" + DEFAULT_FALG),
// 	IsDelete:               values.Get("is_delete"),
// }

// var (
// 	infoPicture string
// )

// if values.Get("info_picture"+DEFAULT_FALG) == "1" {
// 	infoPicture = "" // 清空
// } else
// if values.Get("info_picture") != "" {
// 	infoPicture = values.Get("info_picture")
// }
