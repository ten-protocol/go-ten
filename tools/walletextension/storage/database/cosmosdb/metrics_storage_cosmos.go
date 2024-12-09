package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

const (
	METRICS_CONTAINER_NAME = "metrics"
)

// MetricsStorageCosmosDB implements MetricsStorage interface using CosmosDB
type MetricsStorageCosmosDB struct {
	client           *azcosmos.Client
	metricsContainer *azcosmos.ContainerClient
}

func NewMetricsStorageCosmosDB(connectionString string) (*MetricsStorageCosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	// Ensure database exists
	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Create container for metrics
	metricsContainer, err := client.NewContainer(DATABASE_NAME, METRICS_CONTAINER_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics container: %w", err)
	}

	return &MetricsStorageCosmosDB{
		client:           client,
		metricsContainer: metricsContainer,
	}, nil
}

func (m *MetricsStorageCosmosDB) StoreMetrics(metrics storage.DailyMetrics) error {
	ctx := context.Background()

	// Use date as document ID
	docID := metrics.Date.Format("2006-01-02")
	partitionKey := azcosmos.NewPartitionKeyString(docID)

	// Marshal metrics data to JSON
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	_, err = m.metricsContainer.UpsertItem(ctx, partitionKey, metricsJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert metrics: %w", err)
	}

	return nil
}

func (m *MetricsStorageCosmosDB) GetMetricsForPeriod(from, to time.Time) ([]storage.DailyMetrics, error) {
	ctx := context.Background()

	query := fmt.Sprintf("SELECT * FROM c WHERE c.id >= '%s' AND c.id <= '%s'",
		from.Format("2006-01-02"),
		to.Format("2006-01-02"))

	queryPager := m.metricsContainer.NewQueryItemsPager(query, azcosmos.NewPartitionKeyString(""), nil)

	var metrics []storage.DailyMetrics

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get metrics: %w", err)
		}

		for _, item := range queryResponse.Items {
			var metric storage.DailyMetrics
			if err := json.Unmarshal(item, &metric); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}
