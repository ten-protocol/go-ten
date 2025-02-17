package webserver

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/ten-protocol/go-ten/tools/faucet/faucet"
)

type WebServer struct {
	engine      *gin.Engine
	faucet      *faucet.Faucet
	bindAddress string
	server      *http.Server
}

type requestAddr struct {
	Address string `json:"address" binding:"required"`
}

func NewWebServer(faucetServer *faucet.Faucet, bindAddress string, jwtSecret []byte, defaultAmount *big.Int) *WebServer {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// authed endpoint
	r.POST("/auth/fund/:token", jwtTokenChecker(jwtSecret, faucetServer.Logger), fundingHandler(faucetServer, defaultAmount))

	// todo (@matt) we need to remove this unsecure endpoint before we provide a fully public sepolia faucet
	r.POST("/fund/:token", fundingHandler(faucetServer, defaultAmount))

	r.GET("/balance", balanceReqHandler(faucetServer))

	r.GET("/health", healthReqHandler())

	return &WebServer{
		engine:      r,
		faucet:      faucetServer,
		bindAddress: bindAddress,
	}
}

func jwtTokenChecker(jwtSecret []byte, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			errorHandler(c, err, logger)
			return
		}

		_, err = faucet.ValidateToken(jwtToken, jwtSecret)
		if err != nil {
			errorHandler(c, err, logger)
			return
		}

		c.Next()
	}
}

func (w *WebServer) Start() error {
	w.server = &http.Server{
		Addr:              w.bindAddress,
		Handler:           w.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (w *WebServer) Stop() error {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return w.server.Shutdown(ctx)
}

func errorHandler(c *gin.Context, err error, logger log.Logger) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
	logger.Error(err.Error())
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", fmt.Errorf("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func fundingHandler(faucetServer *faucet.Faucet, defaultAmount *big.Int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReq := c.Params.ByName("token")
		var token string

		// check the token request type
		switch tokenReq {
		case faucet.NativeToken:
			token = faucet.NativeToken
		// we leave this option in temporarily for tools that are still using `/ten` endpoint for native funds
		case faucet.DeprecatedNativeToken:
			token = faucet.NativeToken
		case faucet.WrappedOBX:
			token = faucet.WrappedOBX
		case faucet.WrappedEth:
			token = faucet.WrappedEth
		case faucet.WrappedUSDC:
			token = faucet.WrappedUSDC
		default:
			errorHandler(c, fmt.Errorf("token not recognized: %s", tokenReq), faucetServer.Logger)
			return
		}

		// make sure there's an address
		var req requestAddr
		if err := c.Bind(&req); err != nil {
			errorHandler(c, fmt.Errorf("unable to parse request: %w", err), faucetServer.Logger)
			return
		}

		// make sure the address is valid
		if !common.IsHexAddress(req.Address) {
			errorHandler(c, fmt.Errorf("unexpected address %s", req.Address), faucetServer.Logger)
			return
		}

		// fund the address
		addr := common.HexToAddress(req.Address)
		hash, err := faucetServer.Fund(&addr, token, defaultAmount)
		if err != nil {
			errorHandler(c, fmt.Errorf("unable to fund request %w", err), faucetServer.Logger)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "tx": hash})
	}
}

// returns the remaining native balance of the faucet
func balanceReqHandler(faucetServer *faucet.Faucet) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the balance
		balance, err := faucetServer.Balance(c)
		if err != nil {
			errorHandler(c, fmt.Errorf("unable to get balance %w", err), faucetServer.Logger)
			return
		}

		c.JSON(http.StatusOK, gin.H{"balance": balance.String()})
	}
}

// returns the remaining native balance of the faucet
func healthReqHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healthy": true})
	}
}
