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
	lotteryPictureFields = []PictureField{
		{FieldName: "lottery_bgm_gaming", Path: "lottery/%s/bgm/gaming.mp3"},

		{FieldName: "lottery_jiugongge_classic_h_pic_01", Path: "lottery/classic/lottery_jiugongge_classic_h_pic_01.png"},
		{FieldName: "lottery_jiugongge_classic_h_pic_02", Path: "lottery/classic/lottery_jiugongge_classic_h_pic_02.png"},
		{FieldName: "lottery_jiugongge_classic_h_pic_03", Path: "lottery/classic/lottery_jiugongge_classic_h_pic_03.jpg"},
		{FieldName: "lottery_jiugongge_classic_h_pic_04", Path: "lottery/classic/lottery_jiugongge_classic_h_pic_04.png"},
		{FieldName: "lottery_jiugongge_classic_g_pic_01", Path: "lottery/classic/lottery_jiugongge_classic_g_pic_01.jpg"},
		{FieldName: "lottery_jiugongge_classic_g_pic_02", Path: "lottery/classic/lottery_jiugongge_classic_g_pic_02.png"},
		{FieldName: "lottery_jiugongge_classic_c_pic_01", Path: "lottery/classic/lottery_jiugongge_classic_c_pic_01.png"},
		{FieldName: "lottery_jiugongge_classic_c_pic_02", Path: "lottery/classic/lottery_jiugongge_classic_c_pic_02.png"},
		{FieldName: "lottery_jiugongge_classic_c_pic_03", Path: "lottery/classic/lottery_jiugongge_classic_c_pic_03.png"},
		{FieldName: "lottery_jiugongge_classic_c_pic_04", Path: "lottery/classic/lottery_jiugongge_classic_c_pic_04.png"},
		{FieldName: "lottery_jiugongge_classic_c_ani_01", Path: "lottery/classic/lottery_jiugongge_classic_c_ani_01.png"},
		{FieldName: "lottery_jiugongge_classic_c_ani_02", Path: "lottery/classic/lottery_jiugongge_classic_c_ani_02.png"},
		{FieldName: "lottery_jiugongge_classic_c_ani_03", Path: "lottery/classic/lottery_jiugongge_classic_c_ani_03.png"},

		{FieldName: "lottery_turntable_classic_h_pic_01", Path: "lottery/classic/lottery_turntable_classic_h_pic_01.png"},
		{FieldName: "lottery_turntable_classic_h_pic_02", Path: "lottery/classic/lottery_turntable_classic_h_pic_02.png"},
		{FieldName: "lottery_turntable_classic_h_pic_03", Path: "lottery/classic/lottery_turntable_classic_h_pic_03.jpg"},
		{FieldName: "lottery_turntable_classic_h_pic_04", Path: "lottery/classic/lottery_turntable_classic_h_pic_04.png"},
		{FieldName: "lottery_turntable_classic_g_pic_01", Path: "lottery/classic/lottery_turntable_classic_g_pic_01.jpg"},
		{FieldName: "lottery_turntable_classic_g_pic_02", Path: "lottery/classic/lottery_turntable_classic_g_pic_02.png"},
		{FieldName: "lottery_turntable_classic_c_pic_01", Path: "lottery/classic/lottery_turntable_classic_c_pic_01.png"},
		{FieldName: "lottery_turntable_classic_c_pic_02", Path: "lottery/classic/lottery_turntable_classic_c_pic_02.png"},
		{FieldName: "lottery_turntable_classic_c_pic_03", Path: "lottery/classic/lottery_turntable_classic_c_pic_03.png"},
		{FieldName: "lottery_turntable_classic_c_pic_04", Path: "lottery/classic/lottery_turntable_classic_c_pic_04.png"},
		{FieldName: "lottery_turntable_classic_c_pic_05", Path: "lottery/classic/lottery_turntable_classic_c_pic_05.png"},
		{FieldName: "lottery_turntable_classic_c_pic_06", Path: "lottery/classic/lottery_turntable_classic_c_pic_06.png"},
		{FieldName: "lottery_turntable_classic_c_ani_01", Path: "lottery/classic/lottery_turntable_classic_c_ani_01.png"},
		{FieldName: "lottery_turntable_classic_c_ani_02", Path: "lottery/classic/lottery_turntable_classic_c_ani_02.png"},
		{FieldName: "lottery_turntable_classic_c_ani_03", Path: "lottery/classic/lottery_turntable_classic_c_ani_03.png"},

		{FieldName: "lottery_jiugongge_starrysky_h_pic_01", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_01.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_02", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_02.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_03", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_03.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_04", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_04.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_05", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_05.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_06", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_06.png"},
		{FieldName: "lottery_jiugongge_starrysky_h_pic_07", Path: "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_07.jpg"},
		{FieldName: "lottery_jiugongge_starrysky_g_pic_01", Path: "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_01.png"},
		{FieldName: "lottery_jiugongge_starrysky_g_pic_02", Path: "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_02.png"},
		{FieldName: "lottery_jiugongge_starrysky_g_pic_03", Path: "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_03.jpg"},
		{FieldName: "lottery_jiugongge_starrysky_g_pic_04", Path: "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_04.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_pic_01", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_01.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_pic_02", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_02.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_pic_03", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_03.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_pic_04", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_04.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_01", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_01.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_02", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_02.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_03", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_03.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_04", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_04.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_05", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_05.png"},
		{FieldName: "lottery_jiugongge_starrysky_c_ani_06", Path: "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_06.png"},

		{FieldName: "lottery_turntable_starrysky_h_pic_01", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_01.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_02", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_02.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_03", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_03.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_04", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_04.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_05", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_05.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_06", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_06.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_07", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_07.png"},
		{FieldName: "lottery_turntable_starrysky_h_pic_08", Path: "lottery/starrysky/lottery_turntable_starrysky_h_pic_08.jpg"},
		{FieldName: "lottery_turntable_starrysky_g_pic_01", Path: "lottery/starrysky/lottery_turntable_starrysky_g_pic_01.png"},
		{FieldName: "lottery_turntable_starrysky_g_pic_02", Path: "lottery/starrysky/lottery_turntable_starrysky_g_pic_02.png"},
		{FieldName: "lottery_turntable_starrysky_g_pic_03", Path: "lottery/starrysky/lottery_turntable_starrysky_g_pic_03.png"},
		{FieldName: "lottery_turntable_starrysky_g_pic_04", Path: "lottery/starrysky/lottery_turntable_starrysky_g_pic_04.jpg"},
		{FieldName: "lottery_turntable_starrysky_g_pic_05", Path: "lottery/starrysky/lottery_turntable_starrysky_g_pic_05.png"},
		{FieldName: "lottery_turntable_starrysky_c_pic_01", Path: "lottery/starrysky/lottery_turntable_starrysky_c_pic_01.png"},
		{FieldName: "lottery_turntable_starrysky_c_pic_02", Path: "lottery/starrysky/lottery_turntable_starrysky_c_pic_02.png"},
		{FieldName: "lottery_turntable_starrysky_c_pic_03", Path: "lottery/starrysky/lottery_turntable_starrysky_c_pic_03.png"},
		{FieldName: "lottery_turntable_starrysky_c_pic_04", Path: "lottery/starrysky/lottery_turntable_starrysky_c_pic_04.png"},
		{FieldName: "lottery_turntable_starrysky_c_pic_05", Path: "lottery/starrysky/lottery_turntable_starrysky_c_pic_05.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_01", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_01.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_02", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_02.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_03", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_03.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_04", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_04.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_05", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_05.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_06", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_06.png"},
		{FieldName: "lottery_turntable_starrysky_c_ani_07", Path: "lottery/starrysky/lottery_turntable_starrysky_c_ani_07.png"},
	}
)

// GetLotteryPanel 遊戲抽獎
func (s *SystemTable) GetLotteryPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/lottery/" + id)
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
		picMap := BuildPictureMap(lotteryPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
		SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "lottery", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(lotteryPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "lottery", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增遊戲抽獎遊戲資料(form-data)
// @Tags Lottery
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param game_type formData string true "遊戲類型" Enums(turntable,jiugongge)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param max_times formData integer true "用戶抽獎次數"
// @param topic formData string true "主題樣式" Enums(01_classic, 02_starrysky)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/form [post]
func POSTLottery(ctx *gin.Context) {
}

// @Summary 新增遊戲抽獎獎品資料(form-data)
// @Tags Lottery Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_name formData string true "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string true "類型" Enums(first, second, third, general, thanks)
// @param prize_picture formData file false "照片"
// @param prize_amount formData integer true "數量"
// @param prize_price formData integer true "價值"
// @param prize_method formData string true "兌獎方式" Enums(site, mail, thanks)
// @param prize_password formData string true "兌獎密碼"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/prize/form [post]
func POSTLotteryPrize(ctx *gin.Context) {
}

// @Summary 編輯遊戲抽獎遊戲資料(form-data)
// @Tags Lottery
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param game_type formData string false "遊戲類型" Enums(turntable,jiugongge)
// @param allow formData string false "允許重複中獎" Enums(open, close)
// @param max_times formData integer false "用戶抽獎次數"
// @param topic formData string false "主題樣式" Enums(01_classic, 02_starrysky)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/form [put]
func PUTLottery(ctx *gin.Context) {
}

// @Summary 編輯遊戲抽獎獎品資料(form-data)
// @Tags Lottery Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_id formData string true "獎品ID"
// @param prize_name formData string false "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string false "類型" Enums(first, second, third, general, thanks)
// @param prize_picture formData file false "照片"
// @param prize_amount formData integer false "數量(同時更新剩餘數量)"
// @param prize_price formData integer false "價值"
// @param prize_method formData string false "兌獎方式" Enums(site, mail, thanks)
// @param prize_password formData string false "兌獎密碼"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/prize/form [put]
func PUTLotteryPrize(ctx *gin.Context) {
}

// @Summary 刪除遊戲抽獎資料(form-data)
// @Tags Lottery
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/form [delete]
func DELETELottery(ctx *gin.Context) {
}

// @Summary 刪除遊戲抽獎獎品資料(form-data)
// @Tags Lottery Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/prize/form [delete]
func DELETELotteryPrize(ctx *gin.Context) {
}

// @Summary 遊戲抽獎遊戲JSON資料
// @Tags Lottery
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery [get]
func LotteryJSON(ctx *gin.Context) {
}

// @Summary 遊戲抽獎獎品JSON資料
// @Tags Lottery Prize
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/lottery/prize [get]
func LotteryPrizeJSON(ctx *gin.Context) {
}

// models.NewGameModel{
// 		UserID:        values.Get("user"),
// 		ActivityID:    values.Get("activity_id"),
// 		Title:         values.Get("title"),
// 		GameType:      values.Get("game_type"),
// 		LimitTime:     "close",
// 		Second:        "0",
// 		MaxPeople:     "0",
// 		People:        "0",
// 		MaxTimes:      values.Get("max_times"),
// 		Allow:         values.Get("allow"),
// 		Percent:       "0",
// 		FirstPrize:    "0",
// 		SecondPrize:   "0",
// 		ThirdPrize:    "0",
// 		GeneralPrize:  "0",
// 		Topic:         values.Get("topic"),
// 		Skin:          values.Get("skin"),
// 		Music:         values.Get("music"),
// 		DisplayName:   "open",
// 		BoxReflection: "",
// 		SamePeople:    "",

// 		// 拔河遊戲
// 		AllowChooseTeam:  "",
// 		LeftTeamName:     "",
// 		LeftTeamPicture:  "",
// 		RightTeamName:    "",
// 		RightTeamPicture: "",
// 		Prize:            "",

// 		// 賓果遊戲
// 		MaxNumber:  "0",
// 		BingoLine:  "0",
// 		RoundPrize: "0",

// 		// 扭蛋機遊戲
// 		GachaMachineReflection: "0",
// 		ReflectiveSwitch:       "open",

// 		// 投票遊戲
// 		VoteScreen:       "",
// 		VoteTimes:        "0",
// 		VoteMethod:       "",
// 		VoteMethodPlayer: "",
// 		VoteRestriction:  "",
// 		AvatarShape:      "",
// 		VoteStartTime:    "",
// 		VoteEndTime:      "",
// 		AutoPlay:         "",
// 		ShowRank:         "",
// 		TitleSwitch:      "",
// 		ArrangementGuest: "",

// 		// 遊戲抽獎自定義
// 		LotteryBgmGaming: update[0],

// 		LotteryJiugonggeClassicHPic01: update[1],
// 		LotteryJiugonggeClassicHPic02: update[2],
// 		LotteryJiugonggeClassicHPic03: update[3],
// 		LotteryJiugonggeClassicHPic04: update[4],
// 		LotteryJiugonggeClassicGPic01: update[5],
// 		LotteryJiugonggeClassicGPic02: update[6],
// 		LotteryJiugonggeClassicCPic01: update[7],
// 		LotteryJiugonggeClassicCPic02: update[8],
// 		LotteryJiugonggeClassicCPic03: update[9],
// 		LotteryJiugonggeClassicCPic04: update[10],
// 		LotteryJiugonggeClassicCAni01: update[11],
// 		LotteryJiugonggeClassicCAni02: update[12],
// 		LotteryJiugonggeClassicCAni03: update[13],

// 		LotteryTurntableClassicHPic01: update[14],
// 		LotteryTurntableClassicHPic02: update[15],
// 		LotteryTurntableClassicHPic03: update[16],
// 		LotteryTurntableClassicHPic04: update[17],
// 		LotteryTurntableClassicGPic01: update[18],
// 		LotteryTurntableClassicGPic02: update[19],
// 		LotteryTurntableClassicCPic01: update[20],
// 		LotteryTurntableClassicCPic02: update[21],
// 		LotteryTurntableClassicCPic03: update[22],
// 		LotteryTurntableClassicCPic04: update[23],
// 		LotteryTurntableClassicCPic05: update[24],
// 		LotteryTurntableClassicCPic06: update[25],
// 		LotteryTurntableClassicCAni01: update[26],
// 		LotteryTurntableClassicCAni02: update[27],
// 		LotteryTurntableClassicCAni03: update[28],

// 		LotteryJiugonggeStarryskyHPic01: update[29],
// 		LotteryJiugonggeStarryskyHPic02: update[30],
// 		LotteryJiugonggeStarryskyHPic03: update[31],
// 		LotteryJiugonggeStarryskyHPic04: update[32],
// 		LotteryJiugonggeStarryskyHPic05: update[33],
// 		LotteryJiugonggeStarryskyHPic06: update[34],
// 		LotteryJiugonggeStarryskyHPic07: update[35],
// 		LotteryJiugonggeStarryskyGPic01: update[36],
// 		LotteryJiugonggeStarryskyGPic02: update[37],
// 		LotteryJiugonggeStarryskyGPic03: update[38],
// 		LotteryJiugonggeStarryskyGPic04: update[39],
// 		LotteryJiugonggeStarryskyCPic01: update[40],
// 		LotteryJiugonggeStarryskyCPic02: update[41],
// 		LotteryJiugonggeStarryskyCPic03: update[42],
// 		LotteryJiugonggeStarryskyCPic04: update[43],
// 		LotteryJiugonggeStarryskyCAni01: update[44],
// 		LotteryJiugonggeStarryskyCAni02: update[45],
// 		LotteryJiugonggeStarryskyCAni03: update[46],
// 		LotteryJiugonggeStarryskyCAni04: update[47],
// 		LotteryJiugonggeStarryskyCAni05: update[48],
// 		LotteryJiugonggeStarryskyCAni06: update[49],

// 		LotteryTurntableStarryskyHPic01: update[50],
// 		LotteryTurntableStarryskyHPic02: update[51],
// 		LotteryTurntableStarryskyHPic03: update[52],
// 		LotteryTurntableStarryskyHPic04: update[53],
// 		LotteryTurntableStarryskyHPic05: update[54],
// 		LotteryTurntableStarryskyHPic06: update[55],
// 		LotteryTurntableStarryskyHPic07: update[56],
// 		LotteryTurntableStarryskyHPic08: update[57],
// 		LotteryTurntableStarryskyGPic01: update[58],
// 		LotteryTurntableStarryskyGPic02: update[59],
// 		LotteryTurntableStarryskyGPic03: update[60],
// 		LotteryTurntableStarryskyGPic04: update[61],
// 		LotteryTurntableStarryskyGPic05: update[62],
// 		LotteryTurntableStarryskyCPic01: update[63],
// 		LotteryTurntableStarryskyCPic02: update[64],
// 		LotteryTurntableStarryskyCPic03: update[65],
// 		LotteryTurntableStarryskyCPic04: update[66],
// 		LotteryTurntableStarryskyCPic05: update[67],
// 		LotteryTurntableStarryskyCAni01: update[68],
// 		LotteryTurntableStarryskyCAni02: update[69],
// 		LotteryTurntableStarryskyCAni03: update[70],
// 		LotteryTurntableStarryskyCAni04: update[71],
// 		LotteryTurntableStarryskyCAni05: update[72],
// 		LotteryTurntableStarryskyCAni06: update[73],
// 		LotteryTurntableStarryskyCAni07: update[74],
// 	}

// pics = []string{
// 遊戲抽獎自定義
// "lottery/%s/bgm/gaming.mp3",

// "lottery/classic/lottery_jiugongge_classic_h_pic_01.png",
// "lottery/classic/lottery_jiugongge_classic_h_pic_02.png",
// "lottery/classic/lottery_jiugongge_classic_h_pic_03.jpg",
// "lottery/classic/lottery_jiugongge_classic_h_pic_04.png",
// "lottery/classic/lottery_jiugongge_classic_g_pic_01.jpg",
// "lottery/classic/lottery_jiugongge_classic_g_pic_02.png",
// "lottery/classic/lottery_jiugongge_classic_c_pic_01.png",
// "lottery/classic/lottery_jiugongge_classic_c_pic_02.png",
// "lottery/classic/lottery_jiugongge_classic_c_pic_03.png",
// "lottery/classic/lottery_jiugongge_classic_c_pic_04.png",
// "lottery/classic/lottery_jiugongge_classic_c_ani_01.png",
// "lottery/classic/lottery_jiugongge_classic_c_ani_02.png",
// "lottery/classic/lottery_jiugongge_classic_c_ani_03.png",

// "lottery/classic/lottery_turntable_classic_h_pic_01.png",
// "lottery/classic/lottery_turntable_classic_h_pic_02.png",
// "lottery/classic/lottery_turntable_classic_h_pic_03.jpg",
// "lottery/classic/lottery_turntable_classic_h_pic_04.png",
// "lottery/classic/lottery_turntable_classic_g_pic_01.jpg",
// "lottery/classic/lottery_turntable_classic_g_pic_02.png",
// "lottery/classic/lottery_turntable_classic_c_pic_01.png",
// "lottery/classic/lottery_turntable_classic_c_pic_02.png",
// "lottery/classic/lottery_turntable_classic_c_pic_03.png",
// "lottery/classic/lottery_turntable_classic_c_pic_04.png",
// "lottery/classic/lottery_turntable_classic_c_pic_05.png",
// "lottery/classic/lottery_turntable_classic_c_pic_06.png",
// "lottery/classic/lottery_turntable_classic_c_ani_01.png",
// "lottery/classic/lottery_turntable_classic_c_ani_02.png",
// "lottery/classic/lottery_turntable_classic_c_ani_03.png",

// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_01.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_02.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_03.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_04.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_05.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_06.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_h_pic_07.jpg",
// "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_01.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_02.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_03.jpg",
// "lottery/starrysky/lottery_jiugongge_starrysky_g_pic_04.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_01.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_02.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_03.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_pic_04.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_01.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_02.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_03.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_04.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_05.png",
// "lottery/starrysky/lottery_jiugongge_starrysky_c_ani_06.png",

// "lottery/starrysky/lottery_turntable_starrysky_h_pic_01.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_02.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_03.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_04.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_05.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_06.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_07.png",
// "lottery/starrysky/lottery_turntable_starrysky_h_pic_08.jpg",
// "lottery/starrysky/lottery_turntable_starrysky_g_pic_01.png",
// "lottery/starrysky/lottery_turntable_starrysky_g_pic_02.png",
// "lottery/starrysky/lottery_turntable_starrysky_g_pic_03.png",
// "lottery/starrysky/lottery_turntable_starrysky_g_pic_04.jpg",
// "lottery/starrysky/lottery_turntable_starrysky_g_pic_05.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_pic_01.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_pic_02.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_pic_03.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_pic_04.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_pic_05.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_01.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_02.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_03.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_04.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_05.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_06.png",
// "lottery/starrysky/lottery_turntable_starrysky_c_ani_07.png",
// }
// fields = []string{
// 遊戲抽獎自定義
// "lottery_bgm_gaming",

// "lottery_jiugongge_classic_h_pic_01",
// "lottery_jiugongge_classic_h_pic_02",
// "lottery_jiugongge_classic_h_pic_03",
// "lottery_jiugongge_classic_h_pic_04",
// "lottery_jiugongge_classic_g_pic_01",
// "lottery_jiugongge_classic_g_pic_02",
// "lottery_jiugongge_classic_c_pic_01",
// "lottery_jiugongge_classic_c_pic_02",
// "lottery_jiugongge_classic_c_pic_03",
// "lottery_jiugongge_classic_c_pic_04",
// "lottery_jiugongge_classic_c_ani_01",
// "lottery_jiugongge_classic_c_ani_02",
// "lottery_jiugongge_classic_c_ani_03",

// "lottery_turntable_classic_h_pic_01",
// "lottery_turntable_classic_h_pic_02",
// "lottery_turntable_classic_h_pic_03",
// "lottery_turntable_classic_h_pic_04",
// "lottery_turntable_classic_g_pic_01",
// "lottery_turntable_classic_g_pic_02",
// "lottery_turntable_classic_c_pic_01",
// "lottery_turntable_classic_c_pic_02",
// "lottery_turntable_classic_c_pic_03",
// "lottery_turntable_classic_c_pic_04",
// "lottery_turntable_classic_c_pic_05",
// "lottery_turntable_classic_c_pic_06",
// "lottery_turntable_classic_c_ani_01",
// "lottery_turntable_classic_c_ani_02",
// "lottery_turntable_classic_c_ani_03",

// "lottery_jiugongge_starrysky_h_pic_01",
// "lottery_jiugongge_starrysky_h_pic_02",
// "lottery_jiugongge_starrysky_h_pic_03",
// "lottery_jiugongge_starrysky_h_pic_04",
// "lottery_jiugongge_starrysky_h_pic_05",
// "lottery_jiugongge_starrysky_h_pic_06",
// "lottery_jiugongge_starrysky_h_pic_07",
// "lottery_jiugongge_starrysky_g_pic_01",
// "lottery_jiugongge_starrysky_g_pic_02",
// "lottery_jiugongge_starrysky_g_pic_03",
// "lottery_jiugongge_starrysky_g_pic_04",
// "lottery_jiugongge_starrysky_c_pic_01",
// "lottery_jiugongge_starrysky_c_pic_02",
// "lottery_jiugongge_starrysky_c_pic_03",
// "lottery_jiugongge_starrysky_c_pic_04",
// "lottery_jiugongge_starrysky_c_ani_01",
// "lottery_jiugongge_starrysky_c_ani_02",
// "lottery_jiugongge_starrysky_c_ani_03",
// "lottery_jiugongge_starrysky_c_ani_04",
// "lottery_jiugongge_starrysky_c_ani_05",
// "lottery_jiugongge_starrysky_c_ani_06",

// "lottery_turntable_starrysky_h_pic_01",
// "lottery_turntable_starrysky_h_pic_02",
// "lottery_turntable_starrysky_h_pic_03",
// "lottery_turntable_starrysky_h_pic_04",
// "lottery_turntable_starrysky_h_pic_05",
// "lottery_turntable_starrysky_h_pic_06",
// "lottery_turntable_starrysky_h_pic_07",
// "lottery_turntable_starrysky_h_pic_08",
// "lottery_turntable_starrysky_g_pic_01",
// "lottery_turntable_starrysky_g_pic_02",
// "lottery_turntable_starrysky_g_pic_03",
// "lottery_turntable_starrysky_g_pic_04",
// "lottery_turntable_starrysky_g_pic_05",
// "lottery_turntable_starrysky_c_pic_01",
// "lottery_turntable_starrysky_c_pic_02",
// "lottery_turntable_starrysky_c_pic_03",
// "lottery_turntable_starrysky_c_pic_04",
// "lottery_turntable_starrysky_c_pic_05",
// "lottery_turntable_starrysky_c_ani_01",
// "lottery_turntable_starrysky_c_ani_02",
// "lottery_turntable_starrysky_c_ani_03",
// "lottery_turntable_starrysky_c_ani_04",
// "lottery_turntable_starrysky_c_ani_05",
// "lottery_turntable_starrysky_c_ani_06",
// "lottery_turntable_starrysky_c_ani_07",
// }
// update = make([]string, 300)
// )

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
// 		// 遊戲抽獎自定義
// 		"lottery/%s/bgm/gaming.mp3",

// 		"lottery/classic/lottery_jiugongge_classic_h_pic_01.png",
// 		"lottery/classic/lottery_jiugongge_classic_h_pic_02.png",
// 		"lottery/classic/lottery_jiugongge_classic_h_pic_03.jpg",
// 		"lottery/classic/lottery_jiugongge_classic_h_pic_04.png",
// 		"lottery/classic/lottery_jiugongge_classic_g_pic_01.jpg",
// 		"lottery/classic/lottery_jiugongge_classic_g_pic_02.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_pic_01.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_pic_02.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_pic_03.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_pic_04.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_ani_01.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_ani_02.png",
// 		"lottery/classic/lottery_jiugongge_classic_c_ani_03.png",

// 		"lottery/classic/lottery_turntable_classic_h_pic_01.png",
// 		"lottery/classic/lottery_turntable_classic_h_pic_02.png",
// 		"lottery/classic/lottery_turntable_classic_h_pic_03.jpg",
// 		"lottery/classic/lottery_turntable_classic_h_pic_04.png",
// 		"lottery/classic/lottery_turntable_classic_g_pic_01.jpg",
// 		"lottery/classic/lottery_turntable_classic_g_pic_02.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_01.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_02.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_03.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_04.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_05.png",
// 		"lottery/classic/lottery_turntable_classic_c_pic_06.png",
// 		"lottery/classic/lottery_turntable_classic_c_ani_01.png",
// 		"lottery/classic/lottery_turntable_classic_c_ani_02.png",
// 		"lottery/classic/lottery_turntable_classic_c_ani_03.png",

// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_01.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_02.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_03.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_04.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_05.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_06.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_h_pic_07.jpg",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_g_pic_01.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_g_pic_02.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_g_pic_03.jpg",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_g_pic_04.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_pic_01.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_pic_02.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_pic_03.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_pic_04.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_01.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_02.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_03.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_04.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_05.png",
// 		"lottery/starrysky/lottery_jiugongge_starrysky_c_ani_06.png",

// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_01.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_02.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_03.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_04.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_05.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_06.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_07.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_h_pic_08.jpg",
// 		"lottery/starrysky/lottery_turntable_starrysky_g_pic_01.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_g_pic_02.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_g_pic_03.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_g_pic_04.jpg",
// 		"lottery/starrysky/lottery_turntable_starrysky_g_pic_05.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_pic_01.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_pic_02.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_pic_03.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_pic_04.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_pic_05.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_01.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_02.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_03.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_04.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_05.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_06.png",
// 		"lottery/starrysky/lottery_turntable_starrysky_c_ani_07.png",
// 	}
// 	fields = []string{
// 		// 遊戲抽獎自定義
// 		"lottery_bgm_gaming",

// 		"lottery_jiugongge_classic_h_pic_01",
// 		"lottery_jiugongge_classic_h_pic_02",
// 		"lottery_jiugongge_classic_h_pic_03",
// 		"lottery_jiugongge_classic_h_pic_04",
// 		"lottery_jiugongge_classic_g_pic_01",
// 		"lottery_jiugongge_classic_g_pic_02",
// 		"lottery_jiugongge_classic_c_pic_01",
// 		"lottery_jiugongge_classic_c_pic_02",
// 		"lottery_jiugongge_classic_c_pic_03",
// 		"lottery_jiugongge_classic_c_pic_04",
// 		"lottery_jiugongge_classic_c_ani_01",
// 		"lottery_jiugongge_classic_c_ani_02",
// 		"lottery_jiugongge_classic_c_ani_03",

// 		"lottery_turntable_classic_h_pic_01",
// 		"lottery_turntable_classic_h_pic_02",
// 		"lottery_turntable_classic_h_pic_03",
// 		"lottery_turntable_classic_h_pic_04",
// 		"lottery_turntable_classic_g_pic_01",
// 		"lottery_turntable_classic_g_pic_02",
// 		"lottery_turntable_classic_c_pic_01",
// 		"lottery_turntable_classic_c_pic_02",
// 		"lottery_turntable_classic_c_pic_03",
// 		"lottery_turntable_classic_c_pic_04",
// 		"lottery_turntable_classic_c_pic_05",
// 		"lottery_turntable_classic_c_pic_06",
// 		"lottery_turntable_classic_c_ani_01",
// 		"lottery_turntable_classic_c_ani_02",
// 		"lottery_turntable_classic_c_ani_03",

// 		"lottery_jiugongge_starrysky_h_pic_01",
// 		"lottery_jiugongge_starrysky_h_pic_02",
// 		"lottery_jiugongge_starrysky_h_pic_03",
// 		"lottery_jiugongge_starrysky_h_pic_04",
// 		"lottery_jiugongge_starrysky_h_pic_05",
// 		"lottery_jiugongge_starrysky_h_pic_06",
// 		"lottery_jiugongge_starrysky_h_pic_07",
// 		"lottery_jiugongge_starrysky_g_pic_01",
// 		"lottery_jiugongge_starrysky_g_pic_02",
// 		"lottery_jiugongge_starrysky_g_pic_03",
// 		"lottery_jiugongge_starrysky_g_pic_04",
// 		"lottery_jiugongge_starrysky_c_pic_01",
// 		"lottery_jiugongge_starrysky_c_pic_02",
// 		"lottery_jiugongge_starrysky_c_pic_03",
// 		"lottery_jiugongge_starrysky_c_pic_04",
// 		"lottery_jiugongge_starrysky_c_ani_01",
// 		"lottery_jiugongge_starrysky_c_ani_02",
// 		"lottery_jiugongge_starrysky_c_ani_03",
// 		"lottery_jiugongge_starrysky_c_ani_04",
// 		"lottery_jiugongge_starrysky_c_ani_05",
// 		"lottery_jiugongge_starrysky_c_ani_06",

// 		"lottery_turntable_starrysky_h_pic_01",
// 		"lottery_turntable_starrysky_h_pic_02",
// 		"lottery_turntable_starrysky_h_pic_03",
// 		"lottery_turntable_starrysky_h_pic_04",
// 		"lottery_turntable_starrysky_h_pic_05",
// 		"lottery_turntable_starrysky_h_pic_06",
// 		"lottery_turntable_starrysky_h_pic_07",
// 		"lottery_turntable_starrysky_h_pic_08",
// 		"lottery_turntable_starrysky_g_pic_01",
// 		"lottery_turntable_starrysky_g_pic_02",
// 		"lottery_turntable_starrysky_g_pic_03",
// 		"lottery_turntable_starrysky_g_pic_04",
// 		"lottery_turntable_starrysky_g_pic_05",
// 		"lottery_turntable_starrysky_c_pic_01",
// 		"lottery_turntable_starrysky_c_pic_02",
// 		"lottery_turntable_starrysky_c_pic_03",
// 		"lottery_turntable_starrysky_c_pic_04",
// 		"lottery_turntable_starrysky_c_pic_05",
// 		"lottery_turntable_starrysky_c_ani_01",
// 		"lottery_turntable_starrysky_c_ani_02",
// 		"lottery_turntable_starrysky_c_ani_03",
// 		"lottery_turntable_starrysky_c_ani_04",
// 		"lottery_turntable_starrysky_c_ani_05",
// 		"lottery_turntable_starrysky_c_ani_06",
// 		"lottery_turntable_starrysky_c_ani_07",
// 	}
// 	update = make([]string, 300)
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
// 	GameType:      values.Get("game_type"),
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        "",
// 	MaxPeople:     "",
// 	People:        "",
// 	MaxTimes:      values.Get("max_times"),
// 	Allow:         values.Get("allow"),
// 	Percent:       "",
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

// 	// 遊戲抽獎自定義
// 	LotteryBgmGaming: update[0],

// 	LotteryJiugonggeClassicHPic01: update[1],
// 	LotteryJiugonggeClassicHPic02: update[2],
// 	LotteryJiugonggeClassicHPic03: update[3],
// 	LotteryJiugonggeClassicHPic04: update[4],
// 	LotteryJiugonggeClassicGPic01: update[5],
// 	LotteryJiugonggeClassicGPic02: update[6],
// 	LotteryJiugonggeClassicCPic01: update[7],
// 	LotteryJiugonggeClassicCPic02: update[8],
// 	LotteryJiugonggeClassicCPic03: update[9],
// 	LotteryJiugonggeClassicCPic04: update[10],
// 	LotteryJiugonggeClassicCAni01: update[11],
// 	LotteryJiugonggeClassicCAni02: update[12],
// 	LotteryJiugonggeClassicCAni03: update[13],

// 	LotteryTurntableClassicHPic01: update[14],
// 	LotteryTurntableClassicHPic02: update[15],
// 	LotteryTurntableClassicHPic03: update[16],
// 	LotteryTurntableClassicHPic04: update[17],
// 	LotteryTurntableClassicGPic01: update[18],
// 	LotteryTurntableClassicGPic02: update[19],
// 	LotteryTurntableClassicCPic01: update[20],
// 	LotteryTurntableClassicCPic02: update[21],
// 	LotteryTurntableClassicCPic03: update[22],
// 	LotteryTurntableClassicCPic04: update[23],
// 	LotteryTurntableClassicCPic05: update[24],
// 	LotteryTurntableClassicCPic06: update[25],
// 	LotteryTurntableClassicCAni01: update[26],
// 	LotteryTurntableClassicCAni02: update[27],
// 	LotteryTurntableClassicCAni03: update[28],

// 	LotteryJiugonggeStarryskyHPic01: update[29],
// 	LotteryJiugonggeStarryskyHPic02: update[30],
// 	LotteryJiugonggeStarryskyHPic03: update[31],
// 	LotteryJiugonggeStarryskyHPic04: update[32],
// 	LotteryJiugonggeStarryskyHPic05: update[33],
// 	LotteryJiugonggeStarryskyHPic06: update[34],
// 	LotteryJiugonggeStarryskyHPic07: update[35],
// 	LotteryJiugonggeStarryskyGPic01: update[36],
// 	LotteryJiugonggeStarryskyGPic02: update[37],
// 	LotteryJiugonggeStarryskyGPic03: update[38],
// 	LotteryJiugonggeStarryskyGPic04: update[39],
// 	LotteryJiugonggeStarryskyCPic01: update[40],
// 	LotteryJiugonggeStarryskyCPic02: update[41],
// 	LotteryJiugonggeStarryskyCPic03: update[42],
// 	LotteryJiugonggeStarryskyCPic04: update[43],
// 	LotteryJiugonggeStarryskyCAni01: update[44],
// 	LotteryJiugonggeStarryskyCAni02: update[45],
// 	LotteryJiugonggeStarryskyCAni03: update[46],
// 	LotteryJiugonggeStarryskyCAni04: update[47],
// 	LotteryJiugonggeStarryskyCAni05: update[48],
// 	LotteryJiugonggeStarryskyCAni06: update[49],

// 	LotteryTurntableStarryskyHPic01: update[50],
// 	LotteryTurntableStarryskyHPic02: update[51],
// 	LotteryTurntableStarryskyHPic03: update[52],
// 	LotteryTurntableStarryskyHPic04: update[53],
// 	LotteryTurntableStarryskyHPic05: update[54],
// 	LotteryTurntableStarryskyHPic06: update[55],
// 	LotteryTurntableStarryskyHPic07: update[56],
// 	LotteryTurntableStarryskyHPic08: update[57],
// 	LotteryTurntableStarryskyGPic01: update[58],
// 	LotteryTurntableStarryskyGPic02: update[59],
// 	LotteryTurntableStarryskyGPic03: update[60],
// 	LotteryTurntableStarryskyGPic04: update[61],
// 	LotteryTurntableStarryskyGPic05: update[62],
// 	LotteryTurntableStarryskyCPic01: update[63],
// 	LotteryTurntableStarryskyCPic02: update[64],
// 	LotteryTurntableStarryskyCPic03: update[65],
// 	LotteryTurntableStarryskyCPic04: update[66],
// 	LotteryTurntableStarryskyCPic05: update[67],
// 	LotteryTurntableStarryskyCAni01: update[68],
// 	LotteryTurntableStarryskyCAni02: update[69],
// 	LotteryTurntableStarryskyCAni03: update[70],
// 	LotteryTurntableStarryskyCAni04: update[71],
// 	LotteryTurntableStarryskyCAni05: update[72],
// 	LotteryTurntableStarryskyCAni06: update[73],
// 	LotteryTurntableStarryskyCAni07: update[74],
// }
