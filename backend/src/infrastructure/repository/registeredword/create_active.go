package registeredword

import "context"

func (r *EntRegisteredWordWriteRepo) CreateActive(ctx context.Context, userID, wordID int) error {
	_, err := r.client.RegisteredWord().
		Create().
		SetUserID(userID).
		SetWordID(wordID).
		SetIsActive(true).
		Save(ctx)
	return err
}
