package routers

import (
	"encoding/json"
	"net/http"

	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Error en los datos recibidos"+err.Error(), 400)
		return
	}
	if len(user.Email) == 0 {
		http.Error(w, "Email de usuario requerido"+err.Error(), 400)
		return
	}

	if len(user.Password) < 6 {
		http.Error(w, "Debe especificar password al menos de 6 caracteres"+err.Error(), 400)
		return
	}

	_, encontrado, _ := repository.CheckExistUser(user.Email)
	if encontrado {
		http.Error(w, "Ya existe usuario encontrado"+err.Error(), 400)
	}

	_, status, err := repository.AddUser(user)

	if err != nil {
		http.Error(w, "Error al Registrar usuario"+err.Error(), 400)
		return
	}

	if !status {
		http.Error(w, "Error al guardar usuario"+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
