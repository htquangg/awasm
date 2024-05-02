package entities

import "time"

type ApiKey struct {
	CreatedAt    time.Time  `xorm:"created TIMESTAMPZ created_at"`
	LastUsedAt   *time.Time `xorm:"TIMESTAMPZ last_used_at"`
	ID           string     `xorm:"not null pk VARCHAR(36) id"`
	UserID       string     `xorm:"not null VARCHAR(36) user_id"`
	Key          string     `xorm:"not null VARCHAR(512) key"`
	KeyPreview   string     `xorm:"not null VARCHAR(32) key_preview"`
	FriendlyName string     `xorm:"not null TEXT default '' friendly_name"`
}

func (ApiKey) TableName() string {
	return "api_keys"
}
