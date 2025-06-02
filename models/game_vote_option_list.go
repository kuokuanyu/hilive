package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
	"strings"
)

// GameVoteOptionListModel 資料表欄位
type GameVoteOptionListModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	OptionID   string `json:"option_id" example:"option_id"`
	Name       string `json:"name" example:"name"`
	Avatar     string `json:"avatar" example:"avatar"`
	Leader     string `json:"leader" example:"leader"`
}

// EditGameVoteOptionListModel 資料表欄位
type EditGameVoteOptionListModel struct {
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	OptionID   string `json:"option_id" example:"option_id"`
	Name       string `json:"name" example:"name"`
	Leader     string `json:"leader" example:"leader"`

	People string `json:"people"` // 匯入資料數
}

// DefaultGameVoteOptionListModel 預設GameVoteOptionListModel
func DefaultGameVoteOptionListModel() GameVoteOptionListModel {
	return GameVoteOptionListModel{Base: Base{TableName: config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameVoteOptionListModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameVoteOptionListModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameVoteOptionListModel) SetDbConn(conn db.Connection) GameVoteOptionListModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameVoteOptionListModel) SetRedisConn(conn cache.Connection) GameVoteOptionListModel {
// 	a.RedisConn = conn
// 	return a
// }

// Find 查詢選項名單人員資料
func (a GameVoteOptionListModel) Find(isRedis bool, gameID, optionID, leader string) ([]GameVoteOptionListModel, error) {
	var (
		sql = a.Table(a.Base.TableName).
			Select(
				"activity_game_vote_option_list.id", "activity_game_vote_option_list.activity_id",
				"activity_game_vote_option_list.game_id", "activity_game_vote_option_list.user_id",
				"activity_game_vote_option_list.option_id",
				"activity_game_vote_option_list.leader",

				// 用戶
				"line_users.name", "line_users.avatar",
			).
			LeftJoin(command.Join{
				FieldA:    "activity_game_vote_option_list.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("game_id", "=", gameID).
			OrderBy("id", "asc")
	)

	if optionID != "" { // 查詢特定選項
		sql = sql.Where("option_id", "=", optionID)
	}
	if leader != "" { // 查詢隊長
		sql = sql.Where("leader", "=", leader)
	}

	items, err := sql.All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得選項名單人員資訊，請重新查詢")
	}

	return MapToVoteOptionListModel(items), nil
}

// MapToVoteOptionListModel map轉換[]GameVoteOptionListModel
func MapToVoteOptionListModel(items []map[string]interface{}) []GameVoteOptionListModel {
	var lists = make([]GameVoteOptionListModel, 0)
	for _, item := range items {
		var (
			list GameVoteOptionListModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &list)

		lists = append(lists, list)
	}
	return lists
}

// Adds 批量新增投票選項名單資料
func (a GameVoteOptionListModel) Adds(model EditGameVoteOptionListModel) error {
	// 取得匯入總資料數
	dataAmount, err := strconv.Atoi(model.People)
	if err != nil {
		return errors.New("錯誤: 人數資料發生問題，請重新操作")
	}

	if dataAmount > 0 {
		// 處理匯入的參數資料
		var (
			activityID = make([]string, dataAmount)       // 活動ID
			gameID     = make([]string, dataAmount)       // 遊戲ID
			optionID   = make([]string, dataAmount)       // 選項ID
			userID     = make([]string, dataAmount)       // 用戶ID
			name       = make([]string, dataAmount)       // 用戶名稱
			leader     = strings.Split(model.Leader, ",") // 是否為隊長
		)

		// 處理活動.遊戲.選項ID陣列資料
		for i := range userID {
			activityID[i] = model.ActivityID
			gameID[i] = model.GameID
			optionID[i] = model.OptionID
		}

		if model.UserID == "" && model.Name != "" { // 匯入excel時沒有user_id參數，有name參數
			name = strings.Split(model.Name, ",") // 用戶名稱
			device := make([]string, dataAmount)  // 裝置

			// 處理用戶ID.裝置陣列資料(隨機產生user_id)
			for i := range userID {
				userID[i] = model.ActivityID + "_" + utils.UUID(16)
				device[i] = "vote"
			}

			// 將資料匹量寫入line_users表中
			err = a.Table(config.LINE_USERS_TABLE).BatchInsert(
				dataAmount, "user_id,name,identify,device,activity_id",
				[][]string{userID, name, userID, device, activityID})
			if err != nil {
				return errors.New("錯誤: 匹量新增用戶發生問題(line_users)，請重新操作")
			}
		} else if model.UserID != "" {
			// 從參加人員中匯入時
			userID = strings.Split(model.UserID, ",")
		}

		// 將資料匹量寫入activity_game_vote_option_list表中
		err = a.Table(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).BatchInsert(
			dataAmount, "activity_id,game_id,user_id,option_id,leader",
			[][]string{activityID, gameID, userID, optionID, leader})

		if err != nil {
			return errors.New("錯誤: 匹量新增投票選項名單發生問題(activity_game_vote_option_list)，請重新操作")
		}
	}
	return nil
}

// Add 新增投票選項名單資料(加入資料庫)
func (a GameVoteOptionListModel) Add(model EditGameVoteOptionListModel) error {
	if model.UserID == "" { // 匯入excel時沒有user_id參數
		model.UserID = model.ActivityID + "_" + utils.UUID(16)

		// 將匯入的人員資料寫入line_users表中
		_, err := a.Table(config.LINE_USERS_TABLE).Insert(command.Value{
			"user_id":     model.UserID,
			"name":        model.Name,
			"identify":    model.UserID,
			"device":      "vote",
			"activity_id": model.ActivityID,
		})
		if err != nil {
			return errors.New("錯誤: 新增用戶發生問題(line_users)，請重新操作")
		}
	}

	var (
		fieldValues = command.Value{
			"activity_id": model.ActivityID,
			"game_id":     model.GameID,
			"user_id":     model.UserID,
			"option_id":   model.OptionID,
			// "name":        model.Name,
			"leader": model.Leader,
		}
	)

	if _, err := a.Table(a.TableName).Insert(fieldValues); err != nil {
		return errors.New("錯誤: 新增投票選項名單資料發生問題")
	}

	return nil
}

// Update 更新投票選項名單資料
func (a GameVoteOptionListModel) Update(model EditGameVoteOptionListModel) error {
	var (
		fieldValues = command.Value{"leader": model.Leader}
		fields      = []string{"option_id"}
		values      = []string{model.OptionID}
	)

	if model.Name != "" {
		// 更新用戶名稱資料
		err := a.Table(config.LINE_USERS_TABLE).
			Where("user_id", "=", model.UserID).
			Update(command.Value{
				"name": model.Name,
			})
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 更新用戶發生問題(line_users)，請重新操作")
		}
	}

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	err := a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).
		Where("game_id", "=", model.GameID).
		Where("user_id", "=", model.UserID).
		Update(fieldValues)

	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return err
	}

	return nil
}

// list.ID, _ = item["id"].(int64)
// list.ActivityID, _ = item["activity_id"].(string)
// list.GameID, _ = item["game_id"].(string)
// list.UserID, _ = item["user_id"].(string)
// list.OptionID, _ = item["option_id"].(string)
// list.Name, _ = item["name"].(string)
// list.Avatar, _ = item["avatar"].(string)
// list.Leader, _ = item["leader"].(string)
