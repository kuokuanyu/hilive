package db

import "fmt"

// DatabaseType string
type DatabaseType string

// Value string
type Value string

var (
	StringTypeList = []DatabaseType{Date, Time, Year, Datetime, Timestamptz, Timestamp, Timetz,
		Varchar, Char, Mediumtext, Longtext, Tinytext,
		Text, JSON, Blob, Tinyblob, Mediumblob, Longblob,
		Interval, Point, Bpchar,
		Line, Lseg, Box, Path, Polygon, Circle, Cidr, Inet, Macaddr, Character, Varyingcharacter,
		Nchar, Nativecharacter, Nvarchar, Clob, Binary, Varbinary, Enum, Set, Geometry, Multilinestring,
		Multipolygon, Linestring, Multipoint, Geometrycollection, Name, UUID, Timestamptz,
		Name, UUID, Inet}

	BoolTypeList = []DatabaseType{Bool, Boolean}
	IntTypeList  = []DatabaseType{Int4, Int2, Int8,
		Int,
		Tinyint,
		Mediumint,
		Smallint,
		Smallserial, Serial, Bigserial,
		Integer,
		Bigint}
	FloatTypeList = []DatabaseType{Float, Float4, Float8, Double, Real, Doubleprecision}
	UintTypeList  = []DatabaseType{Decimal, Bit, Money, Numeric}
)

const (
	Int                DatabaseType = "INT"
	Tinyint            DatabaseType = "TINYINT"
	Mediumint          DatabaseType = "MEDIUMINT"
	Smallint           DatabaseType = "SMALLINT"
	Bigint             DatabaseType = "BIGINT"
	Bit                DatabaseType = "BIT"
	Int8               DatabaseType = "INT8"
	Int4               DatabaseType = "INT4"
	Int2               DatabaseType = "INT2"
	Integer            DatabaseType = "INTEGER"
	Numeric            DatabaseType = "NUMERIC"
	Smallserial        DatabaseType = "SMALLSERIAL"
	Serial             DatabaseType = "SERIAL"
	Bigserial          DatabaseType = "BIGSERIAL"
	Money              DatabaseType = "MONEY"
	Real               DatabaseType = "REAL"
	Float              DatabaseType = "FLOAT"
	Float4             DatabaseType = "FLOAT4"
	Float8             DatabaseType = "FLOAT8"
	Double             DatabaseType = "DOUBLE"
	Decimal            DatabaseType = "DECIMAL"
	Doubleprecision    DatabaseType = "DOUBLEPRECISION"
	Date               DatabaseType = "DATE"
	Time               DatabaseType = "TIME"
	Year               DatabaseType = "YEAR"
	Datetime           DatabaseType = "DATETIME"
	Timestamp          DatabaseType = "TIMESTAMP"
	Text               DatabaseType = "TEXT"
	Longtext           DatabaseType = "LONGTEXT"
	Mediumtext         DatabaseType = "MEDIUMTEXT"
	Tinytext           DatabaseType = "TINYTEXT"
	Varchar            DatabaseType = "VARCHAR"
	Char               DatabaseType = "CHAR"
	Bpchar             DatabaseType = "BPCHAR"
	JSON               DatabaseType = "JSON"
	Blob               DatabaseType = "BLOB"
	Tinyblob           DatabaseType = "TINYBLOB"
	Mediumblob         DatabaseType = "MEDIUMBLOB"
	Longblob           DatabaseType = "LONGBLOB"
	Interval           DatabaseType = "INTERVAL"
	Boolean            DatabaseType = "BOOLEAN"
	Bool               DatabaseType = "BOOL"
	Point              DatabaseType = "POINT"
	Line               DatabaseType = "LINE"
	Lseg               DatabaseType = "LSEG"
	Box                DatabaseType = "BOX"
	Path               DatabaseType = "PATH"
	Polygon            DatabaseType = "POLYGON"
	Circle             DatabaseType = "CIRCLE"
	Cidr               DatabaseType = "CIDR"
	Inet               DatabaseType = "INET"
	Macaddr            DatabaseType = "MACADDR"
	Character          DatabaseType = "CHARACTER"
	Varyingcharacter   DatabaseType = "VARYINGCHARACTER"
	Nchar              DatabaseType = "NCHAR"
	Nativecharacter    DatabaseType = "NATIVECHARACTER"
	Nvarchar           DatabaseType = "NVARCHAR"
	Clob               DatabaseType = "CLOB"
	Binary             DatabaseType = "BINARY"
	Varbinary          DatabaseType = "VARBINARY"
	Enum               DatabaseType = "ENUM"
	Set                DatabaseType = "SET"
	Geometry           DatabaseType = "GEOMETRY"
	Multilinestring    DatabaseType = "MULTILINESTRING"
	Multipolygon       DatabaseType = "MULTIPOLYGON"
	Linestring         DatabaseType = "LINESTRING"
	Multipoint         DatabaseType = "MULTIPOINT"
	Geometrycollection DatabaseType = "GEOMETRYCOLLECTION"
	Name               DatabaseType = "NAME"
	UUID               DatabaseType = "UUID"
	Timestamptz        DatabaseType = "TIMESTAMPTZ"
	Timetz             DatabaseType = "TIMETZ"
)

// DT string to DatavaseType
func DT(s string) DatabaseType {
	return DatabaseType(s)
}

// Contains 是否包含
func Contains(v DatabaseType, a []DatabaseType) bool {
	for _, i := range a {
		if v == i {
			return true
		}
	}
	return false
}

// GetValueFromDatabaseType 藉由欄位類型取得data值
func GetValueFromDatabaseType(typ DatabaseType, value interface{}) Value {
	switch {
	case Contains(typ, StringTypeList):
		if v, ok := value.(string); ok {
			return Value(v)
		}
		return ""
	case Contains(typ, BoolTypeList):
		if v, ok := value.(bool); ok {
			if v {
				return "true"
			}
			return "false"
		}
		if v, ok := value.(int64); ok {
			if v == 0 {
				return "false"
			}
			return "true"
		}
		return "false"
	case Contains(typ, IntTypeList):
		if v, ok := value.(int64); ok {
			return Value(fmt.Sprintf("%d", v))
		}
		return "0"
	case Contains(typ, FloatTypeList):
		if v, ok := value.(float64); ok {
			return Value(fmt.Sprintf("%f", v))
		}
		return "0"
	case Contains(typ, UintTypeList):
		if v, ok := value.([]uint8); ok {
			return Value(string(v))
		}
		return "0"
	}
	panic("錯誤databasetype?" + string(typ))
}

// String return the string value.
func (v Value) String() string {
	return string(v)
}
