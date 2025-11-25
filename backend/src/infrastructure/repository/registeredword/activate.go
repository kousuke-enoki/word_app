package registeredword

import (
	"context"

	"word_app/backend/ent/registeredword"
)

func (r *EntRegisteredWordWriteRepo) Activate(ctx context.Context, userID, wordID int) error {
	_, err := r.client.RegisteredWord().
		Update().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDEQ(wordID),
		).
		SetIsActive(true).
		Save(ctx)
	return err
}
