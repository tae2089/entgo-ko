## 스키마 생성해보기
```shell
go run entgo.io/ent/cmd/ent init User
```
### 생성된 스키마 파일 확인하기
entgo-ko/ent/schema/User.go파일에 들어가시면 다음과 같은 코드를 확인하실 수 있습니다.
```go
package schema

import "entgo.io/ent"

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return nil
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
```
User 스키마에 필드를 추가해보겠습니다. 여기서 필드는 테이블의 컬럼이라고 생각하시면 됩니다.
```go

package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.String("name").
            Default("unknown"),
    }
}
```
스키마에 관련한 함수들을 쓰기 위해 <span style="background-color:gray;color:orange;">go generate</span>를 사용하도록 하겠습니다.
생성을 하게 되면 다음과 같은 구조를 보실 수 있습니다.
```
ent
├── client.go
├── config.go
├── context.go
├── ent.go
├── generate.go
├── mutation.go
... truncated
├── schema
│   └── user.go
├── tx.go
├── user
│   ├── user.go
│   └── where.go
├── user.go
├── user_create.go
├── user_delete.go
├── user_query.go
└── user_update.go
```
현재 저희의 코드는 테이블과 연동이 되지 않았을 뿐더러 필드만 정의된 코드입니다.<br/>
이 코드를 DB와 연동하는 과정을 해보겠습니다. <br/>
run이라는 경로를 추가하고 그안에 main.go파일을 생성해줍니다.
```go
package main

import (
    "context"
    "log"

    "<project>/ent"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
    if err != nil {
        log.Fatalf("failed opening connection to sqlite: %v", err)
    }
    defer client.Close()
    // Run the auto migration tool.
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
}
```
위 코드에서 알아야 하는 부분은 ent.Open을 통해서 DB와 커넥션을 맺고 client.Schema.Create를 통해 테이블 생성한 점을 깊게 보시면 됩니다.<br/>

다음으로 User 데이터를 DB에 insert하는 것을 해보도록하겠습니다.<br/>
tutorial이라는 경로를 만들어주고 그 안에 start.go라는 파일을 만들어주고 거기에 다음과 같은 코드를 넣어줍니다.
```go
package tutorial

import (
	"context"
	"entgo-ko/ent"
	"fmt"
	"log"
)

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
```
CreateUser함수를 사용해보도록 하겠습니다. tutorial.go파일로 돌아와 다음과 같이 입력해줍니다.
```go
package main

import (
	"context"
	"entgo-ko/ent"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	savedUser, err := CreateUser(context.Background(), client)
	if err != nil {
		log.Panic(err)
	}
	log.Println(savedUser)
}
```
위 코드를 실행하면 다음과 같은 결과를 얻을 수 있습니다.<br/>
![](../img/스크린샷%202022-03-19%20오후%208.28.53.png)
