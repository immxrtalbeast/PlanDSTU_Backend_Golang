package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"github.com/immxrtalbeast/plandstu/internal/task"
	"gorm.io/datatypes"
)

type Worker struct {
	server  *asynq.Server
	testINT domain.TestInteractor
}

func NewWorker(redisAddr string, concurrency int, testINT domain.TestInteractor) *Worker {
	return &Worker{
		server: asynq.NewServer(
			asynq.RedisClientOpt{Addr: redisAddr},
			asynq.Config{
				Concurrency: concurrency,
			},
		),
		testINT: testINT,
	}
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	w.registerHandlers(mux)
	return w.server.Run(mux)
}
func (w *Worker) handleGenerateTestTask(ctx context.Context, t *asynq.Task) error {
	var payload task.GenerateTestPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	client := &http.Client{Timeout: 60 * time.Minute}
	reqBody := map[string]interface{}{
		"test_id": payload.TestID,
		"themes":  payload.Themes,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest("POST", payload.LLMServiceURL+"api/test-workflow", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "LLM/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("llm request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("llm returned status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	testID, err := uuid.Parse(payload.TestID)
	if err != nil {
		return fmt.Errorf("invalid test ID format: %w", err)
	}
	_, err = w.testINT.CreateTest(ctx, testID, datatypes.JSON(data), payload.HistoryID, false)
	if err != nil {
		return fmt.Errorf("failed to save test: %v", err)
	}

	return nil
}
func (w *Worker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return w.handleGenerateTestTask(ctx, t)
}
func (w *Worker) registerHandlers(mux *asynq.ServeMux) {
	mux.Handle(
		task.QueueGenerateTest,
		asynq.HandlerFunc(w.ProcessTask),
	)
}
