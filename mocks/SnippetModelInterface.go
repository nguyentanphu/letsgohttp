// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	models "letsgohttp/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// SnippetModelInterface is an autogenerated mock type for the SnippetModelInterface type
type SnippetModelInterface struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *SnippetModelInterface) Get(id int) (*models.Snippet, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.Snippet
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.Snippet, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.Snippet); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Snippet)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: title, content, expired
func (_m *SnippetModelInterface) Insert(title string, content string, expired int) (int, error) {
	ret := _m.Called(title, content, expired)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, int) (int, error)); ok {
		return rf(title, content, expired)
	}
	if rf, ok := ret.Get(0).(func(string, string, int) int); ok {
		r0 = rf(title, content, expired)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(title, content, expired)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Latest provides a mock function with given fields:
func (_m *SnippetModelInterface) Latest() ([]*models.Snippet, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Latest")
	}

	var r0 []*models.Snippet
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.Snippet, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.Snippet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Snippet)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSnippetModelInterface creates a new instance of SnippetModelInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSnippetModelInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *SnippetModelInterface {
	mock := &SnippetModelInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}