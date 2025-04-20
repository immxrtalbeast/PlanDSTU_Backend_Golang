package domain

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID // внешний ключ
	CreatedAt time.Time
	Messages  []Message `gorm:"foreignKey:HistoryID"` // 1 ко многим
}
