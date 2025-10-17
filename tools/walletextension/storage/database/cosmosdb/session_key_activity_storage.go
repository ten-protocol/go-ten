package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	gethcommon "github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	sessionKeyActivityContainerName = "session_key_activities"
	sessionKeyActivityDocID         = "activities"
)

type sessionKeyActivityStorageCosmosDB struct {
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
}

type sessionKeyActivityDTO struct {
	ID    string                      `json:"id"`
	Items []sessionKeyActivityItemDTO `json:"items"`
}

type sessionKeyActivityItemDTO struct {
	Addr       string `json:"addr"`
	UserIDHex  string `json:"userId"`
	LastActive string `json:"lastActive"`
}

func NewSessionKeyActivityStorage(connectionString string) (*sessionKeyActivityStorageCosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	container, err := client.NewContainer(DATABASE_NAME, sessionKeyActivityContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session key activities container: %w", err)
	}

	return &sessionKeyActivityStorageCosmosDB{client: client, container: container}, nil
}

func (s *sessionKeyActivityStorageCosmosDB) Load() ([]wecommon.SessionKeyActivity, error) {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(sessionKeyActivityDocID)
	resp, err := s.container.ReadItem(ctx, pk, sessionKeyActivityDocID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		return nil, err
	}

	var dto sessionKeyActivityDTO
	if err := json.Unmarshal(resp.Value, &dto); err != nil {
		return nil, err
	}

	result := make([]wecommon.SessionKeyActivity, 0, len(dto.Items))
	for _, it := range dto.Items {
		addr := gethcommon.HexToAddress(it.Addr)
		userID, _ := hex.DecodeString(strings.TrimPrefix(it.UserIDHex, "0x"))
		t, err := time.Parse(time.RFC3339, it.LastActive)
		if err != nil {
			continue
		}
		result = append(result, wecommon.SessionKeyActivity{Addr: addr, UserID: userID, LastActive: t})
	}
	return result, nil
}

func (s *sessionKeyActivityStorageCosmosDB) Save(items []wecommon.SessionKeyActivity) error {
	ctx := context.Background()
	pk := azcosmos.NewPartitionKeyString(sessionKeyActivityDocID)

	dto := sessionKeyActivityDTO{ID: sessionKeyActivityDocID, Items: make([]sessionKeyActivityItemDTO, 0, len(items))}
	for _, it := range items {
		dto.Items = append(dto.Items, sessionKeyActivityItemDTO{
			Addr:       it.Addr.Hex(),
			UserIDHex:  hex.EncodeToString(it.UserID),
			LastActive: it.LastActive.UTC().Format(time.RFC3339),
		})
	}

	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	_, err = s.container.UpsertItem(ctx, pk, b, nil)
	return err
}
