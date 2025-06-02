package controller

import (
	"hilive/modules/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowCustomizeScene 自定義場景頁面 GET API
func (h *Handler) ShowCustomizeScene(ctx *gin.Context) {
	var (
		host   = ctx.Request.Host
		path   = ctx.Request.URL.Path
		userID = ctx.Query("user_id")
		game   = ctx.Query("game")
		route  string
		// tempName, htmlTmpl string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if userID == "" || game == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識用戶、遊戲資訊，請輸入有效的參數")
		return
	}

	if strings.Contains(path, "customize_scene") {
		route = "./hilives/hilive/views/game/customize_topic.html"
	}

	h.executeHTML(ctx, route, executeParam{})
}
