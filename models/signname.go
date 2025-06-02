package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/command"
	"hilive/modules/mongo"
	"hilive/modules/utils"
)

// SignnameModel 資料表欄位
type SignnameModel struct {
	Base           `json:"-"`
	ID             int64  `json:"id"`
	ActivityUserID string `json:"activity_user_id"`
	UserID         string `json:"user_id"`
	ActivityID     string `json:"activity_id"`
	Picture        string `json:"picture"`

	Name   string `json:"name"`   // 玩家名稱
	Avatar string `json:"avatar"` // 玩家頭像
}

// NewSignnameModel 資料表欄位
type NewSignnameModel struct {
	// ID             string `json:"id"`
	ActivityUserID string `json:"activity_user_id"`
	UserID         string `json:"user_id"`
	ActivityID     string `json:"activity_id"`
	Picture        string `json:"picture"`
}

// SignnameSettingModel 資料表欄位
type SignnameSettingModel struct {
	Base               `json:"-"`
	ID                 int64  `json:"id"`
	ActivityID         string `json:"activity_id"`
	UserID             string `json:"user_id"`
	SignnameMode       string `json:"signname_mode"`
	SignnameTimes      int64  `json:"signname_times"`
	SignnameDisplay    string `json:"signname_display"`
	SignnameLimitTimes string `json:"signname_limit_times"`
	SignnameTopic      string `json:"signname_topic"`
	SignnameMusic      string `json:"signname_music"`
	SignnameLoop       string `json:"signname_loop"`
	SignnameLatest     string `json:"signname_latest"`
	SignnameContent    string `json:"signname_content"`
	SignnameShowName   string `json:"signname_show_name"`

	// SignnameClassicHPic01 string `json:"signname_classic_h_pic_01" example:"picture"`
	// SignnameClassicHPic02 string `json:"signname_classic_h_pic_02" example:"picture"`
	// SignnameClassicHPic03 string `json:"signname_classic_h_pic_03" example:"picture"`
	// SignnameClassicHPic04 string `json:"signname_classic_h_pic_04" example:"picture"`
	// SignnameClassicHPic05 string `json:"signname_classic_h_pic_05" example:"picture"`
	// SignnameClassicCPic01 string `json:"signname_classic_c_pic_01" example:"picture"`

	// // 音樂
	// SignnameBgm string `json:"signname_bgm" example:"picture"`
}

// EditSignnameSettingModel 資料表欄位
type EditSignnameSettingModel struct {
	ActivityID         string `json:"activity_id"`
	SignnameMode       string `json:"signname_mode"`
	SignnameTimes      string `json:"signname_times"`
	SignnameDisplay    string `json:"signname_display"`
	SignnameLimitTimes string `json:"signname_limit_times"`
	SignnameTopic      string `json:"signname_topic"`
	SignnameMusic      string `json:"signname_music"`
	SignnameLoop       string `json:"signname_loop"`
	SignnameLatest     string `json:"signname_latest"`
	SignnameContent    string `json:"signname_content"`
	SignnameShowName   string `json:"signname_show_name"`

	// 簽名牆自定義圖片
	SignnameClassicHPic01 string `json:"signname_classic_h_pic_01" example:"picture"`
	SignnameClassicHPic02 string `json:"signname_classic_h_pic_02" example:"picture"`
	SignnameClassicHPic03 string `json:"signname_classic_h_pic_03" example:"picture"`
	SignnameClassicHPic04 string `json:"signname_classic_h_pic_04" example:"picture"`
	SignnameClassicHPic05 string `json:"signname_classic_h_pic_05" example:"picture"`
	SignnameClassicHPic06 string `json:"signname_classic_h_pic_06" example:"picture"`
	SignnameClassicHPic07 string `json:"signname_classic_h_pic_07" example:"picture"`
	SignnameClassicGPic01 string `json:"signname_classic_g_pic_01" example:"picture"`
	SignnameClassicCPic01 string `json:"signname_classic_c_pic_01" example:"picture"`

	// 音樂
	SignnameBgm string `json:"signname_bgm" example:"picture"`
}

// DefaultSignnameSettingModel 預設SignnameSettingModel
func DefaultSignnameSettingModel() SignnameSettingModel {
	return SignnameSettingModel{Base: Base{TableName: config.ACTIVITY_2_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (s SignnameSettingModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) SignnameSettingModel {
	s.DbConn = dbconn
	s.RedisConn = cacheconn
	s.MongoConn = mongoconn
	return s
}

// SetDbConn 設定connection
// func (s SignnameSettingModel) SetDbConn(conn db.Connection) SignnameSettingModel {
// 	s.DbConn = conn
// 	return s
// }

// // SetRedisConn 設定connection
// func (s SignnameSettingModel) SetRedisConn(conn cache.Connection) SignnameSettingModel {
// 	s.RedisConn = conn
// 	return s
// }

// DefaultSignnameModel 預設SignnameModel
func DefaultSignnameModel() SignnameModel {
	return SignnameModel{Base: Base{TableName: config.ACTIVITY_SIGNNAME_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (s SignnameModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) SignnameModel {
	s.DbConn = dbconn
	s.RedisConn = cacheconn
	s.MongoConn = mongoconn
	return s
}

// SetDbConn 設定connection
// func (s SignnameModel) SetDbConn(conn db.Connection) SignnameModel {
// 	s.DbConn = conn
// 	return s
// }

// // SetRedisConn 設定connection
// func (s SignnameModel) SetRedisConn(conn cache.Connection) SignnameModel {
// 	s.RedisConn = conn
// 	return s
// }

// Find 查詢簽名牆設置資料
func (s SignnameSettingModel) Find(isRedis bool, activityID string) (SignnameSettingModel, error) {
	if isRedis {
		// 判斷redis裡是否有簽名牆設置資訊，有則不執行查詢資料表功能
		dataMap, err := s.RedisConn.HashGetAllCache(config.SIGNNAME_REDIS + activityID)
		if err != nil {
			return SignnameSettingModel{}, errors.New("錯誤: 取得簽名牆快取資料發生問題")
		}

		s.ID = utils.GetInt64FromStringMap(dataMap, "id", 0)
		s.UserID, _ = dataMap["user_id"]
		s.ActivityID, _ = dataMap["activity_id"]
		s.SignnameMode, _ = dataMap["signname_mode"]
		s.SignnameTimes = utils.GetInt64FromStringMap(dataMap, "signname_times", 0)
		s.SignnameDisplay, _ = dataMap["signname_display"]
		s.SignnameLimitTimes, _ = dataMap["signname_limit_times"]
		s.SignnameTopic, _ = dataMap["signname_topic"]
		s.SignnameMusic, _ = dataMap["signname_music"]
		s.SignnameLoop, _ = dataMap["signname_loop"]
		s.SignnameLatest, _ = dataMap["signname_latest"]
		s.SignnameContent, _ = dataMap["signname_content"]
		s.SignnameShowName, _ = dataMap["signname_show_name"]
	}

	// redis中無簽名牆設置快取資料，從資料表中查詢
	if s.ID == 0 {
		item, err := s.Table(s.Base.TableName).Select(
			// 簽名牆
			"activity_2.id",
			"activity_2.activity_id",
			"activity_2.signname_mode",
			"activity_2.signname_times",
			"activity_2.signname_display",
			"activity_2.signname_limit_times",
			"activity_2.signname_topic",
			"activity_2.signname_music",
			"activity_2.signname_loop",
			"activity_2.signname_latest",
			"activity_2.signname_content",
			"activity_2.signname_show_name",

			// 活動資料表
			"activity.user_id",

			// "activity_2.signname_classic_h_pic_01",
			// "activity_2.signname_classic_h_pic_02",
			// "activity_2.signname_classic_h_pic_03",
			// "activity_2.signname_classic_h_pic_04",
			// "activity_2.signname_classic_h_pic_05",
			// "activity_2.signname_classic_c_pic_01",

			// "activity_2.signname_bgm",
		).
			LeftJoin(command.Join{
				FieldA:    "activity.activity_id",
				FieldA1:   "activity_2.activity_id",
				Table:     "activity",
				Operation: "="}).
			Where("activity_2.activity_id", "=", activityID).
			First()
		if err != nil || item == nil {
			return SignnameSettingModel{}, errors.New("錯誤: 無法取得簽名牆資訊，請重新查詢")
		}

		s = s.MapToModel(item)

		if isRedis {
			values := []interface{}{config.SIGNNAME_REDIS + activityID}
			values = append(values, "id", s.ID)
			values = append(values, "user_id", s.UserID)
			values = append(values, "activity_id", s.ActivityID)
			values = append(values, "signname_mode", s.SignnameMode)
			values = append(values, "signname_times", s.SignnameTimes)
			values = append(values, "signname_display", s.SignnameDisplay)
			values = append(values, "signname_limit_times", s.SignnameLimitTimes)
			values = append(values, "signname_topic", s.SignnameTopic)
			values = append(values, "signname_music", s.SignnameMusic)
			values = append(values, "signname_loop", s.SignnameLoop)
			values = append(values, "signname_latest", s.SignnameLatest)
			values = append(values, "signname_content", s.SignnameContent)
			values = append(values, "signname_show_name", s.SignnameShowName)

			if err := s.RedisConn.HashMultiSetCache(values); err != nil {
				return s, errors.New("錯誤: 設置簽名牆快取資料發生問題")
			}

			// 設置過期時間
			// s.RedisConn.SetExpire(config.SIGNNAME_REDIS+activityID,
			// 	config.REDIS_EXPIRE)

		}
	}
	return s, nil
}

// FindSignnameAmount 查詢該活動簽名牆數量
func (s SignnameModel) FindSignnameAmount(activityID string) (int64, error) {
	var (
		sql = s.Table(s.Base.TableName).
			Select("count(*)")
	)

	if activityID != "" {
		sql = sql.Where("activity_signname.activity_id", "=", activityID)
	}

	item, err := sql.First()
	if err != nil {
		return 0, errors.New("錯誤: 無法該活動報名簽到人員數量，請重新查詢")
	}

	// log.Println("item: ", item)

	return item["count(*)"].(int64), nil
}

// Find 尋找資料
func (s SignnameModel) Find(isRedis bool, activityID string,
	limit, offset int64) ([]SignnameModel, error) {
	var (
		signnames = make([]SignnameModel, 0)
		// err       error
	)

	if isRedis {
		// var (
		// 	signnamesOrderByTime = make([]string, 0)
		// 	datas                map[string]string
		// )

		// // 查詢redis裡是否有簽名牆資訊(HASH)，有則不查詢資料表
		// datas, err = s.RedisConn.HashGetAllCache(config.SIGNNAME_DATAS_REDIS + activityID)
		// if err != nil {
		// 	return signnames, errors.New("錯誤: 從redis中取得簽名牆資訊(HASH)發生問題")
		// }

		// // 查詢redis裡是否有簽名牆資訊(LIST)，有則不查詢資料表
		// signnamesOrderByTime, err = s.RedisConn.ListRange(config.SIGNNAME_ORDER_BY_TIME_REDIS+activityID, 0, 0)
		// if err != nil {
		// 	return signnames, errors.New("錯誤: 從redis中取得簽名牆資訊(LIST)發生問題")
		// }
		// // fmt.Println("signnamesOrderByTime: ", signnamesOrderByTime)

		// // 處理簽到人員資訊
		// for i := range signnamesOrderByTime {
		// 	var (
		// 		signname SignnameModel
		// 	)
		// 	// 解碼
		// 	json.Unmarshal([]byte(datas[signnamesOrderByTime[i]]), &signname)

		// 	signnames = append(signnames, signname)
		// }
	}

	if len(signnames) == 0 {
		sql := s.Table(s.Base.TableName).
			Select("activity_signname.id", "activity_signname.activity_id",
				"activity_signname.user_id",
				"activity_signname.picture",

				// 活動資料表
				"activity.user_id as activity_user_id",

				// 用戶
				"line_users.name", "line_users.avatar",
			).
			LeftJoin(command.Join{
				FieldA:    "activity_signname.activity_id",
				FieldA1:   "activity.activity_id",
				Table:     "activity",
				Operation: "="}).
			LeftJoin(command.Join{
				FieldA:    "activity_signname.user_id",
				FieldA1:   "line_users.user_id",
				Table:     "line_users",
				Operation: "="}).
			Where("activity_signname.activity_id", "=", activityID).
			OrderBy("id", "asc")

			// 判斷參數是否為空
		if limit != 0 {
			sql = sql.Limit(limit)
		}
		if offset != 0 {
			sql = sql.Offset(offset)
		}

		items, err := sql.All()
		if err != nil {
			return signnames, errors.New("錯誤: 無法取得簽名牆資訊，請重新查詢")
		}
		signnames = s.MapToModel(items)

		if isRedis && len(signnames) > 0 {
			// var (
			// 	signnamesparams            = []interface{}{config.SIGNNAME_DATAS_REDIS + activityID}         // redis參數(簽名牆資訊，HASH)
			// 	signnamesOrderByTimeParams = []interface{}{config.SIGNNAME_ORDER_BY_TIME_REDIS + activityID} // redis參數(簽名牆資訊，LIST)
			// )

			// for _, signname := range signnames {

			// 	// json編碼
			// 	signnameJson := utils.JSON(signname)

			// 	// 簽名牆資訊
			// 	signnamesparams = append(signnamesparams, signname.ID, signnameJson)
			// 	// 時間排序
			// 	signnamesOrderByTimeParams = append(signnamesOrderByTimeParams, signname.ID)
			// }

			// // 將資料加入redis中(簽名牆資訊，HASH)
			// if err = s.RedisConn.HashMultiSetCache(signnamesparams); err != nil {
			// 	return signnames, errors.New("錯誤: 設置提問快取資料(簽名牆資訊，HASH)發生問題")
			// }

			// // 將資料加入redis中(簽名牆資訊，LIST)
			// if err = s.RedisConn.ListMultiRPush(signnamesOrderByTimeParams); err != nil {
			// 	return signnames, errors.New("錯誤: 設置提問快取資料(簽名牆資訊，LIST)發生問題")
			// }

			// // 設置過期時間
			// s.RedisConn.SetExpire(config.SIGNNAME_DATAS_REDIS+activityID, config.REDIS_EXPIRE)         // 簽名牆資訊，HASH
			// s.RedisConn.SetExpire(config.SIGNNAME_ORDER_BY_TIME_REDIS+activityID, config.REDIS_EXPIRE) // 簽名牆資訊，LIST
		}
	}
	return signnames, nil
}

// Add 增加資料
func (s SignnameModel) Add(isRedis bool, model NewSignnameModel) error {
	var (
		filename = fmt.Sprintf("%s.png", utils.UUID(8))
		// pictureURL = fmt.Sprintf("/admin/uploads/%s/%s/interact/sign/signname/%s",
		// 	model.ActivityUserID, model.ActivityID, filename)
	)

	// 圖片資料處理(base64格式，data:image/png;base64,...)
	// 去掉前缀
	base64Image := strings.TrimPrefix(model.Picture, "data:image/png;base64,")

	// 解碼bese64格式資料
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return errors.New("錯誤: 解碼圖片資料發生問題")
	}

	// 在cloud storage建立檔案
	f, err := os.Create(fmt.Sprintf("./hilives/hilive/uploads/%s/%s/interact/sign/signname/%s",
		model.ActivityUserID, model.ActivityID, filename))
	if err != nil {
		return errors.New("錯誤: 建立圖片檔案發生問題")
	}
	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = err2
		}
	}()

	// 將圖片資料寫入檔案
	_, err = f.Write(imageData)
	if err != nil {
		return errors.New("錯誤: 圖片資料寫入檔案發生問題")
	}

	_, err = s.Table(s.TableName).Insert(command.Value{
		"activity_id": model.ActivityID,
		"user_id":     model.UserID,
		"picture":     filename,
	})

	if isRedis {
		// 將新的簽名牆資料加入redis(HASH)
		// s.RedisConn.HashSetCache(config.SIGNNAME_DATAS_REDIS+model.ActivityID,
		// 	strconv.Itoa(int(id)), utils.JSON(SignnameModel{
		// 		ID:         id,
		// 		ActivityID: model.ActivityID,
		// 		UserID:     model.UserID,
		// 		Picture:    pictureURL,
		// 	}))

		// // 將新的簽名牆資料加入redis(LIST)
		// s.RedisConn.ListRPush(config.SIGNNAME_ORDER_BY_TIME_REDIS+model.ActivityID,
		// 	strconv.Itoa(int(id)))
	}
	return err
}

// Delete 刪除資料
func (s SignnameModel) Delete(isRedis bool, activityID string, ids []string) error {

	err := s.Table(s.TableName).
		Where("activity_id", "=", activityID).
		WhereIn("id", interfaces(ids)).
		Delete()

	if err != nil && err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
		return err
	}

	if isRedis {
	}
	return nil
}

// UpdateSignname 更新簽名牆基本設置資料
func (a ActivityModel) UpdateSignnameSetting(isRedis bool, model EditSignnameSettingModel) error {
	var (
		fieldValues = command.Value{"signname_edit_times": "signname_edit_times + 1"}
		fields      = []string{
			"signname_mode", "signname_times", "signname_display",
			"signname_limit_times",
			"signname_topic", "signname_music",
			"signname_loop", "signname_latest",
			"signname_content",
			"signname_show_name",

			"signname_classic_h_pic_01",
			"signname_classic_h_pic_02",
			"signname_classic_h_pic_03",
			"signname_classic_h_pic_04",
			"signname_classic_h_pic_05",
			"signname_classic_h_pic_06",
			"signname_classic_h_pic_07",
			"signname_classic_g_pic_01",
			"signname_classic_c_pic_01",

			"signname_bgm",
		}
	)

	if model.SignnameMode != "" {
		if model.SignnameMode != "terminal" && model.SignnameMode != "phone" {
			return errors.New("錯誤: 簽名模式資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameTimes != "" {
		if _, err := strconv.Atoi(model.SignnameTimes); err != nil {
			return errors.New("錯誤: 簽名次數資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameDisplay != "" {
		if model.SignnameDisplay != "dynamics" && model.SignnameDisplay != "static" {
			return errors.New("錯誤: 簽名牆顯示方式資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameLimitTimes != "" {
		if model.SignnameLimitTimes != "open" && model.SignnameLimitTimes != "close" {
			return errors.New("錯誤: 限制簽名次數資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameTopic != "" {
		if model.SignnameTopic != "01_classic" {
			return errors.New("錯誤: 主題資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameMusic != "" && (model.SignnameMusic != "classic" && model.SignnameMusic != "customize") {
		return errors.New("錯誤: 音樂資料發生問題，請輸入有效的音樂")
	}

	if model.SignnameLoop != "" {
		if model.SignnameLoop != "open" && model.SignnameLoop != "close" {
			return errors.New("錯誤: 是否輪播資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameLatest != "" {
		if model.SignnameLatest != "open" && model.SignnameLatest != "close" {
			return errors.New("錯誤: 跳轉至最新簽名資料發生問題，請輸入有效的資料")
		}
	}

	if model.SignnameContent != "" {
		if model.SignnameContent != "write" && model.SignnameContent != "message" {
			return errors.New("錯誤: 呈現內容資料發生問題，請輸入有效的資料")
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
	} else if len(fieldValues) != 0 {

	}

	// 清除簽名牆設置資料(redis)
	if isRedis {
		a.RedisConn.DelCache(config.SIGNNAME_REDIS + model.ActivityID)
		a.RedisConn.DelCache(config.ACTIVITY_REDIS + model.ActivityID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_SIGNNAME_EDIT_TIMES_REDIS+model.ActivityID, "修改資料")
	}

	return a.Table(config.ACTIVITY_2_TABLE).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// MapToModel map轉換[]SignnameModel
func (s SignnameModel) MapToModel(items []map[string]interface{}) []SignnameModel {

	var signnames = make([]SignnameModel, 0)

	for _, item := range items {
		var (
			signname SignnameModel
			// activityUserID, _ = item["activity_user_id"].(string) // 活動主辦人user_id
		)

		// json解碼，轉換成strcut
		b, _ := json.Marshal(item)
		json.Unmarshal(b, &signname)

		// 圖片資料處理
		signname.Picture = fmt.Sprintf("/admin/uploads/%s/%s/interact/sign/signname/%s",
			signname.ActivityUserID, signname.ActivityID, signname.Picture) // 圖片資料完整路徑

		signnames = append(signnames, signname)

	}

	return signnames
}

// MapToModel map轉換model
func (s SignnameSettingModel) MapToModel(m map[string]interface{}) SignnameSettingModel {

	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &s)

	return s
}

// log.Println("簽名牆MapToModel換個方式處理")
// log.Println("簽名資料處理前: ", s)

// log.Println("簽名資料處理後: ", s)

// log.Println("簽名牆MapToModel換個方式處理")

// log.Println("處理前: ", s)
// log.Println("處理後: ", s)

// 活動
// s.ID, _ = m["id"].(int64)
// s.UserID, _ = m["user_id"].(string)
// s.ActivityID, _ = m["activity_id"].(string)
// s.SignnameMode, _ = m["signname_mode"].(string)
// s.SignnameTimes, _ = m["signname_times"].(int64)
// s.SignnameDisplay, _ = m["signname_display"].(string)
// s.SignnameLimitTimes, _ = m["signname_limit_times"].(string)
// s.SignnameTopic, _ = m["signname_topic"].(string)
// s.SignnameMusic, _ = m["signname_music"].(string)
// s.SignnameLoop, _ = m["signname_loop"].(string)
// s.SignnameLatest, _ = m["signname_latest"].(string)
// s.SignnameContent, _ = m["signname_content"].(string)

// s.SignnameClassicHPic01, _ = m["signname_classic_h_pic_01"].(string)
// s.SignnameClassicHPic02, _ = m["signname_classic_h_pic_02"].(string)
// s.SignnameClassicHPic03, _ = m["signname_classic_h_pic_03"].(string)
// s.SignnameClassicHPic04, _ = m["signname_classic_h_pic_04"].(string)
// s.SignnameClassicHPic05, _ = m["signname_classic_h_pic_05"].(string)
// s.SignnameClassicCPic01, _ = m["signname_classic_c_pic_01"].(string)

// s.SignnameBgm, _ = m["signname_bgm"].(string)

// var (
// 	signname          SignnameModel
// 	activityUserID, _ = item["activity_user_id"].(string) // 活動主辦人user_id
// 	picture, _        = item["picture"].(string)          // 圖片資料(檔名)
// )

// signname.ID, _ = item["id"].(int64)
// signname.UserID, _ = item["user_id"].(string)
// signname.ActivityID, _ = item["activity_id"].(string)
// signname.Picture = fmt.Sprintf("/admin/uploads/%s/%s/interact/sign/signname/%s",
// 	activityUserID, signname.ActivityID, picture) // 圖片資料完整路徑

// UpdateSignname 更新活動介紹基本設置資料
// func (a ActivityModel) UpdateSignname(model EditSignnameSettingModel) error {
// 	if model.SignnameTitle != "" {
// 		if utf8.RuneCountInString(model.SignnameTitle) > 20 {
// 			return errors.New("錯誤: 手機頁面標題上限為20個字元，請輸入有效的頁面標題")
// 		}
// 		fieldValues := command.Value{
// 			"Signname_title": model.SignnameTitle,
// 		}
// 		return a.Table(a.Base.TableName).
// 			Where("activity_id", "=", model.ActivityID).Update(fieldValues)
// 	}
// 	re

// values = []string{
// 	model.SignnameMode, model.SignnameTimes, model.SignnameDisplay,
// 	model.SignnameLimitTimes,
// 	model.SignnameTopic, model.SignnameMusic,
// 	model.SignnameLoop, model.SignnameLatest,
// 	model.SignnameContent,
// 	model.SignnameShowName,

// 	model.SignnameClassicHPic01,
// 	model.SignnameClassicHPic02,
// 	model.SignnameClassicHPic03,
// 	model.SignnameClassicHPic04,
// 	model.SignnameClassicHPic05,
// 	model.SignnameClassicHPic06,
// 	model.SignnameClassicHPic07,
// 	model.SignnameClassicGPic01,
// 	model.SignnameClassicCPic01,

// 	model.SignnameBgm,
// }

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
