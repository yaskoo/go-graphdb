package graphdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PathLocations = "/rest/locations"
)

type Location struct {
	Uri               string `json:"uri,omitempty"`
	Label             string `json:"label,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	AuthType          string `json:"authType,omitempty"`
	LocationType      string `json:"locationType,omitempty"`
	Active            bool   `json:"active,omitempty"`
	Local             bool   `json:"local,omitempty"`
	System            bool   `json:"system,omitempty"`
	ErrorMsg          string `json:"errorMsg,omitempty"`
	DefaultRepository string `json:"defaultRepository,omitempty"`
}

type LocationClient struct {
	client *Client
}

func (l *LocationClient) List(ctx context.Context, conf ...RequestConfig) ([]Location, error) {
	var locations []Location
	return locations, l.client.get(ctx, PathLocations, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "locations", resp); err != nil {
			return err
		}

		err := json.NewDecoder(resp.Body).Decode(&locations)
		if err != nil {
			return fmt.Errorf("locations: %w", err)
		}
		return nil
	}, conf...)
}

func (l *LocationClient) Add(ctx context.Context, location Location, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(location))
	return l.client.post(ctx, PathLocations, nil, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "locations", resp); err != nil {
			return err
		}
		return nil
	}, conf...)
}

func (l *LocationClient) Update(ctx context.Context, location Location, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(location))
	return l.client.put(ctx, PathLocations, nil, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "locations", resp); err != nil {
			return err
		}
		return nil
	}, conf...)
}

func (l *LocationClient) Delete(ctx context.Context, uri string, conf ...RequestConfig) error {
	conf = append(conf, Query("uri", uri))
	return l.client.delete(ctx, PathLocations, nil, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "locations", resp); err != nil {
			return err
		}
		return nil
	}, conf...)
}
