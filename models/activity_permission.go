package models

import (
	"errors"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
)

// AddPermission 增加權限(多個權限)
func (a ActivityModel) AddPermission(ids []string) error {
	if len(ids) > 0 {
		// 先刪除原有的活動權限資料
		if err := db.Table(config.ACTIVITY_PERMISSIONS_TABLE).
			WithConn(a.DbConn).Where("activity_id", "=", a.ActivityID).
			Delete(); err != nil && err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
			return errors.New("錯誤: 刪除活動權限資料發生問題，請重新操作")
		}

		// 加入新的權限資料
		for _, id := range ids {
			if _, err := db.Table(config.ACTIVITY_PERMISSIONS_TABLE).
				WithConn(a.DbConn).Insert(command.Value{
				"permission_id": id, "activity_id": a.ActivityID,
			}); err != nil {
				return errors.New("錯誤: 新增活動權限資料發生問題，請重新操作")
			}
		}
	}

	return nil
}
