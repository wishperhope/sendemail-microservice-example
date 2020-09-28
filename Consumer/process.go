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

		log.Printf("received message: %s, result: %s", msg.UUID, result)

		msg.Ack()

		query := "UPDATE send_email SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"

		_, err = s.db.Exec(query, result, msg.UUID)
		if err != nil {
			log.Println("Cannot Update Status of Email")
		}

	}
}
