package registeredword

import (
	"context"

	"word_app/backend/ent/registeredword"
)

func (r *EntRegisteredWordReadRepo) FindActiveMapByUserAndWordIDs(ctx context.Context, userID int, wordIDs []int) (map[int]bool, error) {
	if len(wordIDs) == 0 {
		return map[int]bool{}, nil
	}
	rows, err := r.client.RegisteredWord().Query().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDIn(wordIDs...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[int]bool, len(rows))
	for _, rw := range rows {
		m[rw.WordID] = rw.IsActive
	}
	return m, nil
}
