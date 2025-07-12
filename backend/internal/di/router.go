// internal/di/router.go
package di

// func NewRouter(cfg *Config, h *Handlers, jwtMw middleware.JwtMiddleware) *router.RouterImplementation {
// 	r := gin.New()
// 	// ✂ CORS やミドルウェアはここで集中管理
// 	routerConfig.NewRouter(jwtMw, h.Auth, h.User, h.Setting,
// 		h.Word, h.Quiz, h.Result).
// 		SetupRouter(r)
// 	return r
// }
