package authorization

import (
	"math/rand"
	"time"
)

const idChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateID creates an id with n characters
func GenerateID(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	idBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		idBytes[i] = idChars[rand.Intn(len(idChars))]
	}

	return string(idBytes)
}
