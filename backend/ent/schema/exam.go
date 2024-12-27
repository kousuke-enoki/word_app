package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Exam holds the schema definition for the Exam entity.
type Exam struct {
	ent.Schema
}

// Fields of the Exam.
func (Exam) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("total_questions").
			Default(10),
		field.Int("correct_count").
			Default(0),
		field.Bool("is_running").
			Default(false),
		field.String("target_word_types").
			NotEmpty(),
		field.JSON("choices_pos_ids", []int{}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Default(nil),
	}
}

// Edges of the Exam.
func (Exam) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Unique().
			Ref("exams").
			Field("user_id").
			Required(),
		edge.To("exam_questions", ExamQuestion.Type),
	}
}
