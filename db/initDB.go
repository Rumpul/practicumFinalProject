package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func CreateDB() (*sqlx.DB, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	if install {
		err = CreateTable(dbFile)
		if err != nil {
			return nil, err
		}
	}

	db, err := OpenSql(dbFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(path string) error {
	db, err := OpenSql(path)

	if err != nil {
		return err
	}

	createQuery := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "19700101",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment VARCHAR(256) NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX date_scheduler on scheduler (date);
	`

	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("table create error: %w", err)
	}
	return db.Close()
}

func OpenSql(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}
	return db, nil
}
