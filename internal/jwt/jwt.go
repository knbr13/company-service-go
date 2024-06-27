package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a new JWT token with the given claims and expiration time.
func GenerateToken(claims jwt.Claims, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
