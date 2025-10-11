package graphdb

import (
	"context"
	"testing"

	. "go-graphdb/testenv"
)

func TestInfo_Version(t *testing.T) {
	WithEnv(t, func(url string) {
		client := New(url)
		version, err := client.Info().Version(context.Background())
		if err != nil {
			t.Errorf("failed to fetch version info: %v", err)
		}

		if version.Version != "10.8.9" {
			t.Errorf("incorrect version: %s", version.Version)
		}
	})
}
