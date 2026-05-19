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

// GenerateAccessToken รับ expire เป็น parameter
// ไม่ hardcode ค่าใดๆ ไว้ใน function นี้
func GenerateAccessToken(
	userID string,
	email string,
	secret string,
	expire time.Duration, // ✅ รับจากภายนอก
) (string, error) {

	claims := JWTClaim{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(secret))
}
