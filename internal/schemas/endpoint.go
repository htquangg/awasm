package schemas

import (
	"github.com/htquangg/a-wasm/internal/entities"

	"github.com/jinzhu/copier"
)

var Runtimes = map[string]bool{
	"js": true,
	"go": true,
}

func ValidRuntime(runtime string) bool {
	_, ok := Runtimes[runtime]
	return ok
}

// AddEndpointReq holds all necesssary fields to create new run application.
type AddEndpointReq struct {
	// Name of the endpoint
	Name string `json:"name"`
	// Runtime on which the code will be invoked. (go or js for now)
	Runtime string `json:"runtime"`
}

type AddEndpointResp struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Runtime   string `json:"runtime"`
	CreatedAt int64  `json:"createdAt"`
}

func (r *AddEndpointResp) SetFromEndpoint(endpoint *entities.Endpoint) {
	_ = copier.Copy(r, endpoint)

	r.CreatedAt = endpoint.CreatedAt.Unix()
}
