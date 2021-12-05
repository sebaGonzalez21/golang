package jwt

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sagonzalezp/twitt/models"
)

func GenerateJwt(u models.User) (string, error) {
	miPass := []byte("masterDesarrollo")
	payload := jwt.MapClaims{
		"email":     u.Email,
		"name":      u.Name,
		"lastName":  u.LastName,
		"birthYear": u.YearBirth,
		"biography": u.Biography,
		"location":  u.Location,
		"webSite":   u.WebSite,
		"_id":       u.ID.Hex(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(miPass)
	if err != nil {
		return tokenStr, err
	}
	return strings.Join([]string{"Bearer ", tokenStr}, ""), nil
}
