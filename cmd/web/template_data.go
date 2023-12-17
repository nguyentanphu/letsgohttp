package main

import (
	"github.com/justinas/nosurf"
	"letsgohttp/internal/models"
	"net/http"
)

type appTemplateData struct {
	Snippets        []*models.Snippet
	Snippet         *models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func (app *application) newTemplateData(r *http.Request, w http.ResponseWriter) *appTemplateData {
	session, _ := app.sessionStore.Get(r, sessionName)
	var flash string
	if flashes := session.Flashes(); len(flashes) > 0 {
		flash = flashes[0].(string)
	}
	td := &appTemplateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
	session.Save(r, w)
	return td
}
