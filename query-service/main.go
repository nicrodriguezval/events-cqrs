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
	"github.com/nicrodriguezval/cqrs/search"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func mai() {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	initRepository(config)
	initEventStore(config)
	initElasticSearch(config)

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
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

	err = n.OnCreateFeed(onCreatedFeed)
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()
}

func initElasticSearch(config Config) {
	es, err := search.NewElastic(fmt.Sprintf("http://%s", config.ElasticSearchAddress))
	if err != nil {
		log.Fatal(err)
	}

	search.SetRepository(es)

	defer search.Close()
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feeds", listFeedsHandler).Methods(http.MethodGet)
	router.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
	return router
}
