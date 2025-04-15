package api

import (
	"net/http"
	"time"
)

import "github.com/ionov-egor/go_todo_v2/pkg/db"

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getParamId(r)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	task, err := db.GetTaskById(id)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			HandleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		WriteJson(w, r, http.StatusOK, "{}")
		return
	}

	newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.UpdateDate(newDate, id); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, r, http.StatusOK, "{}")
	return
}
