package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ionov-egor/go_todo_v2/pkg/api"
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"net/http"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := config.GetEnv("TODO_PASSWORD", "12345")

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

func Run() error {
	port, _ := config.GetEnvInt("TODO_PORT", 7540)

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/signin", api.SignInHandler)
	router.HandleFunc("/api/task", auth(api.TaskHandler))
	router.HandleFunc("/api/tasks", auth(api.TasksHandler))
	router.HandleFunc("/api/task/done", auth(api.TaskDoneHandler))
	router.HandleFunc("/api/nextdate", api.NextDayHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
