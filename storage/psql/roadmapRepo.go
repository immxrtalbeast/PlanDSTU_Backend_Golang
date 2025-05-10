package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type RoadmapRepository struct {
	db *gorm.DB
}

func NewRoadmapRepository(db *gorm.DB) *RoadmapRepository {
	return &RoadmapRepository{db: db}
}

func (r *RoadmapRepository) History(ctx context.Context, userID uuid.UUID, disciplineID int) (*domain.RoadmapHistory, error) {
	var history domain.RoadmapHistory
	err := r.db.
		Preload("Tests", "status = ?", "pending").
		Where("user_id = ? AND discipline_id = ?", userID, disciplineID).First(&history).Error
	return &history, err
}

func (r *RoadmapRepository) CreateHistory(ctx context.Context, history domain.RoadmapHistory) (*domain.RoadmapHistory, error) {
	result := r.db.WithContext(ctx).Create(&history)
	return &history, result.Error
}

func (r *RoadmapRepository) HistoryByID(ctx context.Context, historyID uuid.UUID) (*domain.RoadmapHistory, error) {
	var history domain.RoadmapHistory
	err := r.db.Where("id = ?", historyID).First(&history).Error
	return &history, err
}

func (r *RoadmapRepository) UpdateHistory(ctx context.Context, history *domain.RoadmapHistory) error {
	result := r.db.WithContext(ctx).Model(&domain.RoadmapHistory{}).
		Where("id = ?", history.ID).
		Omit("id").
		Updates(&history)

	return result.Error
}
