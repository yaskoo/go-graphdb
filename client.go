package graphdb

import (
	"context"
	"io"
	"net/http"
)

type Option func(client *Client)

func WithTransport(transport http.RoundTripper) Option {
	return func(client *Client) {
		client.hc.Transport = transport
	}
}

func WithBasicAuth(username, password string) Option {
	return func(client *Client) {
		client.auth = func(req *http.Request) {
			req.SetBasicAuth(username, password)
		}
	}
}

type Client struct {
	url string
	hc  *http.Client

	auth RequestConfig

	repository   *RepositoryClient
	backup       *BackupClient
	cluster      *ClusterClient
	security     *SecurityClient
	report       *ReportClient
	locations    *LocationClient
	monitoring   *MonitoringClient
	savedQueries *SavedQueriesClient
	info         *InfoClient
	rdf4j        *RDF4J
}

func (c *Client) get(ctx context.Context, path string, rh ResponseHandler, conf ...RequestConfig) error {
	return c.do(ctx, http.MethodGet, path, nil, rh, conf...)
}

func (c *Client) post(ctx context.Context, path string, body io.Reader, rh ResponseHandler, conf ...RequestConfig) error {
	return c.do(ctx, http.MethodPost, path, body, rh, conf...)
}

func (c *Client) patch(ctx context.Context, path string, body io.Reader, rh ResponseHandler, conf ...RequestConfig) error {
	return c.do(ctx, http.MethodPatch, path, body, rh, conf...)
}

func (c *Client) put(ctx context.Context, path string, body io.Reader, rh ResponseHandler, conf ...RequestConfig) error {
	return c.do(ctx, http.MethodPut, path, body, rh, conf...)
}

func (c *Client) delete(ctx context.Context, path string, body io.Reader, rh ResponseHandler, conf ...RequestConfig) error {
	return c.do(ctx, http.MethodDelete, path, body, rh, conf...)
}

func (c *Client) do(ctx context.Context, m string, p string, b io.Reader, rh ResponseHandler, conf ...RequestConfig) error {
	req, err := http.NewRequestWithContext(ctx, m, c.url+p, b)
	if err != nil {
		return err
	}

	for _, cf := range conf {
		cf(req)
	}

	if c.auth != nil {
		c.auth(req)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	}

	err = rh(resp)
	_, _ = io.Copy(io.Discard, resp.Body)
	return err
}

func (c *Client) Repositories() *RepositoryClient {
	return c.repository
}

func (c *Client) Backups() *BackupClient {
	return c.backup
}

func (c *Client) Cluster() *ClusterClient {
	return c.cluster
}

func (c *Client) Security() *SecurityClient {
	return c.security
}

func (c *Client) Report() *ReportClient {
	return c.report
}

func (c *Client) Locations() *LocationClient {
	return c.locations
}

func (c *Client) Monitoring() *MonitoringClient {
	return c.monitoring
}

func (c *Client) SavedQueries() *SavedQueriesClient {
	return c.savedQueries
}

func (c *Client) Info() *InfoClient {
	return c.info
}

func (c *Client) RDF4J() *RDF4J {
	return c.rdf4j
}

func New(url string, opts ...Option) *Client {
	client := &Client{
		url: url,
		hc:  &http.Client{},
	}

	client.repository = &RepositoryClient{client: client}
	client.backup = &BackupClient{client: client}
	client.cluster = &ClusterClient{client: client}
	client.security = &SecurityClient{client: client}
	client.report = &ReportClient{client: client}
	client.locations = &LocationClient{client: client}
	client.monitoring = &MonitoringClient{client: client}
	client.savedQueries = &SavedQueriesClient{client: client}
	client.info = &InfoClient{client: client}
	client.rdf4j = &RDF4J{client: client}

	for _, opt := range opts {
		opt(client)
	}
	return client
}
