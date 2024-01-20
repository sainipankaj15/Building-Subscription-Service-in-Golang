package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func ConnectToDB() *sql.DB {

	countAttempts := 0

	databaseString := os.Getenv("DSN")

	for {
		connection, err := openDB(databaseString)

		if err != nil {
			log.Println("Postgresql database is not connecteed")
		}else{
			log.Println("Postgresql database is connected")
			return connection
		}

		if countAttempts > 10 {
			return nil
		}

		log.Println("Backing off for 1 second")
		time.Sleep( 1 * time.Second)
		countAttempts ++
		continue
	}
}


// Connecting to DB
func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db , nil
}
