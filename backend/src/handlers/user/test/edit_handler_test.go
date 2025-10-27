// src/handlers/user/test/edit_handler_test.go
package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	h "word_app/backend/src/handlers/user"
	"word_app/backend/src/middleware/jwt"
	user_mocks "word_app/backend/src/mocks/usecase/user"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
	user_usecase "word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/************ テスト用ヘルパー ************/

type roles struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
	IsTest  bool
}

func newEditRouter(uc *user_mocks.MockUsecase, rls *roles) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if rls != nil {
		r.Use(func(c *gin.Context) {
			var uid int
			uid = rls.UserID
			p := models.Principal{
				UserID:  uid,
				IsAdmin: rls.IsAdmin,
				IsRoot:  rls.IsRoot,
				IsTest:  rls.IsTest,
			}
			jwt.SetPrincipal(c, p)
			c.Next()
		})
	}

	hd := h.NewHandler(uc, nil)
	r.PUT("/users/:id", hd.EditHandler())
	return r
}

func doPUT(r *gin.Engine, url string, body any) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(http.MethodPut, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

/************ リクエスト型（元の中置きをそのまま外出し） ************/

type reqPwd struct {
	New     *string `json:"new,omitempty"`
	Current *string `json:"current,omitempty"`
}
type reqBody struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *reqPwd `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
}

/************ 各ケースを個別のテスト関数に分割 ************/

func TestEditHandler_Success(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 10})

	name := "Alice"
	email := "alice@example.com"
	role := "admin"
	body := reqBody{
		Name:  &name,
		Email: &email,
		Role:  &role,
	}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 10 &&
			in.TargetID == 123 &&
			in.Name != nil && *in.Name == "Alice" &&
			in.Email != nil && *in.Email == "alice@example.com" &&
			in.Role != nil && *in.Role == "admin" &&
			in.PasswordNew == nil && in.PasswordCurrent == nil
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return(&models.UserDetail{
			ID:      123,
			Name:    "Alice",
			Email:   &email,
			IsAdmin: true,
			IsRoot:  false,
			IsTest:  false,
		}, nil)

	w := doPUT(r, "/users/123", body)

	assert.Equal(t, http.StatusOK, w.Code)
	var got map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, float64(123), got["id"])
	assert.Equal(t, "Alice", got["name"])
	assert.Equal(t, true, got["isAdmin"])
	assert.Equal(t, false, got["isRoot"])
	assert.Equal(t, false, got["isTest"])
	uc.AssertExpectations(t)
}

func TestEditHandler_Unauthorized_NoRoles(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, nil)

	w := doPUT(r, "/users/1", map[string]any{})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_Unauthorized_IsTest(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 1, IsTest: true})

	w := doPUT(r, "/users/1", map[string]any{})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_BindError_InvalidJSON(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 7})

	w := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/7", bytes.NewBufferString(`{ "email": 123 }`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		return w
	}()

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_InvalidID_NonNumeric(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 2})

	w := doPUT(r, "/users/abc", map[string]any{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_InvalidID_NonPositive(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 2})

	w := doPUT(r, "/users/0", map[string]any{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request"}`, w.Body.String())
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_FormValidation_InvalidRole(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 2})

	role := "owner" // invalid
	body := reqBody{Role: &role}

	w := doPUT(r, "/users/5", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var got map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "invalid input", got["error"])
	_, ok := got["fields"]
	assert.True(t, ok, "fields should be present")
	uc.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestEditHandler_Usecase_NotFound(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 9})

	name := "ValidName"
	body := reqBody{Name: &name}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 9 && in.TargetID == 88 && in.Name != nil && *in.Name == "ValidName"
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return((*models.UserDetail)(nil), apperror.NotFoundf("user not found", nil))

	w := doPUT(r, "/users/88", body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"user not found"}`, w.Body.String())
	uc.AssertExpectations(t)
}

func TestEditHandler_Usecase_Forbidden(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 9})

	name := "ValidName"
	body := reqBody{Name: &name}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 9 && in.TargetID == 77 && in.Name != nil && *in.Name == "ValidName"
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return((*models.UserDetail)(nil), apperror.Forbiddenf("forbidden", nil))

	w := doPUT(r, "/users/77", body)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
	uc.AssertExpectations(t)
}

func TestEditHandler_Usecase_Conflict(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 4})

	email := "dup@example.com"
	body := reqBody{Email: &email}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 4 && in.TargetID == 22 && in.Email != nil && *in.Email == "dup@example.com"
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return((*models.UserDetail)(nil), apperror.Conflictf("email already exists", nil))

	w := doPUT(r, "/users/22", body)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.JSONEq(t, `{"error":"email already exists"}`, w.Body.String())
	uc.AssertExpectations(t)
}

func TestEditHandler_Usecase_InvalidCredential(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 5})

	newPw := "NewPass_1!"
	curPw := "Wrong"
	body := reqBody{Password: &reqPwd{New: &newPw, Current: &curPw}}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 5 && in.TargetID == 5 &&
			in.PasswordNew != nil && *in.PasswordNew == newPw &&
			in.PasswordCurrent != nil && *in.PasswordCurrent == curPw
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return((*models.UserDetail)(nil), apperror.InvalidCredentialf("invalid credential", nil))

	w := doPUT(r, "/users/5", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid credential"}`, w.Body.String())
	uc.AssertExpectations(t)
}

func TestEditHandler_Usecase_Internal(t *testing.T) {
	uc := new(user_mocks.MockUsecase)
	r := newEditRouter(uc, &roles{UserID: 6})

	name := "ValidName"
	body := reqBody{Name: &name}

	argMatcher := mock.MatchedBy(func(in user_usecase.UpdateUserInput) bool {
		return in.EditorID == 6 && in.TargetID == 66 && in.Name != nil && *in.Name == "ValidName"
	})
	uc.On("UpdateUser", mock.Anything, argMatcher).
		Return((*models.UserDetail)(nil), apperror.Internalf("internal error", nil))

	w := doPUT(r, "/users/66", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
	uc.AssertExpectations(t)
}
