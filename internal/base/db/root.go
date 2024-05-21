package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"

	"github.com/htquangg/a-wasm/config"
)

type (
	DB interface {
		Engine(ctx context.Context) Engine
		WithTx(
			parentCtx context.Context,
			f func(ctx context.Context) (interface{}, error),
		) (interface{}, error)
		InTransaction(ctx context.Context) bool
		Exec(ctx context.Context, sqlAndArgs ...any) (sql.Result, error)
		Query(ctx context.Context, sqlAndArgs ...any) ([]map[string][]byte, error)
		Shutdown(ctx context.Context) error
	}

	db struct {
		ctx context.Context
		cfg *config.DB
		e   *xorm.Engine
	}
)

func New(ctx context.Context, cfg *config.DB) (DB, error) {
	conn := cfg.Address()

	engine, err := xorm.NewEngine(string(schemas.POSTGRES), conn)
	if err != nil {
		return nil, err
	}

	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	engine.SetMapper(names.GonicMapper{})
	engine.SetLogger(NewXORMLogger(cfg.LogSQL))
	engine.ShowSQL(cfg.LogSQL)
	engine.SetMaxOpenConns(cfg.MaxOpenConns)
	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))
	engine.SetTZDatabase(time.UTC)
	engine.SetDefaultContext(ctx)

	// register the go_sql_stats_connections_* metrics
	if err := prometheus.Register(sqlstats.NewStatsCollector("grafana", engine.DB().DB)); err != nil {
		return nil, err
	}

	return &db{
		ctx: ctx,
		cfg: cfg,
		e:   engine,
	}, nil
}

func (db *db) Engine(ctx context.Context) Engine {
	if e := db.engine(ctx); e != nil {
		return e
	}

	return db.e.Context(ctx)
}

func (db *db) engine(ctx context.Context) Engine {
	if engined, ok := ctx.(Engined); ok {
		return engined.Engine()
	}
	enginedInterface := db.ctx.Value(enginedContextKey)
	if enginedInterface != nil {
		return enginedInterface.(Engined).Engine()
	}
	return nil
}

func (db *db) WithTx(
	parentCtx context.Context,
	f func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	if sess, ok := db.inTransaction(parentCtx); ok {
		result, err := f(newContext(parentCtx, sess, true))
		if err != nil {
			// rollback immediately, in case the caller ignores returned error and tries to commit the transaction.
			_ = sess.Close()
		}
		return result, err
	}
	return db.txWithNoCheck(parentCtx, f)
}

func (db *db) InTransaction(ctx context.Context) bool {
	_, ok := db.inTransaction(ctx)
	return ok
}

func (db *db) Exec(ctx context.Context, sqlAndArgs ...any) (sql.Result, error) {
	return db.Engine(ctx).Exec(sqlAndArgs...)
}

func (db *db) Query(ctx context.Context, sqlAndArgs ...any) ([]map[string][]byte, error) {
	return db.Engine(ctx).Query(sqlAndArgs...)
}

func (db *db) txWithNoCheck(
	parentCtx context.Context,
	f func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	sess := db.e.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return nil, err
	}

	result, err := f(newContext(parentCtx, sess, true))
	if err != nil {
		return result, err
	}

	if err := sess.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

func (db *db) inTransaction(ctx context.Context) (*xorm.Session, bool) {
	e := db.engine(ctx)
	if e == nil {
		return nil, false
	}

	switch t := e.(type) {
	case *xorm.Engine:
		return nil, false
	case *xorm.Session:
		if t.IsInTx() {
			return t, true
		}
		return nil, false
	default:
		return nil, false
	}
}

func (db *db) Shutdown(ctx context.Context) error {
	return db.e.Close()
}
