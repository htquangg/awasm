package db

import (
	"context"
	"database/sql"

	"github.com/htquangg/a-wasm/config"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

// Engine represents a xorm engine or session.
type Engine interface {
	Table(tableNameOrBean any) *xorm.Session
	Count(...any) (int64, error)
	Asc(colNames ...string) *xorm.Session
	Desc(colNames ...string) *xorm.Session
	Incr(column string, arg ...any) *xorm.Session
	Decr(column string, arg ...any) *xorm.Session
	Delete(...any) (int64, error)
	Truncate(...any) (int64, error)
	Exec(...any) (sql.Result, error)
	Find(any, ...any) error
	Get(beans ...any) (bool, error)
	ID(any) *xorm.Session
	In(string, ...any) *xorm.Session
	Insert(...any) (int64, error)
	Iterate(any, xorm.IterFunc) error
	Join(joinOperator string, tablename, condition any, args ...any) *xorm.Session
	SQL(any, ...any) *xorm.Session
	Where(any, ...any) *xorm.Session
	Limit(limit int, start ...int) *xorm.Session
	NoAutoTime() *xorm.Session
	SumInt(bean any, columnName string) (res int64, err error)
	Sync(...any) error
	Select(string) *xorm.Session
	NotIn(string, ...any) *xorm.Session
	OrderBy(any, ...any) *xorm.Session
	Exist(...any) (bool, error)
	Distinct(...string) *xorm.Session
	Query(...any) ([]map[string][]byte, error)
	Cols(...string) *xorm.Session
	Context(ctx context.Context) *xorm.Session
	GroupBy(keys string) *xorm.Session
	Ping() error
}

func NewXORMEngine(cfg *config.DB) (*xorm.Engine, error) {
	conn := cfg.Address()

	engine, err := xorm.NewEngine(string(schemas.POSTGRES), conn)
	if err != nil {
		return nil, err
	}

	return engine, nil
}
