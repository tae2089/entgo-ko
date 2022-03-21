package main

import (
	"ariga.io/sqlcomment"
	"bytes"
	"context"
	"encoding/json"
	"entgo-ko/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type CustomCommenter struct{}

func (mcc CustomCommenter) Tag(ctx context.Context) sqlcomment.Tags {
	return sqlcomment.Tags{
		"key": "value",
	}
}
func main() {
	//client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// Create db driver.
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create sqlcomment driver which wraps sqlite driver.
	drv := sqlcomment.NewDriver(dialect.Debug(db),
		sqlcomment.WithDriverVerTag(),
		sqlcomment.WithTags(sqlcomment.Tags{
			sqlcomment.KeyApplication: "my-app",
			sqlcomment.KeyFramework:   "net/http",
		}),
	)
	client := ent.NewClient(ent.Driver(drv))
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	//_, err = tutorial.CreateUser(context.Background(), client)
	//if err != nil {
	//	log.Panic(err)
	//}
	//findUser, err := tutorial.QueryUser(context.Background(), client)
	//if err != nil {
	//	log.Panic(err)
	//}
	//log.Println(findUser)
	//a8m, err := tutorial.CreateCars(context.Background(), client)
	client.User.
		Create().SetName("sd").
		Save(context.Background())
	//
	//if err != nil {
	//	panic(err)
	//}
	//_, err = tutorial.QueryUser(context.Background(), client)
	//err = tutorial.QueryCars(context.Background(), a8m)
	//if err != nil {
	//	panic(err)
	//}
	//err = tutorial.QueryCarUsers(context.Background(), a8m)
	//if err != nil {
	//	panic(err)
	//}
}

func QueryTest(v ...interface{}) {
	//jsonData := make(map[string]interface{})
	//jsonData["query"] = v
	doc, _ := json.Marshal(v)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, doc, "", "    "); err != nil {
		panic(err)
	}
	_, err := os.Stdout.Write(prettyJSON.Bytes())
	if err != nil {
		panic(err)
	}
}
