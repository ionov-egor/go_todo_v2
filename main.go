package main

import (
	"github.com/joho/godotenv"
)

import (
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"github.com/ionov-egor/go_todo_v2/pkg/db"
	"github.com/ionov-egor/go_todo_v2/pkg/server"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	err := db.Init(config.GetEnv("TODO_DBFILE", "scheduler.db"))
	if err != nil {
		return
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			return
		}
	}()

	err = server.Run()
	if err != nil {
		return
	}

}
