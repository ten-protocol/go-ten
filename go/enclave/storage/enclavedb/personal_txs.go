package enclavedb

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

func CountTransactionsPerAddress(ctx context.Context, stmtCache *PreparedStatementCache, address *uint64, _ bool, _ bool) (uint64, error) {
	var count uint64

	unionQuery, params := createVisibleReceiptsQuery(*address)

	query := "SELECT count(DISTINCT u.id)  FROM (" + unionQuery + ") AS u "

	stmt, err := stmtCache.GetOrPrepare(query)
	if err != nil {
		return 0, fmt.Errorf("could not prepare query: %w", err)
	}

	err = stmt.QueryRowContext(ctx, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetTransactionsPerAddress(ctx context.Context, stmtCache *PreparedStatementCache, address *uint64, pagination *common.QueryPagination, _ bool, _ bool) ([]*core.InternalReceipt, error) {
	receipts, err := loadPersonalTxs(ctx, stmtCache, address, pagination)
	if err != nil {
		return nil, err
	}

	// remove duplicates
	slices.SortFunc(receipts, func(a, b *core.InternalReceipt) int {
		if a.BlockNumber.Uint64() != b.BlockNumber.Uint64() {
			return int(a.BlockNumber.Uint64() - b.BlockNumber.Uint64())
		}
		if a.TransactionIndex != b.TransactionIndex {
			return int(a.TransactionIndex - b.TransactionIndex)
		}
		return 0
	})

	receipts = slices.CompactFunc(receipts, func(a, b *core.InternalReceipt) bool {
		return a.BlockNumber.Uint64() == b.BlockNumber.Uint64() && a.TransactionIndex == b.TransactionIndex
	})

	return receipts, nil
}

func loadPersonalTxs(ctx context.Context, stmtCache *PreparedStatementCache, requestingAccountId *uint64, pagination *common.QueryPagination) ([]*core.InternalReceipt, error) {
	if requestingAccountId == nil {
		return nil, fmt.Errorf("you have to specify requestingAccount")
	}
	var queryParams []any

	visibleReceiptsQuery, visibleReceiptParams := createVisibleReceiptsQuery(*requestingAccountId)

	// apply the pagination directly on the receipts - before fetching all data
	innerQuery := "SELECT u.id  FROM (" + visibleReceiptsQuery + ") AS u ORDER BY u.id DESC LIMIT ? OFFSET ?"

	// fetch all receipt data only for the requested "Page"
	query := "select b.hash, b.height, curr_tx.hash, curr_tx.idx, rec.post_state, rec.status, rec.gas_used, rec.effective_gas_price, rec.created_contract_address, tx_sender.address, tx_contr.address, curr_tx.type "
	query += " from receipt rec " +
		"join (" + innerQuery + ") as inner_query on inner_query.id=rec.id " +
		"left join receipt_viewer rv on rec.id=rv.receipt " +
		"join batch b on rec.batch=b.sequence " +
		"join tx curr_tx on rec.tx=curr_tx.id " +
		"   join externally_owned_account tx_sender on curr_tx.sender_address=tx_sender.id " +
		"   left join contract tx_contr on curr_tx.contract=tx_contr.id "

	queryParams = append(queryParams, visibleReceiptParams...)
	queryParams = append(queryParams, pagination.Size, pagination.Offset)

	stmt, err := stmtCache.GetOrPrepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := make([]*core.InternalReceipt, 0)

	empty := true
	for rows.Next() {
		empty = false
		r, err := onRowWithReceipt(rows)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if empty {
		return nil, errutil.ErrNotFound
	}
	return receipts, nil
}

func onRowWithReceipt(rows *sql.Rows) (*core.InternalReceipt, error) {
	r := core.InternalReceipt{}

	var txIndex *uint
	var blockHash, transactionHash *gethcommon.Hash
	var blockNumber *uint64
	res := []any{&blockHash, &blockNumber, &transactionHash, &txIndex, &r.PostState, &r.Status, &r.GasUsed, &r.EffectiveGasPrice, &r.CreatedContract, &r.From, &r.To, &r.TxType}

	err := rows.Scan(res...)
	if err != nil {
		return nil, fmt.Errorf("could not load receipt from db: %w", err)
	}

	r.BlockHash = *blockHash
	r.BlockNumber = big.NewInt(int64(*blockNumber))
	r.TxHash = *transactionHash
	r.TransactionIndex = *txIndex

	r.Logs = make([]*types.Log, 0)
	return &r, nil
}

// Create a query that returns all receipt Ids that are visible to an address id
func createVisibleReceiptsQuery(address uint64) (string, []any) {
	// receipts visible to the sender
	senderQuery := "select rec.id from tx curr_tx join receipt rec on rec.tx=curr_tx.id WHERE curr_tx.sender_address=? AND curr_tx.is_synthetic=?"

	// receipts visible to the receiver - only applies to native transfers
	receiverQuery := "select rec.id from tx curr_tx join receipt rec on rec.tx=curr_tx.id WHERE curr_tx.to_eoa=? AND curr_tx.is_synthetic=?"

	// receipts visible because they contain an event that is visible
	eventsQuery := "select rec.id from tx curr_tx join receipt rec on rec.tx=curr_tx.id join receipt_viewer rv on rec.id=rv.receipt WHERE rv.eoa=? AND curr_tx.is_synthetic=?"

	unionQuery := senderQuery + " UNION " + receiverQuery + " UNION " + eventsQuery

	var params []any
	params = append(params, address, false, address, false, address, false)
	return unionQuery, params
}
