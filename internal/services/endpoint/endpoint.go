package endpoint

import "github.com/htquangg/a-wasm/internal/db"

type Endpoint struct {
	db db.DB
}

func New(db db.DB) *Endpoint {
	return &Endpoint{
		db: db,
	}
}
