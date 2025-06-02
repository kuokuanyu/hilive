package controller

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ShowSelect 選擇活動管理或聊天室 GET API
func (h *Handler) ShowSelect(ctx *gin.Context) {
	var (
		user       = h.GetLoginUser(ctx.Request, "chatroom_session")
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		params     = url.Values{}
		// *****可以同時登入(暫時拿除)*****
		// ip = utils.ClientIP(ctx.Request)
		// *****可以同時登入(暫時拿除)*****
		userModel models.UserModel
		err       error
		// userID          = ctx.Query("user_id")
		// userModel, err = models.DefaultUserModel().SetDbConn(h.dbConn).Find("user_id", user.UserID)
		// lineModel, err2 = models.DefaultLineModel().SetDbConn(h.dbConn).Find("identify", user.UserID)
		// activityModel, err3 = models.DefaultActivityModel().SetDbConn(h.dbConn).Find("activity_id", activityID)
		// client              = &http.Client{}
		// data                = url.Values{}
		// result              interface{}
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法活動資訊，請輸入有效的活動ID")
		return
	}

	// 先判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
	// dataMap, err := h.redisConn.HashGetAllCache(config.HILIVE_USERS + user.Identify)
	// if err != nil {
	// 	h.executeErrorHTML(ctx, "錯誤: 取得用戶快取資料發生問題")
	// 	return
	// }
	// userModel.UserID, _ = dataMap["user_id"]
	// userModel.Name, _ = dataMap["name"]
	// userModel.Avatar, _ = dataMap["avatar"]
	// userModel.Phone, _ = dataMap["phone"]
	// userModel.Email, _ = dataMap["email"]
	// userModel.Bind, _ = dataMap["bind"]
	// userModel.Cookie, _ = dataMap["cookie"]
	// userModel.Ip, _ = dataMap["ip"]

	// redis中沒有用戶資料，查詢資料表後並加入redis中
	// if userModel.UserID == "" {
	// *****可以同時登入(暫時拿除)*****
	// userModel, err = models.DefaultUserModel().SetDbConn(h.dbConn).
	// 	SetRedisConn(h.redisConn).Find(true, true, ip, "users.user_id", user.Identify)
	// *****可以同時登入(暫時拿除)*****

	// *****可以同時登入(新)*****
	userModel, err = models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, true, "users.user_id", user.Identify)
	// *****可以同時登入(新)*****

	if err != nil || userModel.UserID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
		return
		// }

		// 將用戶資訊加入redis中
		// values := []interface{}{config.HILIVE_USERS + user.Identify}
		// values = append(values, "user_id", userModel.UserID)
		// values = append(values, "name", userModel.Name)
		// values = append(values, "phone", userModel.Phone)
		// values = append(values, "email", userModel.Email)
		// values = append(values, "avatar", userModel.Avatar)
		// values = append(values, "bind", userModel.Bind)
		// values = append(values, "cookie", userModel.Cookie)
		// values = append(values, "ip", userModel.Ip)
		// params = append(params, "table", "users")
		// if err := h.redisConn.HashMultiSetCache(values); err != nil {
		// 	h.executeErrorHTML(ctx, "錯誤: 設置用戶快取資料發生問題")
		// 	return
		// }
		// 設置過期時間
		// if err := h.redisConn.SetExpire(config.HILIVE_USERS+user.Identify,
		// 	config.REDIS_EXPIRE); err != nil {
		// 	h.executeErrorHTML(ctx, "錯誤: 設置快取過期時間發生問題")
		// 	return
		// }
	}

	// *****可以同時登入(暫時拿除)*****
	// 更新用戶ip資訊
	// if err := models.DefaultUserModel().SetDbConn(h.dbConn).
	// 	UpdateIP(userModel.UserID, ip); err != nil {
	// 	h.executeErrorHTML(ctx, err.Error())
	// 	return
	// }
	// *****可以同時登入(暫時拿除)*****

	// 設置hilive_session cookie(session頁面、主持端聊天室頁面也有cookie的設置)
	params.Add("user_id", userModel.UserID)
	// params.Add("name", userModel.Name)
	// params.Add("avatar", userModel.Avatar)
	// params.Add("phone", userModel.Phone)
	// params.Add("email", userModel.Email)
	// params.Add("bind", userModel.Bind)
	// params.Add("cookie", userModel.Cookie)
	// *****可以同時登入(確定拿除)*****
	// params.Add("ip", ip)
	// *****可以同時登入(確定拿除)*****
	params.Add("table", "users")
	params.Add("sign", utils.UserSign(userModel.UserID, config.HiliveCookieSecret))
	ctx.SetSameSite(4)
	ctx.SetCookie("hilive_session", string(utils.Encode([]byte(params.Encode()))),
		config.GetSessionLifeTime(), "/", "", true, false)

	h.executeHTML(ctx, "./hilives/hilive/views/cms/choose_page.html", executeParam{
		Route: route{
			Prefix: h.config.Prefix,
		},
		User:       user,
		ActivityID: activityID,
	})
}

// 將用戶資訊加入redis(目前只儲存ip資料，用於middleware判斷)
// if err = h.redisConn.SetCache(config.HILIVE_USERS+userModel.UserID, ip,
// 	"EX", config.GetSessionLifeTime()); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// if err = h.redisConn.HashSetCache(
// 	"users", userModel.UserID, models.LoginUser{
// 		UserID:   userModel.UserID,
// 		Name:     userModel.Name,
// 		Phone:    userModel.Phone,
// 		Email:    userModel.Email,
// 		Avatar:   userModel.Avatar,
// 		Bind:     userModel.Bind,
// 		Cookie:   userModel.Cookie,
// 		Identify: "",
// 		Friend:   "",
// 		Table:    "users",
// 		Ip:       ip,
// 	}); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// if err := auth.SetCookie(ctx, user, h.dbConn, "hilive_session"); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: cookie發生問題，請重新整理頁面")
// 	return
// }
// if err := auth.SetCookie(ctx, models.UserModel{UserID: lineModel.UserID},
// 	h.dbConn, "chatroom_session"); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: cookie發生問題，請重新整理頁面")
// 	return
// }

// data.Set("phone", user.Phone)
// data.Set("password", user.Password)
// req, err := http.NewRequest("POST", fmt.Sprintf(LOGIN_URL, host), strings.NewReader(data.Encode()))
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 無法登入")
// 	return
// }
// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// resp, err := client.Do(req)
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 無法登入")
// 	return
// }
// defer resp.Body.Close()
// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	response.Error(ctx, "發生錯誤: 讀取內容發生錯誤")
// 	return
// }
// json.Unmarshal(body, &result)
// if result.(map[string]interface{})["code"].(float64) == 400 {
// 	response.Error(ctx, result.(map[string]interface{})["msg"].(string))
// 	return
// }
// for _, cookie := range resp.Cookies() {
// 	ctx.SetCookie(cookie.Name, cookie.Value, config.GetSessionLifeTime(), "/", "", true, false)
// }
