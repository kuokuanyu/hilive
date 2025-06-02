package models

import (
	"errors"
	"hilive/modules/config"
	"hilive/modules/db/command"
	"hilive/modules/utils"
)

// UpdateAttendAndNumber 更新活動人數與號碼資料
func (a ActivityModel) UpdateAttendAndNumber(isRedis bool, activityID string,
	attend, number int, users []string) error {
	var (
		fieldValues = command.Value{"attend": attend}
	)
	if number != 0 {
		// number資料不為0時更新number
		fieldValues["number"] = number
	}

	// 更新資料
	if err := a.Table(a.Base.TableName).
		WhereRaw("`activity_id` = ? and `attend` <= `people`", activityID).
		Update(fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 無法更新活動人數與號碼資料，請重新操作")
	}

	// 更新redis中的number資訊
	if isRedis && number != 0 {
		// 將number資料寫入redis中
		_, err := a.RedisConn.SetCache(config.ACTIVITY_NUMBER_REDIS+activityID, number)
		if err != nil {
			return errors.New("錯誤: 將活動number資料寫入redis發生問題")
		}

		// 設置過期時間
		// a.RedisConn.SetExpire(config.ACTIVITY_NUMBER_REDIS+activityID, config.REDIS_EXPIRE)

		if len(users) > 0 {
			// 從redis取得資料，確定redis中有該場活動報名簽到人員資料(sign_staffs_2_activityID)
			DefaultApplysignModel().
				SetConn(a.DbConn, a.RedisConn, a.MongoConn).
				FindSignStaffs(true, false, config.SIGN_STAFFS_2_REDIS, activityID, 0, 0)

			var params = []interface{}{config.SIGN_STAFFS_2_REDIS + activityID}

			// 將用戶資料加入redis中(新增時處理)
			for _, userID := range users {
				params = append(params, userID)
			}

			// 將簽到人員資訊加入redis中(SET)
			a.RedisConn.SetAdd(params)
		}
	}
	return nil
}

// IncrAttendAndNumber 遞增活動參加人數及號碼資訊(修改活動人數快取資料)
func (a ActivityModel) IncrAttendAndNumber(isRedis bool, activityID, userID string) error {
	if err := a.Table(a.Base.TableName).
		WhereRaw("`activity_id` = ? and `attend` < `people`", activityID).
		Update(command.Value{"attend": "attend + 1", "number": "number + 1"}); err != nil {
		if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
			return errors.New("錯誤: 參加人數已達上限，如要參加活動，請聯絡主辦方")
		}
		return err
	}

	if isRedis {
		// 判斷redis裡是否有活動number資訊
		// data, err := a.RedisConn.GetCache(config.ACTIVITY_NUMBER_REDIS + activityID)
		// if err != nil {
		// 	return errors.New("錯誤: 取得活動number資料發生問題(redis)")
		// }

		// number := utils.GetInt64(data, 0)
		// if number != 0 {
		// 	// redis中已存在活動number快取資訊，遞增number資訊
		// 	a.RedisConn.IncrCache(config.ACTIVITY_NUMBER_REDIS + activityID)
		// }

		// 從redis取得資料，確定redis中有該場活動的number資料(activity_number_activityID)
		DefaultActivityModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			GetActivityNumber(true, activityID)

		a.RedisConn.IncrCache(config.ACTIVITY_NUMBER_REDIS + activityID)

		// 從redis取得資料，確定redis中有該場活動報名簽到人員資料(sign_staffs_2_activityID)
		DefaultApplysignModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			FindSignStaffs(true, false, config.SIGN_STAFFS_2_REDIS, activityID, 0, 0)

		// redis中不存在簽到人員快取資訊，新增簽到人員資訊
		if err := a.RedisConn.SetAdd([]interface{}{config.SIGN_STAFFS_2_REDIS + activityID, userID}); err != nil {
			return errors.New("錯誤: 新增簽到人員快取資料發生問題(sign_staffs_2_activityID)")
		}
	}
	return nil
}

// GetActivityNumber 取得活動的number資料(redis)
func (a ActivityModel) GetActivityNumber(isRedis bool, activityID string) (int64, error) {
	var (
		number int64
	)

	if isRedis {
		// 判斷redis裡是否有活動number資料，有則不執行查詢資料表功能
		data, err := a.RedisConn.GetCache(config.ACTIVITY_NUMBER_REDIS + activityID)
		if err != nil {
			return number, errors.New("錯誤: 取得活動number資料發生問題(redis)")
		}

		number = utils.GetInt64(data, 0)
		if number == 0 {
			// 活動number參數不可能為0，從資料表取的number資料
			activityModel, err := a.Table(a.Base.TableName).
				Select("number").
				Where("activity_id", "=", activityID).
				First()
			if err != nil {
				return number, errors.New("錯誤: 取得活動number資料發生問題(資料表)")
			}

			number = utils.GetInt64(activityModel["number"], 0)

			if isRedis && number != 0 {

				// 將number資料寫入redis中
				_, err = a.RedisConn.SetCache(config.ACTIVITY_NUMBER_REDIS+activityID, number)
				if err != nil {
					return number, errors.New("錯誤: 將活動number資料寫入redis發生問題")
				}

				// 設置過期時間
				// a.RedisConn.SetExpire(config.ACTIVITY_NUMBER_REDIS+activityID, config.REDIS_EXPIRE)
			}
		}

	}

	return number, nil
}