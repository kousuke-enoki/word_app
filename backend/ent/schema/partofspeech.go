package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// PartOfSpeech holds the schema definition for the PartOfSpeech entity.
type PartOfSpeech struct {
	ent.Schema
}

// Fields of the PartOfSpeech.
func (PartOfSpeech) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the PartOfSpeech.
func (PartOfSpeech) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("word_infos", WordInfo.Type),
	}
}
