package hostdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	insertContract      = "INSERT INTO contract_host (address, creator, transparent, custom_config, batch_seq, height, time) VALUES (?, ?, ?, ?, ?, ?, ?)"
	selectContract      = "SELECT id, address, creator, transparent, custom_config, batch_seq, height, time FROM contract_host"
	selectContractCount = "SELECT COUNT(*) FROM contract_host"
	updateContractCount = "UPDATE contract_count SET total=? WHERE id=1"
	whereContractAddr   = " WHERE address = ?"
	orderByDeployed     = " ORDER BY time DESC"
)

// AddContracts adds multiple contracts to the host DB in a single transaction
func AddContracts(db HostDB, contracts []common.PublicContract) error {
	if len(contracts) == 0 {
		return nil
	}

	dbtx, err := db.NewDBTransaction()
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer dbtx.Rollback()

	insertedCount := 0
	for _, contract := range contracts {
		reboundInsert := db.GetSQLDB().Rebind(insertContract)
		_, err := dbtx.Tx.Exec(reboundInsert,
			contract.Address.Bytes(),
			contract.Creator.Bytes(),
			contract.IsTransparent,
			contract.HasCustomConfig,
			contract.BatchSeq,
			contract.Height,
			contract.Time,
		)
		if err != nil {
			if IsRowExistsError(err) {
				continue
			}
			return fmt.Errorf("failed to insert contract %s: %w", contract.Address.Hex(), err)
		}
		insertedCount++
	}

	// update contract count if we inserted any new contracts
	if insertedCount > 0 {
		var currentTotal int
		reboundSelectCount := db.GetSQLDB().Rebind("SELECT total FROM contract_count WHERE id = 1")
		err = dbtx.Tx.QueryRow(reboundSelectCount).Scan(&currentTotal)
		if err != nil {
			return fmt.Errorf("failed to query contract count: %w", err)
		}

		newTotal := currentTotal + insertedCount
		reboundUpdateCount := db.GetSQLDB().Rebind(updateContractCount)
		_, err = dbtx.Tx.Exec(reboundUpdateCount, newTotal)
		if err != nil {
			return fmt.Errorf("failed to update contract count: %w", err)
		}
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetContractListing returns contracts with pagination, ordered by deployment time
func GetContractListing(db HostDB, pagination *common.QueryPagination) (*common.ContractListingResponse, error) {
	query := selectContract + orderByDeployed + paginationQuery
	reboundQuery := db.GetSQLDB().Rebind(query)

	rows, err := db.GetSQLDB().Query(reboundQuery, int64(pagination.Size), int64(pagination.Offset))
	if err != nil {
		return nil, fmt.Errorf("failed to execute contract listing query: %w", err)
	}
	defer rows.Close()

	var contracts []common.PublicContract
	for rows.Next() {
		var contract common.PublicContract
		var addressBytes, creatorBytes []byte

		err = rows.Scan(
			&contract.ID,
			&addressBytes,
			&creatorBytes,
			&contract.IsTransparent,
			&contract.HasCustomConfig,
			&contract.BatchSeq,
			&contract.Height,
			&contract.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract row: %w", err)
		}

		contract.Address = gethcommon.BytesToAddress(addressBytes)
		contract.Creator = gethcommon.BytesToAddress(creatorBytes)

		contracts = append(contracts, common.PublicContract{
			Address:         contract.Address,
			Creator:         contract.Creator,
			IsTransparent:   contract.IsTransparent,
			HasCustomConfig: contract.HasCustomConfig,
			BatchSeq:        contract.BatchSeq,
			Height:          contract.Height,
			Time:            contract.Time,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total uint64
	reboundCountQuery := db.GetSQLDB().Rebind(selectContractCount)
	err = db.GetSQLDB().QueryRow(reboundCountQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total contract count: %w", err)
	}

	return &common.ContractListingResponse{
		Contracts: contracts,
		Total:     total,
	}, nil
}

// GetContractByAddress returns a specific contract by address
func GetContractByAddress(db HostDB, address gethcommon.Address) (*common.PublicContract, error) {
	query := selectContract + whereContractAddr
	reboundQuery := db.GetSQLDB().Rebind(query)

	var contract common.PublicContract
	var addressBytes, creatorBytes []byte

	err := db.GetSQLDB().QueryRow(reboundQuery, address.Bytes()).Scan(
		&contract.ID,
		&addressBytes,
		&creatorBytes,
		&contract.IsTransparent,
		&contract.HasCustomConfig,
		&contract.BatchSeq,
		&contract.Height,
		&contract.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch contract: %w", err)
	}

	contract.Address = gethcommon.BytesToAddress(addressBytes)
	contract.Creator = gethcommon.BytesToAddress(creatorBytes)

	return &common.PublicContract{
		Address:         contract.Address,
		Creator:         contract.Creator,
		IsTransparent:   contract.IsTransparent,
		HasCustomConfig: contract.HasCustomConfig,
		BatchSeq:        contract.BatchSeq,
		Height:          contract.Height,
		Time:            contract.Time,
	}, nil
}

// GetTotalContractCount returns the total number of contracts
func GetTotalContractCount(db HostDB) (uint64, error) {
	var total uint64
	reboundQuery := db.GetSQLDB().Rebind("SELECT total FROM contract_count WHERE id = 1")
	err := db.GetSQLDB().QueryRow(reboundQuery).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get contract count: %w", err)
	}
	return total, nil
}
