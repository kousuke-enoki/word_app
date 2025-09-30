// infra/entrepo/user_list_repo.go
package user

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	usermapper "word_app/backend/src/infrastructure/mapper/user"

	"entgo.io/ent/dialect/sql"
)

type UserQueryRepository struct {
	Client *ent.Client
}

func (r *EntUserRepo) ListUsers(ctx context.Context, f repository.UserListFilter) (*repository.UserListResult, error) {
	base := r.client.User().
		Query().
		Where(user.DeletedAtIsNil())

	// 検索
	if f.Search != "" {
		base = base.Where(
			user.Or(
				user.NameContains(f.Search),
				user.EmailContains(f.Search),
			),
		)
	}

	// 総数
	totalCount, err := base.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	// 一覧取得（LINE連携を判定したいので provider="line" で eager load）
	q := base.Clone().
		Offset(f.Offset).
		Limit(f.Limit).
		WithExternalAuths(func(q *ent.ExternalAuthQuery) {
			q.Where(externalauth.ProviderEqualFold("line"), externalauth.DeletedAtIsNil())
		})

	// ソート
	switch f.SortBy {
	case "name":
		if f.Order == "asc" {
			q = q.Order(ent.Asc(user.FieldName))
		} else {
			q = q.Order(ent.Desc(user.FieldName))
		}
	case "email":
		if f.Order == "asc" {
			q = q.Order(ent.Asc(user.FieldEmail))
		} else {
			q = q.Order(ent.Desc(user.FieldEmail))
		}
	case "role":
		if f.Order == "asc" {
			q = q.Order(func(s *sql.Selector) {
				s.OrderBy(
					sql.Desc(s.C(user.FieldIsRoot)),
					sql.Desc(s.C(user.FieldIsAdmin)),
					sql.Asc(s.C(user.FieldIsTest)),
				)
			})
		} else {
			q = q.Order(func(s *sql.Selector) {
				s.OrderBy(
					sql.Asc(s.C(user.FieldIsRoot)),
					sql.Asc(s.C(user.FieldIsAdmin)),
					sql.Desc(s.C(user.FieldIsTest)),
				)
			})
		}
	default:
		// デフォルトは created_at desc などにしてもOK（任意）
	}

	entUsers, err := q.All(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	// Ent -> Domain（HasLineは eager-loaded auths から算出）
	out := make([]*domain.User, 0, len(entUsers))
	for _, u := range entUsers {
		out = append(out, usermapper.MapEntUser(u, usermapper.WithAuths(u.Edges.ExternalAuths)))
	}

	return &repository.UserListResult{
		Users:      out,
		TotalCount: totalCount,
	}, nil
}
