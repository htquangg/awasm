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
