package controller

import (
	"hilive/modules/config"

	"github.com/gin-gonic/gin"
)

// ShowIndex 首頁 API URL
func (h *Handler) ShowIndex(ctx *gin.Context) {
	var (
		host     = ctx.Request.Host
		htmlTmpl = "./hilives/hilive/views/website/index.html"
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{})
}
