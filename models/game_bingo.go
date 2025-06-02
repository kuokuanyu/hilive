package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/config"
	"hilive/modules/utils"
	"strconv"
)

// -----判斷是否即將中獎相關-----start

// UpdateBingoRound 更新玩家賓果的回合數(redis)
func (s ScoreModel) UpdateGoingBingo(isRedis bool, gameID, userID string, goingBingo bool) (err error) {
	if isRedis {
		// fmt.Println("有接收到訊息嗎?: ", goingBingo)
		if err = s.RedisConn.HashSetCache(config.GAME_BINGO_USER_GOING_BINGO+gameID,
			userID, strconv.FormatBool(goingBingo)); err != nil {
			return errors.New("錯誤: 更新玩家是否即將賓果資料發生問題(redis)，請重新操作")
		}
	}
	return
}

// -----判斷是否即將中獎相關-----end

// -----抽獎號碼相關-----start

// FindBingoNumbers 查詢賓果遊戲所有抽獎號碼(redis)
func (a GameModel) FindBingoNumbers(isRedis bool, gameID string) ([]int64, error) {
	var numbers = make([]int64, 0)

	if isRedis {
		// 取得賓果遊戲所有抽獎號碼(bingo_number_遊戲ID)
		numbersStr, _ := a.RedisConn.ListRange(config.GAME_BINGO_NUMBER_REDIS+gameID, 0, 0)

		for _, numberStr := range numbersStr {
			number, _ := strconv.Atoi(numberStr)
			numbers = append(numbers, int64(number))
		}
	}
	return numbers, nil
}

// UpdateBingoNumber 更新賓果遊戲抽獎號碼資料(redis)
func (a GameModel) UpdateBingoNumber(isRedis bool, gameID string, numbers []int64) error {
	var params = []interface{}{config.GAME_BINGO_NUMBER_REDIS + gameID}
	for _, number := range numbers {
		params = append(params, number)
	}

	if isRedis {
		// 更新賓果遊戲抽獎號碼(bingo_number_遊戲ID)
		a.RedisConn.ListMultiRPush(params)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+gameID, "修改資料")
	}
	return nil
}

// -----抽獎號碼相關-----start

// -----號碼排序相關-----start

// FindUserAmount 查詢賓果遊戲完成選號人數(redis)
func (a GameModel) FindUserAmount(isRedis bool, gameID string) (int64, error) {
	var amount int64

	if isRedis {
		// 玩家的號碼排序(bingo_user_number_遊戲ID，json編碼)
		datas, _ := a.RedisConn.HashGetAllCache(config.GAME_BINGO_USER_NUMBER + gameID)
		amount = int64(len(datas))
	}
	return amount, nil
}

// FindUserNumber 查詢玩家的號碼排序(redis)
func (a GameModel) FindUserNumber(isRedis bool, gameID string, userID string) ([]int64, error) {
	var numbers = make([]int64, 0)

	if isRedis {
		// 玩家的號碼排序(bingo_user_number_遊戲ID，json編碼)
		data, _ := a.RedisConn.HashGetCache(config.GAME_BINGO_USER_NUMBER+gameID, userID)
		json.Unmarshal([]byte(data), &numbers)
	}
	return numbers, nil
}

// UpdateUserNumber 更新玩家的號碼排序(redis)
func (a GameModel) UpdateUserNumber(isRedis bool, gameID string, userID string, numbers []int64) error {

	if isRedis {
		// 更新玩家的號碼排序(bingo_user_number_遊戲ID，json編碼)
		a.RedisConn.HashSetCache(config.GAME_BINGO_USER_NUMBER+gameID, userID,
			utils.JSON(numbers))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, "修改資料")
	}
	return nil
}

// -----號碼排序相關-----end

// -----回合數相關-----start

// UpdateBingoRound 更新玩家賓果的回合數(redis)
func (s ScoreModel) UpdateBingoRound(isRedis bool, gameID, userID string, score int64) (err error) {
	if isRedis {
		if err = s.RedisConn.ZSetAddInt(config.GAME_BINGO_USER_REDIS+gameID,
			userID, score); err != nil {
			return errors.New("錯誤: 更新玩家賓果的回合數資料發生問題(redis)，請重新操作")
		}
	}
	return
}

// UpdateBingoRound 更新賓果遊戲進行回合數資料(redis)
func (a GameModel) UpdateBingoRound(isRedis bool, gameID, round string) error {
	// var roundValue string
	// if round == "1" {
	// 	// 新的一輪，歸零
	// 	roundValue = "1"
	// } else if round == "+1" {
	// 	// 遞增
	// 	roundValue = "qa_round + 1"
	// }

	// if err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
	// 	Where("game_id", "=", gameID).
	// 	// Where(round, "<", "qa_round").
	// 	Update(command.Value{
	// 		"qa_round": roundValue,
	// 	}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return err
	// }

	// 修改redis中的遊戲資訊
	if isRedis {
		if round == "0" {
			// 新的一輪，歸零
			a.RedisConn.HashSetCache(config.GAME_REDIS+gameID,
				"bingo_round", round)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		} else if round == "+1" {
			// 遞增
			a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID,
				"bingo_round")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		}
	}
	return nil
}

// -----回合數相關-----end

// -----賓果人員相關-----start

// FindBingoUser 查詢賓果人員(redis)
func (m ScoreModel) FindBingoUser(isRedis bool, gameID string) ([]ScoreModel, error) {
	var (
		scores = make([]ScoreModel, 0)
		err    error
	)

	if isRedis {
		var users []string
		// 取得已賓果玩家資料
		if users, err = m.RedisConn.ZSetRange(config.GAME_BINGO_USER_REDIS+gameID, 0, 0); err != nil {
			return scores, errors.New("錯誤: 從redis中取得分數由低到高的玩家資訊發生問題")
		}

		// 處理分數排名資訊(只取得分數大於0的資料)
		for _, userID := range users {
			var (
				score  int64
				score2 int64
			)

			score = m.RedisConn.ZSetIntScore(config.SCORES_REDIS+gameID, userID)           // 賓果連線數資料
			score2 = m.RedisConn.ZSetIntScore(config.GAME_BINGO_USER_REDIS+gameID, userID) // 完成賓果的回合數

			if score2 > 0 {
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
					Score2: float64(score2),
				})
			}
		}
	}
	return scores, nil
}

// -----賓果人員相關-----end

// -----賓果分數相關-----start

// FindBingoScore 查詢賓果遊戲所有人員分數資料、即將賓果人員資料(從redis中取得分數由高至低的玩家資訊)
func (m ScoreModel) FindBingoScore(isRedis bool, gameID string, bingoLine int64) ([]ScoreModel, []ScoreModel, error) {
	var (
		scores   = make([]ScoreModel, 0)
		PKStaffs = make([]ScoreModel, 0)
		err      error
	)

	if isRedis {
		var (
			users           []string
			goingBingoUsers map[string]string
		)
		if users, err = m.RedisConn.ZSetRevRange(config.SCORES_REDIS+gameID, 0, 0); err != nil {
			return scores, PKStaffs, errors.New("錯誤: 從redis中取得分數由高至低的玩家資訊發生問題")
		}

		// 玩家是否即將賓果
		if goingBingoUsers, err = m.RedisConn.HashGetAllCache(config.GAME_BINGO_USER_GOING_BINGO + gameID); err != nil {
			return scores, PKStaffs, errors.New("錯誤: 從redis中取得玩家是否即將賓果資訊發生問題")
		}

		// 處理分數排名資訊
		for _, userID := range users {
			var (
				score  int64
				score2 int64
			)
			score = m.RedisConn.ZSetIntScore(config.SCORES_REDIS+gameID, userID)           // 賓果連線數資料
			score2 = m.RedisConn.ZSetIntScore(config.GAME_BINGO_USER_REDIS+gameID, userID) // 完成賓果的回合數

			// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
			user, err := DefaultLineModel().
				SetConn(m.DbConn, m.RedisConn, m.MongoConn).
				Find(true, "", "user_id", userID)
			if err != nil {
				return scores, PKStaffs, errors.New("錯誤: 無法取得用戶資訊")
			}
			// fmt.Println("玩家是否即將中獎: ", goingBingoUsers[userID], goingBingoUsers[userID] == "true")
			// 判斷是否即將賓果
			if goingBingoUsers[userID] == "true" {
				PKStaffs = append(PKStaffs, ScoreModel{
					ID:     user.ID,
					UserID: userID,
					Name:   user.Name,
					Avatar: user.Avatar,
					Score:  score,
					Score2: float64(score2),
				})
			}

			scores = append(scores, ScoreModel{
				ID:     user.ID,
				UserID: userID,
				Name:   user.Name,
				Avatar: user.Avatar,
				Score:  score,
				Score2: float64(score2),
			})
		}
	}
	return scores, PKStaffs, nil
}

// -----賓果分數相關-----end
