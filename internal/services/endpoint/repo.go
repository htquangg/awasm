package endpoint

import "github.com/htquangg/a-wasm/internal/db"

type endpointRepo struct {
	db db.DB
}

func newEndpointRepo(db db.DB) *endpointRepo {
	return &endpointRepo{
		db: db,
	}
}
