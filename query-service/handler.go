package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nicrodriguezval/cqrs/events"
	"github.com/nicrodriguezval/cqrs/models"
	"github.com/nicrodriguezval/cqrs/repository"
	"github.com/nicrodriguezval/cqrs/search"
)

func onCreatedFeed(msg events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          msg.ID,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Println(err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
