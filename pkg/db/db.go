package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    repeat VARCHAR(50),
    comment TEXT
);

CREATE INDEX idx_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return fmt.Errorf("ошибка при проверке файла базы данных: %w", err)
	}

	dbDir := filepath.Dir(dbFile)
	if dbDir == "" {
		dbDir = "."
	}
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("ошибка создания директории для базы данных: %w", err)
		}
	}

	var errOpen error
	db, errOpen = sql.Open("sqlite", dbFile)
	if errOpen != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", errOpen)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка проверки соединения с базой данных: %w", err)
	}

	if install {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("ошибка начала транзакции: %w", err)
		}

		_, err = tx.Exec(schema)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return fmt.Errorf("ошибка создания таблиц: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("ошибка завершения транзакции: %w", err)
		}
	}

	return nil
}

func CloseDB() error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
