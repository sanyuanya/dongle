package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {

	host := os.Getenv("DB_HOST")

	if host == "" {
		host = "81.70.154.116"
		// host = "127.0.0.1"
	}

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=postgres	password=mysecretpassword dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai", host))

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func Transaction() (*sql.Tx, error) {
	return db.Begin()
}
