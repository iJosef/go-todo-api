package main

import (
	"context"
	"github.com/iJosef/go-todo-api/internal/db"
	"github.com/iJosef/go-todo-api/internal/todo"
	"github.com/iJosef/go-todo-api/internal/transport"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number %q: %v", portStr, err)
	}

	// Create a context with timeout for DB connection
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d, err := db.New(user, password, dbName, host, port)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	svc := todo.NewService(d)
	server := transport.NewServer(svc)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

}
