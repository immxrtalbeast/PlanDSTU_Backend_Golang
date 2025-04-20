package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	HistoryID uuid.UUID
	Role      string // "user" или "assistant"
	Content   string
	Timestamp time.Time
}
