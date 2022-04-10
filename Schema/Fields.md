# Fields
## Quick Summary
스키마에서의 필드들은 노드의 속성을 나타냅니다. 예를 들어, User라는 스키마에 age, name, username 그리고 created_at가 있다고 가정해보겠습니다.<br/>
![](../img/er_fields_properties.png)<br/>
코드는 다음과 같이 작성할 수 있습니다.
```go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

// User schema.
type User struct {
    ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age"),
        field.String("name"),
        field.String("username").
            Unique(),
        field.Time("created_at").
            Default(time.Now),
    }
}
```
필드들은 Fields메소드에 정의되어서 사용하게 됩니다. <br/>
required(not null)가 기본 설정이며 null을 허용하고 싶다면 Optional을 통해 사용할 수 있습니다.

## Types
필드에 사용할 수 있는 타입들은 다음과 같습니다.
- go에서 지원해주는 모든 숫자타입. Like int, uint8, float64, etc.
- bool
- string
- time.Time
- UUID
- []byte (SQL only).
- JSON (SQL only).
- Enum (SQL only).
- Other (SQL only).
다음과 같이 사용하실 수 있습니다.
```go
package schema

import (
    "time"
    "net/url"

    "github.com/google/uuid"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

// User schema.
type User struct {
    ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.Float("rank").
            Optional(),
        field.Bool("active").
            Default(false),
        field.String("name").
            Unique(),
        field.Time("created_at").
            Default(time.Now),
        field.JSON("url", &url.URL{}).
            Optional(),
        field.JSON("strings", []string{}).
            Optional(),
        field.Enum("state").
            Values("on", "off").
            Optional(),
        field.UUID("uuid", uuid.UUID{}).
            Default(uuid.New),
    }
}
```

## ID Field
id필드에 경우 선언이 필수적이지 않습니다. 생성 시, 자동으로 int타입으로 생성해줍니다. 만약 설정을 바꿔야한다면 codegen-option을 통해 바꿀 수 있습니다.<br/>
모든 테이블을 통틀어 ID를 다 고유하게 가져가야 한다면 스키마 마이그레이션 작업을 할떄 WithGlobalUniqueID을 사용하시길 바랍니다.<br/>
id필드에 타입을 바꾸길 원하신다면 id를 필드에 추가하신 다음 원하시는 타입으로 설정해주시면 됩니다.
예시)
```go
// Fields of the Group.
func (Group) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            StructTag(`json:"oid,omitempty"`),
    }
}

// Fields of the Blob.
func (Blob) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).
            Default(uuid.New).
            StorageKey("oid"),
    }
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").
            MaxLen(25).
            NotEmpty().
            Unique().
            Immutable(),
    }
}
```
만약 id를 커스텀 해서 쓰고 싶다면 DefaultFunc를 사용하여 구현하시면 됩니다.
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id").
            DefaultFunc(func() int64 {
                // An example of a dumb ID generator - use a production-ready alternative instead.
                return time.Now().Unix() << 8 | atomic.AddInt64(&counter, 1) % 256
            }),
    }
}
```
## DatabaseType
각각의 데이터베이스들은 Go Type을 database Type으로 매핑이 됩니다. 예시로, MYSQL은 데이터 베이스에서 Double컬럼을 Float64로 하여 진행합니다. <br/>
추가적으로 SchemaType 옵션을 통해 기본 동작을 재정의할 수 있습니다.
```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/schema/field"
)

// Card schema.
type Card struct {
    ent.Schema
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
    return []ent.Field{
        field.Float("amount").
            SchemaType(map[string]string{
                dialect.MySQL:    "decimal(6,2)",   // Override MySQL.
                dialect.Postgres: "numeric",        // Override Postgres.
            }),
    }
}
```
## Go Type
필드의 기본 타입은 기본 Go 타입입니다. 예를 들어 문자열 필드의 경우 타입은 string이고 시간 필드의 경우 타입은 time.Time입니다.<br/> 
이 Go Type방법은 기본 ent 타입을 사용자 정의 유형으로 재정의하는 옵션을 제공합니다.<br/>
사용자 정의 유형은 Go 기본 타입으로 변환할 수 있는 타입이거나 ValueScanner 인터페이스를 구현하는 타입이어야 합니다.
```go
package schema

import (
    "database/sql"

    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/schema/field"
    "github.com/shopspring/decimal"
)

// Amount is a custom Go type that's convertible to the basic float64 type.
type Amount float64

// Card schema.
type Card struct {
    ent.Schema
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
    return []ent.Field{
        field.Float("amount").
            GoType(Amount(0)),
        field.String("name").
            Optional().
            // A ValueScanner type.
            GoType(&sql.NullString{}),
        field.Enum("role").
            // A convertible type to string.
            GoType(role.Role("")),
        field.Float("decimal").
            // A ValueScanner type mixed with SchemaType.
            GoType(decimal.Decimal{}).
            SchemaType(map[string]string{
                dialect.MySQL:    "decimal(6,2)",
                dialect.Postgres: "numeric",
            }),
    }
}
```

## Other Field
Postgres 범위 타입 또는 지리공간 타입과 같은 표준 필드 유형에 적합하지 않은 필드를 말합니다.
```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/schema/field"
    
    "github.com/jackc/pgtype"
)

// User schema.
type User struct {
    ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Other("duration", &pgtype.Tstzrange{}).
            SchemaType(map[string]string{
                dialect.Postgres: "tstzrange",
            }),
    }
}
```

## Default values
Non-unique fields Default와 UpdateDefault메서드를 기본값으로 제공해주고 있으며 DefaultFunc을 통해 재정의도 가능합니다.
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Time("created_at").
            Default(time.Now),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now),
        field.String("name").
            Default("unknown"),
        field.String("cuid").
            DefaultFunc(cuid.New),
        field.JSON("dirs", []http.Dir{}).
            Default([]http.Dir{"/tmp"}),
    }
}
```
함수 호출과 같은 SQL 관련 표현식은 entsql.Annotation을 사용하여 기본값 구성에 추가할 수 있습니다.
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        // Add a new field with CURRENT_TIMESTAMP
        // as a default value to all previous rows.
        field.Time("created_at").
            Default(time.Now).
            Annotations(&entsql.Annotation{
                Default: "CURRENT_TIMESTAMP",
            }),
    }
}
```
또한, 오류를 반환하는 경우 schema-hooksDefaultFunc 를 사용하여 올바르게 처리하는 것이 좋습니다

## Validators

## Built-in Validators
entgo는 numeric,string,byte에 관련하여 간단한 내장 validator를 가지고 있습니다.
내장 validator는 다음과 같습니다.

### Numeric types
- Positive()
- Negative()
- NonNegative()
- Min(i) - Validate that the given value is > i.
- Max(i) - Validate that the given value is < i.
- Range(i, j) - Validate that the given value is within the range [i, j].

### string
- MinLen(i)
- MaxLen(i)
- Match(regexp.Regexp)
- NotEmpty

### []byte
- MaxLen(i)
- MinLen(i)
-  NotEmpty

## Optional
엔터티 생성에 필요하지 않은 필드이며 데이터베이스에서 nullable 필드로 설정됩니다. Edge와 달리 field는 required가 Default설정이며
Optional메소드를 붙여야만이 Optional기능을 사용할 수 있습니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("required_name"),
        field.String("optional_name").
            Optional(),
    }
}
```
## Nillable
필드의 값을 0 혹은 null로 처리해줍니다.해당 메소드는 Optional이 있을 때 사용가능합니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("required_name"),
        field.String("optional_name").
            Optional(),
        field.String("nillable_name").
            Optional().
            Nillable(),
    }
}
```


## Immutable
엔터티 생성 시에만 설정할 수 있는 필드입니다. 엔터티가 생성된 이후에는 변경되지 않습니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
    }
}
```

## Unique
해당 필드의 값은 고유한 값을 가져야 합니다. 고유한 값을 가져야 하기에 기본값을 가질수가 없습니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.String("nickname").
            Unique(),
    }
}
```

## Storage Key
저장소에 표현될 field의 이름을 변경해줍니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            StorageKey("old_name"),
    }
}
```
##  Indexes
다중 필드와 일부 유형의 에지에서도 정의할 수 있습니다. 그러나 이것은 현재 SQL 전용 기능이라는 점에 유의해야 합니다

## StructTag
엔티티 생성시 사용자 지정 구조체 태그를 추가할 수 있습니다. json태그에 경우 기본적으로 생성됩니다.
```go
// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            StructTag(`gqlgen:"gql_name"`),
    }
}
```
## Additional Struct Fields
일반적으로 ent를 활용하여 필드를 생성하면 다음과 같습니다.
```go
// User schema.
type User struct {
    ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Optional().
            Nillable(),
        field.String("name").
            StructTag(`gqlgen:"gql_name"`),
    }
}
```

````go
// User is the model entity for the User schema.
type User struct {
    // Age holds the value of the "age" field.
    Age  *int   `json:"age,omitempty"`
    // Name holds the value of the "name" field.
    Name string `json:"name,omitempty" gqlgen:"gql_name"`
}
````


데이터베이스에 저장되지 않는 Additional Struct Fields을 사용하려면 external Template 을 사용해야합니다.
```go
{{ define "model/fields/additional" }}
    {{- if eq $.Name "User" }}
        // StaticField defined by template.
        StaticField string `json:"static,omitempty"`
    {{- end }}
{{ end }}
```
사용하면 다음과 같습니다.
```go
// User is the model entity for the User schema.
type User struct {
    // Age holds the value of the "age" field.
    Age  *int   `json:"age,omitempty"`
    // Name holds the value of the "name" field.
    Name string `json:"name,omitempty" gqlgen:"gql_name"`
    // StaticField defined by template.
    StaticField string `json:"static,omitempty"`
}
```

## Sensitive Fields
패스워드와 같은 민감정보에 대해 사용됩니다. 해당 필드는 인코딩 시 생략이 됩니다. Sensitive Field는 StructTag를 가질 수 없습니다.
```go
// User schema.
type User struct {
    ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("password").
            Sensitive(),
    }
}
```
## Enum Fields
열거형 필드를 만들 수 있습니다 .
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("first_name"),
        field.String("last_name"),
        field.Enum("size").
            Values("big", "small"),
    }
}
```
ValueScanner string을 구현함으로써 GoType을 사용할 수 있습니다.
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("first_name"),
        field.String("last_name"),
        // A convertible type to string.
        field.Enum("shape").
            GoType(property.Shape("")),
    }
}
```
구현은 다음과 같습니다.
```go
package property

type Shape string

const (
    Triangle Shape = "TRIANGLE"
    Circle   Shape = "CIRCLE"
)

// Values provides list valid values for Enum.
func (Shape) Values() (kinds []string) {
    for _, s := range []Shape{Triangle, Circle} {
        kinds = append(kinds, string(s))
    }
    return
}

```
string타입이 아닌 다른 타입으로도 가능합니다.
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("first_name"),
        field.String("last_name"),
        // Add conversion to and from string
        field.Enum("level").
            GoType(property.Level(0)),
    }
}
```
ValueScanner 구현은 다음과 같습니다.
```go
package property

import "database/sql/driver"

type Level int

const (
    Unknown Level = iota
    Low
    High
)

func (p Level) String() string {
    switch p {
    case Low:
        return "LOW"
    case High:
        return "HIGH"
    default:
        return "UNKNOWN"
    }
}

// Values provides list valid values for Enum.
func (Level) Values() []string {
    return []string{Unknown.String(), Low.String(), High.String()}
}

// Value provides the DB a string from int.
func (p Level) Value() (driver.Value, error) {
    return p.String(), nil
}

// Scan tells our code how to read the enum into our type.
func (p *Level) Scan(val interface{}) error {
    var s string
    switch v := val.(type) {
    case nil:
        return nil
    case string:
        s = v
    case []uint8:
        s = string(v)
    }
    switch s {
    case "LOW":
        *p = Low
    case "HIGH":
        *p = High
    default:
        *p = Unknown
    }
    return nil
}
```

## Annotations
코드 생성 시, 필드 개체에 임의의 메타 데이터를 첨부하는 데 사용됩니다. 템플릿 확장은 이 메타데이터를 검색하고 템플릿 내에서 사용할 수 있습니다.
메타데이터 개체는 JSON Raw Value로 직렬화할 수 있어야 합니다.
```go
// User schema.
type User struct {
    ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Time("creation_date").
            Annotations(entgql.Annotation{
                OrderField: "CREATED_AT",
            }),
    }
}
```
## Naming Convention
snake_case로 필드명을 지어야합니다. ent에 생성된 필드들은 go convention에 따라 Pascal Case를 사용합니다.

