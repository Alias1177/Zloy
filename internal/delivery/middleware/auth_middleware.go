package middleware

import (
	"context"
	"github.com/Alias1177/Zloy/internal/domain"
	"net/http"
	"strings"
)

// AuthMiddleware создает middleware для аутентификации
func AuthMiddleware(authUC domain.AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Проверяем формат "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Валидируем токен
			claims, err := authUC.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем claims в контекст
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "login", claims.Login)

			// Передаем запрос дальше с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext получает user_id из контекста
func GetUserIDFromContext(r *http.Request) int {
	if userID, ok := r.Context().Value("user_id").(int); ok {
		return userID
	}
	return 0
}

// GetLoginFromContext получает login из контекста
func GetLoginFromContext(r *http.Request) string {
	if login, ok := r.Context().Value("login").(string); ok {
		return login
	}
	return ""
}
