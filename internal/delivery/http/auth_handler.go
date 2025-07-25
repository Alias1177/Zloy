package http

import (
	"encoding/json"
	"net/http"

	"github.com/Alias1177/Zloy/internal/domain"
)

// AuthHandler обрабатывает HTTP запросы для аутентификации
type AuthHandler struct {
	userUseCase domain.UserUseCase
	captchaUC   domain.CaptchaUseCase
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(userUseCase domain.UserUseCase, captchaUC domain.CaptchaUseCase) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		captchaUC:   captchaUC,
	}
}

// Register обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем CAPTCHA
	if !h.captchaUC.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		http.Error(w, "Invalid CAPTCHA", http.StatusBadRequest)
		return
	}

	user, token, err := h.userUseCase.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			http.Error(w, "User already exists", http.StatusConflict)
		default:
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	response := AuthResponse{
		AccessToken: token,
		User: map[string]interface{}{
			"id":    user.ID,
			"login": user.Login,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, token, err := h.userUseCase.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "Failed to login", http.StatusInternalServerError)
		}
		return
	}

	response := AuthResponse{
		AccessToken: token,
		User: map[string]interface{}{
			"id":    user.ID,
			"login": user.Login,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
