package table

import (
	"encoding/json"
	"errors"
	"strings"

	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// GetAttendPanel 參加遊戲人員
func (s *SystemTable) GetAttendPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_STAFF_GAME_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)

			if err := s.table(config.ACTIVITY_STAFF_GAME_TABLE).WhereIn("id", ids).
				Delete(); err != nil {
				return err
			}
			return nil
		})
	return
}

// GetWinningPanel 中獎人員
func (s *SystemTable) GetWinningPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_STAFF_PRIZE_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)
			if err := s.table(config.ACTIVITY_STAFF_PRIZE_TABLE).WhereIn("id", ids).
				Delete(); err != nil {
				return err
			}
			return nil
		})

	formList := table.GetForm()

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("id", "role") {
			return errors.New("錯誤: ID、角色資料發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 拿掉陣列參數(避免json解碼發生問題)
		delete(flattened, "id")

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditPrizeStaffModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.ID = interfaces(strings.Split(values.Get("id"), ","))

		if err := models.DefaultPrizeStaffModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(model); err != nil &&
			err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetBlackPanel 黑名單人員
func (s *SystemTable) GetBlackPanel() (table Table) {
	table = DefaultTable(DefaultConfig())

	formList := table.GetForm()
	formList.SetTable(config.ACTIVITY_STAFF_BLACK_TABLE)
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game", "line_users") {
			return errors.New("錯誤: ID資料發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 拿掉陣列參數(避免json解碼發生問題)
		delete(flattened, "line_users")

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditBlackStaffModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.LineUsers = interfaces(strings.Split(values.Get("line_users"), ","))

		if err := models.DefaultBlackStaffModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game", "line_users") {
			return errors.New("錯誤: ID資料發生問題，請輸入有效的ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 拿掉陣列參數(避免json解碼發生問題)
		delete(flattened, "line_users")

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditBlackStaffModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 手動處理
		model.LineUsers = interfaces(strings.Split(values.Get("line_users"), ","))

		if err := models.DefaultBlackStaffModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, model); err != nil && err.Error() != "錯誤: 無刪除任何資料，請重新操作" {
			return err
		}
		return nil
	})
	return
}

// GetRecordPanel 遊戲紀錄
func (s *SystemTable) GetRecordPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()
	info.SetTable(config.ACTIVITY_SCORE_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)
			if err := s.table(config.ACTIVITY_SCORE_TABLE).WhereIn("id", ids).
				Delete(); err != nil {
				return err
			}
			return nil
		})

	return
}

// @Summary 新增中獎人員資料(form-data)
// @Tags Winning Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param line_user formData string true "LINE ID"
// @param prize_id formData string true "獎品ID"
// @param game formData string true "遊戲" Enums(bingo)
// @param round formData integer true "輪次"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/winning/form [post]
func POSTWinning(ctx *gin.Context) {
}

// @Summary 新增黑名單人員資料(form-data)
// @Tags Black Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string false "遊戲ID"
// @param line_users formData string true "LINE ID(以,區隔多筆LINE ID)"
// @param game formData string true "遊戲" Enums(activity, message, question, lottery, 3DGachaMachine, redpack, ropepack, whack_mole, monopoly, draw_numbers, QA)
// @param reason formData string false "黑名單原因"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/black/form [post]
func POSTBlack(ctx *gin.Context) {
}

// @Summary 移出黑名單人員資料(form-data)
// @Tags Black Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string false "遊戲ID"
// @param line_users formData string true "LINE ID(以,區隔多筆LINE ID)"
// @param game formData string true "遊戲" Enums(activity, message, question, lottery, 3DGachaMachine, redpack, ropepack, whack_mole, monopoly, draw_numbers, QA)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/black/form [put]
func PUTBlack(ctx *gin.Context) {
}

// @Summary 編輯中獎人員資料(form-data)
// @Tags Winning Staff
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param role formData string true "角色" Enums(admin, guest)
// @param password formData string false "密碼" maxlength(8)
// @param status formData string false "兌獎狀態" Enums(yes, no)
// @@@param white formData string false "白名單" Enums(yes, no)
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/winning/form [put]
func PUTWinning(ctx *gin.Context) {
}

// @Summary 編輯聊天紀錄
// @Tags Chatroom_Record
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "聊天紀錄ID"
// @param activity_id formData string true "活動ID"
// @param message_status formData string false "訊息審核狀態" Enums(yes, no, review)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /chatroom/record [put]
func PUTChatroomRecord(ctx *gin.Context) {
}

// @Summary 刪除遊戲人員資料(form-data)
// @Tags Attend Staff
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/attend/form [delete]
func DELETEAttend(ctx *gin.Context) {
}

// @Summary 刪除中獎人員資料(form-data)
// @Tags Winning Staff
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/winning/form [delete]
func DELETEWinning(ctx *gin.Context) {
}

// @@@Summary 刪除PK人員資料(form-data)
// @@@Tags PK Staff
// @@@version 1.0
// @@@Accept  mpfd
// @@@param id formData string true "ID(以,區隔多筆資料ID)"
// @@@param user_id formData string true "用戶ID"
// @@@param token formData string true "CSRF Token"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/pk/form [delete]
func DELETEPK(ctx *gin.Context) {
}

// @Summary 刪除遊戲紀錄資料(form-data)
// @Tags Record Staff
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/record/form [delete]
func DELETERecord(ctx *gin.Context) {
}

// @Summary 參加遊戲人員JSON資料
// @Tags Attend Staff
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(redpack, ropepack, whack_mole, monopoly, QA, tugofwar, 3DGachaMachine)
// @param user_id query string false "用戶ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/attend [get]
func AttendJSON(ctx *gin.Context) {
}

// @Summary 遊戲中獎人員JSON資料
// @Tags Winning Staff
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(redpack, ropepack, whack_mole, monopoly, draw_numbers, lottery, 3DGachaMachine, QA, tugofwar)
// @param user_id query string false "用戶ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param status query string false "兌獎狀態(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/winning [get]
func WinningJSON(ctx *gin.Context) {
}

// @Summary 黑名單人員JSON資料
// @Tags Black Staff
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID"
// @param game query string false "遊戲" Enums(activity, message, question, redpack, ropepack, whack_mole, monopoly, draw_numbers, lottery, 3DGachaMachine, QA, tugofwar)
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/black [get]
func BlackJSON(ctx *gin.Context) {
}

// @Summary 遊戲紀錄JSON資料
// @Tags Record Staff
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(whack_mole, monopoly, QA, tugofwar, vote, bingo)
// @param user_id query string false "用戶ID(用,間隔)"
// @param option_id query string false "選項ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/record [get]
func RecordJSON(ctx *gin.Context) {
}

// @Summary PK人員JSON資料
// @Tags PK Staff
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(bingo)
// @param user_id query string false "用戶ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/pk [get]
// func PKJSON(ctx *gin.Context) {
// }

// formList.SetInsertFunc(func(values form2.Values) error {
// 	if values.IsEmpty("activity_id", "game_id", "line_user", "prize_id", "round", "game") {
// 		return errors.New("錯誤: ID資料發生問題，請輸入有效的ID")
// 	}

// 	// 賓果遊戲發獎
// 	activityID := values.Get("activity_id")
// 	gameID := values.Get("game_id")
// 	game := values.Get("game")
// 	if game == "bingo" {
// 		// 取得賓果遊戲獎品資訊
// 		prizeModel, err := models.DefaultPrizeModel().SetDbConn(s.dbConn).
// 			SetRedisConn(s.redisConn).FindPrizes(true, gameID)
// 		if err != nil || len(prizeModel) < 1 {
// 			return errors.New("錯誤: 無法取得獎品資訊，請重新操作")
// 		}
// 		prize := prizeModel[0]

// 		// 遞減獎品數量(資料表、redis)
// 		err = models.DefaultPrizeModel().SetDbConn(s.dbConn).
// 			SetRedisConn(s.redisConn).DecrRemain(gameID, prize.PrizeID)
// 		if err != nil {
// 			return err
// 		}

// 		// 刪除玩家遊戲紀錄(中獎.未中獎)
// 		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)

// 		// 刪除玩家遊戲紀錄(中獎.未中獎)
// 		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)
// 	}

// 	if _, err := models.DefaultPrizeStaffModel().SetDbConn(s.dbConn).
// 		Add(models.NewPrizeStaffModel{
// 			ActivityID: values.Get("activity_id"),
// 			GameID:     values.Get("game_id"),
// 			UserID:     values.Get("line_user"),
// 			PrizeID:    values.Get("prize_id"),
// 			Round:      values.Get("round"),
// 			Status:     "no",
// 			Score:      0,
// 			Score2:     0,
// 			Rank:       0,
// 		}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}

// 	return nil
// })

// GetPKPanel PK人員
// func (s *SystemTable) GetPKPanel() (table Table) {
// 	table = DefaultTable(DefaultConfig())
// 	info := table.GetInfo()
// 	info.SetTable(config.ACTIVITY_STAFF_PK_TABLE).
// 		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
// 			var ids = interfaces(idArr)
// 			if err := s.table(config.ACTIVITY_STAFF_PK_TABLE).WhereIn("id", ids).
// 				Delete(); err != nil {
// 				return err
// 			}
// 			return nil
// 		})

// 	return
// }

// models.EditPrizeStaffModel{
// 	IDs:      interfaces(strings.Split(values.Get("id"), ",")),
// 	Role:     values.Get("role"),
// 	Password: values.Get("password"),
// 	Status:   values.Get("status"),
// }

// models.NewBlackStaffModel{
// 	ActivityID: values.Get("activity_id"),
// 	GameID:     values.Get("game_id"),
// 	Game:       values.Get("game"),
// 	LineUsers:  interfaces(strings.Split(values.Get("line_users"), ",")),
// 	Reason:     values.Get("reason"),
// }

// models.EditBlackStaffModel{
// 	ActivityID: values.Get("activity_id"),
// 	GameID:     values.Get("game_id"),
// 	Game:       values.Get("game"),
// 	LineUsers:  interfaces(strings.Split(values.Get("line_users"), ",")),
// 	// Reason:     values.Get("reason"),
// }
