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

### Query문 실행해보기
start.go파일에 a8m이라는 이름의 데이터를 찾아보는 코드를 추가하겠습니다.
```go
func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Query().
        Where(user.Name("a8m")).
        // `Only` fails if no user found,
        // or more than 1 user returned.
        Only(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed querying user: %w", err)
    }
    log.Println("user returned: ", u)
    return u, nil
}
```
추가가 되었다면 run/main.go에 QueryUser문을 추가해서 결과를 확인해보겠습니다.
코드는 다음과 같습니다.
```go
package main

import (
	"context"
	"entgo-ko/ent"
	"entgo-ko/tutorial"
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
	_, err = tutorial.CreateUser(context.Background(), client)
	if err != nil {
		log.Panic(err)
	}
	findUser, err := tutorial.QueryUser(context.Background(), client)
	if err != nil {
		log.Panic(err)
	}
	log.Println(findUser)
}
```
![](../img/스크린샷%202022-03-20%20오후%202.56.48.png)

### Edge 추가하기
유저와 relation을 맺을 두 개의 객체들을 만들겠습니다.<br/>
터미널에 다음과 같은 명령어를 입력해줍니다.
```shell
go run entgo.io/ent/cmd/ent init Car Group
```
ent/schema 폴더 안을 보면 파일이 생긴 것을 확인 하실 수 있습니다.
먼저 car와 group에 필드를 만들어 보도록하겠습니다.
```go
// entgo-ko/ent/schema/car.go
// Fields of the Car.
func (Car) Fields() []ent.Field {
    return []ent.Field{
        field.String("model"),
        field.Time("registered_at"),
    }
}
```
```go
//entgo-ko/ent/schema/group.go
// Fields of the Group.
func (Group) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            // Regexp validation for group name.
            Match(regexp.MustCompile("[a-zA-Z_]+$")),
    }
}
```
필드 정의를 다 했으니 이제 테이블간의 관계를 주도록하겠습니다. 먼저 user와 car의 관계를 만들어보도록 하겠습니다.<br/>
User에 cars라는 edge를 추가해줍니다.
```go

// Edges of the User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("cars", Car.Type),
    }
}
```
추가가 되었다면 다음 명령어를 실행해줍니다.
```shell
go generate ./ent
```

car2대와 user를 생성하는 함수를 만들겠습니다.
```go
import (
	"context"
	"entgo-ko/ent"
	"entgo-ko/ent/user"
	"fmt"
	"log"
	"time"
)
func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// Create a new car with model "Tesla".
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}
```
car와 user를 생성했습니다. 그러면 이제 생성된 카와 user를 호출해보도록하겠습니다. <br/>
start.go 파일에 QueryCars 라는 함수를 추가 하도록 하겠습니다.
```go

import (
    "context"
    "entgo-ko/ent"
    "entgo-ko/ent/car"
    "entgo-ko/ent/user"
    "fmt"
    "log"
    "time"
)

func QueryCars(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println("returned cars:", cars)

    // What about filtering specific cars.
    ford, err := a8m.QueryCars().
        Where(car.Model("Ford")).
        Only(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println(ford)
    return nil
}
```
메인 코드를 실행 해보도록 하겠습니다.
```go
package main

import (
	"context"
	"entgo-ko/ent"
	"entgo-ko/tutorial"
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
	a8m, err := tutorial.CreateCars(context.Background(), client)
	if err != nil {
		panic(err)
	}
	err = tutorial.QueryCars(context.Background(), a8m)
	if err != nil {
		panic(err)
	}
}
```
실행을 하면 다음과 같은 결과를 얻을 수 있습니다.
![](../img/스크린샷%202022-03-20%20오후%205.52.15.png)
### Inverse Edge(BackRef) 실행해보기
유저를 통해 자신의 소유의 차를 찾아 보았습니다. 코드를 실행해보시면서 차를 통해서 소유자를 찾는 거는 안될까?<br/>
라는 생각을 하셨을 겁니다. entgo에서는 소유자를 찾는 것을 할 수 있습니다. 그리고 그걸 inverse Edge라고 부릅니다.

한번 inverse Edge를 만들어 보겠습니다. <br/>

먼저 car 스키마에 owner라는 이름의  inverse edge를 만들겠습니다. 이 edge는 User스키마에 cars를 참조합니다.
```go
// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("owner", User.Type).
			Ref("cars").
			// setting the edge to unique, ensure
			// that a car can have only one owner.
			Unique(),
	}
}
```
go generate를 통해 코드 재생성을 하겠습니다.
````shell
go generate ./ent
````
위 명령어를 진행했다면 이제 inverse edge관련 쿼리를 짜보도록하겠습니다. 코드는 start.go에  추가해줍니다.
```go
func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    // Query the inverse edge.
    for _, ca := range cars {
        owner, err := ca.QueryOwner().Only(ctx)
        if err != nil {
            return fmt.Errorf("failed querying car %q owner: %w", ca.Model, err)
        }
        log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
    }
    return nil
}
```
main.go에 다음 코드를 입력해서 실행해보겠습니다.
````go
package main

import (
	"context"
	"entgo-ko/ent"
	"entgo-ko/tutorial"
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
	a8m, err := tutorial.CreateCars(context.Background(), client)
	if err != nil {
		panic(err)
	}
	err = tutorial.QueryCars(context.Background(), a8m)
	if err != nil {
		panic(err)
	}
	err = tutorial.QueryCarUsers(context.Background(), a8m)
	if err != nil {
		panic(err)
	}
}
````
