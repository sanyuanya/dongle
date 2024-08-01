package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {

	var err error
	db, err = sql.Open("postgres", "user=postgres	password=mysecretpassword dbname=postgres sslmode=disable TimeZone=Asia/Shanghai")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

}
