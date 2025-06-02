package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// SetSession 設置session(hilive_session、chatroom_session)
func (h *Handler) SetSession(ctx *gin.Context) {
	var (
		host        = ctx.Request.Host
		activityID  = ctx.Query("activity_id")
		gameID      = ctx.Query("game_id")
		userID      = ctx.Query("user_id")
		sessionName = ctx.Query("session_name")
		// role        = ctx.Query("role")
		ip          = utils.ClientIP(ctx.Request)
		params      = url.Values{} // cookie參數
		redirectURL string

		// sign        = ctx.Query("sign") // 報名or簽到判斷
		// userModel   models.UserModel
		// liffURL, redirect string
		// err      error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if sessionName == "hilive" {
		redirectURL = fmt.Sprintf(config.ACTIVITY_URL, "true") // 活動頁面

		userModel, err := models.DefaultUserModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, true, "users.user_id", userID)
		if err != nil || userModel.UserID == "" {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
			return
		}

		// 用戶資訊(cookie資料)
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

		if len(userModel.Permissions) == 0 {
			h.executeErrorHTML(ctx, "錯誤: 用戶無瀏覽網頁權限，請聯絡管理員")
			return
		}

		// 權限菜單判斷
		// for _, permission := range userModel.Permissions {
		// 	if redirect == fmt.Sprintf(config.ACTIVITY_URL, "true") {
		// 		break
		// 	}

		// 	for _, path := range permission.HTTPPath{
		// 		if path == "*" || path == "admin" || path == "/admin/activity" {
		// 			redirect = fmt.Sprintf(config.ACTIVITY_URL, "true")// 活動頁面
		// 		}
		// 	}
		// }
	} else if sessionName == "chatroom" {
		lineModel, err := models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, ip, "user_id", userID)
		if err != nil || lineModel.UserID == "" {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
			return
		}

		// 用戶資訊(cookie資料)
		params.Add("activity_id", activityID) // 活動ID，判斷用
		params.Add("user_id", lineModel.UserID)
		// params.Add("name", lineModel.Name)
		// params.Add("avatar", lineModel.Avatar)
		// params.Add("phone", lineModel.Phone)
		// params.Add("email", lineModel.Email)
		// params.Add("identify", lineModel.Identify)
		// params.Add("friend", lineModel.Friend)
		// *****可以同時登入(確定拿除)*****
		// params.Add("ip", ip)
		// *****可以同時登入(確定拿除)*****
		params.Add("table", "line_users")
		params.Add("sign", utils.UserSign(lineModel.UserID, config.ChatroomCookieSecret))

		// if role == "admin" {
		// redirectURL = fmt.Sprintf(config.SELECT_URL, activityID) // 導向選擇頁面(選擇進入後端平台或玩家端聊天室)
		// } else if role == "guest" {
		if gameID == "" {
			// 遊戲ID為空直接導向聊天室
			redirectURL = fmt.Sprintf(config.GUEST_URL, activityID) // 聊天室
		} else if gameID != "" {
			var (
				game string
				err  error
			)

			if gameID != "signname" && gameID != "question" {
				// 遊戲ID不為空，查詢遊戲類型並導向至遊戲頁面
				// 遊戲ID不為signname.question
				game, err = models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindGameType(true, gameID) // 取得遊戲類型
			}

			// log.Println("session取得遊戲種類: ", game)
			if err == nil && game != "" {
				// 遊戲查詢正確
				if game == "redpack" {
					redirectURL = fmt.Sprintf(config.HTTPS_REDPACK_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "ropepack" {
					redirectURL = fmt.Sprintf(config.HTTPS_ROPEPACK_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "whack_mole" {
					redirectURL = fmt.Sprintf(config.HTTPS_WHACK_MOLE_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "lottery" {
					redirectURL = fmt.Sprintf(config.HTTPS_LOTTERY_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "monopoly" {
					redirectURL = fmt.Sprintf(config.HTTPS_MONOPOLY_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "QA" {
					redirectURL = fmt.Sprintf(config.HTTPS_QA_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "tugofwar" {
					redirectURL = fmt.Sprintf(config.HTTPS_TUGOFWAR_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "bingo" {
					redirectURL = fmt.Sprintf(config.HTTPS_BINGO_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				} else if game == "3DGachaMachine" {
					redirectURL = fmt.Sprintf(config.HTTPS_GACHAMACHINE_GAME_URL, config.HILIVES_NET_URL, activityID, "guest") + gameID
				} else if game == "vote" {
					redirectURL = fmt.Sprintf(config.HTTPS_VOTE_GAME_URL, config.HILIVES_NET_URL, activityID, "guest") + gameID
				} else if game == "draw_numbers" {
					redirectURL = fmt.Sprintf(config.HTTPS_DRAW_NUMBERS_GAME_URL, config.HILIVES_NET_URL, activityID) + gameID
				}
			} else if gameID == "signname" {
				// log.Println("自動導向玩家端簽名牆")
				redirectURL = fmt.Sprintf(config.HTTPS_SIGNNAME_URL, config.HILIVES_NET_URL, activityID)
			} else if gameID == "question" {
				// log.Println("自動導向玩家端提問牆")
				// redirectURL = fmt.Sprintf(config.HTTPS_QUESTION_URL, config.HILIVES_NET_URL, activityID)

			} else {
				// 遊戲查詢錯誤直接導向聊天室
				redirectURL = fmt.Sprintf(config.GUEST_URL, activityID) // 聊天室
			}
			// }

			// log.Println("玩家簽到完成，將玩家遊戲紀錄資料寫入redis中")
			// 用戶遊戲紀錄，簽到完成後先寫入reids中，避免玩家進入遊戲時需併發處理大量玩家資料
			models.DefaultPrizeStaffModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindUserGameRecords(true, activityID, userID)

		}
		// if sign == "open" {
		// 	redirect += "&sign=open"
		// }
	}

	// 設置cookie(主持端聊天室頁面、select頁面也有cookie的設置)
	ctx.SetSameSite(4)
	ctx.SetCookie(sessionName+"_session", string(utils.Encode([]byte(params.Encode()))),
		config.GetSessionLifeTime(), "/", "", true, false)

	ctx.Redirect(http.StatusFound, redirectURL)
}

// var userModel models.UserModel
// 先判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
// dataMap, err := h.redisConn.HashGetAllCache(config.HILIVE_USERS + userID)
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

// redis沒有用戶資料，查詢資料表
// if userModel.UserID == "" {

// 將用戶資訊加入redis
// values := []interface{}{config.HILIVE_USERS + userModel.UserID}
// values = append(values, "user_id", userModel.UserID)
// values = append(values, "name", userModel.Name)
// values = append(values, "phone", userModel.Phone)
// values = append(values, "email", userModel.Email)
// values = append(values, "avatar", userModel.Avatar)
// values = append(values, "bind", userModel.Bind)
// values = append(values, "cookie", userModel.Cookie)
// values = append(values, "ip", ip)
// params = append(params, "table", "users")
// if err := h.redisConn.HashMultiSetCache(values); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 設置用戶快取資料發生問題")
// 	return
// }
// 設置過期時間
// if err := h.redisConn.SetExpire(config.HILIVE_USERS+userModel.UserID,
// 	strconv.Itoa(config.GetSessionLifeTime())); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 設置快取過期時間發生問題")
// 	return
// }
// }

// var lineModel models.LineModel
// 先判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
// dataMap, err := h.redisConn.HashGetAllCache(config.LINE_USERS + userID)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 取得用戶快取資料發生問題")
// 	return
// }
// lineModel.UserID, _ = dataMap["user_id"]
// lineModel.Name, _ = dataMap["name"]
// lineModel.Avatar, _ = dataMap["avatar"]
// lineModel.Phone, _ = dataMap["phone"]
// lineModel.Email, _ = dataMap["email"]
// lineModel.Identify, _ = dataMap["identify"]
// lineModel.Friend, _ = dataMap["friend"]
// lineModel.Ip, _ = dataMap["ip"]

// redis沒有用戶資料，查詢資料表
// if lineModel.UserID == "" {

// 將用戶資訊加入redis
// values := []interface{}{config.LINE_USERS + lineModel.UserID}
// values = append(values, "user_id", lineModel.UserID)
// values = append(values, "name", lineModel.Name)
// values = append(values, "phone", lineModel.Phone)
// values = append(values, "email", lineModel.Email)
// values = append(values, "avatar", lineModel.Avatar)
// values = append(values, "identify", lineModel.Identify)
// values = append(values, "friend", lineModel.Friend)
// values = append(values, "ip", ip)
// if err := h.redisConn.HashMultiSetCache(values); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 設置用戶快取資料發生問題")
// 	return
// }
// 設置過期時間
// if err := h.redisConn.SetExpire(config.LINE_USERS+lineModel.UserID,
// 	strconv.Itoa(config.GetSessionLifeTime())); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 設置快取過期時間發生問題")
// 	return
// }
// }

// if err = auth.SetCookie(ctx, userModel, h.dbConn, sessionName+"_session"); err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: cookie發生問題，請重新操作")
// 	return
// }

// // 更新用戶ip資訊
// if err := models.DefaultLineModel().SetDbConn(h.dbConn).
// 	UpdateIP(userID, ip); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }
// 將用戶資訊加入redis(目前只儲存ip資料，用於middleware判斷)
// if err = h.redisConn.SetCache(config.LINE_USERS+lineModel.UserID, ip,
// 	"EX", config.GetSessionLifeTime()); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// if err = h.redisConn.HashSetCache(
// 	"line_users", lineModel.UserID, models.LoginUser{
// 		UserID:   lineModel.UserID,
// 		Name:     lineModel.Name,
// 		Phone:    lineModel.Phone,
// 		Email:    lineModel.Email,
// 		Avatar:   lineModel.Avatar,
// 		Bind:     "",
// 		Cookie:   "",
// 		Identify: lineModel.Identify,
// 		Friend:   lineModel.Friend,
// 		Table:    "line_users",
// 		Ip:       ip,
// 	}); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// 將用戶資訊加入redis(目前只儲存ip資料，用於middleware判斷)
// if err = h.redisConn.SetCache(config.HILIVE_USERS+userModel.UserID, ip,
// 	"EX", config.GetSessionLifeTime()); err != nil {
// 	h.executeErrorHTML(ctx, err.Error())
// 	return
// }

// if err = h.redisConn.HashSetCache(
// 	"hilive_users", userModel.UserID, models.LoginUser{
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
//
