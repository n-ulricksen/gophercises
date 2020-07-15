package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "phone_number"
)

func main() {
	fmt.Println("vim-go")

	// Connect to Postgres
	psqlInfo := fmt.Sprintf("host=%v port=%d user=%v password=%v", host, port,
		user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	// Reset database to demonstrate normalize function
	must(resetDB(db, dbname))
	db.Close()

	// Append database name, connect to database
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	must(db.Ping())
	defer db.Close()
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	createDB(db, name)
	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
