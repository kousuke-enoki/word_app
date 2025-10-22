// src/infrastructure/repository/registeredword/ent_registered_word_read_repo.go
package registeredword

import (
	"context"

	"word_app/backend/ent/registeredword"
)

func (r *EntRegisteredWordReadRepo) ActiveWordIDSetByUser(ctx context.Context, userID int, wordIDs []int) (map[int]struct{}, error) {
	if len(wordIDs) == 0 {
		return map[int]struct{}{}, nil
	}
	rows, err := r.client.RegisteredWord().Query().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDIn(wordIDs...),
			registeredword.IsActive(true),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	set := make(map[int]struct{}, len(rows))
	for _, r := range rows {
		set[r.WordID] = struct{}{}
	}
	return set, nil
}
