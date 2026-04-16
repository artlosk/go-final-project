package db

import (
	"database/sql"
	"errors"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date_idx ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)

	if install {
		if _, err = db.Exec(schema); err != nil {
			_ = db.Close()
			return err
		}
	}

	DB = db
	return nil
}
