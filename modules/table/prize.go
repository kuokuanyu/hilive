package table

import (
	"encoding/json"
	"errors"

	"hilive/models"

	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// GetLotteryPrizePanel 遊戲抽獎
func (s *SystemTable) GetLotteryPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "lottery", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "lottery", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetRedpackPrizePanel 搖紅包獎品
func (s *SystemTable) GetRedpackPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "redpack", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "redpack", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetRopepackPrizePanel 套紅包獎品
func (s *SystemTable) GetRopepackPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "ropepack", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "ropepack", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetWhackMolePrizePanel 敲敲樂獎品
func (s *SystemTable) GetWhackMolePrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "whack_mole", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "whack_mole", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetDrawNumbersPrizePanel 搖號抽獎獎品
func (s *SystemTable) GetDrawNumbersPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {

		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)                            // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)              // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)           // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "draw_numbers", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("prize_id") {
			return errors.New("錯誤: 獎品ID發生問題，請輸入有效的獎品ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "draw_numbers", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetMonopolyPrizePanel 鈔級大富翁獎品
func (s *SystemTable) GetMonopolyPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "monopoly", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "monopoly", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetQAPrizePanel 快問快答獎品
func (s *SystemTable) GetQAPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

		// 刪除redis資訊
		for _, id := range idArr {
			s.redisConn.DelCache(config.PRIZE_REDIS + id)
		}

		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
		// if err := s.table(config.ACTIVITY_GAME_TABLE).
		// 	Where("game_id", "=", gameID).
		// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
		// }

		// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除 Redis Cache
		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

		// 發布 Redis 頻道消息
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

		return nil
	})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id") {
			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "QA", utils.UUID(20), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "QA", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetTugofwarPrizePanel 拔河遊戲獎品
func (s *SystemTable) GetTugofwarPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE)

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE)

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "tugofwar", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetBingoPrizePanel 賓果遊戲獎品
func (s *SystemTable) GetBingoPrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE)

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE)

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "bingo", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// Get3DGachaMachinePrizePanel 扭蛋機遊戲獎品
func (s *SystemTable) Get3DGachaMachinePrizePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_PRIZE_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)
			s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
			s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

			// 刪除redis資訊
			for _, id := range idArr {
				s.redisConn.DelCache(config.PRIZE_REDIS + id)
			}

			// 修改遊戲場次的編輯次數(刷新遊戲頁面)
			// if err := s.table(config.ACTIVITY_GAME_TABLE).
			// 	Where("game_id", "=", gameID).
			// 	Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
			// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			// 	return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
			// }

			// 修改遊戲場次的編輯次數(mongo，刷新遊戲頁面)
			filter := bson.M{"game_id": gameID} // 過濾條件
			// 要更新的值
			update := bson.M{
				// "$set": fieldValues,
				// "$unset": bson.M{},                // 移除不需要的欄位
				"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
			}

			if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
				return err
			}

			// 清除 Redis Cache
			s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
			s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
			// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
			s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

			// 發布 Redis 頻道消息
			s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

			return nil
		})

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_PRIZE_TABLE).
		SetInsertFunc(func(values form2.Values) error {
			if values.IsEmpty("activity_id", "game_id") {
				return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
			}

			var picture string
			if values.Get("prize_picture") != "" {
				picture = values.Get("prize_picture")
			} else { // 預設
				picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditPrizeModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			// 手動處理
			model.PrizePicture = picture

			if err := models.DefaultPrizeModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Add(true, "3DGachaMachine", utils.UUID(20), model); err != nil {
				return err
			}
			return nil
		})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id", "prize_id") {
			return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
		}

		var picture string
		if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
		} else if values.Get("prize_picture") != "" {
			picture = values.Get("prize_picture")
		} else {
			picture = ""
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.PrizePicture = picture

		if err := models.DefaultPrizeModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "3DGachaMachine", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// GetVotePrizePanel 投票遊戲獎品
// func (s *SystemTable) GetVotePrizePanel() (table Table) {
// table = DefaultTable(DefaultConfig())
// info := table.GetInfo()
// info.SetTable(config.ACTIVITY_PRIZE_TABLE).
// 	SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
// 		var ids = interfaces(idArr)
// 		s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
// 		s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

// 		// 刪除redis資訊
// 		for _, id := range idArr {
// 			s.redisConn.DelCache(config.PRIZE_REDIS + id)
// 		}

// 		// 修改遊戲場次的編輯次數(刷新遊戲頁面)
// 		if err := s.table(config.ACTIVITY_GAME_TABLE).
// 			Where("game_id", "=", gameID).
// 			Update(command.Value{"edit_times": "edit_times + 1"}); err != nil &&
// 			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 			return errors.New("錯誤: 更新遊戲場次編輯次數資料發生問題，請重新執行")
// 		}

// 		// 清除 Redis Cache
// 		s.redisConn.DelCache(config.GAME_REDIS + gameID)               // 遊戲設置
// 		s.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID) // 刪除獎品時，需要刪除遊戲獎品數量redis資料
// 		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID) // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
// 		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID) // 刪除玩家遊戲紀錄(中獎.未中獎)

// 		// 發布 Redis 頻道消息
// 		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")              // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料") // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")   // 當接收到 WebSocket 消息時，將其發布到 Redis 頻道

// 		return nil
// 	})

// formList := table.GetForm()
// formList.SetTable(config.ACTIVITY_PRIZE_TABLE).
// 	SetInsertFunc(func(values form2.Values) error {
// 		if values.IsEmpty("activity_id", "game_id") {
// 			return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
// 		}

// 		var picture string
// 		if values.Get("prize_picture") != "" {
// 			picture = values.Get("prize_picture")
// 		} else { // 預設
// 			picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
// 		}

// 		// 將map[string][]string格是資料轉換為map[string]string
// 		flattened := utils.FlattenForm(values)

// 		// 將 map 轉成 JSON
// 		jsonBytes, err := json.Marshal(flattened)
// 		if err != nil {
// 			return err
// 		}

// 		// 轉成 struct
// 		var model models.EditPrizeModel
// 		if err := json.Unmarshal(jsonBytes, &model); err != nil {
// 			return err
// 		}

// 		// 手動處理
// 		model.PrizePicture = picture

// 		if err := models.DefaultPrizeModel().
// 			SetDbConn(s.dbConn).
// 			SetRedisConn(s.redisConn).
// 			Add(true, "vote", utils.UUID(20), model); err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// formList.SetUpdateFunc(func(values form2.Values) error {
// 	if values.IsEmpty("game_id", "prize_id") {
// 		return errors.New("錯誤: 遊戲、獎品ID發生問題，請輸入有效的ID")
// 	}

// 	var picture string
// 	if values.Get("prize_picture"+DEFAULT_FALG) == "1" { // 預設
// 		picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
// 	} else if values.Get("prize_picture") != "" {
// 		picture = values.Get("prize_picture")
// 	} else {
// 		picture = ""
// 	}

// 	// 將map[string][]string格是資料轉換為map[string]string
// 	flattened := utils.FlattenForm(values)

// 	// 將 map 轉成 JSON
// 	jsonBytes, err := json.Marshal(flattened)
// 	if err != nil {
// 		return err
// 	}

// 	// 轉成 struct
// 	var model models.EditPrizeModel
// 	if err := json.Unmarshal(jsonBytes, &model); err != nil {
// 		return err
// 	}

// 	// 手動處理
// 	model.PrizePicture = picture

// 	if err := models.DefaultPrizeModel().
// 		SetDbConn(s.dbConn).
// 		SetRedisConn(s.redisConn).
// 		Update(true, "vote", model); err != nil {
// 		return err
// 	}
// 	return nil
// })
// return
// }

// .SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
// 	var ids = interfaces(idArr)
// 	s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
// 	s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

// 	// 刪除redis資訊
// 	for _, id := range idArr {
// 		s.redisConn.DelCache(config.PRIZE_REDIS + id)
// 	}
// 	return nil
// })

// .SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "game_id") {
// 		return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
// 	}

// 	var picture string
// 	if values.Get("prize_picture") != "" {
// picture = values.Get("prize_picture")
// 	} else { // 預設
// 		picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
// 	}

// 	if err := models.DefaultPrizeModel().SetDbConn(s.dbConn).SetRedisConn(s.redisConn).
// 		Add(true, "tugofwar", utils.UUID(20), models.NewPrizeModel{
// 			ActivityID:    values.Get("activity_id"),
// 			GameID:        values.Get("game_id"),
// 			PrizeName:     values.Get("prize_name"),
// 			PrizeType:     values.Get("prize_type"),
// 			PrizePicture:  picture,
// 			PrizeAmount:   values.Get("prize_amount"),
// 			PrizePrice:    values.Get("prize_price"),
// 			PrizeMethod:   values.Get("prize_method"),
// 			PrizePassword: values.Get("prize_password"),
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// })

// .SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
// 	var ids = interfaces(idArr)
// 	s.table(config.ACTIVITY_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()
// 	s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("prize_id", ids).Delete()

// 	// 刪除redis資訊
// 	for _, id := range idArr {
// 		s.redisConn.DelCache(config.PRIZE_REDIS + id)
// 	}
// 	return nil
// })

// .SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "game_id") {
// 		return errors.New("錯誤: 活動ID、遊戲ID發生問題，請輸入有效的資料")
// 	}

// 	var picture string
// 	if values.Get("prize_picture") != "" {
// picture = values.Get("prize_picture")
// 	} else { // 預設
// 		picture = UPLOAD_SYSTEM_URL + "img-prize-pic.png"
// 	}

// 	if err := models.DefaultPrizeModel().SetDbConn(s.dbConn).SetRedisConn(s.redisConn).
// 		Add(true, "tugofwar", utils.UUID(20), models.NewPrizeModel{
// 			ActivityID:    values.Get("activity_id"),
// 			GameID:        values.Get("game_id"),
// 			PrizeName:     values.Get("prize_name"),
// 			PrizeType:     values.Get("prize_type"),
// 			PrizePicture:  picture,
// 			PrizeAmount:   values.Get("prize_amount"),
// 			PrizePrice:    values.Get("prize_price"),
// 			PrizeMethod:   values.Get("prize_method"),
// 			PrizePassword: values.Get("prize_password"),
// 		});
