package entities

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	TokenLength = 32
)

type ClaimScope string

const (
	UnknownTokenScope ClaimScope = "unknown-token"
	AccessTokenScope  ClaimScope = "access-token"
	RefreshTokenScope ClaimScope = "refresh-token"
	SignupTokenScope  ClaimScope = "signup-token"
)

func (c ClaimScope) Ptr() *ClaimScope {
	return &c
}

type RefreshToken struct {
	ID        string `xorm:"not null pk VARCHAR(36) id"`
	UserID    string `xorm:"not null VARCHAR(36) user_id"`
	Token     string `xorm:"not null VARCHAR(255) token"`
	SessionID string `xorm:"not null VARCHAR(36) session_id"`
	Revoked   bool   `xorm:"not null BOOL default false revoked"`

	CreatedAt time.Time `xorm:"created TIMESTAMPZ created_at"`
	UpdatedAt time.Time `xorm:"updated TIMESTAMPZ updated_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type GrantParams struct {
	FactorID string

	IP        string
	UserAgent string
}

type CommonTokenClaims struct {
	jwt.StandardClaims
	Email string      `json:"email"`
	Scope *ClaimScope `json:"scope"`
}

func (w *CommonTokenClaims) GetScope() ClaimScope {
	if w.Scope == nil {
		return UnknownTokenScope
	}
	return *w.Scope
}

type AccessTokenClaims struct {
	CommonTokenClaims
	AuthenticatorAssuranceLevel   string     `json:"aal,omitempty"`
	AuthenticationMethodReference []AMREntry `json:"amr,omitempty"`
	SessionID                     string     `json:"session_id,omitempty"`
}
