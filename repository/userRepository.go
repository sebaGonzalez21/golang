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
	"golang.org/x/crypto/bcrypt"
)

func AddUser(u models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//instruccion como ultima instruccion de la funcion
	defer cancel()
	db := db.MongoC.Database("testGo")
	col := db.Collection("users")

	u.Password, _ = security.EncryptPass(u.Password)

	result, err := col.InsertOne(ctx, u)

	if err != nil {
		return "", false, err
	}

	ObjID := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

func CheckExistUser(email string) (models.User, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//instruccion como ultima instruccion de la funcion
	defer cancel()
	db := db.MongoC.Database("testGo")
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
	db := db.MongoC.Database("testGo")
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
	db := db.MongoC.Database("testGo")
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
