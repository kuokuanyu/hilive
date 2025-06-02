package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hilive/models"
	"hilive/modules/utils"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CustomizeSceneReadParam 自定義場景接收訊息
type CustomizeSceneReadParam struct {
	ID           int64  `json:"id" `            // id
	TopicID      string `json:"topic_id" `      // topic_id
	TemplateID   string `json:"template_id" `   // template_id
	UserID       string `json:"user_id" `       // user_id
	Game         string `json:"game" `          // game
	Method       string `json:"method"`         // 方法
	Name         string `json:"name" `          // 圖片名稱
	TemplateName string `json:"template_name" ` // 模板名稱
	Picture      string `json:"picture"`        // 圖片
	// TopicName string                    `json:"topic_name"` // 主題
	CustomizeSceneData map[string]interface{} `json:"customize_scene_data"` // 需儲存的資料
	// CustomizeTemplateData models.CustomizeSceneData `json:"customize_template_data"` // 需儲存的資料
}

// CustomizeSceneWriteParam 自定義場景傳送訊息
type CustomizeSceneWriteParam struct {
	IsSuccess             bool                   `json:"is_success"`              // 是否儲存成功
	URL                   string                 `json:"url"`                     // url
	CustomizeTemplateData map[string]interface{} `json:"customize_template_data"` // 儲存的data資料
}

// @Summary 自定義場景
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param user_id query string true "user ID"
// @Param role query string true "role"
// @param body body CustomizeSceneReadParam true "data"
// @Success 200 {array} CustomizeSceneWriteParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/customize_scene [get]
func (h *Handler) CustomizeSceneWebsocket(ctx *gin.Context) {
	var (
		userID            = ctx.Query("user_id")
		role              = ctx.Query("role")
		wsConn, conn, err = NewWebsocketConn(ctx)
	)

	defer wsConn.Close()
	defer conn.Close()

	// 判斷參數是否有效
	if err != nil || userID == "" {
		b, _ := json.Marshal(SignNameParam{
			Error: "錯誤: 無法辨識用戶資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result CustomizeSceneReadParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		json.Unmarshal(data, &result)

		if result.Method == "upload" &&
			result.Name != "" && result.Picture != "" {
			// 上傳圖片
			// var url string

			// 圖片資料處理(base64格式，data:image/png;base64,...)
			// 去掉前缀
			base64Image := strings.TrimPrefix(result.Picture, "data:image/png;base64,")

			// 解碼bese64格式資料
			imageData, err := base64.StdEncoding.DecodeString(base64Image)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: "錯誤: 解碼圖片資料發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 在cloud storage建立檔案
			f, err := os.Create(fmt.Sprintf("./hilives/hilive/uploads/%s/customize_scene/%s.png",
				userID, result.Name))
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: "錯誤: 建立圖片檔案發生問題"})
				conn.WriteMessage(b)
				return
			}
			defer func() {
				if err2 := f.Close(); err2 != nil {
					err = err2
				}
			}()

			// 將圖片資料寫入檔案
			_, err = f.Write(imageData)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: "錯誤: 圖片資料寫入檔案發生問題"})
				conn.WriteMessage(b)
				return
			}

			// 回傳圖片路徑資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				URL: fmt.Sprintf("/admin/uploads/%s/customize_scene/%s.png",
					userID, result.Name),
			})
			conn.WriteMessage(b)
		} else if result.Method == "add" &&
			utils.GetString(result.CustomizeSceneData["topic_name"], "") != "" {

			// 將data寫入mongo中(customize_scene，用戶自定義場景)
			if err := models.DefaultUserCustomizeSceneModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(userID,
					utils.GetString(result.CustomizeSceneData["topic_name"], ""), result.CustomizeSceneData); err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳圖片路徑資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				IsSuccess: true,
			})
			conn.WriteMessage(b)
		} else if result.Method == "update" &&
			result.TopicID != "" && utils.GetString(result.CustomizeSceneData["topic_name"], "") != "" {

			// 將data更新至mongo中(customize_scene，用戶自定義場景)
			if err := models.DefaultUserCustomizeSceneModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Update(result.TopicID, userID,
					utils.GetString(result.CustomizeSceneData["topic_name"], ""), result.CustomizeSceneData); err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳圖片路徑資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				IsSuccess: true,
			})
			conn.WriteMessage(b)
		} else if result.Method == "template" &&
			result.TemplateID != "" && result.TemplateName != "" {
			// 接收前端模板名稱，取得自定義模板data資料
			templateModel, err := models.DefaultCustomizeTemplateModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(result.TemplateID, result.TemplateName)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳模板data資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				CustomizeTemplateData: templateModel.CustomizeTemplateData,
			})
			conn.WriteMessage(b)
		} else if role == "admin" && // 必須是管理員
			result.Method == "template_add" &&
			result.Game != "" && utils.GetString(result.CustomizeSceneData["topic_name"], "") != "" {
			// 將data寫入mongo中(customize_template，管理員自定義模板)
			err := models.DefaultCustomizeTemplateModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Add(result.Game,
					utils.GetString(result.CustomizeSceneData["topic_name"], ""), result.CustomizeSceneData)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳圖片路徑資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				IsSuccess: true,
			})
			conn.WriteMessage(b)
		} else if role == "admin" && // 必須是管理員
			result.Method == "template_update" &&
			result.TemplateID != "" && utils.GetString(result.CustomizeSceneData["topic_name"], "") != "" {
			// 將data更新至mongo中(customize_template，管理員自定義模板)
			err := models.DefaultCustomizeTemplateModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Update(result.TemplateID,
					utils.GetString(result.CustomizeSceneData["topic_name"], ""), result.CustomizeSceneData)
			if err != nil {
				b, _ := json.Marshal(SignNameParam{
					Error: err.Error()})
				conn.WriteMessage(b)
				return
			}

			// 回傳圖片路徑資料
			b, _ := json.Marshal(CustomizeSceneWriteParam{
				IsSuccess: true,
			})
			conn.WriteMessage(b)
		}
	}
}
