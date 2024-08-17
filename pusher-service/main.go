package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/nicrodriguezval/cqrs/events"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	hub := NewHub()

	initEventStore(config, hub)

	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initEventStore(config Config, hub *Hub) {
	address := fmt.Sprintf("nats://%s", config.NatsAddress)

	n, err := events.NewNats(address)
	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreateFeed(func(m events.CreatedFeedMessage) {
		hub.Broadcast(
			newCreatedFeedMessage(m.ID, m.Title, m.Description, m.CreatedAt),
			nil,
		)
	})
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()
}
