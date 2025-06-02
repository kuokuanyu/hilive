package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/file"
	"hilive/modules/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// PUT 編輯 PUT API
func (h *Handler) PUT(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// contentType = ctx.Request.Header.Get("content-type")
		path        = ctx.Request.URL.Path
		prefix      = ctx.Param("__prefix")
		values      = make(map[string][]string)
		activityID  = ctx.Request.FormValue("activity_id")
		gameID      = ctx.Request.FormValue("game_id")
		tokenUserID = ctx.Request.FormValue("user_id")
		token       = ctx.Request.FormValue("token")
		userID      string
		// err           error
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
		// 因為管理員都可以幫忙設置(獎品頁面例外), 所以需要區分user_id參數
		userID = ctx.Request.FormValue("user") // 該活動場次的管理員資料

		// 判斷是否為獎品頁面, 獎品頁面目前沒有user參數
		if strings.Contains(path, "/prize/form") {
			userID = ctx.Request.FormValue("user_id")
		}
	} else {
		userID = ctx.Request.FormValue("user_id")
	}

	if prefix == "" {
		if strings.Contains(path, "/v1/user") {
			prefix = "user"
		} else if strings.Contains(path, "/v1/line_user") {
			prefix = "line_user"
		} else if strings.Contains(path, "activity") {
			prefix = "activity"
		} else if strings.Contains(path, "user") {
			prefix = "user"
			// } else if strings.Contains(path, "/v1/applysign/user") {
			// prefix = "applysign_user"
		} else if strings.Contains(path, "/v1/applysign") {
			prefix = "applysign"
		} else if strings.Contains(path, "/v1/chatroom/record") {
			prefix = "chatroom_record"
		} else if strings.Contains(path, "/v1/question/record") {
			prefix = "question_record"
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

	// 處理表單參數
	param, err := ctx.MultipartForm()
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: fmt.Sprintf("錯誤: 表單參數發生問題，請重新操作，%s", err.Error()),
		})

		// 記錄請求的關鍵資訊
		// log.Println("Request Headers:", ctx.Request.Header)                  // multipart/form-data
		// log.Println("Content-Length:", ctx.Request.ContentLength)            // 不能為0
		// log.Println("Content-Type:", ctx.Request.Header.Get("Content-Type")) // Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryABC123

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Headers: %s, Content-Length: %d, Content-Type: %s",
		// 			ctx.Request.Header, ctx.Request.ContentLength, ctx.Request.Header.Get("Content-Type")),
		// 	})

		// // 嘗試讀取部分 Body，確認是否真的有請求內容
		// bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		// log.Println("Request Body (first 500 bytes):", string(bodyBytes[:500]))

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Body (first 300 bytes): %s", string(bodyBytes[:300])),
		// 	})

		// ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20) // 10MB 限制
		// err := ctx.Request.ParseMultipartForm(10 << 20)                              // 10MB
		// if err != nil {
		// 	log.Println("解析表單時發生錯誤(10MB):", err)

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: fmt.Sprintf("解析表單時發生錯誤(10MB): %s", err.Error()),
		// 		})
		// }

		// if ctx.Request.MultipartForm == nil {
		// 	log.Println("Multipart 表單為空(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "Multipart 表單為空(10MB)",
		// 		})
		// } else {
		// 	fmt.Println("成功解析 Multipart 表單(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "成功解析 Multipart 表單(10MB)",
		// 		})
		// }

		return
	}

	// log.Println("圖片處理")
	// 上傳圖片、檔案
	if len(param.File) > 0 {
		if err = file.GetFileEngine(config.FILE_ENGINE).Upload(ctx.Request.MultipartForm, path,
			userID, activityID, gameID, prefix); err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 上傳檔案發生問題，請重新上傳檔案",
			})
			return

		}
	}
	// log.Println("圖片處理完成")

	// 兌獎、編輯用戶不需要token驗證
	// 修改報名簽到用戶姓名頭像api不需要token驗證
	if prefix != "winning" && prefix != "user" && prefix != "line_user" {
		if !auth.CheckToken(token, tokenUserID) {
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: Token驗證發生問題，請輸入有效的Token值",
			})
			return

		}
	}

	// log.Println("更新")
	values = param.Value
	table, _ := h.GetTable(ctx, prefix)
	if err := table.UpdateData(values); err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return

	}
	// log.Println("完成")
	response.Ok(ctx)
}

// PATCH 更新基本設置 PATCH API
func (h *Handler) PATCH(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// contentType                            = ctx.Request.Header.Get("content-type")
		path        = ctx.Request.URL.Path
		prefix      = ctx.Param("__prefix")
		table, _    = h.GetTable(ctx, prefix)
		values      = make(map[string][]string)
		activityID  = ctx.Request.FormValue("activity_id")
		tokenUserID = ctx.Request.FormValue("user_id")
		token       = ctx.Request.FormValue("token")
		userID      string
	)
	if host != config.API_URL {
		// logger.LoggerError(ctx, "錯誤: 網域請求發生問題")
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return

	}

	if path == "/v1/admin/manager" || path == "/v1/user" ||
		path == "/v1/activity" || strings.Contains(path, "/v1/interact/game") {
		// 因為管理員都可以幫忙設置, 所以需要區分user_id參數
		userID = ctx.Request.FormValue("user")
	} else {
		userID = ctx.Request.FormValue("user_id")
	}

	err := ctx.Request.ParseMultipartForm(1 << 20) // 限制最大 1MB
	if err != nil {
		fmt.Println("表單解析錯誤:", err)
	}

	// 處理表單參數
	param, err := ctx.MultipartForm()
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: fmt.Sprintf("錯誤: 表單參數發生問題，請重新操作，%s", err.Error()),
		})

		// 記錄請求的關鍵資訊
		// log.Println("Request Headers:", ctx.Request.Header)                  // multipart/form-data
		// log.Println("Content-Length:", ctx.Request.ContentLength)            // 不能為0
		// log.Println("Content-Type:", ctx.Request.Header.Get("Content-Type")) // Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryABC123

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Headers: %s, Content-Length: %d, Content-Type: %s",
		// 			ctx.Request.Header, ctx.Request.ContentLength, ctx.Request.Header.Get("Content-Type")),
		// 	})

		// // 嘗試讀取部分 Body，確認是否真的有請求內容
		// bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		// log.Println("Request Body (first 500 bytes):", string(bodyBytes[:500]))

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Body (first 300 bytes): %s", string(bodyBytes[:300])),
		// 	})

		// ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20) // 10MB 限制
		// err := ctx.Request.ParseMultipartForm(10 << 20)                              // 10MB
		// if err != nil {
		// 	log.Println("解析表單時發生錯誤(10MB):", err)

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: fmt.Sprintf("解析表單時發生錯誤(10MB): %s", err.Error()),
		// 		})
		// }

		// if ctx.Request.MultipartForm == nil {
		// 	log.Println("Multipart 表單為空(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "Multipart 表單為空(10MB)",
		// 		})
		// } else {
		// 	fmt.Println("成功解析 Multipart 表單(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "成功解析 Multipart 表單(10MB)",
		// 		})
		// }

		return

	}

	// 上傳圖片、檔案
	if len(param.File) > 0 {
		if err := file.GetFileEngine(config.FILE_ENGINE).Upload(ctx.Request.MultipartForm, path,
			userID, activityID, "", ""); err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 上傳檔案發生問題，請重新上傳檔案",
			})
			return

		}
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

	values = param.Value
	if err := table.UpdateSettingData(values); err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return

	}
	response.Ok(ctx)
}

// @Summary 編輯提問紀錄
// @Tags Question_Record
// @version 1.0
// @Accept  mpfd
// @param id formData integer true "提問紀錄ID"
// @param activity_id formData string true "活動ID"
// @param message_status formData string false "訊息審核狀態" Enums(yes, no, review)
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /question/record [put]
func (h *Handler) PUTQuestionRecord(ctx *gin.Context) {
}

// @@@Summary 編輯自定義報名簽到人員資料
// @@@Tags ApplySign
// @@@version 1.0
// @@@Accept  mpfd
// @@@param applysign_user_id formData string true "自定義報名簽到用戶ID"
// @@@param activity_id formData string true "活動ID"
// @@@param name formData string false "name"
// @@@param phone formData string false "phone"
// @@@param ext_email formData string false "ext_email"
// @@@param ext_password formData string true "ext_password"
// @@@param user_id formData string true "用戶ID(辨識token用)"
// @@@param token formData string true "CSRF Token"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /applysign/user [put]
func (h *Handler) PUTApplysignUser(ctx *gin.Context) {
}

// @Summary 遊戲資料重置(form-data)
// @Tags Game
// @version 1.0
// @Accept  mpfd
// @param game_id formData string true "遊戲ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/reset/form [put]
func (h *Handler) PUTReset(ctx *gin.Context) {
}

// @@@Summary 編輯遊戲人員資料(form-data)
// @@@Tags Attend Staff
// @@@version 1.0
// @@@Accept  mpfd
// @@@param id formData string true "ID"
// @@@param game_id formData string true "遊戲ID"
// @@@param user_id formData string true "用戶ID"
// @@@param status f

// else if contentType == "application/json" {
// if prefix == "setting" {
// 	// 遊戲基本設置
// 	var model models.EditGameSettingModel
// 	err = ctx.BindJSON(&model)
// 	userID = model.UserID
// 	token = model.Token

// 	values["activity_id"] = []string{model.ActivityID}
// 	values["lottery_game_allow"] = []string{model.LotteryGameAllow}
// 	values["redpack_game_allow"] = []string{model.RedpackGameAllow}
// 	values["ropepack_game_allow"] = []string{model.RopepackGameAllow}
// 	values["whack_mole_game_allow"] = []string{model.WhackMoleGameAllow}
// 	values["monopoly_game_allow"] = []string{model.MonopolyGameAllow}
// 	values["qa_game_allow"] = []string{model.QAGameAllow}
// 	values["draw_numbers_game_allow"] = []string{model.DrawNumbersGameAllow}
// 	values["all_game_allow"] = []string{model.AllGameAllow}
// } else if prefix == "black" {
// 	// 黑名單
// 	var model models.EditBlackStaffModel
// 	err = ctx.BindJSON(&model)
// 	userID = model.UserID
// 	token = model.Token

// 	values["activity_id"] = []string{model.ActivityID}
// 	values["game_id"] = []string{model.GameID}
// 	values["game"] = []string{model.Game}
// 	values["user_id"] = []string{model.UserID}
// 	values["line_id"] = []string{model.LINEID}
// 	values["reason"] = []string{model.Reason}
// } else if prefix == "winning" {
// 	// 中獎人員
// 	var model models.EditPrizeStaffModel
// 	err = ctx.BindJSON(&model)

// 	values["id"] = []string{model.ID}
// 	values["status"] = []string{model.Status}
// 	// values["white"] = []string{model.White}
// 	values["role"] = []string{model.Role}
// 	values["password"] = []string{model.Password}
// } else if strings.Contains(prefix, "_prize") {
// 	// 獎品相關參數
// 	var model models.EditPrizeModel
// 	err = ctx.BindJSON(&model)
// 	userID = model.UserID
// 	token = model.Token

// 	values["game_id"] = []string{model.GameID}
// 	values["prize_id"] = []string{model.PrizeID}
// 	values["prize_name"] = []string{model.PrizeName}
// 	values["prize_type"] = []string{model.PrizeType}
// 	values["prize_picture"] = []string{model.PrizePicture}
// 	values["prize_amount"] = []string{model.PrizeAmount}
// 	values["prize_price"] = []string{model.PrizePrice}
// 	values["prize_method"] = []string{model.PrizeMethod}
// 	values["prize_password"] = []string{model.PrizePassword}
// } else if !strings.Contains(prefix, "_prize") {
// 	// 遊戲相關參數
// 	var model models.EditGameModel
// 	err = ctx.BindJSON(&model)
// 	userID = model.UserID
// 	token = model.Token

// 	values["game_id"] = []string{model.GameID}
// 	values["title"] = []string{model.Title}
// 	values["game_type"] = []string{model.GameType}
// 	values["limit_time"] = []string{model.LimitTime}
// 	values["second"] = []string{model.Second}
// 	values["max_people"] = []string{model.MaxPeople}
// 	values["people"] = []string{model.People}
// 	values["max_times"] = []string{model.MaxTimes}
// 	values["allow"] = []string{model.Allow}
// 	values["percent"] = []string{model.Percent}
// 	values["first_prize"] = []string{model.FirstPrize}
// 	values["second_prize"] = []string{model.SecondPrize}
// 	values["third_prize"] = []string{model.ThirdPrize}
// 	values["general_prize"] = []string{model.GeneralPrize}
// 	values["topic"] = []string{model.Topic}
// 	values["skin"] = []string{model.Skin}
// 	values["music"] = []string{model.Music}
// 	values["display_name"] = []string{model.DisplayName}

// 	// 敲敲樂自定義
// 	values["whack_mole_host_background"] = []string{model.WhackMoleHostBackground}
// 	values["whack_mole_guest_background"] = []string{model.WhackMoleGuestBackground}
// 	values["whack_mole_dollar_rat_picture"] = []string{model.WhackMoleDollarRatPicture}
// 	values["whack_mole_redpack_rat_picture"] = []string{model.WhackMoleRedpackRatPicture}
// 	values["whack_mole_bomb_picture"] = []string{model.WhackMoleBombPicture}
// 	values["whack_mole_rat_hole_picture"] = []string{model.WhackMoleRatHolePicture}
// 	values["whack_mole_rock_picture"] = []string{model.WhackMoleRockPicture}
// 	values["whack_mole_rank_picture"] = []string{model.WhackMoleRankPicture}
// 	values["whack_mole_rank_background"] = []string{model.WhackMoleRankBackground}

// 	// 搖號抽獎自定義
// 	values["draw_numbers_background"] = []string{model.DrawNumbersBackground}
// 	values["draw_numbers_title"] = []string{model.DrawNumbersTitle}
// 	values["draw_numbers_gift_inside_picture"] = []string{model.DrawNumbersGiftInsidePicture}
// 	values["draw_numbers_gift_outside_picture"] = []string{model.DrawNumbersGiftOutsidePicture}
// 	values["draw_numbers_prize_left_button"] = []string{model.DrawNumbersPrizeLeftButton}
// 	values["draw_numbers_prize_right_button"] = []string{model.DrawNumbersPrizeRightButton}
// 	values["draw_numbers_prize_leftright_button"] = []string{model.DrawNumbersPrizeLeftrightButton}
// 	values["draw_numbers_addpeople_no_button"] = []string{model.DrawNumbersAddpeopleNoNutton}
// 	values["draw_numbers_addpeople_yes_button"] = []string{model.DrawNumbersAddpeopleYesButton}
// 	values["draw_numbers_people_background"] = []string{model.DrawNumbersPeopleBackground}
// 	values["draw_numbers_add_people"] = []string{model.DrawNumbersAddPeople}
// 	values["draw_numbers_reduce_people"] = []string{model.DrawNumbersReducePeople}
// 	values["draw_numbers_winner_background"] = []string{model.DrawNumbersWinnerBackground}
// 	values["draw_numbers_blackground"] = []string{model.DrawNumbersBlackground}
// 	values["draw_numbers_go_button"] = []string{model.DrawNumbersGoButton}
// 	values["draw_numbers_open_winner_button"] = []string{model.DrawNumbersOpenWinnerButton}
// 	values["draw_numbers_close_winner_button"] = []string{model.DrawNumbersCloseWinnerButton}
// 	values["draw_numbers_current_people"] = []string{model.DrawNumbersCurrentPeople}

// 	// 動圖
// 	values["draw_numbers_gacha_machine"] = []string{model.DrawNumbersGachaMachine}
// 	values["draw_numbers_hood"] = []string{model.DrawNumbersHood}
// 	values["draw_numbers_body"] = []string{model.DrawNumbersBody}
// 	values["draw_numbers_gacha"] = []string{model.DrawNumbersGacha}

// 	// 鑑定師自定義
// 	values["monopoly_screen_again_button"] = []string{model.MonopolyScreenAgainButton}                             // 主持端再來一輪按鈕
// 	values["monopoly_screen_top6_title"] = []string{model.MonopolyScreenTop6Title}                                 // 主持端前六名標題
// 	values["monopoly_screen_end_button"] = []string{model.MonopolyScreenEndButton}                                 // 主持端結束遊戲按鈕
// 	values["monopoly_screen_start_button"] = []string{model.MonopolyScreenStartButton}                             // 主持端開始遊戲按鈕
// 	values["monopoly_screen_gaming_background_png"] = []string{model.MonopolyScreenGamingBackgroundPng}            // 主持端遊戲中背景
// 	values["monopoly_screen_round_countdown"] = []string{model.MonopolyScreenRoundCountdown}                       // 主持端遊戲中輪次倒數
// 	values["monopoly_screen_winner_list"] = []string{model.MonopolyScreenWinnerList}                               // 主持端遊戲和結算中獎列表
// 	values["monopoly_screen_rank_border"] = []string{model.MonopolyScreenRankBorder}                               // 主持端遊戲和結算名次框
// 	values["monopoly_player_carton"] = []string{model.MonopolyPlayerCarton}                                        // 玩家端遊戲中下滑紙箱
// 	values["monopoly_player_any_start_text"] = []string{model.MonopolyPlayerAnyStartText}                          // 玩家端按任意處開始文字
// 	values["monopoly_player_scoreboard"] = []string{model.MonopolyPlayerScoreboard}                                // 玩家端計分和時間板
// 	values["monopoly_player_wait_start_text"] = []string{model.MonopolyPlayerWaitStartText}                        // 玩家端計分和時間板
// 	values["monopoly_player_transparent_background"] = []string{model.MonopolyPlayerTransparentBackground}         // 玩家端開始場景半透明黑底
// 	values["monopoly_player_pile_objects"] = []string{model.MonopolyPlayerPileObjects}                             // 玩家端滑動物件堆
// 	values["monopoly_player_gaming_background"] = []string{model.MonopolyPlayerGamingBackground}                   // 玩家端遊戲中背景
// 	values["monopoly_add_points"] = []string{model.MonopolyAddPoints}                                              // 上滑小標示
// 	values["monopoly_deduct_points"] = []string{model.MonopolyDeductPoints}                                        // 下滑小標示
// 	values["monopoly_player_background_dynamic"] = []string{model.MonopolyPlayerBackgroundDynamic}                 // 玩家端背景
// 	values["monopoly_player_answer_effect"] = []string{model.MonopolyPlayerAnswerEffect}                           // 玩家端遊戲中答對或錯特效
// 	values["monopoly_background_and_gold"] = []string{model.MonopolyBackgroundAndGold}                             // 主持端背景和玩家端金銅條
// 	values["monopoly_screen_redpack_seal"] = []string{model.MonopolyScreenRedpackSeal}                             // 主持端紅包袋封口
// 	values["monopoly_screen_again_button_background"] = []string{model.MonopolyScreenAgainButtonBackground}        // 主持端再來一輪按鈕底
// 	values["monopoly_screen_end_info_skin"] = []string{model.MonopolyScreenEndInfoSkin}                            // 主持端結算3名後玩家資訊木框
// 	values["monopoly_screen_end_npc"] = []string{model.MonopolyScreenEndNpc}                                       // 主持端結算吉祥物
// 	values["monopoly_screen_top_stair"] = []string{model.MonopolyScreenTopStair}                                   // 主持端結算前三名台階
// 	values["monopoly_screen_top_info_skin"] = []string{model.MonopolyScreenTopInfoSkin}                            // 主持端結算前三名資訊框
// 	values["monopoly_screen_top_avatar_skin"] = []string{model.MonopolyScreenTopAvatarSkin}                        // 主持端結算前三名頭像框
// 	values["monopoly_screen_end_background"] = []string{model.MonopolyScreenEndBackground}                         // 主持端結算背景
// 	values["monopoly_screen_start_npc_dialog"] = []string{model.MonopolyScreenStartNpcDialog}                      // 主持端開始畫面人物對話框
// 	values["monopoly_screen_leaderboard"] = []string{model.MonopolyScreenLeaderboard}                              // 主持端遊戲中排行榜
// 	values["monopoly_screen_round_background"] = []string{model.MonopolyScreenRoundBackground}                     // 主持端遊戲中輪次底
// 	values["monopoly_screen_start_end_button_background"] = []string{model.MonopolyScreenStartEndButtonBackground} // 主持端開始和結束按鈕底
// 	values["monopoly_screen_start_background"] = []string{model.MonopolyScreenStartBackground}                     // 主持端開始背景
// 	values["monopoly_screen_start_right_top_decoration"] = []string{model.MonopolyScreenStartRightTopDecoration}   // 主持端開始右上裝飾
// 	values["monopoly_player_tip_arrow"] = []string{model.MonopolyPlayerTipArrow}                                   // 玩家端遊戲中提示箭頭
// 	values["monopoly_player_npc_dialog"] = []string{model.MonopolyPlayerNpcDialog}                                 // 玩家端人物對話框
// 	values["monopoly_player_join_button_background"] = []string{model.MonopolyPlayerJoinButtonBackground}          // 玩家端加入遊戲按鈕底
// 	values["monopoly_player_join_background"] = []string{model.MonopolyPlayerJoinBackground}                       // 玩家端加入遊戲背景
// 	values["monopoly_player_redpack_space"] = []string{model.MonopolyPlayerRedpackSpace}                           // 玩家端紅包袋白底
// 	values["monopoly_player_redpack_seal"] = []string{model.MonopolyPlayerRedpackSeal}                             // 玩家端紅包袋封口
// 	values["monopoly_player_redpack_background"] = []string{model.MonopolyPlayerRedpackBackground}                 // 玩家端紅包袋背景
// 	values["monopoly_player_money_piles"] = []string{model.MonopolyPlayerMoneyPiles}                               // 玩家端鈔票堆
// 	values["monopoly_player_background"] = []string{model.MonopolyPlayerBackground}                                // 玩家端遊戲背景
// 	values["monopoly_player_title"] = []string{model.MonopolyPlayerTitle}                                          // 玩家端遊戲標題
// 	values["monopoly_npc"] = []string{model.MonopolyNpc}                                                           // 代表人物
// 	values["monopoly_button"] = []string{model.MonopolyButton}                                                     // 按鈕
// 	values["monopoly_screen_top_light"] = []string{model.MonopolyScreenTopLight}                                   // 主持端前三名發亮
// 	values["monopoly_screen_end_revolving_light"] = []string{model.MonopolyScreenEndRevolvingLight}                // 主持端結算背景旋轉燈
// 	values["monopoly_screen_end_ribbon"] = []string{model.MonopolyScreenEndRibbon}                                 // 主持端結算彩帶
// 	values["monopoly_player_gaming_redpack"] = []string{model.MonopolyPlayerGamingRedpack}                         // 玩家端遊戲中紅包
// 	values["monopoly_screen_gaming_redpack"] = []string{model.MonopolyScreenGamingRedpack}                         // 主持端遊戲中紅包和玩家端紅包
// 	values["monopoly_screen_top_after_player_info"] = []string{model.MonopolyScreenTopAfterPlayerInfo}             // 主持端結算3名後玩家資訊框
// 	values["monopoly_screen_top_front_player_info"] = []string{model.MonopolyScreenTopFrontPlayerInfo}             // 主持端結算前三名資訊框
// 	values["monopoly_screen_rank"] = []string{model.MonopolyScreenRank}                                            // 主持端遊戲中排行榜
// 	values["monopoly_screen_npc_dialog"] = []string{model.MonopolyScreenNpcDialog}                                 // 主持端對話框
// 	values["monopoly_screen_left_bottom_decoration"] = []string{model.MonopolyScreenLeftBottomDecoration}          // 主持端遊戲中裝飾小物件左下
// 	values["monopoly_player_basket_background"] = []string{model.MonopolyPlayerBasketBackground}                   // 玩家端竹籃背景
// 	values["monopoly_player_gaming_carrots"] = []string{model.MonopolyPlayerGamingCarrots}                         // 玩家端遊戲中紅蘿蔔堆
// 	values["monopoly_button_background"] = []string{model.MonopolyButtonBackground}                                // 按鈕背景
// 	values["monopoly_screen_end_background_dynamic"] = []string{model.MonopolyScreenEndBackgroundDynamic}          // 主持端遊戲中和結算背景
// 	values["monopoly_screen_start_background_dynamic"] = []string{model.MonopolyScreenStartBackgroundDynamic}      // 主持端開始背景
// 	values["monopoly_player_gaming_background_dynamic"] = []string{model.MonopolyPlayerGamingBackgroundDynamic}    // 玩家端遊戲背景
// 	values["monopoly_picking_carrots_and_carrots"] = []string{model.MonopolyPickingCarrotsAndCarrots}              // 主持端遊戲中採蘿蔔和玩家端蘿蔔
// 	values["monopoly_player_top_info"] = []string{model.MonopolyPlayerTopInfo}                                     // 玩家端上方資訊
// 	values["monopoly_player_search_prize_background"] = []string{model.MonopolyPlayerSearchPrizeBackground}        // 玩家端查看獎品背景
// 	values["monopoly_player_food_waste_bin"] = []string{model.MonopolyPlayerFoodWasteBin}                          // 玩家端遊戲中廚餘桶
// 	values["monopoly_screen_end_dynamic"] = []string{model.MonopolyScreenEndDynamic}                               // 主持端結算動圖
// 	values["monopoly_screen_timer"] = []string{model.MonopolyScreenTimer}                                          // 主持端遊戲中計時器
// 	values["monopoly_player_start_gaming_eyecatch"] = []string{model.MonopolyPlayerStartGamingEyecatch}            // 玩家端開始和遊戲過場
// 	values["monopoly_gaming_dynamic_and_fish"] = []string{model.MonopolyGamingDynamicAndFish}                      // 主持端遊戲中動圖和玩家端魚
// 	values["monopoly_screen_gaming_background_jpg"] = []string{model.MonopolyScreenGamingBackgroundJpg}            // 主持端遊戲中背景

// 	// 快問快答自定義
// 	values["qa_mascot"] = []string{model.QAMascot}
// 	values["qa_host_start_background"] = []string{model.QAHostStartBackground}
// 	values["qa_host_game_background"] = []string{model.QAHostGameBackground}
// 	values["qa_host_end_background"] = []string{model.QAhostEndBackground}
// 	values["qa_game_top_1"] = []string{model.QAGameTop1}
// 	values["qa_game_top_2"] = []string{model.QAGameTop2}
// 	values["qa_game_top_3"] = []string{model.QAGameTop3}
// 	values["qa_game_top_4"] = []string{model.QAGameTop4}
// 	values["qa_game_top_5"] = []string{model.QAGameTop5}
// 	values["qa_end_top_1"] = []string{model.QAEndTop1}
// 	values["qa_end_top_2"] = []string{model.QAEndTop2}
// 	values["qa_end_top_3"] = []string{model.QAEndTop3}
// 	values["qa_end_top"] = []string{model.QAEndTop}
// 	values["qa_host_start_game_button"] = []string{model.QAHostStartGameButton}
// 	values["qa_host_pause_countdown_button"] = []string{model.QAHostPauseCountdownButton}
// 	values["qa_host_continue_countdown_button"] = []string{model.QAHostContinueCountdownButton}
// 	values["qa_host_start_answer_button"] = []string{model.QAHostStartAnswerButton}
// 	values["qa_host_see_answer_button"] = []string{model.QAHostSeeAnswerButton}
// 	values["qa_host_next_question_button"] = []string{model.QAHostNextQuestionButton}
// 	values["qa_host_end_game_button"] = []string{model.QAHostEndGameButton}
// 	values["qa_host_again_game_button"] = []string{model.QAHostAgainGameButton}
// 	values["qa_player_start_background"] = []string{model.QAPlayerStartBackground}
// 	values["qa_player_game_background"] = []string{model.QAPlayerGameBackground}
// 	values["qa_player_join_game_button"] = []string{model.QAPlayerJoinGameButton}
// 	values["qa_player_select_answer_button"] = []string{model.QAPlayerSelectAnswerButton}
// 	values["qa_player_confirm_answer_button"] = []string{model.QAPlayerConfirmAnswerButton}
// 	values["qa_player_confirm_status_button"] = []string{model.QAPlayerConfirmStatusButton}

// 	// 搖紅包自定義
// 	values["redpack_screen_again_button"] = []string{model.RedpackScreenAgainButton}
// 	values["redpack_screen_background"] = []string{model.RedpackScreenBackground}
// 	values["redpack_screen_end_button"] = []string{model.RedpackScreenEndButton}
// 	values["redpack_screen_prize_list"] = []string{model.RedpackScreenPrizeList}
// 	values["redpack_screen_prize_redpack"] = []string{model.RedpackScreenPrizeRedpack}
// 	values["redpack_screen_start_button"] = []string{model.RedpackScreenStartButton}
// 	values["redpack_screen_title"] = []string{model.RedpackScreenTitle}
// 	values["redpack_screen_gaming_list"] = []string{model.RedpackScreenGamingList}
// 	values["redpack_screen_gaming_list_background"] = []string{model.RedpackScreenGamingListBackground}
// 	values["redpack_screen_ema"] = []string{model.RedpackScreenEma}
// 	values["redpack_screen_new_list"] = []string{model.RedpackScreenNewList}
// 	values["redpack_screen_lantern1"] = []string{model.RedpackScreenLantern1}
// 	values["redpack_screen_lantern2"] = []string{model.RedpackScreenLantern2}
// 	values["redpack_player_join_button"] = []string{model.RedpackPlayerJoinButton}
// 	values["redpack_player_search_prize_background"] = []string{model.RedpackPlayerSearchPrizeBackground}
// 	values["redpack_player_background"] = []string{model.RedpackPlayerBackground}
// 	values["redpack_player_title"] = []string{model.RedpackPlayerTitle}
// 	values["redpack_player_lantern"] = []string{model.RedpackPlayerLantern}
// 	// 動圖
// 	values["redpack_screen_lucky_bag"] = []string{model.RedpackScreenLuckyBag}
// 	values["redpack_screen_money_piles"] = []string{model.RedpackScreenMoneyPiles}
// 	values["redpack_player_shake"] = []string{model.RedpackPlayerShake}
// 	values["redpack_player_lucky_bag"] = []string{model.RedpackPlayerLuckyBag}
// 	values["redpack_player_money_piles"] = []string{model.RedpackPlayerMoneyPiles}
// 	values["redpack_screen_background_dynamic"] = []string{model.RedpackScreenBackgroundDynamic}
// 	values["redpack_player_background_dynamic"] = []string{model.RedpackPlayerBackgroundDynamic}
// 	values["redpack_player_ready"] = []string{model.RedpackPlayerReady}
// 	// 音樂
// 	values["redpack_bgm_start"] = []string{model.RedpackBgmStart}
// 	values["redpack_bgm_gaming"] = []string{model.RedpackBgmGaming}
// 	values["redpack_bgm_end"] = []string{model.RedpackBgmEnd}

// 	// 套紅包自定義
// 	values["ropepack_screen_prize_list"] = []string{model.RopepackScreenPrizeList}
// 	values["ropepack_screen_again_button"] = []string{model.RopepackScreenAgainButton}
// 	values["ropepack_screen_background"] = []string{model.RopepackScreenBackground}
// 	values["ropepack_screen_decoration"] = []string{model.RopepackScreenDecoration}
// 	values["ropepack_screen_end_button"] = []string{model.RopepackScreenEndButton}
// 	values["ropepack_screen_end_prize_list"] = []string{model.RopepackScreenEndPrizeList}
// 	values["ropepack_screen_prize_redpack"] = []string{model.RopepackScreenPrizeRedpack}
// 	values["ropepack_screen_start_logo"] = []string{model.RopepackScreenStartLogo}
// 	values["ropepack_screen_start_button"] = []string{model.RopepackScreenStartButton}
// 	values["ropepack_screen_prize_skin_red"] = []string{model.RopepackScreenPrizeSkinRed}
// 	values["ropepack_screen_prize_skin_green"] = []string{model.RopepackScreenPrizeSkinGreen}
// 	values["ropepack_player_join_logo"] = []string{model.RopepackPlayerJoinLogo}
// 	values["ropepack_player_join_button"] = []string{model.RopepackPlayerJoinButton}
// 	values["ropepack_player_background"] = []string{model.RopepackPlayerBackground}
// 	values["ropepack_player_ready_redpack1"] = []string{model.RopepackPlayerReadyRedpack1}
// 	values["ropepack_player_ready_redpack2"] = []string{model.RopepackPlayerReadyRedpack2}
// 	values["ropepack_player_ready_background"] = []string{model.RopepackPlayerReadyBackground}
// 	values["ropepack_player_title"] = []string{model.RopepackPlayerTitle}
// 	// 動圖
// 	values["ropepack_screen_background_effect"] = []string{model.RopepackScreenBackgroundEffect}
// 	values["ropepack_player_ropepack_button"] = []string{model.RopepackPlayerRopepackButton}
// 	values["ropepack_player_finger"] = []string{model.RopepackPlayerFinger}
// 	values["ropepack_redpack"] = []string{model.RopepackRedpack}

// 	// 遊戲抽獎自定義
// 	values["lottery_screen_prizer"] = []string{model.LotteryScreenPrizer}
// 	values["lottery_screen_mascot"] = []string{model.LotteryScreenMascot}
// 	values["lottery_screen_background"] = []string{model.LotteryScreenBackground}
// 	values["lottery_screen_prize_notify"] = []string{model.LotteryScreenPrizeNotify}
// 	values["lottery_screen_select_input"] = []string{model.LotteryScreenSelectInput}
// 	values["lottery_screen_close_prize_notify_button"] = []string{model.LotteryScreenClosePrizeNotifyButton}
// 	values["lottery_player_background"] = []string{model.LotteryPlayerBackground}
// 	values["lottery_player_rules"] = []string{model.LotteryPlayerRules}
// 	values["lottery_jiugongge_grid"] = []string{model.LotteryJiugonggeGrid}
// 	values["lottery_jiugongge_start_button"] = []string{model.LotteryJiugonggeStartButton}
// 	values["lottery_turntable_start_button"] = []string{model.LotteryTurntableStartButton}
// 	values["lottery_turntable_roulette"] = []string{model.LotteryTurntableRoulette}
// 	// 動圖
// 	values["lottery_get_prize"] = []string{model.LotteryGetPrize}
// 	values["lottery_jiugongge_border"] = []string{model.LotteryJiugonggeBorder}
// 	values["lottery_jiugongge_title"] = []string{model.LotteryJiugonggeTitle}
// 	values["lottery_turntable_border"] = []string{model.LotteryTurntableBorder}
// 	values["lottery_turntable_title"] = []string{model.LotteryTurntableTitle}

// 	values["qa_1"] = []string{model.QA1}
// 	values["qa_1_options"] = []string{model.QA1Options}
// 	values["qa_1_answer"] = []string{model.QA1Answer}
// 	values["qa_1_score"] = []string{model.QA1Score}

// 	values["qa_2"] = []string{model.QA2}
// 	values["qa_2_options"] = []string{model.QA2Options}
// 	values["qa_2_answer"] = []string{model.QA2Answer}
// 	values["qa_2_score"] = []string{model.QA2Score}

// 	values["qa_3"] = []string{model.QA3}
// 	values["qa_3_options"] = []string{model.QA3Options}
// 	values["qa_3_answer"] = []string{model.QA3Answer}
// 	values["qa_3_score"] = []string{model.QA3Score}

// 	values["qa_4"] = []string{model.QA4}
// 	values["qa_4_options"] = []string{model.QA4Options}
// 	values["qa_4_answer"] = []string{model.QA4Answer}
// 	values["qa_4_score"] = []string{model.QA4Score}

// 	values["qa_5"] = []string{model.QA5}
// 	values["qa_5_options"] = []string{model.QA5Options}
// 	values["qa_5_answer"] = []string{model.QA5Answer}
// 	values["qa_5_score"] = []string{model.QA5Score}

// 	values["qa_6"] = []string{model.QA6}
// 	values["qa_6_options"] = []string{model.QA6Options}
// 	values["qa_6_answer"] = []string{model.QA6Answer}
// 	values["qa_6_score"] = []string{model.QA6Score}

// 	values["qa_7"] = []string{model.QA7}
// 	values["qa_7_options"] = []string{model.QA7Options}
// 	values["qa_7_answer"] = []string{model.QA7Answer}
// 	values["qa_7_score"] = []string{model.QA7Score}

// 	values["qa_8"] = []string{model.QA8}
// 	values["qa_8_options"] = []string{model.QA8Options}
// 	values["qa_8_answer"] = []string{model.QA8Answer}
// 	values["qa_8_score"] = []string{model.QA8Score}

// 	values["qa_9"] = []string{model.QA9}
// 	values["qa_9_options"] = []string{model.QA9Options}
// 	values["qa_9_answer"] = []string{model.QA9Answer}
// 	values["qa_9_score"] = []string{model.QA9Score}

// 	values["qa_10"] = []string{model.QA10}
// 	values["qa_10_options"] = []string{model.QA10Options}
// 	values["qa_10_answer"] = []string{model.QA10Answer}
// 	values["qa_10_score"] = []string{model.QA10Score}

// 	values["qa_11"] = []string{model.QA11}
// 	values["qa_11_options"] = []string{model.QA11Options}
// 	values["qa_11_answer"] = []string{model.QA11Answer}
// 	values["qa_11_score"] = []string{model.QA11Score}

// 	values["qa_12"] = []string{model.QA12}
// 	values["qa_12_options"] = []string{model.QA12Options}
// 	values["qa_12_answer"] = []string{model.QA12Answer}
// 	values["qa_12_score"] = []string{model.QA12Score}

// 	values["qa_13"] = []string{model.QA13}
// 	values["qa_13_options"] = []string{model.QA13Options}
// 	values["qa_13_answer"] = []string{model.QA13Answer}
// 	values["qa_13_score"] = []string{model.QA13Score}

// 	values["qa_14"] = []string{model.QA14}
// 	values["qa_14_options"] = []string{model.QA14Options}
// 	values["qa_14_answer"] = []string{model.QA14Answer}
// 	values["qa_14_score"] = []string{model.QA14Score}

// 	values["qa_15"] = []string{model.QA15}
// 	values["qa_15_options"] = []string{model.QA15Options}
// 	values["qa_15_answer"] = []string{model.QA15Answer}
// 	values["qa_15_score"] = []string{model.QA15Score}

// 	values["qa_16"] = []string{model.QA16}
// 	values["qa_16_options"] = []string{model.QA16Options}
// 	values["qa_16_answer"] = []string{model.QA16Answer}
// 	values["qa_16_score"] = []string{model.QA16Score}

// 	values["qa_17"] = []string{model.QA17}
// 	values["qa_17_options"] = []string{model.QA17Options}
// 	values["qa_17_answer"] = []string{model.QA17Answer}
// 	values["qa_17_score"] = []string{model.QA17Score}

// 	values["qa_18"] = []string{model.QA18}
// 	values["qa_18_options"] = []string{model.QA18Options}
// 	values["qa_18_answer"] = []string{model.QA18Answer}
// 	values["qa_18_score"] = []string{model.QA18Score}

// 	values["qa_19"] = []string{model.QA19}
// 	values["qa_19_options"] = []string{model.QA19Options}
// 	values["qa_19_answer"] = []string{model.QA19Answer}
// 	values["qa_19_score"] = []string{model.QA19Score}

// 	values["qa_20"] = []string{model.QA20}
// 	values["qa_20_options"] = []string{model.QA20Options}
// 	values["qa_20_answer"] = []string{model.QA20Answer}
// 	values["qa_20_score"] = []string{model.QA20Score}

// 	values["total_qa"] = []string{model.TotalQA}
// 	values["qa_second"] = []string{model.QASecond}
// 	}
// 	if err != nil {
// 		if strings.Contains(ctx.GetHeader("Accept"), "json") {
// 			response.Error(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		} else {
// 			h.executeErrorHTML(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		}
// 	}
// }

// @@@Summary 編輯遊戲基本設置資料(json)
// @@@Tags Game
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameSettingModel true "Game Setting Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/setting/json [put]
func (h *Handler) PUTJSONGameSetting(ctx *gin.Context) {
}

// @@@Summary 編輯遊戲抽獎遊戲資料(json)
// @@@Tags Lottery
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "Lottery Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/json [put]
func (h *Handler) PUTJSONLottery(ctx *gin.Context) {
}

// @@@Summary 編輯遊戲抽獎獎品資料(json)
// @@@Tags Lottery Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "Lottery Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/prize/json [put]
func (h *Handler) PUTJSONLotteryPrize(ctx *gin.Context) {
}

// @@@Summary 編輯搖紅包遊戲資料(json)
// @@@Tags Redpack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "Redpack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/json [put]
func (h *Handler) PUTJSONRedpack(ctx *gin.Context) {
}

// @@@Summary 編輯搖紅包獎品資料(json)
// @@@Tags Redpack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "Redpack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/prize/json [put]
func (h *Handler) PUTJSONRedpackPrize(ctx *gin.Context) {
}

// @@@Summary 編輯套紅包遊戲資料(json)
// @@@Tags Ropepack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "Ropepack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/json [put]
func (h *Handler) PUTJSONRopepack(ctx *gin.Context) {
}

// @@@Summary 編輯套紅包獎品資料(json)
// @@@Tags Ropepack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "Ropepack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/prize/json [put]
func (h *Handler) PUTJSONRopepackPrize(ctx *gin.Context) {
}

// @@@Summary 編輯敲敲樂遊戲資料(json)
// @@@Tags Whack_Mole
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "WhackMole Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/json [put]
func (h *Handler) PUTJSONWhackMole(ctx *gin.Context) {
}

// @@@Summary 編輯敲敲樂獎品資料(json)
// @@@Tags Whack_Mole Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "WhackMole Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/prize/json [put]
func (h *Handler) PUTJSONWhackMolePrize(ctx *gin.Context) {
}

// @@@Summary 編輯搖號抽獎遊戲資料(json)
// @@@Tags Draw_Numbers
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "Draw_Numbers Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/json [put]
func (h *Handler) PUTJSONDrawNumber(ctx *gin.Context) {
}

// @@@Summary 編輯搖號抽獎獎品資料(json)
// @@@Tags Draw_Numbers Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "Draw_Numbers Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/prize/json [put]
func (h *Handler) PUTJSONDrawNumbersPrize(ctx *gin.Context) {
}

// @@@Summary 編輯鑑定師遊戲資料(json)
// @@@Tags Monopoly
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "Monopoly Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/json [put]
func (h *Handler) PUTJSONMonopoly(ctx *gin.Context) {
}

// @@@Summary 編輯鑑定師獎品資料(json)
// @@@Tags Monopoly Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "Monopoly Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/prize/json [put]
func (h *Handler) PUTJSONMonopolyPrize(ctx *gin.Context) {
}

// @@@Summary 編輯快問快答遊戲資料(json)
// @@@Tags QA
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditGameModel true "QA Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/json [put]
func (h *Handler) PUTJSONQA(ctx *gin.Context) {
}

// @@@Summary 編輯快問快答獎品資料(json)
// @@@Tags QA Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeModel true "QA Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/prize/json [put]
func (h *Handler) PUTJSONQAPrize(ctx *gin.Context) {
}

// @@@Summary 編輯黑名單人員資料(json)
// @@@Tags Black Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditBlackStaffModel true "Black Staff Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/black/json [put]
func (h *Handler) PUTJSONBlack(ctx *gin.Context) {
}

// @@@Summary 編輯中獎人員資料(json)
// @@@Tags Winning Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.EditPrizeStaffModel true "Winning Staff Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /staffmanage/winning/json [put]
func (h *Handler) PUTJSONWinning(ctx *gin.Context) {
}
