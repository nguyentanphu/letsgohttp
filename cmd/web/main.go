package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"html/template"
	"letsgohttp/internal/models"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	users         *models.UserModel
	templateCache map[string]*template.Template
	formDecoder   *schema.Decoder
	sessionStore  *sessions.CookieStore
	validate      *validator.Validate
}

func main() {
	addr := flag.String("addr", ":4000", "Http network port")
	dsn := flag.String("dsn", "root:sample-password@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte("super-secret"))
	userModel := models.UserModel{DB: db}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("uniq_email", userModel.UniqueEmailValidator)
	formDecoder := schema.NewDecoder()
	formDecoder.IgnoreUnknownKeys(true)
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		users:         &userModel,
		templateCache: templateCache,
		formDecoder:   formDecoder,
		sessionStore:  store,
		validate:      validate,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting http server at port %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
