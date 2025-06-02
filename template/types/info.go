package types

import (
	"strconv"

	"hilive/modules/db"
	"hilive/modules/parameter"
	"hilive/modules/service"
	"hilive/modules/utils"
)

// FieldList 所有欄位資訊
type FieldList []InfoField

// DeleteFunc 刪除POST API
type DeleteFunc func(ids []string, userID, activityID, gameID string) error

// JoinFieldValueDelimiter 分隔符號
var JoinFieldValueDelimiter = utils.UUID(8)

// InfoField 欄位資訊
type InfoField struct {
	Header       string          // 模板欄位名稱
	Field        string          // 資料表欄位名稱
	TypeName     db.DatabaseType // 資料表欄位類型
	Joins        Joins           // 關聯其他表
	CanEdit      bool            // 可編輯
	IsHide       bool            // 隱藏
	FieldDisplay FieldDisplay    // 欄位函式
}

// InfoPanel 資訊面板
type InfoPanel struct {
	FieldList         FieldList  // 欄位資訊
	curFieldListIndex int        // 欄位位置
	Table             string     // 資料表
	PageSizeList      []int      // 頁面顯示資料數
	PageSize          int        // 顯示頁數
	primaryKey        primaryKey // 主鍵
	DeleteFunc        DeleteFunc // 刪除 POST API
}

// TableInfo 資料表資訊
type TableInfo struct {
	Table      string // 資料表
	PrimaryKey string // 主鍵
	Delimiter  string // 分隔符號
	Delimiter2 string // 分隔符號
	Driver     string // 資料表引擎
}

// Joins 所有join資訊
type Joins []Join

// Join join資訊
type Join struct {
	BaseTable string // 原始表
	Field     string // 原始表欄位名稱
	JoinTable string // join表
	JoinField string // join表欄位
}

// InfoList 所有資料
type InfoList []map[string]InfoItem

// InfoItem 資料表資料
type InfoItem struct {
	Content interface{}
	Value   string // string值
}

// primaryKey 主鍵及主鍵type
type primaryKey struct {
	Type db.DatabaseType // type
	Name string          // 主鍵
}

// SetTable 設置資料表
func (i *InfoPanel) SetTable(table string) *InfoPanel {
	i.Table = table
	return i
}

// SetPrimaryKey 設置主鍵至InfoPanel
func (i *InfoPanel) SetPrimaryKey(name string, typ db.DatabaseType) *InfoPanel {
	i.primaryKey = primaryKey{
		Name: name,
		Type: typ,
	}
	return i
}

// DefaultInfoPanel 預設InfoPanel
func DefaultInfoPanel() *InfoPanel {
	return &InfoPanel{
		PageSizeList:      []int{10, 20, 30, 50, 100},
		curFieldListIndex: -1,
		PageSize:          1000,
	}
}

// AddField 增加欄位
func (i *InfoPanel) AddField(header, field string, typeName db.DatabaseType) *InfoPanel {
	i.FieldList = append(i.FieldList, InfoField{
		Header:   header,
		Field:    field,
		TypeName: typeName,
		Joins:    make(Joins, 0),
		CanEdit:  false,
		FieldDisplay: FieldDisplay{
			DisplayFunc: func(value FieldModel) interface{} {
				return value.Value
			},
		},
	})
	i.curFieldListIndex++
	return i
}

// GetFieldList 取得欄位資訊、join語法...等資訊
func (f FieldList) GetFieldList(info TableInfo, params parameter.Parameters, columns []string, sql ...func(services service.List) *db.SQL) (
	FieldList, string, string, string, []string) {
	var (
		fieldList  = make(FieldList, 0)
		fields     = ""
		joinFields = ""                // ex: group_concat(roles.`name` separator 'CkN694kH') as roles_join_name,
		joins      = ""                // ex: left join `role_users` on role_users.`user_id` = users.`id` left join....
		joinTables = make([]string, 0) // ex:{role_users, roles}
		tableName  = info.Delimiter + info.Table + info.Delimiter2
	)

	for _, field := range f {
		if field.Field != info.PrimaryKey && utils.InArray(columns, field.Field) &&
			!field.Joins.Valid() {
			fields += tableName + "." + utils.FilterField(field.Field, info.Delimiter, info.Delimiter2) + ","
		}

		headField := field.Field
		if field.Joins.Valid() {
			headField = field.Joins.Last().GetTableName() + "_join_" + field.Field

			joinFields += db.GetAggregationExpression(info.Driver,
				field.Joins.Last().GetTableName(info.Delimiter, info.Delimiter2)+"."+
					utils.FilterField(field.Field, info.Delimiter, info.Delimiter2),
				headField, JoinFieldValueDelimiter) + ","

			for _, join := range field.Joins {
				if !utils.InArray(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2)) {
					joinTables = append(joinTables, join.GetTableName(info.Delimiter, info.Delimiter2))
					if join.BaseTable == "" {
						join.BaseTable = info.Table
					}

					joins += " left join " + utils.FilterField(join.JoinTable, info.Delimiter, info.Delimiter2) + " on " +
						join.GetTableName(info.Delimiter, info.Delimiter2) + "." +
						utils.FilterField(join.JoinField, info.Delimiter, info.Delimiter2) + " = " +
						utils.Delimiter(info.Delimiter, info.Delimiter2, join.BaseTable) + "." +
						utils.FilterField(join.Field, info.Delimiter, info.Delimiter2)
				}
			}
		}

		if field.IsHide {
			continue
		}

		fieldList = append(fieldList, InfoField{
			Header:  field.Header,
			Field:   headField,
			IsHide:  !utils.InArrayWithoutEmpty(params.Columns, headField),
			CanEdit: field.CanEdit,
		})
	}

	return fieldList, fields, joinFields, joins, joinTables
}

// GetField 取得Field資訊
func (f FieldList) GetField(name string) InfoField {
	for _, field := range f {
		if field.Field == name {
			return field
		}
		if JoinField(field.Joins.Last().JoinTable, field.Field) == name {
			return field
		}
	}
	return InfoField{}
}

// FieldJoin 欄位有關聯其他表
func (i *InfoPanel) FieldJoin(join Join) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Joins = append(i.FieldList[i.curFieldListIndex].Joins, join)
	return i
}

// FieldHide 隱藏欄位
func (i *InfoPanel) FieldHide() *InfoPanel {
	i.FieldList[i.curFieldListIndex].IsHide = true
	return i
}

// SetDisplayFunc 資料處理函式
func (i *InfoPanel) SetDisplayFunc(f FieldFunc) *InfoPanel {
	i.FieldList[i.curFieldListIndex].FieldDisplay.DisplayFunc = f
	return i
}

// SetDeleteFunc 刪除 POST API
func (i *InfoPanel) SetDeleteFunc(fn DeleteFunc) *InfoPanel {
	i.DeleteFunc = fn
	return i
}

// Valid 是否join其他表
func (j Joins) Valid() bool {
	for i := 0; i < len(j); i++ {
		if j[i].JoinTable != "" && j[i].Field != "" && j[i].JoinField != "" && j[i].BaseTable != "" {
			return true
		}
	}
	return false
}

// Last 回傳Join
func (j Joins) Last() Join {
	if len(j) > 0 {
		return j[len(j)-1]
	}
	return Join{}
}

// JoinField return 資料表_join_欄位
func JoinField(table, field string) string {
	return table + "_join_" + field
}

// GetTableName 將table名稱加上分隔符號
func (j Join) GetTableName(delimiter ...string) string {
	if len(delimiter) > 0 {
		return delimiter[0] + j.JoinTable + delimiter[1]
	}
	return j.JoinTable
}

// GetPageSizeList 單頁顯示資料
func (i *InfoPanel) GetPageSizeList() []string {
	var pageSizeList = make([]string, len(i.PageSizeList))
	for j := 0; j < len(i.PageSizeList); j++ {
		pageSizeList[j] = strconv.Itoa(i.PageSizeList[j])
	}
	return pageSizeList
}
