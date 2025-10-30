// backend/src/infrastructure/repository/registeredword/create_active_test.go
package registeredword_test

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"
	entrw "word_app/backend/ent/registeredword"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/infrastructure/repository/registeredword"
	si "word_app/backend/src/interfaces/service_interfaces"
)

type entAdapter struct{ *ent.Client }

func (a entAdapter) EntClient() *ent.Client                    { return a.Client }
func (a entAdapter) ExternalAuth() *ent.ExternalAuthClient     { return a.Client.ExternalAuth }
func (a entAdapter) JapaneseMean() *ent.JapaneseMeanClient     { return a.Client.JapaneseMean }
func (a entAdapter) Quiz() *ent.QuizClient                     { return a.Client.Quiz }
func (a entAdapter) QuizQuestion() *ent.QuizQuestionClient     { return a.Client.QuizQuestion }
func (a entAdapter) RegisteredWord() *ent.RegisteredWordClient { return a.Client.RegisteredWord }
func (a entAdapter) RootConfig() *ent.RootConfigClient         { return a.Client.RootConfig }
func (a entAdapter) Tx(ctx context.Context) (*ent.Tx, error)   { return a.Client.Tx(ctx) }
func (a entAdapter) User() *ent.UserClient                     { return a.Client.User }
func (a entAdapter) UserConfig() *ent.UserConfigClient         { return a.Client.UserConfig }
func (a entAdapter) Word() *ent.WordClient                     { return a.Client.Word }
func (a entAdapter) WordInfo() *ent.WordInfoClient             { return a.Client.WordInfo }
func (a entAdapter) UserDailyUsage() *ent.UserDailyUsageClient { return a.Client.UserDailyUsage }

func newWriteRepo(t *testing.T) (si.EntClientInterface, *registeredword.EntRegisteredWordWriteRepo) {
	t.Helper()

	cli := enttest.Open(t, "sqlite3", "file:test_"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })

	adapter := entAdapter{cli}
	repo := registeredword.NewEntRegisteredWordWriteRepo(adapter)

	return adapter, repo
}

func createUser(t *testing.T, client *ent.Client, email, name string) *ent.User {
	t.Helper()
	u, err := client.User.Create().SetEmail(email).SetName(name).SetPassword("hashed").Save(context.Background())
	require.NoError(t, err)
	return u
}

func createWord(t *testing.T, client *ent.Client, name string) *ent.Word {
	t.Helper()
	w, err := client.Word.Create().SetName(name).Save(context.Background())
	require.NoError(t, err)
	return w
}

func TestEntRegisteredWordWriteRepo_CreateActive(t *testing.T) {
	ctx := context.Background()

	t.Run("success - create active registered word", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "test@example.com", "Test User")
		word := createWord(t, adapter.EntClient(), "apple")

		err := repo.CreateActive(ctx, user.ID, word.ID)
		assert.NoError(t, err)

		// 作成されたか確認
		rw, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user.ID), entrw.WordIDEQ(word.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.True(t, rw.IsActive)
	})

	t.Run("success - create multiple active words", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "multi@example.com", "Multi User")
		w1 := createWord(t, adapter.EntClient(), "apple")
		w2 := createWord(t, adapter.EntClient(), "banana")
		w3 := createWord(t, adapter.EntClient(), "cherry")

		err1 := repo.CreateActive(ctx, user.ID, w1.ID)
		assert.NoError(t, err1)

		err2 := repo.CreateActive(ctx, user.ID, w2.ID)
		assert.NoError(t, err2)

		err3 := repo.CreateActive(ctx, user.ID, w3.ID)
		assert.NoError(t, err3)

		// 全て作成されたか確認
		count, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user.ID), entrw.IsActive(true)).
			Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("success - different users can register same word", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user1 := createUser(t, adapter.EntClient(), "user1@example.com", "User 1")
		user2 := createUser(t, adapter.EntClient(), "user2@example.com", "User 2")
		word := createWord(t, adapter.EntClient(), "shared")

		err1 := repo.CreateActive(ctx, user1.ID, word.ID)
		assert.NoError(t, err1)

		err2 := repo.CreateActive(ctx, user2.ID, word.ID)
		assert.NoError(t, err2)

		// 両方のユーザーが単語を登録できたことを確認
		count, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.WordIDEQ(word.ID)).
			Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("error - duplicate registration (same user and word)", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "duplicate@example.com", "Duplicate User")
		word := createWord(t, adapter.EntClient(), "apple")

		err1 := repo.CreateActive(ctx, user.ID, word.ID)
		assert.NoError(t, err1)

		// 同じ組み合わせで再度作成を試みる
		err2 := repo.CreateActive(ctx, user.ID, word.ID)
		assert.Error(t, err2)
		assert.True(t, ent.IsConstraintError(err2))
	})

	t.Run("error - invalid user ID (foreign key violation)", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		word := createWord(t, adapter.EntClient(), "test")

		err := repo.CreateActive(ctx, 99999, word.ID)
		assert.Error(t, err)
		assert.True(t, ent.IsConstraintError(err))
	})

	t.Run("error - invalid word ID (foreign key violation)", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "fktest@example.com", "FK Test")

		err := repo.CreateActive(ctx, user.ID, 99999)
		assert.Error(t, err)
		assert.True(t, ent.IsConstraintError(err))
	})

	t.Run("error - database closed", func(t *testing.T) {
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := registeredword.NewEntRegisteredWordWriteRepo(closedAdapter)

		user := createUser(t, closedClient, "test@example.com", "Test")
		word := createWord(t, closedClient, "test")

		_ = closedClient.Close()

		err := closedRepo.CreateActive(ctx, user.ID, word.ID)
		assert.Error(t, err)
	})
}
