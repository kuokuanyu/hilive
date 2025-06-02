package models

import (
	"hilive/modules/db/command"
)

// EditTopicModel 資料表欄位
type EditTopicModel struct {
	ActivityID      string `json:"activity_id"`
	TopicBackground string `json:"topic_background"`
}

// UpdateTopic 更新主題牆基本設置資料
func (a ActivityModel) UpdateTopic(model EditTopicModel) error {
	if model.TopicBackground != "" {
		fieldValues := command.Value{
			"topic_background": model.TopicBackground,
		}
		
		return a.Table(a.Base.TableName).
			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
	}
	return nil
}
