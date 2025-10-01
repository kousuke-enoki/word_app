package user

import (
	"testing"

	"word_app/backend/src/mocks"
	"word_app/backend/src/mocks/http/user"

	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	// Arrange: Create mocks for dependencies.
	mockUserUsecase := new(user.MockUsecase)
	mockJWTGenerator := new(mocks.MockJwtGenerator)

	// Act: Create a new UserHandler instance.
	handler := NewHandler(mockUserUsecase, mockJWTGenerator)

	// Assert: Verify that the handler is properly initialized.
	assert.NotNil(t, handler)
	assert.Equal(t, mockUserUsecase, handler.userUsecase, "userUsecase should match the mock instance")
	assert.Equal(t, mockJWTGenerator, handler.jwtGenerator, "jwtGenerator should match the mock instance")
}
