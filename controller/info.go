package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShowGuestInfo 活動資訊頁面(手機用戶端)，包含主頁面、介紹、行程、嘉賓、資料 GET API
func (h *Handler) ShowGuestInfo(ctx *gin.Context) {
	var (
		host                = ctx.Request.Host
		prefix              = ctx.Param("__prefix")
		activityID          = ctx.Query("activity_id")
		activityModel, err  = h.getActivityInfo(false, activityID)
		panelInfo           table.PanelInfo
		executeparam        executeParam
		sortField, sortType []string
		htmlTmpl            string
		// err                 error
	)

	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if err != nil || activityID == "" || activityModel.ID == 0 {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}
	// if activityID == "" {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
	// 	return
	// }

	// template、order
	if prefix == "/" {
		prefix = ""
		// fmt.Println("活動資訊主頁面")
	} else if prefix == "/introduce" {
		prefix = "introduce"
		sortField = []string{"introduce_order"}
		sortType = []string{"asc"}
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/introduce.html"
	} else if prefix == "/schedule" {
		prefix = "schedule"
		sortField = []string{"schedule_date", "start_time"}
		sortType = []string{"asc", "asc"}
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/schedule.html"
	} else if prefix == "/guest" {
		prefix = "guest"
		sortField = []string{"guest_order"}
		sortType = []string{"asc"}
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/guest.html"
	} else if prefix == "/material" {
		prefix = "material"
		sortField = []string{"material_order"}
		sortType = []string{"asc"}
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/material.html"
	}

	// 取得所有資料
	if prefix != "" {
		panel, _ := h.GetTable(ctx, prefix)
		param := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			sortField, sortType)
		// if panelInfo, err = panel.GetData(param, h.services, []string{}); err != nil {
		if panelInfo, err = panel.GetData(param, h.services); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
			return
		}
	}

	// fmt.Println("len(panelInfo.InfoList): ", len(panelInfo.InfoList))
	// fmt.Println("panelInfo.InfoList: ", panelInfo.InfoList)
	executeparam = executeParam{
		PanelInfo:     panelInfo,
		ActivityModel: activityModel,
		// Route: route{
		// Chatroom: fmt.Sprintf(config.GUEST_CHATROOM_URL, activityID),
		// Info:      fmt.Sprintf("%s/?activity_id=%s", config.GUEST_INFO_URL, activityID),
		// GameInfo:  fmt.Sprintf("%s?activity_id=%s", config.GUEST_GAME_URL, activityID),
		// Introduce: fmt.Sprintf(config.GUEST_INTRODUCE_URL, activityID),
		// Schedule:  fmt.Sprintf(config.GUEST_SCHEDULE_URL, activityID),
		// Guest:     fmt.Sprintf(config.GUEST_GUEST_URL, activityID),
		// Material:  fmt.Sprintf(config.GUEST_MATERIAL_URL, activityID),
		// },
	}
	h.executeHTML(ctx, htmlTmpl, executeparam)
}

// ShowAdminInfo 活動資訊頁面(平台) GET API
func (h *Handler) ShowAdminInfo(ctx *gin.Context) {
	var (
		host                       = ctx.Request.Host
		path                       = ctx.Request.URL.Path
		prefix                     = ctx.Param("__prefix")
		activityID                 = ctx.Query("activity_id")
		sidebar                    = ctx.Query("sidebar")
		panel, _                   = h.GetTable(ctx, prefix)
		user                       = h.GetLoginUser(ctx.Request, "hilive_session")
		panelInfo                  table.PanelInfo
		newFormInfo                table.FormInfo
		settingFormInfo            table.FormInfo
		param                      parameter.Parameters
		hasSidebar                 bool
		deleteURL, htmlTmpl        string
		sortField, sortType        []string
		canAdd, canEdit, canDelete bool
		err                        error
	)

	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if prefix == "overview" {
		hasSidebar, err = strconv.ParseBool(sidebar)
		if err != nil {
			hasSidebar = false
		}
	}

	// template
	if prefix == "introduce" {
		htmlTmpl = "./hilives/hilive/views/cms/info_introduce.html"
		sortField = []string{"introduce_order"}
		sortType = []string{"asc"}
	} else if prefix == "schedule" {
		sortField = []string{"schedule_date", "start_time"}
		sortType = []string{"asc", "asc"}
		htmlTmpl = "./hilives/hilive/views/cms/info_schedule.html"
	} else if prefix == "overview" {
		// tempName = "layout"
		// fields = []string{"overview_id"}
		sortField = []string{"id"}
		sortType = []string{"asc"}
		// htmlTmpl = "./hilives/hilive/views/cms/info_overview.html"

		if hasSidebar == false {
			htmlTmpl = "./hilives/hilive/views/cms/info_overview.html"
		} else if hasSidebar == true {
			// 活動頁面
			htmlTmpl = "./hilives/hilive/views/cms/activity_settings.html"
		}
	} else if prefix == "guest" {
		htmlTmpl = "./hilives/hilive/views/cms/info_guest.html"
		sortField = []string{"guest_order"}
		sortType = []string{"asc"}
	} else if prefix == "material" {
		sortField = []string{"material_order"}
		sortType = []string{"asc"}
		htmlTmpl = "./hilives/hilive/views/cms/info_material.html"
	}

	param = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		sortField, sortType).SetPKs(user.UserID, activityID)

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	if prefix != "overview" {
		// 所有資料
		// if panelInfo, err = panel.GetData(param, h.services, fields); err != nil {
		if panelInfo, err = panel.GetData(param, h.services); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
			return
		}

		// 新增資料的表單
		newFormInfo = panel.GetNewFormInfo(h.services, param, []string{"user_id", "activity_id"})
		deleteURL = fmt.Sprintf(config.INFO_API_URL, prefix)

		// 基本設置表單
		if settingFormInfo, err = h.tableList[prefix]().GetSettingFormInfo(
			param, h.services, []string{"user_id", "activity_id"}); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
			return
		}
	} else if prefix == "overview" {
		deleteURL = config.ACTIVITY_API_URL
	}

	// fmt.Println("活動權限長度: ", len(user.ActivityMenus[activityID]))
	// fmt.Println("活動權限: ", user.ActivityMenus[activityID])
	h.executeHTML(ctx, htmlTmpl, executeParam{
		HasSidebar: hasSidebar,
		ActivityID: activityID,
		User:       user,
		UserJSON:   utils.JSON(user),
		PanelInfo:  panelInfo,
		Route: route{
			Prefix: config.Prefix(),
			POST:   fmt.Sprintf(config.INFO_API_URL, prefix),
			PUT:    fmt.Sprintf(config.INFO_API_URL, prefix),
			PATCH:  fmt.Sprintf(config.INFO_API_URL, prefix),
			DELETE: deleteURL,
			Back:   fmt.Sprintf(config.INFO_URL, prefix, activityID),
			Edit:   fmt.Sprintf(config.ACTIVITY_EDIT_URL, activityID), // 編輯活動頁面
		},
		FormInfo:        newFormInfo,
		SettingFormInfo: settingFormInfo,
		Token:           auth.AddToken(user.UserID),
		CanAdd:          canAdd,
		CanEdit:         canEdit,
		CanDelete:       canDelete,
	})
}

// activityModel models.ActivityModel
// activityInfo map[string]interface{}
// editFormInfos       []table.FormInfo
// executeparam    executeParam
// fields              = []string{"id"}
// tempName                   = "content"

// 所有資料編輯的表單
// for _, info := range panelInfo.InfoList {
// 	var (
// 		panel2   = h.tableList[prefix]()
// 		formInfo table.FormInfo
// 	)
// 	editParam := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize, sortField, sortType).
// 		SetPKs(fmt.Sprintf("%s", info["id"].Content), activityID)
// 	if formInfo, err = panel2.GetEditFormInfo(editParam, h.services, []string{"id", "activity_id"}); err != nil {
// 		response.Error(ctx, err.Error())
// 		return
// 	}
// 	editFormInfos = append(editFormInfos, formInfo)
// }

// 活動總覽
// if activityModel, err = models.DefaultActivityModel().SetDbConn(h.dbConn).
// 	Find(false, activityID); activityID == "" ||
// 	err != nil || activityModel.ID == 0 {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
// 	return
// }

// 活動資訊json資料
// activityJSON, err := json.Marshal(activityModel)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
// 	return
// }

// 測試menu參數
// menus := menu.GetMenu(user, h.dbConn)
// for _, menu := range menus.List {
// 	fmt.Println("menu.GetMenu(user, h.dbConn): ", menu)
// 	fmt.Println("--------------------------------------")
// }

// menus, err := models.DefaultMenuModel().SetDbConn(h.dbConn).Find(0)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得側邊菜單資訊，請重新整理頁面")
// 	return
// }
// executeparam = executeParam{

// HasSidebar: hasSidebar,
// PanelInfo:  panelInfo,
// Route: route{
// Prefix: config.Prefix(),
// PATCH:  fmt.Sprintf(config.INFO_API_URL, prefix),
// Edit:   fmt.Sprintf(config.ACTIVITY_EDIT_URL, activityID), // 編輯活動頁面
// DELETE: config.ACTIVITY_API_URL,                           // 刪除活動
// Message:      fmt.Sprintf("%s/message?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// Topic:        fmt.Sprintf("%s/topic?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// Question:     fmt.Sprintf("%s/question?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// Danmu:        fmt.Sprintf("%s/danmu?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// SpecialDanmu: fmt.Sprintf("%s/specialdanmu?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// Picture:      fmt.Sprintf("%s/picture?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// Holdscreen:   fmt.Sprintf("%s/holdscreen?activity_id=%s", config.INTERACT_WALL_URL, activityID),
// General:      fmt.Sprintf("%s/general?activity_id=%s", config.INTERACT_SIGN_URL, activityID),
// Threed:       fmt.Sprintf("%s/threed?activity_id=%s", config.INTERACT_SIGN_URL, activityID),
// Countdown:    fmt.Sprintf("%s/countdown?activity_id=%s", config.INTERACT_SIGN_URL, activityID),
// Lottery:      fmt.Sprintf("%s/lottery?activity_id=%s", config.INTERACT_GAME_URL, activityID),
// Redpack:      fmt.Sprintf("%s/redpack?activity_id=%s", config.INTERACT_GAME_URL, activityID),
// Ropepack:     fmt.Sprintf("%s/ropepack?activity_id=%s", config.INTERACT_GAME_URL, activityID),
// WhackMole:    fmt.Sprintf("%s/whack_mole?activity_id=%s", config.INTERACT_GAME_URL, activityID),
// },
// Menu:       menus,
// ActivityID: activityID,
// User:       user,
// ActivityModel: activityModel,
// ActivityJSON: utils.JSON(activityModel),
// Token:     auth.AddToken( user.UserID),
// CanAdd:    canAdd,
// CanEdit:   canEdit,
// CanDelete: canDelete,
// }
// }

// MessageWebsocket 更新訊息牆資料，即時推送至主持人端
// func (h *Handler) MessageWebsocket(ctx *gin.Context) {
// var (
// 	host       = ctx.Request.Host
// 	activityID = ctx.Query("activity_id")
// 	result     WebsocketMessage
// )
// upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
// 	return true
// }}
// wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// defer wsConn.Close()
// if err != nil {
// 	return
// }
// conn, err := NewWebsocket(wsConn)
// defer conn.Close()
// if err != nil {
// 	conn.Close()
// 	return
// }

// if host != API_HOST {
// 	b, _ := json.Marshal(WebsocketMessage{Error: "錯誤: 錯誤的網域請求"})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
// }
// if activityID == "" {
// 	b, _ := json.Marshal(WebsocketMessage{Error: "錯誤: 無法辨識活動資訊"})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
// }

// for {
// 	data, err := conn.ReadMessage()
// 	if err != nil {
// 		conn.Close()
// 		return
// 	}

// 	log.Println("收到前端訊息")
// 	json.Unmarshal(data, &result)
// 	fmt.Println("data參數: ", string(data), "result解碼參數: ", result)
// 	if result.MessageModel.ActivityID == "" || result.MessageModel.ActivityID != activityID {
// 		b, _ := json.Marshal(WebsocketMessage{Error: "錯誤: 傳遞訊息牆資訊發生錯誤"})
// 		conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 		return
// 	}

// 	b, _ := json.Marshal(result)
// 	if err = conn.WriteMessage(b); err != nil {
// 		conn.Close()
// 		return
// 	}
