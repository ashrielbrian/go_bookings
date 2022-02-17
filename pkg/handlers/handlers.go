package handlers

import (
	"net/http"

	"github.com/ashrielbrian/go_bookings/pkg/config"
	"github.com/ashrielbrian/go_bookings/pkg/models"
	"github.com/ashrielbrian/go_bookings/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepository creates a new repository
func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the Repository for the handlers package
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	var smp = make(map[string]string)
	smp["hello"] = "Brian Tang"
	smp["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: smp,
	})
}
