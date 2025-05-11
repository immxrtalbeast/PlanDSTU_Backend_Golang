package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TeacherTest struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DisciplineID int            `gorm:"not null"`
	DetailsJSONB datatypes.JSON `gorm:"type:jsonb"`
	Answers      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time
}

// DTO для API ответа
type TestResponse struct {
	ID   uuid.UUID   `json:"id"`
	Test []TestBlock `json:"test"`
}

type TestBlock struct {
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

type Option struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}
type TeacherTestInteractor interface {
	CreateTeacherTest(ctx context.Context, detailsData datatypes.JSON, answers datatypes.JSON, disciplineID int) error
	UpdateTeacherTest(ctx context.Context, testID uuid.UUID, detailsData datatypes.JSON, answers datatypes.JSON) error
	DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error
	TeacherTests(ctx context.Context, disciplineID int) ([]*TeacherTest, error)
	TeacherTestByID(ctx context.Context, testID uuid.UUID) (*TeacherTest, error)
	TeacherTestForUser(ctx context.Context, disciplineID int) (*TestResponse, error)
}
type TeacherTestRepository interface {
	CreateTeacherTest(ctx context.Context, test TeacherTest) error
	UpdateTeacherTest(ctx context.Context, test TeacherTest) error
	DeleteTeacherTest(ctx context.Context, testID uuid.UUID) error
	TeacherTests(ctx context.Context, disciplineID int) ([]*TeacherTest, error)
	TeacherTestByID(ctx context.Context, testID uuid.UUID) (*TeacherTest, error)
	TeacherTestWithoutAnswers(ctx context.Context, disciplineID int) ([]*TeacherTest, error)
}
