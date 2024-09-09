package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
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
			field.Int("voice_id").
				Positive(),
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
		}
}
