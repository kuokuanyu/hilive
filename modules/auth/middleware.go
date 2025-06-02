package auth

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"hilive/models"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// MiddlewareCallback is type of callback function.
type MiddlewareCallback func(ctx *gin.Context)

// Invoker 中間件驗證
type Invoker struct {
	// prefix                 string
	sessionName      string
	authFailCallback MiddlewareCallback
	// permissionDenyCallback MiddlewareCallback
	dbConn    db.Connection
	redisConn cache.Connection
}

// DefaultHiliveInvoker 預設Invoker
func DefaultHiliveInvoker(dbConn db.Connection, redisConn cache.Connection, sessionName string) *Invoker {
	return &Invoker{
		// prefix:      config.Prefix(),
		sessionName: sessionName,
		authFailCallback: func(ctx *gin.Context) {
			// fmt.Println("跳出?")

			// if ctx.Request.URL.Path == config.LOGIN_URL {
			// 	return
			// } else if ctx.Request.URL.Path == config.LOGOUT_URL {
			// 	ctx.Redirect(http.StatusFound, config.LOGIN_URL)
			// }

			tmpl, err := template.New("").Parse(`<script>
				alert("登入逾時或無權限瀏覽，請重新登入")
				location.href = "` + config.LOGIN_URL + `"
				</script>`)
			if err != nil {
				log.Println("解析模板發生錯誤")
			}
			if err = tmpl.Execute(ctx.Writer, nil); err != nil {
				log.Println("登入逾時模板時發生錯誤")
			}

			// fmt.Println("跳出結束?")
		},
		// permissionDenyCallback: func(ctx *gin.Context) {
		// },
		dbConn:    dbConn,
		redisConn: redisConn,
	}
}

// DefaultChatroomInvoker 預設Invoker
func DefaultChatroomInvoker(dbConn db.Connection, redisConn cache.Connection, sessionName string) *Invoker {
	return &Invoker{
		// prefix:      config.Prefix(),
		sessionName: sessionName,
		authFailCallback: func(ctx *gin.Context) {
			var (
				params                       = make([]string, 0)
				activityID, gameID, redirect string
			)

			// 當chatroom_session無效時，判斷是否有hilive_session
			// if cookie, err := ctx.Request.Cookie("hilive_session"); err == nil && cookie.Value != "" {
			// 	tmpl, err := template.New("").Parse(`<script>
			// alert("登入逾時或無權限瀏覽，請重新登入")
			// location.href = "` + config.LOGIN_URL + `"
			// </script>`)
			// 	if err != nil {
			// 		log.Println("解析模板發生錯誤")
			// 	}
			// 	if err = tmpl.Execute(ctx.Writer, nil); err != nil {
			// 		log.Println("登入逾時模板時發生錯誤")
			// 	}
			// }

			// 處理activity_id、game_id
			if ctx.Query("activity_id") != "" {
				// log.Println("一般URL")
				// log.Println("ctx.Request.URL: ", ctx.Request.URL)
				// log.Println("ctx.Request.URL.Path: ", ctx.Request.URL.Path)

				// 一般URL
				activityID = ctx.Query("activity_id")
				gameID = ctx.Query("game_id")
				// channelID = ctx.Query("channel_id")

				// 判斷是否為簽名牆或提問牆
				if strings.Contains(ctx.Request.URL.Path, "/signname/signwall") {
					// log.Println("簽名牆")
					gameID = "signname"
				} else if strings.Contains(ctx.Request.URL.Path, "/guest/question") {
					// log.Println("提問牆")
					// gameID = "question"
				}
			} else if ctx.Query("liff.state") != "" {
				// log.Println("middleware function，liff url參數處理")
				// log.Println("ctx.Request.URL.Path: ", ctx.Request.URL.Path)

				// line裝置，liff url會顯示此參數
				liffState := ctx.Query("liff.state") // ex: liff.state=?activity_id=xxx(&role=xxx)&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)
				if len(liffState) > 13 {
					liffState = liffState[13:] // 不讀取?activity_id字串
				}

				// liff.state參數處理
				if strings.Contains(liffState, "&game_id=") {
					// log.Println("遊戲頁面")
					// 遊戲頁面
					// ex: 活動ID(&role=xxx)&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)

					if strings.Contains(liffState, "&role=") {
						// liff url目前沒有設置有channel_id參數的api
						// liff url目前沒有設置有role參數的api
						// 包含role參數(扭蛋機、投票...等遊戲，但liff url裡沒有使用到role參數)
						// 活動ID&role=xxx&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)
						// params = strings.Split(liffState, "&role=") // [活動ID、角色&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)]
						// activityID = params[0]

						// // 拆解params[1]，角色&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)
						// params = strings.Split(params[1], "&game_id=") // [角色、遊戲ID(&channel_id=xxx)(#mst_challenge=xxx)]
						// gameID = params[1]                             // 遊戲ID(&channel_id=xxx)(#mst_challenge=xxx)

						// // 是否包含channel_id參數
						// if strings.Contains(liffState, "&channel_id=") {
						// 	// // 包含channel_id參數
						// 	// // 遊戲ID&channel_id=xxx(#mst_challenge=xxx)
						// 	// params = strings.Split(gameID, "&channel_id=") // [遊戲ID、頻道ID(#mst_challenge=xxx)]

						// 	// // 遊戲ID參數
						// 	// gameID = params[0]

						// 	// // 取得channel_id(可能含有#mst_challenge=xxx)
						// 	// channelID = params[1] // 頻道ID(#mst_challenge=xxx)
						// 	// // 判斷字串#mst_challenge
						// 	// // 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&channel_id=xxx#mst_challenge=xxx
						// 	// if strings.Contains(channelID, "#mst_challenge") {
						// 	// 	channelID = strings.Split(channelID, "#mst_challenge")[0]
						// 	// }
						// } else {
						// 	// 不包含channel_id參數
						// 	// 遊戲ID(#mst_challenge=xxx)
						// 	if strings.Contains(gameID, "#mst_challenge") {
						// 		gameID = strings.Split(gameID, "#mst_challenge")[0]
						// 	}
						// }

					} else {
						// 不包含role參數
						// 活動ID&game_id=xxx(&channel_id=xxx)(#mst_challenge=xxx)

						params = strings.Split(liffState, "&game_id=") // [活動ID、遊戲ID(&channel_id=xxx)(#mst_challenge=xxx)]
						activityID = params[0]

						// 是否包含channel_id參數
						if strings.Contains(liffState, "&channel_id=") {
							// liff url目前沒有設置有channel_id參數的api
							// // log.Println("不包含role參數，包含channel_id參數")
							// // 包含channel_id參數
							// // 遊戲ID&channel_id=xxx(#mst_challenge=xxx)

							// // 取得game_id
							// gameID = params[1]                             // 遊戲ID&channel_id=xxx(#mst_challenge=xxx)
							// params = strings.Split(gameID, "&channel_id=") // ex: [遊戲ID、頻道ID(#mst_challenge=xxx)]
							// gameID = params[0]

							// // 取得channel_id(可能含有#mst_challenge=xxx)
							// channelID = params[1] // 頻道ID(#mst_challenge=xxx)
							// // 判斷字串#mst_challenge
							// // 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&channel_id=xxx#mst_challenge=xxx
							// if strings.Contains(channelID, "#mst_challenge") {
							// 	channelID = strings.Split(channelID, "#mst_challenge")[0]
							// }
						} else {
							// log.Println("不包含role參數，不包含channel_id參數")
							// 不包含channel_id參數
							// 遊戲ID(#mst_challenge=xxx)
							gameID = params[1]

							if strings.Contains(gameID, "#mst_challenge") {
								gameID = strings.Split(params[1], "#mst_challenge")[0]
							}
						}
					}

				} else {
					// log.Println("報名簽到頁面")
					// 報名簽到頁面
					// ex: 活動ID(&sign=xxx)(&channel_id=xxx)(#mst_challenge=xxx)

					if strings.Contains(liffState, "&channel_id=") && strings.Contains(liffState, "&sign=") {
						// liff url目前沒有設置有channel_id參數的api
						// // liff url目前沒有設置有sign參數的api
						// // 包含channel_id、sign參數
						// // 活動ID&sign=xxx&channel_id=xxx(#mst_challenge=xxx)
						// params = strings.Split(liffState, "&sign=") // ex: [活動ID、sign參數&channel_id=xxx(#mst_challenge=xxx)]
						// activityID = params[0]

						// // 取得channel_id(sign參數&channel_id=xxx(#mst_challenge=xxx))
						// params = strings.Split(params[1], "&channel_id=") // ex: [sign參數、頻道ID(#mst_challenge=xxx)]

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
						// liff url目前沒有設置有sign參數的api
						// 包含sign參數
						// 活動ID&sign=xxx(#mst_challenge=xxx)

						// params = strings.Split(liffState, "&sign=") // ex: [活動ID，sign參數(#mst_challenge=xxx)]
						// activityID = params[0]
					} else {
						// log.Println("沒有channel_id參數，沒有sign參數")
						// 兩個參數都不包含
						// 活動ID(#mst_challenge=xxx)
						activityID = liffState

						if strings.Contains(activityID, "#mst_challenge") {
							activityID = strings.Split(activityID, "#mst_challenge")[0]
						}
					}

					// 判斷是否為簽名牆或提問牆
					if strings.Contains(ctx.Request.URL.Path, "/signname/signwall") {
						// log.Println("簽名牆")
						gameID = "signname"
					} else if strings.Contains(ctx.Request.URL.Path, "/guest/question") {
						// log.Println("提問牆")
						// gameID = "question"
					}
				}
			}

			// 如果沒有用戶session資訊則導向驗證頁面
			// activityModel, _ := models.DefaultActivityModel().SetDbConn(dbConn).
			// 	Find(false, activityID)

			// if activityModel.Device == "line" {
			// 	// liff url
			// 	redirect = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, activityID)
			// } else {
			// 一般url
			redirect = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, config.HILIVES_NET_URL, activityID, gameID) + "&is_exist=false"
			// }

			// 簽到審核開關
			redirect += "&sign=open"

			// log.Println("導向: ", redirect)

			ctx.Redirect(http.StatusFound, redirect)
		},
		// permissionDenyCallback: func(ctx *gin.Context) {
		// },
		dbConn:    dbConn,
		redisConn: redisConn,
	}
}

// Middleware 驗證及判斷用戶是否擁有權限登入
func (invoker *Invoker) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authOk := Filter(ctx, invoker.dbConn, invoker.redisConn, invoker.sessionName)
		if authOk {
			// lock.Lock()
			// defer lock.Unlock()
			// userValue["user"] = user

			ctx.Next()
			return
		} else if !authOk {
			invoker.authFailCallback(ctx)
			ctx.Abort()
			return
		}
		// if !permissionOk {
		// 	invoker.permissionDenyCallback(ctx)
		// 	ctx.Abort()
		// 	return
		// }
	}
}

// Filter 驗證用戶角色權限資訊
func Filter(ctx *gin.Context, dbConn db.Connection, redisConn cache.Connection, sessionName string) bool {
	var (
		ip                                      = utils.ClientIP(ctx.Request)
		path                                    = ctx.Request.URL.Path
		user                                    models.UserModel
		cookieParams                            url.Values
		activityID, userID, table, sign, secret string
		ok                                      bool
		// session *Session
		// userID  string
		// err     error
	)

	// url帶有login=false參數，不需要登入也不需要cookie value(不能使用管理平台)
	// 手機.終端機裝置也不需要登入不需要cookie value
	if !strings.Contains(path, "admin") &&
		(ctx.Query("login") == "false" ||
			ctx.Query("device") == "mobile" || ctx.Query("device") == "terminal") {
		return true
	}

	// 判斷cookie資訊
	if cookie, err := ctx.Request.Cookie(sessionName); err == nil && cookie.Value != "" {
		// 解碼cookie值
		decode, err := utils.Decode([]byte(cookie.Value))
		if err != nil {
			models.DefaultErrorLogModel().
				SetConn(dbConn, redisConn, nil).
				Add(models.EditErrorLogModel{
					UserID:    utils.ClientIP(ctx.Request),
					Code:      http.StatusBadRequest,
					Method:    ctx.Request.Method,
					Path:      ctx.Request.URL.Path,
					Message:   fmt.Sprintf("錯誤: 無法取得cookie值，無法辨識用戶(有cookie但無法解碼) , %s", err.Error()),
					PathQuery: ctx.Request.URL.RawQuery,
				})

			// params, _ := url.ParseQuery(string(decode))
			// fmt.Println("錯誤發生: ", cookie.Value)
			return false
		}

		cookieParams, err = url.ParseQuery(string(decode))
		if err != nil {
			models.DefaultErrorLogModel().
				SetConn(dbConn, redisConn, nil).
				Add(models.EditErrorLogModel{
					UserID:    utils.ClientIP(ctx.Request),
					Code:      http.StatusBadRequest,
					Method:    ctx.Request.Method,
					Path:      ctx.Request.URL.Path,
					Message:   "錯誤: 無法解析cookie值，無法辨識用戶",
					PathQuery: ctx.Request.URL.RawQuery,
				})

			return false
		}

		// activityID = cookieParams.Get("activity_id") // 活動ID，判斷用戶參加的活動場次
		userID = cookieParams.Get("user_id") // 用戶ID
		table = cookieParams.Get("table")
		sign = cookieParams.Get("sign")

		// fmt.Println("userID: ", userID)
		// fmt.Println("table: ", table)
	} else {
		// logger.LoggerError(ctx, "錯誤: 沒有cookie值")
		// url帶有login=false參數，則不需要登入也不需要cookie value(不能使用管理平台)
		// if !strings.Contains(path, "admin") && ctx.Query("login") == "false" {
		// 	return true
		// }

		// #####壓力測試暫時使用
		// 解碼cookie值
		// decode, err := utils.Decode([]byte("YXZhdGFyPWh0dHBzJTNBJTJGJTJGcHJvZmlsZS5saW5lLXNjZG4ubmV0JTJGMGhpVlg5X2hla05ubGxRQjMxOUdkSkxsa0ZPQlFTYmpBeEhYWXVHUk5HYmtwTEpIRW9XQ1Y2SFVsRGF4d2RJM1VzRFNGd0dCUkNZQm9iJmVtYWlsPWExNjc4Mjk0MzUlNDB5YWhvby5jb20udHcmZnJpZW5kPXllcyZpZGVudGlmeT1mOHRobmZqSmJScTBKdjNHd0hTQVZlTndJZUoyd0NkSCZpcD0xMjIuMTE3LjM4LjE3OCZuYW1lPSVFNSU4NiVBMHlvJnBob25lPTA5MzI1MzA4MTMmc2lnbj02NGJkYThiOTJkYjQ1OGVmM2FjYzM3OWU5MGYyYzNmNSZ0YWJsZT1saW5lX3VzZXJzJnVzZXJfaWQ9VTMxNDBmZDczY2ZkMzVkZDk5MjY2OGFiM2I2ZWZkYWU5"))
		// if err != nil {
		// 	params, _ := url.ParseQuery(string(decode))
		// 	fmt.Println(params)
		// 	return false, false
		// }
		// params, err := url.ParseQuery(string(decode))
		// if err != nil {
		// 	return false, false
		// }
		// // 用戶資訊
		// ip = "122.117.38.178"
		// userID = params.Get("user_id")
		// table = params.Get("table")
		// sign = params.Get("sign")

		// 網頁端目前沒有cookie value
		return false
	}

	// 聊天室session
	if sessionName == "chatroom_session" {
		secret = config.ChatroomCookieSecret

		params := make([]string, 0)
		// 處理url裡的activity_id
		if ctx.Query("activity_id") != "" {
			activityID = ctx.Query("activity_id")
		} else if ctx.Query("liff.state") != "" {
			liffState := ctx.Query("liff.state")
			if len(liffState) > 13 {
				liffState = liffState[13:] // 不讀取?activity_id字串
			}

			// liff.state參數處理
			if strings.Contains(liffState, "&game_id=") {
				// ex: liff.state=?activity_id=xxx&game_id=xxx#mst_challenge=xxx
				params = strings.Split(liffState, "&game_id=")
			} else {
				// ex: liff.state=?activity_id=xxx(&sign=xxx)#mst_challenge=xxx
				params = strings.Split(liffState, "&sign=")
			}

			activityID = params[0]
			if len(params) == 1 {
				// 安卓系統會有這個參數，ex: liff.state=?activity_id=xxx&...#mst_challenge=xxx
				params = strings.Split(params[0], "#mst_challenge")
				activityID = params[0]
			}
		}

		// 判斷用戶session裡的活動場次ID與URL裡的活動ID是否相符
		if cookieParams.Get("activity_id") != activityID {
			// fmt.Println("用戶不是參加該場次活動，導向LINE驗證流程", cookieParams.Get("activity_id"), activityID)
			// logger.LoggerError(ctx, "錯誤: 用戶不是參加該場次活動，導向LINE驗證流程")
			return false
		}

		// 判斷ip是否與資料表裡的ip資訊相同
		if table == "line_users" {
			// LINE用戶
			if user, ok = GetLineUser(ip, userID, dbConn, redisConn); !ok {

				models.DefaultErrorLogModel().
					SetConn(dbConn, redisConn, nil).
					Add(models.EditErrorLogModel{
						UserID:    userID,
						Code:      http.StatusBadRequest,
						Method:    ctx.Request.Method,
						Path:      ctx.Request.URL.Path,
						Message:   "錯誤: 查詢LINE用戶資訊發生問題",
						PathQuery: ctx.Request.URL.RawQuery,
					})
				return false
			}

			// fmt.Println("user.Avatar: ", user.Avatar)

			// *****可以同時登入(暫時拿除)*****
			// 判斷是否被其他人登入
			// if user.Ip != ip {
			// 	return false
			// }
			// *****可以同時登入(暫時拿除)*****
		} else if table == "users" {
			// 管理員
			if user, ok = GetUser(ip, userID, dbConn, redisConn); !ok {
				models.DefaultErrorLogModel().
					SetConn(dbConn, redisConn, nil).
					Add(models.EditErrorLogModel{
						UserID:    userID,
						Code:      http.StatusBadRequest,
						Method:    ctx.Request.Method,
						Path:      ctx.Request.URL.Path,
						Message:   "錯誤: 查詢平台管理員用戶資訊發生問題(chatroom_session)",
						PathQuery: ctx.Request.URL.RawQuery,
					})
				return false
			}
		}
	} else if sessionName == "hilive_session" {
		secret = config.HiliveCookieSecret

		// 判斷ip是否與資料表裡的ip資訊相同
		if user, ok = GetUser(ip, userID, dbConn, redisConn); !ok {
			models.DefaultErrorLogModel().
				SetConn(dbConn, redisConn, nil).
				Add(models.EditErrorLogModel{
					UserID:    userID,
					Code:      http.StatusBadRequest,
					Method:    ctx.Request.Method,
					Path:      ctx.Request.URL.Path,
					Message:   "錯誤: 查詢平台管理員用戶資訊發生問題(hilive_session)",
					PathQuery: ctx.Request.URL.RawQuery,
				})

			return false
		}
	}

	// ip.簽名不相符，代表帳號被其他人登入
	// *****可以同時登入(暫時拿除)*****
	// if user.Ip != ip || sign != utils.UserSign(user.UserID, secret) {
	// fmt.Println("ip問題: ", user.Ip != ip)
	// fmt.Println("簽名問題: ", sign != utils.UserSign(ip, user.UserID, secret))
	// fmt.Println(user.Ip)
	// fmt.Println(ip)
	// fmt.Println(sign)
	// fmt.Println(utils.UserSign(ip, user.UserID, secret))
	// return false
	// }
	// *****可以同時登入(暫時拿除)*****

	// *****可以同時登入(新)*****
	// 判斷簽名是否相符
	if sign != utils.UserSign(user.UserID, secret) {
		models.DefaultErrorLogModel().
			SetConn(dbConn, redisConn, nil).
			Add(models.EditErrorLogModel{
				UserID:    userID,
				Code:      http.StatusBadRequest,
				Method:    ctx.Request.Method,
				Path:      ctx.Request.URL.Path,
				Message:   "錯誤: 簽名資料不相符",
				PathQuery: ctx.Request.URL.RawQuery,
			})

		return false
	}
	// *****可以同時登入(新)*****

	// *****舊*****
	// // ip.簽名不相符，代表帳號被其他人登入
	// if user.Ip != ip || sign != utils.UserSign(ip, user.UserID, secret) {
	// 	// fmt.Println("ip問題: ", user.Ip != ip)
	// 	// fmt.Println("簽名問題: ", sign != utils.UserSign(ip, user.UserID, secret))
	// 	// fmt.Println(user.Ip)
	// 	// fmt.Println(ip)
	// 	// fmt.Println(sign)
	// 	// fmt.Println(utils.UserSign(ip, user.UserID, secret))
	// 	return false
	// }
	// *****舊*****

	// 檢查用戶權限
	if sessionName == "hilive_session" {
		if path != "/admin/user" && // 用戶頁面不用判斷
			!strings.Contains(path, "/admin/activity") { // 活動頁面不用判斷
			// !strings.Contains(path, "/prize?") && // 獎品頁面不用判斷(目前跟遊戲頁面在一起)
			// !strings.Contains(path, "/attend?") && // 參加人員頁面不用判斷
			// !strings.Contains(path, "/winning?") && // 中獎人員頁面不用判斷
			// !strings.Contains(path, "/black?") && // 黑名單人員頁面不用判斷
			// !strings.Contains(path, "/new?") && // 新增頁面不用判斷
			// !strings.Contains(path, "/edit?") { // 編輯頁面不用判斷

			// 除了管理員頁面，其他頁面的activity_id參數不能為空
			if (!strings.Contains(path, "/admin/manager") &&
				!strings.Contains(path, "/admin/permission") &&
				!strings.Contains(path, "/admin/menu") &&
				!strings.Contains(path, "/admin/overview") &&
				!strings.Contains(path, "/admin/log") &&
				!strings.Contains(path, "/admin/error_log")) &&
				ctx.Query("activity_id") == "" {
				models.DefaultErrorLogModel().
					SetConn(dbConn, redisConn, nil).
					Add(models.EditErrorLogModel{
						UserID:  userID,
						Code:    http.StatusBadRequest,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: activity_id參數為空(除了管理員頁面，其他頁面的activity_id參數不能為空)",
					})

				return false
			}
			return user.CheckPermission(ctx.Query("activity_id"), path)
		}
	}

	return true
	// return user, true,
	// 	user.CheckPermission(ctx.Request.URL.String(), ctx.Request.Method, ctx.Request.PostForm)
}

// GetUser 取得用戶資訊、是否有訪問權限
func GetUser(ip, userID string, dbConn db.Connection, redisConn cache.Connection) (user models.UserModel, ok bool) {
	var err error

	// *****可以同時登入(新)*****
	if user, err = models.DefaultUserModel().
		SetConn(dbConn, redisConn, nil).
		Find(true, true, "users.user_id", userID); err != nil ||
		user.UserID == "" {
		ok = false
		return
	}
	// *****可以同時登入(新)*****

	// *****可以同時登入(暫時拿除)*****
	// if user, err = models.DefaultUserModel().SetDbConn(dbConn).
	// 	SetRedisConn(redisConn).Find(true, true, ip, "users.user_id", userID); err != nil ||
	// 	user.UserID == "" {
	// 	ok = false
	// 	return
	// }
	// *****可以同時登入(暫時拿除)*****

	// TODO: 目前沒有權限相關問題，先拿掉角色權限菜單判斷
	// user = user.GetRoles().GetPermissions().GetMenus()
	// if len(user.MenuIDs) != 0 || user.IsSuperAdmin() || user.Roles[0].Name == "guest" {
	// 	ok = true
	// }
	return user, true
}

// GetLineUser 取得line用戶資訊
func GetLineUser(ip, userID string, dbConn db.Connection, redisConn cache.Connection) (user models.UserModel, ok bool) {
	lineUser, err := models.DefaultLineModel().
		SetConn(dbConn, redisConn, nil).
		Find(true, ip, "user_id", userID)
	if err != nil || lineUser.UserID == "" {
		return user, false
	}

	user = models.UserModel{
		ID:     lineUser.ID,
		UserID: lineUser.UserID,
		Name:   lineUser.Name,
		Avatar: lineUser.Avatar,
		// *****可以同時登入(暫時拿除)*****
		// Ip:     lineUser.Ip,
		// *****可以同時登入(暫時拿除)*****
	}
	return user, true
}

// if user, ok = GetLineUser(userID, conn); !ok {
// 	if user, ok = GetUser(userID, conn); !ok {
// 		return false, false
// 	}
// }

// 先判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
// value, _ := redisConn.HashGetCache(config.HILIVE_USERS+userID, "ip")
// if value != "" {
// 	// redis已有用戶資訊
// 	user = models.UserModel{
// 		UserID: userID,
// 		Ip:     value,
// 	}
// } else if value == "" {
// }

// 將用戶資訊加入redis中
// values := []interface{}{config.HILIVE_USERS + user.UserID}
// values = append(values, "user_id", user.UserID)
// values = append(values, "name", user.Name)
// values = append(values, "phone", user.Phone)
// values = append(values, "email", user.Email)
// values = append(values, "avatar", user.Avatar)
// values = append(values, "bind", user.Bind)
// values = append(values, "cookie", user.Cookie)
// values = append(values, "ip", user.Ip)
// params = append(params, "table", "users")
// if err := redisConn.HashMultiSetCache(values); err != nil {
// 	ok = false
// 	return
// }
// 設置過期時間
// if err := redisConn.SetExpire(config.HILIVE_USERS+user.UserID,
// 	strconv.Itoa(config.GetSessionLifeTime())); err != nil {
// 	ok = false
// 	return
// }
// redis沒有用戶資料，查詢資料表

// 先判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
// value, _ := redisConn.HashGetCache(config.LINE_USERS+userID, "ip")
// if value != "" {
// 	// redis已有用戶資訊
// 	user = models.UserModel{
// 		UserID: userID,
// 		Ip:     value,
// 	}
// } else if value == "" {
// 將用戶資訊加入redis中
// values := []interface{}{config.LINE_USERS + user.UserID}
// values = append(values, "user_id", lineUser.UserID)
// values = append(values, "name", lineUser.Name)
// values = append(values, "phone", lineUser.Phone)
// values = append(values, "email", lineUser.Email)
// values = append(values, "avatar", lineUser.Avatar)
// values = append(values, "identify", lineUser.Identify)
// values = append(values, "friend", lineUser.Friend)
// values = append(values, "ip", lineUser.Ip)
// if err := redisConn.HashMultiSetCache(values); err != nil {
// 	ok = false
// 	return
// }
// 設置過期時間
// if err := redisConn.SetExpire(config.LINE_USERS+user.UserID,
// 	strconv.Itoa(config.GetSessionLifeTime())); err != nil {
// 	ok = false
// 	return
// }

// 查詢到用戶資訊，將用戶資訊加入redis
// if err = redisConn.SetCache(config.HILIVE_USERS+userID, user.Ip,
// 	"EX", config.GetSessionLifeTime()); err != nil {
// 	ok = false
// 	return
// }

// if _, _ = ctx.Request.Cookie(sessionName); ctx.Request.Method != "GET" {
// 	ctx.Redirect(http.StatusFound, config.LOGIN_URL)
// } else {
// ctx.Redirect(http.StatusFound, config.LOGIN_URL)

// h := `<script>
// 	if (typeof(swal) === "function") {
// 		swal({
// 			type: "info",
// 			title: "login info",
// 			text: "登入逾時，請重新登入",
// 			showCancelButton: false,
// 			confirmButtonColor: "#3c8dbc",
// 			confirmButtonText: '` + "got it" + `',
// 		})
// 		setTimeout(function(){ location.href = "` + config.LOGIN_URL + `"; }, 2000);
// 	} else {
// 		alert("登入逾時，請重新登入")
// 		location.href = "` + config.LOGIN_URL + `"
// 	}
// </script>`
// h := `<script>
// 			alert("登入逾時，請重新登入")
// 			location.href = "` + config.LOGIN_URL + `"
// 	</script>`

// ctx.Status(http.StatusOK)
// ctx.Header("Content-Type", "text/html; charset=utf-8")
// ctx.Redirect(http.StatusFound, config.LOGIN_URL)
// }

// if ctx.Request.Method != "GET" {
// 	ctx.JSON(http.StatusForbidden, map[string]interface{}{
// 		"code": http.StatusForbidden,
// 		"msg":  "permission denied",
// 	})
// }
// h := `<div class="missing-content">
// 	<div class="missing-content-title">403</div>
// 	<div class="missing-content-title-subtitle">抱歉, 您沒有權限訪問此網頁.</div>
// </div>
// <style>
// .missing-content {
// 	padding: 48px 32px;
// }
// .missing-content-title {
// 	color: rgba(0,0,0,.85);
// 	font-size: 54px;
// 	line-height: 1.8;
// 	text-align: center;
// }
// .missing-content-title-subtitle {
// 	color: rgba(0,0,0,.45);
// 	font-size: 18px;
// 	line-height: 1.6;
// 	text-align: center;
// }
// </style>`
// tmpl, err := template.New("").Parse(h)
// if err != nil {
// 	log.Println("解析模板發生錯誤")
// }
// if err = tmpl.Execute(ctx.Writer, nil); err != nil {
// 	log.Println("使用權限模板發生錯誤")
// }
// ctx.Status(http.StatusOK)
// ctx.Header("Content-Type", "text/html; charset=utf

// if _, session, err = InitSession(ctx, conn, sessionName); err != nil {
// 	return user, false, false
// }
// session.Context.SetSameSite(4)
// session.Context.SetCookie(session.Cookie, session.SessionID, config.GetSessionLifeTime(), "/", "", true, false)

// if sessionName == "chatroom_session" { // 聊天室session
// 	if userID, ok = session.SessionValues["chatroom"].(string); !ok {
// 		return user, false, false
// 	}

// 	// chatroom_session可能儲存平台user或line user資料表id
// 	if user, ok = GetLineUser(userID, conn); !ok {
// 		if user, ok = GetUser(userID, conn); !ok {
// 			return user, false, false
// 		}
// 	}
// 	return user, true, true
// }

// hilive平台session
// if userID, ok = session.SessionValues["hilive"].(string); !ok {
// 	return user, false, false
// }
// if user, ok = GetUser(userID, conn); !ok {
// 	return user, false, false
// }
