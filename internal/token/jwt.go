package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

func GenerateAccessToken(
	userID string,
	email string,
	secret string,
) (string, error) {

	claims := JWTClaim{
		UserID: userID,
		Email:  email,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),

			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}
