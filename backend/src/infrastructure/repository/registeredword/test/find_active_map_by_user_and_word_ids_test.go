// backend/src/infrastructure/repository/registeredword/find_active_map_by_user_and_word_ids_test.go
package registeredword_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"word_app/backend/src/infrastructure/repository/registeredword"
)

func TestEntRegisteredWordReadRepo_FindActiveMapByUserAndWordIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("success - find all active and inactive", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "find@example.com", "Find User")
		w1 := createWord(t, adapter.EntClient(), "active1")
		w2 := createWord(t, adapter.EntClient(), "inactive1")
		w3 := createWord(t, adapter.EntClient(), "active2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, true)

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, user.ID, []int{w1.ID, w2.ID, w3.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, len(result))
		assert.True(t, result[w1.ID])
		assert.False(t, result[w2.ID])
		assert.True(t, result[w3.ID])
	})

	t.Run("success - find partial words", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "partial@example.com", "Partial User")
		w1 := createWord(t, adapter.EntClient(), "found1")
		w2 := createWord(t, adapter.EntClient(), "notfound")
		w3 := createWord(t, adapter.EntClient(), "found2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, false)

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, user.ID, []int{w1.ID, w2.ID, w3.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result)) // 存在しないword2は含まれない
		assert.True(t, result[w1.ID])
		assert.False(t, result[w3.ID])
	})

	t.Run("success - empty input slice", func(t *testing.T) {
		_, repo := newReadRepo(t)

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, 1, []int{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - no matches", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "nomatch@example.com", "No Match User")
		w1 := createWord(t, adapter.EntClient(), "notregistered")

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, user.ID, []int{w1.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - single word active", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "single@example.com", "Single User")
		w1 := createWord(t, adapter.EntClient(), "apple")
		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, user.ID, []int{w1.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.True(t, result[w1.ID])
	})

	t.Run("success - single word inactive", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "inactive@example.com", "Inactive User")
		w1 := createWord(t, adapter.EntClient(), "banana")
		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, false)

		result, err := repo.FindActiveMapByUserAndWordIDs(ctx, user.ID, []int{w1.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.False(t, result[w1.ID])
	})

	t.Run("success - only for specific user", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user1 := createUser(t, adapter.EntClient(), "user1@example.com", "User 1")
		user2 := createUser(t, adapter.EntClient(), "user2@example.com", "User 2")
		w1 := createWord(t, adapter.EntClient(), "shared")

		createRegisteredWord(t, adapter.EntClient(), user1.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user2.ID, w1.ID, false)

		result1, err := repo.FindActiveMapByUserAndWordIDs(ctx, user1.ID, []int{w1.ID})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result1))
		assert.True(t, result1[w1.ID])

		result2, err := repo.FindActiveMapByUserAndWordIDs(ctx, user2.ID, []int{w1.ID})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result2))
		assert.False(t, result2[w1.ID])
	})

	t.Run("error - database closed", func(t *testing.T) {
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := registeredword.NewEntRegisteredWordReadRepo(closedAdapter)

		_ = closedClient.Close()

		result, err := closedRepo.FindActiveMapByUserAndWordIDs(ctx, 1, []int{1, 2})
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
