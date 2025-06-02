package types

import (
	"fmt"
	"html/template"
	"strings"

	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/form"
)

const (
	UPLOAD_RADIO_URL = "/admin/uploads/system/radio/"

// UPLOAD_SYSTEM_URL = "/admin/uploads/system/"
// UPLOAD_URL        = "/admin/uploads/"
)

// FormPostFunc Post API 功能
type FormPostFunc func(values form2.Values) error

// PostType
type PostType uint8

// FormPanel 紀錄所有欄位資訊、插入更新POST...等功能
type FormPanel struct {
	FieldList         FormFields   // 欄位資訊
	curFieldListIndex int          // 欄位位置
	Table             string       // 資料表
	InsertFunc        FormPostFunc // 插入資料功能
	UpdateFunc        FormPostFunc // 更新資料功能
	primaryKey        primaryKey   // 主鍵
}

// FormFields 所有表單資訊
type FormFields []FormField

// FormField 表單欄位資訊
type FormField struct {
	Field        string          // 資料表欄位名稱
	Name         string          // 前端模板欄位名稱
	TypeName     db.DatabaseType // 資料表欄位類型
	Header       string          // 模板欄位名稱
	FormType     form.Type       // 表單類型
	Unit         string          // 單位
	MaxValue     int64           // 最大值
	Value        template.HTML   // Value1
	Value2       template.HTML   // Value2
	Value3       interface{}     // Value3
	Placeholder  []string        // 表單顯示資訊
	CanEdit      bool            // 允許編輯
	CanAdd       bool            // 允許增加
	Must         bool            // 該欄位必填
	IsHide       bool            // 是否隱藏
	IsOpen       bool            // 是否開啟
	Default      template.HTML   // 預設值
	Joins        Joins           // 關聯表
	FieldDisplay FieldDisplay    // 表單函式
	FieldOptions FieldOptions    // 選單
	OptionTable  OptionTable     // 選單資訊(關聯其他表)
	HelpMsg      template.HTML   // 提示訊息
}

// FieldOptions 所有選單資訊
type FieldOptions []FieldOption

// FieldOption 選單資訊
type FieldOption struct {
	Text          string        // 選單名稱
	Value         string        // 選單值
	Text2         string        // 選單名稱
	IsSelected    bool          // 是否被選擇
	SelectedLabel template.HTML // 選項的label
}

// OptionTable 選單資訊(關聯其他表)
type OptionTable struct {
	Table      string // 關聯表
	TextField  string // 選單名稱
	ValueField string // 選單值
}

// DefaultFormPanel 預設FormPanel
func DefaultFormPanel() *FormPanel {
	return &FormPanel{
		curFieldListIndex: -1,
	}
}

// SetPrimaryKey 設置主鍵
func (f *FormPanel) SetPrimaryKey(name string, typ db.DatabaseType) *FormPanel {
	f.primaryKey = primaryKey{Name: name, Type: typ}
	return f
}

// AddField 增加欄位資訊
func (f *FormPanel) AddField(header, field string, fieldType db.DatabaseType, formType form.Type) *FormPanel {
	form := FormField{
		Header:      header,
		Name:        "",
		Field:       field,
		TypeName:    fieldType,
		CanAdd:      true,
		CanEdit:     true,
		IsHide:      false,
		IsOpen:      true,
		Placeholder: []string{"請輸入" + header},
		Must:        true,
		FormType:    formType,
		FieldDisplay: FieldDisplay{
			DisplayFunc: func(value FieldModel) interface{} {
				return value.Value
			},
		},
	}
	f.FieldList = append(f.FieldList, form)
	f.curFieldListIndex++
	setDefaultDisplayFnOfSelect(f, formType)
	return f
}

// FieldWithValue 取得表單欄位資訊(帶有預設值)
func (f *FormPanel) FieldsWithValue(fields, values, columns []string, res map[string]interface{},
	services service.List, sql func(services service.List) *db.SQL) FormFields {
	var (
		list = make(FormFields, 0)
	)

	for _, field := range f.FieldList {
		dataValue := field.GetDataValue(columns, res[field.Field])
		list = append(list, *(field.UpdateValue(values[0], dataValue, res, sql(services))))
	}

	for i := 0; i < len(fields); i++ {
		list = append(list, FormField{
			Header:   fields[i],
			Field:    fields[i],
			Value:    template.HTML(values[i]),
			FormType: form.Default,
			IsHide:   true,
		})
	}
	return list
}

// FieldsWithDefaultValue 取得表單欄位資訊(允許增加的欄位)，預設值為空
func (f *FormPanel) FieldsWithDefaultValue(keys, value []string, services service.List, sql ...func(services service.List) *db.SQL) FormFields {
	var (
		list = make(FormFields, 0)
	)
	for _, v := range f.FieldList {
		if v.CanAdd {
			v.CanEdit = true
			if len(sql) > 0 {
				v.Value = v.Default
				list = append(list, *(v.UpdateValue("", string(v.Value), make(map[string]interface{}), sql[0](services))))
			} else {
				v.Value = v.Default
				list = append(list, *(v.UpdateValue("", string(v.Value), make(map[string]interface{}), nil)))
			}
		}
	}

	for i := range keys {
		list = append(list, FormField{
			Header:   keys[i],
			Field:    keys[i],
			Value:    template.HTML(value[i]),
			FormType: form.Default,
			IsHide:   true,
		})
	}
	return list
}

// SetTable 設置Table
func (f *FormPanel) SetTable(table string) *FormPanel {
	f.Table = table
	return f
}

// SetName 設置Name
func (f *FormPanel) SetName(name string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Name = name
	return f
}

// SetUnit 設置單位
func (f *FormPanel) SetUnit(unit string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Name = unit
	return f
}

// SetMaxValue 設置最大值
func (f *FormPanel) SetMaxValue(max int64) *FormPanel {
	f.FieldList[f.curFieldListIndex].MaxValue = max
	return f
}

// FieldCanNotAdd 不允許增加
func (f *FormPanel) FieldCanNotAdd() *FormPanel {
	f.FieldList[f.curFieldListIndex].CanAdd = false
	return f
}

// FieldCanNotEdit 不允許編輯
func (f *FormPanel) FieldCanNotEdit() *FormPanel {
	f.FieldList[f.curFieldListIndex].CanEdit = false
	return f
}

// FieldNotMust 表單不是必填
func (f *FormPanel) FieldNotMust() *FormPanel {
	f.FieldList[f.curFieldListIndex].Must = false
	return f
}

// FieldHide 表單隱藏
func (f *FormPanel) FieldHide() *FormPanel {
	f.FieldList[f.curFieldListIndex].IsHide = true
	return f
}

// FieldNotOpen 不開啟
func (f *FormPanel) FieldNotOpen() *FormPanel {
	f.FieldList[f.curFieldListIndex].IsOpen = false
	return f
}

// SetPlaceholder 設置提示訊息
func (f *FormPanel) SetPlaceholder(placeholder []string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Placeholder = placeholder
	return f
}

// SetFieldOptions 選單
func (f *FormPanel) SetFieldOptions(options FieldOptions) *FormPanel {
	f.FieldList[f.curFieldListIndex].FieldOptions = options
	return f
}

// SetInsertFunc 新增POST API
func (f *FormPanel) SetInsertFunc(fn FormPostFunc) *FormPanel {
	f.InsertFunc = fn
	return f
}

// SetUpdateFunc 更新POST API
func (f *FormPanel) SetUpdateFunc(fn FormPostFunc) *FormPanel {
	f.UpdateFunc = fn
	return f
}

// FieldOptionsFromTable 選單(關聯其他表)
func (f *FormPanel) FieldOptionsFromTable(table, textFieldName, valueFieldName string) *FormPanel {
	f.FieldList[f.curFieldListIndex].OptionTable = OptionTable{
		Table:      table,
		TextField:  textFieldName,
		ValueField: valueFieldName,
	}
	return f
}

// SetDisplayFunc 表單funciton
func (f *FormPanel) SetDisplayFunc(fc FieldFunc) *FormPanel {
	f.FieldList[f.curFieldListIndex].FieldDisplay.DisplayFunc = fc
	return f
}

// SetFieldDefault 預設值
func (f *FormPanel) SetFieldDefault(def string) *FormPanel {
	f.FieldList[f.curFieldListIndex].Default = template.HTML(def)
	return f
}

// UpdateValue 處理預設值及選單
func (f *FormField) UpdateValue(id, val string, res map[string]interface{}, s *db.SQL) *FormField {
	m := FieldModel{
		ID:    id,
		Value: val,
		Row:   res,
	}

	if f.Field == "city" {
		// 活動縣市及區域
		f.Value = template.HTML(fmt.Sprintf("%s", res["city"]))
		f.Value2 = template.HTML(fmt.Sprintf("%s", res["town"]))
	} else if f.Field == "start_time" {
		// 活動時間
		start := strings.Replace(fmt.Sprintf("%s", res["start_time"]), " ", "T", 1)
		end := strings.Replace(fmt.Sprintf("%s", res["end_time"]), " ", "T", 1)
		f.Value = template.HTML(start[:len(start)-3])
		f.Value2 = template.HTML(end[:len(start)-3])
	} else if f.Field == "picture_start_time" {
		// 圖片牆時間
		start := fmt.Sprintf("%s", res["picture_start_time"])
		end := fmt.Sprintf("%s", res["picture_end_time"])
		f.Value = template.HTML(start)
		f.Value2 = template.HTML(end)
	} else if f.Field == "message_ban" {
		// 禁止連續傳送訊息
		f.Value = template.HTML(fmt.Sprintf("%s", res["message_ban"]))
		f.Value2 = template.HTML(fmt.Sprintf("%d", res["message_ban_second"]))
	} else if f.Field == "danmu_top" {
		// 彈幕位置
		var position = make([]interface{}, 0)

		f.Value3 = append(position,
			[]string{"頂部", "danmu_top", fmt.Sprintf("%s", res["danmu_top"])},
			[]string{"中間", "danmu_mid", fmt.Sprintf("%s", res["danmu_mid"])},
			[]string{"底部", "danmu_bottom", fmt.Sprintf("%s", res["danmu_bottom"])})
	} else if f.Field == "special_danmu_topic" {
		// 超級彈幕主題
		var topic = make([]interface{}, 0)

		f.Value3 = append(topic,
			[]string{"默認主題", "special_danmu_topic", fmt.Sprintf("%s", res["special_danmu_topic"])})
		// []string{"主題一", "topic1", fmt.Sprintf("%s", model["topic1"]), UPLOAD_RADIO_URL + "superdanmu-style1.png"},
		// []string{"主題二", "topic2", fmt.Sprintf("%s", model["topic2"]), UPLOAD_RADIO_URL + "superdanmu-style2.png"},
		// []string{"主題三", "topic3", fmt.Sprintf("%s", model["topic3"]), UPLOAD_RADIO_URL + "superdanmu-style3.png"}))
	} else if f.Field == "picture" {
		// 圖片牆圖片
		f.Value3 = strings.Split(fmt.Sprintf("%s", res["picture"]), "\n")
	} else if f.Field == "holdscreen_birthday_topic" {
		// 霸屏主題
		var topic = make([]interface{}, 0)

		f.Value3 = append(topic,
			[]string{"預設", "default", "", UPLOAD_RADIO_URL + "danmu-position-style1.pn"},
			[]string{"生日主題", "birthday", fmt.Sprintf("%s", res["holdscreen_birthday_topic"]), UPLOAD_RADIO_URL + "holdscreen-style1.png"})
	} else if f.Field == "threed_background_style" {
		f.setFieldOptions(s)
		f.FieldOptions.SetSelected(f.FieldDisplay.ToDisplayFunc(m), f.FormType.SelectedLabel())
		// 3D簽到牆背景
		f.Value = template.HTML(fmt.Sprintf("%d", res["threed_background_style"]))
		f.Value2 = template.HTML(fmt.Sprintf("%s", res["threed_background"]))
	} else if f.Field == "topic" {
		// 遊戲主題樣式
		f.Value = template.HTML(fmt.Sprintf("%s", res["topic"]))
		f.Value2 = template.HTML(fmt.Sprintf("%s", res["skin"]))
	} else if f.FormType.IsSelect() {
		f.setFieldOptions(s)
		f.FieldOptions.SetSelected(f.FieldDisplay.ToDisplayFunc(m), f.FormType.SelectedLabel())
		// if f.FormType.IsRadioToggle() || f.FormType.IsRadioFile() {
		// 	if valArr, ok := f.FieldDisplay.ToDisplayFunc(m).([]string); ok {
		// 		f.Value = template.HTML(valArr[0])
		// 		f.Value2 = template.HTML(valArr[1])
		// 	}
		// }
		// }else if f.FormType.IsRange() || f.FormType.IsCheckboxNumber() {
		// 	// } else if f.FormType.IsRange() || f.FormType.IsSelectCity() || f.FormType.IsCheckboxNumber() {
		// 	v := f.FieldDisplay.DisplayFunc(m)
		// 	if v != nil {
		// 		vals := strings.Split(fmt.Sprintf("%s", v), ",")
		// 		f.Value = template.HTML(vals[0])
		// 		f.Value2 = template.HTML(vals[1])
		// 	}
		// } else if f.FormType.IsMultiCheckbox() || f.FormType.IsMultiFile() {
		// 	v := f.FieldDisplay.DisplayFunc(m) // 陣列
		// 	if v != nil {
		// 		f.Value3 = v
		// 	}
	} else {
		f.Value = f.FieldDisplay.DisplayFuncToHTML(m)
	}
	return f
}

// setFieldOptions 選單資訊(關聯其他表)
func (f *FormField) setFieldOptions(s *db.SQL) {
	if s != nil && f.OptionTable.Table != "" && len(f.FieldOptions) == 0 {
		queryRes, err := s.Table(f.OptionTable.Table).
			Select(f.OptionTable.ValueField, f.OptionTable.TextField).All() // 取得選單名稱與值
		if err == nil {
			for _, item := range queryRes {
				f.FieldOptions = append(f.FieldOptions, FieldOption{
					Value: fmt.Sprintf("%s", item[f.OptionTable.ValueField]), // 值
					Text:  fmt.Sprintf("%v", item[f.OptionTable.TextField]),  // 名稱
				})
			}
		}
	}
}

// GetDataValue 欄位值
func (f *FormField) GetDataValue(columns []string, v interface{}) string {
	return utils.AorB(utils.InArray(columns, f.Field),
		string(db.GetValueFromDatabaseType(f.TypeName, v)), "")
}

// SetFieldHelpMsg 提示訊息
func (f *FormPanel) SetFieldHelpMsg(s template.HTML) *FormPanel {
	f.FieldList[f.curFieldListIndex].HelpMsg = s
	return f
}

// SetSelected 處理選單(判斷是否被選取)
func (f FieldOptions) SetSelected(val interface{}, labels []template.HTML) FieldOptions {
	if valArr, ok := val.([]string); ok {
		for k := range f {
			text := f[k].Text

			if text == "" {
				text = f[k].Text2
			}
			f[k].IsSelected = utils.InArray(valArr, f[k].Value) || utils.InArray(valArr, text)
			if f[k].IsSelected {
				f[k].SelectedLabel = labels[0]
			} else {
				f[k].SelectedLabel = labels[1]
			}
		}
	} else {
		for k := range f {
			text := f[k].Text
			if text == "" {
				text = f[k].Text2
			}
			f[k].IsSelected = f[k].Value == val || text == val
			if f[k].IsSelected {
				f[k].SelectedLabel = labels[0]
			} else {
				f[k].SelectedLabel = labels[1]
			}
		}
	}
	return f
}
