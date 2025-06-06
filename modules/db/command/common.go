package command

import "fmt"

type commonDialect struct {
	delimiter  string
	delimiter2 string
}

func (c commonDialect) Insert(comp *Filter) string {
	comp.prepareInsert(c.delimiter, c.delimiter2)
	return comp.Statement
}

func (c commonDialect) Inserts(comp *Filter, dataAmount int, fields string, values [][]string) error {
	err := comp.prepareInserts(dataAmount, fields, values)
	return err
}

func (c commonDialect) Delete(comp *Filter) string {
	comp.Statement = "delete from " + c.WrapTableName(comp) + comp.getWheres(c.delimiter, c.delimiter2)
	return comp.Statement
}

func (c commonDialect) Update(comp *Filter) string {
	comp.prepareUpdate(c.delimiter, c.delimiter2)
	return comp.Statement
}

func (c commonDialect) Select(comp *Filter) string {
	comp.Statement = "select " + comp.getFields(c.delimiter, c.delimiter2) + " from " + c.WrapTableName(comp) +
		comp.getJoins(c.delimiter, c.delimiter2) + comp.getWheres(c.delimiter, c.delimiter2) +
		comp.getGroupBy() + comp.getOrderBy() + comp.getLimit() + comp.getOffset()
	return comp.Statement
}

func (c commonDialect) ShowColumns(table string) string {
	return fmt.Sprintf("select * from information_schema.columns where table_name = '%s'", table)
}

func (c commonDialect) GetName() string {
	return "common"
}

func (c commonDialect) WrapTableName(comp *Filter) string {
	return c.delimiter + comp.TableName + c.delimiter2
}

func (c commonDialect) ShowTables() string {
	return "show tables"
}

func (c commonDialect) GetDelimiter() string {
	return c.delimiter
}

func (c commonDialect) GetDelimiter2() string {
	return c.delimiter2
}

func (c commonDialect) GetDelimiters() []string {
	return []string{c.delimiter, c.delimiter2}
}
