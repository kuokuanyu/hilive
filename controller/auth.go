package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	social "hilive/line-login"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"golang.org/x/oauth2/facebook"

	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

var (
	// LoginLock   sync.Mutex
	// LineBotlock sync.Mutex

	oauthStateString = "actionis%sandactivity_idis%sanddeviceis%sandgame_idis%sanduser_idis%s" // state驗證參數，用於防止CSRF攻擊

	// facebook驗證參數
	facebookOauthConfig = &oauth2.Config{
		ClientID:     config.FACEBBOK_ID,
		ClientSecret: config.FACEBOOK_SECRET,
		RedirectURL:  config.FACEBOOK_REDIRECT_URL,
		Scopes:       []string{"email", "public_profile"}, // []string{"email", "public_profile","user_photos", "user_birthday",},
		Endpoint:     facebook.Endpoint,
	}

	// gmail驗證參數
	googleOauthConfig = &oauth2.Config{
		ClientID:     config.GMAIL_ID,
		ClientSecret: config.GMAIL_SECRET,
		RedirectURL:  config.GMAIL_REDIRECT_URL,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		// "https://www.googleapis.com/auth/userinfo.profile" 允许访问用户的基本个人信息，如姓名、头像和公开的 Google 资料信息。
		// "https://www.googleapis.com/auth/userinfo.email" 查看用户的电子邮件地址
		// "https://www.googleapis.com/auth/gmail.readonly" 允许应用程序以只读方式访问用户的 Gmail 数据。
		Endpoint: google.Endpoint,
	}

	// type State []string
// type Host []string
// LineRedirectlock   sync.Mutex
// states = make(State, 0) // state驗證
// hosts  = make(Host, 0)  // 網域
)

// AuthRedirect 導向LINE LOGIN頁面
func (h *Handler) AuthRedirect(ctx *gin.Context) {
	// LoginLock.Lock()
	// defer LoginLock.Unlock()

	var (
		host       = ctx.Request.Host
		liffState  = ctx.Query("liff.state") // line裝置，liff url會顯示此參數，ex: liff.state=?action=xxx，選單會用到
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		userID     = ctx.Query("user_id") // 平台管理員的user_id或自定義人員的user_id
		action     = ctx.Query("action")
		device     = ctx.Query("device")
		scope      = "email profile openid" // email profile openid(用戶資訊)，line login
		botPrompt  = "aggressive"           // normal aggressive(加入官方帳號為好友的顯示方式)，line login
		// prompt                   = "consent"              // 強制執行用戶授權頁面(就算已通過驗證)，line login
		// state                    = social.GenerateNonce() // 用於防止跨站請求偽造的唯一性屬性
		redirectURL                      = fmt.Sprintf(config.HTTPS_AUTH_CALLBACK_URL, host) // callback api驗證url
		lineChannelID, lineChannelSecret string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// line裝置，liff頁面處理(LINE選單用)
	if action == "" && liffState != "" {
		device = "line"

		// liff.state參數處理
		// ex: liff.state=?action=xxx(#mst_challenge=xxx)

		// 處理action參數
		if len(liffState) > 8 {
			liffState = liffState[8:] // 不讀取?action=，xxx(#mst_challenge=xxx)
		}

		if strings.Contains(liffState, "#mst_challenge") {
			// 安卓系統會有這個參數
			action = strings.Split(liffState, "#mst_challenge")[0]
		} else {
			action = liffState
		}
	}

	if action == "" || device == "" {
		h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題")
		return
	}

	// 驗證參數
	// hosts = append(hosts, host)
	// states = append(states, state)

	// 將state驗證參數寫入redis中
	// err := h.redisConn.SetAdd([]interface{}{config.LINE_STATES_REDIS, state})
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: state驗證發生問題")
	// 	return
	// }

	if device == "line" {
		// line裝置
		if action == "bind" {
			// 平台綁定
			if userID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(user_id)")
				return
			}

			// redirectURL += "?user_id=" + userID

			userModel, err := models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(true, true, "users.user_id", userID)
			if err != nil {
				h.executeErrorHTML(ctx, "錯誤: 取得用戶資料發生問題")
				return
			}

			// 判斷是否綁定官方帳號
			if userModel.ChannelID != "" && userModel.ChannelSecret != "" {
				lineChannelID = userModel.ChannelID
				lineChannelSecret = userModel.ChannelSecret
			}
		} else if action == "apply" || action == "sign" {
			// 報名簽到
			if activityID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(activity_id)")
				return
			}

			// 活動報名簽到驗證
			// 取得活動資訊
			activityModel, err := h.getActivityInfo(false, activityID)
			if err != nil || activityModel.ID == 0 {
				h.executeErrorHTML(ctx, "錯誤: 無法取得活動資訊，請重新查詢")
				return
			}

			// 判斷是否綁定官方帳號
			if activityModel.ActivityChannelID != "" && activityModel.ActivityChannelSecret != "" {
				// 活動官方帳號
				lineChannelID = activityModel.ActivityChannelID
				lineChannelSecret = activityModel.ActivityChannelSecret
			} else if activityModel.UserChannelID != "" && activityModel.UserChannelSecret != "" {
				// 用戶活動官方帳號
				lineChannelID = activityModel.UserChannelID
				lineChannelSecret = activityModel.UserChannelSecret
			}
		}
		//  else if action == "search" || action == "create" || action == "case" {
		// 	// LINE 選單
		// }

		// LINE參數為空，補上官方帳號參數
		if lineChannelID == "" || lineChannelSecret == "" {
			lineChannelID = config.CHANNEL_ID
			lineChannelSecret = config.CHANNEL_SECRET
		}

		// 導向line login callback api
		ctx.Redirect(http.StatusFound,
			NewSocialClient(ctx, lineChannelID, lineChannelSecret).GetWebLoinURL(redirectURL,
				fmt.Sprintf(oauthStateString, action, activityID, device, gameID, userID),
				scope, social.AuthRequestOptions{BotPrompt: botPrompt}))
		return
	} else if device == "facebook" {
		if action == "bind" {
			// 平台綁定
			if userID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(user_id)")
				return
			}
		} else {
			// 報名簽到
			if activityID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(activity_id)")
				return
			}
		}

		// facebbok裝置
		ctx.Redirect(http.StatusFound,
			facebookOauthConfig.AuthCodeURL(fmt.Sprintf(oauthStateString, action, activityID, device, gameID, userID)))
		return
	} else if device == "gmail" {
		if action == "bind" {
			// 平台綁定
			if userID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(user_id)")
				return
			}
		} else {
			// 報名簽到
			if activityID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(activity_id)")
				return
			}
		}

		// gmail裝置
		// access_type=online: 适合需要一次性访问用户数据的情况，不需要长期访问。
		// access_type=offline: 适合需要长期访问用户数据的情况，应用程序可以在用户不在线时继续访问数据。
		ctx.Redirect(http.StatusFound,
			googleOauthConfig.AuthCodeURL(fmt.Sprintf(oauthStateString, action, activityID, device, gameID, userID), oauth2.AccessTypeOffline))
		return
	} else if device == "customize" {
		// 自定義功能沒有綁定

		if activityID == "" || userID == "" {
			h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(activity_id、user_id)")
			return
		}

		// 自定義簽到人員(活動主辦方匯入)，直接導向callback api
		ctx.Redirect(http.StatusFound,
			redirectURL+fmt.Sprintf("?action=%s&activity_id=%s&device=%s&user_id=%s&game_id=%s", action, activityID, device, userID, gameID))
		return
	}

	// 將state驗證參數寫入redis中
	// err := h.redisConn.SetAdd([]interface{}{config.LINE_STATES_REDIS, state})
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: state驗證發生問題")
	// 	return
	// }
}

// @@@Summary LOGIN CALLBACK
// @@@Tags Auth
// @@@version 1.0
// @@@Accept  json
// @@@param action query string false "報名、簽到、綁定、LINE選單" Enums(bind, apply, sign, search, create, case)
// @@@param activity_id query string false "活動ID"
// @@@param game_id query string false "遊戲ID"
// @@@param state query string false "actionis%sandactivity_idis%sanddeviceis%sandgame_idis%s"
// @@@param user_id query string false "用戶ID"
// @@@param device query string false "device" Enums(line, facebook, gmail, customize)
// @@@Success 302 {string} redirect
// @@@Failure 404 {array} response.ResponseBadRequest
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /auth/callback [get]
func (h *Handler) LoginCallback(ctx *gin.Context) {
	// LoginLock.Lock()
	// defer LoginLock.Unlock()

	var (
		host                   = ctx.Request.Host
		ip                     = utils.ClientIP(ctx.Request)
		activityID             = ctx.Query("activity_id")
		gameID                 = ctx.Query("game_id")
		state                  = ctx.Query("state")
		action                 = ctx.Query("action")
		device                 = ctx.Query("device")
		authURL                = fmt.Sprintf(config.HTTPS_AUTH_CALLBACK_URL, host)
		activityModel          models.ActivityModel
		lineModel              models.LineModel
		id                     int64
		isPushMessage, isFirst bool
		// activityID             string
		// action                 string
		// device                      string
		lineChannelID, lineChannelSecret string
		chatbotSecret, chatbotToken      string
		line                             string
		redirectURL                      string
		identify                         string
		status                           string
		userID                           string
		adminID                          string // 綁定管理員id
		name                             string
		avatar                           string
		email                            string
		err                              error
		// payload                     *social.Payload
		// activityID                                                                         = ctx.Query("activity_id")
		// action                                                                             = ctx.Query("action")
	)

	if ctx.Request.Host == config.API_URL && state == "" {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// line、fb、gmail有state參數
	// state參數處理，取得activity_ID.action.device.game_id
	// actionis%sandactivity_idis%sanddeviceis%sandgame_idis%sanduser_idis%s
	if len(state) > 8 {
		states := strings.Split(state[8:], "andactivity_idis") // [action資料, activity_id資料anddeviceis%sandgame_idis%sanduser_idis%s]

		// 取得action資料
		action = states[0]

		if len(states) > 1 {
			states = strings.Split(states[1], "anddeviceis") // [activity_id資料, device資料andgame_idis%sanduser_idis%s]

			// 取得activity_id
			activityID = states[0]
		}

		if len(states) > 1 {
			states = strings.Split(states[1], "andgame_idis") // [device資料, game_idis%sanduser_idis%s]

			// 取得device
			device = states[0]
		}

		if len(states) > 1 {
			states = strings.Split(states[1], "anduser_idis") // [game_id資料, user_id資料]

			// 取得game_id.user_id(這是平台管理員的user_id)
			gameID = states[0]
			adminID = states[1] // 平台管理員的user_id
		}
	}

	if action == "" || device == "" {
		h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(state)")
		return
	}

	if action == "bind" {
		// 平台綁定
		if adminID == "" {
			h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(user_id)")
			return
		}
	}

	// line、fb、gmail需要state驗證
	if device == "line" || device == "facebook" || device == "gmail" {
		// 防止被攻擊，state判斷
		if state != fmt.Sprintf(oauthStateString, action, activityID, device, gameID, userID) {
			h.executeErrorHTML(ctx, "錯誤: 驗證流程發生問題，請重新掃描QRcode進行驗證(state)")
			return
		}
	}

	// isMatch := checkState(h, ctx.Query("state"))
	// if isMatch == false {
	// 	h.executeErrorHTML(ctx, "錯誤: LINE驗證流程發生問題，請重新掃描QRcode進行驗證(state)")
	// 	return
	// }

	if activityID != "" {
		// 取得活動狀態資訊
		activityModel, err = models.DefaultActivityModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindActivityLeftJoinCustomize(activityID)
		if err != nil || activityModel.ID == 0 {
			h.executeErrorHTML(ctx, "錯誤: 活動已結束，謝謝參與")
			return
		}
	}

	// 不同裝置下，oauth取得用戶資訊
	if device == "line" {
		// line驗證

		if action == "bind" {
			// 平台綁定

			// 綁定功能(有user_id參數)
			// authURL += "?user_id=" + ctx.Query("user_id")

			userModel, err := models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(true, true, "users.user_id", adminID)
			if err != nil {
				h.executeErrorHTML(ctx, "錯誤: 取得用戶資料發生問題")
				return
			}

			// 判斷用戶是否綁定官方帳號
			if userModel.ChannelID != "" && userModel.ChannelSecret != "" &&
				userModel.ChatbotSecret != "" && userModel.ChatbotToken != "" &&
				userModel.LineID != "" {
				lineChannelID = userModel.ChannelID
				lineChannelSecret = userModel.ChannelSecret
				chatbotSecret = userModel.ChatbotSecret
				chatbotToken = userModel.ChatbotToken
				line = userModel.LineID
				isPushMessage = true
			}
		} else if action == "apply" || action == "sign" {
			// 報名簽到
			if activityID == "" {
				h.executeErrorHTML(ctx, "錯誤: 網頁請求發生問題(activity_id)")
				return
			}

			// 判斷是否綁定官方帳號
			if activityModel.ActivityChannelID != "" && activityModel.ActivityChannelSecret != "" &&
				activityModel.ActivityChatbotSecret != "" && activityModel.ActivityChatbotToken != "" &&
				activityModel.ActivityLineID != "" {
				// 活動官方帳號
				lineChannelID = activityModel.ActivityChannelID
				lineChannelSecret = activityModel.ActivityChannelSecret
				chatbotSecret = activityModel.ActivityChatbotSecret
				chatbotToken = activityModel.ActivityChatbotToken
				line = activityModel.ActivityLineID
				isPushMessage = true
			} else if activityModel.UserChannelID != "" && activityModel.UserChannelSecret != "" &&
				activityModel.UserChatbotSecret != "" && activityModel.UserChatbotToken != "" &&
				activityModel.UserLineID != "" {
				// 用戶官方帳號
				lineChannelID = activityModel.UserChannelID
				lineChannelSecret = activityModel.UserChannelSecret
				chatbotSecret = activityModel.UserChatbotSecret
				chatbotToken = activityModel.UserChatbotToken
				line = activityModel.UserLineID
				isPushMessage = true

			}
		}
		//  else if action == "search" || action == "create" || action == "case" {
		// 	// LINE選單功能
		// }

		// LINE參數為空，補上我方公司官方帳號參數
		if lineChannelID == "" || lineChannelSecret == "" || chatbotSecret == "" || chatbotToken == "" {
			lineChannelID = config.CHANNEL_ID
			lineChannelSecret = config.CHANNEL_SECRET
			chatbotSecret = config.CHATBOT_SECRET
			chatbotToken = config.CHATBOT_TOKEN
			if activityModel.PushMessage == "open" {
				isPushMessage = true
			}
		}

		var (
			// token *social.TokenResponse
			token = &social.TokenResponse{
				AccessToken: "",
				IDToken:     "",
			}
			isSetCookie bool
		)

		// 判斷cookie裡是否有line_session資料
		idToken, err1 := ctx.Request.Cookie("line_session_id_token")
		accessToken, err2 := ctx.Request.Cookie("line_session_access_token")
		if err1 != nil || err2 != nil {
			// log.Println("cookie裡無資料: ", err.Error())
			isSetCookie = true
		} else {
			// log.Println("cookie裡有資料")
			token.IDToken = idToken.Value
			token.AccessToken = accessToken.Value
		}

		client := NewSocialClient(ctx, lineChannelID, lineChannelSecret)

		if err1 != nil || err2 != nil {
			// token處理
			// 執行curl -X POST https://api.line.me/oauth2/v2.1/token 取得token
			if token, err = client.GetAccessToken(authURL, ctx.Query("code")).Do(); err != nil {
				// 官方帳號資訊錯誤會在這裡判斷
				h.executeErrorHTML(ctx, "錯誤: LINE驗證流程發生問題，請重新掃描QRcode進行驗證(GetAccessToken)")
				return
			}

			// 執行GET https://api.line.me/oauth2/v2.1/verify 驗證token
			if _, err = client.TokenVerify(token.AccessToken).Do(); err != nil {
				h.executeErrorHTML(ctx, "錯誤: LINE驗證流程發生問題，請重新掃描QRcode進行驗證(TokenVerify)")
				return
			}
			// 執行curl -v -X POST https://api.line.me/oauth2/v2.1/token 刷新token
			if _, err = client.RefreshToken(token.RefreshToken).Do(); err != nil {
				h.executeErrorHTML(ctx, "錯誤: LINE驗證流程發生問題，請重新掃描QRcode進行驗證(RefreshToken)")
				return
			}
		}

		// 取得用戶資料
		if len(token.IDToken) == 0 {
			var res *social.GetUserProfileResponse
			res, err = client.GetUserProfile(token.AccessToken).Do()

			userID = res.UserID
			name = res.DisplayName
			avatar = res.PictureURL
		} else {
			var payload *social.Payload
			payload, err = token.DecodePayload(config.CHANNEL_ID)

			userID = payload.Sub
			name = payload.Name
			avatar = payload.Picture
			email = payload.Email
		}
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: 取得用戶資料發生問題，請重新掃描QRcode進行驗證(GetUserProfile、DecodePayload)")
			return
		}

		// 將token資料設置至cookie中
		if isSetCookie {
			// 一個小時
			ctx.SetSameSite(4)
			ctx.SetCookie("line_session_id_token", string(token.IDToken),
				3600, "/", "", true, false)
			ctx.SetCookie("line_session_access_token", string(token.AccessToken),
				3600, "/", "", true, false)
		}
	} else if device == "facebook" {
		var (
			// token *oauth2.Token
			token = &oauth2.Token{
				AccessToken: "",
			}
			isSetCookie bool
		)

		// 判斷cookie裡是否有facebook_session資料
		accessToken, err := ctx.Request.Cookie("facebook_session_access_token")
		if err != nil {
			// log.Println("cookie裡無資料: ", err.Error())
			isSetCookie = true
		} else {
			// log.Println("cookie裡有資料")
			token.AccessToken = accessToken.Value
		}

		if err != nil {
			// facebook驗證
			// 第一次facebook api請求
			code := ctx.Request.FormValue("code")
			token, err = facebookOauthConfig.Exchange(context.Background(), code)
			if err != nil {
				h.executeErrorHTML(ctx, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Exchange)")
				return
			}
		}

		// 第二次facebook api請求
		response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email,picture.width(200).height(200)&access_token=" + token.AccessToken)
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Get)")
			return
		}
		defer response.Body.Close()

		var userData map[string]interface{}
		if err = json.NewDecoder(response.Body).Decode(&userData); err != nil {
			h.executeErrorHTML(ctx, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Decode)")
			return
		}

		userID = utils.GetString(userData["id"], "")
		name = utils.GetString(userData["name"], "")
		email = utils.GetString(userData["email"], "")
		avatar = config.HTTP_HILIVES_NET_URL + "/admin/uploads/facebook/" + userID + ".jfif"

		// 頭像資料處理(檔案)
		// ex: picture:map[data:map[ url:xxx ]]
		picture := userData["picture"].(map[string]interface{})
		data := picture["data"].(map[string]interface{})
		pictureURL := data["url"].(string)

		// 將圖片檔案儲存至遠端
		err = utils.DownloadFile(pictureURL, config.STORE_PATH+"/facebook/"+userID+".jfif")
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(圖片處理)")
			return
		}

		// 將token資料設置至cookie中
		if isSetCookie {
			// 一個小時
			ctx.SetSameSite(4)
			ctx.SetCookie("facebook_session_access_token", token.AccessToken,
				3600, "/", "", true, false)
		}
	} else if device == "gmail" {
		var (
			// token *oauth2.Token
			token = &oauth2.Token{
				AccessToken: "",
			}
			isSetCookie bool
		)

		// 判斷cookie裡是否有gmail_session資料
		accessToken, err := ctx.Request.Cookie("gmail_session_access_token")
		if err != nil {
			// log.Println("cookie裡無資料: ", err.Error())
			isSetCookie = true
		} else {
			// log.Println("cookie裡有資料")
			token.AccessToken = accessToken.Value
		}

		if err != nil {
			// gmail驗證
			code := ctx.Request.FormValue("code")

			// 取得token，為了取得用戶資料
			token, err = googleOauthConfig.Exchange(context.Background(), code)
			if err != nil {
				h.executeErrorHTML(ctx, "錯誤: gmail驗證流程發生問題，請重新掃描QRcode進行驗證(Exchange)")
				return
			}
		}

		// 利用token取得用戶資料
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: gmail驗證流程發生問題，請重新掃描QRcode進行驗證(Get)")
			return
		}
		defer response.Body.Close()

		var userData map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&userData); err != nil {
			h.executeErrorHTML(ctx, "錯誤: gmail驗證流程發生問題，請重新掃描QRcode進行驗證(Decode)")
			return
		}

		userID = utils.GetString(userData["id"], "")
		name = utils.GetString(userData["name"], "")
		email = utils.GetString(userData["email"], "")
		avatar = utils.GetString(userData["picture"], "")

		// 將token資料設置至cookie中
		if isSetCookie {
			// 一個小時
			ctx.SetSameSite(4)
			ctx.SetCookie("gmail_session_access_token", token.AccessToken,
				3600, "/", "", true, false)
		}
	} else if device == "customize" {
		// 自定義簽到人員(活動主辦方匯入)
		userID = ctx.Query("user_id")
	}

	for l := 0; l < MaxRetries; l++ {
		// 上鎖
		ok, _ := h.acquireLock(config.LINE_USERS_LOCK_REDIS+userID, LockExpiration)
		if ok == "OK" {
			// 釋放鎖
			// defer h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)

			// 加入funciton

			// 是否存在用戶資料(從資料庫取得資料)
			if lineModel, err = models.DefaultLineModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(false, ip, "user_id", userID); err != nil {
				h.executeErrorHTML(ctx, "錯誤: 取得用戶資料發生問題，請重新掃描QRcode進行驗證(Find)")

				// 釋放鎖
				h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
				return
			}
			id = lineModel.ID
			identify = lineModel.Identify // 字串長度32

			if device == "customize" {
				// 自定義簽到人員(活動主辦方匯入)，避免redis更新資料時為空值
				name = lineModel.Name
				avatar = lineModel.Avatar
				email = lineModel.Email
			}

			// 判斷是否存在用戶資料
			if lineModel.UserID == "" {
				// 不存在用戶資料

				if device == "line" {
					// line驗證
					// var link *linebot.LinkTokenResponse
					if _, err = NewLineBotClient(ctx, chatbotSecret, chatbotToken).IssueLinkToken(userID).Do(); err != nil {
						h.executeErrorHTML(ctx, "錯誤: LINE驗證流程發生問題，請重新掃描QRcode進行驗證(IssueLinkToken)")

						// 釋放鎖
						h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
						return
					}
					identify = userID // identify值
				} else if device == "facebook" {
					// facebook驗證，隨機產生identify
					identify = userID
				} else if device == "gmail" {
					// gmail驗證，隨機產生identify
					identify = userID
				} else if device == "customize" {
					// 自定義簽到人員(活動主辦方匯入)
					h.executeErrorHTML(ctx, "錯誤: 驗證流程發生問題，請重新掃描QRcode進行驗證(無該用戶資料)")

					// 釋放鎖
					h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
					return
				}

				if id, err = models.DefaultLineModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					Add(models.LineModel{
						UserID:   userID,
						Name:     name,
						Avatar:   avatar,
						Email:    email,
						Identify: identify,
						Friend:   "no",
						Line:     line,
						Device:   device,
					}); err != nil {
					h.executeErrorHTML(ctx, err.Error())

					// 釋放鎖
					h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
					return
				}

				// 用戶第一次執行LINE LOGIN
				isFirst = true
			} else {
				// log.Println("是否更新過姓名頭像? ", lineModel.IsModify)
				// log.Println("裝置資料有改變? ", lineModel.Name != name || lineModel.Email != email ||
				// 	lineModel.Avatar != avatar)

				// 存在用戶資料，更新用戶資料
				// 判斷資料是否更新
				// line、facebook、gmail裝置須更新
				if device == "line" || device == "facebook" || device == "gmail" {
					if lineModel.IsModify == "no" && (lineModel.Name != name || lineModel.Email != email ||
						lineModel.Avatar != avatar) {
						// log.Println("更新資料")
						if err = models.DefaultLineModel().
							SetConn(h.dbConn, h.redisConn, h.mongoConn).
							UpdateUser(models.LineModel{
								UserID: userID,
								Name:   name,
								Avatar: avatar,
								Email:  email,
							}); err != nil {
							h.executeErrorHTML(ctx, err.Error())

							// 釋放鎖
							h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
							return
						}
					} else {
						// log.Println("不更新資料")

						name = lineModel.Name
						email = lineModel.Email
						avatar = lineModel.Avatar
					}
				}
			}

			// 釋放鎖
			h.releaseLock(config.LINE_USERS_LOCK_REDIS + userID)
			break
		}

		// 鎖被佔用，稍微延遲後重試
		time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
	}

	// 將用戶資料加入redis中
	h.redisConn.DelCache(config.AUTH_USERS_REDIS + userID)
	values := []interface{}{config.AUTH_USERS_REDIS + userID}
	values = append(values, "id", id)
	values = append(values, "user_id", userID)
	values = append(values, "name", name)
	values = append(values, "phone", lineModel.Phone)
	values = append(values, "email", email)
	values = append(values, "ext_email", lineModel.ExtEmail)
	values = append(values, "avatar", avatar)
	values = append(values, "identify", identify)
	values = append(values, "friend", lineModel.Friend)
	values = append(values, "device", device)
	values = append(values, "ext_password", lineModel.ExtPassword)
	values = append(values, "is_modify", lineModel.IsModify)
	values = append(values, "admin_id", lineModel.AdminID)

	// values = append(values, "ip", ip)
	if err := h.redisConn.HashMultiSetCache(values); err != nil {
		h.executeErrorHTML(ctx, "錯誤: 設置用戶快取資料發生問題")
		return
	}
	// 設置過期時間
	// h.redisConn.SetExpire(config.AUTH_USERS_REDIS+userID,config.REDIS_EXPIRE)

	if action == "bind" {
		// 綁定功能(有user_id參數)

		// 綁定參數
		lineBind := "no"
		fbBind := "no"
		gmailBind := "no"

		var message string                                                                   // 官方帳號要傳遞的訊息
		redirectURL = fmt.Sprintf(config.HTTPS_ACTIVITY_URL, config.HILIVES_NET_URL, "true") // 導向後端活動平台

		// 將用戶admin_id欄位更新為平台管理員的user_id
		if err = models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateAdminID(userID, adminID); err != nil {
			message = "管理員與LINE用戶綁定過程中發生錯誤，可能原因: 該管理員已與其他LINE用戶綁定"
		} else {
			message = "管理員與LINE用戶綁定成功!"

			// 先判斷redis裡是否有用戶資訊，有則修改admin_id資料
			if value, _ := h.redisConn.HashGetCache(config.AUTH_USERS_REDIS+userID, "user_id"); value != "" {
				h.redisConn.HashMultiSetCache([]interface{}{config.AUTH_USERS_REDIS + userID,
					"admin_id", adminID})

				// 設置過期時間
				// h.redisConn.SetExpire(config.AUTH_USERS_REDIS+userID, config.REDIS_EXPIRE)
			}

			if device == "line" {
				lineBind = "yes"
			} else if device == "facebook" {
				fbBind = "yes"
			} else if device == "gmail" {
				gmailBind = "yes"
			}

			// 更新user表的綁定欄位資料(line_bind、fb_bind、gmail_bind)
			if err := models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Update(true, models.EditUserModel{UserID: adminID,
					LineBind:  lineBind,
					FbBind:    fbBind,
					GmailBind: gmailBind,
				},
					"users.user_id", adminID); err != nil {
				message = "管理員與LINE用戶綁定過程中發生錯誤，可能原因: 該管理員已與其他LINE用戶綁定"
			}
		}

		// 使用第三方的官方帳號才需傳送訊息
		if isPushMessage {
			if err = pushMessage(ctx, chatbotSecret, chatbotToken, userID, message); err != nil {
				h.executeErrorHTML(ctx, err.Error())
				return
			}
		}

		ctx.Redirect(http.StatusFound, redirectURL)
		return
	}

	if device == "line" {
		// line驗證裝置

		if action == "search" || action == "create" || action == "case" {
			// LINE選單功能，直接導向頁面
			redirectURL = fmt.Sprintf(config.HTTPS_LINE_RICHMENU_URL, config.HILIVES_NET_URL, action, userID)

			ctx.Redirect(http.StatusFound, redirectURL)
			return
		}
	}

	// 報名簽到人員資料
	applysign, err := models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(0, activityID, userID, false)
	if err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得報名簽到人員資料，請重新報名簽到")
		return
	}

	// 判斷報名審核開關
	if applysign.ID == 0 { // 無人員資料
		status = "no"

		// 新增人員資料(status: no，未補齊資料)
		id, err := models.DefaultApplysignModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Add(true,
				models.NewApplysignModel{
					UserID:     userID,
					ActivityID: activityID,
					Status:     status,
				})
		if err != nil {
			h.executeErrorHTML(ctx, "錯誤: 新增報名簽到人員發生問題，請重新報名簽到")
			return
		}

		applysign.ID = id
		applysign.Status = status
	}

	// 導向補齊資料頁面
	redirectURL = fmt.Sprintf(config.HTTPS_AUTH_USER_URL, config.HILIVES_NET_URL, activityID, userID, gameID)
	if action == "sign" { // 簽到審核開啟
		redirectURL += "&sign=open"
	}
	if isFirst { // 第一次執行LINE LOGIN
		redirectURL += "&isfirst=true"
	}

	// 自定義簽到人員
	if device == "customize" {
		// 自定義簽到人員(活動主辦方匯入)
		// 判斷姓名頭像欄位是否為空
		if lineModel.Name == "" || lineModel.Avatar == "" {
			// 導向補齊資料頁面
			ctx.Redirect(http.StatusFound, redirectURL)
			return
		}
	}

	// 判斷是否必填信箱或電話資料
	if (activityModel.ExtPhoneRequired == "true" && lineModel.Phone == "") ||
		(activityModel.PushPhoneMessage == "open" && lineModel.Phone == "") ||
		(activityModel.ExtEmailRequired == "true" && lineModel.ExtEmail == "") ||
		(activityModel.SendMail == "open" && lineModel.ExtEmail == "") {
		// 導向補齊資料頁面
		ctx.Redirect(http.StatusFound, redirectURL)
		return
	}

	// 判斷自定義欄位的必填資料是否補齊
	var exts = []string{activityModel.Ext1Required, activityModel.Ext2Required, activityModel.Ext3Required,
		activityModel.Ext4Required, activityModel.Ext5Required, activityModel.Ext6Required, activityModel.Ext7Required,
		activityModel.Ext8Required, activityModel.Ext9Required, activityModel.Ext10Required}
	var ExtValues = []string{applysign.Ext1, applysign.Ext2, applysign.Ext3,
		applysign.Ext4, applysign.Ext5, applysign.Ext6, applysign.Ext7,
		applysign.Ext8, applysign.Ext9, applysign.Ext10}
	for i, ext := range exts {
		if ext == "true" && ExtValues[i] == "" {
			// 導向補齊資料頁面
			ctx.Redirect(http.StatusFound, redirectURL)
			return
		}
	}

	// 資料補齊，跳過補齊資料頁面，導向報名簽到處理頁面
	// 補齊資料後將用戶狀態改為審核中或報名成功
	now, _ := time.ParseInLocation("2006-01-02 15:04:05",
		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	if applysign.Status == "no" {
		// 報名審核開啟
		if activityModel.ApplyCheck == "open" {
			status = "review"
		} else {
			// 報名審核關閉
			status = "apply"
		}

		if err = models.DefaultApplysignModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateStatus(true, config.HILIVES_NET_URL, models.EditApplysignModel{
				ActivityID: activityID,
				LineUsers:  []string{userID},
				ReviewTime: now.Format("2006-01-02 15:04:05"),
				Status:     status}, false); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			h.executeErrorHTML(ctx, err.Error())
			return
		}
	} else if applysign.Status == "review" && activityModel.ApplyCheck == "close" {
		// 原本資料為審核中但管理員已將報名審核關閉，因此更新為報名完成狀態
		status = "apply"
		if err = models.DefaultApplysignModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateStatus(true, config.HILIVES_NET_URL, models.EditApplysignModel{
				ActivityID: activityID,
				LineUsers:  []string{userID},
				ReviewTime: now.Format("2006-01-02 15:04:05"),
				Status:     status}, false); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			h.executeErrorHTML(ctx, err.Error())
			return
		}
	}

	// 資料補齊，直接導向報名簽到api
	redirectURL = fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, activityID, userID, gameID)
	// 簽到審核開啟
	if action == "sign" {
		redirectURL += "&sign=open"
	}
	if isFirst { // 第一次執行LINE LOGIN
		redirectURL += "&isfirst=true"
	}
	ctx.Redirect(http.StatusFound, redirectURL)
}

// 判斷是否補齊電話資料
// if lineModel.Phone == "" || phone == "" {
// if activityModel.ExtPhoneRequired == "true" && lineModel.Phone == "" {
// 	fmt.Println("電話資料為空")
// 	// 導向補齊資料頁面
// 	ctx.Redirect(http.StatusFound, redirectURL)
// 	return
// }

// db.Table("activity").WithConn(h.dbConn).
// 	LeftJoin(command.Join{
// 		FieldA:    "activity.activity_id",
// 		FieldB:    "activity_customize.activity_id",
// 		Table:     "activity_customize",
// 		Operation: "=",
// 	}).
// 	Where("end_time", ">", now).Where("activity.activity_id", "=", activityID).First()
// id, _ := activity["id"].(int64)
// applyCheck, _ := activity["apply_check"].(string)
// activityModel, err := models.DefaultActivityModel().SetDbConn(h.dbConn).
// 	FindActivityStatus(activityID)
// if err != nil || activityModel.ID == 0 {
// 	h.executeE
