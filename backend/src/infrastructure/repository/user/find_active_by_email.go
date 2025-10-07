package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	usermapper "word_app/backend/src/infrastructure/mapper/user"
	"word_app/backend/src/infrastructure/repoerr"
)

func (r *EntUserRepo) FindActiveByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(user.EmailEQ(email), user.DeletedAtIsNil()).
		Select(user.FieldID, user.FieldName, user.FieldEmail, user.FieldPassword, user.FieldIsAdmin, user.FieldIsRoot, user.FieldIsTest, user.FieldCreatedAt, user.FieldUpdatedAt).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repoerr.FromEnt(err, "user not found", "duplicate email")
		}
		return nil, repoerr.FromEnt(err, "internal", "internal server error")
	}
	// auths未ロード → 事前計算フラグなしのまま
	return usermapper.MapEntUser(u), nil
}
