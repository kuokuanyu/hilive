package parameter

import (
	"net/url"
	"strconv"
	"strings"

	"hilive/modules/utils"
)

var keys = []string{"__page", "__pageSize", "__sort", "__columns", "id", "__sort_type"}

// 運算符號
var operators = map[string]string{
	"like": "like",
	"gr":   ">",
	"gq":   ">=",
	"eq":   "=",
	"ne":   "!=",
	"le":   "<",
	"lq":   "<=",
}

// Parameters 頁面資訊
type Parameters struct {
	Page      string              // 第n頁
	PageSize  string              // 顯示資料數
	SortField []string            // 排序欄位
	SortType  []string            // 排序欄位類別
	Columns   []string            // 欄位資訊
	Fields    map[string][]string // url儲存資訊
	Path      string              // 路徑(不含篩選參數)
}

// DefaultParameters 預設Parameters
func DefaultParameters() Parameters {
	return Parameters{
		Page:     "1",
		PageSize: "10",
		Fields:   make(map[string][]string),
	}
}

// SetPKs 設置Parameters.Fields["__pk"]值
func (param Parameters) SetPKs(id ...string) Parameters {
	param.Fields["__pk"] = []string{strings.Join(id, ",")}
	return param
}

// SetField 設置Parameters.Fields[field]
func (param Parameters) SetField(field string, value ...string) Parameters {
	param.Fields[field] = value
	return param
}

// SetPage 設置第n頁
func (param Parameters) SetPage(page string) Parameters {
	param.Page = page
	return param
}

// FindPK 取得pk值
func (param Parameters) FindPK() string {
	value, ok := param.Fields["__pk"]
	if ok && len(value) > 0 {
		return strings.Split(value[0], ",")[0]
	}
	return ""
}

// FindPKs 取得pk值(多個)
func (param Parameters) FindPKs() []string {
	value, ok := param.Fields["__pk"]
	if ok && len(value) > 0 {
		return strings.Split(value[0], ",")
	}
	return []string{}
}

// GetParamFromURL 解析URL並設置Parameters
func GetParamFromURL(urlStr string, defaultPageSize int, sortField []string, sortType []string) Parameters {
	u, err := url.Parse(urlStr)
	if err != nil {
		return DefaultParameters()
	}
	return GetParam(u, defaultPageSize, sortField, sortType)
}

// GetParam 解析URL並設置Parameters
func GetParam(u *url.URL, defaultPageSize int, sortField, sortType []string) Parameters {
	values := u.Query()

	page := getDefault(values, "__page", "1")
	pageSize := getDefault(values, "__pageSize", strconv.Itoa(defaultPageSize))
	columns := getDefault(values, "__columns", "")
	columnsArr := make([]string, 0)
	if columns != "" {
		columns, _ = url.QueryUnescape(columns)
		columnsArr = strings.Split(columns, ",")
	}

	fields := make(map[string][]string)
	for key, value := range values {
		if !utils.InArray(keys, key) && len(value) > 0 && value[0] != "" {
			fields[strings.Replace(key, "[]", "", -1)] = value
		}
	}

	return Parameters{
		Page:      page,
		PageSize:  pageSize,
		Path:      u.Path,
		SortField: sortField,
		SortType:  sortType,
		Fields:    fields,
		Columns:   columnsArr,
	}
}

// Statement 處理where語法
func (param Parameters) Statement(wheres, table, delimiter, delimiter2 string,
	whereArgs []interface{}, columns []string) (string, []interface{}) {
	for key, value := range param.Fields {
		var op string
		if strings.Contains(key, "_end") {
			key = strings.Replace(key, "_end", "", -1)
			op = "<="
		} else if strings.Contains(key, "_start") {
			key = strings.Replace(key, "_start", "", -1)
			op = ">="
		} else if len(value) > 1 {
			op = "in"
		} else if !strings.Contains(key, "__operator__") {
			op = "="
		}

		if utils.InArray(columns, key) {
			if op == "in" {
				qmark := ""
				for range value {
					qmark += "?,"
				}
				wheres += table + "." + utils.FilterField(key, delimiter, delimiter2) + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
			} else {
				wheres += table + "." + utils.FilterField(key, delimiter, delimiter2) + " " + op + " ? and "
			}

			if op == "like" && !strings.Contains(value[0], "%") {
				whereArgs = append(whereArgs, "%"+value[0]+"%")
			} else {
				for _, v := range value {
					whereArgs = append(whereArgs, v)
				}
			}
		} else {
			keys := strings.Split(key, "_join_")
			if len(keys) > 1 {
				if op == "in" {
					qmark := ""
					for range value {
						qmark += "?,"
					}
					wheres += keys[0] + "." + keys[1] + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
				} else {
					wheres += keys[0] + "." + keys[1] + " " + op + " ? and "
				}
				if op == "like" && !strings.Contains(value[0], "%") {
					whereArgs = append(whereArgs, "%"+value[0]+"%")
				} else {
					for _, v := range value {
						whereArgs = append(whereArgs, v)
					}
				}
			}
		}
	}
	if len(wheres) > 3 {
		wheres = wheres[:len(wheres)-4]
	}
	return wheres, whereArgs
}

// getDefault 判斷值，若無則回傳預設值
func getDefault(values url.Values, key, def string) string {
	value := values.Get(key)
	if value == "" {
		return def
	}
	return value
}

// GetFieldValue 取得欄位值，若無則回傳空字串
func (param Parameters) GetFieldValue(field string) string {
	value, ok := param.Fields[field]
	if ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

// GetRoute ex: ?__page=1&__pageSize=10&__sort=id&__sort_type=desc
func (param Parameters) GetRoute() string {
	p := param.GetFixedParam()
	p.Add("__page", param.Page)
	return "?" + p.Encode()
}

// GetFixedParam 設置sort、page、field相關資訊
func (param Parameters) GetFixedParam() url.Values {
	p := url.Values{}
	p.Add("__sort", strings.Join(param.SortField, ","))
	p.Add("__pageSize", param.PageSize)
	p.Add("__sort_type", strings.Join(param.SortType, ","))
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return p
}

// GetLastPageRoute 取得上一頁路徑
func (param Parameters) GetLastPageRoute() string {
	p := param.GetFixedParam()
	pageInt, _ := strconv.Atoi(param.Page)
	p.Add("__page", strconv.Itoa(pageInt-1))
	return "?" + p.Encode()
}

// GetNextPageRoute 取得下一頁路徑
func (param Parameters) GetNextPageRoute() string {
	p := param.GetFixedParam()
	pageInt, _ := strconv.Atoi(param.Page)
	p.Add("__page", strconv.Itoa(pageInt+1))
	return "?" + p.Encode()
}

// URL 第n頁的URL
func (param Parameters) URL(page string) string {
	return param.Path + param.SetPage(page).GetRoute()
}

// GetRouteWithoutPageSize 第n頁的URL(沒pagesize)
func (param Parameters) GetRouteWithoutPageSize(page string) string {
	p := url.Values{}
	p.Add("__sort", strings.Join(param.SortField, ","))
	p.Add("__page", page)
	p.Add("__sort_type", strings.Join(param.SortType, ","))
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return "?" + p.Encode()
}

// DeleteField 刪除Fields資訊
func (param Parameters) DeleteField(field string) Parameters {
	delete(param.Fields, field)
	return param
}
