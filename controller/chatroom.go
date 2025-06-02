package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowHostChatroom 主持人聊天室 GET API
func (h *Handler) ShowHostChatroom(ctx *gin.Context) {
	var (
		host          = ctx.Request.Host
		activityID    = ctx.Query("activity_id")
		login         = ctx.Query("login")
		device        = ctx.Query("device")
		activityModel models.ActivityModel
		user          models.LoginUser
		redpacks      = make([]models.GameModel, 0)
		ropepacks     = make([]models.GameModel, 0)
		whackMoles    = make([]models.GameModel, 0)
		lotterys      = make([]models.GameModel, 0)
		monopolys     = make([]models.GameModel, 0)
		qas           = make([]models.GameModel, 0)
		drawNumbers   = make([]models.GameModel, 0)
		tugofwars     = make([]models.GameModel, 0)
		bingos        = make([]models.GameModel, 0)
		gachaMachines = make([]models.GameModel, 0)
		votes         = make([]models.GameModel, 0)
		params        = url.Values{}
		// ip            = utils.ClientIP(ctx.Request)
		applySignQRcode, redpackQRcode, ropepackQRcode, whackMoleQRcode,
		htmlTmpl, lotteryQRcode, questionQRcode, monopolyQRcode, qaQRcode, tugofwarQRcode, bingoQRcode, gachaMachineQRcode, signnameQRcode, voteQRcode string
		err error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 判斷為主持人時，取得所有搖號抽獎場次中獎人員資料，確保寫入redis中(draw_numbers_winning_staffs_活動ID)
	h.getDrawNumbersAllWinningStaffs(true, activityID)

	activityModel, err = h.getActivityInfo(false, activityID)
	if err != nil || activityModel.ID == 0 {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if activityModel.Device == "line" {
		// 只開啟line裝置驗證，liff url
		applySignQRcode = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, activityID)

		redpackQRcode = fmt.Sprintf(config.HILIVES_REDPACK_GAME_LIFF_URL, activityID)
		ropepackQRcode = fmt.Sprintf(config.HILIVES_ROPEPACK_GAME_LIFF_URL, activityID)
		whackMoleQRcode = fmt.Sprintf(config.HILIVES_WHACK_MOLE_GAME_LIFF_URL, activityID)
		lotteryQRcode = fmt.Sprintf(config.HILIVES_LOTTERY_GAME_LIFF_URL, activityID)
		monopolyQRcode = fmt.Sprintf(config.HILIVES_MONOPOLY_GAME_LIFF_URL, activityID)
		qaQRcode = fmt.Sprintf(config.HILIVES_QA_GAME_LIFF_URL, activityID)
		tugofwarQRcode = fmt.Sprintf(config.HILIVES_TUGOFWAR_GAME_LIFF_URL, activityID)
		bingoQRcode = fmt.Sprintf(config.HILIVES_BINGO_GAME_LIFF_URL, activityID)
		gachaMachineQRcode = fmt.Sprintf(config.HILIVES_GACHAMACHINE_GAME_LIFF_URL, activityID)
		voteQRcode = fmt.Sprintf(config.HILIVES_VOTE_GAME_LIFF_URL, activityID)

		questionQRcode = fmt.Sprintf(config.HILIVES_QUESTION_LIFF_URL, activityID)
		signnameQRcode = fmt.Sprintf(config.HILIVES_SIGNNAME_LIFF_URL, activityID)
	} else {
		// 開啟多個裝置驗證，一般url
		applySignQRcode = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, host, activityID, "")

		redpackQRcode = fmt.Sprintf(config.HTTPS_REDPACK_GAME_URL, host, activityID)
		ropepackQRcode = fmt.Sprintf(config.HTTPS_ROPEPACK_GAME_URL, host, activityID)
		whackMoleQRcode = fmt.Sprintf(config.HTTPS_WHACK_MOLE_GAME_URL, host, activityID)
		lotteryQRcode = fmt.Sprintf(config.HTTPS_LOTTERY_GAME_URL, host, activityID)
		monopolyQRcode = fmt.Sprintf(config.HTTPS_MONOPOLY_GAME_URL, host, activityID)
		qaQRcode = fmt.Sprintf(config.HTTPS_QA_GAME_URL, host, activityID)
		tugofwarQRcode = fmt.Sprintf(config.HTTPS_TUGOFWAR_GAME_URL, host, activityID)
		bingoQRcode = fmt.Sprintf(config.HTTPS_BINGO_GAME_URL, host, activityID)
		gachaMachineQRcode = fmt.Sprintf(config.HTTPS_GACHAMACHINE_GAME_URL, host, activityID, "guest")
		voteQRcode = fmt.Sprintf(config.HTTPS_VOTE_GAME_URL, host, activityID, "guest")

		questionQRcode = fmt.Sprintf(config.HTTPS_QUESTION_URL, host, activityID)
		signnameQRcode = fmt.Sprintf(config.HTTPS_SIGNNAME_URL, host, activityID)
	}

	if device == "mobile" {
		// 手機端
		htmlTmpl = "./hilives/hilive/views/chatroom/control.html"
	} else {
		// 電腦端
		htmlTmpl = "./hilives/hilive/views/chatroom/receiver/index.html"
	}

	// 電腦端
	if device != "mobile" {
		games, err := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAll(activityID, "")
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識遊戲資訊，請輸入有效的資料")
			return
		}

		// 簽到審核開關
		if activityModel.SignCheck == "open" {
			applySignQRcode += "&sign=open"
		}

		// 處理遊戲資料
		for i := 0; i < len(games); i++ {
			if games[i].Game == "redpack" {
				redpacks = append(redpacks, games[i])
			} else if games[i].Game == "ropepack" {
				ropepacks = append(ropepacks, games[i])
			} else if games[i].Game == "whack_mole" {
				whackMoles = append(whackMoles, games[i])
			} else if games[i].Game == "lottery" {
				lotterys = append(lotterys, games[i])
			} else if games[i].Game == "monopoly" {
				monopolys = append(monopolys, games[i])
			} else if games[i].Game == "QA" {
				qas = append(qas, games[i])
			} else if games[i].Game == "draw_numbers" {
				drawNumbers = append(drawNumbers, games[i])
			} else if games[i].Game == "tugofwar" {
				tugofwars = append(tugofwars, games[i])
			} else if games[i].Game == "bingo" {
				bingos = append(bingos, games[i])
			} else if games[i].Game == "3DGachaMachine" {
				gachaMachines = append(gachaMachines, games[i])
			} else if games[i].Game == "vote" {
				votes = append(votes, games[i])
			}
		}

		if login == "false" {
			// 需要登入才可執行聊天室
			if activityModel.LoginRequired == "open" {
				h.executeErrorHTML(ctx, "錯誤: 無法開啟主持端聊天室，如要開啟聊天室，請開啟是否啟用自定義密碼登入才可進入聊天室")
				return
			}

			// 不需要驗證管理員，但需要活動主辦方資訊
			userModel, err := models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(true, true,
					"users.user_id", activityModel.UserID)
			if err != nil || userModel.UserID == "" {
				h.executeErrorHTML(ctx, "錯誤: 無法辨識活動主辦方資訊，請重新操作")
				return
			}

			user.UserID = userModel.UserID
			user.Name = userModel.Name
			user.Phone = userModel.Phone
			user.Email = userModel.Email
			user.Avatar = userModel.Avatar
			// user.Bind = userModel.Bind
			user.Cookie = userModel.Cookie
			// *****可以同時登入(確定拿除)*****
			// user.Ip = userModel.Ip
			// *****可以同時登入(確定拿除)*****
			user.MaxActivity = userModel.MaxActivity
			user.Permissions = userModel.Permissions     // 權限
			user.Activitys = userModel.Activitys         // 活動資訊(包含活動權限)
			user.ActivityMenus = userModel.ActivityMenus // 菜單
			user.LineBind = userModel.LineBind
			user.FbBind = userModel.FbBind
			user.GmailBind = userModel.GmailBind
		} else {
			// 需要驗證管理員，設置chatroom_session cookie(進入遊戲需要chatroom_session驗證)
			user = h.GetLoginUser(ctx.Request, "hilive_session")
			login = "true"

			// 用戶資訊
			params.Add("activity_id", activityID) // 活動ID，判斷用
			params.Add("user_id", user.UserID)
			// params.Add("name", user.Name)
			// params.Add("avatar", user.Avatar)
			// params.Add("phone", user.Phone)
			// params.Add("email", user.Email)
			// params.Add("bind", user.Bind)
			// params.Add("cookie", user.Cookie)
			// *****可以同時登入(確定拿除)*****
			// params.Add("ip", ip)
			// *****可以同時登入(確定拿除)*****
			params.Add("table", "users")
			params.Add("sign", utils.UserSign(user.UserID, config.ChatroomCookieSecret))

			// 設置cookie(session頁面、select頁面也有cookie的設置)
			ctx.SetSameSite(4)
			ctx.SetCookie("chatroom_session", string(utils.Encode([]byte(params.Encode()))),
				config.GetSessionLifeTime(), "/", "", true, false)
		}

		// 設置辨識遠端遙控的session資訊(hilive_remote)
		// ctx.SetCookie("hilive_control", utils.UUID(10),
		// 	config.GetSessionLifeTime(), "/", "", true, false)
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		ActivityModel: activityModel,
		User:          user,
		UserJSON:      utils.JSON(user),
		Route: route{
			Index:              config.INDEX_URL,
			Activity:           fmt.Sprintf(config.ACTIVITY_URL, "true"),
			General:            fmt.Sprintf(config.HOST_GENERAL_SIGNWALL_URL, activityID, login),
			Threed:             fmt.Sprintf(config.HOST_THREED_SIGNWALL_URL, activityID, login),
			Question:           fmt.Sprintf(config.HOST_QUESTION_URL, activityID, login),
			QuestionQRcode:     questionQRcode,
			ApplySign:          applySignQRcode,
			Redpack:            fmt.Sprintf(config.REDPACK_GAME_URL, activityID, login),
			RedpackQRcode:      redpackQRcode,
			Ropepack:           fmt.Sprintf(config.ROPEPACK_GAME_URL, activityID, login),
			RopepackQRcode:     ropepackQRcode,
			WhackMole:          fmt.Sprintf(config.WHACK_MOLE_GAME_URL, activityID, login),
			WhackMoleQRcode:    whackMoleQRcode,
			Lottery:            fmt.Sprintf(config.LOTTERY_GAME_URL, activityID, login),
			LotteryQRcode:      lotteryQRcode,
			DrawNumbers:        fmt.Sprintf(config.DRAW_NUMBERS_GAME_URL, activityID, login),
			ThreedDrawNumbers:  fmt.Sprintf(config.THREED_DRAW_NUMBERS_GAME_URL, activityID, login),
			Monopoly:           fmt.Sprintf(config.MONOPOLY_GAME_URL, activityID, login),
			MonopolyQRcode:     monopolyQRcode,
			QA:                 fmt.Sprintf(config.QA_GAME_URL, activityID, login),
			QAQRcode:           qaQRcode,
			Tugofwar:           fmt.Sprintf(config.TUGOFWAR_GAME_URL, activityID, login),
			TugofwarQRcode:     tugofwarQRcode,
			Bingo:              fmt.Sprintf(config.BINGO_GAME_URL, activityID, login),
			BingoQRcode:        bingoQRcode,
			GachaMachine:       fmt.Sprintf(config.GACHAMACHINE_GAME_URL, activityID, login, "host"),
			GachaMachineQRcode: gachaMachineQRcode,
			Vote:               fmt.Sprintf(config.VOTE_GAME_URL, activityID, login, "host"),
			VoteQRcode:         voteQRcode,
			Signname:           fmt.Sprintf(config.HOST_SIGNNAME_SIGNWALL_URL, activityID, login),
			SignnameQRcode:     signnameQRcode,
		},
		Chatroom: Chatroom{
			Redpacks:      redpacks,
			Ropepacks:     ropepacks,
			WhackMoles:    whackMoles,
			Lotterys:      lotterys,
			Monopolys:     monopolys,
			QAs:           qas,
			DrawNumbers:   drawNumbers,
			Tugofwars:     tugofwars,
			Bingos:        bingos,
			GachaMachines: gachaMachines,
			Votes:         votes,
		},
	})
	return
}

// ShowGuest 客戶端首頁 GET API
func (h *Handler) ShowGuest(ctx *gin.Context) {
	var (
		host                = ctx.Request.Host
		activityID          = ctx.Query("activity_id")
		user                models.LoginUser // 活動主辦方資訊
		loginUser           = h.GetLoginUser(ctx.Request, "chatroom_session")
		applySignModel, err = models.DefaultApplysignModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(0, activityID, loginUser.UserID, true)
		qrcodeURL string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	} else if host == config.HILIVES_NET_URL {
		qrcodeURL = fmt.Sprintf(config.HILIVES_QRCODE_LIFF_URL, activityID)
	}

	if err != nil || activityID == "" || applySignModel.ID == 0 {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊或用戶資訊，請輸入有效的活動ID")
		return
	}

	// 取得活動主辦方資訊
	// *****可以同時登入(新)*****
	userModel, err := models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, true,
			"users.user_id", applySignModel.ActivityUserID)
	// *****可以同時登入(新)*****

	// *****可以同時登入(暫時拿除)*****
	// userModel, err := models.DefaultUserModel().SetDbConn(h.dbConn).
	// 	SetRedisConn(h.redisConn).Find(true, true, "",
	// 	"users.user_id", applySignModel.ActivityUserID)
	// *****可以同時登入(暫時拿除)*****
	if err != nil || userModel.UserID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動主辦方資訊，請重新操作")
		return
	}
	user.UserID = userModel.UserID
	user.Name = userModel.Name
	user.Phone = userModel.Phone
	user.Email = userModel.Email
	user.Avatar = userModel.Avatar
	// user.Bind = userModel.Bind
	user.Cookie = userModel.Cookie
	// *****可以同時登入(確定拿除)*****
	// user.Ip = userModel.Ip
	// *****可以同時登入(確定拿除)*****
	user.MaxActivity = userModel.MaxActivity
	user.Permissions = userModel.Permissions     // 權限
	user.Activitys = userModel.Activitys         // 活動資訊(包含活動權限)
	user.ActivityMenus = userModel.ActivityMenus // 菜單
	user.LineBind = userModel.LineBind
	user.FbBind = userModel.FbBind
	user.GmailBind = userModel.GmailBind

	// fmt.Println("用戶端聊天室頁面user.ActivityMenus參數: ", user.ActivityMenus, "user.Ip: ", user.Ip)

	h.executeHTML(ctx, "./hilives/hilive/views/chatroom/style/default/activity_room.html", executeParam{
		ApplysignModel: applySignModel,
		User:           loginUser,
		UserJSON:       utils.JSON(user),
		Route: route{
			Chatroom: fmt.Sprintf(config.GUEST_CHATROOM_URL, activityID),
			POST:     config.CHATROOM_RECORD_API_URL,

			Introduce: fmt.Sprintf(config.GUEST_INTRODUCE_URL, activityID),
			Schedule:  fmt.Sprintf(config.GUEST_SCHEDULE_URL, activityID),
			Guest:     fmt.Sprintf(config.GUEST_GUEST_URL, activityID),
			Material:  fmt.Sprintf(config.GUEST_MATERIAL_URL, activityID),

			Lottery:      fmt.Sprintf(config.GUEST_LOTTERY_GAME_URL, activityID),
			Redpack:      fmt.Sprintf(config.GUEST_REDPACK_GAME_URL, activityID),      // 搖紅包場次資訊
			Ropepack:     fmt.Sprintf(config.GUEST_ROPEPACK_GAME_URL, activityID),     // 套紅包場次資訊
			WhackMole:    fmt.Sprintf(config.GUEST_WHACKMOLE_GAME_URL, activityID),    // 打地鼠場次資訊
			Monopoly:     fmt.Sprintf(config.GUEST_MONOPOLY_GAME_URL, activityID),     // 超級大富翁場次資訊
			QA:           fmt.Sprintf(config.GUEST_QA_GAME_URL, activityID),           // 快問快答場次資訊
			Tugofwar:     fmt.Sprintf(config.GUEST_TUGOFWAR_GAME_URL, activityID),     // 拔河遊戲場次資訊
			Bingo:        fmt.Sprintf(config.GUEST_BINGO_GAME_URL, activityID),        // 賓果遊戲場次資訊
			GachaMachine: fmt.Sprintf(config.GUEST_GACHAMACHINE_GAME_URL, activityID), // 扭蛋機遊戲場次資訊
			Vote:         fmt.Sprintf(config.GUEST_VOTE_GAME_URL, activityID),         // 投票遊戲場次資訊

			Question: fmt.Sprintf(config.GUEST_QUESTION_URL, activityID), // 提問牆
			Signname: fmt.Sprintf(config.GUEST_SIGNNAME_URL, activityID), // 簽名牆
			Winning:  fmt.Sprintf(config.GUEST_WINNING_URL, activityID),  // 個人中獎紀錄頁面
			QRcode:   qrcodeURL,                                          // QRcode頁面
		},
	})
}

// ShowGuestChatroom 客戶端聊天室 GET API
func (h *Handler) ShowGuestChatroom(ctx *gin.Context) {
	var (
		host                = ctx.Request.Host
		activityID          = ctx.Query("activity_id")
		user                models.LoginUser // 活動主辦方資訊
		loginUser           = h.GetLoginUser(ctx.Request, "chatroom_session")
		applySignModel, err = models.DefaultApplysignModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Find(0, activityID, loginUser.UserID, true)
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if err != nil || activityID == "" || applySignModel.ID == 0 {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊或用戶資訊，請輸入有效的活動ID")
		return
	}

	// 取得活動主辦方資訊
	// *****可以同時登入(新)*****
	userModel, err := models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, true,
			"users.user_id", applySignModel.ActivityUserID)
	if err != nil || userModel.UserID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動主辦方資訊，請重新操作")
		return
	}
	// *****可以同時登入(新)*****

	// *****可以同時登入(暫時拿除)*****
	// userModel, err := models.DefaultUserModel().SetDbConn(h.dbConn).
	// 	SetRedisConn(h.redisConn).Find(true, true, "",
	// 	"users.user_id", applySignModel.ActivityUserID)
	// if err != nil || userModel.UserID == "" {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動主辦方資訊，請重新操作")
	// 	return
	// }
	// *****可以同時登入(暫時拿除)*****
	user.UserID = userModel.UserID
	user.Name = userModel.Name
	user.Phone = userModel.Phone
	user.Email = userModel.Email
	user.Avatar = userModel.Avatar
	// user.Bind = userModel.Bind
	user.Cookie = userModel.Cookie
	// *****可以同時登入(確定拿除)*****
	// user.Ip = userModel.Ip
	// *****可以同時登入(確定拿除)*****
	user.MaxActivity = userModel.MaxActivity
	user.Permissions = userModel.Permissions     // 權限
	user.Activitys = userModel.Activitys         // 活動資訊(包含活動權限)
	user.ActivityMenus = userModel.ActivityMenus // 菜單
	user.LineBind = userModel.LineBind
	user.FbBind = userModel.FbBind
	user.GmailBind = userModel.GmailBind

	// fmt.Println("是否為黑名單: ", h.IsBlackStaff(activityID, "", "message", loginUser.UserID), loginUser.UserID)
	h.executeHTML(ctx, "./hilives/hilive/views/chatroom/style/default/chatroom.html", executeParam{
		// ApplysignModel: applySignModel,
		// User:           loginUser,
		// UserJSON:       utils.JSON(user),

		// IsBlack:        h.IsBlackStaff(activityID, "", "message", loginUser.UserID),
		Route: route{
			POST: config.CHATROOM_RECORD_API_URL,
		},
	})
}

// ShowGuestQRcode 客戶端QRcode頁面 GET API
func (h *Handler) ShowGuestQRcode(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		liffState  = ctx.Query("liff.state") // ex: liff.state=?activity_id=xxx
		activityID = ctx.Query("activity_id")
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	// 處理activity_id
	if activityID == "" && liffState != "" {
		if len(liffState) > 13 {
			liffState = liffState[13:] // 不讀取?activity_id字串
		}

		// liff.state參數處理
		// ex: liff.state=?activity_id=xxx&liff.referrer=xxx#mst_challenge=xxx
		params := strings.Split(liffState, "&liff.referrer")
		activityID = params[0]
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊或用戶資訊，請輸入有效的活動ID")
		return
	}
	h.executeHTML(ctx, "./hilives/hilive/views/chatroom/style/default/scan.html", executeParam{
		Route: route{
			Back: fmt.Sprintf(config.GUEST_URL, activityID),
		},
	})
}

// activityModel, err = models.DefaultActivityModel().SetDbConn(h.dbConn).Find("activity_id", activityID)
// records, _          = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).LeftJoinLineUsers(activityID)
// ActivityModel: activityModel,
// Number:         int(applySignModel.Number),
// Chatroom: Chatroom{
// 	// Records: records,
// },

// Chatroom: fmt.Sprintf("%s/?activity_id=%s", config.GUEST_CHATROOM_URL, activityID),
// Info:            fmt.Sprintf("%s/?activity_id=%s", config.GUEST_INFO_URL, activityID),
// GameInfo:        fmt.Sprintf("%s?activity_id=%s", config.GUEST_GAME_URL, activityID),
// Overviews: overviews,
// Records:    records,

// records, _         = models.DefaultChatroomRecordModel().SetDbConn(h.dbConn).LeftJoinLineUsers(activityID)
// htmlTmpl = receiver.Page + receiver.Head + receiver.Body + receiver.Script +
// 	receiver.Loading + receiver.Content + receiver.Interface
// general     = fmt.Sprintf(config.HOST_GENERAL_SIGNWALL_URL, activityID)
// threed      = fmt.Sprintf(config.HOST_THREED_SIGNWALL_URL, activityID)
// question    = fmt.Sprintf(config.HOST_QUESTION_URL, activityID)
// redpack     = fmt.Sprintf(config.REDPACK_GAME_URL, activityID)
// ropepack    = fmt.Sprintf(config.ROPEPACK_GAME_URL, activityID)
// whackMole   = fmt.Sprintf(config.WHACKMOLE_GAME_URL, activityID)
// lottery     = fmt.Sprintf(config.LOTTERY_GAME_URL, activityID)
// drawNumber = fmt.Sprintf(config.DRAW_NUMBERS_GAME_URL, activityID)
// monopoly    = fmt.Sprintf(config.MONOPOLY_GAME_URL, activityID)
// qa          = fmt.Sprintf(config.QA_GAME_URL, activityID)

// params.Add("identify", user.Identify)
// params.Add("friend", user.Friend)

// 將用戶資訊加入redis
// if err = h.redisConn.SetCache(config.LINE_USERS+user.UserID, ip,
// 	"EX", config.GetSessionLifeTime()); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }
// if err = h.redisConn.HashSetCache(
// 	"line_users", user.UserID, models.LoginUser{
// 		UserID:   user.UserID,
// 		Name:     user.Name,
// 		Phone:    user.Phone,
// 		Email:    user.Email,
// 		Avatar:   user.Avatar,
// 		Bind:     user.Bind,
// 		Cookie:   user.Cookie,
// 		Identify: "",
// 		Friend:   "",
// 		Table:    "users",
// 		Ip:       ip,
// 	}); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// if err = auth.SetCookie(ctx, user, h.dbConn,
// 	"chatroom_session"); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: cookie發生問題
