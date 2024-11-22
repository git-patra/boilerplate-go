package middleware

import (
	"boilerplate-go/config"
	"net/http"
	"strings"
)

func BasicAuth(config *config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			companyID := r.Header.Get("company_id")
			userID := r.Header.Get("user_id")
			authorizationString := strings.Split(authorization, " ")

			if len(authorizationString) != 2 {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if strings.ToLower(authorizationString[0]) != "basic" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			username, password, ok := r.BasicAuth()

			if !ok {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if username == config.BasicAuthUsername && password == config.BasicAuthPassword && companyID != "" && userID != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		}
		return http.HandlerFunc(fn)
	}
}
