package table

import (
	"encoding/json"
	"errors"
	"fmt"

	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// GetLineUserPanel 報名簽到用戶
func (s *SystemTable) GetLineUserPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.LINE_USERS_TABLE)

	formList := table.GetForm()

	formList.SetTable(config.LINE_USERS_TABLE).SetUpdateFunc(func(values form2.Values) error {
		// log.Println("更新報名簽到用戶姓名頭像API，is_modify欄位也會更新為yes: ")

		if values.IsEmpty("user_id") {
			return errors.New("錯誤: 用戶ID發生問題，請輸入有效的資料")
		}

		var avatar string
		if values.Get("avatar"+DEFAULT_FALG) == "1" {
			avatar = fmt.Sprintf("%s/admin/uploads/system/img-user-pic.png", config.HTTP_HILIVES_NET_URL)
		} else if values.Get("avatar") != "" {
			avatar = fmt.Sprintf("%s/admin/uploads/line_user/%s",
				config.HTTP_HILIVES_NET_URL, values.Get("avatar"))
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.LineModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.Avatar = avatar
		model.IsModify = "yes"

		if err := models.DefaultLineModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateUser(model); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetUserPanel 用戶
func (s *SystemTable) GetUserPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.USERS_TABLE)

	formList := table.GetForm()
	formList.SetTable(config.USERS_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("name", "phone", "email", "password", "password_again") {
			return errors.New("錯誤: 所有欄位不能為空，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditUserModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.Permissions = "3"
		model.Avatar = UPLOAD_SYSTEM_URL + "img-user-pic.png"
		// model.Bind = "no"
		model.Cookie = "no"
		model.LineBind = "no"
		model.FbBind = "no"
		model.GmailBind = "no"

		// 註冊用
		if _, err := models.DefaultUserModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetTable(config.USERS_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("phone") {
			return errors.New("錯誤: 電話資料發生問題，請輸入有效的電話")
		}

		var avatar string
		if values.Get("avatar"+DEFAULT_FALG) == "1" {
			avatar = UPLOAD_SYSTEM_URL + "img-user-pic.png"
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
		var model models.EditUserModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數
		model.Avatar = avatar

		if err := models.DefaultUserModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, model,
				"users.phone", values.Get("phone")); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// models.EditUserModel{
// 	UserID:        values.Get("user"),
// 	Name:          values.Get("name"),
// 	Phone:         values.Get("phone"),
// 	Email:         values.Get("email"),
// 	Password:      values.Get("password"),
// 	PasswordAgain: values.Get("password_again"),
// 	Avatar:        avatar,
// 	Bind:          values.Get("bind"),
// 	Cookie:        values.Get("cookie"),

// 	// 官方帳號綁定
// 	LineID:        values.Get("line_id"),
// 	ChannelID:     values.Get("channel_id"),
// 	ChannelSecret: values.Get("channel_secret"),
// 	ChatbotSecret: values.Get("chatbot_secret"),
// 	ChatbotToken:  values.Get("chatbot_token"),
// }

// models.EditUserModel{
// 	Name:          values.Get("name"),
// 	Phone:         values.Get("phone"),
// 	Email:         values.Get("email"),
// 	Password:      values.Get("password"),
// 	PasswordAgain: values.Get("password_again"),
// 	Permissions:   "3",
// 	Avatar:        UPLOAD_SYSTEM_URL + "img-user-pic.png",
// 	Bind:          "no",
// 	Cookie:        "no",
// 	// MaxActivity:       "1",
// 	// MaxActivityPeople: "10",
// 	// MaxGamePeople:     "10",

// 	// 官方帳號綁定
// 	LineID:        "",
// 	ChannelID:     "",
// 	ChannelSecret: "",
// 	ChatbotSecret: "",
// 	ChatbotToken:  "",
// }

// var avatar string
// if values.Get("avatar") != "" {
// 	avatar =  values.Get("avatar")
// } else {
// avatar = UPLOAD_SYSTEM_URL + "img-user-pic.png"
// }

// @Summary 編輯報名簽到用戶
// @Tags User
// @version 1.0
// @Accept  mpfd
// @param user_id formData string false "用戶ID"
// @param name formData string false "用戶名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param avatar formData string false "頭像"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /line_user [put]
func PUTLineUser(ctx *gin.Context) {
}

// @Summary 新增用戶(註冊用)
// @Tags User
// @version 1.0
// @Accept  mpfd
// @param name formData string true "用戶名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param phone formData string true "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string true "電子信箱(電子郵件地址中必須包含「@」)"
// @param password formData string true "密碼"
// @param password_again formData string true "再次輸入密碼"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /user [post]
func POSTUser() {
}

// @Summary 編輯用戶
// @Tags User
// @version 1.0
// @Accept  mpfd
// @param user formData string false "用戶ID(編輯的用戶ID)"
// @param phone formData string true "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param name formData string false "用戶名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param email formData string false "電子信箱(請在郵件地址中包含「@」)"
// @param avatar formData file false "頭像"
// @param password formData string false "密碼"
// @param password_again formData string false "再次輸入密碼"
// @param bind formData string false "綁定" Enums(yes, no)
// @param cookie formData string false "cookie訊息" Enums(yes, no)
// @param line_id formData string false "line_id"
// @param channel_id formData string false "channel_id"
// @param channel_secret formData string false "channel_secret"
// @param chatbot_secret formData string false "chatbot_secret"
// @param chatbot_token formData string false "chatbot_token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /user [put]
func PUTUser(ctx *gin.Context) {
}

// models.LineModel{
// 	UserID:   values.Get("user_id"),
// 	Name:     values.Get("name"),
// 	Avatar:   avatar,
// 	IsModify: "yes",
// }
