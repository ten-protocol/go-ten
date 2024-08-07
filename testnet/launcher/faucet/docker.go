package faucet

import (
	"fmt"
	"time"

	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/valyala/fasthttp"
)

type DockerFaucet struct {
	cfg *Config
}

func NewDockerFaucet(cfg *Config) (*DockerFaucet, error) {
	return &DockerFaucet{
		cfg: cfg,
	}, nil // todo (@pedro) - add validation
}

func (n *DockerFaucet) Start() error {
	fmt.Printf("Starting faucet with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{
		"/home/obscuro/go-obscuro/tools/faucet/cmd/faucet",
		"--nodeHost", n.cfg.tenNodeHost,
		"--nodePort", fmt.Sprintf("%d", n.cfg.tenNodePort),
		"--pk", n.cfg.faucetPrivKey,
		"--jwtSecret", "someKey",
		"--serverPort", fmt.Sprintf("%d", n.cfg.faucetPort),
	}

	_, err := docker.StartNewContainer("faucet", n.cfg.dockerImage, cmds, []int{n.cfg.faucetPort}, nil, nil, nil, false)
	return err
}

func (n *DockerFaucet) IsReady() error {
	timeout := time.Minute
	interval := time.Second

	return retry.Do(func() error {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("http://127.0.0.1:%d/health/", n.cfg.faucetPort))
		if err != nil {
			return err
		}

		if statusCode != fasthttp.StatusOK {
			return fmt.Errorf("status not ok - status received: %s", fasthttp.StatusMessage(statusCode))
		}

		return nil
	}, retry.NewTimeoutStrategy(timeout, interval))
}
