package form

import (
	"html/template"
)

// Type uint8
type Type uint8

const (
	Default Type = iota
	Text
	SelectSingle
	Select
	IconPicker
	SelectBox
	File
	Multifile
	Password
	RichText
	Datetime
	DatetimeRange
	Time
	Radio
	Checkbox
	CheckboxNumber
	CheckboxStacked
	CheckboxSingle
	Email
	Date
	DateRange
	URL
	IP
	Color
	Array
	Currency
	Rate
	Number
	Table
	// NumberRange
	TextArea
	Custom
	Switch
	Code
	Slider
	Range
	RadioToggle
	RadioFile
	MultiCheckbox
	SelectCity
	NoUse
)

// IsSelect 是否為select
func (t Type) IsSelect() bool {
	return t == Select || t == SelectSingle || t == SelectBox || t == Radio ||
		t == Switch || t == CheckboxStacked || t == CheckboxSingle || t == RadioToggle || t == RadioFile
}

// SelectedLabel 判斷select與check型態
func (t Type) SelectedLabel() []template.HTML {
	if t == Select || t == SelectSingle || t == SelectBox {
		return []template.HTML{"selected", ""}
	}
	if t == Radio || t == Switch || t == Checkbox || t == CheckboxStacked || t == CheckboxSingle || t == RadioToggle || t == RadioFile {
		return []template.HTML{"checked", ""}
	}
	return []template.HTML{"", ""}
}

// IsRange 是否為範圍
func (t Type) IsRange() bool {
	return t == DatetimeRange
	// return t == DatetimeRange || t == NumberRange
}

// IsSelectCity 是否為城市類型
func (t Type) IsSelectCity() bool {
	return t == SelectCity
}

// IsCheckboxNumber 是否為Checkbox+Number
func (t Type) IsCheckboxNumber() bool {
	return t == CheckboxNumber
}

// IsMultiFile 是否多個檔案
func (t Type) IsMultiFile() bool {
	return t == Multifile
}

// IsRadioToggle 是否為IsRadioToggle
func (t Type) IsRadioToggle() bool {
	return t == RadioToggle
}

// IsRadioFile 是否為IsRadioFile
func (t Type) IsRadioFile() bool {
	return t == RadioFile
}

// IsMultiCheckbox 是否為MultiCheckbox
func (t Type) IsMultiCheckbox() bool {
	return t == MultiCheckbox
}

// String 將type轉換成string
func (t Type) String() string {
	switch t {
	case Default:
		return "default"
	case Text:
		return "text"
	case SelectSingle:
		return "select_single"
	case Select:
		return "select"
	case IconPicker:
		return "iconpicker"
	case SelectBox:
		return "selectbox"
	case File:
		return "file"
	case Table:
		return "table"
	case Multifile:
		return "multi_file"
	case Password:
		return "password"
	case RichText:
		return "richtext"
	case Rate:
		return "rate"
	case Checkbox:
		return "checkbox"
	case CheckboxStacked:
		return "checkbox_stacked"
	case CheckboxSingle:
		return "checkbox_single"
	case CheckboxNumber:
		return "checkbox_number"
	case Date:
		return "date"
	case Time:
		return "time"
	case Range:
		return "range"
	case DateRange:
		return "date_range"
	case Datetime:
		return "datetime"
	case DatetimeRange:
		return "datetime_range"
	case Radio:
		return "radio"
	case Slider:
		return "slider"
	case Array:
		return "array"
	case Email:
		return "email"
	case URL:
		return "url"
	case IP:
		return "ip"
	case Color:
		return "color"
	case Currency:
		return "currency"
	case Number:
		return "number"
	// case NumberRange:
	// 	return "number_range"
	case TextArea:
		return "textarea"
	case Custom:
		return "custom"
	case Switch:
		return "switch"
	case Code:
		return "code"
	case RadioToggle:
		return "radio_toggle"
	case RadioFile:
		return "radio_file"
	case MultiCheckbox:
		return "multi_checkbox"
	case SelectCity:
		return "select_city"
	case NoUse:
		return "no_use"
	default:
		panic("wrong form type")
	}
}
