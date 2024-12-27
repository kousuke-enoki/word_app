package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ExamQuestion holds the schema definition for the ExamQuestion entity.
type ExamQuestion struct {
	ent.Schema
}

// Fields of the ExamQuestion.
func (ExamQuestion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("exam_id"),
		field.Int("correct_jpm_id"),
		field.JSON("choices_jpm_ids", []int{}),
		field.String("answer").
			Nillable().
			Default(""),
		field.Bool("is_correct").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Default(nil),
	}
}

// Edges of the ExamQuestion.
func (ExamQuestion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("exam", Exam.Type).
			Ref("exam_questions").
			Field("exam_id").
			Unique().
			Required(),
		edge.From("japanese_mean", JapaneseMean.Type).
			Ref("exam_questions").
			Unique().
			Field("correct_jpm_id").
			Required(),
	}
}
