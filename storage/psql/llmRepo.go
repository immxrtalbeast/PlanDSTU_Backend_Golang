package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LLMRepository struct {
	db *gorm.DB
}

func NewLLMRepository(db *gorm.DB) *LLMRepository {
	return &LLMRepository{db: db}
}

func (r *LLMRepository) SaveHistory(ctx context.Context, history *domain.History) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"messages_json_b", "created_at"}),
		}).
		Create(history)

	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return history.ID, nil
}
