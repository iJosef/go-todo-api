package main

import (
	"github.com/iJosef/go-todo-api/internal/db"
	"github.com/iJosef/go-todo-api/internal/todo"
	"github.com/iJosef/go-todo-api/internal/transport"
	"log"
)

func main() {
	d, err := db.New("postgres", "example", "postgres", "localhost", 5433)
	if err != nil {
		log.Fatal(err)
	}

	svc := todo.NewService(d)
	server := transport.NewServer(svc)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

}
