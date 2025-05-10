package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type History struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID      `gorm:"uniqueIndex"` // Внешний ключ (можно добавить индекс)
	MessagesJSONB datatypes.JSON `gorm:"type:jsonb"`  // Весь чат как JSONB
	CreatedAt     time.Time
}

type SaveHistoryRequest struct {
	History []struct {
		Content          string         `json:"content"`
		Type             string         `json:"type"`
		AdditionalKwargs map[string]any `json:"additional_kwargs"`
		ResponseMetadata map[string]any `json:"response_metadata"`
		Example          bool           `json:"example"`
		Timestamp        time.Time      `json:"timestamp"`
	} `json:"history"`
	UserID string `json:"user_id"`
}
