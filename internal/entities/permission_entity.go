package entities

type AuthMode int

const (
	JWT AuthMode = iota << 1
	API_KEY
)
