package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware acts as an interceptor to validate JWT access tokens in incoming requests.
type AuthMiddleware struct {
	authConfig *AuthenticationConfig
}

// InitAuthMiddleware initializes an AuthMiddleware with the given AuthenticationConfig.
func InitAuthMiddleware(cfg *AuthenticationConfig) *AuthMiddleware {
	return &AuthMiddleware{
		authConfig: cfg,
	}
}

// ValidateAccessToken checks the signature, expiration, and issuer of a JWT access token.
// Returns an error if the token is invalid, expired, or issued by an unexpected source.
func (am *AuthMiddleware) ValidateAccessToken(accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.authConfig.jwtSecret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid access token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid or expired access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	if iss, ok := claims["iss"].(string); !ok || iss != issuer {
		return fmt.Errorf("invalid issuer")
	}

	return nil
}
