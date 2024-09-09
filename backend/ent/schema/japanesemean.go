package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// JapaneseMean holds the schema definition for the JapaneseMean entity.
type JapaneseMean struct {
	ent.Schema
}

// Fields of the JapaneseMean.
func (JapaneseMean) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.Int("word_info_id").
			Positive(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the JapaneseMean.
func (JapaneseMean) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("word_info", WordInfo.Type).
				Ref("japanese_means").
				Field("word_info_id").
				Unique().
				Required(),
	}
}
