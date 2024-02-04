package uid

import "github.com/google/uuid"

func ID() string {
	return uuid.Must(uuid.NewV7()).String()
}
