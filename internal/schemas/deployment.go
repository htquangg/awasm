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
	CreatedAt  int64  `json:"createdAt"`
	IngressURL string `json:"ingressUrl"`
}

func (r *AddDeploymentResp) SetFromDeployment(deployment *entities.Deployment) {
	_ = copier.Copy(r, deployment)

	r.CreatedAt = deployment.CreatedAt.Unix()
}
