package word

import "word_app/backend/src/models"

func ValidateSaveMemo(req *models.SaveMemoRequest) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, validateMemoLength(req.Memo)...)

	return fieldErrors
}

func validateMemoLength(memo string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	if len(memo) > 200 {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "memo", Message: "memo must be less than 200 characters"})
	}
	return fieldErrors
}
