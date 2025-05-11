package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type TestRepository struct {
	db *gorm.DB
}

func NewTestRepository(db *gorm.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) CreateTest(ctx context.Context, test domain.RoadmapTest) (*domain.RoadmapTest, error) {
	result := r.db.Create(&test)
	if result.Error != nil {
		return nil, result.Error
	}
	return &test, result.Error

}
func (r *TestRepository) Test(ctx context.Context, testID uuid.UUID) (*domain.RoadmapTest, error) {
	var test domain.RoadmapTest
	err := r.db.WithContext(ctx).Where("id = ?", testID).First(&test).Error
	return &test, err
}

func (r *TestRepository) UpdateTest(ctx context.Context, test domain.RoadmapTest) (*domain.RoadmapTest, error) {
	result := r.db.WithContext(ctx).Model(&domain.RoadmapTest{}).
		Where("id = ?", test.ID).
		Omit("id").
		Updates(&test)

	if result.Error != nil {
		return nil, result.Error
	}
	return &test, nil
}

func (r *TestRepository) TestsForReport(ctx context.Context, historyID uuid.UUID) ([]*domain.TestResult, error) {
	var results []*domain.TestResult

	err := r.db.WithContext(ctx).
		Model(&domain.RoadmapTest{}). // Указываем исходную модель
		Select("results_json_b, passed_at").
		Where("roadmap_history_id = ? AND status = ?", historyID, "passed"). //TODO
		Scan(&results).                                                      // Сканируем в DTO
		Error

	return results, err
}
