package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/datatypes"
)

type TestInteractor struct {
	testRepo    domain.TestRepository
	llmURL      string
	roadmapRepo domain.RoadmapRepository
}

func NewTestInteractor(testRepo domain.TestRepository, llmURL string, roadmapRepo domain.RoadmapRepository) *TestInteractor {
	return &TestInteractor{testRepo: testRepo, llmURL: llmURL, roadmapRepo: roadmapRepo}
}

func (ti *TestInteractor) CreateTest(ctx context.Context, generatedTestID uuid.UUID, detailsData datatypes.JSON, roadmapHistoryID uuid.UUID, isFirst bool) (*domain.RoadmapTest, error) {
	const op = "uc.tests.create"
	test := domain.RoadmapTest{
		ID:               generatedTestID,
		DetailsJSONB:     detailsData,
		RoadmapHistoryID: roadmapHistoryID,
	}
	if isFirst {
		test.IsFirst = true
	}
	history, err := ti.testRepo.CreateTest(ctx, test)
	return history, err
}

func (ti *TestInteractor) Answers(ctx context.Context, testID uuid.UUID, answers []string) ([]byte, error) { //map[string]float64
	const op = "uc.tests.answers"
	test, err := ti.testRepo.Test(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	var testDetails domain.TestDetails
	if err := json.Unmarshal(test.DetailsJSONB, &testDetails); err != nil {
		return nil, fmt.Errorf("%s: failed to parse test details: %w", op, err)
	}

	correctAnswers, err := ti.GetCorrectAnswers(ctx, testID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	totalQuestions := 0
	for _, topic := range testDetails.Test {
		totalQuestions += len(topic.Questions)
	}

	if len(answers) != totalQuestions || len(correctAnswers) != totalQuestions {
		return nil, fmt.Errorf("%s: invalid answers count. Answers count %d, Answers given: %d", op, totalQuestions, len(answers))
	}

	results := make(map[string]float64)
	answerIdx := 0

	for _, topic := range testDetails.Test {
		correct := 0
		topicTotal := len(topic.Questions)

		for range topic.Questions {
			if answers[answerIdx] == correctAnswers[answerIdx] {
				correct++
			}
			answerIdx++
		}

		if topicTotal > 0 {
			results[topic.Title] = float64(correct) / float64(topicTotal) * 100
		}
	}
	if test.IsFirst {
		err = ti.SaveAnswersToHistory(ctx, test.RoadmapHistoryID, results)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	var blocksData domain.BlocksData

	// Добавляем новые блоки
	for topic, score := range results {
		blocksData.Blocks = append(blocksData.Blocks, domain.Block{
			Name:  topic,
			Value: math.Round(score*100) / 100, // Округление до 2 знаков
		})
	}

	// Сериализуем и сохраняем
	updatedJSON, _ := json.Marshal(blocksData)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal results: %w", op, err)
	}

	test.PassedAt = time.Now()
	test.Status = "passed"
	test.ResultsJSONB = datatypes.JSON(updatedJSON)
	_, err = ti.testRepo.UpdateTest(ctx, *test)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return updatedJSON, nil
}

func (ti *TestInteractor) GetCorrectAnswers(ctx context.Context, testID uuid.UUID) ([]string, error) {
	type CorrectAnswersResponse struct {
		Answers []string `json:"answers"`
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	url := fmt.Sprintf(ti.llmURL+"test-exmpl-answers/%s", testID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Parser/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error code while getting answers: %s", errorBody)
	}
	var response CorrectAnswersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return response.Answers, nil
}

func (ti *TestInteractor) SaveAnswersToHistory(ctx context.Context, historyID uuid.UUID, results map[string]float64) error {
	// 1. Получаем историю
	history, err := ti.roadmapRepo.HistoryByID(ctx, historyID)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	// 2. Парсим существующие данные
	var blocksData domain.BlocksData

	// Добавляем новые блоки
	for topic, score := range results {
		blocksData.Blocks = append(blocksData.Blocks, domain.Block{
			Name:  topic,
			Value: math.Round(score*100) / 100, // Округление до 2 знаков
		})
	}

	// Сериализуем и сохраняем
	updatedJSON, _ := json.Marshal(blocksData)
	history.BlocksJSONB = updatedJSON

	return ti.roadmapRepo.UpdateHistory(ctx, history)
}
