package schema

import (
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
		field.Int("user_id"),
		field.Int("word_info_id"),
		field.Bool("is_active").
			Default(true),
		field.Int("test_count").
			Default(0),
		field.Int("check_count").
			Default(0),
		field.String("memo").
			Optional().
			Nillable(),
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
		edge.From("word_info", WordInfo.Type).
			Ref("registered_words").
			Unique().
			Field("word_info_id").
			Required(),
		edge.To("test_questions", TestQuestion.Type),
	}
}
