package models

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// 從mysql查詢所有活動資料，寫入mongo中
func Test_Activity_Channel_Mongo(t *testing.T) {
	var data = make([]interface{}, 0)

	// 設定時區為 UTC+8
	// loc, _ := time.LoadLocation("Asia/Taipei")
	// now := time.Now().In(loc)
	// **轉換為 UTC 時間 + 8 小時**
	// now = now.Add(8 * time.Hour)

	// 格式化為 "YYYY-MM-DD HH:MM:SS"
	// formattedTime := now.Format("2006-01-02 15:04:05")

	// 查詢所有活動資料
	activitys, _ := conn.Query("SELECT * FROM activity;")

	for _, activity := range activitys {
		// 取得mongo中的id資料(遞增處理)
		mongoID, _ := mongoConn.GetNextSequence("activity_channel")

		data = append(data, bson.M{
			"id":             mongoID,
			"activity_id":    activity["activity_id"],
			"user_id":        activity["user_id"],
			"channel_1":      "close",
			"channel_2":      "close",
			"channel_3":      "close",
			"channel_4":      "close",
			"channel_5":      "close",
			"channel_6":      "close",
			"channel_7":      "close",
			"channel_8":      "close",
			"channel_9":      "close",
			"channel_10":     "close",
			"channel_amount": int64(5),
			// "created_at":     formattedTime, // 設置 `created_at`
			// "updated_at":     formattedTime, // 設置 `updated_at`
		})
	}

	_, err := mongoConn.InsertMany("activity_channel", data)
	// log.Println("result: ", results)
	log.Println("err: ", err)
}
