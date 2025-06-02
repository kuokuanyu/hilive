package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
	"hilive/modules/utils"
	"log"
	"strconv"
	"strings"
)

// ApplysignUserModel 資料表欄位
type ApplysignUserModel struct {
	Base     `json:"-"`
	ID       int64  `json:"id"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	ExtEmail string `json:"ext_email"`
	Phone    string `json:"phone"`

	ActivityID string `json:"activity_id"` // 自定義簽到人員活動主辦方
	// ExtAccount  string `json:"ext_account"`  // 自定義簽到人員用戶
	ExtPassword string `json:"ext_password"` // 自定義簽到人員密碼

	// 自定義欄位
	Ext1  string `json:"ext_1"`
	Ext2  string `json:"ext_2"`
	Ext3  string `json:"ext_3"`
	Ext4  string `json:"ext_4"`
	Ext5  string `json:"ext_5"`
	Ext6  string `json:"ext_6"`
	Ext7  string `json:"ext_7"`
	Ext8  string `json:"ext_8"`
	Ext9  string `json:"ext_9"`
	Ext10 string `json:"ext_10"`

	Source string `json:"source"` // 來源

	People string `json:"people"` // 匯入資料數
}

// DefaultApplysignUserModel 預設ApplysignUserModel
func DefaultApplysignUserModel() ApplysignUserModel {
	return ApplysignUserModel{Base: Base{TableName: config.LINE_USERS_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a ApplysignUserModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) ApplysignUserModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (user ApplysignUserModel) SetDbConn(conn db.Connection) ApplysignUserModel {
// 	user.DbConn = conn
// 	return user
// }

// // SetRedisConn 設定connection
// func (user ApplysignUserModel) SetRedisConn(conn cache.Connection) ApplysignUserModel {
// 	user.RedisConn = conn
// 	return user
// }

// FindAccount 尋找自定義簽到人員資料
func (user ApplysignUserModel) FindAccount(activityID, password string) (ApplysignUserModel, error) {
	item, err := user.Table(user.Base.TableName).
		Where("activity_id", "=", activityID).
		// Where("ext_account", "=", account).
		Where("ext_password", "=", password).
		Where("device", "=", "customize").
		First()
	if err != nil {
		return ApplysignUserModel{}, errors.New("錯誤: 無法取得用戶資訊，請重新查詢")
	}
	if item == nil {
		return ApplysignUserModel{}, nil
	}
	return user.MapToModel(item), nil
}

// MapToModel 設置ApplysignUserModel
func (user ApplysignUserModel) MapToModel(m map[string]interface{}) ApplysignUserModel {
	// json解碼，轉換成strcut
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &user)

	return user
}

// user.ID, _ = m["id"].(int64)
// user.ActivityID, _ = m["activity_id"].(string)
// user.UserID, _ = m["user_id"].(string)
// user.Name, _ = m["name"].(string)
// user.Avatar, _ = m["avatar"].(string)
// user.Phone, _ = m["phone"].(string)
// user.ExtEmail, _ = m["ext_email"].(string)

// // user.ExtAccount, _ = m["ext_account"].(string)
// user.ExtPassword, _ = m["ext_password"].(string)

// Update 更新簽到人員資料
// func (user ApplysignUserModel) Update(model ApplysignUserModel) error {
// 	var (
// 		fields = []string{"name", "avatar"} // "phone", "ext_email",, "ext_password"

// 		values = []string{model.Name, model.Avatar} // model.Phone, model.ExtEmail,, model.ExtPassword

// 		fieldValues = command.Value{}
// 	)

// 	// if model.Phone != "" {
// 	// 	if len(model.Phone) > 2 {
// 	// 		if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
// 	// 			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 	// 		}
// 	// 	} else {
// 	// 		return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
// 	// 	}
// 	// }

// 	// if model.ExtEmail != "" && !strings.Contains(model.ExtEmail, "@") {
// 	// 	return errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
// 	// }

// 	for i, value := range values {
// 		if value != "" {
// 			fieldValues[fields[i]] = value
// 		}
// 	}
// 	if len(fieldValues) == 0 {
// 		return nil
// 	}

// 	err := user.Table(user.TableName).
// 		Where("activity_id", "=", model.ActivityID).
// 		Where("user_id", "=", model.UserID).
// 		Update(fieldValues)
// 	if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return errors.New("錯誤: 更新用戶資訊發生問題，請重新操作")
// 	}
// 	return nil
// }

// Adds 新增多筆用戶
func (user ApplysignUserModel) Adds(isRedis bool, model ApplysignUserModel) error {
	// log.Println("測試批量新增報名簽到用戶，觀察頭像欄位資料: ")
	// 取得匯入總資料數
	dataAmount, err := strconv.Atoi(model.People)
	if err != nil {
		return errors.New("錯誤: 人數資料發生問題，請重新操作")
	}

	if dataAmount > 0 {
		// 取得活動資料
		activityModel, err := DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			Find(false, model.ActivityID)
		if err != nil || activityModel.ID == 0 {
			return errors.New("錯誤: 取得活動資料發生問題，請重新操作")
		}

		// 處理匯入的參數資料
		var (
			attend           = activityModel.Attend                  // 活動目前參加人數
			maxPeople        = activityModel.People                  // 活動上限人數
			pushPhoneMessage = activityModel.PushPhoneMessage        // 傳送簡訊
			sendMail         = activityModel.SendMail                // 傳送郵件
			messageAmount    = activityModel.MessageAmount           // 簡訊數量
			mailAmount       = activityModel.MailAmount              // 郵件數量
			signCheck        = activityModel.SignCheck               // 簽到審核
			number           = activityModel.Number                  // 抽獎號碼
			activityID       = make([]string, dataAmount)            // 活動ID
			userID           = make([]string, dataAmount)            // 用戶ID
			name             = strings.Split(model.Name, ",")        // 名稱
			avatar           = make([]string, dataAmount)            // 頭像
			phone            = strings.Split(model.Phone, ",")       // 電話
			email            = strings.Split(model.ExtEmail, ",")    // 信箱
			extPassword      = strings.Split(model.ExtPassword, ",") // 驗證碼
			numbers          = make([]string, dataAmount)            // 抽獎號碼
			status           = make([]string, dataAmount)            // 用戶狀態
			modelExts        = []string{
				model.Ext1, model.Ext2, model.Ext3, model.Ext4,
				model.Ext5, model.Ext6, model.Ext7, model.Ext8,
				model.Ext9, model.Ext10,
			} // 前端自定義欄位參數
			exts  = make([][]string, 10) // 所有自定義欄位數值(10個)
			signs = make([]string, 0)    // 簽到人員(寫入redis中)
			// ext1             = make([]string, dataAmount)
			// ext2             = make([]string, dataAmount)
			// ext3             = make([]string, dataAmount)
			// ext4             = make([]string, dataAmount)
			// ext5             = make([]string, dataAmount)
			// ext6             = make([]string, dataAmount)
			// ext7             = make([]string, dataAmount)
			// ext8             = make([]string, dataAmount)
			// ext9             = make([]string, dataAmount)
			// ext10            = make([]string, dataAmount)
		)

		// 處理活動.用戶ID陣列資料
		for i := range userID {
			activityID[i] = model.ActivityID
			userID[i] = model.ActivityID + "_" + utils.UUID(16)
			avatar[i] = activityModel.CustomizeDefaultAvatar
		}

		// 處理自定義陣列資料
		for i, ext := range modelExts {
			if ext == "no data" {
				// 沒有資料，所有資料為空值
				exts[i] = make([]string, dataAmount)
			} else {
				// 將參數資料split
				exts[i] = strings.Split(ext, ",")
			}
		}

		// 處理所有人員status資料
		if signCheck == "close" { // 簽到審核關閉，加入資料status: sign
			if maxPeople-attend < int64(dataAmount) { // 可簽到人數小於匯入資料數，status部分資料為sign.apply
				signPeople := maxPeople - attend // 可簽到人數
				for i := range status {
					if i < int(signPeople) {
						status[i] = "sign"
						numbers[i] = strconv.Itoa(int(number) + i) // 陣列裡的number資料持續遞增
						signs = append(signs, userID[i])           // 將簽到狀態人員寫入陣列中
					} else {
						status[i] = "apply"
						// apply的號碼資料全都是0
						numbers[i] = "0"
					}
				}

				// 增加活動參加人數跟號碼資訊(增加至上限值)
				attend = attend + signPeople
				number = number + signPeople
			} else {
				// 可增加簽到人數大於匯入資料，加入資料全部status: sign
				for i := range status {
					status[i] = "sign"
					numbers[i] = strconv.Itoa(int(number) + i) // 陣列裡的number資料持續遞增
					signs = append(signs, userID[i])           // 將簽到狀態人員寫入陣列中
				}

				// 增加活動參加人數跟號碼資訊
				attend = attend + int64(dataAmount)
				number = number + int64(dataAmount)
			}

			// 更新參加人數跟號碼資訊
			if err = DefaultActivityModel().
				SetConn(user.DbConn, user.RedisConn, user.MongoConn).
				UpdateAttendAndNumber(true, model.ActivityID, int(attend), int(number), signs); err != nil &&
				err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				return err
			}

		} else if signCheck == "open" {
			// 簽到審核開啟，加入資料status: apply
			for i := range status {
				status[i] = "apply"
				// apply的號碼資料全都是0
				numbers[i] = "0"
			}
		}

		// 將資料匹量寫入line_users表中
		err = user.Table(config.LINE_USERS_TABLE).BatchInsert(
			dataAmount, "user_id,name,avatar,phone,ext_email,identify,activity_id,ext_password",
			[][]string{userID, name, avatar, phone, email, userID, activityID, extPassword})
		if err != nil {
			return errors.New("錯誤: 匹量新增用戶發生問題(line_users)，請重新操作")
		}

		// 將資料匹量寫入activity_applysign表中
		err = user.Table(config.ACTIVITY_APPLYSIGN_TABLE).BatchInsert(
			dataAmount, "user_id,activity_id,status,number,ext_1,ext_2,ext_3,ext_4,ext_5,ext_6,ext_7,ext_8,ext_9,ext_10",
			[][]string{userID, activityID, status, numbers,
				exts[0], exts[1], exts[2], exts[3], exts[4], exts[5], exts[6], exts[7], exts[8], exts[9]})

		if err != nil {
			return errors.New("錯誤: 匹量新增用戶發生問題(activity_applysign)，請重新操作")
		}

		// 郵件處理
		if sendMail == "open" && mailAmount > 0 {
			var (
				remain     int64
				sendAmount int64
			)
			// 判斷郵件剩餘數量是否足夠傳給所有用戶
			if mailAmount < int64(dataAmount) { // 匯入資料大於郵件數量
				// 將剩餘郵件數量為0
				remain = 0
				sendAmount = mailAmount
			} else { // 郵件數量大於匯入資料
				remain = mailAmount - int64(dataAmount)
				sendAmount = int64(dataAmount)
			}

			// 更新郵件剩餘數量
			err = DefaultActivityModel().
				SetConn(user.DbConn, user.RedisConn, user.MongoConn).
				UpdateMailAmount(true, model.ActivityID, remain)
			if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				// 其他錯誤
				return err
			}

			for i := 0; i < int(sendAmount); i++ {
				var message string
				url := fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, model.ActivityID, userID[i], "")
				qrcode := fmt.Sprintf(config.HTTPS_APPLYSIGN_QRCODE_URL, config.HILIVES_NET_URL, fmt.Sprintf("activity_id=%s.user_id=%s", model.ActivityID, userID[i]))
				qrcodeMessage := "<p></p>" // qrcode訊息

				subject := fmt.Sprintf("Subject: %s 活動報名簽到訊息\r\n", activityModel.ActivityName)
				mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

				// 判斷主持人掃描qrcode開關是否開啟
				if activityModel.HostScan == "open" {
					qrcodeMessage = fmt.Sprintf(`
				<p>以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到判斷(避免資料洩漏，連結勿提供給他人使用)：</p>
				<img src="https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=%s" alt="簽到QRcode"/>
				`, qrcode)
				}

				body := fmt.Sprintf(`
			<html>
			 	<body>
					<p>%s 您好，</p>
					<p>感謝您報名參加 %s 活動！以下是您的簽到信息：</p>
					<p>該活動的驗證碼為%s，可透過驗證碼進行簽到，</p>
					<p>或利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
					<p><a href="%s"> %s 活動簽到連結</a></p>
					%s
					<p>如有任何問題，請隨時與我們聯繫。</p>
					<p>祝您一切順利！</p>
					<p>此致，</p>
					<p>活動團隊Hilives</p>
			    </body>
			</html>`,
					name[i], activityModel.ActivityName, extPassword[i], url, activityModel.ActivityName, qrcodeMessage)
				message = subject + mime + body

				// 發送郵件
				err = SendMail(email[i], message)
				if err != nil {
					return err
				}
			}
		}

		// 簡訊處理
		if pushPhoneMessage == "open" && messageAmount > 0 {
			var (
				remain     int64
				sendAmount int64
			)
			// 判斷簡訊剩餘數量是否足夠傳給所有用戶
			if messageAmount < int64(dataAmount) { // 匯入資料大於郵件數量
				// 將剩餘郵件數量為0
				remain = 0
				sendAmount = messageAmount
			} else { // 簡訊數量大於匯入資料
				remain = messageAmount - int64(dataAmount)
				sendAmount = int64(dataAmount)
			}

			// 更新簡訊剩餘數量
			err = DefaultActivityModel().
				SetConn(user.DbConn, user.RedisConn, user.MongoConn).
				UpdateMessageAmount(true, model.ActivityID, remain)
			if err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
				// 其他錯誤
				return err
			}

			for i := 0; i < int(sendAmount); i++ {
				// 傳遞訊息(簡訊無法有net的字串連結，因此使用liff url)
				url := fmt.Sprintf(config.HILIVES_APPLYSIGN_URL_LIFF_URL, model.ActivityID, userID[i])

				// XXX 您好: 您已報名XXX活動，該活動的驗證碼為XXX，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): XXX(/applysign?activity_id=xxx&user_id=xxx&sign=open)
				message := fmt.Sprintf("%s 您好: 您已報名 %s 活動，該活動的驗證碼為%s，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): %s",
					name[i], activityModel.ActivityName, extPassword[i], url)

				// 發送簡訊
				err = SendMessage(phone[i], message)
				if err != nil {
					return err
				}

			}

		}

		if isRedis {
			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			user.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "Adds自定義人員")
		}
	}
	return nil
}

// Add 新增單筆用戶
func (user ApplysignUserModel) Add(isRedis bool, model ApplysignUserModel) (string, error) {
	// log.Println("測試單筆新增報名簽到用戶，觀察頭像欄位資料: ")

	var (
		message    string
		userID     = model.ActivityID + "_" + utils.UUID(16)
		password   string
		userfields = []string{
			"user_id",
			"name",
			"avatar",
			"phone",
			"ext_email",
			"identify",
			"activity_id",
			"ext_password",
		}

		applysignfields = []string{
			"user_id",
			"activity_id",
			"status",
			"number",

			"ext_1",
			"ext_2",
			"ext_3",
			"ext_4",
			"ext_5",
			"ext_6",
			"ext_7",
			"ext_8",
			"ext_9",
			"ext_10",
		}
	)

	// 取得自定義欄位資料
	customizeModel, err := DefaultCustomizeModel().
		SetConn(user.DbConn, user.RedisConn, user.MongoConn).
		Find(model.ActivityID)
	if err != nil || customizeModel.ID == 0 {
		return "", errors.New("錯誤: 查詢自定義欄位資料發生問題，請重新操作")
	}

	if (customizeModel.ExtPhoneRequired == "true" || customizeModel.PushPhoneMessage == "open") &&
		model.Phone == "" {
		return "", errors.New("錯誤: 電話欄位為必填，請輸入有效的手機號碼")
	}

	if model.Phone != "" {
		if len(model.Phone) > 2 {
			if !strings.Contains(model.Phone[:2], "09") || len(model.Phone) != 10 {
				return "", errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
			}
		} else {
			return "", errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	}

	if (customizeModel.ExtEmailRequired == "true" || customizeModel.SendMail == "open") &&
		model.ExtEmail == "" {
		return "", errors.New("錯誤: 信箱欄位為必填，請輸入有效的資料「@」，請輸入有效的電子郵件地址")
	}

	if model.ExtEmail != "" && !strings.Contains(model.ExtEmail, "@") {
		return "", errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	var customizepassword string

	if model.Source == "" {
		// 是否自定義輸入驗證碼
		if customizeModel.CustomizePassword == "open" {
			customizepassword = "true"
		} else {
			customizepassword = "false"
		}
	} else if model.Source == "user" {
		customizepassword = "false"
	}

	// 自定義匯入報名簽到人員
	var (
		uniques = []string{
			customizeModel.Ext1Unique, customizeModel.Ext2Unique, customizeModel.Ext3Unique, customizeModel.Ext4Unique,
			customizeModel.Ext5Unique, customizeModel.Ext6Unique, customizeModel.Ext7Unique, customizeModel.Ext8Unique,
			customizeModel.Ext9Unique, customizeModel.Ext10Unique,
			"true", // 驗證碼
		}
		requires = []string{
			customizeModel.Ext1Required, customizeModel.Ext2Required, customizeModel.Ext3Required, customizeModel.Ext4Required,
			customizeModel.Ext5Required, customizeModel.Ext6Required, customizeModel.Ext7Required, customizeModel.Ext8Required,
			customizeModel.Ext9Required, customizeModel.Ext10Required,
			customizepassword, // 驗證碼
		}

		// 判斷自定義欄位是否為唯一值
		values = []string{model.Ext1, model.Ext2, model.Ext3, model.Ext4,
			model.Ext5, model.Ext6, model.Ext7, model.Ext8,
			model.Ext9, model.Ext10,
			model.ExtPassword, // 驗證碼
		}
		exts = make([][]string, 11)
	)

	// 取得該活動所有報名簽到人員資料
	applysigns, err := DefaultApplysignModel().
		SetConn(user.DbConn, user.RedisConn, user.MongoConn).
		FindAll(model.ActivityID, "", "", "", 0, 0)
	if err != nil {
		return "", errors.New("錯誤: 查詢簽到人員資料發生問題，請重新操作")
	}

	// 取得原有報名簽到資料的唯一值
	for _, applysign := range applysigns {
		applysginExts := []string{
			applysign.Ext1, applysign.Ext2, applysign.Ext3, applysign.Ext4,
			applysign.Ext5, applysign.Ext6, applysign.Ext7, applysign.Ext8,
			applysign.Ext9, applysign.Ext10,
			applysign.ExtPassword, // 驗證碼
		}

		for n, ext := range applysginExts {
			// 判斷自定義欄位是否為唯一值，是的話加入陣列中
			if uniques[n] == "true" {
				exts[n] = append(exts[n], ext)
			}
		}
	}

	for n, uniques := range uniques {
		if n == 10 {
			// 驗證碼欄位驗證判斷
			// 是否自定義輸入驗證碼
			if customizepassword == "true" {
				if values[n] == "" {
					return "", errors.New("錯誤: 驗證碼欄位不能為空，請輸入有效資料")
				}

				if utils.InArray(exts[n], values[n]) {
					// 唯一值存在，錯誤
					return "", errors.New("錯誤: 驗證碼欄位不能重複，請輸入有效資料")
				}
				// else {
				// 	exts[n] = append(exts[n], values[n])
				// }

				password = values[n]
			} else if customizepassword == "false" {
				// 不自定義輸入驗證碼，驗證碼隨機產生
				password = utils.RandomNumber(6)

				// 驗證第一次
				if utils.InArray(exts[n], password) {
					// 唯一值存在，錯誤
					// 再隨機產生一次驗證碼
					password = utils.RandomNumber(6)

					// 驗證第二次
					if utils.InArray(exts[n], password) {
						password = utils.RandomNumber(6)
					}

					// 驗證第三次
					if utils.InArray(exts[n], password) {
						password = utils.RandomNumber(6)
					}

					// 驗證第四次
					if utils.InArray(exts[n], password) {
						password = utils.RandomNumber(6)
					}

					// 再驗證一次
					if utils.InArray(exts[n], password) {
						return "", errors.New("錯誤: 隨機產生驗證碼重複，請再匯入一次資料")
					}
				}
			}
		} else {
			// 其它自定義欄位驗證判斷
			if uniques == "true" {
				if values[n] == "" {
					return "", errors.New("錯誤: 唯一值欄位不能為空，請輸入有效資料")
				}

				if utils.InArray(exts[n], values[n]) {
					// 唯一值存在，錯誤
					return "", errors.New("錯誤: 唯一值欄位不能重複，請輸入有效資料")
				}
				// else {
				// 	exts[n] = append(exts[n], values[n])
				// }
			}

			if requires[n] == "true" && values[n] == "" {
				return "", errors.New("錯誤: 必填欄位不能為空，請輸入有效資料")
			}
		}
	}

	var (
		status string
		number int64
	)
	// 判斷簽到審核是否開啟
	if customizeModel.SignCheck == "open" {
		// 開啟，加入資料預設為報名成功
		status = "apply"
		number = 0
	} else if customizeModel.SignCheck == "close" {
		// 關閉，加入資料為簽到成功
		status = "sign"

		// 取得號碼資料
		number, err = DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			GetActivityNumber(true, model.ActivityID)
		if err != nil {
			return "", err
		}

		// 遞增活動參加人數、抽獎號碼
		if err = DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			IncrAttendAndNumber(true, model.ActivityID,
				userID); err != nil {
			// 人數已達上限
			log.Println("參加人數已達上限，改為apply狀態")

			status = "apply"
			number = 0
			// return "", err
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	data["user_id"] = userID
	data["avatar"] = customizeModel.CustomizeDefaultAvatar
	data["identify"] = userID
	data["ext_password"] = password
	data["status"] = status
	data["number"] = number

	// 將資料寫入line_users表中
	_, err = user.Table(config.LINE_USERS_TABLE).Insert(FilterFields(data, userfields))
	if err != nil {
		return "", errors.New("錯誤: 新增用戶發生問題(line_users)，請重新操作")
	}

	// 將資料寫入activity_applysign表中(預設報名成功)
	_, err = user.Table(config.ACTIVITY_APPLYSIGN_TABLE).Insert(FilterFields(data, applysignfields))
	if err != nil {
		return "", errors.New("錯誤: 新增用戶發生問題(activity_applysign)，請重新操作")
	}

	// 郵件
	if model.ExtEmail != "" && customizeModel.SendMail == "open" && customizeModel.MailAmount > 0 {
		url := fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, model.ActivityID, userID, "")
		qrcode := fmt.Sprintf(config.HTTPS_APPLYSIGN_QRCODE_URL, config.HILIVES_NET_URL, fmt.Sprintf("activity_id=%s.user_id=%s", model.ActivityID, userID))
		qrcodeMessage := "<p></p>" // qrcode訊息

		subject := fmt.Sprintf("Subject: %s 活動報名簽到訊息\r\n", customizeModel.ActivityName)
		mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

		// 判斷主持人掃描qrcode開關是否開啟
		if customizeModel.HostScan == "open" {
			qrcodeMessage = fmt.Sprintf(`
			<p>以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到判斷(避免資料洩漏，連結勿提供給他人使用)：</p>
			<img src="https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=%s" alt="簽到QRcode"/>
			`, qrcode)
		}

		body := fmt.Sprintf(`
	    <html>
		 	<body>
				<p>%s 您好，</p>
	    		<p>感謝您報名參加 %s 活動！以下是您的簽到信息：</p>
				<p>該活動的驗證碼為%s，可透過驗證碼進行簽到，</p>
	    		<p>或利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
	    		<p><a href="%s"> %s 活動簽到連結</a></p>
	    		%s
	    		<p>如有任何問題，請隨時與我們聯繫。</p>
	    		<p>祝您一切順利！</p>
	    		<p>此致，</p>
	    		<p>活動團隊Hilives</p>
	        </body>
	    </html>`,
			model.Name, customizeModel.ActivityName, password, url, customizeModel.ActivityName, qrcodeMessage)
		message := subject + mime + body

		// 遞減郵件數量
		err = DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			DecrMailAmount(true, model.ActivityID)
		if err != nil && err.Error() != "錯誤: 數量為0, 不傳遞訊息" {
			// 其他錯誤
			return "", err
		}

		// 無錯誤代表有遞減成功(有錯誤代表可能數量已為0)
		// 郵件發送速度較慢
		if err == nil {
			// 發送郵件
			err = SendMail(model.ExtEmail, message)
			if err != nil {
				return "", err
			}
		}

	}

	// 簡訊
	if model.Phone != "" && customizeModel.PushPhoneMessage == "open" && customizeModel.MessageAmount > 0 {
		// 傳遞訊息(簡訊無法有net的字串連結，因此使用liff url)
		url := fmt.Sprintf(config.HILIVES_APPLYSIGN_URL_LIFF_URL, model.ActivityID, userID)

		// XXX 您好: 您已報名XXX活動，該活動的驗證碼為XXX，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): XXX(/applysign?activity_id=xxx&user_id=xxx&sign=open)
		message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，該活動的驗證碼為%s，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): %s",
			model.Name, customizeModel.ActivityName, password, url)

		// 遞減簡訊數量
		err = DefaultActivityModel().
			SetConn(user.DbConn, user.RedisConn, user.MongoConn).
			DecrMessageAmount(true, model.ActivityID)
		if err != nil && err.Error() != "錯誤: 數量為0, 不傳遞訊息" {
			// 其他錯誤
			return "", err
		}

		// 無錯誤代表有遞減成功(有錯誤代表可能數量已為0)
		if err == nil {
			// 發送簡訊
			err = SendMessage(model.Phone, message)
			if err != nil {
				return "", err
			}
		}
	}

	if isRedis {
		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		user.RedisConn.Publish(config.CHANNEL_SIGN_STAFFS_2_REDIS+model.ActivityID, "Add自定義人員")
	}

	return userID, nil
}

// command.Value{
// 	"user_id":      userID,
// 	"name":         model.Name,
// 	"avatar":       customizeModel.CustomizeDefaultAvatar,
// 	"phone":        model.Phone,
// 	"ext_email":    model.ExtEmail,
// 	"identify":     userID,
// 	"activity_id":  model.ActivityID,
// 	"ext_password": password,

// 	// "avatar":    config.HTTP_HILIVES_NET_URL + "/admin/uploads/system/img-user-pic.png",
// 	// "email":     "",
// 	// "friend":    "no",
// 	// "line":      "晶橙",
// 	// "device":    "customize",
// }

// command.Value{
// 	"user_id":     userID,
// 	"activity_id": model.ActivityID,
// 	"status":      status,
// 	"number":      number,

// 	"ext_1":  model.Ext1,
// 	"ext_2":  model.Ext2,
// 	"ext_3":  model.Ext3,
// 	"ext_4":  model.Ext4,
// 	"ext_5":  model.Ext5,
// 	"ext_6":  model.Ext6,
// 	"ext_7":  model.Ext7,
// 	"ext_8":  model.Ext8,
// 	"ext_9":  model.Ext9,
// 	"ext_10": model.Ext10,
// }

// <html>
// 		<body>
// 		<p>%s 您好，</p>
// 		<p>感謝您報名參加 %s 活動！以下是您的簽到信息：</p>
// 		<p>該活動的驗證碼為%s，可透過驗證碼進行簽到，</p>
// 		<p>或利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
// 		<p><a href="%s"> %s 活動簽到連結</a></p>
// 		<p>以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到判斷(避免資料洩漏，連結勿提供給他人使用)：</p>
// 		<img src="https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=%s" alt="簽到QRcode"/>
// 		<p>如有任何問題，請隨時與我們聯繫。</p>
// 		<p>祝您一切順利！</p>
// 		<p>此致，</p>
// 		<p>活動團隊Hilives</p>
// 	</body>
// </html>

// 取得活動資訊
// activityModel, err := DefaultActivityModel().SetDbConn(user.DbConn).
// 	SetRedisConn(user.RedisConn).Find(false, model.ActivityID)
// if err != nil || activityModel.ID == 0 {
// 	return errors.New("錯誤: 無法取得活動資訊，請重新查詢")
// }

// if model.Source == "excel" {
// 	// 後端平台匯入
// 	// XXX 您好: 您已報名XXX活動，該活動的驗證碼為XXX，可透過驗證碼進行簽到，簽到連結: XXX(/activity/isexist?activity_id=XXX&sign=open&device=customize)
// 	url = fmt.Sprintf(config.HTTPS_ACTIVITY_ISEXIST_URL, config.HILIVES_NET_URL, model.ActivityID) + "&sign=open&device=customize"

// 	// message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，該活動的驗證碼為%s，可透過驗證碼進行簽到，活動簽到連結: %s",
// 	// 	model.Name, customizeModel.ActivityName, password, url)
// } else if model.Source == "user" {
// 	// 用戶填寫表單
// 	url = fmt.Sprintf(config.HTTPS_APPLYSIGN_URL, config.HILIVES_NET_URL, model.ActivityID, userID) + "&sign=open"

// 	// XXX 您好: 您已報名XXX活動，該活動的驗證碼為XXX，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): XXX(/applysign?activity_id=xxx&user_id=xxx&sign=open)
// 	// message = fmt.Sprintf("%s 您好: 您已報名 %s 活動，該活動的驗證碼為%s，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): %s",
// 	// 	model.Name, customizeModel.ActivityName, password, url)
// }

// FindAll 尋找該活動所有簽到人員資料
// func (user ApplysignUserModel) FindAll(activityID string) ([]ApplysignUserModel, error) {

// 	items, err := user.Table(user.Base.TableName).
// 		Where("activity_id", "=", activityID).All()
// 	if err != nil {
// 		return []ApplysignUserModel{}, errors.New("錯誤: 無法取得所有簽到人員資訊，請重新查詢")
// 	}

// 	return MapToApplysignUserModel(items), nil
// }

// MapToApplysignUserModel map轉換[]ApplysignUserModel
// func MapToApplysignUserModel(items []map[string]interface{}) []ApplysignUserModel {
// 	var users = make([]ApplysignUserModel, 0)
// 	for _, item := range items {
// 		var (
// 			user ApplysignUserModel
// 		)

// 		user.ID, _ = item["id"].(int64)
// 		user.UserID, _ = item["user_id"].(string)
// 		user.Name, _ = item["name"].(string)
// 		user.Avatar, _ = item["avatar"].(string)
// 		user.Phone, _ = item["phone"].(string)
// 		user.ExtEmail, _ = item["ext_email"].(string)

// 		user.ActivityID, _ = item["activity_id"].(string)
// 		// user.ExtAccount, _ = item["ext_account"].(string)
// 		user.ExtPassword, _ = item["ext_pass
