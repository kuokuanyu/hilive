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
	whackmolePictureFields = []PictureField{
		{FieldName: "whackmole_bgm_start", Path: "whackmole/%s/bgm/start.mp3"},
		{FieldName: "whackmole_bgm_gaming", Path: "whackmole/%s/bgm/gaming.mp3"},
		{FieldName: "whackmole_bgm_end", Path: "whackmole/%s/bgm/end.mp3"},

		{FieldName: "whackmole_classic_h_pic_01", Path: "whackmole/classic/whackmole_classic_h_pic_01.png"},
		{FieldName: "whackmole_classic_h_pic_02", Path: "whackmole/classic/whackmole_classic_h_pic_02.jpg"},
		{FieldName: "whackmole_classic_h_pic_03", Path: "whackmole/classic/whackmole_classic_h_pic_03.png"},
		{FieldName: "whackmole_classic_h_pic_04", Path: "whackmole/classic/whackmole_classic_h_pic_04.png"},
		{FieldName: "whackmole_classic_h_pic_05", Path: "whackmole/classic/whackmole_classic_h_pic_05.png"},
		{FieldName: "whackmole_classic_h_pic_06", Path: "whackmole/classic/whackmole_classic_h_pic_06.png"},
		{FieldName: "whackmole_classic_h_pic_07", Path: "whackmole/classic/whackmole_classic_h_pic_07.png"},
		{FieldName: "whackmole_classic_h_pic_08", Path: "whackmole/classic/whackmole_classic_h_pic_08.png"},
		{FieldName: "whackmole_classic_h_pic_09", Path: "whackmole/classic/whackmole_classic_h_pic_09.png"},
		{FieldName: "whackmole_classic_h_pic_10", Path: "whackmole/classic/whackmole_classic_h_pic_10.png"},
		{FieldName: "whackmole_classic_h_pic_11", Path: "whackmole/classic/whackmole_classic_h_pic_11.png"},
		{FieldName: "whackmole_classic_h_pic_12", Path: "whackmole/classic/whackmole_classic_h_pic_12.png"},
		{FieldName: "whackmole_classic_h_pic_13", Path: "whackmole/classic/whackmole_classic_h_pic_13.png"},
		{FieldName: "whackmole_classic_h_pic_14", Path: "whackmole/classic/whackmole_classic_h_pic_14.png"},
		{FieldName: "whackmole_classic_h_pic_15", Path: "whackmole/classic/whackmole_classic_h_pic_15.png"},
		{FieldName: "whackmole_classic_g_pic_01", Path: "whackmole/classic/whackmole_classic_g_pic_01.png"},
		{FieldName: "whackmole_classic_g_pic_02", Path: "whackmole/classic/whackmole_classic_g_pic_02.jpg"},
		{FieldName: "whackmole_classic_g_pic_03", Path: "whackmole/classic/whackmole_classic_g_pic_03.png"},
		{FieldName: "whackmole_classic_g_pic_04", Path: "whackmole/classic/whackmole_classic_g_pic_04.png"},
		{FieldName: "whackmole_classic_g_pic_05", Path: "whackmole/classic/whackmole_classic_g_pic_05.png"},
		{FieldName: "whackmole_classic_c_pic_01", Path: "whackmole/classic/whackmole_classic_c_pic_01.png"},
		{FieldName: "whackmole_classic_c_pic_02", Path: "whackmole/classic/whackmole_classic_c_pic_02.png"},
		{FieldName: "whackmole_classic_c_pic_03", Path: "whackmole/classic/whackmole_classic_c_pic_03.png"},
		{FieldName: "whackmole_classic_c_pic_04", Path: "whackmole/classic/whackmole_classic_c_pic_04.png"},
		{FieldName: "whackmole_classic_c_pic_05", Path: "whackmole/classic/whackmole_classic_c_pic_05.png"},
		{FieldName: "whackmole_classic_c_pic_06", Path: "whackmole/classic/whackmole_classic_c_pic_06.png"},
		{FieldName: "whackmole_classic_c_pic_07", Path: "whackmole/classic/whackmole_classic_c_pic_07.png"},
		{FieldName: "whackmole_classic_c_pic_08", Path: "whackmole/classic/whackmole_classic_c_pic_08.png"},
		{FieldName: "whackmole_classic_c_ani_01", Path: "whackmole/classic/whackmole_classic_c_ani_01.png"},

		{FieldName: "whackmole_halloween_h_pic_01", Path: "whackmole/halloween/whackmole_halloween_h_pic_01.png"},
		{FieldName: "whackmole_halloween_h_pic_02", Path: "whackmole/halloween/whackmole_halloween_h_pic_02.jpg"},
		{FieldName: "whackmole_halloween_h_pic_03", Path: "whackmole/halloween/whackmole_halloween_h_pic_03.png"},
		{FieldName: "whackmole_halloween_h_pic_04", Path: "whackmole/halloween/whackmole_halloween_h_pic_04.png"},
		{FieldName: "whackmole_halloween_h_pic_05", Path: "whackmole/halloween/whackmole_halloween_h_pic_05.png"},
		{FieldName: "whackmole_halloween_h_pic_06", Path: "whackmole/halloween/whackmole_halloween_h_pic_06.png"},
		{FieldName: "whackmole_halloween_h_pic_07", Path: "whackmole/halloween/whackmole_halloween_h_pic_07.png"},
		{FieldName: "whackmole_halloween_h_pic_08", Path: "whackmole/halloween/whackmole_halloween_h_pic_08.png"},
		{FieldName: "whackmole_halloween_h_pic_09", Path: "whackmole/halloween/whackmole_halloween_h_pic_09.png"},
		{FieldName: "whackmole_halloween_h_pic_10", Path: "whackmole/halloween/whackmole_halloween_h_pic_10.png"},
		{FieldName: "whackmole_halloween_h_pic_11", Path: "whackmole/halloween/whackmole_halloween_h_pic_11.png"},
		{FieldName: "whackmole_halloween_h_pic_12", Path: "whackmole/halloween/whackmole_halloween_h_pic_12.png"},
		{FieldName: "whackmole_halloween_h_pic_13", Path: "whackmole/halloween/whackmole_halloween_h_pic_13.png"},
		{FieldName: "whackmole_halloween_h_pic_14", Path: "whackmole/halloween/whackmole_halloween_h_pic_14.png"},
		{FieldName: "whackmole_halloween_h_pic_15", Path: "whackmole/halloween/whackmole_halloween_h_pic_15.png"},
		{FieldName: "whackmole_halloween_g_pic_01", Path: "whackmole/halloween/whackmole_halloween_g_pic_01.png"},
		{FieldName: "whackmole_halloween_g_pic_02", Path: "whackmole/halloween/whackmole_halloween_g_pic_02.jpg"},
		{FieldName: "whackmole_halloween_g_pic_03", Path: "whackmole/halloween/whackmole_halloween_g_pic_03.png"},
		{FieldName: "whackmole_halloween_g_pic_04", Path: "whackmole/halloween/whackmole_halloween_g_pic_04.png"},
		{FieldName: "whackmole_halloween_g_pic_05", Path: "whackmole/halloween/whackmole_halloween_g_pic_05.png"},
		{FieldName: "whackmole_halloween_c_pic_01", Path: "whackmole/halloween/whackmole_halloween_c_pic_01.png"},
		{FieldName: "whackmole_halloween_c_pic_02", Path: "whackmole/halloween/whackmole_halloween_c_pic_02.png"},
		{FieldName: "whackmole_halloween_c_pic_03", Path: "whackmole/halloween/whackmole_halloween_c_pic_03.png"},
		{FieldName: "whackmole_halloween_c_pic_04", Path: "whackmole/halloween/whackmole_halloween_c_pic_04.png"},
		{FieldName: "whackmole_halloween_c_pic_05", Path: "whackmole/halloween/whackmole_halloween_c_pic_05.png"},
		{FieldName: "whackmole_halloween_c_pic_06", Path: "whackmole/halloween/whackmole_halloween_c_pic_06.png"},
		{FieldName: "whackmole_halloween_c_pic_07", Path: "whackmole/halloween/whackmole_halloween_c_pic_07.png"},
		{FieldName: "whackmole_halloween_c_pic_08", Path: "whackmole/halloween/whackmole_halloween_c_pic_08.png"},
		{FieldName: "whackmole_halloween_c_ani_01", Path: "whackmole/halloween/whackmole_halloween_c_ani_01.png"},

		{FieldName: "whackmole_christmas_h_pic_01", Path: "whackmole/christmas/whackmole_christmas_h_pic_01.png"},
		{FieldName: "whackmole_christmas_h_pic_02", Path: "whackmole/christmas/whackmole_christmas_h_pic_02.png"},
		{FieldName: "whackmole_christmas_h_pic_03", Path: "whackmole/christmas/whackmole_christmas_h_pic_03.jpg"},
		{FieldName: "whackmole_christmas_h_pic_04", Path: "whackmole/christmas/whackmole_christmas_h_pic_04.png"},
		{FieldName: "whackmole_christmas_h_pic_05", Path: "whackmole/christmas/whackmole_christmas_h_pic_05.png"},
		{FieldName: "whackmole_christmas_h_pic_06", Path: "whackmole/christmas/whackmole_christmas_h_pic_06.png"},
		{FieldName: "whackmole_christmas_h_pic_07", Path: "whackmole/christmas/whackmole_christmas_h_pic_07.png"},
		{FieldName: "whackmole_christmas_h_pic_08", Path: "whackmole/christmas/whackmole_christmas_h_pic_08.png"},
		{FieldName: "whackmole_christmas_h_pic_09", Path: "whackmole/christmas/whackmole_christmas_h_pic_09.png"},
		{FieldName: "whackmole_christmas_h_pic_10", Path: "whackmole/christmas/whackmole_christmas_h_pic_10.png"},
		{FieldName: "whackmole_christmas_h_pic_11", Path: "whackmole/christmas/whackmole_christmas_h_pic_11.png"},
		{FieldName: "whackmole_christmas_h_pic_12", Path: "whackmole/christmas/whackmole_christmas_h_pic_12.png"},
		{FieldName: "whackmole_christmas_h_pic_13", Path: "whackmole/christmas/whackmole_christmas_h_pic_13.png"},
		{FieldName: "whackmole_christmas_h_pic_14", Path: "whackmole/christmas/whackmole_christmas_h_pic_14.png"},
		{FieldName: "whackmole_christmas_g_pic_01", Path: "whackmole/christmas/whackmole_christmas_g_pic_01.png"},
		{FieldName: "whackmole_christmas_g_pic_02", Path: "whackmole/christmas/whackmole_christmas_g_pic_02.jpg"},
		{FieldName: "whackmole_christmas_g_pic_03", Path: "whackmole/christmas/whackmole_christmas_g_pic_03.png"},
		{FieldName: "whackmole_christmas_g_pic_04", Path: "whackmole/christmas/whackmole_christmas_g_pic_04.png"},
		{FieldName: "whackmole_christmas_g_pic_05", Path: "whackmole/christmas/whackmole_christmas_g_pic_05.png"},
		{FieldName: "whackmole_christmas_g_pic_06", Path: "whackmole/christmas/whackmole_christmas_g_pic_06.png"},
		{FieldName: "whackmole_christmas_g_pic_07", Path: "whackmole/christmas/whackmole_christmas_g_pic_07.png"},
		{FieldName: "whackmole_christmas_g_pic_08", Path: "whackmole/christmas/whackmole_christmas_g_pic_08.png"},
		{FieldName: "whackmole_christmas_c_pic_01", Path: "whackmole/christmas/whackmole_christmas_c_pic_01.png"},
		{FieldName: "whackmole_christmas_c_pic_02", Path: "whackmole/christmas/whackmole_christmas_c_pic_02.png"},
		{FieldName: "whackmole_christmas_c_pic_03", Path: "whackmole/christmas/whackmole_christmas_c_pic_03.png"},
		{FieldName: "whackmole_christmas_c_pic_04", Path: "whackmole/christmas/whackmole_christmas_c_pic_04.png"},
		{FieldName: "whackmole_christmas_c_pic_05", Path: "whackmole/christmas/whackmole_christmas_c_pic_05.png"},
		{FieldName: "whackmole_christmas_c_pic_06", Path: "whackmole/christmas/whackmole_christmas_c_pic_06.png"},
		{FieldName: "whackmole_christmas_c_pic_07", Path: "whackmole/christmas/whackmole_christmas_c_pic_07.png"},
		{FieldName: "whackmole_christmas_c_pic_08", Path: "whackmole/christmas/whackmole_christmas_c_pic_08.png"},
		{FieldName: "whackmole_christmas_c_ani_01", Path: "whackmole/christmas/whackmole_christmas_c_ani_01.png"},
		{FieldName: "whackmole_christmas_c_ani_02", Path: "whackmole/christmas/whackmole_christmas_c_ani_02.png"},
	}
)

// GetWhackMolePanel 敲敲樂
func (s *SystemTable) GetWhackMolePanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/whack_mole/" + id)
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
		picMap := BuildPictureMap(whackmolePictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "whack_mole", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(whackmolePictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "whack_mole", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增敲敲樂遊戲資料(form-data)
// @Tags Whack_Mole
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string true "是否限時" Enums(open, close)
// @param second formData integer true "限時秒數"
// @param max_people formData integer true "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param first_prize formData integer true "頭獎中獎人數上限(上限為50人)" maximum(50)
// @param second_prize formData integer true "二獎中獎人數上限(上限為50人)" maximum(50)
// @param third_prize formData integer true "三獎中獎人數上限(上限為100人)" maximum(100)
// @param general_prize formData integer true "普通獎中獎人數上限(上限為800人)" maximum(800)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string true "主題樣式" Enums(01_classic, 02_halloween)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/form [post]
func POSTWhackMole(ctx *gin.Context) {
}

// @Summary 新增敲敲樂獎品資料(form-data)
// @Tags Whack_Mole Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_name formData string true "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string true "類型" Enums(first, second, third, general)
// @param prize_picture formData file false "照片"
// @param prize_method formData string true "兌獎方式" Enums(site, mail)
// @param prize_password formData string true "兌獎密碼"
// @param prize_amount formData integer true "數量"
// @param prize_price formData integer true "價值"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/prize/form [post]
func POSTWhackMolePrize(ctx *gin.Context) {
}

// @Summary 編輯敲敲樂遊戲資料(form-data)
// @Tags Whack_Mole
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string false "是否限時" Enums(open, close)
// @param second formData integer false "限時秒數"
// @param max_people formData integer false "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param first_prize formData integer false "頭獎中獎人數上限(上限為50人)" maximum(50)
// @param second_prize formData integer false "二獎中獎人數上限(上限為50人)" maximum(50)
// @param third_prize formData integer false "三獎中獎人數上限(上限為100人)" maximum(100)
// @param general_prize formData integer false "普通獎中獎人數上限(上限為800人)" maximum(800)
// @param allow formData string false "允許重複中獎" Enums(open, close)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_halloween)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/form [put]
func PUTWhackMole(ctx *gin.Context) {
}

// @Summary 編輯敲敲樂獎品資料(form-data)
// @Tags Whack_Mole Prize
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
// @param prize_password formData string false "兌獎密碼"
// @param prize_amount formData integer false "數量(同時更新剩餘數量)"
// @param prize_price formData integer false "價值"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/prize/form [put]
func PUTWhackMolePrize(ctx *gin.Context) {
}

// @Summary 刪除敲敲樂遊戲資料(form-data)
// @Tags Whack_Mole
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/form [delete]
func DELETEWhackMole(ctx *gin.Context) {
}

// @Summary 刪除敲敲樂獎品資料(form-data)
// @Tags Whack_Mole Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/prize/form [delete]
func DELETEWhackMolePrize(ctx *gin.Context) {
}

// @Summary 敲敲樂遊戲JSON資料
// @Tags Whack_Mole
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole [get]
func WhackMoleJSON(ctx *gin.Context) {
}

// @Summary 敲敲樂獎品JSON資料
// @Tags Whack_Mole Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/whack_mole/prize [get]
func WhackMolePrizeJSON(ctx *gin.Context) {
}

// s.table(config.ACTIVITY_GAME_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_1).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_2).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_BINGO_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_QA_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_STAFF_GAME_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_STAFF_BLACK_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_STAFF_PK_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_SCORE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_QA_RECORD_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// // 投票
// s.table(config.ACTIVITY_GAME_VOTE_PICTURE_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_VOTE_OPTION_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).WhereIn("game_id", ids).Delete()
// s.table(config.ACTIVITY_GAME_VOTE_RECORD_

// 清除遊戲redis資訊
// s.redisConn.DelCache(config.GAME_REDIS + id)
// s.redisConn.DelCache(config.GAME_TYPE_REDIS + id) // 遊戲種類
// s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + id)
// s.redisConn.DelCache(config.BLACK_STAFFS_GAME_REDIS + id)
// s.redisConn.DelCache(config.SCORES_REDIS + id)
// s.redisConn.DelCache(config.SCORES_2_REDIS + id)
// s.redisConn.DelCache(config.CORRECT_REDIS + id)
// s.redisConn.DelCache(config.QA_REDIS + id)
// s.redisConn.DelCache(config.QA_RECORD_REDIS + id)
// s.redisConn.DelCache(config.WINNING_STAFFS_REDIS + id)
// s.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + id) // 未中獎人員
// s.redisConn.DelCache(config.GAME_TEAM_REDIS + id)
// s.redisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + id)               // 紀錄抽過的號碼，LIST
// s.redisConn.DelCache(config.GAME_BINGO_USER_REDIS + id)                 // 賓果中獎人員，ZSET
// s.redisConn.DelCache(config.GAME_BINGO_USER_NUMBER + id)                // 紀錄玩家的號碼排序，HASH
// s.redisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + id)           // 紀錄玩家是否即將中獎，HASH
// s.redisConn.DelCache(config.GAME_ATTEND_REDIS + id)                     // 遊戲人數資訊，SET
// s.redisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + id)  // 拔河遊戲左隊人數資訊，SET
// s.redisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + id) // 拔河遊戲右隊人數資訊，SET

// 投票
// s.redisConn.DelCache(config.GAME_VOTE_RECORDS_REDIS + id)
// s.redisConn.DelCache(config.VOTE_SPECIAL_OFFICER_REDIS + id)

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// s.redisConn.Publish(config.CHANNEL_GAME_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_QA_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_BLACK_STAFFS_GAME_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+id, "修改資料")
// s.redisConn.Publish(config.CHANNEL_SCORES_REDIS+id, "修改資料")

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

// models.NewGameModel{
// 		UserID:        values.Get("user"),
// 		ActivityID:    values.Get("activity_id"),
// 		Title:         values.Get("title"),
// 		GameType:      "",
// 		LimitTime:     values.Get("limit_time"),
// 		Second:        values.Get("second"),
// 		MaxPeople:     values.Get("max_people"),
// 		People:        values.Get("people"),
// 		MaxTimes:      "0",
// 		Allow:         values.Get("allow"),
// 		Percent:       "0",
// 		FirstPrize:    values.Get("first_prize"),
// 		SecondPrize:   values.Get("second_prize"),
// 		ThirdPrize:    values.Get("third_prize"),
// 		GeneralPrize:  values.Get("general_prize"),
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

// 		// 敲敲樂自定義圖片
// 		// 音樂
// 		WhackmoleBgmStart:  update[0],
// 		WhackmoleBgmGaming: update[1],
// 		WhackmoleBgmEnd:    update[2],

// 		WhackmoleClassicHPic01: update[3],
// 		WhackmoleClassicHPic02: update[4],
// 		WhackmoleClassicHPic03: update[5],
// 		WhackmoleClassicHPic04: update[6],
// 		WhackmoleClassicHPic05: update[7],
// 		WhackmoleClassicHPic06: update[8],
// 		WhackmoleClassicHPic07: update[9],
// 		WhackmoleClassicHPic08: update[10],
// 		WhackmoleClassicHPic09: update[11],
// 		WhackmoleClassicHPic10: update[12],
// 		WhackmoleClassicHPic11: update[13],
// 		WhackmoleClassicHPic12: update[14],
// 		WhackmoleClassicHPic13: update[15],
// 		WhackmoleClassicHPic14: update[16],
// 		WhackmoleClassicHPic15: update[17],
// 		WhackmoleClassicGPic01: update[18],
// 		WhackmoleClassicGPic02: update[19],
// 		WhackmoleClassicGPic03: update[20],
// 		WhackmoleClassicGPic04: update[21],
// 		WhackmoleClassicGPic05: update[22],
// 		WhackmoleClassicCPic01: update[23],
// 		WhackmoleClassicCPic02: update[24],
// 		WhackmoleClassicCPic03: update[25],
// 		WhackmoleClassicCPic04: update[26],
// 		WhackmoleClassicCPic05: update[27],
// 		WhackmoleClassicCPic06: update[28],
// 		WhackmoleClassicCPic07: update[29],
// 		WhackmoleClassicCPic08: update[30],
// 		WhackmoleClassicCAni01: update[31],

// 		WhackmoleHalloweenHPic01: update[32],
// 		WhackmoleHalloweenHPic02: update[33],
// 		WhackmoleHalloweenHPic03: update[34],
// 		WhackmoleHalloweenHPic04: update[35],
// 		WhackmoleHalloweenHPic05: update[36],
// 		WhackmoleHalloweenHPic06: update[37],
// 		WhackmoleHalloweenHPic07: update[38],
// 		WhackmoleHalloweenHPic08: update[39],
// 		WhackmoleHalloweenHPic09: update[40],
// 		WhackmoleHalloweenHPic10: update[41],
// 		WhackmoleHalloweenHPic11: update[42],
// 		WhackmoleHalloweenHPic12: update[43],
// 		WhackmoleHalloweenHPic13: update[44],
// 		WhackmoleHalloweenHPic14: update[45],
// 		WhackmoleHalloweenHPic15: update[46],
// 		WhackmoleHalloweenGPic01: update[47],
// 		WhackmoleHalloweenGPic02: update[48],
// 		WhackmoleHalloweenGPic03: update[49],
// 		WhackmoleHalloweenGPic04: update[50],
// 		WhackmoleHalloweenGPic05: update[51],
// 		WhackmoleHalloweenCPic01: update[52],
// 		WhackmoleHalloweenCPic02: update[53],
// 		WhackmoleHalloweenCPic03: update[54],
// 		WhackmoleHalloweenCPic04: update[55],
// 		WhackmoleHalloweenCPic05: update[56],
// 		WhackmoleHalloweenCPic06: update[57],
// 		WhackmoleHalloweenCPic07: update[58],
// 		WhackmoleHalloweenCPic08: update[59],
// 		WhackmoleHalloweenCAni01: update[60],

// 		WhackmoleChristmasHPic01: update[61],
// 		WhackmoleChristmasHPic02: update[62],
// 		WhackmoleChristmasHPic03: update[63],
// 		WhackmoleChristmasHPic04: update[64],
// 		WhackmoleChristmasHPic05: update[65],
// 		WhackmoleChristmasHPic06: update[66],
// 		WhackmoleChristmasHPic07: update[67],
// 		WhackmoleChristmasHPic08: update[68],
// 		WhackmoleChristmasHPic09: update[69],
// 		WhackmoleChristmasHPic10: update[70],
// 		WhackmoleChristmasHPic11: update[71],
// 		WhackmoleChristmasHPic12: update[72],
// 		WhackmoleChristmasHPic13: update[73],
// 		WhackmoleChristmasHPic14: update[74],
// 		WhackmoleChristmasGPic01: update[75],
// 		WhackmoleChristmasGPic02: update[76],
// 		WhackmoleChristmasGPic03: update[77],
// 		WhackmoleChristmasGPic04: update[78],
// 		WhackmoleChristmasGPic05: update[79],
// 		WhackmoleChristmasGPic06: update[80],
// 		WhackmoleChristmasGPic07: update[81],
// 		WhackmoleChristmasGPic08: update[82],
// 		WhackmoleChristmasCPic01: update[83],
// 		WhackmoleChristmasCPic02: update[84],
// 		WhackmoleChristmasCPic03: update[85],
// 		WhackmoleChristmasCPic04: update[86],
// 		WhackmoleChristmasCPic05: update[87],
// 		WhackmoleChristmasCPic06: update[88],
// 		WhackmoleChristmasCPic07: update[89],
// 		WhackmoleChristmasCPic08: update[90],
// 		WhackmoleChristmasCAni01: update[91],
// 		WhackmoleChristmasCAni02: update[92],
// 	}

// pics = []string{
// 敲敲樂自定義
// 音樂
// "whackmole/%s/bgm/start.mp3",
// "whackmole/%s/bgm/gaming.mp3",
// "whackmole/%s/bgm/end.mp3",

// "whackmole/classic/whackmole_classic_h_pic_01.png",
// "whackmole/classic/whackmole_classic_h_pic_02.jpg",
// "whackmole/classic/whackmole_classic_h_pic_03.png",
// "whackmole/classic/whackmole_classic_h_pic_04.png",
// "whackmole/classic/whackmole_classic_h_pic_05.png",
// "whackmole/classic/whackmole_classic_h_pic_06.png",
// "whackmole/classic/whackmole_classic_h_pic_07.png",
// "whackmole/classic/whackmole_classic_h_pic_08.png",
// "whackmole/classic/whackmole_classic_h_pic_09.png",
// "whackmole/classic/whackmole_classic_h_pic_10.png",
// "whackmole/classic/whackmole_classic_h_pic_11.png",
// "whackmole/classic/whackmole_classic_h_pic_12.png",
// "whackmole/classic/whackmole_classic_h_pic_13.png",
// "whackmole/classic/whackmole_classic_h_pic_14.png",
// "whackmole/classic/whackmole_classic_h_pic_15.png",
// "whackmole/classic/whackmole_classic_g_pic_01.png",
// "whackmole/classic/whackmole_classic_g_pic_02.jpg",
// "whackmole/classic/whackmole_classic_g_pic_03.png",
// "whackmole/classic/whackmole_classic_g_pic_04.png",
// "whackmole/classic/whackmole_classic_g_pic_05.png",
// "whackmole/classic/whackmole_classic_c_pic_01.png",
// "whackmole/classic/whackmole_classic_c_pic_02.png",
// "whackmole/classic/whackmole_classic_c_pic_03.png",
// "whackmole/classic/whackmole_classic_c_pic_04.png",
// "whackmole/classic/whackmole_classic_c_pic_05.png",
// "whackmole/classic/whackmole_classic_c_pic_06.png",
// "whackmole/classic/whackmole_classic_c_pic_07.png",
// "whackmole/classic/whackmole_classic_c_pic_08.png",
// "whackmole/classic/whackmole_classic_c_ani_01.png",

// "whackmole/halloween/whackmole_halloween_h_pic_01.png",
// "whackmole/halloween/whackmole_halloween_h_pic_02.jpg",
// "whackmole/halloween/whackmole_halloween_h_pic_03.png",
// "whackmole/halloween/whackmole_halloween_h_pic_04.png",
// "whackmole/halloween/whackmole_halloween_h_pic_05.png",
// "whackmole/halloween/whackmole_halloween_h_pic_06.png",
// "whackmole/halloween/whackmole_halloween_h_pic_07.png",
// "whackmole/halloween/whackmole_halloween_h_pic_08.png",
// "whackmole/halloween/whackmole_halloween_h_pic_09.png",
// "whackmole/halloween/whackmole_halloween_h_pic_10.png",
// "whackmole/halloween/whackmole_halloween_h_pic_11.png",
// "whackmole/halloween/whackmole_halloween_h_pic_12.png",
// "whackmole/halloween/whackmole_halloween_h_pic_13.png",
// "whackmole/halloween/whackmole_halloween_h_pic_14.png",
// "whackmole/halloween/whackmole_halloween_h_pic_15.png",
// "whackmole/halloween/whackmole_halloween_g_pic_01.png",
// "whackmole/halloween/whackmole_halloween_g_pic_02.jpg",
// "whackmole/halloween/whackmole_halloween_g_pic_03.png",
// "whackmole/halloween/whackmole_halloween_g_pic_04.png",
// "whackmole/halloween/whackmole_halloween_g_pic_05.png",
// "whackmole/halloween/whackmole_halloween_c_pic_01.png",
// "whackmole/halloween/whackmole_halloween_c_pic_02.png",
// "whackmole/halloween/whackmole_halloween_c_pic_03.png",
// "whackmole/halloween/whackmole_halloween_c_pic_04.png",
// "whackmole/halloween/whackmole_halloween_c_pic_05.png",
// "whackmole/halloween/whackmole_halloween_c_pic_06.png",
// "whackmole/halloween/whackmole_halloween_c_pic_07.png",
// "whackmole/halloween/whackmole_halloween_c_pic_08.png",
// "whackmole/halloween/whackmole_halloween_c_ani_01.png",

// "whackmole/christmas/whackmole_christmas_h_pic_01.png",
// "whackmole/christmas/whackmole_christmas_h_pic_02.png",
// "whackmole/christmas/whackmole_christmas_h_pic_03.jpg",
// "whackmole/christmas/whackmole_christmas_h_pic_04.png",
// "whackmole/christmas/whackmole_christmas_h_pic_05.png",
// "whackmole/christmas/whackmole_christmas_h_pic_06.png",
// "whackmole/christmas/whackmole_christmas_h_pic_07.png",
// "whackmole/christmas/whackmole_christmas_h_pic_08.png",
// "whackmole/christmas/whackmole_christmas_h_pic_09.png",
// "whackmole/christmas/whackmole_christmas_h_pic_10.png",
// "whackmole/christmas/whackmole_christmas_h_pic_11.png",
// "whackmole/christmas/whackmole_christmas_h_pic_12.png",
// "whackmole/christmas/whackmole_christmas_h_pic_13.png",
// "whackmole/christmas/whackmole_christmas_h_pic_14.png",
// "whackmole/christmas/whackmole_christmas_g_pic_01.png",
// "whackmole/christmas/whackmole_christmas_g_pic_02.jpg",
// "whackmole/christmas/whackmole_christmas_g_pic_03.png",
// "whackmole/christmas/whackmole_christmas_g_pic_04.png",
// "whackmole/christmas/whackmole_christmas_g_pic_05.png",
// "whackmole/christmas/whackmole_christmas_g_pic_06.png",
// "whackmole/christmas/whackmole_christmas_g_pic_07.png",
// "whackmole/christmas/whackmole_christmas_g_pic_08.png",
// "whackmole/christmas/whackmole_christmas_c_pic_01.png",
// "whackmole/christmas/whackmole_christmas_c_pic_02.png",
// "whackmole/christmas/whackmole_christmas_c_pic_03.png",
// "whackmole/christmas/whackmole_christmas_c_pic_04.png",
// "whackmole/christmas/whackmole_christmas_c_pic_05.png",
// "whackmole/christmas/whackmole_christmas_c_pic_06.png",
// "whackmole/christmas/whackmole_christmas_c_pic_07.png",
// "whackmole/christmas/whackmole_christmas_c_pic_08.png",
// "whackmole/christmas/whackmole_christmas_c_ani_01.png",
// "whackmole/christmas/whackmole_christmas_c_ani_02.png",
// }
// fields = []string{
// 敲敲樂自定義
// 音樂
// "whackmole_bgm_start",
// "whackmole_bgm_gaming",
// "whackmole_bgm_end",

// "whackmole_classic_h_pic_01",
// "whackmole_classic_h_pic_02",
// "whackmole_classic_h_pic_03",
// "whackmole_classic_h_pic_04",
// "whackmole_classic_h_pic_05",
// "whackmole_classic_h_pic_06",
// "whackmole_classic_h_pic_07",
// "whackmole_classic_h_pic_08",
// "whackmole_classic_h_pic_09",
// "whackmole_classic_h_pic_10",
// "whackmole_classic_h_pic_11",
// "whackmole_classic_h_pic_12",
// "whackmole_classic_h_pic_13",
// "whackmole_classic_h_pic_14",
// "whackmole_classic_h_pic_15",
// "whackmole_classic_g_pic_01",
// "whackmole_classic_g_pic_02",
// "whackmole_classic_g_pic_03",
// "whackmole_classic_g_pic_04",
// "whackmole_classic_g_pic_05",
// "whackmole_classic_c_pic_01",
// "whackmole_classic_c_pic_02",
// "whackmole_classic_c_pic_03",
// "whackmole_classic_c_pic_04",
// "whackmole_classic_c_pic_05",
// "whackmole_classic_c_pic_06",
// "whackmole_classic_c_pic_07",
// "whackmole_classic_c_pic_08",
// "whackmole_classic_c_ani_01",

// "whackmole_halloween_h_pic_01",
// "whackmole_halloween_h_pic_02",
// "whackmole_halloween_h_pic_03",
// "whackmole_halloween_h_pic_04",
// "whackmole_halloween_h_pic_05",
// "whackmole_halloween_h_pic_06",
// "whackmole_halloween_h_pic_07",
// "whackmole_halloween_h_pic_08",
// "whackmole_halloween_h_pic_09",
// "whackmole_halloween_h_pic_10",
// "whackmole_halloween_h_pic_11",
// "whackmole_halloween_h_pic_12",
// "whackmole_halloween_h_pic_13",
// "whackmole_halloween_h_pic_14",
// "whackmole_halloween_h_pic_15",
// "whackmole_halloween_g_pic_01",
// "whackmole_halloween_g_pic_02",
// "whackmole_halloween_g_pic_03",
// "whackmole_halloween_g_pic_04",
// "whackmole_halloween_g_pic_05",
// "whackmole_halloween_c_pic_01",
// "whackmole_halloween_c_pic_02",
// "whackmole_halloween_c_pic_03",
// "whackmole_halloween_c_pic_04",
// "whackmole_halloween_c_pic_05",
// "whackmole_halloween_c_pic_06",
// "whackmole_halloween_c_pic_07",
// "whackmole_halloween_c_pic_08",
// "whackmole_halloween_c_ani_01",

// "whackmole_christmas_h_pic_01",
// "whackmole_christmas_h_pic_02",
// "whackmole_christmas_h_pic_03",
// "whackmole_christmas_h_pic_04",
// "whackmole_christmas_h_pic_05",
// "whackmole_christmas_h_pic_06",
// "whackmole_christmas_h_pic_07",
// "whackmole_christmas_h_pic_08",
// "whackmole_christmas_h_pic_09",
// "whackmole_christmas_h_pic_10",
// "whackmole_christmas_h_pic_11",
// "whackmole_christmas_h_pic_12",
// "whackmole_christmas_h_pic_13",
// "whackmole_christmas_h_pic_14",
// "whackmole_christmas_g_pic_01",
// "whackmole_christmas_g_pic_02",
// "whackmole_christmas_g_pic_03",
// "whackmole_christmas_g_pic_04",
// "whackmole_christmas_g_pic_05",
// "whackmole_christmas_g_pic_06",
// "whackmole_christmas_g_pic_07",
// "whackmole_christmas_g_pic_08",
// "whackmole_christmas_c_pic_01",
// "whackmole_christmas_c_pic_02",
// "whackmole_christmas_c_pic_03",
// "whackmole_christmas_c_pic_04",
// "whackmole_christmas_c_pic_05",
// "whackmole_christmas_c_pic_06",
// "whackmole_christmas_c_pic_07",
// "whackmole_christmas_c_pic_08",
// "whackmole_christmas_c_ani_01",
// "whackmole_christmas_c_ani_02",
// }
// update = make([]string, 300)

// var (
// 	pics = []string{
// 		// 敲敲樂自定義
// 		// 音樂
// 		"whackmole/%s/bgm/start.mp3",
// 		"whackmole/%s/bgm/gaming.mp3",
// 		"whackmole/%s/bgm/end.mp3",

// 		"whackmole/classic/whackmole_classic_h_pic_01.png",
// 		"whackmole/classic/whackmole_classic_h_pic_02.jpg",
// 		"whackmole/classic/whackmole_classic_h_pic_03.png",
// 		"whackmole/classic/whackmole_classic_h_pic_04.png",
// 		"whackmole/classic/whackmole_classic_h_pic_05.png",
// 		"whackmole/classic/whackmole_classic_h_pic_06.png",
// 		"whackmole/classic/whackmole_classic_h_pic_07.png",
// 		"whackmole/classic/whackmole_classic_h_pic_08.png",
// 		"whackmole/classic/whackmole_classic_h_pic_09.png",
// 		"whackmole/classic/whackmole_classic_h_pic_10.png",
// 		"whackmole/classic/whackmole_classic_h_pic_11.png",
// 		"whackmole/classic/whackmole_classic_h_pic_12.png",
// 		"whackmole/classic/whackmole_classic_h_pic_13.png",
// 		"whackmole/classic/whackmole_classic_h_pic_14.png",
// 		"whackmole/classic/whackmole_classic_h_pic_15.png",
// 		"whackmole/classic/whackmole_classic_g_pic_01.png",
// 		"whackmole/classic/whackmole_classic_g_pic_02.jpg",
// 		"whackmole/classic/whackmole_classic_g_pic_03.png",
// 		"whackmole/classic/whackmole_classic_g_pic_04.png",
// 		"whackmole/classic/whackmole_classic_g_pic_05.png",
// 		"whackmole/classic/whackmole_classic_c_pic_01.png",
// 		"whackmole/classic/whackmole_classic_c_pic_02.png",
// 		"whackmole/classic/whackmole_classic_c_pic_03.png",
// 		"whackmole/classic/whackmole_classic_c_pic_04.png",
// 		"whackmole/classic/whackmole_classic_c_pic_05.png",
// 		"whackmole/classic/whackmole_classic_c_pic_06.png",
// 		"whackmole/classic/whackmole_classic_c_pic_07.png",
// 		"whackmole/classic/whackmole_classic_c_pic_08.png",
// 		"whackmole/classic/whackmole_classic_c_ani_01.png",

// 		"whackmole/halloween/whackmole_halloween_h_pic_01.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_02.jpg",
// 		"whackmole/halloween/whackmole_halloween_h_pic_03.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_04.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_05.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_06.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_07.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_08.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_09.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_10.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_11.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_12.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_13.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_14.png",
// 		"whackmole/halloween/whackmole_halloween_h_pic_15.png",
// 		"whackmole/halloween/whackmole_halloween_g_pic_01.png",
// 		"whackmole/halloween/whackmole_halloween_g_pic_02.jpg",
// 		"whackmole/halloween/whackmole_halloween_g_pic_03.png",
// 		"whackmole/halloween/whackmole_halloween_g_pic_04.png",
// 		"whackmole/halloween/whackmole_halloween_g_pic_05.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_01.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_02.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_03.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_04.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_05.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_06.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_07.png",
// 		"whackmole/halloween/whackmole_halloween_c_pic_08.png",
// 		"whackmole/halloween/whackmole_halloween_c_ani_01.png",

// 		"whackmole/christmas/whackmole_christmas_h_pic_01.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_02.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_03.jpg",
// 		"whackmole/christmas/whackmole_christmas_h_pic_04.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_05.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_06.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_07.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_08.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_09.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_10.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_11.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_12.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_13.png",
// 		"whackmole/christmas/whackmole_christmas_h_pic_14.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_01.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_02.jpg",
// 		"whackmole/christmas/whackmole_christmas_g_pic_03.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_04.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_05.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_06.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_07.png",
// 		"whackmole/christmas/whackmole_christmas_g_pic_08.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_01.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_02.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_03.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_04.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_05.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_06.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_07.png",
// 		"whackmole/christmas/whackmole_christmas_c_pic_08.png",
// 		"whackmole/christmas/whackmole_christmas_c_ani_01.png",
// 		"whackmole/christmas/whackmole_christmas_c_ani_02.png",
// 	}
// 	fields = []string{
// 		// 敲敲樂自定義
// 		// 音樂
// 		"whackmole_bgm_start",
// 		"whackmole_bgm_gaming",
// 		"whackmole_bgm_end",

// 		"whackmole_classic_h_pic_01",
// 		"whackmole_classic_h_pic_02",
// 		"whackmole_classic_h_pic_03",
// 		"whackmole_classic_h_pic_04",
// 		"whackmole_classic_h_pic_05",
// 		"whackmole_classic_h_pic_06",
// 		"whackmole_classic_h_pic_07",
// 		"whackmole_classic_h_pic_08",
// 		"whackmole_classic_h_pic_09",
// 		"whackmole_classic_h_pic_10",
// 		"whackmole_classic_h_pic_11",
// 		"whackmole_classic_h_pic_12",
// 		"whackmole_classic_h_pic_13",
// 		"whackmole_classic_h_pic_14",
// 		"whackmole_classic_h_pic_15",
// 		"whackmole_classic_g_pic_01",
// 		"whackmole_classic_g_pic_02",
// 		"whackmole_classic_g_pic_03",
// 		"whackmole_classic_g_pic_04",
// 		"whackmole_classic_g_pic_05",
// 		"whackmole_classic_c_pic_01",
// 		"whackmole_classic_c_pic_02",
// 		"whackmole_classic_c_pic_03",
// 		"whackmole_classic_c_pic_04",
// 		"whackmole_classic_c_pic_05",
// 		"whackmole_classic_c_pic_06",
// 		"whackmole_classic_c_pic_07",
// 		"whackmole_classic_c_pic_08",
// 		"whackmole_classic_c_ani_01",

// 		"whackmole_halloween_h_pic_01",
// 		"whackmole_halloween_h_pic_02",
// 		"whackmole_halloween_h_pic_03",
// 		"whackmole_halloween_h_pic_04",
// 		"whackmole_halloween_h_pic_05",
// 		"whackmole_halloween_h_pic_06",
// 		"whackmole_halloween_h_pic_07",
// 		"whackmole_halloween_h_pic_08",
// 		"whackmole_halloween_h_pic_09",
// 		"whackmole_halloween_h_pic_10",
// 		"whackmole_halloween_h_pic_11",
// 		"whackmole_halloween_h_pic_12",
// 		"whackmole_halloween_h_pic_13",
// 		"whackmole_halloween_h_pic_14",
// 		"whackmole_halloween_h_pic_15",
// 		"whackmole_halloween_g_pic_01",
// 		"whackmole_halloween_g_pic_02",
// 		"whackmole_halloween_g_pic_03",
// 		"whackmole_halloween_g_pic_04",
// 		"whackmole_halloween_g_pic_05",
// 		"whackmole_halloween_c_pic_01",
// 		"whackmole_halloween_c_pic_02",
// 		"whackmole_halloween_c_pic_03",
// 		"whackmole_halloween_c_pic_04",
// 		"whackmole_halloween_c_pic_05",
// 		"whackmole_halloween_c_pic_06",
// 		"whackmole_halloween_c_pic_07",
// 		"whackmole_halloween_c_pic_08",
// 		"whackmole_halloween_c_ani_01",

// 		"whackmole_christmas_h_pic_01",
// 		"whackmole_christmas_h_pic_02",
// 		"whackmole_christmas_h_pic_03",
// 		"whackmole_christmas_h_pic_04",
// 		"whackmole_christmas_h_pic_05",
// 		"whackmole_christmas_h_pic_06",
// 		"whackmole_christmas_h_pic_07",
// 		"whackmole_christmas_h_pic_08",
// 		"whackmole_christmas_h_pic_09",
// 		"whackmole_christmas_h_pic_10",
// 		"whackmole_christmas_h_pic_11",
// 		"whackmole_christmas_h_pic_12",
// 		"whackmole_christmas_h_pic_13",
// 		"whackmole_christmas_h_pic_14",
// 		"whackmole_christmas_g_pic_01",
// 		"whackmole_christmas_g_pic_02",
// 		"whackmole_christmas_g_pic_03",
// 		"whackmole_christmas_g_pic_04",
// 		"whackmole_christmas_g_pic_05",
// 		"whackmole_christmas_g_pic_06",
// 		"whackmole_christmas_g_pic_07",
// 		"whackmole_christmas_g_pic_08",
// 		"whackmole_christmas_c_pic_01",
// 		"whackmole_christmas_c_pic_02",
// 		"whackmole_christmas_c_pic_03",
// 		"whackmole_christmas_c_pic_04",
// 		"whackmole_christmas_c_pic_05",
// 		"whackmole_christmas_c_pic_06",
// 		"whackmole_christmas_c_pic_07",
// 		"whackmole_christmas_c_pic_08",
// 		"whackmole_christmas_c_ani_01",
// 		"whackmole_christmas_c_ani_02",
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
// 	GameType:      "",
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        values.Get("second"),
// 	MaxPeople:     values.Get("max_people"),
// 	People:        values.Get("people"),
// 	MaxTimes:      "",
// 	Allow:         values.Get("allow"),
// 	Percent:       "",
// 	FirstPrize:    values.Get("first_prize"),
// 	SecondPrize:   values.Get("second_prize"),
// 	ThirdPrize:    values.Get("third_prize"),
// 	GeneralPrize:  values.Get("general_prize"),
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

// 	// 敲敲樂自定義圖片
// 	// 音樂
// 	WhackmoleBgmStart:  update[0],
// 	WhackmoleBgmGaming: update[1],
// 	WhackmoleBgmEnd:    update[2],

// 	WhackmoleClassicHPic01: update[3],
// 	WhackmoleClassicHPic02: update[4],
// 	WhackmoleClassicHPic03: update[5],
// 	WhackmoleClassicHPic04: update[6],
// 	WhackmoleClassicHPic05: update[7],
// 	WhackmoleClassicHPic06: update[8],
// 	WhackmoleClassicHPic07: update[9],
// 	WhackmoleClassicHPic08: update[10],
// 	WhackmoleClassicHPic09: update[11],
// 	WhackmoleClassicHPic10: update[12],
// 	WhackmoleClassicHPic11: update[13],
// 	WhackmoleClassicHPic12: update[14],
// 	WhackmoleClassicHPic13: update[15],
// 	WhackmoleClassicHPic14: update[16],
// 	WhackmoleClassicHPic15: update[17],
// 	WhackmoleClassicGPic01: update[18],
// 	WhackmoleClassicGPic02: update[19],
// 	WhackmoleClassicGPic03: update[20],
// 	WhackmoleClassicGPic04: update[21],
// 	WhackmoleClassicGPic05: update[22],
// 	WhackmoleClassicCPic01: update[23],
// 	WhackmoleClassicCPic02: update[24],
// 	WhackmoleClassicCPic03: update[25],
// 	WhackmoleClassicCPic04: update[26],
// 	WhackmoleClassicCPic05: update[27],
// 	WhackmoleClassicCPic06: update[28],
// 	WhackmoleClassicCPic07: update[29],
// 	WhackmoleClassicCPic08: update[30],
// 	WhackmoleClassicCAni01: update[31],

// 	WhackmoleHalloweenHPic01: update[32],
// 	WhackmoleHalloweenHPic02: update[33],
// 	WhackmoleHalloweenHPic03: update[34],
// 	WhackmoleHalloweenHPic04: update[35],
// 	WhackmoleHalloweenHPic05: update[36],
// 	WhackmoleHalloweenHPic06: update[37],
// 	WhackmoleHalloweenHPic07: update[38],
// 	WhackmoleHalloweenHPic08: update[39],
// 	WhackmoleHalloweenHPic09: update[40],
// 	WhackmoleHalloweenHPic10: update[41],
// 	WhackmoleHalloweenHPic11: update[42],
// 	WhackmoleHalloweenHPic12: update[43],
// 	WhackmoleHalloweenHPic13: update[44],
// 	WhackmoleHalloweenHPic14: update[45],
// 	WhackmoleHalloweenHPic15: update[46],
// 	WhackmoleHalloweenGPic01: update[47],
// 	WhackmoleHalloweenGPic02: update[48],
// 	WhackmoleHalloweenGPic03: update[49],
// 	WhackmoleHalloweenGPic04: update[50],
// 	WhackmoleHalloweenGPic05: update[51],
// 	WhackmoleHalloweenCPic01: update[52],
// 	WhackmoleHalloweenCPic02: update[53],
// 	WhackmoleHalloweenCPic03: update[54],
// 	WhackmoleHalloweenCPic04: update[55],
// 	WhackmoleHalloweenCPic05: update[56],
// 	WhackmoleHalloweenCPic06: update[57],
// 	WhackmoleHalloweenCPic07: update[58],
// 	WhackmoleHalloweenCPic08: update[59],
// 	WhackmoleHalloweenCAni01: update[60],

// 	WhackmoleChristmasHPic01: update[61],
// 	WhackmoleChristmasHPic02: update[62],
// 	WhackmoleChristmasHPic03: update[63],
// 	WhackmoleChristmasHPic04: update[64],
// 	WhackmoleChristmasHPic05: update[65],
// 	WhackmoleChristmasHPic06: update[66],
// 	WhackmoleChristmasHPic07: update[67],
// 	WhackmoleChristmasHPic08: update[68],
// 	WhackmoleChristmasHPic09: update[69],
// 	WhackmoleChristmasHPic10: update[70],
// 	WhackmoleChristmasHPic11: update[71],
// 	WhackmoleChristmasHPic12: update[72],
// 	WhackmoleChristmasHPic13: update[73],
// 	WhackmoleChristmasHPic14: update[74],
// 	WhackmoleChristmasGPic01: update[75],
// 	WhackmoleChristmasGPic02: update[76],
// 	WhackmoleChristmasGPic03: update[77],
// 	WhackmoleChristmasGPic04: update[78],
// 	WhackmoleChristmasGPic05: update[79],
// 	WhackmoleChristmasGPic06: update[80],
// 	WhackmoleChristmasGPic07: update[81],
// 	WhackmoleChristmasGPic08: update[82],
// 	WhackmoleChristmasCPic01: update[83],
// 	WhackmoleChristmasCPic02: update[84],
// 	WhackmoleChristmasCPic03: update[85],
// 	WhackmoleChristmasCPic04: update[86],
// 	WhackmoleChristmasCPic05: update[87],
// 	WhackmoleChristmasCPic06: update[88],
// 	WhackmoleChristmasCPic07: update[89],
// 	WhackmoleChristmasCPic08: update[90],
// 	WhackmoleChristmasCAni01: update[91],
// 	WhackmoleChristmasCAni02: update[92],
// }
