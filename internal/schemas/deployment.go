package schemas

import (
	"github.com/htquangg/a-wasm/internal/entities"

	"github.com/jinzhu/copier"
)

type AddDeploymentReq struct {
	EndpointID string `json:"endpointId"`
	Data       []byte `json:"data"`
}

type AddDeploymentResp struct {
	ID         string `json:"id"`
	EndpointID string `json:"endpointId"`
	Hash       string `json:"hash"`
	IngressURL string `json:"ingressUrl"`
	CreatedAt  int64  `json:"createdAt"`
}

func (r *AddDeploymentResp) SetFromDeployment(deployment *entities.Deployment) {
	_ = copier.Copy(r, deployment)

	r.CreatedAt = deployment.CreatedAt.Unix()
}
