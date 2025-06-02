package logger

import (
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 刪除資料表過期資料
func Logger(conn db.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println("執行Logger")

		// 每日紀錄檔案
		var (
			date   = time.Now().Format("2006-01-02")
			dates  = strings.Split(date, "-")
			day, _ = strconv.Atoi(dates[2])
		)

		// 判斷是否為每個月一號
		if day == 1 {
			// log.Println("每月一號到了")
			var (
				year, _  = strconv.Atoi(dates[0])
				month, _ = strconv.Atoi(dates[1])
			)

			// 判斷要清除的年份、月份
			if month == 1 {
				year-- // 清除前一年的檔案
				if month == 1 {
					month = 12
				}
			} else {
				// 其他月份
				month -= 1
			}

			// 清除過期資料(operation_log資料表)
			db.Conn(conn).Table(config.OPERATION_LOG_TABLE).
				Where("created_at", "<", fmt.Sprintf("%d-%d-01 00:00:00", year, month)).Delete()

			// 清除過期資料(operation_error_log資料表)
			db.Conn(conn).Table(config.OPERATION_ERROR_LOG_TABLE).
				Where("created_at", "<", fmt.Sprintf("%d-%d-01 00:00:00", year, month)).Delete()
		}

	}
}

// Logger 設定輸出檔案的位置、輸出格式、檔案命名等設定，將log紀錄寫入檔案中
// func Logger() gin.HandlerFunc {
// var (
// 	logFilePath string // 檔案路徑(ex: /opt/hilive/src/hilive/logs/)
// 	logger      = logrus.New()
// )

// // 設置log紀錄存放路徑並新建資料夾
// if dir, err := os.Getwd(); err == nil {
// 	logFilePath = dir + "/logs/"
// }
// if err := os.MkdirAll(logFilePath, 0777); err != nil {
// 	// log.Println("錯誤: 建立存放log紀錄的資料夾發生問題")
// }

// 設置logger參數
// logger.SetLevel(logrus.DebugLevel)
// logger.SetFormatter(&logrus.TextFormatter{
// 	TimestampFormat: "2006-01-02 15:04:05",
// })

// return func(c *gin.Context) {
// 每日紀錄檔案
// var (
// 	date        = time.Now().Format("2006-01-02")
// 	dates       = strings.Split(date, "-")
// 	logFileName = date + ".log"                       // 檔案名稱
// 	fileName    = path.Join(logFilePath, logFileName) // ex: /opt/hilive/src/hilive/logs/2023-05-03.log
// 	day, _      = strconv.Atoi(dates[2])
// )

// // 判斷是否為每個月一號
// if day == 1 {
// 	var (
// 		year, _  = strconv.Atoi(dates[0])
// 		month, _ = strconv.Atoi(dates[1])
// 		// year  = 2023
// 		// month = 1
// 		days = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
// 			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
// 			"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
// 			"31"} // 需要清除的檔案日期(1-31號)
// 	)

// 	// fmt.Println("month: ", month)
// 	// 判斷要清除的年份、月份
// 	if month == 1 || month == 2 || month == 3 {
// 		year-- // 清除前一年的檔案
// 		if month == 1 {
// 			month = 10
// 		} else if month == 2 {
// 			month = 11
// 		} else if month == 3 {
// 			month = 12
// 		}
// 	} else {
// 		// 其他月份
// 		month -= 3
// 	}
// 	// fmt.Println("month: ", month)

// 	// 清除過期日誌
// 	for _, dayStr := range days {
// 		deleteFile := path.Join(logFilePath,
// 			fmt.Sprintf("%04d", year)+"-"+fmt.Sprintf("%02d", month)+"-"+dayStr+".log")
// 		// fmt.Println("要刪除的檔案名稱: ", deleteFile)

// 		err := os.Remove(deleteFile)
// 		if err != nil {
// 			// log.Println("錯誤: 清除過期日誌發生問題", err)
// 		}
// 	}
// }

// if !strings.Contains(c.Request.RequestURI, "/undefined") {
// 	// 查詢是否有檔案，沒有則創建一個新的LOG檔
// 	if _, err := os.Stat(fileName); err != nil {
// 		// fmt.Println("os.Stat(fileName)出現問題: ", err)
// 		if _, err := os.Create(fileName); err != nil {
// 			// log.Println("錯誤: 建立每日log紀錄檔案發生問題")
// 		}
// 	}

// 	// 開啟檔案
// 	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		// log.Println("錯誤: 開啟每日log紀錄檔案發生問題")
// 	}
// 	logger.Out = src

// 	// log紀錄
// 	// startTime := time.Now()
// 	c.Next()
// 	// endTime := time.Now()
// 	// latencyTime := endTime.Sub(startTime)
// 	// reqMethod := c.Request.Method
// 	// reqUri := c.Request.RequestURI
// 	// statusCode := c.Writer.Status()
// 	// clientIP := c.ClientIP()

// 	// fmt.Println("statusCode: ", statusCode, errMessage, c.Err(), c.Copy().Errors, c.Errors.Errors(), c.Copy().Err())
// 	logger.Infof("| %3d | %15s | %s | %s |",
// 		c.Writer.Status(),    // API回傳狀態
// 		c.ClientIP(),         // IP位置
// 		c.Request.Method,     // API method
// 		c.Request.RequestURI, // API route
// 	)
// }
// }
// }

// LoggerError 設定輸出檔案的位置、輸出格式、檔案命名等設定，將log錯誤紀錄寫入檔案中
// func LoggerError(c *gin.Context, errMessage string) {
// var (
// 	logFilePath string // 檔案路徑(ex: /opt/hilive/src/hilive/logs/)
// 	logger      = logrus.New()
// )

// // 設置log紀錄存放路徑並新建資料夾
// if dir, err := os.Getwd(); err == nil {
// 	logFilePath = dir + "/logs/"
// }
// if err := os.MkdirAll(logFilePath, 0777); err != nil {
// 	// log.Println("錯誤: 建立存放log紀錄的資料夾發生問題")
// }

// // 設置logger參數
// logger.SetLevel(logrus.DebugLevel)
// logger.SetFormatter(&logrus.TextFormatter{
// 	TimestampFormat: "2006-01-02 15:04:05",
// })

// 每日紀錄檔案
// var (
// 	logFileName = "error.log"                         // 檔案名稱
// 	fileName    = path.Join(logFilePath, logFileName) // ex: /opt/hilive/src/hilive/logs/error.log
// )

// if !strings.Contains(c.Request.RequestURI, "/undefined") {
// 	// 查詢是否有檔案，沒有則創建一個新的LOG檔
// 	if _, err := os.Stat(fileName); err != nil {
// 		// fmt.Println("os.Stat(fileName)出現問題: ", err)
// 		if _, err := os.Create(fileName); err != nil {
// 			// log.Println("錯誤: 建立每日log紀錄檔案發生問題")
// 		}
// 	}

// 	// 開啟檔案
// 	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		// log.Println("錯誤: 開啟每日log紀錄檔案發生問題")
// 	}
// 	logger.Out = src

// 	// log紀錄
// 	// startTime := time.Now()
// 	// c.Next()
// 	// endTime := time.Now()
// 	// latencyTime := endTime.Sub(startTime)
// 	// reqMethod := c.Request.Method
// 	// reqUri := c.Request.RequestURI
// 	// statusCode := c.Writer.Status()
// 	// clientIP := c.ClientIP()

// 	// fmt.Println("statusCode: ", statusCode, errMessage, c.Err(), c.Copy().Errors, c.Errors.Errors(), c.Copy().Err())
// 	logger.Infof("| %15s | %s | %s | %s |",
// 		c.ClientIP(),         // IP位置
// 		c.Request.Method,     // API method
// 		c.Request.RequestURI, // API route,
// 		errMessage,
// 	)
// }

// return
// }

// // Logger 設定輸出檔案的位置、輸出格式、檔案命名等設定
// func Logger() *logrus.Logger {
// 	now := time.Now()
// 	logFilePath := ""

// 	// 設置log紀錄存放路徑並新建資料夾
// 	if dir, err := os.Getwd(); err == nil {
// 		logFilePath = dir + "/logs/"
// 	}
// 	if err := os.MkdirAll(logFilePath, 0777); err != nil {
// 		fmt.Println("錯誤: 建立存放log紀錄的資料夾發生問題")
// 	}

// 	// 每日紀錄檔案
// 	logFileName := now.Format("2006-01-02") + ".log"
// 	fileName := path.Join(logFilePath, logFileName)
// 	if _, err := os.Stat(fileName); err != nil {
// 		if _, err := os.Create(fileName); err != nil {
// 			fmt.Println("錯誤: 建立每日log紀錄檔案發生問題")
// 		}
// 	}

// 	// 開啟檔案
// 	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		fmt.Println("錯誤: 開啟每日log紀錄檔案發生問題")
// 	}

// 	logger := logrus.New()
// 	logger.Out = src
// 	logger.SetLevel(logrus.DebugLevel)
// 	logger.SetFormatter(&logrus.TextFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 	})
// 	return logger
// }

// // LoggerToFile 將log紀錄寫入檔案中
// func LoggerToFile() gin.HandlerFunc {
// 	logger := Logger()
// 	return func(c *gin.Context) {
// 		startTime := time.Now()
// 		c.Next()

// 		endTime := time.Now()
// 		latencyTime := endTime.Sub(startTime)
// 		reqMethod := c.Request.Method
// 		reqUri := c.Request.RequestURI
// 		statusCode := c.Writer.Status()
// 		clientIP := c.ClientIP()
// 		logger.Infof("| %3d | %13v | %15s | %s | %s |",
// 			statusCode,
// 			latencyTime,
// 			clientIP,
// 			reqMethod,
//
