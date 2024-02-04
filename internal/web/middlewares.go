package web

import (
	"context"

	db_internal "github.com/htquangg/a-wasm/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	ContextTxKey string = "tx"
)

func Transactional(db db_internal.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			switch c.Request().Method {
			case "POST", "PUT", "PATCH", "DELETE":
				ctxTx := getContextTx(c)
				if ctxTx != nil {
					return next(c)
				}

				var commiter db_internal.Committer

				ctxTx, commiter, err = db.TxContext(context.Background())
				if err != nil {
					return err
				}

				c.Set(ContextTxKey, ctxTx)

				defer func(committer db_internal.Committer) {
					err = closeTx(committer, err)
				}(commiter)
			}

			return next(c)
		}
	}
}

func GetContextTx(c echo.Context) context.Context {
	ctx := getContextTx(c)
	if ctx == nil {
		log.Warn().Msgf("has no tx context in route: %s", c.Path())
		return context.Background()
	}

	return ctx
}

func getContextTx(c echo.Context) context.Context {
	ctx, _ := c.Get(ContextTxKey).(context.Context)
	if ctx == nil {
		return nil
	}

	return ctx
}

func closeTx(committer db_internal.Committer, err error) error {
	if p := recover(); p != nil {
		_ = committer.Close()
		panic(p)
	} else if err != nil {
		_ = committer.Close()
		return err
	} else {
		return committer.Commit()
	}
}
