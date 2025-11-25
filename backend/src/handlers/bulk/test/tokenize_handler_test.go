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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTokenizeHandler_AllPaths(t *testing.T) {
	type Req struct {
		Text string `json:"text"`
	}

	limits := &config.LimitsCfg{
		BulkMaxBytes: 51200, // 50KB
	}

	t.Run("200 OK - success", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: "Hello world test"}
		reqBody, _ := json.Marshal(req)

		tu.On("Execute", mock.Anything, 1, req.Text).
			Return([]string{"hello", "world"}, []string{}, []string{"test"}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkTokenizeResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"hello", "world"}, resp.Candidates)
		assert.Equal(t, []string{}, resp.Registered)
		assert.Equal(t, []string{"test"}, resp.NotExistWord)
		tu.AssertExpectations(t)
	})

	t.Run("400 - invalid JSON (body parse error)", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		invalidJSON := `{"text": invalid}`

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBufferString(invalidJSON))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid json")
		tu.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("413 - body too large (over 50KB)", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		// 50KB超の大きなテキスト
		largeText := make([]byte, 51201) // 51201 bytes = 50KB + 1
		for i := range largeText {
			largeText[i] = 'a'
		}
		req := Req{Text: string(largeText)}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.Contains(t, w.Body.String(), "text too large")
		tu.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("200 OK - empty text", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: ""}
		reqBody, _ := json.Marshal(req)

		tu.On("Execute", mock.Anything, 1, "").
			Return([]string{}, []string{}, []string{}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkTokenizeResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Empty(t, resp.Candidates)
		assert.Empty(t, resp.Registered)
		assert.Empty(t, resp.NotExistWord)
		tu.AssertExpectations(t)
	})

	t.Run("429 - too many requests (quota exceeded)", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: "Hello world"}
		reqBody, _ := json.Marshal(req)

		tu.On("Execute", mock.Anything, 1, "Hello world").
			Return(nil, nil, nil, apperror.TooManyRequestsf("quota exceeded", nil))

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "quota exceeded")
		tu.AssertExpectations(t)
	})

	t.Run("500 - internal server error", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: "Hello world"}
		reqBody, _ := json.Marshal(req)

		tu.On("Execute", mock.Anything, 1, "Hello world").
			Return(nil, nil, nil, apperror.Internalf("database error", nil))

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
		tu.AssertExpectations(t)
	})

	t.Run("200 OK - mixed candidates, registered, and not exists", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: "apple banana cherry donut elephant"}
		reqBody, _ := json.Marshal(req)

		tu.On("Execute", mock.Anything, 1, req.Text).
			Return(
				[]string{"apple", "banana"}, // candidates
				[]string{"cherry", "donut"}, // registered
				[]string{"elephant"},        // not exists
				nil,
			)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをシミュレート
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq
		p := models.Principal{UserID: 1, IsAdmin: false, IsRoot: false, IsTest: false}
		c.Set("principalKey", p)

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.BulkTokenizeResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, []string{"apple", "banana"}, resp.Candidates)
		assert.Equal(t, []string{"cherry", "donut"}, resp.Registered)
		assert.Equal(t, []string{"elephant"}, resp.NotExistWord)
		tu.AssertExpectations(t)
	})

	t.Run("401 - missing principal", func(t *testing.T) {
		tu := bulk_mocks.NewMockTokenizeUsecase(t)

		req := Req{Text: "Hello world"}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest(http.MethodPost, "/bulk/tokenize", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Principalをセットしない
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		handler := h.NewHandler(tu, nil, limits)
		handler.TokenizeHandler()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		tu.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything, mock.Anything)
	})
}
