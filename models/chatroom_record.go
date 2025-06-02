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
)

// ChatroomRecordModel 資料表欄位
type ChatroomRecordModel struct {
	Base          `json:"-"`
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	ActivityID    string `json:"activity_id"`
	MessageType   string `json:"message_type"`
	MessageStyle  string `json:"message_style"`
	MessagePrice  int64  `json:"message_price"`
	MessageEffect string `json:"message_effect"`
	MessageStatus string `json:"message_status"`
	MessagePlayed string `json:"message_played"`
	Message       string `json:"message"`
	SendTime      string `json:"send_time"`

	// 用戶資訊
	Name   string `json:"name"`
	Avatar string `json:"avatar"`

	// 過濾參數
	YesMessageAmount    int64 `json:"yes_message_amount" example:"10"`    // 通過訊息數量
	NoMessageAmount     int64 `json:"no_message_amount" example:"10"`     // 未通過訊息數量
	ReviewMessageAmount int64 `json:"review_message_amount" example:"10"` // 審核訊息數量

	UserMessageAmount int64 `json:"user_message_amount" example:"10"` // 用戶訊息數量

	YesPlayedAmount int64 `json:"yes_played_amount" example:"10"` // 已播放數量
	NoPlayedAmount  int64 `json:"no_played_amount" example:"10"`  // 未播放數量

	// 不同種類聊天訊息
	YesNormalMessageAmount  int64 `json:"yes_normal_message_amount" example:"10"`
	YesNormalBarrageAmount  int64 `json:"yes_normal_barrage_amount" example:"10"`
	YesSpecialBarrageAmount int64 `json:"yes_special_barrage_amount" example:"10"`
	YesOccupyBarrageAmount  int64 `json:"yes_occupy_barrage_amount" example:"10"`

	NoNormalMessageAmount  int64 `json:"no_normal_message_amount" example:"10"`
	NoNormalBarrageAmount  int64 `json:"no_normal_barrage_amount" example:"10"`
	NoSpecialBarrageAmount int64 `json:"no_special_barrage_amount" example:"10"`
	NoOccupyBarrageAmount  int64 `json:"no_occupy_barrage_amount" example:"10"`

	ReviewNormalMessageAmount  int64 `json:"review_normal_message_amount" example:"10"`
	ReviewNormalBarrageAmount  int64 `json:"review_normal_barrage_amount" example:"10"`
	ReviewSpecialBarrageAmount int64 `json:"review_special_barrage_amount" example:"10"`
	ReviewOccupyBarrageAmount  int64 `json:"review_occupy_barrage_amount" example:"10"`
}

// EditChatroomRecordModel 資料表欄位
type EditChatroomRecordModel struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	// Name         string `json:"name"`
	// Avatar       string `json:"avatar"`
	ActivityID    string `json:"activity_id"`
	MessageType   string `json:"message_type"`
	MessageStyle  string `json:"message_style"`
	MessagePrice  string `json:"message_price"`
	MessageEffect string `json:"message_effect"`
	MessageStatus string `json:"message_status"`
	MessagePlayed string `json:"message_played"`
	Message       string `json:"message"`
}

// DefaultChatroomRecordModel 預設ChatroomRecordModel
func DefaultChatroomRecordModel() ChatroomRecordModel {
	return ChatroomRecordModel{Base: Base{TableName: config.ACTIVITY_CHATROOM_RECORD_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (m ChatroomRecordModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ChatroomRecordModel {
	m.DbConn = dbconn
	m.RedisConn = cacheconn
	m.MongoConn = mongoconn
	return m
}

// SetDbConn 設定connection
// func (m ChatroomRecordModel) SetDbConn(conn db.Connection) ChatroomRecordModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m ChatroomRecordModel) SetRedisConn(conn cache.Connection) ChatroomRecordModel {
// 	m.RedisConn = conn
// 	return m
// }

// FindHostMessageAmount 查詢不同類型聊天訊息數量(主持端)
func (m ChatroomRecordModel) FindHostMessageAmount(activityID string) (ChatroomRecordModel, error) {
	var (
		sql = m.Table(m.Base.TableName).
			Select(
				"SUM(CASE WHEN message_status = 'yes' THEN 1 ELSE 0 END) AS yes_message_amount",
				"SUM(CASE WHEN message_status = 'no' THEN 1 ELSE 0 END) AS no_message_amount",
				"SUM(CASE WHEN message_status = 'review' THEN 1 ELSE 0 END) AS review_message_amount",

				"SUM(CASE WHEN message_played = 'yes' THEN 1 ELSE 0 END) AS yes_played_amount",
				"SUM(CASE WHEN message_played = 'no' THEN 1 ELSE 0 END) AS no_played_amount",

				"SUM(CASE WHEN message_status = 'yes' AND message_type = 'normal-message' THEN 1 ELSE 0 END) AS yes_normal_message_amount",
				"SUM(CASE WHEN message_status = 'yes' AND message_type = 'normal-barrage' THEN 1 ELSE 0 END) AS yes_normal_barrage_amount",
				"SUM(CASE WHEN message_status = 'yes' AND message_type = 'special-barrage' THEN 1 ELSE 0 END) AS yes_special_barrage_amount",
				"SUM(CASE WHEN message_status = 'yes' AND message_type = 'occupy-barrage' THEN 1 ELSE 0 END) AS yes_occupy_barrage_amount",

				"SUM(CASE WHEN message_status = 'no' AND message_type = 'normal-message' THEN 1 ELSE 0 END) AS no_normal_message_amount",
				"SUM(CASE WHEN message_status = 'no' AND message_type = 'normal-barrage' THEN 1 ELSE 0 END) AS no_normal_barrage_amount",
				"SUM(CASE WHEN message_status = 'no' AND message_type = 'special-barrage' THEN 1 ELSE 0 END) AS no_special_barrage_amount",
				"SUM(CASE WHEN message_status = 'no' AND message_type = 'occupy-barrage' THEN 1 ELSE 0 END) AS no_occupy_barrage_amount",

				"SUM(CASE WHEN message_status = 'review' AND message_type = 'normal-message' THEN 1 ELSE 0 END) AS review_normal_message_amount",
				"SUM(CASE WHEN message_status = 'review' AND message_type = 'normal-barrage' THEN 1 ELSE 0 END) AS review_normal_barrage_amount",
				"SUM(CASE WHEN message_status = 'review' AND message_type = 'special-barrage' THEN 1 ELSE 0 END) AS review_special_barrage_amount",
				"SUM(CASE WHEN message_status = 'review' AND message_type = 'occupy-barrage' THEN 1 ELSE 0 END) AS review_occupy_barrage_amount",
			)
	)

	if activityID != "" {
		sql = sql.Where("activity_chatroom_record.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return m, errors.New("錯誤: 無法取得聊天訊息數量，請重新查詢")
	}

	// 回傳的資料為[]byte格式，需轉換為int64
	for key, value := range item {
		item[key] = utils.GetInt64(value, 0)
	}

	m = m.MapToModel(item)

	return m, nil
}

// FindHostMessageAmount 查詢不同類型聊天訊息數量(玩家端)
func (m ChatroomRecordModel) FindGuestMessageAmount(activityID, userID string) (ChatroomRecordModel, error) {
	var (
		sql = m.Table(m.Base.TableName).
			Select(
				fmt.Sprintf("SUM(CASE WHEN user_id = '%s' AND message_status = 'yes' THEN 1 ELSE 0 END) AS yes_message_amount", userID),
				fmt.Sprintf("SUM(CASE WHEN user_id = '%s' AND message_status = 'no' THEN 1 ELSE 0 END) AS no_message_amount", userID),
				fmt.Sprintf("SUM(CASE WHEN user_id = '%s' AND message_status = 'review' THEN 1 ELSE 0 END) AS review_message_amount", userID),

				fmt.Sprintf("SUM(CASE WHEN user_id = '%s' THEN 1 ELSE 0 END) AS user_message_amount", userID),
			)
	)
	if activityID != "" {
		sql = sql.Where("activity_chatroom_record.activity_id", "=", activityID)
	}
	if userID != "" {
		sql = sql.Where("activity_chatroom_record.user_id", "=", userID)
	}

	item, err := sql.First()
	if err != nil {
		return m, errors.New("錯誤: 無法取得聊天訊息數量，請重新查詢")
	}

	// 回傳的資料為[]byte格式，需轉換為int64
	for key, value := range item {
		item[key] = utils.GetInt64(value, 0)
	}
	// log.Println("item: ", item)

	m = m.MapToModel(item)

	return m, nil
}

// Find 尋找聊天紀錄資料，排序由舊至新(包含通過、未通過、審核中)
func (m ChatroomRecordModel) Find(
	isRedis bool, activityID, userID,
	status string, played string, messageType string,
	limit, offset int64) ([]ChatroomRecordModel, error) {
	var (
		records = make([]ChatroomRecordModel, 0)
		// err     error
	)

	// 判斷redis裡是否有聊天紀錄資料
	if isRedis {
		// var (
		// 	idsOrderBytime = make([]string, 0) // 時間排序
		// 	datas          map[string]string
		// )

		// // 所有聊天紀錄資訊(HASH)
		// datas, err = m.RedisConn.HashGetAllCache(config.CHATROOM_REDIS + activityID)
		// if err != nil {
		// 	return records, errors.New("錯誤: 取得聊天紀錄快取資料(所有資訊)發生問題")
		// }

		// // 取得時間排序的聊天紀錄id資訊(LIST)
		// if idsOrderBytime, err = m.RedisConn.ListRange(config.CHATROOM_ORDER_BY_TIME_REDIS+activityID,
		// 	0, 0); err != nil {
		// 	return records, errors.New("錯誤: 從redis中取得時間排序的聊天紀錄資訊發生問題")
		// }

		// for i := 0; i < len(idsOrderBytime); i++ {
		// 	var (
		// 		record ChatroomRecordModel // 時間排序
		// 	)

		// 	// 解碼，取得聊天紀錄資料
		// 	json.Unmarshal([]byte(datas[idsOrderBytime[i]]), &record)

		// 	// 判斷是否查詢正確資料
		// 	if record.ID == 0 {
		// 		// 無法從redis取得聊天紀錄資料時，刪除舊的redis聊天紀錄資訊並重新執行Find function
		// 		m.RedisConn.DelCache(config.CHATROOM_REDIS + activityID)               // 聊天紀錄資訊
		// 		m.RedisConn.DelCache(config.CHATROOM_ORDER_BY_TIME_REDIS + activityID) // 聊天紀錄資訊(時間排序)

		// 		return m.Find(true, activityID, "", "", "", 0, 0)
		// 	}

		// 	records = append(records, record)
		// }
	}

	// redis中無提問快取資料，從資料表中查詢
	if len(records) == 0 {
		sql := m.Table(m.Base.TableName).
			Select("activity_chatroom_record.id", "activity_chatroom_record.activity_id",
				"activity_chatroom_record.user_id", "activity_chatroom_record.message_style",
				"activity_chatroom_record.message_type", "activity_chatroom_record.message",
				"activity_chatroom_record.send_time", "activity_chatroom_record.message_price",
				"activity_chatroom_record.message_effect", "activity_chatroom_record.message_status",
				"activity_chatroom_record.message_played",

				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_chatroom_record.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_chatroom_record.activity_id", "=", activityID).
			// Where("activity_chatroom_record.message_status", "=", "yes").
			OrderBy("activity_chatroom_record.id", "asc")

		// 判斷參數是否為空
		if userID != "" {
			sql = sql.Where("activity_chatroom_record.user_id", "=", userID)
		}

		if messageType != "" {
			sql = sql.WhereIn("activity_chatroom_record.message_type", interfaces(strings.Split(messageType, ",")))
		}

		if status != "" {
			sql = sql.WhereIn("activity_chatroom_record.message_status", interfaces(strings.Split(status, ",")))
		}
		if played != "" {
			sql = sql.WhereIn("activity_chatroom_record.message_played", interfaces(strings.Split(played, ",")))
		}
		if limit != 0 {
			sql = sql.Limit(limit)
		}
		if offset != 0 {
			sql = sql.Offset(offset)
		}

		items, err := sql.All()
		if err != nil {
			return records, errors.New("錯誤: 無法取得聊天紀錄資料(資料表)，請重新查詢")
		}
		records = MapToChatroomRecordModel(items)

		// 將聊天紀錄資訊加入redis中
		if isRedis && len(records) > 0 {
			// var (
			// 	recordsparams            = []interface{}{config.CHATROOM_REDIS + activityID}               // redis參數(聊天紀錄資訊)
			// 	recordsOrderByTimeParams = []interface{}{config.CHATROOM_ORDER_BY_TIME_REDIS + activityID} // redis參數(時間排序聊天紀錄)
			// )
			// for _, record := range records {
			// 	// json編碼
			// 	recordJson := utils.JSON(record)

			// 	// 提問資訊(HASH)
			// 	recordsparams = append(recordsparams, record.ID, recordJson)
			// 	// 時間排序(LIST)
			// 	recordsOrderByTimeParams = append(recordsOrderByTimeParams, record.ID)
			// }

			// // 將資料寫入redis中(聊天紀錄資訊，HASH)
			// if err = m.RedisConn.HashMultiSetCache(recordsparams); err != nil {
			// 	return records, errors.New("錯誤: 設置聊天紀錄快取資料(所有聊天紀錄資訊)發生問題")
			// }
			// // 將資料寫入redis中(時間排序，LIST)
			// if err = m.RedisConn.ListMultiRPush(recordsOrderByTimeParams); err != nil {
			// 	return records, errors.New("錯誤: 設置聊天紀錄快取資料(時間排序)發生問題")
			// }

			// // 設置過期時間
			// m.RedisConn.SetExpire(config.CHATROOM_REDIS+activityID, config.REDIS_EXPIRE)               // 聊天紀錄資訊
			// m.RedisConn.SetExpire(config.CHATROOM_ORDER_BY_TIME_REDIS+activityID, config.REDIS_EXPIRE) // 聊天紀錄資訊時間排序
		}
	}
	return records, nil
}

// FindAll 尋找聊天紀錄資料，排序由舊至新(包含通過、未通過、審核中)
func (m ChatroomRecordModel) FindAll(activityID string) ([]ChatroomRecordModel, error) {
	items, err := m.Table(m.Base.TableName).
		Select("activity_chatroom_record.id", "activity_chatroom_record.activity_id",
			"activity_chatroom_record.user_id", "activity_chatroom_record.message_style",
			"activity_chatroom_record.message_type", "activity_chatroom_record.message",
			"activity_chatroom_record.send_time", "activity_chatroom_record.message_price",
			"activity_chatroom_record.message_status",
			"activity_chatroom_record.message_played",

			"line_users.name", "line_users.avatar").
		LeftJoin(command.Join{
			FieldA:    "activity_chatroom_record.user_id",
			FieldA1:   "line_users.user_id",
			Table:     "line_users",
			Operation: "="}).
		Where("activity_chatroom_record.activity_id", "=", activityID).
		OrderBy("activity_chatroom_record.id", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得聊天室紀錄資料，請重新查詢")
	}
	return MapToChatroomRecordModel(items), nil
}

// Add 增加聊天紀錄資料
func (m ChatroomRecordModel) Add(isRedis bool, model EditChatroomRecordModel) error {
	var (
		fields = []string{
			"activity_id",
			"user_id",
			"message_type",
			"message_style",
			"message_price",
			"message_status",
			"message_effect",
			"message",
			"message_played",
		}
	)

	if model.MessageType != "normal-message" && model.MessageType != "normal-barrage" &&
		model.MessageType != "special-barrage" && model.MessageType != "occupy-barrage" {
		return errors.New("錯誤: 訊息種類資料發生問題，請輸入有效的訊息種類")
	}
	if model.MessageStatus != "yes" && model.MessageStatus != "no" &&
		model.MessageStatus != "review" {
		return errors.New("錯誤: 訊息審核狀態資料發生問題，請輸入有效的訊息審核狀態")
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	data["message_played"] = "no"

	_, err := m.Table(m.TableName).Insert(FilterFields(data, fields))

	// if isRedis {
	// 判斷redis裡是否有用戶資訊，有則不查詢資料表
	// user, err := DefaultLineModel().SetDbConn(m.DbConn).
	// 	SetRedisConn(m.RedisConn).Find(true, "", "user_id", model.UserID)
	// if err != nil {
	// 	return errors.New("錯誤: 無法取得用戶資訊")
	// }

	// // 將新的聊天紀錄加入redis(CHATROOM_REDIS)
	// m.RedisConn.HashSetCache(config.CHATROOM_REDIS+model.ActivityID,
	// 	strconv.Itoa(int(id)), utils.JSON(ChatroomRecordModel{
	// 		ID:            id,
	// 		ActivityID:    model.ActivityID,
	// 		UserID:        model.UserID,
	// 		MessageType:   model.MessageType,
	// 		MessageStyle:  model.MessageStyle,
	// 		MessagePrice:  int64(price),
	// 		MessageStatus: model.MessageStatus,
	// 		MessageEffect: model.MessageEffect,
	// 		Message:       model.Message,
	// 		MessagePlayed: "no",
	// 		SendTime:      now.Format("2006-01-02 15:04:05"),

	// 		Name:   user.Name,
	// 		Avatar: user.Avatar,
	// 	}))

	// // 將聊天資料加入redis中(時間排序，CHATROOM_ORDER_BY_TIME_REDIS)
	// m.RedisConn.ListRPush(config.CHATROOM_ORDER_BY_TIME_REDIS+model.ActivityID,
	// 	strconv.Itoa(int(id)))
	// }
	return err
}

// Update 更新聊天紀錄資料
func (m ChatroomRecordModel) Update(isRedis bool, model EditChatroomRecordModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"message_status", "message_played"}
	)

	if _, err := strconv.Atoi(model.ID); err != nil {
		return errors.New("錯誤: ID資料發生問題，請輸入有效的ID")
	}
	if model.MessageStatus != "" {
		if model.MessageStatus != "yes" && model.MessageStatus != "no" &&
			model.MessageStatus != "review" {
			return errors.New("錯誤: 訊息審核狀態資料發生問題，請輸入有效的訊息審核狀態")
		}
	}
	if model.MessagePlayed != "" {
		if model.MessagePlayed != "yes" && model.MessagePlayed != "no" {
			return errors.New("錯誤: 訊息播放狀態資料發生問題，請輸入有效的訊息播放狀態")
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
		Where("id", "=", model.ID).Update(fieldValues); err != nil {
		return err
	}

	// if isRedis {
	// 清除redis資料
	// m.RedisConn.DelCache(config.CHATROOM_REDIS + model.ActivityID)
	// m.Find(true, model.ActivityID)

	// var (
	// 	data                         string
	// 	record                       ChatroomRecordModel
	// 	messageStatus, messagePlayed string
	// )

	// // 取得聊天紀錄資料
	// data, err := m.RedisConn.HashGetCache(config.CHATROOM_REDIS+model.ActivityID,
	// 	model.ID)
	// if err == nil {
	// 	// 解碼
	// 	json.Unmarshal([]byte(data), &record)
	// }

	// // 無法從redis取得聊天資料時，刪除舊的redis聊天資訊並重新執行Find function
	// if record.ID == 0 {
	// 	m.RedisConn.DelCache(config.CHATROOM_REDIS + model.ActivityID)               // 提問資訊
	// 	m.RedisConn.DelCache(config.CHATROOM_ORDER_BY_TIME_REDIS + model.ActivityID) // 提問資訊(時間)
	// 	m.Find(true, model.ActivityID, "", "", "", 0, 0)
	// 	return nil
	// }

	// // 判斷是否更新該欄位，如果未更新則保持原始資料(record)
	// if model.MessageStatus != "" {
	// 	// 有更新欄位
	// 	messageStatus = model.MessageStatus
	// } else {
	// 	// 未更新欄位
	// 	messageStatus = record.MessageStatus
	// }
	// if model.MessagePlayed != "" {
	// 	// 有更新欄位
	// 	messagePlayed = model.MessagePlayed
	// } else {
	// 	// 未更新欄位
	// 	messagePlayed = record.MessagePlayed
	// }

	// // redis裡已有聊天資料，將更新的聊天資料加入redis
	// m.RedisConn.HashSetCache(config.CHATROOM_REDIS+model.ActivityID,
	// 	model.ID, utils.JSON(ChatroomRecordModel{
	// 		ID:            record.ID,
	// 		ActivityID:    record.ActivityID,
	// 		UserID:        record.UserID,
	// 		MessageType:   record.MessageType,
	// 		MessageStyle:  record.MessageStyle,
	// 		MessagePrice:  record.MessagePrice,
	// 		MessageStatus: messageStatus,
	// 		MessageEffect: record.MessageEffect,
	// 		Message:       record.Message,
	// 		MessagePlayed: messagePlayed,
	// 		SendTime:      record.SendTime,
	// 		Name:          record.Name,
	// 		Avatar:        record.Avatar,
	// 	}))
	// }
	return nil
}

// MapToModel map轉換ChatroomRecordModel
func (m ChatroomRecordModel) MapToModel(item map[string]interface{}) ChatroomRecordModel {
	// log.Println("查詢不同類型訊息數量")

	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &m)

	// log.Println("第一種過濾: ", m.YesMessageAmount, m.NoMessageAmount, m.ReviewMessageAmount)
	// log.Println("第二種過濾: ", m.YesPlayedAmount, m.NoPlayedAmount)
	// log.Println("第三種過濾: ", m.YesNormalMessageAmount, m.YesNormalBarrageAmount, m.YesSpecialBarrageAmount, m.YesOccupyBarrageAmount)
	// log.Println("第四種過濾: ", m.NoNormalMessageAmount, m.NoNormalBarrageAmount, m.NoSpecialBarrageAmount, m.NoOccupyBarrageAmount)
	// log.Println("第五種過濾: ", m.ReviewNormalMessageAmount, m.ReviewNormalBarrageAmount, m.ReviewSpecialBarrageAmount, m.ReviewOccupyBarrageAmount)
	// log.Println("第六種過濾: ", m.UserMessageAmount)

	return m
}

// MapToChatroomRecordModel map轉換[]ChatroomRecordModel
func MapToChatroomRecordModel(items []map[string]interface{}) []ChatroomRecordModel {
	var records = make([]ChatroomRecordModel, 0)
	for _, item := range items {
		var record ChatroomRecordModel

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &record)

		records = append(records, record)
	}
	return records
}

// UpdateUser 更新聊天紀錄人員姓名、頭像資料
// func (m ChatroomRecordModel) UpdateUser(userID, name, avatar string) error {
// 	var (
// 		fieldValues = command.Value{"name": name, "avatar": avatar}
// 	)
// 	if err := m.Table(m.Base.TableName).Where("user_id", "=", userID).
// 		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新聊天紀錄資料發生問題，請重新操作")
// 	}

// 	return nil
// }

// FindMessageAmount 查詢聊天訊息數量
// func (m ChatroomRecordModel) FindMessageAmount(activityID, userID, status, played, messageType string) (int64, error) {
// 	var (
// 		sql = m.Table(m.Base.TableName).
// 			Select("count(*)")
// 	)

// 	if activityID != "" {
// 		sql = sql.Where("activity_chatroom_record.activity_id", "=", activityID)
// 	}
// 	if userID != "" {
// 		sql = sql.Where("activity_chatroom_record.user_id", "=", userID)
// 	}
// 	if status != "" {
// 		sql = sql.WhereIn("activity_chatroom_record.message_status", interfaces(strings.Split(status, ",")))
// 	}
// 	if played != "" {
// 		sql = sql.WhereIn("activity_chatroom_record.message_played", interfaces(strings.Split(played, ",")))
// 	}
// 	if messageType != "" {
// 		sql = sql.WhereIn("activity_chatroom_record.message_type", interfaces(strings.Split(messageType, ",")))
// 	}

// 	item, err := sql.First()
// 	if err != nil {
// 		return 0, errors.New("錯誤: 無法取得聊天訊息數量，請重新查詢")
// 	}

// 	// log.Println("item: ", item)

// 	return item["count(*)"].(int64), nil
// }

// record.ID, _ = item["id"].(int64)
// record.UserID, _ = item["user_id"].(string)
// record.ActivityID, _ = item["activity_id"].(string)
// record.MessageType, _ = item["message_type"].(string)
// record.MessageStyle, _ = item["message_style"].(string)
// record.MessagePrice, _ = item["message_price"].(int64)
// record.MessageStatus, _ = item["message_status"].(string)
// record.MessageEffect, _ = item["message_effect"].(string)
// record.MessagePlayed, _ = item["message_played"].(string)
// record.Message, _ = item["message"].(string)
// record.SendTime, _ = item["send_time"].(string)

// 用戶資訊
// record.Name, _ = item["name"].(string)
// record.Avatar, _ = item["avatar"].(string)
// now, _ = time.ParseInLocation("2006-01-02 15:04:05",
// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
// price, err := strconv.Atoi(model.MessagePrice)
// if err != nil {
// 	return errors.New("錯誤: 訊息價格資料發生問題，請輸入有效的訊息價格")
// }

// command.Value{
// 	"activity_id":    model.ActivityID,
// 	"user_id":        model.UserID,
// 	"message_type":   model.MessageType,
// 	"message_style":  model.MessageStyle,
// 	"message_price":  model.MessagePrice,
// 	"message_status": model.MessageStatus,
// 	"message_effect": model.MessageEffect,
// 	"message":        model.Message,
// 	"message_played": "no",
// 	// "send_time":      now,
// }

// values      = []string{model.MessageStatus, model.MessagePlayed}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
