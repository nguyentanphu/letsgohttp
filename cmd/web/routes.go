package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	r.HandleFunc("/", app.home).Methods("GET")
	r.HandleFunc("/snippet/view/{id:[0-9]+}", app.snippetView).Methods("GET")
	r.HandleFunc("/snippet/create", app.snippetCreateForm).Methods("GET")
	r.HandleFunc("/snippet/create", app.snippetCreate).Methods("POST")

	r.HandleFunc("/user/signup", app.userSignup).Methods("GET")
	r.HandleFunc("/user/signup", app.userSignupPost).Methods("POST")
	r.HandleFunc("/user/login", app.userLogin).Methods("GET")
	r.HandleFunc("/user/login", app.userLoginPost).Methods("POST")
	r.HandleFunc("/user/logout", app.userLogoutPost).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	return standard.Then(r)
}
