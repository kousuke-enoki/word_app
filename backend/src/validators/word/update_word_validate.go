package word

import (
	"word_app/backend/src/models"
)

// バリデーション関数
func ValidateUpdateWordRequest(req *models.UpdateWordRequest) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, validateWordName(req.Name)...)
	fieldErrors = append(fieldErrors, validateWordInfos(req.WordInfos)...)

	return fieldErrors
}
