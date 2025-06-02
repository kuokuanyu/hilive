package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/NebulousLabs/fastrand"
)

// SetInt64IfNotEmpty 當data[key]資料不為空時才會將資料寫入data(int64)
// func SetInt64IfNotEmpty(data map[string]interface{}, key string, val interface{}) {
// 	if v := GetInt64(val, -1); v >= 0 {
// 		data[key] = v
// 	}
// }

// // SetFloat64IfNotEmpty 當data[key]資料不為空時才會將資料寫入data(float64)
// func SetFloat64IfNotEmpty(data map[string]interface{}, key string, val string) {
// 	if v := GetFloat64(val, -1); v >= 0 {
// 		data[key] = v
// 	}
// }

// DetectType 判斷interface{}資料的類型
// func DetectType(i interface{}) {
// 	switch v := i.(type) {
// 	case string:
// 		fmt.Println("It's a 字串:", v)
// 	case int:
// 		fmt.Println("It's an 整數:", v)
// 	case float64:
// 		fmt.Println("It's a 小數:", v)
// 	case []string:
// 		fmt.Println("It's a slice of strings:", v)
// 	case map[string]interface{}:
// 		fmt.Println("It's a map[string]interface{}:", v)
// 	default:
// 		fmt.Println("Unknown type:", reflect.TypeOf(i))
// 	}
// }

// PrintStructFields 將struct中空值的欄位印出
// func PrintStructFields(model interface{}) {
// 	v := reflect.ValueOf(model)
// 	t := reflect.TypeOf(model)

// 	// 若是指標，先取 Elem
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 		t = t.Elem()
// 	}

// 	for i := 0; i < v.NumField(); i++ {
// 		field := t.Field(i)
// 		value := v.Field(i)

// 		if value.Interface() == "" || value.Interface() == 0 {
// 			log.Println("空值欄位名稱: ", field.Name)
// 		}

// 		// fmt.Printf("欄位名稱: %s, 值: %v\n", field.Name, value.Interface())
// 	}
// }

// ApplyInterfaceMapToStruct 將 map[string]interface{} 中的資料自動寫入任意 struct 中（限指標）
// 要求：map 的 key 必須與 struct 的json tag值相同
// struct格式裡的參數類型不一定都是string
func ApplyInterfaceMapToStruct(data map[string]interface{}, target interface{}) error {
	v := reflect.ValueOf(target)
	
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target 必須是非 nil 的指標")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target 必須指向一個 struct")
	}

	t := v.Type()

	for key, value := range data {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
			if jsonTag == "" {
				jsonTag = field.Name // 如果沒有 json tag，退而求其次用欄位名
			}

			if jsonTag == key {
				f := v.Field(i)
				if f.IsValid() && f.CanSet() {
					switch f.Kind() {
					case reflect.String:
						if str, ok := value.(string); ok {
							f.SetString(str)
						} else {
							f.SetString(fmt.Sprintf("%v", value))
						}
					case reflect.Int, reflect.Int64:
						num, err := toInt64(value)
						if err == nil {
							f.SetInt(num)
						}
					case reflect.Float64:
						num, err := toFloat64(value)
						if err == nil {
							f.SetFloat(num)
						}
					case reflect.Bool:
						if b, ok := value.(bool); ok {
							f.SetBool(b)
						}
						// 可擴充更多型別...
					}
				}
			}
		}
	}

	return nil
}



// ApplyMapToStruct 將 map[string]string 中的資料自動寫入任意 struct 中（限指標）
// 要求：map 的 key 必須與 struct 的json tag值相同
// struct格式裡的參數類型都必須是string
func ApplyMapToStruct(data map[string]string, target interface{}) error {
	v := reflect.ValueOf(target) // 取得 target 的 reflect 值（用來操作傳入的變數）

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target 必須是非 nil 的指標") // 檢查是否為非 nil 的指標，否則回傳錯誤
	}

	v = v.Elem() // 取得指標指向的實體（也就是 struct 本身）

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target 必須指向一個 struct") // 確認 target 是否為 struct，若不是則回傳錯誤
	}

	t := v.Type() // 取得 struct 的型別資訊（用來讀取欄位的定義）

	for key, value := range data { // 遍歷 map 的每個 key-value 配對
		for i := 0; i < t.NumField(); i++ { // 遍歷 struct 中的每個欄位
			field := t.Field(i)              // 取得第 i 個欄位的定義
			jsonTag := field.Tag.Get("json") // 取得該欄位的 json tag 標籤

			if strings.Split(jsonTag, ",")[0] == key { // 檢查 json tag 的主名稱是否與 map 的 key 相符
				f := v.Field(i) // 取得該欄位的實際值（可設值的對象）

				if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					// 如果欄位存在、可以設定值、且是 string 類型
					f.SetString(value) // 將 map 的字串值設到 struct 的欄位中
				}
			}
		}
	}

	return nil // 設定成功，回傳 nil 表示無錯誤
}

// StructToMap 將 struct 轉換為 map[string]interface{} 格式，使用 struct 中的 json tag 作為 key(寫入資料表時使用)
func StructToMap(data interface{}) map[string]interface{} {
	// 建立一個結果 map，用來儲存 key-value
	result := make(map[string]interface{})

	// 使用反射取得傳入資料的 Value（值）
	val := reflect.ValueOf(data)

	// 如果是指標，先取得其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 取得 struct 的 Type（類型）
	typ := val.Type()

	// 遍歷 struct 的所有欄位
	for i := 0; i < val.NumField(); i++ {
		// 取得欄位資訊（例如欄位名稱、標籤等）
		field := typ.Field(i)

		// 取得欄位上的 json tag（例如 json:"activity_id"）
		jsonTag := field.Tag.Get("json")

		var key string
		// 如果 json tag 存在且不是 "-", 就用 json tag 作為 map 的 key
		if jsonTag != "" && jsonTag != "-" {
			// 有些 tag 可能寫成 `json:"name,omitempty"`，我們只取 "name"
			key = strings.Split(jsonTag, ",")[0]
		} else {
			// 如果沒有 json tag，就使用欄位名稱轉小寫作為 key
			key = strings.ToLower(field.Name)
		}

		// 取得欄位的值，並放入 result map 中
		result[key] = val.Field(i).Interface()
	}

	// 回傳結果 map
	return result
}

// FlattenForm 將map[string][]string格式資料轉換為map[string]string(取第一個值)
func FlattenForm(values map[string][]string) map[string]string {
	flat := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			flat[k] = v[0] // 只取第一個值
		}
	}
	return flat
}

// AddUniqueStrings 將 B 中不在 A 中的元素加入 A
func AddUniqueStrings(A, B []string) []string {
	// 使用 map 儲存 A 中的元素，方便快速檢查
	existing := make(map[string]bool)
	for _, a := range A {
		existing[a] = true
	}

	// 遍歷 B，將不存在於 A 的元素加入 A
	for _, b := range B {
		if !existing[b] {
			A = append(A, b)
			existing[b] = true
		}
	}

	return A
}

// 將interface格式轉換為陣列(interface裡的資料以,為間隔)
func toStringSliceFromCommaSeparated(data interface{}) ([]string, error) {
	// 確保 data 是字串類型
	str, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("input is not a string")
	}

	// 使用 strings.Split 進行分割
	strSlice := strings.Split(str, ",")

	return strSlice, nil
}

// Interfaces []string轉換成[]interface
func Interfaces(arr []string) []interface{} {
	var iarr = make([]interface{}, len(arr))
	for key, v := range arr {
		iarr[key] = v
	}
	return iarr
}

// 加密秘鑰
var SignSecret = []byte("1234567890abcdefghijvklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 得到客户端IP地址
func ClientIP(request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

// UserSign 用戶資訊加密簽名
func UserSign(userID, secret string) string {
	str := fmt.Sprintf("user_id=%s&secret=%s", userID, secret)
	return CreateSign(str)
}

// *****舊*****
// UserSign 用戶資訊加密簽名
// func UserSign(ip, userID, secret string) string {
// 	str := fmt.Sprintf("ip=%s&user_id=%s&secret=%s", ip, userID, secret)
// 	return CreateSign(str)
// }
// *****舊*****

// CreateSign 對字串進行加密簽名
func CreateSign(str string) string {
	str = string(SignSecret) + str
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}

// Encode 加密
func Encode(password []byte) []byte {
	return []byte(base64.RawStdEncoding.EncodeToString(password))
}

// Decode 解密
func Decode(password []byte) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(string(password))
}

// SetDefault 判斷後回傳值，若無則回傳預設值
func SetDefault(value, condition, def string) string {
	if value == condition {
		return def
	}
	return value
}

// JSON 執行JSON編碼
func JSON(i interface{}) string {
	if i == nil {
		return ""
	}
	j, _ := json.Marshal(i)
	return string(j)
}

// InArray 是否在陣列中
func InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IntInArray int值是否在陣列中
func IntInArray(arr []int64, str int64) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// InArrayWithoutEmpty 是否在陣列中
func InArrayWithoutEmpty(arr []string, str string) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// AorB 判斷條件回傳值
func AorB(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}

// Random 隨機數
func Random(strings []string) ([]string, error) {
	for i := len(strings) - 1; i > 0; i-- {
		num := fastrand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := make([]string, 0)
	for i := 0; i < len(strings); i++ {
		str = append(str, strings[i])
	}
	return str, nil
}

// UUID 設置uuid
func UUID(length int64) string {
	ele := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "v", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ele, _ = Random(ele)
	uuid := ""
	var i int64
	for i = 0; i < length; i++ {
		uuid += ele[fastrand.Intn(59)]
	}
	return uuid
}

// RandomNumber 隨機產生字串(數字)
func RandomNumber(length int) string {
	numbers := "0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(result)
}

// Delimiter 分隔符號
func Delimiter(del, del2, s string) string {
	return del + s + del2
}

// FilterField 將欄位加上分個符號
func FilterField(filed, delimiter, delimiter2 string) string {
	return delimiter + filed + delimiter2
}

// GetString 將interface類型轉變成string類型(支援 int64、float64 的轉換)
func GetString(str interface{}, d string) string {
	if str == nil {
		return d
	}

	switch v := str.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v) // 預設轉換為小數點後 6 位
	}

	return fmt.Sprintf("%s", str)
}

// GetInt64 將interface類型轉變成int64類型(無法轉換時回傳預設值)
func GetInt64(i interface{}, d int64) int64 {
	if i == nil {
		return d
	}
	switch i.(type) {
	case string:
		num, err := strconv.Atoi(i.(string))
		if err != nil {
			return d
		} else {
			return int64(num)
		}
	case []byte:
		bits := i.([]byte)
		if len(bits) == 8 {
			return int64(binary.LittleEndian.Uint64(bits))
		} else if len(bits) <= 4 {
			num, err := strconv.Atoi(string(bits))
			if err != nil {
				return d
			} else {
				return int64(num)
			}
		}
	case uint:
		return int64(i.(uint))
	case uint8:
		return int64(i.(uint8))
	case uint16:
		return int64(i.(uint16))
	case uint32:
		return int64(i.(uint32))
	case uint64:
		return int64(i.(uint64))
	case int:
		return int64(i.(int))
	case int8:
		return int64(i.(int8))
	case int16:
		return int64(i.(int16))
	case int32:
		return int64(i.(int32))
	case int64:
		return i.(int64)
	case float32:
		return int64(i.(float32))
	case float64:
		return int64(i.(float64))
	}
	return d
}

// toInt64 interface{}轉換成int64
func toInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("無法轉換為 int64")
	}
}

// toFloat64 interface{}轉換成float64
func toFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("無法轉換為 float64")
	}
}

// GetInt64FromMap 從map裡取出資料並轉換為int64類型
func GetInt64FromMap(dm map[string]interface{}, key string, dft int64) int64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return GetInt64(data, dft)
}

// GetInt64FromStringMap 從map裡取出資料並轉換為int64類型
func GetInt64FromStringMap(dm map[string]string, key string, dft int64) int64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return GetInt64(data, dft)
}

// GetFloat64FromStringMap 從map裡取出資料並轉換為float64類型
func GetFloat64FromStringMap(dm map[string]string, key string, dft float64) float64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return GetFloat64(data, dft)
}

// GetFloat64 將string轉換為float64(無法轉換時回傳預設值)
func GetFloat64(value string, dft float64) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return dft
	}
	return f
}

// GetStringFromMap 從map裡取出資料並轉換為string類型
func GetStringFromMap(dm map[string]interface{}, key string, dft string) string {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return GetString(data, dft)
}

// GetStringFromStringMap 從map裡取出資料並轉換為string類型
func GetStringFromStringMap(dm map[string]string, key string, dft string) string {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return data
}

// DownloadFile 將圖片檔案的url下載至遠端
func DownloadFile(url string, filepath string) error {
	// 发起 HTTP GET 请求
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 创建本地文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应中的数据写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// ApplyMapToStruct 將 map[string]string 中的資料自動寫入任意 struct 中（限指標）
// 要求：map 的 key 必須與 struct 的欄位名稱相同（區分大小寫）
// func ApplyMapToStruct(data map[string]string, target interface{}) error {
// 	v := reflect.ValueOf(target)                                // 取得 target 的 reflect.Value 表示，用來操作其內部值

// 	if v.Kind() != reflect.Ptr || v.IsNil() {                   // 如果傳進來的不是指標，或是為 nil
// 		return fmt.Errorf("target 必須是非 nil 的指標")               // 回傳錯誤，因為我們需要一個指向 struct 的有效指標
// 	}

// 	v = v.Elem()                                                // 取得指標所指向的實體（也就是 struct 本身）
// 	if v.Kind() != reflect.Struct {                             // 確認這個值是否為 struct
// 		return fmt.Errorf("target 必須指向一個 struct")              // 如果不是 struct，就回傳錯誤
// 	}

// 	for key, value := range data {                              // 遍歷 map 的每個 key-value
// 		field := v.FieldByName(key)                             // 使用 key 去找 struct 中對應名稱的欄位（大小寫需完全符合）
// 		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
// 			// 檢查欄位是否存在、是否可以設值、且必須是 string 類型
// 			field.SetString(value)                              // 將 map 中的值設置到 struct 的欄位上
// 		}
// 	}

// 	return nil                                                  // 完成後回傳 nil，表示成功
// }
