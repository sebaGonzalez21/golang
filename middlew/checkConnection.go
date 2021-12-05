package middlew

import (
	"net/http"

	"github.com/sagonzalezp/twitt/db"
)

func CheckConnection(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.CheckConnection() == 0 {
			http.Error(w, "conecion perdida con la bd", 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}
