package utils

import (
	"log"
	"time"

	"go-echo-api/config"
	"go-echo-api/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID     int                   `json:"user_id"`
	Username   string                `json:"username"`
	Telephone  string                `json:"telephone"`
	Email      string                `json:"email"`
	Level      models.UserLevel      `json:"level"`
	RoleAction models.UserRoleAction `json:"role_action"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the user
func GenerateToken(user models.JwtUser, cfg *config.Config) (string, error) {

	payload := JWTClaims{
		UserID:     user.UserID,
		Username:   user.Username,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Level:      user.Level,
		RoleAction: user.RoleAction,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	if cfg.JWTPrivateKey != "" {
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWTPrivateKey))
		if err != nil {
			log.Fatalf("Error parses a PEM encoded PKCS1 or PKCS8 private key: %s", err.Error())
			return "", err
		}
		token := jwt.NewWithClaims(jwt.SigningMethodPS256, payload)
		return token.SignedString(key)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte("JACKPOT!"))
}

// ValidateToken validates the JWT token and returns the claims
func ValidateToken(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if cfg.JWTPublicKey != "" {
			if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); ok {
				return jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.JWTPublicKey))
			}
		}
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte("JACKPOT!"), nil
		}
		return nil, jwt.ErrTokenInvalidClaims
	})
	if err != nil {
		log.Fatalf("Error ParseWithClaims: %s", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
