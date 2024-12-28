package interfaces

// JWTGenerator インターフェースを定義
type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}
