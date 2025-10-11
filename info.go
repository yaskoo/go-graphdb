package graphdb

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	PathVersion = "/rest/info/version"
)

type VersionInfo struct {
	Product    string `json:"productType,omitempty"`
	Version    string `json:"productVersion,omitempty"`
	Connectors string `json:"connectors,omitempty"`
	Sesame     string `json:"sesame,omitempty"`
	Workbench  string `json:"Workbench,omitempty"`
}

type LicenseInfo struct {
	Version               string   `json:"version,omitempty"`
	Valid                 bool     `json:"valid,omitempty"`
	Message               string   `json:"message,omitempty"`
	InstallationId        string   `json:"installationId,omitempty"`
	Licensee              string   `json:"licensee,omitempty"`
	Product               string   `json:"product,omitempty"`
	TypeOfUse             string   `json:"typeOfUse,omitempty"`
	MaxCpuCores           int      `json:"maxCpuCores,omitempty"`
	ExpiryDate            string   `json:"expiryDate,omitempty"`
	LatestPublicationDate int64    `json:"latestPublicationDate,omitempty"`
	LicenseCapabilities   []string `json:"licenseCapabilities,omitempty"`
	ProductType           string   `json:"productType,omitempty"`
}

type InfoClient struct {
	client *Client
}

func (i *InfoClient) Version(ctx context.Context, conf ...RequestConfig) (VersionInfo, error) {
	var version VersionInfo
	return version, i.client.get(ctx, PathVersion, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "info", resp); err != nil {
			return err
		}

		err := json.NewDecoder(resp.Body).Decode(&version)
		return errors.Wrap(err, "failed to decode version info")
	}, conf...)
}
