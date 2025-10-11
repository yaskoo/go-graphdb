package testenv

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func WithEnv(t *testing.T, callback func(url string)) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "ontotext/graphdb:10.8.9",
		ExposedPorts: []string{"7200/tcp"},
		WaitingFor:   wait.ForListeningPort("7200/tcp"),
	}

	gc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer gc.Terminate(ctx)

	host, _ := gc.Host(ctx)
	port, _ := gc.MappedPort(ctx, "7200")

	url := fmt.Sprintf("http://%s:%s", host, port.Port())
	callback(url)
}
