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
)

// GameSettingModel 資料表欄位
type GameSettingModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id" example:"activity_id"`

	// 基本設置
	LotteryGameAllow      string `json:"lottery_game_allow" example:"open、close"`          // 遊戲抽獎遊戲是否允許重複中獎
	RedpackGameAllow      string `json:"redpack_game_allow" example:"open、close"`          // 搖紅包遊戲是否允許重複中獎
	RopepackGameAllow     string `json:"ropepack_game_allow" example:"open、close"`         // 套紅包遊戲是否允許重複中獎
	WhackMoleGameAllow    string `json:"whack_mole_game_allow" example:"open、close"`       // 敲敲樂遊戲是否允許重複中獎
	MonopolyGameAllow     string `json:"monopoly_game_allow" example:"open、close"`         // 抓偽鈔遊戲是否允許重複中獎
	QAGameAllow           string `json:"qa_game_allow" example:"open、close"`               // 快問快答遊戲是否允許重複中獎
	DrawNumbersGameAllow  string `json:"draw_numbers_game_allow" example:"open、close"`     // 搖號抽獎遊戲是否允許重複中獎
	TugofwarGameAllow     string `json:"tugofwar_game_allow" example:"open、close"`         // 拔河遊戲是否允許重複中獎
	BingoGameAllow        string `json:"bingo_game_allow" example:"open、close"`            // 賓果遊戲是否允許重複中獎
	GachaMachineGameAllow string `json:"3d_gacha_machine_game_allow" example:"open、close"` // 扭蛋機遊戲是否允許重複中獎
	VoteGameAllow         string `json:"vote_game_allow" example:"open、close"`             // 投票遊戲是否允許重複中獎
	AllGameAllow          string `json:"all_game_allow" example:"open、close"`              // 所有遊戲是否允許重複中獎
}

// EditGameSettingModel 資料表欄位
type EditGameSettingModel struct {
	UserID     string `json:"user_id" example:"user_id"`
	ActivityID string `json:"activity_id" example:"activity_id"`

	// 基本設置
	LotteryGameAllow      string `json:"lottery_game_allow" example:"open、close"`          // 遊戲抽獎遊戲是否允許重複中獎
	RedpackGameAllow      string `json:"redpack_game_allow" example:"open、close"`          // 搖紅包遊戲是否允許重複中獎
	RopepackGameAllow     string `json:"ropepack_game_allow" example:"open、close"`         // 套紅包遊戲是否允許重複中獎
	WhackMoleGameAllow    string `json:"whack_mole_game_allow" example:"open、close"`       // 敲敲樂遊戲是否允許重複中獎
	MonopolyGameAllow     string `json:"monopoly_game_allow" example:"open、close"`         // 抓偽鈔遊戲是否允許重複中獎
	QAGameAllow           string `json:"qa_game_allow" example:"open、close"`               // 快問快答遊戲是否允許重複中獎
	DrawNumbersGameAllow  string `json:"draw_numbers_game_allow" example:"open、close"`     // 搖號抽獎遊戲是否允許重複中獎
	TugofwarGameAllow     string `json:"tugofwar_game_allow" example:"open、close"`         // 拔河遊戲是否允許重複中獎
	BingoGameAllow        string `json:"bingo_game_allow" example:"open、close"`            // 賓果遊戲是否允許重複中獎
	GachaMachineGameAllow string `json:"3d_gacha_machine_game_allow" example:"open、close"` // 扭蛋機遊戲是否允許重複中獎
	VoteGameAllow         string `json:"vote_game_allow" example:"open、close"`             // 投票遊戲是否允許重複中獎
	AllGameAllow          string `json:"all_game_allow" example:"open、close"`              // 所有遊戲是否允許重複中獎

	// Token string `json:"token" example:"token"`
}

// DefaultGameSettingModel 預設GameSettingModel
func DefaultGameSettingModel() GameSettingModel {
	return GameSettingModel{Base: Base{TableName: config.ACTIVITY_GAME_SETTING_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameSettingModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameSettingModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameSettingModel) SetDbConn(conn db.Connection) GameSettingModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameSettingModel) SetRedisConn(conn cache.Connection) GameSettingModel {
// 	a.RedisConn = conn
// 	return a
// }

// Find 查詢遊戲基本設置資訊
func (a GameSettingModel) Find(activityID string) (GameSettingModel, error) {
	item, err := a.Table(a.Base.TableName).
		Where("activity_id", "=", activityID).First()
	if err != nil {
		return GameSettingModel{}, errors.New("錯誤: 無法取得遊戲基本設置資訊，請重新查詢")
	}

	return a.MapToModel(item), nil
}

// Add 新增遊戲基本設置資料
func (a GameSettingModel) Add(model EditGameSettingModel) error {
	if _, err := a.Table(a.Base.TableName).Insert(command.Value{
		"activity_id": model.ActivityID,
	}); err != nil {
		return errors.New("錯誤: 無法新增遊戲基本設置資料，請重新操作")
	}
	return nil
}

// Update 更新遊戲基本設置資料
func (a GameSettingModel) Update(isRedis bool, model EditGameSettingModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"lottery_game_allow", "redpack_game_allow",
			"ropepack_game_allow", "whack_mole_game_allow", "monopoly_game_allow",
			"qa_game_allow", "draw_numbers_game_allow", "all_game_allow",
			"tugofwar_game_allow", "bingo_game_allow", "3d_gacha_machine_game_allow",
			"vote_game_allow",
		}
	)

	// 更新遊戲基本設置資料
	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			if val != "open" && val != "close" {
				return errors.New("錯誤: 是否允許重複中獎資料發生問題，請輸入有效的資料")
			}

			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		if err := a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).
			Update(fieldValues); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
	}

	if isRedis {
		// 清除所有遊戲redis資訊，先查詢該活動下的所有遊戲
		games, err := DefaultGameModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			FindAll(model.ActivityID, "")
		if err != nil {
			return err
		}

		deleteKeys := []string{
			config.GAME_REDIS,              // 遊戲設置
			config.GAME_TYPE_REDIS,         // 遊戲種類
			config.SCORES_REDIS,            // 分數
			config.SCORES_2_REDIS,          // 第二分數
			config.CORRECT_REDIS,           // 答對題數
			config.QA_REDIS,                // 快問快答題目資訊
			config.QA_RECORD_REDIS,         // 快問快答答題紀錄
			config.WINNING_STAFFS_REDIS,    // 中獎人員
			config.NO_WINNING_STAFFS_REDIS, // 未中獎人員
			// config.DRAW_NUMBERS_WINNING_STAFFS_REDIS, // 活動下所有場次搖號抽獎的中獎人員 (目前註解保留)
			config.GAME_TEAM_REDIS,                       // 遊戲隊伍資訊，HASH
			config.GAME_BINGO_NUMBER_REDIS,               // 紀錄抽過的號碼，LIST
			config.GAME_BINGO_USER_REDIS,                 // 賓果中獎人員，ZSET
			config.GAME_BINGO_USER_NUMBER,                // 紀錄玩家的號碼排序，HASH
			config.GAME_BINGO_USER_GOING_BINGO,           // 紀錄玩家是否即將賓果，HASH
			config.GAME_ATTEND_REDIS,                     // 遊戲人數資訊，SET
			config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS,  // 拔河遊戲左隊人數資訊，SET
			config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS, // 拔河遊戲右隊人數資訊，SET
		}

		publishChannels := []string{
			config.CHANNEL_GAME_REDIS,
			config.CHANNEL_GUEST_GAME_STATUS_REDIS,
			config.CHANNEL_GAME_BINGO_NUMBER_REDIS,
			config.CHANNEL_QA_REDIS,
			config.CHANNEL_GAME_TEAM_REDIS,
			config.CHANNEL_GAME_EDIT_TIMES_REDIS,
			config.CHANNEL_WINNING_STAFFS_REDIS,
			config.CHANNEL_GAME_BINGO_USER_NUMBER,
			config.CHANNEL_SCORES_REDIS,
		}

		for _, game := range games {
			// 清除遊戲 redis 資訊
			for _, key := range deleteKeys {
				a.RedisConn.DelCache(key + game.GameID)
			}

			// 發佈 WebSocket 訊息
			for _, channel := range publishChannels {
				a.RedisConn.Publish(channel+game.GameID, "修改資料")
			}
		}

	}
	return nil
}

// MapToModel 將值設置至GameSettingModel
func (a GameSettingModel) MapToModel(m map[string]interface{}) GameSettingModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	return a
}

// a.ID, _ = m["id"].(int64)
// a.ActivityID, _ = m["activity_id"].(string)

// 遊戲基本設置
// a.LotteryGameAllow, _ = m["lottery_game_allow"].(string)
// a.RedpackGameAllow, _ = m["redpack_game_allow"].(string)
// a.RopepackGameAllow, _ = m["ropepack_game_allow"].(string)
// a.WhackMoleGameAllow, _ = m["whack_mole_game_allow"].(string)
// a.MonopolyGameAllow, _ = m["monopoly_game_allow"].(string)
// a.QAGameAllow, _ = m["qa_game_allow"].(string)
// a.DrawNumbersGameAllow, _ = m["draw_numbers_game_allow"].(string)
// a.TugofwarGameAllow, _ = m["tugofwar_game_allow"].(string)
// a.BingoGameAllow, _ = m["bingo_game_allow"].(string)
// a.GachaMachineGameAllow, _ = m["3d_gacha_machine_game_allow"].(string)
// a.VoteGameAllow, _ = m["vote_game_allow"].(string)
// a.AllGameAllow, _ = m["all_game_allow"].(string)

// values = []string{model.LotteryGameAllow, model.RedpackGameAllow,
// 	model.RopepackGameAllow, model.WhackMoleGameAllow, model.MonopolyGameAllow,
// 	model.QAGameAllow, model.DrawNumbersGameAllow, model.AllGameAllow,
// 	model.TugofwarGameAllow, model.BingoGameAllow, model.GachaMachineGameAllow,
// 	model.VoteGameAllow}

// for i, value := range values {
// 	if value != "" {
// 		if value != "open" && value != "close" {
// 			return errors.New("錯誤: 是否允許重複中獎資料發生問題，請輸入有效的資料")
// 		}
// 		fieldValues[fields[i]] = value
// 	}
// }

// for _, game := range games {
// 	// 清除遊戲redis資訊(並重新開啟遊戲頁面)
// 	a.RedisConn.DelCache(config.GAME_REDIS + game.GameID)              // 遊戲設置
// 	a.RedisConn.DelCache(config.GAME_TYPE_REDIS + game.GameID)         // 遊戲種類
// 	a.RedisConn.DelCache(config.SCORES_REDIS + game.GameID)            // 分數
// 	a.RedisConn.DelCache(config.SCORES_2_REDIS + game.GameID)          // 第二分數
// 	a.RedisConn.DelCache(config.CORRECT_REDIS + game.GameID)           // 答對題數
// 	a.RedisConn.DelCache(config.QA_REDIS + game.GameID)                // 快問快答題目資訊
// 	a.RedisConn.DelCache(config.QA_RECORD_REDIS + game.GameID)         // 快問快答答題紀錄
// 	a.RedisConn.DelCache(config.WINNING_STAFFS_REDIS + game.GameID)    // 中獎人員
// 	a.RedisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + game.GameID) // 未中獎人員
// 	// a.RedisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + game.ActivityID) // 活動下所有場次搖號抽獎的中獎人員
// 	a.RedisConn.DelCache(config.GAME_TEAM_REDIS + game.GameID)             // 遊戲隊伍資訊，HASH
// 	a.RedisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + game.GameID)     // 紀錄抽過的號碼，LIST
// 	a.RedisConn.DelCache(config.GAME_BINGO_USER_REDIS + game.GameID)       // 賓果中獎人員，ZSET
// 	a.RedisConn.DelCache(config.GAME_BINGO_USER_NUMBER + game.GameID)      // 紀錄玩家的號碼排序，HASH
// 	a.RedisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + game.GameID) // 紀錄玩家是否即將賓果，HASH

// 	a.RedisConn.DelCache(config.GAME_ATTEND_REDIS + game.GameID)                     // 遊戲人數資訊，SET
// 	a.RedisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + game.GameID)  // 拔河遊戲左隊人數資訊，SET
// 	a.RedisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + game.GameID) // 拔河遊戲右隊人數資訊，SET

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_QA_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+game.GameID, "修改資料")

// 	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 	a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+game.GameID, "修改資料")

// }
