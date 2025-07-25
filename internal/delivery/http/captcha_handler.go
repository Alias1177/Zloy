package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Alias1177/Zloy/internal/domain"
)

// CaptchaHandler обрабатывает HTTP запросы для CAPTCHA
type CaptchaHandler struct {
	captchaUC domain.CaptchaUseCase
}

// NewCaptchaHandler создает новый экземпляр CaptchaHandler
func NewCaptchaHandler(captchaUC domain.CaptchaUseCase) *CaptchaHandler {
	return &CaptchaHandler{
		captchaUC: captchaUC,
	}
}

// GenerateCaptcha генерирует CAPTCHA
func (h *CaptchaHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID, imageData, err := h.captchaUC.GenerateCaptcha()
	if err != nil {
		http.Error(w, "Failed to generate CAPTCHA", http.StatusInternalServerError)
		return
	}

	// Кодируем изображение в base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	response := map[string]interface{}{
		"session_id": sessionID,
		"image":      base64Image,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
