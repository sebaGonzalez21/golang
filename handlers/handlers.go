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
	router.HandleFunc("/logins", middlew.CheckConnection(routers.Login)).Methods("POST")
	router.HandleFunc("/users", middlew.CheckConnection(routers.Register)).Methods("POST")
	router.HandleFunc("/users", middlew.CheckConnection(middlew.ValidateJwt(routers.ProfileView))).Methods("GET")
	router.HandleFunc("/users", middlew.CheckConnection(middlew.ValidateJwt(routers.ModifyProfile))).Methods("PUT")
	router.HandleFunc("/tweets", middlew.CheckConnection(middlew.ValidateJwt(routers.SaveTweet))).Methods("POST")
	router.HandleFunc("/tweets", middlew.CheckConnection(middlew.ValidateJwt(routers.ListTweet))).Methods("GET")
	router.HandleFunc("/tweets", middlew.CheckConnection(middlew.ValidateJwt(routers.DeleteTweet))).Methods("DELETE")
	router.HandleFunc("/avatars", middlew.CheckConnection(middlew.ValidateJwt(routers.AddAvatar))).Methods("POST")
	router.HandleFunc("/avatars", middlew.CheckConnection(middlew.ValidateJwt(routers.GetAvatar))).Methods("GET")
	router.HandleFunc("/banners", middlew.CheckConnection(middlew.ValidateJwt(routers.AddBanner))).Methods("POST")
	router.HandleFunc("/banners", middlew.CheckConnection(middlew.ValidateJwt(routers.GetBanner))).Methods("GET")
	router.HandleFunc("/relations", middlew.CheckConnection(middlew.ValidateJwt(routers.SaveRelations))).Methods("POST")
	router.HandleFunc("/relations", middlew.CheckConnection(middlew.ValidateJwt(routers.DeleteRelations))).Methods("DELETE")
	router.HandleFunc("/relations", middlew.CheckConnection(middlew.ValidateJwt(routers.GetRelations))).Methods("GET")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
		jsonLog.Info().Msg("INICIO EN PUERTO.. " + PORT)
	}

	//route de gorrila mux
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
