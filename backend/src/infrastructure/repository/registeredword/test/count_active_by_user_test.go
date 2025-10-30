// backend/src/infrastructure/repository/registeredword/count_active_by_user_test.go
package registeredword_test

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/infrastructure/repository/registeredword"
	si "word_app/backend/src/interfaces/service_interfaces"
)

func newReadRepo(t *testing.T) (si.EntClientInterface, *registeredword.EntRegisteredWordReadRepo) {
	t.Helper()

	cli := enttest.Open(t, "sqlite3", "file:test_"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })

	adapter := entAdapter{cli}
	repo := registeredword.NewEntRegisteredWordReadRepo(adapter)

	return adapter, repo
}

func createRegisteredWord(t *testing.T, client *ent.Client, userID, wordID int, isActive bool) {
	t.Helper()
	_, err := client.RegisteredWord.Create().
		SetUserID(userID).
		SetWordID(wordID).
		SetIsActive(isActive).
		Save(context.Background())
	require.NoError(t, err)
}

func TestEntRegisteredWordReadRepo_CountActiveByUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success - count zero active words", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "zero@example.com", "Zero User")

		count, err := repo.CountActiveByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("success - count single active word", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "single@example.com", "Single User")
		word := createWord(t, adapter.EntClient(), "apple")
		createRegisteredWord(t, adapter.EntClient(), user.ID, word.ID, true)

		count, err := repo.CountActiveByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("success - count multiple active words", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user := createUser(t, adapter.EntClient(), "multi@example.com", "Multi User")
		w1 := createWord(t, adapter.EntClient(), "apple")
		w2 := createWord(t, adapter.EntClient(), "banana")
		w3 := createWord(t, adapter.EntClient(), "cherry")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, true)

		count, err := repo.CountActiveByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("success - count ignores inactive words", func(t *testing.T) {
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

		count, err := repo.CountActiveByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("success - count only for specific user", func(t *testing.T) {
		adapter, repo := newReadRepo(t)

		user1 := createUser(t, adapter.EntClient(), "user1@example.com", "User 1")
		user2 := createUser(t, adapter.EntClient(), "user2@example.com", "User 2")
		w1 := createWord(t, adapter.EntClient(), "apple")
		w2 := createWord(t, adapter.EntClient(), "banana")

		createRegisteredWord(t, adapter.EntClient(), user1.ID, w1.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user1.ID, w2.ID, true)
		createRegisteredWord(t, adapter.EntClient(), user2.ID, w1.ID, true)

		count1, err := repo.CountActiveByUser(ctx, user1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, count1)

		count2, err := repo.CountActiveByUser(ctx, user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count2)
	})

	t.Run("success - count for non-existent user", func(t *testing.T) {
		_, repo := newReadRepo(t)

		count, err := repo.CountActiveByUser(ctx, 99999)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("error - database closed", func(t *testing.T) {
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := registeredword.NewEntRegisteredWordReadRepo(closedAdapter)

		_ = closedClient.Close()

		count, err := closedRepo.CountActiveByUser(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})
}
