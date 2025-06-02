package models

import (
	"fmt"
	"log"
	"testing"
)

func Test_Line_User(t *testing.T) {

	// 取得users表中已綁定人員資料
	users, _ := conn.Query("SELECT * FROM users where bind = 'yes';")

	log.Println("綁定人數: ", len(users))

	for i := 0; i < len(users); i++ {

		userID, _ := users[i]["user_id"].(string) // 管理人員user_id

		// 取得line用戶資料
		// lineusers, _ := conn.Query(fmt.Sprintf("SELECT * FROM line_users where identify = '%s';", userID))

		// 將user_id資料寫入admin_id中
		// for i := 0; i < len(lineusers); i++ {

		conn.Exec(fmt.Sprintf(`update line_users set
				admin_id = '%s' where identify = '%s';`, userID, userID))
		// }
	}

}

// 新增測試人員資料
// func Test_Add_Line_User(t *testing.T) {
// 	names := []string{
// 		"Cindy",
// 		"羽柔",
// 		"陳靜雯",
// 		"Article Flowers",
// 		"甯靖",
// 		"哈哈",
// 		"小熙",
// 		"趙勇翔",
// 		"小魚魚(魚兒)",
// 		"貞美花藝",
// 		"劉皓",
// 		"ss",
// 		"Julia-郁熙",
// 		"Ristretto Freddo",
// 		"阮文聖",
// 		"阿美Amy",
// 		"志遠 David",
// 		"Mạc Thuận Thành",
// 		"_BiBi_",
// 		"文錡",
// 		"張",
// 		"阿兔",
// 		"朱喆強",
// 		"大熊B.",
// 		"王卓盛",
// 		"平孝志",
// 		"Nontapan Jivacate",
// 		"🐰兔兔🐰",
// 		"Beach 海",
// 		"Dream",
// 		"趙家豪",
// 		"娃娃",
// 		"當夜晚來臨時 我在窗外看著你",
// 		"莉莉絲",
// 		"森森",
// 		"C.P.A.",
// 		"克莉絲汀娜",
// 		"Peter豪",
// 		"stk;idw-",
// 		"Princess"}

// 	avatars := []string{
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/1.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/2.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/3.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/4.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/5.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/6.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/7.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/8.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/9.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/10.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/11.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/12.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/13.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/14.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/15.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/16.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/17.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/18.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/19.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/20.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/21.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/22.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/23.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/24.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/25.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/26.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/27.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/28.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/29.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/30.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/31.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/32.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/33.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/34.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/35.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/36.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/37.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/38.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/39.png",
// 		"https://dev.hilives.net/admin/uploads/system/fake_data/40.png"}
// 	// 新增測試人員資料
// 	for i := 1; i <= 40; i++ {
// 		if err := db.Table(config.LINE_USERS_TABLE).WithConn(conn).Where("user_id", "=", strconv.Itoa(i)).Update(command.Value{
// 			"name": names[i-1], "avatar": avatars[i-1],
// 			"email": i, "identify": i,
// 		}); err != nil {
// 			t.Error("錯誤: 新增人員發生問題")
// 		}
// 	}
// }
