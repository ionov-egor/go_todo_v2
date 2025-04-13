package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ionov-egor/go_todo_v2/pkg/api"
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"net/http"
)

func Run() error {
	port, _ := config.GetEnvInt("TODO_PORT", 7540)

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/task", api.TaskHandler)
	router.HandleFunc("/api/tasks", api.TasksHandler)
	router.HandleFunc("/api/nextdate", api.NextDayHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
