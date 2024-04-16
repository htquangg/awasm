package entities

import "time"

type Endpoint struct {
	CreatedAt          time.Time  `xorm:"created TIMESTAMPZ created_at"`
	UpdatedAt          time.Time  `xorm:"updated TIMESTAMPZ updated_at"`
	DeletedAt          *time.Time `xorm:"TIMESTAMPZ deleted_at"`
	ID                 string     `xorm:"not null pk VARCHAR(36) id"`
	UserID             string     `xorm:"not null VARCHAR(36) user_id"`
	Name               string     `xorm:"not null VARCHAR(256) name"`
	Runtime            string     `xorm:"not null VARCHAR(64) runtime"`
	ActiveDeploymentID string     `xorm:"not null default '' VARCHAR(36)"`
}

func (Endpoint) TableName() string {
	return "endpoints"
}

func (e Endpoint) HasActiveDeploy() bool {
	return e.ActiveDeploymentID != ""
}
