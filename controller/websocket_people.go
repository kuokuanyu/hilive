package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 即時回傳遊戲人數、輪次、隊長、分數資訊(主持人端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, monopoly, draw_numbers, QA, tugofwar, bingo)
// @Success 200 {array} GameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/people [get]
func (h *Handler) PeopleWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		gameModel         = h.getGameInfo(gameID, game)
		isOPen            = true
		channels          = []string{config.CHANNEL_GAME_REDIS + gameID} // 遊戲資訊

	)
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" ||
		game == "" || gameModel.ID == 0 {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 開啟時傳送訊息
		var (
			gameModel                               = h.getGameInfo(gameID, game) // 即時人數、輪次
			leftTeamGameAttend, rightTeamGameAttend int64                         // 左右方隊伍人數
			leftTeamLeader, rightTeamLeader         UserModel                     // 左右方隊長資訊
			people, bingoAttend                     int64                         // 參加人數
			leftTeamScore, rightTeamScore           int64                         // 左右方隊伍分數
		)

		if game == "draw_numbers" {
			// 搖號抽獎可抽獎人數資訊
			people = h.drawNumbersPeople(activityID, gameID, game, gameModel)
		} else if game == "tugofwar" {
			// 拔河遊戲雙方隊伍人數、雙方隊長資訊、雙方隊伍分數...等資訊
			leftTeamLeader, rightTeamLeader, leftTeamGameAttend,
				rightTeamGameAttend, people, leftTeamScore, rightTeamScore = h.tugofwarPeople(gameID, gameModel)

		} else if game == "QA" {
			// 快問快答遊戲取qa_people參數
			if gameModel.QAPeople >= 0 {
				people = gameModel.QAPeople
			}
		} else {
			// 其他遊戲都取得game_attend參數
			people = gameModel.GameAttend

			// 查詢賓果遊戲完成選號人數(redis)
			bingoAttend, _ = models.DefaultGameModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindUserAmount(true, gameID)
		}

		b, _ := json.Marshal(GameParam{
			Game: GameModel{
				GameAttend:          people,
				BingoAttend:         bingoAttend,
				GameRound:           gameModel.GameRound,
				QARound:             gameModel.QARound,
				BingoRound:          gameModel.BingoRound,
				LeftTeamGameAttend:  leftTeamGameAttend,
				RightTeamGameAttend: rightTeamGameAttend,
				LeftTeamLeader:      leftTeamLeader,
				RightTeamLeader:     rightTeamLeader,
				LeftTeamScore:       leftTeamScore,
				RightTeamScore:      rightTeamScore,
			},
		})
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	if game == "draw_numbers" {
		channels = append(channels,
			config.CHANNEL_SIGN_STAFFS_2_REDIS+activityID, // 簽到人員(搖號抽獎用)
			config.CHANNEL_WINNING_STAFFS_REDIS+gameID,    // 中獎人員(搖號抽獎用)
			config.CHANNEL_BLACK_STAFFS_GAME_REDIS+gameID, // 黑名單人員(搖號抽獎用)
		)
	} else if game == "tugofwar" {
		channels = append(channels,
			config.CHANNEL_GAME_TEAM_REDIS+gameID, // 隊伍資訊(拔河遊戲用)
			config.CHANNEL_SCORES_REDIS+gameID,    // 分數資訊
		)
	} else if game == "bingo" {
		channels = append(channels,
			config.CHANNEL_GAME_BINGO_USER_NUMBER+gameID, // 紀錄玩家的號碼排序(賓果遊戲用)
		)
	}

	// 啟用redis訂閱
	go h.redisConn.Subscribes(context, channels, func(channel, message string) {
		// log.Println("資料改變? ", message)
		// 資料改變時，回傳最新資料至前端
		var (
			gameModel                               = h.getGameInfo(gameID, game) // 即時人數、輪次
			leftTeamGameAttend, rightTeamGameAttend int64                         // 左右方隊伍人數
			leftTeamLeader, rightTeamLeader         UserModel                     // 左右方隊長資訊
			people, bingoAttend                     int64                         // 參加人數
			leftTeamScore, rightTeamScore           int64                         // 左右方隊伍分數
		)

		if game == "draw_numbers" {
			// 搖號抽獎可抽獎人數資訊
			people = h.drawNumbersPeople(activityID, gameID, game, gameModel)
		} else if game == "tugofwar" {
			// 拔河遊戲雙方隊伍人數、雙方隊長資訊、雙方隊伍分數...等資訊
			leftTeamLeader, rightTeamLeader, leftTeamGameAttend,
				rightTeamGameAttend, people, leftTeamScore, rightTeamScore = h.tugofwarPeople(gameID, gameModel)

		} else if game == "QA" {
			// 快問快答遊戲取qa_people參數
			if gameModel.QAPeople >= 0 {
				people = gameModel.QAPeople
			}
		} else {
			// 其他遊戲都取得game_attend參數
			people = gameModel.GameAttend

			if game == "bingo" {
				// 查詢賓果遊戲完成選號人數(redis)
				bingoAttend, _ = models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					FindUserAmount(true, gameID)
			}
		}

		// log.Println("傳送人數: ", people)
		b, _ := json.Marshal(GameParam{
			Game: GameModel{
				GameAttend:          people,
				BingoAttend:         bingoAttend,
				GameRound:           gameModel.GameRound,
				QARound:             gameModel.QARound,
				BingoRound:          gameModel.BingoRound,
				LeftTeamGameAttend:  leftTeamGameAttend,
				RightTeamGameAttend: rightTeamGameAttend,
				LeftTeamLeader:      leftTeamLeader,
				RightTeamLeader:     rightTeamLeader,
				LeftTeamScore:       leftTeamScore,
				RightTeamScore:      rightTeamScore,
			},
		})
		conn.WriteMessage(b)
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉即時人數輪次資訊(主持端)ws")

			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_SIGN_STAFFS_2_REDIS + activityID)
			h.redisConn.Unsubscribe(config.CHANNEL_WINNING_STAFFS_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_BLACK_STAFFS_GAME_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_TEAM_REDIS + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_BINGO_USER_NUMBER + gameID)
			h.redisConn.Unsubscribe(config.CHANNEL_SCORES_REDIS + gameID)

			conn.Close()
			return
		}
	}
}

// tugofwarPeople 拔河遊戲雙方隊伍人數、雙方隊長資訊、雙方隊伍分數...等資訊
func (h *Handler) tugofwarPeople(gameID string, gameModel models.GameModel) (UserModel, UserModel,
	int64, int64, int64, int64, int64) {
	var (
		teamModel                               = h.getTeamInfo(gameID)      // 隊伍資訊
		leftTeamLeaderID                        = teamModel.LeftTeamLeader   // 左方隊長ID
		rightTeamLeaderID                       = teamModel.RightTeamLeader  // 右方隊長ID
		leftTeamPlayers                         = teamModel.LeftTeamPlayers  // 左方隊伍玩家
		rightTeamPlayers                        = teamModel.RightTeamPlayers // 右方隊伍玩家
		leftTeamLeader, rightTeamLeader         UserModel
		leftTeamGameAttend, rightTeamGameAttend int64 // 左右方隊伍人數
		people                                  int64 // 雙方隊伍總人數
		leftTeamScore, rightTeamScore           int64 // 雙方隊伍分數
	)

	// 左右方隊伍人數
	leftTeamGameAttend = gameModel.LeftTeamGameAttend
	rightTeamGameAttend = gameModel.RightTeamGameAttend
	people = leftTeamGameAttend + rightTeamGameAttend

	// 取得左右方隊長資訊
	leftTeamLeader, rightTeamLeader = h.getTeamLeader(true, leftTeamLeaderID, rightTeamLeaderID)

	// 左右隊伍分數
	leftTeamScore, rightTeamScore = h.getTeamScore(true, gameID, leftTeamPlayers, rightTeamPlayers)

	return leftTeamLeader, rightTeamLeader, leftTeamGameAttend, rightTeamGameAttend, people, leftTeamScore, rightTeamScore
}

// drawNumbersPeople 搖號抽獎可抽獎人數資訊
func (h *Handler) drawNumbersPeople(activityID string, gameID string,
	game string, gameModel models.GameModel) int64 {
	var (
		// staffs   = make([]string, 0) // 可抽獎人員
		noStaffs = make([]string, 0) // 不可抽獎人員
		people   int64
	)

	// 簽到人員資料(不需要用戶詳細資訊)
	signStaffs, _ := h.getSignStaffs(true, false, config.SIGN_STAFFS_2_REDIS, activityID, 0, 0)

	// log.Println("執行1")
	// 取得黑名單人員資料
	blackStaffs, _ := models.DefaultBlackStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindAll(true, activityID, gameID, game)
	// 加黑名單人員加入不可抽獎人員中
	for _, staffModel := range blackStaffs {
		noStaffs = append(noStaffs, staffModel.UserID)
	}
	// log.Println("執行1結束")

	// 判斷是否為黑名單
	// for _, staff := range signStaffs {
	// if !h.IsBlackStaff(activityID, gameID, game, staff.UserID) {
	// 不是黑名單
	// staffs = append(staffs, staff.UserID)
	// }
	// }
	// log.Println("執行2")
	// fmt.Println("判斷黑名單後: ", staffs)

	// 搖號抽獎遊戲需要判斷是否重複中獎，如果不允許重複中獎則需要處理可抽獎人數資訊
	if gameModel.AllGameAllow == "open" && gameModel.DrawNumbersGameAllow == "open" &&
		gameModel.Allow == "open" {
		// 允許重複中獎，取得可抽獎人數

		// people = int64(len(staffs))
		// fmt.Println("可抽獎人數: ", people)
	} else if gameModel.AllGameAllow == "close" || gameModel.DrawNumbersGameAllow == "close" ||
		gameModel.Allow == "close" {

		// log.Println("執行2")
		// 不允許重複中獎，取得已中獎人員資料
		prizeStaffs, _ := h.getDrawNumbersWinningStaffs(true, activityID, gameID, game,
			gameModel.AllGameAllow == "open", gameModel.DrawNumbersGameAllow == "open",
			gameModel.Allow == "open") // 中獎人員資料

		// 將已中獎人員資料加入不可抽獎人員中
		noStaffs = utils.AddUniqueStrings(noStaffs, prizeStaffs)
		// log.Println("執行2結束")

		// 判斷staffs參數裡是否有中獎人員資料，如果有則從陣列中取出
		// for _, staff := range staffs {
		// 	if !h.redisConn.SetIsMember(config.WINNING_STAFFS_REDIS+gameID, staff) {
		// 	newStaffs = append(newStaffs, staff)
		// 	}
		// }
	}

	// log.Println("drawNumbersPeople 執行3")

	for _, staff := range signStaffs {
		// 判斷用戶ID是否為不可抽獎人員
		if utils.InArray(noStaffs, staff.UserID) {
			// 不可抽獎人員
		} else {
			// 可抽獎人員，人數遞增
			people++
		}
	}

	// log.Println("drawNumbersPeople 執行3結束")

	// 扣除已中獎人員
	// people = int64(len(newStaffs))
	// if people < 0 {
	// 	people = 0
	// }

	// log.Println("人數: ", people)

	return people
}

// getTeamScore 取得左右方分數資訊
func (h *Handler) getTeamScore(isRedis bool, gameID string,
	leftTeamPlayers, rightTeamPlayers []string) (int64, int64) {
	var (
		leftTeamScore  int64
		rightTeamScore int64
	)

	// 處理左方隊伍分數資訊
	for _, userID := range leftTeamPlayers {
		// 玩家分數資訊
		scoreModel := models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindUser(true, gameID, userID)

		leftTeamScore += scoreModel.Score
	}

	// 處理右方隊伍分數資訊
	for _, userID := range rightTeamPlayers {
		// 玩家分數資訊
		scoreModel := models.DefaultScoreModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindUser(true, gameID, userID)

		rightTeamScore += scoreModel.Score
	}

	return leftTeamScore, rightTeamScore
}

// getTeamLeader 取得左右方隊長資訊
func (h *Handler) getTeamLeader(isRedis bool,
	leftTeamLeaderID, rightTeamLeaderID string) (UserModel, UserModel) {
	var (
		leftTeamLeader  UserModel
		rightTeamLeader UserModel
	)
	// 左方隊長資訊
	leftUser, _ := models.DefaultLineModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, "", "user_id", leftTeamLeaderID)
	leftTeamLeader = UserModel{
		UserID: leftUser.UserID,
		Name:   leftUser.Name,
		Avatar: leftUser.Avatar,
	}

	// 右方隊長資訊
	rightUser, _ := models.DefaultLineModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		Find(true, "", "user_id", rightTeamLeaderID)
	rightTeamLeader = UserModel{
		UserID: rightUser.UserID,
		Name:   rightUser.Name,
		Avatar: rightUser.Avatar,
	}
	return leftTeamLeader, rightTeamLeader
}

// getDrawNumbersAllWinningStaffs 查詢該活動下所有搖號抽獎場次中獎人員資料(從redis取得，如果沒有才執行資料表查詢)
// 回傳中獎人員陣列名單
func (h *Handler) getDrawNumbersAllWinningStaffs(isRedis bool, activityID string) ([]string, error) {

	return models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindDrawNumbersAllWinningStaffs(isRedis, activityID)
}

// getDrawNumbersWinningStaffs 查詢中獎紀錄(從redis取得，如果沒有才執行資料表查詢)
// 先判斷活動是否允許重複中獎、相同類型遊戲是否允許重複中獎、同場次遊戲是否允許重複中獎
// 回傳中獎人員陣列名單
func (h *Handler) getDrawNumbersWinningStaffs(isRedis bool, activityID, gameID, game string,
	activityAllow, gameAllow, allow bool) ([]string, error) {
	// var (
	// 	staffs = make([]string, 0)
	// 	err    error
	// )

	// staffs, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
	// 					SetRedisConn(h.redisConn).FindWinningStaffs(true, gameID) // 中獎人員資料
	// if err != nil {
	// 	return staffs, err
	// }
	return models.DefaultPrizeStaffModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindDrawNumbersWinningStaffs(isRedis, activityID, gameID, game,
			activityAllow, gameAllow, allow)
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			gameModel                               = h.getGameInfo(gameID) // 即時人數、輪次
// 			leftTeamGameAttend, rightTeamGameAttend int64                   // 左右方隊伍人數
// 			leftTeamLeader, rightTeamLeader         UserModel               // 左右方隊長資訊
// 			people, bingoAttend                     int64                   // 參加人數
// 			leftTeamScore, rightTeamScore           int64                   // 左右方隊伍分數
// 		)

// 		if game == "draw_numbers" {
// 			// 搖號抽獎可抽獎人數資訊
// 			people = h.drawNumbersPeople(activityID, gameID, game, gameModel)
// 		} else if game == "tugofwar" {
// 			// 拔河遊戲雙方隊伍人數、雙方隊長資訊、雙方隊伍分數...等資訊
// 			leftTeamLeader, rightTeamLeader, leftTeamGameAttend,
// 				rightTeamGameAttend, people, leftTeamScore, rightTeamScore = h.tugofwarPeople(gameID, gameModel)

// 		} else if game == "QA" {
// 			// 快問快答遊戲取qa_people參數
// 			if gameModel.QAPeople >= 0 {
// 				people = gameModel.QAPeople
// 			}
// 		} else {
// 			// 其他遊戲都取得game_attend參數
// 			people = gameModel.GameAttend

// 			// 查詢賓果遊戲完成選號人數(redis)
// 			bingoAttend, _ = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
// 				FindUserAmount(true, gameID)
// 		}

// 		// #####加入測試資料#####start
// 		// fmt.Sprintln(leftTeamGameAttend, rightTeamGameAttend, people, bingoAttend)

// 		b, _ := json.Marshal(GameParam{
// 			Game: GameModel{
// 				GameAttend:          people,
// 				BingoAttend:         bingoAttend,
// 				GameRound:           gameModel.GameRound,
// 				QARound:             gameModel.QARound,
// 				BingoRound:          gameModel.BingoRound,
// 				LeftTeamGameAttend:  leftTeamGameAttend,
// 				RightTeamGameAttend: rightTeamGameAttend,
// 				LeftTeamLeader:      leftTeamLeader,
// 				RightTeamLeader:     rightTeamLeader,
// 				LeftTeamScore:       leftTeamScore,
// 				RightTeamScore:      rightTeamScore,
// 			},
// 		})
// 		conn.WriteMessage(b)
// 		// #####加入測試資料#####end

// 		if conn.isClose { // ws關閉
// 			// if game == "draw_numbers" {
// 			// 	// h.redisConn.DelCache(config.GAME_REDIS + gameID)           // 遊戲資訊
// 			// 	h.redisConn.DelCache(config.ACTIVITY_REDIS + activityID)   // 活動資訊
// 			// 	h.redisConn.DelCache(config.WINNING_STAFFS_REDIS + gameID) // 中獎人員資訊
// 			// }
// 			return
// 		}
// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// // 活動資訊
// activityModel, _ := models.DefaultActivityModel().
// 	SetDbConn(h.dbConn).SetRedisConn(h.redisConn).Find(true, activityID)

// attend            = gameModel.GameAttend
// round             = gameModel.GameRound

// 刪除redis所有遊戲相關資訊
// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊

// 初次傳送人數、遊戲輪次訊息
// b, _ := json.Marshal(GameParam{
// 	Game: GameModel{
// 		GameAttend: attend,
// 		GameRound:  round,
// 		// Amount: int64(prizeLength),
// 	},
// })
// if err := conn.WriteMessage(b); err != nil {
// 	return
// }

// ; gameModel.GameAttend != attend ||
// gameModel.GameRound != round { // 如果人數或輪次更新才會回傳前端
// attend = gameModel.GameAttend
// round = gameModel.GameRound

// 刪除redis所有遊戲相關資訊
// h.redisConn.DelCache(config.GAME_REDIS + gameID) // 遊戲資訊

// 獎品數
// prizeLength, err := h.getPrizeAmount(ctx.Request.URL.Path, gameID)
// if err != nil {
// 	b, _ := json.Marshal(WebsocketMessage{Error: err.Error()})
// conn.WriteMessage(b)
// 	return
// }

// 獎品數
// prizeLength, err = h.getPrizeAmount(ctx.Request.URL.Path, gameID)
// if err != nil {
// 	b, _ := json.Marshal(WebsocketMessage{Error: err.Error()})
// 	con
