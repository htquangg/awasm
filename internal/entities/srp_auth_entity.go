package entities

import "time"

type (
	SrpAuth struct {
		CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`
		UpdatedAt time.Time `xorm:"updated TIMESTAMPZ updated_at"`
		ID        string    `xorm:"not null pk VARCHAR(36) id"`
		UserID    string    `xorm:"not null unique VARCHAR(36) user_id"`
		SrpUserID string    `xorm:"not null unique VARCHAR(36) srp_user_id"`
		Salt      string    `xorm:"not null TEXT salt"`
		Verifier  string    `xorm:"not null TEXT verifier"`
	}

	SrpAuthTemp struct {
		CreatedAt      time.Time `xorm:"created TIMESTAMPZ created_at"`
		ID             string    `xorm:"not null pk VARCHAR(36) id"`
		UserID         string    `xorm:"not null VARCHAR(36) user_id"`
		SrpUserID      string    `xorm:"not null VARCHAR(36) srp_user_id"`
		Salt           string    `xorm:"not null TEXT salt"`
		Verifier       string    `xorm:"not null TEXT verifier"`
		SrpChallengeID string    `xorm:"not null VARCHAR(36) srp_challenge_id"`
	}

	SrpChallenge struct {
		CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`

		VerifiedAt *time.Time `xorm:"TIMESTAMPZ verified_at"`

		ID           string `xorm:"not null pk VARCHAR(36) id"`
		SrpUserID    string `xorm:"not null VARCHAR(36) srp_user_id"`
		ServerKey    string `xorm:"not null TEXT server_key"`
		SrpA         string `xorm:"not null TEXT srp_a"`
		AttemptCount int32  `xorm:"not null INT default 0 attempt_count"`
	}
)

func (SrpAuth) TableName() string {
	return "srp_auth"
}

func (SrpAuthTemp) TableName() string {
	return "srp_auth_temp"
}

func (SrpChallenge) TableName() string {
	return "srp_challenges"
}
