package middlew

import "net/http"

func ValidateJwt(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, _, err := routes.ProcessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Error en token"+err.Error(), http.StatusBadRequest)
		}
		next.ServeHTTP(w, r)
	}
}
