package entities

import "time"

type (
	SrpAuth struct {
		ID        string `xorm:"not null pk VARCHAR(36) id"`
		UserID    string `xorm:"not null unique VARCHAR(36) user_id"`
		SrpUserID string `xorm:"not null unique VARCHAR(36) srp_user_id"`
		Salt      string `xorm:"not null TEXT salt"`
		Verifier  string `xorm:"not null TEXT verifier"`

		CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`
		UpdatedAt time.Time `xorm:"updated TIMESTAMPZ updated_at"`
	}

	SrpAuthTemp struct {
		ID        string `xorm:"not null pk VARCHAR(36) id"`
		UserID    string `xorm:"not null VARCHAR(36) user_id"`
		SrpUserID string `xorm:"not null VARCHAR(36) srp_user_id"`
		Salt      string `xorm:"not null TEXT salt"`
		Verifier  string `xorm:"not null TEXT verifier"`

		CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`
	}

	SrpChallenge struct {
		ID            string `xorm:"not null pk VARCHAR(36) id"`
		SrpAuthTempID string `xorm:"not null VARCHAR(36) srp_auth_temp_id"`
		SrpUserID     string `xorm:"not null VARCHAR(36) srp_user_id"`
		ServerKey     string `xorm:"not null TEXT server_key"`
		SrpA          string `xorm:"not null TEXT srp_a"`
		AttemptCount  int32  `xorm:"not null INT default 0 attempt_count"`

		VerifiedAt time.Time `xorm:"TIMESTAMPZ verified_at"`

		CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`
		UpdatedAt time.Time `xorm:"updated TIMESTAMPZ updated_at"`
	}
)
