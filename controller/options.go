package controller

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) OPTIONS(ctx *gin.Context) {
	// ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// // ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, DELETE, GET, PUT, OPTIONS")
}
