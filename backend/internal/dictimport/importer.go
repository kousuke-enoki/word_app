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
)

func ImportJMdict(ctx context.Context, path string, cli *ent.Client, opt Options) ([]ImportErr, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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
	// ------------ 受信 goroutine ------------
	recvWG.Add(1)
	go func() {
		defer recvWG.Done()
		for ie := range errCh {
			mu.Lock()
			errs = append(errs, ie)
			mu.Unlock()
		}
	}()

	jobs := make(chan JMEntry, opt.Workers*2)
	for i := 0; i < opt.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range jobs {
				if err := importOne(ctx, cli, e); err != nil {
					log.Println("id", e.ID)
					// チャネルに渡す
					errCh <- ImportErr{
						ID:      e.ID,
						Message: err.Error(),
					}
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

func importOne(ctx context.Context, cli *ent.Client, e JMEntry) (err error) {
	tx, err := cli.Tx(ctx)
	if err != nil {
		return err
	}
	var txErr error // ← DB処理の成否は txErr に集約
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if txErr != nil {
			_ = tx.Rollback()
			err = txErr
		} else {
			err = tx.Commit() // commit の結果を export
		}
	}()

	// ❶ sense ループ
	for _, s := range e.Sense {
		if len(s.PartOfSpeech) == 0 {
			continue
		}
		posID := mappedPosID(s.PartOfSpeech[0])
		if posID == 11 && len(s.PartOfSpeech) > 1 {
			for _, c := range s.PartOfSpeech[1:] {
				if mapped := mappedPosID(c); mapped != 11 {
					posID = mapped
					break
				}
			}
		}

		// ❷ gloss 英語ごとに Word を生成
		for _, g := range s.Gloss {
			if g.Lang != "eng" {
				continue
			}
			var wID int
			w, e2 := tx.Word.Query().
				Where(word.NameEQ(g.Text)).
				Only(ctx)
			if ent.IsNotFound(e2) {
				spaceRegex := regexp.MustCompile(`\s`)
				specialCharRegex := regexp.MustCompile(`[!?\(\)0-9]`)

				isIdioms := spaceRegex.MatchString(g.Text)
				isSpecialCharacters := specialCharRegex.MatchString(g.Text)

				w, e2 = tx.Word.Create().
					SetName(g.Text).
					SetIsIdioms(isIdioms).
					SetIsSpecialCharacters(isSpecialCharacters).
					SetRegistrationCount(0).
					Save(ctx)
			}
			if e2 != nil {
				txErr = e2
				return
			} // ← ここで return しても defer が commit/rollback

			wID = w.ID

			var wi *ent.WordInfo
			wi, e2 = tx.WordInfo. // ← ここで一度だけ宣言
						Query().
						Where(wordinfo.WordID(wID),
					wordinfo.PartOfSpeechIDEQ(posID)).
				Only(ctx)

			if ent.IsNotFound(e2) {
				wi, e2 = tx.WordInfo.Create().
					SetWordID(wID).
					SetPartOfSpeechID(posID).
					Save(ctx)
			}
			if e2 != nil {
				txErr = e2
				return
			}
			if ent.IsNotFound(e2) {
				wi, e2 = tx.WordInfo.Create().
					SetWordID(wID).
					SetPartOfSpeechID(posID).
					Save(ctx)
			}
			log.Println(wi)
			if e2 != nil {
				txErr = e2
				return
			}
			// c. JapaneseMean を全ての kanji.text で作成（無ければ kana[0]）
			if len(e.Kanji) > 0 {
				for _, kj := range e.Kanji {
					if kj.Text == "" {
						continue
					}
					if err := tx.JapaneseMean.Create().
						SetWordInfoID(wi.ID).
						SetName(kj.Text).
						OnConflict(sql.ConflictColumns("word_info_id", "name")).
						DoNothing().
						Exec(ctx); err != nil {
						return err
					}
				}
			} else if len(e.Kana) > 0 && e.Kana[0].Text != "" {
				if err := tx.JapaneseMean.Create().
					SetWordInfoID(wi.ID).
					SetName(e.Kana[0].Text).
					OnConflict(sql.ConflictColumns("word_info_id", "name")).
					DoNothing().
					Exec(ctx); err != nil {
					return err
				}
			}

		}
	}

	return nil
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

// 結果 0‑11 のどれにも当てはまらなければ 11
func mappedPosID(code string) int {
	if id, ok := posCode2ID[code]; ok {
		return id
	}
	return 11
}
