package controller

import (
	"errors"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// Register 註冊 POST API
// @Summary 註冊
// @Tags Auth
// @version 1.0
// @param name formData string true "用戶名稱(請設置20個字元以內)" minlength(1) maxlength(20)
// @param phone formData string true "手機號碼(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string true "電子信箱(電子郵件地址中必須包含「@」)"
// @param password formData string true "密碼"
// @param password_again formData string true "再次輸入密碼"
// @Success 200 {array} response.Response
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var (
		host          = ctx.Request.Host
		name          = ctx.Request.FormValue("name")
		phone         = ctx.Request.FormValue("phone")
		email         = ctx.Request.FormValue("email")
		password      = ctx.Request.FormValue("password")
		checkPassword = ctx.Request.FormValue("password_again")
		// hash          string
		err error
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

	// 檢查註冊條件
	if err = h.checkRegister(name, phone, email, password, checkPassword); err != nil {
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

	// 新增用戶
	// if err = h.add(ctx, models.EditUserModel{
	// 	// UserID:   utils.UUID(32),
	// 	Name:          name,
	// 	Phone:         phone,
	// 	Email:         email,
	// 	Password:      password,
	// 	PasswordAgain: checkPassword,
	// 	Avatar:        "",
	// 	Bind:          "no",
	// 	Cookie:        "no",
	// 	Permissions:   "3",
	// 	MaxActivity:   "1",
	// }); err != nil {
	// 	response.BadRequest(ctx, err.Error())
	// 	return
	// }

	response.Ok(ctx)
	return
}

// checkRegister 檢查註冊條件
func (h *Handler) checkRegister(name, phone, email, password, passwordAgain string) error {
	// var (
	// 	err error
	// )

	if name == "" || phone == "" || email == "" || password == "" || passwordAgain == "" {
		return errors.New("錯誤: 欄位資訊都不能為空，請輸入有效的欄位資訊")
	}

	// 檢查姓名
	if utf8.RuneCountInString(name) > 20 {
		return errors.New("錯誤: 用戶名稱不能為空並且不能超過20個字元，請輸入有效的用戶名稱")
	}

	// 檢查電話
	if len(phone) > 2 {
		if !strings.Contains(phone[:2], "09") || len(phone) != 10 {
			return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
		}
	} else {
		return errors.New("錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼")
	}

	// 檢查密碼
	if password != passwordAgain {
		return errors.New("錯誤: 輸入密碼不一致，請輸入有效的密碼")
	}

	// 檢查信箱
	if !strings.Contains(email, "@") {
		return errors.New("錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址")
	}

	// if hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
	// 	return "", errors.New("錯誤: 加密發生錯誤，請重新註冊用戶")
	// }

	if user, err := models.DefaultUserModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindPhoneAndEmail(phone, "or", email); user.Phone != "" || err != nil {
		return errors.New("錯誤: 電話號碼或電子郵件地址已被註冊過，請輸入有效的手機號碼與電子郵件地址")
	}

	// return string(hash[:]), nil
	return nil
}

// checkPhone 檢查電話資料
func (h *Handler) checkPhone(phone string) error {
	// 判斷user_phone是否存在電話資料
	phoneModel, err := models.DefaultPhoneModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(phone)
	if err != nil {
		return err
	}

	if phoneModel.ID != 0 {
		// 存在資料
		var (
			now  = time.Now().Format("2006-01-02") // 當日日期
			date = phoneModel.SendTime
			// dates    = strings.Split(date, "-")
			// year, _  = strconv.Atoi(dates[0]) // 年
			// month, _ = strconv.Atoi(dates[1]) // 月
			// day, _   = strconv.Atoi(dates[2]) // 日
			times = phoneModel.Times // 發送次數
		)

		// 判斷發送訊息日期
		if date == now && times >= 3 {
			// 發送日期=今日日期且發送次數已達3次
			return errors.New("錯誤: 電話驗證次數已達單日上限三次")
		} else if date != now {
			// 發送日期!=今日日期，更新資料日期(當日)及發送次數(0)
			if err = models.DefaultPhoneModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Update(
					models.EditPhoneModel{
						Phone:    phone,
						SendTime: now,
						Times:    "0",
					},
				); err != nil {
				return err
			}
		}
	} else {
		// 不存在資料，新增資料
		if err = models.DefaultPhoneModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Add(
				models.EditPhoneModel{
					Phone:  phone,
					Status: "no",
					Times:  "0",
				},
			); err != nil {
			return err
		}
	}

	return nil
}
