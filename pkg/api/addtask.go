package api

import (
	"bytes"
	"encoding/json"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getActualDate(date string, repeat string) (string, error) {
	now := time.Now()

	if strings.TrimSpace(date) == "" {
		date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if afterNow(now, t) {
		if len(repeat) == 0 {
			date = now.Format("20060102")
		} else {
			next, err := NextDate(now, date, repeat)
			if err != nil {
				return "", err
			}
			date = next
		}
	}

	return date, nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = сheckTitle(task.Title); err != nil {
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
