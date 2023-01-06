package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"web-application/internal/models"
	driver "web-application/internal/models/drivers"

	"github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
)

type config struct {
	port   int 
	api    string 
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
	Session 	  *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr: 				fmt.Sprintf(":%d", 4000),
		Handler:        	app.routes(),
		IdleTimeout: 		30 * time.Second,
		ReadTimeout: 		5 * time.Second, 
		ReadHeaderTimeout:  10 * time.Second,
		WriteTimeout: 		10 * time.Second,
	}

	app.infoLog.Printf("Starting CLIENT %s server in port %d", app.cfg.env, app.cfg.port)

	return srv.ListenAndServe()
}

func main() {
	gob.Register(models.User{})
	var cfg config
	
	cfg.api = "4001"
	cfg.port = 4000
	cfg.db.dsn = "user=postgres dbname=potentially_deployed_site password=Matwyenko1_ host=localhost sslmode=disable"
	cfg.env = "development"

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)

	if err != nil {
		fmt.Println(err)
	}

	tc := make(map[string]*template.Template)
	
	session := scs.New()
	session.Lifetime = 24 * time.Hour

	app := &application{
		cfg:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		database: models.DBModel{Db: conn},
		Session: session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}