package word

import (
	"testing"

	"word_app/backend/src/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewWordHandler(t *testing.T) {
	// Arrange: Create mocks for dependencies.
	mockWordService := new(mocks.WordService)

	// Act: Create a new WordHandler instance.
	handler := NewHandler(mockWordService)

	// Assert: Verify that the handler is properly initialized.
	assert.NotNil(t, handler)
	assert.Equal(t, mockWordService, handler.wordService, "WordService should match the mock instance")
}
