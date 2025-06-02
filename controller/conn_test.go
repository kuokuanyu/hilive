package controller

import (
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
)

var (
	DbCfgs = map[string]config.Database{
		"hilive": {
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

	conn  = db.GetConnectionByDriver("mysql").InitDB(DbCfgs)
	redis = cache.GetConnection().InitRedis(RedisCfgs)
)
