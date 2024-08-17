package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/nicrodriguezval/cqrs/database"
	"github.com/nicrodriguezval/cqrs/events"
	"github.com/nicrodriguezval/cqrs/repository"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	initRepository(config)
	initEventStore(config)

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func newRouter() *mux.Router {
  router := mux.NewRouter()
  router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)
  return router
}

func initRepository(config Config) {
	address := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresDB, config.PostgresUser, config.PostgresPassword)

	postgresRepo, err := database.NewPostgresRepository(address)
	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepostiory(postgresRepo)
}

func initEventStore(config Config) {
	address := fmt.Sprintf("nats://%s", config.NatsAddress)

	n, err := events.NewNats(address)
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()
}
