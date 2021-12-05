package main

import (
	"log"

	"github.com/sagonzalezp/twitt/db"
	"github.com/sagonzalezp/twitt/handlers"
)

func main() {
	if db.CheckConnection() == 0 {
		log.Fatal("Sin conexion a la BD")
		return
	}
	handlers.Handlers()
}
