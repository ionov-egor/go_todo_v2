package main

import (
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"github.com/ionov-egor/go_todo_v2/pkg/server"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
	//ToDo написать обработчик ошибок или логирование
}

func main() {
	err := db.Init(config.GetEnv("TODO_DBFILE", "scheduler.db"))
	if err != nil {
		//ToDo написать обработчик ошибок или логирование
		return
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			//ToDo написать обработчик ошибок или логирование
			return
		}
	}()

	err = server.Run()
	if err != nil {
		//ToDo написать обработчик ошибок или логирование
		return
	}

}

//ToDo
//1. Необходимо завести сокращенную ссылку из задания "Запуская web-server" в go.mod и заменить длинные адреса пакетов
//replace go1f => ./
//require go1f v0.0.0
