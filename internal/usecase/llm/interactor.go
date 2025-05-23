package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
)

type LLMInteractor struct {
	userRepo domain.LLMRepository
}

func NewLLMInteractor(userRepo domain.LLMRepository) *LLMInteractor {
	return &LLMInteractor{userRepo: userRepo}
}

func (li *LLMInteractor) SaveHistory(ctx context.Context, req domain.SaveHistoryRequest, userID uuid.UUID) (uuid.UUID, error) {
	const op = "uc.llm.save"
	messagesJSON, err := json.Marshal(req.History)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	history := domain.History{
		UserID:        userID,
		CreatedAt:     time.Now(),
		MessagesJSONB: messagesJSON,
	}
	id, err := li.userRepo.SaveHistory(ctx, &history)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
