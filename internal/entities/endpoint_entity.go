package entities

import "time"

type Endpoint struct {
	ID                 string `xorm:"not null pk VARCHAR(36) id"`
	Name               string `xorm:"not null VARCHAR(256) name"`
	Runtime            string `xorm:"not null VARCHAR(64) runtime"`
	ActiveDeploymentID string `xorm:"not null default '' VARCHAR(36)"`

	CreatedAt time.Time  `xorm:"created TIMESTAMPZ created_at"`
	UpdatedAt time.Time  `xorm:"updated TIMESTAMPZ updated_at"`
	DeletedAt *time.Time `xorm:"TIMESTAMPZ deleted_at"`
}

func (Endpoint) TableName() string {
	return "endpoints"
}

func (e Endpoint) HasActiveDeploy() bool {
	return e.ActiveDeploymentID != ""
}
