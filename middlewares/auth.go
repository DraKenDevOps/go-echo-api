package middleware

import (
	"strings"

	"go-echo-api/config"
	"go-echo-api/utils"
	"go-echo-api/zplogger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type contextKey string

const ClaimsKey contextKey = "user_claims"

// AuthChecker returns an Echo middleware that validates JWT from
// the Authorization (Bearer) header or X-Access-Token header.
func AuthChecker(cfg *config.Config, logger *zplogger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := ExtractToken(c)
			if token == "" {
				return c.JSON(200, map[string]string{
					"status":  "error",
					"message": "Missing authorization token",
				})
			}

			claims, err := utils.ValidateToken(token, cfg)
			if err != nil {
				logger.Error("Error on validate token", zap.Error(err))
				return c.JSON(200, map[string]string{
					"status":  "error",
					"message": "Invalid or expired token",
				})
			}

			c.Set(string(ClaimsKey), claims)
			return next(c)
		}
	}
}

// GetClaims retrieves the JWTClaims stored in the Echo context.
func GetClaims(c echo.Context) *utils.JWTClaims {
	if claims, ok := c.Get(string(ClaimsKey)).(*utils.JWTClaims); ok {
		return claims
	}
	return nil
}

// extractToken pulls the token string from the Authorization header
// (Bearer scheme) or the X-Access-Token header.
func ExtractToken(c echo.Context) string {
	// Try Authorization header
	auth := c.Request().Header.Get("Authorization")
	if auth != "" {
		const prefix = "Bearer "
		if strings.HasPrefix(auth, prefix) {
			return strings.TrimSpace(auth[len(prefix):])
		}
	}

	// Fallback to X-Access-Token header
	return strings.TrimSpace(c.Request().Header.Get("X-Access-Token"))
}
