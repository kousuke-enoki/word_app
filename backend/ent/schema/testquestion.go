package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestQuestion holds the schema definition for the TestQuestion entity.
type TestQuestion struct {
	ent.Schema
}

// Fields of the TestQuestion.
func (TestQuestion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("test_id"),
		field.Int("registered_word_id"),
		field.Bool("is_correct").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the TestQuestion.
func (TestQuestion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("test", Test.Type).
			Ref("test_questions").
			Field("test_id").
			Unique().
			Required(),
		edge.From("registered_word", RegisteredWord.Type).
			Ref("test_questions").
			Unique().
			Field("registered_word_id").
			Required(),
	}
}
