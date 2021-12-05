package routers

import (
	"encoding/json"
	"net/http"
	"time"

	jsonLog "github.com/rs/zerolog/log"
	"github.com/sagonzalezp/twitt/dto"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
	"github.com/sagonzalezp/twitt/security/jwt"
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
		jsonLog.Error().Msg("Email de usuario requerido")
		http.Error(w, "Email de usuario requerido", 400)
		return
	}

	if len(user.Password) < 6 {
		jsonLog.Error().Msg("Debe especificar password al menos de 6 caracteres")
		http.Error(w, "Debe especificar password al menos de 6 caracteres", 400)
		return
	}

	_, encontrado, _ := repository.CheckExistUser(user.Email)
	if encontrado {
		jsonLog.Error().Msg("Ya existe usuario encontrado")
		http.Error(w, "Ya existe usuario encontrado", 400)
		return
	}

	_, status, err := repository.AddUser(user)

	if err != nil {
		jsonLog.Error().Msg("Error al Registrar usuario " + err.Error())
		http.Error(w, "Error al Registrar usuario", 400)
		return
	}

	if !status {
		jsonLog.Error().Msg("Error al guardar usuario ")
		http.Error(w, "Error al guardar usuario", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		jsonLog.Error().Msg("Error al loguearse" + err.Error())
		http.Error(w, "Usuario y/o contraseña invalido ", 400)
		return
	}

	if len(user.Email) == 0 {
		jsonLog.Error().Msg("Error en login")
		http.Error(w, "Email de usuario requerido ", 400)
		return
	}

	documento, existe := repository.Login(user.Email, user.Password)

	if !existe {
		jsonLog.Error().Msg("Error al loguearse")
		http.Error(w, "Usuario y/o contraseña invalido", 400)
		return
	}

	jwtKey, err := jwt.GenerateJwt(documento)

	if err != nil {
		jsonLog.Error().Msg("Error al generar token " + err.Error())
		http.Error(w, "Error al generar token", 400)
		return
	}
	resp := dto.ResponseLogin{
		Token: jwtKey,
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)

	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: expirationTime,
	})
}

func ProfileView(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar parametro id", http.StatusBadRequest)
		return
	}

	profile, err := repository.FindProfile(ID)
	if err != nil {
		http.Error(w, "Ocurrio error al buscar el registro "+err.Error(), 400)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}
