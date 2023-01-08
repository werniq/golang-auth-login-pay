package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"web-application/internal/models"

	_ "github.com/lib/pq"
)

type config struct {
	port   int 
	env    string
	db 	   struct {
		dsn     string
	}
	stripe struct {
		secret  string
		key 	string
	}
}

type application struct {
	cfg			  config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	database      models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr: 				fmt.Sprintf(":%d", 4001),
		Handler:        	app.apiRoutes(),
		IdleTimeout: 		30 * time.Second,
		ReadTimeout: 		5 * time.Second, 
		ReadHeaderTimeout:  10 * time.Second,
		WriteTimeout: 		10 * time.Second,
	}

	app.infoLog.Printf("Starting API %s SERVER in port %d", app.cfg.env, app.cfg.port)

	return srv.ListenAndServe()
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		log.Panic(err)
		return nil, err	
	}

	return db, nil
}

func main() {
	var cfg config

	cfg.port = 4001
	cfg.db.dsn = "db.dsn"
	cfg.env = "development"

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := OpenDB(cfg.db.dsn)

	if err != nil {
		fmt.Println(err)
	}

	tc := make(map[string]*template.Template)
	
	app := &application{
		cfg:        	cfg,
		infoLog:        infoLog,
		errorLog:       errorLog,
		templateCache:  tc,
		database: 		models.DBModel{Db: conn},
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
