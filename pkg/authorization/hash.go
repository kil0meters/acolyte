package authorization

import "golang.org/x/crypto/bcrypt"

// HashString hashes a string
func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyHash verifies a hash
func VerifyHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
