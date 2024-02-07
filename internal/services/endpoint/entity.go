package endpoint

import "time"

type EndpointEntity struct {
	ID        string    `xorm:"not null pk VARCHAR(36) id"`
	CreatedAt time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt time.Time `xorm:"updated TIMESTAMP updated_at"`
	DeletedAt time.Time `xorm:"TIMESTAMP deleted_at"`

	Name               string `xorm:"not null VARCHAR(50) name"`
	Runtime            string `xorm:"not null VARCHAR(50) runtime"`
	ActiveDeploymentID string `xorm:"not null default '' VARCHAR(36)"`
}

func (EndpointEntity) TableName() string {
	return "endpoints"
}
