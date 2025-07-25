package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alias1177/Zloy/internal/delivery/middleware"
	"github.com/Alias1177/Zloy/internal/domain"

	"github.com/go-chi/chi/v5"
)

// ReportHandler обрабатывает HTTP запросы для отчетов
type ReportHandler struct {
	reportUseCase   domain.ReportUseCase
	userUseCase     domain.UserUseCase
	defaultPageSize int
}

// NewReportHandler создает новый экземпляр ReportHandler
func NewReportHandler(reportUseCase domain.ReportUseCase, userUseCase domain.UserUseCase, defaultPageSize int) *ReportHandler {
	return &ReportHandler{
		reportUseCase:   reportUseCase,
		userUseCase:     userUseCase,
		defaultPageSize: defaultPageSize,
	}
}

// LinkAnonymous привязывает анонимные отчеты к пользователю
func (h *ReportHandler) LinkAnonymous(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LinkAnonymousRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	updated, err := h.reportUseCase.LinkAnonymousReports(r.Context(), req.ClientGeneratedID, userID)
	if err != nil {
		http.Error(w, "Failed to link reports", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Reports linked successfully",
		"updated": updated,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserReports получает отчеты пользователя
func (h *ReportHandler) GetUserReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем параметры пагинации
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := h.defaultPageSize // по умолчанию
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	reports, total, err := h.reportUseCase.GetUserReports(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get reports", http.StatusInternalServerError)
		return
	}

	response := ReportsResponse{
		Reports: make([]interface{}, len(reports)),
		Total:   total,
	}

	for i, report := range reports {
		response.Reports[i] = report
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PurchaseReport покупает отчет
func (h *ReportHandler) PurchaseReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reportID := chi.URLParam(r, "report_id")
	if reportID == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	err := h.reportUseCase.PurchaseReport(r.Context(), reportID, userID)
	if err != nil {
		switch err {
		case domain.ErrInsufficientBalance:
			http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		case domain.ErrReportNotFound:
			http.Error(w, "Report not found", http.StatusNotFound)
		case domain.ErrReportAlreadyPurchased:
			http.Error(w, "Report already purchased", http.StatusConflict)
		default:
			http.Error(w, "Failed to purchase report", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Report purchased successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateMockReport создает mock отчет
func (h *ReportHandler) CreateMockReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	report, err := h.reportUseCase.CreateReport(r.Context(), req.ClientGeneratedID)
	if err != nil {
		http.Error(w, "Failed to create report", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":   "Mock report created",
		"report_id": report.ReportID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
