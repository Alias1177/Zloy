package repository

import (
	"context"

	"github.com/Alias1177/Zloy/internal/config"
	"github.com/Alias1177/Zloy/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoReportRepository реализует ReportRepository для MongoDB
type mongoReportRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoReportRepository создает новый экземпляр mongoReportRepository
func NewMongoReportRepository(cfg *config.Config) (domain.ReportRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Database.MongoURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("zl0y_billing")
	collection := db.Collection("reports")

	return &mongoReportRepository{
		db:         db,
		collection: collection,
	}, nil
}

// Create создает новый отчет
func (r *mongoReportRepository) Create(ctx context.Context, report *domain.Report) error {
	_, err := r.collection.InsertOne(ctx, report)
	return err
}

// GetByUserID получает отчеты пользователя с пагинацией
func (r *mongoReportRepository) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*domain.Report, int64, error) {
	filter := bson.M{"user_id": userID}

	// Подсчитываем общее количество
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Получаем отчеты с пагинацией
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []*domain.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// GetByReportID получает отчет по report_id
func (r *mongoReportRepository) GetByReportID(ctx context.Context, reportID string) (*domain.Report, error) {
	filter := bson.M{"report_id": reportID}

	var report domain.Report
	err := r.collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetByClientGeneratedID получает отчеты по client_generated_id
func (r *mongoReportRepository) GetByClientGeneratedID(ctx context.Context, clientGeneratedID string) ([]*domain.Report, error) {
	filter := bson.M{"client_generated_id": clientGeneratedID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []*domain.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

// UpdateUserID обновляет user_id для отчетов с указанным client_generated_id
func (r *mongoReportRepository) UpdateUserID(ctx context.Context, clientGeneratedID string, userID int) (int64, error) {
	filter := bson.M{
		"client_generated_id": clientGeneratedID,
		"user_id":             nil,
	}

	update := bson.M{
		"$set": bson.M{
			"user_id": userID,
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// UpdatePurchased обновляет флаг is_purchased для отчета
func (r *mongoReportRepository) UpdatePurchased(ctx context.Context, reportID string, isPurchased bool) error {
	filter := bson.M{"report_id": reportID}
	update := bson.M{
		"$set": bson.M{
			"is_purchased": isPurchased,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
