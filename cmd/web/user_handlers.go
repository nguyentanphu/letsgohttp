package main

import (
	"fmt"
	"net/http"
)

type userSignupForm struct {
	Name        string            `schema:"name" validate:"required"`
	Email       string            `schema:"email" validate:"required"`
	Password    string            `schema:"password" validate:"required"`
	FieldErrors map[string]string `schema:"-"`
}

type userTemplateData struct {
	Form  userSignupForm
	Flash string
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
