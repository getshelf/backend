package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func newSessionToken() (raw string, hash [32]byte, err error) {
	value := make([]byte, 32)

	if _, err := rand.Read(value); err != nil {
		return "", [32]byte{}, err
	}

	raw = base64.RawURLEncoding.EncodeToString(value)
	hash = sha256.Sum256([]byte(raw))

	return raw, hash, nil
}

func hashSessionToken(raw string) [32]byte {
	return sha256.Sum256([]byte(raw))
}
