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
	"strconv"

	"github.com/xuri/excelize/v2"
)

// BlackStaffModel 資料表欄位
type BlackStaffModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id"`
	GameID     string `json:"game_id"`
	Game       string `json:"game"`
	UserID     string `json:"user_id"`
	Reason     string `json:"reason"`

	// 用戶資訊
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ExtEmail string `json:"ext_email"`

	// 遊戲資訊
	Title string `json:"title" example:"game title"`

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

// // NewBlackStaffModel 資料表欄位
// type NewBlackStaffModel struct {
// 	ActivityID string        `json:"activity_id" example:"activity_id"`
// 	GameID     string        `json:"game_id" example:"game_id"`
// 	LineUsers  []interface{} `json:"line_users"` // 批量新增用戶
// 	Game       string        `json:"game" example:"game"`
// 	Reason     string        `json:"reason" example:"reason"`
// 	// UserID     string   `json:"user_id" example:"user_id"`
// 	// Token      string   `json:"token" example:"token"`
// }

// EditBlackStaffModel 資料表欄位
type EditBlackStaffModel struct {
	ActivityID string        `json:"activity_id" example:"activity_id"`
	GameID     string        `json:"game_id" example:"game_id"`
	LineUsers  []interface{} `json:"line_users"` // 批量更新用戶
	Game       string        `json:"game" example:"game"`
	Reason     string        `json:"reason" example:"reason"`
	// UserID     string `json:"user_id" example:"user_id"`
	// Token      string `json:"token" example:"token"`
}

// DefaultBlackStaffModel 預設BlackStaffModel
func DefaultBlackStaffModel() BlackStaffModel {
	return BlackStaffModel{Base: Base{TableName: config.ACTIVITY_STAFF_BLACK_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m BlackStaffModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) BlackStaffModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m BlackStaffModel) SetDbConn(conn db.Connection) BlackStaffModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m BlackStaffModel) SetRedisConn(conn cache.Connection) BlackStaffModel {
// 	m.RedisConn = conn
// 	return m
// }

// // SetMongoConn 設定connection
// func (m BlackStaffModel) SetMongoConn(conn mongo.Connection) BlackStaffModel {
// 	m.MongoConn = conn
// 	return m
// }

// FindAll 查詢所有黑名單人員資料
func (m BlackStaffModel) FindAll(isRedis bool, activityID, gameID, game string) ([]BlackStaffModel, error) {
	var (
		sql = m.Table(m.Base.TableName).
			Select("activity_staff_black.id", "activity_staff_black.activity_id",
				"activity_staff_black.user_id", "activity_staff_black.game_id",
				"activity_staff_black.game", "activity_staff_black.reason",

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

				// 遊戲
				// "activity_game.title",

				// 用戶
				"line_users.name", "line_users.avatar", "line_users.phone",
				"line_users.email", "line_users.ext_email").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_black.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_staff_black.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_black.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="}).
			LeftJoin(command.Join{
											FieldA:    "activity_staff_black.activity_id",
											FieldA1:   "activity_applysign.activity_id",
											FieldB:    "activity_staff_black.user_id",
											FieldB1:   "activity_applysign.user_id",
											Table:     "activity_applysign",
											Operation: "="}).
			Where("activity_staff_black.activity_id", "=", activityID) // 該活動下所有黑名單資料
		blackStaffsModel = make([]BlackStaffModel, 0) // 過濾後的黑名單資料
		// filter           = bson.M{"activity_id": activityID} // 過濾條件
		// blackstaffs      = make([]string, 0)
		channel, key string
		err          error
	)

	// 從redis取得黑名單資料(set類型)
	if isRedis {
		// redis key名稱
		if game == "activity" {
			// 活動
			key = config.BLACK_STAFFS_ACTIVITY_REDIS + activityID
			channel = config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + activityID
		} else if game == "message" {
			// 訊息
			key = config.BLACK_STAFFS_MESSAGE_REDIS + activityID
			channel = config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + activityID
		} else if game == "question" {
			// 提問
			key = config.BLACK_STAFFS_QUESTION_REDIS + activityID
			channel = config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + activityID
		} else if game == "signname" {
			// 簽名
			key = config.BLACK_STAFFS_SIGNNAME_REDIS + activityID
			channel = config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + activityID
		} else if gameID != "" {
			// 遊戲
			key = config.BLACK_STAFFS_GAME_REDIS + gameID
			channel = config.CHANNEL_BLACK_STAFFS_GAME_REDIS + gameID
		}

		// 判斷redis裡是否有黑名單資訊，有則不執行查詢資料表功能
		blackstaffs, err := m.RedisConn.SetGetMembers(key)
		if err != nil {
			return blackStaffsModel, errors.New("錯誤: 無法取得黑名單人員資訊(redis)，請重新查詢")
		}
		// log.Println("黑名單人員處理前: ", blackstaffs)

		// 添加至BlackStaffModel中
		for _, blackStaff := range blackstaffs {
			// 加入一筆test資料於redis中，避免重複查詢資料庫
			// 過濾測試資料
			if blackStaff != "test" {
				blackStaffsModel = append(blackStaffsModel, BlackStaffModel{
					ActivityID: activityID,
					GameID:     gameID,
					Game:       game,
					UserID:     blackStaff,
				})
			}
		}
		// log.Println("黑名單人員處理後: ", len(blackStaffsModel))

		// redis裡無黑名單資料
		if len(blackstaffs) == 0 {
			// log.Println("黑名單裡無資料，資料庫處理")

			var (
				allBlackStaffs = make([]BlackStaffModel, 0) // 所有黑名單資料
				params         = []interface{}{key}         // redis資料
			)

			// 查詢該活動下所有黑名單資料
			items, err := sql.All()
			if err != nil {
				return blackStaffsModel, errors.New("錯誤: 無法取得黑名單人員資訊，請重新查詢")
			}
			// fmt.Println("items: ", items, "game_id: ", gameID, "game: ", game, "activity_id: ", activityID)

			// 所有黑名單資料
			allBlackStaffs = MapToBlackStaffModel(items)
			// 過濾黑名單
			for _, balckStaff := range allBlackStaffs {
				// 活動黑名單人員必須加入所有其他黑名單裡(訊息、提問、遊戲...等)
				if balckStaff.Game == "activity" ||
					(balckStaff.Game == game && balckStaff.GameID == gameID) {
					blackStaffsModel = append(blackStaffsModel, balckStaff)
					params = append(params, balckStaff.UserID)
				}
			}

			if len(blackStaffsModel) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				params = append(params, "test")
			}

			// 將黑名單資料設置至redis中
			m.RedisConn.SetAdd(params)
			// m.RedisConn.SetExpire(key, config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(channel, "修改資料")
		}
	} else {
		// 不使用redis
		var items []map[string]interface{}

		// 判斷參數是否為空
		if game != "" {
			sql = sql.Where("activity_staff_black.game", "=", game)

			// mongo過濾game參數
			// filter["game"] = bson.M{"game": game}
		}
		if gameID != "" {
			sql = sql.Where("activity_staff_black.game_id", "=", gameID)

			// mongo過濾game參數
			// filter["game"] = bson.M{"game_id": gameID}
		}

		items, err = sql.All()
		if err != nil {
			return blackStaffsModel, errors.New("錯誤: 無法取得黑名單人員資訊，請重新查詢")
		}

		blackStaffsModel = MapToBlackStaffModel(items)

		// 多個遊戲場次資料(mongo)
		gamesModel, err := DefaultGameModel().
			SetConn(m.DbConn, m.RedisConn, m.MongoConn).
			FindAll(activityID, game)
		if err != nil {
			return blackStaffsModel, err
		}

		// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
		dataMap := make(map[string]GameModel)
		for _, game := range gamesModel {
			dataMap[game.GameID] = game
		}

		for i, item := range blackStaffsModel {
			if dataItem, found := dataMap[item.GameID]; found {

				// 合併欄位（可依需要擴充）
				item.Title = dataItem.Title

				blackStaffsModel[i] = item
			}
		}
	}
	return blackStaffsModel, nil
}

// Export 匯出黑名單人員資料
func (m BlackStaffModel) Export(activityID, gameID, game string) (*excelize.File, error) {
	var (
		file                 = excelize.NewFile()      // 開啟EXCEL檔案
		index, _             = file.NewSheet("Sheet1") // 創建一個工作表
		customizeModel, err1 = DefaultCustomizeModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					Find(activityID) // 自定義欄位
		rowNames = []string{"A", "B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
			"S", "T", "U", "V", "W", "X", "Y", "Z"}
		colNames = []string{
			// 用戶資訊
			"用戶名稱", "抽獎號碼",
			// 遊戲資訊
			"遊戲類型", "遊戲名稱",
			// 黑名單資訊
			"黑名單原因",
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
		blackStaffModel, err2 = DefaultBlackStaffModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					FindAll(false, activityID, gameID, game)
		sheet                      = "Sheet1"
		fileName, gameCN, gameIDCN string
	)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err1 != nil || err2 != nil {
		if err1 != nil {
			return file, err1
		} else if err2 != nil {
			return file, err2
		}
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
	for i, staff := range blackStaffModel {
		var (
			newGame   string
			extValues = []string{
				// 電話.信箱
				staff.Phone, staff.ExtEmail,
				// 自定義欄位
				staff.Ext1, staff.Ext2, staff.Ext3,
				staff.Ext4, staff.Ext5, staff.Ext6,
				staff.Ext7, staff.Ext8, staff.Ext9,
				staff.Ext10,
			}
		)

		// 判斷遊戲類型
		if staff.Game == "redpack" {
			newGame = "搖紅包"
		} else if staff.Game == "ropepack" {
			newGame = "套紅包"
		} else if staff.Game == "whack_mole" {
			newGame = "敲敲樂"
		} else if staff.Game == "monopoly" {
			newGame = "鑑定師"
		} else if staff.Game == "draw_numbers" {
			newGame = "搖號抽獎"
		} else if staff.Game == "lottery" {
			newGame = "遊戲抽獎"
		} else if staff.Game == "QA" {
			newGame = "快問快答"
		} else if staff.Game == "activity" {
			newGame = "所有活動"
		} else if staff.Game == "message" {
			newGame = "訊息頁"
		} else if staff.Game == "question" {
			newGame = "提問頁"
		} else if staff.Game == "tugofwar" {
			newGame = "拔河遊戲"
		} else if staff.Game == "bingo" {
			newGame = "賓果遊戲"
		} else if staff.Game == "3DGachaMachine" {
			newGame = "3D扭蛋機遊戲"
		} else if staff.Game == "vote" {
			newGame = "投票遊戲"
		}

		values := []interface{}{
			// 用戶資訊
			staff.Name, staff.Number,
			// 遊戲資訊
			newGame, staff.Title,
			// 黑名單資訊
			staff.Reason,
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

		// 設置存儲格的值
		for n, value := range values {
			file.SetCellValue(sheet, rowNames[n]+strconv.Itoa(i+2), value)
		}
	}

	// 處理excel檔案名稱
	if game != "" {
		gameCN = game
	} else if game == "" {
		gameCN = "所有遊戲類型"
	}

	if gameID != "" {
		gameIDCN = gameID
	} else if gameID == "" {
		gameIDCN = "所有遊戲場次"
	}

	fileName = "黑名單人員-" + activityID + "-" + gameCN + "-" + gameIDCN
	// 儲存EXCEL
	if err := file.SaveAs(fmt.Sprintf(config.STORE_PATH+"/excel/%s.xlsx", fileName)); err != nil {
		return file, err
	}

	return file, nil
}

// IsBlack 是否為黑名單(目前都是從redis裡判斷)
func (m BlackStaffModel) IsBlack(isRedis bool, activityID, gameID,
	game, userID string) (isBlack bool) {
	var (
		// channel,
		key string
		// blackStaffs  = make([]string, 0)
		people int64
	)
	// redis key名稱
	if game == "activity" {
		// 活動
		key = config.BLACK_STAFFS_ACTIVITY_REDIS + activityID
		// channel = config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + activityID
	} else if game == "message" {
		// 訊息
		key = config.BLACK_STAFFS_MESSAGE_REDIS + activityID
		// channel = config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + activityID
	} else if game == "question" {
		// 提問
		key = config.BLACK_STAFFS_QUESTION_REDIS + activityID
		// channel = config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + activityID
	} else if game == "signname" {
		// 簽名
		key = config.BLACK_STAFFS_SIGNNAME_REDIS + activityID
		// channel = config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + activityID
	} else if gameID != "" {
		// 遊戲
		key = config.BLACK_STAFFS_GAME_REDIS + gameID
		// channel = config.CHANNEL_BLACK_STAFFS_GAME_REDIS + gameID
	}

	// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
	people = m.RedisConn.SetCard(key)
	// fmt.Println("黑名單人數: ", people)

	// redis中沒有黑名單人員資訊
	if people == 0 {
		// log.Println("redis中目前沒有資料，執行findall")

		// 從redis取得資料，確定redis中有該場活動黑名單資料
		m.FindAll(true, activityID, gameID, game)
	}
	// else {
	// log.Println("redis中已有資料")
	// }

	// 判斷redis是否有用戶黑名單資訊
	isBlack = m.RedisConn.SetIsMember(key, userID)

	// 判斷是否為黑名單
	// for _, balckStaff := range blackStaffs {
	// 	// 黑名單
	// 	if balckStaff == userID {
	// 		isBlack = true
	// 		break
	// 	}
	// }

	// if isRedis {
	// }
	return
}

// Add 增加資料
func (m BlackStaffModel) Add(isRedis bool, model EditBlackStaffModel) error {
	var (
	// fields = []string{
	// 	"user_id",
	// 	"activity_id",
	// 	"game_id",
	// 	"game",
	// 	"reason",
	// }
	)

	// 先刪除原有的黑名單人員資料
	m.Table(m.Base.TableName).Where("activity_id", "=", model.ActivityID).
		Where("game_id", "=", model.GameID).Where("game", "=", model.Game).
		WhereIn("user_id", model.LineUsers).Delete()

	// 將 struct 轉換為 map[string]interface{} 格式
	// data := utils.StructToMap(model)

	// 將黑名單資訊寫入資料表中(多筆)
	for _, userID := range model.LineUsers {
		if _, err := m.Table(m.TableName).Insert(command.Value{
			"user_id":     userID,
			"activity_id": model.ActivityID,
			"game_id":     model.GameID,
			"game":        model.Game,
			"reason":      model.Reason,
		}); err != nil {
			return err
		}
	}

	if isRedis {
		var channel, key string
		// redis key名稱
		if model.Game == "activity" {
			// 活動
			key = config.BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID
		} else if model.Game == "message" {
			// 訊息
			key = config.BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID
		} else if model.Game == "question" {
			// 提問
			key = config.BLACK_STAFFS_QUESTION_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + model.ActivityID
		} else if model.Game == "signname" {
			// 簽名
			key = config.BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID
		} else if model.GameID != "" {
			// 遊戲
			key = config.BLACK_STAFFS_GAME_REDIS + model.GameID
			channel = config.CHANNEL_BLACK_STAFFS_GAME_REDIS + model.GameID
		}

		if model.Game == "activity" {
			// 更改活動黑名單資訊時，需要清除所有其他的黑名單redis資訊
			m.RedisConn.DelCache(config.BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_QUESTION_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS+model.ActivityID, "修改資料")

			// 清除所有遊戲的黑名單redis資訊，先查詢該活動下的所有遊戲
			games, err := DefaultGameModel().
				SetConn(m.DbConn, m.RedisConn, m.MongoConn).
				FindAll(model.ActivityID, "")
			if err != nil {
				return err
			}
			for _, game := range games {
				m.RedisConn.DelCache(config.BLACK_STAFFS_GAME_REDIS + game.GameID)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_GAME_REDIS+game.GameID, "修改資料")
			}
		} else {
			m.RedisConn.DelCache(key)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(channel, "修改資料")
		}

		// if model.Game == "question" {
		// 	m.RedisConn.DelCache(config.QUESTION_REDIS + model.ActivityID)                   // 提問資訊(時間)
		// 	m.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + model.ActivityID)     // 提問資訊(時間)
		// 	m.RedisConn.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + model.ActivityID)    // 提問資訊(讚數)
		// 	m.RedisConn.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + model.ActivityID) // 提問資訊(用戶按讚紀錄)
		// }
	}
	return nil
}

// Update 更新資料(移出黑名單人員資料)
func (m BlackStaffModel) Update(isRedis bool, model EditBlackStaffModel) error {
	// var (
	// sql = m.Table(m.Base.TableName).Where("activity_id", "=", model.ActivityID).
	// 	Where("game_id", "=", model.GameID).Where("game", "=", model.Game).
	// 	WhereIn("user_id", model.LineUsers)
	// err error
	// )

	// 刪除黑名單人員資料
	err := m.Table(m.Base.TableName).Where("activity_id", "=", model.ActivityID).
		Where("game_id", "=", model.GameID).Where("game", "=", model.Game).
		WhereIn("user_id", model.LineUsers).Delete()
	if err != nil && err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
		return errors.New("錯誤: 移出黑名單人員資料發生問題，請重新操作")
	}

	// if model.Reason != "" {
	// 	// 修改黑名單原因
	// 	err = sql.Update(command.Value{"reason": model.Reason})
	// } else {
	// }

	if isRedis {
		var channel, key string
		// redis key名稱
		if model.Game == "activity" {
			// 活動
			key = config.BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID
		} else if model.Game == "message" {
			// 訊息
			key = config.BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID
		} else if model.Game == "question" {
			// 提問
			key = config.BLACK_STAFFS_QUESTION_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS + model.ActivityID
		} else if model.Game == "signname" {
			// 簽名
			key = config.BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID
			channel = config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID
		} else if model.GameID != "" {
			// 遊戲
			key = config.BLACK_STAFFS_GAME_REDIS + model.GameID
			channel = config.CHANNEL_BLACK_STAFFS_GAME_REDIS + model.GameID
		}

		if model.Game == "activity" {
			// 更改活動黑名單資訊時，需要清除所有其他的黑名單redis資訊
			m.RedisConn.DelCache(config.BLACK_STAFFS_ACTIVITY_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_MESSAGE_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_QUESTION_REDIS + model.ActivityID)
			m.RedisConn.DelCache(config.BLACK_STAFFS_SIGNNAME_REDIS + model.ActivityID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_ACTIVITY_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_MESSAGE_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_QUESTION_REDIS+model.ActivityID, "修改資料")
			m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_SIGNNAME_REDIS+model.ActivityID, "修改資料")

			// 清除所有遊戲的黑名單redis資訊，先查詢該活動下的所有遊戲
			games, err := DefaultGameModel().
				SetConn(m.DbConn, m.RedisConn, m.MongoConn).
				FindAll(model.ActivityID, "")
			if err != nil {
				return err
			}
			for _, game := range games {
				m.RedisConn.DelCache(config.BLACK_STAFFS_GAME_REDIS + game.GameID)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				m.RedisConn.Publish(config.CHANNEL_BLACK_STAFFS_GAME_REDIS+game.GameID, "修改資料")
			}
		} else {
			m.RedisConn.DelCache(key)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(channel, "修改資料")
		}

		// if model.Game == "question" {
		// 	m.RedisConn.DelCache(config.QUESTION_REDIS + model.ActivityID)                   // 提問資訊(時間)
		// 	m.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + model.ActivityID)     // 提問資訊(時間)
		// 	m.RedisConn.DelCache(config.QUESTION_ORDER_BY_LIKES_REDIS + model.ActivityID)    // 提問資訊(讚數)
		// 	m.RedisConn.DelCache(config.QUESTION_USER_LIKE_RECORDS_REDIS + model.ActivityID) // 提問資訊(用戶按讚紀錄)
		// }
	}
	return nil
}

// MapToBlackStaffModel map轉換[]BlackStaffModel
func MapToBlackStaffModel(items []map[string]interface{}) []BlackStaffModel {
	var staffs = make([]BlackStaffModel, 0)
	for _, item := range items {
		var (
			staff BlackStaffModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &staff)

		staffs = append(staffs, staff)
	}
	return staffs
}

// staff.ID, _ = item["id"].(int64)
// staff.ActivityID, _ = item["activity_id"].(string)
// staff.GameID, _ = item["game_id"].(string)
// staff.Game, _ = item["game"].(string)
// staff.UserID, _ = item["user_id"].(string)
// staff.Reason, _ = item["reason"].(string)

// 用戶資訊
// staff.Name, _ = item["name"].(string)
// staff.Avatar, _ = item["avatar"].(string)
// staff.Phone, _ = item["phone"].(string)
// staff.Email, _ = item["email"].(string)
// staff.ExtEmail, _ = item["ext_email"].(string)

// 遊戲資訊
// staff.Title, _ = item["title"].(string)

// 自定義用戶資料
// staff.Ext1, _ = item["ext_1"].(string)
// staff.Ext2, _ = item["ext_2"].(string)
// staff.Ext3, _ = item["ext_3"].(string)
// staff.Ext4, _ = item["ext_4"].(string)
// staff.Ext5, _ = item["ext_5"].(string)
// staff.Ext6, _ = item["ext_6"].(string)
// staff.Ext7, _ = item["ext_7"].(string)
// staff.Ext8, _ = item["ext_8"].(string)
// staff.Ext9, _ = item["ext_9"].(string)
// staff.Ext10, _ = item["ext_10"].(string)
// staff.Number, _ = item["number"].(int64)

// // 自定義欄位資訊
// staff.Ext1Name, _ = item["ext_1_name"].(string)
// staff.Ext1Type, _ = item["ext_1_type"].(string)
// staff.Ext1Options, _ = item["ext_1_options"].(string)
// staff.Ext1Required, _ = item["ext_1_required"].(string)

// staff.Ext2Name, _ = item["ext_2_name"].(string)
// staff.Ext2Type, _ = item["ext_2_type"].(string)
// staff.Ext2Options, _ = item["ext_2_options"].(string)
// staff.Ext2Required, _ = item["ext_2_required"].(string)

// staff.Ext3Name, _ = item["ext_3_name"].(string)
// staff.Ext3Type, _ = item["ext_3_type"].(string)
// staff.Ext3Options, _ = item["ext_3_options"].(string)
// staff.Ext3Required, _ = item["ext_3_required"].(string)

// staff.Ext4Name, _ = item["ext_4_name"].(string)
// staff.Ext4Type, _ = item["ext_4_type"].(string)
// staff.Ext4Options, _ = item["ext_4_options"].(string)
// staff.Ext4Required, _ = item["ext_4_required"].(string)

// staff.Ext5Name, _ = item["ext_5_name"].(string)
// staff.Ext5Type, _ = item["ext_5_type"].(string)
// staff.Ext5Options, _ = item["ext_5_options"].(string)
// staff.Ext5Required, _ = item["ext_5_required"].(string)

// staff.Ext6Name, _ = item["ext_6_name"].(string)
// staff.Ext6Type, _ = item["ext_6_type"].(string)
// staff.Ext6Options, _ = item["ext_6_options"].(string)
// staff.Ext6Required, _ = item["ext_6_required"].(string)

// staff.Ext7Name, _ = item["ext_7_name"].(string)
// staff.Ext7Type, _ = item["ext_7_type"].(string)
// staff.Ext7Options, _ = item["ext_7_options"].(string)
// staff.Ext7Required, _ = item["ext_7_required"].(string)

// staff.Ext8Name, _ = item["ext_8_name"].(string)
// staff.Ext8Type, _ = item["ext_8_type"].(string)
// staff.Ext8Options, _ = item["ext_8_options"].(string)
// staff.Ext8Required, _ = item["ext_8_required"].(string)

// staff.Ext9Name, _ = item["ext_9_name"].(string)
// staff.Ext9Type, _ = item["ext_9_type"].(string)
// staff.Ext9Options, _ = item["ext_9_options"].(string)
// staff.Ext9Required, _ = item["ext_9_required"].(string)

// staff.Ext10Name, _ = item["ext_10_name"].(string)
// staff.Ext10Type, _ = item["ext_10_type"].(string)
// staff.Ext10Options, _ = item["ext_10_options"].(string)
// staff.Ext10Required, _ = item["ext_10_required"].(string)

// staff.ExtEmailRequired, _ = item["ext_email_required"].(string)
// staff.ExtPhoneRequired, _ = item["ext_phone_required"].(string)
// staff.InfoPicture, _ = item["info_picture"].(string)

// MapToBlackStaffModel map轉換[]BlackStaffModel
// func MapToBlackStaffModel(items []map[string]interface{}) []BlackStaffModel {
// 	var staffs = make([]BlackStaffModel, 0)
// 	for _, item := range items {
// 		var (
// 			staff BlackStaffModel
// 		)
// 		staff.ID, _ = item["id"].(int64)
// 		staff.UserID, _ = item["user_id"].(string)
// 		staff.ActivityID, _ = item["activity_id"].(string)
// 		staff.GameID, _ = item["game_id"].(string)
// 		staff.Game, _ = i
