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
	tugofwarPictureFields = []PictureField{
		{FieldName: "left_team_picture", Path: "tugofwar/%s/left_team_picture.png"},
		{FieldName: "right_team_picture", Path: "tugofwar/%s/right_team_picture.png"},

		{FieldName: "tugofwar_bgm_start", Path: "tugofwar/%s/bgm/start.mp3"},
		{FieldName: "tugofwar_bgm_gaming", Path: "tugofwar/%s/bgm/gaming.mp3"},
		{FieldName: "tugofwar_bgm_end", Path: "tugofwar/%s/bgm/end.mp3"},

		// 經典模式圖片
		{FieldName: "tugofwar_classic_h_pic_01", Path: "tugofwar/classic/tugofwar_classic_h_pic_01.png"},
		{FieldName: "tugofwar_classic_h_pic_02", Path: "tugofwar/classic/tugofwar_classic_h_pic_02.png"},
		{FieldName: "tugofwar_classic_h_pic_03", Path: "tugofwar/classic/tugofwar_classic_h_pic_03.png"},
		{FieldName: "tugofwar_classic_h_pic_04", Path: "tugofwar/classic/tugofwar_classic_h_pic_04.png"},
		{FieldName: "tugofwar_classic_h_pic_05", Path: "tugofwar/classic/tugofwar_classic_h_pic_05.png"},
		{FieldName: "tugofwar_classic_h_pic_06", Path: "tugofwar/classic/tugofwar_classic_h_pic_06.png"},
		{FieldName: "tugofwar_classic_h_pic_07", Path: "tugofwar/classic/tugofwar_classic_h_pic_07.png"},
		{FieldName: "tugofwar_classic_h_pic_08", Path: "tugofwar/classic/tugofwar_classic_h_pic_08.jpg"},
		{FieldName: "tugofwar_classic_h_pic_09", Path: "tugofwar/classic/tugofwar_classic_h_pic_09.png"},
		{FieldName: "tugofwar_classic_h_pic_10", Path: "tugofwar/classic/tugofwar_classic_h_pic_10.png"},
		{FieldName: "tugofwar_classic_h_pic_11", Path: "tugofwar/classic/tugofwar_classic_h_pic_11.png"},
		{FieldName: "tugofwar_classic_h_pic_12", Path: "tugofwar/classic/tugofwar_classic_h_pic_12.jpg"},
		{FieldName: "tugofwar_classic_h_pic_13", Path: "tugofwar/classic/tugofwar_classic_h_pic_13.png"},
		{FieldName: "tugofwar_classic_h_pic_14", Path: "tugofwar/classic/tugofwar_classic_h_pic_14.png"},
		{FieldName: "tugofwar_classic_h_pic_15", Path: "tugofwar/classic/tugofwar_classic_h_pic_15.png"},
		{FieldName: "tugofwar_classic_h_pic_16", Path: "tugofwar/classic/tugofwar_classic_h_pic_16.png"},
		{FieldName: "tugofwar_classic_h_pic_17", Path: "tugofwar/classic/tugofwar_classic_h_pic_17.png"},
		{FieldName: "tugofwar_classic_h_pic_18", Path: "tugofwar/classic/tugofwar_classic_h_pic_18.png"},
		{FieldName: "tugofwar_classic_h_pic_19", Path: "tugofwar/classic/tugofwar_classic_h_pic_19.png"},
		{FieldName: "tugofwar_classic_g_pic_01", Path: "tugofwar/classic/tugofwar_classic_g_pic_01.png"},
		{FieldName: "tugofwar_classic_g_pic_02", Path: "tugofwar/classic/tugofwar_classic_g_pic_02.png"},
		{FieldName: "tugofwar_classic_g_pic_03", Path: "tugofwar/classic/tugofwar_classic_g_pic_03.png"},
		{FieldName: "tugofwar_classic_g_pic_04", Path: "tugofwar/classic/tugofwar_classic_g_pic_04.png"},
		{FieldName: "tugofwar_classic_g_pic_05", Path: "tugofwar/classic/tugofwar_classic_g_pic_05.png"},
		{FieldName: "tugofwar_classic_g_pic_06", Path: "tugofwar/classic/tugofwar_classic_g_pic_06.png"},
		{FieldName: "tugofwar_classic_g_pic_07", Path: "tugofwar/classic/tugofwar_classic_g_pic_07.jpg"},
		{FieldName: "tugofwar_classic_g_pic_08", Path: "tugofwar/classic/tugofwar_classic_g_pic_08.png"},
		{FieldName: "tugofwar_classic_g_pic_09", Path: "tugofwar/classic/tugofwar_classic_g_pic_09.png"},
		{FieldName: "tugofwar_classic_h_ani_01", Path: "tugofwar/classic/tugofwar_classic_h_ani_01.png"},
		{FieldName: "tugofwar_classic_h_ani_02", Path: "tugofwar/classic/tugofwar_classic_h_ani_02.png"},
		{FieldName: "tugofwar_classic_h_ani_03", Path: "tugofwar/classic/tugofwar_classic_h_ani_03.png"},
		{FieldName: "tugofwar_classic_c_ani_01", Path: "tugofwar/classic/tugofwar_classic_c_ani_01.png"},

		{FieldName: "tugofwar_school_h_pic_01", Path: "tugofwar/school/tugofwar_school_h_pic_01.png"},
		{FieldName: "tugofwar_school_h_pic_02", Path: "tugofwar/school/tugofwar_school_h_pic_02.png"},
		{FieldName: "tugofwar_school_h_pic_03", Path: "tugofwar/school/tugofwar_school_h_pic_03.png"},
		{FieldName: "tugofwar_school_h_pic_04", Path: "tugofwar/school/tugofwar_school_h_pic_04.png"},
		{FieldName: "tugofwar_school_h_pic_05", Path: "tugofwar/school/tugofwar_school_h_pic_05.png"},
		{FieldName: "tugofwar_school_h_pic_06", Path: "tugofwar/school/tugofwar_school_h_pic_06.png"},
		{FieldName: "tugofwar_school_h_pic_07", Path: "tugofwar/school/tugofwar_school_h_pic_07.png"},
		{FieldName: "tugofwar_school_h_pic_08", Path: "tugofwar/school/tugofwar_school_h_pic_08.png"},
		{FieldName: "tugofwar_school_h_pic_09", Path: "tugofwar/school/tugofwar_school_h_pic_09.png"},
		{FieldName: "tugofwar_school_h_pic_10", Path: "tugofwar/school/tugofwar_school_h_pic_10.png"},
		{FieldName: "tugofwar_school_h_pic_11", Path: "tugofwar/school/tugofwar_school_h_pic_11.png"},
		{FieldName: "tugofwar_school_h_pic_12", Path: "tugofwar/school/tugofwar_school_h_pic_12.png"},
		{FieldName: "tugofwar_school_h_pic_13", Path: "tugofwar/school/tugofwar_school_h_pic_13.png"},
		{FieldName: "tugofwar_school_h_pic_14", Path: "tugofwar/school/tugofwar_school_h_pic_14.png"},
		{FieldName: "tugofwar_school_h_pic_15", Path: "tugofwar/school/tugofwar_school_h_pic_15.png"},
		{FieldName: "tugofwar_school_h_pic_16", Path: "tugofwar/school/tugofwar_school_h_pic_16.png"},
		{FieldName: "tugofwar_school_h_pic_17", Path: "tugofwar/school/tugofwar_school_h_pic_17.png"},
		{FieldName: "tugofwar_school_h_pic_18", Path: "tugofwar/school/tugofwar_school_h_pic_18.png"},
		{FieldName: "tugofwar_school_h_pic_19", Path: "tugofwar/school/tugofwar_school_h_pic_19.png"},
		{FieldName: "tugofwar_school_h_pic_20", Path: "tugofwar/school/tugofwar_school_h_pic_20.png"},
		{FieldName: "tugofwar_school_h_pic_21", Path: "tugofwar/school/tugofwar_school_h_pic_21.png"},
		{FieldName: "tugofwar_school_h_pic_22", Path: "tugofwar/school/tugofwar_school_h_pic_22.png"},
		{FieldName: "tugofwar_school_h_pic_23", Path: "tugofwar/school/tugofwar_school_h_pic_23.png"},
		{FieldName: "tugofwar_school_h_pic_24", Path: "tugofwar/school/tugofwar_school_h_pic_24.png"},
		{FieldName: "tugofwar_school_h_pic_25", Path: "tugofwar/school/tugofwar_school_h_pic_25.png"},
		{FieldName: "tugofwar_school_h_pic_26", Path: "tugofwar/school/tugofwar_school_h_pic_26.png"},
		{FieldName: "tugofwar_school_g_pic_01", Path: "tugofwar/school/tugofwar_school_g_pic_01.png"},
		{FieldName: "tugofwar_school_g_pic_02", Path: "tugofwar/school/tugofwar_school_g_pic_02.jpg"},
		{FieldName: "tugofwar_school_g_pic_03", Path: "tugofwar/school/tugofwar_school_g_pic_03.png"},
		{FieldName: "tugofwar_school_g_pic_04", Path: "tugofwar/school/tugofwar_school_g_pic_04.png"},
		{FieldName: "tugofwar_school_g_pic_05", Path: "tugofwar/school/tugofwar_school_g_pic_05.png"},
		{FieldName: "tugofwar_school_g_pic_06", Path: "tugofwar/school/tugofwar_school_g_pic_06.png"},
		{FieldName: "tugofwar_school_g_pic_07", Path: "tugofwar/school/tugofwar_school_g_pic_07.png"},
		{FieldName: "tugofwar_school_c_pic_01", Path: "tugofwar/school/tugofwar_school_c_pic_01.png"},
		{FieldName: "tugofwar_school_c_pic_02", Path: "tugofwar/school/tugofwar_school_c_pic_02.png"},
		{FieldName: "tugofwar_school_c_pic_03", Path: "tugofwar/school/tugofwar_school_c_pic_03.png"},
		{FieldName: "tugofwar_school_c_pic_04", Path: "tugofwar/school/tugofwar_school_c_pic_04.png"},
		{FieldName: "tugofwar_school_h_ani_01", Path: "tugofwar/school/tugofwar_school_h_ani_01.png"},
		{FieldName: "tugofwar_school_h_ani_02", Path: "tugofwar/school/tugofwar_school_h_ani_02.png"},
		{FieldName: "tugofwar_school_h_ani_03", Path: "tugofwar/school/tugofwar_school_h_ani_03.png"},
		{FieldName: "tugofwar_school_h_ani_04", Path: "tugofwar/school/tugofwar_school_h_ani_04.png"},
		{FieldName: "tugofwar_school_h_ani_05", Path: "tugofwar/school/tugofwar_school_h_ani_05.png"},
		{FieldName: "tugofwar_school_h_ani_06", Path: "tugofwar/school/tugofwar_school_h_ani_06.png"},
		{FieldName: "tugofwar_school_h_ani_07", Path: "tugofwar/school/tugofwar_school_h_ani_07.png"},

		{FieldName: "tugofwar_christmas_h_pic_01", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_01.png"},
		{FieldName: "tugofwar_christmas_h_pic_02", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_02.png"},
		{FieldName: "tugofwar_christmas_h_pic_03", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_03.png"},
		{FieldName: "tugofwar_christmas_h_pic_04", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_04.png"},
		{FieldName: "tugofwar_christmas_h_pic_05", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_05.png"},
		{FieldName: "tugofwar_christmas_h_pic_06", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_06.png"},
		{FieldName: "tugofwar_christmas_h_pic_07", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_07.jpg"},
		{FieldName: "tugofwar_christmas_h_pic_08", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_08.png"},
		{FieldName: "tugofwar_christmas_h_pic_09", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_09.png"},
		{FieldName: "tugofwar_christmas_h_pic_10", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_10.png"},
		{FieldName: "tugofwar_christmas_h_pic_11", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_11.png"},
		{FieldName: "tugofwar_christmas_h_pic_12", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_12.png"},
		{FieldName: "tugofwar_christmas_h_pic_13", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_13.png"},
		{FieldName: "tugofwar_christmas_h_pic_14", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_14.png"},
		{FieldName: "tugofwar_christmas_h_pic_15", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_15.png"},
		{FieldName: "tugofwar_christmas_h_pic_16", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_16.png"},
		{FieldName: "tugofwar_christmas_h_pic_17", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_17.png"},
		{FieldName: "tugofwar_christmas_h_pic_18", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_18.png"},
		{FieldName: "tugofwar_christmas_h_pic_19", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_19.png"},
		{FieldName: "tugofwar_christmas_h_pic_20", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_20.png"},
		{FieldName: "tugofwar_christmas_h_pic_21", Path: "tugofwar/christmas/tugofwar_christmas_h_pic_21.png"},
		{FieldName: "tugofwar_christmas_g_pic_01", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_01.png"},
		{FieldName: "tugofwar_christmas_g_pic_02", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_02.png"},
		{FieldName: "tugofwar_christmas_g_pic_03", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_03.png"},
		{FieldName: "tugofwar_christmas_g_pic_04", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_04.png"},
		{FieldName: "tugofwar_christmas_g_pic_05", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_05.png"},
		{FieldName: "tugofwar_christmas_g_pic_06", Path: "tugofwar/christmas/tugofwar_christmas_g_pic_06.jpg"},
		{FieldName: "tugofwar_christmas_c_pic_01", Path: "tugofwar/christmas/tugofwar_christmas_c_pic_01.png"},
		{FieldName: "tugofwar_christmas_c_pic_02", Path: "tugofwar/christmas/tugofwar_christmas_c_pic_02.png"},
		{FieldName: "tugofwar_christmas_c_pic_03", Path: "tugofwar/christmas/tugofwar_christmas_c_pic_03.png"},
		{FieldName: "tugofwar_christmas_c_pic_04", Path: "tugofwar/christmas/tugofwar_christmas_c_pic_04.png"},
		{FieldName: "tugofwar_christmas_h_ani_01", Path: "tugofwar/christmas/tugofwar_christmas_h_ani_01.png"},
		{FieldName: "tugofwar_christmas_h_ani_02", Path: "tugofwar/christmas/tugofwar_christmas_h_ani_02.png"},
		{FieldName: "tugofwar_christmas_h_ani_03", Path: "tugofwar/christmas/tugofwar_christmas_h_ani_03.png"},
		{FieldName: "tugofwar_christmas_c_ani_01", Path: "tugofwar/christmas/tugofwar_christmas_c_ani_01.png"},
		{FieldName: "tugofwar_christmas_c_ani_02", Path: "tugofwar/christmas/tugofwar_christmas_c_ani_02.png"},
	}
)

// GetTugofwarPanel 拔河遊戲
func (s *SystemTable) GetTugofwarPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/tugofwar/" + id)
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
		picMap := BuildPictureMap(tugofwarPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "tugofwar", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(tugofwarPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "tugofwar", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @ummary 新增拔河遊戲資料(form-data)
// @Tags Tugofwar
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string true "是否限時" Enums(open, close)
// @param second formData integer true "秒數"
// @param max_people formData integer true "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "每隊人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string true "主題樣式" Enums(01_classic)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param allow_choose_team formData string true "允許玩家選擇隊伍" Enums(open, close)
// @param left_team_name formData string true "左邊隊伍名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param left_team_picture formData file false "左邊隊伍照片"
// @param right_team_name formData string true "右邊隊伍名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param right_team_picture formData file false "右邊隊伍照片"
// @param same_people formData string true "隊伍人數是否一致" Enums(open, close)
// @param prize formData string true "獎品發放" Enums(uniform, all)
// @Router /interact/game/tugofwar/form [post]
func POSTTugofwar(ctx *gin.Context) {
}

// @Summary 編輯拔河遊戲資料(form-data)
// @Tags Tugofwar
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string false "是否限時" Enums(open, close)
// @param second formData integer false "秒數"
// @param max_people formData integer false "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "每隊人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_electric)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param allow_choose_team formData string false "允許玩家選擇隊伍" Enums(open, close)
// @param left_team_name formData string false "左邊隊伍名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param left_team_picture formData file false "左邊隊伍照片"
// @param right_team_name formData string false "右邊隊伍名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param right_team_picture formData file false "右邊隊伍照片"
// @param same_people formData string true "隊伍人數是否一致" Enums(open, close)
// @param prize formData string false "獎品發放" Enums(uniform, all)
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/tugofwar/form [put]
func PUTTugofwar(ctx *gin.Context) {
}

// @Summary 編輯拔河遊戲獎品資料(form-data)
// @Tags Tugofwar Prize
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
// @Router /interact/game/tugofwar/prize/form [put]
func PUTTugofwarPrize(ctx *gin.Context) {
}

// @Summary 刪除拔河遊戲資料(form-data)
// @Tags Tugofwar
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/tugofwar/form [delete]
func DELETETugofwar(ctx *gin.Context) {
}

// @Summary 拔河遊戲JSON資料
// @Tags Tugofwar
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/tugofwar [get]
func TugofwarJSON(ctx *gin.Context) {
}

// @Summary 拔河遊戲獎品JSON資料
// @Tags Tugofwar Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/tugofwar/prize [get]
func TugofwarPrizeJSON(ctx *gin.Context) {
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
// 	Percent:       "0",
// 	FirstPrize:    "0",
// 	SecondPrize:   "0",
// 	ThirdPrize:    "0",
// 	GeneralPrize:  "0",
// 	Topic:         values.Get("topic"),
// 	Skin:          values.Get("skin"),
// 	Music:         values.Get("music"),
// 	DisplayName:   "open",
// 	BoxReflection: "",
// 	SamePeople:    values.Get("same_people"),

// 	// 拔河遊戲
// 	AllowChooseTeam:  values.Get("allow_choose_team"),
// 	LeftTeamName:     values.Get("left_team_name"),
// 	LeftTeamPicture:  update[0],
// 	RightTeamName:    values.Get("right_team_name"),
// 	RightTeamPicture: update[1],
// 	Prize:            values.Get("prize"),

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

// 	// 拔河遊戲自定義
// 	// 音樂
// 	TugofwarBgmStart:  update[2], // 遊戲開始
// 	TugofwarBgmGaming: update[3], // 遊戲進行中
// 	TugofwarBgmEnd:    update[4], // 遊戲結束

// 	TugofwarClassicHPic01: update[5],
// 	TugofwarClassicHPic02: update[6],
// 	TugofwarClassicHPic03: update[7],
// 	TugofwarClassicHPic04: update[8],
// 	TugofwarClassicHPic05: update[9],
// 	TugofwarClassicHPic06: update[10],
// 	TugofwarClassicHPic07: update[11],
// 	TugofwarClassicHPic08: update[12],
// 	TugofwarClassicHPic09: update[13],
// 	TugofwarClassicHPic10: update[14],
// 	TugofwarClassicHPic11: update[15],
// 	TugofwarClassicHPic12: update[16],
// 	TugofwarClassicHPic13: update[17],
// 	TugofwarClassicHPic14: update[18],
// 	TugofwarClassicHPic15: update[19],
// 	TugofwarClassicHPic16: update[20],
// 	TugofwarClassicHPic17: update[21],
// 	TugofwarClassicHPic18: update[22],
// 	TugofwarClassicHPic19: update[23],
// 	TugofwarClassicGPic01: update[24],
// 	TugofwarClassicGPic02: update[25],
// 	TugofwarClassicGPic03: update[26],
// 	TugofwarClassicGPic04: update[27],
// 	TugofwarClassicGPic05: update[28],
// 	TugofwarClassicGPic06: update[29],
// 	TugofwarClassicGPic07: update[30],
// 	TugofwarClassicGPic08: update[31],
// 	TugofwarClassicGPic09: update[32],
// 	TugofwarClassicHAni01: update[33],
// 	TugofwarClassicHAni02: update[34],
// 	TugofwarClassicHAni03: update[35],
// 	TugofwarClassicCAni01: update[36],

// 	TugofwarSchoolHPic01: update[37],
// 	TugofwarSchoolHPic02: update[38],
// 	TugofwarSchoolHPic03: update[39],
// 	TugofwarSchoolHPic04: update[40],
// 	TugofwarSchoolHPic05: update[41],
// 	TugofwarSchoolHPic06: update[42],
// 	TugofwarSchoolHPic07: update[43],
// 	TugofwarSchoolHPic08: update[44],
// 	TugofwarSchoolHPic09: update[45],
// 	TugofwarSchoolHPic10: update[46],
// 	TugofwarSchoolHPic11: update[47],
// 	TugofwarSchoolHPic12: update[48],
// 	TugofwarSchoolHPic13: update[49],
// 	TugofwarSchoolHPic14: update[50],
// 	TugofwarSchoolHPic15: update[51],
// 	TugofwarSchoolHPic16: update[52],
// 	TugofwarSchoolHPic17: update[53],
// 	TugofwarSchoolHPic18: update[54],
// 	TugofwarSchoolHPic19: update[55],
// 	TugofwarSchoolHPic20: update[56],
// 	TugofwarSchoolHPic21: update[57],
// 	TugofwarSchoolHPic22: update[58],
// 	TugofwarSchoolHPic23: update[59],
// 	TugofwarSchoolHPic24: update[60],
// 	TugofwarSchoolHPic25: update[61],
// 	TugofwarSchoolHPic26: update[62],
// 	TugofwarSchoolGPic01: update[63],
// 	TugofwarSchoolGPic02: update[64],
// 	TugofwarSchoolGPic03: update[65],
// 	TugofwarSchoolGPic04: update[66],
// 	TugofwarSchoolGPic05: update[67],
// 	TugofwarSchoolGPic06: update[68],
// 	TugofwarSchoolGPic07: update[69],
// 	TugofwarSchoolCPic01: update[70],
// 	TugofwarSchoolCPic02: update[71],
// 	TugofwarSchoolCPic03: update[72],
// 	TugofwarSchoolCPic04: update[73],
// 	TugofwarSchoolHAni01: update[74],
// 	TugofwarSchoolHAni02: update[75],
// 	TugofwarSchoolHAni03: update[76],
// 	TugofwarSchoolHAni04: update[77],
// 	TugofwarSchoolHAni05: update[78],
// 	TugofwarSchoolHAni06: update[79],
// 	TugofwarSchoolHAni07: update[80],

// 	TugofwarChristmasHPic01: update[81],
// 	TugofwarChristmasHPic02: update[82],
// 	TugofwarChristmasHPic03: update[83],
// 	TugofwarChristmasHPic04: update[84],
// 	TugofwarChristmasHPic05: update[85],
// 	TugofwarChristmasHPic06: update[86],
// 	TugofwarChristmasHPic07: update[87],
// 	TugofwarChristmasHPic08: update[88],
// 	TugofwarChristmasHPic09: update[89],
// 	TugofwarChristmasHPic10: update[90],
// 	TugofwarChristmasHPic11: update[91],
// 	TugofwarChristmasHPic12: update[92],
// 	TugofwarChristmasHPic13: update[93],
// 	TugofwarChristmasHPic14: update[94],
// 	TugofwarChristmasHPic15: update[95],
// 	TugofwarChristmasHPic16: update[96],
// 	TugofwarChristmasHPic17: update[97],
// 	TugofwarChristmasHPic18: update[98],
// 	TugofwarChristmasHPic19: update[99],
// 	TugofwarChristmasHPic20: update[100],
// 	TugofwarChristmasHPic21: update[101],
// 	TugofwarChristmasGPic01: update[102],
// 	TugofwarChristmasGPic02: update[103],
// 	TugofwarChristmasGPic03: update[104],
// 	TugofwarChristmasGPic04: update[105],
// 	TugofwarChristmasGPic05: update[106],
// 	TugofwarChristmasGPic06: update[107],
// 	TugofwarChristmasCPic01: update[108],
// 	TugofwarChristmasCPic02: update[109],
// 	TugofwarChristmasCPic03: update[110],
// 	TugofwarChristmasCPic04: update[111],
// 	TugofwarChristmasHAni01: update[112],
// 	TugofwarChristmasHAni02: update[113],
// 	TugofwarChristmasHAni03: update[114],
// 	TugofwarChristmasCAni01: update[115],
// 	TugofwarChristmasCAni02: update[116],
// }

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

// pics = []string{
// "tugofwar/%s/left_team_picture.png",
// "tugofwar/%s/right_team_picture.png",

// 拔河遊戲自定義
// 音樂
// "tugofwar/%s/bgm/start.mp3",
// "tugofwar/%s/bgm/gaming.mp3",
// "tugofwar/%s/bgm/end.mp3",

// "tugofwar/classic/tugofwar_classic_h_pic_01.png",
// "tugofwar/classic/tugofwar_classic_h_pic_02.png",
// "tugofwar/classic/tugofwar_classic_h_pic_03.png",
// "tugofwar/classic/tugofwar_classic_h_pic_04.png",
// "tugofwar/classic/tugofwar_classic_h_pic_05.png",
// "tugofwar/classic/tugofwar_classic_h_pic_06.png",
// "tugofwar/classic/tugofwar_classic_h_pic_07.png",
// "tugofwar/classic/tugofwar_classic_h_pic_08.jpg",
// "tugofwar/classic/tugofwar_classic_h_pic_09.png",
// "tugofwar/classic/tugofwar_classic_h_pic_10.png",
// "tugofwar/classic/tugofwar_classic_h_pic_11.png",
// "tugofwar/classic/tugofwar_classic_h_pic_12.jpg",
// "tugofwar/classic/tugofwar_classic_h_pic_13.png",
// "tugofwar/classic/tugofwar_classic_h_pic_14.png",
// "tugofwar/classic/tugofwar_classic_h_pic_15.png",
// "tugofwar/classic/tugofwar_classic_h_pic_16.png",
// "tugofwar/classic/tugofwar_classic_h_pic_17.png",
// "tugofwar/classic/tugofwar_classic_h_pic_18.png",
// "tugofwar/classic/tugofwar_classic_h_pic_19.png",
// "tugofwar/classic/tugofwar_classic_g_pic_01.png",
// "tugofwar/classic/tugofwar_classic_g_pic_02.png",
// "tugofwar/classic/tugofwar_classic_g_pic_03.png",
// "tugofwar/classic/tugofwar_classic_g_pic_04.png",
// "tugofwar/classic/tugofwar_classic_g_pic_05.png",
// "tugofwar/classic/tugofwar_classic_g_pic_06.png",
// "tugofwar/classic/tugofwar_classic_g_pic_07.jpg",
// "tugofwar/classic/tugofwar_classic_g_pic_08.png",
// "tugofwar/classic/tugofwar_classic_g_pic_09.png",
// "tugofwar/classic/tugofwar_classic_h_ani_01.png",
// "tugofwar/classic/tugofwar_classic_h_ani_02.png",
// "tugofwar/classic/tugofwar_classic_h_ani_03.png",
// "tugofwar/classic/tugofwar_classic_c_ani_01.png",

// "tugofwar/school/tugofwar_school_h_pic_01.png",
// "tugofwar/school/tugofwar_school_h_pic_02.png",
// "tugofwar/school/tugofwar_school_h_pic_03.png",
// "tugofwar/school/tugofwar_school_h_pic_04.png",
// "tugofwar/school/tugofwar_school_h_pic_05.png",
// "tugofwar/school/tugofwar_school_h_pic_06.png",
// "tugofwar/school/tugofwar_school_h_pic_07.png",
// "tugofwar/school/tugofwar_school_h_pic_08.png",
// "tugofwar/school/tugofwar_school_h_pic_09.png",
// "tugofwar/school/tugofwar_school_h_pic_10.png",
// "tugofwar/school/tugofwar_school_h_pic_11.png",
// "tugofwar/school/tugofwar_school_h_pic_12.png",
// "tugofwar/school/tugofwar_school_h_pic_13.png",
// "tugofwar/school/tugofwar_school_h_pic_14.png",
// "tugofwar/school/tugofwar_school_h_pic_15.png",
// "tugofwar/school/tugofwar_school_h_pic_16.png",
// "tugofwar/school/tugofwar_school_h_pic_17.png",
// "tugofwar/school/tugofwar_school_h_pic_18.png",
// "tugofwar/school/tugofwar_school_h_pic_19.png",
// "tugofwar/school/tugofwar_school_h_pic_20.png",
// "tugofwar/school/tugofwar_school_h_pic_21.png",
// "tugofwar/school/tugofwar_school_h_pic_22.png",
// "tugofwar/school/tugofwar_school_h_pic_23.png",
// "tugofwar/school/tugofwar_school_h_pic_24.png",
// "tugofwar/school/tugofwar_school_h_pic_25.png",
// "tugofwar/school/tugofwar_school_h_pic_26.png",
// "tugofwar/school/tugofwar_school_g_pic_01.png",
// "tugofwar/school/tugofwar_school_g_pic_02.jpg",
// "tugofwar/school/tugofwar_school_g_pic_03.png",
// "tugofwar/school/tugofwar_school_g_pic_04.png",
// "tugofwar/school/tugofwar_school_g_pic_05.png",
// "tugofwar/school/tugofwar_school_g_pic_06.png",
// "tugofwar/school/tugofwar_school_g_pic_07.png",
// "tugofwar/school/tugofwar_school_c_pic_01.png",
// "tugofwar/school/tugofwar_school_c_pic_02.png",
// "tugofwar/school/tugofwar_school_c_pic_03.png",
// "tugofwar/school/tugofwar_school_c_pic_04.png",
// "tugofwar/school/tugofwar_school_h_ani_01.png",
// "tugofwar/school/tugofwar_school_h_ani_02.png",
// "tugofwar/school/tugofwar_school_h_ani_03.png",
// "tugofwar/school/tugofwar_school_h_ani_04.png",
// "tugofwar/school/tugofwar_school_h_ani_05.png",
// "tugofwar/school/tugofwar_school_h_ani_06.png",
// "tugofwar/school/tugofwar_school_h_ani_07.png",

// "tugofwar/christmas/tugofwar_christmas_h_pic_01.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_02.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_03.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_04.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_05.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_06.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_07.jpg",
// "tugofwar/christmas/tugofwar_christmas_h_pic_08.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_09.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_10.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_11.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_12.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_13.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_14.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_15.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_16.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_17.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_18.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_19.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_20.png",
// "tugofwar/christmas/tugofwar_christmas_h_pic_21.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_01.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_02.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_03.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_04.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_05.png",
// "tugofwar/christmas/tugofwar_christmas_g_pic_06.jpg",
// "tugofwar/christmas/tugofwar_christmas_c_pic_01.png",
// "tugofwar/christmas/tugofwar_christmas_c_pic_02.png",
// "tugofwar/christmas/tugofwar_christmas_c_pic_03.png",
// "tugofwar/christmas/tugofwar_christmas_c_pic_04.png",
// "tugofwar/christmas/tugofwar_christmas_h_ani_01.png",
// "tugofwar/christmas/tugofwar_christmas_h_ani_02.png",
// "tugofwar/christmas/tugofwar_christmas_h_ani_03.png",
// "tugofwar/christmas/tugofwar_christmas_c_ani_01.png",
// "tugofwar/christmas/tugofwar_christmas_c_ani_02.png",
// }
// fields = []string{
// "left_team_picture",
// "right_team_picture",

// 拔河遊戲自定義
// 音樂
// "tugofwar_bgm_start",  // 遊戲開始
// "tugofwar_bgm_gaming", // 遊戲進行中
// "tugofwar_bgm_end",    // 遊戲結束

// "tugofwar_classic_h_pic_01",
// "tugofwar_classic_h_pic_02",
// "tugofwar_classic_h_pic_03",
// "tugofwar_classic_h_pic_04",
// "tugofwar_classic_h_pic_05",
// "tugofwar_classic_h_pic_06",
// "tugofwar_classic_h_pic_07",
// "tugofwar_classic_h_pic_08",
// "tugofwar_classic_h_pic_09",
// "tugofwar_classic_h_pic_10",
// "tugofwar_classic_h_pic_11",
// "tugofwar_classic_h_pic_12",
// "tugofwar_classic_h_pic_13",
// "tugofwar_classic_h_pic_14",
// "tugofwar_classic_h_pic_15",
// "tugofwar_classic_h_pic_16",
// "tugofwar_classic_h_pic_17",
// "tugofwar_classic_h_pic_18",
// "tugofwar_classic_h_pic_19",
// "tugofwar_classic_g_pic_01",
// "tugofwar_classic_g_pic_02",
// "tugofwar_classic_g_pic_03",
// "tugofwar_classic_g_pic_04",
// "tugofwar_classic_g_pic_05",
// "tugofwar_classic_g_pic_06",
// "tugofwar_classic_g_pic_07",
// "tugofwar_classic_g_pic_08",
// "tugofwar_classic_g_pic_09",
// "tugofwar_classic_h_ani_01",
// "tugofwar_classic_h_ani_02",
// "tugofwar_classic_h_ani_03",
// "tugofwar_classic_c_ani_01",

// "tugofwar_school_h_pic_01",
// "tugofwar_school_h_pic_02",
// "tugofwar_school_h_pic_03",
// "tugofwar_school_h_pic_04",
// "tugofwar_school_h_pic_05",
// "tugofwar_school_h_pic_06",
// "tugofwar_school_h_pic_07",
// "tugofwar_school_h_pic_08",
// "tugofwar_school_h_pic_09",
// "tugofwar_school_h_pic_10",
// "tugofwar_school_h_pic_11",
// "tugofwar_school_h_pic_12",
// "tugofwar_school_h_pic_13",
// "tugofwar_school_h_pic_14",
// "tugofwar_school_h_pic_15",
// "tugofwar_school_h_pic_16",
// "tugofwar_school_h_pic_17",
// "tugofwar_school_h_pic_18",
// "tugofwar_school_h_pic_19",
// "tugofwar_school_h_pic_20",
// "tugofwar_school_h_pic_21",
// "tugofwar_school_h_pic_22",
// "tugofwar_school_h_pic_23",
// "tugofwar_school_h_pic_24",
// "tugofwar_school_h_pic_25",
// "tugofwar_school_h_pic_26",
// "tugofwar_school_g_pic_01",
// "tugofwar_school_g_pic_02",
// "tugofwar_school_g_pic_03",
// "tugofwar_school_g_pic_04",
// "tugofwar_school_g_pic_05",
// "tugofwar_school_g_pic_06",
// "tugofwar_school_g_pic_07",
// "tugofwar_school_c_pic_01",
// "tugofwar_school_c_pic_02",
// "tugofwar_school_c_pic_03",
// "tugofwar_school_c_pic_04",
// "tugofwar_school_h_ani_01",
// "tugofwar_school_h_ani_02",
// "tugofwar_school_h_ani_03",
// "tugofwar_school_h_ani_04",
// "tugofwar_school_h_ani_05",
// "tugofwar_school_h_ani_06",
// "tugofwar_school_h_ani_07",

// "tugofwar_christmas_h_pic_01",
// "tugofwar_christmas_h_pic_02",
// "tugofwar_christmas_h_pic_03",
// "tugofwar_christmas_h_pic_04",
// "tugofwar_christmas_h_pic_05",
// "tugofwar_christmas_h_pic_06",
// "tugofwar_christmas_h_pic_07",
// "tugofwar_christmas_h_pic_08",
// "tugofwar_christmas_h_pic_09",
// "tugofwar_christmas_h_pic_10",
// "tugofwar_christmas_h_pic_11",
// "tugofwar_christmas_h_pic_12",
// "tugofwar_christmas_h_pic_13",
// "tugofwar_christmas_h_pic_14",
// "tugofwar_christmas_h_pic_15",
// "tugofwar_christmas_h_pic_16",
// "tugofwar_christmas_h_pic_17",
// "tugofwar_christmas_h_pic_18",
// "tugofwar_christmas_h_pic_19",
// "tugofwar_christmas_h_pic_20",
// "tugofwar_christmas_h_pic_21",
// "tugofwar_christmas_g_pic_01",
// "tugofwar_christmas_g_pic_02",
// "tugofwar_christmas_g_pic_03",
// "tugofwar_christmas_g_pic_04",
// "tugofwar_christmas_g_pic_05",
// "tugofwar_christmas_g_pic_06",
// "tugofwar_christmas_c_pic_01",
// "tugofwar_christmas_c_pic_02",
// "tugofwar_christmas_c_pic_03",
// "tugofwar_christmas_c_pic_04",
// "tugofwar_christmas_h_ani_01",
// "tugofwar_christmas_h_ani_02",
// "tugofwar_christmas_h_ani_03",
// "tugofwar_christmas_c_ani_01",
// "tugofwar_christmas_c_ani_02",
// }
// update = make([]string, 300)

// var (
// 	pics = []string{
// 		"tugofwar/%s/left_team_picture.png",
// 		"tugofwar/%s/right_team_picture.png",

// 		// 拔河遊戲自定義
// 		// 音樂
// 		"tugofwar/%s/bgm/start.mp3",
// 		"tugofwar/%s/bgm/gaming.mp3",
// 		"tugofwar/%s/bgm/end.mp3",

// 		"tugofwar/classic/tugofwar_classic_h_pic_01.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_02.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_03.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_04.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_05.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_06.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_07.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_08.jpg",
// 		"tugofwar/classic/tugofwar_classic_h_pic_09.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_10.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_11.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_12.jpg",
// 		"tugofwar/classic/tugofwar_classic_h_pic_13.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_14.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_15.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_16.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_17.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_18.png",
// 		"tugofwar/classic/tugofwar_classic_h_pic_19.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_01.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_02.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_03.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_04.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_05.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_06.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_07.jpg",
// 		"tugofwar/classic/tugofwar_classic_g_pic_08.png",
// 		"tugofwar/classic/tugofwar_classic_g_pic_09.png",
// 		"tugofwar/classic/tugofwar_classic_h_ani_01.png",
// 		"tugofwar/classic/tugofwar_classic_h_ani_02.png",
// 		"tugofwar/classic/tugofwar_classic_h_ani_03.png",
// 		"tugofwar/classic/tugofwar_classic_c_ani_01.png",

// 		"tugofwar/school/tugofwar_school_h_pic_01.png",
// 		"tugofwar/school/tugofwar_school_h_pic_02.png",
// 		"tugofwar/school/tugofwar_school_h_pic_03.png",
// 		"tugofwar/school/tugofwar_school_h_pic_04.png",
// 		"tugofwar/school/tugofwar_school_h_pic_05.png",
// 		"tugofwar/school/tugofwar_school_h_pic_06.png",
// 		"tugofwar/school/tugofwar_school_h_pic_07.png",
// 		"tugofwar/school/tugofwar_school_h_pic_08.png",
// 		"tugofwar/school/tugofwar_school_h_pic_09.png",
// 		"tugofwar/school/tugofwar_school_h_pic_10.png",
// 		"tugofwar/school/tugofwar_school_h_pic_11.png",
// 		"tugofwar/school/tugofwar_school_h_pic_12.png",
// 		"tugofwar/school/tugofwar_school_h_pic_13.png",
// 		"tugofwar/school/tugofwar_school_h_pic_14.png",
// 		"tugofwar/school/tugofwar_school_h_pic_15.png",
// 		"tugofwar/school/tugofwar_school_h_pic_16.png",
// 		"tugofwar/school/tugofwar_school_h_pic_17.png",
// 		"tugofwar/school/tugofwar_school_h_pic_18.png",
// 		"tugofwar/school/tugofwar_school_h_pic_19.png",
// 		"tugofwar/school/tugofwar_school_h_pic_20.png",
// 		"tugofwar/school/tugofwar_school_h_pic_21.png",
// 		"tugofwar/school/tugofwar_school_h_pic_22.png",
// 		"tugofwar/school/tugofwar_school_h_pic_23.png",
// 		"tugofwar/school/tugofwar_school_h_pic_24.png",
// 		"tugofwar/school/tugofwar_school_h_pic_25.png",
// 		"tugofwar/school/tugofwar_school_h_pic_26.png",
// 		"tugofwar/school/tugofwar_school_g_pic_01.png",
// 		"tugofwar/school/tugofwar_school_g_pic_02.jpg",
// 		"tugofwar/school/tugofwar_school_g_pic_03.png",
// 		"tugofwar/school/tugofwar_school_g_pic_04.png",
// 		"tugofwar/school/tugofwar_school_g_pic_05.png",
// 		"tugofwar/school/tugofwar_school_g_pic_06.png",
// 		"tugofwar/school/tugofwar_school_g_pic_07.png",
// 		"tugofwar/school/tugofwar_school_c_pic_01.png",
// 		"tugofwar/school/tugofwar_school_c_pic_02.png",
// 		"tugofwar/school/tugofwar_school_c_pic_03.png",
// 		"tugofwar/school/tugofwar_school_c_pic_04.png",
// 		"tugofwar/school/tugofwar_school_h_ani_01.png",
// 		"tugofwar/school/tugofwar_school_h_ani_02.png",
// 		"tugofwar/school/tugofwar_school_h_ani_03.png",
// 		"tugofwar/school/tugofwar_school_h_ani_04.png",
// 		"tugofwar/school/tugofwar_school_h_ani_05.png",
// 		"tugofwar/school/tugofwar_school_h_ani_06.png",
// 		"tugofwar/school/tugofwar_school_h_ani_07.png",

// 		"tugofwar/christmas/tugofwar_christmas_h_pic_01.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_02.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_03.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_04.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_05.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_06.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_07.jpg",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_08.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_09.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_10.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_11.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_12.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_13.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_14.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_15.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_16.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_17.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_18.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_19.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_20.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_pic_21.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_01.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_02.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_03.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_04.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_05.png",
// 		"tugofwar/christmas/tugofwar_christmas_g_pic_06.jpg",
// 		"tugofwar/christmas/tugofwar_christmas_c_pic_01.png",
// 		"tugofwar/christmas/tugofwar_christmas_c_pic_02.png",
// 		"tugofwar/christmas/tugofwar_christmas_c_pic_03.png",
// 		"tugofwar/christmas/tugofwar_christmas_c_pic_04.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_ani_01.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_ani_02.png",
// 		"tugofwar/christmas/tugofwar_christmas_h_ani_03.png",
// 		"tugofwar/christmas/tugofwar_christmas_c_ani_01.png",
// 		"tugofwar/christmas/tugofwar_christmas_c_ani_02.png",
// 	}
// 	fields = []string{
// 		"left_team_picture",
// 		"right_team_picture",

// 		// 拔河遊戲自定義
// 		// 音樂
// 		"tugofwar_bgm_start",  // 遊戲開始
// 		"tugofwar_bgm_gaming", // 遊戲進行中
// 		"tugofwar_bgm_end",    // 遊戲結束

// 		"tugofwar_classic_h_pic_01",
// 		"tugofwar_classic_h_pic_02",
// 		"tugofwar_classic_h_pic_03",
// 		"tugofwar_classic_h_pic_04",
// 		"tugofwar_classic_h_pic_05",
// 		"tugofwar_classic_h_pic_06",
// 		"tugofwar_classic_h_pic_07",
// 		"tugofwar_classic_h_pic_08",
// 		"tugofwar_classic_h_pic_09",
// 		"tugofwar_classic_h_pic_10",
// 		"tugofwar_classic_h_pic_11",
// 		"tugofwar_classic_h_pic_12",
// 		"tugofwar_classic_h_pic_13",
// 		"tugofwar_classic_h_pic_14",
// 		"tugofwar_classic_h_pic_15",
// 		"tugofwar_classic_h_pic_16",
// 		"tugofwar_classic_h_pic_17",
// 		"tugofwar_classic_h_pic_18",
// 		"tugofwar_classic_h_pic_19",
// 		"tugofwar_classic_g_pic_01",
// 		"tugofwar_classic_g_pic_02",
// 		"tugofwar_classic_g_pic_03",
// 		"tugofwar_classic_g_pic_04",
// 		"tugofwar_classic_g_pic_05",
// 		"tugofwar_classic_g_pic_06",
// 		"tugofwar_classic_g_pic_07",
// 		"tugofwar_classic_g_pic_08",
// 		"tugofwar_classic_g_pic_09",
// 		"tugofwar_classic_h_ani_01",
// 		"tugofwar_classic_h_ani_02",
// 		"tugofwar_classic_h_ani_03",
// 		"tugofwar_classic_c_ani_01",

// 		"tugofwar_school_h_pic_01",
// 		"tugofwar_school_h_pic_02",
// 		"tugofwar_school_h_pic_03",
// 		"tugofwar_school_h_pic_04",
// 		"tugofwar_school_h_pic_05",
// 		"tugofwar_school_h_pic_06",
// 		"tugofwar_school_h_pic_07",
// 		"tugofwar_school_h_pic_08",
// 		"tugofwar_school_h_pic_09",
// 		"tugofwar_school_h_pic_10",
// 		"tugofwar_school_h_pic_11",
// 		"tugofwar_school_h_pic_12",
// 		"tugofwar_school_h_pic_13",
// 		"tugofwar_school_h_pic_14",
// 		"tugofwar_school_h_pic_15",
// 		"tugofwar_school_h_pic_16",
// 		"tugofwar_school_h_pic_17",
// 		"tugofwar_school_h_pic_18",
// 		"tugofwar_school_h_pic_19",
// 		"tugofwar_school_h_pic_20",
// 		"tugofwar_school_h_pic_21",
// 		"tugofwar_school_h_pic_22",
// 		"tugofwar_school_h_pic_23",
// 		"tugofwar_school_h_pic_24",
// 		"tugofwar_school_h_pic_25",
// 		"tugofwar_school_h_pic_26",
// 		"tugofwar_school_g_pic_01",
// 		"tugofwar_school_g_pic_02",
// 		"tugofwar_school_g_pic_03",
// 		"tugofwar_school_g_pic_04",
// 		"tugofwar_school_g_pic_05",
// 		"tugofwar_school_g_pic_06",
// 		"tugofwar_school_g_pic_07",
// 		"tugofwar_school_c_pic_01",
// 		"tugofwar_school_c_pic_02",
// 		"tugofwar_school_c_pic_03",
// 		"tugofwar_school_c_pic_04",
// 		"tugofwar_school_h_ani_01",
// 		"tugofwar_school_h_ani_02",
// 		"tugofwar_school_h_ani_03",
// 		"tugofwar_school_h_ani_04",
// 		"tugofwar_school_h_ani_05",
// 		"tugofwar_school_h_ani_06",
// 		"tugofwar_school_h_ani_07",

// 		"tugofwar_christmas_h_pic_01",
// 		"tugofwar_christmas_h_pic_02",
// 		"tugofwar_christmas_h_pic_03",
// 		"tugofwar_christmas_h_pic_04",
// 		"tugofwar_christmas_h_pic_05",
// 		"tugofwar_christmas_h_pic_06",
// 		"tugofwar_christmas_h_pic_07",
// 		"tugofwar_christmas_h_pic_08",
// 		"tugofwar_christmas_h_pic_09",
// 		"tugofwar_christmas_h_pic_10",
// 		"tugofwar_christmas_h_pic_11",
// 		"tugofwar_christmas_h_pic_12",
// 		"tugofwar_christmas_h_pic_13",
// 		"tugofwar_christmas_h_pic_14",
// 		"tugofwar_christmas_h_pic_15",
// 		"tugofwar_christmas_h_pic_16",
// 		"tugofwar_christmas_h_pic_17",
// 		"tugofwar_christmas_h_pic_18",
// 		"tugofwar_christmas_h_pic_19",
// 		"tugofwar_christmas_h_pic_20",
// 		"tugofwar_christmas_h_pic_21",
// 		"tugofwar_christmas_g_pic_01",
// 		"tugofwar_christmas_g_pic_02",
// 		"tugofwar_christmas_g_pic_03",
// 		"tugofwar_christmas_g_pic_04",
// 		"tugofwar_christmas_g_pic_05",
// 		"tugofwar_christmas_g_pic_06",
// 		"tugofwar_christmas_c_pic_01",
// 		"tugofwar_christmas_c_pic_02",
// 		"tugofwar_christmas_c_pic_03",
// 		"tugofwar_christmas_c_pic_04",
// 		"tugofwar_christmas_h_ani_01",
// 		"tugofwar_christmas_h_ani_02",
// 		"tugofwar_christmas_h_ani_03",
// 		"tugofwar_christmas_c_ani_01",
// 		"tugofwar_christmas_c_ani_02",
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
// 	SamePeople:    values.Get("same_people"),

// 	// 拔河遊戲
// 	AllowChooseTeam:  values.Get("allow_choose_team"),
// 	LeftTeamName:     values.Get("left_team_name"),
// 	LeftTeamPicture:  update[0],
// 	RightTeamName:    values.Get("right_team_name"),
// 	RightTeamPicture: update[1],
// 	Prize:            values.Get("prize"),

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

// 	// 拔河遊戲自定義
// 	// 音樂
// 	TugofwarBgmStart:  update[2], // 遊戲開始
// 	TugofwarBgmGaming: update[3], // 遊戲進行中
// 	TugofwarBgmEnd:    update[4], // 遊戲結束

// 	TugofwarClassicHPic01: update[5],
// 	TugofwarClassicHPic02: update[6],
// 	TugofwarClassicHPic03: update[7],
// 	TugofwarClassicHPic04: update[8],
// 	TugofwarClassicHPic05: update[9],
// 	TugofwarClassicHPic06: update[10],
// 	TugofwarClassicHPic07: update[11],
// 	TugofwarClassicHPic08: update[12],
// 	TugofwarClassicHPic09: update[13],
// 	TugofwarClassicHPic10: update[14],
// 	TugofwarClassicHPic11: update[15],
// 	TugofwarClassicHPic12: update[16],
// 	TugofwarClassicHPic13: update[17],
// 	TugofwarClassicHPic14: update[18],
// 	TugofwarClassicHPic15: update[19],
// 	TugofwarClassicHPic16: update[20],
// 	TugofwarClassicHPic17: update[21],
// 	TugofwarClassicHPic18: update[22],
// 	TugofwarClassicHPic19: update[23],
// 	TugofwarClassicGPic01: update[24],
// 	TugofwarClassicGPic02: update[25],
// 	TugofwarClassicGPic03: update[26],
// 	TugofwarClassicGPic04: update[27],
// 	TugofwarClassicGPic05: update[28],
// 	TugofwarClassicGPic06: update[29],
// 	TugofwarClassicGPic07: update[30],
// 	TugofwarClassicGPic08: update[31],
// 	TugofwarClassicGPic09: update[32],
// 	TugofwarClassicHAni01: update[33],
// 	TugofwarClassicHAni02: update[34],
// 	TugofwarClassicHAni03: update[35],
// 	TugofwarClassicCAni01: update[36],

// 	TugofwarSchoolHPic01: update[37],
// 	TugofwarSchoolHPic02: update[38],
// 	TugofwarSchoolHPic03: update[39],
// 	TugofwarSchoolHPic04: update[40],
// 	TugofwarSchoolHPic05: update[41],
// 	TugofwarSchoolHPic06: update[42],
// 	TugofwarSchoolHPic07: update[43],
// 	TugofwarSchoolHPic08: update[44],
// 	TugofwarSchoolHPic09: update[45],
// 	TugofwarSchoolHPic10: update[46],
// 	TugofwarSchoolHPic11: update[47],
// 	TugofwarSchoolHPic12: update[48],
// 	TugofwarSchoolHPic13: update[49],
// 	TugofwarSchoolHPic14: update[50],
// 	TugofwarSchoolHPic15: update[51],
// 	TugofwarSchoolHPic16: update[52],
// 	TugofwarSchoolHPic17: update[53],
// 	TugofwarSchoolHPic18: update[54],
// 	TugofwarSchoolHPic19: update[55],
// 	TugofwarSchoolHPic20: update[56],
// 	TugofwarSchoolHPic21: update[57],
// 	TugofwarSchoolHPic22: update[58],
// 	TugofwarSchoolHPic23: update[59],
// 	TugofwarSchoolHPic24: update[60],
// 	TugofwarSchoolHPic25: update[61],
// 	TugofwarSchoolHPic26: update[62],
// 	TugofwarSchoolGPic01: update[63],
// 	TugofwarSchoolGPic02: update[64],
// 	TugofwarSchoolGPic03: update[65],
// 	TugofwarSchoolGPic04: update[66],
// 	TugofwarSchoolGPic05: update[67],
// 	TugofwarSchoolGPic06: update[68],
// 	TugofwarSchoolGPic07: update[69],
// 	TugofwarSchoolCPic01: update[70],
// 	TugofwarSchoolCPic02: update[71],
// 	TugofwarSchoolCPic03: update[72],
// 	TugofwarSchoolCPic04: update[73],
// 	TugofwarSchoolHAni01: update[74],
// 	TugofwarSchoolHAni02: update[75],
// 	TugofwarSchoolHAni03: update[76],
// 	TugofwarSchoolHAni04: update[77],
// 	TugofwarSchoolHAni05: update[78],
// 	TugofwarSchoolHAni06: update[79],
// 	TugofwarSchoolHAni07: update[80],

// 	TugofwarChristmasHPic01: update[81],
// 	TugofwarChristmasHPic02: update[82],
// 	TugofwarChristmasHPic03: update[83],
// 	TugofwarChristmasHPic04: update[84],
// 	TugofwarChristmasHPic05: update[85],
// 	TugofwarChristmasHPic06: update[86],
// 	TugofwarChristmasHPic07: update[87],
// 	TugofwarChristmasHPic08: update[88],
// 	TugofwarChristmasHPic09: update[89],
// 	TugofwarChristmasHPic10: update[90],
// 	TugofwarChristmasHPic11: update[91],
// 	TugofwarChristmasHPic12: update[92],
// 	TugofwarChristmasHPic13: update[93],
// 	TugofwarChristmasHPic14: update[94],
// 	TugofwarChristmasHPic15: update[95],
// 	TugofwarChristmasHPic16: update[96],
// 	TugofwarChristmasHPic17: update[97],
// 	TugofwarChristmasHPic18: update[98],
// 	TugofwarChristmasHPic19: update[99],
// 	TugofwarChristmasHPic20: update[100],
// 	TugofwarChristmasHPic21: update[101],
// 	TugofwarChristmasGPic01: update[102],
// 	TugofwarChristmasGPic02: update[103],
// 	TugofwarChristmasGPic03: update[104],
// 	TugofwarChristmasGPic04: update[105],
// 	TugofwarChristmasGPic05: update[106],
// 	TugofwarChristmasGPic06: update[107],
// 	TugofwarChristmasCPic01: update[108],
// 	TugofwarChristmasCPic02: update[109],
// 	TugofwarChristmasCPic03: update[110],
// 	TugofwarChristmasCPic04: update[111],
// 	TugofwarChristmasHAni01: update[112],
// 	TugofwarChristmasHAni02: update[113],
// 	TugofwarChristmasHAni03: update[114],
// 	TugofwarChristmasCAni01: update[115],
// 	TugofwarChristmasCAni02: update[116],
// }
