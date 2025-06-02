package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowQuestion 提問牆頁面(主持端、用戶端) GET API
func (h *Handler) ShowQuestion(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		activityID = ctx.Query("activity_id")
		// channelID  = ctx.Query("channel_id")
		liffState = ctx.Query("liff.state") // line裝置，liff url會顯示此參數，ex: liff.state=?activity_id=xxx
		route     string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		if liffState != "" {
			// line裝置
			if len(liffState) > 13 {
				liffState = liffState[13:] // 不讀取?activity_id字串
			}

			// ex: liff.state=?activity_id=xxx(&channel_id=xxx)#mst_challenge=xxx
			// 判斷是否有channel_id參數
			if strings.Contains(liffState, "&channel_id=") {
				// liff url目前沒有設置有channel_id參數的api
				// // 有channel_id參數
				// // xxx&channel_id=xxx(#mst_challenge=xxx)
				// params := strings.Split(liffState, "&channel_id=") // ex: [活動ID，頻道ID(#mst_challenge=xxx)]
				// activityID = params[0]

				// // 取得channel_id(可能含有#mst_challenge=xxx)
				// channelID = params[1]
				// // 判斷字串#mst_challenge
				// // 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&channel_id=xxx#mst_challenge=xxx
				// if strings.Contains(channelID, "#mst_challenge") {
				// 	channelID = strings.Split(channelID, "#mst_challenge")[0]
				// }
			} else {
				// 沒有有channel_id參數
				// xxx(#mst_challenge=xxx)
				activityID = liffState

				if strings.Contains(activityID, "#mst_challenge") {
					activityID = strings.Split(activityID, "#mst_challenge")[0]
				}
			}
		}
	}

	// 判斷主持、用戶端
	if strings.Contains(path, "host") {
		route = "./hilives/hilive/views/game/question.html"

		// user := h.GetLoginUser(ctx.Request, "hilive_session")
		// if activityModel.UserID != user.UserID {
		// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		// 	return
		// }
	} else if strings.Contains(path, "guest") {
		user := h.GetLoginUser(ctx.Request, "chatroom_session")
		route = "./hilives/hilive/views/chatroom/style/default/question.html"
		// isBlack = h.IsBlackStaff(activityID, "", "question", user.UserID)

		// 主持人的session資訊，重新導向報名簽到驗證
		if user.Table != "line_users" {
			// 驗證發生問題時，導向一般url或者liff url判斷
			redirectURL := ""
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
			redirectURL = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, host, activityID, "") + "&is_exist=false"
			// }

			// 簽到審核開關
			redirectURL += "&sign=open"

			ctx.Redirect(http.StatusFound, redirectURL)
			return
		}
	}
	h.executeHTML(ctx, route, executeParam{
		// IsBlack:       isBlack,
		// ActivityModel: activityModel,
	})
}

// 判斷用戶資訊
// if cookie, err := ctx.Request.Cookie("chatroom_session"); err == nil && cookie.Value != "" {
// 	// 解碼cookie值
// 	decode, err := utils.Decode([]byte(cookie.Value))
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound,
// 			fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 		return
// 	}

// 	params, err := url.ParseQuery(string(decode))
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound,
// 			fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 		return
// 	}

// 	// 主持人的session資訊
// 	if params.Get("table") != "line_users" {
// 		ctx.Redirect(http.StatusFound,
// 			fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 		return
// 	}
// } else {
// 	ctx.Redirect(http.StatusFound,
// 		fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 	return
// }

// user := h.GetLoginUser(ctx.Request, "chatroom_session")
// 取得用戶資訊
// lineModel, err := models.DefaultLineModel().SetDbConn(h.dbConn).Find("user_id", user.UserID)
// if err != nil || lineModel.ID == 0 {
// 	ctx.Redirect(http.StatusFound,
// 		fmt.Sprintf(config.LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", activityID))
// 	return
// }
