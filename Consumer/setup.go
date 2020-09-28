package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nats-io/stan.go"
)

// Setup Envoriment Nats
func (s *Server) setup() error {

	// Setup DB
	mysqlAddr := os.Getenv("MYSQL_ADDR")
	mysqlPort := os.Getenv("MSYQL_PORT")
	mysqlDB := os.Getenv("MYSQL_DB")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")

	var err error
	s.db, err = sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@("+mysqlAddr+":"+mysqlPort+")/"+mysqlDB+"?parseTime=true")
	s.db.SetMaxIdleConns(0)
	s.db.SetMaxOpenConns(151)
	s.db.SetConnMaxLifetime(time.Second * 60)
	if err != nil {
		log.Fatal("Cannot Open Connection To DB")
		return err
	}

	// Setup Publisher
	natsAddr := os.Getenv("NATS_ADDR")
	natsPort := os.Getenv("NATS_PORT")
	natsUser := os.Getenv("NATS_USER")
	natsPassword := os.Getenv("NATS_PASSWORD")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsCleintID := os.Getenv("NATS_CLIENT_ID")

	s.emailSubscriber, err = nats.NewStreamingSubscriber(
		nats.StreamingSubscriberConfig{
			ClusterID: natsClusterID,
			// plus random to make sure client id not have confilct if replicated
			ClientID:         natsCleintID + strconv.Itoa(rand.Intn(1000)),
			QueueGroup:       "checkEmail",
			DurableName:      "email-durable",
			SubscribersCount: 4, // how many goroutines should consume messages
			CloseTimeout:     time.Minute,
			AckWaitTimeout:   time.Second * 30,
			StanOptions: []stan.Option{
				stan.NatsURL("nats://" + natsUser + ":" + natsPassword + "@" + natsAddr + ":" + natsPort),
			},
			Unmarshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		log.Fatal("Cannot Open Connection To Nats Server : ", err)
		return err
	}

	return nil
}
