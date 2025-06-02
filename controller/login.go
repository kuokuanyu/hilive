package controller

import (
	"fmt"

	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("hilives") // 用于签名 JWT 的密钥

// Claims 结构体，用于存储 JWT 的声明
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成 JWT 令牌
func GenerateJWT(username string) (string, error) {
	// expirationTime := time.Now().Add(30 * time.Second)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// JWT 验证中间件
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")

// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
// 			c.Abort()
// 			return
// 		}

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// 将用户信息存储到上下文中
// 		c.Set("username", claims.Username)
// 		c.Next()
// 	}
// }

// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 	c.Abort()
// 	return
// }

// token := strings.TrimPrefix(authHeader, "Bearer ")
// if token != "example-valid-token" {  // 验证 Token
// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 	c.Abort()
// 	return
// }

// ShowLogin 登入 GET API
func (h *Handler) ShowLogin(ctx *gin.Context) {
	var host = ctx.Request.Host // 網域
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	// 測試
	// 生成 JWT
	// tokenString, err := GenerateJWT("testuser")
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 1")
	// 	return
	// }

	// // eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.Qovwfn20FmzSU8GG6BI-KXkYT_5a3jpTd7d2OouQHyU
	// log.Println("tokennn: ", tokenString)

	// // 驗證
	// claims := &Claims{}
	// token, _ := jwt.ParseWithClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.Qovwfn20FmzSU8GG6BI-KXkYT_5a3jpTd7d2OouQHyU",
	// 	claims, func(token *jwt.Token) (interface{}, error) {
	// 		return jwtKey, nil
	// 	})
	// log.Println("token.Raw: ", token.Raw)
	// log.Println("token.Signature: ", token.Signature)
	// log.Println("token.Valid: ", token.Valid)
	// log.Println("claims.Username: ", claims.Username)

	// 測試

	h.executeHTML(ctx, "./hilives/hilive/views/cms/login.html", executeParam{
		Route: route{
			Prefix:            h.config.Prefix,
			Login:             config.LOGIN_API_URL,
			Register:          config.REGISTER_API_URL,
			Retrieve:          config.RETRIEVE_API_URL,
			Verification:      config.VERIFICATION_API_URL,
			VerificationCheck: config.VERIFICATION_CHECK_API_URL,
			User:              config.USER_API_URL,
		},
	})
}

// Login 自定義簽到人員登入 POST API
// @Summary 自定義簽到人員登入
// @Tags Auth
// @version 1.0
// @param activity_id formData string true "activity_id"
// @param game_id formData string true "game_id"
// @param ext_password formData string false "密碼"
// @param action formData string true "apply.sign"
// @param is_login formData string true "is_login" Enums(true, false)
// @param is_create formData string true "is_create" Enums(true, false)
// @param name formData string false "name"
// @param phone formData string false "phone"
// @param ext_email formData string false "ext_email"
// @param ext_1 formData string false "ext_1"
// @param ext_2 formData string false "ext_2"
// @param ext_3 formData string false "ext_3"
// @param ext_4 formData string false "ext_4"
// @param ext_5 formData string false "ext_5"
// @param ext_6 formData string false "ext_6"
// @param ext_7 formData string false "ext_7"
// @param ext_8 formData string false "ext_8"
// @param ext_9 formData string false "ext_9"
// @param ext_10 formData string false "ext_10"
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/login [post]
func (h *Handler) ApplysignLogin(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Request.FormValue("activity_id")
		gameID     = ctx.Request.FormValue("game_id")
		isLogin    = ctx.Request.FormValue("is_login")  // 輸入驗證碼
		isCreate   = ctx.Request.FormValue("is_create") // 填寫表單
		password   = ctx.Request.FormValue("ext_password")
		action     = ctx.Request.FormValue("action")
		name       = ctx.Request.FormValue("name")
		phone      = ctx.Request.FormValue("phone")
		email      = ctx.Request.FormValue("ext_email")
		// values     = []string{ctx.Request.FormValue("ext_1"), ctx.Request.FormValue("ext_2"),
		// 	ctx.Request.FormValue("ext_3"), ctx.Request.FormValue("ext_4"), ctx.Request.FormValue("ext_5"),
		// 	ctx.Request.FormValue("ext_6"), ctx.Request.FormValue("ext_7"), ctx.Request.FormValue("ext_8"),
		// 	ctx.Request.FormValue("ext_9"), ctx.Request.FormValue("ext_10")}
		userID string
		err    error
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

	if activityID == "" || action == "" || (isLogin != "true" && isCreate != "true") {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法驗證參數，請輸入有效的資料",
		})
		return
	}

	if isLogin == "true" {
		// 輸入驗證碼
		applysignUserModel, ok := h.checkApplysignUser(activityID, password)
		if !ok {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法驗證用戶，請輸入有效的資料",
			})
			return
		}

		userID = applysignUserModel.UserID
	} else if isCreate == "true" {
		// 用戶填寫表單
		if name == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 姓名、電話不能為空，請輸入有效的資料",
			})
			return
		}

		// 將用戶加入資料表中、傳遞訊息
		if userID, err = models.DefaultApplysignUserModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Add(true,
				models.ApplysignUserModel{
					ActivityID: activityID,
					Name:       name,
					Phone:      phone,
					ExtEmail:   email,

					Ext1:  ctx.Request.FormValue("ext_1"),
					Ext2:  ctx.Request.FormValue("ext_2"),
					Ext3:  ctx.Request.FormValue("ext_3"),
					Ext4:  ctx.Request.FormValue("ext_4"),
					Ext5:  ctx.Request.FormValue("ext_5"),
					Ext6:  ctx.Request.FormValue("ext_6"),
					Ext7:  ctx.Request.FormValue("ext_7"),
					Ext8:  ctx.Request.FormValue("ext_8"),
					Ext9:  ctx.Request.FormValue("ext_9"),
					Ext10: ctx.Request.FormValue("ext_10"),

					// ExtPassword: utils.RandomNumber(6),

					Source: "user",
				}); err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	}

	response.OkWithURL(ctx, fmt.Sprintf(config.AUTH_REDIRECT_URL, action, "customize", activityID, userID, gameID))
}

// Login 登入 POST API
// @Summary 用戶登入
// @Tags Auth
// @version 1.0
// @param phone formData string true "手機號碼(10碼)" minlength(10) maxlength(10)
// @param password formData string true "密碼"
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var (
		host          = ctx.Request.Host
		password      = ctx.Request.FormValue("password")
		phone         = ctx.Request.FormValue("phone")
		userModel, ok = h.checkUser(password, phone)
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

	if !ok || password == "" || phone == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法驗證用戶、密碼與手機號碼不能為空，請輸入有效的手機號碼與密碼",
		})
		return
	}
	if len(userModel.Permissions) == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 用戶無瀏覽網頁權限，請聯絡管理員",
		})
		return
	}

	response.OkWithURL(ctx, fmt.Sprintf(config.HILIVE_SESSION_URL, userModel.UserID))
}

// Logout 登出
func (h *Handler) Logout(ctx *gin.Context) {
	user := h.GetLoginUser(ctx.Request, "hilive_session")

	// 清除用戶redis資訊
	h.redisConn.DelCache(config.HILIVE_USERS_REDIS + user.UserID)
	// h.redisConn.DelCache(config.LINE_USERS + user.UserID)

	// 清除cookie
	ctx.SetCookie("hilive_session", "", -1, "/", "", true, false)
	ctx.SetCookie("chatroom_session", "", -1, "/", "", true, false)

	ctx.Status(302)
	ctx.Header("Location", config.LOGIN_URL)
}

// if len(phone) > 2 {
// 	if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
// 		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  userID,
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
// 		})
// 		return
// 	}
// } else {
// 	response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 		UserID:  userID,
// 		Method:  ctx.Request.Method,
// 		Path:    ctx.Request.URL.Path,
// 		Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
// 	})
// 	return
// }

// if email != "" && !strings.Contains(email, "@") {
// 	response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 		UserID:  userID,
// 		Method:  ctx.Request.Method,
// 		Path:    ctx.Request.URL.Path,
// 		Message: "錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址",
// 	})
// 	return
// }

// *****可以同時登入(暫時拿除)*****
// 清除用戶的ip資料
// if err := models.DefaultUserModel().SetDbConn(h.dbConn).
// 	UpdateIP(user.UserID, ""); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }
// *****可以同時登入(暫時拿除)*****

// if errCookie := auth.DeleteCookie(ctx,
// 	db.GetConnectionFromService(h.services.Get(h.config.Databases["hilive"].Driver)),
// 	"hilive_session"); errCookie != nil &&
// 	errCookie.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	log.Println("登出發生錯誤: ", errCookie)
// }

// *****可以同時登入(確定拿除)*****
// 更新用戶ip資訊
// if err := models.DefaultUserModel().SetDbConn(h.dbConn).
// 	UpdateIP(userModel.UserID, ip); err != nil {
// 	response.Error(ctx, "錯誤: 更新用戶資料發生問題(ip資訊)")
// 	return
// }
// *****可以同時登入(確定拿除)*****

// *****可以同時登入(確定拿除)*****
// 將用戶資訊加入redis中
// values := []interface{}{config.HILIVE_USERS_REDIS + userModel.UserID}
// values = append(values, "id", userModel.ID)
// values = append(values, "user_id", userModel.UserID)
// values = append(values, "name", userModel.Name)
// values = append(values, "phone", userModel.Phone)
// values = append(values, "email", userModel.Email)
// values = append(values, "avatar", userModel.Avatar)
// values = append(values, "bind", userModel.Bind)
// values = append(values, "cookie", userModel.Cookie)
// values = append(values, "ip", ip)
// values = append(values, "max_activity", userModel.MaxActivity)
// values = append(values, "max_activity_people", userModel.MaxActivityPeople)
// values = append(values, "max_game_people", userModel.MaxGamePeople)
// // LINE資訊
// values = append(values, "line_id", userModel.LineID)
// values = append(values, "channel_id", userModel.ChannelID)
// values = append(values, "channel_secret", userModel.ChannelSecret)
// values = append(values, "chatbot_secret", userModel.ChatbotSecret)
// values = append(values, "chatbot_token", userModel.ChatbotToken)

// // 用戶權限
// // json編碼
// permissions := utils.JSON(userModel.Permissions)
// values = append(values, "permissions", permissions)

// // 活動權限
// // json編碼
// activitys := utils.JSON(userModel.Activitys)
// values = append(values, "activitys", activitys)

// // 菜單
// // json編碼
// menus := utils.JSON(userModel.ActivityMenus)
// values = append(values, "activity_menus", menus)

// // params = append(params, "table", "users")
// if err := h.redisConn.HashMultiSetCache(values); err != nil {
// 	response.Error(ctx, "錯誤: 設置用戶快取資料發生問題")
// 	return
// }
// // 設置過期時間
// h.redisConn.SetExpire(config.HILIVE_USERS_REDIS+userModel.UserID,
// 	strconv.Itoa(config.GetSessionLifeTime()))
// *****可以同時登入(確定拿�
