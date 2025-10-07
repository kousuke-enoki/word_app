package jwt

import (
	"testing"

	"word_app/backend/src/mocks/infrastructure/jwt"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthHandler(t *testing.T) {
	mockTokenValidator := new(jwt.MockTokenValidator)

	newMiddleware := NewMiddleware(mockTokenValidator)

	assert.NotNil(t, newMiddleware)
	assert.Equal(t, mockTokenValidator, newMiddleware.tokenValidator, "tokenValidator should match the mock instance")
}
