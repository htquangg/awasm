package repo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/repos/repo_test/container"
	"github.com/htquangg/awasm/pkg/logger"
)

type Database struct {
	DB        db.DB
	Container *container.PostgresContainer
}

func (d *Database) Shutdown(ctx context.Context) {
	if d.DB != nil {
		d.DB.Shutdown(ctx)
	}

	if d.Container != nil {
		_ = d.Container.StopLogProducer()
		_ = d.Container.Terminate(ctx)
	}
}

// testDB used for repo testing
var testDB db.DB

func TestMain(t *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	d, err := initTestDB(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if d != nil {
			d.Shutdown(shutdownCtx)
		}
	}
	defer cleanup()

	testDB = d.DB

	logger.Info("init test database successfully")

	if ret := t.Run(); ret != 0 {
		panic(ret)
	}
}

func initTestDB(ctx context.Context, cfg *config.DB) (*Database, error) {
	container, err := container.NewPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating db container: %w", err)
	}
	cfg.Host = container.Host
	cfg.Port = container.Port

	dbMigrate, err := goose.OpenDBWithDriver("postgres", cfg.Address())
	if err != nil {
		return nil, fmt.Errorf("creating migrate instance: %w", err)
	}
	// run drop to clear target DB (incase we're reusing)
	if err := goose.RunContext(ctx, "reset", dbMigrate, cfg.MigrationDirPath); err != nil &&
		!strings.Contains(err.Error(), "\"goose_db_version\" does not exist") {
		return nil, fmt.Errorf("running drop: %w, %s", err, err.Error())
	}
	if err := dbMigrate.Close(); err != nil {
		return nil, fmt.Errorf("closing db: %w", err)
	}

	// need new instance after drop
	dbMigrate, err = goose.OpenDBWithDriver("postgres", cfg.Address())
	if err != nil {
		return nil, fmt.Errorf("creating migrate instance: %w", err)
	}
	// run drop to clear target DB (incase we're reusing)
	if err := goose.RunContext(ctx, "up", dbMigrate, cfg.MigrationDirPath); err != nil &&
		!strings.Contains(err.Error(), "\"goose_db_version\" does not exist") {
		return nil, fmt.Errorf("running drop: %w, %s", err, err.Error())
	}
	if err := dbMigrate.Close(); err != nil {
		return nil, fmt.Errorf("closing db: %w", err)
	}

	db, err := db.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Database{
		DB:        db,
		Container: container,
	}, nil
}
