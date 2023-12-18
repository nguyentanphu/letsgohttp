package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"letsgohttp/internal/models"
	"net/http"
	"testing"
)

func TestSnippetView(t *testing.T) {
	app, _, snippetMock := newTestApplication()
	ts := newTLSTestServer(t, app.routes())
	defer ts.Close()
	t.Run("Negative Id", func(t *testing.T) {
		url := "/snippet/view/-1"
		status, _, _ := ts.get(t, url)
		assert.Equal(t, status, http.StatusNotFound)
	})
	t.Run("Wrong route", func(t *testing.T) {
		url := "/snippet/view/abc"
		status, _, _ := ts.get(t, url)
		assert.Equal(t, status, http.StatusNotFound)
	})
	t.Run("Wrong route", func(t *testing.T) {
		url := "/snippet/view/abc"
		status, _, _ := ts.get(t, url)
		assert.Equal(t, status, http.StatusNotFound)
	})
	t.Run("Empty id", func(t *testing.T) {
		url := "/snippet/view/"
		status, _, _ := ts.get(t, url)
		assert.Equal(t, status, http.StatusNotFound)
	})
	t.Run("Valid id", func(t *testing.T) {
		id := 123
		content := "Test valid snippet content"
		url := fmt.Sprintf("/snippet/view/%d", id)
		snippetMock.On("Get", id).Return(&models.Snippet{ID: id, Content: content}, nil)
		status, _, body := ts.get(t, url)
		assert.Equal(t, status, http.StatusOK)
		assert.Contains(t, body, content)
	})
	t.Run("Valid id", func(t *testing.T) {
		id := 123
		content := "Test valid snippet content"
		url := fmt.Sprintf("/snippet/view/%d", id)
		snippetMock.On("Get", id).Return(&models.Snippet{ID: id, Content: content}, nil)
		status, _, body := ts.get(t, url)
		assert.Equal(t, status, http.StatusOK)
		assert.Contains(t, body, content)
	})
	t.Run("Not exist id", func(t *testing.T) {
		id := 9999
		url := fmt.Sprintf("/snippet/view/%d", id)
		snippetMock.On("Get", id).Return(nil, models.ErrNoRecords)
		status, _, _ := ts.get(t, url)
		assert.Equal(t, status, http.StatusNotFound)
	})
}
