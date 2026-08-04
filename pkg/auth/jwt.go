package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mstgnz/cronjob/pkg/config"
)

var letterRunes = []rune("0987654321abcçdefgğhıijklmnoöpqrsştuüvwxyzABCÇDEFGĞHIİJKLMNOÖPQRSTUÜVWXYZ-_!?+&%=*")

// TokenTTL is how long an issued token stays valid. The auth cookie is given the
// same lifetime so the browser never holds a cookie that outlives its token.
const TokenTTL = 24 * time.Hour

// GenerateToken token generate
func GenerateToken(userId int) (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
		Issuer:    strconv.Itoa(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.App().SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// ValidateToken token validate
func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(config.App().SecretKey), nil
	})
}

func GetUserIDByToken(token string) (string, error) {
	id, _, err := GetUserIDAndIssuedAt(token)
	return id, err
}

// GetUserIDAndIssuedAt returns the subject and the moment the token was issued.
// The issue time is what makes logout and password changes able to retire a token
// that is otherwise still within its expiry.
func GetUserIDAndIssuedAt(token string) (string, time.Time, error) {
	valid, err := ValidateToken(token)
	if err != nil {
		return "", time.Time{}, err
	}
	claims, ok := valid.Claims.(jwt.MapClaims)
	if !ok {
		return "", time.Time{}, fmt.Errorf("unexpected claims type")
	}
	id := fmt.Sprintf("%v", claims["iss"])

	issuedAt, err := claims.GetIssuedAt()
	if err != nil || issuedAt == nil {
		return "", time.Time{}, fmt.Errorf("token has no issue time")
	}
	return id, issuedAt.Time, nil
}

func RandomString(length int) string {
	s, r := make([]rune, length), []rune(letterRunes)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

func RandomHex(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return RandomString(length)
	}
	return hex.EncodeToString(bytes)
}
