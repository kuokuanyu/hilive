package controller

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Connection websocket需要使用的chan...等參數
type Connection struct {
	wsConn     *websocket.Conn //websocket.Conn物件
	inChan     chan []byte     //用於接收訊息
	outChan    chan []byte     //用於傳送訊息
	closeChan  chan byte       //幫助內部邏輯判斷連線是否被中斷
	closeMutex sync.Mutex      //用於加鎖(關閉時)
	writeMutex sync.Mutex      //用於加鎖(寫入時)
	// heartMutex sync.Mutex      //用於加鎖(心跳時)
	isClose bool //用於避免重複關閉中的不安全因素
}

// NewWebsocket Connection(struct)
func NewWebsocket(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 5000),
		outChan:   make(chan []byte, 5000),
		closeChan: make(chan byte, 1),
	}

	go conn.readLoop()
	go conn.writeLoop()
	return
}

// NewWebsocketConn Connection(struct)
func NewWebsocketConn(ctx *gin.Context) (wsConn *websocket.Conn, conn *Connection, err error) {
	upgrader := websocket.Upgrader{
		// ReadBufferSize和WriteBufferSize指定I/O缓冲区大小（以字节为单位）。
		// 如果缓冲区大小为零，则使用HTTP服务器分配的缓冲区。
		// I/O缓冲区的大小不限制可以发送或接收的消息的大小。
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}

	wsConn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	// defer wsConn.Close()
	if err != nil {
		return
	}

	conn, err = NewWebsocket(wsConn)
	// defer conn.Close()
	if err != nil {
		conn.Close()
		return
	}

	// if ctx.Request.Host != config.API_URL {
	// 	fmt.Println("網域: ", ctx.Request.Host, ", ", config.API_URL)
	// 	err = errors.New("錯誤: 網域請求發生問題")
	// 	return
	// }
	return
}

// readLoop 隨時讀取前端發送的訊息
func (conn *Connection) readLoop() {
	for {
		_, data, err := conn.wsConn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			conn.Close()
			return
		}
	}
}

// writeLoop WriteMessage將值傳給outChan後，觸動writeLoop發送訊息給前端
func (conn *Connection) writeLoop() {
	var data []byte
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			conn.Close()
			return
		}

		// 避免併發寫入錯誤
		conn.writeMutex.Lock()

		err := conn.wsConn.WriteMessage(websocket.TextMessage, []byte(string(data)))
		if err != nil {
			conn.Close()
			conn.writeMutex.Unlock()
			return
		}

		conn.writeMutex.Unlock()
	}
}

// ReadMessage readloop收到前端訊息(將值給inCahn)後，會觸動ReadMessage(inChan有值)
func (conn *Connection) ReadMessage() (data []byte, err error) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	select {
	case data = <-conn.inChan:
		// log.Println("一般訊息")

	case <-ticker.C:
		// log.Println("收到ticker訊息，傳送PingMessage至前端")
		// 心跳判斷

		// 避免併發寫入錯誤
		conn.writeMutex.Lock()

		conn.wsConn.WriteMessage(websocket.PingMessage, []byte{})

		conn.writeMutex.Unlock()
	case <-conn.closeChan:
		err = errors.New("連線已關閉。")
		return
	}
	return
}

// WriteMessage ReadMessage將inChan值給data後，觸動WriteMessage將data值傳給outChan
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("連線已關閉。")
		return
	}
	return
}

// Close 關閉連線
func (conn *Connection) Close() {
	conn.wsConn.Close()

	conn.closeMutex.Lock()

	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}

	conn.closeMutex.Unlock()
}

// #####chatgpt給的#####
// NewWebsocket Connection(struct) 
// func NewWebsocket(wsConn *websocket.Conn) (conn *Connection, err error) {
// 	if wsConn == nil {
// 		return nil, errors.New("WebSocket 連線未初始化")
// 	}

// 	conn = &Connection{
// 		wsConn:    wsConn,
// 		inChan:    make(chan []byte, 5000),
// 		outChan:   make(chan []byte, 5000),
// 		closeChan: make(chan byte, 1),
// 		// ticker:    time.NewTicker(20 * time.Second),
// 	}

// 	go conn.readLoop()
// 	go conn.writeLoop()
// 	return
// }

// // NewWebsocketConn Connection(struct)
// func NewWebsocketConn(ctx *gin.Context) (wsConn *websocket.Conn, conn *Connection, err error) {
// 	upgrader := websocket.Upgrader{
// 		// ReadBufferSize和WriteBufferSize指定I/O缓冲区大小（以字节为单位）。
// 		// 如果缓冲区大小为零，则使用HTTP服务器分配的缓冲区。
// 		// I/O缓冲区的大小不限制可以发送或接收的消息的大小。
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		}}

// 	wsConn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	// defer wsConn.Close()
// 	if err != nil {
// 		return
// 	}

// 	conn, err = NewWebsocket(wsConn)
// 	// defer conn.Close()
// 	if err != nil {
// 		conn.Close()
// 		return
// 	}

// 	// if ctx.Request.Host != config.API_URL {
// 	// 	fmt.Println("網域: ", ctx.Request.Host, ", ", config.API_URL)
// 	// 	err = errors.New("錯誤: 網域請求發生問題")
// 	// 	return
// 	// }
// 	return
// }

// // readLoop 隨時讀取前端發送的訊息
// func (conn *Connection) readLoop() {
// 	for {
// 		_, data, err := conn.wsConn.ReadMessage()
// 		if err != nil {
// 			conn.Close() // 關閉連線
// 			return
// 		}

// 		select {
// 		case conn.inChan <- data:
// 			// 成功寫入通道
// 		case <-conn.closeChan:
// 			// 若 closeChan 已關閉，則直接退出
// 			// conn.Close()
// 			return
// 		default:
// 			// 如果通道已滿，記錄警告並丟棄消息
// 			log.Println("inChan通道已滿，記錄警告並丟棄消息(readLoop)")
// 		}
// 	}
// }

// // writeLoop WriteMessage將值傳給outChan後，觸動writeLoop發送訊息給前端
// func (conn *Connection) writeLoop() {
// 	var data []byte
// 	for {
// 		// 確保連線未關閉
// 		if conn.isClose {
// 			return
// 		}

// 		select {
// 		case data = <-conn.outChan:
// 		case <-conn.closeChan:
// 			conn.Close()
// 			return
// 		}

// 		// 避免併發寫入錯誤
// 		conn.writeMutex.Lock()
// 		err := conn.wsConn.WriteMessage(websocket.TextMessage, []byte(string(data)))
// 		conn.writeMutex.Unlock() // 手動釋放鎖

// 		if err != nil {
// 			conn.Close()
// 			return
// 		}

// 		// conn.writeMutex.Unlock()
// 	}
// }

// // ReadMessage readloop收到前端訊息(將值給inCahn)後，會觸動ReadMessage(inChan有值)
// func (conn *Connection) ReadMessage() (data []byte, err error) {
// 	ticker := time.NewTicker(20 * time.Second)
// 	defer ticker.Stop()

// 	select {
// 	case data = <-conn.inChan:
// 		// 收到正常消息
// 		// log.Println("一般訊息")
// 	case <-ticker.C:
// 		// log.Println("收到ticker訊息，傳送PingMessage至前端")
// 		// 心跳判斷

// 		// 避免併發寫入錯誤
// 		conn.writeMutex.Lock()
// 		err = conn.wsConn.WriteMessage(websocket.PingMessage, []byte{})
// 		conn.writeMutex.Unlock()

// 		if err != nil {
// 			log.Printf("心跳發送失敗: %v\n", err)
// 			conn.Close()
// 			return
// 		}

// 	case <-conn.closeChan:
// 		err = errors.New("連線已關閉。")
// 		return
// 	}
// 	return
// }

// // WriteMessage ReadMessage將inChan值給data後，觸動WriteMessage將data值傳給outChan
// func (conn *Connection) WriteMessage(data []byte) (err error) {
// 	conn.writeMutex.Lock()         // 加鎖
// 	defer conn.writeMutex.Unlock() // 在函數返回時解鎖

// 	select {
// 	case conn.outChan <- data:
// 		// 成功寫入資料到 outChan
// 	case <-conn.closeChan:
// 		err = errors.New("連線已關閉。")
// 		return
// 	default:
// 		// 通道已滿，記錄警告並丟棄消息
// 		log.Println("outChan 通道已滿，丟棄消息")
// 		err = errors.New("outChan 通道已滿")
// 	}
// 	return
// }

// // Close 關閉連線
// func (conn *Connection) Close() {
// 	conn.closeMutex.Lock()         // 加鎖以保護 isClose
// 	defer conn.closeMutex.Unlock() // 解鎖

// 	if !conn.isClose {
// 		if conn.wsConn != nil {
// 			if err := conn.wsConn.Close(); err != nil {
// 				log.Printf("WebSocket 連線關閉失敗: %v\n", err)
// 			}
// 		}
// 		close(conn.closeChan)
// 		conn.isClose = true
// 	}
// }

// #####chatgpt給的#####

// WebsocketMessage webosocet前後端訊息
// type WebsocketMessage struct {
// 	User        UserModel
// 	Game        GameModel
// 	SignStaffs  []SignStaffModel // 簽到人員
// 	SignPeople  int64            // 簽到人數
// 	PrizeStaff  PrizeStaffModel
// 	PrizeStaffs []PrizeStaffModel
// 	SessionID   string `json:"session_id" example:"session"`  // session ID
// 	Message     string `json:"message" example:"message"`     // 訊息
// 	Error       string `json:"error" example:"error message"` // 錯誤訊息
// }

// UserGamePrizeParam 用戶、場次、獎品相關參數資訊
// type UserGamePrizeParam struct {
// 	User        UserModel         // 用戶資訊
// 	Game        GameModel         // 場次資訊
// 	PrizeStaffs []PrizeStaffModel // 中獎紀錄
// 	Error       string            `json:"error" example:"error message"` // 錯誤訊息
// }

// // Session json parameter
// type Session struct {
// 	SessionID string `json:"session_id"  example:"session_id"` // session ID
// }

// // Message json parameter
// type Message struct {
// 	Message string `json:"message" example:"message"` // 訊息
// }

// form-data呼叫API
// url                 = "https://api.hilives.net/v1/staffmanage/attend"
// client              = &http.Client{}
// buf                 = &bytes.Buffer{}
// writer              = multipart.NewWriter(buf)
// writer.WriteField("user_id", userModel.UserID)
// writer.WriteField("activity_id", gameModel.ActivityID)
// writer.WriteField("game_id", gameID)
// writer.WriteField("name", userModel.Name)
// writer.WriteField("avatar", userModel.Avatar)
// writer.WriteField("status", "success")
// writer.WriteField("black", "0")
// writer.Close()
// // 加入遊戲人員資料
// req, err := http.NewRequest("POST", url, buf)
// if err != nil {
// b, _ := json.Marshal(WebsocketMessage{Error: "發生錯誤: 加入遊戲人員資料發生錯誤"})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
// }
// req.Header.Set("Content-Type", writer.FormDataContentType())
// resp, err := client.Do(req)
// if err != nil {
// 	b, _ := json.Marshal(WebsocketMessage{Error: "發生錯誤: 加入遊戲人員資料發生錯誤"})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
// }
// defer resp.Body.Close()
// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	b, _ := json.Marshal(WebsocketMessage{Error: "發生錯誤: 讀取內容發生錯誤"})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
// }

// json.Unmarshal(body, &result)
// if result.(map[string]interface{})["code"].(float64) == 500 {
// 	b, _ := json.Marshal(WebsocketMessage{Error: result.(map[string]interface{})["msg"].(string)})
// 	conn.wsConn.WriteMessage(websocket.TextMessage, b)
// 	return
//
