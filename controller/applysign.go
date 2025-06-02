package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/views/apply"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
// applysignLock sync.Mutex
)

// Applysign 報名簽到處理 GET API
func (h *Handler) Applysign(ctx *gin.Context) {
	// applysignLock.Lock()
	// defer applysignLock.Unlock()

	var (
		host       = ctx.Request.Host
		ip         = utils.ClientIP(ctx.Request)
		userID     = ctx.Query("user_id")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		sign       = ctx.Query("sign")
		isFirst    = ctx.Query("isfirst")
		qrcode     = ctx.Query("qrcode")     // 郵件的url會有qrocde參數，activity_id=xxx.user_id=xxx
		liffState  = ctx.Query("liff.state") // line裝置，liff url會顯示此參數，ex: liff.state=?activity_id=xxx&user_id=xxx
		// redirectURL = fmt.Sprintf(config.CHATROOM_SESSION_URL, activityID, userID, "guest") // session頁面
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
		redirectURL string
		err         error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// 使用liff qrcode，liff url處理
	// 簡訊、LINE通知會使用到的liff url
	// 只會有活動ID、用戶ID(不會有遊戲ID、頻道ID)
	if liffState != "" {
		// 只要用到這個liff url都是已經報名成功，加上簽到參數
		// sign = "open"

		// liff.state參數處理
		// liff.state=?activity_id=xxx&user_id=xxx
		if len(liffState) > 13 {
			liffState = liffState[13:] // 不讀取?activity_id字串
		}

		// 活動ID&user_id=xxx
		params := strings.Split(liffState, "&user_id=") // ex: [活動ID，用戶ID(#mst_challenge=xxx)]
		activityID = params[0]

		// 判斷user_id
		if len(params) == 2 {
			// 取得user_id(可能含有#mst_challenge=xxx)
			userID = params[1]

			// 判斷字串#mst_challenge
			// 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&user_id=xxx#mst_challenge=xxx
			if strings.Contains(userID, "#mst_challenge") {
				userID = strings.Split(userID, "#mst_challenge")[0]
			}
		} else if len(params) == 1 {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識活動用戶資訊，請輸入有效的資料")
			return
		}
	}

	// 信箱郵件裡的報名簽到qrcode連結(讓主辦方掃描)
	// 只會有活動ID、用戶ID
	if qrcode != "" {
		// ex: qrcode=activity_id=xxx.user_id=xxx

		// 只要用到這個參數都是已經報名成功，加上簽到參數
		sign = "open"

		// qrcode參數處理
		params := strings.Split(qrcode, ".user_id=") // ex: [activity_id=活動ID，用戶ID]

		// 判斷user_id
		if len(params) == 2 {
			userID = params[1]

			// 取得活動ID，params[0]="activity_id=xxx"
			params := strings.Split(params[0], "=") // ex: [activity_id，活動ID]
			if len(params) == 2 {
				activityID = params[1]
			}
		}
	}

	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 用戶資訊
	lineModel, err1 := models.DefaultLineModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, ip, "user_id", userID)
	// if err != nil || lineModel.UserID == "" {
	// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
	// 	return
	// }

	// 原本呼叫三次資料表，改成一次呼叫(也會判斷活動是否結束)
	applysignModel, err2 := models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(0, activityID, userID, true)

	if err1 != nil || err2 != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得報名簽到人員資料，請重新報名簽到")
		return
	} else if applysignModel.ID == 0 || lineModel.UserID == "" || applysignModel.Status == "no" {
		// 無用戶報名簽到資料，導向掃描qrcode頁面
		var action = "apply"
		if sign == "open" {
			action = "sign"
		}

		// 驗證發生問題時，導向一般url或者liff url判斷
		activityModel, err := models.DefaultActivityModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(false, activityID)
		if err != nil || activityModel.ID == 0 {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
			return
		}

		// if activityModel.Device == "line" {
		// liff url
		// redirectURL = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, activityID)
		// } else {
		// 一般url
		redirectURL = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, host, activityID, gameID) + "&is_exist=false"
		// }

		// 判斷簽到審核開關
		if action == "sign" {
			redirectURL += "&sign=open"
		}

		ctx.Redirect(http.StatusFound, redirectURL)
		return
	} else if lineModel.Device == "line" && lineModel.Friend != "yes" && isFirst == "" {
		// 判斷line用戶是否加入官方帳號
		if applysignModel.ActivityLineID != "" || applysignModel.UserLineID != "" {
			var lineID string
			if applysignModel.ActivityLineID != "" {
				lineID = applysignModel.ActivityLineID
			} else if applysignModel.UserLineID != "" {
				lineID = applysignModel.UserLineID
			}
			h.executeErrorHTML(ctx, fmt.Sprintf("錯誤: 未將官方帳號加為好友，請將%s加入好友", lineID))
		} else {
			h.executeErrorHTML(ctx, "錯誤: 未將官方帳號加為好友，請將@hilives加入好友")
		}
		return
	}

	// 判斷是否為黑名單
	if h.IsBlackStaff(activityID, "", "activity", lineModel.UserID) {
		h.executeErrorHTML(ctx, "錯誤: 用戶為黑名單人員，無法加入活動。如有疑問，請聯繫主辦單位")
		return
	}

	var status = applysignModel.Status
	// 判斷是否為管理員
	if applysignModel.ActivityUserID == lineModel.AdminID {
		// 管理員，將資料更新為sign狀態
		if status != "sign" {
			if err = applysignModel.UpdateStatus(true, host, models.EditApplysignModel{
				ActivityID: activityID,
				LineUsers:  []string{lineModel.UserID},
				SignTime:   now.Format("2006-01-02 15:04:05"),
				Status:     "sign"}, sign == "open"); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				h.executeErrorHTML(ctx, err.Error())
				return
			}
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf(config.CHATROOM_SESSION_URL, activityID, userID,gameID))
		return
	}

	// 判斷簽到審核開關
	if status == "apply" { // 已報名成功
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", applysignModel.StartTime, time.Local)
		minutes, _ := time.ParseDuration(fmt.Sprintf("-%dm", applysignModel.SignMinutes)) // 活動前n分鐘可以簽到

		start = start.Add(minutes)
		isSign := now.After(start)
		if applysignModel.SignAllow == "open" || isSign { // 已開放簽到行為
			if applysignModel.SignCheck == "open" && sign != "open" { // 簽到審核開啟
				h.executeErrorHTML(ctx, "錯誤: 活動啟用簽到審核功能，參加人員須掃描簽到的QRcode才能完成簽到動作並參加活動")
				return
			}

			// 可簽到狀態
			status = "sign"

			if err = applysignModel.UpdateStatus(true, host, models.EditApplysignModel{
				ActivityID: activityID,
				LineUsers:  []string{lineModel.UserID},
				SignTime:   now.Format("2006-01-02 15:04:05"),
				Status:     status}, sign == "open"); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				h.executeErrorHTML(ctx, err.Error())
				return
			}
		}
	}

	// 簽到成功，導向session頁面
	// 如果api有qrcode=xxx的參數時，顯示已完成簽到的頁面
	if status == "sign" && qrcode == "" {
		redirectURL = fmt.Sprintf(config.CHATROOM_SESSION_URL, activityID, userID,gameID) // session頁面

		ctx.Redirect(http.StatusFound, redirectURL)
		return
	}

	h.execute(ctx, "template", apply.Template+apply.Content, executeParam{
		Status: status,
	})
}

// ShowApplySign 報名簽到頁面(平台) GET API
func (h *Handler) ShowApplysign(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		panel, _   = h.GetTable(ctx, "applysign")
		formInfo   table.FormInfo
		param      = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID)
		customizeModel             models.CustomizeModel
		canAdd, canEdit, canDelete bool
		applySignURL, htmlTmpl     string
		postURL                    string
		putURL                     = config.USER_APPLYSIGN_API_URL
		deleteURL                  = config.USER_APPLYSIGN_API_URL
		err                        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 自定義欄位資訊
	customizeModel, err = models.DefaultCustomizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(activityID)
	if err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得自定義欄位資訊，請重新整理頁面")
		return
	}
	if customizeModel.Device == "line" {
		// 只開啟line裝置驗證，liff url
		applySignURL = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, activityID)
	} else {
		// 開啟多個裝置驗證，一般url
		applySignURL = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, host, activityID, "")
	}

	if prefix == "apply" {
		// 只取得review、apply、refuse狀態的資料
		param.SetField("status", "review", "apply", "refuse", "no")
		htmlTmpl = "./hilives/hilive/views/cms/activity_apply.html"

		postURL = fmt.Sprintf(config.APPLYSIGN_API_URL, "user") // 自定義匯入報名人員api
	} else if prefix == "sign" {
		// 只取得apply、sign、not sign、cancel狀態的資料
		param.SetField("status", "apply", "sign", "not sign", "cancel")
		htmlTmpl = "./hilives/hilive/views/cms/activity_sign.html"
		applySignURL += "&sign=open"

		postURL = fmt.Sprintf(config.APPLYSIGN_API_URL, "user") // 自定義匯入報名人員api
	} else if prefix == "customize" {
		htmlTmpl = "./hilives/hilive/views/cms/custom_field.html"
	} else if prefix == "qrcode" {
		htmlTmpl = "./hilives/hilive/views/cms/custom_qrcode.html"
	}
	// else if prefix == "user" {
	// // 自定義簽到人員
	// htmlTmpl = "./hilives/hilive/views/cms/activity_custom_sign.html"

	// postURL = fmt.Sprintf(config.APPLYSIGN_API_URL, prefix)
	// putURL = fmt.Sprintf(config.APPLYSIGN_API_URL, prefix)
	// deleteURL = fmt.Sprintf(config.APPLYSIGN_API_URL, prefix)
	// }

	if prefix == "apply" || prefix == "sign" {
		// 面板、表單資料
		if formInfo, err = h.tableList[prefix]().GetSettingFormInfo(param, h.services,
			[]string{"user_id", "activity_id"}); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
			return
		}
	}
	// else if prefix == "customize" {
	// }
	// else if prefix == "qrcode" {
	// }

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:                   user,
		ApplysignCustomizeJSON: utils.JSON(customizeModel),
		SettingFormInfo:        formInfo,
		Route: route{
			ApplySign: applySignURL,
			POST:      postURL,
			PATCH:     fmt.Sprintf(config.APPLYSIGN_API_URL, prefix), // 修改基本設置
			PUT:       putURL,
			DELETE:    deleteURL,
			Export:    fmt.Sprintf(config.APPLYSIGN_API_URL, "export"),
			Import:    config.IMPORT_EXCEL_API_URL,
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// if err = auth.SetCookie(ctx, models.UserModel{UserID: userModel.UserID},
// 	h.dbConn, "chatroom_session"); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: cookie發生問題，請重新報名簽到")
// 	return
// }

// } else if host == HILIVES_NET {
// 	liffURL = "https://liff.line.me/1654874788-2BKEx3Q5"
// } else if host == WWW_HILIVES_NET {
// 	liffURL = "https://liff.line.me/1654874788-xL98weA5"
// }

// liff.state參數處理
// if len(liffState) > 13 {
// 	liffState = liffState[13:] // 不讀取?activity_id字串
// }
// params := strings.Split(liffState, "&sign=")
// activityID = params[0]
// if len(params) == 2 {
// 	if strings.Contains(params[1], "open") {
// 		sign = "open"
// 	}
// } else if len(params) == 1 {
// 	// 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx#mst_challenge=xxx
// 	params = strings.Split(params[0], "#mst_challenge")
// 	activityID = params[0]
// }

// 用戶為活動管理員(已綁定)
// if applysignModel.ActivityUserID == applysignModel.Identify {
// 	ctx.Redirect(http.StatusFound, fmt.Sprintf(config.SELECT_URL, activityID))
// 	return
// }

// 取得活動狀態資訊
// activityModel, err := models.DefaultActivityModel().SetDbConn(h.dbConn).
// 	FindActivityStatus(activityID)
// if err != nil || activityModel.ID == 0 {
// 	h.executeErrorHTML(ctx, "錯誤: 活動已結束，謝謝參與")
// 	return
// }

// 取得用戶資訊
// lineModel, err := models.DefaultLineModel().SetDbConn(h.dbConn).Find("user_id", user.UserID)
// if err != nil || lineModel.ID == 0 {
// 	var action = "apply"
// 	if sign == "open" {
// 		action = "sign"
// 	}
// 	ctx.Redirect(http.StatusFound,
// 		fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, action, "activity_id", activityID))
// 	return
// } else if lineModel.ID != 0 && lineModel.Friend == "no" {
// 	h.executeErrorHTML(ctx, "錯誤: 未將hilives官方帳號加為好友，請將@hilives加入好友")
// 	return
// }

// if activityModel.UserID == lineModel.Identify {
// 	ctx.Redirect(http.StatusFound, fmt.Sprintf(config.SELECT_URL, activityID))
// 	return
// }

// 報名簽到人員資料
// applysignModel, err := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	FindByUserID(activityID, user.UserID)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得報名簽到人員資料，請重新報名簽到")
// 	return
// }

// 判斷報名審核開關
// if applysignModel.ID == 0 { // 無人員資料
// 	// ctx.Redirect(http.StatusFound,
// 	// 	fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 	// return

// 	if activityModel.ApplyCheck == "open" { // 報名審核開啟
// 		status = "review"
// 		message = fmt.Sprintf("已提交%s的活動報名，請等待主辦單位審核", activityModel.ActivityName)
// 	} else { // 報名審核關閉
// 		status = "apply"
// 		message = fmt.Sprintf("完成%s的活動報名，請於活動當天依照主辦單位指示完成活動簽到",
// 			activityModel.ActivityName)
// 	}

// 	// 新增人員資料
// 	id, err := models.DefaultApplysignModel().SetDbConn(h.dbConn).Add(
// 		models.NewApplysignModel{
// 			UserID:     user.UserID,
// 			ActivityID: activityID,
// 			Name:       user.Name,
// 			Avatar:     user.Avatar,
// 			Status:     status,
// 		})
// 	if err != nil {
// 		h.executeErrorHTML(ctx, "錯誤: 新增報名簽到人員發生問題，請重新報名簽到")
// 		return
// 	}
// 	if err = pushMessage(ctx, user.UserID, message); err != nil {
// 		h.executeErrorHTML(ctx, "錯誤: 傳送LINE發生問題")
// 		return
// 	}

// 	applysignModel = models.ApplysignModel{
// 		Base: models.Base{
// 			TableName: config.ACTIVITY_APPLYSIGN_TABLE,
// 			Conn:      h.dbConn,
