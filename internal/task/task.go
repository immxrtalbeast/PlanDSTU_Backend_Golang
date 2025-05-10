package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const QueueGenerateTest = "generate_test"

type GenerateTestPayload struct {
	TestID        string    `json:"test_id"`
	Themes        []string  `json:"themes"`
	UserID        string    `json:"user_id"`
	DisciplineID  int       `json:"discipline_id"`
	HistoryID     uuid.UUID `json:"history_id"`
	LLMServiceURL string    `json:"llm_service_url"`
}

var RedisClient *asynq.Client

func Init(redisAddr string) {
	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

func NewGenerateTestTask(payload GenerateTestPayload) (*asynq.Task, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}
	return asynq.NewTask(QueueGenerateTest, payloadJSON, asynq.Retention(1*time.Minute), asynq.MaxRetry(10)), nil
}
