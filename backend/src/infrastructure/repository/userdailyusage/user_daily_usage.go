// src/infrastructure/repository/userddailyusage/ent_user_daily_usage_repo.go
package userdailyusage

import (
	"context"
	"time"

	"word_app/backend/src/domain"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
	sqlexec "word_app/backend/src/interfaces/sqlexec"
)

type Repository interface {
	CreateIfNotExists(ctx context.Context, userID int, now time.Time) error
	IncQuizOr429(ctx context.Context, userID int, now time.Time) (*domain.DailyUsageUpdateResult, error)
	IncBulkOr429(ctx context.Context, userID int, now time.Time) (*domain.DailyUsageUpdateResult, error)
}

type EntUserDailyUsageRepo struct {
	client serviceinterfaces.EntClientInterface
	sql    sqlexec.Runner
	jst    *time.Location
}

func NewEntUserDailyUsageRepo(c serviceinterfaces.EntClientInterface, s sqlexec.Runner) *EntUserDailyUsageRepo {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return &EntUserDailyUsageRepo{
		client: c,
		sql:    s,
		jst:    jst,
	}
}

// JSTのその日0:00に丸める
func (r *EntUserDailyUsageRepo) truncateToJST0(t time.Time) time.Time {
	tt := t.In(r.jst)
	y, m, d := tt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, r.jst)
}
