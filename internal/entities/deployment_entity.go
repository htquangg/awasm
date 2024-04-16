package entities

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/htquangg/a-wasm/pkg/uid"
)

type Deployment struct {
	CreatedAt  time.Time  `xorm:"created TIMESTAMPZ created_at"`
	DeletedAt  *time.Time `xorm:"TIMESTAMPZ deleted_at"`
	ID         string     `xorm:"not null pk VARCHAR(36) id"`
	UserID     string     `xorm:"not null VARCHAR(36) user_id"`
	EndpointID string     `xorm:"not null VARCHAR(36) endpoint_id"`
	Hash       string     `xorm:"not null VARCHAR(32) hash"`
	Data       []byte     `xorm:"not null MEDIUMBLOB data"`
}

func (Deployment) TableName() string {
	return "deployments"
}

func NewFromEndpoint(endpoint *Endpoint, userID string, data []byte) *Deployment {
	hashBytes := md5.Sum(data)
	hashstr := hex.EncodeToString(hashBytes[:])

	return &Deployment{
		ID:         uid.ID(),
		UserID:     userID,
		EndpointID: endpoint.ID,
		Data:       data,
		Hash:       hashstr,
		CreatedAt:  time.Now(),
	}
}
