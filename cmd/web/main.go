package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {

	// connect to the database
	db := initDB()

	// create sessions
	sessionManager := initSessionManager()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	appConfig := Config{
		SessionManager: sessionManager,
		DB:             db,
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		Wait:           &wg,
	}

	// set up mail

	// listen for web connections
	appConfig.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println(("Starting web server ..."))
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database ...")
	}

	return conn
}

func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready ...")
			log.Println(err)
		} else {
			log.Println("connected to databse !")
			return connection
		}

		counts++
		if counts > 10 {
			return nil
		}

		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSessionManager() *scs.SessionManager {
	// set up session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}
