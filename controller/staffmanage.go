package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowGuestWinning 個人中獎紀錄頁面(手機用戶端) GET API
func (h *Handler) ShowGuestWinning(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		activityID = ctx.Query("activity_id")
		user       = h.GetLoginUser(ctx.Request, "chatroom_session")
		htmlTmpl   = "./hilives/hilive/views/chatroom/style/default/my_prize_history.html"
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 取得用戶中獎紀錄
	// prizes, err := models.DefaultPrizeStaffModel().SetDbConn(h.dbConn).
	// 	FindUserWinningRecords(user.UserID, "activity_staff_prize.activity_id", activityID)
	// if err != nil {
	// 	h.executeErrorHTML(ctx, err.Error())
	// 	return
	// }

	// fmt.Println("個人用戶中獎資訊: ", prizes)
	// fmt.Println(utils.JSON(prizes))

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User: user,
		Route: route{
			Chatroom: fmt.Sprintf(config.GUEST_URL, activityID),
			PUT:      fmt.Sprintf(config.STAFFMANAGE_API_URL_FORM, "winning"),
		},
		// WinningSaffsJSON: utils.JSON(prizes),
	})
}

// ShowStaffManage 人員管理頁面，包含參加遊戲人員、中獎人員、黑名單人員、遊戲紀錄(平台) GET API
func (h *Handler) ShowStaffManage(ctx *gin.Context) {
	var (
		host                       = ctx.Request.Host
		path                       = ctx.Request.URL.Path
		staffmanage                = ctx.Param("__staffmanage")
		activityID                 = ctx.Query("activity_id")
		user                       = h.GetLoginUser(ctx.Request, "hilive_session")
		canAdd, canEdit, canDelete bool
		htmlTmpl                   string
		// err                        error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if staffmanage == "attend" {
		htmlTmpl = "./hilives/hilive/views/cms/person_game_list.html"
	} else if staffmanage == "winning" {
		htmlTmpl = "./hilives/hilive/views/cms/person_prize_list.html"
	} else if staffmanage == "black" {
		htmlTmpl = "./hilives/hilive/views/cms/person_black_list.html"
	} else if staffmanage == "record" { // 遊戲紀錄
		htmlTmpl = "./hilives/hilive/views/cms/person_history_list.html"
	}

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:       user,
		ActivityID: activityID,
		Route: route{
			StaffManage: fmt.Sprintf(config.STAFFMANAGE_URL, staffmanage, activityID),
			POST:        fmt.Sprintf(config.STAFFMANAGE_API_URL_FORM, staffmanage),
			PUT:         fmt.Sprintf(config.STAFFMANAGE_API_URL_FORM, staffmanage),
			DELETE:      fmt.Sprintf(config.STAFFMANAGE_API_URL_FORM, staffmanage),
			Export:      fmt.Sprintf(config.STAFFMANAGE_API_URL_EXPORT, staffmanage),
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// ShowGuestGameWinningStaff 遊戲互動裡的所有人員中獎紀錄頁面(手機用戶端) GET API
func (h *Handler) ShowGuestGameWinningStaff(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		gameID     = ctx.Query("game_id")
		panel, _   = h.GetTable(ctx, "winning")
		panelInfo  table.PanelInfo
		params     = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"round", "prize_id"}, []string{"asc", "asc"})
		// staffInfoList types.InfoList
		rounds   []string
		htmlTmpl string
		err      error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" || gameID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID")
		return
	}
	if prefix == "redpack" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/redpack_prize.html"
	} else if prefix == "ropepack" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/ropepack_prize.html"
	} else if prefix == "whack_mole" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/whack_mole_prize.html"
	} else if prefix == "lottery" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/lottery_prize.html"
	} else if prefix == "monopoly" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/monopoly_prize.html"
	} else if prefix == "QA" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/QA_prize.html"
	} else if prefix == "draw_numbers" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/draw_numbers_prize.html"
	} else if prefix == "tugofwar" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/tugofwar_prize.html"
	} else if prefix == "bingo" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/bingo_prize.html"
	} else if prefix == "3DGachaMachine" {
		htmlTmpl = "./hilives/hilive/views/chatroom/style/default/gacha_machine_prize.html"
	}

	// 所有中獎人員資料
	// 中獎判斷
	// params.SetField("activity_prize_join_prize_type", "first", "second", "third", "general")
	// params.SetField("activity_prize_join_prize_method", "site", "mail")
	if panelInfo, err = panel.GetData(params, h.services); err != nil {
		h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
		return
	}

	// fmt.Println("staffInfoList: ", len(staffInfoList))
	// fmt.Println("panelInfo.InfoList: ", panelInfo.InfoList)
	h.executeHTML(ctx, htmlTmpl, executeParam{
		Route: route{},

		StaffManage: StaffManage{
			StaffList: panelInfo.InfoList,
			Rounds:    rounds,
		},
	})
}

// 判斷路徑
// var path string
// if game == "activity" {
// 	path = fmt.Sprintf("/admin/staffmanage/%s/%s/%s", staffmanage, game, staffmanage)
// } else if game == "message" {
// 	path = fmt.Sprintf("/admin/staffmanage/%s/%s/%s", staffmanage, game, staffmanage)
// } else if game == "question" {
// 	path = fmt.Sprintf("/admin/staffmanage/%s/%s/%s", staffmanage, game, staffmanage)
// } else {
// 	path = fmt.Sprintf("/admin/staffmanage/%s/%s", staffmanage, game)
// }

// 判斷是否有新增編輯刪除權限
// for _, permisssion := range user.Permissions {
// 	for _, httppath := range permisssion.HTTPPath {
// 		if httppath == path {
// 			if strings.Contains(permisssion.HTTPMethod, "POST") {
// 				canAdd = true
// 			}
// 			if strings.Contains(permisssion.HTTPMethod, "PUT") {
// 				canEdit = true
// 			}
// 			if strings.Contains(permisssion.HTTPMethod, "DELETE") {
// 				canDelete = true
// 			}
// 		}
// 	}
// }
// CanAdd:    canAdd,
// CanEdit:   canEdit,
// CanDelete: canDelete,

// PanelInfo: panelInfo,
// Chatroom: fmt.Sprintf(config.GUEST_CHATROOM_URL, activityID),
// Info:     fmt.Sprintf("%s/?activity_id=%s", config.GUEST_INFO_URL, activityID),
// GameInfo: fmt.Sprintf("%s?activity_id=%s", config.GUEST_GAME_URL, activityID),

// fmt.Println("panelInfo.InfoList: ", panelInfo.InfoList)

// for i := 0; i < len(panelInfo.InfoList); i++ {
// 	info := panelInfo.InfoList[i]
// 	// 過濾輪次參數
// 	if i == 0 {
// 		rounds = append(rounds, info["round"].Value)
// 	} else {
// 		if info["round"].Value != panelInfo.InfoList[i-1]["round"].Value {
// 			rounds = append(rounds, info["round"].Value)
// 		}
// 	}

// 	// if info["white"].Value == "no" { // 一般名單
// 		staffInfoList = append(staffInfoList, info)
// 	// }
// }

// if (staffmanage != "attend" && staffmanage != "winning") ||
// 	(game != "lottery" && game != "redpack" && game != "ropepack" &&
// 		game != "whack_mole" && game != "draw_numbers" && game != "monopoly" && game != "QA") {
// 	h.executeErrorHTML(ctx, "錯誤: 網址發生問題，請輸入有效的網址")
// 	return
// }

// if game == "lottery" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "redpack" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "ropepack" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "whack_mole" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "monopoly" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "QA" {
// 	htmlTmpl = staff.GameSceneList
// } else if game == "draw_numbers" {
// 	htmlTmpl = staff.GameSceneList
// }
// if game == "lottery" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "redpack" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "ropepack" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "whack_mole" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "monopoly" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "QA" {
// 	htmlTmpl = staff.PrizeSceneList
// } else if game == "draw_numbers" {
// 	htmlTmpl = staff.PrizeSceneList
// }
// if game == "lottery" {
// 	htmlTmpl = ""
// } else if game == "redpack" {
// 	htmlTmpl = ""
// } else if game == "ropepack" {
// 	htmlTmpl = ""
// } else if game == "whack_mole" {
// 	htmlTmpl = ""
// } else if game == "monopoly" {
// 	htmlTmpl = ""
// } else if game == "QA" {
// 	htmlTmpl = ""
// } else if game == "draw_numbers" {
// 	htmlTmpl = ""
// }

// if (staffmanage != "attend" && staffmanage != "winning") ||
// 	(game != "lottery" && game != "redpack" && game != "ropepack" &&
// 		game != "whack_mole" && game != "monopoly" && game != "QA" &&
// 		game != "draw_numbers") {
// 	h.executeErrorHTML(ctx, "錯誤: 網址發生問題，請輸入有效的網址")
// 	return
// }

// model, err                   = models.DefaultActivityModel().SetDbConn(h.dbConn).Find("activity_id", activity)
// fields, rounds            []string
// if game != "draw_numbers" && (activityID == "" || gameID == "") {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動、遊戲資訊，請輸入有效的活動及遊戲ID")
// 	return
// } else if game == "draw_numbers" && activityID == "" {
// 	h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的及遊戲ID")
// 	return
// }

// 新增資料的表單
// newFormInfo table.FormInfo
// model, err = models.DefaultActivityModel().SetDbConn(h.dbConn).Find("activity_id", activityID)
// newFormInfo := panel.GetNewFormInfo(h.services, params, []string{"activity_id", "user_id"})

// addGameStaff call add game staff JSON API
// func addGameStaff(model models.NewGameStaffModel) error {
// 	var (
// 		url    = "https://api.hilives.net/v1/staffmanage/attend/json"
// 		client = &http.Client{}
// 		result interface{}
// 	)

// 	jsonModel, err := json.Marshal(model)
// 	if err != nil {
// 		return err
// 	}

// 	// 加入遊戲人員資料
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonModel))
// 	if err != nil {
// 		return errors.New("錯誤: 呼叫新增遊戲人員API發生問題")
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return errors.New("錯誤: 讀取內容發生錯誤")
// 	}

// 	json.Unmarshal(body, &result)
// 	if result.(m
