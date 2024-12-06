package services

import (
	"context"
	"fmt"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type BackendRPC struct {
	// the OG maintains a connection pool of rpc connections to underlying nodes
	rpcHTTPConnPool *pool.ObjectPool
	rpcWSConnPool   *pool.ObjectPool
	encKey          []byte
	logger          gethlog.Logger
}

// todo - tweak the number of backend connections
const poolSize = 200

func NewBackendRPC(hostAddrHTTP string, hostAddrWS string, logger gethlog.Logger) *BackendRPC {
	factoryHTTP := pool.NewPooledObjectFactory(
		func(context.Context) (interface{}, error) {
			rpcClient, err := gethrpc.Dial(hostAddrHTTP)
			if err != nil {
				return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", hostAddrHTTP, err)
			}
			return rpcClient, nil
		}, func(ctx context.Context, object *pool.PooledObject) error {
			client := object.Object.(*gethrpc.Client)
			client.Close()
			return nil
		}, nil, nil, nil)

	factoryWS := pool.NewPooledObjectFactory(
		func(context.Context) (interface{}, error) {
			rpcClient, err := gethrpc.Dial(hostAddrWS)
			if err != nil {
				return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", hostAddrWS, err)
			}
			return rpcClient, nil
		}, func(ctx context.Context, object *pool.PooledObject) error {
			client := object.Object.(*gethrpc.Client)
			client.Close()
			return nil
		}, nil, nil, nil)

	cfg := pool.NewDefaultPoolConfig()
	cfg.MaxTotal = poolSize

	return &BackendRPC{
		rpcHTTPConnPool: pool.NewObjectPool(context.Background(), factoryHTTP, cfg),
		rpcWSConnPool:   pool.NewObjectPool(context.Background(), factoryWS, cfg),
		encKey:          readEncKey(hostAddrHTTP, logger),
		logger:          logger,
	}
}

func readEncKey(hostAddrHTTP string, logger gethlog.Logger) []byte {
	// read the encryption key
	rpcClient, err := gethrpc.Dial(hostAddrHTTP)
	if err != nil {
		logger.Crit("failed to connect to the node", "err", err)
		return nil
	}
	defer rpcClient.Close()
	n := 0
	for {
		n++
		k, err := tenrpc.ReadEnclaveKey(rpcClient)
		if err != nil {
			logger.Warn("failed to read enc key", "err", err)
			if n > 10 { // wait for ~1m for the backend node to spin up and respond
				logger.Crit("failed to read enc key", "err", err)
				return nil
			}
			time.Sleep(time.Duration(n) * time.Second)
		} else {
			return k
		}
	}
}

func (rpc *BackendRPC) ConnectWS(ctx context.Context, account *wecommon.GWAccount) (*tenrpc.EncRPCClient, error) {
	return connect(ctx, rpc.rpcWSConnPool, account, rpc.encKey, rpc.logger)
}

func (rpc *BackendRPC) ReturnConnWS(conn tenrpc.Client) error {
	return returnConn(rpc.rpcWSConnPool, conn, rpc.logger)
}

func (rpc *BackendRPC) ConnectHttp(ctx context.Context, account *wecommon.GWAccount) (*tenrpc.EncRPCClient, error) {
	return connect(ctx, rpc.rpcHTTPConnPool, account, rpc.encKey, rpc.logger)
}

func (rpc *BackendRPC) PlainConnectWs(ctx context.Context) (*gethrpc.Client, error) {
	return connectPlain(ctx, rpc.rpcWSConnPool, rpc.logger)
}

func (rpc *BackendRPC) ReturnConn(conn tenrpc.Client) error {
	return returnConn(rpc.rpcHTTPConnPool, conn, rpc.logger)
}

func (rpc *BackendRPC) Stop() {
	rpc.rpcHTTPConnPool.Close(context.Background())
	rpc.rpcWSConnPool.Close(context.Background())
}

func WithEncRPCConnection[R any](ctx context.Context, rpc *BackendRPC, acct *wecommon.GWAccount, execute func(*tenrpc.EncRPCClient) (*R, error)) (*R, error) {
	rpcClient, err := connect(ctx, rpc.rpcHTTPConnPool, acct, rpc.encKey, rpc.logger)
	if err != nil {
		return nil, fmt.Errorf("could not connect to backed. Cause: %w", err)
	}
	defer rpc.ReturnConn(rpcClient.BackingClient())
	return execute(rpcClient)
}

func WithPlainRPCConnection[R any](ctx context.Context, b *BackendRPC, execute func(client *rpc.Client) (*R, error)) (*R, error) {
	connectionObj, err := connectPlain(ctx, b.rpcHTTPConnPool, b.logger)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	defer b.ReturnConn(connectionObj)
	return execute(connectionObj)
}

func connectPlain(ctx context.Context, p *pool.ObjectPool, logger gethlog.Logger) (*rpc.Client, error) {
	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "get rpc connection")
	connectionObj, err := p.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	conn := connectionObj.(*rpc.Client)
	return conn, nil
}

func connect(ctx context.Context, p *pool.ObjectPool, account *wecommon.GWAccount, key []byte, logger gethlog.Logger) (*tenrpc.EncRPCClient, error) {
	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "get rpc connection")
	connectionObj, err := p.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	conn := connectionObj.(*rpc.Client)
	encClient, err := wecommon.CreateEncClient(conn, key, account.Address.Bytes(), account.User.UserKey, account.Signature, account.SignatureType, logger)
	if err != nil {
		_ = returnConn(p, conn, logger)
		return nil, fmt.Errorf("error creating new client, %w", err)
	}
	return encClient, nil
}

func returnConn(p *pool.ObjectPool, conn tenrpc.Client, logger gethlog.Logger) error {
	err := p.ReturnObject(context.Background(), conn)
	if err != nil {
		logger.Error("Error returning connection to pool", log.ErrKey, err)
	}
	return err
}
