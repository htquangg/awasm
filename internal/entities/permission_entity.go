package entities

type AuthMode int

const (
	AuthModeJwt AuthMode = iota << 1
	AuthModeApiKey
)
