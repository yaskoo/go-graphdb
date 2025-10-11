package graphdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PathMonitorStructs    = "/rest/monitor/structures"
	PathMonitorInfra      = "/rest/monitor/infrastructure"
	PathMonitorCluster    = "/rest/monitor/cluster"
	PathMonitorRecovery   = "/rest/monitor/backup"
	PathMonitorRepository = "/rest/monitor/repository/%s"
)

type PageCacheStats struct {
	Hits   int `json:"cacheHit,omitempty"`
	Misses int `json:"cacheMiss,omitempty"`
}

type RepositoryStats struct {
	Queries            QueryStats `json:"queries,omitempty"`
	EntityPool         EpoolStats `json:"entityPool,omitempty"`
	ActiveTransactions int        `json:"activeTransactions,omitempty"`
	OpenConnections    int        `json:"openConnections,omitempty"`
}

type QueryStats struct {
	Slow       int `json:"slow,omitempty"`
	Suboptimal int `json:"suboptimal,omitempty"`
}

type EpoolStats struct {
	Reads  int `json:"epoolReads,omitempty"`
	Writes int `json:"epoolWrites,omitempty"`
	Size   int `json:"epoolSize,omitempty"`
}

type InfraStats struct {
	Heap                MemoryStats `json:"heapMemoryUsage,omitempty"`
	OffHeap             MemoryStats `json:"nonHeapMemoryUsage,omitempty"`
	Disk                DiskStats   `json:"storageMemory,omitempty"`
	ThreadCount         int         `json:"threadCount,omitempty"`
	CpuLoad             int         `json:"cpuLoad,omitempty"`
	ClassCount          int         `json:"classCount,omitempty"`
	GcCount             int         `json:"gcCount,omitempty"`
	OpenFileDescriptors int         `json:"openFileDescriptors,omitempty"`
	MaxFileDescriptors  int         `json:"maxFileDescriptors,omitempty"`
}

type MemoryStats struct {
	Max       int `json:"max,omitempty"`
	Committed int `json:"committed,omitempty"`
	Init      int `json:"init,omitempty"`
	Used      int `json:"used,omitempty"`
}

type DiskStats struct {
	DataDirUsed int `json:"dataDirUsed,omitempty"`
	WorkDirUsed int `json:"workDirUsed,omitempty"`
	LogsDirUsed int `json:"logsDirUsed,omitempty"`
	DataDirFree int `json:"dataDirFree,omitempty"`
	WorkDirFree int `json:"workDirFree,omitempty"`
	LogsDirFree int `json:"logsDirFree,omitempty"`
}

type ClusterStats struct {
	Term               int `json:"term,omitempty"`
	FailedRecoveries   int `json:"failureRecoveriesCount,omitempty"`
	FailedTransactions int `json:"failedTransactionsCount,omitempty"`
}

type NodeStats struct {
	Total        int `json:"nodesInCluster,omitempty"`
	InSync       int `json:"nodesInSync,omitempty"`
	OutOfSync    int `json:"nodesOutOfSync,omitempty"`
	Disconnected int `json:"nodesDisconnected,omitempty"`
	Syncing      int `json:"nodesSyncing,omitempty"`
}

type RecoveryStats struct {
	Id                   string              `json:"id,omitempty"`
	Username             string              `json:"username,omitempty"`
	Operation            string              `json:"operation,omitempty"`
	AffectedRepositories []string            `json:"affectedRepositories,omitempty"`
	RunningFor           int                 `json:"msSinceCreated,omitempty"`
	Node                 string              `json:"nodePerformingClusterBackup,omitempty"`
	Options              RecoveryConfigStats `json:"snapshotOptions,omitempty"`
}

type RecoveryConfigStats struct {
	WithRepositoryData bool     `json:"withRepositoryData,omitempty"`
	WithSystemData     bool     `json:"withSystemData,omitempty"`
	CleanDataDir       bool     `json:"cleanDataDir,omitempty"`
	Repositories       []string `json:"repositories,omitempty"`
}

type MonitoringClient struct {
	client *Client
}

func (m *MonitoringClient) PageCache(ctx context.Context, conf ...RequestConfig) (PageCacheStats, error) {
	var stats PageCacheStats
	return stats, m.client.get(ctx, PathMonitorStructs, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "stats", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		return nil
	}, conf...)
}

func (m *MonitoringClient) Repository(ctx context.Context, id string, conf ...RequestConfig) (RepositoryStats, error) {
	var stats RepositoryStats
	return stats, m.client.get(ctx, fmt.Sprintf(PathMonitorRepository, id), func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "stats", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		return nil
	}, conf...)
}

func (m *MonitoringClient) Infrastructure(ctx context.Context, conf ...RequestConfig) (InfraStats, error) {
	var stats InfraStats
	return stats, m.client.get(ctx, PathMonitorInfra, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "stats", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		return nil
	}, conf...)
}

func (m *MonitoringClient) Cluster(ctx context.Context, conf ...RequestConfig) (ClusterStats, error) {
	var stats ClusterStats
	return stats, m.client.get(ctx, PathMonitorCluster, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "stats", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		return nil
	}, conf...)
}

func (m *MonitoringClient) Recovery(ctx context.Context, conf ...RequestConfig) (RecoveryStats, error) {
	var stats RecoveryStats
	return stats, m.client.get(ctx, PathMonitorRecovery, func(resp *http.Response) error {
		if err := ErrNotStatus(http.StatusOK, "stats", resp); err != nil {
			return err
		}

		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		return nil
	}, conf...)
}
