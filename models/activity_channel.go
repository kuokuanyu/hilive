package models

import (
	"encoding/json"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
)

// ActivityChannelModel 資料表欄位
type ActivityChannelModel struct {
	// 活動
	Base          `json:"-"`
	ID            int64  `json:"id"`
	MongoID       string `json:"_id"`
	ActivityID    string `json:"activity_id"`
	UserID        string `json:"user_id"`
	ChannelAmount int64  `json:"channel_amount"` // 頻道數量
	Channel1      string `json:"channel_1"`
	Channel2      string `json:"channel_2"`
	Channel3      string `json:"channel_3"`
	Channel4      string `json:"channel_4"`
	Channel5      string `json:"channel_5"`
	Channel6      string `json:"channel_6"`
	Channel7      string `json:"channel_7"`
	Channel8      string `json:"channel_8"`
	Channel9      string `json:"channel_9"`
	Channel10     string `json:"channel_10"`
}

// DefaultActivityChannelModel 預設ActivityChannelModel
func DefaultActivityChannelModel() ActivityChannelModel {
	return ActivityChannelModel{Base: Base{TableName: config.ACTIVITY_CHANNEL_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a ActivityChannelModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ActivityChannelModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a ActivityChannelModel) SetDbConn(conn db.Connection) ActivityChannelModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a ActivityChannelModel) SetRedisConn(conn cache.Connection) ActivityChannelModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (m ActivityChannelModel) SetMongoConn(conn mongo.Connection) ActivityChannelModel {
// 	m.MongoConn = conn
// 	return m
// }

// MapToModel 將值設置至GameModel
func (a ActivityChannelModel) MapToModel(m map[string]interface{}) ActivityChannelModel {

	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	return a
}

// Find 查詢頻道資訊，redis有資料先從reids取得，沒有則呼叫資料表並加入redis中
// func (a ActivityChannelModel) Find(isRedis bool, activityID string) (ActivityChannelModel, error) {
// 	// log.Println("查詢該活動所有頻道狀態")

// 	if isRedis {
// 		var ()
// 		// 判斷redis裡是否有遊戲資訊，有則不執行查詢資料表功能
// 		dataMap, err := a.RedisConn.HashGetAllCache(config.HOST_CONTROL_CHANNEL_REDIS + activityID)
// 		if err != nil {
// 			return ActivityChannelModel{}, errors.New("錯誤: 取得頻道快取資料發生問題")
// 		}

// 		a.ID = utils.GetInt64FromStringMap(dataMap, "id", 0)
// 		a.UserID = dataMap["user_id"]
// 		a.ActivityID = dataMap["activity_id"]
// 		a.ChannelAmount = utils.GetInt64FromStringMap(dataMap, "channel_amount", 0)
// 		a.Channel1 = dataMap["channel_1"]
// 		a.Channel2 = dataMap["channel_2"]
// 		a.Channel3 = dataMap["channel_3"]
// 		a.Channel4 = dataMap["channel_4"]
// 		a.Channel5 = dataMap["channel_5"]
// 		a.Channel6 = dataMap["channel_6"]
// 		a.Channel7 = dataMap["channel_7"]
// 		a.Channel8 = dataMap["channel_8"]
// 		a.Channel9 = dataMap["channel_9"]
// 		a.Channel10 = dataMap["channel_10"]
// 	}

// 	if a.ID == 0 {
// 		// log.Println("redis無資料，執行資料表查詢該活動頻道資料")

// 		var (
// 			sql = a.Table(a.Base.TableName).
// 				Where("activity_channel.activity_id", "=", activityID)
// 		)

// 		// fmt.Println("使用到資料庫!!!!!!!!!!!!!!!!!!!!!")
// 		item, err := sql.First()
// 		if err != nil {
// 			return ActivityChannelModel{}, errors.New("錯誤: 無法取得頻道資訊，請重新查詢")
// 		}
// 		if item == nil {
// 			return ActivityChannelModel{}, nil
// 		}

// 		a = a.MapToModel(item)
// 		// 將頻道資訊加入redis
// 		if isRedis {
// 			values := []interface{}{config.HOST_CONTROL_CHANNEL_REDIS + activityID}
// 			values = append(values, "id", a.ID)
// 			values = append(values, "user_id", a.UserID)
// 			values = append(values, "activity_id", a.ActivityID)
// 			values = append(values, "channel_amount", a.ChannelAmount)

// 			values = append(values, "channel_1", a.Channel1)
// 			values = append(values, "channel_2", a.Channel2)
// 			values = append(values, "channel_3", a.Channel3)
// 			values = append(values, "channel_4", a.Channel4)
// 			values = append(values, "channel_5", a.Channel5)
// 			values = append(values, "channel_6", a.Channel6)
// 			values = append(values, "channel_7", a.Channel7)
// 			values = append(values, "channel_8", a.Channel8)
// 			values = append(values, "channel_9", a.Channel9)
// 			values = append(values, "channel_10", a.Channel10)

// 			// 頻道資訊
// 			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
// 				return a, errors.New("錯誤: 設置頻道快取資料發生問題")
// 			}

// 			// 設置過期時間
// 			// a.RedisConn.SetExpire(config.HOST_CONTROL_CHANNEL_REDIS+activityID,
// 			// 	config.REDIS_EXPIRE)

// 			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 			a.RedisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+activityID, "修改資料")
// 		}
// 	}

// 	return a, nil
// }

// Update 更新頻道資料
// func (a ActivityChannelModel) Update(isRedis bool, activityID string,
// 	channel string, status string) error {
// 	var (
// 		fieldValues = command.Value{}
// 		fields      = []string{channel}
// 		values      = []string{status}
// 	)

// 	for i, value := range values {
// 		if value != "" {
// 			fieldValues[fields[i]] = value
// 		}
// 	}
// 	if len(fieldValues) == 0 {
// 		return nil
// 	}

// 	// 更新redis
// 	if isRedis {
// 		// 從redis取得資料，確定redis中有該活動的頻道資料
// 		a.Find(true, activityID)

// 		a.RedisConn.HashSetCache(config.HOST_CONTROL_CHANNEL_REDIS+activityID, channel, status)

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 		a.RedisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+activityID, "修改資料")
// 	}

// 	// 更新資料庫
// 	err := a.Table(a.Base.TableName).
// 		Where("activity_id", "=", activityID).
// 		Update(fieldValues)
// 	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}

// 	return nil
// }

// log.Println("頻道換個方式處理")

// log.Println("處理前: ", a)
// log.Println("m: ", m)

// log.Println("處理後: ", a)

// a.ID, _ = m["id"].(int64)
// a.UserID, _ = m["user_id"].(string)
// a.ActivityID, _ = m["activity_id"].(string)
// a.ChannelAmount, _ = m["channel_amount"].(int64)
// a.Channel1, _ = m["channel_1"].(string)
// a.Channel2, _ = m["channel_2"].(string)
// a.Channel3, _ = m["channel_3"].(string)
// a.Channel4, _ = m["channel_4"].(string)
// a.Channel5, _ = m["channel_5"].(string)
// a.Channel6, _ = m["channel_6"].(string)
// a.Channel7, _ = m["channel_7"].(string)
// a.Channel8, _ = m["channel_8"].(string)
// a.Channel9, _ = m["channel_9"].(string)
// a.Channel10, _ = m["channel_10"].(string)
