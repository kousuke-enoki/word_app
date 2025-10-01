// infra/mapper/user/user_mapper.go
package user

import (
	"word_app/backend/ent"
	dm "word_app/backend/src/domain"
)

type opts struct {
	auths           []*ent.ExternalAuth
	computedHasLine *bool
}

// Functional Options
// Line連携などで絞る場合
type Option func(*opts)

func WithAuths(auths []*ent.ExternalAuth) Option {
	return func(o *opts) { o.auths = auths }
}

// eager-loadしていないが、呼び出し側で「LINE連携あり」を知っている場合に使う
func WithComputedFlags(hasLine bool) Option {
	return func(o *opts) { o.computedHasLine = &hasLine }
}

// null 安全な *string 化
func strPtrOrNil(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

// 1件用
func MapEntUser(u *ent.User, opt ...Option) *dm.User {
	if u == nil {
		return nil
	}
	o := &opts{}
	for _, f := range opt {
		if f == nil {
			continue
		}
		f(o)
	}

	hasPwd := u.Password != nil && *u.Password != ""
	hasLine := false

	switch {
	case o.computedHasLine != nil:
		hasLine = *o.computedHasLine
	case len(o.auths) > 0:
		hasLine = true
	}

	return &dm.User{
		ID:          u.ID,
		Email:       strPtrOrNil(u.Email),
		Name:        u.Name,
		Password:    zeroIfNil(u.Password), // 認証用途で必要なら。不要なら空文字で返す or マスク
		IsAdmin:     u.IsAdmin,
		IsRoot:      u.IsRoot,
		IsTest:      u.IsTest,
		HasPassword: hasPwd,
		HasLine:     hasLine,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func zeroIfNil(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// 複数件用（一覧など）
func MapEntUsers(list []*ent.User, opt ...Option) []*dm.User {
	out := make([]*dm.User, 0, len(list))
	for _, u := range list {
		out = append(out, MapEntUser(u, opt...))
	}
	return out
}
