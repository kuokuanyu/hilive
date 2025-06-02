package engine

import (
	"fmt"
	"hilive/controller"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/service"
	"hilive/modules/table"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// Engine 核心組件
type Engine struct {
	dbConn    db.Connection    // 資料庫 CRUD 功能
	redisConn cache.Connection // redis功能
	mongoConn mongo.Connection // mongo功能
	config    *config.Config   // 基本設置
	Gin       *gin.Engine      // gin 框架
	// Gin      *fiber.App          // gin 框架
	Services service.List        // 儲存資料庫引擎、 Config 、 token ...等資訊
	handler  *controller.Handler // API 功能應用
}

// DefaultEngine 預設Engine
func DefaultEngine() *Engine {
	return &Engine{
		Services: service.GetServices(),
	}
}

// InitDatabase 初始化資料庫引擎，加入Services
func (eng *Engine) InitDatabase(cfg config.Config) *Engine {
	eng.config = config.SetGlobalConfig(cfg)
	for driver, databaseCfg := range eng.config.Databases.GroupByDriver() {
		eng.Services.Add(driver, db.GetConnectionByDriver(driver).InitDB(databaseCfg))
	}
	eng.dbConn = eng.GetDBConnection("mysql")

	// #####正式區執行一次後就必須關閉-----start#####

	// #####正式區執行一次後就必須關閉-----end#####

	return eng
}

// InitMongo 初始化mongodb引擎，加入Services
func (eng *Engine) InitMongo(cfg config.Config) *Engine {
	eng.Services.Add("mongo", mongo.GetConnection().InitMongo(cfg.MongoList))
	eng.mongoConn = eng.GetMongoConnection()

	return eng
}

// InitRedis 初始化redis引擎，加入Services，清除舊的資料
func (eng *Engine) InitRedis(cfg config.Config) *Engine {
	eng.Services.Add("redis", cache.GetConnection().InitRedis(cfg.RedisList))
	eng.redisConn = eng.GetRedisConnection()

	// #####正式區執行一次後就必須關閉-----start#####
	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序

	log.Println("清除平台用戶平台用戶redis")
	hiliveusers, _ := eng.dbConn.Query("SELECT * FROM users;")

	// 清除舊的平台用戶資料(redis)
	wg.Add(len(hiliveusers)) //計數器
	for i := 0; i < len(hiliveusers); i++ {
		go func() {
			defer wg.Done()
			userID, _ := hiliveusers[i]["user_id"].(string)

			// 清除所有redis紀錄資料
			eng.redisConn.DelCache(config.HILIVE_USERS_REDIS + userID)

			// log.Println("建立資料夾: ", userID)

			// 建立報名資料夾
			// os.MkdirAll("./hilives/hilive/uploads/"+userID+"/customize_scene", os.ModePerm)

		}()
	}
	wg.Wait() //等待計數器歸0

	// log.Println("清除LINE用戶redis")
	// lineusers, _ := eng.dbConn.Query("SELECT * FROM line_users;")

	// // 清除舊的LINE用戶資料(redis)
	// wg.Add(len(lineusers)) //計數器
	// for i := 0; i < len(lineusers); i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		userID, _ := lineusers[i]["user_id"].(string)

	// 		// 清除所有redis紀錄資料
	// 		eng.redisConn.DelCache(config.AUTH_USERS_REDIS + userID)

	// 	}()
	// }
	// wg.Wait() //等待計數器歸0

	// log.Println("清除活動redis")

	// 查詢所有活動資料
	activitys, _ := eng.dbConn.Query("SELECT * FROM activity;")

	wg.Add(len(activitys)) //計數器

	// 清除舊的活動聊天紀錄資料(redis)
	for _, activity := range activitys {
		go func(activity map[string]interface{}) {
			defer wg.Done()

			activityID, _ := activity["activity_id"].(string)

			// mongo
			filter := bson.M{"activity_id": activityID}
			update := bson.M{
				"$set": bson.M{
					"channel_1": "close", "channel_2": "close", "channel_3": "close", "channel_4": "close", "channel_5": "close",
					"channel_6": "close", "channel_7": "close", "channel_8": "close", "channel_9": "close", "channel_10": "close"},
				// "$unset": bson.M{}, // 移除不需要的欄位
			}

			eng.mongoConn.UpdateOne(config.ACTIVITY_CHANNEL_TABLE, filter, update)

			// 需要清除所有活動 redis 紀錄資料
			redisKeys := []string{
				config.ACTIVITY_REDIS + activityID,                    // 活動資訊，HASH
				config.ACTIVITY_NUMBER_REDIS + activityID,             // 活動 number，STRING
				config.SIGN_STAFFS_2_REDIS + activityID,               // 簽到人員2（更新資料時，修改 redis 裡的資料，james用），SET
				config.SIGNNAME_REDIS + activityID,                    // 簽名牆設置資料，HASH
				config.HOST_CONTROL_CHANNEL_REDIS + activityID,        // 主持端所有可遙控的 session，SET
				config.USER_GAME_RECORDS_REDIS + activityID,           // 玩家在該活動下的遊戲紀錄（中獎、未中獎），HASH
				config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID, // 未中獎人員，LIST
				config.BLACK_STAFFS_ACTIVITY_REDIS + activityID,       // 黑名單人員（活動），SET
				config.BLACK_STAFFS_MESSAGE_REDIS + activityID,        // 黑名單人員（訊息），SET
				config.BLACK_STAFFS_QUESTION_REDIS + activityID,       // 黑名單人員（提問），SET
				config.BLACK_STAFFS_SIGNNAME_REDIS + activityID,       // 黑名單人員（簽名），SET
			}

			// 需要刪除的遙控頻道redis，用迴圈建立(1~10)
			for i := 1; i <= 10; i++ {
				redisKeys = append(redisKeys, fmt.Sprintf("%s%s_channel_%d", config.HOST_CONTROL_REDIS, activityID, i))

			}

		}(activity) // 傳值進 goroutine，避免閉包錯誤
	}
	wg.Wait() //等待計數器歸0

	// 查詢所有遊戲資料(mysql)
	// log.Println("清除遊戲redis(mysql)")

	// games, _ := eng.dbConn.Query("SELECT * FROM activity_game;")

	// wg.Add(len(games)) //計數器

	// // 清除舊的遊戲紀錄資料(redis)
	// for _, gameItem := range games {
	// 	go func(game map[string]interface{}) {
	// 		defer wg.Done()

	// 		gameID, _ := game["game_id"].(string)
	// 		gameType, _ := game["game"].(string)
	// 		gameStatus, _ := game["game_status"].(string)

	// 		// 更新遊戲狀態
	// 		if gameType != "vote" && gameStatus != "close" {
	// 			eng.dbConn.Exec(fmt.Sprintf(`update activity_game set
	// 			game_status = 'close' where game_id = '%s';`, gameID))
	// 		}

	// 		// 要刪除的 Redis keys 統一集中處理
	// 		redisKeys := []string{
	// 			config.GAME_REDIS + gameID,                            // 遊戲資訊，HASH
	// 			config.GAME_TYPE_REDIS + gameID,                       // 遊戲種類資訊，STRING
	// 			config.GAME_PRIZES_AMOUNT_REDIS + gameID,              // 遊戲獎品數量，HASH
	// 			config.GAME_VOTE_RECORDS_REDIS + gameID,               // 玩家投票紀錄，HASH
	// 			config.VOTE_AVATAR_REDIS + gameID,                     // 玩家投票頭像，HASH
	// 			config.BLACK_STAFFS_GAME_REDIS + gameID,               // 黑名單人員(遊戲)，SET
	// 			config.SCORES_REDIS + gameID,                          // 玩家分數，ZSET
	// 			config.SCORES_2_REDIS + gameID,                        // 玩家第二分數，ZSET
	// 			config.CORRECT_REDIS + gameID,                         // 玩家答對題數，ZSET
	// 			config.WINNING_STAFFS_REDIS + gameID,                  // 中獎人員，LIST
	// 			config.NO_WINNING_STAFFS_REDIS + gameID,               // 未中獎人員，LIST
	// 			config.QA_REDIS + gameID,                              // 題數資訊，HASH
	// 			config.QA_RECORD_REDIS + gameID,                       // 答題紀錄資訊，HASH
	// 			config.GAME_TEAM_REDIS + gameID,                       // 遊戲隊伍資訊，HASH
	// 			config.GAME_BINGO_NUMBER_REDIS + gameID,               // 賓果抽過的號碼，LIST
	// 			config.GAME_BINGO_USER_REDIS + gameID,                 // 賓果中獎人員，ZSET
	// 			config.GAME_BINGO_USER_NUMBER + gameID,                // 玩家賓果號碼排序，HASH
	// 			config.GAME_BINGO_USER_GOING_BINGO + gameID,           // 玩家是否即將賓果，HASH
	// 			config.VOTE_SPECIAL_OFFICER_REDIS + gameID,            // 投票遊戲特殊人員，HASH
	// 			config.GAME_ATTEND_REDIS + gameID,                     // 遊戲人數資訊，SET
	// 			config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID,  // 拔河遊戲左隊人數資訊，SET
	// 			config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID, // 拔河遊戲右隊人數資訊，SET
	// 		}

	// 		// 刪除 Redis 中的 key
	// 		for _, key := range redisKeys {
	// 			eng.redisConn.DelCache(key)
	// 		}

	// 	}(gameItem)
	// }
	// wg.Wait() //等待計數器歸0

	// 查詢所有遊戲資料(mongo)
	log.Println("清除遊戲redis(mongo)")

	mongoGames, _ := eng.mongoConn.FindMany("activity_game", bson.M{})

	wg.Add(len(mongoGames)) //計數器

	// 清除舊的遊戲紀錄資料(redis)
	for _, gameItem := range mongoGames {
		go func(game map[string]interface{}) {
			defer wg.Done()

			gameID, _ := game["game_id"].(string)
			gameType, _ := game["game"].(string)
			gameStatus, _ := game["game_status"].(string)

			// 更新遊戲狀態
			// mongo
			if gameType != "vote" && gameStatus != "close" {
				filter := bson.M{"game_id": gameID}
				update := bson.M{
					"$set": bson.M{"game_status": "close"},
					// "$unset": bson.M{}, // 移除不需要的欄位
				}

				eng.mongoConn.UpdateOne("activity_game", filter, update)
			}

			// 要刪除的 Redis keys 統一集中處理
			redisKeys := []string{
				config.GAME_REDIS + gameID,                            // 遊戲資訊，HASH
				config.GAME_TYPE_REDIS + gameID,                       // 遊戲種類資訊，STRING
				config.GAME_PRIZES_AMOUNT_REDIS + gameID,              // 遊戲獎品數量，HASH
				config.GAME_VOTE_RECORDS_REDIS + gameID,               // 玩家投票紀錄，HASH
				config.VOTE_AVATAR_REDIS + gameID,                     // 玩家投票頭像，HASH
				config.BLACK_STAFFS_GAME_REDIS + gameID,               // 黑名單人員(遊戲)，SET
				config.SCORES_REDIS + gameID,                          // 玩家分數，ZSET
				config.SCORES_2_REDIS + gameID,                        // 玩家第二分數，ZSET
				config.CORRECT_REDIS + gameID,                         // 玩家答對題數，ZSET
				config.WINNING_STAFFS_REDIS + gameID,                  // 中獎人員，LIST
				config.NO_WINNING_STAFFS_REDIS + gameID,               // 未中獎人員，LIST
				config.QA_REDIS + gameID,                              // 題數資訊，HASH
				config.QA_RECORD_REDIS + gameID,                       // 答題紀錄資訊，HASH
				config.GAME_TEAM_REDIS + gameID,                       // 遊戲隊伍資訊，HASH
				config.GAME_BINGO_NUMBER_REDIS + gameID,               // 賓果抽過的號碼，LIST
				config.GAME_BINGO_USER_REDIS + gameID,                 // 賓果中獎人員，ZSET
				config.GAME_BINGO_USER_NUMBER + gameID,                // 玩家賓果號碼排序，HASH
				config.GAME_BINGO_USER_GOING_BINGO + gameID,           // 玩家是否即將賓果，HASH
				config.VOTE_SPECIAL_OFFICER_REDIS + gameID,            // 投票遊戲特殊人員，HASH
				config.GAME_ATTEND_REDIS + gameID,                     // 遊戲人數資訊，SET
				config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID,  // 拔河遊戲左隊人數資訊，SET
				config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID, // 拔河遊戲右隊人數資訊，SET
			}

			// 刪除 Redis 中的 key
			for _, key := range redisKeys {
				eng.redisConn.DelCache(key)
			}

		}(gameItem)
	}
	wg.Wait() //等待計數器歸0

	// #####正式區執行一次後就必須關閉-----end#####

	// #####正式區執行一次後就必須關閉-----start#####
	// 查詢所有活動
	// activitysModel, _ := models.DefaultActivityModel().SetDbConn(eng.dbConn).FindOpenActivitys()

	// 舊的活動場次新增資料夾
	// for _, activity := range activitysModel {
	// 	log.Println("建立資料夾: ", activity.UserID, activity.ActivityID)

	// 	// 建立報名資料夾
	// 	os.MkdirAll("./hilives/hilive/uploads/"+activity.UserID+"/"+activity.ActivityID+"/applysign/apply", os.ModePerm)
	// }
	// 正式區執行一次後就必須關閉-----end

	return eng
}

// SetEngine 設置engine資訊
func (eng *Engine) SetEngine() *Engine {
	st := table.NewSystemTable(eng.dbConn, eng.redisConn, eng.mongoConn, eng.config)
	tablelist := map[string]table.Generator{
		// 管理員
		"admin_manager":    st.GetAdminManagerPanel,
		"admin_permission": st.GetAdminPermissionPanel,
		"admin_overview":   st.GetAdminOverviewPanel,
		"admin_menu":       st.GetAdminMenuPanel,

		// 紀錄
		"chatroom_record": st.GetChatroomRecordPanel,
		"question_record": st.GetQuestionRecordPanel,

		"user":                st.GetUserPanel,
		"line_user":           st.GetLineUserPanel,
		"activity":            st.GetActivityPanel,
		"require":             st.GetActivityRequirePanel,
		"overview":            st.GetOverviewPanel,
		"introduce":           st.GetIntroducePanel,
		"schedule":            st.GetSchedulePanel,
		"guest":               st.GetGuestPanel,
		"material":            st.GetMaterialPanel,
		"applysign":           st.GetApplysignPanel,
		"applysign_user":      st.GetApplysignUserPanel,
		"applysign_users":     st.GetApplysignUsersPanel,
		"apply":               st.GetApplyPanel,
		"sign":                st.GetSignPanel,
		"customize":           st.GetCustomizePanel,
		"qrcode":              st.GetQRcodePanel,
		"message":             st.GetMessagePanel,
		"message_check":       st.GetMessageCheckPanel,
		"message_sensitivity": st.GetMessageMessageSensitivityPanel,
		"topic":               st.GetTopicPanel,
		"question":            st.GetQuestionPanel,
		"danmu":               st.GetDanmuPanel,
		"specialdanmu":        st.GetSepcialDanmuPanel,
		// "picture":             st.GetPicturePanel,
		"holdscreen": st.GetHoldScreenPanel,
		"general":    st.GetGeneralPanel,
		"threed":     st.GetThreedPanel,
		"countdown":  st.GetCountdownPanel,
		"signname":   st.GetSignnamePanel,
		// 遊戲
		"setting":            st.GetGameSettingPanel,
		"lottery":            st.GetLotteryPanel,
		"lottery_prize":      st.GetLotteryPrizePanel,
		"redpack":            st.GetRedpackPanel,
		"redpack_prize":      st.GetRedpackPrizePanel,
		"ropepack":           st.GetRopepackPanel,
		"ropepack_prize":     st.GetRopepackPrizePanel,
		"whack_mole":         st.GetWhackMolePanel,
		"whack_mole_prize":   st.GetWhackMolePrizePanel,
		"draw_numbers_prize": st.GetDrawNumbersPrizePanel,
		"monopoly":           st.GetMonopolyPanel,
		"monopoly_prize":     st.GetMonopolyPrizePanel,
		"QA":                 st.GetQAPanel,
		"QA_prize":           st.GetQAPrizePanel,
		// 拔河遊戲
		"tugofwar":       st.GetTugofwarPanel,
		"tugofwar_prize": st.GetTugofwarPrizePanel,

		// 賓果遊戲
		"bingo":       st.GetBingoPanel,
		"bingo_prize": st.GetBingoPrizePanel,

		// 扭蛋機遊戲
		"3DGachaMachine":       st.Get3DGachaMachinePanel,
		"3DGachaMachine_prize": st.Get3DGachaMachinePrizePanel,

		// 投票遊戲
		"vote": st.GetVotePanel,
		// "vote_prize":            st.GetVotePrizePanel,
		"vote_option":           st.GetVoteOptionPanel,          // 投票選項
		"vote_special_officer":  st.GetVoteSpecialOfficerPanel,  // 投票特殊人員
		"vote_option_list":      st.GetVoteOptionListPanel,      // 投票選項名單
		"vote_special_officers": st.GetVoteSpecialOfficersPanel, // 投票特殊人員
		"vote_option_lists":     st.GetVoteOptionListsPanel,     // 投票選項名單

		"attend":  st.GetAttendPanel,
		"winning": st.GetWinningPanel,
		"black":   st.GetBlackPanel,
		// "pk":      st.GetPKPanel,

		"record":       st.GetRecordPanel,
		"draw_numbers": st.GetDrawNumbersPanel,
		// "reset":              st.GetGameResetPanel,

	}
	eng.handler = controller.NewHandler(eng.config, eng.Services, eng.dbConn, eng.redisConn, eng.mongoConn, tablelist)
	return eng
}

// GetDBConnection 取得Connection
func (eng *Engine) GetDBConnection(driver string) db.Connection {
	return db.GetConnectionFromService(eng.Services.Get(driver))
}

// GetRedisConnection 取得Connection
func (eng *Engine) GetRedisConnection() cache.Connection {
	return cache.GetConnectionFromService(eng.Services.Get("redis"))
}

// GetMongoConnection 取得Connection
func (eng *Engine) GetMongoConnection() mongo.Connection {
	return mongo.GetConnectionFromService(eng.Services.Get("mongo"))
}

// 更新所有頻道狀態(改用mongo)
// eng.dbConn.Exec(fmt.Sprintf(`update activity_channel set
// channel_1 = 'close', channel_2 = 'close', channel_3 = 'close',
// channel_4 = 'close', channel_5 = 'close', channel_6 = 'close',
// channel_7 = 'close', channel_8 = 'close', channel_9 = 'close',
// channel_10 = 'close' where activity_id = '%s';`, activityID))

// 查詢所有活動
// activitys, _ := models.DefaultActivityModel().SetDbConn(eng.dbConn).FindOpenActivitys()

// 舊的活動場次新增資料夾
// for _, activity := range activitys {
// 	// log.Println("建立資料夾: ", activity.UserID, activity.ActivityID)

// 	// 建立3D扭蛋機資料夾
// 	// os.MkdirAll("./hilives/hilive/uploads/"+activity.UserID+"/"+activity.ActivityID+"/interact/game/3DGachaMachine", os.ModePerm)
// }

// eng.redisConn.DelCache(config.ACTIVITY_REDIS + activityID)        // 活動
// eng.redisConn.DelCache(config.ACTIVITY_NUMBER_REDIS + activityID) // 活動number，STRING

// eng.redisConn.DelCache(config.SIGN_STAFFS_2_REDIS + activityID) // 簽到人員2(更新資料時，修改redis裡的資料，james用)，SET

// eng.redisConn.DelCache(config.SIGNNAME_REDIS + activityID) // 簽名牆設置資料，HASH

// 頻道
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_1")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_2")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_3")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_4")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_5")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_6")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_7")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_8")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_9")  // 主持端遙控資訊，HASH
// eng.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_10") // 主持端遙控資訊，HASH

// eng.redisConn.DelCache(config.HOST_CONTROL_CHANNEL_REDIS + activityID) // 主持端所有可遙控的session，SET

// eng.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)           // 玩家在該活動下的遊戲紀錄(中獎.未中獎)，HASH
// eng.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 未中獎人員，LIST

// eng.redisConn.DelCache(config.BLACK_STAFFS_ACTIVITY_REDIS + activityID) // 黑名單人員(活動)，SET
// eng.redisConn.DelCache(config.BLACK_STAFFS_MESSAGE_REDIS + activityID)  // 黑名單人員(訊息)，SET
// eng.redisConn.DelCache(config.BLACK_STAFFS_QUESTION_REDIS + activityID) // 黑名單人員(提問)，SET
// eng.redisConn.DelCache(config.BLACK_STAFFS_SIGNNAME_REDIS + activityID) // 黑名單人員(簽名)，SET

// 清除所有舊的遊戲紀錄資料
// eng.redisConn.DelCache(config.GAME_REDIS + gameID)      // 遊戲資訊，HASH
// eng.redisConn.DelCache(config.GAME_TYPE_REDIS + gameID) // 遊戲種類資訊，STRING

// eng.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 遊戲獎品數量，HASH

// eng.redisConn.DelCache(config.GAME_VOTE_RECORDS_REDIS + gameID) // 玩家投票紀錄，HASH
// eng.redisConn.DelCache(config.VOTE_AVATAR_REDIS + gameID)       // 玩家投票紀錄，HASH

// eng.redisConn.DelCache(config.BLACK_STAFFS_GAME_REDIS + gameID) // 黑名單人員(遊戲)，SET

// eng.redisConn.DelCache(config.SCORES_REDIS + gameID)                // 玩家分數，ZSET
// eng.redisConn.DelCache(config.SCORES_2_REDIS + gameID)              // 玩家第二分數，ZSET
// eng.redisConn.DelCache(config.CORRECT_REDIS + gameID)               // 玩家答對題數，ZSET
// eng.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)        // 中獎人員，LIST
// eng.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + gameID)     // 未中獎人員，LIST
// eng.redisConn.DelCache(config.QA_REDIS + gameID)                    // 題數資訊，HASH
// eng.redisConn.DelCache(config.QA_RECORD_REDIS + gameID)             // 答題紀錄資訊，HASH
// eng.redisConn.DelCache(config.GAME_TEAM_REDIS + gameID)             // 遊戲隊伍資訊，HASH
// eng.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + gameID)     // 紀錄抽過的號碼，LIST
// eng.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + gameID)       // 賓果中獎人員，ZSET
// eng.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + gameID)      // 紀錄玩家的號碼排序，HASH
// eng.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + gameID) // 紀錄玩家是否即將賓果，HASH
// eng.redisConn.DelCache(config.VOTE_SPECIAL_OFFICER_REDIS + gameID)  // 投票遊戲特殊人員，HASH

// eng.redisConn.DelCache(config.GAME_ATTEND_REDIS + gameID)                     // 遊戲人數資訊，SET
// eng.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 拔河遊戲左隊人數資訊，SET
// eng.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 拔河遊戲右隊人數資訊，SET
