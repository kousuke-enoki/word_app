package setting

import (
	"testing"

	settingUc "word_app/backend/src/usecase/setting"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthsettingHandler(t *testing.T) {
	// Arrange: Create mocks for dependencies.
	var mockSettingClient settingUc.SettingFacade

	// Act: Create a new Handler instance.
	handler := NewHandler(mockSettingClient)

	// Assert: Verify that the handler is properly initialized.
	assert.NotNil(t, handler)
	assert.Equal(t, mockSettingClient, handler.settingUsecase, "SettingUsecase should match the mock instance")
}
