package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TimelineStat struct {
	AvgScore     float64 `json:"avg_score"`
	MinScore     float64 `json:"min_score"`
	MaxScore     float64 `json:"max_score"`
	ReportsCount int     `json:"reports_count"`
}

type Report struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DisciplineTitle string         `gorm:"not null"`
	DisciplineID    int            `gorm:"not null"`
	Group           string         `gorm:"not null"`
	UserID          uuid.UUID      `gorm:"type:uuid;index"`
	DetailsJSONB    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt       time.Time
}

type DisciplineResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ReportInteractor interface {
	CreateReport(ctx context.Context, discplineID int, resultsJSONB datatypes.JSON, userID uuid.UUID, disciplineTitle string, group string) error
	Report(ctx context.Context, reportID uuid.UUID) (*Report, error)
	ReportsByDisciplineID(ctx context.Context, disciplineID int) ([]*Report, error)
	ReportDisciplines(ctx context.Context) ([]DisciplineResponse, error)
	ReportGroups(ctx context.Context, disciplineName string) ([]string, error)
	ReportsByGroupAndDiscipline(ctx context.Context, disciplineName, group string) ([]*Report, *TimelineStat, error)
}
type ReportRepository interface {
	CreateReport(ctx context.Context, report Report) error
	UpdateReport(ctx context.Context, report Report) error
	Report(ctx context.Context, reportID uuid.UUID) (*Report, error)
	ReportsByDisciplineID(ctx context.Context, disciplineID int) ([]*Report, error)
	ReportByUserAndDisciplineIDs(ctx context.Context, disciplineID int, userID uuid.UUID) (*Report, error)
	ReportDisciplines(ctx context.Context) ([]DisciplineResponse, error)
	ReportGroups(ctx context.Context, disciplineName string) ([]string, error)
	ReportsByGroupAndDiscipline(ctx context.Context, disciplineName string, group string) ([]*Report, *TimelineStat, error)
}
