package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"log"
	"os"
)

var driverName = fmt.Sprintf("%s", database.DialectSQLite3)

type Storage struct {
	DB *sql.DB
}

func MustLoad(storagePath string) *Storage {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		if _, err := os.Create(storagePath); err != nil {
			log.Panicf("fail to create db: %s", err.Error())
		}
	}

	db, err := sql.Open(driverName, storagePath)

	if err != nil {
		log.Panicf("fail to connect db: %s", err.Error())
	}

	return &Storage{DB: db}
}

func (s Storage) Migrate() error {
	if err := goose.SetDialect(driverName); err != nil {
		return fmt.Errorf("fail to set dialect: %s", err.Error())
	}

	if err := goose.Up(s.DB, "migrations"); err != nil {
		return fmt.Errorf("fail to up migration: %s", err.Error())
	}

	return nil
}
