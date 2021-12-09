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
	col := db.Collection("tweets")
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

func Deletetweets(ID string, UserId string) error {

	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("tweets")

	objID, _ := primitive.ObjectIDFromHex(ID)
	condition := bson.M{
		"_id":    objID,
		"userId": UserId,
	}

	_, err := col.DeleteOne(cxt, condition)
	return err

}

/*
func FindFollowTweets(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	pagTemp, err := strconv.Atoi(page)

	if err != nil {
		jsonLog.Error().Msg("Error lista usuario no encontrado " + err.Error())
		http.Error(w, "Error lista usuario no encontrado", 400)
		return
	}

	pag := int64(pagTemp)

	result, status := repository.GetAllUsers(IDUser, pag, search, typeUser)

	if !status {
		jsonLog.Error().Msg("Error leer usuarios " + err.Error())
		http.Error(w, "Error leer usuarios", 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}*/
func ReadFollowTweets(ID string, page int) ([]models.FollowTweets, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("relations")

	skip := (page - 1) * 20
	conditions := make([]bson.M, 0)
	conditions = append(conditions, bson.M{"$match": bson.M{"userId": ID}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "tweets",
			"localField":   "userRelationId",
			"foreignField": "userId",
			"as":           "tweets",
		}})

	conditions = append(conditions, bson.M{"$unwind": "$tweets"})
	conditions = append(conditions, bson.M{"$sort": bson.M{"date": -1}})
	conditions = append(conditions, bson.M{"$skip": skip})
	conditions = append(conditions, bson.M{"$limit": 20})

	//donde se guardan los valores de la bd
	cursor, err := col.Aggregate(cxt, conditions)
	var result []models.FollowTweets

	err = cursor.All(cxt, &result)

	if err != nil {
		return result, false
	}
	return result, true
}

func Readtweets(ID string, page int64) ([]*models.TwitterList, bool) {

	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("tweets")

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
