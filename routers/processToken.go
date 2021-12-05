package routers

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sagonzalezp/twitt/models"
	"github.com/sagonzalezp/twitt/repository"
)

var Email string
var IDUser string

func ProcessToken(token string) (*models.Claim, bool, string, error) {
	miPass := []byte("masterDesarrollo")
	claims := &models.Claim{}

	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return claims, false, string(""), errors.New("Formato Token invalido")
	}

	token = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return miPass, nil
	})

	if err == nil {
		_, finded, _ := repository.CheckExistUser(claims.Email)
		if finded {
			Email = claims.Email
			IDUser = claims.ID.Hex()
		}
		return claims, finded, IDUser, nil
	}

	if !tkn.Valid {
		return claims, false, string(""), errors.New("Token invalido")
	}

	return claims, false, string(""), err
}
