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
	qaPictureFields = []PictureField{
		{FieldName: "qa_bgm_start", Path: "qa/%s/bgm/start.mp3"},
		{FieldName: "qa_bgm_gaming", Path: "qa/%s/bgm/gaming.mp3"},
		{FieldName: "qa_bgm_end", Path: "qa/%s/bgm/end.mp3"},

		{FieldName: "qa_classic_h_pic_01", Path: "qa/classic/qa_classic_h_pic_01.png"},
		{FieldName: "qa_classic_h_pic_02", Path: "qa/classic/qa_classic_h_pic_02.png"},
		{FieldName: "qa_classic_h_pic_03", Path: "qa/classic/qa_classic_h_pic_03.jpg"},
		{FieldName: "qa_classic_h_pic_04", Path: "qa/classic/qa_classic_h_pic_04.jpg"},
		{FieldName: "qa_classic_h_pic_05", Path: "qa/classic/qa_classic_h_pic_05.png"},
		{FieldName: "qa_classic_h_pic_06", Path: "qa/classic/qa_classic_h_pic_06.png"},
		{FieldName: "qa_classic_h_pic_07", Path: "qa/classic/qa_classic_h_pic_07.png"},
		{FieldName: "qa_classic_h_pic_08", Path: "qa/classic/qa_classic_h_pic_08.png"},
		{FieldName: "qa_classic_h_pic_09", Path: "qa/classic/qa_classic_h_pic_09.png"},
		{FieldName: "qa_classic_h_pic_10", Path: "qa/classic/qa_classic_h_pic_10.png"},
		{FieldName: "qa_classic_h_pic_11", Path: "qa/classic/qa_classic_h_pic_11.png"},
		{FieldName: "qa_classic_h_pic_12", Path: "qa/classic/qa_classic_h_pic_12.png"},
		{FieldName: "qa_classic_h_pic_13", Path: "qa/classic/qa_classic_h_pic_13.png"},
		{FieldName: "qa_classic_h_pic_14", Path: "qa/classic/qa_classic_h_pic_14.png"},
		{FieldName: "qa_classic_h_pic_15", Path: "qa/classic/qa_classic_h_pic_15.png"},
		{FieldName: "qa_classic_h_pic_16", Path: "qa/classic/qa_classic_h_pic_16.png"},
		{FieldName: "qa_classic_h_pic_17", Path: "qa/classic/qa_classic_h_pic_17.png"},
		{FieldName: "qa_classic_h_pic_18", Path: "qa/classic/qa_classic_h_pic_18.png"},
		{FieldName: "qa_classic_h_pic_19", Path: "qa/classic/qa_classic_h_pic_19.png"},
		{FieldName: "qa_classic_h_pic_20", Path: "qa/classic/qa_classic_h_pic_20.png"},
		{FieldName: "qa_classic_h_pic_21", Path: "qa/classic/qa_classic_h_pic_21.png"},
		{FieldName: "qa_classic_h_pic_22", Path: "qa/classic/qa_classic_h_pic_22.png"},
		{FieldName: "qa_classic_g_pic_01", Path: "qa/classic/qa_classic_g_pic_01.jpg"},
		{FieldName: "qa_classic_g_pic_02", Path: "qa/classic/qa_classic_g_pic_02.jpg"},
		{FieldName: "qa_classic_g_pic_03", Path: "qa/classic/qa_classic_g_pic_03.png"},
		{FieldName: "qa_classic_g_pic_04", Path: "qa/classic/qa_classic_g_pic_04.png"},
		{FieldName: "qa_classic_g_pic_05", Path: "qa/classic/qa_classic_g_pic_05.png"},
		{FieldName: "qa_classic_c_pic_01", Path: "qa/classic/qa_classic_c_pic_01.png"},
		{FieldName: "qa_classic_h_ani_01", Path: "qa/classic/qa_classic_h_ani_01.png"},
		{FieldName: "qa_classic_h_ani_02", Path: "qa/classic/qa_classic_h_ani_02.png"},
		{FieldName: "qa_classic_g_ani_01", Path: "qa/classic/qa_classic_g_ani_01.png"},
		{FieldName: "qa_classic_g_ani_02", Path: "qa/classic/qa_classic_g_ani_02.png"},

		{FieldName: "qa_electric_h_pic_01", Path: "qa/electric/qa_electric_h_pic_01.png"},
		{FieldName: "qa_electric_h_pic_02", Path: "qa/electric/qa_electric_h_pic_02.png"},
		{FieldName: "qa_electric_h_pic_03", Path: "qa/electric/qa_electric_h_pic_03.png"},
		{FieldName: "qa_electric_h_pic_04", Path: "qa/electric/qa_electric_h_pic_04.jpg"},
		{FieldName: "qa_electric_h_pic_05", Path: "qa/electric/qa_electric_h_pic_05.png"},
		{FieldName: "qa_electric_h_pic_06", Path: "qa/electric/qa_electric_h_pic_06.png"},
		{FieldName: "qa_electric_h_pic_07", Path: "qa/electric/qa_electric_h_pic_07.png"},
		{FieldName: "qa_electric_h_pic_08", Path: "qa/electric/qa_electric_h_pic_08.png"},
		{FieldName: "qa_electric_h_pic_09", Path: "qa/electric/qa_electric_h_pic_09.png"},
		{FieldName: "qa_electric_h_pic_10", Path: "qa/electric/qa_electric_h_pic_10.png"},
		{FieldName: "qa_electric_h_pic_11", Path: "qa/electric/qa_electric_h_pic_11.png"},
		{FieldName: "qa_electric_h_pic_12", Path: "qa/electric/qa_electric_h_pic_12.png"},
		{FieldName: "qa_electric_h_pic_13", Path: "qa/electric/qa_electric_h_pic_13.png"},
		{FieldName: "qa_electric_h_pic_14", Path: "qa/electric/qa_electric_h_pic_14.png"},
		{FieldName: "qa_electric_h_pic_15", Path: "qa/electric/qa_electric_h_pic_15.jpg"},
		{FieldName: "qa_electric_h_pic_16", Path: "qa/electric/qa_electric_h_pic_16.png"},
		{FieldName: "qa_electric_h_pic_17", Path: "qa/electric/qa_electric_h_pic_17.png"},
		{FieldName: "qa_electric_h_pic_18", Path: "qa/electric/qa_electric_h_pic_18.png"},
		{FieldName: "qa_electric_h_pic_19", Path: "qa/electric/qa_electric_h_pic_19.png"},
		{FieldName: "qa_electric_h_pic_20", Path: "qa/electric/qa_electric_h_pic_20.jpg"},
		{FieldName: "qa_electric_h_pic_21", Path: "qa/electric/qa_electric_h_pic_21.png"},
		{FieldName: "qa_electric_h_pic_22", Path: "qa/electric/qa_electric_h_pic_22.png"},
		{FieldName: "qa_electric_h_pic_23", Path: "qa/electric/qa_electric_h_pic_23.png"},
		{FieldName: "qa_electric_h_pic_24", Path: "qa/electric/qa_electric_h_pic_24.png"},
		{FieldName: "qa_electric_h_pic_25", Path: "qa/electric/qa_electric_h_pic_25.png"},
		{FieldName: "qa_electric_h_pic_26", Path: "qa/electric/qa_electric_h_pic_26.png"},
		{FieldName: "qa_electric_g_pic_01", Path: "qa/electric/qa_electric_g_pic_01.png"},
		{FieldName: "qa_electric_g_pic_02", Path: "qa/electric/qa_electric_g_pic_02.png"},
		{FieldName: "qa_electric_g_pic_03", Path: "qa/electric/qa_electric_g_pic_03.png"},
		{FieldName: "qa_electric_g_pic_04", Path: "qa/electric/qa_electric_g_pic_04.png"},
		{FieldName: "qa_electric_g_pic_05", Path: "qa/electric/qa_electric_g_pic_05.jpg"},
		{FieldName: "qa_electric_g_pic_06", Path: "qa/electric/qa_electric_g_pic_06.png"},
		{FieldName: "qa_electric_g_pic_07", Path: "qa/electric/qa_electric_g_pic_07.jpg"},
		{FieldName: "qa_electric_g_pic_08", Path: "qa/electric/qa_electric_g_pic_08.png"},
		{FieldName: "qa_electric_g_pic_09", Path: "qa/electric/qa_electric_g_pic_09.png"},
		{FieldName: "qa_electric_c_pic_01", Path: "qa/electric/qa_electric_c_pic_01.png"},
		{FieldName: "qa_electric_h_ani_01", Path: "qa/electric/qa_electric_h_ani_01.png"},
		{FieldName: "qa_electric_h_ani_02", Path: "qa/electric/qa_electric_h_ani_02.png"},
		{FieldName: "qa_electric_h_ani_03", Path: "qa/electric/qa_electric_h_ani_03.png"},
		{FieldName: "qa_electric_h_ani_04", Path: "qa/electric/qa_electric_h_ani_04.png"},
		{FieldName: "qa_electric_h_ani_05", Path: "qa/electric/qa_electric_h_ani_05.png"},
		{FieldName: "qa_electric_g_ani_01", Path: "qa/electric/qa_electric_g_ani_01.png"},
		{FieldName: "qa_electric_g_ani_02", Path: "qa/electric/qa_electric_g_ani_02.png"},
		{FieldName: "qa_electric_c_ani_01", Path: "qa/electric/qa_electric_c_ani_01.png"},

		{FieldName: "qa_moonfestival_h_pic_01", Path: "qa/moonfestival/qa_moonfestival_h_pic_01.png"},
		{FieldName: "qa_moonfestival_h_pic_02", Path: "qa/moonfestival/qa_moonfestival_h_pic_02.png"},
		{FieldName: "qa_moonfestival_h_pic_03", Path: "qa/moonfestival/qa_moonfestival_h_pic_03.png"},
		{FieldName: "qa_moonfestival_h_pic_04", Path: "qa/moonfestival/qa_moonfestival_h_pic_04.png"},
		{FieldName: "qa_moonfestival_h_pic_05", Path: "qa/moonfestival/qa_moonfestival_h_pic_05.jpg"},
		{FieldName: "qa_moonfestival_h_pic_06", Path: "qa/moonfestival/qa_moonfestival_h_pic_06.png"},
		{FieldName: "qa_moonfestival_h_pic_07", Path: "qa/moonfestival/qa_moonfestival_h_pic_07.png"},
		{FieldName: "qa_moonfestival_h_pic_08", Path: "qa/moonfestival/qa_moonfestival_h_pic_08.png"},
		{FieldName: "qa_moonfestival_h_pic_09", Path: "qa/moonfestival/qa_moonfestival_h_pic_09.png"},
		{FieldName: "qa_moonfestival_h_pic_10", Path: "qa/moonfestival/qa_moonfestival_h_pic_10.png"},
		{FieldName: "qa_moonfestival_h_pic_11", Path: "qa/moonfestival/qa_moonfestival_h_pic_11.png"},
		{FieldName: "qa_moonfestival_h_pic_12", Path: "qa/moonfestival/qa_moonfestival_h_pic_12.png"},
		{FieldName: "qa_moonfestival_h_pic_13", Path: "qa/moonfestival/qa_moonfestival_h_pic_13.png"},
		{FieldName: "qa_moonfestival_h_pic_14", Path: "qa/moonfestival/qa_moonfestival_h_pic_14.png"},
		{FieldName: "qa_moonfestival_h_pic_15", Path: "qa/moonfestival/qa_moonfestival_h_pic_15.png"},
		{FieldName: "qa_moonfestival_h_pic_16", Path: "qa/moonfestival/qa_moonfestival_h_pic_16.png"},
		{FieldName: "qa_moonfestival_h_pic_17", Path: "qa/moonfestival/qa_moonfestival_h_pic_17.png"},
		{FieldName: "qa_moonfestival_h_pic_18", Path: "qa/moonfestival/qa_moonfestival_h_pic_18.png"},
		{FieldName: "qa_moonfestival_h_pic_19", Path: "qa/moonfestival/qa_moonfestival_h_pic_19.png"},
		{FieldName: "qa_moonfestival_h_pic_20", Path: "qa/moonfestival/qa_moonfestival_h_pic_20.png"},
		{FieldName: "qa_moonfestival_h_pic_21", Path: "qa/moonfestival/qa_moonfestival_h_pic_21.png"},
		{FieldName: "qa_moonfestival_h_pic_22", Path: "qa/moonfestival/qa_moonfestival_h_pic_22.png"},
		{FieldName: "qa_moonfestival_h_pic_23", Path: "qa/moonfestival/qa_moonfestival_h_pic_23.png"},
		{FieldName: "qa_moonfestival_h_pic_24", Path: "qa/moonfestival/qa_moonfestival_h_pic_24.png"},
		{FieldName: "qa_moonfestival_g_pic_01", Path: "qa/moonfestival/qa_moonfestival_g_pic_01.png"},
		{FieldName: "qa_moonfestival_g_pic_02", Path: "qa/moonfestival/qa_moonfestival_g_pic_02.png"},
		{FieldName: "qa_moonfestival_g_pic_03", Path: "qa/moonfestival/qa_moonfestival_g_pic_03.jpg"},
		{FieldName: "qa_moonfestival_g_pic_04", Path: "qa/moonfestival/qa_moonfestival_g_pic_04.png"},
		{FieldName: "qa_moonfestival_g_pic_05", Path: "qa/moonfestival/qa_moonfestival_g_pic_05.png"},
		{FieldName: "qa_moonfestival_c_pic_01", Path: "qa/moonfestival/qa_moonfestival_c_pic_01.png"},
		{FieldName: "qa_moonfestival_c_pic_02", Path: "qa/moonfestival/qa_moonfestival_c_pic_02.png"},
		{FieldName: "qa_moonfestival_c_pic_03", Path: "qa/moonfestival/qa_moonfestival_c_pic_03.png"},
		{FieldName: "qa_moonfestival_h_ani_01", Path: "qa/moonfestival/qa_moonfestival_h_ani_01.png"},
		{FieldName: "qa_moonfestival_h_ani_02", Path: "qa/moonfestival/qa_moonfestival_h_ani_02.png"},
		{FieldName: "qa_moonfestival_g_ani_01", Path: "qa/moonfestival/qa_moonfestival_g_ani_01.png"},
		{FieldName: "qa_moonfestival_g_ani_02", Path: "qa/moonfestival/qa_moonfestival_g_ani_02.png"},
		{FieldName: "qa_moonfestival_g_ani_03", Path: "qa/moonfestival/qa_moonfestival_g_ani_03.png"},

		{FieldName: "qa_newyear_dragon_h_pic_01", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_01.png"},
		{FieldName: "qa_newyear_dragon_h_pic_02", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_02.png"},
		{FieldName: "qa_newyear_dragon_h_pic_03", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_03.png"},
		{FieldName: "qa_newyear_dragon_h_pic_04", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_04.png"},
		{FieldName: "qa_newyear_dragon_h_pic_05", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_05.png"},
		{FieldName: "qa_newyear_dragon_h_pic_06", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_06.jpg"},
		{FieldName: "qa_newyear_dragon_h_pic_07", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_07.png"},
		{FieldName: "qa_newyear_dragon_h_pic_08", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_08.png"},
		{FieldName: "qa_newyear_dragon_h_pic_09", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_09.png"},
		{FieldName: "qa_newyear_dragon_h_pic_10", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_10.png"},
		{FieldName: "qa_newyear_dragon_h_pic_11", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_11.png"},
		{FieldName: "qa_newyear_dragon_h_pic_12", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_12.png"},
		{FieldName: "qa_newyear_dragon_h_pic_13", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_13.png"},
		{FieldName: "qa_newyear_dragon_h_pic_14", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_14.png"},
		{FieldName: "qa_newyear_dragon_h_pic_15", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_15.png"},
		{FieldName: "qa_newyear_dragon_h_pic_16", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_16.png"},
		{FieldName: "qa_newyear_dragon_h_pic_17", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_17.png"},
		{FieldName: "qa_newyear_dragon_h_pic_18", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_18.png"},
		{FieldName: "qa_newyear_dragon_h_pic_19", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_19.png"},
		{FieldName: "qa_newyear_dragon_h_pic_20", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_20.png"},
		{FieldName: "qa_newyear_dragon_h_pic_21", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_21.png"},
		{FieldName: "qa_newyear_dragon_h_pic_22", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_22.png"},
		{FieldName: "qa_newyear_dragon_h_pic_23", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_23.png"},
		{FieldName: "qa_newyear_dragon_h_pic_24", Path: "qa/newyear_dragon/qa_newyear_dragon_h_pic_24.png"},
		{FieldName: "qa_newyear_dragon_g_pic_01", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_01.png"},
		{FieldName: "qa_newyear_dragon_g_pic_02", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_02.png"},
		{FieldName: "qa_newyear_dragon_g_pic_03", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_03.png"},
		{FieldName: "qa_newyear_dragon_g_pic_04", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_04.png"},
		{FieldName: "qa_newyear_dragon_g_pic_05", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_05.jpg"},
		{FieldName: "qa_newyear_dragon_g_pic_06", Path: "qa/newyear_dragon/qa_newyear_dragon_g_pic_06.png"},
		{FieldName: "qa_newyear_dragon_c_pic_01", Path: "qa/newyear_dragon/qa_newyear_dragon_c_pic_01.png"},
		{FieldName: "qa_newyear_dragon_h_ani_01", Path: "qa/newyear_dragon/qa_newyear_dragon_h_ani_01.png"},
		{FieldName: "qa_newyear_dragon_h_ani_02", Path: "qa/newyear_dragon/qa_newyear_dragon_h_ani_02.png"},
		{FieldName: "qa_newyear_dragon_g_ani_01", Path: "qa/newyear_dragon/qa_newyear_dragon_g_ani_01.png"},
		{FieldName: "qa_newyear_dragon_g_ani_02", Path: "qa/newyear_dragon/qa_newyear_dragon_g_ani_02.png"},
		{FieldName: "qa_newyear_dragon_g_ani_03", Path: "qa/newyear_dragon/qa_newyear_dragon_g_ani_03.png"},
		{FieldName: "qa_newyear_dragon_c_ani_01", Path: "qa/newyear_dragon/qa_newyear_dragon_c_ani_01.png"},
	}
)

// GetQAPanel 快問快答
func (s *SystemTable) GetQAPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/QA/" + id)
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
		picMap := BuildPictureMap(qaPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "QA", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(qaPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "QA", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增快問快答遊戲資料(form-data)
// @Tags QA
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string true "是否限制答題時間" Enums(open, close)
// @param second formData integer true "答題秒數"
// @param qa_second formData integer true "題目顯示秒數" minimum(1)
// @param max_people formData integer true "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer true "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param first_prize formData integer true "頭獎中獎人數上限(上限為50人)" maximum(50)
// @param second_prize formData integer true "二獎中獎人數上限(上限為50人)" maximum(50)
// @param third_prize formData integer true "三獎中獎人數上限(上限為100人)" maximum(100)
// @param general_prize formData integer true "普通獎中獎人數上限(上限為800人)" maximum(800)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string true "主題樣式" Enums(01_classic, 02_electric, 03_moonfestival)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param qa_1 formData string false "題目一"
// @@@param qa_1_picture formData file false "題目一圖片"
// @param qa_1_options formData string false "題目一選項"
// @param qa_1_answer formData string false "題目一答案" maxlength(1) Enums(0, 1, 2, 3)
// @param qa_1_score formData integer false "題目一分數" minimum(1)
// @param total_qa formData integer true "總題數"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA/form [post]
func POSTQA(ctx *gin.Context) {
}

// @Summary 新增快問快答獎品資料(form-data)
// @Tags QA Prize
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
// @Router /interact/game/QA/prize/form [post]
func POSTQAPrize(ctx *gin.Context) {
}

// @Summary 編輯快問快答遊戲資料(form-data)
// @Tags QA
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string false "是否限制答題時間" Enums(open, close)
// @param second formData integer false "答題秒數"
// @param qa_second formData integer false "題目顯示秒數" minimum(1)
// @param max_people formData integer false "人數上限(依照用戶權限判斷上限)" minimum(1)
// @param people formData integer false "人數上限(依照max_people資料判斷上限)" minimum(1)
// @param first_prize formData integer false "頭獎中獎人數上限(上限為50人)" maximum(50)
// @param second_prize formData integer false "二獎中獎人數上限(上限為50人)" maximum(50)
// @param third_prize formData integer false "三獎中獎人數上限(上限為100人)" maximum(100)
// @param general_prize formData integer false "普通獎中獎人數上限(上限為800人)" maximum(800)
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_electric, 03_moonfestival)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param qa_1 formData string false "題目一"
// @@@param qa_1_picture formData file false "題目一圖片"
// @param qa_1_options formData string false "題目一選項"
// @param qa_1_answer formData string false "題目一答案" maxlength(1) Enums(0, 1, 2, 3)
// @param qa_1_score formData integer false "題目一分數" minimum(1)
// @param total_qa formData integer false "總題數"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA/form [put]
func PUTQA(ctx *gin.Context) {
}

// @Summary 編輯快問快答獎品資料(form-data)
// @Tags QA Prize
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
// @Router /interact/game/QA/prize/form [put]
func PUTQAPrize(ctx *gin.Context) {
}

// @Summary 刪除快問快答遊戲資料(form-data)
// @Tags QA
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA/form [delete]
func DELETEQA(ctx *gin.Context) {
}

// @Summary 刪除快問快答獎品資料(form-data)
// @Tags QA Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA/prize/form [delete]
func DELETEQAPrize(ctx *gin.Context) {
}

// @Summary 快問快答遊戲JSON資料
// @Tags QA
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA [get]
func QAJSON(ctx *gin.Context) {
}

// @Summary 快問快答答題紀錄JSON資料
// @Tags QA
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @param round query integer false "輪次"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /QA/record [get]
func QARecordsJSON(ctx *gin.Context) {
}

// @Summary 快問快答獎品JSON資料
// @Tags QA Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/QA/prize [get]
func QAPrizeJSON(ctx *gin.Context) {
}

// 		// 自定義圖片
// 		QABgmStart:  update[0], // 遊戲開始
// 		QABgmGaming: update[1], // 遊戲進行中
// 		QABgmEnd:    update[2], // 遊戲結束

// 		QAClassicHPic01: update[3],
// 		QAClassicHPic02: update[4],
// 		QAClassicHPic03: update[5],
// 		QAClassicHPic04: update[6],
// 		QAClassicHPic05: update[7],
// 		QAClassicHPic06: update[8],
// 		QAClassicHPic07: update[9],
// 		QAClassicHPic08: update[10],
// 		QAClassicHPic09: update[11],
// 		QAClassicHPic10: update[12],
// 		QAClassicHPic11: update[13],
// 		QAClassicHPic12: update[14],
// 		QAClassicHPic13: update[15],
// 		QAClassicHPic14: update[16],
// 		QAClassicHPic15: update[17],
// 		QAClassicHPic16: update[18],
// 		QAClassicHPic17: update[19],
// 		QAClassicHPic18: update[20],
// 		QAClassicHPic19: update[21],
// 		QAClassicHPic20: update[22],
// 		QAClassicHPic21: update[23],
// 		QAClassicHPic22: update[24],
// 		QAClassicGPic01: update[25],
// 		QAClassicGPic02: update[26],
// 		QAClassicGPic03: update[27],
// 		QAClassicGPic04: update[28],
// 		QAClassicGPic05: update[29],
// 		QAClassicCPic01: update[30],
// 		QAClassicHAni01: update[31],
// 		QAClassicHAni02: update[32],
// 		QAClassicGAni01: update[33],
// 		QAClassicGAni02: update[34],

// 		QAElectricHPic01: update[35],
// 		QAElectricHPic02: update[36],
// 		QAElectricHPic03: update[37],
// 		QAElectricHPic04: update[38],
// 		QAElectricHPic05: update[39],
// 		QAElectricHPic06: update[40],
// 		QAElectricHPic07: update[41],
// 		QAElectricHPic08: update[42],
// 		QAElectricHPic09: update[43],
// 		QAElectricHPic10: update[44],
// 		QAElectricHPic11: update[45],
// 		QAElectricHPic12: update[46],
// 		QAElectricHPic13: update[47],
// 		QAElectricHPic14: update[48],
// 		QAElectricHPic15: update[49],
// 		QAElectricHPic16: update[50],
// 		QAElectricHPic17: update[51],
// 		QAElectricHPic18: update[52],
// 		QAElectricHPic19: update[53],
// 		QAElectricHPic20: update[54],
// 		QAElectricHPic21: update[55],
// 		QAElectricHPic22: update[56],
// 		QAElectricHPic23: update[57],
// 		QAElectricHPic24: update[58],
// 		QAElectricHPic25: update[59],
// 		QAElectricHPic26: update[60],
// 		QAElectricGPic01: update[61],
// 		QAElectricGPic02: update[62],
// 		QAElectricGPic03: update[63],
// 		QAElectricGPic04: update[64],
// 		QAElectricGPic05: update[65],
// 		QAElectricGPic06: update[66],
// 		QAElectricGPic07: update[67],
// 		QAElectricGPic08: update[68],
// 		QAElectricGPic09: update[69],
// 		QAElectricCPic01: update[70],
// 		QAElectricHAni01: update[71],
// 		QAElectricHAni02: update[72],
// 		QAElectricHAni03: update[73],
// 		QAElectricHAni04: update[74],
// 		QAElectricHAni05: update[75],
// 		QAElectricGAni01: update[76],
// 		QAElectricGAni02: update[77],
// 		QAElectricCAni01: update[78],

// 		QAMoonfestivalHPic01: update[79],
// 		QAMoonfestivalHPic02: update[80],
// 		QAMoonfestivalHPic03: update[81],
// 		QAMoonfestivalHPic04: update[82],
// 		QAMoonfestivalHPic05: update[83],
// 		QAMoonfestivalHPic06: update[84],
// 		QAMoonfestivalHPic07: update[85],
// 		QAMoonfestivalHPic08: update[86],
// 		QAMoonfestivalHPic09: update[87],
// 		QAMoonfestivalHPic10: update[88],
// 		QAMoonfestivalHPic11: update[89],
// 		QAMoonfestivalHPic12: update[90],
// 		QAMoonfestivalHPic13: update[91],
// 		QAMoonfestivalHPic14: update[92],
// 		QAMoonfestivalHPic15: update[93],
// 		QAMoonfestivalHPic16: update[94],
// 		QAMoonfestivalHPic17: update[95],
// 		QAMoonfestivalHPic18: update[96],
// 		QAMoonfestivalHPic19: update[97],
// 		QAMoonfestivalHPic20: update[98],
// 		QAMoonfestivalHPic21: update[99],
// 		QAMoonfestivalHPic22: update[100],
// 		QAMoonfestivalHPic23: update[101],
// 		QAMoonfestivalHPic24: update[102],
// 		QAMoonfestivalGPic01: update[103],
// 		QAMoonfestivalGPic02: update[104],
// 		QAMoonfestivalGPic03: update[105],
// 		QAMoonfestivalGPic04: update[106],
// 		QAMoonfestivalGPic05: update[107],
// 		QAMoonfestivalCPic01: update[108],
// 		QAMoonfestivalCPic02: update[109],
// 		QAMoonfestivalCPic03: update[110],
// 		QAMoonfestivalHAni01: update[111],
// 		QAMoonfestivalHAni02: update[112],
// 		QAMoonfestivalGAni01: update[113],
// 		QAMoonfestivalGAni02: update[114],
// 		QAMoonfestivalGAni03: update[115],

// 		QANewyearDragonHPic01: update[116],
// 		QANewyearDragonHPic02: update[117],
// 		QANewyearDragonHPic03: update[118],
// 		QANewyearDragonHPic04: update[119],
// 		QANewyearDragonHPic05: update[120],
// 		QANewyearDragonHPic06: update[121],
// 		QANewyearDragonHPic07: update[122],
// 		QANewyearDragonHPic08: update[123],
// 		QANewyearDragonHPic09: update[124],
// 		QANewyearDragonHPic10: update[125],
// 		QANewyearDragonHPic11: update[126],
// 		QANewyearDragonHPic12: update[127],
// 		QANewyearDragonHPic13: update[128],
// 		QANewyearDragonHPic14: update[129],
// 		QANewyearDragonHPic15: update[130],
// 		QANewyearDragonHPic16: update[131],
// 		QANewyearDragonHPic17: update[132],
// 		QANewyearDragonHPic18: update[133],
// 		QANewyearDragonHPic19: update[134],
// 		QANewyearDragonHPic20: update[135],
// 		QANewyearDragonHPic21: update[136],
// 		QANewyearDragonHPic22: update[137],
// 		QANewyearDragonHPic23: update[138],
// 		QANewyearDragonHPic24: update[139],
// 		QANewyearDragonGPic01: update[140],
// 		QANewyearDragonGPic02: update[141],
// 		QANewyearDragonGPic03: update[142],
// 		QANewyearDragonGPic04: update[143],
// 		QANewyearDragonGPic05: update[144],
// 		QANewyearDragonGPic06: update[145],
// 		QANewyearDragonCPic01: update[146],
// 		QANewyearDragonHAni01: update[147],
// 		QANewyearDragonHAni02: update[148],
// 		QANewyearDragonGAni01: update[149],
// 		QANewyearDragonGAni02: update[150],
// 		QANewyearDragonGAni03: update[151],
// 		QANewyearDragonCAni01: update[152],

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

// update = make([]string, 300)
// qa     = make([]string, 80) // 題目設置
// index  int64
// total  string

// 判斷是否上傳excel檔案
// if values.Get("qa_excel") != "" {
// 	// 開啟excel檔
// 	file, err := excelize.OpenFile("./uploads/excel/" + values.Get("qa_excel"))
// 	if err != nil {
// 		return errors.New("錯誤: 讀取excel檔案發生問題，請重新操作")
// 	}
// 	defer func() {
// 		if err := file.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()

// 	// 題目設置
// 	for i := 0; i < 20; i++ {
// 		rowIndex := strconv.Itoa(i + 2)
// 		a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
// 		b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
// 		c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
// 		d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)
// 		e, _ := file.GetCellValue("Sheet1", "E"+rowIndex)
// 		f, _ := file.GetCellValue("Sheet1", "F"+rowIndex)
// 		g, _ := file.GetCellValue("Sheet1", "G"+rowIndex)

// 		// 某一格欄位為空，停止題目設置
// 		if a == "" || b == "" || c == "" ||
// 			d == "" || e == "" || f == "" || g == "" {
// 			break
// 		}

// 		if a != "" {
// 			total = strconv.Itoa(i + 1)
// 		}

// 		if f == "A" || f == "a" {
// 			f = "0"
// 		} else if f == "B" || f == "b" {
// 			f = "1"
// 		} else if f == "C" || f == "c" {
// 			f = "2"
// 		} else if f == "D" || f == "d" {
// 			f = "3"
// 		} else {
// 			return errors.New("錯誤: 讀取excel檔案發生問題(正確選項只能填寫ABCD)，請重新操作")
// 		}

// 		qa[index] = a
// 		qa[index+1] = strings.Join([]string{b, c, d, e}, "&&&")
// 		qa[index+2] = f
// 		qa[index+3] = g

// 		// 下一題題目設置的index間隔為4
// 		index += 4
// 	}
// } else {
// 	total = values.Get("total_qa")
// 	qa = []string{
// 		values.Get("qa_1"), values.Get("qa_1_options"), values.Get("qa_1_answer"), values.Get("qa_1_score"),
// 		values.Get("qa_2"), values.Get("qa_2_options"), values.Get("qa_2_answer"), values.Get("qa_2_score"),
// 		values.Get("qa_3"), values.Get("qa_3_options"), values.Get("qa_3_answer"), values.Get("qa_3_score"),
// 		values.Get("qa_4"), values.Get("qa_4_options"), values.Get("qa_4_answer"), values.Get("qa_4_score"),
// 		values.Get("qa_5"), values.Get("qa_5_options"), values.Get("qa_5_answer"), values.Get("qa_5_score"),
// 		values.Get("qa_6"), values.Get("qa_6_options"), values.Get("qa_6_answer"), values.Get("qa_6_score"),
// 		values.Get("qa_7"), values.Get("qa_7_options"), values.Get("qa_7_answer"), values.Get("qa_7_score"),
// 		values.Get("qa_8"), values.Get("qa_8_options"), values.Get("qa_8_answer"), values.Get("qa_8_score"),
// 		values.Get("qa_9"), values.Get("qa_9_options"), values.Get("qa_9_answer"), values.Get("qa_9_score"),
// 		values.Get("qa_10"), values.Get("qa_10_options"), values.Get("qa_10_answer"), values.Get("qa_10_score"),
// 		values.Get("qa_11"), values.Get("qa_11_options"), values.Get("qa_11_answer"), values.Get("qa_11_score"),
// 		values.Get("qa_12"), values.Get("qa_12_options"), values.Get("qa_12_answer"), values.Get("qa_12_score"),
// 		values.Get("qa_13"), values.Get("qa_13_options"), values.Get("qa_13_answer"), values.Get("qa_13_score"),
// 		values.Get("qa_14"), values.Get("qa_14_options"), values.Get("qa_14_answer"), values.Get("qa_14_score"),
// 		values.Get("qa_15"), values.Get("qa_15_options"), values.Get("qa_15_answer"), values.Get("qa_15_score"),
// 		values.Get("qa_16"), values.Get("qa_16_options"), values.Get("qa_16_answer"), values.Get("qa_16_score"),
// 		values.Get("qa_17"), values.Get("qa_17_options"), values.Get("qa_17_answer"), values.Get("qa_17_score"),
// 		values.Get("qa_18"), values.Get("qa_18_options"), values.Get("qa_18_answer"), values.Get("qa_18_score"),
// 		values.Get("qa_19"), values.Get("qa_19_options"), values.Get("qa_19_answer"), values.Get("qa_19_score"),
// 		values.Get("qa_20"), values.Get("qa_20_options"), values.Get("qa_20_answer"), values.Get("qa_20_score"),
// 	}
// }
// if total == "" {
// 	return errors.New("錯誤: 題目最少一題，請重新操作")
// }

// models.EditGameModel{
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
// 	FirstPrize:    values.Get("first_prize"),
// 	SecondPrize:   values.Get("second_prize"),
// 	ThirdPrize:    values.Get("third_prize"),
// 	GeneralPrize:  values.Get("general_prize"),
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

// 	// 題目設置
// 	QA1: qa[0], QA1Options: qa[1], QA1Answer: qa[2], QA1Score: qa[3],
// 	QA2: qa[4], QA2Options: qa[5], QA2Answer: qa[6], QA2Score: qa[7],
// 	QA3: qa[8], QA3Options: qa[9], QA3Answer: qa[10], QA3Score: qa[11],
// 	QA4: qa[12], QA4Options: qa[13], QA4Answer: qa[14], QA4Score: qa[15],
// 	QA5: qa[16], QA5Options: qa[17], QA5Answer: qa[18], QA5Score: qa[19],
// 	QA6: qa[20], QA6Options: qa[21], QA6Answer: qa[22], QA6Score: qa[23],
// 	QA7: qa[24], QA7Options: qa[25], QA7Answer: qa[26], QA7Score: qa[27],
// 	QA8: qa[28], QA8Options: qa[29], QA8Answer: qa[30], QA8Score: qa[31],
// 	QA9: qa[32], QA9Options: qa[33], QA9Answer: qa[34], QA9Score: qa[35],
// 	QA10: qa[36], QA10Options: qa[37], QA10Answer: qa[38], QA10Score: qa[39],
// 	QA11: qa[40], QA11Options: qa[41], QA11Answer: qa[42], QA11Score: qa[43],
// 	QA12: qa[44], QA12Options: qa[45], QA12Answer: qa[46], QA12Score: qa[47],
// 	QA13: qa[48], QA13Options: qa[49], QA13Answer: qa[50], QA13Score: qa[51],
// 	QA14: qa[52], QA14Options: qa[53], QA14Answer: qa[54], QA14Score: qa[55],
// 	QA15: qa[56], QA15Options: qa[57], QA15Answer: qa[58], QA15Score: qa[59],
// 	QA16: qa[60], QA16Options: qa[61], QA16Answer: qa[62], QA16Score: qa[63],
// 	QA17: qa[64], QA17Options: qa[65], QA17Answer: qa[66], QA17Score: qa[67],
// 	QA18: qa[68], QA18Options: qa[69], QA18Answer: qa[70], QA18Score: qa[71],
// 	QA19: qa[72], QA19Options: qa[73], QA19Answer: qa[74], QA19Score: qa[75],
// 	QA20: qa[76], QA20Options: qa[77], QA20Answer: qa[78], QA20Score: qa[79],
// 	TotalQA:  values.Get("total_qa"),
// 	QASecond: values.Get("qa_second"),
// }

// pics = []string{
// 快問快答自定義
// "qa/%s/bgm/start.mp3",
// "qa/%s/bgm/gaming.mp3",
// "qa/%s/bgm/end.mp3",

// "qa/classic/qa_classic_h_pic_01.png",
// "qa/classic/qa_classic_h_pic_02.png",
// "qa/classic/qa_classic_h_pic_03.jpg",
// "qa/classic/qa_classic_h_pic_04.jpg",
// "qa/classic/qa_classic_h_pic_05.png",
// "qa/classic/qa_classic_h_pic_06.png",
// "qa/classic/qa_classic_h_pic_07.png",
// "qa/classic/qa_classic_h_pic_08.png",
// "qa/classic/qa_classic_h_pic_09.png",
// "qa/classic/qa_classic_h_pic_10.png",
// "qa/classic/qa_classic_h_pic_11.png",
// "qa/classic/qa_classic_h_pic_12.png",
// "qa/classic/qa_classic_h_pic_13.png",
// "qa/classic/qa_classic_h_pic_14.png",
// "qa/classic/qa_classic_h_pic_15.png",
// "qa/classic/qa_classic_h_pic_16.png",
// "qa/classic/qa_classic_h_pic_17.png",
// "qa/classic/qa_classic_h_pic_18.png",
// "qa/classic/qa_classic_h_pic_19.png",
// "qa/classic/qa_classic_h_pic_20.png",
// "qa/classic/qa_classic_h_pic_21.png",
// "qa/classic/qa_classic_h_pic_22.png",
// "qa/classic/qa_classic_g_pic_01.jpg",
// "qa/classic/qa_classic_g_pic_02.jpg",
// "qa/classic/qa_classic_g_pic_03.png",
// "qa/classic/qa_classic_g_pic_04.png",
// "qa/classic/qa_classic_g_pic_05.png",
// "qa/classic/qa_classic_c_pic_01.png",
// "qa/classic/qa_classic_h_ani_01.png",
// "qa/classic/qa_classic_h_ani_02.png",
// "qa/classic/qa_classic_g_ani_01.png",
// "qa/classic/qa_classic_g_ani_02.png",

// "qa/electric/qa_electric_h_pic_01.png",
// "qa/electric/qa_electric_h_pic_02.png",
// "qa/electric/qa_electric_h_pic_03.png",
// "qa/electric/qa_electric_h_pic_04.jpg",
// "qa/electric/qa_electric_h_pic_05.png",
// "qa/electric/qa_electric_h_pic_06.png",
// "qa/electric/qa_electric_h_pic_07.png",
// "qa/electric/qa_electric_h_pic_08.png",
// "qa/electric/qa_electric_h_pic_09.png",
// "qa/electric/qa_electric_h_pic_10.png",
// "qa/electric/qa_electric_h_pic_11.png",
// "qa/electric/qa_electric_h_pic_12.png",
// "qa/electric/qa_electric_h_pic_13.png",
// "qa/electric/qa_electric_h_pic_14.png",
// "qa/electric/qa_electric_h_pic_15.jpg",
// "qa/electric/qa_electric_h_pic_16.png",
// "qa/electric/qa_electric_h_pic_17.png",
// "qa/electric/qa_electric_h_pic_18.png",
// "qa/electric/qa_electric_h_pic_19.png",
// "qa/electric/qa_electric_h_pic_20.jpg",
// "qa/electric/qa_electric_h_pic_21.png",
// "qa/electric/qa_electric_h_pic_22.png",
// "qa/electric/qa_electric_h_pic_23.png",
// "qa/electric/qa_electric_h_pic_24.png",
// "qa/electric/qa_electric_h_pic_25.png",
// "qa/electric/qa_electric_h_pic_26.png",
// "qa/electric/qa_electric_g_pic_01.png",
// "qa/electric/qa_electric_g_pic_02.png",
// "qa/electric/qa_electric_g_pic_03.png",
// "qa/electric/qa_electric_g_pic_04.png",
// "qa/electric/qa_electric_g_pic_05.jpg",
// "qa/electric/qa_electric_g_pic_06.png",
// "qa/electric/qa_electric_g_pic_07.jpg",
// "qa/electric/qa_electric_g_pic_08.png",
// "qa/electric/qa_electric_g_pic_09.png",
// "qa/electric/qa_electric_c_pic_01.png",
// "qa/electric/qa_electric_h_ani_01.png",
// "qa/electric/qa_electric_h_ani_02.png",
// "qa/electric/qa_electric_h_ani_03.png",
// "qa/electric/qa_electric_h_ani_04.png",
// "qa/electric/qa_electric_h_ani_05.png",
// "qa/electric/qa_electric_g_ani_01.png",
// "qa/electric/qa_electric_g_ani_02.png",
// "qa/electric/qa_electric_c_ani_01.png",

// "qa/moonfestival/qa_moonfestival_h_pic_01.png",
// "qa/moonfestival/qa_moonfestival_h_pic_02.png",
// "qa/moonfestival/qa_moonfestival_h_pic_03.png",
// "qa/moonfestival/qa_moonfestival_h_pic_04.png",
// "qa/moonfestival/qa_moonfestival_h_pic_05.jpg",
// "qa/moonfestival/qa_moonfestival_h_pic_06.png",
// "qa/moonfestival/qa_moonfestival_h_pic_07.png",
// "qa/moonfestival/qa_moonfestival_h_pic_08.png",
// "qa/moonfestival/qa_moonfestival_h_pic_09.png",
// "qa/moonfestival/qa_moonfestival_h_pic_10.png",
// "qa/moonfestival/qa_moonfestival_h_pic_11.png",
// "qa/moonfestival/qa_moonfestival_h_pic_12.png",
// "qa/moonfestival/qa_moonfestival_h_pic_13.png",
// "qa/moonfestival/qa_moonfestival_h_pic_14.png",
// "qa/moonfestival/qa_moonfestival_h_pic_15.png",
// "qa/moonfestival/qa_moonfestival_h_pic_16.png",
// "qa/moonfestival/qa_moonfestival_h_pic_17.png",
// "qa/moonfestival/qa_moonfestival_h_pic_18.png",
// "qa/moonfestival/qa_moonfestival_h_pic_19.png",
// "qa/moonfestival/qa_moonfestival_h_pic_20.png",
// "qa/moonfestival/qa_moonfestival_h_pic_21.png",
// "qa/moonfestival/qa_moonfestival_h_pic_22.png",
// "qa/moonfestival/qa_moonfestival_h_pic_23.png",
// "qa/moonfestival/qa_moonfestival_h_pic_24.png",
// "qa/moonfestival/qa_moonfestival_g_pic_01.png",
// "qa/moonfestival/qa_moonfestival_g_pic_02.png",
// "qa/moonfestival/qa_moonfestival_g_pic_03.jpg",
// "qa/moonfestival/qa_moonfestival_g_pic_04.png",
// "qa/moonfestival/qa_moonfestival_g_pic_05.png",
// "qa/moonfestival/qa_moonfestival_c_pic_01.png",
// "qa/moonfestival/qa_moonfestival_c_pic_02.png",
// "qa/moonfestival/qa_moonfestival_c_pic_03.png",
// "qa/moonfestival/qa_moonfestival_h_ani_01.png",
// "qa/moonfestival/qa_moonfestival_h_ani_02.png",
// "qa/moonfestival/qa_moonfestival_g_ani_01.png",
// "qa/moonfestival/qa_moonfestival_g_ani_02.png",
// "qa/moonfestival/qa_moonfestival_g_ani_03.png",

// "qa/newyear_dragon/qa_newyear_dragon_h_pic_01.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_02.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_03.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_04.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_05.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_06.jpg",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_07.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_08.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_09.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_10.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_11.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_12.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_13.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_14.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_15.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_16.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_17.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_18.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_19.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_20.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_21.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_22.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_23.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_pic_24.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_01.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_02.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_03.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_04.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_05.jpg",
// "qa/newyear_dragon/qa_newyear_dragon_g_pic_06.png",
// "qa/newyear_dragon/qa_newyear_dragon_c_pic_01.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_ani_01.png",
// "qa/newyear_dragon/qa_newyear_dragon_h_ani_02.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_ani_01.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_ani_02.png",
// "qa/newyear_dragon/qa_newyear_dragon_g_ani_03.png",
// "qa/newyear_dragon/qa_newyear_dragon_c_ani_01.png",
// }
// fields = []string{
// 快問快答自定義
// "qa_bgm_start",  // 遊戲開始
// "qa_bgm_gaming", // 遊戲進行中
// "qa_bgm_end",    // 遊戲結束

// "qa_classic_h_pic_01",
// "qa_classic_h_pic_02",
// "qa_classic_h_pic_03",
// "qa_classic_h_pic_04",
// "qa_classic_h_pic_05",
// "qa_classic_h_pic_06",
// "qa_classic_h_pic_07",
// "qa_classic_h_pic_08",
// "qa_classic_h_pic_09",
// "qa_classic_h_pic_10",
// "qa_classic_h_pic_11",
// "qa_classic_h_pic_12",
// "qa_classic_h_pic_13",
// "qa_classic_h_pic_14",
// "qa_classic_h_pic_15",
// "qa_classic_h_pic_16",
// "qa_classic_h_pic_17",
// "qa_classic_h_pic_18",
// "qa_classic_h_pic_19",
// "qa_classic_h_pic_20",
// "qa_classic_h_pic_21",
// "qa_classic_h_pic_22",
// "qa_classic_g_pic_01",
// "qa_classic_g_pic_02",
// "qa_classic_g_pic_03",
// "qa_classic_g_pic_04",
// "qa_classic_g_pic_05",
// "qa_classic_c_pic_01",
// "qa_classic_h_ani_01",
// "qa_classic_h_ani_02",
// "qa_classic_g_ani_01",
// "qa_classic_g_ani_02",

// "qa_electric_h_pic_01",
// "qa_electric_h_pic_02",
// "qa_electric_h_pic_03",
// "qa_electric_h_pic_04",
// "qa_electric_h_pic_05",
// "qa_electric_h_pic_06",
// "qa_electric_h_pic_07",
// "qa_electric_h_pic_08",
// "qa_electric_h_pic_09",
// "qa_electric_h_pic_10",
// "qa_electric_h_pic_11",
// "qa_electric_h_pic_12",
// "qa_electric_h_pic_13",
// "qa_electric_h_pic_14",
// "qa_electric_h_pic_15",
// "qa_electric_h_pic_16",
// "qa_electric_h_pic_17",
// "qa_electric_h_pic_18",
// "qa_electric_h_pic_19",
// "qa_electric_h_pic_20",
// "qa_electric_h_pic_21",
// "qa_electric_h_pic_22",
// "qa_electric_h_pic_23",
// "qa_electric_h_pic_24",
// "qa_electric_h_pic_25",
// "qa_electric_h_pic_26",
// "qa_electric_g_pic_01",
// "qa_electric_g_pic_02",
// "qa_electric_g_pic_03",
// "qa_electric_g_pic_04",
// "qa_electric_g_pic_05",
// "qa_electric_g_pic_06",
// "qa_electric_g_pic_07",
// "qa_electric_g_pic_08",
// "qa_electric_g_pic_09",
// "qa_electric_c_pic_01",
// "qa_electric_h_ani_01",
// "qa_electric_h_ani_02",
// "qa_electric_h_ani_03",
// "qa_electric_h_ani_04",
// "qa_electric_h_ani_05",
// "qa_electric_g_ani_01",
// "qa_electric_g_ani_02",
// "qa_electric_c_ani_01",

// "qa_moonfestival_h_pic_01",
// "qa_moonfestival_h_pic_02",
// "qa_moonfestival_h_pic_03",
// "qa_moonfestival_h_pic_04",
// "qa_moonfestival_h_pic_05",
// "qa_moonfestival_h_pic_06",
// "qa_moonfestival_h_pic_07",
// "qa_moonfestival_h_pic_08",
// "qa_moonfestival_h_pic_09",
// "qa_moonfestival_h_pic_10",
// "qa_moonfestival_h_pic_11",
// "qa_moonfestival_h_pic_12",
// "qa_moonfestival_h_pic_13",
// "qa_moonfestival_h_pic_14",
// "qa_moonfestival_h_pic_15",
// "qa_moonfestival_h_pic_16",
// "qa_moonfestival_h_pic_17",
// "qa_moonfestival_h_pic_18",
// "qa_moonfestival_h_pic_19",
// "qa_moonfestival_h_pic_20",
// "qa_moonfestival_h_pic_21",
// "qa_moonfestival_h_pic_22",
// "qa_moonfestival_h_pic_23",
// "qa_moonfestival_h_pic_24",
// "qa_moonfestival_g_pic_01",
// "qa_moonfestival_g_pic_02",
// "qa_moonfestival_g_pic_03",
// "qa_moonfestival_g_pic_04",
// "qa_moonfestival_g_pic_05",
// "qa_moonfestival_c_pic_01",
// "qa_moonfestival_c_pic_02",
// "qa_moonfestival_c_pic_03",
// "qa_moonfestival_h_ani_01",
// "qa_moonfestival_h_ani_02",
// "qa_moonfestival_g_ani_01",
// "qa_moonfestival_g_ani_02",
// "qa_moonfestival_g_ani_03",

// "qa_newyear_dragon_h_pic_01",
// "qa_newyear_dragon_h_pic_02",
// "qa_newyear_dragon_h_pic_03",
// "qa_newyear_dragon_h_pic_04",
// "qa_newyear_dragon_h_pic_05",
// "qa_newyear_dragon_h_pic_06",
// "qa_newyear_dragon_h_pic_07",
// "qa_newyear_dragon_h_pic_08",
// "qa_newyear_dragon_h_pic_09",
// "qa_newyear_dragon_h_pic_10",
// "qa_newyear_dragon_h_pic_11",
// "qa_newyear_dragon_h_pic_12",
// "qa_newyear_dragon_h_pic_13",
// "qa_newyear_dragon_h_pic_14",
// "qa_newyear_dragon_h_pic_15",
// "qa_newyear_dragon_h_pic_16",
// "qa_newyear_dragon_h_pic_17",
// "qa_newyear_dragon_h_pic_18",
// "qa_newyear_dragon_h_pic_19",
// "qa_newyear_dragon_h_pic_20",
// "qa_newyear_dragon_h_pic_21",
// "qa_newyear_dragon_h_pic_22",
// "qa_newyear_dragon_h_pic_23",
// "qa_newyear_dragon_h_pic_24",
// "qa_newyear_dragon_g_pic_01",
// "qa_newyear_dragon_g_pic_02",
// "qa_newyear_dragon_g_pic_03",
// "qa_newyear_dragon_g_pic_04",
// "qa_newyear_dragon_g_pic_05",
// "qa_newyear_dragon_g_pic_06",
// "qa_newyear_dragon_c_pic_01",
// "qa_newyear_dragon_h_ani_01",
// "qa_newyear_dragon_h_ani_02",
// "qa_newyear_dragon_g_ani_01",
// "qa_newyear_dragon_g_ani_02",
// "qa_newyear_dragon_g_ani_03",
// "qa_newyear_dragon_c_ani_01",
// }

// var (
// 	pics = []string{
// 		// 快問快答自定義
// 		"qa/%s/bgm/start.mp3",
// 		"qa/%s/bgm/gaming.mp3",
// 		"qa/%s/bgm/end.mp3",

// 		"qa/classic/qa_classic_h_pic_01.png",
// 		"qa/classic/qa_classic_h_pic_02.png",
// 		"qa/classic/qa_classic_h_pic_03.jpg",
// 		"qa/classic/qa_classic_h_pic_04.jpg",
// 		"qa/classic/qa_classic_h_pic_05.png",
// 		"qa/classic/qa_classic_h_pic_06.png",
// 		"qa/classic/qa_classic_h_pic_07.png",
// 		"qa/classic/qa_classic_h_pic_08.png",
// 		"qa/classic/qa_classic_h_pic_09.png",
// 		"qa/classic/qa_classic_h_pic_10.png",
// 		"qa/classic/qa_classic_h_pic_11.png",
// 		"qa/classic/qa_classic_h_pic_12.png",
// 		"qa/classic/qa_classic_h_pic_13.png",
// 		"qa/classic/qa_classic_h_pic_14.png",
// 		"qa/classic/qa_classic_h_pic_15.png",
// 		"qa/classic/qa_classic_h_pic_16.png",
// 		"qa/classic/qa_classic_h_pic_17.png",
// 		"qa/classic/qa_classic_h_pic_18.png",
// 		"qa/classic/qa_classic_h_pic_19.png",
// 		"qa/classic/qa_classic_h_pic_20.png",
// 		"qa/classic/qa_classic_h_pic_21.png",
// 		"qa/classic/qa_classic_h_pic_22.png",
// 		"qa/classic/qa_classic_g_pic_01.jpg",
// 		"qa/classic/qa_classic_g_pic_02.jpg",
// 		"qa/classic/qa_classic_g_pic_03.png",
// 		"qa/classic/qa_classic_g_pic_04.png",
// 		"qa/classic/qa_classic_g_pic_05.png",
// 		"qa/classic/qa_classic_c_pic_01.png",
// 		"qa/classic/qa_classic_h_ani_01.png",
// 		"qa/classic/qa_classic_h_ani_02.png",
// 		"qa/classic/qa_classic_g_ani_01.png",
// 		"qa/classic/qa_classic_g_ani_02.png",

// 		"qa/electric/qa_electric_h_pic_01.png",
// 		"qa/electric/qa_electric_h_pic_02.png",
// 		"qa/electric/qa_electric_h_pic_03.png",
// 		"qa/electric/qa_electric_h_pic_04.jpg",
// 		"qa/electric/qa_electric_h_pic_05.png",
// 		"qa/electric/qa_electric_h_pic_06.png",
// 		"qa/electric/qa_electric_h_pic_07.png",
// 		"qa/electric/qa_electric_h_pic_08.png",
// 		"qa/electric/qa_electric_h_pic_09.png",
// 		"qa/electric/qa_electric_h_pic_10.png",
// 		"qa/electric/qa_electric_h_pic_11.png",
// 		"qa/electric/qa_electric_h_pic_12.png",
// 		"qa/electric/qa_electric_h_pic_13.png",
// 		"qa/electric/qa_electric_h_pic_14.png",
// 		"qa/electric/qa_electric_h_pic_15.jpg",
// 		"qa/electric/qa_electric_h_pic_16.png",
// 		"qa/electric/qa_electric_h_pic_17.png",
// 		"qa/electric/qa_electric_h_pic_18.png",
// 		"qa/electric/qa_electric_h_pic_19.png",
// 		"qa/electric/qa_electric_h_pic_20.jpg",
// 		"qa/electric/qa_electric_h_pic_21.png",
// 		"qa/electric/qa_electric_h_pic_22.png",
// 		"qa/electric/qa_electric_h_pic_23.png",
// 		"qa/electric/qa_electric_h_pic_24.png",
// 		"qa/electric/qa_electric_h_pic_25.png",
// 		"qa/electric/qa_electric_h_pic_26.png",
// 		"qa/electric/qa_electric_g_pic_01.png",
// 		"qa/electric/qa_electric_g_pic_02.png",
// 		"qa/electric/qa_electric_g_pic_03.png",
// 		"qa/electric/qa_electric_g_pic_04.png",
// 		"qa/electric/qa_electric_g_pic_05.jpg",
// 		"qa/electric/qa_electric_g_pic_06.png",
// 		"qa/electric/qa_electric_g_pic_07.jpg",
// 		"qa/electric/qa_electric_g_pic_08.png",
// 		"qa/electric/qa_electric_g_pic_09.png",
// 		"qa/electric/qa_electric_c_pic_01.png",
// 		"qa/electric/qa_electric_h_ani_01.png",
// 		"qa/electric/qa_electric_h_ani_02.png",
// 		"qa/electric/qa_electric_h_ani_03.png",
// 		"qa/electric/qa_electric_h_ani_04.png",
// 		"qa/electric/qa_electric_h_ani_05.png",
// 		"qa/electric/qa_electric_g_ani_01.png",
// 		"qa/electric/qa_electric_g_ani_02.png",
// 		"qa/electric/qa_electric_c_ani_01.png",

// 		"qa/moonfestival/qa_moonfestival_h_pic_01.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_02.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_03.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_04.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_05.jpg",
// 		"qa/moonfestival/qa_moonfestival_h_pic_06.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_07.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_08.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_09.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_10.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_11.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_12.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_13.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_14.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_15.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_16.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_17.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_18.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_19.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_20.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_21.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_22.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_23.png",
// 		"qa/moonfestival/qa_moonfestival_h_pic_24.png",
// 		"qa/moonfestival/qa_moonfestival_g_pic_01.png",
// 		"qa/moonfestival/qa_moonfestival_g_pic_02.png",
// 		"qa/moonfestival/qa_moonfestival_g_pic_03.jpg",
// 		"qa/moonfestival/qa_moonfestival_g_pic_04.png",
// 		"qa/moonfestival/qa_moonfestival_g_pic_05.png",
// 		"qa/moonfestival/qa_moonfestival_c_pic_01.png",
// 		"qa/moonfestival/qa_moonfestival_c_pic_02.png",
// 		"qa/moonfestival/qa_moonfestival_c_pic_03.png",
// 		"qa/moonfestival/qa_moonfestival_h_ani_01.png",
// 		"qa/moonfestival/qa_moonfestival_h_ani_02.png",
// 		"qa/moonfestival/qa_moonfestival_g_ani_01.png",
// 		"qa/moonfestival/qa_moonfestival_g_ani_02.png",
// 		"qa/moonfestival/qa_moonfestival_g_ani_03.png",

// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_01.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_02.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_03.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_04.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_05.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_06.jpg",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_07.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_08.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_09.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_10.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_11.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_12.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_13.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_14.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_15.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_16.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_17.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_18.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_19.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_20.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_21.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_22.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_23.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_pic_24.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_01.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_02.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_03.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_04.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_05.jpg",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_pic_06.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_c_pic_01.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_ani_01.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_h_ani_02.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_ani_01.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_ani_02.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_g_ani_03.png",
// 		"qa/newyear_dragon/qa_newyear_dragon_c_ani_01.png",
// 	}
// 	fields = []string{
// 		// 快問快答自定義
// 		"qa_bgm_start",  // 遊戲開始
// 		"qa_bgm_gaming", // 遊戲進行中
// 		"qa_bgm_end",    // 遊戲結束

// 		"qa_classic_h_pic_01",
// 		"qa_classic_h_pic_02",
// 		"qa_classic_h_pic_03",
// 		"qa_classic_h_pic_04",
// 		"qa_classic_h_pic_05",
// 		"qa_classic_h_pic_06",
// 		"qa_classic_h_pic_07",
// 		"qa_classic_h_pic_08",
// 		"qa_classic_h_pic_09",
// 		"qa_classic_h_pic_10",
// 		"qa_classic_h_pic_11",
// 		"qa_classic_h_pic_12",
// 		"qa_classic_h_pic_13",
// 		"qa_classic_h_pic_14",
// 		"qa_classic_h_pic_15",
// 		"qa_classic_h_pic_16",
// 		"qa_classic_h_pic_17",
// 		"qa_classic_h_pic_18",
// 		"qa_classic_h_pic_19",
// 		"qa_classic_h_pic_20",
// 		"qa_classic_h_pic_21",
// 		"qa_classic_h_pic_22",
// 		"qa_classic_g_pic_01",
// 		"qa_classic_g_pic_02",
// 		"qa_classic_g_pic_03",
// 		"qa_classic_g_pic_04",
// 		"qa_classic_g_pic_05",
// 		"qa_classic_c_pic_01",
// 		"qa_classic_h_ani_01",
// 		"qa_classic_h_ani_02",
// 		"qa_classic_g_ani_01",
// 		"qa_classic_g_ani_02",

// 		"qa_electric_h_pic_01",
// 		"qa_electric_h_pic_02",
// 		"qa_electric_h_pic_03",
// 		"qa_electric_h_pic_04",
// 		"qa_electric_h_pic_05",
// 		"qa_electric_h_pic_06",
// 		"qa_electric_h_pic_07",
// 		"qa_electric_h_pic_08",
// 		"qa_electric_h_pic_09",
// 		"qa_electric_h_pic_10",
// 		"qa_electric_h_pic_11",
// 		"qa_electric_h_pic_12",
// 		"qa_electric_h_pic_13",
// 		"qa_electric_h_pic_14",
// 		"qa_electric_h_pic_15",
// 		"qa_electric_h_pic_16",
// 		"qa_electric_h_pic_17",
// 		"qa_electric_h_pic_18",
// 		"qa_electric_h_pic_19",
// 		"qa_electric_h_pic_20",
// 		"qa_electric_h_pic_21",
// 		"qa_electric_h_pic_22",
// 		"qa_electric_h_pic_23",
// 		"qa_electric_h_pic_24",
// 		"qa_electric_h_pic_25",
// 		"qa_electric_h_pic_26",
// 		"qa_electric_g_pic_01",
// 		"qa_electric_g_pic_02",
// 		"qa_electric_g_pic_03",
// 		"qa_electric_g_pic_04",
// 		"qa_electric_g_pic_05",
// 		"qa_electric_g_pic_06",
// 		"qa_electric_g_pic_07",
// 		"qa_electric_g_pic_08",
// 		"qa_electric_g_pic_09",
// 		"qa_electric_c_pic_01",
// 		"qa_electric_h_ani_01",
// 		"qa_electric_h_ani_02",
// 		"qa_electric_h_ani_03",
// 		"qa_electric_h_ani_04",
// 		"qa_electric_h_ani_05",
// 		"qa_electric_g_ani_01",
// 		"qa_electric_g_ani_02",
// 		"qa_electric_c_ani_01",

// 		"qa_moonfestival_h_pic_01",
// 		"qa_moonfestival_h_pic_02",
// 		"qa_moonfestival_h_pic_03",
// 		"qa_moonfestival_h_pic_04",
// 		"qa_moonfestival_h_pic_05",
// 		"qa_moonfestival_h_pic_06",
// 		"qa_moonfestival_h_pic_07",
// 		"qa_moonfestival_h_pic_08",
// 		"qa_moonfestival_h_pic_09",
// 		"qa_moonfestival_h_pic_10",
// 		"qa_moonfestival_h_pic_11",
// 		"qa_moonfestival_h_pic_12",
// 		"qa_moonfestival_h_pic_13",
// 		"qa_moonfestival_h_pic_14",
// 		"qa_moonfestival_h_pic_15",
// 		"qa_moonfestival_h_pic_16",
// 		"qa_moonfestival_h_pic_17",
// 		"qa_moonfestival_h_pic_18",
// 		"qa_moonfestival_h_pic_19",
// 		"qa_moonfestival_h_pic_20",
// 		"qa_moonfestival_h_pic_21",
// 		"qa_moonfestival_h_pic_22",
// 		"qa_moonfestival_h_pic_23",
// 		"qa_moonfestival_h_pic_24",
// 		"qa_moonfestival_g_pic_01",
// 		"qa_moonfestival_g_pic_02",
// 		"qa_moonfestival_g_pic_03",
// 		"qa_moonfestival_g_pic_04",
// 		"qa_moonfestival_g_pic_05",
// 		"qa_moonfestival_c_pic_01",
// 		"qa_moonfestival_c_pic_02",
// 		"qa_moonfestival_c_pic_03",
// 		"qa_moonfestival_h_ani_01",
// 		"qa_moonfestival_h_ani_02",
// 		"qa_moonfestival_g_ani_01",
// 		"qa_moonfestival_g_ani_02",
// 		"qa_moonfestival_g_ani_03",

// 		"qa_newyear_dragon_h_pic_01",
// 		"qa_newyear_dragon_h_pic_02",
// 		"qa_newyear_dragon_h_pic_03",
// 		"qa_newyear_dragon_h_pic_04",
// 		"qa_newyear_dragon_h_pic_05",
// 		"qa_newyear_dragon_h_pic_06",
// 		"qa_newyear_dragon_h_pic_07",
// 		"qa_newyear_dragon_h_pic_08",
// 		"qa_newyear_dragon_h_pic_09",
// 		"qa_newyear_dragon_h_pic_10",
// 		"qa_newyear_dragon_h_pic_11",
// 		"qa_newyear_dragon_h_pic_12",
// 		"qa_newyear_dragon_h_pic_13",
// 		"qa_newyear_dragon_h_pic_14",
// 		"qa_newyear_dragon_h_pic_15",
// 		"qa_newyear_dragon_h_pic_16",
// 		"qa_newyear_dragon_h_pic_17",
// 		"qa_newyear_dragon_h_pic_18",
// 		"qa_newyear_dragon_h_pic_19",
// 		"qa_newyear_dragon_h_pic_20",
// 		"qa_newyear_dragon_h_pic_21",
// 		"qa_newyear_dragon_h_pic_22",
// 		"qa_newyear_dragon_h_pic_23",
// 		"qa_newyear_dragon_h_pic_24",
// 		"qa_newyear_dragon_g_pic_01",
// 		"qa_newyear_dragon_g_pic_02",
// 		"qa_newyear_dragon_g_pic_03",
// 		"qa_newyear_dragon_g_pic_04",
// 		"qa_newyear_dragon_g_pic_05",
// 		"qa_newyear_dragon_g_pic_06",
// 		"qa_newyear_dragon_c_pic_01",
// 		"qa_newyear_dragon_h_ani_01",
// 		"qa_newyear_dragon_h_ani_02",
// 		"qa_newyear_dragon_g_ani_01",
// 		"qa_newyear_dragon_g_ani_02",
// 		"qa_newyear_dragon_g_ani_03",
// 		"qa_newyear_dragon_c_ani_01",
// 	}
// 	update = make([]string, 300)
// 	qa     = make([]string, 80) // 題目設置
// 	// index  int64
// 	// total  string
// )

// 判斷是否上傳excel檔案
// if values.Get("qa_excel") != "" {
// 	// 開啟excel檔
// 	file, err := excelize.OpenFile("./uploads/excel/" + values.Get("qa_excel"))
// 	if err != nil {
// 		return errors.New("錯誤: 讀取excel檔案發生問題(正確選項只能填寫ABCD)，請重新操作")
// 	}
// 	defer func() {
// 		if err := file.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()

// 	// 題目設置
// 	for i := 0; i < 20; i++ {
// 		rowIndex := strconv.Itoa(i + 2)
// 		a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
// 		b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
// 		c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
// 		d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)
// 		e, _ := file.GetCellValue("Sheet1", "E"+rowIndex)
// 		f, _ := file.GetCellValue("Sheet1", "F"+rowIndex)
// 		g, _ := file.GetCellValue("Sheet1", "G"+rowIndex)
// 		// 某一格欄位為空，停止題目設置
// 		if a == "" || b == "" || c == "" ||
// 			d == "" || e == "" || f == "" || g == "" {
// 			break
// 		}
// 		if a != "" {
// 			total = strconv.Itoa(i + 1)
// 		}
// 		if f == "A" || f == "a" {
// 			f = "0"
// 		} else if f == "B" || f == "b" {
// 			f = "1"
// 		} else if f == "C" || f == "c" {
// 			f = "2"
// 		} else if f == "D" || f == "d" {
// 			f = "3"
// 		} else {
// 			return errors.New("錯誤: 讀取excel檔案發生問題，請重新操作")
// 		}

// 		qa[index] = a
// 		qa[index+1] = strings.Join([]string{b, c, d, e}, "&&&")
// 		qa[index+2] = f
// 		qa[index+3] = g

// 		// 下一題題目設置的index間隔為4
// 		index += 4
// 	}
// } else {
// total = values.Get("total_qa")
// qa = []string{
// 	values.Get("qa_1"), values.Get("qa_1_options"), values.Get("qa_1_answer"), values.Get("qa_1_score"),
// 	values.Get("qa_2"), values.Get("qa_2_options"), values.Get("qa_2_answer"), values.Get("qa_2_score"),
// 	values.Get("qa_3"), values.Get("qa_3_options"), values.Get("qa_3_answer"), values.Get("qa_3_score"),
// 	values.Get("qa_4"), values.Get("qa_4_options"), values.Get("qa_4_answer"), values.Get("qa_4_score"),
// 	values.Get("qa_5"), values.Get("qa_5_options"), values.Get("qa_5_answer"), values.Get("qa_5_score"),
// 	values.Get("qa_6"), values.Get("qa_6_options"), values.Get("qa_6_answer"), values.Get("qa_6_score"),
// 	values.Get("qa_7"), values.Get("qa_7_options"), values.Get("qa_7_answer"), values.Get("qa_7_score"),
// 	values.Get("qa_8"), values.Get("qa_8_options"), values.Get("qa_8_answer"), values.Get("qa_8_score"),
// 	values.Get("qa_9"), values.Get("qa_9_options"), values.Get("qa_9_answer"), values.Get("qa_9_score"),
// 	values.Get("qa_10"), values.Get("qa_10_options"), values.Get("qa_10_answer"), values.Get("qa_10_score"),
// 	values.Get("qa_11"), values.Get("qa_11_options"), values.Get("qa_11_answer"), values.Get("qa_11_score"),
// 	values.Get("qa_12"), values.Get("qa_12_options"), values.Get("qa_12_answer"), values.Get("qa_12_score"),
// 	values.Get("qa_13"), values.Get("qa_13_options"), values.Get("qa_13_answer"), values.Get("qa_13_score"),
// 	values.Get("qa_14"), values.Get("qa_14_options"), values.Get("qa_14_answer"), values.Get("qa_14_score"),
// 	values.Get("qa_15"), values.Get("qa_15_options"), values.Get("qa_15_answer"), values.Get("qa_15_score"),
// 	values.Get("qa_16"), values.Get("qa_16_options"), values.Get("qa_16_answer"), values.Get("qa_16_score"),
// 	values.Get("qa_17"), values.Get("qa_17_options"), values.Get("qa_17_answer"), values.Get("qa_17_score"),
// 	values.Get("qa_18"), values.Get("qa_18_options"), values.Get("qa_18_answer"), values.Get("qa_18_score"),
// 	values.Get("qa_19"), values.Get("qa_19_options"), values.Get("qa_19_answer"), values.Get("qa_19_score"),
// 	values.Get("qa_20"), values.Get("qa_20_options"), values.Get("qa_20_answer"), values.Get("qa_20_score"),
// }
// }
// if total == "" {
// 	return errors.New("錯誤: 題目最少一題，請重新操作")
// }

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

// 	// 自定義圖片
// 	QABgmStart:  update[0], // 遊戲開始
// 	QABgmGaming: update[1], // 遊戲進行中
// 	QABgmEnd:    update[2], // 遊戲結束

// 	QAClassicHPic01: update[3],
// 	QAClassicHPic02: update[4],
// 	QAClassicHPic03: update[5],
// 	QAClassicHPic04: update[6],
// 	QAClassicHPic05: update[7],
// 	QAClassicHPic06: update[8],
// 	QAClassicHPic07: update[9],
// 	QAClassicHPic08: update[10],
// 	QAClassicHPic09: update[11],
// 	QAClassicHPic10: update[12],
// 	QAClassicHPic11: update[13],
// 	QAClassicHPic12: update[14],
// 	QAClassicHPic13: update[15],
// 	QAClassicHPic14: update[16],
// 	QAClassicHPic15: update[17],
// 	QAClassicHPic16: update[18],
// 	QAClassicHPic17: update[19],
// 	QAClassicHPic18: update[20],
// 	QAClassicHPic19: update[21],
// 	QAClassicHPic20: update[22],
// 	QAClassicHPic21: update[23],
// 	QAClassicHPic22: update[24],
// 	QAClassicGPic01: update[25],
// 	QAClassicGPic02: update[26],
// 	QAClassicGPic03: update[27],
// 	QAClassicGPic04: update[28],
// 	QAClassicGPic05: update[29],
// 	QAClassicCPic01: update[30],
// 	QAClassicHAni01: update[31],
// 	QAClassicHAni02: update[32],
// 	QAClassicGAni01: update[33],
// 	QAClassicGAni02: update[34],

// 	QAElectricHPic01: update[35],
// 	QAElectricHPic02: update[36],
// 	QAElectricHPic03: update[37],
// 	QAElectricHPic04: update[38],
// 	QAElectricHPic05: update[39],
// 	QAElectricHPic06: update[40],
// 	QAElectricHPic07: update[41],
// 	QAElectricHPic08: update[42],
// 	QAElectricHPic09: update[43],
// 	QAElectricHPic10: update[44],
// 	QAElectricHPic11: update[45],
// 	QAElectricHPic12: update[46],
// 	QAElectricHPic13: update[47],
// 	QAElectricHPic14: update[48],
// 	QAElectricHPic15: update[49],
// 	QAElectricHPic16: update[50],
// 	QAElectricHPic17: update[51],
// 	QAElectricHPic18: update[52],
// 	QAElectricHPic19: update[53],
// 	QAElectricHPic20: update[54],
// 	QAElectricHPic21: update[55],
// 	QAElectricHPic22: update[56],
// 	QAElectricHPic23: update[57],
// 	QAElectricHPic24: update[58],
// 	QAElectricHPic25: update[59],
// 	QAElectricHPic26: update[60],
// 	QAElectricGPic01: update[61],
// 	QAElectricGPic02: update[62],
// 	QAElectricGPic03: update[63],
// 	QAElectricGPic04: update[64],
// 	QAElectricGPic05: update[65],
// 	QAElectricGPic06: update[66],
// 	QAElectricGPic07: update[67],
// 	QAElectricGPic08: update[68],
// 	QAElectricGPic09: update[69],
// 	QAElectricCPic01: update[70],
// 	QAElectricHAni01: update[71],
// 	QAElectricHAni02: update[72],
// 	QAElectricHAni03: update[73],
// 	QAElectricHAni04: update[74],
// 	QAElectricHAni05: update[75],
// 	QAElectricGAni01: update[76],
// 	QAElectricGAni02: update[77],
// 	QAElectricCAni01: update[78],

// 	QAMoonfestivalHPic01: update[79],
// 	QAMoonfestivalHPic02: update[80],
// 	QAMoonfestivalHPic03: update[81],
// 	QAMoonfestivalHPic04: update[82],
// 	QAMoonfestivalHPic05: update[83],
// 	QAMoonfestivalHPic06: update[84],
// 	QAMoonfestivalHPic07: update[85],
// 	QAMoonfestivalHPic08: update[86],
// 	QAMoonfestivalHPic09: update[87],
// 	QAMoonfestivalHPic10: update[88],
// 	QAMoonfestivalHPic11: update[89],
// 	QAMoonfestivalHPic12: update[90],
// 	QAMoonfestivalHPic13: update[91],
// 	QAMoonfestivalHPic14: update[92],
// 	QAMoonfestivalHPic15: update[93],
// 	QAMoonfestivalHPic16: update[94],
// 	QAMoonfestivalHPic17: update[95],
// 	QAMoonfestivalHPic18: update[96],
// 	QAMoonfestivalHPic19: update[97],
// 	QAMoonfestivalHPic20: update[98],
// 	QAMoonfestivalHPic21: update[99],
// 	QAMoonfestivalHPic22: update[100],
// 	QAMoonfestivalHPic23: update[101],
// 	QAMoonfestivalHPic24: update[102],
// 	QAMoonfestivalGPic01: update[103],
// 	QAMoonfestivalGPic02: update[104],
// 	QAMoonfestivalGPic03: update[105],
// 	QAMoonfestivalGPic04: update[106],
// 	QAMoonfestivalGPic05: update[107],
// 	QAMoonfestivalCPic01: update[108],
// 	QAMoonfestivalCPic02: update[109],
// 	QAMoonfestivalCPic03: update[110],
// 	QAMoonfestivalHAni01: update[111],
// 	QAMoonfestivalHAni02: update[112],
// 	QAMoonfestivalGAni01: update[113],
// 	QAMoonfestivalGAni02: update[114],
// 	QAMoonfestivalGAni03: update[115],

// 	QANewyearDragonHPic01: update[116],
// 	QANewyearDragonHPic02: update[117],
// 	QANewyearDragonHPic03: update[118],
// 	QANewyearDragonHPic04: update[119],
// 	QANewyearDragonHPic05: update[120],
// 	QANewyearDragonHPic06: update[121],
// 	QANewyearDragonHPic07: update[122],
// 	QANewyearDragonHPic08: update[123],
// 	QANewyearDragonHPic09: update[124],
// 	QANewyearDragonHPic10: update[125],
// 	QANewyearDragonHPic11: update[126],
// 	QANewyearDragonHPic12: update[127],
// 	QANewyearDragonHPic13: update[128],
// 	QANewyearDragonHPic14: update[129],
// 	QANewyearDragonHPic15: update[130],
// 	QANewyearDragonHPic16: update[131],
// 	QANewyearDragonHPic17: update[132],
// 	QANewyearDragonHPic18: update[133],
// 	QANewyearDragonHPic19: update[134],
// 	QANewyearDragonHPic20: update[135],
// 	QANewyearDragonHPic21: update[136],
// 	QANewyearDragonHPic22: update[137],
// 	QANewyearDragonHPic23: update[138],
// 	QANewyearDragonHPic24: update[139],
// 	QANewyearDragonGPic01: update[140],
// 	QANewyearDragonGPic02: update[141],
// 	QANewyearDragonGPic03: update[142],
// 	QANewyearDragonGPic04: update[143],
// 	QANewyearDragonGPic05: update[144],
// 	QANewyearDragonGPic06: update[145],
// 	QANewyearDragonCPic01: update[146],
// 	QANewyearDragonHAni01: update[147],
// 	QANewyearDragonHAni02: update[148],
// 	QANewyearDragonGAni01: update[149],
// 	QANewyearDragonGAni02: update[150],
// 	QANewyearDragonGAni03: update[151],
// 	QANewyearDragonCAni01: update[152],

// 	QA1: qa[0], QA1Options: qa[1], QA1Answer: qa[2], QA1Score: qa[3],
// 	QA2: qa[4], QA2Options: qa[5], QA2Answer: qa[6], QA2Score: qa[7],
// 	QA3: qa[8], QA3Options: qa[9], QA3Answer: qa[10], QA3Score: qa[11],
// 	QA4: qa[12], QA4Options: qa[13], QA4Answer: qa[14], QA4Score: qa[15],
// 	QA5: qa[16], QA5Options: qa[17], QA5Answer: qa[18], QA5Score: qa[19],
// 	QA6: qa[20], QA6Options: qa[21], QA6Answer: qa[22], QA6Score: qa[23],
// 	QA7: qa[24], QA7Options: qa[25], QA7Answer: qa[26], QA7Score: qa[27],
// 	QA8: qa[28], QA8Options: qa[29], QA8Answer: qa[30], QA8Score: qa[31],
// 	QA9: qa[32], QA9Options: qa[33], QA9Answer: qa[34], QA9Score: qa[35],
// 	QA10: qa[36], QA10Options: qa[37], QA10Answer: qa[38], QA10Score: qa[39],
// 	QA11: qa[40], QA11Options: qa[41], QA11Answer: qa[42], QA11Score: qa[43],
// 	QA12: qa[44], QA12Options: qa[45], QA12Answer: qa[46], QA12Score: qa[47],
// 	QA13: qa[48], QA13Options: qa[49], QA13Answer: qa[50], QA13Score: qa[51],
// 	QA14: qa[52], QA14Options: qa[53], QA14Answer: qa[54], QA14Score: qa[55],
// 	QA15: qa[56], QA15Options: qa[57], QA15Answer: qa[58], QA15Score: qa[59],
// 	QA16: qa[60], QA16Options: qa[61], QA16Answer: qa[62], QA16Score: qa[63],
// 	QA17: qa[64], QA17Options: qa[65], QA17Answer: qa[66], QA17Score: qa[67],
// 	QA18: qa[68], QA18Options: qa[69], QA18Answer: qa[70], QA18Score: qa[71],
// 	QA19: qa[72], QA19Options: qa[73], QA19Answer: qa[74], QA19Score: qa[75],
// 	QA20: qa[76], QA20Options: qa[77], QA20Answer: qa[78], QA20Score: qa[79],
// 	TotalQA:  values.Get("total_qa"),
// 	QASecond: values.Get("qa_second"),
// }
