package config

import (
	"os"
	"strconv"
	"time"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Captcha  CaptchaConfig
	Business BusinessConfig
}

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port            string
	DefaultPageSize int
}

// DatabaseConfig содержит настройки баз данных
type DatabaseConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MongoURI         string
}

// JWTConfig содержит настройки JWT
type JWTConfig struct {
	Secret         []byte
	ExpirationTime time.Duration
}

// CaptchaConfig содержит настройки CAPTCHA
type CaptchaConfig struct {
	Width           int
	Height          int
	NoiseCount      int
	SessionLifetime time.Duration
}

// BusinessConfig содержит бизнес-настройки
type BusinessConfig struct {
	ReportCostCents int
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_LIMIT", 10),
		},
		Database: DatabaseConfig{
			PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
			PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
			PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
			PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
			PostgresDB:       getEnv("POSTGRES_DB", "zl0y_billing"),
			MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		},
		JWT: JWTConfig{
			Secret:         []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production")),
			ExpirationTime: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
		Captcha: CaptchaConfig{
			Width:           getEnvAsInt("CAPTCHA_WIDTH", 200),
			Height:          getEnvAsInt("CAPTCHA_HEIGHT", 80),
			NoiseCount:      getEnvAsInt("CAPTCHA_NOISE_COUNT", 1000),
			SessionLifetime: time.Duration(getEnvAsInt("CAPTCHA_LIFETIME_MINUTES", 5)) * time.Minute,
		},
		Business: BusinessConfig{
			ReportCostCents: getEnvAsInt("REPORT_COST_CENTS", 1000),
		},
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как integer или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
