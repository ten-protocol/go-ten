package postgres

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/common/retry"
)

type Postgres struct{}

func NewDockerPostgres() (*Postgres, error) {
	return &Postgres{}, nil
}

func (p *Postgres) Start() error {
	fmt.Println("starting Postgres container 'pg-ten'")

	// stop and remove any existing pg-ten container in case you forget to clear containers before a re-run
	_ = docker.StopAndRemove("pg-ten")

	envs := map[string]string{
		"POSTGRES_PASSWORD": "postgres",
	}

	_, err := docker.StartNewContainer(
		"pg-ten",
		"postgres:16-alpine",
		nil, // no custom command
		nil, // no exposed ports to host
		envs,
		nil,   // no devices
		nil,   // no volumes
		false, // no auto-restart
	)

	return err
}

func (p *Postgres) IsReady() error {
	timeout := 2 * time.Minute
	interval := 1 * time.Second

	fmt.Println("Waiting for Postgres to become ready...")
	err := retry.Do(func() error {
		// we don't expose any TCP ports so we can't use the existing functions to check for readiness
		return docker.ExecInContainer("pg-ten", []string{"pg_isready", "-U", "postgres", "-h", "localhost"})
	}, retry.NewTimeoutStrategy(timeout, interval))

	if err != nil {
		return fmt.Errorf("postgres did not become ready in time: %w", err)
	}

	fmt.Println("Postgres is ready")
	return nil
}
