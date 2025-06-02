package controller

import (
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
)

func Test_GetLucky(t *testing.T) {

	// cfg := config.Config{
	// 	Databases: config.DatabaseList{
	// 		config.MYSQL_ENGINE: {
	// 			Host:       config.MYSQL_HOST,
	// 			Port:       config.MYSQL_PORT,
	// 			User:       config.MYSQL_USER,
	// 			Pwd:        config.MYSQL_PASSWORD,
	// 			Name:       config.MYSQL_NAME,
	// 			MaxIdleCon: config.MYSQL_MAXIDLECON,
	// 			MaxOpenCon: config.MYSQL_MAXOPENCON,
	// 			Driver:     config.MYSQL_DRIVER,
	// 		},
	// 	},
	// 	RedisList: config.RedisList{
	// 		config.REDIS_ENGINE: {
	// 			Host: config.REDIS_HOST,
	// 			Port: config.REDIS_PORT,
	// 		},
	// 	},
	// 	Prefix: config.PREFIX,
	// 	Store: config.Store{
	// 		Path:   config.STORE_PATH,
	// 		Prefix: config.STORE_PREFIX,
	// 	},
	// }
	// eng := engine.DefaultEngine().InitDatabase(cfg).InitRedis(cfg).
	// 	SetEngine().InitRouter()
	// fmt.Println("eng: ", eng)
}
