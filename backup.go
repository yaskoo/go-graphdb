package graphdb

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

const (
	PathBackup             = "/rest/recovery/backup"
	PathBackupCloud        = "/rest/recovery/cloud-backup"
	PathBackupRestore      = "/rest/recovery/restore"
	PathBackupRestoreCloud = "/rest/recovery/cloud-restore"
)

type BackupOptions struct {
	BucketUri        string   `json:"bucketUri,omitempty"`
	BackupSystemData bool     `json:"backupSystemData,omitempty"`
	Repositories     []string `json:"repositories,omitempty"`
}

type RestoreOptions struct {
	BucketUri               string   `json:"bucketUri,omitempty"`
	RestoreSystemData       bool     `json:"restoreSystemData,omitempty"`
	RemoveStaleRepositories bool     `json:"removeStaleRepositories,omitempty"`
	Repositories            []string `json:"repositories,omitempty"`
}

type BackupClient struct {
	client *Client
}

// Create request for a new backup to be created. If BackupOptions.BucketUri is not set, the backup will be provided,
// to the consumer function. If bucket URI is set, the consumer function is not called and can be nil.
func (b *BackupClient) Create(ctx context.Context, opts BackupOptions, consumer func(filename string, reader io.Reader) error) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	rh := CombinedResponseHandler(ExpectStatusCode(http.StatusOK), func(resp *http.Response) error {
		if opts.BucketUri != "" {
			return nil
		}

		header := resp.Header.Get("content-disposition")
		var filename string
		if header != "" {
			f, err := extractFilename(header)
			if err == nil {
				filename = f
			}
		}

		if filename == "" {
			filename = "backup.tar"
		}
		return consumer(filename, resp.Body)
	})

	r := bytes.NewReader(body)
	if opts.BucketUri != "" {
		return b.client.post(ctx, PathBackupCloud, nil, rh, Header("accept", "application/json"), MultipartFormData(Part{
			Key:   "params",
			Type:  "application/json",
			Value: r,
		}))
	}
	return b.client.post(ctx, PathBackup, r, rh, Header("accept", "application/json"), Header("content-type", "application/json"))
}

func (b *BackupClient) Restore(ctx context.Context, opts RestoreOptions, r io.Reader) error {
	data, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	parts := []Part{{Key: "params", Type: "application/json", Value: bytes.NewBuffer(data)}}
	if opts.BucketUri == "" {
		parts = append(parts, Part{Key: "file", Value: r})
		return b.client.post(ctx, PathBackupRestore, nil, ExpectStatusCode(http.StatusOK), MultipartFormData(parts...))
	}
	return b.client.post(ctx, PathBackupRestoreCloud, nil, ExpectStatusCode(http.StatusOK), MultipartFormData(parts...))
}

// extractFilename parses a Content-Disposition header and returns the filename.
// It supports both "filename" and "filename*" (RFC 5987).
func extractFilename(header string) (string, error) {
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return "", err
	}

	// Prefer RFC 5987 filename* if present
	if fnStar, ok := params["filename*"]; ok {
		// filename* has format: charset''url-encoded
		parts := strings.SplitN(fnStar, "''", 2)
		if len(parts) == 2 {
			// charset := parts[0] // usually UTF-8, can be ignored in most cases
			decoded, err := url.QueryUnescape(parts[1])
			if err != nil {
				return "", err
			}
			return decoded, nil
		}
		return fnStar, nil
	}

	// Fallback to plain filename
	if fn, ok := params["filename"]; ok {
		return fn, nil
	}

	return "", nil // no filename provided
}
