package setting_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settinghdlr "word_app/backend/src/handlers/setting"
	mocks "word_app/backend/src/mocks/usecase/setting"
	"word_app/backend/src/test"
	settingUc "word_app/backend/src/usecase/setting"
)

/* -------------------------------------------------------------------------- */
/*                               helper (POST)                                */
/* -------------------------------------------------------------------------- */

func rawCtx(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	rdr := bytes.NewReader(body)
	req, _ := http.NewRequest(method, path, rdr)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	return c, w
}

/* ========================================================================== */
/*                        GetUserSettingHandler Tests                         */
/* ========================================================================== */

func TestGetUserSettingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	successCfg := &domain.UserConfig{ID: 1, UserID: 99, IsDarkMode: true} // フィールドは例
	tests := []struct {
		name         string
		injectUser   bool
		mockBehavior func(*mocks.MockSettingFacade)
		wantCode     int
	}{
		{
			name:       "正常取得 (200)",
			injectUser: true,
			mockBehavior: func(m *mocks.MockSettingFacade) {
				m.On("GetUser", mock.Anything, mock.Anything).
					Return(&settingUc.OutputGetUserConfig{Config: successCfg}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:         "userID 無 (401)",
			injectUser:   false,
			mockBehavior: func(_ *mocks.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "usecase エラー (500)",
			injectUser: true,
			mockBehavior: func(m *mocks.MockSettingFacade) {
				m.On("GetUser", mock.Anything, mock.Anything).
					Return(nil, errors.New("db err"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockSettingFacade(t)
			tt.mockBehavior(mockUc)

			h := settinghdlr.NewHandler(mockUc)
			c, w := test.NewTestCtx("GET", "/setting/user_config", nil)
			if tt.injectUser {
				test.InjectUser(c, 99, false)
			}

			h.GetUserConfigHandler()(c)
			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var got settingUc.OutputGetUserConfig
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				assert.Equal(t, true, got.Config.IsDarkMode)
			}
			mockUc.AssertExpectations(t)
		})
	}
}

/* ========================================================================== */
/*                       SaveUserSettingHandler Tests                         */
/* ========================================================================== */

func TestSaveUserSettingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	successCfg := &domain.UserConfig{ID: 1, UserID: 99, IsDarkMode: false}
	validBody, _ := json.Marshal(map[string]interface{}{"theme": "light"})
	badJSON := []byte("{invalid-json")

	tests := []struct {
		name         string
		body         []byte
		injectUser   bool
		mockBehavior func(*mocks.MockSettingFacade)
		wantCode     int
	}{
		{
			name:       "正常更新 (200)",
			body:       validBody,
			injectUser: true,
			mockBehavior: func(m *mocks.MockSettingFacade) {
				m.On("UpdateUser", mock.Anything, mock.Anything).
					Return(successCfg, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:         "BindJSON エラー (400)",
			body:         badJSON,
			injectUser:   true,
			mockBehavior: func(_ *mocks.MockSettingFacade) {},
			wantCode:     http.StatusBadRequest,
		},
		{
			name:       "usecase エラー (500)",
			body:       validBody,
			injectUser: true,
			mockBehavior: func(m *mocks.MockSettingFacade) {
				m.On("UpdateUser", mock.Anything, mock.Anything).
					Return(nil, errors.New("db ng"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mocks.NewMockSettingFacade(t)
			tt.mockBehavior(mockUc)

			h := settinghdlr.NewHandler(mockUc)
			c, w := rawCtx("POST", "/setting/user_config", tt.body)
			if tt.injectUser {
				test.InjectUser(c, 99, false)
			}

			h.SaveUserConfigHandler()(c)
			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var got domain.UserConfig
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				assert.Equal(t, successCfg.IsDarkMode, got.IsDarkMode)
			}
			mockUc.AssertExpectations(t)
		})
	}
}
