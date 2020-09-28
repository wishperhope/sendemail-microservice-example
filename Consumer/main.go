package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/joho/godotenv"
)

// Server Struct hold nats connection
type Server struct {
	db              *sql.DB
	emailSubscriber *nats.StreamingSubscriber
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env File")
	}

	s := Server{}
	err = s.setup()
	if err != nil {
		panic(err)
	}

	// Accept topic from producer
	messages, err := s.emailSubscriber.Subscribe(context.Background(), "sendEmail.topic")
	if err != nil {
		panic(err)
	}

	// Process all message
	s.process(messages)
}
