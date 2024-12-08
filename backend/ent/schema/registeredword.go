package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RegisteredWord holds the schema definition for the RegisteredWord entity.
type RegisteredWord struct {
	ent.Schema
}

// Fields of the RegisteredWord.
func (RegisteredWord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Positive(),
		field.Int("word_id"),
		field.Bool("is_active").
			Default(true),
		field.Int("attention_level").
			Default(1).
			Validate(func(level int) error {
				if level < 1 || level > 5 {
					return errors.New("attention_level must be between 1 and 5")
				}
				return nil
			}),
		field.Int("test_count").
			Default(0),
		field.Int("check_count").
			Default(0),
		field.String("memo").
			Optional().
			Nillable().
			Validate(func(memo string) error {
				if len(memo) > 200 {
					return errors.New("memo must be under 200 characters")
				}
				return nil
			}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the RegisteredWord.
func (RegisteredWord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("registered_words").
			Unique().
			Field("user_id").
			Required(),
		edge.From("word", Word.Type).
			Ref("registered_words").
			Unique().
			Field("word_id").
			Required(),
		edge.To("test_questions", TestQuestion.Type),
	}
}
