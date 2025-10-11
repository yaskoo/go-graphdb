package graphdb

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	PathProtocol     = "/protocol"
	PathTransactions = "/repositories/%s/transactions"
	PathTransaction  = "/repositories/%s/transactions/%s"
)

type RDF4J struct {
	client *Client
}

func (r *RDF4J) Protocol(ctx context.Context, config ...RequestConfig) (string, error) {
	var v string
	return v, r.client.get(ctx, PathProtocol, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return ErrNotStatus(http.StatusOK, "rdf4j", resp)
		}

		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("rdf4j: %w", err)
		}
		v = string(all)
		return nil
	}, config...)
}
