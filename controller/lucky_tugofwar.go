package controller

import (
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"hilive/modules/utils"

	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TugofwarParam 拔河遊戲中獎判斷回傳參數
type TugofwarParam struct {
	Game GameModel
}

// @Summary 拔河遊戲發獎判斷，回傳所有中獎資訊、獲勝隊伍、雙方隊長資訊、雙方分數資訊、雙方MVP資訊、雙方所有成員分數資訊
// @Tags Lucky
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Success 200 {array} TugofwarParam
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /game/tugofwar [get]
func (h *Handler) GetTugofwar(ctx *gin.Context) {
	var (
		host                            = ctx.Request.Host
		activityID                      = ctx.Query("activity_id")
		gameID                          = ctx.Query("game_id")
		game                            = "tugofwar"
		gameModel                       = h.getGameInfo(gameID, "tugofwar") // 遊戲資訊
		teamModel                       = h.getTeamInfo(gameID)             // 隊伍資訊
		wg                              sync.WaitGroup
		lock                            sync.Mutex                // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
		leftTeamLeader, rightTeamLeader UserModel                 // 左右方隊長資訊
		leftTeamMVP, rightTeamMVP       UserModel                 // 左右方MVP資訊
		leftTeamPlayers                 = make([]UserModel, 0)    // 左方玩家資訊(不包含隊長、mvp)
		rightTeamPlayers                = make([]UserModel, 0)    // 左方玩家資訊(不包含隊長、mvp)
		leftTeamAllPlayers              = make([]UserModel, 0)    // 左方玩家資訊(包含隊長、mvp，分數由高到低)
		rightTeamAllPlayers             = make([]UserModel, 0)    // 左方玩家資訊(包含隊長、mvp，分數由高到低)
		winLeader, winMVP               UserModel                 // 獲勝隊伍的隊長、mvp
		loseLeader, loseMVP             UserModel                 // 落敗隊伍的隊長、mvp
		winPlayers                      = make([]UserModel, 0)    // 獲勝隊伍的所有隊員
		losePlayers                     = make([]UserModel, 0)    // 落敗隊伍的所有隊員
		round                           = gameModel.GameRound - 1 // 輪次(結算時輪次已經+1)
		leftTeamScore, rightTeamScore   int64                     // 雙方隊伍分數資訊
		winPrizeModel, losePrizeModel   models.PrizeModel
		winPrizeAmount, losePrizeAmount int64  // 獎品數
		winPrizeID, losePrizeID         string // 獲勝隊伍獎品ID、落敗隊伍獎品ID
		winTeam, loseTeam               string // 獲勝隊伍、落敗隊伍
		now, _                          = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)
	// startTime := time.Now()
	// log.Println("拔河遊戲結算")

	// 遊戲人員資料處理
	// 取得該場次所有遊戲人員資訊
	allplayers := make([]string, 0)
	leftplayers, _ := h.redisConn.SetGetMembers(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + gameID)
	rightplayers, _ := h.redisConn.SetGetMembers(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + gameID)
	// log.Println("左隊人員: ", len(leftplayers))
	// log.Println("右隊人員: ", len(rightplayers))

	// 將兩隊資料寫入同一個陣列中
	allplayers = append(allplayers, leftplayers...)
	allplayers = append(allplayers, rightplayers...)

	// log.Println("批量將兩隊人員遊戲紀錄寫入資料表: ", len(allplayers))

	// 批量將所有遊戲人員資料寫入資料表(拔河遊戲的遊戲中輪次已遞增)
	err := models.DefaultGameStaffModel().
	SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Adds(len(allplayers), activityID, gameID, game, allplayers, int(round))
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// log.Println("批量將兩隊人員遊戲紀錄寫入完成")

	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" || gameID == "" || gameModel.ID == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的ID",
		})
		return
	}

	// 取得所有人員分數資料(從redis取得，由高到低，包括0分玩家)
	staffs, _, err := h.getWinningRecords(gameModel, game, gameID, round, 0, 0)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}
	// fmt.Println("staffs: ", staffs)

	// 將所有人員分數資料依照隊伍分至陣列中、取得雙方隊長.MVP.分數資訊
	for _, staff := range staffs {
		userModel := UserModel{
			UserID: staff.UserID,
			Name:   staff.Name,
			Avatar: staff.Avatar,
			Score:  staff.Score,
		}

		// 判斷玩家的隊伍
		if utils.InArray(teamModel.LeftTeamPlayers, staff.UserID) {
			leftTeamAllPlayers = append(leftTeamAllPlayers, userModel)
			if staff.UserID == teamModel.LeftTeamLeader {
				// 判斷是否為隊長
				leftTeamLeader = userModel
			} else if staff.UserID != teamModel.LeftTeamLeader &&
				leftTeamMVP.UserID == "" {
				// 判斷MVP(不是隊長、分數最高的)
				leftTeamMVP = userModel

				// 左方隊伍，將玩家分數資訊加入陣列中
				leftTeamPlayers = append(leftTeamPlayers, userModel)
			} else {
				// 左方隊伍，將玩家分數資訊加入陣列中
				leftTeamPlayers = append(leftTeamPlayers, userModel)
			}

			// 遞增隊伍分數
			leftTeamScore += staff.Score
		} else if utils.InArray(teamModel.RightTeamPlayers, staff.UserID) {
			rightTeamAllPlayers = append(rightTeamAllPlayers, userModel)
			if staff.UserID == teamModel.RightTeamLeader {
				// 判斷是否為隊長
				rightTeamLeader = userModel
			} else if staff.UserID != teamModel.RightTeamLeader &&
				rightTeamMVP.UserID == "" {
				// 判斷MVP(不是隊長、分數最高的)
				rightTeamMVP = userModel

				// 右方隊伍，將玩家分數資訊加入陣列中
				rightTeamPlayers = append(rightTeamPlayers, userModel)
			} else {
				// 右方隊伍，將玩家分數資訊加入陣列中
				rightTeamPlayers = append(rightTeamPlayers, userModel)
			}

			// 遞增隊伍分數
			rightTeamScore += staff.Score
		}
	}
	// fmt.Println("測試1，處理左方隊長: ", leftTeamLeader)
	// fmt.Println("測試1，處理左方mvp: ", leftTeamMVP)
	// fmt.Println("測試1，處理左方玩家: ", leftTeamPlayers)
	// fmt.Println("-------------------------------------------------------")
	// fmt.Println("測試2，處理右方隊長: ", rightTeamLeader)
	// fmt.Println("測試2，處理右方mvp: ", rightTeamMVP)
	// fmt.Println("測試2，處理右方玩家: ", rightTeamPlayers)

	// 判斷獲勝隊伍
	if leftTeamScore > rightTeamScore {
		// 左隊獲勝
		winTeam = "left_team"
		winLeader = leftTeamLeader
		winMVP = leftTeamMVP
		winPlayers = leftTeamAllPlayers
		loseTeam = "right_team"
		loseLeader = rightTeamLeader
		loseMVP = rightTeamMVP
		losePlayers = rightTeamAllPlayers
	} else if leftTeamScore < rightTeamScore {
		// 右隊獲勝
		winTeam = "right_team"
		winLeader = rightTeamLeader
		winMVP = rightTeamMVP
		winPlayers = rightTeamAllPlayers
		loseTeam = "left_team"
		loseLeader = leftTeamLeader
		loseMVP = leftTeamMVP
		losePlayers = leftTeamAllPlayers
	} else if leftTeamScore == rightTeamScore {
		// 兩隊分數平手，比隊長分數
		if leftTeamLeader.Score > rightTeamLeader.Score {
			winTeam = "left_team"
			winLeader = leftTeamLeader
			winMVP = leftTeamMVP
			winPlayers = leftTeamAllPlayers
			loseTeam = "right_team"
			loseLeader = rightTeamLeader
			loseMVP = rightTeamMVP
			losePlayers = rightTeamAllPlayers
		} else if leftTeamLeader.Score < rightTeamLeader.Score {
			winTeam = "right_team"
			winLeader = rightTeamLeader
			winMVP = rightTeamMVP
			winPlayers = rightTeamAllPlayers
			loseTeam = "left_team"
			loseLeader = leftTeamLeader
			loseMVP = leftTeamMVP
			losePlayers = leftTeamAllPlayers
		} else if leftTeamLeader.Score == rightTeamLeader.Score {
			// 兩隊隊長分數平手，比MVP分數
			if leftTeamMVP.Score > rightTeamMVP.Score {
				winTeam = "left_team"
				winLeader = leftTeamLeader
				winMVP = leftTeamMVP
				winPlayers = leftTeamAllPlayers
				loseTeam = "right_team"
				loseLeader = rightTeamLeader
				loseMVP = rightTeamMVP
				losePlayers = rightTeamAllPlayers
			} else if leftTeamMVP.Score < rightTeamMVP.Score {
				winTeam = "right_team"
				winLeader = rightTeamLeader
				winMVP = rightTeamMVP
				winPlayers = rightTeamAllPlayers
				loseTeam = "left_team"
				loseLeader = leftTeamLeader
				loseMVP = leftTeamMVP
				losePlayers = leftTeamAllPlayers
			}
		}
	}
	// fmt.Println("測試3，取得獲勝隊伍資訊: ", winTeam, winLeader, winMVP, winPlayers)

	// 將玩家分數資料寫入遊戲紀錄表中(勝方敗方都要寫入)
	for m := 0; m <= 1; m++ {
		var (
			team        string
			leader, mvp UserModel
			players     = make([]UserModel, 0)
		)
		if m == 0 {
			// 勝方
			team = winTeam
			leader = winLeader
			mvp = winMVP
			players = winPlayers
		} else if m == 1 {
			// 敗方
			team = loseTeam
			leader = loseLeader
			mvp = loseMVP
			players = losePlayers
		}

		// log.Println("拔河遊戲批量寫入分數人員資料(勝敗方): ", len(players))

		if len(players) > 0 {
			var (
				userIDs = make([]string, len(players)) // 用戶ID
				scores1 = make([]string, len(players)) // 分數
				scores2 = make([]string, len(players)) // 分數2(0)
				teams   = make([]string, len(players)) // 隊伍資訊
				leaders = make([]string, len(players)) // 隊長資訊
				mvps    = make([]string, len(players)) // mvp資訊
			)

			// 處理批量寫入資料的相關參數
			for i := 0; i < len(players); i++ {
				var (
					isLeader = "no"
					isMvp    = "no"
				)

				userIDs[i] = players[i].UserID
				scores1[i] = strconv.Itoa(int(players[i].Score))
				scores2[i] = strconv.Itoa(0)

				if players[i].UserID == leader.UserID {
					isLeader = "yes"
				} else if players[i].UserID == mvp.UserID {
					isMvp = "yes"
				}

				teams[i] = team
				leaders[i] = isLeader
				mvps[i] = isMvp
			}

			// 將所有人員分數資料批量寫入資料表中
			if err = models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
				Adds(len(players), activityID, gameID,game,  userIDs,
				int(round), scores1, scores2, teams, leaders, mvps); err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 將所有計分資料寫入資料表發生問題，請重新操作",
				})
				return
			}
		}

		// log.Println("拔河遊戲批量寫入分數人員資料(勝敗方)完成")

	}
	// fmt.Println("測試4，將左右方玩家遊戲紀錄寫入資料表中")

	// 獎品資訊
	prizes, err := models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindPrizes(false, gameID)
	if err != nil || len(prizes) != 2 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 取得拔河遊戲獎品資訊發生問題",
		})
		return
	}
	for _, prize := range prizes {
		if prize.TeamType == "win" {
			winPrizeID = prize.PrizeID
			winPrizeAmount = prize.PrizeRemain
			winPrizeModel = prize
		} else if prize.TeamType == "lose" {
			losePrizeID = prize.PrizeID
			losePrizeAmount = prize.PrizeRemain
			losePrizeModel = prize
		}
	}

	// 發獎，判斷發獎方式
	if gameModel.Prize == "uniform" {
		// log.Println("測試5，統一發獎給隊長")
		// 統一發獎給隊長(勝方敗方都要發)
		for i := 0; i <= 1; i++ {
			var (
				prizeID, team string
				leader        UserModel
				prizeModel    models.PrizeModel
			)
			if i == 0 {
				// 勝方
				prizeID = winPrizeID
				team = winTeam
				leader = winLeader
				prizeModel = winPrizeModel
			} else if i == 1 {
				// 敗方
				prizeID = losePrizeID
				team = loseTeam
				leader = loseLeader
				prizeModel = losePrizeModel
			}

			// 遞減redis獎品數量
			// h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

			// 新增中獎人員名單
			id, err := models.DefaultPrizeStaffModel().
			SetConn(h.dbConn, h.redisConn,h.mongoConn).
				Add(models.NewPrizeStaffModel{
					ActivityID: activityID,
					GameID:     gameID,
					Game: game,
					UserID:     leader.UserID,
					PrizeID:    prizeID,
					Round:      strconv.Itoa(int(round)),
					Status:     "no",
					Score:      leader.Score,
					Score2:     0,
					Rank:       1,
					Team:       team,
					Leader:     "yes",
					Mvp:        "no",
				})
			if err != nil {
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 發獎給隊長時發生問題(統一發獎)，請重新操作",
				})
				return
			}

			// 用戶遊戲紀錄
			record := models.PrizeStaffModel{
				ID:            id,
				ActivityID:    activityID,
				GameID:        gameID,
				PrizeID:       prizeID,
				UserID:        leader.UserID,
				Name:          leader.Name,
				Avatar:        leader.Avatar,
				Game:          game,
				PrizeName:     prizeModel.PrizeName,
				PrizeType:     prizeModel.PrizeType,
				PrizePicture:  prizeModel.PrizePicture,
				PrizePrice:    prizeModel.PrizePrice,
				PrizeMethod:   prizeModel.PrizeMethod,
				PrizePassword: prizeModel.PrizePassword,
				Round:         round,
				WinTime:       now.Format("2006-01-02 15:04:05"),
				Status:        "no",
				Score:         leader.Score,
				Score2:        0,
				Rank:          1,
				Team:          team,
				Leader:        "yes",
				Mvp:           "no",
			}

			// 隊長遊戲紀錄(包含中獎與未中獎)
			allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, leader.UserID)
			if err != nil {
				// log.Println(err)
				return
			}
			allRecords = append(allRecords, record)

			// 更新redis裡隊長遊戲紀錄
			h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
				leader.UserID, utils.JSON(allRecords))

			// 設置過期時間
			// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)
		}

		winPrizeAmount--
		losePrizeAmount--
	} else if gameModel.Prize == "all" {
		// log.Println("測試6，同時發獎(勝方敗方都要發)")

		// 同時發獎(勝方敗方都要發)
		for m := 0; m <= 1; m++ {
			var (
				prizeID, team string
				leader, mvp   UserModel
				prizeModel    models.PrizeModel
				players       = make([]UserModel, 0)
			)
			if m == 0 {
				// 勝方
				prizeID = winPrizeID
				team = winTeam
				leader = winLeader
				mvp = winMVP
				prizeModel = winPrizeModel
				players = winPlayers
			} else if m == 1 {
				// 敗方
				prizeID = losePrizeID
				team = loseTeam
				leader = loseLeader
				mvp = loseMVP
				prizeModel = losePrizeModel
				players = losePlayers
			}

			// 批量寫入所有中獎人員資料
			// log.Println("批量寫入所有中獎人員資料: ", len(players))

			var (
				userIDs  = make([]string, len(players)) // 用戶ID
				peizeIDs = make([]string, len(players)) // 獎品ID
				scores1  = make([]string, len(players)) // 分數
				scores2  = make([]string, len(players)) // 分數2(0)
				ranks    = make([]string, len(players)) // 排名資訊
				teams    = make([]string, len(players)) // 隊伍資訊
				leaders  = make([]string, len(players)) // 隊長資訊
				mvps     = make([]string, len(players)) // mvp資訊
			)

			// 處理批量寫入資料的相關參數
			for i := 0; i < len(players); i++ {
				var (
					isLeader = "no"
					isMvp    = "no"
				)

				if players[i].UserID == leader.UserID {
					isLeader = "yes"
				} else if players[i].UserID == mvp.UserID {
					isMvp = "yes"
				}

				userIDs[i] = players[i].UserID
				peizeIDs[i] = prizeID
				scores1[i] = strconv.Itoa(int(players[i].Score))
				scores2[i] = strconv.Itoa(0)
				ranks[i] = strconv.Itoa(i + 1)
				teams[i] = team
				leaders[i] = isLeader
				mvps[i] = isMvp

				// 獎品數量處理
				lock.Lock() //佔有資源
				if m == 0 {
					winPrizeAmount--
				} else if m == 1 {
					losePrizeAmount--
				}
				lock.Unlock() //釋放資源
			}

			// 將所有中獎人員資料批量寫入資料表中
			if err = models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
				Adds(len(players), activityID, gameID, game, userIDs, peizeIDs,
				int(round), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
				// log.Println("錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作")
				response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  "",
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 將所有中獎人員資料寫入資料表發生問題，請重新操作",
				})
				return
			}

			// log.Println("批量寫入所有中獎人員完成")

			// 多線呈更新所有玩家遊戲紀錄(redis)
			wg.Add(len(players)) //計數器
			for i := 0; i < len(players); i++ {
				go func(i int) {
					defer wg.Done()
					var (
						isLeader = "no"
						isMvp    = "no"
					)
					if players[i].UserID == leader.UserID {
						isLeader = "yes"
					} else if players[i].UserID == mvp.UserID {
						isMvp = "yes"
					}

					// 用戶遊戲紀錄
					record := models.PrizeStaffModel{
						ID:            0,
						ActivityID:    activityID,
						GameID:        gameID,
						PrizeID:       prizeID,
						UserID:        players[i].UserID,
						Name:          players[i].Name,
						Avatar:        players[i].Avatar,
						Game:          game,
						PrizeName:     prizeModel.PrizeName,
						PrizeType:     prizeModel.PrizeType,
						PrizePicture:  prizeModel.PrizePicture,
						PrizePrice:    prizeModel.PrizePrice,
						PrizeMethod:   prizeModel.PrizeMethod,
						PrizePassword: prizeModel.PrizePassword,
						Round:         round,
						WinTime:       now.Format("2006-01-02 15:04:05"),
						Status:        "no",
						Score:         players[i].Score,
						Score2:        0,
						Rank:          int64(i + 1),
						Team:          team,
						Leader:        isLeader,
						Mvp:           isMvp,
					}

					// 用戶遊戲紀錄(包含中獎與未中獎)
					allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, game, players[i].UserID)
					if err != nil {
						// log.Println(err)
						return
					}
					allRecords = append(allRecords, record)

					// 更新redis裡用戶遊戲紀錄
					h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
						players[i].UserID, utils.JSON(allRecords))
					// 設置過期時間
					// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)

				}(i)
			}
			wg.Wait() //等待計數器歸0

			// log.Println("多線呈結束")
		}
	}

	// 更新剩餘獎品數量(勝方，資料表)
	if err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		UpdateRemain(gameID, winPrizeID, winPrizeAmount); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}
	// 更新剩餘獎品數量(敗方，資料表)
	if err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		UpdateRemain(gameID, losePrizeID, losePrizeAmount); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}
	// fmt.Println("獎品數量: ", winPrizeAmount, losePrizeAmount)

	// 清除分數資料
	h.redisConn.DelCache(config.SCORES_REDIS + gameID)   // 分數
	h.redisConn.DelCache(config.SCORES_2_REDIS + gameID) // 第二分數
	// 清除中獎人員資料
	// h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID)
	// 清除遊戲隊伍資訊
	// h.redisConn.DelCache(config.GAME_TEAM_REDIS + gameID)

	// 將獲勝隊伍寫入redis中
	h.redisConn.HashSetCache(config.GAME_TEAM_REDIS+gameID,
		"win_team", winTeam)

	// 設置過期時間
	// h.redisConn.SetExpire(config.GAME_TEAM_REDIS+gameID, config.REDIS_EXPIRE)

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
	h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	h.redisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

	// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道(拔河隊伍人數)
	h.redisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")

	// endTime := time.Now()
	// log.Println("統計中獎人員花費時間: ", endTime.Sub(startTime))
	// fmt.Println("拔河遊戲發獎結束")
	// 回傳隊伍資訊
	response.OkWithData(ctx, TugofwarParam{
		Game: GameModel{
			LeftTeamLeader:   leftTeamLeader,
			RightTeamLeader:  rightTeamLeader,
			LeftTeamMVP:      leftTeamMVP,
			RightTeamMVP:     rightTeamMVP,
			LeftTeamPlayers:  leftTeamPlayers,
			RightTeamPlayers: rightTeamPlayers,
			LeftTeamScore:    leftTeamScore,
			RightTeamScore:   rightTeamScore,
			// PrizeAmount:      prizeAmount,
			WinPrizeAmount:  winPrizeAmount,
			LosePrizeAmount: losePrizeAmount,
			WinTeam:         winTeam,
		},
	})
}

// 遞減redis獎品數量
// h.redisConn.HashDecrCache(config.GAME_PRIZES_AMOUNT_REDIS+gameID, prizeID)

// 新增中獎人員名單
// id, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	Add(models.NewPrizeStaffModel{
// 		ActivityID: activityID,
// 		GameID:     gameID,
// 		UserID:     players[i].UserID,
// 		PrizeID:    prizeID,
// 		Round:      strconv.Itoa(int(round)),
// 		Status:     "no",
// 		Score:      players[i].Score,
// 		Score2:     0,
// 		Rank:       int64(i + 1),
// 		Team:       team,
// 		Leader:     isLeader,
// 		Mvp:        isMvp,
// 	})
// if err != nil {
// 	log.Println(err)
// 	return
// }

// 遊戲獎品數(redis處理)
// allPrizeAmount, err := h.getPrizesAmount(gameID)
// if err != nil {
// 	response.BadRequest(ctx, err.Error())
// 	return
// }

// 更新雙方隊伍人數資料(資料庫.redis)
// if err := models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 	UpdateBothTeamPeople(true, gameID, 0, 0); err != nil {
// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
// 		UserID:  "",
// 		Method:  ctx.Request.Method,
// 		Path:    ctx.Request.URL.Path,
// 		Message: err.Error(),
// 	})
// }

// var prizeAmount int64
// if winPrizeAmount == 0 || losePrizeAmount == 0 {
// 	prizeAmount = 0
// } else {
// 	prizeAmount =winPrizeAmount + losePrizeAmount
// }

// 併發寫入資料
// wg.Add(len(players)) //計數器
// for i := 0; i < len(players); i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		var (
// 			isLeader = "no"
// 			isMvp    = "no"
// 		)
// 		if players[i].UserID == leader.UserID {
// 			isLeader = "yes"
// 		} else if players[i].UserID == mvp.UserID {
// 			isMvp = "yes"
// 		}

// 		if err = models.DefaultScoreModel().SetDbConn(h.dbConn).
// 			SetRedisConn(h.redisConn).Add(models.NewScoreModel{
// 			ActivityID: activityID,
// 			GameID:     gameID,
// 			UserID:     players[i].UserID,
// 			Round:      round,
// 			Score:      players[i].Score,
// 			Score2:     0,
// 			Team:       team,
// 			Leader:     isLeader,
// 			Mvp:        isMvp,
// 		}); err != nil {
// 			log.Println("錯誤: 新增玩家分數資料發生問題，請重新操作")
// 			return
// 		}
// 	}(i)
// }
// wg.Wait() //等待計數器歸0

// 處理左右方隊伍玩家資訊(將隊長、mvp、所有玩家加入陣列)
// 左方隊伍
// if leftTeamLeader.UserID != "" {
// 	leftTeamAllPlayers = append(leftTeamAllPlayers, leftTeamLeader)
// }
// if leftTeamMVP.UserID != "" {
// 	leftTeamAllPlayers = append(leftTeamAllPlayers, leftTeamMVP)
// }
// leftTeamAllPlayers = append(leftTeamAllPlayers, leftTeamPlayers...)
// // 右方隊伍
// if rightTeamLeader.UserID != "" {
// 	rightTeamAllPlayers = append(rightTeamAllPlayers, rightTeamLeader)
// }
// if rightTeamMVP.UserID != "" {
// 	rightTeamAllPlayers = append(rightTeamAllPlayers, rightTeamMVP)
// }
// rightTeamAllPlayers = append(rightTeamAllPlayers, rig
