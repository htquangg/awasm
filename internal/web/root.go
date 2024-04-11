package web

import (
	"context"
	std_errors "errors"
	"fmt"
	std_log "log"
	"net/http"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/controllers"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	middleware_echo "github.com/labstack/echo/v4/middleware"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

type Server struct {
	ctx context.Context
	cfg *config.Server
	e   *echo.Echo
}

func New(
	ctx context.Context,
	cfg *config.Server,
	controllers *controllers.Controllers,
	mws *middleware.Middleware,
) *Server {
	e := echo.New()

	e.Pre(middleware_echo.RemoveTrailingSlashWithConfig(middleware_echo.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			// enable by default only for the API routes
			return !strings.HasPrefix(c.Request().URL.Path, constants.LiveIngressPath) ||
				!strings.HasPrefix(c.Request().URL.Path, constants.PreviewIngressPath) ||
				!strings.HasPrefix(c.Request().URL.Path, "/api/")
		},
	}))
	e.Use(middleware_echo.Recover())
	e.Use(middleware_echo.Secure())

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		if v := new(echo.HTTPError); std_errors.As(err, &v) {
			reason := ""
			if v.Internal != nil {
				reason = v.Error()
			} else {
				reason = fmt.Sprintf("%v", v.Message)
			}
			err = errors.New(v.Code, reason)
		}

		handler.HandleResponse(ctx, err, nil)
	}

	v1Group := e.Group("/api/v1")

	bindHealthApi(v1Group, controllers)
	bindEndpointsApi(v1Group, controllers, mws)
	bindPreviewApi(v1Group, controllers, mws)
	bindLiveApi(v1Group, controllers, mws)
	bindUserApi(v1Group, controllers, mws)

	// catch all any route
	v1Group.Any("/*", func(ctx echo.Context) error {
		return echo.ErrNotFound
	}, middleware_echo.Logger())

	return &Server{
		ctx: ctx,
		cfg: cfg,
		e:   e,
	}
}

func (s *Server) ServeHandler() (execute func() error, interrupt func(error)) {
	s.e.Server.Addr = s.cfg.Addr

	server := &http.Server{
		ReadTimeout:       constants.ReadTimeout,
		ReadHeaderTimeout: constants.ReadHeaderTimeout,
		WriteTimeout:      constants.WriteTimeout,
		Handler:           s.e,
		Addr:              s.cfg.Addr,
	}

	return func() error {
			if s.cfg.ShowStartBanner {
				addr := server.Addr
				schema := "http"

				date := new(strings.Builder)
				std_log.New(date, "", std_log.LstdFlags).Print()

				bold := color.New(color.Bold).Add(color.FgGreen)
				bold.Printf(
					"%s Web server started at %s\n",
					strings.TrimSpace(date.String()),
					color.CyanString("%s://%s", schema, addr),
				)

				regular := color.New()
				regular.Printf("├─ REST API: %s\n", color.CyanString("%s://%s/api/", schema, addr))
			}

			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Warnf("Web server failed to stop gracefully: %v", err)
			}
		}
}
