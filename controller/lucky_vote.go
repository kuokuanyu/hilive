package controller

import (
	"hilive/modules/response"

	"github.com/gin-gonic/gin"
)

// VoteParam 投票遊戲中獎判斷回傳參數
type VoteParam struct {
	Game GameModel
}

// @Summary 投票發獎判斷，回傳所有中獎資訊
// @Tags Lucky
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param prize_id query string true "prize ID"
// @Success 200 {array} TugofwarParam
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /game/vote [get]
func (h *Handler) GetVote(ctx *gin.Context) {
	// var (
	// 	host       = ctx.Request.Host
	// 	activityID = ctx.Query("activity_id")
	// 	gameID     = ctx.Query("game_id")
	// 	prizeID    = ctx.Query("prize_id")
	// 	game       = "vote"
	// 	gameModel  = h.getGameInfo(gameID, "vote") // 遊戲資訊
	// 	wg         sync.WaitGroup
	// 	lock       sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
	// 	now, _     = time.ParseInLocation("2006-01-02 15:04:05",
	// 		time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	// 	winOptionID string
	// )

	// if host != config.API_URL {
	// 	response.Error(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 網域請求發生問題",
	// 	})
	// 	return
	// }
	// if activityID == "" || gameID == "" || gameModel.ID == 0 || prizeID == "" {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 無法辨識活動、遊戲、獎品資訊，請輸入有效的ID",
	// 	})
	// 	return
	// }

	// // 取得所有選項分數資料(從redis取得，由高到低)
	// voteOptionModels, err := models.DefaultGameVoteOptionModel().
	// SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 	FindOrderByScore(true, gameID)
	// if err != nil {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	// // 獲勝隊伍
	// if len(voteOptionModels) > 0 {
	// 	winOptionID = voteOptionModels[0].OptionID
	// } else if len(voteOptionModels) == 0 {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 該投票遊戲場次無任何選項，無法發獎",
	// 	})
	// 	return
	// }

	// // 將投票選項結果寫入遊戲紀錄表中
	// // 併發寫入資料
	// wg.Add(len(voteOptionModels)) //計數器
	// for i := 0; i < len(voteOptionModels); i++ {
	// 	go func(i int) {
	// 		defer wg.Done()

	// 		if err = models.DefaultScoreModel().SetDbConn(h.dbConn).
	// 			SetRedisConn(h.redisConn).Add(models.NewScoreModel{
	// 			ActivityID: activityID,
	// 			GameID:     gameID,
	// 			UserID:     "",
	// 			OptionID:   voteOptionModels[i].OptionID,
	// 			Round:      gameModel.GameRound,
	// 			Score:      voteOptionModels[i].OptionScore,
	// 		}); err != nil {
	// 			log.Println("錯誤: 新增投票選項分數資料發生問題，請重新操作")
	// 			return
	// 		}
	// 	}(i)
	// }
	// wg.Wait() //等待計數器歸0

	// // 獎品資訊
	// prize, err := models.DefaultPrizeModel().SetDbConn(h.dbConn).
	// 	FindPrize(false, prizeID)
	// if err != nil {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 取得拔河遊戲獎品資訊發生問題",
	// 	})
	// 	return
	// }
	// prizeAmount := prize.PrizeRemain // 剩餘獎品數量

	// // 發獎，判斷發獎方式
	// if gameModel.Prize == "uniform" {
	// 	// 判斷剩餘獎品數量
	// 	if prizeAmount < 1 {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 獎品數不夠發放，請設置足夠獎品數",
	// 		})
	// 		return
	// 	}

	// 	// 遞減redis獎品數量
	// 	h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

	// 	// 統一發獎給隊長
	// 	// 取得隊長資訊
	// 	lists, err := models.DefaultGameVoteOptionListModel().
	// 		SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		Find(false, gameID, winOptionID, "leader")
	// 	if err != nil || len(lists) == 0 {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 該選項無隊長，無法發放獎品",
	// 		})
	// 		return
	// 	}

	// 	// 新增中獎人員名單
	// 	id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
	// 		Add(models.NewPrizeStaffModel{
	// 			ActivityID: activityID,
	// 			GameID:     gameID,
	// 			UserID:     lists[0].UserID,
	// 			PrizeID:    prizeID,
	// 			Round:      strconv.Itoa(int(gameModel.GameRound)),
	// 			Status:     "no",
	// 			Score:      voteOptionModels[0].OptionScore,
	// 			Leader:     "yes",
	// 		})
	// 	if err != nil {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 發獎給隊長時發生問題(統一發獎)，請重新操作",
	// 		})
	// 		return
	// 	}

	// 	// 用戶遊戲紀錄
	// 	record := models.PrizeStaffModel{
	// 		ID:            id,
	// 		ActivityID:    activityID,
	// 		GameID:        gameID,
	// 		PrizeID:       prizeID,
	// 		UserID:        lists[0].UserID,
	// 		Name:          lists[0].Name,
	// 		Avatar:        lists[0].Avatar,
	// 		Game:          game,
	// 		PrizeName:     prize.PrizeName,
	// 		PrizeType:     prize.PrizeType,
	// 		PrizePicture:  prize.PrizePicture,
	// 		PrizePrice:    prize.PrizePrice,
	// 		PrizeMethod:   prize.PrizeMethod,
	// 		PrizePassword: prize.PrizePassword,
	// 		Round:         gameModel.GameRound,
	// 		WinTime:       now.Format("2006-01-02 15:04:05"),
	// 		Status:        "no",
	// 		Score:         voteOptionModels[0].OptionScore,
	// 		Leader:        "yes",
	// 	}

	// 	// 隊長遊戲紀錄(包含中獎與未中獎)
	// 	allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, lists[0].UserID)
	// 	if err != nil {
	// 		return
	// 	}
	// 	allRecords = append(allRecords, record)

	// 	// 更新redis裡隊長遊戲紀錄
	// 	h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
	// 		lists[0].UserID, utils.JSON(allRecords))
	// 	h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
	// 		config.REDIS_EXPIRE)

	// } else if gameModel.Prize == "all" {
	// 	// 取得選項名單資訊(所有人員)
	// 	lists, err := models.DefaultGameVoteOptionListModel().
	// 		SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		Find(false, gameID, winOptionID, "")
	// 	if err != nil || len(lists) == 0 {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 該選項無隊長，無法發放獎品",
	// 		})
	// 		return
	// 	}

	// 	if int64(len(lists)) > prizeAmount {
	// 		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 			UserID:  "",
	// 			Method:  ctx.Request.Method,
	// 			Path:    ctx.Request.URL.Path,
	// 			Message: "錯誤: 選項名單數量大於剩餘獎品數，無法發放獎品",
	// 		})
	// 		return
	// 	}

	// 	wg.Add(len(lists)) //計數器
	// 	for i := 0; i < len(lists); i++ {
	// 		go func(i int) {
	// 			defer wg.Done()
	// 			var (
	// 				isLeader = "no"
	// 			)
	// 			if lists[i].Leader == "leader" {
	// 				isLeader = "yes"
	// 			}

	// 			// 遞減redis獎品數量
	// 			h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

	// 			// 新增中獎人員名單
	// 			id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
	// 				Add(models.NewPrizeStaffModel{
	// 					ActivityID: activityID,
	// 					GameID:     gameID,
	// 					UserID:     lists[i].UserID,
	// 					PrizeID:    prizeID,
	// 					Round:      strconv.Itoa(int(gameModel.GameRound)),
	// 					Status:     "no",
	// 					Score:      voteOptionModels[0].OptionScore,
	// 					Leader:     isLeader,
	// 				})
	// 			if err != nil {
	// 				log.Println(err)
	// 				return
	// 			}

	// 			// 用戶遊戲紀錄
	// 			record := models.PrizeStaffModel{
	// 				ID:            id,
	// 				ActivityID:    activityID,
	// 				GameID:        gameID,
	// 				PrizeID:       prizeID,
	// 				UserID:        lists[i].UserID,
	// 				Name:          lists[i].Name,
	// 				Avatar:        lists[i].Avatar,
	// 				Game:          game,
	// 				PrizeName:     prize.PrizeName,
	// 				PrizeType:     prize.PrizeType,
	// 				PrizePicture:  prize.PrizePicture,
	// 				PrizePrice:    prize.PrizePrice,
	// 				PrizeMethod:   prize.PrizeMethod,
	// 				PrizePassword: prize.PrizePassword,
	// 				Round:         gameModel.GameRound,
	// 				WinTime:       now.Format("2006-01-02 15:04:05"),
	// 				Status:        "no",
	// 				Score:         voteOptionModels[0].OptionScore,
	// 				Leader:        isLeader,
	// 			}

	// 			// 用戶遊戲紀錄(包含中獎與未中獎)
	// 			allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, lists[i].UserID)
	// 			if err != nil {
	// 				log.Println(err)
	// 				return
	// 			}
	// 			allRecords = append(allRecords, record)

	// 			// 更新redis裡用戶遊戲紀錄
	// 			h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
	// 				lists[i].UserID, utils.JSON(allRecords))
	// 			h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID,
	// 				config.REDIS_EXPIRE)

	// 			lock.Lock() //佔有資源
	// 			prizeAmount--
	// 			lock.Unlock() //釋放資源
	// 		}(i)
	// 	}
	// 	wg.Wait() //等待計數器歸0
	// }

	// // 更新剩餘獎品數量
	// if err = models.DefaultPrizeModel().SetDbConn(h.dbConn).
	// 	UpdateRemain(prizeID, prizeAmount); err != nil {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	// // 更新遊戲輪次
	// if err = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 	UpdateGameStatus(true, gameID, "+1", "", ""); err != nil {
	// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
	// 		UserID:  "",
	// 		Method:  ctx.Request.Method,
	// 		Path:    ctx.Request.URL.Path,
	// 		Message: "錯誤: 無法更新遊戲輪次資料",
	// 	})
	// 	return
	// }

	// // 清除分數資料
	// h.redisConn.DelCache(config.SCORES_REDIS + gameID) // 分數

	// 回傳隊伍資訊
	response.OkWithData(ctx, VoteParam{
		Game: GameModel{},
	})
}
