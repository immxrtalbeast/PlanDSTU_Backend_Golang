package psql

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateReport(ctx context.Context, report domain.Report) error {
	err := r.db.WithContext(ctx).Create(&report).Error
	return err
}

func (r *ReportRepository) Report(ctx context.Context, reportID uuid.UUID) (*domain.Report, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).Where("id = ?", reportID).First(&report).Error
	return &report, err
}
func (r *ReportRepository) ReportsByDisciplineID(ctx context.Context, disciplineID int) ([]*domain.Report, error) {
	var results []*domain.Report
	err := r.db.WithContext(ctx).
		Model(&domain.Report{}). // Указываем исходную модель
		Where("discipline_id = ?", disciplineID).
		Scan(&results). // Сканируем в DTO
		Error

	return results, err
}
func (r *ReportRepository) ReportByUserAndDisciplineIDs(ctx context.Context, disciplineID int, userID uuid.UUID) (*domain.Report, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).Where("user_id = ? AND discipline_id = ?", userID, disciplineID).First(&report).Error
	return &report, err
}

func (r *ReportRepository) UpdateReport(ctx context.Context, report domain.Report) error {
	result := r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("id = ?", report.ID).
		Omit("id").
		Updates(&report)

	return result.Error
}

func (r *ReportRepository) ReportDisciplines(ctx context.Context) ([]domain.DisciplineResponse, error) {
	var disciplines []domain.DisciplineResponse
	err := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Select("DISTINCT discipline_title as name, discipline_id as id").
		Scan(&disciplines).
		Error

	return disciplines, err
}

func (r *ReportRepository) ReportGroups(ctx context.Context, disciplineName string) ([]string, error) {
	var groups []string
	err := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Where("discipline_title = ?", disciplineName).
		Distinct("group").
		Pluck("group", &groups).
		Error

	return groups, err
}

func (r *ReportRepository) ReportsByGroupAndDiscipline(ctx context.Context, disciplineName string, group string) ([]*domain.Report, *domain.TimelineStat, error) {
	var reports []*domain.Report
	err := r.db.WithContext(ctx).
		Where(`discipline_title = ? AND "group" = ?`, disciplineName, group).
		Find(&reports).
		Error
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}
	stats, err := r.GetTimelineStats(ctx, disciplineName, group)
	return reports, stats, err
}
func (r *ReportRepository) GetTimelineStats(ctx context.Context, discipline, group string) (*domain.TimelineStat, error) {
	// Получаем все отчеты из базы данных
	var reports []domain.Report
	err := r.db.WithContext(ctx).
		Where("discipline_title = ? AND \"group\" = ?", discipline, group).
		Find(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Инициализируем структуру для агрегации
	result := &domain.TimelineStat{
		MinScore: math.MaxFloat64,
		MaxScore: -math.MaxFloat64,
	}

	totalScore := 0.0
	totalBlocks := 0

	for _, report := range reports {
		var details struct {
			Report []struct {
				ResultsJSONB struct {
					Blocks []struct {
						Value float64 `json:"value"`
					} `json:"blocks"`
				} `json:"ResultsJSONB"`
			} `json:"report"`
		}

		if err := json.Unmarshal(report.DetailsJSONB, &details); err != nil {
			return nil, fmt.Errorf("failed to parse DetailsJSONB: %w", err)
		}

		// Обрабатываем все вложенные отчеты
		for _, rpt := range details.Report {
			result.ReportsCount++

			if len(rpt.ResultsJSONB.Blocks) == 0 {
				continue
			}

			// Обновляем общую статистику
			for _, block := range rpt.ResultsJSONB.Blocks {
				val := block.Value
				totalScore += val
				totalBlocks++

				if val < result.MinScore {
					result.MinScore = val
				}
				if val > result.MaxScore {
					result.MaxScore = val
				}
			}
		}
	}

	// Рассчитываем среднее значение
	if totalBlocks > 0 {
		result.AvgScore = math.Round((totalScore/float64(totalBlocks))*100) / 100
	} else {
		result.AvgScore = 0
		result.MinScore = 0
		result.MaxScore = 0
	}

	return result, nil
}
