package entities

import "time"

type Deployment struct {
	ID        string    `xorm:"not null pk VARCHAR(36) id"`
	CreatedAt time.Time `xorm:"created TIMESTAMP created_at"`
	DeletedAt time.Time `xorm:"TIMESTAMP deleted_at"`

	EndpointID string `xorm:"not null VARCHAR(36) endpoint_id"`
	Hash       string `xorm:"not null VARCHAR(32) hash"`
	Data       []byte `xorm:"not null MEDIUMBLOB data"`
}

func (Deployment) Table() string {
	return "deployments"
}
