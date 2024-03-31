package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/base/translator"
	"github.com/htquangg/a-wasm/internal/controllers"
	"github.com/htquangg/a-wasm/internal/protocluster"
	"github.com/htquangg/a-wasm/internal/repos"
	"github.com/htquangg/a-wasm/internal/services"
	"github.com/htquangg/a-wasm/internal/web"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  "Run the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Awasm is starting..........................")
		runApp()
	},
}

func runApp() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	g, err := initApp(ctx, cfg)
	if err != nil {
		panic(err)
	}

	var se run.SignalError
	if err := g.Run(); err != nil && !errors.As(err, &se) {
		panic(err)
	}
}

func initApp(ctx context.Context, cfg *config.Config) (run.Group, error) {
	var g run.Group

	_, err := translator.NewTranslator(cfg.I18n)
	if err != nil {
		return g, err
	}

	db, err := db.New(ctx, cfg.DB)
	if err != nil {
		return g, err
	}

	cache, err := cache.New(ctx, cfg.Redis)
	if err != nil {
		return g, err
	}

	cluster := protocluster.New(ctx, db)

	g.Add(cluster.ServeHandler())

	repos := repos.New(cfg, db, cache)
	services := services.New(cfg, repos, cluster)

	mws := middleware.NewMiddleware(cfg, repos)
	controllers := controllers.New(services)

	g.Add(web.
		New(ctx, cfg.Server, controllers, mws).
		ServeHandler(),
	)

	g.Add(run.SignalHandler(ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	))

	return g, nil
}
