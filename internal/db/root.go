package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // Needed for the MySQL driver
	"xorm.io/xorm"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"
)

var x *xorm.Engine

type (
	DB interface {
		Engine(ctx context.Context) Engine
		TxContext(parentCtx context.Context) (*Context, Committer, error)
		WithTx(parentCtx context.Context, f func(ctx context.Context) error) error
		InTransaction(ctx context.Context) bool
		Exec(ctx context.Context, sqlAndArgs ...any) (sql.Result, error)
		Query(ctx context.Context, sqlAndArgs ...any) ([]map[string][]byte, error)
	}

	db struct {
		ctx context.Context
		cfg *Config
		e   *xorm.Engine
	}
)

func New(ctx context.Context, cfg *Config) (DB, error) {
	conn := cfg.Address()

	engine, err := xorm.NewEngine(string(schemas.MYSQL), conn)
	if err != nil {
		return nil, err
	}

	engine.SetMapper(names.GonicMapper{})
	engine.ShowSQL(cfg.LogSQL)
	engine.SetMaxOpenConns(cfg.MaxOpenConns)
	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))
	engine.SetTZDatabase(time.UTC)
	engine.SetDefaultContext(ctx)

	return &db{
		ctx: ctx,
		cfg: cfg,
		e:   engine,
	}, nil
}

func SetDefaultEngine(ctx context.Context, eng *xorm.Engine) {
	x = eng
	DefaultContext = &Context{
		Context: ctx,
		e:       x,
	}
}

func GetEngine(ctx context.Context) Engine {
	if e := getEngine(ctx); e != nil {
		return e
	}
	return x.Context(ctx)
}

func getEngine(ctx context.Context) Engine {
	if engined, ok := ctx.(Engined); ok {
		return engined.Engine()
	}
	enginedInterface := ctx.Value(enginedContextKey)
	if enginedInterface != nil {
		return enginedInterface.(Engined).Engine()
	}
	return nil
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

func (db *db) TxContext(parentCtx context.Context) (*Context, Committer, error) {
	if sess, ok := db.inTransaction(parentCtx); ok {
		return newContext(parentCtx, sess, true), &halfCommitter{committer: sess}, nil
	}

	sess := db.e.NewSession()
	if err := sess.Begin(); err != nil {
		return nil, nil, err
	}
	return newContext(db.ctx, sess, true), sess, nil
}

func (db *db) WithTx(parentCtx context.Context, f func(ctx context.Context) error) error {
	if sess, ok := db.inTransaction(parentCtx); ok {
		err := f(newContext(parentCtx, sess, true))
		if err != nil {
			// rollback immediately, in case the caller ignores returned error and tries to commit the transaction.
			_ = sess.Close()
		}
		return err
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

func (db *db) txWithNoCheck(parentCtx context.Context, f func(ctx context.Context) error) error {
	sess := db.e.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := f(newContext(parentCtx, sess, true)); err != nil {
		return err
	}

	return sess.Commit()
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
