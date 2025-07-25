package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims представляет JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

// AuthUseCase определяет бизнес-логику для аутентификации
type AuthUseCase interface {
	GenerateToken(userID int, login string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// CaptchaSession представляет сессию CAPTCHA
type CaptchaSession struct {
	ID        string    `json:"id"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

// CaptchaUseCase определяет бизнес-логику для CAPTCHA
type CaptchaUseCase interface {
	GenerateCaptcha() (string, []byte, error)
	VerifyCaptcha(sessionID, userAnswer string) bool
}
