package word_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"word_app/backend/src/handlers/word"
	"word_app/backend/src/mocks"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWordShowHandler_Success(t *testing.T) {
	mockService := new(mocks.WordService)
	handler := word.NewWordHandler(mockService)

	// モックの期待値と戻り値を設定
	mockService.On("GetWordDetails", mock.Anything, 1).Return(&models.WordResponse{
		Name: "test",
	}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/words/:id", handler.WordShowHandler())

	req := httptest.NewRequest(http.MethodGet, "/words/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test")

	// モックの呼び出しが期待通りであることを検証
	mockService.AssertExpectations(t)
}

func TestWordShowHandler_Error(t *testing.T) {
	mockService := new(mocks.WordService)
	handler := word.NewWordHandler(mockService)

	// モックの期待値と戻り値を設定
	mockService.On("GetWordDetails", mock.Anything, 1).Return(nil, errors.New("failed to fetch word details"))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/words/:id", handler.WordShowHandler())

	req := httptest.NewRequest(http.MethodGet, "/words/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to fetch word details")

	// モックの呼び出しが期待通りであることを検証
	mockService.AssertExpectations(t)
}
