package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// EXPORT 匯出檔案 POST API
func (h *Handler) EXPORT(ctx *gin.Context) {
	var (
		host                                                    = ctx.Request.Host
		path                                                    = ctx.Request.URL.Path
		prefix                                                  = ctx.Param("__prefix")
		activityID                                              = ctx.Request.FormValue("activity_id")
		gameID                                                  = ctx.Request.FormValue("game_id")
		userID                                                  = ctx.Request.FormValue("user_id")
		optionID                                                = ctx.Request.FormValue("option_id")
		game                                                    = ctx.Request.FormValue("game")
		round                                                   = ctx.Request.FormValue("round")
		status                                                  = ctx.Request.FormValue("status")
		limit                                                   = ctx.Request.FormValue("limit")
		offset                                                  = ctx.Request.FormValue("offset")
		hostFormValue                                           = ctx.Request.FormValue("host")
		file                                                    *excelize.File
		limitInt, offsetInt                                     int
		fileName, gameCN, gameIDCN, userIDCN, roundCN, statusCN string
		err                                                     error
		// values        = make(map[string][]string)
		// userID, token string
		// err           error
	)

	// 網域判斷
	if host != config.API_URL {
		if strings.Contains(ctx.GetHeader("Accept"), "json") {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 網域請求發生問題",
			})
			return
		} else {
			h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
			return
		}
	}
	if activityID == "" {
		if strings.Contains(ctx.GetHeader("Accept"), "json") {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法辨識活動資訊，請輸入有效的活動ID",
			})
			return
		} else {
			h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
			return
		}
	}

	// 處理excel檔案名稱
	if game != "" {
		gameCN = game
	} else if game == "" {
		gameCN = "所有遊戲類型"
	}
	if gameID != "" {
		gameIDCN = gameID
	} else if gameID == "" {
		gameIDCN = "所有遊戲場次"
	}
	if game != "vote" {
		// 其他遊戲
		if userID != "" {
			userIDCN = userID
		} else if userID == "" {
			userIDCN = "所有用戶"
		}
	} else if game == "vote" {
		// 投票遊戲
		if optionID != "" {
			userIDCN = optionID
		} else if optionID == "" {
			userIDCN = "所有選項"
		}
	}

	// 判斷是否取得特定輪次資料
	if round != "" {
		roundCN = fmt.Sprintf("第%s輪", round)
	} else if round == "" {
		roundCN = "所有輪次"
	}
	// 判斷是否取得特定狀態資料
	if status != "" {
		if status == "yes" {
			statusCN = "已兌獎"
		} else if status == "no" {
			statusCN = "未兌獎"
		} else {
			// 報名簽到
			statusCN = status
		}
	} else if status == "" {
		statusCN = "所有狀態"
	}
	// 判斷是否取得特定資料數(limit.offset參數)
	limitInt, err = strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	offsetInt, err = strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	if offsetInt != 0 && limitInt == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識參數資訊，請輸入有效的limit.offset參數",
		})
		return
	}

	if path == "/v1/applysign/export" {
		// 匯出報名簽到資料
		fileName = "報名簽到-" + activityID + "-" + statusCN + ".xlsx"

		file, err = models.DefaultApplysignModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
			Export(activityID, status)
	} else if prefix == "attend" {
		fileName = "遊戲人員-" + activityID + "-" + gameCN + "-" +
			gameIDCN + "-" + userIDCN + "-" + roundCN + ".xlsx"

		file, err = models.DefaultGameStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			Export(activityID, gameID, game, userID, round,
				int64(limitInt), int64(offsetInt))
	} else if prefix == "winning" {
		fileName = "中獎人員-" + activityID + "-" + gameCN + "-" +
			gameIDCN + "-" + userIDCN + "-" + roundCN + "-" + statusCN + ".xlsx"

		file, err = models.DefaultPrizeStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			Export(activityID, gameID, game, userID, round, status,
				int64(limitInt), int64(offsetInt), hostFormValue)
	} else if prefix == "black" {
		fileName = "黑名單人員-" + activityID + "-" + gameCN + "-" + gameIDCN + ".xlsx"

		file, err = models.DefaultBlackStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
			Export(activityID, gameID, game)
	} else if prefix == "record" {
		fileName = "遊戲紀錄-" + activityID + "-" + gameCN + "-" +
			gameIDCN + "-" + userIDCN + "-" + roundCN + ".xlsx"

		if game != "vote" {
			// 其他遊戲
			file, err = models.DefaultScoreModel().
				SetConn(h.dbConn, h.redisConn,h.mongoConn).
				Export(activityID, gameID, game, userID, round,
					int64(limitInt), int64(offsetInt))
		} else if game == "vote" {
			// 投票遊戲
			file, err = models.DefaultGameVoteRecordModel().
				SetConn(h.dbConn, h.redisConn,h.mongoConn).
				Export(activityID, gameID, optionID, round,
					int64(limitInt), int64(offsetInt))
		}

	}
	if err != nil {
		if strings.Contains(ctx.GetHeader("Accept"), "json") {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: err.Error(),
			})
			return
		} else {
			h.executeErrorHTML(ctx, err.Error())
			return
		}
	}

	buf, err := file.WriteToBuffer()
	if err != nil {
		if strings.Contains(ctx.GetHeader("Accept"), "json") {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 匯出檔案發生問題",
			})
			return
		} else {
			h.executeErrorHTML(ctx, "錯誤: 匯出檔案發生問題")
			return
		}
	}

	ctx.Writer.Header().Set("content-disposition", `attachment; filename=`+fileName)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

// @Summary 匯出報名簽到人員excel檔案
// @Tags ApplySign
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param status query string false "報名簽到狀態(用,間隔)" Enums(sign, apply, review, cancel, refuse, no)
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /applysign/export [post]
func (h *Handler) EXPORTApplySign(ctx *gin.Context) {
}

// @Summary 匯出遊戲人員excel檔案
// @Tags Attend Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(redpack, ropepack, whack_mole, monopoly, QA, tugofwar)
// @param user_id query string false "用戶ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/attend/export [post]
func (h *Handler) EXPORTAttend(ctx *gin.Context) {
}

// @Summary 匯出中獎人員excel檔案
// @Tags Winning Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game formData string false "遊戲種類(用,間隔)" Enums(redpack, ropepack, whack_mole, monopoly, draw_numbers, lottery, QA, tugofwar)
// @param user_id query string false "用戶ID(用,間隔)"
// @param round formData integer false "遊戲輪次(用,間隔)"
// @param status query string false "兌獎狀態(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @param host formData string true "網域" Enums(hilives.net, www.hilives.net, dev.hilives.net)
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/winning/export [post]
func (h *Handler) EXPORTWinning(ctx *gin.Context) {
}

// @Summary 匯出黑名單人員excel檔案
// @Tags Black Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID"
// @param game query string false "遊戲" Enums(activity, message, question, redpack, ropepack, whack_mole, monopoly, draw_numbers, lottery, QA, tugofwar)
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/black/export [post]
func (h *Handler) EXPORTBlack(ctx *gin.Context) {
}

// @Summary 匯出遊戲紀錄excel檔案
// @Tags Record Staff
// @version 1.0
// @Accept  mpfd
// @param activity_id query string true "活動ID"
// @param game_id query string false "遊戲ID(用,間隔)"
// @param game query string false "遊戲種類(用,間隔)" Enums(whack_mole, monopoly, QA, tugofwar, bingo, vote)
// @param user_id query string false "用戶ID(用,間隔)"
// @param option_id query string false "選項ID(用,間隔)"
// @param round query string false "遊戲輪次(用,間隔)"
// @param limit query integer false "取得資料數"
// @param offset query integer false "跳過前n筆資料"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /staffmanage/record/export [post]
func (h *Handler) EXPORTRecord(ctx *gin.Context) {
}

// ctx.Header("Content-Type", "application/octet-stream")
// ctx.Header("Content-Disposition", "attachment; filename="+"123.xlsx")

// ctx.File("./" + "123.xlsx")

// file, _ := os.Open("userInputData.xlsx")
// defer file.Close()
// ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
// io.Copy(ctx.Writer, file)

// buf, err := f.WriteToBuffer()
// if err != nil {
// 	response.Error(ctx, "export error")
// 	return
// }

// ctx.Header("Content-Description", "File Transfer")
// ctx.Header("Content-Transfer-Encoding", "binary")
// ctx.Header("Content-Disposition", "attachment; filename=userInputData.xlsx")
// ctx.Header("Content-Type", "application/octet-stream")
// ctx.File("userInputData.xlsx")

// ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
// ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+"userInputData.xlsx")
// ctx.Writer.Header().Set("File-Name", "userInputData.xlsx")
// ctx.Writer.Header().Set("Content-Transfer-Encoding", "binary")
// ctx.Writer.Header().Set("Expires", "0")
// err := f.Write(ctx.Writer)
// fmt.Println(err)

// ctx.Writer.Header().Set("content-disposition", `attachment; filename=`+"Operation log-1673512426-page-1-pageSize-11.xlsx")
// ctx.Data(200, "application/octet-stream", buf.Bytes())
// fmt.Println("111111111111")
// // ctx.File("./downloads/file.xlsx")
// http.ServeFile(ctx.Writer, ctx.Request, "file")
// fmt.Println("22222")
// ctx.Data(http.StatusOK, "application/octet-stream", buf.Bytes())

// 該遊戲種類下所有資料
// if game == "redpack" {
// 	gameCN = "搖紅包"
// } else if game == "ropepack" {
// 	gameCN = "套紅包"
// } else if game == "whack_mole" {
// 	gameCN = "敲敲樂"
// } else if game == "monopoly" {
// 	gameCN = "抓偽鈔"
// } else if game == "draw_numbers" {
// 	gameCN = "搖號抽獎"
// } else if game == "lottery" {
// 	gameCN = "遊戲抽獎"
// } else if game == "QA" {
// 	gameCN = "快問快答"
// }

// if activityID != "" {
// 	if game != "" {
// 		// 該遊戲種類下所有資料
// 		if game == "redpack" {
// 			game = "搖紅包"
// 		} else if game == "ropepack" {
// 			game = "套紅包"
// 		} else if game == "whack_mole" {
// 			game = "敲敲樂"
// 		} else if game == "monopoly" {
// 			game = "抓偽鈔"
// 		} else if game == "draw_numbers" {
// 			game = "搖號抽獎"
// 		} else if game == "lottery" {
// 			game = "遊戲抽獎"
// 		} else if game == "QA" {
// 			game = "快問快答"
// 		}
// 		fileName = "中獎人員-" + activityID + "-" + game + "-" + round + ".xlsx"
// 	} else if game == "" {
// 		// 該活動下所有資料
// 		fileName = "中獎人員-" + activityID + "-" + round + ".xlsx"
// 	}
// } else if gameID != "" {
// 	fileName = "中獎人員-" + gameID + "-" + round + ".xlsx"
// }
