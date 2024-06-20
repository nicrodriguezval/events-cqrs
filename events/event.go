package events

import (
	"context"

	"github.com/nicrodriguezval/cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, m *models.Feed) error
	SubscribeCreatedFeed(ctx context.Context) (<-chan *models.Feed, error)
	OnCreatedFeed(f func(CreatedFeedMessage)) error
}

var eventStore EventStore

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, m *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, m)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan *models.Feed, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreatedFeed(f func(CreatedFeedMessage)) error {
  return eventStore.OnCreatedFeed(f)
}
