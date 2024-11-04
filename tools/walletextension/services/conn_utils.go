package services

import (
	"context"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

func ConnectWS(ctx context.Context, account *GWAccount, logger gethlog.Logger) (*tenrpc.EncRPCClient, error) {
	return connect(ctx, account.user.services.RpcWSConnPool, account, logger)
}

func connect(ctx context.Context, p *pool.ObjectPool, account *GWAccount, logger gethlog.Logger) (*tenrpc.EncRPCClient, error) {
	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "get rpc connection")
	connectionObj, err := p.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	conn := connectionObj.(*rpc.Client)
	encClient, err := wecommon.CreateEncClient(conn, account.Address.Bytes(), account.user.userKey, account.signature, account.signatureType, logger)
	if err != nil {
		_ = ReturnConn(p, conn, logger)
		return nil, fmt.Errorf("error creating new client, %w", err)
	}
	return encClient, nil
}

func ReturnConn(p *pool.ObjectPool, conn tenrpc.Client, logger gethlog.Logger) error {
	err := p.ReturnObject(context.Background(), conn)
	if err != nil {
		logger.Error("Error returning connection to pool", log.ErrKey, err)
	}
	return err
}

func WithEncRPCConnection[R any](ctx context.Context, w *Services, acct *GWAccount, execute func(*tenrpc.EncRPCClient) (*R, error)) (*R, error) {
	rpcClient, err := connect(ctx, acct.user.services.RpcHTTPConnPool, acct, w.logger)
	if err != nil {
		return nil, fmt.Errorf("could not connect to backed. Cause: %w", err)
	}
	defer ReturnConn(w.RpcHTTPConnPool, rpcClient.BackingClient(), w.logger)
	return execute(rpcClient)
}

func WithPlainRPCConnection[R any](ctx context.Context, w *Services, execute func(client *rpc.Client) (*R, error)) (*R, error) {
	connectionObj, err := w.RpcHTTPConnPool.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	rpcClient := connectionObj.(*rpc.Client)
	defer ReturnConn(w.RpcHTTPConnPool, rpcClient, w.logger)
	return execute(rpcClient)
}
