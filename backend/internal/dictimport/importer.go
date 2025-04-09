package dictimport

import (
	"context"
	"sync"

	"word_app/backend/ent"
	"word_app/backend/ent/word"

	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
)

// ---------- 省略 (Options / ImportJMdict は同じ) ----------

// importOne : 1 エントリを 1 トランザクションで保存
func importOne(ctx context.Context, cli *ent.Client, e JMEntry, posCache *sync.Map) (err error) {
	tx, err := cli.Tx(ctx)
	if err != nil {
		logrus.WithError(err).Error("tx start failed")
		return err
	}
	// txErr は defer から参照する
	var txErr error
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if txErr != nil {
			_ = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}
		err = txErr
	}()

	eng := firstEnglish(e)
	if eng == "" {
		return nil // 英訳が無い
	}

	// ---------- Word ----------
	wc := tx.Word.Create().
		SetName(eng).
		SetRegistrationCount(0).
		OnConflict(
			sql.ConflictColumns(word.FieldName), // UNIQUE(name)
		).
		DoNothing()

	w, txErr := wc.Save(ctx)
	if txErr != nil {
		logrus.WithFields(logrus.Fields{
			"stage": "word.create",
			"name":  eng,
			"id":    e.ID,
		}).WithError(txErr).Error("insert failed")
		return txErr
	}
	// DoNothing で既存なら w == nil
	if w == nil {
		logrus.WithField("name", eng).Debug("duplicate word skipped")
		return nil
	}

	// ---------- WordInfo & JapaneseMean ----------
	for _, s := range e.Sense {
		if len(s.PartOfSpeech) == 0 {
			continue
		}
		posCode := s.PartOfSpeech[0]

		// posCache は並列対応で sync.Map
		var posID int
		if v, ok := posCache.Load(posCode); ok {
			posID = v.(int)
		} else {
			p, perr := tx.PartOfSpeech.Create().
				SetName(posCode).
				Save(ctx)
			if perr != nil {
				txErr = perr
				logrus.WithField("pos", posCode).WithError(perr).
					Error("create part_of_speech failed")
				return txErr
			}
			posID = p.ID
			posCache.Store(posCode, posID)
		}

		wi, werr := tx.WordInfo.Create().
			SetWordID(w.ID).
			SetPartOfSpeechID(posID).
			Save(ctx)
		if werr != nil {
			txErr = werr
			logrus.WithField("stage", "wordinfo.create").WithError(werr).Error("failed")
			return txErr
		}

		if len(e.Kana) > 0 {
			if _, merr := tx.JapaneseMean.Create().
				SetWordInfoID(wi.ID).
				SetName(e.Kana[0].Text).
				Save(ctx); merr != nil {
				txErr = merr
				logrus.WithField("stage", "jmean.create").WithError(merr).Error("failed")
				return txErr
			}
		}
	}

	return nil
}
