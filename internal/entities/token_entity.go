package entities

import "time"

type Token struct {
	UserID    string `xorm:"not null VARCHAR(36) user_id"`
	Token     string `xorm:"not null unique TEXT token"`
	AAL       string `xorm:"not null string aal"`
	IP        string `xorm:"not null TEXT ip"`
	UserAgent string `xorm:"not null TEXT user_agent"`

	LastUsedAt *time.Time `xorm:"TIMESTAMPZ last_used_at"`

	CreatedAt time.Time  `xorm:"created TIMESTAMPZ created_at"`
	DeletedAt *time.Time `xorm:"TIMESTAMPZ deleted_at"`
}

func (Token) TableName() string {
	return "tokens"
}
