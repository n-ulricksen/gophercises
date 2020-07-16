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
	dbname   = "gophercises_phonenumbers"
)

type phoneNumber struct {
	id     string
	number string
}

func main() {
	fmt.Println("vim-go")

	// Connect to local Postgres server
	psqlInfo := fmt.Sprintf("host=%v port=%d user=%v password=%v", host, port,
		user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	// Reset database to demonstrate normalize function
	must(resetDB(db, dbname))
	db.Close()

	// Append database name to the psqlInfo string, reconnect to database
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	must(db.Ping())
	defer db.Close()
	fmt.Println("Connected to psql...")

	// Create the phone numbers table if it does not yet exist
	err = createPhoneNumbersTable(db)
	must(err)

	// Insert the unformatted phone numbers into the phone_number table
	unformattedNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range unformattedNumbers {
		err := insertPhoneNumber(db, number)
		must(err)
	}

	// Normalize all phone numbers in phone_number table
	phoneNumbers, err := getAllPhoneNumbers(db)
	fmt.Println()
	fmt.Println("Not normalized phone numbers:")
	for _, p := range phoneNumbers {
		fmt.Println(p.number)
	}
	must(err)
	updateStatement := `
		UPDATE phone_number
		SET value = $1
		WHERE id = $2`
	for _, p := range phoneNumbers {
		normalized := normalize(p.number)
		_, err = db.Exec(updateStatement, normalized, p.id)
		must(err)
	}
	fmt.Println()
	fmt.Println("Normalized all entries in 'phone_number' table...")
	fmt.Println()

	// Verify numbers are normalized
	phoneNumbers, err = getAllPhoneNumbers(db)
	fmt.Println("Normalized phone numbers:")
	for _, p := range phoneNumbers {
		fmt.Println(p.number)
	}
}

// Normalize a phone number by stripping all non-numeric characters
func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

// Delete and recreate the specified database
func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

// Create a new database with the given name
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	createDB(db, name)
	return nil
}

// Create a phone_number table in the database if it does not yet exist
func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_number (
			id SERIAL,
			value VARCHAR(255)
		)`

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

// Insert a record into the phone_number table in the database
func insertPhoneNumber(db *sql.DB, phone string) error {
	statement := `
		INSERT INTO phone_number(value)
		VALUES ($1)
	`
	_, err := db.Exec(statement, phone)
	if err != nil {
		return err
	}
	fmt.Println("Inserted value:", phone)
	return nil
}

// Get all phone numbers from the phone_number table
func getAllPhoneNumbers(db *sql.DB) ([]phoneNumber, error) {
	query := `
		SELECT * FROM phone_number`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []phoneNumber
	for rows.Next() {
		var p phoneNumber
		err = rows.Scan(&p.id, &p.number)
		if err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, p)
	}

	return phoneNumbers, nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
