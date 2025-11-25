package bulk_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/config"
	h "word_app/backend/src/handlers/bulk"
	bulk_mocks "word_app/backend/src/mocks/usecase/bulk"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	validators.Init()
	binding.Validator = &validators.GinValidator{Validate: validators.V}
}

func TestRegisterHandler_AllPaths(t *testing.T) {
	type Req struct {
		Words []string `json:"words"`
	}

	limits := &config.LimitsCfg{
		BulkRegisterMaxItems: 200,
	}

	t.Run("200 OK - success with all registered", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple", "banana", "cherry"}}
		reqBody, _ := json.Marshal(req)

		ru.On("Register", mock.Anything, 1, req.Words).
			Return(&models.BulkRegisterResponse{
				Success: []string{"apple", "banana", "cherry"},
				Failed:  []models.FailedWord{},
			}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkRegisterResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"apple", "banana", "cherry"}, resp.Success)
		assert.Empty(t, resp.Failed)
		ru.AssertExpectations(t)
	})

	t.Run("200 OK - partial success with some failed", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple", "notexists", "cherry", "limitexceeded"}}
		reqBody, _ := json.Marshal(req)

		ru.On("Register", mock.Anything, 1, req.Words).
			Return(&models.BulkRegisterResponse{
				Success: []string{"apple", "cherry"},
				Failed: []models.FailedWord{
					{Word: "notexists", Reason: "not_exists"},
					{Word: "limitexceeded", Reason: "limit_reached"},
				},
			}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkRegisterResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"apple", "cherry"}, resp.Success)
		assert.Len(t, resp.Failed, 2)
		assert.Equal(t, "notexists", resp.Failed[0].Word)
		assert.Equal(t, "not_exists", resp.Failed[0].Reason)
		assert.Equal(t, "limitexceeded", resp.Failed[1].Word)
		assert.Equal(t, "limit_reached", resp.Failed[1].Reason)
		ru.AssertExpectations(t)
	})

	t.Run("400 - invalid JSON format", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		invalidJSON := `{"words": invalid}`

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBufferString(invalidJSON))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("400 - missing request body", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", nil)
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("400 - empty words array", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{}}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("400 - too many words", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		// 201単語（上限200を超える）
		words := make([]string, 201)
		for i := 0; i < 201; i++ {
			words[i] = "word"
		}
		req := Req{Words: words}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("500 - internal server error", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple", "banana"}}
		reqBody, _ := json.Marshal(req)

		ru.On("Register", mock.Anything, 1, req.Words).
			Return(nil, apperror.Internalf("database error", nil))

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
		ru.AssertExpectations(t)
	})

	t.Run("401 - missing principal", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple", "banana"}}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("200 OK - single word success", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple"}}
		reqBody, _ := json.Marshal(req)

		ru.On("Register", mock.Anything, 1, req.Words).
			Return(&models.BulkRegisterResponse{
				Success: []string{"apple"},
				Failed:  []models.FailedWord{},
			}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkRegisterResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"apple"}, resp.Success)
		assert.Empty(t, resp.Failed)
		ru.AssertExpectations(t)
	})

	t.Run("200 OK - mixed failure reasons", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		req := Req{Words: []string{"apple", "banana", "notexists", "already", "limit", "db"}}
		reqBody, _ := json.Marshal(req)

		ru.On("Register", mock.Anything, 1, req.Words).
			Return(&models.BulkRegisterResponse{
				Success: []string{"apple", "banana"},
				Failed: []models.FailedWord{
					{Word: "notexists", Reason: "not_exists"},
					{Word: "already", Reason: "already_registered"},
					{Word: "limit", Reason: "limit_reached"},
					{Word: "db", Reason: "db_error"},
				},
			}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkRegisterResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"apple", "banana"}, resp.Success)
		assert.Len(t, resp.Failed, 4)
		ru.AssertExpectations(t)
	})

	t.Run("400 - invalid request structure (wrong field type)", func(t *testing.T) {
		ru := bulk_mocks.NewMockRegisterUsecase(t)

		// wordsを文字列として送る（配列である必要がある）
		invalidJSON := `{"words": "invalid"}`

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/register", bytes.NewBufferString(invalidJSON))
		httpReq.Header.Set("Content-Type", "application/json")

		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(nil, ru, limits)
		handler.RegisterHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		ru.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything)
	})
}
