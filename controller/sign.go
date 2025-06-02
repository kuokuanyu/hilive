package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// ShowInteractSignOption 投票選項資訊頁面(平台)，顯示選項資訊 GET API
func (h *Handler) ShowInteractSignOption(ctx *gin.Context) {
	var (
		host                       = ctx.Request.Host
		path                       = ctx.Request.URL.Path
		prefix                     = ctx.Param("__prefix")
		activityID                 = ctx.Query("activity_id")
		gameID                     = ctx.Query("game_id")
		user                       = h.GetLoginUser(ctx.Request, "hilive_session")
		canAdd, canEdit, canDelete bool
		htmlTmpl                   string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" || gameID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的ID")
		return
	}

	if prefix == "vote" {
		// 遊戲連結
		htmlTmpl = "./hilives/hilive/views/cms/vote_option.html"
	}

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		UserJSON:   utils.JSON(user),
		ActivityID: activityID,
		Route: route{
			Back: fmt.Sprintf(config.INTERACT_SIGN_URL, prefix, activityID),
			Host: fmt.Sprintf(config.HOST_GAME_URL, activityID, "false"), // 主持端聊天室頁面
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// ShowSign 簽到展示頁面(平台) GET API
func (h *Handler) ShowSign(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		formInfo   table.FormInfo
		// model, err = models.DefaultActivityModel().SetConn(h.conn).Find("activity_id", activityID)
		canAdd, canEdit, canDelete bool
		htmlTmpl                   string
		err                        error
		// params                     parameter.Parameters
		panelInfo table.PanelInfo
		panel     table.Table
		param     parameter.Parameters
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if prefix == "general" {
		// tempName = "content"
		htmlTmpl = "./hilives/hilive/views/cms/avatar_collage_2D.html"
	} else if prefix == "threed" {
		htmlTmpl = "./hilives/hilive/views/cms/avatar_collage_3D.html"
	} else if prefix == "countdown" {
		// tempName = "content"
		// htmlTmpl = countdown.Content + countdown.Form1 + countdown.Form2
	} else if prefix == "signname" {
		htmlTmpl = "./hilives/hilive/views/cms/signname.html"
	} else if prefix == "vote" {
		htmlTmpl = "./hilives/hilive/views/cms/vote_info.html"
	} else if prefix == "signname_check" {

		htmlTmpl = "./hilives/hilive/views/cms/signname_check.html"
	}

	if prefix != "signname_check" {
		// 簽名牆審核頁面不用執行
		panel, _ = h.GetTable(ctx, prefix)
		param = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID)
	}

	// 投票頁面之外才需要執行
	// 簽名牆審核頁面不用執行
	if prefix != "vote" && prefix != "signname_check" {
		if formInfo, err = panel.GetSettingFormInfo(param, h.services,
			[]string{"user_id", "activity_id"}); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
			return
		}
	}

	// 投票頁面
	if prefix == "vote" {
		// params = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		// 	[]string{"id"}, []string{"asc"})
		// params.SetField("game", prefix) // 遊戲種類判斷參數

		// if panelInfo, err = panel.GetData(params, h.services); err != nil {
		// 	h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
		// 	return
		// }
	}

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	// if prefix != "threed" {
	// h.execute(ctx, tempName, htmlTmpl, executeParam{
	// 	ActivityID:      activityID,
	// 	User:            user,
	// 	SettingFormInfo: formInfo,
	// 	Route: route{
	// 		PATCH: fmt.Sprintf(config.INTERACT_SIGN_API_URL, prefix),
	// 	},
	// 	Token:     auth.AddToken(user.UserID),
	// 	CanAdd:    canAdd,
	// 	CanEdit:   canEdit,
	// 	CanDelete: canDelete,
	// })
	// } else {
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:            user,
		UserJSON:        utils.JSON(user),
		PanelInfo:       panelInfo,
		ActivityID:      activityID,
		SettingFormInfo: formInfo,
		Token:           auth.AddToken(user.UserID),
		CanAdd:          canAdd,
		CanEdit:         canEdit,
		CanDelete:       canDelete,
		Route: route{
			Host:  fmt.Sprintf(config.HOST_GAME_URL, activityID, "false"),          // 主持端聊天室頁面
			Prize: fmt.Sprintf(config.INTERACT_SIGN_PRIZE_URL, prefix, activityID), // 獎品頁面
			// POST:   fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix), // 新增場次
			PUT:        fmt.Sprintf(config.INTERACT_SIGN_API_URL_FORM, prefix),                // 編輯場次
			DELETE:     fmt.Sprintf(config.INTERACT_SIGN_API_URL_FORM, prefix),                // 刪除場次
			New:        fmt.Sprintf(config.INTERACT_SIGN_NEW_URL, prefix, activityID),         // 新增遊戲頁面
			Edit:       fmt.Sprintf(config.INTERACT_SIGN_EDIT_URL, prefix, activityID),        // 編輯遊戲頁面
			PATCH:      fmt.Sprintf(config.INTERACT_SIGN_API_URL, prefix),                     // 編輯設置
			VoteOption: fmt.Sprintf(config.INTERACT_SIGN_VOTE_OPTION_URL, prefix, activityID), // 投票選項設置頁面
			Signname:   fmt.Sprintf(config.INTERACT_SIGN_URL, prefix, activityID),             // 簽名牆訊息審核頁面
		},
	})
	// }
}

// ShowNewInteractSign 新增場次頁面(平台)，新建場次 GET API
func (h *Handler) ShowNewInteractSign(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		panel, _   = h.GetTable(ctx, prefix)
		params     = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"asc"}).SetPKs(user.UserID, activityID)
		htmlTmpl string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if prefix == "vote" {
		htmlTmpl = "./hilives/hilive/views/cms/vote_form.html"
	}

	// fmt.Println("人數: ", user.MaxGamePeople)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		People:     int(user.MaxGamePeople),
		FormInfo:   panel.GetNewFormInfo(h.services, params, []string{"user_id", "activity_id"}),
		ActivityID: activityID,
		Route: route{
			Method: "post",
			POST:   fmt.Sprintf(config.INTERACT_SIGN_API_URL_FORM, prefix),
			Back:   fmt.Sprintf(config.INTERACT_SIGN_URL, prefix, activityID),
			Import: config.IMPORT_EXCEL_API_URL,
		},
		Token: auth.AddToken(user.UserID),
	})
}

// ShowEditInteractSign 編輯場次頁面(平台)，編輯場次 GET API
func (h *Handler) ShowEditInteractSign(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panel, _   = h.GetTable(ctx, prefix)
		// params     = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		// 	[]string{"id"}, []string{"asc"}).SetPKs(gameID)
		htmlTmpl string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" || gameID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動、遊戲ID")
		return
	}

	if prefix == "vote" {
		htmlTmpl = "./hilives/hilive/views/cms/vote_form.html"
	}

	// formInfo, err := panel.GetEditFormInfo(params, h.services, []string{"game_id"})
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
	// 	return
	// }

	// 上限人數
	gameModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(false, gameID, prefix)
	if err != nil || gameModel.ID == 0 {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識遊戲資訊，請輸入有效的遊戲ID")
		return
	}

	// fmt.Println("人數: ", user.MaxGamePeople)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:   user,
		People: int(gameModel.MaxPeople),
		// FormInfo:   formInfo,
		ActivityID: activityID,
		GameID:     gameID,
		Route: route{
			Method: "put",
			PUT:    fmt.Sprintf(config.INTERACT_SIGN_API_URL_FORM, prefix),
			Back:   fmt.Sprintf(config.INTERACT_SIGN_URL, prefix, activityID),
			Import: config.IMPORT_EXCEL_API_URL,
		},
		Token: auth.AddToken(user.UserID),
	})
}
