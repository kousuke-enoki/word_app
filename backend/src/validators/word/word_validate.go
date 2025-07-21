package word

import (
	"regexp"

	"word_app/backend/src/models"
)

func validateWordName(name string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	// 半角アルファベットのみの正規表現
	wordNameRegex := regexp.MustCompile(`^[A-Za-z0-9'’“”"!?(),.:;#@*\-/\s]+$`)
	// word.nameの検証
	if !wordNameRegex.MatchString(name) {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "name", Message: "word.name must contain only alphabetic characters"})
	}
	if len(name) == 0 || len(name) > 100 {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "name", Message: "name must be between 0 and 41 characters"})
	}
	return fieldErrors
}

func validateWordInfos(wordInfos []models.WordInfo) []*models.FieldError {
	var fieldErrors []*models.FieldError
	// WordInfosの検証
	if len(wordInfos) < 1 || len(wordInfos) > 10 {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "wordInfos", Message: "wordInfos must contain between 1 and 10 items"})
	}

	// PartOfSpeechIDの重複チェック用マップ
	partOfSpeechIDMap := make(map[int]bool)

	// 日本語（ひらがな、カタカナ、漢字）と記号「~」のみの正規表現
	japaneseMeanRegex := regexp.MustCompile(`^[ぁ-んァ-ヶ一-龠々ー～]+$`)

	for _, wordInfo := range wordInfos {
		// JapaneseMeansの検証
		if len(wordInfo.JapaneseMeans) < 1 || len(wordInfo.JapaneseMeans) > 10 {
			fieldErrors = append(fieldErrors, &models.FieldError{Field: "japaneseMeans", Message: "japaneseMeans must contain between 1 and 10 items"})
		}

		for _, mean := range wordInfo.JapaneseMeans {
			if !japaneseMeanRegex.MatchString(mean.Name) {
				fieldErrors = append(fieldErrors, &models.FieldError{Field: "japaneseMean", Message: "japaneseMean.name must contain only Japanese characters and the '~' symbol"})
			}
		}

		// PartOfSpeechIDの重複検証
		if partOfSpeechIDMap[wordInfo.PartOfSpeechID] {
			fieldErrors = append(fieldErrors, &models.FieldError{Field: "PartOfSpeechID", Message: "duplicate PartOfSpeechID found in WordInfos"})
		}
		partOfSpeechIDMap[wordInfo.PartOfSpeechID] = true
	}

	return fieldErrors
}
