package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	social "hilive/line-login"
	"hilive/models"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/response"
	"hilive/modules/service"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/template/types"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/crypto/bcrypt"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	MaxRetries     = 10000000000            // 重試次數
	RetryDelay     = 100 * time.Millisecond // 每次重試的延遲
	LockExpiration = 2                      // 鎖的有效期
)

// Handler API 處理程序
type Handler struct {
	config    *config.Config   // 平台基本設置
	dbConn    db.Connection    // 資料庫CRUD功能
	redisConn cache.Connection // redis
	mongoConn mongo.Connection // mongodb
	services  service.List     // 記錄資料庫、token...等資訊
	tableList table.List       // 記錄所有欄位、表單資訊
}

// route 模板需要的路徑
type route struct {
	Prefix            string // 前綴
	Index             string // 首頁
	Method            string // 方法
	Register          string // 註冊
	Login             string // 登入
	Retrieve          string // 忘記密碼
	Verification      string // 手機驗證
	VerificationCheck string // 手機驗證檢查
	Logout            string // 登出
	Back              string // 返回
	Admin             string // 管理員頁面
	Activity          string // 活動頁面
	User              string // 個人頁面
	ApplySign         string // 報名簽到
	// Info            string // 活動資訊頁面
	Overview    string // 總覽頁面
	Introduce   string // 活動介紹頁面
	Schedule    string // 活動行程頁面
	Guest       string // 活動嘉賓頁面
	Material    string // 活動資料頁面
	Sensitivity string // 敏感詞頁面

	Attend      string // 遊戲參加人員頁面
	Winning     string // 中獎人員頁面
	Host        string // 主持人活動頁面
	Bind        string // 綁定頁面
	StaffManage string // 人員管理頁面
	New         string // 新增頁面
	Edit        string // 編輯頁面
	Chatroom    string // 聊天室頁面

	Question       string // 提問牆頁面
	QuestionQRcode string // 提問牆QR code
	Game           string // 遊戲頁面
	QRcode         string // QRcode頁面
	// GameInfo        string // 遊戲資訊頁面
	Prize      string // 獎品設置頁面
	POST       string // 新增資料 API
	PUT        string // 編輯資料 API
	PUT2       string // 編輯資料 API
	PATCH      string // 更新設置 API
	DELETE     string // 刪除資料 API
	Reset      string // 重置遊戲資料 API
	Export     string // 匯出檔案
	Import     string // 匯入檔案
	Team       string // 隊伍設置頁面
	VoteOption string // 投票選項設置頁面
	Line       string // LINE驗證頁面
	Facebook   string // facebook驗證頁面
	Gmail      string // gmail驗證頁面
	Customize  string // 自定義驗證頁面
	QuickStart string // 快速設置頁面
	// Message      string // 訊息牆
	Topic string // 主題牆
	// Danmu        string // 彈幕
	SpecialDanmu string // 特殊彈幕
	// Picture      string // 圖片牆
	Holdscreen string // 霸屏
	Select     string // 選擇項目頁面
	// General      string // 一般簽到
	// Threed       string // 立體簽到
	// Countdown    string // 倒數簽到
	// Lottery      string // 幸運短盤

	// 簽到展示
	General        string // 一般簽到牆頁面
	Threed         string // 立體簽到牆頁面
	Signname       string // 簽名牆頁面
	SignnameCheck  string // 簽名牆審核頁面
	SignnameQRcode string // 簽名牆QR code

	// 遊戲
	Redpack           string // 搖紅包頁面
	Ropepack          string // 套紅包頁面
	WhackMole         string // 打地鼠頁面
	Lottery           string // 遊戲抽獎頁面
	Monopoly          string // 超級大富翁頁面
	QA                string // 快問快答頁面
	Tugofwar          string // 拔河遊戲頁面
	Bingo             string // 賓果遊戲頁面
	GachaMachine      string // 扭蛋機遊戲頁面
	Vote              string // 投票遊戲頁面
	DrawNumbers       string // 抽號碼頁面
	ThreedDrawNumbers string // 3D搖號抽獎頁面

	// 遊戲qrcode
	RedpackQRcode      string // 搖紅包QR code
	RopepackQRcode     string // 套紅包QR code
	WhackMoleQRcode    string // 打地鼠QR code
	LotteryQRcode      string // 遊戲抽獎QR code
	MonopolyQRcode     string // 超級大富翁QR code
	QAQRcode           string // 快問快答QR code
	TugofwarQRcode     string // 拔河遊戲QR code
	BingoQRcode        string // 賓果遊戲QR code
	GachaMachineQRcode string // 扭蛋機遊戲QR code
	VoteQRcode         string // 投票遊戲QR code

	// 遊戲新增頁面
	RedpackNew      string
	RopepackNew     string
	WhackMoleNew    string
	LotteryNew      string
	DrawNumbersNew  string
	MonopolyNew     string
	QANew           string
	TugofwarNew     string
	BingoNew        string
	GachaMachineNew string

	// 遊戲編輯頁面
	RedpackEdit      string
	RopepackEdit     string
	WhackMoleEdit    string
	LotteryEdit      string
	DrawNumbersEdit  string
	MonopolyEdit     string
	QAEdit           string
	TugofwarEdit     string
	BingoEdit        string
	GachaMachineEdit string

	// 遊戲獎品頁面
	RedpackPrize      string
	RopepackPrize     string
	WhackMolePrize    string
	LotteryPrize      string
	DrawNumbersPrize  string
	MonopolyPrize     string
	QAPrize           string
	TugofwarPrize     string
	BingoPrize        string
	GachaMachinePrize string

	// 遊戲增刪改查API
	RedpackAPI      string
	RopepackAPI     string
	WhackMoleAPI    string
	LotteryAPI      string
	DrawNumbersAPI  string
	MonopolyAPI     string
	QAAPI           string
	TugofwarAPI     string
	BingoAPI        string
	GachaMachineAPI string
}

// Chatroom 聊天室設置
type Chatroom struct {
	// 下一次只留Records、Redpacks、Ropepacks、WhackMole
	// Overviews []models.OverviewModel // 活動總覽
	// Introduces []models.IntroduceModel      // 活動介紹
	// Schedules  []models.ScheduleModel       // 活動行程
	// Guests     []models.GuestModel          // 活動嘉賓
	// Materials  []models.MaterialModel       // 活動資料

	// -----上面都刪除
	// Redpacks  []models.RedpackModel        // 搖紅包所有遊戲資訊
	// Ropepacks []models.RopepackModel       // 套紅包所有遊戲資訊
	// WhackMoles []models.WhackMoleModel      // 打地鼠所有遊戲資訊
	Lotterys      []models.GameModel           // 遊戲抽獎所有遊戲資訊
	Redpacks      []models.GameModel           // 搖紅包所有遊戲資訊
	Ropepacks     []models.GameModel           // 套紅包所有遊戲資訊
	WhackMoles    []models.GameModel           // 打地鼠所有遊戲資訊
	Monopolys     []models.GameModel           // 超級大富翁所有遊戲資訊
	QAs           []models.GameModel           // 快問快答所有遊戲資訊
	DrawNumbers   []models.GameModel           // 搖號抽獎所有遊戲資訊
	Tugofwars     []models.GameModel           // 拔河遊戲所有遊戲資訊
	Bingos        []models.GameModel           // 賓果遊戲所有遊戲資訊
	GachaMachines []models.GameModel           // 扭蛋機遊戲所有遊戲資訊
	Votes         []models.GameModel           // 投票遊戲所有遊戲資訊
	Records       []models.ChatroomRecordModel // 歷史紀錄
}

// StaffManage 人員管理頁面參數
type StaffManage struct {
	StaffList    types.InfoList          // 人員資料
	Staffs       []models.ApplysignModel // 人員資料
	Blacks       []models.ApplysignModel // 黑名單人員資料
	IsFilterPage bool                    // 該頁面是否為過濾頁面
	Rounds       []string                // 所有輪次資料(過濾用)
}

// executeParam 模板需要的參數
type executeParam struct {
	Action                 string                // 動作
	ActivityModel          models.ActivityModel  // 活動資訊
	ApplysignModel         models.ApplysignModel // 報名簽到資訊
	ActivityJSON           string                // 活動資訊(json資料)
	ApplysignCustomizeJSON string                // 自定義欄位資訊(json資料)
	GameJSON               string                // 遊戲資訊(json資料)
	WinningSaffsJSON       string                // 中獎人員紀錄(json資料)
	UserJSON               string                // 人員資料(json資料)
	User                   models.LoginUser      // 用戶資訊
	Role                   string                // 角色
	// HasHeader              bool                  // 是否顯示header
	HasSidebar  bool            // 是否顯示sidebar
	PanelInfo   table.PanelInfo // 面板資訊
	FieldList   types.FieldList // 欄位資訊
	StaffManage StaffManage     // 人員管理頁面參數
	// StaffList       types.InfoList   // 人員資料
	// BlackList       types.InfoList   // 黑名單人員資料
	FormInfo        table.FormInfo   // 表單資訊
	SettingFormInfo table.FormInfo   // 表單資訊(用於額外需要設置的表單)
	FormInfos       []table.FormInfo // 多個表單資訊(用於編輯資料)
	Menu            interface{}      // sidebar menu
	Config          *config.Config   // 基本設置
	Route           route            // 路徑設置
	ActivityID      string           // 辨識活動ID
	GameID          string           // 辨識遊戲ID
	ChannelID       string           // 辨識頻道ID
	AttendPeople    int              // 參加人數
	People          int              // 人數
	Status          string           // 狀態
	Chatroom        Chatroom         // 聊天室設置
	IsBlack         bool             // 是否為黑名單人員
	IsSign          bool             // 是否簽到
	IsFirst         bool             // 是否第一次
	IsAdmin         bool             // 是否有管理後台權限
	CanAdd          bool             // 是否有新增權限
	CanEdit         bool             // 是否有編輯權限
	CanDelete       bool             // 是否有刪除權限
	Token           string           // CSRF驗證
	Error           string           // 錯誤訊息
	// Number          int              // 抽獎編號
}

// acquireLock 嘗試獲取 Redis 鎖
func (h *Handler) acquireLock(lockKey string, expiration int) (interface{}, error) {
	return h.redisConn.SetCache(lockKey, "locked", "NX", "EX", expiration)
}

// releaseLock 釋放 Redis 鎖
func (h *Handler) releaseLock(lockKey string) (interface{}, error) {
	return h.redisConn.DelCache(lockKey)
}

// GetLoginUser 利用cookie取得用戶資訊
func (h *Handler) GetLoginUser(request *http.Request, sessionName string) models.LoginUser {
	var (
		user models.LoginUser
		// ip                                      = utils.ClientIP(ctx.Request)
	)
	cookie, err := request.Cookie(sessionName)
	if err != nil {
		return models.LoginUser{}
	}

	// 解碼cookie值
	decode, err := utils.Decode([]byte(cookie.Value))
	if err != nil {
		return models.LoginUser{}
	}

	params, err := url.ParseQuery(string(decode))
	if err != nil {
		return models.LoginUser{}
	}

	if sessionName == "hilive_session" {
		// *****可以同時登入(暫時拿除)*****
		// userModel, err := models.DefaultUserModel().SetDbConn(h.dbConn).
		// 	SetRedisConn(h.redisConn).Find(true, true, utils.ClientIP(request),
		// 	"users.user_id", params.Get("user_id"))
		// *****可以同時登入(暫時拿除)*****

		// *****可以同時登入(新)*****
		userModel, err := models.DefaultUserModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, true,
				"users.user_id", params.Get("user_id"))
		// *****可以同時登入(新)*****
		if err != nil || userModel.UserID == "" {
			return models.LoginUser{}
		}

		user.UserID = userModel.UserID
		user.Name = userModel.Name
		user.Phone = userModel.Phone
		user.Email = userModel.Email
		user.Avatar = userModel.Avatar
		// user.Bind = userModel.Bind
		user.Cookie = userModel.Cookie
		// user.Identify = params.Get("identify")
		// user.Friend = params.Get("friend")
		// *****可以同時登入(確定拿除)*****
		// user.Ip = userModel.Ip
		// *****可以同時登入(確定拿除)*****
		user.Table = params.Get("table")
		user.MaxActivity = userModel.MaxActivity
		user.MaxActivityPeople = userModel.MaxActivityPeople
		user.MaxGamePeople = userModel.MaxGamePeople
		user.Permissions = userModel.Permissions     // 權限
		user.Activitys = userModel.Activitys         // 活動資訊(包含活動權限)
		user.ActivityMenus = userModel.ActivityMenus // 菜單

		user.LineBind = userModel.LineBind
		user.FbBind = userModel.FbBind
		user.GmailBind = userModel.GmailBind
	} else if sessionName == "chatroom_session" {
		lineUser, err := models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, "", "user_id", params.Get("user_id"))
		if err != nil || lineUser.UserID == "" {
			return models.LoginUser{}
		}

		user.UserID = lineUser.UserID
		user.Name = lineUser.Name
		user.Phone = lineUser.Phone
		user.Email = lineUser.Email
		user.Avatar = lineUser.Avatar
		// user.Bind = params.Get("bind")
		// user.Cookie = params.Get("cookie")
		user.Identify = lineUser.Identify
		user.Friend = lineUser.Friend
		// *****可以同時登入(確定拿除)*****
		// user.Ip = params.Get("ip")
		// *****可以同時登入(確定拿除)*****
		user.Table = params.Get("table")
	}
	return user
}

// checkPermissions 判斷是否有新增編輯刪除權限
func (h *Handler) checkPermission(user models.LoginUser, activityID, path string) (canAdd, canEdit, canDelete bool) {
	// 判斷用戶權限
	for _, permisssion := range user.Permissions {
		if canAdd == false || canEdit == false || canDelete == false {
			for _, httppath := range permisssion.HTTPPath {
				// fmt.Println("測試權限: ", path, strings.Contains(path, ""), strings.Contains(path, "123"))
				if httppath == "admin" || httppath == "*" {
					canAdd = true
					canEdit = true
					canDelete = true
					break
				} else if strings.Contains(path, httppath) && httppath != "" {
					if strings.Contains(permisssion.HTTPMethod, "POST") {
						canAdd = true
					}
					if strings.Contains(permisssion.HTTPMethod, "PUT") {
						canEdit = true
					}
					if strings.Contains(permisssion.HTTPMethod, "DELETE") {
						canDelete = true
					}
				}
			}
		} else {
			break
		}
	}

	// 判斷活動權限
	if canAdd == false || canEdit == false || canDelete == false {
		for _, activity := range user.Activitys {
			if activity.ActivityID == activityID {
				for _, permisssion := range activity.Permissions {
					if canAdd == false || canEdit == false || canDelete == false {
						for _, httppath := range permisssion.HTTPPath {
							if httppath == "admin" || httppath == "*" {
								canAdd = true
								canEdit = true
								canDelete = true
								break
							} else if strings.Contains(path, httppath) && httppath != "" {
								if strings.Contains(permisssion.HTTPMethod, "POST") {
									canAdd = true
								}
								if strings.Contains(permisssion.HTTPMethod, "PUT") {
									canEdit = true
								}
								if strings.Contains(permisssion.HTTPMethod, "DELETE") {
									canDelete = true
								}
							}
						}
					} else {
						break
					}
				}
			}
		}
	}
	return
}

// execute 執行模板
func (h *Handler) execute(ctx *gin.Context, tempName string, htmlTmpl string, param executeParam) {
	tmpl, err := template.New(tempName).Parse(htmlTmpl)
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: fmt.Sprintf("解析發生錯誤: %s", err),
		})

		return
	}
	if err = tmpl.ExecuteTemplate(ctx.Writer, tempName, param); err != nil {
		// log.Println(err)
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: fmt.Sprintf("讀取模板發生錯誤: %s", err),
		})

		return
	}
	ctx.Status(http.StatusOK)
	return
}

// executeHTML 執行模板
func (h *Handler) executeHTML(ctx *gin.Context, route string, param executeParam) {
	tmpl := template.Must(template.ParseFiles(route))
	tmpl.Execute(ctx.Writer, param)
	ctx.Status(http.StatusOK)
	return
}

// executeErrorHTML 執行錯誤模板
func (h *Handler) executeErrorHTML(ctx *gin.Context, err string) {
	// logger.LoggerError(ctx, err)

	models.DefaultErrorLogModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Add(models.EditErrorLogModel{
			UserID:    utils.ClientIP(ctx.Request),
			Code:      http.StatusInternalServerError,
			Method:    ctx.Request.Method,
			Path:      ctx.Request.URL.Path,
			Message:   err,
			PathQuery: ctx.Request.URL.RawQuery,
		})

	tmpl := template.Must(template.ParseFiles("./hilives/hilive/views/error-message.html"))
	tmpl.Execute(ctx.Writer, executeParam{
		Error: err,
	})
	ctx.Status(http.StatusOK)
	return
}

// executeAuth 執行登入、註冊模板
// func (h *Handler) executeAuth(ctx *gin.Context, tmpl string) {
// 	h.execute(ctx, "", tmpl, executeParam{
// 		Route: route{
// 			Prefix:   h.config.Prefix,
// 			Login:    config.LOGIN_API_URL,
// 			Register: config.REGISTER_API_URL,
// 			Retrieve: config.RETRIEVE_API_URL,
// 		},
// 	})
// }

// NewHandler 預設Handler
func NewHandler(cfg *config.Config, services service.List, dbConn db.Connection,
	redisConn cache.Connection, mongoConn mongo.Connection, list table.List) *Handler {
	return &Handler{
		config:    cfg,
		services:  services,
		dbConn:    dbConn,
		redisConn: redisConn,
		mongoConn: mongoConn,
		tableList: list,
	}

}

// GetTable 取得table
func (h *Handler) GetTable(ctx *gin.Context, prefix string) (table.Table, string) {
	return h.tableList[prefix](), prefix
}

// NewSocialClient 取得*social.Client
func NewSocialClient(ctx *gin.Context, id, secret string) (client *social.Client) {
	if client, err := social.New(id, secret); err != nil {
		panic("錯誤: 取得*social.Client發生問題")
	} else {
		return client
	}
}

// NewLineBotClient 取得*linebot.Client
func NewLineBotClient(ctx *gin.Context, secret, token string) (bot *linebot.Client) {
	if bot, err := linebot.New(secret, token); err != nil {
		// fmt.Println("錯誤了嗎?")
		panic("錯誤: 取得*linebot.Client發生問題")
	} else {
		return bot
	}
}

// pushMessage 發送訊息
func pushMessage(ctx *gin.Context, secret, token, userID, message string) error {
	if _, err := NewLineBotClient(ctx, secret, token).PushMessage(userID,
		linebot.NewTextMessage(message)).Do(); err != nil {
		return errors.New("錯誤: 傳送Line訊息發生問題，請重新傳送")
	}
	return nil
}

// checkState 檢查是否存在state
// func checkState(h *Handler, state string) bool {
// 	// for i := 0; i < len(states); i++ {
// 	// 	if states[i] == state {
// 	// 		host := hosts[i]

// 	// 		// 清除陣列裡資料
// 	// 		states = append(states[:i], states[i+1:]...)
// 	// 		hosts = append(hosts[:i], hosts[i+1:]...)
// 	// 		return host, true
// 	// 	}
// 	// }

// 	// log.Println("states裡有多少資料", h.redisConn.SetCard(config.LINE_STATES_REDIS))

// 	// 判斷state參數是否存在redis中
// 	isExist := h.redisConn.SetIsMember(config.LINE_STATES_REDIS, state)
// 	if isExist {
// 		// 存在，清除redis裡的資料
// 		h.redisConn.SetRem([]interface{}{config.LINE_STATES_REDIS, state})

// 		// log.Println("存在，清除後剩多少資料: ", h.redisConn.SetCard(config.LINE_STATES_REDIS))
// 		return true
// 	}

// 	return false
// }

// checkUser 檢查手機與密碼是否符合
func (h *Handler) checkUser(password string, phone string) (user models.UserModel, ok bool) {
	// *****可以同時登入(暫時拿除)*****
	// if user, _ = models.DefaultUserModel().SetDbConn(h.dbConn).
	// 	Find(false, true, "", "users.phone", phone); user.ID == int64(0) {
	// 	ok = false
	// 	return
	// }
	// *****可以同時登入(暫時拿除)*****

	// *****可以同時登入(新)*****
	if user, _ = models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(false, true, "users.phone", phone); user.ID == int64(0) {
		ok = false
		return
	}
	// *****可以同時登入(新)*****

	// if comparePassword(password, user.Password) || password == user.Password {
	if comparePassword(password, user.Password) {
		ok = true
		// TODO: 目前沒有權限相關問題，先拿掉角色權限菜單判斷
		// user = user.GetRoles().GetPermissions().GetMenus()
		return
	}
	ok = false
	return
}

// checkApplysignUser 自定義簽到人員檢查用戶與密碼是否符合
func (h *Handler) checkApplysignUser(activityID string, password string) (applysignUserModel models.ApplysignUserModel, ok bool) {

	if applysignUserModel, _ = models.DefaultApplysignUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAccount(activityID, password); applysignUserModel.ID == int64(0) {
		ok = false
		return
	}
	ok = true
	return
}

// comparePassword 檢查密碼是否相符
func comparePassword(comPwd, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}

// getGameInfo 取得活動資訊(從redis取得，如果沒有才執行資料表查詢)
func (h *Handler) getActivityInfo(isRedis bool, activityID string) (activityModel models.ActivityModel, err error) {
	activityModel, err = models.DefaultActivityModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(isRedis, activityID)
	if err != nil || activityModel.ID == 0 {
		return models.ActivityModel{}, errors.New("錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
	}

	return
}

// getGameInfo 取得允許重複中獎、人數上限、中獎機率相關遊戲設置資訊(從redis取得，如果沒有才執行資料表查詢)
func (h *Handler) getGameInfo(gameID, game string) (gameModel models.GameModel) {
	gameModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, gameID, game)
	if err != nil || gameModel.ID == 0 {
		return models.GameModel{}
	}

	return
}

// getTeamInfo 取得隊伍資訊(從redis取得)
func (h *Handler) getTeamInfo(gameID string) models.GameModel {
	gameModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindTeam(true, gameID)
	if err != nil {
		return models.GameModel{}
	}

	return gameModel
}

// getProfile 取得用戶資料
func getProfile(client *linebot.Client, user string) (res *linebot.UserProfileResponse,
	link *linebot.LinkTokenResponse, err error) {
	if res, err = client.GetProfile(user).Do(); err != nil {
		// 官方帳號資訊錯誤會在這裡判斷
		// 封鎖時也會在這裡判斷
		err = errors.New("錯誤: 取得用戶資料發生問題，請重新掃描QRcode進行驗證(錯誤or被封鎖)")
		return
	}
	if link, err = client.IssueLinkToken(user).Do(); err != nil {
		err = errors.New("錯誤: 取得用戶資料發生問題，請重新掃描QRcode進行驗證")
		return
	}
	return
}

// getPrizes 獎品資訊(redis)
func (h *Handler) getPrizes(gameID string) (prizes []models.PrizeModel, err error) {
	// 所有獎品資訊
	prizes, err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindPrizes(true, gameID)
	if err != nil {
		return nil, errors.New("錯誤: 無法取得獎品資訊，請重新操作")
	}
	return prizes, nil
}

// getPrizesAmount 獎品總數
func (h *Handler) getPrizesAmount(gameID string) (amount int64, err error) {
	// 取得剩餘數量大於0的獎品資訊
	amount, err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindPrizesAmount(true, gameID)
	if err != nil {
		return 0, errors.New("錯誤: 無法取得獎品數量資訊，請重新操作")
	}
	return
}

// getUserGameRecords 取得用戶遊戲紀錄(包含中獎與未中獎)
// 回傳該活動的中獎紀錄、相同遊戲種類中獎紀錄、遊戲紀錄(包含中獎與未中獎)、遊戲中獎紀錄
func (h *Handler) getUserGameRecords(activityID, gameID, game, userID string) (
	[]models.PrizeStaffModel, []models.PrizeStaffModel,
	[]models.PrizeStaffModel, []models.PrizeStaffModel, []models.PrizeStaffModel, error) {
	var (
		allRecords     = make([]models.PrizeStaffModel, 0) // 用戶所有遊戲紀錄
		activityPrizes = make([]models.PrizeStaffModel, 0) // 用戶活動中獎紀錄
		gamePrizes     = make([]models.PrizeStaffModel, 0) // 用戶相同類型遊戲中獎紀錄
		gameRecords    = make([]models.PrizeStaffModel, 0) // 用戶遊戲紀錄
		winRecords     = make([]models.PrizeStaffModel, 0) // 用戶該遊戲中獎紀錄
	)
	// 用戶遊戲紀錄
	allRecords, err := models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindUserGameRecords(true, activityID, userID)
	if err != nil {
		return allRecords, activityPrizes, gamePrizes, gameRecords, winRecords, err
	}

	for _, record := range allRecords {
		// 加入遊戲紀錄(不管是否中獎)
		if record.GameID == gameID {
			gameRecords = append(gameRecords, record)
		}

		// 判斷是否中獎
		if (record.PrizeMethod != "" && record.PrizeMethod != "thanks") &&
			(record.PrizeType != "" && record.PrizeType != "thanks") &&
			record.PrizeID != "" {
			// 加入該活動中獎紀錄中
			activityPrizes = append(activityPrizes, record)

			// 判斷是否相同遊戲種類
			if record.Game == game {
				// 加入遊戲種類中獎紀錄中
				gamePrizes = append(gamePrizes, record)

				// 判斷是否同場次遊戲
				if record.GameID == gameID {
					winRecords = append(winRecords, record)
				}
			}
		}
	}
	return allRecords, activityPrizes, gamePrizes, gameRecords, winRecords, nil
}

// getUser 利用session值取得用戶資料、角色
func (h *Handler) getUser(session, userID string) (models.LineModel, string, error) {
	decode, err := utils.Decode([]byte(session))
	if err != nil {
		return models.LineModel{}, "", fmt.Errorf("錯誤: 解碼發生問題，請輸入有效的session值")
	}

	params, err := url.ParseQuery(string(decode))
	if err != nil {
		return models.LineModel{}, "", fmt.Errorf("錯誤: 無法取得用戶資訊，請輸入有效的session值")
	}

	// fmt.Println("params.Get(user_id):", params.Get("user_id"))

	// 判斷是否為管理員
	if params.Get("table") == "users" && params.Get("user_id") == userID {
		// fmt.Println("管理員")
		return models.LineModel{}, "host", nil
	}
	// 可能是平台的其他用戶而不是line用戶
	if params.Get("table") != "line_users" {
		return models.LineModel{}, "", fmt.Errorf("錯誤: 無法取得用戶資訊，請重新報名簽到活動")
	}

	// 取得LINE用戶資訊
	lineUser, err := models.DefaultLineModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, "", "user_id", params.Get("user_id"))
	if err != nil || lineUser.UserID == "" {
		return models.LineModel{}, "", fmt.Errorf("錯誤: 無法取得用戶資訊，請重新報名簽到活動")
	}

	return models.LineModel{
		UserID: lineUser.UserID,
		Name:   lineUser.Name,
		Avatar: lineUser.Avatar,
	}, "guest", nil

}

// getSignStaffs 簽到人員資訊
func (h *Handler) getSignStaffs(isRedis bool, isUserInfo bool,
	redisName, activityID string, limit, offset int64) (
	signStaffs []models.ApplysignModel, err error) {

	signStaffs, err = models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindSignStaffs(isRedis, isUserInfo, redisName, activityID, limit, offset)
	if err != nil {
		return signStaffs, err
	}
	return
}

// IsSign 是否簽到(目前都是從redis判斷，目前只有判斷SIGN_STAFFS_2_REDIS裡的資料)
func (h *Handler) IsSign(redisName, activityID, userID string) (isSign bool) {
	isSign = models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		IsSign(redisName, activityID, userID)
	return
}

// IsBlackStaff 是否為黑名單人員(目前都是從redis裡判斷)
func (h *Handler) IsBlackStaff(activityID, gameID, game, userID string) (isBlack bool) {
	isBlack = models.DefaultBlackStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		IsBlack(true, activityID, gameID,
			game, userID)
	return
}

// sendMessage 發送簡訊
func sendMessage(phone, message string) error {
	// 设置短信参数
	params := &openapi.CreateMessageParams{}
	params.SetTo("+886" + phone[1:])
	params.SetFrom(config.PHONE) // 发送者的 Twilio 电话号码
	params.SetBody(message)

	// 发送短信
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return errors.New("錯誤: 發送簡訊發生問題")
	}

	return nil
}

// // acquireLock 嘗試獲取 Redis 鎖
// func acquireLock(redisConn cache.Connection, lockKey string, expiration int) (interface{}, error) {
// 	return redisConn.SetCache(lockKey, "locked", "NX", "EX", expiration)
// }

// // releaseLock 釋放 Redis 鎖
// func releaseLock(redisConn cache.Connection, lockKey string) (interface{}, error) {
// 	return redisConn.DelCache(lockKey)
// }

// getUserWinningRecords 取得用戶中獎紀錄
// func (h *Handler) getUserWinningRecords(gameID string,
// 	userID string) (staffs []models.PrizeStaffModel, err error) {
// 	staffs, err = models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		FindUserWinningRecords(userID, "activity_staff_prize.game_id", gameID)
// 	if err != nil {
// 		return staffs, err
// 	}
// 	return
// }

// checkPhone 檢查電話號碼是否被使用過
// func (h *Handler) checkPhone(ctx *gin.Context, phone string) error {
// 	if phone != "" {
// 		if len(phone) > 2 {
// 			if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
// 				return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 			}
// 		} else {
// 			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 		}

// 		if user, err := models.DefaultLineModel().SetDbConn(h.dbConn).
// 			FindByPhone(phone); user.Phone != "" || err != nil {
// 			return errors.New("錯誤: 電話號碼已被註冊過，請輸入有效的手機號碼")
// 		}
// 	}
// 	return nil
// }

// add 增加用戶、角色、權限
// func (h *Handler) add(ctx *gin.Context, model models.EditUserModel) error {
// 	_, addErr := models.DefaultUserModel().SetDbConn(h.dbConn).Add(model)

// 	// _, addRoleErr := userModel.SetDbConn(h.dbConn).AddRole(role)
// 	// addPermissionErr := userModel.SetDbConn(h.dbConn).AddPermission(permission)
// 	if addErr != nil {
// 		return addErr
// 		// return errors.New("錯誤: 增加會員資料、權限錯誤，請重新註冊用戶")
// 	}
// 	return nil
// }

// getBlackStaffs 黑名單人員資訊
// func (h *Handler) getBlackStaffs(activityID, gameID,
// 	game string) (staffs []models.BlackStaffModel, err error) {
// 	staffs, err = models.DefaultBlackStaffModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).FindAll(false, activityID, gameID, game)
// 	if err != nil {
// 		return staffs, err
// 	}
// 	return
// }

// if game == "redpack" || game == "ropepack" {
// game == "whack_mole" {
// 取得剩餘數量大於0的獎品資訊
// prizes, err = models.DefaultPrizeModel().SetDbConn(h.dbConn).
// 	FindExistPrizes(game, "game_id", gameID)
// if err != nil {
// 	return nil, errors.New("錯誤: 無法取得獎品資訊，請重新操作")
// }
// } else if game == "lottery" {
// }
// 獎品總數
// for i := 0; i < len(prizes); i++ {
// 	prizeLength += int64(prizes[i].PrizeRemain)
// }

// ID:         score.ID,
// ActivityID: score.ActivityID,
// GameID:     score.GameID,
// Round:      score.Round,
// staffs, err = models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	FindWinningRecordsOrderByScore(gameID, round)

// getQuestions 依照排序取得提問資料
// func (h *Handler) getQuestions(field string, value interface{}, order string) ([]models.QuestionUserModel, error) {
// 	var (
// 		questionModel = models.DefaultQuestionUserModel().SetDbConn(h.dbConn)
// 		// questionModels = make([]models.QuestionUserModel, 0)
// 		questions = make([]models.QuestionUserModel, 0)
// 		err       error
// 	)

// 	if order == "send_time" {
// 		questions, err = questionModel.LeftJoinLineUsersOrderByTime(field, value)
// 	} else if order == "likes" {
// 		questions, err = questionModel.LeftJoinLineUsersOrderByLikes(field, value)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return questions, nil
// }

// getUserQuestions 取得用戶的提問資料
// func (h *Handler) getUserQuestions(activityID, userID string) ([]models.QuestionUserModel, error) {
// 	var (
// 		questionModel = models.DefaultQuestionUserModel().SetDbConn(h.dbConn)
// 		// questionModels = make([]models.QuestionUserModel, 0)
// 		questions = make([]models.QuestionUserModel, 0)
// 		err       error
// 	)
// 	questions, err = questionModel.LeftJoinLineUsersByUser(activityID, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return questions, nil
// }

// getUserLikesRecords 取得用戶按讚紀錄
// func (h *Handler) getUserLikesRecords(activityID, userID string) ([]models.QuestionLikesRecordModel, error) {
// 	var (
// 		questionLikesRecordModel = models.DefaultQuestionLikesRecordModel().SetDbConn(h.dbConn)
// 		// questionLikesRecordModels = make([]models.QuestionLikesRecordModel, 0)
// 		records = make([]models.QuestionLikesRecordModel, 0)
// 		err     error
// 	)
// 	records, err = questionLikesRecordModel.LeftJoinLineUsers(activityID, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return records, nil
// }

// for _, staff := range signStaffs {
// 	staffs = append(staffs, SignStaffModel{
// 		ID:         staff.ID,
// 		UserID:     staff.UserID,
// 		ActivityID: staff.ActivityID,
// 		Name:       staff.Name,
// 		Avatar:     staff.Avatar,
// 		Number:     staff.Number,
// 		// Status:     staff.Status,
// 		// ApplyTime:  staff.ApplyTime,
// 		// ReviewTime: staff.ReviewTime,
// 		// SignTime:   staff.SignTime,
// 	})
// }

// getWhackMoleInfo 取得打地鼠遊戲設置資訊(各獎項設置數量)
// func (h *Handler) getWhackMoleInfo(gameID string) (first, second, third, general int64) {
// 	gameModel, err := models.DefaultGameModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).FindGame(true, gameID)
// 	if err != nil || gameModel.ID == 0 {
// 		return
// 	}

// 	return gameModel.FirstPrize, gameModel.SecondPrize, gameModel.ThirdPrize, gameModel.GeneralPrize
// }

// getTopScoreStaffs 取得分數前n高的人員資料(從redis取得)
// func (h *Handler) getTopScoreStaffs(gameID string, round, limit int64) (staffs []models.PrizeStaffModel, err error) {
// 	// var (
// 	// 	staffs = make([]models.PrizeStaffModel, 0)
// 	// )
// 	scoreStaffs, err := models.DefaultScoreModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).Find(true, gameID, limit)
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得分數前n高的人員資料，請重新操作")
// 	}

// 	for _, score := range scoreStaffs {
// 		staffs = append(staffs, models.PrizeStaffModel{
// 			// ID:         score.ID,
// 			// ActivityID: score.ActivityID,
// 			// GameID:     score.GameID,
// 			UserID: score.UserID,
// 			Name:   score.Name,
// 			Avatar: score.Avatar,
// 			// Round:      score.Round,
// 			Score: score.Score,
// 		})
// 	}
// 	return
// }

// getWinningRecordsByGameID 透過game_id查詢該遊戲場次的中獎紀錄(join line_users、activity_prize join)
// func (h *Handler) getWinningRecordsByGameID(gameID string) ([]models.PrizeStaffModel, error) {
// 	// var (
// 	// 	staffs = make([]models.PrizeStaffModel, 0)
// 	// )

// 	staffs, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		FindWinningRecordsByGameID(gameID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return staffs, nil
// }

// getWinningRecordsOrderByScore 取得中獎紀錄(依照分數排序)
// func (h *Handler) getWinningRecordsOrderByScore(gameID, round string) ([]models.PrizeStaffModel, error) {
// 	// var (
// 	// 	staffs = make([]models.PrizeStaffModel, 0)
// 	// )

// 	staffs, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		FindWinningRecordsOrderByScore(gameID, round)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return staffs, nil
// }

// getWinningRecords 取得中獎紀錄(輪次 or 用戶)
// func (h *Handler) getWinningRecords(gameID string, field string,
// 	value interface{}) ([]models.PrizeStaffModel, error) {
// 	// var (
// 	// 	staffs = make([]models.PrizeStaffModel, 0)
// 	// )

// 	staffs, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		FindWinningRecords(gameID, field, value)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return staffs, nil
// }

// getDrawNumbersPrizes 取得抽號碼獎品資訊
// func (h *Handler) getDrawNumbersPrizes(activityID string) ([]DrawNumbersPrize, error) {
// 	var prizes = make([]DrawNumbersPrize, 0)

// 	prizeModels, err := models.DefaultDrawNumbersPrizeModel().SetDbConn(h.dbConn).
// 		FindAllPrizes(activityID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i := 0; i < len(prizeModels); i++ {
// 		var (
// 			prize DrawNumbersPrize
// 		)
// 		prize.ActivityID = prizeModels[i].ActivityID
// 		prize.PrizeID = prizeModels[i].PrizeID
// 		prize.Name = prizeModels[i].Name
// 		prize.PrizeType = prizeModels[i].PrizeType
// 		prize.Picture = prizeModels[i].Picture
// 		prize.Amount = prizeModels[i].Amount
// 		prize.Remain = prizeModels[i].Remain
// 		prize.Price = prizeModels[i].Price
// 		prize.Method = prizeModels[i].Method
// 		prize.Password = prizeModels[i].Password
// 		prizes = append(prizes, prize)
// 	}
// 	return prizes, nil
// }

// GetData 取得頁面所有資料
// func (h *Handler) GetData(ctx *gin.Context, params parameter.Parameters,
// 	panel table.Table, prefix string, otherFields []string) (table.Table, table.PanelInfo, error) {
// 	if panel == nil {
// 		panel, _ = h.GetTable(ctx, prefix)
// 	}
// 	panelInfo, err := panel.GetData(params, h.services, otherFields)
// 	if err != nil {
// 		return panel, p

// if game == "redpack" {
// 	redpackPrizes, err := models.DefaultRedpackPrizeModel().SetDbConn(h.dbConn).
// 		FindExistPrizres(gameID)
// 	if err != nil {
// 		return nil, prizeLength, errors.New("錯誤: 無法取得搖紅包獎品資訊，請重新操作")
// 	}

// 	for i := 0; i < len(redpackPrizes); i++ {
// 		prizeLength += int64(redpackPrizes[i].Remain)
// 		prizes = append(prizes, PrizeModel{
// 			ID:            redpackPrizes[i].ID,
// 			ActivityID:    redpackPrizes[i].ActivityID,
// 			PrizeID:       redpackPrizes[i].PrizeID,
// 			PrizeName:     redpackPrizes[i].Name,
// 			PrizeType:     redpackPrizes[i].PrizeType,
// 			PrizePicture:  redpackPrizes[i].Picture,
// 			PrizePrice:    redpackPrizes[i].Price,
// 			PrizeAmount:   redpackPrizes[i].Amount,
// 			PrizeRemain:   redpackPrizes[i].Remain,
// 			PrizeMethod:   redpackPrizes[i].Method,
// 			PrizePassword: redpackPrizes[i].Password,
// 		})
// 	}
// } else if game == "ropepack" {
// 	ropepackPrizes, err := models.DefaultRopepackPrizeModel().SetDbConn(h.dbConn).
// 		FindExistPrizres(gameID)
// 	if err != nil {
// 		return nil, prizeLength, errors.New("錯誤: 無法取得套紅包獎品資訊，請重新操作")
// 	}

// 	for i := 0; i < len(ropepackPrizes); i++ {
// 		prizeLength += int64(ropepackPrizes[i].Remain)
// 		prizes = append(prizes, PrizeModel{
// 			ID:            ropepackPrizes[i].ID,
// 			ActivityID:    ropepackPrizes[i].ActivityID,
// 			PrizeID:       ropepackPrizes[i].PrizeID,
// 			PrizeName:     ropepackPrizes[i].Name,
// 			PrizeType:     ropepackPrizes[i].PrizeType,
// 			PrizePicture:  ropepackPrizes[i].Picture,
// 			PrizePrice:    ropepackPrizes[i].Price,
// 			PrizeAmount:   ropepackPrizes[i].Amount,
// 			PrizeRemain:   ropepackPrizes[i].Remain,
// 			PrizeMethod:   ropepackPrizes[i].Method,
// 			PrizePassword: ropepackPrizes[i].Password,
// 		})
// 	}
// } else if game == "whack_mole" {
// 	whackMolePrizes, err := models.DefaultWhackMolePrizeModel().SetDbConn(h.dbConn).
// 		FindExistPrizres(gameID)
// 	if err != nil {
// 		return nil, prizeLength, errors.New("錯誤: 無法取得打地鼠獎品資訊，請重新操作")
// 	}

// 	for i := 0; i < len(whackMolePrizes); i++ {
// 		prizeLength += int64(whackMolePrizes[i].Remain)
// 		prizes = append(prizes, PrizeModel{
// 			ID:            whackMolePrizes[i].ID,
// 			ActivityID:    whackMolePrizes[i].ActivityID,
// 			PrizeID:       whackMolePrizes[i].PrizeID,
// 			PrizeName:     whackMolePrizes[i].Name,
// 			PrizeType:     whackMolePrizes[i].PrizeType,
// 			PrizePicture:  whackMolePrizes[i].Picture,
// 			PrizePrice:    whackMolePrizes[i].Price,
// 			PrizeAmount:   whackMolePrizes[i].Amount,
// 			PrizeRemain:   whackMolePrizes[i].Remain,
// 			PrizeMethod:   whackMolePrizes[i].Method,
// 			PrizePassword: whackMolePrizes[i].Password,
// 		})
// 	}
// } else if game == "draw_numbers" {
// 	drawNumbersModels, err := models.DefaultDrawNumbersPrizeModel().SetDbConn(h.dbConn).
// 		FindExistPrizres(gameID)
// 	if err != nil {
// 		return nil, prizeLength, errors.New("錯誤: 無法取得抽號碼獎品資訊，請重新操作")
// 	}

// 	for i := 0; i < len(drawNumbersModels); i++ {
// 		prizeLength += int64(drawNumbersModels[i].Remain)
// 		prizes = append(prizes, PrizeModel{
// 			ID:            drawNumbersModels[i].ID,
// 			ActivityID:    drawNumbersModels[i].ActivityID,
// 			PrizeID:       drawNumbersModels[i].PrizeID,
// 			PrizeName:     drawNumbersModels[i].Name,
// 			PrizeType:     drawNumbersModels[i].PrizeType,
// 			PrizePicture:  drawNumbersModels[i].Picture,
// 			PrizePrice:    drawNumbersModels[i].Price,
// 			PrizeAmount:   drawNumbersModels[i].Amount,
// 			PrizeRemain:   drawNumbersModels[i].Remain,
// 			PrizeMethod:   drawNumbersModels[i].Method,
// 			PrizePassword: drawNumbersModels[i].Password,
// 		})
// 	}
// } else if game == "lottery" {
// 	lotteryModels, err := models.DefaultLotteryPrizeModel().SetDbConn(h.dbConn).
// 		FindAllPrizes(gameID)
// 	if err != nil {
// 		return nil, prizeLength, errors.New("錯誤: 無法取得抽號碼獎品資訊，請重新操作")
// 	}

// 	for i := 0; i < len(lotteryModels); i++ {
// 		prizeLength += int64(lotteryModels[i].Remain)
// 		prizes = append(prizes, PrizeModel{
// 			ID:            lotteryModels[i].ID,
// 			ActivityID:    lotteryModels[i].ActivityID,
// 			PrizeID:       lotteryModels[i].PrizeID,
// 			PrizeName:     lotteryModels[i].Name,
// 			PrizeType:     lotteryModels[i].PrizeType,
// 			PrizePicture:  lotteryModels[i].Picture,
// 			PrizePrice:    lotteryModels[i].Price,
// 			PrizeAmount:   lotteryModels[i].Amount,
// 			PrizeRemain:   lotteryModels[i].Remain,
// 			PrizeMethod:   lotteryModels[i].Method,
// 			PrizePassword: lotteryModels[i].Password,
// 		})
// 	}
// }
// for _, staff := range prizeStaffs {
// 	staffs = append(staffs, PrizeStaffModel{
// 		ID:            staff.ID,
// 		ActivityID:    staff.ActivityID,
// 		GameID:        staff.GameID,
// 		PrizeID:       staff.PrizeID,
// 		UserID:        staff.UserID,
// 		Name:          staff.Name,
// 		Avatar:        staff.Avatar,
// 		PrizeName:     staff.PrizeName,
// 		PrizeType:     staff.PrizeType,
// 		PrizePicture:  staff.PrizePicture,
// 		PrizePrice:    staff.PrizePrice,
// 		PrizeMethod:   staff.PrizeMethod,
// 		PrizePassword: staff.PrizePassword,
// 		Round:         staff.Round,
// 		WinTime:       staff.WinTime[:len(staff.WinTime)-3], // 不顯示秒
// 		Status:        staff.Status,
// 	})
// }

// allow = gameModel.Allow
// people = gameModel.MaxPeople
// percent = gameModel.Percent
// second = gameModel.Second

// if strings.Contains(game, "redpack") {
// 	redpackModel, err := models.DefaultRedpackModel().SetDbConn(h.dbConn).Find(gameID)
// 	if err != nil || redpackModel.ID == 0 {
// 		return "", 0, 0, 0
// 	}
// 	allow = redpackModel.Allow
// 	people = redpackModel.People
// 	percent = redpackModel.Percent
// 	second = redpackModel.Second
// } else if strings.Contains(game, "ropepack") {
// 	ropepackModel, err := models.DefaultRopepackModel().SetDbConn(h.dbConn).Find(gameID)
// 	if err != nil || ropepackModel.ID == 0 {
// 		return "", 0, 0, 0
// 	}
// 	allow = ropepackModel.Allow
// 	people = ropepackModel.People
// 	percent = ropepackModel.Percent
// 	second = ropepackModel.Second
// } else if strings.Contains(game, "whack_mole") {
// 	whackMoleModel, err := models.DefaultWhackMoleModel().SetDbConn(h.dbConn).Find(gameID)
// 	if err != nil || whackMoleModel.ID == 0 {
// 		return "", 0, 0, 0
// 	}
// 	allow = whackMoleModel.Allow
// 	people = whackMoleModel.MaxPeople
// 	percent = 0
// 	second = whackMoleModel.Second
// }
// first = gameModel.FirstPrize
// second = gameModel.SecondPrize
// third = gameModel.ThirdPrize
// special = gameModel.SpecialPrize
// prize := PrizeModel{
// 	ID:            prizeModel[i].ID,
// 	ActivityID:    prizeModel[i].ActivityID,
// 	PrizeID:       prizeModel[i].PrizeID,
// 	PrizeName:     prizeModel[i].PrizeName,
// 	PrizeType:     prizeModel[i].PrizeType,
// 	PrizePicture:  prizeModel[i].PrizePicture,
// 	PrizePrice:    prizeModel[i].PrizePrice,
// 	PrizeAmount:   prizeModel[i].PrizeAmount,
// 	PrizeRemain:   prizeModel[i].PrizeRemain,
// 	PrizeMethod:   prizeModel[i].PrizeMethod,
// 	PrizePassword: prizeModel[i].PrizePassword,
// }

// 是否中獎
// if random >= len(prizes) {
// 	// 未中獎
// 	return prize, nil
// } else {
// 	// 中獎，取得獎品池中禮物資訊
// 	prize = prizes[random]
// }

// for i := 0; i < len(prizes); i++ {
// 	index += int(prizes[i].Remain)
// 	if index > random {
// 		prize = prizes[i]
// 		break
// 	}
// }
// for i := 0; i < len(prizes); i++ {
// 	index += int(prizes[i].Remain)
// 	if index > random {
// 		prize = prizes[i]
// 		break
// 	}
// }
// for _, staff := range prizeStaffs {
// 	staffs = append(staffs, PrizeStaffModel{
// 		ID:            staff.ID,
// 		ActivityID:    staff.ActivityID,
// 		GameID:        staff.GameID,
// 		PrizeID:       staff.PrizeID,
// 		UserID:        staff.UserID,
// 		Name:          staff.Name,
// 		Avatar:        staff.Avatar,
// 		PrizeName:     staff.PrizeName,
// 		PrizeType:     staff.PrizeType,
// 		PrizePicture:  staff.PrizePicture,
// 		PrizePrice:    staff.PrizePrice,
// 		PrizeMethod:   staff.PrizeMethod,
// 		PrizePassword: staff.PrizePassword,
// 		Round:         staff.Round,
// 		WinTime:       staff.WinTime[:len(staff.WinTime)-3], // 不顯示秒
// 		Status:        staff.Status,
// 	})
// }

// for _, staff := range prizeStaffs {
// 	staffs = append(staffs, PrizeStaffModel{
// 		ID:            staff.ID,
// 		ActivityID:    staff.ActivityID,
// 		GameID:        staff.GameID,
// 		PrizeID:       staff.PrizeID,
// 		UserID:        staff.UserID,
// 		Name:          staff.Name,
// 		Avatar:        staff.Avatar,
// 		PrizeName:     staff.PrizeName,
// 		PrizeType:     staff.PrizeType,
// 		PrizePicture:  staff.PrizePicture,
// 		PrizePrice:    staff.PrizePrice,
// 		PrizeMethod:   staff.PrizeMethod,
// 		PrizePassword: staff.PrizePassword,
// 		Round:         staff.Round,
// 		WinTime:       staff.WinTime[:len(staff.WinTime)-3], // 不顯示秒
// 		Status:        staff.Status,
// 	})
// }
// for _, staff := range prizeStaffs {
// 	staffs = append(staffs, PrizeStaffModel{
// 		ID:            staff.ID,
// 		ActivityID:    staff.ActivityID,
// 		GameID:        staff.GameID,
// 		PrizeID:       staff.PrizeID,
// 		UserID:        staff.UserID,
// 		Name:          staff.Name,
// 		Avatar:        staff.Avatar,
// 		PrizeName:     staff.PrizeName,
// 		PrizeType:     staff.PrizeType,
// 		PrizePicture:  staff.PrizePicture,
// 		PrizePrice:    staff.PrizePrice,
// 		PrizeMethod:   staff.PrizeMethod,
// 		PrizePassword: staff.PrizePassword,
// 		Round:         staff.Round,
// 		WinTime:       staff.WinTime[:len(staff.WinTime)-3], // 不顯示秒
// 		Status:        staff.Status,
// 	})
// }
// for _, question := range questionModels {
// 	questions = append(questions, QuestionUserModel{
// 		ID:         question.ID,
// 		ActivityID: question.ActivityID,
// 		UserID:     question.UserID,
// 		Name:       question.Name,
// 		Avatar:     question.Avatar,
// 		Content:    question.Content,
// 		Likes:      question.Likes,
// 		SendTime:   question.SendTime,
// 		Like:       question.Like,
// 	})
// }
// for _, question := range questionModels {
// 	questions = append(questions, QuestionUserModel{
// 		ID:         question.ID,
// 		ActivityID: question.ActivityID,
// 		UserID:     question.UserID,
// 		Name:       question.Name,
// 		Avatar:     question.Avatar,
// 		Content:    question.Content,
// 		Likes:      question.Likes,
// 		SendTime:   question.SendTime,
// 		Like:       question.Like,
// 	})
// }
// for _, record := range questionLikesRecordModels {
// 	records = append(records, QuestionLikesRecordModel{
// 		ID:         record.ID,
// 		ActivityID: record.ActivityID,
// 		QuestionID: record.QuestionID,
// 		UserID:     record.UserID,
// 		Name:       record.Name,
// 		Avatar:     record.Avatar,
// 	})
//

// var (

// user = GetLoginUser(ctx.Request, "chatroom_session")
// userid string
// ok     bool
// values map[string]interface{}
// )
// if session == "" {
// 	return models.LineModel{}, "", fmt.Errorf("錯誤: session值為空，無法取得用戶資訊，請輸入有效的session值")
// }

// sessionModel, err := db.Table(config.SESSION_TABLE).WithConn(h.dbConn).
// 	Where("session_id", "=", session).First()
// if err != nil || sessionModel == nil {
// 	return models.LineModel{}, "", fmt.Errorf("錯誤: session值無效、過期，請輸入有效的session值")
// }
// if err = json.Unmarshal([]byte(sessionModel["session_values"].(string)), &values); err != nil || len(values) <= 0 {
// 	return models.LineModel{}, "", fmt.Errorf("錯誤: json解碼發生錯誤，請重新操作")
// }
// if userid, ok = values["chatroom"].(string); !ok {
// 	return models.LineModel{}, "", fmt.Errorf("錯誤: 數值轉換錯誤，請重新操作")
// }

// if user, err = models.DefaultLineModel().SetDbConn(h.dbConn).
// 	Find("user_id", userid); err != nil |
