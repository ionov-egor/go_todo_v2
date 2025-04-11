package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ionov-egor/go_todo_v2/pkg/api"
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"net/http"
	"os"
	"strconv"
)

func getPort() (int, error) {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return 7540, nil // значение по умолчанию
	}

	if len(portStr) < 2 || len(portStr) > 5 {
		return 0, fmt.Errorf("неверный формат порта: %s", portStr)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Ошибка при преобразовании:", err)
		return 0, err
	}

	return port, nil
}

func Run() error {
	port, _ := config.GetEnvInt("TODO_PORT", 7540)

	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/api/task", api.TaskHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
