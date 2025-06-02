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
	"strings"

	"github.com/xuri/excelize/v2"
)

// var (
// 	staffLock sync.Mutex
// )

// GameStaffModel 資料表欄位
type GameStaffModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	GameID     string `json:"game_id"`
	Round      int64  `json:"round"`
	Status     string `json:"status"`
	ApplyTime  string `json:"apply_time"`
	// Black      string `json:"-"`

	// 用戶資訊
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ExtEmail string `json:"ext_email"`

	// 遊戲資訊
	Game  string `json:"game" example:"game name"`
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

// NewGameStaffModel 資料表欄位
type NewGameStaffModel struct {
	UserID     string `json:"user_id" example:"user_id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	Game       string `json:"game" example:"game"`
	Team       string `json:"team" example:"team"`
	Round      string `json:"round" example:"1"`
	Status     string `json:"status" example:"success、fail"`
	// Black      string `json:"black" example:"yes、no"`
}

// EditGameStaffModel 資料表欄位
type EditGameStaffModel struct {
	ID     string `json:"id" example:"id"`
	GameID string `json:"game_id" example:"game_id"`
	UserID string `json:"user_id" example:"user_id"`
	Status string `json:"status" example:"success、fail"`
	// Black  string `json:"black" example:"yes、no"`
	Token string `json:"token" example:"token"`
}

// DefaultGameStaffModel 預設GameStaffModel
func DefaultGameStaffModel() GameStaffModel {
	return GameStaffModel{Base: Base{TableName: config.ACTIVITY_STAFF_GAME_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m GameStaffModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameStaffModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m GameStaffModel) SetDbConn(conn db.Connection) GameStaffModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m GameStaffModel) SetRedisConn(conn cache.Connection) GameStaffModel {
// 	m.RedisConn = conn
// 	return m
// }

// // SetMongoConn 設定connection
// func (m GameStaffModel) SetMongoConn(conn mongo.Connection) GameStaffModel {
// 	m.MongoConn = conn
// 	return m
// }

// IsUserExist 判斷用戶資料是否已存在遊戲資料中
func (m GameStaffModel) IsUserExist(gameID, userID string, round int64) bool {
	item, _ := m.Table(m.Base.TableName).
		Where("game_id", "=", gameID).
		Where("user_id", "=", userID).
		Where("round", "=", round).First()
	if item == nil {
		return false
	}
	return true
}

// FindAll 查詢所有遊戲人員資料
func (m GameStaffModel) FindAll(activityID, gameID, userID, game, round string,
	limit, offset int64) ([]GameStaffModel, error) {
	var (
		gameStaffsModel = make([]GameStaffModel, 0)
		sql             = m.Table(m.Base.TableName).
				Select("activity_staff_game.id", "activity_staff_game.activity_id",
				"activity_staff_game.user_id", "activity_staff_game.game_id",
				"activity_staff_game.round", "activity_staff_game.status",
				"activity_staff_game.apply_time",
				"activity_staff_game.game",
				// "activity_staff_game.black",

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
				// "activity_game.game", "activity_game.title",

				// 用戶
				"line_users.name", "line_users.avatar", "line_users.phone",
				"line_users.email", "line_users.ext_email").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_game.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_staff_game.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_game.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_game.activity_id",
				FieldA1:   "activity_applysign.activity_id",
				FieldB:    "activity_staff_game.user_id",
				FieldB1:   "activity_applysign.user_id",
				Table:     "activity_applysign",
				Operation: "="}).
			OrderBy(
				"activity_staff_game.game", "asc",
				"activity_staff_game.game_id", "asc",
				"activity_staff_game.round", "asc",
				"activity_staff_game.status", "asc",
			).
			Where("activity_staff_game.activity_id", "=", activityID)
	)

	// 判斷參數是否為空
	if gameID != "" {
		sql = sql.WhereIn("activity_staff_game.game_id", interfaces(strings.Split(gameID, ",")))
	}
	if userID != "" {
		sql = sql.WhereIn("activity_staff_game.user_id", interfaces(strings.Split(userID, ",")))
	}
	if game != "" {
		sql = sql.WhereIn("activity_staff_game.game", interfaces(strings.Split(game, ",")))
	}
	if round != "" {
		sql = sql.WhereIn("activity_staff_game.round", interfaces(strings.Split(round, ",")))
	}
	if limit != 0 {
		sql = sql.Limit(limit)
	}
	if offset != 0 {
		sql = sql.Offset(offset)
	}

	items, err := sql.All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得遊戲人員資訊，請重新查詢")
	}

	gameStaffsModel = MapToGameStaffModel(items)

	// 多個遊戲場次資料(mongo)
	gamesModel, err := DefaultGameModel().
		SetConn(m.DbConn, m.RedisConn, m.MongoConn).
		FindAll(activityID, game)
	if err != nil {
		return gameStaffsModel, err
	}

	// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
	dataMap := make(map[string]GameModel)
	for _, game := range gamesModel {
		dataMap[game.GameID] = game
	}

	for i, item := range gameStaffsModel {
		if dataItem, found := dataMap[item.GameID]; found {

			// 合併欄位（可依需要擴充）
			item.Title = dataItem.Title

			gameStaffsModel[i] = item
		}
	}

	return gameStaffsModel, nil
}

// DeleteGame 刪除資料(該遊戲場次所有人員資料)
func (m GameStaffModel) DeleteGame(gameID string, round int64) error {
	return m.Table(m.Base.TableName).Where("game_id", "=", gameID).
		Where("round", "=", round).Delete()
}

// DeleteUser 刪除資料(該場次該用戶人員資料)
// func (m GameStaffModel) DeleteUser(gameID, userID string, round int64) error {
// 	var sql = m.Table(m.Base.TableName).
// 		Where("game_id", "=", gameID).
// 		Where("user_id", "=", userID).
// 		Where("round", "=", round)

// 	return sql.Delete()
// }

// Adds 新增多筆遊戲人員
func (m GameStaffModel) Adds(dataAmount int, activityID, gameID, game string, userIDs []string, round int) error {
	// 判斷資料數
	if dataAmount > 0 {
		var (
			activityIDs = make([]string, dataAmount) // 活動ID
			gameIDs     = make([]string, dataAmount) // 活動ID
			rounds      = make([]string, dataAmount) // 輪次
			games       = make([]string, dataAmount) // 遊戲類型
			roundStr    = strconv.Itoa(round)        // 輪次(string)
		)

		// 處理活動.遊戲ID.輪次陣列資料
		for i := range dataAmount {
			activityIDs[i] = activityID
			gameIDs[i] = gameID
			rounds[i] = roundStr
			games[i] = game
		}

		// 將資料匹量寫入activity_staff_game表中
		err := m.Table(m.TableName).BatchInsert(
			dataAmount, "activity_id,game_id,user_id,round,game",
			[][]string{activityIDs, gameIDs, userIDs, rounds, games})
		if err != nil {
			return errors.New("錯誤: 匹量新增遊戲人員發生問題(activity_staff_game)，請重新操作")
		}
	}

	return nil
}

// Add 增加資料
// func (m GameStaffModel) Add(model NewGameStaffModel) error {
// 	// staffLock.Lock()
// 	// defer staffLock.Unlock()

// 	if _, err := strconv.Atoi(model.Round); err != nil {
// 		return errors.New("錯誤: 遊戲輪次資料發生問題，請輸入有效的輪次資料")
// 	}
// 	if model.Status != "success" && model.Status != "fail" {
// 		return errors.New("錯誤: 狀態資料發生問題，請輸入有效的狀態資料")
// 	}
// 	// if model.Black != "yes" && model.Black != "no" {
// 	// 	return errors.New("錯誤: 黑名單資料發生問題，請輸入有效的黑名單資料")
// 	// }

// 	// 判斷遊戲人數是否額滿並新增遊戲人員
// 	// if model.Game == "tugofwar" {
// 	// 	// 拔河遊戲，遞增隊伍人數資料(left_team_game_attend、right_team_game_attend)
// 	// 	// if err := DefaultGameModel().SetDbConn(m.DbConn).SetRedisConn(m.RedisConn).
// 	// 	// 	UpdateTeamPeople(true, model.GameID, model.Team, 1, 1); err != nil {
// 	// 	// 	// log.Println("遊戲報名人數額滿")
// 	// 	// 	return err
// 	// 	// }

// 	// 	// 將玩家寫入隊伍陣列中(redis)
// 	// 	if err := DefaultGameModel().SetDbConn(m.DbConn).SetRedisConn(m.RedisConn).
// 	// 		AddPlayer(true, model.GameID, model.UserID, model.Team); err != nil {
// 	// 		return err
// 	// 	}
// 	// } else {
// 	// 	// 其他遊戲，遞增game_attend資料
// 	// 	if err := DefaultGameModel().SetDbConn(m.DbConn).SetRedisConn(m.RedisConn).
// 	// 		IncrAttend(true, model.GameID); err != nil {
// 	// 		// log.Println("遊戲報名人數額滿")
// 	// 		return err
// 	// 	}
// 	// }

// 	// 新增遊戲人員
// 	_, err := m.Table(m.TableName).Insert(command.Value{
// 		"user_id":     model.UserID,
// 		"activity_id": model.ActivityID,
// 		"game_id":     model.GameID,
// 		"round":       model.Round,
// 		"status":      model.Status,
// 		"game": model.Game,
// 		// "black":       model.Black,
// 	})
// 	return err
// }

// Export 匯出遊戲人員資料
func (m GameStaffModel) Export(activityID, gameID, game, userID, round string,
	limit, offset int64) (*excelize.File, error) {
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
			"遊戲名稱", "遊戲類型",
			// 參加資訊
			"遊戲輪次", "參加時間",
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
		gameStaffModel, err2 = DefaultGameStaffModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					FindAll(activityID, gameID, userID, game, round, limit, offset)
		sheet                                         = "Sheet1"
		fileName, gameCN, gameIDCN, userIDCN, roundCN string
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
	for i, staff := range gameStaffModel {
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
		} else if staff.Game == "tugofwar" {
			newGame = "拔河遊戲"
		} else if staff.Game == "bingo" {
			newGame = "賓果遊戲"
		} else if staff.Game == "3DGachaMachine" {
			newGame = "3D扭蛋機遊戲"
		}

		values := []interface{}{
			// 用戶資訊
			staff.Name, staff.Number,
			// 遊戲資訊
			staff.Title, newGame,
			// 參加資訊
			staff.Round, staff.ApplyTime,
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

	if userID != "" {
		userIDCN = userID
	} else if userID == "" {
		userIDCN = "所有用戶"
	}

	// 判斷是否取得特定輪次資料
	if round != "" {
		roundCN = fmt.Sprintf("第%s輪", round)
	} else if round == "" {
		roundCN = "所有輪次"
	}

	fileName = "遊戲人員-" + activityID + "-" + gameCN + "-" +
		gameIDCN + "-" + userIDCN + "-" + roundCN
	// 儲存EXCEL
	if err := file.SaveAs(fmt.Sprintf(config.STORE_PATH+"/excel/%s.xlsx", fileName)); err != nil {
		return file, err
	}

	return file, nil
}

// MapToGameStaffModel map轉換[]GameStaffModel
func MapToGameStaffModel(items []map[string]interface{}) []GameStaffModel {
	var staffs = make([]GameStaffModel, 0)
	for _, item := range items {
		var (
			staff GameStaffModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &staff)

		staffs = append(staffs, staff)
	}
	return staffs
}

// staff.ID, _ = item["id"].(int64)
// staff.UserID, _ = item["user_id"].(string)
// staff.ActivityID, _ = item["activity_id"].(string)
// staff.GameID, _ = item["game_id"].(string)
// staff.Round, _ = item["round"].(int64)
// staff.Status, _ = item["status"].(string)
// staff.ApplyTime, _ = item["apply_time"].(string)

// // 用戶資訊
// staff.Name, _ = item["name"].(string)
// staff.Avatar, _ = item["avatar"].(string)
// staff.Phone, _ = item["phone"].(string)
// staff.Email, _ = item["email"].(string)
// staff.ExtEmail, _ = item["ext_email"].(string)

// // 遊戲資訊
// staff.Title, _ = item["title"].(string)
// staff.Game, _ = item["game"].(string)

// // 自定義用戶資料
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

// IsBlack 是否為黑名單
// func (m GameStaffModel) IsBlack(isRedis bool, gameID, userID string) (isBlack bool) {
// var count int64
// if isRedis {
// 	// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
// 	count = m.RedisConn.SetCard(config.BLACK_STAFFS_REDIS + gameID)
// }

// // redis中沒有該場次黑名單人員資訊，查詢資料表並加入redis中
// if count == 0 {
// 	items, _ := m.Table(m.Base.TableName).Where("game_id", "=", gameID).
// 		Where("black", "=", "yes").All()

// 	if isRedis && len(items) > 0 {
// 		var staffs = []interface{}{config.BLACK_STAFFS_REDIS + gameID}
// 		for _, item := range items {
// 			// id, _ := item["user_id"].(string)
// 			staffs = append(staffs, item["user_id"])
// 		}
// 		// 將黑名單資料設置至redis中
// 		m.RedisConn.SetAdd(staffs)
// 		m.RedisConn.SetExpire(config.BLACK_STAFFS_REDIS+gameID,
// 			config.REDIS_EXPIRE)
// 	}
// }

// if isRedis {
// 	// 判斷redis是否有用戶黑名單資訊
// 	isBlack = m.RedisConn.SetIsMember(config.BLACK_STAFFS_REDIS+gameID, userID)
// }
// return
// }
// Update 更新資料
// func (m GameStaffModel) Update(isRedis bool, model EditGameStaffModel) error {
// 	var (
// 		fieldValues = command.Value{}
// 		// fields      = []string{"status", "black"}
// 		// values      = []string{model.Status, model.Black}
// 		fields = []string{"status"}
// 		values = []string{model.Status}
// 	)

// 	if model.Status != "" {
// 		if model.Status != "success" && model.Status != "fail" {
// 			return errors.New("錯誤: 狀態資料發生問題，請輸入有效的狀態資料")
// 		}
// 	}
// 	// if model.Black != "" {
// 	// 	if model.Black != "yes" && model.Black != "no" {
// 	// 		return errors.New("錯誤: 黑名單資料發生問題，請輸入有效的黑名單資料")
// 	// 	}
// 	// }

// 	for i, value := range values {
// 		if value != "" {
// 			fieldValues[fields[i]] = value
// 		}
// 	}
// 	if len(fieldValues) == 0 {
// 		return nil
// 	}

// 	if err := m.Table(m.Base.TableName).
// 		Where("id", "=", model.ID).Update(fieldValues); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新遊戲人員資料發生問題，請重新操作")
// 	}

// 	if isRedis {
// 		m.RedisConn.DelCache(config.BLACK_S
