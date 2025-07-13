package auth

import (
	"testing"

	"word_app/backend/src/mocks"
	"word_app/backend/src/mocks/http/auth"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthHandler(t *testing.T) {
	// Arrange: Create mocks for dependencies.
	mockAuthClient := new(auth.MockAuthUsecase)
	mockJWTGenerator := new(mocks.MockJwtGenerator)

	// Act: Create a new AuthHandler instance.
	handler := NewAuthHandler(mockAuthClient, mockJWTGenerator)

	// Assert: Verify that the handler is properly initialized.
	assert.NotNil(t, handler)
	assert.Equal(t, mockAuthClient, handler.AuthUsecase, "AuthUsecase should match the mock instance")
	assert.Equal(t, mockJWTGenerator, handler.jwtGenerator, "jwtGenerator should match the mock instance")
}
