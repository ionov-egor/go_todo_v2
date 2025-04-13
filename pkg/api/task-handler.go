package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"net/http"
	"strconv"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getParamId(r)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.DeleteTask(id); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, r, http.StatusOK, "{}")
	return

}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	var err error

	if _, err := buf.ReadFrom(r.Body); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = checkTitle(task.Title); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = checkRepeat(task.Repeat); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if task.Date, err = getActualDate(task.Date, task.Repeat); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := db.InsertTask(&task)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response := Response{ID: strconv.Itoa(id)}
	WriteJson(w, r, http.StatusOK, response)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	WriteJson(w, r, http.StatusOK, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if task.ID == "" {
		HandleError(w, r, http.StatusBadRequest, fmt.Sprintf("Id задачи не передан"))
		return
	}

	if err = checkTitle(task.Title); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if task.Date, err = getActualDate(task.Date, task.Repeat); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = db.UpdateTask(&task); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, r, http.StatusOK, "{}")
}
