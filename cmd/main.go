package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func main() {
	dbUser := "root"
	dbPass := "password"
	dbHost := "localhost:3306"
	dbName := "haveibeenleaked"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error when open database connection: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error when pinging database: %v", err)
	}

	log.Println("Connection successfully stabilized")

	server := NewApiServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Error when starting API: %v", err)
	}
}
