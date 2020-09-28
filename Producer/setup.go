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
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

// Setup all connection including to nats server and route
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

	s.emailPublisher, err = nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: natsClusterID,
			ClientID:  natsCleintID + strconv.Itoa(rand.Intn(100)),
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
