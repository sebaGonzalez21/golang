package repository

import (
	"context"
	"time"

	jsonLog "github.com/rs/zerolog/log"
	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ReadTweet(ID string, page int64) ([]*models.TwitterList, bool) {

	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("tweet")

	filter := bson.M{
		"userId": ID,
	}

	var result []*models.TwitterList
	options := options.Find() //crear objetos con propiedades que intervienen
	options.SetLimit(20)
	options.SetSort(bson.D{{
		Key:   "date",
		Value: -1, //documentos ordenados en forma decendente
	}})
	options.SetSkip((page - 1) * 20)

	//donde se guardan los valores de la bd
	cursor, err := col.Find(cxt, filter, options)
	if err != nil {
		jsonLog.Error().Msg("Problemas de busqueda " + err.Error())
		return result, false
	}

	for cursor.Next(context.TODO()) {
		var registry models.TwitterList
		err := cursor.Decode(&registry)
		if err != nil {
			jsonLog.Error().Msg("Problemas al obtener lista " + err.Error())
			return result, false
		}
		result = append(result, &registry)
	}
	return result, true

}
