package controller

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// GetCustomizeTemplates 所有自定義模板json資料，GET API
func (h *Handler) GetCustomizeTemplates(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		game = ctx.Query("game")
	)

	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if game == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得遊戲資訊，請輸入有效的參數",
		})
		return
	}

	// 多個自定義模板資料
	templates, err := models.DefaultCustomizeTemplateModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAll(game)

	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, templates)
}

// GetCustomizeScenes 所有自定義場景json資料，GET API
func (h *Handler) GetCustomizeScenes(ctx *gin.Context) {
	var (
		host   = ctx.Request.Host
		userID = ctx.Query("user_id")
	)

	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if userID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得用戶資訊，請輸入有效的參數",
		})
		return
	}

	// 多個自定義場景資料
	scenes, err := models.DefaultUserCustomizeSceneModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAll(userID)

	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, scenes)
}

// TestApi4 測試用api
func (h *Handler) TestApi4(ctx *gin.Context) {
	// log.Println("查詢單筆資料")

	// 自定義場景資料
	game, _ := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, ctx.Query("game_id"), ctx.Query("game"))

	response.OkWithData(ctx, game)

	// 	log.Println("查詢單筆資料")

	// 	filter := bson.M{"name": "Event A"}

	// 	result, err := h.mongoConn.FindOne("activity", filter)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	log.Println("查詢到的活動:", result)
	// 	log.Println("err: ", err)

	// 	log.Println("查詢多筆資料")

	// 	filter = bson.M{"channel_3": "open"}

	// 	results, err := h.mongoConn.FindMany("activity", filter, 1000)
	// 	log.Println("查詢到的活動:", results)
	// 	log.Println("err: ", err)
}

// TestApi 測試用api
// func (h *Handler) TestApi(ctx *gin.Context) {
// 	// applySignModel, _ := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	// 	Find(0, "LevJ8qMfsJb9TlC0UXcj", "U3140fd73cfd35dd992668ab3b6efdae9", true)

// 	// response.OkWithData(ctx, applySignModel)
// 	log.Println("插入單筆資料")

// 	activity := bson.M{
// 		"name":       "Event A",
// 		"channel_1":  "close",
// 		"channel_2":  "close",
// 		"channel_3":  "close",
// 		"channel_4":  "close",
// 		"channel_5":  "close",
// 		"channel_6":  "close",
// 		"channel_7":  "close",
// 		"channel_8":  "close",
// 		"channel_9":  "close",
// 		"channel_10": "close",
// 	}

// 	result, err := h.mongoConn.InsertOne("activity", activity)
// 	log.Println("result: ", result)
// 	log.Println("err: ", err)

// 	log.Println("插入多筆資料")

// 	activitys := []interface{}{
// 		map[string]interface{}{
// 			"name":       "Event B",
// 			"channel_1":  "close",
// 			"channel_2":  "close",
// 			"channel_3":  "close",
// 			"channel_4":  "close",
// 			"channel_5":  "close",
// 			"channel_6":  "close",
// 			"channel_7":  "close",
// 			"channel_8":  "close",
// 			"channel_9":  "close",
// 			"channel_10": "close",
// 		},
// 		map[string]interface{}{
// 			"name":       "Event C",
// 			"channel_1":  "close",
// 			"channel_2":  "close",
// 			"channel_3":  "close",
// 			"channel_4":  "close",
// 			"channel_5":  "close",
// 			"channel_6":  "close",
// 			"channel_7":  "close",
// 			"channel_8":  "close",
// 			"channel_9":  "close",
// 			"channel_10": "close",
// 		},
// 	}

// 	results, err := h.mongoConn.InsertMany("activity", activitys)
// 	log.Println("result: ", results)
// 	log.Println("err: ", err)
// }

// // TestApi2 測試用api
// func (h *Handler) TestApi2(ctx *gin.Context) {
// 	log.Println("更新單筆資料")

// 	filter := bson.M{"name": "Event A"} // 只更新 "活動A"
// 	update := bson.M{
// 		"$set":   bson.M{"channel_1": "open", "channel_2": "open"},
// 		"$unset": bson.M{"channel_6": "", "channel_7": "", "channel_8": "", "channel_9": "", "channel_10": ""}, // 移除不需要的欄位
// 	}

// 	result, err := h.mongoConn.UpdateOne("activity", filter, update)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println("更新的文件數:", result.ModifiedCount)
// 	log.Println("err: ", err)

// 	log.Println("更新多筆資料")

// 	filter = bson.M{"name": bson.M{"$in": []string{"Event B", "Event C"}}}
// 	update = bson.M{
// 		"$set":   bson.M{"channel_1": "open", "channel_3": "open", "channel_5": "open", "channel_7": "open"},
// 		"$unset": bson.M{"channel_8": "", "channel_9": "", "channel_10": ""}, // 移除不需要的欄位
// 	}

// 	results, err := h.mongoConn.UpdateMany("activity", filter, update)
// 	log.Println("更新的文件數:", results.ModifiedCount)
// 	log.Println("err: ", err)
// }

// // TestApi3 測試用api
// func (h *Handler) TestApi3(ctx *gin.Context) {
// 	log.Println("刪除單筆資料")

// 	filter := bson.M{"name": "Event A"}

// 	result, err := h.mongoConn.DeleteOne("activity", filter)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println("刪除的文件數:", result.DeletedCount)
// 	log.Println("err: ", err)

// 	log.Println("刪除多筆資料")

// 	filter = bson.M{"name": bson.M{"$in": []string{"Event B", "Event C"}}}

// 	results, err := h.mongoConn.DeleteMany("activity", filter)
// 	log.Println("刪除的文件數:", results.DeletedCount)
// 	log.Println("err: ", err)
// }

// GetAdmin 管理員json資料 GET API
func (h *Handler) GetAdmin(ctx *gin.Context) {
	var (
		host   = ctx.Request.Host
		prefix = ctx.Param("__prefix")
		id     = ctx.Query("id")
		data   interface{}
		idInt  int
		err    error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	idInt, err = strconv.Atoi(id)
	if err != nil {
		idInt = 0
	}

	if prefix == "manager" {
		userID := ctx.Query("user_id")
		if userID != "" {
			// *****可以同時登入(暫時拿除)*****
			// data, err = models.DefaultUserModel().
			// 	SetDbConn(h.dbConn).Find(false, true, "", "users.user_id", userID)
			// *****可以同時登入(暫時拿除)*****

			// *****可以同時登入(新)*****
			data, err = models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(false, true, "users.user_id", userID)
			// *****可以同時登入(新)*****
		} else if userID == "" {
			data, err = models.DefaultUserModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindUsers()
		}
	} else if prefix == "permission" {
		data, err = models.DefaultPermissionModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(int64(idInt))
	} else if prefix == "menu" {
		data, err = models.DefaultMenuModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(int64(idInt))
	} else if prefix == "overview" {
		data, err = models.DefaultOverviewModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(int64(idInt))

	} else if prefix == "error_log" {
		data, err = models.DefaultErrorLogModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAll()
	} else if prefix == "log" {
		data, err = models.DefaultLogModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAll()
	}
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得資訊，請輸入有效的資料",
		})
		return
	}

	response.OkWithData(ctx, data)
}

// GetApplysign 報名簽到人員json資料 GET API
func (h *Handler) GetApplySign(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		userID     = ctx.Query("user_id")
		status     = ctx.Query("status")
		limit      = ctx.Query("limit")
		offset     = ctx.Query("offset")
		applysigns = make([]models.ApplysignModel, 0)
		err        error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if userID == "" && activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識報名簽到人員資訊，請輸入有效的參數",
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	if offsetInt != 0 && limitInt == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識參數資訊，請輸入有效的limit.offset參數",
		})
		return
	}

	// 查詢用戶所有活動報名簽到狀態: user_id
	// 查詢所有報名簽到人員資料(所有狀態): activity_id
	// 查詢所有報名簽到人員資料(特定狀態): activity_id, status
	applysigns, err = models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAll(activityID, userID, "", status,
			int64(limitInt), int64(offsetInt))
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, applysigns)
}

// GetApplysignCustomize 報名簽到自定義欄位json資料 GET API
func (h *Handler) GetApplysignCustomize(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		err        error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識資訊，請輸入有效的參數",
		})
		return
	}

	// 自定義欄位資訊
	customizeModel, err := models.DefaultCustomizeModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(activityID)
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, customizeModel)
}

// GetRecords 紀錄json資料 GET API
func (h *Handler) GetRecords(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		activityID = ctx.Query("activity_id")
		// gameID     = ctx.Query("game_id")
		// round      = ctx.Query("round")
		// roundInt   int
		data interface{}
		err  error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	// 判斷是否取得特定輪次資料
	// if round != "" {
	// 	roundInt, err = strconv.Atoi(round)
	// 	if err != nil {
	// 		response.BadRequest(ctx, "錯誤: 輪次參數發生問題")
	// 		return
	// 	}
	// }

	if strings.Contains(path, "/v1/chatroom/record") {
		// 聊天室紀錄
		data, err = models.DefaultChatroomRecordModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAll(activityID)
		if activityID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識活動資訊，請輸入有效的活動ID",
			})
			return
		}
	} else if strings.Contains(path, "/v1/question/record") {
		// 提問牆紀錄
		data, err = models.DefaultQuestionUserModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindAll(activityID)
		if activityID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識活動資訊，請輸入有效的活動ID",
			})
			return
		}
	}
	//  else if strings.Contains(path, "/v1/QA/record") {
	// 	// 快問快答答題紀錄
	// 	data, err = models.DefaultGameQARecordModel().
	// 		SetDbConn(h.dbConn).FindAllByDB(gameID, int64(roundInt))
	// 	if gameID == "" {
	// 		response.BadRequest(ctx, "錯誤: 無法辨識遊戲資訊，請輸入有效的遊戲ID")
	// 		return
	// 	}
	// }
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, data)
}

// GetWall 聊天區設定json資料 GET API
func (h *Handler) GetWall(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if prefix == "message_sensitivity" {
		records, err := models.DefaultMessageSensitivityModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(activityID)
		if err != nil || activityID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識敏感詞資料，請輸入有效的活動ID",
			})
			return
		}

		response.OkWithData(ctx, records)
	}
}

// GetActivity 活動資訊json資料 GET API
func (h *Handler) GetActivity(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		data       interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if activityID != "" {
		activitys, err := models.DefaultActivityModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindActivityPermissions("activity.activity_id", activityID)
		if err != nil || len(activitys) == 0 {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識活動資訊，請輸入有效的活動ID",
			})
			return
		}

		data = activitys[0]
	} else if activityID == "" {
		activitys, err := models.DefaultActivityModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindOpenActivitys()
		if err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}

		data = activitys
	}

	response.OkWithData(ctx, data)
}

// GetInfo 活動資訊json資料(介紹，行程，嘉賓，資料) GET API
func (h *Handler) GetInfo(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		data       interface{}
		err        error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if prefix == "introduce" {
		data, err = models.DefaultIntroduceModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find("activity_introduce.activity_id", activityID)
	} else if prefix == "schedule" {
		data, err = models.DefaultScheduleModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find("activity_id", activityID)
	} else if prefix == "guest" {
		data, err = models.DefaultGuestModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find("activity_guest.activity_id", activityID)
	} else if prefix == "material" {
		data, err = models.DefaultMaterialModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find("activity_id", activityID)
	}
	if err != nil || activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得活動資訊(介紹，行程，嘉賓，資料)，請輸入有效的活動ID",
		})
		return
	}
	response.OkWithData(ctx, data)
}

// GetSignGames 所有遊戲json資料(sign分類下)，GET API
func (h *Handler) GetSignGames(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		game       = ctx.Query("game")
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得活動資訊，請輸入有效的活動ID",
		})
		return
	}

	// 多個遊戲場次資料
	gamesModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindAll(activityID, game)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, gamesModel)
}

// GetGames 所有遊戲json資料(或特定遊戲類型) GET API
func (h *Handler) GetGames(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		game       = ctx.Query("game")
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得活動資訊，請輸入有效的活動ID",
		})
		return
	}

	// 多個遊戲場次資料
	gamesModel, err := models.DefaultGameModel().
	SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindAll(activityID, game)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.OkWithData(ctx, gamesModel)
}

// GetVoteOptionReset 投票分數紀錄重新計算 GET API
func (h *Handler) GetVoteOptionReset(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// prefix = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		// optionID = ctx.Query("option_id")
		data interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	// 重新計算所有選項投票紀錄
	err := models.DefaultGameVoteOptionModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Reset(activityID, gameID)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 重新計算所有投票分數紀錄發生問題",
		})
		return
	}

	response.OkWithData(ctx, data)
}

// GetVoteOptionList 選項名單json資料 GET API
func (h *Handler) GetVoteOptionList(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// prefix = ctx.Param("__prefix")
		// activityID = ctx.Query("activity_id")
		gameID   = ctx.Query("game_id")
		optionID = ctx.Query("option_id")
		data     interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	gameModel, err := models.DefaultGameVoteOptionListModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Find(false, gameID, optionID, "")
	if err != nil || gameID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得投票選項資訊，請輸入有效的遊戲ID",
		})
		return
	}

	data = gameModel

	response.OkWithData(ctx, data)
}

// GetVoteSpecialOfficer 投票特殊人員json資料 GET API
func (h *Handler) GetVoteSpecialOfficer(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// prefix = ctx.Param("__prefix")
		// activityID = ctx.Query("activity_id")
		gameID = ctx.Query("game_id")
		userID = ctx.Query("user_id")
		data   interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	gameModel, err := models.DefaultGameVoteSpecialOfficerModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Find(false, gameID, userID)
	if err != nil || gameID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得投票特殊人員資訊，請輸入有效的遊戲ID",
		})
		return
	}

	data = gameModel

	response.OkWithData(ctx, data)
}

// GetVoteOption 投票選項json資料 GET API
func (h *Handler) GetVoteOption(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// prefix = ctx.Param("__prefix")
		// activityID = ctx.Query("activity_id")
		gameID = ctx.Query("game_id")
		data   interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	gameModel, err := models.DefaultGameVoteOptionModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindOrderByID(false, gameID)
	if err != nil || gameID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得投票選項資訊，請輸入有效的遊戲ID",
		})
		return
	}

	data = gameModel

	response.OkWithData(ctx, data)
}

// GetSignGameInfo 遊戲設置json資料(單個遊戲資訊查詢，sign分類下) GET API
func (h *Handler) GetSignGameInfo(ctx *gin.Context) {
	var (
		host   = ctx.Request.Host
		prefix = ctx.Param("__prefix")
		// activityID = ctx.Query("activity_id")
		gameID    = ctx.Query("game_id")
		data      interface{}
		gameModel models.GameModel
		isRedis   bool
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if prefix == "undefined" {
		// 避免錯誤LOG紀錄一直出現
		return
	}

	// 判斷是否使用redis
	if ctx.Query("isredis") == "true" {
		isRedis = true
	}

	// 單個場次
	gameModel, err := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Find(isRedis, gameID, prefix)
	if err != nil || gameModel.ID == 0 || gameID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得遊戲資訊，請輸入有效的遊戲ID",
		})
		return
	}

	// 投票遊戲
	if prefix == "vote" {
		// 判斷遊戲狀態
		var gameStatus string
		// 判斷投票結束時間
		now, _ := time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local) // 目前時間
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteStartTime, time.Local)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gameModel.VoteEndTime, time.Local)

		// 比較時間，判斷遊戲狀態
		if now.Before(startTime) { // 現在時間在開始時間之前
			gameStatus = "close" // 關閉
		} else if now.Before(endTime) { // 現在時間在截止時間之前
			gameStatus = "gaming" // 遊戲中
		} else { // 現在時間在截止時間之後
			gameStatus = "calculate" // 結算狀態
		}

		if gameStatus != gameModel.GameStatus {
			// 狀態不一樣，更新為最新狀態
			gameModel.GameStatus = gameStatus

			// 更新遊戲狀態(mysql)
			// if err := db.Table(config.ACTIVITY_GAME_TABLE).WithConn(h.dbConn).
			// 	Where("game_id", "=", gameModel.GameID).
			// 	Update(command.Value{"game_status": gameStatus}); err != nil &&
			// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			// 		UserID:  "",
			// 		Method:  ctx.Request.Method,
			// 		Path:    ctx.Request.URL.Path,
			// 		Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
			// 	})
			// 	return
			// }

			// 更新遊戲狀態(mongo，activity_game)
			filter := bson.M{"game_id": gameModel.GameID} // 過濾條件
			// 要更新的值
			update := bson.M{
				"$set": bson.M{"game_status": gameStatus},
				// "$unset": bson.M{},                // 移除不需要的欄位
				// "$inc":   bson.M{"edit_times": 1}, // edit 欄位遞增 1
			}

			if _, err := h.mongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 更新遊戲狀態發生問題，請重新操作",
				})
			}

			// 修改redis中的遊戲資訊
			h.redisConn.HashSetCache(config.GAME_REDIS+gameID, "game_status", gameStatus)

			// 設置過期時間
			// h.redisConn.SetExpire(config.GAME_REDIS+gameID, config.REDIS_EXPIRE)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "遊戲狀態改變")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		}
	}

	data = gameModel

	response.OkWithData(ctx, data)
}

// GetGameInfo 遊戲設置json資料(基本設置、單個遊戲資訊查詢) GET API
func (h *Handler) GetGameInfo(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		data       interface{}
		isRedis    bool
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	if prefix == "undefined" || gameID == "undefined" {
		// 避免錯誤LOG紀錄一直出現
		return
	}

	if prefix == "setting" {
		// 遊戲基本設置
		gameSettingModel, err := models.DefaultGameSettingModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Find(activityID)
		if err != nil || gameSettingModel.ID == 0 || activityID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得遊戲資訊，請輸入有效的活動ID",
			})
			return
		}

		data = gameSettingModel
	} else {
		// 判斷是否使用redis
		if ctx.Query("isredis") == "true" {
			isRedis = true
		}

		// 單個場次
		gameModel, err := models.DefaultGameModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			Find(isRedis, gameID, prefix)
		if err != nil || gameModel.ID == 0 || gameID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得遊戲資訊，請輸入有效的遊戲ID",
			})
			return
		}

		data = gameModel
	}
	response.OkWithData(ctx, data)
}

// GetSignPrize 獎品設置json資料(sign分類下) GET API
func (h *Handler) GetSignPrize(ctx *gin.Context) {
	var (
		host   = ctx.Request.Host
		gameID = ctx.Query("game_id")
		data   interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	prizes, err := models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindPrizes(false, gameID)
	if err != nil || gameID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識獎品資訊，請輸入有效的遊戲ID",
		})
		return
	}

	data = prizes

	response.OkWithData(ctx, data)
}

// GetPrize 獎品設置json資料(該遊戲所有獎品) GET API
func (h *Handler) GetPrize(ctx *gin.Context) {
	var (
		host    = ctx.Request.Host
		prefix  = ctx.Param("__prefix")
		gameID  = ctx.Query("game_id")
		isAll   = ctx.Query("isall")
		isArray = ctx.Query("is_array")
		data    interface{}
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	// 搖號抽獎(只取得數量大於0的獎品)
	if prefix == "draw_numbers" && isAll == "" {
		prizes, err := models.DefaultPrizeModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			FindExistPrizes(prefix, gameID)
		if err != nil || gameID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識獎品資訊，請輸入有效的遊戲ID",
			})
			return
		}

		// 判斷有沒有獎品
		if len(prizes) == 0 {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 獎品數為0，無法進入遊戲，請管理員設置獎品後才能開始遊戲",
			})
			return
		}
		data = prizes

		if isArray == "true" {
			response.OkWithDatas(ctx, data)
			return
		}
	} else {
		prizes, err := models.DefaultPrizeModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			FindPrizes(false, gameID)
		if err != nil || gameID == "" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識獎品資訊，請輸入有效的遊戲ID",
			})
			return
		}

		// 幸運抽獎(獎品未滿八個，補未中獎)
		if prefix == "lottery" {
			amount := 8 - len(prizes)
			if amount != 0 {
				for i := 0; i < amount; i++ {
					prizes = append(prizes, models.PrizeModel{
						PrizeID:      "",
						PrizeName:    "未中獎",
						PrizePicture: "/admin/uploads/system/cry.png",
					})
				}
			}
		}
		data = prizes
	}

	response.OkWithData(ctx, data)
}

// GetStaffManage 人員管理json資料(參加遊戲人員、中獎人員、黑名單) GET API
func (h *Handler) GetStaffManage(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		optionID   = ctx.Query("option_id")
		game       = ctx.Query("game")
		userID     = ctx.Query("user_id")
		round      = ctx.Query("round")
		status     = ctx.Query("status")
		limit      = ctx.Query("limit")
		offset     = ctx.Query("offset")
		data       interface{}
		err        error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識活動資訊，請輸入有效的活動ID",
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	if offsetInt != 0 && limitInt == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識參數資訊，請輸入有效的limit.offset參數",
		})
		return
	}

	if prefix == "attend" {
		gameStaffModel, err := models.DefaultGameStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			FindAll(activityID, gameID, userID, game, round,
				int64(limitInt), int64(offsetInt))
		if err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
		data = gameStaffModel
	} else if prefix == "winning" {
		prizeStaffModel, err := models.DefaultPrizeStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			FindAll(activityID, gameID, userID, game, round, status, "",
				int64(limitInt), int64(offsetInt))
		if err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
		data = prizeStaffModel
	} else if prefix == "black" {
		// if game == "" {
		// 	response.BadRequest(ctx, "錯誤: 無法game資訊，請輸入有效的game參數")
		// 	return
		// }

		blackStaffModel, err := models.DefaultBlackStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			FindAll(false, activityID, gameID, game)
		if err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
		data = blackStaffModel
	} else if prefix == "record" {
		if game != "vote" {
			// 其他競技遊戲紀錄
			scoreStaffModel, err := models.DefaultScoreModel().
				SetConn(h.dbConn, h.redisConn,h.mongoConn).
				FindAll(activityID, gameID, userID, game, round,
					int64(limitInt), int64(offsetInt))
			if err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: err.Error(),
				})
				return
			}
			data = scoreStaffModel
		} else if game == "vote" {
			// 投票紀錄
			records, err := models.DefaultGameVoteRecordModel().
				SetConn(h.dbConn, h.redisConn,h.mongoConn).
				FindAll(activityID, gameID, optionID, round, int64(limitInt), int64(offsetInt))
			if err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 無法取得投票紀錄資訊，請輸入有效的遊戲ID",
				})
				return
			}

			data = records
		}
	}
	// else if prefix == "pk" {
	// 	pkStaffModel, err := models.DefaultPKStaffModel().SetDbConn(h.dbConn).
	// 		FindAll(activityID, gameID, userID, game, round,
	// 			int64(limitInt), int64(offsetInt))
	// 	if err != nil {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: err.Error(),
	// 		})
	// 		return
	// 	}
	// 	data = pkStaffModel
	// }

	response.OkWithData(ctx, data)
}

// @Summary 聊天紀錄JSON資料
// @Tags Chatroom_Record
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /chatroom/record [get]
func (h *Handler) ChatroomRecordsJSON(ctx *gin.Context) {
}

// @Summary 提問紀錄JSON資料
// @Tags Question_Record
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /question/record [get]
func (h *Handler) QuestionRecordsJSON(ctx *gin.Context) {
}

// @Summary 所有自定義場景json資料
// @Tags Customize_Scene
// @version 1.0
// @Accept  json
// @param user_id query string true "用戶ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /customize_scene [get]
func (h *Handler) CustomizeScenesJSON(ctx *gin.Context) {
}

// func (h *Handler) TestApi2(ctx *gin.Context) {
// 	var (
// 		host = ctx.Request.Host
// 		// activityID = ctx.Query("activity_id")
// 		gameID = ctx.Query("game_id")
// 		// userID = ctx.Query("user_id")
// 		// wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	)
// 	if host != config.API_URL {
// 		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  "",
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: "錯誤: 網域請求發生問題",
// 		})
// 		return
// 	}

// 	// 判斷玩家是否已加入遊戲
// 	// if models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	// 	IsUserExist(gameID, userID, 200) {

// 	// 遞減隊伍人數(資料庫.redis)
// 	models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 		UpdateTeamPeople(true, gameID, "left_team", -1, -1)
// 	// }

// 	response.Ok(ctx)
// }

// func (h *Handler) TestApi1(ctx *gin.Context) {
// 	var (
// 		host       = ctx.Request.Host
// 		activityID = ctx.Query("activity_id")
// 		// gameID = ctx.Query("game_id")
// 		userID = ctx.Query("user_id")
// 		// wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
// 	)
// 	if host != config.API_URL {
// 		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  "",
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: "錯誤: 網域請求發生問題",
// 		})
// 		return
// 	}

// 	// 用戶遊戲紀錄
// 	allRecords, _ := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 		SetRedisConn(h.redisConn).FindUserGameRecords(true, activityID, userID)
// 	log.Println("資料長度: ", len(allRecords))

// 	// log.Println("刪除0")
// 	// wg.Add(int(500)) //計數器

// 	// for i := 1; i <= 500; i++ {
// 	// go func(i int) {
// 	// defer wg.Done()
// 	// userID := strconv.Itoa(int(i))

// 	// 判斷玩家是否已加入遊戲
// 	// if models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	// IsUserExist(gameID, userID, 200) {
// 	// 刪除加入遊戲人員資料
// 	// models.DefaultGameStaffModel().SetDbConn(h.dbConn).DeleteUser(gameID, userID, 200)
// 	// }

// 	// }(i)
// 	// }
// 	// wg.Wait() //等待計數器歸0

// 	// log.Println("結束")

// 	response.Ok(ctx)
// }

// GetApplysignUser 自定義匯入報名簽到人員json資料 GET API
// func (h *Handler) GetApplysignUser(ctx *gin.Context) {
// 	var (
// 		host       = ctx.Request.Host
// 		activityID = ctx.Query("activity_id")
// 		users      = make([]models.ApplysignUserModel, 0)
// 		err        error
// 	)
// 	if host != config.API_URL {
// 		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  "",
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: "錯誤: 網域請求發生問題",
// 		})
// 		return
// 	}

// 	if activityID == "" {
// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  "",
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: "錯誤: 無法辨識參數，請輸入有效的參數",
// 		})
// 		return
// 	}

// 	// 查詢該活動所有自定義匯入的報名簽到人員資料
// 	users, err = models.DefaultApplysignUserModel().
// 		SetDbConn(h.dbConn).FindAll(activityID)
// 	if err != nil {
// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
// 			UserID:  "",
// 			Method:  ctx.Request.Method,
// 			Path:    ctx.Request.URL.Path,
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	response.OkWithData(ctx, users)
// }

// @@@Summary 自定義匯入報名簽到人員JSON資料
// @@@Tags ApplySign
// @@@version 1.0
// @@@Accept  json
// @@@param activity_id query string false "活動ID"
// @@@Success 200 {array} response.ResponseWithData
// @@@Failure 404 {array} response.ResponseBadRequest
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /applysign/user [get]
func (h *Handler) ApplySignUserJSON(ctx *gin.Context) {
}

// if prefix == "lottery" {
// 	prizes, err := models.DefaultPrizeModel().SetDbConn(h.dbConn).
// 		FindPrizes(false, gameID)
// 	if er
