package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) PingHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
