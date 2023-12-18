package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

func TestUserSignUpForm(t *testing.T) {
	app, _, _ := newTestApplication()
	ts := newTLSTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/signup")
	csfrToken := extractCSFRToken(t, body)

	assert.NotEmpty(t, csfrToken)
}

func TestUserSignUpPost(t *testing.T) {
	app, userMock, _ := newTestApplication()
	ts := newTLSTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/signup")
	csfrToken := extractCSFRToken(t, body)

	const (
		validName     = "Bobby"
		validEmail    = "phutest@gmail.com"
		validPassword = "Passw0rd!@#"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	t.Run("Valid submission", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", validName)
		form.Add("email", validEmail)
		form.Add("password", validPassword)
		form.Add("csrf_token", csfrToken)
		userMock.On("UniqueEmailValidator", mock.Anything).Return(true)
		userMock.On("Insert", validName, validEmail, validPassword).Return(nil)
		status, _, _ := ts.postForm(t, "/user/signup", form)
		assert.Equal(t, status, http.StatusSeeOther)
	})

	//t.Run("Existing user email", func(t *testing.T) {
	//	form := url.Values{}
	//	form.Add("name", validName)
	//	form.Add("email", validEmail)
	//	form.Add("password", validPassword)
	//	form.Add("csrf_token", csfrToken)
	//	userMock.On("UniqueEmailValidator", mock.Anything).Return(false)
	//	userMock.On("Insert", validName, validEmail, validPassword).Return(nil)
	//	status, _, _ := ts.postForm(t, "/user/signup", form)
	//	assert.Equal(t, status, http.StatusUnprocessableEntity)
	//})

	t.Run("Empty name/email", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "")
		form.Add("email", "")
		form.Add("password", validPassword)
		form.Add("csrf_token", csfrToken)
		userMock.On("UniqueEmailValidator", mock.Anything).Return(false)
		userMock.On("Insert", validName, validEmail, validPassword).Return(nil)
		status, _, body := ts.postForm(t, "/user/signup", form)
		assert.Equal(t, status, http.StatusUnprocessableEntity)
		assert.Contains(t, body, "This field is require")
	})

	t.Run("Weak password", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", validName)
		form.Add("email", validEmail)
		form.Add("password", "1234")
		form.Add("csrf_token", csfrToken)
		userMock.On("UniqueEmailValidator", mock.Anything).Return(false)
		userMock.On("Insert", validName, validEmail, validPassword).Return(nil)
		status, _, body := ts.postForm(t, "/user/signup", form)
		assert.Equal(t, status, http.StatusUnprocessableEntity)
		assert.Contains(t, body, "Password must be at least 8 characters")
	})

	t.Run("Missing csrf token", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", validName)
		form.Add("email", validEmail)
		form.Add("password", validPassword)
		userMock.On("UniqueEmailValidator", mock.Anything).Return(false)
		userMock.On("Insert", validName, validEmail, validPassword).Return(nil)
		status, _, _ := ts.postForm(t, "/user/signup", form)
		assert.Equal(t, status, http.StatusBadRequest)
	})
}
