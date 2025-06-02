package models

import (
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	DbCfgs = map[string]config.Database{
		config.MYSQL_ENGINE: {
			Host:       config.MYSQL_HOST,
			Port:       config.MYSQL_PORT,
			User:       config.MYSQL_USER,
			Pwd:        config.MYSQL_PASSWORD,
			Name:       config.MYSQL_NAME,
			MaxIdleCon: config.MYSQL_MAXIDLECON,
			MaxOpenCon: config.MYSQL_MAXOPENCON,
			Driver:     config.MYSQL_DRIVER,
		},
	}
	RedisCfgs = config.RedisList{
		config.REDIS_ENGINE: {
			Host: "127.0.0.1",
			Port: config.REDIS_PORT,
		},
	}
	MongoCfgs = config.MongoList{
		config.MONGO_ENGINE: {
			Host: config.MONGO_HOST,
			Port: config.MONGO_PORT,
			User: config.MONGO_USER,
			Pwd:  config.MONGO_PASSWORD,
			Name: config.MONGO_NAME,
		},
	}

	conn      = db.GetConnectionByDriver("mysql").InitDB(DbCfgs)
	redis     = cache.GetConnection().InitRedis(RedisCfgs)
	mongoConn = mongo.GetConnection().InitMongo(MongoCfgs)
)

// 測試發送簡訊功能
func Test_SendMessage(t *testing.T) {
	// 设置短信参数
	params := &openapi.CreateMessageParams{}
	params.SetTo("+886932530813")
	params.SetFrom(config.PHONE) // 发送者的 Twilio 电话号码
	params.SetBody("XXX 您好: 您已報名 XXX 活動，可利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): https://liff.line.me/1656920628-jwWm55v7?activity_id=XXX&user_id=XXX")

	// 发送短信
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		t.Error("發送簡訊發生問題")
	}
}
