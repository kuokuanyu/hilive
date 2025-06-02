package cache

import (
	"context"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/utils"
	"log"
	"math"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

// var rdsLock sync.Mutex
// var cacheInstance *Redis

// Redis 資料庫引擎
type Redis struct {
	RedisList map[string]*redis.Pool // redis引擎
	Once      sync.Once              // sync.Once為唯一鎖，在代碼需要被執行時，只會被執行一次
	showDebug bool
}

// GetRedis 取得引擎
func (rds *Redis) GetRedis(key string) *redis.Pool {
	return rds.RedisList[key]
}

// GetDefaultRedis 設置redis
func GetDefaultRedis() *Redis {
	return &Redis{
		RedisList: make(map[string]*redis.Pool),
	}
}

// Name 引擎名稱
func (rds *Redis) Name() string {
	return "redis"
}

// ShowDebug 印出操作日誌
func (rds *Redis) ShowDebug(b bool) {
	rds.showDebug = b
}

// InitRedis 初始化redis引擎
func (rds *Redis) InitRedis(cfgs map[string]config.Redis) Connection {
	rds.Once.Do(func() {
		for conn, cfg := range cfgs {
			pool := redis.Pool{
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", cfg.Host+":"+cfg.Port)
					if err != nil {
						// log.Fatal("redis.InitRedis Dial error ", err)
						return nil, err
					}
					return c, nil
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					if time.Since(t) < time.Minute {
						return nil
					}
					_, err := c.Do("PING")
					return err
				},
				MaxIdle:         10000,           // 最大空閒連接數
				MaxActive:       10000,           // 最大連接數，0表示沒限制
				IdleTimeout:     1 * time.Minute, // 最大空閒時間
				Wait:            true,
				MaxConnLifetime: 0,
			}
			rds.RedisList[conn] = &pool
			rds.ShowDebug(false)
		}
	})
	return rds
}

// Do 對外只有一個命令，封裝了redis的命令
func (rds *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := rds.RedisList["hilives"].Get()
	defer conn.Close()

	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("redis Do", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		if err == nil {
			fmt.Printf("[redis] [info] [%dus]cmd=%s, err=nil, args=%v, reply=OK\n", (t2-t1)/1000, commandName, args)
		} else {
			fmt.Printf("[redis] [info] [%dus]cmd=%s, err=%s, args=%v, reply=%s\n", (t2-t1)/1000, commandName, err, args, reply)
		}
	}
	return reply, err
}

// Publish 發布訊息到指定的頻道
func (rds *Redis) Publish(channel, message string) (err error) {
	_, err = rds.Do("PUBLISH", channel, message)
	if err != nil {
		log.Println("錯誤: 發布訊息到Redis頻道發生問題, ", err)
		return err
	}
	return nil
}

// Subscribes 訂閱多個頻道的訊息
func (rds *Redis) Subscribes(ctx context.Context, channels []string, messageHandler func(channel, message string)) error {
	conn := rds.RedisList["hilives"].Get()
	// defer conn.Close()

	pubsubConn := redis.PubSubConn{Conn: conn}

	// 訂閱多個頻道
	for _, channel := range channels {
		err := pubsubConn.Subscribe(channel)
		if err != nil {
			log.Println("錯誤: 訂閱Redis頻道發生問題, ", err)
			conn.Close()
			return err
		}
	}

	go func() {
		defer pubsubConn.Close()
		for {
			select {
			case <-ctx.Done():
				// 如果 context 被取消，则退出
				// log.Println("訂閱被取消，退出")
				return
			default:
				switch v := pubsubConn.Receive().(type) {
				case redis.Message:
					// 當收到訊息時，呼叫 messageHandler 處理訊息
					// log.Println("訂閱者收到訊息: ", v.Channel, string(v.Data))

					messageHandler(v.Channel, string(v.Data))

					// 查看取消訂閱後的訂閱人數
					// numSub, _ := rds.getNumSub(channel)
					// log.Printf("發送訊息，%s 頻道的訂閱人數: %d\n", channel, numSub)

				case redis.Subscription:
					// log.Println("訂閱數?: ", v.Count)，開啟時會一直顯示1
					// 處理訂閱事件（如訂閱、取消訂閱等）
					if v.Count == 0 {
						// 所有頻道都取消訂閱，退出循環
						return
					}
				case error:
					// 處理錯誤並退出
					log.Println("錯誤: 訂閱Redis頻道過程中發生錯誤, ", v)
					return
				}
			}
		}
	}()

	return nil
}

// Subscribe 訂閱指定頻道的訊息
func (rds *Redis) Subscribe(ctx context.Context, channel string, messageHandler func(channel, message string)) error {
	conn := rds.RedisList["hilives"].Get()
	// defer conn.Close()

	pubsubConn := redis.PubSubConn{Conn: conn}
	err := pubsubConn.Subscribe(channel)
	if err != nil {
		log.Println("錯誤: 訂閱Redis頻道發生問題, ", err)
		conn.Close()
		return err
	}

	go func() {
		defer pubsubConn.Close()
		for {
			select {
			case <-ctx.Done():
				// 如果 context 被取消，则退出
				// log.Println("訂閱被取消，退出")
				return
			default:
				switch v := pubsubConn.Receive().(type) {
				case redis.Message:
					// 當收到訊息時，呼叫 messageHandler 處理訊息
					// log.Println("訂閱者收到訊息: ", v.Channel, string(v.Data))

					messageHandler(v.Channel, string(v.Data))

					// 查看取消訂閱後的訂閱人數
					// numSub, _ := rds.getNumSub(channel)
					// log.Printf("發送訊息，%s 頻道的訂閱人數: %d\n", channel, numSub)

				case redis.Subscription:
					// log.Println("訂閱數?: ", v.Count)，開啟時會一直顯示1
					// 處理訂閱事件（如訂閱、取消訂閱等）
					if v.Count == 0 {
						// 所有頻道都取消訂閱，退出循環
						return
					}
				case error:
					// 處理錯誤並退出
					log.Println("錯誤: 訂閱Redis頻道過程中發生錯誤, ", v)
					return
				}
			}
		}
	}()

	return nil
}

// Unsubscribe 取消訂閱指定頻道
func (rds *Redis) Unsubscribe(channel string) error {
	conn := rds.RedisList["hilives"].Get()
	defer conn.Close()

	pubsubConn := redis.PubSubConn{Conn: conn}
	// defer pubsubConn.Close()

	err := pubsubConn.Unsubscribe(channel)
	if err != nil {
		log.Println("錯誤: 取消訂閱Redis頻道發生問題, ", err)
		return err
	}

	// 主動關閉連接，確保不再接收訊息
	if err := pubsubConn.Conn.Close(); err != nil {
		log.Println("錯誤: 無法關閉PubSub連接, ", err)
		return err
	}

	// log.Println("取消訂閱成功: ", channel)
	return nil
}

// getNumSub 查看訂閱數
// func (rds *Redis) getNumSub(channel string) (int, error) {
// 	// 执行 PUBSUB NUMSUB 命令来获取频道的订阅数
// 	res, err := rds.Do("PUBSUB", "NUMSUB", channel)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// 将返回结果解析为 []interface{}
// 	values, err := redis.Values(res, err)
// 	if err != nil {
// 		return 0, err
// 	}

// 	if len(values) < 2 {
// 		return 0, fmt.Errorf("返回结果格式不正确")
// 	}

// 	// 返回的第一个值是频道名，第二个值是订阅数
// 	numSub, ok := values[1].(int64)
// 	if !ok {
// 		return 0, fmt.Errorf("无法转换订阅数为 int64")
// 	}

// 	return int(numSub), nil
// }

// GetCache 從緩存中取得所有資料(string方式)
func (rds *Redis) GetCache(key string) (str string, err error) {
	var rs interface{}
	// 讀取緩存
	rs, err = rds.Do("GET", key)
	if err != nil {
		// log.Println("錯誤: 取得redis(string)發生問題, ", err)
		return
	}
	str = utils.GetString(rs, "")

	return
}

// SetCache 將數據更新至緩存(string方式)
func (rds *Redis) SetCache(value ...interface{}) (reply interface{}, err error) {
	// 更新缓存
	reply, err = rds.Do("SET", value...)
	if err != nil {
		// log.Println("錯誤: 設置redis(string)發生問題, ", err)
		return
	}

	if len(value) == 2 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", value[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	} else {
		// 長度超過2的是設置redis鎖
		// log.Println("redis鎖")
	}

	return
}

// DelCache 清空緩存數據
func (rds *Redis) DelCache(key string) (reply interface{}, err error) {
	// 删除redis中的缓存
	reply, err = rds.Do("DEL", key)
	if err != nil {
		// log.Println("錯誤: 刪除redis發生問題, ", err)
		return
	}
	return
}

// IncrCache 緩存遞增，return遞增後的值(string方式)
func (rds *Redis) IncrCache(key string) int64 {
	// 遞增緩存
	rs, err := rds.Do("INCR", key)
	if err != nil {
		// log.Println("錯誤: 遞增redis資料(string)發生問題, ", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)

		// 設置過期時間
		_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}

		return num
	}
}

// DecrCache 緩存遞減，return遞減後的值(string方式)
func (rds *Redis) DecrCache(key string) int64 {
	// 遞增緩存
	rs, err := rds.Do("DECR", key)
	if err != nil {
		// log.Println("錯誤: 遞減redis資料(string)發生問題, ", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)

		// 設置過期時間
		_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}

		return num
	}
}

// HashGetCache 從緩存中取得單個資料(hash方式)
func (rds *Redis) HashGetCache(key, field string) (data string,
	err error) {
	// 讀取緩存
	data, err = redis.String(rds.Do("HGET", key, field))
	if err != nil {
		// log.Println("錯誤: 取得redis(hash, 單個)發生問題, ", err)
		return
	}
	return
}

// HashGetAllCache 從緩存中取得多個資料(hash方式)
func (rds *Redis) HashGetAllCache(key string) (dataMap map[string]string,
	err error) {
	// 讀取緩存
	dataMap, err = redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		// log.Println("錯誤: 取得redis(hash, 多個)發生問題, ", err)
		return
	}
	return
}

// HashGetAllCacheStrings 從緩存中取得多個資料(hash方式)
func (rds *Redis) HashGetAllCacheStrings(key string) (datas []string,
	err error) {
	// 讀取緩存
	datas, err = redis.Strings(rds.Do("HGETALL", key))
	if err != nil {
		// log.Println("錯誤: 取得redis(hash, 多個)發生問題, ", err)
		return
	}
	return
}

// HashSetCache 將數據更新至緩存(hash方式)，單個
func (rds *Redis) HashSetCache(key, field string, value interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("HSET", key, field, value)
	if err != nil {
		// log.Println("錯誤: 設置redis(hash, 單個)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// HashMultiSetCache 將數據更新至緩存(hash方式)，多個
func (rds *Redis) HashMultiSetCache(params []interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("HMSET", params...)
	if err != nil {
		// log.Println("錯誤: 設置redis(hash, 多個)發生問題, ", err)
		return
	}

	if len(params) > 0 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", params[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	}
	return
}

// HashDelCache 刪除緩存(hash方式)
func (rds *Redis) HashDelCache(key, field string) {
	// 更新缓存
	rds.Do("HDEL", key, field)
	// if err != nil {
	// 	log.Println("錯誤: 刪除redis(hash)發生問題, ", err)
	// 	return
	// }
	return
}

// HashIncrCache 緩存遞增，return遞增後的值(hash方式)
func (rds *Redis) HashIncrCache(key, field string) int64 {
	// 遞增緩存
	rs, err := rds.Do("HINCRBY", key, field, 1)
	if err != nil {
		// log.Println("錯誤: 遞增redis資料(hash)發生問題, ", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)

		// 設置過期時間
		_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}

		return num
	}
}

// HashDecrCache 緩存遞減，return遞減後的值(hash方式)
func (rds *Redis) HashDecrCache(key, field string) int64 {
	// 遞增緩存
	rs, err := rds.Do("HINCRBY", key, field, -1)
	if err != nil {
		// log.Println("錯誤: 遞減redis資料(hash)發生問題, ", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)

		// 設置過期時間
		_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}

		return num
	}
}

// SetExpire 設置過期時間
// func (rds *Redis) SetExpire(key, second string) {
// 	// 更新缓存
// 	rds.Do("EXPIRE", key, second)
// 	// if err != nil {
// 	// 	log.Println("錯誤: 設置redis過期時間發生問題, ", err)
// 	// 	return
// 	// }
// 	return
// }

// ListLPush 將資料更新至緩存，左側(list方式)
func (rds *Redis) ListLPush(key string, value interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("LPUSH", key, value)
	if err != nil {
		// log.Println("錯誤: 將資料更新至緩存，左側(list方式)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// ListRPush 將資料更新至緩存，右側(list方式)
func (rds *Redis) ListRPush(key string, value interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("RPUSH", key, value)
	if err != nil {
		// log.Println("錯誤: 將資料更新至緩存，右側(list方式)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// 將資料更新至緩存(多個)，右側(list方式)
func (rds *Redis) ListMultiRPush(params []interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("RPUSH", params...)
	if err != nil {
		// log.Println("錯誤: 將資料更新至緩存(多個)，右側(list方式)發生問題, ", err)
		return
	}

	if len(params) > 0 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", params[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	}

	return
}

// ListRange 從緩存中取得資料(list方式)
func (rds *Redis) ListRange(key string, start, stop int64) (datas []string, err error) {
	datas, err = redis.Strings(rds.Do("LRANGE", key, start, stop-1))
	if err != nil {
		// log.Println("錯誤: 從緩存中取得資料(list方式)發生問題, ", err)
		return
	}
	return
}

// ListLen 從緩存中取得資料數量(list方式)
func (rds *Redis) ListLen(key string) (count int64) {
	// 讀取緩存
	count, _ = redis.Int64(rds.Do("LLEN", key))
	// if err != nil {
	// 	log.Println("錯誤: 從緩存中取得資料數量(set)發生問題, ", err)
	// 	return
	// }
	return
}

// ListRem 將數據從緩存中清除(list方式)
func (rds *Redis) ListRem(key string, value interface{}) (err error) {
	// 刪除缓存
	_, err = rds.Do("LREM", key, 0, value)
	if err != nil {
		// log.Println("錯誤: 將數據從緩存中清除(list方式)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// SetCard 從緩存中取得資料數量(set方式)
func (rds *Redis) SetCard(key string) (count int64) {
	// 讀取緩存
	count, _ = redis.Int64(rds.Do("SCARD", key))
	// if err != nil {
	// 	log.Println("錯誤: 從緩存中取得資料數量(set)發生問題, ", err)
	// 	return
	// }
	return
}

// SetIsMember 從緩存中判斷是否有資料(set方式)
func (rds *Redis) SetIsMember(key string, value interface{}) (isExist bool) {
	// 讀取緩存
	isExist, _ = redis.Bool(rds.Do("SISMEMBER", key, value))
	// if err != nil {
	// 	log.Println("錯誤: 從緩存中判斷是否有資料(set)發生問題, ", err)
	// 	return
	// }
	return
}

// SetGetMembers 從緩存中取得所有資料(set方式)
func (rds *Redis) SetGetMembers(key string) (datas []string,
	err error) {
	// 讀取緩存
	datas, err = redis.Strings(rds.Do("SMEMBERS", key))
	if err != nil {
		// log.Println("錯誤: 從緩存中取得所有資料(set)發生問題, ", err)
		return
	}
	return
}

// SetAdd 將數據更新至緩存(set方式)
func (rds *Redis) SetAdd(params []interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("SADD", params...)
	if err != nil {
		// log.Println("錯誤: 將數據更新至緩存(set方式)發生問題, ", err)
		return
	}

	if len(params) > 0 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", params[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	}

	return
}

// SetRem 將數據從緩存中清除(set方式)
func (rds *Redis) SetRem(params []interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("SREM", params...)
	if err != nil {
		// log.Println("錯誤: 將數據從緩存中清除(set方式)發生問題, ", err)
		return
	}

	if len(params) > 0 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", params[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	}

	return
}

// ZSetAddInt 將分數更新至緩存(zset方式，整數)
func (rds *Redis) ZSetAddInt(key, member string, score int64) (err error) {
	_, err = rds.Do("ZADD", key, score, member)
	if err != nil {
		// log.Println("錯誤: 將分數更新至緩存(zset方式，整數)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// ZSetAddFloat 將分數更新至緩存(zset方式，小數)
func (rds *Redis) ZSetAddFloat(key, member string, score float64) (err error) {
	_, err = rds.Do("ZADD", key, score, member)
	if err != nil {
		// log.Println("錯誤: 將分數更新至緩存(zset方式，小數)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}

// ZSetMultiAdd 將數據更新至緩存，多個(zset方式)
func (rds *Redis) ZSetMultiAdd(params []interface{}) (err error) {
	// 更新缓存
	_, err = rds.Do("ZADD", params...)
	if err != nil {
		// log.Println("錯誤: 將數據更新至緩存，多個(zset方式)發生問題, ", err)
		return
	}

	if len(params) > 0 {
		// 設置過期時間
		_, err = rds.Do("EXPIRE", params[0], config.REDIS_EXPIRE)
		if err != nil {
			log.Println("錯誤: redis設置過期時間發生問題")
		}
	}

	return
}

// ZSetIncrCache 將分數遞增至緩存(zset方式，整數)
func (rds *Redis) ZSetIncrCache(key, member string, score int64) (err error) {
	// 更新缓存
	_, err = rds.Do("ZINCRBY", key, score, member)
	if err != nil {
		// log.Println("錯誤: 遞增redis資料(zset方式)發生問題, ", err)
		return err
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
	// else {
	// 	num := rs.(int64)
	// 	return num
	// }
}

// ZSetRange 從緩存中取得分數由低到高的資料(zset方式)
func (rds *Redis) ZSetRange(key string, start, stop int64) (datas []string, err error) {
	datas, err = redis.Strings(rds.Do("ZRANGE", key, start, stop-1))
	if err != nil {
		// log.Println("錯誤: 從緩存中取得分數由低到高的資料(zset方式)發生問題, ", err)
		return
	}
	return
}

// ZSetRevRange 從緩存中取得分數由高到低的資料(zset方式)
func (rds *Redis) ZSetRevRange(key string, start, stop int64) (datas []string, err error) {
	datas, err = redis.Strings(rds.Do("ZREVRANGE", key, start, stop-1))
	if err != nil {
		log.Println("錯誤: 從緩存中取得分數由高到低的資料(zset方式)發生問題, ", err)
		return
	}
	return
}

// ZSetIntScore 從緩存中取得分數資料(zset方式，整數)
func (rds *Redis) ZSetIntScore(key, member string) (value int64) {
	value, _ = redis.Int64(rds.Do("ZSCORE", key, member))
	return
}

// ZSetFloatScore 從緩存中取得分數資料(zset方式，小數)
func (rds *Redis) ZSetFloatScore(key, member string) (value float64) {
	value, _ = redis.Float64(rds.Do("ZSCORE", key, member))
	return
}

// ZSetRevRank 從緩存中取得用戶的排名資料(zset方式)
func (rds *Redis) ZSetRevRank(key, member string) (value int64) {
	value, _ = redis.Int64(rds.Do("ZREVRANK", key, member))
	return
}

// ZSetRem 從緩存中移除用戶的分數資料(zset方式)
func (rds *Redis) ZSetRem(key, member string) (err error) {
	// fmt.Println("清除用戶分數資料")
	_, err = rds.Do("ZREM", key, member)
	if err != nil {
		// log.Println("錯誤: 從緩存中移除用戶的分數資料(zset方式)發生問題, ", err)
		return
	}

	// 設置過期時間
	_, err = rds.Do("EXPIRE", key, config.REDIS_EXPIRE)
	if err != nil {
		log.Println("錯誤: redis設置過期時間發生問題")
	}

	return
}
