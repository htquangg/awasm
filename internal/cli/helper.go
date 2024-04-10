package cli

import "github.com/htquangg/a-wasm/internal/schemas"

func ConfigContainsEmail(users []*schemas.LoggedInUser, email string) bool {
	for _, value := range users {
		if value.Email == email {
			return true
		}
	}
	return false
}
