package models

import (
	"errors"
	"hilive/modules/config"
	"hilive/modules/db/command"
)

// UpdateMessageAmount 更新活動簡訊數量(修改簡訊數量快取資料)
func (a ActivityModel) UpdateMessageAmount(isRedis bool, activityID string, amount int64) error {
	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", activityID).
		Where("message_amount", ">", 0).
		Update(command.Value{"message_amount": amount}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 其他錯誤
		return errors.New("錯誤: 更新簡訊數量發生錯誤, 請重新操作")
	}

	if isRedis {
		// 更新簡訊數量快取資料
		a.RedisConn.HashSetCache(config.ACTIVITY_REDIS+activityID, "message_amount", amount)
	}

	return nil
}

// DecrMessageAmount 遞減活動簡訊數量(修改簡訊數量快取資料)
func (a ActivityModel) DecrMessageAmount(isRedis bool, activityID string) error {
	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", activityID).
		Where("message_amount", ">", 0).
		Update(command.Value{"message_amount": "message_amount - 1"}); err != nil {
		if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
			// 數量為0
			return errors.New("錯誤: 數量為0, 不傳遞訊息")
		}
		// 其他錯誤
		return errors.New("錯誤: 減少簡訊數量發生錯誤, 請重新操作")
	}

	if isRedis {
		// 遞減簡訊數量快取資料
		a.RedisConn.HashDecrCache(config.ACTIVITY_REDIS+activityID, "message_amount")
	}

	return nil
}

// UpdateMailAmount 更新活動郵件數量(修改郵件數量快取資料)
func (a ActivityModel) UpdateMailAmount(isRedis bool, activityID string, amount int64) error {
	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", activityID).
		Where("mail_amount", ">", 0).
		Update(command.Value{"mail_amount": amount}); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 其他錯誤
		return errors.New("錯誤: 更新郵件數量發生錯誤, 請重新操作")
	}

	if isRedis {
		// 更新簡訊數量快取資料
		a.RedisConn.HashSetCache(config.ACTIVITY_REDIS+activityID, "mail_amount", amount)
	}

	return nil
}

// DecrMailAmount 遞減活動郵件數量(修改郵件數量快取資料)
func (a ActivityModel) DecrMailAmount(isRedis bool, activityID string) error {
	if err := a.Table(a.Base.TableName).
		Where("activity_id", "=", activityID).
		Where("mail_amount", ">", 0).
		Update(command.Value{"mail_amount": "mail_amount - 1"}); err != nil {
		if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
			// 數量為0
			return errors.New("錯誤: 數量為0, 不傳遞訊息")
		}
		// 其他錯誤
		return errors.New("錯誤: 減少郵件數量發生錯誤, 請重新操作")
	}

	if isRedis {
		// 遞減簡訊數量快取資料
		a.RedisConn.HashDecrCache(config.ACTIVITY_REDIS+activityID, "mail_amount")
	}

	return nil
}
