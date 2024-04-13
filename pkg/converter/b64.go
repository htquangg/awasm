package converter

import "encoding/base64"

func FromB64(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ToB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func FromURLB64(s string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ToURLB64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}
