// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"time"

	"word_app/backend/src/models"

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
		field.String("wordName").
			NotEmpty(),
		field.Int("pos_id").
			Positive(),
		field.Int("correct_jpm_id").
			Comment("correct japanese mean id"),
		field.JSON("choices_jpms", []models.ChoiceJpm{}).
			Comment("4 つの選択肢 (正解 + 誤答)"),
		field.Int("answer_jpm_id").
			Optional().
			Nillable().
			Comment("answered japanese mean ids"),
		field.Bool("is_correct").
			Optional().
			Nillable(),
		field.Time("answered_at").
			Optional().
			Nillable(),
		field.Int("time_ms").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
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
