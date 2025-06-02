package engine

import (
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/logger"
	"hilive/modules/response"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 設置路由
func (eng *Engine) InitRouter() *Engine {
	//template render engine
	// engine := html.New("./templates", ".html")

	// app := fiber.New(fiber.Config{
	// 	Views: engine, //set as render engine
	// })
	// app.Get("/", mainPage)
	// // app.Listen(":3000", fiber.ListenConfig{})

	// eng.Gin = app

	// ------------------------------------------------------

	r := gin.Default()
	// r.Use(limits.RequestSizeLimiter(10000000))
	// r.Use(limit.Limit(10000000))

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := r.Group("").Use(response.CORS)

	// 測試用API
	// router.Static("/test", "./test")
	// router.GET("/test", eng.handler.Test)
	// router.GET("/ws/v1/test", eng.handler.Test)
	// router.GET("/1", eng.handler.Test1)
	// router.GET("/2", eng.handler.Test2)
	// router.GET("/3", eng.handler.Test3)
	// router.GET("/4", eng.handler.Test4)
	// router.Static("/3DGachaMachine/assets", "./hilives/hilive/assets/3DGachaMachine")

	// facebook驗證
	// router.GET("/facebook", eng.handler.HandleFacebookMain)
	// router.GET("/facebook/login", eng.handler.HandleFacebookLogin)
	// router.GET("/facebook/callback", eng.handler.HandleFacebookCallback)
	// gmail驗證
	// router.GET("/gmail", eng.handler.HandleGoogleMain)
	// router.GET("/gmail/login", eng.handler.HandleGoogleLogin)
	// router.GET("/gmail/callback", eng.handler.HandleGoogleCallback)
	// router.GET("/gmail", eng.handler.SendMail)
	// router.GET("/SMTP", eng.handler.SMTP)

	// 靜態檔案
	router.Static("/redpack/assets", "./hilives/hilive/assets/redpack")
	router.Static("/ropepack/assets", "./hilives/hilive/assets/ropepack")
	router.Static("/whack_mole/assets", "./hilives/hilive/assets/whack_mole")
	router.Static("/draw_numbers/assets", "./hilives/hilive/assets/draw_numbers")
	router.Static("/3Ddraw_numbers/assets", "./hilives/hilive/assets/3Ddraw_numbers")
	router.Static("/lottery/assets", "./hilives/hilive/assets/lottery")
	router.Static("/general/assets", "./hilives/hilive/assets/general")
	router.Static("/threed/assets", "./hilives/hilive/assets/threed")
	router.Static("/monopoly/assets", "./hilives/hilive/assets/monopoly")
	router.Static("/QA/assets", "./hilives/hilive/assets/QA")
	router.Static("/tugofwar/assets", "./hilives/hilive/assets/tugofwar")
	router.Static("/bingo/assets", "./hilives/hilive/assets/bingo")
	router.Static("/signname/assets", "./hilives/hilive/assets/signname")
	router.Static("/3DGachaMachine/assets", "./hilives/hilive/assets/3DGachaMachine_guest")
	router.Static("/3DGachaMachine/Build", "./hilives/hilive/assets/3DGachaMachine_host/Build")
	router.Static("/3DGachaMachine/StreamingAssets", "./hilives/hilive/assets/3DGachaMachine_host/StreamingAssets")
	router.Static("/3DGachaMachine/TemplateData", "./hilives/hilive/assets/3DGachaMachine_host/TemplateData")
	router.Static("/3Ddraw_numbers/Build", "./hilives/hilive/assets/3Ddraw_numbers/Build")
	router.Static("/3Ddraw_numbers/StreamingAssets", "./hilives/hilive/assets/3Ddraw_numbers/StreamingAssets")
	router.Static("/3Ddraw_numbers/TemplateData", "./hilives/hilive/assets/3Ddraw_numbers/TemplateData")
	router.Static("/vote/assets", "./hilives/hilive/assets/vote_guest")
	router.Static("/vote/Build", "./hilives/hilive/assets/vote_host/Build")
	router.Static("/vote/StreamingAssets", "./hilives/hilive/assets/vote_host/StreamingAssets")
	router.Static("/vote/TemplateData", "./hilives/hilive/assets/vote_host/TemplateData")
	// r.Use(static.ServeRoot("/QA/assets", "./assets/QA"))

	// 首頁
	// router.GET("/", eng.handler.ShowIndex)

	// Auth
	router.GET("/auth/redirect", eng.handler.AuthRedirect)                    // 導向驗證前的頁面
	router.GET("/auth/user", eng.handler.ShowAuthUser)                        // 報名簽到用戶補齊資料頁面
	router.GET("/v1/auth/callback", eng.handler.LoginCallback)                // LOGIN CALLBACK
	router.POST("/v1/message/line/callback", eng.handler.LineMessageCallback) // LINE MESSAGE CALLBACK

	// facebook驗證
	router.GET("/facebook/delete", eng.handler.DeleteHandler)              // 刪除回呼api(當用戶要求刪除資料時執行)
	router.GET("/facebook/delete/status", eng.handler.DeleteStatusHandler) // 刪除回呼api(當用戶要求刪除資料時執行)

	// LINE 選單
	router.GET("/activity/search", eng.handler.ShowLineRichmenu)           // 活動查詢
	router.GET("/activity/create", eng.handler.ShowLineRichmenu)           // 舉辦活動(填寫需求、聯繫我們)
	router.GET("/activity/create/:__prefix", eng.handler.ShowLineRichmenu) // 填寫需求、聯繫我們
	router.GET("/activity/case", eng.handler.ShowLineRichmenu)             // 場景案例(facebook、youtube、官網介紹)
	router.POST("/v1/activity/create/:__prefix", eng.handler.POST)         // 舉辦活動選單
	// router.GET("/activity/*__prefix", eng.handler.ShowLineRichmenu)
	// router.GET("/activity/search/:__prefix", eng.handler.ShowLineRichmenu) // 進行中活動、用戶報名活動

	// 報名簽到狀態處理
	router.GET("/applysign", eng.handler.Applysign)

	// 手機驗證
	router.POST("/v1/verification", eng.handler.Verification)
	router.POST("/v1/verification/check", eng.handler.VerificationCheck)
	// router.GET("/v1/verification/test", eng.handler.VerificationTest)
	// router.OPTIONS("/v1/verification", eng.handler.OPTIONS)
	// router.OPTIONS("/v1/verification/check", eng.handler.OPTIONS)

	// 註冊、登入、忘記密碼、更新用戶資料
	router.POST("/v1/login", eng.handler.Login)
	router.POST("/v1/register", eng.handler.Register)
	router.POST("/v1/retrieve", eng.handler.Retrieve)
	router.POST("/v1/applysign/login", eng.handler.ApplysignLogin) // 自定義簽到人員登入(活動主辦方匯入)
	router.POST("/v1/auth/user", eng.handler.AuthUser)             // 用戶補齊資料

	// 管理員
	router.GET("/v1/admin/:__prefix", eng.handler.GetAdmin)
	router.POST("/v1/admin/:__prefix", eng.handler.POST)
	router.PUT("/v1/admin/:__prefix", eng.handler.PUT)
	router.DELETE("/v1/admin/:__prefix", eng.handler.DELETE)
	router.OPTIONS("/v1/admin/:__prefix", eng.handler.OPTIONS)

	// 活動是否存在
	router.GET("/activity/isexist", eng.handler.IsExistActivity)

	// 聊天紀錄
	router.GET("/v1/chatroom/record", eng.handler.GetRecords)
	router.POST("/v1/chatroom/record", eng.handler.POST)
	router.PUT("/v1/chatroom/record", eng.handler.PUT)
	router.OPTIONS("/v1/chatroom/record", eng.handler.OPTIONS)

	// 提問紀錄
	router.GET("/v1/question/record", eng.handler.GetRecords)
	router.PUT("/v1/question/record", eng.handler.PUT)
	router.OPTIONS("/v1/question/record", eng.handler.OPTIONS)

	// 平台用戶
	router.GET("/v1/user", eng.handler.GetUserRole)
	router.POST("/v1/user", eng.handler.POST)
	router.PUT("/v1/user", eng.handler.PUT)
	router.OPTIONS("/v1/user", eng.handler.OPTIONS)

	// 自定義場景
	router.GET("/v1/customize_scene", eng.handler.GetCustomizeScenes)
	// 自定義模板
	router.GET("/v1/customize_template", eng.handler.GetCustomizeTemplates)

	// 參加用戶
	router.PUT("/v1/line_user", eng.handler.PUT)
	router.OPTIONS("/v1/line_user", eng.handler.OPTIONS)

	// 活動
	router.GET("/v1/activity", eng.handler.GetActivity)
	router.POST("/v1/activity", eng.handler.POST)
	router.PUT("/v1/activity", eng.handler.PUT)
	router.DELETE("/v1/activity", eng.handler.DELETE)
	router.OPTIONS("/v1/activity", eng.handler.OPTIONS)

	// 活動資訊 info
	router.GET("/v1/info/:__prefix", eng.handler.GetInfo)
	router.POST("/v1/info/:__prefix", eng.handler.POST)
	router.PUT("/v1/info/:__prefix", eng.handler.PUT)
	router.PATCH("/v1/info/:__prefix", eng.handler.PATCH)
	router.DELETE("/v1/info/:__prefix", eng.handler.DELETE)
	router.OPTIONS("/v1/info/:__prefix", eng.handler.OPTIONS)

	// 活動簽到 applysign
	router.GET("/v1/applysign", eng.handler.GetApplySign)
	router.GET("/v1/applysign/customize", eng.handler.GetApplysignCustomize)
	router.POST("/v1/applysign", eng.handler.POST)
	router.POST("/v1/applysign/export", eng.handler.EXPORT)
	router.PUT("/v1/applysign", eng.handler.PUT)
	router.DELETE("/v1/applysign", eng.handler.DELETE)
	router.OPTIONS("/v1/applysign", eng.handler.OPTIONS)
	router.PATCH("/v1/applysign/:__prefix", eng.handler.PATCH)
	router.OPTIONS("/v1/applysign/:__prefix", eng.handler.OPTIONS)

	// 活動簽到自定義人員 applysign
	// router.GET("/v1/applysign/user", eng.handler.GetApplysignUser)
	router.POST("/v1/applysign/user", eng.handler.POST)
	router.POST("/v1/applysign/users", eng.handler.POST)
	// router.PUT("/v1/applysign/user", eng.handler.PUT)
	router.DELETE("/v1/applysign/user", eng.handler.DELETE)
	// router.OPTIONS("/v1/applysign/user", eng.handler.OPTIONS)

	// 聊天區設定 wall
	router.GET("/v1/interact/wall/:__prefix", eng.handler.GetWall)
	router.POST("/v1/interact/wall/:__prefix", eng.handler.POST)
	router.PUT("/v1/interact/wall/:__prefix", eng.handler.PUT)
	router.PATCH("/v1/interact/wall/:__prefix", eng.handler.PATCH)
	router.DELETE("/v1/interact/wall/:__prefix", eng.handler.DELETE)
	router.OPTIONS("/v1/interact/wall/:__prefix", eng.handler.OPTIONS)
	router.POST("/v1/interact/wall/:__prefix/guest", eng.handler.POST)       // 提問嘉賓
	router.PUT("/v1/interact/wall/:__prefix/guest", eng.handler.PUT)         // 提問嘉賓
	router.DELETE("/v1/interact/wall/:__prefix/guest", eng.handler.DELETE)   // 提問嘉賓
	router.OPTIONS("/v1/interact/wall/:__prefix/guest", eng.handler.OPTIONS) // 提問嘉賓

	// 簽到展示 sign
	router.PATCH("/v1/interact/sign/:__prefix", eng.handler.PATCH)
	router.OPTIONS("/v1/interact/sign/:__prefix", eng.handler.OPTIONS)
	// 遊戲場次相關
	router.GET("/v1/interact/sign", eng.handler.GetSignGames)              // 遊戲json資料(多個遊戲場次)
	router.GET("/v1/interact/sign/:__prefix", eng.handler.GetSignGameInfo) // 遊戲json資料(單個遊戲場次)
	router.POST("/v1/interact/sign/:__prefix/form", eng.handler.POST)      // form-data
	router.PUT("/v1/interact/sign/:__prefix/form", eng.handler.PUT)        // form-data
	router.DELETE("/v1/interact/sign/:__prefix/form", eng.handler.DELETE)  // form-data
	router.OPTIONS("/v1/interact/sign/:__prefix/form", eng.handler.OPTIONS)

	// 獎品
	router.GET("/v1/interact/sign/:__prefix/prize", eng.handler.GetSignPrize)   // 獎品json資料
	router.POST("/v1/interact/sign/:__prefix/prize/form", eng.handler.POST)     // form-data
	router.PUT("/v1/interact/sign/:__prefix/prize/form", eng.handler.PUT)       // form-data
	router.DELETE("/v1/interact/sign/:__prefix/prize/form", eng.handler.DELETE) // form-data
	router.OPTIONS("/v1/interact/sign/:__prefix/prize/form", eng.handler.OPTIONS)

	// 投票選項
	router.GET("/v1/interact/sign/:__prefix/option", eng.handler.GetVoteOption) // 投票選項json資料
	router.POST("/v1/interact/sign/:__prefix/option", eng.handler.POST)         // form-data
	// 匹量插入
	router.PUT("/v1/interact/sign/:__prefix/option", eng.handler.PUT)       // form-data
	router.DELETE("/v1/interact/sign/:__prefix/option", eng.handler.DELETE) // form-data
	router.OPTIONS("/v1/interact/sign/:__prefix/option", eng.handler.OPTIONS)
	router.GET("/v1/interact/sign/:__prefix/option/reset", eng.handler.GetVoteOptionReset) // 投票分數紀錄重新計算
	// 投票特殊人員
	router.GET("/v1/interact/sign/:__prefix/special_officer", eng.handler.GetVoteSpecialOfficer) // 投票特殊人員json資料
	router.POST("/v1/interact/sign/:__prefix/special_officer", eng.handler.POST)                 // form-data
	router.POST("/v1/interact/sign/:__prefix/special_officers", eng.handler.POST)                // 匹量插入投票特殊人員
	router.PUT("/v1/interact/sign/:__prefix/special_officer", eng.handler.PUT)                   // form-data
	router.DELETE("/v1/interact/sign/:__prefix/special_officer", eng.handler.DELETE)             // form-data
	router.OPTIONS("/v1/interact/sign/:__prefix/special_officer", eng.handler.OPTIONS)
	// 投票選項名單
	router.GET("/v1/interact/sign/:__prefix/option_list", eng.handler.GetVoteOptionList) // 投票選項名單json資料
	router.POST("/v1/interact/sign/:__prefix/option_list", eng.handler.POST)             // form-data
	router.POST("/v1/interact/sign/:__prefix/option_lists", eng.handler.POST)            // 匹量插入
	router.PUT("/v1/interact/sign/:__prefix/option_list", eng.handler.PUT)               // form-data
	router.DELETE("/v1/interact/sign/:__prefix/option_list", eng.handler.DELETE)         // form-data
	router.OPTIONS("/v1/interact/sign/:__prefix/option_list", eng.handler.OPTIONS)

	// 抽獎遊戲 game
	router.GET("/v1/interact/game", eng.handler.GetGames)                 // 遊戲json資料(多個遊戲場次)
	router.GET("/v1/interact/game/:__prefix", eng.handler.GetGameInfo)    // 遊戲json資料(單個遊戲場次)
	router.POST("/v1/interact/game/:__prefix/form", eng.handler.POST)     // form-data
	router.PUT("/v1/interact/game/:__prefix/form", eng.handler.PUT)       // form-data
	router.DELETE("/v1/interact/game/:__prefix/form", eng.handler.DELETE) // form-data
	router.OPTIONS("/v1/interact/game/:__prefix/form", eng.handler.OPTIONS)
	// router.POST("/v1/interact/game/:__prefix/json", eng.handler.POST)     // json
	// router.PUT("/v1/interact/game/:__prefix/json", eng.handler.PUT)       // json
	// router.DELETE("/v1/interact/game/:__prefix/json", eng.handler.DELETE) // json
	// router.OPTIONS("/v1/interact/game/:__prefix/json", eng.handler.OPTIONS)

	// 遊戲獎品 game prize
	router.GET("/v1/game/lucky", eng.handler.GetLucky)                          // 判斷是否中獎，回傳中獎紀錄及抽獎紀錄(玩家端)
	router.GET("/v1/game/draw_numbers", eng.handler.GetDrawNumbers)             // 搖號抽獎中獎判斷，回傳所有中獎號碼及新的獎品資訊 GET API
	router.GET("/v1/game/tugofwar", eng.handler.GetTugofwar)                    // 拔河遊戲發獎判斷，回傳所有中獎資訊、獲勝隊伍、雙方隊長資訊、雙方分數資訊、雙方MVP資訊、雙方所有成員分數資訊GET API
	router.GET("/v1/game/vote", eng.handler.GetVote)                            // 投票發獎判斷，回傳所有中獎資訊GET API
	router.GET("/v1/interact/game/:__prefix/prize", eng.handler.GetPrize)       // 獎品json資料
	router.POST("/v1/interact/game/:__prefix/prize/form", eng.handler.POST)     // form-data
	router.PUT("/v1/interact/game/:__prefix/prize/form", eng.handler.PUT)       // form-data
	router.DELETE("/v1/interact/game/:__prefix/prize/form", eng.handler.DELETE) // form-data
	router.OPTIONS("/v1/interact/game/:__prefix/prize/form", eng.handler.OPTIONS)
	// router.POST("/v1/interact/game/:__prefix/prize/json", eng.handler.POST)     // json
	// router.PUT("/v1/interact/game/:__prefix/prize/json", eng.handler.PUT)       // json
	// router.DELETE("/v1/interact/game/:__prefix/prize/json", eng.handler.DELETE) // json
	// router.OPTIONS("/v1/interact/game/:__prefix/prize/json", eng.handler.OPTIONS)

	// 人員管理 staffmanage
	router.GET("/v1/staffmanage/:__prefix", eng.handler.GetStaffManage) // 人員json資料
	router.POST("/v1/staffmanage/:__prefix/form", eng.handler.POST)     // form-data
	router.PUT("/v1/staffmanage/:__prefix/form", eng.handler.PUT)       // form-data
	router.DELETE("/v1/staffmanage/:__prefix/form", eng.handler.DELETE) // form-data
	router.POST("/v1/staffmanage/:__prefix/export", eng.handler.EXPORT) // 匯出excel檔案
	router.OPTIONS("/v1/staffmanage/:__prefix/form", eng.handler.OPTIONS)
	router.OPTIONS("/v1/staffmanage/:__prefix/export", eng.handler.OPTIONS)
	// router.POST("/v1/staffmanage/:__prefix/json", eng.handler.POST)     // json
	// router.PUT("/v1/staffmanage/:__prefix/json", eng.handler.PUT)       // json
	// router.DELETE("/v1/staffmanage/:__prefix/json", eng.handler.DELETE) // json
	// router.OPTIONS("/v1/staffmanage/:__prefix/json", eng.handler.OPTIONS)

	// 匯入檔案 excel
	router.POST("/v1/import/excel", eng.handler.POST)

	// websocket功能
	// 遠端控制
	router.GET("/ws/v1/host/control", eng.handler.HostControlWebsocket)                // 主持端遙控螢幕
	router.GET("/ws/v1/host/control/session", eng.handler.HostControlSessionWebsocket) // 該活動所有可控制裝置的頁面資訊
	// -----舊-----
	// router.GET("/ws/v1/host/chatroom", eng.handler.HostChatroomWebsocket) // 主持端手機控制螢幕

	// 聊天紀錄
	router.GET("/ws/v1/chatroom/record", eng.handler.ChatroomRecordWebsocket)

	// 簽到牆
	router.GET("/ws/v1/host/signwall", eng.handler.SignWallWebsocket) // 即時活動簽到人數、人員資訊(主持人端)

	// 簽名牆
	router.GET("/ws/v1/host/signname", eng.handler.SignnameHostWebsocket)   // 即時回傳所有簽名牆資訊(主持人端)
	router.GET("/ws/v1/guest/signname", eng.handler.SignnameGuestWebsocket) // 即時新增簽名牆資訊(玩家端)

	// 提問牆
	router.GET("/ws/v1/host/question", eng.handler.HostQuestionWebsocket)   // 回傳提問資料(新舊排序、熱門排序)
	router.GET("/ws/v1/guest/question", eng.handler.GuestQuestionWebsocket) // 回傳提問資料、我的提問資料、按讚資料

	// 遊戲狀態(共用)
	// router.GET("/ws/v1/activity", eng.handler.ActivityWebsocket)                 // 即時活動資訊
	router.GET("/ws/v1/user", eng.handler.GameUserWebsocket)                     // 即時用戶資訊(玩家端)
	router.GET("/ws/v1/game", eng.handler.GameWebsocket)                         // 即時回傳遊戲狀態資訊(玩家端)
	router.GET("/ws/v1/game/team", eng.handler.GameTeamPlayerWebsocket)          // 即時回傳雙方隊伍所有玩家、隊長資訊(後端平台)
	router.GET("/ws/v1/game/user", eng.handler.GameUserWebsocket)                // 即時回傳用戶資訊(玩家端)
	router.GET("/ws/v1/game/staff", eng.handler.GameStaffWebsocket)              // 處理遊戲人員資訊(玩家端)，回傳用戶中獎紀錄、是否玩過該輪次遊戲
	router.GET("/ws/v1/game/people", eng.handler.PeopleWebsocket)                // 即時回傳遊戲人數、輪次資訊(主持人端)
	router.GET("/ws/v1/game/score", eng.handler.ScoreWebsocket)                  // 即時更新遊戲分數(玩家端)
	router.GET("/ws/v1/game/winning/staff", eng.handler.WinningStaffWebsocket)   // 即時中獎人員與獎品數量資訊(主持人端)
	router.GET("/ws/v1/game/status/host", eng.handler.HostGameStatusWebsocket)   // 主持人端遊戲狀態，更新資料表遊戲狀態
	router.GET("/ws/v1/game/status/guest", eng.handler.GuestGameStatusWebsocket) // 玩家端遊戲狀態，判斷資料表遊戲狀態變化
	router.GET("/ws/v1/game/host/QA", eng.handler.HostGameQAWebsocket)           // 主持端取得總答題人數、各選項答題人數資訊
	router.GET("/ws/v1/game/guest/QA", eng.handler.GuestGameQAWebsocket)         // 玩家端更新總答題人數、各選項答題人數資訊
	router.GET("/ws/v1/game/number/host", eng.handler.HostGameBingoWebsocket)    // 主持端即時更新賓果遊戲號碼
	router.GET("/ws/v1/game/number/guest", eng.handler.GuestGameBingoWebsocket)  // 玩家端即時回傳賓果遊戲所有號碼、更新玩家號碼排序
	router.GET("/ws/v1/game/reset", eng.handler.GameResetWebsocket)              // 即時判斷頁面是否需要重整
	router.GET("/ws/v1/activity/reset", eng.handler.ActivityResetWebsocket)      // 即時判斷頁面是否需要重整
	// 投票遊戲
	router.GET("/ws/v1/game/host/vote", eng.handler.HostGameVoteWebsocket)   // 主持端取得投票選項排名資訊
	router.GET("/ws/v1/game/guest/vote", eng.handler.GuestGameVoteWebsocket) // 玩家端接收前端投票紀錄資訊
	// 自定義場景
	router.GET("/ws/v1/customize_scene", eng.handler.CustomizeSceneWebsocket) // 自定義場景

	// 黑名單
	router.GET("/ws/v1/black/staff", eng.handler.BlackStaffWebsocket) // 即時回傳黑名單人員資訊

	// 操作日誌
	router.GET("/ws/v1/log", eng.handler.LogWebsocket) // 用戶操作日誌

	// 報名簽到人員
	router.GET("/ws/v1/applysign", eng.handler.ApplysignWebsocket) // 即時回傳報名簽到人員資訊(平台管理員端報名簽到頁面)

	// 測試用
	// router.GET("/v1/add", eng.handler.TestApi)
	// router.GET("/v1/update", eng.handler.TestApi2)
	// router.GET("/v1/delete", eng.handler.TestApi3)
	router.GET("/v1/find", eng.handler.TestApi4)
	// router.GET("/v1/test2", eng.handler.TestApi2)
	// router.GET("/ws/v1/test", eng.handler.TestWebsocket)

	// 待刪除------------------------------
	// 搖紅包
	// router.GET("/ws/v1/game/redpack/status/open", eng.handler.GameOpenWebsocket)       // 回傳搖紅包主持端展示的遊戲場次資訊(玩家端)
	// 套紅包
	// router.GET("/ws/v1/game/ropepack/status/open", eng.handler.GameOpenWebsocket)       // 回傳套紅包主持端展示的遊戲場次資訊(玩家端)

	// router.GET("/ws/v1/wall/message", eng.handler.MessageWebsocket)            // 訊息牆

	// middleware
	authChatroomRoute := r.Group("").Use(auth.DefaultChatroomInvoker(eng.dbConn, eng.redisConn, "chatroom_session").
		Middleware(), response.CORS)
	// 選擇活動管理或聊天室
	authChatroomRoute.GET("/select", eng.handler.ShowSelect)

	// 手機用戶端
	// 首頁
	authChatroomRoute.GET("/guest", eng.handler.ShowGuest)
	// 聊天室 chatroom
	authChatroomRoute.GET("/guest/chatroom", eng.handler.ShowGuestChatroom)
	// 個人中獎頁面 winning
	authChatroomRoute.GET("/guest/winning", eng.handler.ShowGuestWinning)
	// QRcode頁面 QRcode
	authChatroomRoute.GET("/guest/QRcode", eng.handler.ShowGuestQRcode)
	// 提問牆 question
	authChatroomRoute.GET("/guest/question", eng.handler.ShowQuestion)
	// 活動資訊 info
	authChatroomRoute.GET("/guest/info/*__prefix", eng.handler.ShowGuestInfo) // ###主頁面目前james沒使用到
	// authChatroomRoute.GET("/guest/info/introduce")
	// authChatroomRoute.GET("/guest/info/schedule")
	// authChatroomRoute.GET("/guest/info/guest")
	// authChatroomRoute.GET("/guest/info/material")

	// 遊戲互動 game
	// authChatroomRoute.GET("/guest/game/*__prefix", eng.handler.ShowGuestGame)
	// authChatroomRoute.GET("/guest/game", eng.handler.ShowGuestGame)                                 // 遊戲互動主頁面(// ###主頁面目前james沒使用到)
	authChatroomRoute.GET("/guest/game/:__prefix", eng.handler.ShowGuestGameInfo)                       // 遊戲資訊
	authChatroomRoute.GET("/guest/game/:__prefix/winning/staff", eng.handler.ShowGuestGameWinningStaff) // 獲獎人員

	// 主持端、手機端共用頁面
	authChatroomRoute.GET("/redpack/game", eng.handler.ShowGame)          // 遊戲頁面
	authChatroomRoute.GET("/ropepack/game", eng.handler.ShowGame)         // 遊戲頁面
	authChatroomRoute.GET("/whack_mole/game", eng.handler.ShowGame)       // 遊戲頁面
	authChatroomRoute.GET("/draw_numbers/game", eng.handler.ShowGame)     // 遊戲頁面
	authChatroomRoute.GET("/lottery/game", eng.handler.ShowGame)          // 遊戲頁面
	authChatroomRoute.GET("/monopoly/game", eng.handler.ShowGame)         // 遊戲頁面
	authChatroomRoute.GET("/QA/game", eng.handler.ShowGame)               // 遊戲頁面
	authChatroomRoute.GET("/tugofwar/game", eng.handler.ShowGame)         // 遊戲頁面
	authChatroomRoute.GET("/bingo/game", eng.handler.ShowGame)            // 遊戲頁面
	authChatroomRoute.GET("/3DGachaMachine/game", eng.handler.ShowGame)   // 遊戲頁面
	authChatroomRoute.GET("/3Ddraw_numbers/game", eng.handler.ShowGame)   // 遊戲頁面
	authChatroomRoute.GET("/vote/game", eng.handler.ShowGame)             // 遊戲頁面
	authChatroomRoute.GET("/signname/signwall", eng.handler.ShowSignWall) // 簽名牆

	// middleware
	authHiliveRoute := r.Group("").Use(
		auth.DefaultHiliveInvoker(eng.dbConn, eng.redisConn, "hilive_session").Middleware(),
		response.CORS)
	// 主持端
	// 簽到牆 signwall
	authHiliveRoute.GET("/general/signwall", eng.handler.ShowSignWall) // 一般簽到
	authHiliveRoute.GET("/threed/signwall", eng.handler.ShowSignWall)  // 立體簽到

	// 提問牆 question
	authHiliveRoute.GET("/host/question", eng.handler.ShowQuestion)
	// 聊天室 chatroom
	authHiliveRoute.GET("/host/chatroom", eng.handler.ShowHostChatroom)

	// 平台管理員
	// 靜態檔案(gzip)
	adminGzipRouter := r.Group("/admin").Use(response.CORS,
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/login", "/session"})))
	adminGzipRouter.Static("/assets", "./hilives/hilive/assets")
	// 遊戲靜態檔案(gzip)
	// gameGzipRouter := r.Group("").Use(response.CORS,
	// 	gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/login", "/session"})))
	// gameGzipRouter.Static("/vote/Build", "./hilives/hilive/assets/vote/Build")
	// gameGzipRouter.Static("/3Ddraw_numbers/Build", "./hilives/hilive/assets/3Ddraw_numbers/Build")

	// 沒有gzip
	adminRouter := r.Group("/admin").Use(response.CORS)
	adminRouter.Static("/uploads", config.STORE_PATH)

	// 加入log middleware fcuntion
	adminRouter = adminRouter.Use(response.CORS, logger.Logger(eng.dbConn))
	adminRouter.GET("/login", eng.handler.ShowLogin)    // 登入
	adminRouter.GET("/logout", eng.handler.Logout)      // 登出
	adminRouter.GET("/session", eng.handler.SetSession) // 設置session

	// middleware
	authAdminRoute := r.Group("/admin").Use(auth.DefaultHiliveInvoker(eng.dbConn, eng.redisConn, "hilive_session").
		Middleware(), response.CORS)

	// 管理員
	authAdminRoute.GET("/manager", eng.handler.ShowAdmin)               // 用戶
	authAdminRoute.GET("/manager/activity", eng.handler.ShowAdmin)      // 活動場次資訊
	authAdminRoute.GET("/manager/activity/game", eng.handler.ShowAdmin) // 遊戲場次資訊
	authAdminRoute.GET("/permission", eng.handler.ShowAdmin)            // 權限
	authAdminRoute.GET("/menu", eng.handler.ShowAdmin)                  // 菜單
	authAdminRoute.GET("/overview", eng.handler.ShowAdmin)              // 總覽
	authAdminRoute.GET("/log", eng.handler.ShowAdmin)                   // 日誌
	authAdminRoute.GET("/error_log", eng.handler.ShowAdmin)             // 錯誤日誌
	// 新增頁面
	authAdminRoute.GET("/manager/new", eng.handler.ShowNewAdmin)               // 用戶
	authAdminRoute.GET("/manager/activity/new", eng.handler.ShowNewAdmin)      // 活動
	authAdminRoute.GET("/manager/activity/game/new", eng.handler.ShowNewAdmin) // 遊戲
	authAdminRoute.GET("/permission/new", eng.handler.ShowNewAdmin)            // 權限
	authAdminRoute.GET("/menu/new", eng.handler.ShowNewAdmin)                  // 菜單
	authAdminRoute.GET("/overview/new", eng.handler.ShowNewAdmin)              // 總覽
	// 編輯頁面
	authAdminRoute.GET("/manager/edit", eng.handler.ShowEditAdmin)               // 用戶
	authAdminRoute.GET("/manager/activity/edit", eng.handler.ShowEditAdmin)      // 活動
	authAdminRoute.GET("/manager/activity/game/edit", eng.handler.ShowEditAdmin) // 遊戲
	authAdminRoute.GET("/permission/edit", eng.handler.ShowEditAdmin)            // 權限
	authAdminRoute.GET("/menu/edit", eng.handler.ShowEditAdmin)                  // 菜單
	authAdminRoute.GET("/overview/edit", eng.handler.ShowEditAdmin)              // 總覽

	authAdminRoute.GET("/user", eng.handler.Show)                         // 用戶
	authAdminRoute.GET("/activity", eng.handler.Show)                     // 活動
	authAdminRoute.GET("/activity/:__page", eng.handler.ShowActivity)     // 活動新增、編輯、快速設置、選擇項目頁面
	authAdminRoute.GET("/info/:__prefix", eng.handler.ShowAdminInfo)      // 活動訊息 info
	authAdminRoute.GET("/applysign/:__prefix", eng.handler.ShowApplysign) // 報名簽到 applysign
	authAdminRoute.GET("/interact/wall/:__prefix", eng.handler.ShowWall)  // 訊息區設定 wall

	// 簽到展示 sign
	authAdminRoute.GET("/interact/sign/:__prefix", eng.handler.ShowSign)                      // 設置頁面
	authAdminRoute.GET("/interact/sign/:__prefix/new", eng.handler.ShowNewInteractSign)       // 新增場次頁面
	authAdminRoute.GET("/interact/sign/:__prefix/edit", eng.handler.ShowEditInteractSign)     // 編輯場次頁面
	authAdminRoute.GET("/interact/sign/:__prefix/prize", eng.handler.ShowSignPrize)           // 獎品頁面
	authAdminRoute.GET("/interact/sign/:__prefix/option", eng.handler.ShowInteractSignOption) // 投票選項頁面

	// 抽獎遊戲 game
	authAdminRoute.GET("/interact/game/:__prefix", eng.handler.ShowInteractGame)          // 遊戲設置頁面
	authAdminRoute.GET("/interact/game/:__prefix/new", eng.handler.ShowNewInteractGame)   // 新增頁面
	authAdminRoute.GET("/interact/game/:__prefix/edit", eng.handler.ShowEditInteractGame) // 編輯頁面
	authAdminRoute.GET("/interact/game/:__prefix/prize", eng.handler.ShowPrize)           // 抽獎頁面
	authAdminRoute.GET("/interact/game/:__prefix/team", eng.handler.ShowInteractGameTeam) // 遊戲隊伍頁面
	// authAdminRoute.GET("/interact/game/:__prefix/pk", eng.handler.ShowInteractGamePK)     // PK頁面

	// 人員管理 staffmanage
	authAdminRoute.GET("/staffmanage/:__staffmanage", eng.handler.ShowStaffManage)

	// 自定義場景頁面
	authAdminRoute.GET("/customize_scene", eng.handler.ShowCustomizeScene)

	eng.Gin = r

	return eng
}

// func mainPage(c *fiber.Ctx) error {
// 	fmt.Println("有?")
// 	return c.Render("mainpage", fiber.Map{
// 		"people": 100,
// 		"People": 10,
// 	})
