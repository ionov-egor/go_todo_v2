package server

import (
	"fmt"
	"net/http"
)

import (
	"github.com/go-chi/chi/v5"
)

import (
	"github.com/ionov-egor/go_todo_v2/pkg/api"
	"github.com/ionov-egor/go_todo_v2/pkg/config"
)

func Run() error {
	port, _ := config.GetEnvInt("TODO_PORT", 7540)
	pass := config.GetEnv("TODO_PASSWORD", "12345")

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/signin", api.SignInHandler)
	router.HandleFunc("/api/task", api.Auth(api.TaskHandler, pass))
	router.HandleFunc("/api/tasks", api.Auth(api.TasksHandler, pass))
	router.HandleFunc("/api/task/done", api.Auth(api.TaskDoneHandler, pass))
	router.HandleFunc("/api/nextdate", api.NextDayHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
