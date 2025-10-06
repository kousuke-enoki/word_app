package user

import (
	"context"

	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/repoerr"
)

func (r *EntUserRepo) Create(ctx context.Context, u *domain.User) (user *domain.User, err error) {
	var emailPtr *string
	if u.Email != nil { // Ent も Nillable にした前提
		Email := *u.Email // string 取り出し
		emailPtr = &Email // ポインタ化（そのまま u.Email でも良い）
	}

	entUser, err := r.client.User().
		Create().
		SetNillableEmail(emailPtr).
		SetName(u.Name).
		SetPassword(u.Password).
		SetIsAdmin(u.IsAdmin).
		SetIsRoot(u.IsRoot).
		SetIsTest(u.IsTest).
		Save(ctx)
	if err != nil {
		return nil, repoerr.FromEnt(err, "failed to user create", "")
	}
	u.ID = entUser.ID
	u.CreatedAt = entUser.CreatedAt
	u.UpdatedAt = entUser.UpdatedAt

	return u, nil
}
