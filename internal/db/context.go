package db

import (
	"context"
)

type contextKey struct {
	name string
}

var (
	enginedContextKey         = &contextKey{"engined"}
	_                 Engined = &Context{}
)

type Context struct {
	context.Context
	e           Engine
	transaction bool
}

func newContext(ctx context.Context, e Engine, transaction bool) *Context {
	return &Context{
		Context:     ctx,
		e:           e,
		transaction: transaction,
	}
}

func (ctx *Context) Engine() Engine {
	return ctx.e
}

func (ctx *Context) Value(key any) any {
	if key == enginedContextKey {
		return ctx
	}
	return ctx.Context.Value(key)
}

type Engined interface {
	Engine() Engine
}

// Committer represents an interface to Commit or Close the Context.
type Committer interface {
	Commit() error
	Close() error
}

// halfCommitter is a wrapper of Committer.
// It can be closed early, but can't be committed early, it is useful for reusing a transaction.
type halfCommitter struct {
	committer Committer
	committed bool
}

func (c *halfCommitter) Commit() error {
	c.committed = true
	// should do nothing, and the parent committer will commit later.
	return nil
}

func (c *halfCommitter) Close() error {
	if c.committed {
		// it's "commit and close", should do nothing, and the parent committer will commit later.
		return nil
	}

	// it's "rollback and close", let the parent committer rollback right now.
	return c.committer.Close()
}
