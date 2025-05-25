package quiz

import (
	"encoding/json"
	"word_app/backend/src/models"
	"word_app/backend/src/repository"
)

// DB に書き込む最低限のレコード
type QuizRecord = repository.QuizRecord // 別名 import でも OK

func toQuizRecord(in *models.CreateQuizDTO) QuizRecord {
	return QuizRecord{
		QuestionCount:      in.QuestionCount,
		IsSaveResult:       in.IsSaveResult,
		IsRegisteredWords:  in.IsRegisteredWords,
		SettingCorrectRate: in.CorrectRate,
		IsIdioms:           in.IsIdioms,
		IsSpecialChars:     in.IsSpecialCharacters,
		AttentionLevels:    in.AttentionLevelList,
		ChoicePosIDs:       in.PartsOfSpeeches,
	}
}

// DB レコード → フロント返却 DTO
func toNextQuestionDTO(q repository.QuizQuestionRecord) models.NextQuestion {
	return models.NextQuestion{
		QuizID:         0, // 呼び出し元で埋める
		QuestionNumber: q.QuestionNumber,
		WordName:       q.WordName, // Word 名は Join 取得 or 別途 map 引数で貰う
		ChoicesJpms:    decodeChoices(q.ChoicesJSON),
	}
}

// JSON []byte → []ChoiceJpm へのユーティリティ
func decodeChoices(b []byte) []models.ChoiceJpm {
	var out []models.ChoiceJpm
	_ = json.Unmarshal(b, &out)
	return out
}

func mustDecodeChoices(b []byte) []models.ChoiceJpm {
	var out []models.ChoiceJpm
	_ = json.Unmarshal(b, &out) // 失敗しても空配列で返す想定
	return out
}
