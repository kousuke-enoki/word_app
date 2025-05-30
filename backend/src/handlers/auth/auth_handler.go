package auth

// type AuthHandler struct {
// 	AuthClient   interfaces.AuthClient
// 	jwtGenerator interfaces.JWTGenerator
// }

// func NewAuthHandler(client interfaces.AuthClient, jwtGen interfaces.JWTGenerator) *AuthHandler {
// 	return &AuthHandler{
// 		AuthClient:   client,
// 		jwtGenerator: jwtGen,
// 	}
// }

// func (h *AuthHandler) LineLogin(c *gin.Context) {
// 	state, nonce := newState(), newNonce()
// 	url := h.uc.StartLogin(c, state, nonce)
// 	c.Redirect(http.StatusFound, url)
// }

// func (h *AuthHandler) LineCallback(c *gin.Context) {
// 	code := c.Query("code")
// 	res, err := h.uc.HandleCallback(c, code, c.Query("state"), loadNonce(c))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err) // temp_token など
// 	}

// 	if res.NeedPassword {
// 		c.JSON(200, res) // temp_token など
// 	} else {
// 		c.JSON(200, gin.H{"token": res.Token})
// 	}
// }

// func (h *AuthHandler) LineComplete(c *gin.Context) {
// 	var req struct{ TempToken, Password string }
// 	_ = c.ShouldBindJSON(&req)
// 	jwt, err := h.uc.CompleteSignUp(c, req.TempToken, req.Password)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err) // temp_token など
// 	}
// 	c.JSON(200, gin.H{"token": jwt})
// }
