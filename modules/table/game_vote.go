package table

import (
	"encoding/json"
	"errors"
	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	votePictureFields = []PictureField{
		{FieldName: "vote_classic_h_pic_01", Path: "vote/classic/vote_classic_h_pic_01.png"},
		{FieldName: "vote_classic_h_pic_02", Path: "vote/classic/vote_classic_h_pic_02.png"},
		{FieldName: "vote_classic_h_pic_03", Path: "vote/classic/vote_classic_h_pic_03.png"},
		{FieldName: "vote_classic_h_pic_04", Path: "vote/classic/vote_classic_h_pic_04.png"},
		{FieldName: "vote_classic_h_pic_05", Path: "vote/classic/vote_classic_h_pic_05.png"},
		{FieldName: "vote_classic_h_pic_06", Path: "vote/classic/vote_classic_h_pic_06.png"},
		{FieldName: "vote_classic_h_pic_07", Path: "vote/classic/vote_classic_h_pic_07.png"},
		{FieldName: "vote_classic_h_pic_08", Path: "vote/classic/vote_classic_h_pic_08.jpg"},
		{FieldName: "vote_classic_h_pic_09", Path: "vote/classic/vote_classic_h_pic_09.png"},
		{FieldName: "vote_classic_h_pic_10", Path: "vote/classic/vote_classic_h_pic_10.png"},
		{FieldName: "vote_classic_h_pic_11", Path: "vote/classic/vote_classic_h_pic_11.png"},
		{FieldName: "vote_classic_h_pic_13", Path: "vote/classic/vote_classic_h_pic_13.png"},
		{FieldName: "vote_classic_h_pic_14", Path: "vote/classic/vote_classic_h_pic_14.png"},
		{FieldName: "vote_classic_h_pic_15", Path: "vote/classic/vote_classic_h_pic_15.png"},
		{FieldName: "vote_classic_h_pic_16", Path: "vote/classic/vote_classic_h_pic_16.png"},
		{FieldName: "vote_classic_h_pic_17", Path: "vote/classic/vote_classic_h_pic_17.png"},
		{FieldName: "vote_classic_h_pic_18", Path: "vote/classic/vote_classic_h_pic_18.png"},
		{FieldName: "vote_classic_h_pic_19", Path: "vote/classic/vote_classic_h_pic_19.png"},
		{FieldName: "vote_classic_h_pic_20", Path: "vote/classic/vote_classic_h_pic_20.png"},
		{FieldName: "vote_classic_h_pic_21", Path: "vote/classic/vote_classic_h_pic_21.png"},
		{FieldName: "vote_classic_h_pic_23", Path: "vote/classic/vote_classic_h_pic_23.png"},
		{FieldName: "vote_classic_h_pic_24", Path: "vote/classic/vote_classic_h_pic_24.png"},
		{FieldName: "vote_classic_h_pic_25", Path: "vote/classic/vote_classic_h_pic_25.png"},
		{FieldName: "vote_classic_h_pic_26", Path: "vote/classic/vote_classic_h_pic_26.png"},
		{FieldName: "vote_classic_h_pic_27", Path: "vote/classic/vote_classic_h_pic_27.png"},
		{FieldName: "vote_classic_h_pic_28", Path: "vote/classic/vote_classic_h_pic_28.png"},
		{FieldName: "vote_classic_h_pic_29", Path: "vote/classic/vote_classic_h_pic_29.png"},
		{FieldName: "vote_classic_h_pic_30", Path: "vote/classic/vote_classic_h_pic_30.png"},
		{FieldName: "vote_classic_h_pic_31", Path: "vote/classic/vote_classic_h_pic_31.png"},
		{FieldName: "vote_classic_h_pic_32", Path: "vote/classic/vote_classic_h_pic_32.png"},
		{FieldName: "vote_classic_h_pic_33", Path: "vote/classic/vote_classic_h_pic_33.png"},
		{FieldName: "vote_classic_h_pic_34", Path: "vote/classic/vote_classic_h_pic_34.png"},
		{FieldName: "vote_classic_h_pic_35", Path: "vote/classic/vote_classic_h_pic_35.png"},
		{FieldName: "vote_classic_h_pic_36", Path: "vote/classic/vote_classic_h_pic_36.png"},
		{FieldName: "vote_classic_h_pic_37", Path: "vote/classic/vote_classic_h_pic_37.png"},
		{FieldName: "vote_classic_g_pic_01", Path: "vote/classic/vote_classic_g_pic_01.png"},
		{FieldName: "vote_classic_g_pic_02", Path: "vote/classic/vote_classic_g_pic_02.png"},
		{FieldName: "vote_classic_g_pic_03", Path: "vote/classic/vote_classic_g_pic_03.png"},
		{FieldName: "vote_classic_g_pic_04", Path: "vote/classic/vote_classic_g_pic_04.jpg"},
		{FieldName: "vote_classic_g_pic_05", Path: "vote/classic/vote_classic_g_pic_05.png"},
		{FieldName: "vote_classic_g_pic_06", Path: "vote/classic/vote_classic_g_pic_06.png"},
		{FieldName: "vote_classic_g_pic_07", Path: "vote/classic/vote_classic_g_pic_07.png"},
		{FieldName: "vote_classic_c_pic_01", Path: "vote/classic/vote_classic_c_pic_01.png"},
		{FieldName: "vote_classic_c_pic_02", Path: "vote/classic/vote_classic_c_pic_02.png"},
		{FieldName: "vote_classic_c_pic_03", Path: "vote/classic/vote_classic_c_pic_03.png"},
		{FieldName: "vote_classic_c_pic_04", Path: "vote/classic/vote_classic_c_pic_04.png"},

		{FieldName: "vote_bgm_gaming", Path: "vote/%s/bgm/gaming.mp3"},
	}
)

// GetVoteOptionListsPanel 投票選項名單人員(匹量匯入)
func (s *SystemTable) GetVoteOptionListsPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "option_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteOptionListModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(models.EditGameVoteOptionListModel{
				ActivityID: values.Get("activity_id"),
				GameID:     values.Get("game_id"),
				UserID:     values.Get("option_list_id"),
				OptionID:   values.Get("option_id"),
				Name:       values.Get("name"),
				Leader:     values.Get("leader"),

				People: values.Get("people"),
			}); err != nil {
			return err
		}
		return nil
	})

	return
}

// GetVoteOptionListPanel 投票選項名單人員
func (s *SystemTable) GetVoteOptionListPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)
			s.table(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).WhereIn("id", ids).Delete() // 投票選項名單資料表

			// 清除遊戲redis資訊
			// s.redisConn.DelCache(config.GAME_REDIS + gameID)
			// s.redisConn.DelCache(config.SCORES_REDIS + gameID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			// s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
			// s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
			// s.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
			return nil
		})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "option_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteOptionListModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(models.EditGameVoteOptionListModel{
				ActivityID: values.Get("activity_id"),
				GameID:     values.Get("game_id"),
				UserID:     values.Get("option_list_id"),
				OptionID:   values.Get("option_id"),
				Name:       values.Get("name"),
				Leader:     values.Get("leader"),
			}); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "option_list_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteOptionListModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(models.EditGameVoteOptionListModel{
				ActivityID: values.Get("activity_id"),
				GameID:     values.Get("game_id"),
				UserID:     values.Get("option_list_id"),
				OptionID:   values.Get("option_id"),
				Name:       values.Get("name"),
				Leader:     values.Get("leader"),
			}); err != nil {
			return err
		}

		return nil
	})
	return
}

// @Summary 新增選項名單人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param option_id formData string true "option_id"
// @param option_list_id formData string true "option_list_id(參加人員加入時)"
// @param name formData integer true "name(匯入excel加入時)"
// @param leader formData integer true "leader"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option_list [post]
func POSTVoteOptionList(ctx *gin.Context) {
}

// @Summary 編輯選項名單人員資料(匹量插入)
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param option_id formData string true "option_id"
// @param people formData integer true "people"
// @param option_list_id formData string true "option_list_id(參加人員加入時，name參數為空)，用逗點間隔，陣列長度需與people一樣"
// @param name formData integer true "name(匯入excel加入時，option_list_id參數為空)，用逗點間隔，陣列長度需與people一樣"
// @param leader formData integer true "leader，用逗點間隔，陣列長度需與people一樣"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option_lists [post]
func POSTVoteOptionLists(ctx *gin.Context) {
}

// @Summary 編輯選項名單人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param option_id formData string true "option_id"
// @param option_list_id formData string true "option_list_id"
// @param name formData integer true "name"
// @param leader formData integer true "leader"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option_list [put]
func PUTVoteOptionList(ctx *gin.Context) {
}

// @Summary 刪除選項名單人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID，id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option_list [delete]
func DELETEVoteOptionList(ctx *gin.Context) {
}

// @Summary 選項名單人員JSON資料
// @Tags Vote
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string true "遊戲ID"
// @param option_id query string true "option_id"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option_list [get]
func VoteOptionListJSON(ctx *gin.Context) {
}

// GetVoteSpecialOfficersPanel 投票特殊人員(匹量匯入)
func (s *SystemTable) GetVoteSpecialOfficersPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	// 面板資訊

	// 表單資訊
	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "special_officer_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteSpecialOfficerModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Adds(true, models.EditGameVoteSpecialOfficerModel{
				ActivityID: values.Get("activity_id"),
				GameID:     values.Get("game_id"),
				People:     values.Get("people"),
				UserID:     values.Get("special_officer_id"),
				Score:      values.Get("score"),
			}); err != nil {
			return err
		}
		return nil
	})

	return
}

// GetVoteSpecialOfficerPanel 投票特殊人員
func (s *SystemTable) GetVoteSpecialOfficerPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).
		SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
			var ids = interfaces(idArr)
			s.table(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).WhereIn("id", ids).Delete() // 投票特殊人員資料表

			// 清除遊戲redis資訊
			s.redisConn.DelCache(config.VOTE_SPECIAL_OFFICER_REDIS + gameID)

			// 更新遊戲場次編輯次數(mongo，會影響遊戲所以需重整)
			filter := bson.M{"game_id": gameID} // 過濾條件
			// 要更新的值
			update := bson.M{
				// "$set": fieldValues,
				// "$unset": bson.M{},                // 移除不需要的欄位
				"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
			}

			if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
				return err
			}

			// 編輯次數更新
			s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")

			// s.redisConn.DelCache(config.GAME_REDIS + gameID)
			// s.redisConn.DelCache(config.SCORES_REDIS + gameID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			// s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
			// s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
			// s.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
			return nil
		})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "special_officer_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteSpecialOfficerModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, models.EditGameVoteSpecialOfficerModel{
				ActivityID:       values.Get("activity_id"),
				GameID:           values.Get("game_id"),
				UserID:           values.Get("special_officer_id"),
				Score:            values.Get("score"),
				Times:            values.Get("times"),
				VoteMethodPlayer: values.Get("vote_method_player"),
			}); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "special_officer_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		if err := models.DefaultGameVoteSpecialOfficerModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, models.EditGameVoteSpecialOfficerModel{
				ActivityID:       values.Get("activity_id"),
				GameID:           values.Get("game_id"),
				UserID:           values.Get("special_officer_id"),
				Score:            values.Get("score"),
				Times:            values.Get("times"),
				VoteMethodPlayer: values.Get("vote_method_player"),
			}); err != nil {
			return err
		}

		return nil
	})
	return
}

// @Summary 新增投票特殊人員資料(匹量插入)
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param people formData integer true "people"
// @param special_officer_id formData string true "special_officer_id，用逗點間隔，陣列長度需與people一樣"
// @param score formData integer true "score，用逗點間隔，陣列長度需與people一樣"
// @param times formData integer true "times，用逗點間隔，陣列長度需與people一樣"
// @param vote_method_player formData integer true "vote_method_player，用逗點間隔，陣列長度需與people一樣"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/special_officers [post]
func POSTVoteSpecialOfficers(ctx *gin.Context) {
}

// @Summary 新增投票特殊人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param special_officer_id formData string true "special_officer_id"
// @param score formData integer true "score"
// @param times formData integer true "times"
// @param vote_method_player formData integer true "vote_method_player"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/special_officer [post]
func POSTVoteSpecialOfficer(ctx *gin.Context) {
}

// @Summary 編輯投票特殊人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param special_officer_id formData string true "special_officer_id"
// @param score formData integer true "score"
// @param times formData integer true "times"
// @param vote_method_player formData integer true "vote_method_player"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/special_officer [put]
func PUTVoteSpecialOfficer(ctx *gin.Context) {
}

// @Summary 刪除投票特殊人員資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID，id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/special_officer [delete]
func DELETEVoteSpecialOfficer(ctx *gin.Context) {
}

// @Summary 投票特殊人員JSON資料
// @Tags Vote
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string true "遊戲ID"
// @param user_id query string true "用戶ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/special_officer [get]
func VoteSpecialOfficerJSON(ctx *gin.Context) {
}

// GetVoteOptionPanel 投票選項
func (s *SystemTable) GetVoteOptionPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_VOTE_OPTION_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)
		s.table(config.ACTIVITY_GAME_VOTE_OPTION_TABLE).WhereIn("option_id", ids).Delete() // 投票選項資料表

		// 更新遊戲場次編輯次數(mongo，會影響遊戲所以需重整)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			// "$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := s.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}

		// 清除遊戲redis資訊
		s.redisConn.DelCache(config.GAME_REDIS + gameID)
		s.redisConn.DelCache(config.SCORES_REDIS + gameID)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		s.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		s.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		s.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

		// 編輯次數更新
		s.redisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+gameID, "修改資料")
		return nil
	})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_VOTE_OPTION_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "option_name") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("option_picture") != "" {
			picture = values.Get("option_picture")
		} else { // 預設
			picture = UPLOAD_SYSTEM_URL + "logo.png"
		}

		if err := models.DefaultGameVoteOptionModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(models.EditGameVoteOptionModel{
				ActivityID:      values.Get("activity_id"),
				GameID:          values.Get("game_id"),
				OptionName:      values.Get("option_name"),
				OptionPicture:   picture,
				OptionIntroduce: values.Get("option_introduce"),
			}); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "option_id") {
			return errors.New("錯誤: 參數資料發生問題，請輸入有效的資料")
		}

		var picture string
		if values.Get("option_picture"+DEFAULT_FALG) == "1" { // 預設
			picture = UPLOAD_SYSTEM_URL + "logo.png"
		} else if values.Get("option_picture") != "" {
			picture = values.Get("option_picture")
		} else {
			picture = ""
		}

		if err := models.DefaultGameVoteOptionModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(models.EditGameVoteOptionModel{
				ActivityID:      values.Get("activity_id"),
				GameID:          values.Get("game_id"),
				OptionID:        values.Get("option_id"),
				OptionName:      values.Get("option_name"),
				OptionPicture:   picture,
				OptionIntroduce: values.Get("option_introduce"),
			}); err != nil {
			return err
		}

		return nil
	})
	return
}

// @Summary 新增投票選項資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param option_name formData string true "option_name"
// @param option_picture formData string true "option_picture"
// @param option_introduce formData string true "option_introduce"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option [post]
func POSTVoteOption(ctx *gin.Context) {
}

// @Summary 編輯投票選項資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param option_id formData string true "option_id"
// @param option_name formData string true "option_name"
// @param option_picture formData string true "option_picture"
// @param option_introduce formData string true "option_introduce"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option [put]
func PUTVoteOption(ctx *gin.Context) {
}

// @Summary 刪除投票選項資料
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID，option_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID(辨識token)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option [delete]
func DELETEVoteOption(ctx *gin.Context) {
}

// @Summary 投票選項JSON資料
// @Tags Vote
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option [get]
func VoteOptionJSON(ctx *gin.Context) {
}

// @Summary 投票選項分數記錄重新計算
// @Tags Vote
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game_id query string true "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/option/reset [get]
func VoteOptionResetJSON(ctx *gin.Context) {
}

// GetVotePanel 投票遊戲
func (s *SystemTable) GetVotePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_TABLE).SetDeleteFunc(func(idArr []string, userID, activityID, gameID string) error {
		var ids = interfaces(idArr)

		// 刪除資料表
		tablesToDelete := []string{
			config.ACTIVITY_PRIZE_TABLE,
			config.ACTIVITY_STAFF_GAME_TABLE,
			config.ACTIVITY_STAFF_PRIZE_TABLE,
			config.ACTIVITY_STAFF_BLACK_TABLE,
			config.ACTIVITY_STAFF_PK_TABLE,
			config.ACTIVITY_SCORE_TABLE,
			config.ACTIVITY_GAME_QA_RECORD_TABLE,

			// 投票
			config.ACTIVITY_GAME_VOTE_OPTION_TABLE,
			config.ACTIVITY_GAME_VOTE_SPECIAL_OFFICER_TABLE,
			config.ACTIVITY_GAME_VOTE_OPTION_LIST_TABLE,
			config.ACTIVITY_GAME_VOTE_RECORD_TABLE,

			// config.ACTIVITY_GAME_TABLE,
			// config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE,
			// config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE,
			// config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE,
			// config.ACTIVITY_GAME_QA_PICTURE_TABLE_1,
			// config.ACTIVITY_GAME_QA_PICTURE_TABLE_2,
			// config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE,
			// config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE,
			// config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE,
			// config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE,
			// config.ACTIVITY_GAME_BINGO_PICTURE_TABLE,
			// config.ACTIVITY_GAME_QA_TABLE,
			// config.ACTIVITY_GAME_VOTE_PICTURE_TABLE,
			// config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE,
		}

		for _, table := range tablesToDelete {
			s.table(table).WhereIn("game_id", ids).Delete()
		}

		// mongo
		mongoTables := []string{
			config.ACTIVITY_GAME_TABLE,
		}
		for _, t := range mongoTables {
			s.mongoConn.DeleteMany(t, bson.M{"game_id": bson.M{"$in": ids}})
		}

		for _, id := range idArr {
			// Redis 要刪除的 key 前綴列表
			delKeys := []string{
				config.GAME_REDIS,
				config.GAME_TYPE_REDIS, // 遊戲種類
				config.GAME_PRIZES_AMOUNT_REDIS,
				config.BLACK_STAFFS_GAME_REDIS,
				config.SCORES_REDIS,
				config.SCORES_2_REDIS,
				config.CORRECT_REDIS,
				config.QA_REDIS,
				config.QA_RECORD_REDIS,
				config.WINNING_STAFFS_REDIS,
				config.NO_WINNING_STAFFS_REDIS, // 未中獎人員
				config.GAME_TEAM_REDIS,
				config.GAME_BINGO_NUMBER_REDIS,               // 紀錄抽過的號碼，LIST
				config.GAME_BINGO_USER_REDIS,                 // 賓果中獎人員，ZSET
				config.GAME_BINGO_USER_NUMBER,                // 紀錄玩家的號碼排序，HASH
				config.GAME_BINGO_USER_GOING_BINGO,           // 紀錄玩家是否即將中獎，HASH
				config.GAME_ATTEND_REDIS,                     // 遊戲人數資訊，SET
				config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS,  // 拔河遊戲左隊人數資訊，SET
				config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS, // 拔河遊戲右隊人數資訊，SET
				// 投票
				config.GAME_VOTE_RECORDS_REDIS,
				config.VOTE_SPECIAL_OFFICER_REDIS,
			}

			for _, key := range delKeys {
				s.redisConn.DelCache(key + id)
			}

			// Redis 要發佈的頻道前綴列表
			pubChannels := []string{
				config.CHANNEL_GAME_REDIS,
				config.CHANNEL_GUEST_GAME_STATUS_REDIS,
				config.CHANNEL_GAME_BINGO_NUMBER_REDIS,
				config.CHANNEL_QA_REDIS,
				config.CHANNEL_GAME_TEAM_REDIS,
				config.CHANNEL_BLACK_STAFFS_GAME_REDIS,
				config.CHANNEL_GAME_EDIT_TIMES_REDIS,
				config.CHANNEL_WINNING_STAFFS_REDIS,
				config.CHANNEL_GAME_BINGO_USER_NUMBER,
				config.CHANNEL_SCORES_REDIS,
			}

			for _, channel := range pubChannels {
				s.redisConn.Publish(channel+id, "修改資料")
			}

			// 刪除遊戲資料夾
			os.RemoveAll(config.STORE_PATH + "/" + userID + "/" + activityID + "/interact/sign/vote/" + id)
		}

		// 刪除遊戲場次時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		// s.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

		// 刪除玩家遊戲紀錄(中獎.未中獎)
		s.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)

		return nil
	})

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_TABLE).SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的資料")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGameModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(votePictureFields, values, true)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		// 手動處理
		model.UserID = values.Get("user") // 抓取的是user參數

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Add(true, "vote", values.Get("game_id"), model); err != nil {
			return err
		}
		return nil
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("game_id") {
			return errors.New("錯誤: 遊戲ID發生問題，請輸入有效的遊戲ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGameModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		// 判斷圖片參數是否為空，將路徑參數寫入map中
		picMap := BuildPictureMap(votePictureFields, values, false)

		// 將圖片資料寫入struct中
		if err := utils.ApplyMapToStruct(picMap, &model); err != nil {
			return errors.New("錯誤: 套用圖片資料失敗")
		}

		if err := models.DefaultGameModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, "vote", model); err != nil {
			return err
		}
		return nil
	})
	return
}

// @Summary 新增投票遊戲資料(form-data)
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param title formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param vote_screen formData string true "投票畫面" Enums(bar_chart, rank, detail_information)
// @param vote_times formData integer true "人員投票次數"
// @param vote_method formData string true "投票模式" Enums(all_vote, single_group, all_group)
// @param vote_method_player formData string true "玩家投票方式" Enums(one_vote, free_vote)
// @param vote_restriction formData string true "投票限制" Enums(all_player, special_officer)
// @param prize formData string true "發獎方式" Enums(uniform, all)
// @param avatar_shape formData string true "選項框" Enums(circle, square)
// @param auto_play formData string false "auto_play" Enums(open, close)
// @param show_rank formData string false "show_rank" Enums(open, close)
// @param title_switch formData string false "title_switch" Enums(open, close)
// @param arrangement_guest formData string false "arrangement_guest" Enums(list, side_by_side)
// @param vote_start_time formData string true "投票開始時間"
// @param vote_end_time formData string true "投票結束時間"
// @param topic formData string true "主題樣式" Enums(01_classic, 02_starrysky)
// @param skin formData string true "樣式選擇" Enums(classic, customize)
// @param music formData string true "音樂選擇" Enums(classic, customize)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/form [post]
func POSTVote(ctx *gin.Context) {
}

// @Summary 新增投票遊戲獎品資料(form-data)
// @Tags Vote Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_name formData string true "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string true "類型" Enums(first, second, third, general, thanks)
// @param prize_picture formData file false "照片"
// @param prize_amount formData integer true "數量"
// @param prize_price formData integer true "價值"
// @param prize_method formData string true "兌獎方式" Enums(site, mail, thanks)
// @param prize_password formData string true "兌獎密碼"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/prize/form [post]
func POSTVotePrize(ctx *gin.Context) {
}

// @Summary 編輯投票遊戲資料(form-data)
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param title formData string false "標題(上限為20個字元)" minlength(1) maxlength(20)
// @param vote_screen formData string false "投票畫面" Enums(bar_chart, rank, detail_information)
// @param vote_times formData integer false "人員投票次數"
// @param vote_method formData string false "投票模式" Enums(all_vote, single_group, all_group)
// @param vote_method_player formData string false "玩家投票方式" Enums(one_vote, free_vote)
// @param vote_restriction formData string false "投票限制" Enums(all_player, special_officer)
// @param prize formData string false "發獎方式" Enums(uniform, all)
// @param avatar_shape formData string false "選項框" Enums(circle, square)
// @param auto_play formData string false "auto_play" Enums(open, close)
// @param show_rank formData string false "show_rank" Enums(open, close)
// @param title_switch formData string false "title_switch" Enums(open, close)
// @param arrangement_guest formData string false "arrangement_guest" Enums(list, side_by_side)
// @param vote_start_time formData string false "投票開始時間"
// @param vote_end_time formData string false "投票結束時間"
// @param topic formData string false "主題樣式" Enums(01_classic, 02_starrysky)
// @param skin formData string false "樣式選擇" Enums(classic, customize)
// @param music formData string false "音樂選擇" Enums(classic, customize)
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/form [put]
func PUTVote(ctx *gin.Context) {
}

// @Summary 編輯投票遊戲獎品資料(form-data)
// @Tags Vote Prize
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param prize_id formData string true "獎品ID"
// @param prize_name formData string false "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param prize_type formData string false "類型" Enums(first, second, third, general, thanks)
// @param prize_picture formData file false "照片"
// @param prize_amount formData integer false "數量"
// @param prize_price formData integer false "價值"
// @param prize_method formData string false "兌獎方式" Enums(site, mail, thanks)
// @param prize_password formData string false "兌獎密碼"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/prize/form [put]
func PUTVotePrize(ctx *gin.Context) {
}

// @Summary 刪除投票資料(form-data)
// @Tags Vote
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param user formData string true "用戶ID(該場次活動的user_id)"
// @param activity_id formData string true "活動ID"
// @param user_id formData string true "用戶ID(辨識token用)"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/form [delete]
func DELETEVote(ctx *gin.Context) {
}

// @Summary 刪除投票獎品資料(form-data)
// @Tags Vote Prize
// @version 1.0
// @Accept  mpfd
// @param id formData string true "ID(以,區隔多筆資料ID)"
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/prize/form [delete]
func DELETEVotePrize(ctx *gin.Context) {
}

// @Summary 投票遊戲JSON資料(sign分類下)
// @Tags Vote
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Param isredis query bool true "redis"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote [get]
func VoteJSON(ctx *gin.Context) {
}

// @Summary 投票遊戲獎品JSON資料(sign分類下)
// @Tags Vote Prize
// @version 1.0
// @Accept  json
// @param activity_id query string false "活動ID"
// @param game_id query string false "遊戲ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/sign/vote/prize [get]
func VotePrizeJSON(ctx *gin.Context) {
}

// for i, field := range fields {
// 	if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	}
// }

// models.NewGameModel{
// 	UserID:        values.Get("user"),
// 	ActivityID:    values.Get("activity_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     "close",
// 	Second:        "0",
// 	MaxPeople:     "0",
// 	People:        "0",
// 	MaxTimes:      "0",
// 	Allow:         "open",
// 	Percent:       "0",
// 	FirstPrize:    "0",
// 	SecondPrize:   "0",
// 	ThirdPrize:    "0",
// 	GeneralPrize:  "0",
// 	Topic:         values.Get("topic"),
// 	Skin:          values.Get("skin"),
// 	Music:         values.Get("music"),
// 	DisplayName:   "open",
// 	BoxReflection: "",
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            values.Get("prize"),

// 	// 賓果遊戲
// 	MaxNumber:  "0",
// 	BingoLine:  "0",
// 	RoundPrize: "0",

// 	// 扭蛋機遊戲
// 	GachaMachineReflection: "0",
// 	ReflectiveSwitch:       "",

// 	// 投票遊戲
// 	VoteScreen:       values.Get("vote_screen"),
// 	VoteTimes:        values.Get("vote_times"),
// 	VoteMethod:       values.Get("vote_method"),
// 	VoteMethodPlayer: values.Get("vote_method_player"),
// 	VoteRestriction:  values.Get("vote_restriction"),
// 	AvatarShape:      values.Get("avatar_shape"),
// 	VoteStartTime:    values.Get("vote_start_time"),
// 	VoteEndTime:      values.Get("vote_end_time"),
// 	AutoPlay:         values.Get("auto_play"),
// 	ShowRank:         values.Get("show_rank"),
// 	TitleSwitch:      values.Get("title_switch"),
// 	ArrangementGuest: values.Get("arrangement_guest"),

// 	// 投票遊戲自定義
// 	VoteClassicHPic01: update[0],
// 	VoteClassicHPic02: update[1],
// 	VoteClassicHPic03: update[2],
// 	VoteClassicHPic04: update[3],
// 	VoteClassicHPic05: update[4],
// 	VoteClassicHPic06: update[5],
// 	VoteClassicHPic07: update[6],
// 	VoteClassicHPic08: update[7],
// 	VoteClassicHPic09: update[8],
// 	VoteClassicHPic10: update[9],
// 	VoteClassicHPic11: update[10],
// 	// VoteClassicHPic12: update[11],
// 	VoteClassicHPic13: update[11],
// 	VoteClassicHPic14: update[12],
// 	VoteClassicHPic15: update[13],
// 	VoteClassicHPic16: update[14],
// 	VoteClassicHPic17: update[15],
// 	VoteClassicHPic18: update[16],
// 	VoteClassicHPic19: update[17],
// 	VoteClassicHPic20: update[18],
// 	VoteClassicHPic21: update[19],
// 	// VoteClassicHPic22: update[21],
// 	VoteClassicHPic23: update[20],
// 	VoteClassicHPic24: update[21],
// 	VoteClassicHPic25: update[22],
// 	VoteClassicHPic26: update[23],
// 	VoteClassicHPic27: update[24],
// 	VoteClassicHPic28: update[25],
// 	VoteClassicHPic29: update[26],
// 	VoteClassicHPic30: update[27],
// 	VoteClassicHPic31: update[28],
// 	VoteClassicHPic32: update[29],
// 	VoteClassicHPic33: update[30],
// 	VoteClassicHPic34: update[31],
// 	VoteClassicHPic35: update[32],
// 	VoteClassicHPic36: update[33],
// 	VoteClassicHPic37: update[34],
// 	VoteClassicGPic01: update[35],
// 	VoteClassicGPic02: update[36],
// 	VoteClassicGPic03: update[37],
// 	VoteClassicGPic04: update[38],
// 	VoteClassicGPic05: update[39],
// 	VoteClassicGPic06: update[40],
// 	VoteClassicGPic07: update[41],
// 	VoteClassicCPic01: update[42],
// 	VoteClassicCPic02: update[43],
// 	VoteClassicCPic03: update[44],
// 	VoteClassicCPic04: update[45],

// 	// 音樂
// 	VoteBgmGaming: update[46],
// }

// pics = []string{
// 	// 投票遊戲自定義
// 	"vote/classic/vote_classic_h_pic_01.png",
// 	"vote/classic/vote_classic_h_pic_02.png",
// 	"vote/classic/vote_classic_h_pic_03.png",
// 	"vote/classic/vote_classic_h_pic_04.png",
// 	"vote/classic/vote_classic_h_pic_05.png",
// 	"vote/classic/vote_classic_h_pic_06.png",
// 	"vote/classic/vote_classic_h_pic_07.png",
// 	"vote/classic/vote_classic_h_pic_08.jpg",
// 	"vote/classic/vote_classic_h_pic_09.png",
// 	"vote/classic/vote_classic_h_pic_10.png",
// 	"vote/classic/vote_classic_h_pic_11.png",
// 	// "vote/classic/vote_classic_h_pic_12.png",
// 	"vote/classic/vote_classic_h_pic_13.png",
// 	"vote/classic/vote_classic_h_pic_14.png",
// 	"vote/classic/vote_classic_h_pic_15.png",
// 	"vote/classic/vote_classic_h_pic_16.png",
// 	"vote/classic/vote_classic_h_pic_17.png",
// 	"vote/classic/vote_classic_h_pic_18.png",
// 	"vote/classic/vote_classic_h_pic_19.png",
// 	"vote/classic/vote_classic_h_pic_20.png",
// 	"vote/classic/vote_classic_h_pic_21.png",
// 	// "vote/classic/vote_classic_h_pic_22.png",
// 	"vote/classic/vote_classic_h_pic_23.png",
// 	"vote/classic/vote_classic_h_pic_24.png",
// 	"vote/classic/vote_classic_h_pic_25.png",
// 	"vote/classic/vote_classic_h_pic_26.png",
// 	"vote/classic/vote_classic_h_pic_27.png",
// 	"vote/classic/vote_classic_h_pic_28.png",
// 	"vote/classic/vote_classic_h_pic_29.png",
// 	"vote/classic/vote_classic_h_pic_30.png",
// 	"vote/classic/vote_classic_h_pic_31.png",
// 	"vote/classic/vote_classic_h_pic_32.png",
// 	"vote/classic/vote_classic_h_pic_33.png",
// 	"vote/classic/vote_classic_h_pic_34.png",
// 	"vote/classic/vote_classic_h_pic_35.png",
// 	"vote/classic/vote_classic_h_pic_36.png",
// 	"vote/classic/vote_classic_h_pic_37.png",
// 	"vote/classic/vote_classic_g_pic_01.png",
// 	"vote/classic/vote_classic_g_pic_02.png",
// 	"vote/classic/vote_classic_g_pic_03.png",
// 	"vote/classic/vote_classic_g_pic_04.jpg",
// 	"vote/classic/vote_classic_g_pic_05.png",
// 	"vote/classic/vote_classic_g_pic_06.png",
// 	"vote/classic/vote_classic_g_pic_07.png",
// 	"vote/classic/vote_classic_c_pic_01.png",
// 	"vote/classic/vote_classic_c_pic_02.png",
// 	"vote/classic/vote_classic_c_pic_03.png",
// 	"vote/classic/vote_classic_c_pic_04.png",

// 	// 音樂
// 	"vote/%s/bgm/gaming.mp3",
// }
// fields = []string{
// 	// 投票遊戲自定義
// 	"vote_classic_h_pic_01",
// 	"vote_classic_h_pic_02",
// 	"vote_classic_h_pic_03",
// 	"vote_classic_h_pic_04",
// 	"vote_classic_h_pic_05",
// 	"vote_classic_h_pic_06",
// 	"vote_classic_h_pic_07",
// 	"vote_classic_h_pic_08",
// 	"vote_classic_h_pic_09",
// 	"vote_classic_h_pic_10",
// 	"vote_classic_h_pic_11",
// 	// "vote_classic_h_pic_12",
// 	"vote_classic_h_pic_13",
// 	"vote_classic_h_pic_14",
// 	"vote_classic_h_pic_15",
// 	"vote_classic_h_pic_16",
// 	"vote_classic_h_pic_17",
// 	"vote_classic_h_pic_18",
// 	"vote_classic_h_pic_19",
// 	"vote_classic_h_pic_20",
// 	"vote_classic_h_pic_21",
// 	// "vote_classic_h_pic_22",
// 	"vote_classic_h_pic_23",
// 	"vote_classic_h_pic_24",
// 	"vote_classic_h_pic_25",
// 	"vote_classic_h_pic_26",
// 	"vote_classic_h_pic_27",
// 	"vote_classic_h_pic_28",
// 	"vote_classic_h_pic_29",
// 	"vote_classic_h_pic_30",
// 	"vote_classic_h_pic_31",
// 	"vote_classic_h_pic_32",
// 	"vote_classic_h_pic_33",
// 	"vote_classic_h_pic_34",
// 	"vote_classic_h_pic_35",
// 	"vote_classic_h_pic_36",
// 	"vote_classic_h_pic_37",
// 	"vote_classic_g_pic_01",
// 	"vote_classic_g_pic_02",
// 	"vote_classic_g_pic_03",
// 	"vote_classic_g_pic_04",
// 	"vote_classic_g_pic_05",
// 	"vote_classic_g_pic_06",
// 	"vote_classic_g_pic_07",
// 	"vote_classic_c_pic_01",
// 	"vote_classic_c_pic_02",
// 	"vote_classic_c_pic_03",
// 	"vote_classic_c_pic_04",

// 	// 音樂
// 	"vote_bgm_gaming", // 遊戲進行中
// }
// update = make([]string, 100)

// var (
// 	pics = []string{
// 		// 投票遊戲自定義
// 		"vote/classic/vote_classic_h_pic_01.png",
// 		"vote/classic/vote_classic_h_pic_02.png",
// 		"vote/classic/vote_classic_h_pic_03.png",
// 		"vote/classic/vote_classic_h_pic_04.png",
// 		"vote/classic/vote_classic_h_pic_05.png",
// 		"vote/classic/vote_classic_h_pic_06.png",
// 		"vote/classic/vote_classic_h_pic_07.png",
// 		"vote/classic/vote_classic_h_pic_08.jpg",
// 		"vote/classic/vote_classic_h_pic_09.png",
// 		"vote/classic/vote_classic_h_pic_10.png",
// 		"vote/classic/vote_classic_h_pic_11.png",
// 		// "vote/classic/vote_classic_h_pic_12.png",
// 		"vote/classic/vote_classic_h_pic_13.png",
// 		"vote/classic/vote_classic_h_pic_14.png",
// 		"vote/classic/vote_classic_h_pic_15.png",
// 		"vote/classic/vote_classic_h_pic_16.png",
// 		"vote/classic/vote_classic_h_pic_17.png",
// 		"vote/classic/vote_classic_h_pic_18.png",
// 		"vote/classic/vote_classic_h_pic_19.png",
// 		"vote/classic/vote_classic_h_pic_20.png",
// 		"vote/classic/vote_classic_h_pic_21.png",
// 		// "vote/classic/vote_classic_h_pic_22.png",
// 		"vote/classic/vote_classic_h_pic_23.png",
// 		"vote/classic/vote_classic_h_pic_24.png",
// 		"vote/classic/vote_classic_h_pic_25.png",
// 		"vote/classic/vote_classic_h_pic_26.png",
// 		"vote/classic/vote_classic_h_pic_27.png",
// 		"vote/classic/vote_classic_h_pic_28.png",
// 		"vote/classic/vote_classic_h_pic_29.png",
// 		"vote/classic/vote_classic_h_pic_30.png",
// 		"vote/classic/vote_classic_h_pic_31.png",
// 		"vote/classic/vote_classic_h_pic_32.png",
// 		"vote/classic/vote_classic_h_pic_33.png",
// 		"vote/classic/vote_classic_h_pic_34.png",
// 		"vote/classic/vote_classic_h_pic_35.png",
// 		"vote/classic/vote_classic_h_pic_36.png",
// 		"vote/classic/vote_classic_h_pic_37.png",
// 		"vote/classic/vote_classic_g_pic_01.png",
// 		"vote/classic/vote_classic_g_pic_02.png",
// 		"vote/classic/vote_classic_g_pic_03.png",
// 		"vote/classic/vote_classic_g_pic_04.jpg",
// 		"vote/classic/vote_classic_g_pic_05.png",
// 		"vote/classic/vote_classic_g_pic_06.png",
// 		"vote/classic/vote_classic_g_pic_07.png",
// 		"vote/classic/vote_classic_c_pic_01.png",
// 		"vote/classic/vote_classic_c_pic_02.png",
// 		"vote/classic/vote_classic_c_pic_03.png",
// 		"vote/classic/vote_classic_c_pic_04.png",

// 		// 音樂
// 		"vote/%s/bgm/gaming.mp3",
// 	}
// 	fields = []string{
// 		// 投票遊戲自定義
// 		"vote_classic_h_pic_01",
// 		"vote_classic_h_pic_02",
// 		"vote_classic_h_pic_03",
// 		"vote_classic_h_pic_04",
// 		"vote_classic_h_pic_05",
// 		"vote_classic_h_pic_06",
// 		"vote_classic_h_pic_07",
// 		"vote_classic_h_pic_08",
// 		"vote_classic_h_pic_09",
// 		"vote_classic_h_pic_10",
// 		"vote_classic_h_pic_11",
// 		// "vote_classic_h_pic_12",
// 		"vote_classic_h_pic_13",
// 		"vote_classic_h_pic_14",
// 		"vote_classic_h_pic_15",
// 		"vote_classic_h_pic_16",
// 		"vote_classic_h_pic_17",
// 		"vote_classic_h_pic_18",
// 		"vote_classic_h_pic_19",
// 		"vote_classic_h_pic_20",
// 		"vote_classic_h_pic_21",
// 		// "vote_classic_h_pic_22",
// 		"vote_classic_h_pic_23",
// 		"vote_classic_h_pic_24",
// 		"vote_classic_h_pic_25",
// 		"vote_classic_h_pic_26",
// 		"vote_classic_h_pic_27",
// 		"vote_classic_h_pic_28",
// 		"vote_classic_h_pic_29",
// 		"vote_classic_h_pic_30",
// 		"vote_classic_h_pic_31",
// 		"vote_classic_h_pic_32",
// 		"vote_classic_h_pic_33",
// 		"vote_classic_h_pic_34",
// 		"vote_classic_h_pic_35",
// 		"vote_classic_h_pic_36",
// 		"vote_classic_h_pic_37",
// 		"vote_classic_g_pic_01",
// 		"vote_classic_g_pic_02",
// 		"vote_classic_g_pic_03",
// 		"vote_classic_g_pic_04",
// 		"vote_classic_g_pic_05",
// 		"vote_classic_g_pic_06",
// 		"vote_classic_g_pic_07",
// 		"vote_classic_c_pic_01",
// 		"vote_classic_c_pic_02",
// 		"vote_classic_c_pic_03",
// 		"vote_classic_c_pic_04",

// 		// 音樂
// 		"vote_bgm_gaming", // 遊戲進行中
// 	}
// 	update = make([]string, 100)
// )

// for i, field := range fields {
// 	if values.Get(field+DEFAULT_FALG) == "1" {
// 		if strings.Contains(pics[i], "%s") {
// 			topics := strings.Split(values.Get("topic"), "_")
// 			topic := ""
// 			if len(topics) == 2 {
// 				topic = topics[1]
// 			} else if len(topics) == 3 {
// 				topic = topics[1] + "_" + topics[2]
// 			}

// 			update[i] = UPLOAD_SYSTEM_URL + fmt.Sprintf(pics[i], topic)
// 		} else {
// 			update[i] = UPLOAD_SYSTEM_URL + pics[i]
// 		}
// 	} else if values.Get(field) != "" {
// 		update[i] = values.Get(field)
// 	} else {
// 		update[i] = ""
// 	}
// }

// models.EditGameModel{
// 	ActivityID:    values.Get("activity_id"),
// 	GameID:        values.Get("game_id"),
// 	Title:         values.Get("title"),
// 	GameType:      "",
// 	LimitTime:     "",
// 	Second:        "",
// 	MaxPeople:     "",
// 	People:        "",
// 	MaxTimes:      "",
// 	Allow:         "",
// 	Percent:       "",
// 	FirstPrize:    "",
// 	SecondPrize:   "",
// 	ThirdPrize:    "",
// 	GeneralPrize:  "",
// 	Topic:         values.Get("topic"),
// 	Skin:          values.Get("skin"),
// 	Music:         values.Get("music"),
// 	DisplayName:   "",
// 	GameOrder:     values.Get("game_order"),
// 	BoxReflection: "",
// 	SamePeople:    "",

// 	// 拔河遊戲
// 	AllowChooseTeam:  "",
// 	LeftTeamName:     "",
// 	LeftTeamPicture:  "",
// 	RightTeamName:    "",
// 	RightTeamPicture: "",
// 	Prize:            values.Get("prize"),

// 	// 賓果遊戲
// 	MaxNumber:  "",
// 	BingoLine:  "",
// 	RoundPrize: "",

// 	// 扭蛋機遊戲
// 	GachaMachineReflection: "",
// 	ReflectiveSwitch:       "",

// 	// 投票遊戲
// 	VoteScreen:       values.Get("vote_screen"),
// 	VoteTimes:        values.Get("vote_times"),
// 	VoteMethod:       values.Get("vote_method"),
// 	VoteMethodPlayer: values.Get("vote_method_player"),
// 	VoteRestriction:  values.Get("vote_restriction"),
// 	AvatarShape:      values.Get("avatar_shape"),
// 	VoteStartTime:    values.Get("vote_start_time"),
// 	VoteEndTime:      values.Get("vote_end_time"),
// 	AutoPlay:         values.Get("auto_play"),
// 	ShowRank:         values.Get("show_rank"),
// 	TitleSwitch:      values.Get("title_switch"),
// 	ArrangementGuest: values.Get("arrangement_guest"),

// 	// 投票遊戲自定義
// 	VoteClassicHPic01: update[0],
// 	VoteClassicHPic02: update[1],
// 	VoteClassicHPic03: update[2],
// 	VoteClassicHPic04: update[3],
// 	VoteClassicHPic05: update[4],
// 	VoteClassicHPic06: update[5],
// 	VoteClassicHPic07: update[6],
// 	VoteClassicHPic08: update[7],
// 	VoteClassicHPic09: update[8],
// 	VoteClassicHPic10: update[9],
// 	VoteClassicHPic11: update[10],
// 	// VoteClassicHPic12: update[11],
// 	VoteClassicHPic13: update[11],
// 	VoteClassicHPic14: update[12],
// 	VoteClassicHPic15: update[13],
// 	VoteClassicHPic16: update[14],
// 	VoteClassicHPic17: update[15],
// 	VoteClassicHPic18: update[16],
// 	VoteClassicHPic19: update[17],
// 	VoteClassicHPic20: update[18],
// 	VoteClassicHPic21: update[19],
// 	// VoteClassicHPic22: update[21],
// 	VoteClassicHPic23: update[20],
// 	VoteClassicHPic24: update[21],
// 	VoteClassicHPic25: update[22],
// 	VoteClassicHPic26: update[23],
// 	VoteClassicHPic27: update[24],
// 	VoteClassicHPic28: update[25],
// 	VoteClassicHPic29: update[26],
// 	VoteClassicHPic30: update[27],
// 	VoteClassicHPic31: update[28],
// 	VoteClassicHPic32: update[29],
// 	VoteClassicHPic33: update[30],
// 	VoteClassicHPic34: update[31],
// 	VoteClassicHPic35: update[32],
// 	VoteClassicHPic36: update[33],
// 	VoteClassicHPic37: update[34],
// 	VoteClassicGPic01: update[35],
// 	VoteClassicGPic02: update[36],
// 	VoteClassicGPic03: update[37],
// 	VoteClassicGPic04: update[38],
// 	VoteClassicGPic05: update[39],
// 	VoteClassicGPic06: update[40],
// 	VoteClassicGPic07: update[41],
// 	VoteClassicCPic01: update[42],
// 	VoteClassicCPic02: update[43],
// 	VoteClassicCPic03: update[44],
// 	VoteClassicCPic04: update[45],

// 	// 音樂
// 	VoteBgmGaming: update[46],
// }
