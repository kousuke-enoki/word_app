package seeder

import (
	"context"
	"log"
	"word_app/ent"
	"word_app/ent/partofspeech"
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

// SeedPartOfSpeech 品詞データのシード
func SeedPartOfSpeech(ctx context.Context, client *ent.Client) {
	partsOfSpeech := []string{"noun", "pronoun", "verb", "djective", "adverb",
		"auxiliary_verb", "preposition", "article", "interjection", "conjunction"}

	for _, name := range partsOfSpeech {
		exists, err := client.PartOfSpeech.Query().Where(partofspeech.Name(name)).Exist(ctx)
		if err != nil {
			log.Fatalf("failed to query part of speech: %v", err)
		}
		if !exists {
			_, err := client.PartOfSpeech.Create().
				SetName(name).
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create part of speech: %v", err)
			}
			log.Printf("Part of speech '%s' seeded\n", name)
		}
	}
}

func SeedWords(ctx context.Context, client *ent.Client) {
	// 単語、品詞、日本語の意味を持つデータセット
	words := []struct {
		name           string
		partOfSpeechId int // 品詞を整数で管理
		japaneseMean   string
	}{
		{"able", 4, "できる"},    // 4 = 形容詞
		{"abroad", 5, "海外で"},  // 5 = 副詞
		{"actually", 5, "実際"}, // 5 = 副詞
		{"add", 3, "加える"},     // 3 = 動詞
		{"agree", 3, "同意する"},  // 3 = 動詞
		{"almost", 5, "もう少しで"},
		{"already", 5, "すでに"},
		{"also", 5, "また"},
		{"always", 5, "いつも"},
		{"amount", 1, "量"},
		{"approach", 1, "方法"},  // 1 = 名詞
		{"approach", 3, "近づく"}, // 3 = 動詞
		{"arrive", 3, "到着する"},
		{"attention", 1, "注意"},
		{"average", 4, "平均的な"},
		{"average", 1, "平均"},
		{"become", 3, "～になる"},
		{"begin", 3, "始める"},
		{"believe", 3, "信じる"},
		{"below", 5, "下に"},
		{"bit", 1, "少し"},
		{"bit", 5, "少し"},
		{"borrow", 3, "借りる"},
	}

	for _, w := range words {
		// まず、word テーブルに単語を追加または取得
		existingWord, err := client.Word.Query().Where(word.Name(w.name)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			log.Fatalf("failed to query word: %v", err)
		}

		var createdWord *ent.Word
		if existingWord != nil {
			// 既存の単語がある場合は、それを使う
			createdWord = existingWord
		} else {
			// ない場合は新しい単語を作成
			createdWord, err = client.Word.Create().
				SetName(w.name).
				SetVoiceID("").
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create word: %v", err)
			}
			log.Printf("Word '%s' seeded\n", w.name)
		}

		// word_info テーブルに品詞情報を追加
		wordInfo, err := client.WordInfo.Create().
			SetWordID(createdWord.ID).
			SetPartOfSpeechID(w.partOfSpeechId).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed to create word info: %v", err)
		}

		// japanese_mean テーブルに日本語の意味を追加
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
	SeedPartOfSpeech(ctx, client)
	SeedWords(ctx, client)
}
