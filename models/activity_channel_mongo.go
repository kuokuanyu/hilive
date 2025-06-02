package models

import (
	"errors"
	"hilive/modules/config"
	"hilive/modules/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// SetMongoConn 設定connection
// func (a ActivityChannelModel) SetMongoConn(conn mongo.Connection) ActivityChannelModel {
// 	a.MongoConn = conn
// 	return a
// }

// FindByMongo 查詢頻道資訊，redis有資料先從reids取得，沒有則呼叫mongo並加入redis中
func (a ActivityChannelModel) FindByMongo(isRedis bool, activityID string) (ActivityChannelModel, error) {
	log.Println("查詢該活動所有頻道狀態(mongo)")

	if isRedis {
		var ()
		// 判斷redis裡是否有遊戲資訊，有則不執行查詢資料表功能
		dataMap, err := a.RedisConn.HashGetAllCache(config.HOST_CONTROL_CHANNEL_REDIS + activityID)
		if err != nil {
			return ActivityChannelModel{}, errors.New("錯誤: 取得頻道快取資料發生問題")
		}

		log.Println("redis? ", utils.GetInt64FromStringMap(dataMap, "id", 0), dataMap["mongo_id"])

		a.ID = utils.GetInt64FromStringMap(dataMap, "id", 0)
		a.MongoID = dataMap["mongo_id"]
		a.UserID = dataMap["user_id"]
		a.ActivityID = dataMap["activity_id"]
		a.ChannelAmount = utils.GetInt64FromStringMap(dataMap, "channel_amount", 0)
		a.Channel1 = dataMap["channel_1"]
		a.Channel2 = dataMap["channel_2"]
		a.Channel3 = dataMap["channel_3"]
		a.Channel4 = dataMap["channel_4"]
		a.Channel5 = dataMap["channel_5"]
		a.Channel6 = dataMap["channel_6"]
		a.Channel7 = dataMap["channel_7"]
		a.Channel8 = dataMap["channel_8"]
		a.Channel9 = dataMap["channel_9"]
		a.Channel10 = dataMap["channel_10"]
	}

	if a.ID == 0 {
		log.Println("redis無資料，執行mongo資料表查詢該活動頻道資料")

		item, err := a.MongoConn.FindOne(a.TableName, bson.M{"activity_id": activityID})
		if err != nil {
			return ActivityChannelModel{}, errors.New("錯誤: 無法取得頻道資訊(mongo)，請重新查詢")
		}
		if item == nil {
			return ActivityChannelModel{}, nil
		}

		a = a.MapToModel(item)
		// 將頻道資訊加入redis
		if isRedis {
			values := []interface{}{config.HOST_CONTROL_CHANNEL_REDIS + activityID}
			values = append(values, "id", a.ID)
			values = append(values, "mongo_id", a.MongoID)
			values = append(values, "user_id", a.UserID)
			values = append(values, "activity_id", a.ActivityID)
			values = append(values, "channel_amount", a.ChannelAmount)

			values = append(values, "channel_1", a.Channel1)
			values = append(values, "channel_2", a.Channel2)
			values = append(values, "channel_3", a.Channel3)
			values = append(values, "channel_4", a.Channel4)
			values = append(values, "channel_5", a.Channel5)
			values = append(values, "channel_6", a.Channel6)
			values = append(values, "channel_7", a.Channel7)
			values = append(values, "channel_8", a.Channel8)
			values = append(values, "channel_9", a.Channel9)
			values = append(values, "channel_10", a.Channel10)

			// 頻道資訊
			if err := a.RedisConn.HashMultiSetCache(values); err != nil {
				return a, errors.New("錯誤: 設置頻道快取資料發生問題")
			}

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+activityID, "修改資料")
		}
	}

	return a, nil
}

// UpdateByMongo 更新頻道資料(mongo)
func (a ActivityChannelModel) UpdateByMongo(isRedis bool, activityID string,
	channel string, status string) error {
	var (
		fieldValues = bson.M{}
		fields      = []string{channel}
		values      = []string{status}
	)

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}

	// 更新redis
	if isRedis {
		// 從redis取得資料，確定redis中有該活動的頻道資料
		a.FindByMongo(true, activityID)

		a.RedisConn.HashSetCache(config.HOST_CONTROL_CHANNEL_REDIS+activityID, channel, status)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+activityID, "修改資料")
	}

	// 更新資料庫
	_, err := a.MongoConn.UpdateOne(a.TableName, bson.M{"activity_id": activityID},
		bson.M{
			"$set": fieldValues,
		})
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 無法更新頻道資料(mongo)，請重新操作")
	}

	return nil
}
