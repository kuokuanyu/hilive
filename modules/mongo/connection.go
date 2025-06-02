package mongo

import (
	"context"
	"hilive/modules/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Connection mongodb處理程序
type Connection interface {
	Name() string
	// GetDbName() string // 取得資料庫名稱

	InitMongo(cfg map[string]config.Mongo) Connection // 初始化連線

	InsertOne(collectionName string, data interface{}) (*mongo.InsertOneResult, error)     // 插入單筆資料
	InsertMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error) // 插入多筆資料

	UpdateOne(collectionName string, filter bson.M, update bson.M) (*mongo.UpdateResult, error)  // 更新單筆資料
	UpdateMany(collectionName string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) // 更新多筆資料

	DeleteOne(collectionName string, filter bson.M) (*mongo.DeleteResult, error)  // 刪除單筆資料
	DeleteMany(collectionName string, filter bson.M) (*mongo.DeleteResult, error) // 刪除多筆資料

	FindOne(collectionName string, filter bson.M, opts ...*options.FindOneOptions) (bson.M, error)    // 查詢單筆資料
	FindMany(collectionName string, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) // 查詢多筆資料

	GetNextSequence(collectionName string) (int64, error) // 執行會遞增_id資料
}

// GetConnection 取得Connection
func GetConnection() Connection {
	return GetDefaultMongo()
}

// GetConnectionFromService 取得Connection
func GetConnectionFromService(s interface{}) Connection {
	if c, ok := s.(Connection); ok {
		return c
	}
	panic("錯誤的Service")
}

// GetNextSequence 執行會遞增_id資料
func (m *Mongo) GetNextSequence(collectionName string) (int64, error) {
	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	err := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection("counters").FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return int64(result.Seq), nil
}

// InsertOne 插入單筆資料
func (m *Mongo) InsertOne(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 自動補上 created_at 與 updated_at
	dataWithTimestamps := appendTimestamps(data)

	result, err := collection.InsertOne(ctx, dataWithTimestamps)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Mongo) InsertMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 自動補上 created_at 與 updated_at
	dataWithTimestamps := make([]interface{}, len(data))
	for i, d := range data {
		dataWithTimestamps[i] = appendTimestamps(d)
	}

	result, err := collection.InsertMany(ctx, dataWithTimestamps)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// appendTimestamps 執行新增時自動加入created_at、updated_at資料
func appendTimestamps(data interface{}) interface{} {
	// 設定時區為 UTC+8
	loc, _ := time.LoadLocation("Asia/Taipei")
	now := time.Now().In(loc)
	// 格式化為 "YYYY-MM-DD HH:MM:SS"
	formattedTime := now.Format("2006-01-02 15:04:05")

	switch d := data.(type) {
	case bson.M:
		d["created_at"] = formattedTime
		d["updated_at"] = formattedTime
		return d
	case map[string]interface{}:
		d["created_at"] = formattedTime
		d["updated_at"] = formattedTime
		return d
	default:
		// 如果是 struct，建議改用 bson.M 傳入，否則這裡無法直接處理
		return data
	}
}

// UpdateOne 更新單筆資料
func (m *Mongo) UpdateOne(collectionName string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 確保有 updated_at
	updateWithTime := mergeUpdateWithUpdatedAt(update)

	result, err := collection.UpdateOne(ctx, filter, updateWithTime)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateMany 更新多筆資料
func (m *Mongo) UpdateMany(collectionName string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 確保有 updated_at
	updateWithTime := mergeUpdateWithUpdatedAt(update)

	result, err := collection.UpdateMany(ctx, filter, updateWithTime)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// mergeUpdateWithUpdatedAt 執行更新時自動更新updated_at資料
func mergeUpdateWithUpdatedAt(update bson.M) bson.M {
	// 設定時區為 UTC+8
	loc, _ := time.LoadLocation("Asia/Taipei")
	now := time.Now().In(loc)
	// 格式化為 "YYYY-MM-DD HH:MM:SS"
	formattedTime := now.Format("2006-01-02 15:04:05")

	// 如果本來就有 $set，我們在裡面加 updated_at
	if setFields, ok := update["$set"].(bson.M); ok {
		setFields["updated_at"] = formattedTime
		update["$set"] = setFields
	} else {
		// 沒有 $set，我們創一個
		update["$set"] = bson.M{"updated_at": formattedTime}
	}
	return update
}

// DeleteOne 刪除單筆資料
func (m *Mongo) DeleteOne(collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteMany 刪除多筆資料
func (m *Mongo) DeleteMany(collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOne 查詢單筆資料
func (m *Mongo) FindOne(collectionName string, filter bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindMany 查詢多筆資料
func (m *Mongo) FindMany(collectionName string, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	collection := m.MongoList["hilives"].Database(config.MONGO_NAME).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// opts := options.Find()
	// if limit > 0 {
	// 	opts.SetLimit(int64(limit)) // limit
	// 	opts.SetSort() // 排序，ex: opts := options.Find().SetSort(bson.D{{欄位名, 1}})，1升冪.-1降冪
	// 	opts.SetProjection() // 特定欄位，ex: options.Find().SetProjection(bson.M{"欄位名": 1})
	// }

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
