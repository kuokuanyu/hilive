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
	ropepackPictureFields = []PictureField{
		{FieldName: "ropepack_bgm_start", Path: "ropepack/%s/bgm/start.mp3"},
		{FieldName: "ropepack_bgm_gaming", Path: "ropepack/%s/bgm/gaming.mp3"},
		{FieldName: "ropepack_bgm_end", Path: "ropepack/%s/bgm/end.mp3"},

		{FieldName: "ropepack_classic_h_pic_01", Path: "ropepack/classic/ropepack_classic_h_pic_01.png"},
		{FieldName: "ropepack_classic_h_pic_02", Path: "ropepack/classic/ropepack_classic_h_pic_02.png"},
		{FieldName: "ropepack_classic_h_pic_03", Path: "ropepack/classic/ropepack_classic_h_pic_03.jpg"},
		{FieldName: "ropepack_classic_h_pic_04", Path: "ropepack/classic/ropepack_classic_h_pic_04.png"},
		{FieldName: "ropepack_classic_h_pic_05", Path: "ropepack/classic/ropepack_classic_h_pic_05.png"},
		{FieldName: "ropepack_classic_h_pic_06", Path: "ropepack/classic/ropepack_classic_h_pic_06.png"},
		{FieldName: "ropepack_classic_h_pic_07", Path: "ropepack/classic/ropepack_classic_h_pic_07.png"},
		{FieldName: "ropepack_classic_h_pic_08", Path: "ropepack/classic/ropepack_classic_h_pic_08.png"},
		{FieldName: "ropepack_classic_h_pic_09", Path: "ropepack/classic/ropepack_classic_h_pic_09.png"},
		{FieldName: "ropepack_classic_h_pic_10", Path: "ropepack/classic/ropepack_classic_h_pic_10.png"},
		{FieldName: "ropepack_classic_g_pic_01", Path: "ropepack/classic/ropepack_classic_g_pic_01.png"},
		{FieldName: "ropepack_classic_g_pic_02", Path: "ropepack/classic/ropepack_classic_g_pic_02.png"},
		{FieldName: "ropepack_classic_g_pic_03", Path: "ropepack/classic/ropepack_classic_g_pic_03.jpg"},
		{FieldName: "ropepack_classic_g_pic_04", Path: "ropepack/classic/ropepack_classic_g_pic_04.png"},
		{FieldName: "ropepack_classic_g_pic_05", Path: "ropepack/classic/ropepack_classic_g_pic_05.png"},
		{FieldName: "ropepack_classic_g_pic_06", Path: "ropepack/classic/ropepack_classic_g_pic_06.jpg"},
		{FieldName: "ropepack_classic_h_ani_01", Path: "ropepack/classic/ropepack_classic_h_ani_01.png"},
		{FieldName: "ropepack_classic_g_ani_01", Path: "ropepack/classic/ropepack_classic_g_ani_01.png"},
		{FieldName: "ropepack_classic_g_ani_02", Path: "ropepack/classic/ropepack_classic_g_ani_02.png"},
		{FieldName: "ropepack_classic_c_ani_01", Path: "ropepack/classic/ropepack_classic_c_ani_01.png"},

		{FieldName: "ropepack_newyear_rabbit_h_pic_01", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_01.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_02", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_02.jpg"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_03", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_03.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_04", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_04.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_05", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_05.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_06", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_06.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_07", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_07.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_08", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_08.png"},
		{FieldName: "ropepack_newyear_rabbit_h_pic_09", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_09.png"},
		{FieldName: "ropepack_newyear_rabbit_g_pic_01", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_01.png"},
		{FieldName: "ropepack_newyear_rabbit_g_pic_02", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_02.jpg"},
		{FieldName: "ropepack_newyear_rabbit_g_pic_03", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_03.png"},
		{FieldName: "ropepack_newyear_rabbit_h_ani_01", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_ani_01.png"},
		{FieldName: "ropepack_newyear_rabbit_g_ani_01", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_01.png"},
		{FieldName: "ropepack_newyear_rabbit_g_ani_02", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_02.png"},
		{FieldName: "ropepack_newyear_rabbit_g_ani_03", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_03.png"},
		{FieldName: "ropepack_newyear_rabbit_c_ani_01", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_01.png"},
		{FieldName: "ropepack_newyear_rabbit_c_ani_02", Path: "ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_02.png"},

		{FieldName: "ropepack_moonfestival_h_pic_01", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_01.png"},
		{FieldName: "ropepack_moonfestival_h_pic_02", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_02.png"},
		{FieldName: "ropepack_moonfestival_h_pic_03", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_03.png"},
		{FieldName: "ropepack_moonfestival_h_pic_04", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_04.png"},
		{FieldName: "ropepack_moonfestival_h_pic_05", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_05.jpg"},
		{FieldName: "ropepack_moonfestival_h_pic_06", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_06.png"},
		{FieldName: "ropepack_moonfestival_h_pic_07", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_07.png"},
		{FieldName: "ropepack_moonfestival_h_pic_08", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_08.jpg"},
		{FieldName: "ropepack_moonfestival_h_pic_09", Path: "ropepack/moonfestival/ropepack_moonfestival_h_pic_09.png"},
		{FieldName: "ropepack_moonfestival_g_pic_01", Path: "ropepack/moonfestival/ropepack_moonfestival_g_pic_01.png"},
		{FieldName: "ropepack_moonfestival_g_pic_02", Path: "ropepack/moonfestival/ropepack_moonfestival_g_pic_02.jpg"},
		{FieldName: "ropepack_moonfestival_c_pic_01", Path: "ropepack/moonfestival/ropepack_moonfestival_c_pic_01.png"},
		{FieldName: "ropepack_moonfestival_h_ani_01", Path: "ropepack/moonfestival/ropepack_moonfestival_h_ani_01.png"},
		{FieldName: "ropepack_moonfestival_g_ani_01", Path: "ropepack/moonfestival/ropepack_moonfestival_g_ani_01.png"},
		{FieldName: "ropepack_moonfestival_g_ani_02", Path: "ropepack/moonfestival/ropepack_moonfestival_g_ani_02.png"},
		{FieldName: "ropepack_moonfestival_c_ani_01", Path: "ropepack/moonfestival/ropepack_moonfestival_c_ani_01.png"},
		{FieldName: "ropepack_moonfestival_c_ani_02", Path: "ropepack/moonfestival/ropepack_moonfestival_c_ani_02.png"},

		{FieldName: "ropepack_3D_h_pic_01", Path: "ropepack/3D/ropepack_3D_h_pic_01.png"},
		{FieldName: "ropepack_3D_h_pic_02", Path: "ropepack/3D/ropepack_3D_h_pic_02.png"},
		{FieldName: "ropepack_3D_h_pic_03", Path: "ropepack/3D/ropepack_3D_h_pic_03.png"},
		{FieldName: "ropepack_3D_h_pic_04", Path: "ropepack/3D/ropepack_3D_h_pic_04.jpg"},
		{FieldName: "ropepack_3D_h_pic_05", Path: "ropepack/3D/ropepack_3D_h_pic_05.png"},
		{FieldName: "ropepack_3D_h_pic_06", Path: "ropepack/3D/ropepack_3D_h_pic_06.png"},
		{FieldName: "ropepack_3D_h_pic_07", Path: "ropepack/3D/ropepack_3D_h_pic_07.png"},
		{FieldName: "ropepack_3D_h_pic_08", Path: "ropepack/3D/ropepack_3D_h_pic_08.png"},
		{FieldName: "ropepack_3D_h_pic_09", Path: "ropepack/3D/ropepack_3D_h_pic_09.png"},
		{FieldName: "ropepack_3D_h_pic_10", Path: "ropepack/3D/ropepack_3D_h_pic_10.png"},
		{FieldName: "ropepack_3D_h_pic_11", Path: "ropepack/3D/ropepack_3D_h_pic_11.png"},
		{FieldName: "ropepack_3D_h_pic_12", Path: "ropepack/3D/ropepack_3D_h_pic_12.png"},
		{FieldName: "ropepack_3D_h_pic_13", Path: "ropepack/3D/ropepack_3D_h_pic_13.png"},
		{FieldName: "ropepack_3D_h_pic_14", Path: "ropepack/3D/ropepack_3D_h_pic_14.png"},
		{FieldName: "ropepack_3D_h_pic_15", Path: "ropepack/3D/ropepack_3D_h_pic_15.png"},
		{FieldName: "ropepack_3D_g_pic_01", Path: "ropepack/3D/ropepack_3D_g_pic_01.png"},
		{FieldName: "ropepack_3D_g_pic_02", Path: "ropepack/3D/ropepack_3D_g_pic_02.jpg"},
		{FieldName: "ropepack_3D_g_pic_03", Path: "ropepack/3D/ropepack_3D_g_pic_03.png"},
		{FieldName: "ropepack_3D_g_pic_04", Path: "ropepack/3D/ropepack_3D_g_pic_04.png"},
		{FieldName: "ropepack_3D_h_ani_01", Path: "ropepack/3D/ropepack_3D_h_ani_01.png"},
		{FieldName: "ropepack_3D_h_ani_02", Path: "ropepack/3D/ropepack_3D_h_ani_02.png"},
		{FieldName: "ropepack_3D_h_ani_03", Path: "ropepack/3D/ropepack_3D_h_ani_03.png"},
		{FieldName: "ropepack_3D_g_ani_01", Path: "ropepack/3D/ropepack_3D_g_ani_01.png"},
		{FieldName: "ropepack_3D_g_ani_02", Path: "ropepack/3D/ropepack_3D_g_ani_02.png"},
		{FieldName: "ropepack_3D_c_ani_01", Path: "ropepack/3D/ropepack_3D_c_ani_01.png"},
	}
)

// GetRopepackPanel 套紅包
func (s *SystemTable) GetRopepackPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/ropepack/" + id)
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
		picMap := BuildPictureMap(ropepackPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "ropepack", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(ropepackPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "ropepack", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增套紅包遊戲資料(form-data)
// @Tags Ropepack
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
// @param topic formData string true "主題樣式" Enums(01_classic, 02_newyear_rabbit)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack/form [post]
func POSTRopepack(ctx *gin.Context) {
}

// @Summary 新增套紅包獎品資料(form-data)
// @Tags Ropepack Prize
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
// @Router /interact/game/ropepack/prize/form [post]
func POSTRopepackPrize(ctx *gin.Context) {
}

// @Summary 編輯套紅包遊戲資料(form-data)
// @Tags Ropepack
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
// @param topic formData string false "主題樣式" Enums(01_classic, 02_newyear_rabbit)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack/form [put]
func PUTRopepack(ctx *gin.Context) {
}

// @Summary 編輯套紅包獎品資料(form-data)
// @Tags Ropepack Prize
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
// @Router /interact/game/ropepack/prize/form [put]
func PUTRopepackPrize(ctx *gin.Context) {
}

// @Summary 刪除套紅包遊戲資料(form-data)
// @Tags Ropepack
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack/form [delete]
func DELETERopepack(ctx *gin.Context) {
}

// @Summary 刪除套紅包獎品資料(form-data)
// @Tags Ropepack Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack/prize/form [delete]
func DELETERopepackPrize(ctx *gin.Context) {
}

// @Summary 套紅包遊戲JSON資料
// @Tags Ropepack
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack [get]
func RopepackJSON(ctx *gin.Context) {
}

// @Summary 套紅包獎品JSON資料
// @Tags Ropepack Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/ropepack/prize [get]
func RopepackPrizeJSON(ctx *gin.Context) {
}

// pics = []string{
// 套紅包自定義
// 音樂
// "ropepack/%s/bgm/start.mp3",
// "ropepack/%s/bgm/gaming.mp3",
// "ropepack/%s/bgm/end.mp3",

// "ropepack/classic/ropepack_classic_h_pic_01.png",
// "ropepack/classic/ropepack_classic_h_pic_02.png",
// "ropepack/classic/ropepack_classic_h_pic_03.jpg",
// "ropepack/classic/ropepack_classic_h_pic_04.png",
// "ropepack/classic/ropepack_classic_h_pic_05.png",
// "ropepack/classic/ropepack_classic_h_pic_06.png",
// "ropepack/classic/ropepack_classic_h_pic_07.png",
// "ropepack/classic/ropepack_classic_h_pic_08.png",
// "ropepack/classic/ropepack_classic_h_pic_09.png",
// "ropepack/classic/ropepack_classic_h_pic_10.png",
// "ropepack/classic/ropepack_classic_g_pic_01.png",
// "ropepack/classic/ropepack_classic_g_pic_02.png",
// "ropepack/classic/ropepack_classic_g_pic_03.jpg",
// "ropepack/classic/ropepack_classic_g_pic_04.png",
// "ropepack/classic/ropepack_classic_g_pic_05.png",
// "ropepack/classic/ropepack_classic_g_pic_06.jpg",
// "ropepack/classic/ropepack_classic_h_ani_01.png",
// "ropepack/classic/ropepack_classic_g_ani_01.png",
// "ropepack/classic/ropepack_classic_g_ani_02.png",
// "ropepack/classic/ropepack_classic_c_ani_01.png",

// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_01.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_02.jpg",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_03.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_04.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_05.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_06.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_07.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_08.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_09.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_01.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_02.jpg",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_03.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_ani_01.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_01.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_02.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_03.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_01.png",
// "ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_02.png",

// "ropepack/moonfestival/ropepack_moonfestival_h_pic_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_02.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_03.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_04.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_05.jpg",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_06.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_07.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_08.jpg",
// "ropepack/moonfestival/ropepack_moonfestival_h_pic_09.png",
// "ropepack/moonfestival/ropepack_moonfestival_g_pic_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_g_pic_02.jpg",
// "ropepack/moonfestival/ropepack_moonfestival_c_pic_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_h_ani_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_g_ani_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_g_ani_02.png",
// "ropepack/moonfestival/ropepack_moonfestival_c_ani_01.png",
// "ropepack/moonfestival/ropepack_moonfestival_c_ani_02.png",

// "ropepack/3D/ropepack_3D_h_pic_01.png",
// "ropepack/3D/ropepack_3D_h_pic_02.png",
// "ropepack/3D/ropepack_3D_h_pic_03.png",
// "ropepack/3D/ropepack_3D_h_pic_04.jpg",
// "ropepack/3D/ropepack_3D_h_pic_05.png",
// "ropepack/3D/ropepack_3D_h_pic_06.png",
// "ropepack/3D/ropepack_3D_h_pic_07.png",
// "ropepack/3D/ropepack_3D_h_pic_08.png",
// "ropepack/3D/ropepack_3D_h_pic_09.png",
// "ropepack/3D/ropepack_3D_h_pic_10.png",
// "ropepack/3D/ropepack_3D_h_pic_11.png",
// "ropepack/3D/ropepack_3D_h_pic_12.png",
// "ropepack/3D/ropepack_3D_h_pic_13.png",
// "ropepack/3D/ropepack_3D_h_pic_14.png",
// "ropepack/3D/ropepack_3D_h_pic_15.png",
// "ropepack/3D/ropepack_3D_g_pic_01.png",
// "ropepack/3D/ropepack_3D_g_pic_02.jpg",
// "ropepack/3D/ropepack_3D_g_pic_03.png",
// "ropepack/3D/ropepack_3D_g_pic_04.png",
// "ropepack/3D/ropepack_3D_h_ani_01.png",
// "ropepack/3D/ropepack_3D_h_ani_02.png",
// "ropepack/3D/ropepack_3D_h_ani_03.png",
// "ropepack/3D/ropepack_3D_g_ani_01.png",
// "ropepack/3D/ropepack_3D_g_ani_02.png",
// "ropepack/3D/ropepack_3D_c_ani_01.png",
// }

// fields = []string{
// 套紅包自定義
// 音樂
// "ropepack_bgm_start",
// "ropepack_bgm_gaming",
// "ropepack_bgm_end",

// "ropepack_classic_h_pic_01",
// "ropepack_classic_h_pic_02",
// "ropepack_classic_h_pic_03",
// "ropepack_classic_h_pic_04",
// "ropepack_classic_h_pic_05",
// "ropepack_classic_h_pic_06",
// "ropepack_classic_h_pic_07",
// "ropepack_classic_h_pic_08",
// "ropepack_classic_h_pic_09",
// "ropepack_classic_h_pic_10",
// "ropepack_classic_g_pic_01",
// "ropepack_classic_g_pic_02",
// "ropepack_classic_g_pic_03",
// "ropepack_classic_g_pic_04",
// "ropepack_classic_g_pic_05",
// "ropepack_classic_g_pic_06",
// "ropepack_classic_h_ani_01",
// "ropepack_classic_g_ani_01",
// "ropepack_classic_g_ani_02",
// "ropepack_classic_c_ani_01",

// "ropepack_newyear_rabbit_h_pic_01",
// "ropepack_newyear_rabbit_h_pic_02",
// "ropepack_newyear_rabbit_h_pic_03",
// "ropepack_newyear_rabbit_h_pic_04",
// "ropepack_newyear_rabbit_h_pic_05",
// "ropepack_newyear_rabbit_h_pic_06",
// "ropepack_newyear_rabbit_h_pic_07",
// "ropepack_newyear_rabbit_h_pic_08",
// "ropepack_newyear_rabbit_h_pic_09",
// "ropepack_newyear_rabbit_g_pic_01",
// "ropepack_newyear_rabbit_g_pic_02",
// "ropepack_newyear_rabbit_g_pic_03",
// "ropepack_newyear_rabbit_h_ani_01",
// "ropepack_newyear_rabbit_g_ani_01",
// "ropepack_newyear_rabbit_g_ani_02",
// "ropepack_newyear_rabbit_g_ani_03",
// "ropepack_newyear_rabbit_c_ani_01",
// "ropepack_newyear_rabbit_c_ani_02",

// "ropepack_moonfestival_h_pic_01",
// "ropepack_moonfestival_h_pic_02",
// "ropepack_moonfestival_h_pic_03",
// "ropepack_moonfestival_h_pic_04",
// "ropepack_moonfestival_h_pic_05",
// "ropepack_moonfestival_h_pic_06",
// "ropepack_moonfestival_h_pic_07",
// "ropepack_moonfestival_h_pic_08",
// "ropepack_moonfestival_h_pic_09",
// "ropepack_moonfestival_g_pic_01",
// "ropepack_moonfestival_g_pic_02",
// "ropepack_moonfestival_c_pic_01",
// "ropepack_moonfestival_h_ani_01",
// "ropepack_moonfestival_g_ani_01",
// "ropepack_moonfestival_g_ani_02",
// "ropepack_moonfestival_c_ani_01",
// "ropepack_moonfestival_c_ani_02",

// "ropepack_3D_h_pic_01",
// "ropepack_3D_h_pic_02",
// "ropepack_3D_h_pic_03",
// "ropepack_3D_h_pic_04",
// "ropepack_3D_h_pic_05",
// "ropepack_3D_h_pic_06",
// "ropepack_3D_h_pic_07",
// "ropepack_3D_h_pic_08",
// "ropepack_3D_h_pic_09",
// "ropepack_3D_h_pic_10",
// "ropepack_3D_h_pic_11",
// "ropepack_3D_h_pic_12",
// "ropepack_3D_h_pic_13",
// "ropepack_3D_h_pic_14",
// "ropepack_3D_h_pic_15",
// "ropepack_3D_g_pic_01",
// "ropepack_3D_g_pic_02",
// "ropepack_3D_g_pic_03",
// "ropepack_3D_g_pic_04",
// "ropepack_3D_h_ani_01",
// "ropepack_3D_h_ani_02",
// "ropepack_3D_h_ani_03",
// "ropepack_3D_g_ani_01",
// "ropepack_3D_g_ani_02",
// "ropepack_3D_c_ani_01",
// }
// update = make([]string, 300)

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

// 	// 套紅包自定義
// 	// 音樂
// 	RopepackBgmStart:  update[0],
// 	RopepackBgmGaming: update[1],
// 	RopepackBgmEnd:    update[2],

// 	RopepackClassicHPic01: update[3],
// 	RopepackClassicHPic02: update[4],
// 	RopepackClassicHPic03: update[5],
// 	RopepackClassicHPic04: update[6],
// 	RopepackClassicHPic05: update[7],
// 	RopepackClassicHPic06: update[8],
// 	RopepackClassicHPic07: update[9],
// 	RopepackClassicHPic08: update[10],
// 	RopepackClassicHPic09: update[11],
// 	RopepackClassicHPic10: update[12],
// 	RopepackClassicGPic01: update[13],
// 	RopepackClassicGPic02: update[14],
// 	RopepackClassicGPic03: update[15],
// 	RopepackClassicGPic04: update[16],
// 	RopepackClassicGPic05: update[17],
// 	RopepackClassicGPic06: update[18],
// 	RopepackClassicHAni01: update[19],
// 	RopepackClassicGAni01: update[20],
// 	RopepackClassicGAni02: update[21],
// 	RopepackClassicCAni01: update[22],

// 	RopepackNewyearRabbitHPic01: update[23],
// 	RopepackNewyearRabbitHPic02: update[24],
// 	RopepackNewyearRabbitHPic03: update[25],
// 	RopepackNewyearRabbitHPic04: update[26],
// 	RopepackNewyearRabbitHPic05: update[27],
// 	RopepackNewyearRabbitHPic06: update[28],
// 	RopepackNewyearRabbitHPic07: update[29],
// 	RopepackNewyearRabbitHPic08: update[30],
// 	RopepackNewyearRabbitHPic09: update[31],
// 	RopepackNewyearRabbitGPic01: update[32],
// 	RopepackNewyearRabbitGPic02: update[33],
// 	RopepackNewyearRabbitGPic03: update[34],
// 	RopepackNewyearRabbitHAni01: update[35],
// 	RopepackNewyearRabbitGAni01: update[36],
// 	RopepackNewyearRabbitGAni02: update[37],
// 	RopepackNewyearRabbitGAni03: update[38],
// 	RopepackNewyearRabbitCAni01: update[39],
// 	RopepackNewyearRabbitCAni02: update[40],

// 	RopepackMoonfestivalHPic01: update[41],
// 	RopepackMoonfestivalHPic02: update[42],
// 	RopepackMoonfestivalHPic03: update[43],
// 	RopepackMoonfestivalHPic04: update[44],
// 	RopepackMoonfestivalHPic05: update[45],
// 	RopepackMoonfestivalHPic06: update[46],
// 	RopepackMoonfestivalHPic07: update[47],
// 	RopepackMoonfestivalHPic08: update[48],
// 	RopepackMoonfestivalHPic09: update[49],
// 	RopepackMoonfestivalGPic01: update[50],
// 	RopepackMoonfestivalGPic02: update[51],
// 	RopepackMoonfestivalCPic01: update[52],
// 	RopepackMoonfestivalHAni01: update[53],
// 	RopepackMoonfestivalGAni01: update[54],
// 	RopepackMoonfestivalGAni02: update[55],
// 	RopepackMoonfestivalCAni01: update[56],
// 	RopepackMoonfestivalCAni02: update[57],

// 	Ropepack3DHPic01: update[58],
// 	Ropepack3DHPic02: update[59],
// 	Ropepack3DHPic03: update[60],
// 	Ropepack3DHPic04: update[61],
// 	Ropepack3DHPic05: update[62],
// 	Ropepack3DHPic06: update[63],
// 	Ropepack3DHPic07: update[64],
// 	Ropepack3DHPic08: update[65],
// 	Ropepack3DHPic09: update[66],
// 	Ropepack3DHPic10: update[67],
// 	Ropepack3DHPic11: update[68],
// 	Ropepack3DHPic12: update[69],
// 	Ropepack3DHPic13: update[70],
// 	Ropepack3DHPic14: update[71],
// 	Ropepack3DHPic15: update[72],
// 	Ropepack3DGPic01: update[73],
// 	Ropepack3DGPic02: update[74],
// 	Ropepack3DGPic03: update[75],
// 	Ropepack3DGPic04: update[76],
// 	Ropepack3DHAni01: update[77],
// 	Ropepack3DHAni02: update[78],
// 	Ropepack3DHAni03: update[79],
// 	Ropepack3DGAni01: update[80],
// 	Ropepack3DGAni02: update[81],
// 	Ropepack3DCAni01: update[82],
// }

// var (
// 	pics = []string{
// 		// 套紅包自定義
// 		// 音樂
// 		"ropepack/%s/bgm/start.mp3",
// 		"ropepack/%s/bgm/gaming.mp3",
// 		"ropepack/%s/bgm/end.mp3",

// 		"ropepack/classic/ropepack_classic_h_pic_01.png",
// 		"ropepack/classic/ropepack_classic_h_pic_02.png",
// 		"ropepack/classic/ropepack_classic_h_pic_03.jpg",
// 		"ropepack/classic/ropepack_classic_h_pic_04.png",
// 		"ropepack/classic/ropepack_classic_h_pic_05.png",
// 		"ropepack/classic/ropepack_classic_h_pic_06.png",
// 		"ropepack/classic/ropepack_classic_h_pic_07.png",
// 		"ropepack/classic/ropepack_classic_h_pic_08.png",
// 		"ropepack/classic/ropepack_classic_h_pic_09.png",
// 		"ropepack/classic/ropepack_classic_h_pic_10.png",
// 		"ropepack/classic/ropepack_classic_g_pic_01.png",
// 		"ropepack/classic/ropepack_classic_g_pic_02.png",
// 		"ropepack/classic/ropepack_classic_g_pic_03.jpg",
// 		"ropepack/classic/ropepack_classic_g_pic_04.png",
// 		"ropepack/classic/ropepack_classic_g_pic_05.png",
// 		"ropepack/classic/ropepack_classic_g_pic_06.jpg",
// 		"ropepack/classic/ropepack_classic_h_ani_01.png",
// 		"ropepack/classic/ropepack_classic_g_ani_01.png",
// 		"ropepack/classic/ropepack_classic_g_ani_02.png",
// 		"ropepack/classic/ropepack_classic_c_ani_01.png",

// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_01.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_02.jpg",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_03.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_04.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_05.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_06.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_07.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_08.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_pic_09.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_01.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_02.jpg",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_pic_03.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_h_ani_01.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_01.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_02.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_g_ani_03.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_01.png",
// 		"ropepack/newyear_rabbit/ropepack_newyear_rabbit_c_ani_02.png",

// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_02.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_03.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_04.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_05.jpg",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_06.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_07.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_08.jpg",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_pic_09.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_g_pic_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_g_pic_02.jpg",
// 		"ropepack/moonfestival/ropepack_moonfestival_c_pic_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_h_ani_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_g_ani_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_g_ani_02.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_c_ani_01.png",
// 		"ropepack/moonfestival/ropepack_moonfestival_c_ani_02.png",

// 		"ropepack/3D/ropepack_3D_h_pic_01.png",
// 		"ropepack/3D/ropepack_3D_h_pic_02.png",
// 		"ropepack/3D/ropepack_3D_h_pic_03.png",
// 		"ropepack/3D/ropepack_3D_h_pic_04.jpg",
// 		"ropepack/3D/ropepack_3D_h_pic_05.png",
// 		"ropepack/3D/ropepack_3D_h_pic_06.png",
// 		"ropepack/3D/ropepack_3D_h_pic_07.png",
// 		"ropepack/3D/ropepack_3D_h_pic_08.png",
// 		"ropepack/3D/ropepack_3D_h_pic_09.png",
// 		"ropepack/3D/ropepack_3D_h_pic_10.png",
// 		"ropepack/3D/ropepack_3D_h_pic_11.png",
// 		"ropepack/3D/ropepack_3D_h_pic_12.png",
// 		"ropepack/3D/ropepack_3D_h_pic_13.png",
// 		"ropepack/3D/ropepack_3D_h_pic_14.png",
// 		"ropepack/3D/ropepack_3D_h_pic_15.png",
// 		"ropepack/3D/ropepack_3D_g_pic_01.png",
// 		"ropepack/3D/ropepack_3D_g_pic_02.jpg",
// 		"ropepack/3D/ropepack_3D_g_pic_03.png",
// 		"ropepack/3D/ropepack_3D_g_pic_04.png",
// 		"ropepack/3D/ropepack_3D_h_ani_01.png",
// 		"ropepack/3D/ropepack_3D_h_ani_02.png",
// 		"ropepack/3D/ropepack_3D_h_ani_03.png",
// 		"ropepack/3D/ropepack_3D_g_ani_01.png",
// 		"ropepack/3D/ropepack_3D_g_ani_02.png",
// 		"ropepack/3D/ropepack_3D_c_ani_01.png",
// 	}

// 	fields = []string{
// 		// 套紅包自定義
// 		// 音樂
// 		"ropepack_bgm_start",
// 		"ropepack_bgm_gaming",
// 		"ropepack_bgm_end",

// 		"ropepack_classic_h_pic_01",
// 		"ropepack_classic_h_pic_02",
// 		"ropepack_classic_h_pic_03",
// 		"ropepack_classic_h_pic_04",
// 		"ropepack_classic_h_pic_05",
// 		"ropepack_classic_h_pic_06",
// 		"ropepack_classic_h_pic_07",
// 		"ropepack_classic_h_pic_08",
// 		"ropepack_classic_h_pic_09",
// 		"ropepack_classic_h_pic_10",
// 		"ropepack_classic_g_pic_01",
// 		"ropepack_classic_g_pic_02",
// 		"ropepack_classic_g_pic_03",
// 		"ropepack_classic_g_pic_04",
// 		"ropepack_classic_g_pic_05",
// 		"ropepack_classic_g_pic_06",
// 		"ropepack_classic_h_ani_01",
// 		"ropepack_classic_g_ani_01",
// 		"ropepack_classic_g_ani_02",
// 		"ropepack_classic_c_ani_01",

// 		"ropepack_newyear_rabbit_h_pic_01",
// 		"ropepack_newyear_rabbit_h_pic_02",
// 		"ropepack_newyear_rabbit_h_pic_03",
// 		"ropepack_newyear_rabbit_h_pic_04",
// 		"ropepack_newyear_rabbit_h_pic_05",
// 		"ropepack_newyear_rabbit_h_pic_06",
// 		"ropepack_newyear_rabbit_h_pic_07",
// 		"ropepack_newyear_rabbit_h_pic_08",
// 		"ropepack_newyear_rabbit_h_pic_09",
// 		"ropepack_newyear_rabbit_g_pic_01",
// 		"ropepack_newyear_rabbit_g_pic_02",
// 		"ropepack_newyear_rabbit_g_pic_03",
// 		"ropepack_newyear_rabbit_h_ani_01",
// 		"ropepack_newyear_rabbit_g_ani_01",
// 		"ropepack_newyear_rabbit_g_ani_02",
// 		"ropepack_newyear_rabbit_g_ani_03",
// 		"ropepack_newyear_rabbit_c_ani_01",
// 		"ropepack_newyear_rabbit_c_ani_02",

// 		"ropepack_moonfestival_h_pic_01",
// 		"ropepack_moonfestival_h_pic_02",
// 		"ropepack_moonfestival_h_pic_03",
// 		"ropepack_moonfestival_h_pic_04",
// 		"ropepack_moonfestival_h_pic_05",
// 		"ropepack_moonfestival_h_pic_06",
// 		"ropepack_moonfestival_h_pic_07",
// 		"ropepack_moonfestival_h_pic_08",
// 		"ropepack_moonfestival_h_pic_09",
// 		"ropepack_moonfestival_g_pic_01",
// 		"ropepack_moonfestival_g_pic_02",
// 		"ropepack_moonfestival_c_pic_01",
// 		"ropepack_moonfestival_h_ani_01",
// 		"ropepack_moonfestival_g_ani_01",
// 		"ropepack_moonfestival_g_ani_02",
// 		"ropepack_moonfestival_c_ani_01",
// 		"ropepack_moonfestival_c_ani_02",

// 		"ropepack_3D_h_pic_01",
// 		"ropepack_3D_h_pic_02",
// 		"ropepack_3D_h_pic_03",
// 		"ropepack_3D_h_pic_04",
// 		"ropepack_3D_h_pic_05",
// 		"ropepack_3D_h_pic_06",
// 		"ropepack_3D_h_pic_07",
// 		"ropepack_3D_h_pic_08",
// 		"ropepack_3D_h_pic_09",
// 		"ropepack_3D_h_pic_10",
// 		"ropepack_3D_h_pic_11",
// 		"ropepack_3D_h_pic_12",
// 		"ropepack_3D_h_pic_13",
// 		"ropepack_3D_h_pic_14",
// 		"ropepack_3D_h_pic_15",
// 		"ropepack_3D_g_pic_01",
// 		"ropepack_3D_g_pic_02",
// 		"ropepack_3D_g_pic_03",
// 		"ropepack_3D_g_pic_04",
// 		"ropepack_3D_h_ani_01",
// 		"ropepack_3D_h_ani_02",
// 		"ropepack_3D_h_ani_03",
// 		"ropepack_3D_g_ani_01",
// 		"ropepack_3D_g_ani_02",
// 		"ropepack_3D_c_ani_01",
// 	}
// 	update = make([]string, 300)
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

// 	// 套紅包自定義
// 	// 音樂
// 	RopepackBgmStart:  update[0],
// 	RopepackBgmGaming: update[1],
// 	RopepackBgmEnd:    update[2],

// 	RopepackClassicHPic01: update[3],
// 	RopepackClassicHPic02: update[4],
// 	RopepackClassicHPic03: update[5],
// 	RopepackClassicHPic04: update[6],
// 	RopepackClassicHPic05: update[7],
// 	RopepackClassicHPic06: update[8],
// 	RopepackClassicHPic07: update[9],
// 	RopepackClassicHPic08: update[10],
// 	RopepackClassicHPic09: update[11],
// 	RopepackClassicHPic10: update[12],
// 	RopepackClassicGPic01: update[13],
// 	RopepackClassicGPic02: update[14],
// 	RopepackClassicGPic03: update[15],
// 	RopepackClassicGPic04: update[16],
// 	RopepackClassicGPic05: update[17],
// 	RopepackClassicGPic06: update[18],
// 	RopepackClassicHAni01: update[19],
// 	RopepackClassicGAni01: update[20],
// 	RopepackClassicGAni02: update[21],
// 	RopepackClassicCAni01: update[22],

// 	RopepackNewyearRabbitHPic01: update[23],
// 	RopepackNewyearRabbitHPic02: update[24],
// 	RopepackNewyearRabbitHPic03: update[25],
// 	RopepackNewyearRabbitHPic04: update[26],
// 	RopepackNewyearRabbitHPic05: update[27],
// 	RopepackNewyearRabbitHPic06: update[28],
// 	RopepackNewyearRabbitHPic07: update[29],
// 	RopepackNewyearRabbitHPic08: update[30],
// 	RopepackNewyearRabbitHPic09: update[31],
// 	RopepackNewyearRabbitGPic01: update[32],
// 	RopepackNewyearRabbitGPic02: update[33],
// 	RopepackNewyearRabbitGPic03: update[34],
// 	RopepackNewyearRabbitHAni01: update[35],
// 	RopepackNewyearRabbitGAni01: update[36],
// 	RopepackNewyearRabbitGAni02: update[37],
// 	RopepackNewyearRabbitGAni03: update[38],
// 	RopepackNewyearRabbitCAni01: update[39],
// 	RopepackNewyearRabbitCAni02: update[40],

// 	RopepackMoonfestivalHPic01: update[41],
// 	RopepackMoonfestivalHPic02: update[42],
// 	RopepackMoonfestivalHPic03: update[43],
// 	RopepackMoonfestivalHPic04: update[44],
// 	RopepackMoonfestivalHPic05: update[45],
// 	RopepackMoonfestivalHPic06: update[46],
// 	RopepackMoonfestivalHPic07: update[47],
// 	RopepackMoonfestivalHPic08: update[48],
// 	RopepackMoonfestivalHPic09: update[49],
// 	RopepackMoonfestivalGPic01: update[50],
// 	RopepackMoonfestivalGPic02: update[51],
// 	RopepackMoonfestivalCPic01: update[52],
// 	RopepackMoonfestivalHAni01: update[53],
// 	RopepackMoonfestivalGAni01: update[54],
// 	RopepackMoonfestivalGAni02: update[55],
// 	RopepackMoonfestivalCAni01: update[56],
// 	RopepackMoonfestivalCAni02: update[57],

// 	Ropepack3DHPic01: update[58],
// 	Ropepack3DHPic02: update[59],
// 	Ropepack3DHPic03: update[60],
// 	Ropepack3DHPic04: update[61],
// 	Ropepack3DHPic05: update[62],
// 	Ropepack3DHPic06: update[63],
// 	Ropepack3DHPic07: update[64],
// 	Ropepack3DHPic08: update[65],
// 	Ropepack3DHPic09: update[66],
// 	Ropepack3DHPic10: update[67],
// 	Ropepack3DHPic11: update[68],
// 	Ropepack3DHPic12: update[69],
// 	Ropepack3DHPic13: update[70],
// 	Ropepack3DHPic14: update[71],
// 	Ropepack3DHPic15: update[72],
// 	Ropepack3DGPic01: update[73],
// 	Ropepack3DGPic02: update[74],
// 	Ropepack3DGPic03: update[75],
// 	Ropepack3DGPic04: update[76],
// 	Ropepack3DHAni01: update[77],
// 	Ropepack3DHAni02: update[78],
// 	Ropepack3DHAni03: update[79],
// 	Ropepack3DGAni01: update[80],
// 	Ropepack3DGAni02: update[81],
// 	Ropepack3DCAni01: update[82],
// }
