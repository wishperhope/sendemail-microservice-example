package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

		jobJSON, err := json.Marshal(&job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uuid := watermill.NewUUID()

		_, err = s.db.Exec("INSERT INTO send_email (id, `to`, `from`, subject, status, note) VALUES  (?, ?, ?, ?, ?, ?)", uuid, job.To, job.From, job.Subject, "waiting", job.Note)
		if err != nil {
			log.Fatal("Cannot Insert New Record", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		result := Job{
			ID:        uuid,
			To:        job.To,
			From:      job.From,
			Status:    "waiting",
			Subject:   job.Subject,
			CreatedAt: time.Now(),
			Note:      job.Note,
		}

		// send to nats
		msg := message.NewMessage(uuid, jobJSON)
		if err = s.emailPublisher.Publish("sendEmail.topic", msg); err != nil {
			log.Fatal("Cannot Send Content To Nats Server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		res, err := json.Marshal(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
