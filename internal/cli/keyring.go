package cli

import (
	"github.com/zalando/go-keyring"
)

const MainKeyringService = "awasm-cli"

func SetValueInKeyring(key, value string) error {
	return keyring.Set(MainKeyringService, key, value)
}

func GetValueInKeyring(key string) (string, error) {
	return keyring.Get(MainKeyringService, key)
}

func DeleteValueInKeyring(key string) error {
	return keyring.Delete(MainKeyringService, key)
}
