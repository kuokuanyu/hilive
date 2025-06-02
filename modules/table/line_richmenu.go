package table

import (
	"encoding/json"
	"errors"

	"hilive/models"
	"hilive/modules/config"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
)

// GetActivityRequirePanel 舉辦活動需求
func (s *SystemTable) GetActivityRequirePanel() (table Table) {
	table = DefaultTable(DefaultConfig())
	formList := table.GetForm()

	// 新增舉辦活動需求
	formList.SetTable(config.ACTIVITY_REQUIRE_TABLE).
		SetInsertFunc(func(values form2.Values) error {
			if values.IsEmpty("user_id", "name", "phone", "email", "company_name",
				"activity_type", "people", "start_time", "end_time") {
				return errors.New("錯誤: 填寫資料發生問題，請輸入有效的資料")
			}

			// 將map[string][]string格是資料轉換為map[string]string
			flattened := utils.FlattenForm(values)

			// 將 map 轉成 JSON
			jsonBytes, err := json.Marshal(flattened)
			if err != nil {
				return err
			}

			// 轉成 struct
			var model models.EditActivityRequireModel
			if err := json.Unmarshal(jsonBytes, &model); err != nil {
				return err
			}

			if err := models.DefaultActivityRequireModel().
				SetConn(s.dbConn, s.redisConn, s.mongoConn).
				Add(model); err != nil {
				return err
			}
			return nil
		})

	return
}

// @Summary 新增舉辦活動需求
// @Tags LINE選單
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "LINE用戶ID"
// @param name formData string true "聯絡人(上限為20個字元)" minlength(1) maxlength(20)
// @param phone formData string true "聯絡電話(格式必須為台灣地區的手機號碼，如:09XXXXXXXX)" minlength(10) maxlength(10)
// @param email formData string true "聯絡信箱(電子郵件地址中必須包含「@」)"
// @param company_name formData string true "公司名稱(上限為20個字元)" minlength(1) maxlength(20)
// @param activity_type formData string true "活動類型" Enums(企業會議, 其他, 商業活動, 培訓/教育, 婚禮, 年會, 校園活動, 競技賽事, 論壇會議, 酒吧/餐飲娛樂, 電視/媒體)
// @param people formData integer true "活動人數" minimum(1)
// @param start_time formData string true "活動開始時間(西元年-月-日T時:分)"
// @param end_time formData string true "活動結束時間(西元年-月-日T時:分)"
// @param needs formData string false "需求"
// @param other formData string false "其他"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /activity/create/require [post]
func POSTActivityRequire() {
}

// models.EditActivityRequireModel{
// 	UserID:       values.Get("user_id"),
// 	Name:         values.Get("name"),
// 	Phone:        values.Get("phone"),
// 	Email:        values.Get("email"),
// 	CompanyName:  values.Get("company_name"),
// 	ActivityType: values.Get("activity_type"),
// 	People:       values.Get("people"),
// 	StartTime:    values.Get("start_time"),
// 	EndTime:      values.Get("end_time"),
// 	Needs:        values.Get("needs"),
// 	Other:        values.Get("other"),
// }
