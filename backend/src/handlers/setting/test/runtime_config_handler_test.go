// runtime_config_handler_test.go
package setting_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/handlers/setting"
	mockSettinguc "word_app/backend/src/mocks/usecase/setting"
	"word_app/backend/src/test"
	settingUc "word_app/backend/src/usecase/setting"
)

/* ========================================================================== */
/*                    GetRuntimeConfigHandler Tests                           */
/* ========================================================================== */

func TestGetRuntimeConfigHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	successConfig := &settingUc.RuntimeConfigDTO{
		IsTestUserMode:       true,
		IsLineAuthentication: false,
		Version:              "2025-01-06T12:00:00Z",
	}

	tests := []struct {
		name         string
		mockBehavior func(*mockSettinguc.MockSettingFacade)
		wantCode     int
		wantResponse interface{}
		checkHeaders func(*testing.T, *http.Response)
	}{
		{
			name: "正常取得 (200)",
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("GetRuntimeConfig", mock.Anything).
					Return(successConfig, nil)
			},
			wantCode: http.StatusOK,
			wantResponse: map[string]interface{}{
				"is_test_user_mode":      true,
				"is_line_authentication": false,
				"version":                "2025-01-06T12:00:00Z",
			},
			checkHeaders: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, "public, max-age=60, stale-while-revalidate=300", resp.Header.Get("Cache-Control"))
			},
		},
		{
			name: "usecase エラー (500)",
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("GetRuntimeConfig", mock.Anything).
					Return(nil, errors.New("db err"))
			},
			wantCode: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error": "db err",
			},
			checkHeaders: func(t *testing.T, resp *http.Response) {
				// エラー時は Cache-Control ヘッダは設定されない（または設定されてもOK）
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mockSettinguc.NewMockSettingFacade(t)
			tt.mockBehavior(mockUc)

			h := setting.NewHandler(mockUc)

			c, w := test.NewTestCtx("GET", "/public/runtime-config", nil)
			h.GetRuntimeConfigHandler()(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantResponse != nil {
				var got map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)

				expected := tt.wantResponse.(map[string]interface{})
				for key, expectedValue := range expected {
					if key == "version" {
						// version は動的に生成されるため、存在確認のみ
						assert.Contains(t, got, key)
					} else {
						assert.Equal(t, expectedValue, got[key], "key: %s", key)
					}
				}
			}

			if tt.checkHeaders != nil {
				tt.checkHeaders(t, w.Result())
			}

			mockUc.AssertExpectations(t)
		})
	}
}
