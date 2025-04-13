package api

import (
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"net/http"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	}
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(10) // в параметре максимальное количество записей

	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, r, http.StatusOK, TasksResp{Tasks: tasks})
}
