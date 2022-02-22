package main

import (
	"net/http"

	"github.com/ashrielbrian/go_bookings/internal/config"
	"github.com/ashrielbrian/go_bookings/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(a *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)

	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	// serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
