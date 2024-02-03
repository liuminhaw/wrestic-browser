package encryptor

import (
	"encoding/base64"
	"fmt"
)

// UrlDecodeKey decodes a URL-encoded key string and returns a fixed-size byte array.
// The key string is expected to be encoded using base64.URLEncoding.
// If decoding fails, an error is returned.
func UrlDecodeKey(key string) ([32]byte, error) {
	encKeyBytes, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return [32]byte{}, fmt.Errorf("url decode key: %w", err)
	}

	return [32]byte(encKeyBytes), nil
}
