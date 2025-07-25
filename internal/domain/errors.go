package domain

import (
	"errors"
	"time"
)

// Ошибки домена
var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInsufficientBalance    = errors.New("insufficient balance")
	ErrReportNotFound         = errors.New("report not found")
	ErrReportAlreadyPurchased = errors.New("report already purchased")
	ErrInvalidCaptcha         = errors.New("invalid captcha")
)

// Now возвращает текущее время
func Now() time.Time {
	return time.Now()
}
