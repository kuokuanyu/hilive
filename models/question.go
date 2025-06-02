package models

import (
	"errors"
	"unicode/utf8"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// EditQuestionModel 資料表欄位
type EditQuestionModel struct {
	ActivityID         string `json:"activity_id"`
	QuestionAnonymous  string `json:"question_anonymous"`
	QuestionQrcode     string `json:"question_qrcode"`
	QuestionBackground string `json:"question_background"`
	// QuestionMessageCheck  string `json:"question_message_check"`
	// QuestionHideAnswered  string `json:"question_hide_answered"`
	// QuestionGuestAnswered string `json:"question_guest_answered"`
}

// QuestionGuestModel 資料表欄位
type QuestionGuestModel struct {
	Base       `json:"-"`
	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

// EditQuestionGuestModel 資料表欄位
type EditQuestionGuestModel struct {
	ID         string `json:"id"`
	ActivityID string `json:"activity_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

// DefaultQuestionGuestModel 預設QuestionGuestModel
func DefaultQuestionGuestModel() QuestionGuestModel {
	return QuestionGuestModel{Base: Base{TableName: config.ACTIVITY_QUESTION_GUEST_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (g QuestionGuestModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) QuestionGuestModel {
	g.DbConn = dbconn
	g.RedisConn = cacheconn
	g.MongoConn = mongoconn
	return g
}

// SetDbConn 設定connection
// func (g QuestionGuestModel) SetDbConn(conn db.Connection) QuestionGuestModel {
// 	g.DbConn = conn
// 	return g
// }

// Add 增加資料
func (g QuestionGuestModel) Add(model EditQuestionGuestModel) error {
	if utf8.RuneCountInString(model.Name) > 10 {
		return errors.New("錯誤: 姓名上限為10個字元，請輸入有效的提問嘉賓姓名")
	}

	_, err := g.Table(g.TableName).Insert(command.Value{
		"activity_id": model.ActivityID,
		"name":        model.Name,
		"avatar":      model.Avatar,
	})
	return err
}

// UpdateQuestion 更新提問牆基本設置資料
func (a ActivityModel) UpdateQuestion(model EditQuestionModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"question_anonymous",
			"question_qrcode",
			"question_background"}
	)

	if model.QuestionAnonymous != "" {
		if model.QuestionAnonymous != "open" && model.QuestionAnonymous != "close" {
			return errors.New("錯誤: 匿名提問資料發生問題，請輸入有效的資料")
		}
	}

	if model.QuestionQrcode != "" {
		if model.QuestionQrcode != "open" && model.QuestionQrcode != "close" {
			return errors.New("錯誤: 二維碼資料發生問題，請輸入有效的資料")
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
	return a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// Update 更新提問嘉賓資料
func (g QuestionGuestModel) Update(model EditQuestionGuestModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"name", "avatar"}
		values      = []string{model.Name, model.Avatar}
	)
	if utf8.RuneCountInString(model.Name) > 10 {
		return errors.New("錯誤: 姓名上限為10個字元，請輸入有效的提問嘉賓姓名")
	}

	for i, value := range values {
		if value != "" {
			fieldValues[fields[i]] = value
		}
	}
	if len(fieldValues) == 0 {
		return nil
	}
	return g.Table(g.Base.TableName).Where("activity_id", "=", model.ActivityID).
		Where("id", "=", model.ID).Update(fieldValues)
}

// DecrLikes 遞減按讚數
// func (g QuestionUserModel) DecrLikes(id int64, activityID string) error {
// 	if err := g.Table(g.Base.TableName).Where("activity_id", "=", activityID).
// 		Where("id", "=", id).Where("likes", ">", 0).
// 		Update(command.Value{"likes": "likes - 1"}); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// 	return nil
// }

// Name       string `json:"name"`
// Avatar     string `json:"avatar"`
// Name       string `json:"name"`
// Avatar     string `json:"avatar"`

// QuestionsOrderByLikes 尋找該活動所有提問資料(讚數排序)
// func (g QuestionUserModel) QuestionsOrderByLikes(isRedis bool,
// 	activityID string) (questions []QuestionUserModel, err error) {
// 	// var (
// 	// 	item []map[string]interface{}
// 	// 	// err  error
// 	// )
// 	// 判斷redis裡是否有提問資訊
// 	if isRedis {
// 		var ids = make([]string, 0)
// 		if ids, err = g.RedisConn.ZSetRevRange(config.QUESTION_LIKES_REDIS+activityID,
// 			0, 1000); err != nil {
// 			return questions, errors.New("錯誤: 從redis中取得讚數由高至低的提問資訊發生問題")
// 		}

// 		for _, id := range ids {
// 			var question QuestionUserModel
// 			data, err := g.RedisConn.HashGetCache(config.QUESTION_TIME_REDIS+activityID, id)
// 			if err != nil {
// 				return questions, errors.New("錯誤: 從redis中取得提問資訊發生問題")
// 			}
// 			// 解碼
// 			json.Unmarshal([]byte(data), &question)

// 			questions = append(questions, question)
// 		}
// 	}

// 	// redis中無提問快取資料，從資料表中查詢
// 	if len(questions) == 0 {
// 		item, err := g.Table(g.Base.TableName).
// 			Select("activity_question_user.id", "activity_question_user.activity_id",
// 				"activity_question_user.user_id", "activity_question_user.content",
// 				"activity_question_user.send_time",
// 				"activity_question_user.likes", "activity_question_user.like",

// 				"line_users.name", "line_users.avatar").
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_question_user.user_id",
// 				FieldB:    "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			Where("activity_question_user.activity_id", "=", activityID).
// 			OrderBy("activity_question_user.likes", "desc").All()
// 		if err != nil {
// 			return nil, errors.New("錯誤: 無法取得用戶提問資料，請重新查詢")
// 		}

// 		questions = MapToQuestionUserModel(item)

// 		// 將提問資訊加入redis中
// 		if isRedis && len(questions) > 0 {
// 			for _, question := range questions {
// 				// 將按讚資料加入redis中
// 				g.RedisConn.ZSetAdd(config.QUESTION_LIKES_REDIS+activityID,
// 					strconv.Itoa(int(question.ID)), question.Likes)
// 			}
// 			// 設置過期時間
// 			g.RedisConn.SetExpire(config.QUESTION_LIKES_REDIS+activityID, config.REDIS_EXPIRE)
// 		}
// 	}
// 	return
// }

// LeftJoinLineUsersByUser 尋找該用戶在活動中所有提問資料
// func (g QuestionUserModel) LeftJoinLineUsersByUser(activityID, userID string) ([]QuestionUserModel, error) {
// 	var (
// 		item []map[string]interface{}
// 		err  error
// 	)

// 	item, err = g.Table(g.Base.TableName).
// 		Select("activity_question_user.id", "activity_question_user.activity_id",
// 			"activity_question_user.user_id", "activity_question_user.content",
// 			"activity_question_user.send_time",
// 			"activity_question_user.likes", "activity_question_user.like",

// 			"line_users.name", "line_users.avatar").
// 		LeftJoin(command.Join{
// 			FieldA:    "activity_question_user.user_id",
// 			FieldB:    "line_users.user_id",
// 			Table:     "line_users",
// 			Operation: "="}).
// 		Where("activity_question_user.activity_id", "=", activityID).
// 		Where("activity_question_user.user_id", "=", userID).
// 		OrderBy("activity_question_user.send_time", "desc").All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得用戶提問資料，請重新查詢")
// 	}
// 	return MapToQuestionUserModel(item), nil
// }

// "name":        model.Name,
// "avatar":      model.Avatar,
// fields      = []string{"question_message_check", "question_anonymous",
// 	"question_hide_answered", "question_qrcode", "question_background",
// 	"question_guest_answered"}
// values = []string{model.QuestionMessageCheck, model.QuestionAnonymous,
// 	model.QuestionHideAnswered, model.QuestionQrcode,
// 	model.QuestionBackground, model.QuestionGuestAnswered}
// if model.QuestionMessageCheck != "" {
// 	if model.QuestionMessageCheck != "open" && model.QuestionMessageCheck != "close" {
// 		return errors.New("錯誤: 提問訊息審核資料發生問題，請輸入有效的資料")
// 	}
// }
// if model.QuestionHideAnswered != "" {
// 	if model.QuestionHideAnswered != "open" && model.QuestionHideAnswered != "close" {
// 		return errors.New("錯誤: 隱藏已解答問題資料發生問題，請輸入有效的資料")
// 	}
// }
// if model.QuestionGuestAnswered != "" {
// 	if model.QuestionGuestAnswered != "open" && model.QuestionGuestAnswered != "close" {
// 		return errors.New("錯誤: 嘉賓解答資料發生問題，請輸入有效的資料")
// 	}
// }

// UpdateUser 更新提問訊息姓名、頭像資料
// func (g QuestionUserModel) UpdateUser(userID, name, avatar string) error {
// 	var (
// 		fieldValues = command.Value{"name": name, "avatar": avatar}
// 	)
// 	if err := g.Table(g.Base.TableName).Where("user_id", "=", userID).
// 		Update(fieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新提問訊息人員資料發生問題，請重新操作")
// 	}
// 	return nil
// }

// UpdateUser 更新按讚紀錄人員姓名、頭像資料
// func (g QuestionLikesRecordModel) UpdateUser(userID, name, avatar string) error {
// 	var (
// 		fieldValues = command.Value{"name": name, "avatar": avatar}
// 	)
// 	if err := g.Table(g.Base.TableName).Where("u

// values      = []string{model.QuestionAnonymous, model.QuestionQrcode, model.QuestionBackground}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
