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
			Validate(func(s string) error {
				// 半角アルファベットのみ許可
				match, _ := regexp.MatchString(`^[a-zA-Z]+$`, s)
				if !match {
					return errors.New("name must contain only alphabetic characters")
				}
				return nil
			}),
		field.String("voice_id").
			Optional().
			Nillable(),
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
	}
}
