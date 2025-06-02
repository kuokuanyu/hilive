package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/response"
	"hilive/modules/utils"

	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// DrawNumbersParam 搖號抽獎中獎判斷回傳參數
type DrawNumbersParam struct {
	PrizeStaffs []models.PrizeStaffModel // 中獎人員資訊
	// Numbers     []int64 `json:"numbers" example:"1,2,3"`   // 所有中獎號碼
	PrizeRemain int64 `json:"prize_remain" example:"10"` // 抽獎後剩餘數量
	// Error   string  `json:"error" example:"error message"` // 錯誤訊息
}

// @Summary 搖號抽獎中獎判斷，回傳所有中獎號碼及新的獎品資訊
// @Tags Lucky
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game_id query string true "game ID"
// @Param prize_id query string true "prize ID"
// @Param people query string true "單輪抽獎人數"
// @Success 200 {array} DrawNumbersParam
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /game/draw_numbers [get]
func (h *Handler) GetDrawNumbers(ctx *gin.Context) {
	// log.Println("抽獎開始")
	var (
		host               = ctx.Request.Host
		activityID         = ctx.Query("activity_id")
		gameID             = ctx.Query("game_id")
		prizeID            = ctx.Query("prize_id")
		game               = "draw_numbers"
		peopleStr          = ctx.Query("people")
		signStaffs         = make([]models.ApplysignModel, 0)      // 簽到人員資訊
		noStaffs           = make([]string, 0)                     // 不可抽獎人員
		gameModel          = h.getGameInfo(gameID, "draw_numbers") // 遊戲資訊
		staffs             = make([]models.ApplysignModel, 0)
		numbers            = make([]int64, 0)
		winningStaffparams = []interface{}{config.WINNING_STAFFS_REDIS + gameID}                  // redis參數(不允許重複中獎時需要添加中獎人員進redis中)
		allStaffparams     = []interface{}{config.DRAW_NUMBERS_WINNING_STAFFS_REDIS + activityID} // redis參數(所有搖號抽獎場次中獎人員資料)
		wg                 sync.WaitGroup
		lock               sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
		paramLock          sync.Mutex // 宣告Lock用以資源佔有與解鎖(避免共用變數計算錯誤)
		messageAmount      int64      // 簡訊數量
		// now, _        = time.ParseInLocation("2006-01-02 15:04:05",
		// 	time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local)
	)

	// startTime := time.Now()
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return
	}
	if activityID == "" || gameID == "" || prizeID == "" {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識活動、遊戲、獎品資訊，請輸入有效的ID",
		})
		return
	}
	people, err := strconv.Atoi(peopleStr)
	if err != nil || people == 0 {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法辨識中獎人數資訊，請輸入有效的中獎人數",
		})
		return
	}

	// 獎品資訊
	prizeModel, err := models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindPrize(false, prizeID)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// 簽到人員資訊
	if signStaffs, err = h.getSignStaffs(false, false, "", activityID, 0, 0); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// 取得所有搖號抽獎場次中獎人員資料，確保寫入redis中
	_, err = h.getDrawNumbersAllWinningStaffs(true, activityID)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return
	}

	// log.Println("目前所有搖號抽獎場次中獎人數: ", len(allPrizeStaffs))

	// 取得黑名單人員資料
	blackStaffs, _ := models.DefaultBlackStaffModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		FindAll(true, activityID, gameID, "draw_numbers")
	// 加黑名單人員加入不可抽獎人員中
	for _, staffModel := range blackStaffs {
		noStaffs = append(noStaffs, staffModel.UserID)
	}

	// 不允許重複中獎，將簽到人員裡已中獎人員資料去除
	if gameModel.AllGameAllow == "close" || gameModel.DrawNumbersGameAllow == "close" ||
		gameModel.Allow == "close" {
		// 不允許重複中獎，取得已中獎人員資料
		prizeStaffs, _ := h.getDrawNumbersWinningStaffs(true, activityID, gameID, "draw_numbers",
			gameModel.AllGameAllow == "open", gameModel.DrawNumbersGameAllow == "open",
			gameModel.Allow == "open") // 中獎人員資料

		// 將已中獎人員資料加入不可抽獎人員中
		noStaffs = utils.AddUniqueStrings(noStaffs, prizeStaffs)

		// log.Println("不能抽獎的人數: ", len(noStaffs))
	}

	// log.Println("GetDrawNumbers 執行可抽獎人員判斷")
	for _, staff := range signStaffs {
		// 判斷用戶ID是否為不可抽獎人員
		if utils.InArray(noStaffs, staff.UserID) {
			// 不可抽獎人員
		} else {
			// 可抽獎人員
			staffs = append(staffs, staff)
		}
	}
	// log.Println("GetDrawNumbers 執行可抽獎人員判斷結束")
	// log.Println("GetDrawNumbers 可抽獎人數: ", len(staffs))

	// 抽到的index值(不是抽中的號碼，而是先抽報名簽到人員陣列中的index值)
	for i := 0; len(numbers) < people; i++ {
		// 隨機值
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(len(staffs))

		// 判斷是否已經抽過該值，如果沒抽過則加入陣列中
		if !utils.IntInArray(numbers, int64(random)) {
			numbers = append(numbers, int64(random))
		}
	}

	// 取得簡訊數量
	activityModel, err := models.DefaultActivityModel().
	SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Find(false, activityID)
	if err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法取得活動資訊，請重新查詢",
		})
		return
	}
	messageAmount = activityModel.MessageAmount

	// 發獎
	prizeStaffs := make([]models.PrizeStaffModel, people) // 中獎人員資訊

	// 批量寫入所有中獎人員資料
	// log.Println("批量寫入所有中獎人員資料: ", len(numbers))

	var (
		userIDs  = make([]string, len(numbers)) // 用戶ID
		peizeIDs = make([]string, len(numbers)) // 獎品ID
		scores1  = make([]string, len(numbers)) // 分數(空)
		scores2  = make([]string, len(numbers)) // 分數2(空)
		ranks    = make([]string, len(numbers)) // 排名資訊(空)
		teams    = make([]string, len(numbers)) // 隊伍資訊(空)
		leaders  = make([]string, len(numbers)) // 隊長資訊(空)
		mvps     = make([]string, len(numbers)) // mvp資訊(空)
	)

	// 處理批量寫入資料的相關參數
	for i := 0; i < len(numbers); i++ {
		userIDs[i] = staffs[numbers[i]].UserID
		peizeIDs[i] = prizeID
		scores1[i] = strconv.Itoa(0)
		scores2[i] = strconv.Itoa(0)
		ranks[i] = strconv.Itoa(0)
	}

	// 將所有中獎人員資料批量寫入資料表中
	if err = models.DefaultPrizeStaffModel().
	SetConn(h.dbConn, h.redisConn,h.mongoConn).
		Adds(len(numbers), activityID, gameID, game, userIDs, peizeIDs,
		int(gameModel.GameRound), scores1, scores2, ranks, teams, leaders, mvps); err != nil {
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

	// log.Println("更新獎品數量")

	// 更新獎品數量
	if err = models.DefaultPrizeModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		UpdateRemain(gameID, prizeID, prizeModel.PrizeRemain-int64(people)); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 獎品無庫存，請重新抽獎",
		})
		return
	}

	// log.Println("更新獎品數量完成")

	// 更新遊戲輪次
	if err = models.DefaultGameModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
		UpdateGameStatus(true, gameID, "+1",
			"", ""); err != nil {
		response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  "",
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 無法更新遊戲狀態",
		})
		return
	}

	// log.Println("更新輪次完成")

	// 多線呈處理簡訊判斷、中獎人員資料處理
	wg.Add(len(numbers)) //計數器
	for i := 0; i < len(numbers); i++ {
		go func(i int) {
			// log.Println("i: ", i)
			defer wg.Done()
			staff := staffs[numbers[i]] // 中獎人員

			// 用戶遊戲紀錄(將中獎紀錄寫入redis)
			// record := models.PrizeStaffModel{
			// 	ID:            0,
			// 	ActivityID:    activityID,
			// 	GameID:        gameID,
			// 	PrizeID:       prizeID,
			// 	UserID:        staff.UserID,
			// 	Name:          staff.Name,
			// 	Avatar:        staff.Avatar,
			// 	Game:          "draw_numbers",
			// 	PrizeName:     prizeModel.PrizeName,
			// 	PrizeType:     prizeModel.PrizeType,
			// 	PrizePicture:  prizeModel.PrizePicture,
			// 	PrizePrice:    prizeModel.PrizePrice,
			// 	PrizeMethod:   prizeModel.PrizeMethod,
			// 	PrizePassword: prizeModel.PrizePassword,
			// 	Round:         gameModel.GameRound - 1, // 輪次已遞增
			// 	WinTime:       now.Format("2006-01-02 15:04:05"),
			// 	Status:        "no",
			// 	Score:         0,
			// 	Score2:        0,
			// 	Rank:          0,
			// }

			// 用戶遊戲紀錄(包含中獎與未中獎)
			// allRecords, _, _, _, _, err := h.getUserGameRecords(activityID, gameID, "draw_numbers", staff.UserID)
			// if err != nil {
			// 	// log.Println(err)
			// 	return
			// }
			// allRecords = append(allRecords, record)

			// 更新redis裡用戶遊戲紀錄
			// h.redisConn.HashSetCache(config.USER_GAME_RECORDS_REDIS+activityID,
			// 	staff.UserID, utils.JSON(allRecords))

			// 設置過期時間
			// h.redisConn.SetExpire(config.USER_GAME_RECORDS_REDIS+activityID, config.REDIS_EXPIRE)

			// 中獎人員添加至陣列中(回傳給前端的資料)
			prizeStaffs[i] = models.PrizeStaffModel{
				UserID: staff.UserID,
				Name:   staff.Name,
				Avatar: staff.Avatar,
				Number: staff.Number,
			}

			paramLock.Lock()                                              //佔有資源
			winningStaffparams = append(winningStaffparams, staff.UserID) // 將中獎人員加入redis中(不允許重複中獎時須執行)

			// 判斷中獎人員是否已經在redis中(DRAW_NUMBERS_WINNING_STAFFS_REDIS)
			// if !utils.InArray(allPrizeStaffs, staff.UserID) {
			allStaffparams = append(allStaffparams, staff.UserID)
			// }
			paramLock.Unlock() //釋放資源

			// for l := 0; l < MaxRetries; l++ {
			// 上鎖
			// ok, _ := h.acquireLock(config.LUCKY_DRAW_NUMBERS_LOCK_REDIS+gameID, LockExpiration)
			// if ok == "OK" {
			// 釋放鎖
			// defer h.releaseLock(config.LUCKY_DRAW_NUMBERS_LOCK_REDIS + gameID)

			// 傳送中獎簡訊
			// %s 您好: 恭喜你於 %s 活動遊戲中獲得%s獎品，請聯繫主辦單位領取獎品，謝謝
			if staff.Phone != "" && activityModel.PushPhoneMessage == "open" && messageAmount > 0 {
				lock.Lock() //佔有資源
				// 遞減訊息數量
				messageAmount--
				lock.Unlock() //釋放資源

				message := fmt.Sprintf("%s 您好: 恭喜你於 %s 活動遊戲中獲得%s獎品，請聯繫主辦單位領取獎品，謝謝",
					staff.Name, prizeModel.ActivityName, prizeModel.PrizeName)

				// 發送簡訊
				err = sendMessage(staff.Phone, message)
				if err != nil {
					response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  "",
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: err.Error(),
					})
				}
			}

			// 釋放鎖
			// h.releaseLock(config.LUCKY_DRAW_NUMBERS_LOCK_REDIS + gameID)
			// break
			// }

			// 鎖被佔用，稍微延遲後重試
			// time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
			// }
		}(i)
	}
	wg.Wait() //等待計數器歸0

	// log.Println("要新增至redis的人數): ", len(allStaffparams)-1, len(winningStaffparams)-1)

	// 有開啟簡訊功能才更新數量
	if activityModel.PushPhoneMessage == "open" {
		// 更新簡訊數量
		if err = models.DefaultActivityModel().
		SetConn(h.dbConn, h.redisConn,h.mongoConn).
			UpdateActivity(
			true, models.EditActivityModel{
				ActivityID:    activityID,
				MessageAmount: strconv.Itoa(int(messageAmount)),

				// 避免更新後以下數值被清空
				LineID:        activityModel.ActivityLineID,
				ChannelID:     activityModel.ActivityChannelID,
				ChannelSecret: activityModel.ActivityChannelSecret,
				ChatbotSecret: activityModel.ActivityChatbotSecret,
				ChatbotToken:  activityModel.ActivityChatbotToken,
			}); err != nil {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  "",
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 更新簡訊數量發生問題",
			})
			return
		}
	}

	// log.Println("多線呈結束")

	// 不允許重複中獎，將中獎人員加入redis中(WINNING_STAFFS_REDIS)
	if gameModel.AllGameAllow == "close" || gameModel.DrawNumbersGameAllow == "close" ||
		gameModel.Allow == "close" {
		h.redisConn.SetAdd(winningStaffparams)

		// 設置過期時間
		// h.redisConn.SetExpire(config.WINNING_STAFFS_REDIS+gameID, config.REDIS_EXPIRE)

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		h.redisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+gameID, "修改資料")
	}

	if len(allStaffparams)-1 > 0 { // 陣列裡包含redis名稱，實際數量應減一
		// 將所有搖號抽獎場次中獎人員資料加入redis中(DRAW_NUMBERS_WINNING_STAFFS_REDIS)
		h.redisConn.SetAdd(allStaffparams)
	}

	// allPrizeStaffs, _ := h.getDrawNumbersAllWinningStaffs(true, activityID)
	// log.Println("redis新增後的的資料數: ", len(allPrizeStaffs))

	// 刪除玩家遊戲紀錄(中獎.未中獎)，不確定redis是否有所有用戶的遊戲紀錄
	// 搖號抽獎執行沒有判斷角色api，如果有執行判斷角色api那redis裡一定會有玩家紀錄
	// h.redisConn.DelCache(config.USER_GAME_RECORDS_REDIS + activityID)

	// endTime := time.Now()
	// log.Println("統計中獎人員花費時間: ", endTime.Sub(startTime))

	// 回傳中獎號碼與新的獎品資訊
	response.OkWithData(ctx, DrawNumbersParam{
		PrizeStaffs: prizeStaffs,
		// Numbers:     numbers,
		PrizeRemain: prizeModel.PrizeRemain - int64(people),
	})
}

// 新增中獎人員名單
// if _, err = models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
// 	Add(models.NewPrizeStaffModel{
// 		UserID:     staff.UserID,
// 		ActivityID: activityID,
// 		GameID:     gameID,
// 		PrizeID:    prizeID,
// 		Round:      strconv.Itoa(int(gameModel.GameRound)),
// 		Status:     "no",
// 		Score:      0,
// 		Score2:     0,
// 		Rank:       0,
// 		// White:      "no",
// 	}); err != nil {
// 	response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
// 		UserID:  "",
// 		Method:  ctx.Request.Method,
// 		Path:    ctx.Request.URL.Path,
// 		Message: "錯誤: 新增中獎人員發生問題，請重新操作",
// 	})
// 	return
// }

// newStaffs := make([]models.ApplysignModel, 0) // 簽到人員資訊(整理過後)
// for _, staff := range staffs {
// 	if !h.redisConn.SetIsMember(config.WINNING_STAFFS_REDIS+gameID, staff.UserID) {
// 		// 未中獎人員，加入抽獎陣列中
// 		newStaffs = append(newStaffs, staff)
// 	}
// }
// staffs = newStaffs
// fmt.Println("不允許重複中獎，處理簽到人員後的抽獎人數: ", len(staffs))

// 判斷是否為黑名單
// for _, staff := range signStaffs {
// 	if !h.IsBlackStaff(activityID, gameID, "draw_numbers", staff.UserID) {
// 		// 不為黑名單
// 		staffs = append(staffs, staff)
// 	}
// }

// numbers[i] = staff.Number             // 抽到的號碼，更新至陣列中

// people  int
// prize      models.PrizeModel                  // 獎品資訊

// -----如果這一段威翔有做判斷，可以拿掉
// 參加抽獎人數>=要發放的獎品數
// if len(staffs) < people {
// 	response.BadRequest(ctx, "錯誤: 活動簽到人數必須大於或等於要發放獎品數，請重新執行")
// 	return
// }

// prizeModel, err := models.DefaultPrizeModel().
// 	SetDbConn(h.dbConn).FindPrize(prizeID)
// if err != nil || prizeModel.ID == 0 {
// 	response.BadRequest(ctx, "錯誤: 無法取得獎品資訊，請重新查詢")
// 	return
// }
// 後端獎品數>=要發放的獎品數
// if prizeModel.Remain < int64(people) {
// 	response.BadRequest(ctx, "錯誤: 後端的獎品數必須大於或等於要發放的獎品數，請重新執行")
// 	return
// }
// -----如果這一段威翔有做判斷，可以拿掉

// DrawNumbersPrize 資料表欄位
// type DrawNumbersPrize struct {
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	PrizeID    string `json:"prize_id" example:"prize_id"`
// 	Name       string `json:"name" example:"獎品名稱"`
// 	PrizeType  string `json:"prize_type" example:"獎品類型"`
// 	Picture    string `json:"picture" example:"獎品圖片"`
// 	Amount     int64  `json:"amount" example:"5"`
// 	Remain     int64  `json:"remain" example:"5"`
// 	Price      int64  `json:"price" example:"1000"`
// 	Method     string `json:"method" example:"兌換方式"`
// 	Password   string `json:"password" example:"兌獎密碼"`
// }

// Name:       staff.Name,
// Avatar:     staff.Avatar,
// PrizeName:  prizeModel.Name,
// PrizeType:  prizeModel.PrizeType,
// Picture:    prizeModel.Picture,
// Price:      strconv.Itoa(int(prizeModel.Price)),
// Method:     prizeModel.Method,
// Password:   prizeModel.Password,
// WinTime:    now.Format("2006-01-02 15:
