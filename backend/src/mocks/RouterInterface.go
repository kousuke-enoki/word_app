package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockRouter is a mock implementation of the Router interface
type MockRouter struct {
	mock.Mock
}

// Use implements the Router interface
func (m *MockRouter) Use(middleware ...gin.HandlerFunc) {
	m.Called(middleware)
}

// GET implements the Router interface
func (m *MockRouter) GET(path string, handler gin.HandlerFunc) {
	m.Called(path, handler)
}

// POST implements the Router interface
func (m *MockRouter) POST(path string, handler gin.HandlerFunc) {
	m.Called(path, handler)
}

// Group implements the Router interface
func (m *MockRouter) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	args := m.Called(relativePath, handlers)
	return args.Get(0).(*gin.RouterGroup)
}
