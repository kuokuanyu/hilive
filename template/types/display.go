package types

import (
	"fmt"
	"html/template"
	"strings"

	"hilive/template/form"
)

// FieldFunc 欄位處理函式
type FieldFunc func(value FieldModel) interface{}

// FieldDisplay 欄位處理
type FieldDisplay struct {
	DisplayFunc FieldFunc // 欄位處理
}

// FieldModel 主鍵值、資料表...等資料
type FieldModel struct {
	ID    string                 // 主鍵值
	Value string                 // 資料表值
	Row   map[string]interface{} // 資料表資料
}

// setDefaultDisplayFnOfSelect 設置DisplayFunc(類型為select)
func setDefaultDisplayFnOfSelect(f *FormPanel, typ form.Type) {
	if typ.IsSelect() {
		f.FieldList[f.curFieldListIndex].FieldDisplay.DisplayFunc = func(value FieldModel) interface{} {
			return strings.Split(value.Value, ",")
		}
	}
}

// ToDisplayFunc 執行欄位處理函式
func (f FieldDisplay) ToDisplayFunc(value FieldModel) interface{} {
	val := f.DisplayFunc(value)
	if f.IsNotSelect(val) {
		return FieldModel{
			Row:   value.Row,
			Value: fmt.Sprintf("%v", val),
			ID:    value.ID,
		}
	}
	return val
}

// DisplayFuncToHTML 執行欄位處理函式後轉換hmtl
func (f FieldDisplay) DisplayFuncToHTML(value FieldModel) template.HTML {
	v := f.DisplayFunc(value)
	if h, ok := v.(template.HTML); ok {
		return h
	} else if s, ok := v.(string); ok {
		return template.HTML(s)
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return template.HTML(arr[0])
	} else if v != nil {
		return ""
	} else {
		return ""
	}
}

// IsNotSelect 是否為選單欄位
func (f FieldDisplay) IsNotSelect(v interface{}) bool {
	switch v.(type) {
	case template.HTML:
		return false
	case []string:
		return false
	case [][]string:
		return false
	default:
		return true
	}
}
