package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"html"
	"io"
	"letsgohttp/mocks"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func newTLSTestServer(t *testing.T, handler http.Handler) *testServer {
	ts := httptest.NewTLSServer(handler)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, url string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

func newTestApplication() (*application, *mocks.UserModelInterface, *mocks.SnippetModelInterface) {
	store := sessions.NewCookieStore([]byte("test-secret"))
	formDecoder := schema.NewDecoder()
	formDecoder.IgnoreUnknownKeys(true)
	userMock := &mocks.UserModelInterface{}
	snippetMock := &mocks.SnippetModelInterface{}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("uniq_email", userMock.UniqueEmailValidator)
	templateCache, _ := newTemplateCache()
	return &application{
		users:         userMock,
		snippets:      snippetMock,
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		sessionStore:  store,
		validate:      validate,
		formDecoder:   formDecoder,
		templateCache: templateCache,
	}, userMock, snippetMock
}

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSFRToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no CSRF token found")
	}

	return html.UnescapeString(matches[1])
}

func (ts *testServer) postForm(t *testing.T, url string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+url, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}
