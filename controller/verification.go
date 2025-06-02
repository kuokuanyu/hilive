package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"net/http"
	"time"

	"github.com/twilio/twilio-go"

	"github.com/gin-gonic/gin"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	client = twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username:   config.ACCOUNT_SID,
			Password:   config.AUTH_TOKEN,
			AccountSid: config.ACCOUNT_SID,
		})
)

// Verification 發送手機驗證碼 POST API
// @Summary 發送手機驗證碼
// @Tags Verification
// @version 1.0
// @param phone formData string true "手機號碼(10碼)" minlength(10) maxlength(10)
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /verification [post]
func (h *Handler) Verification(ctx *gin.Context) {
	var (
		host  = ctx.Request.Host
		phone = ctx.Request.FormValue("phone")
		now   = time.Now().Format("2006-01-02") // 當日日期
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

	// 取得phone資料
	phoneModel, err := models.DefaultPhoneModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(phone)
	if err != nil {
		// logger.LoggerError(ctx, err.Error())
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	} else if phoneModel.ID != 0 {
		// 判斷發送訊息上限
		var (
			date  = phoneModel.SendTime
			times = phoneModel.Times // 發送次數
		)

		// 判斷發送訊息日期
		if date == now && times >= 3 {
			// 發送日期=今日日期且發送次數已達3次
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 電話驗證次數已達單日上限三次",
			})
			return
		}
	}

	// 更新發送次數資料(遞增)
	if err = models.DefaultPhoneModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Update(
			models.EditPhoneModel{
				Phone:    phone,
				SendTime: now,
				Times:    "times + 1",
			},
		); err != nil {
		// logger.LoggerError(ctx, err.Error())
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// 發送驗證碼
	// fmt.Println("電話: ", "+886"+phone[1:])
	params := &verify.CreateVerificationParams{}
	params.SetTo("+886" + phone[1:])
	params.SetChannel("sms")

	_, err = client.VerifyV2.CreateVerification(config.SERVICE_SID, params)
	if err != nil {
		fmt.Println("錯誤: ", err)
		// logger.LoggerError(ctx, "錯誤: 發送驗證碼發生問題，請重新操作")
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 發送驗證碼發生問題，請重新操作",
		})
		return
	}

	response.Ok(ctx)
	return
}

// VerificationCheck 檢查手機驗證碼 POST API
// @Summary 檢查手機驗證碼
// @Tags Verification
// @version 1.0
// @param phone formData string true "手機號碼(10碼)" minlength(10) maxlength(10)
// @param code formData string true "驗證碼(6碼)" minlength(6) maxlength(6)
// @Success 200 {array} response.ResponseWithURL
// @Failure 400 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /verification/check [post]
func (h *Handler) VerificationCheck(ctx *gin.Context) {
	var (
		host  = ctx.Request.Host
		phone = ctx.Request.FormValue("phone")
		code  = ctx.Request.FormValue("code")
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

	// 檢查驗證碼
	// fmt.Println("電話: ", "+886"+phone[1:])
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo("+886" + phone[1:])
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(config.SERVICE_SID, params)
	if err != nil {
		// logger.LoggerError(ctx, "錯誤: 檢查驗證碼發生問題，請重新輸入")
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 檢查驗證碼發生問題，請重新輸入",
		})
		return
	} else if resp.Status != nil && *resp.Status == "pending" {
		// logger.LoggerError(ctx, "錯誤: 驗證碼錯誤，請重新輸入")
		// 驗證回傳pending(失敗) or approved(成功)
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 驗證碼錯誤，請重新輸入",
		})
		return
	}

	// 完成手機驗證，更新驗證狀態資料(status:yes)
	if err = models.DefaultPhoneModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Update(
			models.EditPhoneModel{
				Phone:  phone,
				Status: "yes",
			},
		); err != nil {
		// logger.LoggerError(ctx, err.Error())
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	response.Ok(ctx)
	return
}

func (h *Handler) VerificationTest(ctx *gin.Context) {
	// https://dev.hilives.net/applysign?activity_id=Bfon6SaV6ORhmuDQUioI&user_id=Bfon6SaV6ORhmuDQUioI_tAjOCwMPI7vllBnm&sign=open
	// 郭 您好: 您已報名 測試傳送訊息close123 活動，該活動的驗證碼為asMDMM，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): https://dev.hilives.net/applysign?activity_id=Bfon6SaV6ORhmuDQUioI&user_id=Bfon6SaV6ORhmuDQUioI_tAjOCwMPI7vllBnm&sign=open
	// https://liff.line.me/1656920628-zJOEMMRl?activity_id=%s
	message := "郭 您好: 您已報名 測試傳送訊息close123 活動，該活動的驗證碼為asMDMM，可透過驗證碼進行簽到，或利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): https://dev.hilives.nnn/applysign?activity_id=Bfon6SaV6ORhmuDQUioI&user_id=Bfon6SaV6ORhmuDQUioI_tAjOCwMPI7vllBnm&sign=open"
	// 设置短信参数
	params := &openapi.CreateMessageParams{}
	params.SetTo("+886932530813")
	params.SetFrom("+18596966103") // 发送者的 Twilio 电话号码
	params.SetBody(message)

	// 发送短信
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		fmt.Println("Message sent successfully!")
		fmt.Println("Message SID: " + message)
	}

	ctx.JSON(http.StatusOK, message)
	return
}
