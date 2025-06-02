package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"hilive/modules/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ShowAuthUser 報名簽到用戶補齊資料頁面 GET API
func (h *Handler) ShowAuthUser(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		userID     = ctx.Query("user_id")
		sign       = ctx.Query("sign") // 簽到審核開關
		isFirst    = ctx.Query("isfirst")
		htmlTmpl   = "./hilives/hilive/views/chatroom/custom_form.html"
		err        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" || userID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、用戶資訊，請輸入有效的活動、用戶ID")
		return
	}

	// 報名簽到、自定義、用戶資訊
	applysignModel, err := models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(0, activityID, userID, false)
	if err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得報名簽到人員資料，請重新報名簽到")
		return
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		ActivityID:             activityID,
		GameID:                 gameID,
		IsSign:                 sign == "open",
		IsFirst:                isFirst == "true",
		ApplysignCustomizeJSON: utils.JSON(applysignModel),
		Route: route{
			POST: config.AUTH_USER_API_URL,
		},
	})
}

// AuthUser 用戶補齊資料 POST API
// @Summary 用戶補齊個人資料
// @Tags Auth
// @version 1.0
// @param activity_id formData string true "活動ID"
// @param game_id formData string true "遊戲ID"
// @param user_id formData string true "用戶ID"
// @param name formData string false "name"
// @param avatar formData string false "avatar"
// @param avatar__default_flag formData string false "avatar__default_flag"
// @param phone formData string false "手機號碼(10碼)" minlength(10) maxlength(10)
// @param ext_email formData string false "電子信箱"
// @param ext_1 formData string false "ext_1"
// @param ext_2 formData string false "ext_2"
// @param ext_3 formData string false "ext_3"
// @param ext_4 formData string false "ext_4"
// @param ext_5 formData string false "ext_5"
// @param ext_6 formData string false "ext_6"
// @param ext_7 formData string false "ext_7"
// @param ext_8 formData string false "ext_8"
// @param ext_9 formData string false "ext_9"
// @param ext_10 formData string false "ext_10"
// @param sign formData string true "是否開啟簽到" Enums(true, false)
// @param isfirst formData string true "是否第一次執行LINE授權" Enums(true, false)
// @param host formData string true "網域" Enums(hilives.net, www.hilives.net, dev.hilives.net)
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /auth/user [post]
func (h *Handler) AuthUser(ctx *gin.Context) {
	var (
		userID     = ctx.Request.FormValue("user_id")
		phone      = ctx.Request.FormValue("phone")
		email      = ctx.Request.FormValue("ext_email")
		sign       = ctx.Request.FormValue("sign")
		isFirst    = ctx.Request.FormValue("isfirst")
		activityID = ctx.Request.FormValue("activity_id")
		gameID     = ctx.Request.FormValue("game_id")
		// name       = ctx.Request.FormValue("name")
		// avatar     = ctx.Request.FormValue("avatar")
		// avatarDefault = ctx.Request.FormValue("avatar__default_flag")
		values = []string{ctx.Request.FormValue("ext_1"), ctx.Request.FormValue("ext_2"),
			ctx.Request.FormValue("ext_3"), ctx.Request.FormValue("ext_4"), ctx.Request.FormValue("ext_5"),
			ctx.Request.FormValue("ext_6"), ctx.Request.FormValue("ext_7"), ctx.Request.FormValue("ext_8"),
			ctx.Request.FormValue("ext_9"), ctx.Request.FormValue("ext_10")}
		redirectURL = fmt.Sprintf(config.APPLYSIGN_URL, activityID, userID, gameID)
		host        = ctx.Request.FormValue("host")
		now, _      = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
		status string
	)
	if ctx.Request.Host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" || userID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識活動、用戶資訊，請輸入有效的活動、用戶ID",
		})
		return
	}

	// 是否開啟簽到
	if sign == "true" {
		redirectURL += "&sign=open"
	}
	// 是否第一次執行LINE授權
	if isFirst == "true" {
		redirectURL += "&isfirst=true"
	}

	// 檢查電話與電子信箱
	if phone != "" {
		if len(phone) > 2 {
			if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
				})
				return
			}
		} else {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
			})
			return
		}
	}
	if email != "" && !strings.Contains(email, "@") {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址",
		})
		return
	}

	// 判斷信箱電話欄位是否為空
	if phone != "" || email != "" {
		// 更新電話以及自定義電子信箱資訊(有數值才更新)
		if err := models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdatePhoneAndEmail(models.LineModel{
				UserID: userID, Phone: phone, ExtEmail: email}); err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	}

	// 判斷是否上傳圖片
	// 處理表單參數
	// param, err := ctx.MultipartForm()
	// if err != nil {
	// 	response.Error(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  userID,
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 表單參數發生問題，請重新操作",
	// 	})
	// 	return
	// }

	// 上傳圖片、檔案
	// if len(param.File) > 0 {
	// 	// 有上傳檔案
	// 	if err := file.GetFileEngine(config.FILE_ENGINE).Upload(
	// 		ctx.Request.MultipartForm, ctx.Request.URL.Path,
	// 		userID, "", "", ""); err != nil {
	// 		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  userID,
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 上傳檔案發生問題，請重新上傳檔案",
	// 		})
	// 		return
	// 	}
	// 	avatar = config.HTTP_HILIVES_NET_URL + "/admin/uploads/applysign_user/" + param.Value["avatar"][0]
	// }
	// if avatarDefault == "1" {
	// 	// 未上傳頭像資料，使用預設的頭像
	// 	avatar = config.HTTP_HILIVES_NET_URL + "/admin/uploads/system/img-user-pic.png"
	// }

	// 判斷姓名頭像欄位是否為空
	// if name != "" || avatar != "" {
	// 	// 更新資料
	// 	if err := models.DefaultApplysignUserModel().SetDbConn(h.dbConn).
	// 		Update(models.ApplysignUserModel{
	// 			ActivityID: activityID,
	// 			UserID:     userID,
	// 			Name:       name,
	// 			Avatar:     avatar}); err != nil {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  userID,
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: err.Error(),
	// 		})
	// 		return
	// 	}
	// }

	// 判斷redis裡是否有用戶資訊，有則修改name.avatar.phone資料
	if value, _ := h.redisConn.HashGetCache(config.AUTH_USERS_REDIS+userID, "user_id"); value != "" {
		params := []interface{}{config.AUTH_USERS_REDIS + userID}
		fields := []string{
			// "name", "avatar",
			"phone", "ext_email",
		}
		uservalues := []string{
			// name, avatar,
			phone, email,
		}

		for i, value := range uservalues {
			if value != "" {
				params = append(params, fields[i], value)
			}
		}

		// 更新redis中的用戶資料
		if len(params) > 1 {
			h.redisConn.HashMultiSetCache(params)

			// 設置過期時間
			// h.redisConn.SetExpire(config.AUTH_USERS_REDIS+userID,config.REDIS_EXPIRE)
		}
	}

	// 更新報名簽到人員ext欄位資料
	if err := models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		UpdateExt(
			activityID, userID, values); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}
	// 報名簽到人員資料
	applysign, err := models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(0, activityID, userID, false)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得報名簽到人員資料，請重新報名簽到",
		})
		return
	}

	if applysign.Status != "sign" {
		// 補齊資料後將用戶狀態改為審核中或報名成功
		// 取得活動狀態資訊
		activityModel, err := h.getActivityInfo(false, activityID)
		if err != nil || activityModel.ID == 0 {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得活動資訊，請重新查詢",
			})
			return
		}
		// 報名審核開啟
		if activityModel.ApplyCheck == "open" {
			status = "review"
		} else {
			// 報名審核關閉
			status = "apply"
		}

		if err = models.DefaultApplysignModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			UpdateStatus(true, host, models.EditApplysignModel{
				ActivityID: activityID,
				LineUsers:  []string{userID},
				ReviewTime: now.Format("2006-01-02 15:04:05"),
				Status:     status,
			}, false); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		}
	}

	response.OkWithURL(ctx, redirectURL)
}

// host       = ctx.Request.FormValue("host")

// customizeModel, err := models.DefaultCustomizeModel().SetDbConn(h.dbConn).FindToModel(activityID)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得自定義欄位資訊，請重新整理頁面")
// 	return
// }
// LINE用戶資訊
// lineModel, err := models.DefaultLineModel().SetDbConn(h.dbConn).Find("identify", userID)
// if err != nil || lineModel.ID == 0 {
// 	response.BadRequest(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
// 	return
// }
// 報名簽到人員資料
// applysign, err := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	FindByUserID(activityID, lineModel.UserID)
// applysign, err := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	LeftJoinByIdentify(activityID, userID)
// if err != nil {
// 	h.executeErrorHTML(ctx, "錯誤: 無法取得報名簽到人員資料，請重新報名簽到")
// 	return
// }

// LINE用戶資訊
// lineModel, err := models.DefaultLineModel().SetDbConn(h.dbConn).Find("identify", userID)
// if err != nil || lineModel.ID == 0 {
// 	response.BadRequest(ctx, "錯誤: 無法辨識用戶資訊，請輸入有效的用戶ID")
// 	return
//
