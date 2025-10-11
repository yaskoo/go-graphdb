package graphdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PathClusterConfig        = "/rest/cluster/config"
	PathClusterConfigNodes   = "/rest/cluster/config/node"
	PathClusterNodeStatus    = "/rest/cluster/node/status"
	PathClusterGroupStatus   = "/rest/cluster/group/status"
	PathClusterTruncateLog   = "/rest/cluster/truncate-log"
	PathClusterTag           = "/rest/cluster/config/tag"
	PathClusterSecondaryMode = "/rest/cluster/config/secondary-mode"
)

const (
	ClusterAddNodes         = "cluster_add_node"
	ClusterDeleteNodes      = "cluster_delete_node"
	ClusterReplaceNodes     = "cluster_replace_node"
	ClusterCreateConfig     = "cluster_create_config"
	ClusterUpdateConfig     = "cluster_update_config"
	ClusterDeleteConfig     = "cluster_delete_config"
	ClusterGetConfig        = "cluster_get_config"
	ClusterNodeStatus       = "cluster_node_status"
	ClusterGroupStatus      = "cluster_group_status"
	ClusterTruncateLog      = "cluster_truncate_log"
	ClusterAddTag           = "cluster_add_tag"
	ClusterRemoveTag        = "cluster_remove_tag"
	ClusterEnableSecondary  = "cluster_enable_secondary"
	ClusterDisableSecondary = "cluster_disable_secondary"
)

type ClusterError struct {
	Op       string
	NotFound bool
	Messages StringSliceMap
	Err      error
}

func (c ClusterError) Error() string {
	if c.Err != nil {
		return fmt.Sprintf("%s: %s", c.Op, c.Err.Error())
	}
	return c.Op
}

func (c ClusterError) Unwrap() error {
	return c.Err
}

type ClusterProperties struct {
	ElectionMinTimeout          int     `json:"electionMinTimeout,omitempty"`
	ElectionRangeTimeout        int     `json:"electionRangeTimeout,omitempty"`
	HeartbeatInterval           int     `json:"heartbeatInterval,omitempty"`
	MessageSizeKB               int     `json:"messageSizeKB,omitempty"`
	VerificationTimeout         int     `json:"verificationTimeout,omitempty"`
	TransactionLogMaximumSizeGB float32 `json:"transactionLogMaximumSizeGB,omitempty"`
	BatchUpdateInterval         int     `json:"batchUpdateInterval,omitempty"`
}

type RecoveryStatus struct {
	State         RecoveryState `json:"state,omitempty"`
	AffectedNodes []string      `json:"affectedNodes,omitempty"`
	Message       string        `json:"message,omitempty"`
}

type RecoveryState struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type TopologyStatus struct {
	State         string         `json:"state,omitempty"`
	PrimaryTags   map[string]int `json:"primaryTags,omitempty"`
	PrimaryIndex  int            `json:"primaryIndex,omitempty"`
	PrimaryLeader string         `json:"primaryLeader,omitempty"`
}

type NodeStatus struct {
	ClusterEnabled bool              `json:"clusterEnabled,omitempty"`
	Address        string            `json:"address,omitempty"`
	NodeState      string            `json:"nodeState,omitempty"`
	Term           int               `json:"term,omitempty"`
	SyncStatus     map[string]string `json:"syncStatus,omitempty"`
	LastLogTerm    int               `json:"lastLogTerm,omitempty"`
	LastLogIndex   int               `json:"lastLogIndex,omitempty"`
	Endpoint       string            `json:"endpoint,omitempty"`
	RecoveryStatus RecoveryStatus    `json:"recoveryStatus,omitempty"`
	TopologyStatus TopologyStatus    `json:"topologyStatus,omitempty"`
}

type ClusterConfig struct {
	ClusterProperties
	Nodes []string `json:"nodes,omitempty"`
}

type StringSliceMap struct {
	l StringSlice
	m map[string]string
}

func (s *StringSliceMap) IsList() bool {
	return len(s.l) > 0
}

func (s *StringSliceMap) List() []string {
	return s.l
}

func (s *StringSliceMap) Map() map[string]string {
	return s.m
}

func (s *StringSliceMap) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.l); err == nil {
		return nil
	}

	if err := json.Unmarshal(data, &s.m); err != nil {
		return err
	}
	return nil
}

// StringSlice is a slice of strings, but can unmarshal from a single string json literal.
type StringSlice []string

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringSlice{str}
		return nil
	}

	var all []string
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}

	*s = all
	return nil
}

type ClusterClient struct {
	client *Client
}

func (c *ClusterClient) Config(ctx context.Context, config ...RequestConfig) (ClusterConfig, error) {
	var cc ClusterConfig
	rh := clusterResponseHandler(ClusterGetConfig, http.StatusOK, &cc)
	return cc, c.client.get(ctx, PathClusterConfig, rh, config...)
}

func (c *ClusterClient) Create(ctx context.Context, cc ClusterConfig, config ...RequestConfig) (map[string]string, error) {
	config = append(config, JsonBody(cc))

	var messages map[string]string
	rh := clusterResponseHandler(ClusterCreateConfig, http.StatusCreated, &messages)
	return messages, c.client.post(ctx, PathClusterConfig, nil, rh, config...)
}

func (c *ClusterClient) Update(ctx context.Context, props ClusterProperties, config ...RequestConfig) (ClusterConfig, error) {
	config = append(config, JsonBody(props))
	var clusterConfig ClusterConfig
	rh := clusterResponseHandler(ClusterUpdateConfig, http.StatusOK, &clusterConfig)
	return clusterConfig, c.client.patch(ctx, PathClusterConfig, nil, rh, config...)
}

func (c *ClusterClient) Delete(ctx context.Context, config ...RequestConfig) error {
	rh := clusterResponseHandler(ClusterDeleteConfig, http.StatusOK, nil)
	return c.client.delete(ctx, PathClusterConfig, nil, rh, config...)
}

func (c *ClusterClient) AddNodes(ctx context.Context, nodes []string, config ...RequestConfig) (map[string]string, error) {
	config = append(config, JsonBody(map[string][]string{
		"nodes": nodes,
	}))

	var messages map[string]string
	rh := clusterResponseHandler(ClusterAddNodes, http.StatusOK, &messages)
	return messages, c.client.post(ctx, PathClusterConfigNodes, nil, rh, config...)
}

func (c *ClusterClient) DeleteNodes(ctx context.Context, nodes []string, config ...RequestConfig) (map[string]string, error) {
	config = append(config, JsonBody(map[string][]string{
		"nodes": nodes,
	}))

	var messages map[string]string
	rh := clusterResponseHandler(ClusterDeleteNodes, http.StatusOK, &messages)
	return messages, c.client.delete(ctx, PathClusterConfigNodes, nil, rh, config...)
}

func (c *ClusterClient) ReplaceNodes(ctx context.Context, add []string, remove []string, config ...RequestConfig) (map[string]string, error) {
	config = append(config, JsonBody(map[string][]string{
		"addNodes":    add,
		"removeNodes": remove,
	}))

	var messages map[string]string
	rh := clusterResponseHandler(ClusterReplaceNodes, http.StatusOK, &messages)
	return messages, c.client.patch(ctx, PathClusterConfigNodes, nil, rh, config...)
}

func (c *ClusterClient) NodeStatus(ctx context.Context, config ...RequestConfig) (NodeStatus, error) {
	var status NodeStatus
	rh := clusterResponseHandler(ClusterNodeStatus, http.StatusOK, &status)
	return status, c.client.get(ctx, PathClusterNodeStatus, rh, config...)
}

func (c *ClusterClient) Status(ctx context.Context, config ...RequestConfig) ([]NodeStatus, error) {
	var status []NodeStatus
	rh := clusterResponseHandler(ClusterGroupStatus, http.StatusOK, &status)
	return status, c.client.get(ctx, PathClusterGroupStatus, rh, config...)
}

func (c *ClusterClient) Truncate(ctx context.Context, config ...RequestConfig) error {
	rh := clusterResponseHandler(ClusterTruncateLog, http.StatusOK, nil)
	return c.client.post(ctx, PathClusterTruncateLog, nil, rh, config...)
}

func (c *ClusterClient) AddTag(ctx context.Context, tag string, config ...RequestConfig) error {
	config = append(config, JsonBody(map[string]string{
		"tag": tag,
	}))

	rh := clusterResponseHandler(ClusterAddTag, http.StatusOK, nil)
	return c.client.post(ctx, PathClusterTag, nil, rh, config...)
}

func (c *ClusterClient) RemoveTag(ctx context.Context, tag string, config ...RequestConfig) error {
	config = append(config, JsonBody(map[string]string{
		"tag": tag,
	}))

	rh := clusterResponseHandler(ClusterRemoveTag, http.StatusOK, nil)
	return c.client.delete(ctx, PathClusterTag, nil, rh, config...)
}

func (c *ClusterClient) EnableSecondaryMode(ctx context.Context, primary, tag string, config ...RequestConfig) error {
	config = append(config, JsonBody(map[string]string{
		"primaryNode": primary,
		"tag":         tag,
	}))

	rh := clusterResponseHandler(ClusterEnableSecondary, http.StatusOK, nil)
	return c.client.post(ctx, PathClusterSecondaryMode, nil, rh, config...)
}

func (c *ClusterClient) DisableSecondaryMode(ctx context.Context, config ...RequestConfig) error {
	rh := clusterResponseHandler(ClusterDisableSecondary, http.StatusOK, nil)
	return c.client.post(ctx, PathClusterSecondaryMode, nil, rh, config...)
}

func clusterResponseHandler(op string, okStatus int, re any) ResponseHandler {
	return func(resp *http.Response) error {
		if resp.StatusCode == okStatus {
			if re == nil {
				return nil
			}

			err := json.NewDecoder(resp.Body).Decode(&re)
			if err != nil {
				return ClusterError{
					Op:  op,
					Err: err,
				}
			}
			return nil
		}

		clusterErr := ClusterError{
			Op:       op,
			NotFound: resp.StatusCode == http.StatusNotFound,
		}
		err := json.NewDecoder(resp.Body).Decode(&clusterErr.Messages)
		if err != nil {
			clusterErr.Err = err
		}
		return clusterErr
	}
}
