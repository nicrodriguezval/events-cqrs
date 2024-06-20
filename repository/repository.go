package repository

import (
	"context"

	"github.com/nicrodriguezval/cqrs/models"
)

type Repository interface {
  Close()
  InsertFeed(ctx context.Context, feed *models.Feed) error
  ListFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepostiory(r Repository) {
  repository = r
}

func Close() {
  repository.Close()
}

func InsertFeed(ctx context.Context, feed *models.Feed) error {
  return repository.InsertFeed(ctx, feed)
}

func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
  return repository.ListFeeds(ctx)
}
