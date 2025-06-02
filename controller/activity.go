package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowActivity 新增、編輯活動、快速設置、選擇項目頁面(平台) GET API
func (h *Handler) ShowActivity(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// path      = ctx.Request.URL.Path
		page       = ctx.Param("__page")
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		formInfo   table.FormInfo
		people     int
		// canAdd, canEdit, canDelete bool
		htmlTmpl string
		method   string
		err      error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if page == "new" {
		htmlTmpl = "./hilives/hilive/views/cms/activity_form.html"
		method = "post"

		panel, _ := h.GetTable(ctx, "activity")
		param := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"desc"})
		param.SetPKs(user.UserID)

		formInfo = panel.GetNewFormInfo(h.services, param, []string{"user_id"})

		// 上限人數
		people = int(user.MaxActivityPeople)

	} else if page == "edit" {
		htmlTmpl = "./hilives/hilive/views/cms/activity_form.html"
		method = "put"

		panel, _ := h.GetTable(ctx, "activity")
		param := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"desc"})
		param.SetPKs(user.UserID, activityID)

		if formInfo, err = panel.GetEditFormInfo(param, h.services, []string{"user_id", "activity_id"}); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
			return
		}

		// 上限人數
		activityModel, err := h.getActivityInfo(false, activityID)
		if err != nil || activityModel.ID == 0 {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
			return
		}
		people = int(activityModel.MaxPeople)
	} else if page == "quick_start" {
		htmlTmpl = "./hilives/hilive/views/cms/quick_start.html"
	} else if page == "select" {
		htmlTmpl = "./hilives/hilive/views/cms/select_page.html"
	}

	// 判斷是否有新增編輯刪除權限
	// 判斷是否為試用用戶(活動場次上限為1)，試用用戶沒有新增編輯刪除功能
	// if user.MaxActivity > 1 {
	// 	canAdd = true
	// 	canEdit = true
	// 	canDelete = true
	// }
	// fmt.Println("人數: ", user.MaxActivityPeople)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:     user,
		People:   people,
		FormInfo: formInfo,
		Route: route{
			Method:       method,
			Back:         fmt.Sprintf(config.ACTIVITY_URL, "false"),
			POST:         config.ACTIVITY_API_URL,
			PUT:          config.ACTIVITY_API_URL,
			Overview:     config.OVERVIEW_URL + activityID,                    // 活動總覽頁面
			Activity:     fmt.Sprintf(config.ACTIVITY_URL, "false"),           // 活動頁面
			Select:       fmt.Sprintf(config.ACTIVITY_SELECT_URL, activityID), // 選擇項目頁面，ACTIVITY_SELECT_URL
			Topic:        fmt.Sprintf(config.INTERACT_WALL_URL, "topic", activityID),
			Question:     fmt.Sprintf(config.INTERACT_WALL_URL, "question", activityID),
			SpecialDanmu: fmt.Sprintf(config.INTERACT_WALL_URL, "specialdanmu", activityID),
			Holdscreen:   fmt.Sprintf(config.INTERACT_WALL_URL, "holdscreen", activityID),
			General:      fmt.Sprintf(config.INTERACT_SIGN_URL, "general", activityID),
			Threed:       fmt.Sprintf(config.INTERACT_SIGN_URL, "threed", activityID),
			Signname:     fmt.Sprintf(config.INTERACT_SIGN_URL, "signname", activityID),

			// 遊戲頁面
			Redpack:      fmt.Sprintf(config.GAME_URL, "redpack", activityID),
			Ropepack:     fmt.Sprintf(config.GAME_URL, "ropepack", activityID),
			WhackMole:    fmt.Sprintf(config.GAME_URL, "whack_mole", activityID),
			Lottery:      fmt.Sprintf(config.GAME_URL, "lottery", activityID),
			DrawNumbers:  fmt.Sprintf(config.GAME_URL, "draw_numbers", activityID),
			Monopoly:     fmt.Sprintf(config.GAME_URL, "monopoly", activityID),
			QA:           fmt.Sprintf(config.GAME_URL, "QA", activityID),
			Tugofwar:     fmt.Sprintf(config.GAME_URL, "tugofwar", activityID),
			Bingo:        fmt.Sprintf(config.GAME_URL, "bingo", activityID),
			GachaMachine: fmt.Sprintf(config.GAME_URL, "3DGachaMachine", activityID),

			// 遊戲新增頁面
			RedpackNew:      fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "redpack", activityID),
			RopepackNew:     fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "ropepack", activityID),
			WhackMoleNew:    fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "whack_mole", activityID),
			LotteryNew:      fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "lottery", activityID),
			DrawNumbersNew:  fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "draw_numbers", activityID),
			MonopolyNew:     fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "monopoly", activityID),
			QANew:           fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "QA", activityID),
			TugofwarNew:     fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "tugofwar", activityID),
			BingoNew:        fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "bingo", activityID),
			GachaMachineNew: fmt.Sprintf(config.INTERACT_GAME_NEW_URL, "3DGachaMachine", activityID),

			// 遊戲編輯頁面
			RedpackEdit:      fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "redpack", activityID),
			RopepackEdit:     fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "ropepack", activityID),
			WhackMoleEdit:    fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "whack_mole", activityID),
			LotteryEdit:      fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "lottery", activityID),
			DrawNumbersEdit:  fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "draw_numbers", activityID),
			MonopolyEdit:     fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "monopoly", activityID),
			QAEdit:           fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "QA", activityID),
			TugofwarEdit:     fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "tugofwar", activityID),
			BingoEdit:        fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "bingo", activityID),
			GachaMachineEdit: fmt.Sprintf(config.INTERACT_GAME_EDIT_URL, "3DGachaMachine", activityID),

			// 遊戲獎品頁面
			RedpackPrize:      fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "redpack", activityID),
			RopepackPrize:     fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "ropepack", activityID),
			WhackMolePrize:    fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "whack_mole", activityID),
			LotteryPrize:      fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "lottery", activityID),
			DrawNumbersPrize:  fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "draw_numbers", activityID),
			MonopolyPrize:     fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "monopoly", activityID),
			QAPrize:           fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "QA", activityID),
			TugofwarPrize:     fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "tugofwar", activityID),
			BingoPrize:        fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "bingo", activityID),
			GachaMachinePrize: fmt.Sprintf(config.INTERACT_GAME_PRIZE_URL, "3DGachaMachine", activityID),

			// 遊戲增刪改查api
			RedpackAPI:      fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "redpack"),
			RopepackAPI:     fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "ropepack"),
			WhackMoleAPI:    fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "whack_mole"),
			LotteryAPI:      fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "lottery"),
			DrawNumbersAPI:  fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "draw_numbers"),
			MonopolyAPI:     fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "monopoly"),
			QAAPI:           fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "QA"),
			TugofwarAPI:     fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "tugofwar"),
			BingoAPI:        fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "bingo"),
			GachaMachineAPI: fmt.Sprintf(config.INTERACT_GAME_API_URL_FORM, "3DGachaMachine"),
		},
		Token: auth.AddToken(user.UserID),
		// CanAdd:    canAdd,
		// CanEdit:   canEdit,
		// CanDelete: canDelete,
	})
}

// IsExistActivity 是否存在活動(掃描報名簽到QRCODE) GET API
func (h *Handler) IsExistActivity(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Request.FormValue("activity_id")
		gameID     = ctx.Request.FormValue("game_id")
		sign       = ctx.Request.FormValue("sign")     // 簽到審核開關
		isExist    = ctx.Request.FormValue("is_exist") // 用戶報名簽到資料是否存在
		liffState  = ctx.Query("liff.state")           // line裝置，liff url會顯示此參數，ex: liff.state=?activity_id=xxx&sign=open
		// device     = ctx.Query("device")           // 有此參數時直接顯示驗證碼
		action = "apply"
		// lineURL    = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, activityID)
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// 使用liff qrcode，liff url處理
	// 判斷活動是否存在QRcode(沒有遊戲ID)
	// 可能有sign參數
	if liffState != "" {
		// liff.state參數處理
		var params = make([]string, 0)

		if len(liffState) > 13 {
			liffState = liffState[13:] // 不讀取?activity_id字串
		}

		// ex: 活動ID(&sign=xxx)(&channel_id=xxx)(#mst_challenge=xxx)
		if strings.Contains(liffState, "&channel_id=") && strings.Contains(liffState, "&sign=") {
			// liff url目前沒有設置有channel_id參數的api
			// // log.Println("有channel_id參數，有sign參數")

			// // 包含channel_id、sign參數
			// // 活動ID&sign=xxx&channel_id=xxx(#mst_challenge=xxx)
			// params = strings.Split(liffState, "&sign=") // ex: [活動ID、sign參數&channel_id=xxx(#mst_challenge=xxx)]
			// activityID = params[0]

			// // 取得channel_id(sign參數&channel_id=xxx(#mst_challenge=xxx))
			// params = strings.Split(params[1], "&channel_id=") // ex: [sign參數、頻道ID(#mst_challenge=xxx)]

			// // 判斷簽到參數
			// if params[0] == "open" {
			// 	action = "sign"
			// }

			// channelID = params[1] // 頻道ID(#mst_challenge=xxx)
			// // 判斷字串#mst_challenge
			// // 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&sign=xxx&channel_id=xxx#mst_challenge=xxx
			// if strings.Contains(channelID, "#mst_challenge") {
			// 	channelID = strings.Split(channelID, "#mst_challenge")[0]
			// }
		} else if strings.Contains(liffState, "&channel_id=") {
			// liff url目前沒有設置有channel_id參數的api
			// // log.Println("只有channel_id參數")

			// // 包含channel_id參數
			// // 活動ID&channel_id=xxx(#mst_challenge=xxx)
			// params = strings.Split(liffState, "&channel_id=") // ex: [活動ID、頻道ID(#mst_challenge=xxx)]
			// activityID = params[0]

			// // 取得channel_id(可能含有#mst_challenge=xxx)
			// channelID = params[1]
			// // 判斷字串#mst_challenge
			// // 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&channel_id=xxx#mst_challenge=xxx
			// if strings.Contains(channelID, "#mst_challenge") {
			// 	channelID = strings.Split(channelID, "#mst_challenge")[0]
			// }
		} else if strings.Contains(liffState, "&sign=") {
			// log.Println("只有sign參數")

			// 包含sign參數
			// 活動ID&sign=xxx(#mst_challenge=xxx)
			params = strings.Split(liffState, "&sign=") // ex: [活動ID，sign參數(#mst_challenge=xxx)]
			activityID = params[0]

			// 判斷簽到參數
			if strings.Contains(params[1], "open") {
				action = "sign"
			}
		} else {
			// log.Println("沒有channel_id參數，沒有sign參數")

			// 兩個參數都不包含
			// 活動ID(#mst_challenge=xxx)
			activityID = liffState

			if strings.Contains(activityID, "#mst_challenge") {
				activityID = strings.Split(activityID, "#mst_challenge")[0]
			}
		}

		// 不顯示頁面，直接導向auth redirect url
		ctx.Redirect(http.StatusFound, fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL,
			host, action, "line", "activity_id", activityID, gameID))
		return
	}

	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 判斷活動是否結束
	activityModel, isEnd := models.DefaultActivityModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		IsEnd(activityID)
	if isEnd {
		h.executeErrorHTML(ctx, "錯誤: 活動已結束，謝謝參與")
		return
	}

	// 取得session資料，判斷瀏覽器是否存在用戶資訊，存在用戶資訊則直接導向報名簽到頁面
	cookie, err := ctx.Request.Cookie("chatroom_session")
	if err == nil && cookie.Value != "" {
		// 解碼cookie值
		decode, err := utils.Decode([]byte(cookie.Value))
		if err != nil {
			h.executeErrorHTML(ctx, fmt.Sprintf("錯誤: 無法取得cookie值，無法辨識用戶(有cookie但無法解碼) , %s", err.Error()))
			return
		}

		cookieParams, err := url.ParseQuery(string(decode))
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法解析cookie值，無法辨識用戶")
			return
		}

		if cookieParams.Get("user_id") != "" && cookieParams.Get("activity_id") == activityID {
			// 存在用戶資訊，直接導向報名簽到頁面
			redirectURL := fmt.Sprintf(config.APPLYSIGN_URL,
				activityID, cookieParams.Get("user_id"), gameID)

			// 判斷簽到審核開關
			if sign == "open" {
				redirectURL += "&sign=open"
			}

			// 如果用戶資料不存在則不直接導向報名簽到頁面
			if isExist != "false" {
				// 用戶資料存在，導向報名簽到頁面
				ctx.Redirect(http.StatusFound, redirectURL)
				return
			}
		}
	}

	// 判斷簽到審核開關
	if sign == "open" {
		// liff url加上sign=open參數路徑
		// lineURL += "&sign=open"

		action = "sign"
	}

	// 判斷驗證裝置是否為單一裝置，是的話則直接導向驗證頁面
	if activityModel.Device == "line" {
		// 不顯示頁面，直接導向auth redirect url
		ctx.Redirect(http.StatusFound, fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "line", "activity_id", activityID, gameID))
		return
	} else if activityModel.Device == "facebook" {
		// 不顯示頁面，直接導向auth redirect url
		ctx.Redirect(http.StatusFound, fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "facebook", "activity_id", activityID, gameID))
		return
	} else if activityModel.Device == "gmail" {
		// 不顯示頁面，直接導向auth redirect url
		ctx.Redirect(http.StatusFound, fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "gmail", "activity_id", activityID, gameID))
		return
	}
	//  else if activityModel.Device == "customize" {
	// 自定義頁面，直接顯示驗證碼或報名簽到表單
	// }

	h.executeHTML(ctx, "./hilives/hilive/views/cms/other_verification.html", executeParam{
		Action: action,
		Route: route{
			Line: fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "line", "activity_id", activityID, gameID), // 一般url
			// Line:      lineURL, // liff url
			Facebook:  fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "facebook", "activity_id", activityID, gameID),
			Gmail:     fmt.Sprintf(config.HTTPS_AUTH_REDIRECT_URL, host, action, "gmail", "activity_id", activityID, gameID),
			Customize: config.APPLYSIGN_LOGIN_API_URL,
			// ApplySign: fmt.Sprintf(config.APPLYSIGN_API_URL, "user"),
		},
	})
}

// activity         = ctx.Request.FormValue("activity_id")
// sign             = ctx.Request.FormValue("sign")
// sessionModel  map[string]interface{}
// sessionValues map[string]interface{}
// activityID, userID, redirect, liffURL string
// ok                                    bool
// err                                   error

// else if host == HILIVES_NET {
// 	liffURL = "https://liff.line.me/1654874788-8QNE2Grk"
// } else if host == WWW_HILIVES_NET {
// 	liffURL = "https://liff.line.me/1654874788-3j8y6kGZ"
// }

// // 判斷是否已經存在用戶session資訊
// cookie, err := ctx.Request.Cookie("chatroom_session")
// // 有cookie值
// if err == nil && cookie.Value != "" {
// 	sessionModel, err = db.Table("session").WithConn(h.conn).
// 		Where("session_id", "=", cookie.Value).First()
// 	// 資料表有session資料
// 	if err == nil && sessionModel != nil {
// 		// json解碼資料表session_values欄位資料
// 		err = json.Unmarshal([]byte(sessionModel["session_values"].(string)),
// 			&sessionValues)
// 		if err == nil {
// 			// 取得用戶ID
// 			if userID, ok = sessionValues["chatroom"].(string); ok {
// 				// 取得LINE用戶資料
// 				if user, err := models.DefaultLineModel().SetDbConn(h.conn).
// 					Find("user_id", userID); err == nil && user.ID != int64(0) {
// 					redirect = fmt.Sprintf("%s?activity_id=%s", liffURL, activityID)
// 					// 直接導向報名簽到頁面，不用LINE驗證
// 					if len(params) ==
