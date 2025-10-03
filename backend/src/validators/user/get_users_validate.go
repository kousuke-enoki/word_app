package user

import (
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/models"
)

func ValidateUserListRequest(req *user.ListUsersInput) []models.FieldError {
	var fieldErrors []models.FieldError

	// Validate each field
	fieldErrors = append(fieldErrors, validateUserID(req.ViewerID)...)
	fieldErrors = append(fieldErrors, validateSearch(req.Search)...)
	fieldErrors = append(fieldErrors, validateSortBy(req.SortBy)...)
	fieldErrors = append(fieldErrors, validateOrder(req.Order)...)
	fieldErrors = append(fieldErrors, validatePagination(req.Page, req.Limit)...)

	return fieldErrors
}

func validateUserID(userID int) []models.FieldError {
	var fieldErrors []models.FieldError
	if userID <= 0 {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "userID", Message: "userID must be a positive integer"})
	}
	return fieldErrors
}

func validateSearch(search string) []models.FieldError {
	var fieldErrors []models.FieldError
	if len(search) > 100 {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "search", Message: "search query must not exceed 100 characters"})
	}
	return fieldErrors
}

func validateSortBy(sortBy string) []models.FieldError {
	var fieldErrors []models.FieldError
	allowedSortBy := map[string]bool{"name": true, "role": true, "email": true}
	if !allowedSortBy[sortBy] {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "sortBy", Message: "sortBy must be 'name' or 'registration_count'"})
	}
	return fieldErrors
}

func validateOrder(order string) []models.FieldError {
	var fieldErrors []models.FieldError
	allowedOrder := map[string]bool{"asc": true, "desc": true}
	if !allowedOrder[order] {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "order", Message: "order must be 'asc' or 'desc'"})
	}
	return fieldErrors
}

func validatePagination(page, limit int) []models.FieldError {
	var fieldErrors []models.FieldError
	if page <= 0 {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "page", Message: "page must be a positive integer"})
	}
	if limit <= 0 || limit > 100 {
		fieldErrors = append(fieldErrors, models.FieldError{Field: "limit", Message: "limit must be between 1 and 100"})
	}
	return fieldErrors
}
