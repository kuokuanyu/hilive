package engine

import (
	"hilive/modules/config"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
)

var (
	cfg = config.Config{
		Databases: config.DatabaseList{
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
		},
		RedisList: config.RedisList{
			config.REDIS_ENGINE: {
				Host: config.REDIS_HOST,
				Port: config.REDIS_PORT,
			},
		},
		Prefix: config.PREFIX,
		Store: config.Store{
			Path:   config.STORE_PATH,
			Prefix: config.STORE_PREFIX,
		},
	}
	eng = DefaultEngine().InitDatabase(cfg).InitRedis(cfg).
		SetEngine().InitRouter()
)
