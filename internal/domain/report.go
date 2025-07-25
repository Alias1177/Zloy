package domain

import (
	"context"
	"errors"
	"time"
)

// Report представляет сущность отчета
type Report struct {
	ID                string    `json:"_id" bson:"_id"`
	ReportID          string    `json:"report_id" bson:"report_id"`
	UserID            *int      `json:"user_id" bson:"user_id"`
	ClientGeneratedID string    `json:"client_generated_id" bson:"client_generated_id"`
	IsPurchased       bool      `json:"is_purchased" bson:"is_purchased"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
}

// Validate проверяет валидность отчета
func (r *Report) Validate() error {
	if r.ClientGeneratedID == "" {
		return errors.New("client_generated_id is required")
	}
	if r.ReportID == "" {
		return errors.New("report_id is required")
	}
	return nil
}

// ReportRepository определяет интерфейс для работы с отчетами
type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*Report, int64, error)
	GetByReportID(ctx context.Context, reportID string) (*Report, error)
	GetByClientGeneratedID(ctx context.Context, clientGeneratedID string) ([]*Report, error)
	UpdateUserID(ctx context.Context, clientGeneratedID string, userID int) (int64, error)
	UpdatePurchased(ctx context.Context, reportID string, isPurchased bool) error
}

// ReportUseCase определяет бизнес-логику для работы с отчетами
type ReportUseCase interface {
	CreateReport(ctx context.Context, clientGeneratedID string) (*Report, error)
	GetUserReports(ctx context.Context, userID int, limit, offset int) ([]*Report, int64, error)
	LinkAnonymousReports(ctx context.Context, clientGeneratedID string, userID int) (int64, error)
	PurchaseReport(ctx context.Context, reportID string, userID int) error
}
