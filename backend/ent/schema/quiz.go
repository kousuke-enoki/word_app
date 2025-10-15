// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Quiz holds the schema definition for the Quiz entity.
type Quiz struct {
	ent.Schema
}

// Fields of the Quiz.
func (Quiz) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("quiz_number"),
		field.Bool("is_running").
			Default(false),
		field.Int("total_questions_count").
			Default(10),
		field.Int("correct_count").
			Default(0),
		field.Float("result_correct_rate").
			Default(0),
		field.Bool("is_save_result").
			Default(false),
		field.Int("is_registered_words").
			Default(0).
			Comment("0: all, 1: registered only, 2: unregistered only").
			Validate(func(i int) error {
				if i < 0 || i > 2 {
					return errors.New("is_registered_words must be 0, 1, or 2")
				}
				return nil
			}),
		field.Int("setting_correct_rate").
			Default(0),
		field.Int("is_idioms").
			Default(0).
			Comment("0: all, 1: idioms only, 2: not idioms only").
			Validate(func(i int) error {
				if i < 0 || i > 2 {
					return errors.New("is_idioms must be 0, 1, or 2")
				}
				return nil
			}),
		field.Int("is_special_characters").
			Default(0).
			Comment("0: all, 1: special characters only, 2: not special characters only").
			Validate(func(i int) error {
				if i < 0 || i > 2 {
					return errors.New("is_special_characters must be 0, 1, or 2")
				}
				return nil
			}),
		field.JSON("attention_level_list", []int{}),
		field.JSON("choices_pos_ids", []int{}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the Quiz.
func (Quiz) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Unique().
			Ref("quizzes").
			Field("user_id").
			Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("quiz_questions", QuizQuestion.Type),
	}
}
