package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mateoschiro8/morfeo/server/handlers"
	"github.com/mateoschiro8/morfeo/server/types"
)

var (
	idx    = 0
	tokens = make(map[string]*types.UserInput)
)

func StartServer() {
	r := gin.Default()

	handlers.LoadTokenControler(&tokens)

	r.GET("/", func(c *gin.Context) {
		fmt.Println("HICIERON GET")
	})

	r.POST("/tokens", handleNewToken)

	handlers.HandleQRs(r)
	handlers.HandleIMGs(r)
	handlers.HandleCSS(r)
	handlers.HandlePDFs(r)

	// COMO EXTRAER EL ID DE LA URL:
	// r.GET("/:id", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	c.String(200, "id = %s", id)
	// 	fmt.Println(tokens[id].Msg)
	// })

	// Esto es vital para que ngrok pase la IP real en el header X-Forwarded-For
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	r.Run(":8000")
}

func handleNewToken(c *gin.Context) {

	var input types.UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tokens[strconv.Itoa(idx)] = &input

	c.String(200, "%d", idx)

	idx++
}