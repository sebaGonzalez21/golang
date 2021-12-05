package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sagonzalezp/twitt/middlew"
	"github.com/sagonzalezp/twitt/routers"
)

/**
Handlers de servicios
*/
func Handlers() {
	//captura el http
	router := mux.NewRouter()
	router.HandleFunc("/registro", middlew.CheckConnection(routers.Register)).Methods("POST")
	PORT := os.Getenv("PORT")
	if PORT != "" {
		PORT = "8080"
	}

	//route de gorrila mux
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
