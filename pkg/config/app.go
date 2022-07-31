package config

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Connect() {
	d, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("couldn't open the database")
	}
	db = d
	if err = db.Ping(); err != nil {
		log.Fatal("connection to the database is dead")
	}
}

func GetDB() *sql.DB {
	query, err := ioutil.ReadFile("./database/schema/setup.sql")
	if err != nil {
		log.Fatal("couldn't read setup.sql")
	}
	if _, err = db.Exec(string(query)); err != nil {
		log.Fatal("database setup wasn't successful")
	}
	return db
}

func Logger() (*log.Logger, *log.Logger) {
	colorRed, colorGreen, colorReset := "\033[31m", "\033[32m", "\033[0m"
	logInf := log.New(os.Stdout, colorGreen+"INFO\t"+colorReset, log.Ldate|log.Ltime)
	logErr := log.New(os.Stderr, colorRed+"ERROR\t"+colorReset, log.Ldate|log.Ltime|log.Lshortfile)
	return logInf, logErr
}
