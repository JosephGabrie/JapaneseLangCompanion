package main


import (
	"fmt"
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

var DB *sql.DB

func Init() {
    var err error
    DB, err = sql.Open("postgres", "postgres://username:password@localhost/japlearning?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatal(err)
    }

	fmt.Print("Database connected")
    log.Println("Database connected successfully!")
}
