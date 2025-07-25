package usecase

import (
	"github.com/Alias1177/Zloy/internal/domain"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

// authUseCase реализует бизнес-логику для аутентификации
type authUseCase struct {
	jwtSecret      []byte
	expirationTime time.Duration
}

// NewAuthUseCase создает новый экземпляр authUseCase
func NewAuthUseCase(jwtSecret []byte, expirationTime time.Duration) domain.AuthUseCase {
	return &authUseCase{
		jwtSecret:      jwtSecret,
		expirationTime: expirationTime,
	}
}

// GenerateToken генерирует JWT токен
func (uc *authUseCase) GenerateToken(userID int, login string) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(uc.expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}

// ValidateToken валидирует JWT токен
func (uc *authUseCase) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidCredentials
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, domain.ErrInvalidCredentials
	}

	return claims, nil
}
