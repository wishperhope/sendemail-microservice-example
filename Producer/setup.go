package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

// Setup all connection including to nats server and route
func (s *Server) setup() error {

	// Setup Publisher
	natsAddr := os.Getenv("NATS_ADDR")
	natsPort := os.Getenv("NATS_PORT")
	natsUser := os.Getenv("NATS_USER")
	natsPassword := os.Getenv("NATS_PASSWORD")
	natsClusterId := os.Getenv("NATS_CLUSTER_ID")
	natsCleintId := os.Getenv("NATS_CLIENT_ID")

	var err error
	s.emailPublisher, err = nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: natsClusterId,
			ClientID:  natsCleintId + strconv.Itoa(rand.Intn(100)),
			StanOptions: []stan.Option{
				stan.NatsURL("nats://" + natsUser + ":" + natsPassword + "@" + natsAddr + ":" + natsPort),
			},
			Marshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		log.Fatal("Cannot Open Connection To Nats Server : ", err)
		return err
	}

	s.token = os.Getenv("APP_KEY")
	if s.token == "" {
		log.Fatal("APP KEY NOT SET")
		panic("APP KEY NOT SET")
	}
	s.router = mux.NewRouter()
	s.route()
	return nil
}
