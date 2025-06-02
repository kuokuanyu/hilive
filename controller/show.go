package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Show 用戶、活動頁面(平台) GET API
func (h *Handler) Show(ctx *gin.Context) {
	var (
		host                                = ctx.Request.Host
		user                                = h.GetLoginUser(ctx.Request, "hilive_session")
		panelInfo                           table.PanelInfo
		header                              = ctx.Query("header")
		isAdmin, canAdd, canEdit, canDelete bool
		prefix, htmlTmpl                    string
		// err                                 error
		// ip        = utils.ClientIP(ctx.Request)
		// path      = ctx.Request.URL.Path
		// formInfo  table.FormInfo
		// fields   = []string{}
		// tmplName = "layout_info"
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if strings.Contains(ctx.Request.URL.Path, "activity") {
		// 活動頁面
		prefix = "activity"
		htmlTmpl = "./hilives/hilive/views/cms/index.html"
		// fields = []string{"activity_id"}
	} else if strings.Contains(ctx.Request.URL.Path, "user") {
		//用戶頁面
		prefix = "user"
		htmlTmpl = "./hilives/hilive/views/cms/manager.html"
	}

	// table
	panel, _ := h.GetTable(ctx, prefix)
	param := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		[]string{"id"}, []string{"desc"})
	if prefix == "user" {
		param.SetPKs(user.UserID)
	}
	// 只取得特定user資料(url隱藏user_id資料)
	param.SetField("user_id", user.UserID)

	hasHeader, err := strconv.ParseBool(header)
	if err != nil {
		// h.executeErrorHTML(ctx, "錯誤: header參數發生問題，請重新整理頁面")
		return
	}
	if hasHeader == false {
		if prefix == "activity" {
			htmlTmpl = "./hilives/hilive/views/cms/activity_info.html"
			// tmplName = "activity_content"
		} else if prefix == "user" {
			htmlTmpl = "./hilives/hilive/views/cms/manager.html"
			// tmplName = "profile_content"
		}
	} else if hasHeader == true {
		// 活動頁面
		htmlTmpl = "./hilives/hilive/views/cms/index.html"
	}

	// 面板、表單資訊
	// if panelInfo, err = panel.GetData(param, h.services, fields); err != nil {
	if prefix == "activity" {
		if panelInfo, err = panel.GetData(param, h.services); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
			return
		}
	} else if prefix == "user" {
		// var params = url.Values{}
		// if formInfo, err = panel.GetEditFormInfo(param, h.services,
		// 	[]string{"user_id"}); err != nil {
		// 	h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
		// 	return
		// }
	}

	// 是否有管理後台權限
	for _, permission := range user.Permissions {
		if len(permission.HTTPPath) > 0 &&
			permission.HTTPPath[0] == "admin" {
			isAdmin = true
			break
		}
	}

	// 判斷是否有新增編輯刪除權限
	// 判斷是否為試用用戶(活動場次上限為1)，試用用戶沒有新增編輯刪除功能
	if user.MaxActivity > 1 {
		canEdit = true
		// canDelete = true
		if int(user.MaxActivity) <= len(panelInfo.InfoList) {
			canAdd = false
		} else {
			canAdd = true
		}
	}

	executeparam := executeParam{
		User: user,
		// HasHeader: hasHeader,
		// Role:      user.Roles[0].Name,
		PanelInfo: panelInfo,
		// FormInfo:  formInfo,
		Route: route{
			Logout:     config.LOGOUT_URL,
			Admin:      fmt.Sprintf(config.ADMIN_URL, "manager"),                                                  // 管理員頁面
			User:       fmt.Sprintf(config.USER_URL, "false"),                                                     // 用戶頁面
			Activity:   fmt.Sprintf(config.ACTIVITY_URL, "false"),                                                 // 活動頁面
			Overview:   config.OVERVIEW_URL,                                                                       // 活動總覽頁面
			Host:       config.HOST_CHATROOM_URL,                                                                  // 主持端聊天室頁面
			Bind:       fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, "bind", "line", "user_id", user.UserID, ""), // 綁定用戶頁面
			New:        config.ACTIVITY_NEW_URL,                                                                   // 新增活動頁面
			Edit:       config.ACTIVITY_EDIT_NO_ACTIVITYID_URL,                                                    // 編輯活動頁面
			PUT:        config.USER_API_URL,                                                                       // 編輯用戶資訊
			DELETE:     config.ACTIVITY_API_URL,
			QuickStart: config.ACTIVITY_QUICK_START_URL, // 刪除活動
		},
		Token:     auth.AddToken(user.UserID),
		IsAdmin:   isAdmin,
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	}

	// fmt.Println("canAdd: ", canAdd, "canEdit: ", canEdit, "canDelete: ", canDelete)
	// if (prefix == "activity" && hasHeader == false) || prefix == "user" {
	h.executeHTML(ctx, htmlTmpl, executeparam)
	// } else {
	// h.execute(ctx, tmplName, htmlTmpl, executeparam)
	// }
	return
}

// if prefix == "user" {
// 	executeparam.Token = auth.GetTokenService(h.services.Get("token_csrf")).AddToken(h.conn, user)
// }

// userModel, err := models.DefaultUserModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).Find(false, "", "user_id", user.UserID)
// if err != nil || userModel.UserID == "" {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
// 	return
// }

// // 用戶資訊
// params.Add("user_id", userModel.UserID)
// params.Add("name", userModel.Name)
// params.Add("avatar", userModel.Avatar)
// params.Add("phone", userModel.Phone)
// params.Add("email", userModel.Email)
// params.Add("bind", userModel.Bind)
// params.Add("cookie", userModel.Cookie)
// params.Add("ip", ip)
// params.Add("table", "users")
// params.Add("sign", utils.UserSign(ip, user.UserID, config.HiliveCookieSecret))
// // 設置cookie
// ctx.SetSameSite(4)
// ctx.SetCookie("hilive_session", string(utils.Encode([]byte(params.Encode()))),
// 	config.GetSessionLifeTime(), "/", "", true, false)
