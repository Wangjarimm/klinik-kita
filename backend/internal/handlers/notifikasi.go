package handlers

import (
	"klinik/internal/middleware"
	"klinik/internal/ws"

	"github.com/gin-gonic/gin"
)

func WSNotifikasi(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	hubWS.ServeNotifikasi(c.Writer, c.Request, ws.RoomUser(claims.UserID))
}
