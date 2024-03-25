package entities

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/htquangg/a-wasm/pkg/uid"
)

type Deployment struct {
	ID         string `xorm:"not null pk VARCHAR(36) id"`
	EndpointID string `xorm:"not null VARCHAR(36) endpoint_id"`
	Hash       string `xorm:"not null VARCHAR(32) hash"`
	Data       []byte `xorm:"not null MEDIUMBLOB data"`

	CreatedAt time.Time  `xorm:"created TIMESTAMPZ created_at"`
	DeletedAt *time.Time `xorm:"TIMESTAMPZ deleted_at"`
}

func (Deployment) TableName() string {
	return "deployments"
}

func NewFromEndpoint(endpoint *Endpoint, data []byte) *Deployment {
	hashBytes := md5.Sum(data)
	hashstr := hex.EncodeToString(hashBytes[:])

	return &Deployment{
		ID:         uid.ID(),
		EndpointID: endpoint.ID,
		Data:       data,
		Hash:       hashstr,
		CreatedAt:  time.Now(),
	}
}
