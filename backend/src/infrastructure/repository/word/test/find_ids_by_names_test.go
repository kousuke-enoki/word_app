// backend/src/infrastructure/repository/word/find_ids_by_names_test.go
package word_test

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/infrastructure/repository/word"
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

func newRepo(t *testing.T) (si.EntClientInterface, *word.EntWordReadRepo) {
	t.Helper()

	cli := enttest.Open(t, "sqlite3", "file:test_"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })

	adapter := entAdapter{cli}
	repo := word.NewEntWordReadRepo(adapter)

	return adapter, repo
}

func createWord(t *testing.T, client *ent.Client, name string) *ent.Word {
	t.Helper()
	w, err := client.Word.Create().SetName(name).Save(context.Background())
	require.NoError(t, err)
	return w
}

func TestEntWordReadRepo_FindIDsByNames(t *testing.T) {
	ctx := context.Background()

	t.Run("success - find multiple words", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 単語を作成
		w1 := createWord(t, adapter.EntClient(), "apple")
		w2 := createWord(t, adapter.EntClient(), "banana")
		w3 := createWord(t, adapter.EntClient(), "cherry")

		// FindIDsByNamesを実行
		result, err := repo.FindIDsByNames(ctx, []string{"apple", "banana", "cherry"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, w1.ID, result["apple"])
		assert.Equal(t, w2.ID, result["banana"])
		assert.Equal(t, w3.ID, result["cherry"])
	})

	t.Run("success - find partial words", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 単語を作成
		w1 := createWord(t, adapter.EntClient(), "apple")
		_ = createWord(t, adapter.EntClient(), "banana")

		// 存在する単語と存在しない単語で検索
		result, err := repo.FindIDsByNames(ctx, []string{"apple", "orange"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result)) // 存在するものだけが返される
		assert.Equal(t, w1.ID, result["apple"])
		_, exists := result["orange"]
		assert.False(t, exists) // 存在しないキーは結果に含まれない
	})

	t.Run("success - find single word", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 単語を作成
		w1 := createWord(t, adapter.EntClient(), "hello")

		// FindIDsByNamesを実行
		result, err := repo.FindIDsByNames(ctx, []string{"hello"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, w1.ID, result["hello"])
	})

	t.Run("success - empty slice input", func(t *testing.T) {
		_, repo := newRepo(t)

		// 空のスライスで検索
		result, err := repo.FindIDsByNames(ctx, []string{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - no matches found", func(t *testing.T) {
		_, repo := newRepo(t)

		// 存在しない単語で検索
		result, err := repo.FindIDsByNames(ctx, []string{"nonexistent"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("success - case sensitive", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 小文字で単語を作成
		_ = createWord(t, adapter.EntClient(), "apple")

		// 大文字で検索
		result, err := repo.FindIDsByNames(ctx, []string{"Apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result)) // ケースは厳密に一致する必要がある
	})

	t.Run("success - duplicate names in input", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 単語を作成
		w1 := createWord(t, adapter.EntClient(), "apple")

		// 重複した名前で検索
		result, err := repo.FindIDsByNames(ctx, []string{"apple", "apple"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result)) // 重複は1つにまとめられる
		assert.Equal(t, w1.ID, result["apple"])
	})

	t.Run("success - large input slice", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 複数の単語を作成
		words := make([]*ent.Word, 0, 50)
		names := make([]string, 0, 50)
		for i := 0; i < 50; i++ {
			name := "word" + string(rune('A'+i%26)) + string(rune('a'+i/26))
			w := createWord(t, adapter.EntClient(), name)
			words = append(words, w)
			names = append(names, name)
		}

		// FindIDsByNamesを実行
		result, err := repo.FindIDsByNames(ctx, names)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 50, len(result))
		// 全てのIDが正しく返されているか確認
		for i, name := range names {
			assert.Equal(t, words[i].ID, result[name], "Mismatch for word: %s", name)
		}
	})

	t.Run("success - mixed existing and non-existing", func(t *testing.T) {
		adapter, repo := newRepo(t)

		// 一部の単語のみ作成
		w1 := createWord(t, adapter.EntClient(), "exist1")
		w2 := createWord(t, adapter.EntClient(), "exist2")

		// 存在するものと存在しないものを混在させて検索
		result, err := repo.FindIDsByNames(ctx, []string{"exist1", "nonexist1", "exist2", "nonexist2"})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result)) // 存在するものだけが返される
		assert.Equal(t, w1.ID, result["exist1"])
		assert.Equal(t, w2.ID, result["exist2"])
	})

	t.Run("error - database closed", func(t *testing.T) {
		// 新しいクライアントを作成して閉じる
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := word.NewEntWordReadRepo(closedAdapter)

		// クライアントを閉じる
		_ = closedClient.Close()

		// 閉じたクライアントでFindIDsByNamesを実行
		result, err := closedRepo.FindIDsByNames(ctx, []string{"test"})

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNewEntWordReadRepo(t *testing.T) {
	t.Run("success - initialize repository", func(t *testing.T) {
		adapter, repo := newRepo(t)

		assert.NotNil(t, repo)
		assert.NotNil(t, adapter)
	})
}
