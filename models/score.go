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

// ScoreModel 資料表欄位
type ScoreModel struct {
	Base       `json:"-"`
	ID         int64   `json:"id"`          // ID
	ActivityID string  `json:"activity_id"` // 活動ID
	GameID     string  `json:"game_id"`     //遊戲ID
	UserID     string  `json:"user_id"`     // 用戶ID
	Round      int64   `json:"round"`       // 輪次
	Score      int64   `json:"score"`       // 分數
	Score2     float64 `json:"score_2"`     // 分數(第二判斷依據)

	// 隊伍資訊
	Team   string `json:"team"`   // 隊伍
	Leader string `json:"leader"` // 隊長
	Mvp    string `json:"mvp"`    // mvp
	// 名次
	// Rank int64 `json:"rank" example:"1"`

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

	// 答題紀錄
	QA1Option       string `json:"qa_1_option"`
	QA1OriginScore  int64  `json:"qa_1_origin_score"`
	QA1AddScore     int64  `json:"qa_1_add_score"`
	QA2Option       string `json:"qa_2_option"`
	QA2OriginScore  int64  `json:"qa_2_origin_score"`
	QA2AddScore     int64  `json:"qa_2_add_score"`
	QA3Option       string `json:"qa_3_option"`
	QA3OriginScore  int64  `json:"qa_3_origin_score"`
	QA3AddScore     int64  `json:"qa_3_add_score"`
	QA4Option       string `json:"qa_4_option"`
	QA4OriginScore  int64  `json:"qa_4_origin_score"`
	QA4AddScore     int64  `json:"qa_4_add_score"`
	QA5Option       string `json:"qa_5_option"`
	QA5OriginScore  int64  `json:"qa_5_origin_score"`
	QA5AddScore     int64  `json:"qa_5_add_score"`
	QA6Option       string `json:"qa_6_option"`
	QA6OriginScore  int64  `json:"qa_6_origin_score"`
	QA6AddScore     int64  `json:"qa_6_add_score"`
	QA7Option       string `json:"qa_7_option"`
	QA7OriginScore  int64  `json:"qa_7_origin_score"`
	QA7AddScore     int64  `json:"qa_7_add_score"`
	QA8Option       string `json:"qa_8_option"`
	QA8OriginScore  int64  `json:"qa_8_origin_score"`
	QA8AddScore     int64  `json:"qa_8_add_score"`
	QA9Option       string `json:"qa_9_option"`
	QA9OriginScore  int64  `json:"qa_9_origin_score"`
	QA9AddScore     int64  `json:"qa_9_add_score"`
	QA10Option      string `json:"qa_10_option"`
	QA10OriginScore int64  `json:"qa_10_origin_score"`
	QA10AddScore    int64  `json:"qa_10_add_score"`
	QA11Option      string `json:"qa_11_option"`
	QA11OriginScore int64  `json:"qa_11_origin_score"`
	QA11AddScore    int64  `json:"qa_11_add_score"`
	QA12Option      string `json:"qa_12_option"`
	QA12OriginScore int64  `json:"qa_12_origin_score"`
	QA12AddScore    int64  `json:"qa_12_add_score"`
	QA13Option      string `json:"qa_13_option"`
	QA13OriginScore int64  `json:"qa_13_origin_score"`
	QA13AddScore    int64  `json:"qa_13_add_score"`
	QA14Option      string `json:"qa_14_option"`
	QA14OriginScore int64  `json:"qa_14_origin_score"`
	QA14AddScore    int64  `json:"qa_14_add_score"`
	QA15Option      string `json:"qa_15_option"`
	QA15OriginScore int64  `json:"qa_15_origin_score"`
	QA15AddScore    int64  `json:"qa_15_add_score"`
	QA16Option      string `json:"qa_16_option"`
	QA16OriginScore int64  `json:"qa_16_origin_score"`
	QA16AddScore    int64  `json:"qa_16_add_score"`
	QA17Option      string `json:"qa_17_option"`
	QA17OriginScore int64  `json:"qa_17_origin_score"`
	QA17AddScore    int64  `json:"qa_17_add_score"`
	QA18Option      string `json:"qa_18_option"`
	QA18OriginScore int64  `json:"qa_18_origin_score"`
	QA18AddScore    int64  `json:"qa_18_add_score"`
	QA19Option      string `json:"qa_19_option"`
	QA19OriginScore int64  `json:"qa_19_origin_score"`
	QA19AddScore    int64  `json:"qa_19_add_score"`
	QA20Option      string `json:"qa_20_option"`
	QA20OriginScore int64  `json:"qa_20_origin_score"`
	QA20AddScore    int64  `json:"qa_20_add_score"`
}

// NewScoreModel 資料表欄位
type NewScoreModel struct {
	ActivityID string  `json:"activity_id" example:"activity_id"`
	GameID     string  `json:"game_id" example:"game_id"`
	Game       string  `json:"game" example:"game"`
	UserID     string  `json:"user_id" example:"user_id"`
	OptionID   string  `json:"option_id" example:"option_id"`
	Round      int64   `json:"round" example:"1"`
	Score      int64   `json:"score" example:"100"`
	Score2     float64 `json:"score_2" example:"100"`
	// 名次
	// Rank int64 `json:"rank" example:"1"`

	Team   string `json:"team" example:"left_team"` // 隊伍
	Leader string `json:"leader" example:"yes"`     // 隊長
	Mvp    string `json:"mvp" example:"yes"`        // mvp
}

// DefaultScoreModel 預設ScoreModel
func DefaultScoreModel() ScoreModel {
	return ScoreModel{Base: Base{TableName: config.ACTIVITY_SCORE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m ScoreModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ScoreModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m ScoreModel) SetDbConn(conn db.Connection) ScoreModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m ScoreModel) SetRedisConn(conn cache.Connection) ScoreModel {
// 	m.RedisConn = conn
// 	return m
// }

// // SetMongoConn 設定connection
// func (m ScoreModel) SetMongoConn(conn mongo.Connection) ScoreModel {
// 	m.MongoConn = conn
// 	return m
// }

// FindAll 查詢所有遊戲紀錄
func (m ScoreModel) FindAll(activityID, gameID, userID, game, round string,
	limit, offset int64) ([]ScoreModel, error) {
	var (
		scoresModel = make([]ScoreModel, 0)
		sql         = m.Table(m.Base.TableName).
				Select("activity_score.id", "activity_score.activity_id",
				"activity_score.game_id", "activity_score.user_id",
				"activity_score.round", "activity_score.score",
				"activity_score.score_2", "activity_score.team",
				"activity_score.leader", "activity_score.mvp",
				"activity_score.game",

				// 遊戲場次
				// "activity_game.title", "activity_game.game",
				// 隊伍資訊
				// "activity_game.left_team_name", "activity_game.right_team_name",

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

				// 答題記錄
				"activity_game_qa_record.qa_1_option", "activity_game_qa_record.qa_1_origin_score", "activity_game_qa_record.qa_1_add_score",
				"activity_game_qa_record.qa_2_option", "activity_game_qa_record.qa_2_origin_score", "activity_game_qa_record.qa_2_add_score",
				"activity_game_qa_record.qa_3_option", "activity_game_qa_record.qa_3_origin_score", "activity_game_qa_record.qa_3_add_score",
				"activity_game_qa_record.qa_4_option", "activity_game_qa_record.qa_4_origin_score", "activity_game_qa_record.qa_4_add_score",
				"activity_game_qa_record.qa_5_option", "activity_game_qa_record.qa_5_origin_score", "activity_game_qa_record.qa_5_add_score",
				"activity_game_qa_record.qa_6_option", "activity_game_qa_record.qa_6_origin_score", "activity_game_qa_record.qa_6_add_score",
				"activity_game_qa_record.qa_7_option", "activity_game_qa_record.qa_7_origin_score", "activity_game_qa_record.qa_7_add_score",
				"activity_game_qa_record.qa_8_option", "activity_game_qa_record.qa_8_origin_score", "activity_game_qa_record.qa_8_add_score",
				"activity_game_qa_record.qa_9_option", "activity_game_qa_record.qa_9_origin_score", "activity_game_qa_record.qa_9_add_score",
				"activity_game_qa_record.qa_10_option", "activity_game_qa_record.qa_10_origin_score", "activity_game_qa_record.qa_10_add_score",
				"activity_game_qa_record.qa_11_option", "activity_game_qa_record.qa_11_origin_score", "activity_game_qa_record.qa_11_add_score",
				"activity_game_qa_record.qa_12_option", "activity_game_qa_record.qa_12_origin_score", "activity_game_qa_record.qa_12_add_score",
				"activity_game_qa_record.qa_13_option", "activity_game_qa_record.qa_13_origin_score", "activity_game_qa_record.qa_13_add_score",
				"activity_game_qa_record.qa_14_option", "activity_game_qa_record.qa_14_origin_score", "activity_game_qa_record.qa_14_add_score",
				"activity_game_qa_record.qa_15_option", "activity_game_qa_record.qa_15_origin_score", "activity_game_qa_record.qa_15_add_score",
				"activity_game_qa_record.qa_16_option", "activity_game_qa_record.qa_16_origin_score", "activity_game_qa_record.qa_16_add_score",
				"activity_game_qa_record.qa_17_option", "activity_game_qa_record.qa_17_origin_score", "activity_game_qa_record.qa_17_add_score",
				"activity_game_qa_record.qa_18_option", "activity_game_qa_record.qa_18_origin_score", "activity_game_qa_record.qa_18_add_score",
				"activity_game_qa_record.qa_19_option", "activity_game_qa_record.qa_19_origin_score", "activity_game_qa_record.qa_19_add_score",
				"activity_game_qa_record.qa_20_option", "activity_game_qa_record.qa_20_origin_score", "activity_game_qa_record.qa_20_add_score",
			).
			LeftJoin(command.Join{
				FieldA:    "activity_score.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			// LeftJoin(command.Join{
			// 	FieldA:    "activity_score.game_id",
			// 	FieldA1:   "activity_game.game_id",
			// 	Table:     "activity_game",
			// 	Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_score.activity_id",
				FieldA1:   "activity_customize.activity_id",
				Table:     "activity_customize",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_score.activity_id",
				FieldA1:   "activity_applysign.activity_id",
				FieldB:    "activity_score.user_id",
				FieldB1:   "activity_applysign.user_id",
				Table:     "activity_applysign",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_score.game_id",
				FieldA1:   "activity_game_qa_record.game_id",
				FieldB:    "activity_score.round",
				FieldB1:   "activity_game_qa_record.round",
				FieldC:    "activity_score.user_id",
				FieldC1:   "activity_game_qa_record.user_id",
				Table:     "activity_game_qa_record",
				Operation: "="}).
			Where("activity_score.activity_id", "=", activityID).
			OrderBy(
				"activity_score.game", "asc",
				"activity_score.game_id", "asc",
				"activity_score.round", "asc",
				"activity_score.score", "desc",
			)
	)

	// 判斷參數是否為空
	if gameID != "" {
		sql = sql.WhereIn("activity_score.game_id", interfaces(strings.Split(gameID, ",")))
	}
	if userID != "" {
		sql = sql.WhereIn("activity_score.user_id", interfaces(strings.Split(userID, ",")))
	}
	if game != "" {
		sql = sql.WhereIn("activity_score.game", interfaces(strings.Split(game, ",")))
	}
	if round != "" {
		sql = sql.WhereIn("activity_score.round", interfaces(strings.Split(round, ",")))
	}

	if limit != 0 {
		sql = sql.Limit(limit)
	}
	if offset != 0 {
		sql = sql.Offset(offset)
	}

	items, err := sql.All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得遊戲紀錄，請重新查詢")
	}

	scoresModel = MapToScoreModel(items)

	// 多個遊戲場次資料(mongo)
	gamesModel, err := DefaultGameModel().
		SetConn(m.DbConn, m.RedisConn, m.MongoConn).
		FindAll(activityID, game)
	if err != nil {
		return scoresModel, err
	}

	// 把mysql.mongo資料轉成 map(要合併的欄位較少)，以便快速查找
	dataMap := make(map[string]GameModel)
	for _, game := range gamesModel {
		dataMap[game.GameID] = game
	}

	for i, item := range scoresModel {
		if dataItem, found := dataMap[item.GameID]; found {

			// 合併欄位（可依需要擴充）
			item.Title = dataItem.Title

			scoresModel[i] = item
		}
	}

	return scoresModel, nil
}

// FindTopUser 查詢分數前n名的人員(從redis中取得分數由高至低的玩家資訊)
func (m ScoreModel) FindTopUser(isRedis bool, gameID string, game string, round int64,
	limit int64, order string) ([]ScoreModel, error) {
	var (
		scores = make([]ScoreModel, 0)
		err    error
	)
	if isRedis {
		var users []string
		if users, err = m.RedisConn.ZSetRevRange(config.SCORES_REDIS+gameID, 0, limit); err != nil {
			return scores, errors.New("錯誤: 從redis中取得分數由高至低的玩家資訊發生問題")
		}
		// fmt.Println("redis裡的數: ", users)

		// 處理分數排名資訊(只取得分數大於0的資料)
		for _, userID := range users {
			var score2 float64

			score := m.RedisConn.ZSetIntScore(config.SCORES_REDIS+gameID, userID) // 分數資料
			// 第二分數資料(敲敲樂、鑑定師、快問快答)
			if game == "whack_mole" || game == "monopoly" || game == "QA" {
				score2 = m.RedisConn.ZSetFloatScore(config.SCORES_2_REDIS+gameID, userID)
				// correct := m.RedisConn.ZSetIntScore(config.CORRECT_REDIS+gameID, userID)   // 分數資料
			}

			// 敲敲樂、鑑定師、快問快答不取得0分玩家資料
			if (game == "whack_mole" || game == "monopoly" || game == "QA") &&
				score == 0 {
				break
			}

			// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
			user, err := DefaultLineModel().
				SetConn(m.DbConn, m.RedisConn, m.MongoConn).
				Find(true, "", "user_id", userID)
			if err != nil {
				return scores, errors.New("錯誤: 無法取得用戶資訊")
			}

			scores = append(scores, ScoreModel{
				ID:     user.ID,
				UserID: userID,
				Name:   user.Name,
				Avatar: user.Avatar,
				Score:  score,
				Score2: score2,
				// Correct: correct,
			})
			// fmt.Println("分數: ", score)
		}
		// fmt.Println("處理完排名後的數量: ", len(scores))
	} else {
		items, err := m.Table(m.Base.TableName).
			Select("activity_score.id", "activity_score.activity_id",
				"activity_score.game_id", "activity_score.user_id",
				"activity_score.round", "activity_score.score",
				"activity_score.score_2", "activity_score.team",
				"activity_score.leader", "activity_score.mvp",
				// "activity_score.rank",

				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_score.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_score.game_id", "=", gameID).
			Where("activity_score.round", "=", round).
			Where("activity_score.score", ">", 0).
			Limit(limit).
			OrderBy("activity_score.score", "desc", "activity_score.score_2", order).All() // 分數、花費時間排序
		// OrderBy("activity_score.score", "desc", "activity_score.id", "asc").All()
		if err != nil {
			return scores, errors.New("錯誤: 無法取得用戶分數資訊，請重新查詢")
		}
		scores = MapToScoreModel(items)
	}
	return scores, nil
}

// FindUser 查詢玩家分數(從redis中取得)
func (m ScoreModel) FindUser(isRedis bool, gameID string, userID string) ScoreModel {
	var (
		scoreModel ScoreModel
	)
	if isRedis {
		// 取得玩家分數資訊
		score := m.RedisConn.ZSetIntScore(config.SCORES_REDIS+gameID, userID) // 分數資料

		scoreModel = ScoreModel{
			UserID: userID,
			Score:  score,
		}
	}
	return scoreModel
}

// Adds 新增多筆分數資料
func (s ScoreModel) Adds(dataAmount int, activityID, gameID, game string, userIDs []string,
	round int, scores []string, scores2 []string,
	teams []string, leaders []string, mvps []string) error {
	// 判斷資料數
	if dataAmount > 0 {
		var (
			activityIDs = make([]string, dataAmount) // 活動ID
			gameIDs     = make([]string, dataAmount) // 活動ID
			rounds      = make([]string, dataAmount) // 輪次
			roundStr    = strconv.Itoa(round)        // 輪次(string)
			games       = make([]string, dataAmount) // 遊戲類型
		)

		// 處理活動.遊戲ID.輪次陣列資料
		for i := range dataAmount {
			activityIDs[i] = activityID
			gameIDs[i] = gameID
			rounds[i] = roundStr
			games[i] = game
		}

		// 將資料匹量寫入activity_score表中
		err := s.Table(s.TableName).BatchInsert(
			dataAmount, "activity_id,game_id,user_id,round,score,score_2,team,leader,mvp,game",
			[][]string{activityIDs, gameIDs, userIDs, rounds, scores, scores2, teams, leaders, mvps, games})
		if err != nil {
			return errors.New("錯誤: 匹量新增分數資料發生問題(activity_score)，請重新操作")
		}
	}

	return nil
}

// Add 增加資料(將玩家分數資料加入資料表中)
// func (s ScoreModel) Addd(model NewScoreModel) (err error) {
// 	if _, err = s.Table(s.TableName).Insert(command.Value{
// 		"activity_id": model.ActivityID,
// 		"game_id":     model.GameID,
// 		"game":        model.Game,
// 		"user_id":     model.UserID,
// 		"option_id":   model.OptionID,
// 		"round":       model.Round,
// 		"score":       model.Score,
// 		"score_2":     model.Score2,
// 		"team":        model.Team,
// 		"leader":      model.Leader,
// 		"mvp":         model.Mvp,
// 		// "rank":        model.Rank,
// 	}); err != nil {
// 		return errors.New("錯誤: 新增玩家分數資料發生問題，請重新操作")
// 	}
// 	return
// }

// UpdateScore 更新分數資料(更新redis中玩家分數資料)
func (s ScoreModel) UpdateScore(isRedis bool, gameID, userID string, score int64) (err error) {
	if isRedis {
		if err = s.RedisConn.ZSetAddInt(config.SCORES_REDIS+gameID,
			userID, score); err != nil {
			return errors.New("錯誤: 更新玩家分數發生問題(redis)，請重新操作")
		}

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
	}
	return
}

// UpdateScore2 更新玩家第二分數資料(更新redis中玩家第二分數資料)
func (s ScoreModel) UpdateScore2(isRedis bool, gameID, userID string, score float64) (err error) {
	if isRedis {
		if err = s.RedisConn.ZSetAddFloat(config.SCORES_2_REDIS+gameID,
			userID, score); err != nil {
			return errors.New("錯誤: 更新玩家花費時間資料發生問題(redis)，請重新操作")
		}
	}
	return
}

// UpdateCorrect 更新答對題數資料(更新redis中玩家答對題數資料)
func (s ScoreModel) UpdateCorrect(isRedis bool, gameID, userID string, correct int64) (err error) {
	if isRedis {
		if err = s.RedisConn.ZSetAddInt(config.CORRECT_REDIS+gameID,
			userID, correct); err != nil {
			return errors.New("錯誤: 更新玩家答對題數發生問題(redis)，請重新操作")
		}
	}
	return
}

// Export 匯出遊戲紀錄資料
func (m ScoreModel) Export(activityID, gameID, game, userID, round string,
	limit, offset int64) (*excelize.File, error) {
	var (
		file                 = excelize.NewFile()      // 開啟EXCEL檔案
		index, _             = file.NewSheet("Sheet1") // 創建一個工作表
		customizeModel, err1 = DefaultCustomizeModel().
					SetConn(m.DbConn, m.RedisConn, m.MongoConn).
					Find(activityID) // 自定義欄位
		rowNames = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K",
			"L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK",
			"AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV",
			"AW", "AX", "AY", "AZ",
			"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK",
			"BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV",
			"BW", "BX", "BY", "BZ",
			"CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK",
			"CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV",
			"CW", "CX", "CY", "CZ",
			"DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK",
			"DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV",
			"DW", "DX", "DY", "DZ",
		}
		colNames = []string{
			// 用戶資訊
			"用戶名稱", "抽獎號碼",
			// 遊戲資訊
			"遊戲名稱", "遊戲類型",
			// 參加資訊
			"遊戲輪次", "分數",
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
		// 自定義欄位是否必填
		qas = []string{
			"第一題選項", "第一題原始分數", "第一題加成分數",
			"第二題選項", "第二題原始分數", "第二題加成分數",
			"第三題選項", "第三題原始分數", "第三題加成分數",
			"第四題選項", "第四題原始分數", "第四題加成分數",
			"第五題選項", "第五題原始分數", "第五題加成分數",
			"第六題選項", "第六題原始分數", "第六題加成分數",
			"第七題選項", "第七題原始分數", "第七題加成分數",
			"第八題選項", "第八題原始分數", "第八題加成分數",
			"第九題選項", "第九題原始分數", "第九題加成分數",
			"第十題選項", "第十題原始分數", "第十題加成分數",
			"第十一題選項", "第十一題原始分數", "第十一題加成分數",
			"第十二題選項", "第十二題原始分數", "第十二題加成分數",
			"第十三題選項", "第十三題原始分數", "第十三題加成分數",
			"第十四題選項", "第十四題原始分數", "第十四題加成分數",
			"第十五題選項", "第十五題原始分數", "第十五題加成分數",
			"第十六題選項", "第十六題原始分數", "第十六題加成分數",
			"第十七題選項", "第十七題原始分數", "第十七題加成分數",
			"第十八題選項", "第十八題原始分數", "第十八題加成分數",
			"第十九題選項", "第十九題原始分數", "第十九題加成分數",
			"第二十題選項", "第二十題原始分數", "第二十題加成分數",
		}
		recordModel, err2 = DefaultScoreModel().
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

	// 判斷匯出遊戲是否包含拔河遊戲
	if game == "" || strings.Contains(game, "tugofwar") {
		colNames = append(colNames, "隊伍", "隊長", "mvp")
	}

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

	// 判斷匯出遊戲是否包含快問快答
	if game == "" || strings.Contains(game, "QA") {
		for _, qa := range qas {
			colNames = append(colNames, qa)
		}
	}

	// 將欄位中文名稱寫入EXCEL第一行
	for i, name := range colNames {
		file.SetCellValue(sheet, rowNames[i]+"1", name) // 設置存儲格的值
	}

	// 將所有中獎人員寫入EXCEL中(從第二行開始)
	for i, record := range recordModel {
		var (
			newGame, newTeam, newLeader, newMvp string
			extValues                           = []string{
				// 電話.信箱
				record.Phone, record.ExtEmail,
				// 自定義欄位
				record.Ext1, record.Ext2, record.Ext3,
				record.Ext4, record.Ext5, record.Ext6,
				record.Ext7, record.Ext8, record.Ext9,
				record.Ext10,
			}
			scores = []interface{}{
				record.QA1Option, record.QA1OriginScore, record.QA1AddScore,
				record.QA2Option, record.QA2OriginScore, record.QA2AddScore,
				record.QA3Option, record.QA3OriginScore, record.QA3AddScore,
				record.QA4Option, record.QA4OriginScore, record.QA4AddScore,
				record.QA5Option, record.QA5OriginScore, record.QA5AddScore,
				record.QA6Option, record.QA6OriginScore, record.QA6AddScore,
				record.QA7Option, record.QA7OriginScore, record.QA7AddScore,
				record.QA8Option, record.QA8OriginScore, record.QA8AddScore,
				record.QA9Option, record.QA9OriginScore, record.QA9AddScore,
				record.QA10Option, record.QA10OriginScore, record.QA10AddScore,
				record.QA11Option, record.QA11OriginScore, record.QA11AddScore,
				record.QA12Option, record.QA12OriginScore, record.QA12AddScore,
				record.QA13Option, record.QA13OriginScore, record.QA13AddScore,
				record.QA14Option, record.QA14OriginScore, record.QA14AddScore,
				record.QA15Option, record.QA15OriginScore, record.QA15AddScore,
				record.QA16Option, record.QA16OriginScore, record.QA16AddScore,
				record.QA17Option, record.QA17OriginScore, record.QA17AddScore,
				record.QA18Option, record.QA18OriginScore, record.QA18AddScore,
				record.QA19Option, record.QA19OriginScore, record.QA19AddScore,
				record.QA20Option, record.QA20OriginScore, record.QA20AddScore,
			}
		)

		// 判斷遊戲類型
		if record.Game == "redpack" {
			newGame = "搖紅包"
		} else if record.Game == "ropepack" {
			newGame = "套紅包"
		} else if record.Game == "whack_mole" {
			newGame = "敲敲樂"
		} else if record.Game == "monopoly" {
			newGame = "鑑定師"
		} else if record.Game == "draw_numbers" {
			newGame = "搖號抽獎"
		} else if record.Game == "lottery" {
			newGame = "遊戲抽獎"
		} else if record.Game == "QA" {
			newGame = "快問快答"
		} else if record.Game == "tugofwar" {
			newGame = "拔河遊戲"
		} else if record.Game == "bingo" {
			newGame = "賓果遊戲"
		} else if record.Game == "3DGachaMachine" {
			newGame = "3D扭蛋機遊戲"
		}

		// 判斷隊伍資訊
		if record.Team == "" {
			// newTeam = "無隊伍"
		} else if record.Team == "left_team" {
			newTeam = record.LeftTeamName
		} else if record.Team == "right_team" {
			newTeam = record.RightTeamName
		}

		// 判斷隊長資訊
		if record.Leader == "" {
		} else if record.Leader == "yes" {
			newLeader = "是"
		} else if record.Leader == "no" {
			newLeader = "否"
		}

		// 判斷隊長資訊
		if record.Mvp == "" {
		} else if record.Mvp == "yes" {
			newMvp = "是"
		} else if record.Mvp == "no" {
			newMvp = "否"
		}

		values := []interface{}{
			// 用戶資訊
			record.Name, record.Number,
			// 遊戲資訊
			record.Title, newGame,
			// 參加資訊
			record.Round, record.Score,
		}

		// 判斷匯出遊戲是否包含快問快答
		if game == "" || strings.Contains(game, "tugofwar") {
			values = append(values, newTeam, newLeader, newMvp)
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

		// 判斷匯出遊戲是否包含快問快答
		if game == "" || strings.Contains(game, "QA") {
			for _, score := range scores {
				values = append(values, score)
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

	fileName = "遊戲紀錄-" + activityID + "-" + gameCN + "-" +
		gameIDCN + "-" + userIDCN + "-" + roundCN
	// 儲存EXCEL
	if err := file.SaveAs(fmt.Sprintf(config.STORE_PATH+"/excel/%s.xlsx", fileName)); err != nil {
		return file, err
	}

	return file, nil
}

// MapToScoreModel map轉換[]MapToScoreModel
func MapToScoreModel(items []map[string]interface{}) []ScoreModel {
	var scores = make([]ScoreModel, 0)
	for _, item := range items {
		var (
			score ScoreModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &score)

		scores = append(scores, score)
	}
	return scores
}

// score.ID, _ = item["id"].(int64)
// score.ActivityID, _ = item["activity_id"].(string)
// score.GameID, _ = item["game_id"].(string)
// score.UserID, _ = item["user_id"].(string)
// score.Round, _ = item["round"].(int64)
// score.Score, _ = item["score"].(int64)
// score.Score2, _ = item["score_2"].(float64)
// score.Team, _ = item["team"].(string)
// score.Leader, _ = item["leader"].(string)
// score.Mvp, _ = item["mvp"].(string)
// score.Rank, _ = item["rank"].(int64)
// score.Correct, _ = item["correct"].(int64)

// 用戶資訊
// score.Name, _ = item["name"].(string)
// score.Avatar, _ = item["avatar"].(string)
// score.Phone, _ = item["phone"].(string)
// score.Email, _ = item["email"].(string)
// score.ExtEmail, _ = item["ext_email"].(string)

// 遊戲資訊
// score.Title, _ = item["title"].(string)
// score.Game, _ = item["game"].(string)
// score.LeftTeamName, _ = item["left_team_name"].(string)
// score.RightTeamName, _ = item["right_team_name"].(string)

// 自定義用戶資料
// score.Ext1, _ = item["ext_1"].(string)
// score.Ext2, _ = item["ext_2"].(string)
// score.Ext3, _ = item["ext_3"].(string)
// score.Ext4, _ = item["ext_4"].(string)
// score.Ext5, _ = item["ext_5"].(string)
// score.Ext6, _ = item["ext_6"].(string)
// score.Ext7, _ = item["ext_7"].(string)
// score.Ext8, _ = item["ext_8"].(string)
// score.Ext9, _ = item["ext_9"].(string)
// score.Ext10, _ = item["ext_10"].(string)
// score.Number, _ = item["number"].(int64)

// 自定義欄位資訊
// score.Ext1Name, _ = item["ext_1_name"].(string)
// score.Ext1Type, _ = item["ext_1_type"].(string)
// score.Ext1Options, _ = item["ext_1_options"].(string)
// score.Ext1Required, _ = item["ext_1_required"].(string)

// score.Ext2Name, _ = item["ext_2_name"].(string)
// score.Ext2Type, _ = item["ext_2_type"].(string)
// score.Ext2Options, _ = item["ext_2_options"].(string)
// score.Ext2Required, _ = item["ext_2_required"].(string)

// score.Ext3Name, _ = item["ext_3_name"].(string)
// score.Ext3Type, _ = item["ext_3_type"].(string)
// score.Ext3Options, _ = item["ext_3_options"].(string)
// score.Ext3Required, _ = item["ext_3_required"].(string)

// score.Ext4Name, _ = item["ext_4_name"].(string)
// score.Ext4Type, _ = item["ext_4_type"].(string)
// score.Ext4Options, _ = item["ext_4_options"].(string)
// score.Ext4Required, _ = item["ext_4_required"].(string)

// score.Ext5Name, _ = item["ext_5_name"].(string)
// score.Ext5Type, _ = item["ext_5_type"].(string)
// score.Ext5Options, _ = item["ext_5_options"].(string)
// score.Ext5Required, _ = item["ext_5_required"].(string)

// score.Ext6Name, _ = item["ext_6_name"].(string)
// score.Ext6Type, _ = item["ext_6_type"].(string)
// score.Ext6Options, _ = item["ext_6_options"].(string)
// score.Ext6Required, _ = item["ext_6_required"].(string)

// score.Ext7Name, _ = item["ext_7_name"].(string)
// score.Ext7Type, _ = item["ext_7_type"].(string)
// score.Ext7Options, _ = item["ext_7_options"].(string)
// score.Ext7Required, _ = item["ext_7_required"].(string)

// score.Ext8Name, _ = item["ext_8_name"].(string)
// score.Ext8Type, _ = item["ext_8_type"].(string)
// score.Ext8Options, _ = item["ext_8_options"].(string)
// score.Ext8Required, _ = item["ext_8_required"].(string)

// score.Ext9Name, _ = item["ext_9_name"].(string)
// score.Ext9Type, _ = item["ext_9_type"].(string)
// score.Ext9Options, _ = item["ext_9_options"].(string)
// score.Ext9Required, _ = item["ext_9_required"].(string)

// score.Ext10Name, _ = item["ext_10_name"].(string)
// score.Ext10Type, _ = item["ext_10_type"].(string)
// score.Ext10Options, _ = item["ext_10_options"].(string)
// score.Ext10Required, _ = item["ext_10_required"].(string)

// score.ExtEmailRequired, _ = item["ext_email_required"].(string)
// score.ExtPhoneRequired, _ = item["ext_phone_required"].(string)
// score.InfoPicture, _ = item["info_picture"].(string)

// 答題記錄
// score.QA1Option, _ = item["qa_1_option"].(string)
// score.QA1OriginScore, _ = item["qa_1_origin_score"].(int64)
// score.QA1AddScore, _ = item["qa_1_add_score"].(int64)

// score.QA2Option, _ = item["qa_2_option"].(string)
// score.QA2OriginScore, _ = item["qa_2_origin_score"].(int64)
// score.QA2AddScore, _ = item["qa_2_add_score"].(int64)

// score.QA3Option, _ = item["qa_3_option"].(string)
// score.QA3OriginScore, _ = item["qa_3_origin_score"].(int64)
// score.QA3AddScore, _ = item["qa_3_add_score"].(int64)

// score.QA4Option, _ = item["qa_4_option"].(string)
// score.QA4OriginScore, _ = item["qa_4_origin_score"].(int64)
// score.QA4AddScore, _ = item["qa_4_add_score"].(int64)

// score.QA5Option, _ = item["qa_5_option"].(string)
// score.QA5OriginScore, _ = item["qa_5_origin_score"].(int64)
// score.QA5AddScore, _ = item["qa_5_add_score"].(int64)

// score.QA6Option, _ = item["qa_6_option"].(string)
// score.QA6OriginScore, _ = item["qa_6_origin_score"].(int64)
// score.QA6AddScore, _ = item["qa_6_add_score"].(int64)

// score.QA7Option, _ = item["qa_7_option"].(string)
// score.QA7OriginScore, _ = item["qa_7_origin_score"].(int64)
// score.QA7AddScore, _ = item["qa_7_add_score"].(int64)

// score.QA8Option, _ = item["qa_8_option"].(string)
// score.QA8OriginScore, _ = item["qa_8_origin_score"].(int64)
// score.QA8AddScore, _ = item["qa_8_add_score"].(int64)

// score.QA9Option, _ = item["qa_9_option"].(string)
// score.QA9OriginScore, _ = item["qa_9_origin_score"].(int64)
// score.QA9AddScore, _ = item["qa_9_add_score"].(int64)

// score.QA10Option, _ = item["qa_10_option"].(string)
// score.QA10OriginScore, _ = item["qa_10_origin_score"].(int64)
// score.QA10AddScore, _ = item["qa_10_add_score"].(int64)

// score.QA11Option, _ = item["qa_11_option"].(string)
// score.QA11OriginScore, _ = item["qa_11_origin_score"].(int64)
// score.QA11AddScore, _ = item["qa_11_add_score"].(int64)

// score.QA12Option, _ = item["qa_12_option"].(string)
// score.QA12OriginScore, _ = item["qa_12_origin_score"].(int64)
// score.QA12AddScore, _ = item["qa_12_add_score"].(int64)

// score.QA13Option, _ = item["qa_13_option"].(string)
// score.QA13OriginScore, _ = item["qa_13_origin_score"].(int64)
// score.QA13AddScore, _ = item["qa_13_add_score"].(int64)

// score.QA14Option, _ = item["qa_14_option"].(string)
// score.QA14OriginScore, _ = item["qa_14_origin_score"].(int64)
// score.QA14AddScore, _ = item["qa_14_add_score"].(int64)

// score.QA15Option, _ = item["qa_15_option"].(string)
// score.QA15OriginScore, _ = item["qa_15_origin_score"].(int64)
// score.QA15AddScore, _ = item["qa_15_add_score"].(int64)

// score.QA16Option, _ = item["qa_16_option"].(string)
// score.QA16OriginScore, _ = item["qa_16_origin_score"].(int64)
// score.QA16AddScore, _ = item["qa_16_add_score"].(int64)

// score.QA17Option, _ = item["qa_17_option"].(string)
// score.QA17OriginScore, _ = item["qa_17_origin_score"].(int64)
// score.QA17AddScore, _ = item["qa_17_add_score"].(int64)

// score.QA18Option, _ = item["qa_18_option"].(string)
// score.QA18OriginScore, _ = item["qa_18_origin_score"].(int64)
// score.QA18AddScore, _ = item["qa_18_add_score"].(int64)

// score.QA19Option, _ = item["qa_19_option"].(string)
// score.QA19OriginScore, _ = item["qa_19_origin_score"].(int64)
// score.QA19AddScore, _ = item["qa_19_add_score"].(int64)

// score.QA20Option, _ = item["qa_20_option"].(string)
// score.QA20OriginScore, _ = item["qa_20_origin_score"].(int64)
// score.QA20AddScore, _ = item["qa_20_add_score"].(int64)

// if _, err := strconv.Atoi(model.Round); err != nil {
// 	return errors.New("錯誤: 遊戲輪次資料發生問題，請輸入有效的輪次資料")
// }

// if isRedis {
// 	err = s.RedisConn.ZSetAdd(config.SCORES_REDIS+gameID, userID, 0)
// 	// s.RedisConn.SetExpire(config.SCORES_REDIS+gameID,
// 	// 	config.REDIS_EXPIRE)
// }

// scoreInt, err := strconv.Atoi(score)
// if err != nil {
// 	return errors.New("錯誤: 遊戲輪次資料發生問題，請輸入有效的輪次資料")
// }

// if err := s.Table(s.Base.TableName).
// 	Where("game_id", "=", model.GameID).
// 	Where("user_id", "=", model.UserID).
// 	Where("round", "=", model.Round).
// 	Update(command.Value{"score": model.Score}); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return errors.New("錯誤: 更新計分表人員資料發生問題，請重新操作")
// }

// UpdateUser 更新計分表人員姓名、頭像資料
// func (s ScoreModel) UpdateUser(userID, name, avatar string) error {
// 	var (
// 		fieldValues = command.Value{"name": name, "avatar": avatar}
// 	)
// 	if err := s.Table(s.Base.TableName).Where("user_id", "=", userID).
// 		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新計分表人員資料發生問題，

// IsUserExist 判斷用戶分數資料是否存在遊戲中
// func (m ScoreModel) IsUserExist(gameID, userID string, round int64) bool {
// 	item, _ := m.Table(m.Base.TableName).Where("game_id", "=", gameID).
// 		Where("user_id", "=", userID).Where("round", "=", round).First()
// 	if item == nil {
// 		return false
// 	}
// 	return true
// }

// EditScoreModel 資料表欄位
// type EditScoreModel struct {
// 	// GameID string `json:"game_id" example:"game_id"`
// 	// UserID string `json:"user_id" example:"user_id"`
// 	// Round  string `json:"round" example:"1"`
// 	// Score  string `json:"score
