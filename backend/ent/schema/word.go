package schema

import (
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
			NotEmpty(),
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
