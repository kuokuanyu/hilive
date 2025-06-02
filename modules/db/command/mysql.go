package command

// mysql
type mysql struct {
	commonDialect // 分隔符號
}

// GetName 取得引擎名稱
func (mysql) GetName() string {
	return "mysql"
}

// ShowColumns 取得所有欄位
func (mysql) ShowColumns(table string) string {
	return "show columns in " + table
}

// ShowTables 回傳所有資料表
func (mysql) ShowTables() string {
	return "show tables"
}
