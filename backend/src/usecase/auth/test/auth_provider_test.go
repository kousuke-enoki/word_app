// auth_usecase_test.go
package auth_test

// ★AuthUsecase 構造体を定義しているパッケージ

/* -------------------------------------------------------------------------- */
/*                               モック定義                                   */
/* -------------------------------------------------------------------------- */

// Provider インタフェース（AuthURL / Exchange）用モック
// type providerMock struct{ mock.Mock }

// func (m *providerMock) AuthURL(state, nonce string) string {
// 	args := m.Called(state, nonce)
// 	return args.String(0)
// }

// func (m *providerMock) Exchange(ctx context.Context, code string) (*domain.IDToken, error) {
// 	args := m.Called(ctx, code)
// 	token, _ := args.Get(0).(*domain.IDToken)
// 	return token, args.Error(1)
// }

// // UserRepository (FindByProvider / Create) 用モック
// type userRepoMock struct{ mock.Mock }

// func (m *userRepoMock) FindByProvider(ctx context.Context, provider, subject string) (*domain.User, error) {
// 	args := m.Called(ctx, provider, subject)
// 	user, _ := args.Get(0).(*domain.User)
// 	return user, args.Error(1)
// }

// func (m *userRepoMock) Create(ctx context.Context, u *domain.User, ext *domain.ExternalAuth) error {
// 	return m.Called(ctx, u, ext).Error(0)
// }

// // JWTGenerator (GenerateJWT) 用モック
// type jwtGenMock struct{ mock.Mock }

// func (m *jwtGenMock) GenerateJWT(sub string) (string, error) {
// 	args := m.Called(sub)
// 	return args.String(0), args.Error(1)
// }

// // TempTokenGenerator (GenerateTemp / ParseTemp) 用モック
// type tempGenMock struct{ mock.Mock }

// func (m *tempGenMock) GenerateTemp(id *domain.IDToken, ttl time.Duration) (string, error) {
// 	args := m.Called(id, ttl)
// 	return args.String(0), args.Error(1)
// }
// func (m *tempGenMock) ParseTemp(token string) (*domain.IDToken, error) {
// 	args := m.Called(token)
// 	id, _ := args.Get(0).(*domain.IDToken)
// 	return id, args.Error(1)
// }

// /* -------------------------------------------------------------------------- */
// /*                                  Helper                                    */
// /* -------------------------------------------------------------------------- */

// func newUsecase(p *providerMock, r *userRepoMock, j *jwtGenMock, t *tempGenMock) *auth.AuthUsecase {
// 	return &auth.AuthUsecase{
// 		Provider:     p,
// 		UserRepo:     r,
// 		JwtGenerator: j,
// 		TempJwtGen:   t,
// 	} // フィールド名は実装に合わせてください
// }

// /* -------------------------------------------------------------------------- */
// /*                               StartLogin                                   */
// /* -------------------------------------------------------------------------- */

// func TestStartLogin(t *testing.T) {
// 	p := &providerMock{}
// 	r := &userRepoMock{}
// 	j := &jwtGenMock{}
// 	tmp := &tempGenMock{}

// 	state, nonce := "s123", "n456"
// 	wantURL := "https://example.com/auth?state=s123"

// 	p.On("AuthURL", state, nonce).Return(wantURL)

// 	uc := newUsecase(p, r, j, tmp)
// 	got := uc.StartLogin(context.Background(), state, nonce)

// 	assert.Equal(t, wantURL, got)
// 	p.AssertExpectations(t)
// }

// /* -------------------------------------------------------------------------- */
// /*                             HandleCallback                                 */
// /* -------------------------------------------------------------------------- */

// func TestHandleCallback_AllCases(t *testing.T) {
// 	ctx := context.Background()
// 	baseID := &domain.IDToken{
// 		Provider: "line",
// 		Subject:  "sub123",
// 		Email:    "taro@example.com",
// 		Name:     "Taro",
// 	}

// 	tests := []struct {
// 		name          string
// 		prepare       func(*providerMock, *userRepoMock, *jwtGenMock, *tempGenMock)
// 		wantErr       bool
// 		wantNeedPass  bool
// 		wantToken     string
// 		wantTempToken string
// 	}{
// 		{
// 			name: "既存ユーザ => JWT 発行",
// 			prepare: func(p *providerMock, r *userRepoMock, j *jwtGenMock, tmp *tempGenMock) {
// 				p.On("Exchange", mock.Anything, "code").Return(baseID, nil)
// 				user := &domain.User{ID: 1, Email: baseID.Email, Name: baseID.Name}
// 				r.On("FindByProvider", mock.Anything, "line", "sub123").Return(user, nil)
// 				j.On("GenerateJWT", "1").Return("jwt_token", nil)
// 			},
// 			wantNeedPass: false,
// 			wantToken:    "jwt_token",
// 		},
// 		{
// 			name: "新規ユーザ => TempToken + NeedPassword",
// 			prepare: func(p *providerMock, r *userRepoMock, j *jwtGenMock, tmp *tempGenMock) {
// 				p.On("Exchange", mock.Anything, "code").Return(baseID, nil)
// 				r.On("FindByProvider", mock.Anything, "line", "sub123").Return(nil, nil)
// 				tmp.On("GenerateTemp", baseID, mock.AnythingOfType("time.Duration")).Return("TMP123", nil)
// 			},
// 			wantNeedPass:  true,
// 			wantTempToken: "TMP123",
// 		},
// 		{
// 			name: "Exchange エラー",
// 			prepare: func(p *providerMock, r *userRepoMock, j *jwtGenMock, tmp *tempGenMock) {
// 				p.On("Exchange", mock.Anything, "code").Return(nil, errors.New("exchange NG"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		p := &providerMock{}
// 		r := &userRepoMock{}
// 		j := &jwtGenMock{}
// 		tmp := &tempGenMock{}
// 		tt.prepare(p, r, j, tmp)

// 		uc := newUsecase(p, r, j, tmp)
// 		res, err := uc.HandleCallback(ctx, "code", "state", "nonce")

// 		if tt.wantErr {
// 			assert.Error(t, err, tt.name)
// 			continue
// 		}
// 		assert.NoError(t, err, tt.name)
// 		assert.Equal(t, tt.wantNeedPass, res.NeedPassword, tt.name)
// 		if tt.wantNeedPass {
// 			assert.Equal(t, tt.wantTempToken, res.TempToken, tt.name)
// 			assert.Equal(t, baseID.Email, res.SuggestedMail, tt.name)
// 		} else {
// 			assert.Equal(t, tt.wantToken, res.Token, tt.name)
// 		}
// 		p.AssertExpectations(t)
// 		r.AssertExpectations(t)
// 		j.AssertExpectations(t)
// 		tmp.AssertExpectations(t)
// 	}
// }

// /* -------------------------------------------------------------------------- */
// /*                             CompleteSignUp                                 */
// /* -------------------------------------------------------------------------- */

// func TestCompleteSignUp_AllCases(t *testing.T) {
// 	ctx := context.Background()
// 	idToken := &domain.IDToken{
// 		Provider: "line", Subject: "sub123",
// 		Email: "taro@example.com", Name: "Taro",
// 	}

// 	tests := []struct {
// 		name      string
// 		pass      string
// 		prepare   func(*tempGenMock, *userRepoMock, *jwtGenMock)
// 		wantErr   bool
// 		wantToken string
// 	}{
// 		{
// 			name: "ParseTemp エラー",
// 			pass: "Password1!",
// 			prepare: func(tmp *tempGenMock, repo *userRepoMock, j *jwtGenMock) {
// 				tmp.On("ParseTemp", "BAD").Return(nil, errors.New("parse error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "ユーザ作成バリデーション NG",
// 			pass: "bad", // ← domain.NewUser が弾く想定
// 			prepare: func(tmp *tempGenMock, repo *userRepoMock, j *jwtGenMock) {
// 				tmp.On("ParseTemp", "TOKEN").Return(idToken, nil)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "DB 作成エラー",
// 			pass: "Password1!",
// 			prepare: func(tmp *tempGenMock, repo *userRepoMock, j *jwtGenMock) {
// 				tmp.On("ParseTemp", "TOKEN").Return(idToken, nil)
// 				repo.On("Create", mock.Anything, mock.Anything, mock.Anything).
// 					Return(errors.New("db error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "正常完了 => JWT",
// 			pass: "Password1!",
// 			prepare: func(tmp *tempGenMock, repo *userRepoMock, j *jwtGenMock) {
// 				tmp.On("ParseTemp", "TOKEN").Return(idToken, nil)
// 				repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 				j.On("GenerateJWT", mock.AnythingOfType("string")).Return("JWT_OK", nil)
// 			},
// 			wantToken: "JWT_OK",
// 		},
// 	}

// 	for _, tt := range tests {
// 		p := &providerMock{} // 未使用
// 		tmp := &tempGenMock{}
// 		repo := &userRepoMock{}
// 		j := &jwtGenMock{}

// 		tt.prepare(tmp, repo, j)

// 		uc := newUsecase(p, repo, j, tmp)
// 		token, err := uc.CompleteSignUp(ctx, "TOKEN", tt.pass)

// 		if tt.wantErr {
// 			assert.Error(t, err, tt.name)
// 			continue
// 		}
// 		assert.NoError(t, err, tt.name)
// 		assert.Equal(t, tt.wantToken, token, tt.name)
// 		tmp.AssertExpectations(t)
// 		repo.AssertExpectations(t)
// 		j.AssertExpectations(t)
// 	}
// }
