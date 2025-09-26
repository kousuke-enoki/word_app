package setting

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/userconfig"
	"word_app/backend/src/domain"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntUserConfigRepo struct {
	client serviceinterfaces.EntClientInterface
}

type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int) (*domain.UserConfig, error)
	Upsert(ctx context.Context, cfg *domain.UserConfig) (*domain.UserConfig, error)
	CreateDefault(ctx context.Context, userID int) error
	SoftDeleteByUserID(ctx context.Context, userID int, t time.Time) error
}

func NewEntUserConfigRepo(client serviceinterfaces.EntClientInterface) *EntUserConfigRepo {
	return &EntUserConfigRepo{client: client}
}

func (r *EntUserConfigRepo) GetByUserID(ctx context.Context, uid int) (*domain.UserConfig, error) {
	uc, err := r.client.UserConfig().
		Query().
		Where(userconfig.UserID(uid)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.UserConfig{
		ID:         uc.ID,
		UserID:     uc.UserID,
		IsDarkMode: uc.IsDarkMode,
	}, nil
}

func (r *EntUserConfigRepo) Upsert(ctx context.Context, cfg *domain.UserConfig) (*domain.UserConfig, error) {
	uc, err := r.client.UserConfig().
		Query().
		Where(userconfig.UserIDEQ(cfg.UserID)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if uc == nil {
		uc, err = r.client.UserConfig().
			Create().
			SetUserID(cfg.UserID).
			SetIsDarkMode(cfg.IsDarkMode).
			Save(ctx)
	} else {
		uc, err = r.client.UserConfig().
			UpdateOne(uc).
			SetIsDarkMode(cfg.IsDarkMode).
			Save(ctx)
	}
	if err != nil {
		return nil, err
	}
	return &domain.UserConfig{UserID: uc.UserID, IsDarkMode: uc.IsDarkMode}, nil
}
