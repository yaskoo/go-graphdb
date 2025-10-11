package graphdb

import (
	"context"
	"go-graphdb/testenv"
	"io"
	"os"
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		ctx := context.Background()

		ok, err := client.Report().Generate(ctx)
		if !ok || err != nil {
			t.Error("error generating report", err)
		}

		var retries int
		var status ReportStatus
		for {
			status, err = client.Report().Status(ctx)
			if err != nil {
				t.Error("error getting report status", err)
				break
			}

			if status.State == StateReady {
				break
			}

			if retries >= 5 {
				t.Error("report retries exceeded waiting to be ready", retries)
				break
			}
			retries++
			time.Sleep(1 * time.Second)
		}

		temp, err := os.CreateTemp(os.TempDir(), "graphdb-")
		if err != nil {
			t.Fatal("error creating temp file for report", err)
		}
		defer temp.Close()

		err = client.Report().Download(ctx, func(filename string, r io.Reader) error {
			i, cpErr := io.Copy(temp, r)
			if cpErr != nil {
				t.Fatal("failed to copy report", cpErr)
			}

			if i == 0 {
				t.Fatal("zero bytes copied")
			}
			return nil
		})

		if err != nil {
			t.Fatal("error downloading report", err)
		}
	})
}
