package graphdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PathSavedQueries = "/rest/sparql/saved-queries"
)

type SavedQueriesClient struct {
	client *Client
}

type SavedQuery struct {
	Name   string `json:"name,omitempty"`
	Body   string `json:"body,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

func (s *SavedQueriesClient) SavedQueries(ctx context.Context, conf ...RequestConfig) ([]SavedQuery, error) {
	var queries []SavedQuery
	return queries, s.client.get(ctx, PathSavedQueries, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "saved_queries", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&queries); err != nil {
			return fmt.Errorf("saved_queries: %w", err)
		}
		return nil
	}, conf...)
}

func (s *SavedQueriesClient) UpdateSavedQueries(ctx context.Context, name string, query SavedQuery, conf ...RequestConfig) error {
	conf = append(conf, Query("oldQueryName", name), JsonBody(query))
	return s.client.put(ctx, PathSavedQueries, nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusOK, "saved_queries", resp)
	}, conf...)
}

func (s *SavedQueriesClient) CreateSavedQueries(ctx context.Context, query SavedQuery, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(query))
	return s.client.post(ctx, PathSavedQueries, nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusCreated, "saved_queries", resp)
	}, conf...)
}

func (s *SavedQueriesClient) DeleteSavedQueries(ctx context.Context, name string, conf ...RequestConfig) error {
	conf = append(conf, Query("name", name))
	return s.client.put(ctx, PathSavedQueries, nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusOK, "saved_queries", resp)
	}, conf...)
}
