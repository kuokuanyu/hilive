package models

import (
	"errors"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"

	"github.com/twilio/twilio-go"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/rickb777/date"
	"golang.org/x/crypto/bcrypt"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	MaxRetries     = 10000000000            // 重試次數
	RetryDelay     = 100 * time.Millisecond // 每次重試的延遲
	LockExpiration = 2                      // 鎖的有效期
)

var (
	client = twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username:   config.ACCOUNT_SID,
			Password:   config.AUTH_TOKEN,
			AccountSid: config.ACCOUNT_SID,
		})
)

// Base 紀錄資料表資訊
type Base struct {
	TableName string
	DbConn    db.Connection
	RedisConn cache.Connection
	MongoConn mongo.Connection
	// Tx        *sql.Tx
}

// DeleteModel 資料表欄位
type DeleteModel struct {
	ID         string `json:"id" example:"id1,id2,id3..."`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	UserID     string `json:"user_id" example:"user_id"`
	Token      string `json:"token" example:"token"`
}

// FilterFields 過濾欄位資料，只取得需要的欄位資料寫入map中(寫入資料表時使用)
func FilterFields(source map[string]interface{}, fields []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		if val, ok := source[field]; ok {
			result[field] = val
		}
	}

	// log.Println("過濾欄位後的資料: ", result)
	return result
}

// acquireLock 嘗試獲取 Redis 鎖
func acquireLock(redisConn cache.Connection, lockKey string, expiration int) (interface{}, error) {
	return redisConn.SetCache(lockKey, "locked", "NX", "EX", expiration)
}

// releaseLock 釋放 Redis 鎖
func releaseLock(redisConn cache.Connection, lockKey string) (interface{}, error) {
	return redisConn.DelCache(lockKey)
}

// Table 設置SQL
func (b Base) Table(table string, conn ...string) *db.SQL {
	if len(conn) == 0 {
		return db.Table(table).WithConn(b.DbConn)
	}
	return db.Table(table).WithConn(b.DbConn).WithConnName(conn[0])
}

// SendMessage 發送簡訊
func SendMessage(phone, message string) error {
	// 设置短信参数
	params := &openapi.CreateMessageParams{}
	params.SetTo("+886" + phone[1:])
	params.SetFrom(config.PHONE) // 发送者的 Twilio 电话号码
	params.SetBody(message)

	// 发送短信
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return errors.New(fmt.Sprintf("錯誤: 發送簡訊發生問題, %s", err.Error()))
	}

	return nil
}

// SendMail 發送郵件
func SendMail(to, message string) error {
	// 設置認證信息
	auth := smtp.PlainAuth("", config.GMAIL, config.GMAIL_PASSWORD, config.GMAIL_SMTP_HOST)

	// 發送郵件
	err := smtp.SendMail(config.GMAIL_SMTP_HOST+":"+config.GMAIL_SMTP_PORT,
		auth, config.GMAIL, []string{to}, []byte(message))
	if err != nil {
		return errors.New(fmt.Sprintf("錯誤: 發送郵件發生問題, %s", err.Error()))
	}

	return nil
}

// CompareDate 日期比較
func CompareDate(start, end string) bool {
	var (
		s             = strings.Split(start, "-")
		startYear, _  = strconv.Atoi(s[0])
		startMonth, _ = strconv.Atoi(s[1])
		startDay, _   = strconv.Atoi(s[2])
		startTime     = date.New(startYear, time.Month(startMonth), startDay)
		e             = strings.Split(end, "-")
		endYear, _    = strconv.Atoi(e[0])
		endMonth, _   = strconv.Atoi(e[1])
		endDay, _     = strconv.Atoi(e[2])
		endTime       = date.New(endYear, time.Month(endMonth), endDay)
	)
	if endTime.Sub(startTime) < 0 {
		return false
	}
	return true
}

// CompareTime 時間比較
func CompareTime(start, end string) bool {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+end, time.Local)
	boolTime := endTime.After(startTime) && startTime.Before(endTime)
	if boolTime == false && start != end {
		return false
	}
	return true

}

// CompareDatetime 日期時間比較
func CompareDatetime(start, end string) bool {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	boolTime := endTime.After(startTime) && startTime.Before(endTime)
	if boolTime == false && start != end {
		return false
	}
	return true
}

// CompareTDatetime 日期時間比較
func CompareTDatetime(start, end string) bool {
	startTime, _ := time.ParseInLocation("2006-01-02T15:04", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02T15:04", end, time.Local)
	boolTime := endTime.After(startTime) && startTime.Before(endTime)
	if boolTime == false && start != end {
		return false
	}
	return true
}

// lineBot 開啟line bot
func lineBot(secret, token string) (bot *linebot.Client) {
	if bot, err := linebot.New(secret, token); err != nil {
		panic("錯誤: 開啟LINE BOT發生問題")
	} else {
		return bot
	}
}

// pushMessage 發送訊息
func pushMessage(secret, token, userid, message string) error {
	if _, err := lineBot(secret, token).PushMessage(userid, linebot.NewTextMessage(message)).Do(); err != nil {
		return errors.New("錯誤: 傳送Line訊息發生問題，請重新傳送")
	}
	return nil
}

// interfaces []string轉換成[]interface
func interfaces(arr []string) []interface{} {
	var iarr = make([]interface{}, len(arr))
	for key, v := range arr {
		iarr[key] = v
	}
	return iarr
}

// GetURLParam 處理後的url
// func GetURLParam(u string) (string, url.Values) {
// 	m := make(url.Values)
// 	urr := strings.Split(u, "?")
// 	if len(urr) > 1 {
// 		m, _ = url.ParseQuery(urr[1])
// 	}
// 	return urr[0], m
// }

// methodInSlice string是否在slice裡
// func methodInSlice(arr []string, str string) bool {
// 	for i := 0; i < len(arr); i++ {
// 		if strings.EqualFold(arr[i], str) {
// 			return true
// 		}
// 	}
// 	return false
// }

// checkParam 檢查url是否符合
// func checkParam(src, comp url.Values) bool {
// 	if len(comp) == 0 {
// 		return true
// 	}
// 	if len(src) == 0 {
// 		return false
// 	}
// 	for key, value := range comp {
// 		v, find := src[key]
// 		if !find {
// 			return false
// 		}
// 		if len(value) == 0 {
// 			continue
// 		}
// 		if len(v) == 0 {
// 			return false
// 		}
// 		for i := 0; i < len(v); i++ {
// 			if v[i] == value[i] {
// 				continue
// 			} else {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// ReleaseConn 清空user.DbConn
func (user UserModel) ReleaseConn() UserModel {
	user.DbConn = nil
	return user
}

// EncodePassword 加密
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}
