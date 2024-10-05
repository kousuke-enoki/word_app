package seeder

import (
	"context"
	"log"
	"word_app/ent"
	"word_app/ent/partofspeech"
	"word_app/ent/user"
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
	partsOfSpeech := []string{"Noun", "Verb", "Adjective", "Adverb"}

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

// RunSeeder 初回のみシード実行
func RunSeeder(ctx context.Context, client *ent.Client) {
	SeedAdminUsers(ctx, client)
	SeedPartOfSpeech(ctx, client)
}
