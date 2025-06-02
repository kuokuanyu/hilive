package mongo

import (
	"context"
	"fmt"
	"hilive/modules/config"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo 資料庫引擎
type Mongo struct {
	Once      sync.Once                // sync.Once為唯一鎖，在代碼需要被執行時，只會被執行一次
	MongoList map[string]*mongo.Client // mongodb引擎
}

// GetMongo 取得引擎
func (m *Mongo) GetMongo(key string) *mongo.Client {
	return m.MongoList[key]
}

// GetDefaultMongo 設置mongodb
func GetDefaultMongo() *Mongo {
	return &Mongo{
		MongoList: make(map[string]*mongo.Client),
	}
}

// Name 引擎名稱
func (m *Mongo) Name() string {
	return "mongo"
}

// InitMongo 初始化 MongoDB 連線
func (m *Mongo) InitMongo(cfgs map[string]config.Mongo) Connection {
	m.Once.Do(func() {
		m.MongoList = make(map[string]*mongo.Client)

		for conn, cfg := range cfgs {
			// clientOptions := options.Client().ApplyURI(
			// 	"mongodb+srv://a167829435:199XoHi9hV9XrkU0@hilives-cluster.tbsz6.mongodb.net/hilive_dev?authSource=admin")

			clientOptions := options.Client().ApplyURI(
				fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
					cfg.User, cfg.Pwd, cfg.Host, cfg.Name))

			// log.Println(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
			// 	cfg.User, cfg.Pwd, cfg.Host, cfg.Name))

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, clientOptions)
			if err != nil {
				log.Printf("MongoDB連線失敗: %s\n %s", conn, err.Error())
				continue
			}

			// 測試連線
			if err := client.Ping(ctx, nil); err != nil {
				log.Printf("MongoDB測試連線發生問題: %s\n %s", conn, err.Error())
				continue
			}

			m.MongoList[conn] = client
			log.Printf("MongoDB成功連線: %s", conn)
		}
	})
	return m
}

// 取得資料庫名稱
// func (m *Mongo) GetDbName() string {
// 	db := m.MongoList["hilive"].Database("hilive_dev")
// 	log.Println("Database Name:", db.Name()) // 顯示資料庫名稱

// 	return db.Name()

// }

// GetDatabase 取得 MongoDB 資料庫
// func (m *Mongo) GetDatabase(connName string, dbName string) (*mongo.Database, error) {
// 	client, exists := m.MongoList[connName]
// 	if !exists {
// 		return nil, fmt.Errorf("MongoDB資料庫無法取得", connName)
// 	}
// 	return client.Database(dbName), nil
// }
