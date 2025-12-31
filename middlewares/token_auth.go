package middlewares

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var PublicKey *rsa.PublicKey

func LoadPublicKey() error {
	keyData, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return err
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(keyData)
	return err
}
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Expecting: "Bearer <token>"
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		// Parse and verify
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure token uses RS256
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return PublicKey, nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		// Attach claims to context if needed
		c.Locals("jwt_claims", token.Claims)
		return c.Next()
	}
}
