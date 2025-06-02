package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowGame 遊戲頁面 GET API
func (h *Handler) ShowGame(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		liffState  = ctx.Query("liff.state") // line裝置，liff url會顯示此參數，ex: liff.state=?activity_id=xxx&game_id=xxx
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		role       = ctx.Query("role") // 角色
		route      string
		// tempName, htmlTmpl string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// 處理activity_id、game_id
	if activityID == "" || gameID == "" {
		if liffState != "" {
			if len(liffState) > 13 {
				liffState = liffState[13:] // 不讀取?activity_id字串
			}

			// qrcode都是給玩家掃描用的
			role = "guest"

			// liff.state參數處理
			if strings.Contains(liffState, "&game_id=") {
				// ex: liff.state=?activity_id=xxx&game_id=xxx#mst_challenge=xxx
				params := strings.Split(liffState, "&game_id=")
				activityID = params[0]

				if len(params) == 2 {
					params = strings.Split(params[1], "#mst_challenge")
					gameID = params[0]
				}
			}
		}
	}
	if activityID == "" || gameID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID")
		return
	}

	if strings.Contains(path, "redpack") {
		route = "./hilives/hilive/views/game/redpack.html"
	} else if strings.Contains(path, "ropepack") {
		route = "./hilives/hilive/views/game/ropepack.html"
	} else if strings.Contains(path, "whack_mole") {
		route = "./hilives/hilive/views/game/whack_mole.html"
	} else if strings.Contains(path, "3Ddraw_numbers") {
		route = "./hilives/hilive/views/game/3Ddraw_numbers.html"
	} else if strings.Contains(path, "draw_numbers") {
		route = "./hilives/hilive/views/game/draw_numbers.html"
	} else if strings.Contains(path, "lottery") {
		route = "./hilives/hilive/views/game/lottery.html"
	} else if strings.Contains(path, "monopoly") {
		route = "./hilives/hilive/views/game/monopoly.html"
	} else if strings.Contains(path, "QA") {
		route = "./hilives/hilive/views/game/QA.html"
	} else if strings.Contains(path, "tugofwar") {
		route = "./hilives/hilive/views/game/tugofwar.html"
	} else if strings.Contains(path, "bingo") {
		route = "./hilives/hilive/views/game/bingo.html"
	} else if strings.Contains(path, "3DGachaMachine") && role == "host" {
		route = "./hilives/hilive/views/game/3DGachaMachine_host.html"
	} else if strings.Contains(path, "3DGachaMachine") && role == "guest" {
		route = "./hilives/hilive/views/game/3DGachaMachine_guest.html"
	} else if strings.Contains(path, "vote") && role == "host" {
		route = "./hilives/hilive/views/game/vote_host.html"
	} else if strings.Contains(path, "vote") && role == "guest" {
		route = "./hilives/hilive/views/game/vote_guest.html"
	} else if strings.Contains(path, "vote") && role == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識角色資訊，請輸入有效的資料")
		return
	} else if strings.Contains(path, "3DGachaMachine") && role == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識角色資訊，請輸入有效的資料")
		return
	}

	h.executeHTML(ctx, route, executeParam{})
}

// ShowGuestGameInfo 遊戲互動裡的遊戲資訊頁面(手機用戶端)，顯示所有場次資訊 GET API
func (h *Handler) ShowGuestGameInfo(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		// panelInfo  table.PanelInfo
		htmlTmpl string
		// err        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	gameURL := fmt.Sprintf(config.GAME_URL, prefix, activityID) // 進入遊戲按鈕

	if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/redpack.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/ropepack.html"
	} else if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/whack_mole.html"
	} else if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/lottery.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/monopoly.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/QA.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/draw_numbers.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/tugofwar.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/bingo.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/gacha_machine.html"
		gameURL = fmt.Sprintf(config.GAME_ROLE_URL, prefix, activityID, "guest") // 進入遊戲按鈕
	} else if prefix == "vote" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/vote.html"
		gameURL = fmt.Sprintf(config.GAME_ROLE_URL, prefix, activityID, "guest") // 進入遊戲按鈕
	}

	// 所有遊戲資料
	// panel, _ := h.GetTable(ctx, prefix)
	// params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
	// 	[]string{"id"}, []string{"asc"})
	// params.SetField("game", prefix) // 遊戲種類判斷參數
	// if panelInfo, err = panel.GetData(params, h.services); err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
	// 	return
	// }

	// fmt.Println("panelInfo.InfoList: ", panelInfo.InfoList)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		// PanelInfo: panelInfo,
		Route: route{
			// Chatroom: fmt.Sprintf(config.GUEST_CHATROOM_URL, activityID),
			// Info:     fmt.Sprintf("%s/?activity_id=%s", config.GUEST_INFO_URL, activityID),
			// GameInfo: fmt.Sprintf("%s?activity_id=%s", config.GUEST_GAME_URL, activityID),
			Game:    gameURL,                                                        //進入遊戲
			Winning: fmt.Sprintf(config.GUEST_GAME_WINNING_URL, prefix, activityID), // 中獎人員資訊(手機用戶端)
		}})
}

// ShowInteractGame 遊戲資訊頁面(平台)，顯示所有場次資訊 GET API
func (h *Handler) ShowInteractGame(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panel      table.Table
		// params     parameter.Parameters
		// panelInfo  table.PanelInfo
		// formInfo   table.FormInfo
		canAdd, canEdit, canDelete bool
		htmlTmpl                   string
		// err                        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/cms/lottery_game_info.html"
	} else if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/cms/redpack_game_info.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/cms/ropepack_game_info.html"
	} else if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/cms/whackmole_game_info.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/cms/monopoly_game.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/cms/QA_game.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/cms/drawnumbers_game_info.html"
	} else if prefix == "setting" {
		// 遊戲設置
		htmlTmpl = "./hilives/hilive/views/cms/game_options.html"
	} else if prefix == "link" {
		// 遊戲連結
		htmlTmpl = "./hilives/hilive/views/cms/game_href.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/cms/tugofwar_info.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/cms/bingo_info.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/cms/gacha_machine_info.html"
	}

	// 遊戲連結、遊戲基本設置頁面不用panel
	if prefix != "link" && prefix != "setting" {
		// panel, _ = h.GetTable(ctx, prefix)
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

	// if prefix != "setting" {
	// 	params.SetField("game", prefix) // 遊戲種類判斷參數
	// } else if prefix == "setting" {
	// params.SetPKs(activityID)
	// formInfo, err = panel.GetEditFormInfo(params, h.services, []string{"activity_id"})
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
	// 	return
	// }
	// }

	// if prefix != "link" && prefix != "setting" {
	// if panelInfo, err = panel.GetData(params, h.services); err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
	// 	return
	// }
	// }

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:     user,
		UserJSON: utils.JSON(user),
		// PanelInfo: panelInfo,
		// FormInfo:   formInfo,
		ActivityID: activityID,
		Route: route{
			Host:  fmt.Sprintf(config.HOST_GAME_URL, activityID, "false"), // 主持端聊天室頁面
			Prize: fmt.Sprintf(config.INTERACT_PRIZE_URL, prefix, activityID),
			// POST:   fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
			PUT:    fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
			DELETE: fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
			New:    fmt.Sprintf(config.INTERACT_GAME_NEW_URL, prefix, activityID),  // 新增遊戲頁面
			Edit:   fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, prefix, activityID), // 編輯遊戲頁面
			Team:   fmt.Sprintf(config.INTERACT_GAME_TEAM_URL, prefix, activityID), // 隊伍頁面
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// ShowNewInteractGame 遊戲資訊新增頁面(平台)，新建場次 GET API
func (h *Handler) ShowNewInteractGame(ctx *gin.Context) {
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

	if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/cms/whackmole_game_form.html"
	} else if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/cms/lottery_game_form.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/cms/monopoly_game_form.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/cms/QA_game_form.html"
	} else if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/cms/redpack_game_form.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/cms/ropepack_game_form.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/cms/drawnumbers_game_form.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/cms/tugofwar_form.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/cms/bingo_form.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/cms/gacha_machine_form.html"
	}

	// fmt.Println("人數: ", user.MaxGamePeople)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		People:     int(user.MaxGamePeople),
		FormInfo:   panel.GetNewFormInfo(h.services, params, []string{"user_id", "activity_id"}),
		ActivityID: activityID,
		Route: route{
			Method: "post",
			POST:   fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
			Back:   fmt.Sprintf(config.INTERACT_GAME_URL, prefix, activityID),
			Import: config.IMPORT_EXCEL_API_URL,
		},
		Token: auth.AddToken(user.UserID),
	})
}

// ShowEditInteractGame 遊戲資訊編輯頁面(平台)，編輯場次 GET API
func (h *Handler) ShowEditInteractGame(ctx *gin.Context) {
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

	if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/cms/whackmole_game_form.html"
	} else if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/cms/lottery_game_form.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/cms/monopoly_game_form.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/cms/QA_game_form.html"
	} else if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/cms/redpack_game_form.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/cms/ropepack_game_form.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/cms/drawnumbers_game_form.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/cms/tugofwar_form.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/cms/bingo_form.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/cms/gacha_machine_form.html"
	}

	// formInfo, err := panel.GetEditFormInfo(params, h.services, []string{"game_id"})
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
	// 	return
	// }

	// 上限人數
	gameModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
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
			PUT:    fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
			Back:   fmt.Sprintf(config.INTERACT_GAME_URL, prefix, activityID),
			Import: config.IMPORT_EXCEL_API_URL,
		},
		Token: auth.AddToken(user.UserID),
	})
}

// ShowInteractGamePK 遊戲PK資訊頁面(平台)，顯示所有PK資訊 GET API
// func (h *Handler) ShowInteractGamePK(ctx *gin.Context) {
// 	var (
// 		host                       = ctx.Request.Host
// 		path                       = ctx.Request.URL.Path
// 		prefix                     = ctx.Param("__prefix")
// 		activityID                 = ctx.Query("activity_id")
// 		gameID                     = ctx.Query("game_id")
// 		user                       = h.GetLoginUser(ctx.Request, "hilive_session")
// 		canAdd, canEdit, canDelete bool
// 		htmlTmpl                   string
// 	)
// 	if host == config.API_URL {
// 		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
// 		return
// 	}
// 	if activityID == "" || gameID == "" {
// 		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的ID")
// 		return
// 	}

// 	if prefix == "bingo" {
// 		// 遊戲連結
// 		htmlTmpl = "./hilives/hilive/views/cms/bingo_staff.html"
// 	}

// 	// 判斷是否有新增編輯刪除權限
// 	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

// 	h.executeHTML(ctx, htmlTmpl, executeParam{
// 		User:       user,
// 		UserJSON:   utils.JSON(user),
// 		ActivityID: activityID,
// 		Route: route{
// 			Host: fmt.Sprintf(config.HOST_GAME_URL, activityID, "false"), // 主持端聊天室頁面
// 		},
// 		Token:     auth.AddToken(user.UserID),
// 		CanAdd:    canAdd,
// 		CanEdit:   canEdit,
// 		CanDelete: canDelete,
// 	})
// }

// ShowInteractGameTeam 遊戲隊伍資訊頁面(平台)，顯示所有隊伍資訊 GET API
func (h *Handler) ShowInteractGameTeam(ctx *gin.Context) {
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

	if prefix == "tugofwar" {
		// 遊戲連結
		htmlTmpl = "./hilives/hilive/views/cms/tugofwar_team.html"
	}

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		UserJSON:   utils.JSON(user),
		ActivityID: activityID,
		Route: route{
			Host: fmt.Sprintf(config.HOST_GAME_URL, activityID, "false"), // 主持端聊天室頁面
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// ShowGuestGame 遊戲互動主頁面(手機用戶端)，選擇遊戲種類 GET API
// func (h *Handler) ShowGuestGame(ctx *gin.Context) {
// var (
// 	host               = ctx.Request.Host
// 	activityID         = ctx.Query("activity_id")
// 	tempName, htmlTmpl string
// )
// fmt.Println("遊戲互動主頁面")
// if host == config.API_URL {
// 	h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
// 	return
// }
// if activityID == "" {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
// 	return
// }

// h.execute(ctx, tempName, htmlTmpl, executeParam{
// 	Route: route{
// 		Chatroom: fmt.Sprintf( config.GUEST_CHATROOM_URL, activityID),
// 		// Info:      fmt.Sprintf("%s/?activity_id=%s", config.GUEST_INFO_URL, activityID),
// 		// GameInfo:  fmt.Sprintf("%s?activity_id=%s", config.GUEST_GAME_URL, activityID),
// 		Redpack:   fmt.Sprintf("%s/redpack?activity_id=%s", config.GUEST_INFO_URL, activityID),    // 搖紅包場次資訊
// 		Ropepack:  fmt.Sprintf("%s/ropepack?activity_id=%s", config.GUEST_GAME_URL, activityID),   // 套紅包場次資訊
// 		WhackMole: fmt.Sprintf("%s/whack_mole?activity_id=%s", config.GUEST_GAME_URL, activityID), // 打地鼠場次資訊
// 		Lottery:   fmt.Sprintf("%s/lottery?activity_id=%s", config.GUEST_GAME_URL, activityID),    // 遊戲抽獎場次資訊,
// 	}})
// }

// if prefix != "whack_mole" && prefix != "lottery" && prefix != "monopoly" &&
// 		prefix != "QA" && prefix != "draw_numbers" {
// 		h.execute(ctx, "content", htmlTmpl, executeParam{
// 			User:      user,
// 			PanelInfo: panelInfo,

// 			ActivityID: activityID,
// 			Route: route{
// 				Prize:  fmt.Sprintf(config.INTERACT_PRIZE_URL, prefix, activityID),
// 				POST:   fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
// 				PUT:    fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),
// 				DELETE: fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, prefix),

// 				New:  fmt.Sprintf(config.INTERACT_GAME_NEW_URL, prefix, activityID),  // 新增遊戲頁面
// 				Edit: fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, prefix, activityID), // 編輯遊戲頁面

// 			},
// 			Token: auth.AddToken( user.UserID),
// 		})
// 	} else {
// 	}

// if prefix == "QA" {
// 	// 快問快答需要遊戲JSON資料
// 	gameModel, err = models.DefaultGameModel().SetDbConn(h.dbConn).
// 		FindGame(false, gameID)
// 	if err != nil {
// 		h.executeErrorHTML(ctx, "錯誤: 無法取得遊戲資訊，請重新整理頁面")
// 		return
// 	}
// }

// FormInfo:   panel.GetNewFormInfo(h.services, params, []string{"user_id", "activity_id"}),
// Back:   fmt.Sprintf("%s/%s?activity_id=%s", config.INTERACT_GAME_URL, prefix, activityID),
// Reset: fmt.Sprintf(config.INTERACT_GAME_RESET_API_URL_FORM),           // 遊戲重置API
// Back:   fmt.Sprintf("%s/%s?activity_id=%s", config.INTERACT_GAME_URL, prefix, activityID),
// Reset: fmt.Sprintf(config.INTERACT_GAME_RESET_API_URL_FORM),           // 遊戲重置API

// newFormInfo table.FormInfo
// model, err = models.DefaultActivityModel().SetConn(h.conn).Find("activity_id", activityID)
// 新增資料的表單
// newFormInfo := panel.GetNewFormInfo(h.services, params, []string{"activity_id", "user_id"})

// FormInfo:   panel.GetNewFormInfo(h.services, params, []string{"user_id", "activity_id"}),

// if !strings.Contains(path, "draw_numbers") {
// } else if strings.Contains(path, "draw_numbers") {
// 	if activityID == "" {
// 		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
// 		return
// 	}
// }

// template
// if strings.Contains(path, "redpack") {
// 	tempName = "redpack_content"
// 	htmlTmpl = game.RedpackContent
// } else if strings.Contains(path, "ropepack") {
// 	tempName = "ropepack_content"
// 	htmlTmpl = game.RopepackContent
// }

// h.execute(ctx, tempName, htmlTmpl, executeParam{
// 	Role: "我想測試是否能讀到這個參數",
// }
