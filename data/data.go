package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {

	// host := os.Getenv("DB_HOST")

	// if host == "" {

	// }
	var err error
	db, err = sql.Open("postgres", "host=127.0.0.1 user=postgres	password=mysecretpassword dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

}
