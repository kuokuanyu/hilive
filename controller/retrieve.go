package controller

import (
	"errors"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// Retrieve 找回密碼 POST API
// @Summary 找回密碼
// @Tags Auth
// @version 1.0
// @param phone formData string true "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string true "電子信箱(電子郵件地址中必須包含「@」)"
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /retrieve [post]
func (h *Handler) Retrieve(ctx *gin.Context) {
	var (
		host  = ctx.Request.Host
		phone = ctx.Request.FormValue("phone")
		email = ctx.Request.FormValue("email")
		err   error
	)
	if host != config.API_URL {
		// logger.LoggerError(ctx, "錯誤: 網域請求發生問題")
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}

	// 檢查忘記密碼條件
	if err = h.checkRetrieve(phone, email); err != nil {
		// logger.LoggerError(ctx, err.Error())
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// 檢查電話資料
	if err = h.checkPhone(phone); err != nil {
		// logger.LoggerError(ctx, err.Error())
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// 更新一組新密碼: admin
	// if err = user.Update(true, models.EditUserModel{
	// 	UserID:        user.UserID,
	// 	Password:      "admin",
	// 	PasswordAgain: "admin",
	// }, "users.user_id", user.UserID); err != nil {
	// 	response.BadRequest(ctx, "錯誤: 設置新密碼失敗，請重新操作")
	// 	return
	// }

	response.Ok(ctx)
	return
}

// checkRetrieve 檢查忘記密碼條件
func (h *Handler) checkRetrieve(phone, email string) error {
	// 判斷欄位是否為空
	if phone == "" || email == "" {
		return errors.New("錯誤: 手機號碼以及電子郵件地址欄位不能為空，請輸入有效的手機號碼以及電子郵件地址")
	}

	// 電話
	if len(phone) > 2 {
		if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	} else {
		return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
	}

	// 信箱
	if !strings.Contains(email, "@") {
		return errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	// 驗證用戶
	user, err := models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindPhoneAndEmail(phone, "and", email)
	if err != nil || user.Phone == "" {
		return errors.New("錯誤: 無法辨識用戶資訊，請輸入有效的手機號碼以及電子郵件地址")
	}

	return nil
}
