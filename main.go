package main

import (
	"database/sql"
	"log"
	"sms/app"
)

func main() {
	DB, error := InitDBWithDSN("sms.db")
	if error != nil {
		log.Fatal(error.Error())
	}
	app.Start(DB)
}
func InitDBWithDSN(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
