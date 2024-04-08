package user

import (
	"encoding/base64"
)

func convertStringToBytes(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func convertBytesToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
