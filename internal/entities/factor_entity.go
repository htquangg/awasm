package entities

import "time"

type FactorState int

const (
	FactorStateUnverified FactorState = iota
	FactorStateVerified
)

func (factorState FactorState) String() string {
	switch factorState {
	case FactorStateUnverified:
		return "unverified"
	case FactorStateVerified:
		return "verified"
	}
	return ""
}

const TOTP = "totp"

type AuthenticationMethod int

const (
	OAuth AuthenticationMethod = iota
	PasswordGrant
	OTP
	TOTPSignIn
	Recovery
	Invite
	MagicLink
	EmailSignup
	TokenRefresh
	Anonymous
)

func (authMethod AuthenticationMethod) String() string {
	switch authMethod {
	case OAuth:
		return "oauth"
	case PasswordGrant:
		return "password"
	case OTP:
		return "otp"
	case TOTPSignIn:
		return "totp"
	case Recovery:
		return "recovery"
	case Invite:
		return "invite"
	case MagicLink:
		return "magiclink"
	case EmailSignup:
		return "email/signup"
	case TokenRefresh:
		return "token_refresh"
	case Anonymous:
		return "anonymous"
	}
	return ""
}

type MFAFactor struct {
	CreatedAt    time.Time  `xorm:"created TIMESTAMPZ created_at"`
	UpdatedAt    time.Time  `xorm:"updated TIMESTAMPZ updated_at"`
	LastUsedAt   *time.Time `xorm:"TIMESTAMPZ last_used_at"`
	ID           string     `xorm:"not null pk VARCHAR(36) id"`
	UserID       string     `xorm:"not null VARCHAR(36) user_id"`
	Status       string     `xorm:"not null TEXT status"`
	FriendlyName string     `xorm:"not null TEXT default '' friendly_name"`
	FactorType   string     `xorm:"not null TEXT factor_type"` // totp, webauthn
	Secret       string     `xorm:"not null TEXT default '' secret"`
}

func (MFAFactor) TableName() string {
	return "mfa_factors"
}

type MFAAMRClaim struct {
	CreatedAt            time.Time `xorm:"created TIMESTAMPZ created_at"`
	UpdatedAt            time.Time `xorm:"updated TIMESTAMPZ updated_at"`
	SessionID            string    `xorm:"not null VARCHAR(36) session_id"`
	AuthenticationMethod string    `xorm:"not null authentication_method"`
}

func (MFAAMRClaim) TableName() string {
	return "mfa_amr_claims"
}

func (a *MFAAMRClaim) GetAuthenticationMethod() string {
	return a.AuthenticationMethod
}

type MFAChallenge struct {
	CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`

	VerifiedAt *time.Time `xorm:"TIMESTAMPZ verified_at"`

	ID       string `xorm:"not null pk VARCHAR(36) id"`
	FactorID string `xorm:"not null VARCHAR(36) facor_id"`
	IP       string `xorm:"not null default '' TEXT ip"`
}
