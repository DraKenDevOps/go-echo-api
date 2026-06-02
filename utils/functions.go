package utils

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

// HashPassword creates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a bcrypt-hashed password with a plaintext candidate.
// Returns nil on match, or an error on mismatch.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func FormatUptime(d time.Duration) string {
	seconds := int(d.Seconds())

	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", seconds)

	case seconds < 3600:
		return fmt.Sprintf("%dm", seconds/60)

	case seconds < 86400:
		return fmt.Sprintf("%dh %dm",
			seconds/3600,
			(seconds%3600)/60)

	default:
		return fmt.Sprintf("%dd %dh %dm",
			seconds/86400,
			(seconds%86400)/3600,
			(seconds%3600)/60)
	}
}
