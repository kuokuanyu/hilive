package table

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAdminManagerPanel 用戶
func (s *SystemTable) GetAdminManagerPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 面板資訊
	info := table.GetInfo()

	info.SetTable(config.USERS_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var (
				activitys = make([]interface{}, 0)
				// users     = make([]interface{}, 0)
				tables = []string{
					config.ACTIVITY_PERMISSIONS_TABLE, // 活動權限
					config.ACTIVITY_TABLE, config.ACTIVITY_2_TABLE,
					// config.ACTIVITY_CHANNEL_TABLE,
					config.ACTIVITY_INTRODUCE_TABLE,
					config.ACTIVITY_SCHEDULE_TABLE, config.ACTIVITY_GUEST_TABLE, config.ACTIVITY_MATERIAL_TABLE,
					config.ACTIVITY_QUESTION_GUEST_TABLE,
					config.ACTIVITY_STAFF_GAME_TABLE, config.ACTIVITY_STAFF_PRIZE_TABLE,
					config.ACTIVITY_STAFF_BLACK_TABLE, config.ACTIVITY_STAFF_PK_TABLE,
					// config.ACTIVITY_GAME_TABLE,
					config.ACTIVITY_PRIZE_TABLE,
					config.ACTIVITY_APPLYSIGN_TABLE, config.ACTIVITY_CHATROOM_RECORD_TABLE,
					config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE, config.ACTIVITY_QUESTION_USER_TABLE,
					config.ACTIVITY_SCORE_TABLE, config.ACTIVITY_CUSTOMIZE_TABLE, config.ACTIVITY_GAME_SETTING_TABLE,
					config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE,
					// 簽名牆
					config.ACTIVITY_SIGNNAME_TABLE} // 需要刪除的資料表
			)

			// 查詢所有用戶下的所有活動
			items, err := s.table(config.ACTIVITY_TABLE).WhereIn("user_id", interfaces(idArr)).All()
			if err != nil {
				return errors.New("錯誤: 查詢所有用戶創建的活動資訊發生問題")
			}
			for _, item := range items {
				activityID, _ := item["activity_id"].(string)
				activitys = append(activitys, activityID)
			}

			// 刪除活動資訊
			if len(activitys) > 0 {
				for _, t := range tables {
					s.table(t).WhereIn("activity_id", activitys).Delete()
				}
			}

			mongoTables := []string{
				config.ACTIVITY_CHANNEL_TABLE, // 活動頻道
				config.ACTIVITY_GAME_TABLE,    // 遊戲資訊
			}
			for _, t := range mongoTables {
				s.mongoConn.DeleteMany(t, bson.M{"activity_id": bson.M{"$in": activitys}})
			}

			// 刪除用戶(先把用戶刪除會無法查詢用戶資訊)
			if err := s.table(config.USERS_TABLE).
				WhereIn("user_id", interfaces(idArr)).Delete(); err != nil {
				return err
			}

			// 刪除用戶權限資料
			if err := s.table(config.USER_PERMISSIONS_TABLE).
				WhereIn("user_id", interfaces(idArr)).Delete(); err != nil {
				return err
			}

			for _, userID := range idArr {
				// 刪除用戶資料夾
				os.RemoveAll(config.STORE_PATH + "/" + userID)
			}

			return nil
		})

	// 表單資訊
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
		model.Avatar = UPLOAD_SYSTEM_URL + "img-user-pic.png"
		// model.Bind = "no"
		model.Cookie = "no"
		model.LineBind = "no"
		model.FbBind = "no"
		model.GmailBind = "no"

		if _, err := models.DefaultUserModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("user") {
			return errors.New("錯誤: 用戶ID發生問題，請輸入有效的用戶ID")
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
				"users.user_id", values.Get("user")); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})

	return
}

// GetAdminPermissionPanel 權限
func (s *SystemTable) GetAdminPermissionPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 面板資訊
	info := table.GetInfo()

	info.SetTable(config.PERMISSION_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			// 刪除權限資料
			if err := s.table(config.PERMISSION_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}

			// 刪除用戶權限表裡的資料
			if err := s.table(config.USER_PERMISSIONS_TABLE).
				WhereIn("permission_id", interfaces(idArr)).Delete(); err != nil &&
				err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.PERMISSION_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("permission") {
			return errors.New("錯誤: 權限名稱欄位不能為空，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPermissionModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultPermissionModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPermissionModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.
			DefaultPermissionModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil {
			return err
		}
		return nil
	})

	return
}

// GetAdminMenuPanel 菜單
func (s *SystemTable) GetAdminMenuPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 面板資訊
	info := table.GetInfo()

	info.SetTable(config.MENU_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.MENU_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.MENU_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("parent_id", "title", "sidebarbtnl") {
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
		var model models.EditMenuModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultMenuModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditMenuModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultMenuModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil {
			return err
		}
		return nil
	})

	return
}

// GetAdminOverviewPanel 總覽
func (s *SystemTable) GetAdminOverviewPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	// 面板資訊
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_OVERVIEW_GAME_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			if err := s.table(config.ACTIVITY_OVERVIEW_GAME_TABLE).
				WhereIn("id", interfaces(idArr)).Delete(); err != nil {
				return err
			}
			return nil
		})

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_OVERVIEW_GAME_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("overview_type", "overview_name",
			"name", "description", "class_id", "div_id", "url") {
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
		var model models.EditOverviewModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultOverviewModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditOverviewModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultOverviewModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil {
			return err
		}
		return nil
	})

	return
}

// models.EditOverviewModel{
// 	OverviewType: values.Get("overview_type"),
// 	OverviewName: values.Get("overview_name"),
// 	Name:         values.Get("name"),
// 	Description:  values.Get("description"),
// 	ClassID:      values.Get("class_id"),
// 	DivID:        values.Get("div_id"),
// 	URL:          values.Get("url"),
// }

// models.EditOverviewModel{
// 	ID:           values.Get("id"),
// 	OverviewType: values.Get("overview_type"),
// 	OverviewName: values.Get("overview_name"),
// 	Name:         values.Get("name"),
// 	Description:  values.Get("description"),
// 	ClassID:      values.Get("class_id"),
// 	DivID:        values.Get("div_id"),
// 	URL:          values.Get("url"),
// }

// @Summary 新增用戶(管理員用)
// @Tags Admin_Manager
// @version 1.0
// @Accept  mpfd
// @param name formData string true "用戶名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param phone formData string true "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string true "電子信箱(電子郵件地址中必須包含「@」)"
// @param password formData string true "密碼"
// @param password_again formData string true "再次輸入密碼"
// @param max_activity formData integer true "活動場次上限"
// @param max_activity_people formData integer true "活動人數上限"
// @param max_game_people formData integer true "遊戲人數上限"
// @param permissions formData string false "用戶權限(以,分隔權限)"
// @param avatar formData file false "用戶頭像"
// @param line_id formData string false "line_id"
// @param channel_id formData string false "channel_id"
// @param channel_secret formData string false "channel_secret"
// @param chatbot_secret formData string false "chatbot_secret"
// @param chatbot_token formData string false "chatbot_token"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/manager [post]
func POSTAdminManager(ctx *gin.Context) {
}

// @Summary 新增權限
// @Tags Admin_Permission
// @version 1.0
// @Accept  mpfd
// @param permission formData string true "權限名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param http_method formData string true "可使用方法(多選，用,分隔方法)" Enums(POST, PUT, DELETE)
// @param http_path formData string true "可訪問路徑(用,分隔路徑)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/permission [post]
func POSTAdminPermission(ctx *gin.Context) {
}

// @Summary 新增總覽
// @Tags Admin_Overview
// @version 1.0
// @Accept  mpfd
// @param overview_type formData string true "總覽類型" minlength(1) maxlength(50) Enums(聊天牆展示, 暖場互動, 抽獎類型, 競技類型)
// @param overview_name formData string true "總覽名稱(中文)"  minlength(1) maxlength(50)
// @param name formData string true "總覽名稱(英文)"  minlength(1) maxlength(50)
// @param description formData string true "總覽描述" minlength(1) maxlength(50)
// @param class_id formData string true "class_id" minlength(1) maxlength(50)
// @param div_id formData string true "div_id" minlength(1) maxlength(50)
// @param url formData string true "url連結" minlength(1) maxlength(100)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/overview [post]
func POSTAdminOverivew(ctx *gin.Context) {
}

// @Summary 新增菜單
// @Tags Admin_Menu
// @version 1.0
// @Accept  mpfd
// @param parent_id formData integer true "子選單"
// @param title formData string true "菜單名稱" minlength(1) maxlength(50)
// @param url formData string false "url連結" maxlength(100)
// @param sidebarbtnl formData string true "sidebarbtnl" minlength(1) maxlength(50)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/menu [post]
func POSTAdminMenu(ctx *gin.Context) {
}

// @Summary 編輯用戶
// @Tags Admin_Manager
// @version 1.0
// @Accept  mpfd
// @@@param id formData integer true "用戶ID"
// @param user formData string true "用戶ID(編輯的用戶ID)"
// @param name formData string false "用戶名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param phone formData string false "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string false "電子信箱(電子郵件地址中必須包含「@」)"
// @param password formData string false "密碼"
// @param password_again formData string false "再次輸入密碼"
// @param max_activity formData integer false "活動場次上限"
// @param max_activity_people formData integer false "活動人數上限"
// @param max_game_people formData integer false "遊戲人數上限"
// @param permissions formData string false "用戶權限(以,分隔權限)"
// @param avatar formData file false "用戶頭像"
// @param line_id formData string false "line_id"
// @param channel_id formData string false "channel_id"
// @param channel_secret formData string false "channel_secret"
// @param chatbot_secret formData string false "chatbot_secret"
// @param chatbot_token formData string false "chatbot_token"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/manager [put]
func PUTAdminManager(ctx *gin.Context) {
}

// @Summary 編輯權限
// @Tags Admin_Permission
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "用戶ID"
// @param permission formData string false "權限名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param http_method formData string false "可使用方法(多選，用,分隔方法)" (POST, PUT, DELETE)
// @param http_path formData string false "可訪問路徑(用,分隔路徑)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/permission [put]
func PUTAdminPermission(ctx *gin.Context) {
}

// @Summary 編輯總覽
// @Tags Admin_Overview
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "總覽ID"
// @param overview_type formData string false "總覽類型" minlength(1) maxlength(50) Enums(聊天牆展示, 暖場互動, 抽獎類型, 競技類型)
// @param overview_name formData string false "總覽名稱(中文)"  minlength(1) maxlength(50)
// @param name formData string false "總覽名稱(英文)"  minlength(1) maxlength(50)
// @param description formData string false "總覽描述" minlength(1) maxlength(50)
// @param class_id formData string false "class_id" minlength(1) maxlength(50)
// @param div_id formData string false "div_id" minlength(1) maxlength(50)
// @param url formData string false "url連結" minlength(1) maxlength(100)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/overview [put]
func PUTAdminOverivew(ctx *gin.Context) {
}

// @Summary 編輯菜單
// @Tags Admin_Menu
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "總覽ID"
// @param parent_id formData integer false "子選單"
// @param title formData string false "菜單名稱" minlength(1) maxlength(50)
// @param url formData string false "url連結" maxlength(100)
// @param sidebarbtnl formData string false "sidebarbtnl" minlength(1) maxlength(50)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/menu [put]
func PUTAdminMenu(ctx *gin.Context) {
}

// @Summary 刪除用戶
// @Tags Admin_Manager
// @version 1.0
// @Accept  mpfd
// @param id formData string true "用戶ID(以,區隔多筆用戶ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/manager [DELETE]
func DELETEAdminManager(ctx *gin.Context) {
}

// @Summary 刪除權限
// @Tags Admin_Permission
// @version 1.0
// @Accept  mpfd
// @param id formData string true "權限ID(以,區隔多筆權限ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/permission [DELETE]
func DELETEAdminPermission(ctx *gin.Context) {
}

// @Summary 刪除總覽
// @Tags Admin_Overview
// @version 1.0
// @Accept  mpfd
// @param id formData string true "總覽ID(以,區隔多筆總覽ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/overview [DELETE]
func DELETEAdminOverivew(ctx *gin.Context) {
}

// @Summary 刪除菜單
// @Tags Admin_Menu
// @version 1.0
// @Accept  mpfd
// @param id formData string true "菜單ID(以,區隔多筆菜單ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/menu [DELETE]
func DELETEAdminMenu(ctx *gin.Context) {
}

// @Summary 用戶JSON資料(包含用戶裡所有活動的權限)
// @Tags Admin_Manager
// @version 1.0
// @Accept  json
// @param user_id query string false "用戶ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/manager [get]
func AdminManagerJSON(ctx *gin.Context) {
}

// @Summary 權限JSON資料
// @Tags Admin_Permission
// @version 1.0
// @Accept  json
// @param id query integer false "用戶ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/permission [get]
func AdminPermissionJSON(ctx *gin.Context) {
}

// @Summary 總覽JSON資料
// @Tags Admin_Overview
// @version 1.0
// @Accept  json
// @param id query integer false "總覽ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/overview [get]
func AdminOverivewJSON(ctx *gin.Context) {
}

// @Summary 菜單JSON資料
// @Tags Admin_Menu
// @version 1.0
// @Accept  json
// @param id query integer false "菜單ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/menu [get]
func AdminMenuJSON(ctx *gin.Context) {
}

// @Summary 操作紀錄JSON資料
// @Tags Admin_Menu
// @version 1.0
// @Accept  json
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /admin/log [get]
func AdminLogJSON(ctx *gin.Context) {
}

// 取得所有用戶ID
// items, err := s.table(config.USERS_TABLE).
// 	WhereIn("id", interfaces(idArr)).All()
// if err != nil {
// 	return errors.New("錯誤: 查詢所有用戶資訊發生問題")
// }
// for _, id := range idArr {
// 	users = append(users, id)
// }

// var avatar string
// if values.Get("avatar") != "" {
// 	avatar = values.Get("avatar")
// } else {
// 	avatar = UPLOAD_SYSTEM_URL + "img-user-pic.png"

// models.EditUserModel{
// 	UserID:            values.Get("user"),
// 	Name:              values.Get("name"),
// 	Phone:             values.Get("phone"),
// 	Email:             values.Get("email"),
// 	Password:          values.Get("password"),
// 	PasswordAgain:     values.Get("password_again"),
// 	Permissions:       values.Get("permissions"),
// 	Avatar:            avatar,
// 	Bind:              "",
// 	Cookie:            "",
// 	MaxActivity:       values.Get("max_activity"),
// 	MaxActivityPeople: values.Get("max_activity_people"),
// 	MaxGamePeople:     values.Get("max_game_people"),

// 	// 官方帳號綁定
// 	LineID:        values.Get("line_id"),
// 	ChannelID:     values.Get("channel_id"),
// 	ChannelSecret: values.Get("channel_secret"),
// 	ChatbotSecret: values.Get("chatbot_secret"),
// 	ChatbotToken:  values.Get("chatbot_token"),
// }

// models.EditUserModel{
// 	Name:              values.Get("name"),
// 	Phone:             values.Get("phone"),
// 	Email:             values.Get("email"),
// 	Password:          values.Get("password"),
// 	PasswordAgain:     values.Get("password_again"),
// 	Permissions:       values.Get("permissions"),
// 	Avatar:            UPLOAD_SYSTEM_URL + "img-user-pic.png",
// 	Bind:              "no",
// 	Cookie:            "no",
// 	MaxActivity:       values.Get("max_activity"),
// 	MaxActivityPeople: values.Get("max_activity_people"),
// 	MaxGamePeople:     values.Get("max_game_people"),

// 	// 官方帳號綁定
// 	LineID:        values.Get("line_id"),
// 	ChannelID:     values.Get("channel_id"),
// 	ChannelSecret: values.Get("channel_secret"),
// 	ChatbotSecret: values.Get("chatbot_secret"),
// 	ChatbotToken:  values.Get("chatbot_token"),
// }

// models.EditPermissionModel{
// 	Permission: values.Get("permission"),
// 	HTTPMethod: values.Get("http_method"),
// 	HTTPPath:   values.Get("http_path"),
// }

// models.EditPermissionModel{
// 	ID:         values.Get("id"),
// 	Permission: values.Get("permission"),
// 	HTTPMethod: values.Get("http_method"),
// 	HTTPPath:   values.Get("http_path"),
// }

// models.EditMenuModel{
// 	ParentID:    values.Get("parent_id"),
// 	Title:       values.Get("title"),
// 	URL:         values.Get("url"),
// 	SidebarBtnL: values.Get("sidebarbtnl"),
// }

// models.EditMenuModel{
// 	ID:          values.Get("id"),
// 	ParentID:    values.Get("parent_id"),
// 	Title:       values.Get("title"),
// 	URL:         values.Get("url"),
// 	SidebarBtnL: values.Get("sidebarbtnl"),
// }
