package encryptor

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

// Encrypt takes a message and a key, encrypts the message using the key,
// and returns the encrypted message as a base64-encoded string.
// It generates a random nonce and uses it along with the key to encrypt the message.
// The encrypted message is then returned along with a nil error if encryption is successful.
func Encrypt(message []byte, key [32]byte) (string, error) {
	var nonce [24]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return "", fmt.Errorf("encrypt: nonce: %w", err)
	}
	encrypted := secretbox.Seal(nonce[:], message, &nonce, &key)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts the given encrypted string using the provided key.
// It returns the decrypted byte slice or an error if decryption fails.
func Decrypt(encrypted string, key [32]byte) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt: base64 decode: %w", err)
	}
	var nonce [24]byte
	copy(nonce[:], encryptedBytes[:24])
	decrypted, ok := secretbox.Open(nil, encryptedBytes[24:], &nonce, &key)
	if !ok {
		return nil, fmt.Errorf("decrypt: failed to open")
	}

	return decrypted, nil
}
