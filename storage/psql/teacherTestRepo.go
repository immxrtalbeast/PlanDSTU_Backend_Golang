package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type TeacherTestRepository struct {
	db *gorm.DB
}

func NewTeacherTestRepository(db *gorm.DB) *TeacherTestRepository {
	return &TeacherTestRepository{db: db}
}

func (r *TeacherTestRepository) TeacherTests(ctx context.Context, disciplineID int) ([]*domain.TeacherTest, error) {
	var tests []*domain.TeacherTest
	err := r.db.WithContext(ctx).Model(&domain.TeacherTest{}).Where("discipline_id = ?", disciplineID).Scan(&tests).Error
	return tests, err

}
func (r *TeacherTestRepository) TeacherTestByID(ctx context.Context, testID uuid.UUID) (*domain.TeacherTest, error) {
	var test domain.TeacherTest
	err := r.db.WithContext(ctx).Where("id = ?", testID).First(&test).Error
	return &test, err

}
func (r *TeacherTestRepository) TeacherTestWithoutAnswers(ctx context.Context, disciplineID int) ([]*domain.TeacherTest, error) {
	var tests []*domain.TeacherTest
	err := r.db.WithContext(ctx).
		Model(&domain.TeacherTest{}).
		Where("discipline_id = ?", disciplineID).
		Omit("Answers").
		Scan(&tests).
		Error
	return tests, err

}

func (r *TeacherTestRepository) CreateTeacherTest(ctx context.Context, test domain.TeacherTest) error {
	err := r.db.WithContext(ctx).Create(&test).Error
	return err
}

// UpdateTeacherTest(ctx context.Context, test TeacherTest) error
// DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error

func (r *TeacherTestRepository) UpdateTeacherTest(ctx context.Context, test domain.TeacherTest) error {
	result := r.db.WithContext(ctx).Model(&domain.TeacherTest{}).
		Where("id = ?", test.ID).
		Omit("id").
		Updates(&test)

	return result.Error
}
func (r *TeacherTestRepository) DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", testID).Delete(&domain.TeacherTest{}).Error
}
