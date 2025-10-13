// schema/user_daily_usage.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserDailyUsage struct{ ent.Schema }

func (UserDailyUsage) Fields() []ent.Field {
	return []ent.Field{
		field.Time("last_reset_date").
			Comment("JSTでの当日0時。これと今日JSTを比較してリセット判定する"),
		field.Int("quiz_count").Default(0),
		field.Int("bulk_count").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (UserDailyUsage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("user_daily_usage").
			Unique().
			Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}
