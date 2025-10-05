// src/handlers/user/test/sign_up_handler_test.go
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
)

/************ テスト用ヘルパー ************/

func newSignUpRouter(uc *user_mocks.MockUsecase, jwt *mocks.MockJwtGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	hd := h.NewHandler(uc, jwt)
	r.POST("/signup", hd.SignUpHandler())
	return r
}

func postJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
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

/************ 本体テスト ************/

func TestSignUpHandler_AllPaths(t *testing.T) {
	type Req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	t.Run("200 OK - success", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		req := Req{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "Secret_123!",
		}

		// Usecase: 入力DTOが正しく詰め替えられていることをざっくり検証
		argMatcher := mock.MatchedBy(func(in user_interface.SignUpInput) bool {
			return in.Name == req.Name && in.Email == req.Email && in.Password == req.Password
		})
		uc.On("SignUp", mock.Anything, argMatcher).
			Return(&user_interface.SignUpOutput{UserID: 42}, nil)
		jwt.On("GenerateJWT", "42").Return("tok_abc", nil)

		w := postJSON(r, "/signup", req)

		assert.Equal(t, http.StatusOK, w.Code)
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "Authentication successful", got["message"])
		assert.Equal(t, "tok_abc", got["token"])
		uc.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("500 - invalid JSON (bind error at parseRequest)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		// email を number にして JSON バインドエラーを発生させる
		w := func() *httptest.ResponseRecorder {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(`{"name":"A","email":123,"password":"Secret_123!"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			return w
		}()

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		uc.AssertNotCalled(t, "SignUp", mock.Anything, mock.Anything)
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
	})

	t.Run("400 - validation error (ValidateSignUp)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		// わざと不正: 短すぎる名前/パスワードなど（実装側の validator に合わせて調整）
		req := Req{
			Name:     "",                 // 必須
			Email:    "bad_email_format", // 不正フォーマット
			Password: "short",            // 要件未達
		}

		w := postJSON(r, "/signup", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		assert.Equal(t, "invalid input", got["error"])
		_, hasFields := got["fields"] // fields が返ってくることのみ確認（詳細は validator 側の責務）
		assert.True(t, hasFields)
		uc.AssertNotCalled(t, "SignUp", mock.Anything, mock.Anything)
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
	})

	t.Run("409 - conflict from usecase (duplicate email)", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		req := Req{
			Name:     "Bob",
			Email:    "dup@example.com",
			Password: "Secret_123!",
		}

		argMatcher := mock.MatchedBy(func(in user_interface.SignUpInput) bool {
			return in.Email == req.Email
		})
		uc.On("SignUp", mock.Anything, argMatcher).
			Return((*user_interface.SignUpOutput)(nil), apperror.Conflictf("email already exists", nil))

		w := postJSON(r, "/signup", req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"error":"email already exists"}`, w.Body.String())
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
		uc.AssertExpectations(t)
	})

	t.Run("500 - internal from usecase", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		req := Req{
			Name:     "Carol",
			Email:    "carol@example.com",
			Password: "Secret_123!",
		}

		uc.On("SignUp", mock.Anything, mock.Anything).
			Return((*user_interface.SignUpOutput)(nil), apperror.Internalf("internal error", nil))

		w := postJSON(r, "/signup", req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
		jwt.AssertNotCalled(t, "GenerateJWT", mock.Anything)
		uc.AssertExpectations(t)
	})

	t.Run("400 - JWT generation failure", func(t *testing.T) {
		uc := new(user_mocks.MockUsecase)
		jwt := &mocks.MockJwtGenerator{}
		r := newSignUpRouter(uc, jwt)

		req := Req{
			Name:     "Dave",
			Email:    "dave@example.com",
			Password: "Secret_123!",
		}

		uc.On("SignUp", mock.Anything, mock.Anything).
			Return(&user_interface.SignUpOutput{UserID: 777}, nil)
		jwt.On("GenerateJWT", "777").Return("", assert.AnError)

		w := postJSON(r, "/signup", req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// ここはハンドラー側が Validation で「Failed to generate token」というメッセージを返す
		var got map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &got)
		// メッセージ文字列はコードに合わせる（固定文言にしている場合は JSONEq でもOK）
		assert.Equal(t, "Failed to generate token", got["error"])
		uc.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})
}
