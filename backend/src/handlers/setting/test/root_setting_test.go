// setting_handler_test.go
package setting_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settinghdlr "word_app/backend/src/handlers/setting"
	"word_app/backend/src/middleware/jwt"
	mockSettinguc "word_app/backend/src/mocks/usecase/setting"
	"word_app/backend/src/models"
	"word_app/backend/src/test"
	settingUc "word_app/backend/src/usecase/setting"
)

/* -------------------------------------------------------------------------- */
/*                        helpers (raw ctx when必要)                          */
/* -------------------------------------------------------------------------- */

func newRawCtx(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var rdr *bytes.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req, _ := http.NewRequest(method, path, rdr)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	return c, w
}

/* ========================================================================== */
/*                         GetRootSettingHandler Tests                        */
/* ========================================================================== */

func TestGetRootSettingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// successCfg := &ent.RootConfig{
	// 	ID:                         1,
	// 	EditingPermission:          "admin",
	// 	IsTestUserMode:             false,
	// 	IsEmailAuthenticationCheck: false,
	// 	IsLineAuthentication:       false,
	// }
	successCfg := &domain.RootConfig{
		ID: 1, EditingPermission: "admin",
		IsTestUserMode: false, IsEmailAuthenticationCheck: false,
		IsLineAuthentication: false,
	}
	tests := []struct {
		name         string
		injectUser   bool // true = use test.InjectUser
		isRoot       bool // if injectUser
		needUserID   bool // include userID but NO roles
		mockBehavior func(*mockSettinguc.MockSettingFacade)
		wantCode     int
	}{
		{
			name:       "正常取得 (200)",
			injectUser: true, isRoot: true,
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("GetRoot", mock.Anything, mock.Anything).
					Return(&settingUc.OutputGetRootConfig{Config: successCfg}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:       "userID 無 (401)",
			injectUser: false, needUserID: false,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "userRoles 取得失敗 (401)",
			injectUser: false, needUserID: true,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "root 権限なし (401)",
			injectUser: true, isRoot: false,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "usecase エラー (500)",
			injectUser: true, isRoot: true,
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("GetRoot", mock.Anything, mock.Anything).
					Return(nil, errors.New("db err"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUc := mockSettinguc.NewMockSettingFacade(t)
			tt.mockBehavior(mockUc)

			h := settinghdlr.NewHandler(mockUc)

			c, w := test.NewTestCtx("GET", "/setting/root_config", nil)
			switch {
			case tt.injectUser:
				test.InjectUser(c, 99, tt.isRoot)
			case tt.needUserID:
				// userID only, no roles (IsRoot=false) for testing unauthorized access
				p := models.Principal{
					UserID:  99,
					IsAdmin: false,
					IsRoot:  false,
					IsTest:  false,
				}
				jwt.SetPrincipal(c, p)
			}

			h.GetRootConfigHandler()(c)
			assert.Equal(t, tt.wantCode, w.Code)
			if w.Code == http.StatusOK {
				var got settingUc.OutputGetRootConfig
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				logrus.Info("got:", got)
				assert.Equal(t, successCfg.EditingPermission, got.Config.EditingPermission)
			}
			mockUc.AssertExpectations(t)
		})
	}
}

/* ========================================================================== */
/*                        SaveRootSettingHandler Tests                        */
/* ========================================================================== */

func TestSaveRootSettingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validJSON, _ := json.Marshal(map[string]interface{}{
		"editing_permission":            "admin",
		"is_test_user_mode":             true,
		"is_email_authentication_check": false,
		"is_line_authentication":        false,
	})
	invalidPermJSON, _ := json.Marshal(map[string]interface{}{
		"editing_permission":            "hacker", // ← validation NG
		"is_test_user_mode":             true,
		"is_email_authentication_check": false,
		"is_line_authentication":        false,
	})

	okCfg := &domain.RootConfig{
		ID:                         1,
		EditingPermission:          "admin",
		IsTestUserMode:             true,
		IsEmailAuthenticationCheck: false,
		IsLineAuthentication:       false,
	}

	tests := []struct {
		name         string
		body         []byte
		injectUser   bool
		isRoot       bool
		needUserID   bool // include userID but no roles
		mockBehavior func(*mockSettinguc.MockSettingFacade)
		wantCode     int
	}{
		{
			name:       "正常更新 (200)",
			body:       validJSON,
			injectUser: true, isRoot: true,
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("UpdateRoot", mock.Anything, mock.AnythingOfType("settinguc.InputUpdateRootConfig")).
					Return(okCfg, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:         "userID 無 (401)",
			body:         validJSON,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:         "userRoles 取得失敗 (401)",
			body:         validJSON,
			needUserID:   true,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "root 権限なし (401)",
			body:       validJSON,
			injectUser: true, isRoot: false,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusUnauthorized,
		},
		{
			name:       "BindJSON エラー (400)",
			body:       []byte("{invalid-json"),
			injectUser: true, isRoot: true,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusBadRequest,
		},
		{
			name:       "ValidateRootConfig エラー (400)",
			body:       invalidPermJSON,
			injectUser: true, isRoot: true,
			mockBehavior: func(_ *mockSettinguc.MockSettingFacade) {},
			wantCode:     http.StatusBadRequest,
		},
		{
			name:       "usecase.UpdateRoot エラー (500)",
			body:       validJSON,
			injectUser: true, isRoot: true,
			mockBehavior: func(m *mockSettinguc.MockSettingFacade) {
				m.On("UpdateRoot", mock.Anything, mock.Anything).
					Return(nil, errors.New("db err"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mockUc := &mockSettingUsecase{}
			// tt.mockBehavior(mockUc)

			mockUc := mockSettinguc.NewMockSettingFacade(t)
			tt.mockBehavior(mockUc) // ① 期待呼び出しを登録
			h := settinghdlr.NewHandler(mockUc)

			c, w := newRawCtx("POST", "/setting/root_config", tt.body)
			switch {
			case tt.injectUser:
				test.InjectUser(c, 99, tt.isRoot)
			case tt.needUserID:
				// userID only, no roles (IsRoot=false) for testing unauthorized access
				p := models.Principal{
					UserID:  99,
					IsAdmin: false,
					IsRoot:  false,
					IsTest:  false,
				}
				jwt.SetPrincipal(c, p)
			}

			h.SaveRootConfigHandler()(c)
			assert.Equal(t, tt.wantCode, w.Code)

			if w.Code == http.StatusOK {
				var got domain.RootConfig // または ent.RootConfig
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				assert.Equal(t, okCfg.IsTestUserMode, got.IsTestUserMode)
			}
			mockUc.AssertExpectations(t)
		})
	}
}
