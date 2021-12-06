package repository

import (
	"context"
	"time"

	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTwett(tw models.Twitter) (string, bool, error) {

	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("tweet")
	registry := bson.M{
		"userId":  tw.UserId,
		"message": tw.Message,
		"date":    tw.Date,
	}
	result, err := col.InsertOne(cxt, registry)
	if err != nil {
		return "", false, err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	//hex o string
	return objID.Hex(), true, nil
}
