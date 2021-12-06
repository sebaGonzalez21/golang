package routers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
)

func SaveTweet(w http.ResponseWriter, r *http.Request) {
	var twit models.Twitter
	err := json.NewDecoder(r.Body).Decode(&twit)

	registry := models.Twitter{
		UserId:  IDUser,
		Message: twit.Message,
		Date:    time.Now(),
	}

	_, status, err := repository.AddTwett(registry)
	if err != nil {

		http.Error(w, "Error al agregar twitter", 400)
		return
	}

	if !status {
		http.Error(w, "No se logro crear twitter", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
