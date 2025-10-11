package graphdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-graphdb/testenv"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRepository_List(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		_, err := createRepository(t, client)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}

		repos, err := client.Repositories().Infos(context.Background())
		if err != nil {
			t.Fatalf("failed to get repositories: %v", err)
		}

		if len(repos) == 0 {
			t.Fatalf("failed to find any repositories")
		}
	})
}

func TestRepository_CreateEditDelete(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		id, err := createRepositoryWithTTL(t, client)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}

		config, err := repositoryJsonConfig()
		if err != nil {
			t.Fatal("cannot read defaul repository config")
		}
		config.Id = uuid.New().String()

		err = client.Repositories().Edit(context.Background(), id, config)
		if err != nil {
			fmt.Println(err)
		}

		if err = client.Repositories().Delete(context.Background(), id); err != nil {
			t.Fatalf("failed to delete repository: %v", err)
		}
	})
}

func TestRepository_Size(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		id, err := createRepository(t, client)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}

		size, err := client.Repositories().Size(context.Background(), id)
		if err != nil {
			t.Fatalf("failed to delete repository: %v", err)
		}

		if size.Total != 70 || size.Inferred != 70 || size.Explicit != 0 {
			t.Errorf("unexpected repository size: %v", size)
		}
	})
}

func TestRepository_Restart(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		id, err := createRepository(t, client)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}

		err = client.Repositories().Restart(context.Background(), id, Query("sync", "true"))
		if err != nil {
			t.Fatalf("failed to restart repository with sync query param: %v", err)
		}

		err = client.Repositories().Restart(context.Background(), id)
		if err != nil {
			t.Fatalf("failed to restart repository: %v", err)
		}
	})
}

func createRepository(t *testing.T, client *Client) (string, error) {
	config, err := repositoryJsonConfig()
	if err != nil {
		return "", err
	}
	config.Id = uuid.New().String()
	return config.Id, client.Repositories().Create(context.Background(), JsonBody(config))
}

func createRepositoryWithTTL(t *testing.T, client *Client) (string, error) {
	file, err := os.Open("testdata/repository-config.ttl")
	if err != nil {
		t.Fatalf("failed to open repository config file: %v", err)
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read repository config file: %v", err)
	}

	id := uuid.New().String()
	config := strings.ReplaceAll(string(all), "[[REPOSITORY_ID]]", id)
	part := Part{
		Key:      "config",
		Filename: "config.ttl",
		Value:    bytes.NewReader([]byte(config)),
	}
	return id, client.Repositories().Create(context.Background(), MultipartFormData(part))
}

func repositoryJsonConfig() (RepositoryConfig, error) {
	var v RepositoryConfig
	file, err := os.Open("testdata/repository-config.json")
	if err != nil {
		return v, err
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return v, err
	}

	return v, json.Unmarshal(all, &v)
}
