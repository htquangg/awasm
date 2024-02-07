package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/cluster"
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services"
	"github.com/htquangg/a-wasm/internal/web"

	"github.com/oklog/run"
	"github.com/pkg/errors"
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

	db, err := db.New(ctx, cfg.DB)
	if err != nil {
		return g, err
	}

	svc := services.New(db)

	g.Add(web.
		New(ctx, &web.Config{
			ShowStartBanner: true,
			Addr:            cfg.Server.HTTP.Addr,
		},
			svc,
			db,
		).
		ServeHandler(),
	)

	g.Add(cluster.New(ctx).ServeHandler())

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
