package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowSignPrize 遊戲獎品頁面(平台) GET API
func (h *Handler) ShowSignPrize(ctx *gin.Context) {
	var (
		host                       = ctx.Request.Host
		path                       = ctx.Request.URL.Path
		prefix                     = ctx.Param("__prefix")
		activityID                 = ctx.Query("activity_id")
		gameID                     = ctx.Query("game_id")
		user                       = h.GetLoginUser(ctx.Request, "hilive_session")
		panelInfo                  table.PanelInfo
		newFormInfo                table.FormInfo
		panel, _                   = h.GetTable(ctx, prefix+"_prize")
		params                     parameter.Parameters
		canAdd, canEdit, canDelete bool
		htmlTmpl                   string
		err                        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" || gameID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID")
		return
	}

	if prefix == "vote" {
		htmlTmpl = "./hilives/hilive/views/cms/vote_prize.html"
	}

	params = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID, gameID)
	newFormInfo = panel.GetNewFormInfo(h.services, params,
		[]string{"user_id", "activity_id", "game_id"}) // 新增資料的表單

	params.SetField("game", prefix) // 遊戲種類判斷參數

	if panelInfo, err = panel.GetData(params, h.services); err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
		return
	}

	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		PanelInfo:  panelInfo,
		FormInfo:   newFormInfo,
		ActivityID: activityID,
		GameID:     gameID,
		Route: route{
			POST:   fmt.Sprintf(config.INTERACT_SIGN_PRIZE_API_URL_FORM, prefix),
			PUT:    fmt.Sprintf(config.INTERACT_SIGN_PRIZE_API_URL_FORM, prefix),
			DELETE: fmt.Sprintf(config.INTERACT_SIGN_PRIZE_API_URL_FORM, prefix),
		},
		Token: auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}
