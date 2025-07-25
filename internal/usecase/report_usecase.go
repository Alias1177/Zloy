package usecase

import (
	"context"

	"github.com/Alias1177/Zloy/internal/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// reportUseCase реализует бизнес-логику для отчетов
type reportUseCase struct {
	reportRepo      domain.ReportRepository
	userUC          domain.UserUseCase
	reportCostCents int
}

// NewReportUseCase создает новый экземпляр reportUseCase
func NewReportUseCase(reportRepo domain.ReportRepository, userUC domain.UserUseCase, reportCostCents int) domain.ReportUseCase {
	return &reportUseCase{
		reportRepo:      reportRepo,
		userUC:          userUC,
		reportCostCents: reportCostCents,
	}
}

// CreateReport создает новый отчет
func (uc *reportUseCase) CreateReport(ctx context.Context, clientGeneratedID string) (*domain.Report, error) {
	report := &domain.Report{
		ID:                primitive.NewObjectID().Hex(),
		ReportID:          uuid.New().String(),
		UserID:            nil, // Анонимный отчет
		ClientGeneratedID: clientGeneratedID,
		IsPurchased:       false,
		CreatedAt:         domain.Now(),
	}

	err := uc.reportRepo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetUserReports получает отчеты пользователя с пагинацией
func (uc *reportUseCase) GetUserReports(ctx context.Context, userID int, limit, offset int) ([]*domain.Report, int64, error) {
	return uc.reportRepo.GetByUserID(ctx, userID, limit, offset)
}

// LinkAnonymousReports привязывает анонимные отчеты к пользователю
func (uc *reportUseCase) LinkAnonymousReports(ctx context.Context, clientGeneratedID string, userID int) (int64, error) {
	return uc.reportRepo.UpdateUserID(ctx, clientGeneratedID, userID)
}

// PurchaseReport покупает отчет
func (uc *reportUseCase) PurchaseReport(ctx context.Context, reportID string, userID int) error {
	// Получаем баланс пользователя
	balance, err := uc.userUC.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем, что у пользователя достаточно средств
	if balance < uc.reportCostCents {
		return domain.ErrInsufficientBalance
	}

	// Получаем отчет
	report, err := uc.reportRepo.GetByReportID(ctx, reportID)
	if err != nil {
		return domain.ErrReportNotFound
	}

	if report.IsPurchased {
		return domain.ErrReportAlreadyPurchased
	}

	// Обновляем баланс пользователя
	newBalance := balance - uc.reportCostCents
	err = uc.userUC.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return err
	}

	// Обновляем флаг покупки отчета
	err = uc.reportRepo.UpdatePurchased(ctx, reportID, true)
	if err != nil {
		// В случае ошибки пытаемся откатить баланс
		uc.userUC.UpdateBalance(ctx, userID, balance)
		return err
	}

	return nil
}
