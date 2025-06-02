package auth

import (
	"hilive/modules/config"
	"hilive/modules/utils"
)

func AddToken(userID string) string {
	return utils.UserSign(userID, config.TokenSecret)
}

// CheckToken 驗證token
func CheckToken(CheckToken, userID string) bool {
	var (
		token = utils.UserSign(userID, config.TokenSecret)
	)

	if token != CheckToken {
		return false
	}
	return true
}

// *****舊*****
// AddToken 新增Token
// func AddToken(ip, userID string) string {
// 	return utils.UserSign(ip, userID, config.TokenSecret)
// }

// CheckToken 驗證token
// func CheckToken(CheckToken, ip, userID string) bool {
// 	var (
// 		token = utils.UserSign(ip, userID, config.TokenSecret)
// 	)

// 	if token != CheckToken {
// 		return false
// 	}
// 	return true
// }
// *****舊*****

// CSRFToken token list
// type CSRFToken []string

// TokenService struct
// type TokenService struct {
// 	Tokens CSRFToken
// 	// lock   sync.Mutex
// }

// func init() {
// 	service.Register("token_csrf", func() (service.Service, error) {
// 		return &TokenService{
// 			Tokens: make(CSRFToken, 0),
// 		}, nil
// 	})
// }

// GetTokenService 取得TokenService
// func GetTokenService(s interface{}) *TokenService {
// 	if srv, ok := s.(*TokenService); ok {
// 		return srv
// 	}
// 	panic("錯誤的Service")
// }

// Name TokenService方法
// func (s *TokenService) Name() string {
// 	return "token_csrf"
// }

// AddToken 新增Token
// func (s *TokenService) AddToken(conn db.Connection, user models.UserModel) string {
// 	var (
// 		tokenModel map[string]interface{}
// 		uuid       string
// 		err        error
// 	)
// 	// s.lock.Lock()
// 	// defer s.lock.Unlock()

// 	// 檢查是否有token資料
// 	if tokenModel, err = db.Table(config.TOKEN_TABLE).WithConn(conn).
// 		Where("user_id", "=", user.ID).First(); err != nil {
// 		return ""
// 	}
// 	if tokenModel == nil {
// 		uuid = utils.UUID(32)
// 		if _, err := db.Table(config.TOKEN_TABLE).WithConn(conn).Insert(command.Value{
// 			"token_values": uuid, "user_id": user.ID}); err != nil {
// 			log.Fatal("錯誤: 新增csrf token資料發生問題，請重新操作")
// 		}
// 	} else {
// 		uuid = fmt.Sprintf("%s", tokenModel["token_values"])
// 		// if err = db.Table("token").WithConn(conn).Where("user_id", "=", user.ID).Update(
// 		// 	command.Value{"updated_at": time.Now().Add(time.Hour * time.Duration(8))}); err != nil &&
// 		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		// 	log.Fatal("更新csrf token資料發生錯誤")
// 		// }
// 	}

// 	// s.Tokens = append(s.Tokens, uuid)
// 	return uuid
// }

// // CheckToken 檢查是否存在token
// func (s *TokenService) CheckToken(conn db.Connection, CheckToken string) bool {
// 	list, err := db.Table(config.TOKEN_TABLE).WithConn(conn).All()
// 	if err != nil {
// 		log.Fatal("錯誤: 取得csrf token資料發生問題，請重新操作")
// 	}

// 	for i := 0; i < len(list); i++ {
// 		if list[i]["token_values"].(string) == CheckToken {
// 			// s.Tokens = append((s.Tokens)[:i], (s.Tokens)[i+1:]...)
// 			return true
// 		}
// 	}
// 	return false
// }

// deleteOverdueToken 刪除資料表中超過時效的token資料
// func (driver *DBDriver) deleteOverdueToken() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	if config.GetDatabases().GetHilive().Driver == "mysql" {
// 		db.Table(config.TOKEN_TABLE).WithConn(driver.conn).
// 			WhereRaw(`unix_timestamp(updated_at) < unix_timestamp() - ` +
// 				strconv.Itoa(config.GetSessionLifeTime())).Delete()
// 	}
// }
