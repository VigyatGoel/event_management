package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsnWithoutDb := "root:1234@tcp(127.0.0.1:3306)/"

	db, err := sql.Open("mysql", dsnWithoutDb)
	if err != nil {
		log.Fatalf("Error connecting to MySQL server: %v", err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS event_management")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
	log.Println("Database 'event_management' created")

	dsnWithDB := "root:1234@tcp(127.0.0.1:3306)/event_management"

	DB, err = sql.Open("mysql", dsnWithDB)

	if err != nil {
		log.Fatalf("Failed to connect to event_mangement database")
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE,
		password VARCHAR(255),
		isalive BOOLEAN DEFAULT TRUE
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	log.Println("Table 'users' created")
}