package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	"github.com/nicrodriguezval/cqrs/models"
)

type NatsEventStore struct {
  conn *nats.Conn
  feedCreatedSub *nats.Subscription
  feedCreatedChan chan CreatedFeedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
  conn, err := nats.Connect(url)
  if err != nil {
    return nil, err
  }
  return &NatsEventStore{
    conn: conn,
    feedCreatedChan: make(chan CreatedFeedMessage),
  }, nil
}

func (n *NatsEventStore) Close() {
  if n.conn != nil {
    n.conn.Close()
  }
  if n.feedCreatedSub != nil {
    n.feedCreatedSub.Unsubscribe()
  }
  close(n.feedCreatedChan)
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
  b := bytes.Buffer{}
  if err := gob.NewEncoder(&b).Encode(m); err != nil {
    return nil, err
  }

  return b.Bytes(), nil
}

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
  msg := CreatedFeedMessage{
    ID: feed.ID,
    Title: feed.Title,
    Description: feed.Description,
    CreatedAt: feed.CreatedAt,
  }

  data, err := n.encodeMessage(msg)
  if err != nil {
    return err
  }

  return n.conn.Publish(msg.Type(), data)
}
