package graphdb

import (
	"context"
	"go-graphdb/testenv"
	"testing"
)

func TestRDF4J_Protocol(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		protocol, err := client.RDF4J().Protocol(context.Background())
		if err != nil {
			t.Errorf("unable to get RDF4J protocol: %s", err)
		}

		if protocol != "12" {
			t.Errorf("unexpected protocol version '%s'", protocol)
		}
	})
}
