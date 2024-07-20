package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/fatih/color"
	"github.com/segmentfault/pacman"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/cache"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/middleware"
	"github.com/htquangg/awasm/internal/base/translator"
	"github.com/htquangg/awasm/internal/constants"
	"github.com/htquangg/awasm/internal/controllers"
	"github.com/htquangg/awasm/internal/protocluster"
	"github.com/htquangg/awasm/internal/repos"
	"github.com/htquangg/awasm/internal/services"
	"github.com/htquangg/awasm/internal/web"
	"github.com/htquangg/awasm/pkg/logger"
)

var runCmd = &cobra.Command{
	Use:                   "run",
	Short:                 "Run the application",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Awasm is starting..........................")
		return runApp()
	},
}

func init() {
	runCmd.PersistentFlags().String("config-path", "", "Specify the config path of the application")
	ensure(viper.BindPFlag("server.config-path", runCmd.PersistentFlags().Lookup("config-path")))
}

func runApp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load the configuration from the file
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Setup application logger
	logger.SetLogger(logger.NewZapLogger(
		logger.WithZapFilename(cfg.Logging.Filename),
		logger.WithZapLevel(cfg.Logging.Level),
		logger.WithZapMaxSize(cfg.Logging.MaxSize),
		logger.WithZapMaxBackups(cfg.Logging.MaxBackups),
		logger.WithZapMaxAge(cfg.Logging.MaxAge),
		logger.WithZapLocalTime(cfg.Logging.UseLocalTime),
		logger.WithZapCompress(cfg.Logging.UseCompress),
	))

	// Setup dependencies and application
	app, err := initApp(ctx, cfg)
	if err != nil {
		return err
	}

	constants.Version = Version
	constants.Revision = Revision
	constants.GoVersion = GoVersion
	regular := color.New()
	regular.Println("awasm Version:", constants.Version, " Revision:", constants.Revision)

	if err := app.Run(context.Background()); err != nil {
		return err
	}

	return nil
}

func initApp(ctx context.Context, cfg *config.Config) (*pacman.Application, error) {
	_, err := translator.NewTranslator(cfg.I18n)
	if err != nil {
		return nil, err
	}

	db, err := db.New(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	cluster := protocluster.New(ctx, db)

	repos := repos.New(cfg, db, cache)

	services := services.New(cfg, repos, cluster)

	mws := middleware.NewMiddleware(cfg, cache, services, repos)
	controllers := controllers.New(services)

	return pacman.NewApp(
		pacman.WithName(Name),
		pacman.WithVersion(Version),
		pacman.WithSignals(
			[]os.Signal{
				os.Interrupt,
				os.Kill,
				syscall.SIGTERM,
				syscall.SIGQUIT,
				syscall.SIGINT,
				syscall.SIGHUP,
			},
		),
		pacman.WithServer(
			web.New(ctx, cfg.HTTP, controllers, mws),
			cluster,
		),
	), nil
}
