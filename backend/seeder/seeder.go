package seeder

import (
	"context"
	"log"
	"word_app/ent"
	"word_app/ent/user"
	"word_app/ent/word"
)

// SeedAdminUsers シードデータを流す
func SeedAdminUsers(ctx context.Context, client *ent.Client) {
	exists, err := client.User.Query().Where(user.Email("admin@example.com")).Exist(ctx)
	if err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	if !exists {
		_, err := client.User.Create().
			SetEmail("admin@example.com").
			SetName("Admin User").
			SetPassword("hashed_password").
			SetAdmin(true).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed to create admin user: %v", err)
		}
		log.Println("Admin user seeded")
	}
}

func SeedWords(ctx context.Context, client *ent.Client) {
	// 単語、品詞、日本語の意味を持つデータセット
	words := []struct {
		name         string
		partOfSpeech int // 品詞を整数で管理 (例: 0 = 形容詞, 1 = 動詞, 2 = 副詞など)
		japaneseMean string
	}{
		{"able", 0, "できる"},    // 0 = 形容詞
		{"abroad", 2, "海外で"},  // 2 = 副詞
		{"actually", 2, "実際"}, // 2 = 副詞
		{"add", 1, "加える"},     // 1 = 動詞
		{"agree", 1, "同意する"},  // 1 = 動詞
	}

	for _, w := range words {
		// まず、word テーブルに単語を追加する
		exists, err := client.Word.Query().Where(word.Name(w.name)).Exist(ctx)
		if err != nil {
			log.Fatalf("failed to query word: %v", err)
		}
		if !exists {
			createdWord, err := client.Word.Create().
				SetName(w.name).
				SetVoiceID("").
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create word: %v", err)
			}
			log.Printf("Word '%s' seeded\n", w.name)

			// 次に、word_info テーブルに品詞情報を追加
			wordInfo, err := client.WordInfo.Create().
				SetWordID(createdWord.ID).
				SetPartOfSpeech(w.partOfSpeech).
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create word info: %v", err)
			}

			// 最後に、japanese_mean テーブルに日本語の意味を追加
			_, err = client.JapaneseMean.Create().
				SetWordInfoID(wordInfo.ID).
				SetName(w.japaneseMean).
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create japanese mean: %v", err)
			}
			log.Printf("Japanese meaning for '%s' seeded\n", w.name)
		}
	}
}

// 1	able
// 形 できる
// 2
// ■■ abroad
// 副 海外で
// 3
// ■■ actually
// 副 実際
// 4
// ■■ add
// 動 加える
// 5
// ■■ agree
// 動 同意する

// RunSeeder 初回のみシード実行
func RunSeeder(ctx context.Context, client *ent.Client) {
	SeedAdminUsers(ctx, client)
	SeedWords(ctx, client)
}
