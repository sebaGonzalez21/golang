package routers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func FindAllUsers(w http.ResponseWriter, r *http.Request) {
	typeUser := r.URL.Query().Get("type")
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	pagTemp, err := strconv.Atoi(page)

	if err != nil {
		jsonLog.Error().Msg("Error lista usuario no encontrado " + err.Error())
		http.Error(w, "Error lista usuario no encontrado", 400)
		return
	}

	pag := int64(pagTemp)

	result, status := repository.GetAllUsers(IDUser, pag, search, typeUser)

	if !status {
		jsonLog.Error().Msg("Error leer usuarios " + err.Error())
		http.Error(w, "Error leer usuarios", 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
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
		jsonLog.Error().Msg("Debe enviar parametro id")
		http.Error(w, "Debe enviar parametro id", http.StatusBadRequest)
		return
	}

	profile, err := repository.FindProfile(ID)
	if err != nil {
		jsonLog.Error().Msg("Ocurrio error al buscar el registro " + err.Error())
		http.Error(w, "Ocurrio error al buscar el registro", 400)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func AddAvatar(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("avatar")
	var extension = strings.Split(handler.Filename, ".")[1]
	var archive string = "uploads/avatars/" + IDUser + "." + extension

	f, err := os.OpenFile(archive, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		jsonLog.Error().Msg("Error subir imagen " + err.Error())
		http.Error(w, "Error subir image", 400)
	}
	_, err = io.Copy(f, file)
	if err != nil {
		jsonLog.Error().Msg("Error al copiar la imagen " + err.Error())
		http.Error(w, "Error al copiar la imagen", 400)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var user models.User
	user.Avatar = IDUser + "." + extension
	var status bool
	status, err = repository.ModifyUser(user, IDUser)

	if err != nil || !status {
		jsonLog.Error().Msg("Error al guardar imagen en la bd " + err.Error())
		http.Error(w, "Error al guardar imagen en la bd ", 400)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func GetAvatar(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar parametro id", http.StatusBadRequest)
		return
	}

	profile, err := repository.FindProfile(ID)

	if err != nil {
		jsonLog.Error().Msg("Error usuario no encontrado " + err.Error())
		http.Error(w, "Error usuario no encontrado", 400)
		return
	}

	OpenFile, err := os.Open("uploads/avatars/" + profile.Avatar)

	if err != nil {
		jsonLog.Error().Msg("Error imagen no encontrada " + err.Error())
		http.Error(w, "Error imagen no encontrada", 400)
		return
	}

	_, err = io.Copy(w, OpenFile)

	if err != nil {
		jsonLog.Error().Msg("Error al copiar imagen " + err.Error())
		http.Error(w, "Error al copiar imagen ", 400)
	}

}

func AddBanner(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("banner")
	var extension = strings.Split(handler.Filename, ".")[1]
	var archive string = "uploads/banners/" + IDUser + "." + extension

	f, err := os.OpenFile(archive, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		jsonLog.Error().Msg("Error al subir banner " + err.Error())
		http.Error(w, "Error al subir banner", 400)
	}
	_, err = io.Copy(f, file)
	if err != nil {
		jsonLog.Error().Msg("Error al copiar el banner " + err.Error())
		http.Error(w, "Error al copiar el banner", 400)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var user models.User
	user.Banner = IDUser + "." + extension
	var status bool
	status, err = repository.ModifyUser(user, IDUser)

	if err != nil || !status {
		jsonLog.Error().Msg("Error al guardar banner en la bd " + err.Error())
		http.Error(w, "Error al guardar banner en la bd ", 400)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func GetBanner(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		http.Error(w, "Debe enviar parametro id", http.StatusBadRequest)
		return
	}

	profile, err := repository.FindProfile(ID)

	if err != nil {
		jsonLog.Error().Msg("Error usuario no encontrado " + err.Error())
		http.Error(w, "Error usuario no encontrado", 400)
		return
	}

	OpenFile, err := os.Open("uploads/banners/" + profile.Banner)

	if err != nil {
		jsonLog.Error().Msg("Error banner no encontrada " + err.Error())
		http.Error(w, "Error banner no encontrada", 400)
		return
	}

	_, err = io.Copy(w, OpenFile)

	if err != nil {
		jsonLog.Error().Msg("Error al copiar banner " + err.Error())
		http.Error(w, "Error al copiar banner ", 400)
	}

}

func ModifyProfile(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonLog.Error().Msg("Datos Incorrectos " + err.Error())
		http.Error(w, "Datos Incorrectos", 400)
		return
	}
	status, err := repository.ModifyUser(user, IDUser)
	if err != nil {
		jsonLog.Error().Msg("Error al modificar el registro" + err.Error())
		http.Error(w, "Error al modificar el registro", 400)
		return
	}

	if !status {
		jsonLog.Error().Msg("No se encontro el usuario " + err.Error())
		http.Error(w, "No se encontro el usuario", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}
