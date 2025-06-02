package command

import (
	"errors"
	"fmt"
	"hilive/modules/utils"
	"log"
	"strings"
)

// Where where條件
type Where struct {
	Operation   string
	Field       string
	Value       string
	Combination string
}

// Join join條件
type Join struct {
	Table           string
	FieldA          string
	FieldA1         string
	FieldB          string
	FieldB1         string
	FieldC          string
	FieldC1         string
	Operation       string
	UnionAliasTable string // union其他表後的表名
	UnionStatement  string // 先union其他表之後才left join
}

// Value 欄位: 數值
type Value map[string]interface{}

// Filter 過濾條件
type Filter struct {
	Fields    []string      // 欄位
	TableName string        // 資料表
	Wheres    []Where       // where 條件
	WhereRaws string        // where條件
	Leftjoins []Join        // join條件
	Args      []interface{} // 篩選條件值
	Order     string        // 排序
	Offset    string        // 從第幾筆開始
	Limit     string        // 筆數
	Group     string        // 分群
	Statement string        // 資料庫命令
	Values    Value         // insert、update值
}

// Command 資料庫CRUD方法
type Command interface {
	GetName() string                                                                   // 取得引擎名稱
	ShowColumns(table string) string                                                   // 取得所有欄位
	ShowTables() string                                                                // 回傳所有資料表
	Insert(condition *Filter) string                                                   // 處理新增資料命令
	Delete(condition *Filter) string                                                   // 處理刪除資料命令
	Update(condition *Filter) string                                                   // 處理更新資料命令
	Select(condition *Filter) string                                                   // 處理查詢資料命令
	Inserts(condition *Filter, dataAmount int, fields string, values [][]string) error // 處理匹量新增資料命令
}

// GetCommand 取得Command
func GetCommand(driver string) Command {
	switch driver {
	case "mysql":
		return mysql{
			commonDialect: commonDialect{delimiter: "`", delimiter2: "`"},
		}
	default:
		return commonDialect{delimiter: "`", delimiter2: "`"}
	}
}

func (f *Filter) prepareInserts(dataAmount int, fields string, values [][]string) error {
	if dataAmount > 0 && len(values) > 0 {
		// ex: field = column1,column2,column3,column4
		// 解析欄位並去除多餘空白
		fieldList := strings.Split(fields, ",")
		for i := range fieldList {
			fieldList[i] = strings.TrimSpace(fieldList[i])
		}

		// 比對value參數陣列長度是否跟欄位數量相等
		// if len(values) != len(fieldList) {
		// 	log.Println("錯誤: 欄位數量與值數量不一致")
		// 	return errors.New("錯誤: 欄位數量與值數量不一致")
		// }

		// 比對value陣列裡的陣列資料長度是否跟dataAmount數量相等
		for i := 0; i < len(values); i++ {
			if len(values[i]) != dataAmount {
				log.Println("錯誤: 每個值陣列的長度必須等於資料數量")
				return errors.New("錯誤: 每個值陣列的長度必須等於資料數量")
			}
		}

		// 構建 SQL 語法
		placeholders := "(" + strings.Repeat("?,", len(fieldList))
		placeholders = placeholders[:len(placeholders)-1] + ")" // ex: (?,?,?,?)
		f.Statement = "INSERT INTO " + f.TableName + " (" + strings.Join(fieldList, ", ") + ") VALUES "

		// 多筆資料處理
		var args []interface{}
		placeholderList := make([]string, dataAmount)
		for i := 0; i < dataAmount; i++ {
			placeholderList[i] = placeholders

			// 插入資料
			for j := 0; j < len(fieldList); j++ {
				args = append(args, values[j][i])
			}
		}

		// 最終語法
		f.Statement += strings.Join(placeholderList, ",")
		f.Args = args
	}
	return nil
}

func (f *Filter) prepareInsert(delimiter, delimiter2 string) {
	fields := " ("
	quesMark := "("

	for key, value := range f.Values {
		fields += wrap(delimiter, delimiter2, key) + ","
		quesMark += "?,"
		f.Args = append(f.Args, value)
	}
	fields = fields[:len(fields)-1] + ")"
	quesMark = quesMark[:len(quesMark)-1] + ")"

	f.Statement = "insert into " + delimiter + f.TableName + delimiter2 +
		fields + " values " + quesMark
}

func (f *Filter) prepareUpdate(delimiter, delimiter2 string) {
	fields := ""
	args := make([]interface{}, 0)

	if len(f.Values) != 0 {
		for key, value := range f.Values {
			if strings.Contains(utils.GetString(value, ""), "+ 1") {
				fields += fmt.Sprintf("%s = %s + 1, ", wrap(delimiter, delimiter2, key), key)
			} else if strings.Contains(utils.GetString(value, ""), "- 1") {
				fields += fmt.Sprintf("%s = %s - 1, ", wrap(delimiter, delimiter2, key), key)
			} else if key == "prize_amount" {
				// 減少獎品數量，更新獎品總數時避免map格式發生順序不一的情況導致更新數值發生問題
				// 先更新剩餘獎品再更新獎品總數
				// ex: prize_remain = prire_remain + model.PrizeAmount - prize_amount, `prize_amount` = model.PrizeAmount
				fields += fmt.Sprintf("`prize_remain` = `prize_remain` + %s - `prize_amount`, `prize_amount` = %s, ",
					value, value)
			} else if key == "option_score" {
				// 投票選項分數遞增
				// ex: option_score = option_score + ?
				fields += fmt.Sprintf("`option_score` = %s, ", value)
			} else {
				fields += wrap(delimiter, delimiter2, key) + " = ?, "
				args = append(args, value)
			}

			// if value == "game_round + 1" {
			// 	// 增加遊戲輪次
			// 	fields += wrap(delimiter, delimiter2, key) + " = game_round + 1, "
			// 	// args = append(args, value)
			// } else if value == "qa_round + 1" {
			// 	// 遞增快問快答題數
			// 	fields += wrap(delimiter, delimiter2, key) + " = qa_round + 1, "
			// } else if value == "attend + 1" {
			// 	// 遞增活動人數
			// 	fields += wrap(delimiter, delimiter2, key) + " = attend + 1, "
			// } else if value == "attend - 1" {
			// 	// 遞減活動人數
			// 	fields += wrap(delimiter, delimiter2, key) + " = attend - 1, "
			// } else if value == "game_attend + 1" {
			// 	// 遞增遊戲人數
			// 	fields += wrap(delimiter, delimiter2, key) + " = game_attend + 1, "
			// } else if value == "prize_remain - 1" {
			// 	// 遞減獎品數量
			// 	fields += wrap(delimiter, delimiter2, key) + " = prize_remain - 1, "
			// } else if key == "prize_amount" {
			// 	// 減少獎品數量，更新獎品總數時避免map格式發生順序不一的情況導致更新數值發生問題
			// 	// 先更新剩餘獎品再更新獎品總數
			// 	// ex: prize_remain = prire_remain + model.PrizeAmount - prize_amount, `prize_amount` = model.PrizeAmount
			// 	fields += fmt.Sprintf("`prize_remain` = `prize_remain` + %s - `prize_amount`, `prize_amount` = %s, ",
			// 		value, value)
			// } else if value == "likes + 1" {
			// 	// 遞增按讚數量
			// 	fields += wrap(delimiter, delimiter2, key) + " = likes + 1, "
			// } else if value == "likes - 1" {
			// 	// 遞減按讚數量
			// 	fields += wrap(delimiter, delimiter2, key) + " = likes - 1, "
			// } else if value == "number + 1" {
			// 	// 遞增抽號碼資料
			// 	fields += wrap(delimiter, delimiter2, key) + " = number + 1, "
			// 	// args = append(args, strings.Split(fmt.Sprintf("%s", value), " + ")[1]) // value = attend + n
			// } else {
			// 	fields += wrap(delimiter, delimiter2, key) + " = ?, "
			// 	args = append(args, value)
			// }

			// else if strings.Contains(fmt.Sprintf("%s", value), "prize_remain") &&
			// 	strings.Contains(fmt.Sprintf("%s", value), "prize_amount") { // 減少獎品數量(更新獎品總數時)
			// 	// ex: prize_remain = prire_remain + model.PrizeAmount - prize_amount
			// 	fields += wrap(delimiter, delimiter2, key) + " = " + fmt.Sprintf("%s", value) + ", "
			// } else if value == "likes + 1" { // 遞增按讚數量
			// 	fields += wrap(delimiter, delimiter2, key) + " = likes + 1, "
			// } else if value == "likes - 1" { // 遞減按讚數量
			// 	fields += wrap(delimiter, delimiter2, key) + " = likes - 1, "
			// } else if value == "number + 1" { // 遞增抽號碼資料
			// 	fields += wrap(delimiter, delimiter2, key) + " = number + 1, "
			// 	// args = append(args, strings.Split(fmt.Sprintf("%s", value), " + ")[1]) // value = attend + n
			// } else {
			// 	fields += wrap(delimiter, delimiter2, key) + " = ?, "
			// 	args = append(args, value)
			// }
		}

		fields = fields[:len(fields)-2]
		f.Args = append(args, f.Args...)
	} else {
		panic("執行prepareUpdate函式更新資料發生錯誤: 更新必須設置參數")
	}

	f.Statement = "update " + delimiter + f.TableName + delimiter2 + f.getJoins(delimiter, delimiter2) +
		" set " + fields + f.getWheres(delimiter, delimiter2)
}

// 處理limit
func (f *Filter) getLimit() string {
	if f.Limit == "" {
		return ""
	}
	return " limit " + f.Limit + " "
}

// 處理offset
func (f *Filter) getOffset() string {
	if f.Offset == "" {
		return ""
	}
	return " offset " + f.Offset + " "
}

// 處理排序
func (f *Filter) getOrderBy() string {
	if f.Order == "" {
		return ""
	}
	return " order by " + f.Order + " "
}

// 處理分群
func (f *Filter) getGroupBy() string {
	if f.Group == "" {
		return ""
	}
	return " group by " + f.Group + " "
}

// 處理join
func (f *Filter) getJoins(delimiter, delimiter2 string) string {
	if len(f.Leftjoins) == 0 {
		return ""
	}
	joins := ""
	for _, join := range f.Leftjoins {
		if join.UnionAliasTable != "" && join.UnionStatement != "" {
			// ex: left join(SELECT `game_id`,`title` FROM `activity_redpack`
			// union all
			// SELECT `game_id`,`title` FROM `activity_ropepack`
			// union all
			// SELECT `game_id`,`title` FROM `activity_lottery`
			// union all
			// SELECT `game_id`,`title` FROM `activity_whack_mole`) as `game`
			joins += " left join (" + join.UnionStatement + ") as " +
				wrap(delimiter, delimiter2, join.UnionAliasTable)
		} else {
			// left join table
			joins += " left join " + wrap(delimiter, delimiter2, join.Table)
		}

		joins += " on (" +
			f.processLeftJoinField(join.FieldA, delimiter, delimiter2) + " " +
			join.Operation + " " +
			f.processLeftJoinField(join.FieldA1, delimiter, delimiter2)

		if join.FieldB != "" && join.FieldB1 != "" {
			joins += " and " +
				f.processLeftJoinField(join.FieldB, delimiter, delimiter2) + " " +
				join.Operation + " " +
				f.processLeftJoinField(join.FieldB1, delimiter, delimiter2)
		} else {
			joins += ") "
			continue
		}

		if join.FieldC != "" && join.FieldC1 != "" {
			joins += " and " +
				f.processLeftJoinField(join.FieldC, delimiter, delimiter2) + " " +
				join.Operation + " " +
				f.processLeftJoinField(join.FieldC1, delimiter, delimiter2) + ") "
		} else {
			joins += ") "
			continue
		}
	}

	return joins
}

// 處理欄位
func (f *Filter) getFields(delimiter, delimiter2 string) string {
	if len(f.Fields) == 0 {
		return "*"
	}

	if f.Fields[0] == "count(*)" {
		// 查詢某個條件下的資料筆數
		// ex: select count(*) from activity
		return "count(*)"
	}

	if strings.Contains(f.Fields[0], "sum") {
		// 查詢某個欄位的和
		// ex: sum(likes)
		return f.Fields[0]
	}

	fields := ""
	if len(f.Leftjoins) == 0 {
		// 無join其他表
		for _, field := range f.Fields {
			if strings.Contains(field, "SUM") || strings.Contains(field, "COUNT(*)") {
				// 計算不同種類的資料數量
				// ex: COUNT(*) AS all_amount
				// ex: SUM(CASE WHEN message_status = 'yes' THEN 1 ELSE 0 END) AS yes_message_amount
				fields += field + ","
			} else {
				arr := strings.Split(field, ".")
				if len(arr) > 1 {
					// ex: `activity_game`.`status` or `activity_game`.`status` as `status1`
					if strings.Contains(arr[1], " as ") {
						newArr := strings.Split(arr[1], " as ")
						fields += wrap(delimiter, delimiter2, arr[0]) + "." +
							wrap(delimiter, delimiter2, newArr[0]) + " as " +
							wrap(delimiter, delimiter2, newArr[1]) + ","
					} else {
						fields += wrap(delimiter, delimiter2, arr[0]) + "." +
							wrap(delimiter, delimiter2, arr[1]) + ","
					}
				} else {
					fields += wrap(delimiter, delimiter2, field) + ","
				}

				// 一般欄位，處理後變成: `欄位名稱`
				// fields += wrap(delimiter, delimiter2, field) + ","
			}
		}
	} else {
		// 有join其他表
		for _, field := range f.Fields {
			arr := strings.Split(field, ".")
			if len(arr) > 1 {
				// ex: `activity_game`.`status` or `activity_game`.`status` as `status1`
				if strings.Contains(arr[1], " as ") {
					newArr := strings.Split(arr[1], " as ")
					fields += wrap(delimiter, delimiter2, arr[0]) + "." +
						wrap(delimiter, delimiter2, newArr[0]) + " as " +
						wrap(delimiter, delimiter2, newArr[1]) + ","
				} else {
					fields += wrap(delimiter, delimiter2, arr[0]) + "." +
						wrap(delimiter, delimiter2, arr[1]) + ","
				}
			} else {
				fields += wrap(delimiter, delimiter2, field) + ","
			}
		}
	}

	return fields[:len(fields)-1]
}

// 處理where
func (f *Filter) getWheres(delimiter, delimiter2 string) string {
	if len(f.Wheres) == 0 {
		if f.WhereRaws != "" {
			return " where " + f.WhereRaws
		}
		return ""
	}
	wheres := " where "
	var arr []string

	// for _, where := range f.Wheres {
	// 	arr = strings.Split(where.Field, ".")
	// 	if len(arr) > 1 {
	// 		wheres += arr[0] + "." + wrap(delimiter, delimiter2, arr[1]) + " " + where.Operation + " " + where.Value + " and "
	// 	} else {
	// 		wheres += wrap(delimiter, delimiter2, where.Field) + " " + where.Operation + " " + where.Value + " and "
	// 	}
	// }
	// if f.WhereRaws != "" {
	// 	return wheres + f.WhereRaws
	// }

	for i := 0; i < len(f.Wheres); i++ {
		arr = strings.Split(f.Wheres[i].Field, ".")
		if len(arr) > 1 {
			wheres += arr[0] + "." + wrap(delimiter, delimiter2, arr[1]) + " " +
				f.Wheres[i].Operation + " " + f.Wheres[i].Value + " "
		} else {
			wheres += wrap(delimiter, delimiter2, f.Wheres[i].Field) + " " +
				f.Wheres[i].Operation + " " + f.Wheres[i].Value + " "
		}
		if i < len(f.Wheres)-1 {
			wheres += f.Wheres[i].Combination + " "
		}
	}
	if f.WhereRaws != "" {
		return wheres + f.WhereRaws
	}
	return wheres[:len(wheres)-1]
}

// wrap 欄位名稱加上分隔符號
func wrap(delimiter, delimiter2, field string) string {
	if field == "*" {
		return "*"
	}

	// 判斷欄位是否為tablename.field
	arr := strings.Split(field, ".")
	if len(arr) == 2 {
		return delimiter + arr[0] + delimiter2 + "." + delimiter + arr[1] + delimiter2
	}
	return delimiter + field + delimiter2
}

// processLeftJoinField 處理left join語法
func (f *Filter) processLeftJoinField(field, delimiter, delimiter2 string) string {
	arr := strings.Split(field, ".")
	if len(arr) == 2 {
		return delimiter + arr[0] + delimiter2 + "." + delimiter + arr[1] + delimiter2
	}
	return delimiter + field + delimiter2
	// return field
}
