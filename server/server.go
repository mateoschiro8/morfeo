package server

import (
	"context"
	"fmt"
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
		fmt.Println("HICIERON GET")
	})

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGOURL")))
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

	r.Run(":8000")
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
