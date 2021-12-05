package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	jsonLog "github.com/rs/zerolog/log"
	"github.com/sagonzalezp/twitt/middlew"
	"github.com/sagonzalezp/twitt/routers"
)

/**
Handlers de servicios
*/
func Handlers() {
	//captura el http
	router := mux.NewRouter()
	router.HandleFunc("/login", middlew.CheckConnection(routers.Login)).Methods("POST")
	router.HandleFunc("/register", middlew.CheckConnection(routers.Register)).Methods("POST")
	router.HandleFunc("/profile", middlew.CheckConnection(middlew.ValidateJwt(routers.ProfileView))).Methods("GET")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
		jsonLog.Info().Msg("INICIO EN PUERTO.. " + PORT)
	}

	//route de gorrila mux
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
