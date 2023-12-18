package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	ts := newTLSTestServer(t, app.routes())
	defer ts.Close()
	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "OK")
}
