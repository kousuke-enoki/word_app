package word

import (
	"regexp"
	"strconv"

	"word_app/backend/src/models"

	"github.com/go-playground/validator/v10"
)

var reWord = regexp.MustCompile(`^[A-Za-z]+(?:'[A-Za-z]+)?$`)

func BulkRegisterStructLevel(sl validator.StructLevel) {
	req := sl.Current().Interface().(models.BulkRegisterRequest)

	for i, w := range req.Words {
		if len(w) == 0 || len(w) > 40 || !reWord.MatchString(w) {
			fieldName := "Words[" + strconv.Itoa(i) + "]"
			sl.ReportError(w, fieldName, fieldName, "word_length", "")
		}
	}
}
