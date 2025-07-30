package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	METRICS_CONTAINER_NAME       = "metrics"
	METRICS_DOC_ID               = "global_metrics"
	SHARD_PREFIX                 = "user_shard_"
	ACTIVITY_STATS_DOC_ID        = "activity_stats"
	ACTIVITY_STATS_DOC_ID_PREFIX = "activity_stats_"
	DEFAULT_SHARD_COUNT          = 50 // Number of shards to distribute users across
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

// DailyStats represents a single day's statistics
type DailyStats struct {
	Date          string `json:"date"`          // ISO date format YYYY-MM-DD
	DailyActive   int    `json:"dailyActive"`   // Users active on this day
	WeeklyActive  int    `json:"weeklyActive"`  // Users active in the past 7 days from this date
	MonthlyActive int    `json:"monthlyActive"` // Users active in the past 30 days from this date
	NewUsers      int    `json:"newUsers"`      // New users that joined on this day
	TotalUsers    uint64 `json:"totalUsers"`    // Total users as of this date (for calculating new users)
}

// ActivityStatsDocument contains historical user activity statistics
type ActivityStatsDocument struct {
	ID          string       `json:"id"`
	LastUpdated string       `json:"lastUpdated"`
	DailyStats  []DailyStats `json:"dailyStats"`
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

// getCurrentActivityStatsDocumentID returns the document ID for the current year's activity stats
func getCurrentActivityStatsDocumentID() string {
	currentYear := time.Now().UTC().Year()
	return fmt.Sprintf("%s%d", ACTIVITY_STATS_DOC_ID_PREFIX, currentYear)
}

// getActivityStatsDocumentIDForYear returns the document ID for a specific year's activity stats
func getActivityStatsDocumentIDForYear(year int) string {
	return fmt.Sprintf("%s%d", ACTIVITY_STATS_DOC_ID_PREFIX, year)
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

// LoadActivityStats loads the activity statistics document for the current year
func (m *MetricsStorageCosmosDB) LoadActivityStats() (*ActivityStatsDocument, error) {
	ctx := context.Background()
	docID := getCurrentActivityStatsDocumentID()
	partitionKey := azcosmos.NewPartitionKeyString(docID)

	response, err := m.metricsContainer.ReadItem(ctx, partitionKey, docID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// Initialize with empty stats if not found
			return &ActivityStatsDocument{
				ID:         docID,
				DailyStats: []DailyStats{},
			}, nil
		}
		return nil, err
	}

	var doc ActivityStatsDocument
	if err := json.Unmarshal(response.Value, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// LoadActivityStatsByYear loads the activity statistics document for a specific year
func (m *MetricsStorageCosmosDB) LoadActivityStatsByYear(year int) (*ActivityStatsDocument, error) {
	ctx := context.Background()
	docID := getActivityStatsDocumentIDForYear(year)
	partitionKey := azcosmos.NewPartitionKeyString(docID)

	response, err := m.metricsContainer.ReadItem(ctx, partitionKey, docID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// Initialize with empty stats if not found
			return &ActivityStatsDocument{
				ID:         docID,
				DailyStats: []DailyStats{},
			}, nil
		}
		return nil, err
	}

	var doc ActivityStatsDocument
	if err := json.Unmarshal(response.Value, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// SaveActivityStats saves the activity statistics document
func (m *MetricsStorageCosmosDB) SaveActivityStats(stats *ActivityStatsDocument) error {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(stats.ID)

	stats.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	docJSON, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	_, err = m.metricsContainer.UpsertItem(ctx, partitionKey, docJSON, nil)
	return err
}

// UpdateDailyStats updates the activity stats document with current day's statistics
func (m *MetricsStorageCosmosDB) UpdateDailyStats() error {
	// Get current date in ISO format (YYYY-MM-DD)
	today := time.Now().UTC().Format("2006-01-02")

	// Load existing stats and global metrics
	stats, err := m.LoadActivityStats()
	if err != nil {
		return err
	}

	globalMetrics, err := m.LoadGlobalMetrics()
	if err != nil {
		return fmt.Errorf("failed to load global metrics: %w", err)
	}

	// Calculate time thresholds
	now := time.Now().UTC()
	dailyStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weeklyStart := dailyStart.AddDate(0, 0, -7)
	monthlyStart := dailyStart.AddDate(0, 0, -30)

	// Count daily active users
	dailyCount, err := m.CountActiveUsers(dailyStart)
	if err != nil {
		return fmt.Errorf("failed to count daily active users: %w", err)
	}

	// Count weekly active users
	weeklyCount, err := m.CountActiveUsers(weeklyStart)
	if err != nil {
		return fmt.Errorf("failed to count weekly active users: %w", err)
	}

	// Count monthly active users
	monthlyCount, err := m.CountActiveUsers(monthlyStart)
	if err != nil {
		return fmt.Errorf("failed to count monthly active users: %w", err)
	}

	// Calculate new users for today
	currentTotalUsers := globalMetrics.TotalUsers
	newUsersToday := 0

	// Find the most recent entry before today to calculate new users
	if len(stats.DailyStats) > 0 {
		// Sort entries by date to find the most recent one before today
		sortedStats := make([]DailyStats, len(stats.DailyStats))
		copy(sortedStats, stats.DailyStats)
		sort.Slice(sortedStats, func(i, j int) bool {
			return sortedStats[i].Date > sortedStats[j].Date
		})

		// Find the most recent entry before today
		var previousTotalUsers uint64 = 0
		for _, stat := range sortedStats {
			if stat.Date < today {
				previousTotalUsers = stat.TotalUsers
				break
			}
		}

		// Calculate new users as the difference
		if currentTotalUsers >= previousTotalUsers {
			newUsersToday = int(currentTotalUsers - previousTotalUsers)
		} else {
			// This shouldn't happen unless there's a data inconsistency
			newUsersToday = 0
		}
	} else {
		// First day ever - all current users are "new" today
		newUsersToday = int(currentTotalUsers)
	}

	// Create today's stats
	todayStats := DailyStats{
		Date:          today,
		DailyActive:   dailyCount,
		WeeklyActive:  weeklyCount,
		MonthlyActive: monthlyCount,
		NewUsers:      newUsersToday,
		TotalUsers:    currentTotalUsers,
	}

	// Check if we already have an entry for today
	found := false
	for i, stat := range stats.DailyStats {
		if stat.Date == today {
			// Update existing entry but preserve NewUsers if it was already calculated
			if stat.NewUsers == 0 {
				stats.DailyStats[i] = todayStats
			} else {
				// Preserve existing NewUsers count, update other fields
				stats.DailyStats[i].DailyActive = dailyCount
				stats.DailyStats[i].WeeklyActive = weeklyCount
				stats.DailyStats[i].MonthlyActive = monthlyCount
				stats.DailyStats[i].TotalUsers = currentTotalUsers
			}
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		stats.DailyStats = append(stats.DailyStats, todayStats)
	}

	// Sort entries by date (newest first)
	sort.Slice(stats.DailyStats, func(i, j int) bool {
		return stats.DailyStats[i].Date > stats.DailyStats[j].Date
	})

	// Save the updated stats (no trimming - keep all data forever in yearly documents)
	return m.SaveActivityStats(stats)
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
	LoadActivityStats() (*ActivityStatsDocument, error)
	LoadActivityStatsByYear(int) (*ActivityStatsDocument, error)
	SaveActivityStats(*ActivityStatsDocument) error
	UpdateDailyStats() error
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

func (n *noOpMetricsStorage) LoadActivityStats() (*ActivityStatsDocument, error) {
	return &ActivityStatsDocument{DailyStats: []DailyStats{}}, nil
}

func (n *noOpMetricsStorage) LoadActivityStatsByYear(int) (*ActivityStatsDocument, error) {
	return &ActivityStatsDocument{DailyStats: []DailyStats{}}, nil
}

func (n *noOpMetricsStorage) SaveActivityStats(*ActivityStatsDocument) error {
	return nil
}

func (n *noOpMetricsStorage) UpdateDailyStats() error {
	return nil
}
