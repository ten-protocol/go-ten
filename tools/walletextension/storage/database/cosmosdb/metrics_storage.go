package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	METRICS_CONTAINER_NAME = "metrics"
	METRICS_DOC_ID         = "global_metrics"
	SHARD_PREFIX           = "user_shard_"
	DEFAULT_SHARD_COUNT    = 50 // Number of shards to distribute users across
)

// GlobalMetricsDocument contains global metrics without user activity data
type GlobalMetricsDocument struct {
	ID                 string `json:"id"`
	TotalUsers         uint64 `json:"totalUsers"`
	AccountsRegistered uint64 `json:"accountsRegistered"`
	ActiveUsersCount   int    `json:"activeUsersCount"`
	LastUpdated        string `json:"lastUpdated"`
	ShardCount         int    `json:"shardCount"` // Number of shards being used
}

// UserShardDocument contains activity data for a subset of users
type UserShardDocument struct {
	ID          string            `json:"id"`
	ActiveUsers map[string]string `json:"activeUsers"` // double-hashed userID -> ISO timestamp
	ShardIndex  int               `json:"shardIndex"`
	LastUpdated string            `json:"lastUpdated"`
}

// MetricsStorageCosmosDB handles metrics persistence in CosmosDB
type MetricsStorageCosmosDB struct {
	client           *azcosmos.Client
	metricsContainer *azcosmos.ContainerClient
	shardCount       int
}

func NewMetricsStorage(connectionString string) (*MetricsStorageCosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	ctx := context.Background()

	// Ensure database exists
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	metricsContainer, err := client.NewContainer(DATABASE_NAME, METRICS_CONTAINER_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics container: %w", err)
	}

	return &MetricsStorageCosmosDB{
		client:           client,
		metricsContainer: metricsContainer,
		shardCount:       DEFAULT_SHARD_COUNT,
	}, nil
}

// getShardDocumentID determines which shard a user belongs to
func (m *MetricsStorageCosmosDB) getShardDocumentID(userID string) string {
	// Use hash function to determine shard
	h := fnv.New32a()
	h.Write([]byte(userID))
	shardIndex := int(h.Sum32()) % m.shardCount

	return fmt.Sprintf("%s%d", SHARD_PREFIX, shardIndex)
}

// getShardIndexFromID extracts shard index from document ID
func (m *MetricsStorageCosmosDB) getShardIndexFromID(shardID string) int {
	var index int
	fmt.Sscanf(shardID, SHARD_PREFIX+"%d", &index)
	return index
}

// LoadGlobalMetrics loads the global metrics document
func (m *MetricsStorageCosmosDB) LoadGlobalMetrics() (*GlobalMetricsDocument, error) {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(METRICS_DOC_ID)

	response, err := m.metricsContainer.ReadItem(ctx, partitionKey, METRICS_DOC_ID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// Initialize with empty global metrics if not found
			return &GlobalMetricsDocument{
				ID:         METRICS_DOC_ID,
				ShardCount: m.shardCount,
			}, nil
		}
		return nil, err
	}

	var doc GlobalMetricsDocument
	if err := json.Unmarshal(response.Value, &doc); err != nil {
		return nil, err
	}

	// Update shard count if it has changed
	if doc.ShardCount > 0 {
		m.shardCount = doc.ShardCount
	}

	return &doc, nil
}

// LoadUserShard loads a specific user shard document
func (m *MetricsStorageCosmosDB) LoadUserShard(shardID string) (*UserShardDocument, error) {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(shardID)

	response, err := m.metricsContainer.ReadItem(ctx, partitionKey, shardID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// Initialize empty shard if not found
			shardIndex := m.getShardIndexFromID(shardID)
			return &UserShardDocument{
				ID:          shardID,
				ActiveUsers: make(map[string]string),
				ShardIndex:  shardIndex,
			}, nil
		}
		return nil, err
	}

	var doc UserShardDocument
	if err := json.Unmarshal(response.Value, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// LoadAllShards loads all user shard documents
func (m *MetricsStorageCosmosDB) LoadAllShards() ([]*UserShardDocument, error) {
	var shards []*UserShardDocument

	// Load global metrics first to get current shard count
	global, err := m.LoadGlobalMetrics()
	if err != nil {
		return nil, err
	}

	if global.ShardCount > 0 {
		m.shardCount = global.ShardCount
	}

	// Load each shard
	for i := 0; i < m.shardCount; i++ {
		shardID := fmt.Sprintf("%s%d", SHARD_PREFIX, i)
		shard, err := m.LoadUserShard(shardID)
		if err != nil {
			return nil, err
		}
		shards = append(shards, shard)
	}

	return shards, nil
}

// SaveGlobalMetrics saves the global metrics document
func (m *MetricsStorageCosmosDB) SaveGlobalMetrics(metrics *GlobalMetricsDocument) error {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(METRICS_DOC_ID)

	metrics.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	metrics.ShardCount = m.shardCount

	docJSON, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = m.metricsContainer.UpsertItem(ctx, partitionKey, docJSON, nil)
	return err
}

// SaveUserShard saves a specific user shard document
func (m *MetricsStorageCosmosDB) SaveUserShard(shard *UserShardDocument) error {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(shard.ID)

	shard.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	docJSON, err := json.Marshal(shard)
	if err != nil {
		return err
	}

	_, err = m.metricsContainer.UpsertItem(ctx, partitionKey, docJSON, nil)
	return err
}

// GetUserActivityTimestamp retrieves the activity timestamp for a specific user
func (m *MetricsStorageCosmosDB) GetUserActivityTimestamp(userID string) (time.Time, bool, error) {
	// Determine which shard this user belongs to
	shardID := m.getShardDocumentID(userID)

	// Load the appropriate shard
	shard, err := m.LoadUserShard(shardID)
	if err != nil {
		return time.Time{}, false, err
	}

	// Find the user in the shard
	if timestampStr, exists := shard.ActiveUsers[userID]; exists {
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			return time.Time{}, true, err
		}
		return timestamp, true, nil
	}

	return time.Time{}, false, nil
}

// UpdateUserActivity updates activity for a specific user
func (m *MetricsStorageCosmosDB) UpdateUserActivity(userID string, timestamp time.Time) error {
	// Determine which shard this user belongs to
	shardID := m.getShardDocumentID(userID)

	// Load the appropriate shard
	shard, err := m.LoadUserShard(shardID)
	if err != nil {
		return err
	}

	// Update the user's activity timestamp
	shard.ActiveUsers[userID] = timestamp.UTC().Format(time.RFC3339)

	// Save the updated shard
	return m.SaveUserShard(shard)
}

// CountActiveUsers counts active users across all shards
func (m *MetricsStorageCosmosDB) CountActiveUsers(activeThreshold time.Time) (int, error) {
	shards, err := m.LoadAllShards()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, shard := range shards {
		for _, timestampStr := range shard.ActiveUsers {
			if timestamp, err := time.Parse(time.RFC3339, timestampStr); err == nil {
				if timestamp.After(activeThreshold) {
					count++
				}
			}
		}
	}

	return count, nil
}

// ResizeShards changes the number of shards and redistributes users
// Warning: This is a costly operation, should be used rarely
func (m *MetricsStorageCosmosDB) ResizeShards(newShardCount int) error {
	if newShardCount <= 0 {
		return fmt.Errorf("new shard count must be positive")
	}

	// Load all current shards to get all active users
	allShards, err := m.LoadAllShards()
	if err != nil {
		return err
	}

	// Collect all active users
	allActiveUsers := make(map[string]string)
	for _, shard := range allShards {
		for userID, timestamp := range shard.ActiveUsers {
			allActiveUsers[userID] = timestamp
		}
	}

	// Change shard count
	oldShardCount := m.shardCount
	m.shardCount = newShardCount

	// Update global metrics with new shard count
	global, err := m.LoadGlobalMetrics()
	if err != nil {
		m.shardCount = oldShardCount
		return err
	}

	global.ShardCount = newShardCount
	if err := m.SaveGlobalMetrics(global); err != nil {
		m.shardCount = oldShardCount
		return err
	}

	// Create new shards with redistributed users
	newShards := make(map[string]*UserShardDocument)

	// Redistribute users to new shards
	for userID, timestamp := range allActiveUsers {
		// Determine which new shard this user belongs to
		shardID := m.getShardDocumentID(userID)

		// Create shard document if it doesn't exist
		if _, exists := newShards[shardID]; !exists {
			shardIndex := m.getShardIndexFromID(shardID)
			newShards[shardID] = &UserShardDocument{
				ID:          shardID,
				ActiveUsers: make(map[string]string),
				ShardIndex:  shardIndex,
			}
		}

		// Add user to appropriate shard
		newShards[shardID].ActiveUsers[userID] = timestamp
	}

	// Save new shards
	for _, shard := range newShards {
		if err := m.SaveUserShard(shard); err != nil {
			// If there's an error, don't revert - just log and continue
			// because we're already committed to the new shard count
			return fmt.Errorf("error saving shard %s during resize: %w", shard.ID, err)
		}
	}

	return nil
}

// MetricsStorage interface defines the metrics storage operations
type MetricsStorage interface {
	LoadGlobalMetrics() (*GlobalMetricsDocument, error)
	SaveGlobalMetrics(*GlobalMetricsDocument) error
	LoadUserShard(string) (*UserShardDocument, error)
	SaveUserShard(*UserShardDocument) error
	LoadAllShards() ([]*UserShardDocument, error)
	UpdateUserActivity(string, time.Time) error
	GetUserActivityTimestamp(string) (time.Time, bool, error)
	CountActiveUsers(time.Time) (int, error)
	ResizeShards(int) error
}

// NoOpMetricsStorage is a no-op implementation of metrics storage
type noOpMetricsStorage struct{}

func NewNoOpMetricsStorage() MetricsStorage {
	return &noOpMetricsStorage{}
}

func (n *noOpMetricsStorage) LoadGlobalMetrics() (*GlobalMetricsDocument, error) {
	return &GlobalMetricsDocument{}, nil
}

func (n *noOpMetricsStorage) SaveGlobalMetrics(*GlobalMetricsDocument) error {
	return nil
}

func (n *noOpMetricsStorage) LoadUserShard(string) (*UserShardDocument, error) {
	return &UserShardDocument{ActiveUsers: make(map[string]string)}, nil
}

func (n *noOpMetricsStorage) SaveUserShard(*UserShardDocument) error {
	return nil
}

func (n *noOpMetricsStorage) LoadAllShards() ([]*UserShardDocument, error) {
	return []*UserShardDocument{}, nil
}

func (n *noOpMetricsStorage) UpdateUserActivity(string, time.Time) error {
	return nil
}

func (n *noOpMetricsStorage) GetUserActivityTimestamp(string) (time.Time, bool, error) {
	return time.Time{}, false, nil
}

func (n *noOpMetricsStorage) CountActiveUsers(time.Time) (int, error) {
	return 0, nil
}

func (n *noOpMetricsStorage) ResizeShards(int) error {
	return nil
}
