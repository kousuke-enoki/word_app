package schema

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Unique().
			NotEmpty().
			Validate(func(email string) error {
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRegex.MatchString(email) {
					return errors.New("invalid email format")
				}
				return nil
			}),
		field.String("password").
			Sensitive().
			NotEmpty().
			Validate(func(password string) error {
				if len(password) < 8 || len(password) > 30 {
					return errors.New("password must be between 8 and 30 characters")
				}
				var hasUpper, hasLower, hasNumber, hasSpecial bool
				for _, ch := range password {
					switch {
					case unicode.IsUpper(ch):
						hasUpper = true
					case unicode.IsLower(ch):
						hasLower = true
					case unicode.IsNumber(ch):
						hasNumber = true
					case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
						hasSpecial = true
					}
				}
				if !(hasUpper && hasLower && hasNumber && hasSpecial) {
					return errors.New("password must include at least one uppercase letter, one lowercase letter, one number, and one special character")
				}
				return nil
			}),
		field.String("name").
			NotEmpty().
			Validate(func(name string) error {
				if len(name) < 3 || len(name) > 20 {
					return errors.New("name must be between 3 and 20 characters")
				}
				return nil
			}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Bool("admin").
			Default(false),
		field.Bool("root").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("registered_words", RegisteredWord.Type),
		edge.To("tests", Test.Type),
	}
}
