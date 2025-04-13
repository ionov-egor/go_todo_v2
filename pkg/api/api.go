package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"net/http"
	"strings"
)

const (
	paramNow    = "now"
	paramDate   = "date"
	paramRepeat = "repeat"
	paramId     = "id"
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
	dataJSON, err := json.Marshal(data)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dataJSON)
}

func сheckTitle(title string) error {
	err := сheckTitle4empty(title)
	if err != nil {
		return err
	}
	return nil
}

func сheckTitle4empty(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("строка пустая или содержит только пробелы")
	}
	return nil
}
