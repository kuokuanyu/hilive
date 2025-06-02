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
	drawnumbersPictureFields = []PictureField{
		{FieldName: "draw_numbers_bgm_gaming", Path: "draw_numbers/%s/bgm/gaming.mp3"},

		{FieldName: "draw_numbers_classic_h_pic_01", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_01.jpg"},
		{FieldName: "draw_numbers_classic_h_pic_02", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_02.png"},
		{FieldName: "draw_numbers_classic_h_pic_03", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_03.png"},
		{FieldName: "draw_numbers_classic_h_pic_04", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_04.png"},
		{FieldName: "draw_numbers_classic_h_pic_05", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_05.png"},
		{FieldName: "draw_numbers_classic_h_pic_06", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_06.png"},
		{FieldName: "draw_numbers_classic_h_pic_07", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_07.png"},
		{FieldName: "draw_numbers_classic_h_pic_08", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_08.png"},
		{FieldName: "draw_numbers_classic_h_pic_09", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_09.png"},
		{FieldName: "draw_numbers_classic_h_pic_10", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_10.png"},
		{FieldName: "draw_numbers_classic_h_pic_11", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_11.png"},
		{FieldName: "draw_numbers_classic_h_pic_12", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_12.png"},
		{FieldName: "draw_numbers_classic_h_pic_13", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_13.png"},
		{FieldName: "draw_numbers_classic_h_pic_14", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_14.png"},
		{FieldName: "draw_numbers_classic_h_pic_15", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_15.png"},
		{FieldName: "draw_numbers_classic_h_pic_16", Path: "draw_numbers/classic/draw_numbers_classic_h_pic_16.png"},
		{FieldName: "draw_numbers_classic_h_ani_01", Path: "draw_numbers/classic/draw_numbers_classic_h_ani_01.png"},

		{FieldName: "draw_numbers_gold_h_pic_01", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_01.jpg"},
		{FieldName: "draw_numbers_gold_h_pic_02", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_02.png"},
		{FieldName: "draw_numbers_gold_h_pic_03", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_03.png"},
		{FieldName: "draw_numbers_gold_h_pic_04", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_04.png"},
		{FieldName: "draw_numbers_gold_h_pic_05", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_05.png"},
		{FieldName: "draw_numbers_gold_h_pic_06", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_06.png"},
		{FieldName: "draw_numbers_gold_h_pic_07", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_07.png"},
		{FieldName: "draw_numbers_gold_h_pic_08", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_08.png"},
		{FieldName: "draw_numbers_gold_h_pic_09", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_09.png"},
		{FieldName: "draw_numbers_gold_h_pic_10", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_10.png"},
		{FieldName: "draw_numbers_gold_h_pic_11", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_11.png"},
		{FieldName: "draw_numbers_gold_h_pic_12", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_12.png"},
		{FieldName: "draw_numbers_gold_h_pic_13", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_13.png"},
		{FieldName: "draw_numbers_gold_h_pic_14", Path: "draw_numbers/gold/draw_numbers_gold_h_pic_14.png"},
		{FieldName: "draw_numbers_gold_h_ani_01", Path: "draw_numbers/gold/draw_numbers_gold_h_ani_01.png"},
		{FieldName: "draw_numbers_gold_h_ani_02", Path: "draw_numbers/gold/draw_numbers_gold_h_ani_02.png"},
		{FieldName: "draw_numbers_gold_h_ani_03", Path: "draw_numbers/gold/draw_numbers_gold_h_ani_03.png"},

		{FieldName: "draw_numbers_newyear_dragon_h_pic_01", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_01.jpg"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_02", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_02.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_03", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_03.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_04", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_04.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_05", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_05.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_06", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_06.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_07", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_07.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_08", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_08.jpg"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_09", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_09.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_10", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_10.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_11", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_11.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_12", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_12.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_13", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_13.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_14", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_14.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_15", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_15.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_16", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_16.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_17", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_17.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_18", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_18.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_19", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_19.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_pic_20", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_20.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_ani_01", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_01.png"},
		{FieldName: "draw_numbers_newyear_dragon_h_ani_02", Path: "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_02.png"},

		{FieldName: "draw_numbers_cherry_h_pic_01", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_01.png"},
		{FieldName: "draw_numbers_cherry_h_pic_02", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_02.png"},
		{FieldName: "draw_numbers_cherry_h_pic_03", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_03.png"},
		{FieldName: "draw_numbers_cherry_h_pic_04", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_04.png"},
		{FieldName: "draw_numbers_cherry_h_pic_05", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_05.png"},
		{FieldName: "draw_numbers_cherry_h_pic_06", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_06.png"},
		{FieldName: "draw_numbers_cherry_h_pic_07", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_07.png"},
		{FieldName: "draw_numbers_cherry_h_pic_08", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_08.png"},
		{FieldName: "draw_numbers_cherry_h_pic_09", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_09.png"},
		{FieldName: "draw_numbers_cherry_h_pic_10", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_10.jpg"},
		{FieldName: "draw_numbers_cherry_h_pic_11", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_11.png"},
		{FieldName: "draw_numbers_cherry_h_pic_12", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_12.png"},
		{FieldName: "draw_numbers_cherry_h_pic_13", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_13.png"},
		{FieldName: "draw_numbers_cherry_h_pic_14", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_14.png"},
		{FieldName: "draw_numbers_cherry_h_pic_15", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_15.png"},
		{FieldName: "draw_numbers_cherry_h_pic_16", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_16.png"},
		{FieldName: "draw_numbers_cherry_h_pic_17", Path: "draw_numbers/cherry/draw_numbers_cherry_h_pic_17.png"},
		{FieldName: "draw_numbers_cherry_h_ani_01", Path: "draw_numbers/cherry/draw_numbers_cherry_h_ani_01.png"},
		{FieldName: "draw_numbers_cherry_h_ani_02", Path: "draw_numbers/cherry/draw_numbers_cherry_h_ani_02.png"},
		{FieldName: "draw_numbers_cherry_h_ani_03", Path: "draw_numbers/cherry/draw_numbers_cherry_h_ani_03.png"},
		{FieldName: "draw_numbers_cherry_h_ani_04", Path: "draw_numbers/cherry/draw_numbers_cherry_h_ani_04.png"},

		{FieldName: "draw_numbers_3D_space_h_pic_01", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_01.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_02", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_02.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_03", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_03.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_04", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_04.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_05", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_05.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_06", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_06.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_07", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_07.png"},
		{FieldName: "draw_numbers_3D_space_h_pic_08", Path: "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_08.png"},
	}
)

// GetDrawNumbersPanel 搖號抽獎
func (s *SystemTable) GetDrawNumbersPanel() (table Table) {
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
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/game/draw_numbers/" + id)
		}

		// 刪除遊戲場次時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

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
		picMap := BuildPictureMap(drawnumbersPictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "draw_numbers", values.Get("game_id"), model); err != nil {
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
		picMap := BuildPictureMap(drawnumbersPictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "draw_numbers", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增搖號抽獎遊戲資料(form-data)
// @Tags Draw_Numbers
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string true "是否開啟公布中獎名單倒數功能" Enums(open, close)
// @param second formData integer true "公布中獎名單倒數時間"
// @param allow formData string true "允許重複中獎" Enums(open, close)
// @param display_name formData string true "是否顯示中獎人員姓名頭像" Enums(open, close)
// @param topic formData string true "主題樣式" Enums(01_classic, 02_gold, 04_cherry)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers/form [post]
func POSTDrawNumbers(ctx *gin.Context) {
}

// @Summary 新增搖號抽獎獎品資料(form-data)
// @Tags Draw_Numbers Prize
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
// @Router /interact/game/draw_numbers/prize/form [post]
func POSTDrawNumbersPrize(ctx *gin.Context) {
}

// @Summary 編輯搖號抽獎遊戲資料(form-data)
// @Tags Draw_Numbers
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param limit_time formData string false "是否開啟公布中獎名單倒數功能" Enums(open, close)
// @param second formData integer false "公布中獎名單倒數時間"
// @param allow formData string false "允許重複中獎" Enums(open, close)
// @param display_name formData string false "是否顯示中獎人員姓名頭像" Enums(open, close)
// @param topic formData string false "主題樣式" Enums(01_classic, 02_gold, 04_cherry)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param game_order formData integer false "game_order"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers/form [put]
func PUTDrawNumber(ctx *gin.Context) {
}

// @Summary 編輯搖號抽獎獎品資料(form-data)
// @Tags Draw_Numbers Prize
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
// @Router /interact/game/draw_numbers/prize/form [put]
func PUTDrawNumbersPrize(ctx *gin.Context) {
}

// @Summary 刪除搖號抽獎遊戲資料(form-data)
// @Tags Draw_Numbers
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers/form [delete]
func DELETEDrawNumbers(ctx *gin.Context) {
}

// @Summary 刪除搖號抽獎獎品資料(form-data)
// @Tags Draw_Numbers Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers/prize/form [delete]
func DELETEDrawNumbersPrize(ctx *gin.Context) {
}

// @Summary 搖號抽獎遊戲JSON資料
// @Tags Draw_Numbers
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers [get]
func DrawNumbersJSON(ctx *gin.Context) {
}

// @Summary 搖號抽獎獎品JSON資料(剩餘數量大於0的獎品)
// @Tags Draw_Numbers Prize
// @version 1.0
// @Accept  json
// @param game_id query string true "活動ID"
// @param isall query string false "是否顯示全部資料"
// @param is_array query string false "是否為陣列"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/draw_numbers/prize [get]
func DrawNumbersPrizeJSON(ctx *gin.Context) {
}

// pics = []string{
// 搖號抽獎自定義
// 音樂
// "draw_numbers/%s/bgm/gaming.mp3",

// "draw_numbers/classic/draw_numbers_classic_h_pic_01.jpg",
// "draw_numbers/classic/draw_numbers_classic_h_pic_02.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_03.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_04.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_05.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_06.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_07.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_08.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_09.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_10.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_11.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_12.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_13.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_14.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_15.png",
// "draw_numbers/classic/draw_numbers_classic_h_pic_16.png",
// "draw_numbers/classic/draw_numbers_classic_h_ani_01.png",

// "draw_numbers/gold/draw_numbers_gold_h_pic_01.jpg",
// "draw_numbers/gold/draw_numbers_gold_h_pic_02.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_03.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_04.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_05.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_06.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_07.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_08.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_09.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_10.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_11.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_12.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_13.png",
// "draw_numbers/gold/draw_numbers_gold_h_pic_14.png",
// "draw_numbers/gold/draw_numbers_gold_h_ani_01.png",
// "draw_numbers/gold/draw_numbers_gold_h_ani_02.png",
// "draw_numbers/gold/draw_numbers_gold_h_ani_03.png",

// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_01.jpg",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_02.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_03.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_04.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_05.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_06.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_07.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_08.jpg",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_09.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_10.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_11.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_12.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_13.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_14.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_15.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_16.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_17.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_18.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_19.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_20.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_01.png",
// "draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_02.png",

// "draw_numbers/cherry/draw_numbers_cherry_h_pic_01.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_02.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_03.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_04.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_05.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_06.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_07.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_08.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_09.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_10.jpg",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_11.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_12.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_13.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_14.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_15.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_16.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_pic_17.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_ani_01.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_ani_02.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_ani_03.png",
// "draw_numbers/cherry/draw_numbers_cherry_h_ani_04.png",

// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_01.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_02.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_03.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_04.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_05.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_06.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_07.png",
// "draw_numbers/3D_space/draw_numbers_3D_space_h_pic_08.png",
// }
// fields = []string{
// 搖號抽獎自定義
// 音樂
// "draw_numbers_bgm_gaming",

// "draw_numbers_classic_h_pic_01",
// "draw_numbers_classic_h_pic_02",
// "draw_numbers_classic_h_pic_03",
// "draw_numbers_classic_h_pic_04",
// "draw_numbers_classic_h_pic_05",
// "draw_numbers_classic_h_pic_06",
// "draw_numbers_classic_h_pic_07",
// "draw_numbers_classic_h_pic_08",
// "draw_numbers_classic_h_pic_09",
// "draw_numbers_classic_h_pic_10",
// "draw_numbers_classic_h_pic_11",
// "draw_numbers_classic_h_pic_12",
// "draw_numbers_classic_h_pic_13",
// "draw_numbers_classic_h_pic_14",
// "draw_numbers_classic_h_pic_15",
// "draw_numbers_classic_h_pic_16",
// "draw_numbers_classic_h_ani_01",

// "draw_numbers_gold_h_pic_01",
// "draw_numbers_gold_h_pic_02",
// "draw_numbers_gold_h_pic_03",
// "draw_numbers_gold_h_pic_04",
// "draw_numbers_gold_h_pic_05",
// "draw_numbers_gold_h_pic_06",
// "draw_numbers_gold_h_pic_07",
// "draw_numbers_gold_h_pic_08",
// "draw_numbers_gold_h_pic_09",
// "draw_numbers_gold_h_pic_10",
// "draw_numbers_gold_h_pic_11",
// "draw_numbers_gold_h_pic_12",
// "draw_numbers_gold_h_pic_13",
// "draw_numbers_gold_h_pic_14",
// "draw_numbers_gold_h_ani_01",
// "draw_numbers_gold_h_ani_02",
// "draw_numbers_gold_h_ani_03",

// "draw_numbers_newyear_dragon_h_pic_01",
// "draw_numbers_newyear_dragon_h_pic_02",
// "draw_numbers_newyear_dragon_h_pic_03",
// "draw_numbers_newyear_dragon_h_pic_04",
// "draw_numbers_newyear_dragon_h_pic_05",
// "draw_numbers_newyear_dragon_h_pic_06",
// "draw_numbers_newyear_dragon_h_pic_07",
// "draw_numbers_newyear_dragon_h_pic_08",
// "draw_numbers_newyear_dragon_h_pic_09",
// "draw_numbers_newyear_dragon_h_pic_10",
// "draw_numbers_newyear_dragon_h_pic_11",
// "draw_numbers_newyear_dragon_h_pic_12",
// "draw_numbers_newyear_dragon_h_pic_13",
// "draw_numbers_newyear_dragon_h_pic_14",
// "draw_numbers_newyear_dragon_h_pic_15",
// "draw_numbers_newyear_dragon_h_pic_16",
// "draw_numbers_newyear_dragon_h_pic_17",
// "draw_numbers_newyear_dragon_h_pic_18",
// "draw_numbers_newyear_dragon_h_pic_19",
// "draw_numbers_newyear_dragon_h_pic_20",
// "draw_numbers_newyear_dragon_h_ani_01",
// "draw_numbers_newyear_dragon_h_ani_02",

// "draw_numbers_cherry_h_pic_01",
// "draw_numbers_cherry_h_pic_02",
// "draw_numbers_cherry_h_pic_03",
// "draw_numbers_cherry_h_pic_04",
// "draw_numbers_cherry_h_pic_05",
// "draw_numbers_cherry_h_pic_06",
// "draw_numbers_cherry_h_pic_07",
// "draw_numbers_cherry_h_pic_08",
// "draw_numbers_cherry_h_pic_09",
// "draw_numbers_cherry_h_pic_10",
// "draw_numbers_cherry_h_pic_11",
// "draw_numbers_cherry_h_pic_12",
// "draw_numbers_cherry_h_pic_13",
// "draw_numbers_cherry_h_pic_14",
// "draw_numbers_cherry_h_pic_15",
// "draw_numbers_cherry_h_pic_16",
// "draw_numbers_cherry_h_pic_17",
// "draw_numbers_cherry_h_ani_01",
// "draw_numbers_cherry_h_ani_02",
// "draw_numbers_cherry_h_ani_03",
// "draw_numbers_cherry_h_ani_04",

// "draw_numbers_3D_space_h_pic_01",
// "draw_numbers_3D_space_h_pic_02",
// "draw_numbers_3D_space_h_pic_03",
// "draw_numbers_3D_space_h_pic_04",
// "draw_numbers_3D_space_h_pic_05",
// "draw_numbers_3D_space_h_pic_06",
// "draw_numbers_3D_space_h_pic_07",
// "draw_numbers_3D_space_h_pic_08",
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

// models.EditGameModel{
// 	UserID:        values.Get("user"),
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        values.Get("second"),
// 	MaxPeople:     "0",
// 	People:        "0",
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
// 	DisplayName:   values.Get("display_name"),
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

// 	// 自定義圖片
// 	// 音樂
// 	DrawNumbersBgmGaming: update[0],

// 	DrawNumbersClassicHPic01: update[1],
// 	DrawNumbersClassicHPic02: update[2],
// 	DrawNumbersClassicHPic03: update[3],
// 	DrawNumbersClassicHPic04: update[4],
// 	DrawNumbersClassicHPic05: update[5],
// 	DrawNumbersClassicHPic06: update[6],
// 	DrawNumbersClassicHPic07: update[7],
// 	DrawNumbersClassicHPic08: update[8],
// 	DrawNumbersClassicHPic09: update[9],
// 	DrawNumbersClassicHPic10: update[10],
// 	DrawNumbersClassicHPic11: update[11],
// 	DrawNumbersClassicHPic12: update[12],
// 	DrawNumbersClassicHPic13: update[13],
// 	DrawNumbersClassicHPic14: update[14],
// 	DrawNumbersClassicHPic15: update[15],
// 	DrawNumbersClassicHPic16: update[16],
// 	DrawNumbersClassicHAni01: update[17],

// 	DrawNumbersGoldHPic01: update[18],
// 	DrawNumbersGoldHPic02: update[19],
// 	DrawNumbersGoldHPic03: update[20],
// 	DrawNumbersGoldHPic04: update[21],
// 	DrawNumbersGoldHPic05: update[22],
// 	DrawNumbersGoldHPic06: update[23],
// 	DrawNumbersGoldHPic07: update[24],
// 	DrawNumbersGoldHPic08: update[25],
// 	DrawNumbersGoldHPic09: update[26],
// 	DrawNumbersGoldHPic10: update[27],
// 	DrawNumbersGoldHPic11: update[28],
// 	DrawNumbersGoldHPic12: update[29],
// 	DrawNumbersGoldHPic13: update[30],
// 	DrawNumbersGoldHPic14: update[31],
// 	DrawNumbersGoldHAni01: update[32],
// 	DrawNumbersGoldHAni02: update[33],
// 	DrawNumbersGoldHAni03: update[34],

// 	DrawNumbersNewyearDragonHPic01: update[35],
// 	DrawNumbersNewyearDragonHPic02: update[36],
// 	DrawNumbersNewyearDragonHPic03: update[37],
// 	DrawNumbersNewyearDragonHPic04: update[38],
// 	DrawNumbersNewyearDragonHPic05: update[39],
// 	DrawNumbersNewyearDragonHPic06: update[40],
// 	DrawNumbersNewyearDragonHPic07: update[41],
// 	DrawNumbersNewyearDragonHPic08: update[42],
// 	DrawNumbersNewyearDragonHPic09: update[43],
// 	DrawNumbersNewyearDragonHPic10: update[44],
// 	DrawNumbersNewyearDragonHPic11: update[45],
// 	DrawNumbersNewyearDragonHPic12: update[46],
// 	DrawNumbersNewyearDragonHPic13: update[47],
// 	DrawNumbersNewyearDragonHPic14: update[48],
// 	DrawNumbersNewyearDragonHPic15: update[49],
// 	DrawNumbersNewyearDragonHPic16: update[50],
// 	DrawNumbersNewyearDragonHPic17: update[51],
// 	DrawNumbersNewyearDragonHPic18: update[52],
// 	DrawNumbersNewyearDragonHPic19: update[53],
// 	DrawNumbersNewyearDragonHPic20: update[54],
// 	DrawNumbersNewyearDragonHAni01: update[55],
// 	DrawNumbersNewyearDragonHAni02: update[56],

// 	DrawNumbersCherryHPic01: update[57],
// 	DrawNumbersCherryHPic02: update[58],
// 	DrawNumbersCherryHPic03: update[59],
// 	DrawNumbersCherryHPic04: update[60],
// 	DrawNumbersCherryHPic05: update[61],
// 	DrawNumbersCherryHPic06: update[62],
// 	DrawNumbersCherryHPic07: update[63],
// 	DrawNumbersCherryHPic08: update[64],
// 	DrawNumbersCherryHPic09: update[65],
// 	DrawNumbersCherryHPic10: update[66],
// 	DrawNumbersCherryHPic11: update[67],
// 	DrawNumbersCherryHPic12: update[68],
// 	DrawNumbersCherryHPic13: update[69],
// 	DrawNumbersCherryHPic14: update[70],
// 	DrawNumbersCherryHPic15: update[71],
// 	DrawNumbersCherryHPic16: update[72],
// 	DrawNumbersCherryHPic17: update[73],
// 	DrawNumbersCherryHAni01: update[74],
// 	DrawNumbersCherryHAni02: update[75],
// 	DrawNumbersCherryHAni03: update[76],
// 	DrawNumbersCherryHAni04: update[77],

// 	DrawNumbers3DSpaceHPic01: update[78],
// 	DrawNumbers3DSpaceHPic02: update[79],
// 	DrawNumbers3DSpaceHPic03: update[80],
// 	DrawNumbers3DSpaceHPic04: update[81],
// 	DrawNumbers3DSpaceHPic05: update[82],
// 	DrawNumbers3DSpaceHPic06: update[83],
// 	DrawNumbers3DSpaceHPic07: update[84],
// 	DrawNumbers3DSpaceHPic08: update[85],
// }

// var (
// 	pics = []string{
// 		// 搖號抽獎自定義
// 		// 音樂
// 		"draw_numbers/%s/bgm/gaming.mp3",

// 		"draw_numbers/classic/draw_numbers_classic_h_pic_01.jpg",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_02.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_03.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_04.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_05.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_06.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_07.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_08.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_09.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_10.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_11.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_12.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_13.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_14.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_15.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_pic_16.png",
// 		"draw_numbers/classic/draw_numbers_classic_h_ani_01.png",

// 		"draw_numbers/gold/draw_numbers_gold_h_pic_01.jpg",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_02.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_03.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_04.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_05.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_06.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_07.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_08.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_09.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_10.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_11.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_12.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_13.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_pic_14.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_ani_01.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_ani_02.png",
// 		"draw_numbers/gold/draw_numbers_gold_h_ani_03.png",

// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_01.jpg",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_02.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_03.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_04.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_05.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_06.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_07.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_08.jpg",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_09.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_10.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_11.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_12.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_13.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_14.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_15.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_16.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_17.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_18.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_19.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_pic_20.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_01.png",
// 		"draw_numbers/newyear_dragon/draw_numbers_newyear_dragon_h_ani_02.png",

// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_01.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_02.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_03.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_04.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_05.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_06.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_07.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_08.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_09.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_10.jpg",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_11.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_12.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_13.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_14.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_15.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_16.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_pic_17.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_ani_01.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_ani_02.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_ani_03.png",
// 		"draw_numbers/cherry/draw_numbers_cherry_h_ani_04.png",

// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_01.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_02.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_03.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_04.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_05.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_06.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_07.png",
// 		"draw_numbers/3D_space/draw_numbers_3D_space_h_pic_08.png",
// 	}
// 	fields = []string{
// 		// 搖號抽獎自定義
// 		// 音樂
// 		"draw_numbers_bgm_gaming",

// 		"draw_numbers_classic_h_pic_01",
// 		"draw_numbers_classic_h_pic_02",
// 		"draw_numbers_classic_h_pic_03",
// 		"draw_numbers_classic_h_pic_04",
// 		"draw_numbers_classic_h_pic_05",
// 		"draw_numbers_classic_h_pic_06",
// 		"draw_numbers_classic_h_pic_07",
// 		"draw_numbers_classic_h_pic_08",
// 		"draw_numbers_classic_h_pic_09",
// 		"draw_numbers_classic_h_pic_10",
// 		"draw_numbers_classic_h_pic_11",
// 		"draw_numbers_classic_h_pic_12",
// 		"draw_numbers_classic_h_pic_13",
// 		"draw_numbers_classic_h_pic_14",
// 		"draw_numbers_classic_h_pic_15",
// 		"draw_numbers_classic_h_pic_16",
// 		"draw_numbers_classic_h_ani_01",

// 		"draw_numbers_gold_h_pic_01",
// 		"draw_numbers_gold_h_pic_02",
// 		"draw_numbers_gold_h_pic_03",
// 		"draw_numbers_gold_h_pic_04",
// 		"draw_numbers_gold_h_pic_05",
// 		"draw_numbers_gold_h_pic_06",
// 		"draw_numbers_gold_h_pic_07",
// 		"draw_numbers_gold_h_pic_08",
// 		"draw_numbers_gold_h_pic_09",
// 		"draw_numbers_gold_h_pic_10",
// 		"draw_numbers_gold_h_pic_11",
// 		"draw_numbers_gold_h_pic_12",
// 		"draw_numbers_gold_h_pic_13",
// 		"draw_numbers_gold_h_pic_14",
// 		"draw_numbers_gold_h_ani_01",
// 		"draw_numbers_gold_h_ani_02",
// 		"draw_numbers_gold_h_ani_03",

// 		"draw_numbers_newyear_dragon_h_pic_01",
// 		"draw_numbers_newyear_dragon_h_pic_02",
// 		"draw_numbers_newyear_dragon_h_pic_03",
// 		"draw_numbers_newyear_dragon_h_pic_04",
// 		"draw_numbers_newyear_dragon_h_pic_05",
// 		"draw_numbers_newyear_dragon_h_pic_06",
// 		"draw_numbers_newyear_dragon_h_pic_07",
// 		"draw_numbers_newyear_dragon_h_pic_08",
// 		"draw_numbers_newyear_dragon_h_pic_09",
// 		"draw_numbers_newyear_dragon_h_pic_10",
// 		"draw_numbers_newyear_dragon_h_pic_11",
// 		"draw_numbers_newyear_dragon_h_pic_12",
// 		"draw_numbers_newyear_dragon_h_pic_13",
// 		"draw_numbers_newyear_dragon_h_pic_14",
// 		"draw_numbers_newyear_dragon_h_pic_15",
// 		"draw_numbers_newyear_dragon_h_pic_16",
// 		"draw_numbers_newyear_dragon_h_pic_17",
// 		"draw_numbers_newyear_dragon_h_pic_18",
// 		"draw_numbers_newyear_dragon_h_pic_19",
// 		"draw_numbers_newyear_dragon_h_pic_20",
// 		"draw_numbers_newyear_dragon_h_ani_01",
// 		"draw_numbers_newyear_dragon_h_ani_02",

// 		"draw_numbers_cherry_h_pic_01",
// 		"draw_numbers_cherry_h_pic_02",
// 		"draw_numbers_cherry_h_pic_03",
// 		"draw_numbers_cherry_h_pic_04",
// 		"draw_numbers_cherry_h_pic_05",
// 		"draw_numbers_cherry_h_pic_06",
// 		"draw_numbers_cherry_h_pic_07",
// 		"draw_numbers_cherry_h_pic_08",
// 		"draw_numbers_cherry_h_pic_09",
// 		"draw_numbers_cherry_h_pic_10",
// 		"draw_numbers_cherry_h_pic_11",
// 		"draw_numbers_cherry_h_pic_12",
// 		"draw_numbers_cherry_h_pic_13",
// 		"draw_numbers_cherry_h_pic_14",
// 		"draw_numbers_cherry_h_pic_15",
// 		"draw_numbers_cherry_h_pic_16",
// 		"draw_numbers_cherry_h_pic_17",
// 		"draw_numbers_cherry_h_ani_01",
// 		"draw_numbers_cherry_h_ani_02",
// 		"draw_numbers_cherry_h_ani_03",
// 		"draw_numbers_cherry_h_ani_04",

// 		"draw_numbers_3D_space_h_pic_01",
// 		"draw_numbers_3D_space_h_pic_02",
// 		"draw_numbers_3D_space_h_pic_03",
// 		"draw_numbers_3D_space_h_pic_04",
// 		"draw_numbers_3D_space_h_pic_05",
// 		"draw_numbers_3D_space_h_pic_06",
// 		"draw_numbers_3D_space_h_pic_07",
// 		"draw_numbers_3D_space_h_pic_08",
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
// 	ActivityID:    values.Get("activity_id"),
// 	GameID:        values.Get("game_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     values.Get("limit_time"),
// 	Second:        values.Get("second"),
// 	MaxPeople:     "",
// 	People:        "",
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
// 	DisplayName:   values.Get("display_name"),
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
// 	// 音樂
// 	DrawNumbersBgmGaming: update[0],

// 	DrawNumbersClassicHPic01: update[1],
// 	DrawNumbersClassicHPic02: update[2],
// 	DrawNumbersClassicHPic03: update[3],
// 	DrawNumbersClassicHPic04: update[4],
// 	DrawNumbersClassicHPic05: update[5],
// 	DrawNumbersClassicHPic06: update[6],
// 	DrawNumbersClassicHPic07: update[7],
// 	DrawNumbersClassicHPic08: update[8],
// 	DrawNumbersClassicHPic09: update[9],
// 	DrawNumbersClassicHPic10: update[10],
// 	DrawNumbersClassicHPic11: update[11],
// 	DrawNumbersClassicHPic12: update[12],
// 	DrawNumbersClassicHPic13: update[13],
// 	DrawNumbersClassicHPic14: update[14],
// 	DrawNumbersClassicHPic15: update[15],
// 	DrawNumbersClassicHPic16: update[16],
// 	DrawNumbersClassicHAni01: update[17],

// 	DrawNumbersGoldHPic01: update[18],
// 	DrawNumbersGoldHPic02: update[19],
// 	DrawNumbersGoldHPic03: update[20],
// 	DrawNumbersGoldHPic04: update[21],
// 	DrawNumbersGoldHPic05: update[22],
// 	DrawNumbersGoldHPic06: update[23],
// 	DrawNumbersGoldHPic07: update[24],
// 	DrawNumbersGoldHPic08: update[25],
// 	DrawNumbersGoldHPic09: update[26],
// 	DrawNumbersGoldHPic10: update[27],
// 	DrawNumbersGoldHPic11: update[28],
// 	DrawNumbersGoldHPic12: update[29],
// 	DrawNumbersGoldHPic13: update[30],
// 	DrawNumbersGoldHPic14: update[31],
// 	DrawNumbersGoldHAni01: update[32],
// 	DrawNumbersGoldHAni02: update[33],
// 	DrawNumbersGoldHAni03: update[34],

// 	DrawNumbersNewyearDragonHPic01: update[35],
// 	DrawNumbersNewyearDragonHPic02: update[36],
// 	DrawNumbersNewyearDragonHPic03: update[37],
// 	DrawNumbersNewyearDragonHPic04: update[38],
// 	DrawNumbersNewyearDragonHPic05: update[39],
// 	DrawNumbersNewyearDragonHPic06: update[40],
// 	DrawNumbersNewyearDragonHPic07: update[41],
// 	DrawNumbersNewyearDragonHPic08: update[42],
// 	DrawNumbersNewyearDragonHPic09: update[43],
// 	DrawNumbersNewyearDragonHPic10: update[44],
// 	DrawNumbersNewyearDragonHPic11: update[45],
// 	DrawNumbersNewyearDragonHPic12: update[46],
// 	DrawNumbersNewyearDragonHPic13: update[47],
// 	DrawNumbersNewyearDragonHPic14: update[48],
// 	DrawNumbersNewyearDragonHPic15: update[49],
// 	DrawNumbersNewyearDragonHPic16: update[50],
// 	DrawNumbersNewyearDragonHPic17: update[51],
// 	DrawNumbersNewyearDragonHPic18: update[52],
// 	DrawNumbersNewyearDragonHPic19: update[53],
// 	DrawNumbersNewyearDragonHPic20: update[54],
// 	DrawNumbersNewyearDragonHAni01: update[55],
// 	DrawNumbersNewyearDragonHAni02: update[56],

// 	DrawNumbersCherryHPic01: update[57],
// 	DrawNumbersCherryHPic02: update[58],
// 	DrawNumbersCherryHPic03: update[59],
// 	DrawNumbersCherryHPic04: update[60],
// 	DrawNumbersCherryHPic05: update[61],
// 	DrawNumbersCherryHPic06: update[62],
// 	DrawNumbersCherryHPic07: update[63],
// 	DrawNumbersCherryHPic08: update[64],
// 	DrawNumbersCherryHPic09: update[65],
// 	DrawNumbersCherryHPic10: update[66],
// 	DrawNumbersCherryHPic11: update[67],
// 	DrawNumbersCherryHPic12: update[68],
// 	DrawNumbersCherryHPic13: update[69],
// 	DrawNumbersCherryHPic14: update[70],
// 	DrawNumbersCherryHPic15: update[71],
// 	DrawNumbersCherryHPic16: update[72],
// 	DrawNumbersCherryHPic17: update[73],
// 	DrawNumbersCherryHAni01: update[74],
// 	DrawNumbersCherryHAni02: update[75],
// 	DrawNumbersCherryHAni03: update[76],
// 	DrawNumbersCherryHAni04: update[77],

// 	DrawNumbers3DSpaceHPic01: update[78],
// 	DrawNumbers3DSpaceHPic02: update[79],
// 	DrawNumbers3DSpaceHPic03: update[80],
// 	DrawNumbers3DSpaceHPic04: update[81],
// 	DrawNumbers3DSpaceHPic05: update[82],
// 	DrawNumbers3DSpaceHPic06: update[83],
// 	DrawNumbers3DSpaceHPic07: update[84],
// 	DrawNumbers3DSpaceHPic08: update[85],
// }
