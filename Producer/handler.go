package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/monoculum/formam"
)

func (s *Server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err = r.ParseForm(); err != nil {
			log.Println(r.FormValue("host"))
			http.Error(w, "Bad Request", 400)
			return
		}

		// validate field
		field := [8]string{"To", "From", "Subject", "Body", "Host", "Port", "Username", "Password"}
		for _, value := range field {
			if r.FormValue(value) == "" {
				log.Println("Missing Parameter'", value, "' = ", r.FormValue(value))
				http.Error(w, "Bad Request", 400)
				return
			}
		}

		job := Email{}

		dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
		if err := dec.Decode(r.Form, &job); err != nil {
			log.Fatal("Wrong Parameter Request", err)
		}

		res, err := json.Marshal(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send to nats
		msg := message.NewMessage(watermill.NewUUID(), res)
		if err = s.emailPublisher.Publish("sendEmail.topic", msg); err != nil {
			log.Fatal("Cannot Send Content To Nats Server")
		}
		log.Println("Sending Message : ", string(res))

		w.Header().Set("Content-Type", "application/json")
		// Return fixed response change this if nescessary
		_, err = w.Write([]byte("\"{success:true}\""))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
