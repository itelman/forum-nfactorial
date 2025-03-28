package encoder

import (
	"encoding/base64"
)

func EncodeAccessToken(token string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(token))
	return encoded
}

func DecodeAccessToken(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
