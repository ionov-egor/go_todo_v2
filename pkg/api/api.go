package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"log"
	"net/http"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := fmt.Sprintf("ID: %s\nDate: %s\nComment: %s\n Repeat: %s\n", task.ID, task.Date, task.Comment, task.Repeat)
	log.Println(s)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
