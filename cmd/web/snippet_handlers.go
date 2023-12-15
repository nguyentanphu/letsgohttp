package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"letsgohttp/internal/models"
	"net/http"
	"strconv"
)

var sessionStoreName = "Session-Store"

type snippetTemplateData struct {
	Snippets []*models.Snippet
	Snippet  *models.Snippet
	Form     snippetCreateForm
	Flash    string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	session, err := app.sessionStore.Get(r, sessionStoreName)
	var flash string
	if flashes := session.Flashes(); len(flashes) > 0 {
		flash = flashes[0].(string)
	}
	templateData := snippetTemplateData{
		Snippets: snippets,
		Flash:    flash,
	}
	session.Save(r, w)
	app.render(w, http.StatusOK, "home.tmpl.html", templateData)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	session, err := app.sessionStore.Get(r, sessionStoreName)
	var flash string
	if flashes := session.Flashes(); len(flashes) > 0 {
		flash = flashes[0].(string)
	}
	templateData := snippetTemplateData{
		Snippet: snippet,
		Flash:   flash,
	}
	session.Save(r, w)
	app.render(w, http.StatusOK, "view.tmpl.html", templateData)
}

type snippetCreateForm struct {
	Title       string            `validate:"required,max=100" schema:"title"`
	Content     string            `validate:"required" schema:"content"`
	Expires     int               `validate:"oneof=1 7 365" schema:"expires"`
	FieldErrors map[string]string `form:"-"`
}

func snippetErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Title":
		return "This field is require"
	case "Content":
		return "This field is required and must be less than 100 character"
	case "Expires":
		return "This field must be 1, 7 or 365"
	}
	return fe.Error()
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	form := snippetCreateForm{
		FieldErrors: make(map[string]string),
	}
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.validate.Struct(form)
	if err != nil {
		var validErr validator.ValidationErrors
		if errors.As(err, &validErr) {
			errors.As(err, &validErr)
			for _, fe := range validErr {
				if _, ok := form.FieldErrors[fe.Field()]; !ok {
					form.FieldErrors[fe.Field()] = snippetErrorMessage(fe)
				}
			}
			data := snippetTemplateData{Form: form}
			app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		}
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	session, _ := app.sessionStore.Get(r, sessionStoreName)
	session.AddFlash("Snippet successfully created!")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := snippetTemplateData{
		Form: snippetCreateForm{
			Expires: 365,
		}}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}
