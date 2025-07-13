package word_validate_test

import (
	"testing"

	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"
)

func TestValidateSaveMemo(t *testing.T) {
	tests := []struct {
		name       string
		req        *models.SaveMemoRequest
		wantErrors []*models.FieldError
	}{
		{
			name: "Valid memo",
			req: &models.SaveMemoRequest{
				Memo: "This is a valid memo.",
			},
			wantErrors: nil,
		},
		{
			name: "Empty memo",
			req: &models.SaveMemoRequest{
				Memo: "",
			},
			wantErrors: nil, // 空のメモはエラーにならない仕様
		},
		{
			name: "Memo exceeds max length",
			req: &models.SaveMemoRequest{
				Memo: string(make([]byte, 201)), // 201文字のメモ
			},
			wantErrors: []*models.FieldError{
				{Field: "memo", Message: "memo must be less than 200 characters"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrors := word.ValidateSaveMemo(tt.req)

			if len(gotErrors) != len(tt.wantErrors) {
				t.Errorf("unexpected number of errors, got %d, want %d", len(gotErrors), len(tt.wantErrors))
			}

			for i, err := range gotErrors {
				if err.Field != tt.wantErrors[i].Field || err.Message != tt.wantErrors[i].Message {
					t.Errorf("unexpected error at index %d, got %+v, want %+v", i, err, tt.wantErrors[i])
				}
			}
		})
	}
}
