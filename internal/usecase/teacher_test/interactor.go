package teachertest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/datatypes"
)

type TeacherTestInteractor struct {
	teacherTestRepo domain.TeacherTestRepository
}

func NewTeacherTestInteractor(teacherTestRepo domain.TeacherTestRepository) *TeacherTestInteractor {
	return &TeacherTestInteractor{teacherTestRepo: teacherTestRepo}
}

func (ti *TeacherTestInteractor) TeacherTests(ctx context.Context, disciplineID int) ([]*domain.TeacherTest, error) {
	const op = "uc.teacher_test.get"
	tests, err := ti.teacherTestRepo.TeacherTests(ctx, disciplineID)
	return tests, err
}
func (ti *TeacherTestInteractor) TeacherTestByID(ctx context.Context, testID uuid.UUID) (*domain.TeacherTest, error) {
	const op = "uc.teacher_test.get_id"
	test, err := ti.teacherTestRepo.TeacherTestByID(ctx, testID)
	return test, err
}
func (ti *TeacherTestInteractor) CreateTeacherTest(ctx context.Context, detailsData datatypes.JSON, answers datatypes.JSON, disciplineID int) error {
	const op = "uc.teacher_test.create"
	test := domain.TeacherTest{
		DisciplineID: disciplineID,
		DetailsJSONB: detailsData,
		Answers:      answers,
	}
	err := ti.teacherTestRepo.CreateTeacherTest(ctx, test)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (ti *TeacherTestInteractor) TeacherTestForUser(ctx context.Context, disciplineID int) (*domain.TestResponse, error) {
	const op = "uc.teacher_test.for_user"
	tests, err := ti.teacherTestRepo.TeacherTestWithoutAnswers(ctx, disciplineID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests available")
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := rand.Intn(len(tests))
	randomTest := tests[randomIndex]
	var details []domain.TestBlock
	if err := json.Unmarshal(randomTest.DetailsJSONB, &details); err != nil {
		return nil, fmt.Errorf("%s: failed to unmarshal details: %w", op, err)
	}

	// Формируем финальный ответ
	response := &domain.TestResponse{
		ID:   randomTest.ID,
		Test: details,
	}

	return response, nil
}

// UpdateTeacherTest(ctx context.Context, testID uuid.UUID, detailsData datatypes.JSON, answers datatypes.JSON) error
// DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error
func (ti *TeacherTestInteractor) UpdateTeacherTest(ctx context.Context, testID uuid.UUID, detailsData datatypes.JSON, answers datatypes.JSON) error {
	const op = "uc.teacher_test.update"
	existingTest, err := ti.teacherTestRepo.TeacherTestByID(ctx, testID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	existingTest.Answers = answers
	existingTest.DetailsJSONB = detailsData
	err = ti.teacherTestRepo.UpdateTeacherTest(ctx, *existingTest)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (ti *TeacherTestInteractor) DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error {
	const op = "uc.teacher_test.delele"
	if err := ti.teacherTestRepo.DeleteTeacherTest(ctx, testID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
