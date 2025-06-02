package models

import (
	"errors"
	"hilive/modules/config"
)

// FindTeam 查詢隊伍資訊，從reids取得
func (a GameModel) FindTeam(isRedis bool, gameID string) (GameModel, error) {
	if isRedis {
		var (
			leftTeamPlayers  = make([]string, 0)
			rightTeamPlayers = make([]string, 0)
			// leftTeamLeader, rightTeamLeader string
		)

		// 查詢redis裡隊伍資訊(hash，隊長、leader、win)
		dataMap, err := a.RedisConn.HashGetAllCache(config.GAME_TEAM_REDIS + gameID)
		if err != nil {
			return GameModel{}, errors.New("錯誤: 取得隊伍快取資料發生問題")
		}

		// 查詢redis裡左隊隊員資訊(set)
		leftTeamPlayers, err = a.RedisConn.SetGetMembers(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)
		if err != nil {
			return GameModel{}, errors.New("錯誤: 取得左隊隊伍人員快取資料發生問題")
		}

		// 查詢redis裡右隊隊員資訊(set)
		rightTeamPlayers, err = a.RedisConn.SetGetMembers(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID)
		if err != nil {
			return GameModel{}, errors.New("錯誤: 取得右隊隊伍人員快取資料發生問題")
		}

		a.LeftTeamPlayers = leftTeamPlayers
		a.RightTeamPlayers = rightTeamPlayers
		a.LeftTeamLeader = dataMap["left_team_leader"]
		a.RightTeamLeader = dataMap["right_team_leader"]
		a.WinTeam = dataMap["win_team"]

		// 解碼
		// json.Unmarshal([]byte(dataMap["left_team_players"]), &leftTeamPlayers)   // 左方隊伍玩家
		// json.Unmarshal([]byte(dataMap["right_team_players"]), &rightTeamPlayers) // 右方隊伍玩家
		// json.Unmarshal([]byte(dataMap["left_team_leader"]), &leftTeamLeader)     // 左方隊長
		// json.Unmarshal([]byte(dataMap["right_team_leader"]), &rightTeamLeader)   // 右方隊長
		// if len(leftTeamPlayers) == 0 {
		// 	leftTeamPlayers = make([]string, 0)
		// fmt.Println(leftTeamPlayers)
		// }
		// if len(rightTeamPlayers) == 0 {
		// 	rightTeamPlayers = make([]string, 0)
		// fmt.Println(rightTeamPlayers)
		// }

		// fmt.Println("左方隊伍玩家: ", leftTeamPlayers)
		// fmt.Println("右方隊伍玩家: ", rightTeamPlayers)
	}

	return a, nil
}

// AddPlayer 更新redis中的隊員資訊
func (a GameModel) AddPlayer(isRedis bool, gameID string, userID string,
	team string) error {
	if isRedis {
		var (
			key string // redis key
			// players = make([]string, 0)
			// field   = team + "_players"
		)

		if team == "left_team" {
			// log.Println("加入遊戲人員，左隊")
			key = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲左隊人數資訊，SET
		} else if team == "right_team" {
			// log.Println("加入遊戲人員，右隊")
			key = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
		}

		// log.Println("加入前隊伍人數: ", a.RedisConn.SetCard(key))

		// 將玩家資料加入隊伍(redis)
		a.RedisConn.SetAdd([]interface{}{key, userID})

		// log.Println("加入後隊伍人數: ", a.RedisConn.SetCard(key))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資訊)
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}

	return nil
}

// UpdateTeamPlayers 更新redis中的雙方隊伍玩家資訊
func (a GameModel) UpdateTeamPlayers(isRedis bool, gameID string,
	leftTeamPlayers []string, rightTeamPlayers []string) error {
	if isRedis {
		var (
			leftTeamRedisKey  = []interface{}{config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID}  // 左隊reids key
			rightTeamRedisKey = []interface{}{config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID} // 右隊reids key
			err1, err2        error
		)

		// 左隊
		for _, userID := range leftTeamPlayers {
			leftTeamRedisKey = append(leftTeamRedisKey, userID)
		}

		// 右隊
		for _, userID := range rightTeamPlayers {
			rightTeamRedisKey = append(rightTeamRedisKey, userID)
		}

		// 先清除舊的雙分隊伍資訊
		a.RedisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)  // 左隊
		a.RedisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID) // 右隊

		// 將新的隊伍人員資訊寫入redis中
		if len(leftTeamRedisKey) > 1 {
			err1 = a.RedisConn.SetAdd(leftTeamRedisKey) // 左隊
		}
		if len(rightTeamRedisKey) > 1 {
			err2 = a.RedisConn.SetAdd(rightTeamRedisKey) // 右隊
		}
		if err1 != nil || err2 != nil {
			return errors.New("錯誤: 更新雙方隊伍玩家快取資料發生問題")
		}

		// log.Println("遊戲頁面端調整隊伍後左隊人數: ", a.RedisConn.SetCard(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID))
		// log.Println("遊戲頁面端調整隊伍後右隊人數: ", a.RedisConn.SetCard(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資訊)
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}

	return nil
}

// UpdateTeamLeader 更新redis中的雙方隊長資訊
func (a GameModel) UpdateTeamLeader(isRedis bool, gameID string,
	leftTeamLeader string, rightTeamLeader string) error {
	if isRedis {
		// 更新redis裡隊伍資訊
		err := a.RedisConn.HashMultiSetCache([]interface{}{
			config.GAME_TEAM_REDIS + gameID,
			"left_team_leader", leftTeamLeader, // 左方隊長
			"right_team_leader", rightTeamLeader, // 右方隊長
		})
		if err != nil {
			return errors.New("錯誤: 更新隊長快取資料發生問題")
		}

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")
	}

	return nil
}

// ChangePlayer 更新redis中的雙方隊員資訊
func (a GameModel) ChangePlayer(isRedis bool, gameID string, userID string,
	newTeam string) error {
	if isRedis {
		// log.Println("平台端修改雙方隊伍人員資料")

		var (
			newTeamRedisKey, oldTeamRedisKey string // 隊伍redis key
		)

		if newTeam == "left_team" {
			newTeamRedisKey = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID  // 拔河遊戲左隊人數資訊，SET
			oldTeamRedisKey = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
		} else if newTeam == "right_team" {
			newTeamRedisKey = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
			oldTeamRedisKey = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID  // 拔河遊戲左隊人數資訊，SET
		}

		// 處理雙方隊伍玩家資訊(redis)
		// 將玩家加入新的隊伍中(redis)
		a.RedisConn.SetAdd([]interface{}{newTeamRedisKey, userID})

		// 將玩家從舊的隊伍中清除(redis)
		a.RedisConn.SetRem([]interface{}{oldTeamRedisKey, userID})

		// log.Println("平台端修改後左隊人數: ", a.RedisConn.SetCard(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID))
		// log.Println("平台端修改後右隊人數: ", a.RedisConn.SetCard(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資訊)
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")
	}

	return nil
}

// UpdateLeader 更新redis中的隊長資訊(其中一隊)
func (a GameModel) UpdateLeader(isRedis bool, gameID string, userID string,
	team string) error {
	// fmt.Println("當隊長")
	if isRedis {
		// 更新redis裡隊伍資訊
		err := a.RedisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID, team+"_leader", userID)
		if err != nil {
			return errors.New("錯誤: 更新隊長快取資料發生問題")
		}

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")
	}

	return nil
}

// DeletePlayer 清除玩家隊伍陣列資訊(redis)
func (a GameModel) DeletePlayer(isRedis bool, gameID string,
	team string, userID string) {
	if isRedis {
		var (
			key string // redis key
		)

		if team == "left_team" {
			// log.Println("刪除遊戲人員，左隊")
			key = config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲左隊人數資訊，SET
		} else if team == "right_team" {
			// log.Println("刪除遊遊戲人員，右隊")
			key = config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID // 拔河遊戲右隊人數資訊，SET
		}

		// log.Println("平台端刪除前隊伍人數: ", a.RedisConn.SetCard(key))

		// 清除隊伍中的玩家資料(redis)
		a.RedisConn.SetRem([]interface{}{key, userID})

		// log.Println("平台端刪除後隊伍人數: ", a.RedisConn.SetCard(key))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(隊伍資訊)
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}

	// return
}

// DeleteLeader 清除隊長資訊(redis)
func (a GameModel) DeleteLeader(isRedis bool, gameID string, team string) {
	if isRedis {
		// var (
		// 	field       = team + "_leader"
		// )

		// 將redis中的隊長資訊清空
		a.RedisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID, team+"_leader", "")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")
	}

	// return
}

// UpdateTeamPeople 更新隊伍人數資料(資料庫.redis)，遞增、遞減、歸零
// func (a GameModel) UpdateTeamPeople(isRedis bool, gameID string, team string,
// sqlValue int64, redisValue int64) error {
// field := team + "_game_attend"

// if sqlValue == 1 {
// log.Println("資料庫不處理遞增隊伍資料")
// 遞增隊伍人數
// if err := a.Table(config.ACTIVITY_GAME_TABLE).
// 	WhereRaw(fmt.Sprintf("`game_id` = ? and `%s` < `people`", field), gameID).
// 	Update(command.Value{field: field + " + 1"}); err != nil {
// 	if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 該隊伍遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝")
// 	}
// 	return err
// }
// } else if sqlValue == -1 {
// log.Println("資料庫不處理遞減隊伍資料")
// 遞減人員
// if err := a.Table(config.ACTIVITY_GAME_TABLE).
// 	WhereRaw(fmt.Sprintf("`game_id` = ? and `%s` <= `people`", field), gameID).
// 	Update(command.Value{field: field + " - 1"}); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return errors.New("錯誤: 減少隊伍人數發生錯誤，請重新操作")
// }
// } else if sqlValue == 0 {
// log.Println("資料庫不處理歸零隊伍資料")
// 左右方隊伍人數歸零
// if err := a.Table(config.ACTIVITY_GAME_TABLE).
// 	Where("game_id", "=", gameID).
// 	Update(command.Value{
// 		"left_team_game_attend":  0,
// 		"right_team_game_attend": 0}); err != nil &&
// 	err.Error() == "錯誤: 無更新任何資料，請重新操作" {
// 	return errors.New("錯誤: 歸零隊伍人數發生錯誤(資料表)，請重新操作")
// }
// }

// 修改redis中的遊戲資訊
// if isRedis {
// 	if redisValue == 0 {
// 		a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "left_team_game_attend", 0)  // 歸零
// 		a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "right_team_game_attend", 0) // 歸零

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
// 	} else if redisValue == 1 {
// 		a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID, field) // 遞增

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
// 	} else if redisValue == -1 {
// 		people := a.RedisConn.HashDecrCache(config.GAME_REDIS+gameID, field) // 遞減
// 		if people < 0 {
// 			// 人數不能為負，歸零
// 			// fmt.Println("負的，歸零")
// 			a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, field, 0)
// 		}

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// 		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
// 	}
// }

// return nil
// }

// UpdateBothTeamPeople 更新雙方隊伍人數資料(資料庫.redis)
// func (a GameModel) UpdateBothTeamPeople(isRedis bool, gameID string,
// leftTeamAttend int64, rightTeamAttend int64) error {

// 更新資料表雙方隊伍人數
// if err := a.Table(config.ACTIVITY_GAME_TABLE).
// 	Where("game_id", "=", gameID).
// 	Update(command.Value{
// 		"left_team_game_attend":  leftTeamAttend,
// 		"right_team_game_attend": rightTeamAttend}); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return errors.New("錯誤: 更新雙方隊伍人數發生錯誤(資料表)，請重新操作")
// }

// 更新redis雙方隊伍人數
// err1 := a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "left_team_game_attend", leftTeamAttend)
// err2 := a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "right_team_game_attend", rightTeamAttend)
// if err1 != nil || err2 != nil {
// 	return errors.New("錯誤: 更新雙方隊伍人數發生錯誤(redis)，請重新操作")
// }

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
// a.RedisConn.Publish(config.CH
