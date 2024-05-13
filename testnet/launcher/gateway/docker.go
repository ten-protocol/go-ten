package gateway

import (
	"fmt"
	"time"

	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/valyala/fasthttp"
)

type DockerGateway struct {
	cfg *Config
}

func NewDockerGateway(cfg *Config) (*DockerGateway, error) {
	return &DockerGateway{
		cfg: cfg,
	}, nil // todo (@pedro) - add validation
}

func (n *DockerGateway) Start() error {
	fmt.Printf("Starting gateway with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{
		"/home/obscuro/go-obscuro/tools/walletextension/bin/wallet_extension_linux",
		"--host", "0.0.0.0",
		"--port", fmt.Sprintf("%d", n.cfg.gatewayHTTPPort),
		"--portWS", fmt.Sprintf("%d", n.cfg.gatewayWSPort),
		"--nodePortHTTP", fmt.Sprintf("%d", n.cfg.tenNodeHTTPPort),
		"--nodePortWS", fmt.Sprintf("%d", n.cfg.tenNodeWSPort),
		"--nodeHost", n.cfg.tenNodeHost,
		"--dbType", "sqlite",
		"--logPath", "sys_out",
	}

	_, err := docker.StartNewContainer("gateway", n.cfg.dockerImage, cmds, []int{n.cfg.gatewayHTTPPort, n.cfg.gatewayWSPort}, nil, nil, nil)
	return err
}

func (n *DockerGateway) IsReady() error {
	timeout := time.Minute
	interval := time.Second

	return retry.Do(func() error {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("http://127.0.0.1:%d/v1/health/", n.cfg.gatewayHTTPPort))
		if err != nil {
			return err
		}

		if statusCode != fasthttp.StatusOK {
			return fmt.Errorf("status not ok - status received: %s", fasthttp.StatusMessage(statusCode))
		}

		return nil
	}, retry.NewTimeoutStrategy(timeout, interval))
}
