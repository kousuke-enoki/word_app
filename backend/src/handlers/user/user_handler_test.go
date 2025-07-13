package user

import (
	"testing"

	"word_app/backend/src/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	// Arrange: Create mocks for dependencies.
	mockUserClient := new(mocks.UserClient)
	mockJWTGenerator := new(mocks.MockJwtGenerator)

	// Act: Create a new UserHandler instance.
	handler := NewUserHandler(mockUserClient, mockJWTGenerator)

	// Assert: Verify that the handler is properly initialized.
	assert.NotNil(t, handler)
	assert.Equal(t, mockUserClient, handler.userClient, "userClient should match the mock instance")
	assert.Equal(t, mockJWTGenerator, handler.jwtGenerator, "jwtGenerator should match the mock instance")
}
