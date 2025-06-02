package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// GameVoteRecordModel 資料表欄位
type GameVoteRecordModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	OptionID   string `json:"option_id" example:"option_id"`
	Round      int64  `json:"round" example:"0"`
	Score      int64  `json:"score" example:"0"`
	CreatedAt  string `json:"created_at" example:"0"`

	// 投票選項資訊
	OptionName string `json:"option_name" example:"option_name"`

	// 用戶資訊
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ExtEmail string `json:"ext_email"`

	// 遊戲資訊
	Title string `json:"title" example:"Game Title"`

	// 用戶自定義資料
	Ext1   string `json:"ext_1"`
	Ext2   string `json:"ext_2"`
	Ext3   string `json:"ext_3"`
	Ext4   string `json:"ext_4"`
	Ext5   string `json:"ext_5"`
	Ext6   string `json:"ext_6"`
	Ext7   string `json:"ext_7"`
	Ext8   string `json:"ext_8"`
	Ext9   string `json:"ext_9"`
	Ext10  string `json:"ext_10"`
	Number int64  `json:"number"` // 抽獎號碼

	// 自定義欄位資訊
	Ext1Name     string `json:"ext_1_name"`
	Ext1Type     string `json:"ext_1_type"`
	Ext1Options  string `json:"ext_1_options"`
	Ext1Required string `json:"ext_1_required"`

	Ext2Name     string `json:"ext_2_name"`
	Ext2Type     string `json:"ext_2_type"`
	Ext2Options  string `json:"ext_2_options"`
	Ext2Required string `json:"ext_2_required"`

	Ext3Name     string `json:"ext_3_name"`
	Ext3Type     string `json:"ext_3_type"`
	Ext3Options  string `json:"ext_3_options"`
	Ext3Required string `json:"ext_3_required"`

	Ext4Name     string `json:"ext_4_name"`
	Ext4Type     string `json:"ext_4_type"`
	Ext4Options  string `json:"ext_4_options"`
	Ext4Required string `json:"ext_4_required"`

	Ext5Name     string `json:"ext_5_name"`
	Ext5Type     string `json:"ext_5_type"`
	Ext5Options  string `json:"ext_5_options"`
	Ext5Required string `json:"ext_5_required"`

	Ext6Name     string `json:"ext_6_name"`
	Ext6Type     string `json:"ext_6_type"`
	Ext6Options  string `json:"ext_6_options"`
	Ext6Required string `json:"ext_6_required"`

	Ext7Name     string `json:"ext_7_name"`
	Ext7Type     string `json:"ext_7_type"`
	Ext7Options  string `json:"ext_7_options"`
	Ext7Required string `json:"ext_7_required"`

	Ext8Name     string `json:"ext_8_name"`
	Ext8Type     string `json:"ext_8_type"`
	Ext8Options  string `json:"ext_8_options"`
	Ext8Required string `json:"ext_8_required"`

	Ext9Name     string `json:"ext_9_name"`
	Ext9Type     string `json:"ext_9_type"`
	Ext9Options  string `json:"ext_9_options"`
	Ext9Required string `json:"ext_9_required"`

	Ext10Name     string `json:"ext_10_name"`
	Ext10Type     string `json:"ext_10_type"`
	Ext10Options  string `json:"ext_10_options"`
	Ext10Required string `json:"ext_10_required"`

	ExtEmailRequired string `json:"ext_email_required"`
	ExtPhoneRequired string `json:"ext_phone_required"`
	InfoPicture      string `json:"info_picture"`
}

// EditGameVoteRecordModel 資料表欄位
type EditGameVoteRecordModel struct {
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	OptionID   string `json:"option_id" example:"option_id"`
	Round      int64  `json:"round" example:"0"`
	Score      int64  `json:"score" example:"0"`
}

// DefaultGameVoteRecordModel 預設GameVoteRecordModel
func DefaultGameVoteRecordModel() GameVoteRecordModel {
	return GameVoteRecordModel{Base: Base{TableName: config.ACTIVITY_GAME_VOTE_RECORD_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameVoteRecordModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameVoteRecordModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameVoteRecordModel) SetDbConn(conn db.Connection) GameVoteRecordModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameVoteRecordModel) SetRedisConn(conn cache.Connection) GameVoteRecordModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (a GameVoteRecordModel) SetMongoConn(conn mongo.Connection) GameVoteRecordModel {
// 	a.MongoConn = conn
// 	return a
// }

// Export 匯出投票紀錄資料
func (a GameVoteRecordModel) Export(activityID, gameID, optionID, round string,
	limit, offset int64) (*excelize.File, error) {
	var (
		file                 = excelize.NewFile()      // 開啟EXCEL檔案
		index, _             = file.NewSheet("Sheet1") // 創建一個工作表
		customizeModel, err1 = DefaultCustomizeModel().
					SetConn(a.DbConn, a.RedisConn, a.MongoConn).
					Find(activityID) // 自定義欄位
		rowNames = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
			"K", "L", "M", "N", "O", "P", "Q",
		}
		colNames = []string{
			// 用戶資訊
			"用戶名稱",
			// 遊戲資訊
			"遊戲名稱",
			// 參加資訊
			"遊戲輪次", "投票選項", "分數",
		}
		// 自定義欄位是否必填
		extRequireds = []string{
			// 電話.信箱是否必填
			customizeModel.ExtPhoneRequired, customizeModel.ExtEmailRequired,
			// 自定義欄位
			customizeModel.Ext1Name, customizeModel.Ext2Name, customizeModel.Ext3Name,
			customizeModel.Ext4Name, customizeModel.Ext5Name, customizeModel.Ext6Name,
			customizeModel.Ext7Name, customizeModel.Ext8Name, customizeModel.Ext9Name,
			customizeModel.Ext10Name,
		}

		// 投票紀錄
		records, err = DefaultGameVoteRecordModel().
				SetConn(a.DbConn, a.RedisConn, a.MongoConn).
				FindAll(activityID, gameID, optionID, round, limit, offset)

		sheet                                 = "Sheet1"
		fileName, gameIDCN, userIDCN, roundCN string
	)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		return file, err1
	}

	// 設置活頁簿的默認工作表
	file.SetActiveSheet(index)

	// 判斷自定義欄位
	for i, extRequired := range extRequireds {
		if i == 0 {
			if extRequired == "true" {
				colNames = append(colNames, "電話")
			}
		} else if i == 1 {
			if extRequired == "true" {
				colNames = append(colNames, "信箱")
			}
		} else {
			if extRequired != "" {
				colNames = append(colNames, extRequired)
			}
		}
	}

	// 將欄位中文名稱寫入EXCEL第一行
	for i, name := range colNames {
		file.SetCellValue(sheet, rowNames[i]+"1", name) // 設置存儲格的值
	}

	// 將所有中獎人員寫入EXCEL中(從第二行開始)
	for i, record := range records {
		var (
			extValues = []string{
				// 電話.信箱
				record.Phone, record.ExtEmail,
				// 自定義欄位
				record.Ext1, record.Ext2, record.Ext3,
				record.Ext4, record.Ext5, record.Ext6,
				record.Ext7, record.Ext8, record.Ext9,
				record.Ext10,
			}
		)

		values := []interface{}{
			// 用戶資訊
			record.Name,
			// 遊戲資訊
			record.Title,
			// 參加資訊
			record.Round, record.OptionName, record.Score,
		}

		// 判斷自定義欄位
		for i, extRequired := range extRequireds {
			if i == 0 {
				if extRequired == "true" {
					values = append(values, extValues[i])
				}
			} else if i == 1 {
				if extRequired == "true" {
					values = append(values, extValues[i])
				}
			} else {
				if extRequired != "" {
					values = append(values, extValues[i])
				}
			}
		}

		for n, value := range values {
			file.SetCellValue(sheet, rowNames[n]+strconv.Itoa(i+2), value) // 設置存儲格的值
		}
	}

	// 處理excel檔案名稱
	if gameID != "" {
		gameIDCN = gameID
	} else if gameID == "" {
		gameIDCN = "所有遊戲場次"
	}
	// 投票遊戲
	if optionID != "" {
		userIDCN = optionID
	} else if optionID == "" {
		userIDCN = "所有選項"
	}
	// 判斷是否取得特定輪次資料
	if round != "" {
		roundCN = fmt.Sprintf("第%s輪", round)
	} else if round == "" {
		roundCN = "所有輪次"
	}

	fileName = "遊戲紀錄-" + activityID + "-" + "vote" + "-" +
		gameIDCN + "-" + userIDCN + "-" + roundCN
	// 儲存EXCEL
	if err := file.SaveAs(fmt.Sprintf(config.STORE_PATH+"/excel/%s.xlsx", fileName)); err != nil {
		return file, err
	}

	return file, nil
}

// FindAll 查詢投票遊戲所有投票紀錄
func (a GameVoteRecordModel) FindAll(activityID, gameID, optionID, round string,
	limit, offset int64) ([]GameVoteRecordModel, error) {
	var (
		recordsModel = make([]GameVoteRecordModel, 0)
		sql          = a.Table(a.Base.TableName).
				Select("activity_game_vote_record.id", "activity_game_vote_record.activity_id",
				"activity_game_vote_record.game_id", "activity_game_vote_record.user_id",
				"activity_game_vote_record.option_id", "activity_game_vote_record.score",
				"activity_game_vote_record.created_at", "activity_game_vote_record.round",

				// 遊戲場次
				// "activity_game.title",

				// 用戶
				"line_users.name", "line_users.avatar", "line_users.phone",
				"line_users.email", "line_users.ext_email",

				// 自定義欄位
				"activity_applysign.ext_1", "activity_customize.ext_1_name", "activity_customize.ext_1_type",
				"activity_customize.ext_1_options", "activity_customize.ext_1_required",
				"activity_applysign.ext_2", "activity_customize.ext_2_name", "activity_customize.ext_2_type",
				"activity_customize.ext_2_options", "activity_customize.ext_2_required",
				"activity_applysign.ext_3", "activity_customize.ext_3_name", "activity_customize.ext_3_type",
				"activity_customize.ext_3_options", "activity_customize.ext_3_required",
				"activity_applysign.ext_4", "activity_customize.ext_4_name", "activity_customize.ext_4_type",
				"activity_customize.ext_4_options", "activity_customize.ext_4_required",
				"activity_applysign.ext_5", "activity_customize.ext_5_name", "activity_customize.ext_5_type",
				"activity_customize.ext_5_options", "activity_customize.ext_5_required",
				"activity_applysign.ext_6", "activity_customize.ext_6_name", "activity_customize.ext_6_type",
				"activity_customize.ext_6_options", "activity_customize.ext_6_required",
				"activity_applysign.ext_7", "activity_customize.ext_7_name", "activity_customize.ext_7_type",
				"activity_customize.ext_7_options", "activity_customize.ext_7_required",
				"activity_applysign.ext_8", "activity_customize.ext_8_name", "activity_customize.ext_8_type",
				"activity_customize.ext_8_options", "activity_customize.ext_8_required",
				"activity_applysign.ext_9", "activity_customize.ext_9_name", "activity_customize.ext_9_type",
				"activity_customize.ext_9_options", "activity_customize.ext_9_required",
				"activity_applysign.ext_10", "activity_customize.ext_10_name", "activity_customize.ext_10_type",
				"activity_customize.ext_10_options", "activity_customize.ext_10_required",
				"activity_applysign.number",

				"activity_customize.ext_email_required", "activity_customize.ext_phone_required",
				"activity_customize.info_picture",

				// 投票選項資訊
				"activity_game_vote_option.option_name",
			// "activity_game_vote_option.option_picture",
			// "activity_game_vote_option.option_introduce",
			// "activity_game_vote_option.option_score",
			).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_record.option_id",
				FieldA1:   "activity_game_vote_option.option_id",
				Table:     "activity_game_vote_option",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_record.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_game_vote_record.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_record.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_record.activity_id",
				FieldA1:   "activity_applysign.activity_id",
				FieldB:    "activity_game_vote_record.user_id",
				FieldB1:   "activity_applysign.user_id",
				Table:     "activity_applysign",
				Operation: "="}).
			Where("activity_game_vote_record.activity_id", "=", activityID).
			OrderBy(
				"activity_game_vote_record.game_id", "asc",
				"activity_game_vote_record.round", "asc",
				"activity_game_vote_record.score", "desc",
			)
	)

	// 判斷參數是否為空
	if gameID != "" {
		sql = sql.WhereIn("activity_game_vote_record.game_id", interfaces(strings.Split(gameID, ",")))
	}
	if optionID != "" {
		sql = sql.WhereIn("activity_game_vote_record.option_id", interfaces(strings.Split(optionID, ",")))
	}

	if round != "" {
		sql = sql.WhereIn("activity_game_vote_record.round", interfaces(strings.Split(round, ",")))
	}

	if limit != 0 {
		sql = sql.Limit(limit)
	}
	if offset != 0 {
		sql = sql.Offset(offset)
	}

	items, err := sql.All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得投票紀錄，請重新查詢")
	}

	recordsModel = MapToVoteRecordModel(items)

	// 多個遊戲場次資料(mongo)
	gamesModel, err := DefaultGameModel().
		SetConn(a.DbConn, a.RedisConn, a.MongoConn).
		FindAll(activityID, "vote")
	if err != nil {
		return recordsModel, err
	}

	// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
	dataMap := make(map[string]GameModel)
	for _, game := range gamesModel {
		dataMap[game.GameID] = game
	}

	for i, item := range recordsModel {
		if dataItem, found := dataMap[item.GameID]; found {

			// 合併欄位（可依需要擴充）
			item.Title = dataItem.Title

			recordsModel[i] = item
		}
	}

	return recordsModel, nil
}

// FindUserRecord 查詢玩家投票紀錄資料
func (a GameVoteRecordModel) FindUserRecord(isRedis bool, gameID, userID string) ([]GameVoteRecordModel, error) {
	// log.Println("執行FindUserRecord")

	var (
		records = make([]GameVoteRecordModel, 0)
		sql     = a.Table(a.Base.TableName).
			Where("game_id", "=", gameID).
			OrderBy("id", "asc")
	)

	if userID != "" { // 查詢特定用戶
		sql = sql.Where("user_id", "=", userID)
	}

	// 從redis取得玩家投票遊戲紀錄資料(hash類型)
	if isRedis && userID != "" {
		// redis處理
		// 判斷redis裡是否有玩家投票紀錄資訊，有則不執行查詢資料表功能
		recordsJson, err := a.RedisConn.HashGetCache(config.GAME_VOTE_RECORDS_REDIS+gameID, userID) // 包含test測試資料，避免重複執行資料庫
		if err != nil || recordsJson == "" {
			// redis查詢不到玩家投票紀錄資料
			// 從資料庫讀取得玩家投票紀錄資料
			// log.Println("redis無法取得人員投票紀錄，查詢資料表")

			items, err := sql.All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得投票紀錄資訊，請重新查詢")
			}
			records = MapToVoteRecordModel(items) // 多筆資料

			// 將玩家投票紀錄資料添加至redis中(hash類型)
			if len(records) > 0 {
				// log.Println("資料表有資料，加入人員投票紀錄至redis")

				err = a.RedisConn.HashSetCache(config.GAME_VOTE_RECORDS_REDIS+gameID,
					userID, utils.JSON(records))
				if err != nil {
					return nil, errors.New("錯誤: 將人員投票紀錄寫入redis中發生問題")
				}

				// a.RedisConn.SetExpire(config.GAME_VOTE_RECORDS_REDIS+gameID,
				// 	config.REDIS_EXPIRE)
			} else if len(records) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				err = a.RedisConn.HashSetCache(config.GAME_VOTE_RECORDS_REDIS+gameID, userID,
					utils.JSON([]GameVoteRecordModel{{
						UserID: "test",
					}}))
				if err != nil {
					return nil, errors.New("錯誤: 將人員投票紀錄寫入redis中發生問題")
				}

				// a.RedisConn.SetExpire(config.GAME_VOTE_RECORDS_REDIS+gameID,
				// 	config.REDIS_EXPIRE)
			}
		} else {
			// log.Println("redis有資料")

			var recordModels = make([]GameVoteRecordModel, 0)
			// 成功從redis取得資料，解碼
			json.Unmarshal([]byte(recordsJson), &recordModels) // 包含test測試資料，避免重複執行資料庫

			// log.Println("資料: ", recordModels)

			for _, recordModel := range recordModels {
				if recordModel.UserID != "test" {
					// log.Println("不是測試資料")
					// 過濾測試資料
					records = append(records, recordModel)
				} else {
					// log.Println("目前為測試資料，跳過")
				}
			}

		}
	}

	// 查詢資料表
	if !isRedis {
		// log.Println("查詢資料表: ")

		items, err := sql.All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得投票紀錄資訊，請重新查詢")
		}
		records = MapToVoteRecordModel(items)

		// 判斷是查詢個人資料還是所有資料(因為redis的value值是陣列)
		if userID != "" {
			if len(records) > 0 {
				// log.Println("資料庫有資料")

				// 將玩家投票紀錄資料添加至redis中(hash類型)
				err = a.RedisConn.HashSetCache(config.GAME_VOTE_RECORDS_REDIS+gameID,
					userID, utils.JSON(records))

			} else if len(records) == 0 {
				// log.Println("資料庫無資料，加入test")

				// 將玩家投票紀錄資料添加至redis中(hash類型)
				err = a.RedisConn.HashSetCache(config.GAME_VOTE_RECORDS_REDIS+gameID,
					userID, utils.JSON([]GameVoteRecordModel{{
						UserID: "test",
					}}))
			}

			if err != nil {
				return nil, errors.New("錯誤: 將人員投票紀錄寫入redis中發生問題")
			}

			// a.RedisConn.SetExpire(config.GAME_VOTE_RECORDS_REDIS+gameID,
			// 	config.REDIS_EXPIRE)
		}
	}

	return records, nil
}

// MapToVoteRecordModel map轉換[]GameVoteRecordModel
func MapToVoteRecordModel(items []map[string]interface{}) []GameVoteRecordModel {
	var records = make([]GameVoteRecordModel, 0)
	for _, item := range items {
		var (
			record GameVoteRecordModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &record)

		records = append(records, record)
	}
	return records
}

// Add 新增投票紀錄資料(加入資料庫)
func (a GameVoteRecordModel) Add(isRedis bool, model EditGameVoteRecordModel) error {
	var (
		fieldValues = command.Value{
			"activity_id": model.ActivityID,
			"game_id":     model.GameID,
			"user_id":     model.UserID,
			"option_id":   model.OptionID,
			"round":       model.Round,
			"score":       model.Score,
		}
		records = make([]GameVoteRecordModel, 0)
		err     error
	)

	// score, err := strconv.Atoi(model.Score)
	// if err != nil {
	// 	return errors.New("錯誤: 分數資料發生問題")
	// }

	if isRedis {
		// 從redis取得資料，確定redis中有該玩家投票紀錄資料
		records, err = a.FindUserRecord(true, model.GameID, model.UserID)
		if err != nil {
			return errors.New("錯誤: 取得玩家投票紀錄資料發生問題")
		}
	}

	// 將投票紀錄寫入資料表
	id, err := a.Table(a.TableName).Insert(fieldValues)
	if err != nil {
		return errors.New("錯誤: 新增投票紀錄資料發生問題(資料庫)")
	}

	if isRedis {
		// 加入新的投票紀錄
		records = append(records, GameVoteRecordModel{
			ID:         id,
			ActivityID: model.ActivityID,
			GameID:     model.GameID,
			UserID:     model.UserID,
			OptionID:   model.OptionID,
			Round:      model.Round,
			Score:      model.Score,
		})

		// 將新的投票紀錄寫入redis中(hash類型)
		a.RedisConn.HashSetCache(config.GAME_VOTE_RECORDS_REDIS+model.GameID,
			model.UserID, utils.JSON(records))
		// a.RedisConn.SetExpire(config.GAME_VOTE_RECORDS_REDIS+model.GameID,
		// 	config.REDIS_EXPIRE)
	}

	return nil
}

// record.ID, _ = item["id"].(int64)
// record.ActivityID, _ = item["activity_id"].(string)
// record.GameID, _ = item["game_id"].(string)
// record.UserID, _ = item["user_id"].(string)
// record.OptionID, _ = item["option_id"].(string)
// record.Round, _ = item["round"].(int64)
// record.Score, _ = item["score"].(int64)

// record.OptionName, _ = item["option_name"].(string)

// // 用戶資訊
// record.Name, _ = item["name"].(string)
// record.Avatar, _ = item["avatar"].(string)
// record.Phone, _ = item["phone"].(string)
// record.Email, _ = item["email"].(string)
// record.ExtEmail, _ = item["ext_email"].(string)

// // 遊戲資訊
// record.Title, _ = item["title"].(string)

// // 自定義用戶資料
// record.Ext1, _ = item["ext_1"].(string)
// record.Ext2, _ = item["ext_2"].(string)
// record.Ext3, _ = item["ext_3"].(string)
// record.Ext4, _ = item["ext_4"].(string)
// record.Ext5, _ = item["ext_5"].(string)
// record.Ext6, _ = item["ext_6"].(string)
// record.Ext7, _ = item["ext_7"].(string)
// record.Ext8, _ = item["ext_8"].(string)
// record.Ext9, _ = item["ext_9"].(string)
// record.Ext10, _ = item["ext_10"].(string)
// record.Number, _ = item["number"].(int64)

// // 自定義欄位資訊
// record.Ext1Name, _ = item["ext_1_name"].(string)
// record.Ext1Type, _ = item["ext_1_type"].(string)
// record.Ext1Options, _ = item["ext_1_options"].(string)
// record.Ext1Required, _ = item["ext_1_required"].(string)

// record.Ext2Name, _ = item["ext_2_name"].(string)
// record.Ext2Type, _ = item["ext_2_type"].(string)
// record.Ext2Options, _ = item["ext_2_options"].(string)
// record.Ext2Required, _ = item["ext_2_required"].(string)

// record.Ext3Name, _ = item["ext_3_name"].(string)
// record.Ext3Type, _ = item["ext_3_type"].(string)
// record.Ext3Options, _ = item["ext_3_options"].(string)
// record.Ext3Required, _ = item["ext_3_required"].(string)

// record.Ext4Name, _ = item["ext_4_name"].(string)
// record.Ext4Type, _ = item["ext_4_type"].(string)
// record.Ext4Options, _ = item["ext_4_options"].(string)
// record.Ext4Required, _ = item["ext_4_required"].(string)

// record.Ext5Name, _ = item["ext_5_name"].(string)
// record.Ext5Type, _ = item["ext_5_type"].(string)
// record.Ext5Options, _ = item["ext_5_options"].(string)
// record.Ext5Required, _ = item["ext_5_required"].(string)

// record.Ext6Name, _ = item["ext_6_name"].(string)
// record.Ext6Type, _ = item["ext_6_type"].(string)
// record.Ext6Options, _ = item["ext_6_options"].(string)
// record.Ext6Required, _ = item["ext_6_required"].(string)

// record.Ext7Name, _ = item["ext_7_name"].(string)
// record.Ext7Type, _ = item["ext_7_type"].(string)
// record.Ext7Options, _ = item["ext_7_options"].(string)
// record.Ext7Required, _ = item["ext_7_required"].(string)

// record.Ext8Name, _ = item["ext_8_name"].(string)
// record.Ext8Type, _ = item["ext_8_type"].(string)
// record.Ext8Options, _ = item["ext_8_options"].(string)
// record.Ext8Required, _ = item["ext_8_required"].(string)

// record.Ext9Name, _ = item["ext_9_name"].(string)
// record.Ext9Type, _ = item["ext_9_type"].(string)
// record.Ext9Options, _ = item["ext_9_options"].(string)
// record.Ext9Required, _ = item["ext_9_required"].(string)

// record.Ext10Name, _ = item["ext_10_name"].(string)
// record.Ext10Type, _ = item["ext_10_type"].(string)
// record.Ext10Options, _ = item["ext_10_options"].(string)
// record.Ext10Required, _ = item["ext_10_required"].(string)

// record.ExtEmailRequired, _ = item["ext_email_required"].(string)
// record.ExtPhoneRequired, _ = item["ext_phone_required"].(string)
// record.InfoPicture, _ = item["info_picture"].(string)

// Update 更新投票紀錄資料
// func (a GameVoteRecordModel) Update(model EditGameVoteRecordModel) error {
// 	var (
// 		fieldValues = command.Value{}
// 		fields      = []string{"name", "leader", "option_id"}
// 		values      = []string{model.Name, model.Leader, model.OptionID}
// 	)

// 	for i, value := range values {
// 		if value != "" {
// 			fieldValues[fields[i]] = value
// 		}
// 	}
// 	if len(fieldValues) == 0 {
// 		return nil
// 	}

// 	return a.Table(a.Base.TableName).
// 		Where("activity_id",
