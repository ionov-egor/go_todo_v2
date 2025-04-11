package main

import (
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"github.com/ionov-egor/go_todo_v2/pkg/server"
)

func main() {
	err := db.Init(config.GetEnv("TODO_DBFILE", "scheduler.db"))
	if err != nil {
		//ToDo написать обработчик ошибок или логирование
		return
	}

	err = server.Run()
	if err != nil {
		//ToDo написать обработчик ошибок или логирование
		return
	}

}
