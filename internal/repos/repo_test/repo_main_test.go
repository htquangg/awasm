package repo_test

import (
	"context"
	"testing"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/db"

	"github.com/segmentfault/pacman/log"
)

type TestDBSetting struct {
	Driver       string
	ImageName    string
	ImageVersion string
	ENV          []string
	PortID       string
	Connection   string
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

	if err := initTestDB(ctx, cfg.DB); err != nil {
		panic(err)
	}
	log.Info("init test database successfully")

	if ret := t.Run(); ret != 0 {
		panic(ret)
	}
}

func initTestDB(ctx context.Context, cfg *config.DB) error {
	db, err := db.New(ctx, cfg)
	if err != nil {
		return err
	}

	testDB = db

	return nil
}
