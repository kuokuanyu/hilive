package table

import (
	"encoding/json"
	"errors"

	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
)

// GetChatroomRecordPanel 聊天室紀錄
func (s *SystemTable) GetChatroomRecordPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	formList := table.GetForm()
	// 新增聊天紀錄
	formList.SetTable(config.ACTIVITY_CHATROOM_RECORD_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("user_id", "activity_id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditChatroomRecordModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultChatroomRecordModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, model); err != nil {
			return err
		}
		return nil
	})

	// 編輯聊天紀錄
	formList.SetTable(config.ACTIVITY_CHATROOM_RECORD_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "activity_id") {
			return errors.New("錯誤: ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditChatroomRecordModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultChatroomRecordModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetQuestionRecordPanel 提問紀錄
func (s *SystemTable) GetQuestionRecordPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	formList := table.GetForm()

	// 編輯提問紀錄
	formList.SetTable(config.ACTIVITY_QUESTION_USER_TABLE).
		SetUpdateFunc(func(values form2.Values) error {
			if values.IsEmpty("id", "activity_id") {
				return errors.New("錯誤: ID發生問題，請輸入有效的資料")
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditQuestionUserModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultQuestionUserModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Update(true, model); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}
			return nil
		})
	return
}

// models.EditChatroomRecordModel{
// 	ActivityID:    values.Get("activity_id"),
// 	UserID:        values.Get("user_id"),
// 	MessageType:   values.Get("message_type"),
// 	MessageStyle:  values.Get("message_style"),
// 	MessagePrice:  values.Get("message_price"),
// 	MessageStatus: values.Get("message_status"),
// 	MessageEffect: values.Get("message_effect"),
// 	Message:       values.Get("message"),
// }

// models.EditChatroomRecordModel{
// 	ID:            values.Get("id"),
// 	ActivityID:    values.Get("activity_id"),
// 	MessageStatus: values.Get("message_status"),
// 	MessagePlayed: values.Get("message_played"),
// }

// models.EditQuestionUserModel{
// 	ID:            values.Get("id"),
// 	ActivityID:    values.Get("activity_id"),
// 	MessageStatus: values.Get("message_status"),
// }
