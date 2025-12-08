package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleBINs(r *gin.Engine) {
	r.GET("/bins/:tokenID", func(c *gin.Context) {
		tokenID := c.Param("tokenID")

		controller := c.MustGet("tokenController").(*TokenController)
		token, err := controller.GetToken(tokenID)
		if err != nil {
			c.JSON(404, gin.H{"error": "token not found"})
			return
		}
		alertText := "Fue activado el token " + strings.ToLower(token.Msg) + " desde la IP: " + c.ClientIP()
		Alert(alertText, token.Chat)
	})
}
