package controller

import (
	"encoding/json"
	"hilive/models"

	"github.com/gin-gonic/gin"
)

// LogParam 接收日誌訊息
type LogParam struct {
	UserID string `json:"user_id" example:"user_id"`                  // 用戶ID
	Method string `json:"method" example:"GET、POST、PATCH、PUT、DELETE"` // 方法
	Path   string `json:"path" example:"path"`                        // 路徑
	Error  string `json:"error" example:"error message"`              // 錯誤訊息
}

// @Summary 新增用戶操作日誌資料
// @Tags Websocket
// @version 1.0
// @Accept  json
// @param body body LogParam true "user_id、method、path... param"
// @Success 200 {array} LogParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/log [get]
func (h *Handler) LogWebsocket(ctx *gin.Context) {
	var (
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            LogParam
	)
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil {
		b, _ := json.Marshal(LogParam{
			Error: "錯誤: 錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result LogParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// log.Println("關閉? ", err)
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)
		// if result.UserID == "" || result.Method == "" ||
		// 	result.Path == "" {
		// 	log.Println("LOG錯誤?")
		// 	b, _ := json.Marshal(LogParam{
		// 		Error: "錯誤: 取得操作日誌資料發生問題"})
		// 	conn.WriteMessage(b)
		// 	return
		// }

		// log.Println("收到LOG資料訊息: ", result)

		if result.UserID != "" && result.Method != "" &&
			result.Path != "" {
			// log.Println("LOG執行")

			// 新增操作日誌資料
			if err = models.DefaultLogModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(models.EditLogModel{
					UserID: result.UserID,
					Method: result.Method,
					Path:   result.Path,
				}); err != nil {
				b, _ := json.Marshal(LogParam{
					Error: "錯誤: 新增操作日誌資料發生問題"})
				conn.WriteMessage(b)
				return
			}
		}

		// b, _ := json.Marshal(ChatroomRecordWriteParam{
		// 	ChatroomRecords: newQuestions,
		// })
		// conn.WriteMessage(b)
	}
}
