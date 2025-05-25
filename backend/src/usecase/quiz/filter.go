package quiz

import (
	"word_app/backend/src/domain/quiz"
	"word_app/backend/src/models"
)

func buildFilter(in *models.CreateQuizDTO) quiz.WordFilter {
	return quiz.WordFilter{
		PartsOfSpeech:   in.PartsOfSpeeches,
		RegisteredMode:  in.IsRegisteredWords,
		MaxCorrectRate:  in.CorrectRate,
		AttentionLevels: in.AttentionLevelList,
		IncludeIdioms:   in.IsIdioms,
		IncludeSpecial:  in.IsSpecialCharacters,
	}
}
