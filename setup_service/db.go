package setup_service

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
)

var Db *sql.DB

func InitDB() {
	var err error
	Db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	query, err := ioutil.ReadFile("./migration.sql")

	if err != nil {
		log.Fatal("can't create table", err)
	}
	_, err = Db.Exec(string(query))

	if err != nil {
		log.Fatal("can't create table", err)
	}
}
