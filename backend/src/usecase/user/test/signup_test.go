package user_test

// // --- helper: tx done() 呼び出しの記録 ---
// type doneRecorder struct {
// 	lastCommit bool
// 	calls      int
// 	err        error
// }

// func (d *doneRecorder) fn(commit bool) error {
// 	d.lastCommit = commit
// 	d.calls++
// 	return d.err
// }

// // --- UC をモックで組み立て ---
// func makeUC_SignUp(t *testing.T, tx *txmock.MockManager, ur *usermock.MockRepository, sr *settingmock.MockUserConfigRepository) *uc.UserUsecase {
// 	return uc.NewUserUsecase(
// 		tx,
// 		ur,
// 		sr,
// 		authmock.NewMockExternalAuthRepository(t), // SignUp では使わないが必須引数
// 	)
// }

// func TestUserUsecase_SignUp_WithMocks(t *testing.T) {
// 	ctx := context.Background()

// 	// ========== 成功系 ==========
// 	t.Run("OK: creates user, creates default settings, commits, returns ID", func(t *testing.T) {
// 		tx := txmock.NewMockManager(t)
// 		ur := usermock.NewMockRepository(t)
// 		sr := settingmock.NewMockUserConfigRepository(t)
// 		ucase := makeUC_SignUp(t, tx, ur, sr)

// 		rec := &doneRecorder{}
// 		tx.EXPECT().
// 			Begin(mock.Anything).
// 			Return(ctx, rec.fn, nil).
// 			Once()

// 		in := uc.SignUpInput{
// 			Name:     "Alice",
// 			Email:    "alice@example.com",
// 			Password: "pw123",
// 		}

// 		// Create 呼び出しで組み立てられた domain.User をざっくり検証
// 		ur.EXPECT().
// 			Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
// 				require.Equal(t, "Alice", u.Name)
// 				require.NotNil(t, u.Email)
// 				require.Equal(t, "alice@example.com", *u.Email)
// 				// パスワードの中身（ハッシュ）は domain.NewUser の仕様次第なのでここでは触れない
// 				return true
// 			})).
// 			Return(&domain.User{ID: 42}, nil).
// 			Once()

// 		sr.EXPECT().
// 			CreateDefault(mock.Anything, 42).
// 			Return(nil).
// 			Once()

// 		out, err := ucase.SignUp(ctx, in)
// 		require.NoError(t, err)
// 		require.NotNil(t, out)
// 		require.Equal(t, 42, out.UserID)

// 		// SignUp の実装は defer と明示呼びで done を2回呼ぶ想定
// 		require.Equal(t, true, rec.lastCommit)
// 		require.Equal(t, 2, rec.calls)
// 	})

// 	// ========== 失敗系 ==========
// 	t.Run("NG: tx.Begin returns error -> no repo call", func(t *testing.T) {
// 		tx := txmock.NewMockManager(t)
// 		ur := usermock.NewMockRepository(t)
// 		sr := settingmock.NewMockUserConfigRepository(t)
// 		ucase := makeUC_SignUp(t, tx, ur, sr)

// 		tx.EXPECT().
// 			Begin(mock.Anything).
// 			Return(ctx, nil, errors.New("tx-begin")).
// 			Once()

// 		in := uc.SignUpInput{
// 			Name:     "Bob",
// 			Email:    "bob@example.com",
// 			Password: "pw",
// 		}

// 		out, err := ucase.SignUp(ctx, in)
// 		require.Error(t, err)
// 		require.Contains(t, err.Error(), "tx-begin")
// 		require.Nil(t, out)
// 	})

// 	t.Run("NG: userRepo.Create returns error (duplicate, etc.) -> rollback (commit=false)", func(t *testing.T) {
// 		tx := txmock.NewMockManager(t)
// 		ur := usermock.NewMockRepository(t)
// 		sr := settingmock.NewMockUserConfigRepository(t)
// 		ucase := makeUC_SignUp(t, tx, ur, sr)

// 		rec := &doneRecorder{}
// 		tx.EXPECT().
// 			Begin(mock.Anything).
// 			Return(ctx, rec.fn, nil).
// 			Once()

// 		in := uc.SignUpInput{
// 			Name:     "Carol",
// 			Email:    "dup@example.com",
// 			Password: "pw",
// 		}

// 		ur.EXPECT().
// 			Create(mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(nil, errors.New("duplicate email")).
// 			Once()

// 		// 設定作成は呼ばれない
// 		out, err := ucase.SignUp(ctx, in)
// 		require.Error(t, err)
// 		require.Contains(t, err.Error(), "duplicate")
// 		require.Nil(t, out)

// 		// commit=false で defer の1回のみ
// 		require.Equal(t, false, rec.lastCommit)
// 		require.Equal(t, 1, rec.calls)
// 	})

// 	t.Run("NG: settingRepo.CreateDefault returns error -> rollback (commit=false)", func(t *testing.T) {
// 		tx := txmock.NewMockManager(t)
// 		ur := usermock.NewMockRepository(t)
// 		sr := settingmock.NewMockUserConfigRepository(t)
// 		ucase := makeUC_SignUp(t, tx, ur, sr)

// 		rec := &doneRecorder{}
// 		tx.EXPECT().
// 			Begin(mock.Anything).
// 			Return(ctx, rec.fn, nil).
// 			Once()

// 		in := uc.SignUpInput{
// 			Name:     "Dave",
// 			Email:    "dave@example.com",
// 			Password: "pw",
// 		}

// 		ur.EXPECT().
// 			Create(mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(&domain.User{ID: 7}, nil).
// 			Once()

// 		sr.EXPECT().
// 			CreateDefault(mock.Anything, 7).
// 			Return(errors.New("settings-fail")).
// 			Once()

// 		out, err := ucase.SignUp(ctx, in)
// 		require.Error(t, err)
// 		require.Contains(t, err.Error(), "settings-fail")
// 		require.Nil(t, out)

// 		require.Equal(t, false, rec.lastCommit)
// 		require.Equal(t, 1, rec.calls)
// 	})

// 	t.Run("NG: commit phase done() returns error -> surfaced", func(t *testing.T) {
// 		tx := txmock.NewMockManager(t)
// 		ur := usermock.NewMockRepository(t)
// 		sr := settingmock.NewMockUserConfigRepository(t)
// 		ucase := makeUC_SignUp(t, tx, ur, sr)

// 		rec := &doneRecorder{err: errors.New("commit-fail")}
// 		tx.EXPECT().
// 			Begin(mock.Anything).
// 			Return(ctx, rec.fn, nil).
// 			Once()

// 		in := uc.SignUpInput{
// 			Name:     "Eve",
// 			Email:    "eve@example.com",
// 			Password: "pw",
// 		}

// 		ur.EXPECT().
// 			Create(mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(&domain.User{ID: 99}, nil).
// 			Once()

// 		sr.EXPECT().
// 			CreateDefault(mock.Anything, 99).
// 			Return(nil).
// 			Once()

// 		out, err := ucase.SignUp(ctx, in)
// 		require.Error(t, err)
// 		require.Contains(t, err.Error(), "commit-fail")
// 		require.Nil(t, out)

// 		// 成功フローなので commit=true、done は2回呼ばれる実装
// 		require.Equal(t, true, rec.lastCommit)
// 		require.Equal(t, 2, rec.calls)
// 	})
// }
