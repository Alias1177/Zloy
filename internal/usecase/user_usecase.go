package usecase

import (
	"context"
	"errors"

	"github.com/Alias1177/Zloy/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// userUseCase реализует бизнес-логику для пользователей
type userUseCase struct {
	userRepo domain.UserRepository
	authUC   domain.AuthUseCase
}

// NewUserUseCase создает новый экземпляр userUseCase
func NewUserUseCase(userRepo domain.UserRepository, authUC domain.AuthUseCase) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		authUC:   authUC,
	}
}

// Register регистрирует нового пользователя
func (uc *userUseCase) Register(ctx context.Context, login, password string) (*domain.User, string, error) {
	// Проверяем, существует ли пользователь
	existingUser, _ := uc.userRepo.GetByLogin(ctx, login)
	if existingUser != nil {
		return nil, "", domain.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Создаем пользователя
	user, err := uc.userRepo.Create(ctx, login, string(hashedPassword))
	if err != nil {
		return nil, "", err
	}

	// Генерируем токен
	token, err := uc.authUC.GenerateToken(user.ID, user.Login)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login выполняет вход пользователя
func (uc *userUseCase) Login(ctx context.Context, login, password string) (*domain.User, string, error) {
	// Получаем пользователя
	user, err := uc.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	// Генерируем токен
	token, err := uc.authUC.GenerateToken(user.ID, user.Login)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetBalance получает баланс пользователя
func (uc *userUseCase) GetBalance(ctx context.Context, userID int) (int, error) {
	return uc.userRepo.GetBalance(ctx, userID)
}

// UpdateBalance обновляет баланс пользователя
func (uc *userUseCase) UpdateBalance(ctx context.Context, userID int, newBalance int) error {
	return uc.userRepo.UpdateBalance(ctx, userID, newBalance)
}

// TopUpBalance пополняет баланс пользователя
func (uc *userUseCase) TopUpBalance(ctx context.Context, userID int, amount int) (int, error) {
	// Проверяем, что сумма положительная
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}

	// Получаем текущий баланс
	currentBalance, err := uc.userRepo.GetBalance(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Обновляем баланс
	newBalance := currentBalance + amount
	err = uc.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}
