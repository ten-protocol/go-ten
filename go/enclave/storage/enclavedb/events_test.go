package enclavedb

import (
	"context"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

func TestGetTransactionsPerAddress(t *testing.T) {
	db, err := createTestSQLiteDB()
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}
	defer db.Close()

	setupTestData(t, db)

	testAddress := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("should return both sent and received transactions", func(t *testing.T) {
		pagination := &common.QueryPagination{Size: 10, Offset: 0}

		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, receipts, 2)

		// verify both tx are presnet
		txHashes := make(map[gethcommon.Hash]bool)
		for _, receipt := range receipts {
			txHashes[receipt.TxHash] = true
		}

		expectedTx1 := gethcommon.HexToHash("0x4444444444444444444444444444444444444444444444444444444444444444")
		expectedTx2 := gethcommon.HexToHash("0x5555555555555555555555555555555555555555555555555555555555555555")

		require.True(t, txHashes[expectedTx1], "Should include transaction where testAddress is sender")
		require.True(t, txHashes[expectedTx2], "Should include transaction where testAddress is recipient")
	})

	t.Run("should not return transactions where address is neither sender nor recipient", func(t *testing.T) {
		pagination := &common.QueryPagination{Size: 10, Offset: 0}

		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)

		unexpectedTx := gethcommon.HexToHash("0x6666666666666666666666666666666666666666666666666666666666666666")
		for _, receipt := range receipts {
			require.NotEqual(t, unexpectedTx, receipt.TxHash, "Should not include transaction where address is neither sender nor recipient")
		}
	})

	t.Run("should handle pagination correctly", func(t *testing.T) {
		pagination := &common.QueryPagination{Size: 1, Offset: 0}
		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		pagination = &common.QueryPagination{Size: 1, Offset: 1}
		receipts, err = GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		// verify different results
		pagination = &common.QueryPagination{Size: 10, Offset: 0}
		allReceipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, allReceipts, 2)
		require.NotEqual(t, allReceipts[0].TxHash, allReceipts[1].TxHash)
	})

	t.Run("should return empty result for non-existent address", func(t *testing.T) {
		nonExistentAddress := gethcommon.HexToAddress("0x9999999999999999999999999999999999999999")
		pagination := &common.QueryPagination{Size: 10, Offset: 0}

		receipts, err := GetTransactionsPerAddress(context.Background(), db, &nonExistentAddress, pagination)
		require.Error(t, err)
		require.ErrorIs(t, err, errutil.ErrNotFound)
		require.Nil(t, receipts)
	})

	t.Run("should handle contract transactions", func(t *testing.T) {
		contractAddress := gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccc")
		_, err := db.Exec("INSERT INTO contract (id, address, creator, auto_visibility, transparent, tx) VALUES (1, ?, 1, true, false, 1)",
			contractAddress.Bytes())
		require.NoError(t, err)

		pagination := &common.QueryPagination{Size: 10, Offset: 0}
		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, receipts, 2) // should still return 2 transactions
	})

	t.Run("should verify receipt data integrity", func(t *testing.T) {
		pagination := &common.QueryPagination{Size: 10, Offset: 0}
		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)

		for _, receipt := range receipts {
			require.NotNil(t, receipt)
			require.NotEqual(t, gethcommon.Hash{}, receipt.TxHash)
			require.NotEqual(t, gethcommon.Hash{}, receipt.BlockHash)
			require.NotNil(t, receipt.BlockNumber)
			require.GreaterOrEqual(t, receipt.Status, uint64(0))
			require.GreaterOrEqual(t, receipt.GasUsed, uint64(0))
			require.NotEqual(t, gethcommon.Address{}, receipt.From)
		}
	})

	t.Run("should handle multiple transactions per address", func(t *testing.T) {
		// add another tx including testAddress as to_address
		_, err := db.Exec("INSERT INTO tx (id, hash, content, to_address, type, sender_address, idx, batch_height, is_synthetic, time) VALUES (4, ?, ?, 1, 0, 2, 3, 100, false, 1234567893)",
			gethcommon.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa").Bytes(),
			[]byte("tx_content_4"))
		require.NoError(t, err)

		// another random tx
		_, err = db.Exec("INSERT INTO receipt (id, post_state, status, gas_used, effective_gas_price, created_contract_address, tx, batch) VALUES (4, ?, 1, 21000, 20000000000, NULL, 4, 1)",
			gethcommon.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb").Bytes())
		require.NoError(t, err)

		pagination := &common.QueryPagination{Size: 10, Offset: 0}
		receipts, err := GetTransactionsPerAddress(context.Background(), db, &testAddress, pagination)
		require.NoError(t, err)
		require.Len(t, receipts, 3) // Now should have 3 transactions
	})
}

func createTestSQLiteDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// manually creating the tables to avoid cyclical dependency if we import the sqlite package
	_, err = db.Exec(`
		CREATE TABLE externally_owned_account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			address BLOB(20) NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE batch (
			sequence INTEGER PRIMARY KEY,
			converted_hash BLOB(32) NOT NULL,
			hash BLOB(32) NOT NULL,
			height INTEGER NOT NULL,
			is_canonical BOOLEAN NOT NULL,
			header BLOB NOT NULL,
			l1_proof_hash BLOB(32) NOT NULL,
			is_executed BOOLEAN NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE tx (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hash BLOB(32) NOT NULL,
			content BLOB NOT NULL,
			to_address INTEGER,
			type INTEGER NOT NULL,
			sender_address INTEGER NOT NULL REFERENCES externally_owned_account(id),
			idx INTEGER NOT NULL,
			batch_height INTEGER NOT NULL,
			is_synthetic BOOLEAN NOT NULL,
			time INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE contract (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			address BLOB(20) NOT NULL,
			creator INTEGER NOT NULL REFERENCES externally_owned_account(id),
			auto_visibility BOOLEAN NOT NULL,
			transparent BOOLEAN,
			tx INTEGER NOT NULL REFERENCES tx(id)
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE receipt (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_state BLOB(32),
			status INTEGER NOT NULL,
			gas_used INTEGER NOT NULL,
			effective_gas_price INTEGER,
			created_contract_address BLOB(20),
			tx INTEGER NOT NULL REFERENCES tx(id),
			batch INTEGER NOT NULL REFERENCES batch(sequence)
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupTestData(t *testing.T, db *sqlx.DB) {
	testAddress := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	otherAddress := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	_, err := db.Exec("INSERT INTO externally_owned_account (id, address) VALUES (1, ?), (2, ?)",
		testAddress.Bytes(), otherAddress.Bytes())
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO batch (sequence, converted_hash, hash, height, is_canonical, header, l1_proof_hash, is_executed) VALUES (1, ?, ?, 100, true, ?, ?, true)",
		gethcommon.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111").Bytes(),
		gethcommon.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222").Bytes(),
		[]byte("header_data"),
		gethcommon.HexToHash("0x3333333333333333333333333333333333333333333333333333333333333333").Bytes())
	require.NoError(t, err)

	// tx 1: testAddress sends to otherAddress
	_, err = db.Exec("INSERT INTO tx (id, hash, content, to_address, type, sender_address, idx, batch_height, is_synthetic, time) VALUES (1, ?, ?, 2, 0, 1, 0, 100, false, 1234567890)",
		gethcommon.HexToHash("0x4444444444444444444444444444444444444444444444444444444444444444").Bytes(),
		[]byte("tx_content_1"))
	require.NoError(t, err)

	// tx 2: otherAddress sends to testAddress
	_, err = db.Exec("INSERT INTO tx (id, hash, content, to_address, type, sender_address, idx, batch_height, is_synthetic, time) VALUES (2, ?, ?, 1, 0, 2, 1, 100, false, 1234567891)",
		gethcommon.HexToHash("0x5555555555555555555555555555555555555555555555555555555555555555").Bytes(),
		[]byte("tx_content_2"))
	require.NoError(t, err)

	// tx 3: otherAddress sends to otherAddress
	_, err = db.Exec("INSERT INTO tx (id, hash, content, to_address, type, sender_address, idx, batch_height, is_synthetic, time) VALUES (3, ?, ?, 2, 0, 2, 2, 100, false, 1234567892)",
		gethcommon.HexToHash("0x6666666666666666666666666666666666666666666666666666666666666666").Bytes(),
		[]byte("tx_content_3"))
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO receipt (id, post_state, status, gas_used, effective_gas_price, created_contract_address, tx, batch) VALUES (1, ?, 1, 21000, 20000000000, NULL, 1, 1)",
		gethcommon.HexToHash("0x7777777777777777777777777777777777777777777777777777777777777777").Bytes())
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO receipt (id, post_state, status, gas_used, effective_gas_price, created_contract_address, tx, batch) VALUES (2, ?, 1, 21000, 20000000000, NULL, 2, 1)",
		gethcommon.HexToHash("0x8888888888888888888888888888888888888888888888888888888888888888").Bytes())
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO receipt (id, post_state, status, gas_used, effective_gas_price, created_contract_address, tx, batch) VALUES (3, ?, 1, 21000, 20000000000, NULL, 3, 1)",
		gethcommon.HexToHash("0x9999999999999999999999999999999999999999999999999999999999999999").Bytes())
	require.NoError(t, err)
}
