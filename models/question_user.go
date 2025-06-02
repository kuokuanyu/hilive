package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
// questionLock sync.Mutex
)

// QuestionUserModel 資料表欄位
type QuestionUserModel struct {
	Base          `json:"-"`
	ID            int64  `json:"id" example:"1"`
	ActivityID    string `json:"activity_id" example:"activity_id"`
	UserID        string `json:"user_id" example:"user_id"`
	Message       string `json:"message" example:"message"`
	MessageStatus string `json:"message_status" example:"yes"`
	Likes         int64  `json:"likes" example:"1"`
	SendTime      string `json:"send_time" example:"2022/01/01 00:00:00"`
	Like          string `json:"like" example:"yes、no"`

	// 用戶資訊
	Name   string `json:"name" example:"name"`
	Avatar string `json:"avatar" example:"avatar"`

	// 過濾參數
	YesQuestionAmount    int64 `json:"yes_question_amount" example:"10"`    // 已通過提問數量
	NoQuestionAmount     int64 `json:"no_question_amount" example:"10"`     // 未通過提問數量
	ReviewQuestionAmount int64 `json:"review_question_amount" example:"10"` // 審核提問數量

	UserQuestionAmount int64 `json:"user_question_amount" example:"10"` // 該用戶提問數量

	QuestionLikesAmount int64 `json:"question_likes_amount" example:"10"` // 所有提問總讚數數量

	HostQuestionLikeAmount   int64 `json:"host_question_like_amount" example:"10"`   // 主持人已按讚提問數量
	HostQuestionUnlikeAmount int64 `json:"host_question_unlike_amount" example:"10"` // 主持人未按讚提問數量
}

// EditQuestionUserModel 資料表欄位
type EditQuestionUserModel struct {
	ID            string `json:"id"`
	ActivityID    string `json:"activity_id"`
	UserID        string `json:"user_id"`
	Message       string `json:"message"`
	MessageStatus string `json:"message_status"`
}

// DefaultQuestionUserModel 預設QuestionUserModel
func DefaultQuestionUserModel() QuestionUserModel {
	return QuestionUserModel{Base: Base{TableName: config.ACTIVITY_QUESTION_USER_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (g QuestionUserModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) QuestionUserModel {
	g.DbConn = dbconn
	g.RedisConn = cacheconn
	g.MongoConn = mongoconn
	return g
}

// SetDbConn 設定connection
// func (g QuestionUserModel) SetDbConn(conn db.Connection) QuestionUserModel {
// 	g.DbConn = conn
// 	return g
// }

// // SetRedisConn 設定connection
// func (g QuestionUserModel) SetRedisConn(conn cache.Connection) QuestionUserModel {
// 	g.RedisConn = conn
// 	return g
// }

// FindHostQuestionAmount 查詢不同類型提問訊息數量(主持端)
func (g QuestionUserModel) FindHostQuestionAmount(activityID string) (QuestionUserModel, error) {
	var (
		sql = g.Table(g.Base.TableName).
			Select(
				"SUM(CASE WHEN message_status = 'yes' THEN 1 ELSE 0 END) AS yes_question_amount",
				"SUM(CASE WHEN message_status = 'no' THEN 1 ELSE 0 END) AS no_question_amount",
				"SUM(CASE WHEN message_status = 'review' THEN 1 ELSE 0 END) AS review_question_amount",

				"SUM(likes) AS question_likes_amount",

				"SUM(CASE WHEN `like` = 'yes' THEN 1 ELSE 0 END) AS host_question_like_amount",
				"SUM(CASE WHEN `like` = 'no' THEN 1 ELSE 0 END) AS host_question_unlike_amount",
			)
	)

	if activityID != "" {
		sql = sql.Where("activity_question_user.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return g, errors.New("錯誤: 無法取得提問訊息數量，請重新查詢")
	}

	// 回傳的資料為[]byte格式，需轉換為int64
	for key, value := range item {
		item[key] = utils.GetInt64(value, 0)
	}

	g = g.MapToModel(item)

	return g, nil
}

// FindGuestQuestionAmount 查詢不同類型提問訊息數量(玩家端)
func (g QuestionUserModel) FindGuestQuestionAmount(activityID, userID string) (QuestionUserModel, error) {
	// log.Println("userID: ", userID)
	var (
		sql = g.Table(g.Base.TableName).
			Select(
				"SUM(CASE WHEN message_status = 'yes' THEN 1 ELSE 0 END) AS yes_question_amount",
				"SUM(CASE WHEN message_status = 'no' THEN 1 ELSE 0 END) AS no_question_amount",
				"SUM(CASE WHEN message_status = 'review' THEN 1 ELSE 0 END) AS review_question_amount",

				fmt.Sprintf("SUM(CASE WHEN user_id = '%s' THEN 1 ELSE 0 END) AS user_question_amount", userID),

				"SUM(likes) AS question_likes_amount",

				"SUM(CASE WHEN `like` = 'yes' THEN 1 ELSE 0 END) AS host_question_like_amount",
				"SUM(CASE WHEN `like` = 'no' THEN 1 ELSE 0 END) AS host_question_unlike_amount",
			)
	)

	if activityID != "" {
		sql = sql.Where("activity_question_user.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return g, errors.New("錯誤: 無法取得提問訊息數量，請重新查詢")
	}

	// 回傳的資料為[]byte格式，需轉換為int64
	for key, value := range item {
		item[key] = utils.GetInt64(value, 0)
	}

	g = g.MapToModel(item)

	return g, nil
}

// FindQuestionAmount 查詢提問訊息數量
func (g QuestionUserModel) FindQuestionAmount(activityID, userID, status, like string) (int64, error) {
	var (
		sql = g.Table(g.Base.TableName).
			Select("count(*)")
	)

	if activityID != "" {
		sql = sql.Where("activity_question_user.activity_id", "=", activityID)
	}
	if userID != "" {
		sql = sql.Where("activity_question_user.user_id", "=", userID)
	}
	if status != "" {
		sql = sql.WhereIn("activity_question_user.message_status", interfaces(strings.Split(status, ",")))
	}
	if like != "" {
		sql = sql.WhereIn("activity_question_user.like", interfaces(strings.Split(like, ",")))
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法取得提問訊息數量，請重新查詢")
	}

	// log.Println("item: ", item)

	return item["count(*)"].(int64), nil
}

// FindQuestionLikesAmount 查詢提問訊息總按讚數量
func (g QuestionUserModel) FindQuestionLikesAmount(activityID string) (int64, error) {
	var (
		sql = g.Table(g.Base.TableName).
			Select("sum(likes)")
	)

	if activityID != "" {
		sql = sql.Where("activity_question_user.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法取得提問訊息總按讚數量，請重新查詢")
	}

	// log.Println("item: ", item)

	return utils.GetInt64(item["sum(likes)"], 0), nil
}

// Find 尋找該活動所有提問資料(時間、讚數排序)，包含通過、未通過、審核中
func (g QuestionUserModel) Find(isRedis bool, activityID, userID string,
	status, like string,
	limit, offset int64) ([]QuestionUserModel, error) {
	// questionLock.Lock()
	// defer questionLock.Unlock()
	var (
		questionsOrderByTime = make([]QuestionUserModel, 0)
		// err                  error
	)

	// 判斷redis裡是否有提問資訊
	if isRedis {
		// var (
		// 	idsOrderBytime = make([]string, 0) // 時間排序
		// 	datas          map[string]string
		// )

		// // 所有提問資訊(HASH)
		// datas, err = g.RedisConn.HashGetAllCache(config.QUESTION_REDIS + activityID)
		// if err != nil {
		// 	return questionsOrderByTime, errors.New("錯誤: 取得提問快取資料(所有資訊)發生問題")
		// }

		// // 取得時間排序的提問id資訊(LIST)
		// if idsOrderBytime, err = g.RedisConn.ListRange(config.QUESTION_ORDER_BY_TIME_REDIS+activityID,
		// 	0, 0); err != nil {
		// 	return questionsOrderByTime, errors.New("錯誤: 從redis中取得時間排序的提問資訊發生問題")
		// }

		// for i := 0; i < len(idsOrderBytime); i++ {
		// 	var (
		// 		questionOrderBytime QuestionUserModel // 時間排序
		// 	)

		// 	// 解碼，取得提問資料
		// 	json.Unmarshal([]byte(datas[idsOrderBytime[i]]), &questionOrderBytime)

		// 	// 判斷是否查詢正確資料
		// 	if questionOrderBytime.ID == 0 {
		// 		// 無法從redis取得提問資料時，刪除舊的redis提問資訊並重新執行Find function
		// 		g.RedisConn.DelCache(config.QUESTION_REDIS + activityID)               // 提問資訊
		// 		g.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID) // 提問資訊(時間)

		// 		return g.Find(true, activityID,"","","", 0, 0)
		// 	}

		// 	questionsOrderByTime = append(questionsOrderByTime, questionOrderBytime)
		// }
	}

	// redis中無提問快取資料，從資料表中查詢(資料排序由新到舊)
	if len(questionsOrderByTime) == 0 {
		sql := g.Table(g.Base.TableName).
			Select("activity_question_user.id", "activity_question_user.activity_id",
				"activity_question_user.user_id",
				"activity_question_user.message", "activity_question_user.message_status",
				"activity_question_user.send_time",
				"activity_question_user.likes", "activity_question_user.like",

				// 用戶資料
				"line_users.name", "line_users.avatar").
			LeftJoin(command.Join{
				FieldA:    "activity_question_user.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_question_user.activity_id", "=", activityID).
			OrderBy("activity_question_user.id", "desc")
		if userID != "" {
			sql = sql.Where("activity_question_user.user_id", "=", userID)
		}
		if status != "" {
			sql = sql.WhereIn("activity_question_user.message_status", interfaces(strings.Split(status, ",")))
		}
		if like != "" {
			sql = sql.WhereIn("activity_question_user.like", interfaces(strings.Split(like, ",")))
		}
		if limit != 0 {
			sql = sql.Limit(limit)
		}
		if offset != 0 {
			sql = sql.Offset(offset)
		}

		items, err := sql.All()
		if err != nil {
			return questionsOrderByTime, errors.New("錯誤: 無法取得用戶提問資料，請重新查詢")
		}

		questionsOrderByTime = MapToQuestionUserModel(items)

		// 將提問資訊加入redis中
		if isRedis && len(questionsOrderByTime) > 0 {
			// var (
			// 	questionsparams            = []interface{}{config.QUESTION_REDIS + activityID}               // redis參數(提問資訊)
			// 	questionsOrderByTimeParams = []interface{}{config.QUESTION_ORDER_BY_TIME_REDIS + activityID} // redis參數(時間排序提問)
			// )
			// for _, question := range questionsOrderByTime {
			// 	// json編碼
			// 	questionJson := utils.JSON(question)

			// 	// 提問資訊(HASH)
			// 	questionsparams = append(questionsparams, question.ID, questionJson)
			// 	// 時間排序(LIST)
			// 	questionsOrderByTimeParams = append(questionsOrderByTimeParams, question.ID)
			// }

			// // 將資料寫入redis中(提問資訊，HASH)
			// if err = g.RedisConn.HashMultiSetCache(questionsparams); err != nil {
			// 	return questionsOrderByTime, errors.New("錯誤: 設置提問快取資料(所有提問資訊)發生問題")
			// }
			// // 將資料寫入redis中(時間排序，LIST)
			// if err = g.RedisConn.ListMultiRPush(questionsOrderByTimeParams); err != nil {
			// 	return questionsOrderByTime, errors.New("錯誤: 設置提問快取資料(時間排序)發生問題")
			// }

			// // 設置過期時間
			// g.RedisConn.SetExpire(config.QUESTION_REDIS+activityID, config.REDIS_EXPIRE)               // 提問資訊
			// g.RedisConn.SetExpire(config.QUESTION_ORDER_BY_TIME_REDIS+activityID, config.REDIS_EXPIRE) // 時間
		}
	}

	// questionLock.Unlock()
	return questionsOrderByTime, nil
}

// FindAll 尋找所有提問資料(包含通過、未通過、審核中)
func (g QuestionUserModel) FindAll(activityID string) ([]QuestionUserModel, error) {
	items, err := g.Table(g.Base.TableName).
		Select("activity_question_user.id", "activity_question_user.activity_id",
			"activity_question_user.user_id",
			"activity_question_user.message", "activity_question_user.message_status",
			"activity_question_user.send_time",
			"activity_question_user.likes", "activity_question_user.like",

			// 用戶資料
			"line_users.name", "line_users.avatar").
		LeftJoin(command.Join{
			FieldA:    "activity_question_user.user_id",
			FieldA1:   "line_users.user_id",
			Table:     "line_users",
			Operation: "="}).
		Where("activity_question_user.activity_id", "=", activityID).
		OrderBy("activity_question_user.id", "asc").All()
	if err != nil {
		return nil, errors.New("錯誤: 無法取得提問區紀錄資料，請重新查詢")
	}
	return MapToQuestionUserModel(items), nil
}

// Add 增加資料(用戶提問時)
func (g QuestionUserModel) Add(isRedis bool, model EditQuestionUserModel) (id int64, err error) {
	var (
		now, _ = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)

	if utf8.RuneCountInString(model.Message) > 100 {
		return id, errors.New("錯誤: 提問訊息上限為100個字元，請輸入有效的提問提問訊息")
	}
	if model.MessageStatus != "yes" && model.MessageStatus != "no" &&
		model.MessageStatus != "review" {
		return id, errors.New("錯誤: 訊息審核狀態資料發生問題，請輸入有效的訊息審核狀態")
	}

	if id, err = g.Table(g.TableName).Insert(command.Value{
		"activity_id":    model.ActivityID,
		"user_id":        model.UserID,
		"message":        model.Message,
		"message_status": model.MessageStatus,
		"likes":          0,
		"like":           "no",
		"send_time":      now,
	}); err != nil {
		return id, err
	}

	// 新增提問資訊至redis中
	// if isRedis {
	// 判斷redis裡是否有用戶資訊，有則不查詢資料表
	// user, err := DefaultLineModel().SetDbConn(g.DbConn).
	// 	SetRedisConn(g.RedisConn).Find(true, "", "user_id", model.UserID)
	// if err != nil {
	// 	return id, errors.New("錯誤: 無法取得用戶資訊")
	// }

	// // 將新的提問加入redis(所有提問資訊)
	// g.RedisConn.HashSetCache(config.QUESTION_REDIS+model.ActivityID,
	// 	strconv.Itoa(int(id)), utils.JSON(QuestionUserModel{
	// 		ID:            id,
	// 		ActivityID:    model.ActivityID,
	// 		UserID:        model.UserID,
	// 		Message:       model.Message,
	// 		MessageStatus: model.MessageStatus,
	// 		Likes:         0,
	// 		SendTime:      now.Format("2006-01-02 15:04:05"),
	// 		Like:          "no",
	// 		Name:          user.Name,
	// 		Avatar:        user.Avatar,
	// 	}))

	// // 將提問資料加入redis中(時間排序)
	// g.RedisConn.ListLPush(config.QUESTION_ORDER_BY_TIME_REDIS+model.ActivityID,
	// 	strconv.Itoa(int(id)))
	// }
	return id, err
}

// Update 更新提問訊息狀態資料
func (g QuestionUserModel) Update(isRedis bool, model EditQuestionUserModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{"message_status"}
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

	if err := g.Table(g.Base.TableName).
		Where("id", "=", model.ID).Update(fieldValues); err != nil {
		return err
	}

	// if isRedis {
	// 清除redis資料，重新寫入redis中
	// g.RedisConn.DelCache(config.QUESTION_REDIS + model.ActivityID)               // 提問資訊
	// g.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + model.ActivityID) // 提問資訊(時間)
	// g.Find(true, model.ActivityID)

	// var (
	// 	data     string
	// 	question QuestionUserModel
	// )
	// // 取得提問資料
	// data, err := g.RedisConn.HashGetCache(config.QUESTION_REDIS+model.ActivityID,
	// 	model.ID)
	// if err == nil {
	// 	// 解碼
	// 	json.Unmarshal([]byte(data), &question)
	// }

	// // 無法從redis取得提問資料時，刪除舊的redis提問資訊並重新執行Find function
	// if question.ID == 0 {
	// 	g.RedisConn.DelCache(config.QUESTION_REDIS + model.ActivityID)               // 提問資訊
	// 	g.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + model.ActivityID) // 提問資訊(時間)
	// 	g.Find(true, model.ActivityID,"","","", 0, 0)
	// 	return nil
	// }

	// // redis裡已有提問資料，將更新的提問加入redis
	// g.RedisConn.HashSetCache(config.QUESTION_REDIS+model.ActivityID,
	// 	model.ID, utils.JSON(QuestionUserModel{
	// 		ID:            question.ID,
	// 		ActivityID:    question.ActivityID,
	// 		UserID:        question.UserID,
	// 		Message:       question.Message,
	// 		MessageStatus: model.MessageStatus,
	// 		Likes:         question.Likes,
	// 		SendTime:      question.SendTime,
	// 		Like:          question.Like,
	// 		Name:          question.Name,
	// 		Avatar:        question.Avatar,
	// 	}))
	// }
	return nil
}

// UpdateHostLike 更新提問資料、遞增遞減按讚數(主持端)
func (g QuestionUserModel) UpdateHostLikes(isRedis bool, id int64,
	activityID string, like string) (err error) {
	var (
		value       int64
		fieldValues = command.Value{"like": like}
	)
	if like != "yes" && like != "no" {
		return errors.New("錯誤: 是否按讚資料發生問題，請輸入有效的資料")
	}

	if like == "yes" {
		fieldValues["likes"] = "likes + 1" // 遞增
		value++
	} else if like == "no" {
		fieldValues["likes"] = "likes - 1" // 遞減
		value--
	}

	// 更新資料表
	if err = g.Table(g.Base.TableName).Where("activity_id", "=", activityID).
		Where("id", "=", id).Update(fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新提問資料發生問題")
	}

	// 判斷redis裡是否有提問資訊
	// if isRedis {
	// var (
	// 	data     string
	// 	question QuestionUserModel
	// )
	// // 所有提問資訊
	// data, err = g.RedisConn.HashGetCache(config.QUESTION_REDIS+activityID,
	// 	strconv.Itoa(int(id)))
	// if err == nil {
	// 	// 有取得資料，解碼
	// 	json.Unmarshal([]byte(data), &question)
	// }

	// // 無法從redis取得提問資料時，刪除舊的redis提問資訊並重新執行Find function
	// if question.ID == 0 {
	// 	g.RedisConn.DelCache(config.QUESTION_REDIS + activityID)               // 提問資訊
	// 	g.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID) // 提問資訊(時間)
	// 	g.Find(true, activityID,"","","", 0, 0)
	// 	return
	// }

	// // redis裡已有提問資料，將更新的提問加入redis(時間排序)
	// g.RedisConn.HashSetCache(config.QUESTION_REDIS+activityID,
	// 	strconv.Itoa(int(id)), utils.JSON(QuestionUserModel{
	// 		ID:            question.ID,
	// 		ActivityID:    question.ActivityID,
	// 		UserID:        question.UserID,
	// 		Message:       question.Message,
	// 		MessageStatus: question.MessageStatus,
	// 		Likes:         question.Likes + value,
	// 		SendTime:      question.SendTime,
	// 		Like:          like,
	// 		Name:          question.Name,
	// 		Avatar:        question.Avatar,
	// 	}))
	// }
	return
}

// UpdateGuestLikes 遞增遞減按讚數資料(用戶端)
func (g QuestionUserModel) UpdateGuestLikes(isRedis bool, id int64,
	activityID string, like string) (err error) {
	var (
		value       int64
		fieldValues = command.Value{}
	)
	if like != "like" && like != "unlike" {
		return errors.New("錯誤: 按讚資料發生問題，請輸入有效的資料")
	}

	if like == "like" {
		fieldValues["likes"] = "likes + 1" // 遞增
		value++
	} else if like == "unlike" {
		fieldValues["likes"] = "likes - 1" // 遞減
		value--
	}

	// 資料表更新
	if err = g.Table(g.Base.TableName).Where("activity_id", "=", activityID).
		Where("id", "=", id).Update(fieldValues); err != nil &&
		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 更新提問資料發生問題")
	}

	// 判斷redis裡是否有提問資訊
	// if isRedis {
	// var (
	// 	data     string
	// 	question QuestionUserModel
	// )
	// // 所有提問資訊
	// data, err = g.RedisConn.HashGetCache(config.QUESTION_REDIS+activityID,
	// 	strconv.Itoa(int(id)))
	// if err == nil {
	// 	// 有取得資料，解碼
	// 	json.Unmarshal([]byte(data), &question)
	// }

	// // 無法從redis取得提問資料時，刪除舊的redis提問資訊並重新執行Find function
	// if question.ID == 0 {
	// 	g.RedisConn.DelCache(config.QUESTION_REDIS + activityID)               // 提問資訊
	// 	g.RedisConn.DelCache(config.QUESTION_ORDER_BY_TIME_REDIS + activityID) // 提問資訊(時間)
	// 	g.Find(true, activityID,"","","", 0, 0)
	// 	return
	// }

	// g.RedisConn.HashSetCache(config.QUESTION_REDIS+activityID,
	// 	strconv.Itoa(int(id)), utils.JSON(QuestionUserModel{
	// 		ID:            question.ID,
	// 		ActivityID:    question.ActivityID,
	// 		UserID:        question.UserID,
	// 		Message:       question.Message,
	// 		MessageStatus: question.MessageStatus,
	// 		Likes:         question.Likes + value,
	// 		Like:          question.Like,
	// 		SendTime:      question.SendTime,
	// 		Name:          question.Name,
	// 		Avatar:        question.Avatar,
	// 	}))
	// }
	return nil
}

// MapToModel map轉換QuestionUserModel
func (g QuestionUserModel) MapToModel(item map[string]interface{}) QuestionUserModel {
	// log.Println("查詢不同類型訊息數量")

	// json解碼，轉換成strcut
	b, _ := json.Marshal(item)
	json.Unmarshal(b, &g)

	// log.Println("第一種過濾: ", g.YesQuestionAmount, g.NoQuestionAmount, g.ReviewQuestionAmount)
	// log.Println("第二種過濾: ", g.QuestionLikesAmount)
	// log.Println("第三種過濾: ", g.HostQuestionLikeAmount, g.HostQuestionUnlikeAmount)
	// log.Println("第四種過濾: ", g.UserQuestionAmount)

	return g
}

// MapToQuestionUserModel map轉換[]QuestionUserModel
func MapToQuestionUserModel(items []map[string]interface{}) []QuestionUserModel {
	var questions = make([]QuestionUserModel, 0)
	for _, item := range items {
		var (
			question QuestionUserModel
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &question)

		questions = append(questions, question)
	}
	return questions
}

// question.ID, _ = item["id"].(int64)
// question.UserID, _ = item["user_id"].(string)
// question.ActivityID, _ = item["activity_id"].(string)
// question.Message, _ = item["message"].(string)
// question.MessageStatus, _ = item["message_status"].(string)
// question.Likes, _ = item["likes"].(int64)
// question.SendTime, _ = item["send_time"].(string)
// question.Like, _ = item["like"].(string)

// // 用戶資訊
// question.Name, _ = item["name"].(string)
// question.Avatar, _ = item["avatar"].(string)

// values      = []string{model.MessageStatus}
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
