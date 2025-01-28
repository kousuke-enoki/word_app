package word_validate_test

import (
	"testing"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"
)

func TestValidateCreateWordRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *models.CreateWordRequest
		expectErr bool
		errFields []string // エラーが期待されるフィールド
	}{
		{
			name: "ValidRequest",
			request: &models.CreateWordRequest{
				Name: "example",
				WordInfos: []models.WordInfo{
					{
						JapaneseMeans:  []models.JapaneseMean{{Name: "例"}},
						PartOfSpeechID: 1,
					},
				},
			},
			expectErr: false,
		},
		{
			name: "InvalidName",
			request: &models.CreateWordRequest{
				Name: "example123", // 無効な名前
				WordInfos: []models.WordInfo{
					{
						JapaneseMeans:  []models.JapaneseMean{{Name: "例"}},
						PartOfSpeechID: 1,
					},
				},
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "InvalidWordInfosTooShort",
			request: &models.CreateWordRequest{
				Name:      "example",
				WordInfos: []models.WordInfo{}, // WordInfosが空
			},
			expectErr: true,
			errFields: []string{"wordInfos"},
		},
		{
			name: "DuplicatePartOfSpeechID",
			request: &models.CreateWordRequest{
				Name: "example",
				WordInfos: []models.WordInfo{
					{
						JapaneseMeans:  []models.JapaneseMean{{Name: "例"}},
						PartOfSpeechID: 1,
					},
					{
						JapaneseMeans:  []models.JapaneseMean{{Name: "例"}},
						PartOfSpeechID: 1, // 重複
					},
				},
			},
			expectErr: true,
			errFields: []string{"PartOfSpeechID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := word.ValidateCreateWordRequest(tt.request)
			if tt.expectErr && len(errors) == 0 {
				t.Errorf("Expected errors, but got none")
			} else if !tt.expectErr && len(errors) > 0 {
				t.Errorf("Expected no errors, but got: %v", errors)
			}

			// フィールドごとのエラー検証
			for _, field := range tt.errFields {
				found := false
				for _, err := range errors {
					if err.Field == field {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error for field %q, but not found", field)
				}
			}
		})
	}
}
