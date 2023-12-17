package main

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"letsgohttp/internal/models"
	"net/http"
)

type userSignupForm struct {
	Name        string            `schema:"name" validate:"required"`
	Email       string            `schema:"email" validate:"required,email,uniq_email"`
	Password    string            `schema:"password" validate:"required,min=8"`
	FieldErrors map[string]string `schema:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, w)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func userSignUpErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Name":
		return "This field is require"
	case "Email":
		return "This field is required and must be an unique email"
	case "Password":
		return "Password must be at least 8 characters"
	}
	return fe.Error()
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	form := userSignupForm{
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
			for _, fe := range validErr {
				if _, ok := form.FieldErrors[fe.Field()]; !ok {
					form.FieldErrors[fe.Field()] = userSignUpErrorMessage(fe)
				}
			}
			data := app.newTemplateData(r, w)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		}
		return
	}

	app.users.Insert(form.Name, form.Email, form.Password)
	session, _ := app.sessionStore.Get(r, sessionName)
	session.AddFlash("Your signup was successful. Please log in.")
	session.Save(r, w)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email          string            `schema:"email" validate:"required,email"`
	Password       string            `schema:"password" validate:"required,min=8"`
	FieldErrors    map[string]string `schema:"-"`
	NonFieldErrors []string          `schema:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, w)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func userLoginErrorMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		return "This field is required and must be an valid email"
	case "Password":
		return "Password must be at least 8 characters"
	}
	return fe.Error()
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	form := userLoginForm{
		FieldErrors: make(map[string]string),
	}
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	data := app.newTemplateData(r, w)
	data.Form = &form
	err = app.validate.Struct(form)
	if err != nil {
		var validErr validator.ValidationErrors
		if errors.As(err, &validErr) {
			for _, fe := range validErr {
				if _, ok := form.FieldErrors[fe.Field()]; !ok {
					form.FieldErrors[fe.Field()] = userLoginErrorMessage(fe)
				}
			}
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		}
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.NonFieldErrors = append(form.NonFieldErrors, "Email or password is incorrect")
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	session, _ := app.sessionStore.Get(r, sessionName)
	session.Values[authUserIdKey] = id
	session.Save(r, w)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	session, err := app.sessionStore.Get(r, sessionName)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	delete(session.Values, "authenticatedUserID")
	session.Save(r, w)
	session.AddFlash("You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
