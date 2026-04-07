package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Options struct {
	JWTSecret string
}

func AuthMiddleware(opt Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		println(authHeader)
		println("cek =", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid header"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(opt.JWTSecret), nil
		})

		if err != nil {
			println("ERROR:", err.Error())
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}