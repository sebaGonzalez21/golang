package routers

import (
	"encoding/json"
	"net/http"

	"github.com/sagonzalezp/twitt/dto"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
)

func SaveRelations(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		http.Error(w, "Error parametro obligatorio", 400)
		return
	}

	var relation models.Relation
	relation.UserID = IDUser
	relation.UserRelationId = ID

	status, err := repository.AddRelation(relation)
	if err != nil {

		http.Error(w, "Error al agregar relacion", 400)
		return
	}

	if !status {
		http.Error(w, "No se logro crear relacion", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteRelations(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		http.Error(w, "Error parametro obligatorio", 400)
		return
	}

	var relation models.Relation
	relation.UserID = IDUser
	relation.UserRelationId = ID

	status, err := repository.DeleteRelation(relation)
	if err != nil {

		http.Error(w, "Error al eliminar relacion", 400)
		return
	}

	if !status {
		http.Error(w, "No se logro eliminar relacion", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetRelations(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		http.Error(w, "Error parametro obligatorio", 400)
		return
	}

	var relation models.Relation
	relation.UserID = IDUser
	relation.UserRelationId = ID

	var resp dto.ResponseRelation

	status, err := repository.GetRelation(relation)
	if err != nil || !status {
		resp.Status = false
	} else {
		resp.Status = true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
