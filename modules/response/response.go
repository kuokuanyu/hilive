package response

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"ok"`
}

type ResponseBadRequest struct {
	Code    int    `json:"code" example:"400"`
	Role    string `json:"role" example:"host"`
	Message string `json:"message" example:"error message"`
}

type ResponseInternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"error message"`
}

type ResponseWithData struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"message"`
	Data    interface{} `json:"data"`
	Datas   interface{} `json:"datas"`
}

type ResponseWithURL struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"message"`
	URL     string `json:"url" example:"redirect URL"`
}

type ResponseDeleteWithURL struct {
	URL              string `json:"url"`
	ConfirmationCode string `json:"confirmation_code"`
}

// CORS 跨來源資源共用
func CORS(ctx *gin.Context) {
	method := ctx.Request.Method
	origin := ctx.Request.Header.Get("Origin")
	// log.Println("執行CORS origin: ", origin)
	if origin != "" {
		// log.Println("執行CORS origin: ", origin, ctx.Request.URL.Path, ctx.Request.Host)
		if strings.Contains(origin, config.HILIVES_NET_URL) {
			// log.Println("執行CORS", origin, ctx.Request.URL.Path)
			// ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, DELETE, GET, PUT, OPTIONS")
		}
	}
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	// ctx.Next()
}

// Ok 回傳成功
func Ok(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "ok",
	})
}

// OkWithData 回傳成功以及data
func OkWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, ResponseWithData{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})
}

// OkWithDatas 回傳成功以及datas
func OkWithDatas(ctx *gin.Context, datas interface{}) {
	ctx.JSON(http.StatusOK, ResponseWithData{
		Code:    http.StatusOK,
		Message: "ok",
		Datas:   datas,
	})
}

// OkWithURL 回傳成功以及url
func OkWithURL(ctx *gin.Context, url string) {
	ctx.JSON(http.StatusOK, ResponseWithURL{
		Code:    http.StatusOK,
		Message: "ok",
		URL:     url,
	})
}

// OkDeleteWithURL 回傳url及user_id
func OkDeleteWithURL(ctx *gin.Context, userID, url string) {
	ctx.JSON(http.StatusOK, ResponseDeleteWithURL{
		URL:              url,
		ConfirmationCode: userID,
	})
}

// OkWithMsg 回傳成功以及msg
// func OkWithMsg(ctx *gin.Context, msg string) {
// 	ctx.JSON(http.StatusOK, map[string]interface{}{
// 		"code": http.StatusOK,
// 		"msg":  msg,
// 	})
// }

// BadRequest 400錯誤
func BadRequest(ctx *gin.Context, conn db.Connection, model models.EditErrorLogModel) {
	// logger.LoggerError(ctx, msg)

	if model.UserID == "" {
		// 如果user_id為空則紀錄ip
		model.UserID = utils.ClientIP(ctx.Request)
	}

	models.DefaultErrorLogModel().
		SetConn(conn, nil, nil).
		Add(models.EditErrorLogModel{
			UserID:    model.UserID,
			Code:      http.StatusBadRequest,
			Method:    ctx.Request.Method,
			Path:      ctx.Request.URL.Path,
			Message:   model.Message,
			PathQuery: ctx.Request.URL.RawQuery,
		})

	ctx.JSON(http.StatusBadRequest, ResponseBadRequest{
		Code:    http.StatusBadRequest,
		Message: model.Message,
	})
}

// BadRequestWithRole 400錯誤(帶有角色參數)
func BadRequestWithRole(ctx *gin.Context, role string, conn db.Connection, model models.EditErrorLogModel) {
	// logger.LoggerError(ctx, msg)

	if model.UserID == "" {
		// 如果user_id為空則紀錄ip
		model.UserID = utils.ClientIP(ctx.Request)
	}

	// 判斷錯誤訊息是否需寫入資料表(error log)
	if model.Message != "錯誤: 主持人未開啟遊戲，無法參加遊戲，請等待主持人開放該場次遊戲" &&
		model.Message != "錯誤: 玩家未完成報名簽到，請重新報名簽到" &&
		model.Message != "錯誤: 其他玩家正在遊戲中，無法參加遊戲，請等待遊戲結束" &&
		model.Message != "錯誤: 獎品數為0，無法進入遊戲，請管理員設置獎品後才能開始遊戲" &&
		model.Message != "錯誤: 獎品數量必須大於每輪發獎人數，請重新設置" &&
		model.Message != "錯誤: 該場次遊戲已被開啟，無法重複開啟遊戲" &&
		model.Message != "錯誤: 該場次抽獎次數已達上限，無法參加抽獎" &&
		model.Message != "錯誤: 玩家剩餘投票次數為0，請等待遊戲結束" &&
		model.Message != "錯誤: 用戶為黑名單人員(活動、遊戲黑名單)，無法參加遊戲。如有疑問，請聯繫主辦單位" &&
		model.Message != "錯誤: 用戶為黑名單人員(活動、簽名牆黑名單)，無法參加遊戲。如有疑問，請聯繫主辦單位" {
		models.DefaultErrorLogModel().
			SetConn(conn, nil, nil).
			Add(models.EditErrorLogModel{
				UserID:    model.UserID,
				Code:      http.StatusBadRequest,
				Method:    ctx.Request.Method,
				Path:      ctx.Request.URL.Path,
				Message:   model.Message,
				PathQuery: ctx.Request.URL.RawQuery,
			})
	}

	ctx.JSON(http.StatusBadRequest, ResponseBadRequest{
		Code:    http.StatusBadRequest,
		Role:    role,
		Message: model.Message,
	})
}

// Error 回傳code:500 and msg
func Error(ctx *gin.Context, conn db.Connection, model models.EditErrorLogModel) {
	// logger.LoggerError(ctx, msg)

	if model.UserID == "" {
		// 如果user_id為空則紀錄ip
		model.UserID = utils.ClientIP(ctx.Request)
	}

	// 判斷錯誤訊息是否需寫入資料表(error log)
	if model.Message != "錯誤: 活動已結束，謝謝參與" &&
		model.Message != "錯誤: 未將官方帳號加為好友，請將@hilives加入好友" &&
		model.Message != "錯誤: 無法辨識活動資訊，請輸入有效的活動ID" &&
		model.Message != "錯誤: 取得用戶資料發生問題，請重新掃描QRcode進行驗證(錯誤or被封鎖)" {
		models.DefaultErrorLogModel().
			SetConn(conn, nil, nil).
			Add(models.EditErrorLogModel{
				UserID:    model.UserID,
				Code:      http.StatusInternalServerError,
				Method:    ctx.Request.Method,
				Path:      ctx.Request.URL.Path,
				Message:   model.Message,
				PathQuery: ctx.Request.URL.RawQuery,
			})
	}

	ctx.JSON(http.StatusInternalServerError, ResponseInternalServerError{
		Code:    http.StatusInternalServerError,
		Message: model.Message,
	})
}
