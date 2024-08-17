package events

import (
	"context"

	"github.com/nicrodriguezval/cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, m *models.Feed) error
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	OnCreateFeed(f func(CreatedFeedMessage)) error
}

var eventStore EventStore

func SetEventStore(es EventStore) {
  eventStore = es
}

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, m *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, m)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreateFeed(f func(CreatedFeedMessage)) error {
  return eventStore.OnCreateFeed(f)
}
