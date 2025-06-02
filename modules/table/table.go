package table

import (
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/form"
	form2 "hilive/modules/form"

	"hilive/modules/mongo"
	"hilive/modules/paginator"
	"hilive/modules/parameter"
	"hilive/modules/service"
	"hilive/modules/utils"
	"hilive/template/types"

	"fmt"
	"strconv"
	"strings"
)

const (
	UPLOAD_RADIO_URL  = "/admin/uploads/system/radio/"
	UPLOAD_SYSTEM_URL = "/admin/uploads/system/"
	UPLOAD_URL        = "/admin/uploads/"
	DEFAULT_FALG      = "__default_flag"
)

// BaseTable 包含面板、表單...等資訊
type BaseTable struct {
	Info        *types.InfoPanel // 面板資訊
	Form        *types.FormPanel // 表單資訊
	SettingForm *types.FormPanel // 基本設置的表單資訊
	canAdd      bool             // 可新增
	canEdit     bool             // 可編輯
	canDelete   bool             // 可刪除
	primaryKey  PrimaryKey       // 主鍵
	driver      string           // 資料庫引擎
	connName    string           // 資料庫名稱
	connection  db.Connection    // 資料庫處理程序
}

// PrimaryKey 紀錄主鍵
type PrimaryKey struct {
	Name string
	Type db.DatabaseType
}

// PanelInfo 頁面資訊
type PanelInfo struct {
	FieldList types.FieldList     // 欄位資訊，可編輯、編輯選項、是否隱藏...等資訊
	InfoList  types.InfoList      // 所有資料
	Paginator paginator.Paginator // 分頁
}

// FormInfo 表單資訊
type FormInfo struct {
	FieldList types.FormFields // 所有欄位資訊
}

// Generator 取得方法
type Generator func() Table

// List 儲存所有頁面table
type List map[string]Generator

// Table interface
type Table interface {
	GetInfo() *types.InfoPanel                                                                // 設置Info主鍵
	GetForm() *types.FormPanel                                                                // 設置Form主鍵
	GetSettingForm() *types.FormPanel                                                         // 設置SettingForm主鍵                                                        // 設置FormInfo主鍵
	GetNewFormInfo(services service.List, param parameter.Parameters, keys []string) FormInfo // 取得新增的表單資訊
	GetEditFormInfo(param parameter.Parameters, services service.List,
		pkFields []string) (FormInfo, error) // 取得編輯的表單資訊並設置預設值
	GetSettingFormInfo(param parameter.Parameters, services service.List,
		pkFields []string) (FormInfo, error) // 取得基本設置的表單資訊並設置預設值
	GetData(params parameter.Parameters, services service.List) (PanelInfo, error) // 從資料庫取得資訊面板需要顯示的資料
	InsertData(dataList form.Values) error                                         // 插入資料
	UpdateData(dataList form.Values) error                                         // 更新資料
	UpdateSettingData(dataList form.Values) error                                  // 更新基本設置資料
	DeleteData(pk string, userID, activityID, gameID string) error                 // 刪除資料
}

// SystemTable 設置模板資料、欄位資訊
type SystemTable struct {
	dbConn    db.Connection
	redisConn cache.Connection
	mongoConn mongo.Connection
	config    *config.Config
}

type PictureField struct {
	FieldName string // 給資料庫的欄位名，ex: 3d_gacha_machine_classic_h_pic_02
	Path      string // 圖片或音樂檔案路徑模板，ex: 3DGachaMachine/classic/3d_gacha_machine_classic_h_pic_02.png
}

// BuildPictureMap 新增或編輯時，判斷圖片參數是否為空，將路徑參數寫入map中
func BuildPictureMap(fields []PictureField, values form2.Values, isNew bool) map[string]string {
	pictureMap := make(map[string]string)
	topic := values.Get("topic")

	for _, f := range fields {
		val := values.Get(f.FieldName)
		defaultFlag := values.Get(f.FieldName + DEFAULT_FALG)

		// 編輯時使用 defaultFlag = "1" 才使用預設圖片
		// 新增時當欄位為空就使用預設圖片
		if (isNew && val == "") || (!isNew && defaultFlag == "1") {
			// 使用預設圖
			if strings.Contains(f.Path, "%s") {
				parts := strings.Split(topic, "_")
				suffix := ""
				if len(parts) == 2 {
					suffix = parts[1]
				} else if len(parts) == 3 {
					suffix = parts[1] + "_" + parts[2]
				}
				pictureMap[f.FieldName] = UPLOAD_SYSTEM_URL + fmt.Sprintf(f.Path, suffix)
			} else {
				pictureMap[f.FieldName] = UPLOAD_SYSTEM_URL + f.Path
			}
		} else if val != "" {
			// 使用使用者上傳圖
			pictureMap[f.FieldName] = val
		}
	}
	return pictureMap
}

// NewSystemTable 預設SystemTable
func NewSystemTable(conn db.Connection, redis cache.Connection, mongoConn mongo.Connection, c *config.Config) *SystemTable {
	return &SystemTable{dbConn: conn, redisConn: redis, mongoConn: mongoConn, config: c}
}

// table 取得SQL
func (s *SystemTable) table(table string) *db.SQL {
	return s.connection().Table(table)
}

// connection 設置SQL
func (s *SystemTable) connection() *db.SQL {
	return db.Conn(s.dbConn)
}

// DefaultTable 預設Table
func DefaultTable(cfgs ...Config) Table {
	var cfg Config
	if len(cfgs) > 0 && cfgs[0].primaryKey.Name != "" {
		cfg = cfgs[0]
	} else {
		cfg = DefaultConfig()
	}
	return &BaseTable{
		Info:        types.DefaultInfoPanel(),
		SettingForm: types.DefaultFormPanel(),
		Form:        types.DefaultFormPanel(),
		canAdd:      cfg.canAdd,
		canEdit:     cfg.canEdit,
		canDelete:   cfg.canDelete,
		primaryKey:  cfg.primaryKey,
		driver:      cfg.driver,
		connName:    cfg.conn,
	}
}

// GetInfo 設置InfoPanel主鍵
func (base *BaseTable) GetInfo() *types.InfoPanel {
	return base.Info.SetPrimaryKey(base.primaryKey.Name, base.primaryKey.Type)
}

// GetForm 設置FormPanel主鍵
func (base *BaseTable) GetForm() *types.FormPanel {
	return base.Form.SetPrimaryKey(base.primaryKey.Name, base.primaryKey.Type)
}

// GetSettingForm 設置SettingForm主鍵
func (base *BaseTable) GetSettingForm() *types.FormPanel {
	return base.SettingForm.SetPrimaryKey(base.primaryKey.Name, base.primaryKey.Type)
}

// GetNewFormInfo 取得新增的表單資訊
func (base *BaseTable) GetNewFormInfo(services service.List, param parameter.Parameters, keys []string) FormInfo {
	return FormInfo{FieldList: base.Form.FieldsWithDefaultValue(keys, param.FindPKs(), services, base.sql)}
}

// GetEditFormInfo 取得編輯的表單資訊並設置預設值
func (base *BaseTable) GetEditFormInfo(param parameter.Parameters, services service.List, pkFields []string) (FormInfo, error) {
	var (
		connection                             = base.db(services)
		delimiter                              = connection.GetDelimiter()
		delimiter2                             = connection.GetDelimiter2()
		tableName                              = utils.Delimiter(delimiter, delimiter2, base.GetForm().Table)
		pks                                    = param.FindPKs()
		args                                   = []interface{}{}
		fields, joins, groupBy, queryStatement string
		columns, _                             = base.getColumns(base.Form.Table, services)
		result                                 []map[string]interface{}
		err                                    error
	)

	if len(pks) == 2 {
		args = []interface{}{pks[0], pks[1]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? and " + pkFields[1] + " = ? %s "
	} else if len(pks) == 1 {
		args = []interface{}{pks[0]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? %s "
	} else if len(pks) == 3 {
		args = []interface{}{pks[0], pks[1], pks[2]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? and " +
			pkFields[1] + " = ? and " + pkFields[2] + " = ? %s "
	}

	for i, field := range base.Form.FieldList {
		if utils.InArray(columns, field.Field) && !field.Joins.Valid() {
			fields += tableName + "." +
				utils.FilterField(base.Form.FieldList[i].Field, delimiter, delimiter2) + ","
		}
	}
	for i := range pkFields {
		fields += tableName + "." + utils.FilterField(pkFields[i], delimiter, delimiter2) + ", "
	}

	queryCmd := fmt.Sprintf(queryStatement, fields[:len(fields)-2], tableName, joins, groupBy)
	if result, err = connection.QueryWithConnection(base.connName,
		queryCmd, args...); err != nil || len(result) == 0 {
		return FormInfo{}, err
	}

	return FormInfo{
		FieldList: base.Form.FieldsWithValue(pkFields, pks,
			columns, result[0], services, base.sql),
	}, nil
}

// GetSettingFormInfo 取得基本設置的表單資訊並設置預設值
func (base *BaseTable) GetSettingFormInfo(param parameter.Parameters, services service.List, pkFields []string) (FormInfo, error) {
	var (
		connection                             = base.db(services)
		delimiter                              = connection.GetDelimiter()
		delimiter2                             = connection.GetDelimiter2()
		tableName                              = utils.Delimiter(delimiter, delimiter2, base.GetSettingForm().Table)
		pks                                    = param.FindPKs()
		args                                   = []interface{}{}
		fields, joins, groupBy, queryStatement string
		columns, _                             = base.getColumns(base.SettingForm.Table, services)
		result                                 []map[string]interface{}
		err                                    error
	)

	if len(pks) == 2 {
		args = []interface{}{pks[0], pks[1]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? and " + pkFields[1] + " = ? %s "
	} else if len(pks) == 1 {
		args = []interface{}{pks[0]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? %s "
	} else if len(pks) == 3 {
		args = []interface{}{pks[0], pks[1], pks[2]}
		queryStatement = "select %s from %s" + " %s where " + pkFields[0] + " = ? and " +
			pkFields[1] + " = ? and " + pkFields[2] + " = ? %s "
	}

	for i, field := range base.SettingForm.FieldList {
		if utils.InArray(columns, field.Field) && !field.Joins.Valid() {
			fields += tableName + "." +
				utils.FilterField(base.SettingForm.FieldList[i].Field, delimiter, delimiter2) + ","
		}
	}
	for i := range pkFields {
		fields += tableName + "." + utils.FilterField(pkFields[i], delimiter, delimiter2) + ", "
	}

	queryCmd := fmt.Sprintf(queryStatement, fields[:len(fields)-2], tableName, joins, groupBy)
	if result, err = connection.QueryWithConnection(base.connName,
		queryCmd, args...); err != nil || len(result) == 0 {
		return FormInfo{}, err
	}

	return FormInfo{
		FieldList: base.SettingForm.FieldsWithValue(pkFields, pks,
			columns, result[0], services, base.sql),
	}, nil
}

// GetData 取得頁面所有資料
func (base *BaseTable) GetData(params parameter.Parameters, services service.List) (PanelInfo, error) {
	return base.getDataFromDatabase(params, services)
}

// InsertData 插入資料
func (base *BaseTable) InsertData(dataList form.Values) error {
	if base.Form.InsertFunc != nil {
		if err := base.Form.InsertFunc(dataList); err != nil {
			return err
		}
	}
	return nil
}

// UpdateData 更新資料
func (base *BaseTable) UpdateData(dataList form.Values) error {
	if base.Form.UpdateFunc != nil {
		if err := base.Form.UpdateFunc(dataList); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSettingData 更新基本設置資料
func (base *BaseTable) UpdateSettingData(dataList form.Values) error {
	if base.SettingForm.UpdateFunc != nil {
		if err := base.SettingForm.UpdateFunc(dataList); err != nil {
			return err
		}
	}
	return nil
}

// DeleteData 刪除資料
func (base *BaseTable) DeleteData(id string, userID, activityID, gameID string) error {
	if base.Info.DeleteFunc != nil {
		if err := base.Info.DeleteFunc(strings.Split(id, ","), userID, activityID, gameID); err != nil {
			return err
		}
	}
	return nil
}

// sql 取得db.SQL
func (base *BaseTable) sql(services service.List) *db.SQL {
	return db.Conn(base.db(services))
}

// getColumns 取得所有欄位
func (base *BaseTable) getColumns(table string, services service.List) ([]string, bool) {
	var auto = false
	columnsModel, _ := base.sql(services).Table(table).ShowColumns()
	// fmt.Println("所有欄位參數columnsModel: ", columnsModel)
	columns := make([]string, len(columnsModel))
	for key, model := range columnsModel {
		columns[key] = model["Field"].(string)
		if columns[key] == base.primaryKey.Name {
			if v, ok := model["Extra"].(string); ok {
				if v == "auto_increment" {
					auto = true
				}
			}
		}
	}

	return columns, auto
}

// getDataFromDatabase 取得頁面所有資料
func (base *BaseTable) getDataFromDatabase(params parameter.Parameters, services service.List) (PanelInfo, error) {
	var (
		connection  = base.db(services)
		delimiter   = connection.GetDelimiter()
		delimiter2  = connection.GetDelimiter2()
		placeholder = utils.Delimiter(delimiter, delimiter2, "%s")
		// countStatement = "select count(*) from %s %s %s"
		queryStatement string
		queryCmd       string
		table          = utils.Delimiter(delimiter, delimiter2, base.Info.Table)
		primaryKey     = table + "." + utils.Delimiter(delimiter, delimiter2, base.primaryKey.Name)
		wheres         = ""
		args           = make([]interface{}, 0)
		whereArgs      = make([]interface{}, 0)
		// size           int
	)

	if len(params.SortField) == 1 {
		// %s means: fields, table, join table, wheres, group by, order by field, order by type
		queryStatement = "select %s from " + placeholder + "%s %s %s order by " +
			placeholder + "." + placeholder + " %s LIMIT ? OFFSET ?"
	} else if len(params.SortField) == 2 {
		// %s means: fields, table, join table, wheres, group by, order by field, order by type
		queryStatement = "select %s from " + placeholder + "%s %s %s order by " +
			placeholder + "." + placeholder + " %s, " + placeholder + "." + placeholder + " %s LIMIT ? OFFSET ?"
	}

	columns, _ := base.getColumns(base.Info.Table, services)
	fieldList, fields, joinFields, joins, joinTables := base.GetFieldList(params, columns)

	fields += primaryKey
	allFields := fields
	if joinFields != "" {
		allFields += "," + joinFields[:len(joinFields)-1]
	}

	// parameter
	wheres, whereArgs = params.Statement(wheres, base.Info.Table,
		connection.GetDelimiter(), connection.GetDelimiter2(), whereArgs, columns)
	if wheres != "" {
		wheres = "where " + wheres
	}

	if connection.Name() == "mysql" {
		pageSizeInt, _ := strconv.Atoi(params.PageSize)
		pageInt, _ := strconv.Atoi(params.Page)
		args = append(whereArgs, pageSizeInt, (pageInt-1)*(pageSizeInt))
	}

	groupBy := ""
	if len(joinTables) > 0 {
		if connection.Name() == "mysql" {
			groupBy = " GROUP BY " + primaryKey
		}
	}

	if len(params.SortField) == 1 {
		queryCmd = fmt.Sprintf(queryStatement, allFields, base.Info.Table, joins, wheres, groupBy,
			base.Info.Table, params.SortField[0], params.SortType[0])
	} else if len(params.SortField) == 2 {
		queryCmd = fmt.Sprintf(queryStatement, allFields, base.Info.Table, joins, wheres, groupBy,
			base.Info.Table, params.SortField[0], params.SortType[0],
			base.Info.Table, params.SortField[1], params.SortType[1])
	}

	res, err := connection.QueryWithConnection(base.connName, queryCmd, args...)
	if err != nil {
		return PanelInfo{}, err
	}

	infoList := make([]map[string]types.InfoItem, 0)
	for i := 0; i < len(res); i++ {
		infoList = append(infoList, base.getTemplateDataModel(res[i], params, columns))
	}

	// countCmd := fmt.Sprintf(countStatement, base.Info.Table, joins, wheres)
	// total, err := connection.QueryWithConnection(base.connName, countCmd, whereArgs...)
	// if err != nil {
	// 	return PanelInfo{}, err
	// }
	// if base.driver == "mysql" {
	// 	size = int(total[0]["count(*)"].(int64))
	// }

	// paginator := paginator.GetPaginatorInformation(size, params)
	// paginator.PageSizeList = base.Info.GetPageSizeList()
	// paginator.Option = make(map[string]template.HTML, len(paginator.PageSizeList))
	// for i := 0; i < len(paginator.PageSizeList); i++ {
	// 	paginator.Option[paginator.PageSizeList[i]] = template.HTML("")
	// }
	// paginator.Option[params.PageSize] = template.HTML("select")

	return PanelInfo{
		InfoList:  infoList,
		FieldList: fieldList,
		// Paginator: paginator,
	}, nil
}

// getTemplateDataModel 處理每一筆資料
func (base *BaseTable) getTemplateDataModel(res map[string]interface{}, params parameter.Parameters,
	columns []string) map[string]types.InfoItem {
	var (
		templateDataModel = make(map[string]types.InfoItem)
		headField         = ""
		primaryKeyValue   = db.GetValueFromDatabaseType(base.primaryKey.Type, res[base.primaryKey.Name])
	)

	for _, field := range base.Info.FieldList {
		headField = field.Field

		if field.Joins.Valid() {
			headField = field.Joins.Last().JoinTable + "_join_" + field.Field
		}
		// if field.IsHide {
		// 	continue
		// }
		if !utils.InArrayWithoutEmpty(params.Columns, headField) {
			continue
		}

		typeName := field.TypeName
		if field.Joins.Valid() {
			typeName = db.Varchar
		}

		combineValue := db.GetValueFromDatabaseType(typeName, res[headField]).String()
		var value interface{}
		if len(columns) == 0 || utils.InArray(columns, headField) || field.Joins.Valid() {
			value = field.FieldDisplay.DisplayFunc(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: combineValue,
				Row:   res,
			})
		} else {
			value = field.FieldDisplay.DisplayFunc(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: "",
				Row:   res,
			})
		}

		if valueStr, ok := value.(string); ok {
			templateDataModel[headField] = types.InfoItem{
				Content: valueStr,
				Value:   combineValue,
			}
		} else {
			templateDataModel[headField] = types.InfoItem{
				Content: value,
				Value:   combineValue,
			}
		}
	}

	// if len(others) > 0 {
	// 	for i := 0; i < len(others); i++ {
	// 		field := base.Info.FieldList.GetField(others[i])
	// 		fieldValue := db.GetValueFromDatabaseType(field.TypeName, res[field.Field])

	// 		value := field.FieldDisplay.DisplayFunc(types.FieldModel{
	// 			ID:    fieldValue.String(),
	// 			Value: fieldValue.String(),
	// 			Row:   res,
	// 		})
	// 		if valueStr, ok := value.(string); ok {
	// 			templateDataModel[field.Field] = types.InfoItem{
	// 				Content: valueStr,
	// 				Value:   fieldValue.String(),
	// 			}
	// 		} else {
	// 			templateDataModel[field.Field] = types.InfoItem{
	// 				Content: value,
	// 				Value:   primaryKeyValue.String(),
	// 			}
	// 		}
	// 	}
	// }
	return templateDataModel
}

// GetFieldList 取得欄位資訊、join語法...等資訊
func (base *BaseTable) GetFieldList(params parameter.Parameters, columns []string) (types.FieldList,
	string, string, string, []string) {
	return base.Info.FieldList.GetFieldList(types.TableInfo{
		Table:      base.Info.Table,
		Delimiter:  base.getDelimiter(),
		Delimiter2: base.getDelimiter2(),
		Driver:     base.driver,
		PrimaryKey: base.primaryKey.Name,
	}, params, columns)
}

// db 取得connection
func (base *BaseTable) db(services service.List) db.Connection {
	if base.connection == nil {
		base.connection = db.GetConnectionFromService(services.Get(base.driver))
	}
	return base.connection
}

// getDelimiter 取得分隔符號
func (base *BaseTable) getDelimiter() string {
	if base.driver == "mysql" {
		return "`"
	}
	return ""
}

// getDelimiter2 取得分隔符號
func (base *BaseTable) getDelimiter2() string {
	if base.driver == "mysql" {
		return "`"
	}
	return ""
}

// interfaces []string轉換成[]interface
func interfaces(arr []string) []interface{} {
	var iarr = make([]interface{}, len(arr))
	for key, v := range arr {
		iarr[key] = v
	}
	return iarr
}
