package http

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var PrivateKey *rsa.PrivateKey

func LoadPrivateKey() error {
	key := os.Getenv("PRIVATE_KEY")
	if key == "" {
		log.Fatal("PRIVATE_KEY not set")
	}

	// Replace literal "\n" with actual newlines
	key = strings.ReplaceAll(key, `\n`, "\n")

	fmt.Println("Loaded private key successfully!")
	var err error
	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	return err
}

// TokenRequest represents the request body for token generation
type TokenRequest struct {
	Username string `json:"username" example:"harrison"`
	Password string `json:"password" example:"password123"`
}

// TokenResponse represents the JSON response with JWT
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int64  `json:"expires_in" example:"299"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Bad Request"`
}

// generateToken godoc
// @Summary Generate a mock JWT token
// @Description Generates a mock JWT token for a given username and password. For testing purposes only.
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body TokenRequest true "Username and password"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users/mock-authoriz-user [post]
func generateToken(c *fiber.Ctx) error {
	var req TokenRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// MOCK AUTH — DO NOT USE IN PROD
	if req.Username != "harrison" || req.Password != "password123" {
		return fiber.ErrUnauthorized
	}

	expiration := time.Now().Add(5 * time.Minute)

	claims := jwt.MapClaims{
		"iss": "mock-auth-service",
		"sub": req.Username,
		"aud": "internal-tools",
		"exp": expiration.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(PrivateKey)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(TokenResponse{
		AccessToken: signedToken,
		ExpiresIn:   int64(time.Until(expiration).Seconds()),
	})
}

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(app *fiber.App) {
	userRoutes := app.Group("/v1/users")
	userRoutes.Post("/mock-authoriz-user", generateToken)
}
