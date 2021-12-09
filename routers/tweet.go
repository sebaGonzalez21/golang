package routers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func ListTweet(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar el parametro id", http.StatusBadRequest)
		return
	}

	pageQ := r.URL.Query().Get("page")

	if len(pageQ) < 1 {
		http.Error(w, "Debe enviar el parametro de la pagina", http.StatusBadRequest)
		return
	}

	//conversion de alfa a integer
	page, err := strconv.Atoi(pageQ)
	if err != nil {
		http.Error(w, "Debe enviar el parametro de la pagina con valor mayor a 0", http.StatusBadRequest)
		return
	}

	pag := int64(page)
	response, correct := repository.ReadTweet(ID, pag)
	if !correct {
		http.Error(w, "Error al leer tweets", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
