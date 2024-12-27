package schema

import (
	"errors"
	"regexp"
	"time"

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
			NotEmpty(),
		field.String("name").
			Default("JohnDoe").
			Comment("Name of the user.\n If not specified, defaults to \"John Doe\".").
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
		field.Bool("isAdmin").
			Default(false),
		field.Bool("isRoot").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("registered_words", RegisteredWord.Type),
		edge.To("user_config", UserConfig.Type).
			Unique(),
		edge.To("exams", Exam.Type),
	}
}
