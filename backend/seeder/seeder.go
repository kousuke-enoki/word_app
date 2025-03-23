package seeder

import (
	"context"
	"log"

	"word_app/backend/ent"
	"word_app/backend/ent/partofspeech"
	"word_app/backend/ent/user"
	"word_app/backend/ent/word"
	"word_app/backend/src/interfaces"

	"golang.org/x/crypto/bcrypt"
)

// RunSeeder 初回のみシード実行
func RunSeeder(ctx context.Context, client interfaces.ClientInterface) {
	SeedAdminUsers(ctx, client)
	SeedPartOfSpeech(ctx, client)
	SeedWords(ctx, client)
}

// SeedAdminUsers シードデータを流す
func SeedAdminUsers(ctx context.Context, client interfaces.ClientInterface) {
	entClient := client.EntClient()
	exists, err := entClient.User.Query().Where(user.Email("root@example.com")).Exist(ctx)
	if err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	if !exists {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123$"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password")
			return
		}

		_, err = entClient.User.Create().
			SetEmail("root@example.com").
			SetName("Root User").
			SetPassword(string(hashedPassword)).
			SetIsAdmin(true).
			SetIsRoot(true).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed to create root user: %v", err)
		}
		log.Println("Root user seeded")
	}
}

// SeedPartOfSpeech 品詞データのシード
func SeedPartOfSpeech(ctx context.Context, client interfaces.ClientInterface) {
	entClient := client.EntClient()
	partsOfSpeech := []string{"名詞", "代名詞", "動詞", "形容詞", "副詞",
		"助動詞", "前置詞", "冠詞", "間投詞", "接続詞"}

	for _, name := range partsOfSpeech {
		exists, err := entClient.PartOfSpeech.Query().Where(partofspeech.Name(name)).Exist(ctx)
		if err != nil {
			log.Fatalf("failed to query part of speech: %v", err)
		}
		if !exists {
			_, err := entClient.PartOfSpeech.Create().
				SetName(name).
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create part of speech: %v", err)
			}
			log.Printf("Part of speech '%s' seeded\n", name)
		}
	}
}

func SeedWords(ctx context.Context, client interfaces.ClientInterface) {
	entClient := client.EntClient()
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
		{"break", 3, "壊す"},
		{"break", 1, "休憩"},
		{"broad", 4, "幅広い"},
		{"busy", 4, "忙しい"},
		{"cancel", 3, "中止する"},
		{"careful", 4, "丁寧な"},
		{"carefully", 5, "丁寧に"},
		{"certain", 4, "確かな"},
		{"choose", 3, "選ぶ"},
		{"clean", 4, "きれいな"},
		{"clean", 3, "きれいにする"},
		{"clear", 4, "はっきりした"},
		{"clear", 3, "片付ける"},
		{"clerk", 1, "事務員"},
		{"close", 4, "近い"},
		{"close", 3, "閉める"},
		{"collect", 3, "集める"},
		{"common", 4, "共通の"},
		{"company", 1, "会社"},
		{"compare", 3, "比べる"},
		{"condition", 1, "状態"},
		{"connect", 3, "接続する"},
		{"contact", 3, "連絡する"},
		{"contact", 1, "連絡先"},
		{"continue", 3, "続く"},
		{"convenient", 4, "便利な"},
		{"conversation", 1, "会話"},
		{"corporation", 1, "会社"},
		{"customer", 1, "客"},
		{"damage", 1, "損害"},
		{"damage", 3, "損害を与える"},
		{"deal", 3, "扱う"},
		{"deal", 1, "取引"},
		{"decide", 3, "決める"},
		{"decorate", 3, "飾る"},
		{"demonstration", 1, "実演"},
		{"different", 4, "異なる"},
		{"difficult", 4, "難しい"},
		{"discover", 3, "発見する"},
		{"discuss", 3, "話し合う"},
		{"double", 3, "二倍にする"},
		{"double", 4, "二倍の"},
		{"easy", 4, "簡単な"},
		{"education", 1, "教育"},
		{"effect", 1, "効果"},
		{"effort", 1, "努力"},
		{"electricity", 1, "電気"},
		{"empty", 4, "空の"},
		{"empty", 3, "空にする"},
		{"enjoy", 3, "楽しむ"},
		{"enough", 4, "十分な"},
		{"enough", 5, "十分に"},
		{"entrance", 1, "入口"},
		{"especially", 5, "特に"},
		{"example", 1, "例"},
		{"exchange", 3, "交換する"},
		{"exchange", 1, "交換"},
		{"exciting", 4, "人を興奮させるような"},
		{"expensive", 4, "値段が高い"},
		{"experience", 1, "経験"},
		{"experience", 3, "経験する"},
		{"explain", 3, "説明する"},
		{"extra", 4, "余分な"},
		{"extra", 1, "余分な物"},
		{"extra", 5, "余分に"},
		{"fact", 1, "事実"},
		{"fast", 4, "速い"},
		{"fast", 5, "速く"},
		{"favorite", 4, "好きな"},
		{"favorite", 1, "お気に入り"},
		{"feedback", 1, "感想"},
		{"film", 1, "映画"},
		{"film", 3, "撮影する"},
		{"final", 4, "最終的な"},
		{"forecast", 1, "予想"},
		{"forecast", 3, "予想する"},
		{"form", 1, "用紙"},
		{"form", 3, "結成する"},
		{"furniture", 1, "家具"},
		{"further", 4, "さらなる"},
		{"further", 5, "さらに"},
		{"gallery", 1, "美術館"},
		{"glad", 4, "嬉しい"},
		{"government", 1, "政府"},
		{"graduate", 3, "卒業する"},
		{"growth", 1, "成長"},
		{"happen", 3, "起こる"},
		{"healthy", 4, "健康的な"},
		{"afraid", 4, "恐れている"},
		{"important", 4, "重要な"},
		{"instead", 5, "代わりに"},
		{"interested", 4, "関心がある"},
		{"international", 4, "国際的な"},
		{"interview", 1, "面接"},
		{"interview", 3, "面接する"},
		{"item", 1, "物"},
		{"join", 3, "加わる"},
		{"law", 1, "法律"},
		{"lecture", 1, "講義"},
		{"let", 3, "させる"},
		{"library", 1, "図書館"},
		{"loan", 1, "融資"},
		{"loss", 1, "損失"},
		{"lower", 4, "より低い"},
		{"lower", 3, "下げる"},
		{"main", 4, "主な"},
		{"major", 4, "主要な"},
		{"material", 1, "材料"},
		{"meal", 1, "食事"},
		{"meaning", 1, "意味"},
		{"medical", 4, "医療の"},
		{"monthly", 4, "毎月の"},
		{"monthly", 5, "毎月"},
		{"monthly", 1, "月刊誌"},
		{"mostly", 5, "主に"},
		{"national", 4, "全国的な"},
		{"nearby", 4, "近くの"},
		{"nearby", 5, "近くに"},
		{"noise", 1, "騒音"},
		{"notice", 1, "お知らせ"},
		{"notice", 3, "気付く"},
		{"novel", 1, "小説"},
		{"official", 4, "公式の"},
		{"official", 1, "（政府などの）担当者"},
		{"often", 5, "しばしば"},
		{"opinion", 1, "意見"},
		{"own", 4, "自身の"},
		{"own", 3, "所有する"},
		{"park", 3, "駐車する"},
		{"park", 1, "公園"},
		{"part", 1, "部分"},
		{"passenger", 1, "乗客"},
		{"patient", 1, "患者"},
		{"perform", 3, "演じる"},
		{"perhaps", 5, "たぶん"},
		{"personal", 4, "個人の"},
		{"position", 1, "職"},
		{"possible", 4, "可能な"},
		{"press", 1, "報道機関"},
		{"press", 3, "押す"},
		{"pretty", 5, "かなり"},
		{"pretty", 4, "かわいい"},
		{"probably", 5, "おそらく"},
		{"problem", 1, "問題"},
		{"product", 1, "製品"},
		{"progress", 1, "進歩"},
		{"progress", 3, "進む"},
		{"promise", 3, "約束する"},
		{"promise", 1, "約束"},
		{"protect", 3, "保護する"},
		{"proud", 4, "誇りに思う"},
		{"public", 4, "公共"},
		{"public", 1, "一般の人"},
		{"purpose", 1, "目的"},
		{"quality", 1, "質"},
		{"quality", 4, "質が高い"},
		{"quantity", 1, "量"},
		{"quickly", 5, "速く"},
		{"quite", 5, "かなり"},
		{"rather", 5, "むしろ"},
		{"ready", 4, "準備ができて"},
		{"realize", 3, "気付く"},
		{"receive", 3, "受け取る"},
		{"relationship", 1, "関係"},
		{"remain", 3, "～のままでいる"},
		{"remember", 3, "思い出す"},
		{"reply", 3, "返事をする"},
		{"request", 3, "依頼する"},
		{"request", 1, "依頼"},
		{"rest", 1, "残り"},
		{"rest", 3, "休む"},
		{"return", 3, "戻す"},
		{"return", 1, "返品"},
		{"ride", 3, "乗る"},
		{"ride", 1, "乗ること"},
		{"rise", 3, "上がる"},
		{"rise", 1, "上昇"},
		{"role", 1, "役割"},
		{"safe", 4, "安全な"},
		{"safety", 1, "安全"},
		{"save", 3, "節約する"},
		{"science", 1, "科学"},
		{"seem", 3, "～のようだ"},
		{"select", 3, "選ぶ"},
		{"select", 4, "厳選された"},
		{"series", 1, "連続"},
		{"sightseeing", 1, "観光"},
		{"sign", 3, "署名する"},
		{"sign", 1, "看板"},
		{"similar", 4, "似た"},
		{"situation", 1, "状況"},
		{"society", 1, "社会"},
		{"sometime", 5, "いつか"},
		{"soon", 5, "もうすぐ"},
		{"source", 1, "源"},
		{"spend", 3, "費やす"},
		{"stay", 3, "とどまる"},
		{"stay", 1, "滞在"},
		{"still", 5, "まだ"},
		{"succeed", 3, "成功する"},
		{"success", 1, "成功"},
		{"support", 3, "支持する"},
		{"support", 1, "支持"},
		{"sure", 4, "確かな"},
		{"surprising", 4, "驚くべき"},
		{"task", 1, "任務"},
		{"tax", 1, "税"},
		{"technology", 1, "技術"},
		{"tourist", 1, "観光客"},
		{"traffic", 1, "交通量"},
		{"true", 4, "本当の"},
		{"useful", 4, "役に立つ"},
		{"usually", 5, "普段"},
		{"variety", 1, "種類"},
		{"visit", 3, "訪れる"},
		{"visit", 1, "訪問"},
		{"visitor", 1, "訪問者"},
		{"warm", 4, "温かい"},
		{"weather", 1, "天気"},
		{"welcome", 3, "歓迎する"},
		{"welcome", 1, "歓迎"},
		{"whole", 4, "全体の"},
		{"whole", 1, "全体"},
		{"wide", 4, "広い"},
		{"workshop", 1, "研修会"},
		{"worried", 4, "心配して"},
		{"yet", 5, "まだ"},
		{"yet", 10, "けれども"},
	}

	for _, w := range words {
		// まず、word テーブルに単語を追加または取得
		existingWord, err := entClient.Word.Query().Where(word.Name(w.name)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			log.Fatalf("failed to query word: %v", err)
		}

		var createdWord *ent.Word
		if existingWord != nil {
			// 既存の単語がある場合は、それを使う
			createdWord = existingWord
		} else {
			// ない場合は新しい単語を作成
			createdWord, err = entClient.Word.Create().
				SetName(w.name).
				SetVoiceID("").
				Save(ctx)
			if err != nil {
				log.Fatalf("failed to create word: %v", err)
			}
			log.Printf("Word '%s' seeded\n", w.name)
		}

		// word_info テーブルに品詞情報を追加
		wordInfo, err := entClient.WordInfo.Create().
			SetWordID(createdWord.ID).
			SetPartOfSpeechID(w.partOfSpeechId).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed to create word info: %v", err)
		}

		// japanese_mean テーブルに日本語の意味を追加
		_, err = entClient.JapaneseMean.Create().
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
