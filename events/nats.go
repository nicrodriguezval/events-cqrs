package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	"github.com/nicrodriguezval/cqrs/models"
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{
		conn:            conn,
		feedCreatedChan: make(chan CreatedFeedMessage, 64),
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

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

func (n *NatsEventStore) OnCreateFeed(f func(CreatedFeedMessage)) (err error) {
	msg := CreatedFeedMessage{}

	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		if err := n.DecodeMessage(m.Data, &msg); err != nil {
			return
		}
		f(msg)
	})

	return
}

func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	m := CreatedFeedMessage{}
	ch := make(chan *nats.Msg, 64)
	var err error

	n.feedCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				if err := n.DecodeMessage(msg.Data, &m); err != nil {
					continue
				}
				n.feedCreatedChan <- m
			case <-ctx.Done():
				return
			}
		}
	}()

	return n.feedCreatedChan, nil
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	if err := gob.NewEncoder(&b).Encode(m); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (n *NatsEventStore) DecodeMessage(data []byte, m any) error {
	b := bytes.NewBuffer(data)
	if err := gob.NewDecoder(b).Decode(m); err != nil {
		return err
	}

	return nil
}
