package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"regexp"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Positive().Unique().Comment("User ID"),
		field.UUID("uuid", uuid.UUID{}).Default(uuid.New),
		field.String("username").Unique().MaxLen(40).Comment("UserName").Match(regexp.MustCompile("[a-zA-Z_]+$")),
		field.String("nickname").Default("").MaxLen(35).Comment("NickName"),
		field.String("password").Default("").MaxLen(50).Comment("Password"),
		field.Bool("active").Default(false),
		field.Enum("state").Values("on", "off").Optional(),
		field.Time("created_at").Default(time.Now).Comment("create time"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "username").
			Unique(),
	}
}
