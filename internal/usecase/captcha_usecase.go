package usecase

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/Alias1177/Zloy/internal/domain"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// captchaUseCase реализует бизнес-логику для CAPTCHA
type captchaUseCase struct {
	sessions        map[string]domain.CaptchaSession
	mutex           sync.RWMutex
	width           int
	height          int
	noiseCount      int
	sessionLifetime time.Duration
}

// NewCaptchaUseCase создает новый экземпляр captchaUseCase
func NewCaptchaUseCase(width, height, noiseCount int, sessionLifetime time.Duration) domain.CaptchaUseCase {
	uc := &captchaUseCase{
		sessions:        make(map[string]domain.CaptchaSession),
		width:           width,
		height:          height,
		noiseCount:      noiseCount,
		sessionLifetime: sessionLifetime,
	}

	// Запускаем очистку старых сессий
	go uc.cleanupSessions()
	return uc
}

// GenerateCaptcha генерирует CAPTCHA изображение
func (uc *captchaUseCase) GenerateCaptcha() (string, []byte, error) {
	// Генерируем случайное 4-значное число
	num, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return "", nil, err
	}
	answer := strconv.FormatInt(num.Int64()+1000, 10)

	// Создаем изображение
	img := image.NewRGBA(image.Rect(0, 0, uc.width, uc.height))

	// Заполняем фон
	for y := 0; y < uc.height; y++ {
		for x := 0; x < uc.width; x++ {
			img.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}

	// Добавляем шум
	for i := 0; i < uc.noiseCount; i++ {
		x := int(num.Int64()) % uc.width
		y := int(num.Int64()) % uc.height
		img.Set(x, y, color.RGBA{200, 200, 200, 255})
	}

	// Рисуем цифры
	point := fixed.Point26_6{X: fixed.Int26_6(20 * 64), Y: fixed.Int26_6(50 * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{50, 50, 50, 255}),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(answer)

	// Кодируем в PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", nil, err
	}

	// Создаем сессию CAPTCHA
	sessionID := fmt.Sprintf("captcha_%d", time.Now().UnixNano())
	session := domain.CaptchaSession{
		ID:        sessionID,
		Answer:    answer,
		CreatedAt: domain.Now(),
	}

	uc.mutex.Lock()
	uc.sessions[sessionID] = session
	uc.mutex.Unlock()

	return sessionID, buf.Bytes(), nil
}

// VerifyCaptcha проверяет CAPTCHA
func (uc *captchaUseCase) VerifyCaptcha(sessionID, userAnswer string) bool {
	uc.mutex.RLock()
	session, exists := uc.sessions[sessionID]
	uc.mutex.RUnlock()

	if !exists {
		return false
	}

	// Проверяем время жизни
	if time.Since(session.CreatedAt) > uc.sessionLifetime {
		uc.mutex.Lock()
		delete(uc.sessions, sessionID)
		uc.mutex.Unlock()
		return false
	}

	// Проверяем ответ
	if session.Answer == userAnswer {
		uc.mutex.Lock()
		delete(uc.sessions, sessionID)
		uc.mutex.Unlock()
		return true
	}

	return false
}

// cleanupSessions очищает старые сессии
func (uc *captchaUseCase) cleanupSessions() {
	ticker := time.NewTicker(uc.sessionLifetime)
	defer ticker.Stop()

	for range ticker.C {
		uc.mutex.Lock()
		now := domain.Now()
		for sessionID, session := range uc.sessions {
			if now.Sub(session.CreatedAt) > uc.sessionLifetime {
				delete(uc.sessions, sessionID)
			}
		}
		uc.mutex.Unlock()
	}
}
