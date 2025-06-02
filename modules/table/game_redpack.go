package table

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	redpackPictureFields = []PictureField{
		{FieldName: "redpack_bgm_start", Path: "redpack/%s/bgm/start.mp3"},
		{FieldName: "redpack_bgm_gaming", Path: "redpack/%s/bgm/gaming.mp3"},
		{FieldName: "redpack_bgm_end", Path: "redpack/%s/bgm/end.mp3"},

		{FieldName: "redpack_classic_h_pic_01", Path: "redpack/classic/redpack_classic_h_pic_01.png"},
		{FieldName: "redpack_classic_h_pic_02", Path: "redpack/classic/redpack_classic_h_pic_02.jpg"},
		{FieldName: "redpack_classic_h_pic_03", Path: "redpack/classic/redpack_classic_h_pic_03.png"},
		{FieldName: "redpack_classic_h_pic_04", Path: "redpack/classic/redpack_classic_h_pic_04.png"},
		{FieldName: "redpack_classic_h_pic_05", Path: "redpack/classic/redpack_classic_h_pic_05.png"},
		{FieldName: "redpack_classic_h_pic_06", Path: "redpack/classic/redpack_classic_h_pic_06.png"},
		{FieldName: "redpack_classic_h_pic_07", Path: "redpack/classic/redpack_classic_h_pic_07.png"},
		{FieldName: "redpack_classic_h_pic_08", Path: "redpack/classic/redpack_classic_h_pic_08.png"},
		{FieldName: "redpack_classic_h_pic_09", Path: "redpack/classic/redpack_classic_h_pic_09.png"},
		{FieldName: "redpack_classic_h_pic_10", Path: "redpack/classic/redpack_classic_h_pic_10.png"},
		{FieldName: "redpack_classic_h_pic_11", Path: "redpack/classic/redpack_classic_h_pic_11.png"},
		{FieldName: "redpack_classic_h_pic_12", Path: "redpack/classic/redpack_classic_h_pic_12.png"},
		{FieldName: "redpack_classic_h_pic_13", Path: "redpack/classic/redpack_classic_h_pic_13.jpg"},
		{FieldName: "redpack_classic_g_pic_01", Path: "redpack/classic/redpack_classic_g_pic_01.png"},
		{FieldName: "redpack_classic_g_pic_02", Path: "redpack/classic/redpack_classic_g_pic_02.jpg"},
		{FieldName: "redpack_classic_g_pic_03", Path: "redpack/classic/redpack_classic_g_pic_03.png"},
		{FieldName: "redpack_classic_h_ani_01", Path: "redpack/classic/redpack_classic_h_ani_01.png"},
		{FieldName: "redpack_classic_h_ani_02", Path: "redpack/classic/redpack_classic_h_ani_02.png"},
		{FieldName: "redpack_classic_g_ani_01", Path: "redpack/classic/redpack_classic_g_ani_01.png"},
		{FieldName: "redpack_classic_g_ani_02", Path: "redpack/classic/redpack_classic_g_ani_02.png"},
		{FieldName: "redpack_classic_g_ani_03", Path: "redpack/classic/redpack_classic_g_ani_03.png"},

		{FieldName: "redpack_cherry_h_pic_01", Path: "redpack/cherry/redpack_cherry_h_pic_01.png"},
		{FieldName: "redpack_cherry_h_pic_02", Path: "redpack/cherry/redpack_cherry_h_pic_02.png"},
		{FieldName: "redpack_cherry_h_pic_03", Path: "redpack/cherry/redpack_cherry_h_pic_03.png"},
		{FieldName: "redpack_cherry_h_pic_04", Path: "redpack/cherry/redpack_cherry_h_pic_04.png"},
		{FieldName: "redpack_cherry_h_pic_05", Path: "redpack/cherry/redpack_cherry_h_pic_05.png"},
		{FieldName: "redpack_cherry_h_pic_06", Path: "redpack/cherry/redpack_cherry_h_pic_06.png"},
		{FieldName: "redpack_cherry_h_pic_07", Path: "redpack/cherry/redpack_cherry_h_pic_07.png"},
		{FieldName: "redpack_cherry_g_pic_01", Path: "redpack/cherry/redpack_cherry_g_pic_01.png"},
		{FieldName: "redpack_cherry_g_pic_02", Path: "redpack/cherry/redpack_cherry_g_pic_02.jpg"},
		{FieldName: "redpack_cherry_h_ani_01", Path: "redpack/cherry/redpack_cherry_h_ani_01.png"},
		{FieldName: "redpack_cherry_h_ani_02", Path: "redpack/cherry/redpack_cherry_h_ani_02.png"},
		{FieldName: "redpack_cherry_g_ani_01", Path: "redpack/cherry/redpack_cherry_g_ani_01.png"},
		{FieldName: "redpack_cherry_g_ani_02", Path: "redpack/cherry/redpack_cherry_g_ani_02.png"},

		{FieldName: "redpack_christmas_h_pic_01", Path: "redpack/christmas/redpack_christmas_h_pic_01.png"},
		{FieldName: "redpack_christmas_h_pic_02", Path: "redpack/christmas/redpack_christmas_h_pic_02.jpg"},
		{FieldName: "redpack_christmas_h_pic_03", Path: "redpack/christmas/redpack_christmas_h_pic_03.png"},
		{FieldName: "redpack_christmas_h_pic_04", Path: "redpack/christmas/redpack_christmas_h_pic_04.png"},
		{FieldName: "redpack_christmas_h_pic_05", Path: "redpack/christmas/redpack_christmas_h_pic_05.png"},
		{FieldName: "redpack_christmas_h_pic_06", Path: "redpack/christmas/redpack_christmas_h_pic_06.png"},
		{FieldName: "redpack_christmas_h_pic_07", Path: "redpack/christmas/redpack_christmas_h_pic_07.png"},
		{FieldName: "redpack_christmas_h_pic_08", Path: "redpack/christmas/redpack_christmas_h_pic_08.png"},
		{FieldName: "redpack_christmas_h_pic_09", Path: "redpack/christmas/redpack_christmas_h_pic_09.png"},
		{FieldName: "redpack_christmas_h_pic_10", Path: "redpack/christmas/redpack_christmas_h_pic_10.png"},
		{FieldName: "redpack_christmas_h_pic_11", Path: "redpack/christmas/redpack_christmas_h_pic_11.png"},
		{FieldName: "redpack_christmas_h_pic_12", Path: "redpack/christmas/redpack_christmas_h_pic_12.png"},
		{FieldName: "redpack_christmas_h_pic_13", Path: "redpack/christmas/redpack_christmas_h_pic_13.png"},
		{FieldName: "redpack_christmas_g_pic_01", Path: "redpack/christmas/redpack_christmas_g_pic_01.png"},
		{FieldName: "redpack_christmas_g_pic_02", Path: "redpack/christmas/redpack_christmas_g_pic_02.jpg"},
		{FieldName: "redpack_christmas_g_pic_03", Path: "redpack/christmas/redpack_christmas_g_pic_03.png"},
		{FieldName: "redpack_christmas_g_pic_04", Path: "redpack/christmas/redpack_christmas_g_pic_04.png"},
		{FieldName: "redpack_christmas_c_pic_01", Path: "redpack/christmas/redpack_christmas_c_pic_01.png"},
		{FieldName: "redpack_christmas_c_pic_02", Path: "redpack/christmas/redpack_christmas_c_pic_02.png"},
		{FieldName: "redpack_christmas_c_ani_01", Path: "redpack/christmas/redpack_christmas_c_ani_01.png"},
	}
)

// GetRedpackPanel 搖紅包
func (s *SystemTable) GetRedpackPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)

		// 刪除資料表
		tablesToDelete := []string{
			config.ACTIVITY_PRIZE_TABLE,
			config.ACTIVITY_STAFF_GAME_TABLE,
			config.ACTIVITY_STAFF_PRIZE_TABLE,
			config.ACTIVITY_STAFF_BLACK_TABLE,
			config.ACTIVITY_STAFF_PK_TABLE,
			config.ACTIVITY_SCORE_TABLE,
			config.ACTIVITY_GAME_QA_RECORD_TABLE,

			// 投票
			config.ACTIVITY_GAME_VOTE_OPTION_TABLE,
			config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE,
			config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE,
			config.ACTIVITY_GAME_VOTE_RECORD_TABLE,

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
			// config.ACTIVITY_GAME_QA_TABLE,
			// config.ACTIVITY_GAME_VOTE_PICTURE_TABLE,
			// config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE,
		}

		for _, table := range tablesToDelete {
			s.table(table).WhereIn("game_id", ids).Delete()
		}

		// mongo
		mongoTables := []string{
			config.ACTIVITY_GAME_TABLE,
		}
		for _, t := range mongoTables {
			s.mongoConn.DeleteMany(t, bson.M{"game_id": bson.M{"$in": ids}})
		}

		for _, id := range idArr {
			// Redis 要刪除的 key 前綴列表
			delKeys := []string{
				config.GAME_REDIS,
				config.GAME_TYPE_REDIS, // 遊戲種類
				config.GAME_PRIZES_AMOUNT_REDIS,
				config.BLACK_STAFFS_GAME_REDIS,
				config.SCORES_REDIS,
				config.SCORES_2_REDIS,
				config.CORRECT_REDIS,
				config.QA_REDIS,
				config.QA_RECORD_REDIS,
				config.WINNING_STAFFS_REDIS,
				config.NO_WINNING_STAFFS_REDIS, // 未中獎人員
				config.GAME_TEAM_REDIS,
				config.GAME_BINGO_NUMBER_REDIS,               // 紀錄抽過的號碼，LIST
				config.GAME_BINGO_USER_REDIS,                 // 賓果中獎人員，ZSET
				config.GAME_BINGO_USER_NUMBER,                // 紀錄玩家的號碼排序，HASH
				config.GAME_BINGO_USER_GOING_BINGO,           // 紀錄玩家是否即將中獎，HASH
				config.GAME_ATTEND_REDIS,                     // 遊戲人數資訊，SET
				config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS,  // 拔河遊戲左隊人數資訊，SET
				config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS, // 拔河遊戲右隊人數資訊，SET
				// 投票
				config.GAME_VOTE_RECORDS_REDIS,
				config.VOTE_SPECIAL_OFFICER_REDIS,
			}

			for _, key := range delKeys {
				s.redisConn.DelCache(key + id)
			}

			// Redis 要發佈的頻道前綴列表
			pubChannels := []string{
				config.CHANNEL_GAME_REDIS,
				config.CHANNEL_GUEST_GAME_STATUS_REDIS,
				config.CHANNEL_GAME_BINGO_NUMBER_REDIS,
				config.CHANNEL_QA_REDIS,
				config.CHANNEL_GAME_TEAM_REDIS,
				config.CHANNEL_BLACK_STAFFS_GAME_REDIS,
				config.CHANNEL_GAME_EDIT_TIMES_REDIS,
				config.CHANNEL_WINNING_STAFFS_REDIS,
				config.CHANNEL_GAME_BINGO_USER_NUMBER,
				config.CHANNEL_SCORES_REDIS,
			}

			for _, channel := range pubChannels {
				s.redisConn.Publish(channel+id, "修改資料")
			}

			// 刪除遊戲資料夾
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/redpack/" + id)
		}

		// 刪除遊戲場次時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

		// 刪除玩家遊戲紀錄(中獎.未中獎)
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)
		return nil
	})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGameModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(redpackPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "redpack", values.Get("game_id"), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id") {
			return errors.New("錯誤: 遊戲ID發生問題，請輸入有效的遊戲ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGameModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(redpackPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "redpack", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增搖紅包遊戲資料(form-data)
// @Tags Redpack
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param max_people formData integer true "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param limit_time formData string true "是否限時" Enums(open, close)
// @param second formData integer true "限時秒數"
// @param percent formData integer true "中獎機率(0-100%)" minimum(0) maximum(100)
// @param topic formData string true "主題樣式" Enums(01_classic, 02_cherry)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/form [post]
func POSTRedpack(ctx *gin.Context) {
}

// @Summary 新增搖紅包獎品資料(form-data)
// @Tags Redpack Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_name formData string true "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string true "類型" Enums(first, second, third, general)
// @param prize_picture formData file false "照片"
// @param prize_method formData string true "兌獎方式" Enums(site, mail)
// @param prize_password formData string true "兌獎密碼(最多八個字元)" minlength(1) maxlength(8)
// @param prize_amount formData integer true "數量"
// @param prize_price formData integer true "價值"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/prize/form [post]
func POSTRedpackPrize(ctx *gin.Context) {
}

// @Summary 編輯搖紅包遊戲資料(form-data)
// @Tags Redpack
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param max_people formData integer false "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string false "允許重複中獎" Enums(open, close)
// @param limit_time formData string false "是否限時" Enums(open, close)
// @param second formData integer false "限時秒數"
// @param percent formData integer false "中獎機率(0-100%)" minimum(0) maximum(100)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_cherry)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/form [put]
func PUTRedpack(ctx *gin.Context) {
}

// @Summary 編輯搖紅包獎品資料(form-data)
// @Tags Redpack Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_id formData string true "獎品ID"
// @param prize_name formData string false "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string false "類型" Enums(first, second, third, general)
// @param prize_picture formData file false "照片"
// @param prize_method formData string false "兌獎方式" Enums(site, mail)
// @param prize_password formData string false "兌獎密碼(最多八個字元)" minlength(0) maxlength(8)
// @param prize_amount formData integer false "數量(同時更新剩餘數量)"
// @param prize_price formData integer false "價值"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/prize/form [put]
func PUTRedpackPrize(ctx *gin.Context) {
}

// @Summary 刪除搖紅包遊戲資料(form-data)
// @Tags Redpack
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/form [delete]
func DELETERedpack(ctx *gin.Context) {
}

// @Summary 刪除搖紅包獎品資料(form-data)
// @Tags Redpack Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/prize/form [delete]
func DELETERedpackPrize(ctx *gin.Context) {
}

// @Summary 搖紅包遊戲JSON資料
// @Tags Redpack
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack [get]
func RedpackJSON(ctx *gin.Context) {
}

// @Summary 搖紅包獎品JSON資料
// @Tags Redpack Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/redpack/prize [get]
func RedpackPrizeJSON(ctx *gin.Context) {
}

// models.NewGameModel{
// 	UserID:        values.Get("user"),
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        values.Get("second"),
// 	MaxPeople:     values.Get("max_people"),
// 	People:        values.Get("people"),
// 	MaxTimes:      "0",
// 	Allow:         values.Get("allow"),
// 	Percent:       values.Get("percent"),
// 	FirstPrize:    "0",
// 	SecondPrize:   "0",
// 	ThirdPrize:    "0",
// 	GeneralPrize:  "0",
// 	Topic:         values.Get("topic"),
// 	Skin:          values.Get("skin"),
// 	Music:         values.Get("music"),
// 	DisplayName:   "open",
// 	BoxReflection: "",
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            "",

// 	// 賓果遊戲
// 	MaxNumber:  "0",
// 	BingoLine:  "0",
// 	RoundPrize: "0",

// 	// 扭蛋機遊戲
// 	GachaMachineReflection: "0",
// 	ReflectiveSwitch:       "open",

// 	// 投票遊戲
// 	VoteScreen:       "",
// 	VoteTimes:        "0",
// 	VoteMethod:       "",
// 	VoteMethodPlayer: "",
// 	VoteRestriction:  "",
// 	AvatarShape:      "",
// 	VoteStartTime:    "",
// 	VoteEndTime:      "",
// 	AutoPlay:         "",
// 	ShowRank:         "",
// 	TitleSwitch:      "",
// 	ArrangementGuest: "",

// 	// 搖紅包自定義
// 	// 音樂
// 	RedpackBgmStart:  update[0],
// 	RedpackBgmGaming: update[1],
// 	RedpackBgmEnd:    update[2],

// 	RedpackClassicHPic01: update[3],
// 	RedpackClassicHPic02: update[4],
// 	RedpackClassicHPic03: update[5],
// 	RedpackClassicHPic04: update[6],
// 	RedpackClassicHPic05: update[7],
// 	RedpackClassicHPic06: update[8],
// 	RedpackClassicHPic07: update[9],
// 	RedpackClassicHPic08: update[10],
// 	RedpackClassicHPic09: update[11],
// 	RedpackClassicHPic10: update[12],
// 	RedpackClassicHPic11: update[13],
// 	RedpackClassicHPic12: update[14],
// 	RedpackClassicHPic13: update[15],
// 	RedpackClassicGPic01: update[16],
// 	RedpackClassicGPic02: update[17],
// 	RedpackClassicGPic03: update[18],
// 	RedpackClassicHAni01: update[19],
// 	RedpackClassicHAni02: update[20],
// 	RedpackClassicGAni01: update[21],
// 	RedpackClassicGAni02: update[22],
// 	RedpackClassicGAni03: update[23],

// 	RedpackCherryHPic01: update[24],
// 	RedpackCherryHPic02: update[25],
// 	RedpackCherryHPic03: update[26],
// 	RedpackCherryHPic04: update[27],
// 	RedpackCherryHPic05: update[28],
// 	RedpackCherryHPic06: update[29],
// 	RedpackCherryHPic07: update[30],
// 	RedpackCherryGPic01: update[31],
// 	RedpackCherryGPic02: update[32],
// 	RedpackCherryHAni01: update[33],
// 	RedpackCherryHAni02: update[34],
// 	RedpackCherryGAni01: update[35],
// 	RedpackCherryGAni02: update[36],

// 	RedpackChristmasHPic01: update[37],
// 	RedpackChristmasHPic02: update[38],
// 	RedpackChristmasHPic03: update[39],
// 	RedpackChristmasHPic04: update[40],
// 	RedpackChristmasHPic05: update[41],
// 	RedpackChristmasHPic06: update[42],
// 	RedpackChristmasHPic07: update[43],
// 	RedpackChristmasHPic08: update[44],
// 	RedpackChristmasHPic09: update[45],
// 	RedpackChristmasHPic10: update[46],
// 	RedpackChristmasHPic11: update[47],
// 	RedpackChristmasHPic12: update[48],
// 	RedpackChristmasHPic13: update[49],
// 	RedpackChristmasGPic01: update[50],
// 	RedpackChristmasGPic02: update[51],
// 	RedpackChristmasGPic03: update[52],
// 	RedpackChristmasGPic04: update[53],
// 	RedpackChristmasCPic01: update[

// pics = []string{
// 搖紅包自定義
// "redpack/%s/bgm/start.mp3",
// "redpack/%s/bgm/gaming.mp3",
// "redpack/%s/bgm/end.mp3",

// "redpack/classic/redpack_classic_h_pic_01.png",
// "redpack/classic/redpack_classic_h_pic_02.jpg",
// "redpack/classic/redpack_classic_h_pic_03.png",
// "redpack/classic/redpack_classic_h_pic_04.png",
// "redpack/classic/redpack_classic_h_pic_05.png",
// "redpack/classic/redpack_classic_h_pic_06.png",
// "redpack/classic/redpack_classic_h_pic_07.png",
// "redpack/classic/redpack_classic_h_pic_08.png",
// "redpack/classic/redpack_classic_h_pic_09.png",
// "redpack/classic/redpack_classic_h_pic_10.png",
// "redpack/classic/redpack_classic_h_pic_11.png",
// "redpack/classic/redpack_classic_h_pic_12.png",
// "redpack/classic/redpack_classic_h_pic_13.jpg",
// "redpack/classic/redpack_classic_g_pic_01.png",
// "redpack/classic/redpack_classic_g_pic_02.jpg",
// "redpack/classic/redpack_classic_g_pic_03.png",
// "redpack/classic/redpack_classic_h_ani_01.png",
// "redpack/classic/redpack_classic_h_ani_02.png",
// "redpack/classic/redpack_classic_g_ani_01.png",
// "redpack/classic/redpack_classic_g_ani_02.png",
// "redpack/classic/redpack_classic_g_ani_03.png",

// "redpack/cherry/redpack_cherry_h_pic_01.png",
// "redpack/cherry/redpack_cherry_h_pic_02.png",
// "redpack/cherry/redpack_cherry_h_pic_03.png",
// "redpack/cherry/redpack_cherry_h_pic_04.png",
// "redpack/cherry/redpack_cherry_h_pic_05.png",
// "redpack/cherry/redpack_cherry_h_pic_06.png",
// "redpack/cherry/redpack_cherry_h_pic_07.png",
// "redpack/cherry/redpack_cherry_g_pic_01.png",
// "redpack/cherry/redpack_cherry_g_pic_02.jpg",
// "redpack/cherry/redpack_cherry_h_ani_01.png",
// "redpack/cherry/redpack_cherry_h_ani_02.png",
// "redpack/cherry/redpack_cherry_g_ani_01.png",
// "redpack/cherry/redpack_cherry_g_ani_02.png",

// "redpack/christmas/redpack_christmas_h_pic_01.png",
// "redpack/christmas/redpack_christmas_h_pic_02.jpg",
// "redpack/christmas/redpack_christmas_h_pic_03.png",
// "redpack/christmas/redpack_christmas_h_pic_04.png",
// "redpack/christmas/redpack_christmas_h_pic_05.png",
// "redpack/christmas/redpack_christmas_h_pic_06.png",
// "redpack/christmas/redpack_christmas_h_pic_07.png",
// "redpack/christmas/redpack_christmas_h_pic_08.png",
// "redpack/christmas/redpack_christmas_h_pic_09.png",
// "redpack/christmas/redpack_christmas_h_pic_10.png",
// "redpack/christmas/redpack_christmas_h_pic_11.png",
// "redpack/christmas/redpack_christmas_h_pic_12.png",
// "redpack/christmas/redpack_christmas_h_pic_13.png",
// "redpack/christmas/redpack_christmas_g_pic_01.png",
// "redpack/christmas/redpack_christmas_g_pic_02.jpg",
// "redpack/christmas/redpack_christmas_g_pic_03.png",
// "redpack/christmas/redpack_christmas_g_pic_04.png",
// "redpack/christmas/redpack_christmas_c_pic_01.png",
// "redpack/christmas/redpack_christmas_c_pic_02.png",
// "redpack/christmas/redpack_christmas_c_ani_01.png",
// }

// fields = []string{
// 搖紅包自定義
// "redpack_bgm_start",
// "redpack_bgm_gaming",
// "redpack_bgm_end",

// "redpack_classic_h_pic_01",
// "redpack_classic_h_pic_02",
// "redpack_classic_h_pic_03",
// "redpack_classic_h_pic_04",
// "redpack_classic_h_pic_05",
// "redpack_classic_h_pic_06",
// "redpack_classic_h_pic_07",
// "redpack_classic_h_pic_08",
// "redpack_classic_h_pic_09",
// "redpack_classic_h_pic_10",
// "redpack_classic_h_pic_11",
// "redpack_classic_h_pic_12",
// "redpack_classic_h_pic_13",
// "redpack_classic_g_pic_01",
// "redpack_classic_g_pic_02",
// "redpack_classic_g_pic_03",
// "redpack_classic_h_ani_01",
// "redpack_classic_h_ani_02",
// "redpack_classic_g_ani_01",
// "redpack_classic_g_ani_02",
// "redpack_classic_g_ani_03",

// "redpack_cherry_h_pic_01",
// "redpack_cherry_h_pic_02",
// "redpack_cherry_h_pic_03",
// "redpack_cherry_h_pic_04",
// "redpack_cherry_h_pic_05",
// "redpack_cherry_h_pic_06",
// "redpack_cherry_h_pic_07",
// "redpack_cherry_g_pic_01",
// "redpack_cherry_g_pic_02",
// "redpack_cherry_h_ani_01",
// "redpack_cherry_h_ani_02",
// "redpack_cherry_g_ani_01",
// "redpack_cherry_g_ani_02",

// "redpack_christmas_h_pic_01",
// "redpack_christmas_h_pic_02",
// "redpack_christmas_h_pic_03",
// "redpack_christmas_h_pic_04",
// "redpack_christmas_h_pic_05",
// "redpack_christmas_h_pic_06",
// "redpack_christmas_h_pic_07",
// "redpack_christmas_h_pic_08",
// "redpack_christmas_h_pic_09",
// "redpack_christmas_h_pic_10",
// "redpack_christmas_h_pic_11",
// "redpack_christmas_h_pic_12",
// "redpack_christmas_h_pic_13",
// "redpack_christmas_g_pic_01",
// "redpack_christmas_g_pic_02",
// "redpack_christmas_g_pic_03",
// "redpack_christmas_g_pic_04",
// "redpack_christmas_c_pic_01",
// "redpack_christmas_c_pic_02",
// "redpack_christmas_c_ani_01",
// }
// update = make([]string, 200)

// for i, field := range fields {
// 	if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	}
// }

// var (
// 	pics = []string{
// 		// 搖紅包自定義
// 		"redpack/%s/bgm/start.mp3",
// 		"redpack/%s/bgm/gaming.mp3",
// 		"redpack/%s/bgm/end.mp3",

// 		"redpack/classic/redpack_classic_h_pic_01.png",
// 		"redpack/classic/redpack_classic_h_pic_02.jpg",
// 		"redpack/classic/redpack_classic_h_pic_03.png",
// 		"redpack/classic/redpack_classic_h_pic_04.png",
// 		"redpack/classic/redpack_classic_h_pic_05.png",
// 		"redpack/classic/redpack_classic_h_pic_06.png",
// 		"redpack/classic/redpack_classic_h_pic_07.png",
// 		"redpack/classic/redpack_classic_h_pic_08.png",
// 		"redpack/classic/redpack_classic_h_pic_09.png",
// 		"redpack/classic/redpack_classic_h_pic_10.png",
// 		"redpack/classic/redpack_classic_h_pic_11.png",
// 		"redpack/classic/redpack_classic_h_pic_12.png",
// 		"redpack/classic/redpack_classic_h_pic_13.jpg",
// 		"redpack/classic/redpack_classic_g_pic_01.png",
// 		"redpack/classic/redpack_classic_g_pic_02.jpg",
// 		"redpack/classic/redpack_classic_g_pic_03.png",
// 		"redpack/classic/redpack_classic_h_ani_01.png",
// 		"redpack/classic/redpack_classic_h_ani_02.png",
// 		"redpack/classic/redpack_classic_g_ani_01.png",
// 		"redpack/classic/redpack_classic_g_ani_02.png",
// 		"redpack/classic/redpack_classic_g_ani_03.png",

// 		"redpack/cherry/redpack_cherry_h_pic_01.png",
// 		"redpack/cherry/redpack_cherry_h_pic_02.png",
// 		"redpack/cherry/redpack_cherry_h_pic_03.png",
// 		"redpack/cherry/redpack_cherry_h_pic_04.png",
// 		"redpack/cherry/redpack_cherry_h_pic_05.png",
// 		"redpack/cherry/redpack_cherry_h_pic_06.png",
// 		"redpack/cherry/redpack_cherry_h_pic_07.png",
// 		"redpack/cherry/redpack_cherry_g_pic_01.png",
// 		"redpack/cherry/redpack_cherry_g_pic_02.jpg",
// 		"redpack/cherry/redpack_cherry_h_ani_01.png",
// 		"redpack/cherry/redpack_cherry_h_ani_02.png",
// 		"redpack/cherry/redpack_cherry_g_ani_01.png",
// 		"redpack/cherry/redpack_cherry_g_ani_02.png",

// 		"redpack/christmas/redpack_christmas_h_pic_01.png",
// 		"redpack/christmas/redpack_christmas_h_pic_02.jpg",
// 		"redpack/christmas/redpack_christmas_h_pic_03.png",
// 		"redpack/christmas/redpack_christmas_h_pic_04.png",
// 		"redpack/christmas/redpack_christmas_h_pic_05.png",
// 		"redpack/christmas/redpack_christmas_h_pic_06.png",
// 		"redpack/christmas/redpack_christmas_h_pic_07.png",
// 		"redpack/christmas/redpack_christmas_h_pic_08.png",
// 		"redpack/christmas/redpack_christmas_h_pic_09.png",
// 		"redpack/christmas/redpack_christmas_h_pic_10.png",
// 		"redpack/christmas/redpack_christmas_h_pic_11.png",
// 		"redpack/christmas/redpack_christmas_h_pic_12.png",
// 		"redpack/christmas/redpack_christmas_h_pic_13.png",
// 		"redpack/christmas/redpack_christmas_g_pic_01.png",
// 		"redpack/christmas/redpack_christmas_g_pic_02.jpg",
// 		"redpack/christmas/redpack_christmas_g_pic_03.png",
// 		"redpack/christmas/redpack_christmas_g_pic_04.png",
// 		"redpack/christmas/redpack_christmas_c_pic_01.png",
// 		"redpack/christmas/redpack_christmas_c_pic_02.png",
// 		"redpack/christmas/redpack_christmas_c_ani_01.png",
// 	}

// 	fields = []string{
// 		// 搖紅包自定義
// 		"redpack_bgm_start",
// 		"redpack_bgm_gaming",
// 		"redpack_bgm_end",

// 		"redpack_classic_h_pic_01",
// 		"redpack_classic_h_pic_02",
// 		"redpack_classic_h_pic_03",
// 		"redpack_classic_h_pic_04",
// 		"redpack_classic_h_pic_05",
// 		"redpack_classic_h_pic_06",
// 		"redpack_classic_h_pic_07",
// 		"redpack_classic_h_pic_08",
// 		"redpack_classic_h_pic_09",
// 		"redpack_classic_h_pic_10",
// 		"redpack_classic_h_pic_11",
// 		"redpack_classic_h_pic_12",
// 		"redpack_classic_h_pic_13",
// 		"redpack_classic_g_pic_01",
// 		"redpack_classic_g_pic_02",
// 		"redpack_classic_g_pic_03",
// 		"redpack_classic_h_ani_01",
// 		"redpack_classic_h_ani_02",
// 		"redpack_classic_g_ani_01",
// 		"redpack_classic_g_ani_02",
// 		"redpack_classic_g_ani_03",

// 		"redpack_cherry_h_pic_01",
// 		"redpack_cherry_h_pic_02",
// 		"redpack_cherry_h_pic_03",
// 		"redpack_cherry_h_pic_04",
// 		"redpack_cherry_h_pic_05",
// 		"redpack_cherry_h_pic_06",
// 		"redpack_cherry_h_pic_07",
// 		"redpack_cherry_g_pic_01",
// 		"redpack_cherry_g_pic_02",
// 		"redpack_cherry_h_ani_01",
// 		"redpack_cherry_h_ani_02",
// 		"redpack_cherry_g_ani_01",
// 		"redpack_cherry_g_ani_02",

// 		"redpack_christmas_h_pic_01",
// 		"redpack_christmas_h_pic_02",
// 		"redpack_christmas_h_pic_03",
// 		"redpack_christmas_h_pic_04",
// 		"redpack_christmas_h_pic_05",
// 		"redpack_christmas_h_pic_06",
// 		"redpack_christmas_h_pic_07",
// 		"redpack_christmas_h_pic_08",
// 		"redpack_christmas_h_pic_09",
// 		"redpack_christmas_h_pic_10",
// 		"redpack_christmas_h_pic_11",
// 		"redpack_christmas_h_pic_12",
// 		"redpack_christmas_h_pic_13",
// 		"redpack_christmas_g_pic_01",
// 		"redpack_christmas_g_pic_02",
// 		"redpack_christmas_g_pic_03",
// 		"redpack_christmas_g_pic_04",
// 		"redpack_christmas_c_pic_01",
// 		"redpack_christmas_c_pic_02",
// 		"redpack_christmas_c_ani_01",
// 	}
// 	update = make([]string, 200)
// )

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditGameModel{
// 	// UserID:           values.Get("user_id"),
// 	ActivityID:    values.Get("activity_id"),
// 	GameID:        values.Get("game_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        values.Get("second"),
// 	MaxPeople:     values.Get("max_people"),
// 	People:        values.Get("people"),
// 	MaxTimes:      "",
// 	Allow:         values.Get("allow"),
// 	Percent:       values.Get("percent"),
// 	FirstPrize:    "",
// 	SecondPrize:   "",
// 	ThirdPrize:    "",
// 	GeneralPrize:  "",
// 	Topic:         values.Get("topic"),
// 	Skin:          values.Get("skin"),
// 	Music:         values.Get("music"),
// 	DisplayName:   "",
// 	GameOrder:     values.Get("game_order"),
// 	BoxReflection: "",
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            "",

// 	// 賓果遊戲
// 	MaxNumber:  "",
// 	BingoLine:  "",
// 	RoundPrize: "",

// 	// 扭蛋機遊戲
// 	GachaMachineReflection: "",
// 	ReflectiveSwitch:       "",

// 	// 投票遊戲
// 	VoteScreen:       "",
// 	VoteTimes:        "",
// 	VoteMethod:       "",
// 	VoteMethodPlayer: "",
// 	VoteRestriction:  "",
// 	AvatarShape:      "",
// 	VoteStartTime:    "",
// 	VoteEndTime:      "",
// 	AutoPlay:         "",
// 	ShowRank:         "",
// 	TitleSwitch:      "",
// 	ArrangementGuest: "",

// 	// 搖紅包自定義
// 	// 音樂
// 	RedpackBgmStart:  update[0],
// 	RedpackBgmGaming: update[1],
// 	RedpackBgmEnd:    update[2],

// 	RedpackClassicHPic01: update[3],
// 	RedpackClassicHPic02: update[4],
// 	RedpackClassicHPic03: update[5],
// 	RedpackClassicHPic04: update[6],
// 	RedpackClassicHPic05: update[7],
// 	RedpackClassicHPic06: update[8],
// 	RedpackClassicHPic07: update[9],
// 	RedpackClassicHPic08: update[10],
// 	RedpackClassicHPic09: update[11],
// 	RedpackClassicHPic10: update[12],
// 	RedpackClassicHPic11: update[13],
// 	RedpackClassicHPic12: update[14],
// 	RedpackClassicHPic13: update[15],
// 	RedpackClassicGPic01: update[16],
// 	RedpackClassicGPic02: update[17],
// 	RedpackClassicGPic03: update[18],
// 	RedpackClassicHAni01: update[19],
// 	RedpackClassicHAni02: update[20],
// 	RedpackClassicGAni01: update[21],
// 	RedpackClassicGAni02: update[22],
// 	RedpackClassicGAni03: update[23],

// 	RedpackCherryHPic01: update[24],
// 	RedpackCherryHPic02: update[25],
// 	RedpackCherryHPic03: update[26],
// 	RedpackCherryHPic04: update[27],
// 	RedpackCherryHPic05: update[28],
// 	RedpackCherryHPic06: update[29],
// 	RedpackCherryHPic07: update[30],
// 	RedpackCherryGPic01: update[31],
// 	RedpackCherryGPic02: update[32],
// 	RedpackCherryHAni01: update[33],
// 	RedpackCherryHAni02: update[34],
// 	RedpackCherryGAni01: update[35],
// 	RedpackCherryGAni02: update[36],

// 	RedpackChristmasHPic01: update[37],
// 	RedpackChristmasHPic02: update[38],
// 	RedpackChristmasHPic03: update[39],
// 	RedpackChristmasHPic04: update[40],
// 	RedpackChristmasHPic05: update[41],
// 	RedpackChristmasHPic06: update[42],
// 	RedpackChristmasHPic07: update[43],
// 	RedpackChristmasHPic08: update[44],
// 	RedpackChristmasHPic09: update[45],
// 	RedpackChristmasHPic10: update[46],
// 	RedpackChristmasHPic11: update[47],
// 	RedpackChristmasHPic12: update[48],
// 	RedpackChristmasHPic13: update[49],
// 	RedpackChristmasGPic01: update[50],
// 	RedpackChristmasGPic02: update[51],
// 	RedpackChristmasGPic03: update[52],
// 	RedpackChristmasGPic04: update[53],
// 	RedpackChristmasCPic01: update[54],
// 	RedpackChristmasCPic02: update[55],
// 	RedpackChristmasCAni01: update[56],
//
