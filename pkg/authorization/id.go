package authorization

import (
	"math/rand"
	"time"
)

// NOTABLE EXCEPTIONS: a, c, p for Accounts, Comments, and Posts, respectively
const idChars = "bdefghijklmnoqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateID creates an id with n characters
func GenerateID(prefix string, length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	idBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		idBytes[i] = idChars[rand.Intn(len(idChars))]
	}

	return prefix + string(idBytes)
}
