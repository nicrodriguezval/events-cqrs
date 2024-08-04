package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/nicrodriguezval/cqrs/models"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client: client}, nil
}

func (r *ElasticSearchRepository) Close() {
	// Close the client
}

func (r *ElasticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := r.client.Index(
		"feeds",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(feed.ID),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)

	return err
}

func (r *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) (results []models.Feed, err error) {
	var buf bytes.Buffer

	searchQuery := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,      // Allow up to 3 characters to be different
				"cutoff_frequency": 0.0001, // Ignore terms that appear in more than 0.01% of documents
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("feeds"),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes map[string]any
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	for _, hit := range eRes["hits"].(map[string]any)["hits"].([]any) {
		source := hit.(map[string]any)["_source"].(map[string]any)
		sourceBuf, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var feed models.Feed
		if err := json.Unmarshal(sourceBuf, &feed); err != nil {
			return nil, err
		}

		results = append(results, feed)
	}

  return results, nil
}
