// backend/src/usecase/user/test/delete_with_mocks_test.go
package user_test

// func TestUserUsecase_Delete_WithMocks(t *testing.T) {
// 	type ids struct {
// 		editor int
// 		target int
// 	}
// 	type actors struct {
// 		editor *domain.User
// 		target *domain.User
// 	}
// 	type fails struct {
// 		txBegin error
// 		doneErr error
// 		authErr error
// 		setErr  error
// 		userErr error
// 	}

// 	type expect struct {
// 		ok                 bool
// 		commitTrue         bool
// 		expectForbidden    bool
// 		expectNotFoundLike bool // FindByID エラー想定時
// 		checkSameTime      bool // 3つの削除に同一時刻が渡されるか
// 	}

// 	run := func(name string, id ids, a actors, f fails, e expect) {
// 		t.Run(name, func(t *testing.T) {
// 			ctx := context.Background()

// 			// --- モックの用意 ---
// 			tx := txmock.NewMockManager(t)
// 			userRepo := usermock.NewMockRepository(t)
// 			authRepo := authmock.NewMockExternalAuthRepository(t)
// 			setRepo := settingmock.NewMockUserConfigRepository(t)

// 			// Begin の戻りの done(commit) をキャプチャして検証
// 			var capturedCommit *bool
// 			tx.On("Begin", mock.Anything).Return(
// 				ctx,
// 				func(commit bool) error {
// 					capturedCommit = &commit
// 					return f.doneErr
// 				},
// 				f.txBegin,
// 			).Once()

// 			// txBegin で即エラー帰るパス
// 			if f.txBegin != nil {
// 				// userRepo := usermock.NewMockRepository(t)
// 				// authRepo := authmock.NewMockExternalAuthRepository(t)
// 				// setRepo := settingmock.NewMockUserConfigRepository(t)
// 				// tx := txmock.NewMockManager(t)
// 				ucase := uc.NewUserUsecase(tx, userRepo, setRepo, authRepo)

// 				// ucase := uc.NewUserUsecase(tx, userRepo, setRepo, authRepo)
// 				err := ucase.Delete(ctx, uc.DeleteUserInput{EditorID: id.editor, TargetID: id.target})
// 				require.Error(t, err)
// 				if capturedCommit != nil {
// 					// Begin が失敗しているので通常は done は呼ばれない想定
// 					require.False(t, *capturedCommit)
// 				}
// 				return
// 			}

// 			// --- FindByID の期待値 ---
// 			// editor
// 			if a.editor != nil {
// 				userRepo.EXPECT().
// 					FindByID(mock.Anything, id.editor).
// 					Return(&domain.User{ID: a.editor.ID, IsRoot: a.editor.IsRoot}, nil).
// 					Once()
// 			} else {
// 				userRepo.EXPECT().
// 					FindByID(mock.Anything, id.editor).
// 					Return(nil, errors.New("editor not found")).
// 					Once()
// 			}
// 			// target（editor が見つからずに return の実装なら、呼ばれないが Once にしておいてOK）
// 			if a.target != nil {
// 				userRepo.EXPECT().
// 					FindByID(mock.Anything, id.target).
// 					Return(&domain.User{ID: a.target.ID, IsRoot: a.target.IsRoot}, nil).
// 					Maybe() // editor 側で失敗するケースもあるので Maybe にしておく
// 			} else {
// 				userRepo.EXPECT().
// 					FindByID(mock.Anything, id.target).
// 					Return(nil, errors.New("target not found")).
// 					Maybe()
// 			}

// 			// --- 関連削除の呼び出し（成功系や後段の失敗系でのみ期待）---
// 			var t0 time.Time
// 			authCall := authRepo.EXPECT().
// 				SoftDeleteByUserID(mock.Anything, id.target, mock.AnythingOfType("time.Time"))
// 			authCall.Run(func(ctx context.Context, userID int, tt time.Time) {
// 				t0 = tt
// 			})
// 			authCall.Return(f.authErr).Maybe()

// 			setRepo.EXPECT().
// 				SoftDeleteByUserID(mock.Anything, id.target,
// 					mock.MatchedBy(func(tt time.Time) bool {
// 						if e.checkSameTime && !t0.IsZero() {
// 							return tt.Equal(t0)
// 						}
// 						return true
// 					}),
// 				).
// 				Return(f.setErr).Maybe()

// 			userRepo.EXPECT().
// 				SoftDeleteByID(mock.Anything, id.target,
// 					mock.MatchedBy(func(tt time.Time) bool {
// 						if e.checkSameTime && !t0.IsZero() {
// 							return tt.Equal(t0)
// 						}
// 						return true
// 					}),
// 				).
// 				Return(f.userErr).Maybe()

// 			// --- 実行 ---
// 			// userRepo = usermock.NewMockRepository(t)
// 			// authRepo = authmock.NewMockExternalAuthRepository(t)
// 			// setRepo = settingmock.NewMockUserConfigRepository(t)
// 			// tx = txmock.NewMockManager(t)
// 			ucase := uc.NewUserUsecase(tx, userRepo, setRepo, authRepo)
// 			err := ucase.Delete(ctx, uc.DeleteUserInput{EditorID: id.editor, TargetID: id.target})

// 			// --- 検証 ---
// 			if e.ok {
// 				require.NoError(t, err)
// 			} else {
// 				require.Error(t, err)
// 				if e.expectForbidden {
// 					require.Contains(t, strings.ToLower(err.Error()), "forbidden")
// 				}
// 				if e.expectNotFoundLike {
// 					require.Contains(t, err.Error(), "not found")
// 				}
// 			}

// 			// commit フラグの検証
// 			if capturedCommit != nil {
// 				require.Equal(t, e.commitTrue, *capturedCommit)
// 			}

// 			// done エラーケースは、成功フローの末尾で commit=true のまま done でエラー返す
// 			if f.doneErr != nil {
// 				require.Error(t, err)
// 				require.Contains(t, err.Error(), f.doneErr.Error())
// 			}
// 		})
// 	}

// 	// ---------------- 成功系 ----------------
// 	// root が root でない他人を削除
// 	run("OK_root_deletes_non_root_other_user",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{},
// 		expect{ok: true, commitTrue: true, checkSameTime: true},
// 	)

// 	// 非rootが自分自身を削除
// 	run("OK_non_root_deletes_self",
// 		ids{editor: 10, target: 10},
// 		actors{
// 			editor: &domain.User{ID: 10, IsRoot: false},
// 			target: &domain.User{ID: 10, IsRoot: false},
// 		},
// 		fails{},
// 		expect{ok: true, commitTrue: true, checkSameTime: true},
// 	)

// 	// ---------------- 失敗系：Tx開始 ----------------
// 	run("NG_tx_begin_failed",
// 		ids{editor: 1, target: 2},
// 		actors{}, // 参照されない
// 		fails{txBegin: errors.New("tx-begin")},
// 		expect{ok: false, commitTrue: false},
// 	)

// 	// ---------------- 失敗系：FindByID ----------------
// 	run("NG_editor_not_found",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: nil,
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{},
// 		expect{ok: false, commitTrue: false, expectNotFoundLike: true},
// 	)

// 	run("NG_target_not_found",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: nil,
// 		},
// 		fails{},
// 		expect{ok: false, commitTrue: false, expectNotFoundLike: true},
// 	)

// 	// ---------------- 失敗系：ポリシー（authorize） ----------------
// 	run("NG_root_cannot_delete_self",
// 		ids{editor: 1, target: 1},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 1, IsRoot: false},
// 		},
// 		fails{},
// 		expect{ok: false, commitTrue: false, expectForbidden: true},
// 	)

// 	run("NG_root_cannot_delete_another_root",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: true},
// 		},
// 		fails{},
// 		expect{ok: false, commitTrue: false, expectForbidden: true},
// 	)

// 	run("NG_non_root_cannot_delete_others",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: false},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{},
// 		expect{ok: false, commitTrue: false, expectForbidden: true},
// 	)

// 	// ---------------- 失敗系：関連削除のどこかで失敗 ----------------
// 	run("NG_authRepo_error",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{authErr: errors.New("auth-fail")},
// 		expect{ok: false, commitTrue: false},
// 	)

// 	run("NG_settingRepo_error",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{setErr: errors.New("setting-fail")},
// 		expect{ok: false, commitTrue: false},
// 	)

// 	run("NG_userRepo_soft_delete_error",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{userErr: errors.New("user-fail")},
// 		expect{ok: false, commitTrue: false},
// 	)

// 	// ---------------- 失敗系：コミット時に done エラー ----------------
// 	run("NG_commit_done_error",
// 		ids{editor: 1, target: 2},
// 		actors{
// 			editor: &domain.User{ID: 1, IsRoot: true},
// 			target: &domain.User{ID: 2, IsRoot: false},
// 		},
// 		fails{doneErr: errors.New("commit-fail")},
// 		expect{ok: false, commitTrue: true}, // commit=true で done を叩いてエラーを表に返す
// 	)
// }
