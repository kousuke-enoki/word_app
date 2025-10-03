package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "word_app/backend/src/handlers/user"
	user_interface "word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/mocks"
	user_mocks "word_app/backend/src/mocks/http/user"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func newSignInRouter(uc *user_mocks.MockUsecase, jwt *mocks.MockJwtGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	hd := h.NewHandler(uc, jwt)
	r.POST("/signin", hd.SignInHandler())
	return r
}

func performJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(http.MethodPost, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestSignInHandler_AllPaths(t *testing.T) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 便利ヘルパー: bcrypt ハッシュを作る
	hash := func(pw string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(b)
	}

	t.Run("200 OK - success", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		req := Req{Email: "alice@example.com", Password: "Secret_123!"}
		uc.On("FindByEmail", mock.Anything, req.Email).
			Return(&user_interface.FindByEmailOutput{
				UserID:         42,
				HashedPassword: hash(req.Password),
			}, nil)
		jwt.On("GenerateJWT", "42").Return("tok_abc", nil)

		w := performJSON(r, "/signin", req)

		assert.Equal(t, http.StatusOK, w.Code)
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "Authentication successful", got["message"])
		assert.Equal(t, "tok_abc", got["token"])
		uc.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("500 - invalid JSON (bind error)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		// 壊れたJSON
		w := func() *httptest.ResponseRecorder {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(`{ "email": 123 }`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			return w
		}()

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		uc.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
	})

	t.Run("400 - validation error (missing fields)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		// 空フィールド → ValidateSignIn が FieldErrors を返す想定
		w := performJSON(r, "/signin", Req{Email: "", Password: ""})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// fields が入ってくる（メッセージはあなたの validator 実装に依存）
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid input", got["error"])
		_, hasFields := got["fields"]
		assert.True(t, hasFields)
		uc.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
	})

	t.Run("404 - user not found", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		req := Req{Email: "unknown@example.com", Password: "x"}
		uc.On("FindByEmail", mock.Anything, req.Email).
			Return((*user_interface.FindByEmailOutput)(nil), apperror.NotFoundf("user not found", nil))

		w := performJSON(r, "/signin", req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"user not found"}`, w.Body.String())
		uc.AssertExpectations(t)
	})

	t.Run("400 - password mismatch", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		req := Req{Email: "bob@example.com", Password: "Wrong_123!"}
		uc.On("FindByEmail", mock.Anything, req.Email).
			Return(&user_interface.FindByEmailOutput{
				UserID:         7,
				HashedPassword: hash("Correct_123!"),
			}, nil)

		w := performJSON(r, "/signin", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// comparePasswordPtrエラー時はメッセージを固定化（情報リーク防止）
		assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
	})

	t.Run("400 - password hash not set (external-only account)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		req := Req{Email: "line-only@example.com", Password: "anything"}
		uc.On("FindByEmail", mock.Anything, req.Email).
			Return(&user_interface.FindByEmailOutput{
				UserID:         100,
				HashedPassword: "", // 未設定
			}, nil)

		w := performJSON(r, "/signin", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
	})

	t.Run("400 - jwt generation failure", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignInRouter(uc, jwt)

		req := Req{Email: "ok@example.com", Password: "Secret_123!"}
		uc.On("FindByEmail", mock.Anything, req.Email).
			Return(&user_interface.FindByEmailOutput{
				UserID:         55,
				HashedPassword: hash(req.Password),
			}, nil)
		jwt.On("GenerateJWT", "55").Return("", assert.AnError)

		w := performJSON(r, "/signin", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
	})
}
