// backend/src/infrastructure/repository/registeredword/active_word_id_set_by_user_test.go
package registeredword_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"word_app/backend/src/infrastructure/repository/registeredword"
)

func TestEntRegisteredWordReadRepo_ActiveWordIDSetByUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success - find active word IDs only", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "active@example.com", "Active User")
		w1 := createWord(t, adapter.EntClient(), "active1")
		w2 := createWord(t, adapter.EntClient(), "inactive1")
		w3 := createWord(t, adapter.EntClient(), "active2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, true)

		result, err := repo.ActiveWordIDSetByUser(ctx, user.ID, []int{w1.ID, w2.ID, w3.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result)) // アクティブなものだけ
		_, exists1 := result[w1.ID]
		assert.True(t, exists1)
		_, exists2 := result[w2.ID]
		assert.False(t, exists2)
		_, exists3 := result[w3.ID]
		assert.True(t, exists3)
	})

	t.Run("success - find partial active words", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "partial@example.com", "Partial User")
		w1 := createWord(t, adapter.EntClient(), "found1")
		w2 := createWord(t, adapter.EntClient(), "notfound")
		w3 := createWord(t, adapter.EntClient(), "found2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, true)

		result, err := repo.ActiveWordIDSetByUser(ctx, user.ID, []int{w1.ID, w2.ID, w3.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result)) // 存在しないword2は含まれない
		_, exists1 := result[w1.ID]
		assert.True(t, exists1)
		_, exists3 := result[w3.ID]
		assert.True(t, exists3)
	})

	t.Run("success - empty input slice", func(t *testing.T) {
		_, repo := newReadRepo(t)

		result, err := repo.ActiveWordIDSetByUser(ctx, 1, []int{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - no active words", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "noactive@example.com", "No Active User")
		w1 := createWord(t, adapter.EntClient(), "inactive1")
		w2 := createWord(t, adapter.EntClient(), "inactive2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, false)

		result, err := repo.ActiveWordIDSetByUser(ctx, user.ID, []int{w1.ID, w2.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - single active word", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "single@example.com", "Single User")
		w1 := createWord(t, adapter.EntClient(), "apple")
		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)

		result, err := repo.ActiveWordIDSetByUser(ctx, user.ID, []int{w1.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		_, exists := result[w1.ID]
		assert.True(t, exists)
	})

	t.Run("success - mixed active and inactive", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "mixed@example.com", "Mixed User")
		w1 := createWord(t, adapter.EntClient(), "active1")
		w2 := createWord(t, adapter.EntClient(), "inactive1")
		w3 := createWord(t, adapter.EntClient(), "active2")
		w4 := createWord(t, adapter.EntClient(), "inactive2")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w4.ID, false)

		result, err := repo.ActiveWordIDSetByUser(ctx, user.ID, []int{w1.ID, w2.ID, w3.ID, w4.ID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result))
		_, exists1 := result[w1.ID]
		assert.True(t, exists1)
		_, exists2 := result[w2.ID]
		assert.False(t, exists2)
		_, exists3 := result[w3.ID]
		assert.True(t, exists3)
		_, exists4 := result[w4.ID]
		assert.False(t, exists4)
	})

	t.Run("success - only for specific user", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user1 := createUser(t, adapter.EntClient(), "user1@example.com", "User 1")
		user2 := createUser(t, adapter.EntClient(), "user2@example.com", "User 2")
		w1 := createWord(t, adapter.EntClient(), "shared")

		createRegisteredWord(t, adapter.EntClient(), user1.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user2.ID, w1.ID, false)

		result1, err := repo.ActiveWordIDSetByUser(ctx, user1.ID, []int{w1.ID})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result1))
		_, exists := result1[w1.ID]
		assert.True(t, exists)

		result2, err := repo.ActiveWordIDSetByUser(ctx, user2.ID, []int{w1.ID})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(result2))
	})

	t.Run("error - database closed", func(t *testing.T) {
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := registeredword.NewEntRegisteredWordReadRepo(closedAdapter)

		_ = closedClient.Close()

		result, err := closedRepo.ActiveWordIDSetByUser(ctx, 1, []int{1, 2})
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
