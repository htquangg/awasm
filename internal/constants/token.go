package constants

import "time"

const (
	ExpiresInTokenSignup    = 15 * 60 * time.Second
	ExpiresInOTPEmailSignup = 60 * time.Second
)
