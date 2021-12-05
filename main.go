package main

import (
	"log"

	"github.com/sagonzalezp/twitt/bd"
	"github.com/sagonzalezp/twitt/handlers"
)

func main() {
	if bd.CheckConnection() == 0 {
		log.Fatal("Sin conexion a la BD")
		return
	}
	handlers.Handlers()
}
