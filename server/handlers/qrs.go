package handlers

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
)

func HandleQRs(r *gin.Engine) {
	r.GET("/qs", func(c *gin.Context) {
		info, _ := base64.RawURLEncoding.DecodeString(c.Query("data"))
		fmt.Println(string(info))
		c.Redirect(302, "https://google.com")
	})

}
