package repository

import (
	"context"
	"time"

	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	ObjId, _ := result.InsertedID.(primitive.ObjectID)
	return ObjId.String(), true, nil
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
