package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() error {

	db, err := sql.Open("sqlite", "data/banco.db")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	DB = db

	fmt.Println("Banco conectado com sucesso!")

	return nil
}

func GetDB() *sql.DB {
    return DB
}

func Close() error {
    return DB.Close()
}