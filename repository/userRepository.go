package repository

import (
	"context"
	"time"

	jsonLog "github.com/rs/zerolog/log"
	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(u models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//instruccion como ultima instruccion de la funcion
	defer cancel()
	db := db.MongoC.Database("twitter")
	col := db.Collection("users")

	u.Password, _ = security.EncryptPass(u.Password)

	result, err := col.InsertOne(ctx, u)

	if err != nil {
		return "", false, err
	}

	ObjID := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

func GetAllUsers(ID string, page int64, search string, types string) ([]*models.User, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := db.MongoC.Database("twitter")
	col := db.Collection("users")

	var result []*models.User

	options := options.Find() //crear objetos con propiedades que intervienen
	options.SetLimit(20)
	options.SetSkip((page - 1) * 20)

	filter := bson.M{
		"name": bson.M{"$regex": `(?i)` + search},
	}

	//donde se guardan los valores de la bd
	cursor, err := col.Find(cxt, filter, options)
	if err != nil {
		jsonLog.Error().Msg("Problemas de busqueda en usuarios " + err.Error())
		return result, false
	}

	for cursor.Next(context.TODO()) {
		var users models.User
		err := cursor.Decode(&users)
		if err != nil {
			jsonLog.Error().Msg("Problemas al obtener lista de usuarios " + err.Error())
			return result, false
		}

		var rel models.Relation
		rel.UserID = ID
		rel.UserRelationId = users.ID.Hex()

		var finded, include bool
		finded, err = GetRelation(rel)

		if types == "new" && !finded {
			include = true
		}
		if types == "follow" && finded {
			include = true
		}

		if rel.UserRelationId == ID {
			include = false
		}

		if include {
			users.Password = ""
			users.Biography = ""
			users.WebSite = ""
			users.Banner = ""
			users.Email = ""
			result = append(result, &users)
		}

	}

	err = cursor.Err()
	if err != nil {
		jsonLog.Error().Msg("Error presentado al listar todos los usuarios ")
		return result, false
	}
	return result, true
}

func CheckExistUser(email string) (models.User, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//instruccion como ultima instruccion de la funcion
	defer cancel()
	db := db.MongoC.Database("twitter")
	col := db.Collection("users")

	condition := bson.M{"email": email}
	var result models.User

	err := col.FindOne(ctx, condition).Decode(&result)
	//hexadecimal string
	ID := result.ID.Hex()
	if err != nil {
		return result, false, ID
	}
	return result, true, ID
}

func FindProfile(ID string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//instruccion como ultima instruccion de la funcion
	defer cancel()
	db := db.MongoC.Database("twitter")
	col := db.Collection("users")

	var profile models.User
	objID, _ := primitive.ObjectIDFromHex(ID)

	condition := bson.M{
		"_id": objID,
	}

	err := col.FindOne(ctx, condition).Decode(&profile)
	profile.Password = ""
	if err != nil {
		jsonLog.Error().Msg("Registro no encontrado " + err.Error())
		return profile, err
	}
	return profile, nil
}

func Login(email string, password string) (models.User, bool) {
	user, exist, _ := CheckExistUser(email)
	if !exist {
		return user, false
	}
	passwordBytes := []byte(password)
	passwordBd := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(passwordBd, passwordBytes)

	if err != nil {
		return user, false
	}
	return user, true

}

func ModifyUser(u models.User, ID string) (bool, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db := db.MongoC.Database("twitter")
	col := db.Collection("users")

	registre := make(map[string]interface{})

	if len(u.Name) > 0 {
		registre["name"] = u.Name
	}

	if len(u.Name) > 0 {
		registre["lastName"] = u.LastName
	}

	registre["yearBirth"] = u.YearBirth

	if len(u.Avatar) > 0 {
		registre["avatar"] = u.Avatar
	}

	if len(u.Banner) > 0 {
		registre["banner"] = u.Banner
	}

	if len(u.Biography) > 0 {
		registre["biography"] = u.Biography
	}

	if len(u.WebSite) > 0 {
		registre["webSite"] = u.WebSite
	}

	if len(u.Location) > 0 {
		registre["location"] = u.Location
	}

	updateString := bson.M{
		"$set": registre,
	}

	objID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": bson.M{"$eq": objID}}

	_, err := col.UpdateOne(cxt, filter, updateString)

	if err != nil {
		return false, err
	}

	return true, nil
}
