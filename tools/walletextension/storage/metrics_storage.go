package storage

import (
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

// MetricsStorage defines the interface for metrics storage
type MetricsStorage interface {
	// StoreMetrics stores the daily metrics
	StoreMetrics(metrics DailyMetrics) error
	// GetMetricsForPeriod returns metrics for the specified time range
	GetMetricsForPeriod(from, to time.Time) ([]DailyMetrics, error)
}

type DailyMetrics struct {
	Date                    time.Time `json:"date"`
	RegisteredUsers         int       `json:"registered_users"`
	AuthenticatedUsers      int       `json:"authenticated_users"`
	MonthlyActiveUsers      int       `json:"monthly_active_users"`
	LastCalculatedUserCount int       `json:"last_calculated_user_count"`
}

// NewMetricsStorage creates a new metrics storage instance based on the database type
func NewMetricsStorage(dbType, dbConnectionURL string, logger gethlog.Logger) (MetricsStorage, error) {
	switch dbType {
	case "cosmosDB":
		return cosmosdb.NewMetricsStorageCosmosDB(dbConnectionURL)
	default:
		// For sqlite we return nil as we don't want to collect metrics
		return nil, nil
	}
}
