package db

import (
	"context"
	"log"

	jsonLog "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* MongoC objecto de conexion a la bd */
var MongoC = ConnectBd()
var clientOptions = options.Client().ApplyURI("mongodb+srv://sagonzalez:mongotest@clustertest.itmur.mongodb.net/testGo?retryWrites=true&w=majority")

/* ConnectBd es la funcion que se conecta la bd */
func ConnectBd() *mongo.Client {
	//coneccion a la bd
	//contextos espacio en memoria donde se comparten cosas
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	jsonLog.Info().Msg("Conexion exitosa con la bd")
	return client
}

/* CheckConnection es ping a la bd*/
func CheckConnection() int {
	err := MongoC.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}
