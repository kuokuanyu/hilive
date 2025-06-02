package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/response"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// SignedRequestPayload 用來解析 signed_request 資料的結構
type SignedRequestPayload struct {
	UserID string `json:"user_id"`
	// 可以根據需求添加其他欄位，例如算法類型、過期時間等
}

// parseSignedRequest 解析並驗證 Facebook 的 signed_request
func parseSignedRequest(signedRequest, appSecret string) (*SignedRequestPayload, error) {
	parts := strings.Split(signedRequest, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("錯誤: signed_request格式無效")
	}

	// 解碼數據部分
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("錯誤: 解碼signed_request負載失敗: %v", err)
	}

	// 解析 JSON
	var data SignedRequestPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("錯誤: 解析JSON負載失敗: %v", err)
	}

	return &data, nil
}

// @Summary 刪除回呼api(當用戶要求刪除資料時執行)
// @Tags Auth
// @version 1.0
// @Accept  json
// @Success 200 {array} response.ResponseDeleteWithURL
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /facebook/delete [get]
func (h *Handler) DeleteHandler(ctx *gin.Context) {
	var (
		host          = ctx.Request.Host
		signedRequest = ctx.Query("signed_request")
		appSecret     = config.FACEBOOK_SECRET // 替換為你的 Facebook 應用程式的 App Secret
	)

	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if signedRequest == "" {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 缺少 signed_request 參數",
		})
		return
	}

	// 解析 signed_request
	payload, err := parseSignedRequest(signedRequest, appSecret)
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: signed_request 無效",
		})
		return
	}

	log.Println("用戶ID: ", payload.UserID)
	// 刪除line_users資料表用戶資料
	if err := db.Conn(h.dbConn).Table(config.LINE_USERS_TABLE).
		Where("user_id", "=", payload.UserID).Delete(); err != nil &&
		err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 刪除用戶資料發生問題",
		})
		return
	}

	// 刪除用戶redis資料
	h.redisConn.DelCache(config.AUTH_USERS_REDIS + payload.UserID) // 簽到人員資訊

	// 返回刪除狀態的 URL 和確認碼
	response.OkDeleteWithURL(ctx, payload.UserID, fmt.Sprintf("https://%s/delete/status", config.HILIVES_NET_URL))
}

// DeleteStatusHandler 顯示用戶已刪除資料訊息
func (h *Handler) DeleteStatusHandler(ctx *gin.Context) {
	h.executeHTML(ctx, "./hilives/hilive/views/fb_msg.html", executeParam{})
}

// *****限速處理*****
// 設定 Facebook OAuth 認證配置
// var facebookOauthConfig = &oauth2.Config{
// 	ClientID:     "YOUR_FACEBOOK_APP_ID",           // 替換為你的 Facebook App ID
// 	ClientSecret: "YOUR_FACEBOOK_APP_SECRET",       // 替換為你的 Facebook App Secret
// 	RedirectURL:  "http://localhost:8080/auth/facebook", // 替換為你的回調 URL
// 	Scopes:       []string{"email", "public_profile"},
// 	Endpoint:     facebook.Endpoint,
// }

// var (
// 	ctx = context.Background()
// 	// 初始化 Redis 客戶端
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379", // 替換為你的 Redis 伺服器地址
// 	})

// 	requestLimit = 200  // 每小時最大請求次數
// 	timeWindow   = time.Hour
// )

// // limitRequest 檢查當前請求是否超過限制
// func limitRequest(key string) (bool, error) {
// 	// 獲取當前請求計數
// 	currentCount, err := redisClient.Get(ctx, key).Int()
// 	if err != nil && err != redis.Nil {
// 		return false, err
// 	}

// 	// 檢查是否超過限速
// 	if currentCount >= requestLimit {
// 		return false, nil
// 	}

// 	// 增加請求計數
// 	if err := redisClient.Incr(ctx, key).Err(); err != nil {
// 		return false, err
// 	}

// 	// 設置過期時間，如果是第一次請求則設置
// 	if currentCount == 0 {
// 		if err := redisClient.Expire(ctx, key, timeWindow).Err(); err != nil {
// 			return false, err
// 		}
// 	}

// 	return true, nil
// }

// // facebookAuthHandler 處理 Facebook 認證請求
// func facebookAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	// 使用用戶 IP 地址作為 Redis 鍵
// 	key := "request_count:" + r.RemoteAddr

// 	// 限制請求
// 	allowed, err := limitRequest(key)
// 	if err != nil {
// 		http.Error(w, "錯誤: 限速檢查發生問題", http.StatusInternalServerError)
// 		return
// 	}
// 	if !allowed {
// 		http.Error(w, "錯誤: 超過請求限制，稍後再試。", http.StatusTooManyRequests)
// 		return
// 	}

// 	// 獲取認證碼
// 	code := r.FormValue("code")
// 	token, err := facebookOauthConfig.Exchange(ctx, code)
// 	if err != nil {
// 		http.Error(w, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Exchange)", http.StatusInternalServerError)
// 		return
// 	}

// 	// 請求用戶資料
// 	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email,picture.width(200).height(200)&access_token=" + token.AccessToken)
// 	if err != nil {
// 		http.Error(w, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Get)", http.StatusInternalServerError)
// 		return
// 	}
// 	defer response.Body.Close()

// 	var userData map[string]interface{}
// 	if err = json.NewDecoder(response.Body).Decode(&userData); err != nil {
// 		http.Error(w, "錯誤: Facebook驗證流程發生問題，請重新掃描QRcode進行驗證(Decode)", http.StatusInternalServerError)
// 		return
// 	}

// 	// 這裡可以處理 userData，根據需求返回用戶資料或執行其他操作
// 	fmt.Fprintf(w, "用戶數據: %+v\n", userData) // 示例：返回用戶數據
// }

// func main() {
// 	http.HandleFunc("/auth/facebook", facebookAuthHandler)
// 	http.ListenAndServe(":8080", nil)
// }

// *****限速處理*****

// var (
// 	facebookOauthConfig = &oauth2.Config{
// 		ClientID:     "814026164217592",
// 		ClientSecret: "e06412863dbb12daec68d54c0da47734",
// 		RedirectURL:  "https://dev.hilives.net/admin/login",
// 		Scopes:       []string{"email"}, // []string{"email", "public_profile","user_photos", "user_birthday",},
// 		Endpoint:     facebook.Endpoint,
// 	}
// 	// oauthStateString = "action=%s&activity_id=%s&device=%s" // 用於防止CSRF攻擊
// )

// func (h *Handler) HandleFacebookMain(ctx *gin.Context) {
// 	var html = `<html><body><a href="/facebook/login">Facebook 登入</a></body></html>`
// 	fmt.Fprint(ctx.Writer, html)
// }

// func (h *Handler) HandleFacebookLogin(ctx *gin.Context) {
// 	activityID := "12345"
// 	action := "sign"
// 	url := facebookOauthConfig.AuthCodeURL(fmt.Sprintf(oauthStateString, action, activityID, "facebook"))

// 	// url += "&action=xxx&activity_id=xxx&device=facebook"

// 	http.Redirect(ctx.Writer, ctx.Request, url, http.StatusTemporaryRedirect)
// }

// func (h *Handler) HandleFacebookCallback(ctx *gin.Context) {
// 	activityID := "12345"
// 	action := "sign"

// 	state := ctx.Request.FormValue("state")
// 	if state != fmt.Sprintf(oauthStateString, action, activityID, "facebook") {
// 		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
// 		http.Redirect(ctx.Writer, ctx.Request, "/facebook", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	log.Println("state: ", state)

// 	code := ctx.Request.FormValue("code")
// 	token, err := facebookOauthConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		log.Printf("code exchange failed: %s\n", err.Error())
// 		http.Redirect(ctx.Writer, ctx.Request, "/facebook", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email,picture.width(200).height(200)&access_token=" + token.AccessToken)
// 	if err != nil {
// 		log.Printf("failed getting user info: %s\n", err.Error())
// 		http.Redirect(ctx.Writer, ctx.Request, "/facebook", http.StatusTemporaryRedirect)
// 		return
// 	}
// 	defer response.Body.Close()

// 	var userData map[string]interface{}
// 	if err := json.NewDecoder(response.Body).Decode(&userData); err != nil {
// 		log.Printf("failed decoding user info: %s\n", err.Error())
// 		http.Redirect(ctx.Writer, ctx.Request, "/facebook", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	fmt.Fprintf(ctx.Writer, "User Info: %s", userData)
// }
