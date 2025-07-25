package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alias1177/Zloy/internal/delivery/middleware"
	"github.com/Alias1177/Zloy/internal/domain"
)

// UserHandler обрабатывает HTTP запросы для пользователей
type UserHandler struct {
	userUseCase domain.UserUseCase
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// TopUpBalance пополняет баланс пользователя
func (h *UserHandler) TopUpBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем сумму из query параметров
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		http.Error(w, "Amount parameter is required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	newBalance, err := h.userUseCase.TopUpBalance(r.Context(), userID, amount)
	if err != nil {
		if err.Error() == "amount must be positive" {
			http.Error(w, "Amount must be positive", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to top up balance", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Balance topped up successfully",
		"new_balance": newBalance,
		"added":       amount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBalance получает текущий баланс пользователя
func (h *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.userUseCase.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
