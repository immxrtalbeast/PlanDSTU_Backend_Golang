package domain

import (
	"context"

	"github.com/google/uuid"
)

type LLMInteractor interface {
	SaveHistory(ctx context.Context, req SaveHistoryRequest, userID uuid.UUID) (uuid.UUID, error)
}
type LLMRepository interface {
	SaveHistory(ctx context.Context, history *History) (uuid.UUID, error)
}
