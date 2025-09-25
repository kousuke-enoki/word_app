// handler_test.go
package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_handler "word_app/backend/src/handlers/auth"
	"word_app/backend/src/interfaces/http/auth"
	auth_mock "word_app/backend/src/mocks/http/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthLineHandler(t *testing.T) {
	t.Run("TestLineLogin", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		tests := []struct {
			name        string
			redirectURL string
			wantHeader  string
		}{
			{"redirect_success", "https://line.me/oauth", "https://line.me/oauth"},
			{"redirect_empty", "", "/line/"}, // 修正
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockUC := new(auth_mock.MockUsecase)
				mockUC.
					On("StartLogin", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(tt.redirectURL)

				mockJWTGen := new(auth_mock.MockJWTGenerator)
				h := auth_handler.NewHandler(mockUC, mockJWTGen)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				req := httptest.NewRequest(http.MethodGet, "/line/login", nil)
				c.Request = req

				h.LineLogin()(c)

				assert.Equal(t, http.StatusFound, w.Code)
				assert.Equal(t, tt.wantHeader, w.Header().Get("Location"))
				mockUC.AssertExpectations(t)
			})
		}
	})

	t.Run("TestLineCallback", func(t *testing.T) {
		// func TestLineCallback(t *testing.T) {
		gin.SetMode(gin.TestMode)

		successRes := &auth.CallbackResult{Token: "jwt123"}
		usecaseErr := errors.New("db error")

		tests := []struct {
			name           string
			query          string
			mockReturn     *auth.CallbackResult
			mockErr        error
			wantStatusCode int
			wantContains   string
		}{
			{
				name:           "success",
				query:          "?code=abc&state=xyz",
				mockReturn:     successRes,
				mockErr:        nil,
				wantStatusCode: http.StatusOK,
				wantContains:   `"token":"jwt123"`,
			},
			{
				name:           "handle_callback_error",
				query:          "?code=abc&state=xyz",
				mockReturn:     nil,
				mockErr:        usecaseErr,
				wantStatusCode: http.StatusInternalServerError,
				wantContains:   `error`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockUC := new(auth_mock.MockUsecase)
				mockUC.
					On("HandleCallback", mock.Anything, "abc", mock.Anything).
					Return(tt.mockReturn, tt.mockErr)
				mockJWTGenerator := new(auth_mock.MockJWTGenerator)

				userHandler := auth_handler.NewHandler(mockUC, mockJWTGenerator)
				// h := &auth.Handler{Usecase: mockUC, jwtGenerator: mockJWTGenerator}

				//  := &auth.Handler{Usecase: mockUC}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				req := httptest.NewRequest(http.MethodGet, "/callback"+tt.query, nil)
				c.Request = req

				userHandler.LineCallback()(c)

				assert.Equal(t, tt.wantStatusCode, w.Code)
				assert.Contains(t, w.Body.String(), tt.wantContains)
				mockUC.AssertExpectations(t)
			})
		}
	})

	t.Run("TestLineComplete", func(t *testing.T) {
		// func TestLineComplete(t *testing.T) {
		gin.SetMode(gin.TestMode)

		type body struct {
			TempToken string `json:"temp_token"`
			Password  string `json:"password"`
		}
		validBody, _ := json.Marshal(body{TempToken: "tmp123", Password: "pass1"})
		invalidJSON := []byte(`{"temp_token":1`) // bind エラー用

		completeErr := errors.New("complete signup failed")

		tests := []struct {
			name           string
			requestBody    []byte
			mockJWT        string
			mockErr        error
			wantStatusCode int
			wantContains   string
		}{
			{
				name:           "success",
				requestBody:    validBody,
				mockJWT:        "jwt456",
				mockErr:        nil,
				wantStatusCode: http.StatusOK,
				wantContains:   `"token":"jwt456"`,
			},
			{
				name:           "json_bind_error",
				requestBody:    invalidJSON,
				wantStatusCode: http.StatusBadRequest,
				wantContains:   `error`,
			},
			{
				name:           "complete_signup_error",
				requestBody:    validBody,
				mockJWT:        "",
				mockErr:        completeErr,
				wantStatusCode: http.StatusInternalServerError,
				wantContains:   `error`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockUC := new(auth_mock.MockUsecase)
				// Bind エラーケースでは CompleteSignUp は呼ばれない想定
				if tt.name != "json_bind_error" {
					pass := "pass1"
					mockUC.
						On("CompleteSignUp", mock.Anything, "tmp123", &pass).
						Return(tt.mockJWT, tt.mockErr)
				}
				mockJWTGenerator := new(auth_mock.MockJWTGenerator)

				userHandler := auth_handler.NewHandler(mockUC, mockJWTGenerator)
				// h := &auth.Handler{Usecase: mockUC}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				req := httptest.NewRequest(http.MethodPost, "/complete", bytes.NewReader(tt.requestBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req

				userHandler.LineComplete()(c)

				assert.Equal(t, tt.wantStatusCode, w.Code)
				assert.Contains(t, w.Body.String(), tt.wantContains)
				mockUC.AssertExpectations(t)
			})
		}
	})
}
