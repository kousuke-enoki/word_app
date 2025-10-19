// auth_middleware_test.go
package jwt

// func TestAuthMiddleware(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	newRouter := func(mw gin.HandlerFunc) *gin.Engine {
// 		r := gin.New()
// 		r.Use(mw)
// 		r.GET("/ping", func(c *gin.Context) {
// 			// ミドルウェアがセットした値をそのまま echo
// 			uid, _ := c.Get("userID")
// 			admin, _ := c.Get("isAdmin")
// 			root, _ := c.Get("isRoot")
// 			c.JSON(200, gin.H{"uid": uid, "admin": admin, "root": root})
// 		})
// 		return r
// 	}

// 	t.Run("header missing → 401", func(t *testing.T) {
// 		mockVal := mockJwt.NewMockTokenValidator(t)
// 		mw := jwt.NewMiddleware(mockVal).AuthMiddleware()
// 		// mw := (&jwt.Middleware{tokenValidator: mockVal}).AuthMiddleware()
// 		r := newRouter(mw)

// 		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		require.Equal(t, http.StatusUnauthorized, w.Code)
// 		require.Contains(t, w.Body.String(), "authorization header required")
// 	})

// 	t.Run("invalid token → 401", func(t *testing.T) {
// 		mockVal := mockJwt.NewMockTokenValidator(t)
// 		mockVal.
// 			EXPECT().
// 			Validate(mock.Anything, "badtoken").
// 			Return(contextutil.UserRoles{}, errors.New("token_invalid parse_error")).
// 			Once()
// 		mw := jwt.NewMiddleware(mockVal).AuthMiddleware()
// 		// mw := (&jwt.Middleware{tokenValidator: mockVal}).AuthMiddleware()
// 		r := newRouter(mw)

// 		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
// 		req.Header.Set("Authorization", "Bearer badtoken")
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		require.Equal(t, http.StatusUnauthorized, w.Code)
// 		require.Contains(t, w.Body.String(), "token_invalid parse_error")
// 	})

// 	t.Run("valid token → 200 & context set", func(t *testing.T) {
// 		mockVal := mockJwt.NewMockTokenValidator(t)

// 		mockVal.
// 			EXPECT().
// 			Validate(mock.Anything, "good"). // ← 修正
// 			Return(contextutil.UserRoles{UserID: 7, IsAdmin: true, IsRoot: false}, nil).
// 			Once()

// 		mw := jwt.NewMiddleware(mockVal).AuthMiddleware()
// 		r := newRouter(mw)

// 		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
// 		req.Header.Set("Authorization", "Bearer good") // ← これと一致
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		require.Equal(t, http.StatusOK, w.Code)
// 		require.JSONEq(t,
// 			`{"uid":7,"admin":true,"root":false}`,
// 			string(bytes.TrimSpace(w.Body.Bytes())),
// 		)
// 	})
// }
