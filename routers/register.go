package routers

import (
	"encoding/json"
	"net/http"

	jsonLog "github.com/rs/zerolog/log"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		jsonLog.Error().Msg("Error en los datos recibidos " + err.Error())
		http.Error(w, "Error en los datos recibidos", 400)
		return
	}
	if len(user.Email) == 0 {
		jsonLog.Error().Msg("Email de usuario requerido" + err.Error())
		http.Error(w, "Email de usuario requerido", 400)
		return
	}

	if len(user.Password) < 6 {
		jsonLog.Error().Msg("Debe especificar password al menos de 6 caracteres" + err.Error())
		http.Error(w, "Debe especificar password al menos de 6 caracteres", 400)
		return
	}

	_, encontrado, _ := repository.CheckExistUser(user.Email)
	if encontrado {
		jsonLog.Error().Msg("Ya existe usuario encontrado" + err.Error())
		http.Error(w, "Ya existe usuario encontrado", 400)
	}

	_, status, err := repository.AddUser(user)

	if err != nil {
		jsonLog.Error().Msg("Error al Registrar usuario" + err.Error())
		http.Error(w, "Error al Registrar usuario", 400)
		return
	}

	if !status {
		jsonLog.Error().Msg("Error al guardar usuario" + err.Error())
		http.Error(w, "Error al guardar usuario", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
