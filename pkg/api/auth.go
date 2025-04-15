package api

import "net/http"

import "github.com/golang-jwt/jwt/v5"

func Auth(next http.HandlerFunc, pass string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var jwtCookie string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}

			jwtToken := jwt.New(jwt.SigningMethodHS256)
			signedToken, err := jwtToken.SignedString([]byte("my_secret_key"))
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if jwtCookie != signedToken {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

		}
		next(w, r)
	})
}
