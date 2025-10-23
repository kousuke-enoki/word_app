package bulk

import (
	"word_app/backend/config"
	"word_app/backend/src/usecase/bulk"

	"github.com/gin-gonic/gin"
)

type BulkHandler struct {
	tokenizeUsecase bulk.TokenizeUsecase
	registerUsecase bulk.RegisterUsecase
	limits          *config.LimitsCfg
}

func NewHandler(
	tokenizeUsecase bulk.TokenizeUsecase,
	registerUsecase bulk.RegisterUsecase,
	limits *config.LimitsCfg,
) *BulkHandler {
	return &BulkHandler{
		tokenizeUsecase: tokenizeUsecase,
		registerUsecase: registerUsecase,
		limits:          limits,
	}
}

// type Service interface {
// 	BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error)
// 	BulkRegister(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error)
// }

type Handler interface {
	TokenizeHandler() gin.HandlerFunc
	RegisterHandler() gin.HandlerFunc
}
