package models

import (
	"errors"
	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateActivity 更新活動資料
func (a ActivityModel) UpdateActivity(isRedis bool, model EditActivityModel) error {
	var (
		// 第一步：一定要加入的欄位（即使為空）
		fieldValues = map[string]interface{}{
			// 官方帳號(可能會清空數值)
			"line_id":        model.LineID,
			"channel_id":     model.ChannelID,
			"channel_secret": model.ChannelSecret,
			"chatbot_secret": model.ChatbotSecret,
			"chatbot_token":  model.ChatbotToken,
		}

		fields = []string{
			"activity_name", "activity_type",
			"login_required", "login_password", "password_required",
			"max_people", "people",
			"push_message",
			"device",
			"message_amount", "push_phone_message",
			"send_mail", "mail_amount",
			"host_scan",
		}

		fieldValues2 = command.Value{}
		fields2      = []string{
			"channel_amount",
		}

		activityTypes = []string{"企業會議", "其他", "商業活動", "培訓/教育", "婚禮", "尾牙春酒",
			"校園活動", "競技賽事", "論壇會議", "酒吧/餐飲娛樂", "電視/媒體"}
		Isbool bool
	)

	// 欄位值判斷
	if utf8.RuneCountInString(model.ActivityName) > 20 {
		return errors.New("錯誤: 活動名稱上限為20個字元，請輸入有效的活動名稱")
	}
	if model.ActivityType != "" {
		for i := range activityTypes {
			if model.ActivityType == activityTypes[i] {
				Isbool = true
				break
			}
		}
		if Isbool == false {
			return errors.New("錯誤: 活動類型發生問題，請輸入有效的活動類型")
		}
	}
	if model.City != "" && model.Town != "" {
		fieldValues["city"] = model.City
		fieldValues["town"] = model.Town
	}
	if model.StartTime != "" && model.EndTime != "" {
		if !CompareTDatetime(model.StartTime, model.EndTime) {
			return errors.New("錯誤: 活動時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
		}
		fieldValues["start_time"] = model.StartTime
		fieldValues["end_time"] = model.EndTime
	}
	if model.LoginRequired != "" {
		if model.LoginRequired != "open" && model.LoginRequired != "close" {
			return errors.New("錯誤: 是否需要登入才可進入聊天室資料發生問題，請輸入有效的資料")
		}
	}
	if model.PasswordRequired != "" {
		if model.PasswordRequired != "open" && model.PasswordRequired != "close" {
			return errors.New("錯誤: 是否需要設置密碼才可進入聊天室資料發生問題，請輸入有效的資料")
		}
	}
	if utf8.RuneCountInString(model.LoginPassword) > 8 {
		return errors.New("錯誤: 聊天室驗證碼上限為8個字元，請輸入有效的密碼")
	}

	if model.PushMessage != "" {
		if model.PushMessage != "open" && model.PushMessage != "close" {
			return errors.New("錯誤: 官方帳號發送訊息資料發生問題，請輸入有效的資料")
		}
	}

	if model.PushPhoneMessage != "" {
		if model.PushPhoneMessage != "open" && model.PushPhoneMessage != "close" {
			return errors.New("錯誤: 發送手機簡訊資料發生問題，請輸入有效的資料")
		}
	}

	if model.SendMail != "" {
		if model.SendMail != "open" && model.SendMail != "close" {
			return errors.New("錯誤: 發送郵件資料發生問題，請輸入有效的資料")
		}
	}

	if model.ChannelAmount != "" {
		if _, err := strconv.Atoi(model.ChannelAmount); err != nil {
			return errors.New("錯誤: 頻道數量資料發生問題，請輸入有效的頻道數量")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) != 0 {
		if model.People != "" {
			// 有更新人數資訊，必須有人數判斷的條件
			if err := a.Table(a.Base.TableName).
				Where("attend", "<=", model.People).
				Where("activity_id", "=", model.ActivityID).
				Update(fieldValues); err != nil {
				if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
					return errors.New("錯誤: 活動人數上限數值必須大於目前的參加人數")
				}
				return errors.New("錯誤: 無法更新活動資料，請重新操作")
			}
		} else {
			if err := a.Table(a.Base.TableName).
				Where("activity_id", "=", model.ActivityID).
				Update(fieldValues); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return errors.New("錯誤: 無法更新活動資料，請重新操作")
			}
		}
	}

	for _, key := range fields2 {
		if val, ok := data[key]; ok && val != "" {
			fieldValues2[key] = val
		}
	}

	if len(fieldValues2) != 0 {
		log.Println("更新活動api, 更新頻道數量(mongo): ", model.ChannelAmount)

		// mongo
		filter := bson.M{"activity_id": model.ActivityID}
		update := bson.M{
			"$set": fieldValues2,
			// "$unset": bson.M{}, // 移除不需要的欄位
		}

		_, err := a.MongoConn.UpdateOne(config.ACTIVITY_CHANNEL_TABLE, filter, update)
		if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
	}

	// 新增權限
	if model.Permissions != "" {
		a.ActivityID = model.ActivityID
		if err := a.AddPermission(strings.Split(model.Permissions, ",")); err != nil {
			return err
		}
	}

	// 編輯活動時需要清除用戶redis資訊(重新判斷活動權限資訊)
	if isRedis {
		a.RedisConn.DelCache(config.HILIVE_USERS_REDIS + model.UserID)
		a.RedisConn.DelCache(config.ACTIVITY_REDIS + model.ActivityID)

		a.RedisConn.DelCache(config.HOST_CONTROL_CHANNEL_REDIS + model.ActivityID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_HOST_CONTROL_CHANNEL_REDIS+model.ActivityID, "修改資料")
	}
	return nil
}

// mysql
// err := a.Table(config.ACTIVITY_CHANNEL_TABLE).
// 	Where("activity_id", "=", model.ActivityID).Update(fieldValues2)
// if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return err
// }

// values2 = []string{model.ChannelAmount}

// if len(fieldValues) == 0 {
// 	return nil
// }

// values = []string{model.ActivityName, model.ActivityType,
// 	model.LoginRequired, model.LoginPassword, model.PasswordRequired,
// 	model.MaxPeople, model.People,
// 	model.PushMessage,
// 	model.Device,
// 	model.MessageAmount, model.PushPhoneMessage,
// 	model.SendMail, model.MailAmount,
// 	model.HostScan,
// 	// 官方帳號
// 	// model.LineID, model.ChannelID, model.ChannelSecret, model.ChatbotSecret, model.ChatbotToken,
// }

// if model.People != "" {
// 	// 判斷活動人數上限
// 	maxPeopleInt, err1 := strconv.Atoi(model.MaxPeople)
// 	peopleInt, err2 := strconv.Atoi(model.People)
// 	if err1 != nil || err2 != nil || peopleInt > maxPeopleInt {
// 		return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
// 	}
// 	fieldValues["max_people"] = model.MaxPeople
// 	fieldValues["people"] = model.People
// 	// 判斷活動人數上限
// 	// 用戶資訊
// 	// userModel, err := DefaultUserModel().SetDbConn(a.DbConn).
// 	// 	SetRedisConn(a.RedisConn).Find(true, true, "",
// 	// 	"users.user_id", model.UserID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// _, err := strconv.Atoi(model.People)
// 	// if err != nil {
// 	// 	return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
// 	// }
// 	// if err != nil || peopleInt > int(userModel.MaxActivityPeople) {
// 	// 	return errors.New("錯誤: 活動人數上限資料發生問題，請輸入有效的活動人數上限")
// 	// }
// 	// if peopleInt, err := strconv.Atoi(model.People); peopleInt > 10000 || err != nil {
// 	// 	return errors.New("錯誤: 活動人數上限為10000，請輸入有效的活動人數")
// 	// }
// }
