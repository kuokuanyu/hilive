package table

import (
	"encoding/json"
	"errors"

	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"

	"github.com/gin-gonic/gin"
)

// GetGameSettingPanel 基本遊戲設置
func (s *SystemTable) GetGameSettingPanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	info := table.GetInfo()

	info.SetTable(config.ACTIVITY_GAME_SETTING_TABLE)

	formList := table.GetForm()

	formList.SetTable(config.ACTIVITY_GAME_SETTING_TABLE).SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("錯誤: 活動ID發生問題，請輸入有效的活動ID")
		}

		// 將map[string][]string格是資料轉換為map[string]string
		flattened := utils.FlattenForm(values)

		// 將 map 轉成 JSON
		jsonBytes, err := json.Marshal(flattened)
		if err != nil {
			return err
		}

		// 轉成 struct
		var model models.EditGameSettingModel
		if err := json.Unmarshal(jsonBytes, &model); err != nil {
			return err
		}

		if err := models.DefaultGameSettingModel().
			SetConn(s.dbConn, s.redisConn, s.mongoConn).
			Update(true, model); err != nil {
			return err
		}
		return nil
	})
	return
}

// models.EditGameSettingModel{
// 	ActivityID:            values.Get("activity_id"),
// 	LotteryGameAllow:      values.Get("lottery_game_allow"),
// 	RedpackGameAllow:      values.Get("redpack_game_allow"),
// 	RopepackGameAllow:     values.Get("ropepack_game_allow"),
// 	WhackMoleGameAllow:    values.Get("whack_mole_game_allow"),
// 	MonopolyGameAllow:     values.Get("monopoly_game_allow"),
// 	QAGameAllow:           values.Get("qa_game_allow"),
// 	DrawNumbersGameAllow:  values.Get("draw_numbers_game_allow"),
// 	TugofwarGameAllow:     values.Get("tugofwar_game_allow"),
// 	BingoGameAllow:        values.Get("bingo_game_allow"),
// 	GachaMachineGameAllow: values.Get("3d_gacha_machine_game_allow"),
// 	VoteGameAllow:         values.Get("vote_game_allow"),
// 	AllGameAllow:          values.Get("all_game_allow"),
// }

// applyPictureToNewGameModel 將map中所有圖片資料自動寫入struct中
// func applyPictureToNewGameModel(model *models.NewGameModel, picture map[string]string) {
// 	v := reflect.ValueOf(model).Elem() // 透過 reflect 取得結構體的實體（不是指標）
// 	for field, val := range picture {    // 遍歷 media 的每個欄位與值
// 		f := v.FieldByName(field)      // 根據欄位名稱找出 struct 對應的欄位
// 		if f.IsValid() && f.CanSet() { // 確保這個欄位是存在且可以設值的
// 			f.SetString(val) // 設定這個欄位為 val（字串型）
// 		}
// 	}
// }

// @Summary 編輯遊戲基本設置資料(form-data)
// @Tags Game
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @param lottery_game_allow formData string false "遊戲抽獎遊戲是否允許重複中獎" Enums(open, close)
// @param redpack_game_allow formData string false "搖紅包遊戲是否允許重複中獎" Enums(open, close)
// @param ropepack_game_allow formData string false "套紅包遊戲是否允許重複中獎" Enums(open, close)
// @param whack_mole_game_allow formData string false "敲敲樂遊戲是否允許重複中獎" Enums(open, close)
// @param monopoly_game_allow formData string false "鑑定師遊戲是否允許重複中獎" Enums(open, close)
// @param qa_game_allow formData string false "快問快答遊戲是否允許重複中獎" Enums(open, close)
// @param draw_numbers_game_allow formData string false "搖號抽獎遊戲是否允許重複中獎" Enums(open, close)
// @param tugofwar_game_allow formData string false "拔河遊戲是否允許重複中獎" Enums(open, close)
// @param bingo_game_allow formData string false "賓果遊戲是否允許重複中獎" Enums(open, close)
// @param 3d_gacha_machine_game_allow formData string false "扭蛋機遊戲是否允許重複中獎" Enums(open, close)
// @param vote_game_allow formData string false "投票遊戲是否允許重複中獎" Enums(open, close)
// @param all_game_allow formData string false "所有遊戲是否允許重複中獎" Enums(open, close)
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/setting/form [put]
func PUTGameSetting(ctx *gin.Context) {
}

// @Summary 所有遊戲JSON資料(或特定遊戲類型)
// @Tags Game
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @param game query string false "遊戲種類(用,間隔)"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game [get]
func GamesJSON(ctx *gin.Context) {
}

// @Summary 遊戲基本設置JSON資料
// @Tags Game
// @version 1.0
// @Accept  json
// @param activity_id query string true "活動ID"
// @Success 200 {array} response.ResponseWithData
// @Failure 404 {array} response.ResponseBadRequest
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /interact/game/setting [get]
func GameSettingJSON(ctx *gin.Context) {
}

// 原始讀取檔案的題目設置
// a2, _ := file.GetCellValue("Sheet1", "A2")
// if a2 != "" {
// 	total = "1"
// }
// b2, _ := file.GetCellValue("Sheet1", "B2")
// c2, _ := file.GetCellValue("Sheet1", "C2")
// d2, _ := file.GetCellValue("Sheet1", "D2")
// e2, _ := file.GetCellValue("Sheet1", "E2")
// f2, _ := file.GetCellValue("Sheet1", "F2")
// g2, _ := file.GetCellValue("Sheet1", "G2")

// a3, _ := file.GetCellValue("Sheet1", "A3")
// if a3 != "" {
// 	total = "2"
// }
// b3, _ := file.GetCellValue("Sheet1", "B3")
// c3, _ := file.GetCellValue("Sheet1", "C3")
// d3, _ := file.GetCellValue("Sheet1", "D3")
// e3, _ := file.GetCellValue("Sheet1", "E3")
// f3, _ := file.GetCellValue("Sheet1", "F3")
// g3, _ := file.GetCellValue("Sheet1", "G3")

// a4, _ := file.GetCellValue("Sheet1", "A4")
// if a4 != "" {
// 	total = "3"
// }
// b4, _ := file.GetCellValue("Sheet1", "B4")
// c4, _ := file.GetCellValue("Sheet1", "C4")
// d4, _ := file.GetCellValue("Sheet1", "D4")
// e4, _ := file.GetCellValue("Sheet1", "E4")
// f4, _ := file.GetCellValue("Sheet1", "F4")
// g4, _ := file.GetCellValue("Sheet1", "G4")

// a5, _ := file.GetCellValue("Sheet1", "A5")
// if a5 != "" {
// 	total = "4"
// }
// b5, _ := file.GetCellValue("Sheet1", "B5")
// c5, _ := file.GetCellValue("Sheet1", "C5")
// d5, _ := file.GetCellValue("Sheet1", "D5")
// e5, _ := file.GetCellValue("Sheet1", "E5")
// f5, _ := file.GetCellValue("Sheet1", "F5")
// g5, _ := file.GetCellValue("Sheet1", "G5")

// a6, _ := file.GetCellValue("Sheet1", "A6")
// if a6 != "" {
// 	total = "5"
// }
// b6, _ := file.GetCellValue("Sheet1", "B6")
// c6, _ := file.GetCellValue("Sheet1", "C6")
// d6, _ := file.GetCellValue("Sheet1", "D6")
// e6, _ := file.GetCellValue("Sheet1", "E6")
// f6, _ := file.GetCellValue("Sheet1", "F6")
// g6, _ := file.GetCellValue("Sheet1", "G6")

// a7, _ := file.GetCellValue("Sheet1", "A7")
// if a7 != "" {
// 	total = "6"
// }
// b7, _ := file.GetCellValue("Sheet1", "B7")
// c7, _ := file.GetCellValue("Sheet1", "C7")
// d7, _ := file.GetCellValue("Sheet1", "D7")
// e7, _ := file.GetCellValue("Sheet1", "E7")
// f7, _ := file.GetCellValue("Sheet1", "F7")
// g7, _ := file.GetCellValue("Sheet1", "G7")

// a8, _ := file.GetCellValue("Sheet1", "A8")
// if a8 != "" {
// 	total = "7"
// }
// b8, _ := file.GetCellValue("Sheet1", "B8")
// c8, _ := file.GetCellValue("Sheet1", "C8")
// d8, _ := file.GetCellValue("Sheet1", "D8")
// e8, _ := file.GetCellValue("Sheet1", "E8")
// f8, _ := file.GetCellValue("Sheet1", "F8")
// g8, _ := file.GetCellValue("Sheet1", "G8")

// a9, _ := file.GetCellValue("Sheet1", "A9")
// if a9 != "" {
// 	total = "8"
// }
// b9, _ := file.GetCellValue("Sheet1", "B9")
// c9, _ := file.GetCellValue("Sheet1", "C9")
// d9, _ := file.GetCellValue("Sheet1", "D9")
// e9, _ := file.GetCellValue("Sheet1", "E9")
// f9, _ := file.GetCellValue("Sheet1", "F9")
// g9, _ := file.GetCellValue("Sheet1", "G9")

// a10, _ := file.GetCellValue("Sheet1", "A10")
// if a10 != "" {
// 	total = "9"
// }
// b10, _ := file.GetCellValue("Sheet1", "B10")
// c10, _ := file.GetCellValue("Sheet1", "C10")
// d10, _ := file.GetCellValue("Sheet1", "D10")
// e10, _ := file.GetCellValue("Sheet1", "E10")
// f10, _ := file.GetCellValue("Sheet1", "F10")
// g10, _ := file.GetCellValue("Sheet1", "G10")

// a11, _ := file.GetCellValue("Sheet1", "A11")
// if a11 != "" {
// 	total = "10"
// }
// b11, _ := file.GetCellValue("Sheet1", "B11")
// c11, _ := file.GetCellValue("Sheet1", "C11")
// d11, _ := file.GetCellValue("Sheet1", "D11")
// e11, _ := file.GetCellValue("Sheet1", "E11")
// f11, _ := file.GetCellValue("Sheet1", "F11")
// g11, _ := file.GetCellValue("Sheet1", "G11")

// a12, _ := file.GetCellValue("Sheet1", "A12")
// if a12 != "" {
// 	total = "11"
// }
// b12, _ := file.GetCellValue("Sheet1", "B12")
// c12, _ := file.GetCellValue("Sheet1", "C12")
// d12, _ := file.GetCellValue("Sheet1", "D12")
// e12, _ := file.GetCellValue("Sheet1", "E12")
// f12, _ := file.GetCellValue("Sheet1", "F12")
// g12, _ := file.GetCellValue("Sheet1", "G12")

// a13, _ := file.GetCellValue("Sheet1", "A13")
// if a13 != "" {
// 	total = "12"
// }
// b13, _ := file.GetCellValue("Sheet1", "B13")
// c13, _ := file.GetCellValue("Sheet1", "C13")
// d13, _ := file.GetCellValue("Sheet1", "D13")
// e13, _ := file.GetCellValue("Sheet1", "E13")
// f13, _ := file.GetCellValue("Sheet1", "F13")
// g13, _ := file.GetCellValue("Sheet1", "G13")

// a14, _ := file.GetCellValue("Sheet1", "A14")
// if a14 != "" {
// 	total = "13"
// }
// b14, _ := file.GetCellValue("Sheet1", "B14")
// c14, _ := file.GetCellValue("Sheet1", "C14")
// d14, _ := file.GetCellValue("Sheet1", "D14")
// e14, _ := file.GetCellValue("Sheet1", "E14")
// f14, _ := file.GetCellValue("Sheet1", "F14")
// g14, _ := file.GetCellValue("Sheet1", "G14")

// a15, _ := file.GetCellValue("Sheet1", "A15")
// if a15 != "" {
// 	total = "14"
// }
// b15, _ := file.GetCellValue("Sheet1", "B15")
// c15, _ := file.GetCellValue("Sheet1", "C15")
// d15, _ := file.GetCellValue("Sheet1", "D15")
// e15, _ := file.GetCellValue("Sheet1", "E15")
// f15, _ := file.GetCellValue("Sheet1", "F15")
// g15, _ := file.GetCellValue("Sheet1", "G15")

// a16, _ := file.GetCellValue("Sheet1", "A16")
// if a16 != "" {
// 	total = "15"
// }
// b16, _ := file.GetCellValue("Sheet1", "B16")
// c16, _ := file.GetCellValue("Sheet1", "C16")
// d16, _ := file.GetCellValue("Sheet1", "D16")
// e16, _ := file.GetCellValue("Sheet1", "E16")
// f16, _ := file.GetCellValue("Sheet1", "F16")
// g16, _ := file.GetCellValue("Sheet1", "G16")

// a17, _ := file.GetCellValue("Sheet1", "A17")
// if a17 != "" {
// 	total = "16"
// }
// b17, _ := file.GetCellValue("Sheet1", "B17")
// c17, _ := file.GetCellValue("Sheet1", "C17")
// d17, _ := file.GetCellValue("Sheet1", "D17")
// e17, _ := file.GetCellValue("Sheet1", "E17")
// f17, _ := file.GetCellValue("Sheet1", "F17")
// g17, _ := file.GetCellValue("Sheet1", "G17")

// a18, _ := file.GetCellValue("Sheet1", "A18")
// if a18 != "" {
// 	total = "17"
// }
// b18, _ := file.GetCellValue("Sheet1", "B18")
// c18, _ := file.GetCellValue("Sheet1", "C18")
// d18, _ := file.GetCellValue("Sheet1", "D18")
// e18, _ := file.GetCellValue("Sheet1", "E18")
// f18, _ := file.GetCellValue("Sheet1", "F18")
// g18, _ := file.GetCellValue("Sheet1", "G18")

// a19, _ := file.GetCellValue("Sheet1", "A19")
// if a19 != "" {
// 	total = "18"
// }
// b19, _ := file.GetCellValue("Sheet1", "B19")
// c19, _ := file.GetCellValue("Sheet1", "C19")
// d19, _ := file.GetCellValue("Sheet1", "D19")
// e19, _ := file.GetCellValue("Sheet1", "E19")
// f19, _ := file.GetCellValue("Sheet1", "F19")
// g19, _ := file.GetCellValue("Sheet1", "G19")

// a20, _ := file.GetCellValue("Sheet1", "A20")
// if a20 != "" {
// 	total = "19"
// }
// b20, _ := file.GetCellValue("Sheet1", "B20")
// c20, _ := file.GetCellValue("Sheet1", "C20")
// d20, _ := file.GetCellValue("Sheet1", "D20")
// e20, _ := file.GetCellValue("Sheet1", "E20")
// f20, _ := file.GetCellValue("Sheet1", "F20")
// g20, _ := file.GetCellValue("Sheet1", "G20")

// a21, _ := file.GetCellValue("Sheet1", "A21")
// if a21 != "" {
// 	total = "20"
// }
// b21, _ := file.GetCellValue("Sheet1", "B21")
// c21, _ := file.GetCellValue("Sheet1", "C21")
// d21, _ := file.GetCellValue("Sheet1", "D21")
// e21, _ := file.GetCellValue("Sheet1", "E21")
// f21, _ := file.GetCellValue("Sheet1", "F21")
// g21, _ := file.GetCellValue("Sheet1", "G21")

//
// qa = []string{
// 	a2, strings.Join([]string{b2, c2, d2, e2}, "&&&"), f2, g2,
// 	a3, strings.Join([]string{b3, c3, d3, e3}, "&&&"), f3, g3,
// 	a4, strings.Join([]string{b4, c4, d4, e4}, "&&&"), f4, g4,
// 	a5, strings.Join([]string{b5, c5, d5, e5}, "&&&"), f5, g5,
// 	a6, strings.Join([]string{b6, c6, d6, e6}, "&&&"), f6, g6,
// 	a7, strings.Join([]string{b7, c7, d7, e7}, "&&&"), f7, g7,
// 	a8, strings.Join([]string{b8, c8, d8, e8}, "&&&"), f8, g8,
// 	a9, strings.Join([]string{b9, c9, d9, e9}, "&&&"), f9, g9,
// 	a10, strings.Join([]string{b10, c10, d10, e10}, "&&&"), f10, g10,
// 	a11, strings.Join([]string{b11, c11, d11, e11}, "&&&"), f11, g11,
// 	a12, strings.Join([]string{b12, c12, d12, e12}, "&&&"), f12, g12,
// 	a13, strings.Join([]string{b13, c13, d13, e13}, "&&&"), f13, g13,
// 	a14, strings.Join([]string{b14, c14, d14, e14}, "&&&"), f14, g14,
// 	a15, strings.Join([]string{b15, c15, d15, e15}, "&&&"), f15, g15,
// 	a16, strings.Join([]string{b16, c16, d16, e16}, "&&&"), f16, g16,
// 	a17, strings.Join([]string{b17, c17, d17, e17}, "&&&"), f17, g17,
// 	a18, strings.Join([]string{b18, c18, d18, e18}, "&&&"), f18, g18,
// 	a19, strings.Join([]string{b19, c19, d19, e19}, "&&&"), f19, g19,
// 	a20, strings.Join([]string{b20, c20, d20, e20}, "&&&"), f20, g20,
// 	a21, strings.Join([]string{b21, c21, d21, e21}, "&&&"), f21, g21,
// }
// 原始讀取檔案的題目設置

// GetGameResetPanel 遊戲資料重置
// func (s *SystemTable) GetGameResetPanel() (table Table) {
// 	table = DefaultTable(DefaultConfig())

// 	formList := table.GetForm()
// 	formList.SetUpdateFunc(func(values form2.Values) error {
// 		if values.IsEmpty("game_id") {
// 			return errors.New("錯誤: 遊戲ID發生問題，請輸入有效的遊戲ID")
// 		}

// 		if err := models.DefaultGameModel().SetDbConn(s.dbConn).SetRedisConn(s.redisConn).
// 			ResetGame(true, values.Get("game_id")); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return
// }

// QA1:          values.Get("qa_1"), QA1Picture: values.Get("qa_1_picture"), QA1Options: values.Get("qa_1_options"), QA1Answer: values.Get("qa_1_answer"),
// 				QA2: values.Get("qa_2"), QA2Picture: values.Get("qa_2_picture"), QA2Options: values.Get("qa_2_options"), QA2Answer: values.Get("qa_2_answer"),
// 				QA3: values.Get("qa_3"), QA3Picture: values.Get("qa_3_picture"), QA3Options: values.Get("qa_3_options"), QA3Answer: values.Get("qa_3_answer"),
// 				QA4: values.Get("qa_4"), QA4Picture: values.Get("qa_4_picture"), QA4Options: values.Get("qa_4_options"), QA4Answer: values.Get("qa_4_answer"),
// 				QA5: values.Get("qa_5"), QA5Picture: values.Get("qa_5_picture"), QA5Options: values.Get("qa_5_options"), QA5Answer: values.Get("qa_5_answer"),
// 				QA6: values.Get("qa_6"), QA6Picture: values.Get("qa_6_picture"), QA6Options: values.Get("qa_6_options"), QA6Answer: values.Get("qa_6_answer"),
// 				QA7: values.Get("qa_7"), QA7Picture: values.Get("qa_7_picture"), QA7Options: values.Get("qa_7_options"), QA7Answer: values.Get("qa_7_answer"),
// 				QA8: values.Get("qa_8"), QA8Picture: values.Get("qa_8_picture"), QA8Options: values.Get("qa_8_options"), QA8Answer: values.Get("qa_8_answer"),
// 				QA9: values.Get("qa_9"), QA9Picture: values.Get("qa_9_picture"), QA9Options: values.Get("qa_9_options"), QA9Answer: values.Get("qa_9_answer"),
// 				QA10: values.Get("qa_10"), QA10Picture: values.Get("qa_10_picture"), QA10Options: values.Get("qa_10_options"), QA10Answer: values.Get("qa_10_answer"),
// 				QA11: values.Get("qa_11"), QA11Picture: values.Get("qa_11_picture"), QA11Options: values.Get("qa_11_options"), QA11Answer: values.Get("qa_11_answer"),
// 				QA12: values.Get("qa_12"), QA12Picture: values.Get("qa_12_picture"), QA12Options: values.Get("qa_12_options"), QA12Answer: values.Get("qa_12_answer"),
// 				QA13: values.Get("qa_13"), QA13Picture: values.Get("qa_13_picture"), QA13Options: values.Get("qa_13_options"), QA13Answer: values.Get("qa_13_answer"),
// 				QA14: values.Get("qa_14"), QA14Picture: values.Get("qa_14_picture"), QA14Options: values.Get("qa_14_options"), QA14Answer: values.Get("qa_14_answer"),
// 				QA15: values.Get("qa_15"), QA15Picture: values.Get("qa_15_picture"), QA15Options: values.Get("qa_15_options"), QA15Answer: values.Get("qa_15_answe
