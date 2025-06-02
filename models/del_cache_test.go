package models

import "testing"

// 刪除redis資料
func Test_DelCache_Data(t *testing.T) {
	redis.DelCache("QA_record_game")
}
