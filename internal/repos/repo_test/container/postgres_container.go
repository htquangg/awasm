package container

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/htquangg/a-wasm/config"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	testcontainers.Container
	Host string
	Port int
}

func NewPostgresContainer(ctx context.Context, cfg *config.DB) (*PostgresContainer, error) {
	var (
		req  testcontainers.ContainerRequest
		port nat.Port
	)

	port = nat.Port(fmt.Sprintf("%d/tcp", cfg.Port))
	req = testcontainers.ContainerRequest{
		Image:        "postgres:12.2-alpine",
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", cfg.Port)},
		WaitingFor: wait.ForSQL(port, "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, host, port.Port(), cfg.Schema)
		}),
		Env: map[string]string{
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_DB":       cfg.Schema,
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	if err := container.StartLogProducer(ctx); err != nil {
		return nil, err
	}

	verbose, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("FLIPT_TEST_DATABASE_VERBOSE")))
	if verbose {
		var logger testContainerLogger
		container.FollowOutput(&logger)
	}

	mappedPort, err := container.MappedPort(ctx, port)
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{Container: container, Host: hostIP, Port: mappedPort.Int()}, nil
}

type testContainerLogger struct{}

func (t testContainerLogger) Accept(entry testcontainers.Log) {
	log.Println(entry.LogType, ":", string(entry.Content))
}
