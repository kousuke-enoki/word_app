package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QuizQuestion holds the schema definition for the QuizQuestion entity.
type QuizQuestion struct {
	ent.Schema
}

// Fields of the QuizQuestion.
func (QuizQuestion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quiz_id"),
		field.Int("question_number"),
		field.Int("word_id").
			Positive(),
		field.Int("correct_jpm_id"),
		field.JSON("choices_jpm_ids", []int{}),
		field.Int("answer_jpm_id").
			Nillable(),
		field.Bool("is_correct").
			Nillable(),
		field.Time("answered_at").
			Nillable(),
		field.Int("time_ms").
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Default(nil),
	}
}

// Edges of the QuizQuestion.
func (QuizQuestion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("quiz", Quiz.Type).
			Ref("quiz_questions").
			Field("quiz_id").
			Unique().
			Required(),
		edge.From("word", Word.Type).
			Ref("quiz_questions").
			Field("word_id").
			Unique().
			Required(),
		edge.From("japanese_mean", JapaneseMean.Type).
			Ref("quiz_questions").
			Unique().
			Field("correct_jpm_id").
			Required(),
	}
}
