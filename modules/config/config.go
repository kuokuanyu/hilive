package config

import (
	"hilive/modules/utils"
	"sync"
	"sync/atomic"
)

const (
	// redis名稱
	// LINE驗證
	// LINE_STATES_REDIS = "line_states" // 紀錄state資訊，SET
	// LINE_HOSTS_REDIS  = "line_hosts"  // 紀錄網域，LIST

	// redis pub/sub channel
	// pub/sub簽到
	CHANNEL_SIGN_STAFFS_2_REDIS = "channel:sign_staffs_2_" // 報名簽到人員，james用

	// pub/sub遊戲
	CHANNEL_GAME_REDIS              = "channel:game_"              // 遊戲資訊(會判斷game_attend.qa_people變化)，人數更新時會觸發
	CHANNEL_GUEST_GAME_STATUS_REDIS = "channel:guest_game_status_" // 遊戲狀態資訊(不判斷game_attend.qa_people變化)，玩家端遊戲狀態ws使用，人數更新時不會觸發(除了拔河)
	CHANNEL_GAME_BINGO_NUMBER_REDIS = "channel:bingo_number_"      // 紀錄抽過的號碼
	CHANNEL_QA_REDIS                = "channel:QA_"                // 題數資訊
	CHANNEL_GAME_TEAM_REDIS         = "channel:game_team_"         // 遊戲隊伍資訊
	CHANNEL_WINNING_STAFFS_REDIS    = "channel:winning_staffs_"    // 中獎人員
	CHANNEL_GAME_BINGO_USER_NUMBER  = "channel:bingo_user_number_" // 紀錄玩家的號碼排序
	CHANNEL_SCORES_REDIS            = "channel:score_"             // 玩家分數

	// pub/sub黑名單
	CHANNEL_BLACK_STAFFS_GAME_REDIS     = "channel:black_staffs_game_"     // 黑名單人員(遊戲)
	CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS = "channel:black_staffs_activity_" // 黑名單人員(活動)
	CHANNEL_BLACK_STAFFS_MESSAGE_REDIS  = "channel:black_staffs_message_"  // 黑名單人員(訊息)
	CHANNEL_BLACK_STAFFS_QUESTION_REDIS = "channel:black_staffs_question_" // 黑名單人員(提問)
	CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS = "channel:black_staffs_signname_" // 黑名單人員(簽名)

	// pub/sub編輯次數
	CHANNEL_SIGNNAME_EDIT_TIMES_REDIS = "channel:signname_edit_times_" // 簽名牆編輯次數
	CHANNEL_GENERAL_EDIT_TIMES_REDIS  = "channel:general_edit_times_"  // 一般簽到牆編輯次數
	CHANNEL_THREED_EDIT_TIMES_REDIS   = "channel:threed_edit_times_"   // 3D簽到牆編輯次數
	CHANNEL_GAME_EDIT_TIMES_REDIS     = "channel:game_edit_times_"     // 遊戲編輯次數

	// pub/sub遙控
	CHANNEL_HOST_CONTROL_REDIS         = "channel:host_control_"         // 主持端遙控資訊
	CHANNEL_HOST_CONTROL_CHANNEL_REDIS = "channel:host_control_channel_" // 主持端所有可遙控的channel

	// redis鎖
	LUCKY_LOCK_REDIS       = "lucky_lock_"       // 抽獎
	ADD_STAFF_LOCK_REDIS   = "add_staff_lock_"   // 加入遊戲人員
	BINGO_LOCK_REDIS       = "bingo_lock_"       // 賓果
	LINE_USERS_LOCK_REDIS  = "line_users_lock_"  // LINE用戶
	VOTE_OPTION_LOCK_REDIS = "vote_option_lock_" // 投票選項
	// 目前未用到的redis鎖
	// APPLYSIGN_LOCK_REDIS          = "applysign_lock_"          // 報名簽到處理
	// MAIL_LOCK_REDIS               = "mail_lock_"               // 郵件
	// MESSAGE_LOCK_REDIS            = "message_lock_"            // 簡訊
	// LUCKY_DRAW_NUMBERS_LOCK_REDIS = "lucky_draw_numbers_lock_" // 搖號抽獎
	// PRIZE_LOCK_REDIS              = "prize_lock_"              // 獎品

	// 用戶資訊
	HILIVE_USERS_REDIS = "hilive_users_" // 平台管理員資訊，HASH
	AUTH_USERS_REDIS   = "auth_users_"   // 參加用戶資訊，HASH

	// 活動
	ACTIVITY_REDIS        = "activity_"        // 活動資訊，HASH
	ACTIVITY_NUMBER_REDIS = "activity_number_" // 活動number，STRING

	// 簽到人員
	// SIGN_STAFFS_1_REDIS = "sign_staffs_1_" // 簽到人員1(更新資料時，不修改redis裡的資料，威翔用)，LIST
	SIGN_STAFFS_2_REDIS = "sign_staffs_2_" // 簽到人員2(更新資料時，修改redis裡的資料，james用)，SET

	// 簽名牆
	SIGNNAME_REDIS = "signname_" // 簽名牆設置資料，HASH
	// SIGNNAME_DATAS_REDIS         = "signname_datas_"         // 簽名牆資料，HASH
	// SIGNNAME_ORDER_BY_TIME_REDIS = "signname_order_by_time_" // 簽名牆資料(時間)，LIST

	// 搖控
	HOST_CONTROL_REDIS         = "host_control_"         // 主持端遙控資訊，HASH，host_control_id_channel_n
	HOST_CONTROL_CHANNEL_REDIS = "host_control_channel_" // 主持端所有可遙控的channel，HASH

	// 聊天室
	// CHATROOM_REDIS               = "chatroom_"               // 聊天紀錄資訊，HASH
	// CHATROOM_ORDER_BY_TIME_REDIS = "chatroom_order_by_time_" // 聊天紀錄資訊(時間)，LIST

	// 提問區
	// QUESTION_REDIS               = "question_"               // 提問資訊(所有提問資訊)，HASH
	// QUESTION_ORDER_BY_TIME_REDIS = "question_order_by_time_" // 提問資訊(時間)，LIST
	// QUESTION_ORDER_BY_LIKES_REDIS    = "question_order_by_likes_"    // 提問資訊(讚數)，ZSET
	// QUESTION_USER_LIKE_RECORDS_REDIS = "question_user_like_records_" // 提問資訊(用戶按讚紀錄)，HASH

	// 遊戲
	GAME_REDIS                            = "game_"                            // 遊戲資訊，HASH
	GAME_TYPE_REDIS                       = "game_type_"                       // 遊戲種類資訊，STRING
	GAME_PRIZES_AMOUNT_REDIS              = "game_prizes_amount_"              // 遊戲獎品數量，HASH
	PRIZE_REDIS                           = "prize_"                           // 獎品資訊，HASH
	USER_GAME_RECORDS_REDIS               = "user_game_records_"               // 玩家在該活動下的遊戲紀錄(中獎.未中獎)，HASH
	GAME_VOTE_RECORDS_REDIS               = "game_vote_records_"               // 玩家投票紀錄，HASH
	GAME_ATTEND_REDIS                     = "game_attend_"                     // 遊戲人數資訊，SET
	GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS  = "game_tugofwar_left_team_attend_"  // 拔河遊戲左隊人數資訊，SET
	GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS = "game_tugofwar_right_team_attend_" // 拔河遊戲右隊人數資訊，SET

	// -----舊-----
	// HOST_CHATROOM_REDIS = "host_chatroom_" // 遙控主持端資訊

	// 黑名單人員
	// BLACK_STAFFS_REDIS          = ""
	BLACK_STAFFS_ACTIVITY_REDIS = "black_staffs_activity_" // 黑名單人員(活動)，SET
	BLACK_STAFFS_MESSAGE_REDIS  = "black_staffs_message_"  // 黑名單人員(訊息)，SET
	BLACK_STAFFS_QUESTION_REDIS = "black_staffs_question_" // 黑名單人員(提問)，SET
	BLACK_STAFFS_SIGNNAME_REDIS = "black_staffs_signname_" // 黑名單人員(簽名)，SET
	BLACK_STAFFS_GAME_REDIS     = "black_staffs_game_"     // 黑名單人員(遊戲)，SET

	// 分數紀錄相關
	SCORES_REDIS                      = "score_"                       // 玩家分數，ZSET
	SCORES_2_REDIS                    = "score_2_"                     // 玩家第二分數，ZSET
	CORRECT_REDIS                     = "correct_"                     // 玩家答對題數，ZSET
	WINNING_STAFFS_REDIS              = "winning_staffs_"              // 中獎人員，LIST、SET
	NO_WINNING_STAFFS_REDIS           = "no_winning_staffs_"           // 未中獎人員(紅包遊戲用)，LIST
	DRAW_NUMBERS_WINNING_STAFFS_REDIS = "draw_numbers_winning_staffs_" // 活動下所有搖號抽獎場次中獎人員資料，SET
	QA_REDIS                          = "QA_"                          // 題數資訊，HASH
	QA_RECORD_REDIS                   = "QA_record_"                   // 答題紀錄資訊，HASH
	GAME_TEAM_REDIS                   = "game_team_"                   // 遊戲隊伍資訊(隊長、leader、win)，HASH
	GAME_BINGO_NUMBER_REDIS           = "bingo_number_"                // 紀錄抽過的號碼，LIST
	GAME_BINGO_USER_REDIS             = "bingo_user_"                  // 賓果中獎人員，ZSET
	GAME_BINGO_USER_NUMBER            = "bingo_user_number_"           // 紀錄玩家的號碼排序，HASH
	GAME_BINGO_USER_GOING_BINGO       = "bingo_user_going_bingo_"      // 紀錄玩家是否即將賓果，HASH
	VOTE_SPECIAL_OFFICER_REDIS        = "vote_special_officer_"        // 投票遊戲特殊人員，HASH
	VOTE_AVATAR_REDIS                 = "vote_avatar_"                 // 投票人員頭像，LIST

	// VOTE_OPTION_REDIS  = "vote_option_"    // 投票選項資料，HASH

	// cookie加密秘鑰
	HiliveCookieSecret   = "hilive_session"
	ChatroomCookieSecret = "chatroom_session"

	// token加密秘鑰
	TokenSecret = "token_csrf"

	// 報名簽到驗證
	HTTPS_ACTIVITY_ISEXIST_URL = "https://%s/activity/isexist?activity_id=%s&game_id=%s"                   // 一般url，報名簽到qrcode
	HTTPS_AUTH_REDIRECT_URL    = "https://%s/auth/redirect?action=%s&device=%s&%s=%s&game_id=%s"           // 一般url，auth redirect url
	AUTH_REDIRECT_URL          = "/auth/redirect?action=%s&device=%s&activity_id=%s&user_id=%s&game_id=%s" // 一般url，自定義匯入簽到人員登入
	HTTPS_AUTH_CALLBACK_URL    = "https://%s/v1/auth/callback"                                             // callback url

	// LINE選單
	HTTPS_LINE_RICHMENU_URL  = "https://%s/activity/%s?user_id=%s" // 選單頁面(要網域)
	LINE_RICHMENU_URL        = "/activity/%s?user_id=%s"           // 選單頁面(不要網域)
	LINE_RICHMENU_SEARCH_URL = "/activity/search/%s?user_id=%s"    // 活動查詢選單頁面

	// 測試區
	// 網域配置
	HTTPS_API_URL        = "https://apidev.hilives.net"
	API_URL              = "apidev.hilives.net"
	HTTP_HILIVES_NET_URL = "https://dev.hilives.net"
	HILIVES_NET_URL      = "dev.hilives.net"

	// LINE LOGIN、MESSAGE相關ID、SECRET
	CHANNEL_ID     = "1656920628"
	CHANNEL_SECRET = "19f076025659a5e50bccd931da4641d8"
	CHATBOT_SECRET = "88d68c3fa9c25e99e6c777a791043a48"
	CHATBOT_TOKEN  = "k6z4Z2jvoYdc75I801ZZUc+dicI+oZe+WezNBh1Wrk5E2+gUWiUeT2I7XEO44OycU8jnGwNsWsq2pHVb7EDndhurcltSlhLTfWtfDkiHLyL32g+2QoQwdtFfVRv/ar5MNwiNPhWEBZBArNLSfGY0FwdB04t89/1O/w1cDnyilFU="

	// facebook驗證參數
	FACEBBOK_ID           = "814026164217592"
	FACEBOOK_SECRET       = "e06412863dbb12daec68d54c0da47734"
	FACEBOOK_REDIRECT_URL = "https://dev.hilives.net/v1/auth/callback"

	// gmail驗證參數
	GMAIL_ID           = "804432621213-vcs0h63r2uslcn9jeuutgsjokorvjl3k.apps.googleusercontent.com"
	GMAIL_SECRET       = "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"
	GMAIL_REDIRECT_URL = "https://dev.hilives.net/v1/auth/callback"
	GMAIL              = "sales@hilives.net"
	GMAIL_PASSWORD     = "kgdl nacn ltkd ycir"
	// gmail SMTP 服務器配置
	GMAIL_SMTP_HOST = "smtp.gmail.com"
	GMAIL_SMTP_PORT = "587"

	// 活動是否存在(不會執行middleware function)
	HILIVES_ACTIVITY_ISEXIST_LIFF_URL = "https://liff.line.me/1656920628-zJOEMMRl?activity_id=%s"
	// 報名簽到(不會執行middleware function)
	HILIVES_APPLYSIGN_URL_LIFF_URL = "https://liff.line.me/1656920628-jwWm55v7?activity_id=%s&user_id=%s"
	// 搖紅包
	HILIVES_REDPACK_GAME_LIFF_URL = "https://liff.line.me/1656920628-ZXAlXXa7?activity_id=%s&game_id="
	// 套紅包
	HILIVES_ROPEPACK_GAME_LIFF_URL = "https://liff.line.me/1656920628-YW3mWWxl?activity_id=%s&game_id="
	// 打地鼠
	HILIVES_WHACK_MOLE_GAME_LIFF_URL = "https://liff.line.me/1656920628-Yoe0EEpy?activity_id=%s&game_id="
	// 遊戲抽獎
	HILIVES_LOTTERY_GAME_LIFF_URL = "https://liff.line.me/1656920628-n8RGYY8b?activity_id=%s&game_id="
	// 超級大富翁
	HILIVES_MONOPOLY_GAME_LIFF_URL = "https://liff.line.me/1656920628-GXrJ22qw?activity_id=%s&game_id="
	// 提問區
	HILIVES_QUESTION_LIFF_URL = "https://liff.line.me/1656920628-WOvQqqJ1?activity_id=%s"
	// 快問快答
	HILIVES_QA_GAME_LIFF_URL = "https://liff.line.me/1656920628-Rq90oo8G?activity_id=%s&game_id="
	// 拔河遊戲
	HILIVES_TUGOFWAR_GAME_LIFF_URL = "https://liff.line.me/1656920628-dP3k2250?activity_id=%s&game_id="
	// 賓果遊戲
	HILIVES_BINGO_GAME_LIFF_URL = "https://liff.line.me/1656920628-5rAD11ME?activity_id=%s&game_id="
	// QRcode(手機玩家端使用)
	HILIVES_QRCODE_LIFF_URL = "https://liff.line.me/1656920628-yYgRWWd5?activity_id=%s"
	// LINE richmenu
	HILIVES_LINE_RICHMENU_LIFF_URL = "https://liff.line.me/1656920628-5ZWdGGBo?action=%s"
	// 簽名牆
	HILIVES_SIGNNAME_LIFF_URL = "https://liff.line.me/1656920628-x7dzooDR?activity_id=%s"
	// 扭蛋機遊戲
	HILIVES_GACHAMACHINE_GAME_LIFF_URL = "https://liff.line.me/1656920628-08N9rr6Q?activity_id=%s&game_id="
	// 投票遊戲
	HILIVES_VOTE_GAME_LIFF_URL = "https://liff.line.me/1656920628-Ra4MbbBo?activity_id=%s&game_id="

	// Redis相關配置
	REDIS_ENGINE = "hilives"
	REDIS_HOST   = "10.12.225.3"
	REDIS_EXPIRE = "86400" //86400
	REDIS_PORT   = "6379"

	// Mysql相關配置
	MYSQL_ENGINE   = "hilives"
	MYSQL_HOST     = "10.12.224.11" // 35.221.150.128
	MYSQL_PORT     = "3306"
	MYSQL_USER     = "hilives"
	MYSQL_PASSWORD = "Cco@53383499"
	MYSQL_NAME     = "hilive_dev"
	// SHOW VARIABLES LIKE 'max_connections'; 目前值為4030
	MYSQL_MAXOPENCON = 4000
	MYSQL_MAXIDLECON = 2000
	MYSQL_DRIVER     = "mysql"

	// mongodb相關配置
	MONGO_ENGINE   = "hilives"
	MONGO_HOST     = "hilives002.jadard7.mongodb.net"
	MONGO_PORT     = "27017"
	MONGO_USER     = "hilives_dev"
	MONGO_PASSWORD = "Cco%4053383499"
	MONGO_NAME     = "hilives_dev"

	// 測試區

	// 手機驗證碼
	ACCOUNT_SID = "ACc80317eb45de47031773d4afdad4252e"
	AUTH_TOKEN  = "6d7349e4bb734eb142916bf4832374a1"
	SERVICE_SID = "VA9d864fb293777afe046b5d8e26ad862e"
	PHONE       = "+18596966103"

	// 檔案引擎配置
	FILE_ENGINE = "hilives"

	// 圖片相關
	UPLOAD_RADIO_URL  = "/admin/uploads/system/radio/"
	UPLOAD_SYSTEM_URL = "/admin/uploads/system/"
	UPLOAD_URL        = "/admin/uploads/"
	DEFAULT_FALG      = "__default_flag"

	// 註冊、登入、忘記密碼、登出、首頁
	INDEX_URL               = "/"                                // 首頁URL
	LOGIN_URL               = "/admin/login"                     // 登入頁面URL
	LOGIN_API_URL           = HTTPS_API_URL + "/v1/login"        // 登入API
	APPLYSIGN_LOGIN_API_URL = HTTPS_API_URL + "/applysign/login" // 自定義匯入簽到人員登入API
	REGISTER_API_URL        = HTTPS_API_URL + "/v1/register"     // 註冊API
	RETRIEVE_API_URL        = HTTPS_API_URL + "/v1/retrieve"     // 忘記密碼API
	LOGOUT_URL              = "/admin/logout"                    // 登出頁面URL

	// 手機驗證
	VERIFICATION_API_URL       = HTTPS_API_URL + "/v1/verification"       // 手機驗證API
	VERIFICATION_CHECK_API_URL = HTTPS_API_URL + "/v1/verification/check" // 手機驗證碼檢查API

	// session相關
	// HTTPS_CHATROOM_SESSION_URL = "%s/admin/session?session_name=chatroom&activity_id=%s&user_id=%s" // session頁面(要網域)
	CHATROOM_SESSION_URL = "/admin/session?session_name=chatroom&activity_id=%s&user_id=%s&game_id=%s" // session頁面(不須網域)
	HILIVE_SESSION_URL   = "/admin/session?session_name=hilive&user_id=%s"                                     // session頁面(不須網域)

	// 補齊資料
	AUTH_USER_API_URL   = HTTPS_API_URL + "/v1/auth/user"                             // 補齊資料API
	HTTPS_AUTH_USER_URL = "https://%s/auth/user?activity_id=%s&user_id=%s&game_id=%s" // 補齊資料頁面

	// 報名簽到
	HTTPS_APPLYSIGN_URL        = "https://%s/applysign?activity_id=%s&user_id=%s&game_id=%s" // 報名簽到(要網域)
	HTTPS_APPLYSIGN_QRCODE_URL = "https://%s/applysign?qrcode=%s"                            // 報名簽到(要網域，郵件中的qrcode連結)
	APPLYSIGN_URL              = "/applysign?activity_id=%s&user_id=%s&game_id=%s"           // 報名簽到(不須網域)

	// 管理員綁定選擇頁面
	SELECT_URL = "/select?activity_id=%s"

	// 主持端
	HOST_CHATROOM_URL          = "/host/chatroom?activity_id="                  // 主持端聊天室頁面URL
	HOST_GAME_URL              = "/host/chatroom?activity_id=%s&login=%s&game=" // 主持端聊天室頁面URL，直接導向某個遊戲
	HOST_GENERAL_SIGNWALL_URL  = "/general/signwall?activity_id=%s&login=%s"    // 一般簽到牆URL
	HOST_THREED_SIGNWALL_URL   = "/threed/signwall?activity_id=%s&login=%s"     // 立體簽到牆URL
	HOST_SIGNNAME_SIGNWALL_URL = "/signname/signwall?activity_id=%s&login=%s"   // 簽名牆URL
	HOST_QUESTION_URL          = "/host/question?activity_id=%s&login=%s"       // 提問牆URL

	// 手機端
	GUEST_WINNING_URL = "/guest/winning?activity_id=%s" // 手機端用戶中獎紀錄URL

	GUEST_INTRODUCE_URL = "/guest/info/introduce?activity_id=%s" // 介紹
	GUEST_SCHEDULE_URL  = "/guest/info/schedule?activity_id=%s"  // 行程
	GUEST_GUEST_URL     = "/guest/info/guest?activity_id=%s"     // 嘉賓
	GUEST_MATERIAL_URL  = "/guest/info/material?activity_id=%s"  // 資料

	GUEST_LOTTERY_GAME_URL      = "/guest/game/lottery?activity_id=%s"        // 遊戲抽獎場次資訊
	GUEST_REDPACK_GAME_URL      = "/guest/game/redpack?activity_id=%s"        // 搖紅包場次資訊
	GUEST_ROPEPACK_GAME_URL     = "/guest/game/ropepack?activity_id=%s"       // 套紅包場次資訊
	GUEST_WHACKMOLE_GAME_URL    = "/guest/game/whack_mole?activity_id=%s"     // 敲敲樂場次資訊
	GUEST_MONOPOLY_GAME_URL     = "/guest/game/monopoly?activity_id=%s"       // 超級大富翁場次資訊
	GUEST_QA_GAME_URL           = "/guest/game/QA?activity_id=%s"             // 快問快答場次資訊
	GUEST_TUGOFWAR_GAME_URL     = "/guest/game/tugofwar?activity_id=%s"       // 拔河遊戲場次資訊
	GUEST_BINGO_GAME_URL        = "/guest/game/bingo?activity_id=%s"          // 賓果遊戲場次資訊
	GUEST_GACHAMACHINE_GAME_URL = "/guest/game/3DGachaMachine?activity_id=%s" // 扭蛋機遊戲場次資訊
	GUEST_VOTE_GAME_URL         = "/guest/game/vote?activity_id=%s"           // 投票遊戲場次資訊

	GUEST_QUESTION_URL     = "/guest/question?activity_id=%s"                       // 提問牆
	GUEST_SIGNNAME_URL     = "/signname/signwall?activity_id=%s"                    // 簽名牆
	GUEST_GAME_WINNING_URL = "/guest/game/%s/winning/staff?activity_id=%s&game_id=" // 中獎資訊

	// GUEST_LOTTERY_WINNING_URL   = "/guest/game/lottery/winning/staff?activity_id=%s&game_id="    // 遊戲抽獎中獎資訊
	// GUEST_REDPACK_WINNING_URL   = "/guest/game/redpack/winning/staff?activity_id=%s&game_id="    // 搖紅包中獎資訊
	// GUEST_ROPEPACK_WINNING_URL  = "/guest/game/ropepack/winning/staff?activity_id=%s&game_id="   // 套紅包中獎資訊
	// GUEST_WHACKMOLE_WINNING_URL = "/guest/game/whack_mole/winning/staff?activity_id=%s&game_id=" // 敲敲樂中獎資訊

	// 聊天室
	GUEST_URL               = "/guest?activity_id=%s"               // 手機端首頁URL
	GUEST_CHATROOM_URL      = "/guest/chatroom?activity_id=%s"      // 手機端聊天室頁面URL
	CHATROOM_RECORD_API_URL = HTTPS_API_URL + "/v1/chatroom/record" // 聊天室紀錄API
	GUEST_QRCODE_URL        = "/guest/QRcode?activity_id=%s"        // 手機端QRcode URL

	// 提問區
	QUESTION_RECORD_API_URL = HTTPS_API_URL + "/v1/question/record" // 提問紀錄API

	// 遊戲頁面(無網域)，主持端使用
	GAME_URL                     = "/%s/game?activity_id=%s&game_id="                              // 遊戲頁面
	GAME_ROLE_URL                = "/%s/game?activity_id=%s&role=%s&game_id="                      // 遊戲頁面(包含角色參數)
	REDPACK_GAME_URL             = "/redpack/game?activity_id=%s&login=%s&game_id="                // 搖紅包
	ROPEPACK_GAME_URL            = "/ropepack/game?activity_id=%s&login=%s&game_id="               // 套紅包
	WHACK_MOLE_GAME_URL          = "/whack_mole/game?activity_id=%s&login=%s&game_id="             // 敲敲樂
	DRAW_NUMBERS_GAME_URL        = "/draw_numbers/game?activity_id=%s&login=%s&game_id="           // 搖號抽獎
	THREED_DRAW_NUMBERS_GAME_URL = "/3Ddraw_numbers/game?activity_id=%s&login=%s&game_id="         // 3D搖號抽獎
	LOTTERY_GAME_URL             = "/lottery/game?activity_id=%s&login=%s&game_id="                // 遊戲抽獎
	MONOPOLY_GAME_URL            = "/monopoly/game?activity_id=%s&login=%s&game_id="               // 超級大富翁
	QA_GAME_URL                  = "/QA/game?activity_id=%s&login=%s&game_id="                     // 快問快答
	TUGOFWAR_GAME_URL            = "/tugofwar/game?activity_id=%s&login=%s&game_id="               // 拔河遊戲
	BINGO_GAME_URL               = "/bingo/game?activity_id=%s&login=%s&game_id="                  // 賓果遊戲
	GACHAMACHINE_GAME_URL        = "/3DGachaMachine/game?activity_id=%s&login=%s&role=%s&game_id=" // 扭蛋機遊戲
	VOTE_GAME_URL                = "/vote/game?activity_id=%s&login=%s&role=%s&game_id="           // 投票遊戲

	// 遊戲頁面(有網域)，qrcode使用
	HTTPS_REDPACK_GAME_URL      = "https://%s/redpack/game?activity_id=%s&game_id="                // 搖紅包
	HTTPS_ROPEPACK_GAME_URL     = "https://%s/ropepack/game?activity_id=%s&game_id="               // 套紅包
	HTTPS_WHACK_MOLE_GAME_URL   = "https://%s/whack_mole/game?activity_id=%s&game_id="             // 敲敲樂
	HTTPS_DRAW_NUMBERS_GAME_URL = "https://%s/draw_numbers/game?activity_id=%s&game_id="           // 搖號抽獎
	HTTPS_LOTTERY_GAME_URL      = "https://%s/lottery/game?activity_id=%s&game_id="                // 遊戲抽獎
	HTTPS_MONOPOLY_GAME_URL     = "https://%s/monopoly/game?activity_id=%s&game_id="               // 超級大富翁
	HTTPS_QA_GAME_URL           = "https://%s/QA/game?activity_id=%s&game_id="                     // 快問快答
	HTTPS_TUGOFWAR_GAME_URL     = "https://%s/tugofwar/game?activity_id=%s&game_id="               // 拔河遊戲
	HTTPS_BINGO_GAME_URL        = "https://%s/bingo/game?activity_id=%s&game_id="                  // 賓果遊戲
	HTTPS_GACHAMACHINE_GAME_URL = "https://%s/3DGachaMachine/game?activity_id=%s&role=%s&game_id=" // 扭蛋機遊戲遊戲
	HTTPS_VOTE_GAME_URL         = "https://%s/vote/game?activity_id=%s&role=%s&game_id="           // 投票遊戲遊戲

	HTTPS_SIGNNAME_URL = "https://%s/signname/signwall?activity_id=%s" // 簽名牆URL
	HTTPS_QUESTION_URL = "https://%s/guest/question?activity_id=%s"    // 提問牆URL

	// ------------------平台---------------

	// 管理員
	ADMIN_URL              = "/admin/%s"                    // 資訊頁面URL
	ADMIN_NEW_URL          = "/admin/%s/new"                // 新增頁面URL
	ADMIN_EDIT_URL         = "/admin/%s/edit?id="           // 編輯頁面URL
	ADMIN_MANAGER_EDIT_URL = "/admin/%s/edit?user_id="      // 用戶編輯頁面URL
	ADMIN_API_URL          = HTTPS_API_URL + "/v1/admin/%s" // 管理API

	// 管理員頁面裡的活動場次頁面
	ADMIN_ACTIVITY_URL      = "/admin/manager/activity?user_id=%s"                     // 活動資訊頁面URL
	ADMIN_ACTIVITY_NEW_URL  = "/admin/manager/activity/new?user_id=%s"                 // 新增活動頁面URL
	ADMIN_ACTIVITY_EDIT_URL = "/admin/manager/activity/edit?user_id=%s&activity_id=%s" // 編輯活動頁面URL
	// 管理員頁面裡的遊戲場次頁面
	ADMIN_GAME_TYPE_URL = "/admin/manager/activity/game?user_id=%s&activity_id=%s"                         // 遊戲種類頁面URL
	ADMIN_GAME_URL      = "/admin/manager/activity/game?user_id=%s&activity_id=%s&game=%s"                 // 遊戲資訊頁面URL
	ADMIN_GAME_NEW_URL  = "/admin/manager/activity/game/new?user_id=%s&activity_id=%s&game=%s"             // 新增遊戲頁面URL
	ADMIN_GAME_EDIT_URL = "/admin/manager/activity/game/edit?user_id=%s&activity_id=%s&game=%s&game_id=%s" // 編輯遊戲頁面URL

	// 用戶
	USER_URL     = "/admin/user?header=%s"    // 用戶頁面URL
	USER_API_URL = HTTPS_API_URL + "/v1/user" // 用戶API

	// 活動
	HTTPS_ACTIVITY_URL              = "https://%s/admin/activity?header=%s"      // 活動頁面URL(網域)
	ACTIVITY_URL                    = "/admin/activity?header=%s"                // 活動頁面URL(不須網域)
	ACTIVITY_NEW_URL                = "/admin/activity/new"                      // 新增活動頁面URL
	ACTIVITY_EDIT_URL               = "/admin/activity/edit?activity_id=%s"      // 編輯活動頁面URL(有activity_id)
	ACTIVITY_EDIT_NO_ACTIVITYID_URL = "/admin/activity/edit?activity_id="        // 編輯活動頁面URL(沒有activity_id)
	ACTIVITY_API_URL                = HTTPS_API_URL + "/v1/activity"             // 活動API
	ACTIVITY_QUICK_START_URL        = "/admin/activity/quick_start?activity_id=" // 快速設置頁面URL
	ACTIVITY_SELECT_URL             = "/admin/activity/select?activity_id=%s"    // 選擇項目頁面URL

	// 報名簽到
	USER_APPLYSIGN_API_URL = HTTPS_API_URL + "/v1/applysign"    // 用戶報名簽到API
	APPLYSIGN_API_URL      = HTTPS_API_URL + "/v1/applysign/%s" // 報名簽到基本設置API
	// APPLY_API_URL     = HTTPS_API_URL + "/v1/applysign/apply"     // 報名API
	// SIGN_API_URL      = HTTPS_API_URL + "/v1/applysign/sign"      // 簽到API
	// CUSTOMIZE_API_URL = HTTPS_API_URL + "/v1/applysign/customize" // 自訂義API

	// 活動資訊
	INFO_URL     = "/admin/info/%s?activity_id=%s"                  // 活動資訊URL
	INFO_API_URL = HTTPS_API_URL + "/v1/info/%s"                    // 活動資訊API
	OVERVIEW_URL = "/admin/info/overview?sidebar=true&activity_id=" // 活動總覽頁面URL

	// 聊天區設定
	INTERACT_WALL_URL      = "/admin/interact/wall/%s?activity_id=%s"           // 聊天區設定頁面URL
	INTERACT_WALL_API_URL  = HTTPS_API_URL + "/v1/interact/wall/%s"             // 聊天區設定API
	QUESTION_GUEST_API_URL = HTTPS_API_URL + "/v1/interact/wall/question/guest" // 提問嘉賓API

	// 簽到展示
	INTERACT_SIGN_API_URL            = HTTPS_API_URL + "/v1/interact/sign/%s"                   // 簽到展示API
	INTERACT_SIGN_URL                = "/admin/interact/sign/%s?activity_id=%s"                 // 簽到展示頁面API
	INTERACT_SIGN_API_URL_FORM       = HTTPS_API_URL + "/v1/interact/sign/%s/form"              // 遊戲API
	INTERACT_SIGN_PRIZE_URL          = "/admin/interact/sign/%s/prize?activity_id=%s&game_id="  // 獎品頁面URL
	INTERACT_SIGN_NEW_URL            = "/admin/interact/sign/%s/new?activity_id=%s"             // 新增遊戲頁面URL
	INTERACT_SIGN_EDIT_URL           = "/admin/interact/sign/%s/edit?activity_id=%s&game_id="   // 編輯遊戲頁面URL
	INTERACT_SIGN_PRIZE_API_URL_FORM = HTTPS_API_URL + "/v1/interact/sign/%s/prize/form"        // 獎品API
	INTERACT_SIGN_VOTE_OPTION_URL    = "/admin/interact/sign/%s/option?activity_id=%s&game_id=" // 投票選項頁面URL

	// 互動遊戲
	INTERACT_GAME_URL                = "/admin/interact/game/%s?activity_id=%s"                // 遊戲頁面URL
	INTERACT_GAME_TEAM_URL           = "/admin/interact/game/%s/team?activity_id=%s&game_id="  // 隊伍頁面URL
	INTERACT_PRIZE_URL               = "/admin/interact/game/%s/prize?activity_id=%s&game_id=" // 獎品頁面URL
	INTERACT_GAME_NEW_URL            = "/admin/interact/game/%s/new?activity_id=%s"            // 新增遊戲頁面URL
	INTERACT_GAME_EDIT_URL           = "/admin/interact/game/%s/edit?activity_id=%s&game_id="  // 編輯遊戲頁面URL
	INTERACT_GAME_PRIZE_URL          = "/admin/interact/game/%s/prize?activity_id=%s&game_id=" // 獎品頁面URL
	INTERACT_GAME_API_URL_FORM       = HTTPS_API_URL + "/v1/interact/game/%s/form"             // 遊戲API
	INTERACT_PRIZE_API_URL_FORM      = HTTPS_API_URL + "/v1/interact/game/%s/prize/form"       // 獎品API
	INTERACT_GAME_RESET_API_URL_FORM = HTTPS_API_URL + "/v1/interact/game/reset/form"          // 遊戲重置API

	// 人員管理
	// STAFFMANAGE_NO_GAMEID_URL  = "/admin/staffmanage/%s/%s/%s?activity_id=%s&game_id="   // 人員管理頁面URL(沒有game_id)
	STAFFMANAGE_URL          = "/admin/staffmanage/%s?activity_id=%s"    // 人員管理頁面URL
	STAFFMANAGE_API_URL_FORM = HTTPS_API_URL + "/v1/staffmanage/%s/form" // 人員管理API(form)
	// STAFFMANAGE_API_URL_JSON   = HTTPS_API_URL + "/v1/staffmanage/%s/json"   // 中獎人員API(json)
	STAFFMANAGE_API_URL_EXPORT = HTTPS_API_URL + "/v1/staffmanage/%s/export" // 匯出人員API
	// EXPORT_QA_RECORD_API_URL         = HTTPS_API_URL + "/v1/QA/record/export"            // 匯出答題紀錄API
	// WINNING_API_URL_FORM      = HTTPS_API_URL + "/v1/staffmanage/winning/form"          // 中獎人員API(form)

	// 匯入excel檔案
	IMPORT_EXCEL_API_URL = HTTPS_API_URL + "/v1/import/excel"

	// 舉辦活動選單
	ACTIVITY_CREATE_REQUIRE_API_URL = HTTPS_API_URL + "/v1/activity/create/require" // 填寫活動需求

	// 前綴
	PREFIX = "/admin"

	// 檔案存放位置
	STORE_PATH   = "./hilives/hilive/uploads"
	STORE_PREFIX = "uploads"

// NEW_URL        = "/admin/api/new/"
// EDIT_URL       = "/admin/api/edit/"
// DELETE_URL     = "/admin/api/delete/"
)

const (
	// 自定義場景
	CUSTOMIZE_SCENE = "customize_scene"
	// 自定義模板
	CUSTOMIZE_TEMPLATE = "customize_template"

	// 活動
	ACTIVITY_TABLE         = "activity"
	ACTIVITY_2_TABLE       = "activity_2"
	ACTIVITY_CHANNEL_TABLE = "activity_channel"
	ACTIVITY_REQUIRE_TABLE = "activity_require"
	ACTIVITY_TYPE_TABLE    = "activity_type"

	ACTIVITY_APPLYSIGN_TABLE             = "activity_applysign"             // 報名簽到
	ACTIVITY_CUSTOMIZE_TABLE             = "activity_customize"             // 報名簽到自定義
	ACTIVITY_CHATROOM_RECORD_TABLE       = "activity_chatroom_record"       // 聊天紀錄
	ACTIVITY_MESSAGE_SENSITIVITY_TABLE   = "activity_message_sensitivity"   // 敏感詞
	ACTIVITY_QUESTION_GUEST_TABLE        = "activity_question_guest"        // 提問嘉賓
	ACTIVITY_QUESTION_USER_TABLE         = "activity_question_user"         // 提問用戶
	ACTIVITY_QUESTION_LIKES_RECORD_TABLE = "activity_question_likes_record" // 用戶按讚紀錄
	ACTIVITY_SIGNNAME_TABLE              = "activity_signname"              // 簽名牆資料

	// 活動總覽
	// ACTIVITY_OVERVIEW_TABLE      = "activity_overview" // 到時候刪除(修改活動總覽時)
	ACTIVITY_OVERVIEW_GAME_TABLE = "activity_overview_game"
	ACTIVITY_INTRODUCE_TABLE     = "activity_introduce" // 活動介紹
	ACTIVITY_SCHEDULE_TABLE      = "activity_schedule"  // 活動行程
	ACTIVITY_GUEST_TABLE         = "activity_guest"     // 活動嘉賓
	ACTIVITY_MATERIAL_TABLE      = "activity_material"  // 活動資料

	// 所有遊戲、獎品
	ACTIVITY_GAME_SETTING_TABLE = "activity_game_setting"
	ACTIVITY_GAME_TABLE         = "activity_game"

	// 獎品
	ACTIVITY_PRIZE_TABLE = "activity_prize"

	// 快問快答
	ACTIVITY_GAME_QA_RECORD_TABLE = "activity_game_qa_record" // 快問快答答題紀錄

	// 未用到的資料表
	// ACTIVITY_GAME_QA_PICTURE_TABLE_1 = "activity_game_qa_picture"
	// ACTIVITY_GAME_QA_PICTURE_TABLE_2 = "activity_game_qa_picture_2"
	// ACTIVITY_GAME_QA_TABLE           = "activity_game_qa"        // 快問快答問題
	// ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE     = "activity_game_ropepack_picture"
	// ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE     = "activity_game_monopoly_picture"
	// ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE   = "activity_game_whack_mole_picture"
	// ACTIVITY_GAME_REDPACK_PICTURE_TABLE      = "activity_game_redpack_picture"
	// ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE = "activity_game_draw_numbers_picture"
	// ACTIVITY_GAME_LOTTERY_PICTURE_TABLE      = "activity_game_lottery_picture"
	// ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE     = "activity_game_tugofwar_picture"
	// ACTIVITY_GAME_BINGO_PICTURE_TABLE        = "activity_game_bingo_picture"
	// ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE = "activity_game_3d_gacha_machine_picture"
	// ACTIVITY_GAME_VOTE_PICTURE_TABLE         = "activity_game_vote_picture"         // 投票自定義

	// 投票
	ACTIVITY_GAME_VOTE_OPTION_TABLE          = "activity_game_vote_option"          // 投票選項
	ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE = "activity_game_vote_special_officer" // 投票特殊人員
	ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE     = "activity_game_vote_option_list"     // 投票選像名單
	ACTIVITY_GAME_VOTE_RECORD_TABLE          = "activity_game_vote_record"          // 投票紀錄

	// 幸運轉盤
	// ACTIVITY_LOTTERY_TABLE       = "activity_lottery"
	// ACTIVITY_LOTTERY_PRIZE_TABLE = "activity_lottery_prize"
	// 搖紅包
	// ACTIVITY_REDPACK_TABLE       = "activity_redpack"
	// ACTIVITY_REDPACK_PRIZE_TABLE = "activity_redpack_prize"
	// 套紅包
	// ACTIVITY_ROPEPACK_TABLE       = "activity_ropepack"
	// ACTIVITY_ROPEPACK_PRIZE_TABLE = "activity_ropepack_prize"
	// 打地鼠
	// ACTIVITY_WHACK_MOLE_TABLE       = "activity_whack_mole"
	// ACTIVITY_WHACK_MOLE_PRIZE_TABLE = "activity_whack_mole_prize"
	ACTIVITY_SCORE_TABLE = "activity_score"
	// 抽號碼
	// ACTIVITY_DRAW_NUMBERS_PRIZE_TABLE = "activity_draw_numbers_prize"

	// 遊戲、中獎人員、黑名單人員、pk人員
	ACTIVITY_STAFF_GAME_TABLE  = "activity_staff_game"
	ACTIVITY_STAFF_PRIZE_TABLE = "activity_staff_prize"
	ACTIVITY_STAFF_BLACK_TABLE = "activity_staff_black"
	ACTIVITY_STAFF_PK_TABLE    = "activity_staff_pk"

	// 用戶、角色、權限、菜單、操作日至
	USERS_TABLE                = "users"
	USER_PHONE_TABLE           = "user_phone"
	LINE_USERS_TABLE           = "line_users"
	PERMISSION_TABLE           = "permission"
	USER_PERMISSIONS_TABLE     = "user_permissions"     // 用戶權限
	ACTIVITY_PERMISSIONS_TABLE = "activity_permissions" // 活動權限
	OPERATION_LOG_TABLE        = "operation_log"
	OPERATION_ERROR_LOG_TABLE  = "operation_error_log"
	// ROLES_TABLE            = "roles"
	// ROLE_MENU_TABLE        = "role_menu"
	// ROLE_PERMISSIONS_TABLE = "role_permissions"
	// ROLE_USERS_TABLE       = "role_users"
	// SESSION_TABLE = "session"
	// TOKEN_TABLE   = "token"
	MENU_TABLE = "menu"
)

var (
	globalCfg = new(Config) // 設置資料庫...等資訊
	count     uint32
	// updateLock sync.Mutex
	lock sync.Mutex
)

// DatabaseList 資料庫引擎儲存
type DatabaseList map[string]Database

// RedisList redis引擎儲存
type RedisList map[string]Redis

// MongoList mongodb引擎儲存
type MongoList map[string]Mongo

// Database 資料庫引擎設置
type Database struct {
	Host       string            `json:"host"`
	Port       string            `json:"port"`
	User       string            `json:"user"`
	Pwd        string            `json:"pwd"`
	Name       string            `json:"name"`
	MaxIdleCon int               `json:"max_idle_con"`
	MaxOpenCon int               `json:"max_open_con"`
	Driver     string            `json:"driver"`
	File       string            `json:"file"`
	Dsn        string            `json:"dsn"`
	Params     map[string]string `json:"params"`
}

// Redis redis引擎設置
type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
	// User string `json:"user"`
	// Pwd  string `json:"pwd"`
}

// Mongo mongo引擎設置
type Mongo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Name string `json:"name"`
}

// Service service互相轉換
type Service struct {
	Config *Config
}

// Config 基本配置
type Config struct {
	Databases       DatabaseList `json:"database"`          // 支持多個資料庫連接
	RedisList       RedisList    `json:"redis"`             // 支持redis連接
	MongoList       MongoList    `json:"mongodb"`           // 支持mongodb連接
	Store           Store        `json:"store"`             // 檔案儲存位置
	Prefix          string       `json:"prefix"`            // 前綴
	SessionLifeTime int          `json:"session_life_time"` // session 存在時效
	NoLimitLoginIP  bool         `json:"no_limit_login_ip"` // 不限制登入IP
}

// Store 儲存文件
type Store struct {
	Path   string `json:"path" `  // 路徑
	Prefix string `json:"prefix"` // 前綴
}

// SetGlobalConfig 設置全局變數globalCfg
func SetGlobalConfig(cfg Config) *Config {
	lock.Lock()
	defer lock.Unlock()
	if atomic.LoadUint32(&count) != 0 {
		panic("基本設置不能重複設置")
	} else {
		atomic.StoreUint32(&count, 1)
	}
	cfg = SetDefault(cfg)
	globalCfg = &cfg
	return globalCfg
}

// SetDefault 設置Config預設值
func SetDefault(cfg Config) Config {
	if cfg.SessionLifeTime == 0 {
		cfg.SessionLifeTime = 86400 * 7 // 1天
		// cfg.SessionLifeTime = 1800
	}
	return cfg
}

// GetService 先將參數轉換成Service後回傳config
func GetService(s interface{}) *Config {
	if srv, ok := s.(*Service); ok {
		return srv.Config
	}
	panic("錯誤的service")
}

// Name Service方法
func (c *Service) Name() string {
	return "config"
}

// ServiceWithConfig 將Config轉換成Service
func ServiceWithConfig(c *Config) *Service {
	return &Service{c}
}

// GetDatabases 取得全局變數(globalCfg)資訊
func GetDatabases() DatabaseList {
	var list = make(DatabaseList, len(globalCfg.Databases))
	for key := range globalCfg.Databases {
		list[key] = Database{
			Driver: globalCfg.Databases[key].Driver,
		}
	}
	return list
}

// GroupByDriver 依照引擎分組
func (d DatabaseList) GroupByDriver() map[string]DatabaseList {
	drivers := make(map[string]DatabaseList)
	for key, item := range d {
		if driverList, ok := drivers[item.Driver]; ok {
			driverList.Add(key, item)
		} else {
			drivers[item.Driver] = make(DatabaseList)
			drivers[item.Driver].Add(key, item)
		}
	}
	return drivers
}

// Add 將資料庫引擎加入DatabaseList
func (d DatabaseList) Add(key string, db Database) {
	d[key] = db
}

// GetHilive DatabaseList["hilive"]
func (d DatabaseList) GetHilive() Database {
	return d["hilive"]
}

// URL 處理url
func (c *Config) URL(suffix string) string {
	if c.Prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return c.Prefix
	}
	return c.Prefix + suffix
}

// URL 處理儲存資料的路徑
func (s Store) URL(suffix string) string {
	if len(suffix) > 4 && suffix[:4] == "http" {
		return suffix
	}
	if s.Prefix == "" {
		if suffix[0] == '/' {
			return suffix
		}
		return "/" + suffix
	}
	if s.Prefix[0] == '/' {
		if suffix[0] == '/' {
			return s.Prefix + suffix
		}
		return s.Prefix + "/" + suffix
	}
	if suffix[0] == '/' {
		if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
			return s.Prefix + suffix
		}
		return "/" + s.Prefix + suffix
	}
	if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
		return s.Prefix + "/" + suffix
	}
	return "/" + s.Prefix + "/" + suffix
}

// JSON JSON編碼
func (s Store) JSON() string {
	if s.Path == "" && s.Prefix == "" {
		return ""
	}
	return utils.JSON(s)
}

// Prefix 前綴
func Prefix() string {
	return globalCfg.Prefix
}

// GetSessionLifeTime cookie時間
func GetSessionLifeTime() int {
	return globalCfg.SessionLifeTime
}

// GetNoLimitLoginIP 是否限制登入IP
func GetNoLimitLoginIP() bool {
	return globalCfg.NoLimitLoginIP
}

// GetStore Store
func GetStore() Store {
	return globalCfg.Store
}
