package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RoadmapTest struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoadmapHistoryID uuid.UUID      `gorm:"type:uuid;index"` // Внешний ключ
	Status           string         `gorm:"size:50;default:'pending'"`
	DetailsJSONB     datatypes.JSON `gorm:"type:jsonb"`
	ResultsJSONB     datatypes.JSON `gorm:"type:jsonb"`
	IsFirst          bool           `gorm:"default:false"`
	CreatedAt        time.Time
	PassedAt         time.Time
}
type TestDetails struct {
	Test []struct {
		Title     string `json:"title"`
		Questions []struct {
			Text    string   `json:"text"`
			Options []Option `json:"options"`
		} `json:"questions"`
	} `json:"test"`
}

type TestResult struct {
	ResultsJSONB datatypes.JSON `gorm:"column:results_json_b"`
	PassedAt     time.Time      `gorm:"column:passed_at"`
}
type TestInteractor interface {
	CreateTest(ctx context.Context, generatedTestID uuid.UUID, detailsData datatypes.JSON, roadmapHistoryID uuid.UUID, isFirst bool) (*RoadmapTest, error)
	Answers(ctx context.Context, testID uuid.UUID, answers []string) ([]byte, error)
	GetCorrectAnswers(ctx context.Context, testID uuid.UUID) ([]string, error)
}

type TestRepository interface {
	CreateTest(ctx context.Context, test RoadmapTest) (*RoadmapTest, error)
	Test(ctx context.Context, testID uuid.UUID) (*RoadmapTest, error)
	UpdateTest(ctx context.Context, test RoadmapTest) (*RoadmapTest, error)
	TestsForReport(ctx context.Context, historyID uuid.UUID) ([]*TestResult, error)
}
