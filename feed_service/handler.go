package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nicrodriguezval/cqrs/events"
	"github.com/nicrodriguezval/cqrs/models"
	"github.com/nicrodriguezval/cqrs/repository"
	"github.com/segmentio/ksuid"
)

type createdFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createdFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now()
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		ID:        id.String(),
		Title:     req.Description,
		CreatedAt: createdAt,
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("error publishing event: %v", err)
	}

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(feed)
}
