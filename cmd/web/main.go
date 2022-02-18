package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashrielbrian/go_bookings/internal/config"
	"github.com/ashrielbrian/go_bookings/internal/handlers"
	"github.com/ashrielbrian/go_bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app = config.AppConfig{}
var session *scs.SessionManager

func main() {

	// change this to true when deploying to production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // allows user session to remain after browser window closes
	session.Cookie.SameSite = http.SameSiteLaxMode // go default
	session.Cookie.Secure = app.InProduction

	app.Session = session

	app.UseCache = false

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Error loading app config...")
	}

	app.TemplateCache = tc

	render.NewTemplates(&app)
	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// _ = http.ListenAndServe(portNumber, nil)

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
