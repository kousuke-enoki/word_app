// validators/init.go
package validators

import (
	"regexp"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/word"

	"github.com/go-playground/validator/v10"
)

var V *validator.Validate
var reWord = regexp.MustCompile(`^[A-Za-z]+(?:'[A-Za-z]+)?$`)

func Init() {
	V = validator.New()
	V.SetTagName("binding")
	// フィールド用タグ (word) も一応残す
	_ = V.RegisterValidation("bulk", func(fl validator.FieldLevel) bool {
		return reWord.MatchString(fl.Field().String())
	})

	// 構造体レベル登録
	V.RegisterStructValidation(word.BulkRegisterStructLevel, models.BulkRegisterRequest{})

}

/* --- Gin が求める StructValidator --- */
type GinValidator struct{ *validator.Validate }

func (gv *GinValidator) ValidateStruct(obj any) error { return gv.Struct(obj) }
func (gv *GinValidator) Engine() any                  { return gv.Validate }
