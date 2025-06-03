package jwt

// JwtAuth : JWT検証 & ユーザー情報(ロール)取得
// func (m *JwtMiddleware) JwtCheckMiddleware() gin.HandlerFunc {
// 	logrus.Info("JwtCheckMiddleware")
// 	jwtSecret := os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable is required")
// 	}

// 	return func(c *gin.Context) {
// 		// Authorization ヘッダーからトークン取得
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		// JWTトークンを解析
// 		token, err := jwt.ParseWithClaims(tokenString, &auth_jwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method")
// 			}
// 			return []byte(jwtSecret), nil
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
// 			c.Abort()
// 			return
// 		}

// 		// トークンが有効かつ claims に userID があるかを確認
// 		if claims, ok := token.Claims.(*auth_jwt.Claims); ok && token.Valid {
// 			userID := claims.UserID
// 			if userID == "" {
// 				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: userID not found"})
// 				c.Abort()
// 				return
// 			}

// 			userIDInt, err := strconv.Atoi(userID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
// 				c.Abort()
// 				return
// 			}

// 			// ここで Ent を用いてユーザー情報を取得し、admin/root を確認する
// 			entClient := database.GetEntClient()
// 			u, err := entClient.User.
// 				Query().
// 				Where(entUser.ID(userIDInt)).
// 				Select(user.FieldIsAdmin, user.FieldIsRoot).
// 				Only(c)
// 			if err != nil {
// 				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or DB error"})
// 				c.Abort()
// 				return
// 			}

// 			isAdmin := u.IsAdmin
// 			isRoot := u.IsRoot

// 			// gin.Context に格納して後続ハンドラーで利用できるようにする
// 			c.Set("userID", userIDInt)
// 			c.Set("isAdmin", isAdmin)
// 			c.Set("isRoot", isRoot)
// 			c.Next()

// 		} else {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}
// 	}
// }
