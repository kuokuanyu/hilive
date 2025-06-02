package cache

import (
	"context"
	"hilive/modules/config"
)

// Connection redis處理程序
type Connection interface {
	Name() string
	InitRedis(cfg map[string]config.Redis) Connection
	// SetExpire(key, second string) // 設置過期時間

	// pub.sub
	Publish(channel, message string) error                                                                 // 發布
	Subscribe(ctx context.Context, channel string, messageHandler func(channel, message string)) error     // 訂閱
	Subscribes(ctx context.Context, channels []string, messageHandler func(channel, message string)) error // 訂閱多個
	Unsubscribe(channel string) error                                                                      // 取消訂閱

	// string
	GetCache(key string) (string, error)                 // 從緩存中取得所有資料(string方式)
	SetCache(params ...interface{}) (interface{}, error) // 將數據更新至緩存(string方式)
	DelCache(key string) (interface{}, error)            // 清空緩存數據
	IncrCache(key string) int64                          // 緩存遞增，return遞增後的值(string方式)
	DecrCache(key string) int64                          // 緩存遞減，return遞減後的值(string方式)

	// hash
	HashGetCache(key, field string) (string, error)          // 從緩存中取得單個資料(hash方式)，查詢不到會有錯誤
	HashGetAllCache(key string) (map[string]string, error)   // 從緩存中取得多個資料(hash方式)
	HashGetAllCacheStrings(key string) ([]string, error)     // 從緩存中取得多個資料(hash方式)
	HashSetCache(key, field string, value interface{}) error // 將數據更新至緩存(hash方式)，單個
	HashMultiSetCache(params []interface{}) error            // 將數據更新至緩存(hash方式)，多個
	HashDelCache(key, field string)                          // 刪除緩存(hash方式)
	HashIncrCache(key, field string) int64                   // 緩存遞增，return遞增後的值(hash方式)
	HashDecrCache(key, field string) int64                   // 緩存遞減，return遞減後的值(hash方式)

	// list
	ListLPush(key string, value interface{}) error             // 將資料更新至緩存，左側(list方式)
	ListRPush(key string, value interface{}) error             // 將資料更新至緩存，右側(list方式)
	ListMultiRPush(params []interface{}) error                 // 將資料更新至緩存(多個)，右側(list方式)
	ListRange(key string, start, stop int64) ([]string, error) // 從緩存中取得資料(list方式)
	ListLen(key string) int64                                  // 從緩存中取得資料數量(list方式)
	ListRem(key string, value interface{}) error               // 將數據從緩存中清除(list方式)

	// set
	SetCard(key string) int64                       // 從緩存中取得資料數量(set方式)
	SetIsMember(key string, value interface{}) bool // 從緩存中判斷是否有key的資料(set方式)
	SetGetMembers(key string) ([]string, error)     // 從緩存中取得所有資料(set方式)
	SetAdd(params []interface{}) error              // 將數據更新至緩存(set方式)
	SetRem(params []interface{}) error              // 將數據從緩存中清除(set方式)

	// zset
	ZSetAddInt(key, member string, score int64) error             // 將分數更新至緩存(zset方式，整數)
	ZSetAddFloat(key, member string, score float64) error         // 將分數更新至緩存(zset方式，小數)
	ZSetMultiAdd(params []interface{}) error                      // 將分數更新至緩存，多個(zset方式)
	ZSetIncrCache(key, member string, score int64) error          // 將分數遞增至緩存(zset方式，整數)
	ZSetRange(key string, start, stop int64) ([]string, error)    // 從緩存中取得分數由低到高的資料(zset方式)
	ZSetRevRange(key string, start, stop int64) ([]string, error) // 從緩存中取得分數由高到低的資料(zset方式)
	ZSetIntScore(key, member string) int64                        // 從緩存中取得分數資料(zset方式，整數)
	ZSetFloatScore(key, member string) float64                    // 從緩存中取得分數資料(zset方式，小數)
	ZSetRevRank(key string, member string) int64                  // 從緩存中取得用戶的排名資料(zset方式)
	ZSetRem(key string, member string) error                      // 從緩存中移除用戶的分數資料(zset方式)
}

// GetConnection 取得Connection
func GetConnection() Connection {
	return GetDefaultRedis()
}

// GetConnectionFromService 取得Connection
func GetConnectionFromService(s interface{}) Connection {
	if c, ok := s.(Connection); ok {
		return c
	}
	panic("錯誤的Service")
}
