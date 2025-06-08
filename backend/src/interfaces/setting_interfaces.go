// interfaces/handlers.go
package interfaces

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

type SettingHandler interface {
	GetUserSettingHandler() gin.HandlerFunc
	SaveUserSettingHandler() gin.HandlerFunc
	GetRootSettingHandler() gin.HandlerFunc
	SaveRootSettingHandler() gin.HandlerFunc
	GetAuthSettingHandler() gin.HandlerFunc
}

// 他の箇所はmodelsを定義してレスポンスを返すが、ここではentの型をそのまま使用してレスポンスを生成する。
// entの型にある、omitemptyを外すことでレスポンスは返せる。（omitemptyはnull,falseなどの場合なにも返さないようにするもの）
type SettingClient interface {
	GetUserConfig(ctx context.Context, userID int) (*ent.UserConfig, error)
	UpdateUserConfig(ctx context.Context, userID int, isDarkMode bool) (*ent.UserConfig, error)
	GetRootConfig(ctx context.Context, userID int) (*ent.RootConfig, error)
	UpdateRootConfig(
		ctx context.Context,
		userID int,
		editingPermission string,
		isTestUserMode bool,
		isEmailAuthCheck bool,
		isLineAuth bool,
	) (*ent.RootConfig, error)
	GetAuthConfig(ctx context.Context) (*models.AuthSettingResponse, error)
}

type SettingValidator interface {
	ValidateRootConfig(SignUpRequest *models.RootConfig) []*models.FieldError
}
