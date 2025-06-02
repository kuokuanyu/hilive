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
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
)

// PrizeModel 資料表欄位
type PrizeModel struct {
	Base          `json:"-"`
	ID            int64  `json:"id"`
	UserID        string `json:"user_id" example:"user_id"`
	ActivityID    string `json:"activity_id" example:"activity_id"`
	GameID        string `json:"game_id" example:"game_id"`
	PrizeID       string `json:"prize_id" example:"prize_id"`
	Game          string `json:"game" example:"game name"`
	PrizeName     string `json:"prize_name" example:"prize name"`
	PrizeType     string `json:"prize_type" example:"prize type"`
	PrizePicture  string `json:"prize_picture" example:"prize picture"`
	PrizeAmount   int64  `json:"prize_amount" example:"prize amount"`
	PrizeRemain   int64  `json:"prize_remain" example:"prize remain"`
	PrizePrice    int64  `json:"prize_price" example:"prize prize"`
	PrizeMethod   string `json:"prize_method" example:"prize method"`
	PrizePassword string `json:"prize_password" example:"prize password"`
	TeamType      string `json:"team_type" example:"win、lose"`

	ActivityName string `json:"activity_name" example:"activity_name"`
}

// NewPrizeModel 資料表欄位
// type NewPrizeModel struct {
// 	UserID       string `json:"user_id" example:"user_id"`
// 	ActivityID   string `json:"activity_id" example:"activity_id"`
// 	GameID       string `json:"game_id" example:"game_id"`
// 	PrizeName    string `json:"prize_name" example:"prize name"`
// 	PrizeType    string `json:"prize_type" example:"first、second、third、general、thanks"`
// 	PrizePicture string `json:"prize_picture" example:"https://..."`
// 	PrizeAmount  string `json:"prize_amount" example:"30"`
// 	// PrizeRemain   int64  `json:"prize_remain" example:"30"`
// 	PrizePrice    string `json:"prize_price" example:"100"`
// 	PrizeMethod   string `json:"prize_method" example:"site、mail、thanks"`
// 	PrizePassword string `json:"prize_password" example:"prize password(最多為8個字元)"`
// 	Token         string `json:"token" example:"token"`
// 	TeamType      string `json:"team_type" example:"win、lose"`
// }

// EditPrizeModel 資料表欄位
type EditPrizeModel struct {
	UserID       string `json:"user_id" example:"user_id"`
	ActivityID   string `json:"activity_id" example:"activity_id"`
	GameID       string `json:"game_id" example:"game_id"`
	PrizeID      string `json:"prize_id" example:"prize_id"`
	PrizeName    string `json:"prize_name" example:"prize name"`
	PrizeType    string `json:"prize_type" example:"first、second、third、general、thanks"`
	PrizePicture string `json:"prize_picture" example:"https://..."`
	PrizeAmount  string `json:"prize_amount" example:"30"`
	// PrizeRemain   int64  `json:"prize_remain" example:"30"`
	PrizePrice    string `json:"prize_price" example:"100"`
	PrizeMethod   string `json:"prize_method" example:"site、mail、thanks"`
	PrizePassword string `json:"prize_password" example:"prize password(最多為8個字元)"`
	// Token         string `json:"token" example:"token"`
	TeamType string `json:"team_type" example:"win、lose"`
}

// DefaultPrizeModel 預設PrizeModel
func DefaultPrizeModel() PrizeModel {
	return PrizeModel{Base: Base{TableName: config.ACTIVITY_PRIZE_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a PrizeModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) PrizeModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a PrizeModel) SetDbConn(conn db.Connection) PrizeModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a PrizeModel) SetRedisConn(conn cache.Connection) PrizeModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (a PrizeModel) SetMongoConn(conn mongo.Connection) PrizeModel {
// 	a.MongoConn = conn
// 	return a
// }

// FindPrize 查詢獎品資訊(單個)
func (a PrizeModel) FindPrize(isRedis bool, prizeID string) (PrizeModel, error) {
	var (
		prize PrizeModel
	)

	// 從redis取得用戶遊戲紀錄資料(hash類型)
	if isRedis {
		// 判斷redis裡是否有獎品資訊，有則不執行查詢資料表功能
		data, err := a.RedisConn.HashGetAllCache(config.PRIZE_REDIS + prizeID)
		if err != nil {
			return prize, errors.New("錯誤: 取得獎品快取資料發生問題")
		}

		price, _ := strconv.Atoi(data["prize_price"])
		prize.PrizeID, _ = data["prize_id"]
		prize.Game, _ = data["game"]
		prize.PrizeName, _ = data["prize_name"]
		prize.PrizeType, _ = data["prize_type"]
		prize.PrizePicture, _ = data["prize_picture"]
		prize.PrizePrice = int64(price)
		prize.PrizeMethod, _ = data["prize_method"]
		prize.PrizePassword, _ = data["prize_password"]
	}

	if prize.PrizeID == "" {
		item, err := a.Table(a.Base.TableName).
			Select(
				// 獎品
				"activity_prize.id",
				"activity_prize.activity_id",
				"activity_prize.game_id",
				"activity_prize.prize_id",
				"activity_prize.game",
				"activity_prize.prize_name",
				"activity_prize.prize_type",
				"activity_prize.prize_picture",
				"activity_prize.prize_amount",
				"activity_prize.prize_remain",
				"activity_prize.prize_price",
				"activity_prize.prize_method",
				"activity_prize.team_type",
				"activity_prize.prize_password",

				// 活動
				"activity.user_id", "activity.activity_name").
			LeftJoin(command.Join{
				FieldA:    "activity_prize.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			Where("prize_id", "=", prizeID).First()
		if err != nil || item == nil {
			return PrizeModel{}, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
		}
		prize = a.MapToModel(item)

		// 將用戶資訊加入redis
		if isRedis {
			values := []interface{}{config.PRIZE_REDIS + prizeID}
			values = append(values, "prize_id", prize.PrizeID)
			values = append(values, "game", prize.Game)
			values = append(values, "prize_name", prize.PrizeName)
			values = append(values, "prize_type", prize.PrizeType)
			values = append(values, "prize_picture", prize.PrizePicture)
			values = append(values, "prize_price", prize.PrizePrice)
			values = append(values, "prize_method", prize.PrizeMethod)
			values = append(values, "prize_password", prize.PrizePassword)

			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
				return prize, errors.New("錯誤: 設置獎品快取資料發生問題")
			}
			// 設置過期時間
			// a.RedisConn.SetExpire(config.PRIZE_REDIS+prizeID,
			// 	config.REDIS_EXPIRE)
		}
	}

	return prize, nil
}

// FindPrizes 查詢獎品資訊
func (a PrizeModel) FindPrizes(isRedis bool, gameID string) (prizes []PrizeModel, err error) {
	var datas = make(map[string]string, 0)
	if isRedis {
		// 判斷redis裡是否有獎品資訊，有則不執行查詢資料表功能
		datas, err = a.RedisConn.HashGetAllCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
		if err != nil {
			return prizes, errors.New("錯誤: 取得獎品數量快取資料發生問題")
		}

		for id, data := range datas {
			value, err := strconv.Atoi(data)
			if err != nil {
				return prizes, errors.New("錯誤: 轉換快取資料發生問題")
			}

			prizes = append(prizes, PrizeModel{
				PrizeID:     id,
				PrizeRemain: int64(value),
			})
		}
	}

	// redis無獎品資訊，從資料表查詢獎品資訊
	if len(datas) == 0 {
		// 從資料表查詢獎品資訊
		items, err := a.Table(a.Base.TableName).
			Select(
				// 獎品
				"activity_prize.id",
				"activity_prize.activity_id",
				"activity_prize.game_id",
				"activity_prize.prize_id",
				"activity_prize.game",
				"activity_prize.prize_name",
				"activity_prize.prize_type",
				"activity_prize.prize_picture",
				"activity_prize.prize_amount",
				"activity_prize.prize_remain",
				"activity_prize.prize_price",
				"activity_prize.prize_method",
				"activity_prize.team_type",
				"activity_prize.prize_password",

				// 活動
				"activity.user_id").
			LeftJoin(command.Join{
				FieldA:    "activity_prize.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).

			// Where("game", "=", game).
			Where("game_id", "=", gameID).All()
		if err != nil {
			return prizes, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
		}
		prizes = MapToPrizeModel(items)

		if isRedis && len(prizes) > 0 {
			values := []interface{}{config.GAME_PRIZES_AMOUNT_REDIS + gameID} // redis參數資訊

			// 獎品資訊加入redis參數中
			for i := 0; i < len(prizes); i++ {
				values = append(values, prizes[i].PrizeID, prizes[i].PrizeRemain)
			}

			// 將獎品資訊設置至redis中
			err = a.RedisConn.HashMultiSetCache(values)
			if err != nil {
				return prizes, errors.New("錯誤: 設置獎品資訊快取發生問題")
			}
			// 設置過期時間
			// a.RedisConn.SetExpire(config.GAME_PRIZES_AMOUNT_REDIS+gameID,
			// 	config.REDIS_EXPIRE)
		}
	}
	return
}

// FindExistPrizes 尋找剩餘數量大於0的獎品
func (a PrizeModel) FindExistPrizes(game, gameID string) ([]PrizeModel, error) {
	items, err := a.Table(a.Base.TableName).
		Select(
			// 獎品
			"activity_prize.id",
			"activity_prize.activity_id",
			"activity_prize.game_id",
			"activity_prize.prize_id",
			"activity_prize.game",
			"activity_prize.prize_name",
			"activity_prize.prize_type",
			"activity_prize.prize_picture",
			"activity_prize.prize_amount",
			"activity_prize.prize_remain",
			"activity_prize.prize_price",
			"activity_prize.prize_method",
			"activity_prize.team_type",
			"activity_prize.prize_password",

			// 活動
			"activity.user_id").
		LeftJoin(command.Join{
			FieldA:    "activity_prize.activity_id",
			FieldA1:   "activity.activity_id",
			Table:     "activity",
			Operation: "="}).
		Where("game", "=", game).
		Where("game_id", "=", gameID).
		Where("prize_remain", ">", 0).All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
	}
	return MapToPrizeModel(items), nil
}

// FindPrizesAmount 利用redis取得獎品數，redis沒有資料則從資料表查詢並將獎品資訊加入redis
func (a PrizeModel) FindPrizesAmount(isRedis bool, gameID string) (count int64, err error) {
	var datas = make(map[string]string, 0)
	// 判斷redis裡是否有獎品資訊，有則不執行查詢資料表功能
	if isRedis {
		datas, err = a.RedisConn.HashGetAllCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)
		if err != nil {
			return 0, errors.New("錯誤: 取得獎品數量快取資料發生問題")
		}

		for _, data := range datas {
			value, err := strconv.Atoi(data)
			if err != nil {
				return 0, errors.New("錯誤: 轉換快取資料發生問題")
			}
			count += int64(value)
		}
		// fmt.Println("datas: ", datas)
	}

	// redis無獎品資訊，從資料表查詢獎品資訊
	if len(datas) == 0 {
		values := []interface{}{config.GAME_PRIZES_AMOUNT_REDIS + gameID} // redis參數資訊

		// 從資料表查詢獎品資訊
		prizes, err := a.FindPrizes(false, gameID)
		if err != nil {
			return 0, err
		}

		// 獎品總數
		for i := 0; i < len(prizes); i++ {
			count += prizes[i].PrizeRemain
			values = append(values, prizes[i].PrizeID, prizes[i].PrizeRemain)
		}

		if isRedis && len(prizes) > 0 {
			// 將獎品資訊設置至redis中
			err = a.RedisConn.HashMultiSetCache(values)
			if err != nil {
				return 0, errors.New("錯誤: 設置獎品資訊快取發生問題")
			}
			// 設置過期時間
			// a.RedisConn.SetExpire(config.GAME_PRIZES_AMOUNT_REDIS+gameID,
			// 	config.REDIS_EXPIRE)
		}
	}
	return int64(count), nil
}

// Add 增加獎品資料
func (a PrizeModel) Add(isRedis bool, game, prizeID string, model EditPrizeModel) error {
	var (
		fields = []string{
			"activity_id",
			"game_id",
			"prize_id",
			"game",
			"prize_name",
			"prize_type",
			"prize_picture",
			"prize_amount",
			"prize_remain",
			"prize_price",
			"prize_method",
			"prize_password", // 加密
			"team_type",
		}
	)

	if model.PrizeName == "" || utf8.RuneCountInString(model.PrizeName) > 20 {
		return errors.New("錯誤: 獎品名稱上限為20個字元，請輸入有效的獎品名稱")
	}

	if (game == "lottery" || game == "tugofwar") &&
		(model.PrizeType != "first" && model.PrizeType != "second" &&
			model.PrizeType != "third" && model.PrizeType != "general" &&
			model.PrizeType != "thanks") {
		return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
	} else if (game == "redpack" || game == "ropepack" || game == "whack_mole" ||
		game == "draw_numbers" || game == "monopoly" || game == "QA" ||
		game == "3DGachaMachine" ||
		game == "vote") &&
		(model.PrizeType != "first" && model.PrizeType != "second" &&
			model.PrizeType != "third" && model.PrizeType != "general") {
		// 紅包遊戲、敲敲樂、搖號抽獎沒有謝謝參與獎
		return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
	}

	if _, err := strconv.Atoi(model.PrizeAmount); err != nil {
		return errors.New("錯誤: 獎品數量資料發生問題，請輸入有效的獎品數量")
	}

	if _, err := strconv.Atoi(model.PrizePrice); err != nil {
		return errors.New("錯誤: 獎值價值資料發生問題，請輸入有效的獎品價值")
	}

	if (game == "lottery") &&
		(model.PrizeMethod != "site" && model.PrizeMethod != "mail" &&
			model.PrizeMethod != "thanks") {
		return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
	} else if (game == "redpack" || game == "ropepack" || game == "whack_mole" ||
		game == "draw_numbers" || game == "monopoly" || game == "QA" ||
		game == "tugofwar" || game == "3DGachaMachine" ||
		game == "vote") &&
		(model.PrizeMethod != "site" && model.PrizeMethod != "mail") {
		// 紅包遊戲、敲敲樂、搖號抽獎沒有謝謝參與獎
		return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
	}
	if model.PrizePassword == "" || utf8.RuneCountInString(model.PrizePassword) > 8 {
		return errors.New("錯誤: 獎品密碼不能為空並且上限為8個字元，請輸入有效的獎品密碼")
	}

	// 隊伍類型
	if (game == "tugofwar") &&
		(model.TeamType != "win" && model.TeamType != "lose") {
		return errors.New("錯誤: 隊伍類型資料發生問題，請輸入有效的隊伍類型")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	data["prize_id"] = prizeID
	data["game"] = game
	data["prize_password"] = string(utils.Encode([]byte(model.PrizePassword)))
	data["prize_remain"] = model.PrizeAmount

	if _, err := a.Table(a.TableName).
		Insert(FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 新增遊戲獎品發生問題")
	}

	// 修改遊戲場次的編輯次數(mysql，刷新遊戲頁面)
	// if err := a.Table(config.ACTIVITY_GAME_TABLE).
	// 	Where("game_id", "=", model.GameID).
	// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
	// }

	// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
	filter := bson.M{"game_id": model.GameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		// "$set": fieldValues,
		// "$unset": bson.M{},                // 移除不需要的欄位
		"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return err
	}


	if isRedis {
		// 要刪除的 Redis Key
		delKeys := []string{
			config.GAME_REDIS + model.GameID,                            // 遊戲設置
			config.SCORES_REDIS + model.GameID,                          // 分數
			config.SCORES_2_REDIS + model.GameID,                        // 第二分數
			config.CORRECT_REDIS + model.GameID,                         // 答對題數
			config.QA_REDIS + model.GameID,                              // 快問快答題目資訊
			config.QA_RECORD_REDIS + model.GameID,                       // 快問快答答題紀錄
			config.WINNING_STAFFS_REDIS + model.GameID,                  // 中獎人員
			config.NO_WINNING_STAFFS_REDIS + model.GameID,               // 未中獎人員
			config.GAME_TEAM_REDIS + model.GameID,                       // 遊戲隊伍資訊
			config.GAME_BINGO_NUMBER_REDIS + model.GameID,               // 紀錄抽過的號碼
			config.GAME_BINGO_USER_REDIS + model.GameID,                 // 賓果中獎人員
			config.GAME_BINGO_USER_NUMBER + model.GameID,                // 紀錄玩家號碼排序
			config.GAME_BINGO_USER_GOING_BINGO + model.GameID,           // 玩家即將賓果
			config.GAME_ATTEND_REDIS + model.GameID,                     // 遊戲人數資訊
			config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + model.GameID,  // 拔河左隊人數
			config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + model.GameID, // 拔河右隊人數
			config.GAME_PRIZES_AMOUNT_REDIS + model.GameID,              // 遊戲獎品數量
		}

		for _, key := range delKeys {
			a.RedisConn.DelCache(key)
		}

		// 要發布的 Redis Channel
		publishChannels := []string{
			config.CHANNEL_GAME_REDIS + model.GameID,              // 遊戲設置更新
			config.CHANNEL_GUEST_GAME_STATUS_REDIS + model.GameID, // 遊戲狀態更新
			config.CHANNEL_GAME_BINGO_NUMBER_REDIS + model.GameID, // 賓果號碼更新
			config.CHANNEL_QA_REDIS + model.GameID,                // 快問快答題目更新
			config.CHANNEL_GAME_TEAM_REDIS + model.GameID,         // 隊伍資訊更新
			config.CHANNEL_GAME_EDIT_TIMES_REDIS + model.GameID,   // 編輯次數更新
			config.CHANNEL_WINNING_STAFFS_REDIS + model.GameID,    // 中獎人員更新
			config.CHANNEL_GAME_BINGO_USER_NUMBER + model.GameID,  // 玩家號碼更新
			config.CHANNEL_SCORES_REDIS + model.GameID,            // 分數更新
		}

		for _, channel := range publishChannels {
			a.RedisConn.Publish(channel, "修改資料")
		}

	}
	return nil
}

// Update 更新獎品資料
func (a PrizeModel) Update(isRedis bool, game string, model EditPrizeModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"prize_name",
			"prize_type",
			"prize_picture",
			"prize_amount",
			"prize_price",
			"prize_method"}
	)
	if utf8.RuneCountInString(model.PrizeName) > 20 {
		return errors.New("錯誤: 獎品名稱上限為20個字元，請輸入有效的獎品名稱")
	}

	if model.PrizeType != "" {
		if (game == "lottery" || game == "tugofwar") &&
			(model.PrizeType != "first" && model.PrizeType != "second" &&
				model.PrizeType != "third" && model.PrizeType != "general" &&
				model.PrizeType != "thanks") {
			return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
		} else if (game == "redpack" || game == "ropepack" || game == "whack_mole" ||
			game == "draw_numbers" || game == "monopoly" || game == "QA" ||
			game == "bingo" || game == "3DGachaMachine" ||
			game == "vote") &&
			(model.PrizeType != "first" && model.PrizeType != "second" &&
				model.PrizeType != "third" && model.PrizeType != "general") {
			// 紅包遊戲、敲敲樂、搖號抽獎、賓果遊戲沒有謝謝參與獎
			return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
		}
	}

	if model.PrizeAmount != "" {
		_, err := strconv.Atoi(model.PrizeAmount)
		if err != nil {
			return errors.New("錯誤: 獎品數量資料發生問題，請輸入有效的獎品數量")
		}
	}

	if model.PrizePrice != "" {
		if _, err := strconv.Atoi(model.PrizePrice); err != nil {
			return errors.New("錯誤: 獎值價值資料發生問題，請輸入有效的獎品價值")
		}
	}

	if model.PrizeMethod != "" {
		if (game == "lottery") &&
			(model.PrizeMethod != "site" && model.PrizeMethod != "mail" &&
				model.PrizeMethod != "thanks") {
			return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
		} else if (game == "redpack" || game == "ropepack" || game == "whack_mole" ||
			game == "draw_numbers" || game == "monopoly" || game == "QA" ||
			game == "bingo" || game == "tugofwar" || game == "3DGachaMachine" ||
			game == "vote") &&
			(model.PrizeMethod != "site" && model.PrizeMethod != "mail") {
			// 紅包遊戲、敲敲樂、搖號抽獎沒有謝謝參與獎
			return errors.New("錯誤: 獎品類型資料發生問題，請輸入有效的獎品類型")
		}
	}
	if model.PrizePassword != "" {
		if utf8.RuneCountInString(model.PrizePassword) > 8 {
			return errors.New("錯誤: 獎品密碼上限為8個字元，請輸入有效的獎品密碼")
		}
		fieldValues["prize_password"] = string(utils.Encode([]byte(model.PrizePassword)))
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

	if err := a.Table(a.Base.TableName).
		Where("prize_id", "=", model.PrizeID).
		Update(fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New(`錯誤: 獎品數量資料發生問題(剩餘數量會隨著獎品數量遞增或遞減，遞減後的剩餘數量不能小於0)，請輸入有效的獎品數量`)
	}

	// 修改遊戲場次的編輯次數(刷新遊戲頁面)
	// if err := a.Table(config.ACTIVITY_GAME_TABLE).
	// 	Where("game_id", "=", model.GameID).
	// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
	// }

	// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
	filter := bson.M{"game_id": model.GameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		// "$set": fieldValues,
		// "$unset": bson.M{},                // 移除不需要的欄位
		"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return err
	}

	if isRedis {
		// 清除遊戲redis資訊(並重新開啟遊戲頁面)

		// 要刪除的 Redis Key
		delKeys := []string{
			config.GAME_REDIS + model.GameID,                            // 遊戲設置
			config.SCORES_REDIS + model.GameID,                          // 分數
			config.SCORES_2_REDIS + model.GameID,                        // 第二分數
			config.CORRECT_REDIS + model.GameID,                         // 答對題數
			config.QA_REDIS + model.GameID,                              // 快問快答題目資訊
			config.QA_RECORD_REDIS + model.GameID,                       // 快問快答答題紀錄
			config.WINNING_STAFFS_REDIS + model.GameID,                  // 中獎人員
			config.NO_WINNING_STAFFS_REDIS + model.GameID,               // 未中獎人員
			config.GAME_TEAM_REDIS + model.GameID,                       // 遊戲隊伍資訊
			config.GAME_BINGO_NUMBER_REDIS + model.GameID,               // 紀錄抽過的號碼
			config.GAME_BINGO_USER_REDIS + model.GameID,                 // 賓果中獎人員
			config.GAME_BINGO_USER_NUMBER + model.GameID,                // 紀錄玩家號碼排序
			config.GAME_BINGO_USER_GOING_BINGO + model.GameID,           // 玩家即將賓果
			config.GAME_ATTEND_REDIS + model.GameID,                     // 遊戲人數資訊
			config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + model.GameID,  // 拔河左隊人數
			config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + model.GameID, // 拔河右隊人數
			config.GAME_PRIZES_AMOUNT_REDIS + model.GameID,              // 遊戲獎品數量
			config.PRIZE_REDIS + model.PrizeID,                          // 獎品資訊
		}

		for _, key := range delKeys {
			a.RedisConn.DelCache(key)
		}

		// 要發布的 Redis Channel
		publishChannels := []string{
			config.CHANNEL_GAME_REDIS + model.GameID,              // 遊戲設置更新
			config.CHANNEL_GUEST_GAME_STATUS_REDIS + model.GameID, // 遊戲狀態更新
			config.CHANNEL_GAME_BINGO_NUMBER_REDIS + model.GameID, // 賓果號碼更新
			config.CHANNEL_QA_REDIS + model.GameID,                // 快問快答題目更新
			config.CHANNEL_GAME_TEAM_REDIS + model.GameID,         // 隊伍資訊更新
			config.CHANNEL_GAME_EDIT_TIMES_REDIS + model.GameID,   // 編輯次數更新
			config.CHANNEL_WINNING_STAFFS_REDIS + model.GameID,    // 中獎人員更新
			config.CHANNEL_GAME_BINGO_USER_NUMBER + model.GameID,  // 玩家號碼更新
			config.CHANNEL_SCORES_REDIS + model.GameID,            // 分數更新
		}

		for _, channel := range publishChannels {
			a.RedisConn.Publish(channel, "修改資料")
		}

	}
	return nil
}

// UpdateRemain 更新剩餘獎品(資料表、redis)
func (a PrizeModel) UpdateRemain(gameID, prizeID string, remain int64) error {
	if err := a.Table(a.Base.TableName).
		Where("prize_id", "=", prizeID).
		// Where("prize_remain", ">", 0).
		Update(command.Value{"prize_remain": remain}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// log.Println(prizeID, remain)
		return errors.New("錯誤: 更新獎品剩餘數量發生問題")
	}

	// 更新redis獎品數量
	a.RedisConn.HashSetCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID, remain)
	return nil
}

// DecrRemain 遞減剩餘獎品(資料表、redis)
func (a PrizeModel) DecrRemain(gameID, prizeID string) error {
	if err := a.Table(a.Base.TableName).
		Where("prize_id", "=", prizeID).
		Where("prize_remain", ">", 0).
		Update(command.Value{"prize_remain": "prize_remain - 1"}); err != nil {
		if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 抱歉，獎品無庫存")
		}
		return err
	}

	// 遞減redis獎品數量
	a.RedisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)
	return nil
}

// MapToModel map轉換model
func (a PrizeModel) MapToModel(item map[string]interface{}) PrizeModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &a)

	// a.PrizePicture, _ = item["prize_picture"].(string)
	if !strings.Contains(a.PrizePicture, "system") && a.Game != "vote" {
		a.PrizePicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/game/" + a.Game + "/" + a.GameID + "/" + a.PrizePicture
	} else if !strings.Contains(a.PrizePicture, "system") && a.Game == "vote" {
		a.PrizePicture = "/admin/uploads/" + a.UserID + "/" + a.ActivityID + "/interact/sign/" + a.Game + "/" + a.GameID + "/" + a.PrizePicture
	}

	// 獎品密碼解碼
	// password, _ := item["prize_password"].(string)
	passwordByte, _ := utils.Decode([]byte(utils.GetString(a.PrizePassword, "")))
	// password = string(passwordByte)
	a.PrizePassword = string(passwordByte)

	return a
}

// MapToPrizeModel map將值設置至[]PrizeModel
func MapToPrizeModel(items []map[string]interface{}) []PrizeModel {
	var prizes = make([]PrizeModel, 0)
	for _, item := range items {
		var (
			prize PrizeModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &prize)

		// prize.PrizePicture, _ = item["prize_picture"].(string)
		if !strings.Contains(prize.PrizePicture, "system") && prize.Game != "vote" {
			prize.PrizePicture = "/admin/uploads/" + prize.UserID + "/" + prize.ActivityID + "/interact/game/" + prize.Game + "/" + prize.GameID + "/" + prize.PrizePicture
		} else if !strings.Contains(prize.PrizePicture, "system") && prize.Game == "vote" {
			prize.PrizePicture = "/admin/uploads/" + prize.UserID + "/" + prize.ActivityID + "/interact/sign/" + prize.Game + "/" + prize.GameID + "/" + prize.PrizePicture
		}

		// 獎品密碼解碼
		// password, _ := item["prize_password"].(string)
		passwordByte, _ := utils.Decode([]byte(utils.GetString(prize.PrizePassword, "")))
		// password = string(passwordByte)
		prize.PrizePassword = string(passwordByte)

		prizes = append(prizes, prize)
	}
	return prizes
}

// command.Value{
// 	"activity_id":    model.ActivityID,
// 	"game_id":        model.GameID,
// 	"prize_id":       prizeID,
// 	"game":           game,
// 	"prize_name":     model.PrizeName,
// 	"prize_type":     model.PrizeType,
// 	"prize_picture":  model.PrizePicture,
// 	"prize_amount":   model.PrizeAmount,
// 	"prize_remain":   model.PrizeAmount,
// 	"prize_price":    model.PrizePrice,
// 	"prize_method":   model.PrizeMethod,
// 	"prize_password": string(utils.Encode([]byte(model.PrizePassword))), // 加密
// 	"team_type":      model.TeamType,
// }

// prize.ID, _ = item["id"].(int64)
// prize.UserID, _ = item["user_id"].(string)
// prize.ActivityID, _ = item["activity_id"].(string)
// prize.GameID, _ = item["game_id"].(string)
// prize.PrizeID, _ = item["prize_id"].(string)
// prize.Game, _ = item["game"].(string)
// prize.PrizeName, _ = item["prize_name"].(string)
// prize.PrizeType, _ = item["prize_type"].(string)
// *****舊*****
// prize.PrizePicture, _ = item["prize_picture"].(string)
// *****舊*****
// *****新*****
// *****新*****

// prize.PrizeAmount, _ = item["prize_amount"].(int64)
// prize.PrizeRemain, _ = item["prize_remain"].(int64)
// prize.PrizePrice, _ = item["prize_price"].(int64)
// prize.PrizeMethod, _ = item["prize_method"].(string)
// prize.PrizePassword, _ = item["prize_password"].(string)
// prize.TeamType, _ = item["team_type"].(string)

// a.ID, _ = item["id"].(int64)
// a.UserID, _ = item["user_id"].(string)
// a.ActivityID, _ = item["activity_id"].(string)
// a.GameID, _ = item["game_id"].(string)
// a.PrizeID, _ = item["prize_id"].(string)
// a.Game, _ = item["game"].(string)
// a.PrizeName, _ = item["prize_name"].(string)
// a.PrizeType, _ = item["prize_type"].(string)
// a.PrizeAmount, _ = item["prize_amount"].(int64)
// a.PrizeRemain, _ = item["prize_remain"].(int64)
// a.PrizePrice, _ = item["prize_price"].(int64)
// a.PrizeMethod, _ = item["prize_method"].(string)
// a.PrizePassword, _ = item["prize_password"].(string)
// a.TeamType, _ = item["team_type"].(string)

// a.ActivityName, _ = item["activity_name"].(string)

// FindPrizesByGameID 查詢獎品資訊(多個)
// func (a PrizeModel) FindPrizesByGameID(game string) ([]PrizeModel, error) {
// 	items, err := a.Table(a.Base.TableName).
// 		Where("game_id", "=", game).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
// 	}
// 	return MapToPrizeModel(items), nil
// }

// FindExistPrizesByGameID 尋找剩餘數量大於0的獎品
// func (a PrizeModel) FindExistPrizesByGameID(game string) ([]PrizeModel, error) {
// 	items, err := a.Table(a.Base.TableName).
// 		Where("game_id", "=", game).
// 		Where("prize_remain", ">", 0).All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
// 	}
// 	return MapToPrizeModel(items), nil
// }

// Find 尋找資料
// func (a PrizeModel) Find(prize interface{}) (PrizeModel, error) {
// 	item, err := a.Table(a.Base.TableName).Where("prize_id", "=", prize).First()
// 	if err != nil {
// 		return PrizeModel{}, errors.New("錯誤: 無法取得獎品資訊，請重新查詢")
// 	}
// 	if item == nil {
// 		return PrizeModel{}, nil
// 	}
// 	return a.MapToModel(item), nil
// }

// ex: prize_remain = prire_remain + model.PrizeAmount - prize_amount, `prize_amount` = model.PrizeAmount
// fieldValues["prize_remain"] = model.PrizeAmount

// 查詢原始獎品數量、剩餘數量
// prizeModel, err := a.FindPrize(model.PrizeID)
// if err != nil {
// 	return err
// }

// remain := int(prizeModel.PrizeRemain) + amountInt - int(prizeModel.PrizeAmount)
// if remain < 0 {
// 	return errors.New("錯誤: 獎品數量資料發生問題(剩餘數量會隨著獎品數量遞增或遞減，遞減後的剩餘數量不能小於0)，請輸入有效的獎品數量")
// }
// fieldValues["prize_amount"] = amountInt

// amountInt, err := strconv.Atoi(model.Amount)
// if err != nil {
// 	return errors.New("錯誤: 獎品數量資料發生問題，請輸入有效的獎品數量")
// }
// remainInt := amountInt - int(prizeModel.Amount) // 變更數量

// 如果數量未更改，amount、remain不變
// if amountInt == int(prizeModel.Amount) {
// 	fmt.Println("數量不變")
// 	fieldValues["amount"] = prizeModel.Amount
// 	fieldValues["remain"] = prizeModel.Remain
// } else {
// 遞增or遞減，剩餘數量不能小於0
// fmt.Println("遞增or遞減")
// remain := int(prizeModel.Remain) + amountInt - int(prizeModel.Amount)
// if remain < 0 {
// 	return errors.New("錯誤: 獎品數量資料發生問題(剩餘數量會隨著獎品數量遞增或遞減，遞減後的剩餘數量不能小於0)，請輸入有效的獎品數量")
// }

// ex: prize_remain = prire_remain + model.PrizeAmount - prize_amount
// fieldValues["prize_remain"] = fmt.Sprintf("prize_remain + %s - prize_amount", model.PrizeAm

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_QA_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+model.GameID, "修改資料")

// a.RedisConn.DelCache(config.GAME_REDIS + model.GameID)                            // 遊戲設置
// a.RedisConn.DelCache(config.SCORES_REDIS + model.GameID)                          // 分數
// a.RedisConn.DelCache(config.SCORES_2_REDIS + model.GameID)                        // 第二分數
// a.RedisConn.DelCache(config.CORRECT_REDIS + model.GameID)                         // 答對題數
// a.RedisConn.DelCache(config.QA_REDIS + model.GameID)                              // 快問快答題目資訊
// a.RedisConn.DelCache(config.QA_RECORD_REDIS + model.GameID)                       // 快問快答答題紀錄
// a.RedisConn.DelCache(config.WINNING_STAFFS_REDIS + model.GameID)                  // 中獎人員
// a.RedisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + model.GameID)               // 未中獎人員
// a.RedisConn.DelCache(config.GAME_TEAM_REDIS + model.GameID)                       // 遊戲隊伍資訊，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + model.GameID)               // 紀錄抽過的號碼，LIST
// a.RedisConn.DelCache(config.GAME_BINGO_USER_REDIS + model.GameID)                 // 賓果中獎人員，ZSET
// a.RedisConn.DelCache(config.GAME_BINGO_USER_NUMBER + model.GameID)                // 紀錄玩家的號碼排序，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + model.GameID)           // 紀錄玩家是否即將賓果，HASH
// a.RedisConn.DelCache(config.GAME_ATTEND_REDIS + model.GameID)                     // 遊戲人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + model.GameID)  // 拔河遊戲左隊人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + model.GameID) // 拔河遊戲右隊人數資訊，SET

// 清除獎品redis資訊
// a.RedisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + model.GameID)
// a.RedisConn.DelCache(config.PRIZE_REDIS + model.PrizeID)

// 清除遊戲redis資訊(並重新開啟遊戲頁面)
// a.RedisConn.DelCache(config.GAME_REDIS + model.GameID)                            // 遊戲設置
// a.RedisConn.DelCache(config.SCORES_REDIS + model.GameID)                          // 分數
// a.RedisConn.DelCache(config.SCORES_2_REDIS + model.GameID)                        // 第二分數
// a.RedisConn.DelCache(config.CORRECT_REDIS + model.GameID)                         // 答對題數
// a.RedisConn.DelCache(config.QA_REDIS + model.GameID)                              // 快問快答題目資訊
// a.RedisConn.DelCache(config.QA_RECORD_REDIS + model.GameID)                       // 快問快答答題紀錄
// a.RedisConn.DelCache(config.WINNING_STAFFS_REDIS + model.GameID)                  // 中獎人員
// a.RedisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + model.GameID)               // 未中獎人員
// a.RedisConn.DelCache(config.GAME_TEAM_REDIS + model.GameID)                       // 遊戲隊伍資訊，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + model.GameID)               // 紀錄抽過的號碼，LIST
// a.RedisConn.DelCache(config.GAME_BINGO_USER_REDIS + model.GameID)                 // 賓果中獎人員，ZSET
// a.RedisConn.DelCache(config.GAME_BINGO_USER_NUMBER + model.GameID)                // 紀錄玩家的號碼排序，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + model.GameID)           // 紀錄玩家是否即將賓果，HASH
// a.RedisConn.DelCache(config.GAME_ATTEND_REDIS + model.GameID)                     // 遊戲人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + model.GameID)  // 拔河遊戲左隊人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + model.GameID) // 拔河遊戲右隊人數資訊，SET

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_QA_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+model.GameID, "修改資料")

// // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+model.GameID, "修改資料")

// // 清除獎品redis資訊
// a.RedisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + model.GameID)

// values = []string{model.PrizeName, model.PrizeType, model.PrizePicture,
// 	model.PrizeAmount, model.PrizePrice, model.PrizeMethod}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
