package http

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LinkAnonymousRequest представляет запрос на привязку анонимных отчетов
type LinkAnonymousRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

// CreateReportRequest представляет запрос на создание отчета
type CreateReportRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

// ReportsResponse представляет ответ с отчетами
type ReportsResponse struct {
	Reports []interface{} `json:"reports"`
	Total   int64         `json:"total"`
}

// AuthResponse представляет ответ с токеном
type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
