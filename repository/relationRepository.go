package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddRelation(t models.Relation) (bool, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("relations")
	_, err := col.InsertOne(cxt, t)

	if err != nil {
		return false, err
	}

	//hex o string
	return true, nil
}

func GetRelation(t models.Relation) (bool, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("relations")

	condition := bson.M{
		"userId":         t.UserID,
		"userRelationId": t.UserRelationId,
	}

	var result models.Relation
	fmt.Println(result)

	err := col.FindOne(cxt, condition).Decode(&result)

	if err != nil {
		return false, err
	}

	//hex o string
	return true, nil
}

func DeleteRelation(t models.Relation) (bool, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("relations")
	_, err := col.DeleteOne(cxt, t)

	if err != nil {
		return false, err
	}
	//hex o string
	return true, nil
}
