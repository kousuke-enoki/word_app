package dictimport

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"sync"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"

	"log"

	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
)

// ImportJMdict は巨大な JMdictJSON を並列インポートするエントリポイント
func ImportJMdict(ctx context.Context, path string, cli *ent.Client, opt Options) ([]ImportErr, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Fatalf("failed to close ent client: %v", err)
		}
	}()

	dec := json.NewDecoder(f)
	var root JMdictJSON
	if err := dec.Decode(&root); err != nil {
		return nil, err
	}

	errCh := make(chan ImportErr, 1024)

	var (
		errs   []ImportErr
		mu     sync.Mutex
		recvWG sync.WaitGroup
		wg     sync.WaitGroup
	)

	// ---- エラー集約 goroutine ----
	recvWG.Add(1)
	go func() {
		defer recvWG.Done()
		for ie := range errCh {
			mu.Lock()
			errs = append(errs, ie)
			mu.Unlock()
		}
	}()

	// ---- ワーカープール ----
	jobs := make(chan JMEntry, opt.Workers*2)
	for i := 0; i < opt.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range jobs {
				if err := importOne(ctx, cli, e); err != nil {
					log.Println("id", e.ID)
					errCh <- ImportErr{ID: e.ID, Message: err.Error()}
				}
			}
		}()
	}

	for _, e := range root.Words {
		jobs <- e
	}
	close(jobs)

	wg.Wait()
	close(errCh)
	recvWG.Wait()

	return errs, nil
}

// ---------------- 内部処理 ----------------

// withTx は Tx を張って関数 fn を実行し、commit / rollback を一元管理する。
func withTx(ctx context.Context, cli *ent.Client, fn func(*ent.Tx) error) error {
	tx, err := cli.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func importOne(ctx context.Context, cli *ent.Client, e JMEntry) error {
	return withTx(ctx, cli, func(tx *ent.Tx) error {
		for _, s := range e.Sense {
			if len(s.PartOfSpeech) == 0 {
				continue
			}
			posID := pickPosID(s.PartOfSpeech)

			for _, g := range s.Gloss {
				if g.Lang != "eng" {
					continue
				}
				if err := processGloss(ctx, tx, g.Text, posID, e); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// pickPosID は Priority 付き part‑of‑speech スライスから最適な posID を決定する
func pickPosID(parts []string) int {
	id := mappedPosID(parts[0])
	if id == 11 && len(parts) > 1 {
		for _, c := range parts[1:] {
			if m := mappedPosID(c); m != 11 {
				return m
			}
		}
	}
	return id
}

// processGloss は 1 つの英語 Gloss に対して Word / WordInfo / JapaneseMean を確立する
func processGloss(ctx context.Context, tx *ent.Tx, text string, posID int, entry JMEntry) error {
	w, err := ensureWord(ctx, tx, text)
	if err != nil {
		return err
	}

	wi, err := ensureWordInfo(ctx, tx, w.ID, posID)
	if err != nil {
		return err
	}

	return ensureJapaneseMeans(ctx, tx, wi.ID, entry)
}

// ensureWord は text に対応する ent.Word を取得 / 作成する
func ensureWord(ctx context.Context, tx *ent.Tx, text string) (*ent.Word, error) {
	w, err := tx.Word.Query().Where(word.NameEQ(text)).Only(ctx)
	if ent.IsNotFound(err) {
		reSpace := regexp.MustCompile(`\s`)
		reSymbol := regexp.MustCompile(`[!?\(\)0-9]`)
		return tx.Word.Create().
			SetName(text).
			SetIsIdioms(reSpace.MatchString(text)).
			SetIsSpecialCharacters(reSymbol.MatchString(text)).
			SetRegistrationCount(0).
			Save(ctx)
	}
	return w, err
}

// ensureWordInfo は (wordID,posID) に対応する WordInfo を取得 / 作成する
func ensureWordInfo(ctx context.Context, tx *ent.Tx, wordID, posID int) (*ent.WordInfo, error) {
	wi, err := tx.WordInfo.Query().
		Where(wordinfo.WordID(wordID), wordinfo.PartOfSpeechIDEQ(posID)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return tx.WordInfo.Create().
			SetWordID(wordID).
			SetPartOfSpeechID(posID).
			Save(ctx)
	}
	return wi, err
}

// ensureJapaneseMeans は Entry 情報から日本語訳を upsert する
func ensureJapaneseMeans(ctx context.Context, tx *ent.Tx, wiID int, entry JMEntry) error {
	// kanji.text 優先、無い場合 kana[0]
	if len(entry.Kanji) > 0 {
		for _, kj := range entry.Kanji {
			if kj.Text == "" {
				continue
			}
			if err := upsertJapaneseMean(ctx, tx, wiID, kj.Text); err != nil {
				return err
			}
		}
	} else if len(entry.Kana) > 0 && entry.Kana[0].Text != "" {
		if err := upsertJapaneseMean(ctx, tx, wiID, entry.Kana[0].Text); err != nil {
			return err
		}
	}
	return nil
}

func upsertJapaneseMean(ctx context.Context, tx *ent.Tx, wiID int, name string) error {
	return tx.JapaneseMean.Create().
		SetWordInfoID(wiID).
		SetName(name).
		OnConflict(sql.ConflictColumns("word_info_id", "name")).
		DoNothing().
		Exec(ctx)
}

var posCode2ID = map[string]int{
	"n":  1,
	"pn": 2,
	"vs": 3, "v5r": 3, "v1": 3, "vi": 3, "vt": 3,
	"adj-i": 4, "adj-na": 4, "adj-f": 4,
	"adv": 5, "adv-to": 5,
	"aux-v": 6,
	"prep":  7,
	"art":   8,
	"int":   9,
	"conj":  10,
	"exp":   11,
	"unc":   12,
}

// mappedPosID は 品詞コード → 内部 ID 変換（未知コードは 11）
func mappedPosID(code string) int {
	if id, ok := posCode2ID[code]; ok {
		return id
	}
	return 11
}
