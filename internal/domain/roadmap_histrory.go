package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RoadmapHistory struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DisciplineID int            `gorm:"not null"`
	UserID       uuid.UUID      `gorm:"type:uuid;index"`
	BlocksJSONB  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time
	Tests        []RoadmapTest `gorm:"foreignKey:RoadmapHistoryID"` // Связь с тестами

}

// Структура для работы с BlocksJSONB

type BlocksData struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
type RoadmapInteractor interface {
	History(ctx context.Context, userID uuid.UUID, disciplineID int) (*RoadmapHistory, error)
	// CreateHistoryWithFirstTest(ctx context.Context, userID uuid.UUID, discplineID int, firstTest HistoryTest) error
	CreateHistory(ctx context.Context, userID uuid.UUID, discplineID int) (*RoadmapHistory, error)
	Report(ctx context.Context, userID uuid.UUID, disciplineID int) ([]*TestResult, error)
}

type RoadmapRepository interface {
	History(ctx context.Context, userID uuid.UUID, disciplineID int) (*RoadmapHistory, error)
	HistoryByID(ctx context.Context, historyID uuid.UUID) (*RoadmapHistory, error)
	CreateHistory(ctx context.Context, history RoadmapHistory) (*RoadmapHistory, error)
	UpdateHistory(ctx context.Context, history *RoadmapHistory) error
}
