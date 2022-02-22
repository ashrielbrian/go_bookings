package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashrielbrian/go_bookings/internal/config"
	"github.com/ashrielbrian/go_bookings/internal/driver"
	"github.com/ashrielbrian/go_bookings/internal/handlers"
	"github.com/ashrielbrian/go_bookings/internal/helpers"
	"github.com/ashrielbrian/go_bookings/internal/models"
	"github.com/ashrielbrian/go_bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app = config.AppConfig{}
var session *scs.SessionManager

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Printf("Application listening on port %s", portNumber)
	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	// to be placed in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when deploying to production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // allows user session to remain after browser window closes
	session.Cookie.SameSite = http.SameSiteLaxMode // go default
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=briant password=")

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	app.UseCache = false
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Error loading app config...")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepository(&app, db)
	render.NewRenderer(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	return db, nil
}
