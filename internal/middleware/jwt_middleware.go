package middleware

import (
	"auth-service/internal/config"
	"auth-service/internal/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})

			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		claims := &token.JWTClaim{}

		tkn, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			},
		)

		if err != nil || !tkn.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})

			c.Abort()

			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
