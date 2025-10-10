package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	authUc "word_app/backend/src/usecase/auth"
	"word_app/backend/src/utils/tempjwt"

	mockLine "word_app/backend/src/mocks/infrastructure/auth/line"
	mockJwt "word_app/backend/src/mocks/infrastructure/jwt"
	mockAuth "word_app/backend/src/mocks/infrastructure/repository/auth"
	mockSetting "word_app/backend/src/mocks/infrastructure/repository/setting"
	mockTx "word_app/backend/src/mocks/infrastructure/repository/tx"
	mockUser "word_app/backend/src/mocks/infrastructure/repository/user"
	mockudu "word_app/backend/src/mocks/infrastructure/repository/userdailyusage"
)

/* -------------------------------------------------------------------------- */
/*                              helper to build UC                            */
/* -------------------------------------------------------------------------- */

type mockClock struct{ mock.Mock }

func (m *mockClock) Now() time.Time {
	return time.Now()
}

func newUC(
	txm *mockTx.MockManager,
	p *mockLine.MockProvider,
	r *mockUser.MockRepository,
	s *mockSetting.MockUserConfigRepository,
	a *mockAuth.MockExternalAuthRepository,
	j *mockJwt.MockJWTGenerator,
	t *mockJwt.MockTempTokenGenerator,
	rs *mockSetting.MockRootConfigRepository,
	udu *mockudu.MockRepository,
	c *mockClock,
	tHelper testing.TB, // *testing.T を渡すため追加
) *authUc.AuthUsecase {
	ext := mockAuth.NewMockExternalAuthRepository(tHelper)

	return authUc.NewUsecase(
		txm,
		p,   // AuthProvider
		r,   // UserRepository
		s,   // UserConfigRepository
		ext, // ExternalAuthRepository
		j,   // JWTGenerator
		t,   // TempTokenGenerator
		rs,  // rootSettingRepo
		udu, // userDailyUsageRepo
		c,   // clock
	)
}

/* ========================================================================== */
/*                                  StartLogin                                */
/* ========================================================================== */

func TestStartLogin(t *testing.T) {
	// p := &providerMock{}
	p := new(mockLine.MockProvider)
	mockTx := new(mockTx.MockManager)
	mockUserSetting := new(mockSetting.MockUserConfigRepository)
	mockAuth := new(mockAuth.MockExternalAuthRepository)
	mockJwtg := new(mockJwt.MockJWTGenerator)
	mockTempJwt := new(mockJwt.MockTempTokenGenerator)
	mockRootSetting := new(mockSetting.MockRootConfigRepository)
	mockUdu := new(mockudu.MockRepository)
	c := &mockClock{}

	p.On("AuthURL", "st", "no").Return("https://example/auth")
	uc := newUC(mockTx, p, mockUser.NewMockRepository(t), mockUserSetting, mockAuth, mockJwtg, mockTempJwt,
		mockRootSetting, mockUdu, c, t)

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
		name  string
		setup func(p *mockLine.MockProvider, r *mockUser.MockRepository,
			j *mockJwt.MockJWTGenerator, tmp *mockJwt.MockTempTokenGenerator)

		wantErr  bool
		wantJWT  string
		wantTemp string
		needPass bool
	}{
		{
			name: "Exchange error",
			setup: func(p *mockLine.MockProvider, _ *mockUser.MockRepository,
				_ *mockJwt.MockJWTGenerator, _ *mockJwt.MockTempTokenGenerator,
			) {
				p.On("Exchange", mock.Anything, "code").Return(nil, errors.New("x"))
			},
			wantErr: true,
		},
		{
			name: "既存ユーザ → JWT",
			setup: func(p *mockLine.MockProvider, r *mockUser.MockRepository,
				j *mockJwt.MockJWTGenerator, _ *mockJwt.MockTempTokenGenerator,
			) {
				p.On("Exchange", mock.Anything, "code").Return(idTok, nil)
				r.On("FindByProvider", mock.Anything, "line", "sub").
					Return(&domain.User{ID: 123}, nil)
				j.On("GenerateJWT", "123").Return("JWT123", nil)
			},
			wantJWT: "JWT123",
		},
		{
			name: "新規ユーザ → Temp",
			setup: func(p *mockLine.MockProvider, r *mockUser.MockRepository,
				_ *mockJwt.MockJWTGenerator, tmp *mockJwt.MockTempTokenGenerator,
			) {
				p.On("Exchange", mock.Anything, "code").Return(idTok, nil)
				r.On("FindByProvider", mock.Anything, "line", "sub").Return(nil, nil)
				tmp.On("GenerateTemp", idTok, mock.AnythingOfType("time.Duration")).
					Return("TMP42", nil)
			},
			wantTemp: "TMP42", needPass: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := new(mockLine.MockProvider)
			r := mockUser.NewMockRepository(t)
			mockTx := new(mockTx.MockManager)
			mockUserSetting := new(mockSetting.MockUserConfigRepository)
			mockAuth := new(mockAuth.MockExternalAuthRepository)
			j := new(mockJwt.MockJWTGenerator)
			tmp := new(mockJwt.MockTempTokenGenerator)
			mockRootSetting := new(mockSetting.MockRootConfigRepository)
			mockUdu := new(mockudu.MockRepository)
			c := &mockClock{}

			tc.setup(p, r, j, tmp)
			uc := newUC(mockTx, p, r, mockUserSetting, mockAuth, j, tmp,
				mockRootSetting, mockUdu, c, t)
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

// func TestCompleteSignUp(t *testing.T) {
// 	ctx := context.Background()
// 	email := "m@x.com"
// 	idTok := &tempjwt.Identity{
// 		Provider: "line",
// 		Subject:  "sub",
// 		Email:    &email,
// 		Name:     "Mika",
// 	}
// 	cases := []struct {
// 		name     string
// 		setup    func(*tempMock, *mockUser.MockRepository, *jwtMock)
// 		tokenArg string
// 		wantErr  bool
// 		wantJWT  string
// 	}{
// 		{
// 			name: "ParseTemp error",
// 			setup: func(tmp *tempMock, _ *mockUser.MockRepository, _ *jwtMock) {
// 				tmp.On("ParseTemp", "BAD").Return(nil, errors.New("bad"))
// 			},
// 			tokenArg: "BAD", wantErr: true,
// 		},
// 		{
// 			name: "Repo.Create error",
// 			setup: func(tmp *tempMock, r *mockUser.MockRepository, _ *jwtMock) {
// 				tmp.On("ParseTemp", "OK").Return(idTok, nil)
// 				r.On("Create", ctx, mock.Anything, mock.Anything).Return(errors.New("db"))
// 			},
// 			tokenArg: "OK", wantErr: true,
// 		},
// 		{
// 			name: "success",
// 			setup: func(tmp *tempMock, r *mockUser.MockRepository, j *jwtMock) {
// 				tmp.On("ParseTemp", "OK").Return(idTok, nil)
// 				r.On("Create", ctx, mock.Anything, mock.Anything).Return(nil)
// 				j.On("GenerateJWT", mock.AnythingOfType("string")).Return("JWT_OK", nil)
// 			},
// 			tokenArg: "OK", wantJWT: "JWT_OK",
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			p := &providerMock{}
// 			r := mockUser.NewMockRepository(t)
// 			j := &jwtMock{}
// 			tmp := &tempMock{}
// 			mockTx := new(mockTx.MockManager)
// 			mockSetting := new(mockSetting.MockUserConfigRepository)
// 			tc.setup(tmp, r, j)

// 			uc := newUC(mockTx, p, r, mockSetting, j, tmp, t)
// 			pass := "Passw0rd!"
// 			jwt, err := uc.CompleteSignUp(ctx, tc.tokenArg, &pass)

// 			if tc.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tc.wantJWT, jwt)
// 			}
// 		})
// 	}
// }

/* ========================================================================== */
/*                                 StartLogin                                 */
/* ========================================================================== */

func TestProviderAuthURLDelegation(t *testing.T) {
	p := new(mockLine.MockProvider)
	mockTx := new(mockTx.MockManager)
	mockUserSetting := new(mockSetting.MockUserConfigRepository)
	mockAuth := new(mockAuth.MockExternalAuthRepository)
	mockJwtg := new(mockJwt.MockJWTGenerator)
	mockTempJwt := new(mockJwt.MockTempTokenGenerator)
	mockRootSetting := new(mockSetting.MockRootConfigRepository)
	mockUdu := new(mockudu.MockRepository)
	c := &mockClock{}

	p.On("AuthURL", "s", "n").Return("url")
	uc := newUC(mockTx, p, mockUser.NewMockRepository(t), mockUserSetting,
		mockAuth, mockJwtg, mockTempJwt, mockRootSetting, mockUdu, c, t)

	assert.Equal(t, "url", uc.StartLogin(context.Background(), "s", "n"))
	p.AssertCalled(t, "AuthURL", "s", "n")
}

/* ========================================================================== */
