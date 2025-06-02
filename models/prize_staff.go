package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"strconv"
	"strings"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"

	"github.com/xuri/excelize/v2"
)

// PrizeStaffModel 資料表欄位
type PrizeStaffModel struct {
	Base           `json:"-"`
	ID             int64  `json:"id" example:"1"`
	ActivityID     string `json:"activity_id" example:"activity_id"`
	ActivityUserID string `json:"activity_user_id" example:"activity_id"`
	GameID         string `json:"game_id" example:"game_id"`
	PrizeID        string `json:"prize_id" example:"prize_id"`
	UserID         string `json:"user_id" example:"user_id"`
	Round          int64  `json:"round" example:"1"`
	WinTime        string `json:"win_time" example:"2021-01-01 00:00"`
	Status         string `json:"status" example:"yes、no"`
	// White      string `json:"white" example:"yes、no"`

	// 分數
	Score  int64   `json:"score" example:"1"`
	Score2 float64 `json:"score_2" example:"1"`
	// 名次
	Rank int64 `json:"rank" example:"1"`

	// 隊伍資訊
	Team   string `json:"team"`   // 隊伍
	Leader string `json:"leader"` // 隊長
	Mvp    string `json:"mvp"`    // mvp

	// 用戶資訊
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ExtEmail string `json:"ext_email"`

	// 遊戲資訊
	Title         string `json:"title" example:"Game Title"`
	Game          string `json:"game" example:"redpack、ropepack、whack_mole、lottery"`
	LeftTeamName  string `json:"left_team_name" example:"name"`  // 左邊隊伍名稱
	RightTeamName string `json:"right_team_name" example:"name"` // 右邊隊伍名稱

	// 用戶號碼
	Number int64 `json:"number" example:"1"`

	// 獎品資訊
	PrizeName     string `json:"prize_name" example:"Prize Name"`
	PrizeType     string `json:"prize_type" example:"money、prize、thanks、first、second、third、special"`
	PrizePicture  string `json:"prize_picture" example:"https://..."`
	PrizePrice    int64  `json:"prize_price" example:"100"`
	PrizeMethod   string `json:"prize_method" example:"site、mail、thanks"`
	PrizePassword string `json:"prize_password" example:"password"`

	// 用戶自定義資料
	Ext1  string `json:"ext_1"`
	Ext2  string `json:"ext_2"`
	Ext3  string `json:"ext_3"`
	Ext4  string `json:"ext_4"`
	Ext5  string `json:"ext_5"`
	Ext6  string `json:"ext_6"`
	Ext7  string `json:"ext_7"`
	Ext8  string `json:"ext_8"`
	Ext9  string `json:"ext_9"`
	Ext10 string `json:"ext_10"`

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

// NewPrizeStaffModel 資料表欄位
type NewPrizeStaffModel struct {
	UserID     string `json:"user_id" example:"user_id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	Game       string `json:"game" example:"game"`
	GameID     string `json:"game_id" example:"game_id"`
	PrizeID    string `json:"prize_id" example:"prize_id"`
	Round      string `json:"round" example:"1"`
	Status     string `json:"status" example:"yes、no"`
	// 分數
	Score  int64   `json:"score" example:"1"`
	Score2 float64 `json:"score_2" example:"1"`
	// 名次
	Rank int64 `json:"rank" example:"1"`
	// White      string `json:"white" example:"yes、no"`

	// 隊伍資訊
	Team   string `json:"team"`   // 隊伍
	Leader string `json:"leader"` // 隊長
	Mvp    string `json:"mvp"`    // mvp
}

// EditPrizeStaffModel 資料表欄位
type EditPrizeStaffModel struct {
	ID       []interface{} `json:"id" example:"id"`
	Role     string        `json:"role" example:"admin、guest"`
	Password string        `json:"password" example:"password"`
	Status   string        `json:"status" example:"yes、no"`
	// White    string `json:"White" example:"yes、no"`
}

// DefaultPrizeStaffModel 預設PrizeStaffModel
func DefaultPrizeStaffModel() PrizeStaffModel {
	return PrizeStaffModel{Base: Base{TableName: config.ACTIVITY_STAFF_PRIZE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m PrizeStaffModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) PrizeStaffModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m PrizeStaffModel) SetDbConn(conn db.Connection) PrizeStaffModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m PrizeStaffModel) SetRedisConn(conn cache.Connection) PrizeStaffModel {
// 	m.RedisConn = conn
// 	return m
// }

// // SetMongoConn 設定connection
// func (m PrizeStaffModel) SetMongoConn(conn mongo.Connection) PrizeStaffModel {
// 	m.MongoConn = conn
// 	return m
// }

// FindAll 查詢所有中獎紀錄
func (m PrizeStaffModel) FindAll(activityID, gameID, userID, game, round, status, teamType string,
	limit, offset int64) ([]PrizeStaffModel, error) {
	var (
		prizeStaffsModel = make([]PrizeStaffModel, 0)
		sql              = m.Table(m.Base.TableName).
					Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
				"activity_staff_prize.user_id", "activity_staff_prize.round",
				"activity_staff_prize.win_time", "activity_staff_prize.status",
				"activity_staff_prize.score",
				"activity_staff_prize.score_2",
				"activity_staff_prize.rank",
				"activity_staff_prize.team",
				"activity_staff_prize.leader",
				"activity_staff_prize.mvp",
				"activity_staff_prize.game",

				// 活動資訊
				"activity.user_id as activity_user_id",

				// 遊戲場次
				// "activity_game.title", "activity_game.game",
				// 隊伍資訊
				// "activity_game.left_team_name", "activity_game.right_team_name",

				// 獎品資訊
				"activity_prize.prize_name", "activity_prize.prize_type",
				"activity_prize.prize_picture", "activity_prize.prize_price",
				"activity_prize.prize_method", "activity_prize.prize_password",
				"activity_prize.team_type",

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
				"activity_customize.ext_email_required", "activity_customize.ext_phone_required",
				"activity_customize.info_picture",
				"activity_applysign.number").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_staff_prize.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.activity_id",
				FieldA1:   "activity_applysign.activity_id",
				FieldB:    "activity_staff_prize.user_id",
				FieldB1:   "activity_applysign.user_id",
				Table:     "activity_applysign",
				Operation: "="}).
			Where("activity_staff_prize.prize_id", "!=", "").
			Where("activity_prize.prize_type", "!=", "thanks").
			Where("activity_prize.prize_type", "!=", "").
			Where("activity_prize.prize_method", "!=", "thanks").
			Where("activity_prize.prize_method", "!=", "").
			Where("activity_staff_prize.activity_id", "=", activityID).
			OrderBy(
				"activity_staff_prize.game", "asc",
				"activity_staff_prize.game_id", "asc",
				"activity_staff_prize.round", "asc",
				"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
				"activity_staff_prize.rank", "asc",
				"activity_staff_prize.score", "desc",
				"activity_staff_prize.prize_id", "asc",
			)
	)

	// 判斷參數是否為空
	if gameID != "" {
		sql = sql.WhereIn("activity_staff_prize.game_id", interfaces(strings.Split(gameID, ",")))
	}
	if userID != "" {
		sql = sql.WhereIn("activity_staff_prize.user_id", interfaces(strings.Split(userID, ",")))
	}
	if game != "" {
		sql = sql.WhereIn("activity_staff_prize.game", interfaces(strings.Split(game, ",")))
	}
	if round != "" {
		sql = sql.WhereIn("activity_staff_prize.round", interfaces(strings.Split(round, ",")))
	}
	if status != "" {
		sql = sql.WhereIn("activity_staff_prize.status", interfaces(strings.Split(status, ",")))
	}
	if teamType != "" {
		sql = sql.Where("activity_prize.team_type", "=", teamType)
	}
	if limit != 0 {
		sql = sql.Limit(limit)
	}
	if offset != 0 {
		sql = sql.Offset(offset)
	}

	items, err := sql.All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
	}

	prizeStaffsModel = MapToPrizeStaffModel(items)

	// 多個遊戲場次資料
	gamesModel, err := DefaultGameModel().
		SetConn(m.DbConn, m.RedisConn, m.MongoConn).
		FindAll(activityID, game)
	if err != nil {
		return prizeStaffsModel, err
	}

	// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
	dataMap := make(map[string]GameModel)
	for _, game := range gamesModel {
		dataMap[game.GameID] = game
	}

	for i, item := range prizeStaffsModel {
		if dataItem, found := dataMap[item.GameID]; found {

			// 合併欄位（可依需要擴充）
			item.Title = dataItem.Title

			prizeStaffsModel[i] = item
		}
	}

	return prizeStaffsModel, nil
}

// FindUserGameRecords 查詢用戶該活動的中獎及未中獎紀錄資料
// 目前都是使用redis處理
// 玩家遊戲紀錄redis中不一定有所有搖號抽獎中獎紀錄，因為抽獎時是將中獎紀錄直接寫入該活動所有搖號抽獎場次中獎人員redis中(draw_numbers_winning_staffs_ 活動ID)
func (m PrizeStaffModel) FindUserGameRecords(isRedis bool, activityID, userID string) ([]PrizeStaffModel, error) {
	// log.Println("執行查詢用戶該活動的中獎及未中獎紀錄資料")

	var (
		records = make([]PrizeStaffModel, 0)
		sql     = m.Table(m.Base.TableName).
			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
				"activity_staff_prize.user_id", "activity_staff_prize.round",
				"activity_staff_prize.win_time", "activity_staff_prize.status",
				"activity_staff_prize.score",
				"activity_staff_prize.score_2",
				"activity_staff_prize.rank",
				"activity_staff_prize.team",
				"activity_staff_prize.leader",
				"activity_staff_prize.mvp",

				// 活動資訊
				"activity.user_id as activity_user_id",

				// 獎品
				"activity_prize.game",
				"activity_prize.prize_name", "activity_prize.prize_type",
				"activity_prize.prize_picture", "activity_prize.prize_price",
				"activity_prize.prize_method", "activity_prize.prize_password",

				// 用戶
				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "="}).
			Where("activity_staff_prize.activity_id", "=", activityID).
			Where("activity_staff_prize.user_id", "=", userID)
	)

	// 從redis取得玩家投票遊戲紀錄資料(hash類型)
	if isRedis && userID != "" {
		// redis處理
		// 判斷redis裡是否有玩家遊戲紀錄資訊，有則不執行查詢資料表功能
		recordsJson, err := m.RedisConn.HashGetCache(config.USER_GAME_RECORDS_REDIS+activityID, userID) // 包含test測試資料，避免重複執行資料庫
		if err != nil || recordsJson == "" {
			// redis查詢不到玩家遊戲紀錄資料
			// 從資料庫讀取得玩家遊戲紀錄資料
			// log.Println("redis無法取得玩家遊戲紀錄，查詢資料表")

			items, err := sql.All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得玩家遊戲紀錄資訊，請重新查詢")
			}
			records = MapToPrizeStaffModel(items) // 多筆資料

			// 將玩家遊戲紀錄資料添加至redis中(hash類型)
			if len(records) > 0 {
				// log.Println("資料表有資料，加入玩家遊戲紀錄至redis")

				err = m.RedisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
					userID, utils.JSON(records))
				if err != nil {
					return nil, errors.New("錯誤: 將玩家遊戲紀錄寫入redis中發生問題")
				}

				// m.RedisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
				// 	config.REDIS_EXPIRE)
			} else if len(records) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 將測試資料寫入redis，避免重複查詢資料庫
				err = m.RedisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID, userID,
					utils.JSON([]PrizeStaffModel{{
						UserID: "test",
					}}))
				if err != nil {
					return nil, errors.New("錯誤: 將人員投票紀錄寫入redis中發生問題")
				}

				// m.RedisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
				// 	config.REDIS_EXPIRE)
			}
		} else {
			// log.Println("redis有資料，過濾資料並回傳")

			var recordModels = make([]PrizeStaffModel, 0)
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

	// if len(records) == 0 {
	// 	items, err := sql.All()
	// 	if err != nil {
	// 		return nil, errors.New("錯誤: 無法取得中獎人員資訊，請重新查詢")
	// 	}
	// 	records = MapToPrizeStaffModel(items)
	// 	// fmt.Println("models package FindUserGameRecords function資料表查詢，長度: ", len(records))

	// 	// 將用戶遊戲紀錄資料添加至redis中(hash類型)
	// 	if isRedis && len(records) > 0 {
	// 		m.RedisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
	// 			userID, utils.JSON(records))
	// 		m.RedisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
	// 			config.REDIS_EXPIRE)
	// 	}
	// }
	return records, nil
}

// FindDrawNumbersAllWinningStaffs 查詢該活動下所有搖號抽獎場次中獎人員資料(join activity_prize)
// 回傳中獎人員陣列名單
// 搖號抽獎用(redis: DRAW_NUMBERS_WINNING_STAFFS_REDIS)
func (m PrizeStaffModel) FindDrawNumbersAllWinningStaffs(isRedis bool, activityID string) ([]string, error) {
	// log.Println("執行FindDrawNumbersAllWinningStaffs")

	var (
		sql = m.Table(m.Base.TableName).
			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
				"activity_staff_prize.user_id", "activity_staff_prize.round",
				"activity_staff_prize.win_time", "activity_staff_prize.status",
				"activity_staff_prize.score",
				"activity_staff_prize.score_2",
				"activity_staff_prize.rank",
				"activity_staff_prize.team",
				"activity_staff_prize.leader",
				"activity_staff_prize.mvp",

				// 獎品
				"activity_prize.game",
				"activity_prize.prize_name", "activity_prize.prize_type",
				"activity_prize.prize_picture", "activity_prize.prize_price",
				"activity_prize.prize_method", "activity_prize.prize_password").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "="}).
			Where("activity_staff_prize.activity_id", "=", activityID).
			Where("activity_staff_prize.prize_id", "!=", "").
			Where("activity_prize.prize_type", "!=", "thanks").
			Where("activity_prize.prize_type", "!=", "").
			Where("activity_prize.prize_method", "!=", "thanks").
			Where("activity_prize.prize_method", "!=", "").
			Where("activity_prize.game", "=", "draw_numbers")
		items  = make([]map[string]interface{}, 0)
		staffs = make([]string, 0)
		// err    error
	)

	// 從redis取得該活動下所有搖號抽獎場次中獎人員資料(set類型)
	if isRedis {
		// 判斷redis裡是否有中獎人員資訊，有則不執行查詢資料表功能
		redisStaffs, err := m.RedisConn.SetGetMembers(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)
		if err != nil {
			return nil, errors.New("錯誤: 無法該活動下所有搖號抽獎場次中獎人員資料(redis)，請重新查詢")
		}

		// log.Println("查詢redis中中獎紀錄(包含測試資料): ", len(redisStaffs))

		// 處理中獎資訊
		for _, staff := range redisStaffs {
			// 過濾測試資料
			if staff != "test" {
				staffs = append(staffs, staff)
			}
		}

		// log.Println("查詢redis後的中獎紀錄(不包含測試資料): ", len(staffs))

		if len(redisStaffs) == 0 {
			// log.Println("redis裡連測試資料都沒有，查詢資料表: ")

			items, err = sql.All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
			}

			var params = []interface{}{config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID} // redis參數

			if len(items) > 0 {
				for _, item := range items {
					userID, _ := item["user_id"].(string)

					// 判斷陣列裡是否已經有用戶資訊，如果沒有則加入陣列中
					if !utils.InArray(staffs, userID) {
						params = append(params, userID)
						staffs = append(staffs, userID)
					}
				}
			} else if len(items) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				params = append(params, "test")
			}

			// 將中獎人員資料設置至redis中
			m.RedisConn.SetAdd(params)

		}
	}

	return staffs, nil
}

// FindDrawNumbersWinningStaffs 查詢遊戲中獎人員資訊(join activity_prize)
// 先判斷活動是否允許重複中獎、相同類型遊戲是否允許重複中獎、同場次遊戲是否允許重複中獎
// 回傳中獎人員陣列名單
// 搖號抽獎用(redis: WINNING_STAFFS_REDIS)
func (m PrizeStaffModel) FindDrawNumbersWinningStaffs(isRedis bool, activityID, gameID, game string,
	activityAllow, gameAllow, allow bool) ([]string, error) {
	// log.Println("執行FindDrawNumbersWinningStaffs")
	var (
		sql = m.Table(m.Base.TableName).
			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
				"activity_staff_prize.user_id", "activity_staff_prize.round",
				"activity_staff_prize.win_time", "activity_staff_prize.status",
				"activity_staff_prize.score",
				"activity_staff_prize.score_2",
				"activity_staff_prize.rank",
				"activity_staff_prize.team",
				"activity_staff_prize.leader",
				"activity_staff_prize.mvp",

				// 獎品
				"activity_prize.game",
				"activity_prize.prize_name", "activity_prize.prize_type",
				"activity_prize.prize_picture", "activity_prize.prize_price",
				"activity_prize.prize_method", "activity_prize.prize_password").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "="}).
			Where("activity_staff_prize.activity_id", "=", activityID).
			Where("activity_staff_prize.prize_id", "!=", "").
			Where("activity_prize.prize_type", "!=", "thanks").
			Where("activity_prize.prize_type", "!=", "").
			Where("activity_prize.prize_method", "!=", "thanks").
			Where("activity_prize.prize_method", "!=", "")
		items  = make([]map[string]interface{}, 0)
		staffs = make([]string, 0)
		// err    error
	)

	// 從redis取得中獎人員資料(set類型)
	if isRedis {
		// 判斷redis裡是否有中獎資訊，有則不執行查詢資料表功能
		redisStaffs, err := m.RedisConn.SetGetMembers(config.WINNING_STAFFS_REDIS + gameID)
		if err != nil {
			return nil, errors.New("錯誤: 無法取得中獎人員紀錄(redis)，請重新查詢")
		}

		// log.Println("查詢redis中中獎紀錄(包含測試資料): ", len(redisStaffs))

		// 處理中獎資訊
		for _, staff := range redisStaffs {
			// 過濾測試資料
			if staff != "test" {
				staffs = append(staffs, staff)
			}
		}

		// log.Println("查詢redis後的中獎紀錄(不包含測試資料): ", len(staffs))

		if len(redisStaffs) == 0 {
			// log.Println("redis裡連測試資料都沒有，查詢資料表: ")

			// redis中無中獎人員資料，查詢資料表
			if !activityAllow {
				// 活動不允許重複中獎
			} else if !gameAllow {
				// 相同類型遊戲不允許重複中獎
				sql = sql.Where("activity_prize.game", "=", game)
			} else if !allow {
				// 同場次遊戲不允許重複中獎
				sql = sql.Where("activity_staff_prize.game_id", "=", gameID)
			} else {
				// 允許重複中獎
			}

			items, err = sql.All()
			if err != nil {
				return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
			}

			var params = []interface{}{config.WINNING_STAFFS_REDIS + gameID} // redis參數

			if len(items) > 0 {
				for _, item := range items {
					userID, _ := item["user_id"].(string)

					// 判斷陣列裡是否已經有用戶資訊，如果沒有則加入陣列中
					if !utils.InArray(staffs, userID) {
						params = append(params, userID)
						staffs = append(staffs, userID)
					}
				}
			} else if len(items) == 0 {
				// log.Println("資料表一樣沒資料，加入測試資料")

				// 資料庫目前無資料，加入test資料，避免頻繁查詢資料庫
				// 將測試資料寫入redis，避免重複查詢資料庫
				params = append(params, "test")
			}

			// 將中獎人員資料設置至redis中
			m.RedisConn.SetAdd(params)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
		}
	}

	// 直接查詢資料庫
	// if !isRedis {
	// 	log.Println("查詢資料表")

	// 	// redis中無中獎人員資料，查詢資料表
	// 	if !activityAllow {
	// 		// 活動不允許重複中獎
	// 	} else if !gameAllow {
	// 		// 相同類型遊戲不允許重複中獎
	// 		sql = sql.Where("activity_prize.game", "=", game)
	// 	} else if !allow {
	// 		// 同場次遊戲不允許重複中獎
	// 		sql = sql.Where("activity_staff_prize.game_id", "=", gameID)
	// 	} else {
	// 		// 允許重複中獎
	// 	}

	// 	items, err = sql.All()
	// 	if err != nil {
	// 		return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
	// 	}

	// 	// 將中獎人員資料添加至redis中(list類型)
	// 	// if isRedis && len(items) > 0 {
	// 	// 	var params = []interface{}{config.WINNING_STAFFS_REDIS + gameID} // redis參數
	// 	// 	for _, item := range items {
	// 	// 		userID, _ := item["user_id"].(string)

	// 	// 		// 判斷陣列裡是否已經有用戶資訊，如果沒有則加入陣列中
	// 	// 		if !utils.InArray(staffs, userID) {
	// 	// 			params = append(params, userID)
	// 	// 			staffs = append(staffs, userID)
	// 	// 		}
	// 	// 	}

	// 	// 	// 將中獎人員資料設置至redis中
	// 	// 	m.RedisConn.SetAdd(params)
	// 	// 	// m.RedisConn.SetExpire(config.WINNING_STAFFS_REDIS+gameID,
	// 	// 	// 	config.REDIS_EXPIRE)

	// 	// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	// 	// 	m.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
	// 	// }
	// }
	return staffs, nil
}

// FindRedpackAndLotteryWinningStaffs 查詢該遊戲場次的中獎紀錄
// 紅包遊戲、遊戲抽獎用(需照順序推進陣列中)
func (m PrizeStaffModel) FindRedpackAndLotteryWinningStaffs(isRedis bool, redisCount int64, gameID string,
	round, limit int64) ([]PrizeStaffModel, error) {
	// log.Println("測試FindRedpackAndLotteryWinningStaffs: ", redisCount)

	var (
		records = make([]PrizeStaffModel, 0)
		err     error
	)
	// 從redis取得中獎人員資料(list類型)
	if isRedis {
		var recordsJson []string // 中獎人員json資料(獎品類型排列，1頭獎.2二獎...依序排列)
		recordsJson, err = m.RedisConn.ListRange(config.WINNING_STAFFS_REDIS+gameID, 0, limit)
		if err != nil {
			return records, errors.New("錯誤: 從redis中取得中獎人員資訊發生問題")
		}

		if len(recordsJson) > 0 {
			for _, record := range recordsJson {
				var staff PrizeStaffModel
				// 解碼
				json.Unmarshal([]byte(record), &staff)

				records = append(records, staff)
			}
		}
		// fmt.Println("records: ", records)
	}
	// 舊版: 依照獎品類型分類至redis中(zset類型)
	// if isRedis {
	// 	var recordsJson []string // 中獎人員json資料(獎品類型排列，1頭獎.2二獎...依序排列)
	// 	recordsJson, err = m.RedisConn.ZSetRange(config.WINNING_STAFFS_REDIS+gameID, 0, limit)
	// 	if err != nil {
	// 		return records, errors.New("錯誤: 從redis中取得中獎人員資訊發生問題")
	// 	}

	// 	if len(recordsJson) > 0 {
	// 		for _, record := range recordsJson {
	// 			var staff PrizeStaffModel
	// 			// 解碼
	// 			json.Unmarshal([]byte(record), &staff)

	// 			records = append(records, staff)
	// 		}
	// 	}
	// }

	// fmt.Println("models package FindWinningRecordsByGameID function使用到資料庫: ", !isRedis, isRedis && len(records) == 0, redisCount == 0)

	// 如果第一次使用redis並且查詢不到資料, 查詢資料表
	if !isRedis || (isRedis && len(records) == 0 && redisCount == 0) {

		// if len(records) == 0 {

		// redis中無中獎人員資料，查詢資料表
		sql := m.Table(m.Base.TableName).
			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
				"activity_staff_prize.user_id", "activity_staff_prize.round",
				"activity_staff_prize.win_time", "activity_staff_prize.status",
				"activity_staff_prize.score",
				"activity_staff_prize.score_2",
				"activity_staff_prize.rank",
				"activity_staff_prize.team",
				"activity_staff_prize.leader",
				"activity_staff_prize.mvp",
				"activity_staff_prize.game",

				// 活動資訊
				"activity.user_id as activity_user_id",

				// 遊戲
				// "activity_game.game", "activity_game.title",
				// 隊伍資訊
				// "activity_game.left_team_name", "activity_game.right_team_name",

				// 獎品
				"activity_prize.prize_name", "activity_prize.prize_type",
				"activity_prize.prize_picture", "activity_prize.prize_price",
				"activity_prize.prize_method", "activity_prize.prize_password",

				// 用戶
				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_staff_prize.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			Where("activity_staff_prize.game_id", "=", gameID).
			Where("activity_staff_prize.prize_id", "!=", "").
			Where("activity_prize.prize_type", "!=", "thanks").
			Where("activity_prize.prize_type", "!=", "").
			Where("activity_prize.prize_method", "!=", "thanks").
			Where("activity_prize.prize_method", "!=", "").
			OrderBy("activity_staff_prize.game", "asc",
				"activity_staff_prize.game_id", "asc",
				"activity_staff_prize.round", "asc",
				"activity_staff_prize.status", "asc",
				"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
				"activity_staff_prize.prize_id", "asc",
			)
		if round == 0 {
			// 不分輪次(例如遊戲抽獎)
		} else {
			// 特定輪次
			sql = sql.Where("activity_staff_prize.round", "=", round)
		}
		// fmt.Println("有執行嗎? ")

		var items = make([]map[string]interface{}, 0)
		items, err = sql.All()
		if err != nil {
			return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
		}
		records = MapToPrizeStaffModel(items)
		// fmt.Println("資料長度: ", len(records))

		// 將中獎人員資料添加至redis中(list類型)
		if isRedis && len(records) > 0 {
			var params = []interface{}{config.WINNING_STAFFS_REDIS + gameID}
			for _, record := range records {
				params = append(params, utils.JSON(record))
			}

			// 將中獎人員資訊加入redis中
			m.RedisConn.ListMultiRPush(params)
			// m.RedisConn.SetExpire(config.WINNING_STAFFS_REDIS+gameID,
			// 	config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			m.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
		}

		// 舊版: 依照獎品類型分類至redis中(zset類型)
		// if isRedis && len(records) > 0 {
		// 	for _, record := range records {
		// 		var prizeType int64
		// 		if record.PrizeType == "first" {
		// 			prizeType = 1
		// 		} else if record.PrizeType == "second" {
		// 			prizeType = 2
		// 		} else if record.PrizeType == "third" {
		// 			prizeType = 3
		// 		} else if record.PrizeType == "general" {
		// 			prizeType = 4
		// 		}

		// 		recordJson := utils.JSON(record) // 中獎人員json編碼資料
		// 		// 加入redis
		// 		m.RedisConn.ZSetAdd(config.WINNING_STAFFS_REDIS+gameID, recordJson, prizeType)
		// 	}
		// }
	}
	return records, nil
}

// Export 匯出中獎人員資料
func (m PrizeStaffModel) Export(activityID, gameID, game, userID, round, status string,
	limit, offset int64, host string) (*excelize.File, error) {
	var (
		file                 = excelize.NewFile()      // 開啟EXCEL檔案
		index, _             = file.NewSheet("Sheet1") // 創建一個工作表
		customizeModel, err1 = DefaultCustomizeModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					Find(activityID) // 自定義欄位
		rowNames = []string{"A", "B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
			"S", "T", "U", "V", "W", "X", "Y", "Z",
			"AA", "AB", "AC", "AD", "AE", "AF", "AG",
			"AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR",
			"AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}
		colNames = []string{
			// 用戶資訊
			"用戶名稱", "抽獎號碼",
			// 中獎資訊
			"兌獎狀況", "中獎時間", "名次", "分數", "中獎輪次",
			// 隊伍資訊
			"隊伍", "隊長", "mvp",
			// 遊戲資訊
			"遊戲名稱", "遊戲類型",
			// 獎品資訊
			"獎品名稱", "獎品類型", "獎品圖片", "獎品價值",
			"兌獎方式", "兌獎密碼",
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
		prizeStaffModel, err2 = DefaultPrizeStaffModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					FindAll(activityID, gameID, userID, game, round, status, "", limit, offset)
		sheet                                                   = "Sheet1"
		fileName, gameCN, gameIDCN, userIDCN, roundCN, statusCN string
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
	for i, staff := range prizeStaffModel {
		var (
			newGame, newPrizeType, newStatus, newPicture,
			newMethod, newTeam, newLeader, newMvp string
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

		// 判斷兌獎狀況(yes、no)
		if staff.Status == "yes" {
			newStatus = "已領獎"
		} else if staff.Status == "no" {
			newStatus = "尚未領獎"
		}
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
		// 判斷獎品類型(first、second、third、general)
		if staff.PrizeType == "first" {
			newPrizeType = "頭獎"
		} else if staff.PrizeType == "second" {
			newPrizeType = "二獎"
		} else if staff.PrizeType == "third" {
			newPrizeType = "三獎"
		} else if staff.PrizeType == "general" {
			newPrizeType = "普通獎"
		}
		// 判斷兌獎方式(site、mail)
		if staff.PrizeMethod == "site" {
			newMethod = "現場兌獎"
		} else if staff.PrizeMethod == "mail" {
			newMethod = "郵寄兌獎"
		}
		// 判斷獎品圖片
		if !strings.Contains(utils.GetString(staff.PrizePicture, "/admin/uploads/system/img-prize-pic.png"), "https") &&
			!strings.Contains(utils.GetString(staff.PrizePicture, "/admin/uploads/system/img-prize-pic.png"), "http") {
			newPicture = "https://" + host + utils.GetString(staff.PrizePicture, "")
		}
		// 獎品密碼解碼
		// password, _ := utils.Decode([]byte(utils.GetString(staff.PrizePassword, "")))
		// newPassword = string(password)
		// 判斷隊伍資訊
		if staff.Team == "" {
			// newTeam = "無隊伍"
		} else if staff.Team == "left_team" {
			newTeam = staff.LeftTeamName
		} else if staff.Team == "right_team" {
			newTeam = staff.RightTeamName
		}

		// 判斷隊長資訊
		if staff.Leader == "" {
		} else if staff.Leader == "yes" {
			newLeader = "是"
		} else if staff.Leader == "no" {
			newLeader = "否"
		}

		// 判斷隊長資訊
		if staff.Mvp == "" {
		} else if staff.Mvp == "yes" {
			newMvp = "是"
		} else if staff.Mvp == "no" {
			newMvp = "否"
		}

		values := []interface{}{
			// 用戶資訊
			staff.Name, staff.Number,
			// 中獎資訊
			newStatus, staff.WinTime, staff.Rank, staff.Score, staff.Round,
			// 隊伍資訊
			newTeam, newLeader, newMvp,
			// 遊戲資訊
			staff.Title, newGame,
			// 獎品資訊
			staff.PrizeName, newPrizeType, newPicture,
			staff.PrizePrice, newMethod, staff.PrizePassword,
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

	// 判斷是否取得特定狀態資料
	if status != "" {
		if status == "yes" {
			statusCN = "已兌獎"
		} else if status == "no" {
			statusCN = "未兌獎"
		}
	} else if status == "" {
		statusCN = "所有狀態"
	}

	fileName = "中獎人員-" + activityID + "-" + gameCN + "-" +
		gameIDCN + "-" + userIDCN + "-" + roundCN + "-" + statusCN
	// 儲存EXCEL
	if err := file.SaveAs(fmt.Sprintf(config.STORE_PATH+"/excel/%s.xlsx", fileName)); err != nil {
		return file, err
	}

	return file, nil
}

// Adds 新增多筆中獎人員資料
func (m PrizeStaffModel) Adds(dataAmount int, activityID, gameID, game string, userIDs []string, prizeIDs []string,
	round int, scores []string, scores2 []string, ranks []string,
	teams []string, leaders []string, mvps []string) error {
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

		// 將資料匹量寫入activity_staff_prize表中
		err := m.Table(m.TableName).BatchInsert(
			dataAmount, "activity_id,game_id,user_id,prize_id,round,score,score_2,`rank`,team,leader,mvp,game",
			[][]string{activityIDs, gameIDs, userIDs, prizeIDs, rounds, scores, scores2, ranks, teams, leaders, mvps, games})
		if err != nil {
			return errors.New("錯誤: 匹量新增中獎人員資料發生問題(activity_staff_prize)，請重新操作")
		}
	}

	return nil
}

// Add 增加資料
func (m PrizeStaffModel) Add(model NewPrizeStaffModel) (int64, error) {

	if model.Status != "yes" && model.Status != "no" {
		return 0, errors.New("錯誤: 兌獎狀態資料發生問題，請輸入有效的兌獎狀態")
	}
	// if model.White != "yes" && model.White != "no" {
	// 	return 0, errors.New("錯誤: 白名單資料發生問題，請輸入有效的白名單資料")
	// }
	if _, err := strconv.Atoi(model.Round); err != nil {
		return 0, errors.New("錯誤: 遊戲輪次資料發生問題，請輸入有效的輪次資料")
	}

	id, err := m.Table(m.TableName).Insert(command.Value{
		"user_id":     model.UserID,
		"activity_id": model.ActivityID,
		"game_id":     model.GameID,
		"game":        model.Game,
		"prize_id":    model.PrizeID,
		"round":       model.Round,
		"status":      model.Status,
		"score":       model.Score,
		"score_2":     model.Score2,
		"rank":        model.Rank,
		"team":        model.Team,
		"leader":      model.Leader,
		"mvp":         model.Mvp,
		// "white":       model.White,
	})
	return id, err
}

// Update 更新資料(管理員、手機用戶兌獎)
func (m PrizeStaffModel) Update(model EditPrizeStaffModel) error {
	if model.Role == "admin" {
		var (
			fieldValues = command.Value{}
			fields      = []string{"status"}
		)
		if model.Status != "" {
			if model.Status != "yes" && model.Status != "no" {
				return errors.New("錯誤: 兌獎狀態資料發生問題，請輸入有效的兌獎狀態")
			}
		}

		// 將 struct 轉換為 map[string]interface{} 格式
		data := utils.StructToMap(model)

		for _, key := range fields {
			if val, ok := data[key]; ok && val != "" {
				fieldValues[key] = val
			}
		}
		if len(fieldValues) == 0 {
			return nil
		}

		if err := m.Table(m.Base.TableName).
			WhereIn("activity_staff_prize.id", model.ID).
			Update(fieldValues); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 兌獎發生問題，請重新操作")
		}
	} else if model.Role == "guest" {
		if err := m.Table(m.Base.TableName).
			LeftJoin(command.Join{
				FieldA:    "activity_staff_prize.prize_id",
				FieldA1:   "activity_prize.prize_id",
				Table:     "activity_prize",
				Operation: "=",
			}).
			WhereIn("activity_staff_prize.id", model.ID).
			Where("activity_prize.prize_password", "=", string(utils.Encode([]byte(model.Password)))).
			Update(command.Value{"activity_staff_prize.status": "yes"}); err != nil {
			if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
				return errors.New("錯誤: 兌獎密碼發生問題，請輸入正確的密碼資料")
			}
			return errors.New("錯誤: 兌獎發生問題，請重新操作")
		}
	}
	return nil
}

// MapToPrizeStaffModel map轉換[]PrizeStaffModel
func MapToPrizeStaffModel(items []map[string]interface{}) []PrizeStaffModel {
	var staffs = make([]PrizeStaffModel, 0)
	for _, item := range items {
		var (
			staff PrizeStaffModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &staff)

		staff.WinTime = staff.WinTime[:len(staff.WinTime)-3] // 不顯示秒

		// 獎品密碼解碼
		// password, _ := item["prize_password"].(string)
		passwordByte, _ := utils.Decode([]byte(utils.GetString(staff.PrizePassword, "")))
		// password := string(passwordByte)
		staff.PrizePassword = string(passwordByte)

		// staff.PrizePicture, _ = item["prize_picture"].(string)
		if !strings.Contains(staff.PrizePicture, "system") {
			staff.PrizePicture = "/admin/uploads/" + staff.ActivityUserID + "/" + staff.ActivityID + "/interact/game/" + staff.Game + "/" + staff.GameID + "/" + staff.PrizePicture
		}

		staffs = append(staffs, staff)
	}
	return staffs
}

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

// 自定義欄位資訊
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

// staff.ID, _ = item["id"].(int64)
// staff.ActivityID, _ = item["activity_id"].(string)
// staff.ActivityUserID, _ = item["activity_user_id"].(string)
// staff.GameID, _ = item["game_id"].(string)
// staff.PrizeID, _ = item["prize_id"].(string)
// staff.UserID, _ = item["user_id"].(string)
// staff.Round, _ = item["round"].(int64)
// staff.WinTime, _ = item["win_time"].(string)
// staff.Status, _ = item["status"].(string)
// staff.White, _ = item["white"].(string)

// 隊伍資訊
// staff.Team, _ = item["team"].(string)
// staff.Leader, _ = item["leader"].(string)
// staff.Mvp, _ = item["mvp"].(string)

// 分數
// staff.Score, _ = item["score"].(int64)
// staff.Score2, _ = item["score_2"].(float64)
// staff.Rank, _ = item["rank"].(int64)

// staff.PrizePrice, _ = item["prize_price"].(int64)
// staff.PrizeMethod, _ = item["prize_method"].(string)

// 用戶資訊
// staff.Name, _ = item["name"].(string)
// staff.Avatar, _ = item["avatar"].(string)
// staff.Phone, _ = item["phone"].(string)
// staff.Email, _ = item["email"].(string)
// staff.ExtEmail, _ = item["ext_email"].(string)

// 遊戲資訊
// staff.Title, _ = item["title"].(string)
// staff.Game, _ = item["game"].(string)
// staff.LeftTeamName, _ = item["left_team_name"].(string)
// staff.RightTeamName, _ = item["right_team_name"].(string)

// 獎品資訊
// staff.PrizeName, _ = item["prize_name"].(string)
// staff.PrizeType, _ = item["prize_type"].(string)

// FindUserWinningRecords 查詢用戶的中獎紀錄(join activity_game、activity_prize join)
// func (m PrizeStaffModel) FindUserWinningRecords(userID, field, value string) ([]PrizeStaffModel, error) {
// 	items, err := m.Table(m.Base.TableName).
// 		Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 			"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 			"activity_staff_prize.user_id", "activity_staff_prize.round",
// 			"activity_staff_prize.win_time", "activity_staff_prize.status",
// 			"activity_staff_prize.score",
// 			"activity_staff_prize.score_2",
// 			"activity_staff_prize.rank",

// 			"activity_game.title", "activity_game.game",

// 			"activity_prize.prize_name", "activity_prize.prize_type",
// 			"activity_prize.prize_picture", "activity_prize.prize_price",
// 			"activity_prize.prize_method", "activity_prize.prize_password",

// 			"line_users.name", "line_users.avatar").
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:   "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.game_id",
// 			FieldA1:   "activity_game.game_id",
// 			Table:     "activity_game",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:   "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		Where(field, "=", value).
// 		Where("activity_staff_prize.user_id", "=", userID).
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_type", "!=", "").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "").All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得中獎人員資訊，請重新查詢")
// 	}
// 	return MapToPrizeStaffModel(items), nil
// }

// FindWinningRecordsByActivityID 查詢該活動下或該遊戲種類下所有的中獎紀錄(join line_users、activity_game、activity_prize join)
// func (m PrizeStaffModel) FindWinningRecordsByActivityID(activityID, game string) ([]PrizeStaffModel, error) {
// 	var (
// 		sql = m.Table(m.Base.TableName).
// 			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 				"activity_staff_prize.user_id", "activity_staff_prize.round",
// 				"activity_staff_prize.win_time", "activity_staff_prize.status",
// 				"activity_staff_prize.score",
// 				"activity_staff_prize.score_2",
// 				"activity_staff_prize.rank",

// 				// 遊戲
// 				"activity_game.game", "activity_game.title",

// 				// 獎品
// 				"activity_prize.prize_name", "activity_prize.prize_type",
// 				"activity_prize.prize_picture", "activity_prize.prize_price",
// 				"activity_prize.prize_method", "activity_prize.prize_password",

// 				// 用戶
// 				"line_users.name", "line_users.avatar").
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.user_id",
// 				FieldA1:    "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.prize_id",
// 				FieldA1:    "activity_prize.prize_id",
// 				Table:     "activity_prize",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.game_id",
// 				FieldA1:    "activity_game.game_id",
// 				Table:     "activity_game",
// 				Operation: "="}).
// 			Where("activity_staff_prize.activity_id", "=", activityID).
// 			Where("activity_staff_prize.prize_id", "!=", "").
// 			Where("activity_prize.prize_type", "!=", "thanks").
// 			Where("activity_prize.prize_type", "!=", "").
// 			Where("activity_prize.prize_method", "!=", "thanks").
// 			Where("activity_prize.prize_method", "!=", "").
// 			OrderBy("activity_game.game", "asc",
// 				"activity_staff_prize.game_id", "asc",
// 				"activity_staff_prize.round", "asc",
// 				"activity_staff_prize.status", "asc",
// 				"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 				"activity_staff_prize.prize_id", "asc",
// 			)
// 		items   = make([]map[string]interface{}, 0)
// 		records = make([]PrizeStaffModel, 0)
// 		err     error
// 	)
// 	if game != "" {
// 		// 該遊戲種類下所有中獎紀錄
// 		sql = sql.Where("activity_game.game", "=", game).
// 			Where("activity_prize.game", "=", game)
// 	}

// 	items, err = sql.All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
// 	}
// 	records = MapToPrizeStaffModel(items)

// 	return records, nil
// }

// FindWinningRecordsByGame 查詢該遊戲種類下所有的中獎紀錄(join line_users、activity_game、activity_prize join)
// func (m PrizeStaffModel) FindWinningRecordsByGame(activityID, game string) ([]PrizeStaffModel, error) {
// 	var (
// 		sql = m.Table(m.Base.TableName).
// 			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 				"activity_staff_prize.user_id", "activity_staff_prize.round",
// 				"activity_staff_prize.win_time", "activity_staff_prize.status",

// 				// 遊戲
// 				"activity_game.game", "activity_game.title",

// 				// 獎品
// 				"activity_prize.prize_name", "activity_prize.prize_type",
// 				"activity_prize.prize_picture", "activity_prize.prize_price",
// 				"activity_prize.prize_method", "activity_prize.prize_password",

// 				// 用戶
// 				"line_users.name", "line_users.avatar").
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.user_id",
// 				FieldA1:    "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.prize_id",
// 				FieldA1:    "activity_prize.prize_id",
// 				Table:     "activity_prize",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.game_id",
// 				FieldA1:    "activity_game.game_id",
// 				Table:     "activity_game",
// 				Operation: "="}).
// 			Where("activity_staff_prize.activity_id", "=", activityID).
// 			Where("activity_game.game", "=", game).
// 			Where("activity_prize.game", "=", game).
// 			Where("activity_staff_prize.prize_id", "!=", "").
// 			Where("activity_prize.prize_type", "!=", "thanks").
// 			Where("activity_prize.prize_type", "!=", "").
// 			Where("activity_prize.prize_method", "!=", "thanks").
// 			Where("activity_prize.prize_method", "!=", "").
// 			OrderBy("activity_game.game", "asc",
// 				"activity_staff_prize.game_id", "asc",
// 				"activity_staff_prize.round", "asc",
// 				"activity_staff_prize.status", "asc",
// 				"field", "(activity_prize.prize_type, 'first', 'second', 'third', 'general')",
// 				"activity_staff_prize.prize_id", "asc",
// 			)
// 		items   = make([]map[string]interface{}, 0)
// 		records = make([]PrizeStaffModel, 0)
// 		err     error
// 	)

// 	items, err = sql.All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
// 	}
// 	records = MapToPrizeStaffModel(items)

// 	return records, nil
// }

// FindWinningRecordsByGameID 查詢該遊戲場次的中獎紀錄(join line_users、activity_prize join)，用在沒有輪次分別的遊戲(遊戲抽獎)
// func (m PrizeStaffModel) FindWinningRecordsByGameID(isRedis bool,
// 	gameID string) (staffs []PrizeStaffModel, err error) {
// 	// var (
// 	// 	items = make([]map[string]interface{}, 0)
// 	// 	err   error
// 	// )

// 	// 從redis取得中獎人員資料
// 	if isRedis {
// 		var staffsJson []string // 中獎人員json資料(獎品類型排列，1頭獎.2二獎...依序排列)
// 		staffsJson, err = m.RedisConn.ZSetRange(config.WINNING_STAFFS_REDIS+gameID, 0, 1000)
// 		if err != nil {
// 			return staffs, errors.New("錯誤: 從redis中取得中獎人員資訊發生問題")
// 		}

// 		if len(staffsJson) > 0 {
// 			for _, staffJson := range staffsJson {
// 				var staff PrizeStaffModel
// 				// 解碼
// 				json.Unmarshal([]byte(staffJson), &staff)

// 				staffs = append(staffs, staff)
// 			}
// 		}
// 	}

// 	if len(staffs) == 0 {
// 		items, err := m.Table(m.Base.TableName).
// 			Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 				"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 				"activity_staff_prize.user_id", "activity_staff_prize.round",
// 				"activity_staff_prize.win_time", "activity_staff_prize.status",

// 				"activity_prize.prize_name", "activity_prize.prize_type",
// 				"activity_prize.prize_picture", "activity_prize.prize_price",
// 				"activity_prize.prize_method", "activity_prize.prize_password",

// 				"line_users.name", "line_users.avatar").
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.user_id",
// 				FieldA1:    "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_prize.prize_id",
// 				FieldA1:    "activity_prize.prize_id",
// 				Table:     "activity_prize",
// 				Operation: "="}).
// 			Where("activity_staff_prize.game_id", "=", gameID).
// 			Where("activity_prize.prize_type", "!=", "thanks").
// 			Where("activity_prize.prize_method", "!=", "thanks").
// 			Where("activity_staff_prize.prize_id", "!=", "").All()
// 		if err != nil {
// 			return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
// 		}

// 		staffs = MapToPrizeStaffModel(items)
// 		// 將中獎人員資料添加至redis
// 		if isRedis {
// 			for _, staff := range staffs {
// 				var prizeType int64
// 				if staff.PrizeType == "first" {
// 					prizeType = 1
// 				} else if staff.PrizeType == "second" {
// 					prizeType = 2
// 				} else if staff.PrizeType == "third" {
// 					prizeType = 3
// 				} else if staff.PrizeType == "general" {
// 					prizeType = 4
// 				}

// 				staffJson := utils.JSON(staff) // 中獎人員json編碼資料
// 				// 加入redis
// 				m.RedisConn.ZSetAdd(config.WINNING_STAFFS_REDIS+gameID, staffJson, prizeType)
// 			}
// 		}
// 	}
// 	return
// }

// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:    "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:    "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:    "activity_score.user_id",
// 			Table:     "activity_score",
// 			Operation: "="}).
// 		Where("activity_staff_prize.game_id", "=", gameID).
// 		Where("activity_staff_prize.round", "=", round).
// 		Where("activity_score.game_id", "=", gameID).
// 		Where("activity_score.round", "=", round).
// 		Where("activity_score.score", ">", 0).
// 		Where("activity_prize.prize_type", "!=", "thanks").
// 		Where("activity_prize.prize_method", "!=", "thanks").
// 		Where("activity_staff_prize.prize_id", "!=", "").
// 		OrderBy("activity_score.score", "desc", "activity_score.id", "asc").All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得中獎紀錄，請重新查詢")
// 	}
// 	return MapToPrizeStaffModel(items), nil
// }

// res, err := m.DbConn.Exec(
// 	fmt.Sprintf("update `activity_staff_prize` INNER JOIN `activity_prize` ON `activity_staff_prize`.`prize_id` = `activity_prize`.`prize_id` set `activity_staff_prize`.`status` = 'yes' where `activity_prize`.`prize_password` = '%s' and `activity_staff_prize`.`id` = %s",
// 		string(utils.Encode([]byte(model.Password))), model.ID))
// if err != nil {
// 	return errors.New("錯誤: 兌獎發生問題，請重新操作")
// }
// if affectRow, _ := res.RowsAffected(); affectRow < 1 {
// 	return errors.New("錯誤: 兌獎密碼發生問題，請輸入正確的密碼資料")
// }

// FindWinningRecordsOrderByScore 透過輪次查詢中獎紀錄並依照分數高低排序(join line_users、activity_prize、activity_whack_mole_score join)
// func (m PrizeStaffModel) FindWinningRecordsOrderByScore(gameID string, round int64) ([]PrizeStaffModel, error) {
// 	// var (
// 	// 	items = make([]map[string]interface{}, 0)
// 	// 	err   error
// 	// )
// 	items, err := m.Table(m.Base.TableName).
// 		Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 			"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 			"activity_staff_prize.user_id", "activity_staff_prize.round",
// 			"activity_staff_prize.win_time", "activity_staff_prize.status",

// 			"activity_prize.prize_name", "activity_prize.prize_type",
// 			"activity_prize.prize_picture", "activity_prize.prize_price",
// 			"activity_prize.prize_method", "activity_prize.prize_password",

// 			"line_users.name", "line_users.avatar",

// 			"activity_score.score").

// 用戶資訊
// Name     string `json:"name" example:"Name"`
// Avatar   string `json:"avatar" example:"https://..."`
// Phone    string `json:"phone" example:"0912345678"`
// Email    string `json:"email" example:"test@gmail.com"`
// ExtEmail string `json:"ext_email" example:"test@gmail.com"`

// 自定義
// Ext1       string `json:"ext_1" example:"ext1"`
// Ext2       string `json:"ext_2" example:"ext2"`
// Ext3       string `json:"ext_3" example:"ext3"`
// Ext4       string `json:"ext_4" example:"ext4"`
// Ext5       string `json:"ext_5" example:"ext5"`
// Ext6       string `json:"ext_6" example:"ext6"`
// Ext7       string `json:"ext_7" example:"ext7"`
// Ext8       string `json:"ext_8" example:"ext8"`
// Ext9       string `json:"ext_9" example:"ext9"`
// Ext10      string `json:"ext_10" example:"ext10"`

// Score      int64  `json:"score"`
// Phone    string `json:"phone"`
// Email    string `json:"email"`
// ExtEmail string `json:"ext_email"`
// 自定義
// Ext1       string `json:"ext_1"`
// Ext2       string `json:"ext_2"`
// Ext3       string `json:"ext_3"`
// Ext4       string `json:"ext_4"`
// Ext5       string `json:"ext_5"`
// Ext6       string `json:"ext_6"`
// Ext7       string `json:"ext_7"`
// Ext8       string `json:"ext_8"`
// Ext9       string `json:"ext_9"`
// Ext10      string `json:"ext_10"`

// Name       string `json:"name" example:"User Name"`
// Avatar     string `json:"avatar" example:"https://..."`
// PrizeName  string `json:"prize_name" example:"Prize Name"`
// PrizeType  string `json:"prize_type" example:"money、prize、thanks、first、second、third、special"`
// Picture    string `json:"picture" example:"https://..."`
// Price      string `json:"price" example:"1000"`
// Method     string `json:"method" example:"site、mail、thanks"`
// Password   string `json:"password" example:"password(最多為8個字元)"`
// WinTime  string `json:"win_time" example:"2021-01-01 00:00"`
// Phone    string `json:"phone" example:"0912345678"`
// Email    string `json:"email" example:"test@gmail.com"`
// ExtEmail string `json:"ext_email" example:"test@gmail.com"`
// Ext1     string `json:"ext_1" example:"ext1"`
// Ext2     string `json:"ext_2" example:"ext2"`
// Ext3     string `json:"ext_3" example:"ext3"`
// Ext4     string `json:"ext_4" example:"ext4"`
// Ext5     string `json:"ext_5" example:"ext5"`
// Ext6     string `json:"ext_6" example:"ext6"`
// Ext7     string `json:"ext_7" example:"ext7"`
// Ext8     string `json:"ext_8" example:"ext8"`
// Ext9     string `json:"ext_9" example:"ext9"`
// Ext10    string `json:"ext_10" example:"ext10"`

// LeftJoinPrizeAndLineUsersByUserID 查詢用戶該遊戲場次中獎及未中獎紀錄資料(join line_users、activity_prize join)
// func (m PrizeStaffModel) LeftJoinPrizeAndLineUsersByUserID(gameID, userID string) ([]PrizeStaffModel, error) {
// 	items, err := m.Table(m.Base.TableName).
// 		Select("activity_staff_prize.id", "activity_staff_prize.activity_id",
// 			"activity_staff_prize.game_id", "activity_staff_prize.prize_id",
// 			"activity_staff_prize.user_id", "activity_staff_prize.round",
// 			"activity_staff_prize.win_time", "activity_staff_prize.status",

// 			"activity_prize.prize_name", "activity_prize.prize_type",
// 			"activity_prize.picture", "activity_prize.price",
// 			"activity_prize.method", "activity_prize.password",

// 			"line_users.name", "line_users.avatar").
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.user_id",
// 			FieldA1:    "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_staff_prize.prize_id",
// 			FieldA1:    "activity_prize.prize_id",
// 			Table:     "activity_prize",
// 			Operation: "="}).
// 		Where("activity_staff_prize.game_id", "=", gameID).
// 		Where("activity_staff_prize.user_id", "=", userID).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得用戶遊戲資料，請重新查詢")
// 	}
// 	return MapToPrizeStaffModel(items), nil
// }

// GetPeople 遊戲報名人數
// func (m GameStaffModel) GetPeople(gameID string) (int64, error) {
// 	items, err := m.Table(m.Base.TableName).Where("game_id", "=", gameID).
// 		Where("status", "=", "success").All()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(len(items)), nil
// }

// if utf8.RuneCountInString(model.Name) > 20 {
// 	return errors.New("錯誤: 姓名上限為20個字元，請輸入有效的姓名")
// }
// 取得報名簽到人員資料
// applysignModel, err := DefaultApplysignModel().SetDbConn(m.DbConn).
// 	FindByUserID(model.ActivityID, model.UserID)
// if err != nil || applysignModel.ID == 0 {
// 	return errors.New("錯誤: 無法查詢報名簽到人員資料，請重新操作")
// }
// "name":        model.Name,
// "avatar":      model.Avatar,
// "phone":       applysignModel.Phone,
// "email":       applysignModel.Email,
// "ext_email":   applysignModel.ExtEmail,
// "ext_1":       applysignModel.Ext1,
// "ext_2":       applysignModel.Ext2,
// "ext_3":       applysignModel.Ext3,
// "ext_4":       applysignModel.Ext4,
// "ext_5":       applysignModel.Ext5,
// "ext_6":       applysignModel.Ext6,
// "ext_7":       applysignModel.Ext7,
// "ext_8":       applysignModel.Ext8,
// "ext_9":       applysignModel.Ext9,
// "ext_10":      applysignModel.Ext10,
// 更新遊戲人數資訊
// gameModel, err := DefaultGameModel().SetDbConn(m.DbConn).Find(model.GameID)
// if err != nil || gameModel.ID == 0 {
// 	return errors.New("錯誤: 無法辨識遊戲資訊")
// }
// if err = DefaultGameModel().SetDbConn(m.DbConn).Update(EditGameModel{
// 	GameID: model.GameID,
// 	Attend: "1",
// }); err != nil {
// 	return err
// }

// UpdateUser 更新遊戲參加人員姓名、頭像資料
// func (m GameStaffModel) UpdateUser(userID, name, avatar, email, extEmail, phone string) error {
// 	var (
// 		fields      = []string{"phone", "ext_email"}
// 		values      = []string{phone, extEmail}
// 		fieldValues = command.Value{"name": name, "avatar": avatar, "email": email}
// 	)
// 	for i, field := range fields {
// 		if values[i] != "" {
// 			fieldValues[field] = values[i]
// 		}
// 	}

// 	if err := m.Table(m.Base.TableName).Where("user_id", "=", userID).
// 		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新遊戲參加人員資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// UpdateExt 更新參加遊戲人員Ext欄位資料
// func (m GameStaffModel) UpdateExt(activityID, userID string, values []string) error {
// 	var (
// 		fields = []string{"phone", "ext_email", "ext_1", "ext_2", "ext_3",
// 			"ext_4", "ext_5", "ext_6", "ext_7", "ext_8", "ext_9", "ext_10"}
// 		fieldValues = command.Value{}
// 	)
// 	for i, field := range fields {
// 		if values[i] != "" {
// 			fieldValues[field] = values[i]
// 		}
// 	}

// 	if err := m.Table(m.Base.TableName).Where("activity_id", "=", activityID).
// 		Where("user_id", "=", userID).Update(fieldValues); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新參加遊戲人員Ext欄位資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// DeleteExt 清空參加遊戲人員Ext欄位資料
// func (m GameStaffModel) DeleteExt(activityID, field string) error {
// 	if err := m.Table(m.Base.TableName).Where("activity_id", "=", activityID).
// 		Update(command.Value{field: ""}); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新參加遊戲人員Ext欄位資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// if utf8.RuneCountInString(model.Name) > 20 {
// 	return 0, errors.New("錯誤: 姓名上限為20個字元，請輸入有效的姓名")
// }
// if utf8.RuneCountInString(model.PrizeName) > 20 {
// 	return 0, errors.New("錯誤: 獎品名稱上限為20個字元，請輸入有效的獎品名稱")
// }
// if model.Avatar == "" || model.Picture == "" {
// 	return 0, errors.New("錯誤: 用戶及獎品照片資料不能為空")
// }
// if model.Method != "site" && model.Method != "mail" &&
// 	model.Method != "thanks" {
// 	return 0, errors.New("錯誤: 兌獎方式資料發生問題，請輸入有效的兌獎方式")
// }
// if model.Password == "" || utf8.RuneCountInString(model.Password) > 8 {
// 	return 0, errors.New("錯誤: 獎品密碼不能為空並且上限為8個字元，請輸入有效的獎品密碼")
// }
// if _, err := strconv.Atoi(model.Price); err != nil {
// 	return 0, errors.New("錯誤: 價格資料發生問題，請輸入有效的獎品價格")
// }

// 取得報名簽到人員資料
// applysignModel, err := DefaultApplysignModel().SetDbConn(m.DbConn).
// 	FindByUserID(model.ActivityID, model.UserID)
// if err != nil || applysignModel.ID == 0 {
// 	return 0, errors.New("錯誤: 無法查詢報名簽到人員資料，請重新操作")
// }

// "score":       model.Score,

// "name":        model.Name,
// "avatar":      model.Avatar,
// "prize_name":  model.PrizeName,
// "prize_type":  model.PrizeType,
// "picture":     model.Picture,
// "price":       model.Price,
// "method":      model.Method,
// "password":    model.Password,
// "win_time":    model.WinTime,
// "phone":     applysignModel.Phone,
// "email":     applysignModel.Email,
// "ext_email": applysignModel.ExtEmail,
// "ext_1":     applysignModel.Ext1,
// "ext_2":     applysignModel.Ext2,
// "ext_3":     applysignModel.Ext3,
// "ext_4":     applysignModel.Ext4,
// "ext_5":     applysignModel.Ext5,
// "ext_6":     applysignModel.Ext6,
// "ext_7":     applysignModel.Ext7,
// "ext_8":     applysignModel.Ext8,
// "ext_9":     applysignModel.Ext9,
// "ext_10":    applysignModel.Ext10,
// UpdateUser 更新中獎人員姓名、頭像資料
// func (m PrizeStaffModel) UpdateUser(userID, name, avatar, email, extEmail, phone string) error {
// 	var (
// 		fields      = []string{"phone", "ext_email"}
// 		values      = []string{phone, extEmail}
// 		fieldValues = command.Value{"name": name, "avatar": avatar, "email": email}
// 	)
// 	for i, field := range fields {
// 		if values[i] != "" {
// 			fieldValues[field] = values[i]
// 		}
// 	}
// 	if err := m.Table(m.Base.TableName).Where("user_id", "=", userID).
// 		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新中獎人員資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// UpdateExt 更新中獎人員Ext欄位資料
// func (m PrizeStaffModel) UpdateExt(activityID, userID string, values []string) error {
// 	var (
// 		fields = []string{"phone", "ext_email", "ext_1", "ext_2", "ext_3",
// 			"ext_4", "ext_5", "ext_6", "ext_7", "ext_8", "ext_9", "ext_10"}
// 		fieldValues = command.Value{}
// 	)
// 	for i, field := range fields {
// 		if values[i] != "" {
// 			fieldValues[field] = values[i]
// 		}
// 	}

// 	if err := m.Table(m.Base.TableName).Where("activity_id", "=", activityID).
// 		Where("user_id", "=", userID).Update(fieldValues); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新中獎人員Ext欄位資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// DeleteExt 清空中獎人員Ext欄位資料
// func (m PrizeStaffModel) DeleteExt(activityID, field string) error {
// 	if err := m.Table(m.Base.TableName).Where("activity_id", "=", activityID).
// 		Update(command.Value{field: ""}); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新中獎人員Ext欄位資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// 自定義
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
// staff.Score, _ = item["score"].(int64)
// staff.Phone, _ = item["phone"].(string)
// staff.Email, _ = item["email"].(string)
// staff.ExtEmail, _ = item["ext_email"].(string)
// 自定義
// staff.Ext1, _ = item["ext_1"].(string)
// staff.Ext2, _ = item["ext_2"].(string)
// staff.Ext3, _ = item["ext_3"].(string)
// staff.Ext4, _ = item["ext_4"].(string)
// staff.Ext5, _ = item["ext_5"].(string)
// staff.Ext6, _ = item["ext_6"].(string)
// staff.Ext7, _ = item["ext_7"].(string)
// staff.Ext8, _ = item["ext_8"].(string)
// staff.Ext9, _ = item["ext_9"].(string)
// staff.Ext10, _ = i

// fields      = []string{"status", "white"}
// values      = []string{model.Status, model.White}

// if model.White != "" {
// 	if model.White != "yes" && model.White != "no" {
// 		return errors.New("錯誤: 白名單資料發生問題，請輸入有效的白名單資料")
// 	}
// }

// values      = []string{model.Status}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fiel
