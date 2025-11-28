package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mateoschiro8/morfeo/server/types"
)

func HandleBINs(r *gin.Engine) {
	r.GET("/bins/:tokenID", func(c *gin.Context) {
		tokenID := c.Param("tokenID")
		var token *types.UserInput = TC.GetToken(tokenID)
		Alert(token.Msg)
	})
}
