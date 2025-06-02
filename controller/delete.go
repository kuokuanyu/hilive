package controller

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// DELETE 刪除 DELETE API
func (h *Handler) DELETE(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// contentType = ctx.Request.Header.Get("content-type")
		path        = ctx.Request.URL.Path
		prefix      = ctx.Param("__prefix")
		id          = ctx.Request.FormValue("id")
		userID      = ctx.Request.FormValue("user")
		activityID  = ctx.Request.FormValue("activity_id")
		gameID      = ctx.Request.FormValue("game_id")
		tokenUserID = ctx.Request.FormValue("user_id")
		token       = ctx.Request.FormValue("token")
		// id, tokenUserID, userID, activityID, gameID, token string
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return

	}

	if path == "/v1/admin/manager" || path == "/v1/user" ||
		path == "/v1/activity" || strings.Contains(path, "/v1/interact/game") ||
		path == "/v1/interact/sign/vote/form" {
		// 因為管理員都可以幫忙設置, 所以需要區分user_id參數
		userID = ctx.Request.FormValue("user") // 該活動場次的管理員資料
	} else {
		userID = ctx.Request.FormValue("user_id")
	}

	if prefix == "" {
		if strings.Contains(path, "activity") {
			prefix = "activity"
		} else if strings.Contains(path, "/v1/applysign/user") {
			prefix = "applysign_user"
		} else if strings.Contains(path, "/v1/applysign") {
			prefix = "applysign"
		}
	} else if strings.Contains(path, "option_list") {
		prefix += "_option_list"
	} else if strings.Contains(path, "option") {
		prefix += "_option"
	} else if strings.Contains(path, "special_officer") {
		prefix += "_special_officer"
	} else if strings.Contains(path, "prize") {
		prefix += "_prize"
	} else if strings.Contains(path, "admin") {
		prefix = "admin_" + prefix
	}

	if !auth.CheckToken(token, tokenUserID) {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: Token驗證發生問題，請輸入有效的Token值",
		})

		return

	}

	table, _ := h.GetTable(ctx, prefix)
	if err := table.DeleteData(id, userID, activityID, gameID); err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})

		return

	}

	// 刪除redis資訊
	if strings.Contains(path, "prize") {
		// 獎品頁面
		if activityID == "" || gameID == "" {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得活動ID或遊戲ID參數",
			})
			return

		}

		// // 刪除獎品時，需要刪除遊戲獎品數量redis資料
		// h.redisConn.DelCache(config.GAME_PRIZES_AMOUNT_REDIS + gameID)

		// // 刪除獎品時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		// h.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

		// // 刪除玩家遊戲紀錄(中獎.未中獎)
		// h.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)
	} else if strings.Contains(path, "game") {
		// 遊戲頁面
		if activityID == "" {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得活動ID或遊戲ID參數",
			})
			return

		}

		// // 刪除遊戲場次時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
		// h.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

		// // 刪除玩家遊戲紀錄(中獎.未中獎)
		// h.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)
	} else if prefix == "winning" {
		if activityID == "" || gameID == "" {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 無法取得活動ID或遊戲ID參數",
			})
			return

		}

		// 中獎人員頁面
		gameIDs := strings.Split(gameID, ",") // 多個game_id結合
		for _, id := range gameIDs {
			// 刪除中獎人員時，需要刪除中獎人員redis資料
			h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + id)

			// 刪除中獎人員時，需要刪除未中獎人員redis資料
			h.redisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + id)

			// 刪除中獎人員時，需要刪除活動下所有場次搖號抽獎的中獎人員redis資料
			h.redisConn.DelCache(config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID)

			// 刪除玩家遊戲紀錄(中獎.未中獎)
			h.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+id, "修改資料")
		}
	}

	response.Ok(ctx)
}

// @@@Summary 刪除拔河遊戲獎品資料(form-data)
// @@@Tags Tugofwar Prize
// @@@version 1.0
// @@@Accept  mpfd
// @@@param id formData string true "ID(以,區隔多筆資料ID)"
// @@@param activity_id formData string true "活動ID"
// @@@param game_id formData string true "遊戲ID"
// @@@param user_id formData string true "用戶ID"
// @@@param token formData string true "CSRF Token"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/tugofwar/prize/form [delete]
func (h *Handler) DELETETugofwarPrize(ctx *gin.Context) {
}

// else if contentType == "application/json" {
// 	var model models.DeleteModel
// 	err := ctx.BindJSON(&model)
// 	if err != nil {
// 		if strings.Contains(ctx.GetHeader("Accept"), "json") {
// 			response.Error(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		} else {
// 			h.executeErrorHTML(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		}
// 	}

// 	id = model.ID
// 	activityID = model.ActivityID
// 	gameID = model.GameID
// 	userID = model.UserID
// 	token = model.Token
// }

// @@@Summary 刪除遊戲抽獎資料(json)
// @@@Tags Lottery
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Lottery Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/json [delete]
func (h *Handler) DELETEJSONLottery(ctx *gin.Context) {
}

// @@@Summary 刪除遊戲抽獎獎品資料(json)
// @@@Tags Lottery Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Lottery Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/prize/json [delete]
func (h *Handler) DELETEJSONLotteryPrize(ctx *gin.Context) {
}

// @@@Summary 刪除搖紅包遊戲資料(json)
// @@@Tags Redpack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Redpack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/json [delete]
func (h *Handler) DELETEJSONRedpack(ctx *gin.Context) {
}

// @@@Summary 刪除套紅包遊戲資料(json)
// @@@Tags Ropepack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Ropepack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/json [delete]
func (h *Handler) DELETEJSONRopepack(ctx *gin.Context) {
}

// @@@Summary 刪除敲敲樂遊戲資料(json)
// @@@Tags Whack_Mole
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "WhackMole Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/json [delete]
func (h *Handler) DELETEJSONWhackMole(ctx *gin.Context) {
}

// @@@Summary 刪除搖紅包獎品資料(json)
// @@@Tags Redpack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Redpack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/prize/json [delete]
func (h *Handler) DELETEJSONRedpackPrize(ctx *gin.Context) {
}

// @@@Summary 刪除套紅包獎品資料(json)
// @@@Tags Ropepack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Ropepack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/prize/json [delete]
func (h *Handler) DELETEJSONRopepackPrize(ctx *gin.Context) {
}

// @@@Summary 刪除敲敲樂獎品資料(json)
// @@@Tags Whack_Mole Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "WhackMole Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/prize/json [delete]
func (h *Handler) DELETEJSONWhackMolePrize(ctx *gin.Context) {
}

// @@@Summary 刪除搖號抽獎遊戲資料(json)
// @@@Tags Draw_Numbers
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Draw_Numbers Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/json [delete]
func (h *Handler) DELETEJSONDrawNumbers(ctx *gin.Context) {
}

// @@@Summary 刪除搖號抽獎獎品資料(json)
// @@@Tags Draw_Numbers Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Draw_Numbers Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/prize/json [delete]
func (h *Handler) DELETEJSONDrawNumbersPrize(ctx *gin.Context) {
}

// @@@Summary 刪除鑑定師遊戲資料(json)
// @@@Tags Monopoly
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Monopoly Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/json [delete]
func (h *Handler) DELETEJSONMonopoly(ctx *gin.Context) {
}

// @@@Summary 刪除鑑定師獎品資料(json)
// @@@Tags Monopoly Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Monopoly Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/prize/json [delete]
func (h *Handler) DELETEJSONMonopolyPrize(ctx *gin.Context) {
}

// @@@Summary 刪除快問快答遊戲資料(json)
// @@@Tags QA
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "QA Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/json [delete]
func (h *Handler) DELETEJSONQA(ctx *gin.Context) {
}

// @@@Summary 刪除遊戲人員資料(json)
// @@@Tags Attend Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Attend Staff Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/attend/json [delete]
func (h *Handler) DELETEJSONAttend(ctx *gin.Context) {
}

// @@@Summary 刪除快問快答獎品資料(json)
// @@@Tags QA Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "QA Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/prize/json [delete]
func (h *Handler) DELETEJSONQAPrize(ctx *gin.Context) {
}

// @@@Summary 刪除中獎人員資料(json)
// @@@Tags Winning Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Winning Staff Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/winning/json [delete]
func (h *Handler) DELETEJSONWinning(ctx *gin.Context) {
}

// @@@Summary 刪除遊戲紀錄資料(json)
// @@@Tags Record Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.DeleteModel true "Record Staff Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/record/json [delete]
func (h *Handler) DELETEJSONRecord(ctx *gin.Context) {
}

// @@@param body body models.DeleteModel true "Rec
