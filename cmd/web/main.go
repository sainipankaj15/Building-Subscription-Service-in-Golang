package main

import (
	"database/sql"
	"log"
)

const webPort = ":80"

func main(){
	// connecting database

	db := intiDB()

	db.Ping()

	// Simple web application : Page will render from server side = Creating sessions

	// creating channels

	// creating waitgroup 

	// setp the application config
	
	// When user takes a subscirption , sends emails : Will use gorutine probabaly for this 

	// listeing for web connections 

}

func intiDB() *sql.DB {
	conn := ConnectToDB()

	if conn == nil {
		log.Panic("Cannot connect to database")
	}

	return conn
}
