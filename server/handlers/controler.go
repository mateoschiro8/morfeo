package handlers

import (
	"context"

	"github.com/mateoschiro8/morfeo/server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenController struct {
	Collection *mongo.Collection
}

func NewTokenController(col *mongo.Collection) *TokenController {
	return &TokenController{Collection: col}
}

func (tc *TokenController) GetToken(id string) (*types.UserInput, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var token types.UserInput
	err = tc.Collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
