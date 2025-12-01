package server

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mateoschiro8/morfeo/server/handlers"
	"github.com/mateoschiro8/morfeo/server/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartServer() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(morfeoString))
	})

	mongoURL := "mongodb+srv://" + os.Getenv("MONGOUSER") + ":" + os.Getenv("MONGOPW") + "@" + os.Getenv("MONGOCLUSTER") + "/?appName=" + os.Getenv("MONGOAPP")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}

	collection := client.Database("fcen").Collection("tokens")
	tokenController := handlers.NewTokenController(collection)

	// Middleware para que el controller esté disponible en todos los handlers
	r.Use(func(c *gin.Context) {
		c.Set("tokenController", tokenController)
		c.Next()
	})

	r.POST("/tokens", handleNewToken)

	handlers.HandleQRs(r)
	handlers.HandleIMGs(r)
	handlers.HandleCSS(r)
	handlers.HandlePDFs(r)
	handlers.HandleBINs(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r.Run(":" + port)
}

func handleNewToken(c *gin.Context) {

	var input types.UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	controller := c.MustGet("tokenController").(*handlers.TokenController)

	res, err := controller.Collection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(500, gin.H{"error": "could not cast inserted ID"})
		return
	}

	c.String(200, "%s", oid.Hex())
}

var morfeoString = `
<html>
	<head>
		<meta charset="utf-8">
		<style>
			html,body{height:100%;margin:0}
			body{display:flex;align-items:center;justify-content:center;background:#000}
			span{font-size:10vw;font-weight:800;color:#fff;font-family:system-ui,Segoe UI,Arial,Helvetica,sans-serif}
		</style>
	</head>
	<body>
		<span>
			M<span style="color:blue">o</span>rfe<span style="color:red">o</span>
		</span>
	</body>
</html>`
