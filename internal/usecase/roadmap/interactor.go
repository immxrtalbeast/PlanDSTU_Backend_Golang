package roadmap

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type RoadmapInteractor struct {
	roadmapRepo domain.RoadmapRepository
	testsRepo   domain.TestRepository
}

func NewRoadmapInteractor(roadmapRepo domain.RoadmapRepository, testsRepo domain.TestRepository) *RoadmapInteractor {
	return &RoadmapInteractor{roadmapRepo: roadmapRepo, testsRepo: testsRepo}
}

func (ri *RoadmapInteractor) History(ctx context.Context, userID uuid.UUID, disciplineID int) (*domain.RoadmapHistory, error) {
	const op = "uc.roadmap.history"
	history, err := ri.roadmapRepo.History(ctx, userID, disciplineID)
	return history, err

}

func (ri *RoadmapInteractor) CreateHistory(ctx context.Context, userID uuid.UUID, discplineID int) (*domain.RoadmapHistory, error) {
	const op = "uc.roadmap.create_history"
	_, err := ri.roadmapRepo.History(ctx, userID, discplineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			history := domain.RoadmapHistory{
				UserID:       userID,
				DisciplineID: discplineID,
			}
			result_history, err := ri.roadmapRepo.CreateHistory(ctx, history)
			return result_history, err
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	} else {
		return nil, fmt.Errorf("%s", "History already exists.")
	}
}

func (ri *RoadmapInteractor) Report(ctx context.Context, userID uuid.UUID, disciplineID int) ([]*domain.TestResult, error) {
	const op = "uc.roadmap.report"
	history, err := ri.roadmapRepo.History(ctx, userID, disciplineID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	tests, err := ri.testsRepo.TestsForReport(ctx, history.ID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	// TODO
	return tests, nil
}
