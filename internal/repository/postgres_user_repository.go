package repository

import (
	"context"
	"fmt"

	"github.com/Alias1177/Zloy/internal/config"
	"github.com/Alias1177/Zloy/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// postgresUserRepository реализует UserRepository для PostgreSQL
type postgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository создает новый экземпляр postgresUserRepository
func NewPostgresUserRepository(cfg *config.Config) (domain.UserRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.PostgresHost, cfg.Database.PostgresPort, cfg.Database.PostgresUser,
		cfg.Database.PostgresPassword, cfg.Database.PostgresDB)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &postgresUserRepository{db: db}, nil
}

// Create создает нового пользователя
func (r *postgresUserRepository) Create(ctx context.Context, login, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (login, password_hash, balance)
		VALUES ($1, $2, 0)
		RETURNING id, login, password_hash, balance, created_at`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, login, passwordHash).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByLogin получает пользователя по логину
func (r *postgresUserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
		SELECT id, login, password_hash, balance, created_at
		FROM users WHERE login = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID получает пользователя по ID
func (r *postgresUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, login, password_hash, balance, created_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateBalance обновляет баланс пользователя
func (r *postgresUserRepository) UpdateBalance(ctx context.Context, userID int, newBalance int) error {
	query := `UPDATE users SET balance = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newBalance, userID)
	return err
}

// GetBalance получает баланс пользователя
func (r *postgresUserRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	query := `SELECT balance FROM users WHERE id = $1`
	var balance int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	return balance, err
}

// Close закрывает соединение с базой данных
func (r *postgresUserRepository) Close() error {
	return r.db.Close()
}
