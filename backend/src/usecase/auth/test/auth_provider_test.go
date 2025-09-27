package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	mockExt "word_app/backend/src/mocks/infrastructure/repository/auth"
	authUc "word_app/backend/src/usecase/auth"
	"word_app/backend/src/utils/tempjwt"

	mockUser "word_app/backend/src/mocks/infrastructure/repository/user"
)

/* -------------------------------------------------------------------------- */
/*                            Simple hand-made mocks                          */
/* -------------------------------------------------------------------------- */

type providerMock struct{ mock.Mock }

func (m *providerMock) AuthURL(state, nonce string) string {
	return m.Called(state, nonce).String(0)
}

func (m *providerMock) Exchange(ctx context.Context, code string) (*tempjwt.Identity, error) {
	args := m.Called(ctx, code)
	id, _ := args.Get(0).(*tempjwt.Identity)
	return id, args.Error(1)
}

func (m *providerMock) ValidateNonce(idTok *oidc.IDToken, expected string) error { // ★追加
	args := m.Called(idTok, expected)
	return args.Error(0)
}

type jwtMock struct{ mock.Mock }

func (m *jwtMock) GenerateJWT(sub string) (string, error) {
	args := m.Called(sub)
	return args.String(0), args.Error(1)
}

type tempMock struct{ mock.Mock }

func (m *tempMock) GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error) {
	args := m.Called(id, ttl)
	return args.String(0), args.Error(1)
}

func (m *tempMock) ParseTemp(tok string) (*tempjwt.Identity, error) {
	args := m.Called(tok)
	id, _ := args.Get(0).(*tempjwt.Identity)
	return id, args.Error(1)
}

/* -------------------------------------------------------------------------- */
/*                              helper to build UC                            */
/* -------------------------------------------------------------------------- */

func newUC(
	p *providerMock,
	r *mockUser.MockRepository,
	j *jwtMock,
	t *tempMock,
	tHelper testing.TB, // *testing.T を渡すため追加
) *authUc.Usecase {
	ext := mockExt.NewMockExternalAuthRepository(tHelper)

	return authUc.NewUsecase(
		p,   // AuthProvider
		r,   // Repository
		ext, // ExternalAuthRepository
		j,   // JWTGenerator
		t,   // TempTokenGenerator
	)
}

/* ========================================================================== */
/*                                  StartLogin                                */
/* ========================================================================== */

func TestStartLogin(t *testing.T) {
	p := &providerMock{}
	p.On("AuthURL", "st", "no").Return("https://example/auth")
	uc := newUC(p, mockUser.NewMockRepository(t), &jwtMock{}, &tempMock{},
		t)

	got := uc.StartLogin(context.Background(), "st", "no")
	assert.Equal(t, "https://example/auth", got)
	p.AssertExpectations(t)
}

/* ========================================================================== */
/*                               HandleCallback                               */
/* ========================================================================== */

func TestHandleCallback(t *testing.T) {
	ctx := context.Background()
	email := "a@b.com"
	idTok := &tempjwt.Identity{
		Provider: "line",
		Subject:  "sub",
		Email:    &email,
		Name:     "Taro",
	}
	cases := []struct {
		name     string
		setup    func(*providerMock, *mockUser.MockRepository, *jwtMock, *tempMock)
		wantErr  bool
		wantJWT  string
		wantTemp string
		needPass bool
	}{
		{
			name: "Exchange error",
			setup: func(p *providerMock, _ *mockUser.MockRepository, _ *jwtMock, _ *tempMock) {
				p.On("Exchange", ctx, "code").Return(nil, errors.New("x"))
			},
			wantErr: true,
		},
		{
			name: "既存ユーザ → JWT",
			setup: func(p *providerMock, r *mockUser.MockRepository, j *jwtMock, _ *tempMock) {
				p.On("Exchange", ctx, "code").Return(idTok, nil)
				r.On("FindByProvider", ctx, "line", "sub").Return(&domain.User{ID: 123}, nil)
				j.On("GenerateJWT", "123").Return("JWT123", nil)
			},
			wantJWT: "JWT123",
		},
		{
			name: "新規ユーザ → Temp",
			setup: func(p *providerMock, r *mockUser.MockRepository, _ *jwtMock, tmp *tempMock) {
				p.On("Exchange", ctx, "code").Return(idTok, nil)
				r.On("FindByProvider", ctx, "line", "sub").Return(nil, nil)
				tmp.On("GenerateTemp", idTok, mock.AnythingOfType("time.Duration")).Return("TMP42", nil)
			},
			wantTemp: "TMP42", needPass: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &providerMock{}
			r := mockUser.NewMockRepository(t)
			j := &jwtMock{}
			tmp := &tempMock{}
			tc.setup(p, r, j, tmp)

			uc := newUC(p, r, j, tmp, t)
			res, err := uc.HandleCallback(ctx, "code")

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tc.needPass {
				assert.True(t, res.NeedPassword)
				assert.Equal(t, tc.wantTemp, res.TempToken)
			} else {
				assert.Equal(t, tc.wantJWT, res.Token)
			}
		})
	}
}

/* ========================================================================== */
/*                               CompleteSignUp                               */
/* ========================================================================== */

func TestCompleteSignUp(t *testing.T) {
	ctx := context.Background()
	email := "m@x.com"
	idTok := &tempjwt.Identity{
		Provider: "line",
		Subject:  "sub",
		Email:    &email,
		Name:     "Mika",
	}
	cases := []struct {
		name     string
		setup    func(*tempMock, *mockUser.MockRepository, *jwtMock)
		tokenArg string
		wantErr  bool
		wantJWT  string
	}{
		{
			name: "ParseTemp error",
			setup: func(tmp *tempMock, _ *mockUser.MockRepository, _ *jwtMock) {
				tmp.On("ParseTemp", "BAD").Return(nil, errors.New("bad"))
			},
			tokenArg: "BAD", wantErr: true,
		},
		{
			name: "Repo.Create error",
			setup: func(tmp *tempMock, r *mockUser.MockRepository, _ *jwtMock) {
				tmp.On("ParseTemp", "OK").Return(idTok, nil)
				r.On("Create", ctx, mock.Anything, mock.Anything).Return(errors.New("db"))
			},
			tokenArg: "OK", wantErr: true,
		},
		{
			name: "success",
			setup: func(tmp *tempMock, r *mockUser.MockRepository, j *jwtMock) {
				tmp.On("ParseTemp", "OK").Return(idTok, nil)
				r.On("Create", ctx, mock.Anything, mock.Anything).Return(nil)
				j.On("GenerateJWT", mock.AnythingOfType("string")).Return("JWT_OK", nil)
			},
			tokenArg: "OK", wantJWT: "JWT_OK",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &providerMock{}
			r := mockUser.NewMockRepository(t)
			j := &jwtMock{}
			tmp := &tempMock{}
			tc.setup(tmp, r, j)

			uc := newUC(p, r, j, tmp, t)
			pass := "Passw0rd!"
			jwt, err := uc.CompleteSignUp(ctx, tc.tokenArg, &pass)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantJWT, jwt)
			}
		})
	}
}

/* ========================================================================== */
/*                                 StartLogin                                 */
/* ========================================================================== */

func TestProviderAuthURLDelegation(t *testing.T) {
	p := &providerMock{}
	p.On("AuthURL", "s", "n").Return("url")
	uc := newUC(p, mockUser.NewMockRepository(t), &jwtMock{}, &tempMock{}, t)

	assert.Equal(t, "url", uc.StartLogin(context.Background(), "s", "n"))
	p.AssertCalled(t, "AuthURL", "s", "n")
}

/* ========================================================================== */
