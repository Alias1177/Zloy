package domain

import (
	"context"
	"errors"
	"time"
)

// User представляет сущность пользователя
type User struct {
	ID           int       `json:"id" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      int       `json:"balance" db:"balance"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Validate проверяет валидность пользователя
func (u *User) Validate() error {
	if u.Login == "" {
		return errors.New("login is required")
	}
	if len(u.Login) < 3 {
		return errors.New("login must be at least 3 characters")
	}
	return nil
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, login, passwordHash string) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	UpdateBalance(ctx context.Context, userID int, newBalance int) error
	GetBalance(ctx context.Context, userID int) (int, error)
}

// UserUseCase определяет бизнес-логику для работы с пользователями
type UserUseCase interface {
	Register(ctx context.Context, login, password string) (*User, string, error)
	Login(ctx context.Context, login, password string) (*User, string, error)
	GetBalance(ctx context.Context, userID int) (int, error)
	UpdateBalance(ctx context.Context, userID int, newBalance int) error
	TopUpBalance(ctx context.Context, userID int, amount int) (int, error)
}
