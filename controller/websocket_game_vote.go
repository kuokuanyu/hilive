package controller

import (
	"context"
	"encoding/json"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// HostVoteWriteParam 主持端投票遊戲回傳前端訊息
type HostVoteWriteParam struct {
	Game  GameModel // 場次資訊
	Error string    `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 主持端取得投票選項排名資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(vote)
// @Success 200 {array} HostVoteWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/host/vote [get]
func (h *Handler) HostGameVoteWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		isOPen            = true
	)

	defer wsConn.Close()
	defer conn.Close()

	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	if isOPen {
		// 回傳前端之後將舊的redis資料清除，下一次回傳前端時只需要傳新的
		h.redisConn.DelCache(config.VOTE_AVATAR_REDIS + gameID)

		// 開啟時傳送訊息
		var (
			voteOptions = make([]GameVoteOptionModel, 0)
		)

		// 取得投票選項資訊(redis取得)
		voteOptionModels, _ := models.DefaultGameVoteOptionModel().
			SetConn(h.dbConn, h.redisConn, h.mongoConn).
			FindOrderByScore(true, gameID)

		// 只傳option_id跟option_score，介紹資訊在一開始就有傳送了
		for _, optionModel := range voteOptionModels {
			voteOptions = append(voteOptions, GameVoteOptionModel{
				OptionID:    optionModel.OptionID,
				OptionScore: optionModel.OptionScore,
			})
		}

		// 傳遞資料至前端
		conn.WriteMessage([]byte(utils.JSON(
			HostVoteWriteParam{
				Game: GameModel{
					VoteAvatars: []GameVoteAvatarModel{},
					VoteOptions: voteOptions,
				}})))

		isOPen = false
	}

	// #####加入測試資料#####start
	// if gameID == "pk4BGSyzBIXNyCn792fC" || gameID == "B8SuaaQsZ78r6VDSnFBH" {
	// 	// 取得投票選項資訊
	// 	voteOptionModels, _ := models.DefaultGameVoteOptionModel().
	// 		SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 		FindOrderByID(true, gameID)

	// 	go func() {
	// 		for {
	// 			avatars := []string{
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/1.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/2.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/3.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/4.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/5.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/6.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/7.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/8.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/9.png",
	// 				"https://dev.hilives.net/admin/uploads/system/fake_data/10.png",
	// 			}

	// 			for _, option := range voteOptionModels {
	// 				// 將投票紀錄、用戶頭像資料寫入redis中(LIST，config.VOTE_AVATAR_REDIS+gameID)
	// 				h.redisConn.ListRPush(config.VOTE_AVATAR_REDIS+gameID,
	// 					utils.JSON(GameVoteAvatarModel{
	// 						OptionID: option.OptionID,
	// 						Avatar:   avatars[int64(rand.Intn(10))],
	// 					}))

	// 				// 每個投票資料加入隨機票數
	// 				// 遞增投票選項分數資料
	// 				models.DefaultGameVoteOptionModel().SetDbConn(h.dbConn).SetRedisConn(h.redisConn).
	// 					IncrOptionScore(gameID, option.OptionID, int64(rand.Intn(10)))
	// 			}

	// 			// ws關閉
	// 			if conn.isClose {
	// 				return
	// 			}

	// 			time.Sleep(time.Second * 1)
	// 		}
	// 	}()
	// }
	// #####加入測試資料#####end

	// 用來控制 Redis 訂閱開關判斷
	context, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel() // 確保函式退出時取消訂閱

	// 啟用redis訂閱
	go h.redisConn.Subscribes(context,
		[]string{config.CHANNEL_SCORES_REDIS + gameID}, func(channel, message string) {
			// 資料改變時，回傳最新資料至前端
			var (
				voteOptions = make([]GameVoteOptionModel, 0)
				voteAvatars = make([]GameVoteAvatarModel, 0)
			)

			// 取得投票選項資訊(redis取得)
			voteOptionModels, _ := models.DefaultGameVoteOptionModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindOrderByScore(true, gameID)

			// 只傳option_id跟option_score，介紹資訊在一開始就有傳送了
			for _, optionModel := range voteOptionModels {
				voteOptions = append(voteOptions, GameVoteOptionModel{
					OptionID:    optionModel.OptionID,
					OptionScore: optionModel.OptionScore,
				})
			}

			// 取得投票投像紀錄(config.VOTE_AVATAR_REDIS+gameID，LIST)
			avatarsJson, _ := h.redisConn.ListRange(config.VOTE_AVATAR_REDIS+gameID, 0, 0)
			if len(avatarsJson) > 0 {
				for _, avatar := range avatarsJson {
					var voteAvatar GameVoteAvatarModel
					// 解碼
					json.Unmarshal([]byte(avatar), &voteAvatar)

					voteAvatars = append(voteAvatars, voteAvatar)
				}
			}
			// 回傳前端之後將舊的redis資料清除，下一次回傳前端時只需要傳新的
			h.redisConn.DelCache(config.VOTE_AVATAR_REDIS + gameID)

			// 傳遞資料至前端
			conn.WriteMessage([]byte(utils.JSON(
				HostVoteWriteParam{
					Game: GameModel{
						VoteOptions: voteOptions,
						VoteAvatars: voteAvatars,
					}})))
		})

	for {
		_, err := conn.ReadMessage()
		if err != nil {
			// 取消訂閱
			h.redisConn.Unsubscribe(config.CHANNEL_SCORES_REDIS + gameID)

			// 回傳前端之後將舊的redis資料清除，下一次回傳前端時只需要傳新的
			h.redisConn.DelCache(config.VOTE_AVATAR_REDIS + gameID)
			conn.Close()
			return
		}
	}
}

// @Summary 玩家端接收前端投票紀錄資訊
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param game query string true "game" Enums(vote)
// @param body body UserGameParam true "user、game param"
// @Success 200 {array} UserGameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game/guest/vote [get]
func (h *Handler) GuestGameVoteWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		gameID            = ctx.Query("game_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            UserGameParam
	)
	// fmt.Println("開啟玩家端更新總答題人數、各選項答題人數、答題紀錄資訊ws")
	defer wsConn.Close()
	defer conn.Close()
	if err != nil || activityID == "" || gameID == "" || game == "" {
		b, _ := json.Marshal(UserGameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result UserGameParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Game.UserID != "" &&
			result.Game.OptionID != "" && result.Game.Score != 0 {

			// 取得投票輪次
			gameModel := h.getGameInfo(gameID, game) // 即時輪次

			// 新增用戶投票紀錄
			err = models.DefaultGameVoteRecordModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(true, models.EditGameVoteRecordModel{
					ActivityID: activityID,
					GameID:     gameID,
					UserID:     result.Game.UserID,
					OptionID:   result.Game.OptionID,
					Score:      result.Game.Score,
					Round:      gameModel.GameRound,
				})
			if err != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: "錯誤: 新增用戶投票紀錄發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 將投票紀錄、用戶頭像資料寫入redis中(LIST，config.VOTE_AVATAR_REDIS+gameID)
			// 取得用戶頭像，判斷redis裡是否有用戶資訊，有則不執行查詢資料表功能
			// user, _ := models.DefaultLineModel().SetDbConn(h.dbConn).
			// 	SetRedisConn(h.redisConn).Find(true, "", "user_id", result.Game.UserID)

			h.redisConn.ListRPush(config.VOTE_AVATAR_REDIS+gameID,
				utils.JSON(GameVoteAvatarModel{
					OptionID: result.Game.OptionID,
					Avatar:   result.User.Avatar,
				}))

			// 設置過期時間
			// h.redisConn.SetExpire(config.VOTE_AVATAR_REDIS+gameID, config.REDIS_EXPIRE)

			// 遞增投票選項分數資料(會觸發pub/sub，config.CHANNEL_SCORES_REDIS+gameID)
			if err = models.DefaultGameVoteOptionModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				IncrOptionScore(gameID, result.Game.OptionID, result.Game.Score); err != nil {
				b, _ := json.Marshal(UserGameParam{
					Error: "錯誤: 更新投票選項分數資料發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 傳遞資料至前端
			conn.WriteMessage([]byte(utils.JSON(
				UserGameParam{})))
		}
	}
}
