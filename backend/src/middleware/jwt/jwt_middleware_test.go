package jwt

import (
	"testing"
	"word_app/backend/src/mocks/http/middleware"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthHandler(t *testing.T) {
	mockTokenValidator := new(middleware.MockTokenValidator)

	new_middleware := NewJwtMiddleware(mockTokenValidator)

	assert.NotNil(t, new_middleware)
	assert.Equal(t, mockTokenValidator, new_middleware.tokenValidator, "tokenValidator should match the mock instance")
}
