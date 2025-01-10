package main

import (
	"go-api/internal/db"
	"go-api/internal/todo"
	"go-api/internal/transport"
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
