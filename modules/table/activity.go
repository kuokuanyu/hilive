package table

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// GetActivityPanel 活動
func (s *SystemTable) GetActivityPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)

		tables := []string{
			// 活動權限
			config.ACTIVITY_PERMISSIONS_TABLE,

			// 活動資料表
			config.ACTIVITY_TABLE,
			config.ACTIVITY_2_TABLE,
			// config.ACTIVITY_CHANNEL_TABLE,

			// 活動資訊
			config.ACTIVITY_INTRODUCE_TABLE, config.ACTIVITY_SCHEDULE_TABLE,
			config.ACTIVITY_GUEST_TABLE, config.ACTIVITY_MATERIAL_TABLE,

			// 報名簽到
			config.ACTIVITY_APPLYSIGN_TABLE, config.ACTIVITY_CUSTOMIZE_TABLE,

			// 聊天.提問紀錄
			config.ACTIVITY_CHATROOM_RECORD_TABLE, config.ACTIVITY_QUESTION_LIKES_RECORD_TABLE,
			config.ACTIVITY_QUESTION_USER_TABLE, config.ACTIVITY_QUESTION_GUEST_TABLE,

			// 聊天區設定
			config.ACTIVITY_MESSAGE_SENSITIVITY_TABLE,

			// 簽名牆
			config.ACTIVITY_SIGNNAME_TABLE,

			// 遊戲
			config.ACTIVITY_GAME_SETTING_TABLE,
			// config.ACTIVITY_GAME_TABLE,
			// config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE,
			// config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE,
			// config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE,
			// config.ACTIVITY_GAME_QA_PICTURE_TABLE_1,
			// config.ACTIVITY_GAME_QA_PICTURE_TABLE_2,
			// config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE,
			// config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE,
			// config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE,
			// config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE,
			// config.ACTIVITY_GAME_BINGO_PICTURE_TABLE,
			// config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE,

			// 投票
			// config.ACTIVITY_GAME_VOTE_PICTURE_TABLE,
			config.ACTIVITY_GAME_VOTE_OPTION_TABLE,
			config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE,
			config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE,
			config.ACTIVITY_GAME_VOTE_RECORD_TABLE,

			// config.ACTIVITY_GAME_QA_TABLE,

			// 獎品
			config.ACTIVITY_PRIZE_TABLE,

			// 人員管理
			config.ACTIVITY_STAFF_GAME_TABLE, config.ACTIVITY_STAFF_PRIZE_TABLE,
			config.ACTIVITY_STAFF_BLACK_TABLE, config.ACTIVITY_STAFF_PK_TABLE,

			// 分數.遊戲紀錄
			config.ACTIVITY_SCORE_TABLE, config.ACTIVITY_GAME_QA_RECORD_TABLE,
		}
		for _, t := range tables {
			s.table(t).WhereIn("activity_id", ids).Delete()
		}

		mongoTables := []string{
			config.ACTIVITY_CHANNEL_TABLE, // 活動頻道
			config.ACTIVITY_GAME_TABLE,    // 遊戲資訊
		}
		for _, t := range mongoTables {
			s.mongoConn.DeleteMany(t, bson.M{"activity_id": bson.M{"$in": ids}})
		}

		for _, id := range ids {
			idStr := fmt.Sprintf("%s", id)

			// log.Println("idStr: ", idStr, "activity_id: ", activityID)
			// log.Println("%s%s_channel_%d: ", fmt.Sprintf("%s%s_channel_%d", config.HOST_CONTROL_REDIS, idStr, 1))

			// 需要刪除的key清單
			cacheKeys := []string{
				config.ACTIVITY_REDIS + idStr,                    // 活動資訊
				config.ACTIVITY_NUMBER_REDIS + idStr,             // 活動number
				config.SIGN_STAFFS_2_REDIS + idStr,               // 簽到人員資訊
				config.HOST_CONTROL_CHANNEL_REDIS + idStr,        // 主持端所有可遙控的session
				config.BLACK_STAFFS_ACTIVITY_REDIS + idStr,       // 黑名單人員(活動)
				config.BLACK_STAFFS_MESSAGE_REDIS + idStr,        // 黑名單人員(訊息)
				config.BLACK_STAFFS_QUESTION_REDIS + idStr,       // 黑名單人員(提問)
				config.BLACK_STAFFS_SIGNNAME_REDIS + idStr,       // 黑名單人員(簽名)
				config.USER_GAME_RECORDS_REDIS + idStr,           // 玩家遊戲紀錄(中獎.未中獎)
				config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + idStr, // 活動下所有場次搖號抽獎的中獎人員，SET
			}

			// 需要刪除的遙控頻道redis，用迴圈建立(1~10)
			for i := 1; i <= 10; i++ {
				cacheKeys = append(cacheKeys, fmt.Sprintf("%s%s_channel_%d", config.HOST_CONTROL_REDIS, idStr, i))
			}

			// 刪除 Redis cache
			for _, key := range cacheKeys {
				s.redisConn.DelCache(key)
			}

			// 需要pub的頻道發送清單
			channels := []string{
				config.CHANNEL_SIGN_STAFFS_2_REDIS + idStr,
				config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + idStr,
				config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + idStr,
				config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + idStr,
				config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + idStr,
				config.CHANNEL_SIGNNAME_EDIT_TIMES_REDIS + idStr,
				config.CHANNEL_GENERAL_EDIT_TIMES_REDIS + idStr,
				config.CHANNEL_THREED_EDIT_TIMES_REDIS + idStr,
				config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS + idStr,
			}

			// 需要pub的遙控頻道redis，用迴圈建立(1~10)
			for i := 1; i <= 10; i++ {
				channels = append(channels, fmt.Sprintf("%s%s_channel_%d", config.CHANNEL_HOST_CONTROL_REDIS, idStr, i))
			}

			// 發布 Redis 頻道
			for _, ch := range channels {
				s.redisConn.Publish(ch, "修改資料")
			}

			// 刪除活動資料夾
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + idStr)

		}
		return nil
	})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("user", "activity_name", "activity_type", "people",
			"city", "town", "start_time", "end_time") {
			return errors.New("錯誤: ID、活動參數不能為空")
		}

		activityID := utils.UUID(20)

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditActivityModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		// 自己加上ActivityID因為 values 裡面沒有
		model.ActivityID = activityID

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, model); err != nil {
			return err
		}

		// 將struct中空值的欄位印出
		// utils.PrintStructFields(model)

		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("user", "activity_id") {
			return errors.New("錯誤: 用戶、活動ID不能為空，請輸入有效的活動ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditActivityModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultActivityModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			UpdateActivity(true, model); err != nil {
			return err
		}

		// 將struct中空值的欄位印出
		// utils.PrintStructFields(model)

		return nil
	})
	return
}

// func applyValuesToStruct(model interface{}, values url.Values) {
// 	v := reflect.ValueOf(model).Elem()
// 	t := v.Type()

// 	for i := 0; i < v.NumField(); i++ {
// 		field := t.Field(i)
// 		fieldName := field.Name

// 		// 試著用欄位名稱當作 key 去找值
// 		if val, ok := values[fieldName]; ok && len(val) > 0 {
// 			f := v.FieldByName(fieldName)
// 			if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
// 				f.SetString(val[0])
// 			}
// 			// 這邊你也可以再延伸支援 int、bool 轉換等
// 		}
// 	}
// }

// @Summary 新增活動
// @Tags Activity
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_name formData string true "活動名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param activity_type formData string true "活動類型" Enums(企業會議, 其他, 商業活動, 培訓/教育, 婚禮, 年會, 校園活動, 競技賽事, 論壇會議, 酒吧/餐飲娛樂, 電視/媒體)
// @param max_people formData integer true "活動人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "活動人數上限(依照max_people資料判斷上限)" minimum(1)
// @param city formData string true "活動舉辦縣市" Enums(基隆市, 臺北市, 新北市, 宜蘭縣, 新竹市, 新竹縣, 桃園市, 苗栗縣, 臺中市, 彰化縣, 南投縣, 嘉義市, 嘉義縣, 雲林縣, 臺南市, 高雄市, 屏東縣, 臺東縣, 花蓮縣, 金門縣, 連江縣, 澎湖縣)
// @param town formData string true "活動舉辦區域"
// @param start_time formData string true "活動開始時間(西元年-月-日T時:分)"
// @param end_time formData string true "活動結束時間(西元年-月-日T時:分)"
// @param login_required formData string true "是否需要登入才可進入聊天室" Enums(open, close)
// @param login_password formData string false "登入密碼" maxlength(8)
// @param password_required formData string true "是否需要設置密碼才可進入聊天室" Enums(open, close)
// @param permissions formData string false "用戶權限(以,分隔權限)"
// @param device formData string true "device" Enums(line, facebook, gmail, customize)
// @param line_id formData string false "line_id"
// @param channel_id formData string false "channel_id"
// @param channel_secret formData string false "channel_secret"
// @param chatbot_secret formData string false "chatbot_secret"
// @param chatbot_token formData string false "chatbot_token"
// @param push_message formData string true "官方帳號傳送訊息" Enums(open, close)
// @param message_amount formData integer false "簡訊數量" minimum(1)
// @param push_phone_message formData string true "手機傳送簡訊" Enums(open, close)
// @param mail_amount formData integer false "郵件數量" minimum(1)
// @param send_mail formData string true "發送郵件" Enums(open, close)
// @param host_scan formData string true "主持人掃描qrcode" Enums(open, close)
// @param channel_amount formData integer false "channel_amount" minimum(1) maximum(10)
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /activity [post]
func POSTActivity(ctx *gin.Context) {
}

// @Summary 編輯活動
// @Tags Activity
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param activity_name formData string false "活動名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param activity_type formData string false "活動類型" Enums(企業會議, 其他, 商業活動, 培訓/教育, 婚禮, 年會, 校園活動, 競技賽事, 論壇會議, 酒吧/餐飲娛樂, 電視/媒體)
// @param max_people formData integer false "活動人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "活動人數上限(依照max_people資料判斷上限)" minimum(1)
// @param city formData string false "活動舉辦縣市" Enums(基隆市, 臺北市, 新北市, 宜蘭縣, 新竹市, 新竹縣, 桃園市, 苗栗縣, 臺中市, 彰化縣, 南投縣, 嘉義市, 嘉義縣, 雲林縣, 臺南市, 高雄市, 屏東縣, 臺東縣, 花蓮縣, 金門縣, 連江縣, 澎湖縣)
// @param town formData string false "活動舉辦區域"
// @param start_time formData string false "活動開始時間(西元年-月-日T時:分)"
// @param end_time formData string false "活動結束時間(西元年-月-日T時:分)"
// @param login_required formData string false "是否需要登入才可進入聊天室" Enums(open, close)
// @param login_password formData string false "聊天室驗證碼" minlength(1) maxlength(8)
// @param password_required formData string false "是否需要設置密碼才可進入聊天室" Enums(open, close)
// @param permissions formData string false "用戶權限(以,分隔權限)"
// @param device formData string true "device" Enums(line, facebook, gmail, customize)
// @param line_id formData string false "line_id"
// @param channel_id formData string false "channel_id"
// @param channel_secret formData string false "channel_secret"
// @param chatbot_secret formData string false "chatbot_secret"
// @param chatbot_token formData string false "chatbot_token"
// @param push_message formData string false "官方帳號傳送訊息" Enums(open, close)
// @param message_amount formData integer false "簡訊數量" minimum(1)
// @param push_phone_message formData string true "手機傳送簡訊" Enums(open, close)
// @param mail_amount formData integer false "郵件數量" minimum(1)
// @param send_mail formData string true "發送郵件" Enums(open, close)
// @param host_scan formData string true "主持人掃描qrcode" Enums(open, close)
// @param channel_amount formData integer false "channel_amount" minimum(1) maximum(10)
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /activity [put]
func PUTActivity(ctx *gin.Context) {
}

// @Summary 刪除活動
// @Tags Activity
// @version 1.0
// @Accept  mpfd
// @param id formData string true "活動ID(以,區隔多筆活動ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /activity [DELETE]
func DELETEActivity(ctx *gin.Context) {
}

// @Summary 活動資訊JSON資料(包含活動權限)，未填寫activity_id則查詢所有進行中的活動
// @Tags Activity
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /activity [get]
func ActivityJSON(ctx *gin.Context) {
}

// SetDisplayFunc(func(model types.FieldModel) interface{} {
// 	if model.ID == "" {
// 		return "請輸入時間, 請輸入時間"
// 	}
// 	res, _ := s.table(config.ACTIVITY_TABLE).Select("start_time", "end_time").Find("activity_id", model.ID)
// 	start := strings.Replace(fmt.Sprintf("%s", res["start_time"]), " ", "T", 1)
// 	end := strings.Replace(fmt.Sprintf("%s", res["end_time"]), " ", "T", 1)
// 	return fmt.Sprintf("%s,%s", start[:len(start)-3], end[:len(start)-3])
// })

// SetDisplayFunc(func(model types.FieldModel) interface{} {
// 	if model.ID == "" {
// 		return nil
// 	}
// 	res, _ := s.table(config.ACTIVITY_TABLE).Select("city", "town").
// 		Where("activity_id", "=", model.ID).First()
// 	return fmt.Sprintf("%s,%s", res["city"], res["town"])
// })

// tables := []string{config.ACTIVITY_TABLE, config.ACTIVITY_OVERVIEW_TABLE, config.ACTIVITY_INTRODUCE_TABLE,
// 	config.ACTIVITY_SCHEDULE_TABLE, config.ACTIVITY_GUEST_TABLE, config.ACTIVITY_MATERIAL_TABLE,
// 	config.ACTIVITY_QUESTION_GUEST_TABLE, config.ACTIVITY_LOTTERY_TABLE, config.ACTIVITY_LOTTERY_PRIZE_TABLE,
// 	config.ACTIVITY_REDPACK_TABLE, config.ACTIVITY_REDPACK_PRIZE_TABLE,
// 	config.ACTIVITY_ROPEPACK_TABLE, config.ACTIVITY_ROPEPACK_PRIZE_TABLE, config.ACTIVITY_STAFF_GAME_TABLE,
// 	config.ACTIVITY_STAFF_PRIZE_TABLE, config.ACTIVITY_GAME_TABLE, config.ACTIVITY_PRIZE_TABLE,
// 	config.ACTIVITY_APPLYSIGN_TABLE}
// config.ACTIVITY_LOTTERY_TABLE, config.ACTIVITY_LOTTERY_PRIZE_TABLE,
// config.ACTIVITY_REDPACK_TABLE, config.ACTIVITY_REDPACK_PRIZE_TABLE,
// config.ACTIVITY_ROPEPACK_TABLE, config.ACTIVITY_ROPEPACK_PRIZE_TABLE, config.ACTIVITY_STAFF_GAME_TABLE,
// config.ACTIVITY_DRAW_NUMBERS_PRIZE_TABLE,
// config.ACTIVITY_WHACK_MOLE_TABLE,config.ACTIVITY_WHACK_MOLE_PRIZE_TABLE,

// 加入所有相關資料表
// ---------------下次修改活動總覽時，刪除下半部
// for i := 1; i < 15; i++ {
// 	models.DefaultOverviewModel().SetDbConn(s.conn).Add(models.NewOverviewModel{
// 		ActivityID: activityID,
// 		OverviewID: strconv.Itoa(i),
// 		Open:       "open",
// 	})
// }

// models.DefaultIntroduceModel().SetDbConn(s.conn).Add(models.EditIntroduceModel{
// 	ActivityID: activityID, Title: "簡介", IntroduceType: "文字",
// 	Content: "活動介紹", IntroduceOrder: "1",
// })
// models.DefaultIntroduceSettingModel().SetDbConn(s.conn).Add(activityID, "活動介紹")
// models.DefaultScheduleModel().SetDbConn(s.conn).Add(activityID, "簽到", "簽到", values.Get("start_time"), "00:00:00", "00:00:00")
// models.DefaultScheduleSettingModel().SetDbConn(s.conn).Add(activityID, "行程", "open", "open")
// models.DefaultGuestSettingModel().SetDbConn(s.conn).Add(activityID, "嘉賓", UPLOAD_SYSTEM_URL+"img-guest-bg.png")
// models.DefaultMessageModel().SetDbConn(s.conn).Add(activityID, "close", "close", "close", 1, 20, "open", "歡迎來到活動聊天室!", 1)
// models.DefaultTopicModel().SetDbConn(s.conn).Add(activityID, UPLOAD_SYSTEM_URL+"img-topic-bg.png")
// models.DefaultQuestionModel().SetDbConn(s.conn).Add(activityID, "close", "close", "close", "open", UPLOAD_SYSTEM_URL+"img-request-bg.png", "close")
// models.DefaultDanmuModel().SetDbConn(s.conn).Add(activityID, "close", "open", "open", "open", "open", "open", 100, 100, 100, 100, 1)
// models.DefaultSpecialdanmuModel().SetDbConn(s.conn).Add(activityID, "close", 0, 0, 0, "open")
// models.DefaultPictureModel().SetDbConn(s.conn).Add(activityID, values.Get("start_time"), values.Get("end_time"), "0", "3", "1", "system/img-pic.png", UPLOAD_SYSTEM_URL+"img-picture-bg.png", "1")
// models.DefaultHoldScreenModel().SetDbConn(s.conn).Add(activityID, 0, "close", "close", 10, "open")
// models.DefaultGeneralModel().SetDbConn(s.conn).Add(activityID, "open", 1, UPLOAD_SYSTEM_URL+"img-sign-bg.png")
// models.DefaultThreeDModel().SetDbConn(s.conn).Add(activityID, UPLOAD_SYSTEM_URL+"img-sign3d-headpic.png", "circle", "open", 1, UPLOAD_SYSTEM_URL+"img-sign3d-bg.png", UPLOAD_SYSTEM_URL+"img-sign3d-design.png", 50)
// models.DefaultCountdownModel().SetDbConn(s.conn).Add(activityID, "current", UPLOAD_SYSTEM_URL+"img-countdown-headpic.png", "circle", 1, 5)
// models.DefaultApplyModel().SetDbConn(s.conn).Add(activityID, "open")
// models.DefaultSignModel().SetDbConn(s.conn).Add(activityID, "open", "close", "30", "open")
// models.DefaultMaterialSettingMod

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+idStr, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS+idStr, "修改資料")
// s.redisConn.Publish(config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS+idStr, "修改資料")
// s.redisConn.Publish(config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS+idStr, "修改資料")
// s.redisConn.Publish(config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS+idStr, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_SIGNNAME_EDIT_TIMES_REDIS+idStr, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GENERAL_EDIT_TIMES_REDIS+idStr, "修改資料")
// s.redisConn.Publish(config.CHANNEL_THREED_EDIT_TIMES_REDIS+idStr, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_1", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_2", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_3", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_4", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_5", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_6", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_7", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_8", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_9", "修改資料")
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_REDIS+idStr+"_channel_10", "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+idStr, "修改資料")

// 刪除活動資料夾
// os.RemoveAll("uploads/" + userID + "/" + idStr)

// 刪除redis活動相關資訊
// s.redisConn.DelCache(config.ACTIVITY_REDIS + idStr)        // 活動資訊
// s.redisConn.DelCache(config.ACTIVITY_NUMBER_REDIS + idStr) // 活動number
// s.redisConn.DelCache(config.SIGN_STAFFS_1_REDIS + idStr) // 簽到人員資訊
// s.redisConn.DelCache(config.SIGN_STAFFS_2_REDIS + idStr) // 簽到人員資訊
// s.redisConn.DelCache(config.QUESTION_REDIS + idStr)               // 提問資訊(所有資訊)
// s.redisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + idStr) // 提問資訊(時間)
// s.redisConn.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + idStr)    // 提問資訊(讚數)
// s.redisConn.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + idStr) // 提問資訊(用戶按讚紀錄)
// s.redisConn.DelCache(config.CHATROOM_REDIS + idStr)                   // 聊天紀錄資訊
// s.redisConn.DelCache(config.CHATROOM_ORDER_BY_TIME_REDIS + idStr)     // 聊天紀錄資訊(時間)

// 頻道
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_1")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_2")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_3")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_4")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_5")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_6")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_7")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_8")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_9")  // 主持端遙控資訊，HASH
// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + activityID + "_channel_10") // 主持端遙控資訊，HASH

// s.redisConn.DelCache(config.HOST_CONTROL_REDIS + idStr)          // 主持端遙控資訊
// s.redisConn.DelCache(config.HOST_CONTROL_CHANNEL_REDIS + idStr)        // 主持端所有可遙控的session
// s.redisConn.DelCache(config.BLACK_STAFFS_ACTIVITY_REDIS + idStr)       // 黑名單人員(活動)
// s.redisConn.DelCache(config.BLACK_STAFFS_MESSAGE_REDIS + idStr)        // 黑名單人員(訊息)
// s.redisConn.DelCache(config.BLACK_STAFFS_QUESTION_REDIS + idStr)       // 黑名單人員(提問)
// s.redisConn.DelCache(config.BLACK_STAFFS_SIGNNAME_REDIS + idStr)       // 黑名單人員(簽名)
// s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + idStr)           // 玩家遊戲紀錄(中獎.未中獎)
// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + idStr) // 活動

// log.Println("新增api: ", model.ActivityID)
// log.Println("model.UserID: ", model.UserID)
// log.Println("model.ActivityName: ", model.ActivityName)
// log.Println("model.People: ", model.People)
// log.Println("model.StartTime: ", model.StartTime)

// models.EditActivityModel{
// 	ActivityID:       activityID,
// 	UserID:           values.Get("user"),
// 	ActivityName:     values.Get("activity_name"),
// 	ActivityType:     values.Get("activity_type"),
// 	MaxPeople:        values.Get("max_people"),
// 	People:           values.Get("people"),
// 	City:             values.Get("city"),
// 	Town:             values.Get("town"),
// 	StartTime:        values.Get("start_time"),
// 	EndTime:          values.Get("end_time"),
// 	LoginRequired:    values.Get("login_required"),
// 	LoginPassword:    values.Get("login_password"),
// 	PasswordRequired: values.Get("password_required"),
// 	Permissions:      values.Get("permissions"),
// 	MessageAmount:    values.Get("message_amount"),
// 	PushPhoneMessage: values.Get("push_phone_message"),

// 	MailAmount: values.Get("mail_amount"),
// 	SendMail:   values.Get("send_mail"),

// 	HostScan: values.Get("host_scan"),

// 	// 官方帳號綁定
// 	LineID:        values.Get("line_id"),
// 	ChannelID:     values.Get("channel_id"),
// 	ChannelSecret: values.Get("channel_secret"),
// 	ChatbotSecret: values.Get("chatbot_secret"),
// 	ChatbotToken:  values.Get("chatbot_token"),

// 	PushMessage: values.Get("push_message"),
// 	Device:      values.Get("device"),

// 	ChannelAmount: values.Get("channel_amount"),
// }
// log.Println("編輯api: ", model.ActivityID)
// log.Println("model.UserID: ", model.UserID)
// log.Println("model.ActivityName: ", model.ActivityName)
// log.Println("model.People: ", model.People)
// log.Println("model.StartTime: ", model.StartTime)

// models.EditActivityModel{
// 	ActivityID:       values.Get("activity_id"),
// 	UserID:           values.Get("user"),
// 	ActivityName:     values.Get("activity_name"),
// 	ActivityType:     values.Get("activity_type"),
// 	MaxPeople:        values.Get("max_people"),
// 	People:           values.Get("people"),
// 	City:             values.Get("city"),
// 	Town:             values.Get("town"),
// 	StartTime:        values.Get("start_time"),
// 	EndTime:          values.Get("end_time"),
// 	LoginRequired:    values.Get("login_required"),
// 	LoginPassword:    values.Get("login_password"),
// 	PasswordRequired: values.Get("password_required"),
// 	Permissions:      values.Get("permissions"),
// 	MessageAmount:    values.Get("message_amount"),
// 	PushPhoneMessage: values.Get("push_phone_message"),

// 	MailAmount: values.Get("mail_amount"),
// 	SendMail:   values.Get("send_mail"),

// 	HostScan: values.Get("host_scan"),

// 	// 官方帳號綁定
// 	LineID:        values.Get("line_id"),
// 	ChannelID:     values.Get("channel_id"),
// 	ChannelSecret: values.Get("channel_secret"),
// 	ChatbotSecret: values.Get("chatbot_secret"),
// 	ChatbotToken:  values.Get("chatbot_token"),

// 	PushMessage: values.Get("push_message"),

// 	Device: values.Get("device"),

// 	ChannelAmount: values.Get("channel_amount"),
// }
