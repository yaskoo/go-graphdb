package graphdb

import (
	"context"
	"encoding/json"
	fmt "fmt"
	"io"
	"net/http"
	"strings"
)

const (
	PathRepositories                  = "/rest/repositories"
	PathRepository                    = PathRepositories + "/%s"
	PathRepositorySize                = PathRepository + "/size"
	PathRepositoryRestart             = PathRepository + "/restart"
	PathRepositoryImport              = PathRepository + "/import"
	PathRepositoryImportServer        = PathRepositoryImport + "/server"
	PathRepositorySparqlTemplates     = PathRepository + "/sparql-templates"
	PathRepositorySparqlTemplatesExec = PathRepository + "/sparql-templates/execute"
	PathRepositorySparqlTemplatesConf = PathRepository + "/sparql-templates/configuration"

	PathSqlViews      = "/rest/sql-views/tables"
	PathSqlViewsTable = "/rest/sql-views/tables/%s"
)

type ImportSettings struct {
	Name                   string         `json:"name"`
	Status                 string         `json:"status"`
	Message                string         `json:"message"`
	Context                string         `json:"context"`
	ReplaceGraphs          []string       `json:"replaceGraphs"`
	BaseURI                string         `json:"baseURI"`
	ForceSerial            bool           `json:"forceSerial"`
	Type                   string         `json:"type"`
	Format                 string         `json:"format"`
	Data                   string         `json:"data"`
	ParserSettings         ParserSettings `json:"parserSettings"`
	Size                   string         `json:"size"`
	LastModified           int            `json:"lastModified"`
	Imported               int            `json:"imported"`
	AddedStatements        int            `json:"addedStatements"`
	RemovedStatements      int            `json:"removedStatements"`
	NumReplacedGraphs      int            `json:"numReplacedGraphs"`
	FileSize               int            `json:"fileSize"`
	FileLastModified       string         `json:"fileLastModified"`
	AddedRemovedStatements int            `json:"addedRemovedStatements"`
}

type ParserSettings struct {
	PreserveBNodeIds          bool   `json:"preserveBNodeIds"`
	FailOnUnknownDataTypes    bool   `json:"failOnUnknownDataTypes"`
	VerifyDataTypeValues      bool   `json:"verifyDataTypeValues"`
	NormalizeDataTypeValues   bool   `json:"normalizeDataTypeValues"`
	FailOnUnknownLanguageTags bool   `json:"failOnUnknownLanguageTags"`
	VerifyLanguageTags        bool   `json:"verifyLanguageTags"`
	NormalizeLanguageTags     bool   `json:"normalizeLanguageTags"`
	StopOnError               bool   `json:"stopOnError"`
	ContextLink               string `json:"contextLink"`
}

type RepositoryInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Uri         string `json:"uri"`
	ExternalUrl string `json:"externalUrl"`
	Local       bool   `json:"local"`
	Type        string `json:"type"`
	SesameType  string `json:"sesameType"`
	Location    string `json:"location"`
	Readable    bool   `json:"readable"`
	Writable    bool   `json:"writable"`
	Unsupported bool   `json:"unsupported"`
	State       string `json:"state"`
}

type RepositorySize struct {
	Inferred int `json:"inferred"`
	Total    int `json:"total"`
	Explicit int `json:"explicit"`
}

type RepositoryConfig struct {
	Id         string                               `json:"id,omitempty"`
	Title      string                               `json:"title,omitempty"`
	Type       string                               `json:"type,omitempty"`
	SesameType string                               `json:"sesameType,omitempty"`
	Location   string                               `json:"location,omitempty"`
	Params     map[string]RepositoryConfigParameter `json:"params,omitempty"`
}

type RepositoryConfigParameter struct {
	Name  string     `json:"name,omitempty"`
	Label string     `json:"label,omitempty"`
	Value *IntString `json:"value,omitempty"`
}

type IntString struct {
	i *int
	s *string
}

func (is *IntString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		is.s = &s
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		is.i = &i
		return nil
	}
	return fmt.Errorf("IntString: not an int or string: %s", string(data))
}

func (is *IntString) MarshalJSON() ([]byte, error) {
	if is.s != nil {
		return json.Marshal(*is.s)
	}

	if is.i != nil {
		return json.Marshal(*is.i)
	}

	return []byte("null"), nil
}

func (is *IntString) IsZero() bool {
	return is.i == nil && is.s == nil
}

type RepositoryClient struct {
	client *Client
}

func (r *RepositoryClient) Infos(ctx context.Context, conf ...RequestConfig) ([]RepositoryInfo, error) {
	var infos []RepositoryInfo
	rh := CombinedResponseHandler(ExpectStatusCode(http.StatusOK), UnmarshalJson(&infos))
	return infos, r.client.get(ctx, PathRepositories, rh, conf...)
}

func (r *RepositoryClient) Size(ctx context.Context, id string, conf ...RequestConfig) (RepositorySize, error) {
	var size RepositorySize
	rh := CombinedResponseHandler(ExpectStatusCode(http.StatusOK), UnmarshalJson(&size))
	return size, r.client.get(ctx, fmt.Sprintf(PathRepositorySize, id), rh, conf...)
}

func (r *RepositoryClient) Create(ctx context.Context, config RequestConfig, other ...RequestConfig) error {
	rc := []RequestConfig{config}
	rc = append(rc, other...)
	return r.client.post(ctx, PathRepositories, nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusCreated {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("repo: %s", string(b))
		}
		return nil
	}, rc...)
}

func (r *RepositoryClient) Edit(ctx context.Context, id string, config RepositoryConfig, other ...RequestConfig) error {
	rc := []RequestConfig{JsonBody(config)}
	rc = append(rc, other...)
	return r.client.put(ctx, fmt.Sprintf(PathRepository, id), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusCreated {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("repo: %s", string(b))
		}
		return nil
	}, rc...)
}

func (r *RepositoryClient) Delete(ctx context.Context, id string, conf ...RequestConfig) error {
	rh := CombinedResponseHandler(ExpectStatusCode(http.StatusOK))
	return r.client.delete(ctx, fmt.Sprintf(PathRepository, id), nil, rh, conf...)
}

func (r *RepositoryClient) Restart(ctx context.Context, id string, conf ...RequestConfig) error {
	rh := CombinedResponseHandler(ExpectOneOfStatusCode(http.StatusOK, http.StatusAccepted))
	return r.client.post(ctx, fmt.Sprintf(PathRepositoryRestart, id), nil, rh, conf...)
}

func (r *RepositoryClient) ServerFiles(ctx context.Context, id string, conf ...RequestConfig) ([]ImportSettings, error) {
	var available []ImportSettings
	return available, r.client.get(ctx, fmt.Sprintf(PathRepositoryImportServer, id), func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("import server: %s", string(b))
		}

		err := json.NewDecoder(resp.Body).Decode(&available)
		if err != nil {
			return fmt.Errorf("import server: %w", err)
		}
		return nil
	}, conf...)
}

func (r *RepositoryClient) ImportServerFiles(ctx context.Context, id string, files []string, settings ImportSettings, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(map[string]interface{}{
		"fileNames":      files,
		"importSettings": settings,
	}))
	return r.client.post(ctx, fmt.Sprintf(PathRepositoryImportServer, id), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusAccepted {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("import server: %s", string(b))
		}
		return nil
	}, conf...)
}

func (r *RepositoryClient) CancelServerFile(ctx context.Context, id string, file string, conf ...RequestConfig) error {
	conf = append(conf, Query("name", file))
	return r.client.delete(ctx, fmt.Sprintf(PathRepositoryImportServer, id), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusAccepted {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("import server: %s", string(b))
		}
		return nil
	}, conf...)
}

// todo: imports from text, url, and upload are hidden... why?

func (r *RepositoryClient) SparqlTemplates(ctx context.Context, repo string, conf ...RequestConfig) ([]string, error) {
	var ids []string
	return ids, r.client.get(ctx, fmt.Sprintf(PathRepositorySparqlTemplates, repo), func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "sparql_templates", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
			return fmt.Errorf("sparql_templates: %w", err)
		}
		return nil
	}, conf...)
}

type SparqlTemplate struct {
	Id    string `json:"id,omitempty"`
	Query string `json:"query,omitempty"`
}

func (r *RepositoryClient) CreateSparqlTemplates(ctx context.Context, repo, template SparqlTemplate, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(template))
	return r.client.put(ctx, fmt.Sprintf(PathRepositorySparqlTemplates, repo), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusCreated, "sparql_templates", resp)
	}, conf...)
}

func (r *RepositoryClient) UpdateSparqlTemplates(ctx context.Context, repo string, template SparqlTemplate, conf ...RequestConfig) error {
	conf = append(conf, Query("templateID", template.Id))
	body := strings.NewReader(template.Query)

	return r.client.put(ctx, fmt.Sprintf(PathRepositorySparqlTemplates, repo), body, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusOK, "sparql_templates", resp)
	}, conf...)
}

func (r *RepositoryClient) DeleteSparqlTemplates(ctx context.Context, repo, template string, conf ...RequestConfig) error {
	conf = append(conf, Query("templateID", template))
	return r.client.delete(ctx, fmt.Sprintf(PathRepositorySparqlTemplates, repo), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusNoContent, "sparql_templates", resp)
	}, conf...)
}

func (r *RepositoryClient) RunSparqlTemplates(ctx context.Context, repo string, params map[string]interface{}, consumer func(r io.Reader) error, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(params))
	return r.client.post(ctx, fmt.Sprintf(PathRepositorySparqlTemplatesExec, repo), nil, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusNoContent, "sparql_templates", resp); err != nil {
			return err
		}
		return consumer(resp.Body)
	}, conf...)
}

func (r *RepositoryClient) SparqlTemplate(ctx context.Context, repo, template string, conf ...RequestConfig) (SparqlTemplate, error) {
	conf = append(conf, Query("templateID", template))

	var v SparqlTemplate
	return v, r.client.get(ctx, fmt.Sprintf(PathRepositorySparqlTemplatesConf, repo), func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusNoContent, "sparql_templates", resp); err != nil {
			return err
		}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return fmt.Errorf("sparql_templates: %w", err)
		}
		return nil
	}, conf...)
}

type SqlView struct {
	Name    string    `json:"name,omitempty"`
	Query   string    `json:"query,omitempty"`
	Columns SqlColumn `json:"columns,omitempty"`
}

type SqlColumn struct {
	Name             string `json:"column_name,omitempty"`
	Type             string `json:"column_type,omitempty"`
	SqlTypePrecision int    `json:"sql_type_precision,omitempty"`
	SqlTypeScale     int    `json:"sql_type_scale,omitempty"`
	Nullable         bool   `json:"nullable,omitempty"`
	SparqlType       string `json:"sparql_type,omitempty"`
}

func (r *RepositoryClient) SqlView(ctx context.Context, repo, view string, conf ...RequestConfig) (SqlView, error) {
	conf = append(conf, Header("x-graphdb-repository", repo))

	var v SqlView
	return v, r.client.get(ctx, fmt.Sprintf(PathSqlViewsTable, view), func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "sql_views", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return fmt.Errorf("sql_views: %w", err)
		}
		return nil
	}, conf...)
}

func (r *RepositoryClient) UpdateSqlView(ctx context.Context, repo string, view SqlView, conf ...RequestConfig) error {
	conf = append(conf, Header("x-graphdb-repository", repo), JsonBody(view))

	return r.client.put(ctx, fmt.Sprintf(PathSqlViewsTable, view.Name), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusOK, "sql_views", resp)
	}, conf...)
}

func (r *RepositoryClient) DeleteSqlView(ctx context.Context, repo string, view string, conf ...RequestConfig) error {
	conf = append(conf, Header("x-graphdb-repository", repo), JsonBody(view))

	return r.client.delete(ctx, fmt.Sprintf(PathSqlViewsTable, view), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusNoContent, "sql_views", resp)
	}, conf...)
}

func (r *RepositoryClient) SqlViews(ctx context.Context, repo string, conf ...RequestConfig) ([]string, error) {
	conf = append(conf, Header("x-graphdb-repository", repo))

	var ids []string
	return ids, r.client.get(ctx, PathSqlViews, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "sql_views", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
			return fmt.Errorf("sql_views: %w", err)
		}
		return nil
	}, conf...)
}

func (r *RepositoryClient) CreateSqlView(ctx context.Context, repo string, view SqlView, conf ...RequestConfig) error {
	conf = append(conf, Header("x-graphdb-repository", repo), JsonBody(view))

	return r.client.post(ctx, PathSqlViews, nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusCreated, "sql_views", resp)
	}, conf...)
}
