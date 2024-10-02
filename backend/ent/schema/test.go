package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Test holds the schema definition for the Test entity.
type Test struct {
	ent.Schema
}

// Fields of the Test.
func (Test) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("total_questions").
			Default(10),
		field.Int("correct_count").
			Default(0),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Test.
func (Test) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Unique().
			Ref("tests").
			Field("user_id").
			Required(),
		edge.To("test_questions", TestQuestion.Type),
	}
}
