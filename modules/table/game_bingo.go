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
	bingoPictureFields = []PictureField{
		{FieldName: "bingo_bgm_start", Path: "bingo/%s/bgm/start.mp3"},
		{FieldName: "bingo_bgm_gaming", Path: "bingo/%s/bgm/gaming.mp3"},
		{FieldName: "bingo_bgm_end", Path: "bingo/%s/bgm/end.mp3"},

		{FieldName: "bingo_classic_h_pic_01", Path: "bingo/classic/bingo_classic_h_pic_01.png"},
		{FieldName: "bingo_classic_h_pic_02", Path: "bingo/classic/bingo_classic_h_pic_02.png"},
		{FieldName: "bingo_classic_h_pic_03", Path: "bingo/classic/bingo_classic_h_pic_03.png"},
		{FieldName: "bingo_classic_h_pic_04", Path: "bingo/classic/bingo_classic_h_pic_04.png"},
		{FieldName: "bingo_classic_h_pic_05", Path: "bingo/classic/bingo_classic_h_pic_05.jpg"},
		{FieldName: "bingo_classic_h_pic_06", Path: "bingo/classic/bingo_classic_h_pic_06.png"},
		{FieldName: "bingo_classic_h_pic_07", Path: "bingo/classic/bingo_classic_h_pic_07.png"},
		{FieldName: "bingo_classic_h_pic_08", Path: "bingo/classic/bingo_classic_h_pic_08.png"},
		{FieldName: "bingo_classic_h_pic_09", Path: "bingo/classic/bingo_classic_h_pic_09.png"},
		{FieldName: "bingo_classic_h_pic_10", Path: "bingo/classic/bingo_classic_h_pic_10.png"},
		{FieldName: "bingo_classic_h_pic_11", Path: "bingo/classic/bingo_classic_h_pic_11.png"},
		{FieldName: "bingo_classic_h_pic_12", Path: "bingo/classic/bingo_classic_h_pic_12.png"},
		{FieldName: "bingo_classic_h_pic_13", Path: "bingo/classic/bingo_classic_h_pic_13.png"},
		{FieldName: "bingo_classic_h_pic_14", Path: "bingo/classic/bingo_classic_h_pic_14.png"},
		{FieldName: "bingo_classic_h_pic_15", Path: "bingo/classic/bingo_classic_h_pic_15.png"},
		{FieldName: "bingo_classic_h_pic_16", Path: "bingo/classic/bingo_classic_h_pic_16.png"},
		{FieldName: "bingo_classic_g_pic_01", Path: "bingo/classic/bingo_classic_g_pic_01.png"},
		{FieldName: "bingo_classic_g_pic_02", Path: "bingo/classic/bingo_classic_g_pic_02.png"},
		{FieldName: "bingo_classic_g_pic_03", Path: "bingo/classic/bingo_classic_g_pic_03.png"},
		{FieldName: "bingo_classic_g_pic_04", Path: "bingo/classic/bingo_classic_g_pic_04.jpg"},
		{FieldName: "bingo_classic_g_pic_05", Path: "bingo/classic/bingo_classic_g_pic_05.png"},
		{FieldName: "bingo_classic_g_pic_06", Path: "bingo/classic/bingo_classic_g_pic_06.png"},
		{FieldName: "bingo_classic_g_pic_07", Path: "bingo/classic/bingo_classic_g_pic_07.png"},
		{FieldName: "bingo_classic_g_pic_08", Path: "bingo/classic/bingo_classic_g_pic_08.png"},
		{FieldName: "bingo_classic_c_pic_01", Path: "bingo/classic/bingo_classic_c_pic_01.png"},
		{FieldName: "bingo_classic_c_pic_02", Path: "bingo/classic/bingo_classic_c_pic_02.png"},
		{FieldName: "bingo_classic_c_pic_03", Path: "bingo/classic/bingo_classic_c_pic_03.png"},
		{FieldName: "bingo_classic_c_pic_04", Path: "bingo/classic/bingo_classic_c_pic_04.png"},
		{FieldName: "bingo_classic_h_ani_01", Path: "bingo/classic/bingo_classic_h_ani_01.png"},
		{FieldName: "bingo_classic_g_ani_01", Path: "bingo/classic/bingo_classic_g_ani_01.png"},
		{FieldName: "bingo_classic_c_ani_01", Path: "bingo/classic/bingo_classic_c_ani_01.png"},
		{FieldName: "bingo_classic_c_ani_02", Path: "bingo/classic/bingo_classic_c_ani_02.png"},

		{FieldName: "bingo_newyear_dragon_h_pic_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_01.jpg"},
		{FieldName: "bingo_newyear_dragon_h_pic_02", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_02.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_03", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_03.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_04", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_04.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_05", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_05.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_06", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_06.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_07", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_07.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_08", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_08.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_09", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_09.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_10", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_10.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_11", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_11.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_12", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_12.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_13", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_13.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_14", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_14.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_16", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_16.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_17", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_17.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_18", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_18.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_19", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_19.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_20", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_20.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_21", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_21.png"},
		{FieldName: "bingo_newyear_dragon_h_pic_22", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_22.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_01.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_02", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_02.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_03", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_03.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_04", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_04.jpg"},
		{FieldName: "bingo_newyear_dragon_g_pic_05", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_05.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_06", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_06.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_07", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_07.png"},
		{FieldName: "bingo_newyear_dragon_g_pic_08", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_pic_08.png"},
		{FieldName: "bingo_newyear_dragon_c_pic_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_pic_01.png"},
		{FieldName: "bingo_newyear_dragon_c_pic_02", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_pic_02.png"},
		{FieldName: "bingo_newyear_dragon_c_pic_03", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_pic_03.png"},
		{FieldName: "bingo_newyear_dragon_h_ani_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_h_ani_01.png"},
		{FieldName: "bingo_newyear_dragon_g_ani_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_g_ani_01.png"},
		{FieldName: "bingo_newyear_dragon_c_ani_01", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_ani_01.png"},
		{FieldName: "bingo_newyear_dragon_c_ani_02", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_ani_02.png"},
		{FieldName: "bingo_newyear_dragon_c_ani_03", Path: "bingo/newyear_dragon/bingo_newyear_dragon_c_ani_03.png"},

		{FieldName: "bingo_cherry_h_pic_01", Path: "bingo/cherry/bingo_cherry_h_pic_01.png"},
		{FieldName: "bingo_cherry_h_pic_02", Path: "bingo/cherry/bingo_cherry_h_pic_02.png"},
		{FieldName: "bingo_cherry_h_pic_03", Path: "bingo/cherry/bingo_cherry_h_pic_03.png"},
		{FieldName: "bingo_cherry_h_pic_04", Path: "bingo/cherry/bingo_cherry_h_pic_04.png"},
		{FieldName: "bingo_cherry_h_pic_05", Path: "bingo/cherry/bingo_cherry_h_pic_05.jpg"},
		{FieldName: "bingo_cherry_h_pic_06", Path: "bingo/cherry/bingo_cherry_h_pic_06.png"},
		{FieldName: "bingo_cherry_h_pic_07", Path: "bingo/cherry/bingo_cherry_h_pic_07.png"},
		{FieldName: "bingo_cherry_h_pic_08", Path: "bingo/cherry/bingo_cherry_h_pic_08.png"},
		{FieldName: "bingo_cherry_h_pic_09", Path: "bingo/cherry/bingo_cherry_h_pic_09.png"},
		{FieldName: "bingo_cherry_h_pic_10", Path: "bingo/cherry/bingo_cherry_h_pic_10.png"},
		{FieldName: "bingo_cherry_h_pic_11", Path: "bingo/cherry/bingo_cherry_h_pic_11.png"},
		{FieldName: "bingo_cherry_h_pic_12", Path: "bingo/cherry/bingo_cherry_h_pic_12.png"},
		{FieldName: "bingo_cherry_h_pic_14", Path: "bingo/cherry/bingo_cherry_h_pic_14.png"},
		{FieldName: "bingo_cherry_h_pic_15", Path: "bingo/cherry/bingo_cherry_h_pic_15.png"},
		{FieldName: "bingo_cherry_h_pic_17", Path: "bingo/cherry/bingo_cherry_h_pic_17.png"},
		{FieldName: "bingo_cherry_h_pic_18", Path: "bingo/cherry/bingo_cherry_h_pic_18.png"},
		{FieldName: "bingo_cherry_h_pic_19", Path: "bingo/cherry/bingo_cherry_h_pic_19.png"},
		{FieldName: "bingo_cherry_g_pic_01", Path: "bingo/cherry/bingo_cherry_g_pic_01.png"},
		{FieldName: "bingo_cherry_g_pic_02", Path: "bingo/cherry/bingo_cherry_g_pic_02.png"},
		{FieldName: "bingo_cherry_g_pic_03", Path: "bingo/cherry/bingo_cherry_g_pic_03.png"},
		{FieldName: "bingo_cherry_g_pic_04", Path: "bingo/cherry/bingo_cherry_g_pic_04.jpg"},
		{FieldName: "bingo_cherry_g_pic_05", Path: "bingo/cherry/bingo_cherry_g_pic_05.png"},
		{FieldName: "bingo_cherry_g_pic_06", Path: "bingo/cherry/bingo_cherry_g_pic_06.png"},
		{FieldName: "bingo_cherry_c_pic_01", Path: "bingo/cherry/bingo_cherry_c_pic_01.png"},
		{FieldName: "bingo_cherry_c_pic_02", Path: "bingo/cherry/bingo_cherry_c_pic_02.png"},
		{FieldName: "bingo_cherry_c_pic_03", Path: "bingo/cherry/bingo_cherry_c_pic_03.png"},
		{FieldName: "bingo_cherry_c_pic_04", Path: "bingo/cherry/bingo_cherry_c_pic_04.png"},
		{FieldName: "bingo_cherry_h_ani_02", Path: "bingo/cherry/bingo_cherry_h_ani_02.png"},
		{FieldName: "bingo_cherry_h_ani_03", Path: "bingo/cherry/bingo_cherry_h_ani_03.png"},
		{FieldName: "bingo_cherry_g_ani_01", Path: "bingo/cherry/bingo_cherry_g_ani_01.png"},
		{FieldName: "bingo_cherry_g_ani_02", Path: "bingo/cherry/bingo_cherry_g_ani_02.png"},
		{FieldName: "bingo_cherry_c_ani_01", Path: "bingo/cherry/bingo_cherry_c_ani_01.png"},
		{FieldName: "bingo_cherry_c_ani_02", Path: "bingo/cherry/bingo_cherry_c_ani_02.png"},
	}
)

// GetBingoPanel 賓果遊戲
func (s *SystemTable) GetBingoPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/bingo/" + id)
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
		picMap := BuildPictureMap(bingoPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "bingo", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(bingoPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "bingo", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// fields = []string{
// 	// 賓果遊戲自定義
// 	// 音樂
// 	"bingo_bgm_start",  // 遊戲開始
// 	"bingo_bgm_gaming", // 遊戲進行中
// 	"bingo_bgm_end",    // 遊戲結束

// 	"bingo_classic_h_pic_01",
// 	"bingo_classic_h_pic_02",
// 	"bingo_classic_h_pic_03",
// 	"bingo_classic_h_pic_04",
// 	"bingo_classic_h_pic_05",
// 	"bingo_classic_h_pic_06",
// 	"bingo_classic_h_pic_07",
// 	"bingo_classic_h_pic_08",
// 	"bingo_classic_h_pic_09",
// 	"bingo_classic_h_pic_10",
// 	"bingo_classic_h_pic_11",
// 	"bingo_classic_h_pic_12",
// 	"bingo_classic_h_pic_13",
// 	"bingo_classic_h_pic_14",
// 	"bingo_classic_h_pic_15",
// 	"bingo_classic_h_pic_16",
// 	"bingo_classic_g_pic_01",
// 	"bingo_classic_g_pic_02",
// 	"bingo_classic_g_pic_03",
// 	"bingo_classic_g_pic_04",
// 	"bingo_classic_g_pic_05",
// 	"bingo_classic_g_pic_06",
// 	"bingo_classic_g_pic_07",
// 	"bingo_classic_g_pic_08",
// 	"bingo_classic_c_pic_01",
// 	"bingo_classic_c_pic_02",
// 	"bingo_classic_c_pic_03",
// 	"bingo_classic_c_pic_04",
// 	"bingo_classic_h_ani_01",
// 	"bingo_classic_g_ani_01",
// 	"bingo_classic_c_ani_01",
// 	"bingo_classic_c_ani_02",

// 	"bingo_newyear_dragon_h_pic_01",
// 	"bingo_newyear_dragon_h_pic_02",
// 	"bingo_newyear_dragon_h_pic_03",
// 	"bingo_newyear_dragon_h_pic_04",
// 	"bingo_newyear_dragon_h_pic_05",
// 	"bingo_newyear_dragon_h_pic_06",
// 	"bingo_newyear_dragon_h_pic_07",
// 	"bingo_newyear_dragon_h_pic_08",
// 	"bingo_newyear_dragon_h_pic_09",
// 	"bingo_newyear_dragon_h_pic_10",
// 	"bingo_newyear_dragon_h_pic_11",
// 	"bingo_newyear_dragon_h_pic_12",
// 	"bingo_newyear_dragon_h_pic_13",
// 	"bingo_newyear_dragon_h_pic_14",
// 	"bingo_newyear_dragon_h_pic_16",
// 	"bingo_newyear_dragon_h_pic_17",
// 	"bingo_newyear_dragon_h_pic_18",
// 	"bingo_newyear_dragon_h_pic_19",
// 	"bingo_newyear_dragon_h_pic_20",
// 	"bingo_newyear_dragon_h_pic_21",
// 	"bingo_newyear_dragon_h_pic_22",
// 	"bingo_newyear_dragon_g_pic_01",
// 	"bingo_newyear_dragon_g_pic_02",
// 	"bingo_newyear_dragon_g_pic_03",
// 	"bingo_newyear_dragon_g_pic_04",
// 	"bingo_newyear_dragon_g_pic_05",
// 	"bingo_newyear_dragon_g_pic_06",
// 	"bingo_newyear_dragon_g_pic_07",
// 	"bingo_newyear_dragon_g_pic_08",
// 	"bingo_newyear_dragon_c_pic_01",
// 	"bingo_newyear_dragon_c_pic_02",
// 	"bingo_newyear_dragon_c_pic_03",
// 	"bingo_newyear_dragon_h_ani_01",
// 	"bingo_newyear_dragon_g_ani_01",
// 	"bingo_newyear_dragon_c_ani_01",
// 	"bingo_newyear_dragon_c_ani_02",
// 	"bingo_newyear_dragon_c_ani_03",

// 	"bingo_cherry_h_pic_01",
// 	"bingo_cherry_h_pic_02",
// 	"bingo_cherry_h_pic_03",
// 	"bingo_cherry_h_pic_04",
// 	"bingo_cherry_h_pic_05",
// 	"bingo_cherry_h_pic_06",
// 	"bingo_cherry_h_pic_07",
// 	"bingo_cherry_h_pic_08",
// 	"bingo_cherry_h_pic_09",
// 	"bingo_cherry_h_pic_10",
// 	"bingo_cherry_h_pic_11",
// 	"bingo_cherry_h_pic_12",
// 	"bingo_cherry_h_pic_14",
// 	"bingo_cherry_h_pic_15",
// 	"bingo_cherry_h_pic_17",
// 	"bingo_cherry_h_pic_18",
// 	"bingo_cherry_h_pic_19",
// 	"bingo_cherry_g_pic_01",
// 	"bingo_cherry_g_pic_02",
// 	"bingo_cherry_g_pic_03",
// 	"bingo_cherry_g_pic_04",
// 	"bingo_cherry_g_pic_05",
// 	"bingo_cherry_g_pic_06",
// 	"bingo_cherry_c_pic_01",
// 	"bingo_cherry_c_pic_02",
// 	"bingo_cherry_c_pic_03",
// 	"bingo_cherry_c_pic_04",
// 	"bingo_cherry_h_ani_02",
// 	"bingo_cherry_h_ani_03",
// 	"bingo_cherry_g_ani_01",
// 	"bingo_cherry_g_ani_02",
// 	"bingo_cherry_c_ani_01",
// 	"bingo_cherry_c_ani_02",
// }
// update = make([]string, 300)

// pics = []string{
// 	// 賓果遊戲自定義
// 	// 音樂
// 	"bingo/%s/bgm/start.mp3",
// 	"bingo/%s/bgm/gaming.mp3",
// 	"bingo/%s/bgm/end.mp3",

// 	"bingo/classic/bingo_classic_h_pic_01.png",
// 	"bingo/classic/bingo_classic_h_pic_02.png",
// 	"bingo/classic/bingo_classic_h_pic_03.png",
// 	"bingo/classic/bingo_classic_h_pic_04.png",
// 	"bingo/classic/bingo_classic_h_pic_05.jpg",
// 	"bingo/classic/bingo_classic_h_pic_06.png",
// 	"bingo/classic/bingo_classic_h_pic_07.png",
// 	"bingo/classic/bingo_classic_h_pic_08.png",
// 	"bingo/classic/bingo_classic_h_pic_09.png",
// 	"bingo/classic/bingo_classic_h_pic_10.png",
// 	"bingo/classic/bingo_classic_h_pic_11.png",
// 	"bingo/classic/bingo_classic_h_pic_12.png",
// 	"bingo/classic/bingo_classic_h_pic_13.png",
// 	"bingo/classic/bingo_classic_h_pic_14.png",
// 	"bingo/classic/bingo_classic_h_pic_15.png",
// 	"bingo/classic/bingo_classic_h_pic_16.png",
// 	"bingo/classic/bingo_classic_g_pic_01.png",
// 	"bingo/classic/bingo_classic_g_pic_02.png",
// 	"bingo/classic/bingo_classic_g_pic_03.png",
// 	"bingo/classic/bingo_classic_g_pic_04.jpg",
// 	"bingo/classic/bingo_classic_g_pic_05.png",
// 	"bingo/classic/bingo_classic_g_pic_06.png",
// 	"bingo/classic/bingo_classic_g_pic_07.png",
// 	"bingo/classic/bingo_classic_g_pic_08.png",
// 	"bingo/classic/bingo_classic_c_pic_01.png",
// 	"bingo/classic/bingo_classic_c_pic_02.png",
// 	"bingo/classic/bingo_classic_c_pic_03.png",
// 	"bingo/classic/bingo_classic_c_pic_04.png",
// 	// "bingo/classic/bingo_classic_c_pic_05.png",
// 	"bingo/classic/bingo_classic_h_ani_01.png",
// 	"bingo/classic/bingo_classic_g_ani_01.png",
// 	"bingo/classic/bingo_classic_c_ani_01.png",
// 	"bingo/classic/bingo_classic_c_ani_02.png",

// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_01.jpg",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_02.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_03.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_04.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_05.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_06.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_07.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_08.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_09.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_10.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_11.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_12.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_13.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_14.png",
// 	// "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_15.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_16.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_17.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_18.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_19.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_20.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_21.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_22.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_01.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_02.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_03.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_04.jpg",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_05.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_06.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_07.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_08.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_01.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_02.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_03.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_h_ani_01.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_g_ani_01.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_01.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_02.png",
// 	"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_03.png",

//		"bingo/cherry/bingo_cherry_h_pic_01.png",
//		"bingo/cherry/bingo_cherry_h_pic_02.png",
//		"bingo/cherry/bingo_cherry_h_pic_03.png",
//		"bingo/cherry/bingo_cherry_h_pic_04.png",
//		"bingo/cherry/bingo_cherry_h_pic_05.jpg",
//		"bingo/cherry/bingo_cherry_h_pic_06.png",
//		"bingo/cherry/bingo_cherry_h_pic_07.png",
//		"bingo/cherry/bingo_cherry_h_pic_08.png",
//		"bingo/cherry/bingo_cherry_h_pic_09.png",
//		"bingo/cherry/bingo_cherry_h_pic_10.png",
//		"bingo/cherry/bingo_cherry_h_pic_11.png",
//		"bingo/cherry/bingo_cherry_h_pic_12.png",
//		"bingo/cherry/bingo_cherry_h_pic_14.png",
//		"bingo/cherry/bingo_cherry_h_pic_15.png",
//		"bingo/cherry/bingo_cherry_h_pic_17.png",
//		"bingo/cherry/bingo_cherry_h_pic_18.png",
//		"bingo/cherry/bingo_cherry_h_pic_19.png",
//		"bingo/cherry/bingo_cherry_g_pic_01.png",
//		"bingo/cherry/bingo_cherry_g_pic_02.png",
//		"bingo/cherry/bingo_cherry_g_pic_03.png",
//		"bingo/cherry/bingo_cherry_g_pic_04.jpg",
//		"bingo/cherry/bingo_cherry_g_pic_05.png",
//		"bingo/cherry/bingo_cherry_g_pic_06.png",
//		"bingo/cherry/bingo_cherry_c_pic_01.png",
//		"bingo/cherry/bingo_cherry_c_pic_02.png",
//		"bingo/cherry/bingo_cherry_c_pic_03.png",
//		"bingo/cherry/bingo_cherry_c_pic_04.png",
//		"bingo/cherry/bingo_cherry_h_ani_02.png",
//		"bingo/cherry/bingo_cherry_h_ani_03.png",
//		"bingo/cherry/bingo_cherry_g_ani_01.png",
//		"bingo/cherry/bingo_cherry_g_ani_02.png",
//		"bingo/cherry/bingo_cherry_c_ani_01.png",
//		"bingo/cherry/bingo_cherry_c_ani_02.png",
//	}
//

// @ummary 新增賓果遊戲資料(form-data)
// @Tags Bingo
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param max_people formData integer true "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string true "主題樣式" Enums(01_classic)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param max_number formData integer true "最大號碼" maximum(99)
// @param bingo_line formData integer true "賓果連線數"
// @param round_prize formData integer true "每輪發獎人數"
// @Router /interact/game/bingo/form [post]
func POSTBingo(ctx *gin.Context) {
}

// @Summary 編輯賓果遊戲資料(form-data)
// @Tags Bingo
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param max_people formData integer false "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_electric)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param max_number formData integer false "最大號碼" maximum(99)
// @param bingo_line formData integer false "賓果連線數"
// @param round_prize formData integer false "每輪發獎人數"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/bingo/form [put]
func PUTBingo(ctx *gin.Context) {
}

// @Summary 編輯賓果遊戲獎品資料(form-data)
// @Tags Bingo Prize
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
// @Router /interact/game/bingo/prize/form [put]
func PUTBingoPrize(ctx *gin.Context) {
}

// @Summary 刪除賓果遊戲資料(form-data)
// @Tags Bingo
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/bingo/form [delete]
func DELETEBingo(ctx *gin.Context) {
}

// @Summary 賓果遊戲JSON資料
// @Tags Bingo
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/bingo [get]
func BingoJSON(ctx *gin.Context) {
}

// @Summary 賓果遊戲獎品JSON資料
// @Tags Bingo Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/bingo/prize [get]
func BingoPrizeJSON(ctx *gin.Context) {
}

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
// 	UserID:        values.Get("user"),
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     "close",
// 	Second:        "0",
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
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            "",

// 	// 賓果遊戲
// 	MaxNumber:  values.Get("max_number"),
// 	BingoLine:  values.Get("bingo_line"),
// 	RoundPrize: values.Get("round_prize"),

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

// 	// 賓果遊戲自定義
// 	// 音樂
// 	BingoBgmStart:  update[0], // 遊戲開始
// 	BingoBgmGaming: update[1], // 遊戲進行中
// 	BingoBgmEnd:    update[2], // 遊戲結束

// 	BingoClassicHPic01: update[3],
// 	BingoClassicHPic02: update[4],
// 	BingoClassicHPic03: update[5],
// 	BingoClassicHPic04: update[6],
// 	BingoClassicHPic05: update[7],
// 	BingoClassicHPic06: update[8],
// 	BingoClassicHPic07: update[9],
// 	BingoClassicHPic08: update[10],
// 	BingoClassicHPic09: update[11],
// 	BingoClassicHPic10: update[12],
// 	BingoClassicHPic11: update[13],
// 	BingoClassicHPic12: update[14],
// 	BingoClassicHPic13: update[15],
// 	BingoClassicHPic14: update[16],
// 	BingoClassicHPic15: update[17],
// 	BingoClassicHPic16: update[18],
// 	BingoClassicGPic01: update[19],
// 	BingoClassicGPic02: update[20],
// 	BingoClassicGPic03: update[21],
// 	BingoClassicGPic04: update[22],
// 	BingoClassicGPic05: update[23],
// 	BingoClassicGPic06: update[24],
// 	BingoClassicGPic07: update[25],
// 	BingoClassicGPic08: update[26],
// 	BingoClassicCPic01: update[27],
// 	BingoClassicCPic02: update[28],
// 	BingoClassicCPic03: update[29],
// 	BingoClassicCPic04: update[30],
// 	// BingoClassicCPic05: update[30],
// 	BingoClassicHAni01: update[31],
// 	BingoClassicGAni01: update[32],
// 	BingoClassicCAni01: update[33],
// 	BingoClassicCAni02: update[34],

// 	BingoNewyearDragonHPic01: update[35],
// 	BingoNewyearDragonHPic02: update[36],
// 	BingoNewyearDragonHPic03: update[37],
// 	BingoNewyearDragonHPic04: update[38],
// 	BingoNewyearDragonHPic05: update[39],
// 	BingoNewyearDragonHPic06: update[40],
// 	BingoNewyearDragonHPic07: update[41],
// 	BingoNewyearDragonHPic08: update[42],
// 	BingoNewyearDragonHPic09: update[43],
// 	BingoNewyearDragonHPic10: update[44],
// 	BingoNewyearDragonHPic11: update[45],
// 	BingoNewyearDragonHPic12: update[46],
// 	BingoNewyearDragonHPic13: update[47],
// 	BingoNewyearDragonHPic14: update[48],
// 	// BingoNewyearDragonHPic15: update[46],
// 	BingoNewyearDragonHPic16: update[49],
// 	BingoNewyearDragonHPic17: update[50],
// 	BingoNewyearDragonHPic18: update[51],
// 	BingoNewyearDragonHPic19: update[52],
// 	BingoNewyearDragonHPic20: update[53],
// 	BingoNewyearDragonHPic21: update[54],
// 	BingoNewyearDragonHPic22: update[55],
// 	BingoNewyearDragonGPic01: update[56],
// 	BingoNewyearDragonGPic02: update[57],
// 	BingoNewyearDragonGPic03: update[58],
// 	BingoNewyearDragonGPic04: update[59],
// 	BingoNewyearDragonGPic05: update[60],
// 	BingoNewyearDragonGPic06: update[61],
// 	BingoNewyearDragonGPic07: update[62],
// 	BingoNewyearDragonGPic08: update[63],
// 	BingoNewyearDragonCPic01: update[64],
// 	BingoNewyearDragonCPic02: update[65],
// 	BingoNewyearDragonCPic03: update[66],
// 	BingoNewyearDragonHAni01: update[67],
// 	BingoNewyearDragonGAni01: update[68],
// 	BingoNewyearDragonCAni01: update[69],
// 	BingoNewyearDragonCAni02: update[70],
// 	BingoNewyearDragonCAni03: update[71],

// 	BingoCherryHPic01: update[72],
// 	BingoCherryHPic02: update[73],
// 	BingoCherryHPic03: update[74],
// 	BingoCherryHPic04: update[75],
// 	BingoCherryHPic05: update[76],
// 	BingoCherryHPic06: update[77],
// 	BingoCherryHPic07: update[78],
// 	BingoCherryHPic08: update[79],
// 	BingoCherryHPic09: update[80],
// 	BingoCherryHPic10: update[81],
// 	BingoCherryHPic11: update[82],
// 	BingoCherryHPic12: update[83],
// 	// BingoCherryHPic13: update[83],
// 	BingoCherryHPic14: update[84],
// 	BingoCherryHPic15: update[85],
// 	// BingoCherryHPic16: update[86],
// 	BingoCherryHPic17: update[86],
// 	BingoCherryHPic18: update[87],
// 	BingoCherryHPic19: update[88],
// 	BingoCherryGPic01: update[89],
// 	BingoCherryGPic02: update[90],
// 	BingoCherryGPic03: update[91],
// 	BingoCherryGPic04: update[92],
// 	BingoCherryGPic05: update[93],
// 	BingoCherryGPic06: update[94],
// 	BingoCherryCPic01: update[95],
// 	BingoCherryCPic02: update[96],
// 	BingoCherryCPic03: update[97],
// 	BingoCherryCPic04: update[98],
// 	// BingoCherryHAni01: update[100],
// 	BingoCherryHAni02: update[99],
// 	BingoCherryHAni03: update[100],
// 	BingoCherryGAni01: update[101],
// 	BingoCherryGAni02: update[102],
// 	BingoCherryCAni01: update[103],
// 	BingoCherryCAni02: update[104],
// }

// var (
// 	pics = []string{
// 		// 賓果遊戲自定義
// 		// 音樂
// 		"bingo/%s/bgm/start.mp3",
// 		"bingo/%s/bgm/gaming.mp3",
// 		"bingo/%s/bgm/end.mp3",

// 		"bingo/classic/bingo_classic_h_pic_01.png",
// 		"bingo/classic/bingo_classic_h_pic_02.png",
// 		"bingo/classic/bingo_classic_h_pic_03.png",
// 		"bingo/classic/bingo_classic_h_pic_04.png",
// 		"bingo/classic/bingo_classic_h_pic_05.jpg",
// 		"bingo/classic/bingo_classic_h_pic_06.png",
// 		"bingo/classic/bingo_classic_h_pic_07.png",
// 		"bingo/classic/bingo_classic_h_pic_08.png",
// 		"bingo/classic/bingo_classic_h_pic_09.png",
// 		"bingo/classic/bingo_classic_h_pic_10.png",
// 		"bingo/classic/bingo_classic_h_pic_11.png",
// 		"bingo/classic/bingo_classic_h_pic_12.png",
// 		"bingo/classic/bingo_classic_h_pic_13.png",
// 		"bingo/classic/bingo_classic_h_pic_14.png",
// 		"bingo/classic/bingo_classic_h_pic_15.png",
// 		"bingo/classic/bingo_classic_h_pic_16.png",
// 		"bingo/classic/bingo_classic_g_pic_01.png",
// 		"bingo/classic/bingo_classic_g_pic_02.png",
// 		"bingo/classic/bingo_classic_g_pic_03.png",
// 		"bingo/classic/bingo_classic_g_pic_04.jpg",
// 		"bingo/classic/bingo_classic_g_pic_05.png",
// 		"bingo/classic/bingo_classic_g_pic_06.png",
// 		"bingo/classic/bingo_classic_g_pic_07.png",
// 		"bingo/classic/bingo_classic_g_pic_08.png",
// 		"bingo/classic/bingo_classic_c_pic_01.png",
// 		"bingo/classic/bingo_classic_c_pic_02.png",
// 		"bingo/classic/bingo_classic_c_pic_03.png",
// 		"bingo/classic/bingo_classic_c_pic_04.png",
// 		// "bingo/classic/bingo_classic_c_pic_05.png",
// 		"bingo/classic/bingo_classic_h_ani_01.png",
// 		"bingo/classic/bingo_classic_g_ani_01.png",
// 		"bingo/classic/bingo_classic_c_ani_01.png",
// 		"bingo/classic/bingo_classic_c_ani_02.png",

// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_01.jpg",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_02.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_03.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_04.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_05.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_06.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_07.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_08.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_09.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_10.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_11.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_12.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_13.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_14.png",
// 		// "bingo/newyear_dragon/bingo_newyear_dragon_h_pic_15.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_16.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_17.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_18.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_19.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_20.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_21.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_pic_22.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_01.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_02.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_03.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_04.jpg",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_05.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_06.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_07.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_pic_08.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_01.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_02.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_pic_03.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_h_ani_01.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_g_ani_01.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_01.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_02.png",
// 		"bingo/newyear_dragon/bingo_newyear_dragon_c_ani_03.png",

// 		"bingo/cherry/bingo_cherry_h_pic_01.png",
// 		"bingo/cherry/bingo_cherry_h_pic_02.png",
// 		"bingo/cherry/bingo_cherry_h_pic_03.png",
// 		"bingo/cherry/bingo_cherry_h_pic_04.png",
// 		"bingo/cherry/bingo_cherry_h_pic_05.jpg",
// 		"bingo/cherry/bingo_cherry_h_pic_06.png",
// 		"bingo/cherry/bingo_cherry_h_pic_07.png",
// 		"bingo/cherry/bingo_cherry_h_pic_08.png",
// 		"bingo/cherry/bingo_cherry_h_pic_09.png",
// 		"bingo/cherry/bingo_cherry_h_pic_10.png",
// 		"bingo/cherry/bingo_cherry_h_pic_11.png",
// 		"bingo/cherry/bingo_cherry_h_pic_12.png",
// 		// "bingo/cherry/bingo_cherry_h_pic_13.png",
// 		"bingo/cherry/bingo_cherry_h_pic_14.png",
// 		"bingo/cherry/bingo_cherry_h_pic_15.png",
// 		// "bingo/cherry/bingo_cherry_h_pic_16.png",
// 		"bingo/cherry/bingo_cherry_h_pic_17.png",
// 		"bingo/cherry/bingo_cherry_h_pic_18.png",
// 		"bingo/cherry/bingo_cherry_h_pic_19.png",
// 		"bingo/cherry/bingo_cherry_g_pic_01.png",
// 		"bingo/cherry/bingo_cherry_g_pic_02.png",
// 		"bingo/cherry/bingo_cherry_g_pic_03.png",
// 		"bingo/cherry/bingo_cherry_g_pic_04.jpg",
// 		"bingo/cherry/bingo_cherry_g_pic_05.png",
// 		"bingo/cherry/bingo_cherry_g_pic_06.png",
// 		"bingo/cherry/bingo_cherry_c_pic_01.png",
// 		"bingo/cherry/bingo_cherry_c_pic_02.png",
// 		"bingo/cherry/bingo_cherry_c_pic_03.png",
// 		"bingo/cherry/bingo_cherry_c_pic_04.png",
// 		// "bingo/cherry/bingo_cherry_h_ani_01.png",
// 		"bingo/cherry/bingo_cherry_h_ani_02.png",
// 		"bingo/cherry/bingo_cherry_h_ani_03.png",
// 		"bingo/cherry/bingo_cherry_g_ani_01.png",
// 		"bingo/cherry/bingo_cherry_g_ani_02.png",
// 		"bingo/cherry/bingo_cherry_c_ani_01.png",
// 		"bingo/cherry/bingo_cherry_c_ani_02.png",
// 	}
// 	fields = []string{
// 		// 賓果遊戲自定義
// 		// 音樂
// 		"bingo_bgm_start",  // 遊戲開始
// 		"bingo_bgm_gaming", // 遊戲進行中
// 		"bingo_bgm_end",    // 遊戲結束

// 		"bingo_classic_h_pic_01",
// 		"bingo_classic_h_pic_02",
// 		"bingo_classic_h_pic_03",
// 		"bingo_classic_h_pic_04",
// 		"bingo_classic_h_pic_05",
// 		"bingo_classic_h_pic_06",
// 		"bingo_classic_h_pic_07",
// 		"bingo_classic_h_pic_08",
// 		"bingo_classic_h_pic_09",
// 		"bingo_classic_h_pic_10",
// 		"bingo_classic_h_pic_11",
// 		"bingo_classic_h_pic_12",
// 		"bingo_classic_h_pic_13",
// 		"bingo_classic_h_pic_14",
// 		"bingo_classic_h_pic_15",
// 		"bingo_classic_h_pic_16",
// 		"bingo_classic_g_pic_01",
// 		"bingo_classic_g_pic_02",
// 		"bingo_classic_g_pic_03",
// 		"bingo_classic_g_pic_04",
// 		"bingo_classic_g_pic_05",
// 		"bingo_classic_g_pic_06",
// 		"bingo_classic_g_pic_07",
// 		"bingo_classic_g_pic_08",
// 		"bingo_classic_c_pic_01",
// 		"bingo_classic_c_pic_02",
// 		"bingo_classic_c_pic_03",
// 		"bingo_classic_c_pic_04",
// 		// "bingo_classic_c_pic_05",
// 		"bingo_classic_h_ani_01",
// 		"bingo_classic_g_ani_01",
// 		"bingo_classic_c_ani_01",
// 		"bingo_classic_c_ani_02",

// 		"bingo_newyear_dragon_h_pic_01",
// 		"bingo_newyear_dragon_h_pic_02",
// 		"bingo_newyear_dragon_h_pic_03",
// 		"bingo_newyear_dragon_h_pic_04",
// 		"bingo_newyear_dragon_h_pic_05",
// 		"bingo_newyear_dragon_h_pic_06",
// 		"bingo_newyear_dragon_h_pic_07",
// 		"bingo_newyear_dragon_h_pic_08",
// 		"bingo_newyear_dragon_h_pic_09",
// 		"bingo_newyear_dragon_h_pic_10",
// 		"bingo_newyear_dragon_h_pic_11",
// 		"bingo_newyear_dragon_h_pic_12",
// 		"bingo_newyear_dragon_h_pic_13",
// 		"bingo_newyear_dragon_h_pic_14",
// 		// "bingo_newyear_dragon_h_pic_15",
// 		"bingo_newyear_dragon_h_pic_16",
// 		"bingo_newyear_dragon_h_pic_17",
// 		"bingo_newyear_dragon_h_pic_18",
// 		"bingo_newyear_dragon_h_pic_19",
// 		"bingo_newyear_dragon_h_pic_20",
// 		"bingo_newyear_dragon_h_pic_21",
// 		"bingo_newyear_dragon_h_pic_22",
// 		"bingo_newyear_dragon_g_pic_01",
// 		"bingo_newyear_dragon_g_pic_02",
// 		"bingo_newyear_dragon_g_pic_03",
// 		"bingo_newyear_dragon_g_pic_04",
// 		"bingo_newyear_dragon_g_pic_05",
// 		"bingo_newyear_dragon_g_pic_06",
// 		"bingo_newyear_dragon_g_pic_07",
// 		"bingo_newyear_dragon_g_pic_08",
// 		"bingo_newyear_dragon_c_pic_01",
// 		"bingo_newyear_dragon_c_pic_02",
// 		"bingo_newyear_dragon_c_pic_03",
// 		"bingo_newyear_dragon_h_ani_01",
// 		"bingo_newyear_dragon_g_ani_01",
// 		"bingo_newyear_dragon_c_ani_01",
// 		"bingo_newyear_dragon_c_ani_02",
// 		"bingo_newyear_dragon_c_ani_03",

// 		"bingo_cherry_h_pic_01",
// 		"bingo_cherry_h_pic_02",
// 		"bingo_cherry_h_pic_03",
// 		"bingo_cherry_h_pic_04",
// 		"bingo_cherry_h_pic_05",
// 		"bingo_cherry_h_pic_06",
// 		"bingo_cherry_h_pic_07",
// 		"bingo_cherry_h_pic_08",
// 		"bingo_cherry_h_pic_09",
// 		"bingo_cherry_h_pic_10",
// 		"bingo_cherry_h_pic_11",
// 		"bingo_cherry_h_pic_12",
// 		// "bingo_cherry_h_pic_13",
// 		"bingo_cherry_h_pic_14",
// 		"bingo_cherry_h_pic_15",
// 		// "bingo_cherry_h_pic_16",
// 		"bingo_cherry_h_pic_17",
// 		"bingo_cherry_h_pic_18",
// 		"bingo_cherry_h_pic_19",
// 		"bingo_cherry_g_pic_01",
// 		"bingo_cherry_g_pic_02",
// 		"bingo_cherry_g_pic_03",
// 		"bingo_cherry_g_pic_04",
// 		"bingo_cherry_g_pic_05",
// 		"bingo_cherry_g_pic_06",
// 		"bingo_cherry_c_pic_01",
// 		"bingo_cherry_c_pic_02",
// 		"bingo_cherry_c_pic_03",
// 		"bingo_cherry_c_pic_04",
// 		// "bingo_cherry_h_ani_01",
// 		"bingo_cherry_h_ani_02",
// 		"bingo_cherry_h_ani_03",
// 		"bingo_cherry_g_ani_01",
// 		"bingo_cherry_g_ani_02",
// 		"bingo_cherry_c_ani_01",
// 		"bingo_cherry_c_ani_02",
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
// 	LimitTime:     "",
// 	Second:        "",
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
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            "",

// 	// 賓果遊戲
// 	MaxNumber:  values.Get("max_number"),
// 	BingoLine:  values.Get("bingo_line"),
// 	RoundPrize: values.Get("round_prize"),

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

// 	// 賓果遊戲自定義
// 	// 音樂
// 	BingoBgmStart:  update[0], // 遊戲開始
// 	BingoBgmGaming: update[1], // 遊戲進行中
// 	BingoBgmEnd:    update[2], // 遊戲結束

// 	BingoClassicHPic01: update[3],
// 	BingoClassicHPic02: update[4],
// 	BingoClassicHPic03: update[5],
// 	BingoClassicHPic04: update[6],
// 	BingoClassicHPic05: update[7],
// 	BingoClassicHPic06: update[8],
// 	BingoClassicHPic07: update[9],
// 	BingoClassicHPic08: update[10],
// 	BingoClassicHPic09: update[11],
// 	BingoClassicHPic10: update[12],
// 	BingoClassicHPic11: update[13],
// 	BingoClassicHPic12: update[14],
// 	BingoClassicHPic13: update[15],
// 	BingoClassicHPic14: update[16],
// 	BingoClassicHPic15: update[17],
// 	BingoClassicHPic16: update[18],
// 	BingoClassicGPic01: update[19],
// 	BingoClassicGPic02: update[20],
// 	BingoClassicGPic03: update[21],
// 	BingoClassicGPic04: update[22],
// 	BingoClassicGPic05: update[23],
// 	BingoClassicGPic06: update[24],
// 	BingoClassicGPic07: update[25],
// 	BingoClassicGPic08: update[26],
// 	BingoClassicCPic01: update[27],
// 	BingoClassicCPic02: update[28],
// 	BingoClassicCPic03: update[29],
// 	BingoClassicCPic04: update[30],
// 	// BingoClassicCPic05: update[30],
// 	BingoClassicHAni01: update[31],
// 	BingoClassicGAni01: update[32],
// 	BingoClassicCAni01: update[33],
// 	BingoClassicCAni02: update[34],

// 	BingoNewyearDragonHPic01: update[35],
// 	BingoNewyearDragonHPic02: update[36],
// 	BingoNewyearDragonHPic03: update[37],
// 	BingoNewyearDragonHPic04: update[38],
// 	BingoNewyearDragonHPic05: update[39],
// 	BingoNewyearDragonHPic06: update[40],
// 	BingoNewyearDragonHPic07: update[41],
// 	BingoNewyearDragonHPic08: update[42],
// 	BingoNewyearDragonHPic09: update[43],
// 	BingoNewyearDragonHPic10: update[44],
// 	BingoNewyearDragonHPic11: update[45],
// 	BingoNewyearDragonHPic12: update[46],
// 	BingoNewyearDragonHPic13: update[47],
// 	BingoNewyearDragonHPic14: update[48],
// 	// BingoNewyearDragonHPic15: update[46],
// 	BingoNewyearDragonHPic16: update[49],
// 	BingoNewyearDragonHPic17: update[50],
// 	BingoNewyearDragonHPic18: update[51],
// 	BingoNewyearDragonHPic19: update[52],
// 	BingoNewyearDragonHPic20: update[53],
// 	BingoNewyearDragonHPic21: update[54],
// 	BingoNewyearDragonHPic22: update[55],
// 	BingoNewyearDragonGPic01: update[56],
// 	BingoNewyearDragonGPic02: update[57],
// 	BingoNewyearDragonGPic03: update[58],
// 	BingoNewyearDragonGPic04: update[59],
// 	BingoNewyearDragonGPic05: update[60],
// 	BingoNewyearDragonGPic06: update[61],
// 	BingoNewyearDragonGPic07: update[62],
// 	BingoNewyearDragonGPic08: update[63],
// 	BingoNewyearDragonCPic01: update[64],
// 	BingoNewyearDragonCPic02: update[65],
// 	BingoNewyearDragonCPic03: update[66],
// 	BingoNewyearDragonHAni01: update[67],
// 	BingoNewyearDragonGAni01: update[68],
// 	BingoNewyearDragonCAni01: update[69],
// 	BingoNewyearDragonCAni02: update[70],
// 	BingoNewyearDragonCAni03: update[71],

// 	BingoCherryHPic01: update[72],
// 	BingoCherryHPic02: update[73],
// 	BingoCherryHPic03: update[74],
// 	BingoCherryHPic04: update[75],
// 	BingoCherryHPic05: update[76],
// 	BingoCherryHPic06: update[77],
// 	BingoCherryHPic07: update[78],
// 	BingoCherryHPic08: update[79],
// 	BingoCherryHPic09: update[80],
// 	BingoCherryHPic10: update[81],
// 	BingoCherryHPic11: update[82],
// 	BingoCherryHPic12: update[83],
// 	// BingoCherryHPic13: update[83],
// 	BingoCherryHPic14: update[84],
// 	BingoCherryHPic15: update[85],
// 	// BingoCherryHPic16: update[86],
// 	BingoCherryHPic17: update[86],
// 	BingoCherryHPic18: update[87],
// 	BingoCherryHPic19: update[88],
// 	BingoCherryGPic01: update[89],
// 	BingoCherryGPic02: update[90],
// 	BingoCherryGPic03: update[91],
// 	BingoCherryGPic04: update[92],
// 	BingoCherryGPic05: update[93],
// 	BingoCherryGPic06: update[94],
// 	BingoCherryCPic01: update[95],
// 	BingoCherryCPic02: update[96],
// 	BingoCherryCPic03: update[97],
// 	BingoCherryCPic04: update[98],
// 	// BingoCherryHAni01: update[100],
// 	BingoCherryHAni02: update[99],
// 	BingoCherryHAni03: update[100],
// 	BingoCherryGAni01: update[101],
// 	BingoCherryGAni02: update[102],
// 	BingoCherryCAni01: update[103],
// 	BingoCherryCAni02: update[104],
// }
