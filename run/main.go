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
	savedUser, err := tutorial.CreateUser(context.Background(), client)
	if err != nil {
		log.Panic(err)
	}
	log.Println(savedUser)
}
