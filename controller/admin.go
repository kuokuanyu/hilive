package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowEditAdmin 管理員編輯頁面(平台) GET API
func (h *Handler) ShowEditAdmin(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		id         = ctx.Query("id")
		userID     = ctx.Query("user_id")
		activityID = ctx.Query("activity_id")
		game       = ctx.Query("game")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panelInfo        table.PanelInfo
		// activityModel                     models.ActivityModel
		apiURL, backURL, prefix, htmlTmpl string
		// err                               error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if id == "" && userID == "" && activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識資訊，請輸入有效的ID")
		return
	}

	if strings.Contains(path, "/manager/activity/game") {
		prefix = "game"
		if game == "redpack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_redpack_form.html"
		} else if game == "ropepack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_ropepack_form.html"
		} else if game == "draw_numbers" {
			htmlTmpl = "./hilives/hilive/views/admin/game_drawnumbers_form.html"
		} else if game == "lottery" {
			htmlTmpl = "./hilives/hilive/views/admin/game_lottery_form.html"
		} else if game == "QA" {
			htmlTmpl = "./hilives/hilive/views/admin/game_qa_form.html"
		} else if game == "monopoly" {
			htmlTmpl = "./hilives/hilive/views/admin/game_monopoly_form.html"
		} else if game == "whack_mole" {
			htmlTmpl = "./hilives/hilive/views/admin/game_whackmole_form.html"
		} else if game == "tugofwar" {
			htmlTmpl = "./hilives/hilive/views/admin/game_tugofwar_form.html"
		} else if game == "bingo" {
			htmlTmpl = "./hilives/hilive/views/admin/game_bingo_form.html"
		}
	} else if strings.Contains(path, "/manager/activity") {
		prefix = "activity"
		htmlTmpl = "./hilives/hilive/views/admin/admin_activity_form.html"
	} else if strings.Contains(path, "manager") {
		prefix = "manager"
		htmlTmpl = "./hilives/hilive/views/admin/manager_form.html"
	} else if strings.Contains(path, "permission") {
		prefix = "permission"
		htmlTmpl = "./hilives/hilive/views/admin/permission_form.html"
	} else if strings.Contains(path, "menu") {
		prefix = "menu"
		htmlTmpl = "./hilives/hilive/views/admin/menu_form.html"
	} else if strings.Contains(path, "overview") {
		prefix = "overview"
		htmlTmpl = "./hilives/hilive/views/admin/overview_form.html"
	}

	if prefix == "game" {
		if game != "" {
			// 各別遊戲頁面
			apiURL = fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, game)
			backURL = fmt.Sprintf(config.ADMIN_GAME_URL, userID, activityID, game)
		}
	} else if prefix == "activity" {
		// 活動頁面
		// activityModel, err = models.DefaultActivityModel().SetDbConn(h.dbConn).
		// 	Find(false, activityID)
		// if err != nil || activityModel.ID == 0 {
		// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		// 	return
		// }

		apiURL = config.ACTIVITY_API_URL
		backURL = fmt.Sprintf(config.ADMIN_ACTIVITY_URL, userID)
	} else {
		// 用戶、權限、菜單、總覽頁面
		apiURL = fmt.Sprintf(config.ADMIN_API_URL, prefix)
		backURL = fmt.Sprintf(config.ADMIN_URL, prefix)
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		// ActivityModel: activityModel,
		User: user,
		// FormInfo:   formInfo,
		Route: route{
			Method: "put",
			PUT:    apiURL,
			Back:   backURL,
		},
		Token: auth.AddToken(user.UserID),
	})
	return
}

// ShowNewAdmin 管理員新增頁面(平台) GET API
func (h *Handler) ShowNewAdmin(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		userID     = ctx.Query("user_id")
		activityID = ctx.Query("activity_id")
		game       = ctx.Query("game")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panelInfo        table.PanelInfo
		apiURL, backURL, prefix, htmlTmpl string
		// err              error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if strings.Contains(path, "/manager/activity/game") {
		prefix = "game"
		if game == "redpack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_redpack_form.html"
		} else if game == "ropepack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_ropepack_form.html"
		} else if game == "draw_numbers" {
			htmlTmpl = "./hilives/hilive/views/admin/game_drawnumbers_form.html"
		} else if game == "lottery" {
			htmlTmpl = "./hilives/hilive/views/admin/game_lottery_form.html"
		} else if game == "QA" {
			htmlTmpl = "./hilives/hilive/views/admin/game_qa_form.html"
		} else if game == "monopoly" {
			htmlTmpl = "./hilives/hilive/views/admin/game_monopoly_form.html"
		} else if game == "whack_mole" {
			htmlTmpl = "./hilives/hilive/views/admin/game_whackmole_form.html"
		} else if game == "tugofwar" {
			htmlTmpl = "./hilives/hilive/views/admin/game_tugofwar_form.html"
		} else if game == "bingo" {
			htmlTmpl = "./hilives/hilive/views/admin/game_bingo_form.html"
		}
	} else if strings.Contains(path, "/manager/activity") {
		prefix = "activity"
		htmlTmpl = "./hilives/hilive/views/admin/admin_activity_form.html"
	} else if strings.Contains(path, "manager") {
		prefix = "manager"
		htmlTmpl = "./hilives/hilive/views/admin/manager_form.html"
	} else if strings.Contains(path, "permission") {
		prefix = "permission"
		htmlTmpl = "./hilives/hilive/views/admin/permission_form.html"
	} else if strings.Contains(path, "menu") {
		prefix = "menu"
		htmlTmpl = "./hilives/hilive/views/admin/menu_form.html"
	} else if strings.Contains(path, "overview") {
		prefix = "overview"
		htmlTmpl = "./hilives/hilive/views/admin/overview_form.html"
	}

	if prefix == "game" {
		if game != "" {
			// 各別遊戲頁面
			apiURL = fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, game)
			backURL = fmt.Sprintf(config.ADMIN_GAME_URL, userID, activityID, game)
		}
	} else if prefix == "activity" {
		// 活動頁面
		if userID == "" {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識資訊，請輸入有效的ID")
			return
		}

		apiURL = config.ACTIVITY_API_URL
		backURL = fmt.Sprintf(config.ADMIN_ACTIVITY_URL, userID)
	} else {
		// 用戶、權限、菜單、總覽頁面
		apiURL = fmt.Sprintf(config.ADMIN_API_URL, prefix)
		backURL = fmt.Sprintf(config.ADMIN_URL, prefix)
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User: user,
		// FormInfo:   panel.GetNewFormInfo(h.services, params, []string{"user_id"}),
		Route: route{
			Method: "post",
			POST:   apiURL,
			Back:   backURL,
		},
		Token: auth.AddToken(user.UserID),
	})
	return
}

// ShowAdmin 管理員頁面(平台) GET API
func (h *Handler) ShowAdmin(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		activityID = ctx.Query("activity_id")
		userID     = ctx.Query("user_id")
		game       = ctx.Query("game")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panelInfo        table.PanelInfo
		prefix, htmlTmpl, newURL, editURL, deleteURL, activityURL, gameURL string
		// err              error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if strings.Contains(path, "/manager/activity/game") {
		prefix = "game"
		if game == "" {
			htmlTmpl = "./hilives/hilive/views/admin/game_type.html"
		} else if game == "redpack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_redpack.html"
		} else if game == "ropepack" {
			htmlTmpl = "./hilives/hilive/views/admin/game_ropepack.html"
		} else if game == "draw_numbers" {
			htmlTmpl = "./hilives/hilive/views/admin/game_drawnumbers.html"
		} else if game == "lottery" {
			htmlTmpl = "./hilives/hilive/views/admin/game_lottery.html"
		} else if game == "QA" {
			htmlTmpl = "./hilives/hilive/views/admin/game_qa.html"
		} else if game == "monopoly" {
			htmlTmpl = "./hilives/hilive/views/admin/game_monopoly.html"
		} else if game == "whack_mole" {
			htmlTmpl = "./hilives/hilive/views/admin/game_whackmole.html"
		} else if game == "tugofwar" {
			htmlTmpl = "./hilives/hilive/views/admin/game_tugofwar.html"
		} else if game == "bingo" {
			htmlTmpl = "./hilives/hilive/views/admin/game_bingo.html"
		}
	} else if strings.Contains(path, "/manager/activity") {
		prefix = "activity"
		htmlTmpl = "./hilives/hilive/views/admin/admin_activity.html"
	} else if strings.Contains(path, "manager") {
		prefix = "manager"
		htmlTmpl = "./hilives/hilive/views/admin/manager.html"
	} else if strings.Contains(path, "permission") {
		prefix = "permission"
		htmlTmpl = "./hilives/hilive/views/admin/permission.html"
	} else if strings.Contains(path, "menu") {
		prefix = "menu"
		htmlTmpl = "./hilives/hilive/views/admin/menu.html"
	} else if strings.Contains(path, "overview") {
		prefix = "overview"
		htmlTmpl = "./hilives/hilive/views/admin/overview.html"
	} else if strings.Contains(path, "error_log") {
		// 操作日至
		prefix = "error_log"
		htmlTmpl = "./hilives/hilive/views/admin/log_error.html"
	}else if strings.Contains(path, "log") {
		// 操作日至
		prefix = "log"
		htmlTmpl = "./hilives/hilive/views/admin/log.html"
	}

	if prefix == "game" {
		// 遊戲頁面
		if game == "" {
			// 種類頁面
			gameURL = fmt.Sprintf(config.ADMIN_GAME_URL, userID, activityID, "")
		} else if game != "" {
			// 各別遊戲頁面
			newURL = fmt.Sprintf(config.ADMIN_GAME_NEW_URL, userID, activityID, game)
			editURL = fmt.Sprintf(config.ADMIN_GAME_EDIT_URL, userID, activityID, game, "")
			deleteURL = fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, game)
		}
	} else if prefix == "activity" {
		// 活動頁面
		newURL = fmt.Sprintf(config.ADMIN_ACTIVITY_NEW_URL, userID)
		editURL = fmt.Sprintf(config.ADMIN_ACTIVITY_EDIT_URL, userID, "")
		deleteURL = config.ACTIVITY_API_URL
		gameURL = fmt.Sprintf(config.ADMIN_GAME_TYPE_URL, userID, "")
	} else {
		// 用戶、權限、菜單、總覽頁面
		newURL = fmt.Sprintf(config.ADMIN_NEW_URL, prefix)
		deleteURL = fmt.Sprintf(config.ADMIN_API_URL, prefix)

		if prefix == "manager" {
			editURL = fmt.Sprintf(config.ADMIN_MANAGER_EDIT_URL, prefix)
			activityURL = fmt.Sprintf(config.ADMIN_ACTIVITY_URL, "")
		} else {
			editURL = fmt.Sprintf(config.ADMIN_EDIT_URL, prefix)
		}
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User: user,
		// PanelInfo: panelInfo,
		Route: route{
			New:      newURL,  // 新增頁面
			Edit:     editURL, // 編輯頁面
			DELETE:   deleteURL,
			Activity: activityURL,
			Game:     gameURL,
		},
		Token: auth.AddToken(user.UserID),
	})
	return
}

// panel, _ := h.GetTable(ctx, "admin_"+prefix)
// params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
// 	[]string{"id"}, []string{"asc"}).SetPKs(id)
// formInfo, err := panel.GetEditFormInfo(params, h.services, []string{"id"})
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
// 	return
// }

// panel, _ := h.GetTable(ctx, "admin_"+prefix)
// params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
// 	[]string{"id"}, []string{"asc"}).SetPKs(user.UserID)

// panel, _ := h.GetTable(ctx, "admin_"+prefix)
// param := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
// 	[]string{"id"}, []string{"asc"})
// if panelInfo, err = panel.GetData(param, h.services); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
// 	retur
