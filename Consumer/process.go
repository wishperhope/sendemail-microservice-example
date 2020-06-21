package main

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

// Process Recived Message
func (s *Server) process(messages <-chan *message.Message) {
	for msg := range messages {

		// send Email
		result, err := send(string(msg.Payload))
		if err != nil {
			log.Println(err)
		}

		log.Printf("received message: %s, email: %s result: %s", msg.UUID, string(msg.Payload), result)

		msg.Ack()

	}
}
