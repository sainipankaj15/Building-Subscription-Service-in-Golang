package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

const webPort = ":80"

func main(){

	// connecting database
	db := intiDB()
	db.Ping()
 
	// Creating session 
	session := initSession()

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

func initSession() *scs.SessionManager{
	
	session := scs.New()

	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24*time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session

}

func initRedis() *redis.Pool{
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error)  {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}