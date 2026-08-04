package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// passwordHashCost is above bcrypt.DefaultCost (10), which is below what OWASP
// currently recommends. Existing hashes keep working: the cost is embedded in them.
const passwordHashCost = 12

// DummyHash is a valid bcrypt hash of an unguessable value. It is compared against
// when the address is unknown, so a login attempt for a non-existent account costs
// the same time as one for a real account and cannot be told apart by timing.
const DummyHash = "$2a$12$fcdBNm3OAC91c8tcAIhsAOQmatoVR2oQgmJfeJ5tpgMlX.5c8GBP2"

func HashAndSalt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
