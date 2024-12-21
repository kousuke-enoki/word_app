package word_test

import (
	"encoding/json"
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

func TestAllWordListHandler_success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// モックサービスの初期化
	mockWordService := new(mocks.WordService)

	// モックデータ
	mockWords := []models.Word{
		{
			ID:   1,
			Name: "example",
			WordInfos: []models.WordInfo{
				{
					ID: 1,
					PartOfSpeech: models.PartOfSpeech{
						ID:   1,
						Name: "noun",
					},
					JapaneseMeans: []models.JapaneseMean{
						{
							ID:   1,
							Name: "例",
						},
					},
				},
			},
		},
		{
			ID:   2,
			Name: "test",
			WordInfos: []models.WordInfo{
				{
					ID: 2,
					PartOfSpeech: models.PartOfSpeech{
						ID:   2,
						Name: "verb",
					},
					JapaneseMeans: []models.JapaneseMean{
						{
							ID:   2,
							Name: "テスト",
						},
					},
				},
			},
		},
	}
	total := 2
	page := 1
	limit := 10

	// モックの振る舞いを設定
	mockWordService.On("GetWords", mock.Anything, "", "id", "asc", page, limit).
		Return(mockWords, total, page, nil)

	// ハンドラーを初期化
	wordHandler := word.NewWordHandler(mockWordService)

	// テスト用のリクエストとレスポンス
	req := httptest.NewRequest(http.MethodGet, "/words/all_list?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/words/all_list", wordHandler.AllWordListHandler())

	// ハンドラーを実行
	router.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Words      []models.Word `json:"words"`
		TotalPages int           `json:"totalPages"`
		TotalCount int           `json:"totalCount"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, mockWords, response.Words)
	assert.Equal(t, total, response.TotalCount)
	assert.Equal(t, page, response.TotalPages)

	// モックが期待された呼び出しを受けたかを確認
	mockWordService.AssertExpectations(t)
}
func TestAllWordListHandler_Failure(t *testing.T) {
	// gin.SetMode(gin.TestMode)

	// // モックサービスの初期化
	// mockWordService := new(mocks.WordService)

	// // 失敗用のモックデータ
	// errorMessage := "failed to fetch words"

	// // モックの振る舞いを設定（エラーを返す）
	// mockWordService.On("GetWords", mock.Anything, "", "id", "asc", 1, 10).
	// 	Return(nil, 0, 0, errors.New(errorMessage))

	// // ハンドラーを初期化
	// wordHandler := word.NewWordHandler(mockWordService)

	// // テスト用のリクエストとレスポンス
	// req := httptest.NewRequest(http.MethodGet, "/words/all_list?page=1&limit=10", nil)
	// w := httptest.NewRecorder()
	// router := gin.Default()
	// router.GET("/words/all_list", wordHandler.AllWordListHandler())

	// // ハンドラーを実行
	// router.ServeHTTP(w, req)

	// // レスポンスの検証
	// assert.Equal(t, http.StatusInternalServerError, w.Code)

	// var errorResponse struct {
	// 	Message string `json:"error"`
	// }

	// err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	// assert.NoError(t, err)
	// assert.Equal(t, errorMessage, errorResponse.Message)

	// // モックが期待された呼び出しを受けたかを確認
	// mockWordService.AssertExpectations(t)
}

func TestAllWordListHandler_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// ハンドラーを初期化
	mockWordService := new(mocks.WordService)
	wordHandler := word.NewWordHandler(mockWordService)

	// 不正なパラメータを含むリクエスト
	req := httptest.NewRequest(http.MethodGet, "/words/all_list?page=0&limit=0", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/words/all_list", wordHandler.AllWordListHandler())

	// ハンドラーを実行
	router.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse struct {
		Message string `json:"error"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid query parameters", errorResponse.Message)
}
