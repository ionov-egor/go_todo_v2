package api

import (
	"net/http"
	"time"
)

import "github.com/ionov-egor/go_todo_v2/pkg/db"

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	}
}

const rowLimit = 10

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue(paramSearch)

	var tasks = []*db.Task{}
	var err error

	switch search {
	case "":
		tasks, err = db.Tasks(rowLimit)

	default:
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			tasks, err = db.TasksByDate(rowLimit, parsedDate.Format(dateLayout))
		} else {
			tasks, err = db.TasksBySearch(rowLimit, search)
		}
	}

	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, r, http.StatusOK, TasksResp{Tasks: tasks})
}
