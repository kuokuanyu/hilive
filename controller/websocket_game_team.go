package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"

	"github.com/gin-gonic/gin"
)

// UserGameParam 用戶、場次參數資訊
type UserGameParam struct {
	User  UserModel // 用戶資訊
	Game  GameModel // 場次資訊
	Error string    `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 即時回傳雙方隊伍所有玩家、隊長資訊(後端平台)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "activity ID"
// @param body body UserGameParam true "user、game param"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/team [get]
func (h *Handler) GameTeamPlayerWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		wsConn, conn, err = NewWebsocketConn(ctx)
		isOPen            = true
		// result            UserGameParam
	)
	// fmt.Println("開啟即時回傳雙方隊伍所有玩家、隊長資訊(後端平台)ws api")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" {
		b, _ := json.Marshal(UserGameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 開啟時傳送訊息
		// 取得雙方隊伍玩家資訊、隊長資訊(從redis取得)
		b, _ := json.Marshal(UserGameParam{
			Game: h.getTeamPlayers(gameID),
		})
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, config.CHANNEL_GAME_TEAM_REDIS+gameID, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		// 取得雙方隊伍玩家資訊、隊長資訊(從redis取得)
		b, _ := json.Marshal(UserGameParam{
			Game: h.getTeamPlayers(gameID),
		})
		conn.WriteMessage(b)
	})

	for {
		var (
			result UserGameParam
			// round     int64
			err, err1 error
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉即時回傳雙方隊伍所有玩家、隊長資訊(後端平台)ws api")
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_TEAM_REDIS + gameID)
			conn.Close()
			return
		}

		// gameModel := h.getGameInfo(gameID, "tugofwar") // 遊戲資訊(redis處理)
		// 如果遊戲狀態是start or gaming時，輪次參數需要減1(因為start狀態時輪次已加1)
		// if gameModel.GameStatus == "start" || gameModel.GameStatus == "gaming" {
		// 	round = gameModel.GameRound - 1
		// } else {
		// 	round = gameModel.GameRound
		// }

		// 解碼
		json.Unmarshal(data, &result)

		if result.User.UserID != "" &&
			(result.User.Team == "left_team" || result.User.Team == "right_team") &&
			(result.User.Action == "change" || result.User.Action == "leader" ||
				result.User.Action == "delete" || result.User.Action == "cancel leader") {

			if result.User.Action == "change" {
				// 換隊
				err = models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					ChangePlayer(true, gameID, result.User.UserID, result.User.Team)
			} else if result.User.Action == "leader" {
				// 當隊長
				err = models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateLeader(true, gameID, result.User.UserID, result.User.Team)
			} else if result.User.Action == "cancel leader" {
				// 取消隊長(將隊長資料更新為空值)
				err = models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					UpdateLeader(true, gameID, "", result.User.Team)
			} else if result.User.Action == "delete" {
				// 刪除隊員
				// 將已報名遊戲的人員資料刪除
				// err = models.DefaultGameStaffModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	DeleteUser(gameID, result.User.UserID, round)

				// 遞減隊伍人數(資料庫.redis)
				// err1 = models.DefaultGameModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
				// 	UpdateTeamPeople(true, gameID, result.User.Team, -1, -1)

				// 資料處理(redis)
				models.DefaultGameModel().
					SetConn(h.dbConn, h.redisConn, h.mongoConn).
					DeletePlayer(true, gameID, result.User.Team, result.User.UserID)

				// 判斷是否為隊長
				var (
					teamModel = h.getTeamInfo(gameID) // 隊伍資訊
					leader    string
				)
				if result.User.Team == "left_team" {
					leader = teamModel.LeftTeamLeader
				} else if result.User.Team == "right_team" {
					leader = teamModel.RightTeamLeader
				}

				if result.User.UserID == leader {
					// 清除隊長資訊(redis)
					models.DefaultGameModel().
						SetConn(h.dbConn, h.redisConn, h.mongoConn).
						DeleteLeader(true, gameID, result.User.Team)
				}

				// 清除玩家分數資料
				h.redisConn.ZSetRem(config.SCORES_REDIS+gameID, result.User.UserID)

				// 設置過期時間
				// h.redisConn.SetExpire(config.SCORES_REDIS+gameID, config.REDIS_EXPIRE)

				// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
				h.redisConn.Publish(config.CHANNEL_SCORES_REDIS+gameID, "修改資料")
			}
			if err != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}
			if err1 != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: err1.Error()})
				conn.WriteMessage(b)
				return
			}

			// log.Println("優化後")
			// 回傳隊伍資訊給前端
			// b, _ := json.Marshal(UserGameParam{
			// 	Game: h.getTeamPlayers(gameID),
			// })
			// conn.WriteMessage(b)
		}
	}
}

// @Summary 即時回傳玩家隊伍資訊(玩家端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(tugofwar)
// @Param user_id query string true "user ID"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/user [get]
func (h *Handler) GameUserWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		userID            = ctx.Query("user_id")
		wsConn, conn, err = NewWebsocketConn(ctx)
		isOPen            = true
	)

	// fmt.Println("開啟即時回傳用戶資訊(玩家端)ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" ||
		game == "" || userID == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 開啟時傳送訊息
		var (
			teamModel = h.getTeamInfo(gameID) // 隊伍資訊
			team      string
			isLeader  bool
		)

		// 判斷左方隊伍
		// 判斷用戶的隊伍(redis處理，game_tugofwar_left_team_attend_遊戲ID、game_tugofwar_right_team_attend_遊戲ID)
		if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID, userID) {
			team = "left_team"

			if teamModel.LeftTeamLeader == userID {
				isLeader = true
			}
		}

		// 判斷右方隊伍
		if team == "" {
			if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID, userID) {
				team = "right_team"

				if teamModel.RightTeamLeader == userID {
					isLeader = true
				}
			}
		}

		b, _ := json.Marshal(UserGameParam{
			User: UserModel{
				Team:     team,
				IsLeader: isLeader,
			},
			Game: GameModel{},
		})
		conn.WriteMessage(b)

		isOPen = false
	}

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribe(context, config.CHANNEL_GAME_TEAM_REDIS+gameID, func(channel, message string) {
		// 資料改變時，回傳最新資料至前端
		// 開啟時傳送訊息
		var (
			teamModel = h.getTeamInfo(gameID) // 隊伍資訊
			team      string
			isLeader  bool
		)

		// 判斷左方隊伍
		// 判斷用戶的隊伍(redis處理，game_tugofwar_left_team_attend_遊戲ID、game_tugofwar_right_team_attend_遊戲ID)
		if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS+gameID, userID) {
			team = "left_team"

			if teamModel.LeftTeamLeader == userID {
				isLeader = true
			}
		}

		// 判斷右方隊伍
		if team == "" {
			if h.redisConn.SetIsMember(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS+gameID, userID) {
				team = "right_team"

				if teamModel.RightTeamLeader == userID {
					isLeader = true
				}
			}
		}

		b, _ := json.Marshal(UserGameParam{
			User: UserModel{
				Team:     team,
				IsLeader: isLeader,
			},
			Game: GameModel{},
		})
		conn.WriteMessage(b)
	})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉即時回傳用戶資訊(玩家端)ws")
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_GAME_TEAM_REDIS + gameID)
			conn.Close()
			return
		}
	}
}

// ###優化前，定頻###
// go func() {
// 	for {
// 		var (
// 			teamModel = h.getTeamInfo(gameID) // 隊伍資訊
// 			team      string
// 			isLeader  bool
// 		)

// 		// 判斷左方隊伍
// 		if utils.InArray(teamModel.LeftTeamPlayers, userID) {
// 			team = "left_team"

// 			if teamModel.LeftTeamLeader == userID {
// 				isLeader = true
// 			}
// 		}

// 		// 判斷右方隊伍
// 		if team == "" {
// 			if utils.InArray(teamModel.RightTeamPlayers, userID) {
// 				team = "right_team"

// 				if teamModel.RightTeamLeader == userID {
// 					isLeader = true
// 				}
// 			}
// 		}
// 		// fmt.Println("玩家隊伍資訊: ", team)

// 		b, _ := json.Marshal(UserGameParam{
// 			User: UserModel{
// 				Team:     team,
// 				IsLeader: isLeader,
// 			},
// 			Game: GameModel{},
// 		})
// 		conn.WriteMessage(b)

// 		// ws關閉
// 		if conn.isClose {
// 			return
// 		}

// 		time.Sleep(time.Second * 1)
// 	}
// }()
// ###優化前，定頻###

// if result.User.UserID == "" ||
// 	(result.User.Team != "left_team" && result.User.Team != "right_team") ||
// 	(result.User.Action != "change" && result.User.Action != "leader"&&
// 	  result.User.Action != "delete" && result.User.Action != "cancel leader") {
// 	// fmt.Println("接收參數? ", result)
// 	b, _ := json.Marshal(UserGameParam{
// 		Error: "錯誤: 無法辨識用戶資訊"})
// 	conn.WriteMessage(b)
// 	return
// }

// ###優化前，定頻###
// go func() {
// 	for {
// 		// 取得雙方隊伍玩家資訊、隊長資訊(從redis取得)
// 		b, _ := json.Marshal(UserGameParam{
// 			Game: h.getTeamPlayers(gameID),
// 		})
// 		conn.WriteMessage(b)

//			// ws關閉
//			if conn.isClose {
//				return
//			}
//			time.Sleep(time.Second * 1)
//		}
//	}()
//
// ###優化前，定頻###
// getTeamPlayers 取得雙方隊伍玩家資訊、隊長資訊(從redis取得)
func (h *Handler) getTeamPlayers(gameID string) GameModel {
	var (
		leftTeamPlayers  = make([]UserModel, 0) // 左方隊伍玩家資訊
		rightTeamPlayers = make([]UserModel, 0) // 右方隊伍玩家資訊
		leftTeamLeader   UserModel              // 左方隊伍隊長資訊
		rightTeamLeader  UserModel              // 右方隊伍隊長資訊
	)

	// 取得雙方隊伍資訊
	teamModel, _ := models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn, h.mongoConn).
		FindTeam(true, gameID)

	// fmt.Println("左方隊伍玩家teamModel.LeftTeamPlayers: ", teamModel.LeftTeamPlayers)
	// fmt.Println("右方隊伍玩家: ", teamModel.RightTeamPlayers)
	// fmt.Println("左方隊長: ", teamModel.LeftTeamLeader)
	// fmt.Println("右方隊長: ", teamModel.RightTeamLeader)
	// 處理左方隊伍玩家資訊
	for _, userID := range teamModel.LeftTeamPlayers {
		// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
		user, _ := models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, "", "user_id", userID)

		// 判斷是否為隊長
		if userID == teamModel.LeftTeamLeader {
			leftTeamLeader = UserModel{
				UserID: userID,
				Name:   user.Name,
				Avatar: user.Avatar,
				Team:   "left_team",
			}
		}

		leftTeamPlayers = append(leftTeamPlayers, UserModel{
			UserID: userID,
			Name:   user.Name,
			Avatar: user.Avatar,
			Team:   "left_team",
		})

	}
	// fmt.Println("左方隊伍玩家leftTeamPlayers: ", teamModel.LeftTeamPlayers)
	// 處理右方隊伍玩家資訊
	for _, userID := range teamModel.RightTeamPlayers {
		// 判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
		user, _ := models.DefaultLineModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			Find(true, "", "user_id", userID)

		// 判斷是否為隊長
		if userID == teamModel.RightTeamLeader {
			rightTeamLeader = UserModel{
				UserID: userID,
				Name:   user.Name,
				Avatar: user.Avatar,
				Team:   "right_team",
			}
		}

		rightTeamPlayers = append(rightTeamPlayers, UserModel{
			UserID: userID,
			Name:   user.Name,
			Avatar: user.Avatar,
			Team:   "right_team",
		})

	}

	return GameModel{
		LeftTeamPlayers:  leftTeamPlayers,
		RightTeamPlayers: rightTeamPlayers,
		LeftTeamLeader:   leftTeamLeader,
		RightTeamLeader:  rightTeamLeader,
	}
}
