package models

import (
	"encoding/json"
	"errors"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// UserCustomizeSceneModel 用戶自定義場景資料
type UserCustomizeSceneModel struct {
	Base               `json:"-" bson:"-"`
	ID                 int64                  `json:"id" bson:"id"`                                     // id
	TopicID            string                 `json:"topic_id" bson:"topic_id"`                         // 主題id
	UserID             string                 `json:"user_id" bson:"user_id" `                          // 用戶ID
	TopicName          string                 `json:"topic_name" bson:"topic_name"`                     // 主題名稱
	CustomizeSceneData map[string]interface{} `json:"customize_scene_data" bson:"customize_scene_data"` // 這裡會包含所有畫面設定資料
}

// // CustomizeSceneData 自定義場景Data欄位資料
// type CustomizeSceneData struct {
// 	TopicName  string         `json:"topic_name" bson:"topic_name"` // 主題名稱
// 	Pictures   []MediaItem    `json:"pictures" bson:"pictures"`     // 圖片素材
// 	Animations []MediaItem    `json:"animations" bson:"animations"` // 動畫素材
// 	Videos     []MediaItem    `json:"videos" bson:"videos"`         // 影片素材
// 	MyMedia    []MyMediaItem  `json:"my_media" bson:"my_media"`     // 使用者自定義素材
// 	Scenes     [][]SceneItem  `json:"scenes" bson:"scenes"`         // 多場景物件
// 	Template   []TemplateItem `json:"template" bson:"template"`     // 模板
// 	Prefabs    []SceneItem    `json:"prefabs" bson:"prefabs"`
// }

// type MediaItem struct {
// 	Name string `json:"name" bson:"name"`
// 	URL  string `json:"url" bson:"url"`
// }

// type MyMediaItem struct {
// 	Name         string `json:"name" bson:"name"`
// 	URL          string `json:"url" bson:"url"`
// 	Type         string `json:"type" bson:"type"`
// 	URLAnimation string `json:"url_animation" bson:"url_animation"`
// }

// type SceneItem struct {
// 	ObjectName        string   `json:"object_name" bson:"object_name"`
// 	ObjectType        string   `json:"object_type" bson:"object_type"`
// 	ObjectMediaIndex  int64    `json:"object_media_index" bson:"object_media_index"`
// 	ObjectMediaURL    string   `json:"object_media_url" bson:"object_media_url"`
// 	AnimationURL      string   `json:"animation_url" bson:"animation_url"`
// 	AnimationPlay     string   `json:"animation_play" bson:"animation_play"`
// 	ObjectPositionX   int64    `json:"object_positionX" bson:"object_positionX"`
// 	ObjectPositionY   int64    `json:"object_positionY" bson:"object_positionY"`
// 	ObjectRotation    int64    `json:"object_rotation" bson:"object_rotation"`
// 	ObjectScaleX      int64    `json:"object_scaleX" bson:"object_scaleX"`
// 	ObjectScaleY      int64    `json:"object_scaleY" bson:"object_scaleY"`
// 	ObjectNecessary   bool     `json:"object_necessary" bson:"object_necessary"`
// 	TextString        string   `json:"text_string" bson:"text_string"`
// 	TextColor         string   `json:"text_color" bson:"text_color"`
// 	TextSize          int64    `json:"text_size" bson:"text_size"`
// 	ButtonImgsURL     []string `json:"button_imgs_url" bson:"button_imgs_url"`
// 	ScrollTypesetting string   `json:"scroll_typesetting" bson:"scroll_typesetting"`
// 	ScrollSpacingX    int64    `json:"scroll_spacingX" bson:"scroll_spacingX"`
// 	ScrollSpacingY    int64    `json:"scroll_spacingY" bson:"scroll_spacingY"`
// 	IsActive          bool     `json:"is_active" bson:"is_active"`
// 	PrefabIndex       []int64  `json:"prefabIndex" bson:"prefabIndex"`
// }

// type TemplateItem struct {
// 	Name string `json:"name" bson:"name"`
// }

// DefaultUserCustomizeSceneModel 預設UserCustomizeSceneModel
func DefaultUserCustomizeSceneModel() UserCustomizeSceneModel {
	return UserCustomizeSceneModel{Base: Base{TableName: config.CUSTOMIZE_SCENE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a UserCustomizeSceneModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) UserCustomizeSceneModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetMongoConn 設定connection
// func (a UserCustomizeSceneModel) SetMongoConn(conn mongo.Connection) UserCustomizeSceneModel {
// 	a.MongoConn = conn
// 	return a
// }

// FindAll 查詢所有自定義場景資訊
func (a UserCustomizeSceneModel) FindAll(userID string) ([]UserCustomizeSceneModel, error) {
	// log.Println("查詢該用戶所有自定義場景狀態(mongo，多筆)")

	items, err := a.MongoConn.FindMany(a.TableName, bson.M{"user_id": userID})
	if err != nil {
		return []UserCustomizeSceneModel{}, errors.New("錯誤: 無法取得自定義場景資訊(mongo)，請重新查詢")
	}

	return MapToUserCustomizeSceneModel(items), nil
}

// Find 查詢自定義場景資訊
func (a UserCustomizeSceneModel) Find(userID, topicID string) (UserCustomizeSceneModel, error) {
	// log.Println("查詢該用戶所有自定義場景狀態(mongo)，單筆")

	item, err := a.MongoConn.FindOne(a.TableName,
		bson.M{"user_id": userID, "topic_id": topicID})
	if err != nil {
		return UserCustomizeSceneModel{}, errors.New("錯誤: 無法取得自定義場景資訊(mongo)，請重新查詢")
	}

	a = a.MapToModel(item)

	return a, nil
}

// MapToModel map轉換model
func (a UserCustomizeSceneModel) MapToModel(m map[string]interface{}) UserCustomizeSceneModel {

	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &a)

	return a
}

// MapToUserCustomizeSceneModel map轉換[]UserCustomizeSceneModel
func MapToUserCustomizeSceneModel(items []bson.M) []UserCustomizeSceneModel {
	var scenes = make([]UserCustomizeSceneModel, 0)

	// json解碼，轉換成strcut
	b, _ := json.Marshal(items)
	json.Unmarshal(b, &scenes)

	// for _, item := range items {
	// 	var (
	// 		scene UserCustomizeSceneModel
	// 	)

	// 	// json解碼，轉換成strcut
	// 	b, _ := json.Marshal(item)
	// 	json.Unmarshal(b, &scene)

	// 	scenes = append(scenes, scene)
	// }
	return scenes
}

// Add 新增自定義場景資料(mongo)
func (a UserCustomizeSceneModel) Add(userID string, topicName string,
	data map[string]interface{}) error {

	// 取得mongo中的id資料(遞增處理)
	mongoID, _ := a.MongoConn.GetNextSequence(config.CUSTOMIZE_SCENE)

	_, err := a.MongoConn.InsertOne(config.CUSTOMIZE_SCENE, bson.M{
		"id":                   mongoID,
		"topic_id":             utils.UUID(20),
		"user_id":              userID,
		"topic_name":           topicName,
		"customize_scene_data": data,
	})

	if err != nil {
		return errors.New("錯誤: 無法新增自定義場景資料(mongo)，請重新操作")
	}

	return nil
}

// Update 更新自定義場景資料(mongo)
func (a UserCustomizeSceneModel) Update(topicID string, userID string, topicName string,
	data map[string]interface{}) error {

	// 更新資料庫
	_, err := a.MongoConn.UpdateOne(a.TableName,
		bson.M{"topic_id": topicID, "user_id": userID}, // 過濾參數
		bson.M{
			"$set": bson.M{
				"topic_name":           topicName,
				"customize_scene_data": data,
			}, // 更新參數
		})
	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		return errors.New("錯誤: 無法更新自定義場景資料(mongo)，請重新操作")
	}

	return nil
}
