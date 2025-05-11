package report

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReportInteractor struct {
	reportRepo domain.ReportRepository
}

func NewReportInteractor(reportRepo domain.ReportRepository) *ReportInteractor {
	return &ReportInteractor{reportRepo: reportRepo}
}

func (ri *ReportInteractor) CreateReport(ctx context.Context, discplineID int, resultsJSONB datatypes.JSON, userID uuid.UUID, disciplineTitle string, group string) error {
	const op = "uc.report.create"
	existingReport, err := ri.reportRepo.ReportByUserAndDisciplineIDs(ctx, discplineID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			report := domain.Report{
				DisciplineTitle: disciplineTitle,
				DisciplineID:    discplineID,
				DetailsJSONB:    resultsJSONB,
				UserID:          userID,
				Group:           group,
			}
			if err := ri.reportRepo.CreateReport(ctx, report); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		} else {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		existingReport.DetailsJSONB = resultsJSONB
		if err := ri.reportRepo.UpdateReport(ctx, *existingReport); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	}
	return nil
}

func (ri *ReportInteractor) Report(ctx context.Context, reportID uuid.UUID) (*domain.Report, error) {
	const op = "uc.report.get"
	report, err := ri.reportRepo.Report(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return report, nil
}

func (ri *ReportInteractor) ReportsByDisciplineID(ctx context.Context, disciplineID int) ([]*domain.Report, error) {
	const op = "uc.report.all.by.DisciplineID"
	reports, err := ri.reportRepo.ReportsByDisciplineID(ctx, disciplineID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return reports, nil
}

func (ri *ReportInteractor) ReportDisciplines(ctx context.Context) ([]domain.DisciplineResponse, error) {
	const op = "uc.report.disciplines"
	disciplines, err := ri.reportRepo.ReportDisciplines(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return disciplines, nil
}

func (ri *ReportInteractor) ReportGroups(ctx context.Context, disciplineName string) ([]string, error) {
	const op = "uc.report.groups"
	groups, err := ri.reportRepo.ReportGroups(ctx, disciplineName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return groups, nil
}

func (ri *ReportInteractor) ReportsByGroupAndDiscipline(ctx context.Context, disciplineName, group string) ([]*domain.Report, *domain.TimelineStat, error) {
	const op = "uc.report.reportsByGroup"
	reports, stats, err := ri.reportRepo.ReportsByGroupAndDiscipline(ctx, disciplineName, group)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	return reports, stats, nil
}
