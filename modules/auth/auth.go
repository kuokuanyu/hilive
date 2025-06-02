package auth

var (
// userValue = make(map[string]interface{}) // UserValue 登入用戶資訊
// lock sync.Mutex
)

// Auth 取得目前登入用戶資訊
// func Auth() models.UserModel {
// 	return userValue["user"].(models.UserModel)
// }

// // CheckUser 檢查手機與密碼是否符合
// func CheckUser(password string, phone string, conn db.Connection) (user models.UserModel, ok bool) {
// 	if user, _ = models.DefaultUserModel().SetConn(conn).Find("phone", phone); user.ID == int64(0) {
// 		ok = false
// 		return
// 	}
// 	// if comparePassword(password, user.Password) || password == user.Password {
// 	if comparePassword(password, user.Password) {
// 		ok = true
// 		// TODO: 目前沒有權限相關問題，先拿掉角色權限菜單判斷
// 		// user = user.GetRoles().GetPermissions().GetMenus()
// 		return
// 	}
// 	ok = false
// 	return
// }

// SetCookie 設置cookie並加入header
// func SetCookie(ctx *gin.Context, user models.UserModel, conn db.Connection, sessionName string) error {
// 	var (
// 		cookie  *Session
// 		isExist bool
// 		err     error
// 	)
// 	if isExist, cookie, err = InitSession(ctx, conn, sessionName); err != nil {
// 		return err
// 	}

// 	// 判斷session名稱
// 	if sessionName == "hilive_session" {
// 		cookie.SessionValues["hilive"] = user.UserID
// 	} else if sessionName == "chatroom_session" {
// 		cookie.SessionValues["chatroom"] = user.UserID
// 	}

// 	if err = cookie.Driver.Update(isExist, cookie.SessionID, cookie.SessionValues); err != nil {
// 		return err
// 	}

// 	cookie.Context.SetSameSite(4)
// 	cookie.Context.SetCookie(sessionName, cookie.SessionID, config.GetSessionLifeTime(), "/", "", true, false)
// 	return nil
// }

// DeleteCookie delete the cookie
// func DeleteCookie(ctx *gin.Context, conn db.Connection, sessionName string) error {
// 	_, cookie, err := InitSession(ctx, conn, sessionName)
// 	if err != nil {
// 		return err
// 	}
// 	return cookie.Clear()
// }

// Clear 清除Session
// func (cookie *Session) Clear() error {
// 	cookie.SessionValues = map[string]interface{}{}
// 	return cookie.Driver.Update(false, cookie.SessionID, cookie.SessionValues)
// }

// // comparePassword 檢查密碼是否相符
// func comparePassword(comPwd, pwdHash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
// 	return err == nil
// }

// EncodePassword 加密
// func EncodePassword(pwd []byte) string {
// 	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(hash)
// }
