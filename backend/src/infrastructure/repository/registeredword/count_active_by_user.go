package registeredword

import (
	"context"

	"word_app/backend/ent/registeredword"
)

func (r *EntRegisteredWordReadRepo) CountActiveByUser(ctx context.Context, userID int) (int, error) {
	return r.client.RegisteredWord().Query().
		Where(registeredword.UserIDEQ(userID), registeredword.IsActive(true)).
		Count(ctx)
}
