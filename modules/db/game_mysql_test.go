package db

import (
	"fmt"
	"hilive/modules/cache"
	"hilive/modules/config"
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
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
	conn  = GetConnectionByDriver("mysql").InitDB(DbCfgs)
	redis = cache.GetConnection().InitRedis(RedisCfgs)
)

// 資料表增加遊戲自定義主題欄位
func Test_Mysql_Add_Column(t *testing.T) {
	// ex : ALTER TABLE activity_game_lottery_picture ADD lottery_starrysky_35 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';
	var (
		table       = "activity_game_lottery_picture" // 資料表
		game        = "lottery"                       // 遊戲
		topic       = "starrysky"                     // 主題
		start       = 36                              // 開始欄位
		end         = 40                              // 結束欄位
		startMusic  = true                            // 是否建立開始音樂欄位
		gamingMusic = false                           // 是否建立遊戲中音樂欄位
		endMusic    = true                            // 是否建立結束音樂欄位
		query       = "ALTER TABLE %s ADD %s_%s_%d varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"
		query2      = "ALTER TABLE %s ADD %s_%s_%s varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"
	)

	// 自定義欄位
	// ex : ALTER TABLE lottery_bgm_gaming ADD lottery_starrysky_35 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';
	for i := start; i <= end; i++ {
		// conn.ExecWithConnection("hilive", fmt.Sprintf(query, table, game, topic, i))
		conn.Exec(fmt.Sprintf(query, table, game, topic, i))
	}

	// 開始音樂
	if startMusic {
		conn.Exec(fmt.Sprintf(query2, table, game, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		conn.Exec(fmt.Sprintf(query2, table, game, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		conn.Exec(fmt.Sprintf(query2, table, game, "bgm", "end"))
	}
}
