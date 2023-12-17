package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"letsgohttp/ui"
	"net/http"
)

func (app *application) routes() http.Handler {
	fileServer := http.FileServer(http.FS(ui.Files))
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(fileServer)

	dynamic := alice.New(noSurf, app.authenticate)
	r.Handle("/", dynamic.ThenFunc(app.home)).Methods("GET")
	r.Handle("/snippet/view/{id:[0-9]+}", dynamic.ThenFunc(app.snippetView)).Methods("GET")

	r.Handle("/user/signup", dynamic.ThenFunc(app.userSignup)).Methods("GET")
	r.Handle("/user/signup", dynamic.ThenFunc(app.userSignupPost)).Methods("POST")
	r.Handle("/user/login", dynamic.ThenFunc(app.userLogin)).Methods("GET")
	r.Handle("/user/login", dynamic.ThenFunc(app.userLoginPost)).Methods("POST")

	auth := dynamic.Append(app.requireAuthentication)
	r.Handle("/snippet/create", auth.ThenFunc(app.snippetCreate)).Methods("GET")
	r.Handle("/snippet/create", auth.ThenFunc(app.snippetCreatePost)).Methods("POST")
	r.Handle("/user/logout", auth.ThenFunc(app.userLogoutPost)).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	return standard.Then(r)
}
