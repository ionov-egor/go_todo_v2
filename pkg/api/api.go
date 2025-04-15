package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

import "github.com/go-chi/render"

import "github.com/ionov-egor/go_todo_v2/pkg/db"

const (
	paramNow    = "now"
	paramDate   = "date"
	paramRepeat = "repeat"
	paramId     = "id"
	paramSearch = "search"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	ID string `json:"id"`
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func HandleError(w http.ResponseWriter, r *http.Request, status int, message string) {
	errResponse := ErrorResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	render.JSON(w, r, errResponse)
}

func WriteJson(w http.ResponseWriter, r *http.Request, status int, data any) {
	var dataJSON []byte
	var err error

	if data != "{}" {
		dataJSON, err = json.Marshal(data)
		if err != nil {
			HandleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		dataJSON = []byte(`{}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(dataJSON); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func checkTitle(title string) error {
	err := checkTitle4empty(title)
	if err != nil {
		return err
	}
	return nil
}

func checkTitle4empty(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("строка пустая или содержит только пробелы")
	}
	return nil
}

func getParamId(r *http.Request) (int, error) {
	id := r.FormValue(paramId)

	idNum, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if idNum == 0 {
		return 0, err
	}

	return idNum, nil
}
