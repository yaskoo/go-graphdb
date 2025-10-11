package graphdb

import (
	"context"
	"fmt"
	"net/http"
)

// todo: move to repositories

const (
	PathAccessControlLists = "/rest/repositories/%s/acl"
)

type Policy interface {
	PolicyType() string
}

type SystemPolicy struct {
	PolicyName string `json:"policy,omitempty"`
	Role       string `json:"role,omitempty"`
	Scope      string `json:"scope,omitempty"`
	Operation  string `json:"operation,omitempty"`
}

func (s SystemPolicy) PolicyType() string {
	return s.PolicyName
}

type StatementStatement struct {
	SystemPolicy
	Subject   string `json:"subject,omitempty"`
	Predicate string `json:"predicate,omitempty"`
	Object    string `json:"object,omitempty"`
	Context   string `json:"context,omitempty"`
}

func (s StatementStatement) PolicyType() string {
	return s.PolicyName
}

type PluginPolicy struct {
	SystemPolicy
	Plugin string `json:"plugin,omitempty"`
}

func (s PluginPolicy) PolicyType() string {
	return s.PolicyName
}

type ClearGraphEntry struct {
	SystemPolicy
	Context string `json:"context,omitempty"`
}

func (s ClearGraphEntry) PolicyType() string {
	return s.PolicyName
}

type AclClient struct {
	client *Client
}

func (a *AclClient) List(ctx context.Context, id string, conf ...RequestConfig) error {
	return a.client.get(ctx, fmt.Sprintf(PathAccessControlLists, id), func(resp *http.Response) error {
		return ErrNotStatus(http.StatusCreated, "acl", resp)
	}, conf...)
}

func (a *AclClient) Add(ctx context.Context, id string, policies []Policy, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(policies))
	return a.client.post(ctx, fmt.Sprintf(PathAccessControlLists, id), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusCreated, "acl", resp)
	}, conf...)
}

func (a *AclClient) Replace(ctx context.Context, id string, policies []Policy, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(policies))
	return a.client.put(ctx, fmt.Sprintf(PathAccessControlLists, id), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusOK, "acl", resp)
	}, conf...)
}

func (a *AclClient) Delete(ctx context.Context, id string, policies []Policy, conf ...RequestConfig) error {
	conf = append(conf, JsonBody(policies))
	return a.client.delete(ctx, fmt.Sprintf(PathAccessControlLists, id), nil, func(resp *http.Response) error {
		return ErrNotStatus(http.StatusNoContent, "acl", resp)
	}, conf...)
}
