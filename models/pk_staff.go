package models

// PKStaffModel 資料表欄位
// type PKStaffModel struct {
// 	Base       `json:"-"`
// 	ID         int64  `json:"id" example:"1"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	GameID     string `json:"game_id" example:"game_id"`
// 	UserID     string `json:"user_id" example:"user_id"`
// 	Round      int64  `json:"round" example:"1"`

// 	// 用戶資訊
// 	Name     string `json:"name"`
// 	Avatar   string `json:"avatar"`
// 	Phone    string `json:"phone"`
// 	Email    string `json:"email"`
// 	ExtEmail string `json:"ext_email"`

// 	// 遊戲資訊
// 	Title string `json:"title" example:"Game Title"`
// 	Game  string `json:"game" example:"redpack、ropepack、whack_mole、lottery"`
// }

// // NewPKStaffModel 資料表欄位
// type NewPKStaffModel struct {
// 	UserID     string `json:"user_id" example:"user_id"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	GameID     string `json:"game_id" example:"game_id"`
// 	Round      int64  `json:"round" example:"1"`
// }

// // DefaultPKStaffModel 預設PKStaffModel
// func DefaultPKStaffModel() PKStaffModel {
// 	return PKStaffModel{Base: Base{TableName: config.ACTIVITY_STAFF_PK_TABLE}}
// }

// // SetDbConn 設定connection
// func (m PKStaffModel) SetDbConn(conn db.Connection) PKStaffModel {
// 	m.DbConn = conn
// 	return m
// }

// // SetRedisConn 設定connection
// func (m PKStaffModel) SetRedisConn(conn cache.Connection) PKStaffModel {
// 	m.RedisConn = conn
// 	return m
// }

// // FindAll 查詢所有PK紀錄
// func (m PKStaffModel) FindAll(activityID, gameID, userID, game, round string, limit, offset int64) ([]PKStaffModel, error) {
// 	var (
// 		sql = m.Table(m.Base.TableName).
// 			Select("activity_staff_pk.id", "activity_staff_pk.activity_id",
// 				"activity_staff_pk.game_id",
// 				"activity_staff_pk.user_id", "activity_staff_pk.round",

// 				// 遊戲場次
// 				"activity_game.title", "activity_game.game",

// 				// 用戶
// 				"line_users.name", "line_users.avatar", "line_users.phone",
// 				"line_users.email", "line_users.ext_email",
// 			).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_pk.user_id",
// 				FieldA1:   "line_users.user_id",
// 				Table:     "line_users",
// 				Operation: "="}).
// 			LeftJoin(command.Join{
// 				FieldA:    "activity_staff_pk.game_id",
// 				FieldA1:   "activity_game.game_id",
// 				Table:     "activity_game",
// 				Operation: "="}).
// 			Where("activity_staff_pk.activity_id", "=", activityID).
// 			OrderBy("activity_game.game", "asc",
// 				"activity_staff_pk.game_id", "asc",
// 			)
// 	)

// 	// 判斷參數是否為空
// 	if gameID != "" {
// 		sql = sql.WhereIn("activity_staff_pk.game_id", interfaces(strings.Split(gameID, ",")))
// 	}
// 	if userID != "" {
// 		sql = sql.WhereIn("activity_staff_pk.user_id", interfaces(strings.Split(userID, ",")))
// 	}
// 	if game != "" {
// 		sql = sql.WhereIn("activity_game.game", interfaces(strings.Split(game, ",")))
// 	}
// 	if round != "" {
// 		sql = sql.WhereIn("activity_staff_pk.round", interfaces(strings.Split(round, ",")))
// 	}

// 	if limit != 0 {
// 		sql = sql.Limit(limit)
// 	}
// 	if offset != 0 {
// 		sql = sql.Offset(offset)
// 	}

// 	items, err := sql.All()
// 	if err != nil {
// 		return nil, errors.New("錯誤: 無法取得PK紀錄，請重新查詢")
// 	}

// 	return MapToPKStaffModel(items), nil
// }

// // Add 增加資料
// func (m PKStaffModel) Add(model NewPKStaffModel) error {
// 	_, err := m.Table(m.TableName).Insert(command.Value{
// 		"user_id":     model.UserID,
// 		"activity_id": model.ActivityID,
// 		"game_id":     model.GameID,
// 		"round":       model.Round,
// 	})
// 	return err
// }

// // MapToPKStaffModel map轉換[]PKStaffModel
// func MapToPKStaffModel(items []map[string]interface{}) []PKStaffModel {
// 	var staffs = make([]PKStaffModel, 0)
// 	for _, item := range items {
// 		var (
// 			staff PKStaffModel
// 		)
// 		staff.ID, _ = item["id"].(int64)
// 		staff.ActivityID, _ = item["activity_id"].(string)
// 		staff.GameID, _ = item["game_id"].(string)
// 		staff.UserID, _ = item["user_id"].(string)
// 		staff.Round, _ = item["round"].(int64)

// 		// 用戶資訊
// 		staff.Name, _ = item["name"].(string)
// 		staff.Avatar, _ = item["avatar"].(string)
// 		staff.Phone, _ = item["phone"].(string)
// 		staff.Email, _ = item["email"].(string)
// 		staff.ExtEmail, _ = item["ext_email"].(string)

// 		// 遊戲資訊
// 		staff.Title, _ = item["title"].(string)
// 		staff.Game, _ = item["game"].(string)

// 		staffs = append(staffs, staff)
// 	}
// 	return staffs
// }
