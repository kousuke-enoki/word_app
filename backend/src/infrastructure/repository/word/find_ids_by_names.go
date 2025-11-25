// src/infrastructure/repository/word/ent_word_read_repo.go
package word

import (
	"context"

	"word_app/backend/ent/word"
)

func (r *EntWordReadRepo) FindIDsByNames(ctx context.Context, names []string) (map[string]int, error) {
	found, err := r.client.Word().Query().Where(word.NameIn(names...)).All(ctx)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int, len(found))
	for _, w := range found {
		res[w.Name] = w.ID
	}
	return res, nil
}
