package word_validate_test

import (
	"testing"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"
)

func TestValidateWordListRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *models.WordListRequest
		expectErr bool
		errFields []string // エラーが期待されるフィールド
	}{
		{
			name: "ValidRequest",
			request: &models.WordListRequest{
				UserID: 1,
				Search: "example",
				SortBy: "name",
				Order:  "asc",
				Page:   1,
				Limit:  10,
			},
			expectErr: false,
		},
		{
			name: "InvalidUserID",
			request: &models.WordListRequest{
				UserID: -1,
				Search: "example",
				SortBy: "name",
				Order:  "asc",
				Page:   1,
				Limit:  10,
			},
			expectErr: true,
			errFields: []string{"userID"},
		},
		{
			name: "SearchTooLong",
			request: &models.WordListRequest{
				UserID: 1,
				Search: "a very long search string that exceeds the allowed limit of 100 characters ...........................................",
				SortBy: "name",
				Order:  "asc",
				Page:   1,
				Limit:  10,
			},
			expectErr: true,
			errFields: []string{"search"},
		},
		{
			name: "InvalidSortBy",
			request: &models.WordListRequest{
				UserID: 1,
				Search: "example",
				SortBy: "invalid_sort",
				Order:  "asc",
				Page:   1,
				Limit:  10,
			},
			expectErr: true,
			errFields: []string{"sortBy"},
		},
		{
			name: "InvalidOrder",
			request: &models.WordListRequest{
				UserID: 1,
				Search: "example",
				SortBy: "name",
				Order:  "random",
				Page:   1,
				Limit:  10,
			},
			expectErr: true,
			errFields: []string{"order"},
		},
		{
			name: "InvalidPagination",
			request: &models.WordListRequest{
				UserID: 1,
				Search: "example",
				SortBy: "name",
				Order:  "asc",
				Page:   -1,
				Limit:  200,
			},
			expectErr: true,
			errFields: []string{"page", "limit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := word.ValidateWordListRequest(tt.request)
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
