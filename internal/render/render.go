package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ashrielbrian/go_bookings/internal/config"
	"github.com/ashrielbrian/go_bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates sets the config for the render package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// reading from memory
	t, ok := tc[tmpl]

	if !ok {
		log.Println("No template set found for ", tmpl)
		return errors.New("Can't get template from cache")
	}

	// write template in-memory to buffer (can also write straight to ResponseWriter)

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	log.Println(td.Error)
	_ = t.Execute(buf, td)

	// could equally do: _ = t.Execute(w, nil)

	// write buffer data to ResponseWriter w
	_, err := buf.WriteTo(w)

	if err != nil {
		log.Fatal("Error writing to ResponseWriter")
		return err
	}
	return nil
}

// CreateTemplateCache creates the go templates and layouts
func CreateTemplateCache() (map[string]*template.Template, error) {
	var cache = map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	// this helps *.page.tmpl files to find layout files, ie
	// parses and replaces "base" with the layout page that defines it
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts

	}
	return cache, nil
}
