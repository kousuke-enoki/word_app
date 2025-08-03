// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"errors"
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Word holds the schema definition for the Word entity.
type Word struct {
	ent.Schema
}

// Fields of the Word.
func (Word) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Validate(func(name string) error {
				valid := regexp.MustCompile(`^[A-Za-z0-9'’“”"!?(),.:;#@*\-/\s]+$`).MatchString
				if !valid(name) {
					return errors.New("invalid word name")
				}
				if len(name) == 0 || len(name) > 100 {
					return errors.New("name must be between 0 and 100 characters")
				}
				return nil
			}).
			Unique(),
		field.String("voice_id").
			Optional().
			Nillable(),
		field.Bool("is_idioms").
			Default(false),
		field.Bool("is_special_characters").
			Default(false),
		field.Int("registration_count").
			Default(0),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Word.
func (Word) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("word_infos", WordInfo.Type),
		edge.To("registered_words", RegisteredWord.Type),
		edge.To("quiz_questions", QuizQuestion.Type),
	}
}
