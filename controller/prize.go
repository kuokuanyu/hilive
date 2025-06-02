package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowPrize 遊戲獎品頁面(平台) GET API
func (h *Handler) ShowPrize(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// model, err  = models.DefaultGameModel().SetConn(h.conn).Find(gameID)
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
	// if prefix != "draw_numbers" && (activityID == "" || gameID == "") {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID")
	// 	return
	// } else if prefix == "draw_numbers" && activityID == "" {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的及遊戲ID")
	// 	return
	// }

	if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/cms/lottery_prize.html"
	} else if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/cms/redpack_prize.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/cms/ropepack_prize.html"
	} else if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/cms/whack_mole_prize.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/cms/draw_numbers_prize.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/cms/monopoly_prize.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/cms/QA_prize.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/cms/tugofwar_prize.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/cms/bingo_prize.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/cms/gacha_machine_prize.html"
	}

	// if prefix != "draw_numbers" {
	params = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID, gameID)
	newFormInfo = panel.GetNewFormInfo(h.services, params,
		[]string{"user_id", "activity_id", "game_id"}) // 新增資料的表單
	// }
	// else if prefix == "draw_numbers" {
	// 	params = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
	// 		[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID)
	// 	newFormInfo = panel.GetNewFormInfo(h.services, params,
	// 		[]string{"user_id", "activity_id"}) // 新增資料的表單
	// }
	params.SetField("game", prefix) // 遊戲種類判斷參數

	// 所有獎品資料
	// if panelInfo, err = panel.GetData(params, h.services,[]string{"prize_id"}); err != nil {
	if panelInfo, err = panel.GetData(params, h.services); err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
		return
	}

	// 新增資料的表單
	// newFormInfo = panel.GetNewFormInfo(h.services, params,
	// 	[]string{"user_id", "activity_id", "game_id"})

	// 判斷是否有新增編輯刪除權限
	// path := fmt.Sprintf("/admin/interact/game/%s", prefix)
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	// if prefix != "whack_mole" && prefix != "draw_numbers" && prefix != "lottery" &&
	// 	prefix != "monopoly" && prefix != "QA" {
	// 	h.execute(ctx, "prize_content", htmlTmpl, executeParam{
	// 		User:       user,
	// 		PanelInfo:  panelInfo,
	// 		FormInfo:   newFormInfo,
	// 		ActivityID: activityID,
	// 		GameID:     gameID,
	// 		Route: route{
	// 			POST:   fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
	// 			PUT:    fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
	// 			DELETE: fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
	// 		},
	// 		Token:     auth.AddToken(user.UserID),
	// 		CanAdd:    canAdd,
	// 		CanEdit:   canEdit,
	// 		CanDelete: canDelete,
	// 	})
	// } else {
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		PanelInfo:  panelInfo,
		FormInfo:   newFormInfo,
		ActivityID: activityID,
		GameID:     gameID,
		Route: route{
			POST:   fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
			PUT:    fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
			DELETE: fmt.Sprintf(config.INTERACT_PRIZE_API_URL_FORM, prefix),
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
	// }
}
