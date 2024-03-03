package mws

import (
	"context"

	db_internal "github.com/htquangg/a-wasm/internal/base/db"
	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/log"
)

const (
	ContextTxKey string = "tx"
)

type TransactionalMiddleware = echo.MiddlewareFunc

func Transactional(db db_internal.DB) TransactionalMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			switch ctx.Request().Method {
			case "POST", "PUT", "PATCH", "DELETE":
				ctxTx := getCtxTxFromContext(ctx)
				if ctxTx != nil {
					return next(ctx)
				}

				var committer db_internal.Committer

				ctxTx, committer, err = db.TxContext(context.Background())
				if err != nil {
					return err
				}

				ctx.Set(ContextTxKey, ctxTx)

				defer func(committer db_internal.Committer) {
					err = closeTx(committer, err)
				}(committer)
			}

			return next(ctx)
		}
	}
}

func GetCtxTxFromContext(ctx echo.Context) context.Context {
	ctxTx := getCtxTxFromContext(ctx)
	if ctxTx == nil {
		log.Warnf("has no tx context in route: %s", ctx.Path())
		return context.Background()
	}

	return ctxTx
}

func getCtxTxFromContext(ctx echo.Context) context.Context {
	ctxTx, _ := ctx.Get(ContextTxKey).(context.Context)
	if ctxTx == nil {
		return nil
	}

	return ctxTx
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
